package consensus_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    ValidatorCheckInterval        = 5 * time.Minute  // Interval for checking validator performance
    ValidatorRotationEncryptionKey = "validator_rotation_key" // Encryption key for logging validator rotation events
    MinimumValidatorPerformance   = 0.8             // Minimum performance threshold for validators
    MaxValidatorLoad              = 0.9             // Maximum load threshold before rotation
)

// ValidatorRotationAutomation handles automatic rotation of PoS validators in the Synnergy Consensus system
type ValidatorRotationAutomation struct {
    ledgerInstance  *ledger.Ledger                    // Blockchain ledger for tracking validator rotation events
    consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validator performance checks
    stateMutex      *sync.RWMutex                     // Mutex for thread-safe ledger access
    validatorPool   []synnergy_consensus.Validator     // Pool of available validators for rotation
}

// NewValidatorRotationAutomation initializes the automation for validator rotation
func NewValidatorRotationAutomation(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex, validatorPool []synnergy_consensus.Validator) *ValidatorRotationAutomation {
    return &ValidatorRotationAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        validatorPool:   validatorPool,
    }
}

// StartValidatorRotation begins the continuous monitoring and rotation of validators based on performance
func (automation *ValidatorRotationAutomation) StartValidatorRotation() {
    ticker := time.NewTicker(ValidatorCheckInterval)

    go func() {
        for range ticker.C {
            fmt.Println("Monitoring validators for performance and load...")
            automation.checkAndRotateValidators()
        }
    }()
}

// checkAndRotateValidators checks the performance of all active validators and rotates out underperformers
func (automation *ValidatorRotationAutomation) checkAndRotateValidators() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the current list of active validators from the consensus engine
    activeValidators := automation.consensusEngine.GetActiveValidators()

    // Iterate over each active validator to check performance and load
    for _, validator := range activeValidators {
        performance := automation.consensusEngine.GetValidatorPerformance(validator)
        load := automation.consensusEngine.GetValidatorLoad(validator)

        if performance < MinimumValidatorPerformance || load > MaxValidatorLoad {
            fmt.Printf("Validator %s is underperforming or overloaded (Performance: %.2f, Load: %.2f). Rotating...\n", validator.ID, performance, load)
            automation.rotateValidator(validator)
        }
    }
}

// rotateValidator replaces the given underperforming validator with a new one from the pool
func (automation *ValidatorRotationAutomation) rotateValidator(validator synnergy_consensus.Validator) {
    // Remove the underperforming validator from the active pool
    automation.consensusEngine.RemoveValidator(validator)

    // Select a new validator from the available pool
    newValidator := automation.selectNewValidator()

    // Ensure the new validator has the required stake before adding
    if newValidator.Stake < automation.consensusEngine.GetRequiredPoSStake() {
        fmt.Printf("Validator %s does not meet the minimum stake requirement. Adjusting stake...\n", newValidator.ID)
        newValidator.Stake = automation.consensusEngine.GetRequiredPoSStake()
    }

    // Add the new validator to the active pool
    automation.consensusEngine.AddValidator(newValidator)
    fmt.Printf("Validator %s has been added to the active pool.\n", newValidator.ID)

    // Log the rotation event in the ledger
    automation.logValidatorRotation(validator, newValidator)
}

// selectNewValidator picks a new validator from the available pool
func (automation *ValidatorRotationAutomation) selectNewValidator() synnergy_consensus.Validator {
    for _, validator := range automation.validatorPool {
        if !automation.consensusEngine.IsValidatorActive(validator) {
            return validator
        }
    }
    fmt.Println("No available validators in the pool. Re-using an inactive validator.")
    return automation.validatorPool[0] // Default to the first validator if no others are available
}

// logValidatorRotation logs the rotation event into the ledger securely
func (automation *ValidatorRotationAutomation) logValidatorRotation(oldValidator, newValidator synnergy_consensus.Validator) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    rotationLog := common.LedgerEntry{
        ID:        fmt.Sprintf("validator-rotation-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Validator Rotation",
        Status:    "Completed",
        Details:   fmt.Sprintf("Validator %s replaced by %s due to performance/load issues.", oldValidator.ID, newValidator.ID),
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(rotationLog, []byte(ValidatorRotationEncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting validator rotation log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Validator rotation log stored in the ledger.")
}

// MonitorValidatorPool ensures the validator pool is healthy and replenishes if needed
func (automation *ValidatorRotationAutomation) MonitorValidatorPool() {
    fmt.Println("Monitoring validator pool for availability...")

    if len(automation.validatorPool) == 0 {
        fmt.Println("Validator pool is empty! Replenishing validator pool from the available validators...")
        automation.validatorPool = automation.consensusEngine.GetAvailableValidators()
    }
}

// ForceRotateValidator triggers a forced rotation of a specific validator, typically used for maintenance or manual overrides
func (automation *ValidatorRotationAutomation) ForceRotateValidator(validator synnergy_consensus.Validator) {
    fmt.Printf("Force rotating validator %s...\n", validator.ID)
    automation.rotateValidator(validator)
}
