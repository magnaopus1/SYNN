package common

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)

// VMInterface defines the general interface for Virtual Machine functionalities.
type VMInterface interface {
	ExecuteContract(contractID, contractSource, language string, parameters map[string]interface{}, encryptionKey []byte) (map[string]interface{}, error)
	Compile(contractID, contractCode, language string, parameters map[string]interface{}) (string, error)
	ExecuteTransaction(transaction *Transaction) error
	CreateSubBlock(subBlock *SubBlock) error
	FinalizeBlock(block *Block) error
	LoadContract(code string, parameters map[string]interface{}) error
	ExecuteInstruction(instr Instruction) error // Added this method to the interface
	ValidateContractSyntax(contractID, contractSource, language string) (bool, error) // New method
	ParseBytecode(bytecode string) ([]string, error)
    ExecuteInstructions(contractID string, instructions []interface{}, parameters map[string]interface{}) (map[string]interface{}, error)
}

// NetworkNode represents a node in the network.
type NetworkNode interface {
	GetAddress() string
	IsNodeActive() bool
}

// Define logging levels, typically in a separate constants or logging package
const (
	LogLevelDebug = iota // 0 - Debug level for detailed internal information
	LogLevelInfo         // 1 - Info level for general information
	LogLevelWarning      // 2 - Warning level for potential issues
	LogLevelError        // 3 - Error level for serious issues
)



// VMFactory dynamically creates either a LightVM or HeavyVM based on task complexity.
func VMFactory(taskComplexity int, ledgerInstance *ledger.Ledger, consensus *SynnergyConsensus, loggingEnabled bool) (VMInterface, error) {
	// Define complexity threshold (e.g., tasks with complexity <= 5 use LightVM)
	if taskComplexity <= 5 {
		return NewLightVM(ledgerInstance, consensus, loggingEnabled)
	}
	return NewHeavyVM(ledgerInstance, consensus, loggingEnabled)
}

// LightVM represents the streamlined VM for low-resource tasks.
type LightVM struct {
	Contracts                 map[string]*SmartContract
	BytecodeInterpreter       *BytecodeInterpreter
	BytecodeGenerator         *BytecodeGenerator
	CodeQualityAssurance      *CodeQualityAssurance
	CompilationDebugger       *CompilationDebugger
	SandboxManager            *SandboxManager
	SyntaxChecker             *SyntaxChecker
	GasManager                *GasManager
	Registers                 *RegisterBank
	Memory                    *VMMemory
	ProgramCounter            int
	Stack                     *VMStack
	LoggingEnabled            bool
	EncryptionModule          *Encryption
	EventQueue                *EventQueue
	Scheduler                 *EventScheduler
	mutex                     sync.Mutex
	SolidityCompiler          *SolidityCompiler
	RustCompiler              *RustCompiler
	YulCompiler               *YulCompiler
	JavascriptCompiler        *JavaScriptCompiler
	GolangCompiler            *GoContractCompiler
	SoliditySupport           *SoliditySupport
	RustSupport               *RustSupport
	YulSupport                *YulSupport
	GolangSupport             *GoSupport
	JavascriptSupport         *JavaScriptSupport
	InstructionFilter         *InstructionFilter
	InstructionSet            *InstructionSet
	BottleneckSharder         *BottleneckSharder
	VirtualMachineRouter      *VMRouter
	VirtualMachineConcurrency *VirtualMachineConcurrency
	PersistentState           *PersistentStateManager
	Logs                      *LogManager
	Cache                     *Cache
	SpeedOptimizer            *SpeedOptimizer
	ProcessingPowerOptimizer  *ProcessingPowerOptimizer
	State                     *VMState
}

// NewLightVM initializes a new LightVM instance.
func NewLightVM(ledgerInstance *ledger.Ledger, consensus *SynnergyConsensus, loggingEnabled bool) (*LightVM, error) {
	// Initialize Encryption Module
	encryptionModule, err := NewEncryption(256) // Assuming 256 as key size
	if err != nil {
		return nil, fmt.Errorf("failed to initialize encryption module: %v", err)
	}

	// Initialize LogManager
	logManager, err := NewLogManager(LogLevelInfo, true, "logs/vm.log", 10*1024*1024, true)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize LogManager: %v", err)
	}

	// Use the singleton PersistentStateManager
	persistentStateManager, err := GetPersistentStateManager(ledgerInstance, "state_file_path", logManager)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PersistentStateManager: %v", err)
	}

	// Create an array of NetworkNode-compatible nodes
	node := &Node{}
	nodes := []NetworkNode{node}

	// Initialize VMRouter
	vmRouter := NewVMRouter(node, nodes, nil, nil, 100, 200)

	// Initialize BottleneckSharder
	bottleneckSharder := NewBottleneckSharder(2, 4, 6, 8, 10, 12, ledgerInstance, consensus)

	return &LightVM{
		Contracts:                 make(map[string]*SmartContract),
		GasManager:                NewGasManager(ledgerInstance, consensus, 0.01),
		BytecodeInterpreter:       NewBytecodeInterpreter(ledgerInstance),
		BytecodeGenerator:         NewBytecodeGenerator(ledgerInstance),
		CodeQualityAssurance:      NewCodeQualityAssurance(ledgerInstance),
		CompilationDebugger:       NewCompilationDebugger(ledgerInstance),
		SandboxManager:            NewSandboxManager(ledgerInstance),
		SyntaxChecker:             NewSyntaxChecker(ledgerInstance),
		Registers:                 NewRegisterBank("R32"),
		Memory:                    NewVMMemory(512*1024, nil),
		ProgramCounter:            0,
		Stack:                     NewVMStack(256),
		LoggingEnabled:            loggingEnabled,
		EncryptionModule:          encryptionModule,
		EventQueue:                NewEventQueue(nil),
		Scheduler:                 NewEventScheduler(nil, time.Second),
		mutex:                     sync.Mutex{},
		SolidityCompiler:          NewSolidityCompiler(ledgerInstance),
		RustCompiler:              NewRustCompiler(ledgerInstance),
		YulCompiler:               NewYulCompiler(ledgerInstance),
		JavascriptCompiler:        NewJavaScriptCompiler(ledgerInstance),
		GolangCompiler:            NewGoContractCompiler(ledgerInstance),
		SoliditySupport:           NewSoliditySupport(ledgerInstance),
		RustSupport:               NewRustSupport(ledgerInstance),
		YulSupport:                NewYulSupport(ledgerInstance),
		GolangSupport:             NewGoSupport(ledgerInstance),
		JavascriptSupport:         NewJavaScriptSupport(ledgerInstance),
		InstructionFilter:         NewInstructionFilter(10, 20, 50, vmRouter),
		InstructionSet:            NewOpcodeInstructionSet(),
		BottleneckSharder:         bottleneckSharder,
		VirtualMachineRouter:      vmRouter,
		VirtualMachineConcurrency: NewVirtualMachineConcurrency(ledgerInstance, nil, nil, consensus),
		PersistentState:           persistentStateManager,
		Logs:                      logManager,
		Cache:                     NewCache(1000, time.Minute),
		State:                     persistentStateManager.state,
		SpeedOptimizer:            NewSpeedOptimizer(100, time.Second),
		ProcessingPowerOptimizer:  NewProcessingPowerOptimizer(5, 20, 0.75, time.Second*5),
	}, nil
}



func (vm *LightVM) ExecuteContract(contractID, contractSource, language string, parameters map[string]interface{}, encryptionKey []byte) (map[string]interface{}, error) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	// Simplified execution for light tasks
	isValid, err := vm.SyntaxChecker.ValidateContractSyntax(contractID, contractSource, language)
	if err != nil || !isValid {
		return nil, fmt.Errorf("contract syntax error: %v", err)
	}

	// Compile the contract
	compiledCode, err := vm.Compile(contractID, contractSource, language, parameters)
	if err != nil {
		return nil, fmt.Errorf("contract compilation error: %v", err)
	}

	// Execute the contract
	result := map[string]interface{}{"status": "success", "output": compiledCode}
	return result, nil
}

func (vm *LightVM) Compile(contractID, contractCode, language string, parameters map[string]interface{}) (string, error) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	switch language {
	case "solidity":
		err := vm.SolidityCompiler.CompileSolidityContract(contractID, contractCode)
		if err != nil {
			return "", err
		}
		return "Compiled Solidity Contract", nil

	// Include other languages as needed

	default:
		return "", fmt.Errorf("unsupported contract language: %s", language)
	}
}

func (vm *LightVM) ExecuteTransaction(transaction *Transaction) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	// Basic transaction execution
	isValid, err := vm.SyntaxChecker.ValidateTransactionSyntax(transaction)
	if err != nil || !isValid {
		return fmt.Errorf("transaction syntax error: %v", err)
	}

	// Charge gas
	err = vm.GasManager.ChargeGas(transaction.TransactionID, transaction.FromAddress, "validatorID123", time.Second*1)
	if err != nil {
		return fmt.Errorf("failed to charge gas: %v", err)
	}

	// Process transaction (simplified)
	return nil
}

// ValidateContractSyntax validates the syntax of a contract in LightVM.
func (vm *LightVM) ValidateContractSyntax(contractID, contractSource, language string) (bool, error) {
	if vm.SyntaxChecker == nil {
		return false, fmt.Errorf("syntax checker not initialized")
	}
	return vm.SyntaxChecker.ValidateContractSyntax(contractID, contractSource, language)
}

// ExecuteInstruction processes a single instruction in the LightVM.
func (vm *LightVM) ExecuteInstruction(instr Instruction) error {
	// Lock at the LightVM level to ensure thread-safe access to vm.State
	vm.mutex.Lock()  
	defer vm.mutex.Unlock()

	// Update program counter in VMState
	vm.State.ProgramCounter++

	// Convert Opcode from string to Opcode type
	opcode, err := stringToOpcode(instr.Opcode)
	if err != nil {
		return fmt.Errorf("invalid opcode: %v", err)
	}

	// Use Payload as operands, assuming Payload is intended to store operand data
	operands, ok := instr.Payload.([]interface{}) // Ensure Payload contains operands as a slice of interfaces
	if !ok {
		return fmt.Errorf("invalid operand format in instruction payload")
	}

	// Execute the instruction using InstructionSet
	result, err := vm.InstructionSet.ExecuteOpcodeWithOperands("TransactionInstructionSet", opcode, operands...)
	if err != nil {
		return fmt.Errorf("failed to execute instruction %s: %v", instr.Opcode, err)
	}

	fmt.Printf("LightVM executed instruction: %v at PC: %d\n", result, vm.State.ProgramCounter)
	return nil
}



// stringToOpcode converts a string to an Opcode, mapping known strings to Opcode constants.
func stringToOpcode(opcodeStr string) (Opcode, error) {
	switch opcodeStr {
	case "some_opcode_name": // Example case, replace with real opcode mappings
		return Opcode(0x01), nil
	// Add cases for other known opcodes
	default:
		return Opcode(0), fmt.Errorf("unknown opcode: %s", opcodeStr)
	}
}




func (vm *LightVM) CreateSubBlock(subBlock *SubBlock) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	// Simplified sub-block creation
	return nil
}

func (vm *LightVM) FinalizeBlock(block *Block) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	// Simplified block finalization
	return nil
}

func (vm *LightVM) LoadContract(code string, parameters map[string]interface{}) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	contractID := generateContractID("owner", code)
	contract := &SmartContract{
		ID:         contractID,
		Code:       code,
		Parameters: parameters,
		State:      make(map[string]interface{}),
	}

	vm.Contracts[contractID] = contract
	if vm.LoggingEnabled {
		fmt.Printf("Smart Contract %s loaded into LightVM.\n", contractID)
	}

	return nil
}

func (vm *LightVM) ParseBytecode(bytecode string) ([]string, error) {
    return vm.BytecodeInterpreter.parseBytecode(bytecode)
}

func (vm *LightVM) ExecuteInstructions(contractID string, instructions []interface{}, parameters map[string]interface{}) (map[string]interface{}, error) {
    return vm.BytecodeInterpreter.executeInstructions(contractID, instructions, parameters)
}


// HeavyVM represents the VM optimized for complex, high-resource tasks.
type HeavyVM struct {
	LedgerInstance            *ledger.Ledger
	Contracts                 map[string]*SmartContract
	BytecodeInterpreter       *BytecodeInterpreter
	BytecodeGenerator         *BytecodeGenerator
	CodeQualityAssurance      *CodeQualityAssurance
	CompilationDebugger       *CompilationDebugger
	SandboxManager            *SandboxManager
	SyntaxChecker             *SyntaxChecker
	GasManager                *GasManager
	SubBlockManager           *SubBlockManager
	Registers                 *RegisterBank
	Memory                    *VMMemory
	ProgramCounter            int
	Stack                     *VMStack
	ConditionFlags            *ConditionFlags
	EventQueue                *EventQueue
	Scheduler                 *EventScheduler
	LoggingEnabled            bool
	EncryptionModule          *Encryption
	mutex                     sync.Mutex
	SolidityCompiler          *SolidityCompiler
	RustCompiler              *RustCompiler
	YulCompiler               *YulCompiler
	JavascriptCompiler        *JavaScriptCompiler
	GolangCompiler            *GoContractCompiler
	SoliditySupport           *SoliditySupport
	RustSupport               *RustSupport
	YulSupport                *YulSupport
	GolangSupport             *GoSupport
	JavascriptSupport         *JavaScriptSupport
	VirtualMachineRouter      *VMRouter
	VirtualMachineConcurrency *VirtualMachineConcurrency
	InstructionFilter         *InstructionFilter
	InstructionSet            *InstructionSet
	BottleneckSharder         *BottleneckSharder
	Oracles                   *OracleManager
	CrossChain                *CrossChainHandler
	PersistentState           *PersistentStateManager
	Logs                      *LogManager
	Cache                     *Cache
	SpeedOptimizer            *SpeedOptimizer
	ProcessingPowerOptimizer  *ProcessingPowerOptimizer
	State                     *VMState
}

// NewHeavyVM initializes a new HeavyVM instance.
func NewHeavyVM(ledgerInstance *ledger.Ledger, consensus *SynnergyConsensus, loggingEnabled bool) (*HeavyVM, error) {
	// Initialize Encryption Module with a key size of 256 bits
	encryptionModule, err := NewEncryption(256)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize encryption module: %v", err)
	}

	// Initialize LogManager with correct parameters
	logManager, err := NewLogManager(LogLevelInfo, true, "logs/vm.log", 10*1024*1024, true)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize LogManager: %v", err)
	}

	// Use the singleton PersistentStateManager
	persistentStateManager, err := GetPersistentStateManager(ledgerInstance, "state_file_path", logManager)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PersistentStateManager: %v", err)
	}

	// Initialize Cache with a size of 1000 and 1-minute TTL
	cache := NewCache(1000, time.Minute)

	// Create VM Router with a placeholder *Node and VM instances
	node := &Node{} // Create a pointer to a Node instance
	vmRouter := NewVMRouter(node, []NetworkNode{node}, nil, nil, 100, 200)

	// Initialize BottleneckSharder with the correct number of parameters
	bottleneckSharder := NewBottleneckSharder(2, 4, 6, 8, 10, 12, ledgerInstance, consensus)

	return &HeavyVM{
		LedgerInstance:            ledgerInstance,
		Contracts:                 make(map[string]*SmartContract),
		BytecodeInterpreter:       NewBytecodeInterpreter(ledgerInstance),
		BytecodeGenerator:         NewBytecodeGenerator(ledgerInstance),
		CodeQualityAssurance:      NewCodeQualityAssurance(ledgerInstance),
		CompilationDebugger:       NewCompilationDebugger(ledgerInstance),
		SandboxManager:            NewSandboxManager(ledgerInstance),
		SyntaxChecker:             NewSyntaxChecker(ledgerInstance),
		GasManager:                NewGasManager(ledgerInstance, consensus, 0.05),
		SubBlockManager:           NewSubBlockManager(ledgerInstance, encryptionModule),
		Registers:                 NewRegisterBank("R32"),
		Memory:                    NewVMMemory(1024*1024, nil),
		ProgramCounter:            0,
		Stack:                     NewVMStack(1024),
		ConditionFlags:            NewConditionFlags(),
		EventQueue:                NewEventQueue(ledgerInstance),
		Scheduler:                 NewEventScheduler(NewEventQueue(ledgerInstance), time.Second),
		LoggingEnabled:            loggingEnabled,
		EncryptionModule:          encryptionModule,
		mutex:                     sync.Mutex{},
		SolidityCompiler:          NewSolidityCompiler(ledgerInstance),
		RustCompiler:              NewRustCompiler(ledgerInstance),
		YulCompiler:               NewYulCompiler(ledgerInstance),
		JavascriptCompiler:        NewJavaScriptCompiler(ledgerInstance),
		GolangCompiler:            NewGoContractCompiler(ledgerInstance),
		SoliditySupport:           NewSoliditySupport(ledgerInstance),
		RustSupport:               NewRustSupport(ledgerInstance),
		YulSupport:                NewYulSupport(ledgerInstance),
		GolangSupport:             NewGoSupport(ledgerInstance),
		JavascriptSupport:         NewJavaScriptSupport(ledgerInstance),
		InstructionFilter:         NewInstructionFilter(10, 20, 50, vmRouter),
		InstructionSet:            NewOpcodeInstructionSet(),
		BottleneckSharder:         bottleneckSharder,
		VirtualMachineRouter:      vmRouter,
		VirtualMachineConcurrency: NewVirtualMachineConcurrency(ledgerInstance, nil, NewSubBlockManager(ledgerInstance, encryptionModule), nil),
		Oracles:                   NewOracleManager(time.Minute, time.Second*10, 3, 5, []string{"trustedOracle1", "trustedOracle2"}),
		CrossChain:                NewCrossChainHandler(time.Minute),
		PersistentState:           persistentStateManager,
		Logs:                      logManager,
		Cache:                     cache,
		State:                     persistentStateManager.state,
		SpeedOptimizer:            NewSpeedOptimizer(100, time.Second),
		ProcessingPowerOptimizer:  NewProcessingPowerOptimizer(5, 20, 0.75, time.Second*5),
	}, nil
}

// Implementing VMInterface for HeavyVM

func (vm *HeavyVM) ExecuteContract(contractID, contractSource, language string, parameters map[string]interface{}, encryptionKey []byte) (map[string]interface{}, error) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	// Step 1: Syntax check
	isValid, err := vm.SyntaxChecker.ValidateContractSyntax(contractID, contractSource, language)
	if err != nil || !isValid {
		return nil, fmt.Errorf("contract syntax error: %v", err)
	}

	// Step 2: Compile contract
	compiledBytecode, err := vm.CompilationDebugger.CompileBytecode(contractID, contractSource, vm.EncryptionModule, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("contract compilation error: %v", err)
	}

	// Step 3: Code quality checks
	_, err = vm.CodeQualityAssurance.ValidateBytecodeQuality(contractID, compiledBytecode, vm.EncryptionModule, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("bytecode quality validation failed: %v", err)
	}

	// Step 4: Execute bytecode in sandbox
	executionResult, err := vm.BytecodeInterpreter.ExecuteBytecode(contractID, compiledBytecode, parameters, vm.EncryptionModule, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("contract execution failed: %v", err)
	}

	// Step 5: Charge gas fees
	err = vm.GasManager.ChargeGas(contractID, "txID123", "validatorID123", time.Second*5)
	if err != nil {
		return nil, fmt.Errorf("failed to charge gas fees: %v", err)
	}

	// Step 6: Log the execution into the ledger
	err = vm.LedgerInstance.VirtualMachineLedger.LogExecution(contractID, "bytecodeQualityResult", "validatorID123")
	if err != nil {
		return nil, fmt.Errorf("failed to log execution: %v", err)
	}

	return executionResult, nil
}

func (vm *HeavyVM) Compile(contractID, contractCode, language string, parameters map[string]interface{}) (string, error) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	switch language {
	case "solidity":
		err := vm.SolidityCompiler.CompileSolidityContract(contractID, contractCode)
		if err != nil {
			return "", err
		}
		return vm.SoliditySupport.DeployAndExecuteContract(contractID, contractCode, parameters)

	// Include other languages as needed

	default:
		return "", fmt.Errorf("unsupported contract language: %s", language)
	}
}

func (vm *HeavyVM) ExecuteTransaction(transaction *Transaction) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	// Step 1: Syntax validation
	isValid, err := vm.SyntaxChecker.ValidateTransactionSyntax(transaction)
	if err != nil || !isValid {
		return fmt.Errorf("transaction syntax error: %v", err)
	}

	// Step 2: Charge gas
	err = vm.GasManager.ChargeGas(transaction.TransactionID, transaction.FromAddress, "validatorID123", time.Second*5)
	if err != nil {
		return fmt.Errorf("failed to charge gas: %v", err)
	}

	// Step 3: Add transaction to sub-block
	err = vm.SubBlockManager.AddTransactionToSubBlock(transaction)
	if err != nil {
		return fmt.Errorf("failed to add transaction to sub-block: %v", err)
	}

	// Step 4: Log transaction
	err = vm.LedgerInstance.BlockchainConsensusCoinLedger.LogTransaction(transaction.TransactionID, transaction.Signature)
	if err != nil {
		return fmt.Errorf("failed to log transaction: %v", err)
	}

	return nil
}

// ValidateContractSyntax validates the syntax of a contract in HeavyVM.
func (vm *HeavyVM) ValidateContractSyntax(contractID, contractSource, language string) (bool, error) {
	if vm.SyntaxChecker == nil {
		return false, fmt.Errorf("syntax checker not initialized")
	}
	return vm.SyntaxChecker.ValidateContractSyntax(contractID, contractSource, language)
}

func (vm *HeavyVM) ParseBytecode(bytecode string) ([]string, error) {
    return vm.BytecodeInterpreter.parseBytecode(bytecode)
}

func (vm *HeavyVM) ExecuteInstructions(contractID string, instructions []interface{}, parameters map[string]interface{}) (map[string]interface{}, error) {
    return vm.BytecodeInterpreter.executeInstructions(contractID, instructions, parameters)
}

// ExecuteInstruction processes a single instruction in the HeavyVM.
func (vm *HeavyVM) ExecuteInstruction(instr Instruction) error {
	// Use HeavyVM's mutex to synchronize access to vm.State
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	// Update program counter in VMState
	vm.State.ProgramCounter++

	// Convert Opcode from string to Opcode type
	opcode, err := stringToOpcode(instr.Opcode)
	if err != nil {
		return fmt.Errorf("invalid opcode: %v", err)
	}

	// Use Payload as operands, assuming Payload is intended to store operand data
	operands, ok := instr.Payload.([]interface{}) // Ensure Payload contains operands as a slice of interfaces
	if !ok {
		return fmt.Errorf("invalid operand format in instruction payload")
	}

	// Execute the instruction using InstructionSet
	result, err := vm.InstructionSet.ExecuteOpcodeWithOperands("AutomationInstructionSet", opcode, operands...)
	if err != nil {
		return fmt.Errorf("failed to execute instruction %s: %v", instr.Opcode, err)
	}

	fmt.Printf("HeavyVM executed instruction: %v at PC: %d\n", result, vm.State.ProgramCounter)
	return nil
}




func (vm *HeavyVM) CreateSubBlock(subBlock *SubBlock) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	// Validate all transactions
	for _, tx := range subBlock.Transactions {
		err := vm.ExecuteTransaction(&tx)
		if err != nil {
			return fmt.Errorf("sub-block validation failed: %v", err)
		}
	}

	// Convert SubBlock to ledger.SubBlock by dereferencing
	ledgerSubBlock := ConvertToLedgerSubBlock(*subBlock)

	// Store sub-block in ledger by passing a pointer to ledgerSubBlock
	err := vm.LedgerInstance.BlockchainConsensusCoinLedger.LogSubBlock(&ledgerSubBlock)
	if err != nil {
		return fmt.Errorf("failed to log sub-block: %v", err)
	}

	return nil
}



func (vm *HeavyVM) FinalizeBlock(block *Block) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	if len(block.SubBlocks) < 1000 {
		return errors.New("block must contain 1000 sub-blocks")
	}

	// Log the finalized block
	err := vm.LedgerInstance.BlockchainConsensusCoinLedger.LogBlock(fmt.Sprintf("Block %d", block.Index), block.Hash)
	if err != nil {
		return fmt.Errorf("failed to log finalized block: %v", err)
	}

	return nil
}

func (vm *HeavyVM) LoadContract(code string, parameters map[string]interface{}) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	contractID := generateContractID("owner", code)
	contract := &SmartContract{
		ID:         contractID,
		Code:       code,
		Parameters: parameters,
		State:      make(map[string]interface{}),
	}

	vm.Contracts[contractID] = contract
	if vm.LoggingEnabled {
		fmt.Printf("Smart Contract %s loaded into HeavyVM.\n", contractID)
	}

	return nil
}

// VirtualMachine is a wrapper that holds a VMInterface.
// This allows other parts of the system to interact with the VM without knowing whether it's a LightVM or HeavyVM.
type VirtualMachine struct {
	vm VMInterface
}

// NewVirtualMachine creates a new VirtualMachine instance, selecting either LightVM or HeavyVM based on task complexity.
func NewVirtualMachine(taskComplexity int, ledgerInstance *ledger.Ledger, consensus *SynnergyConsensus, loggingEnabled bool) (*VirtualMachine, error) {
	vmInstance, err := VMFactory(taskComplexity, ledgerInstance, consensus, loggingEnabled)
	if err != nil {
		return nil, err
	}
	return &VirtualMachine{vm: vmInstance}, nil
}



func (vm *VirtualMachine) ExecuteContract(contractID, contractSource, language string, parameters map[string]interface{}, encryptionKey []byte) (map[string]interface{}, error) {
	return vm.vm.ExecuteContract(contractID, contractSource, language, parameters, encryptionKey)
}

func (vm *VirtualMachine) Compile(contractID, contractCode, language string, parameters map[string]interface{}) (string, error) {
	return vm.vm.Compile(contractID, contractCode, language, parameters)
}

func (vm *VirtualMachine) ExecuteTransaction(transaction *Transaction) error {
	return vm.vm.ExecuteTransaction(transaction)
}

func (vm *VirtualMachine) CreateSubBlock(subBlock *SubBlock) error {
	return vm.vm.CreateSubBlock(subBlock)
}

func (vm *VirtualMachine) FinalizeBlock(block *Block) error {
	return vm.vm.FinalizeBlock(block)
}

func (vm *VirtualMachine) LoadContract(code string, parameters map[string]interface{}) error {
	return vm.vm.LoadContract(code, parameters)
}




func (vm *VirtualMachine) ExecuteInstruction(instr Instruction) error {
	return vm.vm.ExecuteInstruction(instr)
}

func (vm *VirtualMachine) ValidateContractSyntax(contractID, contractSource, language string) (bool, error) {
	return vm.vm.ValidateContractSyntax(contractID, contractSource, language)
}

func (vm *VirtualMachine) ParseBytecode(bytecode string) ([]string, error) {
	return vm.vm.ParseBytecode(bytecode)
}

func (vm *VirtualMachine) ExecuteInstructions(contractID string, instructions []interface{}, parameters map[string]interface{}) (map[string]interface{}, error) {
	return vm.vm.ExecuteInstructions(contractID, instructions, parameters)
}
