package consensus

import (
	"fmt"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// ConsensusEnableValidatorBans enables the banning mechanism for validators.
func ConsensusEnableValidatorBans(ledgerInstance *ledger.Ledger) error {
    // Validate the ledger instance
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Attempting to enable validator bans.")

    // Enable the banning mechanism via the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.EnableValidatorBans(); err != nil {
        log.Printf("[ERROR] Failed to enable validator bans: %v", err)
        return fmt.Errorf("failed to enable validator bans: %w", err)
    }

    log.Printf("[SUCCESS] Validator banning mechanism enabled successfully.")
    return nil
}


// ConsensusDisableValidatorBans disables the banning mechanism for validators.
func ConsensusDisableValidatorBans(ledgerInstance *ledger.Ledger) error {
    // Validate the ledger instance
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Attempting to disable validator bans.")

    // Disable the banning mechanism via the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.DisableValidatorBans(); err != nil {
        log.Printf("[ERROR] Failed to disable validator bans: %v", err)
        return fmt.Errorf("failed to disable validator bans: %w", err)
    }

    log.Printf("[SUCCESS] Validator banning mechanism disabled successfully.")
    return nil
}


// ConsensusBanValidator bans a validator and records the ban in the ledger.
func ConsensusBanValidator(validatorID string, reason string, ledgerInstance *ledger.Ledger) error {
    // Input validation
    if validatorID == "" {
        return fmt.Errorf("validator ID cannot be empty")
    }
    if reason == "" {
        return fmt.Errorf("ban reason cannot be empty")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Initiating ban for validator ID: %s", validatorID)

    // Initialize encryption for secure recording
    encryptionInstance, err := common.NewEncryption(256)
    if err != nil {
        log.Printf("[ERROR] Failed to create encryption instance: %v", err)
        return fmt.Errorf("failed to create encryption instance: %w", err)
    }

    // Encrypt the reason for the ban
    encryptedReason, err := encryptionInstance.EncryptData("AES", []byte(reason), common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt ban reason: %v", err)
        return fmt.Errorf("failed to encrypt ban reason: %w", err)
    }

    // Record the ban in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.BanValidator(validatorID, string(encryptedReason)); err != nil {
        log.Printf("[ERROR] Failed to ban validator %s: %v", validatorID, err)
        return fmt.Errorf("failed to ban validator %s: %w", validatorID, err)
    }

    log.Printf("[SUCCESS] Validator %s banned successfully.", validatorID)
    return nil
}


// ConsensusUnbanValidator removes a ban on a validator.
func ConsensusUnbanValidator(validatorID string, ledgerInstance *ledger.Ledger) error {
    // Input validation
    if validatorID == "" {
        return fmt.Errorf("validator ID cannot be empty")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Initiating unban for validator ID: %s", validatorID)

    // Remove the ban in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.UnbanValidator(validatorID); err != nil {
        log.Printf("[ERROR] Failed to unban validator %s: %v", validatorID, err)
        return fmt.Errorf("failed to unban validator %s: %w", validatorID, err)
    }

    log.Printf("[SUCCESS] Validator %s unbanned successfully.", validatorID)
    return nil
}


// ConsensusGetBannedValidators fetches the list of currently banned validators.
func ConsensusGetBannedValidators(ledgerInstance *ledger.Ledger) ([]string, error) {
    // Input validation
    if ledgerInstance == nil {
        return nil, fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Fetching list of banned validators.")

    // Fetch banned validators from the ledger
    bannedValidators, err := ledgerInstance.BlockchainConsensusCoinLedger.FetchBannedValidators()
    if err != nil {
        log.Printf("[ERROR] Failed to fetch banned validators: %v", err)
        return nil, fmt.Errorf("failed to fetch banned validators: %w", err)
    }

    log.Printf("[SUCCESS] Retrieved %d banned validators.", len(bannedValidators))
    return bannedValidators, nil
}


// ConsensusAuditPunishments audits the punishment records for validators.
func ConsensusAuditPunishments(ledgerInstance *ledger.Ledger) error {
    // Input validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Initiating audit of validator punishment records.")

    // Perform the punishment audit via the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.AuditValidatorPunishments(); err != nil {
        log.Printf("[ERROR] Failed to audit validator punishments: %v", err)
        return fmt.Errorf("failed to audit validator punishments: %w", err)
    }

    log.Printf("[SUCCESS] Validator punishment records audited successfully.")
    return nil
}

// ConsensusFetchPunishmentHistory retrieves the punishment history for validators.
func ConsensusFetchPunishmentHistory(validatorID string, ledgerInstance *ledger.Ledger) ([]ledger.PunishmentRecord, error) {
    // Validate inputs
    if validatorID == "" {
        return nil, fmt.Errorf("validatorID cannot be empty")
    }
    if ledgerInstance == nil {
        return nil, fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Fetching punishment history for validator %s.", validatorID)

    // Fetch punishment history from the ledger
    punishmentHistory, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPunishmentHistory(validatorID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch punishment history for validator %s: %v", validatorID, err)
        return nil, fmt.Errorf("failed to fetch punishment history for validator %s: %w", validatorID, err)
    }

    log.Printf("[SUCCESS] Retrieved punishment history for validator %s with %d records.", validatorID, len(punishmentHistory))
    return punishmentHistory, nil
}


// ConsensusResetPunishmentCount resets the punishment count for a validator.
func ConsensusResetPunishmentCount(validatorID string, ledgerInstance *ledger.Ledger) error {
    // Validate inputs
    if validatorID == "" {
        return fmt.Errorf("validatorID cannot be empty")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Resetting punishment count for validator %s.", validatorID)

    // Reset the punishment count in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.ResetPunishmentCount(validatorID); err != nil {
        log.Printf("[ERROR] Failed to reset punishment count for validator %s: %v", validatorID, err)
        return fmt.Errorf("failed to reset punishment count for validator %s: %w", validatorID, err)
    }

    log.Printf("[SUCCESS] Punishment count reset for validator %s.", validatorID)
    return nil
}


// ConsensusSetAutoPunishmentRate sets the rate for automatic punishment escalation.
func ConsensusSetAutoPunishmentRate(rate float64, ledgerInstance *ledger.Ledger) error {
    // Input validation
    if rate < 0 {
        return fmt.Errorf("auto punishment rate must be non-negative")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Setting auto punishment rate to %.2f.", rate)

    // Set the auto punishment rate in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetAutoPunishmentRate(rate); err != nil {
        log.Printf("[ERROR] Failed to set auto punishment rate: %v", err)
        return fmt.Errorf("failed to set auto punishment rate: %w", err)
    }

    log.Printf("[SUCCESS] Auto punishment rate set to %.2f.", rate)
    return nil
}


// ConsensusGetAutoPunishmentRate retrieves the current automatic punishment rate.
func ConsensusGetAutoPunishmentRate(ledgerInstance *ledger.Ledger) (float64, error) {
    // Validate ledger instance
    if ledgerInstance == nil {
        return 0, fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Retrieving auto punishment rate.")

    // Retrieve the auto punishment rate from the ledger
    rate, err := ledgerInstance.BlockchainConsensusCoinLedger.GetAutoPunishmentRate()
    if err != nil {
        log.Printf("[ERROR] Failed to retrieve auto punishment rate: %v", err)
        return 0, fmt.Errorf("failed to get auto punishment rate: %w", err)
    }

    log.Printf("[SUCCESS] Retrieved auto punishment rate: %.2f.", rate)
    return rate, nil
}


// ConsensusEnableValidatorReinforcement enables reinforcement mechanisms for validators.
func ConsensusEnableValidatorReinforcement(ledgerInstance *ledger.Ledger) error {
    // Validate ledger instance
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Enabling validator reinforcement mechanisms.")

    // Enable validator reinforcement in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.EnableValidatorReinforcement(); err != nil {
        log.Printf("[ERROR] Failed to enable validator reinforcement: %v", err)
        return fmt.Errorf("failed to enable validator reinforcement: %w", err)
    }

    log.Printf("[SUCCESS] Validator reinforcement mechanisms enabled.")
    return nil
}


// ConsensusDisableValidatorReinforcement disables reinforcement mechanisms for validators.
func ConsensusDisableValidatorReinforcement(ledgerInstance *ledger.Ledger) error {
    // Validate ledger instance
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Disabling validator reinforcement mechanisms.")

    // Disable validator reinforcement in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.DisableValidatorReinforcement(); err != nil {
        log.Printf("[ERROR] Failed to disable validator reinforcement: %v", err)
        return fmt.Errorf("failed to disable validator reinforcement: %w", err)
    }

    log.Printf("[SUCCESS] Validator reinforcement mechanisms disabled.")
    return nil
}


// ConsensusTrackValidatorRewards tracks rewards earned by validators.
func ConsensusTrackValidatorRewards(validatorID string, reward float64, ledgerInstance *ledger.Ledger) error {
    // Input validation
    if validatorID == "" {
        return fmt.Errorf("validator ID cannot be empty")
    }
    if reward < 0 {
        return fmt.Errorf("reward cannot be negative")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Tracking reward for validator %s with reward amount %.2f.", validatorID, reward)

    // Initialize encryption instance
    encryptionInstance, err := common.NewEncryption(256)
    if err != nil {
        log.Printf("[ERROR] Failed to create encryption instance: %v", err)
        return fmt.Errorf("failed to create encryption instance: %w", err)
    }

    // Encrypt the reward amount
    encryptedReward, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%f", reward)), common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt reward for validator %s: %v", validatorID, err)
        return fmt.Errorf("failed to encrypt reward: %w", err)
    }

    // Record the encrypted reward in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.RecordValidatorReward(validatorID, encryptedReward); err != nil {
        log.Printf("[ERROR] Failed to record reward for validator %s: %v", validatorID, err)
        return fmt.Errorf("failed to track rewards for validator %s: %w", validatorID, err)
    }

    log.Printf("[SUCCESS] Reward of %.2f tracked for validator %s.", reward, validatorID)
    return nil
}


// ConsensusAuditRewardDistributions audits the distribution of rewards among validators.
func ConsensusAuditRewardDistributions(ledgerInstance *ledger.Ledger) error {
    // Input validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Auditing reward distributions.")

    // Perform the audit operation
    if err := ledgerInstance.BlockchainConsensusCoinLedger.AuditRewardDistributions(); err != nil {
        log.Printf("[ERROR] Failed to audit reward distributions: %v", err)
        return fmt.Errorf("failed to audit reward distributions: %w", err)
    }

    log.Printf("[SUCCESS] Reward distributions successfully audited.")
    return nil
}


// ConsensusFetchRewardHistory retrieves the reward history for a specific validator.
func ConsensusFetchRewardHistory(validatorID string, ledgerInstance *ledger.Ledger) ([]ledger.RewardRecord, error) {
    // Input validation
    if validatorID == "" {
        return nil, fmt.Errorf("validator ID cannot be empty")
    }
    if ledgerInstance == nil {
        return nil, fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Fetching reward history for validator %s.", validatorID)

    // Fetch the reward history from the ledger
    rewardHistory, err := ledgerInstance.BlockchainConsensusCoinLedger.GetRewardHistory(validatorID)
    if err != nil {
        log.Printf("[ERROR] Failed to fetch reward history for validator %s: %v", validatorID, err)
        return nil, fmt.Errorf("failed to fetch reward history for validator %s: %w", validatorID, err)
    }

    log.Printf("[SUCCESS] Reward history fetched for validator %s. Records count: %d", validatorID, len(rewardHistory))
    return rewardHistory, nil
}


// ConsensusSetPoHValidationWindow sets the time window for Proof of History (PoH) validation.
func ConsensusSetPoHValidationWindow(window time.Duration, ledgerInstance *ledger.Ledger) error {
    // Input validation
    if window <= 0 {
        return fmt.Errorf("validation window must be greater than zero")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Setting PoH validation window to %s.", window)

    // Set the PoH validation window in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetPoHValidationWindow(window); err != nil {
        log.Printf("[ERROR] Failed to set PoH validation window: %v", err)
        return fmt.Errorf("failed to set PoH validation window: %w", err)
    }

    log.Printf("[SUCCESS] PoH validation window set to %s.", window)
    return nil
}


// ConsensusGetPoHValidationWindow retrieves the PoH validation window.
func ConsensusGetPoHValidationWindow(ledgerInstance *ledger.Ledger) (time.Duration, error) {
    // Validate the ledger instance
    if ledgerInstance == nil {
        return 0, fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Retrieving PoH validation window.")

    // Retrieve the PoH validation window from the ledger
    window, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPoHValidationWindow()
    if err != nil {
        log.Printf("[ERROR] Failed to retrieve PoH validation window: %v", err)
        return 0, fmt.Errorf("failed to get PoH validation window: %w", err)
    }

    log.Printf("[SUCCESS] PoH validation window retrieved: %s.", window)
    return window, nil
}


// ConsensusAuditPoHValidation audits the PoH validation process.
func ConsensusAuditPoHValidation(ledgerInstance *ledger.Ledger) error {
    // Validate the ledger instance
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Auditing PoH validation process.")

    // Perform the PoH validation audit
    if err := ledgerInstance.BlockchainConsensusCoinLedger.AuditPoHValidation(); err != nil {
        log.Printf("[ERROR] Failed to audit PoH validation: %v", err)
        return fmt.Errorf("failed to audit PoH validation: %w", err)
    }

    log.Printf("[SUCCESS] PoH validation audit completed.")
    return nil
}


// ConsensusFetchPoHValidationLogs retrieves logs of PoH validation.
func ConsensusFetchPoHValidationLogs(ledgerInstance *ledger.Ledger) ([]ledger.PoHLog, error) {
    // Validate the ledger instance
    if ledgerInstance == nil {
        err := fmt.Errorf("ledger instance cannot be nil")
        log.Printf("[ERROR] Failed to fetch PoH validation logs: %v", err)
        return nil, err
    }

    log.Printf("[INFO] Fetching PoH validation logs.")

    // Fetch PoH validation logs from the ledger
    validationLogs, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPoHValidationLogs()
    if err != nil {
        log.Printf("[ERROR] Failed to fetch PoH validation logs: %v", err)
        return nil, fmt.Errorf("failed to fetch PoH validation logs: %w", err)
    }

    log.Printf("[SUCCESS] Successfully fetched %d PoH validation logs.", len(validationLogs))
    return validationLogs, nil
}

// ConsensusSetPoHFailureThreshold sets the failure threshold for PoH validations.
func ConsensusSetPoHFailureThreshold(threshold int, ledgerInstance *ledger.Ledger) error {
    // Validate inputs
    if threshold <= 0 {
        err := fmt.Errorf("failure threshold must be greater than zero")
        log.Printf("[ERROR] Invalid PoH failure threshold: %v", err)
        return err
    }
    if ledgerInstance == nil {
        err := fmt.Errorf("ledger instance cannot be nil")
        log.Printf("[ERROR] Failed to set PoH failure threshold: %v", err)
        return err
    }

    log.Printf("[INFO] Setting PoH failure threshold to %d.", threshold)

    // Set the PoH failure threshold in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetPoHFailureThreshold(threshold); err != nil {
        log.Printf("[ERROR] Failed to set PoH failure threshold: %v", err)
        return fmt.Errorf("failed to set PoH failure threshold: %w", err)
    }

    log.Printf("[SUCCESS] PoH failure threshold set to %d.", threshold)
    return nil
}


// ConsensusGetPoHFailureThreshold retrieves the current PoH failure threshold.
func ConsensusGetPoHFailureThreshold(ledgerInstance *ledger.Ledger) (int, error) {
    // Validate the ledger instance
    if ledgerInstance == nil {
        err := fmt.Errorf("ledger instance cannot be nil")
        log.Printf("[ERROR] Failed to get PoH failure threshold: %v", err)
        return 0, err
    }

    log.Printf("[INFO] Fetching PoH failure threshold.")

    // Fetch the PoH failure threshold from the ledger
    threshold, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPoHFailureThreshold()
    if err != nil {
        log.Printf("[ERROR] Failed to fetch PoH failure threshold: %v", err)
        return 0, fmt.Errorf("failed to get PoH failure threshold: %w", err)
    }

    log.Printf("[SUCCESS] Retrieved PoH failure threshold: %d.", threshold)
    return threshold, nil
}


// ConsensusEnablePoWHalving enables the PoW reward halving mechanism.
func ConsensusEnablePoWHalving(ledgerInstance *ledger.Ledger) error {
    // Validate the ledger instance
    if ledgerInstance == nil {
        err := fmt.Errorf("ledger instance cannot be nil")
        log.Printf("[ERROR] Failed to enable PoW halving: %v", err)
        return err
    }

    log.Printf("[INFO] Enabling PoW halving mechanism.")

    // Enable the PoW halving mechanism in the ledger
    ledgerInstance.BlockchainConsensusCoinLedger.PoWHalvingEnabled = true

    // Verify that the mechanism has been successfully enabled
    if !ledgerInstance.BlockchainConsensusCoinLedger.PoWHalvingEnabled {
        err := fmt.Errorf("failed to enable PoW halving mechanism due to internal error")
        log.Printf("[ERROR] %v", err)
        return err
    }

    log.Printf("[SUCCESS] PoW halving mechanism enabled.")
    return nil
}

// ConsensusDisablePoWHalving disables the PoW reward halving mechanism.
func ConsensusDisablePoWHalving(ledgerInstance *ledger.Ledger) error {
    // Validate the ledger instance
    if ledgerInstance == nil {
        err := fmt.Errorf("ledger instance cannot be nil")
        log.Printf("[ERROR] Failed to disable PoW halving: %v", err)
        return err
    }

    log.Printf("[INFO] Disabling PoW halving mechanism.")

    // Disable the PoW halving mechanism in the ledger
    ledgerInstance.BlockchainConsensusCoinLedger.PoWHalvingEnabled = false

    // Verify the operation
    if ledgerInstance.BlockchainConsensusCoinLedger.PoWHalvingEnabled {
        err := fmt.Errorf("failed to disable PoW halving due to internal error")
        log.Printf("[ERROR] %v", err)
        return err
    }

    log.Printf("[SUCCESS] PoW halving mechanism disabled.")
    return nil
}


// ConsensusSetPoWHalvingInterval sets the interval for PoW halving events.
func ConsensusSetPoWHalvingInterval(interval time.Duration, ledgerInstance *ledger.Ledger) error {
    // Validate inputs
    if ledgerInstance == nil {
        err := fmt.Errorf("ledger instance cannot be nil")
        log.Printf("[ERROR] Failed to set PoW halving interval: %v", err)
        return err
    }
    if interval <= 0 {
        err := fmt.Errorf("interval must be greater than zero")
        log.Printf("[ERROR] Invalid PoW halving interval: %v", err)
        return err
    }

    log.Printf("[INFO] Setting PoW halving interval to %s.", interval)

    // Set the PoW halving interval in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetPoWHalvingInterval(interval); err != nil {
        log.Printf("[ERROR] Failed to set PoW halving interval: %v", err)
        return fmt.Errorf("failed to set PoW halving interval: %w", err)
    }

    log.Printf("[SUCCESS] PoW halving interval set to %s.", interval)
    return nil
}


// ConsensusGetPoWHalvingInterval retrieves the interval for PoW halving events.
func ConsensusGetPoWHalvingInterval(ledgerInstance *ledger.Ledger) (time.Duration, error) {
    // Step 1: Input Validation
    if ledgerInstance == nil {
        err := fmt.Errorf("ledger instance cannot be nil")
        log.Printf("[ERROR] Failed to get PoW halving interval: %v", err)
        return 0, err
    }

    log.Printf("[INFO] Fetching PoW halving interval from the ledger.")

    // Step 2: Fetch PoW Halving Interval
    interval, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPoWHalvingInterval()
    if err != nil {
        log.Printf("[ERROR] Failed to fetch PoW halving interval: %v", err)
        return 0, fmt.Errorf("failed to get PoW halving interval: %w", err)
    }

    // Step 3: Verification
    if interval <= 0 {
        err := fmt.Errorf("invalid PoW halving interval retrieved: %d", interval)
        log.Printf("[ERROR] %v", err)
        return 0, err
    }

    log.Printf("[SUCCESS] Retrieved PoW halving interval: %s.", interval)
    return interval, nil
}

