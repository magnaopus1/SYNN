package syn131

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// NewStorageManager initializes a new StorageManager instance.
func NewStorageManager(storagePath string, ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *StorageManager {
	return &StorageManager{
		storagePath:      storagePath,
		Ledger:           ledger,
		ConsensusEngine:  consensusEngine,
		EncryptionService: encryptionService,
	}
}

// SaveToken saves a SYN131 token with all its metadata to persistent storage.
func (sm *StorageManager) SaveToken(token *Syn131Token) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Serialize the token and encrypt it before saving
	tokenData, err := json.Marshal(token)
	if err != nil {
		return errors.New("failed to serialize token data")
	}

	encryptedTokenData, encryptionKey, err := sm.EncryptionService.EncryptData(tokenData)
	if err != nil {
		return errors.New("failed to encrypt token data")
	}

	// Store the encrypted token data to a file
	filePath := sm.storagePath + "/tokens/" + token.ID + ".json"
	err = os.WriteFile(filePath, []byte(encryptedTokenData), 0644)
	if err != nil {
		return errors.New("failed to write token data to file")
	}

	// Store the encryption key securely in the ledger
	err = sm.Ledger.StoreEncryptionKey(token.ID, encryptionKey)
	if err != nil {
		return errors.New("failed to store encryption key in the ledger")
	}

	return nil
}

// LoadToken loads a SYN131 token from persistent storage.
func (sm *StorageManager) LoadToken(tokenID string) (*Syn131Token, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the encrypted token data from the file
	filePath := sm.storagePath + "/tokens/" + tokenID + ".json"
	encryptedTokenData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("failed to read token data from file")
	}

	// Retrieve the encryption key from the ledger
	encryptionKey, err := sm.Ledger.GetEncryptionKey(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve encryption key from the ledger")
	}

	// Decrypt the token data
	decryptedTokenData, err := sm.EncryptionService.DecryptData(encryptedTokenData, encryptionKey)
	if err != nil {
		return nil, errors.New("failed to decrypt token data")
	}

	// Deserialize the token data
	var token Syn131Token
	err = json.Unmarshal(decryptedTokenData, &token)
	if err != nil {
		return nil, errors.New("failed to deserialize token data")
	}

	return &token, nil
}

// SaveTransaction saves a SYN131 transaction to persistent storage.
func (sm *StorageManager) SaveTransaction(transaction interface{}) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Serialize the transaction and encrypt it before saving
	transactionData, err := json.Marshal(transaction)
	if err != nil {
		return errors.New("failed to serialize transaction data")
	}

	encryptedTransactionData, encryptionKey, err := sm.EncryptionService.EncryptData(transactionData)
	if err != nil {
		return errors.New("failed to encrypt transaction data")
	}

	// Store the encrypted transaction data to a file
	transactionID := ""
	switch t := transaction.(type) {
	case *OwnershipTransaction:
		transactionID = t.TransactionID
	case *ShardedOwnershipTransaction:
		transactionID = t.TransactionID
	case *RentalTransaction:
		transactionID = t.TransactionID
	case *LeaseTransaction:
		transactionID = t.TransactionID
	case *PurchaseTransaction:
		transactionID = t.TransactionID
	default:
		return errors.New("unknown transaction type")
	}

	filePath := sm.storagePath + "/transactions/" + transactionID + ".json"
	err = os.WriteFile(filePath, []byte(encryptedTransactionData), 0644)
	if err != nil {
		return errors.New("failed to write transaction data to file")
	}

	// Store the encryption key securely in the ledger
	err = sm.Ledger.StoreEncryptionKey(transactionID, encryptionKey)
	if err != nil {
		return errors.New("failed to store encryption key in the ledger")
	}

	return nil
}

// LoadTransaction loads a SYN131 transaction from persistent storage.
func (sm *StorageManager) LoadTransaction(transactionID string) (interface{}, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the encrypted transaction data from the file
	filePath := sm.storagePath + "/transactions/" + transactionID + ".json"
	encryptedTransactionData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("failed to read transaction data from file")
	}

	// Retrieve the encryption key from the ledger
	encryptionKey, err := sm.Ledger.GetEncryptionKey(transactionID)
	if err != nil {
		return nil, errors.New("failed to retrieve encryption key from the ledger")
	}

	// Decrypt the transaction data
	decryptedTransactionData, err := sm.EncryptionService.DecryptData(encryptedTransactionData, encryptionKey)
	if err != nil {
		return nil, errors.New("failed to decrypt transaction data")
	}

	// Deserialize the transaction data
	var transaction interface{}
	switch {
	case isOwnershipTransaction(decryptedTransactionData):
		var t OwnershipTransaction
		err = json.Unmarshal(decryptedTransactionData, &t)
		transaction = t
	case isShardedOwnershipTransaction(decryptedTransactionData):
		var t ShardedOwnershipTransaction
		err = json.Unmarshal(decryptedTransactionData, &t)
		transaction = t
	case isRentalTransaction(decryptedTransactionData):
		var t RentalTransaction
		err = json.Unmarshal(decryptedTransactionData, &t)
		transaction = t
	case isLeaseTransaction(decryptedTransactionData):
		var t LeaseTransaction
		err = json.Unmarshal(decryptedTransactionData, &t)
		transaction = t
	case isPurchaseTransaction(decryptedTransactionData):
		var t PurchaseTransaction
		err = json.Unmarshal(decryptedTransactionData, &t)
		transaction = t
	default:
		return nil, errors.New("unknown transaction type")
	}

	if err != nil {
		return nil, errors.New("failed to deserialize transaction data")
	}

	return transaction, nil
}

// Helper functions to identify transaction types based on data
func isOwnershipTransaction(data []byte) bool {
	var transaction OwnershipTransaction
	return json.Unmarshal(data, &transaction) == nil
}

func isShardedOwnershipTransaction(data []byte) bool {
	var transaction ShardedOwnershipTransaction
	return json.Unmarshal(data, &transaction) == nil
}

func isRentalTransaction(data []byte) bool {
	var transaction RentalTransaction
	return json.Unmarshal(data, &transaction) == nil
}

func isLeaseTransaction(data []byte) bool {
	var transaction LeaseTransaction
	return json.Unmarshal(data, &transaction) == nil
}

func isPurchaseTransaction(data []byte) bool {
	var transaction PurchaseTransaction
	return json.Unmarshal(data, &transaction) == nil
}

// SaveMetadata saves asset metadata to persistent storage.
func (sm *StorageManager) SaveMetadata(metadata *AssetMetadata) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Serialize the metadata and encrypt it before saving
	metadataData, err := json.Marshal(metadata)
	if err != nil {
		return errors.New("failed to serialize metadata")
	}

	encryptedMetadataData, encryptionKey, err := sm.EncryptionService.EncryptData(metadataData)
	if err != nil {
		return errors.New("failed to encrypt metadata")
	}

	// Store the encrypted metadata data to a file
	filePath := sm.storagePath + "/metadata/" + metadata.ID + ".json"
	err = os.WriteFile(filePath, []byte(encryptedMetadataData), 0644)
	if err != nil {
		return errors.New("failed to write metadata to file")
	}

	// Store the encryption key securely in the ledger
	err = sm.Ledger.StoreEncryptionKey(metadata.ID, encryptionKey)
	if err != nil {
		return errors.New("failed to store encryption key in the ledger")
	}

	return nil
}

// LoadMetadata loads asset metadata from persistent storage.
func (sm *StorageManager) LoadMetadata(metadataID string) (*AssetMetadata, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the encrypted metadata data from the file
	filePath := sm.storagePath + "/metadata/" + metadataID + ".json"
	encryptedMetadataData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("failed to read metadata from file")
	}

	// Retrieve the encryption key from the ledger
	encryptionKey, err := sm.Ledger.GetEncryptionKey(metadataID)
	if err != nil {
		return nil, errors.New("failed to retrieve encryption key from the ledger")
	}

	// Decrypt the metadata data
	decryptedMetadataData, err := sm.EncryptionService.DecryptData(encryptedMetadataData, encryptionKey)
	if err != nil {
		return nil, errors.New("failed to decrypt metadata")
	}

	// Deserialize the metadata data
	var metadata AssetMetadata
	err = json.Unmarshal(decryptedMetadataData, &metadata)
	if err != nil {
		return nil, errors.New("failed to deserialize metadata")
	}

	return &metadata, nil
}
