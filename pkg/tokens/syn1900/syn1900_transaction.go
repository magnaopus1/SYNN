package syn1900

import (
	"errors"
	"time"
)


// TokenTransactionService provides transaction-related operations for SYN1900 tokens.
type TokenTransactionService struct {
	transactionStore TransactionStoreInterface // Interface for interacting with the storage backend
	mutex            sync.Mutex                // To handle concurrent access to the transaction system
}

// TransactionStoreInterface defines methods for interacting with the storage layer.
type TransactionStoreInterface interface {
	RecordTransaction(tx common.SYN1900Transaction) error
	GetTransactionByID(transactionID string) (common.SYN1900Transaction, error)
	UpdateTransaction(tx common.SYN1900Transaction) error
	ListAllTransactions() ([]common.SYN1900Transaction, error)
}

// SYN1900Transaction represents an educational credit token transaction.
type SYN1900Transaction struct {
	TransactionID    string    // Unique identifier for the transaction
	TokenID          string    // The ID of the token being transferred or updated
	Sender           string    // ID of the sender (e.g., the original issuer of the education credit)
	Recipient        string    // ID of the recipient (e.g., the student or employee receiving the credit)
	TransactionDate  time.Time // Timestamp of when the transaction occurred
	TransactionType  string    // Type of transaction (e.g., "Issue", "Transfer", "Revoke")
	Amount           float64   // The amount of credit involved in the transaction
	Metadata         string    // Additional metadata (e.g., course information, certificate details)
	EncryptedMetadata []byte   // Encrypted version of the metadata for security
}

// RecordTransaction records a new transaction for a SYN1900 token.
func (s *TokenTransactionService) RecordTransaction(tx common.SYN1900Transaction, encryptionKey string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Encrypt sensitive metadata before recording the transaction
	encryptedMetadata, err := encryption.Encrypt([]byte(tx.Metadata), encryptionKey)
	if err != nil {
		return errors.New("failed to encrypt transaction metadata")
	}
	tx.EncryptedMetadata = encryptedMetadata
	tx.Metadata = "" // Clear the plain metadata for security

	// Record the transaction in storage
	err = s.transactionStore.RecordTransaction(tx)
	if err != nil {
		return errors.New("failed to record transaction in storage")
	}

	return nil
}

// GetTransaction retrieves a transaction by its ID and decrypts its metadata.
func (s *TokenTransactionService) GetTransaction(transactionID, decryptionKey string) (common.SYN1900Transaction, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Retrieve the transaction from the store
	tx, err := s.transactionStore.GetTransactionByID(transactionID)
	if err != nil {
		return common.SYN1900Transaction{}, errors.New("transaction not found in storage")
	}

	// Decrypt the metadata
	decryptedMetadata, err := encryption.Decrypt(tx.EncryptedMetadata, decryptionKey)
	if err != nil {
		return common.SYN1900Transaction{}, errors.New("failed to decrypt transaction metadata")
	}
	tx.Metadata = string(decryptedMetadata)

	return tx, nil
}

// UpdateTransaction updates a transaction's details in the storage backend.
func (s *TokenTransactionService) UpdateTransaction(tx common.SYN1900Transaction, encryptionKey string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Re-encrypt the metadata before updating
	encryptedMetadata, err := encryption.Encrypt([]byte(tx.Metadata), encryptionKey)
	if err != nil {
		return errors.New("failed to encrypt transaction metadata")
	}
	tx.EncryptedMetadata = encryptedMetadata
	tx.Metadata = "" // Clear the plain metadata for security

	// Update the transaction in storage
	err = s.transactionStore.UpdateTransaction(tx)
	if err != nil {
		return errors.New("failed to update transaction in storage")
	}

	return nil
}

// ListAllTransactions retrieves a list of all SYN1900 transactions.
func (s *TokenTransactionService) ListAllTransactions() ([]common.SYN1900Transaction, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Retrieve all transactions from storage
	txs, err := s.transactionStore.ListAllTransactions()
	if err != nil {
		return nil, errors.New("failed to list transactions from storage")
	}

	return txs, nil
}

// ValidateTransaction checks the integrity of a transaction before recording it.
func (s *TokenTransactionService) ValidateTransaction(tx common.SYN1900Transaction) error {
	// Ensure the transaction contains valid data (e.g., non-empty fields, valid IDs)
	if tx.TransactionID == "" || tx.TokenID == "" || tx.Sender == "" || tx.Recipient == "" {
		return errors.New("invalid transaction data: missing essential fields")
	}

	// Additional validation logic can be added here (e.g., checking credit amounts)

	return nil
}

// RevokeTransaction marks a transaction as revoked and updates its status.
func (s *TokenTransactionService) RevokeTransaction(transactionID, reason string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Retrieve the transaction from storage
	tx, err := s.transactionStore.GetTransactionByID(transactionID)
	if err != nil {
		return errors.New("transaction not found in storage")
	}

	// Mark the transaction as revoked
	tx.TransactionType = "Revoke"
	tx.Metadata = reason
	tx.TransactionDate = time.Now()

	// Update the transaction in storage
	err = s.transactionStore.UpdateTransaction(tx)
	if err != nil {
		return errors.New("failed to revoke transaction in storage")
	}

	return nil
}
