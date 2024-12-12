package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/dapps"
)

// Configuration for dApp upgrade enforcement automation
const (
	UpgradeCheckInterval        = 15 * time.Second // Interval to check for dApp upgrades
	MaxAllowedUpgradeFrequency  = 3                // Maximum upgrades allowed per dApp in a 24-hour period
	MaxUpgradeViolations        = 2                // Maximum upgrade violations before restricting dApp
)

// DAppUpgradeEnforcementAutomation monitors and enforces compliance of dApp upgrades
type DAppUpgradeEnforcementAutomation struct {
	dAppManager       *dapps.DAppManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	upgradeFrequency  map[string]int // Tracks upgrade frequency per dApp within a time period
	violationCount    map[string]int // Tracks upgrade violations per dApp
}

// NewDAppUpgradeEnforcementAutomation initializes the dApp upgrade enforcement automation
func NewDAppUpgradeEnforcementAutomation(dAppManager *dapps.DAppManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *DAppUpgradeEnforcementAutomation {
	return &DAppUpgradeEnforcementAutomation{
		dAppManager:      dAppManager,
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		enforcementMutex: enforcementMutex,
		upgradeFrequency: make(map[string]int),
		violationCount:   make(map[string]int),
	}
}

// StartUpgradeEnforcement begins continuous monitoring and enforcement of dApp upgrade compliance
func (automation *DAppUpgradeEnforcementAutomation) StartUpgradeEnforcement() {
	ticker := time.NewTicker(UpgradeCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkUpgradeCompliance()
		}
	}()
}

// checkUpgradeCompliance monitors and validates upgrades to ensure compliance with upgrade frequency and standards
func (automation *DAppUpgradeEnforcementAutomation) checkUpgradeCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, dAppID := range automation.dAppManager.GetUpgradedDApps() {
		if !automation.validateUpgrade(dAppID) {
			automation.handleUpgradeViolation(dAppID)
		} else {
			automation.trackUpgradeFrequency(dAppID)
		}
	}
}

// validateUpgrade verifies if an upgrade meets compliance standards, including validation through consensus
func (automation *DAppUpgradeEnforcementAutomation) validateUpgrade(dAppID string) bool {
	valid, err := automation.consensusEngine.ValidateUpgrade(dAppID)
	if err != nil {
		fmt.Printf("Consensus validation error for dApp %s upgrade: %v\n", dAppID, err)
		automation.logUpgradeAction(dAppID, "Upgrade Validation Error")
		return false
	}

	return valid
}

// trackUpgradeFrequency records the upgrade frequency and enforces restrictions if thresholds are exceeded
func (automation *DAppUpgradeEnforcementAutomation) trackUpgradeFrequency(dAppID string) {
	automation.upgradeFrequency[dAppID]++

	if automation.upgradeFrequency[dAppID] > MaxAllowedUpgradeFrequency {
		fmt.Printf("Upgrade frequency violation detected for dApp %s.\n", dAppID)
		automation.handleUpgradeViolation(dAppID)
	}
}

// handleUpgradeViolation restricts dApps with repeated upgrade violations
func (automation *DAppUpgradeEnforcementAutomation) handleUpgradeViolation(dAppID string) {
	automation.violationCount[dAppID]++

	if automation.violationCount[dAppID] >= MaxUpgradeViolations {
		err := automation.dAppManager.RestrictDApp(dAppID)
		if err != nil {
			fmt.Printf("Failed to restrict dApp %s due to upgrade violations: %v\n", dAppID, err)
			automation.logUpgradeAction(dAppID, "Failed Upgrade Restriction")
		} else {
			fmt.Printf("dApp %s restricted due to repeated upgrade violations.\n", dAppID)
			automation.logUpgradeAction(dAppID, "dApp Restricted for Upgrade Violations")
			automation.violationCount[dAppID] = 0
		}
	} else {
		fmt.Printf("Upgrade compliance violation detected for dApp %s.\n", dAppID)
		automation.logUpgradeAction(dAppID, "Upgrade Compliance Violation Detected")
	}
}

// logUpgradeAction securely logs actions related to dApp upgrade enforcement
func (automation *DAppUpgradeEnforcementAutomation) logUpgradeAction(dAppID, action string) {
	entryDetails := fmt.Sprintf("Action: %s, dApp ID: %s", action, dAppID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("dapp-upgrade-enforcement-%s-%d", dAppID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "dApp Upgrade Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log upgrade enforcement action for dApp %s: %v\n", dAppID, err)
	} else {
		fmt.Println("Upgrade enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DAppUpgradeEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualUpgradeValidation allows administrators to manually validate an upgrade for a specific dApp
func (automation *DAppUpgradeEnforcementAutomation) TriggerManualUpgradeValidation(dAppID string) {
	fmt.Printf("Manually triggering upgrade validation for dApp: %s\n", dAppID)

	if !automation.validateUpgrade(dAppID) {
		automation.handleUpgradeViolation(dAppID)
	} else {
		fmt.Printf("dApp %s upgrade is compliant with standards.\n", dAppID)
		automation.logUpgradeAction(dAppID, "Manual Upgrade Validation Passed")
	}
}
