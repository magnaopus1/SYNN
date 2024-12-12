package loanpool

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/consensus"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewUnsecuredLoanDisbursementManager initializes a new disbursement manager for the unsecured loan pool.
func NewUnsecuredLoanDisbursementManager(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *encryption.Encryption, initialBalance float64) *common.UnsecuredLoanDisbursementManager {
	return &common.UnsecuredLoanDisbursementManager{
		Ledger:            ledgerInstance,
		Consensus:         consensusEngine,
		FundBalance:       initialBalance,
		DisbursementQueue: []*common.UnsecuredLoanDisbursementQueueEntry{},
		QueueMaxTime:      48 * time.Hour, // Maximum of 48 hours to disburse
		EncryptionService: encryptionService,
		IssuerFeePercentage: 0.5, // 0.5% issuer fee
	}
}

// RequestDisbursement handles the request for disbursing funds for an approved loan proposal.
func (ldm *common.UnsecuredLoanDisbursementManager) RequestDisbursement(proposalID, proposerWallet string, amount, avgInterestRate float64) error {
	ldm.mutex.Lock()
	defer ldm.mutex.Unlock()

	// Check if the funds are available.
	if ldm.FundBalance >= amount {
		// Deduct the amount from the pool.
		ldm.FundBalance -= amount

		// Apply the issuer fee and distribute it to the authority nodes.
		issuerFee := (ldm.IssuerFeePercentage / 100) * amount
		ldm.distributeIssuerFee(proposalID, issuerFee)

		// Log the disbursement in the ledger with the applied average interest rate.
		err := ldm.Ledger.RecordDisbursementWithInterest(proposalID, proposerWallet, amount, avgInterestRate)
		if err != nil {
			return fmt.Errorf("failed to record disbursement in ledger: %v", err)
		}

		fmt.Printf("Disbursement of %.2f with interest rate %.2f%% for proposal %s to wallet %s completed successfully.\n", amount, avgInterestRate, proposalID, proposerWallet)
		return nil
	}

	// If funds are unavailable, add the proposal to the disbursement queue.
	queueEntry := &common.UnsecuredLoanDisbursementQueueEntry{
		ProposalID:        proposalID,
		ProposerWallet:    proposerWallet,
		RequestedAmount:   amount,
		DisbursementStart: time.Now(),
		AverageInterest:   avgInterestRate,
	}

	ldm.DisbursementQueue = append(ldm.DisbursementQueue, queueEntry)
	fmt.Printf("Proposal %s added to the disbursement queue due to insufficient funds.\n", proposalID)

	return nil
}

// ProcessDisbursementQueue processes the disbursement queue and disburses funds when available.
func (ldm *common.UnsecuredLoanDisbursementManager) ProcessDisbursementQueue() {
	ldm.mutex.Lock()
	defer ldm.mutex.Unlock()

	newQueue := []*common.UnsecuredLoanDisbursementQueueEntry{}

	for _, entry := range ldm.DisbursementQueue {
		// Check if funds are available.
		if ldm.FundBalance >= entry.RequestedAmount {
			// Disburse the funds.
			ldm.FundBalance -= entry.RequestedAmount

			// Apply the issuer fee and distribute it to the authority nodes.
			issuerFee := (ldm.IssuerFeePercentage / 100) * entry.RequestedAmount
			ldm.distributeIssuerFee(entry.ProposalID, issuerFee)

			// Log the disbursement in the ledger with the applied average interest rate.
			err := ldm.Ledger.RecordDisbursementWithInterest(entry.ProposalID, entry.ProposerWallet, entry.RequestedAmount, entry.AverageInterest)
			if err != nil {
				fmt.Printf("Failed to record disbursement for proposal %s: %v\n", entry.ProposalID, err)
				newQueue = append(newQueue, entry) // Re-add to the queue if there's an error.
				continue
			}

			fmt.Printf("Disbursement of %.2f with interest rate %.2f%% for proposal %s to wallet %s completed from the queue.\n", entry.RequestedAmount, entry.AverageInterest, entry.ProposalID, entry.ProposerWallet)
		} else if time.Since(entry.DisbursementStart) <= ldm.QueueMaxTime {
			// Still within the 48-hour window, so keep it in the queue.
			newQueue = append(newQueue, entry)
		} else {
			// Time expired, mark the proposal as paused due to unavailable funds.
			err := ldm.Ledger.RecordProposalPaused(entry.ProposalID, "Funds unavailable after 48 hours in queue.")
			if err != nil {
				fmt.Printf("Failed to record proposal pause for %s: %v\n", entry.ProposalID, err)
			}

			fmt.Printf("Proposal %s marked as paused due to unavailable funds after 48 hours.\n", entry.ProposalID)
		}
	}

	// Update the queue.
	ldm.DisbursementQueue = newQueue
}

// distributeIssuerFee distributes the issuer fee equally among all authority nodes that voted on the proposal.
func (ldm *common.UnsecuredLoanDisbursementManager) distributeIssuerFee(proposalID string, issuerFee float64) {
	nodes, err := ldm.Ledger.GetVotingAuthorityNodes(proposalID)
	if err != nil {
		fmt.Printf("Failed to retrieve voting nodes for proposal %s: %v\n", proposalID, err)
		return
	}

	// Split the issuer fee equally between the voting nodes.
	equalFee := issuerFee / float64(len(nodes))
	for _, node := range nodes {
		err := ldm.Ledger.RecordFeeDistribution(node.NodeID, equalFee)
		if err != nil {
			fmt.Printf("Failed to distribute fee to node %s: %v\n", node.NodeID, err)
		} else {
			fmt.Printf("Distributed fee of %.2f to node %s for proposal %s.\n", equalFee, node.NodeID, proposalID)
		}
	}
}

// GetFundBalance returns the current balance of the unsecured loan pool.
func (ldm *common.UnsecuredLoanDisbursementManager) GetFundBalance() float64 {
	ldm.mutex.Lock()
	defer ldm.mutex.Unlock()

	return ldm.FundBalance
}

// GetDisbursementQueue returns the current proposals in the disbursement queue.
func (ldm *common.UnsecuredLoanDisbursementManager) GetDisbursementQueue() []*common.UnsecuredLoanDisbursementQueueEntry {
	ldm.mutex.Lock()
	defer ldm.mutex.Unlock()

	return ldm.DisbursementQueue
}
