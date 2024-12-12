package high_availability

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewDataBackupManager initializes a DataBackupManager with a backup interval and location
func NewDataBackupManager(ledgerInstance *ledger.Ledger, backupInterval time.Duration, backupLocation string) *DataBackupManager {
    return &DataBackupManager{
        LedgerInstance: ledgerInstance,
        BackupInterval: backupInterval,
        BackupLocation: backupLocation,
    }
}

// StartAutoBackup starts an automated backup process at a specified interval
func (dbm *DataBackupManager) StartAutoBackup() {
    fmt.Printf("Starting auto-backup every %s...\n", dbm.BackupInterval)
    
    ticker := time.NewTicker(dbm.BackupInterval)
    go func() {
        for range ticker.C {
            dbm.mutex.Lock()
            err := dbm.BackupLedger()
            dbm.mutex.Unlock()
            if err != nil {
                fmt.Printf("Backup failed: %v\n", err)
            } else {
                fmt.Println("Backup successful.")
            }
        }
    }()
}

// BackupLedger creates a backup of the current ledger state and stores it in a file
func (dbm *DataBackupManager) BackupLedger() error {
    fmt.Println("Backing up ledger...")

    // Create backup file path with timestamp
    backupFilePath := fmt.Sprintf("%s/ledger_backup_%d.json", dbm.BackupLocation, time.Now().Unix())

    // Serialize the ledger state
    ledgerData, err := json.MarshalIndent(dbm.LedgerInstance, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal ledger data: %v", err)
    }

    // Write the ledger data to the backup file
    err = os.WriteFile(backupFilePath, ledgerData, 0644)
    if err != nil {
        return fmt.Errorf("failed to write backup file: %v", err)
    }

    fmt.Printf("Ledger backup created at %s.\n", backupFilePath)
    return nil
}

// RestoreLedger restores the ledger from a backup file
func (dbm *DataBackupManager) RestoreLedger(backupFilePath string) error {
    fmt.Printf("Restoring ledger from backup: %s...\n", backupFilePath)

    // Read the backup file
    backupData, err := os.ReadFile(backupFilePath)
    if err != nil {
        return fmt.Errorf("failed to read backup file: %v", err)
    }

    // Deserialize the backup data into the ledger instance
    dbm.mutex.Lock()
    defer dbm.mutex.Unlock()
    err = json.Unmarshal(backupData, dbm.LedgerInstance)
    if err != nil {
        return fmt.Errorf("failed to unmarshal backup data: %v", err)
    }

    fmt.Println("Ledger restored successfully.")
    return nil
}

// ListBackups lists all available backups in the backup location
func (dbm *DataBackupManager) ListBackups() ([]string, error) {
    fmt.Println("Listing available backups...")

    files, err := os.ReadDir(dbm.BackupLocation)
    if err != nil {
        return nil, fmt.Errorf("failed to list backup files: %v", err)
    }

    var backups []string
    for _, file := range files {
        if !file.IsDir() {
            backups = append(backups, file.Name())
        }
    }

    return backups, nil
}

// RemoveOldBackups removes backups older than a specified duration
func (dbm *DataBackupManager) RemoveOldBackups(retentionDuration time.Duration) error {
    fmt.Printf("Removing backups older than %s...\n", retentionDuration)

    files, err := os.ReadDir(dbm.BackupLocation)
    if err != nil {
        return fmt.Errorf("failed to list backup files: %v", err)
    }

    for _, file := range files {
        if !file.IsDir() {
            fileInfo, err := os.Stat(fmt.Sprintf("%s/%s", dbm.BackupLocation, file.Name()))
            if err != nil {
                return fmt.Errorf("failed to get file info for %s: %v", file.Name(), err)
            }

            // If the backup is older than the retention duration, delete it
            if time.Since(fileInfo.ModTime()) > retentionDuration {
                err := os.Remove(fmt.Sprintf("%s/%s", dbm.BackupLocation, file.Name()))
                if err != nil {
                    return fmt.Errorf("failed to remove backup file %s: %v", file.Name(), err)
                }
                fmt.Printf("Removed old backup: %s\n", file.Name())
            }
        }
    }

    return nil
}

// RetrieveLatestBackup retrieves the latest backup for a given node
func (dbm *DataBackupManager) RetrieveLatestBackup(nodeID string) (*BlockchainBackup, error) {
    dbm.mutex.Lock()
    defer dbm.mutex.Unlock()

    // Check if backups exist for the given node
    backups, exists := dbm.Backups[nodeID]
    if !exists || len(backups) == 0 {
        return nil, fmt.Errorf("no backups found for node %s", nodeID)
    }

    // Sort the backups by timestamp (if not already sorted)
    sort.Slice(backups, func(i, j int) bool {
        return backups[i].Timestamp.After(backups[j].Timestamp)
    })

    // Return the latest backup (which should be the first after sorting)
    return backups[0], nil
}