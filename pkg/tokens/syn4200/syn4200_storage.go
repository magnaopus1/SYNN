package syn4200

import (
	"errors"
	"time"
	"sync"
)

// Syn4200StorageManager handles secure storage and retrieval of SYN4200 tokens.
type Syn4200StorageManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex // Ensures thread-safe operations
}

// NewSyn4200StorageManager initializes a new storage manager.
func NewSyn4200StorageManager(ledgerService *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *Syn4200StorageManager {
	return &Syn4200StorageManager{
		ledgerService:     ledgerService,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// StoreToken securely stores a Syn4200Token to the ledger.
func (sm *Syn4200StorageManager) StoreToken(token *Syn4200Token) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Serialize the token data for storage
	serializedToken, err := sm.serializeToken(token)
	if err != nil {
		return fmt.Errorf("failed to serialize token: %w", err)
	}

	// Encrypt the token data
	encryptedData, err := sm.encryptionService.EncryptData(serializedToken)
	if err != nil {
		return fmt.Errorf("failed to encrypt token data: %w", err)
	}

	// Store the encrypted token in the ledger
	if err := sm.ledgerService.StoreData(token.TokenID, encryptedData); err != nil {
		return fmt.Errorf("failed to store token in ledger: %w", err)
	}

	// Log the event and validate the storage using Synnergy Consensus
	sm.logStorageEvent(token.TokenID, "TokenStored")
	if err := sm.consensusService.ValidateSubBlock(token.TokenID); err != nil {
		return fmt.Errorf("failed to validate token storage: %w", err)
	}

	return nil
}

// RetrieveToken retrieves a stored Syn4200Token from the ledger.
func (sm *Syn4200StorageManager) RetrieveToken(tokenID string) (*Syn4200Token, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve encrypted token data from the ledger
	encryptedData, err := sm.ledgerService.RetrieveData(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token from ledger: %w", err)
	}

	// Decrypt the token data
	decryptedData, err := sm.encryptionService.DecryptData(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token data: %w", err)
	}

	// Deserialize the token data into a Syn4200Token
	token, err := sm.deserializeToken(decryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize token: %w", err)
	}

	return token, nil
}

// DeleteToken securely deletes a Syn4200Token from the ledger.
func (sm *Syn4200StorageManager) DeleteToken(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Remove token data from the ledger
	if err := sm.ledgerService.DeleteData(tokenID); err != nil {
		return fmt.Errorf("failed to delete token from ledger: %w", err)
	}

	// Log the event and validate the deletion using Synnergy Consensus
	sm.logStorageEvent(tokenID, "TokenDeleted")
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return fmt.Errorf("failed to validate token deletion: %w", err)
	}

	return nil
}

// serializeToken serializes a Syn4200Token into a storable format.
func (sm *Syn4200StorageManager) serializeToken(token *Syn4200Token) ([]byte, error) {
	// Convert the Syn4200Token to JSON or another serializable format
	// In production, you'd use more advanced serialization
	data, err := json.Marshal(token)
	if err != nil {
		return nil, fmt.Errorf("serialization error: %w", err)
	}
	return data, nil
}

// deserializeToken deserializes the token data back into a Syn4200Token.
func (sm *Syn4200StorageManager) deserializeToken(data []byte) (*Syn4200Token, error) {
	var token Syn4200Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("deserialization error: %w", err)
	}
	return &token, nil
}

// logStorageEvent logs the event of storing, deleting, or modifying tokens.
func (sm *Syn4200StorageManager) logStorageEvent(tokenID, eventType string) {
	// Log the event in the ledger
	_ = sm.ledgerService.LogEvent(eventType, time.Now(), tokenID)
}
