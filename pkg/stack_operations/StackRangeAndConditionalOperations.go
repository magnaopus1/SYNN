package stack_operations

import (
    "errors"
    "fmt"
    "synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// StackRange retrieves all values within a specified index range from the stack.
func StackRange(stackID string, startIndex, endIndex int) ([]common.StackValue, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return nil, fmt.Errorf("stack not found: %v", err)
    }

    if startIndex < 0 || endIndex >= len(stack.Values) || startIndex > endIndex {
        return nil, errors.New("invalid range specified")
    }

    rangeValues := stack.Values[startIndex : endIndex+1]
    decryptedValues := make([]StackValue, len(rangeValues))

    for i, val := range rangeValues {
        decryptedData, err := encryption.DecryptData(val.Data)
        if err != nil {
            return nil, fmt.Errorf("failed to decrypt stack value: %v", err)
        }
        decryptedValues[i] = common.StackValue{Data: decryptedData, Value: val.Value}
    }

    return decryptedValues, nil
}

// SetStackValueConditionally sets a specific value at positions in the stack where a condition is met.
func SetStackValueConditionally(stackID string, condition func(common.StackValue) bool, newValue common.StackValue) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    encryptedData, err := encryption.EncryptData(newValue.Data)
    if err != nil {
        return fmt.Errorf("failed to encrypt new value: %v", err)
    }
    newValue.Data = encryptedData

    for i, val := range stack.Values {
        if condition(val) {
            stack.Values[i] = newValue
        }
    }
    return ledger.UpdateStack(stack)
}

// GetStackSubsection extracts a subsection of the stack based on a specified condition.
func GetStackSubsection(stackID string, condition func(common.StackValue) bool) ([]common.StackValue, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return nil, fmt.Errorf("stack not found: %v", err)
    }

    subsection := []common.StackValue{}
    for _, val := range stack.Values {
        if condition(val) {
            decryptedData, err := encryption.DecryptData(val.Data)
            if err != nil {
                return nil, fmt.Errorf("failed to decrypt stack value: %v", err)
            }
            subsection = append(subsection, common.StackValue{Data: decryptedData, Value: val.Value})
        }
    }
    return subsection, nil
}

// CountStackMatches counts the number of values in the stack that match a given condition.
func CountStackMatches(stackID string, condition func(common.StackValue) bool) (int, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return 0, fmt.Errorf("stack not found: %v", err)
    }

    count := 0
    for _, val := range stack.Values {
        if condition(val) {
            count++
        }
    }
    return count, nil
}

// CheckStackUniqueness verifies if all values in the stack are unique.
func CheckStackUniqueness(stackID string) (bool, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return false, fmt.Errorf("stack not found: %v", err)
    }

    valueSet := make(map[string]struct{})
    for _, val := range stack.Values {
        decryptedData, err := encryption.DecryptData(val.Data)
        if err != nil {
            return false, fmt.Errorf("failed to decrypt stack value: %v", err)
        }

        if _, exists := valueSet[decryptedData]; exists {
            return false, nil // Duplicate found
        }
        valueSet[decryptedData] = struct{}{}
    }
    return true, nil
}

// CheckStackBalance checks if the stack has balanced pairs based on a given pair-matching function.
func CheckStackBalance(stackID string, isPair func(a, b common.StackValue) bool) (bool, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return false, fmt.Errorf("stack not found: %v", err)
    }

    if len(stack.Values)%2 != 0 {
        return false, nil // Odd number of elements, cannot be balanced
    }

    for i := 0; i < len(stack.Values); i += 2 {
        if !isPair(stack.Values[i], stack.Values[i+1]) {
            return false // Pair not matching
        }
    }
    return true, nil
}
