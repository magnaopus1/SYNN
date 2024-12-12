package math_and_logical

import (
	"errors"
	"math"
)

// Addition returns the sum of two numbers
func Addition(a, b float64) float64 {
	return a + b
}

// Subtraction returns the result of subtracting the second number from the first
func Subtraction(a, b float64) float64 {
	return a - b
}

// Multiplication returns the product of two numbers
func Multiplication(a, b float64) float64 {
	return a * b
}

// Division returns the result of dividing the first number by the second
func Division(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero is not allowed")
	}
	return a / b, nil
}

// Modulus returns the remainder of dividing the first number by the second
func Modulus(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("modulus by zero is not allowed")
	}
	return a % b, nil
}

// Exponentiation returns the result of raising the base to the power of exponent
func Exponentiation(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

// Increment increases a number by 1
func Increment(a int) int {
	return a + 1
}

// Decrement decreases a number by 1
func Decrement(a int) int {
	return a - 1
}

// AbsoluteValue returns the absolute value of a number
func AbsoluteValue(a float64) float64 {
	return math.Abs(a)
}

// Negation returns the negation of a number
func Negation(a float64) float64 {
	return -a
}

// SquareRoot returns the square root of a number
func SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, errors.New("square root of negative number is not allowed")
	}
	return math.Sqrt(a), nil
}

// Logarithm returns the natural logarithm of a number
func Logarithm(a float64) (float64, error) {
	if a <= 0 {
		return 0, errors.New("logarithm of non-positive numbers is not defined")
	}
	return math.Log(a), nil
}
