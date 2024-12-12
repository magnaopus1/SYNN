package automations

import (
    "fmt"
    "os"
    "time"
    "sync"
    "compress/gzip"
    "bytes"
    "io/ioutil"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    CompressionCheckInterval = 5 * time.Minute  // Interval to check for compression
    MaxFileSizeForCompression = 1048576         // 1 MB, files larger than this are compressed
)

// FileCompressionAutomation handles file compression automation
type FileCompressionAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus
    ledgerInstance   *ledger.Ledger
    stateMutex       *sync.RWMutex
    compressionCycle int
}

// NewFileCompressionAutomation initializes the file compression automation
func NewFileCompressionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *FileCompressionAutomation {
    return &FileCompressionAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        compressionCycle: 0,
    }
}

// StartCompressionMonitoring starts the file compression monitoring in a continuous loop
func (automation *FileCompressionAutomation) StartCompressionMonitoring() {
    ticker := time.NewTicker(CompressionCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkAndCompressFiles()
        }
    }()
}

// checkAndCompressFiles checks for files that need to be compressed
func (automation *FileCompressionAutomation) checkAndCompressFiles() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    files := automation.fetchFilesForCompression()

    for _, file := range files {
        if file.Size >= MaxFileSizeForCompression {
            fmt.Printf("File %s exceeds size limit (%d bytes). Compressing...\n", file.Name, file.Size)
            err := automation.compressFile(file)
            if err != nil {
                fmt.Printf("Failed to compress file %s: %v\n", file.Name, err)
                automation.logCompressionFailure(file, err)
            } else {
                automation.logCompressionSuccess(file)
            }
        }
    }

    automation.compressionCycle++
    fmt.Printf("Compression cycle #%d completed.\n", automation.compressionCycle)

    if automation.compressionCycle % 1000 == 0 {
        automation.finalizeCompressionCycle()
    }
}

// fetchFilesForCompression fetches a list of files that need to be checked for compression
func (automation *FileCompressionAutomation) fetchFilesForCompression() []common.FileInfo {
    // Placeholder logic, replace with logic to fetch real file information
    return common.GetFileList()
}

// compressFile compresses the given file
func (automation *FileCompressionAutomation) compressFile(file common.FileInfo) error {
    // Read the file content
    data, err := ioutil.ReadFile(file.Path)
    if err != nil {
        return fmt.Errorf("failed to read file: %v", err)
    }

    // Compress the file data
    compressedData, err := automation.performCompression(data)
    if err != nil {
        return fmt.Errorf("compression failed: %v", err)
    }

    // Write the compressed data back to a new file
    compressedFilePath := fmt.Sprintf("%s.gz", file.Path)
    err = ioutil.WriteFile(compressedFilePath, compressedData, os.ModePerm)
    if err != nil {
        return fmt.Errorf("failed to write compressed file: %v", err)
    }

    // Encrypt the compressed file for storage
    err = automation.encryptAndStoreFile(compressedFilePath)
    if err != nil {
        return fmt.Errorf("failed to encrypt and store file: %v", err)
    }

    return nil
}

// performCompression performs the compression logic using gzip
func (automation *FileCompressionAutomation) performCompression(data []byte) ([]byte, error) {
    var buf bytes.Buffer
    writer := gzip.NewWriter(&buf)
    _, err := writer.Write(data)
    if err != nil {
        return nil, err
    }
    err = writer.Close()
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

// encryptAndStoreFile encrypts the compressed file and stores it securely
func (automation *FileCompressionAutomation) encryptAndStoreFile(filePath string) error {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read compressed file for encryption: %v", err)
    }

    // Encrypt the file data before storage
    encryptedData, err := encryption.EncryptData(data)
    if err != nil {
        return fmt.Errorf("failed to encrypt file: %v", err)
    }

    // Store the encrypted file
    err = ioutil.WriteFile(filePath, encryptedData, os.ModePerm)
    if err != nil {
        return fmt.Errorf("failed to store encrypted file: %v", err)
    }

    return nil
}

// finalizeCompressionCycle logs the finalization of a compression cycle
func (automation *FileCompressionAutomation) finalizeCompressionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeCompressionCycle()
    if success {
        fmt.Println("Compression cycle finalized successfully.")
        automation.logCycleFinalization()
    } else {
        fmt.Println("Error finalizing compression cycle.")
    }
}

// logCompressionSuccess logs successful file compression in the ledger
func (automation *FileCompressionAutomation) logCompressionSuccess(file common.FileInfo) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-%s", file.Name),
        Timestamp: time.Now().Unix(),
        Type:      "File Compression",
        Status:    "Success",
        Details:   fmt.Sprintf("File %s compressed successfully.", file.Name),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with file compression success for %s.\n", file.Name)
}

// logCompressionFailure logs failed file compression attempts in the ledger
func (automation *FileCompressionAutomation) logCompressionFailure(file common.FileInfo, err error) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-failure-%s", file.Name),
        Timestamp: time.Now().Unix(),
        Type:      "File Compression",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to compress file %s: %v", file.Name, err),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with file compression failure for %s.\n", file.Name)
}

// logCycleFinalization logs the finalization of the compression cycle in the ledger
func (automation *FileCompressionAutomation) logCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Compression Cycle",
        Status:    "Finalized",
        Details:   "Compression cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with compression cycle finalization.")
}
