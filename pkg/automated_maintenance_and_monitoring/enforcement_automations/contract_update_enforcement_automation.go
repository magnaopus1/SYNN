package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/contracts"
)

// Configuration for contract update enforcement automation
const (
	UpdateCheckInterval           = 10 * time.Second // Interval to check for contract updates
	MaxAllowedUpdateFrequency     = 5                // Maximum allowed updates per contract within an hour
	MaxAllowedViolations          = 3                // Maximum violations before restricting contract
)

// ContractUpdateEnforcementAutomation monitors and enforces contract update compliance
type ContractUpdateEnforcementAutomation struct {
	contractManager   *contracts.ContractManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	updateFrequency   map[string]int // Tracks the number of updates per contract within a time frame
	violationCount    map[string]int // Tracks update violations per contract
}

// NewContractUpdateEnforcementAutomation initializes the contract update enforcement automation
func NewContractUpdateEnforcementAutomation(contractManager *contracts.ContractManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *ContractUpdateEnforcementAutomation {
	return &ContractUpdateEnforcementAutomation{
		contractManager:   contractManager,
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		enforcementMutex:  enforcementMutex,
		updateFrequency:   make(map[string]int),
		violationCount:    make(map[string]int),
	}
}

// StartUpdateEnforcement begins continuous monitoring and enforcement of contract update compliance
func (automation *ContractUpdateEnforcementAutomation) StartUpdateEnforcement() {
	ticker := time.NewTicker(UpdateCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkUpdateCompliance()
		}
	}()
}

// checkUpdateCompliance monitors update frequencies and validates update compliance for each contract
func (automation *ContractUpdateEnforcementAutomation) checkUpdateCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, contractID := range automation.contractManager.GetUpdatedContracts() {
		updateValid := automation.validateUpdate(contractID)
		if !updateValid {
			automation.handleUpdateViolation(contractID)
		} else {
			automation.trackUpdateFrequency(contractID)
		}
	}
}

// validateUpdate checks if the recent update to a contract meets compliance requirements
func (automation *ContractUpdateEnforcementAutomation) validateUpdate(contractID string) bool {
	// Logic to validate contract update content, e.g., compliance with standards, change scope, etc.
	compliant, err := automation.contractManager.ValidateUpdateContent(contractID)
	if err != nil {
		fmt.Printf("Error validating update for contract %s: %v\n", contractID, err)
		return false
	}
	return compliant
}

// trackUpdateFrequency logs the frequency of updates for a contract and restricts if it exceeds allowed limits
func (automation *ContractUpdateEnforcementAutomation) trackUpdateFrequency(contractID string) {
	automation.updateFrequency[contractID]++

	if automation.updateFrequency[contractID] > MaxAllowedUpdateFrequency {
		fmt.Printf("Update frequency violation for contract %s.\n", contractID)
		automation.handleUpdateViolation(contractID)
	}
}

// handleUpdateViolation manages contracts that exceed allowed update limits or fail compliance checks
func (automation *ContractUpdateEnforcementAutomation) handleUpdateViolation(contractID string) {
	automation.violationCount[contractID]++

	if automation.violationCount[contractID] >= MaxAllowedViolations {
		err := automation.contractManager.RestrictContract(contractID)
		if err != nil {
			fmt.Printf("Failed to restrict contract %s for update violations: %v\n", contractID, err)
			automation.logUpdateAction(contractID, "Failed Update Restriction")
		} else {
			fmt.Printf("Contract %s restricted due to repeated update violations.\n", contractID)
			automation.logUpdateAction(contractID, "Contract Restricted for Update Violations")
			automation.violationCount[contractID] = 0
		}
	} else {
		fmt.Printf("Compliance violation detected for contract update %s.\n", contractID)
		automation.logUpdateAction(contractID, "Update Compliance Violation Detected")
	}
}

// logUpdateAction securely logs actions related to contract update enforcement
func (automation *ContractUpdateEnforcementAutomation) logUpdateAction(contractID, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Contract ID: %s", action, contractID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("update-enforcement-%s-%d", contractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Contract Update Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log update enforcement action for contract %s: %v\n", contractID, err)
	} else {
		fmt.Println("Update enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ContractUpdateEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualUpdateCheck allows administrators to manually check update compliance for a specific contract
func (automation *ContractUpdateEnforcementAutomation) TriggerManualUpdateCheck(contractID string) {
	fmt.Printf("Manually triggering update compliance check for contract: %s\n", contractID)

	valid := automation.validateUpdate(contractID)
	if !valid {
		automation.handleUpdateViolation(contractID)
	} else {
		fmt.Printf("Contract %s is compliant with update standards.\n", contractID)
		automation.logUpdateAction(contractID, "Manual Update Compliance Check Passed")
	}
}
