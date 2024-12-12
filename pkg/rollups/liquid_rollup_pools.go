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

// NewLiquidRollupPool initializes a new Liquid Rollup Pool (LRP)
func NewLiquidRollupPool(poolID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus, networkManager *common.NetworkManager) *common.LiquidRollupPool {
	return &common.LiquidRollupPool{
		PoolID:        poolID,
		Assets:        make(map[string]map[string]float64),
		YieldRates:    make(map[string]float64),
		Transactions:  []*common.Transaction{},
		IsFinalized:   false,
		Ledger:        ledgerInstance,
		Encryption:    encryptionService,
		Consensus:     consensus,
		NetworkManager: networkManager,
	}
}

// AddLiquidity allows a participant to add liquidity to the pool for a specific asset
func (lrp *common.LiquidRollupPool) AddLiquidity(asset string, participant string, amount float64) error {
	lrp.mu.Lock()
	defer lrp.mu.Unlock()

	// Initialize asset pool if it does not exist
	if _, exists := lrp.Assets[asset]; !exists {
		lrp.Assets[asset] = make(map[string]float64)
	}

	// Update the participant's balance for the asset
	lrp.Assets[asset][participant] += amount

	// Log liquidity addition in the ledger
	err := lrp.Ledger.RecordLiquidityAddition(lrp.PoolID, asset, participant, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log liquidity addition: %v", err)
	}

	fmt.Printf("Participant %s added %f of asset %s to pool %s\n", participant, amount, asset, lrp.PoolID)
	return nil
}

// RemoveLiquidity allows a participant to withdraw liquidity from the pool
func (lrp *common.LiquidRollupPool) RemoveLiquidity(asset string, participant string, amount float64) error {
	lrp.mu.Lock()
	defer lrp.mu.Unlock()

	// Check if the participant has enough liquidity
	if lrp.Assets[asset][participant] < amount {
		return errors.New("insufficient liquidity to withdraw")
	}

	// Deduct liquidity from participant's balance
	lrp.Assets[asset][participant] -= amount

	// Log liquidity removal in the ledger
	err := lrp.Ledger.RecordLiquidityRemoval(lrp.PoolID, asset, participant, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log liquidity removal: %v", err)
	}

	fmt.Printf("Participant %s removed %f of asset %s from pool %s\n", participant, amount, asset, lrp.PoolID)
	return nil
}

// DistributeYield automatically redistributes the yield for each participant based on their liquidity and asset
func (lrp *common.LiquidRollupPool) DistributeYield() error {
	lrp.mu.Lock()
	defer lrp.mu.Unlock()

	for asset, participants := range lrp.Assets {
		// Retrieve yield rate for the asset
		yieldRate, exists := lrp.YieldRates[asset]
		if !exists {
			continue
		}

		// Distribute yield to each participant proportionally
		for participant, balance := range participants {
			yield := balance * yieldRate
			lrp.Assets[asset][participant] += yield

			// Log yield distribution in the ledger
			err := lrp.Ledger.RecordYieldDistribution(lrp.PoolID, asset, participant, yield, time.Now())
			if err != nil {
				return fmt.Errorf("failed to log yield distribution: %v", err)
			}

			fmt.Printf("Distributed %f yield of asset %s to participant %s in pool %s\n", yield, asset, participant, lrp.PoolID)
		}
	}

	return nil
}

// AddTransaction adds a new transaction related to the liquidity pool
func (lrp *common.LiquidRollupPool) AddTransaction(tx *common.Transaction) error {
	lrp.mu.Lock()
	defer lrp.mu.Unlock()

	if lrp.IsFinalized {
		return errors.New("liquidity pool is finalized, no new transactions can be added")
	}

	// Encrypt the transaction before adding it
	encryptedTx, err := lrp.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add the transaction to the pool
	lrp.Transactions = append(lrp.Transactions, tx)

	// Log the transaction addition in the ledger
	err = lrp.Ledger.RecordTransactionAddition(lrp.PoolID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to liquidity pool %s\n", tx.TxID, lrp.PoolID)
	return nil
}

// FinalizePool finalizes the pool, redistributes yield, and validates the pool using Synnergy Consensus
func (lrp *common.LiquidRollupPool) FinalizePool() error {
	lrp.mu.Lock()
	defer lrp.mu.Unlock()

	if lrp.IsFinalized {
		return errors.New("liquidity pool is already finalized")
	}

	// Perform yield distribution before finalization
	err := lrp.DistributeYield()
	if err != nil {
		return fmt.Errorf("failed to distribute yield: %v", err)
	}

	// Validate the final state of the pool using Synnergy Consensus
	finalStateRoot := common.GenerateMerkleRoot(lrp.Transactions)
	err = lrp.Consensus.ValidateRollup(lrp.PoolID, finalStateRoot)
	if err != nil {
		return fmt.Errorf("failed to validate pool finalization: %v", err)
	}

	lrp.IsFinalized = true

	// Log pool finalization in the ledger
	err = lrp.Ledger.RecordRollupFinalization(lrp.PoolID, finalStateRoot, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log pool finalization: %v", err)
	}

	fmt.Printf("Liquidity pool %s finalized with state root %s\n", lrp.PoolID, finalStateRoot)
	return nil
}

// RetrieveLiquidity retrieves a participant's liquidity for a given asset
func (lrp *common.LiquidRollupPool) RetrieveLiquidity(asset string, participant string) (float64, error) {
	lrp.mu.Lock()
	defer lrp.mu.Unlock()

	if balance, exists := lrp.Assets[asset][participant]; exists {
		fmt.Printf("Participant %s has %f of asset %s in pool %s\n", participant, balance, asset, lrp.PoolID)
		return balance, nil
	}

	return 0, fmt.Errorf("no liquidity found for participant %s in asset %s", participant, asset)
}

// BroadcastPool sends the finalized liquidity pool data to the network
func (lrp *common.LiquidRollupPool) BroadcastPool() error {
	lrp.mu.Lock()
	defer lrp.mu.Unlock()

	if !lrp.IsFinalized {
		return errors.New("liquidity pool is not finalized, cannot broadcast")
	}

	// Broadcast pool data to the network
	err := lrp.NetworkManager.BroadcastData(lrp.PoolID, []byte("Finalized pool data"))
	if err != nil {
		return fmt.Errorf("failed to broadcast pool: %v", err)
	}

	fmt.Printf("Liquidity pool %s broadcasted to the network\n", lrp.PoolID)
	return nil
}

// SetYieldRate sets the yield rate for a specific asset
func (lrp *common.LiquidRollupPool) SetYieldRate(asset string, rate float64) error {
	lrp.mu.Lock()
	defer lrp.mu.Unlock()

	// Set the yield rate for the asset
	lrp.YieldRates[asset] = rate

	// Log the yield rate change in the ledger
	err := lrp.Ledger.RecordYieldRateChange(lrp.PoolID, asset, rate, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log yield rate change: %v", err)
	}

	fmt.Printf("Set yield rate of %f for asset %s in pool %s\n", rate, asset, lrp.PoolID)
	return nil
}
