package common

import (
	"errors"
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// Sandbox is an isolated virtual environment for testing and executing smart contracts.
type Sandbox struct {
	ID               string                     // Unique identifier for the sandbox instance
	SmartContracts   map[string]*SmartContract  // Loaded smart contracts in the sandbox
	ExecutionHistory []SmartContractExecutionRecord // Records of contract executions within the sandbox
	LedgerInstance   *ledger.Ledger             // Sandbox-associated ledger for tracking state
	VirtualMachine   VMInterface                // Use VMInterface to allow flexibility
	mutex            sync.Mutex                 // Mutex for thread-safe operations
	IsActive         bool                       // Indicates if the sandbox is active or paused
	CreatedAt        time.Time                  // Time of sandbox creation
}



// SandboxManager manages multiple sandboxes for isolated smart contract testing and execution.
type SandboxManager struct {
	LedgerInstance   *ledger.Ledger             // Ledger for logging sandbox actions
	Sandboxes        map[string]*Sandbox        // Map of sandbox instances
	mutex            sync.Mutex                 // Mutex for thread-safety
}

// NewSandboxManager initializes a new SandboxManager instance.
func NewSandboxManager(ledgerInstance *ledger.Ledger) *SandboxManager {
	return &SandboxManager{
		LedgerInstance: ledgerInstance,
		Sandboxes:      make(map[string]*Sandbox), // Initialize the map for sandbox instances
		mutex:          sync.Mutex{},
	}
}


// NewSandbox creates a new isolated sandbox instance for smart contract execution.
func NewSandbox(ledgerInstance *ledger.Ledger, sandboxManager *SandboxManager, synnergyConsensus *SynnergyConsensus) (*Sandbox, error) {
    // Set a task complexity level (e.g., 5 for basic sandbox execution)
    taskComplexity := 5
    loggingEnabled := true

    // Initialize the virtual machine with the required parameters
    virtualMachine, err := NewVirtualMachine(taskComplexity, ledgerInstance, synnergyConsensus, loggingEnabled)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize virtual machine: %v", err)
    }

    // Return the sandbox instance with the initialized virtual machine and other configurations
    return &Sandbox{
        ID:               generateSandboxID(),
        SmartContracts:   make(map[string]*SmartContract),
        ExecutionHistory: []SmartContractExecutionRecord{},
        LedgerInstance:   ledgerInstance,
        VirtualMachine:   virtualMachine, // Now it fully implements VMInterface
        IsActive:         true,
        CreatedAt:        time.Now(),
    }, nil
}


// LoadSmartContract loads a smart contract into the sandbox for testing.
func (sb *Sandbox) LoadSmartContract(contract *SmartContract) error {
    sb.mutex.Lock()
    defer sb.mutex.Unlock()

    if !sb.IsActive {
        return errors.New("sandbox is inactive")
    }

    if _, exists := sb.SmartContracts[contract.ID]; exists {
        return fmt.Errorf("contract with ID %s already loaded", contract.ID)
    }

    sb.SmartContracts[contract.ID] = contract
    fmt.Printf("Contract %s loaded into sandbox %s.\n", contract.ID, sb.ID)
    return nil
}

// CompileSmartContract compiles a smart contract code into bytecode in the sandbox environment.
func (sb *Sandbox) CompileSmartContract(contractID, contractCode, language string, parameters map[string]interface{}) (string, error) {
    sb.mutex.Lock()
    defer sb.mutex.Unlock()

    if !sb.IsActive {
        return "", errors.New("sandbox is inactive")
    }

    // Use the virtual machine's compiler to compile the smart contract code into bytecode
    bytecode, err := sb.VirtualMachine.Compile(contractID, contractCode, language, parameters)
    if err != nil {
        return "", fmt.Errorf("compilation error: %v", err)
    }

    // Store the compiled bytecode in the smart contract
    if contract, exists := sb.SmartContracts[contractID]; exists {
        contract.Code = contractCode
        contract.Bytecode = bytecode
        fmt.Printf("Contract %s compiled in sandbox %s.\n", contractID, sb.ID)
        return bytecode, nil
    }

    return "", fmt.Errorf("contract with ID %s not found in sandbox", contractID)
}



// DebugSmartContract runs the smart contract in the debugger in the sandbox.
func (sb *Sandbox) DebugSmartContract(contractID string, parameters map[string]interface{}) (string, error) {
    sb.mutex.Lock()
    defer sb.mutex.Unlock()

    if !sb.IsActive {
        return "", errors.New("sandbox is inactive")
    }

    contract, exists := sb.SmartContracts[contractID]
    if !exists {
        return "", fmt.Errorf("contract with ID %s not found in sandbox", contractID)
    }

    if contract.Bytecode == "" {
        return "", errors.New("contract bytecode not available, please compile the contract first")
    }

    // Step 1: Parse the bytecode into individual instructions using VMInterface
    instructions, err := sb.VirtualMachine.ParseBytecode(contract.Bytecode)
    if err != nil {
        return "", fmt.Errorf("failed to parse bytecode: %v", err)
    }

    // Convert []string to []interface{}
    instructionsInterface := make([]interface{}, len(instructions))
    for i, instruction := range instructions {
        instructionsInterface[i] = instruction
    }

    // Step 2: Simulate execution of bytecode instructions using VMInterface
    debugResult, err := sb.VirtualMachine.ExecuteInstructions(contractID, instructionsInterface, parameters)
    if err != nil {
        return "", fmt.Errorf("debugging error: %v", err)
    }

    fmt.Printf("Contract %s debugged in sandbox %s.\n", contractID, sb.ID)

    // Convert debugResult to string for returning
    debugResultStr := fmt.Sprintf("%v", debugResult) // Convert map[string]interface{} to string

    return debugResultStr, nil
}






// ExecuteSmartContract executes the contract in the sandbox and records the results.
func (sb *Sandbox) ExecuteSmartContract(contractID string, parameters map[string]interface{}) (string, error) {
    sb.mutex.Lock()
    defer sb.mutex.Unlock()

    if !sb.IsActive {
        return "", errors.New("sandbox is inactive")
    }

    contract, exists := sb.SmartContracts[contractID]
    if !exists {
        return "", fmt.Errorf("contract with ID %s not found in sandbox", contractID)
    }

    if contract.Bytecode == "" {
        return "", errors.New("contract bytecode not available, please compile the contract first")
    }

    // Step 1: Parse the bytecode into individual instructions using VMInterface
    instructions, err := sb.VirtualMachine.ParseBytecode(contract.Bytecode)
    if err != nil {
        return "", fmt.Errorf("failed to parse bytecode: %v", err)
    }

    // Convert []string to []interface{} for execution
    instructionsInterface := make([]interface{}, len(instructions))
    for i, instruction := range instructions {
        instructionsInterface[i] = instruction
    }

    // Step 2: Execute the instructions using VMInterface
    executionResult, err := sb.VirtualMachine.ExecuteInstructions(contractID, instructionsInterface, parameters)
    if err != nil {
        return "", fmt.Errorf("execution error: %v", err)
    }

    // Record the execution
    record := SmartContractExecutionRecord{
        ContractID: contract.ID,
        Timestamp:  time.Now(),
        Result:     fmt.Sprintf("%v", executionResult), // Convert execution result to string
    }
    sb.ExecutionHistory = append(sb.ExecutionHistory, record)

    fmt.Printf("Contract %s executed in sandbox %s.\n", contractID, sb.ID)
    return fmt.Sprintf("%v", executionResult), nil
}



// PauseSandbox pauses all activities in the sandbox, preventing further executions.
func (sb *Sandbox) PauseSandbox() {
    sb.mutex.Lock()
    defer sb.mutex.Unlock()

    sb.IsActive = false
    fmt.Printf("Sandbox %s paused.\n", sb.ID)
}

// ResumeSandbox resumes the paused sandbox, allowing contract activities to continue.
func (sb *Sandbox) ResumeSandbox() {
    sb.mutex.Lock()
    defer sb.mutex.Unlock()

    sb.IsActive = true
    fmt.Printf("Sandbox %s resumed.\n", sb.ID)
}

// ResetSandbox clears the loaded contracts and execution history in the sandbox.
func (sb *Sandbox) ResetSandbox() {
    sb.mutex.Lock()
    defer sb.mutex.Unlock()

    sb.SmartContracts = make(map[string]*SmartContract)
    sb.ExecutionHistory = []SmartContractExecutionRecord{}
    fmt.Printf("Sandbox %s reset.\n", sb.ID)
}

// TerminateSandbox shuts down the sandbox environment.
func (sb *Sandbox) TerminateSandbox() {
    sb.mutex.Lock()
    defer sb.mutex.Unlock()

    sb.IsActive = false
    sb.SmartContracts = nil
    sb.ExecutionHistory = nil
    fmt.Printf("Sandbox %s terminated.\n", sb.ID)
}

// generateSandboxID generates a unique ID for each sandbox instance.
func generateSandboxID() string {
    return fmt.Sprintf("sandbox-%d", time.Now().UnixNano())
}
