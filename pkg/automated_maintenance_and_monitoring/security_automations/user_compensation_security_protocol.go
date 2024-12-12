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
    CompensationMonitoringInterval = 5 * time.Second  // Interval for monitoring compensation events
    MaxCompensationRetries         = 3                // Maximum retries for issuing compensation
    SubBlocksPerBlock              = 1000             // Number of sub-blocks in a block
)

// UserCompensationSecurityProtocol automates the process of compensating users based on security events
type UserCompensationSecurityProtocol struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger for logging compensation-related events
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    compensationRetryCount map[string]int             // Counter for retrying compensation
    compensationCycleCount int                        // Counter for compensation monitoring cycles
}

// NewUserCompensationSecurityProtocol initializes the compensation protocol
func NewUserCompensationSecurityProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *UserCompensationSecurityProtocol {
    return &UserCompensationSecurityProtocol{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        compensationRetryCount: make(map[string]int),
        compensationCycleCount: 0,
    }
}

// StartCompensationMonitoring starts the continuous loop for monitoring compensation events
func (protocol *UserCompensationSecurityProtocol) StartCompensationMonitoring() {
    ticker := time.NewTicker(CompensationMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorCompensationEvents()
        }
    }()
}

// monitorCompensationEvents checks for events that trigger user compensation and takes action
func (protocol *UserCompensationSecurityProtocol) monitorCompensationEvents() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch the list of compensatable security events from the consensus system
    compensationReports := protocol.consensusSystem.FetchCompensationReports()

    for _, report := range compensationReports {
        if protocol.isCompensationRequired(report) {
            fmt.Printf("Compensation required for user %s due to event %s. Taking action.\n", report.UserID, report.EventID)
            protocol.handleCompensationEvent(report)
        } else {
            fmt.Printf("No compensation required for user %s.\n", report.UserID)
        }
    }

    protocol.compensationCycleCount++
    fmt.Printf("User compensation monitoring cycle #%d completed.\n", protocol.compensationCycleCount)

    if protocol.compensationCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeCompensationCycle()
    }
}

// isCompensationRequired checks if compensation is required based on the security event
func (protocol *UserCompensationSecurityProtocol) isCompensationRequired(report common.CompensationReport) bool {
    return report.CompensationAmount > 0
}

// handleCompensationEvent responds to a compensatable event by issuing compensation
func (protocol *UserCompensationSecurityProtocol) handleCompensationEvent(report common.CompensationReport) {
    protocol.compensationRetryCount[report.UserID]++

    if protocol.compensationRetryCount[report.UserID] >= MaxCompensationRetries {
        fmt.Printf("Multiple compensation attempts failed for user %s. Escalating.\n", report.UserID)
        protocol.escalateCompensationEvent(report)
    } else {
        fmt.Printf("Issuing compensation of %f to user %s.\n", report.CompensationAmount, report.UserID)
        protocol.issueCompensation(report)
    }
}

// issueCompensation processes the compensation and logs the event in the ledger
func (protocol *UserCompensationSecurityProtocol) issueCompensation(report common.CompensationReport) {
    encryptedCompensationData := protocol.encryptCompensationData(report)

    // Issue compensation through the Synnergy Consensus system
    compensationSuccess := protocol.consensusSystem.IssueCompensation(report.UserID, encryptedCompensationData)

    if compensationSuccess {
        fmt.Printf("Compensation issued successfully to user %s.\n", report.UserID)
        protocol.logCompensationEvent(report, "Compensation Issued")
        protocol.resetCompensationRetry(report.UserID)
    } else {
        fmt.Printf("Error issuing compensation to user %s. Retrying...\n", report.UserID)
        protocol.retryCompensationResponse(report)
    }
}

// escalateCompensationEvent escalates the response to compensation issues
func (protocol *UserCompensationSecurityProtocol) escalateCompensationEvent(report common.CompensationReport) {
    encryptedEscalationData := protocol.encryptCompensationData(report)

    // Attempt to escalate the compensation response through the Synnergy Consensus system
    escalationSuccess := protocol.consensusSystem.EscalateCompensationResponse(report.UserID, encryptedEscalationData)

    if escalationSuccess {
        fmt.Printf("Compensation event escalated for user %s.\n", report.UserID)
        protocol.logCompensationEvent(report, "Response Escalated")
        protocol.resetCompensationRetry(report.UserID)
    } else {
        fmt.Printf("Error escalating compensation response for user %s. Retrying...\n", report.UserID)
        protocol.retryCompensationResponse(report)
    }
}

// retryCompensationResponse retries the compensation process if the initial action fails
func (protocol *UserCompensationSecurityProtocol) retryCompensationResponse(report common.CompensationReport) {
    protocol.compensationRetryCount[report.UserID]++
    if protocol.compensationRetryCount[report.UserID] < MaxCompensationRetries {
        protocol.escalateCompensationEvent(report)
    } else {
        fmt.Printf("Max retries reached for compensation to user %s. Response failed.\n", report.UserID)
        protocol.logCompensationFailure(report)
    }
}

// resetCompensationRetry resets the retry count for compensation for a specific user
func (protocol *UserCompensationSecurityProtocol) resetCompensationRetry(userID string) {
    protocol.compensationRetryCount[userID] = 0
}

// finalizeCompensationCycle finalizes the compensation cycle and logs the result in the ledger
func (protocol *UserCompensationSecurityProtocol) finalizeCompensationCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeCompensationCycle()
    if success {
        fmt.Println("Compensation cycle finalized successfully.")
        protocol.logCompensationCycleFinalization()
    } else {
        fmt.Println("Error finalizing compensation cycle.")
    }
}

// logCompensationEvent logs a compensation-related event into the ledger
func (protocol *UserCompensationSecurityProtocol) logCompensationEvent(report common.CompensationReport, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compensation-event-%s-%s", report.UserID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "User Compensation Event",
        Status:    eventType,
        Details:   fmt.Sprintf("User %s received compensation due to event %s.", report.UserID, report.EventID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compensation event for user %s.\n", report.UserID)
}

// logCompensationFailure logs the failure to respond to a compensation event into the ledger
func (protocol *UserCompensationSecurityProtocol) logCompensationFailure(report common.CompensationReport) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compensation-failure-%s", report.UserID),
        Timestamp: time.Now().Unix(),
        Type:      "Compensation Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to compensate user %s after maximum retries.", report.UserID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compensation failure for user %s.\n", report.UserID)
}

// logCompensationCycleFinalization logs the finalization of a compensation cycle into the ledger
func (protocol *UserCompensationSecurityProtocol) logCompensationCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compensation-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Compensation Cycle Finalization",
        Status:    "Finalized",
        Details:   "User compensation cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with compensation cycle finalization.")
}

// encryptCompensationData encrypts compensation-related data before taking action or logging events
func (protocol *UserCompensationSecurityProtocol) encryptCompensationData(report common.CompensationReport) common.CompensationReport {
    encryptedData, err := encryption.EncryptData(report.CompensationData)
    if err != nil {
        fmt.Println("Error encrypting compensation data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("Compensation data successfully encrypted for user:", report.UserID)
    return report
}
