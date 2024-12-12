package environment_and_system_core

import (
	"fmt"
	"log"
	"synnergy_network/pkg/ledger"
	"time"
)

// allocatePriorityResources reserves resources of a specific type for a context.
func allocatePriorityResources(ledgerInstance *ledger.Ledger, contextID, resourceType string, amount int) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("allocatePriorityResources: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("allocatePriorityResources: contextID cannot be empty")
	}
	if resourceType == "" {
		return fmt.Errorf("allocatePriorityResources: resourceType cannot be empty")
	}
	if amount <= 0 {
		return fmt.Errorf("allocatePriorityResources: amount must be greater than zero")
	}

	// Lock for thread-safe resource allocation
	contextLock.Lock()
	defer contextLock.Unlock()

	// Reserve resources in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ResourceManager.ReserveResources(contextID, resourceType, amount); err != nil {
		return fmt.Errorf("allocatePriorityResources: failed to reserve resources for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Resources of type %s successfully allocated to context %s (amount: %d).", resourceType, contextID, amount)
	return nil
}


// releasePriorityResources releases resources of a specific type for a context.
func releasePriorityResources(ledgerInstance *ledger.Ledger, contextID, resourceType string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("releasePriorityResources: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("releasePriorityResources: contextID cannot be empty")
	}
	if resourceType == "" {
		return fmt.Errorf("releasePriorityResources: resourceType cannot be empty")
	}

	// Lock for thread-safe resource release
	contextLock.Lock()
	defer contextLock.Unlock()

	// Release resources in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ResourceManager.ReleaseResources(contextID, resourceType); err != nil {
		return fmt.Errorf("releasePriorityResources: failed to release resources for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Resources of type %s successfully released for context %s.", resourceType, contextID)
	return nil
}


// synchronizeContextClock synchronizes the context's clock with the system time.
func synchronizeContextClock(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("synchronizeContextClock: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("synchronizeContextClock: contextID cannot be empty")
	}

	// Synchronize clock
	systemTime := time.Now().UTC()
	if err := ledgerInstance.EnvironmentSystemCoreLedger.TimeManager.SetContextClock(contextID, systemTime); err != nil {
		return fmt.Errorf("synchronizeContextClock: failed to synchronize clock for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Clock synchronized for context %s to system time %s.", contextID, systemTime)
	return nil
}


// resetContextClock resets the context's clock to its default state.
func resetContextClock(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("resetContextClock: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("resetContextClock: contextID cannot be empty")
	}

	// Reset the clock in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.TimeManager.ResetContextClock(contextID); err != nil {
		return fmt.Errorf("resetContextClock: failed to reset clock for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Clock reset for context %s.", contextID)
	return nil
}


// setContextMemoryLimit sets the memory limit for a specific context.
func setContextMemoryLimit(ledgerInstance *ledger.Ledger, contextID string, memoryLimit int) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("setContextMemoryLimit: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("setContextMemoryLimit: contextID cannot be empty")
	}
	if memoryLimit <= 0 {
		return fmt.Errorf("setContextMemoryLimit: memoryLimit must be greater than zero")
	}

	// Update memory limit in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.MemoryManager.UpdateMemoryLimit(contextID, memoryLimit); err != nil {
		return fmt.Errorf("setContextMemoryLimit: failed to set memory limit for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Memory limit for context %s set to %d.", contextID, memoryLimit)
	return nil
}


// adjustExecutionCapacity adjusts the execution capacity for a specific context.
func adjustExecutionCapacity(ledgerInstance *ledger.Ledger, contextID string, capacity int) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("adjustExecutionCapacity: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("adjustExecutionCapacity: contextID cannot be empty")
	}
	if capacity <= 0 {
		return fmt.Errorf("adjustExecutionCapacity: capacity must be greater than zero")
	}

	// Update execution capacity in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.UpdateExecutionCapacity(contextID, capacity); err != nil {
		return fmt.Errorf("adjustExecutionCapacity: failed to adjust execution capacity for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution capacity for context %s adjusted to %d.", contextID, capacity)
	return nil
}


// initializeContextVariables sets variables for a specific context.
func initializeContextVariables(ledgerInstance *ledger.Ledger, contextID string, variables map[string]interface{}) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("initializeContextVariables: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("initializeContextVariables: contextID cannot be empty")
	}
	if variables == nil || len(variables) == 0 {
		return fmt.Errorf("initializeContextVariables: variables cannot be nil or empty")
	}

	// Set context variables in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.SetContextVariables(contextID, variables); err != nil {
		return fmt.Errorf("initializeContextVariables: failed to set variables for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Context variables initialized for context %s.", contextID)
	return nil
}


// clearContextVariables removes all variables for a specific context.
func clearContextVariables(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("clearContextVariables: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("clearContextVariables: contextID cannot be empty")
	}

	// Clear variables in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ClearContextVariables(contextID); err != nil {
		return fmt.Errorf("clearContextVariables: failed to clear variables for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Context variables cleared for context %s.", contextID)
	return nil
}

// captureContextDiagnostics logs diagnostics information for a context.
func captureContextDiagnostics(ledgerInstance *ledger.Ledger, contextID, diagnostics string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("captureContextDiagnostics: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("captureContextDiagnostics: contextID cannot be empty")
	}
	if diagnostics == "" {
		return fmt.Errorf("captureContextDiagnostics: diagnostics cannot be empty")
	}

	// Create diagnostics entry
	entry := ledger.ContextDiagnosticsLog{
		ContextID:   contextID,
		Timestamp:   time.Now().UTC(),
		Diagnostics: diagnostics,
	}

	// Record diagnostics in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordContextDiagnostics(entry); err != nil {
		return fmt.Errorf("captureContextDiagnostics: failed to capture diagnostics for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Diagnostics captured for context %s.", contextID)
	return nil
}


// generateContextReport retrieves a detailed report for a context.
func generateContextReport(ledgerInstance *ledger.Ledger, contextID string) (ledger.ContextReport, error) {
	// Validate inputs
	if ledgerInstance == nil {
		return ledger.ContextReport{}, fmt.Errorf("generateContextReport: ledger instance cannot be nil")
	}
	if contextID == "" {
		return ledger.ContextReport{}, fmt.Errorf("generateContextReport: contextID cannot be empty")
	}

	// Retrieve report from the ledger
	report, err := ledgerInstance.EnvironmentSystemCoreLedger.GetContextReport(contextID)
	if err != nil {
		return ledger.ContextReport{}, fmt.Errorf("generateContextReport: failed to generate report for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Report generated for context %s.", contextID)
	return report, nil
}


// checkResourceAllocationPolicy verifies compliance with the resource allocation policy.
func checkResourceAllocationPolicy(ledgerInstance *ledger.Ledger, contextID string) (bool, error) {
	// Validate inputs
	if ledgerInstance == nil {
		return false, fmt.Errorf("checkResourceAllocationPolicy: ledger instance cannot be nil")
	}
	if contextID == "" {
		return false, fmt.Errorf("checkResourceAllocationPolicy: contextID cannot be empty")
	}

	// Verify policy compliance in the ledger
	compliant, err := ledgerInstance.EnvironmentSystemCoreLedger.VerifyResourcePolicy(contextID)
	if err != nil {
		return false, fmt.Errorf("checkResourceAllocationPolicy: failed to check policy compliance for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Resource allocation policy compliance checked for context %s: %t.", contextID, compliant)
	return compliant, nil
}


// adjustContextConcurrency updates the concurrency level for a specific context.
func adjustContextConcurrency(ledgerInstance *ledger.Ledger, contextID string, concurrencyLevel int) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("adjustContextConcurrency: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("adjustContextConcurrency: contextID cannot be empty")
	}
	if concurrencyLevel <= 0 {
		return fmt.Errorf("adjustContextConcurrency: concurrencyLevel must be greater than zero")
	}

	// Update concurrency level in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.UpdateContextConcurrency(contextID, concurrencyLevel); err != nil {
		return fmt.Errorf("adjustContextConcurrency: failed to adjust concurrency level for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Concurrency level for context %s adjusted to %d.", contextID, concurrencyLevel)
	return nil
}


// createContextCheckpoint creates a checkpoint for the specified context.
func createContextCheckpoint(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("createContextCheckpoint: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("createContextCheckpoint: contextID cannot be empty")
	}

	// Save checkpoint in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.SaveContextCheckpoint(contextID); err != nil {
		return fmt.Errorf("createContextCheckpoint: failed to create checkpoint for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Checkpoint created for context %s.", contextID)
	return nil
}


// restoreContextFromCheckpoint restores the context from its last checkpoint.
func restoreContextFromCheckpoint(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("restoreContextFromCheckpoint: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("restoreContextFromCheckpoint: contextID cannot be empty")
	}

	// Restore context from checkpoint in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.RestoreContextCheckpoint(contextID); err != nil {
		return fmt.Errorf("restoreContextFromCheckpoint: failed to restore context %s from checkpoint: %w", contextID, err)
	}

	// Log success
	log.Printf("Context %s restored from checkpoint.", contextID)
	return nil
}


// markContextForCleanup flags a context for cleanup.
func markContextForCleanup(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("markContextForCleanup: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("markContextForCleanup: contextID cannot be empty")
	}

	// Flag context for cleanup in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.FlagContextForCleanup(contextID); err != nil {
		return fmt.Errorf("markContextForCleanup: failed to flag context %s for cleanup: %w", contextID, err)
	}

	// Log success
	log.Printf("Context %s flagged for cleanup.", contextID)
	return nil
}


// logContextTermination logs the termination of a context.
func logContextTermination(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("logContextTermination: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("logContextTermination: contextID cannot be empty")
	}

	// Create termination log entry
	entry := ledger.ContextTerminationLog{
		ContextID: contextID,
		Timestamp: time.Now().UTC(),
	}

	// Record termination log in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordContextTermination(entry); err != nil {
		return fmt.Errorf("logContextTermination: failed to log termination for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Context termination logged for context %s.", contextID)
	return nil
}
