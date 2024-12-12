package syn2600

import (

	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

)

// SYN2600TokenStorage represents storage and data retrieval functions for SYN2600 tokens.
type SYN2600TokenStorage struct {
	TokenID        string
	AssetDetails   string
	Owner          string
	EncryptedToken string
	IssuedDate     time.Time
	ExpiryDate     time.Time
	ActiveStatus   bool
}

// StoreToken stores a new SYN2600 token in the ledger after encrypting it.
func StoreToken(token *SYN2600TokenStorage) (string, error) {
	// Encrypt token data before storage
	encryptedToken, err := encryption.EncryptTokenData(token)
	if err != nil {
		return "", errors.New("failed to encrypt token data before storage")
	}

	// Store encrypted token in the ledger
	err = ledger.StoreInvestorToken(encryptedToken)
	if err != nil {
		return "", errors.New("failed to store encrypted token in the ledger")
	}

	// Validate the transaction using Synnergy Consensus (sub-block validation)
	err = synconsensus.ValidateSubBlockTransaction(token.TokenID, encryptedToken)
	if err != nil {
		return "", errors.New("failed to validate token storage transaction")
	}

	return token.TokenID, nil
}

// FetchToken retrieves an encrypted SYN2600 token from the ledger and decrypts it.
func FetchToken(tokenID string) (*SYN2600TokenStorage, error) {
	// Fetch the encrypted token from the ledger
	encryptedToken, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to fetch token from ledger")
	}

	// Decrypt the token data for usage
	decryptedToken, err := encryption.DecryptTokenData(encryptedToken)
	if err != nil {
		return nil, errors.New("failed to decrypt token data")
	}

	return decryptedToken, nil
}

// UpdateToken updates an existing SYN2600 token in the ledger after encrypting new data.
func UpdateToken(token *SYN2600TokenStorage) (string, error) {
	// Encrypt the updated token data
	encryptedToken, err := encryption.EncryptTokenData(token)
	if err != nil {
		return "", errors.New("failed to encrypt token data before update")
	}

	// Update the encrypted token in the ledger
	err = ledger.UpdateInvestorToken(encryptedToken)
	if err != nil {
		return "", errors.New("failed to update token in the ledger")
	}

	// Validate the transaction using Synnergy Consensus (sub-block validation)
	err = synconsensus.ValidateSubBlockTransaction(token.TokenID, encryptedToken)
	if err != nil {
		return "", errors.New("failed to validate token update transaction")
	}

	return token.TokenID, nil
}

// DeleteToken removes an existing SYN2600 token from the ledger and records the action.
func DeleteToken(tokenID string) (string, error) {
	// Delete the token from the ledger
	err := ledger.DeleteInvestorToken(tokenID)
	if err != nil {
		return "", errors.New("failed to delete token from ledger")
	}

	// Record the token deletion event
	eventID, err := RecordSecurityEvent(tokenID, "DELETE", "Token deleted from ledger", "")
	if err != nil {
		return "", errors.New("failed to record token deletion event")
	}

	return eventID, nil
}

// RecordTransactionLog records the transaction log of the token in the ledger for audit purposes.
func RecordTransactionLog(tokenID string, transactionType string, details string) error {
	// Create a transaction log with a timestamp
	log := common.TransactionLog{
		TokenID:        tokenID,
		TransactionType: transactionType,
		Details:        details,
		Timestamp:      time.Now(),
	}

	// Store the transaction log in the ledger
	err := ledger.RecordTransactionLog(log)
	if err != nil {
		return errors.New("failed to record transaction log in the ledger")
	}

	return nil
}

// TokenExists checks if a given token exists in the ledger.
func TokenExists(tokenID string) (bool, error) {
	// Check existence in the ledger
	exists, err := ledger.TokenExists(tokenID)
	if err != nil {
		return false, errors.New("failed to check if token exists")
	}

	return exists, nil
}

// ListAllTokens returns a list of all stored SYN2600 tokens from the ledger.
func ListAllTokens() ([]*SYN2600TokenStorage, error) {
	// Fetch all encrypted tokens from the ledger
	encryptedTokens, err := ledger.FetchAllInvestorTokens()
	if err != nil {
		return nil, errors.New("failed to fetch all tokens from ledger")
	}

	// Decrypt each token and append to the result list
	var tokens []*SYN2600TokenStorage
	for _, encryptedToken := range encryptedTokens {
		decryptedToken, err := encryption.DecryptTokenData(encryptedToken)
		if err != nil {
			return nil, errors.New("failed to decrypt token data for list")
		}
		tokens = append(tokens, decryptedToken)
	}

	return tokens, nil
}
