package syn12

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// Syn12Token represents the core structure for Treasury Bill tokens.
type Syn12Token struct {
    TokenID         string
    Metadata        Syn12Metadata
    Issuer          string
    Ledger          *ledger.Ledger
    Consensus       *consensus.SynnergyConsensus
    DiscountRate    float64
    Locked          bool
    Wallets         map[string]uint64 // Wallet ID to balance
    mutex           sync.Mutex
}

// TRANSFER_TBILL transfers T-Bill tokens from one wallet to another.
func (token *Syn12Token) TRANSFER_TBILL(fromWallet, toWallet string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Locked {
        return fmt.Errorf("T-Bill transactions are currently locked")
    }

    if amount > token.Wallets[fromWallet] {
        return fmt.Errorf("insufficient balance in wallet %s", fromWallet)
    }

    if !token.Consensus.ValidateTransaction("TRANSFER") {
        return fmt.Errorf("transfer failed consensus validation")
    }

    token.Wallets[fromWallet] -= amount
    token.Wallets[toWallet] += amount

    return token.Ledger.RecordTransaction("TBillTransfer", fromWallet, toWallet, amount)
}

// CHECK_TBILL_BALANCE retrieves the T-Bill balance of a specified wallet.
func (token *Syn12Token) CHECK_TBILL_BALANCE(walletID string) (uint64, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    balance, exists := token.Wallets[walletID]
    if !exists {
        return 0, fmt.Errorf("wallet %s not found", walletID)
    }
    return balance, nil
}

// GET_TBILL_METADATA retrieves metadata associated with the T-Bill token.
func (token *Syn12Token) GET_TBILL_METADATA() (*Syn12Metadata, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    metadataCopy := token.Metadata
    return &metadataCopy, nil
}

// UPDATE_TBILL_METADATA updates the metadata for the T-Bill token.
func (token *Syn12Token) UPDATE_TBILL_METADATA(newMetadata Syn12Metadata) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    encryptedMetadata, err := token.Encryption.Encrypt(fmt.Sprintf("%v", newMetadata))
    if err != nil {
        return fmt.Errorf("metadata encryption failed: %v", err)
    }

    token.Metadata = newMetadata
    token.Metadata.EncryptedData = encryptedMetadata

    return token.Ledger.RecordLog("MetadataUpdated", "T-Bill metadata updated")
}

// SET_DISCOUNT_RATE sets the discount rate for the T-Bill token.
func (token *Syn12Token) SET_DISCOUNT_RATE(rate float64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.DiscountRate = rate
    return token.Ledger.RecordLog("DiscountRateSet", fmt.Sprintf("Discount rate set to %.2f%%", rate))
}

// FETCH_DISCOUNT_RATE retrieves the current discount rate for the T-Bill token.
func (token *Syn12Token) FETCH_DISCOUNT_RATE() float64 {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.DiscountRate
}

// LOCK_TBILL locks all T-Bill transactions, preventing transfers.
func (token *Syn12Token) LOCK_TBILL() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Locked = true
    return token.Ledger.RecordLog("TBillLocked", "T-Bill transactions locked")
}

// UNLOCK_TBILL unlocks T-Bill transactions, allowing transfers.
func (token *Syn12Token) UNLOCK_TBILL() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Locked = false
    return token.Ledger.RecordLog("TBillUnlocked", "T-Bill transactions unlocked")
}

// CREATE_TBILL_WALLET creates a new wallet for holding T-Bill tokens.
func (token *Syn12Token) CREATE_TBILL_WALLET(walletID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if _, exists := token.Wallets[walletID]; exists {
        return fmt.Errorf("wallet %s already exists", walletID)
    }

    token.Wallets[walletID] = 0
    return token.Ledger.RecordLog("TBillWalletCreated", fmt.Sprintf("Wallet %s created for T-Bill", walletID))
}

// DELETE_TBILL_WALLET deletes a T-Bill wallet if its balance is zero.
func (token *Syn12Token) DELETE_TBILL_WALLET(walletID string) error {
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
    return token.Ledger.RecordLog("TBillWalletDeleted", fmt.Sprintf("Wallet %s deleted", walletID))
}

// QUERY_TBILL_WALLET_STATUS checks if a wallet exists and is active for T-Bill operations.
func (token *Syn12Token) QUERY_TBILL_WALLET_STATUS(walletID string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    _, exists := token.Wallets[walletID]
    return exists, nil
}

// DEPOSIT_TBILL_TO_WALLET deposits T-Bill tokens into a specified wallet.
func (token *Syn12Token) DEPOSIT_TBILL_TO_WALLET(walletID string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Locked {
        return fmt.Errorf("T-Bill transactions are currently locked")
    }

    token.Wallets[walletID] += amount
    return token.Ledger.RecordTransaction("TBillDeposit", "", walletID, amount)
}

// WITHDRAW_TBILL_FROM_WALLET withdraws T-Bill tokens from a specified wallet if balance allows.
func (token *Syn12Token) WITHDRAW_TBILL_FROM_WALLET(walletID string, amount uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Locked {
        return fmt.Errorf("T-Bill transactions are currently locked")
    }

    balance := token.Wallets[walletID]
    if amount > balance {
        return fmt.Errorf("insufficient balance in wallet %s", walletID)
    }

    token.Wallets[walletID] -= amount
    return token.Ledger.RecordTransaction("TBillWithdrawal", walletID, "", amount)
}
