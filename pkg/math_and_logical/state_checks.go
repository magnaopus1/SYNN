package math_and_logical

// IsPositive checks if a given integer or float value is positive
func IsPositive(value float64) bool {
	return value > 0
}

// IsNegative checks if a given integer or float value is negative
func IsNegative(value float64) bool {
	return value < 0
}

// IsOdd checks if an integer value is odd
func IsOdd(value int) bool {
	return value%2 != 0
}

// IsEven checks if an integer value is even
func IsEven(value int) bool {
	return value%2 == 0
}

// IsNull checks if a given pointer or interface is nil (null)
func IsNull(value interface{}) bool {
	return value == nil
}

// IsNotNull checks if a given pointer or interface is not nil (non-null)
func IsNotNull(value interface{}) bool {
	return value != nil
}

// IsNonZero checks if a given float or integer value is non-zero
func IsNonZero(value float64) bool {
	return value != 0
}

// IsEmpty checks if a slice, string, or map is empty
func IsEmpty(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return len(v) == 0
	case []interface{}:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case map[interface{}]interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		return false
	}
}
