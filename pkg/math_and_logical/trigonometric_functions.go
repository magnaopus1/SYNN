package math_and_logical

import (
	"errors"
	"math"
)

// Sine calculates the sine of a given angle in radians
func Sine(angle float64) float64 {
	return math.Sin(angle)
}

// Cosine calculates the cosine of a given angle in radians
func Cosine(angle float64) float64 {
	return math.Cos(angle)
}

// Tangent calculates the tangent of a given angle in radians
func Tangent(angle float64) float64 {
	return math.Tan(angle)
}

// ArcSine calculates the inverse sine (arcsin) of a value
// Returns the angle in radians whose sine is the given value
func ArcSine(value float64) (float64, error) {
	if value < -1 || value > 1 {
		return 0, errors.New("input out of range for arcsine")
	}
	return math.Asin(value), nil
}

// ArcCosine calculates the inverse cosine (arccos) of a value
// Returns the angle in radians whose cosine is the given value
func ArcCosine(value float64) (float64, error) {
	if value < -1 || value > 1 {
		return 0, errors.New("input out of range for arccosine")
	}
	return math.Acos(value), nil
}

// ArcTangent calculates the inverse tangent (arctan) of a value
// Returns the angle in radians whose tangent is the given value
func ArcTangent(value float64) float64 {
	return math.Atan(value)
}

