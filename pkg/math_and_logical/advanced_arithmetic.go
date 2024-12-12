package math_and_logical

import (
	"errors"
	"math"
	"sort"
)

// Factorial calculates the factorial of a given number
func Factorial(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("factorial is not defined for negative numbers")
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result, nil
}

// Permutation calculates the number of permutations of n items taken r at a time
func Permutation(n, r int) (int, error) {
	if n < 0 || r < 0 || r > n {
		return 0, errors.New("invalid values for n or r")
	}
	num, _ := Factorial(n)
	den, _ := Factorial(n - r)
	return num / den, nil
}

// Combination calculates the number of combinations of n items taken r at a time
func Combination(n, r int) (int, error) {
	if n < 0 || r < 0 || r > n {
		return 0, errors.New("invalid values for n or r")
	}
	num, _ := Factorial(n)
	den1, _ := Factorial(r)
	den2, _ := Factorial(n - r)
	return num / (den1 * den2), nil
}

// Hypotenuse calculates the hypotenuse of a right-angled triangle given sides a and b
func Hypotenuse(a, b float64) float64 {
	return math.Sqrt(a*a + b*b)
}

// Mean calculates the mean of a slice of numbers
func Mean(numbers []float64) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("empty slice provided")
	}
	total := 0.0
	for _, num := range numbers {
		total += num
	}
	return total / float64(len(numbers)), nil
}

// Median calculates the median of a slice of numbers
func Median(numbers []float64) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("empty slice provided")
	}
	sort.Float64s(numbers)
	mid := len(numbers) / 2
	if len(numbers)%2 == 0 {
		return (numbers[mid-1] + numbers[mid]) / 2.0, nil
	}
	return numbers[mid], nil
}

// Mode calculates the mode of a slice of numbers
func Mode(numbers []float64) ([]float64, error) {
	if len(numbers) == 0 {
		return nil, errors.New("empty slice provided")
	}
	counts := make(map[float64]int)
	for _, num := range numbers {
		counts[num]++
	}
	maxCount := 0
	var modes []float64
	for num, count := range counts {
		if count > maxCount {
			maxCount = count
			modes = []float64{num}
		} else if count == maxCount {
			modes = append(modes, num)
		}
	}
	return modes, nil
}

// Power calculates the power of base raised to the exponent
func Power(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

// Root calculates the n-th root of a number
func Root(number, n float64) (float64, error) {
	if n == 0 {
		return 0, errors.New("cannot calculate the zero-th root")
	}
	return math.Pow(number, 1/n), nil
}

// LogBase10 calculates the logarithm of a number to base 10
func LogBase10(x float64) (float64, error) {
	if x <= 0 {
		return 0, errors.New("logarithm undefined for non-positive values")
	}
	return math.Log10(x), nil
}

// LogBase2 calculates the logarithm of a number to base 2
func LogBase2(x float64) (float64, error) {
	if x <= 0 {
		return 0, errors.New("logarithm undefined for non-positive values")
	}
	return math.Log2(x), nil
}

// LogBaseN calculates the logarithm of a number to a custom base
func LogBaseN(x, base float64) (float64, error) {
	if x <= 0 || base <= 0 || base == 1 {
		return 0, errors.New("invalid input for logarithm calculation")
	}
	return math.Log(x) / math.Log(base), nil
}

// StandardDeviation calculates the standard deviation of a slice of numbers
func StandardDeviation(numbers []float64) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("empty slice provided")
	}
	mean, _ := Mean(numbers)
	variance := 0.0
	for _, num := range numbers {
		variance += (num - mean) * (num - mean)
	}
	variance /= float64(len(numbers))
	return math.Sqrt(variance), nil
}

// Variance calculates the variance of a slice of numbers
func Variance(numbers []float64) (float64, error) {
	if len(numbers) == 0 {
		return 0, errors.New("empty slice provided")
	}
	mean, _ := Mean(numbers)
	variance := 0.0
	for _, num := range numbers {
		variance += (num - mean) * (num - mean)
	}
	return variance / float64(len(numbers)), nil
}
