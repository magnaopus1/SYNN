package syn4900

import (
	"errors"
	"sync"
	"time"
)

// StorageManager handles persistent storage operations for SYN5000 tokens.
type StorageManager struct {
	mu           sync.RWMutex
	storagePath  string                     // Path to the storage file for persistence
	tokens       map[string]*SYN5000Token    // In-memory storage of tokens
	ledger       *ledger.Ledger              // Ledger integration for tracking storage-related operations
	encryptor    *encryption.Encryptor       // Encryption manager for securing token data
	consensus    *consensus.SynnergyConsensus // Synnergy Consensus for validating storage blocks
}

// NewStorageManager creates a new instance of StorageManager.
func NewStorageManager(storagePath string, ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *StorageManager {
	return &StorageManager{
		storagePath: storagePath,
		tokens:      make(map[string]*SYN5000Token),
		ledger:      ledger,
		encryptor:   encryptor,
		consensus:   consensus,
	}
}

// LoadTokens loads the tokens from persistent storage into memory.
func (sm *StorageManager) LoadTokens() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	file, err := os.Open(sm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No storage file exists yet, nothing to load
		}
		return err
	}
	defer file.Close()

	var storedTokens map[string]*SYN5000Token
	err = json.NewDecoder(file).Decode(&storedTokens)
	if err != nil {
		return err
	}

	sm.tokens = storedTokens

	// Log storage load in the ledger.
	sm.ledger.LogStorageEvent("TokensLoaded", time.Now(), "Loaded tokens from storage")
	return nil
}

// SaveTokens saves the current state of tokens to persistent storage.
func (sm *StorageManager) SaveTokens() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	file, err := os.Create(sm.storagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(sm.tokens)
	if err != nil {
		return err
	}

	// Log storage save in the ledger.
	sm.ledger.LogStorageEvent("TokensSaved", time.Now(), "Saved tokens to storage")
	return nil
}

// AddToken adds a new SYN5000Token to the storage and ledger.
func (sm *StorageManager) AddToken(token *SYN5000Token) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.tokens[token.TokenID]; exists {
		return errors.New("token already exists")
	}

	// Encrypt token data before storing.
	encryptedToken, err := sm.encryptor.EncryptData(token)
	if err != nil {
		return err
	}

	// Add token to the in-memory map.
	sm.tokens[token.TokenID] = encryptedToken

	// Log the addition in the ledger.
	sm.ledger.LogStorageEvent("TokenAdded", time.Now(), token.TokenID)
	return sm.SaveTokens()
}

// GetToken retrieves a token from the storage based on its TokenID.
func (sm *StorageManager) GetToken(tokenID string) (*SYN5000Token, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	encryptedToken, exists := sm.tokens[tokenID]
	if !exists {
		return nil, errors.New("token not found")
	}

	// Decrypt token data.
	token, err := sm.encryptor.DecryptData(encryptedToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// UpdateToken updates the information of an existing token in the storage.
func (sm *StorageManager) UpdateToken(token *SYN5000Token) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.tokens[token.TokenID]; !exists {
		return errors.New("token not found")
	}

	// Encrypt token data before updating.
	encryptedToken, err := sm.encryptor.EncryptData(token)
	if err != nil {
		return err
	}

	// Update token in memory.
	sm.tokens[token.TokenID] = encryptedToken

	// Log the update in the ledger.
	sm.ledger.LogStorageEvent("TokenUpdated", time.Now(), token.TokenID)
	return sm.SaveTokens()
}

// RemoveToken deletes a token from the storage.
func (sm *StorageManager) RemoveToken(tokenID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.tokens[tokenID]; !exists {
		return errors.New("token not found")
	}

	delete(sm.tokens, tokenID)

	// Log the removal in the ledger.
	sm.ledger.LogStorageEvent("TokenRemoved", time.Now(), tokenID)
	return sm.SaveTokens()
}

// SyncWithConsensus validates storage with the consensus mechanism for integrity.
func (sm *StorageManager) SyncWithConsensus() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Generate a snapshot of the current tokens.
	tokensSnapshot := sm.generateSnapshot()

	// Validate the snapshot with the Synnergy Consensus.
	if err := sm.consensus.ValidateSubBlock(tokensSnapshot); err != nil {
		return err
	}

	// Log the consensus sync in the ledger.
	sm.ledger.LogStorageEvent("ConsensusSynced", time.Now(), "Synchronized storage with consensus")
	return nil
}

// generateSnapshot generates a snapshot of the current token state for consensus validation.
func (sm *StorageManager) generateSnapshot() string {
	snapshotData, _ := json.Marshal(sm.tokens)
	return string(snapshotData)
}
