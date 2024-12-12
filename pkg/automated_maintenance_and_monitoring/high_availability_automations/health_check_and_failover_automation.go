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
    HealthCheckInterval      = 3000 * time.Millisecond // Interval for checking node health
    SubBlocksPerBlock        = 1000                    // Number of sub-blocks in a block
    MaxHealthCheckFailures   = 3                       // Maximum consecutive health check failures before triggering failover
)

// HealthCheckAndFailoverAutomation automates the health check and failover process across nodes for high availability
type HealthCheckAndFailoverAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger to store health check and failover actions
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    healthCheckCount      int                          // Counter for health check cycles
    failedNodeChecks      map[string]int               // Map to track consecutive failed health checks per node
}

// NewHealthCheckAndFailoverAutomation initializes the automation for health check and failover across nodes
func NewHealthCheckAndFailoverAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *HealthCheckAndFailoverAutomation {
    return &HealthCheckAndFailoverAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        healthCheckCount: 0,
        failedNodeChecks: make(map[string]int),
    }
}

// StartHealthCheck starts the continuous loop for monitoring node health and triggering failover if necessary
func (automation *HealthCheckAndFailoverAutomation) StartHealthCheck() {
    ticker := time.NewTicker(HealthCheckInterval)

    go func() {
        for range ticker.C {
            automation.performHealthCheck()
        }
    }()
}

// performHealthCheck checks the health of all nodes in the network and triggers failover for failing nodes
func (automation *HealthCheckAndFailoverAutomation) performHealthCheck() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch node health status from the Synnergy Consensus
    nodeHealth := automation.consensusSystem.GetNodeHealthStatus()

    for _, node := range nodeHealth {
        if automation.isNodeHealthy(node) {
            fmt.Printf("Node %s passed health check.\n", node.ID)
            automation.failedNodeChecks[node.ID] = 0 // Reset the failure count on a successful check
        } else {
            automation.failedNodeChecks[node.ID]++
            fmt.Printf("Node %s failed health check. Consecutive failures: %d\n", node.ID, automation.failedNodeChecks[node.ID])

            if automation.failedNodeChecks[node.ID] >= MaxHealthCheckFailures {
                fmt.Printf("Node %s has exceeded maximum failures. Triggering failover.\n", node.ID)
                automation.triggerFailoverForNode(node)
                automation.failedNodeChecks[node.ID] = 0 // Reset failure count after failover
            }
        }
    }

    automation.healthCheckCount++
    fmt.Printf("Health check cycle #%d executed.\n", automation.healthCheckCount)

    if automation.healthCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeHealthCheckCycle()
    }
}

// isNodeHealthy checks whether a node passes the health check based on various parameters like response time, load, etc.
func (automation *HealthCheckAndFailoverAutomation) isNodeHealthy(node common.Node) bool {
    return automation.consensusSystem.IsNodeHealthy(node.ID)
}

// triggerFailoverForNode triggers the failover process for a node that has failed the health checks
func (automation *HealthCheckAndFailoverAutomation) triggerFailoverForNode(node common.Node) {
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

// finalizeHealthCheckCycle finalizes the health check cycle and logs the result in the ledger
func (automation *HealthCheckAndFailoverAutomation) finalizeHealthCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeHealthCheckCycle()
    if success {
        fmt.Println("Health check cycle finalized successfully.")
        automation.logHealthCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing health check cycle.")
    }
}

// logFailoverEvent logs the failover event for a specific node into the ledger for traceability
func (automation *HealthCheckAndFailoverAutomation) logFailoverEvent(node common.Node) {
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

// logHealthCheckCycleFinalization logs the finalization of a health check cycle into the ledger
func (automation *HealthCheckAndFailoverAutomation) logHealthCheckCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("health-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Health Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with health check cycle finalization.")
}

// AddEncryptionToNodeData encrypts the node data before triggering failover
func (automation *HealthCheckAndFailoverAutomation) AddEncryptionToNodeData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node data successfully encrypted.")
    return node
}

// ensureHealthCheckIntegrity checks the integrity of health check and failover data and triggers actions if necessary
func (automation *HealthCheckAndFailoverAutomation) ensureHealthCheckIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateHealthCheckIntegrity()
    if !integrityValid {
        fmt.Println("Health check integrity breach detected. Re-triggering health checks.")
        automation.performHealthCheck()
    } else {
        fmt.Println("Health check integrity is valid.")
    }
}
