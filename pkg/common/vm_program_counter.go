package common


import (
	"errors"
	"fmt"
)

// ProgramCounter represents the program counter in the VM, controlling bytecode execution flow.
type ProgramCounter struct {
	PC         int // Current position in the bytecode
	Bytecode   []byte // The bytecode being executed
	BytecodeLen int   // Length of the bytecode for bounds checking
}

// NewProgramCounter initializes a ProgramCounter with the given bytecode.
func NewProgramCounter(bytecode []byte) *ProgramCounter {
	return &ProgramCounter{
		PC:         0,
		Bytecode:   bytecode,
		BytecodeLen: len(bytecode),
	}
}

// Set sets the program counter to a specific address.
func (pc *ProgramCounter) Set(address int) error {
	if address < 0 || address >= pc.BytecodeLen {
		return fmt.Errorf("program counter set out of bounds: %d", address)
	}
	pc.PC = address
	return nil
}

// Increment advances the program counter by a specified number of bytes.
func (pc *ProgramCounter) Increment(steps int) error {
	if pc.PC+steps >= pc.BytecodeLen {
		return errors.New("program counter increment out of bounds")
	}
	pc.PC += steps
	return nil
}

// Jump sets the program counter to a specific address if within bounds.
func (pc *ProgramCounter) Jump(address int) error {
	return pc.Set(address)
}

// ConditionalJump performs a jump to a specified address based on a condition.
func (pc *ProgramCounter) ConditionalJump(address int, condition bool) error {
	if condition {
		return pc.Jump(address)
	}
	return nil
}

// JumpIfZero performs a jump if the Zero condition flag is set.
func (pc *ProgramCounter) JumpIfZero(flags *ConditionFlags, address int) error {
	return pc.ConditionalJump(address, flags.Zero)
}

// JumpIfNotZero performs a jump if the Zero condition flag is not set.
func (pc *ProgramCounter) JumpIfNotZero(flags *ConditionFlags, address int) error {
	return pc.ConditionalJump(address, !flags.Zero)
}

// JumpIfNegative performs a jump if the Negative condition flag is set.
func (pc *ProgramCounter) JumpIfNegative(flags *ConditionFlags, address int) error {
	return pc.ConditionalJump(address, flags.Negative)
}

// JumpIfOverflow performs a jump if the Overflow condition flag is set.
func (pc *ProgramCounter) JumpIfOverflow(flags *ConditionFlags, address int) error {
	return pc.ConditionalJump(address, flags.Overflow)
}

// JumpIfCarry performs a jump if the Carry condition flag is set.
func (pc *ProgramCounter) JumpIfCarry(flags *ConditionFlags, address int) error {
	return pc.ConditionalJump(address, flags.Carry)
}

// Reset resets the program counter to the beginning of the bytecode.
func (pc *ProgramCounter) Reset() {
	pc.PC = 0
}

// Fetch retrieves the next byte in the bytecode and increments the program counter.
func (pc *ProgramCounter) Fetch() (byte, error) {
	if pc.PC >= pc.BytecodeLen {
		return 0, errors.New("program counter fetch out of bounds")
	}
	byteCode := pc.Bytecode[pc.PC]
	pc.PC++
	return byteCode, nil
}
