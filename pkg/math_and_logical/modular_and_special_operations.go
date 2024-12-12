package math_and_logical

import (
	"errors"
)

// ModuloAdd performs modular addition: (a + b) % mod
func ModuloAdd(a, b, mod int) (int, error) {
	if mod <= 0 {
		return 0, errors.New("modulus must be positive")
	}
	return (a + b) % mod, nil
}

// ModuloSubtract performs modular subtraction: (a - b) % mod
func ModuloSubtract(a, b, mod int) (int, error) {
	if mod <= 0 {
		return 0, errors.New("modulus must be positive")
	}
	return (a - b + mod) % mod, nil
}

// ModuloMultiply performs modular multiplication: (a * b) % mod
func ModuloMultiply(a, b, mod int) (int, error) {
	if mod <= 0 {
		return 0, errors.New("modulus must be positive")
	}
	return (a * b) % mod, nil
}

// ModuloDivide performs modular division: (a / b) % mod
// Uses modular multiplicative inverse for division in a modular system
func ModuloDivide(a, b, mod int) (int, error) {
	if mod <= 0 {
		return 0, errors.New("modulus must be positive")
	}
	inverse, err := modularInverse(b, mod)
	if err != nil {
		return 0, err
	}
	return (a * inverse) % mod, nil
}

// Remainder returns the remainder of division: a % b
func Remainder(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero is undefined")
	}
	return a % b, nil
}

// Wrap wraps an integer within a specified range [0, max)
// Returns an integer between 0 and max - 1
func Wrap(value, max int) (int, error) {
	if max <= 0 {
		return 0, errors.New("maximum bound must be positive")
	}
	return ((value % max) + max) % max, nil
}

// Helper function to find the modular multiplicative inverse of a number
// Uses the extended Euclidean algorithm
func modularInverse(a, mod int) (int, error) {
	t, newT := 0, 1
	r, newR := mod, a

	for newR != 0 {
		quotient := r / newR
		t, newT = newT, t-quotient*newT
		r, newR = newR, r-quotient*newR
	}

	if r > 1 {
		return 0, errors.New("no modular inverse exists")
	}
	if t < 0 {
		t += mod
	}
	return t, nil
}
