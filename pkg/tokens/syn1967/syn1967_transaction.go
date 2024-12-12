package syn1967

import (
	"fmt"
	"time"
	"sync"
)

// TransactionManager handles all SYN1967 token transactions
type TransactionManager struct {
	mu sync.Mutex
}

// TransferToken handles the secure transfer of ownership for a SYN1967 token
func (tm *TransactionManager) TransferToken(tokenID, from, to string, transferAmount float64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Retrieve the token from storage
	token, err := storage.RetrieveToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token for transfer: %v", err)
	}

	// Validate that the sender is the current owner
	if token.Owner != from {
		return fmt.Errorf("sender does not own the token: %s", from)
	}

	// Fractional ownership checks, if applicable
	if token.Fractionalized {
		if transferAmount > token.Amount {
			return fmt.Errorf("transfer amount exceeds the available token balance")
		}
		// Adjust the token balance
		token.Amount -= transferAmount
		// Create a new fractional token for the recipient
		newToken := *token // Copy token
		newToken.TokenID = generateUniqueID()
		newToken.Owner = to
		newToken.Amount = transferAmount

		// Store the new fractional token for the recipient
		if err := storage.StoreToken(&newToken); err != nil {
			return fmt.Errorf("failed to store fractionalized token for recipient: %v", err)
		}
	} else {
		// Full token transfer
		token.Owner = to
	}

	// Update the token in storage
	if err := storage.UpdateToken(token); err != nil {
		return fmt.Errorf("failed to update token ownership: %v", err)
	}

	// Log the transfer in the ledger
	transactionLog := common.TransactionLog{
		TransactionID:    generateUniqueID(),
		TokenID:          token.TokenID,
		Sender:           from,
		Recipient:        to,
		Amount:           transferAmount,
		TransactionType:  "Transfer",
		TransactionDate:  time.Now(),
	}
	err = ledger.RecordTransaction(transactionLog.TransactionID, transactionLog, "Token Transfer")
	if err != nil {
		return fmt.Errorf("failed to record token transfer in the ledger: %v", err)
	}

	return nil
}

// ValidateTransaction validates a SYN1967 token transaction within the context of the Synnergy Consensus
func (tm *TransactionManager) ValidateTransaction(transactionID, subBlockID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Validate the transaction under the Synnergy Consensus
	valid, err := subblock.ValidateSubBlockTransaction(subBlockID, transactionID)
	if err != nil {
		return fmt.Errorf("failed to validate transaction: %v", err)
	}

	if !valid {
		return fmt.Errorf("transaction %s is not part of the validated sub-block", transactionID)
	}

	return nil
}

// RecordTransaction handles general SYN1967 transactions, such as purchase or trade
func (tm *TransactionManager) RecordTransaction(tokenID, from, to string, transactionType string, amount float64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Retrieve the token
	token, err := storage.RetrieveToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token for transaction: %v", err)
	}

	// Encrypt transaction data
	transactionData := common.TransactionLog{
		TransactionID:    generateUniqueID(),
		TokenID:          token.TokenID,
		Sender:           from,
		Recipient:        to,
		Amount:           amount,
		TransactionType:  transactionType,
		TransactionDate:  time.Now(),
	}
	encryptedTransactionData, err := encryption.Encrypt(transactionData)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction data: %v", err)
	}

	// Store the encrypted transaction in the ledger
	err = ledger.RecordTransaction(transactionData.TransactionID, encryptedTransactionData, transactionType)
	if err != nil {
		return fmt.Errorf("failed to record transaction in the ledger: %v", err)
	}

	return nil
}

// GenerateTransactionHistory retrieves the complete transaction history for a SYN1967 token
func (tm *TransactionManager) GenerateTransactionHistory(tokenID string) ([]common.TransactionLog, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Retrieve the transaction history from the ledger
	transactionLogs, err := ledger.RetrieveTransactionLogs(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction history for token: %v", err)
	}

	return transactionLogs, nil
}

// RevokeTransaction reverses or cancels a previously recorded SYN1967 token transaction
func (tm *TransactionManager) RevokeTransaction(transactionID, reason string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Retrieve the transaction from the ledger
	transaction, err := ledger.RetrieveTransaction(transactionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction for revocation: %v", err)
	}

	// Ensure that the transaction can be revoked based on business logic
	if time.Since(transaction.TransactionDate).Hours() > 24 {
		return fmt.Errorf("transactions older than 24 hours cannot be revoked")
	}

	// Log the revocation in the ledger
	err = ledger.RecordTransaction(generateUniqueID(), nil, fmt.Sprintf("Transaction %s revoked: %s", transactionID, reason))
	if err != nil {
		return fmt.Errorf("failed to record transaction revocation: %v", err)
	}

	return nil
}

// generateUniqueID generates a unique identifier for transactions and tokens
func generateUniqueID() string {
	return fmt.Sprintf("TX-%d", time.Now().UnixNano())
}
