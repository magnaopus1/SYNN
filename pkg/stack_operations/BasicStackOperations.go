package stack_operations

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)


// PushStack adds a new value to the top of the stack.
func PushStack(stackID string, value common.StackValue) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    encryptedValue, err := encryption.EncryptData(value.Data)
    if err != nil {
        return fmt.Errorf("failed to encrypt stack value: %v", err)
    }
    value.Data = encryptedValue

    stack.Values = append(stack.Values, value)
    return ledger.UpdateStack(stack)
}

// PopStack removes and returns the top value from the stack.
func PopStack(stackID string) (common.StackValue, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return common.StackValue{}, fmt.Errorf("stack not found: %v", err)
    }

    if len(stack.Values) == 0 {
        return common.StackValue{}, errors.New("stack is empty")
    }

    topValue := stack.Values[len(stack.Values)-1]
    stack.Values = stack.Values[:len(stack.Values)-1]
    decryptedData, err := encryption.DecryptData(topValue.Data)
    if err != nil {
        return common.StackValue{}, fmt.Errorf("failed to decrypt stack value: %v", err)
    }
    topValue.Data = decryptedData

    return topValue, ledger.UpdateStack(stack)
}

// PeekStack returns the top value from the stack without removing it.
func PeekStack(stackID string) (common.StackValue, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return common.StackValue{}, fmt.Errorf("stack not found: %v", err)
    }

    if len(stack.Values) == 0 {
        return common.StackValue{}, errors.New("stack is empty")
    }

    topValue := stack.Values[len(stack.Values)-1]
    decryptedData, err := encryption.DecryptData(topValue.Data)
    if err != nil {
        return common.StackValue{}, fmt.Errorf("failed to decrypt stack value: %v", err)
    }
    topValue.Data = decryptedData

    return topValue, nil
}

// ClearStack removes all values from the stack.
func ClearStack(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    stack.Values = nil
    return ledger.UpdateStack(stack)
}

// DuplicateStack duplicates the top value of the stack.
func DuplicateStack(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    if len(stack.Values) == 0 {
        return errors.New("stack is empty")
    }

    topValue := stack.Values[len(stack.Values)-1]
    stack.Values = append(stack.Values, topValue)
    return ledger.UpdateStack(stack)
}

// SwapStackTop swaps the top two values on the stack.
func SwapStackTop(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    if len(stack.Values) < 2 {
        return errors.New("not enough values to swap")
    }

    stack.Values[len(stack.Values)-1], stack.Values[len(stack.Values)-2] = stack.Values[len(stack.Values)-2], stack.Values[len(stack.Values)-1]
    return ledger.UpdateStack(stack)
}

// RotateStack rotates the top three values on the stack.
func RotateStack(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    if len(stack.Values) < 3 {
        return errors.New("not enough values to rotate")
    }

    stack.Values[len(stack.Values)-3], stack.Values[len(stack.Values)-2], stack.Values[len(stack.Values)-1] =
        stack.Values[len(stack.Values)-2], stack.Values[len(stack.Values)-1], stack.Values[len(stack.Values)-3]

    return ledger.UpdateStack(stack)
}

// CopyToStack copies a specific value from the stack by index to the top of the stack.
func CopyToStack(stackID string, index int) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    if index < 0 || index >= len(stack.Values) {
        return errors.New("index out of bounds")
    }

    stack.Values = append(stack.Values, stack.Values[index])
    return ledger.UpdateStack(stack)
}

// ReplaceStackValue replaces a specific value on the stack at a given index.
func ReplaceStackValue(stackID string, index int, newValue common.StackValue) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    if index < 0 || index >= len(stack.Values) {
        return errors.New("index out of bounds")
    }

    encryptedValue, err := encryption.EncryptData(newValue.Data)
    if err != nil {
        return fmt.Errorf("failed to encrypt new value: %v", err)
    }
    newValue.Data = encryptedValue
    stack.Values[index] = newValue

    return ledger.UpdateStack(stack)
}

// InsertStackValue inserts a new value at a specified position in the stack.
func InsertStackValue(stackID string, index int, value common.StackValue) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    if index < 0 || index > len(stack.Values) {
        return errors.New("index out of bounds")
    }

    encryptedValue, err := encryption.EncryptData(value.Data)
    if err != nil {
        return fmt.Errorf("failed to encrypt value: %v", err)
    }
    value.Data = encryptedValue

    stack.Values = append(stack.Values[:index], append([]common.StackValue{value}, stack.Values[index:]...)...)
    return ledger.UpdateStack(stack)
}

// RemoveStackValue removes a specific value from the stack at a given index.
func RemoveStackValue(stackID string, index int) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    if index < 0 || index >= len(stack.Values) {
        return errors.New("index out of bounds")
    }

    stack.Values = append(stack.Values[:index], stack.Values[index+1:]...)
    return ledger.UpdateStack(stack)
}
