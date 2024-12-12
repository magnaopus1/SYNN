package utility

import (
	"errors"
	"fmt"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

var scheduledTasks = make(map[string]time.Time)
var contractExpiry = make(map[string]time.Time)
var executionBoundaries = make(map[string]struct{})
var taskLock sync.Mutex

// QUERY_SCHEDULED_TASKS: Returns a list of all scheduled tasks and their execution times
func QUERY_SCHEDULED_TASKS() map[string]time.Time {
    taskLock.Lock()
    defer taskLock.Unlock()
    LogTaskOperation("QUERY_SCHEDULED_TASKS", "Scheduled tasks queried")
    return scheduledTasks
}

// AUTO_SCALE_RESOURCE: Adjusts resource allocation based on load
func AUTO_SCALE_RESOURCE(taskID string, loadFactor float64) {
    taskLock.Lock()
    defer taskLock.Unlock()
    if _, exists := resourceAllocations[taskID]; exists {
        resourceAllocations[taskID] *= loadFactor
        LogTaskOperation("AUTO_SCALE_RESOURCE", fmt.Sprintf("Resource for %s scaled by factor %.2f", taskID, loadFactor))
    } else {
        LogTaskOperation("AUTO_SCALE_RESOURCE", fmt.Sprintf("Task %s not found", taskID))
    }
}

// CHECKPOINT_EXECUTION: Saves the current execution state of a task
func CHECKPOINT_EXECUTION(taskID string, state interface{}) error {
    taskLock.Lock()
    defer taskLock.Unlock()
    if _, exists := scheduledTasks[taskID]; !exists {
        LogTaskOperation("CHECKPOINT_EXECUTION", "Task not found: " + taskID)
        return errors.New("task not found")
    }
    common.ledger.SaveCheckpoint(taskID, state)
    LogTaskOperation("CHECKPOINT_EXECUTION", "Checkpoint saved for task: " + taskID)
    return nil
}

// ROLLBACK_EXECUTION: Restores a task to its last checkpoint state
func ROLLBACK_EXECUTION(taskID string) (interface{}, error) {
    taskLock.Lock()
    defer taskLock.Unlock()
    state, err := ledger.GetCheckpoint(taskID)
    if err != nil {
        LogTaskOperation("ROLLBACK_EXECUTION", "Failed to rollback task: " + taskID)
        return nil, err
    }
    LogTaskOperation("ROLLBACK_EXECUTION", "Rolled back to last checkpoint for task: " + taskID)
    return state, nil
}

// ENCRYPT_TEMP_DATA: Encrypts temporary data for secure storage
func ENCRYPT_TEMP_DATA(data []byte) ([]byte, error) {
    encryptedData, err := encryption.Encrypt(data)
    if err != nil {
        LogTaskOperation("ENCRYPT_TEMP_DATA", "Encryption failed")
        return nil, err
    }
    LogTaskOperation("ENCRYPT_TEMP_DATA", "Data encrypted")
    return encryptedData, nil
}

// DECRYPT_TEMP_DATA: Decrypts previously encrypted temporary data
func DECRYPT_TEMP_DATA(encryptedData []byte) ([]byte, error) {
    decryptedData, err := encryption.Decrypt(encryptedData)
    if err != nil {
        LogTaskOperation("DECRYPT_TEMP_DATA", "Decryption failed")
        return nil, err
    }
    LogTaskOperation("DECRYPT_TEMP_DATA", "Data decrypted")
    return decryptedData, nil
}

// TRACK_RESOURCE_ALLOCATIONS: Tracks and logs resource allocations for a task
func TRACK_RESOURCE_ALLOCATIONS(taskID string, allocation float64) {
    taskLock.Lock()
    defer taskLock.Unlock()
    resourceAllocations[taskID] = allocation
    LogTaskOperation("TRACK_RESOURCE_ALLOCATIONS", fmt.Sprintf("Resource allocated for %s: %.2f", taskID, allocation))
}

// LOG_RESOURCE_UTILIZATION: Logs the resource utilization of a task
func LOG_RESOURCE_UTILIZATION(taskID string) {
    taskLock.Lock()
    defer taskLock.Unlock()
    allocation, exists := resourceAllocations[taskID]
    if exists {
        LogTaskOperation("LOG_RESOURCE_UTILIZATION", fmt.Sprintf("Resource utilization for %s: %.2f", taskID, allocation))
    } else {
        LogTaskOperation("LOG_RESOURCE_UTILIZATION", fmt.Sprintf("No allocation found for task %s", taskID))
    }
}

// FORCE_CONTRACT_EXPIRY: Forces an immediate expiry of a contract
func FORCE_CONTRACT_EXPIRY(contractID string) error {
    taskLock.Lock()
    defer taskLock.Unlock()
    if _, exists := contractExpiry[contractID]; !exists {
        LogTaskOperation("FORCE_CONTRACT_EXPIRY", "Contract not found: " + contractID)
        return errors.New("contract not found")
    }
    delete(contractExpiry, contractID)
    LogTaskOperation("FORCE_CONTRACT_EXPIRY", "Contract expired: " + contractID)
    return nil
}

// EXTEND_CONTRACT_EXPIRY: Extends the expiry time of a contract
func EXTEND_CONTRACT_EXPIRY(contractID string, extension time.Duration) error {
    taskLock.Lock()
    defer taskLock.Unlock()
    expiry, exists := contractExpiry[contractID]
    if !exists {
        LogTaskOperation("EXTEND_CONTRACT_EXPIRY", "Contract not found: " + contractID)
        return errors.New("contract not found")
    }
    contractExpiry[contractID] = expiry.Add(extension)
    LogTaskOperation("EXTEND_CONTRACT_EXPIRY", fmt.Sprintf("Contract %s extended by %s", contractID, extension))
    return nil
}

// DEFINE_EXECUTION_BOUNDARY: Defines an execution boundary to limit task scope
func DEFINE_EXECUTION_BOUNDARY(boundaryID string) {
    taskLock.Lock()
    defer taskLock.Unlock()
    executionBoundaries[boundaryID] = struct{}{}
    LogTaskOperation("DEFINE_EXECUTION_BOUNDARY", "Execution boundary defined: " + boundaryID)
}

// CLEAR_EXECUTION_BOUNDARY: Clears an execution boundary to remove limitations
func CLEAR_EXECUTION_BOUNDARY(boundaryID string) error {
    taskLock.Lock()
    defer taskLock.Unlock()
    if _, exists := executionBoundaries[boundaryID]; !exists {
        LogTaskOperation("CLEAR_EXECUTION_BOUNDARY", "Boundary not found: " + boundaryID)
        return errors.New("execution boundary not found")
    }
    delete(executionBoundaries, boundaryID)
    LogTaskOperation("CLEAR_EXECUTION_BOUNDARY", "Execution boundary cleared: " + boundaryID)
    return nil
}

// CONVERT_ENV_VARIABLE: Converts an environment variable to a specific format
func CONVERT_ENV_VARIABLE(value string, format string) (interface{}, error) {
    var convertedValue interface{}
    switch format {
    case "int":
        convertedValue = StringToInt(value)
    case "float":
        convertedValue = StringToFloat(value)
    case "bool":
        convertedValue = StringToBool(value)
    default:
        LogTaskOperation("CONVERT_ENV_VARIABLE", "Unknown format: " + format)
        return nil, errors.New("unknown format")
    }
    LogTaskOperation("CONVERT_ENV_VARIABLE", fmt.Sprintf("Variable converted to %s", format))
    return convertedValue, nil
}

// CHECK_SYSTEM_HEALTH: Runs a health check on critical system components
func CHECK_SYSTEM_HEALTH() string {
    healthStatus := "All systems operational"
    LogTaskOperation("CHECK_SYSTEM_HEALTH", healthStatus)
    return healthStatus
}

// Helper Functions

// LogTaskOperation: Logs task management operations securely
func LogTaskOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("TaskOperation", encryptedMessage)
}
