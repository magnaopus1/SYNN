package syn1000

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

// SYN1000Transaction represents a transaction involving SYN1000 tokens
type SYN1000Transaction struct {
	TransactionID string            `json:"transaction_id"`
	TokenID       string            `json:"token_id"`
	From          string            `json:"from"`
	To            string            `json:"to"`
	Amount        float64           `json:"amount"`
	Timestamp     time.Time         `json:"timestamp"`
	ExtraData     map[string]string `json:"extra_data"`
}

// SYN1000TransactionManager manages the creation, validation, and storage of SYN1000 transactions
type SYN1000TransactionManager struct {
	transactions   map[string]string // TransactionID -> Encrypted Transaction Data
	mutex          sync.RWMutex
	Ledger         *ledger.Ledger
	Consensus      *consensus.SynnergyConsensus
	Encryption     *encryption.EncryptionService
}

// NewSYN1000TransactionManager initializes a new transaction manager for SYN1000 tokens
func NewSYN1000TransactionManager(ledger *ledger.Ledger, consensus *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *SYN1000TransactionManager {
	return &SYN1000TransactionManager{
		transactions: make(map[string]string),
		Ledger:       ledger,
		Consensus:    consensus,
		Encryption:   encryptionService,
	}
}

// CreateTransaction creates a new transaction and stores it securely
func (tm *SYN1000TransactionManager) CreateTransaction(tokenID, from, to string, amount float64, extraData map[string]string) (string, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Generate a new transaction ID
	transactionID := generateUniqueID()

	// Create the transaction object
	tx := SYN1000Transaction{
		TransactionID: transactionID,
		TokenID:       tokenID,
		From:          from,
		To:            to,
		Amount:        amount,
		Timestamp:     time.Now(),
		ExtraData:     extraData,
	}

	// Validate the transaction using Synnergy Consensus
	if err := tm.Consensus.ValidateTransaction(&tx); err != nil {
		return "", errors.New("transaction validation failed via Synnergy Consensus")
	}

	// Encrypt the transaction data
	txData, err := json.Marshal(tx)
	if err != nil {
		return "", errors.New("failed to serialize transaction")
	}

	encryptedTx, encryptionKey, err := tm.Encryption.EncryptData(txData)
	if err != nil {
		return "", errors.New("failed to encrypt transaction data")
	}

	// Store the encrypted transaction in memory
	tm.transactions[transactionID] = encryptedTx

	// Record the transaction in the ledger
	if err := tm.Ledger.StoreTransaction(transactionID, encryptedTx, encryptionKey); err != nil {
		return "", errors.New("failed to store transaction in ledger")
	}

	return transactionID, nil
}

// RetrieveTransaction retrieves a transaction by its ID and decrypts it
func (tm *SYN1000TransactionManager) RetrieveTransaction(transactionID string) (*SYN1000Transaction, error) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	// Get encrypted transaction data from memory
	encryptedTx, exists := tm.transactions[transactionID]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	// Retrieve encryption key from the ledger
	encryptionKey, err := tm.Ledger.GetEncryptionKey(transactionID)
	if err != nil {
		return nil, errors.New("failed to retrieve encryption key from ledger")
	}

	// Decrypt the transaction data
	decryptedTx, err := tm.Encryption.DecryptData(encryptedTx, encryptionKey)
	if err != nil {
		return nil, errors.New("failed to decrypt transaction data")
	}

	var tx SYN1000Transaction
	if err := json.Unmarshal(decryptedTx, &tx); err != nil {
		return nil, errors.New("failed to deserialize transaction")
	}

	return &tx, nil
}

// ValidateTransaction checks if a transaction is valid based on business rules
func (tm *SYN1000TransactionManager) ValidateTransaction(tx *SYN1000Transaction) error {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	// Ensure the transaction amount is greater than zero
	if tx.Amount <= 0 {
		return errors.New("transaction amount must be greater than zero")
	}

	// Check that the sender and receiver are not the same
	if tx.From == tx.To {
		return errors.New("transaction sender and receiver must be different")
	}

	// Use the Synnergy Consensus to validate the transaction further
	if err := tm.Consensus.ValidateTransaction(tx); err != nil {
		return errors.New("Synnergy Consensus validation failed")
	}

	return nil
}

// GetTransactionHistory retrieves the transaction history for a specific token, with encryption
func (tm *SYN1000TransactionManager) GetTransactionHistory(tokenID string) ([]SYN1000Transaction, error) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	// Retrieve transaction logs from the ledger
	encryptedLogs, err := tm.Ledger.GetTransactionLogs(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve transaction logs from ledger")
	}

	// Decrypt the transaction logs
	var history []SYN1000Transaction
	for _, encryptedLog := range encryptedLogs {
		decryptedLog, err := tm.Encryption.DecryptData(encryptedLog, "")
		if err != nil {
			return nil, err
		}

		var tx SYN1000Transaction
		if err := json.Unmarshal(decryptedLog, &tx); err != nil {
			return nil, err
		}

		history = append(history, tx)
	}

	return history, nil
}

// generateUniqueID generates a unique ID for transactions
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
