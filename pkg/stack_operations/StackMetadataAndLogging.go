package stack_operations

import (
    "errors"
    "fmt"
    "time"
"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// StackLastAccessedIndex returns the index of the last accessed value on the stack.
func StackLastAccessedIndex(stackID string) (int, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return -1, fmt.Errorf("stack not found: %v", err)
    }
    
    return stack.LastAccessedIndex, nil
}

// LogStackOperations logs a specific operation performed on the stack.
func LogStackOperations(stackID string, operation string, details string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    encryptedDetails, err := encryption.EncryptData(details)
    if err != nil {
        return fmt.Errorf("failed to encrypt operation details: %v", err)
    }

    logEntry := common.StackLog{
        Operation:   operation,
        Details:     encryptedDetails,
        Timestamp:   time.Now(),
    }

    stack.OperationLogs = append(stack.OperationLogs, logEntry)
    return ledger.UpdateStack(stack)
}

// LogStackOverflow records an overflow error when the stack exceeds its maximum capacity.
func LogStackOverflow(stackID string) error {
    overflowMessage := "Stack overflow: attempt to push beyond stack capacity."
    return LogStackOperations(stackID, "Overflow", overflowMessage)
}

// LogStackUnderflow records an underflow error when a pop operation is attempted on an empty stack.
func LogStackUnderflow(stackID string) error {
    underflowMessage := "Stack underflow: attempt to pop from an empty stack."
    return LogStackOperations(stackID, "Underflow", underflowMessage)
}

// LogStackMetadata logs metadata for stack operations, like size and last access time.
func LogStackMetadata(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    metadataMessage := fmt.Sprintf("Stack size: %d, Last accessed index: %d", len(stack.Values), stack.LastAccessedIndex)
    encryptedMetadata, err := encryption.EncryptData(metadataMessage)
    if err != nil {
        return fmt.Errorf("failed to encrypt metadata: %v", err)
    }

    logEntry := common.StackLog{
        Operation: "Metadata",
        Details:   encryptedMetadata,
        Timestamp: time.Now(),
    }

    stack.OperationLogs = append(stack.OperationLogs, logEntry)
    return ledger.UpdateStack(stack)
}
