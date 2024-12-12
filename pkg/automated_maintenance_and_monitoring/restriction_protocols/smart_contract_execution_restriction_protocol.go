package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

// Constants for contract monitoring and restrictions
const (
	ContractCheckInterval         = 10 * time.Second // Frequency of contract execution checks
	MaxFailedExecutions           = 5                // Maximum allowed failed executions before restricting
	UnauthorizedExecutionMessage  = "Unauthorized smart contract execution detected"
)

// SmartContractExecutionRestrictionAutomation handles restrictions for smart contract executions in the network
type SmartContractExecutionRestrictionAutomation struct {
	consensusSystem         *consensus.SynnergyConsensus
	ledgerInstance          *ledger.Ledger
	stateMutex              *sync.RWMutex
	failedExecutionCount    map[string]int // Track the number of failed executions for each contract
	restrictedContracts     map[string]bool // Track restricted contracts
}

// NewSmartContractExecutionRestrictionAutomation initializes the restriction automation
func NewSmartContractExecutionRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SmartContractExecutionRestrictionAutomation {
	return &SmartContractExecutionRestrictionAutomation{
		consensusSystem:      consensusSystem,
		ledgerInstance:       ledgerInstance,
		stateMutex:           stateMutex,
		failedExecutionCount: make(map[string]int),
		restrictedContracts:  make(map[string]bool),
	}
}

// StartContractExecutionMonitoring starts the continuous monitoring of smart contract executions
func (automation *SmartContractExecutionRestrictionAutomation) StartContractExecutionMonitoring() {
	ticker := time.NewTicker(ContractCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateSmartContractExecutions()
		}
	}()
}

// evaluateSmartContractExecutions checks contract executions for failures or unauthorized activities
func (automation *SmartContractExecutionRestrictionAutomation) evaluateSmartContractExecutions() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	contracts := automation.consensusSystem.GetActiveContracts()

	for contractID, executionStatus := range contracts {
		if executionStatus == "failed" {
			automation.failedExecutionCount[contractID]++
			automation.logFailedExecution(contractID)

			// Restrict the contract if failure count exceeds the threshold
			if automation.failedExecutionCount[contractID] >= MaxFailedExecutions {
				automation.restrictSmartContract(contractID)
			}
		} else if executionStatus == "unauthorized" {
			automation.logUnauthorizedExecution(contractID)
			automation.restrictSmartContract(contractID)
		} else {
			// Reset failed execution count if contract executed successfully
			automation.failedExecutionCount[contractID] = 0
		}
	}
}

// logFailedExecution logs smart contract execution failures in the ledger
func (automation *SmartContractExecutionRestrictionAutomation) logFailedExecution(contractID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("contract-failure-%s-%d", contractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Smart Contract Execution Failure",
		Status:    "Failure",
		Details:   fmt.Sprintf("Smart contract %s failed to execute correctly.", contractID),
	}

	// Encrypt the failure details before logging in the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log contract execution failure:", err)
	} else {
		fmt.Println("Contract execution failure logged for:", contractID)
	}
}

// logUnauthorizedExecution logs unauthorized contract execution attempts in the ledger
func (automation *SmartContractExecutionRestrictionAutomation) logUnauthorizedExecution(contractID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("unauthorized-execution-%s-%d", contractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Unauthorized Smart Contract Execution",
		Status:    "Critical",
		Details:   fmt.Sprintf("Unauthorized execution detected for smart contract %s.", contractID),
	}

	// Encrypt the unauthorized access details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log unauthorized contract execution:", err)
	} else {
		fmt.Println("Unauthorized contract execution logged for:", contractID)
	}
}

// restrictSmartContract restricts a smart contract based on failures or unauthorized actions
func (automation *SmartContractExecutionRestrictionAutomation) restrictSmartContract(contractID string) {
	fmt.Printf("Smart contract %s has exceeded failure threshold or unauthorized execution. Access restricted.\n", contractID)

	// Log restriction in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("contract-restriction-%s-%d", contractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Smart Contract Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Smart contract %s restricted due to failure threshold or unauthorized execution.", contractID),
	}

	// Encrypt the restriction details before logging into the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log contract restriction:", err)
	} else {
		fmt.Println("Smart contract restricted:", contractID)
	}

	// Update the consensus system to restrict contract activity
	automation.consensusSystem.RestrictContract(contractID)
}

// encryptData encrypts sensitive information before storing in the ledger
func (automation *SmartContractExecutionRestrictionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting contract data:", err)
		return data
	}
	return string(encryptedData)
}

// SimulateUnauthorizedExecution simulates an unauthorized contract execution for testing purposes
func (automation *SmartContractExecutionRestrictionAutomation) SimulateUnauthorizedExecution(contractID string) {
	fmt.Println(UnauthorizedExecutionMessage)

	// Log unauthorized execution in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("unauthorized-execution-%s-%d", contractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Unauthorized Execution",
		Status:    "Critical",
		Details:   fmt.Sprintf("Unauthorized execution simulated for smart contract %s.", contractID),
	}

	// Encrypt the details before adding to the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log unauthorized execution:", err)
	} else {
		fmt.Println("Simulated unauthorized execution logged for:", contractID)
	}

	// Immediately restrict smart contract due to simulated unauthorized execution
	automation.restrictSmartContract(contractID)
}
