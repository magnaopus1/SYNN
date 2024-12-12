package utility

import (
	"errors"
	"math"
	"synnergy_network/pkg/common"
)

// StackRange: Returns the range (max - min) of values in the stack
func (s *Stack) StackRange() (float64, error) {
    s.lock.Lock()
    defer s.lock.Unlock()
    if len(s.items) == 0 {
        return 0, errors.New("stack is empty")
    }
    min, max := s.items[0], s.items[0]
    for _, v := range s.items {
        if v < min {
            min = v
        } else if v > max {
            max = v
        }
    }
    rangeValue := max - min
    LogStackUtility("StackRange", rangeValue)
    return rangeValue, nil
}

// SetStackValueConditionally: Sets a value in the stack if it meets a condition
func (s *Stack) SetStackValueConditionally(index int, value float64, condition func(float64) bool) error {
    s.lock.Lock()
    defer s.lock.Unlock()
    if index < 0 || index >= len(s.items) {
        return errors.New("index out of range")
    }
    if condition(s.items[index]) {
        s.items[index] = value
        LogStackUtility("SetStackValueConditionally", value)
    }
    return nil
}

// GetStackSubsection: Returns a subsection of the stack based on start and end indices
func (s *Stack) GetStackSubsection(start, end int) ([]float64, error) {
    s.lock.Lock()
    defer s.lock.Unlock()
    if start < 0 || end >= len(s.items) || start > end {
        return nil, errors.New("invalid range")
    }
    subsection := s.items[start : end+1]
    LogStackUtility("GetStackSubsection", subsection)
    return subsection, nil
}

// CountStackMatches: Counts occurrences of values in the stack that match a specified value
func (s *Stack) CountStackMatches(value float64) int {
    s.lock.Lock()
    defer s.lock.Unlock()
    count := 0
    for _, v := range s.items {
        if v == value {
            count++
        }
    }
    LogStackUtility("CountStackMatches", count)
    return count
}

// CheckStackUniqueness: Checks if all values in the stack are unique
func (s *Stack) CheckStackUniqueness() bool {
    s.lock.Lock()
    defer s.lock.Unlock()
    seen := make(map[float64]struct{})
    for _, v := range s.items {
        if _, exists := seen[v]; exists {
            LogStackUtility("CheckStackUniqueness", false)
            return false
        }
        seen[v] = struct{}{}
    }
    LogStackUtility("CheckStackUniqueness", true)
    return true
}

// StackModuloOperation: Applies modulo operation to each element in the stack
func (s *Stack) StackModuloOperation(mod float64) {
    s.lock.Lock()
    defer s.lock.Unlock()
    for i, v := range s.items {
        s.items[i] = math.Mod(v, mod)
    }
    LogStackUtility("StackModuloOperation", "applied modulo operation")
}

// StackAbsoluteValues: Replaces each element in the stack with its absolute value
func (s *Stack) StackAbsoluteValues() {
    s.lock.Lock()
    defer s.lock.Unlock()
    for i, v := range s.items {
        s.items[i] = math.Abs(v)
    }
    LogStackUtility("StackAbsoluteValues", "applied absolute values")
}

// StackPowerValues: Raises each element in the stack to a specified power
func (s *Stack) StackPowerValues(power float64) {
    s.lock.Lock()
    defer s.lock.Unlock()
    for i, v := range s.items {
        s.items[i] = math.Pow(v, power)
    }
    LogStackUtility("StackPowerValues", "applied power values")
}

// StackExponentialValues: Replaces each element in the stack with its exponential value
func (s *Stack) StackExponentialValues() {
    s.lock.Lock()
    defer s.lock.Unlock()
    for i, v := range s.items {
        s.items[i] = math.Exp(v)
    }
    LogStackUtility("StackExponentialValues", "applied exponential values")
}

// StackFloorValues: Replaces each element in the stack with its floor value
func (s *Stack) StackFloorValues() {
    s.lock.Lock()
    defer s.lock.Unlock()
    for i, v := range s.items {
        s.items[i] = math.Floor(v)
    }
    LogStackUtility("StackFloorValues", "applied floor values")
}

// StackCeilingValues: Replaces each element in the stack with its ceiling value
func (s *Stack) StackCeilingValues() {
    s.lock.Lock()
    defer s.lock.Unlock()
    for i, v := range s.items {
        s.items[i] = math.Ceil(v)
    }
    LogStackUtility("StackCeilingValues", "applied ceiling values")
}

// CheckStackBalance: Checks if the sum of stack values is zero
func (s *Stack) CheckStackBalance() bool {
    s.lock.Lock()
    defer s.lock.Unlock()
    sum := 0.0
    for _, v := range s.items {
        sum += v
    }
    result := sum == 0
    LogStackUtility("CheckStackBalance", result)
    return result
}

// FindStackMedian: Finds the median value of the stack
func (s *Stack) FindStackMedian() (float64, error) {
    s.lock.Lock()
    defer s.lock.Unlock()
    if len(s.items) == 0 {
        return 0, errors.New("stack is empty")
    }
    sortedItems := append([]float64(nil), s.items...)
    sort.Float64s(sortedItems)
    median := 0.0
    if len(sortedItems)%2 == 1 {
        median = sortedItems[len(sortedItems)/2]
    } else {
        mid := len(sortedItems) / 2
        median = (sortedItems[mid-1] + sortedItems[mid]) / 2
    }
    LogStackUtility("FindStackMedian", median)
    return median, nil
}

// NormalizeStack: Normalizes stack values to a 0-1 range
func (s *Stack) NormalizeStack() error {
    s.lock.Lock()
    defer s.lock.Unlock()
    if len(s.items) == 0 {
        return errors.New("stack is empty")
    }
    min, max := s.items[0], s.items[0]
    for _, v := range s.items {
        if v < min {
            min = v
        } else if v > max {
            max = v
        }
    }
    rangeValue := max - min
    if rangeValue == 0 {
        return errors.New("cannot normalize stack with identical values")
    }
    for i, v := range s.items {
        s.items[i] = (v - min) / rangeValue
    }
    LogStackUtility("NormalizeStack", "stack normalized")
    return nil
}

// StackRounding: Rounds each value in the stack to the nearest integer
func (s *Stack) StackRounding() {
    s.lock.Lock()
    defer s.lock.Unlock()
    for i, v := range s.items {
        s.items[i] = math.Round(v)
    }
    LogStackUtility("StackRounding", "applied rounding")
}

// Helper Functions

// LogStackUtility: Helper function to log encrypted stack utility operations
func LogStackUtility(operation string, details interface{}) error {
    message := "Operation: " + operation + " - Details: " + fmt.Sprintf("%v", details)
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("StackUtility", encryptedMessage)
}
