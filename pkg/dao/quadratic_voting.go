package dao

import (
	"errors"
	"fmt"
	"math"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/tokens/syn900"
)



// NewQuadraticVotingSystem initializes a new QuadraticVotingSystem.
func NewQuadraticVotingSystem(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, syn800Token *syn800.Token, verifier *syn900.Verifier) *QuadraticVotingSystem {
	return &QuadraticVotingSystem{
		Proposals:         make(map[string]*QuadraticProposal),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
		Syn800Token:       syn800Token,
		Syn900Verifier:    verifier,
	}
}

// CreateQuadraticProposal allows a user to create a new quadratic voting proposal.
func (qv *QuadraticVotingSystem) CreateQuadraticProposal(proposerWallet, proposalText string, votingDuration time.Duration) (*QuadraticProposal, error) {
	qv.mutex.Lock()
	defer qv.mutex.Unlock()

	// Generate a unique ID for the proposal.
	proposalID := generateUniqueID()

	// Create the quadratic proposal.
	proposal := &QuadraticProposal{
		ProposalID:   proposalID,
		ProposalText: proposalText,
		CreationTime: time.Now(),
		Deadline:     time.Now().Add(votingDuration),
		SubmittedBy:  proposerWallet,
		Status:       "Open",
		VoterRecords: make(map[string]float64),
	}

	// Store the proposal in the voting system.
	qv.Proposals[proposalID] = proposal

	// Record the proposal creation in the ledger.
	err := qv.Ledger.DAOLedger.RecordProposalCreation(proposalID, proposalText, proposerWallet, proposal.CreationTime)
	if err != nil {
		return nil, fmt.Errorf("failed to record proposal creation in ledger: %v", err)
	}

	fmt.Printf("Quadratic voting proposal %s created by %s\n", proposalID, proposerWallet)
	return proposal, nil
}

// VoteOnQuadraticProposal allows a user to vote on a quadratic proposal using Syn800 governance tokens.
func (qv *QuadraticVotingSystem) VoteOnQuadraticProposal(voterWallet, proposalID, voteOption string, tokenAmount float64) error {
	qv.mutex.Lock()
	defer qv.mutex.Unlock()

	// Retrieve the proposal from the system.
	proposal, exists := qv.Proposals[proposalID]
	if !exists {
		return errors.New("proposal not found")
	}

	// Ensure the voting is still open.
	if time.Now().After(proposal.Deadline) {
		return errors.New("voting period for this proposal has ended")
	}
	if proposal.Status != "Open" {
		return errors.New("voting on this proposal is closed")
	}

	// Check if the voter has already voted.
	if _, voted := proposal.VoterRecords[voterWallet]; voted {
		return errors.New("you have already voted on this proposal")
	}

	// Verify the voter's identity using Syn900.
	verified, err := qv.Syn900Verifier.VerifyIdentity(voterWallet)
	if err != nil || !verified {
		return fmt.Errorf("identity verification failed for wallet %s", voterWallet)
	}

	// Check if the voter has sufficient Syn800 tokens.
	if !qv.Syn800Token.HasSufficientBalance(voterWallet, tokenAmount) {
		return fmt.Errorf("insufficient token balance for wallet %s", voterWallet)
	}

	// Calculate the quadratic cost (tokens are squared).
	voteWeight := math.Sqrt(tokenAmount)

	// Encrypt the vote before submission.
	encryptedVote, err := qv.EncryptionService.EncryptData([]byte(voteOption), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt vote: %v", err)
	}

	// Add the encrypted vote to the ledger.
	err = qv.Ledger.DAOLedger.RecordVote(proposalID, voterWallet, encryptedVote, voteWeight)
	if err != nil {
		return fmt.Errorf("failed to record vote in ledger: %v", err)
	}

	// Update the proposal with the quadratic vote.
	if voteOption == "yes" {
		proposal.YesVotes += voteWeight
	} else if voteOption == "no" {
		proposal.NoVotes += voteWeight
	} else {
		return errors.New("invalid vote option, must be 'yes' or 'no'")
	}

	// Mark the voter as having voted and store the quadratic vote weight.
	proposal.VoterRecords[voterWallet] = voteWeight
	proposal.TotalVotes += voteWeight

	// Deduct the quadratic token amount from the voter's balance.
	err = qv.Syn800Token.DeductTokens(voterWallet, tokenAmount)
	if err != nil {
		return fmt.Errorf("failed to deduct tokens from wallet %s: %v", voterWallet, err)
	}

	fmt.Printf("User %s successfully voted %s on proposal %s with %f tokens (vote weight: %f)\n", voterWallet, voteOption, proposalID, tokenAmount, voteWeight)
	return nil
}

// TallyQuadraticVotes checks if a quadratic proposal has met the deadline and calculates the final result.
func (qv *QuadraticVotingSystem) TallyQuadraticVotes(proposalID string) (*QuadraticProposal, error) {
	qv.mutex.Lock()
	defer qv.mutex.Unlock()

	// Retrieve the proposal.
	proposal, exists := qv.Proposals[proposalID]
	if !exists {
		return nil, errors.New("proposal not found")
	}

	// Ensure the voting period has ended.
	if time.Now().Before(proposal.Deadline) {
		return nil, errors.New("voting period has not yet ended")
	}

	// Ensure the proposal is still open for tallying.
	if proposal.Status != "Open" {
		return proposal, errors.New("proposal has already been tallied")
	}

	// Determine the outcome.
	if proposal.YesVotes > proposal.NoVotes {
		proposal.Status = "Passed"
	} else {
		proposal.Status = "Rejected"
	}

	// Record the final result in the ledger.
	err := qv.Ledger.DAOLedger.RecordProposalResult(proposalID, proposal.Status, proposal.YesVotes, proposal.NoVotes, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to record proposal result in ledger: %v", err)
	}

	fmt.Printf("Proposal %s has been %s. Yes: %f, No: %f\n", proposalID, proposal.Status, proposal.YesVotes, proposal.NoVotes)
	return proposal, nil
}

// ViewQuadraticProposalResult allows users to view the final result of a quadratic proposal after the voting period ends.
func (qv *QuadraticVotingSystem) ViewQuadraticProposalResult(proposalID string) (*QuadraticProposal, error) {
	qv.mutex.Lock()
	defer qv.mutex.Unlock()

	// Retrieve the proposal.
	proposal, exists := qv.Proposals[proposalID]
	if !exists {
		return nil, errors.New("proposal not found")
	}

	// Ensure the voting period has ended.
	if proposal.Status == "Open" {
		return nil, errors.New("proposal voting is still open, results are not yet available")
	}

	return proposal, nil
}
