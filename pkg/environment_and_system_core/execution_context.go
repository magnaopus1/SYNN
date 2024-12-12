package environment_and_system_core

import (
	"fmt"
	"log"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// Mutex for concurrent context operations
var contextLock sync.Mutex

// createExecutionContext creates a new execution context with specified resources and priority.
func createExecutionContext(ledgerInstance *ledger.Ledger, contextID string, resources ledger.ResourceRequirements, priority int) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("createExecutionContext: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("createExecutionContext: contextID cannot be empty")
	}
	if priority < 0 {
		return fmt.Errorf("createExecutionContext: priority cannot be negative")
	}

	// Lock to ensure concurrency safety
	contextLock.Lock()
	defer contextLock.Unlock()

	// Encrypt the resource requirements
	encryptedResources, err := common.EncryptData(resources)
	if err != nil {
		return fmt.Errorf("createExecutionContext: failed to encrypt resources: %w", err)
	}

	// Record the execution context in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.RecordExecutionContext(contextID, encryptedResources, priority); err != nil {
		return fmt.Errorf("createExecutionContext: failed to create execution context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution context %s created successfully with priority %d.", contextID, priority)
	return nil
}


// switchExecutionContext switches the active execution context to a new one.
func switchExecutionContext(ledgerInstance *ledger.Ledger, newContextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("switchExecutionContext: ledger instance cannot be nil")
	}
	if newContextID == "" {
		return fmt.Errorf("switchExecutionContext: newContextID cannot be empty")
	}

	// Lock to ensure concurrency safety
	contextLock.Lock()
	defer contextLock.Unlock()

	// Set the active context in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.SetActiveContext(newContextID); err != nil {
		return fmt.Errorf("switchExecutionContext: failed to switch to context %s: %w", newContextID, err)
	}

	// Log success
	log.Printf("Switched active execution context to %s.", newContextID)
	return nil
}


// pauseExecutionContext pauses the execution of a given context.
func pauseExecutionContext(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("pauseExecutionContext: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("pauseExecutionContext: contextID cannot be empty")
	}

	// Lock to ensure concurrency safety
	contextLock.Lock()
	defer contextLock.Unlock()

	// Update the context status to "Paused"
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.UpdateContextStatus(contextID, "Paused"); err != nil {
		return fmt.Errorf("pauseExecutionContext: failed to pause execution context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution context %s paused successfully.", contextID)
	return nil
}


// resumeExecutionContext resumes a paused execution context.
func resumeExecutionContext(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("resumeExecutionContext: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("resumeExecutionContext: contextID cannot be empty")
	}

	// Lock to ensure concurrency safety
	contextLock.Lock()
	defer contextLock.Unlock()

	// Update the context status to "Active"
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.UpdateContextStatus(contextID, "Active"); err != nil {
		return fmt.Errorf("resumeExecutionContext: failed to resume execution context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution context %s resumed successfully.", contextID)
	return nil
}


// destroyExecutionContext deletes a given execution context.
func destroyExecutionContext(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("destroyExecutionContext: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("destroyExecutionContext: contextID cannot be empty")
	}

	// Lock to ensure concurrency safety
	contextLock.Lock()
	defer contextLock.Unlock()

	// Delete the execution context
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.DeleteExecutionContext(contextID); err != nil {
		return fmt.Errorf("destroyExecutionContext: failed to destroy execution context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution context %s destroyed successfully.", contextID)
	return nil
}


// setExecutionPriority sets the priority for a given execution context.
func setExecutionPriority(ledgerInstance *ledger.Ledger, contextID string, priority int) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("setExecutionPriority: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("setExecutionPriority: contextID cannot be empty")
	}
	if priority < 0 {
		return fmt.Errorf("setExecutionPriority: priority cannot be negative")
	}

	// Update the priority in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.UpdateContextPriority(contextID, priority); err != nil {
		return fmt.Errorf("setExecutionPriority: failed to set priority for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution priority for context %s set to %d.", contextID, priority)
	return nil
}


// queryExecutionStatus retrieves the current status of an execution context.
func queryExecutionStatus(ledgerInstance *ledger.Ledger, contextID string) (ledger.ExecutionStatus, error) {
	// Validate inputs
	if ledgerInstance == nil {
		return ledger.ExecutionStatus{}, fmt.Errorf("queryExecutionStatus: ledger instance cannot be nil")
	}
	if contextID == "" {
		return ledger.ExecutionStatus{}, fmt.Errorf("queryExecutionStatus: contextID cannot be empty")
	}

	// Get the context status from the ledger
	status, err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.GetContextStatus(contextID)
	if err != nil {
		return ledger.ExecutionStatus{}, fmt.Errorf("queryExecutionStatus: failed to query status for context %s: %w", contextID, err)
	}

	// Return the status
	return status, nil
}


// adjustResourceAllocation adjusts the resources allocated to a given execution context.
func adjustResourceAllocation(ledgerInstance *ledger.Ledger, contextID string, newResources ledger.ResourceRequirements) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("adjustResourceAllocation: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("adjustResourceAllocation: contextID cannot be empty")
	}

	// Encrypt the new resources
	encryptedResources, err := common.EncryptData(newResources)
	if err != nil {
		return fmt.Errorf("adjustResourceAllocation: failed to encrypt resources: %w", err)
	}

	// Update the resource allocation in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.UpdateContextResources(contextID, encryptedResources); err != nil {
		return fmt.Errorf("adjustResourceAllocation: failed to adjust resources for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Resource allocation adjusted for context %s.", contextID)
	return nil
}


// setTimeoutForExecution sets a timeout for the given execution context.
func setTimeoutForExecution(ledgerInstance *ledger.Ledger, contextID string, timeout time.Duration) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("setTimeoutForExecution: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("setTimeoutForExecution: contextID cannot be empty")
	}
	if timeout <= 0 {
		return fmt.Errorf("setTimeoutForExecution: timeout must be greater than 0")
	}

	// Set the timeout in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.SetContextTimeout(contextID, timeout); err != nil {
		return fmt.Errorf("setTimeoutForExecution: failed to set timeout for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Timeout of %v set for context %s.", timeout, contextID)
	return nil
}

// lockExecutionContext locks an execution context to prevent modifications.
func lockExecutionContext(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("lockExecutionContext: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("lockExecutionContext: contextID cannot be empty")
	}

	// Lock to ensure concurrency safety
	contextLock.Lock()
	defer contextLock.Unlock()

	// Lock the context in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.LockContext(contextID); err != nil {
		return fmt.Errorf("lockExecutionContext: failed to lock context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution context %s locked successfully.", contextID)
	return nil
}


// unlockExecutionContext unlocks an execution context to allow modifications.
func unlockExecutionContext(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("unlockExecutionContext: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("unlockExecutionContext: contextID cannot be empty")
	}

	// Lock for concurrency safety
	contextLock.Lock()
	defer contextLock.Unlock()

	// Unlock the context in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.UnlockContext(contextID); err != nil {
		return fmt.Errorf("unlockExecutionContext: failed to unlock context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution context %s unlocked successfully.", contextID)
	return nil
}

// saveExecutionState saves the current state of an execution context securely.
func saveExecutionState(ledgerInstance *ledger.Ledger, contextID string, state ledger.ExecutionState) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("saveExecutionState: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("saveExecutionState: contextID cannot be empty")
	}

	// Encrypt the execution state
	encryptedState, err := common.EncryptData(state)
	if err != nil {
		return fmt.Errorf("saveExecutionState: failed to encrypt execution state for context %s: %w", contextID, err)
	}

	// Save the encrypted state in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.SaveContextState(contextID, encryptedState); err != nil {
		return fmt.Errorf("saveExecutionState: failed to save execution state for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution state for context %s saved successfully.", contextID)
	return nil
}


// loadExecutionState retrieves and decrypts the execution state of a context.
func loadExecutionState(ledgerInstance *ledger.Ledger, contextID string) (ledger.ExecutionState, error) {
	// Validate inputs
	if ledgerInstance == nil {
		return ledger.ExecutionState{}, fmt.Errorf("loadExecutionState: ledger instance cannot be nil")
	}
	if contextID == "" {
		return ledger.ExecutionState{}, fmt.Errorf("loadExecutionState: contextID cannot be empty")
	}

	// Retrieve the encrypted state from the ledger
	encryptedState, err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.RetrieveContextState(contextID)
	if err != nil {
		return ledger.ExecutionState{}, fmt.Errorf("loadExecutionState: failed to retrieve execution state for context %s: %w", contextID, err)
	}

	// Decrypt the state
	var state ledger.ExecutionState
	if err := common.DecryptData(encryptedState, &state); err != nil {
		return ledger.ExecutionState{}, fmt.Errorf("loadExecutionState: failed to decrypt execution state for context %s: %w", contextID, err)
	}

	// Return the state
	return state, nil
}


// isolateExecutionContext isolates a context to prevent interference from other processes.
func isolateExecutionContext(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("isolateExecutionContext: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("isolateExecutionContext: contextID cannot be empty")
	}

	// Set context isolation in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.SetContextIsolation(contextID, true); err != nil {
		return fmt.Errorf("isolateExecutionContext: failed to isolate context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution context %s isolated successfully.", contextID)
	return nil
}


// mergeExecutionContext merges two execution contexts into a new context.
func mergeExecutionContext(ledgerInstance *ledger.Ledger, contextID1, contextID2, newContextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("mergeExecutionContext: ledger instance cannot be nil")
	}
	if contextID1 == "" || contextID2 == "" || newContextID == "" {
		return fmt.Errorf("mergeExecutionContext: all context IDs must be provided")
	}

	// Load the states of the two contexts
	state1, err := loadExecutionState(ledgerInstance, contextID1)
	if err != nil {
		return fmt.Errorf("mergeExecutionContext: failed to load execution state for context %s: %w", contextID1, err)
	}

	state2, err := loadExecutionState(ledgerInstance, contextID2)
	if err != nil {
		return fmt.Errorf("mergeExecutionContext: failed to load execution state for context %s: %w", contextID2, err)
	}

	// Merge the states
	mergedState := ledger.ExecutionState{
		ContextID: newContextID,
		Resources: state1.Resources + ";" + state2.Resources,
		Status:    "Merged",
		Timestamp: time.Now().UTC(),
	}

	// Save the merged state
	if err := saveExecutionState(ledgerInstance, newContextID, mergedState); err != nil {
		return fmt.Errorf("mergeExecutionContext: failed to save merged state for context %s: %w", newContextID, err)
	}

	// Log success
	log.Printf("Contexts %s and %s merged into new context %s successfully.", contextID1, contextID2, newContextID)
	return nil
}


// suspendExecutionContext suspends an execution context temporarily.
func suspendExecutionContext(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("suspendExecutionContext: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("suspendExecutionContext: contextID cannot be empty")
	}

	// Update the context status to "Suspended"
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.UpdateContextStatus(contextID, "Suspended"); err != nil {
		return fmt.Errorf("suspendExecutionContext: failed to suspend context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution context %s suspended successfully.", contextID)
	return nil
}

