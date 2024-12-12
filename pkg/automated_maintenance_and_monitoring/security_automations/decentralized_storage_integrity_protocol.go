package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/storage"
)

const (
    StorageIntegrityCheckInterval = 15 * time.Second // Interval for checking the integrity of decentralized storage
    MaxIntegrityRetries           = 3                // Maximum number of retry attempts for fixing integrity issues
)

// DecentralizedStorageIntegrityAutomation automates the process of monitoring and ensuring the integrity of decentralized storage
type DecentralizedStorageIntegrityAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging integrity events
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    integrityRetryMap map[string]int               // Counter for retrying integrity fixes
}

// NewDecentralizedStorageIntegrityAutomation initializes the automation for storage integrity checks
func NewDecentralizedStorageIntegrityAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DecentralizedStorageIntegrityAutomation {
    return &DecentralizedStorageIntegrityAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        integrityRetryMap: make(map[string]int),
    }
}

// StartStorageIntegrityCheck starts the continuous loop for regularly checking the integrity of decentralized storage
func (automation *DecentralizedStorageIntegrityAutomation) StartStorageIntegrityCheck() {
    ticker := time.NewTicker(StorageIntegrityCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkAndEnsureStorageIntegrity()
        }
    }()
}

// checkAndEnsureStorageIntegrity checks storage for integrity issues and attempts to fix them
func (automation *DecentralizedStorageIntegrityAutomation) checkAndEnsureStorageIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch list of data that may have integrity issues
    integrityIssues := automation.consensusSystem.GetStorageIntegrityIssues()

    if len(integrityIssues) > 0 {
        for _, data := range integrityIssues {
            fmt.Printf("Integrity issue detected for data ID %s. Attempting to resolve.\n", data.ID)
            automation.fixIntegrityIssue(data)
        }
    } else {
        fmt.Println("No storage integrity issues detected.")
    }
}

// fixIntegrityIssue attempts to resolve the integrity issue for the specified data
func (automation *DecentralizedStorageIntegrityAutomation) fixIntegrityIssue(data common.StorageData) {
    // Validate the integrity using the consensus system
    integrityFixed := automation.consensusSystem.FixStorageIntegrity(data)

    if integrityFixed {
        fmt.Printf("Integrity issue resolved for data ID %s.\n", data.ID)
        automation.logIntegrityFixEvent(data.ID)
        automation.resetIntegrityRetry(data.ID)
    } else {
        fmt.Printf("Failed to resolve integrity issue for data ID %s. Retrying...\n", data.ID)
        automation.retryIntegrityFix(data)
    }
}

// retryIntegrityFix attempts to retry fixing the integrity of data a limited number of times
func (automation *DecentralizedStorageIntegrityAutomation) retryIntegrityFix(data common.StorageData) {
    automation.integrityRetryMap[data.ID]++
    if automation.integrityRetryMap[data.ID] < MaxIntegrityRetries {
        automation.fixIntegrityIssue(data)
    } else {
        fmt.Printf("Max retries reached for resolving integrity issue with data ID %s. Marking as failed.\n", data.ID)
        automation.logIntegrityFailure(data.ID)
    }
}

// resetIntegrityRetry resets the retry counter for a specific data ID
func (automation *DecentralizedStorageIntegrityAutomation) resetIntegrityRetry(dataID string) {
    automation.integrityRetryMap[dataID] = 0
}

// logIntegrityFixEvent logs a successful integrity fix event into the ledger
func (automation *DecentralizedStorageIntegrityAutomation) logIntegrityFixEvent(dataID string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("storage-integrity-fix-%s", dataID),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Integrity Fix",
        Status:    "Resolved",
        Details:   fmt.Sprintf("Integrity issue for data ID %s successfully resolved.", dataID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with integrity fix event for data ID %s.\n", dataID)
}

// logIntegrityFailure logs a failure event into the ledger if integrity cannot be restored
func (automation *DecentralizedStorageIntegrityAutomation) logIntegrityFailure(dataID string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("storage-integrity-failure-%s", dataID),
        Timestamp: time.Now().Unix(),
        Type:      "Storage Integrity Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to resolve integrity issue for data ID %s after maximum retries.", dataID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with integrity failure for data ID %s.\n", dataID)
}

// validateStorageIntegrity continuously validates the integrity of the decentralized storage
func (automation *DecentralizedStorageIntegrityAutomation) validateStorageIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateStorageIntegrity()
    if !integrityValid {
        fmt.Println("Storage integrity breach detected. Initiating corrective measures.")
        automation.checkAndEnsureStorageIntegrity()
    } else {
        fmt.Println("All decentralized storage data integrity is valid.")
    }
}
