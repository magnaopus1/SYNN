package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    EncryptionCheckInterval    = 10 * time.Second  // Interval for checking data encryption status
    SubBlocksPerBlock          = 1000              // Number of sub-blocks in a block
    MaxEncryptionRetries       = 5                 // Maximum number of retry attempts for encryption
)

// DataEncryptionSecurityAutomation handles encryption for data before it's processed in the blockchain
type DataEncryptionSecurityAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger for logging encryption events
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    encryptionRetryCount map[string]int               // Counter for retrying encryption on failure
    encryptionCycleCount int                          // Counter for encryption cycles
}

// NewDataEncryptionSecurityAutomation initializes the automation for data encryption
func NewDataEncryptionSecurityAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataEncryptionSecurityAutomation {
    return &DataEncryptionSecurityAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        encryptionRetryCount: make(map[string]int),
        encryptionCycleCount: 0,
    }
}

// StartEncryptionMonitoring starts the continuous loop for monitoring and enforcing encryption
func (automation *DataEncryptionSecurityAutomation) StartEncryptionMonitoring() {
    ticker := time.NewTicker(EncryptionCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEncryptData()
        }
    }()
}

// monitorAndEncryptData checks all data and ensures it's encrypted before processing
func (automation *DataEncryptionSecurityAutomation) monitorAndEncryptData() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of all data that needs encryption
    dataList := automation.consensusSystem.GetDataForEncryption()

    if len(dataList) > 0 {
        for _, dataItem := range dataList {
            fmt.Printf("Encrypting data for item %s.\n", dataItem.ID)
            automation.encryptDataItem(dataItem)
        }
    } else {
        fmt.Println("No data needs encryption at this time.")
    }

    automation.encryptionCycleCount++
    fmt.Printf("Data encryption cycle #%d executed.\n", automation.encryptionCycleCount)

    if automation.encryptionCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeEncryptionCycle()
    }
}

// encryptDataItem handles the encryption of a specific data item
func (automation *DataEncryptionSecurityAutomation) encryptDataItem(dataItem common.DataItem) {
    // Encrypt the data before processing
    encryptedData := automation.encryptData(dataItem)

    // Trigger encrypted data processing through the Synnergy Consensus system
    encryptionSuccess := automation.consensusSystem.ProcessEncryptedData(dataItem, encryptedData)

    if encryptionSuccess {
        fmt.Printf("Data %s encrypted and processed successfully.\n", dataItem.ID)
        automation.logEncryptionEvent(dataItem)
        automation.resetEncryptionRetry(dataItem.ID)
    } else {
        fmt.Printf("Error encrypting data %s. Retrying...\n", dataItem.ID)
        automation.retryEncryption(dataItem)
    }
}

// retryEncryption attempts to retry a failed encryption a limited number of times
func (automation *DataEncryptionSecurityAutomation) retryEncryption(dataItem common.DataItem) {
    automation.encryptionRetryCount[dataItem.ID]++
    if automation.encryptionRetryCount[dataItem.ID] < MaxEncryptionRetries {
        automation.encryptDataItem(dataItem)
    } else {
        fmt.Printf("Max retries reached for data %s. Encryption failed.\n", dataItem.ID)
        automation.logEncryptionFailure(dataItem)
    }
}

// resetEncryptionRetry resets the retry count for a data encryption process
func (automation *DataEncryptionSecurityAutomation) resetEncryptionRetry(dataID string) {
    automation.encryptionRetryCount[dataID] = 0
}

// finalizeEncryptionCycle finalizes the encryption check cycle and logs the result in the ledger
func (automation *DataEncryptionSecurityAutomation) finalizeEncryptionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeEncryptionCycle()
    if success {
        fmt.Println("Data encryption check cycle finalized successfully.")
        automation.logEncryptionCycleFinalization()
    } else {
        fmt.Println("Error finalizing data encryption check cycle.")
    }
}

// logEncryptionEvent logs the encryption event into the ledger
func (automation *DataEncryptionSecurityAutomation) logEncryptionEvent(dataItem common.DataItem) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("data-encryption-%s", dataItem.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Data Encryption",
        Status:    "Completed",
        Details:   fmt.Sprintf("Encryption successfully completed for data %s.", dataItem.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with encryption event for data %s.\n", dataItem.ID)
}

// logEncryptionFailure logs the failure of a data encryption event into the ledger
func (automation *DataEncryptionSecurityAutomation) logEncryptionFailure(dataItem common.DataItem) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("data-encryption-failure-%s", dataItem.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Data Encryption Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Encryption failed for data %s after maximum retries.", dataItem.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with encryption failure for data %s.\n", dataItem.ID)
}

// encryptData handles the encryption of a data item
func (automation *DataEncryptionSecurityAutomation) encryptData(dataItem common.DataItem) common.EncryptedDataItem {
    encryptedData, err := encryption.EncryptData(dataItem.Data)
    if err != nil {
        fmt.Println("Error encrypting data:", err)
        return common.EncryptedDataItem{}
    }

    fmt.Println("Data successfully encrypted.")
    return common.EncryptedDataItem{
        ID:            dataItem.ID,
        EncryptedData: encryptedData,
    }
}

// logEncryptionCycleFinalization logs the finalization of an encryption cycle into the ledger
func (automation *DataEncryptionSecurityAutomation) logEncryptionCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("data-encryption-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Data Encryption Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with data encryption cycle finalization.")
}
