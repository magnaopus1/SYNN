package storage

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NewFileIndexer initializes a new FileIndexer
func NewFileIndexer(ledgerInstance *ledger.Ledger) *FileIndexer {
    return &FileIndexer{
        Indexes:        make(map[string]*FileIndex),
        LedgerInstance: ledgerInstance,
    }
}

// AddFileToIndex indexes a new file with its metadata, encrypting the IPFS CID and logging it in the ledger
func (indexer *FileIndexer) AddFileToIndex(fileName string, fileSize int64, owner string, cid string) (string, error) {
    indexer.mutex.Lock()
    defer indexer.mutex.Unlock()

    // Create a unique FileID based on the file and metadata
    fileID := indexer.generateFileID(fileName, owner, fileSize)
    
    encryptedCID, err := encryption.EncryptData(cid, common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt CID: %v", err)
    }

    fileIndex := &FileIndex{
        FileID:       fileID,
        FileName:     fileName,
        FileSize:     fileSize,
        UploadedAt:   time.Now(),
        Owner:        owner,
        EncryptedCID: encryptedCID,
    }

    // Add the file metadata to the index
    indexer.Indexes[fileID] = fileIndex

    // Log the indexed file to the ledger
    err = indexer.logFileIndexToLedger(fileIndex)
    if err != nil {
        return "", fmt.Errorf("failed to log file index to ledger: %v", err)
    }

    fmt.Printf("File %s added to index with FileID: %s\n", fileName, fileID)
    return fileID, nil
}

// RetrieveFileFromIndex fetches the file's metadata based on its FileID
func (indexer *FileIndexer) RetrieveFileFromIndex(fileID string) (*FileIndex, error) {
    indexer.mutex.Lock()
    defer indexer.mutex.Unlock()

    fileIndex, exists := indexer.Indexes[fileID]
    if !exists {
        return nil, fmt.Errorf("file with FileID %s not found", fileID)
    }

    fmt.Printf("File with FileID %s retrieved from index.\n", fileID)
    return fileIndex, nil
}

// RemoveFileFromIndex removes the file metadata from the index and logs the removal in the ledger
func (indexer *FileIndexer) RemoveFileFromIndex(fileID string) error {
    indexer.mutex.Lock()
    defer indexer.mutex.Unlock()

    _, exists := indexer.Indexes[fileID]
    if !exists {
        return fmt.Errorf("file with FileID %s not found", fileID)
    }

    // Remove the file metadata from the index
    delete(indexer.Indexes, fileID)

    // Log the removal operation in the ledger
    err := indexer.LedgerInstance.RecordFileOperation(fileID, "file_remove", "", "")
    if err != nil {
        return fmt.Errorf("failed to log file removal to ledger: %v", err)
    }

    fmt.Printf("File with FileID %s removed from index and logged in the ledger.\n", fileID)
    return nil
}

// generateFileID creates a unique file identifier based on the file's metadata
func (indexer *FileIndexer) generateFileID(fileName, owner string, fileSize int64) string {
    hashInput := fmt.Sprintf("%s%s%d%d", fileName, owner, fileSize, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// logFileIndexToLedger logs the file's metadata and encrypted CID to the ledger
func (indexer *FileIndexer) logFileIndexToLedger(fileIndex *FileIndex) error {
    fileMetadata := fmt.Sprintf("FileID: %s, FileName: %s, FileSize: %d, UploadedAt: %s, Owner: %s",
        fileIndex.FileID, fileIndex.FileName, fileIndex.FileSize, fileIndex.UploadedAt, fileIndex.Owner)

    encryptedMetadata, err := encryption.EncryptData(fileMetadata, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt file metadata: %v", err)
    }

    err = indexer.LedgerInstance.RecordFileOperation(fileIndex.FileID, "file_add", fileIndex.FileName, encryptedMetadata)
    if err != nil {
        return fmt.Errorf("failed to log file index to ledger: %v", err)
    }

    fmt.Printf("File metadata for FileID %s logged to the ledger.\n", fileIndex.FileID)
    return nil
}
