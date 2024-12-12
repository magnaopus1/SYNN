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

// Configuration for carbon offset allocation automation
const (
	OffsetAllocationCheckInterval = 1 * time.Hour // Interval to check carbon offset allocations
	MinimumOffsetRequirement      = 500           // Minimum required carbon offsets per entity
	OffsetWarningThreshold        = 600           // Threshold for issuing a warning before reaching minimum requirement
	MaxOffsetViolations           = 3             // Maximum violations before restricting entity activity
)

// CarbonOffsetAllocationEnforcementAutomation monitors and enforces carbon offset allocations for compliance
type CarbonOffsetAllocationEnforcementAutomation struct {
	greenManager     *green.GreenManager
	consensusEngine  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	enforcementMutex *sync.RWMutex
	violationCount   map[string]int // Tracks offset violations per entity
}

// NewCarbonOffsetAllocationEnforcementAutomation initializes the carbon offset allocation automation
func NewCarbonOffsetAllocationEnforcementAutomation(greenManager *green.GreenManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *CarbonOffsetAllocationEnforcementAutomation {
	return &CarbonOffsetAllocationEnforcementAutomation{
		greenManager:     greenManager,
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		enforcementMutex: enforcementMutex,
		violationCount:   make(map[string]int),
	}
}

// StartCarbonOffsetEnforcement begins continuous monitoring and enforcement of carbon offset allocations
func (automation *CarbonOffsetAllocationEnforcementAutomation) StartCarbonOffsetEnforcement() {
	ticker := time.NewTicker(OffsetAllocationCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkOffsetAllocations()
		}
	}()
}

// checkOffsetAllocations monitors and enforces minimum carbon offset allocations for each registered entity
func (automation *CarbonOffsetAllocationEnforcementAutomation) checkOffsetAllocations() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, entityID := range automation.greenManager.GetRegisteredEntities() {
		offsetAllocated := automation.greenManager.GetEntityCarbonOffsets(entityID)

		if offsetAllocated < MinimumOffsetRequirement {
			automation.handleOffsetViolation(entityID, offsetAllocated)
		} else if offsetAllocated < OffsetWarningThreshold {
			fmt.Printf("Warning: Entity %s is approaching minimum offset allocation with %d credits.\n", entityID, offsetAllocated)
			automation.logOffsetAction(entityID, "Low Carbon Offset Allocation Warning", offsetAllocated)
		}
	}
}

// handleOffsetViolation applies restrictions if an entity consistently fails to meet offset requirements
func (automation *CarbonOffsetAllocationEnforcementAutomation) handleOffsetViolation(entityID string, offsetAllocated int) {
	automation.violationCount[entityID]++

	if automation.violationCount[entityID] >= MaxOffsetViolations {
		err := automation.greenManager.RestrictEntityActivity(entityID)
		if err != nil {
			fmt.Printf("Failed to restrict activity for entity %s: %v\n", entityID, err)
			automation.logOffsetAction(entityID, "Failed Activity Restriction", offsetAllocated)
		} else {
			fmt.Printf("Activity restriction applied to entity %s after %d violations.\n", entityID, automation.violationCount[entityID])
			automation.logOffsetAction(entityID, "Activity Restricted Due to Insufficient Carbon Offsets", offsetAllocated)
			automation.violationCount[entityID] = 0
		}
	} else {
		fmt.Printf("Entity %s does not meet minimum offset requirement with only %d credits.\n", entityID, offsetAllocated)
		automation.logOffsetAction(entityID, "Minimum Offset Requirement Not Met", offsetAllocated)
	}
}

// logOffsetAction securely logs actions related to carbon offset allocation enforcement
func (automation *CarbonOffsetAllocationEnforcementAutomation) logOffsetAction(entityID, action string, offsetAllocated int) {
	entryDetails := fmt.Sprintf("Action: %s, Entity: %s, Carbon Offsets Allocated: %d", action, entityID, offsetAllocated)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("carbon-offset-enforcement-%s-%d", entityID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Carbon Offset Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log carbon offset enforcement action for entity %s: %v\n", entityID, err)
	} else {
		fmt.Println("Carbon offset enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *CarbonOffsetAllocationEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualOffsetRestriction allows administrators to manually restrict activity for an entity not meeting offset requirements
func (automation *CarbonOffsetAllocationEnforcementAutomation) TriggerManualOffsetRestriction(entityID string) {
	fmt.Printf("Manually triggering activity restriction for entity: %s due to insufficient carbon offsets.\n", entityID)

	err := automation.greenManager.RestrictEntityActivity(entityID)
	if err != nil {
		fmt.Printf("Failed to manually restrict activity for entity %s: %v\n", entityID, err)
		automation.logOffsetAction(entityID, "Manual Activity Restriction Failed", 0)
	} else {
		fmt.Printf("Manual activity restriction applied to entity %s.\n", entityID)
		automation.logOffsetAction(entityID, "Manual Activity Restriction Due to Insufficient Carbon Offsets", 0)
	}
}
