package utility

// BITWISE_INSERT inserts a specified bit pattern into a value at a given position.
func BITWISE_INSERT(value, pattern uint32, pos uint) uint32 {
	mask := pattern << pos
	return (value &^ mask) | mask
}

// BITWISE_EXTRACT extracts a sequence of bits from a value starting from a given position.
func BITWISE_EXTRACT(value uint32, pos, length uint) uint32 {
	mask := (1<<length - 1) << pos
	return (value & mask) >> pos
}

// BITWISE_MERGE merges two values by combining bits based on a specified mask.
func BITWISE_MERGE(value1, value2, mask uint32) uint32 {
	return (value1 &^ mask) | (value2 & mask)
}

// BITWISE_SEPARATE separates bits based on a mask, keeping specified bits in the result.
func BITWISE_SEPARATE(value, mask uint32) uint32 {
	return value & mask
}

// BITWISE_SWAP_BITS swaps individual bits at two positions in a value.
func BITWISE_SWAP_BITS(value uint32, pos1, pos2 uint) uint32 {
	bit1 := (value >> pos1) & 1
	bit2 := (value >> pos2) & 1
	if bit1 != bit2 {
		value ^= (1 << pos1) | (1 << pos2)
	}
	return value
}

// BITWISE_BYTE_SWAP swaps bytes within a 32-bit integer.
func BITWISE_BYTE_SWAP(value uint32) uint32 {
	return (value&0xFF)<<24 | (value&0xFF00)<<8 | (value&0xFF0000)>>8 | (value&0xFF000000)>>24
}

// BITWISE_NIBBLE_SWAP swaps nibbles (4-bit segments) within a byte.
func BITWISE_NIBBLE_SWAP(value uint8) uint8 {
	return (value&0xF)<<4 | (value&0xF0)>>4
}

// BITWISE_CLEAR_LEFTMOST clears the leftmost 1-bit in a 32-bit integer.
func BITWISE_CLEAR_LEFTMOST(value uint32) uint32 {
	if value == 0 {
		return 0
	}
	position := 31
	for (value & (1 << position)) == 0 {
		position--
	}
	return value &^ (1 << position)
}

// BITWISE_CLEAR_RIGHTMOST clears the rightmost 1-bit in a 32-bit integer.
func BITWISE_CLEAR_RIGHTMOST(value uint32) uint32 {
	return value & (value - 1)
}

// BITWISE_PROPAGATE_RIGHTMOST_ONE propagates the rightmost 1-bit to the right.
func BITWISE_PROPAGATE_RIGHTMOST_ONE(value uint32) uint32 {
	return value | (value - 1)
}

// BITWISE_PROPAGATE_LEFTMOST_ONE propagates the leftmost 1-bit to the left.
func BITWISE_PROPAGATE_LEFTMOST_ONE(value uint32) uint32 {
	if value == 0 {
		return 0
	}
	for i := uint(1); i < 32; i <<= 1 {
		value |= value >> i
	}
	return value
}

// BITWISE_FIND_FIRST_ONE finds the position of the first 1-bit (from the left).
func BITWISE_FIND_FIRST_ONE(value uint32) int {
	for i := 31; i >= 0; i-- {
		if (value & (1 << uint(i))) != 0 {
			return i
		}
	}
	return -1 // No 1-bit found
}

// BITWISE_FIND_LAST_ONE finds the position of the last 1-bit (from the right).
func BITWISE_FIND_LAST_ONE(value uint32) int {
	for i := 0; i < 32; i++ {
		if (value & (1 << uint(i))) != 0 {
			return i
		}
	}
	return -1 // No 1-bit found
}

// BITWISE_RESET_FIRST_ONE resets the first 1-bit (from the left) in a 32-bit integer.
func BITWISE_RESET_FIRST_ONE(value uint32) uint32 {
	pos := BITWISE_FIND_FIRST_ONE(value)
	if pos == -1 {
		return value
	}
	return value &^ (1 << uint(pos))
}

// BITWISE_RESET_LAST_ONE resets the last 1-bit (from the right) in a 32-bit integer.
func BITWISE_RESET_LAST_ONE(value uint32) uint32 {
	pos := BITWISE_FIND_LAST_ONE(value)
	if pos == -1 {
		return value
	}
	return value &^ (1 << uint(pos))
}
