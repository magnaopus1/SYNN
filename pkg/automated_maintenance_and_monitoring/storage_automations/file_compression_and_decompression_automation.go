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
    CompressionMonitoringInterval  = 30 * time.Second // Interval for monitoring compression and decompression
    MaxCompressionRetries          = 3                // Maximum retries for compression/decompression failures
    SubBlocksPerBlock              = 1000             // Number of sub-blocks in a block
    CompressionThreshold           = 80               // Threshold percentage for compression trigger
)

// FileCompressionAndDecompressionAutomation manages file compression and decompression in real time
type FileCompressionAndDecompressionAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger for logging compression-related events
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    compressionRetryCount map[string]int            // Counter for retrying compression/decompression actions
    compressionCycleCount int                       // Counter for compression monitoring cycles
}

// NewFileCompressionAndDecompressionAutomation initializes the automation
func NewFileCompressionAndDecompressionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *FileCompressionAndDecompressionAutomation {
    return &FileCompressionAndDecompressionAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        compressionRetryCount: make(map[string]int),
        compressionCycleCount: 0,
    }
}

// StartCompressionMonitoring starts the continuous loop for monitoring file compression and decompression
func (automation *FileCompressionAndDecompressionAutomation) StartCompressionMonitoring() {
    ticker := time.NewTicker(CompressionMonitoringInterval)

    go func() {
        for range ticker.C {
            automation.monitorFileCompression()
        }
    }()
}

// monitorFileCompression checks for files that require compression or decompression and handles the process
func (automation *FileCompressionAndDecompressionAutomation) monitorFileCompression() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch file compression reports
    fileReports := automation.consensusSystem.FetchFileCompressionReports()

    for _, report := range fileReports {
        if automation.isCompressionRequired(report) {
            fmt.Printf("File compression required for file %s. Initiating compression.\n", report.FileID)
            automation.applyFileCompression(report)
        } else {
            fmt.Printf("No compression required for file %s.\n", report.FileID)
        }
    }

    automation.compressionCycleCount++
    fmt.Printf("File compression cycle #%d completed.\n", automation.compressionCycleCount)

    if automation.compressionCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeCompressionCycle()
    }
}

// isCompressionRequired checks if a file requires compression based on size or usage
func (automation *FileCompressionAndDecompressionAutomation) isCompressionRequired(report common.FileCompressionReport) bool {
    if report.UsagePercentage >= CompressionThreshold {
        fmt.Printf("File %s size at %d%%, triggering compression.\n", report.FileID, report.UsagePercentage)
        return true
    }
    return false
}

// applyFileCompression attempts to compress a file that exceeds the threshold
func (automation *FileCompressionAndDecompressionAutomation) applyFileCompression(report common.FileCompressionReport) {
    encryptedCompressionData := automation.encryptCompressionData(report)

    // Simulating raw compression action here
    compressedData := compressData(report.FileData)

    // Log the compression action
    fmt.Printf("File %s compressed successfully. File size reduced.\n", report.FileID)
    automation.logCompressionEvent(report, "File Compressed", compressedData)

    // Reset retry counter after successful compression
    automation.resetCompressionRetry(report.FileID)
}

// compressData handles the raw compression process (simulated)
func compressData(data []byte) []byte {
    // Simulate compression logic: placeholder for actual compression algorithm
    compressedData := make([]byte, len(data)/2) // Example: file size reduced by 50%
    copy(compressedData, data[:len(compressedData)])
    return compressedData
}

// retryFileCompression retries the compression process in case of failure
func (automation *FileCompressionAndDecompressionAutomation) retryFileCompression(report common.FileCompressionReport) {
    automation.compressionRetryCount[report.FileID]++
    if automation.compressionRetryCount[report.FileID] < MaxCompressionRetries {
        automation.applyFileCompression(report)
    } else {
        fmt.Printf("Max retries reached for compressing file %s. Compression failed.\n", report.FileID)
        automation.logCompressionFailure(report)
    }
}

// resetCompressionRetry resets the retry count for file compression actions
func (automation *FileCompressionAndDecompressionAutomation) resetCompressionRetry(fileID string) {
    automation.compressionRetryCount[fileID] = 0
}

// finalizeCompressionCycle finalizes the compression cycle and logs the result in the ledger
func (automation *FileCompressionAndDecompressionAutomation) finalizeCompressionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeCompressionCycle()
    if success {
        fmt.Println("Compression cycle finalized successfully.")
        automation.logCompressionCycleFinalization()
    } else {
        fmt.Println("Error finalizing compression cycle.")
    }
}

// logCompressionEvent logs a file compression event into the ledger
func (automation *FileCompressionAndDecompressionAutomation) logCompressionEvent(report common.FileCompressionReport, eventType string, compressedData []byte) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-event-%s-%s", report.FileID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "File Compression Event",
        Status:    eventType,
        Details:   fmt.Sprintf("File %s compressed successfully. Compressed size: %d bytes.", report.FileID, len(compressedData)),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compression event for file %s.\n", report.FileID)
}

// logCompressionFailure logs the failure of a file compression attempt into the ledger
func (automation *FileCompressionAndDecompressionAutomation) logCompressionFailure(report common.FileCompressionReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-failure-%s", report.FileID),
        Timestamp: time.Now().Unix(),
        Type:      "File Compression Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("File compression failed for file %s after maximum retries.", report.FileID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compression failure for file %s.\n", report.FileID)
}

// logCompressionCycleFinalization logs the finalization of a compression cycle into the ledger
func (automation *FileCompressionAndDecompressionAutomation) logCompressionCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Compression Cycle Finalization",
        Status:    "Finalized",
        Details:   "File compression cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with compression cycle finalization.")
}

// encryptCompressionData encrypts compression-related data before taking action or logging events
func (automation *FileCompressionAndDecompressionAutomation) encryptCompressionData(report common.FileCompressionReport) common.FileCompressionReport {
    encryptedData, err := encryption.EncryptData(report.FileData)
    if err != nil {
        fmt.Println("Error encrypting file compression data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("File compression data successfully encrypted for file:", report.FileID)
    return report
}
