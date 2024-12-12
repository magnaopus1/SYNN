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
    BottleneckCheckInterval     = 5 * time.Minute // Interval for checking bottlenecks
    MaxBottleneckRetries        = 3               // Maximum retries for resolving bottlenecks
    BottleneckThreshold         = 90              // Percentage threshold for identifying bottlenecks
    BottleneckResolutionTimeout = 15 * time.Minute // Timeout for resolving bottlenecks
)

// BottleneckDetectionAndResolution handles the detection and resolution of network bottlenecks
type BottleneckDetectionAndResolution struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger for logging bottleneck-related events
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    bottleneckRetryCount map[string]int               // Counter for retrying bottleneck resolutions
    detectionCycle       int                          // Counter for bottleneck detection cycles
}

// NewBottleneckDetectionAndResolution initializes the automation for bottleneck detection and resolution
func NewBottleneckDetectionAndResolution(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *BottleneckDetectionAndResolution {
    return &BottleneckDetectionAndResolution{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        bottleneckRetryCount: make(map[string]int),
        detectionCycle:       0,
    }
}

// StartBottleneckDetectionAndResolution starts the continuous loop for detecting and resolving bottlenecks
func (automation *BottleneckDetectionAndResolution) StartBottleneckDetectionAndResolution() {
    ticker := time.NewTicker(BottleneckCheckInterval)

    go func() {
        for range ticker.C {
            automation.detectAndResolveBottlenecks()
        }
    }()
}

// detectAndResolveBottlenecks checks for bottlenecks in the network and attempts to resolve them
func (automation *BottleneckDetectionAndResolution) detectAndResolveBottlenecks() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch reports on system bottlenecks from the consensus system
    bottleneckReports := automation.consensusSystem.FetchBottleneckReports()

    for _, report := range bottleneckReports {
        if report.Utilization >= BottleneckThreshold {
            fmt.Printf("Bottleneck detected on node %s with utilization at %d%%. Resolving...\n", report.NodeID, report.Utilization)
            automation.resolveBottleneck(report)
        } else {
            fmt.Printf("Node %s is within normal limits with utilization at %d%%.\n", report.NodeID, report.Utilization)
        }
    }

    automation.detectionCycle++
    fmt.Printf("Bottleneck detection cycle #%d completed.\n", automation.detectionCycle)

    if automation.detectionCycle%SubBlocksPerBlock == 0 {
        automation.finalizeDetectionCycle()
    }
}

// resolveBottleneck attempts to resolve bottlenecks on a given node
func (automation *BottleneckDetectionAndResolution) resolveBottleneck(report common.BottleneckReport) {
    encryptedData := automation.encryptBottleneckData(report)

    // Attempt to resolve the bottleneck through the Synnergy Consensus system
    success := automation.consensusSystem.ResolveBottleneck(report.NodeID, encryptedData)
    if success {
        fmt.Printf("Bottleneck resolution successful for node %s.\n", report.NodeID)
        automation.logResolutionEvent(report, "Resolved")
        automation.resetBottleneckRetry(report.NodeID)
    } else {
        fmt.Printf("Bottleneck resolution failed for node %s. Retrying...\n", report.NodeID)
        automation.retryBottleneckResolution(report)
    }
}

// retryBottleneckResolution retries the resolution of bottlenecks if initial attempts fail
func (automation *BottleneckDetectionAndResolution) retryBottleneckResolution(report common.BottleneckReport) {
    automation.bottleneckRetryCount[report.NodeID]++
    if automation.bottleneckRetryCount[report.NodeID] < MaxBottleneckRetries {
        automation.resolveBottleneck(report)
    } else {
        fmt.Printf("Max retries reached for node %s. Bottleneck resolution failed.\n", report.NodeID)
        automation.logResolutionFailure(report)
    }
}

// resetBottleneckRetry resets the retry count for bottleneck resolutions
func (automation *BottleneckDetectionAndResolution) resetBottleneckRetry(nodeID string) {
    automation.bottleneckRetryCount[nodeID] = 0
}

// encryptBottleneckData encrypts the bottleneck data before attempting to resolve it
func (automation *BottleneckDetectionAndResolution) encryptBottleneckData(report common.BottleneckReport) []byte {
    encryptedData, err := encryption.EncryptData(report.Data)
    if err != nil {
        fmt.Printf("Error encrypting data for node %s: %v\n", report.NodeID, err)
        return report.Data
    }

    fmt.Printf("Bottleneck data successfully encrypted for node %s.\n", report.NodeID)
    return encryptedData
}

// logResolutionEvent logs a successful bottleneck resolution event into the ledger
func (automation *BottleneckDetectionAndResolution) logResolutionEvent(report common.BottleneckReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("bottleneck-resolution-%s-%s", report.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Bottleneck Resolution Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s bottleneck %s successfully.", report.NodeID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with bottleneck resolution event for node %s.\n", report.NodeID)
}

// logResolutionFailure logs the failure to resolve a bottleneck for a specific node into the ledger
func (automation *BottleneckDetectionAndResolution) logResolutionFailure(report common.BottleneckReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("bottleneck-resolution-failure-%s", report.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Bottleneck Resolution Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Bottleneck resolution failed for node %s after maximum retries.", report.NodeID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with bottleneck resolution failure for node %s.\n", report.NodeID)
}

// finalizeDetectionCycle finalizes the bottleneck detection cycle and logs the result in the ledger
func (automation *BottleneckDetectionAndResolution) finalizeDetectionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeBottleneckCycle()
    if success {
        fmt.Println("Bottleneck detection cycle finalized successfully.")
        automation.logCycleFinalization()
    } else {
        fmt.Println("Error finalizing bottleneck detection cycle.")
    }
}

// logCycleFinalization logs the finalization of a bottleneck detection cycle into the ledger
func (automation *BottleneckDetectionAndResolution) logCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("bottleneck-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Bottleneck Detection Cycle",
        Status:    "Finalized",
        Details:   "Bottleneck detection cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with bottleneck detection cycle finalization.")
}
