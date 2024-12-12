package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/storage"
)

const (
    StorageEncryptionCheckInterval = 10 * time.Second // Interval for checking decentralized storage encryption
    MaxRetries                     = 3               // Maximum number of retry attempts for encryption
)

// DecentralizedStorageEncryptionAutomation automates the process of encrypting decentralized storage data
type DecentralizedStorageEncryptionAutomation struct {
    consensusSystem *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance  *ledger.Ledger               // Ledger for logging encryption events
    stateMutex      *sync.RWMutex                // Mutex for thread-safe access
    encryptionRetry map[string]int               // Counter for retrying encryption on failure
}

// NewDecentralizedStorageEncryptionAutomation initializes the automation for decentralized storage encryption
func NewDecentralizedStorageEncryptionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DecentralizedStorageEncryptionAutomation {
    return &DecentralizedStorageEncryptionAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        encryptionRetry: make(map[string]int),
    }
}

// StartStorageEncryptionCheck starts the continuous loop for regularly checking and enforcing encryption on decentralized storage
func (automation *DecentralizedStorageEncryptionAutomation) StartStorageEncryptionCheck() {
    ticker := time.NewTicker(StorageEncryptionCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkAndEncryptStorageData()
        }
    }()
}

// checkAndEncryptStorageData retrieves unencrypted data and applies encryption
func (automation *DecentralizedStorageEncryptionAutomation) checkAndEncryptStorageData() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch list of unencrypted data from decentralized storage
    unencryptedDataList := automation.consensusSystem.GetUnencryptedStorageData()

    if len(unencryptedDataList) > 0 {
        for _, data := range unencryptedDataList {
            fmt.Printf("Encrypting storage data with ID %s.\n", data.ID)
            automation.encryptStorageData(data)
        }
    } else {
        fmt.Println("No unencrypted storage data found.")
    }
}

// encryptStorageData encrypts the provided data and stores it back into decentralized storage
func (automation *DecentralizedStorageEncryptionAutomation) encryptStorageData(data common.StorageData) {
    // Encrypt the data
    encryptedData, err := encryption.EncryptData(data.Content)
    if err != nil {
        fmt.Printf("Error encrypting data with ID %s: %v. Retrying...\n", data.ID, err)
        automation.retryEncryption(data)
        return
    }

    // Store encrypted data in decentralized storage
    success := storage.StoreEncryptedData(data.ID, encryptedData)
    if success {
        fmt.Printf("Data with ID %s encrypted and stored successfully.\n", data.ID)
        automation.logEncryptionEvent(data.ID)
        automation.resetEncryptionRetry(data.ID)
    } else {
        fmt.Printf("Failed to store encrypted data with ID %s. Retrying...\n", data.ID)
        automation.retryEncryption(data)
    }
}

// retryEncryption attempts to retry encryption for data a limited number of times
func (automation *DecentralizedStorageEncryptionAutomation) retryEncryption(data common.StorageData) {
    automation.encryptionRetry[data.ID]++
    if automation.encryptionRetry[data.ID] < MaxRetries {
        automation.encryptStorageData(data)
    } else {
        fmt.Printf("Max retries reached for encrypting data with ID %s. Encryption failed.\n", data.ID)
        automation.logEncryptionFailure(data.ID)
    }
}

// resetEncryptionRetry resets the retry counter for encryption of specific data
func (automation *DecentralizedStorageEncryptionAutomation) resetEncryptionRetry(dataID string) {
    automation.encryptionRetry[dataID] = 0
}

// logEncryptionEvent logs a successful encryption event into the ledger
func (automation *DecentralizedStorageEncryptionAutomation) logEncryptionEvent(dataID string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("storage-encryption-%s", dataID),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Encryption",
        Status:    "Completed",
        Details:   fmt.Sprintf("Data with ID %s was successfully encrypted and stored.", dataID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with encryption event for data ID %s.\n", dataID)
}

// logEncryptionFailure logs an encryption failure event into the ledger
func (automation *DecentralizedStorageEncryptionAutomation) logEncryptionFailure(dataID string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("storage-encryption-failure-%s", dataID),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Encryption Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to encrypt data with ID %s after maximum retries.", dataID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with encryption failure for data ID %s.\n", dataID)
}

// checkStorageIntegrity ensures the integrity of encrypted data in storage
func (automation *DecentralizedStorageEncryptionAutomation) checkStorageIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Validate the integrity of all encrypted storage data
    integrityValid := automation.consensusSystem.ValidateStorageIntegrity()
    if !integrityValid {
        fmt.Println("Storage integrity breach detected. Re-encrypting affected data.")
        automation.checkAndEncryptStorageData()
    } else {
        fmt.Println("All storage data integrity is valid.")
    }
}
