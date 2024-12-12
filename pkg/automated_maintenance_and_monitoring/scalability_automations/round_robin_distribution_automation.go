package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
)

const (
    DistributionInterval = 30 * time.Second  // Interval to trigger round-robin distribution
    MaxRetries           = 3                 // Maximum retries in case of distribution failure
)

// RoundRobinDistributionAutomation handles distributing transactions and sub-blocks across network nodes using round-robin.
type RoundRobinDistributionAutomation struct {
    ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging
    consensusSystem *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    stateMutex      *sync.RWMutex                // Mutex for concurrency control
    nodeList        []common.Node                // List of active nodes in the network
    currentIndex    int                          // Current index for round-robin distribution
}

// NewRoundRobinDistributionAutomation creates and initializes a new RoundRobinDistributionAutomation.
func NewRoundRobinDistributionAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, nodeList []common.Node, stateMutex *sync.RWMutex) *RoundRobinDistributionAutomation {
    return &RoundRobinDistributionAutomation{
        ledgerInstance:  ledgerInstance,
        consensusSystem: consensusSystem,
        nodeList:        nodeList,
        stateMutex:      stateMutex,
        currentIndex:    0,
    }
}

// StartDistributionAutomation starts the round-robin distribution in a continuous loop.
func (automation *RoundRobinDistributionAutomation) StartDistributionAutomation() {
    ticker := time.NewTicker(DistributionInterval)
    go func() {
        for range ticker.C {
            automation.distributeLoad()
        }
    }()
}

// distributeLoad performs the round-robin distribution of transactions and sub-blocks across network nodes.
func (automation *RoundRobinDistributionAutomation) distributeLoad() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    subBlockID := automation.generateSubBlockID()

    fmt.Printf("Distributing sub-block with ID: %s\n", subBlockID)

    success := automation.performRoundRobinDistribution(subBlockID)
    if !success {
        automation.retryDistribution(subBlockID)
        return
    }

    // Encrypt distribution details before logging
    encryptedDistributionDetails, err := automation.encryptDistributionData(subBlockID)
    if err != nil {
        fmt.Printf("Error encrypting distribution data: %v\n", err)
        return
    }

    // Log the distribution result in the ledger
    automation.logDistributionResult(subBlockID, "Success", encryptedDistributionDetails)

    fmt.Printf("Distribution completed successfully for sub-block ID: %s\n", subBlockID)
}

// performRoundRobinDistribution distributes the sub-block to the next node in the round-robin cycle.
func (automation *RoundRobinDistributionAutomation) performRoundRobinDistribution(subBlockID string) bool {
    if len(automation.nodeList) == 0 {
        fmt.Printf("No available nodes for distribution.\n")
        return false
    }

    // Select the next node in the round-robin sequence
    selectedNode := automation.nodeList[automation.currentIndex]

    fmt.Printf("Distributing sub-block %s to node %s\n", subBlockID, selectedNode.ID)

    err := automation.consensusSystem.ValidateAndDistributeSubBlock(subBlockID, selectedNode)
    if err != nil {
        fmt.Printf("Failed to distribute sub-block %s to node %s: %v\n", subBlockID, selectedNode.ID, err)
        return false
    }

    // Move to the next node in the round-robin sequence
    automation.currentIndex = (automation.currentIndex + 1) % len(automation.nodeList)

    return true
}

// retryDistribution retries the distribution if it fails, up to MaxRetries.
func (automation *RoundRobinDistributionAutomation) retryDistribution(subBlockID string) {
    retryAttempts := 0
    for retryAttempts < MaxRetries {
        fmt.Printf("Retrying distribution for sub-block %s (attempt %d)\n", subBlockID, retryAttempts+1)
        if automation.performRoundRobinDistribution(subBlockID) {
            encryptedData, err := automation.encryptDistributionData(subBlockID)
            if err != nil {
                fmt.Printf("Error encrypting distribution data during retry: %v\n", err)
                return
            }
            automation.logDistributionResult(subBlockID, "Success (after retry)", encryptedData)
            fmt.Printf("Distribution completed successfully for sub-block ID %s on retry.\n", subBlockID)
            return
        }
        retryAttempts++
    }
    fmt.Printf("Max retries reached for sub-block %s. Logging failure.\n", subBlockID)
    automation.logDistributionResult(subBlockID, "Failed", "")
}

// logDistributionResult logs the result of the sub-block distribution into the ledger.
func (automation *RoundRobinDistributionAutomation) logDistributionResult(subBlockID string, status string, encryptedDetails string) {
    entry := common.LedgerEntry{
        ID:        subBlockID,
        Timestamp: time.Now().Unix(),
        Type:      "SubBlock Distribution",
        Status:    status,
        Details:   fmt.Sprintf("Sub-block distribution %s status: %s, Details: %s", subBlockID, status, encryptedDetails),
    }

    if err := automation.ledgerInstance.AddEntry(entry); err != nil {
        fmt.Printf("Failed to log sub-block distribution result for %s: %v\n", subBlockID, err)
    }
}

// encryptDistributionData encrypts the sub-block distribution data before logging to the ledger.
func (automation *RoundRobinDistributionAutomation) encryptDistributionData(subBlockID string) (string, error) {
    encryptedData, err := encryption.EncryptData([]byte(fmt.Sprintf("SubBlock ID: %s", subBlockID)))
    if err != nil {
        return "", fmt.Errorf("error encrypting distribution data: %v", err)
    }
    return string(encryptedData), nil
}

// generateSubBlockID generates a unique identifier for each sub-block to be distributed.
func (automation *RoundRobinDistributionAutomation) generateSubBlockID() string {
    return fmt.Sprintf("SUBBLOCK-%d", time.Now().UnixNano())
}
