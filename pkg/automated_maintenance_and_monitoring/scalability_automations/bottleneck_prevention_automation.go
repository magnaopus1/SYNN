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
    PreventionCheckInterval  = 2 * time.Minute // Interval for bottleneck prevention checks
    PreventionThreshold      = 80              // Percentage threshold to trigger bottleneck prevention
    MaxPreventionRetries     = 3               // Maximum retries for prevention measures
)

// BottleneckPreventionAutomation handles bottleneck prevention by taking preemptive measures
type BottleneckPreventionAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger for logging prevention-related events
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    preventionRetryCount  map[string]int               // Counter for retrying bottleneck prevention
    preventionCycleCount  int                          // Counter for bottleneck prevention cycles
}

// NewBottleneckPreventionAutomation initializes the automation for bottleneck prevention
func NewBottleneckPreventionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *BottleneckPreventionAutomation {
    return &BottleneckPreventionAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        preventionRetryCount: make(map[string]int),
        preventionCycleCount: 0,
    }
}

// StartBottleneckPrevention starts the continuous loop for detecting and preventing bottlenecks
func (automation *BottleneckPreventionAutomation) StartBottleneckPrevention() {
    ticker := time.NewTicker(PreventionCheckInterval)

    go func() {
        for range ticker.C {
            automation.detectAndPreventBottlenecks()
        }
    }()
}

// detectAndPreventBottlenecks checks for potential bottlenecks and triggers prevention measures
func (automation *BottleneckPreventionAutomation) detectAndPreventBottlenecks() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch system reports from the consensus system
    systemReports := automation.consensusSystem.FetchUsageReports()

    for _, report := range systemReports {
        if report.Utilization >= PreventionThreshold {
            fmt.Printf("High utilization detected on node %s with %d%% usage. Taking preventive action.\n", report.NodeID, report.Utilization)
            automation.preventBottleneck(report)
        } else {
            fmt.Printf("Node %s is within normal utilization at %d%%.\n", report.NodeID, report.Utilization)
        }
    }

    automation.preventionCycleCount++
    fmt.Printf("Bottleneck prevention cycle #%d completed.\n", automation.preventionCycleCount)

    if automation.preventionCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizePreventionCycle()
    }
}

// preventBottleneck attempts to take preemptive measures to prevent bottlenecks
func (automation *BottleneckPreventionAutomation) preventBottleneck(report common.UsageReport) {
    encryptedData := automation.encryptPreventionData(report)

    // Attempt to take preventive action via the Synnergy Consensus system
    success := automation.consensusSystem.InitiatePreventionAction(report.NodeID, encryptedData)
    if success {
        fmt.Printf("Prevention action successful for node %s.\n", report.NodeID)
        automation.logPreventionEvent(report, "Prevented")
        automation.resetPreventionRetry(report.NodeID)
    } else {
        fmt.Printf("Prevention action failed for node %s. Retrying...\n", report.NodeID)
        automation.retryPreventionAction(report)
    }
}

// retryPreventionAction retries bottleneck prevention actions if initial attempts fail
func (automation *BottleneckPreventionAutomation) retryPreventionAction(report common.UsageReport) {
    automation.preventionRetryCount[report.NodeID]++
    if automation.preventionRetryCount[report.NodeID] < MaxPreventionRetries {
        automation.preventBottleneck(report)
    } else {
        fmt.Printf("Max retries reached for node %s. Prevention action failed.\n", report.NodeID)
        automation.logPreventionFailure(report)
    }
}

// resetPreventionRetry resets the retry count for bottleneck prevention
func (automation *BottleneckPreventionAutomation) resetPreventionRetry(nodeID string) {
    automation.preventionRetryCount[nodeID] = 0
}

// encryptPreventionData encrypts the prevention data before initiating action
func (automation *BottleneckPreventionAutomation) encryptPreventionData(report common.UsageReport) []byte {
    encryptedData, err := encryption.EncryptData(report.Data)
    if err != nil {
        fmt.Printf("Error encrypting prevention data for node %s: %v\n", report.NodeID, err)
        return report.Data
    }

    fmt.Printf("Prevention data successfully encrypted for node %s.\n", report.NodeID)
    return encryptedData
}

// logPreventionEvent logs a successful prevention action into the ledger
func (automation *BottleneckPreventionAutomation) logPreventionEvent(report common.UsageReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("bottleneck-prevention-%s-%s", report.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Bottleneck Prevention Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s bottleneck %s successfully.", report.NodeID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with bottleneck prevention event for node %s.\n", report.NodeID)
}

// logPreventionFailure logs a failure to prevent a bottleneck for a specific node into the ledger
func (automation *BottleneckPreventionAutomation) logPreventionFailure(report common.UsageReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("bottleneck-prevention-failure-%s", report.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Bottleneck Prevention Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Prevention action failed for node %s after maximum retries.", report.NodeID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with prevention failure for node %s.\n", report.NodeID)
}

// finalizePreventionCycle finalizes the bottleneck prevention cycle and logs the result in the ledger
func (automation *BottleneckPreventionAutomation) finalizePreventionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizePreventionCycle()
    if success {
        fmt.Println("Bottleneck prevention cycle finalized successfully.")
        automation.logPreventionCycleFinalization()
    } else {
        fmt.Println("Error finalizing bottleneck prevention cycle.")
    }
}

// logPreventionCycleFinalization logs the finalization of a bottleneck prevention cycle into the ledger
func (automation *BottleneckPreventionAutomation) logPreventionCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("bottleneck-prevention-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Bottleneck Prevention Cycle",
        Status:    "Finalized",
        Details:   "Bottleneck prevention cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with bottleneck prevention cycle finalization.")
}
