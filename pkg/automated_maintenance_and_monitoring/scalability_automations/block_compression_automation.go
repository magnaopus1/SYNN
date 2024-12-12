package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "bytes"
    "compress/gzip"
    "io/ioutil"
)

const (
    CompressionCheckInterval   = 15 * time.Minute // Interval for checking compression requirements
    MaxCompressionRetries      = 3                // Maximum retry attempts for failed compression
    SubBlocksPerBlock          = 1000             // Number of sub-blocks in a block
    CompressionThreshold       = 85               // Compression threshold for blocks (e.g., compress if block size exceeds this percentage)
)

// BlockCompressionAutomation manages the compression of blocks to optimize storage
type BlockCompressionAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger for logging compression events
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    compressionRetryMap map[string]int              // Counter for retrying failed compression
    compressionCycle   int                          // Counter for compression cycles
}

// NewBlockCompressionAutomation initializes the automation for block compression
func NewBlockCompressionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *BlockCompressionAutomation {
    return &BlockCompressionAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        compressionRetryMap: make(map[string]int),
        compressionCycle:    0,
    }
}

// StartBlockCompression starts the continuous loop for compressing blocks
func (automation *BlockCompressionAutomation) StartBlockCompression() {
    ticker := time.NewTicker(CompressionCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkAndCompressBlocks()
        }
    }()
}

// checkAndCompressBlocks checks if blocks exceed the size threshold and compresses them if needed
func (automation *BlockCompressionAutomation) checkAndCompressBlocks() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    blockReports := automation.consensusSystem.FetchBlockReports()

    for _, report := range blockReports {
        if report.SizePercent >= CompressionThreshold {
            fmt.Printf("Block %s exceeds compression threshold with size %d%%. Initiating compression.\n", report.BlockID, report.SizePercent)
            automation.compressBlock(report)
        } else {
            fmt.Printf("Block %s is within size limits: %d%%.\n", report.BlockID, report.SizePercent)
        }
    }

    automation.compressionCycle++
    fmt.Printf("Block compression cycle #%d completed.\n", automation.compressionCycle)

    if automation.compressionCycle%SubBlocksPerBlock == 0 {
        automation.finalizeCompressionCycle()
    }
}

// compressBlock handles the compression of a block
func (automation *BlockCompressionAutomation) compressBlock(report common.BlockReport) {
    encryptedData := automation.encryptBlockData(report)

    compressedData, err := automation.rawCompressBlock(encryptedData)
    if err != nil {
        fmt.Printf("Compression failed for block %s: %v. Retrying...\n", report.BlockID, err)
        automation.retryCompression(report)
        return
    }

    fmt.Printf("Block %s successfully compressed.\n", report.BlockID)
    automation.logCompressionEvent(report, "Compressed")
    automation.resetCompressionRetry(report.BlockID)
}

// rawCompressBlock handles the actual compression logic for block data
func (automation *BlockCompressionAutomation) rawCompressBlock(data []byte) ([]byte, error) {
    var buffer bytes.Buffer
    writer := gzip.NewWriter(&buffer)
    _, err := writer.Write(data)
    if err != nil {
        return nil, err
    }

    err = writer.Close()
    if err != nil {
        return nil, err
    }

    fmt.Println("Data successfully compressed.")
    return buffer.Bytes(), nil
}

// retryCompression retries the compression process if it fails
func (automation *BlockCompressionAutomation) retryCompression(report common.BlockReport) {
    automation.compressionRetryMap[report.BlockID]++
    if automation.compressionRetryMap[report.BlockID] < MaxCompressionRetries {
        automation.compressBlock(report)
    } else {
        fmt.Printf("Max retries reached for block %s. Compression failed.\n", report.BlockID)
        automation.logCompressionFailure(report)
    }
}

// resetCompressionRetry resets the retry count for block compression operations
func (automation *BlockCompressionAutomation) resetCompressionRetry(blockID string) {
    automation.compressionRetryMap[blockID] = 0
}

// encryptBlockData encrypts the block data before compression
func (automation *BlockCompressionAutomation) encryptBlockData(report common.BlockReport) []byte {
    encryptedData, err := encryption.EncryptData(report.Data)
    if err != nil {
        fmt.Printf("Error encrypting data for block %s: %v\n", report.BlockID, err)
        return report.Data
    }

    fmt.Printf("Block data successfully encrypted for block %s.\n", report.BlockID)
    return encryptedData
}

// logCompressionEvent logs a successful compression event into the ledger
func (automation *BlockCompressionAutomation) logCompressionEvent(report common.BlockReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-event-%s-%s", report.BlockID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Block Compression Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Block %s %s successfully.", report.BlockID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compression event for block %s.\n", report.BlockID)
}

// logCompressionFailure logs the failure to compress a specific block into the ledger
func (automation *BlockCompressionAutomation) logCompressionFailure(report common.BlockReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-failure-%s", report.BlockID),
        Timestamp: time.Now().Unix(),
        Type:      "Block Compression Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to compress block %s after maximum retries.", report.BlockID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compression failure for block %s.\n", report.BlockID)
}

// finalizeCompressionCycle finalizes the block compression cycle and logs the result in the ledger
func (automation *BlockCompressionAutomation) finalizeCompressionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeCompressionCycle()
    if success {
        fmt.Println("Block compression cycle finalized successfully.")
        automation.logCompressionCycleFinalization()
    } else {
        fmt.Println("Error finalizing block compression cycle.")
    }
}

// logCompressionCycleFinalization logs the finalization of a compression cycle into the ledger
func (automation *BlockCompressionAutomation) logCompressionCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Compression Cycle Finalization",
        Status:    "Finalized",
        Details:   "Block compression cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with compression cycle finalization.")
}
