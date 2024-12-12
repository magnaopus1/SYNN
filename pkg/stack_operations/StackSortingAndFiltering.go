package stack_operations

import (
    "errors"
    "fmt"
    "sort"
"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// SortStack sorts the stack in ascending or descending order.
func SortStack(stackID string, ascending bool) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    if ascending {
        sort.Slice(stack.Values, func(i, j int) bool {
            return stack.Values[i].Value < stack.Values[j].Value
        })
    } else {
        sort.Slice(stack.Values, func(i, j int) bool {
            return stack.Values[i].Value > stack.Values[j].Value
        })
    }
    return ledger.UpdateStack(stack)
}

// FilterStack filters the stack values based on a given condition.
func FilterStack(stackID string, condition func(common.StackValue) bool) ([]common.StackValue, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return nil, fmt.Errorf("stack not found: %v", err)
    }

    filteredValues := []common.StackValue{}
    for _, val := range stack.Values {
        if condition(val) {
            decryptedData, err := encryption.DecryptData(val.Data)
            if err != nil {
                return nil, fmt.Errorf("failed to decrypt stack value: %v", err)
            }
            val.Data = decryptedData
            filteredValues = append(filteredValues, val)
        }
    }
    return filteredValues, nil
}

// MapStackValues applies a mapping function to each stack value.
func MapStackValues(stackID string, mapper func(common.StackValue) common.StackValue) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    for i, val := range stack.Values {
        mappedValue := mapper(val)
        encryptedData, err := encryption.EncryptData(mappedValue.Data)
        if err != nil {
            return fmt.Errorf("failed to encrypt mapped value: %v", err)
        }
        stack.Values[i].Data = encryptedData
        stack.Values[i].Value = mappedValue.Value
    }
    return ledger.UpdateStack(stack)
}

// AccumulateStack applies an accumulator function to aggregate stack values into a single result.
func AccumulateStack(stackID string, accumulator func(float64, float64) float64, initialValue float64) (float64, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return 0, fmt.Errorf("stack not found: %v", err)
    }

    result := initialValue
    for _, val := range stack.Values {
        result = accumulator(result, val.Value)
    }
    return result, nil
}

// IsStackSorted checks if the stack is sorted in ascending or descending order.
func IsStackSorted(stackID string, ascending bool) (bool, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return false, fmt.Errorf("stack not found: %v", err)
    }

    for i := 1; i < len(stack.Values); i++ {
        if ascending && stack.Values[i-1].Value > stack.Values[i].Value {
            return false, nil
        } else if !ascending && stack.Values[i-1].Value < stack.Values[i].Value {
            return false, nil
        }
    }
    return true, nil
}

// RemoveStackDuplicates removes duplicate values from the stack, keeping the first occurrence.
func RemoveStackDuplicates(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    uniqueValues := make(map[string]struct{})
    uniqueStack := []common.StackValue{}

    for _, val := range stack.Values {
        decryptedData, err := encryption.DecryptData(val.Data)
        if err != nil {
            return fmt.Errorf("failed to decrypt stack value: %v", err)
        }
        if _, exists := uniqueValues[decryptedData]; !exists {
            uniqueValues[decryptedData] = struct{}{}
            val.Data = decryptedData
            uniqueStack = append(uniqueStack, val)
        }
    }

    stack.Values = uniqueStack
    return ledger.UpdateStack(stack)
}

// NormalizeStack scales all values on the stack to be between 0 and 1.
func NormalizeStack(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    if len(stack.Values) == 0 {
        return errors.New("stack is empty")
    }

    min, max := stack.Values[0].Value, stack.Values[0].Value
    for _, val := range stack.Values {
        if val.Value < min {
            min = val.Value
        }
        if val.Value > max {
            max = val.Value
        }
    }

    if min == max {
        return errors.New("cannot normalize stack with all identical values")
    }

    for i := range stack.Values {
        stack.Values[i].Value = (stack.Values[i].Value - min) / (max - min)
    }
    return ledger.UpdateStack(stack)
}
