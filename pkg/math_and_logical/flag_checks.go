package math_and_logical

import (
	"errors"
)

// CheckZeroFlag checks if a value is zero, returning true if it is zero
func CheckZeroFlag(value int) bool {
	return value == 0
}

// CheckNegativeFlag checks if a value is negative, returning true if it is negative
func CheckNegativeFlag(value int) bool {
	return value < 0
}

// CheckOverflowFlag checks for overflow in arithmetic operations
// It returns true if an overflow has occurred for int32 operations
func CheckOverflowFlag(a, b int32, operation string) (bool, error) {
	switch operation {
	case "add":
		if (b > 0 && a > (1<<31-1)-b) || (b < 0 && a < (-1<<31)-b) {
			return true, nil
		}
	case "subtract":
		if (b > 0 && a < (-1<<31)+b) || (b < 0 && a > (1<<31-1)+b) {
			return true, nil
		}
	case "multiply":
		if a != 0 && ((a*b)/a != b || (a*b)/b != a) {
			return true, nil
		}
	default:
		return false, errors.New("unsupported operation for overflow check")
	}
	return false, nil
}

// CheckCarryFlag checks for a carry in an addition operation
// It returns true if a carry has occurred for unsigned 32-bit integers
func CheckCarryFlag(a, b uint32) bool {
	return a > (1<<32-1)-b
}

// CheckParityFlag checks if the number of set bits (1s) in a value is even (returns true) or odd (returns false)
func CheckParityFlag(value int) bool {
	count := 0
	for value != 0 {
		count += value & 1
		value >>= 1
	}
	return count%2 == 0
}

// CheckSignFlag checks the sign of a value, returning true if the value is negative (sign flag set)
func CheckSignFlag(value int) bool {
	return value < 0
}

// CheckLogicalIntegrity verifies the logical integrity of a value
// For blockchain applications, this could involve ensuring that the value is within expected bounds
func CheckLogicalIntegrity(value int, min int, max int) bool {
	return value >= min && value <= max
}
