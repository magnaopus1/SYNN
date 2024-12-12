package environment_and_system_core

import (
	"fmt"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// PauseProcess pauses the specified process and records the state in the ledger.
func PauseProcess(processID string, ledger *ledger.Ledger) error {
	// Validate input
	if processID == "" {
		return fmt.Errorf("processID cannot be empty")
	}

	// Pause the process in the system
	if err := ledger.EnvironmentSystemCoreLedger.ProcessManager.Pause(processID); err != nil {
		return fmt.Errorf("failed to pause process %s: %w", processID, err)
	}

	// Record the paused state in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordProcessState(processID, "paused", 0); err != nil {
		return fmt.Errorf("failed to record process state for process %s in ledger: %w", processID, err)
	}

	// Log success
	log.Printf("Process %s has been paused and recorded in ledger.", processID)
	return nil
}


// ResumeProcess resumes the specified process and records the state in the ledger.
func ResumeProcess(processID string, ledger *ledger.Ledger) error {
	// Validate input
	if processID == "" {
		return fmt.Errorf("processID cannot be empty")
	}

	// Resume the process in the system
	if err := ledger.EnvironmentSystemCoreLedger.ProcessManager.Resume(processID); err != nil {
		return fmt.Errorf("failed to resume process %s: %w", processID, err)
	}

	// Record the running state in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordProcessState(processID, "running", 0); err != nil {
		return fmt.Errorf("failed to record process state for process %s in ledger: %w", processID, err)
	}

	// Log success
	log.Printf("Process %s has been resumed and recorded in ledger.", processID)
	return nil
}


// SetProcessTimeout sets a timeout for the specified process and records the update in the ledger.
func SetProcessTimeout(processID string, duration time.Duration, ledger *ledger.Ledger) error {
	// Validate input
	if processID == "" {
		return fmt.Errorf("processID cannot be empty")
	}
	if duration <= 0 {
		return fmt.Errorf("timeout duration must be greater than zero")
	}

	// Set the timeout in the system
	if err := ledger.EnvironmentSystemCoreLedger.ProcessManager.SetTimeout(processID, duration); err != nil {
		return fmt.Errorf("failed to set timeout for process %s: %w", processID, err)
	}

	// Record the timeout in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordProcessState(processID, "running", duration); err != nil {
		return fmt.Errorf("failed to record process timeout for process %s in ledger: %w", processID, err)
	}

	// Log success
	log.Printf("Timeout set for process %s with duration %v and recorded in ledger.", processID, duration)
	return nil
}


// GetProcessInfo retrieves process information and encrypts it for secure transmission.
func GetProcessInfo(processID string) (string, error) {
	// Validate input
	if processID == "" {
		return "", fmt.Errorf("processID cannot be empty")
	}

	// Retrieve process information
	ledger := &ledger.Ledger{}
	info, err := ledger.EnvironmentSystemCoreLedger.ProcessManager.GetInfo(processID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve process info for %s: %w", processID, err)
	}

	// Encrypt the information
	encryptedInfo, err := encryption.Encrypt(info)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt process info for %s: %w", processID, err)
	}

	// Log success
	log.Printf("Process info for %s retrieved and encrypted.", processID)
	return encryptedInfo, nil
}


// AllocateSystemResource allocates system resources and records the allocation in the ledger.
func AllocateSystemResource(resourceType string, quantity int, ledger *ledger.Ledger) (string, error) {
	// Validate input
	if resourceType == "" {
		return "", fmt.Errorf("resourceType cannot be empty")
	}
	if quantity <= 0 {
		return "", fmt.Errorf("quantity must be greater than zero")
	}

	// Allocate the resource in the system
	resourceID, err := ledger.EnvironmentSystemCoreLedger.ResourceManager.Allocate(resourceType, quantity)
	if err != nil {
		return "", fmt.Errorf("failed to allocate resource of type %s: %w", resourceType, err)
	}

	// Record the allocation in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordResourceAllocation(resourceType, resourceID, quantity); err != nil {
		return "", fmt.Errorf("failed to record resource allocation for %s in ledger: %w", resourceID, err)
	}

	// Log success
	log.Printf("Resource of type %s allocated with ID %s and quantity %d, recorded in ledger.", resourceType, resourceID, quantity)
	return resourceID, nil
}


// ReleaseSystemResource releases a system resource and updates the ledger.
func ReleaseSystemResource(resourceID string, ledger *ledger.Ledger) error {
	// Validate input
	if resourceID == "" {
		return fmt.Errorf("resourceID cannot be empty")
	}

	// Release the resource in the system
	if err := ledger.EnvironmentSystemCoreLedger.ResourceManager.Release(resourceID); err != nil {
		return fmt.Errorf("failed to release resource %s: %w", resourceID, err)
	}

	// Remove the resource allocation from the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RemoveResourceAllocation(resourceID); err != nil {
		return fmt.Errorf("failed to remove resource allocation for %s in ledger: %w", resourceID, err)
	}

	// Log success
	log.Printf("Resource %s has been released and recorded in ledger.", resourceID)
	return nil
}



// MonitorResourceAllocation generates a report on resource allocation, encrypts it, and logs it in the ledger.
func MonitorResourceAllocation(ledger *ledger.Ledger) (string, error) {
	// Generate the resource monitoring report
	report, err := ledger.EnvironmentSystemCoreLedger.ResourceManager.GenerateAllocationReport()
	if err != nil {
		return "", fmt.Errorf("failed to generate resource monitoring report: %w", err)
	}

	// Encrypt the report
	encryptedReport, err := common.EncryptData(report)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt resource monitoring report: %w", err)
	}

	// Log the encrypted report in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.LogResourceMonitor(encryptedReport); err != nil {
		return "", fmt.Errorf("failed to log resource monitoring report in ledger: %w", err)
	}

	// Log success
	log.Println("Resource monitoring report generated, encrypted, and logged in ledger.")
	return report, nil
}



// LockSystemState locks the system state for a specified reason and logs the update in the ledger.
func LockSystemState(reason string, ledger *ledger.Ledger) error {
	// Validate input
	if reason == "" {
		return fmt.Errorf("lock reason cannot be empty")
	}

	// Update the system state to locked
	if err := ledger.EnvironmentSystemCoreLedger.SystemManager.LockState(reason); err != nil {
		return fmt.Errorf("failed to lock system state: %w", err)
	}

	// Record the lock state in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.UpdateSystemLockState(reason, true); err != nil {
		return fmt.Errorf("failed to update system lock state in ledger: %w", err)
	}

	// Log success
	log.Printf("System state locked with reason: %s and recorded in ledger.", reason)
	return nil
}



// UnlockSystemState unlocks the system state and updates the ledger.
func UnlockSystemState(ledger *ledger.Ledger) error {
	// Update the system state to unlocked
	if err := ledger.EnvironmentSystemCoreLedger.SystemManager.UnlockState(); err != nil {
		return fmt.Errorf("failed to unlock system state: %w", err)
	}

	// Record the unlock state in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.UpdateSystemLockState("", false); err != nil {
		return fmt.Errorf("failed to update system unlock state in ledger: %w", err)
	}

	// Log success
	log.Println("System state unlocked and recorded in ledger.")
	return nil
}


// DefineSystemConstant defines a new system constant and logs it in the ledger.
func DefineSystemConstant(name string, value string, ledger *ledger.Ledger) error {
	// Validate input
	if name == "" || value == "" {
		return fmt.Errorf("constant name and value cannot be empty")
	}

	// Encrypt the constant value
	encryptedValue, err := common.EncryptData(value)
	if err != nil {
		return fmt.Errorf("failed to encrypt system constant value for %s: %w", name, err)
	}

	// Log the constant definition in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.LogConstantDefinition(name, encryptedValue); err != nil {
		return fmt.Errorf("failed to log system constant definition for %s in ledger: %w", name, err)
	}

	// Log success
	log.Printf("System constant %s defined and recorded in ledger.", name)
	return nil
}

// UpdateSystemConstant updates the value of an existing system constant and logs the change in the ledger.
func UpdateSystemConstant(name string, newValue string, ledger *ledger.Ledger) error {
	// Validate input
	if name == "" || newValue == "" {
		return fmt.Errorf("constant name and new value cannot be empty")
	}

	// Encrypt the new constant value
	encryptedValue, err := common.EncryptData(newValue)
	if err != nil {
		return fmt.Errorf("failed to encrypt new system constant value for %s: %w", name, err)
	}

	// Log the constant update in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.LogConstantUpdate(name, encryptedValue); err != nil {
		return fmt.Errorf("failed to log system constant update for %s in ledger: %w", name, err)
	}

	// Log success
	log.Printf("System constant %s updated and recorded in ledger.", name)
	return nil
}

// ResetSystemConstants resets all system constants to their default values and logs the event in the ledger.
func ResetSystemConstants(ledger *ledger.Ledger) error {
	// Reset constants in the system
	if err := ledger.EnvironmentSystemCoreLedger.ConstantsManager.ResetToDefaults(); err != nil {
		return fmt.Errorf("failed to reset system constants: %w", err)
	}

	// Log the reset in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.LogConstantsReset(); err != nil {
		return fmt.Errorf("failed to log system constants reset: %w", err)
	}

	// Log success
	log.Println("System constants have been reset to default and logged in ledger.")
	return nil
}

// CheckSystemSecurity performs a security check and logs the encrypted report in the ledger.
func CheckSystemSecurity(ledger *ledger.Ledger) (string, error) {
	// Perform security check
	securityReport, err := ledger.EnvironmentSystemCoreLedger.SecurityManager.GenerateReport()
	if err != nil {
		return "", fmt.Errorf("failed to generate security report: %w", err)
	}

	// Encrypt the report
	encryptedReport, err := common.EncryptData(securityReport)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt security report: %w", err)
	}

	// Log the report in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.LogSecurityCheck(encryptedReport); err != nil {
		return "", fmt.Errorf("failed to log security check in ledger: %w", err)
	}

	// Log success
	log.Println("Security check performed, report encrypted, and logged in ledger.")
	return securityReport, nil
}

// LogSecurityEvent logs a security event with encrypted details in the ledger.
func LogSecurityEvent(eventType string, details string, ledger *ledger.Ledger) error {
	// Validate input
	if eventType == "" || details == "" {
		return fmt.Errorf("event type and details cannot be empty")
	}

	// Encrypt the event details
	encryptedDetails, err := common.EncryptData(details)
	if err != nil {
		return fmt.Errorf("failed to encrypt security event details: %w", err)
	}

	// Log the event in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.LogSecurityEvent(eventType, encryptedDetails); err != nil {
		return fmt.Errorf("failed to log security event of type %s in ledger: %w", eventType, err)
	}

	// Log success
	log.Printf("Security event logged: %s.", eventType)
	return nil
}

// QuerySecurityLog retrieves security events of a specified type from the ledger.
func QuerySecurityLog(eventType string, ledger *ledger.Ledger) (string, error) {
	// Validate input
	if eventType == "" {
		return "", fmt.Errorf("event type cannot be empty")
	}

	// Retrieve events from the ledger
	events, err := ledger.EnvironmentSystemCoreLedger.RetrieveSecurityEvents(eventType)
	if err != nil {
		return "", fmt.Errorf("failed to query security log for event type %s: %w", eventType, err)
	}

	// Log success
	log.Printf("Security log queried for event type: %s.", eventType)
	return events, nil
}

// PerformSystemMaintenance performs system maintenance and logs the event in the ledger.
func PerformSystemMaintenance(ledger *ledger.Ledger) error {
	// Execute system maintenance tasks
	if err := ledger.EnvironmentSystemCoreLedger.MaintenanceManager.ExecuteTasks(); err != nil {
		return fmt.Errorf("failed to perform system maintenance: %w", err)
	}

	// Log the maintenance event in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.LogMaintenanceEvent("SystemMaintenanceCompleted"); err != nil {
		return fmt.Errorf("failed to log maintenance event in ledger: %w", err)
	}

	// Log success
	log.Println("System maintenance completed and recorded in ledger.")
	return nil
}

