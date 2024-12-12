package syn12

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// Syn12Metadata defines the metadata associated with the SYN12 Treasury Bill token.
type Syn12Metadata struct {
    TokenID           string
    Name              string
    Symbol            string
    TBillCode         string
    Issuer            IssuerInfo
    MaturityDate      time.Time
    DiscountRate      float64
    CreationDate      time.Time
    TotalSupply       uint64
    CirculatingSupply uint64
    LegalCompliance   LegalInfo
}

// Syn12Token represents a Treasury Bill token with metadata and compliance tracking.
type Syn12Token struct {
    TokenID       string
    Metadata      Syn12Metadata
    Issuer        string
    Ledger        *ledger.Ledger
    Consensus     *consensus.SynnergyConsensus
    Compliance    *compliance.KYCAmlService
    CentralBank   string
    Encrypted     bool
    mutex         sync.Mutex
}

// CHECK_TBILL_COMPLIANCE_AUDIT performs a compliance audit on a given T-Bill transaction.
func (token *Syn12Token) CHECK_TBILL_COMPLIANCE_AUDIT(transactionID string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if !token.Consensus.ValidateTransaction(transactionID) {
        return false, fmt.Errorf("transaction failed consensus validation")
    }

    auditResult, err := token.Compliance.CheckCompliance(transactionID)
    if err != nil || !auditResult {
        return false, fmt.Errorf("compliance audit failed for transaction %s: %v", transactionID, err)
    }

    err = token.LOG_TBILL_COMPLIANCE_ACTIVITY(transactionID, "Audit Passed")
    if err != nil {
        return false, fmt.Errorf("failed to log compliance activity: %v", err)
    }

    return true, nil
}

// REVIEW_TBILL_COMPLIANCE_AUDIT reviews historical compliance data for a T-Bill transaction.
func (token *Syn12Token) REVIEW_TBILL_COMPLIANCE_AUDIT(transactionID string) (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    history, err := token.Ledger.GetTransactionHistory(transactionID)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve transaction history for %s: %v", transactionID, err)
    }

    complianceVerified := token.Compliance.VerifyHistory(history)
    if !complianceVerified {
        return "", fmt.Errorf("historical compliance review failed for transaction %s", transactionID)
    }

    return history, nil
}

// LOG_TBILL_COMPLIANCE_ACTIVITY securely logs compliance-related activities for a transaction.
func (token *Syn12Token) LOG_TBILL_COMPLIANCE_ACTIVITY(transactionID, activity string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    logEntry := fmt.Sprintf("Transaction %s: %s at %v", transactionID, activity, time.Now())
    encryptedLog, err := token.Encryption.Encrypt(logEntry)
    if err != nil {
        return fmt.Errorf("encryption failed for compliance log: %v", err)
    }

    return token.Ledger.RecordLog("TBillComplianceLog", encryptedLog)
}
