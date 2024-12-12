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
    DataDistributionCheckInterval = 2500 * time.Millisecond // Interval for checking data distribution across nodes
    SubBlocksPerBlock             = 1000                    // Number of sub-blocks in a block
)

// DataDistributionEnforcementAutomation automates the enforcement of data distribution across all nodes for redundancy and high availability
type DataDistributionEnforcementAutomation struct {
    consensusSystem          *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance           *ledger.Ledger               // Ledger to store data distribution actions
    stateMutex               *sync.RWMutex                // Mutex for thread-safe access
    distributionCheckCount   int                          // Counter for distribution check cycles
}

// NewDataDistributionEnforcementAutomation initializes the automation for enforcing data distribution across nodes
func NewDataDistributionEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataDistributionEnforcementAutomation {
    return &DataDistributionEnforcementAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        distributionCheckCount: 0,
    }
}

// StartDataDistributionEnforcement starts the continuous loop for monitoring and enforcing data distribution across nodes
func (automation *DataDistributionEnforcementAutomation) StartDataDistributionEnforcement() {
    ticker := time.NewTicker(DataDistributionCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceDataDistribution()
        }
    }()
}

// monitorAndEnforceDataDistribution checks if data is distributed across all nodes and enforces proper distribution if needed
func (automation *DataDistributionEnforcementAutomation) monitorAndEnforceDataDistribution() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch all data that should be distributed across nodes
    dataItems := automation.consensusSystem.GetDataToDistribute()

    for _, dataItem := range dataItems {
        if automation.checkDataDistribution(dataItem) {
            fmt.Printf("Data item %s is correctly distributed across all nodes.\n", dataItem.ID)
        } else {
            fmt.Printf("Data item %s is not properly distributed. Enforcing distribution to remaining nodes.\n", dataItem.ID)
            automation.distributeDataToNodes(dataItem)
        }
    }

    automation.distributionCheckCount++
    fmt.Printf("Data distribution enforcement cycle #%d executed.\n", automation.distributionCheckCount)

    if automation.distributionCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeDistributionCycle()
    }
}

// checkDataDistribution checks if the data item is distributed across all nodes
func (automation *DataDistributionEnforcementAutomation) checkDataDistribution(dataItem common.DataItem) bool {
    return automation.consensusSystem.IsDataDistributedAcrossNodes(dataItem.ID)
}

// distributeDataToNodes enforces data distribution to nodes where the data is missing
func (automation *DataDistributionEnforcementAutomation) distributeDataToNodes(dataItem common.DataItem) {
    // Encrypt the data item before distributing
    encryptedDataItem := automation.AddEncryptionToDataItem(dataItem)

    // Distribute the data item to all nodes via the Synnergy Consensus
    distributionSuccess := automation.consensusSystem.EnforceDataDistributionAcrossNodes(encryptedDataItem)

    if distributionSuccess {
        fmt.Printf("Data item %s successfully distributed across all nodes.\n", dataItem.ID)
        automation.logDataDistribution(dataItem)
    } else {
        fmt.Printf("Error distributing data item %s to all nodes.\n", dataItem.ID)
    }
}

// finalizeDistributionCycle finalizes the data distribution check cycle and logs the result in the ledger
func (automation *DataDistributionEnforcementAutomation) finalizeDistributionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeDataDistributionCycle()
    if success {
        fmt.Println("Data distribution check cycle finalized successfully.")
        automation.logDistributionCycleFinalization()
    } else {
        fmt.Println("Error finalizing data distribution check cycle.")
    }
}

// logDataDistribution logs each data distribution action into the ledger for traceability
func (automation *DataDistributionEnforcementAutomation) logDataDistribution(dataItem common.DataItem) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("data-distribution-%s", dataItem.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Data Distribution",
        Status:    "Distributed",
        Details:   fmt.Sprintf("Data item %s successfully distributed across all nodes.", dataItem.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with data distribution event for data item %s.\n", dataItem.ID)
}

// logDistributionCycleFinalization logs the finalization of a data distribution check cycle into the ledger
func (automation *DataDistributionEnforcementAutomation) logDistributionCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("data-distribution-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Data Distribution Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with data distribution cycle finalization.")
}

// AddEncryptionToDataItem encrypts the data item before distribution
func (automation *DataDistributionEnforcementAutomation) AddEncryptionToDataItem(dataItem common.DataItem) common.DataItem {
    encryptedData, err := encryption.EncryptData(dataItem)
    if err != nil {
        fmt.Println("Error encrypting data item:", err)
        return dataItem
    }
    dataItem.EncryptedData = encryptedData
    fmt.Println("Data item successfully encrypted.")
    return dataItem
}

// ensureDataDistributionIntegrity checks the integrity of data distribution and triggers enforcement if necessary
func (automation *DataDistributionEnforcementAutomation) ensureDataDistributionIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateDataDistributionIntegrity()
    if !integrityValid {
        fmt.Println("Data distribution integrity breach detected. Re-triggering data distribution enforcement.")
        automation.monitorAndEnforceDataDistribution()
    } else {
        fmt.Println("Data distribution integrity is valid.")
    }
}
