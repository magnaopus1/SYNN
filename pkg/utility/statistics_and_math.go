package utility

import (
	"errors"
	"fmt"
	"math"
	"synnergy_network/pkg/common"
)

// StackStandardDeviation: Calculates the standard deviation of the stack values
func (s *Stack) StackStandardDeviation() (float64, error) {
    variance, err := s.StackVariance()
    if err != nil {
        return 0, err
    }
    stdDev := math.Sqrt(variance)
    LogMathOperation("StackStandardDeviation", stdDev)
    return stdDev, nil
}

// StackAverage: Calculates the average of the stack values
func (s *Stack) StackAverage() (float64, error) {
    s.lock.Lock()
    defer s.lock.Unlock()
    if len(s.items) == 0 {
        return 0, errors.New("stack is empty")
    }
    sum := 0.0
    for _, v := range s.items {
        sum += v
    }
    avg := sum / float64(len(s.items))
    LogMathOperation("StackAverage", avg)
    return avg, nil
}

// StackVariance: Calculates the variance of the stack values
func (s *Stack) StackVariance() (float64, error) {
    avg, err := s.StackAverage()
    if err != nil {
        return 0, err
    }
    varianceSum := 0.0
    for _, v := range s.items {
        varianceSum += math.Pow(v - avg, 2)
    }
    variance := varianceSum / float64(len(s.items))
    LogMathOperation("StackVariance", variance)
    return variance, nil
}

// GetStackMin: Retrieves the minimum value in the stack
func (s *Stack) GetStackMin() (float64, error) {
    s.lock.Lock()
    defer s.lock.Unlock()
    if len(s.items) == 0 {
        return 0, errors.New("stack is empty")
    }
    min := s.items[0]
    for _, v := range s.items {
        if v < min {
            min = v
        }
    }
    LogMathOperation("GetStackMin", min)
    return min, nil
}

// GetStackMax: Retrieves the maximum value in the stack
func (s *Stack) GetStackMax() (float64, error) {
    s.lock.Lock()
    defer s.lock.Unlock()
    if len(s.items) == 0 {
        return 0, errors.New("stack is empty")
    }
    max := s.items[0]
    for _, v := range s.items {
        if v > max {
            max = v
        }
    }
    LogMathOperation("GetStackMax", max)
    return max, nil
}

// StackSum: Calculates the sum of all values in the stack
func (s *Stack) StackSum() float64 {
    s.lock.Lock()
    defer s.lock.Unlock()
    sum := 0.0
    for _, v := range s.items {
        sum += v
    }
    LogMathOperation("StackSum", sum)
    return sum
}

// StackProduct: Calculates the product of all values in the stack
func (s *Stack) StackProduct() float64 {
    s.lock.Lock()
    defer s.lock.Unlock()
    product := 1.0
    for _, v := range s.items {
        product *= v
    }
    LogMathOperation("StackProduct", product)
    return product
}

// CheckStackContains: Checks if the stack contains a specified value
func (s *Stack) CheckStackContains(value float64) bool {
    s.lock.Lock()
    defer s.lock.Unlock()
    for _, v := range s.items {
        if v == value {
            LogMathOperation("CheckStackContains", true)
            return true
        }
    }
    LogMathOperation("CheckStackContains", false)
    return false
}

// RemoveStackDuplicates: Removes duplicate values from the stack
func (s *Stack) RemoveStackDuplicates() {
    s.lock.Lock()
    defer s.lock.Unlock()
    uniqueItems := []float64{}
    seen := make(map[float64]struct{})
    for _, v := range s.items {
        if _, exists := seen[v]; !exists {
            seen[v] = struct{}{}
            uniqueItems = append(uniqueItems, v)
        }
    }
    s.items = uniqueItems
    LogMathOperation("RemoveStackDuplicates", "duplicates removed")
}

// CheckStackEmptiness: Checks if the stack is empty
func (s *Stack) CheckStackEmptiness() bool {
    s.lock.Lock()
    defer s.lock.Unlock()
    isEmpty := len(s.items) == 0
    LogMathOperation("CheckStackEmptiness", isEmpty)
    return isEmpty
}

// StackLastAccessedIndex: Retrieves the last accessed index in the stack
func (s *Stack) StackLastAccessedIndex() int {
    s.lock.Lock()
    defer s.lock.Unlock()
    LogMathOperation("StackLastAccessedIndex", s.lastAccessedIndex)
    return s.lastAccessedIndex
}

// LogStackOperations: Logs details of specific stack operations
func (s *Stack) LogStackOperations(operation string, details interface{}) error {
    message := "Operation: " + operation + " - Details: " + fmt.Sprintf("%v", details)
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("StackOperation", encryptedMessage)
}

// ValidateStackBounds: Validates that an index is within the bounds of the stack
func (s *Stack) ValidateStackBounds(index int) error {
    s.lock.Lock()
    defer s.lock.Unlock()
    if index < 0 || index >= len(s.items) {
        LogStackOverflow("ValidateStackBounds", "index out of bounds")
        return errors.New("index out of bounds")
    }
    LogMathOperation("ValidateStackBounds", "index within bounds")
    return nil
}

// LogStackOverflow: Logs a stack overflow event
func LogStackOverflow(context, message string) error {
    return logStackEvent("StackOverflow", context, message)
}

// LogStackUnderflow: Logs a stack underflow event
func LogStackUnderflow(context, message string) error {
    return logStackEvent("StackUnderflow", context, message)
}

// Helper Functions

// LogMathOperation: Logs mathematical operations securely
func LogMathOperation(operation string, result interface{}) error {
    message := "Operation: " + operation + " - Result: " + fmt.Sprintf("%v", result)
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogDiagnostic("MathOperation", encryptedMessage)
}

// logStackEvent: Logs stack overflow/underflow events securely
func logStackEvent(eventType, context, message string) error {
    logMessage := "Event: " + eventType + " - Context: " + context + " - Message: " + message
    encryptedMessage, err := encryption.Encrypt([]byte(logMessage))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent(eventType, encryptedMessage)
}
