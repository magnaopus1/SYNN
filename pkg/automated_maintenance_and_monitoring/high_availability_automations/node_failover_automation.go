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
    NodeFailoverCheckInterval = 2500 * time.Millisecond // Interval for checking node health and failover status
    SubBlocksPerBlock         = 1000                    // Number of sub-blocks in a block
)

// NodeFailoverAutomation automates the failover process across nodes to ensure high availability in case of failure
type NodeFailoverAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger to store failover actions
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    failoverCheckCount    int                          // Counter for failover check cycles
    failedNodeChecks      map[string]int               // Map to track failed nodes and their status
    maxFailoverAttempts   int                          // Maximum number of failover attempts before marking a node as failed
}

// NewNodeFailoverAutomation initializes the automation for failover across nodes
func NewNodeFailoverAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex, maxFailoverAttempts int) *NodeFailoverAutomation {
    return &NodeFailoverAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        failoverCheckCount:  0,
        failedNodeChecks:    make(map[string]int),
        maxFailoverAttempts: maxFailoverAttempts,
    }
}

// StartNodeFailoverCheck starts the continuous loop for monitoring node health and triggering failover actions if necessary
func (automation *NodeFailoverAutomation) StartNodeFailoverCheck() {
    ticker := time.NewTicker(NodeFailoverCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndTriggerFailover()
        }
    }()
}

// monitorAndTriggerFailover checks the health of all nodes in the network and triggers failover actions for failing nodes
func (automation *NodeFailoverAutomation) monitorAndTriggerFailover() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch node health status from the Synnergy Consensus
    nodeHealth := automation.consensusSystem.GetNodeHealthStatus()

    for _, node := range nodeHealth {
        if automation.isNodeHealthy(node) {
            fmt.Printf("Node %s passed health check.\n", node.ID)
            automation.failedNodeChecks[node.ID] = 0 // Reset failure count on success
        } else {
            automation.failedNodeChecks[node.ID]++
            fmt.Printf("Node %s failed health check. Consecutive failures: %d\n", node.ID, automation.failedNodeChecks[node.ID])

            if automation.failedNodeChecks[node.ID] >= automation.maxFailoverAttempts {
                fmt.Printf("Node %s exceeded maximum failover attempts. Marking node as failed.\n", node.ID)
                automation.markNodeAsFailed(node)
            } else {
                fmt.Printf("Triggering failover for node %s.\n", node.ID)
                automation.triggerFailoverForNode(node)
            }
        }
    }

    automation.failoverCheckCount++
    fmt.Printf("Failover check cycle #%d executed.\n", automation.failoverCheckCount)

    if automation.failoverCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeFailoverCycle()
    }
}

// isNodeHealthy checks if a node is operating within acceptable limits and passes health checks
func (automation *NodeFailoverAutomation) isNodeHealthy(node common.Node) bool {
    return automation.consensusSystem.IsNodeHealthy(node.ID)
}

// triggerFailoverForNode triggers the failover process for a node that has failed a health check but has not yet exceeded failover attempts
func (automation *NodeFailoverAutomation) triggerFailoverForNode(node common.Node) {
    // Encrypt node data before triggering failover
    encryptedNodeData := automation.AddEncryptionToNodeData(node)

    // Trigger failover through the Synnergy Consensus
    failoverSuccess := automation.consensusSystem.TriggerFailoverForNode(encryptedNodeData)

    if failoverSuccess {
        fmt.Printf("Failover successfully triggered for node %s.\n", node.ID)
        automation.logFailoverEvent(node)
    } else {
        fmt.Printf("Error triggering failover for node %s.\n", node.ID)
    }
}

// markNodeAsFailed marks a node as permanently failed if it exceeds the maximum failover attempts
func (automation *NodeFailoverAutomation) markNodeAsFailed(node common.Node) {
    // Log the node failure event
    automation.logNodeFailureEvent(node)

    // Handle any actions for removing or isolating the failed node from the network
    removalSuccess := automation.consensusSystem.RemoveFailedNode(node)

    if removalSuccess {
        fmt.Printf("Node %s has been successfully removed from the network.\n", node.ID)
    } else {
        fmt.Printf("Error removing node %s from the network.\n", node.ID)
    }
}

// finalizeFailoverCycle finalizes the failover check cycle and logs the result in the ledger
func (automation *NodeFailoverAutomation) finalizeFailoverCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeFailoverCycle()
    if success {
        fmt.Println("Failover check cycle finalized successfully.")
        automation.logFailoverCycleFinalization()
    } else {
        fmt.Println("Error finalizing failover check cycle.")
    }
}

// logFailoverEvent logs the failover event for a specific node into the ledger for traceability
func (automation *NodeFailoverAutomation) logFailoverEvent(node common.Node) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("node-failover-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Node Failover",
        Status:    "Triggered",
        Details:   fmt.Sprintf("Failover successfully triggered for node %s.", node.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with failover event for node %s.\n", node.ID)
}

// logNodeFailureEvent logs the permanent failure of a node that exceeded the maximum failover attempts
func (automation *NodeFailoverAutomation) logNodeFailureEvent(node common.Node) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("node-failure-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Node Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Node %s has been marked as failed and removed from the network.", node.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with node failure event for node %s.\n", node.ID)
}

// logFailoverCycleFinalization logs the finalization of a failover check cycle into the ledger
func (automation *NodeFailoverAutomation) logFailoverCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("failover-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Failover Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with failover cycle finalization.")
}

// AddEncryptionToNodeData encrypts the node data before triggering failover or removal
func (automation *NodeFailoverAutomation) AddEncryptionToNodeData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node data successfully encrypted.")
    return node
}

// ensureFailoverIntegrity checks the integrity of failover processes and data, and triggers actions if necessary
func (automation *NodeFailoverAutomation) ensureFailoverIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateFailoverIntegrity()
    if !integrityValid {
        fmt.Println("Failover integrity breach detected. Re-triggering failover checks.")
        automation.monitorAndTriggerFailover()
    } else {
        fmt.Println("Failover integrity is valid.")
    }
}
