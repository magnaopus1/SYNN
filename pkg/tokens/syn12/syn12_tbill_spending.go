package syn12

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// Syn12Token represents the core structure for Treasury Bill tokens.
type Syn12Token struct {
    TokenID               string
    Metadata              Syn12Metadata
    Issuer                string
    Ledger                *ledger.Ledger
    Consensus             *consensus.SynnergyConsensus
    ComplianceService     *compliance.KYCAmlService
    KYCEnabled            bool
    IssuancePolicy        string
    SpendingAllowances    map[string]uint64 // Account ID to spending allowance
    mutex                 sync.Mutex
}

// APPROVE_TBILL_SPENDING sets a spending allowance for a specified account.
func (token *Syn12Token) APPROVE_TBILL_SPENDING(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.SpendingAllowances[account] = amount
    return token.LOG_TBILL_TRANSACTION("T-Bill spending approved", account, amount)
}

// CHECK_TBILL_ALLOWANCE retrieves the approved spending allowance for an account.
func (token *Syn12Token) CHECK_TBILL_ALLOWANCE(account string) (uint64, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    allowance, exists := token.SpendingAllowances[account]
    if !exists {
        return 0, fmt.Errorf("allowance not set for account %s", account)
    }
    return allowance, nil
}

// SET_TBILL_ISSUER sets the issuer of the T-Bill tokens, typically a central bank.
func (token *Syn12Token) SET_TBILL_ISSUER(issuer string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Issuer = issuer
    return token.LOG_TBILL_TRANSACTION("T-Bill issuer set", issuer, 0)
}

// GET_TBILL_ISSUER retrieves the current issuer of the T-Bill tokens.
func (token *Syn12Token) GET_TBILL_ISSUER() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Issuer
}

// ENABLE_KYC enables KYC requirements for T-Bill transactions.
func (token *Syn12Token) ENABLE_KYC() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.KYCEnabled = true
    return token.LOG_TBILL_TRANSACTION("KYC enabled for T-Bill transactions", "", 0)
}

// DISABLE_KYC disables KYC requirements for T-Bill transactions.
func (token *Syn12Token) DISABLE_KYC() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.KYCEnabled = false
    return token.LOG_TBILL_TRANSACTION("KYC disabled for T-Bill transactions", "", 0)
}

// VALIDATE_KYC checks if an account meets KYC requirements.
func (token *Syn12Token) VALIDATE_KYC(account string) (bool, error) {
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

// SET_TBILL_ISSUANCE_POLICY sets the issuance policy for T-Bill tokens.
func (token *Syn12Token) SET_TBILL_ISSUANCE_POLICY(policy string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.IssuancePolicy = policy
    return token.LOG_TBILL_TRANSACTION("T-Bill issuance policy set", "", 0)
}

// GET_TBILL_ISSUANCE_POLICY retrieves the current issuance policy for T-Bills.
func (token *Syn12Token) GET_TBILL_ISSUANCE_POLICY() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.IssuancePolicy
}

// REDEEM_TBILLS redeems T-Bill tokens, reducing the total supply.
func (token *Syn12Token) REDEEM_TBILLS(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.ReduceSupply(account, amount)
    if err != nil {
        return fmt.Errorf("failed to redeem T-Bills: %v", err)
    }

    return token.LOG_TBILL_TRANSACTION("T-Bills redeemed", account, amount)
}

// ISSUE_NEW_TBILLS mints new T-Bill tokens according to the issuance policy.
func (token *Syn12Token) ISSUE_NEW_TBILLS(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Issuer == "" {
        return fmt.Errorf("T-Bill issuer not set")
    }

    err := token.Ledger.IncreaseSupply(account, amount)
    if err != nil {
        return fmt.Errorf("failed to issue new T-Bills: %v", err)
    }

    return token.LOG_TBILL_TRANSACTION("New T-Bills issued", account, amount)
}

// REVOKE_TBILLS removes a specified amount of T-Bill tokens from circulation.
func (token *Syn12Token) REVOKE_TBILLS(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.BurnTokens(account, amount)
    if err != nil {
        return fmt.Errorf("failed to revoke T-Bills: %v", err)
    }

    return token.LOG_TBILL_TRANSACTION("T-Bills revoked", account, amount)
}

// CHECK_TBILL_COMPLIANCE_STATUS checks if an account complies with T-Bill-specific standards.
func (token *Syn12Token) CHECK_TBILL_COMPLIANCE_STATUS(account string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.ComplianceService.CheckCompliance(account)
}

// LOG_TBILL_TRANSACTION records T-Bill transactions and activities in the ledger.
func (token *Syn12Token) LOG_TBILL_TRANSACTION(description, account string, amount uint64) error {
    logEntry := fmt.Sprintf("Description: %s, Account: %s, Amount: %d, Timestamp: %v", description, account, amount, time.Now())
    encryptedLog, err := token.Encryption.Encrypt(logEntry)
    if err != nil {
        return fmt.Errorf("encryption failed for T-Bill transaction log: %v", err)
    }

    return token.Ledger.RecordLog("TBillTransactionLog", encryptedLog)
}

// VIEW_TBILL_TRANSACTION_HISTORY retrieves the transaction history for a specified account.
func (token *Syn12Token) VIEW_TBILL_TRANSACTION_HISTORY(account string) (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    history, err := token.Ledger.GetTransactionHistory(account)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve transaction history for account %s: %v", account, err)
    }
    return history, nil
}
