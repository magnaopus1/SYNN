package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/system"
	"synnergy_network_demo/synnergy_consensus"
)

const (
	SystemUpdateCheckInterval   = 1 * time.Hour // Check for system updates every hour
	SystemUpdateLedgerEntryType = "System Update"
)

// SystemUpdateExecutionAutomation handles automatic checking and applying of system updates
type SystemUpdateExecutionAutomation struct {
	consensusEngine  *synnergy_consensus.SynnergyConsensus // Consensus engine for system validation
	ledgerInstance   *ledger.Ledger                        // Ledger for logging updates
	systemManager    *system.Manager                       // System manager for handling updates
	updateMutex      *sync.RWMutex                         // Mutex for thread-safe execution
}

// NewSystemUpdateExecutionAutomation initializes the automation for system updates
func NewSystemUpdateExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, systemManager *system.Manager, updateMutex *sync.RWMutex) *SystemUpdateExecutionAutomation {
	return &SystemUpdateExecutionAutomation{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		systemManager:   systemManager,
		updateMutex:     updateMutex,
	}
}

// StartSystemUpdateMonitor starts the continuous monitoring for system updates
func (automation *SystemUpdateExecutionAutomation) StartSystemUpdateMonitor() {
	ticker := time.NewTicker(SystemUpdateCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkAndApplySystemUpdate()
		}
	}()
}

// checkAndApplySystemUpdate checks for available system updates and applies them if required
func (automation *SystemUpdateExecutionAutomation) checkAndApplySystemUpdate() {
	automation.updateMutex.Lock()
	defer automation.updateMutex.Unlock()

	// Check for available updates in the system
	availableUpdates, err := automation.systemManager.CheckForUpdates()
	if err != nil {
		fmt.Printf("Error checking for system updates: %v\n", err)
		return
	}

	// If updates are available, proceed with the update process
	for _, update := range availableUpdates {
		automation.validateAndApplyUpdate(update)
	}
}

// validateAndApplyUpdate validates and applies the system update after consensus validation
func (automation *SystemUpdateExecutionAutomation) validateAndApplyUpdate(update *system.SystemUpdate) {
	// Validate the update with the Synnergy Consensus
	valid, err := automation.consensusEngine.ValidateSystemUpdate(update)
	if err != nil {
		fmt.Printf("Failed to validate system update %s: %v\n", update.ID, err)
		return
	}

	if !valid {
		fmt.Printf("System update %s failed consensus validation.\n", update.ID)
		automation.logUpdateFailure(update, "Consensus validation failed")
		return
	}

	// Apply the system update
	err = automation.systemManager.ApplyUpdate(update)
	if err != nil {
		fmt.Printf("Failed to apply system update %s: %v\n", update.ID, err)
		automation.logUpdateFailure(update, "Update application failed")
		return
	}

	// Log the successful update application in the ledger
	automation.logUpdateSuccess(update)
}

// logUpdateSuccess logs the successful application of a system update into the ledger
func (automation *SystemUpdateExecutionAutomation) logUpdateSuccess(update *system.SystemUpdate) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("system-update-success-%s-%d", update.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      SystemUpdateLedgerEntryType,
		Status:    "Success",
		Details:   fmt.Sprintf("System update %s successfully applied.", update.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log successful update for %s in the ledger: %v\n", update.ID, err)
	} else {
		fmt.Println("System update successfully logged in the ledger.")
	}
}

// logUpdateFailure logs a failure event when a system update fails validation or application
func (automation *SystemUpdateExecutionAutomation) logUpdateFailure(update *system.SystemUpdate, reason string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("system-update-failure-%s-%d", update.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      SystemUpdateLedgerEntryType,
		Status:    "Failure",
		Details:   fmt.Sprintf("System update %s failed: %s", update.ID, reason),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log update failure for %s in the ledger: %v\n", update.ID, err)
	} else {
		fmt.Println("Update failure successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *SystemUpdateExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualSystemUpdate allows administrators to manually trigger a system update
func (automation *SystemUpdateExecutionAutomation) TriggerManualSystemUpdate(updateID string) {
	fmt.Printf("Manually triggering system update %s...\n", updateID)

	update := automation.systemManager.GetUpdateByID(updateID)
	if update != nil {
		automation.validateAndApplyUpdate(update)
	} else {
		fmt.Printf("System update %s not found.\n", updateID)
	}
}
