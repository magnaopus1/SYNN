package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/smartcontract"
)

// Configuration for smart contract rule enforcement automation
const (
	ContractCheckInterval             = 20 * time.Second // Interval to check smart contract compliance
	MaxAllowedViolations              = 2                // Allowed violations before enforcement action
)

// SmartContractRuleEnforcementAutomation monitors and enforces rules for deployed smart contracts
type SmartContractRuleEnforcementAutomation struct {
	contractManager     *smartcontract.ContractManager
	consensusEngine     *consensus.SynnergyConsensus
	ledgerInstance      *ledger.Ledger
	enforcementMutex    *sync.RWMutex
	contractViolationMap map[string]int // Tracks violation count for each contract
}

// NewSmartContractRuleEnforcementAutomation initializes the smart contract rule enforcement automation
func NewSmartContractRuleEnforcementAutomation(contractManager *smartcontract.ContractManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *SmartContractRuleEnforcementAutomation {
	return &SmartContractRuleEnforcementAutomation{
		contractManager:      contractManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		contractViolationMap: make(map[string]int),
	}
}

// StartSmartContractRuleEnforcement begins continuous monitoring and enforcement of smart contract rules
func (automation *SmartContractRuleEnforcementAutomation) StartSmartContractRuleEnforcement() {
	ticker := time.NewTicker(ContractCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkContractCompliance()
		}
	}()
}

// checkContractCompliance monitors each smart contract for rule violations and enforces actions if necessary
func (automation *SmartContractRuleEnforcementAutomation) checkContractCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.evaluateContractRules()
	automation.enforceContractCompliance()
}

// evaluateContractRules checks each contract for rule violations and flags contracts exceeding limits
func (automation *SmartContractRuleEnforcementAutomation) evaluateContractRules() {
	for _, contractID := range automation.contractManager.GetAllContracts() {
		violations := automation.contractManager.GetContractViolations(contractID)
		automation.contractViolationMap[contractID] = violations

		if violations > 0 {
			fmt.Printf("Contract %s has %d rule violation(s).\n", contractID, violations)
		}
	}
}

// enforceContractCompliance takes action on contracts that exceed allowed rule violations
func (automation *SmartContractRuleEnforcementAutomation) enforceContractCompliance() {
	for contractID, violations := range automation.contractViolationMap {
		if violations > MaxAllowedViolations {
			fmt.Printf("Enforcing compliance action on contract %s due to excessive rule violations.\n", contractID)
			automation.disableContract(contractID)
		}
	}
}

// disableContract disables a contract that has exceeded the rule violation limit
func (automation *SmartContractRuleEnforcementAutomation) disableContract(contractID string) {
	err := automation.contractManager.DisableContract(contractID)
	if err != nil {
		fmt.Printf("Failed to disable contract %s: %v\n", contractID, err)
		automation.logContractAction(contractID, "Disable Failed", fmt.Sprintf("Violations: %d", automation.contractViolationMap[contractID]))
	} else {
		fmt.Printf("Contract %s has been disabled due to rule violations.\n", contractID)
		automation.logContractAction(contractID, "Contract Disabled", fmt.Sprintf("Violations: %d", automation.contractViolationMap[contractID]))
	}
}

// logContractAction securely logs actions related to smart contract rule enforcement
func (automation *SmartContractRuleEnforcementAutomation) logContractAction(contractID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Contract ID: %s, Details: %s", action, contractID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("contract-rule-enforcement-%s-%d", contractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Smart Contract Rule Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log contract rule enforcement action for contract %s: %v\n", contractID, err)
	} else {
		fmt.Println("Smart contract rule enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *SmartContractRuleEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualComplianceEnforcement allows administrators to manually enforce compliance on a specific contract
func (automation *SmartContractRuleEnforcementAutomation) TriggerManualComplianceEnforcement(contractID string) {
	fmt.Printf("Manually enforcing compliance for contract: %s\n", contractID)

	if automation.contractViolationMap[contractID] > MaxAllowedViolations {
		automation.disableContract(contractID)
	} else {
		fmt.Printf("Contract %s is within compliance limits, no action taken.\n", contractID)
		automation.logContractAction(contractID, "Manual Enforcement Skipped", "Within Compliance Limits")
	}
}
