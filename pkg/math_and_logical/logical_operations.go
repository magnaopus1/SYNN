package math_and_logical

// LogicalAnd performs a logical AND operation on two boolean values
func LogicalAnd(a, b bool) bool {
	return a && b
}

// LogicalOr performs a logical OR operation on two boolean values
func LogicalOr(a, b bool) bool {
	return a || b
}

// LogicalXor performs a logical XOR (exclusive OR) operation on two boolean values
func LogicalXor(a, b bool) bool {
	return a != b
}

// LogicalNot performs a logical NOT operation, inverting a single boolean value
func LogicalNot(a bool) bool {
	return !a
}

// LogicalNand performs a logical NAND operation on two boolean values (NOT AND)
func LogicalNand(a, b bool) bool {
	return !(a && b)
}

// LogicalNor performs a logical NOR operation on two boolean values (NOT OR)
func LogicalNor(a, b bool) bool {
	return !(a || b)
}

// LogicalXnor performs a logical XNOR operation on two boolean values (NOT XOR)
func LogicalXnor(a, b bool) bool {
	return a == b
}
