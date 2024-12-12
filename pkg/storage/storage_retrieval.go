package storage

import (
    "fmt"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NewStorageRetrievalManager initializes a new storage retrieval manager
func NewStorageRetrievalManager(ledgerInstance *ledger.Ledger) *StorageRetrievalManager {
    return &StorageRetrievalManager{
        OnChainData:    make(map[string][]byte),
        OffChainUnits:  make(map[string]*OffchainStorage),
        LedgerInstance: ledgerInstance,
    }
}

// RetrieveOnChainData retrieves and decrypts data stored on-chain
func (srm *StorageRetrievalManager) RetrieveOnChainData(transactionHash string) ([]byte, error) {
    srm.mutex.Lock()
    defer srm.mutex.Unlock()

    data, exists := srm.OnChainData[transactionHash]
    if !exists {
        return nil, fmt.Errorf("no data found for transaction hash: %s", transactionHash)
    }

    decryptedData, err := common.DecryptData(string(data), common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt on-chain data: %v", err)
    }

    fmt.Printf("On-chain data retrieved for transaction hash: %s\n", transactionHash)
    return []byte(decryptedData), nil
}

// StoreOnChainData stores data on-chain and encrypts it before recording the transaction
func (srm *StorageRetrievalManager) StoreOnChainData(data []byte) (string, error) {
    srm.mutex.Lock()
    defer srm.mutex.Unlock()

    // Encrypt the data before storing it
    encryptedData, err := common.EncryptData(string(data), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt data: %v", err)
    }

    // Generate a transaction hash for the on-chain storage
    transactionHash := GenerateTransactionHash(encryptedData)

    // Store the encrypted data by its transaction hash
    srm.OnChainData[transactionHash] = []byte(encryptedData)

    // Record the storage transaction on the ledger
    err = srm.LedgerInstance.RecordTransaction(transactionHash, "on_chain_storage", "system", encryptedData)
    if err != nil {
        return "", fmt.Errorf("failed to record storage transaction in ledger: %v", err)
    }

    fmt.Printf("On-chain data stored with transaction hash: %s\n", transactionHash)
    return transactionHash, nil
}

// RetrieveOffChainData retrieves data from off-chain storage units
func (srm *StorageRetrievalManager) RetrieveOffChainData(storageID string) (*OffchainStorage, error) {
    srm.mutex.Lock()
    defer srm.mutex.Unlock()

    storageUnit, exists := srm.OffChainUnits[storageID]
    if !exists {
        return nil, fmt.Errorf("no off-chain storage unit found with ID: %s", storageID)
    }

    fmt.Printf("Off-chain data retrieved for storage ID: %s\n", storageID)
    return storageUnit, nil
}

// StoreOffChainData links an off-chain storage unit to the retrieval system
func (srm *StorageRetrievalManager) StoreOffChainData(storageUnit *OffchainStorage) error {
    srm.mutex.Lock()
    defer srm.mutex.Unlock()

    // Ensure the storage ID is unique
    if _, exists := srm.OffChainUnits[storageUnit.StorageID]; exists {
        return fmt.Errorf("off-chain storage unit with ID %s already exists", storageUnit.StorageID)
    }

    // Store the off-chain storage unit in the retrieval system
    srm.OffChainUnits[storageUnit.StorageID] = storageUnit

    // Record the off-chain storage unit in the ledger
    err := srm.LedgerInstance.RecordTransaction(storageUnit.StorageID, "off_chain_storage", storageUnit.Owner, storageUnit.EncryptedDetails)
    if err != nil {
        return fmt.Errorf("failed to record off-chain storage in ledger: %v", err)
    }

    fmt.Printf("Off-chain storage unit with ID %s linked to the system.\n", storageUnit.StorageID)
    return nil
}

// ValidateStorageTransaction ensures that storage transactions (on-chain or off-chain) are verified using the ledger
func (srm *StorageRetrievalManager) ValidateStorageTransaction(transactionHash string) (bool, error) {
    srm.mutex.Lock()
    defer srm.mutex.Unlock()

    valid, err := srm.LedgerInstance.VerifyTransaction(transactionHash)
    if err != nil {
        return false, fmt.Errorf("failed to validate storage transaction: %v", err)
    }

    if !valid {
        return false, fmt.Errorf("storage transaction with hash %s is invalid", transactionHash)
    }

    fmt.Printf("Storage transaction with hash %s validated successfully.\n", transactionHash)
    return true, nil
}
