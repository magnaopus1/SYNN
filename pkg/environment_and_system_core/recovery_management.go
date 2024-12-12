package environment_and_system_core

import (
	"fmt"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// terminateRecoveryMode stops all recovery-related processes and logs the termination in the ledger.
func terminateRecoveryMode() error {
	// Stop recovery mode processes
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.RecoveryManager.TerminateRecoveryProcesses(); err != nil {
		return fmt.Errorf("failed to terminate recovery processes: %w", err)
	}

	// Record the recovery termination in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordRecoveryEvent("RecoveryModeTerminated"); err != nil {
		return fmt.Errorf("failed to record recovery termination in ledger: %w", err)
	}

	// Log success
	log.Println("Recovery mode terminated and recorded in ledger.")
	return nil
}

// startInterruptDiagnostics initiates diagnostics on interrupts and logs the action.
func startInterruptDiagnostics() error {
	// Start interrupt diagnostics in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.DiagnosticManager.StartDiagnostics("Interrupts"); err != nil {
		return fmt.Errorf("failed to start interrupt diagnostics: %w", err)
	}

	// Record the event in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordDiagnosticEvent("InterruptDiagnosticsStarted"); err != nil {
		return fmt.Errorf("failed to record interrupt diagnostics start in ledger: %w", err)
	}

	// Log success
	log.Println("Interrupt diagnostics started and recorded in ledger.")
	return nil
}

// stopInterruptDiagnostics stops diagnostics on interrupts and updates the ledger.
func stopInterruptDiagnostics() error {
	// Stop interrupt diagnostics in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.DiagnosticManager.StopDiagnostics("Interrupts"); err != nil {
		return fmt.Errorf("failed to stop interrupt diagnostics: %w", err)
	}

	// Record the event in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordDiagnosticEvent("InterruptDiagnosticsStopped"); err != nil {
		return fmt.Errorf("failed to record interrupt diagnostics stop in ledger: %w", err)
	}

	// Log success
	log.Println("Interrupt diagnostics stopped and recorded in ledger.")
	return nil
}

// retryOperationWithDelay retries a failed operation after a specified delay.
func retryOperationWithDelay(operationID string, delay time.Duration) error {
	// Validate input
	if operationID == "" {
		return fmt.Errorf("operationID cannot be empty")
	}
	if delay <= 0 {
		return fmt.Errorf("delay must be greater than 0")
	}

	// Log the retry action
	log.Printf("Retrying operation %s after a delay of %v...", operationID, delay)
	time.Sleep(delay)

	// Retry the operation
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.OperationManager.RetryOperation(operationID); err != nil {
		return fmt.Errorf("failed to retry operation %s: %w", operationID, err)
	}

	// Log success
	log.Printf("Operation %s retried successfully.", operationID)
	return nil
}

// querySystemHaltStatus retrieves the system halt status for monitoring purposes.
func querySystemHaltStatus() (string, error) {
	// Get the system halt status
	ledger := &ledger.Ledger{}
	status, err := ledger.EnvironmentSystemCoreLedger.SystemManager.GetHaltStatus()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve system halt status: %w", err)
	}

	// Encrypt the status
	encryptedStatus, err := common.EncryptData(status)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt system halt status: %w", err)
	}

	// Log success
	log.Println("System halt status retrieved and encrypted.")
	return encryptedStatus, nil
}

// enableRecoveryLogging starts logging recovery events and registers it in the ledger.
func enableRecoveryLogging() error {
	// Enable recovery logging in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.RecoveryManager.EnableLogging(); err != nil {
		return fmt.Errorf("failed to enable recovery logging: %w", err)
	}

	// Record the event in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordRecoveryEvent("RecoveryLoggingEnabled"); err != nil {
		return fmt.Errorf("failed to record recovery logging enablement in ledger: %w", err)
	}

	// Log success
	log.Println("Recovery logging enabled and recorded in ledger.")
	return nil
}

// disableRecoveryLogging stops logging recovery events and updates the ledger.
func disableRecoveryLogging() error {
	// Stop recovery logging in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.RecoveryManager.DisableLogging(); err != nil {
		return fmt.Errorf("failed to disable recovery logging: %w", err)
	}

	// Record the event in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordRecoveryEvent("RecoveryLoggingDisabled"); err != nil {
		return fmt.Errorf("failed to record recovery logging disablement in ledger: %w", err)
	}

	// Log success
	log.Println("Recovery logging disabled and recorded in ledger.")
	return nil
}

// getLastTrapLog retrieves the most recent trap event log and encrypts it for secure access.
func getLastTrapLog() (string, error) {
	// Retrieve the last trap log
	ledger := &ledger.Ledger{}
	lastLog, err := ledger.EnvironmentSystemCoreLedger.LogManager.GetLastTrapLog()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve last trap log: %w", err)
	}

	// Encrypt the log for secure access
	encryptedLog, err := common.EncryptData(lastLog)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt last trap log: %w", err)
	}

	// Log success
	log.Println("Last trap log retrieved and encrypted.")
	return encryptedLog, nil
}

// setPanicHandler establishes a handler for panic events and logs the configuration.
func setPanicHandler(handlerFunc func()) error {
	// Validate input
	if handlerFunc == nil {
		return fmt.Errorf("panic handler function cannot be nil")
	}

	// Set the panic handler in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.HandlerManager.SetPanicHandler(handlerFunc); err != nil {
		return fmt.Errorf("failed to set panic handler: %w", err)
	}

	// Record the event in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordHandlerEvent("PanicHandlerConfigured"); err != nil {
		return fmt.Errorf("failed to record panic handler configuration in ledger: %w", err)
	}

	// Log success
	log.Println("Panic handler set and recorded in ledger.")
	return nil
}

// clearPanicHandler removes the panic handler and logs the change.
func clearPanicHandler() error {
	// Clear the panic handler in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.HandlerManager.ClearPanicHandler(); err != nil {
		return fmt.Errorf("failed to clear panic handler: %w", err)
	}

	// Record the event in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordHandlerEvent("PanicHandlerRemoved"); err != nil {
		return fmt.Errorf("failed to record panic handler removal in ledger: %w", err)
	}

	// Log success
	log.Println("Panic handler cleared and recorded in ledger.")
	return nil
}

// verifyInterruptHandler checks that the interrupt handler is functioning as expected.
func verifyInterruptHandler() (bool, error) {
	// Verify the interrupt handler status
	ledger := &ledger.Ledger{}
	status, err := ledger.EnvironmentSystemCoreLedger.HandlerManager.VerifyInterruptHandler()
	if err != nil {
		return false, fmt.Errorf("failed to verify interrupt handler: %w", err)
	}

	// Log success
	log.Printf("Interrupt handler verification status: %v", status)
	return status, nil
}

// registerEmergencyOverride activates an emergency override and logs it in the ledger.
func registerEmergencyOverride(reason string) error {
	// Validate input
	if reason == "" {
		return fmt.Errorf("reason for emergency override cannot be empty")
	}

	// Activate the emergency override in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.OverrideManager.ActivateOverride(reason); err != nil {
		return fmt.Errorf("failed to activate emergency override: %w", err)
	}

	// Record the event in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordOverrideEvent("Activated", reason); err != nil {
		return fmt.Errorf("failed to record emergency override activation in ledger: %w", err)
	}

	// Log success
	log.Printf("Emergency override registered: %s", reason)
	return nil
}

// cancelEmergencyOverride deactivates the emergency override and updates the ledger.
func cancelEmergencyOverride() error {
	// Deactivate the emergency override in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.OverrideManager.DeactivateOverride(); err != nil {
		return fmt.Errorf("failed to deactivate emergency override: %w", err)
	}

	// Record the event in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordOverrideEvent("Cancelled", ""); err != nil {
		return fmt.Errorf("failed to record emergency override cancellation in ledger: %w", err)
	}

	// Log success
	log.Println("Emergency override cancelled and recorded in ledger.")
	return nil
}

// queryPanicStatus checks the current panic status in the system for recovery assessment.
func queryPanicStatus() (bool, error) {
	// Retrieve the panic status from the system
	ledger := &ledger.Ledger{}
	status, err := ledger.EnvironmentSystemCoreLedger.SystemManager.GetPanicStatus()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve panic status: %w", err)
	}

	// Log the status
	log.Printf("Panic status retrieved: %v", status)
	return status, nil
}

// cleanupFailedOperation clears resources and logs for a specific failed operation.
func cleanupFailedOperation(operationID string) error {
	// Validate input
	if operationID == "" {
		return fmt.Errorf("operation ID cannot be empty")
	}

	// Perform cleanup in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.CleanupManager.CleanupOperation(operationID); err != nil {
		return fmt.Errorf("failed to clean up resources for operation %s: %w", operationID, err)
	}

	// Record the cleanup event in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordCleanupEvent(operationID); err != nil {
		return fmt.Errorf("failed to record cleanup for operation %s in ledger: %w", operationID, err)
	}

	// Log success
	log.Printf("Failed operation %s cleaned up and recorded in ledger.", operationID)
	return nil
}

// initiateGracefulShutdown initiates a safe shutdown sequence and logs each critical step.
func initiateGracefulShutdown() error {
	// Start graceful shutdown in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.ShutdownManager.InitiateShutdownSequence(); err != nil {
		return fmt.Errorf("failed to initiate graceful shutdown: %w", err)
	}

	// Record the shutdown event in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordShutdownEvent("GracefulShutdownInitiated"); err != nil {
		return fmt.Errorf("failed to record shutdown event in ledger: %w", err)
	}

	// Log success
	log.Println("Graceful shutdown initiated and recorded in ledger.")
	return nil
}
