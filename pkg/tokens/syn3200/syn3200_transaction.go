package syn3200

import (
	"time"
	"errors"
	"sync"

)

// TransactionType represents the type of bill transaction.
type TransactionType string

const (
	Payment        TransactionType = "Payment"
	Refund         TransactionType = "Refund"
	Fee            TransactionType = "Fee"
	Adjustment     TransactionType = "Adjustment"
	Automated      TransactionType = "Automated Payment"
	Fractional     TransactionType = "Fractional Payment"
)

// BillTransaction represents a transaction associated with a SYN3200 bill token.
type BillTransaction struct {
	TransactionID   string          `json:"transaction_id"`
	BillID          string          `json:"bill_id"`
	TransactionType TransactionType `json:"transaction_type"`
	Amount          float64         `json:"amount"`
	Timestamp       time.Time       `json:"timestamp"`
	Payer           string          `json:"payer"`
	Receiver        string          `json:"receiver"`
	Status          string          `json:"status"` // Pending, Completed, Failed
}

// TransactionManager manages all SYN3200 bill token transactions.
type TransactionManager struct {
	transactions     map[string]*BillTransaction
	ledgerService    *ledger.Ledger
	encryptionService *encryption.Encryptor
	consensusService *consensus.SynnergyConsensus
	mutex            sync.Mutex
}

// NewTransactionManager creates a new instance of TransactionManager.
func NewTransactionManager(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *TransactionManager {
	return &TransactionManager{
		transactions:     make(map[string]*BillTransaction),
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// CreateTransaction creates a new transaction for a bill and logs it in the ledger.
func (tm *TransactionManager) CreateTransaction(billID string, txType TransactionType, amount float64, payer, receiver string) (*BillTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Generate a unique transaction ID.
	transactionID := generateTransactionID()

	// Create the transaction.
	transaction := &BillTransaction{
		TransactionID:   transactionID,
		BillID:          billID,
		TransactionType: txType,
		Amount:          amount,
		Timestamp:       time.Now(),
		Payer:           payer,
		Receiver:        receiver,
		Status:          "Pending",
	}

	// Encrypt the transaction before storing it.
	encryptedTransaction, err := tm.encryptionService.EncryptData(transaction)
	if err != nil {
		return nil, err
	}

	// Store the encrypted transaction.
	tm.transactions[transactionID] = encryptedTransaction.(*BillTransaction)

	// Log the transaction creation in the ledger.
	if err := tm.ledgerService.LogEvent("TransactionCreated", time.Now(), transactionID); err != nil {
		return nil, err
	}

	// Validate the transaction with Synnergy Consensus.
	if err := tm.consensusService.ValidateSubBlock(transactionID); err != nil {
		return nil, err
	}

	// Mark the transaction as completed after consensus validation.
	transaction.Status = "Completed"
	return transaction, nil
}

// RetrieveTransaction retrieves a transaction by its ID.
func (tm *TransactionManager) RetrieveTransaction(transactionID string) (*BillTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the transaction from storage.
	transaction, exists := tm.transactions[transactionID]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	// Decrypt the transaction.
	decryptedTransaction, err := tm.encryptionService.DecryptData(transaction)
	if err != nil {
		return nil, err
	}

	return decryptedTransaction.(*BillTransaction), nil
}

// UpdateTransactionStatus updates the status of a transaction and logs the update in the ledger.
func (tm *TransactionManager) UpdateTransactionStatus(transactionID, status string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the transaction.
	transaction, exists := tm.transactions[transactionID]
	if !exists {
		return errors.New("transaction not found")
	}

	// Update the transaction status.
	transaction.Status = status

	// Encrypt the updated transaction.
	encryptedTransaction, err := tm.encryptionService.EncryptData(transaction)
	if err != nil {
		return err
	}

	// Update the transaction in storage.
	tm.transactions[transactionID] = encryptedTransaction.(*BillTransaction)

	// Log the status update in the ledger.
	if err := tm.ledgerService.LogEvent("TransactionStatusUpdated", time.Now(), transactionID); err != nil {
		return err
	}

	// Validate the status update with Synnergy Consensus.
	if err := tm.consensusService.ValidateSubBlock(transactionID); err != nil {
		return err
	}

	return nil
}

// ProcessAutomatedPayment processes an automated payment and validates it.
func (tm *TransactionManager) ProcessAutomatedPayment(billID string, amount float64, payer, receiver string) (*BillTransaction, error) {
	return tm.CreateTransaction(billID, Automated, amount, payer, receiver)
}

// ProcessFractionalPayment processes a fractional payment for a bill.
func (tm *TransactionManager) ProcessFractionalPayment(billID string, amount float64, payer, receiver string) (*BillTransaction, error) {
	return tm.CreateTransaction(billID, Fractional, amount, payer, receiver)
}

// GenerateTransactionID creates a unique transaction ID for each transaction.
func generateTransactionID() string {
	// Implement logic to generate a unique transaction ID (e.g., using timestamp, random strings, etc.)
	return "TX" + time.Now().Format("20060102150405") + "_" + generateRandomString(6)
}

// Helper function to generate random string for transaction ID.
func generateRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[i%len(letterBytes)]
	}
	return string(b)
}
