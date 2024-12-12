package syn131

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// LOG_ASSET_SECURITY_EVENT logs a security-related event for the asset.
func (token *Syn131Token) LOG_ASSET_SECURITY_EVENT(eventDescription string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    logEntry := fmt.Sprintf("Event: %s, Timestamp: %v", eventDescription, time.Now())
    encryptedLog, err := token.Encryption.Encrypt(logEntry)
    if err != nil {
        return fmt.Errorf("encryption failed for security event: %v", err)
    }

    return token.Ledger.RecordLog("AssetSecurityEvent", encryptedLog)
}

// SET_ENCRYPTION_FOR_ASSET enables or disables encryption for the asset terms.
func (token *Syn131Token) SET_ENCRYPTION_FOR_ASSET(enable bool) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if enable {
        encryptedTerms, err := token.Encryption.Encrypt(token.Terms)
        if err != nil {
            return fmt.Errorf("encryption failed for asset terms: %v", err)
        }
        token.EncryptedTerms = encryptedTerms
    } else {
        token.EncryptedTerms = ""
    }
    
    return token.Ledger.RecordLog("AssetEncryptionStatus", fmt.Sprintf("Encryption %t for asset %s", enable, token.ID))
}

// GET_ENCRYPTION_STATUS retrieves the encryption status of the asset terms.
func (token *Syn131Token) GET_ENCRYPTION_STATUS() bool {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.EncryptedTerms != ""
}

// VALIDATE_ASSET_SECURITY checks if the asset meets the required security standards.
func (token *Syn131Token) VALIDATE_ASSET_SECURITY() (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.EncryptedTerms == "" {
        return false, fmt.Errorf("asset terms are not encrypted, fails security validation")
    }
    return true, nil
}

// CHECK_ASSET_COMPLIANCE_STATUS verifies the asset's compliance status.
func (token *Syn131Token) CHECK_ASSET_COMPLIANCE_STATUS() (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.IntangibleAssetStatus.Compliant, nil
}

// LOG_ASSET_COMPLIANCE_ACTIVITY records compliance-related activities.
func (token *Syn131Token) LOG_ASSET_COMPLIANCE_ACTIVITY(activity string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    logEntry := fmt.Sprintf("Compliance Activity: %s, Timestamp: %v", activity, time.Now())
    encryptedLog, err := token.Encryption.Encrypt(logEntry)
    if err != nil {
        return fmt.Errorf("failed to encrypt compliance activity log: %v", err)
    }

    return token.Ledger.RecordLog("AssetComplianceActivity", encryptedLog)
}

// ENABLE_ASSET_COMPLIANCE_LOGGING enables compliance activity logging.
func (token *Syn131Token) ENABLE_ASSET_COMPLIANCE_LOGGING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceLoggingEnabled = true
    return token.Ledger.RecordLog("ComplianceLoggingEnabled", fmt.Sprintf("Compliance logging enabled for asset %s", token.ID))
}

// DISABLE_ASSET_COMPLIANCE_LOGGING disables compliance activity logging.
func (token *Syn131Token) DISABLE_ASSET_COMPLIANCE_LOGGING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceLoggingEnabled = false
    return token.Ledger.RecordLog("ComplianceLoggingDisabled", fmt.Sprintf("Compliance logging disabled for asset %s", token.ID))
}

// FETCH_ASSET_COMPLIANCE_REPORT retrieves the latest compliance report for the asset.
func (token *Syn131Token) FETCH_ASSET_COMPLIANCE_REPORT() (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    report, err := token.Ledger.GetComplianceReport(token.ID)
    if err != nil {
        return "", fmt.Errorf("failed to fetch compliance report: %v", err)
    }
    return report, nil
}

// INITIATE_ASSET_REVIEW initiates a formal compliance review for the asset.
func (token *Syn131Token) INITIATE_ASSET_REVIEW() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.IntangibleAssetStatus.ReviewStatus = "In Progress"
    return token.Ledger.RecordLog("AssetReviewInitiated", fmt.Sprintf("Compliance review initiated for asset %s", token.ID))
}

// COMPLETE_ASSET_REVIEW completes the compliance review and records the outcome.
func (token *Syn131Token) COMPLETE_ASSET_REVIEW(outcome string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.IntangibleAssetStatus.ReviewStatus = "Completed"
    token.IntangibleAssetStatus.ReviewOutcome = outcome
    return token.Ledger.RecordLog("AssetReviewCompleted", fmt.Sprintf("Compliance review completed for asset %s with outcome: %s", token.ID, outcome))
}

// GET_ASSET_REVIEW_STATUS retrieves the current compliance review status.
func (token *Syn131Token) GET_ASSET_REVIEW_STATUS() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.IntangibleAssetStatus.ReviewStatus
}

// ENABLE_ASSET_MARKET_INTEGRATION allows the asset to be integrated with external markets.
func (token *Syn131Token) ENABLE_ASSET_MARKET_INTEGRATION() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.MarketIntegrationEnabled = true
    return token.Ledger.RecordLog("MarketIntegrationEnabled", fmt.Sprintf("Market integration enabled for asset %s", token.ID))
}

// DISABLE_ASSET_MARKET_INTEGRATION prevents the asset from being integrated with external markets.
func (token *Syn131Token) DISABLE_ASSET_MARKET_INTEGRATION() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.MarketIntegrationEnabled = false
    return token.Ledger.RecordLog("MarketIntegrationDisabled", fmt.Sprintf("Market integration disabled for asset %s", token.ID))
}
