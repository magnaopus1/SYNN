package common

import (
	"fmt"
	"sync"
)

// Define an opcode type
type Opcode int

// InstructionGroup represents a set of opcode instructions.
type InstructionGroup map[Opcode]func(params ...interface{}) (interface{}, error)

// Define the OpcodeInstructionSet struct to hold multiple instruction groups.
type InstructionSet struct {
	instructionGroups map[string]InstructionGroup // Map to hold all instruction sets
	mutex             sync.Mutex                   // Mutex for thread-safe operations
}

// NewOpcodeInstructionSet initializes an OpcodeInstructionSet.
func NewOpcodeInstructionSet() *InstructionSet {
	ois := &InstructionSet{
		instructionGroups: make(map[string]InstructionGroup),
	}
	ois.loadInstructionGroups() // Load predefined instruction sets
	return ois
}

// loadInstructionGroups initializes instruction sets within the opcode instruction set.
func (ois *InstructionSet) loadInstructionGroups() {
	ois.instructionGroups = make(map[string]InstructionGroup)

	// Initialize real instruction sets
	ois.instructionGroups["TransactionInstructionSet"] = make(InstructionGroup)
	ois.instructionGroups["AutomationInstructionSet"] = make(InstructionGroup)

	// Register the actual opcodes
	ois.RegisterOpcode("TransactionInstructionSet", 0x01, ois.handleLoad) // Example opcode for loading
	ois.RegisterOpcode("AutomationInstructionSet", 0x10, ois.handleAutomate) // Example opcode for automation
}

// RegisterOpcode registers a single opcode within a specified instruction group.
func (ois *InstructionSet) RegisterOpcode(groupName string, opcode Opcode, handler func(params ...interface{}) (interface{}, error)) error {
	ois.mutex.Lock()
	defer ois.mutex.Unlock()

	// Check if the instruction group exists
	group, exists := ois.instructionGroups[groupName]
	if !exists {
		// If not, create a new instruction group and add it to instructionGroups
		group = make(InstructionGroup)
		ois.instructionGroups[groupName] = group
	}

	// Check if the opcode already exists within the group
	if _, opcodeExists := group[opcode]; opcodeExists {
		return fmt.Errorf("opcode %d already exists in instruction group %s", opcode, groupName)
	}

	// Register the opcode and its handler function within the instruction group
	group[opcode] = handler
	return nil
}

// ExecuteOpcodeWithOperands is the dispatcher for calling opcodes from any group.
func (ois *InstructionSet) ExecuteOpcodeWithOperands(groupName string, opcode Opcode, operands ...interface{}) (interface{}, error) {
	ois.mutex.Lock()
	defer ois.mutex.Unlock()

	group, exists := ois.instructionGroups[groupName]
	if !exists {
		return nil, fmt.Errorf("instruction group %s does not exist", groupName)
	}

	function, opcodeExists := group[opcode]
	if !opcodeExists {
		return nil, fmt.Errorf("opcode %d is not defined in the instruction group %s", opcode, groupName)
	}

	return function(operands...) // Execute the function with provided operands
}

// Example opcode handler for loading a value into a register.
func (ois *InstructionSet) handleLoad(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("LOAD expects 2 operands")
	}
	registerName, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("first operand must be a register name")
	}
	value, ok := params[1].(int)
	if !ok {
		return nil, fmt.Errorf("second operand must be an integer")
	}
	// Here you would set the value in the register (to be implemented)
	fmt.Printf("Loaded %d into %s\n", value, registerName)
	return nil, nil
}

// Example opcode handler for an automation task.
func (ois *InstructionSet) handleAutomate(params ...interface{}) (interface{}, error) {
	// Implementation of automation logic goes here
	fmt.Println("Executed automation task")
	return nil, nil
}
