package utility

import (
	"math"
	"strings"
	"synnergy_network/pkg/common"
)

// IsNotNull: Checks if a value is not nil
func IsNotNull(value interface{}) bool {
    result := value != nil
    LogDiagnostic("IsNotNull", "Checked for non-null value: result = " + boolToString(result))
    return result
}

// IsZero: Checks if a numeric value is zero
func IsZero(value float64) bool {
    result := value == 0
    LogDiagnostic("IsZero", "Checked if value is zero: result = " + boolToString(result))
    return result
}

// IsNonZero: Checks if a numeric value is non-zero
func IsNonZero(value float64) bool {
    result := value != 0
    LogDiagnostic("IsNonZero", "Checked if value is non-zero: result = " + boolToString(result))
    return result
}

// ApproxEqual: Checks if two float64 values are approximately equal within a tolerance
func ApproxEqual(a, b, tolerance float64) bool {
    result := math.Abs(a-b) <= tolerance
    LogDiagnostic("ApproxEqual", "Checked approximate equality: result = " + boolToString(result))
    return result
}

// ApproxNotEqual: Checks if two float64 values are not approximately equal within a tolerance
func ApproxNotEqual(a, b, tolerance float64) bool {
    result := math.Abs(a-b) > tolerance
    LogDiagnostic("ApproxNotEqual", "Checked approximate non-equality: result = " + boolToString(result))
    return result
}

// StringEqual: Checks if two strings are equal
func StringEqual(a, b string) bool {
    result := a == b
    LogDiagnostic("StringEqual", "Checked string equality: result = " + boolToString(result))
    return result
}

// StringNotEqual: Checks if two strings are not equal
func StringNotEqual(a, b string) bool {
    result := a != b
    LogDiagnostic("StringNotEqual", "Checked string non-equality: result = " + boolToString(result))
    return result
}

// StringLessThan: Checks if string a is lexicographically less than string b
func StringLessThan(a, b string) bool {
    result := strings.Compare(a, b) < 0
    LogDiagnostic("StringLessThan", "Checked if string a < string b: result = " + boolToString(result))
    return result
}

// StringGreaterThan: Checks if string a is lexicographically greater than string b
func StringGreaterThan(a, b string) bool {
    result := strings.Compare(a, b) > 0
    LogDiagnostic("StringGreaterThan", "Checked if string a > string b: result = " + boolToString(result))
    return result
}

// StringLessThanOrEqual: Checks if string a is lexicographically less than or equal to string b
func StringLessThanOrEqual(a, b string) bool {
    result := strings.Compare(a, b) <= 0
    LogDiagnostic("StringLessThanOrEqual", "Checked if string a <= string b: result = " + boolToString(result))
    return result
}

// StringGreaterThanOrEqual: Checks if string a is lexicographically greater than or equal to string b
func StringGreaterThanOrEqual(a, b string) bool {
    result := strings.Compare(a, b) >= 0
    LogDiagnostic("StringGreaterThanOrEqual", "Checked if string a >= string b: result = " + boolToString(result))
    return result
}

// IsSubset: Checks if setA is a subset of setB
func IsSubset(setA, setB []string) bool {
    result := isSubsetHelper(setA, setB, false)
    LogDiagnostic("IsSubset", "Checked if setA is a subset of setB: result = " + boolToString(result))
    return result
}

// IsSuperset: Checks if setA is a superset of setB
func IsSuperset(setA, setB []string) bool {
    result := isSubsetHelper(setB, setA, false)
    LogDiagnostic("IsSuperset", "Checked if setA is a superset of setB: result = " + boolToString(result))
    return result
}

// IsStrictSubset: Checks if setA is a strict subset of setB
func IsStrictSubset(setA, setB []string) bool {
    result := isSubsetHelper(setA, setB, true)
    LogDiagnostic("IsStrictSubset", "Checked if setA is a strict subset of setB: result = " + boolToString(result))
    return result
}

// IsStrictSuperset: Checks if setA is a strict superset of setB
func IsStrictSuperset(setA, setB []string) bool {
    result := isSubsetHelper(setB, setA, true)
    LogDiagnostic("IsStrictSuperset", "Checked if setA is a strict superset of setB: result = " + boolToString(result))
    return result
}

// Helper Functions

// isSubsetHelper: Checks if setA is a subset of setB; strict specifies if strict comparison is required
func isSubsetHelper(setA, setB []string, strict bool) bool {
    setBMap := make(map[string]struct{})
    for _, item := range setB {
        setBMap[item] = struct{}{}
    }

    for _, item := range setA {
        if _, exists := setBMap[item]; !exists {
            return false
        }
    }
    
    if strict && len(setA) == len(setB) {
        return false
    }
    return true
}

// LogDiagnostic: Helper function to log encrypted diagnostic messages
func LogDiagnostic(context, message string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogDiagnostic(context, encryptedMessage)
}

// boolToString: Converts a boolean value to a string ("true" or "false")
func boolToString(value bool) string {
    if value {
        return "true"
    }
    return "false"
}
