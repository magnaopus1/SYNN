package syn2600

import (

	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

)

// SYN2600Transaction represents a transaction structure for investor tokens.
type SYN2600Transaction struct {
	TransactionID  string
	TokenID        string
	FromAddress    string
	ToAddress      string
	Amount         float64
	TransactionFee float64
	EncryptedData  string
	Timestamp      time.Time
	Status         string
}

// CreateTransaction creates a new transaction for SYN2600 tokens and validates it with Synnergy Consensus.
func CreateTransaction(tokenID string, fromAddress string, toAddress string, amount float64, fee float64) (*SYN2600Transaction, error) {
	// Check if the token exists in the ledger
	exists, err := TokenExists(tokenID)
	if err != nil || !exists {
		return nil, errors.New("token does not exist")
	}

	// Build the transaction object
	transaction := &SYN2600Transaction{
		TransactionID:  common.GenerateUniqueID(),
		TokenID:        tokenID,
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		Amount:         amount,
		TransactionFee: fee,
		Timestamp:      time.Now(),
		Status:         "Pending",
	}

	// Encrypt the transaction data before further operations
	encryptedData, err := encryption.EncryptTransactionData(transaction)
	if err != nil {
		return nil, errors.New("failed to encrypt transaction data")
	}
	transaction.EncryptedData = encryptedData

	// Record the transaction in the ledger
	err = ledger.RecordTransaction(transaction.TransactionID, encryptedData)
	if err != nil {
		return nil, errors.New("failed to record transaction in the ledger")
	}

	// Validate the transaction through Synnergy Consensus (sub-block validation)
	err = synconsensus.ValidateSubBlockTransaction(transaction.TransactionID, encryptedData)
	if err != nil {
		transaction.Status = "Failed"
		return nil, errors.New("transaction validation failed")
	}

	// Set the transaction status to confirmed
	transaction.Status = "Confirmed"

	// Update the ledger with the confirmed transaction
	err = ledger.UpdateTransactionStatus(transaction.TransactionID, "Confirmed")
	if err != nil {
		return nil, errors.New("failed to update transaction status in ledger")
	}

	return transaction, nil
}

// FetchTransaction retrieves an existing SYN2600 transaction by transaction ID and decrypts it.
func FetchTransaction(transactionID string) (*SYN2600Transaction, error) {
	// Fetch encrypted transaction from the ledger
	encryptedData, err := ledger.FetchTransaction(transactionID)
	if err != nil {
		return nil, errors.New("failed to fetch transaction from ledger")
	}

	// Decrypt the transaction data
	decryptedTransaction, err := encryption.DecryptTransactionData(encryptedData)
	if err != nil {
		return nil, errors.New("failed to decrypt transaction data")
	}

	return decryptedTransaction, nil
}

// CancelTransaction allows the user to cancel a pending transaction before it gets confirmed.
func CancelTransaction(transactionID string) (string, error) {
	// Fetch the transaction from the ledger
	transaction, err := FetchTransaction(transactionID)
	if err != nil {
		return "", errors.New("transaction not found")
	}

	// Ensure the transaction is still pending
	if transaction.Status != "Pending" {
		return "", errors.New("only pending transactions can be cancelled")
	}

	// Update the transaction status to "Cancelled"
	transaction.Status = "Cancelled"

	// Update the ledger with the cancellation
	err = ledger.UpdateTransactionStatus(transaction.TransactionID, "Cancelled")
	if err != nil {
		return "", errors.New("failed to update transaction status in ledger")
	}

	// Record the cancellation event
	eventID, err := RecordSecurityEvent(transaction.TokenID, "CANCEL", "Transaction cancelled", transaction.FromAddress)
	if err != nil {
		return "", errors.New("failed to record cancellation event")
	}

	return eventID, nil
}

// ListTransactionsByToken lists all transactions related to a specific SYN2600 token.
func ListTransactionsByToken(tokenID string) ([]*SYN2600Transaction, error) {
	// Fetch all transactions related to the token from the ledger
	encryptedTransactions, err := ledger.FetchTransactionsByToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to fetch transactions from ledger")
	}

	// Decrypt each transaction
	var transactions []*SYN2600Transaction
	for _, encryptedData := range encryptedTransactions {
		decryptedTransaction, err := encryption.DecryptTransactionData(encryptedData)
		if err != nil {
			return nil, errors.New("failed to decrypt transaction data")
		}
		transactions = append(transactions, decryptedTransaction)
	}

	return transactions, nil
}

// RecordTransactionLog logs transaction events to the ledger for audit purposes.
func RecordTransactionLog(transactionID string, transactionType string, details string) error {
	// Create transaction log
	log := common.TransactionLog{
		TransactionID:  transactionID,
		TransactionType: transactionType,
		Details:        details,
		Timestamp:      time.Now(),
	}

	// Store log in the ledger
	err := ledger.RecordTransactionLog(log)
	if err != nil {
		return errors.New("failed to record transaction log in the ledger")
	}

	return nil
}

// ValidateTransaction ensures the integrity and accuracy of the transaction details.
func ValidateTransaction(transaction *SYN2600Transaction) error {
	// Validate transaction fields
	if transaction.Amount <= 0 || transaction.TransactionFee < 0 {
		return errors.New("invalid transaction details")
	}

	// Ensure both FromAddress and ToAddress are provided
	if transaction.FromAddress == "" || transaction.ToAddress == "" {
		return errors.New("transaction addresses must be provided")
	}

	// Validate transaction status
	if transaction.Status == "" {
		transaction.Status = "Pending"
	}

	// Record a validation log
	err := RecordTransactionLog(transaction.TransactionID, "VALIDATE", "Transaction validated successfully")
	if err != nil {
		return errors.New("failed to log transaction validation")
	}

	return nil
}
