package ledger

import (
	"errors"
	"fmt"
	"time"
)

// RecordDelegation records a delegation from one account to another.
func (l *GovernanceLedger) RecordDelegation(A *AccountsWalletLedger, delegateFrom, delegateTo string) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := A.AccountsWalletLedgerState.Accounts[delegateFrom]; !exists {
		return errors.New("delegator account does not exist")
	}

	if _, exists := A.AccountsWalletLedgerState.Accounts[delegateTo]; !exists {
		return errors.New("delegatee account does not exist")
	}

	if l.GovernanceRecords[delegateFrom].Delegations == nil {
		l.GovernanceRecords[delegateFrom] = GovernanceRecord{
			Delegations: make(map[string]string),
		}
	}

	l.GovernanceRecords[delegateFrom].Delegations[delegateFrom] = delegateTo
	fmt.Printf("Delegation recorded from %s to %s\n", delegateFrom, delegateTo)
	return nil
}


// RecordVote records a vote for a specific proposal.
func (l *GovernanceLedger) RecordVote(A *AccountsWalletLedger, proposalID, voter, vote string) error {
	l.Lock()
	defer l.Unlock()

	proposal, exists := l.GovernanceRecords[proposalID].Proposals[proposalID]
	if !exists {
		return errors.New("proposal does not exist")
	}

	if _, exists := A.AccountsWalletLedgerState.Accounts[voter]; !exists {
		return errors.New("voter account does not exist")
	}

	if vote != "Yes" && vote != "No" && vote != "Abstain" {
		return errors.New("invalid vote option")
	}

	// Record the vote in the proposal
	proposal.VotesFor += map[string]int{"Yes": 1, "No": 0}[vote]
	proposal.VotesAgainst += map[string]int{"No": 1, "Yes": 0}[vote]

	l.GovernanceRecords[proposalID].Proposals[proposalID] = proposal
	fmt.Printf("Vote recorded for proposal %s by voter %s: %s\n", proposalID, voter, vote)
	return nil
}


// RecordProposal records the creation of a governance proposal.
func (l *GovernanceLedger) RecordProposal(A *AccountsWalletLedger, proposalID, proposer, details string, creationFee float64) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := A.AccountsWalletLedgerState.Accounts[proposer]; !exists {
		return errors.New("proposer account does not exist")
	}

	if _, exists := l.GovernanceRecords[proposalID]; exists {
		return errors.New("proposal already exists")
	}

	newProposal := GovernanceProposal{
		ProposalID:  proposalID,
		Creator:     proposer,
		Title:       "New Proposal",
		Description: details,
		CreatedAt:   time.Now(),
		Status:      Pending,
		CreationFee: creationFee,
		VotesFor:    0,
		VotesAgainst: 0,
	}

	l.GovernanceRecords[proposalID] = GovernanceRecord{
		Proposals: map[string]GovernanceProposal{
			proposalID: newProposal,
		},
	}

	fmt.Printf("Proposal %s recorded by proposer %s\n", proposalID, proposer)
	return nil
}


// RecordExecution records the execution of a proposal.
func (l *GovernanceLedger) RecordExecution(proposalID string) error {
	l.Lock()
	defer l.Unlock()

	proposal, exists := l.GovernanceRecords[proposalID].Proposals[proposalID]
	if !exists {
		return errors.New("proposal does not exist")
	}

	if proposal.Status != Approved {
		return errors.New("proposal is not approved for execution")
	}

	// Ensure the proposal has not already been executed by checking its status
	if proposal.Status == "Executed" {
		return errors.New("proposal already executed")
	}

	// Mark the proposal as executed
	proposal.Status = "Executed"
	l.GovernanceRecords[proposalID].Proposals[proposalID] = proposal

	fmt.Printf("Proposal %s executed\n", proposalID)
	return nil
}



// GetTotalTransactionFeesForLastBlocks returns the total transaction fees for the last N blocks.
func (l *GovernanceLedger) GetTotalTransactionFeesForLastBlocks(C *BlockchainConsensusCoinLedger, blockCount uint64) (float64, error) {
	l.Lock()
	defer l.Unlock()

	if blockCount > uint64(C.BlockIndex) {
		return 0, errors.New("block count exceeds total blocks")
	}

	var totalFees float64
	startBlock := C.BlockIndex - int(blockCount)

	for block := startBlock; block <= C.BlockIndex; block++ {
		totalFees += l.GovernanceRecords[fmt.Sprintf("block_%d", block)].TransactionFees[uint64(block)]
	}

	fmt.Printf("Total transaction fees for the last %d blocks: %.2f\n", blockCount, totalFees)
	return totalFees, nil
}


// RecordTransactionFee records the transaction fee for a specific block.
func (l *GovernanceLedger) RecordGovernanceTransactionFee(blockIndex uint64, fee float64) error {
	l.Lock()
	defer l.Unlock()

	blockKey := fmt.Sprintf("block_%d", blockIndex)

	// Check if the governance record for the block exists
	if _, exists := l.GovernanceRecords[blockKey]; !exists {
		l.GovernanceRecords[blockKey] = GovernanceRecord{
			TransactionFees: make(map[uint64]float64),
		}
	}

	// Add the fee to the transaction fees for the block
	l.GovernanceRecords[blockKey].TransactionFees[blockIndex] += fee

	fmt.Printf("Transaction fee of %.2f recorded for block %d\n", fee, blockIndex)
	return nil
}


