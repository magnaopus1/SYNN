package utility

import (
	"errors"
	"sync"
	"synnergy_network/pkg/common"
	"time"
)

var lastExecutedOpcode string
var opcodeTracingEnabled bool
var opcodeFrequencyMap sync.Map
var executionDelay time.Duration
var memorySnapshots sync.Map
var opcodeDependencies sync.Map
var registerStates sync.Map

// GetLastExecutedOpcode: Retrieves the last executed opcode
func GetLastExecutedOpcode() string {
    LogDiagnostic("Opcode Execution", "Last executed opcode: " + lastExecutedOpcode)
    return lastExecutedOpcode
}

// GetOpcodeCondition: Retrieves conditions under which an opcode can execute
func GetOpcodeCondition(opcode string) (string, error) {
    condition, err := common.ledger.GetOpcodeCondition(opcode)
    if err != nil {
        LogDiagnostic("Opcode Condition", "Failed to retrieve condition for opcode " + opcode)
        return "", err
    }
    LogDiagnostic("Opcode Condition", "Condition for opcode " + opcode + ": " + condition)
    return condition, nil
}

// ViewMemorySnapshot: Retrieves a memory snapshot for a specific opcode
func ViewMemorySnapshot(opcode string) (interface{}, error) {
    snapshot, exists := memorySnapshots.Load(opcode)
    if !exists {
        return nil, errors.New("no memory snapshot available for opcode")
    }
    LogDiagnostic("Memory Snapshot", "Viewed memory snapshot for opcode " + opcode)
    return snapshot, nil
}

// LogRegisterState: Logs the state of a register at a specific point in execution
func LogRegisterState(opcode string, state interface{}) {
    registerStates.Store(opcode, state)
    LogDiagnostic("Register State", "Register state logged for opcode " + opcode)
}

// RevertToPreviousState: Reverts the system to a prior state
func RevertToPreviousState(opcode string) error {
    state, exists := registerStates.Load(opcode)
    if !exists {
        return errors.New("no previous state found for opcode")
    }
    err := common.ledger.RestoreState(state)
    if err != nil {
        LogDiagnostic("State Revert", "Failed to revert state for opcode " + opcode)
        return err
    }
    LogDiagnostic("State Revert", "State reverted for opcode " + opcode)
    return nil
}

// EnableOpcodeTracing: Enables tracing for all opcode executions
func EnableOpcodeTracing() {
    opcodeTracingEnabled = true
    LogDiagnostic("Opcode Tracing", "Opcode tracing enabled")
}

// DisableOpcodeTracing: Disables tracing for all opcode executions
func DisableOpcodeTracing() {
    opcodeTracingEnabled = false
    LogDiagnostic("Opcode Tracing", "Opcode tracing disabled")
}

// RecordOpcodeFrequency: Tracks the frequency of opcode executions
func RecordOpcodeFrequency(opcode string) {
    count, _ := opcodeFrequencyMap.LoadOrStore(opcode, 0)
    opcodeFrequencyMap.Store(opcode, count.(int)+1)
    LogDiagnostic("Opcode Frequency", "Opcode " + opcode + " executed, frequency updated")
}

// SetOpcodeExecutionDelay: Sets a delay before executing opcodes
func SetOpcodeExecutionDelay(delay time.Duration) {
    executionDelay = delay
    LogDiagnostic("Execution Delay", "Opcode execution delay set to " + delay.String())
}

// ClearExecutionDelay: Clears any set execution delay for opcodes
func ClearExecutionDelay() {
    executionDelay = 0
    LogDiagnostic("Execution Delay", "Opcode execution delay cleared")
}

// AnalyzeOpcodePerformance: Analyzes the performance of an opcode
func AnalyzeOpcodePerformance(opcode string) (string, error) {
    performanceData, err := common.ledger.GetOpcodePerformance(opcode)
    if err != nil {
        LogDiagnostic("Opcode Performance", "Failed to retrieve performance data for opcode " + opcode)
        return "", err
    }
    LogDiagnostic("Opcode Performance", "Performance data retrieved for opcode " + opcode)
    return performanceData, nil
}

// ViewOpcodeExecutionPath: Retrieves the execution path for an opcode
func ViewOpcodeExecutionPath(opcode string) ([]string, error) {
    path, err := common.ledger.GetExecutionPath(opcode)
    if err != nil {
        LogDiagnostic("Execution Path", "Failed to retrieve execution path for opcode " + opcode)
        return nil, err
    }
    LogDiagnostic("Execution Path", "Execution path retrieved for opcode " + opcode)
    return path, nil
}

// CheckOpcodeDependencies: Checks dependencies for a specific opcode
func CheckOpcodeDependencies(opcode string) ([]string, error) {
    deps, exists := opcodeDependencies.Load(opcode)
    if !exists {
        return nil, errors.New("no dependencies found for opcode")
    }
    LogDiagnostic("Opcode Dependencies", "Dependencies checked for opcode " + opcode)
    return deps.([]string), nil
}

// LogOpcodeParameters: Logs the parameters passed to an opcode securely
func LogOpcodeParameters(opcode string, params map[string]interface{}) error {
    paramString := formatParams(params)
    encryptedParams, err := encryption.Encrypt([]byte(paramString))
    if err != nil {
        return err
    }
    return common.ledger.LogParameters(opcode, encryptedParams)
}

// StoreOpcodeResult: Stores the result of an opcode execution securely
func StoreOpcodeResult(opcode string, result interface{}) error {
    encryptedResult, err := encryption.Encrypt([]byte(fmt.Sprintf("%v", result)))
    if err != nil {
        return err
    }
    return common.ledger.StoreOpcodeResult(opcode, encryptedResult)
}

// Helper Functions

// formatParams: Helper to format parameters as a string for logging
func formatParams(params map[string]interface{}) string {
    var formatted string
    for key, value := range params {
        formatted += key + ": " + fmt.Sprintf("%v", value) + "; "
    }
    return formatted
}

// LogDiagnostic: Helper function for encrypted diagnostic logging
func LogDiagnostic(context, message string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogDiagnostic(context, encryptedMessage)
}
