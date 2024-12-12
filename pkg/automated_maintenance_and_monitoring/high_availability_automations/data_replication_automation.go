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
    DataReplicationCheckInterval = 3000 * time.Millisecond // Interval for checking data replication across nodes
    SubBlocksPerBlock            = 1000                    // Number of sub-blocks in a block
)

// DataReplicationAutomation automates the replication of data across all nodes for high availability and fault tolerance
type DataReplicationAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store data replication actions
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    replicationCheckCount  int                          // Counter for replication check cycles
}

// NewDataReplicationAutomation initializes the automation for data replication across nodes
func NewDataReplicationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataReplicationAutomation {
    return &DataReplicationAutomation{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        replicationCheckCount: 0,
    }
}

// StartDataReplicationCheck starts the continuous loop for monitoring and enforcing data replication across nodes
func (automation *DataReplicationAutomation) StartDataReplicationCheck() {
    ticker := time.NewTicker(DataReplicationCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceDataReplication()
        }
    }()
}

// monitorAndEnforceDataReplication checks if data is replicated across all nodes and enforces replication if needed
func (automation *DataReplicationAutomation) monitorAndEnforceDataReplication() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch all data items that should be replicated across nodes
    dataItems := automation.consensusSystem.GetDataToReplicate()

    for _, dataItem := range dataItems {
        if automation.checkDataReplication(dataItem) {
            fmt.Printf("Data item %s is properly replicated across all nodes.\n", dataItem.ID)
        } else {
            fmt.Printf("Data item %s is not properly replicated. Enforcing replication to remaining nodes.\n", dataItem.ID)
            automation.replicateDataToNodes(dataItem)
        }
    }

    automation.replicationCheckCount++
    fmt.Printf("Data replication check cycle #%d executed.\n", automation.replicationCheckCount)

    if automation.replicationCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeReplicationCycle()
    }
}

// checkDataReplication checks if the data item is replicated across all nodes
func (automation *DataReplicationAutomation) checkDataReplication(dataItem common.DataItem) bool {
    return automation.consensusSystem.IsDataReplicatedAcrossNodes(dataItem.ID)
}

// replicateDataToNodes enforces data replication to any nodes where the data is missing
func (automation *DataReplicationAutomation) replicateDataToNodes(dataItem common.DataItem) {
    // Encrypt the data item before replication
    encryptedDataItem := automation.AddEncryptionToDataItem(dataItem)

    // Replicate the data item to all nodes via the Synnergy Consensus
    replicationSuccess := automation.consensusSystem.ReplicateDataAcrossNodes(encryptedDataItem)

    if replicationSuccess {
        fmt.Printf("Data item %s successfully replicated across all nodes.\n", dataItem.ID)
        automation.logDataReplication(dataItem)
    } else {
        fmt.Printf("Error replicating data item %s to all nodes.\n", dataItem.ID)
    }
}

// finalizeReplicationCycle finalizes the data replication check cycle and logs the result in the ledger
func (automation *DataReplicationAutomation) finalizeReplicationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeDataReplicationCycle()
    if success {
        fmt.Println("Data replication check cycle finalized successfully.")
        automation.logReplicationCycleFinalization()
    } else {
        fmt.Println("Error finalizing data replication check cycle.")
    }
}

// logDataReplication logs each data replication action into the ledger for traceability
func (automation *DataReplicationAutomation) logDataReplication(dataItem common.DataItem) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("data-replication-%s", dataItem.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Data Replication",
        Status:    "Replicated",
        Details:   fmt.Sprintf("Data item %s successfully replicated across all nodes.", dataItem.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with data replication event for data item %s.\n", dataItem.ID)
}

// logReplicationCycleFinalization logs the finalization of a data replication check cycle into the ledger
func (automation *DataReplicationAutomation) logReplicationCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("data-replication-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Data Replication Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with data replication cycle finalization.")
}

// AddEncryptionToDataItem encrypts the data item before replication
func (automation *DataReplicationAutomation) AddEncryptionToDataItem(dataItem common.DataItem) common.DataItem {
    encryptedData, err := encryption.EncryptData(dataItem)
    if err != nil {
        fmt.Println("Error encrypting data item:", err)
        return dataItem
    }
    dataItem.EncryptedData = encryptedData
    fmt.Println("Data item successfully encrypted.")
    return dataItem
}

// ensureDataReplicationIntegrity checks the integrity of replicated data and triggers enforcement if necessary
func (automation *DataReplicationAutomation) ensureDataReplicationIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateDataReplicationIntegrity()
    if !integrityValid {
        fmt.Println("Data replication integrity breach detected. Re-triggering data replication enforcement.")
        automation.monitorAndEnforceDataReplication()
    } else {
        fmt.Println("Data replication integrity is valid.")
    }
}
