package syn300

import (
	"errors"
	"sync"
	"time"
)

// Syn300StorageManager handles storage, retrieval, and encryption of token data for SYN300 tokens.
type Syn300StorageManager struct {
	mutex   sync.RWMutex
	Ledger  *ledger.Ledger
	Storage map[string]TokenStorage // Maps token ID to its encrypted storage data
}

// TokenStorage stores the encrypted metadata and balance of a SYN300 token.
type TokenStorage struct {
	TokenID      string
	EncryptedData string
	LastUpdated  time.Time
}

// NewSyn300StorageManager creates a new instance of Syn300StorageManager.
func NewSyn300StorageManager(ledger *ledger.Ledger) *Syn300StorageManager {
	return &Syn300StorageManager{
		Ledger:  ledger,
		Storage: make(map[string]TokenStorage),
	}
}

// SaveTokenData securely encrypts and stores the metadata and balance of a SYN300 token in the storage.
func (sm *Syn300StorageManager) SaveTokenData(tokenID, metadata string, balance map[string]uint64, encryptionKey string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Combine metadata and balance for encryption
	combinedData := struct {
		Metadata string
		Balance  map[string]uint64
	}{
		Metadata: metadata,
		Balance:  balance,
	}

	// Encrypt the combined data
	encryptedData, err := encryption.EncryptStruct(combinedData, encryptionKey)
	if err != nil {
		return errors.New("failed to encrypt token data")
	}

	// Store the encrypted data in the manager
	sm.Storage[tokenID] = TokenStorage{
		TokenID:      tokenID,
		EncryptedData: encryptedData,
		LastUpdated:  time.Now(),
	}

	// Store in the ledger for auditability and redundancy
	if err := sm.Ledger.StoreData(tokenID, encryptedData); err != nil {
		return errors.New("failed to store token data in the ledger")
	}

	return nil
}

// RetrieveTokenData decrypts and retrieves the metadata and balance of a SYN300 token from storage.
func (sm *Syn300StorageManager) RetrieveTokenData(tokenID, encryptionKey string) (string, map[string]uint64, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Retrieve the encrypted data from storage
	storage, exists := sm.Storage[tokenID]
	if !exists {
		return "", nil, errors.New("token data not found")
	}

	// Decrypt the data
	var decryptedData struct {
		Metadata string
		Balance  map[string]uint64
	}
	err := encryption.DecryptStruct(storage.EncryptedData, encryptionKey, &decryptedData)
	if err != nil {
		return "", nil, errors.New("failed to decrypt token data")
	}

	return decryptedData.Metadata, decryptedData.Balance, nil
}

// ValidateStorageConsistency ensures the consistency of stored data by verifying it with the Synnergy Consensus mechanism.
func (sm *Syn300StorageManager) ValidateStorageConsistency(tokenID string) (bool, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Verify data consistency using Synnergy Consensus
	isValid, err := consensus.ValidateStorage(tokenID, "")
	if err != nil || !isValid {
		return false, errors.New("storage validation failed")
	}

	return true, nil
}

// UpdateTokenData updates the metadata and balance of a SYN300 token in storage.
func (sm *Syn300StorageManager) UpdateTokenData(tokenID, metadata string, balance map[string]uint64, encryptionKey string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve current data for verification before updating
	_, _, err := sm.RetrieveTokenData(tokenID, encryptionKey)
	if err != nil {
		return errors.New("failed to retrieve token data for update")
	}

	// Encrypt and save the updated data
	return sm.SaveTokenData(tokenID, metadata, balance, encryptionKey)
}

// DeleteTokenData securely removes the token data from storage and ledger.
func (sm *Syn300StorageManager) DeleteTokenData(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Remove from storage
	delete(sm.Storage, tokenID)

	// Remove from ledger for full deletion
	if err := sm.Ledger.DeleteData(tokenID); err != nil {
		return errors.New("failed to delete token data from ledger")
	}

	return nil
}

// BackupTokenData creates an encrypted backup of token data for disaster recovery purposes.
func (sm *Syn300StorageManager) BackupTokenData(tokenID, backupPath, encryptionKey string) error {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Retrieve the current token data
	storage, exists := sm.Storage[tokenID]
	if !exists {
		return errors.New("token data not found for backup")
	}

	// Backup the encrypted data to the specified path
	err := encryption.BackupData(storage.EncryptedData, backupPath, encryptionKey)
	if err != nil {
		return errors.New("failed to create backup for token data")
	}

	return nil
}

// RestoreTokenData restores token data from an encrypted backup.
func (sm *Syn300StorageManager) RestoreTokenData(backupPath, tokenID, encryptionKey string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Restore the encrypted data from the backup
	encryptedData, err := encryption.RestoreData(backupPath, encryptionKey)
	if err != nil {
		return errors.New("failed to restore token data from backup")
	}

	// Store the restored data
	sm.Storage[tokenID] = TokenStorage{
		TokenID:      tokenID,
		EncryptedData: encryptedData,
		LastUpdated:  time.Now(),
	}

	// Store in ledger for redundancy
	if err := sm.Ledger.StoreData(tokenID, encryptedData); err != nil {
		return errors.New("failed to store restored token data in ledger")
	}

	return nil
}

// AuditStorage performs a security audit on all stored SYN300 token data.
func (sm *Syn300StorageManager) AuditStorage() ([]string, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Fetch audit logs from the ledger
	auditLogs, err := sm.Ledger.GetAuditLogs()
	if err != nil {
		return nil, errors.New("failed to retrieve audit logs")
	}

	return auditLogs, nil
}
