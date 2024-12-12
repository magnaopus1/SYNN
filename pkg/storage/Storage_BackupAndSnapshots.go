package storage

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
)

// StorageBackupData backs up data and stores it in the specified location
func StorageBackupData(data []byte, location string) error {
    encryptedData, err := encryption.EncryptData(data, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("encryption failed: %v", err)
    }
    if err := common.SaveToStorage(location, encryptedData); err != nil {
        return fmt.Errorf("backup storage failed: %v", err)
    }
    ledger.RecordBackup("data", location, time.Now())
    fmt.Printf("Data backed up to %s\n", location)
    return nil
}

// StorageRestoreData restores data from the backup location
func StorageRestoreData(location string) ([]byte, error) {
    encryptedData, err := common.LoadFromStorage(location)
    if err != nil {
        return nil, fmt.Errorf("data retrieval failed: %v", err)
    }
    data, err := encryption.DecryptData(encryptedData, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %v", err)
    }
    fmt.Printf("Data restored from %s\n", location)
    return data, nil
}

// StorageSetBackupFrequency sets the frequency for automated backups
func StorageSetBackupFrequency(frequency time.Duration) {
    common.SetBackupFrequency(frequency)
    fmt.Printf("Backup frequency set to %s\n", frequency)
}

// StorageFetchBackupFrequency retrieves the current backup frequency
func StorageFetchBackupFrequency() time.Duration {
    frequency := common.GetBackupFrequency()
    fmt.Printf("Current backup frequency: %s\n", frequency)
    return frequency
}

// StorageAuditBackup audits backups to verify integrity
func StorageAuditBackup(location string) error {
    if !common.VerifyBackupIntegrity(location) {
        return fmt.Errorf("backup integrity check failed for %s", location)
    }
    fmt.Printf("Backup at %s passed integrity check\n", location)
    return nil
}

// StorageGenerateBackupReport generates a report of all backup activities
func StorageGenerateBackupReport() error {
    report, err := ledger.GenerateBackupReport()
    if err != nil {
        return fmt.Errorf("failed to generate backup report: %v", err)
    }
    fmt.Println("Backup Report:\n", report)
    return nil
}

// StorageSetBackupLocation sets the storage location for backups
func StorageSetBackupLocation(location string) {
    common.SetBackupLocation(location)
    fmt.Printf("Backup location set to %s\n", location)
}

// StorageFetchBackupLocation retrieves the current backup location
func StorageFetchBackupLocation() string {
    location := common.GetBackupLocation()
    fmt.Printf("Current backup location: %s\n", location)
    return location
}

// StorageOptimizeBackup optimizes storage utilization for backups
func StorageOptimizeBackup(location string) error {
    if err := common.OptimizeStorage(location); err != nil {
        return fmt.Errorf("storage optimization failed: %v", err)
    }
    fmt.Printf("Backup storage at %s optimized\n", location)
    return nil
}

// StorageDeleteBackup deletes a specified backup
func StorageDeleteBackup(location string) error {
    if err := common.DeleteFromStorage(location); err != nil {
        return fmt.Errorf("backup deletion failed: %v", err)
    }
    ledger.RecordBackupDeletion(location, time.Now())
    fmt.Printf("Backup deleted from %s\n", location)
    return nil
}

// StorageFetchBackupHistory retrieves the history of backups
func StorageFetchBackupHistory() ([]string, error) {
    history, err := ledger.FetchBackupHistory()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch backup history: %v", err)
    }
    fmt.Println("Backup History:\n", history)
    return history, nil
}

// StorageMonitorBackupUsage monitors usage of backup storage space
func StorageMonitorBackupUsage() error {
    usage, err := common.MonitorStorageUsage()
    if err != nil {
        return fmt.Errorf("storage usage monitoring failed: %v", err)
    }
    fmt.Printf("Current backup storage usage: %v\n", usage)
    return nil
}

// StorageCreateSnapshot creates a data snapshot
func StorageCreateSnapshot(data []byte) (string, error) {
    snapshotID, err := common.CreateSnapshot(data)
    if err != nil {
        return "", fmt.Errorf("snapshot creation failed: %v", err)
    }
    ledger.RecordSnapshotCreation(snapshotID, time.Now())
    fmt.Printf("Snapshot created with ID: %s\n", snapshotID)
    return snapshotID, nil
}

// StorageRestoreSnapshot restores data from a snapshot
func StorageRestoreSnapshot(snapshotID string) ([]byte, error) {
    snapshot, err := common.RestoreSnapshot(snapshotID)
    if err != nil {
        return nil, fmt.Errorf("snapshot restore failed: %v", err)
    }
    fmt.Printf("Snapshot %s restored\n", snapshotID)
    return snapshot, nil
}

// StorageDeleteSnapshot deletes a snapshot
func StorageDeleteSnapshot(snapshotID string) error {
    if err := common.DeleteSnapshot(snapshotID); err != nil {
        return fmt.Errorf("snapshot deletion failed: %v", err)
    }
    ledger.RecordSnapshotDeletion(snapshotID, time.Now())
    fmt.Printf("Snapshot %s deleted\n", snapshotID)
    return nil
}

// StorageFetchSnapshotHistory retrieves the history of snapshots
func StorageFetchSnapshotHistory() ([]string, error) {
    history, err := ledger.FetchSnapshotHistory()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch snapshot history: %v", err)
    }
    fmt.Println("Snapshot History:\n", history)
    return history, nil
}

// StorageSetSnapshotFrequency sets the frequency for creating snapshots
func StorageSetSnapshotFrequency(frequency time.Duration) {
    common.SetSnapshotFrequency(frequency)
    fmt.Printf("Snapshot frequency set to %s\n", frequency)
}

// StorageGetSnapshotFrequency retrieves the snapshot creation frequency
func StorageGetSnapshotFrequency() time.Duration {
    frequency := common.GetSnapshotFrequency()
    fmt.Printf("Current snapshot frequency: %s\n", frequency)
    return frequency
}

// StorageMonitorSnapshotUsage monitors usage of snapshot storage space
func StorageMonitorSnapshotUsage() error {
    usage, err := common.MonitorSnapshotStorageUsage()
    if err != nil {
        return fmt.Errorf("snapshot storage usage monitoring failed: %v", err)
    }
    fmt.Printf("Current snapshot storage usage: %v\n", usage)
    return nil
}

// StorageAuditSnapshot audits a snapshot for data integrity
func StorageAuditSnapshot(snapshotID string) error {
    if !common.VerifySnapshotIntegrity(snapshotID) {
        return fmt.Errorf("snapshot integrity check failed for %s", snapshotID)
    }
    fmt.Printf("Snapshot %s passed integrity check\n", snapshotID)
    return nil
}

// StorageBackupVirtualDisk backs up a virtual disk
func StorageBackupVirtualDisk(diskID string, location string) error {
    data, err := common.GetVirtualDiskData(diskID)
    if err != nil {
        return fmt.Errorf("failed to retrieve virtual disk data: %v", err)
    }
    return StorageBackupData(data, location)
}

// StorageRestoreVirtualDisk restores a virtual disk from a backup
func StorageRestoreVirtualDisk(location string, diskID string) error {
    data, err := StorageRestoreData(location)
    if err != nil {
        return fmt.Errorf("failed to restore virtual disk: %v", err)
    }
    if err := common.SetVirtualDiskData(diskID, data); err != nil {
        return fmt.Errorf("failed to set virtual disk data: %v", err)
    }
    fmt.Printf("Virtual disk %s restored from %s\n", diskID, location)
    return nil
}

// StorageAuditVirtualDisk audits a virtual disk for integrity
func StorageAuditVirtualDisk(diskID string) error {
    if !common.VerifyVirtualDiskIntegrity(diskID) {
        return fmt.Errorf("virtual disk integrity check failed for %s", diskID)
    }
    fmt.Printf("Virtual disk %s passed integrity check\n", diskID)
    return nil
}

// StorageBackupFile backs up an individual file
func StorageBackupFile(filePath string, location string) error {
    data, err := common.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("file read failed: %v", err)
    }
    return StorageBackupData(data, location)
}

// StorageRestoreFile restores an individual file from a backup
func StorageRestoreFile(location string, filePath string) error {
    data, err := StorageRestoreData(location)
    if err != nil {
        return fmt.Errorf("file restore failed: %v", err)
    }
    if err := common.WriteFile(filePath, data); err != nil {
        return fmt.Errorf("file write failed: %v", err)
    }
    fmt.Printf("File restored to %s from %s\n", filePath, location)
    return nil
}
