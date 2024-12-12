package math_and_logical

import (
	"errors"
)

// AssertTrue checks if a condition is true, returning an error if it’s false
func AssertTrue(condition bool, message string) error {
	if !condition {
		return errors.New(message)
	}
	return nil
}

// AssertFalse checks if a condition is false, returning an error if it’s true
func AssertFalse(condition bool, message string) error {
	if condition {
		return errors.New(message)
	}
	return nil
}

// LogicalTernary implements a ternary operation: returns `trueValue` if `condition` is true, else `falseValue`
func LogicalTernary(condition bool, trueValue, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

// LogicalImply returns the result of the logical implication `if a then b`
// In implication, the statement is false only if a is true and b is false
func LogicalImply(a, b bool) bool {
	return !a || b
}

// LogicalBiconditional returns true if both conditions are either true or false
func LogicalBiconditional(a, b bool) bool {
	return a == b
}

// LogicalAllTrue checks if all provided boolean values are true
func LogicalAllTrue(conditions ...bool) bool {
	for _, condition := range conditions {
		if !condition {
			return false
		}
	}
	return true
}

// LogicalAnyTrue checks if any of the provided boolean values is true
func LogicalAnyTrue(conditions ...bool) bool {
	for _, condition := range conditions {
		if condition {
			return true
		}
	}
	return false
}

// LogicalInvert inverts the provided boolean value
func LogicalInvert(value bool) bool {
	return !value
}

// SwitchLogicalState toggles a boolean value from true to false or vice versa
func SwitchLogicalState(value *bool) {
	*value = !*value
}
