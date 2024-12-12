package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// ContractInvocation represents the invocation of a smart contract, including the contract, method, parameters, and the caller's details.
type ContractInvocation struct {
    ContractAddress string            // Address of the smart contract being invoked
    Method          string            // Method or function name being invoked in the contract
    Params          map[string]string // Parameters passed to the contract method (key-value pairs)
    CallerAddress   string            // Address of the caller who is invoking the contract
    GasLimit        uint64            // Maximum gas allowed for the invocation
    GasPrice        uint64            // Gas price to pay for the execution
    Timestamp       time.Time         // Time when the invocation was made
    Signature       string            // Signature of the caller (for authorization)
    Nonce           uint64            // Unique value to prevent replay attacks
}

// ExecutionRecord represents the log of a smart contract execution.
type SmartContractExecutionRecord struct {
	ContractID  string    // ID of the executed contract
	Timestamp   time.Time // Timestamp of the execution
	Result      string    // Result of the execution (success, failure, etc.)
}

// SmartContract represents a smart contract in the Synnergy blockchain.
type SmartContract struct {
    ID             string                 // Unique ID of the contract
    Consensus      *SynnergyConsensus        // Consensus mechanism for validating transactions
    Code           string                 // Smart contract code (e.g., bytecode or source)
    Parameters     map[string]interface{} // Default parameters (initial values or configuration)
    Bytecode       string                // Compiled bytecode of the smart contrac
    State          map[string]interface{} // Current state (storage variables, balances, etc.)
    Owner          string                 // Owner of the smart contract (e.g., the creator’s wallet address)
    Executions     []ContractExecution    // History of contract executions
    CreationTime   time.Time              // Time the contract was created
    LastModified   time.Time              // Last time the contract was modified or executed
    IsActive       bool                   // Indicates if the contract is active or not
    mutex          sync.Mutex             // Mutex for safe concurrency
    LedgerInstance *ledger.Ledger         // Ledger instance for storing contract data and actions
}

// ContractExecution represents the execution log for a smart contract.
type ContractExecution struct {
    ExecutionID string                 // Unique ID for the execution
    ContractID  string                 // The contract ID that was executed
    Executor    string                 // The entity executing the contract
    Timestamp   time.Time              // When the execution happened
    Result      map[string]interface{} // The result of the execution
}

// SmartContractManager handles the lifecycle of smart contracts
type SmartContractManager struct {
    Contracts      map[string]*SmartContract // Mapping of contract ID to smart contract
    LedgerInstance *ledger.Ledger            // Reference to the ledger for recording contract actions
    mutex          sync.Mutex                // Mutex for thread-safe operations
    Consensus      *SynnergyConsensus        // Consensus mechanism for validating transactions

}

// NewSmartContractManager initializes the smart contract manager.
func NewSmartContractManager(ledgerInstance *ledger.Ledger) *SmartContractManager {
    return &SmartContractManager{
        Contracts:      make(map[string]*SmartContract),
        LedgerInstance: ledgerInstance,
    }
}

// DeployContract deploys a new smart contract to the blockchain and the virtual machine.
func (scm *SmartContractManager) DeployContract(owner, code string, parameters map[string]interface{}) (*SmartContract, error) {
    scm.mutex.Lock()
    defer scm.mutex.Unlock()

    // Generate a unique contract ID
    contractID := generateContractID(owner, code)

    // Create a new SmartContract struct
    contract := &SmartContract{
        ID:             contractID,
        Code:           code,
        Parameters:     parameters,
        State:          make(map[string]interface{}),
        Owner:          owner,
        LedgerInstance: scm.LedgerInstance,
    }

    // Store the contract in the manager's map
    scm.Contracts[contract.ID] = contract
    fmt.Printf("Smart Contract %s deployed by %s.\n", contract.ID, owner)

    // Encrypt and store the contract deployment on the ledger
    encryptionKey := []byte("your-encryption-key") // Replace with a secure key
    encryptedContract, err := EncryptContract(contract, encryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt contract: %v", err)
    }

    // Record the deployment in the ledger (direct blockchain deployment)
    err = scm.LedgerInstance.SmartContractLedger.RecordContractDeployment(contract.ID, owner, string(encryptedContract))
    if err != nil {
        return nil, fmt.Errorf("failed to record contract deployment in the ledger: %v", err)
    }

    // Initialize the virtual machine with correct parameters for deployment
    taskComplexity := 5 // Set complexity level for VM selection
    loggingEnabled := true
    virtualMachine, err := NewVirtualMachine(
        taskComplexity,          // Complexity level
        scm.LedgerInstance,      // Ledger instance
        scm.Consensus,           // Consensus mechanism, if applicable
        loggingEnabled,          // Enable logging
    )
    if err != nil {
        return nil, fmt.Errorf("failed to initialize virtual machine: %v", err)
    }

    // Load the contract code into the virtual machine
    err = virtualMachine.LoadContract(contract.Code, contract.Parameters)
    if err != nil {
        return nil, fmt.Errorf("failed to load contract into virtual machine: %v", err)
    }

    return contract, nil
}


// ExecuteContract executes the code of the smart contract using the virtual machine.
func (sc *SmartContract) ExecuteContract(sender string, parameters map[string]interface{}) (map[string]interface{}, error) {
    sc.mutex.Lock()
    defer sc.mutex.Unlock()

    // Validate contract owner or authorized executor
    if sender != sc.Owner {
        return nil, fmt.Errorf("unauthorized sender: %s is not the owner of the contract", sender)
    }

    // Define task complexity and logging requirements
    taskComplexity := 5 // Set a suitable task complexity level
    loggingEnabled := true

    // Initialize the virtual machine with correct parameters for execution
    vm, err := NewVirtualMachine(taskComplexity, sc.LedgerInstance, sc.Consensus, loggingEnabled)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize virtual machine: %v", err)
    }

    // Execute the contract code in the VM using the VMInterface
    result, err := vm.ExecuteContract(sc.ID, sc.Code, "solidity", parameters, []byte("encryption-key")) // Replace "encryption-key" with the actual key
    if err != nil {
        return nil, fmt.Errorf("contract execution failed in virtual machine: %v", err)
    }

    // Update the contract's state
    sc.State = result

    // Create an execution log for the contract
    execution := ContractExecution{
        ExecutionID: fmt.Sprintf("%s-exec-%d", sc.ID, time.Now().UnixNano()), // Unique ID for the execution
        ContractID:  sc.ID,
        Executor:    sender,
        Timestamp:   time.Now(),
        Result:      result,
    }
    sc.Executions = append(sc.Executions, execution)

    // Encrypt and store the execution result on the ledger
    encryptedExecution, err := EncryptContractExecution(execution, []byte("encryption-key")) // Use the same encryption key
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt contract execution: %v", err)
    }

    // Prepare the execution data to store in the ledger
    executionData := map[string]interface{}{
        "execution_id":     execution.ExecutionID,
        "contract_id":      execution.ContractID,
        "executor":         execution.Executor,
        "timestamp":        execution.Timestamp,
        "encrypted_result": string(encryptedExecution),
    }

    // Record the contract execution on the ledger
    err = sc.LedgerInstance.VirtualMachineLedger.RecordContractExecution(sc.ID, executionData)
    if err != nil {
        return nil, fmt.Errorf("failed to record contract execution in the ledger: %v", err)
    }

    fmt.Printf("Smart Contract %s executed by %s with result: %v\n", sc.ID, sender, result)
    return result, nil
}




// generateContractID generates a unique ID for a smart contract based on its owner and code.
func generateContractID(owner, code string) string {
    hashInput := fmt.Sprintf("%s%s%d", owner, code, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}
