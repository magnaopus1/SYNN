package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
	"synnergy_network/pkg/transaction" // Importing the transaction package for transaction cancellation
)

// NewOptimisticRollup initializes a new optimistic rollup process
func NewOptimisticRollup(rollupID, nodeID string, transactions []*common.Transaction, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, consensus *common.SynnergyConsensus, cancellationManager *common.TransactionCancellationManager) *common.OptimisticRollup {
	return &common.OptimisticRollup{
		RollupID:        rollupID,
		NodeID:          nodeID,
		Transactions:    transactions,
		Ledger:          ledgerInstance,
		Encryption:      encryptionService,
		NetworkManager:  networkManager,
		Consensus:       consensus,
		SubmittedProofs: make(map[string]*common.FraudProof),
		CancellationMgr: cancellationManager, // Initialize cancellation manager
	}
}

// SubmitRollup submits the optimistic rollup transactions to the network
func (or *common.OptimisticRollup) SubmitRollup() error {
	or.mu.Lock()
	defer or.mu.Unlock()

	// Encrypt all transactions before submission
	for _, tx := range or.Transactions {
		encryptedTx, err := or.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt transaction: %v", err)
		}
		tx.TxID = string(encryptedTx)
	}

	// Submit transactions to the network
	for _, tx := range or.Transactions {
		for _, node := range or.NetworkManager.GetConnectedNodes() {
			err := or.NetworkManager.SendData(node, []byte(tx.TxID))
			if err != nil {
				return fmt.Errorf("failed to submit transaction %s to node %s: %v", tx.TxID, node.NodeID, err)
			}
			fmt.Printf("Transaction %s submitted to node %s\n", tx.TxID, node.NodeID)
		}
	}

	// Log the rollup submission in the ledger
	err := or.Ledger.RecordRollupSubmission(or.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup submission: %v", err)
	}

	fmt.Printf("Optimistic rollup %s submitted\n", or.RollupID)
	return nil
}

// ValidateTransactions optimistically validates the transactions in the rollup (assumes validity unless disputed)
func (or *common.OptimisticRollup) ValidateTransactions() error {
	or.mu.Lock()
	defer or.mu.Unlock()

	// Validate each transaction using the consensus mechanism
	for _, tx := range or.Transactions {
		err := or.Consensus.ValidateTransaction(tx)
		if err != nil {
			return fmt.Errorf("transaction %s validation failed: %v", tx.TxID, err)
		}
	}

	// Log the validation in the ledger
	err := or.Ledger.RecordTransactionValidation(or.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction validation: %v", err)
	}

	fmt.Printf("Transactions in rollup %s validated optimistically\n", or.RollupID)
	return nil
}

// SubmitFraudProof submits a fraud proof to challenge the validity of a specific transaction
func (or *common.OptimisticRollup) SubmitFraudProof(txID, challenger, evidence string) (*common.FraudProof, error) {
	or.mu.Lock()
	defer or.mu.Unlock()

	// Create and log the fraud proof
	fraudProof := &common.FraudProof{
		ProofID:    common.GenerateUniqueID(),
		TxID:       txID,
		Challenger: challenger,
		Evidence:   evidence,
		Timestamp:  time.Now(),
		IsResolved: false,
	}

	// Store the fraud proof
	or.SubmittedProofs[fraudProof.ProofID] = fraudProof

	// Log the fraud proof submission in the ledger
	err := or.Ledger.RecordFraudProofSubmission(fraudProof.ProofID, fraudProof.TxID, fraudProof.Challenger, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log fraud proof submission: %v", err)
	}

	fmt.Printf("Fraud proof %s submitted by challenger %s for transaction %s\n", fraudProof.ProofID, fraudProof.Challenger, fraudProof.TxID)
	return fraudProof, nil
}

// ResolveFraudProof resolves a fraud proof by verifying its evidence and deciding on its validity
func (or *common.OptimisticRollup) ResolveFraudProof(proofID string) error {
	or.mu.Lock()
	defer or.mu.Unlock()

	// Retrieve the fraud proof
	fraudProof, exists := or.SubmittedProofs[proofID]
	if !exists {
		return fmt.Errorf("fraud proof %s not found", proofID)
	}

	// Verify the provided evidence (this is a simplified example)
	// In practice, more advanced verification would take place here
	if fraudProof.Evidence == "" {
		return errors.New("invalid fraud proof: evidence is empty")
	}

	// If the fraud proof is valid, cancel the fraudulent transaction using the cancellation manager
	tx, err := or.findTransactionByID(fraudProof.TxID)
	if err != nil {
		return fmt.Errorf("transaction %s not found: %v", fraudProof.TxID, err)
	}

	// Use the transaction cancellation manager to cancel the fraudulent transaction
	err = or.CancellationMgr.CancelTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to cancel fraudulent transaction %s: %v", fraudProof.TxID, err)
	}

	// Mark the fraud proof as resolved
	fraudProof.IsResolved = true

	// Log the fraud proof resolution in the ledger
	err = or.Ledger.RecordFraudProofResolution(proofID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log fraud proof resolution: %v", err)
	}

	fmt.Printf("Fraud proof %s resolved for transaction %s and transaction cancelled\n", fraudProof.ProofID, fraudProof.TxID)
	return nil
}

// MonitorFraudProofs continuously monitors submitted fraud proofs for resolution
func (or *common.OptimisticRollup) MonitorFraudProofs(interval time.Duration) {
	for {
		time.Sleep(interval)
		for _, proof := range or.SubmittedProofs {
			if !proof.IsResolved {
				err := or.ResolveFraudProof(proof.ProofID)
				if err != nil {
					fmt.Printf("Failed to resolve fraud proof %s: %v\n", proof.ProofID, err)
				}
			}
		}
	}
}

// RetrieveFraudProof retrieves a fraud proof by its ID
func (or *common.OptimisticRollup) RetrieveFraudProof(proofID string) (*common.FraudProof, error) {
	or.mu.Lock()
	defer or.mu.Unlock()

	fraudProof, exists := or.SubmittedProofs[proofID]
	if !exists {
		return nil, fmt.Errorf("fraud proof %s not found", proofID)
	}

	fmt.Printf("Retrieved fraud proof %s\n", proofID)
	return fraudProof, nil
}

// findTransactionByID is a helper function to retrieve a transaction by its ID from the rollup
func (or *common.OptimisticRollup) findTransactionByID(txID string) (*common.Transaction, error) {
	for _, tx := range or.Transactions {
		if tx.TxID == txID {
			return tx, nil
		}
	}
	return nil, fmt.Errorf("transaction %s not found", txID)
}
