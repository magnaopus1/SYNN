package utility

import "math/bits"

// BITWISE_FIND_FIRST_ZERO finds the position of the first zero bit from the least significant bit (LSB).
func BITWISE_FIND_FIRST_ZERO(value uint32) int {
	for i := 0; i < 32; i++ {
		if value&(1<<i) == 0 {
			return i
		}
	}
	return -1 // Returns -1 if no zero bit is found
}

// BITWISE_FIND_LAST_ZERO finds the position of the last zero bit from the most significant bit (MSB).
func BITWISE_FIND_LAST_ZERO(value uint32) int {
	for i := 31; i >= 0; i-- {
		if value&(1<<i) == 0 {
			return i
		}
	}
	return -1 // Returns -1 if no zero bit is found
}

// BITWISE_MIRROR_BITS reverses the order of bits in a 32-bit integer.
func BITWISE_MIRROR_BITS(value uint32) uint32 {
	return bits.Reverse32(value)
}

// BITWISE_POP_COUNT counts the number of 1 bits in a 32-bit integer.
func BITWISE_POP_COUNT(value uint32) int {
	return bits.OnesCount32(value)
}

// BITWISE_IS_POWER_OF_TWO checks if the given value is a power of two.
func BITWISE_IS_POWER_OF_TWO(value uint32) bool {
	return value != 0 && (value&(value-1)) == 0
}

// BITWISE_NEXT_POWER_OF_TWO calculates the next power of two greater than or equal to the value.
func BITWISE_NEXT_POWER_OF_TWO(value uint32) uint32 {
	if value == 0 {
		return 1
	}
	value--
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	return value + 1
}

// BITWISE_PREV_POWER_OF_TWO calculates the previous power of two less than or equal to the value.
func BITWISE_PREV_POWER_OF_TWO(value uint32) uint32 {
	if value == 0 {
		return 0
	}
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	return value - (value >> 1)
}

// BITWISE_HAMMING_DISTANCE calculates the Hamming distance between two 32-bit integers.
func BITWISE_HAMMING_DISTANCE(a, b uint32) int {
	return bits.OnesCount32(a ^ b)
}

// BITWISE_ADD performs addition using bitwise operations.
func BITWISE_ADD(a, b uint32) uint32 {
	for b != 0 {
		carry := (a & b) << 1
		a = a ^ b
		b = carry
	}
	return a
}

// BITWISE_SUBTRACT performs subtraction using bitwise operations.
func BITWISE_SUBTRACT(a, b uint32) uint32 {
	for b != 0 {
		borrow := (^a & b) << 1
		a = a ^ b
		b = borrow
	}
	return a
}

// BITWISE_MULTIPLY performs multiplication using bitwise operations.
func BITWISE_MULTIPLY(a, b uint32) uint32 {
	var result uint32 = 0
	for b != 0 {
		if b&1 != 0 {
			result = BITWISE_ADD(result, a)
		}
		a <<= 1
		b >>= 1
	}
	return result
}

// BITWISE_DIVIDE performs division using bitwise operations.
func BITWISE_DIVIDE(dividend, divisor uint32) (uint32, uint32) {
	var quotient uint32 = 0
	var remainder uint32 = 0
	for i := 31; i >= 0; i-- {
		remainder <<= 1
		remainder |= (dividend >> i) & 1
		if remainder >= divisor {
			remainder -= divisor
			quotient |= (1 << i)
		}
	}
	return quotient, remainder
}

// BITWISE_SIGN_EXTEND sign-extends an n-bit value to a full 32-bit signed integer.
func BITWISE_SIGN_EXTEND(value uint32, bits uint) int32 {
	mask := int32(1) << (bits - 1)
	return int32(value) ^ mask - mask
}
