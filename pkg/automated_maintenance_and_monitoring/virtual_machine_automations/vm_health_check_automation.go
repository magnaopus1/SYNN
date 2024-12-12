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
    HealthCheckInterval    = 5000 * time.Millisecond // Interval for checking VM health
    SubBlocksPerBlock      = 1000                    // Number of sub-blocks in a block
    MaxRetryAttempts       = 3                       // Max retry attempts before marking VM as failed
    HealthRecoveryInterval = 10000 * time.Millisecond // Interval for attempting VM recovery
)

// VMHealthCheckAutomation automates the process of monitoring and ensuring VM health in the blockchain
type VMHealthCheckAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance   *ledger.Ledger               // Ledger to store VM health events
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    healthCheckCount int                          // Counter for health check cycles
    failedVMs        map[string]int               // Map to track VMs and retry attempts
}

// NewVMHealthCheckAutomation initializes the automation for VM health check
func NewVMHealthCheckAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMHealthCheckAutomation {
    return &VMHealthCheckAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        healthCheckCount: 0,
        failedVMs:        make(map[string]int),
    }
}

// StartHealthCheck starts the continuous loop for monitoring and ensuring VM health
func (automation *VMHealthCheckAutomation) StartHealthCheck() {
    ticker := time.NewTicker(HealthCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndCheckVMHealth()
        }
    }()
}

// monitorAndCheckVMHealth checks the health of all VMs and triggers recovery if needed
func (automation *VMHealthCheckAutomation) monitorAndCheckVMHealth() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of all VMs and their health status
    vmHealthStats := automation.consensusSystem.GetVMHealthStatistics()

    for _, vmStat := range vmHealthStats {
        if !vmStat.IsHealthy {
            fmt.Printf("VM %s is unhealthy. Triggering recovery attempts.\n", vmStat.ID)
            automation.triggerVMRecovery(vmStat)
        } else {
            fmt.Printf("VM %s is healthy.\n", vmStat.ID)
            automation.resetVMFailure(vmStat.ID)
        }
    }

    automation.healthCheckCount++
    fmt.Printf("Health check cycle #%d executed.\n", automation.healthCheckCount)

    if automation.healthCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeHealthCheckCycle()
    }
}

// triggerVMRecovery attempts to recover an unhealthy VM
func (automation *VMHealthCheckAutomation) triggerVMRecovery(vm common.VirtualMachine) {
    retryCount, exists := automation.failedVMs[vm.ID]

    if exists && retryCount >= MaxRetryAttempts {
        fmt.Printf("VM %s has exceeded max recovery attempts. Marking as failed.\n", vm.ID)
        automation.markVMAsFailed(vm)
        return
    }

    success := automation.consensusSystem.AttemptVMRecovery(vm)

    if success {
        fmt.Printf("VM %s successfully recovered.\n", vm.ID)
        automation.resetVMFailure(vm.ID)
        automation.logVMRecovery(vm)
    } else {
        automation.failedVMs[vm.ID] = retryCount + 1
        fmt.Printf("VM %s recovery failed. Retry attempt #%d.\n", vm.ID, automation.failedVMs[vm.ID])
    }
}

// markVMAsFailed marks a VM as failed after exceeding retry attempts and logs the event in the ledger
func (automation *VMHealthCheckAutomation) markVMAsFailed(vm common.VirtualMachine) {
    // Encrypt VM data before marking it as failed
    encryptedVMData := encryption.EncryptVMData(vm)

    // Trigger the failure protocol in SynnergyConsensus
    automation.consensusSystem.MarkVMAsFailed(encryptedVMData)

    automation.logVMFailure(vm)
}

// resetVMFailure resets the failure counter for a healthy VM
func (automation *VMHealthCheckAutomation) resetVMFailure(vmID string) {
    delete(automation.failedVMs, vmID)
}

// finalizeHealthCheckCycle finalizes the health check cycle and logs the result in the ledger
func (automation *VMHealthCheckAutomation) finalizeHealthCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeHealthCheckCycle()
    if success {
        fmt.Println("Health check cycle finalized successfully.")
        automation.logHealthCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing health check cycle.")
    }
}

// logVMRecovery logs the successful recovery of a VM into the ledger
func (automation *VMHealthCheckAutomation) logVMRecovery(vm common.VirtualMachine) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-recovery-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Recovery",
        Status:    "Recovered",
        Details:   fmt.Sprintf("VM %s successfully recovered.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with VM recovery for %s.\n", vm.ID)
}

// logVMFailure logs the failure of a VM into the ledger for traceability
func (automation *VMHealthCheckAutomation) logVMFailure(vm common.VirtualMachine) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-failure-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("VM %s marked as failed after exceeding recovery attempts.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with VM failure for %s.\n", vm.ID)
}

// logHealthCheckCycleFinalization logs the finalization of a health check cycle into the ledger
func (automation *VMHealthCheckAutomation) logHealthCheckCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("health-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Health Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with health check cycle finalization.")
}

// ensureVMHealthIntegrity checks the integrity of VM health data and triggers recovery if necessary
func (automation *VMHealthCheckAutomation) ensureVMHealthIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateVMHealthIntegrity()
    if !integrityValid {
        fmt.Println("VM health data integrity breach detected. Re-triggering health checks.")
        automation.monitorAndCheckVMHealth()
    } else {
        fmt.Println("VM health data integrity is valid.")
    }
}
