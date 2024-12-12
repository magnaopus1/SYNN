package environment_and_system_core

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

var (
	currentLoggingLevel string
	mu                  sync.Mutex // Ensures thread safety for logging level changes
)

type SystemHookRegistry struct {
	hooks    map[string]func() error
	priorities map[string]int
	mu       sync.Mutex
}



// GetSystemLoadAverage retrieves the average system load over a specified period.
func getSystemLoadAverage(duration time.Duration) (float64, error) {
	// Validate the duration input
	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration: must be greater than zero, got %s", duration)
	}

	// Read the load average from the system's /proc/loadavg file (Linux-based systems)
	data, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, fmt.Errorf("failed to read system load average: %w", err)
	}

	// Parse the first value from the file, which represents the 1-minute load average
	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return 0, fmt.Errorf("unexpected format in /proc/loadavg: %s", string(data))
	}

	loadAverage, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse load average: %w", err)
	}

	// Log the load average
	log.Printf("Retrieved system load average over %s: %.2f\n", duration, loadAverage)

	// Record the load average metric in the ledger
	ledger := &ledger.Ledger{}
	err = ledger.EnvironmentSystemCoreLedger.RecordSystemMetric("LoadAverage", loadAverage)
	if err != nil {
		return 0, fmt.Errorf("failed to record system load average in ledger: %w", err)
	}

	return loadAverage, nil
}

// EnableLogging enables system-wide logging.
func enableLogging() error {
	// Enable the logging system
	log.Println("Enabling system-wide logging.")

	// Record the event in the ledger
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("LoggingEnabled", "System-wide logging has been enabled.")
	if err != nil {
		return fmt.Errorf("failed to record logging enable event in ledger: %w", err)
	}

	log.Println("System-wide logging successfully enabled.")
	return nil
}

// DisableLogging disables system-wide logging.
func disableLogging() error {
	// Disable the logging system
	log.Println("Disabling system-wide logging.")

	// Record the event in the ledger
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("LoggingDisabled", "System-wide logging has been disabled.")
	if err != nil {
		return fmt.Errorf("failed to record logging disable event in ledger: %w", err)
	}

	log.Println("System-wide logging successfully disabled.")
	return nil
}


// SetLoggingLevel sets the logging level for system operations.
func setLoggingLevel(level string) error {
	// Define valid logging levels
	validLevels := map[string]struct{}{
		"DEBUG": {},
		"INFO":  {},
		"WARN":  {},
		"ERROR": {},
	}

	// Normalize and validate the logging level
	normalizedLevel := strings.ToUpper(level)
	if _, exists := validLevels[normalizedLevel]; !exists {
		return fmt.Errorf("invalid logging level: %s", level)
	}

	// Update the global logging level with thread safety
	mu.Lock()
	defer mu.Unlock()
	currentLoggingLevel = normalizedLevel

	// Record the logging level change in the ledger
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent(
		"LoggingLevelSet",
		fmt.Sprintf("Logging level set to %s.", normalizedLevel),
	)
	if err != nil {
		return fmt.Errorf("failed to record logging level set event in ledger: %w", err)
	}

	// Log the successful application of the logging level
	log.Printf("Logging level successfully set to %s.\n", normalizedLevel)

	return nil
}


// GetCurrentLoggingLevel retrieves the current system logging level.
func getCurrentLoggingLevel() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	if currentLoggingLevel == "" {
		return "", fmt.Errorf("logging level is not set")
	}

	log.Printf("Current logging level retrieved: %s", currentLoggingLevel)
	return currentLoggingLevel, nil
}


// CheckLoggingLevel retrieves and logs the current logging level.
func checkLoggingLevel() (string, error) {
	level, err := getCurrentLoggingLevel()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve logging level: %w", err)
	}

	log.Printf("Current logging level: %s", level)
	return level, nil
}


// RegisterHook registers a system hook for a specified event.
func registerHook(eventID string, handlerFunc func() error) error {
	// Validate input
	if eventID == "" {
		return fmt.Errorf("eventID cannot be empty")
	}
	if handlerFunc == nil {
		return fmt.Errorf("handler function cannot be nil for eventID: %s", eventID)
	}

	// Register the hook in the system's memory or handler registry
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.SystemHookRegistry.Add(eventID, handlerFunc)
	if err != nil {
		return fmt.Errorf("failed to register hook for event %s: %w", eventID, err)
	}

	// Record the event in the ledger
	err = ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("HookRegistered", fmt.Sprintf("EventID: %s", eventID))
	if err != nil {
		return fmt.Errorf("failed to record hook registration for event %s in ledger: %w", eventID, err)
	}

	log.Printf("Hook successfully registered for event: %s", eventID)
	return nil
}


// UnregisterHook removes a system hook for a specified event.
func unregisterHook(eventID string) error {
	// Validate input
	if eventID == "" {
		return fmt.Errorf("eventID cannot be empty")
	}

	// Remove the hook from the system's memory or handler registry
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.SystemHookRegistry.Remove(eventID)
	if err != nil {
		return fmt.Errorf("failed to unregister hook for event %s: %w", eventID, err)
	}

	// Record the event in the ledger
	err = ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("HookUnregistered", fmt.Sprintf("EventID: %s", eventID))
	if err != nil {
		return fmt.Errorf("failed to record hook unregistration for event %s in ledger: %w", eventID, err)
	}

	log.Printf("Hook successfully unregistered for event: %s", eventID)
	return nil
}


// SetHookPriority sets the priority for a specific system hook.
func setHookPriority(eventID string, priority int) error {
	// Validate input
	if eventID == "" {
		return fmt.Errorf("eventID cannot be empty")
	}
	if priority < 0 {
		return fmt.Errorf("priority must be a non-negative integer, got %d", priority)
	}

	// Update the priority in the system's hook registry
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.SystemHookRegistry.SetPriority(eventID, priority)
	if err != nil {
		return fmt.Errorf("failed to set priority for event %s: %w", eventID, err)
	}

	// Record the priority change in the ledger
	err = ledger.EnvironmentSystemCoreLedger.AddSystemHook(eventID, priority)
	if err != nil {
		return fmt.Errorf("failed to record priority change for event %s in ledger: %w", eventID, err)
	}

	log.Printf("Hook priority successfully set for event %s: %d", eventID, priority)
	return nil
}


// ClearHookPriority clears the priority for a specific system hook.
func clearHookPriority(eventID string) error {
	// Validate input
	if eventID == "" {
		return fmt.Errorf("eventID cannot be empty")
	}

	// Clear the hook priority in the system's hook registry
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.SystemHookRegistry.RemovePriority(eventID)
	if err != nil {
		return fmt.Errorf("failed to clear hook priority for event %s: %w", eventID, err)
	}

	// Record the removal in the ledger
	err = ledger.EnvironmentSystemCoreLedger.RemoveSystemHook(eventID)
	if err != nil {
		return fmt.Errorf("failed to remove hook record from ledger for event %s: %w", eventID, err)
	}

	log.Printf("Hook priority successfully cleared for event: %s", eventID)
	return nil
}


// SyncSystemState synchronizes the current system state with all nodes.
func syncSystemState() error {
	// Perform synchronization across nodes
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.SystemStateSynchronizer.SyncAllNodes()
	if err != nil {
		return fmt.Errorf("failed to synchronize system state across nodes: %w", err)
	}

	// Record the synchronization in the ledger
	err = ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("SystemStateSynchronized", "System state synchronized across all nodes")
	if err != nil {
		return fmt.Errorf("failed to record system synchronization event in ledger: %w", err)
	}

	log.Println("System state successfully synchronized across all nodes.")
	return nil
}


// FlushSystemCache clears the system cache.
func flushSystemCache() error {
	// Clear the system cache
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.SystemCacheManager.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush system cache: %w", err)
	}

	// Record the cache flush event in the ledger
	err = ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("SystemCacheFlushed", "System cache cleared successfully")
	if err != nil {
		return fmt.Errorf("failed to record system cache flush event in ledger: %w", err)
	}

	log.Println("System cache successfully flushed.")
	return nil
}



// DefineSystemRole defines a new role within the system with specific permissions.
func DefineSystemRole(roleName string, permissions []string) error {
	// Validate input
	if roleName == "" {
		return fmt.Errorf("roleName cannot be empty")
	}
	if len(permissions) == 0 {
		return fmt.Errorf("permissions cannot be empty for role: %s", roleName)
	}

	// Add the role to the system
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.RoleManager.AddRole(roleName, permissions)
	if err != nil {
		return fmt.Errorf("failed to define system role %s: %w", roleName, err)
	}

	// Record the role in the ledger
	err = ledger.EnvironmentSystemCoreLedger.AddRoleRecord(roleName, permissions)
	if err != nil {
		return fmt.Errorf("failed to record system role %s in ledger: %w", roleName, err)
	}

	log.Printf("System role %s successfully defined with permissions %v.", roleName, permissions)
	return nil
}

// RevokeSystemRole revokes a specified system role, removing its permissions.
func revokeSystemRole(roleName string) error {
	// Validate input
	if roleName == "" {
		return fmt.Errorf("roleName cannot be empty")
	}

	// Revoke the role in the system
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.RoleManager.RemoveRole(roleName)
	if err != nil {
		return fmt.Errorf("failed to revoke system role %s: %w", roleName, err)
	}

	// Remove the role record from the ledger
	err = ledger.EnvironmentSystemCoreLedger.RemoveRoleRecord(roleName)
	if err != nil {
		return fmt.Errorf("failed to remove role record for %s from ledger: %w", roleName, err)
	}

	log.Printf("System role %s successfully revoked.", roleName)
	return nil
}


// AssignRolePermissions assigns permissions to a specific role.
func assignRolePermissions(roleName string, permissions []string) error {
	// Validate input
	if roleName == "" {
		return fmt.Errorf("roleName cannot be empty")
	}
	if len(permissions) == 0 {
		return fmt.Errorf("permissions cannot be empty for role: %s", roleName)
	}

	// Assign permissions to the role in the system
	ledger := &ledger.Ledger{}
	err := ledger.EnvironmentSystemCoreLedger.RoleManager.UpdateRolePermissions(roleName, permissions)
	if err != nil {
		return fmt.Errorf("failed to assign permissions to role %s: %w", roleName, err)
	}

	// Update the permissions in the ledger
	err = ledger.EnvironmentSystemCoreLedger.UpdateRolePermissions(roleName, permissions)
	if err != nil {
		return fmt.Errorf("failed to record permission updates for role %s in ledger: %w", roleName, err)
	}

	log.Printf("Permissions %v successfully assigned to role %s.", permissions, roleName)
	return nil
}


// QueryRolePermissions retrieves the permissions associated with a specified role.
func queryRolePermissions(roleName string) ([]string, error) {
	// Validate input
	if roleName == "" {
		return nil, fmt.Errorf("roleName cannot be empty")
	}

	// Query the permissions from the ledger
	ledger := &ledger.Ledger{}
	ledgerPermissions, err := ledger.EnvironmentSystemCoreLedger.QueryRoleRecord(roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to query role record for %s: %w", roleName, err)
	}

	// Ensure permissions are not nil
	if ledgerPermissions == nil {
		return nil, fmt.Errorf("no permissions found for role %s", roleName)
	}

	log.Printf("Permissions queried for role %s: %v", roleName, ledgerPermissions)
	return ledgerPermissions, nil
}


// LoadSystemProfile loads the system profile and configurations into memory.
func loadSystemProfile(profileID string) error {
	// Validate input
	if profileID == "" {
		return fmt.Errorf("profileID cannot be empty")
	}

	// Check if the profile exists in the ledger
	ledger := &ledger.Ledger{}
	exists, err := ledger.EnvironmentSystemCoreLedger.CheckProfileExists(profileID)
	if err != nil {
		return fmt.Errorf("failed to check profile existence for %s: %w", profileID, err)
	}
	if !exists {
		return fmt.Errorf("profile with ID %s does not exist", profileID)
	}

	// Load the profile configurations into memory
	err = ledger.EnvironmentSystemCoreLedger.SystemProfileManager.Load(profileID)
	if err != nil {
		return fmt.Errorf("failed to load system profile %s: %w", profileID, err)
	}

	// Record the event in the ledger
	err = ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("SystemProfileLoaded", fmt.Sprintf("ProfileID: %s", profileID))
	if err != nil {
		return fmt.Errorf("failed to record system profile load event in ledger: %w", err)
	}

	log.Printf("System profile successfully loaded: %s", profileID)
	return nil
}



// Add registers a new hook.
func (r *SystemHookRegistry) Add(eventID string, handler func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.hooks[eventID]; exists {
		return fmt.Errorf("hook already exists for eventID: %s", eventID)
	}

	r.hooks[eventID] = handler
	r.priorities[eventID] = 0 // Default priority
	return nil
}

// Remove unregisters an existing hook.
func (r *SystemHookRegistry) Remove(eventID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.hooks[eventID]; !exists {
		return fmt.Errorf("no hook found for eventID: %s", eventID)
	}

	delete(r.hooks, eventID)
	delete(r.priorities, eventID)
	return nil
}

// SetPriority updates the priority of an existing hook.
func (r *SystemHookRegistry) SetPriority(eventID string, priority int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.hooks[eventID]; !exists {
		return fmt.Errorf("no hook found for eventID: %s", eventID)
	}

	r.priorities[eventID] = priority
	return nil
}