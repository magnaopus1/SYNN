package defi

import (
	"fmt"
	"log"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"time"
)

// StakingCreateProgram creates a new staking program with specified parameters.
func StakingCreateProgram(programID string, rewardRate float64, minStakeAmount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Creating staking program. Program ID: %s, Reward Rate: %.2f, Min Stake: %.2f", programID, rewardRate, minStakeAmount)

    // Step 1: Input validation
    if programID == "" {
        err := fmt.Errorf("programID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if rewardRate <= 0 {
        err := fmt.Errorf("rewardRate must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if minStakeAmount <= 0 {
        err := fmt.Errorf("minStakeAmount must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Create staking program object
    program := ledger.StakingProgram{
        ProgramID:    programID,
        RewardRate:   rewardRate,
        MinStake:     minStakeAmount,
        Status:       "Active",
        LockedTokens: make(map[string]float64),
    }

    // Step 3: Record staking program in ledger
    startTime := time.Now()
    if err := ledgerInstance.DeFiLedger.CreateStakingProgram(program); err != nil {
        log.Printf("[ERROR] Failed to create staking program %s: %v", programID, err)
        return fmt.Errorf("failed to create staking program %s: %w", programID, err)
    }

    // Step 4: Log success
    log.Printf("[INFO] Staking program created successfully. Program ID: %s, Duration: %v", programID, time.Since(startTime))
    return nil
}


// StakingStakeTokens stakes tokens for a user in a specified staking program.
func StakingStakeTokens(programID, userID string, amount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Staking tokens. Program ID: %s, User ID: %s, Amount: %.2f", programID, userID, amount)

    // Step 1: Input validation
    if programID == "" || userID == "" {
        err := fmt.Errorf("programID and userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if amount <= 0 {
        err := fmt.Errorf("amount must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt user ID
    encryptedUserID, err := encryption.EncryptString(userID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt user ID: %v", err)
        return fmt.Errorf("failed to encrypt userID: %w", err)
    }

    // Step 3: Stake tokens in the ledger
    startTime := time.Now()
    if err := ledgerInstance.DeFiLedger.StakeTokens(programID, encryptedUserID, amount); err != nil {
        log.Printf("[ERROR] Failed to stake tokens for program %s: %v", programID, err)
        return fmt.Errorf("failed to stake tokens for program %s: %w", programID, err)
    }

    // Step 4: Log success
    log.Printf("[INFO] Tokens staked successfully. Program ID: %s, User ID: %s, Amount: %.2f, Duration: %v", programID, userID, amount, time.Since(startTime))
    return nil
}


// StakingUnstakeTokens unstakes tokens for a user in a specified staking program.
func StakingUnstakeTokens(programID, userID string, amount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating unstaking process. Program ID: %s, User ID: %s, Amount: %.2f", programID, userID, amount)

    // Step 1: Input validation
    if programID == "" || userID == "" {
        err := fmt.Errorf("programID and userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if amount <= 0 {
        err := fmt.Errorf("amount must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt user ID for security
    encryptedUserID, err := encryption.EncryptString(userID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt user ID: %v", err)
        return fmt.Errorf("failed to encrypt userID: %w", err)
    }

    // Step 3: Unstake tokens
    startTime := time.Now()
    if err := ledgerInstance.DeFiLedger.UnstakeTokens(programID, encryptedUserID, amount); err != nil {
        log.Printf("[ERROR] Failed to unstake tokens for program %s: %v", programID, err)
        return fmt.Errorf("failed to unstake tokens for program %s: %w", programID, err)
    }

    // Step 4: Log success and return
    log.Printf("[INFO] Tokens unstaked successfully. Program ID: %s, User ID: %s, Amount: %.2f, Duration: %v", programID, userID, amount, time.Since(startTime))
    return nil
}


// StakingCalculateRewards calculates rewards for a user in a specified staking program.
func StakingCalculateRewards(programID, userID string, ledgerInstance *ledger.Ledger) (float64, error) {
    log.Printf("[INFO] Calculating staking rewards. Program ID: %s, User ID: %s", programID, userID)

    // Step 1: Input validation
    if programID == "" || userID == "" {
        err := fmt.Errorf("programID and userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Step 2: Encrypt user ID for security
    encryptedUserID, err := encryption.EncryptString(userID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt user ID: %v", err)
        return 0, fmt.Errorf("failed to encrypt userID: %w", err)
    }

    // Step 3: Calculate rewards
    startTime := time.Now()
    rewards, err := ledgerInstance.DeFiLedger.CalculateStakingRewards(programID, encryptedUserID)
    if err != nil {
        log.Printf("[ERROR] Failed to calculate rewards for program %s: %v", programID, err)
        return 0, fmt.Errorf("failed to calculate rewards for program %s: %w", programID, err)
    }

    // Step 4: Log success and return rewards
    log.Printf("[INFO] Rewards calculated successfully. Program ID: %s, User ID: %s, Rewards: %.2f, Duration: %v", programID, userID, rewards, time.Since(startTime))
    return rewards, nil
}


// StakingDistributeRewards distributes rewards for a specified staking program.
func StakingDistributeRewards(programID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating reward distribution for program: %s", programID)

    // Step 1: Input validation
    if programID == "" {
        err := fmt.Errorf("programID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Distribute rewards
    startTime := time.Now()
    if err := ledgerInstance.DeFiLedger.DistributeStakingRewards(programID); err != nil {
        log.Printf("[ERROR] Failed to distribute rewards for program %s: %v", programID, err)
        return fmt.Errorf("failed to distribute rewards for program %s: %w", programID, err)
    }

    // Step 3: Log success and return
    log.Printf("[INFO] Rewards distributed successfully for program %s. Duration: %v", programID, time.Since(startTime))
    return nil
}


// StakingLockTokens locks tokens for a user in a specified staking program.
func StakingLockTokens(programID, userID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating token lock for program: %s, User ID: %s", programID, userID)

    // Step 1: Input validation
    if programID == "" || userID == "" {
        err := fmt.Errorf("programID and userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt user ID for security
    encryptedUserID, err := encryption.EncryptString(userID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt user ID: %v", err)
        return fmt.Errorf("failed to encrypt userID: %w", err)
    }

    // Step 3: Lock tokens
    startTime := time.Now()
    if err := ledgerInstance.DeFiLedger.LockTokens(programID, encryptedUserID); err != nil {
        log.Printf("[ERROR] Failed to lock tokens for program %s: %v", programID, err)
        return fmt.Errorf("failed to lock tokens for program %s: %w", programID, err)
    }

    // Step 4: Log success and return
    log.Printf("[INFO] Tokens locked successfully for program %s, User ID: %s. Duration: %v", programID, userID, time.Since(startTime))
    return nil
}



// StakingUnlockTokens unlocks staked tokens for a user in a specified staking program.
func StakingUnlockTokens(programID, userID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating token unlock for Program ID: %s, User ID: %s", programID, userID)

    // Step 1: Input validation
    if programID == "" || userID == "" {
        err := fmt.Errorf("programID and userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt user ID for secure processing
    encryptedUserID, err := encryption.EncryptString(userID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt userID: %v", err)
        return fmt.Errorf("failed to encrypt userID: %w", err)
    }

    // Step 3: Unlock tokens
    startTime := time.Now()
    if err := ledgerInstance.DeFiLedger.UnlockTokens(programID, encryptedUserID); err != nil {
        log.Printf("[ERROR] Failed to unlock tokens for Program ID: %s, User ID: %s. Error: %v", programID, userID, err)
        return fmt.Errorf("failed to unlock tokens for program %s: %w", programID, err)
    }

    // Step 4: Log success and return
    log.Printf("[INFO] Tokens unlocked successfully for Program ID: %s, User ID: %s. Duration: %v", programID, userID, time.Since(startTime))
    return nil
}


// StakingAudit audits a specified staking program to ensure compliance and correctness.
func StakingAudit(programID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Starting audit for Staking Program ID: %s", programID)

    // Step 1: Input validation
    if programID == "" {
        err := fmt.Errorf("programID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Perform audit
    startTime := time.Now()
    if err := ledgerInstance.DeFiLedger.AuditStakingProgram(programID); err != nil {
        log.Printf("[ERROR] Failed to audit Staking Program ID: %s. Error: %v", programID, err)
        return fmt.Errorf("failed to audit staking program %s: %w", programID, err)
    }

    // Step 3: Log success and return
    log.Printf("[INFO] Audit completed successfully for Staking Program ID: %s. Duration: %v", programID, time.Since(startTime))
    return nil
}


// StakingMonitor monitors the status of a specified staking program.
func StakingMonitor(programID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Starting monitoring for Staking Program ID: %s", programID)

    // Step 1: Input validation
    if programID == "" {
        err := fmt.Errorf("programID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Start monitoring
    startTime := time.Now()
    if err := ledgerInstance.DeFiLedger.MonitorStakingProgram(programID); err != nil {
        log.Printf("[ERROR] Failed to monitor Staking Program ID: %s. Error: %v", programID, err)
        return fmt.Errorf("failed to monitor staking program %s: %w", programID, err)
    }

    // Step 3: Log success
    log.Printf("[INFO] Staking Program ID: %s monitored successfully. Duration: %v", programID, time.Since(startTime))
    return nil
}


// StakingSnapshot takes a snapshot of a specified staking program.
func StakingSnapshot(programID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating snapshot for Staking Program ID: %s", programID)

    // Step 1: Input validation
    if programID == "" {
        err := fmt.Errorf("programID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Take snapshot
    startTime := time.Now()
    if err := ledgerInstance.DeFiLedger.TakeStakingSnapshot(programID); err != nil {
        log.Printf("[ERROR] Failed to take snapshot for Staking Program ID: %s. Error: %v", programID, err)
        return fmt.Errorf("failed to take snapshot for program %s: %w", programID, err)
    }

    // Step 3: Log success
    log.Printf("[INFO] Snapshot taken successfully for Staking Program ID: %s. Duration: %v", programID, time.Since(startTime))
    return nil
}


// StakingFetchStakeAmount retrieves the staked amount for a user in a specified staking program.
func StakingFetchStakeAmount(programID, userID string, ledgerInstance *ledger.Ledger) (float64, error) {
    log.Printf("[INFO] Fetching staked amount for Program ID: %s, User ID: %s", programID, userID)

    // Step 1: Input validation
    if programID == "" || userID == "" {
        err := fmt.Errorf("programID and userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Step 2: Encrypt user ID
    encryptedUserID, err := encryption.EncryptString(userID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt User ID: %s. Error: %v", userID, err)
        return 0, fmt.Errorf("failed to encrypt userID: %w", err)
    }

    // Step 3: Fetch stake amount
    stakeAmount, err := ledgerInstance.DeFiLedger.FetchStakeAmount(programID, encryptedUserID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch stake amount for Program ID: %s, User ID: %s. Error: %v", programID, userID, err)
        return 0, fmt.Errorf("failed to fetch stake amount for program %s: %w", programID, err)
    }

    // Step 4: Log success and return
    log.Printf("[SUCCESS] Staked amount fetched successfully. Program ID: %s, User ID: %s, Amount: %.2f", programID, userID, stakeAmount)
    return stakeAmount, nil
}


// StakingUpdateStakeAmount updates the staked amount for a user in a specified staking program.
func StakingUpdateStakeAmount(programID, userID string, newAmount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Updating staked amount for Program ID: %s, User ID: %s", programID, userID)

    // Step 1: Input validation
    if programID == "" || userID == "" {
        err := fmt.Errorf("programID and userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if newAmount <= 0 {
        err := fmt.Errorf("newAmount must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt user ID
    encryptedUserID, err := encryption.EncryptString(userID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt User ID: %s. Error: %v", userID, err)
        return fmt.Errorf("failed to encrypt userID: %w", err)
    }

    // Step 3: Update stake amount
    if err := ledgerInstance.DeFiLedger.UpdateStakeAmount(programID, encryptedUserID, newAmount); err != nil {
        log.Printf("[ERROR] Failed to update stake amount for Program ID: %s, User ID: %s. Error: %v", programID, userID, err)
        return fmt.Errorf("failed to update stake amount for program %s: %w", programID, err)
    }

    // Step 4: Log success
    log.Printf("[SUCCESS] Stake amount updated successfully. Program ID: %s, User ID: %s, New Amount: %.2f", programID, userID, newAmount)
    return nil
}


// StakingFetchRewardHistory retrieves the reward history for a user in a specified staking program.
func StakingFetchRewardHistory(programID, userID string, ledgerInstance *ledger.Ledger) ([]ledger.RewardRecord, error) {
    log.Printf("[INFO] Fetching reward history for Program ID: %s, User ID: %s", programID, userID)

    // Step 1: Input validation
    if programID == "" || userID == "" {
        err := fmt.Errorf("programID and userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return nil, err
    }

    // Step 2: Encrypt user ID
    encryptedUserID, err := encryption.EncryptString(userID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt User ID: %s. Error: %v", userID, err)
        return nil, fmt.Errorf("failed to encrypt userID: %w", err)
    }

    // Step 3: Fetch reward history
    history, err := ledgerInstance.DeFiLedger.FetchRewardHistory(programID, encryptedUserID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch reward history for Program ID: %s, User ID: %s. Error: %v", programID, userID, err)
        return nil, fmt.Errorf("failed to fetch reward history for program %s: %w", programID, err)
    }

    // Step 4: Log success
    log.Printf("[SUCCESS] Reward history fetched successfully. Program ID: %s, User ID: %s", programID, userID)
    return history, nil
}


// StakingDistributeBonuses distributes bonuses to all eligible participants in a staking program.
func StakingDistributeBonuses(programID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Distributing bonuses for Program ID: %s", programID)

    // Step 1: Input validation
    if programID == "" {
        err := fmt.Errorf("programID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Distribute bonuses
    if err := ledgerInstance.DeFiLedger.DistributeStakingBonuses(programID); err != nil {
        log.Printf("[ERROR] Failed to distribute bonuses for Program ID: %s. Error: %v", programID, err)
        return fmt.Errorf("failed to distribute bonuses for program %s: %w", programID, err)
    }

    // Step 3: Log success
    log.Printf("[SUCCESS] Bonuses distributed successfully for Program ID: %s", programID)
    return nil
}


// StakingReclaimRewards allows a user to reclaim unclaimed rewards from a staking program.
func StakingReclaimRewards(programID, userID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Initiating reward reclamation. Program ID: %s, User ID: %s", programID, userID)

    // Step 1: Input Validation
    if programID == "" || userID == "" {
        err := fmt.Errorf("programID and userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt User ID
    encryptedUserID, err := encryption.EncryptString(userID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt User ID: %s. Error: %v", userID, err)
        return fmt.Errorf("failed to encrypt userID: %w", err)
    }

    // Step 3: Reclaim Rewards
    if err := ledgerInstance.DeFiLedger.ReclaimStakingRewards(programID, encryptedUserID); err != nil {
        log.Printf("[ERROR] Failed to reclaim rewards for Program ID: %s, User ID: %s. Error: %v", programID, userID, err)
        return fmt.Errorf("failed to reclaim rewards for program %s: %w", programID, err)
    }

    // Step 4: Log Success
    log.Printf("[SUCCESS] Rewards reclaimed successfully. Program ID: %s, User ID: %s", programID, userID)
    return nil
}


// StakingSetMinimumStake sets the minimum stake amount for a staking program.
func StakingSetMinimumStake(programID string, minAmount float64, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Setting minimum stake for Program ID: %s. Minimum Amount: %.2f", programID, minAmount)

    // Step 1: Input Validation
    if programID == "" {
        err := fmt.Errorf("programID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if minAmount <= 0 {
        err := fmt.Errorf("minAmount must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Update Minimum Stake in Ledger
    if err := ledgerInstance.DeFiLedger.SetMinimumStake(programID, minAmount); err != nil {
        log.Printf("[ERROR] Failed to set minimum stake for Program ID: %s. Error: %v", programID, err)
        return fmt.Errorf("failed to set minimum stake for program %s: %w", programID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Minimum stake set successfully. Program ID: %s, Minimum Amount: %.2f", programID, minAmount)
    return nil
}


// StakingFetchMinimumStake retrieves the minimum stake amount for a staking program.
func StakingFetchMinimumStake(programID string, ledgerInstance *ledger.Ledger) (float64, error) {
    log.Printf("[INFO] Fetching minimum stake for Program ID: %s", programID)

    // Step 1: Input Validation
    if programID == "" {
        err := fmt.Errorf("programID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Step 2: Fetch Minimum Stake
    minStake, err := ledgerInstance.DeFiLedger.FetchMinimumStake(programID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch minimum stake for Program ID: %s. Error: %v", programID, err)
        return 0, fmt.Errorf("failed to fetch minimum stake for program %s: %w", programID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Minimum stake fetched successfully. Program ID: %s, Minimum Amount: %.2f", programID, minStake)
    return minStake, nil
}


// StakingSetRewardDistribution sets the reward distribution frequency for a staking program.
func StakingSetRewardDistribution(programID string, frequency time.Duration, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Setting reward distribution frequency. Program ID: %s, Frequency: %v", programID, frequency)

    // Step 1: Input Validation
    if programID == "" {
        err := fmt.Errorf("programID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if frequency <= 0 {
        err := fmt.Errorf("frequency must be greater than zero")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Update Reward Distribution Frequency
    if err := ledgerInstance.DeFiLedger.SetRewardDistributionFrequency(programID, frequency); err != nil {
        log.Printf("[ERROR] Failed to set reward distribution frequency for Program ID: %s. Error: %v", programID, err)
        return fmt.Errorf("failed to set reward distribution frequency for program %s: %w", programID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Reward distribution frequency set successfully. Program ID: %s, Frequency: %v", programID, frequency)
    return nil
}


// StakingFetchRewardDistribution retrieves the reward distribution frequency for a staking program.
func StakingFetchRewardDistribution(programID string, ledgerInstance *ledger.Ledger) (time.Duration, error) {
    log.Printf("[INFO] Fetching reward distribution frequency for Program ID: %s", programID)

    // Step 1: Input Validation
    if programID == "" {
        err := fmt.Errorf("programID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    // Step 2: Fetch Reward Distribution Frequency
    distributionFrequency, err := ledgerInstance.DeFiLedger.FetchRewardDistributionFrequency(programID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch reward distribution frequency for Program ID: %s. Error: %v", programID, err)
        return 0, fmt.Errorf("failed to fetch reward distribution frequency for program %s: %w", programID, err)
    }

    // Step 3: Log Success
    log.Printf("[SUCCESS] Reward distribution frequency fetched successfully. Program ID: %s, Frequency: %v", programID, distributionFrequency)
    return distributionFrequency, nil
}


// StakingAutoReinvest enables auto-reinvestment for a user in a staking program.
func StakingAutoReinvest(programID, userID string, ledgerInstance *ledger.Ledger) error {
    log.Printf("[INFO] Enabling auto-reinvestment for Program ID: %s, User ID: %s", programID, userID)

    // Step 1: Input Validation
    if programID == "" {
        err := fmt.Errorf("programID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }
    if userID == "" {
        err := fmt.Errorf("userID cannot be empty")
        log.Printf("[ERROR] %v", err)
        return err
    }

    // Step 2: Encrypt User ID
    encryptedUserID, err := encryption.EncryptString(userID)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt User ID: %s. Error: %v", userID, err)
        return fmt.Errorf("failed to encrypt userID: %w", err)
    }

    // Step 3: Enable Auto-Reinvestment
    if err := ledgerInstance.DeFiLedger.EnableAutoReinvestment(programID, encryptedUserID); err != nil {
        log.Printf("[ERROR] Failed to enable auto-reinvestment for Program ID: %s, User ID: %s. Error: %v", programID, userID, err)
        return fmt.Errorf("failed to enable auto-reinvestment for program %s: %w", programID, err)
    }

    // Step 4: Log Success
    log.Printf("[SUCCESS] Auto-reinvestment enabled successfully. Program ID: %s, User ID: %s", programID, userID)
    return nil
}

