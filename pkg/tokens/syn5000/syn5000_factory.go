package syn4900

import (
	"errors"
	"sync"
	"time"
)

// Metadata defines the basic information associated with the SYN5000 token.
type Metadata struct {
	Name      string
	Symbol    string
	Supply    uint64
	GameType  string
}

// SYN5000Token represents a gambling token in the SYN5000 standard.
type SYN5000Token struct {
	TokenID      string
	Owner        string
	Metadata     Metadata
	IssuedDate   time.Time
	ExpiryDate   time.Time
	ActiveStatus bool
	SecureHash   string
	mutex        sync.RWMutex
}

// TokenFactory manages the creation, minting, and burning of SYN5000 tokens.
type TokenFactory struct {
	Ledger   *ledger.Ledger
	Storage  map[string]*SYN5000Token
	mutex    sync.RWMutex
}

// NewTokenFactory creates a new instance of TokenFactory.
func NewTokenFactory(ledger *ledger.Ledger) *TokenFactory {
	return &TokenFactory{
		Ledger:  ledger,
		Storage: make(map[string]*SYN5000Token),
	}
}

// CreateToken initializes a new SYN5000 gambling token.
func (tf *TokenFactory) CreateToken(owner, name, symbol, gameType string, supply uint64, issuedDate, expiryDate time.Time, encryptionKey string) (*SYN5000Token, error) {
	tf.mutex.Lock()
	defer tf.mutex.Unlock()

	// Generate unique Token ID using hash or any secure method.
	tokenID := generateTokenID(owner, name, issuedDate)

	// Create token metadata
	metadata := Metadata{
		Name:     name,
		Symbol:   symbol,
		Supply:   supply,
		GameType: gameType,
	}

	// Create the new token
	token := &SYN5000Token{
		TokenID:      tokenID,
		Owner:        owner,
		Metadata:     metadata,
		IssuedDate:   issuedDate,
		ExpiryDate:   expiryDate,
		ActiveStatus: true,
	}

	// Encrypt the token data
	encryptedData, err := encryption.EncryptStruct(token, encryptionKey)
	if err != nil {
		return nil, errors.New("failed to encrypt token data")
	}
	token.SecureHash = encryptedData

	// Store the token in the storage and ledger
	tf.Storage[tokenID] = token
	err = tf.Ledger.StoreData(tokenID, encryptedData)
	if err != nil {
		return nil, errors.New("failed to store token in the ledger")
	}

	// Return the newly created token
	return token, nil
}

// MintToken increases the supply of a specific SYN5000 token.
func (tf *TokenFactory) MintToken(tokenID string, amount uint64, encryptionKey string) error {
	tf.mutex.Lock()
	defer tf.mutex.Unlock()

	token, exists := tf.Storage[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	// Decrypt the token data before updating it
	err := tf.decryptTokenData(tokenID, encryptionKey)
	if err != nil {
		return err
	}

	// Update the token supply
	token.Metadata.Supply += amount

	// Re-encrypt the updated token data
	return tf.encryptAndUpdateToken(tokenID, token, encryptionKey)
}

// BurnToken decreases the supply of a specific SYN5000 token.
func (tf *TokenFactory) BurnToken(tokenID string, amount uint64, encryptionKey string) error {
	tf.mutex.Lock()
	defer tf.mutex.Unlock()

	token, exists := tf.Storage[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	// Decrypt the token data before updating it
	err := tf.decryptTokenData(tokenID, encryptionKey)
	if err != nil {
		return err
	}

	// Ensure the supply doesn't go below zero
	if token.Metadata.Supply < amount {
		return errors.New("insufficient supply to burn")
	}

	// Update the token supply
	token.Metadata.Supply -= amount

	// Re-encrypt the updated token data
	return tf.encryptAndUpdateToken(tokenID, token, encryptionKey)
}

// DeactivateToken deactivates a token, making it unusable.
func (tf *TokenFactory) DeactivateToken(tokenID string, encryptionKey string) error {
	tf.mutex.Lock()
	defer tf.mutex.Unlock()

	token, exists := tf.Storage[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	// Decrypt the token data before updating it
	err := tf.decryptTokenData(tokenID, encryptionKey)
	if err != nil {
		return err
	}

	// Deactivate the token
	token.ActiveStatus = false

	// Re-encrypt the updated token data
	return tf.encryptAndUpdateToken(tokenID, token, encryptionKey)
}

// ValidateToken ensures the token's data is valid using Synnergy Consensus.
func (tf *TokenFactory) ValidateToken(tokenID string) (bool, error) {
	tf.mutex.RLock()
	defer tf.mutex.RUnlock()

	// Check the token data's validity via Synnergy Consensus
	isValid, err := consensus.ValidateStorage(tokenID, "")
	if err != nil || !isValid {
		return false, errors.New("token validation failed via Synnergy Consensus")
	}

	return true, nil
}

// DeleteToken removes a token from storage and ledger.
func (tf *TokenFactory) DeleteToken(tokenID string) error {
	tf.mutex.Lock()
	defer tf.mutex.Unlock()

	// Delete from storage
	delete(tf.Storage, tokenID)

	// Remove from ledger
	if err := tf.Ledger.DeleteData(tokenID); err != nil {
		return errors.New("failed to delete token from ledger")
	}

	return nil
}

// encryptAndUpdateToken encrypts the updated token data and stores it in the ledger and local storage.
func (tf *TokenFactory) encryptAndUpdateToken(tokenID string, token *SYN5000Token, encryptionKey string) error {
	encryptedData, err := encryption.EncryptStruct(token, encryptionKey)
	if err != nil {
		return errors.New("failed to re-encrypt token data")
	}
	token.SecureHash = encryptedData

	// Store the updated token in the storage and ledger
	tf.Storage[tokenID] = token
	err = tf.Ledger.StoreData(tokenID, encryptedData)
	if err != nil {
		return errors.New("failed to update token in ledger")
	}

	return nil
}

// decryptTokenData decrypts the token data for modifications.
func (tf *TokenFactory) decryptTokenData(tokenID string, encryptionKey string) error {
	token, exists := tf.Storage[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	// Decrypt the token data
	err := encryption.DecryptStruct(token.SecureHash, encryptionKey, &token)
	if err != nil {
		return errors.New("failed to decrypt token data")
	}

	return nil
}

// generateTokenID generates a unique token ID based on the owner's address, token name, and issuance date.
func generateTokenID(owner, name string, issuedDate time.Time) string {
	// Example: Generate a unique hash or identifier based on inputs
	return encryption.HashString(owner + name + issuedDate.String())
}
