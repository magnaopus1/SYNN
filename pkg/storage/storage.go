package storage

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NewStorageManager initializes a new StorageManager
func NewStorageManager(ledgerInstance *ledger.Ledger) *StorageManager {
    return &StorageManager{
        LedgerInstance: ledgerInstance,
        StorageMap:     make(map[string]StorageEntry),
    }
}

// CreateStorageEntry creates and stores a new encrypted storage entry, and logs it to the ledger
func (sm *StorageManager) CreateStorageEntry(ownerID string, data []byte, expirationDuration time.Duration) (string, error) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    // Generate unique storage ID
    storageID := common.GenerateStorageID(ownerID, time.Now().UnixNano())

    // Encrypt the data
    encryptedData, err := common.EncryptData(string(data), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt storage data: %v", err)
    }

    // Create new storage entry
    newEntry := StorageEntry{
        StorageID:    storageID,
        OwnerID:      ownerID,
        Data:         []byte(encryptedData),
        CreatedAt:    time.Now(),
        ExpiresAt:    time.Now().Add(expirationDuration),
    }

    // Add to storage map
    sm.StorageMap[storageID] = newEntry

    // Log to ledger
    err = sm.logStorageToLedger(newEntry)
    if err != nil {
        return "", fmt.Errorf("failed to log storage entry to ledger: %v", err)
    }

    fmt.Printf("Storage entry %s created for owner %s.\n", storageID, ownerID)
    return storageID, nil
}

// RetrieveStorageEntry retrieves and decrypts the storage entry for a given storage ID
func (sm *StorageManager) RetrieveStorageEntry(storageID string) ([]byte, error) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    entry, exists := sm.StorageMap[storageID]
    if !exists {
        return nil, fmt.Errorf("storage entry %s not found", storageID)
    }

    if time.Now().After(entry.ExpiresAt) {
        return nil, fmt.Errorf("storage entry %s has expired", storageID)
    }

    // Decrypt the data
    decryptedData, err := common.DecryptData(string(entry.Data), common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt storage data: %v", err)
    }

    fmt.Printf("Storage entry %s retrieved for owner %s.\n", storageID, entry.OwnerID)
    return []byte(decryptedData), nil
}

// DeleteStorageEntry deletes a storage entry for a given storage ID and logs it in the ledger
func (sm *StorageManager) DeleteStorageEntry(storageID string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    entry, exists := sm.StorageMap[storageID]
    if !exists {
        return fmt.Errorf("storage entry %s not found", storageID)
    }

    // Remove entry from map
    delete(sm.StorageMap, storageID)

    // Log deletion to ledger
    err := sm.logStorageDeletionToLedger(storageID, entry.OwnerID)
    if err != nil {
        return fmt.Errorf("failed to log storage deletion to ledger: %v", err)
    }

    fmt.Printf("Storage entry %s deleted for owner %s.\n", storageID, entry.OwnerID)
    return nil
}

// logStorageToLedger logs the creation of a storage entry in the ledger
func (sm *StorageManager) logStorageToLedger(entry StorageEntry) error {
    storageRecord := fmt.Sprintf("StorageID: %s, OwnerID: %s, CreatedAt: %s, ExpiresAt: %s", 
        entry.StorageID, entry.OwnerID, entry.CreatedAt.String(), entry.ExpiresAt.String())

    encryptedRecord, err := common.EncryptData(storageRecord, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt storage log: %v", err)
    }

    return sm.LedgerInstance.RecordTransaction(entry.StorageID, "storage_creation", entry.OwnerID, encryptedRecord)
}

// logStorageDeletionToLedger logs the deletion of a storage entry in the ledger
func (sm *StorageManager) logStorageDeletionToLedger(storageID, ownerID string) error {
    deletionRecord := fmt.Sprintf("StorageID: %s, OwnerID: %s, DeletionTime: %s", 
        storageID, ownerID, time.Now().String())

    encryptedRecord, err := common.EncryptData(deletionRecord, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt storage deletion log: %v", err)
    }

    return sm.LedgerInstance.RecordTransaction(storageID, "storage_deletion", ownerID, encryptedRecord)
}
