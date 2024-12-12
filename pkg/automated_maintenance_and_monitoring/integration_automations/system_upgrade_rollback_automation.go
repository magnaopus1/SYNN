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
    RollbackCheckInterval = 5000 * time.Millisecond // Interval for checking system stability after upgrades
    SubBlocksPerBlock     = 1000                    // Number of sub-blocks in a block
)

// SystemUpgradeRollbackAutomation automates the process of rolling back system upgrades in case of failure
type SystemUpgradeRollbackAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger to store rollback actions
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    rollbackCheckCount int                          // Counter for rollback check cycles
}

// NewSystemUpgradeRollbackAutomation initializes the automation for system upgrade rollbacks
func NewSystemUpgradeRollbackAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SystemUpgradeRollbackAutomation {
    return &SystemUpgradeRollbackAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        rollbackCheckCount: 0,
    }
}

// StartRollbackCheck starts the continuous loop for monitoring system stability after upgrades
func (automation *SystemUpgradeRollbackAutomation) StartRollbackCheck() {
    ticker := time.NewTicker(RollbackCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorSystemStability()
        }
    }()
}

// monitorSystemStability checks system stability and triggers rollback if an upgrade fails
func (automation *SystemUpgradeRollbackAutomation) monitorSystemStability() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch system stability data after upgrade
    systemStability, err := automation.consensusSystem.CheckSystemStability()
    if err != nil {
        fmt.Printf("Error checking system stability: %v\n", err)
        return
    }

    if !systemStability.IsStable {
        fmt.Println("System instability detected. Initiating rollback.")
        automation.initiateRollback(systemStability.FailureDetails)
    } else {
        fmt.Println("System is stable after upgrade.")
    }

    automation.rollbackCheckCount++
    fmt.Printf("Rollback check cycle #%d executed.\n", automation.rollbackCheckCount)

    if automation.rollbackCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeRollbackCheckCycle()
    }
}

// initiateRollback triggers the rollback process for failed system upgrades
func (automation *SystemUpgradeRollbackAutomation) initiateRollback(failureDetails string) {
    fmt.Println("Encrypting rollback data for failed system upgrade.")

    // Encrypt rollback details before proceeding
    encryptedDetails, err := encryption.EncryptData(failureDetails)
    if err != nil {
        fmt.Printf("Error encrypting rollback data: %v\n", err)
        return
    }

    // Trigger rollback in the Synnergy Consensus
    success := automation.consensusSystem.TriggerRollback(encryptedDetails)
    if success {
        fmt.Println("Rollback successfully triggered.")
        automation.logRollbackEvent(failureDetails)
    } else {
        fmt.Println("Error triggering rollback.")
        automation.logRollbackEvent("Rollback Failed: " + failureDetails)
    }
}

// logRollbackEvent logs the result of the rollback action into the ledger for traceability
func (automation *SystemUpgradeRollbackAutomation) logRollbackEvent(details string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("rollback-event-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "System Rollback",
        Status:    "Triggered",
        Details:   fmt.Sprintf("System rollback triggered. Details: %s", details),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with rollback event details.")
}

// finalizeRollbackCheckCycle finalizes the rollback check cycle and logs the result in the ledger
func (automation *SystemUpgradeRollbackAutomation) finalizeRollbackCheckCycle() {
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
func (automation *SystemUpgradeRollbackAutomation) logRollbackCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("rollback-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Rollback Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with rollback cycle finalization.")
}
