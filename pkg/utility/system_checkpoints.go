package utility

import (
	"errors"
	"fmt"
	"sync"
	"synnergy_network/pkg/common"
	"time"
)

var systemCheckpoints = make(map[string]interface{})
var systemCache = make(map[string]interface{})
var resourceQuota float64
var resourceLimit float64
var checkpointLock sync.Mutex
var cacheLock sync.Mutex
var flagLock sync.Mutex

// DEFINE_SYSTEM_CHECKPOINT: Defines a checkpoint for the current system state
func DEFINE_SYSTEM_CHECKPOINT(name string, state interface{}) {
    checkpointLock.Lock()
    defer checkpointLock.Unlock()
    systemCheckpoints[name] = state
    LogSystemOperation("DEFINE_SYSTEM_CHECKPOINT", "Checkpoint defined: " + name)
}

// ROLLBACK_TO_CHECKPOINT: Rolls back the system to a previously defined checkpoint
func ROLLBACK_TO_CHECKPOINT(name string) (interface{}, error) {
    checkpointLock.Lock()
    defer checkpointLock.Unlock()
    state, exists := systemCheckpoints[name]
    if !exists {
        LogSystemOperation("ROLLBACK_TO_CHECKPOINT", "Checkpoint not found: " + name)
        return nil, errors.New("checkpoint not found")
    }
    LogSystemOperation("ROLLBACK_TO_CHECKPOINT", "Rolled back to checkpoint: " + name)
    return state, nil
}

// EXECUTE_WITH_DELAY: Executes a function after a specified delay
func EXECUTE_WITH_DELAY(delay time.Duration, function func()) {
    time.AfterFunc(delay, func() {
        start := time.Now()
        function()
        duration := time.Since(start)
        LogSystemOperation("EXECUTE_WITH_DELAY", "Execution completed in " + duration.String())
    })
}

// LOG_EXECUTION_TIME: Logs the execution time of a function
func LOG_EXECUTION_TIME(function func()) {
    start := time.Now()
    function()
    duration := time.Since(start)
    LogSystemOperation("LOG_EXECUTION_TIME", "Function executed in " + duration.String())
}

// GENERATE_SYSTEM_REPORT: Generates a report of the current system state
func GENERATE_SYSTEM_REPORT() map[string]interface{} {
    checkpointLock.Lock()
    cacheLock.Lock()
    defer checkpointLock.Unlock()
    defer cacheLock.Unlock()
    report := map[string]interface{}{
        "checkpoints": systemCheckpoints,
        "cache":       systemCache,
        "flags":       systemFlags,
        "resourceQuota": resourceQuota,
        "resourceLimit": resourceLimit,
    }
    LogSystemOperation("GENERATE_SYSTEM_REPORT", "System report generated")
    return report
}

// CLEAR_SYSTEM_CACHE: Clears all cached data
func CLEAR_SYSTEM_CACHE() {
    cacheLock.Lock()
    defer cacheLock.Unlock()
    systemCache = make(map[string]interface{})
    LogSystemOperation("CLEAR_SYSTEM_CACHE", "System cache cleared")
}

// CACHE_EXECUTION_RESULT: Caches the result of an execution
func CACHE_EXECUTION_RESULT(key string, result interface{}) {
    cacheLock.Lock()
    defer cacheLock.Unlock()
    systemCache[key] = result
    LogSystemOperation("CACHE_EXECUTION_RESULT", "Result cached with key: " + key)
}

// RETRIEVE_CACHED_RESULT: Retrieves a cached result by key
func RETRIEVE_CACHED_RESULT(key string) (interface{}, error) {
    cacheLock.Lock()
    defer cacheLock.Unlock()
    result, exists := systemCache[key]
    if !exists {
        LogSystemOperation("RETRIEVE_CACHED_RESULT", "Cache miss for key: " + key)
        return nil, errors.New("cached result not found")
    }
    LogSystemOperation("RETRIEVE_CACHED_RESULT", "Cache hit for key: " + key)
    return result, nil
}


// SYSTEM_SHUTDOWN_TIMER: Schedules a system shutdown after a specified delay
func SYSTEM_SHUTDOWN_TIMER(delay time.Duration, shutdownFunc func()) {
    EXECUTE_WITH_DELAY(delay, shutdownFunc)
    LogSystemOperation("SYSTEM_SHUTDOWN_TIMER", "Shutdown scheduled after delay: " + delay.String())
}

// SYSTEM_STARTUP_TIMER: Schedules a system startup after a specified delay
func SYSTEM_STARTUP_TIMER(delay time.Duration, startupFunc func()) {
    EXECUTE_WITH_DELAY(delay, startupFunc)
    LogSystemOperation("SYSTEM_STARTUP_TIMER", "Startup scheduled after delay: " + delay.String())
}

// CHECK_RESOURCE_LIMITS: Checks if a specified resource usage is within the limit
func CHECK_RESOURCE_LIMITS(usage float64) bool {
    result := usage <= resourceLimit
    LogSystemOperation("CHECK_RESOURCE_LIMITS", "Resource usage " + floatToString(usage) + " within limits: " + boolToString(result))
    return result
}

// SET_RESOURCE_QUOTA: Sets a resource quota for the system
func SET_RESOURCE_QUOTA(quota float64) {
    resourceQuota = quota
    LogSystemOperation("SET_RESOURCE_QUOTA", "Resource quota set to " + floatToString(quota))
}

// Helper Functions

// LogSystemOperation: Helper function to log encrypted system operations
func LogSystemOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("SystemOperation", encryptedMessage)
}


// floatToString: Converts a float64 value to string for logging
func floatToString(value float64) string {
    return fmt.Sprintf("%.2f", value)
}
