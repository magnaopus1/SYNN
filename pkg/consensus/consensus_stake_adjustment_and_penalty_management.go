package consensus

import (
	"fmt"
	"log"
	"strconv"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// StakeChange represents a change in a validator's stake.
type StakeChange struct {
	ValidatorID  string
	ChangeAmount float64 // Encrypted stake change amount
	Timestamp    time.Time
}

// StakeLog represents a log entry for auditing stake adjustments.
type StakeLog struct {
	LogID       string
	ValidatorID string
	Adjustment  float64 // Encrypted adjustment amount
	Timestamp   time.Time
}

// ValidatorPenalty represents a penalty imposed on a validator.
type ValidatorPenalty struct {
	ValidatorID   string
	PenaltyAmount float64 // Encrypted penalty amount
	Timestamp     time.Time
}

// EpochLog represents a historical record of epoch changes for auditing.
type EpochLog struct {
	EpochID   string
	Duration  time.Duration // Time taken for the epoch
	Timestamp time.Time
}

// ReinforcementPolicy defines the rules for reinforcing consensus.
type ReinforcementPolicy struct {
	PolicyID    string
	Description string
	Details     map[string]interface{} // Specifics of the policy
	Timestamp   time.Time
}

// HealthLog represents health metrics of the consensus for analysis.
type HealthLog struct {
	HealthID  string
	Metric    string
	Value     float64
	Timestamp time.Time
}

// consensusTrackValidatorStakeChanges tracks changes in validator stakes with encryption.
func ConsensusTrackValidatorStakeChanges(validatorID string, stakeChange float64, ledgerInstance *ledger.Ledger) error {
    // Input Validation
    if validatorID == "" {
        return fmt.Errorf("validator ID cannot be empty")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Tracking stake change for validator: %s", validatorID)

    // Create Encryption Instance
    encryptionInstance, err := common.NewEncryption(256)
    if err != nil {
        log.Printf("[ERROR] Failed to create encryption instance: %v", err)
        return fmt.Errorf("failed to create encryption instance: %w", err)
    }

    // Encrypt the stake change
    stakeChangeStr := fmt.Sprintf("%f", stakeChange)
    encryptedStakeChange, err := encryptionInstance.EncryptData("AES", []byte(stakeChangeStr), common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt stake change for validator %s: %v", validatorID, err)
        return fmt.Errorf("failed to encrypt stake change: %w", err)
    }

    // Create Stake Change Record
    stakeChangeRecord := ledger.StakeChangeRecord{
        ValidatorID:          validatorID,
        EncryptedStakeChange: encryptedStakeChange,
        Timestamp:            time.Now(),
    }

    // Record the Stake Change in the Ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.RecordStakeChange(stakeChangeRecord); err != nil {
        log.Printf("[ERROR] Failed to record stake change for validator %s: %v", validatorID, err)
        return fmt.Errorf("failed to record stake change: %w", err)
    }

    log.Printf("[SUCCESS] Successfully tracked encrypted stake change for validator: %s", validatorID)
    return nil
}


// consensusLogStakeAdjustments logs stake adjustments for auditing purposes.
func ConsensusLogStakeAdjustments(stakeLog ledger.StakeLog, ledgerInstance *ledger.Ledger) error {
    // Input Validation
    if stakeLog.ValidatorID == "" {
        return fmt.Errorf("validator ID cannot be empty in stake log")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Logging stake adjustment for validator: %s", stakeLog.ValidatorID)

    // Create Encryption Instance
    encryptionInstance, err := common.NewEncryption(256)
    if err != nil {
        log.Printf("[ERROR] Failed to create encryption instance: %v", err)
        return fmt.Errorf("failed to create encryption instance: %w", err)
    }

    // Encrypt the stake adjustment
    adjustmentStr := fmt.Sprintf("%f", stakeLog.Adjustment)
    encryptedAdjustment, err := encryptionInstance.EncryptData("AES", []byte(adjustmentStr), common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt stake adjustment for validator %s: %v", stakeLog.ValidatorID, err)
        return fmt.Errorf("failed to encrypt stake adjustment: %w", err)
    }
    stakeLog.EncryptedAdjustment = encryptedAdjustment

    // Log the Stake Adjustment in the Ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.LogStakeAdjustment(stakeLog); err != nil {
        log.Printf("[ERROR] Failed to log stake adjustment for validator %s: %v", stakeLog.ValidatorID, err)
        return fmt.Errorf("failed to log stake adjustments: %w", err)
    }

    log.Printf("[SUCCESS] Successfully logged stake adjustment for validator: %s", stakeLog.ValidatorID)
    return nil
}


// consensusSetValidatorPenalty sets a penalty for a specific validator.
func ConsensusSetValidatorPenalty(validatorID string, penaltyAmount float64, ledgerInstance *ledger.Ledger) error {
    // Input Validation
    if validatorID == "" {
        return fmt.Errorf("validator ID cannot be empty")
    }
    if penaltyAmount <= 0 {
        return fmt.Errorf("penalty amount must be greater than zero")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Setting penalty for validator %s. Amount: %.2f", validatorID, penaltyAmount)

    // Create Encryption Instance
    encryptionInstance, err := common.NewEncryption(256)
    if err != nil {
        log.Printf("[ERROR] Failed to create encryption instance: %v", err)
        return fmt.Errorf("failed to create encryption instance: %w", err)
    }

    // Encrypt the Penalty Amount
    penaltyStr := fmt.Sprintf("%f", penaltyAmount)
    encryptedPenalty, err := encryptionInstance.EncryptData("AES", []byte(penaltyStr), common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt penalty amount for validator %s: %v", validatorID, err)
        return fmt.Errorf("failed to encrypt penalty amount: %w", err)
    }

    // Set the Penalty in the Ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetValidatorPenalty(validatorID, encryptedPenalty); err != nil {
        log.Printf("[ERROR] Failed to set penalty for validator %s: %v", validatorID, err)
        return fmt.Errorf("failed to set penalty for validator %s: %w", validatorID, err)
    }

    log.Printf("[SUCCESS] Penalty of %.2f set for validator %s.", penaltyAmount, validatorID)
    return nil
}


// consensusGetValidatorPenalty retrieves the penalty amount for a specific validator.
func ConsensusGetValidatorPenalty(validatorID string, ledgerInstance *ledger.Ledger) (float64, error) {
    // Input Validation
    if validatorID == "" {
        return 0, fmt.Errorf("validator ID cannot be empty")
    }
    if ledgerInstance == nil {
        return 0, fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Retrieving penalty for validator %s", validatorID)

    // Retrieve the Encrypted Penalty from the Ledger
    encryptedPenalty, err := ledgerInstance.BlockchainConsensusCoinLedger.GetValidatorPenalty(validatorID)
    if err != nil {
        log.Printf("[ERROR] Failed to retrieve penalty for validator %s: %v", validatorID, err)
        return 0, fmt.Errorf("failed to retrieve penalty for validator %s: %w", validatorID, err)
    }

    // Create Encryption Instance
    encryptionInstance, err := common.NewEncryption(256)
    if err != nil {
        log.Printf("[ERROR] Failed to create encryption instance: %v", err)
        return 0, fmt.Errorf("failed to create encryption instance: %w", err)
    }

    // Decrypt the Penalty Amount
    decryptedPenalty, err := encryptionInstance.DecryptData(encryptedPenalty, common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to decrypt penalty for validator %s: %v", validatorID, err)
        return 0, fmt.Errorf("failed to decrypt penalty: %w", err)
    }

    // Parse the Decrypted Penalty into a Float
    penalty, err := strconv.ParseFloat(string(decryptedPenalty), 64)
    if err != nil {
        log.Printf("[ERROR] Failed to parse decrypted penalty for validator %s: %v", validatorID, err)
        return 0, fmt.Errorf("failed to parse decrypted penalty: %w", err)
    }

    log.Printf("[SUCCESS] Retrieved penalty for validator %s. Amount: %.2f", validatorID, penalty)
    return penalty, nil
}


// consensusSetEpochTimeout sets the timeout period for an epoch.
func ConsensusSetEpochTimeout(timeout time.Duration, ledgerInstance *ledger.Ledger) error {
    // Input Validation
    if timeout <= 0 {
        return fmt.Errorf("timeout must be greater than zero")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Setting epoch timeout to %s", timeout)

    // Set the timeout in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetEpochTimeout(timeout); err != nil {
        log.Printf("[ERROR] Failed to set epoch timeout: %v", err)
        return fmt.Errorf("failed to set epoch timeout: %w", err)
    }

    log.Printf("[SUCCESS] Epoch timeout set to %s", timeout)
    return nil
}


// consensusGetEpochTimeout retrieves the timeout period for the current epoch.
func ConsensusGetEpochTimeout(ledgerInstance *ledger.Ledger) (time.Duration, error) {
    // Input Validation
    if ledgerInstance == nil {
        return 0, fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Retrieving epoch timeout")

    // Get the timeout from the ledger
    timeout, err := ledgerInstance.BlockchainConsensusCoinLedger.GetEpochTimeout()
    if err != nil {
        log.Printf("[ERROR] Failed to retrieve epoch timeout: %v", err)
        return 0, fmt.Errorf("failed to get epoch timeout: %w", err)
    }

    log.Printf("[SUCCESS] Retrieved epoch timeout: %s", timeout)
    return timeout, nil
}


// ConsensusTrackEpochTime tracks the time taken for each epoch.
func ConsensusTrackEpochTime(epochID string, duration time.Duration, ledgerInstance *ledger.Ledger) error {
    // Input Validation
    if epochID == "" {
        return fmt.Errorf("epochID cannot be empty")
    }
    if duration <= 0 {
        return fmt.Errorf("duration must be greater than zero")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Tracking time for epoch %s with duration %s", epochID, duration)

    // Record the epoch time in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.RecordEpochTime(epochID, duration); err != nil {
        log.Printf("[ERROR] Failed to track epoch time for epoch %s: %v", epochID, err)
        return fmt.Errorf("failed to track epoch time for epoch %s: %w", epochID, err)
    }

    log.Printf("[SUCCESS] Tracked time for epoch %s: %s", epochID, duration)
    return nil
}



// consensusLogEpochChanges logs changes in epochs for historical auditing.
func ConsensusLogEpochChanges(epochLog ledger.EpochLog, ledgerInstance *ledger.Ledger) error {
    // Input Validation
    if epochLog.EpochID == "" {
        return fmt.Errorf("epochID in epochLog cannot be empty")
    }
    if epochLog.Duration <= 0 {
        return fmt.Errorf("epoch duration must be greater than zero")
    }
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Logging changes for epoch %s", epochLog.EpochID)

    // Initialize the encryption instance
    encryptionInstance, err := common.NewEncryption(256)
    if err != nil {
        log.Printf("[ERROR] Failed to create encryption instance: %v", err)
        return fmt.Errorf("failed to create encryption instance: %w", err)
    }

    // Encrypt the epoch duration (as a string representation of float64 seconds)
    encryptedDuration, err := encryptionInstance.EncryptData(
        "AES",
        []byte(fmt.Sprintf("%f", float64(epochLog.Duration.Seconds()))),
        common.EncryptionKey,
    )
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt duration for epoch %s: %v", epochLog.EpochID, err)
        return fmt.Errorf("failed to encrypt duration: %w", err)
    }

    // Assign the encrypted duration to the epoch log
    epochLog.EncryptedDuration = encryptedDuration

    // Log the epoch changes in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.LogEpochChange(epochLog); err != nil {
        log.Printf("[ERROR] Failed to log epoch changes for epoch %s: %v", epochLog.EpochID, err)
        return fmt.Errorf("failed to log epoch changes for epoch %s: %w", epochLog.EpochID, err)
    }

    log.Printf("[SUCCESS] Logged changes for epoch %s", epochLog.EpochID)
    return nil
}


// consensusSetConsensusReinforcementPolicy sets the policy for consensus reinforcement.
func ConsensusSetConsensusReinforcementPolicy(policy ledger.ReinforcementPolicy, ledgerInstance *ledger.Ledger) error {
    // Validate input
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }
    if len(policy.Details) == 0 {
        return fmt.Errorf("policy details cannot be empty")
    }

    log.Printf("[INFO] Setting consensus reinforcement policy")

    // Initialize encryption instance
    encryptionInstance, err := common.NewEncryption(256)
    if err != nil {
        log.Printf("[ERROR] Failed to create encryption instance: %v", err)
        return fmt.Errorf("failed to create encryption instance: %w", err)
    }

    // Encrypt each detail in the policy
    for key, value := range policy.Details {
        floatValue, ok := value.(float64)
        if !ok {
            return fmt.Errorf("policy detail value must be a float64 for key %s", key)
        }

        encryptedValue, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%f", floatValue)), common.EncryptionKey)
        if err != nil {
            log.Printf("[ERROR] Failed to encrypt policy detail for key %s: %v", key, err)
            return fmt.Errorf("failed to encrypt policy detail for key %s: %w", key, err)
        }
        policy.Details[key] = encryptedValue
    }

    // Store the encrypted policy in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetReinforcementPolicy(policy); err != nil {
        log.Printf("[ERROR] Failed to set reinforcement policy: %v", err)
        return fmt.Errorf("failed to set reinforcement policy: %w", err)
    }

    log.Printf("[SUCCESS] Consensus reinforcement policy set successfully")
    return nil
}


// consensusGetConsensusReinforcementPolicy retrieves the current reinforcement policy.
func ConsensusGetConsensusReinforcementPolicy(ledgerInstance *ledger.Ledger) (ledger.ReinforcementPolicy, error) {
    // Validate input
    if ledgerInstance == nil {
        return ledger.ReinforcementPolicy{}, fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Retrieving consensus reinforcement policy")

    // Retrieve the encrypted policy from the ledger
    encryptedPolicy, err := ledgerInstance.BlockchainConsensusCoinLedger.GetReinforcementPolicy()
    if err != nil {
        log.Printf("[ERROR] Failed to get reinforcement policy: %v", err)
        return ledger.ReinforcementPolicy{}, fmt.Errorf("failed to get reinforcement policy: %w", err)
    }

    // Initialize encryption instance
    encryptionInstance, err := common.NewEncryption(256)
    if err != nil {
        log.Printf("[ERROR] Failed to create encryption instance: %v", err)
        return ledger.ReinforcementPolicy{}, fmt.Errorf("failed to create encryption instance: %w", err)
    }

    // Decrypt each policy detail value
    for key, value := range encryptedPolicy.Details {
        encryptedValue, ok := value.([]byte)
        if !ok {
            return ledger.ReinforcementPolicy{}, fmt.Errorf("policy detail must be of type []byte for key %s", key)
        }

        decryptedValue, err := encryptionInstance.DecryptData(encryptedValue, common.EncryptionKey)
        if err != nil {
            log.Printf("[ERROR] Failed to decrypt policy detail for key %s: %v", key, err)
            return ledger.ReinforcementPolicy{}, fmt.Errorf("failed to decrypt policy detail for key %s: %w", key, err)
        }

        // Convert decrypted data back to float64
        floatValue, err := strconv.ParseFloat(string(decryptedValue), 64)
        if err != nil {
            log.Printf("[ERROR] Failed to parse decrypted value for key %s: %v", key, err)
            return ledger.ReinforcementPolicy{}, fmt.Errorf("failed to parse decrypted policy detail for key %s: %w", key, err)
        }
        encryptedPolicy.Details[key] = floatValue
    }

    log.Printf("[SUCCESS] Consensus reinforcement policy retrieved successfully")
    return encryptedPolicy, nil
}


// consensusLogConsensusHealth logs the health metrics of the consensus for future analysis.
func consensusLogConsensusHealth(healthLog ledger.HealthLog, ledgerInstance *ledger.Ledger) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }
    if healthLog.Metric == "" {
        return fmt.Errorf("health log metric name cannot be empty")
    }
    if healthLog.Timestamp.IsZero() {
        return fmt.Errorf("health log timestamp cannot be zero")
    }

    log.Printf("[INFO] Logging consensus health metrics for metric: %s", healthLog.Metric)

    // Initialize encryption instance
    encryptionInstance, err := common.NewEncryption(256)
    if err != nil {
        log.Printf("[ERROR] Failed to create encryption instance: %v", err)
        return fmt.Errorf("failed to create encryption instance: %w", err)
    }

    // Encrypt the health metric value
    encryptedValue, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%f", healthLog.Value)), common.EncryptionKey)
    if err != nil {
        log.Printf("[ERROR] Failed to encrypt health metric value for metric: %s, error: %v", healthLog.Metric, err)
        return fmt.Errorf("failed to encrypt health metric value: %w", err)
    }

    // Assign the encrypted value to the health log
    healthLog.EncryptedValue = encryptedValue

    // Log the encrypted health metrics in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.LogHealthMetrics(healthLog); err != nil {
        log.Printf("[ERROR] Failed to log health metrics for metric: %s, error: %v", healthLog.Metric, err)
        return fmt.Errorf("failed to log consensus health metrics: %w", err)
    }

    log.Printf("[SUCCESS] Consensus health metrics logged successfully for metric: %s", healthLog.Metric)
    return nil
}

