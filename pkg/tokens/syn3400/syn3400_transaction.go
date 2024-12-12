package syn3400

import (
	"errors"
	"sync"
	"time"

)

// ForexTransaction represents a transaction involving SYN3400 tokens
type ForexTransaction struct {
	TransactionID   string    // Unique ID of the transaction
	PositionID      string    // The ID of the associated position
	TokenID         string    // The Forex token involved in the transaction
	TransactionType string    // Type of transaction (Buy/Sell/Hedge)
	Amount          float64   // Amount of Forex involved
	Rate            float64   // Rate at which the transaction occurred
	Timestamp       time.Time // Time of the transaction
	TransactionHash string    // Encrypted hash for security
}

// ForexTransactionManager handles all Forex transactions for SYN3400 tokens
type ForexTransactionManager struct {
	transactions map[string]*ForexTransaction
	mutex        sync.Mutex
	ledger       *ledger.Ledger
	encryptor    *encryption.Encryptor
	consensus    *consensus.SynnergyConsensus
}

// NewForexTransactionManager initializes a new instance of ForexTransactionManager
func NewForexTransactionManager(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *ForexTransactionManager {
	return &ForexTransactionManager{
		transactions: make(map[string]*ForexTransaction),
		ledger:       ledger,
		encryptor:    encryptor,
		consensus:    consensus,
	}
}

// CreateTransaction creates a new Forex transaction for SYN3400 tokens
func (ftm *ForexTransactionManager) CreateTransaction(positionID, tokenID, transactionType string, amount, rate float64) (*ForexTransaction, error) {
	ftm.mutex.Lock()
	defer ftm.mutex.Unlock()

	// Validate transaction parameters
	if amount <= 0 || rate <= 0 {
		return nil, errors.New("invalid transaction parameters")
	}

	// Create the transaction
	transactionID := generateUniqueTransactionID()
	timestamp := time.Now()
	transaction := &ForexTransaction{
		TransactionID:   transactionID,
		PositionID:      positionID,
		TokenID:         tokenID,
		TransactionType: transactionType,
		Amount:          amount,
		Rate:            rate,
		Timestamp:       timestamp,
	}

	// Encrypt the transaction
	encryptedTransaction, err := ftm.encryptor.EncryptData(transaction)
	if err != nil {
		return nil, err
	}

	ftm.transactions[transactionID] = encryptedTransaction.(*ForexTransaction)

	// Log the transaction in the ledger
	ftm.ledger.LogEvent("ForexTransactionCreated", timestamp, transactionID)

	// Validate the transaction using Synnergy Consensus
	err = ftm.consensus.ValidateSubBlock(transactionID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransaction retrieves a Forex transaction by its ID
func (ftm *ForexTransactionManager) GetTransaction(transactionID string) (*ForexTransaction, error) {
	ftm.mutex.Lock()
	defer ftm.mutex.Unlock()

	transaction, exists := ftm.transactions[transactionID]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	// Decrypt the transaction before returning it
	decryptedTransaction, err := ftm.encryptor.DecryptData(transaction)
	if err != nil {
		return nil, err
	}

	return decryptedTransaction.(*ForexTransaction), nil
}

// DeleteTransaction deletes a Forex transaction by its ID
func (ftm *ForexTransactionManager) DeleteTransaction(transactionID string) error {
	ftm.mutex.Lock()
	defer ftm.mutex.Unlock()

	if _, exists := ftm.transactions[transactionID]; !exists {
		return errors.New("transaction not found")
	}

	// Remove the transaction from storage
	delete(ftm.transactions, transactionID)

	// Log the deletion in the ledger
	ftm.ledger.LogEvent("ForexTransactionDeleted", time.Now(), transactionID)

	// Validate the deletion using consensus
	return ftm.consensus.ValidateSubBlock(transactionID)
}

// ListTransactions lists all Forex transactions
func (ftm *ForexTransactionManager) ListTransactions() ([]*ForexTransaction, error) {
	ftm.mutex.Lock()
	defer ftm.mutex.Unlock()

	var allTransactions []*ForexTransaction
	for _, transaction := range ftm.transactions {
		// Decrypt each transaction before adding it to the list
		decryptedTransaction, err := ftm.encryptor.DecryptData(transaction)
		if err != nil {
			return nil, err
		}
		allTransactions = append(allTransactions, decryptedTransaction.(*ForexTransaction))
	}

	return allTransactions, nil
}

// UpdateTransaction updates a Forex transaction and logs it in the ledger
func (ftm *ForexTransactionManager) UpdateTransaction(updatedTransaction *ForexTransaction) error {
	ftm.mutex.Lock()
	defer ftm.mutex.Unlock()

	// Encrypt the updated transaction
	encryptedTransaction, err := ftm.encryptor.EncryptData(updatedTransaction)
	if err != nil {
		return err
	}

	ftm.transactions[updatedTransaction.TransactionID] = encryptedTransaction.(*ForexTransaction)

	// Log the update in the ledger
	ftm.ledger.LogEvent("ForexTransactionUpdated", time.Now(), updatedTransaction.TransactionID)

	// Validate the update with consensus
	return ftm.consensus.ValidateSubBlock(updatedTransaction.TransactionID)
}

// generateUniqueTransactionID generates a unique ID for each transaction
func generateUniqueTransactionID() string {
	return time.Now().Format("20060102150405") + "-" + RandomString(10)
}

// RandomString generates a random string of the specified length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
