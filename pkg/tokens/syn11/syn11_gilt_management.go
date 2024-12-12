package syn11

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SYN11Token represents the core structure for Central Bank Digital Gilt tokens.
type SYN11Token struct {
    TokenID         string
    Metadata        Syn11Metadata
    Issuer          string
    Ledger          *ledger.Ledger
    Consensus       *consensus.SynnergyConsensus
    Encrypted       bool
    Locked          bool
    Wallets         map[string]uint64 // Wallet ID to balance
    CouponRate      float64
    mutex           sync.Mutex
}

// TRANSFER_GILT transfers gilt tokens between wallets with validation.
func (token *SYN11Token) TRANSFER_GILT(fromWallet, toWallet string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Locked {
        return fmt.Errorf("gilt transactions are currently locked")
    }

    if amount > token.Wallets[fromWallet] {
        return fmt.Errorf("insufficient balance in wallet %s", fromWallet)
    }

    if !token.Consensus.ValidateTransaction("TRANSFER") {
        return fmt.Errorf("transfer failed consensus validation")
    }

    token.Wallets[fromWallet] -= amount
    token.Wallets[toWallet] += amount

    return token.Ledger.RecordTransaction("GiltTransfer", fromWallet, toWallet, amount)
}

// CHECK_GILT_BALANCE retrieves the balance of a specified gilt wallet.
func (token *SYN11Token) CHECK_GILT_BALANCE(walletID string) (uint64, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    balance, exists := token.Wallets[walletID]
    if !exists {
        return 0, fmt.Errorf("wallet %s not found", walletID)
    }
    return balance, nil
}

// GET_GILT_METADATA retrieves metadata information for the gilt token.
func (token *SYN11Token) GET_GILT_METADATA() (*Syn11Metadata, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    metadataCopy := token.Metadata
    return &metadataCopy, nil
}

// UPDATE_GILT_METADATA updates metadata for the gilt token.
func (token *SYN11Token) UPDATE_GILT_METADATA(newMetadata Syn11Metadata) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    encryptedMetadata, err := token.Encryption.Encrypt(fmt.Sprintf("%v", newMetadata))
    if err != nil {
        return fmt.Errorf("metadata encryption failed: %v", err)
    }

    token.Metadata.EncryptedData = encryptedMetadata
    token.Metadata = newMetadata

    return token.Ledger.RecordLog("MetadataUpdated", "Gilt metadata updated")
}

// SET_COUPON_RATE sets the coupon rate for the gilt, affecting future payouts.
func (token *SYN11Token) SET_COUPON_RATE(rate float64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.CouponRate = rate
    return token.Ledger.RecordLog("CouponRateSet", fmt.Sprintf("Coupon rate set to %.2f%%", rate))
}

// FETCH_COUPON_RATE retrieves the current coupon rate of the gilt.
func (token *SYN11Token) FETCH_COUPON_RATE() float64 {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.CouponRate
}

// LOCK_GILT locks all gilt-related transactions.
func (token *SYN11Token) LOCK_GILT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Locked = true
    return token.Ledger.RecordLog("GiltLocked", "Gilt transactions locked")
}

// UNLOCK_GILT unlocks all gilt-related transactions.
func (token *SYN11Token) UNLOCK_GILT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Locked = false
    return token.Ledger.RecordLog("GiltUnlocked", "Gilt transactions unlocked")
}

// CREATE_GILT_WALLET creates a new wallet for storing gilt tokens.
func (token *SYN11Token) CREATE_GILT_WALLET(walletID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if _, exists := token.Wallets[walletID]; exists {
        return fmt.Errorf("wallet %s already exists", walletID)
    }
    token.Wallets[walletID] = 0
    return token.Ledger.RecordLog("GiltWalletCreated", fmt.Sprintf("Wallet %s created for gilt", walletID))
}

// DELETE_GILT_WALLET deletes a gilt wallet if the balance is zero.
func (token *SYN11Token) DELETE_GILT_WALLET(walletID string) error {
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
    return token.Ledger.RecordLog("GiltWalletDeleted", fmt.Sprintf("Wallet %s deleted", walletID))
}

// QUERY_GILT_WALLET_STATUS checks if a wallet exists and is active.
func (token *SYN11Token) QUERY_GILT_WALLET_STATUS(walletID string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    _, exists := token.Wallets[walletID]
    return exists, nil
}

// DEPOSIT_GILT_TO_WALLET deposits gilt tokens into a specified wallet.
func (token *SYN11Token) DEPOSIT_GILT_TO_WALLET(walletID string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Locked {
        return fmt.Errorf("gilt transactions are currently locked")
    }

    token.Wallets[walletID] += amount
    return token.Ledger.RecordTransaction("GiltDeposit", "", walletID, amount)
}

// WITHDRAW_GILT_FROM_WALLET withdraws gilt tokens from a specified wallet if balance allows.
func (token *SYN11Token) WITHDRAW_GILT_FROM_WALLET(walletID string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Locked {
        return fmt.Errorf("gilt transactions are currently locked")
    }

    balance := token.Wallets[walletID]
    if amount > balance {
        return fmt.Errorf("insufficient balance in wallet %s", walletID)
    }

    token.Wallets[walletID] -= amount
    return token.Ledger.RecordTransaction("GiltWithdrawal", walletID, "", amount)
}
