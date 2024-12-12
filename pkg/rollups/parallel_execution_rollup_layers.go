package rollups

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)

// NewParallelExecutionRollupLayer initializes a new PERL layer
func NewParallelExecutionRollupLayer(layerID, rollupID string, encryptionService *encryption.Encryption, ledgerInstance *ledger.Ledger, consensus *common.SynnergyConsensus, networkManager *common.NetworkManager) *common.ParallelExecutionRollupLayer {
	return &common.ParallelExecutionRollupLayer{
		LayerID:       layerID,
		RollupID:      rollupID,
		Transactions:  []*common.Transaction{},
		IsFinalized:   false,
		Encryption:    encryptionService,
		Ledger:        ledgerInstance,
		Consensus:     consensus,
		NetworkManager: networkManager,
	}
}

// AddTransaction adds a transaction to the parallel execution layer
func (pel *common.ParallelExecutionRollupLayer) AddTransaction(tx *common.Transaction) error {
	pel.mu.Lock()
	defer pel.mu.Unlock()

	if pel.IsFinalized {
		return errors.New("cannot add transactions, layer execution is already finalized")
	}

	// Encrypt transaction data before adding
	encryptedTx, err := pel.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	pel.Transactions = append(pel.Transactions, tx)

	// Record the transaction addition in the ledger
	err = pel.Ledger.RecordTransactionAddition(pel.LayerID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to record transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to parallel execution layer %s\n", tx.TxID, pel.LayerID)
	return nil
}

// ExecuteTransactions runs the transactions in parallel within the rollup layer
func (pel *common.ParallelExecutionRollupLayer) ExecuteTransactions() error {
	pel.mu.Lock()
	defer pel.mu.Unlock()

	if pel.IsFinalized {
		return errors.New("layer execution is already finalized")
	}

	// Initialize a wait group to handle parallel execution
	var wg sync.WaitGroup

	for _, tx := range pel.Transactions {
		wg.Add(1)

		// Execute each transaction in a Goroutine for parallelism
		go func(tx *common.Transaction) {
			defer wg.Done()

			// Simulate transaction processing (this would be more complex in a real scenario)
			time.Sleep(100 * time.Millisecond) // Simulate processing delay

			fmt.Printf("Processed transaction %s in layer %s\n", tx.TxID, pel.LayerID)
		}(tx)
	}

	// Wait for all Goroutines to finish processing
	wg.Wait()

	fmt.Printf("Executed %d transactions in parallel for layer %s\n", len(pel.Transactions), pel.LayerID)
	return nil
}

// FinalizeLayer finalizes the parallel execution layer and marks it as complete
func (pel *common.ParallelExecutionRollupLayer) FinalizeLayer() error {
	pel.mu.Lock()
	defer pel.mu.Unlock()

	if pel.IsFinalized {
		return errors.New("layer is already finalized")
	}

	// Validate all transactions using Synnergy Consensus
	for _, tx := range pel.Transactions {
		err := pel.Consensus.ValidateTransaction(tx)
		if err != nil {
			return fmt.Errorf("failed to validate transaction %s: %v", tx.TxID, err)
		}
	}

	// Mark the layer as finalized
	pel.IsFinalized = true

	// Record the layer finalization in the ledger
	err := pel.Ledger.RecordLayerFinalization(pel.LayerID, pel.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to record layer finalization: %v", err)
	}

	fmt.Printf("Parallel execution layer %s finalized for rollup %s\n", pel.LayerID, pel.RollupID)
	return nil
}

// BroadcastLayer broadcasts the finalized layer to the network
func (pel *common.ParallelExecutionRollupLayer) BroadcastLayer() error {
	pel.mu.Lock()
	defer pel.mu.Unlock()

	if !pel.IsFinalized {
		return errors.New("layer must be finalized before broadcasting")
	}

	// Serialize the finalized layer data
	layerData := fmt.Sprintf("LayerID: %s, RollupID: %s, Transactions: %d", pel.LayerID, pel.RollupID, len(pel.Transactions))

	// Broadcast the layer data
	err := pel.NetworkManager.BroadcastData(pel.LayerID, []byte(layerData))
	if err != nil {
		return fmt.Errorf("failed to broadcast layer data: %v", err)
	}

	// Log the broadcast event in the ledger
	err = pel.Ledger.RecordLayerBroadcast(pel.LayerID, pel.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log layer broadcast: %v", err)
	}

	fmt.Printf("Layer %s broadcasted for rollup %s\n", pel.LayerID, pel.RollupID)
	return nil
}

// RetrieveTransaction retrieves a specific transaction from the parallel execution layer by ID
func (pel *common.ParallelExecutionRollupLayer) RetrieveTransaction(txID string) (*common.Transaction, error) {
	pel.mu.Lock()
	defer pel.mu.Unlock()

	for _, tx := range pel.Transactions {
		if tx.TxID == txID {
			fmt.Printf("Retrieved transaction %s from layer %s\n", txID, pel.LayerID)
			return tx, nil
		}
	}

	return nil, fmt.Errorf("transaction %s not found in layer %s", txID, pel.LayerID)
}
