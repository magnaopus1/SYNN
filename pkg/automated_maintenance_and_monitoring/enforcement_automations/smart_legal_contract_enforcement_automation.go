package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/legalcontract"
)

// Configuration for smart legal contract enforcement automation
const (
	LegalContractCheckInterval      = 30 * time.Second // Interval to check legal contract compliance
	MaxLegalViolationsAllowed       = 1                // Maximum allowed legal violations before enforcement action
)

// SmartLegalContractEnforcementAutomation monitors and enforces legal standards for deployed smart contracts
type SmartLegalContractEnforcementAutomation struct {
	legalContractManager  *legalcontract.LegalContractManager
	consensusEngine       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	enforcementMutex      *sync.RWMutex
	contractViolationMap  map[string]int // Tracks violation count for each legal contract
}

// NewSmartLegalContractEnforcementAutomation initializes the smart legal contract enforcement automation
func NewSmartLegalContractEnforcementAutomation(legalContractManager *legalcontract.LegalContractManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *SmartLegalContractEnforcementAutomation {
	return &SmartLegalContractEnforcementAutomation{
		legalContractManager: legalContractManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		contractViolationMap: make(map[string]int),
	}
}

// StartSmartLegalContractEnforcement begins continuous monitoring and enforcement of legal standards for smart contracts
func (automation *SmartLegalContractEnforcementAutomation) StartSmartLegalContractEnforcement() {
	ticker := time.NewTicker(LegalContractCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkLegalCompliance()
		}
	}()
}

// checkLegalCompliance monitors each smart legal contract for rule violations and enforces actions if necessary
func (automation *SmartLegalContractEnforcementAutomation) checkLegalCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.evaluateLegalCompliance()
	automation.enforceLegalStandards()
}

// evaluateLegalCompliance checks each contract for legal compliance violations and flags non-compliant contracts
func (automation *SmartLegalContractEnforcementAutomation) evaluateLegalCompliance() {
	for _, contractID := range automation.legalContractManager.GetAllLegalContracts() {
		violations := automation.legalContractManager.GetLegalViolations(contractID)
		automation.contractViolationMap[contractID] = violations

		if violations > 0 {
			fmt.Printf("Legal contract %s has %d legal violation(s).\n", contractID, violations)
		}
	}
}

// enforceLegalStandards takes action on contracts that exceed allowed legal violations
func (automation *SmartLegalContractEnforcementAutomation) enforceLegalStandards() {
	for contractID, violations := range automation.contractViolationMap {
		if violations > MaxLegalViolationsAllowed {
			fmt.Printf("Enforcing legal compliance action on contract %s due to legal violations.\n", contractID)
			automation.suspendContract(contractID)
		}
	}
}

// suspendContract suspends a contract that has exceeded legal compliance violations
func (automation *SmartLegalContractEnforcementAutomation) suspendContract(contractID string) {
	err := automation.legalContractManager.SuspendContract(contractID)
	if err != nil {
		fmt.Printf("Failed to suspend contract %s: %v\n", contractID, err)
		automation.logLegalContractAction(contractID, "Suspension Failed", fmt.Sprintf("Legal Violations: %d", automation.contractViolationMap[contractID]))
	} else {
		fmt.Printf("Contract %s has been suspended due to legal compliance violations.\n", contractID)
		automation.logLegalContractAction(contractID, "Contract Suspended", fmt.Sprintf("Legal Violations: %d", automation.contractViolationMap[contractID]))
	}
}

// logLegalContractAction securely logs actions related to smart legal contract enforcement
func (automation *SmartLegalContractEnforcementAutomation) logLegalContractAction(contractID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Contract ID: %s, Details: %s", action, contractID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("legal-contract-enforcement-%s-%d", contractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Smart Legal Contract Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log legal contract enforcement action for contract %s: %v\n", contractID, err)
	} else {
		fmt.Println("Legal contract enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *SmartLegalContractEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualLegalEnforcement allows administrators to manually enforce legal standards on a specific contract
func (automation *SmartLegalContractEnforcementAutomation) TriggerManualLegalEnforcement(contractID string) {
	fmt.Printf("Manually enforcing legal compliance for contract: %s\n", contractID)

	if automation.contractViolationMap[contractID] > MaxLegalViolationsAllowed {
		automation.suspendContract(contractID)
	} else {
		fmt.Printf("Contract %s is within legal compliance limits, no action taken.\n", contractID)
		automation.logLegalContractAction(contractID, "Manual Enforcement Skipped", "Within Compliance Limits")
	}
}
