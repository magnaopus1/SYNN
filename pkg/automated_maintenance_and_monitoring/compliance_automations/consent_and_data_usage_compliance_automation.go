package compliance_automations

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
)

const (
    ConsentComplianceCheckInterval = 30 * time.Minute // Interval for checking consent-based data usage compliance
)

// ConsentAndDataUsageComplianceAutomation manages compliance with user consent agreements and data usage policies
type ConsentAndDataUsageComplianceAutomation struct {
    ledgerInstance *ledger.Ledger // Blockchain ledger instance
    stateMutex     *sync.RWMutex  // Mutex for thread-safe ledger access
    apiURL         string         // API URL for compliance and data protection
}

// NewConsentAndDataUsageComplianceAutomation initializes the automation handler
func NewConsentAndDataUsageComplianceAutomation(apiURL string, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ConsentAndDataUsageComplianceAutomation {
    return &ConsentAndDataUsageComplianceAutomation{
        ledgerInstance: ledgerInstance,
        stateMutex:     stateMutex,
        apiURL:         apiURL,
    }
}

// StartConsentComplianceMonitoring starts continuous monitoring of consent-based data usage
func (automation *ConsentAndDataUsageComplianceAutomation) StartConsentComplianceMonitoring() {
    ticker := time.NewTicker(ConsentComplianceCheckInterval)
    for range ticker.C {
        fmt.Println("Starting user consent compliance monitoring...")
        automation.monitorDataUsageCompliance()
    }
}

// monitorDataUsageCompliance retrieves and validates data usage against user consent agreements
func (automation *ConsentAndDataUsageComplianceAutomation) monitorDataUsageCompliance() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    // Retrieve validated sub-blocks from the ledger containing transaction data
    subBlocks := automation.ledgerInstance.GetValidatedSubBlocks()
    for _, subBlock := range subBlocks {
        for _, tx := range subBlock.Transactions {
            log := automation.extractDataUsageLog(tx)
            if !automation.isCompliantWithConsent(log) {
                fmt.Printf("Data usage violation detected for transaction ID: %s\n", tx.ID)
                automation.triggerDataViolationResponse(log)
            }
        }
    }
}

// extractDataUsageLog extracts data usage log from a validated transaction
func (automation *ConsentAndDataUsageComplianceAutomation) extractDataUsageLog(tx common.Transaction) common.DataUsageLog {
    var log common.DataUsageLog
    err := json.Unmarshal([]byte(tx.Data), &log)
    if err != nil {
        fmt.Printf("Error extracting data usage log from transaction ID %s: %v\n", tx.ID, err)
    }
    return log
}

// isCompliantWithConsent checks if the data usage complies with user consent conditions
func (automation *ConsentAndDataUsageComplianceAutomation) isCompliantWithConsent(log common.DataUsageLog) bool {
    // Retrieve user consent from the ledger for the specific user ID
    consentRecord := automation.ledgerInstance.GetConsentRecord(log.UserID)

    if consentRecord == nil {
        fmt.Printf("No consent record found for user ID: %s\n", log.UserID)
        return false
    }

    // Check if data usage complies with consent agreements
    if log.DataUsageType != consentRecord.AllowedDataUsageType {
        fmt.Printf("Data usage type %s violates user consent for log ID: %s\n", log.DataUsageType, log.ID)
        return false
    }

    if log.SharedWith != "" && !consentRecord.AllowedDataSharing {
        fmt.Printf("Data sharing not allowed for log ID: %s\n", log.ID)
        return false
    }

    return true
}

// triggerDataViolationResponse triggers an automatic response for data usage violations
func (automation *ConsentAndDataUsageComplianceAutomation) triggerDataViolationResponse(log common.DataUsageLog) {
    fmt.Printf("Triggering data violation response for log ID: %s\n", log.ID)

    // Log the violation in the audit trail within the ledger
    automation.logDataUsageViolation(log)

    // Optionally notify the user about the violation
    automation.notifyUserOfViolation(log.UserID, log.ID)
}

// logDataUsageViolation logs the data usage violation in the blockchain ledger audit trail
func (automation *ConsentAndDataUsageComplianceAutomation) logDataUsageViolation(log common.DataUsageLog) {
    auditEntry := common.LedgerEntry{
        ID:        log.ID,
        UserID:    log.UserID,
        Event:     "Data usage violation",
        Timestamp: time.Now().Unix(),
        Data:      fmt.Sprintf("Data usage type %s violated user consent", log.DataUsageType),
    }

    // Encrypt the audit entry before storing it in the ledger
    encryptedAuditEntry, err := encryption.EncryptLedgerEntry(auditEntry, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting audit entry for log ID %s: %v\n", log.ID, err)
        return
    }

    // Add the encrypted entry to the ledger using Synnergy Consensus
    automation.ledgerInstance.AddEntry(encryptedAuditEntry)
    fmt.Printf("Data usage violation logged in ledger for log ID: %s\n", log.ID)
}

// notifyUserOfViolation notifies the user about the data usage violation
func (automation *ConsentAndDataUsageComplianceAutomation) notifyUserOfViolation(userID, logID string) {
    // Integrate with external services to notify the user (email, SMS, etc.)
    fmt.Printf("Notifying user ID %s of data usage violation for log ID %s\n", userID, logID)
}

// StartConsentAgreementValidation continuously validates consent agreements and ensures compliance
func (automation *ConsentAndDataUsageComplianceAutomation) StartConsentAgreementValidation() {
    ticker := time.NewTicker(ConsentComplianceCheckInterval)
    for range ticker.C {
        fmt.Println("Validating user consent agreements...")
        automation.validateConsentAgreements()
    }
}

// validateConsentAgreements checks all consent records in the ledger and ensures compliance
func (automation *ConsentAndDataUsageComplianceAutomation) validateConsentAgreements() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    consentRecords := automation.ledgerInstance.GetAllConsentRecords()
    for _, record := range consentRecords {
        if record.Expiry.Before(time.Now()) {
            fmt.Printf("Consent for user ID %s has expired.\n", record.UserID)
            automation.revokeConsent(record.UserID)
        }
    }
}

// revokeConsent revokes the user's consent if it has expired or been violated
func (automation *ConsentAndDataUsageComplianceAutomation) revokeConsent(userID string) {
    fmt.Printf("Revoking consent for user ID %s\n", userID)

    // Log consent revocation in the audit trail in the ledger
    automation.logConsentRevocation(userID)

    // Apply data restriction after consent revocation
    automation.applyDataRestriction(userID)
}

// logConsentRevocation logs the consent revocation in the blockchain ledger
func (automation *ConsentAndDataUsageComplianceAutomation) logConsentRevocation(userID string) {
    auditEntry := common.LedgerEntry{
        ID:        fmt.Sprintf("consent-revocation-%s", userID),
        UserID:    userID,
        Event:     "Consent revoked",
        Timestamp: time.Now().Unix(),
        Data:      "User consent was revoked due to expiration or violation",
    }

    encryptedAuditEntry, err := encryption.EncryptLedgerEntry(auditEntry, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting consent revocation for user ID %s: %v\n", userID, err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedAuditEntry)
    fmt.Printf("Consent revocation logged for user ID %s\n", userID)
}

// applyDataRestriction applies data restrictions to the user's data after consent revocation
func (automation *ConsentAndDataUsageComplianceAutomation) applyDataRestriction(userID string) {
    url := fmt.Sprintf("%s/api/compliance/data_protection/apply", automation.apiURL)
    body, _ := json.Marshal(map[string]string{
        "user_id":     userID,
        "restriction": "No data usage allowed after consent revocation",
    })

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error applying data restrictions for user ID %s: %v\n", userID, err)
    } else {
        fmt.Printf("Data usage restriction applied for user ID %s.\n", userID)
    }
}
