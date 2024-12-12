package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/update"
)

// Configuration for hotfix deployment enforcement automation
const (
	HotfixCheckInterval        = 30 * time.Second // Interval to check for pending hotfix deployments
	MaxAllowedDowntime         = 5 * time.Minute  // Maximum allowed downtime before hotfix deployment is enforced
)

// HotfixDeploymentEnforcementAutomation monitors and enforces deployment of critical network hotfixes
type HotfixDeploymentEnforcementAutomation struct {
	updateManager       *update.UpdateManager
	consensusEngine     *consensus.SynnergyConsensus
	ledgerInstance      *ledger.Ledger
	enforcementMutex    *sync.RWMutex
	pendingHotfixMap    map[string]bool // Tracks pending hotfixes
}

// NewHotfixDeploymentEnforcementAutomation initializes the hotfix deployment enforcement automation
func NewHotfixDeploymentEnforcementAutomation(updateManager *update.UpdateManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *HotfixDeploymentEnforcementAutomation {
	return &HotfixDeploymentEnforcementAutomation{
		updateManager:      updateManager,
		consensusEngine:    consensusEngine,
		ledgerInstance:     ledgerInstance,
		enforcementMutex:   enforcementMutex,
		pendingHotfixMap:   make(map[string]bool),
	}
}

// StartHotfixEnforcement begins continuous monitoring and enforcement of hotfix deployments
func (automation *HotfixDeploymentEnforcementAutomation) StartHotfixEnforcement() {
	ticker := time.NewTicker(HotfixCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkHotfixCompliance()
		}
	}()
}

// checkHotfixCompliance monitors pending hotfixes and enforces deployment if necessary
func (automation *HotfixDeploymentEnforcementAutomation) checkHotfixCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.verifyPendingHotfixes()
}

// verifyPendingHotfixes ensures that all critical hotfixes are deployed within the allowable timeframe
func (automation *HotfixDeploymentEnforcementAutomation) verifyPendingHotfixes() {
	for _, hotfixID := range automation.updateManager.GetPendingHotfixes() {
		if !automation.pendingHotfixMap[hotfixID] {
			automation.deployHotfix(hotfixID)
		}
	}
}

// deployHotfix triggers the deployment of a pending hotfix and logs the action
func (automation *HotfixDeploymentEnforcementAutomation) deployHotfix(hotfixID string) {
	err := automation.updateManager.DeployHotfix(hotfixID)
	if err != nil {
		fmt.Printf("Failed to deploy hotfix %s: %v\n", hotfixID, err)
		automation.logHotfixAction(hotfixID, "Deployment Failed")
	} else {
		fmt.Printf("Hotfix %s deployed successfully.\n", hotfixID)
		automation.pendingHotfixMap[hotfixID] = true
		automation.logHotfixAction(hotfixID, "Deployed")
	}
}

// logHotfixAction securely logs actions related to hotfix deployment enforcement
func (automation *HotfixDeploymentEnforcementAutomation) logHotfixAction(hotfixID, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Hotfix ID: %s", action, hotfixID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("hotfix-enforcement-%s-%d", hotfixID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Hotfix Deployment Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log hotfix deployment enforcement action for hotfix %s: %v\n", hotfixID, err)
	} else {
		fmt.Println("Hotfix deployment enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *HotfixDeploymentEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualHotfixDeployment allows administrators to manually deploy a specific hotfix
func (automation *HotfixDeploymentEnforcementAutomation) TriggerManualHotfixDeployment(hotfixID string) {
	fmt.Printf("Manually triggering deployment for hotfix: %s\n", hotfixID)

	if automation.pendingHotfixMap[hotfixID] {
		fmt.Printf("Hotfix %s has already been deployed.\n", hotfixID)
		automation.logHotfixAction(hotfixID, "Manual Deployment Attempt Skipped - Already Deployed")
	} else {
		automation.deployHotfix(hotfixID)
	}
}
