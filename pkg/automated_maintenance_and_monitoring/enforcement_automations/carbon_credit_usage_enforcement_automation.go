package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/green"
)

// Configuration for carbon credit enforcement automation
const (
	CarbonCreditCheckInterval      = 30 * time.Second // Interval to check carbon credit usage
	CarbonCreditUsageThreshold     = 1000             // Max carbon credits allowed per entity
	CarbonCreditWarningThreshold   = 800              // Warning threshold for high carbon credit usage
	MaxCarbonCreditViolations      = 3                // Max violations before restriction
)

// CarbonCreditUsageEnforcementAutomation monitors and enforces carbon credit limits
type CarbonCreditUsageEnforcementAutomation struct {
	greenManager      *green.GreenManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	violationCount    map[string]int // Track carbon credit violations per entity
}

// NewCarbonCreditUsageEnforcementAutomation initializes the carbon credit usage automation
func NewCarbonCreditUsageEnforcementAutomation(greenManager *green.GreenManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *CarbonCreditUsageEnforcementAutomation {
	return &CarbonCreditUsageEnforcementAutomation{
		greenManager:      greenManager,
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		enforcementMutex:  enforcementMutex,
		violationCount:    make(map[string]int),
	}
}

// StartCarbonCreditUsageEnforcement begins continuous monitoring and enforcement of carbon credit usage
func (automation *CarbonCreditUsageEnforcementAutomation) StartCarbonCreditUsageEnforcement() {
	ticker := time.NewTicker(CarbonCreditCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkCarbonCreditUsage()
		}
	}()
}

// checkCarbonCreditUsage monitors entity carbon credit usage and applies restrictions if limits are exceeded
func (automation *CarbonCreditUsageEnforcementAutomation) checkCarbonCreditUsage() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, entityID := range automation.greenManager.GetRegisteredEntities() {
		creditsUsed := automation.greenManager.GetEntityCarbonCreditUsage(entityID)

		if creditsUsed >= CarbonCreditUsageThreshold {
			automation.handleCreditViolation(entityID, creditsUsed)
		} else if creditsUsed >= CarbonCreditWarningThreshold {
			fmt.Printf("Warning: Entity %s is approaching carbon credit limit with usage of %d credits.\n", entityID, creditsUsed)
			automation.logCreditAction(entityID, "High Carbon Credit Usage Warning", creditsUsed)
		}
	}
}

// handleCreditViolation applies restrictions on entities exceeding carbon credit limits
func (automation *CarbonCreditUsageEnforcementAutomation) handleCreditViolation(entityID string, creditsUsed int) {
	automation.violationCount[entityID]++

	if automation.violationCount[entityID] >= MaxCarbonCreditViolations {
		err := automation.greenManager.RestrictEntityCarbonCreditUsage(entityID)
		if err != nil {
			fmt.Printf("Failed to restrict carbon credit usage for entity %s: %v\n", entityID, err)
			automation.logCreditAction(entityID, "Failed Carbon Credit Restriction", creditsUsed)
		} else {
			fmt.Printf("Carbon credit restriction applied to entity %s after %d violations.\n", entityID, automation.violationCount[entityID])
			automation.logCreditAction(entityID, "Carbon Credit Usage Restricted", creditsUsed)
			automation.violationCount[entityID] = 0
		}
	} else {
		fmt.Printf("Entity %s has exceeded carbon credit limit with usage of %d credits.\n", entityID, creditsUsed)
		automation.logCreditAction(entityID, "Carbon Credit Limit Exceeded", creditsUsed)
	}
}

// logCreditAction securely logs actions related to carbon credit enforcement
func (automation *CarbonCreditUsageEnforcementAutomation) logCreditAction(entityID, action string, creditsUsed int) {
	entryDetails := fmt.Sprintf("Action: %s, Entity: %s, Carbon Credits Used: %d", action, entityID, creditsUsed)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("carbon-credit-enforcement-%s-%d", entityID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Carbon Credit Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log carbon credit enforcement action for entity %s: %v\n", entityID, err)
	} else {
		fmt.Println("Carbon credit enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *CarbonCreditUsageEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualCreditRestriction allows administrators to manually restrict carbon credit usage for an entity
func (automation *CarbonCreditUsageEnforcementAutomation) TriggerManualCreditRestriction(entityID string) {
	fmt.Printf("Manually triggering carbon credit restriction for entity: %s\n", entityID)

	err := automation.greenManager.RestrictEntityCarbonCreditUsage(entityID)
	if err != nil {
		fmt.Printf("Failed to manually restrict carbon credit usage for entity %s: %v\n", entityID, err)
		automation.logCreditAction(entityID, "Manual Carbon Credit Restriction Failed", 0)
	} else {
		fmt.Printf("Manual carbon credit restriction applied to entity %s.\n", entityID)
		automation.logCreditAction(entityID, "Manual Carbon Credit Restriction", 0)
	}
}
