package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/storage"
)

// Configuration for backup automation
const (
	BackupInterval        = 24 * time.Hour // Regular interval for scheduled backups
	MaxBackupRetries      = 3              // Maximum retry attempts for failed backups
	BackupIntegrityCheck  = 72 * time.Hour // Interval for checking backup integrity
)

// BackupScheduleEnforcementAutomation manages scheduled and triggered backups of blockchain data
type BackupScheduleEnforcementAutomation struct {
	storageManager    *storage.StorageManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	retryAttempts     map[string]int // Track retry attempts for each backup
}

// NewBackupScheduleEnforcementAutomation initializes the backup schedule automation
func NewBackupScheduleEnforcementAutomation(storageManager *storage.StorageManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *BackupScheduleEnforcementAutomation {
	return &BackupScheduleEnforcementAutomation{
		storageManager:   storageManager,
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		enforcementMutex: enforcementMutex,
		retryAttempts:    make(map[string]int),
	}
}

// StartBackupScheduleEnforcement begins continuous scheduling and monitoring of backups
func (automation *BackupScheduleEnforcementAutomation) StartBackupScheduleEnforcement() {
	backupTicker := time.NewTicker(BackupInterval)
	integrityCheckTicker := time.NewTicker(BackupIntegrityCheck)

	go func() {
		for {
			select {
			case <-backupTicker.C:
				automation.performScheduledBackup()
			case <-integrityCheckTicker.C:
				automation.verifyBackupIntegrity()
			}
		}
	}()
}

// performScheduledBackup initiates a scheduled backup and validates its completion
func (automation *BackupScheduleEnforcementAutomation) performScheduledBackup() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	// Trigger a backup and record success or failure
	backupID, err := automation.storageManager.CreateBackup()
	if err != nil {
		fmt.Println("Scheduled backup failed, retrying.")
		automation.handleBackupFailure(backupID)
		return
	}

	automation.logBackupAction(backupID, "Scheduled Backup Completed")
}

// handleBackupFailure manages backup retries and triggers additional backups if necessary
func (automation *BackupScheduleEnforcementAutomation) handleBackupFailure(backupID string) {
	automation.retryAttempts[backupID]++

	if automation.retryAttempts[backupID] <= MaxBackupRetries {
		err := automation.storageManager.CreateBackup()
		if err != nil {
			fmt.Printf("Backup retry attempt %d for backup ID %s failed.\n", automation.retryAttempts[backupID], backupID)
		} else {
			fmt.Printf("Backup retry attempt %d for backup ID %s succeeded.\n", automation.retryAttempts[backupID], backupID)
			automation.logBackupAction(backupID, "Backup Retry Successful")
			automation.retryAttempts[backupID] = 0
		}
	} else {
		fmt.Printf("Max backup retries reached for backup ID %s, marking as failed.\n", backupID)
		automation.logBackupAction(backupID, "Backup Failed After Retries")
	}
}

// verifyBackupIntegrity periodically checks the integrity of backups
func (automation *BackupScheduleEnforcementAutomation) verifyBackupIntegrity() {
	for backupID := range automation.storageManager.ListBackups() {
		integrityValid := automation.storageManager.CheckBackupIntegrity(backupID)
		if !integrityValid {
			fmt.Printf("Backup integrity check failed for backup ID %s, triggering a new backup.\n", backupID)
			automation.performBackupOnIntegrityFailure(backupID)
		} else {
			fmt.Printf("Backup integrity verified for backup ID %s.\n", backupID)
			automation.logBackupAction(backupID, "Backup Integrity Verified")
		}
	}
}

// performBackupOnIntegrityFailure triggers a backup if an integrity failure is detected
func (automation *BackupScheduleEnforcementAutomation) performBackupOnIntegrityFailure(backupID string) {
	newBackupID, err := automation.storageManager.CreateBackup()
	if err != nil {
		fmt.Printf("Failed to create new backup after integrity failure for backup ID %s: %v\n", backupID, err)
		automation.handleBackupFailure(newBackupID)
	} else {
		fmt.Printf("New backup created after integrity failure for backup ID %s. New backup ID: %s\n", backupID, newBackupID)
		automation.logBackupAction(newBackupID, "New Backup Created After Integrity Failure")
	}
}

// logBackupAction securely logs actions related to backups
func (automation *BackupScheduleEnforcementAutomation) logBackupAction(backupID string, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Backup ID: %s", action, backupID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("backup-action-%s-%d", backupID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Backup",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log backup action for backup ID %s: %v\n", backupID, err)
	} else {
		fmt.Println("Backup action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *BackupScheduleEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualBackup allows administrators to manually initiate a backup
func (automation *BackupScheduleEnforcementAutomation) TriggerManualBackup() {
	fmt.Println("Manually triggering a backup.")

	backupID, err := automation.storageManager.CreateBackup()
	if err != nil {
		fmt.Println("Manual backup failed, retrying.")
		automation.handleBackupFailure(backupID)
	} else {
		fmt.Println("Manual backup completed successfully.")
		automation.logBackupAction(backupID, "Manual Backup Completed")
	}
}
