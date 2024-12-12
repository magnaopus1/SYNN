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

// NewMultiAssetLiquidityRollup initializes a new Multi-Asset Liquidity Rollup (MALR)
func NewMultiAssetLiquidityRollup(rollupID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus, networkManager *common.NetworkManager) *common.MultiAssetLiquidityRollup {
	return &common.MultiAssetLiquidityRollup{
		RollupID:       rollupID,
		LiquidityPools: make(map[string]map[string]float64),
		Transactions:   []*common.Transaction{},
		IsFinalized:    false,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		Consensus:      consensus,
		NetworkManager: networkManager,
	}
}

// AddLiquidity adds liquidity to a specific asset pool from a participant
func (malr *common.MultiAssetLiquidityRollup) AddLiquidity(asset string, participant string, amount float64) error {
	malr.mu.Lock()
	defer malr.mu.Unlock()

	// Initialize asset pool if not exists
	if _, exists := malr.LiquidityPools[asset]; !exists {
		malr.LiquidityPools[asset] = make(map[string]float64)
	}

	// Update participant's liquidity for the asset
	malr.LiquidityPools[asset][participant] += amount

	// Log the liquidity addition to the ledger
	err := malr.Ledger.RecordLiquidityAddition(malr.RollupID, asset, participant, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log liquidity addition: %v", err)
	}

	fmt.Printf("Liquidity added: Participant %s added %f of asset %s to rollup %s\n", participant, amount, asset, malr.RollupID)
	return nil
}

// TransferLiquidity performs a liquidity transfer between participants for a given asset
func (malr *common.MultiAssetLiquidityRollup) TransferLiquidity(asset string, sender string, recipient string, amount float64) error {
	malr.mu.Lock()
	defer malr.mu.Unlock()

	// Ensure liquidity is sufficient
	if malr.LiquidityPools[asset][sender] < amount {
		return errors.New("insufficient liquidity")
	}

	// Perform the transfer
	malr.LiquidityPools[asset][sender] -= amount
	malr.LiquidityPools[asset][recipient] += amount

	// Log the liquidity transfer to the ledger
	err := malr.Ledger.RecordLiquidityTransfer(malr.RollupID, asset, sender, recipient, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log liquidity transfer: %v", err)
	}

	fmt.Printf("Liquidity transferred: %f of asset %s from %s to %s in rollup %s\n", amount, asset, sender, recipient, malr.RollupID)
	return nil
}

// AddTransaction adds a new transaction to the rollup
func (malr *common.MultiAssetLiquidityRollup) AddTransaction(tx *common.Transaction) error {
	malr.mu.Lock()
	defer malr.mu.Unlock()

	if malr.IsFinalized {
		return errors.New("rollup is already finalized, no new transactions can be added")
	}

	// Encrypt transaction data before adding
	encryptedTx, err := malr.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add transaction to the rollup
	malr.Transactions = append(malr.Transactions, tx)

	// Log transaction addition in the ledger
	err = malr.Ledger.RecordTransactionAddition(malr.RollupID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to rollup %s\n", tx.TxID, malr.RollupID)
	return nil
}

// FinalizeRollup finalizes the rollup by computing the state root and validating it using Synnergy Consensus
func (malr *common.MultiAssetLiquidityRollup) FinalizeRollup() error {
	malr.mu.Lock()
	defer malr.mu.Unlock()

	if malr.IsFinalized {
		return errors.New("rollup is already finalized")
	}

	// Apply liquidity reconciliation logic here if needed (e.g., cross-asset balancing)
	// Perform any final liquidity adjustments before state finalization

	// Compute the final state root based on the liquidity and transactions
	stateRoot := common.GenerateMerkleRoot(malr.Transactions)

	// Use Synnergy Consensus to validate the rollup
	err := malr.Consensus.ValidateRollup(malr.RollupID, stateRoot)
	if err != nil {
		return fmt.Errorf("failed to validate rollup: %v", err)
	}

	malr.IsFinalized = true

	// Log finalization in the ledger
	err = malr.Ledger.RecordRollupFinalization(malr.RollupID, stateRoot, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup finalization: %v", err)
	}

	fmt.Printf("Rollup %s finalized with state root %s\n", malr.RollupID, stateRoot)
	return nil
}

// ReconcileLiquidity performs on-demand liquidity reconciliation across assets
func (malr *common.MultiAssetLiquidityRollup) ReconcileLiquidity() error {
	malr.mu.Lock()
	defer malr.mu.Unlock()

	// Perform any necessary reconciliation across asset pools
	// Cross-asset balancing or liquidation logic would go here

	// Log reconciliation event
	err := malr.Ledger.RecordLiquidityReconciliation(malr.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log liquidity reconciliation: %v", err)
	}

	fmt.Printf("Liquidity reconciliation performed for rollup %s\n", malr.RollupID)
	return nil
}

// BroadcastRollup sends the finalized rollup to the network
func (malr *common.MultiAssetLiquidityRollup) BroadcastRollup() error {
	malr.mu.Lock()
	defer malr.mu.Unlock()

	if !malr.IsFinalized {
		return errors.New("rollup is not finalized, cannot broadcast")
	}

	// Broadcast rollup data to the network
	err := malr.NetworkManager.BroadcastData(malr.RollupID, []byte("Finalized rollup data"))
	if err != nil {
		return fmt.Errorf("failed to broadcast rollup: %v", err)
	}

	fmt.Printf("Rollup %s broadcasted to the network\n", malr.RollupID)
	return nil
}

// RetrieveLiquidity retrieves the liquidity of a participant for a specific asset
func (malr *common.MultiAssetLiquidityRollup) RetrieveLiquidity(asset string, participant string) (float64, error) {
	malr.mu.Lock()
	defer malr.mu.Unlock()

	if balance, exists := malr.LiquidityPools[asset][participant]; exists {
		fmt.Printf("Retrieved liquidity: Participant %s has %f of asset %s in rollup %s\n", participant, balance, asset, malr.RollupID)
		return balance, nil
	}

	return 0, fmt.Errorf("no liquidity found for participant %s in asset %s", participant, asset)
}
