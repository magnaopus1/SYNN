package common


// These flags determine conditional jumps, loops, and other operations based on the results of prior instructions.
type ConditionFlags struct {
	Zero      bool // Set if the result of an operation is zero
	Overflow  bool // Set if an operation resulted in overflow
	Negative  bool // Set if the result is negative
	Carry     bool // Set if an operation produced a carry bit
}

// NewConditionFlags initializes the condition flags with default values.
func NewConditionFlags() *ConditionFlags {
	return &ConditionFlags{
		Zero:      false,
		Overflow:  false,
		Negative:  false,
		Carry:     false,
	}
}

// UpdateZero sets the Zero flag based on whether the result is zero.
func (cf *ConditionFlags) UpdateZero(result int) {
	cf.Zero = (result == 0)
}

// UpdateOverflow sets the Overflow flag if an overflow occurred in a signed operation.
func (cf *ConditionFlags) UpdateOverflow(result int, operand1, operand2 int) {
	// Check for overflow in a 32-bit signed integer
	cf.Overflow = (operand1 > 0 && operand2 > 0 && result < 0) || (operand1 < 0 && operand2 < 0 && result > 0)
}

// UpdateNegative sets the Negative flag based on whether the result is negative.
func (cf *ConditionFlags) UpdateNegative(result int) {
	cf.Negative = (result < 0)
}

// UpdateCarry sets the Carry flag for operations where a carry occurs in unsigned addition/subtraction.
func (cf *ConditionFlags) UpdateCarry(result uint, operand1, operand2 uint, operation string) {
	switch operation {
	case "add":
		cf.Carry = (result < operand1) || (result < operand2)
	case "subtract":
		cf.Carry = (operand1 < operand2)
	default:
		cf.Carry = false
	}
}

// Reset resets all flags to their default false state.
func (cf *ConditionFlags) Reset() {
	cf.Zero = false
	cf.Overflow = false
	cf.Negative = false
	cf.Carry = false
}

// VM Condition Flag Methods for Execution Control

// ShouldJumpIfZero returns true if the Zero flag is set, used for conditional jumps based on zero results.
func (cf *ConditionFlags) ShouldJumpIfZero() bool {
	return cf.Zero
}

// ShouldJumpIfNotZero returns true if the Zero flag is not set, used for non-zero conditional jumps.
func (cf *ConditionFlags) ShouldJumpIfNotZero() bool {
	return !cf.Zero
}

// ShouldJumpIfOverflow returns true if the Overflow flag is set, used for overflow-based conditional jumps.
func (cf *ConditionFlags) ShouldJumpIfOverflow() bool {
	return cf.Overflow
}

// ShouldJumpIfNegative returns true if the Negative flag is set, used for conditional jumps based on negative results.
func (cf *ConditionFlags) ShouldJumpIfNegative() bool {
	return cf.Negative
}

// ShouldJumpIfCarry returns true if the Carry flag is set, used for conditional jumps based on carry results.
func (cf *ConditionFlags) ShouldJumpIfCarry() bool {
	return cf.Carry
}

// Integration with the VM Execution Cycle

// SetFlagsForArithmetic updates all relevant flags based on the result and operands of an arithmetic operation.
// This function is typically called after executing an arithmetic opcode.
func (cf *ConditionFlags) SetFlagsForArithmetic(result int, operand1, operand2 int) {
	cf.UpdateZero(result)
	cf.UpdateOverflow(result, operand1, operand2)
	cf.UpdateNegative(result)
}

// SetFlagsForUnsignedArithmetic updates flags relevant to unsigned arithmetic operations (e.g., carry).
func (cf *ConditionFlags) SetFlagsForUnsignedArithmetic(result uint, operand1, operand2 uint, operation string) {
	cf.UpdateZero(int(result))
	cf.UpdateCarry(result, operand1, operand2, operation)
}
