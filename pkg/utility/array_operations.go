package utility

import (
	"math"
	"reflect"
)

// ArrayEqual checks if two arrays have equal elements and order.
func ArrayEqual(arr1, arr2 []interface{}) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for i, v := range arr1 {
		if v != arr2[i] {
			return false
		}
	}
	return true
}

// ArrayNotEqual checks if two arrays do not have equal elements or order.
func ArrayNotEqual(arr1, arr2 []interface{}) bool {
	return !ArrayEqual(arr1, arr2)
}

// ArrayLengthEqual checks if two arrays have the same length.
func ArrayLengthEqual(arr1, arr2 []interface{}) bool {
	return len(arr1) == len(arr2)
}

// ArrayLengthNotEqual checks if two arrays do not have the same length.
func ArrayLengthNotEqual(arr1, arr2 []interface{}) bool {
	return !ArrayLengthEqual(arr1, arr2)
}

// ListContains checks if a list contains a specific value.
func ListContains(list []interface{}, value interface{}) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// ListNotContains checks if a list does not contain a specific value.
func ListNotContains(list []interface{}, value interface{}) bool {
	return !ListContains(list, value)
}

// IsDivisibleBy checks if a number is divisible by another.
func IsDivisibleBy(num, divisor int) bool {
	if divisor == 0 {
		return false
	}
	return num%divisor == 0
}

// IsMultipleOf checks if a number is a multiple of another.
func IsMultipleOf(num, factor int) bool {
	return IsDivisibleBy(num, factor)
}

// IsPrime checks if a number is prime.
func IsPrime(num int) bool {
	if num <= 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(num))); i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

// IsComposite checks if a number is composite (non-prime and greater than 1).
func IsComposite(num int) bool {
	return num > 1 && !IsPrime(num)
}

// IdentityCheck checks if two variables have identical values and types.
func IdentityCheck(value1, value2 interface{}) bool {
	return reflect.DeepEqual(value1, value2)
}

// ValueInSet checks if a value is present in a predefined set.
func ValueInSet(value interface{}, set []interface{}) bool {
	return ListContains(set, value)
}

// ValueNotInSet checks if a value is not present in a predefined set.
func ValueNotInSet(value interface{}, set []interface{}) bool {
	return !ListContains(set, value)
}

// IsPowerOfTwo checks if a number is a power of two.
func IsPowerOfTwo(num int) bool {
	return num > 0 && (num&(num-1)) == 0
}

// IsNotPowerOfTwo checks if a number is not a power of two.
func IsNotPowerOfTwo(num int) bool {
	return !IsPowerOfTwo(num)
}

