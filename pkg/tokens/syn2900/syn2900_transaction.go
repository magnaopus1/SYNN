package syn2900

import (

	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"time"
	"sync"
)

// TransactionManager handles all SYN2900 token-related transactions.
type TransactionManager struct {
	mu sync.Mutex
}

// NewTransactionManager creates a new instance of TransactionManager.
func NewTransactionManager() *TransactionManager {
	return &TransactionManager{}
}

// InitiateTransaction starts a new insurance token transaction between two parties.
func (tm *TransactionManager) InitiateTransaction(sender common.InsuranceTokenOwner, recipient common.InsuranceTokenOwner, token common.InsuranceToken, transferAmount float64) (*common.InsuranceTransaction, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Ensure sender owns the token and has sufficient balance
	if sender.TokenID != token.TokenID || sender.Balance < transferAmount {
		return nil, errors.New("insufficient token balance or invalid token ownership")
	}

	// Create a new insurance transaction
	transaction := &common.InsuranceTransaction{
		TransactionID:  common.GenerateTransactionID(),
		TokenID:        token.TokenID,
		Sender:         sender.OwnerID,
		Recipient:      recipient.OwnerID,
		Amount:         transferAmount,
		Timestamp:      time.Now(),
		Status:         "Pending",
	}

	// Encrypt transaction data before storing
	encryptedData, err := encryptTransactionData(transaction)
	if err != nil {
		return nil, err
	}

	// Store encrypted transaction in the ledger
	err = ledger.Store("transaction_"+transaction.TransactionID, encryptedData)
	if err != nil {
		return nil, err
	}

	// Update balances
	err = tm.updateBalances(sender, recipient, token, transferAmount)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// CompleteTransaction completes the pending insurance transaction.
func (tm *TransactionManager) CompleteTransaction(transactionID string) (*common.InsuranceTransaction, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Retrieve encrypted transaction data from the ledger
	encryptedData, err := ledger.Retrieve("transaction_" + transactionID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	// Decrypt the transaction data
	decryptedData, err := decryptTransactionData(encryptedData)
	if err != nil {
		return nil, err
	}

	// Unmarshal the decrypted data into the InsuranceTransaction struct
	var transaction common.InsuranceTransaction
	err = json.Unmarshal(decryptedData, &transaction)
	if err != nil {
		return nil, err
	}

	// Mark the transaction as completed
	transaction.Status = "Completed"

	// Encrypt and store the updated transaction data
	updatedEncryptedData, err := encryptTransactionData(&transaction)
	if err != nil {
		return nil, err
	}

	err = ledger.Store("transaction_"+transaction.TransactionID, updatedEncryptedData)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// RevertTransaction reverts a failed or invalid insurance transaction.
func (tm *TransactionManager) RevertTransaction(transactionID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Retrieve encrypted transaction data from the ledger
	encryptedData, err := ledger.Retrieve("transaction_" + transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Decrypt the transaction data
	decryptedData, err := decryptTransactionData(encryptedData)
	if err != nil {
		return err
	}

	// Unmarshal the decrypted data into the InsuranceTransaction struct
	var transaction common.InsuranceTransaction
	err = json.Unmarshal(decryptedData, &transaction)
	if err != nil {
		return err
	}

	// Revert balances by restoring the sender's and recipient's original balances
	err = tm.revertBalances(transaction)
	if err != nil {
		return err
	}

	// Mark the transaction as "Reverted"
	transaction.Status = "Reverted"

	// Encrypt and store the updated transaction data
	updatedEncryptedData, err := encryptTransactionData(&transaction)
	if err != nil {
		return err
	}

	err = ledger.Store("transaction_"+transaction.TransactionID, updatedEncryptedData)
	if err != nil {
		return err
	}

	return nil
}

// updateBalances handles balance updates for the sender and recipient during a transaction.
func (tm *TransactionManager) updateBalances(sender common.InsuranceTokenOwner, recipient common.InsuranceTokenOwner, token common.InsuranceToken, transferAmount float64) error {
	// Decrease sender's balance
	sender.Balance -= transferAmount

	// Increase recipient's balance
	recipient.Balance += transferAmount

	// Update the ledger with the new balances
	err := ledger.UpdateBalance(sender.OwnerID, token.TokenID, sender.Balance)
	if err != nil {
		return err
	}

	err = ledger.UpdateBalance(recipient.OwnerID, token.TokenID, recipient.Balance)
	if err != nil {
		return err
	}

	return nil
}

// revertBalances restores the original balances after a transaction is reverted.
func (tm *TransactionManager) revertBalances(transaction common.InsuranceTransaction) error {
	// Retrieve sender and recipient's original balances from the ledger
	senderBalance, err := ledger.GetBalance(transaction.Sender, transaction.TokenID)
	if err != nil {
		return err
	}

	recipientBalance, err := ledger.GetBalance(transaction.Recipient, transaction.TokenID)
	if err != nil {
		return err
	}

	// Reverse the transaction amounts
	senderBalance += transaction.Amount
	recipientBalance -= transaction.Amount

	// Update the ledger with the reverted balances
	err = ledger.UpdateBalance(transaction.Sender, transaction.TokenID, senderBalance)
	if err != nil {
		return err
	}

	err = ledger.UpdateBalance(transaction.Recipient, transaction.TokenID, recipientBalance)
	if err != nil {
		return err
	}

	return nil
}

// encryptTransactionData encrypts transaction data before storage.
func encryptTransactionData(transaction *common.InsuranceTransaction) (string, error) {
	data, err := json.Marshal(transaction)
	if err != nil {
		return "", err
	}

	encryptedData, err := encryptData(data)
	if err != nil {
		return "", err
	}

	return encryptedData, nil
}

// decryptTransactionData decrypts transaction data after retrieval.
func decryptTransactionData(encrypted string) ([]byte, error) {
	return decryptData(encrypted)
}

