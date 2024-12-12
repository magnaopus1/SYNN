package defi

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewYieldFarmingManager initializes the yield farming manager with a ledger instance and encryption service.
func NewYieldFarmingManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *YieldFarmingManager {
    log.Printf("[INFO] Initializing Yield Farming Manager with ledger and encryption service")
    return &YieldFarmingManager{
        FarmingPools:      make(map[string]*FarmingPool),
        StakingRecords:    make(map[string]*StakingRecord),
        Ledger:            ledgerInstance,
        EncryptionService: encryptionService,
    }
}


// CreateFarmingPool creates a new yield farming pool for liquidity providers.
func (yfm *YieldFarmingManager) CreateFarmingPool(tokenPair string, rewardRate, initialLiquidity, rewards float64) (*FarmingPool, error) {
    log.Printf("[INFO] Starting farming pool creation for token pair: %s", tokenPair)

    // Step 1: Acquire Lock for Concurrent Safety
    yfm.mu.Lock()
    defer yfm.mu.Unlock()

    // Step 2: Validate Inputs
    if tokenPair == "" {
        err := fmt.Errorf("tokenPair cannot be empty")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }
    if rewardRate <= 0 {
        err := fmt.Errorf("rewardRate must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }
    if initialLiquidity < 0 {
        err := fmt.Errorf("initialLiquidity cannot be negative")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }
    if rewards < 0 {
        err := fmt.Errorf("rewards cannot be negative")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Step 3: Generate Unique Pool ID
    poolID := generateUniqueID()
    log.Printf("[INFO] Generated unique PoolID: %s for token pair: %s", poolID, tokenPair)

    // Step 4: Encrypt Pool Data
    poolData := fmt.Sprintf("PoolID: %s, TokenPair: %s, Liquidity: %f, RewardRate: %f", poolID, tokenPair, initialLiquidity, rewardRate)
    encryptedData, err := yfm.EncryptionService.EncryptData("AES", []byte(poolData), common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt pool data for PoolID %s: %v", poolID, err)
        return nil, fmt.Errorf("failed to encrypt pool data: %w", err)
    }

    // Step 5: Create the Farming Pool
    pool := &FarmingPool{
        PoolID:         poolID,
        TokenPair:      tokenPair,
        TotalLiquidity: initialLiquidity,
        RewardRate:     rewardRate,
        Rewards:        rewards,
        CreatedAt:      time.Now(),
        Status:         "Active",
        EncryptedData:  string(encryptedData), // Store encrypted data as a string
    }

    // Step 6: Add Pool to Manager
    yfm.FarmingPools[poolID] = pool
    log.Printf("[INFO] Farming pool added to manager: PoolID %s", poolID)

    // Step 7: Log Pool Creation in Ledger
    err = yfm.Ledger.DeFiLedger.RecordFarmingPoolCreation(poolID, initialLiquidity, rewardRate)
    if err != nil {
        log.Printf("[ERROR] Failed to log farming pool creation in ledger for PoolID %s: %v", poolID, err)
        return nil, fmt.Errorf("failed to log pool creation in the ledger: %w", err)
    }

    // Step 8: Final Success Log
    log.Printf("[SUCCESS] Farming pool %s created for token pair %s with initial liquidity %f and reward rate %f", poolID, tokenPair, initialLiquidity, rewardRate)
    return pool, nil
}


// StakeLiquidity allows a user to stake liquidity in a farming pool and start earning rewards.
func (yfm *YieldFarmingManager) StakeLiquidity(poolID, stakerAddress string, amount float64) (*StakingRecord, error) {
    log.Printf("[INFO] Starting liquidity staking for pool: %s by staker: %s", poolID, stakerAddress)

    // Acquire lock for thread-safety
    yfm.mu.Lock()
    defer yfm.mu.Unlock()

    // Validate inputs
    if poolID == "" {
        return nil, fmt.Errorf("poolID cannot be empty")
    }
    if stakerAddress == "" {
        return nil, fmt.Errorf("stakerAddress cannot be empty")
    }
    if amount <= 0 {
        return nil, fmt.Errorf("amount must be greater than zero")
    }

    // Retrieve the farming pool
    pool, exists := yfm.FarmingPools[poolID]
    if !exists {
        return nil, fmt.Errorf("farming pool %s not found", poolID)
    }

    // Generate a unique stake ID
    stakeID := generateUniqueID()
    log.Printf("[INFO] Generated stake ID: %s for staker: %s", stakeID, stakerAddress)

    // Encrypt staking data
    stakeData := fmt.Sprintf("StakeID: %s, Staker: %s, Amount: %f", stakeID, stakerAddress, amount)
    encryptedData, err := yfm.EncryptionService.EncryptData("AES", []byte(stakeData), common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt staking data for stake ID %s: %v", stakeID, err)
        return nil, fmt.Errorf("failed to encrypt staking data: %w", err)
    }

    // Create staking record
    stakingRecord := &StakingRecord{
        StakeID:        stakeID,
        StakerAddress:  stakerAddress,
        AmountStaked:   amount,
        StakeTimestamp: time.Now(),
        RewardEarned:   0, // Initial reward is 0
        EncryptedData:  string(encryptedData),
    }

    // Add to staking records and update pool liquidity
    yfm.StakingRecords[stakeID] = stakingRecord
    pool.TotalLiquidity += amount
    log.Printf("[INFO] Updated total liquidity for pool %s: %f", poolID, pool.TotalLiquidity)

    // Log staking event in the ledger
    err = yfm.Ledger.DeFiLedger.RecordLiquidityStake(stakeID, poolID, amount)
    if err != nil {
        log.Printf("[ERROR] Failed to log staking event in ledger for stake ID %s: %v", stakeID, err)
        return nil, fmt.Errorf("failed to log staking event in the ledger: %w", err)
    }

    log.Printf("[SUCCESS] Liquidity of %f staked in pool %s by %s", amount, poolID, stakerAddress)
    return stakingRecord, nil
}



// ClaimRewards allows a staker to claim the rewards they have earned from staking liquidity
func (yfm *YieldFarmingManager) ClaimRewards(stakeID string) (float64, error) {
    log.Printf("[INFO] Processing reward claim for stake ID: %s", stakeID)

    // Acquire lock for thread-safety
    yfm.mu.Lock()
    defer yfm.mu.Unlock()

    // Retrieve the staking record
    stakingRecord, exists := yfm.StakingRecords[stakeID]
    if !exists {
        err := fmt.Errorf("staking record %s not found", stakeID)
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Retrieve the associated pool
    pool, exists := yfm.FarmingPools[stakingRecord.StakerAddress]
    if !exists {
        err := fmt.Errorf("farming pool for staker %s not found", stakingRecord.StakerAddress)
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Calculate rewards
    duration := time.Since(stakingRecord.StakeTimestamp).Hours() / 24 // Staking duration in days
    reward := stakingRecord.AmountStaked * pool.RewardRate * duration
    stakingRecord.RewardEarned += reward

    // Log the reward claim in the ledger
    rewardStr := fmt.Sprintf("%.2f", reward)
    currentTime := float64(time.Now().Unix()) // Current timestamp as float64
    err := yfm.Ledger.DeFiLedger.RecordRewardClaim(stakeID, rewardStr, currentTime)
    if err != nil {
        log.Printf("[ERROR] Failed to log reward claim in ledger for stake ID %s: %v", stakeID, err)
        return 0, fmt.Errorf("failed to log reward claim in the ledger: %w", err)
    }

    log.Printf("[SUCCESS] Rewards of %f claimed for stake ID %s by staker %s", reward, stakeID, stakingRecord.StakerAddress)
    return reward, nil
}


// UnstakeLiquidity allows a user to withdraw their liquidity and stop earning rewards
func (yfm *YieldFarmingManager) UnstakeLiquidity(stakeID string) (float64, error) {
    log.Printf("[INFO] Starting unstaking process for stake ID: %s", stakeID)

    // Acquire lock for thread safety
    yfm.mu.Lock()
    defer yfm.mu.Unlock()

    // Validate the staking record
    stakingRecord, exists := yfm.StakingRecords[stakeID]
    if !exists {
        err := fmt.Errorf("staking record %s not found", stakeID)
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Retrieve the associated farming pool
    pool, exists := yfm.FarmingPools[stakingRecord.StakerAddress]
    if !exists {
        err := fmt.Errorf("farming pool for staker %s not found", stakingRecord.StakerAddress)
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Calculate and claim rewards before unstaking
    reward, err := yfm.ClaimRewards(stakeID)
    if err != nil {
        log.Printf("[ERROR] Failed to claim rewards for stake ID %s: %v", stakeID, err)
        return 0, err
    }
    log.Printf("[INFO] Claimed rewards of %f for stake ID %s", reward, stakeID)

    // Unstake the liquidity and update pool's liquidity
    stakedAmount := stakingRecord.AmountStaked
    pool.TotalLiquidity -= stakedAmount
    delete(yfm.StakingRecords, stakeID) // Remove the staking record

    // Log the unstaking event in the ledger
    stakedAmountStr := fmt.Sprintf("%.2f", stakedAmount) // Format staked amount as a string
    currentTime := float64(time.Now().Unix())            // Current timestamp
    err = yfm.Ledger.DeFiLedger.RecordLiquidityUnstake(stakeID, stakedAmountStr, currentTime)
    if err != nil {
        log.Printf("[ERROR] Failed to log unstaking event for stake ID %s: %v", stakeID, err)
        return 0, fmt.Errorf("failed to log unstaking event in the ledger: %w", err)
    }

    log.Printf("[SUCCESS] Unstaked %f liquidity from pool by %s", stakedAmount, stakingRecord.StakerAddress)
    return stakedAmount, nil
}



// generateUniqueID creates a cryptographically secure unique ID
func generateUniqueID() string {
	id := make([]byte, 16)
	rand.Read(id)  // Corrected: removed extra 'r'
	return hex.EncodeToString(id)
}


// CloseFarmingPool allows an admin to close a farming pool, preventing new stakes but allowing existing stakers to withdraw
func (yfm *YieldFarmingManager) CloseFarmingPool(poolID string) error {
    log.Printf("[INFO] Initiating closure for farming pool: %s", poolID)

    // Acquire lock for thread safety
    yfm.mu.Lock()
    defer yfm.mu.Unlock()

    // Validate the farming pool
    pool, exists := yfm.FarmingPools[poolID]
    if !exists {
        err := fmt.Errorf("farming pool %s not found", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Mark the pool as inactive
    pool.Status = "Inactive"
    log.Printf("[INFO] Farming pool %s marked as inactive", poolID)

    // Log the pool closure in the ledger
    err := yfm.Ledger.DeFiLedger.RecordFarmingPoolClosure(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to log pool closure for pool ID %s: %v", poolID, err)
        return fmt.Errorf("failed to log pool closure in the ledger: %w", err)
    }

    log.Printf("[SUCCESS] Farming pool %s closed successfully", poolID)
    return nil
}



// ViewPoolDetails allows a user to view details of a specific farming pool
func (yfm *YieldFarmingManager) ViewPoolDetails(poolID string) (*FarmingPool, error) {
    log.Printf("[INFO] Request to view details for farming pool ID: %s", poolID)

    // Acquire lock for thread safety
    yfm.mu.Lock()
    defer yfm.mu.Unlock()

    // Validate input
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Retrieve the farming pool
    pool, exists := yfm.FarmingPools[poolID]
    if !exists {
        err := fmt.Errorf("farming pool %s not found", poolID)
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    log.Printf("[SUCCESS] Farming pool %s retrieved successfully", poolID)
    return pool, nil
}


// ViewStakingRecord allows a user to view their staking record by stake ID.
func (yfm *YieldFarmingManager) ViewStakingRecord(stakeID string) (*StakingRecord, error) {
    log.Printf("[INFO] Request to view staking record for stake ID: %s", stakeID)

    // Acquire lock for thread safety
    yfm.mu.Lock()
    defer yfm.mu.Unlock()

    // Validate input
    if stakeID == "" {
        err := fmt.Errorf("stakeID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Retrieve the staking record
    stakingRecord, exists := yfm.StakingRecords[stakeID]
    if !exists {
        err := fmt.Errorf("staking record %s not found", stakeID)
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    log.Printf("[SUCCESS] Staking record %s retrieved successfully", stakeID)
    return stakingRecord, nil
}


// YieldFarmAddLiquidity adds liquidity to a yield farming pool.
func YieldFarmAddLiquidity(poolID string, amount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Request to add liquidity: PoolID=%s, Amount=%.2f", poolID, amount)

    // Validate input
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if amount <= 0 {
        err := fmt.Errorf("amount must be greater than 0")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Attempt to add liquidity to the pool
    err := ledgerInstance.DeFiLedger.AddLiquidityToPool(poolID, amount)
    if err != nil {
        log.Printf("[ERROR] Failed to add liquidity to pool %s: %v", poolID, err)
        return fmt.Errorf("failed to add liquidity to pool %s: %w", poolID, err)
    }

    // Log the successful operation
    log.Printf("[SUCCESS] Liquidity of %.2f added to pool %s", amount, poolID)
    return nil
}


// YieldFarmRemoveLiquidity removes liquidity from a yield farming pool.
func YieldFarmRemoveLiquidity(poolID string, amount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Request to remove liquidity: PoolID=%s, Amount=%.2f", poolID, amount)

    // Validate input
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if amount <= 0 {
        err := fmt.Errorf("amount must be greater than 0")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Attempt to remove liquidity from the pool
    err := ledgerInstance.DeFiLedger.RemoveLiquidityFromPool(poolID, amount)
    if err != nil {
        log.Printf("[ERROR] Failed to remove liquidity from pool %s: %v", poolID, err)
        return fmt.Errorf("failed to remove liquidity from pool %s: %w", poolID, err)
    }

    // Log the successful operation
    log.Printf("[SUCCESS] Liquidity of %.2f removed from pool %s", amount, poolID)
    return nil
}


// YieldFarmStakeTokens allows a user to stake tokens in a yield farming pool.
func YieldFarmStakeTokens(poolID, userID string, amount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Request to stake tokens: PoolID=%s, UserID=%s, Amount=%.2f", poolID, userID, amount)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if userID == "" {
        err := fmt.Errorf("userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if amount <= 0 {
        err := fmt.Errorf("amount must be greater than 0")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Check if the pool exists and is active
    if !ledgerInstance.DeFiLedger.IsPoolActive(poolID) {
        err := fmt.Errorf("pool %s is not active or does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Stake tokens in the pool
    err := ledgerInstance.DeFiLedger.StakeTokensInPool(poolID, userID, amount)
    if err != nil {
        log.Printf("[ERROR] Failed to stake tokens in pool %s for user %s: %v", poolID, userID, err)
        return fmt.Errorf("failed to stake tokens in pool %s for user %s: %w", poolID, userID, err)
    }

    // Log successful operation
    log.Printf("[SUCCESS] User %s staked %.2f tokens in pool %s", userID, amount, poolID)
    return nil
}


// YieldFarmUnstakeTokens allows a user to unstake tokens from a yield farming pool.
func YieldFarmUnstakeTokens(poolID, userID string, amount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Request to unstake tokens: PoolID=%s, UserID=%s, Amount=%.2f", poolID, userID, amount)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if userID == "" {
        err := fmt.Errorf("userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if amount <= 0 {
        err := fmt.Errorf("amount must be greater than 0")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Check if the user has sufficient staked tokens
    stakedAmount, err := ledgerInstance.DeFiLedger.GetStakedTokens(poolID, userID)
    if err != nil {
        log.Printf("[ERROR] Failed to retrieve staked tokens for user %s in pool %s: %v", userID, poolID, err)
        return fmt.Errorf("failed to retrieve staked tokens for user %s in pool %s: %w", userID, poolID, err)
    }
    if amount > stakedAmount {
        err := fmt.Errorf("user %s has insufficient staked tokens in pool %s. Requested=%.2f, Available=%.2f", userID, poolID, amount, stakedAmount)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Unstake tokens from the pool
    err = ledgerInstance.DeFiLedger.UnstakeTokensFromPool(poolID, userID, amount)
    if err != nil {
        log.Printf("[ERROR] Failed to unstake tokens from pool %s for user %s: %v", poolID, userID, err)
        return fmt.Errorf("failed to unstake tokens from pool %s for user %s: %w", poolID, userID, err)
    }

    // Log successful operation
    log.Printf("[SUCCESS] User %s unstaked %.2f tokens from pool %s", userID, amount, poolID)
    return nil
}


// YieldFarmHarvestRewards allows a user to harvest rewards from a yield farming pool.
func YieldFarmHarvestRewards(poolID, userID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Request to harvest rewards: PoolID=%s, UserID=%s", poolID, userID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if userID == "" {
        err := fmt.Errorf("userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Check if the pool exists and is active
    if !ledgerInstance.DeFiLedger.IsPoolActive(poolID) {
        err := fmt.Errorf("pool %s is not active or does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Calculate rewards dynamically
    rewards, err := ledgerInstance.DeFiLedger.CalculateUserRewards(poolID, userID)
    if err != nil {
        log.Printf("[ERROR] Failed to calculate rewards for user %s in pool %s: %v", userID, poolID, err)
        return fmt.Errorf("failed to calculate rewards for user %s in pool %s: %w", userID, poolID, err)
    }

    if rewards <= 0 {
        log.Printf("[INFO] No rewards available for user %s in pool %s", userID, poolID)
        return nil // No rewards to harvest
    }

    // Harvest rewards
    err = ledgerInstance.DeFiLedger.HarvestYieldFarmRewards(poolID, userID)
    if err != nil {
        log.Printf("[ERROR] Failed to harvest rewards from pool %s for user %s: %v", poolID, userID, err)
        return fmt.Errorf("failed to harvest rewards from pool %s for user %s: %w", poolID, userID, err)
    }

    // Log successful operation
    log.Printf("[SUCCESS] User %s harvested %.2f rewards from pool %s", userID, rewards, poolID)
    return nil
}


// YieldFarmDistributeRewards distributes rewards to participants in a yield farming pool.
func YieldFarmDistributeRewards(poolID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Request to distribute rewards: PoolID=%s", poolID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Check if the pool exists and is active
    if !ledgerInstance.DeFiLedger.IsPoolActive(poolID) {
        err := fmt.Errorf("pool %s is not active or does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Distribute rewards
    err := ledgerInstance.DeFiLedger.DistributePoolRewards(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to distribute rewards for pool %s: %v", poolID, err)
        return fmt.Errorf("failed to distribute rewards for pool %s: %w", poolID, err)
    }

    // Log successful operation
    log.Printf("[SUCCESS] Rewards distributed for pool %s", poolID)
    return nil
}

// YieldFarmCalculateAPY calculates the Annual Percentage Yield (APY) for a yield farming pool.
func YieldFarmCalculateAPY(poolID string, ledgerInstance *ledger.Ledger) (float64, error) {
    log.Printf("[INFO] Request to calculate APY: PoolID=%s", poolID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Verify that the pool exists and is active
    if !ledgerInstance.DeFiLedger.IsPoolActive(poolID) {
        err := fmt.Errorf("pool %s is not active or does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Calculate APY
    apy, err := ledgerInstance.DeFiLedger.CalculateAPY(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to calculate APY for pool %s: %v", poolID, err)
        return 0, fmt.Errorf("failed to calculate APY for pool %s: %w", poolID, err)
    }

    // Log successful calculation
    log.Printf("[SUCCESS] APY for pool %s calculated successfully: %.2f%%", poolID, apy)
    return apy, nil
}


// YieldFarmLockFunds locks funds in a yield farming pool to prevent withdrawals.
func YieldFarmLockFunds(poolID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Request to lock funds: PoolID=%s", poolID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Verify that the pool exists and is active
    if !ledgerInstance.DeFiLedger.IsPoolActive(poolID) {
        err := fmt.Errorf("pool %s is not active or does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Lock funds in the pool
    err := ledgerInstance.DeFiLedger.LockPoolFunds(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to lock funds in pool %s: %v", poolID, err)
        return fmt.Errorf("failed to lock funds in pool %s: %w", poolID, err)
    }

    // Log successful operation
    log.Printf("[SUCCESS] Funds in pool %s locked successfully", poolID)
    return nil
}


// YieldFarmUnlockFunds unlocks previously locked funds in a yield farming pool.
func YieldFarmUnlockFunds(poolID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Request to unlock funds: PoolID=%s", poolID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Verify that the pool exists and funds are locked
    if !ledgerInstance.DeFiLedger.IsPoolLocked(poolID) {
        err := fmt.Errorf("funds in pool %s are not locked or pool does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Attempt to unlock funds
    err := ledgerInstance.DeFiLedger.UnlockPoolFunds(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to unlock funds in pool %s: %v", poolID, err)
        return fmt.Errorf("failed to unlock funds in pool %s: %w", poolID, err)
    }

    // Log success
    log.Printf("[SUCCESS] Funds in pool %s unlocked successfully", poolID)
    return nil
}


// YieldFarmMonitorPool enables real-time monitoring for a yield farming pool.
func YieldFarmMonitorPool(poolID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Enabling monitoring for pool: PoolID=%s", poolID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Verify that the pool exists
    if !ledgerInstance.DeFiLedger.DoesPoolExist(poolID) {
        err := fmt.Errorf("pool %s does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Enable monitoring (This could include tracking metrics, transactions, liquidity changes, etc.)
    err := ledgerInstance.DeFiLedger.StartMonitoringPool(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to enable monitoring for pool %s: %v", poolID, err)
        return fmt.Errorf("failed to enable monitoring for pool %s: %w", poolID, err)
    }

    // Log success
    log.Printf("[SUCCESS] Monitoring enabled for yield farming pool %s", poolID)
    return nil
}


// YieldFarmAuditPool performs an audit of the specified yield farming pool to verify its integrity.
func YieldFarmAuditPool(poolID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Starting audit for pool: PoolID=%s", poolID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Verify pool existence
    exists, err := ledgerInstance.DeFiLedger.DoesPoolExist(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to check existence of pool %s: %v", poolID, err)
        return fmt.Errorf("failed to verify pool existence: %w", err)
    }
    if !exists {
        err := fmt.Errorf("pool %s does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Perform the audit
    err = ledgerInstance.DeFiLedger.AuditYieldFarmPool(poolID)
    if err != nil {
        log.Printf("[ERROR] Audit failed for pool %s: %v", poolID, err)
        return fmt.Errorf("failed to audit pool %s: %w", poolID, err)
    }

    // Log success
    log.Printf("[SUCCESS] Audit completed successfully for pool %s", poolID)
    return nil
}


// YieldFarmLockPool locks a yield farming pool, preventing new stakes or changes.
func YieldFarmLockPool(poolID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Request to lock pool: PoolID=%s", poolID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Verify pool existence and current status
    exists, err := ledgerInstance.DeFiLedger.DoesPoolExist(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to verify existence of pool %s: %v", poolID, err)
        return fmt.Errorf("failed to verify pool existence: %w", err)
    }
    if !exists {
        err := fmt.Errorf("pool %s does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    isLocked, err := ledgerInstance.DeFiLedger.IsPoolLocked(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to check lock status for pool %s: %v", poolID, err)
        return fmt.Errorf("failed to check lock status: %w", err)
    }
    if isLocked {
        log.Printf("[INFO] Pool %s is already locked", poolID)
        return nil
    }

    // Lock the pool
    err = ledgerInstance.DeFiLedger.LockYieldFarmPool(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to lock pool %s: %v", poolID, err)
        return fmt.Errorf("failed to lock yield farming pool %s: %w", poolID, err)
    }

    // Log success
    log.Printf("[SUCCESS] Yield farming pool %s locked successfully", poolID)
    return nil
}


// YieldFarmUnlockPool unlocks a previously locked yield farming pool, allowing new stakes and operations.
func YieldFarmUnlockPool(poolID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Request to unlock yield farming pool: PoolID=%s", poolID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Check if the pool exists
    exists, err := ledgerInstance.DeFiLedger.DoesPoolExist(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to verify existence of pool %s: %v", poolID, err)
        return fmt.Errorf("failed to verify pool existence: %w", err)
    }
    if !exists {
        err := fmt.Errorf("pool %s does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Check if the pool is already unlocked
    isLocked, err := ledgerInstance.DeFiLedger.IsPoolLocked(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to retrieve lock status for pool %s: %v", poolID, err)
        return fmt.Errorf("failed to retrieve lock status: %w", err)
    }
    if !isLocked {
        log.Printf("[INFO] Pool %s is already unlocked", poolID)
        return nil
    }

    // Unlock the pool
    err = ledgerInstance.DeFiLedger.UnlockYieldFarmPool(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to unlock pool %s: %v", poolID, err)
        return fmt.Errorf("failed to unlock pool %s: %w", poolID, err)
    }

    // Log success
    log.Printf("[SUCCESS] Yield farming pool %s unlocked successfully", poolID)
    return nil
}


// YieldFarmTrackPerformance tracks the performance of a yield farming pool for analytics and optimization.
func YieldFarmTrackPerformance(poolID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating performance tracking for yield farming pool: PoolID=%s", poolID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Validate pool existence
    exists, err := ledgerInstance.DeFiLedger.DoesPoolExist(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to verify existence of pool %s: %v", poolID, err)
        return fmt.Errorf("failed to verify pool existence: %w", err)
    }
    if !exists {
        err := fmt.Errorf("pool %s does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Track pool performance
    err = ledgerInstance.DeFiLedger.TrackPoolPerformance(poolID)
    if err != nil {
        log.Printf("[ERROR] Performance tracking failed for pool %s: %v", poolID, err)
        return fmt.Errorf("failed to track performance for pool %s: %w", poolID, err)
    }

    // Log success
    log.Printf("[SUCCESS] Performance tracking enabled for yield farming pool %s", poolID)
    return nil
}


// YieldFarmFetchPerformanceMetrics retrieves detailed performance metrics for a yield farming pool.
func YieldFarmFetchPerformanceMetrics(poolID string, ledgerInstance *ledger.Ledger) (ledger.PoolPerformanceMetrics, error) {
    log.Printf("[INFO] Fetching performance metrics for yield farming pool: PoolID=%s", poolID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return ledger.PoolPerformanceMetrics{}, err
    }

    // Check if the pool exists
    exists, err := ledgerInstance.DeFiLedger.DoesPoolExist(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to verify existence of pool %s: %v", poolID, err)
        return ledger.PoolPerformanceMetrics{}, fmt.Errorf("failed to verify pool existence: %w", err)
    }
    if !exists {
        err := fmt.Errorf("pool %s does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return ledger.PoolPerformanceMetrics{}, err
    }

    // Fetch performance metrics
    metrics, err := ledgerInstance.DeFiLedger.GetPoolPerformanceMetrics(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch performance metrics for pool %s: %v", poolID, err)
        return ledger.PoolPerformanceMetrics{}, fmt.Errorf("failed to fetch performance metrics: %w", err)
    }

    // Log success
    log.Printf("[SUCCESS] Performance metrics fetched successfully for pool %s", poolID)
    return metrics, nil
}


// YieldFarmIncreaseRewards increases the rewards for a yield farming pool by a specified increment.
func YieldFarmIncreaseRewards(poolID string, increment float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Increasing rewards for yield farming pool: PoolID=%s, Increment=%.2f", poolID, increment)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if increment <= 0 {
        err := fmt.Errorf("increment must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Check if the pool exists
    exists, err := ledgerInstance.DeFiLedger.DoesPoolExist(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to verify existence of pool %s: %v", poolID, err)
        return fmt.Errorf("failed to verify pool existence: %w", err)
    }
    if !exists {
        err := fmt.Errorf("pool %s does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Increase rewards
    err = ledgerInstance.DeFiLedger.IncreasePoolRewards(poolID, increment)
    if err != nil {
        log.Printf("[ERROR] Failed to increase rewards for pool %s: %v", poolID, err)
        return fmt.Errorf("failed to increase rewards for pool %s: %w", poolID, err)
    }

    // Log success
    log.Printf("[SUCCESS] Rewards increased by %.2f for yield farming pool %s", increment, poolID)
    return nil
}


// YieldFarmDecreaseRewards decreases the rewards for a yield farming pool by a specified decrement.
func YieldFarmDecreaseRewards(poolID string, decrement float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating reward decrement for pool: PoolID=%s, Decrement=%.2f", poolID, decrement)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if decrement <= 0 {
        err := fmt.Errorf("decrement must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Check if the pool exists
    exists, err := ledgerInstance.DeFiLedger.DoesPoolExist(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to verify existence of pool %s: %v", poolID, err)
        return fmt.Errorf("failed to verify pool existence: %w", err)
    }
    if !exists {
        err := fmt.Errorf("pool %s does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Perform the reward decrement
    err = ledgerInstance.DeFiLedger.DecreasePoolRewards(poolID, decrement)
    if err != nil {
        log.Printf("[ERROR] Failed to decrease rewards for pool %s: %v", poolID, err)
        return fmt.Errorf("failed to decrease rewards for pool %s: %w", poolID, err)
    }

    // Log success
    log.Printf("[SUCCESS] Rewards decreased by %.2f for yield farming pool %s", decrement, poolID)
    return nil
}


// YieldFarmCompoundRewards compounds the rewards in a yield farming pool for reinvestment.
func YieldFarmCompoundRewards(poolID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating reward compounding for pool: PoolID=%s", poolID)

    // Input validation
    if poolID == "" {
        err := fmt.Errorf("poolID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Check if the pool exists
    exists, err := ledgerInstance.DeFiLedger.DoesPoolExist(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to verify existence of pool %s: %v", poolID, err)
        return fmt.Errorf("failed to verify pool existence: %w", err)
    }
    if !exists {
        err := fmt.Errorf("pool %s does not exist", poolID)
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Perform reward compounding
    err = ledgerInstance.DeFiLedger.CompoundPoolRewards(poolID)
    if err != nil {
        log.Printf("[ERROR] Failed to compound rewards for pool %s: %v", poolID, err)
        return fmt.Errorf("failed to compound rewards for pool %s: %w", poolID, err)
    }

    // Log success
    log.Printf("[SUCCESS] Rewards compounded successfully for yield farming pool %s", poolID)
    return nil
}

