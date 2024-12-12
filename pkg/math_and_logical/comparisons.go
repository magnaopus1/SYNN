package math_and_logical

import (
	"errors"
)

// CompareEqual checks if two values are equal
func CompareEqual(a, b interface{}) (bool, error) {
	switch a := a.(type) {
	case int:
		if b, ok := b.(int); ok {
			return a == b, nil
		}
	case float64:
		if b, ok := b.(float64); ok {
			return a == b, nil
		}
	case string:
		if b, ok := b.(string); ok {
			return a == b, nil
		}
	default:
		return false, errors.New("unsupported data type for comparison")
	}
	return false, errors.New("type mismatch for comparison")
}

// CompareNotEqual checks if two values are not equal
func CompareNotEqual(a, b interface{}) (bool, error) {
	equal, err := CompareEqual(a, b)
	if err != nil {
		return false, err
	}
	return !equal, nil
}

// CompareGreaterThan checks if the first value is greater than the second
func CompareGreaterThan(a, b interface{}) (bool, error) {
	switch a := a.(type) {
	case int:
		if b, ok := b.(int); ok {
			return a > b, nil
		}
	case float64:
		if b, ok := b.(float64); ok {
			return a > b, nil
		}
	default:
		return false, errors.New("unsupported data type for greater-than comparison")
	}
	return false, errors.New("type mismatch for greater-than comparison")
}

// CompareLessThan checks if the first value is less than the second
func CompareLessThan(a, b interface{}) (bool, error) {
	switch a := a.(type) {
	case int:
		if b, ok := b.(int); ok {
			return a < b, nil
		}
	case float64:
		if b, ok := b.(float64); ok {
			return a < b, nil
		}
	default:
		return false, errors.New("unsupported data type for less-than comparison")
	}
	return false, errors.New("type mismatch for less-than comparison")
}

// CompareGreaterEqual checks if the first value is greater than or equal to the second
func CompareGreaterEqual(a, b interface{}) (bool, error) {
	greater, err := CompareGreaterThan(a, b)
	if err != nil {
		return false, err
	}
	if greater {
		return true, nil
	}
	return CompareEqual(a, b)
}

// CompareLessEqual checks if the first value is less than or equal to the second
func CompareLessEqual(a, b interface{}) (bool, error) {
	less, err := CompareLessThan(a, b)
	if err != nil {
		return false, err
	}
	if less {
		return true, nil
	}
	return CompareEqual(a, b)
}
