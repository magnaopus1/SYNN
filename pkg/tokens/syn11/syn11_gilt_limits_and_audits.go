package syn11

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SYN11Token represents the core structure for Central Bank Digital Gilt tokens.
type SYN11Token struct {
    TokenID               string
    Metadata              Syn11Metadata
    Issuer                string
    Ledger                *ledger.Ledger
    Consensus             *consensus.SynnergyConsensus
    Encrypted             bool
    GiltTransactionLimits map[string]uint64
    AutoMintingEnabled    bool
    TransactionLogging    bool
    AuditStatus           string
    GiltHistory           []string
    mutex                 sync.Mutex
}

// SET_GILT_TRANSACTION_LIMIT sets a limit on the amount that can be transacted for specific gilt operations.
func (token *SYN11Token) SET_GILT_TRANSACTION_LIMIT(operation string, limit uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.GiltTransactionLimits[operation] = limit
    return token.Ledger.RecordLog("GiltTransactionLimitSet", fmt.Sprintf("Limit set for %s: %d", operation, limit))
}

// GET_GILT_TRANSACTION_LIMIT retrieves the transaction limit for a specified operation.
func (token *SYN11Token) GET_GILT_TRANSACTION_LIMIT(operation string) (uint64, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    limit, exists := token.GiltTransactionLimits[operation]
    if !exists {
        return 0, fmt.Errorf("no limit set for operation %s", operation)
    }
    return limit, nil
}

// ENABLE_AUTO_MINTING enables automatic minting of new gilt tokens based on specified criteria.
func (token *SYN11Token) ENABLE_AUTO_MINTING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AutoMintingEnabled = true
    return token.Ledger.RecordLog("AutoMintingEnabled", "Automatic minting enabled")
}

// DISABLE_AUTO_MINTING disables the automatic minting feature.
func (token *SYN11Token) DISABLE_AUTO_MINTING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AutoMintingEnabled = false
    return token.Ledger.RecordLog("AutoMintingDisabled", "Automatic minting disabled")
}

// GENERATE_GILT_REPORT generates a comprehensive report on gilt transactions and audits.
func (token *SYN11Token) GENERATE_GILT_REPORT() (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    report, err := token.Ledger.GenerateTransactionReport()
    if err != nil {
        return "", fmt.Errorf("failed to generate gilt report: %v", err)
    }

    encryptedReport, err := token.Encryption.Encrypt(report)
    if err != nil {
        return "", fmt.Errorf("encryption failed for gilt report: %v", err)
    }
    return encryptedReport, nil
}

// FETCH_GILT_AUDIT_TRAIL retrieves the full audit trail of gilt transactions.
func (token *SYN11Token) FETCH_GILT_AUDIT_TRAIL() ([]string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    auditTrail, err := token.Ledger.GetAuditTrail()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch gilt audit trail: %v", err)
    }
    return auditTrail, nil
}

// INITIATE_GILT_AUDIT initiates an audit process for gilt transactions.
func (token *SYN11Token) INITIATE_GILT_AUDIT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AuditStatus = "In Progress"
    return token.Ledger.RecordLog("GiltAuditInitiated", "Gilt audit initiated")
}

// COMPLETE_GILT_AUDIT completes the ongoing audit process for gilt transactions.
func (token *SYN11Token) COMPLETE_GILT_AUDIT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AuditStatus = "Completed"
    return token.Ledger.RecordLog("GiltAuditCompleted", "Gilt audit completed")
}

// CHECK_GILT_AUDIT_STATUS checks the current status of the gilt audit process.
func (token *SYN11Token) CHECK_GILT_AUDIT_STATUS() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.AuditStatus
}

// ENABLE_GILT_TRANSACTIONS_LOGGING enables logging for all gilt transactions.
func (token *SYN11Token) ENABLE_GILT_TRANSACTIONS_LOGGING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TransactionLogging = true
    return token.Ledger.RecordLog("GiltTransactionLoggingEnabled", "Transaction logging enabled for gilt")
}

// DISABLE_GILT_TRANSACTIONS_LOGGING disables logging for gilt transactions.
func (token *SYN11Token) DISABLE_GILT_TRANSACTIONS_LOGGING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TransactionLogging = false
    return token.Ledger.RecordLog("GiltTransactionLoggingDisabled", "Transaction logging disabled for gilt")
}

// GET_LOGGING_STATUS retrieves the current logging status for gilt transactions.
func (token *SYN11Token) GET_LOGGING_STATUS() bool {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.TransactionLogging
}

// RESET_GILT resets the gilt token to its default configuration, clearing settings and transaction history.
func (token *SYN11Token) RESET_GILT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.GiltTransactionLimits = make(map[string]uint64)
    token.AutoMintingEnabled = false
    token.TransactionLogging = false
    token.AuditStatus = "Not Started"
    token.GiltHistory = []string{}

    return token.Ledger.RecordLog("GiltReset", "Gilt token reset to default state")
}

// FETCH_GILT_HISTORY retrieves the history of all gilt-related transactions.
func (token *SYN11Token) FETCH_GILT_HISTORY() ([]string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    history, err := token.Ledger.GetTokenHistory()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch gilt history: %v", err)
    }
    return history, nil
}

// SET_GILT_LIMITS defines limits for the maximum supply and transfer volumes for gilt tokens.
func (token *SYN11Token) SET_GILT_LIMITS(maxSupply, maxTransfer uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.GiltTransactionLimits["MaxSupply"] = maxSupply
    token.GiltTransactionLimits["MaxTransfer"] = maxTransfer

    return token.Ledger.RecordLog("GiltLimitsSet", fmt.Sprintf("Gilt limits set: MaxSupply=%d, MaxTransfer=%d", maxSupply, maxTransfer))
}
