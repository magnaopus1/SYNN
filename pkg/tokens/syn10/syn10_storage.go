package syn10

import (
    "errors"
    "sync"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)



// NewStorageManager initializes a new StorageManager instance.
func NewStorageManager(ledgerInstance *ledger.Ledger, consensusInstance *consensus.SynnergyConsensus, encryptionService *encryption.Service) *StorageManager {
    return &StorageManager{
        Ledger:       ledgerInstance,
        Consensus:    consensusInstance,
        balances:     make(map[string]uint64),
        transactions: make(map[string][]Transaction),
        Encryption:   encryptionService,
    }
}

// GetBalance retrieves the token balance for a specific user.
func (sm *StorageManager) GetBalance(userID string) (uint64, error) {
    sm.mutex.RLock()
    defer sm.mutex.RUnlock()

    balance, exists := sm.balances[userID]
    if !exists {
        return 0, errors.New("user balance not found")
    }

    return balance, nil
}

// UpdateBalance updates the token balance for a specific user.
func (sm *StorageManager) UpdateBalance(userID string, amount uint64) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    sm.balances[userID] = amount
    return sm.Ledger.UpdateBalance(userID, amount)
}

// RecordTransaction records a token transaction between two users.
func (sm *StorageManager) RecordTransaction(senderID, receiverID string, amount uint64, encryptedData string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    transactionID := common.GenerateUUID()
    transaction := Transaction{
        TransactionID: transactionID,
        Sender:        senderID,
        Receiver:      receiverID,
        Amount:        amount,
        Timestamp:     time.Now(),
        Status:        "pending",
        EncryptedData: encryptedData,
    }

    sm.transactions[senderID] = append(sm.transactions[senderID], transaction)
    sm.transactions[receiverID] = append(sm.transactions[receiverID], transaction)

    // Update ledger with the new transaction (validated in sub-blocks)
    err := sm.Ledger.AddTransaction(transactionID, senderID, receiverID, amount)
    if err != nil {
        return err
    }

    // Validate the transaction using Synnergy Consensus
    if valid, err := sm.Consensus.ValidateTransaction(transaction); !valid || err != nil {
        transaction.Status = "failed"
        return errors.New("transaction validation failed")
    }

    transaction.Status = "completed"
    return nil
}

// GetTransactionHistory retrieves the transaction history for a specific user.
func (sm *StorageManager) GetTransactionHistory(userID string) ([]Transaction, error) {
    sm.mutex.RLock()
    defer sm.mutex.RUnlock()

    transactions, exists := sm.transactions[userID]
    if !exists {
        return nil, errors.New("transaction history not found")
    }

    return transactions, nil
}

// EncryptTransactionData encrypts sensitive transaction data using AES encryption.
func (sm *StorageManager) EncryptTransactionData(data []byte, key []byte) ([]byte, error) {
    return sm.Encryption.EncryptData(data, key)
}

// DecryptTransactionData decrypts AES-encrypted transaction data.
func (sm *StorageManager) DecryptTransactionData(ciphertext []byte, key []byte) ([]byte, error) {
    return sm.Encryption.DecryptData(ciphertext, key)
}

// MintTokens mints new tokens and credits the central bank's account.
// Only the central bank authority node can mint tokens.
func (sm *StorageManager) MintTokens(issuerID string, amount uint64) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    if !sm.isCentralBankNode(issuerID) {
        return errors.New("only central bank can mint tokens")
    }

    balance, exists := sm.balances[issuerID]
    if !exists {
        balance = 0
    }

    sm.balances[issuerID] = balance + amount
    return sm.Ledger.UpdateBalance(issuerID, sm.balances[issuerID])
}

// BurnTokens burns tokens from the central bank's account.
// Only the central bank authority node can burn tokens.
func (sm *StorageManager) BurnTokens(issuerID string, amount uint64) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    if !sm.isCentralBankNode(issuerID) {
        return errors.New("only central bank can burn tokens")
    }

    balance, exists := sm.balances[issuerID]
    if !exists || balance < amount {
        return errors.New("insufficient balance to burn")
    }

    sm.balances[issuerID] = balance - amount
    return sm.Ledger.UpdateBalance(issuerID, sm.balances[issuerID])
}

// isCentralBankNode checks if the user is a central bank authority node.
func (sm *StorageManager) isCentralBankNode(userID string) bool {
    // This function checks if the user has the Central Bank role from the ledger or security module
    return sm.Ledger.IsCentralBankNode(userID)
}
