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
    PartitionCreationInterval    = 30 * time.Second  // Time interval for checking if partition creation is necessary
    MaxPartitionRetries          = 3                 // Max retries for partition creation
)

// PartitionCreationAutomation manages dynamic partition creation to ensure load balancing and scalability.
type PartitionCreationAutomation struct {
    ledgerInstance        *ledger.Ledger               // Reference to the ledger for logging partition actions
    consensusSystem       *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    stateMutex            *sync.RWMutex                // Mutex for thread safety
    partitionAttempts     map[string]int               // Tracks retry attempts for partitions
}

// NewPartitionCreationAutomation initializes the partition creation automation.
func NewPartitionCreationAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *PartitionCreationAutomation {
    return &PartitionCreationAutomation{
        ledgerInstance:  ledgerInstance,
        consensusSystem: consensusSystem,
        stateMutex:      stateMutex,
        partitionAttempts: make(map[string]int),
    }
}

// StartPartitionCreationAutomation continuously checks for the need for partition creation and enforces it.
func (automation *PartitionCreationAutomation) StartPartitionCreationAutomation() {
    ticker := time.NewTicker(PartitionCreationInterval)
    go func() {
        for range ticker.C {
            automation.enforcePartitionCreation()
        }
    }()
}

// enforcePartitionCreation checks if new partitions are needed and performs creation if necessary.
func (automation *PartitionCreationAutomation) enforcePartitionCreation() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Retrieve the current load distribution from the consensus
    currentLoad, err := automation.consensusSystem.GetLoadDistribution()
    if err != nil {
        fmt.Printf("Error retrieving load distribution: %v\n", err)
        return
    }

    // Identify if a new partition is needed for load balancing
    if automation.isPartitionNeeded(currentLoad) {
        partitionID, err := automation.createPartition()
        if err != nil {
            fmt.Printf("Partition creation failed: %v\n", err)
            automation.retryPartitionCreation(partitionID)
        } else {
            fmt.Printf("Partition %s successfully created.\n", partitionID)
        }
    }
}

// isPartitionNeeded checks whether the system load requires partition creation.
func (automation *PartitionCreationAutomation) isPartitionNeeded(load common.LoadDistribution) bool {
    // Logic to determine if the current load exceeds threshold, requiring a new partition.
    return load.AverageLoad > common.PartitionCreationLoadThreshold
}

// createPartition performs the partition creation and logs the event.
func (automation *PartitionCreationAutomation) createPartition() (string, error) {
    partitionID := automation.generatePartitionID()

    encryptedPartitionData, err := automation.encryptPartitionData(partitionID)
    if err != nil {
        return "", fmt.Errorf("error encrypting partition data: %v", err)
    }

    // Perform partition creation in the consensus system
    err = automation.consensusSystem.CreatePartition(partitionID, encryptedPartitionData)
    if err != nil {
        return "", fmt.Errorf("error creating partition: %v", err)
    }

    // Log the partition creation event into the ledger
    automation.logPartitionCreation(partitionID, "Created")
    return partitionID, nil
}

// retryPartitionCreation retries partition creation if it fails.
func (automation *PartitionCreationAutomation) retryPartitionCreation(partitionID string) {
    automation.partitionAttempts[partitionID]++
    if automation.partitionAttempts[partitionID] < MaxPartitionRetries {
        _, err := automation.createPartition()
        if err != nil {
            fmt.Printf("Retry failed for partition creation %s: %v\n", partitionID, err)
        }
    } else {
        fmt.Printf("Max retries reached for partition creation %s. Logging as failure.\n", partitionID)
        automation.logPartitionCreation(partitionID, "Failed")
    }
}

// logPartitionCreation logs the partition creation event into the ledger.
func (automation *PartitionCreationAutomation) logPartitionCreation(partitionID, status string) {
    entry := common.LedgerEntry{
        ID:        partitionID,
        Timestamp: time.Now().Unix(),
        Type:      "Partition Creation",
        Status:    status,
        Details:   fmt.Sprintf("Partition %s has status: %s", partitionID, status),
    }

    if err := automation.ledgerInstance.AddEntry(entry); err != nil {
        fmt.Printf("Failed to log partition creation event for %s: %v\n", partitionID, err)
    }
}

// encryptPartitionData encrypts the partition data before storing or processing.
func (automation *PartitionCreationAutomation) encryptPartitionData(partitionID string) (string, error) {
    data := fmt.Sprintf("PartitionData-%s", partitionID)
    encryptedData, err := encryption.EncryptData([]byte(data))
    if err != nil {
        return "", fmt.Errorf("error encrypting partition data: %v", err)
    }
    return string(encryptedData), nil
}

// generatePartitionID generates a new unique partition ID.
func (automation *PartitionCreationAutomation) generatePartitionID() string {
    return fmt.Sprintf("partition-%d", time.Now().UnixNano())
}
