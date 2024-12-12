package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/network"
)

const (
	VersionCheckInterval      = 15 * time.Minute // Interval for checking available version updates
	UpdateValidationTimeout   = 10 * time.Minute // Timeout for validating an update with consensus
	UpdateConfirmationTimeout = 20 * time.Minute // Timeout for receiving confirmation of a successful update
)

// AutomaticVersionUpdateExecutionAutomation handles automatic version updates of nodes
type AutomaticVersionUpdateExecutionAutomation struct {
	consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine
	ledgerInstance  *ledger.Ledger                        // Ledger instance for recording updates
	stateMutex      *sync.RWMutex                         // Mutex for thread-safe operations
	networkManager  *network.NetworkManager               // Network manager for fetching version updates
}

// NewAutomaticVersionUpdateExecutionAutomation initializes the automation for version updates
func NewAutomaticVersionUpdateExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex, networkManager *network.NetworkManager) *AutomaticVersionUpdateExecutionAutomation {
	return &AutomaticVersionUpdateExecutionAutomation{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		stateMutex:      stateMutex,
		networkManager:  networkManager,
	}
}

// StartVersionUpdateMonitor begins the continuous monitoring for version updates
func (automation *AutomaticVersionUpdateExecutionAutomation) StartVersionUpdateMonitor() {
	versionTicker := time.NewTicker(VersionCheckInterval)

	go func() {
		for range versionTicker.C {
			automation.checkForVersionUpdates()
		}
	}()
}

// checkForVersionUpdates fetches available updates and initiates the upgrade process if needed
func (automation *AutomaticVersionUpdateExecutionAutomation) checkForVersionUpdates() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	latestVersion, err := automation.networkManager.FetchLatestVersion()
	if err != nil {
		fmt.Println("Error fetching latest version:", err)
		return
	}

	currentVersion := automation.consensusEngine.GetCurrentVersion()

	if latestVersion != currentVersion {
		fmt.Printf("New version detected: %s (current version: %s)\n", latestVersion, currentVersion)
		automation.validateAndUpdateVersion(latestVersion)
	} else {
		fmt.Println("No new version available.")
	}
}

// validateAndUpdateVersion validates the new version and triggers the update process
func (automation *AutomaticVersionUpdateExecutionAutomation) validateAndUpdateVersion(newVersion string) {
	validationSuccess := automation.consensusEngine.ValidateVersionUpdate(newVersion, UpdateValidationTimeout)

	if validationSuccess {
		fmt.Printf("Version %s validated successfully. Initiating update.\n", newVersion)
		automation.initiateVersionUpdate(newVersion)
	} else {
		fmt.Printf("Version %s failed validation.\n", newVersion)
		automation.logVersionUpdate("Failed", newVersion, "Validation failed.")
	}
}

// initiateVersionUpdate handles the update process and logs it to the ledger
func (automation *AutomaticVersionUpdateExecutionAutomation) initiateVersionUpdate(newVersion string) {
	updateSuccess := automation.networkManager.ApplyVersionUpdate(newVersion, UpdateConfirmationTimeout)

	if updateSuccess {
		fmt.Printf("Version %s successfully applied.\n", newVersion)
		automation.consensusEngine.SetCurrentVersion(newVersion)
		automation.logVersionUpdate("Success", newVersion, "Version update applied successfully.")
	} else {
		fmt.Printf("Version %s update failed.\n", newVersion)
		automation.logVersionUpdate("Failed", newVersion, "Version update failed during application.")
	}
}

// logVersionUpdate records the version update process in the ledger with encryption
func (automation *AutomaticVersionUpdateExecutionAutomation) logVersionUpdate(status, version, details string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("version-update-%s-%d", version, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Version Update",
		Status:    status,
		Details:   fmt.Sprintf("Version: %s, Details: %s", version, details),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	if err := automation.ledgerInstance.AddEntry(entry); err != nil {
		fmt.Printf("Error logging version update: %v\n", err)
	} else {
		fmt.Printf("Version update logged successfully: %s\n", version)
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *AutomaticVersionUpdateExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
