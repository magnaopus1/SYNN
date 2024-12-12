package storage

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NewOffchainStorageManager initializes a new OffchainStorageManager
func NewOffchainStorageManager(ledgerInstance *ledger.Ledger) *OffchainStorageManager {
    return &OffchainStorageManager{
        StorageUnits:   make(map[string]*OffchainStorage),
        LedgerInstance: ledgerInstance,
    }
}

// AddStorageUnit allows a user to add a new off-chain storage unit
func (osm *OffchainStorageManager) AddStorageUnit(owner, location string, capacityGB int, details string) (string, error) {
    osm.mutex.Lock()
    defer osm.mutex.Unlock()

    // Generate a unique ID for the storage unit
    storageID := osm.generateStorageID(owner, location, capacityGB)

    // Encrypt the storage details
    encryptedDetails, err := encryption.EncryptData(details, common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt storage details: %v", err)
    }

    // Create the new storage unit
    newStorage := &OffchainStorage{
        StorageID:       storageID,
        Owner:           owner,
        Location:        location,
        CapacityGB:      capacityGB,
        UsedCapacityGB:  0,
        EncryptedDetails: encryptedDetails,
        Active:          true,
    }

    // Add the storage unit to the manager
    osm.StorageUnits[storageID] = newStorage

    // Log the storage unit creation in the ledger
    err = osm.logStorageUnitToLedger(newStorage)
    if err != nil {
        return "", fmt.Errorf("failed to log storage unit to ledger: %v", err)
    }

    fmt.Printf("Off-chain storage unit %s created by %s.\n", storageID, owner)
    return storageID, nil
}

// UpdateStorageCapacity updates the used capacity of an off-chain storage unit
func (osm *OffchainStorageManager) UpdateStorageCapacity(storageID string, usedCapacityGB int) error {
    osm.mutex.Lock()
    defer osm.mutex.Unlock()

    storage, exists := osm.StorageUnits[storageID]
    if !exists {
        return fmt.Errorf("storage unit not found")
    }

    if usedCapacityGB > storage.CapacityGB {
        return fmt.Errorf("used capacity exceeds total capacity")
    }

    storage.UsedCapacityGB = usedCapacityGB

    // Log the capacity update in the ledger
    err := osm.logStorageUnitToLedger(storage)
    if err != nil {
        return fmt.Errorf("failed to log storage update to ledger: %v", err)
    }

    fmt.Printf("Off-chain storage unit %s updated to %d GB used.\n", storageID, usedCapacityGB)
    return nil
}

// ViewStorageDetails allows viewing the details of an off-chain storage unit
func (osm *OffchainStorageManager) ViewStorageDetails(storageID string) (*OffchainStorage, error) {
    osm.mutex.Lock()
    defer osm.mutex.Unlock()

    storage, exists := osm.StorageUnits[storageID]
    if !exists {
        return nil, fmt.Errorf("storage unit not found")
    }

    fmt.Printf("Storage details for unit %s viewed.\n", storageID)
    return storage, nil
}

// RemoveStorageUnit deactivates a storage unit from the off-chain system
func (osm *OffchainStorageManager) RemoveStorageUnit(storageID string) error {
    osm.mutex.Lock()
    defer osm.mutex.Unlock()

    storage, exists := osm.StorageUnits[storageID]
    if !exists {
        return fmt.Errorf("storage unit not found")
    }

    // Mark the storage unit as inactive
    storage.Active = false

    // Log the removal in the ledger
    err := osm.logStorageUnitToLedger(storage)
    if err != nil {
        return fmt.Errorf("failed to log storage removal to ledger: %v", err)
    }

    fmt.Printf("Off-chain storage unit %s deactivated.\n", storageID)
    return nil
}

// generateStorageID generates a unique ID for each off-chain storage unit
func (osm *OffchainStorageManager) generateStorageID(owner, location string, capacityGB int) string {
    hashInput := fmt.Sprintf("%s%s%d%d", owner, location, capacityGB, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// logStorageUnitToLedger logs the off-chain storage details to the ledger
func (osm *OffchainStorageManager) logStorageUnitToLedger(storage *OffchainStorage) error {
    storageDetails := fmt.Sprintf("StorageID: %s, Owner: %s, Location: %s, Capacity: %d GB, Used: %d GB",
        storage.StorageID, storage.Owner, storage.Location, storage.CapacityGB, storage.UsedCapacityGB)

    encryptedDetails, err := encryption.EncryptData(storageDetails, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt storage details: %v", err)
    }

    err = osm.LedgerInstance.RecordTransaction(storage.StorageID, "offchain_storage", storage.Owner, encryptedDetails)
    if err != nil {
        return fmt.Errorf("failed to log storage unit to ledger: %v", err)
    }

    fmt.Printf("Off-chain storage unit %s logged to the ledger.\n", storage.StorageID)
    return nil
}
