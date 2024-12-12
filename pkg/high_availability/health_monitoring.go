package high_availability

import (
    "fmt"
    "time"
    "synnergy_network/pkg/ledger"

)


// NewHealthMonitoringManager initializes the HealthMonitoringManager
func NewHealthMonitoringManager(nodes []string, ledgerInstance *ledger.Ledger) *HealthMonitoringManager {
    return &HealthMonitoringManager{
        Nodes:            nodes,
        NodeHealthStatus: make(map[string]bool),
        SubBlockLatency:  make(map[int]time.Duration),
        BlockLatency:     make(map[int]time.Duration),
        LedgerInstance:   ledgerInstance,
    }
}

// MonitorNodeHealth continuously checks the health of nodes in the network
func (hm *HealthMonitoringManager) MonitorNodeHealth() {
    for _, node := range hm.Nodes {
        go func(node string) {
            for {
                hm.mutex.Lock()
                healthy := hm.checkNodeHealth(node)
                hm.NodeHealthStatus[node] = healthy
                hm.mutex.Unlock()

                if healthy {
                    fmt.Printf("Node %s is healthy.\n", node)
                } else {
                    fmt.Printf("Node %s is unhealthy.\n", node)
                    hm.TriggerNodeAlert(node)
                }

                time.Sleep(time.Second * 10) // Check node health every 10 seconds
            }
        }(node)
    }
}

// checkNodeHealth simulates a health check of the node and returns true if the node is healthy
func (hm *HealthMonitoringManager) checkNodeHealth(node string) bool {
    // Simulated node health check logic (e.g., network latency, response time, etc.)
    fmt.Printf("Performing health check for node %s...\n", node)
    return true // In a real implementation, this would include actual health checks
}

// MonitorSubBlockValidation tracks the time taken to validate each sub-block
func (hm *HealthMonitoringManager) MonitorSubBlockValidation(subBlockIndex int, start time.Time) {
    hm.mutex.Lock()
    defer hm.mutex.Unlock()

    latency := time.Since(start)
    hm.SubBlockLatency[subBlockIndex] = latency
    fmt.Printf("Sub-block %d validated in %s.\n", subBlockIndex, latency)
}

// MonitorBlockValidation tracks the time taken to validate each block
func (hm *HealthMonitoringManager) MonitorBlockValidation(blockIndex int, start time.Time) {
    hm.mutex.Lock()
    defer hm.mutex.Unlock()

    latency := time.Since(start)
    hm.BlockLatency[blockIndex] = latency
    fmt.Printf("Block %d validated in %s.\n", blockIndex, latency)
}

// TriggerNodeAlert raises an alert if a node is found to be unhealthy
func (hm *HealthMonitoringManager) TriggerNodeAlert(node string) {
    fmt.Printf("ALERT: Node %s is experiencing issues. Please investigate.\n", node)
    // In a real-world system, this could be integrated with an alerting system (e.g., email, SMS, etc.)
}

// GetNodeHealthStatus returns the current health status of a node
func (hm *HealthMonitoringManager) GetNodeHealthStatus(node string) bool {
    hm.mutex.Lock()
    defer hm.mutex.Unlock()

    return hm.NodeHealthStatus[node]
}

// GetSubBlockLatency retrieves the validation latency for a specific sub-block
func (hm *HealthMonitoringManager) GetSubBlockLatency(subBlockIndex int) time.Duration {
    hm.mutex.Lock()
    defer hm.mutex.Unlock()

    return hm.SubBlockLatency[subBlockIndex]
}

// GetBlockLatency retrieves the validation latency for a specific block
func (hm *HealthMonitoringManager) GetBlockLatency(blockIndex int) time.Duration {
    hm.mutex.Lock()
    defer hm.mutex.Unlock()

    return hm.BlockLatency[blockIndex]
}
