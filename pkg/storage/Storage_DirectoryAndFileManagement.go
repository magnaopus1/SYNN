package storage

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
)

// StorageCreateDirectory creates a new directory with encryption
func StorageCreateDirectory(path string) error {
    if err := common.CreateDirectory(path); err != nil {
        return fmt.Errorf("directory creation failed: %v", err)
    }
    ledger.RecordDirectoryCreation(path, time.Now())
    fmt.Printf("Directory created at %s\n", path)
    return nil
}

// StorageDeleteDirectory deletes a directory and its contents
func StorageDeleteDirectory(path string) error {
    if err := common.DeleteDirectory(path); err != nil {
        return fmt.Errorf("directory deletion failed: %v", err)
    }
    ledger.RecordDirectoryDeletion(path, time.Now())
    fmt.Printf("Directory deleted at %s\n", path)
    return nil
}

// StorageListDirectoryContents lists the contents of a directory
func StorageListDirectoryContents(path string) ([]string, error) {
    contents, err := common.ListDirectoryContents(path)
    if err != nil {
        return nil, fmt.Errorf("failed to list directory contents: %v", err)
    }
    fmt.Printf("Contents of directory %s: %v\n", path, contents)
    return contents, nil
}

// StorageSetDirectoryPermissions sets permissions for a directory
func StorageSetDirectoryPermissions(path string, permissions string) error {
    if err := common.SetDirectoryPermissions(path, permissions); err != nil {
        return fmt.Errorf("setting directory permissions failed: %v", err)
    }
    ledger.RecordDirectoryPermissionsSet(path, permissions, time.Now())
    fmt.Printf("Permissions set for directory %s to %s\n", path, permissions)
    return nil
}

// StorageGetDirectoryPermissions retrieves the permissions for a directory
func StorageGetDirectoryPermissions(path string) (string, error) {
    permissions, err := common.GetDirectoryPermissions(path)
    if err != nil {
        return "", fmt.Errorf("failed to get directory permissions: %v", err)
    }
    fmt.Printf("Permissions for directory %s: %s\n", path, permissions)
    return permissions, nil
}

// StorageMonitorDirectoryUsage monitors usage for a specific directory
func StorageMonitorDirectoryUsage(path string) (float64, error) {
    usage, err := common.MonitorDirectoryUsage(path)
    if err != nil {
        return 0, fmt.Errorf("directory usage monitoring failed: %v", err)
    }
    fmt.Printf("Directory %s usage: %.2f%%\n", path, usage)
    return usage, nil
}

// StorageAuditDirectory audits the directory for security and integrity
func StorageAuditDirectory(path string) error {
    if !common.VerifyDirectoryIntegrity(path) {
        return fmt.Errorf("directory audit failed for %s", path)
    }
    fmt.Printf("Directory %s passed integrity audit\n", path)
    return nil
}

// StorageTrackDirectoryChanges tracks changes within a directory
func StorageTrackDirectoryChanges(path string) error {
    if err := common.TrackDirectoryChanges(path); err != nil {
        return fmt.Errorf("tracking directory changes failed: %v", err)
    }
    ledger.RecordDirectoryChangeTracking(path, time.Now())
    fmt.Printf("Changes in directory %s are being tracked\n", path)
    return nil
}

// StorageFetchDirectoryChangeLog retrieves the change log for a directory
func StorageFetchDirectoryChangeLog(path string) ([]string, error) {
    log, err := ledger.FetchDirectoryChangeLog(path)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch directory change log: %v", err)
    }
    fmt.Println("Directory Change Log:\n", log)
    return log, nil
}

// StorageCreateFile creates a new file with encryption
func StorageCreateFile(filePath string, data []byte) error {
    encryptedData, err := encryption.EncryptData(data, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("encryption failed: %v", err)
    }
    if err := common.WriteFile(filePath, encryptedData); err != nil {
        return fmt.Errorf("file creation failed: %v", err)
    }
    ledger.RecordFileCreation(filePath, time.Now())
    fmt.Printf("File created at %s\n", filePath)
    return nil
}

// StorageDeleteFile deletes a specified file
func StorageDeleteFile(filePath string) error {
    if err := common.DeleteFile(filePath); err != nil {
        return fmt.Errorf("file deletion failed: %v", err)
    }
    ledger.RecordFileDeletion(filePath, time.Now())
    fmt.Printf("File deleted at %s\n", filePath)
    return nil
}

// StorageReadFile reads a file and decrypts its contents
func StorageReadFile(filePath string) ([]byte, error) {
    encryptedData, err := common.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("file read failed: %v", err)
    }
    data, err := encryption.DecryptData(encryptedData, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %v", err)
    }
    fmt.Printf("File read from %s\n", filePath)
    return data, nil
}

// StorageWriteFile writes encrypted data to a file
func StorageWriteFile(filePath string, data []byte) error {
    encryptedData, err := encryption.EncryptData(data, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("encryption failed: %v", err)
    }
    if err := common.WriteFile(filePath, encryptedData); err != nil {
        return fmt.Errorf("file write failed: %v", err)
    }
    ledger.RecordFileWrite(filePath, time.Now())
    fmt.Printf("File written to %s\n", filePath)
    return nil
}

// StorageTrackFileChanges tracks changes to a specific file
func StorageTrackFileChanges(filePath string) error {
    if err := common.TrackFileChanges(filePath); err != nil {
        return fmt.Errorf("tracking file changes failed: %v", err)
    }
    ledger.RecordFileChangeTracking(filePath, time.Now())
    fmt.Printf("Changes to file %s are being tracked\n", filePath)
    return nil
}

// StorageFetchFileChangeHistory retrieves the change history for a file
func StorageFetchFileChangeHistory(filePath string) ([]string, error) {
    history, err := ledger.FetchFileChangeHistory(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch file change history: %v", err)
    }
    fmt.Println("File Change History:\n", history)
    return history, nil
}

// StorageAuditFileAccess audits access logs for a specific file
func StorageAuditFileAccess(filePath string) error {
    if !common.VerifyFileAccessLogIntegrity(filePath) {
        return fmt.Errorf("file access audit failed for %s", filePath)
    }
    fmt.Printf("File access for %s passed audit\n", filePath)
    return nil
}

// StorageSetFileVersionLimit sets the maximum number of versions retained for a file
func StorageSetFileVersionLimit(filePath string, limit int) error {
    if err := common.SetFileVersionLimit(filePath, limit); err != nil {
        return fmt.Errorf("setting file version limit failed: %v", err)
    }
    ledger.RecordFileVersionLimitSet(filePath, limit, time.Now())
    fmt.Printf("Version limit for file %s set to %d\n", filePath, limit)
    return nil
}

// StorageGetFileVersionLimit retrieves the version limit for a file
func StorageGetFileVersionLimit(filePath string) (int, error) {
    limit, err := common.GetFileVersionLimit(filePath)
    if err != nil {
        return 0, fmt.Errorf("failed to get file version limit: %v", err)
    }
    fmt.Printf("Current version limit for file %s: %d\n", filePath, limit)
    return limit, nil
}

// StorageFetchFileVersions fetches the versions of a file
func StorageFetchFileVersions(filePath string) ([]string, error) {
    versions, err := ledger.FetchFileVersions(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch file versions: %v", err)
    }
    fmt.Println("File Versions:\n", versions)
    return versions, nil
}

// StorageRevertFileVersion reverts a file to a specific version
func StorageRevertFileVersion(filePath, versionID string) error {
    if err := common.RevertFileVersion(filePath, versionID); err != nil {
        return fmt.Errorf("reverting file version failed: %v", err)
    }
    ledger.RecordFileReversion(filePath, versionID, time.Now())
    fmt.Printf("File %s reverted to version %s\n", filePath, versionID)
    return nil
}

// StorageMonitorFileAccess monitors access logs for a specific file
func StorageMonitorFileAccess(filePath string) error {
    if err := common.MonitorFileAccess(filePath); err != nil {
        return fmt.Errorf("file access monitoring failed: %v", err)
    }
    fmt.Printf("Monitoring access for file %s\n", filePath)
    return nil
}

// StorageFetchFileAccessLog retrieves the access log for a specific file
func StorageFetchFileAccessLog(filePath string) ([]string, error) {
    log, err := ledger.FetchFileAccessLog(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch file access log: %v", err)
    }
    fmt.Println("File Access Log:\n", log)
    return log, nil
}

// StorageAuditFileAccessLog audits the file access log for integrity
func StorageAuditFileAccessLog(filePath string) error {
    if !common.VerifyFileAccessLogIntegrity(filePath) {
        return fmt.Errorf("file access log audit failed for %s", filePath)
    }
    fmt.Printf("File access log for %s passed audit\n", filePath)
    return nil
}

// StorageSetFileCompression enables file compression
func StorageSetFileCompression(filePath string) error {
    if err := common.EnableFileCompression(filePath); err != nil {
        return fmt.Errorf("file compression failed: %v", err)
    }
    ledger.RecordFileCompressionSet(filePath, time.Now())
    fmt.Printf("File compression enabled for %s\n", filePath)
    return nil
}

// StorageDisableFileCompression disables file compression
func StorageDisableFileCompression(filePath string) error {
    if err := common.DisableFileCompression(filePath); err != nil {
        return fmt.Errorf("disabling file compression failed: %v", err)
    }
    ledger.RecordFileCompressionDisabled(filePath, time.Now())
    fmt.Printf("File compression disabled for %s\n", filePath)
    return nil
}
