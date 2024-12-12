package syn10

import (
    "sync"
    "time"
    "fmt"
    "path/to/ledger"
    "path/to/synnergy_consensus"
    "path/to/encryption"
)


// TRANSFER_TOKEN transfers tokens between wallets with full validation.
func (token *SYN10Token) TRANSFER_TOKEN(fromWallet, toWallet string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    if token.Locked {
        return fmt.Errorf("token is currently locked for transfers")
    }
    
    if amount > token.Wallets[fromWallet] {
        return fmt.Errorf("insufficient balance in wallet %s", fromWallet)
    }
    
    if !token.Consensus.ValidateSubBlock("TRANSFER") {
        return fmt.Errorf("transfer failed consensus validation")
    }

    token.Wallets[fromWallet] -= amount
    token.Wallets[toWallet] += amount
    
    return token.Ledger.RecordTransaction("Transfer", fromWallet, toWallet, amount)
}

// CHECK_BALANCE retrieves the balance of a specified wallet.
func (token *SYN10Token) CHECK_BALANCE(walletID string) (uint64, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    balance, exists := token.Wallets[walletID]
    if !exists {
        return 0, fmt.Errorf("wallet not found")
    }
    return balance, nil
}

// GET_TOKEN_METADATA retrieves metadata for the token.
func (token *SYN10Token) GET_TOKEN_METADATA() (*SYN10Metadata, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    metadataCopy := *token.Metadata
    return &metadataCopy, nil
}

// UPDATE_TOKEN_METADATA updates the metadata for the token.
func (token *SYN10Token) UPDATE_TOKEN_METADATA(newMetadata SYN10Metadata) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    encryptedMetadata, err := token.Encryption.Encrypt(fmt.Sprintf("%v", newMetadata))
    if err != nil {
        return fmt.Errorf("metadata encryption failed: %v", err)
    }
    
    token.Metadata.EncryptedData = encryptedMetadata
    token.Metadata = &newMetadata
    return nil
}

// SET_EXCHANGE_RATE updates the exchange rate for the token.
func (token *SYN10Token) SET_EXCHANGE_RATE(rate float64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    token.ExchangeRate = rate
    return token.Ledger.RecordLog("ExchangeRateUpdate", fmt.Sprintf("Exchange rate set to %f", rate))
}

// FETCH_EXCHANGE_RATE retrieves the current exchange rate for the token.
func (token *SYN10Token) FETCH_EXCHANGE_RATE() float64 {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    return token.ExchangeRate
}

// LOCK_TOKEN locks the token to prevent transfers.
func (token *SYN10Token) LOCK_TOKEN() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    token.Locked = true
    return token.Ledger.RecordLog("TokenLock", "Token is locked for transfers")
}

// UNLOCK_TOKEN unlocks the token to allow transfers.
func (token *SYN10Token) UNLOCK_TOKEN() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    token.Locked = false
    return token.Ledger.RecordLog("TokenUnlock", "Token is unlocked for transfers")
}

// CREATE_WALLET initializes a new wallet with a zero balance.
func (token *SYN10Token) CREATE_WALLET(walletID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    if _, exists := token.Wallets[walletID]; exists {
        return fmt.Errorf("wallet %s already exists", walletID)
    }
    token.Wallets[walletID] = 0
    return token.Ledger.RecordLog("WalletCreation", fmt.Sprintf("Wallet %s created", walletID))
}

// DELETE_WALLET removes a wallet from the system if balance is zero.
func (token *SYN10Token) DELETE_WALLET(walletID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    balance, exists := token.Wallets[walletID]
    if !exists {
        return fmt.Errorf("wallet %s not found", walletID)
    }
    
    if balance != 0 {
        return fmt.Errorf("wallet %s must have zero balance to be deleted", walletID)
    }
    
    delete(token.Wallets, walletID)
    return token.Ledger.RecordLog("WalletDeletion", fmt.Sprintf("Wallet %s deleted", walletID))
}

// QUERY_WALLET_STATUS checks if a wallet exists and is active.
func (token *SYN10Token) QUERY_WALLET_STATUS(walletID string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    _, exists := token.Wallets[walletID]
    return exists, nil
}

// DEPOSIT_TO_WALLET adds tokens to a specified wallet.
func (token *SYN10Token) DEPOSIT_TO_WALLET(walletID string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    if token.Locked {
        return fmt.Errorf("token is currently locked for deposits")
    }
    
    token.Wallets[walletID] += amount
    return token.Ledger.RecordTransaction("Deposit", walletID, "", amount)
}

// WITHDRAW_FROM_WALLET deducts tokens from a specified wallet if sufficient balance exists.
func (token *SYN10Token) WITHDRAW_FROM_WALLET(walletID string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()
    
    if token.Locked {
        return fmt.Errorf("token is currently locked for withdrawals")
    }
    
    balance := token.Wallets[walletID]
    if amount > balance {
        return fmt.Errorf("insufficient funds in wallet %s", walletID)
    }
    
    token.Wallets[walletID] -= amount
    return token.Ledger.RecordTransaction("Withdrawal", walletID, "", amount)
}
