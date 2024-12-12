package utility

import (
	"errors"
)

// BITWISE_ZERO_EXTEND extends a smaller bit-width integer to a larger one by zero-padding.
func BITWISE_ZERO_EXTEND(value uint32, targetBits int) (uint64, error) {
	if targetBits != 32 && targetBits != 64 {
		return 0, errors.New("targetBits must be 32 or 64")
	}
	return uint64(value), nil
}

// BITWISE_LEFT_ROTATE_WITH_CARRY performs a left rotation with carry.
func BITWISE_LEFT_ROTATE_WITH_CARRY(value uint32, rotateBy uint, carryIn uint8) (uint32, uint8) {
	rotateBy %= 32
	newCarry := (value >> (31 - rotateBy)) & 1
	rotated := (value << rotateBy) | uint32(carryIn)
	return rotated | (value >> (32 - rotateBy)), uint8(newCarry)
}

// BITWISE_RIGHT_ROTATE_WITH_CARRY performs a right rotation with carry.
func BITWISE_RIGHT_ROTATE_WITH_CARRY(value uint32, rotateBy uint, carryIn uint8) (uint32, uint8) {
	rotateBy %= 32
	newCarry := (value >> (rotateBy - 1)) & 1
	rotated := (value >> rotateBy) | uint32(carryIn)<<31
	return rotated | (value << (32 - rotateBy)), uint8(newCarry)
}

// BITWISE_REVERSE reverses the bits in a 32-bit integer.
func BITWISE_REVERSE(value uint32) uint32 {
	var reversed uint32
	for i := 0; i < 32; i++ {
		reversed |= ((value >> i) & 1) << (31 - i)
	}
	return reversed
}

// BITWISE_FILL_PARITY fills the last bit based on parity (even/odd).
func BITWISE_FILL_PARITY(value uint32) uint32 {
	if BITWISE_CHECK_EVEN_PARITY(value) {
		return value &^ 1
	}
	return value | 1
}

// BITWISE_SUM_BITS counts the number of bits set to 1 in a 32-bit integer.
func BITWISE_SUM_BITS(value uint32) int {
	count := 0
	for value > 0 {
		count += int(value & 1)
		value >>= 1
	}
	return count
}

// BITWISE_CHECK_EVEN_PARITY returns true if the number of 1-bits is even.
func BITWISE_CHECK_EVEN_PARITY(value uint32) bool {
	return BITWISE_SUM_BITS(value)%2 == 0
}

// BITWISE_CHECK_ODD_PARITY returns true if the number of 1-bits is odd.
func BITWISE_CHECK_ODD_PARITY(value uint32) bool {
	return BITWISE_SUM_BITS(value)%2 != 0
}

// BITWISE_COMPARE compares two 32-bit values bitwise.
func BITWISE_COMPARE(value1, value2 uint32) bool {
	return value1 == value2
}

// BITWISE_SHUFFLE_BITS rearranges the bits of a 32-bit integer based on a provided pattern.
func BITWISE_SHUFFLE_BITS(value uint32, pattern []int) (uint32, error) {
	if len(pattern) != 32 {
		return 0, errors.New("pattern length must be 32")
	}
	var shuffled uint32
	for i, pos := range pattern {
		if pos < 0 || pos > 31 {
			return 0, errors.New("invalid bit position in pattern")
		}
		shuffled |= ((value >> pos) & 1) << i
	}
	return shuffled, nil
}

// BITWISE_REORDER_BYTES reverses the byte order of a 32-bit integer.
func BITWISE_REORDER_BYTES(value uint32) uint32 {
	return (value&0xFF)<<24 | (value&0xFF00)<<8 | (value&0xFF0000)>>8 | (value&0xFF000000)>>24
}

