package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewTemporalRollup initializes a new Temporal Rollup.
func NewTemporalRollup(rollupID string, pruneThreshold time.Duration, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, consensus *common.SynnergyConsensus) *common.TemporalRollup {
	return &common.TemporalRollup{
		RollupID:       rollupID,
		Transactions:   []*common.Transaction{},
		CreationTime:   time.Now(),
		PruneThreshold: pruneThreshold,
		IsFinalized:    false,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
		SynnergyConsensus: consensus,
	}
}

// AddTransaction adds a new transaction to the rollup.
func (r *common.TemporalRollup) AddTransaction(tx *common.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsFinalized {
		return errors.New("rollup is already finalized, no new transactions can be added")
	}

	// Encrypt the transaction data.
	encryptedTx, err := r.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add the transaction to the rollup.
	r.Transactions = append(r.Transactions, tx)

	// Log the transaction addition in the ledger.
	err = r.Ledger.RecordTransactionAddition(r.RollupID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to temporal rollup %s\n", tx.TxID, r.RollupID)
	return nil
}

// FinalizeRollup finalizes the rollup by computing the final state root.
func (r *common.TemporalRollup) FinalizeRollup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsFinalized {
		return errors.New("rollup is already finalized")
	}

	// Compute the final state root based on all transactions.
	r.StateRoot = common.GenerateMerkleRoot(r.Transactions)
	r.IsFinalized = true

	// Log the rollup finalization in the ledger.
	err := r.Ledger.RecordRollupFinalization(r.RollupID, r.StateRoot, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup finalization: %v", err)
	}

	fmt.Printf("Temporal rollup %s finalized with state root %s\n", r.RollupID, r.StateRoot)
	return nil
}

// PruneOldTransactions prunes transactions older than the defined prune threshold.
func (r *common.TemporalRollup) PruneOldTransactions() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Prune based on the prune threshold.
	now := time.Now()
	prunedCount := 0

	var remainingTransactions []*common.Transaction
	for _, tx := range r.Transactions {
		txTime := tx.Timestamp
		if now.Sub(txTime) <= r.PruneThreshold {
			remainingTransactions = append(remainingTransactions, tx)
		} else {
			prunedCount++
		}
	}

	r.Transactions = remainingTransactions

	// Log the pruning event in the ledger.
	err := r.Ledger.RecordTransactionPruning(r.RollupID, prunedCount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction pruning: %v", err)
	}

	fmt.Printf("Pruned %d old transactions from temporal rollup %s\n", prunedCount, r.RollupID)
	return nil
}

// VerifyRollup verifies the rollup using Synnergy Consensus.
func (r *common.TemporalRollup) VerifyRollup() (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verify the rollup using the Synnergy Consensus mechanism.
	valid, err := r.SynnergyConsensus.VerifyRollup(r.StateRoot, r.Transactions)
	if err != nil {
		return false, fmt.Errorf("failed to verify rollup: %v", err)
	}

	// Log the rollup verification in the ledger.
	err = r.Ledger.RecordRollupVerification(r.RollupID, time.Now())
	if err != nil {
		return false, fmt.Errorf("failed to log rollup verification: %v", err)
	}

	fmt.Printf("Temporal rollup %s verified successfully\n", r.RollupID)
	return valid, nil
}

// RetrieveTransaction retrieves a transaction by its ID from the rollup.
func (r *common.TemporalRollup) RetrieveTransaction(txID string) (*common.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, tx := range r.Transactions {
		if tx.TxID == txID {
			fmt.Printf("Retrieved transaction %s from temporal rollup %s\n", txID, r.RollupID)
			return tx, nil
		}
	}

	return nil, fmt.Errorf("transaction %s not found in temporal rollup %s", txID, r.RollupID)
}

// BroadcastRollup sends the finalized rollup data to the network.
func (r *common.TemporalRollup) BroadcastRollup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.IsFinalized {
		return errors.New("rollup is not finalized, cannot broadcast")
	}

	// Broadcast the rollup data to the network.
	err := r.NetworkManager.BroadcastData(r.RollupID, []byte(r.StateRoot))
	if err != nil {
		return fmt.Errorf("failed to broadcast temporal rollup: %v", err)
	}

	fmt.Printf("Temporal rollup %s broadcasted to the network\n", r.RollupID)
	return nil
}
