package utility

import (
	"sort"
	"synnergy_network/pkg/ledger"
)

// PushMultiple pushes multiple elements onto the stack in a single operation.
func PushMultiple(stack *[]interface{}, elements ...interface{}) {
	*stack = append(*stack, elements...)
	ledger.RecordStackOperation("PushMultiple", elements)
}

// PopMultiple removes a specified number of elements from the stack.
func PopMultiple(stack *[]interface{}, count int) ([]interface{}, error) {
	if count > len(*stack) {
		return nil, errors.New("count exceeds stack size")
	}
	popped := (*stack)[len(*stack)-count:]
	*stack = (*stack)[:len(*stack)-count]
	ledger.RecordStackOperation("PopMultiple", popped)
	return popped, nil
}

// SplitStack divides the stack into two stacks at a specified index.
func SplitStack(stack *[]interface{}, index int) ([]interface{}, []interface{}, error) {
	if index < 0 || index > len(*stack) {
		return nil, nil, errors.New("index out of bounds")
	}
	stack1 := (*stack)[:index]
	stack2 := (*stack)[index:]
	ledger.RecordStackOperation("SplitStack", index)
	return stack1, stack2, nil
}

// MergeStacks merges two stacks into one, appending stack2 to stack1.
func MergeStacks(stack1, stack2 *[]interface{}) []interface{} {
	merged := append(*stack1, *stack2...)
	ledger.RecordStackOperation("MergeStacks", merged)
	return merged
}

// CloneStack creates a clone of the stack.
func CloneStack(stack *[]interface{}) []interface{} {
	cloned := append([]interface{}{}, *stack...)
	ledger.RecordStackOperation("CloneStack", cloned)
	return cloned
}

// SortStack sorts the stack in ascending order if values are comparable.
func SortStack(stack *[]int) {
	sort.Ints(*stack)
	ledger.RecordStackOperation("SortStack", *stack)
}

// FilterStack filters stack elements based on a filter function.
func FilterStack(stack *[]interface{}, filter func(interface{}) bool) []interface{} {
	var filtered []interface{}
	for _, elem := range *stack {
		if filter(elem) {
			filtered = append(filtered, elem)
		}
	}
	ledger.RecordStackOperation("FilterStack", filtered)
	return filtered
}

// MapStackValues applies a transformation function to each element in the stack.
func MapStackValues(stack *[]interface{}, transform func(interface{}) interface{}) []interface{} {
	mapped := make([]interface{}, len(*stack))
	for i, elem := range *stack {
		mapped[i] = transform(elem)
	}
	ledger.RecordStackOperation("MapStackValues", mapped)
	return mapped
}

// AccumulateStack reduces the stack to a single accumulated value based on a function.
func AccumulateStack(stack *[]interface{}, initial interface{}, accumulator func(interface{}, interface{}) interface{}) interface{} {
	result := initial
	for _, elem := range *stack {
		result = accumulator(result, elem)
	}
	ledger.RecordStackOperation("AccumulateStack", result)
	return result
}

// StackTopFrequency returns the frequency of the top element in the stack.
func StackTopFrequency(stack *[]interface{}) int {
	if len(*stack) == 0 {
		return 0
	}
	top := (*stack)[len(*stack)-1]
	count := 0
	for _, elem := range *stack {
		if elem == top {
			count++
		}
	}
	ledger.RecordStackOperation("StackTopFrequency", count)
	return count
}

// IsStackSorted checks if the stack is sorted in ascending order.
func IsStackSorted(stack *[]int) bool {
	for i := 1; i < len(*stack); i++ {
		if (*stack)[i-1] > (*stack)[i] {
			return false
		}
	}
	ledger.RecordStackOperation("IsStackSorted", *stack)
	return true
}

// StackDifference returns elements in stack1 not found in stack2.
func StackDifference(stack1, stack2 *[]interface{}) []interface{} {
	diff := []interface{}{}
	for _, elem1 := range *stack1 {
		found := false
		for _, elem2 := range *stack2 {
			if elem1 == elem2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, elem1)
		}
	}
	ledger.RecordStackOperation("StackDifference", diff)
	return diff
}

// StackUnion returns a union of elements in stack1 and stack2, without duplicates.
func StackUnion(stack1, stack2 *[]interface{}) []interface{} {
	union := append([]interface{}{}, *stack1...)
	for _, elem2 := range *stack2 {
		found := false
		for _, elem1 := range union {
			if elem1 == elem2 {
				found = true
				break
			}
		}
		if !found {
			union = append(union, elem2)
		}
	}
	ledger.RecordStackOperation("StackUnion", union)
	return union
}

// StackIntersection returns common elements in stack1 and stack2.
func StackIntersection(stack1, stack2 *[]interface{}) []interface{} {
	intersection := []interface{}{}
	for _, elem1 := range *stack1 {
		for _, elem2 := range *stack2 {
			if elem1 == elem2 {
				intersection = append(intersection, elem1)
				break
			}
		}
	}
	ledger.RecordStackOperation("StackIntersection", intersection)
	return intersection
}

// StackExclusiveElements returns elements exclusive to each stack (elements in one but not both).
func StackExclusiveElements(stack1, stack2 *[]interface{}) []interface{} {
	diff1 := StackDifference(stack1, stack2)
	diff2 := StackDifference(stack2, stack1)
	exclusive := append(diff1, diff2...)
	ledger.RecordStackOperation("StackExclusiveElements", exclusive)
	return exclusive
}
