package environment_and_system_core

import (
	"fmt"
	"log"
	"synnergy_network/pkg/ledger"
	"time"
)


// raiseTrap sets a critical system trap condition using the TrapManager and logs the event in the ledger.
func raiseTrap(ledgerInstance *ledger.Ledger, errorCode int, message string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("raiseTrap: ledgerInstance cannot be nil")
	}
	if errorCode <= 0 {
		return fmt.Errorf("raiseTrap: errorCode must be greater than 0")
	}
	if message == "" {
		return fmt.Errorf("raiseTrap: message cannot be empty")
	}

	// Parameters for the trap condition
	trapParams := map[string]interface{}{
		"ErrorCode": errorCode,
		"Message":   message,
		"Timestamp": time.Now(),
	}

	// Use the TrapManager to set the trap condition
	if err := ledgerInstance.EnvironmentSystemCoreLedger.TrapManager.SetCondition("CriticalSystemTrap", trapParams); err != nil {
		return fmt.Errorf("raiseTrap: failed to set critical system trap: %w", err)
	}

	// Record the trap condition in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordSystemEvent("CriticalSystemTrapRaised", trapParams); err != nil {
		return fmt.Errorf("raiseTrap: failed to record system trap event in ledger: %w", err)
	}

	// Log the trap and halt the system
	log.Fatalf("Critical System Trap Raised: Code %d - %s", errorCode, message)
	return nil
}

// handleInterrupt processes an interrupt using its registered handler and records the event in the ledger.
func handleInterrupt(ledgerInstance *ledger.Ledger, interruptID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("handleInterrupt: ledgerInstance cannot be nil")
	}
	if interruptID == "" {
		return fmt.Errorf("handleInterrupt: interruptID cannot be empty")
	}

	// Retrieve and execute the handler using InterruptManager
	if err := ledgerInstance.EnvironmentSystemCoreLedger.InterruptManager.ExecuteHandler(interruptID); err != nil {
		return fmt.Errorf("handleInterrupt: failed to handle interrupt ID %s: %w", interruptID, err)
	}

	// Record the interrupt handling event in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordSystemEvent("InterruptHandled", map[string]interface{}{"InterruptID": interruptID}); err != nil {
		return fmt.Errorf("handleInterrupt: failed to record interrupt handling event: %w", err)
	}

	// Log success
	log.Printf("Interrupt %s handled successfully and recorded in ledger.", interruptID)
	return nil
}


// registerInterruptHandler registers a handler for a specific interrupt ID and ensures it is recorded in the ledger.
func registerInterruptHandler(ledgerInstance *ledger.Ledger, interruptID string, handler func() error) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("registerInterruptHandler: ledgerInstance cannot be nil")
	}
	if interruptID == "" {
		return fmt.Errorf("registerInterruptHandler: interruptID cannot be empty")
	}
	if handler == nil {
		return fmt.Errorf("registerInterruptHandler: handler cannot be nil")
	}

	// Register the handler using InterruptManager
	if err := ledgerInstance.EnvironmentSystemCoreLedger.InterruptManager.RegisterHandler(interruptID, handler); err != nil {
		return fmt.Errorf("registerInterruptHandler: failed to register handler for interrupt ID %s: %w", interruptID, err)
	}

	// Record the registration in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordSystemEvent("InterruptHandlerRegistered", map[string]interface{}{"InterruptID": interruptID}); err != nil {
		return fmt.Errorf("registerInterruptHandler: failed to record handler registration in ledger: %w", err)
	}

	// Log success
	log.Printf("Handler for interrupt ID %s successfully registered and recorded in ledger.", interruptID)
	return nil
}


// clearInterruptHandler removes a handler for a specific interrupt ID and records the event in the ledger.
func clearInterruptHandler(ledgerInstance *ledger.Ledger, interruptID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("clearInterruptHandler: ledgerInstance cannot be nil")
	}
	if interruptID == "" {
		return fmt.Errorf("clearInterruptHandler: interruptID cannot be empty")
	}

	// Deregister the handler using InterruptManager
	if err := ledgerInstance.EnvironmentSystemCoreLedger.InterruptManager.DeregisterHandler(interruptID); err != nil {
		return fmt.Errorf("clearInterruptHandler: failed to deregister handler for interrupt ID %s: %w", interruptID, err)
	}

	// Record the deregistration in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordSystemEvent("InterruptHandlerDeregistered", map[string]interface{}{"InterruptID": interruptID}); err != nil {
		return fmt.Errorf("clearInterruptHandler: failed to record handler deregistration in ledger: %w", err)
	}

	// Log success
	log.Printf("Handler for interrupt ID %s successfully deregistered and recorded in ledger.", interruptID)
	return nil
}


// systemHalt securely stops all operations using ShutdownManager and logs the halt event in the ledger.
func systemHalt(ledgerInstance *ledger.Ledger, reason string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("systemHalt: ledgerInstance cannot be nil")
	}
	if reason == "" {
		return fmt.Errorf("systemHalt: reason cannot be empty")
	}

	// Perform the halt operation
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ShutdownManager.InitiateHalt(reason); err != nil {
		return fmt.Errorf("systemHalt: failed to halt system: %w", err)
	}

	// Record the halt event in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordSystemEvent("SystemHalt", map[string]interface{}{"Reason": reason}); err != nil {
		return fmt.Errorf("systemHalt: failed to record system halt event in ledger: %w", err)
	}

	// Log success
	log.Fatalf("System halt initiated due to: %s", reason)
	return nil
}


// logException records an exception in the ledger for further analysis.
func logException(ledgerInstance *ledger.Ledger, exceptionID string, description string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("logException: ledgerInstance cannot be nil")
	}
	if exceptionID == "" {
		return fmt.Errorf("logException: exceptionID cannot be empty")
	}
	if description == "" {
		return fmt.Errorf("logException: description cannot be empty")
	}

	// Record the exception using LogManager
	exceptionEntry := map[string]interface{}{
		"ExceptionID": exceptionID,
		"Description": description,
		"Timestamp":   time.Now(),
	}
	if err := ledgerInstance.EnvironmentSystemCoreLedger.LogManager.RecordException(exceptionEntry); err != nil {
		return fmt.Errorf("logException: failed to record exception: %w", err)
	}

	// Log success
	log.Printf("Exception %s logged successfully in ledger.", exceptionID)
	return nil
}


// retryFailedOperation retries a failed operation using the OperationManager, up to a maximum number of attempts.
func retryFailedOperation(ledgerInstance *ledger.Ledger, operationID string, maxAttempts int) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("retryFailedOperation: ledgerInstance cannot be nil")
	}
	if operationID == "" {
		return fmt.Errorf("retryFailedOperation: operationID cannot be empty")
	}
	if maxAttempts <= 0 {
		return fmt.Errorf("retryFailedOperation: maxAttempts must be greater than 0")
	}

	// Retry the operation up to the maximum attempts
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ledgerInstance.EnvironmentSystemCoreLedger.OperationManager.RetryOperation(operationID); err == nil {
			// Record the successful retry in the ledger
			if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordSystemEvent("OperationRetried", map[string]interface{}{"OperationID": operationID, "Attempt": attempt}); err != nil {
				return fmt.Errorf("retryFailedOperation: failed to record retry success: %w", err)
			}
			log.Printf("Operation %s successfully retried on attempt %d.", operationID, attempt)
			return nil
		}
		time.Sleep(time.Duration(attempt) * time.Second) // Backoff
	}

	return fmt.Errorf("retryFailedOperation: all retry attempts failed for operation %s", operationID)
}


// emergencyShutdown performs an emergency system shutdown using ShutdownManager and records the event in the ledger.
func emergencyShutdown(ledgerInstance *ledger.Ledger, reason string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("emergencyShutdown: ledgerInstance cannot be nil")
	}
	if reason == "" {
		return fmt.Errorf("emergencyShutdown: reason cannot be empty")
	}

	// Perform the emergency shutdown
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ShutdownManager.InitiateEmergencyShutdown(reason); err != nil {
		return fmt.Errorf("emergencyShutdown: failed to initiate emergency shutdown: %w", err)
	}

	// Record the shutdown event in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordSystemEvent("EmergencyShutdown", map[string]interface{}{"Reason": reason}); err != nil {
		return fmt.Errorf("emergencyShutdown: failed to record emergency shutdown event in ledger: %w", err)
	}

	// Log success
	log.Fatalf("Emergency shutdown initiated due to: %s", reason)
	return nil
}



// panicRecovery initiates recovery procedures after a system panic using the RecoveryManager.
func panicRecovery(ledger *ledger.Ledger) error {
	// Validate input
	if ledger == nil {
		return fmt.Errorf("panicRecovery: ledger instance cannot be nil")
	}

	// Initiate recovery procedure
	if err := ledger.EnvironmentSystemCoreLedger.RecoveryManager.InitiateRecoveryProcedure(); err != nil {
		return fmt.Errorf("panicRecovery: failed to initiate recovery procedure: %w", err)
	}

	// Record the recovery event
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("RecoveryInitiated", nil); err != nil {
		return fmt.Errorf("panicRecovery: failed to record recovery event: %w", err)
	}

	// Log success
	log.Println("System recovery successfully initiated after panic.")
	return nil
}


// resumeOperation resumes system operations after a disruption using the SystemManager.
func resumeOperation(ledger *ledger.Ledger) error {
	// Validate input
	if ledger == nil {
		return fmt.Errorf("resumeOperation: ledger instance cannot be nil")
	}

	// Resume operations
	if err := ledger.EnvironmentSystemCoreLedger.SystemManager.ResumeSystemOperations(); err != nil {
		return fmt.Errorf("resumeOperation: failed to resume system operations: %w", err)
	}

	// Record the operation resumption
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("OperationsResumed", nil); err != nil {
		return fmt.Errorf("resumeOperation: failed to record resumption event: %w", err)
	}

	// Log success
	log.Println("System operations resumed successfully.")
	return nil
}


// softReset performs a soft reset of the system using the SystemManager.
func softReset(ledger *ledger.Ledger) error {
	// Validate input
	if ledger == nil {
		return fmt.Errorf("softReset: ledger instance cannot be nil")
	}

	// Perform soft reset
	if err := ledger.EnvironmentSystemCoreLedger.SystemManager.PerformSoftReset(); err != nil {
		return fmt.Errorf("softReset: failed to perform soft reset: %w", err)
	}

	// Record the soft reset
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("SoftResetPerformed", nil); err != nil {
		return fmt.Errorf("softReset: failed to record soft reset event: %w", err)
	}

	// Log success
	log.Println("Soft reset completed successfully.")
	return nil
}


// hardReset performs a hard reset of the system using the SystemManager.
func hardReset(ledger *ledger.Ledger) error {
	// Validate input
	if ledger == nil {
		return fmt.Errorf("hardReset: ledger instance cannot be nil")
	}

	// Perform hard reset
	if err := ledger.EnvironmentSystemCoreLedger.SystemManager.PerformHardReset(); err != nil {
		return fmt.Errorf("hardReset: failed to perform hard reset: %w", err)
	}

	// Record the hard reset
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("HardResetPerformed", nil); err != nil {
		return fmt.Errorf("hardReset: failed to record hard reset event: %w", err)
	}

	// Log success
	log.Println("Hard reset completed successfully.")
	return nil
}


// queryInterruptStatus retrieves the status of a specific interrupt using the InterruptManager.
func queryInterruptStatus(ledger *ledger.Ledger, interruptID string) (string, error) {
	// Validate input
	if ledger == nil {
		return "", fmt.Errorf("queryInterruptStatus: ledger instance cannot be nil")
	}
	if interruptID == "" {
		return "", fmt.Errorf("queryInterruptStatus: interruptID cannot be empty")
	}

	// Get interrupt status
	status, err := ledger.EnvironmentSystemCoreLedger.InterruptManager.GetInterruptStatus(interruptID)
	if err != nil {
		return "", fmt.Errorf("queryInterruptStatus: failed to retrieve interrupt status: %w", err)
	}

	return status, nil
}


// setInterruptPriority sets the priority of a specific interrupt using the InterruptManager.
func setInterruptPriority(ledger *ledger.Ledger, interruptID string, priority int) error {
	// Validate inputs
	if ledger == nil {
		return fmt.Errorf("setInterruptPriority: ledger instance cannot be nil")
	}
	if interruptID == "" {
		return fmt.Errorf("setInterruptPriority: interruptID cannot be empty")
	}
	if priority < 0 {
		return fmt.Errorf("setInterruptPriority: priority cannot be negative")
	}

	// Set interrupt priority
	if err := ledger.EnvironmentSystemCoreLedger.InterruptManager.SetInterruptPriority(interruptID, priority); err != nil {
		return fmt.Errorf("setInterruptPriority: failed to set interrupt priority: %w", err)
	}

	// Record the priority change
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("InterruptPrioritySet", map[string]interface{}{
		"InterruptID": interruptID,
		"Priority":    priority,
	}); err != nil {
		return fmt.Errorf("setInterruptPriority: failed to record interrupt priority change: %w", err)
	}

	return nil
}


// getInterruptPriority retrieves the priority of a specific interrupt using the InterruptManager.
func getInterruptPriority(ledger *ledger.Ledger, interruptID string) (int, error) {
	// Validate inputs
	if ledger == nil {
		return 0, fmt.Errorf("getInterruptPriority: ledger instance cannot be nil")
	}
	if interruptID == "" {
		return 0, fmt.Errorf("getInterruptPriority: interruptID cannot be empty")
	}

	// Get interrupt priority
	priority, err := ledger.EnvironmentSystemCoreLedger.InterruptManager.GetInterruptPriority(interruptID)
	if err != nil {
		return 0, fmt.Errorf("getInterruptPriority: failed to get interrupt priority: %w", err)
	}

	return priority, nil
}


// disableInterrupts disables all interrupts in the system using the InterruptManager.
func disableInterrupts(ledger *ledger.Ledger) error {
	// Validate input
	if ledger == nil {
		return fmt.Errorf("disableInterrupts: ledger instance cannot be nil")
	}

	// Disable all interrupts
	if err := ledger.EnvironmentSystemCoreLedger.InterruptManager.DisableAllInterrupts(); err != nil {
		return fmt.Errorf("disableInterrupts: failed to disable all interrupts: %w", err)
	}

	// Record the event
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("AllInterruptsDisabled", nil); err != nil {
		return fmt.Errorf("disableInterrupts: failed to record interrupts disable event: %w", err)
	}

	// Log success
	log.Println("All interrupts disabled successfully.")
	return nil
}

