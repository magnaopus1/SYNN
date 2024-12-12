package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/scalability/compression"
)

const (
    OnChainStorageOptimizationInterval  = 6 * time.Hour  // Interval for checking on-chain data storage optimization
    MaxOptimizationRetries              = 3              // Maximum retries for failed optimization operations
    SubBlocksPerBlock                   = 1000           // Number of sub-blocks in a block
    CompressionThreshold                = 50 * 1024 * 1024 // 50 MB threshold for on-chain data compression
)

// OnChainDataStorageOptimizationAutomation handles on-chain data optimization for efficient storage
type OnChainDataStorageOptimizationAutomation struct {
    consensusSystem         *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance          *ledger.Ledger               // Ledger for logging storage optimization events
    stateMutex              *sync.RWMutex                // Mutex for thread-safe access
    optimizationRetryCount  map[string]int               // Counter for retrying failed optimization operations
    optimizationCycleCount  int                          // Counter for optimization cycles
}

// NewOnChainDataStorageOptimizationAutomation initializes the automation for on-chain data storage optimization
func NewOnChainDataStorageOptimizationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *OnChainDataStorageOptimizationAutomation {
    return &OnChainDataStorageOptimizationAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        optimizationRetryCount: make(map[string]int),
        optimizationCycleCount: 0,
    }
}

// StartOnChainStorageOptimization starts the continuous loop for on-chain storage optimization
func (automation *OnChainDataStorageOptimizationAutomation) StartOnChainStorageOptimization() {
    ticker := time.NewTicker(OnChainStorageOptimizationInterval)

    go func() {
        for range ticker.C {
            automation.monitorOnChainStorageOptimization()
        }
    }()
}

// monitorOnChainStorageOptimization checks for data that needs compression or restructuring to optimize on-chain storage
func (automation *OnChainDataStorageOptimizationAutomation) monitorOnChainStorageOptimization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch reports for data that potentially needs optimization
    optimizationReports := automation.consensusSystem.FetchOnChainOptimizationReports()

    for _, report := range optimizationReports {
        if automation.isCompressionRequired(report) {
            fmt.Printf("Compression required for data %s. Initiating compression process.\n", report.DataID)
            automation.applyCompression(report)
        } else {
            fmt.Printf("No optimization required for data %s.\n", report.DataID)
        }
    }

    automation.optimizationCycleCount++
    fmt.Printf("On-chain storage optimization cycle #%d completed.\n", automation.optimizationCycleCount)

    if automation.optimizationCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeOptimizationCycle()
    }
}

// isCompressionRequired checks if the data exceeds the threshold for on-chain data compression
func (automation *OnChainDataStorageOptimizationAutomation) isCompressionRequired(report common.StorageReport) bool {
    if report.DataSize >= CompressionThreshold {
        fmt.Printf("Data %s exceeds compression threshold (size: %d bytes), triggering compression.\n", report.DataID, report.DataSize)
        return true
    }
    return false
}

// applyCompression compresses the data and updates on-chain references
func (automation *OnChainDataStorageOptimizationAutomation) applyCompression(report common.StorageReport) {
    compressedData, err := compression.CompressData(report.Data)
    if err != nil {
        fmt.Printf("Error compressing data %s: %v. Retrying...\n", report.DataID, err)
        automation.retryOptimization(report)
        return
    }

    fmt.Printf("Data %s successfully compressed.\n", report.DataID)
    automation.logOptimizationEvent(report, "Compressed", len(compressedData))
    automation.resetOptimizationRetry(report.DataID)
}

// retryOptimization retries the optimization process in case of failure
func (automation *OnChainDataStorageOptimizationAutomation) retryOptimization(report common.StorageReport) {
    automation.optimizationRetryCount[report.DataID]++
    if automation.optimizationRetryCount[report.DataID] < MaxOptimizationRetries {
        automation.applyCompression(report)
    } else {
        fmt.Printf("Max retries reached for compressing data %s. Optimization failed.\n", report.DataID)
        automation.logOptimizationFailure(report)
    }
}

// resetOptimizationRetry resets the retry count for optimization operations
func (automation *OnChainDataStorageOptimizationAutomation) resetOptimizationRetry(dataID string) {
    automation.optimizationRetryCount[dataID] = 0
}

// finalizeOptimizationCycle finalizes the optimization cycle and logs the result in the ledger
func (automation *OnChainDataStorageOptimizationAutomation) finalizeOptimizationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeOptimizationCycle()
    if success {
        fmt.Println("On-chain storage optimization cycle finalized successfully.")
        automation.logOptimizationCycleFinalization()
    } else {
        fmt.Println("Error finalizing on-chain storage optimization cycle.")
    }
}

// logOptimizationEvent logs a successful optimization event into the ledger
func (automation *OnChainDataStorageOptimizationAutomation) logOptimizationEvent(report common.StorageReport, eventType string, compressedSize int) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("optimization-event-%s-%s", report.DataID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "On-Chain Storage Optimization Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Data %s compressed from %d bytes to %d bytes.", report.DataID, report.DataSize, compressedSize),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with optimization event for data %s.\n", report.DataID)
}

// logOptimizationFailure logs the failure of an optimization attempt into the ledger
func (automation *OnChainDataStorageOptimizationAutomation) logOptimizationFailure(report common.StorageReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("optimization-failure-%s", report.DataID),
        Timestamp: time.Now().Unix(),
        Type:      "On-Chain Storage Optimization Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Optimization failed for data %s after maximum retries.", report.DataID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with optimization failure for data %s.\n", report.DataID)
}

// logOptimizationCycleFinalization logs the finalization of an optimization cycle into the ledger
func (automation *OnChainDataStorageOptimizationAutomation) logOptimizationCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("optimization-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Optimization Cycle Finalization",
        Status:    "Finalized",
        Details:   "On-chain storage optimization cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with optimization cycle finalization.")
}
