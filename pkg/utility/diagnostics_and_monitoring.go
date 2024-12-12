package utility

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"synnergy_network/pkg/common"
	"time"
)

var taskStatuses sync.Map
var resourceSharingEnabled = false

// QUERY_EXECUTION_STATUS: Retrieve the current status of a running task
func QUERY_EXECUTION_STATUS(taskID string) (string, error) {
    status, exists := taskStatuses.Load(taskID)
    if !exists {
        return "", errors.New("task ID not found")
    }
    return status.(string), nil
}

// SET_TIMEOUT_FOR_TASK: Set timeout for a specific task
func SET_TIMEOUT_FOR_TASK(taskID string, duration time.Duration) error {
    go func() {
        time.Sleep(duration)
        taskStatuses.Store(taskID, "Timed out")
        // Log timeout event
        LOG_SYSTEM_EVENT("Timeout", "Task " + taskID + " timed out after " + duration.String())
    }()
    return nil
}

// RESET_CONTRACT_ENVIRONMENT: Reset the environment for a specific contract
func RESET_CONTRACT_ENVIRONMENT(contractID string) error {
    // Reset environment variables, cleanup operations for contract
    err := common.ResetEnvironment(contractID)
    if err != nil {
        return err
    }
    LOG_SYSTEM_EVENT("Environment Reset", "Environment for contract " + contractID + " has been reset")
    return nil
}

// LOG_SYSTEM_EVENT: Logs system events in the ledger for accountability and audit
func LOG_SYSTEM_EVENT(eventType, message string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent(eventType, encryptedMessage)
}

// ADJUST_EXECUTION_INTERVAL: Adjusts the interval for executing repetitive tasks
func ADJUST_EXECUTION_INTERVAL(taskID string, newInterval time.Duration) error {
    // Update task interval in ledger
    return common.ledger.UpdateTaskInterval(taskID, newInterval)
}

// SET_CONTRACT_LIFETIME: Set contract lifetime within the blockchain
func SET_CONTRACT_LIFETIME(contractID string, lifetime time.Duration) error {
    return common.ledger.SetContractLifetime(contractID, lifetime)
}

// UPDATE_CONTRACT_TIMEOUT: Updates the timeout for contract tasks
func UPDATE_CONTRACT_TIMEOUT(contractID string, newTimeout time.Duration) error {
    return common.ledger.UpdateTimeout(contractID, newTimeout)
}

// CHECK_TASK_STATUS: Check the status of a task using the ledger and sub-blocks
func CHECK_TASK_STATUS(taskID string) (string, error) {
    return QUERY_EXECUTION_STATUS(taskID)
}

// TRACK_MEMORY_USAGE: Logs and monitors memory usage over time
func TRACK_MEMORY_USAGE() error {
    usage, err := common.GetMemoryUsage()
    if err != nil {
        return err
    }
    LOG_SYSTEM_EVENT("Memory Usage", "Current memory usage: " + usage)
    return nil
}

// MANAGE_TEMPORARY_FILES: Manages temp files for contracts
func MANAGE_TEMPORARY_FILES(contractID string) error {
    tempDir := filepath.Join("temp", contractID)
    err := os.MkdirAll(tempDir, os.ModePerm)
    return err
}

// RETRIEVE_TEMPORARY_FILES: Retrieve files for a given contract from temp storage
func RETRIEVE_TEMPORARY_FILES(contractID string) ([]string, error) {
    tempDir := filepath.Join("temp", contractID)
    files, err := ioutil.ReadDir(tempDir)
    if err != nil {
        return nil, err
    }
    var fileNames []string
    for _, file := range files {
        fileNames = append(fileNames, file.Name())
    }
    return fileNames, nil
}

// DELETE_TEMPORARY_FILES: Deletes temporary files for a contract
func DELETE_TEMPORARY_FILES(contractID string) error {
    tempDir := filepath.Join("temp", contractID)
    return os.RemoveAll(tempDir)
}

// ENABLE_RESOURCE_SHARING: Enables resource sharing among tasks
func ENABLE_RESOURCE_SHARING() {
    resourceSharingEnabled = true
    LOG_SYSTEM_EVENT("Resource Sharing", "Resource sharing enabled")
}

// DISABLE_RESOURCE_SHARING: Disables resource sharing
func DISABLE_RESOURCE_SHARING() {
    resourceSharingEnabled = false
    LOG_SYSTEM_EVENT("Resource Sharing", "Resource sharing disabled")
}
