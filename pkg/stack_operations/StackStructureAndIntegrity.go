package stack_operations

import (
    "errors"
    "fmt"
"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)
)

// ReverseStack reverses the order of elements in the stack.
func ReverseStack(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    for i, j := 0, len(stack.Values)-1; i < j; i, j = i+1, j-1 {
        stack.Values[i], stack.Values[j] = stack.Values[j], stack.Values[i]
    }
    return ledger.UpdateStack(stack)
}

// StackSearch searches for a specific value in the stack and returns its index.
func StackSearch(stackID string, target common.StackValue) (int, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return -1, fmt.Errorf("stack not found: %v", err)
    }

    decryptedTarget, err := encryption.DecryptData(target.Data)
    if err != nil {
        return -1, fmt.Errorf("failed to decrypt target value: %v", err)
    }

    for i, val := range stack.Values {
        decryptedData, err := encryption.DecryptData(val.Data)
        if err != nil {
            return -1, fmt.Errorf("failed to decrypt stack value: %v", err)
        }
        if decryptedData == decryptedTarget {
            return i, nil
        }
    }
    return -1, errors.New("value not found in stack")
}

// StackLength returns the number of elements in the stack.
func StackLength(stackID string) (int, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return 0, fmt.Errorf("stack not found: %v", err)
    }

    return len(stack.Values), nil
}

// ValidateStackIntegrity verifies the integrity of stack data based on a predefined hash or checksum.
func ValidateStackIntegrity(stackID string) (bool, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return false, fmt.Errorf("stack not found: %v", err)
    }

    for _, val := range stack.Values {
        decryptedData, err := encryption.DecryptData(val.Data)
        if err != nil || !encryption.VerifyChecksum(decryptedData, val.Checksum) {
            return false, fmt.Errorf("integrity check failed: %v", err)
        }
    }
    return true, nil
}

// StackTopFrequency returns the frequency of the top element in the stack.
func StackTopFrequency(stackID string) (int, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return 0, fmt.Errorf("stack not found: %v", err)
    }

    if len(stack.Values) == 0 {
        return 0, errors.New("stack is empty")
    }

    topValue := stack.Values[len(stack.Values)-1]
    count := 0

    for _, val := range stack.Values {
        if val.Value == topValue.Value {
            count++
        }
    }
    return count, nil
}

// CheckStackContains verifies if a specific value exists in the stack.
func CheckStackContains(stackID string, target common.StackValue) (bool, error) {
    index, err := StackSearch(stackID, target)
    if err != nil {
        if err.Error() == "value not found in stack" {
            return false, nil
        }
        return false, err
    }
    return index != -1, nil
}

// CheckStackEmptiness returns true if the stack is empty.
func CheckStackEmptiness(stackID string) (bool, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return false, fmt.Errorf("stack not found: %v", err)
    }

    return len(stack.Values) == 0, nil
}

// ValidateStackBounds checks that the stack size is within permissible bounds.
func ValidateStackBounds(stackID string, minSize, maxSize int) (bool, error) {
    stackLength, err := StackLength(stackID)
    if err != nil {
        return false, err
    }

    if stackLength < minSize || stackLength > maxSize {
        return false, nil
    }
    return true, nil
}

// LogStackOverflow logs an overflow error if the stack exceeds its maximum allowable size.
func LogStackOverflow(stackID string) error {
    overflowMessage := "Stack overflow: attempt to push beyond stack capacity."
    return logStackOperation(stackID, "Overflow", overflowMessage)
}

// LogStackUnderflow logs an underflow error if an operation is attempted on an empty stack.
func LogStackUnderflow(stackID string) error {
    underflowMessage := "Stack underflow: attempt to pop from an empty stack."
    return logStackOperation(stackID, "Underflow", underflowMessage)
}

// logStackOperation logs a specific operation or error message for the stack.
func logStackOperation(stackID, operation, message string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    encryptedMessage, err := encryption.EncryptData(message)
    if err != nil {
        return fmt.Errorf("failed to encrypt log message: %v", err)
    }

    logEntry := common.StackLog{
        Operation: operation,
        Details:   encryptedMessage,
        Timestamp: time.Now(),
    }

    stack.OperationLogs = append(stack.OperationLogs, logEntry)
    return ledger.UpdateStack(stack)
}
