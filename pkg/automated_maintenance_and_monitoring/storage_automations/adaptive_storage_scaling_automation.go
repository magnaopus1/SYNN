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
    StorageMonitoringInterval      = 30 * time.Second // Interval for monitoring storage usage
    MaxScalingRetries              = 3                // Maximum retries for applying storage scaling
    SubBlocksPerBlock              = 1000             // Number of sub-blocks in a block
    StorageUsageThreshold          = 80               // Threshold percentage for triggering storage scaling
    AdaptiveStorageScalingLimit    = 90               // Limit for adaptive storage scaling
)

// AdaptiveStorageScalingAutomation monitors and adapts storage scaling based on real-time requirements
type AdaptiveStorageScalingAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger for logging storage scaling-related events
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    scalingRetryCount  map[string]int               // Counter for retrying storage scaling actions
    scalingCycleCount  int                          // Counter for storage scaling monitoring cycles
}

// NewAdaptiveStorageScalingAutomation initializes the adaptive storage scaling automation
func NewAdaptiveStorageScalingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AdaptiveStorageScalingAutomation {
    return &AdaptiveStorageScalingAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        scalingRetryCount:  make(map[string]int),
        scalingCycleCount:  0,
    }
}

// StartStorageMonitoring begins the continuous loop for monitoring storage usage
func (automation *AdaptiveStorageScalingAutomation) StartStorageMonitoring() {
    ticker := time.NewTicker(StorageMonitoringInterval)

    go func() {
        for range ticker.C {
            automation.monitorStorageUsage()
        }
    }()
}

// monitorStorageUsage checks the current storage usage and applies scaling if necessary
func (automation *AdaptiveStorageScalingAutomation) monitorStorageUsage() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch storage usage reports
    storageReports := automation.consensusSystem.FetchStorageUsageReports()

    for _, report := range storageReports {
        if automation.isScalingRequired(report) {
            fmt.Printf("Storage scaling required for node %s. Taking action.\n", report.NodeID)
            automation.applyStorageScaling(report)
        } else {
            fmt.Printf("Storage usage is within safe limits for node %s.\n", report.NodeID)
        }
    }

    automation.scalingCycleCount++
    fmt.Printf("Storage scaling monitoring cycle #%d completed.\n", automation.scalingCycleCount)

    if automation.scalingCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeScalingCycle()
    }
}

// isScalingRequired checks if storage scaling is required based on current usage
func (automation *AdaptiveStorageScalingAutomation) isScalingRequired(report common.StorageUsageReport) bool {
    if report.UsagePercentage >= StorageUsageThreshold && report.UsagePercentage <= AdaptiveStorageScalingLimit {
        fmt.Printf("Storage usage for node %s is %d%%, above the threshold.\n", report.NodeID, report.UsagePercentage)
        return true
    }
    return false
}

// applyStorageScaling attempts to scale storage resources for a given node
func (automation *AdaptiveStorageScalingAutomation) applyStorageScaling(report common.StorageUsageReport) {
    encryptedScalingData := automation.encryptScalingData(report)

    // Attempt to scale storage through the Synnergy Consensus system
    scalingSuccess := automation.consensusSystem.ApplyStorageScaling(report.NodeID, encryptedScalingData)

    if scalingSuccess {
        fmt.Printf("Storage scaling applied successfully for node %s.\n", report.NodeID)
        automation.logScalingEvent(report, "Storage Scaled")
        automation.resetScalingRetry(report.NodeID)
    } else {
        fmt.Printf("Error applying storage scaling for node %s. Retrying...\n", report.NodeID)
        automation.retryStorageScaling(report)
    }
}

// retryStorageScaling retries storage scaling if the initial attempt fails
func (automation *AdaptiveStorageScalingAutomation) retryStorageScaling(report common.StorageUsageReport) {
    automation.scalingRetryCount[report.NodeID]++
    if automation.scalingRetryCount[report.NodeID] < MaxScalingRetries {
        automation.applyStorageScaling(report)
    } else {
        fmt.Printf("Max retries reached for storage scaling on node %s. Scaling failed.\n", report.NodeID)
        automation.logScalingFailure(report)
    }
}

// resetScalingRetry resets the retry count for storage scaling actions
func (automation *AdaptiveStorageScalingAutomation) resetScalingRetry(nodeID string) {
    automation.scalingRetryCount[nodeID] = 0
}

// finalizeScalingCycle finalizes the storage scaling cycle and logs the result in the ledger
func (automation *AdaptiveStorageScalingAutomation) finalizeScalingCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeStorageScalingCycle()
    if success {
        fmt.Println("Storage scaling cycle finalized successfully.")
        automation.logScalingCycleFinalization()
    } else {
        fmt.Println("Error finalizing storage scaling cycle.")
    }
}

// logScalingEvent logs a storage scaling event into the ledger
func (automation *AdaptiveStorageScalingAutomation) logScalingEvent(report common.StorageUsageReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("storage-scaling-%s-%s", report.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Scaling Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s triggered %s due to storage usage issues.", report.NodeID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with storage scaling event for node %s.\n", report.NodeID)
}

// logScalingFailure logs the failure of a storage scaling attempt into the ledger
func (automation *AdaptiveStorageScalingAutomation) logScalingFailure(report common.StorageUsageReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("storage-scaling-failure-%s", report.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Scaling Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Storage scaling failed for node %s after maximum retries.", report.NodeID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with storage scaling failure for node %s.\n", report.NodeID)
}

// logScalingCycleFinalization logs the finalization of a storage scaling cycle into the ledger
func (automation *AdaptiveStorageScalingAutomation) logScalingCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("storage-scaling-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Scaling Cycle Finalization",
        Status:    "Finalized",
        Details:   "Storage scaling cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with storage scaling cycle finalization.")
}

// encryptScalingData encrypts storage scaling-related data before applying scaling or logging events
func (automation *AdaptiveStorageScalingAutomation) encryptScalingData(report common.StorageUsageReport) common.StorageUsageReport {
    encryptedData, err := encryption.EncryptData(report.UsageData)
    if err != nil {
        fmt.Println("Error encrypting storage usage data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("Storage usage data successfully encrypted for node:", report.NodeID)
    return report
}
