package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)

// NewRollup initializes a new rollup with essential parameters
func NewRollup(rollupID, validatorAddress string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, cancellationManager *common.TransactionCancellationManager) *common.Rollup {
	return &common.Rollup{
		RollupID:        rollupID,
		Transactions:    []*common.Transaction{},
		IsFinalized:     false,
		Ledger:          ledgerInstance,
		Encryption:      encryptionService,
		NetworkManager:  networkManager,
		CancellationMgr: cancellationManager,
		TotalFees:       0,
		ValidatorAddress: validatorAddress,
		CreationTime:    time.Now(),
	}
}

// AddTransaction adds a new transaction to the rollup and encrypts it
func (r *common.Rollup) AddTransaction(tx *common.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsFinalized {
		return errors.New("rollup is already finalized, no new transactions can be added")
	}

	// Encrypt the transaction data
	encryptedTx, err := r.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add the transaction to the rollup
	r.Transactions = append(r.Transactions, tx)
	r.TotalFees += tx.Fee

	// Log the transaction addition in the ledger
	err = r.Ledger.RecordTransactionAddition(r.RollupID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to rollup %s\n", tx.TxID, r.RollupID)
	return nil
}

// FinalizeRollup finalizes the rollup by computing the final state root and closing it
func (r *common.Rollup) FinalizeRollup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsFinalized {
		return errors.New("rollup is already finalized")
	}

	// Compute the final state root based on all transactions
	r.StateRoot = common.GenerateMerkleRoot(r.Transactions)
	r.IsFinalized = true

	// Log the rollup finalization in the ledger
	err := r.Ledger.RecordRollupFinalization(r.RollupID, r.StateRoot, time.Now(), r.TotalFees)
	if err != nil {
		return fmt.Errorf("failed to log rollup finalization: %v", err)
	}

	fmt.Printf("Rollup %s finalized with state root %s and total fees %f\n", r.RollupID, r.StateRoot, r.TotalFees)
	return nil
}

// RetrieveTransaction retrieves a transaction by its ID from the rollup
func (r *common.Rollup) RetrieveTransaction(txID string) (*common.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, tx := range r.Transactions {
		if tx.TxID == txID {
			fmt.Printf("Retrieved transaction %s from rollup %s\n", txID, r.RollupID)
			return tx, nil
		}
	}

	return nil, fmt.Errorf("transaction %s not found in rollup %s", txID, r.RollupID)
}

// BroadcastRollup sends the finalized rollup data, including fees and state root, to the network
func (r *common.Rollup) BroadcastRollup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.IsFinalized {
		return errors.New("rollup is not finalized, cannot broadcast")
	}

	// Prepare data for broadcasting
	data := fmt.Sprintf("RollupID: %s, StateRoot: %s, TotalFees: %f", r.RollupID, r.StateRoot, r.TotalFees)
	encryptedData, err := r.Encryption.EncryptData([]byte(data), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt rollup data: %v", err)
	}

	// Broadcast the rollup data to the network
	err = r.NetworkManager.BroadcastData(r.RollupID, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to broadcast rollup: %v", err)
	}

	fmt.Printf("Rollup %s broadcasted to the network with state root %s\n", r.RollupID, r.StateRoot)
	return nil
}

// CancelTransaction cancels a transaction using the transaction cancellation manager if fraud is detected
func (r *common.Rollup) CancelTransaction(txID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Find the transaction to be cancelled
	tx, err := r.findTransactionByID(txID)
	if err != nil {
		return fmt.Errorf("transaction %s not found: %v", txID, err)
	}

	// Cancel the transaction using the cancellation manager
	err = r.CancellationMgr.CancelTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to cancel transaction %s: %v", txID, err)
	}

	// Log the cancellation event in the ledger
	err = r.Ledger.RecordTransactionCancellation(txID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction cancellation: %v", err)
	}

	fmt.Printf("Transaction %s cancelled in rollup %s\n", txID, r.RollupID)
	return nil
}

// findTransactionByID is a helper function to find a transaction by its ID
func (r *common.Rollup) findTransactionByID(txID string) (*common.Transaction, error) {
	for _, tx := range r.Transactions {
		if tx.TxID == txID {
			return tx, nil
		}
	}
	return nil, fmt.Errorf("transaction %s not found", txID)
}
