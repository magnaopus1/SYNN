package automations

import (
    "fmt"
    "sync"
    "time"
    "os"
    "crypto/sha256"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    IntegrityCheckInterval      = 12 * time.Hour // Interval for checking storage integrity
    SubBlocksPerBlock           = 1000           // Number of sub-blocks in a block
    IntegrityRetryLimit         = 3              // Maximum retry attempts for integrity check failures
)

// StorageIntegrityCheckAutomation handles periodic storage integrity checks for on-chain and off-chain storage
type StorageIntegrityCheckAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger for logging integrity check events
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    integrityRetryCount   map[string]int               // Counter for retrying failed integrity checks
    checkCycleCount       int                          // Counter for integrity check cycles
}

// NewStorageIntegrityCheckAutomation initializes the automation for storage integrity checks
func NewStorageIntegrityCheckAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *StorageIntegrityCheckAutomation {
    return &StorageIntegrityCheckAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        integrityRetryCount: make(map[string]int),
        checkCycleCount:     0,
    }
}

// StartStorageIntegrityCheck starts the continuous loop for checking storage integrity
func (automation *StorageIntegrityCheckAutomation) StartStorageIntegrityCheck() {
    ticker := time.NewTicker(IntegrityCheckInterval)

    go func() {
        for range ticker.C {
            automation.runStorageIntegrityCheck()
        }
    }()
}

// runStorageIntegrityCheck verifies the integrity of on-chain and off-chain storage by comparing hashes
func (automation *StorageIntegrityCheckAutomation) runStorageIntegrityCheck() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityReports := automation.consensusSystem.FetchStorageIntegrityReports()

    for _, report := range integrityReports {
        if automation.isIntegrityCompromised(report) {
            fmt.Printf("Integrity compromised for data %s. Initiating integrity recovery.\n", report.DataID)
            automation.recoverIntegrity(report)
        } else {
            fmt.Printf("Integrity verified for data %s.\n", report.DataID)
        }
    }

    automation.checkCycleCount++
    fmt.Printf("Storage integrity check cycle #%d completed.\n", automation.checkCycleCount)

    if automation.checkCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeIntegrityCheckCycle()
    }
}

// isIntegrityCompromised checks if the data has been altered by comparing its stored hash with a newly computed hash
func (automation *StorageIntegrityCheckAutomation) isIntegrityCompromised(report common.StorageReport) bool {
    computedHash := automation.computeHash(report.Data)

    if report.StoredHash != computedHash {
        fmt.Printf("Integrity check failed for data %s. Expected hash: %s, computed hash: %s.\n", report.DataID, report.StoredHash, computedHash)
        return true
    }

    return false
}

// computeHash computes the SHA-256 hash of the given data
func (automation *StorageIntegrityCheckAutomation) computeHash(data []byte) string {
    hash := sha256.Sum256(data)
    return fmt.Sprintf("%x", hash)
}

// recoverIntegrity attempts to recover data integrity by restoring from backups or reverting to a valid state
func (automation *StorageIntegrityCheckAutomation) recoverIntegrity(report common.StorageReport) {
    err := automation.triggerDataRecovery(report)
    if err != nil {
        fmt.Printf("Data recovery failed for data %s: %v. Retrying...\n", report.DataID, err)
        automation.retryIntegrityRecovery(report)
        return
    }

    fmt.Printf("Integrity successfully recovered for data %s.\n", report.DataID)
    automation.logIntegrityEvent(report, "Recovered")
    automation.resetIntegrityRetry(report.DataID)
}

// triggerDataRecovery handles the actual recovery of data integrity, either through consensus rollback or backups
func (automation *StorageIntegrityCheckAutomation) triggerDataRecovery(report common.StorageReport) error {
    fmt.Println("Triggering data recovery for integrity compromise...")

    // Simulate recovery logic
    success := automation.consensusSystem.TriggerDataRecovery(report)

    if success {
        fmt.Println("Data recovery triggered successfully.")
        return nil
    }

    return fmt.Errorf("data recovery failed for data %s", report.DataID)
}

// retryIntegrityRecovery retries the integrity recovery process in case of failure
func (automation *StorageIntegrityCheckAutomation) retryIntegrityRecovery(report common.StorageReport) {
    automation.integrityRetryCount[report.DataID]++
    if automation.integrityRetryCount[report.DataID] < IntegrityRetryLimit {
        automation.recoverIntegrity(report)
    } else {
        fmt.Printf("Max retries reached for data recovery on data %s. Integrity compromised.\n", report.DataID)
        automation.logIntegrityFailure(report)
    }
}

// resetIntegrityRetry resets the retry count for data integrity recovery
func (automation *StorageIntegrityCheckAutomation) resetIntegrityRetry(dataID string) {
    automation.integrityRetryCount[dataID] = 0
}

// logIntegrityEvent logs a successful integrity recovery event into the ledger
func (automation *StorageIntegrityCheckAutomation) logIntegrityEvent(report common.StorageReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("integrity-event-%s-%s", report.DataID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Integrity Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Integrity %s for data %s.", eventType, report.DataID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with integrity event for data %s.\n", report.DataID)
}

// logIntegrityFailure logs a failure to recover data integrity into the ledger
func (automation *StorageIntegrityCheckAutomation) logIntegrityFailure(report common.StorageReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("integrity-failure-%s", report.DataID),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Integrity Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to recover integrity for data %s after maximum retries.", report.DataID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with integrity failure for data %s.\n", report.DataID)
}

// finalizeIntegrityCheckCycle finalizes the storage integrity check cycle and logs the result in the ledger
func (automation *StorageIntegrityCheckAutomation) finalizeIntegrityCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeIntegrityCheckCycle()
    if success {
        fmt.Println("Storage integrity check cycle finalized successfully.")
        automation.logIntegrityCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing storage integrity check cycle.")
    }
}

// logIntegrityCheckCycleFinalization logs the finalization of an integrity check cycle into the ledger
func (automation *StorageIntegrityCheckAutomation) logIntegrityCheckCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("integrity-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Integrity Check Cycle Finalization",
        Status:    "Finalized",
        Details:   "Storage integrity check cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with integrity check cycle finalization.")
}
