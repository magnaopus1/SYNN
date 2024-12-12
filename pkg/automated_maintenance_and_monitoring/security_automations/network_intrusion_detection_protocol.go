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
    NetworkIntrusionMonitoringInterval = 15 * time.Second // Interval for checking network intrusions
    MaxIntrusionResponseRetries        = 3                // Maximum retries for responding to intrusions
    SubBlocksPerBlock                  = 1000             // Number of sub-blocks in a block
    IntrusionAlertThreshold            = 5                // Threshold for alerting on multiple suspicious activities
)

// NetworkIntrusionDetectionProtocol handles the detection, alerting, and prevention of network intrusions
type NetworkIntrusionDetectionProtocol struct {
    consensusSystem           *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance            *ledger.Ledger               // Ledger for logging intrusion-related events
    stateMutex                *sync.RWMutex                // Mutex for thread-safe access
    intrusionRetryCount       map[string]int               // Counter for retrying responses to intrusion attempts
    intrusionCycleCount       int                          // Counter for intrusion detection cycles
    suspiciousActivityTracker map[string]int               // Tracks suspicious activity per node
}

// NewNetworkIntrusionDetectionProtocol initializes the automation for network intrusion detection
func NewNetworkIntrusionDetectionProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *NetworkIntrusionDetectionProtocol {
    return &NetworkIntrusionDetectionProtocol{
        consensusSystem:         consensusSystem,
        ledgerInstance:          ledgerInstance,
        stateMutex:              stateMutex,
        intrusionRetryCount:     make(map[string]int),
        suspiciousActivityTracker: make(map[string]int),
        intrusionCycleCount:     0,
    }
}

// StartIntrusionMonitoring starts the continuous loop for monitoring network intrusions
func (protocol *NetworkIntrusionDetectionProtocol) StartIntrusionMonitoring() {
    ticker := time.NewTicker(NetworkIntrusionMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorForIntrusions()
        }
    }()
}

// monitorForIntrusions checks for suspicious network activities and takes action if an intrusion is detected
func (protocol *NetworkIntrusionDetectionProtocol) monitorForIntrusions() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch suspicious activities from the consensus system
    suspiciousActivities := protocol.consensusSystem.DetectSuspiciousActivities()

    if len(suspiciousActivities) > 0 {
        for _, activity := range suspiciousActivities {
            fmt.Printf("Suspicious activity detected for node %s. Investigating.\n", activity.NodeID)
            protocol.handleSuspiciousActivity(activity)
        }
    } else {
        fmt.Println("No suspicious activities detected this cycle.")
    }

    protocol.intrusionCycleCount++
    fmt.Printf("Network intrusion detection cycle #%d completed.\n", protocol.intrusionCycleCount)

    if protocol.intrusionCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeIntrusionDetectionCycle()
    }
}

// handleSuspiciousActivity takes action based on suspicious activity, such as alerting or blocking
func (protocol *NetworkIntrusionDetectionProtocol) handleSuspiciousActivity(activity common.NetworkActivity) {
    protocol.suspiciousActivityTracker[activity.NodeID]++

    // If suspicious activity exceeds the alert threshold, take action
    if protocol.suspiciousActivityTracker[activity.NodeID] >= IntrusionAlertThreshold {
        fmt.Printf("Multiple suspicious activities detected for node %s. Triggering response.\n", activity.NodeID)
        protocol.respondToIntrusion(activity)
    } else {
        fmt.Printf("Suspicious activity count for node %s: %d\n", activity.NodeID, protocol.suspiciousActivityTracker[activity.NodeID])
    }
}

// respondToIntrusion handles the response to an intrusion attempt by applying security measures
func (protocol *NetworkIntrusionDetectionProtocol) respondToIntrusion(activity common.NetworkActivity) {
    encryptedIntrusionData := protocol.encryptIntrusionData(activity)

    // Attempt to block or mitigate the intrusion through the Synnergy Consensus system
    responseSuccess := protocol.consensusSystem.BlockSuspiciousNode(activity.NodeID, encryptedIntrusionData)

    if responseSuccess {
        fmt.Printf("Intrusion response successful for node %s.\n", activity.NodeID)
        protocol.logIntrusionEvent(activity, "Blocked")
        protocol.resetIntrusionRetry(activity.NodeID)
        protocol.clearSuspiciousActivityCount(activity.NodeID)
    } else {
        fmt.Printf("Error responding to intrusion for node %s. Retrying...\n", activity.NodeID)
        protocol.retryIntrusionResponse(activity)
    }
}

// retryIntrusionResponse retries the response to an intrusion attempt if the first response fails
func (protocol *NetworkIntrusionDetectionProtocol) retryIntrusionResponse(activity common.NetworkActivity) {
    protocol.intrusionRetryCount[activity.NodeID]++
    if protocol.intrusionRetryCount[activity.NodeID] < MaxIntrusionResponseRetries {
        protocol.respondToIntrusion(activity)
    } else {
        fmt.Printf("Max retries reached for responding to intrusion by node %s. Response failed.\n", activity.NodeID)
        protocol.logIntrusionFailure(activity)
    }
}

// resetIntrusionRetry resets the retry count for a node after a successful response
func (protocol *NetworkIntrusionDetectionProtocol) resetIntrusionRetry(nodeID string) {
    protocol.intrusionRetryCount[nodeID] = 0
}

// clearSuspiciousActivityCount clears the suspicious activity count for a node after an action is taken
func (protocol *NetworkIntrusionDetectionProtocol) clearSuspiciousActivityCount(nodeID string) {
    delete(protocol.suspiciousActivityTracker, nodeID)
}

// finalizeIntrusionDetectionCycle finalizes the intrusion detection cycle and logs the result in the ledger
func (protocol *NetworkIntrusionDetectionProtocol) finalizeIntrusionDetectionCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeIntrusionDetectionCycle()
    if success {
        fmt.Println("Intrusion detection cycle finalized successfully.")
        protocol.logIntrusionCycleFinalization()
    } else {
        fmt.Println("Error finalizing intrusion detection cycle.")
    }
}

// logIntrusionEvent logs an intrusion response event into the ledger
func (protocol *NetworkIntrusionDetectionProtocol) logIntrusionEvent(activity common.NetworkActivity, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("intrusion-event-%s-%s", activity.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Intrusion Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s was %s due to suspicious activity.", activity.NodeID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with intrusion event for node %s.\n", activity.NodeID)
}

// logIntrusionFailure logs the failure of an intrusion response into the ledger
func (protocol *NetworkIntrusionDetectionProtocol) logIntrusionFailure(activity common.NetworkActivity) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("intrusion-failure-%s", activity.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Intrusion Response Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to respond to intrusion from node %s after maximum retries.", activity.NodeID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with intrusion failure for node %s.\n", activity.NodeID)
}

// logIntrusionCycleFinalization logs the finalization of an intrusion detection cycle into the ledger
func (protocol *NetworkIntrusionDetectionProtocol) logIntrusionCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("intrusion-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Intrusion Cycle Finalization",
        Status:    "Finalized",
        Details:   "Intrusion detection cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with intrusion detection cycle finalization.")
}

// encryptIntrusionData encrypts intrusion data before applying security actions
func (protocol *NetworkIntrusionDetectionProtocol) encryptIntrusionData(activity common.NetworkActivity) common.NetworkActivity {
    encryptedData, err := encryption.EncryptData(activity.Data)
    if err != nil {
        fmt.Println("Error encrypting intrusion data:", err)
        return activity
    }

    activity.EncryptedData = encryptedData
    fmt.Println("Intrusion data successfully encrypted for node ID:", activity.NodeID)
    return activity
}

// triggerEmergencyIntrusionLockdown triggers an emergency lockdown on a node in response to critical suspicious activities
func (protocol *NetworkIntrusionDetectionProtocol) triggerEmergencyIntrusionLockdown(nodeID string) {
    fmt.Printf("Emergency lockdown triggered for node ID: %s.\n", nodeID)
    nodeActivity := protocol.consensusSystem.GetNodeActivityByID(nodeID)
    encryptedData := protocol.encryptIntrusionData(nodeActivity)

    success := protocol.consensusSystem.TriggerEmergencyNodeLockdown(nodeID, encryptedData)

    if success {
        protocol.logIntrusionEvent(nodeActivity, "Emergency Locked")
        fmt.Println("Emergency intrusion lockdown executed successfully.")
    } else {
        fmt.Println("Emergency intrusion lockdown failed.")
    }
}
