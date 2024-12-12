package utility

import (
	"sync"
	"synnergy_network/pkg/common"
	"time"
)

var systemFlags = make(map[string]bool)
var maintenanceMode bool
var lastHealthCheckStatus string
var networkStatus string
var configLock sync.Mutex

// GET_TIMESTAMP: Retrieves the current system timestamp
func GET_TIMESTAMP() time.Time {
    timestamp := time.Now()
    LogConfigOperation("GET_TIMESTAMP", "Timestamp retrieved: " + timestamp.String())
    return timestamp
}

// RETRIEVE_SYSTEM_CONFIG: Retrieves the current system configuration settings
func RETRIEVE_SYSTEM_CONFIG() map[string]interface{} {
    configLock.Lock()
    defer configLock.Unlock()
    config := map[string]interface{}{
        "systemFlags":    systemFlags,
        "debugMode":      debugMode,
        "maintenanceMode": maintenanceMode,
        "lastHealthCheckStatus": lastHealthCheckStatus,
        "networkStatus": networkStatus,
    }
    LogConfigOperation("RETRIEVE_SYSTEM_CONFIG", "System configuration retrieved")
    return config
}

// SET_SYSTEM_FLAG: Sets a specified system flag
func SET_SYSTEM_FLAG(flag string, value bool) {
    configLock.Lock()
    defer configLock.Unlock()
    systemFlags[flag] = value
    LogConfigOperation("SET_SYSTEM_FLAG", "Flag set: " + flag + " = " + boolToString(value))
}

// CLEAR_SYSTEM_FLAG: Clears a specified system flag
func CLEAR_SYSTEM_FLAG(flag string) {
    configLock.Lock()
    defer configLock.Unlock()
    delete(systemFlags, flag)
    LogConfigOperation("CLEAR_SYSTEM_FLAG", "Flag cleared: " + flag)
}

// DUMP_DIAGNOSTIC_INFO: Dumps diagnostic information for system analysis
func DUMP_DIAGNOSTIC_INFO() map[string]interface{} {
    diagnostics := map[string]interface{}{
        "timestamp":               GET_TIMESTAMP(),
        "lastHealthCheckStatus":   lastHealthCheckStatus,
        "systemFlags":             systemFlags,
        "maintenanceMode":         maintenanceMode,
        "debugMode":               debugMode,
        "networkStatus":           networkStatus,
    }
    LogConfigOperation("DUMP_DIAGNOSTIC_INFO", "Diagnostic information dumped")
    return diagnostics
}

// RETRIEVE_ERROR_LOG: Retrieves the error log entries from the ledger
func RETRIEVE_ERROR_LOG() ([]string, error) {
    errorLog, err := common.ledger.GetErrorLog()
    if err != nil {
        LogConfigOperation("RETRIEVE_ERROR_LOG", "Error retrieving log")
        return nil, err
    }
    LogConfigOperation("RETRIEVE_ERROR_LOG", "Error log retrieved")
    return errorLog, nil
}

// SET_DEBUG_MODE: Enables debug mode for the system
func SET_DEBUG_MODE() {
    configLock.Lock()
    defer configLock.Unlock()
    debugMode = true
    LogConfigOperation("SET_DEBUG_MODE", "Debug mode enabled")
}

// CLEAR_DEBUG_MODE: Disables debug mode for the system
func CLEAR_DEBUG_MODE() {
    configLock.Lock()
    defer configLock.Unlock()
    debugMode = false
    LogConfigOperation("CLEAR_DEBUG_MODE", "Debug mode disabled")
}

// ACTIVATE_MAINTENANCE_MODE: Activates maintenance mode for the system
func ACTIVATE_MAINTENANCE_MODE() {
    configLock.Lock()
    defer configLock.Unlock()
    maintenanceMode = true
    LogConfigOperation("ACTIVATE_MAINTENANCE_MODE", "Maintenance mode activated")
}

// DEACTIVATE_MAINTENANCE_MODE: Deactivates maintenance mode for the system
func DEACTIVATE_MAINTENANCE_MODE() {
    configLock.Lock()
    defer configLock.Unlock()
    maintenanceMode = false
    LogConfigOperation("DEACTIVATE_MAINTENANCE_MODE", "Maintenance mode deactivated")
}

// SYSTEM_RESET: Performs a system reset to restore initial configurations
func SYSTEM_RESET() {
    configLock.Lock()
    defer configLock.Unlock()
    systemFlags = make(map[string]bool)
    debugMode = false
    maintenanceMode = false
    lastHealthCheckStatus = "Not Checked"
    networkStatus = "Initializing"
    LogConfigOperation("SYSTEM_RESET", "System reset to initial configuration")
}

// RUN_HEALTH_CHECK: Executes a health check on the system and logs the result
func RUN_HEALTH_CHECK() {
    // Simulate health check process
    lastHealthCheckStatus = "Healthy"
    LogConfigOperation("RUN_HEALTH_CHECK", "Health check executed - Status: " + lastHealthCheckStatus)
}

// UPDATE_NETWORK_STATUS: Updates the current network status
func UPDATE_NETWORK_STATUS(status string) {
    configLock.Lock()
    defer configLock.Unlock()
    networkStatus = status
    LogConfigOperation("UPDATE_NETWORK_STATUS", "Network status updated to " + status)
}


// LogConfigOperation: Logs a configuration operation with encryption
func LogConfigOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("ConfigOperation", encryptedMessage)
}

