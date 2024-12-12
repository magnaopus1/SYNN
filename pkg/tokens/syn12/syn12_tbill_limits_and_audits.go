package syn12

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// Syn12Token represents the core structure for Treasury Bill tokens.
type Syn12Token struct {
    TokenID                string
    Metadata               Syn12Metadata
    Issuer                 string
    Ledger                 *ledger.Ledger
    Consensus              *consensus.SynnergyConsensus
    Encrypted              bool
    TBillTransactionLimits map[string]uint64
    AutoMintingEnabled     bool
    TransactionLogging     bool
    AuditStatus            string
    TBillHistory           []string
    mutex                  sync.Mutex
}

// SET_TBILL_TRANSACTION_LIMIT sets transaction limits for specific Treasury Bill operations.
func (token *Syn12Token) SET_TBILL_TRANSACTION_LIMIT(operation string, limit uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TBillTransactionLimits[operation] = limit
    return token.Ledger.RecordLog("TBillTransactionLimitSet", fmt.Sprintf("Set limit for %s: %d", operation, limit))
}

// GET_TBILL_TRANSACTION_LIMIT retrieves the transaction limit for a specified operation.
func (token *Syn12Token) GET_TBILL_TRANSACTION_LIMIT(operation string) (uint64, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    limit, exists := token.TBillTransactionLimits[operation]
    if !exists {
        return 0, fmt.Errorf("no limit set for operation %s", operation)
    }
    return limit, nil
}

// ENABLE_AUTO_MINTING activates auto-minting functionality based on defined criteria.
func (token *Syn12Token) ENABLE_AUTO_MINTING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AutoMintingEnabled = true
    return token.Ledger.RecordLog("AutoMintingEnabled", "Automatic minting enabled for Treasury Bills")
}

// DISABLE_AUTO_MINTING deactivates auto-minting functionality.
func (token *Syn12Token) DISABLE_AUTO_MINTING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AutoMintingEnabled = false
    return token.Ledger.RecordLog("AutoMintingDisabled", "Automatic minting disabled for Treasury Bills")
}

// GENERATE_TBILL_REPORT generates an audit report on Treasury Bill transactions and activities.
func (token *Syn12Token) GENERATE_TBILL_REPORT() (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    report, err := token.Ledger.GenerateTransactionReport()
    if err != nil {
        return "", fmt.Errorf("failed to generate Treasury Bill report: %v", err)
    }

    encryptedReport, err := token.Encryption.Encrypt(report)
    if err != nil {
        return "", fmt.Errorf("encryption failed for Treasury Bill report: %v", err)
    }
    return encryptedReport, nil
}

// FETCH_TBILL_AUDIT_TRAIL retrieves the full audit trail for Treasury Bill transactions.
func (token *Syn12Token) FETCH_TBILL_AUDIT_TRAIL() ([]string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    auditTrail, err := token.Ledger.GetAuditTrail()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch Treasury Bill audit trail: %v", err)
    }
    return auditTrail, nil
}

// INITIATE_TBILL_AUDIT begins the audit process for Treasury Bill transactions.
func (token *Syn12Token) INITIATE_TBILL_AUDIT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AuditStatus = "In Progress"
    return token.Ledger.RecordLog("TBillAuditInitiated", "Audit process initiated for Treasury Bills")
}

// COMPLETE_TBILL_AUDIT finalizes the current audit process for Treasury Bill transactions.
func (token *Syn12Token) COMPLETE_TBILL_AUDIT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AuditStatus = "Completed"
    return token.Ledger.RecordLog("TBillAuditCompleted", "Audit process completed for Treasury Bills")
}

// CHECK_TBILL_AUDIT_STATUS provides the current status of the Treasury Bill audit process.
func (token *Syn12Token) CHECK_TBILL_AUDIT_STATUS() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.AuditStatus
}

// ENABLE_TBILL_TRANSACTIONS_LOGGING enables logging for all Treasury Bill transactions.
func (token *Syn12Token) ENABLE_TBILL_TRANSACTIONS_LOGGING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TransactionLogging = true
    return token.Ledger.RecordLog("TBillTransactionLoggingEnabled", "Transaction logging enabled for Treasury Bills")
}

// DISABLE_TBILL_TRANSACTIONS_LOGGING disables logging for Treasury Bill transactions.
func (token *Syn12Token) DISABLE_TBILL_TRANSACTIONS_LOGGING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TransactionLogging = false
    return token.Ledger.RecordLog("TBillTransactionLoggingDisabled", "Transaction logging disabled for Treasury Bills")
}

// GET_LOGGING_STATUS retrieves the current status of transaction logging.
func (token *Syn12Token) GET_LOGGING_STATUS() bool {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.TransactionLogging
}

// RESET_TBILL resets the Treasury Bill token to its default state, clearing limits and logs.
func (token *Syn12Token) RESET_TBILL() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TBillTransactionLimits = map[string]uint64{}
    token.AutoMintingEnabled = false
    token.TransactionLogging = false
    token.AuditStatus = "Not Started"
    token.TBillHistory = []string{}

    return token.Ledger.RecordLog("TBillReset", "Treasury Bill reset to default state")
}

// FETCH_TBILL_HISTORY retrieves the historical transaction record for Treasury Bills.
func (token *Syn12Token) FETCH_TBILL_HISTORY() ([]string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    history, err := token.Ledger.GetTokenHistory()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch Treasury Bill history: %v", err)
    }
    return history, nil
}

// SET_TBILL_LIMITS establishes maximum supply and transfer limits for Treasury Bill tokens.
func (token *Syn12Token) SET_TBILL_LIMITS(maxSupply, maxTransfer uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TBillTransactionLimits["MaxSupply"] = maxSupply
    token.TBillTransactionLimits["MaxTransfer"] = maxTransfer

    return token.Ledger.RecordLog("TBillLimitsSet", fmt.Sprintf("Treasury Bill limits set: MaxSupply=%d, MaxTransfer=%d", maxSupply, maxTransfer))
}
