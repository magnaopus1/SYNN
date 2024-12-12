package math_and_logical

import (
	"errors"
	"math"
	"math/cmplx"
)

// FFTCalculate performs a Fast Fourier Transform on a slice of complex numbers
func FFTCalculate(data []complex128) ([]complex128, error) {
	n := len(data)
	if n == 0 || (n&(n-1)) != 0 {
		return nil, errors.New("data length must be a power of 2")
	}
	return fft(data), nil
}

// IFFTCalculate performs an Inverse Fast Fourier Transform on a slice of complex numbers
func IFFTCalculate(data []complex128) ([]complex128, error) {
	n := len(data)
	if n == 0 || (n&(n-1)) != 0 {
		return nil, errors.New("data length must be a power of 2")
	}
	data = fft(data)
	for i := range data {
		data[i] /= complex(float64(n), 0)
	}
	return data, nil
}

// FourierTransform performs a discrete Fourier Transform on a slice of complex numbers
func FourierTransform(data []complex128) []complex128 {
	n := len(data)
	result := make([]complex128, n)
	for k := 0; k < n; k++ {
		for t := 0; t < n; t++ {
			angle := -2 * math.Pi * float64(k) * float64(t) / float64(n)
			result[k] += data[t] * cmplx.Exp(complex(0, angle))
		}
	}
	return result
}

// InverseFourierTransform performs an inverse discrete Fourier Transform
func InverseFourierTransform(data []complex128) []complex128 {
	n := len(data)
	result := make([]complex128, n)
	for k := 0; k < n; k++ {
		for t := 0; t < n; t++ {
			angle := 2 * math.Pi * float64(k) * float64(t) / float64(n)
			result[k] += data[t] * cmplx.Exp(complex(0, angle))
		}
		result[k] /= complex(float64(n), 0)
	}
	return result
}

// ExponentialFunction calculates e raised to the power of x
func ExponentialFunction(x float64) float64 {
	return math.Exp(x)
}

// LogarithmBaseN calculates the logarithm of a number to a custom base
func LogarithmBaseN(x, base float64) (float64, error) {
	if x <= 0 || base <= 0 || base == 1 {
		return 0, errors.New("invalid input for logarithm calculation")
	}
	return math.Log(x) / math.Log(base), nil
}

// TaylorSeriesExpansion calculates the Taylor series expansion for e^x at a specified order
func TaylorSeriesExpansion(x float64, order int) float64 {
	result := 1.0
	factorial := 1.0
	power := 1.0
	for i := 1; i <= order; i++ {
		factorial *= float64(i)
		power *= x
		result += power / factorial
	}
	return result
}

// HyperbolicSine calculates the hyperbolic sine of a number
func HyperbolicSine(x float64) float64 {
	return math.Sinh(x)
}

// HyperbolicCosine calculates the hyperbolic cosine of a number
func HyperbolicCosine(x float64) float64 {
	return math.Cosh(x)
}

// HyperbolicTangent calculates the hyperbolic tangent of a number
func HyperbolicTangent(x float64) float64 {
	return math.Tanh(x)
}

// InverseHyperbolicSine calculates the inverse hyperbolic sine of a number
func InverseHyperbolicSine(x float64) float64 {
	return math.Asinh(x)
}

// InverseHyperbolicCosine calculates the inverse hyperbolic cosine of a number
func InverseHyperbolicCosine(x float64) (float64, error) {
	if x < 1 {
		return 0, errors.New("inverse hyperbolic cosine undefined for values less than 1")
	}
	return math.Acosh(x), nil
}

// InverseHyperbolicTangent calculates the inverse hyperbolic tangent of a number
func InverseHyperbolicTangent(x float64) (float64, error) {
	if x <= -1 || x >= 1 {
		return 0, errors.New("inverse hyperbolic tangent undefined for values <= -1 or >= 1")
	}
	return math.Atanh(x), nil
}

// EllipticIntegral calculates the complete elliptic integral of the first kind for a given modulus
func EllipticIntegral(k float64) (float64, error) {
	if k < 0 || k > 1 {
		return 0, errors.New("modulus must be between 0 and 1 for elliptic integral")
	}
	const epsilon = 1e-10
	a, b := 1.0, math.Sqrt(1-k*k)
	sum := math.Pi / 2.0
	for math.Abs(a-b) > epsilon {
		a, b = (a+b)/2, math.Sqrt(a*b)
		sum -= (a - b) / (2 * math.Pi)
	}
	return sum, nil
}

// Helper function: fft calculates the Fast Fourier Transform using the Cooley-Tukey algorithm
func fft(data []complex128) []complex128 {
	n := len(data)
	if n <= 1 {
		return data
	}
	even := fft(data[0:n:2])
	odd := fft(data[1:n:2])

	combined := make([]complex128, n)
	for k := 0; k < n/2; k++ {
		t := cmplx.Exp(complex(0, -2*math.Pi*float64(k)/float64(n))) * odd[k]
		combined[k] = even[k] + t
		combined[k+n/2] = even[k] - t
	}
	return combined
}
