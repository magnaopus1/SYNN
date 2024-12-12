package utility

import (
	"errors"
	"sync"
	"synnergy_network/pkg/common"
)

type Stack struct {
    items []interface{}
    lock  sync.Mutex
}

// PushStack: Pushes a value onto the stack
func (s *Stack) PushStack(value interface{}) error {
    s.lock.Lock()
    defer s.lock.Unlock()
    s.items = append(s.items, value)
    LogStackOperation("PushStack", value)
    return nil
}

// PopStack: Pops a value from the stack and returns it
func (s *Stack) PopStack() (interface{}, error) {
    s.lock.Lock()
    defer s.lock.Unlock()
    if len(s.items) == 0 {
        return nil, errors.New("stack is empty")
    }
    value := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    LogStackOperation("PopStack", value)
    return value, nil
}

// PeekStack: Peeks at the top value of the stack without removing it
func (s *Stack) PeekStack() (interface{}, error) {
    s.lock.Lock()
    defer s.lock.Unlock()
    if len(s.items) == 0 {
        return nil, errors.New("stack is empty")
    }
    value := s.items[len(s.items)-1]
    LogStackOperation("PeekStack", value)
    return value, nil
}

// ClearStack: Clears all values from the stack
func (s *Stack) ClearStack() {
    s.lock.Lock()
    defer s.lock.Unlock()
    s.items = []interface{}{}
    LogStackOperation("ClearStack", "stack cleared")
}

// DuplicateStack: Duplicates the top value of the stack
func (s *Stack) DuplicateStack() error {
    s.lock.Lock()
    defer s.lock.Unlock()
    if len(s.items) == 0 {
        return errors.New("stack is empty")
    }
    s.items = append(s.items, s.items[len(s.items)-1])
    LogStackOperation("DuplicateStack", "top value duplicated")
    return nil
}

// SwapStackTop: Swaps the top two values of the stack
func (s *Stack) SwapStackTop() error {
    s.lock.Lock()
    defer s.lock.Unlock()
    if len(s.items) < 2 {
        return errors.New("not enough elements to swap")
    }
    s.items[len(s.items)-1], s.items[len(s.items)-2] = s.items[len(s.items)-2], s.items[len(s.items)-1]
    LogStackOperation("SwapStackTop", "top two values swapped")
    return nil
}

// RotateStack: Rotates the stack by moving the top value to the bottom
func (s *Stack) RotateStack() {
    s.lock.Lock()
    defer s.lock.Unlock()
    if len(s.items) > 1 {
        top := s.items[len(s.items)-1]
        s.items = append([]interface{}{top}, s.items[:len(s.items)-1]...)
        LogStackOperation("RotateStack", "stack rotated")
    }
}

// CopyToStack: Copies an element at a specific index to the top of the stack
func (s *Stack) CopyToStack(index int) error {
    s.lock.Lock()
    defer s.lock.Unlock()
    if index < 0 || index >= len(s.items) {
        return errors.New("index out of range")
    }
    s.items = append(s.items, s.items[index])
    LogStackOperation("CopyToStack", "element copied to top")
    return nil
}

// ReplaceStackValue: Replaces a value at a specific index in the stack
func (s *Stack) ReplaceStackValue(index int, newValue interface{}) error {
    s.lock.Lock()
    defer s.lock.Unlock()
    if index < 0 || index >= len(s.items) {
        return errors.New("index out of range")
    }
    s.items[index] = newValue
    LogStackOperation("ReplaceStackValue", "value replaced at index "+string(index))
    return nil
}

// ReverseStack: Reverses the order of the stack
func (s *Stack) ReverseStack() {
    s.lock.Lock()
    defer s.lock.Unlock()
    for i, j := 0, len(s.items)-1; i < j; i, j = i+1, j-1 {
        s.items[i], s.items[j] = s.items[j], s.items[i]
    }
    LogStackOperation("ReverseStack", "stack reversed")
}

// InsertStackValue: Inserts a value at a specified index in the stack
func (s *Stack) InsertStackValue(index int, value interface{}) error {
    s.lock.Lock()
    defer s.lock.Unlock()
    if index < 0 || index > len(s.items) {
        return errors.New("index out of range")
    }
    s.items = append(s.items[:index], append([]interface{}{value}, s.items[index:]...)...)
    LogStackOperation("InsertStackValue", "value inserted at index "+string(index))
    return nil
}

// RemoveStackValue: Removes a value at a specified index in the stack
func (s *Stack) RemoveStackValue(index int) error {
    s.lock.Lock()
    defer s.lock.Unlock()
    if index < 0 || index >= len(s.items) {
        return errors.New("index out of range")
    }
    s.items = append(s.items[:index], s.items[index+1:]...)
    LogStackOperation("RemoveStackValue", "value removed at index "+string(index))
    return nil
}

// StackSearch: Searches for a value in the stack and returns its index, or -1 if not found
func (s *Stack) StackSearch(value interface{}) int {
    s.lock.Lock()
    defer s.lock.Unlock()
    for i := len(s.items) - 1; i >= 0; i-- {
        if s.items[i] == value {
            LogStackOperation("StackSearch", "value found at index "+string(i))
            return i
        }
    }
    LogStackOperation("StackSearch", "value not found")
    return -1
}

// StackLength: Returns the current length of the stack
func (s *Stack) StackLength() int {
    s.lock.Lock()
    defer s.lock.Unlock()
    length := len(s.items)
    LogStackOperation("StackLength", "current stack length: "+string(length))
    return length
}

// ValidateStackIntegrity: Validates the integrity of the stack by ensuring non-nil values
func (s *Stack) ValidateStackIntegrity() bool {
    s.lock.Lock()
    defer s.lock.Unlock()
    for _, item := range s.items {
        if item == nil {
            LogStackOperation("ValidateStackIntegrity", false)
            return false
        }
    }
    LogStackOperation("ValidateStackIntegrity", true)
    return true
}

// Helper Functions

// LogStackOperation: Helper function to log encrypted stack operations
func LogStackOperation(operation string, details interface{}) error {
    message := "Operation: " + operation + " - Details: " + fmt.Sprintf("%v", details)
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("StackOperation", encryptedMessage)
}
