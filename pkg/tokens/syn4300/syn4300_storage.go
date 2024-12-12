package syn4300

import (
	"errors"
	"sync"
	"time"
)

// StorageManager handles storing, retrieving, and managing SYN4300 tokens.
type StorageManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// NewStorageManager creates a new instance of StorageManager.
func NewStorageManager(ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *StorageManager {
	return &StorageManager{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// StoreToken stores a SYN4300 token securely into the ledger after encrypting the data.
func (sm *StorageManager) StoreToken(token *Syn4300Token) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Encrypt the token data before storage
	encryptedData, err := sm.encryptionService.EncryptData(token)
	if err != nil {
		return err
	}

	// Store the encrypted token data in the ledger
	if err := sm.ledgerService.StoreToken(token.TokenID, encryptedData); err != nil {
		return err
	}

	// Log the storage event in the ledger
	if err := sm.ledgerService.LogEvent("TokenStored", time.Now(), token.TokenID); err != nil {
		return err
	}

	// Validate the storage operation with Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(token.TokenID); err != nil {
		return err
	}

	return nil
}

// RetrieveToken retrieves a SYN4300 token from the ledger and decrypts it.
func (sm *StorageManager) RetrieveToken(tokenID string) (*Syn4300Token, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the encrypted token data from the ledger
	encryptedData, err := sm.ledgerService.RetrieveToken(tokenID)
	if err != nil {
		return nil, err
	}

	// Decrypt the token data
	decryptedData, err := sm.encryptionService.DecryptData(encryptedData)
	if err != nil {
		return nil, err
	}

	// Cast the decrypted data back to SYN4300Token struct
	token := decryptedData.(*Syn4300Token)

	// Validate the retrieval with Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return nil, err
	}

	return token, nil
}

// UpdateToken updates the details of a SYN4300 token, re-encrypts it, and stores it in the ledger.
func (sm *StorageManager) UpdateToken(token *Syn4300Token) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Encrypt the updated token data
	encryptedData, err := sm.encryptionService.EncryptData(token)
	if err != nil {
		return err
	}

	// Update the encrypted token data in the ledger
	if err := sm.ledgerService.UpdateToken(token.TokenID, encryptedData); err != nil {
		return err
	}

	// Log the update event in the ledger
	if err := sm.ledgerService.LogEvent("TokenUpdated", time.Now(), token.TokenID); err != nil {
		return err
	}

	// Validate the update with Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(token.TokenID); err != nil {
		return err
	}

	return nil
}

// DeleteToken deletes a SYN4300 token from the ledger and invalidates it in the consensus.
func (sm *StorageManager) DeleteToken(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Delete the token from the ledger
	if err := sm.ledgerService.DeleteToken(tokenID); err != nil {
		return err
	}

	// Log the deletion event in the ledger
	if err := sm.ledgerService.LogEvent("TokenDeleted", time.Now(), tokenID); err != nil {
		return err
	}

	// Invalidate the token in Synnergy Consensus
	if err := sm.consensusService.InvalidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// ListTokens retrieves a list of all token IDs from the ledger.
func (sm *StorageManager) ListTokens() ([]string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve all token IDs from the ledger
	tokenIDs, err := sm.ledgerService.ListAllTokens()
	if err != nil {
		return nil, err
	}

	return tokenIDs, nil
}

// EncryptToken encrypts the token data before storage or transmission.
func (sm *StorageManager) EncryptToken(token *Syn4300Token) ([]byte, error) {
	// Encrypt the token using the encryption service
	return sm.encryptionService.EncryptData(token)
}

// DecryptToken decrypts the token data after retrieval from storage.
func (sm *StorageManager) DecryptToken(encryptedData []byte) (*Syn4300Token, error) {
	// Decrypt the token using the encryption service
	decryptedData, err := sm.encryptionService.DecryptData(encryptedData)
	if err != nil {
		return nil, err
	}

	return decryptedData.(*Syn4300Token), nil
}
