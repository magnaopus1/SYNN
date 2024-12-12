package math_and_logical

import (
	"errors"
	"math"
)

// Fibonacci calculates the n-th Fibonacci number
func Fibonacci(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("fibonacci is not defined for negative numbers")
	}
	a, b := 0, 1
	for i := 0; i < n; i++ {
		a, b = b, a+b
	}
	return a, nil
}

// BinomialCoefficient calculates the binomial coefficient "n choose k"
func BinomialCoefficient(n, k int) (int, error) {
	if k < 0 || n < k {
		return 0, errors.New("invalid values for n and k in binomial coefficient")
	}
	num, _ := Factorial(n)
	den1, _ := Factorial(k)
	den2, _ := Factorial(n - k)
	return num / (den1 * den2), nil
}

// PolynomialEvaluation evaluates a polynomial at a given x value
// Coefficients are provided as an array where the i-th element is the coefficient for x^i
func PolynomialEvaluation(coefficients []float64, x float64) float64 {
	result := 0.0
	for i, coef := range coefficients {
		result += coef * math.Pow(x, float64(i))
	}
	return result
}

// PolynomialDerivative calculates the derivative of a polynomial
// Returns a new set of coefficients for the derived polynomial
func PolynomialDerivative(coefficients []float64) []float64 {
	if len(coefficients) <= 1 {
		return []float64{0}
	}
	derivative := make([]float64, len(coefficients)-1)
	for i := 1; i < len(coefficients); i++ {
		derivative[i-1] = coefficients[i] * float64(i)
	}
	return derivative
}

// PolynomialIntegral calculates the indefinite integral of a polynomial
// Returns a new set of coefficients with an integration constant of 0
func PolynomialIntegral(coefficients []float64) []float64 {
	integral := make([]float64, len(coefficients)+1)
	for i := 0; i < len(coefficients); i++ {
		integral[i+1] = coefficients[i] / float64(i+1)
	}
	return integral
}

// SplineInterpolation performs linear spline interpolation for a given set of points
func SplineInterpolation(xValues, yValues []float64, x float64) (float64, error) {
	if len(xValues) != len(yValues) || len(xValues) < 2 {
		return 0, errors.New("invalid input: xValues and yValues must have the same length and contain at least two points")
	}
	for i := 0; i < len(xValues)-1; i++ {
		if x >= xValues[i] && x <= xValues[i+1] {
			t := (x - xValues[i]) / (xValues[i+1] - xValues[i])
			return (1-t)*yValues[i] + t*yValues[i+1], nil
		}
	}
	return 0, errors.New("x is out of bounds for the given interpolation points")
}

// BezierCurveEvaluation evaluates a Bezier curve at a parameter t (0 <= t <= 1)
// Control points are provided as an array of points, each represented as a [2]float64 for (x, y)
func BezierCurveEvaluation(controlPoints [][2]float64, t float64) ([2]float64, error) {
	if t < 0 || t > 1 {
		return [2]float64{0, 0}, errors.New("parameter t must be in the range [0, 1]")
	}
	n := len(controlPoints)
	points := make([][2]float64, n)
	copy(points, controlPoints)

	for k := 1; k < n; k++ {
		for i := 0; i < n-k; i++ {
			points[i][0] = (1-t)*points[i][0] + t*points[i+1][0]
			points[i][1] = (1-t)*points[i][1] + t*points[i+1][1]
		}
	}
	return points[0], nil
}
