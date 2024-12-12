package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/transaction"
)

const (
    FeeMonitoringInterval = 5 * time.Second  // Interval for monitoring transaction fees
    MaxFeeRetries         = 3                // Maximum retries for enforcing fee ceilings/floors
    SubBlocksPerBlock     = 1000             // Number of sub-blocks in a block
)

// TransactionFeeCeilingProtocol manages the enforcement of fee ceilings and floors in transactions
type TransactionFeeCeilingProtocol struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger for logging fee ceiling/floor events
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    feeRetryCount        map[string]int               // Counter for retrying fee enforcement actions
    feeCycleCount        int                          // Counter for fee monitoring cycles
    feeAnomalyCounter    map[string]int               // Tracks detected fee anomalies
}

// NewTransactionFeeCeilingProtocol initializes the fee ceiling/floor enforcement protocol
func NewTransactionFeeCeilingProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *TransactionFeeCeilingProtocol {
    return &TransactionFeeCeilingProtocol{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        feeRetryCount:     make(map[string]int),
        feeAnomalyCounter: make(map[string]int),
        feeCycleCount:     0,
    }
}

// StartFeeMonitoring starts the continuous loop for monitoring and enforcing transaction fees
func (protocol *TransactionFeeCeilingProtocol) StartFeeMonitoring() {
    ticker := time.NewTicker(FeeMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorTransactionFees()
        }
    }()
}

// monitorTransactionFees checks for transactions that exceed the defined fee ceiling or fall below the fee floor
func (protocol *TransactionFeeCeilingProtocol) monitorTransactionFees() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch transaction fee reports from the transaction package
    feeReports := transaction.FetchTransactionFeeReports()

    for _, report := range feeReports {
        if protocol.isFeeAnomalyDetected(report) {
            fmt.Printf("Fee anomaly detected for transaction ID %s. Taking action.\n", report.TransactionID)
            protocol.handleFeeAnomaly(report)
        } else {
            fmt.Printf("No fee anomaly detected for transaction ID %s.\n", report.TransactionID)
        }
    }

    protocol.feeCycleCount++
    fmt.Printf("Transaction fee monitoring cycle #%d completed.\n", protocol.feeCycleCount)

    if protocol.feeCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeFeeMonitoringCycle()
    }
}

// isFeeAnomalyDetected checks if the transaction fee exceeds the ceiling or falls below the floor
func (protocol *TransactionFeeCeilingProtocol) isFeeAnomalyDetected(report common.FeeReport) bool {
    // Detect if the transaction fee exceeds the ceiling
    if report.TransactionFee > transaction.EnforceFeeCeiling {
        fmt.Printf("Transaction fee exceeds ceiling for transaction ID: %s. Fee: %f\n", report.TransactionID, report.TransactionFee)
        return true
    }

    // Detect if the transaction fee falls below the floor
    if report.TransactionFee < transaction.EnforceFeeFloor {
        fmt.Printf("Transaction fee falls below floor for transaction ID: %s. Fee: %f\n", report.TransactionID, report.TransactionFee)
        return true
    }

    // If none of the above conditions match, no anomaly is detected
    fmt.Printf("No fee anomaly detected for transaction ID: %s\n", report.TransactionID)
    return false
}

// handleFeeAnomaly takes action when a transaction fee anomaly is detected
func (protocol *TransactionFeeCeilingProtocol) handleFeeAnomaly(report common.FeeReport) {
    protocol.feeAnomalyCounter[report.TransactionID]++

    if protocol.feeAnomalyCounter[report.TransactionID] >= MaxFeeRetries {
        fmt.Printf("Multiple fee anomalies detected for transaction ID %s. Escalating response.\n", report.TransactionID)
        protocol.escalateFeeAnomalyResponse(report)
    } else {
        fmt.Printf("Issuing alert for fee anomaly in transaction ID %s.\n", report.TransactionID)
        protocol.alertForFeeAnomaly(report)
    }
}

// alertForFeeAnomaly issues an alert regarding a detected fee anomaly
func (protocol *TransactionFeeCeilingProtocol) alertForFeeAnomaly(report common.FeeReport) {
    encryptedAlertData := protocol.encryptFeeData(report)

    // Issue an alert through the Synnergy Consensus system
    alertSuccess := protocol.consensusSystem.IssueFeeAnomalyAlert(report.TransactionID, encryptedAlertData)

    if alertSuccess {
        fmt.Printf("Fee anomaly alert issued for transaction ID %s.\n", report.TransactionID)
        protocol.logFeeEvent(report, "Alert Issued")
        protocol.resetFeeRetry(report.TransactionID)
    } else {
        fmt.Printf("Error issuing fee anomaly alert for transaction ID %s. Retrying...\n", report.TransactionID)
        protocol.retryFeeResponse(report)
    }
}

// escalateFeeAnomalyResponse escalates the response to a detected fee anomaly
func (protocol *TransactionFeeCeilingProtocol) escalateFeeAnomalyResponse(report common.FeeReport) {
    encryptedEscalationData := protocol.encryptFeeData(report)

    // Attempt to escalate the fee anomaly response through the Synnergy Consensus system
    escalationSuccess := protocol.consensusSystem.EscalateFeeAnomalyResponse(report.TransactionID, encryptedEscalationData)

    if escalationSuccess {
        fmt.Printf("Fee anomaly response escalated for transaction ID %s.\n", report.TransactionID)
        protocol.logFeeEvent(report, "Response Escalated")
        protocol.resetFeeRetry(report.TransactionID)
    } else {
        fmt.Printf("Error escalating fee anomaly response for transaction ID %s. Retrying...\n", report.TransactionID)
        protocol.retryFeeResponse(report)
    }
}

// retryFeeResponse retries the response to a fee anomaly if the initial action fails
func (protocol *TransactionFeeCeilingProtocol) retryFeeResponse(report common.FeeReport) {
    protocol.feeRetryCount[report.TransactionID]++
    if protocol.feeRetryCount[report.TransactionID] < MaxFeeRetries {
        protocol.escalateFeeAnomalyResponse(report)
    } else {
        fmt.Printf("Max retries reached for fee anomaly response for transaction ID %s. Response failed.\n", report.TransactionID)
        protocol.logFeeFailure(report)
    }
}

// resetFeeRetry resets the retry count for fee anomaly responses on a specific transaction ID
func (protocol *TransactionFeeCeilingProtocol) resetFeeRetry(transactionID string) {
    protocol.feeRetryCount[transactionID] = 0
}

// finalizeFeeMonitoringCycle finalizes the transaction fee monitoring cycle and logs the result in the ledger
func (protocol *TransactionFeeCeilingProtocol) finalizeFeeMonitoringCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeFeeMonitoringCycle()
    if success {
        fmt.Println("Transaction fee monitoring cycle finalized successfully.")
        protocol.logFeeMonitoringCycleFinalization()
    } else {
        fmt.Println("Error finalizing transaction fee monitoring cycle.")
    }
}

// logFeeEvent logs a fee-related event into the ledger
func (protocol *TransactionFeeCeilingProtocol) logFeeEvent(report common.FeeReport, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("fee-event-%s-%s", report.TransactionID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Transaction Fee Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Transaction %s triggered %s due to fee anomaly detection.", report.TransactionID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with fee anomaly event for transaction ID %s.\n", report.TransactionID)
}

// logFeeFailure logs the failure to respond to a fee anomaly into the ledger
func (protocol *TransactionFeeCeilingProtocol) logFeeFailure(report common.FeeReport) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("fee-failure-%s", report.TransactionID),
        Timestamp: time.Now().Unix(),
        Type:      "Fee Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to respond to fee anomaly for transaction ID %s after maximum retries.", report.TransactionID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with fee failure for transaction ID %s.\n", report.TransactionID)
}

// logFeeMonitoringCycleFinalization logs the finalization of a fee monitoring cycle into the ledger
func (protocol *TransactionFeeCeilingProtocol) logFeeMonitoringCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("fee-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Fee Monitoring Cycle Finalization",
        Status:    "Finalized",
        Details:   "Transaction fee monitoring cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with fee monitoring cycle finalization.")
}

// encryptFeeData encrypts fee-related data before taking action or logging events
func (protocol *TransactionFeeCeilingProtocol) encryptFeeData(report common.FeeReport) common.FeeReport {
    encryptedData, err := encryption.EncryptData(report.FeeData)
    if err != nil {
        fmt.Println("Error encrypting fee data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("Fee data successfully encrypted for transaction ID:", report.TransactionID)
    return report
}
