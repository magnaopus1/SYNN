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
    ComplianceService     *compliance.KYCAmlService
    KYCEnabled            bool
    IssuancePolicy        string
    SpendingAllowances    map[string]uint64 // Account ID to spending allowance
    mutex                 sync.Mutex
}

// APPROVE_GILT_SPENDING sets a spending allowance for a specified account.
func (token *SYN11Token) APPROVE_GILT_SPENDING(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.SpendingAllowances[account] = amount
    return token.LOG_GILT_TRANSACTION("Gilt spending approved", account, amount)
}

// CHECK_GILT_ALLOWANCE retrieves the approved spending allowance for an account.
func (token *SYN11Token) CHECK_GILT_ALLOWANCE(account string) (uint64, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    allowance, exists := token.SpendingAllowances[account]
    if !exists {
        return 0, fmt.Errorf("allowance not set for account %s", account)
    }
    return allowance, nil
}

// SET_GILT_ISSUER sets the issuer of the gilt tokens, typically a central bank.
func (token *SYN11Token) SET_GILT_ISSUER(issuer string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Issuer = issuer
    return token.LOG_GILT_TRANSACTION("Gilt issuer set", issuer, 0)
}

// GET_GILT_ISSUER retrieves the current issuer of the gilt tokens.
func (token *SYN11Token) GET_GILT_ISSUER() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Issuer
}

// ENABLE_KYC enables KYC requirements for transactions involving gilt tokens.
func (token *SYN11Token) ENABLE_KYC() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.KYCEnabled = true
    return token.LOG_GILT_TRANSACTION("KYC enabled for gilt transactions", "", 0)
}

// DISABLE_KYC disables KYC requirements for gilt transactions.
func (token *SYN11Token) DISABLE_KYC() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.KYCEnabled = false
    return token.LOG_GILT_TRANSACTION("KYC disabled for gilt transactions", "", 0)
}

// VALIDATE_KYC checks if an account meets KYC requirements.
func (token *SYN11Token) VALIDATE_KYC(account string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if !token.KYCEnabled {
        return true, nil // No KYC required
    }

    kycStatus, err := token.ComplianceService.CheckKYC(account)
    if err != nil || !kycStatus {
        return false, fmt.Errorf("KYC validation failed for account %s: %v", account, err)
    }
    return true, nil
}

// SET_GILT_ISSUANCE_POLICY sets the issuance policy for gilt tokens.
func (token *SYN11Token) SET_GILT_ISSUANCE_POLICY(policy string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.IssuancePolicy = policy
    return token.LOG_GILT_TRANSACTION("Gilt issuance policy set", "", 0)
}

// GET_GILT_ISSUANCE_POLICY retrieves the current issuance policy.
func (token *SYN11Token) GET_GILT_ISSUANCE_POLICY() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.IssuancePolicy
}

// REDEEM_GILTS reduces the total supply by a specified amount, simulating redemption.
func (token *SYN11Token) REDEEM_GILTS(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.ReduceSupply(account, amount)
    if err != nil {
        return fmt.Errorf("failed to redeem gilts: %v", err)
    }

    return token.LOG_GILT_TRANSACTION("Gilts redeemed", account, amount)
}

// ISSUE_NEW_GILTS mints new gilt tokens according to the issuance policy.
func (token *SYN11Token) ISSUE_NEW_GILTS(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Issuer == "" {
        return fmt.Errorf("gilt issuer not set")
    }

    err := token.Ledger.IncreaseSupply(account, amount)
    if err != nil {
        return fmt.Errorf("failed to issue new gilts: %v", err)
    }

    return token.LOG_GILT_TRANSACTION("New gilts issued", account, amount)
}

// REVOKE_GILTS removes a specific amount of gilt tokens from circulation.
func (token *SYN11Token) REVOKE_GILTS(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.BurnTokens(account, amount)
    if err != nil {
        return fmt.Errorf("failed to revoke gilts: %v", err)
    }

    return token.LOG_GILT_TRANSACTION("Gilts revoked", account, amount)
}

// CHECK_GILT_COMPLIANCE_STATUS checks if an account complies with gilt-specific standards.
func (token *SYN11Token) CHECK_GILT_COMPLIANCE_STATUS(account string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.ComplianceService.CheckCompliance(account)
}

// LOG_GILT_TRANSACTION records gilt transactions and activities in the ledger.
func (token *SYN11Token) LOG_GILT_TRANSACTION(description, account string, amount uint64) error {
    logEntry := fmt.Sprintf("Description: %s, Account: %s, Amount: %d, Timestamp: %v", description, account, amount, time.Now())
    encryptedLog, err := token.Encryption.Encrypt(logEntry)
    if err != nil {
        return fmt.Errorf("encryption failed for gilt transaction log: %v", err)
    }

    return token.Ledger.RecordLog("GiltTransactionLog", encryptedLog)
}

// VIEW_GILT_TRANSACTION_HISTORY retrieves the transaction history for a specified account.
func (token *SYN11Token) VIEW_GILT_TRANSACTION_HISTORY(account string) (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    history, err := token.Ledger.GetTransactionHistory(account)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve transaction history for account %s: %v", account, err)
    }
    return history, nil
}
