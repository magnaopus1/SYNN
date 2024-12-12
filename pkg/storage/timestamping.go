package storage

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "synnergy_network/pkg/ledger"
)



// NewTimestampManager initializes a new TimestampManager instance
func NewTimestampManager(ledgerInstance *ledger.Ledger) *TimestampManager {
    return &TimestampManager{
        LedgerInstance: ledgerInstance,
    }
}

// TimestampData creates a timestamp for the given data and stores it immutably in the ledger
func (tm *TimestampManager) TimestampData(ownerID string, data []byte) (string, error) {
    // Generate a hash for the data
    dataHash := tm.generateDataHash(data)

    // Create a unique timestamp entry
    timestamp := TimestampEntry{
        DataHash:   dataHash,
        Timestamp:  time.Now(),
        OwnerID:    ownerID,
    }

    // Encrypt the timestamp data
    encryptedTimestamp, err := common.EncryptData([]byte(fmt.Sprintf("%+v", timestamp)), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt timestamp data: %v", err)
    }

    // Log the timestamp to the ledger
    err = tm.logTimestampToLedger(timestamp, encryptedTimestamp)
    if err != nil {
        return "", fmt.Errorf("failed to log timestamp to ledger: %v", err)
    }

    fmt.Printf("Timestamp for data hash %s recorded by owner %s at %s.\n", timestamp.DataHash, ownerID, timestamp.Timestamp)
    return dataHash, nil
}

// ValidateTimestamp checks if a timestamped data hash exists in the ledger
func (tm *TimestampManager) ValidateTimestamp(data []byte) (bool, error) {
    dataHash := tm.generateDataHash(data)

    // Check the ledger for the existence of the timestamp
    exists, err := tm.LedgerInstance.CheckRecord(dataHash)
    if err != nil {
        return false, fmt.Errorf("failed to check timestamp in ledger: %v", err)
    }

    if exists {
        fmt.Printf("Data hash %s is valid and was previously timestamped.\n", dataHash)
        return true, nil
    }

    fmt.Printf("Data hash %s was not found in the ledger.\n", dataHash)
    return false, nil
}

// generateDataHash generates a SHA-256 hash of the provided data
func (tm *TimestampManager) generateDataHash(data []byte) string {
    hash := sha256.New()
    hash.Write(data)
    return hex.EncodeToString(hash.Sum(nil))
}

// logTimestampToLedger logs the timestamp data into the ledger
func (tm *TimestampManager) logTimestampToLedger(timestamp TimestampEntry, encryptedTimestamp []byte) error {
    record := fmt.Sprintf("DataHash: %s, OwnerID: %s, Timestamp: %s",
        timestamp.DataHash, timestamp.OwnerID, timestamp.Timestamp.String())

    encryptedRecord, err := commonn.EncryptData([]byte(record), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt ledger record: %v", err)
    }

    // Record the timestamp in the ledger as a transaction
    return tm.LedgerInstance.RecordTransaction(timestamp.DataHash, "timestamp_record", timestamp.OwnerID, encryptedRecord)
}

