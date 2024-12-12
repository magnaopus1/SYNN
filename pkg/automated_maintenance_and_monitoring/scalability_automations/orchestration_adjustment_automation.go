package automations

import (
    "fmt"
    "time"
    "sync"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/common"
)

const (
    OrchestrationCheckInterval       = 10 * time.Second  // Interval for checking orchestration status
    MaxOrchestrationRetries          = 5                 // Max number of retries for failed orchestration adjustments
    OrchestrationTimeout             = 20 * time.Second  // Timeout for orchestration actions
    PerformanceThreshold             = 0.85              // Performance threshold for triggering adjustments
)

// OrchestrationAdjustmentAutomation handles dynamic orchestration adjustments based on network conditions
type OrchestrationAdjustmentAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to the consensus system
    ledgerInstance        *ledger.Ledger               // Reference to the ledger for tracking events
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    orchestrationAttempts map[string]int               // Tracks retries for orchestration adjustments
    activeOrchestration   map[string]bool              // Active orchestration sessions
}

// NewOrchestrationAdjustmentAutomation initializes orchestration adjustment automation
func NewOrchestrationAdjustmentAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *OrchestrationAdjustmentAutomation {
    return &OrchestrationAdjustmentAutomation{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        orchestrationAttempts: make(map[string]int),
        activeOrchestration:   make(map[string]bool),
    }
}

// StartOrchestrationAdjustment begins continuous monitoring and adjustment of orchestration
func (automation *OrchestrationAdjustmentAutomation) StartOrchestrationAdjustment() {
    ticker := time.NewTicker(OrchestrationCheckInterval)
    go func() {
        for range ticker.C {
            automation.monitorAndAdjustOrchestration()
        }
    }()
}

// monitorAndAdjustOrchestration checks the system’s performance and orchestrates necessary adjustments
func (automation *OrchestrationAdjustmentAutomation) monitorAndAdjustOrchestration() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch performance metrics from the consensus system
    performanceData := automation.consensusSystem.GetPerformanceMetrics()

    for shard, metrics := range performanceData {
        if metrics.PerformanceScore < PerformanceThreshold {
            fmt.Printf("Performance for shard %s below threshold. Initiating orchestration adjustment.\n", shard)
            automation.adjustOrchestration(shard, metrics)
        }
    }
}

// adjustOrchestration triggers orchestration adjustments based on the performance score
func (automation *OrchestrationAdjustmentAutomation) adjustOrchestration(shard string, metrics common.PerformanceMetrics) {
    encryptedMetrics := automation.encryptMetrics(metrics)

    adjustmentSuccess := automation.consensusSystem.ExecuteOrchestrationAdjustment(shard, encryptedMetrics)

    if adjustmentSuccess {
        fmt.Printf("Orchestration adjustment successful for shard %s.\n", shard)
        automation.logOrchestrationEvent(shard, "Adjusted")
        automation.resetOrchestrationRetry(shard)
    } else {
        fmt.Printf("Failed to adjust orchestration for shard %s. Retrying...\n", shard)
        automation.retryOrchestrationAdjustment(shard, metrics)
    }
}

// retryOrchestrationAdjustment retries orchestration adjustment if it initially fails
func (automation *OrchestrationAdjustmentAutomation) retryOrchestrationAdjustment(shard string, metrics common.PerformanceMetrics) {
    automation.orchestrationAttempts[shard]++

    if automation.orchestrationAttempts[shard] < MaxOrchestrationRetries {
        automation.adjustOrchestration(shard, metrics)
    } else {
        fmt.Printf("Max retries reached for orchestration adjustment on shard %s. Marking as failed.\n", shard)
        automation.logAdjustmentFailure(shard)
    }
}

// resetOrchestrationRetry resets the retry count for orchestration adjustments
func (automation *OrchestrationAdjustmentAutomation) resetOrchestrationRetry(shard string) {
    automation.orchestrationAttempts[shard] = 0
}

// encryptMetrics encrypts performance metrics before orchestration adjustments
func (automation *OrchestrationAdjustmentAutomation) encryptMetrics(metrics common.PerformanceMetrics) common.PerformanceMetrics {
    encryptedData, err := encryption.EncryptData(metrics.RawData)
    if err != nil {
        fmt.Println("Error encrypting performance metrics:", err)
        return metrics
    }
    metrics.EncryptedData = encryptedData
    fmt.Println("Performance metrics successfully encrypted.")
    return metrics
}

// logOrchestrationEvent logs a successful orchestration adjustment event to the ledger
func (automation *OrchestrationAdjustmentAutomation) logOrchestrationEvent(shard string, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("orchestration-%s-%s", shard, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Orchestration Adjustment",
        Status:    eventType,
        Details:   fmt.Sprintf("Orchestration adjustment for shard %s marked as %s.", shard, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with orchestration event for shard %s.\n", shard)
}

// logAdjustmentFailure logs a failed orchestration adjustment to the ledger
func (automation *OrchestrationAdjustmentAutomation) logAdjustmentFailure(shard string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("orchestration-failure-%s", shard),
        Timestamp: time.Now().Unix(),
        Type:      "Orchestration Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to adjust orchestration for shard %s after maximum retries.", shard),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with orchestration failure for shard %s.\n", shard)
}

// emergencyOrchestrationAdjustment triggers an emergency orchestration adjustment in critical situations
func (automation *OrchestrationAdjustmentAutomation) emergencyOrchestrationAdjustment(shard string) {
    fmt.Printf("Emergency orchestration adjustment triggered for shard %s.\n", shard)
    success := automation.consensusSystem.TriggerEmergencyOrchestrationAdjustment(shard)

    if success {
        automation.logOrchestrationEvent(shard, "Emergency Adjusted")
        fmt.Println("Emergency orchestration adjustment executed successfully.")
    } else {
        fmt.Println("Emergency orchestration adjustment failed.")
    }
}
