package math_and_logical

import (
	"errors"
	"math"
)

// VarianceCalculation calculates the variance of a given dataset
func VarianceCalculation(data []float64) (float64, error) {
	if len(data) < 2 {
		return 0, errors.New("at least two data points are required to calculate variance")
	}
	mean, _ := Average(data)
	variance := 0.0
	for _, value := range data {
		variance += math.Pow(value-mean, 2)
	}
	return variance / float64(len(data)-1), nil
}

// Average calculates the mean (average) of a given dataset
func Average(data []float64) (float64, error) {
	if len(data) == 0 {
		return 0, errors.New("data slice cannot be empty")
	}
	sum := Sum(data)
	return sum / float64(len(data)), nil
}

// Sum calculates the sum of elements in a dataset
func Sum(data []float64) float64 {
	total := 0.0
	for _, value := range data {
		total += value
	}
	return total
}

// Product calculates the product of elements in a dataset
func Product(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	result := 1.0
	for _, value := range data {
		result *= value
	}
	return result
}

// Difference calculates the difference between two numbers
func Difference(a, b float64) float64 {
	return a - b
}

// Clamp restricts a value within a specified range [min, max]
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}



// AbsoluteDifference calculates the absolute difference between two numbers
func AbsoluteDifference(a, b float64) float64 {
	return math.Abs(a - b)
}
