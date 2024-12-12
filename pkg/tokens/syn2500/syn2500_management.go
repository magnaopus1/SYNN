package syn2500

import (
	"time"
	"errors"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rsa"
	"crypto/rand"
)

// SYN2500Management handles DAO membership, voting, and governance for SYN2500 tokens.
type SYN2500Management struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewSYN2500Management initializes the DAO management system with RSA keys.
func NewSYN2500Management() (*SYN2500Management, error) {
	// Generate encryption keys
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &SYN2500Management{
		privateKey: privKey,
		publicKey:  &privKey.PublicKey,
	}, nil
}

// AddMember adds a new member to the DAO and issues a SYN2500 token to them.
func (manager *SYN2500Management) AddMember(member common.DAOMember, daoID string) (*common.SYN2500Token, error) {
	// Create a new SYN2500 token for the member
	token := &common.SYN2500Token{
		TokenID:     generateUniqueID(),
		Owner:       member.MemberID,
		DAOID:       daoID,
		VotingPower: member.InitialVotingPower,
		IssuedDate:  time.Now(),
		Active:      true,
	}

	// Register the token with the ledger using Synnergy Consensus
	err := ledger.StoreDAOToken(token, synconsensus.SubBlockValidation)
	if err != nil {
		return nil, err
	}

	// Add token to the member's record
	member.TokenID = token.TokenID
	member.Status = "Active"

	// Return the generated token
	return token, nil
}

// RemoveMember removes a member from the DAO and deactivates their SYN2500 token.
func (manager *SYN2500Management) RemoveMember(memberID string, daoID string) error {
	// Retrieve the member's token from the ledger
	token, err := ledger.GetDAOTokenByMemberID(memberID, daoID)
	if err != nil {
		return err
	}

	// Deactivate the token
	token.Active = false

	// Update the ledger
	err = ledger.UpdateDAOToken(token, synconsensus.SubBlockValidation)
	if err != nil {
		return err
	}

	// Remove the member from the DAO
	err = ledger.RemoveMemberFromDAO(memberID, daoID, synconsensus.SubBlockValidation)
	if err != nil {
		return err
	}

	return nil
}

// CastVote allows a member to cast a vote on a proposal in the DAO.
func (manager *SYN2500Management) CastVote(vote common.DAOVote, token *common.SYN2500Token) error {
	if !token.Active {
		return errors.New("token is not active, voting is not allowed")
	}

	// Register the vote
	vote.VoteID = generateUniqueID()
	vote.Timestamp = time.Now()

	// Encrypt the vote for secure processing
	encryptedVote, err := manager.encryptVote(vote)
	if err != nil {
		return err
	}

	// Store the encrypted vote in the ledger
	err = ledger.StoreDAOVote(encryptedVote, synconsensus.SubBlockValidation)
	if err != nil {
		return err
	}

	// Append the vote to the token's voting log
	token.VoteLog = append(token.VoteLog, vote)

	return nil
}

// Propose creates a new proposal in the DAO.
func (manager *SYN2500Management) Propose(proposal common.DAOProposal, token *common.SYN2500Token) error {
	if !token.Active {
		return errors.New("only active members can propose")
	}

	// Register the proposal
	proposal.ProposalID = generateUniqueID()
	proposal.Timestamp = time.Now()

	// Encrypt the proposal for secure storage
	encryptedProposal, err := manager.encryptProposal(proposal)
	if err != nil {
		return err
	}

	// Store the encrypted proposal in the ledger
	err = ledger.StoreDAOProposal(encryptedProposal, synconsensus.SubBlockValidation)
	if err != nil {
		return err
	}

	// Append the proposal to the token's proposal log
	token.ProposalLog = append(token.ProposalLog, proposal)

	return nil
}

// TransferMembership transfers DAO membership by transferring the SYN2500 token to another user.
func (manager *SYN2500Management) TransferMembership(token *common.SYN2500Token, newOwner string) error {
	if !token.Active {
		return errors.New("token is not active, transfer not allowed")
	}

	// Update the token ownership
	token.Owner = newOwner

	// Store the updated token in the ledger
	err := ledger.UpdateDAOToken(token, synconsensus.SubBlockValidation)
	if err != nil {
		return err
	}

	return nil
}

// encryptVote encrypts the vote before storing it in the ledger.
func (manager *SYN2500Management) encryptVote(vote common.DAOVote) ([]byte, error) {
	voteBytes := serializeVote(vote)
	hashed := sha512.Sum512(voteBytes)
	encryptedBytes, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, manager.publicKey, hashed[:], nil)
	if err != nil {
		return nil, err
	}
	return encryptedBytes, nil
}

// encryptProposal encrypts the proposal before storing it in the ledger.
func (manager *SYN2500Management) encryptProposal(proposal common.DAOProposal) ([]byte, error) {
	proposalBytes := serializeProposal(proposal)
	hashed := sha512.Sum512(proposalBytes)
	encryptedBytes, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, manager.publicKey, hashed[:], nil)
	if err != nil {
		return nil, err
	}
	return encryptedBytes, nil
}

// serializeVote serializes the DAO vote into a byte slice for encryption.
func serializeVote(vote common.DAOVote) []byte {
	voteBytes, _ := json.Marshal(vote)
	return voteBytes
}

// serializeProposal serializes the DAO proposal into a byte slice for encryption.
func serializeProposal(proposal common.DAOProposal) []byte {
	proposalBytes, _ := json.Marshal(proposal)
	return proposalBytes
}

// generateUniqueID creates a unique ID for DAO tokens, proposals, and votes.
func generateUniqueID() string {
	timestamp := time.Now().UnixNano()
	hash := sha256.New()
	hash.Write([]byte(string(timestamp)))
	return hex.EncodeToString(hash.Sum(nil))
}
