package dao

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/tokens/syn900"
)


// NewStakingManager initializes a new StakingManager.
func NewDAOStakingManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption, syn900Verifier *syn900.Verifier) *StakingManager {
	return &StakingManager{
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
		Syn900Verifier:    syn900Verifier,
		GovernanceStakes:  make(map[string]*GovernanceStakingSystem),
	}
}

// InitializeGovernanceStaking sets up the governance staking system for a DAO.
func (sm *StakingManager) InitializeGovernanceStaking(daoID string, minStakeAmount float64, stakingDuration time.Duration) (*GovernanceStakingSystem, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Ensure the DAO doesn't already have a staking system
	if _, exists := sm.GovernanceStakes[daoID]; exists {
		return nil, errors.New("governance staking system already exists for this DAO")
	}

	// Initialize the governance staking system
	stakingSystem := &GovernanceStakingSystem{
		DAOID:             daoID,
		TotalStakedTokens: 0,
		StakingRecords:    make(map[string]*GovernanceStake),
		MinStakeAmount:    minStakeAmount,
		StakingDuration:   stakingDuration,
	}

	// Record the staking system initialization in the ledger
	if err := sm.Ledger.DAOLedger.RecordGovernanceStakingInitialization(stakingSystem); err != nil {
		return nil, fmt.Errorf("failed to record governance staking initialization in ledger: %v", err)
	}

	sm.GovernanceStakes[daoID] = stakingSystem
	fmt.Printf("Governance staking system initialized for DAO %s\n", daoID)
	return stakingSystem, nil
}

// StakeTokensForGovernance allows a user to stake tokens for governance in the DAO.
func (sm *StakingManager) StakeTokensForGovernance(daoID, stakerWallet string, amount float64) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the governance staking system for the DAO
	stakingSystem, exists := sm.GovernanceStakes[daoID]
	if !exists {
		return errors.New("governance staking system not found for this DAO")
	}

	// Ensure the amount to stake is above the minimum threshold
	if amount < stakingSystem.MinStakeAmount {
		return errors.New("staking amount is below the minimum required for governance participation")
	}

	// Create or update the staking record for the user
	stakeRecord, exists := stakingSystem.StakingRecords[stakerWallet]
	if exists && stakeRecord.IsActive {
		return errors.New("user already has an active stake for governance")
	}

	votingPower := sm.calculateVotingPower(amount, stakingSystem.TotalStakedTokens)

	// Create a new staking record
	stakeRecord = &GovernanceStake{
		StakerWallet:   stakerWallet,
		Amount:         amount,
		VotingPower:    votingPower,
		StakeTimestamp: time.Now(),
		IsActive:       true,
	}

	// Add the record to the staking system
	stakingSystem.StakingRecords[stakerWallet] = stakeRecord
	stakingSystem.TotalStakedTokens += amount

	// Record the staking transaction in the ledger
	err := sm.Ledger.DAOLedger.RecordStakeTransaction(daoID, stakerWallet, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to record staking transaction in ledger: %v", err)
	}

	fmt.Printf("User %s successfully staked %f tokens for governance in DAO %s\n", stakerWallet, amount, daoID)
	return nil
}

// UnstakeTokensForGovernance allows a user to withdraw their staked tokens after the lock-in period.
func (sm *StakingManager) UnstakeTokensForGovernance(daoID, stakerWallet string) (float64, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the governance staking system for the DAO
	stakingSystem, exists := sm.GovernanceStakes[daoID]
	if !exists {
		return 0, errors.New("governance staking system not found for this DAO")
	}

	// Retrieve the user's staking record
	stakeRecord, exists := stakingSystem.StakingRecords[stakerWallet]
	if !exists || !stakeRecord.IsActive {
		return 0, errors.New("no active governance stake found for this user")
	}

	// Check if the lock-in period has expired
	if time.Now().Before(stakeRecord.StakeTimestamp.Add(stakingSystem.StakingDuration)) {
		return 0, errors.New("staked tokens are still locked and cannot be unstaked")
	}

	// Remove the stake and update the total staked tokens
	unstakeAmount := stakeRecord.Amount
	stakeRecord.IsActive = false
	stakingSystem.TotalStakedTokens -= unstakeAmount

	// Record the unstaking transaction in the ledger
	err := sm.Ledger.DAOLedger.RecordUnstakeTransaction(daoID, stakerWallet, unstakeAmount, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to record unstaking transaction in ledger: %v", err)
	}

	fmt.Printf("User %s successfully unstaked %f tokens from governance in DAO %s\n", stakerWallet, unstakeAmount, daoID)
	return unstakeAmount, nil
}

// Calculate the voting power based on the amount staked.
func (sm *StakingManager) calculateVotingPower(stakedAmount, totalStaked float64) float64 {
	if totalStaked == 0 {
		return stakedAmount
	}
	return stakedAmount / totalStaked * 100 // Voting power as a percentage
}

// GetVotingPower retrieves the voting power of a specific user in a DAO.
func (sm *StakingManager) GetVotingPower(daoID, stakerWallet string) (float64, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the governance staking system for the DAO
	stakingSystem, exists := sm.GovernanceStakes[daoID]
	if !exists {
		return 0, errors.New("governance staking system not found for this DAO")
	}

	// Retrieve the user's staking record
	stakeRecord, exists := stakingSystem.StakingRecords[stakerWallet]
	if !exists || !stakeRecord.IsActive {
		return 0, errors.New("no active governance stake found for this user")
	}

	return stakeRecord.VotingPower, nil
}
