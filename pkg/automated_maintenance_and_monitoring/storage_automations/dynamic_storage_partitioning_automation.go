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
    PartitioningCheckInterval     = 20 * time.Second // Interval for checking partitioning requirements
    MaxPartitioningRetries        = 3                // Maximum retries for partitioning issues
    SubBlocksPerBlock             = 1000             // Number of sub-blocks in a block
    StoragePartitioningThreshold  = 85               // Storage usage percentage threshold for partitioning
)

// DynamicStoragePartitioningAutomation manages real-time partitioning of storage across nodes
type DynamicStoragePartitioningAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger for logging partitioning-related events
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    partitionRetryCount map[string]int              // Counter for retrying partitioning actions
    partitionCycleCount int                         // Counter for partitioning monitoring cycles
}

// NewDynamicStoragePartitioningAutomation initializes the dynamic storage partitioning automation
func NewDynamicStoragePartitioningAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DynamicStoragePartitioningAutomation {
    return &DynamicStoragePartitioningAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        partitionRetryCount: make(map[string]int),
        partitionCycleCount: 0,
    }
}

// StartPartitioningMonitoring starts the continuous loop for monitoring storage partitioning
func (automation *DynamicStoragePartitioningAutomation) StartPartitioningMonitoring() {
    ticker := time.NewTicker(PartitioningCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorStoragePartitioning()
        }
    }()
}

// monitorStoragePartitioning checks for nodes that require storage partitioning and handles the process
func (automation *DynamicStoragePartitioningAutomation) monitorStoragePartitioning() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch storage partitioning reports
    partitionReports := automation.consensusSystem.FetchStoragePartitioningReports()

    for _, report := range partitionReports {
        if automation.isPartitioningRequired(report) {
            fmt.Printf("Storage partitioning required for node %s. Initiating partitioning process.\n", report.NodeID)
            automation.applyStoragePartitioning(report)
        } else {
            fmt.Printf("No partitioning required for node %s.\n", report.NodeID)
        }
    }

    automation.partitionCycleCount++
    fmt.Printf("Storage partitioning cycle #%d completed.\n", automation.partitionCycleCount)

    if automation.partitionCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizePartitioningCycle()
    }
}

// isPartitioningRequired checks if a node requires storage partitioning based on usage
func (automation *DynamicStoragePartitioningAutomation) isPartitioningRequired(report common.StoragePartitioningReport) bool {
    if report.UsagePercentage >= StoragePartitioningThreshold {
        fmt.Printf("Node %s storage usage at %d%%, requiring partitioning.\n", report.NodeID, report.UsagePercentage)
        return true
    }
    return false
}

// applyStoragePartitioning attempts to partition storage for a node that exceeds the threshold
func (automation *DynamicStoragePartitioningAutomation) applyStoragePartitioning(report common.StoragePartitioningReport) {
    encryptedPartitioningData := automation.encryptPartitioningData(report)

    // Attempt to partition storage through the Synnergy Consensus system
    partitioningSuccess := automation.consensusSystem.ApplyStoragePartitioning(report.NodeID, encryptedPartitioningData)

    if partitioningSuccess {
        fmt.Printf("Storage partitioning successfully applied for node %s.\n", report.NodeID)
        automation.logPartitioningEvent(report, "Partitioning Successful")
        automation.resetPartitioningRetry(report.NodeID)
    } else {
        fmt.Printf("Error partitioning storage for node %s. Retrying...\n", report.NodeID)
        automation.retryStoragePartitioning(report)
    }
}

// retryStoragePartitioning retries the partitioning process in case of failure
func (automation *DynamicStoragePartitioningAutomation) retryStoragePartitioning(report common.StoragePartitioningReport) {
    automation.partitionRetryCount[report.NodeID]++
    if automation.partitionRetryCount[report.NodeID] < MaxPartitioningRetries {
        automation.applyStoragePartitioning(report)
    } else {
        fmt.Printf("Max retries reached for partitioning storage on node %s. Partitioning failed.\n", report.NodeID)
        automation.logPartitioningFailure(report)
    }
}

// resetPartitioningRetry resets the retry count for storage partitioning actions
func (automation *DynamicStoragePartitioningAutomation) resetPartitioningRetry(nodeID string) {
    automation.partitionRetryCount[nodeID] = 0
}

// finalizePartitioningCycle finalizes the partitioning cycle and logs the result in the ledger
func (automation *DynamicStoragePartitioningAutomation) finalizePartitioningCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizePartitioningCycle()
    if success {
        fmt.Println("Partitioning cycle finalized successfully.")
        automation.logPartitioningCycleFinalization()
    } else {
        fmt.Println("Error finalizing partitioning cycle.")
    }
}

// logPartitioningEvent logs a storage partitioning event into the ledger
func (automation *DynamicStoragePartitioningAutomation) logPartitioningEvent(report common.StoragePartitioningReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("partitioning-event-%s-%s", report.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Partitioning Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s triggered %s due to storage partitioning requirements.", report.NodeID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with partitioning event for node %s.\n", report.NodeID)
}

// logPartitioningFailure logs the failure of a storage partitioning attempt into the ledger
func (automation *DynamicStoragePartitioningAutomation) logPartitioningFailure(report common.StoragePartitioningReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("partitioning-failure-%s", report.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Partitioning Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Storage partitioning failed for node %s after maximum retries.", report.NodeID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with partitioning failure for node %s.\n", report.NodeID)
}

// logPartitioningCycleFinalization logs the finalization of a partitioning cycle into the ledger
func (automation *DynamicStoragePartitioningAutomation) logPartitioningCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("partitioning-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Partitioning Cycle Finalization",
        Status:    "Finalized",
        Details:   "Storage partitioning cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with partitioning cycle finalization.")
}

// encryptPartitioningData encrypts storage partitioning data before taking action or logging events
func (automation *DynamicStoragePartitioningAutomation) encryptPartitioningData(report common.StoragePartitioningReport) common.StoragePartitioningReport {
    encryptedData, err := encryption.EncryptData(report.PartitioningData)
    if err != nil {
        fmt.Println("Error encrypting storage partitioning data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("Storage partitioning data successfully encrypted for node:", report.NodeID)
    return report
}
