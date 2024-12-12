package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewZeroLatencyRollupBridge initializes a new zero-latency rollup bridge.
func NewZeroLatencyRollupBridge(bridgeID, sourceRollupID, destinationRollupID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, consensus *common.SynnergyConsensus) *common.ZeroLatencyRollupBridge {
	return &common.ZeroLatencyRollupBridge{
		BridgeID:        bridgeID,
		SourceRollupID:  sourceRollupID,
		DestinationRollupID: destinationRollupID,
		Transactions:   []*common.Transaction{},
		IsFinalized:    false,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
		SynnergyConsensus: consensus,
	}
}

// AddTransaction adds a new transaction to the bridge for transfer.
func (zlb *common.ZeroLatencyRollupBridge) AddTransaction(tx *common.Transaction) error {
	zlb.mu.Lock()
	defer zlb.mu.Unlock()

	if zlb.IsFinalized {
		return errors.New("bridge is already finalized, no new transactions can be added")
	}

	// Encrypt the transaction data.
	encryptedTx, err := zlb.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add the transaction to the bridge.
	zlb.Transactions = append(zlb.Transactions, tx)

	// Log the transaction addition in the ledger.
	err = zlb.Ledger.RecordTransactionAddition(zlb.BridgeID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to bridge %s\n", tx.TxID, zlb.BridgeID)
	return nil
}

// FinalizeBridge finalizes the rollup bridge, computing the state root and marking the sync as complete.
func (zlb *common.ZeroLatencyRollupBridge) FinalizeBridge() error {
	zlb.mu.Lock()
	defer zlb.mu.Unlock()

	if zlb.IsFinalized {
		return errors.New("bridge is already finalized")
	}

	// Compute the final state root based on all transactions.
	zlb.StateRoot = common.GenerateMerkleRoot(zlb.Transactions)
	zlb.IsFinalized = true

	// Log the bridge finalization in the ledger.
	err := zlb.Ledger.RecordBridgeFinalization(zlb.BridgeID, zlb.StateRoot, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log bridge finalization: %v", err)
	}

	fmt.Printf("Bridge %s finalized with state root %s\n", zlb.BridgeID, zlb.StateRoot)
	return nil
}

// SyncBridge synchronizes the data across the bridge to the destination rollup.
func (zlb *common.ZeroLatencyRollupBridge) SyncBridge() error {
	zlb.mu.Lock()
	defer zlb.mu.Unlock()

	if !zlb.IsFinalized {
		return errors.New("bridge is not finalized, cannot sync")
	}

	// Broadcast the finalized state root and transactions to the destination rollup.
	err := zlb.NetworkManager.BroadcastData(zlb.DestinationRollupID, []byte(zlb.StateRoot))
	if err != nil {
		return fmt.Errorf("failed to broadcast bridge state root: %v", err)
	}

	for _, tx := range zlb.Transactions {
		err := zlb.NetworkManager.BroadcastData(zlb.DestinationRollupID, []byte(tx.TxID))
		if err != nil {
			return fmt.Errorf("failed to broadcast transaction %s: %v", tx.TxID, err)
		}
	}

	// Log the bridge synchronization event in the ledger.
	err = zlb.Ledger.RecordBridgeSync(zlb.BridgeID, zlb.SourceRollupID, zlb.DestinationRollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log bridge synchronization: %v", err)
	}

	fmt.Printf("Bridge %s synced between rollup %s and rollup %s\n", zlb.BridgeID, zlb.SourceRollupID, zlb.DestinationRollupID)
	return nil
}

// VerifyBridge uses Synnergy Consensus to verify the rollup bridge data.
func (zlb *common.ZeroLatencyRollupBridge) VerifyBridge() (bool, error) {
	zlb.mu.Lock()
	defer zlb.mu.Unlock()

	// Verify the bridge using Synnergy Consensus.
	valid, err := zlb.SynnergyConsensus.VerifyBridge(zlb.StateRoot, zlb.Transactions)
	if err != nil {
		return false, fmt.Errorf("failed to verify bridge: %v", err)
	}

	// Log the bridge verification in the ledger.
	err = zlb.Ledger.RecordBridgeVerification(zlb.BridgeID, time.Now())
	if err != nil {
		return false, fmt.Errorf("failed to log bridge verification: %v", err)
	}

	fmt.Printf("Bridge %s verified successfully\n", zlb.BridgeID)
	return valid, nil
}

// RetrieveTransaction retrieves a transaction by its ID from the bridge.
func (zlb *common.ZeroLatencyRollupBridge) RetrieveTransaction(txID string) (*common.Transaction, error) {
	zlb.mu.Lock()
	defer zlb.mu.Unlock()

	for _, tx := range zlb.Transactions {
		if tx.TxID == txID {
			fmt.Printf("Retrieved transaction %s from bridge %s\n", txID, zlb.BridgeID)
			return tx, nil
		}
	}

	return nil, fmt.Errorf("transaction %s not found in bridge %s", txID, zlb.BridgeID)
}
