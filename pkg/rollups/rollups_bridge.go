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

// NewRollupBridge initializes a new bridge for rollups, capable of handling both tokens and SYNN
func NewRollupBridge(bridgeID, rollupID, targetChainID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, consensus *common.SynnergyConsensus) *common.RollupBridge {
	return &common.RollupBridge{
		BridgeID:       bridgeID,
		RollupID:       rollupID,
		TargetChainID:  targetChainID,
		Transactions:   make(map[string]*common.BridgeTx),
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
		Consensus:      consensus,
	}
}

// AddTransaction adds a new transaction to the bridge, specifying whether it's for a token or SYNN coin
func (rb *common.RollupBridge) AddTransaction(txID, sourceChain, destination string, amount float64, assetType common.AssetType) (*common.BridgeTx, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if _, exists := rb.Transactions[txID]; exists {
		return nil, errors.New("transaction already exists")
	}

	tx := &BridgeTx{
		TxID:        txID,
		SourceChain: sourceChain,
		Destination: destination,
		Amount:      amount,
		AssetType:   assetType,
		Timestamp:   time.Now(),
		IsFinalized: false,
	}

	// Encrypt the transaction details
	encryptedTxID, err := rb.Encryption.EncryptData([]byte(txID), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTxID)

	rb.Transactions[txID] = tx

	// Log the transaction in the ledger
	err = rb.Ledger.RecordBridgeTransaction(rb.BridgeID, txID, sourceChain, destination, amount, string(assetType), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log transaction: %v", err)
	}

	fmt.Printf("Transaction %s (Asset: %s) added to bridge %s\n", txID, assetType, rb.BridgeID)
	return tx, nil
}

// FinalizeTransaction marks a transaction as finalized across the bridge
func (rb *common.RollupBridge) FinalizeTransaction(txID string) error {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	tx, exists := rb.Transactions[txID]
	if !exists {
		return fmt.Errorf("transaction %s not found", txID)
	}

	if tx.IsFinalized {
		return errors.New("transaction already finalized")
	}

	tx.IsFinalized = true

	// Log the transaction finalization in the ledger
	err := rb.Ledger.RecordTransactionFinalization(txID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction finalization: %v", err)
	}

	fmt.Printf("Transaction %s finalized on bridge %s\n", txID, rb.BridgeID)
	return nil
}

// SyncBridge synchronizes bridge data across rollups and chains
func (rb *common.RollupBridge) SyncBridge() error {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	// Example sync logic - encrypt and transfer data between chains/rollups
	for txID, tx := range rb.Transactions {
		// Encrypt transaction details for syncing
		encryptedData, err := rb.Encryption.EncryptData([]byte(txID), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt transaction data: %v", err)
		}

		// Transfer the transaction data between the rollup and target chain
		err = rb.NetworkManager.TransferDataBetweenChains(rb.RollupID, rb.TargetChainID, encryptedData)
		if err != nil {
			return fmt.Errorf("failed to sync bridge transaction %s: %v", txID, err)
		}

		fmt.Printf("Synchronized transaction %s across rollup %s and chain %s\n", txID, rb.RollupID, rb.TargetChainID)
	}

	// Log the synchronization event in the ledger
	err := rb.Ledger.RecordBridgeSync(rb.BridgeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log bridge sync: %v", err)
	}

	fmt.Printf("Bridge %s synchronized between rollup %s and chain %s\n", rb.BridgeID, rb.RollupID, rb.TargetChainID)
	return nil
}

// ValidateBridge ensures that the bridge operates correctly, checking the state of all transactions
func (rb *common.RollupBridge) ValidateBridge() error {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	// Example validation logic for checking the state of all transactions
	for txID, tx := range rb.Transactions {
		if !tx.IsFinalized {
			return fmt.Errorf("unfinalized transaction found: %s", txID)
		}
	}

	// Use the consensus mechanism to validate the integrity of the bridge
	err := rb.Consensus.ValidateBridge(rb.BridgeID, rb.RollupID, rb.TargetChainID)
	if err != nil {
		return fmt.Errorf("bridge validation failed: %v", err)
	}

	// Log the bridge validation in the ledger
	err = rb.Ledger.RecordBridgeValidation(rb.BridgeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log bridge validation: %v", err)
	}

	fmt.Printf("Bridge %s validated successfully\n", rb.BridgeID)
	return nil
}

// RetrieveTransaction retrieves a transaction from the bridge by its ID
func (rb *common.RollupBridge) RetrieveTransaction(txID string) (*common.BridgeTx, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	tx, exists := rb.Transactions[txID]
	if !exists {
		return nil, fmt.Errorf("transaction %s not found", txID)
	}

	fmt.Printf("Retrieved transaction %s from bridge %s\n", txID, rb.BridgeID)
	return tx, nil
}
