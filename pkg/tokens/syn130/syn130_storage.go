package syn130

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// StorageManager handles the storage, encryption, and retrieval of SYN130 tokens and related assets.
type StorageManager struct {
	StorageLedger     *ledger.StorageLedger         // Ledger to track all storage-related events
	EncryptionService *encryption.EncryptionService // Encryption service to secure storage
	Consensus         *consensus.SynnergyConsensus  // Synnergy Consensus for storage validation
	mutex             sync.Mutex                    // Mutex for thread-safe storage operations
}

// NewStorageManager initializes a new StorageManager.
func NewStorageManager(storageLedger *ledger.StorageLedger, encryptionService *encryption.EncryptionService, consensusEngine *consensus.SynnergyConsensus) *StorageManager {
	return &StorageManager{
		StorageLedger:     storageLedger,
		EncryptionService: encryptionService,
		Consensus:         consensusEngine,
	}
}

// StoreTokenData securely stores the token data and metadata in the ledger, ensuring encryption and validation.
func (sm *StorageManager) StoreTokenData(tokenID string, tokenData map[string]interface{}) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Convert token data to JSON for storage
	tokenDataJSON, err := json.Marshal(tokenData)
	if err != nil {
		return errors.New("failed to serialize token data")
	}

	// Encrypt the token data before storage
	encryptedData, err := sm.EncryptionService.EncryptData(tokenDataJSON)
	if err != nil {
		return errors.New("failed to encrypt token data")
	}

	// Validate the storage event with Synnergy Consensus
	if err := sm.Consensus.ValidateStorage(tokenID, encryptedData); err != nil {
		return errors.New("storage validation failed through Synnergy Consensus")
	}

	// Store the encrypted data in the ledger
	event := ledger.StorageEvent{
		EventID:    tokenID,
		Data:       encryptedData,
		Timestamp:  time.Now(),
		EventType:  "Token Data Storage",
	}
	if err := sm.StorageLedger.RecordStorageEvent(&event); err != nil {
		return errors.New("failed to record storage event in ledger")
	}

	return nil
}

// RetrieveTokenData fetches and decrypts token data based on the token ID.
func (sm *StorageManager) RetrieveTokenData(tokenID string) (map[string]interface{}, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the encrypted data from the ledger
	event, err := sm.StorageLedger.GetStorageEvent(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve storage event from ledger")
	}

	// Decrypt the data
	decryptedData, err := sm.EncryptionService.DecryptData(event.Data)
	if err != nil {
		return nil, errors.New("failed to decrypt token data")
	}

	// Unmarshal the decrypted data into a map
	var tokenData map[string]interface{}
	if err := json.Unmarshal(decryptedData, &tokenData); err != nil {
		return nil, errors.New("failed to unmarshal token data")
	}

	return tokenData, nil
}

// StoreAssetMetadata securely stores asset metadata associated with a token.
func (sm *StorageManager) StoreAssetMetadata(tokenID string, metadata map[string]interface{}) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Serialize metadata to JSON
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return errors.New("failed to serialize asset metadata")
	}

	// Encrypt the metadata
	encryptedMetadata, err := sm.EncryptionService.EncryptData(metadataJSON)
	if err != nil {
		return errors.New("failed to encrypt asset metadata")
	}

	// Validate the metadata storage with Synnergy Consensus
	if err := sm.Consensus.ValidateStorage(tokenID, encryptedMetadata); err != nil {
		return errors.New("metadata storage validation failed through Synnergy Consensus")
	}

	// Store the encrypted metadata in the ledger
	event := ledger.StorageEvent{
		EventID:    tokenID,
		Data:       encryptedMetadata,
		Timestamp:  time.Now(),
		EventType:  "Asset Metadata Storage",
	}
	if err := sm.StorageLedger.RecordStorageEvent(&event); err != nil {
		return errors.New("failed to record asset metadata storage event in ledger")
	}

	return nil
}

// RetrieveAssetMetadata fetches and decrypts the asset metadata associated with a token.
func (sm *StorageManager) RetrieveAssetMetadata(tokenID string) (map[string]interface{}, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the encrypted metadata from the ledger
	event, err := sm.StorageLedger.GetStorageEvent(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve asset metadata from ledger")
	}

	// Decrypt the metadata
	decryptedMetadata, err := sm.EncryptionService.DecryptData(event.Data)
	if err != nil {
		return nil, errors.New("failed to decrypt asset metadata")
	}

	// Unmarshal the decrypted metadata into a map
	var metadata map[string]interface{}
	if err := json.Unmarshal(decryptedMetadata, &metadata); err != nil {
		return nil, errors.New("failed to unmarshal asset metadata")
	}

	return metadata, nil
}

// StoreTransactionHistory securely stores a token's transaction history.
func (sm *StorageManager) StoreTransactionHistory(tokenID string, transactions []map[string]interface{}) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Serialize transaction history
	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		return errors.New("failed to serialize transaction history")
	}

	// Encrypt the transaction history
	encryptedData, err := sm.EncryptionService.EncryptData(transactionsJSON)
	if err != nil {
		return errors.New("failed to encrypt transaction history")
	}

	// Validate the transaction history storage through Synnergy Consensus
	if err := sm.Consensus.ValidateStorage(tokenID, encryptedData); err != nil {
		return errors.New("transaction history validation failed through Synnergy Consensus")
	}

	// Store the encrypted transaction history in the ledger
	event := ledger.StorageEvent{
		EventID:    tokenID,
		Data:       encryptedData,
		Timestamp:  time.Now(),
		EventType:  "Transaction History Storage",
	}
	if err := sm.StorageLedger.RecordStorageEvent(&event); err != nil {
		return errors.New("failed to record transaction history storage event in ledger")
	}

	return nil
}

// RetrieveTransactionHistory fetches and decrypts a token's transaction history.
func (sm *StorageManager) RetrieveTransactionHistory(tokenID string) ([]map[string]interface{}, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the encrypted transaction history from the ledger
	event, err := sm.StorageLedger.GetStorageEvent(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve transaction history from ledger")
	}

	// Decrypt the transaction history
	decryptedData, err := sm.EncryptionService.DecryptData(event.Data)
	if err != nil {
		return nil, errors.New("failed to decrypt transaction history")
	}

	// Unmarshal the decrypted data into a list of transactions
	var transactions []map[string]interface{}
	if err := json.Unmarshal(decryptedData, &transactions); err != nil {
		return nil, errors.New("failed to unmarshal transaction history")
	}

	return transactions, nil
}

// ValidateStorageIntegrity audits and verifies the integrity of stored data using Synnergy Consensus.
func (sm *StorageManager) ValidateStorageIntegrity(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the encrypted data from the ledger
	event, err := sm.StorageLedger.GetStorageEvent(tokenID)
	if err != nil {
		return errors.New("failed to retrieve storage event for validation")
	}

	// Validate the integrity of the stored data using Synnergy Consensus
	if err := sm.Consensus.ValidateStorage(tokenID, event.Data); err != nil {
		return errors.New("storage integrity validation failed through Synnergy Consensus")
	}

	return nil
}
