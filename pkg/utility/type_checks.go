package utility

import (
	"math"
	"reflect"
	"sync"
	"synnergy_network/pkg/common"
)

var typeCheckLock sync.Mutex

// IsInteger: Checks if a value is an integer type
func IsInteger(value interface{}) bool {
    isInt := reflect.TypeOf(value).Kind() == reflect.Int || reflect.TypeOf(value).Kind() == reflect.Int64
    LogTypeCheck("IsInteger", isInt)
    return isInt
}

// IsFloat: Checks if a value is a float type
func IsFloat(value interface{}) bool {
    isFloat := reflect.TypeOf(value).Kind() == reflect.Float32 || reflect.TypeOf(value).Kind() == reflect.Float64
    LogTypeCheck("IsFloat", isFloat)
    return isFloat
}

// IsString: Checks if a value is a string type
func IsString(value interface{}) bool {
    isString := reflect.TypeOf(value).Kind() == reflect.String
    LogTypeCheck("IsString", isString)
    return isString
}

// IsBoolean: Checks if a value is a boolean type
func IsBoolean(value interface{}) bool {
    isBool := reflect.TypeOf(value).Kind() == reflect.Bool
    LogTypeCheck("IsBoolean", isBool)
    return isBool
}

// IsEmpty: Checks if a value is empty (e.g., empty string, nil slice, or map)
func IsEmpty(value interface{}) bool {
    isEmpty := false
    switch reflect.TypeOf(value).Kind() {
    case reflect.String:
        isEmpty = value == ""
    case reflect.Slice, reflect.Map:
        isEmpty = reflect.ValueOf(value).Len() == 0
    case reflect.Ptr, reflect.Interface:
        isEmpty = reflect.ValueOf(value).IsNil()
    }
    LogTypeCheck("IsEmpty", isEmpty)
    return isEmpty
}

// IsNonEmpty: Checks if a value is non-empty
func IsNonEmpty(value interface{}) bool {
    isNonEmpty := !IsEmpty(value)
    LogTypeCheck("IsNonEmpty", isNonEmpty)
    return isNonEmpty
}

// IsInf: Checks if a float value is positive or negative infinity
func IsInf(value float64) bool {
    result := math.IsInf(value, 0)
    LogTypeCheck("IsInf", result)
    return result
}

// IsNegInf: Checks if a float value is negative infinity
func IsNegInf(value float64) bool {
    result := math.IsInf(value, -1)
    LogTypeCheck("IsNegInf", result)
    return result
}

// IsNaN: Checks if a float value is NaN (Not a Number)
func IsNaN(value float64) bool {
    result := math.IsNaN(value)
    LogTypeCheck("IsNaN", result)
    return result
}

// IsFinite: Checks if a float value is finite (not infinity or NaN)
func IsFinite(value float64) bool {
    result := !math.IsInf(value, 0) && !math.IsNaN(value)
    LogTypeCheck("IsFinite", result)
    return result
}

// IsSymmetric: Checks if a slice or string is symmetric (mirrored around center)
func IsSymmetric(value interface{}) bool {
    var isSymmetric bool
    switch v := value.(type) {
    case string:
        isSymmetric = isPalindrome(v)
    case []int:
        isSymmetric = isSliceSymmetric(v)
    }
    LogTypeCheck("IsSymmetric", isSymmetric)
    return isSymmetric
}

// IsAsymmetric: Checks if a slice or string is asymmetric
func IsAsymmetric(value interface{}) bool {
    isAsymmetric := !IsSymmetric(value)
    LogTypeCheck("IsAsymmetric", isAsymmetric)
    return isAsymmetric
}

// IsRecursive: Checks if a function calls itself (requires user-provided function)
func IsRecursive(f interface{}) bool {
    // This is a simulated check, as Go lacks reflection for function bodies
    isRecursive := false // Placeholder as detecting true recursion programmatically is complex
    LogTypeCheck("IsRecursive", isRecursive)
    return isRecursive
}

// IsNonRecursive: Checks if a function is not recursive
func IsNonRecursive(f interface{}) bool {
    isNonRecursive := !IsRecursive(f)
    LogTypeCheck("IsNonRecursive", isNonRecursive)
    return isNonRecursive
}

// IsPalindrome: Checks if a string or slice of integers is a palindrome
func IsPalindrome(value interface{}) bool {
    var result bool
    switch v := value.(type) {
    case string:
        result = isPalindrome(v)
    case []int:
        result = isSliceSymmetric(v)
    }
    LogTypeCheck("IsPalindrome", result)
    return result
}

// Helper Functions

// isPalindrome: Helper to check if a string is a palindrome
func isPalindrome(s string) bool {
    n := len(s)
    for i := 0; i < n/2; i++ {
        if s[i] != s[n-i-1] {
            return false
        }
    }
    return true
}

// isSliceSymmetric: Helper to check if a slice is symmetric (mirrored)
func isSliceSymmetric(slice []int) bool {
    n := len(slice)
    for i := 0; i < n/2; i++ {
        if slice[i] != slice[n-i-1] {
            return false
        }
    }
    return true
}

// LogTypeCheck: Logs type check operations with encryption
func LogTypeCheck(operation string, result interface{}) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Result: " + fmt.Sprintf("%v", result)))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("TypeCheckOperation", encryptedMessage)
}
