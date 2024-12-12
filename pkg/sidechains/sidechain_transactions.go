package sidechains

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/consensus"
)

// NewSidechainTransactionPool initializes a new transaction pool
func NewSidechainTransactionPool(ledgerInstance *ledger.Ledger, consensus *common.SynnergyConsensus, encryptionService *encryption.Encryption) *common.SidechainTransactionPool {
	return &common.SidechainTransactionPool{
		PendingTransactions: make(map[string]*common.SidechainTransaction),
		Ledger:              ledgerInstance,
		Consensus:           consensus,
		Encryption:          encryptionService,
	}
}

// CreateTransaction creates a new transaction within the sidechain
func (stp *common.SidechainTransactionPool) CreateTransaction(from, to string, amount, fee float64) (*common.SidechainTransaction, error) {
	stp.mu.Lock()
	defer stp.mu.Unlock()

	// Create new transaction
	tx := &common.SidechainTransaction{
		TxID:      common.GenerateUUID(), // Assuming a function to generate unique transaction IDs
		From:      from,
		To:        to,
		Amount:    amount,
		Fee:       fee,
		Timestamp: time.Now(),
	}

	// Encrypt transaction data before adding it to the pool
	encryptedTxID, err := stp.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt transaction data: %v", err)
	}
	tx.TxID = string(encryptedTxID)

	stp.PendingTransactions[tx.TxID] = tx

	// Log the transaction creation in the ledger
	err = stp.Ledger.RecordTransactionCreation(tx.TxID, from, to, amount, fee, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log transaction creation: %v", err)
	}

	fmt.Printf("Transaction %s created from %s to %s\n", tx.TxID, from, to)
	return tx, nil
}

// ValidateTransaction validates a pending transaction using Synnergy Consensus
func (stp *common.SidechainTransactionPool) ValidateTransaction(txID string) error {
	stp.mu.Lock()
	defer stp.mu.Unlock()

	tx, exists := stp.PendingTransactions[txID]
	if !exists {
		return fmt.Errorf("transaction %s not found in the pool", txID)
	}

	// Use Synnergy Consensus to validate the transaction
	err := stp.Consensus.ValidateTransaction(tx.TxID, tx.From, tx.To, tx.Amount)
	if err != nil {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Mark the transaction as validated
	tx.IsValidated = true

	// Log the transaction validation in the ledger
	err = stp.Ledger.RecordTransactionValidation(tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction validation: %v", err)
	}

	fmt.Printf("Transaction %s validated\n", tx.TxID)
	return nil
}

// FinalizeTransaction marks a validated transaction as completed
func (stp *common.SidechainTransactionPool) FinalizeTransaction(txID string) error {
	stp.mu.Lock()
	defer stp.mu.Unlock()

	tx, exists := stp.PendingTransactions[txID]
	if !exists {
		return fmt.Errorf("transaction %s not found in the pool", txID)
	}

	if !tx.IsValidated {
		return fmt.Errorf("transaction %s is not validated", txID)
	}

	// Mark the transaction as finalized
	tx.IsFinalized = true

	// Log the transaction finalization in the ledger
	err := stp.Ledger.RecordTransactionFinalization(tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction finalization: %v", err)
	}

	// Remove the transaction from the pending pool
	delete(stp.PendingTransactions, tx.TxID)

	fmt.Printf("Transaction %s finalized\n", tx.TxID)
	return nil
}

// RetrieveTransaction retrieves a transaction by its ID from the pending pool
func (stp *common.SidechainTransactionPool) RetrieveTransaction(txID string) (*common.SidechainTransaction, error) {
	stp.mu.Lock()
	defer stp.mu.Unlock()

	tx, exists := stp.PendingTransactions[txID]
	if !exists {
		return nil, fmt.Errorf("transaction %s not found", txID)
	}

	fmt.Printf("Retrieved transaction %s\n", txID)
	return tx, nil
}

// GetPendingTransactions returns all pending transactions
func (stp *common.SidechainTransactionPool) GetPendingTransactions() []*common.SidechainTransaction {
	stp.mu.Lock()
	defer stp.mu.Unlock()

	var pendingTxs []*common.SidechainTransaction
	for _, tx := range stp.PendingTransactions {
		pendingTxs = append(pendingTxs, tx)
	}

	return pendingTxs
}

// SyncTransaction synchronizes a transaction across sidechain nodes
func (stp *common.SidechainTransactionPool) SyncTransaction(txID string, destinationNodeID string) error {
	stp.mu.Lock()
	defer stp.mu.Unlock()

	tx, exists := stp.PendingTransactions[txID]
	if !exists {
		return fmt.Errorf("transaction %s not found", txID)
	}

	// Encrypt the transaction data before synchronizing
	encryptedTxData, err := stp.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction data: %v", err)
	}

	// Sync the transaction data (assuming a network manager or consensus sync function)
	err = stp.Consensus.SyncTransaction(tx.TxID, destinationNodeID, encryptedTxData)
	if err != nil {
		return fmt.Errorf("failed to sync transaction %s: %v", txID, err)
	}

	// Log the transaction sync in the ledger
	err = stp.Ledger.RecordTransactionSync(txID, destinationNodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction sync: %v", err)
	}

	fmt.Printf("Transaction %s synced to node %s\n", txID, destinationNodeID)
	return nil
}
