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
    OrchestrationInterval           = 10 * time.Second  // Frequency of orchestration checks
    MaxRetries                      = 3                 // Maximum retries for orchestration actions
    OrchestrationThreshold          = 0.9               // Performance threshold for triggering orchestration
    OrchestrationTimeout            = 30 * time.Second  // Timeout for orchestration actions
)

// OrchestrationAutomation handles dynamic orchestration tasks such as resource allocation and optimization.
type OrchestrationAutomation struct {
    ledgerInstance        *ledger.Ledger               // Ledger instance for logging
    consensusSystem       *consensus.SynnergyConsensus // Consensus system instance
    orchestrationAttempts map[string]int               // Track orchestration attempts
    stateMutex            *sync.RWMutex                // Mutex for safe state handling
}

// NewOrchestrationAutomation initializes the orchestration automation process.
func NewOrchestrationAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *OrchestrationAutomation {
    return &OrchestrationAutomation{
        ledgerInstance:        ledgerInstance,
        consensusSystem:       consensusSystem,
        orchestrationAttempts: make(map[string]int),
        stateMutex:            stateMutex,
    }
}

// StartAutomation triggers the orchestration automation loop
func (automation *OrchestrationAutomation) StartAutomation() {
    ticker := time.NewTicker(OrchestrationInterval)
    go func() {
        for range ticker.C {
            automation.monitorOrchestration()
        }
    }()
}

// monitorOrchestration checks network performance and adjusts orchestration accordingly.
func (automation *OrchestrationAutomation) monitorOrchestration() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    performanceMetrics := automation.consensusSystem.GetNetworkMetrics()

    for shard, metrics := range performanceMetrics {
        if metrics.PerformanceScore < OrchestrationThreshold {
            fmt.Printf("Orchestration adjustment needed for shard %s\n", shard)
            automation.executeOrchestration(shard, metrics)
        }
    }
}

// executeOrchestration performs the orchestration tasks for underperforming shards.
func (automation *OrchestrationAutomation) executeOrchestration(shard string, metrics common.PerformanceMetrics) {
    encryptedMetrics := automation.encryptMetrics(metrics)

    success := automation.consensusSystem.AdjustOrchestration(shard, encryptedMetrics)

    if success {
        fmt.Printf("Orchestration adjustment successful for shard %s\n", shard)
        automation.logOrchestrationSuccess(shard)
        automation.resetRetries(shard)
    } else {
        fmt.Printf("Orchestration adjustment failed for shard %s, retrying...\n", shard)
        automation.handleRetry(shard, metrics)
    }
}

// handleRetry manages retries for orchestration failures.
func (automation *OrchestrationAutomation) handleRetry(shard string, metrics common.PerformanceMetrics) {
    automation.orchestrationAttempts[shard]++

    if automation.orchestrationAttempts[shard] < MaxRetries {
        automation.executeOrchestration(shard, metrics)
    } else {
        fmt.Printf("Max retries reached for shard %s, marking as failed.\n", shard)
        automation.logOrchestrationFailure(shard)
    }
}

// encryptMetrics encrypts performance metrics before sending for orchestration adjustments.
func (automation *OrchestrationAutomation) encryptMetrics(metrics common.PerformanceMetrics) common.PerformanceMetrics {
    encryptedData, err := encryption.EncryptData(metrics.RawData)
    if err != nil {
        fmt.Println("Error encrypting performance metrics:", err)
        return metrics
    }
    metrics.EncryptedData = encryptedData
    fmt.Println("Performance metrics encrypted.")
    return metrics
}

// logOrchestrationSuccess logs successful orchestration adjustments to the ledger.
func (automation *OrchestrationAutomation) logOrchestrationSuccess(shard string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("orchestration-success-%s", shard),
        Timestamp: time.Now().Unix(),
        Type:      "Orchestration Success",
        Status:    "Completed",
        Details:   fmt.Sprintf("Orchestration adjustment completed for shard %s.", shard),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Orchestration success logged for shard %s\n", shard)
}

// logOrchestrationFailure logs failed orchestration adjustments to the ledger.
func (automation *OrchestrationAutomation) logOrchestrationFailure(shard string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("orchestration-failure-%s", shard),
        Timestamp: time.Now().Unix(),
        Type:      "Orchestration Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to adjust orchestration for shard %s after retries.", shard),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Orchestration failure logged for shard %s\n", shard)
}

// resetRetries resets the retry count for orchestration.
func (automation *OrchestrationAutomation) resetRetries(shard string) {
    automation.orchestrationAttempts[shard] = 0
}

// EmergencyOrchestration triggers an emergency orchestration when critical conditions are detected.
func (automation *OrchestrationAutomation) EmergencyOrchestration(shard string) {
    fmt.Printf("Emergency orchestration triggered for shard %s\n", shard)
    success := automation.consensusSystem.TriggerEmergencyOrchestration(shard)

    if success {
        automation.logOrchestrationSuccess(shard)
        fmt.Println("Emergency orchestration executed successfully.")
    } else {
        automation.logOrchestrationFailure(shard)
        fmt.Println("Emergency orchestration failed.")
    }
}
