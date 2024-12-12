package syn3300

import (
	"sync"
	"time"

)

// Syn3300StorageManager handles storage operations for SYN3300 tokens.
type Syn3300StorageManager struct {
	ledgerService    *ledger.Ledger           // Ledger for event logging and validation
	encryptionService *encryption.Encryptor    // Encryption service to secure stored data
	consensusService *consensus.SynnergyConsensus // Consensus service to validate sub-blocks
	mutex            sync.Mutex               // Mutex to ensure thread-safe operations
	storage          map[string]*Syn3300Token // In-memory storage for SYN3300 tokens
}

// NewSyn3300StorageManager creates a new instance of Syn3300StorageManager.
func NewSyn3300StorageManager(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *Syn3300StorageManager {
	return &Syn3300StorageManager{
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService: consensus,
		storage:          make(map[string]*Syn3300Token),
	}
}

// StoreToken stores a new SYN3300 token securely and logs the event.
func (sm *Syn3300StorageManager) StoreToken(token *Syn3300Token) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Encrypt the token data before storing.
	encryptedToken, err := sm.encryptionService.EncryptData(token)
	if err != nil {
		return err
	}

	// Log the storage event in the ledger.
	if err := sm.ledgerService.LogEvent("TokenStored", time.Now(), token.ID); err != nil {
		return err
	}

	// Validate the token creation using Synnergy Consensus.
	if err := sm.consensusService.ValidateSubBlock(token.ID); err != nil {
		return err
	}

	// Store the encrypted token.
	sm.storage[token.ID] = encryptedToken.(*Syn3300Token)

	return nil
}

// RetrieveToken retrieves a SYN3300 token by its ID.
func (sm *Syn3300StorageManager) RetrieveToken(tokenID string) (*Syn3300Token, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the encrypted token.
	encryptedToken, exists := sm.storage[tokenID]
	if !exists {
		return nil, errors.New("token not found")
	}

	// Decrypt the token data.
	decryptedToken, err := sm.encryptionService.DecryptData(encryptedToken)
	if err != nil {
		return nil, err
	}

	return decryptedToken.(*Syn3300Token), nil
}

// UpdateToken updates the details of an existing SYN3300 token and logs the event.
func (sm *Syn3300StorageManager) UpdateToken(updatedToken *Syn3300Token) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if the token exists.
	if _, exists := sm.storage[updatedToken.ID]; !exists {
		return errors.New("token not found")
	}

	// Encrypt the updated token data.
	encryptedToken, err := sm.encryptionService.EncryptData(updatedToken)
	if err != nil {
		return err
	}

	// Log the update event in the ledger.
	if err := sm.ledgerService.LogEvent("TokenUpdated", time.Now(), updatedToken.ID); err != nil {
		return err
	}

	// Validate the update using Synnergy Consensus.
	if err := sm.consensusService.ValidateSubBlock(updatedToken.ID); err != nil {
		return err
	}

	// Update the token in storage.
	sm.storage[updatedToken.ID] = encryptedToken.(*Syn3300Token)

	return nil
}

// DeleteToken deletes a SYN3300 token by its ID and logs the event.
func (sm *Syn3300StorageManager) DeleteToken(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if the token exists.
	if _, exists := sm.storage[tokenID]; !exists {
		return errors.New("token not found")
	}

	// Log the deletion event in the ledger.
	if err := sm.ledgerService.LogEvent("TokenDeleted", time.Now(), tokenID); err != nil {
		return err
	}

	// Validate the deletion using Synnergy Consensus.
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	// Delete the token from storage.
	delete(sm.storage, tokenID)

	return nil
}

// ListAllTokens lists all SYN3300 tokens currently stored.
func (sm *Syn3300StorageManager) ListAllTokens() ([]*Syn3300Token, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve and decrypt all stored tokens.
	var allTokens []*Syn3300Token
	for _, encryptedToken := range sm.storage {
		decryptedToken, err := sm.encryptionService.DecryptData(encryptedToken)
		if err != nil {
			return nil, err
		}
		allTokens = append(allTokens, decryptedToken.(*Syn3300Token))
	}

	return allTokens, nil
}

// ValidateTokenIntegrity verifies the integrity of a stored SYN3300 token by comparing its secure hash.
func (sm *Syn3300StorageManager) ValidateTokenIntegrity(tokenID string) (bool, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the encrypted token.
	token, exists := sm.storage[tokenID]
	if !exists {
		return false, errors.New("token not found")
	}

	// Recalculate the secure hash and compare it to the stored value.
	currentHash := sm.encryptionService.GenerateHash(token)
	originalHash := sm.encryptionService.GetStoredHash(tokenID) // Assuming there's a method for this.

	return currentHash == originalHash, nil
}
