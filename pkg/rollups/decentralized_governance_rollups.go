package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewDecentralizedGovernanceRollup initializes a new rollup with decentralized governance functionality
func NewDecentralizedGovernanceRollup(rollupID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.DecentralizedGovernanceRollup {
	return &common.DecentralizedGovernanceRollup{
		RollupID:        rollupID,
		Transactions:    []*common.Transaction{},
		IsFinalized:     false,
		Ledger:          ledgerInstance,
		Encryption:      encryptionService,
		NetworkManager:  networkManager,
		VotingProposals: make(map[string]*common.GovernanceProposal),
		Participants:    participants,
	}
}

// AddTransaction adds a new transaction to the rollup
func (dgr *common.DecentralizedGovernanceRollup) AddTransaction(tx *common.Transaction) error {
	dgr.mu.Lock()
	defer dgr.mu.Unlock()

	if dgr.IsFinalized {
		return errors.New("rollup is already finalized, no new transactions can be added")
	}

	// Encrypt the transaction data
	encryptedTx, err := dgr.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add the transaction to the rollup
	dgr.Transactions = append(dgr.Transactions, tx)

	// Log the transaction addition in the ledger
	err = dgr.Ledger.RecordTransactionAddition(dgr.RollupID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to rollup %s\n", tx.TxID, dgr.RollupID)
	return nil
}

// CreateProposal creates a new governance proposal within the rollup
func (dgr *common.DecentralizedGovernanceRollup) CreateProposal(proposalID, title, description string, startTime, endTime time.Time) (*common.GovernanceProposal, error) {
	dgr.mu.Lock()
	defer dgr.mu.Unlock()

	// Check if the proposal already exists
	if _, exists := dgr.VotingProposals[proposalID]; exists {
		return nil, errors.New("governance proposal already exists")
	}

	// Create the new governance proposal
	proposal := &common.GovernanceProposal{
		ProposalID:  proposalID,
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		Voters:      make(map[string]bool),
	}

	dgr.VotingProposals[proposalID] = proposal

	// Log the proposal creation in the ledger
	err := dgr.Ledger.RecordGovernanceProposal(dgr.RollupID, proposalID, title, description, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to log proposal creation: %v", err)
	}

	fmt.Printf("Governance proposal %s created in rollup %s\n", proposalID, dgr.RollupID)
	return proposal, nil
}

// CastVote casts a vote on a governance proposal
func (dgr *common.DecentralizedGovernanceRollup) CastVote(proposalID, participant string, voteFor bool) error {
	dgr.mu.Lock()
	defer dgr.mu.Unlock()

	proposal, exists := dgr.VotingProposals[proposalID]
	if !exists {
		return fmt.Errorf("proposal %s not found in rollup %s", proposalID, dgr.RollupID)
	}

	// Check if the participant has already voted
	if _, hasVoted := proposal.Voters[participant]; hasVoted {
		return errors.New("participant has already voted on this proposal")
	}

	// Record the vote
	if voteFor {
		proposal.VotesFor++
	} else {
		proposal.VotesAgainst++
	}

	proposal.Voters[participant] = true

	// Log the vote in the ledger
	err := dgr.Ledger.RecordGovernanceVote(dgr.RollupID, proposalID, participant, voteFor, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log governance vote: %v", err)
	}

	fmt.Printf("Participant %s cast a vote on proposal %s in rollup %s\n", participant, proposalID, dgr.RollupID)
	return nil
}

// FinalizeProposal finalizes a governance proposal once the voting period ends
func (dgr *common.DecentralizedGovernanceRollup) FinalizeProposal(proposalID string) error {
	dgr.mu.Lock()
	defer dgr.mu.Unlock()

	proposal, exists := dgr.VotingProposals[proposalID]
	if !exists {
		return fmt.Errorf("proposal %s not found in rollup %s", proposalID, dgr.RollupID)
	}

	if proposal.IsFinalized {
		return errors.New("governance proposal is already finalized")
	}

	// Finalize the proposal if the voting period has ended
	if time.Now().After(proposal.EndTime) {
		proposal.IsFinalized = true
	} else {
		return errors.New("voting period for this proposal has not ended yet")
	}

	// Log the proposal finalization in the ledger
	err := dgr.Ledger.RecordGovernanceFinalization(dgr.RollupID, proposalID, proposal.VotesFor, proposal.VotesAgainst, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log governance finalization: %v", err)
	}

	fmt.Printf("Governance proposal %s finalized in rollup %s\n", proposalID, dgr.RollupID)
	return nil
}

// RetrieveProposal retrieves a governance proposal by its ID
func (dgr *common.DecentralizedGovernanceRollup) RetrieveProposal(proposalID string) (*common.GovernanceProposal, error) {
	dgr.mu.Lock()
	defer dgr.mu.Unlock()

	proposal, exists := dgr.VotingProposals[proposalID]
	if !exists {
		return nil, fmt.Errorf("proposal %s not found in rollup %s", proposalID, dgr.RollupID)
	}

	fmt.Printf("Retrieved governance proposal %s from rollup %s\n", proposalID, dgr.RollupID)
	return proposal, nil
}

// FinalizeRollup finalizes the rollup and computes the final state root
func (dgr *common.DecentralizedGovernanceRollup) FinalizeRollup() error {
	dgr.mu.Lock()
	defer dgr.mu.Unlock()

	if dgr.IsFinalized {
		return errors.New("rollup is already finalized")
	}

	// Compute the final state root
	dgr.StateRoot = common.GenerateMerkleRoot(dgr.Transactions)
	dgr.IsFinalized = true

	// Log the rollup finalization in the ledger
	err := dgr.Ledger.RecordRollupFinalization(dgr.RollupID, dgr.StateRoot, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup finalization: %v", err)
	}

	fmt.Printf("Rollup %s finalized with state root %s\n", dgr.RollupID, dgr.StateRoot)
	return nil
}

// BroadcastRollup broadcasts the finalized rollup and proposals to the network
func (dgr *common.DecentralizedGovernanceRollup) BroadcastRollup() error {
	dgr.mu.Lock()
	defer dgr.mu.Unlock()

	if !dgr.IsFinalized {
		return errors.New("rollup is not finalized, cannot broadcast")
	}

	// Broadcast the rollup data to the network
	err := dgr.NetworkManager.BroadcastData(dgr.RollupID, []byte(dgr.StateRoot))
	if err != nil {
		return fmt.Errorf("failed to broadcast rollup: %v", err)
	}

	fmt.Printf("Rollup %s broadcasted to the network\n", dgr.RollupID)
	return nil
}
