package automations

import (
    "fmt"
    "time"
    "sync"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/common"
)

const (
    PartitionDataRetrievalInterval = 20 * time.Second  // Time interval for data retrieval
    MaxDataRetrievalRetries        = 3                 // Max retries for partition data retrieval
)

// PartitionDataRetrievalAutomation handles retrieving data from blockchain partitions.
type PartitionDataRetrievalAutomation struct {
    ledgerInstance        *ledger.Ledger               // Reference to the ledger
    consensusSystem       *consensus.SynnergyConsensus // Reference to the consensus mechanism
    stateMutex            *sync.RWMutex                // Mutex for concurrency control
    dataRetrievalAttempts map[string]int               // Tracks data retrieval retry attempts per partition
}

// NewPartitionDataRetrievalAutomation initializes the partition data retrieval automation.
func NewPartitionDataRetrievalAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *PartitionDataRetrievalAutomation {
    return &PartitionDataRetrievalAutomation{
        ledgerInstance:  ledgerInstance,
        consensusSystem: consensusSystem,
        stateMutex:      stateMutex,
        dataRetrievalAttempts: make(map[string]int),
    }
}

// StartPartitionDataRetrievalAutomation starts the automated partition data retrieval process in a continuous loop.
func (automation *PartitionDataRetrievalAutomation) StartPartitionDataRetrievalAutomation() {
    ticker := time.NewTicker(PartitionDataRetrievalInterval)
    go func() {
        for range ticker.C {
            automation.retrievePartitionData()
        }
    }()
}

// retrievePartitionData checks for data in partitions and retrieves it if necessary.
func (automation *PartitionDataRetrievalAutomation) retrievePartitionData() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch active partitions that may have data available
    partitions, err := automation.consensusSystem.GetActivePartitions()
    if err != nil {
        fmt.Printf("Error fetching active partitions: %v\n", err)
        return
    }

    for _, partition := range partitions {
        err := automation.handlePartitionData(partition)
        if err != nil {
            fmt.Printf("Error handling data retrieval for partition %s: %v\n", partition.ID, err)
            automation.retryDataRetrieval(partition.ID)
        }
    }
}

// handlePartitionData retrieves and decrypts the data from a specific partition.
func (automation *PartitionDataRetrievalAutomation) handlePartitionData(partition common.Partition) error {
    // Fetch data from the partition
    encryptedData, err := automation.consensusSystem.RetrievePartitionData(partition.ID)
    if err != nil {
        return fmt.Errorf("error retrieving data from partition %s: %v", partition.ID, err)
    }

    // Decrypt the partition data
    decryptedData, err := automation.decryptPartitionData(encryptedData)
    if err != nil {
        return fmt.Errorf("error decrypting data from partition %s: %v", partition.ID, err)
    }

    fmt.Printf("Successfully retrieved and decrypted data from partition %s.\n", partition.ID)

    // Log the partition data retrieval in the ledger
    automation.logDataRetrieval(partition.ID, decryptedData)

    return nil
}

// retryDataRetrieval retries partition data retrieval if the first attempt fails.
func (automation *PartitionDataRetrievalAutomation) retryDataRetrieval(partitionID string) {
    automation.dataRetrievalAttempts[partitionID]++
    if automation.dataRetrievalAttempts[partitionID] < MaxDataRetrievalRetries {
        // Retry data retrieval for the partition
        fmt.Printf("Retrying data retrieval for partition %s...\n", partitionID)
        err := automation.handlePartitionData(common.Partition{ID: partitionID})
        if err != nil {
            fmt.Printf("Retry failed for partition %s: %v\n", partitionID, err)
        }
    } else {
        fmt.Printf("Max retries reached for partition %s. Logging as failure.\n", partitionID)
        automation.logDataRetrieval(partitionID, "Failed")
    }
}

// logDataRetrieval logs the data retrieval event into the ledger.
func (automation *PartitionDataRetrievalAutomation) logDataRetrieval(partitionID string, data interface{}) {
    entry := common.LedgerEntry{
        ID:        partitionID,
        Timestamp: time.Now().Unix(),
        Type:      "Partition Data Retrieval",
        Status:    "Success",
        Details:   fmt.Sprintf("Data retrieved for partition %s: %v", partitionID, data),
    }

    if err := automation.ledgerInstance.AddEntry(entry); err != nil {
        fmt.Printf("Failed to log data retrieval event for partition %s: %v\n", partitionID, err)
    }
}

// decryptPartitionData decrypts the partition data after retrieval.
func (automation *PartitionDataRetrievalAutomation) decryptPartitionData(encryptedData string) (string, error) {
    decryptedData, err := encryption.DecryptData([]byte(encryptedData))
    if err != nil {
        return "", fmt.Errorf("error decrypting partition data: %v", err)
    }
    return string(decryptedData), nil
}
