package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/system"
)

// Configuration for automatic upgrade checks
const (
	UpgradeCheckInterval        = 24 * time.Hour // Interval to check for upgrades
	RequiredConsensusPercentage = 0.8            // 80% consensus required for upgrade
	MaxRollbackAttempts         = 3              // Max number of rollbacks if upgrade fails
)

// AutomaticSystemUpgradeEnforcement manages and enforces system upgrades
type AutomaticSystemUpgradeEnforcement struct {
	systemManager     *system.SystemManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	rollbackAttempts  int
}

// NewAutomaticSystemUpgradeEnforcement initializes the system upgrade automation
func NewAutomaticSystemUpgradeEnforcement(systemManager *system.SystemManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *AutomaticSystemUpgradeEnforcement {
	return &AutomaticSystemUpgradeEnforcement{
		systemManager:     systemManager,
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		enforcementMutex:  enforcementMutex,
		rollbackAttempts:  0,
	}
}

// StartAutomaticUpgradeEnforcement begins continuous monitoring and enforcement for system upgrades
func (automation *AutomaticSystemUpgradeEnforcement) StartAutomaticUpgradeEnforcement() {
	ticker := time.NewTicker(UpgradeCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkForUpgrade()
		}
	}()
}

// checkForUpgrade initiates the upgrade process if conditions are met
func (automation *AutomaticSystemUpgradeEnforcement) checkForUpgrade() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	// Check if an upgrade is available
	if upgradeAvailable := automation.systemManager.IsUpgradeAvailable(); upgradeAvailable {
		automation.initiateUpgrade()
	}
}

// initiateUpgrade initiates a system upgrade and logs the action
func (automation *AutomaticSystemUpgradeEnforcement) initiateUpgrade() {
	consensusAchieved := automation.consensusEngine.CheckConsensus(RequiredConsensusPercentage)

	if consensusAchieved {
		err := automation.systemManager.ExecuteUpgrade()
		if err != nil {
			fmt.Println("Upgrade failed, initiating rollback.")
			automation.handleUpgradeFailure()
		} else {
			fmt.Println("System upgrade executed successfully.")
			automation.logUpgradeAction("Upgrade Successful")
		}
	} else {
		fmt.Println("Consensus not reached for system upgrade.")
	}
}

// handleUpgradeFailure manages failed upgrades, including rollback attempts
func (automation *AutomaticSystemUpgradeEnforcement) handleUpgradeFailure() {
	automation.rollbackAttempts++

	if automation.rollbackAttempts <= MaxRollbackAttempts {
		err := automation.systemManager.RollbackUpgrade()
		if err != nil {
			fmt.Printf("Rollback attempt %d failed: %v\n", automation.rollbackAttempts, err)
		} else {
			fmt.Printf("Rollback successful on attempt %d.\n", automation.rollbackAttempts)
			automation.logUpgradeAction("Upgrade Rolled Back")
		}
	} else {
		fmt.Println("Max rollback attempts reached, marking system as critical.")
		automation.logUpgradeAction("Critical: Upgrade Failure")
	}
}

// logUpgradeAction securely logs actions related to system upgrades
func (automation *AutomaticSystemUpgradeEnforcement) logUpgradeAction(action string) {
	entryDetails := fmt.Sprintf("Action: %s", action)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("system-upgrade-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "System Upgrade",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log system upgrade action: %v\n", err)
	} else {
		fmt.Println("System upgrade action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *AutomaticSystemUpgradeEnforcement) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualUpgrade allows administrators to manually trigger a system upgrade
func (automation *AutomaticSystemUpgradeEnforcement) TriggerManualUpgrade() {
	fmt.Println("Manually triggering system upgrade.")

	err := automation.systemManager.ExecuteUpgrade()
	if err != nil {
		fmt.Println("Manual upgrade failed, initiating rollback.")
		automation.handleUpgradeFailure()
	} else {
		fmt.Println("Manual system upgrade executed successfully.")
		automation.logUpgradeAction("Manual Upgrade Successful")
	}
}
