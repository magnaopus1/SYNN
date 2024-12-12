package utility

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"synnergy_network/pkg/common"
)

// IsAnagram: Checks if two strings are anagrams of each other
func IsAnagram(str1, str2 string) bool {
    normalizedStr1 := strings.ReplaceAll(str1, " ", "")
    normalizedStr2 := strings.ReplaceAll(str2, " ", "")
    if len(normalizedStr1) != len(normalizedStr2) {
        LogSequenceOperation("IsAnagram", false)
        return false
    }
    sortedStr1 := sortString(normalizedStr1)
    sortedStr2 := sortString(normalizedStr2)
    result := sortedStr1 == sortedStr2
    LogSequenceOperation("IsAnagram", result)
    return result
}

// CompareLength: Compares the length of two strings or sequences
func CompareLength(seq1, seq2 string) int {
    result := len(seq1) - len(seq2)
    LogSequenceOperation("CompareLength", result)
    return result
}

// CompareSum: Compares the sum of two numeric arrays
func CompareSum(arr1, arr2 []float64) float64 {
    sum1, sum2 := arraySum(arr1), arraySum(arr2)
    result := sum1 - sum2
    LogSequenceOperation("CompareSum", result)
    return result
}

// CompareProduct: Compares the product of two numeric arrays
func CompareProduct(arr1, arr2 []float64) float64 {
    product1, product2 := arrayProduct(arr1), arrayProduct(arr2)
    result := product1 - product2
    LogSequenceOperation("CompareProduct", result)
    return result
}

// MatchesPattern: Checks if a string matches a specified pattern
func MatchesPattern(str, pattern string) bool {
    matched := strings.Contains(str, pattern)
    LogSequenceOperation("MatchesPattern", matched)
    return matched
}

// NotMatchesPattern: Checks if a string does not match a specified pattern
func NotMatchesPattern(str, pattern string) bool {
    notMatched := !strings.Contains(str, pattern)
    LogSequenceOperation("NotMatchesPattern", notMatched)
    return notMatched
}

// IsPositiveInfinity: Checks if a value is positive infinity
func IsPositiveInfinity(value float64) bool {
    result := math.IsInf(value, 1)
    LogSequenceOperation("IsPositiveInfinity", result)
    return result
}

// IsNegativeInfinity: Checks if a value is negative infinity
func IsNegativeInfinity(value float64) bool {
    result := math.IsInf(value, -1)
    LogSequenceOperation("IsNegativeInfinity", result)
    return result
}

// IsUnique: Checks if all elements in a string or array are unique
func IsUnique(values []string) bool {
    unique := make(map[string]struct{})
    for _, v := range values {
        if _, exists := unique[v]; exists {
            LogSequenceOperation("IsUnique", false)
            return false
        }
        unique[v] = struct{}{}
    }
    LogSequenceOperation("IsUnique", true)
    return true
}

// IsDuplicate: Checks if a value is duplicated within a string or array
func IsDuplicate(values []string) bool {
    seen := make(map[string]struct{})
    for _, v := range values {
        if _, exists := seen[v]; exists {
            LogSequenceOperation("IsDuplicate", true)
            return true
        }
        seen[v] = struct{}{}
    }
    LogSequenceOperation("IsDuplicate", false)
    return false
}

// CompareType: Compares the type of two values and returns true if they match
func CompareType(val1, val2 interface{}) bool {
    result := getType(val1) == getType(val2)
    LogSequenceOperation("CompareType", result)
    return result
}

// IsFloatEqual: Checks if two float64 values are equal within a specified tolerance
func IsFloatEqual(a, b, tolerance float64) bool {
    result := math.Abs(a-b) <= tolerance
    LogSequenceOperation("IsFloatEqual", result)
    return result
}

// IsFloatNotEqual: Checks if two float64 values are not equal within a specified tolerance
func IsFloatNotEqual(a, b, tolerance float64) bool {
    result := math.Abs(a-b) > tolerance
    LogSequenceOperation("IsFloatNotEqual", result)
    return result
}

// ArrayContainsElement: Checks if an array contains a specific element
func ArrayContainsElement(arr []string, element string) bool {
    for _, v := range arr {
        if v == element {
            LogSequenceOperation("ArrayContainsElement", true)
            return true
        }
    }
    LogSequenceOperation("ArrayContainsElement", false)
    return false
}

// ArrayDoesNotContainElement: Checks if an array does not contain a specific element
func ArrayDoesNotContainElement(arr []string, element string) bool {
    result := !ArrayContainsElement(arr, element)
    LogSequenceOperation("ArrayDoesNotContainElement", result)
    return result
}

// Helper Functions

// sortString: Sorts a string's characters in alphabetical order
func sortString(str string) string {
    sortedRunes := []rune(str)
    sort.Slice(sortedRunes, func(i, j int) bool { return sortedRunes[i] < sortedRunes[j] })
    return string(sortedRunes)
}

// arraySum: Calculates the sum of all elements in a float64 array
func arraySum(arr []float64) float64 {
    sum := 0.0
    for _, v := range arr {
        sum += v
    }
    return sum
}

// arrayProduct: Calculates the product of all elements in a float64 array
func arrayProduct(arr []float64) float64 {
    product := 1.0
    for _, v := range arr {
        product *= v
    }
    return product
}



// LogSequenceOperation: Logs a string or sequence operation securely
func LogSequenceOperation(operation string, result interface{}) error {
    message := "Operation: " + operation + " - Result: " + fmt.Sprintf("%v", result)
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogDiagnostic("SequenceOperation", encryptedMessage)
}
