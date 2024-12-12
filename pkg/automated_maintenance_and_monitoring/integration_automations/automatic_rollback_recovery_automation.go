package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    RollbackCheckInterval = 3000 * time.Millisecond // Interval for checking rollback/recovery conditions
    SubBlocksPerBlock     = 1000                    // Number of sub-blocks per block
)

// AutomaticRollbackRecoveryAutomation handles automatic rollback and recovery in case of errors or unauthorized changes
type AutomaticRollbackRecoveryAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance   *ledger.Ledger               // Ledger to log rollback and recovery events
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    rollbackCount    int                          // Counter for rollback cycles
}

// NewAutomaticRollbackRecoveryAutomation initializes the automation for rollback and recovery
func NewAutomaticRollbackRecoveryAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AutomaticRollbackRecoveryAutomation {
    return &AutomaticRollbackRecoveryAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        rollbackCount:   0,
    }
}

// StartRollbackRecoveryMonitoring starts the continuous loop for monitoring rollback conditions
func (automation *AutomaticRollbackRecoveryAutomation) StartRollbackRecoveryMonitoring() {
    ticker := time.NewTicker(RollbackCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkForRollbackConditions()
        }
    }()
}

// checkForRollbackConditions checks for conditions that may require rollback or recovery
func (automation *AutomaticRollbackRecoveryAutomation) checkForRollbackConditions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch any new errors, unauthorized changes, or failures from the consensus system
    rollbackNeeded := automation.consensusSystem.CheckForRollbackTrigger()

    if rollbackNeeded {
        fmt.Println("Rollback conditions detected. Initiating rollback.")
        success := automation.initiateRollbackProcess()

        if success {
            fmt.Println("Rollback process completed successfully.")
            automation.logRollbackEvent("Success")
        } else {
            fmt.Println("Error during rollback process.")
            automation.logRollbackEvent("Failed")
        }
    } else {
        fmt.Println("No rollback conditions detected. System stable.")
    }

    automation.rollbackCount++
    fmt.Printf("Rollback check cycle #%d executed.\n", automation.rollbackCount)

    if automation.rollbackCount%SubBlocksPerBlock == 0 {
        automation.finalizeRollbackCheckCycle()
    }
}

// initiateRollbackProcess initiates the rollback process to revert to a stable state
func (automation *AutomaticRollbackRecoveryAutomation) initiateRollbackProcess() bool {
    // Simulate rollback process, this could involve reversing blocks, removing unauthorized changes, etc.
    rollbackSuccess := automation.consensusSystem.RevertToLastStableState()

    // Encrypt any sensitive data during rollback
    if rollbackSuccess {
        err := automation.encryptRollbackData()
        if err != nil {
            fmt.Println("Error encrypting rollback data:", err)
            return false
        }
        fmt.Println("Rollback data successfully encrypted.")
    }

    return rollbackSuccess
}

// encryptRollbackData encrypts sensitive data during the rollback process
func (automation *AutomaticRollbackRecoveryAutomation) encryptRollbackData() error {
    sensitiveData := automation.consensusSystem.GetSensitiveRollbackData()

    encryptedData, err := encryption.EncryptData(sensitiveData)
    if err != nil {
        return err
    }

    automation.consensusSystem.StoreEncryptedData(encryptedData)
    return nil
}

// logRollbackEvent logs rollback events in the ledger for audit and traceability
func (automation *AutomaticRollbackRecoveryAutomation) logRollbackEvent(result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("rollback-event-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Rollback Event",
        Status:    result,
        Details:   fmt.Sprintf("Rollback event with result: %s", result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with rollback event result: %s.\n", result)
}

// finalizeRollbackCheckCycle finalizes the rollback check cycle and logs the result in the ledger
func (automation *AutomaticRollbackRecoveryAutomation) finalizeRollbackCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeRollbackCycle()
    if success {
        fmt.Println("Rollback check cycle finalized successfully.")
        automation.logRollbackCycleFinalization()
    } else {
        fmt.Println("Error finalizing rollback check cycle.")
    }
}

// logRollbackCycleFinalization logs the finalization of a rollback check cycle into the ledger
func (automation *AutomaticRollbackRecoveryAutomation) logRollbackCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("rollback-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Rollback Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with rollback cycle finalization.")
}
