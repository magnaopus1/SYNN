package syn10

import (
	"fmt"
	"time"
)

// APPROVE_SPENDING approves a spending allowance for a specific account.
func (token *SYN10Token) ApproveSpending(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Allowances[account] = amount
    log := SYN10TransactionLog{
        LogID:      generateUniqueID(),
        Description: "Spending approved",
        Account:     account,
        Amount:      amount,
        Timestamp:   time.Now(),
    }

    return token.Ledger.RecordLog(log)
}


// CHECK_ALLOWANCE retrieves the spending allowance for a specific account.
func (token *SYN10Token) checkAllowance(account string) (uint64, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    allowance, exists := token.Allowances[account]
    if !exists {
        return 0, fmt.Errorf("allowance not found for account %s", account)
    }
    return allowance, nil
}


// SET_TOKEN_ISSUER defines the issuer of the token, typically the central bank.
func (token *SYN10Token) setTokenIssuer(issuer string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TokenIssuer = issuer
    log := SYN10TransactionLog{
        LogID:      generateUniqueID(),
        Description: "Token issuer set",
        Account:     issuer,
        Amount:      0,
        Timestamp:   time.Now(),
    }

    return token.Ledger.RecordLog(log)
}


func (token *SYN10Token) redeemTokens(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.reduceSupply(account, amount)
    if err != nil {
        return fmt.Errorf("failed to redeem tokens: %v", err)
    }

    log := SYN10TransactionLog{
        LogID:      generateUniqueID(),
        Description: "Tokens redeemed",
        Account:     account,
        Amount:      amount,
        Timestamp:   time.Now(),
    }

    return token.Ledger.recordLog(log)
}

func (token *SYN10Token) issueNewTokens(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.TokenIssuer == "" {
        return fmt.Errorf("token issuer not set")
    }

    err := token.Ledger.increaseSupply(account, amount)
    if err != nil {
        return fmt.Errorf("failed to issue new tokens: %v", err)
    }

    log := SYN10TransactionLog{
        LogID:      generateUniqueID(),
        Description: "New tokens issued",
        Account:     account,
        Amount:      amount,
        Timestamp:   time.Now(),
    }

    return token.Ledger.recordLog(log)
}


func (token *SYN10Token) ViewTransactionHistory(account string) (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Ledger.GetTransactionHistory(account)
}


// GET_TOKEN_ISSUER retrieves the current token issuer.
func (token *SYN10Token) GET_TOKEN_ISSUER() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    return token.TokenIssuer
}

// ENABLE_KYC enables the KYC process for transactions involving this token.
func (token *SYN10Token) ENABLE_KYC() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    token.KYCEnabled = true
    return token.LOG_TRANSACTION("KYC enabled", "", 0)
}

// DISABLE_KYC disables the KYC process for transactions involving this token.
func (token *SYN10Token) DISABLE_KYC() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    token.KYCEnabled = false
    return token.LOG_TRANSACTION("KYC disabled", "", 0)
}

// VALIDATE_KYC verifies KYC compliance for a given account.
func (token *SYN10Token) VALIDATE_KYC(account string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    if !token.KYCEnabled {
        return true, nil // No KYC required
    }
    
    kycStatus, err := token.Compliance.CheckKYC(account)
    if err != nil || !kycStatus {
        return false, fmt.Errorf("KYC validation failed for account %s: %v", account, err)
    }
    return true, nil
}

// SET_MINTING_POLICY sets the minting policy for this token.
func (token *SYN10Token) SET_MINTING_POLICY(policy string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    token.MintingPolicy = policy
    return token.LOG_TRANSACTION("Minting policy set", "", 0)
}

// GET_MINTING_POLICY retrieves the current minting policy.
func (token *SYN10Token) GET_MINTING_POLICY() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    return token.MintingPolicy
}

// REDEEM_TOKENS redeems a specified amount of tokens, reducing total supply.
func (token *SYN10Token) REDEEM_TOKENS(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    err := token.Ledger.ReduceSupply(account, amount)
    if err != nil {
        return fmt.Errorf("failed to redeem tokens: %v", err)
    }
    
    return token.LOG_TRANSACTION("Tokens redeemed", account, amount)
}

// ISSUE_NEW_TOKENS mints new tokens according to the minting policy.
func (token *SYN10Token) ISSUE_NEW_TOKENS(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    if token.TokenIssuer == "" {
        return fmt.Errorf("token issuer not set")
    }
    
    err := token.Ledger.IncreaseSupply(account, amount)
    if err != nil {
        return fmt.Errorf("failed to issue new tokens: %v", err)
    }
    
    return token.LOG_TRANSACTION("New tokens issued", account, amount)
}

// REVOKE_TOKENS removes a specific amount of tokens from circulation.
func (token *SYN10Token) REVOKE_TOKENS(account string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    err := token.Ledger.BurnTokens(account, amount)
    if err != nil {
        return fmt.Errorf("failed to revoke tokens: %v", err)
    }
    
    return token.LOG_TRANSACTION("Tokens revoked", account, amount)
}

// CHECK_COMPLIANCE_STATUS checks if an account complies with set standards.
func (token *SYN10Token) CHECK_COMPLIANCE_STATUS(account string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    return token.Compliance.CheckCompliance(account)
}

// LOG_TRANSACTION records a transaction activity in the ledger.
func (token *SYN10Token) LOG_TRANSACTION(description, account string, amount uint64) error {
    logEntry := fmt.Sprintf("Description: %s, Account: %s, Amount: %d, Timestamp: %v", description, account, amount, time.Now())
    encryptedLog, err := token.Encryption.Encrypt(logEntry)
    if err != nil {
        return fmt.Errorf("encryption failed for log entry: %v", err)
    }
    
    return token.Ledger.RecordLog("TransactionLog", encryptedLog)
}

// VIEW_TRANSACTION_HISTORY retrieves transaction history for an account.
func (token *SYN10Token) VIEW_TRANSACTION_HISTORY(account string) (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    history, err := token.Ledger.GetTransactionHistory(account)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve transaction history for account %s: %v", account, err)
    }
    return history, nil
}
