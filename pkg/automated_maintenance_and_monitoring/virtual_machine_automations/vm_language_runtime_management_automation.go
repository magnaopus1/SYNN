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
    LanguageRuntimeCheckInterval = 4000 * time.Millisecond // Interval for checking runtime environments
    SubBlocksPerBlock            = 1000                    // Number of sub-blocks in a block
    MaxRetryAttempts             = 5                       // Max retry attempts for initializing runtime environments
)

// VMLanguageRuntimeManagementAutomation automates the process of managing multiple language runtimes for the VMs
type VMLanguageRuntimeManagementAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger to store runtime management actions
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    runtimeCheckCount    int                          // Counter for runtime check cycles
    failedRuntimes       map[string]int               // Map to track failed language runtimes and retry attempts
    supportedLanguages   []string                     // Supported language runtimes in the VM (Solidity, Yul, Rust, etc.)
}

// NewVMLanguageRuntimeManagementAutomation initializes the automation for managing language runtimes across VMs
func NewVMLanguageRuntimeManagementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMLanguageRuntimeManagementAutomation {
    return &VMLanguageRuntimeManagementAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        runtimeCheckCount:  0,
        failedRuntimes:     make(map[string]int),
        supportedLanguages: []string{"Solidity", "Yul", "Rust", "Golang", "JavaScript"},
    }
}

// StartRuntimeCheck starts the continuous loop for checking and managing language runtimes in the VMs
func (automation *VMLanguageRuntimeManagementAutomation) StartRuntimeCheck() {
    ticker := time.NewTicker(LanguageRuntimeCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndManageRuntimes()
        }
    }()
}

// monitorAndManageRuntimes checks the language runtimes and attempts to initialize or recover failed runtimes
func (automation *VMLanguageRuntimeManagementAutomation) monitorAndManageRuntimes() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    for _, lang := range automation.supportedLanguages {
        runtimeStatus := automation.consensusSystem.GetRuntimeStatus(lang)

        if !runtimeStatus.IsHealthy {
            fmt.Printf("Runtime for language %s is unhealthy. Attempting recovery.\n", lang)
            automation.triggerRuntimeRecovery(lang)
        } else {
            fmt.Printf("Runtime for language %s is healthy.\n", lang)
            automation.resetRuntimeFailure(lang)
        }
    }

    automation.runtimeCheckCount++
    fmt.Printf("Runtime check cycle #%d executed.\n", automation.runtimeCheckCount)

    if automation.runtimeCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeRuntimeCheckCycle()
    }
}

// triggerRuntimeRecovery attempts to recover a failed runtime
func (automation *VMLanguageRuntimeManagementAutomation) triggerRuntimeRecovery(language string) {
    retryCount, exists := automation.failedRuntimes[language]

    if exists && retryCount >= MaxRetryAttempts {
        fmt.Printf("Runtime for language %s has exceeded max recovery attempts. Marking as failed.\n", language)
        automation.markRuntimeAsFailed(language)
        return
    }

    success := automation.consensusSystem.AttemptRuntimeRecovery(language)

    if success {
        fmt.Printf("Runtime for language %s successfully recovered.\n", language)
        automation.resetRuntimeFailure(language)
        automation.logRuntimeRecovery(language)
    } else {
        automation.failedRuntimes[language] = retryCount + 1
        fmt.Printf("Recovery for language %s runtime failed. Retry attempt #%d.\n", language, automation.failedRuntimes[language])
    }
}

// markRuntimeAsFailed marks a language runtime as failed after exceeding retry attempts and logs the event
func (automation *VMLanguageRuntimeManagementAutomation) markRuntimeAsFailed(language string) {
    // Encrypt runtime data before marking it as failed
    encryptedRuntimeData := encryption.EncryptData([]byte(language))

    // Mark runtime as failed in Synnergy Consensus
    automation.consensusSystem.MarkRuntimeAsFailed(encryptedRuntimeData)

    automation.logRuntimeFailure(language)
}

// resetRuntimeFailure resets the failure counter for a healthy language runtime
func (automation *VMLanguageRuntimeManagementAutomation) resetRuntimeFailure(language string) {
    delete(automation.failedRuntimes, language)
}

// finalizeRuntimeCheckCycle finalizes the runtime check cycle and logs the result in the ledger
func (automation *VMLanguageRuntimeManagementAutomation) finalizeRuntimeCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeRuntimeCheckCycle()
    if success {
        fmt.Println("Runtime check cycle finalized successfully.")
        automation.logRuntimeCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing runtime check cycle.")
    }
}

// logRuntimeRecovery logs the successful recovery of a runtime into the ledger
func (automation *VMLanguageRuntimeManagementAutomation) logRuntimeRecovery(language string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("runtime-recovery-%s", language),
        Timestamp: time.Now().Unix(),
        Type:      "Runtime Recovery",
        Status:    "Recovered",
        Details:   fmt.Sprintf("Runtime for language %s successfully recovered.", language),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with runtime recovery for language %s.\n", language)
}

// logRuntimeFailure logs the failure of a runtime into the ledger for traceability
func (automation *VMLanguageRuntimeManagementAutomation) logRuntimeFailure(language string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("runtime-failure-%s", language),
        Timestamp: time.Now().Unix(),
        Type:      "Runtime Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Runtime for language %s marked as failed after exceeding recovery attempts.", language),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with runtime failure for language %s.\n", language)
}

// logRuntimeCheckCycleFinalization logs the finalization of a runtime check cycle into the ledger
func (automation *VMLanguageRuntimeManagementAutomation) logRuntimeCheckCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("runtime-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Runtime Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with runtime check cycle finalization.")
}

// ensureRuntimeIntegrity checks the integrity of language runtime data and triggers recovery if necessary
func (automation *VMLanguageRuntimeManagementAutomation) ensureRuntimeIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateRuntimeIntegrity()
    if !integrityValid {
        fmt.Println("Runtime integrity breach detected. Re-triggering runtime checks.")
        automation.monitorAndManageRuntimes()
    } else {
        fmt.Println("Runtime data integrity is valid.")
    }
}
