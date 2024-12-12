package math_and_logical

// IsTrue checks if a condition is true
func IsTrue(condition bool) bool {
	return condition
}

// IsFalse checks if a condition is false
func IsFalse(condition bool) bool {
	return !condition
}

// LogicalIf evaluates a condition and returns one of two values depending on whether the condition is true
func LogicalIf(condition bool, trueValue, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

// LogicalElse is designed to be used with LogicalIf to evaluate an alternative condition
// In this setup, LogicalElse simply returns the value passed if no previous condition was met
func LogicalElse(falseValue interface{}) interface{} {
	return falseValue
}

// LogicalSwitch performs a switch-case operation by evaluating a set of cases and returning the value of the first matching case
func LogicalSwitch(input interface{}, cases map[interface{}]interface{}) (interface{}, bool) {
	for caseValue, result := range cases {
		if input == caseValue {
			return result, true
		}
	}
	return nil, false
}

// LogicalCase evaluates a case in a switch-case structure and returns a boolean indicating if the case matched
func LogicalCase(input, caseValue interface{}) bool {
	return input == caseValue
}

// LogicalBreak is used to signal the end of a case in a switch-case structure
// In this setup, LogicalBreak returns a signal boolean to indicate a break
func LogicalBreak() bool {
	return true // Acts as a placeholder to represent a break in switch-case logic
}

// LogicalContinue is used to signal a continuation to the next iteration in a loop
// In this setup, LogicalContinue returns a signal boolean to indicate continuation
func LogicalContinue() bool {
	return true // Acts as a placeholder to represent a continue in loop logic
}

// LogicalEndIf is used to signal the end of a conditional structure
// In this setup, LogicalEndIf returns a signal boolean to indicate the end of an if structure
func LogicalEndIf() bool {
	return true // Placeholder for conditional logic structure closure
}
