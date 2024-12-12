package syn10

import (
	"fmt"

)

// ENABLE_SECURITY_AUDIT activates enhanced security monitoring.
func (token *SYN10Token) enableSecurityAudit() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.SecurityProtocols["audit"] = "enabled"
    err := token.Ledger.LogSecurityActivity("Security audit enabled")
    if err != nil {
        return fmt.Errorf("failed to enable security audit: %v", err)
    }
    return nil
}

// DISABLE_SECURITY_AUDIT deactivates security monitoring.
func (token *SYN10Token) disableSecurityAudit() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.SecurityProtocols["audit"] = "disabled"
    return token.Ledger.LogSecurityActivity("Security audit disabled")
}


// INTEGRATE_SECURITY_MEASURES adds new security measures to the token.
func (token *SYN10Token) integrateSecurityMeasures(measure string, description string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.SecurityProtocols[measure] = description
    return token.Ledger.LogSecurityActivity(fmt.Sprintf("Integrated security measure: %s", measure))
}


// REGISTER_SECURITY_EVENT logs a security event for auditing.
func (token *SYN10Token) registerSecurityEvent(event string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    encryptedEvent, err := token.Encryption.Encrypt(event)
    if err != nil {
        return fmt.Errorf("encryption failed: %v", err)
    }

    return token.Ledger.LogSecurityActivity(fmt.Sprintf("Registered event: %s", encryptedEvent))
}


// FETCH_SECURITY_REPORT retrieves a summary of recent security events.
func (token *SYN10Token) FetchSecurityReport() (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    report, err := token.Ledger.getSecurityLog()
    if err != nil {
        return "", fmt.Errorf("failed to fetch security report: %v", err)
    }
    return report, nil
}


// UPDATE_SECURITY_MEASURES updates existing security measures.
func (token *SYN10Token) UpdateSecurityMeasures(measure string, newDescription string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if _, exists := token.SecurityProtocols[measure]; exists {
        token.SecurityProtocols[measure] = newDescription
        return token.Ledger.LogSecurityActivity(fmt.Sprintf("Updated security measure: %s", measure))
    }
    return fmt.Errorf("security measure not found")
}


// SET_COMPLIANCE_PARAMETERS establishes parameters for regulatory compliance.
func (token *SYN10Token) SetComplianceParameters(param string, value string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceParams[param] = value
    return token.Ledger.LogComplianceActivity(fmt.Sprintf("Set compliance parameter: %s", param))
}


// GET_COMPLIANCE_PARAMETERS retrieves compliance parameters.
func (token *SYN10Token) GetComplianceParameters(param string) (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    value, exists := token.ComplianceParams[param]
    if !exists {
        return "", fmt.Errorf("compliance parameter not found")
    }
    return value, nil
}



// ARCHIVE_TRANSACTION_LOG moves old logs into a secure archive.
func (token *SYN10Token) ArchiveTransactionLog(transactionID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.ArchiveTransactionLog(transactionID)
    if err != nil {
        return fmt.Errorf("failed to archive transaction log: %v", err)
    }

    return nil
}


// CLEAR_ARCHIVE removes logs from the archive securely.
func (token *SYN10Token) ClearArchive() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Ledger.ClearArchive()
}

// ENABLE_COMPLIANCE_LOGGING activates logging for all compliance activities.
func (token *SYN10Token) EnableComplianceLogging() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceLogging = true
    return token.Ledger.LogComplianceActivity("Compliance logging enabled")
}


// DISABLE_COMPLIANCE_LOGGING deactivates logging for compliance activities.
func (token *SYN10Token) disableComplianceLogging() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceLogging = false
    return token.Ledger.LogComplianceActivity("Compliance logging disabled")
}


// NOTIFY_REGULATORS sends notifications to regulators for critical compliance events.
func (token *SYN10Token) NotifyRegulators(event string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    encryptedEvent, err := token.Encryption.Encrypt(event)
    if err != nil {
        return fmt.Errorf("encryption failed: %v", err)
    }

    return token.Ledger.NotifyRegulators(encryptedEvent)
}
