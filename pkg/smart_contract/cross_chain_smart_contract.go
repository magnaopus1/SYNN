package smart_contract

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// CrossChainSmartContract represents a cross-chain smart contract that interacts with multiple blockchains.
type CrossChainSmartContract struct {
	ID              string                  // Unique ID of the contract
	Code            string                  // Smart contract code (bytecode or source)
	Parameters      map[string]interface{}  // Contract parameters
	State           map[string]interface{}  // Current state of the contract
	Owner           string                  // Owner of the smart contract
	ConnectedChains []string                // List of connected blockchains
	Executions      []ledger.ContractExecution // History of contract executions
	mutex           sync.Mutex              // Mutex for concurrency safety
	LedgerInstance  *ledger.Ledger          // Ledger for storing contract data
}

// CrossChainContractManager manages cross-chain smart contracts.
type CrossChainContractManager struct {
	Contracts      map[string]*CrossChainSmartContract // All deployed cross-chain contracts
	LedgerInstance *ledger.Ledger                      // Ledger for contract deployments and executions
	mutex          sync.Mutex                          // Mutex for concurrency safety
}

// NewCrossChainContractManager initializes a new CrossChainContractManager.
func NewCrossChainContractManager(ledgerInstance *ledger.Ledger) *CrossChainContractManager {
	return &CrossChainContractManager{
		Contracts:      make(map[string]*CrossChainSmartContract),
		LedgerInstance: ledgerInstance,
	}
}

// DeployCrossChainContract deploys a new cross-chain smart contract with security and monitoring.
func (ccm *CrossChainContractManager) DeployCrossChainContract(owner, code string, parameters map[string]interface{}, connectedChains []string) (*CrossChainSmartContract, error) {
	ccm.mutex.Lock()
	defer ccm.mutex.Unlock()

	contractID := generateContractID(owner, code)
	contract := &CrossChainSmartContract{
		ID:              contractID,
		Code:            code,
		Parameters:      parameters,
		State:           make(map[string]interface{}),
		Owner:           owner,
		ConnectedChains: connectedChains,
		LedgerInstance:  ccm.LedgerInstance,
	}

	// Store contract in memory
	ccm.Contracts[contract.ID] = contract
	fmt.Printf("Cross-Chain Smart Contract %s deployed by %s.\n", contract.ID, owner)

	// Convert CrossChainSmartContract to common.SmartContract (assuming the conversion is required)
	commonContract := &common.SmartContract{
		ID:         contract.ID,
		Code:       contract.Code,
		Parameters: contract.Parameters,
		State:      contract.State,
		Owner:      contract.Owner,
	}

	// Encrypt and store the contract deployment on the ledger
	encryptedContract, err := common.EncryptContract(commonContract, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt cross-chain contract: %v", err)
	}

	// Convert connectedChains slice to a string (if needed, depending on ledger implementation)
	connectedChainsStr := strings.Join(connectedChains, ",")

	// Store connectedChains metadata in the ledger
	err = ccm.LedgerInstance.RecordContractDeployment(contract.ID, string(encryptedContract), connectedChainsStr)
	if err != nil {
		return nil, fmt.Errorf("failed to record cross-chain contract deployment in the ledger: %v", err)
	}

	// Real-time monitoring of contract deployment (ensure LogDeployment is properly implemented)
	ledger.LogDeployment(contractID, owner, connectedChainsStr)

	return contract, nil
}

// ExecuteCrossChainContract executes the smart contract across multiple blockchains using the virtual machine.
func (cc *CrossChainSmartContract) ExecuteCrossChainContract(vm common.VirtualMachine, sender string, parameters map[string]interface{}, encryption *common.Encryption, encryptedData []byte) (map[string]interface{}, error) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// Authorization check
	if sender != cc.Owner {
		return nil, fmt.Errorf("unauthorized sender: %s is not the owner of the contract", sender)
	}

	// Pre-execution validation (check connectivity, sync state across chains)
	for _, chain := range cc.ConnectedChains {
		err := ValidateBlockchainConnection(chain)
		if err != nil {
			// Use fallback mechanism to handle chain failures
			TriggerFallback(chain, cc.ID)
			return nil, fmt.Errorf("failed to validate connection to chain %s: %v", chain, err)
		}
	}

	// Execute the smart contract inside the virtual machine
	result, err := vm.ExecuteContract(cc.Code, cc.ID, sender, parameters, encryption, encryptedData)
	if err != nil {
		ledger.LogExecutionFailure(cc.ID, err.Error())
		return nil, fmt.Errorf("cross-chain contract execution failed: %v", err)
	}

	// Update contract state based on execution result
	cc.State = result

	// Record the execution in the ledger
	execution := ledger.ContractExecution{
		ExecutionID: fmt.Sprintf("%s-exec-%d", cc.ID, time.Now().UnixNano()),
		ContractID:  cc.ID,
		Executor:    sender,
		Timestamp:   time.Now(),
		Result:      result,
	}
	cc.Executions = append(cc.Executions, execution)

	// Encrypt the execution data
	encryptedExecution, err := common.EncryptContractExecution(execution, []byte{})
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt contract execution: %v", err)
	}

	// Prepare the map for the encrypted execution result
	encryptedExecutionData := map[string]interface{}{
		"execution_id": execution.ExecutionID,
		"contract_id":  cc.ID,
		"executor":     sender,
		"timestamp":    execution.Timestamp,
		"result":       string(encryptedExecution), // Use encrypted result as needed
	}

	err = cc.LedgerInstance.RecordContractExecution(
		execution.ExecutionID,
		encryptedExecutionData,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to record contract execution in the ledger: %v", err)
	}

	// Monitoring execution success
	ledger.LogExecution(execution.ExecutionID, cc.ID)

	fmt.Printf("Cross-Chain Smart Contract %s executed by %s.\n", cc.ID, sender)
	return result, nil
}



// QueryCrossChainContract queries the state of the cross-chain smart contract.
func (cc *CrossChainSmartContract) QueryCrossChainContract() ([]byte, error) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// Encrypt the contract state before returning
	encryptedState, err := common.EncryptContractState(cc.State, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt contract state: %v", err)
	}

	return encryptedState, nil
}


// ValidateContract ensures that the contract code and state remain consistent.
func (ccm *CrossChainContractManager) ValidateContract(contractID string) error {
	ccm.mutex.Lock()
	defer ccm.mutex.Unlock()

	contract, exists := ccm.Contracts[contractID]
	if !exists {
		return fmt.Errorf("contract %s not found", contractID)
	}

	// Initialize the Encryption object
	encryption := &common.Encryption{}

	// Decrypt the contract code for validation
	decryptedCode, err := encryption.DecryptData([]byte(contract.Code), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt contract code: %v", err)
	}

	// Generate the hash of the decrypted code
	hash := sha256.New()
	hash.Write(decryptedCode)
	calculatedHash := hex.EncodeToString(hash.Sum(nil))

	// Retrieve the ledger record for the contract
	ledgerRecord, err := ccm.LedgerInstance.RetrieveContractDeployment(contractID)
	if err != nil {
		return fmt.Errorf("failed to retrieve contract deployment from ledger: %v", err)
	}

	// Compare the ledger hash with the calculated hash
	if ledgerRecord.StoredHash != calculatedHash {
		_ = ccm.LedgerInstance.LogValidationFailure(contractID)
		return fmt.Errorf("contract validation failed: ledger data does not match")
	}

	_ = ccm.LedgerInstance.LogValidationSuccess(contractID)
	fmt.Printf("Contract %s validated successfully.\n", contractID)
	return nil
}





// simulateCrossChainExecution simulates smart contract execution on connected blockchains.
func simulateCrossChainExecution(code string, parameters map[string]interface{}, connectedChains []string) (map[string]interface{}, error) {
	// Simulate interaction with multiple blockchains
	fmt.Printf("Simulating cross-chain execution on chains: %v\n", connectedChains)

	// Validate the connection to external blockchains
	for _, chain := range connectedChains {
		err := ValidateBlockchainConnection(chain)
		if err != nil {
			return nil, fmt.Errorf("failed to validate connection to chain %s: %v", chain, err)
		}
	}

	// Simulate contract execution and return result
	result := make(map[string]interface{})
	result["status"] = "success"
	result["executed"] = true

	return result, nil
}

// generateContractID generates a unique contract ID based on the owner and contract code.
// You can customize this function to meet your requirements.
func generateContractID(owner, code string) string {
	// Generate an ID using the owner and contract code with the current Unix timestamp in nanoseconds.
	return fmt.Sprintf("%s-%s-%d", owner, code, time.Now().UnixNano())
}

// Modify ValidateBlockchainConnection to accept a single string
func ValidateBlockchainConnection(chain string) error {
	// Validate single chain connection logic here
	if chain == "" {
		return fmt.Errorf("invalid blockchain connection to chain: %s", chain)
	}
	// Add other validation logic here if needed
	return nil
}



// TriggerFallback handles chain failures by logging the issue or initiating a fallback mechanism.
func TriggerFallback(chain, contractID string) {
	// Log the failure or perform a rollback/recovery mechanism.
	fmt.Printf("Fallback triggered for chain %s on contract %s due to connection failure.\n", chain, contractID)
}
