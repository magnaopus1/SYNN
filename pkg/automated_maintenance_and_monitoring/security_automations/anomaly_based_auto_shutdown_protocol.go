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
    AnomalyDetectionInterval = 10 * time.Minute // Time interval for detecting anomalies
    MaxAllowedAnomalies      = 5                // Maximum number of anomalies before triggering auto-shutdown
    SubBlocksPerBlock        = 1000             // Number of sub-blocks in a block
)

// AnomalyBasedAutoShutdownAutomation automates shutdown upon detecting anomalies
type AnomalyBasedAutoShutdownAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging shutdown events
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    anomalyCounter    int                          // Counter for detected anomalies
}

// NewAnomalyBasedAutoShutdownAutomation initializes the automation for detecting and responding to anomalies
func NewAnomalyBasedAutoShutdownAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AnomalyBasedAutoShutdownAutomation {
    return &AnomalyBasedAutoShutdownAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        anomalyCounter:  0,
    }
}

// StartAnomalyDetection starts the continuous loop for anomaly detection and auto-shutdown
func (automation *AnomalyBasedAutoShutdownAutomation) StartAnomalyDetection() {
    ticker := time.NewTicker(AnomalyDetectionInterval)

    go func() {
        for range ticker.C {
            automation.detectAndHandleAnomalies()
        }
    }()
}

// detectAndHandleAnomalies detects system anomalies and triggers a shutdown if the threshold is exceeded
func (automation *AnomalyBasedAutoShutdownAutomation) detectAndHandleAnomalies() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the number of detected anomalies from the consensus system
    detectedAnomalies := automation.consensusSystem.DetectAnomalies()

    if detectedAnomalies > 0 {
        fmt.Printf("Detected %d anomalies. Monitoring...\n", detectedAnomalies)
        automation.anomalyCounter += detectedAnomalies

        if automation.anomalyCounter >= MaxAllowedAnomalies {
            fmt.Println("Maximum allowed anomalies exceeded. Initiating auto-shutdown...")
            automation.shutdownSystem()
        }
    } else {
        fmt.Println("No anomalies detected in this cycle.")
    }

    // Log anomaly detection in the ledger
    automation.logAnomalyDetection(detectedAnomalies)
}

// shutdownSystem triggers a system shutdown and logs the event in the ledger
func (automation *AnomalyBasedAutoShutdownAutomation) shutdownSystem() {
    success := automation.consensusSystem.Shutdown()
    if success {
        fmt.Println("System successfully shut down due to anomaly detection.")
        automation.logShutdownEvent()
        automation.anomalyCounter = 0
    } else {
        fmt.Println("Error during system shutdown.")
    }
}

// logAnomalyDetection logs the number of anomalies detected into the ledger
func (automation *AnomalyBasedAutoShutdownAutomation) logAnomalyDetection(detectedAnomalies int) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("anomaly-detection-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Anomaly Detection",
        Status:    "Detected",
        Details:   fmt.Sprintf("Detected %d anomalies.", detectedAnomalies),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with anomaly detection event.")
}

// logShutdownEvent logs the system shutdown event into the ledger
func (automation *AnomalyBasedAutoShutdownAutomation) logShutdownEvent() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("system-shutdown-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "System Shutdown",
        Status:    "Completed",
        Details:   "System shutdown initiated due to exceeding anomaly threshold.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with system shutdown event.")
}

// finalizeAnomalyDetectionCycle finalizes the anomaly detection cycle and logs it in the ledger
func (automation *AnomalyBasedAutoShutdownAutomation) finalizeAnomalyDetectionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeAnomalyDetectionCycle()
    if success {
        fmt.Println("Anomaly detection cycle finalized successfully.")
        automation.logFinalization()
    } else {
        fmt.Println("Error finalizing anomaly detection cycle.")
    }
}

// logFinalization logs the finalization of an anomaly detection cycle into the ledger
func (automation *AnomalyBasedAutoShutdownAutomation) logFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("anomaly-detection-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Anomaly Detection Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with anomaly detection cycle finalization.")
}
