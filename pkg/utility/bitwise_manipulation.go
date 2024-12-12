package utility

// BITWISE_CLEAR_BITS clears the specified bits in the value based on a mask.
func BITWISE_CLEAR_BITS(value, mask uint32) uint32 {
	return value &^ mask
}

// BITWISE_SET_BITS sets the specified bits in the value based on a mask.
func BITWISE_SET_BITS(value, mask uint32) uint32 {
	return value | mask
}

// BITWISE_TOGGLE_BITS toggles the specified bits in the value based on a mask.
func BITWISE_TOGGLE_BITS(value, mask uint32) uint32 {
	return value ^ mask
}

// BITWISE_COUNT_ONES counts the number of 1 bits in a 32-bit integer.
func BITWISE_COUNT_ONES(value uint32) int {
	count := 0
	for value > 0 {
		count += int(value & 1)
		value >>= 1
	}
	return count
}

// BITWISE_COUNT_ZEROS counts the number of 0 bits in a 32-bit integer.
func BITWISE_COUNT_ZEROS(value uint32) int {
	return 32 - BITWISE_COUNT_ONES(value)
}

// BITWISE_PARITY returns 1 if there is an odd number of 1 bits, 0 if even.
func BITWISE_PARITY(value uint32) int {
	return BITWISE_COUNT_ONES(value) % 2
}

// BITWISE_CHECK_BIT checks if a specific bit is set (1) in a 32-bit integer.
func BITWISE_CHECK_BIT(value uint32, position uint) bool {
	if position >= 32 {
		return false
	}
	return (value & (1 << position)) != 0
}

// BITWISE_SET_BIT sets a specific bit (to 1) in a 32-bit integer.
func BITWISE_SET_BIT(value uint32, position uint) uint32 {
	if position >= 32 {
		return value
	}
	return value | (1 << position)
}

// BITWISE_CLEAR_BIT clears a specific bit (to 0) in a 32-bit integer.
func BITWISE_CLEAR_BIT(value uint32, position uint) uint32 {
	if position >= 32 {
		return value
	}
	return value &^ (1 << position)
}

// BITWISE_TOGGLE_BIT toggles a specific bit in a 32-bit integer.
func BITWISE_TOGGLE_BIT(value uint32, position uint) uint32 {
	if position >= 32 {
		return value
	}
	return value ^ (1 << position)
}

// BITWISE_FLIP_ALL inverts all bits in a 32-bit integer.
func BITWISE_FLIP_ALL(value uint32) uint32 {
	return ^value
}

// BITWISE_FILL_ONES sets all bits in a 32-bit integer to 1.
func BITWISE_FILL_ONES() uint32 {
	return 0xFFFFFFFF
}

// BITWISE_FILL_ZEROS sets all bits in a 32-bit integer to 0.
func BITWISE_FILL_ZEROS() uint32 {
	return 0
}

// BITWISE_ISOLATE_LEFTMOST_ONE isolates the leftmost 1 bit in a 32-bit integer.
func BITWISE_ISOLATE_LEFTMOST_ONE(value uint32) uint32 {
	if value == 0 {
		return 0
	}
	position := 31
	for (value & (1 << position)) == 0 {
		position--
	}
	return 1 << position
}

// BITWISE_ISOLATE_RIGHTMOST_ONE isolates the rightmost 1 bit in a 32-bit integer.
func BITWISE_ISOLATE_RIGHTMOST_ONE(value uint32) uint32 {
	return value & -value
}

