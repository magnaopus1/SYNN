package automations

import (
    "fmt"
    "time"
    "sync"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    DecompressionCheckInterval = 2 * time.Minute // Interval for checking compressed data
    DecompressionThreshold     = 85              // Decompression trigger when data is compressed over this percentage
    MaxDecompressionRetries    = 3               // Maximum number of retry attempts for decompression
)

// DataDecompressionAutomation handles automated data decompression in the network
type DataDecompressionAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Integration with Synnergy Consensus
    ledgerInstance    *ledger.Ledger               // Integration with the ledger
    stateMutex        *sync.RWMutex                // Mutex for synchronization
    retryCounts       map[string]int               // Track retry attempts per node
    decompressionCycle int                         // Counter for decompression cycles
}

// NewDataDecompressionAutomation initializes the automation
func NewDataDecompressionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataDecompressionAutomation {
    return &DataDecompressionAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        retryCounts:        make(map[string]int),
        decompressionCycle: 0,
    }
}

// StartDecompressionMonitoring begins the monitoring process in a loop
func (automation *DataDecompressionAutomation) StartDecompressionMonitoring() {
    ticker := time.NewTicker(DecompressionCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkAndTriggerDecompression()
        }
    }()
}

// checkAndTriggerDecompression monitors compression and triggers decompression if needed
func (automation *DataDecompressionAutomation) checkAndTriggerDecompression() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch compressed data reports from Synnergy Consensus
    compressionReports := automation.consensusSystem.FetchCompressionReports()

    for _, report := range compressionReports {
        if report.CompressionRate >= DecompressionThreshold {
            fmt.Printf("Data decompression needed for node %s (Compression: %d%%). Initiating decompression.\n", report.NodeID, report.CompressionRate)
            automation.decompressData(report)
        } else {
            fmt.Printf("Compression on node %s is below threshold: %d%%.\n", report.NodeID, report.CompressionRate)
        }
    }

    automation.decompressionCycle++
    fmt.Printf("Decompression cycle #%d completed.\n", automation.decompressionCycle)

    if automation.decompressionCycle%SubBlocksPerBlock == 0 {
        automation.finalizeDecompressionCycle()
    }
}

// decompressData handles the data decompression process for a node
func (automation *DataDecompressionAutomation) decompressData(report common.CompressionReport) {
    // Decrypt the data before decompression
    decryptedData := automation.decryptData(report.CompressedData)

    // Perform decompression process
    decompressedData, err := automation.performDecompression(decryptedData)
    if err != nil {
        fmt.Printf("Decompression failed for node %s: %v. Retrying...\n", report.NodeID, err)
        automation.retryDecompression(report)
        return
    }

    success := automation.consensusSystem.InitiateDataDecompression(report.NodeID, decompressedData)
    if success {
        fmt.Printf("Data decompression successful for node %s.\n", report.NodeID)
        automation.logDecompressionEvent(report, "Decompressed")
        automation.resetRetryCount(report.NodeID)
    } else {
        fmt.Printf("Data decompression failed for node %s. Retrying...\n", report.NodeID)
        automation.retryDecompression(report)
    }
}

// retryDecompression handles retries for failed decompression attempts
func (automation *DataDecompressionAutomation) retryDecompression(report common.CompressionReport) {
    automation.retryCounts[report.NodeID]++
    if automation.retryCounts[report.NodeID] < MaxDecompressionRetries {
        automation.decompressData(report)
    } else {
        fmt.Printf("Max retry attempts reached for node %s. Decompression failed.\n", report.NodeID)
        automation.logDecompressionFailure(report)
    }
}

// resetRetryCount resets the retry count for a node after a successful decompression
func (automation *DataDecompressionAutomation) resetRetryCount(nodeID string) {
    automation.retryCounts[nodeID] = 0
}

// finalizeDecompressionCycle logs the finalization of a decompression cycle
func (automation *DataDecompressionAutomation) finalizeDecompressionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeDecompressionCycle()
    if success {
        fmt.Println("Decompression cycle finalized successfully.")
        automation.logCycleFinalization()
    } else {
        fmt.Println("Error finalizing decompression cycle.")
    }
}

// logDecompressionEvent logs successful decompression events in the ledger
func (automation *DataDecompressionAutomation) logDecompressionEvent(report common.CompressionReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("decompression-%s-%s", report.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Decompression Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s data %s successfully.", report.NodeID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with decompression event for node %s.\n", report.NodeID)
}

// logDecompressionFailure logs failed decompression attempts in the ledger
func (automation *DataDecompressionAutomation) logDecompressionFailure(report common.CompressionReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("decompression-failure-%s", report.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Decompression Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Decompression failed for node %s after maximum retries.", report.NodeID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with decompression failure for node %s.\n", report.NodeID)
}

// logCycleFinalization logs the finalization of the decompression cycle in the ledger
func (automation *DataDecompressionAutomation) logCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("decompression-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Decompression Cycle",
        Status:    "Finalized",
        Details:   "Decompression cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with decompression cycle finalization.")
}

// performDecompression performs actual data decompression logic
func (automation *DataDecompressionAutomation) performDecompression(data []byte) ([]byte, error) {
    // Implement your decompression logic here, such as gzip, zlib, etc.
    decompressedData := common.DecompressData(data)
    if decompressedData == nil {
        return nil, fmt.Errorf("decompression error")
    }
    return decompressedData, nil
}

// decryptData decrypts the compressed data before decompression
func (automation *DataDecompressionAutomation) decryptData(data []byte) []byte {
    decryptedData, err := encryption.DecryptData(data)
    if err != nil {
        fmt.Printf("Error decrypting data: %v\n", err)
        return data
    }
    return decryptedData
}
