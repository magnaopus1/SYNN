package stack_operations

import (
    "errors"
    "fmt"
    "math"
    "sort"
	"synnergy_network/pkg/ledger"
)

// StackModuloOperation applies modulo operation on the top two stack values.
func StackModuloOperation(stackID string) (float64, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return 0, fmt.Errorf("stack not found: %v", err)
    }
    
    if len(stack.Values) < 2 {
        return 0, errors.New("not enough values on stack for modulo operation")
    }

    topValue, secondValue := stack.Values[len(stack.Values)-1], stack.Values[len(stack.Values)-2]
    modResult := math.Mod(secondValue.Value, topValue.Value)

    stack.Values[len(stack.Values)-2].Value = modResult
    stack.Values = stack.Values[:len(stack.Values)-1]
    return modResult, ledger.UpdateStack(stack)
}

// StackAbsoluteValues converts all values on the stack to their absolute values.
func StackAbsoluteValues(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    for i := range stack.Values {
        stack.Values[i].Value = math.Abs(stack.Values[i].Value)
    }
    return ledger.UpdateStack(stack)
}

// StackPowerValues raises each stack value to the power of a given exponent.
func StackPowerValues(stackID string, exponent float64) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    for i := range stack.Values {
        stack.Values[i].Value = math.Pow(stack.Values[i].Value, exponent)
    }
    return ledger.UpdateStack(stack)
}

// StackExponentialValues applies the exponential function to each value on the stack.
func StackExponentialValues(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    for i := range stack.Values {
        stack.Values[i].Value = math.Exp(stack.Values[i].Value)
    }
    return ledger.UpdateStack(stack)
}

// StackFloorValues applies the floor function to each value on the stack.
func StackFloorValues(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    for i := range stack.Values {
        stack.Values[i].Value = math.Floor(stack.Values[i].Value)
    }
    return ledger.UpdateStack(stack)
}

// StackCeilingValues applies the ceiling function to each value on the stack.
func StackCeilingValues(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    for i := range stack.Values {
        stack.Values[i].Value = math.Ceil(stack.Values[i].Value)
    }
    return ledger.UpdateStack(stack)
}

// FindStackMedian finds the median value of the stack values.
func FindStackMedian(stackID string) (float64, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return 0, fmt.Errorf("stack not found: %v", err)
    }

    if len(stack.Values) == 0 {
        return 0, errors.New("stack is empty")
    }

    values := make([]float64, len(stack.Values))
    for i, v := range stack.Values {
        values[i] = v.Value
    }

    sort.Float64s(values)
    mid := len(values) / 2
    median := values[mid]
    if len(values)%2 == 0 {
        median = (values[mid-1] + values[mid]) / 2
    }

    return median, nil
}

// StackRounding rounds all values on the stack to the nearest integer.
func StackRounding(stackID string) error {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return fmt.Errorf("stack not found: %v", err)
    }

    for i := range stack.Values {
        stack.Values[i].Value = math.Round(stack.Values[i].Value)
    }
    return ledger.UpdateStack(stack)
}

// StackStandardDeviation calculates the standard deviation of the stack values.
func StackStandardDeviation(stackID string) (float64, error) {
    variance, err := StackVariance(stackID)
    if err != nil {
        return 0, err
    }
    return math.Sqrt(variance), nil
}

// StackAverage calculates the average of the values on the stack.
func StackAverage(stackID string) (float64, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return 0, fmt.Errorf("stack not found: %v", err)
    }

    if len(stack.Values) == 0 {
        return 0, errors.New("stack is empty")
    }

    sum := 0.0
    for _, value := range stack.Values {
        sum += value.Value
    }
    return sum / float64(len(stack.Values)), nil
}

// StackVariance calculates the variance of the values on the stack.
func StackVariance(stackID string) (float64, error) {
    average, err := StackAverage(stackID)
    if err != nil {
        return 0, err
    }

    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return 0, fmt.Errorf("stack not found: %v", err)
    }

    variance := 0.0
    for _, value := range stack.Values {
        diff := value.Value - average
        variance += diff * diff
    }
    return variance / float64(len(stack.Values)), nil
}

// StackSum calculates the sum of all values on the stack.
func StackSum(stackID string) (float64, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return 0, fmt.Errorf("stack not found: %v", err)
    }

    sum := 0.0
    for _, value := range stack.Values {
        sum += value.Value
    }
    return sum, nil
}

// StackProduct calculates the product of all values on the stack.
func StackProduct(stackID string) (float64, error) {
    stack, err := ledger.GetStack(stackID)
    if err != nil {
        return 0, fmt.Errorf("stack not found: %v", err)
    }

    product := 1.0
    for _, value := range stack.Values {
        product *= value.Value
    }
    return product, nil
}
