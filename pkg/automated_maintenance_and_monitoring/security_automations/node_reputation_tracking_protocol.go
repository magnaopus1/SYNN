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
    ReputationMonitoringInterval  = 20 * time.Second // Interval for checking node reputation
    MaxReputationRetries          = 3                // Maximum retries for adjusting reputation
    SubBlocksPerBlock             = 1000             // Number of sub-blocks in a block
    ReputationThresholdWarning    = 60               // Threshold for issuing a warning to low-reputation nodes
    ReputationThresholdSuspension = 40               // Threshold for suspending a low-reputation node
    ReputationThresholdBan        = 20               // Threshold for banning a node based on reputation
)

// NodeReputationTrackingProtocol handles the monitoring and enforcement of node reputation policies
type NodeReputationTrackingProtocol struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger for logging reputation-related events
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    reputationRetryCount   map[string]int               // Counter for retrying reputation adjustments
    reputationCycleCount   int                          // Counter for reputation monitoring cycles
    nodeReputationRecords  map[string]int               // Tracks the reputation score of each node
}

// NewNodeReputationTrackingProtocol initializes the automation for node reputation tracking
func NewNodeReputationTrackingProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *NodeReputationTrackingProtocol {
    return &NodeReputationTrackingProtocol{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        reputationRetryCount:  make(map[string]int),
        nodeReputationRecords: make(map[string]int),
        reputationCycleCount:  0,
    }
}

// StartReputationMonitoring starts the continuous loop for monitoring node reputations
func (protocol *NodeReputationTrackingProtocol) StartReputationMonitoring() {
    ticker := time.NewTicker(ReputationMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorNodeReputation()
        }
    }()
}

// monitorNodeReputation checks node reputations and applies warnings, suspensions, or bans based on threshold conditions
func (protocol *NodeReputationTrackingProtocol) monitorNodeReputation() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch reputation scores for active nodes from the consensus system
    activeNodes := protocol.consensusSystem.FetchNodeReputations()

    for _, node := range activeNodes {
        reputation := node.ReputationScore
        protocol.nodeReputationRecords[node.ID] = reputation

        if reputation <= ReputationThresholdBan {
            fmt.Printf("Node %s has a reputation of %d. Triggering ban.\n", node.ID, reputation)
            protocol.banNode(node)
        } else if reputation <= ReputationThresholdSuspension {
            fmt.Printf("Node %s has a reputation of %d. Triggering suspension.\n", node.ID, reputation)
            protocol.suspendNode(node)
        } else if reputation <= ReputationThresholdWarning {
            fmt.Printf("Node %s has a reputation of %d. Issuing warning.\n", node.ID, reputation)
            protocol.warnNode(node)
        } else {
            fmt.Printf("Node %s has a stable reputation of %d. No action needed.\n", node.ID, reputation)
        }
    }

    protocol.reputationCycleCount++
    fmt.Printf("Reputation monitoring cycle #%d completed.\n", protocol.reputationCycleCount)

    if protocol.reputationCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeReputationCycle()
    }
}

// warnNode issues a warning to a node with a reputation below the warning threshold
func (protocol *NodeReputationTrackingProtocol) warnNode(node common.Node) {
    encryptedReputationData := protocol.encryptReputationData(node)

    // Issue a warning through the Synnergy Consensus system
    warningSuccess := protocol.consensusSystem.WarnNode(node, encryptedReputationData)

    if warningSuccess {
        fmt.Printf("Warning issued to node %s.\n", node.ID)
        protocol.logReputationEvent(node, "Warning Issued")
        protocol.resetReputationRetry(node.ID)
    } else {
        fmt.Printf("Error issuing warning to node %s. Retrying...\n", node.ID)
        protocol.retryReputationAdjustment(node)
    }
}

// suspendNode suspends a node with a reputation below the suspension threshold
func (protocol *NodeReputationTrackingProtocol) suspendNode(node common.Node) {
    encryptedReputationData := protocol.encryptReputationData(node)

    // Attempt to suspend the node through the Synnergy Consensus system
    suspensionSuccess := protocol.consensusSystem.SuspendNode(node, encryptedReputationData)

    if suspensionSuccess {
        fmt.Printf("Node %s suspended.\n", node.ID)
        protocol.logReputationEvent(node, "Suspended")
        protocol.resetReputationRetry(node.ID)
    } else {
        fmt.Printf("Error suspending node %s. Retrying...\n", node.ID)
        protocol.retryReputationAdjustment(node)
    }
}

// banNode bans a node with a reputation below the banning threshold
func (protocol *NodeReputationTrackingProtocol) banNode(node common.Node) {
    encryptedReputationData := protocol.encryptReputationData(node)

    // Attempt to ban the node through the Synnergy Consensus system
    banSuccess := protocol.consensusSystem.BanNode(node, encryptedReputationData)

    if banSuccess {
        fmt.Printf("Node %s banned.\n", node.ID)
        protocol.logReputationEvent(node, "Banned")
        protocol.resetReputationRetry(node.ID)
    } else {
        fmt.Printf("Error banning node %s. Retrying...\n", node.ID)
        protocol.retryReputationAdjustment(node)
    }
}

// retryReputationAdjustment retries a failed reputation adjustment action
func (protocol *NodeReputationTrackingProtocol) retryReputationAdjustment(node common.Node) {
    protocol.reputationRetryCount[node.ID]++
    if protocol.reputationRetryCount[node.ID] < MaxReputationRetries {
        protocol.adjustReputation(node)
    } else {
        fmt.Printf("Max retries reached for adjusting reputation of node %s. Adjustment failed.\n", node.ID)
        protocol.logReputationFailure(node)
    }
}

// adjustReputation is a helper function that determines the appropriate action to retry based on the node's reputation
func (protocol *NodeReputationTrackingProtocol) adjustReputation(node common.Node) {
    if node.ReputationScore <= ReputationThresholdBan {
        protocol.banNode(node)
    } else if node.ReputationScore <= ReputationThresholdSuspension {
        protocol.suspendNode(node)
    } else if node.ReputationScore <= ReputationThresholdWarning {
        protocol.warnNode(node)
    }
}

// resetReputationRetry resets the retry count for a node's reputation adjustment
func (protocol *NodeReputationTrackingProtocol) resetReputationRetry(nodeID string) {
    protocol.reputationRetryCount[nodeID] = 0
}

// finalizeReputationCycle finalizes the reputation monitoring cycle and logs the result in the ledger
func (protocol *NodeReputationTrackingProtocol) finalizeReputationCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeReputationCycle()
    if success {
        fmt.Println("Reputation monitoring cycle finalized successfully.")
        protocol.logReputationCycleFinalization()
    } else {
        fmt.Println("Error finalizing reputation monitoring cycle.")
    }
}

// logReputationEvent logs a reputation-related event into the ledger
func (protocol *NodeReputationTrackingProtocol) logReputationEvent(node common.Node, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("reputation-event-%s-%s", node.ID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Reputation Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s was %s based on reputation score.", node.ID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with reputation event for node %s.\n", node.ID)
}

// logReputationFailure logs a failure to adjust reputation into the ledger
func (protocol *NodeReputationTrackingProtocol) logReputationFailure(node common.Node) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("reputation-adjustment-failure-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Reputation Adjustment Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to adjust reputation for node %s after maximum retries.", node.ID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with reputation adjustment failure for node %s.\n", node.ID)
}

// logReputationCycleFinalization logs the finalization of a reputation monitoring cycle into the ledger
func (protocol *NodeReputationTrackingProtocol) logReputationCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("reputation-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Reputation Cycle Finalization",
        Status:    "Finalized",
        Details:   "Reputation monitoring cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with reputation monitoring cycle finalization.")
}

// encryptReputationData encrypts the reputation data before applying warnings, suspensions, or bans
func (protocol *NodeReputationTrackingProtocol) encryptReputationData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node.ReputationData)
    if err != nil {
        fmt.Println("Error encrypting reputation data:", err)
        return node
    }

    node.EncryptedData = encryptedData
    fmt.Println("Reputation data successfully encrypted for node ID:", node.ID)
    return node
}

// triggerEmergencyReputationAction triggers an emergency reputation-related action for a node in case of critical behavior
func (protocol *NodeReputationTrackingProtocol) triggerEmergencyReputationAction(nodeID string) {
    fmt.Printf("Emergency action triggered for node ID: %s.\n", nodeID)
    node := protocol.consensusSystem.GetNodeByID(nodeID)
    encryptedData := protocol.encryptReputationData(node)

    success := protocol.consensusSystem.TriggerEmergencyNodeBan(nodeID, encryptedData)

    if success {
        protocol.logReputationEvent(node, "Emergency Banned")
        fmt.Println("Emergency reputation action executed successfully.")
    } else {
        fmt.Println("Emergency reputation action failed.")
    }
}
