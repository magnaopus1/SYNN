package defi

import (
	"fmt"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewDeFiManagement initializes the DeFi management system
// Creates a new instance of the DeFi management system with initialized maps and services.
func NewDeFiManagement(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *DeFiManagement {
	if ledgerInstance == nil || encryptionService == nil {
		log.Fatal("Ledger and Encryption Service cannot be nil when initializing DeFiManagement")
	}

	return &DeFiManagement{
		LiquidityPools:       make(map[string]*LiquidityPool),
		AssetPools:           make(map[string]*AssetPool),
		YieldFarmingRecords:  make(map[string]*FarmingRecord),
		LoanManagement:       make(map[string]*Loan),
		SyntheticAssets:      make(map[string]*SyntheticAsset),
		InsurancePolicies
		DataOracccles
		DefiBetting
		PredictionMarkets
		LendingPools
		CrowdfundingCampaigns
		Ledger:               ledgerInstance,
		EncryptionService:    encryptionService,
	}
}

// CreateLiquidityPool creates a new liquidity pool for DeFi management
// Validates inputs, initializes the pool, and logs the operation in the ledger.
func (dm *DeFiManagement) CreateLiquidityPool(poolID string, initialLiquidity, rewardRate float64) (*LiquidityPool, error) {
	// Input validation
	if poolID == "" {
		return nil, fmt.Errorf("poolID cannot be empty")
	}
	if initialLiquidity <= 0 {
		return nil, fmt.Errorf("initial liquidity must be greater than zero")
	}
	if rewardRate < 0 {
		return nil, fmt.Errorf("reward rate cannot be negative")
	}

	// Thread-safe lock
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Check for duplicate pool
	if _, exists := dm.LiquidityPools[poolID]; exists {
		return nil, fmt.Errorf("liquidity pool with ID %s already exists", poolID)
	}

	// Create the liquidity pool
	pool := &LiquidityPool{
		PoolID:             poolID,
		TotalLiquidity:     initialLiquidity,
		AvailableLiquidity: initialLiquidity,
		RewardRate:         rewardRate,
		CreatedAt:          time.Now(),
		Status:             "Active",
	}

	// Add pool to management
	dm.LiquidityPools[poolID] = pool

	// Log pool creation in the ledger
	if err := dm.Ledger.DeFiLedger.RecordLiquidityPoolCreation(poolID, initialLiquidity, rewardRate); err != nil {
		log.Printf("Failed to log liquidity pool creation for %s: %v", poolID, err)
		return nil, fmt.Errorf("failed to log liquidity pool creation: %w", err)
	}

	// Log success
	log.Printf("Liquidity pool %s created successfully with initial liquidity: %.2f and reward rate: %.2f", poolID, initialLiquidity, rewardRate)
	return pool, nil
}

// CreateAssetPool creates an asset pool for managing synthetic or DeFi assets
// Validates inputs, encrypts sensitive data, initializes the pool, and logs the operation.
func (dm *DeFiManagement) CreateAssetPool(poolID, assetType string, totalAssets, rewardRate float64) (*AssetPool, error) {
	// Input validation
	if poolID == "" {
		return nil, fmt.Errorf("poolID cannot be empty")
	}
	if assetType == "" {
		return nil, fmt.Errorf("assetType cannot be empty")
	}
	if totalAssets <= 0 {
		return nil, fmt.Errorf("total assets must be greater than zero")
	}
	if rewardRate < 0 {
		return nil, fmt.Errorf("reward rate cannot be negative")
	}

	// Thread-safe lock
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Check for duplicate pool
	if _, exists := dm.AssetPools[poolID]; exists {
		return nil, fmt.Errorf("asset pool with ID %s already exists", poolID)
	}

	// Encrypt asset pool data
	assetData := fmt.Sprintf("PoolID: %s, AssetType: %s, TotalAssets: %.2f, RewardRate: %.2f", poolID, assetType, totalAssets, rewardRate)
	encryptedData, err := dm.EncryptionService.EncryptData("AES", []byte(assetData), common.EncryptionKey)
	if err != nil {
		log.Printf("Failed to encrypt asset pool data for %s: %v", poolID, err)
		return nil, fmt.Errorf("failed to encrypt asset pool data: %w", err)
	}

	// Create the asset pool
	assetPool := &AssetPool{
		PoolID:      poolID,
		TotalAssets: totalAssets,
		AssetType:   assetType,
		RewardRate:  rewardRate,
		CreatedAt:   time.Now(),
		Status:      "Active",
		EncryptedData: string(encryptedData),
	}

	// Add pool to management
	dm.AssetPools[poolID] = assetPool

	// Log asset pool creation in the ledger
	if err := dm.Ledger.DeFiLedger.RecordAssetPoolCreation(poolID, assetType, totalAssets, rewardRate); err != nil {
		log.Printf("Failed to log asset pool creation for %s: %v", poolID, err)
		return nil, fmt.Errorf("failed to log asset pool creation: %w", err)
	}

	// Log success
	log.Printf("Asset pool %s created successfully for asset type %s with total assets: %.2f and reward rate: %.2f", poolID, assetType, totalAssets, rewardRate)
	return assetPool, nil
}



// ManageYieldFarming adds a new yield farming record for a user staking liquidity
// Validates inputs, updates liquidity, and logs the operation in the ledger.
func (dm *DeFiManagement) ManageYieldFarming(farmingID, userID string, amountStaked float64, poolID string) (*FarmingRecord, error) {
	// Input validation
	if farmingID == "" || userID == "" || poolID == "" {
		return nil, fmt.Errorf("farmingID, userID, and poolID cannot be empty")
	}
	if amountStaked <= 0 {
		return nil, fmt.Errorf("amountStaked must be greater than zero")
	}

	// Thread-safe lock
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Retrieve the liquidity pool
	pool, exists := dm.LiquidityPools[poolID]
	if !exists {
		return nil, fmt.Errorf("liquidity pool %s not found", poolID)
	}

	// Ensure sufficient liquidity
	if amountStaked > pool.AvailableLiquidity {
		return nil, fmt.Errorf("insufficient liquidity in pool %s", poolID)
	}

	// Create the farming record
	farmingRecord := &FarmingRecord{
		FarmingID:      farmingID,
		UserID:         userID,
		AmountStaked:   amountStaked,
		RewardsEarned:  0,
		StakeTimestamp: time.Now(),
		Status:         "Active",
	}

	// Add farming record and update pool
	dm.YieldFarmingRecords[farmingID] = farmingRecord
	pool.AvailableLiquidity -= amountStaked

	// Log operation in the ledger
	err := dm.Ledger.DeFiLedger.RecordYieldFarming(farmingID, amountStaked, pool.RewardRate)
	if err != nil {
		return nil, fmt.Errorf("failed to log farming record: %w", err)
	}

	// Log success
	log.Printf("Yield farming record %s created for user %s with staked amount %.2f", farmingID, userID, amountStaked)
	return farmingRecord, nil
}

// DistributeRewards distributes rewards for a user participating in yield farming
// Calculates rewards based on time staked and logs the distribution.
func (dm *DeFiManagement) DistributeRewards(farmingID string) (float64, error) {
	// Input validation
	if farmingID == "" {
		return 0, fmt.Errorf("farmingID cannot be empty")
	}

	// Thread-safe lock
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Retrieve the farming record
	farmingRecord, exists := dm.YieldFarmingRecords[farmingID]
	if !exists {
		return 0, fmt.Errorf("farming record %s not found", farmingID)
	}

	// Retrieve the pool
	pool, exists := dm.LiquidityPools[farmingRecord.FarmingID]
	if !exists {
		return 0, fmt.Errorf("liquidity pool for farming record %s not found", farmingRecord.FarmingID)
	}

	// Calculate rewards
	duration := time.Since(farmingRecord.StakeTimestamp).Hours() / 24 // Days
	reward := farmingRecord.AmountStaked * pool.RewardRate * duration
	farmingRecord.RewardsEarned += reward

	// Log reward distribution in the ledger
	rewardDistribution := ledger.RewardDistribution{
		DistributionID: farmingID,
		Amount:         reward,
		Timestamp:      time.Now(),
	}
	err := dm.Ledger.RecordRewardDistribution(farmingID, rewardDistribution)
	if err != nil {
		return 0, fmt.Errorf("failed to log reward distribution: %w", err)
	}

	// Log success
	log.Printf("Rewards of %.2f distributed for farming record %s", reward, farmingID)
	return reward, nil
}

// CloseLiquidityPool closes a liquidity pool, preventing new stakes
// Marks the pool as inactive and logs the operation.
func (dm *DeFiManagement) CloseLiquidityPool(poolID string) error {
	// Input validation
	if poolID == "" {
		return fmt.Errorf("poolID cannot be empty")
	}

	// Thread-safe lock
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Retrieve the pool
	pool, exists := dm.LiquidityPools[poolID]
	if !exists {
		return fmt.Errorf("liquidity pool %s not found", poolID)
	}

	// Mark the pool as inactive
	pool.Status = "Paused"

	// Log operation in the ledger
	err := dm.Ledger.DeFiLedger.RecordLiquidityPoolClosure(poolID)
	if err != nil {
		return fmt.Errorf("failed to log pool closure: %w", err)
	}

	// Log success
	log.Printf("Liquidity pool %s has been closed", poolID)
	return nil
}

// GetLiquidityPoolDetails retrieves details of a liquidity pool by ID
// Returns the pool details if found, otherwise returns an error.
func (dm *DeFiManagement) GetLiquidityPoolDetails(poolID string) (*LiquidityPool, error) {
	// Input validation
	if poolID == "" {
		return nil, fmt.Errorf("poolID cannot be empty")
	}

	// Thread-safe lock
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Retrieve the pool
	pool, exists := dm.LiquidityPools[poolID]
	if !exists {
		return nil, fmt.Errorf("liquidity pool %s not found", poolID)
	}

	return pool, nil
}

// GetAssetPoolDetails retrieves details of an asset pool by ID
// Returns the asset pool details if found, otherwise returns an error.
func (dm *DeFiManagement) GetAssetPoolDetails(poolID string) (*AssetPool, error) {
	// Input validation
	if poolID == "" {
		return nil, fmt.Errorf("poolID cannot be empty")
	}

	// Thread-safe lock
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Retrieve the asset pool
	assetPool, exists := dm.AssetPools[poolID]
	if !exists {
		return nil, fmt.Errorf("asset pool %s not found", poolID)
	}

	return assetPool, nil
}

// GetFarmingRecord retrieves a user's yield farming record by ID
// Returns the farming record if found, otherwise returns an error.
func (dm *DeFiManagement) GetFarmingRecord(farmingID string) (*FarmingRecord, error) {
	// Input validation
	if farmingID == "" {
		return nil, fmt.Errorf("farmingID cannot be empty")
	}

	// Thread-safe lock
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Retrieve the farming record
	farmingRecord, exists := dm.YieldFarmingRecords[farmingID]
	if !exists {
		return nil, fmt.Errorf("farming record %s not found", farmingID)
	}

	return farmingRecord, nil
}
