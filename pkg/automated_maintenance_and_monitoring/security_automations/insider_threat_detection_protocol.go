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
    InsiderThreatMonitoringInterval = 20 * time.Second // Interval for monitoring insider threat activities
    MaxInsiderAlertRetries          = 5                // Maximum retries for alerting about insider threats
    SubBlocksPerBlock               = 1000             // Number of sub-blocks in a block
)

// InsiderThreatDetectionProtocol monitors for insider threats within the network
type InsiderThreatDetectionProtocol struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger for logging insider threat-related events
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    insiderAlertRetries   map[string]int               // Counter for retrying insider threat alerts
    detectionCycleCount   int                          // Counter for insider threat detection cycles
    insiderSuspiciousData map[string]common.Node       // Cache of suspicious activity for validation
}

// NewInsiderThreatDetectionProtocol initializes the automation for detecting insider threats
func NewInsiderThreatDetectionProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *InsiderThreatDetectionProtocol {
    return &InsiderThreatDetectionProtocol{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        insiderAlertRetries:   make(map[string]int),
        insiderSuspiciousData: make(map[string]common.Node),
        detectionCycleCount:   0,
    }
}

// StartInsiderThreatMonitoring starts the continuous loop for monitoring insider threats within the network
func (protocol *InsiderThreatDetectionProtocol) StartInsiderThreatMonitoring() {
    ticker := time.NewTicker(InsiderThreatMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.detectAndRespondToInsiderThreats()
        }
    }()
}

// detectAndRespondToInsiderThreats checks the blockchain for insider threat indicators and takes action accordingly
func (protocol *InsiderThreatDetectionProtocol) detectAndRespondToInsiderThreats() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch insider threat indicators through consensus
    insiderThreatActivities := protocol.consensusSystem.DetectInsiderThreats()

    if len(insiderThreatActivities) > 0 {
        for _, activity := range insiderThreatActivities {
            fmt.Printf("Insider threat detected for node %s. Investigating.\n", activity.NodeID)
            protocol.handleInsiderThreat(activity)
        }
    } else {
        fmt.Println("No insider threats detected this cycle.")
    }

    protocol.detectionCycleCount++
    fmt.Printf("Insider threat detection cycle #%d completed.\n", protocol.detectionCycleCount)

    if protocol.detectionCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeInsiderThreatCycle()
    }
}

// handleInsiderThreat takes action in response to detected insider threats, such as blocking or alerting
func (protocol *InsiderThreatDetectionProtocol) handleInsiderThreat(activity common.Node) {
    encryptedInsiderData := protocol.encryptInsiderData(activity)

    // Attempt to apply insider threat response through the Synnergy Consensus system
    responseSuccess := protocol.consensusSystem.RespondToInsiderThreat(activity, encryptedInsiderData)

    if responseSuccess {
        fmt.Printf("Insider threat response applied successfully for node %s.\n", activity.NodeID)
        protocol.logInsiderThreatEvent(activity, "Responded")
        protocol.resetInsiderAlertRetry(activity.NodeID)
    } else {
        fmt.Printf("Error responding to insider threat for node %s. Retrying...\n", activity.NodeID)
        protocol.retryInsiderThreatAlert(activity)
    }
}

// retryInsiderThreatAlert retries an insider threat alert if the initial response fails
func (protocol *InsiderThreatDetectionProtocol) retryInsiderThreatAlert(activity common.Node) {
    protocol.insiderAlertRetries[activity.NodeID]++
    if protocol.insiderAlertRetries[activity.NodeID] < MaxInsiderAlertRetries {
        protocol.handleInsiderThreat(activity)
    } else {
        fmt.Printf("Max retries reached for insider threat alert on node %s. Alert failed.\n", activity.NodeID)
        protocol.logInsiderThreatFailure(activity)
    }
}

// resetInsiderAlertRetry resets the retry count for insider threat alerts
func (protocol *InsiderThreatDetectionProtocol) resetInsiderAlertRetry(nodeID string) {
    protocol.insiderAlertRetries[nodeID] = 0
}

// finalizeInsiderThreatCycle finalizes the insider threat detection cycle and logs the result in the ledger
func (protocol *InsiderThreatDetectionProtocol) finalizeInsiderThreatCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeInsiderThreatCycle()
    if success {
        fmt.Println("Insider threat detection cycle finalized successfully.")
        protocol.logInsiderThreatCycleFinalization()
    } else {
        fmt.Println("Error finalizing insider threat detection cycle.")
    }
}

// logInsiderThreatEvent logs an insider threat event into the ledger
func (protocol *InsiderThreatDetectionProtocol) logInsiderThreatEvent(activity common.Node, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("insider-threat-%s-%s", activity.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Insider Threat Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s %s due to insider threat activity.", activity.NodeID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with insider threat event for node %s.\n", activity.NodeID)
}

// logInsiderThreatFailure logs the failure of an insider threat response into the ledger
func (protocol *InsiderThreatDetectionProtocol) logInsiderThreatFailure(activity common.Node) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("insider-threat-failure-%s", activity.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Insider Threat Response Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to respond to insider threat on node %s after maximum retries.", activity.NodeID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with insider threat response failure for node %s.\n", activity.NodeID)
}

// logInsiderThreatCycleFinalization logs the finalization of an insider threat detection cycle into the ledger
func (protocol *InsiderThreatDetectionProtocol) logInsiderThreatCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("insider-threat-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Insider Threat Cycle Finalization",
        Status:    "Finalized",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with insider threat detection cycle finalization.")
}

// encryptInsiderData encrypts insider threat data before applying protective actions
func (protocol *InsiderThreatDetectionProtocol) encryptInsiderData(activity common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(activity.NodeData)
    if err != nil {
        fmt.Println("Error encrypting insider threat data:", err)
        return activity
    }

    activity.EncryptedData = encryptedData
    fmt.Println("Insider threat data successfully encrypted for protective actions.")
    return activity
}

// triggerEmergencyLockdown triggers an emergency lockdown of a node in case of critical insider threat detection
func (protocol *InsiderThreatDetectionProtocol) triggerEmergencyLockdown(node common.Node) {
    fmt.Printf("Emergency lockdown triggered for node %s.\n", node.NodeID)
    success := protocol.consensusSystem.TriggerEmergencyLockdown(node)

    if success {
        protocol.logInsiderThreatEvent(node, "Locked Down")
        fmt.Println("Emergency lockdown executed successfully.")
    } else {
        fmt.Println("Emergency lockdown failed.")
    }
}
