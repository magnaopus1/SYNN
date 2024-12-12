package stack_operations

import (
    "errors"
    "fmt"
"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)
)

// PushMultiple pushes multiple values onto the stack.
func PushMultiple(stackID string, values []common.StackValue) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    for _, value := range values {
        encryptedData, err := encryption.EncryptData(value.Data)
        if err != nil {
            return fmt.Errorf("failed to encrypt value: %v", err)
        }
        value.Data = encryptedData
        stack.Values = append(stack.Values, value)
    }
    return ledger.UpdateStack(stack)
}

// PopMultiple removes and returns multiple values from the top of the stack.
func PopMultiple(stackID string, count int) ([]common.StackValue, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return nil, fmt.Errorf("stack not found: %v", err)
    }

    if count > len(stack.Values) {
        return nil, errors.New("not enough values on stack")
    }

    poppedValues := stack.Values[len(stack.Values)-count:]
    stack.Values = stack.Values[:len(stack.Values)-count]

    for i, val := range poppedValues {
        decryptedData, err := encryption.DecryptData(val.Data)
        if err != nil {
            return nil, fmt.Errorf("failed to decrypt value: %v", err)
        }
        poppedValues[i].Data = decryptedData
    }

    return poppedValues, ledger.UpdateStack(stack)
}

// SplitStack splits the stack into two at a specified index.
func SplitStack(stackID string, index int) (string, string, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return "", "", fmt.Errorf("stack not found: %v", err)
    }

    if index < 0 || index > len(stack.Values) {
        return "", "", errors.New("index out of bounds")
    }

    stack1Values := stack.Values[:index]
    stack2Values := stack.Values[index:]

    stack1ID := ledger.CreateStack(stack1Values)
    stack2ID := ledger.CreateStack(stack2Values)

    return stack1ID, stack2ID, nil
}

// MergeStacks merges two stacks into one.
func MergeStacks(stackID1, stackID2 string) (string, error) {
    stack1, err := ledger.GetStack(stackID1)
    if err != nil {
        return "", fmt.Errorf("stack 1 not found: %v", err)
    }
    stack2, err := ledger.GetStack(stackID2)
    if err != nil {
        return "", fmt.Errorf("stack 2 not found: %v", err)
    }

    mergedValues := append(stack1.Values, stack2.Values...)
    mergedStackID := ledger.CreateStack(mergedValues)

    return mergedStackID, nil
}

// CloneStack creates a copy of the specified stack.
func CloneStack(stackID string) (string, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return "", fmt.Errorf("stack not found: %v", err)
    }

    cloneValues := make([]common.StackValue, len(stack.Values))
    copy(cloneValues, stack.Values)

    cloneStackID := ledger.CreateStack(cloneValues)
    return cloneStackID, nil
}

// StackDifference returns the difference of values in the two stacks.
func StackDifference(stackID1, stackID2 string) ([]common.StackValue, error) {
    stack1, err := ledger.GetStack(stackID1)
    if err != nil {
        return nil, fmt.Errorf("stack 1 not found: %v", err)
    }
    stack2, err := ledger.GetStack(stackID2)
    if err != nil {
        return nil, fmt.Errorf("stack 2 not found: %v", err)
    }

    difference := []common.StackValue{}
    stack2Values := make(map[string]struct{})
    for _, v := range stack2.Values {
        stack2Values[v.Data] = struct{}{}
    }

    for _, v := range stack1.Values {
        if _, exists := stack2Values[v.Data]; !exists {
            difference = append(difference, v)
        }
    }

    return difference, nil
}

// StackUnion combines all unique values from both stacks.
func StackUnion(stackID1, stackID2 string) ([]common.StackValue, error) {
    stack1, err := ledger.GetStack(stackID1)
    if err != nil {
        return nil, fmt.Errorf("stack 1 not found: %v", err)
    }
    stack2, err := ledger.GetStack(stackID2)
    if err != nil {
        return nil, fmt.Errorf("stack 2 not found: %v", err)
    }

    unionMap := make(map[string]common.StackValue)
    for _, v := range stack1.Values {
        unionMap[v.Data] = v
    }
    for _, v := range stack2.Values {
        unionMap[v.Data] = v
    }

    union := []common.StackValue{}
    for _, value := range unionMap {
        union = append(union, value)
    }

    return union, nil
}

// StackIntersection returns the common values between two stacks.
func StackIntersection(stackID1, stackID2 string) ([]common.StackValue, error) {
    stack1, err := ledger.GetStack(stackID1)
    if err != nil {
        return nil, fmt.Errorf("stack 1 not found: %v", err)
    }
    stack2, err := ledger.GetStack(stackID2)
    if err != nil {
        return nil, fmt.Errorf("stack 2 not found: %v", err)
    }

    stack2Values := make(map[string]struct{})
    for _, v := range stack2.Values {
        stack2Values[v.Data] = struct{}{}
    }

    intersection := []common.StackValue{}
    for _, v := range stack1.Values {
        if _, exists := stack2Values[v.Data]; exists {
            intersection = append(intersection, v)
        }
    }

    return intersection, nil
}

// StackExclusiveElements returns values that are only in one of the two stacks.
func StackExclusiveElements(stackID1, stackID2 string) ([]common.StackValue, error) {
    stack1, err := ledger.GetStack(stackID1)
    if err != nil {
        return nil, fmt.Errorf("stack 1 not found: %v", err)
    }
    stack2, err := ledger.GetStack(stackID2)
    if err != nil {
        return nil, fmt.Errorf("stack 2 not found: %v", err)
    }

    exclusive := []common.StackValue{}
    stack1Values := make(map[string]struct{})
    stack2Values := make(map[string]struct{})

    for _, v := range stack1.Values {
        stack1Values[v.Data] = struct{}{}
    }
    for _, v := range stack2.Values {
        stack2Values[v.Data] = struct{}{}
    }

    for _, v := range stack1.Values {
        if _, exists := stack2Values[v.Data]; !exists {
            exclusive = append(exclusive, v)
        }
    }

    for _, v := range stack2.Values {
        if _, exists := stack1Values[v.Data]; !exists {
            exclusive = append(exclusive, v)
        }
    }

    return exclusive, nil
}
