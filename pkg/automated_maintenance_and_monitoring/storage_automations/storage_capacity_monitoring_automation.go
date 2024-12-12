package automations

import (
    "fmt"
    "sync"
    "time"
    "os"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    StorageMonitoringInterval   = 1 * time.Hour  // Interval for checking storage capacity
    MaxStorageThreshold         = 500 * 1024 * 1024 * 1024 // 500 GB maximum threshold for on-chain storage
    WarningStorageThreshold     = 450 * 1024 * 1024 * 1024 // 450 GB threshold to issue warnings
    SubBlocksPerBlock           = 1000          // Number of sub-blocks in a block
    StorageRetryLimit           = 3             // Maximum retry attempts for handling storage overflow
)

// StorageCapacityMonitoringAutomation monitors on-chain storage usage and triggers appropriate actions
type StorageCapacityMonitoringAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger for logging storage-related events
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    storageRetryCount      map[string]int               // Counter for retrying storage overflow operations
    monitoringCycleCount   int                          // Counter for storage monitoring cycles
}

// NewStorageCapacityMonitoringAutomation initializes the storage monitoring automation
func NewStorageCapacityMonitoringAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *StorageCapacityMonitoringAutomation {
    return &StorageCapacityMonitoringAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        storageRetryCount:    make(map[string]int),
        monitoringCycleCount: 0,
    }
}

// StartStorageCapacityMonitoring starts the continuous monitoring of storage capacity
func (automation *StorageCapacityMonitoringAutomation) StartStorageCapacityMonitoring() {
    ticker := time.NewTicker(StorageMonitoringInterval)

    go func() {
        for range ticker.C {
            automation.monitorStorageCapacity()
        }
    }()
}

// monitorStorageCapacity checks the current storage usage and determines if any action is required
func (automation *StorageCapacityMonitoringAutomation) monitorStorageCapacity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    currentStorageUsage, err := automation.getStorageUsage()
    if err != nil {
        fmt.Printf("Error retrieving storage usage: %v. Retrying...\n", err)
        return
    }

    if currentStorageUsage >= MaxStorageThreshold {
        fmt.Printf("Storage capacity exceeded. Current usage: %d bytes. Initiating overflow handling.\n", currentStorageUsage)
        automation.handleStorageOverflow(currentStorageUsage)
    } else if currentStorageUsage >= WarningStorageThreshold {
        fmt.Printf("Warning: Storage usage nearing capacity. Current usage: %d bytes.\n", currentStorageUsage)
        automation.logStorageWarning(currentStorageUsage)
    } else {
        fmt.Printf("Storage usage is within normal limits. Current usage: %d bytes.\n", currentStorageUsage)
    }

    automation.monitoringCycleCount++
    fmt.Printf("Storage monitoring cycle #%d completed.\n", automation.monitoringCycleCount)

    if automation.monitoringCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeStorageMonitoringCycle()
    }
}

// getStorageUsage retrieves the current usage of on-chain storage
func (automation *StorageCapacityMonitoringAutomation) getStorageUsage() (int64, error) {
    var stat os.FileInfo
    stat, err := os.Stat("/data/onchain_storage") // Adjust path to the actual on-chain storage directory
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve on-chain storage info: %v", err)
    }
    return stat.Size(), nil
}

// handleStorageOverflow handles the scenario where on-chain storage exceeds the allowed threshold
func (automation *StorageCapacityMonitoringAutomation) handleStorageOverflow(currentStorageUsage int64) {
    err := automation.triggerStorageOptimization()
    if err != nil {
        fmt.Printf("Storage optimization failed: %v. Retrying...\n", err)
        automation.retryStorageOverflow(currentStorageUsage)
        return
    }

    fmt.Printf("Storage overflow handled successfully.\n")
    automation.logStorageEvent(currentStorageUsage, "Overflow handled", "Optimization triggered")
    automation.resetStorageRetry()
}

// triggerStorageOptimization optimizes or offloads data to handle storage overflow
func (automation *StorageCapacityMonitoringAutomation) triggerStorageOptimization() error {
    fmt.Println("Triggering on-chain storage optimization...")
    success := automation.consensusSystem.TriggerOnChainOptimization()

    if success {
        fmt.Println("On-chain optimization triggered successfully.")
        return nil
    }

    return fmt.Errorf("on-chain optimization failed")
}

// retryStorageOverflow retries the storage overflow handling process in case of failure
func (automation *StorageCapacityMonitoringAutomation) retryStorageOverflow(currentStorageUsage int64) {
    automation.storageRetryCount["overflow"]++
    if automation.storageRetryCount["overflow"] < StorageRetryLimit {
        automation.handleStorageOverflow(currentStorageUsage)
    } else {
        fmt.Printf("Max retries reached for storage overflow. Operation failed.\n")
        automation.logStorageFailure(currentStorageUsage)
    }
}

// resetStorageRetry resets the retry count for storage overflow handling
func (automation *StorageCapacityMonitoringAutomation) resetStorageRetry() {
    automation.storageRetryCount["overflow"] = 0
}

// logStorageWarning logs a warning when the storage usage is nearing its capacity
func (automation *StorageCapacityMonitoringAutomation) logStorageWarning(currentStorageUsage int64) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("storage-warning-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Capacity Warning",
        Status:    "Warning",
        Details:   fmt.Sprintf("Storage usage nearing capacity. Current usage: %d bytes.", currentStorageUsage),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with storage warning.\n")
}

// logStorageEvent logs a successful storage handling event into the ledger
func (automation *StorageCapacityMonitoringAutomation) logStorageEvent(currentStorageUsage int64, eventType, action string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("storage-event-%s-%d", eventType, time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Capacity Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Storage event triggered at %d bytes. Action: %s.", currentStorageUsage, action),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with storage event.\n")
}

// logStorageFailure logs the failure to handle storage overflow into the ledger
func (automation *StorageCapacityMonitoringAutomation) logStorageFailure(currentStorageUsage int64) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("storage-failure-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Capacity Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to handle storage overflow. Current usage: %d bytes.", currentStorageUsage),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with storage overflow failure.\n")
}

// finalizeStorageMonitoringCycle finalizes the storage monitoring cycle and logs the result in the ledger
func (automation *StorageCapacityMonitoringAutomation) finalizeStorageMonitoringCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeStorageCycle()
    if success {
        fmt.Println("Storage monitoring cycle finalized successfully.")
        automation.logStorageCycleFinalization()
    } else {
        fmt.Println("Error finalizing storage monitoring cycle.")
    }
}

// logStorageCycleFinalization logs the finalization of a storage monitoring cycle into the ledger
func (automation *StorageCapacityMonitoringAutomation) logStorageCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("storage-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Monitoring Cycle Finalization",
        Status:    "Finalized",
        Details:   "Storage monitoring cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with storage cycle finalization.")
}
