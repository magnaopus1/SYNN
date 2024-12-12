package rollups

import (

	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)

// NewMultiLayerTransactionPruning initializes a new Multi-Layer Transaction Pruning (MLTP) system
func NewMultiLayerTransactionPruning(rollupID string, pruningInterval, maxRetention time.Duration, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus, networkManager *common.NetworkManager) *common.MultiLayerTransactionPruning {
	return &common.MultiLayerTransactionPruning{
		RollupID:        rollupID,
		Transactions:    []*common.Transaction{},
		PrunedLayers:    make(map[int][]*common.Transaction),
		PruningInterval: pruningInterval,
		MaxRetention:    maxRetention,
		Ledger:          ledgerInstance,
		Consensus:       consensus,
		Encryption:      encryptionService,
		NetworkManager:  networkManager,
	}
}

// AddTransaction adds a transaction to the current rollup
func (mltp *common.MultiLayerTransactionPruning) AddTransaction(tx *common.Transaction) error {
	mltp.mu.Lock()
	defer mltp.mu.Unlock()

	// Encrypt transaction before adding
	encryptedTx, err := mltp.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	mltp.Transactions = append(mltp.Transactions, tx)

	// Log transaction addition in the ledger
	err = mltp.Ledger.RecordTransactionAddition(mltp.RollupID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to rollup %s\n", tx.TxID, mltp.RollupID)
	return nil
}

// PruneTransactions performs multi-layer pruning based on transaction age
func (mltp *common.MultiLayerTransactionPruning) PruneTransactions() error {
	mltp.mu.Lock()
	defer mltp.mu.Unlock()

	currentTime := time.Now()

	// Identify transactions to prune based on the MaxRetention period
	var pruned []*common.Transaction
	var remaining []*common.Transaction

	for _, tx := range mltp.Transactions {
		if currentTime.Sub(tx.Timestamp) > mltp.MaxRetention {
			pruned = append(pruned, tx)
		} else {
			remaining = append(remaining, tx)
		}
	}

	if len(pruned) == 0 {
		return nil // No transactions to prune
	}

	// Move pruned transactions into a new pruning layer
	pruningLevel := len(mltp.PrunedLayers) + 1
	mltp.PrunedLayers[pruningLevel] = pruned
	mltp.Transactions = remaining

	// Log the pruning event in the ledger
	err := mltp.Ledger.RecordPruningEvent(mltp.RollupID, pruningLevel, len(pruned), time.Now())
	if err != nil {
		return fmt.Errorf("failed to log pruning event: %v", err)
	}

	// Perform consensus validation on the pruned layer
	err = mltp.Consensus.ValidatePrunedLayer(mltp.RollupID, pruningLevel, pruned)
	if err != nil {
		return fmt.Errorf("failed to validate pruned layer: %v", err)
	}

	fmt.Printf("Pruned %d transactions from rollup %s into layer %d\n", len(pruned), mltp.RollupID, pruningLevel)
	return nil
}

// RetrievePrunedLayer retrieves a specific pruned layer by its level
func (mltp *common.MultiLayerTransactionPruning) RetrievePrunedLayer(layer int) ([]*common.Transaction, error) {
	mltp.mu.Lock()
	defer mltp.mu.Unlock()

	prunedLayer, exists := mltp.PrunedLayers[layer]
	if !exists {
		return nil, fmt.Errorf("pruned layer %d not found", layer)
	}

	fmt.Printf("Retrieved pruned layer %d for rollup %s\n", layer, mltp.RollupID)
	return prunedLayer, nil
}

// BroadcastPrunedLayer broadcasts pruned layer data to the network
func (mltp *common.MultiLayerTransactionPruning) BroadcastPrunedLayer(layer int) error {
	mltp.mu.Lock()
	defer mltp.mu.Unlock()

	prunedLayer, exists := mltp.PrunedLayers[layer]
	if !exists {
		return fmt.Errorf("pruned layer %d not found", layer)
	}

	// Serialize the pruned layer data
	layerData := fmt.Sprintf("RollupID: %s, Layer: %d, PrunedTransactions: %d", mltp.RollupID, layer, len(prunedLayer))

	// Broadcast the pruned layer data
	err := mltp.NetworkManager.BroadcastData(fmt.Sprintf("Layer-%d", layer), []byte(layerData))
	if err != nil {
		return fmt.Errorf("failed to broadcast pruned layer: %v", err)
	}

	// Log the broadcast event in the ledger
	err = mltp.Ledger.RecordLayerBroadcast(fmt.Sprintf("Layer-%d", layer), mltp.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log layer broadcast: %v", err)
	}

	fmt.Printf("Pruned layer %d broadcasted for rollup %s\n", layer, mltp.RollupID)
	return nil
}

// StartPruningScheduler starts an automatic pruning scheduler that runs at regular intervals
func (mltp *common.MultiLayerTransactionPruning) StartPruningScheduler() {
	go func() {
		for {
			time.Sleep(mltp.PruningInterval)

			err := mltp.PruneTransactions()
			if err != nil {
				fmt.Printf("Error pruning transactions for rollup %s: %v\n", mltp.RollupID, err)
			}
		}
	}()
}
