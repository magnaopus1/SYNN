package syn131

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// ENABLE_AUTO_VALUATION enables automated valuation updates for the asset.
func (token *Syn131Token) ENABLE_AUTO_VALUATION() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AutoValuationEnabled = true
    return token.Ledger.RecordLog("AutoValuationEnabled", fmt.Sprintf("Automatic valuation enabled for asset %s", token.ID))
}

// DISABLE_AUTO_VALUATION disables automated valuation updates for the asset.
func (token *Syn131Token) DISABLE_AUTO_VALUATION() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AutoValuationEnabled = false
    return token.Ledger.RecordLog("AutoValuationDisabled", fmt.Sprintf("Automatic valuation disabled for asset %s", token.ID))
}

// INITIATE_ASSET_AUDIT begins an audit process for the asset.
func (token *Syn131Token) INITIATE_ASSET_AUDIT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.IntangibleAssetValuation.AuditStatus = "In Progress"
    return token.Ledger.RecordLog("AssetAuditInitiated", fmt.Sprintf("Audit initiated for asset %s", token.ID))
}

// COMPLETE_ASSET_AUDIT finalizes the audit process for the asset.
func (token *Syn131Token) COMPLETE_ASSET_AUDIT(outcome string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.IntangibleAssetValuation.AuditStatus = "Completed"
    token.IntangibleAssetValuation.AuditOutcome = outcome
    return token.Ledger.RecordLog("AssetAuditCompleted", fmt.Sprintf("Audit completed for asset %s with outcome: %s", token.ID, outcome))
}

// CHECK_ASSET_AUDIT_STATUS retrieves the current status of the asset audit.
func (token *Syn131Token) CHECK_ASSET_AUDIT_STATUS() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.IntangibleAssetValuation.AuditStatus
}

// ENABLE_AUDIT_LOGGING enables logging for all audit-related activities.
func (token *Syn131Token) ENABLE_AUDIT_LOGGING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceLoggingEnabled = true
    return token.Ledger.RecordLog("AuditLoggingEnabled", fmt.Sprintf("Audit logging enabled for asset %s", token.ID))
}

// DISABLE_AUDIT_LOGGING disables logging for audit activities.
func (token *Syn131Token) DISABLE_AUDIT_LOGGING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceLoggingEnabled = false
    return token.Ledger.RecordLog("AuditLoggingDisabled", fmt.Sprintf("Audit logging disabled for asset %s", token.ID))
}

// FETCH_AUDIT_LOGGING_STATUS checks if audit logging is currently enabled.
func (token *Syn131Token) FETCH_AUDIT_LOGGING_STATUS() bool {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.ComplianceLoggingEnabled
}

// UPDATE_LEASE_TERMS modifies the lease terms associated with the asset.
func (token *Syn131Token) UPDATE_LEASE_TERMS(newTerms LeaseTerms) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.LeaseTerms = append(token.LeaseTerms, newTerms)
    return token.Ledger.RecordLog("LeaseTermsUpdated", fmt.Sprintf("Lease terms updated for asset %s", token.ID))
}

// GET_LEASE_TERMS retrieves the current lease terms for the asset.
func (token *Syn131Token) GET_LEASE_TERMS() []LeaseTerms {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.LeaseTerms
}

// REGISTER_ASSET_NOTIFICATION sets up a notification setting for the asset.
func (token *Syn131Token) REGISTER_ASSET_NOTIFICATION(notification NotificationSettings) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.NotificationSettings = notification
    return token.Ledger.RecordLog("NotificationRegistered", fmt.Sprintf("Notification registered for asset %s", token.ID))
}

// TRIGGER_ASSET_NOTIFICATION sends a notification based on predefined triggers.
func (token *Syn131Token) TRIGGER_ASSET_NOTIFICATION(trigger string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    notificationMessage := fmt.Sprintf("Notification for asset %s: Triggered by %s", token.ID, trigger)
    encryptedMessage, err := token.Encryption.Encrypt(notificationMessage)
    if err != nil {
        return fmt.Errorf("failed to encrypt notification: %v", err)
    }

    return token.Ledger.RecordNotification("AssetNotification", encryptedMessage)
}

// ENABLE_AUTO_NOTIFICATIONS enables automated notifications for the asset.
func (token *Syn131Token) ENABLE_AUTO_NOTIFICATIONS() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AutoNotificationEnabled = true
    return token.Ledger.RecordLog("AutoNotificationEnabled", fmt.Sprintf("Automatic notifications enabled for asset %s", token.ID))
}

// DISABLE_AUTO_NOTIFICATIONS disables automated notifications for the asset.
func (token *Syn131Token) DISABLE_AUTO_NOTIFICATIONS() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AutoNotificationEnabled = false
    return token.Ledger.RecordLog("AutoNotificationDisabled", fmt.Sprintf("Automatic notifications disabled for asset %s", token.ID))
}

// CONNECT_AUTOVALUATION_AI integrates AI for automated valuation updates.
func (token *Syn131Token) CONNECT_AUTOVALUATION_AI(modelID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.ConnectAIModel(token.ID, modelID, "valuation")
    if err != nil {
        return fmt.Errorf("failed to connect auto-valuation AI model: %v", err)
    }
    return token.Ledger.RecordLog("AutoValuationAIConnected", fmt.Sprintf("AI model %s connected for auto-valuation of asset %s", modelID, token.ID))
}