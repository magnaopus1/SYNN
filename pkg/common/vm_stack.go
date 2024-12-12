package common


import (
	"errors"
	"fmt"
)

// VMStack represents the virtual machine’s stack, used for storing temporary data and managing function calls.
type VMStack struct {
	Stack []int64 // The stack memory, using int64 for flexibility in handling large values
	Size  int     // Maximum stack size to prevent overflow
	Top   int     // Points to the current top of the stack
}

// NewVMStack initializes a new VM stack with a given maximum size.
func NewVMStack(size int) *VMStack {
	return &VMStack{
		Stack: make([]int64, size),
		Size:  size,
		Top:   -1, // Top is initialized to -1 to indicate an empty stack
	}
}

// Push adds a value to the top of the stack.
func (s *VMStack) Push(value int64) error {
	if s.Top >= s.Size-1 {
		return errors.New("stack overflow: unable to push, stack is full")
	}
	s.Top++
	s.Stack[s.Top] = value
	return nil
}

// Pop removes and returns the top value from the stack.
func (s *VMStack) Pop() (int64, error) {
	if s.Top == -1 {
		return 0, errors.New("stack underflow: unable to pop, stack is empty")
	}
	value := s.Stack[s.Top]
	s.Stack[s.Top] = 0 // Clear the value for security purposes
	s.Top--
	return value, nil
}

// Peek returns the top value from the stack without removing it.
func (s *VMStack) Peek() (int64, error) {
	if s.Top == -1 {
		return 0, errors.New("stack is empty: unable to peek")
	}
	return s.Stack[s.Top], nil
}

// Dup duplicates the top value on the stack.
func (s *VMStack) Dup() error {
	value, err := s.Peek()
	if err != nil {
		return err
	}
	return s.Push(value)
}

// Swap swaps the top two values on the stack.
func (s *VMStack) Swap() error {
	if s.Top < 1 {
		return errors.New("stack underflow: unable to swap, not enough elements")
	}
	s.Stack[s.Top], s.Stack[s.Top-1] = s.Stack[s.Top-1], s.Stack[s.Top]
	return nil
}

// Clear resets the stack, removing all elements.
func (s *VMStack) Clear() {
	s.Stack = make([]int64, s.Size)
	s.Top = -1
}

// IsEmpty checks if the stack is empty.
func (s *VMStack) IsEmpty() bool {
	return s.Top == -1
}

// Dump prints the current stack contents for debugging.
func (s *VMStack) Dump() {
	fmt.Println("Current Stack State:")
	for i := s.Top; i >= 0; i-- {
		fmt.Printf("[%d]: %d\n", i, s.Stack[i])
	}
	fmt.Println("End of Stack")
}
