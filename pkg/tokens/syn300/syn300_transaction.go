package syn300

import (
	"errors"
	"sync"
	"time"
)

// TransactionType defines the types of possible governance token transactions
type TransactionType string

const (
	Transfer      TransactionType = "Transfer"
	Delegate      TransactionType = "Delegate"
	Vote          TransactionType = "Vote"
	Mint          TransactionType = "Mint"
	Burn          TransactionType = "Burn"
)

// GovernanceTransaction represents a transaction in the SYN300 governance token system
type GovernanceTransaction struct {
	ID              string            `json:"id"`         // Transaction ID
	TokenID         string            `json:"token_id"`   // ID of the governance token being transacted
	From            string            `json:"from"`       // Sender's address
	To              string            `json:"to"`         // Receiver's address (for transfer or delegation)
	Amount          uint64            `json:"amount"`     // Amount involved in the transaction
	TransactionType TransactionType   `json:"type"`       // Type of transaction (e.g., transfer, delegate, vote)
	Signature       string            `json:"signature"`  // Cryptographic signature
	Encrypted       bool              `json:"encrypted"`  // Flag for encrypted transaction details
	Timestamp       time.Time         `json:"timestamp"`  // Time of the transaction
	Validated       bool              `json:"validated"`  // Validation status
}

// Syn300TransactionManager manages governance token transactions
type Syn300TransactionManager struct {
	Ledger       *ledger.Ledger
	Transactions map[string]GovernanceTransaction
	mutex        sync.RWMutex
}

// NewSyn300TransactionManager creates a new transaction manager for SYN300 governance tokens
func NewSyn300TransactionManager(ledger *ledger.Ledger) *Syn300TransactionManager {
	return &Syn300TransactionManager{
		Ledger:       ledger,
		Transactions: make(map[string]GovernanceTransaction),
	}
}

// CreateTransaction creates a new governance token transaction
func (tm *Syn300TransactionManager) CreateTransaction(tokenID string, from string, to string, amount uint64, txType TransactionType, signature string) (string, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	txID := generateTransactionID()

	// Encrypt transaction details for security
	encryptedFrom, err := encryption.Encrypt(from)
	if err != nil {
		return "", errors.New("failed to encrypt sender address")
	}
	encryptedTo, err := encryption.Encrypt(to)
	if err != nil {
		return "", errors.New("failed to encrypt receiver address")
	}

	transaction := GovernanceTransaction{
		ID:              txID,
		TokenID:         tokenID,
		From:            encryptedFrom,
		To:              encryptedTo,
		Amount:          amount,
		TransactionType: txType,
		Signature:       signature,
		Encrypted:       true,
		Timestamp:       time.Now(),
		Validated:       false,
	}

	// Store the transaction in the local map
	tm.Transactions[txID] = transaction

	// Log the transaction in the ledger
	if err := tm.Ledger.StoreTransaction(transaction); err != nil {
		return "", errors.New("failed to store transaction in the ledger")
	}

	return txID, nil
}

// ValidateTransaction validates a transaction using Synnergy Consensus
func (tm *Syn300TransactionManager) ValidateTransaction(txID string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	transaction, exists := tm.Transactions[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	if transaction.Validated {
		return errors.New("transaction already validated")
	}

	// Validate transaction using Synnergy Consensus
	if err := consensus.ValidateTransaction(txID, transaction.Signature); err != nil {
		return errors.New("failed to validate transaction under Synnergy Consensus")
	}

	// Mark transaction as validated
	transaction.Validated = true
	tm.Transactions[txID] = transaction

	// Update the ledger with the validated transaction
	if err := tm.Ledger.UpdateTransaction(transaction); err != nil {
		return errors.New("failed to update transaction in the ledger")
	}

	return nil
}

// GetTransaction retrieves a transaction by its ID
func (tm *Syn300TransactionManager) GetTransaction(txID string) (GovernanceTransaction, error) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	transaction, exists := tm.Transactions[txID]
	if !exists {
		return GovernanceTransaction{}, errors.New("transaction not found")
	}

	// Decrypt transaction details before returning
	decryptedFrom, err := encryption.Decrypt(transaction.From)
	if err != nil {
		return GovernanceTransaction{}, errors.New("failed to decrypt sender address")
	}
	decryptedTo, err := encryption.Decrypt(transaction.To)
	if err != nil {
		return GovernanceTransaction{}, errors.New("failed to decrypt receiver address")
	}

	transaction.From = decryptedFrom
	transaction.To = decryptedTo
	return transaction, nil
}

// GetValidatedTransactions retrieves all validated transactions
func (tm *Syn300TransactionManager) GetValidatedTransactions() ([]GovernanceTransaction, error) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	var validatedTransactions []GovernanceTransaction
	for _, transaction := range tm.Transactions {
		if transaction.Validated {
			// Decrypt sender and receiver addresses for each validated transaction
			decryptedFrom, err := encryption.Decrypt(transaction.From)
			if err != nil {
				return nil, errors.New("failed to decrypt sender address")
			}
			decryptedTo, err := encryption.Decrypt(transaction.To)
			if err != nil {
				return nil, errors.New("failed to decrypt receiver address")
			}
			transaction.From = decryptedFrom
			transaction.To = decryptedTo
			validatedTransactions = append(validatedTransactions, transaction)
		}
	}
	return validatedTransactions, nil
}



// Helper function to generate unique transaction IDs
func generateTransactionID() string {
	// Placeholder for a proper unique ID generator
	return "tx_" + time.Now().Format("20060102150405")
}
