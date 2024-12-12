package utility

// Equal checks if two integers are equal.
func Equal(a, b int) bool {
	return a == b
}

// NotEqual checks if two integers are not equal.
func NotEqual(a, b int) bool {
	return a != b
}

// LessThan checks if the first integer is less than the second.
func LessThan(a, b int) bool {
	return a < b
}

// GreaterThan checks if the first integer is greater than the second.
func GreaterThan(a, b int) bool {
	return a > b
}

// LessThanOrEqual checks if the first integer is less than or equal to the second.
func LessThanOrEqual(a, b int) bool {
	return a <= b
}

// GreaterThanOrEqual checks if the first integer is greater than or equal to the second.
func GreaterThanOrEqual(a, b int) bool {
	return a >= b
}

// ZeroCheck returns true if the provided integer is zero.
func ZeroCheck(a int) bool {
	return a == 0
}

// IsPositive checks if the provided integer is positive.
func IsPositive(a int) bool {
	return a > 0
}

// IsNegative checks if the provided integer is negative.
func IsNegative(a int) bool {
	return a < 0
}

// IsOdd returns true if the provided integer is odd.
func IsOdd(a int) bool {
	return a%2 != 0
}

// IsEven returns true if the provided integer is even.
func IsEven(a int) bool {
	return a%2 == 0
}

// Between checks if a value is strictly between two bounds (exclusive).
func Between(value, lower, upper int) bool {
	return value > lower && value < upper
}

// NotBetween checks if a value is not strictly between two bounds (exclusive).
func NotBetween(value, lower, upper int) bool {
	return value < lower || value > upper
}

// InRange checks if a value is within a range, inclusive of both bounds.
func InRange(value, lower, upper int) bool {
	return value >= lower && value <= upper
}

// OutOfRange checks if a value is outside a specified range, inclusive of both bounds.
func OutOfRange(value, lower, upper int) bool {
	return value < lower || value > upper
}

// IsNull checks if a given pointer is nil.
func IsNull(ptr interface{}) bool {
	return ptr == nil
}
