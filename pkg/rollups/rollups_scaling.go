package rollups

import (
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewRollupScalingManager initializes a new Rollup Scaling Manager
func NewRollupScalingManager(rollupID string, scalingFactor, maxScalingLimit, minScalingLimit float64, transactions []*common.Transaction, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, cancellationManager *common.TransactionCancellationManager) *common.RollupScalingManager {
	return &common.RollupScalingManager{
		RollupID:        rollupID,
		ScalingFactor:   scalingFactor,
		MaxScalingLimit: maxScalingLimit,
		MinScalingLimit: minScalingLimit,
		Transactions:    transactions,
		Ledger:          ledgerInstance,
		Encryption:      encryptionService,
		NetworkManager:  networkManager,
		CancellationMgr: cancellationManager,
	}
}

// AdjustScalingFactor dynamically adjusts the scaling factor based on current network load and transaction volume
func (rsm *common.RollupScalingManager) AdjustScalingFactor(networkLoad float64, transactionVolume float64) error {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	// Calculate scaling factor based on network load and transaction volume
	newScalingFactor := rsm.ScalingFactor * (1 + networkLoad*0.05 + transactionVolume*0.1)

	// Enforce the scaling limits
	if newScalingFactor > rsm.MaxScalingLimit {
		newScalingFactor = rsm.MaxScalingLimit
	} else if newScalingFactor < rsm.MinScalingLimit {
		newScalingFactor = rsm.MinScalingLimit
	}

	rsm.ScalingFactor = newScalingFactor

	// Log the scaling adjustment in the ledger
	err := rsm.Ledger.RecordScalingAdjustment(rsm.RollupID, rsm.ScalingFactor, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log scaling adjustment: %v", err)
	}

	fmt.Printf("Rollup %s adjusted scaling factor to %f\n", rsm.RollupID, rsm.ScalingFactor)
	return nil
}

// OptimizeScaling continuously monitors network and transaction activity to optimize scaling in real-time
func (rsm *common.RollupScalingManager) OptimizeScaling(interval time.Duration) {
	for {
		time.Sleep(interval)

		// Example: Using dummy values for network load and transaction volume for demonstration
		networkLoad := 0.8  // Example network load
		transactionVolume := 0.5 // Example transaction volume

		err := rsm.AdjustScalingFactor(networkLoad, transactionVolume)
		if err != nil {
			fmt.Printf("Error adjusting scaling factor for rollup %s: %v\n", rsm.RollupID, err)
			return
		}
	}
}

// EncryptTransactions encrypts all transactions before submission
func (rsm *common.RollupScalingManager) EncryptTransactions() error {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	for _, tx := range rsm.Transactions {
		encryptedTx, err := rsm.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt transaction: %v", err)
		}
		tx.TxID = string(encryptedTx)
	}

	fmt.Printf("Transactions for rollup %s encrypted\n", rsm.RollupID)
	return nil
}

// SubmitTransactions submits the encrypted transactions to the network for further processing
func (rsm *common.RollupScalingManager) SubmitTransactions() error {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	for _, tx := range rsm.Transactions {
		for _, node := range rsm.NetworkManager.GetConnectedNodes() {
			err := rsm.NetworkManager.SendData(node, []byte(tx.TxID))
			if err != nil {
				return fmt.Errorf("failed to submit transaction %s to node %s: %v", tx.TxID, node.NodeID, err)
			}
			fmt.Printf("Transaction %s submitted to node %s\n", tx.TxID, node.NodeID)
		}
	}

	// Log transaction submission in the ledger
	err := rsm.Ledger.RecordTransactionSubmission(rsm.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction submission: %v", err)
	}

	return nil
}

// CancelFraudulentTransaction cancels a transaction using the transaction cancellation manager if fraud is detected
func (rsm *common.RollupScalingManager) CancelFraudulentTransaction(txID string) error {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	// Find the transaction to be cancelled
	tx, err := rsm.findTransactionByID(txID)
	if err != nil {
		return fmt.Errorf("transaction %s not found: %v", txID, err)
	}

	// Cancel the transaction using the cancellation manager
	err = rsm.CancellationMgr.CancelTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to cancel fraudulent transaction %s: %v", txID, err)
	}

	// Log the cancellation event in the ledger
	err = rsm.Ledger.RecordTransactionCancellation(txID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction cancellation: %v", err)
	}

	fmt.Printf("Fraudulent transaction %s cancelled for rollup %s\n", txID, rsm.RollupID)
	return nil
}

// findTransactionByID is a helper function to find a transaction by its ID
func (rsm *common.RollupScalingManager) findTransactionByID(txID string) (*common.Transaction, error) {
	for _, tx := range rsm.Transactions {
		if tx.TxID == txID {
			return tx, nil
		}
	}
	return nil, fmt.Errorf("transaction %s not found", txID)
}
