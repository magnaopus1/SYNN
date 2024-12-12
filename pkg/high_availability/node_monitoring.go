package high_availability

import (
    "fmt"
    "time"
    "synnergy_network/pkg/ledger"
)


// NewNodeMonitoringService initializes the node monitoring service
func NewNodeMonitoringService(nodes []string, ledger *ledger.Ledger, checkInterval time.Duration, faultThreshold int) *NodeMonitoringService {
    nodeMetrics := make(map[string]*NodeMetrics)
    for _, node := range nodes {
        nodeMetrics[node] = &NodeMetrics{
            NodeID:        node,
            FaultCount:    0,
            LastChecked:   time.Now(),
            IsHealthy:     true,
            Faulty:        false,
            RecoveryState: false,
        }
    }

    return &NodeMonitoringService{
        Nodes:          nodeMetrics,
        LedgerInstance: ledger,
        CheckInterval:  checkInterval,
        FaultThreshold: faultThreshold,
    }
}

// StartMonitoring begins monitoring nodes, checking their health and recording metrics
func (nms *NodeMonitoringService) StartMonitoring() {
    go func() {
        for {
            nms.mutex.Lock()
            for nodeID, metrics := range nms.Nodes {
                nms.monitorNode(nodeID, metrics)
            }
            nms.mutex.Unlock()
            time.Sleep(nms.CheckInterval)
        }
    }()
}

// monitorNode checks health, latency, and performance of a node, logs data to the ledger
func (nms *NodeMonitoringService) monitorNode(nodeID string, metrics *NodeMetrics) {
    fmt.Printf("Monitoring health for node %s...\n", nodeID)
    
    // Simulate real-world health check and fault detection
    metrics.LastChecked = time.Now()
    metrics.Latency = nms.simulateLatency(nodeID)
    metrics.CPUUsage = nms.simulateCPUUsage(nodeID)

    if nms.isNodeFaulty(metrics) {
        metrics.FaultCount++
        if metrics.FaultCount >= nms.FaultThreshold {
            metrics.IsHealthy = false
            metrics.Faulty = true
            fmt.Printf("Node %s marked as faulty after %d faults.\n", nodeID, metrics.FaultCount)
            nms.triggerRecovery(nodeID, metrics)
        }
    } else {
        metrics.IsHealthy = true
        metrics.FaultCount = 0
        metrics.Faulty = false
    }

    // Log metrics to the ledger, passing the required arguments
    nms.LedgerInstance.LogNodeMetrics(
        metrics.NodeID,        // Node identifier
        metrics.CPUUsage,      // CPU usage
        metrics.MemoryUsage,   // Memory usage (you may need to simulate or add this as well)
        metrics.Latency,       // Latency
    )
}


// isNodeFaulty checks if the node is behaving within acceptable thresholds (simulated)
func (nms *NodeMonitoringService) isNodeFaulty(metrics *NodeMetrics) bool {
    if metrics.Latency > 200 || metrics.CPUUsage > 80 {
        fmt.Printf("Node %s is exhibiting faulty behavior. Latency: %.2fms, CPU: %.2f%%\n", 
                    metrics.NodeID, metrics.Latency, metrics.CPUUsage)
        return true
    }
    return false
}

// simulateLatency simulates network latency for a node (in real-world, use actual ping)
func (nms *NodeMonitoringService) simulateLatency(nodeID string) float64 {
    return float64(50 + len(nodeID)*10) // Example: Add randomness based on nodeID
}

// simulateCPUUsage simulates CPU usage for a node (in real-world, gather actual metrics)
func (nms *NodeMonitoringService) simulateCPUUsage(nodeID string) float64 {
    return float64(10 + len(nodeID)*5) // Example: Add randomness based on nodeID
}

// triggerRecovery initiates the recovery process for a faulty node
func (nms *NodeMonitoringService) triggerRecovery(nodeID string, metrics *NodeMetrics) {
    metrics.RecoveryState = true
    fmt.Printf("Triggering recovery for node %s...\n", nodeID)

    // Perform the recovery process without simulating a delay
    metrics.RecoveryState = false
    metrics.Faulty = false
    metrics.IsHealthy = true
    metrics.FaultCount = 0
    fmt.Printf("Node %s successfully recovered.\n", nodeID)

    // Log recovery to the ledger, formatting time.Now() to a string
    recoveryTime := time.Now().Format(time.RFC3339) // Use RFC3339 or another suitable format
    nms.LedgerInstance.LogNodeRecovery(nodeID, recoveryTime)
}


// GetNodeStatus provides a detailed health report of a node
func (nms *NodeMonitoringService) GetNodeStatus(nodeID string) (*NodeMetrics, error) {
    nms.mutex.Lock()
    defer nms.mutex.Unlock()

    if metrics, exists := nms.Nodes[nodeID]; exists {
        return metrics, nil
    }
    return nil, fmt.Errorf("Node %s not found", nodeID)
}
