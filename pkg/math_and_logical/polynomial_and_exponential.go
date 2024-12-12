package math_and_logical

import (
	"errors"
	"math"
)

// LogarithmicSum calculates the sum of the logarithms of each element in the input slice with a specified base
func LogarithmicSum(values []float64, base float64) (float64, error) {
	if base <= 0 || base == 1 {
		return 0, errors.New("invalid base for logarithm")
	}
	sum := 0.0
	for _, v := range values {
		if v <= 0 {
			return 0, errors.New("logarithm of non-positive values is undefined")
		}
		sum += math.Log(v) / math.Log(base)
	}
	return sum, nil
}

// ExponentialSum calculates the sum of the exponential values of each element in the input slice
func ExponentialSum(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += math.Exp(v)
	}
	return sum
}

// SquaredSum calculates the sum of the squares of each element in the input slice
func SquaredSum(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v * v
	}
	return sum
}

// Square calculates the square of a given value
func Square(value float64) float64 {
	return value * value
}

// Cube calculates the cube of a given value
func Cube(value float64) float64 {
	return value * value * value
}

// FourthRoot calculates the fourth root of a given value
func FourthRoot(value float64) (float64, error) {
	if value < 0 {
		return 0, errors.New("fourth root of negative values is undefined in real numbers")
	}
	return math.Pow(value, 0.25), nil
}

// FifthRoot calculates the fifth root of a given value
func FifthRoot(value float64) (float64, error) {
	if value < 0 {
		return -math.Pow(-value, 0.2), nil // Allows for real roots of negative values
	}
	return math.Pow(value, 0.2), nil
}
