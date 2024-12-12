package common

import (
	"encoding/binary"
	"errors"
	"sync"
)

// AddressingMode specifies the addressing mode for opcode execution
type AddressingMode int

const (
	Immediate AddressingMode = iota // Direct value, e.g., constants or literals
	Direct                          // Direct memory address
	Indirect                        // Address points to another address in memory
	Indexed                         // Uses a base address plus an index
)

// AddressResolver is responsible for interpreting and resolving different addressing modes.
type AddressResolver struct {
	memory    *VMMemory      // Virtual machine's memory instance
	registers *RegisterBank  // Register bank with VM registers
	mutex     sync.Mutex     // Mutex for thread-safe memory access
}

// RegisterBank is a struct for managing VM registers.
type RegisterBank struct {
	registers map[string]int
	mutex     sync.Mutex // Optional mutex for thread-safe operations
}

// NewRegisterBank initializes a new RegisterBank with predefined register names or an empty bank.
func NewRegisterBank(registerNames ...string) *RegisterBank {
	rb := &RegisterBank{registers: make(map[string]int)}
	// Initialize with given register names if any, else leave it as empty
	for _, name := range registerNames {
		rb.registers[name] = 0
	}
	return rb
}

// SetRegisterValue sets the value of a specified register.
func (rb *RegisterBank) SetRegisterValue(register string, value int) error {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()

	// Check if the register exists
	if _, exists := rb.registers[register]; !exists {
		return errors.New("register not found")
	}

	// Set the register value
	rb.registers[register] = value
	return nil
}

// GetRegisterValue retrieves the value of a specified register.
func (rb *RegisterBank) GetRegisterValue(register string) (int, error) {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()

	// Check if the register exists
	value, exists := rb.registers[register]
	if !exists {
		return 0, errors.New("register not found")
	}

	return value, nil
}

// NewAddressResolver initializes a new AddressResolver.
func NewAddressResolver(memory *VMMemory, registers *RegisterBank) *AddressResolver {
	return &AddressResolver{
		memory:    memory,
		registers: registers,
	}
}

// Resolve fetches the value or address based on the given addressing mode.
func (ar *AddressResolver) Resolve(mode AddressingMode, operand int) (int, error) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()

	switch mode {
	case Immediate:
		return operand, nil // Return the operand directly as it's an immediate value

	case Direct:
		data, err := ar.memory.Read(operand, 4) // Read 4 bytes (or size required) from memory
		if err != nil {
			return 0, err
		}
		return bytesToInt(data), nil

	case Indirect:
		data, err := ar.memory.Read(operand, 4) // Read address from memory
		if err != nil {
			return 0, err
		}
		address := bytesToInt(data)
		valueData, err := ar.memory.Read(address, 4) // Then read actual value
		if err != nil {
			return 0, err
		}
		return bytesToInt(valueData), nil

	case Indexed:
		baseAddress := operand
		index, err := ar.registers.GetRegisterValue("Index") // Fetch the index register for offset
		if err != nil {
			return 0, err
		}
		data, err := ar.memory.Read(baseAddress+index, 4) // Read from memory with the indexed offset
		if err != nil {
			return 0, err
		}
		return bytesToInt(data), nil

	default:
		return 0, errors.New("invalid addressing mode")
	}
}

// bytesToInt converts a byte slice to an integer.
func bytesToInt(data []byte) int {
	return int(binary.BigEndian.Uint32(data)) // Assuming 4 bytes for simplicity
}

// Get retrieves the value of the specified register.
func (rb *RegisterBank) Get(register string) (int, error) {
	value, exists := rb.registers[register]
	if !exists {
		return 0, errors.New("register not found")
	}
	return value, nil
}
