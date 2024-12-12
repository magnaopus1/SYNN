package utility

import (
	"fmt"
	"math"
	"synnergy_network/pkg/common"
)

// HasOverlap: Checks if two sets have any overlapping elements
func HasOverlap(setA, setB []string) bool {
    result := findOverlap(setA, setB) != nil
    LogSetOperation("HasOverlap", result)
    return result
}

// NoOverlap: Checks if two sets have no overlapping elements
func NoOverlap(setA, setB []string) bool {
    result := findOverlap(setA, setB) == nil
    LogSetOperation("NoOverlap", result)
    return result
}

// IsRepeatingPattern: Checks if a set has a repeating pattern
func IsRepeatingPattern(set []string) bool {
    for i := 1; i <= len(set)/2; i++ {
        if len(set)%i == 0 && isPattern(set, i) {
            LogSetOperation("IsRepeatingPattern", true)
            return true
        }
    }
    LogSetOperation("IsRepeatingPattern", false)
    return false
}

// CompareMagnitude: Compares the magnitude (length) of two sets
func CompareMagnitude(setA, setB []string) int {
    result := len(setA) - len(setB)
    LogSetOperation("CompareMagnitude", result)
    return result
}

// IsSequence: Checks if a set of numbers forms a consecutive sequence
func IsSequence(nums []int) bool {
    if len(nums) < 2 {
        LogSetOperation("IsSequence", false)
        return false
    }
    isSeq := true
    for i := 1; i < len(nums); i++ {
        if nums[i] != nums[i-1]+1 {
            isSeq = false
            break
        }
    }
    LogSetOperation("IsSequence", isSeq)
    return isSeq
}

// IsMonotonic: Checks if a sequence is monotonic (either increasing or decreasing)
func IsMonotonic(nums []int) bool {
    inc, dec := true, true
    for i := 1; i < len(nums); i++ {
        inc = inc && nums[i] >= nums[i-1]
        dec = dec && nums[i] <= nums[i-1]
    }
    result := inc || dec
    LogSetOperation("IsMonotonic", result)
    return result
}

// CompareBitCount: Compares the number of 1-bits in two integers
func CompareBitCount(a, b int) int {
    countA := bitCount(a)
    countB := bitCount(b)
    result := countA - countB
    LogSetOperation("CompareBitCount", result)
    return result
}

// HasDuplicates: Checks if a set contains duplicate elements
func HasDuplicates(set []string) bool {
    result := containsDuplicates(set)
    LogSetOperation("HasDuplicates", result)
    return result
}

// NoDuplicates: Checks if a set contains no duplicate elements
func NoDuplicates(set []string) bool {
    result := !containsDuplicates(set)
    LogSetOperation("NoDuplicates", result)
    return result
}

// HasCommonElements: Checks if two sets have common elements
func HasCommonElements(setA, setB []string) bool {
    result := HasOverlap(setA, setB)
    LogSetOperation("HasCommonElements", result)
    return result
}

// NoCommonElements: Checks if two sets have no common elements
func NoCommonElements(setA, setB []string) bool {
    result := NoOverlap(setA, setB)
    LogSetOperation("NoCommonElements", result)
    return result
}

// IsSymmetricDifference: Checks if two sets are symmetric differences
func IsSymmetricDifference(setA, setB []string) bool {
    symDiff := make(map[string]bool)
    for _, v := range setA {
        symDiff[v] = !symDiff[v]
    }
    for _, v := range setB {
        symDiff[v] = !symDiff[v]
    }
    result := true
    for _, v := range symDiff {
        if v {
            result = false
            break
        }
    }
    LogSetOperation("IsSymmetricDifference", result)
    return result
}

// IsEqualSum: Checks if the sums of two integer sets are equal
func IsEqualSum(setA, setB []int) bool {
    result := sum(setA) == sum(setB)
    LogSetOperation("IsEqualSum", result)
    return result
}

// IsNotEqualSum: Checks if the sums of two integer sets are not equal
func IsNotEqualSum(setA, setB []int) bool {
    result := sum(setA) != sum(setB)
    LogSetOperation("IsNotEqualSum", result)
    return result
}

// CompareRange: Compares the range (difference between max and min) of two sets
func CompareRange(setA, setB []int) int {
    rangeA := rangeOf(setA)
    rangeB := rangeOf(setB)
    result := rangeA - rangeB
    LogSetOperation("CompareRange", result)
    return result
}

// IsPositiveModulo: Checks if all elements in a set have a positive modulo with a given number
func IsPositiveModulo(set []int, mod int) bool {
    result := true
    for _, v := range set {
        if v%mod <= 0 {
            result = false
            break
        }
    }
    LogSetOperation("IsPositiveModulo", result)
    return result
}

// IsNegativeModulo: Checks if all elements in a set have a negative modulo with a given number
func IsNegativeModulo(set []int, mod int) bool {
    result := true
    for _, v := range set {
        if v%mod >= 0 {
            result = false
            break
        }
    }
    LogSetOperation("IsNegativeModulo", result)
    return result
}

// IsIdenticalType: Checks if all elements in two sets are of the same type
func IsIdenticalType(setA, setB []interface{}) bool {
    result := getType(setA) == getType(setB)
    LogSetOperation("IsIdenticalType", result)
    return result
}

// CompareDepth: Compares the depth (number of levels) of two nested structures
func CompareDepth(structA, structB interface{}) int {
    depthA := depth(structA)
    depthB := depth(structB)
    result := depthA - depthB
    LogSetOperation("CompareDepth", result)
    return result
}

// IsEqualMagnitude: Checks if the magnitudes of two values are equal
func IsEqualMagnitude(a, b float64) bool {
    result := math.Abs(a) == math.Abs(b)
    LogSetOperation("IsEqualMagnitude", result)
    return result
}

// Helper Functions

// findOverlap: Finds overlapping elements between two sets
func findOverlap(setA, setB []string) map[string]bool {
    overlap := make(map[string]bool)
    setBMap := make(map[string]struct{})
    for _, item := range setB {
        setBMap[item] = struct{}{}
    }
    for _, item := range setA {
        if _, exists := setBMap[item]; exists {
            overlap[item] = true
        }
    }
    return overlap
}

// isPattern: Checks if a repeating pattern exists with a given interval
func isPattern(set []string, interval int) bool {
    for i := 0; i < interval; i++ {
        for j := i; j < len(set); j += interval {
            if set[j] != set[i] {
                return false
            }
        }
    }
    return true
}

// bitCount: Counts the number of 1 bits in an integer
func bitCount(n int) int {
    count := 0
    for n != 0 {
        count++
        n &= n - 1
    }
    return count
}

// containsDuplicates: Checks if a set contains duplicate elements
func containsDuplicates(set []string) bool {
    seen := make(map[string]struct{})
    for _, item := range set {
        if _, exists := seen[item]; exists {
            return true
        }
        seen[item] = struct{}{}
    }
    return false
}

// sum: Sums up all elements in an integer set
func sum(set []int) int {
    total := 0
    for _, v := range set {
        total += v
    }
    return total
}

// rangeOf: Returns the range (max - min) of an integer set
func rangeOf(set []int) int {
    min, max := set[0], set[0]
    for _, v := range set {
        if v < min {
            min = v
        } else if v > max {
            max = v
        }
    }
    return max - min
}

// depth: Determines the depth of a nested structure
func depth(value interface{}) int {
    if _, ok := value.([]interface{}); !ok {
        return 0
    }
    maxDepth := 0
    for _, item := range value.([]interface{}) {
        d := depth(item)
        if d > maxDepth {
            maxDepth = d
        }
    }
    return maxDepth + 1
}

// getType: Returns a string representing the type of elements in a set
func getType(set []interface{}) string {
    if len(set) == 0 {
        return "empty"
    }
    return fmt.Sprintf("%T", set[0])
}

// LogSetOperation: Logs the result of a set operation securely
func LogSetOperation(operation string, result interface{}) error {
    message := fmt.Sprintf("Operation: %s - Result: %v", operation, result)
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogDiagnostic("SetOperation", encryptedMessage)
}
