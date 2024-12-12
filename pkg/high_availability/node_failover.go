package high_availability

import (
    "fmt"
    "time"
    "synnergy_network/pkg/ledger"
)


// NewNodeFailoverManager initializes the NodeFailoverManager with primary and backup nodes
func NewNodeFailoverManager(primaryNodes []string, backupNodes []string, ledger *ledger.Ledger) *NodeFailoverManager {
    return &NodeFailoverManager{
        PrimaryNodes:     primaryNodes,
        BackupNodes:      backupNodes,
        NodeHealthStatus: make(map[string]bool),
        CurrentPrimary:   primaryNodes[0],  // Start with the first primary node
        LedgerInstance:   ledger,
    }
}

// StartMonitoring continuously monitors the health of primary nodes and triggers failover if needed
func (fm *NodeFailoverManager) StartMonitoring() {
    go func() {
        for {
            fm.mutex.Lock()
            currentPrimaryHealthy := fm.checkNodeHealth(fm.CurrentPrimary)
            fm.NodeHealthStatus[fm.CurrentPrimary] = currentPrimaryHealthy
            fm.mutex.Unlock()

            if !currentPrimaryHealthy {
                fmt.Printf("Primary node %s is down. Initiating failover...\n", fm.CurrentPrimary)
                fm.failover()
            } else {
                fmt.Printf("Primary node %s is healthy.\n", fm.CurrentPrimary)
            }

            time.Sleep(time.Second * 10) // Check health every 10 seconds
        }
    }()
}

// checkNodeHealth simulates a health check of a node and returns true if the node is healthy
func (fm *NodeFailoverManager) checkNodeHealth(node string) bool {
    // Simulated node health check logic (e.g., network latency, resource usage, etc.)
    fmt.Printf("Performing health check for node %s...\n", node)
    return true // In real implementation, this would include actual health checks
}

// failover switches to the next available backup node
func (fm *NodeFailoverManager) failover() {
    for _, backupNode := range fm.BackupNodes {
        if fm.checkNodeHealth(backupNode) {
            fm.mutex.Lock()
            fm.CurrentPrimary = backupNode
            fm.mutex.Unlock()
            fmt.Printf("Failover complete. Switched to backup node: %s\n", backupNode)
            return
        }
    }

    fmt.Println("ALERT: No available backup nodes. System is at risk.")
}

// AddBackupNode allows adding a new backup node to the system
func (fm *NodeFailoverManager) AddBackupNode(node string) {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    fm.BackupNodes = append(fm.BackupNodes, node)
    fmt.Printf("Added %s as a backup node.\n", node)
}

// RemovePrimaryNode allows removing an unhealthy primary node from the list
func (fm *NodeFailoverManager) RemovePrimaryNode(node string) {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    for i, primaryNode := range fm.PrimaryNodes {
        if primaryNode == node {
            fm.PrimaryNodes = append(fm.PrimaryNodes[:i], fm.PrimaryNodes[i+1:]...)
            fmt.Printf("Removed %s from primary nodes.\n", node)
            return
        }
    }

    fmt.Printf("Node %s not found in primary nodes.\n", node)
}

// GetCurrentPrimaryNode returns the current active primary node
func (fm *NodeFailoverManager) GetCurrentPrimaryNode() string {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    return fm.CurrentPrimary
}
