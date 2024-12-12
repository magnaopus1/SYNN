package high_availability

import (
    "fmt"
    "time"
    "synnergy_network/pkg/ledger"
)


// NewHeartbeatService initializes a new HeartbeatService
func NewHeartbeatService(nodes []string, interval time.Duration, ledger *ledger.Ledger) *HeartbeatService {
    return &HeartbeatService{
        Nodes:         nodes,
        HeartbeatLogs: make(map[string]time.Time),
        Interval:      interval,
        LedgerInstance: ledger,
    }
}

// Start sends heartbeats to all nodes and monitors their responses
func (hb *HeartbeatService) Start() {
    for _, node := range hb.Nodes {
        go func(node string) {
            for {
                hb.sendHeartbeat(node)
                time.Sleep(hb.Interval)
            }
        }(node)
    }

    go hb.monitorHeartbeats() // Start monitoring heartbeats for responsiveness
}

// sendHeartbeat simulates sending a heartbeat signal to a node
func (hb *HeartbeatService) sendHeartbeat(node string) {
    fmt.Printf("Sending heartbeat to node: %s\n", node)

    // Simulating heartbeat response from node (In a real-world system, this would be an actual network request)
    hb.mutex.Lock()
    hb.HeartbeatLogs[node] = time.Now()
    hb.mutex.Unlock()
}

// monitorHeartbeats checks if any node has missed its heartbeat within the allowed time frame
func (hb *HeartbeatService) monitorHeartbeats() {
    for {
        time.Sleep(hb.Interval)

        hb.mutex.Lock()
        for node, lastHeartbeat := range hb.HeartbeatLogs {
            if time.Since(lastHeartbeat) > hb.Interval*2 {
                fmt.Printf("ALERT: Node %s is unresponsive.\n", node)
                hb.triggerNodeRecovery(node)
            }
        }
        hb.mutex.Unlock()
    }
}

// triggerNodeRecovery simulates initiating a recovery process for an unresponsive node
func (hb *HeartbeatService) triggerNodeRecovery(node string) {
    fmt.Printf("Initiating recovery process for node: %s\n", node)
    // In a real-world system, this would integrate with node recovery protocols or alert mechanisms
}

// GetLastHeartbeat returns the last time a heartbeat was received from a node
func (hb *HeartbeatService) GetLastHeartbeat(node string) time.Time {
    hb.mutex.Lock()
    defer hb.mutex.Unlock()

    return hb.HeartbeatLogs[node]
}

// IsNodeResponsive checks if a node is responsive based on the last received heartbeat
func (hb *HeartbeatService) IsNodeResponsive(node string) bool {
    hb.mutex.Lock()
    defer hb.mutex.Unlock()

    lastHeartbeat, exists := hb.HeartbeatLogs[node]
    if !exists || time.Since(lastHeartbeat) > hb.Interval*2 {
        return false
    }
    return true
}
