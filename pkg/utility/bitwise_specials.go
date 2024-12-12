package utility

// BITWISE_DUPLICATE duplicates the value by concatenating it with itself in a 64-bit space.
func BITWISE_DUPLICATE(value uint32) uint64 {
	return uint64(value)<<32 | uint64(value)
}

// BITWISE_QUAD_SWAP swaps every 4 bits in the value.
func BITWISE_QUAD_SWAP(value uint32) uint32 {
	return ((value & 0xF0F0F0F0) >> 4) | ((value & 0x0F0F0F0F) << 4)
}

// BITWISE_TRIP_SWAP swaps every 3 bits in the value.
func BITWISE_TRIP_SWAP(value uint32) uint32 {
	mask := uint32(0x49249249) // 001001001 repeated
	return ((value & mask) << 2) | ((value >> 1) & mask) | ((value >> 2) & mask)
}

// BITWISE_REVERSE_NIBBLES reverses the order of nibbles (4-bit units) in the value.
func BITWISE_REVERSE_NIBBLES(value uint32) uint32 {
	return ((value & 0xF0F0F0F0) >> 4) | ((value & 0x0F0F0F0F) << 4)
}

// BITWISE_ISOLATE_BYTE isolates a specific byte (0-3) in the value.
func BITWISE_ISOLATE_BYTE(value uint32, bytePos uint) uint32 {
	return (value >> (bytePos * 8)) & 0xFF
}

// BITWISE_EXTRACT_BYTE extracts and returns a specific byte (0-3) in the value.
func BITWISE_EXTRACT_BYTE(value uint32, bytePos uint) uint32 {
	return (value >> (8 * bytePos)) & 0xFF
}

// BITWISE_PACK packs four bytes into a 32-bit value.
func BITWISE_PACK(b1, b2, b3, b4 uint8) uint32 {
	return (uint32(b1) << 24) | (uint32(b2) << 16) | (uint32(b3) << 8) | uint32(b4)
}

// BITWISE_UNPACK unpacks a 32-bit value into four separate bytes.
func BITWISE_UNPACK(value uint32) (uint8, uint8, uint8, uint8) {
	return uint8(value >> 24), uint8(value >> 16 & 0xFF), uint8(value >> 8 & 0xFF), uint8(value & 0xFF)
}

// BITWISE_SET_HIGH sets the high 16 bits of a 32-bit value.
func BITWISE_SET_HIGH(value uint32, high uint16) uint32 {
	return (value & 0x0000FFFF) | (uint32(high) << 16)
}

// BITWISE_SET_LOW sets the low 16 bits of a 32-bit value.
func BITWISE_SET_LOW(value uint32, low uint16) uint32 {
	return (value & 0xFFFF0000) | uint32(low)
}

// BITWISE_MERGE_BYTES merges two 16-bit halves into a 32-bit value.
func BITWISE_MERGE_BYTES(high, low uint16) uint32 {
	return (uint32(high) << 16) | uint32(low)
}

// BITWISE_SEPARATE_BYTES separates a 32-bit value into two 16-bit halves.
func BITWISE_SEPARATE_BYTES(value uint32) (uint16, uint16) {
	return uint16(value >> 16), uint16(value & 0xFFFF)
}

// BITWISE_BROADCAST broadcasts the least significant byte across all bytes in a 32-bit value.
func BITWISE_BROADCAST(value uint8) uint32 {
	return uint32(value) | uint32(value)<<8 | uint32(value)<<16 | uint32(value)<<24
}

// BITWISE_AND_MASK performs an AND operation on value with a specific mask.
func BITWISE_AND_MASK(value, mask uint32) uint32 {
	return value & mask
}

// BITWISE_OR_MASK performs an OR operation on value with a specific mask.
func BITWISE_OR_MASK(value, mask uint32) uint32 {
	return value | mask
}

// BITWISE_XOR_MASK performs an XOR operation on value with a specific mask.
func BITWISE_XOR_MASK(value, mask uint32) uint32 {
	return value ^ mask
}

// BITWISE_LEFT_SHIFT_DOUBLE performs a left shift by double the specified bits.
func BITWISE_LEFT_SHIFT_DOUBLE(value uint32, bits uint) uint32 {
	return value << (2 * bits)
}

// BITWISE_RIGHT_SHIFT_DOUBLE performs a right shift by double the specified bits.
func BITWISE_RIGHT_SHIFT_DOUBLE(value uint32, bits uint) uint32 {
	return value >> (2 * bits)
}

// BITWISE_LEFT_SHIFT_MASKED shifts left by specified bits and applies a mask.
func BITWISE_LEFT_SHIFT_MASKED(value uint32, bits uint, mask uint32) uint32 {
	return (value << bits) & mask
}

// BITWISE_RIGHT_SHIFT_MASKED shifts right by specified bits and applies a mask.
func BITWISE_RIGHT_SHIFT_MASKED(value uint32, bits uint, mask uint32) uint32 {
	return (value >> bits) & mask
}

// BITWISE_SHIFT_CYCLIC performs a cyclic shift by a specified number of bits.
func BITWISE_SHIFT_CYCLIC(value uint32, bits uint) uint32 {
	return (value >> bits) | (value << (32 - bits))
}

// BITWISE_SET_TO_ZERO sets the value to zero.
func BITWISE_SET_TO_ZERO() uint32 {
	return 0
}

// BITWISE_SET_TO_ONE sets the value to one.
func BITWISE_SET_TO_ONE() uint32 {
	return 1
}

// BITWISE_RESET_TO_ZERO resets a specified bit to zero.
func BITWISE_RESET_TO_ZERO(value uint32, bitPos uint) uint32 {
	return value &^ (1 << bitPos)
}
