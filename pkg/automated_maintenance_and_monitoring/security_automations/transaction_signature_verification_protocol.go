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
    SignatureVerificationInterval  = 5 * time.Second // Interval for checking transaction signatures
    MaxSignatureRetries            = 3               // Maximum retries for responding to signature verification failures
    SubBlocksPerBlock              = 1000            // Number of sub-blocks in a block
)

// TransactionSignatureVerificationProtocol handles the verification of transaction signatures
type TransactionSignatureVerificationProtocol struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger for logging signature verification events
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    signatureRetryCount  map[string]int               // Counter for retrying signature verification
    signatureCycleCount  int                          // Counter for monitoring cycles
}

// NewTransactionSignatureVerificationProtocol initializes the signature verification protocol
func NewTransactionSignatureVerificationProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *TransactionSignatureVerificationProtocol {
    return &TransactionSignatureVerificationProtocol{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        signatureRetryCount: make(map[string]int),
        signatureCycleCount: 0,
    }
}

// StartSignatureVerification starts the continuous loop for monitoring transaction signatures
func (protocol *TransactionSignatureVerificationProtocol) StartSignatureVerification() {
    ticker := time.NewTicker(SignatureVerificationInterval)

    go func() {
        for range ticker.C {
            protocol.verifyTransactionSignatures()
        }
    }()
}

// verifyTransactionSignatures checks the validity of transaction signatures and takes action if invalid
func (protocol *TransactionSignatureVerificationProtocol) verifyTransactionSignatures() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch the list of transactions from the consensus system
    transactionReports := protocol.consensusSystem.FetchTransactionReports()

    for _, report := range transactionReports {
        if !protocol.isSignatureValid(report) {
            fmt.Printf("Invalid signature detected for transaction ID %s. Taking action.\n", report.TransactionID)
            protocol.handleInvalidSignature(report)
        } else {
            fmt.Printf("Valid signature for transaction ID %s.\n", report.TransactionID)
        }
    }

    protocol.signatureCycleCount++
    fmt.Printf("Signature verification cycle #%d completed.\n", protocol.signatureCycleCount)

    if protocol.signatureCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeSignatureVerificationCycle()
    }
}

// isSignatureValid verifies if the transaction signature is valid
func (protocol *TransactionSignatureVerificationProtocol) isSignatureValid(report common.TransactionReport) bool {
    // Use the consensus system to verify the transaction signature
    validSignature := protocol.consensusSystem.VerifyTransactionSignature(report.TransactionID, report.Signature)

    if !validSignature {
        fmt.Printf("Signature verification failed for transaction ID: %s.\n", report.TransactionID)
        return false
    }
    return true
}

// handleInvalidSignature takes action when an invalid signature is detected
func (protocol *TransactionSignatureVerificationProtocol) handleInvalidSignature(report common.TransactionReport) {
    protocol.signatureRetryCount[report.TransactionID]++

    if protocol.signatureRetryCount[report.TransactionID] >= MaxSignatureRetries {
        fmt.Printf("Multiple invalid signatures detected for transaction ID %s. Escalating response.\n", report.TransactionID)
        protocol.escalateInvalidSignatureResponse(report)
    } else {
        fmt.Printf("Issuing alert for invalid signature in transaction ID %s.\n", report.TransactionID)
        protocol.alertForInvalidSignature(report)
    }
}

// alertForInvalidSignature issues an alert for an invalid transaction signature
func (protocol *TransactionSignatureVerificationProtocol) alertForInvalidSignature(report common.TransactionReport) {
    encryptedAlertData := protocol.encryptSignatureData(report)

    // Issue an alert through the Synnergy Consensus system
    alertSuccess := protocol.consensusSystem.IssueInvalidSignatureAlert(report.TransactionID, encryptedAlertData)

    if alertSuccess {
        fmt.Printf("Invalid signature alert issued for transaction ID %s.\n", report.TransactionID)
        protocol.logSignatureEvent(report, "Alert Issued")
        protocol.resetSignatureRetry(report.TransactionID)
    } else {
        fmt.Printf("Error issuing invalid signature alert for transaction ID %s. Retrying...\n", report.TransactionID)
        protocol.retrySignatureResponse(report)
    }
}

// escalateInvalidSignatureResponse escalates the response to a persistent invalid signature
func (protocol *TransactionSignatureVerificationProtocol) escalateInvalidSignatureResponse(report common.TransactionReport) {
    encryptedEscalationData := protocol.encryptSignatureData(report)

    // Attempt to escalate the invalid signature response through the Synnergy Consensus system
    escalationSuccess := protocol.consensusSystem.EscalateInvalidSignatureResponse(report.TransactionID, encryptedEscalationData)

    if escalationSuccess {
        fmt.Printf("Invalid signature response escalated for transaction ID %s.\n", report.TransactionID)
        protocol.logSignatureEvent(report, "Response Escalated")
        protocol.resetSignatureRetry(report.TransactionID)
    } else {
        fmt.Printf("Error escalating invalid signature response for transaction ID %s. Retrying...\n", report.TransactionID)
        protocol.retrySignatureResponse(report)
    }
}

// retrySignatureResponse retries the response to an invalid signature if the initial action fails
func (protocol *TransactionSignatureVerificationProtocol) retrySignatureResponse(report common.TransactionReport) {
    protocol.signatureRetryCount[report.TransactionID]++
    if protocol.signatureRetryCount[report.TransactionID] < MaxSignatureRetries {
        protocol.escalateInvalidSignatureResponse(report)
    } else {
        fmt.Printf("Max retries reached for invalid signature response for transaction ID %s. Response failed.\n", report.TransactionID)
        protocol.logSignatureFailure(report)
    }
}

// resetSignatureRetry resets the retry count for invalid signatures for a specific transaction ID
func (protocol *TransactionSignatureVerificationProtocol) resetSignatureRetry(transactionID string) {
    protocol.signatureRetryCount[transactionID] = 0
}

// finalizeSignatureVerificationCycle finalizes the signature verification cycle and logs the result in the ledger
func (protocol *TransactionSignatureVerificationProtocol) finalizeSignatureVerificationCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeSignatureVerificationCycle()
    if success {
        fmt.Println("Signature verification cycle finalized successfully.")
        protocol.logSignatureVerificationCycleFinalization()
    } else {
        fmt.Println("Error finalizing signature verification cycle.")
    }
}

// logSignatureEvent logs a signature-related event into the ledger
func (protocol *TransactionSignatureVerificationProtocol) logSignatureEvent(report common.TransactionReport, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("signature-event-%s-%s", report.TransactionID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Signature Verification Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Transaction %s triggered %s due to signature verification failure.", report.TransactionID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with signature verification event for transaction ID %s.\n", report.TransactionID)
}

// logSignatureFailure logs the failure to respond to an invalid signature into the ledger
func (protocol *TransactionSignatureVerificationProtocol) logSignatureFailure(report common.TransactionReport) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("signature-failure-%s", report.TransactionID),
        Timestamp: time.Now().Unix(),
        Type:      "Signature Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to respond to invalid signature for transaction ID %s after maximum retries.", report.TransactionID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with signature failure for transaction ID %s.\n", report.TransactionID)
}

// logSignatureVerificationCycleFinalization logs the finalization of a signature verification cycle into the ledger
func (protocol *TransactionSignatureVerificationProtocol) logSignatureVerificationCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("signature-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Signature Verification Cycle Finalization",
        Status:    "Finalized",
        Details:   "Transaction signature verification cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with signature verification cycle finalization.")
}

// encryptSignatureData encrypts signature-related data before taking action or logging events
func (protocol *TransactionSignatureVerificationProtocol) encryptSignatureData(report common.TransactionReport) common.TransactionReport {
    encryptedData, err := encryption.EncryptData(report.SignatureData)
    if err != nil {
        fmt.Println("Error encrypting signature data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("Signature data successfully encrypted for transaction ID:", report.TransactionID)
    return report
}
