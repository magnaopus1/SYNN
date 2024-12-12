package syn4900

import (
	"errors"
	"sync"
	"time"
)

// TokenStorage represents the storage manager for SYN4900 tokens.
type TokenStorage struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	mutex             sync.Mutex
}

// NewTokenStorage initializes a new TokenStorage instance.
func NewTokenStorage(ledgerService *ledger.LedgerService, encryptionService *encryption.Encryptor) *TokenStorage {
	return &TokenStorage{
		ledgerService:     ledgerService,
		encryptionService: encryptionService,
	}
}

// StoreToken stores a new SYN4900 token in the ledger.
func (ts *TokenStorage) StoreToken(token *Syn4900Token) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	// Encrypt the token data before storage.
	encryptedToken, err := ts.encryptionService.EncryptData(token)
	if err != nil {
		return err
	}

	// Log the event in the ledger.
	if err := ts.ledgerService.LogEvent("TokenStored", time.Now(), token.TokenID); err != nil {
		return err
	}

	// Store the token in the ledger.
	return ts.ledgerService.StoreToken(token.TokenID, encryptedToken)
}

// RetrieveToken retrieves a SYN4900 token from the ledger.
func (ts *TokenStorage) RetrieveToken(tokenID string) (*Syn4900Token, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	// Retrieve the encrypted token data from the ledger.
	encryptedToken, err := ts.ledgerService.RetrieveToken(tokenID)
	if err != nil {
		return nil, err
	}

	// Decrypt the token data.
	decryptedToken, err := ts.encryptionService.DecryptData(encryptedToken)
	if err != nil {
		return nil, err
	}

	// Return the decrypted token.
	return decryptedToken.(*Syn4900Token), nil
}

// UpdateToken updates the details of an existing SYN4900 token in the ledger.
func (ts *TokenStorage) UpdateToken(token *Syn4900Token) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	// Encrypt the updated token data.
	encryptedToken, err := ts.encryptionService.EncryptData(token)
	if err != nil {
		return err
	}

	// Log the token update event in the ledger.
	if err := ts.ledgerService.LogEvent("TokenUpdated", time.Now(), token.TokenID); err != nil {
		return err
	}

	// Update the token data in the ledger.
	return ts.ledgerService.StoreToken(token.TokenID, encryptedToken)
}

// DeleteToken removes a SYN4900 token from the ledger.
func (ts *TokenStorage) DeleteToken(tokenID string) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	// Log the deletion event in the ledger.
	if err := ts.ledgerService.LogEvent("TokenDeleted", time.Now(), tokenID); err != nil {
		return err
	}

	// Remove the token from the ledger.
	return ts.ledgerService.DeleteToken(tokenID)
}

// VerifyToken ensures the integrity of a stored token by comparing it with the ledger.
func (ts *TokenStorage) VerifyToken(tokenID string) (bool, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	// Retrieve the token from the ledger.
	token, err := ts.RetrieveToken(tokenID)
	if err != nil {
		return false, err
	}

	// Verify the token's integrity.
	if token == nil {
		return false, errors.New("token not found")
	}

	// Additional verification logic can be implemented here if needed.
	return true, nil
}

// ListAllTokens retrieves a list of all tokens stored in the ledger.
func (ts *TokenStorage) ListAllTokens() ([]*Syn4900Token, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	// Retrieve all token IDs from the ledger.
	tokenIDs, err := ts.ledgerService.ListAllTokenIDs()
	if err != nil {
		return nil, err
	}

	// Retrieve each token and append to the list.
	var tokens []*Syn4900Token
	for _, tokenID := range tokenIDs {
		token, err := ts.RetrieveToken(tokenID)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}
