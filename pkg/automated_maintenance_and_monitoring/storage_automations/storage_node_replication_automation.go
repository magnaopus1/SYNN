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
    ReplicationInterval     = 1 * time.Hour  // Interval for checking replication status
    MaxReplicationRetries   = 3              // Maximum retry attempts for failed replications
    SubBlocksPerBlock       = 1000           // Number of sub-blocks in a block
)

// StorageNodeReplicationAutomation handles the replication of storage nodes to ensure data redundancy
type StorageNodeReplicationAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance      *ledger.Ledger               // Ledger for logging replication events
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    replicationRetryMap map[string]int               // Counter for retrying failed replication operations
    replicationCycle    int                          // Counter for replication cycles
}

// NewStorageNodeReplicationAutomation initializes the automation for storage node replication
func NewStorageNodeReplicationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *StorageNodeReplicationAutomation {
    return &StorageNodeReplicationAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        replicationRetryMap: make(map[string]int),
        replicationCycle:   0,
    }
}

// StartNodeReplication starts the continuous loop for replicating storage nodes
func (automation *StorageNodeReplicationAutomation) StartNodeReplication() {
    ticker := time.NewTicker(ReplicationInterval)

    go func() {
        for range ticker.C {
            automation.checkAndReplicateStorageNodes()
        }
    }()
}

// checkAndReplicateStorageNodes checks the replication status of storage nodes and replicates them if needed
func (automation *StorageNodeReplicationAutomation) checkAndReplicateStorageNodes() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    replicationReports := automation.consensusSystem.FetchReplicationReports()

    for _, report := range replicationReports {
        if report.IsReplicationNeeded {
            fmt.Printf("Replication needed for node %s. Initiating replication.\n", report.NodeID)
            automation.replicateNode(report)
        } else {
            fmt.Printf("Replication not needed for node %s.\n", report.NodeID)
        }
    }

    automation.replicationCycle++
    fmt.Printf("Storage node replication cycle #%d completed.\n", automation.replicationCycle)

    if automation.replicationCycle%SubBlocksPerBlock == 0 {
        automation.finalizeReplicationCycle()
    }
}

// replicateNode handles the replication of storage data to other nodes
func (automation *StorageNodeReplicationAutomation) replicateNode(report common.ReplicationReport) {
    encryptedData := automation.encryptReplicationData(report)

    err := automation.triggerReplication(report, encryptedData)
    if err != nil {
        fmt.Printf("Replication failed for node %s: %v. Retrying...\n", report.NodeID, err)
        automation.retryReplication(report)
        return
    }

    fmt.Printf("Replication successful for node %s.\n", report.NodeID)
    automation.logReplicationEvent(report, "Replicated")
    automation.resetReplicationRetry(report.NodeID)
}

// triggerReplication triggers the actual replication process through consensus
func (automation *StorageNodeReplicationAutomation) triggerReplication(report common.ReplicationReport, encryptedData []byte) error {
    fmt.Println("Triggering storage node replication...")

    success := automation.consensusSystem.ReplicateNode(report, encryptedData)
    if success {
        fmt.Println("Node replication triggered successfully.")
        return nil
    }

    return fmt.Errorf("node replication failed for node %s", report.NodeID)
}

// retryReplication retries the replication process in case of failure
func (automation *StorageNodeReplicationAutomation) retryReplication(report common.ReplicationReport) {
    automation.replicationRetryMap[report.NodeID]++
    if automation.replicationRetryMap[report.NodeID] < MaxReplicationRetries {
        automation.replicateNode(report)
    } else {
        fmt.Printf("Max retries reached for replicating node %s. Replication failed.\n", report.NodeID)
        automation.logReplicationFailure(report)
    }
}

// resetReplicationRetry resets the retry count for replication operations
func (automation *StorageNodeReplicationAutomation) resetReplicationRetry(nodeID string) {
    automation.replicationRetryMap[nodeID] = 0
}

// encryptReplicationData encrypts the storage data before replication
func (automation *StorageNodeReplicationAutomation) encryptReplicationData(report common.ReplicationReport) []byte {
    encryptedData, err := encryption.EncryptData(report.Data)
    if err != nil {
        fmt.Printf("Error encrypting data for node %s replication: %v\n", report.NodeID, err)
        return report.Data
    }

    fmt.Printf("Data successfully encrypted for replication of node %s.\n", report.NodeID)
    return encryptedData
}

// logReplicationEvent logs a successful replication event into the ledger
func (automation *StorageNodeReplicationAutomation) logReplicationEvent(report common.ReplicationReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("replication-event-%s-%s", report.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Node Replication Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Replication %s for node %s.", eventType, report.NodeID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with replication event for node %s.\n", report.NodeID)
}

// logReplicationFailure logs the failure to replicate a specific node into the ledger
func (automation *StorageNodeReplicationAutomation) logReplicationFailure(report common.ReplicationReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("replication-failure-%s", report.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Node Replication Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to replicate node %s after maximum retries.", report.NodeID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with replication failure for node %s.\n", report.NodeID)
}

// finalizeReplicationCycle finalizes the storage replication cycle and logs the result in the ledger
func (automation *StorageNodeReplicationAutomation) finalizeReplicationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeReplicationCycle()
    if success {
        fmt.Println("Storage replication cycle finalized successfully.")
        automation.logReplicationCycleFinalization()
    } else {
        fmt.Println("Error finalizing storage replication cycle.")
    }
}

// logReplicationCycleFinalization logs the finalization of a replication cycle into the ledger
func (automation *StorageNodeReplicationAutomation) logReplicationCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("replication-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Replication Cycle Finalization",
        Status:    "Finalized",
        Details:   "Storage replication cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with replication cycle finalization.")
}
