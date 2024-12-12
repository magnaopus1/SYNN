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
    SanctionListMonitoringInterval = 5 * time.Second  // Interval for monitoring transaction sanctions
    MaxSanctionListRetries         = 3                // Maximum retries for responding to sanction list violations
    SubBlocksPerBlock              = 1000             // Number of sub-blocks in a block
)

// TransactionSanctionListProtocol checks transactions against a sanctioned entities list
type TransactionSanctionListProtocol struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger for logging sanction-related events
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    sanctionListRetryCount map[string]int               // Counter for retrying sanction list responses
    sanctionCycleCount     int                          // Counter for monitoring cycles
    sanctionViolationCounter map[string]int             // Tracks detected sanction list violations
}

// NewTransactionSanctionListProtocol initializes the sanction list protocol
func NewTransactionSanctionListProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *TransactionSanctionListProtocol {
    return &TransactionSanctionListProtocol{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        sanctionListRetryCount: make(map[string]int),
        sanctionViolationCounter: make(map[string]int),
        sanctionCycleCount:     0,
    }
}

// StartSanctionListMonitoring starts the continuous loop for monitoring sanction list violations
func (protocol *TransactionSanctionListProtocol) StartSanctionListMonitoring() {
    ticker := time.NewTicker(SanctionListMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorSanctionListViolations()
        }
    }()
}

// monitorSanctionListViolations checks for transactions involving sanctioned entities and takes action
func (protocol *TransactionSanctionListProtocol) monitorSanctionListViolations() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch transaction reports from the consensus system
    transactionReports := protocol.consensusSystem.FetchTransactionReports()

    for _, report := range transactionReports {
        if protocol.isSanctionedEntityInvolved(report) {
            fmt.Printf("Sanctioned entity detected in transaction ID %s. Taking action.\n", report.TransactionID)
            protocol.handleSanctionViolation(report)
        } else {
            fmt.Printf("No sanction list violation detected for transaction ID %s.\n", report.TransactionID)
        }
    }

    protocol.sanctionCycleCount++
    fmt.Printf("Sanction list monitoring cycle #%d completed.\n", protocol.sanctionCycleCount)

    if protocol.sanctionCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeSanctionMonitoringCycle()
    }
}

// isSanctionedEntityInvolved checks if any of the entities in the transaction are on the sanctioned list
func (protocol *TransactionSanctionListProtocol) isSanctionedEntityInvolved(report common.TransactionReport) bool {
    for _, addr := range report.InvolvedAddresses {
        if protocol.consensusSystem.IsSanctionedAddress(addr) {
            fmt.Printf("Sanctioned address %s detected in transaction ID %s.\n", addr, report.TransactionID)
            return true
        }
    }
    return false
}

// handleSanctionViolation handles the response to a sanction list violation
func (protocol *TransactionSanctionListProtocol) handleSanctionViolation(report common.TransactionReport) {
    protocol.sanctionViolationCounter[report.TransactionID]++

    if protocol.sanctionViolationCounter[report.TransactionID] >= MaxSanctionListRetries {
        fmt.Printf("Multiple sanction violations detected for transaction ID %s. Escalating response.\n", report.TransactionID)
        protocol.escalateSanctionViolation(report)
    } else {
        fmt.Printf("Issuing alert for sanction list violation in transaction ID %s.\n", report.TransactionID)
        protocol.alertForSanctionViolation(report)
    }
}

// alertForSanctionViolation issues an alert for a detected sanction list violation
func (protocol *TransactionSanctionListProtocol) alertForSanctionViolation(report common.TransactionReport) {
    encryptedAlertData := protocol.encryptSanctionData(report)

    // Issue an alert through the Synnergy Consensus system
    alertSuccess := protocol.consensusSystem.IssueSanctionViolationAlert(report.TransactionID, encryptedAlertData)

    if alertSuccess {
        fmt.Printf("Sanction violation alert issued for transaction ID %s.\n", report.TransactionID)
        protocol.logSanctionEvent(report, "Alert Issued")
        protocol.resetSanctionRetry(report.TransactionID)
    } else {
        fmt.Printf("Error issuing sanction violation alert for transaction ID %s. Retrying...\n", report.TransactionID)
        protocol.retrySanctionResponse(report)
    }
}

// escalateSanctionViolation escalates the response to a detected sanction list violation
func (protocol *TransactionSanctionListProtocol) escalateSanctionViolation(report common.TransactionReport) {
    encryptedEscalationData := protocol.encryptSanctionData(report)

    // Attempt to escalate the sanction violation response through the Synnergy Consensus system
    escalationSuccess := protocol.consensusSystem.EscalateSanctionViolationResponse(report.TransactionID, encryptedEscalationData)

    if escalationSuccess {
        fmt.Printf("Sanction violation response escalated for transaction ID %s.\n", report.TransactionID)
        protocol.logSanctionEvent(report, "Response Escalated")
        protocol.resetSanctionRetry(report.TransactionID)
    } else {
        fmt.Printf("Error escalating sanction violation response for transaction ID %s. Retrying...\n", report.TransactionID)
        protocol.retrySanctionResponse(report)
    }
}

// retrySanctionResponse retries the response to a sanction list violation if the initial action fails
func (protocol *TransactionSanctionListProtocol) retrySanctionResponse(report common.TransactionReport) {
    protocol.sanctionListRetryCount[report.TransactionID]++
    if protocol.sanctionListRetryCount[report.TransactionID] < MaxSanctionListRetries {
        protocol.escalateSanctionViolation(report)
    } else {
        fmt.Printf("Max retries reached for sanction violation response for transaction ID %s. Response failed.\n", report.TransactionID)
        protocol.logSanctionFailure(report)
    }
}

// resetSanctionRetry resets the retry count for sanction violations for a specific transaction ID
func (protocol *TransactionSanctionListProtocol) resetSanctionRetry(transactionID string) {
    protocol.sanctionListRetryCount[transactionID] = 0
}

// finalizeSanctionMonitoringCycle finalizes the sanction list monitoring cycle and logs the result in the ledger
func (protocol *TransactionSanctionListProtocol) finalizeSanctionMonitoringCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeSanctionMonitoringCycle()
    if success {
        fmt.Println("Sanction list monitoring cycle finalized successfully.")
        protocol.logSanctionMonitoringCycleFinalization()
    } else {
        fmt.Println("Error finalizing sanction list monitoring cycle.")
    }
}

// logSanctionEvent logs a sanction-related event into the ledger
func (protocol *TransactionSanctionListProtocol) logSanctionEvent(report common.TransactionReport, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sanction-event-%s-%s", report.TransactionID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Sanction List Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Transaction %s triggered %s due to sanction list violation.", report.TransactionID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with sanction violation event for transaction ID %s.\n", report.TransactionID)
}

// logSanctionFailure logs the failure to respond to a sanction list violation into the ledger
func (protocol *TransactionSanctionListProtocol) logSanctionFailure(report common.TransactionReport) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sanction-failure-%s", report.TransactionID),
        Timestamp: time.Now().Unix(),
        Type:      "Sanction List Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to respond to sanction list violation for transaction ID %s after maximum retries.", report.TransactionID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with sanction failure for transaction ID %s.\n", report.TransactionID)
}

// logSanctionMonitoringCycleFinalization logs the finalization of a sanction monitoring cycle into the ledger
func (protocol *TransactionSanctionListProtocol) logSanctionMonitoringCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sanction-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Sanction Monitoring Cycle Finalization",
        Status:    "Finalized",
        Details:   "Sanction list monitoring cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with sanction monitoring cycle finalization.")
}

// encryptSanctionData encrypts sanction-related data before taking action or logging events
func (protocol *TransactionSanctionListProtocol) encryptSanctionData(report common.TransactionReport) common.TransactionReport {
    encryptedData, err := encryption.EncryptData(report.SanctionData)
    if err != nil {
        fmt.Println("Error encrypting sanction data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("Sanction data successfully encrypted for transaction ID:", report.TransactionID)
    return report
}
