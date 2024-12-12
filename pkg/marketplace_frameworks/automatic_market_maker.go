package marketplace

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// NewAMMManager initializes the AMM Manager
func NewAMMManager(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.AMMManager {
	return &common.AMMManager{
		Pools:             make(map[string]*common.LiquidityPool),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// CreateLiquidityPool creates a new liquidity pool for two assets
func (amm *common.AMMManager) CreateLiquidityPool(poolID, assetA, assetB string, reserveA, reserveB float64) (*common.LiquidityPool, error) {
	amm.mu.Lock()
	defer amm.mu.Unlock()

	// Encrypt pool data
	poolData := fmt.Sprintf("PoolID: %s, AssetA: %s, AssetB: %s, ReserveA: %f, ReserveB: %f", poolID, assetA, assetB, reserveA, reserveB)
	encryptedData, err := amm.EncryptionService.EncryptData([]byte(poolData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt pool data: %v", err)
	}

	// Create the liquidity pool
	pool := &common.LiquidityPool{
		PoolID:         poolID,
		AssetA:         assetA,
		AssetB:         assetB,
		ReserveA:       reserveA,
		ReserveB:       reserveB,
		TotalLiquidity: reserveA * reserveB,
		CreatedAt:      time.Now(),
		PriceRatio:     reserveA / reserveB,
		Active:         true,
	}

	// Add the pool to the AMM
	amm.Pools[poolID] = pool

	// Log the pool creation in the ledger
	err = amm.Ledger.RecordLiquidityPoolCreation(poolID, assetA, assetB, reserveA, reserveB, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log liquidity pool creation: %v", err)
	}

	fmt.Printf("Liquidity pool %s created with assets %s and %s\n", poolID, assetA, assetB)
	return pool, nil
}

// AddLiquidity allows liquidity providers to add liquidity to a pool
func (amm *common.AMMManager) AddLiquidity(poolID string, amountA, amountB float64) error {
	amm.mu.Lock()
	defer amm.mu.Unlock()

	// Retrieve the liquidity pool
	pool, exists := amm.Pools[poolID]
	if !exists {
		return fmt.Errorf("liquidity pool %s not found", poolID)
	}

	// Update reserves and total liquidity
	pool.ReserveA += amountA
	pool.ReserveB += amountB
	pool.TotalLiquidity += amountA * amountB
	pool.PriceRatio = pool.ReserveA / pool.ReserveB

	// Log the liquidity addition in the ledger
	err := amm.Ledger.RecordLiquidityAddition(poolID, amountA, amountB, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log liquidity addition: %v", err)
	}

	fmt.Printf("Liquidity added to pool %s: %f of %s and %f of %s\n", poolID, amountA, pool.AssetA, amountB, pool.AssetB)
	return nil
}

// RemoveLiquidity allows liquidity providers to remove liquidity from a pool
func (amm *common.AMMManager) RemoveLiquidity(poolID string, amountA, amountB float64) error {
	amm.mu.Lock()
	defer amm.mu.Unlock()

	// Retrieve the liquidity pool
	pool, exists := amm.Pools[poolID]
	if !exists {
		return fmt.Errorf("liquidity pool %s not found", poolID)
	}

	// Ensure there's enough liquidity to remove
	if amountA > pool.ReserveA || amountB > pool.ReserveB {
		return fmt.Errorf("insufficient liquidity in pool %s to remove", poolID)
	}

	// Update reserves and total liquidity
	pool.ReserveA -= amountA
	pool.ReserveB -= amountB
	pool.TotalLiquidity -= amountA * amountB
	pool.PriceRatio = pool.ReserveA / pool.ReserveB

	// Log the liquidity removal in the ledger
	err := amm.Ledger.RecordLiquidityRemoval(poolID, amountA, amountB, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log liquidity removal: %v", err)
	}

	fmt.Printf("Liquidity removed from pool %s: %f of %s and %f of %s\n", poolID, amountA, pool.AssetA, amountB, pool.AssetB)
	return nil
}

// Swap executes an asset swap within a pool, updating reserves and applying fees
func (amm *common.AMMManager) Swap(poolID, fromAsset string, amountIn float64) (float64, error) {
	amm.mu.Lock()
	defer amm.mu.Unlock()

	// Retrieve the liquidity pool
	pool, exists := amm.Pools[poolID]
	if !exists {
		return 0, fmt.Errorf("liquidity pool %s not found", poolID)
	}

	// Calculate the output amount based on constant product formula
	var amountOut float64
	const feeRate = 0.003 // 0.3% fee
	if fromAsset == pool.AssetA {
		amountOut = (pool.ReserveB * amountIn) / (pool.ReserveA + amountIn*(1-feeRate))
		pool.ReserveA += amountIn
		pool.ReserveB -= amountOut
	} else if fromAsset == pool.AssetB {
		amountOut = (pool.ReserveA * amountIn) / (pool.ReserveB + amountIn*(1-feeRate))
		pool.ReserveB += amountIn
		pool.ReserveA -= amountOut
	} else {
		return 0, fmt.Errorf("invalid asset for swap: %s", fromAsset)
	}

	// Update fees collected
	fee := amountIn * feeRate
	pool.FeesCollected += fee

	// Log the swap in the ledger
	err := amm.Ledger.RecordSwapTransaction(poolID, fromAsset, amountIn, amountOut, fee, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to log swap transaction: %v", err)
	}

	fmt.Printf("Swap in pool %s: %f of %s for %f of the other asset\n", poolID, amountIn, fromAsset, amountOut)
	return amountOut, nil
}

// GetPoolDetails retrieves the details of a liquidity pool by its ID
func (amm *common.AMMManager) GetPoolDetails(poolID string) (*common.LiquidityPool, error) {
	amm.mu.Lock()
	defer amm.mu.Unlock()

	// Retrieve the pool
	pool, exists := amm.Pools[poolID]
	if !exists {
		return nil, fmt.Errorf("liquidity pool %s not found", poolID)
	}

	return pool, nil
}
