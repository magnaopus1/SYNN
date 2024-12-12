package consensus_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    HealthReportInterval    = 10 * time.Minute  // Interval for generating health reports
    HealthReportKey         = "health_report_key" // Encryption key for health reports
    BlockFinalizationLimit  = 15 * time.Second  // Limit for acceptable block finalization time
    ValidatorPerformanceThreshold = 0.75        // Threshold for acceptable validator performance
)

// SynnergyConsensusHealthReportingAutomation generates periodic health reports for the Synnergy Consensus
type SynnergyConsensusHealthReportingAutomation struct {
    ledgerInstance  *ledger.Ledger                    // Blockchain ledger for logging health reports
    consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine to gather metrics
    stateMutex      *sync.RWMutex                     // Mutex for thread-safe ledger access
}

// NewSynnergyConsensusHealthReportingAutomation initializes the health reporting automation
func NewSynnergyConsensusHealthReportingAutomation(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *SynnergyConsensusHealthReportingAutomation {
    return &SynnergyConsensusHealthReportingAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
    }
}

// StartHealthReporting starts the periodic generation of health reports for the Synnergy Consensus
func (automation *SynnergyConsensusHealthReportingAutomation) StartHealthReporting() {
    ticker := time.NewTicker(HealthReportInterval)

    go func() {
        for range ticker.C {
            fmt.Println("Generating health report for Synnergy Consensus...")
            automation.generateHealthReport()
        }
    }()
}

// generateHealthReport collects and compiles metrics from PoH, PoS, and PoW stages
func (automation *SynnergyConsensusHealthReportingAutomation) generateHealthReport() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Gather PoH metrics (e.g., block finalization time)
    pohFinalizationTime := automation.consensusEngine.GetPoHFinalizationTime()

    // Gather PoS metrics (e.g., validator performance)
    validatorPerformance := automation.consensusEngine.GetValidatorPerformance()

    // Gather PoW metrics (e.g., mining difficulty)
    powDifficulty := automation.consensusEngine.GetPoWDifficulty()

    // Assess the health status based on gathered metrics
    isConsensusHealthy := automation.evaluateConsensusHealth(pohFinalizationTime, validatorPerformance, powDifficulty)

    // Compile the report with all gathered metrics
    healthReport := common.HealthReport{
        Timestamp:           time.Now().Unix(),
        PoHFinalizationTime: pohFinalizationTime,
        ValidatorPerformance: validatorPerformance,
        PoWDifficulty:       powDifficulty,
        ConsensusHealthy:    isConsensusHealthy,
    }

    // Log and encrypt the report into the ledger
    automation.logHealthReport(healthReport)
}

// evaluateConsensusHealth assesses the overall consensus health based on the metrics
func (automation *SynnergyConsensusHealthReportingAutomation) evaluateConsensusHealth(pohFinalizationTime time.Duration, validatorPerformance, powDifficulty float64) bool {
    // Check if block finalization time exceeds the acceptable limit
    if pohFinalizationTime > BlockFinalizationLimit {
        fmt.Println("Warning: Block finalization time exceeds the limit.")
        return false
    }

    // Check if validator performance falls below the threshold
    if validatorPerformance < ValidatorPerformanceThreshold {
        fmt.Println("Warning: Validator performance is below acceptable thresholds.")
        return false
    }

    // Check if PoW difficulty is too high (customizable logic based on network state)
    if powDifficulty > 10 {
        fmt.Println("Warning: Mining difficulty is too high.")
        return false
    }

    return true
}

// logHealthReport encrypts and logs the health report into the ledger for record-keeping
func (automation *SynnergyConsensusHealthReportingAutomation) logHealthReport(healthReport common.HealthReport) {
    encryptedReport, err := encryption.EncryptHealthReport(healthReport, []byte(HealthReportKey))
    if err != nil {
        fmt.Printf("Error encrypting health report: %v\n", err)
        return
    }

    // Add the encrypted report to the ledger
    automation.ledgerInstance.AddEntry(common.LedgerEntry{
        ID:        fmt.Sprintf("health-report-%d", time.Now().Unix()),
        Timestamp: healthReport.Timestamp,
        Type:      "Health Report",
        Status:    "Generated",
        Details:   "Consensus health report generated and encrypted.",
    })

    fmt.Println("Health report successfully generated and stored in the ledger.")
}

// triggerHealthReport allows manual or triggered health report generation
func (automation *SynnergyConsensusHealthReportingAutomation) TriggerHealthReport() {
    fmt.Println("Triggering manual health report generation...")
    automation.generateHealthReport()
}

// sendHealthReport allows sending the health report to system administrators or external systems
func (automation *SynnergyConsensusHealthReportingAutomation) sendHealthReport(report common.HealthReport) {
    // Code to send the health report to system administrators or external systems (e.g., via email, dashboard, etc.)
    fmt.Printf("Sending health report: %+v\n", report)
}

// Additional trigger that can generate a health report when certain conditions are met
func (automation *SynnergyConsensusHealthReportingAutomation) TriggerReportOnIssueDetected() {
    // Example: Trigger a report if PoW difficulty exceeds a certain threshold or validator performance drops
    if automation.consensusEngine.GetPoWDifficulty() > 15 || automation.consensusEngine.GetValidatorPerformance() < 0.5 {
        fmt.Println("Critical condition detected, generating health report.")
        automation.generateHealthReport()
    }
}
