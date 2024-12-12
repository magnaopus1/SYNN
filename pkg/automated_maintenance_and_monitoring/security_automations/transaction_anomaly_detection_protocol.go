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
    TransactionAnomalyMonitoringInterval = 5 * time.Second  // Interval for monitoring transaction anomalies
    MaxTransactionRetries                = 3                // Maximum retries for responding to detected anomalies
    SubBlocksPerBlock                    = 1000             // Number of sub-blocks in a block
    TransactionAnomalyThreshold          = 0.15             // Threshold for detecting transaction anomalies
)

// TransactionAnomalyDetectionProtocol manages the detection and handling of transaction anomalies
type TransactionAnomalyDetectionProtocol struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger for logging anomaly-related events
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    transactionRetryCount  map[string]int               // Counter for retrying anomaly handling
    transactionCycleCount  int                          // Counter for transaction monitoring cycles
    transactionAnomalyCounter map[string]int            // Tracks detected transaction anomalies
}

// NewTransactionAnomalyDetectionProtocol initializes the transaction anomaly detection protocol
func NewTransactionAnomalyDetectionProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *TransactionAnomalyDetectionProtocol {
    return &TransactionAnomalyDetectionProtocol{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        transactionRetryCount:  make(map[string]int),
        transactionAnomalyCounter: make(map[string]int),
        transactionCycleCount:  0,
    }
}

// StartTransactionAnomalyMonitoring starts the continuous loop for monitoring transaction anomalies
func (protocol *TransactionAnomalyDetectionProtocol) StartTransactionAnomalyMonitoring() {
    ticker := time.NewTicker(TransactionAnomalyMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorTransactionAnomalies()
        }
    }()
}

// monitorTransactionAnomalies checks for transaction anomalies and takes action accordingly
func (protocol *TransactionAnomalyDetectionProtocol) monitorTransactionAnomalies() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch transaction reports from the consensus system
    transactionReports := protocol.consensusSystem.FetchTransactionReports()

    for _, report := range transactionReports {
        if protocol.isTransactionAnomalyDetected(report) {
            fmt.Printf("Transaction anomaly detected for transaction ID %s. Taking action.\n", report.TransactionID)
            protocol.handleTransactionAnomaly(report)
        } else {
            fmt.Printf("No anomaly detected for transaction ID %s.\n", report.TransactionID)
        }
    }

    protocol.transactionCycleCount++
    fmt.Printf("Transaction anomaly monitoring cycle #%d completed.\n", protocol.transactionCycleCount)

    if protocol.transactionCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeTransactionMonitoringCycle()
    }
}

// isTransactionAnomalyDetected checks if there is an anomaly in the transaction report
func (protocol *TransactionAnomalyDetectionProtocol) isTransactionAnomalyDetected(report common.TransactionReport) bool {
    // Check if the anomaly score exceeds the threshold
    if report.AnomalyScore >= TransactionAnomalyThreshold {
        fmt.Printf("Anomaly score exceeded threshold for transaction ID: %s. Score: %f\n", report.TransactionID, report.AnomalyScore)
        return true
    }

    // Detect suspicious transaction patterns (e.g., unusually large amounts or abnormal transaction frequency)
    if report.IsSuspiciousAmount || report.IsAbnormalTransactionFrequency {
        fmt.Printf("Suspicious transaction pattern detected for transaction ID: %s\n", report.TransactionID)
        return true
    }

    // Detect if malicious or blacklisted addresses are involved
    if protocol.isMaliciousAddressInvolved(report) {
        fmt.Printf("Malicious address involvement detected for transaction ID: %s\n", report.TransactionID)
        return true
    }

    // If none of the conditions match, no anomaly is detected
    fmt.Printf("No anomaly detected for transaction ID: %s\n", report.TransactionID)
    return false
}

// isMaliciousAddressInvolved checks if any malicious or blacklisted addresses are involved in the transaction
func (protocol *TransactionAnomalyDetectionProtocol) isMaliciousAddressInvolved(report common.TransactionReport) bool {
    for _, addr := range report.InvolvedAddresses {
        if protocol.consensusSystem.IsBlacklistedAddress(addr) {
            fmt.Printf("Blacklisted address %s detected in transaction ID %s.\n", addr, report.TransactionID)
            return true
        }
    }
    return false
}

// handleTransactionAnomaly takes action when a transaction anomaly is detected
func (protocol *TransactionAnomalyDetectionProtocol) handleTransactionAnomaly(report common.TransactionReport) {
    protocol.transactionAnomalyCounter[report.TransactionID]++

    if protocol.transactionAnomalyCounter[report.TransactionID] >= MaxTransactionRetries {
        fmt.Printf("Multiple anomalies detected for transaction ID %s. Escalating response.\n", report.TransactionID)
        protocol.escalateTransactionAnomalyResponse(report)
    } else {
        fmt.Printf("Issuing alert for anomaly in transaction ID %s.\n", report.TransactionID)
        protocol.alertForTransactionAnomaly(report)
    }
}

// alertForTransactionAnomaly issues an alert regarding a detected transaction anomaly
func (protocol *TransactionAnomalyDetectionProtocol) alertForTransactionAnomaly(report common.TransactionReport) {
    encryptedAlertData := protocol.encryptTransactionData(report)

    // Issue an alert through the Synnergy Consensus system
    alertSuccess := protocol.consensusSystem.IssueTransactionAnomalyAlert(report.TransactionID, encryptedAlertData)

    if alertSuccess {
        fmt.Printf("Anomaly alert issued for transaction ID %s.\n", report.TransactionID)
        protocol.logTransactionEvent(report, "Alert Issued")
        protocol.resetTransactionRetry(report.TransactionID)
    } else {
        fmt.Printf("Error issuing anomaly alert for transaction ID %s. Retrying...\n", report.TransactionID)
        protocol.retryTransactionResponse(report)
    }
}

// escalateTransactionAnomalyResponse escalates the response to a detected transaction anomaly
func (protocol *TransactionAnomalyDetectionProtocol) escalateTransactionAnomalyResponse(report common.TransactionReport) {
    encryptedEscalationData := protocol.encryptTransactionData(report)

    // Attempt to escalate the response through the Synnergy Consensus system
    escalationSuccess := protocol.consensusSystem.EscalateTransactionAnomalyResponse(report.TransactionID, encryptedEscalationData)

    if escalationSuccess {
        fmt.Printf("Anomaly response escalated for transaction ID %s.\n", report.TransactionID)
        protocol.logTransactionEvent(report, "Response Escalated")
        protocol.resetTransactionRetry(report.TransactionID)
    } else {
        fmt.Printf("Error escalating anomaly response for transaction ID %s. Retrying...\n", report.TransactionID)
        protocol.retryTransactionResponse(report)
    }
}

// retryTransactionResponse retries the response to a transaction anomaly if the initial action fails
func (protocol *TransactionAnomalyDetectionProtocol) retryTransactionResponse(report common.TransactionReport) {
    protocol.transactionRetryCount[report.TransactionID]++
    if protocol.transactionRetryCount[report.TransactionID] < MaxTransactionRetries {
        protocol.escalateTransactionAnomalyResponse(report)
    } else {
        fmt.Printf("Max retries reached for transaction anomaly response for transaction ID %s. Response failed.\n", report.TransactionID)
        protocol.logTransactionFailure(report)
    }
}

// resetTransactionRetry resets the retry count for a specific transaction ID
func (protocol *TransactionAnomalyDetectionProtocol) resetTransactionRetry(transactionID string) {
    protocol.transactionRetryCount[transactionID] = 0
}

// finalizeTransactionMonitoringCycle finalizes the transaction anomaly monitoring cycle and logs the result in the ledger
func (protocol *TransactionAnomalyDetectionProtocol) finalizeTransactionMonitoringCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeTransactionMonitoringCycle()
    if success {
        fmt.Println("Transaction anomaly monitoring cycle finalized successfully.")
        protocol.logTransactionMonitoringCycleFinalization()
    } else {
        fmt.Println("Error finalizing transaction anomaly monitoring cycle.")
    }
}

// logTransactionEvent logs a transaction anomaly event into the ledger
func (protocol *TransactionAnomalyDetectionProtocol) logTransactionEvent(report common.TransactionReport, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("transaction-event-%s-%s", report.TransactionID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Transaction Anomaly Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Transaction %s triggered %s due to anomaly detection.", report.TransactionID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with transaction anomaly event for transaction ID %s.\n", report.TransactionID)
}

// logTransactionFailure logs a failure in responding to a transaction anomaly into the ledger
func (protocol *TransactionAnomalyDetectionProtocol) logTransactionFailure(report common.TransactionReport) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("transaction-failure-%s", report.TransactionID),
        Timestamp: time.Now().Unix(),
        Type:      "Transaction Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to respond to transaction anomaly for transaction ID %s after maximum retries.", report.TransactionID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with transaction failure for transaction ID %s.\n", report.TransactionID)
}

// logTransactionMonitoringCycleFinalization logs the finalization of a transaction monitoring cycle into the ledger
func (protocol *TransactionAnomalyDetectionProtocol) logTransactionMonitoringCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("transaction-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Transaction Monitoring Cycle Finalization",
        Status:    "Finalized",
        Details:   "Transaction anomaly monitoring cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with transaction monitoring cycle finalization.")
}

// encryptTransactionData encrypts transaction-related data before taking action or logging events
func (protocol *TransactionAnomalyDetectionProtocol) encryptTransactionData(report common.TransactionReport) common.TransactionReport {
    encryptedData, err := encryption.EncryptData(report.TransactionData)
    if err != nil {
        fmt.Println("Error encrypting transaction data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("Transaction data successfully encrypted for transaction ID:", report.TransactionID)
    return report
}
