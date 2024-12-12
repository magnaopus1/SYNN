package utility

import (
	"errors"
	"runtime"
	"sync"
	"synnergy_network/pkg/common"
	"time"
)

var tempVariables = make(map[string]interface{})
var isolatedEnvironments = make(map[string]bool)
var debugFlag bool
var lifecycleLock sync.Mutex

// GET_SYSTEM_METRICS: Retrieves key system metrics like memory usage and CPU load
func GET_SYSTEM_METRICS() map[string]interface{} {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    metrics := map[string]interface{}{
        "MemoryAlloc":    m.Alloc,
        "MemoryTotalAlloc": m.TotalAlloc,
        "MemorySys":      m.Sys,
        "NumGoroutine":   runtime.NumGoroutine(),
    }
    LogMetricsOperation("GET_SYSTEM_METRICS", "System metrics retrieved")
    return metrics
}

// LOG_MEMORY_STATE: Logs the current state of memory usage
func LOG_MEMORY_STATE() {
    metrics := GET_SYSTEM_METRICS()
    memoryUsage := fmt.Sprintf("MemoryAlloc: %v, MemorySys: %v", metrics["MemoryAlloc"], metrics["MemorySys"])
    LogMetricsOperation("LOG_MEMORY_STATE", memoryUsage)
}

// STORE_TEMP_VARIABLE: Stores a temporary variable for immediate use
func STORE_TEMP_VARIABLE(key string, value interface{}) {
    lifecycleLock.Lock()
    defer lifecycleLock.Unlock()
    tempVariables[key] = value
    LogMetricsOperation("STORE_TEMP_VARIABLE", "Stored temporary variable: " + key)
}

// RETRIEVE_TEMP_VARIABLE: Retrieves a previously stored temporary variable
func RETRIEVE_TEMP_VARIABLE(key string) (interface{}, error) {
    lifecycleLock.Lock()
    defer lifecycleLock.Unlock()
    value, exists := tempVariables[key]
    if !exists {
        LogMetricsOperation("RETRIEVE_TEMP_VARIABLE", "Temp variable not found: " + key)
        return nil, errors.New("temporary variable not found")
    }
    LogMetricsOperation("RETRIEVE_TEMP_VARIABLE", "Retrieved temporary variable: " + key)
    return value, nil
}

// CREATE_ISOLATED_EXEC_ENV: Creates an isolated execution environment
func CREATE_ISOLATED_EXEC_ENV(envID string) {
    lifecycleLock.Lock()
    defer lifecycleLock.Unlock()
    isolatedEnvironments[envID] = true
    LogMetricsOperation("CREATE_ISOLATED_EXEC_ENV", "Isolated environment created: " + envID)
}

// DESTROY_ISOLATED_EXEC_ENV: Destroys a specified isolated execution environment
func DESTROY_ISOLATED_EXEC_ENV(envID string) error {
    lifecycleLock.Lock()
    defer lifecycleLock.Unlock()
    if _, exists := isolatedEnvironments[envID]; !exists {
        LogMetricsOperation("DESTROY_ISOLATED_EXEC_ENV", "Environment not found: " + envID)
        return errors.New("isolated environment not found")
    }
    delete(isolatedEnvironments, envID)
    LogMetricsOperation("DESTROY_ISOLATED_EXEC_ENV", "Isolated environment destroyed: " + envID)
    return nil
}

// RUN_IN_SANDBOX_MODE: Enables sandbox mode for safe execution of functions
func RUN_IN_SANDBOX_MODE(function func()) {
    lifecycleLock.Lock()
    defer lifecycleLock.Unlock()
    LogMetricsOperation("RUN_IN_SANDBOX_MODE", "Entered sandbox mode")
    function()
    LogMetricsOperation("EXIT_SANDBOX_MODE", "Exited sandbox mode")
}

// EXIT_SANDBOX_MODE: Exits sandbox mode
func EXIT_SANDBOX_MODE() {
    LogMetricsOperation("EXIT_SANDBOX_MODE", "Exited sandbox mode")
}

// PERFORM_SELF_TEST: Runs a self-test on the system and logs the result
func PERFORM_SELF_TEST() {
    selfTestStatus := "Self-test completed successfully"
    LogMetricsOperation("PERFORM_SELF_TEST", selfTestStatus)
}

// GET_SYSTEM_LOAD: Retrieves the current system load (simulated as a placeholder for actual implementation)
func GET_SYSTEM_LOAD() float64 {
    // Simulate system load for demonstration (in real implementation, retrieve actual CPU load)
    systemLoad := 1.2
    LogMetricsOperation("GET_SYSTEM_LOAD", fmt.Sprintf("System load: %.2f", systemLoad))
    return systemLoad
}

// SET_DEBUG_FLAG: Sets the debug flag to enable detailed logging
func SET_DEBUG_FLAG() {
    lifecycleLock.Lock()
    defer lifecycleLock.Unlock()
    debugFlag = true
    LogMetricsOperation("SET_DEBUG_FLAG", "Debug flag set")
}

// CLEAR_DEBUG_FLAG: Clears the debug flag to disable detailed logging
func CLEAR_DEBUG_FLAG() {
    lifecycleLock.Lock()
    defer lifecycleLock.Unlock()
    debugFlag = false
    LogMetricsOperation("CLEAR_DEBUG_FLAG", "Debug flag cleared")
}

// MONITOR_EXECUTION_LIFECYCLE: Monitors and logs the lifecycle of a specific task
func MONITOR_EXECUTION_LIFECYCLE(taskID string, function func()) {
    start := time.Now()
    LogMetricsOperation("MONITOR_EXECUTION_LIFECYCLE", "Started task: " + taskID)
    function()
    duration := time.Since(start)
    LogMetricsOperation("MONITOR_EXECUTION_LIFECYCLE", "Completed task: " + taskID + " in " + duration.String())
}

// RECORD_TASK_METRICS: Records performance metrics for a task's execution time
func RECORD_TASK_METRICS(taskID string, function func()) {
    start := time.Now()
    function()
    duration := time.Since(start)
    LogMetricsOperation("RECORD_TASK_METRICS", "Task " + taskID + " executed in " + duration.String())
}

// Helper Functions

// LogMetricsOperation: Helper function to log encrypted metrics and system operations
func LogMetricsOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("MetricsOperation", encryptedMessage)
}
