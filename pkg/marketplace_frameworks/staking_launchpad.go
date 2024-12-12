package marketplace

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewStakingLaunchpad initializes a new staking launchpad
func NewStakingLaunchpad(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.StakingLaunchpad {
	return &common.StakingLaunchpad{
		Pools:            make(map[string]*common.StakingPool),
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// CreateStakingPool allows a project owner to create a new staking pool
func (sl *common.StakingLaunchpad) CreateStakingPool(poolID, projectName, tokenAddress, owner string, rewardRate float64, startTime, endTime time.Time) (*common.StakingPool, error) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	// Ensure the end time is after the start time
	if endTime.Before(startTime) {
		return nil, errors.New("end time must be after start time")
	}

	// Encrypt staking pool data
	poolData := fmt.Sprintf("PoolID: %s, Project: %s, Owner: %s, Token: %s", poolID, projectName, owner, tokenAddress)
	encryptedData, err := sl.EncryptionService.EncryptData([]byte(poolData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt pool data: %v", err)
	}

	// Create the staking pool
	pool := &common.StakingPool{
		PoolID:       poolID,
		ProjectName:  projectName,
		TokenAddress: tokenAddress,
		Owner:        owner,
		RewardRate:   rewardRate,
		StartTime:    startTime,
		EndTime:      endTime,
		IsActive:     true,
		Participants: make(map[string]float64),
	}

	// Add the pool to the launchpad
	sl.Pools[poolID] = pool

	// Log the staking pool creation in the ledger
	err = sl.Ledger.RecordStakingPoolCreation(poolID, projectName, owner, rewardRate, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to log staking pool creation: %v", err)
	}

	fmt.Printf("Staking pool for project %s created by %s with reward rate %f\n", projectName, owner, rewardRate)
	return pool, nil
}

// StakeTokens allows a user to stake tokens in a staking pool
func (sl *common.StakingLaunchpad) StakeTokens(poolID, staker string, amount float64) error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	// Retrieve the staking pool
	pool, exists := sl.Pools[poolID]
	if !exists || !pool.IsActive {
		return fmt.Errorf("staking pool %s is not active", poolID)
	}

	// Update the participant's stake and the total staked amount
	pool.Participants[staker] += amount
	pool.StakedAmount += amount

	// Log the staking transaction in the ledger
	err := sl.Ledger.RecordStakeTransaction(poolID, staker, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log staking transaction: %v", err)
	}

	fmt.Printf("%f tokens staked by %s in pool %s\n", amount, staker, poolID)
	return nil
}

// UnstakeTokens allows a user to withdraw staked tokens from the pool
func (sl *common.StakingLaunchpad) UnstakeTokens(poolID, staker string) (float64, error) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	// Retrieve the staking pool
	pool, exists := sl.Pools[poolID]
	if !exists {
		return 0, fmt.Errorf("staking pool %s not found", poolID)
	}

	// Ensure the user has staked tokens
	stakedAmount, participantExists := pool.Participants[staker]
	if !participantExists || stakedAmount == 0 {
		return 0, fmt.Errorf("no tokens staked by %s in pool %s", staker, poolID)
	}

	// Remove the user's stake
	delete(pool.Participants, staker)
	pool.StakedAmount -= stakedAmount

	// Log the unstaking transaction in the ledger
	err := sl.Ledger.RecordUnstakeTransaction(poolID, staker, stakedAmount, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to log unstaking transaction: %v", err)
	}

	fmt.Printf("%f tokens unstaked by %s from pool %s\n", stakedAmount, staker, poolID)
	return stakedAmount, nil
}

// DistributeRewards distributes rewards to all participants in a staking pool
func (sl *common.StakingLaunchpad) DistributeRewards(poolID string) error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	// Retrieve the staking pool
	pool, exists := sl.Pools[poolID]
	if !exists {
		return fmt.Errorf("staking pool %s not found", poolID)
	}

	// Ensure the pool is active and within the staking period
	if time.Now().Before(pool.StartTime) || time.Now().After(pool.EndTime) {
		return fmt.Errorf("staking pool %s is not within the reward distribution period", poolID)
	}

	// Distribute rewards to each participant
	for staker, stakedAmount := range pool.Participants {
		reward := stakedAmount * pool.RewardRate
		// Log the reward distribution
		err := sl.Ledger.RecordRewardDistribution(poolID, staker, reward, time.Now())
		if err != nil {
			return fmt.Errorf("failed to log reward distribution: %v", err)
		}
		fmt.Printf("Distributed %f rewards to %s in pool %s\n", reward, staker, poolID)
	}

	return nil
}

// EndStakingPool ends a staking pool and prevents further staking
func (sl *common.StakingLaunchpad) EndStakingPool(poolID string) error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	// Retrieve the staking pool
	pool, exists := sl.Pools[poolID]
	if !exists {
		return fmt.Errorf("staking pool %s not found", poolID)
	}

	// Mark the pool as inactive
	pool.IsActive = false

	// Log the pool closure in the ledger
	err := sl.Ledger.RecordStakingPoolClosure(poolID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log staking pool closure: %v", err)
	}

	fmt.Printf("Staking pool %s has been closed\n", poolID)
	return nil
}

// generateUniqueID creates a cryptographically secure unique ID
func generateUniqueID() string {
	id := make([]byte, 16)
	rand.Read(id)
	return hex.EncodeToString(id)
}
