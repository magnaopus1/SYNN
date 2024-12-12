package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/smart_contracts"
	"synnergy_network_demo/synnergy_consensus"
)

const (
	SmartContractCheckInterval        = 1 * time.Minute  // Interval for checking and validating pending smart contracts
	SmartContractExecutionLedgerType  = "Smart Contract Execution"
)

// SmartContractExecutionAutomation handles the validation and execution of smart contracts
type SmartContractExecutionAutomation struct {
	consensusEngine   *synnergy_consensus.SynnergyConsensus   // Synnergy Consensus engine for validating contracts
	ledgerInstance    *ledger.Ledger                         // Ledger instance for logging contract events
	smartContractManager *smart_contracts.Manager             // Manager for handling smart contracts
	executionMutex    *sync.RWMutex                          // Mutex for thread-safe operations
}

// NewSmartContractExecutionAutomation initializes the smart contract execution automation
func NewSmartContractExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, smartContractManager *smart_contracts.Manager, executionMutex *sync.RWMutex) *SmartContractExecutionAutomation {
	return &SmartContractExecutionAutomation{
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		smartContractManager: smartContractManager,
		executionMutex:    executionMutex,
	}
}

// StartSmartContractMonitor begins the continuous monitoring for smart contract validation and execution
func (automation *SmartContractExecutionAutomation) StartSmartContractMonitor() {
	ticker := time.NewTicker(SmartContractCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkAndExecuteSmartContracts()
		}
	}()
}

// checkAndExecuteSmartContracts fetches pending smart contracts and triggers validation and execution
func (automation *SmartContractExecutionAutomation) checkAndExecuteSmartContracts() {
	automation.executionMutex.Lock()
	defer automation.executionMutex.Unlock()

	pendingContracts := automation.smartContractManager.GetPendingContracts()

	for _, contract := range pendingContracts {
		automation.validateAndExecuteContract(contract)
	}
}

// validateAndExecuteContract validates a smart contract via the Synnergy Consensus and executes it
func (automation *SmartContractExecutionAutomation) validateAndExecuteContract(contract *smart_contracts.SmartContract) {
	// Validate the smart contract using the Synnergy Consensus engine
	valid, err := automation.consensusEngine.ValidateContract(contract)
	if err != nil {
		fmt.Printf("Failed to validate contract %s: %v\n", contract.ID, err)
		return
	}

	if !valid {
		fmt.Printf("Smart contract %s validation failed.\n", contract.ID)
		automation.logContractFailure(contract, "Validation failed")
		return
	}

	// Execute the smart contract after validation
	err = automation.smartContractManager.ExecuteContract(contract)
	if err != nil {
		fmt.Printf("Failed to execute contract %s: %v\n", contract.ID, err)
		automation.logContractFailure(contract, "Execution failed")
		return
	}

	// Log successful execution into the ledger
	automation.logContractSuccess(contract)
}

// logContractSuccess logs the successful execution of a smart contract in the ledger
func (automation *SmartContractExecutionAutomation) logContractSuccess(contract *smart_contracts.SmartContract) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("smart-contract-success-%s-%d", contract.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      SmartContractExecutionLedgerType,
		Status:    "Success",
		Details:   fmt.Sprintf("Smart contract %s successfully executed.", contract.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log successful contract execution for contract %s: %v\n", contract.ID, err)
	} else {
		fmt.Println("Smart contract execution successfully logged in the ledger.")
	}
}

// logContractFailure logs a failure event when a smart contract fails validation or execution
func (automation *SmartContractExecutionAutomation) logContractFailure(contract *smart_contracts.SmartContract, reason string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("smart-contract-failure-%s-%d", contract.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      SmartContractExecutionLedgerType,
		Status:    "Failure",
		Details:   fmt.Sprintf("Smart contract %s failed: %s", contract.ID, reason),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log contract failure for contract %s: %v\n", contract.ID, err)
	} else {
		fmt.Println("Smart contract failure logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *SmartContractExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualContractExecution allows administrators to manually trigger the execution of a smart contract
func (automation *SmartContractExecutionAutomation) TriggerManualContractExecution(contractID string) {
	fmt.Printf("Manually triggering smart contract execution for contract %s...\n", contractID)

	contract := automation.smartContractManager.GetContractByID(contractID)
	if contract != nil {
		automation.validateAndExecuteContract(contract)
	} else {
		fmt.Printf("Smart contract %s not found.\n", contractID)
	}
}
