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
    TransactionRateMonitoringInterval = 5 * time.Second  // Interval for monitoring transaction rates
    MaxTransactionRateRetries         = 3                // Maximum retries for enforcing rate limits
    SubBlocksPerBlock                 = 1000             // Number of sub-blocks in a block
    TransactionRateThreshold          = 100              // Maximum allowed transactions per second (TPS)
)

// TransactionRateLimitingProtocol monitors and enforces transaction rate limits
type TransactionRateLimitingProtocol struct {
    consensusSystem         *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance          *ledger.Ledger               // Ledger for logging transaction rate events
    stateMutex              *sync.RWMutex                // Mutex for thread-safe access
    transactionRateRetryCount map[string]int             // Counter for retrying transaction rate enforcement
    rateLimitingCycleCount  int                          // Counter for rate limiting monitoring cycles
    transactionRateCounter  map[string]int               // Tracks transaction rate anomalies
}

// NewTransactionRateLimitingProtocol initializes the transaction rate limiting protocol
func NewTransactionRateLimitingProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *TransactionRateLimitingProtocol {
    return &TransactionRateLimitingProtocol{
        consensusSystem:         consensusSystem,
        ledgerInstance:          ledgerInstance,
        stateMutex:              stateMutex,
        transactionRateRetryCount: make(map[string]int),
        transactionRateCounter:  make(map[string]int),
        rateLimitingCycleCount:  0,
    }
}

// StartRateLimitingMonitoring starts the continuous loop for monitoring transaction rates
func (protocol *TransactionRateLimitingProtocol) StartRateLimitingMonitoring() {
    ticker := time.NewTicker(TransactionRateMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorTransactionRates()
        }
    }()
}

// monitorTransactionRates checks for rate limit violations and takes appropriate actions
func (protocol *TransactionRateLimitingProtocol) monitorTransactionRates() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch the current transaction rate reports from the consensus system
    rateReports := protocol.consensusSystem.FetchTransactionRateReports()

    for _, report := range rateReports {
        if protocol.isTransactionRateExceeded(report) {
            fmt.Printf("Transaction rate limit exceeded for node %s. Taking action.\n", report.NodeID)
            protocol.handleRateLimitViolation(report)
        } else {
            fmt.Printf("Transaction rate within acceptable limits for node %s.\n", report.NodeID)
        }
    }

    protocol.rateLimitingCycleCount++
    fmt.Printf("Transaction rate monitoring cycle #%d completed.\n", protocol.rateLimitingCycleCount)

    if protocol.rateLimitingCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeRateMonitoringCycle()
    }
}

// isTransactionRateExceeded checks if the transaction rate exceeds the threshold
func (protocol *TransactionRateLimitingProtocol) isTransactionRateExceeded(report common.TransactionRateReport) bool {
    if report.TransactionRate > TransactionRateThreshold {
        fmt.Printf("Transaction rate exceeded for node %s. Current rate: %d TPS\n", report.NodeID, report.TransactionRate)
        return true
    }
    return false
}

// handleRateLimitViolation handles a transaction rate limit violation by throttling or alerting
func (protocol *TransactionRateLimitingProtocol) handleRateLimitViolation(report common.TransactionRateReport) {
    protocol.transactionRateCounter[report.NodeID]++

    if protocol.transactionRateCounter[report.NodeID] >= MaxTransactionRateRetries {
        fmt.Printf("Multiple rate violations detected for node %s. Escalating response.\n", report.NodeID)
        protocol.escalateRateLimitViolation(report)
    } else {
        fmt.Printf("Issuing alert for rate limit violation in node %s.\n", report.NodeID)
        protocol.alertForRateLimitViolation(report)
    }
}

// alertForRateLimitViolation issues an alert for a detected rate limit violation
func (protocol *TransactionRateLimitingProtocol) alertForRateLimitViolation(report common.TransactionRateReport) {
    encryptedAlertData := protocol.encryptRateLimitData(report)

    // Issue an alert through the Synnergy Consensus system
    alertSuccess := protocol.consensusSystem.IssueRateLimitViolationAlert(report.NodeID, encryptedAlertData)

    if alertSuccess {
        fmt.Printf("Rate limit violation alert issued for node %s.\n", report.NodeID)
        protocol.logRateLimitEvent(report, "Alert Issued")
        protocol.resetRateRetry(report.NodeID)
    } else {
        fmt.Printf("Error issuing rate limit violation alert for node %s. Retrying...\n", report.NodeID)
        protocol.retryRateResponse(report)
    }
}

// escalateRateLimitViolation escalates the response to a persistent rate limit violation
func (protocol *TransactionRateLimitingProtocol) escalateRateLimitViolation(report common.TransactionRateReport) {
    encryptedEscalationData := protocol.encryptRateLimitData(report)

    // Attempt to escalate the response through the Synnergy Consensus system
    escalationSuccess := protocol.consensusSystem.EscalateRateLimitViolationResponse(report.NodeID, encryptedEscalationData)

    if escalationSuccess {
        fmt.Printf("Rate limit violation response escalated for node %s.\n", report.NodeID)
        protocol.logRateLimitEvent(report, "Response Escalated")
        protocol.resetRateRetry(report.NodeID)
    } else {
        fmt.Printf("Error escalating rate limit violation response for node %s. Retrying...\n", report.NodeID)
        protocol.retryRateResponse(report)
    }
}

// retryRateResponse retries the response to a rate limit violation if the initial action fails
func (protocol *TransactionRateLimitingProtocol) retryRateResponse(report common.TransactionRateReport) {
    protocol.transactionRateRetryCount[report.NodeID]++
    if protocol.transactionRateRetryCount[report.NodeID] < MaxTransactionRateRetries {
        protocol.escalateRateLimitViolation(report)
    } else {
        fmt.Printf("Max retries reached for rate limit violation response for node %s. Response failed.\n", report.NodeID)
        protocol.logRateLimitFailure(report)
    }
}

// resetRateRetry resets the retry count for rate limit violations for a specific node ID
func (protocol *TransactionRateLimitingProtocol) resetRateRetry(nodeID string) {
    protocol.transactionRateRetryCount[nodeID] = 0
}

// finalizeRateMonitoringCycle finalizes the transaction rate monitoring cycle and logs the result in the ledger
func (protocol *TransactionRateLimitingProtocol) finalizeRateMonitoringCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeRateMonitoringCycle()
    if success {
        fmt.Println("Transaction rate monitoring cycle finalized successfully.")
        protocol.logRateMonitoringCycleFinalization()
    } else {
        fmt.Println("Error finalizing transaction rate monitoring cycle.")
    }
}

// logRateLimitEvent logs a rate limit event into the ledger
func (protocol *TransactionRateLimitingProtocol) logRateLimitEvent(report common.TransactionRateReport, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("rate-limit-event-%s-%s", report.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Transaction Rate Limit Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s triggered %s due to transaction rate violation.", report.NodeID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with rate limit violation event for node %s.\n", report.NodeID)
}

// logRateLimitFailure logs the failure to respond to a rate limit violation into the ledger
func (protocol *TransactionRateLimitingProtocol) logRateLimitFailure(report common.TransactionRateReport) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("rate-limit-failure-%s", report.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Rate Limit Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to respond to transaction rate violation for node %s after maximum retries.", report.NodeID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with rate limit failure for node %s.\n", report.NodeID)
}

// logRateMonitoringCycleFinalization logs the finalization of a rate monitoring cycle into the ledger
func (protocol *TransactionRateLimitingProtocol) logRateMonitoringCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("rate-monitoring-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Rate Monitoring Cycle Finalization",
        Status:    "Finalized",
        Details:   "Transaction rate monitoring cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with rate monitoring cycle finalization.")
}

// encryptRateLimitData encrypts rate limit-related data before taking action or logging events
func (protocol *TransactionRateLimitingProtocol) encryptRateLimitData(report common.TransactionRateReport) common.TransactionRateReport {
    encryptedData, err := encryption.EncryptData(report.RateLimitData)
    if err != nil {
        fmt.Println("Error encrypting rate limit data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("Rate limit data successfully encrypted for node ID:", report.NodeID)
    return report
}
