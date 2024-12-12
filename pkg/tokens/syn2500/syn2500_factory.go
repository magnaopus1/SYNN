package syn2500

import (
	"time"
	"errors"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rsa"
	"crypto/rand"
)

// SYN2500Token represents a token in a decentralized autonomous organization (DAO)
type SYN2500Token struct {
	TokenID          string               // Unique ID of the DAO token
	Owner            string               // The current owner of the token
	DAOID            string               // Identifier for the DAO to which the token belongs
	VotingPower      int                  // The voting power assigned to this token
	IssuedDate       time.Time            // The date the token was issued
	ActiveStatus     bool                 // Indicates if the token is active
	MembershipStatus string               // Current status of the DAO membership (e.g., "active", "revoked")
	Transferable     bool                 // Whether the token can be transferred between users
	ImmutableRecords []DAOTransactionLog  // List of immutable records for all the token transactions
	VotingRecords    []VotingRecord       // List of all voting activities associated with this token
	DelegatedVoting  string               // Address of the delegated voter if voting rights are delegated
	Proposals        []GovernanceProposal // Proposals created by this token holder
	ReputationScore  float64              // Optional: Reputation score of the token holder in the DAO ecosystem
}

// DAOTransactionLog represents a log entry for a transaction or activity involving the DAO token
type DAOTransactionLog struct {
	Timestamp    time.Time // The time the transaction took place
	Action       string    // The type of transaction (e.g., "Transfer", "Revocation", "Delegation")
	PerformedBy  string    // The entity who performed the action
	Details      string    // Additional details regarding the action
}

// VotingRecord keeps track of voting activities associated with the DAO token
type VotingRecord struct {
	VoteID       string    // Unique identifier for the vote
	ProposalID   string    // Identifier for the proposal being voted on
	VoteCast     string    // The vote cast (e.g., "yes", "no", "abstain")
	VotePower    int       // The voting power used for this vote
	VoteDate     time.Time // The date the vote was cast
	Delegated    bool      // Indicates if the vote was delegated to another party
}

// GovernanceProposal represents a proposal made by a DAO member
type GovernanceProposal struct {
	ProposalID    string    // Unique identifier for the governance proposal
	ProposalTitle string    // The title or subject of the proposal
	CreatedBy     string    // The DAO member who created the proposal
	CreationDate  time.Time // Date when the proposal was created
	Status        string    // Status of the proposal (e.g., "open", "approved", "rejected")
	VoteCount     int       // Number of votes received for the proposal
	VotingDeadline time.Time // Deadline for voting on the proposal
}

// CreateDAOProposal allows a member to submit a new proposal to the DAO
func (token *SYN2500Token) CreateDAOProposal(proposalID, title string, votingDeadline time.Time) GovernanceProposal {
	proposal := GovernanceProposal{
		ProposalID:    proposalID,
		ProposalTitle: title,
		CreatedBy:     token.Owner,
		CreationDate:  time.Now(),
		Status:        "open",
		VoteCount:     0,
		VotingDeadline: votingDeadline,
	}
	token.Proposals = append(token.Proposals, proposal)

	// Log the proposal creation
	token.ImmutableRecords = append(token.ImmutableRecords, DAOTransactionLog{
		Timestamp:   time.Now(),
		Action:      "Create Proposal",
		PerformedBy: token.Owner,
		Details:     "Created new proposal with ID " + proposalID,
	})

	return proposal
}

// CastVote allows the token owner to cast a vote on a proposal
func (token *SYN2500Token) CastVote(voteID, proposalID, vote string) error {
	if !token.ActiveStatus {
		return errors.New("membership inactive, voting not allowed")
	}

	votingRecord := VotingRecord{
		VoteID:     voteID,
		ProposalID: proposalID,
		VoteCast:   vote,
		VotePower:  token.VotingPower,
		VoteDate:   time.Now(),
		Delegated:  false,
	}

	token.VotingRecords = append(token.VotingRecords, votingRecord)

	// Log the voting action
	token.ImmutableRecords = append(token.ImmutableRecords, DAOTransactionLog{
		Timestamp:   time.Now(),
		Action:      "Cast Vote",
		PerformedBy: token.Owner,
		Details:     "Cast vote on proposal " + proposalID + " with vote " + vote,
	})

	return nil
}

// DelegateVoting delegates voting power to another DAO member
func (token *SYN2500Token) DelegateVoting(delegate string) error {
	if !token.Transferable {
		return errors.New("delegation not allowed for this token")
	}

	token.DelegatedVoting = delegate

	// Log the delegation action
	token.ImmutableRecords = append(token.ImmutableRecords, DAOTransactionLog{
		Timestamp:   time.Now(),
		Action:      "Delegate Voting",
		PerformedBy: token.Owner,
		Details:     "Delegated voting power to " + delegate,
	})

	return nil
}

// TransferOwnership transfers the token to another member
func (token *SYN2500Token) TransferOwnership(newOwner string) error {
	if !token.Transferable {
		return errors.New("transfer not allowed for this token")
	}

	previousOwner := token.Owner
	token.Owner = newOwner

	// Log the ownership transfer
	token.ImmutableRecords = append(token.ImmutableRecords, DAOTransactionLog{
		Timestamp:   time.Now(),
		Action:      "Transfer Ownership",
		PerformedBy: previousOwner,
		Details:     "Transferred ownership to " + newOwner,
	})

	return nil
}

// RevokeMembership revokes the membership of the current token holder
func (token *SYN2500Token) RevokeMembership(reason string) error {
	if !token.ActiveStatus {
		return errors.New("membership already revoked")
	}

	token.ActiveStatus = false
	token.MembershipStatus = "revoked"

	// Log the revocation
	token.ImmutableRecords = append(token.ImmutableRecords, DAOTransactionLog{
		Timestamp:   time.Now(),
		Action:      "Revoke Membership",
		PerformedBy: "DAO Admin",
		Details:     "Membership revoked for reason: " + reason,
	})

	return nil
}

// ReactivateMembership reactivates a revoked membership
func (token *SYN2500Token) ReactivateMembership() error {
	if token.ActiveStatus {
		return errors.New("membership already active")
	}

	token.ActiveStatus = true
	token.MembershipStatus = "active"

	// Log the reactivation
	token.ImmutableRecords = append(token.ImmutableRecords, DAOTransactionLog{
		Timestamp:   time.Now(),
		Action:      "Reactivate Membership",
		PerformedBy: "DAO Admin",
		Details:     "Membership reactivated",
	})

	return nil
}

// SYN2500Factory is responsible for the creation, management, and transaction of SYN2500 DAO tokens
type SYN2500Factory struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewSYN2500Factory initializes the factory with encryption keys
func NewSYN2500Factory() (*SYN2500Factory, error) {
	// Generate RSA keys for encryption
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &SYN2500Factory{
		privateKey: privKey,
		publicKey:  &privKey.PublicKey,
	}, nil
}

// IssueNewDAO generates a new SYN2500Token and registers it on the blockchain
func (factory *SYN2500Factory) IssueNewDAO(owner string, daoID string, votingPower int, isTransferable bool) (*common.SYN2500Token, error) {
	// Create new token metadata
	token := &common.SYN2500Token{
		TokenID:          generateUniqueID(),
		Owner:            owner,
		DAOID:            daoID,
		VotingPower:      votingPower,
		IssuedDate:       time.Now(),
		ActiveStatus:     true,
		Transferable:     isTransferable,
		MembershipStatus: "active",
		ImmutableRecords: []common.DAOTransactionLog{},
		VotingRecords:    []common.VotingRecord{},
		Proposals:        []common.GovernanceProposal{},
		ReputationScore:  0.0,
	}

	// Log the issuance
	token.ImmutableRecords = append(token.ImmutableRecords, common.DAOTransactionLog{
		Timestamp:   time.Now(),
		Action:      "Issue Token",
		PerformedBy: owner,
		Details:     "New DAO token issued for DAO: " + daoID,
	})

	// Encrypt token metadata
	encryptedToken, err := factory.encryptToken(token)
	if err != nil {
		return nil, err
	}

	// Register the token on the ledger via Synnergy Consensus
	err = ledger.AddNewTokenToLedger(encryptedToken, synconsensus.SubBlockValidation)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// TransferDAOToken transfers ownership of the SYN2500Token to a new owner
func (factory *SYN2500Factory) TransferDAOToken(token *common.SYN2500Token, newOwner string) error {
	if !token.Transferable {
		return errors.New("this DAO token is not transferable")
	}

	// Log the transfer
	token.ImmutableRecords = append(token.ImmutableRecords, common.DAOTransactionLog{
		Timestamp:   time.Now(),
		Action:      "Transfer Ownership",
		PerformedBy: token.Owner,
		Details:     "Ownership transferred to " + newOwner,
	})

	// Update token owner
	token.Owner = newOwner

	// Encrypt updated token and update ledger
	encryptedToken, err := factory.encryptToken(token)
	if err != nil {
		return err
	}
	return ledger.UpdateTokenInLedger(token.TokenID, encryptedToken)
}

// VoteOnProposal allows a DAO member to cast a vote on a specific proposal
func (factory *SYN2500Factory) VoteOnProposal(token *common.SYN2500Token, proposalID string, vote string) error {
	if !token.ActiveStatus {
		return errors.New("token membership inactive, cannot vote")
	}

	// Log the vote
	token.VotingRecords = append(token.VotingRecords, common.VotingRecord{
		VoteID:     generateUniqueID(),
		ProposalID: proposalID,
		VoteCast:   vote,
		VotePower:  token.VotingPower,
		VoteDate:   time.Now(),
		Delegated:  false,
	})

	// Log the voting action
	token.ImmutableRecords = append(token.ImmutableRecords, common.DAOTransactionLog{
		Timestamp:   time.Now(),
		Action:      "Cast Vote",
		PerformedBy: token.Owner,
		Details:     "Voted on proposal " + proposalID,
	})

	// Encrypt updated token and update ledger
	encryptedToken, err := factory.encryptToken(token)
	if err != nil {
		return err
	}
	return ledger.UpdateTokenInLedger(token.TokenID, encryptedToken)
}

// RevokeDAOToken revokes the DAO membership represented by the SYN2500Token
func (factory *SYN2500Factory) RevokeDAOToken(token *common.SYN2500Token, reason string) error {
	if !token.ActiveStatus {
		return errors.New("membership already revoked")
	}

	// Update token status
	token.ActiveStatus = false
	token.MembershipStatus = "revoked"

	// Log the revocation
	token.ImmutableRecords = append(token.ImmutableRecords, common.DAOTransactionLog{
		Timestamp:   time.Now(),
		Action:      "Revoke Membership",
		PerformedBy: "DAO Admin",
		Details:     "Membership revoked due to: " + reason,
	})

	// Encrypt updated token and update ledger
	encryptedToken, err := factory.encryptToken(token)
	if err != nil {
		return err
	}
	return ledger.UpdateTokenInLedger(token.TokenID, encryptedToken)
}

// Encrypt the token data before adding it to the ledger
func (factory *SYN2500Factory) encryptToken(token *common.SYN2500Token) ([]byte, error) {
	tokenBytes := serializeToken(token)
	hashed := sha256.Sum256(tokenBytes)
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, factory.publicKey, hashed[:], nil)
	if err != nil {
		return nil, err
	}
	return encryptedBytes, nil
}

// serializeToken converts the SYN2500Token struct into a byte slice for encryption
func serializeToken(token *common.SYN2500Token) []byte {
	// Convert the token struct into a serialized format (JSON, protobuf, etc.)
	// Assuming a simple JSON serialization for this example
	tokenBytes, _ := json.Marshal(token)
	return tokenBytes
}

// generateUniqueID creates a unique ID for DAO tokens, votes, or transactions
func generateUniqueID() string {
	timestamp := time.Now().UnixNano()
	hash := sha256.New()
	hash.Write([]byte(string(timestamp)))
	return hex.EncodeToString(hash.Sum(nil))
}
