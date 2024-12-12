package syn12

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// Syn12Token represents the core structure for Treasury Bill tokens.
type Syn12Token struct {
    TokenID             string
    Metadata            Syn12Metadata
    Issuer              string
    Ledger              *ledger.Ledger
    Consensus           *consensus.SynnergyConsensus
    ComplianceService   *compliance.KYCAmlService
    SecurityEnabled     bool
    ComplianceLogging   bool
    SecurityMeasures    map[string]string
    ComplianceParams    map[string]string
    mutex               sync.Mutex
}

// ENABLE_TBILL_SECURITY_AUDIT enables security auditing for Treasury Bill transactions.
func (token *Syn12Token) ENABLE_TBILL_SECURITY_AUDIT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.SecurityEnabled = true
    return token.Ledger.RecordLog("TBillSecurityAuditEnabled", "Security audit enabled for Treasury Bill transactions")
}

// DISABLE_TBILL_SECURITY_AUDIT disables security auditing for Treasury Bill transactions.
func (token *Syn12Token) DISABLE_TBILL_SECURITY_AUDIT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.SecurityEnabled = false
    return token.Ledger.RecordLog("TBillSecurityAuditDisabled", "Security audit disabled for Treasury Bill transactions")
}

// INTEGRATE_TBILL_SECURITY_MEASURES adds or updates security measures for T-Bill operations.
func (token *Syn12Token) INTEGRATE_TBILL_SECURITY_MEASURES(measure, description string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.SecurityMeasures[measure] = description
    return token.Ledger.RecordLog("TBillSecurityMeasureIntegrated", fmt.Sprintf("Integrated security measure: %s", measure))
}

// REGISTER_TBILL_SECURITY_EVENT logs security-related events with encryption for protection.
func (token *Syn12Token) REGISTER_TBILL_SECURITY_EVENT(event string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    encryptedEvent, err := token.Encryption.Encrypt(event)
    if err != nil {
        return fmt.Errorf("failed to encrypt security event: %v", err)
    }

    return token.Ledger.RecordLog("TBillSecurityEvent", encryptedEvent)
}

// FETCH_TBILL_SECURITY_REPORT retrieves a report of recent security events for T-Bill operations.
func (token *Syn12Token) FETCH_TBILL_SECURITY_REPORT() (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    report, err := token.Ledger.GetSecurityLog()
    if err != nil {
        return "", fmt.Errorf("failed to fetch security report: %v", err)
    }
    return report, nil
}

// UPDATE_TBILL_SECURITY_MEASURES modifies existing security measures for T-Bill operations.
func (token *Syn12Token) UPDATE_TBILL_SECURITY_MEASURES(measure, newDescription string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if _, exists := token.SecurityMeasures[measure]; exists {
        token.SecurityMeasures[measure] = newDescription
        return token.Ledger.RecordLog("TBillSecurityMeasureUpdated", fmt.Sprintf("Updated security measure: %s", measure))
    }
    return fmt.Errorf("security measure %s not found", measure)
}

// SET_TBILL_COMPLIANCE_PARAMETERS defines compliance requirements for T-Bill operations.
func (token *Syn12Token) SET_TBILL_COMPLIANCE_PARAMETERS(param, value string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceParams[param] = value
    return token.Ledger.RecordLog("TBillComplianceParameterSet", fmt.Sprintf("Set compliance parameter %s to %s", param, value))
}

// GET_TBILL_COMPLIANCE_PARAMETERS retrieves the value of a specific compliance parameter.
func (token *Syn12Token) GET_TBILL_COMPLIANCE_PARAMETERS(param string) (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    value, exists := token.ComplianceParams[param]
    if !exists {
        return "", fmt.Errorf("compliance parameter %s not found", param)
    }
    return value, nil
}

// VALIDATE_TBILL_COMPLIANCE_MEASURES checks if a transaction meets T-Bill compliance standards.
func (token *Syn12Token) VALIDATE_TBILL_COMPLIANCE_MEASURES(transactionID string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if !token.Consensus.ValidateTransaction(transactionID) {
        return false, fmt.Errorf("transaction failed consensus validation")
    }

    return token.ComplianceService.ValidateTransaction(transactionID)
}

// ARCHIVE_TBILL_TRANSACTION_LOG archives older T-Bill transaction logs for compliance.
func (token *Syn12Token) ARCHIVE_TBILL_TRANSACTION_LOG(transactionID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.ArchiveLog(transactionID)
    if err != nil {
        return fmt.Errorf("failed to archive transaction log %s: %v", transactionID, err)
    }
    return nil
}

// CLEAR_TBILL_ARCHIVE deletes older archived logs as per retention policies.
func (token *Syn12Token) CLEAR_TBILL_ARCHIVE() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Ledger.ClearArchive()
}

// ENABLE_TBILL_COMPLIANCE_LOGGING starts logging all compliance-related activities for T-Bills.
func (token *Syn12Token) ENABLE_TBILL_COMPLIANCE_LOGGING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceLogging = true
    return token.Ledger.RecordLog("TBillComplianceLoggingEnabled", "Compliance logging enabled for Treasury Bill transactions")
}

// DISABLE_TBILL_COMPLIANCE_LOGGING stops logging compliance activities for T-Bills.
func (token *Syn12Token) DISABLE_TBILL_COMPLIANCE_LOGGING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceLogging = false
    return token.Ledger.RecordLog("TBillComplianceLoggingDisabled", "Compliance logging disabled for Treasury Bill transactions")
}

// NOTIFY_TBILL_REGULATORS sends encrypted notifications to regulators for significant compliance events.
func (token *Syn12Token) NOTIFY_TBILL_REGULATORS(event string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    encryptedEvent, err := token.Encryption.Encrypt(event)
    if err != nil {
        return fmt.Errorf("encryption failed for regulator notification: %v", err)
    }

    return token.Ledger.NotifyRegulators(encryptedEvent)
}
