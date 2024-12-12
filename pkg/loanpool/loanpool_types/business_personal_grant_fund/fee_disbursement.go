package loanpool

import (

	"fmt"
	"time"

	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/common"
)


// NewBusinessPersonalGrantDisbursementManager initializes a new disbursement manager for the Business Personal Grant Fund.
func NewBusinessPersonalGrantDisbursementManager(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, initialBalance float64) *BusinessPersonalGrantDisbursementManager {
	return &BusinessPersonalGrantDisbursementManager{
		Ledger:            ledgerInstance,
		Consensus:         consensusEngine,
		FundBalance:       initialBalance,
		DisbursementQueue: []*BusinessPersonalGrantDisbursementQueueEntry{},
		QueueMaxTime:      30 * 24 * time.Hour, // 30-day queue duration
	}
}

// RequestDisbursement is called when the proposer wishes to disburse funds after confirmation for the Business Personal Grant Fund.
func (fdm *BusinessPersonalGrantDisbursementManager) RequestDisbursement(proposalID string, proposerWallet string, amount float64) error {
	fdm.mutex.Lock()
	defer fdm.mutex.Unlock()

	// Check if the funds are available
	if fdm.FundBalance >= amount {
		// Deduct the amount from the pool
		fdm.FundBalance -= amount

		// Log the disbursement in the ledger
		err := fdm.Ledger.RecordDisbursement(proposalID, proposerWallet, amount)
		if err != nil {
			return fmt.Errorf("failed to record disbursement in ledger: %v", err)
		}

		fmt.Printf("Disbursement of %f for proposal %s to wallet %s completed successfully.\n", amount, proposalID, proposerWallet)
		return nil
	}

	// If funds are unavailable, add the proposal to the disbursement queue
	queueEntry := &BusinessPersonalGrantDisbursementQueueEntry{
		ProposalID:        proposalID,
		ProposerWallet:    proposerWallet,
		RequestedAmount:   amount,
		DisbursementStart: time.Now(),
	}

	fdm.DisbursementQueue = append(fdm.DisbursementQueue, queueEntry)
	fmt.Printf("Proposal %s added to the disbursement queue due to insufficient funds.\n", proposalID)

	return nil
}

// ProcessDisbursementQueue processes the disbursement queue and disburses funds if they become available.
func (fdm *BusinessPersonalGrantDisbursementManager) ProcessDisbursementQueue() {
	fdm.mutex.Lock()
	defer fdm.mutex.Unlock()

	newQueue := []*BusinessPersonalGrantDisbursementQueueEntry{}

	for _, entry := range fdm.DisbursementQueue {
		// Check if funds are available
		if fdm.FundBalance >= entry.RequestedAmount {
			// Disburse the funds
			fdm.FundBalance -= entry.RequestedAmount

			// Log the disbursement in the ledger
			err := fdm.Ledger.RecordDisbursement(entry.ProposalID, entry.ProposerWallet, entry.RequestedAmount)
			if err != nil {
				fmt.Printf("Failed to record disbursement for proposal %s: %v\n", entry.ProposalID, err)
				newQueue = append(newQueue, entry) // Re-add to the queue if there's an error
				continue
			}

			fmt.Printf("Disbursement of %f for proposal %s to wallet %s completed from the queue.\n", entry.RequestedAmount, entry.ProposalID, entry.ProposerWallet)
		} else if time.Since(entry.DisbursementStart) <= fdm.QueueMaxTime {
			// Still within the 30-day window, so keep it in the queue
			newQueue = append(newQueue, entry)
		} else {
			// Time expired, mark the proposal as paused due to unavailable funds
			err := fdm.Ledger.RecordProposalPaused(entry.ProposalID, "Funds unavailable after 30 days in queue.")
			if err != nil {
				fmt.Printf("Failed to record proposal pause for %s: %v\n", entry.ProposalID, err)
			}

			fmt.Printf("Proposal %s marked as paused due to unavailable funds after 30 days.\n", entry.ProposalID)
		}
	}

	// Update the queue
	fdm.DisbursementQueue = newQueue
}

// GetFundBalance returns the current balance of the Business Personal Grant Fund.
func (fdm *BusinessPersonalGrantDisbursementManager) GetFundBalance() float64 {
	fdm.mutex.Lock()
	defer fdm.mutex.Unlock()

	return fdm.FundBalance
}

// GetDisbursementQueue returns the current proposals in the disbursement queue.
func (fdm *BusinessPersonalGrantDisbursementManager) GetDisbursementQueue() []*BusinessPersonalGrantDisbursementQueueEntry {
	fdm.mutex.Lock()
	defer fdm.mutex.Unlock()

	return fdm.DisbursementQueue
}