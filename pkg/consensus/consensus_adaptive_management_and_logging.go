package consensus

import (
	"fmt"
	"log"
	"synnergy_network/pkg/ledger"
	"time"
)

// PunitiveMeasureRecord represents a record of punitive actions taken against validators.
type PunitiveMeasureRecord struct {
	ActionID    string    // Unique identifier for the action
	ValidatorID string    // Validator subject to punitive action
	Reason      string    // Reason for the punitive action
	Timestamp   time.Time // Time the punitive action was taken
	Status      string    // Status of the action (e.g., active, reverted)
}

// PunishmentAdjustmentLog represents a log entry for adjustments made to punitive measures.
type PunishmentAdjustmentLog struct {
	AdjustmentID string    // Unique identifier for the adjustment
	ActionID     string    // Related punitive action ID
	AdjustedBy   string    // Who made the adjustment
	Timestamp    time.Time // Time the adjustment was made
	Details      string    // Details of the adjustment
}

// consensusFetchPunitiveMeasureLogs retrieves logs of punitive measures applied to validators.
func ConsensusFetchPunitiveMeasureLogs(ledgerInstance *ledger.Ledger) ([]ledger.PunitiveMeasureRecord, error) {
    if ledgerInstance == nil {
        return nil, fmt.Errorf("ledger instance is nil")
    }

    log.Printf("[Info] Fetching punitive measure logs from ledger...")

    // Locking the ledger for safe access
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Fetch logs and validate response
    logs, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPunitiveMeasureLogs()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch punitive measure logs: %w", err)
    }

    if len(logs) == 0 {
        log.Printf("[Warning] No punitive measure logs found in the ledger.")
        return logs, nil
    }

    log.Printf("[Success] Fetched %d punitive measure logs.", len(logs))
    return logs, nil
}


// consensusRevertPunitiveActions reverts punitive actions based on provided criteria, updating the ledger.
func ConsensusRevertPunitiveActions(actionID string, ledgerInstance *ledger.Ledger) error {
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }

    if actionID == "" {
        return fmt.Errorf("actionID is empty")
    }

    log.Printf("[Info] Starting punitive action revert for actionID: %s", actionID)

    // Lock the ledger for secure updates
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Validate if the punitive action exists
    existingAction, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPunitiveActionByID(actionID)
    if err != nil {
        return fmt.Errorf("failed to validate existence of punitive action with ID %s: %w", actionID, err)
    }

    if existingAction == nil {
        return fmt.Errorf("no punitive action found with ID: %s", actionID)
    }

    log.Printf("[Info] Valid punitive action found. Proceeding with revert for actionID: %s", actionID)

    // Revert the punitive action
    if err := ledgerInstance.BlockchainConsensusCoinLedger.RevertPunitiveAction(actionID); err != nil {
        return fmt.Errorf("failed to revert punitive action with ID %s: %w", actionID, err)
    }

    // Validate post-revert state
    updatedAction, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPunitiveActionByID(actionID)
    if err != nil {
        return fmt.Errorf("failed to validate post-revert state for actionID %s: %w", actionID, err)
    }

    if updatedAction != nil && updatedAction.Status != "reverted" {
        return fmt.Errorf("reversion failed; actionID %s is not marked as reverted", actionID)
    }

    log.Printf("[Success] Punitive action with ID %s successfully reverted.", actionID)
    return nil
}


// consensusSetPunishmentReevaluationInterval sets the interval at which punishments are re-evaluated.
func ConsensusSetPunishmentReevaluationInterval(interval time.Duration, ledgerInstance *ledger.Ledger) error {
    // Input validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }
    if interval <= 0 {
        return fmt.Errorf("invalid interval: %v. Interval must be greater than zero", interval)
    }

    log.Printf("[Info] Setting punishment reevaluation interval to: %v", interval)

    // Lock the ledger to prevent race conditions
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Set the punishment reevaluation interval
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetPunishmentReevaluationInterval(interval); err != nil {
        return fmt.Errorf("failed to set punishment reevaluation interval: %w", err)
    }

    // Verify that the interval was successfully updated
    updatedInterval, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPunishmentReevaluationInterval()
    if err != nil {
        return fmt.Errorf("failed to verify updated punishment reevaluation interval: %w", err)
    }

    if updatedInterval != interval {
        return fmt.Errorf("verification failed: expected interval %v, got %v", interval, updatedInterval)
    }

    log.Printf("[Success] Punishment reevaluation interval set to: %v", interval)
    return nil
}


// consensusGetPunishmentReevaluationInterval retrieves the current punishment reevaluation interval.
func ConsensusGetPunishmentReevaluationInterval(ledgerInstance *ledger.Ledger) (time.Duration, error) {
    // Input validation
    if ledgerInstance == nil {
        return 0, fmt.Errorf("ledger instance is nil")
    }

    log.Println("[Info] Retrieving punishment reevaluation interval...")

    // Lock the ledger to ensure safe access
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Fetch the punishment reevaluation interval
    interval, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPunishmentReevaluationInterval()
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve punishment reevaluation interval: %w", err)
    }

    // Validate that the interval is a positive duration
    if interval <= 0 {
        return 0, fmt.Errorf("invalid punishment reevaluation interval retrieved: %v", interval)
    }

    log.Printf("[Success] Current punishment reevaluation interval: %v", interval)
    return interval, nil
}


// consensusLogPunishmentAdjustments logs all adjustments made to punitive measures.
func ConsensusLogPunishmentAdjustments(actionID, adjustedBy, details string, ledgerInstance *ledger.Ledger) error {
    // Input Validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }
    if actionID == "" || adjustedBy == "" || details == "" {
        return fmt.Errorf("invalid input: actionID, adjustedBy, and details must not be empty")
    }

    log.Printf("[Info] Logging punishment adjustment: actionID=%s, adjustedBy=%s, details=%s", actionID, adjustedBy, details)

    // Lock the ledger to ensure thread safety
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Log the punishment adjustment in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.LogPunishmentAdjustments(actionID, adjustedBy, details); err != nil {
        return fmt.Errorf("failed to log punishment adjustments: %w", err)
    }

    // Verify that the adjustment has been logged
    adjustmentLogs, err := ledgerInstance.BlockchainConsensusCoinLedger.GetAdjustmentLogs(actionID)
    if err != nil {
        return fmt.Errorf("failed to retrieve adjustment logs for verification: %w", err)
    }

    log.Printf("[Success] Punishment adjustment logged successfully for actionID=%s. Total adjustments logged: %d", actionID, len(adjustmentLogs))
    return nil
}


// consensusSetAdaptiveDifficulty dynamically adjusts the difficulty based on network conditions.
func ConsensusSetAdaptiveDifficulty(difficultyLevel int, ledgerInstance *ledger.Ledger) error {
    // Input Validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }
    if difficultyLevel <= 0 {
        return fmt.Errorf("invalid difficulty level: %d. Difficulty must be greater than zero", difficultyLevel)
    }

    log.Printf("[Info] Setting adaptive difficulty to: %d", difficultyLevel)

    // Lock the ledger to ensure thread-safe updates
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Update the adaptive difficulty in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetAdaptiveDifficulty(difficultyLevel); err != nil {
        return fmt.Errorf("failed to set adaptive difficulty: %w", err)
    }

    // Verify that the difficulty was updated
    updatedDifficulty, err := ledgerInstance.BlockchainConsensusCoinLedger.GetCurrentDifficulty()
    if err != nil {
        return fmt.Errorf("failed to retrieve updated difficulty for verification: %w", err)
    }

    if updatedDifficulty != difficultyLevel {
        return fmt.Errorf("verification failed: expected difficulty %d, got %d", difficultyLevel, updatedDifficulty)
    }

    log.Printf("[Success] Adaptive difficulty set to: %d", difficultyLevel)
    return nil
}

// consensusGetAdaptiveDifficulty retrieves the current adaptive difficulty level.
func ConsensusGetAdaptiveDifficulty(ledgerInstance *ledger.Ledger) (int, error) {
    // Validate ledger instance
    if ledgerInstance == nil {
        return 0, fmt.Errorf("ledger instance is nil")
    }

    log.Printf("[Info] Retrieving adaptive difficulty level...")

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Retrieve the adaptive difficulty level
    difficulty, err := ledgerInstance.BlockchainConsensusCoinLedger.GetAdaptiveDifficulty()
    if err != nil {
        return 0, fmt.Errorf("failed to get adaptive difficulty: %w", err)
    }

    log.Printf("[Success] Retrieved adaptive difficulty level: %d", difficulty)
    return difficulty, nil
}


// consensusEnableAdaptiveRewardDistribution enables the dynamic distribution of rewards based on network needs.
func ConsensusEnableAdaptiveRewardDistribution(ledgerInstance *ledger.Ledger) error {
    // Validate ledger instance
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }

    log.Printf("[Info] Enabling adaptive reward distribution...")

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Enable adaptive reward distribution
    if err := ledgerInstance.BlockchainConsensusCoinLedger.EnableAdaptiveRewardDistribution(); err != nil {
        return fmt.Errorf("failed to enable adaptive reward distribution: %w", err)
    }

    // Verify that adaptive reward distribution is enabled
    isEnabled, err := ledgerInstance.BlockchainConsensusCoinLedger.IsAdaptiveRewardDistributionEnabled()
    if err != nil {
        return fmt.Errorf("failed to verify adaptive reward distribution state: %w", err)
    }

    if !isEnabled {
        return fmt.Errorf("verification failed: adaptive reward distribution is not enabled")
    }

    log.Printf("[Success] Adaptive reward distribution enabled successfully.")
    return nil
}


// consensusDisableAdaptiveRewardDistribution disables adaptive reward distribution.
func ConsensusDisableAdaptiveRewardDistribution(ledgerInstance *ledger.Ledger) error {
    // Validate ledger instance
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }

    log.Printf("[Info] Disabling adaptive reward distribution...")

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Disable adaptive reward distribution
    if err := ledgerInstance.BlockchainConsensusCoinLedger.DisableAdaptiveRewardDistribution(); err != nil {
        return fmt.Errorf("failed to disable adaptive reward distribution: %w", err)
    }

    // Verify that adaptive reward distribution is disabled
    isEnabled, err := ledgerInstance.BlockchainConsensusCoinLedger.IsAdaptiveRewardDistributionEnabled()
    if err != nil {
        return fmt.Errorf("failed to verify adaptive reward distribution state: %w", err)
    }

    if isEnabled {
        return fmt.Errorf("verification failed: adaptive reward distribution is still enabled")
    }

    log.Printf("[Success] Adaptive reward distribution disabled successfully.")
    return nil
}
