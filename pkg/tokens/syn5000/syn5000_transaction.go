package syn4900

import (
	"errors"
	"sync"
	"time"
)

// SYN5000Transaction represents a transaction involving the SYN5000 gambling token.
type SYN5000Transaction struct {
	TransactionID    string    // Unique identifier for the transaction.
	TokenID          string    // Identifier of the token involved in the transaction.
	From             string    // Address of the sender.
	To               string    // Address of the recipient.
	Amount           float64   // Amount of tokens transferred.
	TransactionType  string    // Type of transaction (e.g., "bet", "win", "loss", "transfer").
	Timestamp        time.Time // Time of the transaction.
	TransactionHash  string    // Secure hash of the transaction for validation.
	SubBlockHash     string    // Hash of the sub-block for validation.
	BlockHash        string    // Hash of the full block for validation.
	Status           string    // Status of the transaction (pending/complete).
}

// TransactionType constants.
const (
	TransactionBet      = "bet"
	TransactionWin      = "win"
	TransactionLoss     = "loss"
	TransactionTransfer = "transfer"
)

// TransactionManager manages SYN5000 token transactions.
type TransactionManager struct {
	mu            sync.RWMutex
	transactions  map[string]*SYN5000Transaction // In-memory storage of transactions.
	ledger        *ledger.TransactionLedger      // Integration with the ledger for transaction storage.
	security      *encryption.Security           // Encryption and hashing for secure transaction storage.
	consensus     *consensus.SynnergyConsensus   // Synnergy Consensus for validating transactions.
	subBlockCount int                            // Number of sub-blocks validated.
}

// NewTransactionManager initializes a new TransactionManager.
func NewTransactionManager(ledger *ledger.TransactionLedger, security *encryption.Security, consensus *consensus.SynnergyConsensus) *TransactionManager {
	return &TransactionManager{
		transactions:  make(map[string]*SYN5000Transaction),
		ledger:        ledger,
		security:      security,
		consensus:     consensus,
		subBlockCount: 0,
	}
}

// CreateTransaction creates a new transaction involving SYN5000 tokens.
func (tm *TransactionManager) CreateTransaction(tokenID, from, to string, amount float64, transactionType string) (*SYN5000Transaction, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Validate input.
	if amount <= 0 {
		return nil, errors.New("invalid transaction amount")
	}

	// Generate unique transaction ID and hash.
	transactionID := generateUniqueTransactionID()
	timestamp := time.Now()
	transactionHash := tm.security.GenerateHash(fmt.Sprintf("%s%s%s%f%s%s", tokenID, from, to, amount, transactionType, timestamp.String()))

	// Create the transaction object.
	transaction := &SYN5000Transaction{
		TransactionID:   transactionID,
		TokenID:         tokenID,
		From:            from,
		To:              to,
		Amount:          amount,
		TransactionType: transactionType,
		Timestamp:       timestamp,
		TransactionHash: transactionHash,
		Status:          "pending",
	}

	// Store the transaction in-memory.
	tm.transactions[transactionID] = transaction

	// Validate the transaction using Synnergy Consensus.
	subBlockHash, err := tm.consensus.ValidateSubBlock(transactionHash)
	if err != nil {
		return nil, err
	}
	transaction.SubBlockHash = subBlockHash

	// Increment sub-block count.
	tm.subBlockCount++
	if tm.subBlockCount == 1000 {
		blockHash, err := tm.consensus.ValidateBlock(transaction.SubBlockHash)
		if err != nil {
			return nil, err
		}
		transaction.BlockHash = blockHash
		tm.subBlockCount = 0 // Reset sub-block count after a block is validated.
	}

	// Mark the transaction as complete and update the ledger.
	transaction.Status = "complete"
	tm.ledger.StoreTransaction(transactionID, transaction)

	return transaction, nil
}

// GetTransaction retrieves a transaction by its ID.
func (tm *TransactionManager) GetTransaction(transactionID string) (*SYN5000Transaction, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	transaction, exists := tm.transactions[transactionID]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	return transaction, nil
}

// generateUniqueTransactionID generates a unique identifier for transactions using SHA-256.
func generateUniqueTransactionID() string {
	return hex.EncodeToString(sha256.New().Sum([]byte(fmt.Sprintf("%d", time.Now().UnixNano()))))
}

// ConfirmTransaction allows confirmation of the transaction, ensuring full block validation.
func (tm *TransactionManager) ConfirmTransaction(transactionID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	transaction, exists := tm.transactions[transactionID]
	if !exists {
		return errors.New("transaction not found")
	}

	if transaction.Status == "complete" {
		return errors.New("transaction already confirmed")
	}

	// Final validation before confirming the transaction.
	blockHash, err := tm.consensus.ValidateBlock(transaction.SubBlockHash)
	if err != nil {
		return err
	}
	transaction.BlockHash = blockHash

	// Mark the transaction as confirmed and store it in the ledger.
	transaction.Status = "confirmed"
	tm.ledger.StoreTransaction(transactionID, transaction)

	return nil
}

// ValidateTransaction verifies the integrity of a transaction using its hash.
func (tm *TransactionManager) ValidateTransaction(transactionID, expectedHash string) (bool, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	transaction, exists := tm.transactions[transactionID]
	if !exists {
		return false, errors.New("transaction not found")
	}

	// Validate the transaction hash.
	if transaction.TransactionHash != expectedHash {
		return false, errors.New("transaction hash mismatch")
	}

	return true, nil
}
