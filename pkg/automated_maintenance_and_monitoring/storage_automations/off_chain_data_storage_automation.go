package automations

import (
    "fmt"
    "sync"
    "time"
    "os"
    "io/ioutil"
    "path/filepath"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    OffChainStorageInterval     = 12 * time.Hour  // Interval for checking off-chain storage needs
    MaxOffChainStorageRetries   = 3               // Maximum retries for off-chain storage issues
    SubBlocksPerBlock           = 1000            // Number of sub-blocks in a block
    OffChainStorageThreshold    = 100 * 1024 * 1024 // 100 MB threshold for off-chain storage
    OffChainStorageDir          = "/data/offchain_storage" // Directory to store off-chain files
)

// OffChainDataStorageAutomation manages the off-chain storage process
type OffChainDataStorageAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger for logging off-chain storage events
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    offChainRetryCount     map[string]int               // Counter for retrying failed off-chain storage operations
    storageCycleCount      int                          // Counter for off-chain storage monitoring cycles
}

// NewOffChainDataStorageAutomation initializes the automation for off-chain data storage
func NewOffChainDataStorageAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *OffChainDataStorageAutomation {
    return &OffChainDataStorageAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        offChainRetryCount: make(map[string]int),
        storageCycleCount:  0,
    }
}

// StartOffChainStorageMonitoring starts the continuous loop for monitoring off-chain storage needs
func (automation *OffChainDataStorageAutomation) StartOffChainStorageMonitoring() {
    ticker := time.NewTicker(OffChainStorageInterval)

    go func() {
        for range ticker.C {
            automation.monitorOffChainStorage()
        }
    }()
}

// monitorOffChainStorage checks for large on-chain data requiring off-chain storage
func (automation *OffChainDataStorageAutomation) monitorOffChainStorage() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch reports for data that potentially needs off-chain storage
    storageReports := automation.consensusSystem.FetchOffChainStorageReports()

    for _, report := range storageReports {
        if automation.isOffChainStorageRequired(report) {
            fmt.Printf("Off-chain storage required for data %s. Initiating transfer to off-chain storage.\n", report.DataID)
            automation.applyOffChainStorage(report)
        } else {
            fmt.Printf("No off-chain storage required for data %s.\n", report.DataID)
        }
    }

    automation.storageCycleCount++
    fmt.Printf("Off-chain storage cycle #%d completed.\n", automation.storageCycleCount)

    if automation.storageCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeOffChainStorageCycle()
    }
}

// isOffChainStorageRequired checks if the data exceeds the threshold for off-chain storage
func (automation *OffChainDataStorageAutomation) isOffChainStorageRequired(report common.StorageReport) bool {
    if report.DataSize >= OffChainStorageThreshold {
        fmt.Printf("Data %s exceeds size threshold (size: %d bytes), triggering off-chain storage.\n", report.DataID, report.DataSize)
        return true
    }
    return false
}

// applyOffChainStorage transfers data to off-chain storage and maintains on-chain references
func (automation *OffChainDataStorageAutomation) applyOffChainStorage(report common.StorageReport) {
    encryptedData := automation.encryptOffChainData(report)

    // Store the encrypted data in the off-chain storage directory
    err := automation.storeDataOffChain(report.DataID, encryptedData.EncryptedData)
    if err != nil {
        fmt.Printf("Error storing data %s off-chain: %v. Retrying...\n", report.DataID, err)
        automation.retryOffChainStorage(report)
        return
    }

    storageReference := filepath.Join(OffChainStorageDir, fmt.Sprintf("%s.enc", report.DataID))

    fmt.Printf("Data %s successfully stored off-chain. Reference: %s\n", report.DataID, storageReference)
    automation.logOffChainStorageEvent(report, "Stored off-chain", storageReference)
    automation.resetOffChainStorageRetry(report.DataID)
}

// storeDataOffChain writes the encrypted data to a file in the off-chain storage directory
func (automation *OffChainDataStorageAutomation) storeDataOffChain(dataID string, data []byte) error {
    if _, err := os.Stat(OffChainStorageDir); os.IsNotExist(err) {
        err := os.MkdirAll(OffChainStorageDir, os.ModePerm)
        if err != nil {
            return fmt.Errorf("failed to create off-chain storage directory: %v", err)
        }
    }

    filePath := filepath.Join(OffChainStorageDir, fmt.Sprintf("%s.enc", dataID))
    err := ioutil.WriteFile(filePath, data, 0644)
    if err != nil {
        return fmt.Errorf("failed to write data to off-chain storage: %v", err)
    }

    return nil
}

// retryOffChainStorage retries the off-chain storage process in case of failure
func (automation *OffChainDataStorageAutomation) retryOffChainStorage(report common.StorageReport) {
    automation.offChainRetryCount[report.DataID]++
    if automation.offChainRetryCount[report.DataID] < MaxOffChainStorageRetries {
        automation.applyOffChainStorage(report)
    } else {
        fmt.Printf("Max retries reached for off-chain storage of data %s. Operation failed.\n", report.DataID)
        automation.logOffChainStorageFailure(report)
    }
}

// resetOffChainStorageRetry resets the retry count for off-chain storage operations
func (automation *OffChainDataStorageAutomation) resetOffChainStorageRetry(dataID string) {
    automation.offChainRetryCount[dataID] = 0
}

// finalizeOffChainStorageCycle finalizes the off-chain storage cycle and logs the result in the ledger
func (automation *OffChainDataStorageAutomation) finalizeOffChainStorageCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeOffChainStorageCycle()
    if success {
        fmt.Println("Off-chain storage cycle finalized successfully.")
        automation.logOffChainStorageCycleFinalization()
    } else {
        fmt.Println("Error finalizing off-chain storage cycle.")
    }
}

// logOffChainStorageEvent logs a successful off-chain storage event into the ledger
func (automation *OffChainDataStorageAutomation) logOffChainStorageEvent(report common.StorageReport, eventType, reference string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("offchain-storage-event-%s-%s", report.DataID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Off-Chain Storage Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Data %s stored off-chain. Reference: %s", report.DataID, reference),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with off-chain storage event for data %s.\n", report.DataID)
}

// logOffChainStorageFailure logs the failure of an off-chain storage attempt into the ledger
func (automation *OffChainDataStorageAutomation) logOffChainStorageFailure(report common.StorageReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("offchain-storage-failure-%s", report.DataID),
        Timestamp: time.Now().Unix(),
        Type:      "Off-Chain Storage Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Off-chain storage failed for data %s after maximum retries.", report.DataID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with off-chain storage failure for data %s.\n", report.DataID)
}

// logOffChainStorageCycleFinalization logs the finalization of an off-chain storage cycle into the ledger
func (automation *OffChainDataStorageAutomation) logOffChainStorageCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("offchain-storage-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Off-Chain Storage Cycle Finalization",
        Status:    "Finalized",
        Details:   "Off-chain storage cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with off-chain storage cycle finalization.")
}

// encryptOffChainData encrypts data before sending it to off-chain storage
func (automation *OffChainDataStorageAutomation) encryptOffChainData(report common.StorageReport) common.StorageReport {
    encryptedData, err := encryption.EncryptData(report.Data)
    if err != nil {
        fmt.Println("Error encrypting data for off-chain storage:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("Data successfully encrypted for off-chain storage. DataID:", report.DataID)
    return report
}
