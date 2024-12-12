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
    UptimeMonitoringInterval  = 15 * time.Second // Interval for monitoring node uptime
    MaxUptimeRetries          = 3                // Maximum retries for node uptime alerts
    SubBlocksPerBlock         = 1000             // Number of sub-blocks in a block
    DowntimeWarningThreshold  = 5 * time.Minute  // Time threshold before a downtime warning is issued
    DowntimeSuspensionThreshold = 15 * time.Minute // Time threshold before a node is suspended
    DowntimeBanThreshold      = 30 * time.Minute // Time threshold before a node is banned for excessive downtime
)

// NodeUptimeMonitoringProtocol tracks node uptime and handles warnings, suspensions, or bans based on downtime
type NodeUptimeMonitoringProtocol struct {
    consensusSystem         *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance          *ledger.Ledger               // Ledger for logging uptime-related events
    stateMutex              *sync.RWMutex                // Mutex for thread-safe access
    uptimeRetryCount        map[string]int               // Counter for retrying uptime alerts
    uptimeCycleCount        int                          // Counter for uptime monitoring cycles
    nodeUptimeTracker       map[string]time.Time         // Tracks the last active time of each node
}

// NewNodeUptimeMonitoringProtocol initializes the automation for node uptime monitoring
func NewNodeUptimeMonitoringProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *NodeUptimeMonitoringProtocol {
    return &NodeUptimeMonitoringProtocol{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        uptimeRetryCount:   make(map[string]int),
        nodeUptimeTracker:  make(map[string]time.Time),
        uptimeCycleCount:   0,
    }
}

// StartUptimeMonitoring starts the continuous loop for monitoring node uptime
func (protocol *NodeUptimeMonitoringProtocol) StartUptimeMonitoring() {
    ticker := time.NewTicker(UptimeMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorNodeUptime()
        }
    }()
}

// monitorNodeUptime checks node uptime and triggers warnings, suspensions, or bans based on downtime thresholds
func (protocol *NodeUptimeMonitoringProtocol) monitorNodeUptime() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch active nodes and their last known active times from the consensus system
    activeNodes := protocol.consensusSystem.FetchNodeUptimeRecords()

    currentTime := time.Now()

    for _, node := range activeNodes {
        lastActive, exists := protocol.nodeUptimeTracker[node.ID]
        if !exists {
            protocol.nodeUptimeTracker[node.ID] = currentTime
            continue
        }

        downtime := currentTime.Sub(lastActive)

        if downtime >= DowntimeBanThreshold {
            fmt.Printf("Node %s has been down for %s. Triggering ban.\n", node.ID, downtime)
            protocol.banNode(node)
        } else if downtime >= DowntimeSuspensionThreshold {
            fmt.Printf("Node %s has been down for %s. Triggering suspension.\n", node.ID, downtime)
            protocol.suspendNode(node)
        } else if downtime >= DowntimeWarningThreshold {
            fmt.Printf("Node %s has been down for %s. Issuing warning.\n", node.ID, downtime)
            protocol.warnNode(node)
        } else {
            fmt.Printf("Node %s is active. Uptime is within acceptable limits.\n", node.ID)
        }
    }

    protocol.uptimeCycleCount++
    fmt.Printf("Uptime monitoring cycle #%d completed.\n", protocol.uptimeCycleCount)

    if protocol.uptimeCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeUptimeCycle()
    }
}

// warnNode issues a warning to a node that has been down for too long
func (protocol *NodeUptimeMonitoringProtocol) warnNode(node common.Node) {
    encryptedUptimeData := protocol.encryptUptimeData(node)

    // Issue a warning through the Synnergy Consensus system
    warningSuccess := protocol.consensusSystem.WarnNode(node, encryptedUptimeData)

    if warningSuccess {
        fmt.Printf("Warning issued to node %s for extended downtime.\n", node.ID)
        protocol.logUptimeEvent(node, "Warning Issued")
        protocol.resetUptimeRetry(node.ID)
    } else {
        fmt.Printf("Error issuing warning to node %s. Retrying...\n", node.ID)
        protocol.retryUptimeAction(node)
    }
}

// suspendNode suspends a node that has exceeded the suspension threshold for downtime
func (protocol *NodeUptimeMonitoringProtocol) suspendNode(node common.Node) {
    encryptedUptimeData := protocol.encryptUptimeData(node)

    // Attempt to suspend the node through the Synnergy Consensus system
    suspensionSuccess := protocol.consensusSystem.SuspendNode(node, encryptedUptimeData)

    if suspensionSuccess {
        fmt.Printf("Node %s suspended for extended downtime.\n", node.ID)
        protocol.logUptimeEvent(node, "Suspended")
        protocol.resetUptimeRetry(node.ID)
    } else {
        fmt.Printf("Error suspending node %s. Retrying...\n", node.ID)
        protocol.retryUptimeAction(node)
    }
}

// banNode bans a node that has exceeded the ban threshold for downtime
func (protocol *NodeUptimeMonitoringProtocol) banNode(node common.Node) {
    encryptedUptimeData := protocol.encryptUptimeData(node)

    // Attempt to ban the node through the Synnergy Consensus system
    banSuccess := protocol.consensusSystem.BanNode(node, encryptedUptimeData)

    if banSuccess {
        fmt.Printf("Node %s banned for excessive downtime.\n", node.ID)
        protocol.logUptimeEvent(node, "Banned")
        protocol.resetUptimeRetry(node.ID)
    } else {
        fmt.Printf("Error banning node %s. Retrying...\n", node.ID)
        protocol.retryUptimeAction(node)
    }
}

// retryUptimeAction retries a failed uptime-related action, such as warnings, suspensions, or bans
func (protocol *NodeUptimeMonitoringProtocol) retryUptimeAction(node common.Node) {
    protocol.uptimeRetryCount[node.ID]++
    if protocol.uptimeRetryCount[node.ID] < MaxUptimeRetries {
        protocol.adjustUptimeAction(node)
    } else {
        fmt.Printf("Max retries reached for node %s. Action failed.\n", node.ID)
        protocol.logUptimeFailure(node)
    }
}

// adjustUptimeAction determines the appropriate action to retry based on the node's downtime
func (protocol *NodeUptimeMonitoringProtocol) adjustUptimeAction(node common.Node) {
    lastActive := protocol.nodeUptimeTracker[node.ID]
    downtime := time.Since(lastActive)

    if downtime >= DowntimeBanThreshold {
        protocol.banNode(node)
    } else if downtime >= DowntimeSuspensionThreshold {
        protocol.suspendNode(node)
    } else if downtime >= DowntimeWarningThreshold {
        protocol.warnNode(node)
    }
}

// resetUptimeRetry resets the retry count for a node's uptime-related action
func (protocol *NodeUptimeMonitoringProtocol) resetUptimeRetry(nodeID string) {
    protocol.uptimeRetryCount[nodeID] = 0
}

// finalizeUptimeCycle finalizes the uptime monitoring cycle and logs the result in the ledger
func (protocol *NodeUptimeMonitoringProtocol) finalizeUptimeCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeUptimeCycle()
    if success {
        fmt.Println("Uptime monitoring cycle finalized successfully.")
        protocol.logUptimeCycleFinalization()
    } else {
        fmt.Println("Error finalizing uptime monitoring cycle.")
    }
}

// logUptimeEvent logs a downtime-related event into the ledger
func (protocol *NodeUptimeMonitoringProtocol) logUptimeEvent(node common.Node, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("uptime-event-%s-%s", node.ID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Uptime Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s was %s due to extended downtime.", node.ID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with uptime event for node %s.\n", node.ID)
}

// logUptimeFailure logs a failure to handle downtime-related actions into the ledger
func (protocol *NodeUptimeMonitoringProtocol) logUptimeFailure(node common.Node) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("uptime-action-failure-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Uptime Action Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to take action for node %s after maximum retries.", node.ID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with uptime action failure for node %s.\n", node.ID)
}

// logUptimeCycleFinalization logs the finalization of an uptime monitoring cycle into the ledger
func (protocol *NodeUptimeMonitoringProtocol) logUptimeCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("uptime-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Uptime Cycle Finalization",
        Status:    "Finalized",
        Details:   "Uptime monitoring cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with uptime cycle finalization.")
}

// encryptUptimeData encrypts the node uptime data before applying warnings, suspensions, or bans
func (protocol *NodeUptimeMonitoringProtocol) encryptUptimeData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node.UptimeData)
    if err != nil {
        fmt.Println("Error encrypting uptime data:", err)
        return node
    }

    node.EncryptedData = encryptedData
    fmt.Println("Uptime data successfully encrypted for node ID:", node.ID)
    return node
}

// triggerEmergencyNodeBan triggers an emergency ban for a node in case of a critical uptime issue
func (protocol *NodeUptimeMonitoringProtocol) triggerEmergencyNodeBan(nodeID string) {
    fmt.Printf("Emergency ban triggered for node ID: %s.\n", nodeID)
    node := protocol.consensusSystem.GetNodeByID(nodeID)
    encryptedData := protocol.encryptUptimeData(node)

    success := protocol.consensusSystem.TriggerEmergencyNodeBan(nodeID, encryptedData)

    if success {
        protocol.logUptimeEvent(node, "Emergency Banned")
        fmt.Println("Emergency ban executed successfully.")
    } else {
        fmt.Println("Emergency ban failed.")
    }
}
