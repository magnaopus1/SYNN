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
    BandwidthMonitoringInterval = 10 * time.Second // Interval for monitoring bandwidth usage
    MaxThrottleRetries          = 3                // Maximum retries for applying throttling
    SubBlocksPerBlock           = 1000             // Number of sub-blocks in a block
)

// BandwidthThrottlingSecurityAutomation manages bandwidth usage to ensure security and performance
type BandwidthThrottlingSecurityAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger for logging bandwidth-related events
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    throttleRetryCount map[string]int               // Counter for retrying bandwidth throttling on failure
    throttleCycleCount int                          // Counter for bandwidth monitoring cycles
}

// NewBandwidthThrottlingSecurityAutomation initializes the automation for bandwidth throttling security
func NewBandwidthThrottlingSecurityAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *BandwidthThrottlingSecurityAutomation {
    return &BandwidthThrottlingSecurityAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        throttleRetryCount: make(map[string]int),
        throttleCycleCount: 0,
    }
}

// StartBandwidthMonitoring starts the continuous loop for monitoring bandwidth usage
func (automation *BandwidthThrottlingSecurityAutomation) StartBandwidthMonitoring() {
    ticker := time.NewTicker(BandwidthMonitoringInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndThrottleBandwidth()
        }
    }()
}

// monitorAndThrottleBandwidth continuously monitors the network for bandwidth overuse and applies throttling if necessary
func (automation *BandwidthThrottlingSecurityAutomation) monitorAndThrottleBandwidth() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of nodes or components with bandwidth issues
    bandwidthOveruseList := automation.consensusSystem.DetectBandwidthOveruse()

    if len(bandwidthOveruseList) > 0 {
        for _, node := range bandwidthOveruseList {
            fmt.Printf("Bandwidth overuse detected for node %s. Applying throttling.\n", node.ID)
            automation.applyThrottling(node)
        }
    } else {
        fmt.Println("No bandwidth issues detected this cycle.")
    }

    automation.throttleCycleCount++
    fmt.Printf("Bandwidth monitoring cycle #%d completed.\n", automation.throttleCycleCount)

    if automation.throttleCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeThrottleCycle()
    }
}

// applyThrottling attempts to throttle bandwidth for the given node or component
func (automation *BandwidthThrottlingSecurityAutomation) applyThrottling(node common.Node) {
    encryptedThrottleData := automation.encryptThrottleData(node)

    // Attempt to throttle the node through the Synnergy Consensus system
    throttleSuccess := automation.consensusSystem.ThrottleNode(node, encryptedThrottleData)

    if throttleSuccess {
        fmt.Printf("Bandwidth throttling applied successfully for node %s.\n", node.ID)
        automation.logThrottleEvent(node, "Throttled")
        automation.resetThrottleRetry(node.ID)
    } else {
        fmt.Printf("Error throttling node %s. Retrying...\n", node.ID)
        automation.retryThrottling(node)
    }
}

// retryThrottling attempts to retry a failed throttling action a limited number of times
func (automation *BandwidthThrottlingSecurityAutomation) retryThrottling(node common.Node) {
    automation.throttleRetryCount[node.ID]++
    if automation.throttleRetryCount[node.ID] < MaxThrottleRetries {
        automation.applyThrottling(node)
    } else {
        fmt.Printf("Max retries reached for throttling node %s. Throttling failed.\n", node.ID)
        automation.logThrottleFailure(node)
    }
}

// resetThrottleRetry resets the retry count for throttling actions
func (automation *BandwidthThrottlingSecurityAutomation) resetThrottleRetry(nodeID string) {
    automation.throttleRetryCount[nodeID] = 0
}

// finalizeThrottleCycle finalizes the bandwidth throttling cycle and logs the result in the ledger
func (automation *BandwidthThrottlingSecurityAutomation) finalizeThrottleCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeThrottleCycle()
    if success {
        fmt.Println("Bandwidth throttling cycle finalized successfully.")
        automation.logThrottleCycleFinalization()
    } else {
        fmt.Println("Error finalizing bandwidth throttling cycle.")
    }
}

// logThrottleEvent logs a throttling event into the ledger
func (automation *BandwidthThrottlingSecurityAutomation) logThrottleEvent(node common.Node, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("throttle-%s-%s", node.ID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Bandwidth Throttling Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s %s successfully.", node.ID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with throttling event for node %s.\n", node.ID)
}

// logThrottleFailure logs the failure of throttling a specific node into the ledger
func (automation *BandwidthThrottlingSecurityAutomation) logThrottleFailure(node common.Node) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("throttle-failure-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Bandwidth Throttling Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Throttling failed for node %s after maximum retries.", node.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with throttling failure for node %s.\n", node.ID)
}

// logThrottleCycleFinalization logs the finalization of a throttling cycle into the ledger
func (automation *BandwidthThrottlingSecurityAutomation) logThrottleCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("throttle-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Throttling Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with throttling cycle finalization.")
}

// encryptThrottleData encrypts throttling data before applying throttling
func (automation *BandwidthThrottlingSecurityAutomation) encryptThrottleData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node.Data)
    if err != nil {
        fmt.Println("Error encrypting node data for throttling:", err)
        return node
    }

    node.EncryptedData = encryptedData
    fmt.Println("Node data successfully encrypted for throttling.")
    return node
}

// manualIntervention allows for manual intervention in the throttling process
func (automation *BandwidthThrottlingSecurityAutomation) manualIntervention(node common.Node, action string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    if action == "throttle" {
        fmt.Printf("Manually throttling node %s.\n", node.ID)
        automation.applyThrottling(node)
    } else if action == "ignore" {
        fmt.Printf("Manually ignoring bandwidth overuse for node %s.\n", node.ID)
    } else {
        fmt.Println("Invalid action for manual intervention.")
    }
}

// emergencyUnthrottling triggers emergency unthrottling of a node in case of critical needs
func (automation *BandwidthThrottlingSecurityAutomation) emergencyUnthrottling(node common.Node) {
    fmt.Printf("Emergency unthrottling triggered for node %s.\n", node.ID)
    success := automation.consensusSystem.TriggerEmergencyUnthrottling(node)

    if success {
        automation.logThrottleEvent(node, "Unthrottled")
        fmt.Println("Emergency unthrottling executed successfully.")
    } else {
        fmt.Println("Emergency unthrottling failed.")
    }
}
