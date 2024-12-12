package utility

// BITWISE_AND performs a bitwise AND operation between two values.
func BITWISE_AND(value1, value2 uint32) uint32 {
	return value1 & value2
}

// BITWISE_OR performs a bitwise OR operation between two values.
func BITWISE_OR(value1, value2 uint32) uint32 {
	return value1 | value2
}

// BITWISE_XOR performs a bitwise XOR operation between two values.
func BITWISE_XOR(value1, value2 uint32) uint32 {
	return value1 ^ value2
}

// BITWISE_NOT performs a bitwise NOT operation, flipping all bits in a value.
func BITWISE_NOT(value uint32) uint32 {
	return ^value
}

// BITWISE_SHIFT_LEFT shifts bits in a value to the left by a specified number of positions.
func BITWISE_SHIFT_LEFT(value uint32, positions uint) uint32 {
	return value << positions
}

// BITWISE_SHIFT_RIGHT shifts bits in a value to the right by a specified number of positions.
func BITWISE_SHIFT_RIGHT(value uint32, positions uint) uint32 {
	return value >> positions
}

// BITWISE_ROTATE_LEFT rotates bits in a value to the left by a specified number of positions.
func BITWISE_ROTATE_LEFT(value uint32, positions uint) uint32 {
	return (value << positions) | (value >> (32 - positions))
}

// BITWISE_ROTATE_RIGHT rotates bits in a value to the right by a specified number of positions.
func BITWISE_ROTATE_RIGHT(value uint32, positions uint) uint32 {
	return (value >> positions) | (value << (32 - positions))
}

// BITWISE_CLEAR clears specific bits in a value based on a mask.
func BITWISE_CLEAR(value, mask uint32) uint32 {
	return value &^ mask
}

// BITWISE_SET sets specific bits in a value based on a mask.
func BITWISE_SET(value, mask uint32) uint32 {
	return value | mask
}

// BITWISE_TOGGLE toggles (inverts) specific bits in a value based on a mask.
func BITWISE_TOGGLE(value, mask uint32) uint32 {
	return value ^ mask
}

// BITWISE_MASK applies a mask to a value, returning only the bits allowed by the mask.
func BITWISE_MASK(value, mask uint32) uint32 {
	return value & mask
}

// BITWISE_CHECK checks if specific bits (defined by a mask) are set in a value.
func BITWISE_CHECK(value, mask uint32) bool {
	return value&mask == mask
}

// BITWISE_AND_NOT performs a bitwise AND-NOT operation, clearing bits in value1 that are set in value2.
func BITWISE_AND_NOT(value1, value2 uint32) uint32 {
	return value1 &^ value2
}

// BITWISE_OR_NOT performs a bitwise OR with the negation of the second value.
func BITWISE_OR_NOT(value1, value2 uint32) uint32 {
	return value1 | (^value2)
}

// BITWISE_XOR_NOT performs a bitwise XOR with the negation of the second value.
func BITWISE_XOR_NOT(value1, value2 uint32) uint32 {
	return value1 ^ (^value2)
}
