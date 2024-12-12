package dao

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/token/syn800"
	"synnergy_network/pkg/token/syn900"
)


// NewGovernanceTokenVotingSystem initializes a new GovernanceTokenVotingSystem.
func NewGovernanceTokenVotingSystem(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, syn800Token *syn800.Token, verifier *syn900.Verifier) *GovernanceTokenVotingSystem {
	return &GovernanceTokenVotingSystem{
		Proposals:         make(map[string]*common.GovernanceProposal),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
		Syn800Token:       syn800Token,
		Syn900Verifier:    verifier,
	}
}

// CreateProposal allows a user to create a new governance proposal.
func (gv *GovernanceTokenVotingSystem) CreateProposal(proposerWallet, proposalText string, votingDuration time.Duration) (*common.GovernanceProposal, error) {
	gv.mutex.Lock()
	defer gv.mutex.Unlock()

	// Generate a unique ID for the proposal.
	proposalID := generateUniqueID()

	// Create the governance proposal.
	proposal := &common.GovernanceProposal{
		ProposalID:   proposalID,
		ProposalText: proposalText,
		CreationTime: time.Now(),
		Deadline:     time.Now().Add(votingDuration),
		SubmittedBy:  proposerWallet,
		Status:       "Open",
		VoterRecords: make(map[string]bool),
	}

	// Store the proposal in the voting system.
	gv.Proposals[proposalID] = proposal

	// Record the proposal creation in the ledger.
	err := gv.Ledger.DAOLedger.RecordProposalCreation(proposalID, proposalText, proposerWallet, proposal.CreationTime)
	if err != nil {
		return nil, fmt.Errorf("failed to record proposal creation in ledger: %v", err)
	}

	fmt.Printf("Governance proposal %s created by %s\n", proposalID, proposerWallet)
	return proposal, nil
}

// VoteOnProposal allows a user to vote on a proposal using Syn800 governance tokens.
func (gv *GovernanceTokenVotingSystem) VoteOnProposal(voterWallet, proposalID, voteOption string, tokenAmount float64) error {
	gv.mutex.Lock()
	defer gv.mutex.Unlock()

	// Retrieve the proposal from the system.
	proposal, exists := gv.Proposals[proposalID]
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
	if proposal.VoterRecords[voterWallet] {
		return errors.New("you have already voted on this proposal")
	}

	// Verify the voter's identity using Syn900.
	verified, err := gv.Syn900Verifier.VerifyIdentity(voterWallet)
	if err != nil || !verified {
		return fmt.Errorf("identity verification failed for wallet %s", voterWallet)
	}

	// Check if the voter has sufficient Syn800 tokens.
	if !gv.Syn800Token.HasSufficientBalance(voterWallet, tokenAmount) {
		return fmt.Errorf("insufficient token balance for wallet %s", voterWallet)
	}

	// Encrypt the vote before submission.
	encryptedVote, err := gv.EncryptionService.EncryptData([]byte(voteOption), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt vote: %v", err)
	}

	// Add the encrypted vote to the ledger.
	err = gv.Ledger.DAOLedger.RecordVote(proposalID, voterWallet, encryptedVote, tokenAmount)
	if err != nil {
		return fmt.Errorf("failed to record vote in ledger: %v", err)
	}

	// Update the proposal with the vote.
	if voteOption == "yes" {
		proposal.YesVotes += tokenAmount
	} else if voteOption == "no" {
		proposal.NoVotes += tokenAmount
	} else {
		return errors.New("invalid vote option, must be 'yes' or 'no'")
	}

	// Mark the voter as having voted.
	proposal.VoterRecords[voterWallet] = true
	proposal.TotalVotes += tokenAmount

	// Deduct the token amount from the voter's balance (governance stake).
	err = gv.Syn800Token.DeductTokens(voterWallet, tokenAmount)
	if err != nil {
		return fmt.Errorf("failed to deduct tokens from wallet %s: %v", voterWallet, err)
	}

	fmt.Printf("User %s successfully voted %s on proposal %s with %f tokens\n", voterWallet, voteOption, proposalID, tokenAmount)
	return nil
}

// TallyVotes checks if a proposal has met the deadline and calculates the final result.
func (gv *GovernanceTokenVotingSystem) TallyVotes(proposalID string) (*common.GovernanceProposal, error) {
	gv.mutex.Lock()
	defer gv.mutex.Unlock()

	// Retrieve the proposal.
	proposal, exists := gv.Proposals[proposalID]
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
	err := gv.Ledger.DAOLedger.RecordProposalResult(proposalID, proposal.Status, proposal.YesVotes, proposal.NoVotes, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to record proposal result in ledger: %v", err)
	}

	fmt.Printf("Proposal %s has been %s. Yes: %f, No: %f\n", proposalID, proposal.Status, proposal.YesVotes, proposal.NoVotes)
	return proposal, nil
}

// ViewProposalResult allows users to view the final result of a proposal after the voting period ends.
func (gv *GovernanceTokenVotingSystem) ViewProposalResult(proposalID string) (*common.GovernanceProposal, error) {
	gv.mutex.Lock()
	defer gv.mutex.Unlock()

	// Retrieve the proposal.
	proposal, exists := gv.Proposals[proposalID]
	if !exists {
		return nil, errors.New("proposal not found")
	}

	// Ensure the voting period has ended.
	if proposal.Status == "Open" {
		return nil, errors.New("proposal voting is still open, results are not yet available")
	}

	return proposal, nil
}

