package syn721

import (
    "sync"
    "fmt"
    "time"
    "path/to/ledger"
    "path/to/consensus"
    "path/to/encryption"
)

// SYN721Token struct represents an NFT token with advanced feature capabilities.
type SYN721Token struct {
    mutex                   sync.Mutex
    TokenID                 string
    Owner                   string
    ExpirationDate          *time.Time               // Optional expiration date
    TransactionHistory      []TransactionRecord      // Transaction history log
    ComplianceMonitoring    bool                     // Compliance monitoring enabled status
    NonFungibleFeatures     map[string]bool          // Enabled non-fungible token features
    Ledger                  *ledger.Ledger
    Consensus               *consensus.SynnergyConsensus
    Encryption              *encryption.Encryption
}

// LOG_TOKEN_PURCHASE_EVENT logs a token purchase event in the transaction history.
func (token *SYN721Token) LOG_TOKEN_PURCHASE_EVENT(buyer string, purchaseDetails string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    transaction := TransactionRecord{
        TokenID:    token.TokenID,
        Buyer:      buyer,
        Details:    purchaseDetails,
        Timestamp:  time.Now(),
    }
    token.TransactionHistory = append(token.TransactionHistory, transaction)
    return token.Ledger.RecordLog("TokenPurchaseEvent", fmt.Sprintf("Token %s purchased by %s", token.TokenID, buyer))
}

// GET_TOKEN_TRANSACTION_HISTORY retrieves the full transaction history for the token.
func (token *SYN721Token) GET_TOKEN_TRANSACTION_HISTORY() []TransactionRecord {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.TransactionHistory
}

// CHECK_TOKEN_EXPIRATION_STATUS checks if the token is expired based on the expiration date.
func (token *SYN721Token) CHECK_TOKEN_EXPIRATION_STATUS() (bool, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    if token.ExpirationDate == nil {
        return false, fmt.Errorf("expiration date not set")
    }
    return time.Now().After(*token.ExpirationDate), nil
}

// ENABLE_NON_FUNGIBLE_TOKEN_FEATURES enables specific non-fungible token features.
func (token *SYN721Token) ENABLE_NON_FUNGIBLE_TOKEN_FEATURES(feature string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.NonFungibleFeatures == nil {
        token.NonFungibleFeatures = make(map[string]bool)
    }
    token.NonFungibleFeatures[feature] = true
    return token.Ledger.RecordLog("FeatureEnabled", fmt.Sprintf("Feature %s enabled for token %s", feature, token.TokenID))
}

// DISABLE_NON_FUNGIBLE_TOKEN_FEATURES disables specific non-fungible token features.
func (token *SYN721Token) DISABLE_NON_FUNGIBLE_TOKEN_FEATURES(feature string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.NonFungibleFeatures == nil {
        return fmt.Errorf("no features to disable")
    }
    delete(token.NonFungibleFeatures, feature)
    return token.Ledger.RecordLog("FeatureDisabled", fmt.Sprintf("Feature %s disabled for token %s", feature, token.TokenID))
}

// GET_NON_FUNGIBLE_TOKEN_FEATURES retrieves the current enabled non-fungible features for the token.
func (token *SYN721Token) GET_NON_FUNGIBLE_TOKEN_FEATURES() map[string]bool {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.NonFungibleFeatures
}

// ENABLE_COMPLIANCE_MONITORING enables compliance monitoring for the token.
func (token *SYN721Token) ENABLE_COMPLIANCE_MONITORING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceMonitoring = true
    return token.Ledger.RecordLog("ComplianceMonitoringEnabled", fmt.Sprintf("Compliance monitoring enabled for token %s", token.TokenID))
}

// DISABLE_COMPLIANCE_MONITORING disables compliance monitoring for the token.
func (token *SYN721Token) DISABLE_COMPLIANCE_MONITORING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ComplianceMonitoring = false
    return token.Ledger.RecordLog("ComplianceMonitoringDisabled", fmt.Sprintf("Compliance monitoring disabled for token %s", token.TokenID))
}

// LOG_TOKEN_COMPLIANCE_UPDATE logs any updates or changes in the token's compliance status.
func (token *SYN721Token) LOG_TOKEN_COMPLIANCE_UPDATE(updateDetails string) error {
    return token.Ledger.RecordLog("TokenComplianceUpdate", fmt.Sprintf("Compliance update for token %s: %s", token.TokenID, updateDetails))
}
