package syn3900

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
	"sync"

)

// StorageManager handles secure storage and retrieval of SYN3900 tokens.
type StorageManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// NewStorageManager creates a new StorageManager instance.
func NewStorageManager(ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *StorageManager {
	return &StorageManager{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// StoreToken securely stores the token in the ledger after encrypting its data.
func (sm *StorageManager) StoreToken(token *Syn3900Token) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Encrypt token data
	encryptedData, err := sm.encryptionService.EncryptData(token)
	if err != nil {
		return fmt.Errorf("failed to encrypt token data: %w", err)
	}

	// Store encrypted token in the ledger
	if err := sm.ledgerService.StoreToken(token.TokenID, encryptedData); err != nil {
		return fmt.Errorf("failed to store token in ledger: %w", err)
	}

	// Log storage event
	if err := sm.ledgerService.LogEvent("TokenStored", time.Now(), token.TokenID); err != nil {
		return fmt.Errorf("failed to log storage event: %w", err)
	}

	// Validate token storage using Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(token.TokenID); err != nil {
		return fmt.Errorf("failed to validate token storage with consensus: %w", err)
	}

	return nil
}

// RetrieveToken securely retrieves and decrypts a token from the ledger.
func (sm *StorageManager) RetrieveToken(tokenID string) (*Syn3900Token, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve encrypted token data from ledger
	encryptedData, err := sm.ledgerService.RetrieveToken(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token from ledger: %w", err)
	}

	// Decrypt token data
	decryptedToken, err := sm.encryptionService.DecryptData(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token data: %w", err)
	}

	token, ok := decryptedToken.(*Syn3900Token)
	if !ok {
		return nil, errors.New("retrieved data does not match expected token format")
	}

	return token, nil
}

// UpdateToken securely updates an existing token in the ledger.
func (sm *StorageManager) UpdateToken(token *Syn3900Token) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Encrypt updated token data
	encryptedData, err := sm.encryptionService.EncryptData(token)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated token data: %w", err)
	}

	// Update encrypted token in the ledger
	if err := sm.ledgerService.UpdateToken(token.TokenID, encryptedData); err != nil {
		return fmt.Errorf("failed to update token in ledger: %w", err)
	}

	// Log update event
	if err := sm.ledgerService.LogEvent("TokenUpdated", time.Now(), token.TokenID); err != nil {
		return fmt.Errorf("failed to log update event: %w", err)
	}

	// Validate token update using Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(token.TokenID); err != nil {
		return fmt.Errorf("failed to validate token update with consensus: %w", err)
	}

	return nil
}

// DeleteToken securely removes a token from the ledger and logs the action.
func (sm *StorageManager) DeleteToken(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Delete token from ledger
	if err := sm.ledgerService.DeleteToken(tokenID); err != nil {
		return fmt.Errorf("failed to delete token from ledger: %w", err)
	}

	// Log deletion event
	if err := sm.ledgerService.LogEvent("TokenDeleted", time.Now(), tokenID); err != nil {
		return fmt.Errorf("failed to log token deletion: %w", err)
	}

	// Validate token deletion using Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return fmt.Errorf("failed to validate token deletion with consensus: %w", err)
	}

	return nil
}

// RetrieveAllTokens retrieves all tokens from the ledger, decrypting each one securely.
func (sm *StorageManager) RetrieveAllTokens() ([]*Syn3900Token, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve all encrypted tokens from the ledger
	encryptedTokens, err := sm.ledgerService.RetrieveAllTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all tokens from ledger: %w", err)
	}

	// Decrypt each token
	var tokens []*Syn3900Token
	for _, encryptedToken := range encryptedTokens {
		decryptedToken, err := sm.encryptionService.DecryptData(encryptedToken)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt token data: %w", err)
		}

		token, ok := decryptedToken.(*Syn3900Token)
		if !ok {
			return nil, errors.New("retrieved data does not match expected token format")
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

// VerifyTokenIntegrity verifies the integrity of a token by rehashing and checking against the stored value.
func (sm *StorageManager) VerifyTokenIntegrity(tokenID string) (bool, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the token
	token, err := sm.RetrieveToken(tokenID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve token for integrity check: %w", err)
	}

	// Rehash token metadata and compare with stored hash
	metadataHash := sm.encryptionService.HashData(token.Metadata)
	storedHash, err := sm.ledgerService.GetTokenMetadataHash(tokenID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve stored hash from ledger: %w", err)
	}

	return metadataHash == storedHash, nil
}
