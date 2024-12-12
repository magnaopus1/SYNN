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
    PerformanceMonitoringInterval   = 10 * time.Second // Interval for monitoring validator performance
    MaxPerformanceViolationRetries  = 3                // Maximum retries for handling performance violations
    SubBlocksPerBlock               = 1000             // Number of sub-blocks in a block
    ValidatorPerformanceThreshold   = 70               // Performance score threshold for acceptable validators
)

// ValidatorPerformanceSecurityProtocol monitors validator performance and enforces network policies
type ValidatorPerformanceSecurityProtocol struct {
    consensusSystem      *consensus.SynnergyConsensus  // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger                // Ledger for logging performance-related events
    stateMutex           *sync.RWMutex                 // Mutex for thread-safe access
    violationRetryCount  map[string]int                // Counter for retrying responses to performance violations
    performanceCycleCount int                          // Counter for performance monitoring cycles
}

// NewValidatorPerformanceSecurityProtocol initializes the performance protocol
func NewValidatorPerformanceSecurityProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ValidatorPerformanceSecurityProtocol {
    return &ValidatorPerformanceSecurityProtocol{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        violationRetryCount:  make(map[string]int),
        performanceCycleCount: 0,
    }
}

// StartPerformanceMonitoring begins the continuous loop for monitoring validator performance
func (protocol *ValidatorPerformanceSecurityProtocol) StartPerformanceMonitoring() {
    ticker := time.NewTicker(PerformanceMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorValidatorPerformance()
        }
    }()
}

// monitorValidatorPerformance checks for validators performing below the required threshold
func (protocol *ValidatorPerformanceSecurityProtocol) monitorValidatorPerformance() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch performance reports for validators
    performanceReports := protocol.consensusSystem.FetchValidatorPerformanceReports()

    for _, report := range performanceReports {
        if protocol.isPerformanceViolationDetected(report) {
            fmt.Printf("Performance violation detected for validator %s. Taking action.\n", report.ValidatorID)
            protocol.handlePerformanceViolation(report)
        } else {
            fmt.Printf("Validator %s performance is within acceptable limits.\n", report.ValidatorID)
        }
    }

    protocol.performanceCycleCount++
    fmt.Printf("Validator performance monitoring cycle #%d completed.\n", protocol.performanceCycleCount)

    if protocol.performanceCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizePerformanceCycle()
    }
}

// isPerformanceViolationDetected checks if a validator's performance violates the network policy
func (protocol *ValidatorPerformanceSecurityProtocol) isPerformanceViolationDetected(report common.ValidatorPerformanceReport) bool {
    if report.PerformanceScore < ValidatorPerformanceThreshold {
        fmt.Printf("Validator %s has a performance score of %d, which is below the threshold.\n", report.ValidatorID, report.PerformanceScore)
        return true
    }
    return false
}

// handlePerformanceViolation responds to a validator performance violation and logs the event
func (protocol *ValidatorPerformanceSecurityProtocol) handlePerformanceViolation(report common.ValidatorPerformanceReport) {
    protocol.violationRetryCount[report.ValidatorID]++

    if protocol.violationRetryCount[report.ValidatorID] >= MaxPerformanceViolationRetries {
        fmt.Printf("Multiple performance violations detected for validator %s. Escalating response.\n", report.ValidatorID)
        protocol.escalatePerformanceViolation(report)
    } else {
        fmt.Printf("Handling performance violation for validator %s.\n", report.ValidatorID)
        protocol.issuePerformanceSanction(report)
    }
}

// issuePerformanceSanction takes corrective actions on validators performing below the required threshold
func (protocol *ValidatorPerformanceSecurityProtocol) issuePerformanceSanction(report common.ValidatorPerformanceReport) {
    encryptedSanctionData := protocol.encryptPerformanceData(report)

    // Apply performance sanction through the Synnergy Consensus system
    sanctionSuccess := protocol.consensusSystem.ApplyPerformanceSanction(report.ValidatorID, encryptedSanctionData)

    if sanctionSuccess {
        fmt.Printf("Performance sanction applied to validator %s.\n", report.ValidatorID)
        protocol.logPerformanceEvent(report, "Performance Sanction Applied")
        protocol.resetViolationRetry(report.ValidatorID)
    } else {
        fmt.Printf("Error applying performance sanction to validator %s. Retrying...\n", report.ValidatorID)
        protocol.retryPerformanceViolationResponse(report)
    }
}

// escalatePerformanceViolation escalates performance violations that remain unresolved
func (protocol *ValidatorPerformanceSecurityProtocol) escalatePerformanceViolation(report common.ValidatorPerformanceReport) {
    encryptedEscalationData := protocol.encryptPerformanceData(report)

    // Escalate the performance violation through the Synnergy Consensus system
    escalationSuccess := protocol.consensusSystem.EscalatePerformanceViolation(report.ValidatorID, encryptedEscalationData)

    if escalationSuccess {
        fmt.Printf("Performance violation escalated for validator %s.\n", report.ValidatorID)
        protocol.logPerformanceEvent(report, "Performance Violation Escalated")
        protocol.resetViolationRetry(report.ValidatorID)
    } else {
        fmt.Printf("Error escalating performance violation for validator %s. Retrying...\n", report.ValidatorID)
        protocol.retryPerformanceViolationResponse(report)
    }
}

// retryPerformanceViolationResponse retries the response if initial attempts fail
func (protocol *ValidatorPerformanceSecurityProtocol) retryPerformanceViolationResponse(report common.ValidatorPerformanceReport) {
    protocol.violationRetryCount[report.ValidatorID]++
    if protocol.violationRetryCount[report.ValidatorID] < MaxPerformanceViolationRetries {
        protocol.issuePerformanceSanction(report)
    } else {
        fmt.Printf("Max retries reached for handling performance violation of validator %s. Response failed.\n", report.ValidatorID)
        protocol.logPerformanceFailure(report)
    }
}

// resetViolationRetry resets the retry count for handling performance violations for a specific validator
func (protocol *ValidatorPerformanceSecurityProtocol) resetViolationRetry(validatorID string) {
    protocol.violationRetryCount[validatorID] = 0
}

// finalizePerformanceCycle finalizes the monitoring cycle and logs the result in the ledger
func (protocol *ValidatorPerformanceSecurityProtocol) finalizePerformanceCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizePerformanceMonitoringCycle()
    if success {
        fmt.Println("Validator performance monitoring cycle finalized successfully.")
        protocol.logPerformanceCycleFinalization()
    } else {
        fmt.Println("Error finalizing performance monitoring cycle.")
    }
}

// logPerformanceEvent logs a performance-related event into the ledger
func (protocol *ValidatorPerformanceSecurityProtocol) logPerformanceEvent(report common.ValidatorPerformanceReport, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("performance-event-%s-%s", report.ValidatorID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Validator Performance Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Validator %s triggered %s due to performance issues.", report.ValidatorID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with performance violation event for validator %s.\n", report.ValidatorID)
}

// logPerformanceFailure logs the failure to resolve a performance violation into the ledger
func (protocol *ValidatorPerformanceSecurityProtocol) logPerformanceFailure(report common.ValidatorPerformanceReport) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("performance-failure-%s", report.ValidatorID),
        Timestamp: time.Now().Unix(),
        Type:      "Performance Violation Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to resolve performance violation for validator %s after maximum retries.", report.ValidatorID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with performance violation failure for validator %s.\n", report.ValidatorID)
}

// logPerformanceCycleFinalization logs the finalization of a performance monitoring cycle into the ledger
func (protocol *ValidatorPerformanceSecurityProtocol) logPerformanceCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("performance-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Performance Monitoring Cycle Finalization",
        Status:    "Finalized",
        Details:   "Validator performance monitoring cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with performance monitoring cycle finalization.")
}

// encryptPerformanceData encrypts performance-related data before taking action or logging events
func (protocol *ValidatorPerformanceSecurityProtocol) encryptPerformanceData(report common.ValidatorPerformanceReport) common.ValidatorPerformanceReport {
    encryptedData, err := encryption.EncryptData(report.ValidatorData)
    if err != nil {
        fmt.Println("Error encrypting validator data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("Validator data successfully encrypted for validator:", report.ValidatorID)
    return report
}
