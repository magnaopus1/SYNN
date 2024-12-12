package storage

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
)

// StorageEncryptData encrypts data using the current encryption key
func StorageEncryptData(data []byte) ([]byte, error) {
    encryptedData, err := encryption.EncryptData(data, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("encryption failed: %v", err)
    }
    fmt.Println("Data encrypted successfully")
    return encryptedData, nil
}

// StorageDecryptData decrypts data using the current encryption key
func StorageDecryptData(data []byte) ([]byte, error) {
    decryptedData, err := encryption.DecryptData(data, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %v", err)
    }
    fmt.Println("Data decrypted successfully")
    return decryptedData, nil
}

// StorageSetEncryptionKey sets a new encryption key
func StorageSetEncryptionKey(newKey string) {
    common.EncryptionKey = newKey
    ledger.RecordEncryptionKeySet(newKey, time.Now())
    fmt.Println("Encryption key set successfully")
}

// StorageGetEncryptionKey retrieves the current encryption key
func StorageGetEncryptionKey() string {
    fmt.Println("Encryption key retrieved successfully")
    return common.EncryptionKey
}

// StorageRemoveEncryptionKey removes the encryption key
func StorageRemoveEncryptionKey() {
    common.EncryptionKey = ""
    ledger.RecordEncryptionKeyRemoval(time.Now())
    fmt.Println("Encryption key removed successfully")
}

// StorageEnableEncryption enables encryption for storage operations
func StorageEnableEncryption() {
    common.EncryptionEnabled = true
    ledger.RecordEncryptionEnabled(time.Now())
    fmt.Println("Encryption enabled for storage operations")
}

// StorageDisableEncryption disables encryption for storage operations
func StorageDisableEncryption() {
    common.EncryptionEnabled = false
    ledger.RecordEncryptionDisabled(time.Now())
    fmt.Println("Encryption disabled for storage operations")
}

// StorageFetchEncryptionStatus checks if encryption is enabled
func StorageFetchEncryptionStatus() bool {
    fmt.Printf("Encryption status: %v\n", common.EncryptionEnabled)
    return common.EncryptionEnabled
}

// StorageSetAccessPermissions sets access permissions for a data item
func StorageSetAccessPermissions(itemID, permissions string) error {
    if err := common.SetAccessPermissions(itemID, permissions); err != nil {
        return fmt.Errorf("setting access permissions failed: %v", err)
    }
    ledger.RecordAccessPermissionsSet(itemID, permissions, time.Now())
    fmt.Printf("Permissions set for item %s to %s\n", itemID, permissions)
    return nil
}

// StorageGetAccessPermissions retrieves access permissions for a data item
func StorageGetAccessPermissions(itemID string) (string, error) {
    permissions, err := common.GetAccessPermissions(itemID)
    if err != nil {
        return "", fmt.Errorf("failed to get access permissions: %v", err)
    }
    fmt.Printf("Permissions for item %s: %s\n", itemID, permissions)
    return permissions, nil
}

// StorageRevokeAccessPermissions revokes access permissions for a data item
func StorageRevokeAccessPermissions(itemID string) error {
    if err := common.RevokeAccessPermissions(itemID); err != nil {
        return fmt.Errorf("revoking access permissions failed: %v", err)
    }
    ledger.RecordAccessPermissionsRevoked(itemID, time.Now())
    fmt.Printf("Access permissions revoked for item %s\n", itemID)
    return nil
}

// StorageMonitorAccessPermissions monitors access permissions across storage
func StorageMonitorAccessPermissions() error {
    permissionsStatus, err := common.MonitorAccessPermissions()
    if err != nil {
        return fmt.Errorf("monitoring access permissions failed: %v", err)
    }
    fmt.Printf("Access Permissions Status: %v\n", permissionsStatus)
    return nil
}

// StorageAuditPermissions audits permissions to ensure compliance
func StorageAuditPermissions() error {
    if !common.VerifyPermissionsCompliance() {
        return fmt.Errorf("permissions audit failed")
    }
    fmt.Println("Permissions passed compliance audit")
    return nil
}

// StorageEncryptVirtualDisk encrypts a virtual disk
func StorageEncryptVirtualDisk(diskID string, data []byte) error {
    encryptedData, err := StorageEncryptData(data)
    if err != nil {
        return fmt.Errorf("disk encryption failed: %v", err)
    }
    if err := common.SaveToVirtualDisk(diskID, encryptedData); err != nil {
        return fmt.Errorf("saving encrypted data to disk failed: %v", err)
    }
    ledger.RecordVirtualDiskEncryption(diskID, time.Now())
    fmt.Printf("Virtual disk %s encrypted successfully\n", diskID)
    return nil
}

// StorageDecryptVirtualDisk decrypts a virtual disk
func StorageDecryptVirtualDisk(diskID string) ([]byte, error) {
    encryptedData, err := common.ReadFromVirtualDisk(diskID)
    if err != nil {
        return nil, fmt.Errorf("reading encrypted data from disk failed: %v", err)
    }
    data, err := StorageDecryptData(encryptedData)
    if err != nil {
        return nil, fmt.Errorf("disk decryption failed: %v", err)
    }
    fmt.Printf("Virtual disk %s decrypted successfully\n", diskID)
    return data, nil
}

// StorageFetchVirtualDiskEncryptionStatus retrieves the encryption status of a virtual disk
func StorageFetchVirtualDiskEncryptionStatus(diskID string) bool {
    status := common.CheckVirtualDiskEncryptionStatus(diskID)
    fmt.Printf("Encryption status for virtual disk %s: %v\n", diskID, status)
    return status
}

// StorageEncryptFile encrypts a file
func StorageEncryptFile(filePath string, data []byte) error {
    encryptedData, err := StorageEncryptData(data)
    if err != nil {
        return fmt.Errorf("file encryption failed: %v", err)
    }
    if err := common.WriteFile(filePath, encryptedData); err != nil {
        return fmt.Errorf("saving encrypted file failed: %v", err)
    }
    ledger.RecordFileEncryption(filePath, time.Now())
    fmt.Printf("File %s encrypted successfully\n", filePath)
    return nil
}

// StorageDecryptFile decrypts a file
func StorageDecryptFile(filePath string) ([]byte, error) {
    encryptedData, err := common.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("reading encrypted file failed: %v", err)
    }
    data, err := StorageDecryptData(encryptedData)
    if err != nil {
        return nil, fmt.Errorf("file decryption failed: %v", err)
    }
    fmt.Printf("File %s decrypted successfully\n", filePath)
    return data, nil
}

// StorageFetchFileEncryptionStatus retrieves the encryption status of a file
func StorageFetchFileEncryptionStatus(filePath string) bool {
    status := common.CheckFileEncryptionStatus(filePath)
    fmt.Printf("Encryption status for file %s: %v\n", filePath, status)
    return status
}

// StorageSetFilePermissions sets permissions for a file
func StorageSetFilePermissions(filePath, permissions string) error {
    if err := common.SetFilePermissions(filePath, permissions); err != nil {
        return fmt.Errorf("setting file permissions failed: %v", err)
    }
    ledger.RecordFilePermissionsSet(filePath, permissions, time.Now())
    fmt.Printf("Permissions set for file %s to %s\n", filePath, permissions)
    return nil
}

// StorageGetFilePermissions retrieves permissions for a file
func StorageGetFilePermissions(filePath string) (string, error) {
    permissions, err := common.GetFilePermissions(filePath)
    if err != nil {
        return "", fmt.Errorf("failed to get file permissions: %v", err)
    }
    fmt.Printf("Permissions for file %s: %s\n", filePath, permissions)
    return permissions, nil
}

// StorageAuditFileAccess audits access to a file
func StorageAuditFileAccess(filePath string) error {
    if !common.VerifyFileAccessLogIntegrity(filePath) {
        return fmt.Errorf("file access audit failed for %s", filePath)
    }
    fmt.Printf("File access for %s passed audit\n", filePath)
    return nil
}

// StorageEnableVersionControl enables version control for a file
func StorageEnableVersionControl(filePath string) error {
    if err := common.EnableFileVersionControl(filePath); err != nil {
        return fmt.Errorf("enabling version control failed: %v", err)
    }
    ledger.RecordVersionControlEnabled(filePath, time.Now())
    fmt.Printf("Version control enabled for file %s\n", filePath)
    return nil
}

// StorageDisableVersionControl disables version control for a file
func StorageDisableVersionControl(filePath string) error {
    if err := common.DisableFileVersionControl(filePath); err != nil {
        return fmt.Errorf("disabling version control failed: %v", err)
    }
    ledger.RecordVersionControlDisabled(filePath, time.Now())
    fmt.Printf("Version control disabled for file %s\n", filePath)
    return nil
}
