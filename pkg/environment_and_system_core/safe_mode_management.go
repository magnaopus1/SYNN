package environment_and_system_core

import (
	"fmt"
	"log"
	"synnergy_network/pkg/ledger"
)

// exitSafeMode transitions the system out of safe mode, ensuring all related states are updated and recorded in the ledger.
func exitSafeMode() error {
	// Ensure the system is currently in safe mode
	ledger := &ledger.Ledger{}
	if !ledger.EnvironmentSystemCoreLedger.SystemState.IsInSafeMode() {
		return fmt.Errorf("system is not in safe mode, cannot exit")
	}

	// Update system state to reflect the exit from safe mode
	if err := ledger.EnvironmentSystemCoreLedger.SystemState.ExitSafeMode(); err != nil {
		return fmt.Errorf("failed to exit safe mode: %w", err)
	}

	// Record the event in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("ExitedSafeMode", "SAFE_EXIT"); err != nil {
		return fmt.Errorf("failed to log safe mode exit in ledger: %w", err)
	}

	// Log the successful operation
	log.Println("System exited safe mode successfully.")
	return nil
}



// logInterruptEvent logs a blockchain interrupt event with a unique ID and relevant details.
func logInterruptEvent(eventID string, details string) error {
	// Validate input
	if eventID == "" {
		return fmt.Errorf("eventID cannot be empty")
	}
	if details == "" {
		return fmt.Errorf("details cannot be empty for eventID: %s", eventID)
	}

	// Add the interrupt event to the system state
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.InterruptManager.AddInterruptEvent(eventID, details); err != nil {
		return fmt.Errorf("failed to add interrupt event to system state: %w", err)
	}

	// Record the event in the ledger
	
	if err := ledger.EnvironmentSystemCoreLedger.RecordInterruptEvent(eventID, details); err != nil {
		return fmt.Errorf("failed to record interrupt event in ledger: %w", err)
	}

	// Log the successful operation
	log.Printf("Interrupt event logged successfully: %s - %s", eventID, details)
	return nil
}


// queryInterruptQueue retrieves the current interrupt queue from the blockchain system state.
func queryInterruptQueue() ([]string, error) {
	// Retrieve the interrupt queue from the system state
	ledger := &ledger.Ledger{}
	queue, err := ledger.EnvironmentSystemCoreLedger.InterruptManager.GetInterruptQueue()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve interrupt queue: %w", err)
	}

	// Log the size of the queue
	log.Printf("Interrupt queue retrieved with %d items.", len(queue))

	// Record the retrieval event in the ledger
	
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("InterruptQueueQueried", fmt.Sprintf("%d items retrieved", len(queue))); err != nil {
		return nil, fmt.Errorf("failed to log interrupt queue retrieval in ledger: %w", err)
	}

	return queue, nil
}

// clearInterruptQueue clears the interrupt queue from the system state and logs the operation.
func clearInterruptQueue() error {
	// Clear the interrupt queue from the system state
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.InterruptManager.ClearInterruptQueue(); err != nil {
		return fmt.Errorf("failed to clear interrupt queue: %w", err)
	}

	// Record the event in the ledger
	
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("InterruptQueueCleared", "Queue successfully cleared"); err != nil {
		return fmt.Errorf("failed to log interrupt queue clearance in ledger: %w", err)
	}

	// Log the successful operation
	log.Println("Interrupt queue cleared successfully.")
	return nil
}


// interruptStatusMonitor monitors the interrupt status of the blockchain system and logs the findings.
func interruptStatusMonitor() error {
	// Retrieve the current interrupt status
	ledger := &ledger.Ledger{}
	status, err := ledger.EnvironmentSystemCoreLedger.InterruptManager.GetInterruptStatus()
	if err != nil {
		return fmt.Errorf("failed to monitor interrupt status: %w", err)
	}

	// Record the interrupt status in the ledger

	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("InterruptStatusMonitored", fmt.Sprintf("Status: %v", status)); err != nil {
		return fmt.Errorf("failed to log interrupt status monitoring in ledger: %w", err)
	}

	// Log the monitored status
	log.Printf("Interrupt status monitored successfully: %v", status)
	return nil
}


// registerSystemInterrupt registers a new system interrupt with a unique ID and source details.
func registerSystemInterrupt(interruptID string, source string) error {
	// Validate input
	if interruptID == "" {
		return fmt.Errorf("interruptID cannot be empty")
	}
	if source == "" {
		return fmt.Errorf("source cannot be empty for interruptID: %s", interruptID)
	}

	// Add the interrupt to the system state
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.InterruptManager.RegisterInterrupt(interruptID, source); err != nil {
		return fmt.Errorf("failed to register system interrupt: %w", err)
	}

	// Record the interrupt in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordInterruptEvent(interruptID, source); err != nil {
		return fmt.Errorf("failed to record interrupt registration in ledger: %w", err)
	}

	// Log the successful operation
	log.Printf("System interrupt registered successfully: %s from %s", interruptID, source)
	return nil
}


// clearSystemInterrupt clears a specific interrupt by its ID.
func clearSystemInterrupt(interruptID string) error {
	// Validate input
	if interruptID == "" {
		return fmt.Errorf("interruptID cannot be empty")
	}

	// Remove the interrupt from the system state
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.InterruptManager.ClearInterrupt(interruptID); err != nil {
		return fmt.Errorf("failed to clear interrupt in system state: %w", err)
	}

	// Record the event in the ledger
	
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("InterruptCleared", interruptID); err != nil {
		return fmt.Errorf("failed to log interrupt clearance in ledger: %w", err)
	}

	// Log success
	log.Printf("System interrupt successfully cleared: %s", interruptID)
	return nil
}


// backupOnInterrupt triggers a backup operation during an interrupt.
func backupOnInterrupt() error {
	// Perform the backup operation
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.SystemBackupManager.TriggerBackup(); err != nil {
		return fmt.Errorf("failed to trigger system backup: %w", err)
	}

	// Record the backup event in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordBackupEvent("BACKUP_INT", "Backup triggered on interrupt"); err != nil {
		return fmt.Errorf("failed to log backup event in ledger: %w", err)
	}

	// Log success
	log.Println("System backup successfully triggered during interrupt.")
	return nil
}


// logSystemHaltReason logs the reason for a system halt.
func logSystemHaltReason(reason string) error {
	// Validate input
	if reason == "" {
		return fmt.Errorf("halt reason cannot be empty")
	}

	// Record the halt reason in the ledger
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent(reason, "SYS_HALT"); err != nil {
		return fmt.Errorf("failed to log system halt reason in ledger: %w", err)
	}

	// Log success
	log.Printf("System halt reason successfully logged: %s", reason)
	return nil
}


// defineEmergencyRecovery defines an emergency recovery protocol.
func defineEmergencyRecovery(protocolID string, recoverySteps []string) error {
	// Validate input
	if protocolID == "" {
		return fmt.Errorf("protocolID cannot be empty")
	}
	if len(recoverySteps) == 0 {
		return fmt.Errorf("recoverySteps cannot be empty for protocolID: %s", protocolID)
	}

	// Define the recovery protocol in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.RecoveryManager.DefineProtocol(protocolID, recoverySteps); err != nil {
		return fmt.Errorf("failed to define recovery protocol: %w", err)
	}

	// Record the protocol in the ledger
	
	if err := ledger.EnvironmentSystemCoreLedger.RecordRecoveryProtocol(protocolID, recoverySteps); err != nil {
		return fmt.Errorf("failed to log recovery protocol in ledger: %w", err)
	}

	// Log success
	log.Printf("Emergency recovery protocol defined successfully: %s", protocolID)
	return nil
}


// setInterruptMask sets an interrupt mask in the system.
func setInterruptMask(maskID string) error {
	// Validate input
	if maskID == "" {
		return fmt.Errorf("maskID cannot be empty")
	}

	// Apply the interrupt mask in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.InterruptManager.ApplyMask(maskID); err != nil {
		return fmt.Errorf("failed to apply interrupt mask: %w", err)
	}

	// Record the event in the ledger
	
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("InterruptMaskSet", maskID); err != nil {
		return fmt.Errorf("failed to log interrupt mask application in ledger: %w", err)
	}

	// Log success
	log.Printf("Interrupt mask successfully set: %s", maskID)
	return nil
}


// clearInterruptMask clears an interrupt mask from the system.
func clearInterruptMask(maskID string) error {
	// Validate input
	if maskID == "" {
		return fmt.Errorf("maskID cannot be empty")
	}

	// Remove the interrupt mask from the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.InterruptManager.ClearMask(maskID); err != nil {
		return fmt.Errorf("failed to clear interrupt mask: %w", err)
	}

	// Record the event in the ledger
	
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("InterruptMaskCleared", maskID); err != nil {
		return fmt.Errorf("failed to log interrupt mask clearance in ledger: %w", err)
	}

	// Log success
	log.Printf("Interrupt mask successfully cleared: %s", maskID)
	return nil
}


// setTrapHandler sets a trap handler for a specified trap ID and logs the event.
func setTrapHandler(trapID string, handlerFunc func() error) error {
	// Validate input
	if trapID == "" {
		return fmt.Errorf("trapID cannot be empty")
	}
	if handlerFunc == nil {
		return fmt.Errorf("handlerFunc cannot be nil for trapID: %s", trapID)
	}

	// Register the trap handler in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.TrapManager.RegisterHandler(trapID, handlerFunc); err != nil {
		return fmt.Errorf("failed to register trap handler for %s: %w", trapID, err)
	}

	// Record the event in the ledger
	
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("TrapHandlerSet", trapID); err != nil {
		return fmt.Errorf("failed to log trap handler registration in ledger: %w", err)
	}

	// Log success
	log.Printf("Trap handler set successfully for trap %s", trapID)
	return nil
}


// clearTrapHandler clears a trap handler for a specified trap ID and logs the event.
func clearTrapHandler(trapID string) error {
	// Validate input
	if trapID == "" {
		return fmt.Errorf("trapID cannot be empty")
	}

	// Remove the trap handler from the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.TrapManager.RemoveHandler(trapID); err != nil {
		return fmt.Errorf("failed to clear trap handler for %s: %w", trapID, err)
	}

	// Record the event in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("TrapHandlerCleared", trapID); err != nil {
		return fmt.Errorf("failed to log trap handler removal in ledger: %w", err)
	}

	// Log success
	log.Printf("Trap handler cleared successfully for trap %s", trapID)
	return nil
}


// initiateFallbackMechanism initiates a fallback mechanism for the specified ID and logs the event.
func initiateFallbackMechanism(mechanismID string) error {
	// Validate input
	if mechanismID == "" {
		return fmt.Errorf("mechanismID cannot be empty")
	}

	// Trigger the fallback mechanism in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.FallbackManager.InitiateMechanism(mechanismID); err != nil {
		return fmt.Errorf("failed to initiate fallback mechanism %s: %w", mechanismID, err)
	}

	// Record the event in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("FallbackMechanismInitiated", mechanismID); err != nil {
		return fmt.Errorf("failed to log fallback mechanism activation in ledger: %w", err)
	}

	// Log success
	log.Printf("Fallback mechanism successfully initiated: %s", mechanismID)
	return nil
}


// logBackupEvent logs a backup event with a specified ID and description in the ledger.
func logBackupEvent(eventID string, description string) error {
	// Validate input
	if eventID == "" {
		return fmt.Errorf("eventID cannot be empty")
	}
	if description == "" {
		return fmt.Errorf("description cannot be empty for eventID: %s", eventID)
	}

	// Record the backup event in the ledger
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.RecordBackupEvent(eventID, description); err != nil {
		return fmt.Errorf("failed to log backup event in ledger: %w", err)
	}

	// Log success
	log.Printf("Backup event successfully logged: %s - %s", eventID, description)
	return nil
}


// monitorInterruptStatus monitors the current interrupt status and logs the findings in the ledger.
func monitorInterruptStatus() error {
	// Retrieve the current interrupt status
	ledger := &ledger.Ledger{}
	status, err := ledger.EnvironmentSystemCoreLedger.InterruptManager.GetStatus()
	if err != nil {
		return fmt.Errorf("failed to retrieve interrupt status: %w", err)
	}

	// Record the status in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordSystemEvent("InterruptStatusMonitored", fmt.Sprintf("Status: %v", status)); err != nil {
		return fmt.Errorf("failed to log interrupt status monitoring in ledger: %w", err)
	}

	// Log success
	log.Printf("Interrupt status successfully monitored: %v", status)
	return nil
}


// setupRecoveryProtocol sets up a recovery protocol with specified steps and logs the event.
func setupRecoveryProtocol(protocolID string, recoverySteps []string) error {
	// Validate input
	if protocolID == "" {
		return fmt.Errorf("protocolID cannot be empty")
	}
	if len(recoverySteps) == 0 {
		return fmt.Errorf("recoverySteps cannot be empty for protocolID: %s", protocolID)
	}

	// Define the recovery protocol in the system
	ledger := &ledger.Ledger{}
	if err := ledger.EnvironmentSystemCoreLedger.RecoveryManager.DefineProtocol(protocolID, recoverySteps); err != nil {
		return fmt.Errorf("failed to define recovery protocol %s: %w", protocolID, err)
	}

	// Record the protocol in the ledger
	
	if err := ledger.EnvironmentSystemCoreLedger.RecordRecoveryProtocol(protocolID, recoverySteps); err != nil {
		return fmt.Errorf("failed to log recovery protocol in ledger: %w", err)
	}

	// Log success
	log.Printf("Recovery protocol successfully set: %s", protocolID)
	return nil
}
