package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/smartcontracts"
)

const (
	BatchExecutionInterval   = 10 * time.Minute // Interval for batch smart contract execution
	MaxBatchSize             = 100              // Maximum number of contracts to execute in a batch
	ExecutionTimeout         = 5 * time.Minute  // Timeout for each contract execution
	ContractValidationPeriod = 2 * time.Minute  // Time allocated for consensus validation of each contract
)

// BatchSmartContractExecutionAutomation handles batch execution of smart contracts
type BatchSmartContractExecutionAutomation struct {
	consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine
	ledgerInstance    *ledger.Ledger                        // Ledger instance for recording contract executions
	stateMutex        *sync.RWMutex                         // Mutex for thread-safe operations
	contractQueue     *smartcontracts.ContractQueue         // Queue holding pending contracts for batch execution
	batchSize         int                                   // Number of contracts per batch execution
}

// NewBatchSmartContractExecutionAutomation initializes the batch execution automation for smart contracts
func NewBatchSmartContractExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex, contractQueue *smartcontracts.ContractQueue) *BatchSmartContractExecutionAutomation {
	return &BatchSmartContractExecutionAutomation{
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		stateMutex:        stateMutex,
		contractQueue:     contractQueue,
		batchSize:         MaxBatchSize,
	}
}

// StartBatchExecutionMonitor starts the automation for periodically executing smart contracts in batches
func (automation *BatchSmartContractExecutionAutomation) StartBatchExecutionMonitor() {
	ticker := time.NewTicker(BatchExecutionInterval)

	go func() {
		for range ticker.C {
			automation.executeBatchContracts()
		}
	}()
}

// executeBatchContracts retrieves contracts from the queue, validates them with consensus, and executes them
func (automation *BatchSmartContractExecutionAutomation) executeBatchContracts() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	contracts := automation.contractQueue.GetNextBatch(automation.batchSize)
	if len(contracts) == 0 {
		fmt.Println("No contracts to execute in the current batch.")
		return
	}

	fmt.Printf("Executing batch of %d contracts...\n", len(contracts))

	for _, contract := range contracts {
		automation.validateAndExecuteContract(contract)
	}
}

// validateAndExecuteContract performs consensus validation and executes the smart contract if valid
func (automation *BatchSmartContractExecutionAutomation) validateAndExecuteContract(contract smartcontracts.SmartContract) {
	validationSuccess := automation.consensusEngine.ValidateContractExecution(contract, ContractValidationPeriod)
	if validationSuccess {
		fmt.Printf("Contract %s validated successfully. Executing contract.\n", contract.ID)

		executionSuccess := automation.executeSmartContract(contract)
		if executionSuccess {
			automation.logContractExecution("Success", contract.ID, "Contract executed successfully.")
		} else {
			automation.logContractExecution("Failed", contract.ID, "Contract execution failed.")
		}
	} else {
		fmt.Printf("Contract %s failed validation.\n", contract.ID)
		automation.logContractExecution("Failed", contract.ID, "Contract validation failed.")
	}
}

// executeSmartContract performs the actual execution of the smart contract
func (automation *BatchSmartContractExecutionAutomation) executeSmartContract(contract smartcontracts.SmartContract) bool {
	err := contract.Execute(ExecutionTimeout)
	if err != nil {
		fmt.Printf("Error executing contract %s: %v\n", contract.ID, err)
		return false
	}
	return true
}

// logContractExecution logs the execution result of a smart contract into the ledger
func (automation *BatchSmartContractExecutionAutomation) logContractExecution(status, contractID, details string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("contract-execution-%s-%d", contractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Smart Contract Execution",
		Status:    status,
		Details:   fmt.Sprintf("Contract ID: %s, Details: %s", contractID, details),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	if err := automation.ledgerInstance.AddEntry(entry); err != nil {
		fmt.Printf("Error logging contract execution: %v\n", err)
	} else {
		fmt.Printf("Contract execution logged successfully: %s\n", contractID)
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *BatchSmartContractExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
