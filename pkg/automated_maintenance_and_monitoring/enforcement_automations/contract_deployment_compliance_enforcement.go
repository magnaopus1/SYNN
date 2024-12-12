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

// Configuration for contract deployment compliance automation
const (
	ComplianceCheckInterval         = 15 * time.Second // Interval to check compliance of deployed contracts
	MaxComplianceViolations         = 3                // Maximum compliance violations before restriction
)

// ContractDeploymentComplianceEnforcement monitors and enforces compliance standards for contract deployments
type ContractDeploymentComplianceEnforcement struct {
	contractManager  *contracts.ContractManager
	consensusEngine  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	enforcementMutex *sync.RWMutex
	violationCount   map[string]int // Tracks compliance violations per contract
}

// NewContractDeploymentComplianceEnforcement initializes the contract compliance enforcement automation
func NewContractDeploymentComplianceEnforcement(contractManager *contracts.ContractManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *ContractDeploymentComplianceEnforcement {
	return &ContractDeploymentComplianceEnforcement{
		contractManager:  contractManager,
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		enforcementMutex: enforcementMutex,
		violationCount:   make(map[string]int),
	}
}

// StartContractComplianceEnforcement begins continuous monitoring and enforcement of contract deployment compliance
func (automation *ContractDeploymentComplianceEnforcement) StartContractComplianceEnforcement() {
	ticker := time.NewTicker(ComplianceCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkContractCompliance()
		}
	}()
}

// checkContractCompliance verifies that each deployed contract complies with network standards
func (automation *ContractDeploymentComplianceEnforcement) checkContractCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, contractID := range automation.contractManager.GetDeployedContracts() {
		compliant, err := automation.contractManager.CheckCompliance(contractID)
		if err != nil {
			fmt.Printf("Error checking compliance for contract %s: %v\n", contractID, err)
			continue
		}

		if !compliant {
			automation.handleComplianceViolation(contractID)
		}
	}
}

// handleComplianceViolation manages contracts that do not meet compliance standards
func (automation *ContractDeploymentComplianceEnforcement) handleComplianceViolation(contractID string) {
	automation.violationCount[contractID]++

	if automation.violationCount[contractID] >= MaxComplianceViolations {
		err := automation.contractManager.RestrictContract(contractID)
		if err != nil {
			fmt.Printf("Failed to restrict contract %s for compliance violations: %v\n", contractID, err)
			automation.logComplianceAction(contractID, "Failed Compliance Restriction")
		} else {
			fmt.Printf("Contract %s restricted due to repeated compliance violations.\n", contractID)
			automation.logComplianceAction(contractID, "Contract Restricted for Compliance Violations")
			automation.violationCount[contractID] = 0
		}
	} else {
		fmt.Printf("Compliance violation detected for contract %s.\n", contractID)
		automation.logComplianceAction(contractID, "Compliance Violation Detected")
	}
}

// logComplianceAction securely logs actions related to contract compliance enforcement
func (automation *ContractDeploymentComplianceEnforcement) logComplianceAction(contractID, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Contract ID: %s", action, contractID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("contract-compliance-enforcement-%s-%d", contractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Contract Compliance Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log contract compliance action for contract %s: %v\n", contractID, err)
	} else {
		fmt.Println("Contract compliance enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ContractDeploymentComplianceEnforcement) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualComplianceCheck allows administrators to manually check compliance for a specific contract
func (automation *ContractDeploymentComplianceEnforcement) TriggerManualComplianceCheck(contractID string) {
	fmt.Printf("Manually triggering compliance check for contract: %s\n", contractID)

	compliant, err := automation.contractManager.CheckCompliance(contractID)
	if err != nil {
		fmt.Printf("Failed to manually check compliance for contract %s: %v\n", contractID, err)
		automation.logComplianceAction(contractID, "Manual Compliance Check Failed")
		return
	}

	if !compliant {
		automation.handleComplianceViolation(contractID)
	} else {
		fmt.Printf("Contract %s is compliant.\n", contractID)
		automation.logComplianceAction(contractID, "Manual Compliance Check Passed")
	}
}
