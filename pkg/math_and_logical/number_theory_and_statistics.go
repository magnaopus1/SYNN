package math_and_logical

import (
	"errors"
	"math"
)

// PrimeCheck checks if a number is prime
func PrimeCheck(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// PrimeFactorize returns the prime factors of a given number
func PrimeFactorize(n int) ([]int, error) {
	if n <= 1 {
		return nil, errors.New("number must be greater than 1 for prime factorization")
	}
	factors := []int{}
	for n%2 == 0 {
		factors = append(factors, 2)
		n /= 2
	}
	for i := 3; i*i <= n; i += 2 {
		for n%i == 0 {
			factors = append(factors, i)
			n /= i
		}
	}
	if n > 2 {
		factors = append(factors, n)
	}
	return factors, nil
}

// GCDCalculate computes the greatest common divisor (GCD) of two numbers
func GCDCalculate(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCMCalculate computes the least common multiple (LCM) of two numbers
func LCMCalculate(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return int(math.Abs(float64(a*b))) / GCDCalculate(a, b)
}

// LinearRegression calculates the slope and intercept of the best-fit line for a set of points
func LinearRegression(x, y []float64) (float64, float64, error) {
	if len(x) != len(y) || len(x) < 2 {
		return 0, 0, errors.New("x and y must have the same length and contain at least two points")
	}
	n := float64(len(x))
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0
	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n
	return slope, intercept, nil
}

// PolynomialFit fits a polynomial to a set of points using least squares, returning the coefficients
func PolynomialFit(x, y []float64, degree int) ([]float64, error) {
	if len(x) != len(y) || len(x) < degree+1 {
		return nil, errors.New("insufficient data points for the requested polynomial degree")
	}
	n := len(x)
	matrix := make([][]float64, degree+1)
	for i := range matrix {
		matrix[i] = make([]float64, degree+2)
		for j := 0; j <= degree; j++ {
			for k := 0; k < n; k++ {
				matrix[i][j] += math.Pow(x[k], float64(i+j))
			}
		}
		for k := 0; k < n; k++ {
			matrix[i][degree+1] += y[k] * math.Pow(x[k], float64(i))
		}
	}
	return gaussJordan(matrix, degree+1), nil
}

// GammaFunction computes the gamma function for a given value
func GammaFunction(x float64) float64 {
	if x <= 0 {
		return math.NaN()
	}
	const p = 7
	const a = 0.99999999999980993
	coefficients := []float64{
		676.5203681218851,
		-1259.1392167224028,
		771.32342877765313,
		-176.61502916214059,
		12.507343278686905,
		-0.13857109526572012,
		9.9843695780195716e-6,
		1.5056327351493116e-7,
	}
	y := x + float64(p) - 0.5
	sum := a
	for i := 0; i < len(coefficients); i++ {
		sum += coefficients[i] / (x + float64(i))
	}
	return math.Sqrt(2*math.Pi) * math.Pow(y, x-0.5) * math.Exp(-y) * sum
}

// BetaFunction computes the beta function for given values x and y
func BetaFunction(x, y float64) float64 {
	return GammaFunction(x) * GammaFunction(y) / GammaFunction(x+y)
}

// SigmoidFunction calculates the sigmoid function for a given value
func SigmoidFunction(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// SoftmaxFunction computes the softmax function for a slice of values
func SoftmaxFunction(values []float64) []float64 {
	expValues := make([]float64, len(values))
	sumExp := 0.0
	for i, v := range values {
		expValues[i] = math.Exp(v)
		sumExp += expValues[i]
	}
	for i := range expValues {
		expValues[i] /= sumExp
	}
	return expValues
}

// Helper function for PolynomialFit: Gaussian-Jordan elimination for solving linear equations
func gaussJordan(matrix [][]float64, n int) []float64 {
	for i := 0; i < n; i++ {
		max := i
		for j := i + 1; j < n; j++ {
			if math.Abs(matrix[j][i]) > math.Abs(matrix[max][i]) {
				max = j
			}
		}
		matrix[i], matrix[max] = matrix[max], matrix[i]
		for j := i + 1; j < n+1; j++ {
			matrix[i][j] /= matrix[i][i]
		}
		for j := 0; j < n; j++ {
			if j != i {
				factor := matrix[j][i]
				for k := i; k < n+1; k++ {
					matrix[j][k] -= factor * matrix[i][k]
				}
			}
		}
	}
	coefficients := make([]float64, n)
	for i := 0; i < n; i++ {
		coefficients[i] = matrix[i][n]
	}
	return coefficients
}
