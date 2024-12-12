package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/common"
)

const (
    PartitionRebalanceCheckInterval = 30 * time.Second  // Interval to check for partition rebalancing
    MaxRebalancingRetries           = 3                 // Maximum number of retries for partition rebalancing
)

// PartitionRebalancingAutomation manages the rebalancing of blockchain partitions based on system load and triggers.
type PartitionRebalancingAutomation struct {
    ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging
    consensusSystem *consensus.SynnergyConsensus // Reference to the Synnergy Consensus
    stateMutex      *sync.RWMutex                // Mutex for concurrency control
    rebalanceAttempts map[string]int             // Tracks rebalancing retries for each partition
}

// NewPartitionRebalancingAutomation creates and initializes a new PartitionRebalancingAutomation.
func NewPartitionRebalancingAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *PartitionRebalancingAutomation {
    return &PartitionRebalancingAutomation{
        ledgerInstance:    ledgerInstance,
        consensusSystem:   consensusSystem,
        stateMutex:        stateMutex,
        rebalanceAttempts: make(map[string]int),
    }
}

// StartPartitionRebalancingAutomation starts the partition rebalancing automation in a continuous loop.
func (automation *PartitionRebalancingAutomation) StartPartitionRebalancingAutomation() {
    ticker := time.NewTicker(PartitionRebalanceCheckInterval)
    go func() {
        for range ticker.C {
            automation.rebalancePartitions()
        }
    }()
}

// rebalancePartitions checks for partitions that need rebalancing and performs rebalancing operations.
func (automation *PartitionRebalancingAutomation) rebalancePartitions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    partitions, err := automation.consensusSystem.GetOverloadedPartitions()
    if err != nil {
        fmt.Printf("Error fetching overloaded partitions: %v\n", err)
        return
    }

    for _, partition := range partitions {
        err := automation.handlePartitionRebalance(partition)
        if err != nil {
            fmt.Printf("Error rebalancing partition %s: %v\n", partition.ID, err)
            automation.retryPartitionRebalance(partition.ID)
        }
    }
}

// handlePartitionRebalance performs rebalancing operations on a specific partition.
func (automation *PartitionRebalancingAutomation) handlePartitionRebalance(partition common.Partition) error {
    // Fetch the current state of the partition
    partitionState, err := automation.consensusSystem.GetPartitionState(partition.ID)
    if err != nil {
        return fmt.Errorf("error retrieving partition state for %s: %v", partition.ID, err)
    }

    // Encrypt partition state before rebalancing
    encryptedState, err := automation.encryptPartitionState(partitionState)
    if err != nil {
        return fmt.Errorf("error encrypting partition state for %s: %v", partition.ID, err)
    }

    // Perform rebalancing operation
    err = automation.consensusSystem.RebalancePartition(partition.ID, encryptedState)
    if err != nil {
        return fmt.Errorf("error during rebalancing of partition %s: %v", partition.ID, err)
    }

    fmt.Printf("Successfully rebalanced partition %s.\n", partition.ID)

    // Log the rebalancing event in the ledger
    automation.logPartitionRebalance(partition.ID, "Rebalanced")

    return nil
}

// retryPartitionRebalance retries the rebalancing process if it fails the first time.
func (automation *PartitionRebalancingAutomation) retryPartitionRebalance(partitionID string) {
    automation.rebalanceAttempts[partitionID]++
    if automation.rebalanceAttempts[partitionID] < MaxRebalancingRetries {
        fmt.Printf("Retrying rebalancing for partition %s...\n", partitionID)
        err := automation.handlePartitionRebalance(common.Partition{ID: partitionID})
        if err != nil {
            fmt.Printf("Retry failed for partition %s: %v\n", partitionID, err)
        }
    } else {
        fmt.Printf("Max retries reached for partition %s. Logging failure.\n", partitionID)
        automation.logPartitionRebalance(partitionID, "Failed")
    }
}

// logPartitionRebalance logs the partition rebalancing event into the ledger.
func (automation *PartitionRebalancingAutomation) logPartitionRebalance(partitionID string, status string) {
    entry := common.LedgerEntry{
        ID:        partitionID,
        Timestamp: time.Now().Unix(),
        Type:      "Partition Rebalancing",
        Status:    status,
        Details:   fmt.Sprintf("Partition %s rebalancing status: %s", partitionID, status),
    }

    if err := automation.ledgerInstance.AddEntry(entry); err != nil {
        fmt.Printf("Failed to log partition rebalancing event for partition %s: %v\n", partitionID, err)
    }
}

// encryptPartitionState encrypts the state of a partition before rebalancing.
func (automation *PartitionRebalancingAutomation) encryptPartitionState(state string) (string, error) {
    encryptedState, err := encryption.EncryptData([]byte(state))
    if err != nil {
        return "", fmt.Errorf("error encrypting partition state: %v", err)
    }
    return string(encryptedState), nil
}

