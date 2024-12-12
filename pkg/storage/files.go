package storage

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "time"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/common"
)

// NewFileManager initializes a new FileManager
func NewFileManager(storageDir string, ledgerInstance *ledger.Ledger) *FileManager {
    if _, err := os.Stat(storageDir); os.IsNotExist(err) {
        os.MkdirAll(storageDir, os.ModePerm)
    }

    return &FileManager{
        Files:          make(map[string]*FileEntry),
        storageDir:     storageDir,
        LedgerInstance: ledgerInstance,
    }
}

// StoreFile stores a file securely, encrypts it, and logs the operation in the ledger
func (fm *FileManager) StoreFile(fileName string, fileData []byte) (string, error) {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    encryptedData, err := encryption.EncryptData(string(fileData), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt file data: %v", err)
    }

    fileHash := fm.generateFileHash(fileData)
    filePath := filepath.Join(fm.storageDir, fileHash)

    err = ioutil.WriteFile(filePath, []byte(encryptedData), 0644)
    if err != nil {
        return "", fmt.Errorf("failed to write file to disk: %v", err)
    }

    fm.Files[fileHash] = &FileEntry{
        FileName:     fileName,
        FilePath:     filePath,
        Encrypted:    true,
        UploadedAt:   time.Now(),
        LastAccessed: time.Now(),
    }

    fmt.Printf("File %s stored successfully with hash: %s\n", fileName, fileHash)

    // Log the operation in the ledger
    err = fm.LedgerInstance.RecordFileOperation(fileHash, "store", fileName, encryptedData)
    if err != nil {
        return "", fmt.Errorf("failed to log file operation to ledger: %v", err)
    }

    return fileHash, nil
}

// RetrieveFile retrieves and decrypts a file by its hash
func (fm *FileManager) RetrieveFile(fileHash string) ([]byte, error) {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    entry, exists := fm.Files[fileHash]
    if !exists {
        return nil, errors.New("file not found")
    }

    fileData, err := ioutil.ReadFile(entry.FilePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read file from disk: %v", err)
    }

    decryptedData, err := encryption.DecryptData(string(fileData), common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt file data: %v", err)
    }

    entry.LastAccessed = time.Now()

    fmt.Printf("File with hash %s retrieved successfully.\n", fileHash)

    // Log the retrieval in the ledger
    err = fm.LedgerInstance.RecordFileOperation(fileHash, "retrieve", entry.FileName, string(fileData))
    if err != nil {
        return nil, fmt.Errorf("failed to log file retrieval to ledger: %v", err)
    }

    return []byte(decryptedData), nil
}

// DeleteFile deletes a file by its hash and logs the operation in the ledger
func (fm *FileManager) DeleteFile(fileHash string) error {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    entry, exists := fm.Files[fileHash]
    if !exists {
        return errors.New("file not found")
    }

    err := os.Remove(entry.FilePath)
    if err != nil {
        return fmt.Errorf("failed to delete file: %v", err)
    }

    delete(fm.Files, fileHash)

    fmt.Printf("File with hash %s deleted successfully.\n", fileHash)

    // Log the deletion in the ledger
    err = fm.LedgerInstance.RecordFileOperation(fileHash, "delete", entry.FileName, "")
    if err != nil {
        return fmt.Errorf("failed to log file deletion to ledger: %v", err)
    }

    return nil
}

// ListFiles lists all the files currently stored in the system
func (fm *FileManager) ListFiles() []FileEntry {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    files := make([]FileEntry, 0, len(fm.Files))
    for _, entry := range fm.Files {
        files = append(files, *entry)
    }

    return files
}

// generateFileHash generates a hash for a file's content
func (fm *FileManager) generateFileHash(fileData []byte) string {
    hash := sha256.New()
    hash.Write(fileData)
    return hex.EncodeToString(hash.Sum(nil))
}

// CleanUpExpiredFiles removes files older than the given duration
func (fm *FileManager) CleanUpExpiredFiles(duration time.Duration) {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    now := time.Now()

    for hash, entry := range fm.Files {
        if now.Sub(entry.UploadedAt) > duration {
            os.Remove(entry.FilePath)
            delete(fm.Files, hash)
            fmt.Printf("File %s expired and removed.\n", hash)
            fm.LedgerInstance.RecordFileOperation(hash, "cleanup", entry.FileName, "")
        }
    }
}
