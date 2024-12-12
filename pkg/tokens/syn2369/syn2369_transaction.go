package syn2369

import (
	"time"
	"errors"
)


// InitiateTransaction starts a transaction to transfer ownership of a SYN2369 token.
func InitiateTransaction(tokenID, fromOwner, toOwner string, customAttributes map[string]interface{}) (string, error) {
	// Retrieve the token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return "", errors.New("token not found: " + err.Error())
	}

	// Ensure the sender owns the token
	if token.Owner != fromOwner {
		return "", errors.New("ownership validation failed: token does not belong to the sender")
	}

	// Validate the transaction through Synnergy Consensus
	err = synnergy.ValidateTransaction(tokenID, fromOwner, toOwner)
	if err != nil {
		return "", errors.New("transaction validation failed: " + err.Error())
	}

	// Encrypt the custom attributes (if any) before processing the transaction
	encryptedAttributes, err := encryption.EncryptData(customAttributes)
	if err != nil {
		return "", errors.New("failed to encrypt custom attributes: " + err.Error())
	}

	// Create a transaction record
	tx := common.SYN2369Transaction{
		TokenID:           tokenID,
		FromOwner:         fromOwner,
		ToOwner:           toOwner,
		TransactionTime:   time.Now(),
		EncryptedMetadata: encryptedAttributes,
	}

	// Store the transaction in the ledger
	txID, err := ledger.StoreTransaction(tx)
	if err != nil {
		return "", errors.New("failed to store transaction in ledger: " + err.Error())
	}

	// Update token ownership in the ledger
	token.Owner = toOwner
	err = ledger.UpdateToken(token)
	if err != nil {
		return "", errors.New("failed to update token ownership: " + err.Error())
	}

	// Log the transaction event
	err = LogTransactionEvent(txID, "TokenTransferred", "Token ID "+tokenID+" transferred from "+fromOwner+" to "+toOwner)
	if err != nil {
		return "", err
	}

	return txID, nil
}

// FinalizeTransaction finalizes the pending transaction after consensus validation and applies the ownership change.
func FinalizeTransaction(txID string) error {
	// Retrieve the transaction from the ledger
	tx, err := ledger.GetTransaction(txID)
	if err != nil {
		return errors.New("transaction not found: " + err.Error())
	}

	// Validate the transaction through Synnergy Consensus
	err = synnergy.ValidateTransactionFinalization(txID)
	if err != nil {
		return errors.New("transaction finalization validation failed: " + err.Error())
	}

	// Mark the transaction as complete in the ledger
	tx.Completed = true
	err = ledger.UpdateTransaction(tx)
	if err != nil {
		return errors.New("failed to finalize transaction: " + err.Error())
	}

	// Log the event of transaction finalization
	err = LogTransactionEvent(txID, "TransactionFinalized", "Transaction ID "+txID+" has been finalized.")
	if err != nil {
		return err
	}

	return nil
}

// CancelTransaction cancels a pending transaction before it is finalized.
func CancelTransaction(txID string) error {
	// Retrieve the transaction from the ledger
	tx, err := ledger.GetTransaction(txID)
	if err != nil {
		return errors.New("transaction not found: " + err.Error())
	}

	// Validate if the transaction can be canceled
	if tx.Completed {
		return errors.New("transaction has already been finalized and cannot be canceled")
	}

	// Perform any necessary pre-cancellation validations using Synnergy Consensus
	err = synnergy.ValidateTransactionCancellation(txID)
	if err != nil {
		return errors.New("transaction cancellation validation failed: " + err.Error())
	}

	// Mark the transaction as canceled in the ledger
	tx.Canceled = true
	err = ledger.UpdateTransaction(tx)
	if err != nil {
		return errors.New("failed to cancel transaction: " + err.Error())
	}

	// Log the transaction cancellation event
	err = LogTransactionEvent(txID, "TransactionCanceled", "Transaction ID "+txID+" has been canceled.")
	if err != nil {
		return err
	}

	return nil
}

// LogTransactionEvent logs the transaction-related event in the ledger.
func LogTransactionEvent(txID, eventType, eventDescription string) error {
	eventLog := common.SYN2369Event{
		TokenID:          txID,
		EventType:        eventType,
		EventDescription: eventDescription,
		EventTime:        time.Now(),
	}

	// Store the event in the ledger
	err := ledger.StoreEvent(eventLog)
	if err != nil {
		return err
	}

	return nil
}

// GetTransactionDetails retrieves the details of a specific transaction using its transaction ID.
func GetTransactionDetails(txID string) (common.SYN2369Transaction, error) {
	// Retrieve the transaction details from the ledger
	tx, err := ledger.GetTransaction(txID)
	if err != nil {
		return common.SYN2369Transaction{}, errors.New("transaction not found: " + err.Error())
	}

	// Decrypt custom attributes for the transaction
	decryptedAttributes, err := encryption.DecryptData(tx.EncryptedMetadata)
	if err != nil {
		return common.SYN2369Transaction{}, errors.New("failed to decrypt custom attributes: " + err.Error())
	}
	tx.DecryptedMetadata = decryptedAttributes

	return tx, nil
}

// GetTransactionHistory retrieves the transaction history for a specific token.
func GetTransactionHistory(tokenID string) ([]common.SYN2369Transaction, error) {
	// Retrieve the transaction history from the ledger
	transactions, err := ledger.GetTransactionHistoryByTokenID(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve transaction history: " + err.Error())
	}

	// Decrypt custom attributes for each transaction
	for i := range transactions {
		decryptedAttributes, err := encryption.DecryptData(transactions[i].EncryptedMetadata)
		if err != nil {
			return nil, errors.New("failed to decrypt custom attributes for transaction ID " + transactions[i].TransactionID)
		}
		transactions[i].DecryptedMetadata = decryptedAttributes
	}

	return transactions, nil
}
