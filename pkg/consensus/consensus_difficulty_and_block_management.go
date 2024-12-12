package consensus

import (
	"fmt"
	"log"
	"synnergy_network/pkg/ledger"
	"time"
)

// DifficultyAdjustmentLog records adjustments made to consensus difficulty.
type DifficultyAdjustmentLog struct {
	AdjustmentID       string
	Timestamp          time.Time
	NewDifficultyLevel int
	Reason             string
}

// BlockGenerationLog records block generation times.
type BlockGenerationLog struct {
	BlockID        string
	GenerationTime time.Duration
	Timestamp      time.Time
}

// ConsensusAuditLog tracks consensus audit details.
type ConsensusAuditLog struct {
	AuditID             string
	ValidatorID         string
	Timestamp           time.Time
	ParticipationStatus string
}

// FinalityCheckLog records finality checks for blocks.
type FinalityCheckLog struct {
	BlockID        string
	FinalityStatus bool
	Timestamp      time.Time
}

// ValidatorActivityLog tracks validator activity.
type ValidatorActivityLog struct {
	ValidatorID string
	Action      string
	Timestamp   time.Time
	Details     string
}

// RewardDistributionMode and ValidatorSelectionMode settings
type RewardDistributionMode struct {
	ModeID      string
	Description string
	Active      bool
}

type ValidatorSelectionMode struct {
	ModeID      string
	Description string
	Active      bool
}

// consensusAdjustDifficultyBasedOnTime adjusts difficulty based on recent block generation times.
func ConsensusAdjustDifficultyBasedOnTime(newLevel int, reason string, ledgerInstance *ledger.Ledger) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }
    if newLevel < 0 {
        return fmt.Errorf("invalid difficulty level: %d", newLevel)
    }
    if reason == "" {
        return fmt.Errorf("reason for difficulty adjustment is required")
    }

    log.Printf("[Info] Adjusting difficulty level to %d. Reason: %s", newLevel, reason)

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Adjust difficulty level in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetDifficultyLevel(newLevel, reason); err != nil {
        return fmt.Errorf("failed to adjust difficulty: %w", err)
    }

    // Post-action verification
    currentDifficulty, err := ledgerInstance.BlockchainConsensusCoinLedger.GetCurrentDifficultyLevel()
    if err != nil {
        return fmt.Errorf("failed to verify difficulty adjustment: %w", err)
    }
    if currentDifficulty != newLevel {
        return fmt.Errorf("verification failed: difficulty not set to the intended level (%d)", newLevel)
    }

    log.Printf("[Success] Difficulty adjusted to %d successfully.", newLevel)
    return nil
}


// consensusMonitorBlockGenerationTime logs the block generation time.
func consensusMonitorBlockGenerationTime(blockID string, generationTime time.Duration, ledgerInstance *ledger.Ledger) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }
    if blockID == "" {
        return fmt.Errorf("block ID is required")
    }
    if generationTime <= 0 {
        return fmt.Errorf("invalid generation time: %v", generationTime)
    }

    log.Printf("[Info] Logging block generation time. BlockID: %s, Time: %v", blockID, generationTime)

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Log block generation time in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.LogBlockGeneration(blockID, generationTime); err != nil {
        return fmt.Errorf("failed to log block generation time: %w", err)
    }

    // Post-action verification
    logEntry, err := ledgerInstance.BlockchainConsensusCoinLedger.GetBlockGenerationLog(blockID)
    if err != nil {
        return fmt.Errorf("failed to verify block generation log: %w", err)
    }
    if logEntry.GenerationTime != generationTime {
        return fmt.Errorf("verification failed: logged generation time does not match the provided value")
    }

    log.Printf("[Success] Block generation time logged successfully. BlockID: %s, Time: %v", blockID, generationTime)
    return nil
}


// consensusEnableConsensusAudit enables auditing of consensus participation.
func consensusEnableConsensusAudit(ledgerInstance *ledger.Ledger) error {
    // Input validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }

    log.Printf("[Info] Enabling consensus audit...")

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Enable consensus audit
    err := ledgerInstance.BlockchainConsensusCoinLedger.EnableConsensusAudit()
    if err != nil {
        return fmt.Errorf("failed to enable consensus audit: %w", err)
    }

    // Post-action verification
    isAuditEnabled, err := ledgerInstance.BlockchainConsensusCoinLedger.IsConsensusAuditEnabled()
    if err != nil {
        return fmt.Errorf("failed to verify consensus audit status: %w", err)
    }
    if !isAuditEnabled {
        return fmt.Errorf("verification failed: consensus audit not enabled as expected")
    }

    log.Printf("[Success] Consensus audit enabled successfully.")
    return nil
}


// consensusDisableConsensusAudit disables auditing of consensus participation.
func consensusDisableConsensusAudit(ledgerInstance *ledger.Ledger) error {
    // Input validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }

    log.Printf("[Info] Disabling consensus audit...")

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Disable consensus audit
    err := ledgerInstance.BlockchainConsensusCoinLedger.DisableConsensusAudit()
    if err != nil {
        return fmt.Errorf("failed to disable consensus audit: %w", err)
    }

    // Post-action verification
    isAuditEnabled, err := ledgerInstance.BlockchainConsensusCoinLedger.IsConsensusAuditEnabled()
    if err != nil {
        return fmt.Errorf("failed to verify consensus audit status: %w", err)
    }
    if isAuditEnabled {
        return fmt.Errorf("verification failed: consensus audit still enabled after attempted disable")
    }

    log.Printf("[Success] Consensus audit disabled successfully.")
    return nil
}


// consensusSetRewardDistributionMode sets the reward distribution mode.
func consensusSetRewardDistributionMode(mode ledger.RewardDistributionMode, ledgerInstance *ledger.Ledger) error {
    // Input validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }
    if !isValidRewardDistributionMode(mode) {
        return fmt.Errorf("invalid reward distribution mode: %v", mode)
    }

    log.Printf("[Info] Setting reward distribution mode to: %v", mode)

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Set the reward distribution mode
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetRewardDistributionMode(mode); err != nil {
        return fmt.Errorf("failed to set reward distribution mode: %w", err)
    }

    // Post-action verification
    currentMode, err := ledgerInstance.BlockchainConsensusCoinLedger.GetRewardDistributionMode()
    if err != nil {
        return fmt.Errorf("failed to verify reward distribution mode: %w", err)
    }
    if currentMode != mode {
        return fmt.Errorf("verification failed: expected mode %v, but found %v", mode, currentMode)
    }

    log.Printf("[Success] Reward distribution mode successfully set to: %v", mode)
    return nil
}


// consensusGetRewardDistributionMode retrieves the current reward distribution mode.
func consensusGetRewardDistributionMode(ledgerInstance *ledger.Ledger) (ledger.RewardDistributionMode, error) {
    // Input validation
    if ledgerInstance == nil {
        return ledger.RewardDistributionMode{}, fmt.Errorf("ledger instance is nil")
    }

    log.Printf("[Info] Retrieving current reward distribution mode...")

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Get the reward distribution mode
    mode, err := ledgerInstance.BlockchainConsensusCoinLedger.GetRewardDistributionMode()
    if err != nil {
        return ledger.RewardDistributionMode{}, fmt.Errorf("failed to retrieve reward distribution mode: %w", err)
    }

    log.Printf("[Success] Current reward distribution mode: %v", mode)
    return mode, nil
}


// consensusTrackConsensusParticipation logs validator participation in consensus.
func ConsensusTrackConsensusParticipation(validatorID, status string, ledgerInstance *ledger.Ledger) error {
    // Input validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }
    if validatorID == "" {
        return fmt.Errorf("validator ID cannot be empty")
    }
    if status == "" {
        return fmt.Errorf("status cannot be empty")
    }
    if !isValidConsensusStatus(status) {
        return fmt.Errorf("invalid consensus status: %s", status)
    }

    log.Printf("[Info] Logging consensus participation: ValidatorID=%s, Status=%s", validatorID, status)

    // Lock ledger for thread-safe updates
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Log participation in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.LogConsensusParticipation(validatorID, status); err != nil {
        return fmt.Errorf("failed to log consensus participation: %w", err)
    }

    // Post-action verification
    participationLogs := ledgerInstance.BlockchainConsensusCoinLedger.GetParticipationLogs()
    found := false
    for _, logEntry := range participationLogs {
        if logEntry.ValidatorID == validatorID && logEntry.Status == status {
            found = true
            break
        }
    }
    if !found {
        return fmt.Errorf("post-verification failed: participation log not found for ValidatorID=%s, Status=%s", validatorID, status)
    }

    log.Printf("[Success] Consensus participation logged successfully: ValidatorID=%s, Status=%s", validatorID, status)
    return nil
}

// Validates the status for consensus participation
func isValidConsensusStatus(status string) bool {
    validStatuses := []string{"Active", "Inactive", "Disqualified", "Suspended"}
    for _, validStatus := range validStatuses {
        if status == validStatus {
            return true
        }
    }
    return false
}


// consensusFetchConsensusLogs retrieves consensus audit logs.
func ConsensusFetchConsensusLogs(ledgerInstance *ledger.Ledger) ([]ledger.ConsensusAuditLog, error) {
    // Input validation
    if ledgerInstance == nil {
        return nil, fmt.Errorf("ledger instance is nil")
    }

    log.Printf("[Info] Fetching consensus audit logs...")

    // Lock ledger for thread-safe access
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Retrieve audit logs
    logs := ledgerInstance.BlockchainConsensusCoinLedger.ConsensusAuditLogs
    if logs == nil || len(logs) == 0 {
        log.Printf("[Warning] No consensus audit logs found.")
        return nil, fmt.Errorf("no consensus audit logs available")
    }

    log.Printf("[Success] Retrieved %d consensus audit logs.", len(logs))
    return logs, nil
}


// consensusSetValidatorSelectionMode sets the validator selection mode.
func ConsensusSetValidatorSelectionMode(mode ledger.ValidatorSelectionMode, ledgerInstance *ledger.Ledger) error {
    // Input validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }
    if !isValidValidatorSelectionMode(mode) {
        return fmt.Errorf("invalid validator selection mode: %v", mode)
    }

    log.Printf("[Info] Setting validator selection mode: %v", mode)

    // Lock the ledger for thread-safe updates
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Attempt to set the validator selection mode
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetValidatorSelectionMode(mode); err != nil {
        return fmt.Errorf("failed to set validator selection mode: %w", err)
    }

    // Post-action verification
    currentMode, err := ledgerInstance.BlockchainConsensusCoinLedger.GetValidatorSelectionMode()
    if err != nil {
        return fmt.Errorf("post-verification failed: unable to retrieve validator selection mode: %w", err)
    }
    if currentMode != mode {
        return fmt.Errorf("post-verification failed: expected %v, got %v", mode, currentMode)
    }

    log.Printf("[Success] Validator selection mode set to: %v", mode)
    return nil
}

// Helper function: Validates the provided validator selection mode
func isValidValidatorSelectionMode(mode ledger.ValidatorSelectionMode) bool {
    validModes := []ledger.ValidatorSelectionMode{
        ledger.RandomSelection,
        ledger.StakeBasedSelection,
        ledger.RotationBasedSelection,
    }
    for _, validMode := range validModes {
        if mode == validMode {
            return true
        }
    }
    return false
}


// consensusGetValidatorSelectionMode retrieves the current validator selection mode.
func ConsensusGetValidatorSelectionMode(ledgerInstance *ledger.Ledger) (ledger.ValidatorSelectionMode, error) {
    // Input validation
    if ledgerInstance == nil {
        return ledger.ValidatorSelectionMode{}, fmt.Errorf("ledger instance is nil")
    }

    log.Printf("[Info] Retrieving current validator selection mode...")

    // Lock the ledger for thread-safe access
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Retrieve the current selection mode
    mode, err := ledgerInstance.BlockchainConsensusCoinLedger.GetValidatorSelectionMode()
    if err != nil {
        return ledger.ValidatorSelectionMode{}, fmt.Errorf("failed to retrieve validator selection mode: %w", err)
    }

    log.Printf("[Success] Retrieved validator selection mode: %v", mode)
    return mode, nil
}


// consensusSetPoHParticipationThreshold sets the participation threshold for Proof of History.
func ConsensusSetPoHParticipationThreshold(threshold float64, ledgerInstance *ledger.Ledger) error {
    // Input validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }
    if threshold < 0 || threshold > 1 {
        return fmt.Errorf("invalid threshold: %.2f, must be between 0 and 1", threshold)
    }

    log.Printf("[Info] Setting PoH participation threshold to %.2f", threshold)

    // Lock the ledger for thread-safe updates
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Attempt to set the PoH participation threshold
    if err := ledgerInstance.BlockchainConsensusCoinLedger.SetPoHParticipationThreshold(threshold); err != nil {
        return fmt.Errorf("failed to set PoH participation threshold: %w", err)
    }

    // Post-action verification
    currentThreshold, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPoHParticipationThreshold()
    if err != nil {
        return fmt.Errorf("post-verification failed: unable to retrieve PoH participation threshold: %w", err)
    }
    if currentThreshold != threshold {
        return fmt.Errorf("post-verification failed: expected %.2f, got %.2f", threshold, currentThreshold)
    }

    log.Printf("[Success] PoH participation threshold successfully set to %.2f", threshold)
    return nil
}


// consensusGetPoHParticipationThreshold retrieves the current PoH participation threshold.
func ConsensusGetPoHParticipationThreshold(ledgerInstance *ledger.Ledger) (float64, error) {
    // Input validation
    if ledgerInstance == nil {
        return 0, fmt.Errorf("ledger instance is nil")
    }

    log.Printf("[Info] Retrieving current PoH participation threshold...")

    // Lock the ledger for thread-safe access
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Retrieve the current threshold
    threshold, err := ledgerInstance.BlockchainConsensusCoinLedger.GetPoHParticipationThreshold()
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve PoH participation threshold: %w", err)
    }

    log.Printf("[Success] Retrieved PoH participation threshold: %.2f", threshold)
    return threshold, nil
}


// consensusValidateValidatorActivity validates recent activity of a validator.
func ConsensusValidateValidatorActivity(validatorID, action, details string, ledgerInstance *ledger.Ledger) error {
    // Input validation
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }
    if validatorID == "" {
        return fmt.Errorf("validator ID cannot be empty")
    }
    if action == "" {
        return fmt.Errorf("action cannot be empty")
    }

    log.Printf("[Info] Validating activity for validator: %s, Action: %s, Details: %s", validatorID, action, details)

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Log the validator's activity in the ledger
    if err := ledgerInstance.BlockchainConsensusCoinLedger.LogValidatorActivity(validatorID, action, details); err != nil {
        return fmt.Errorf("failed to log validator activity for %s: %w", validatorID, err)
    }

    // Post-action verification
    logs := ledgerInstance.BlockchainConsensusCoinLedger.ValidatorActivityLogs
    found := false
    for _, logEntry := range logs {
        if logEntry.ValidatorID == validatorID && logEntry.Action == action && logEntry.Details == details {
            found = true
            break
        }
    }
    if !found {
        return fmt.Errorf("post-verification failed: activity log for validator %s not found", validatorID)
    }

    log.Printf("[Success] Validator activity validated and logged for %s: Action: %s", validatorID, action)
    return nil
}


// consensusFetchValidatorActivityLogs fetches logs of validator activities.
func ConsensusFetchValidatorActivityLogs(ledgerInstance *ledger.Ledger) ([]ledger.ValidatorActivityLog, error) {
    // Input validation
    if ledgerInstance == nil {
        return nil, fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[Info] Fetching validator activity logs...")

    // Lock the ledger for thread-safe access
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Fetch logs from the ledger
    logs := ledgerInstance.BlockchainConsensusCoinLedger.ValidatorActivityLogs
    if logs == nil {
        return nil, fmt.Errorf("validator activity logs are empty")
    }

    log.Printf("[Success] Retrieved %d validator activity logs", len(logs))
    return logs, nil
}


// consensusEnableDynamicStakeAdjustment enables dynamic stake adjustment.
func ConsensusEnableDynamicStakeAdjustment(ledgerInstance *ledger.Ledger) error {
    // Validate input
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Enabling dynamic stake adjustment...")

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Enable dynamic stake adjustment
    if err := ledgerInstance.BlockchainConsensusCoinLedger.EnableDynamicStakeAdjustment(); err != nil {
        log.Printf("[ERROR] Failed to enable dynamic stake adjustment: %v", err)
        return fmt.Errorf("failed to enable dynamic stake adjustment: %w", err)
    }

    // Post-action verification
    if !ledgerInstance.BlockchainConsensusCoinLedger.IsDynamicStakeAdjustmentEnabled() {
        return fmt.Errorf("post-verification failed: dynamic stake adjustment not enabled")
    }

    log.Printf("[SUCCESS] Dynamic stake adjustment enabled.")
    return nil
}


// consensusDisableDynamicStakeAdjustment disables dynamic stake adjustment.
func ConsensusDisableDynamicStakeAdjustment(ledgerInstance *ledger.Ledger) error {
    // Validate input
    if ledgerInstance == nil {
        return fmt.Errorf("ledger instance cannot be nil")
    }

    log.Printf("[INFO] Disabling dynamic stake adjustment...")

    // Lock the ledger for thread-safe operations
    ledgerInstance.Lock()
    defer ledgerInstance.Unlock()

    // Disable dynamic stake adjustment
    if err := ledgerInstance.BlockchainConsensusCoinLedger.DisableDynamicStakeAdjustment(); err != nil {
        log.Printf("[ERROR] Failed to disable dynamic stake adjustment: %v", err)
        return fmt.Errorf("failed to disable dynamic stake adjustment: %w", err)
    }

    // Post-action verification
    if ledgerInstance.BlockchainConsensusCoinLedger.IsDynamicStakeAdjustmentEnabled() {
        return fmt.Errorf("post-verification failed: dynamic stake adjustment still enabled")
    }

    log.Printf("[SUCCESS] Dynamic stake adjustment disabled.")
    return nil
}

