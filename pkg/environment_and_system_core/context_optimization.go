package environment_and_system_core

import (
	"fmt"
	"log"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// Mutex to manage concurrent access to context optimization operations
var contextOptLock sync.Mutex

// setExecutionRestartPolicy sets the restart policy for a specific execution context.
func setExecutionRestartPolicy(ledgerInstance *ledger.Ledger, contextID string, policy ledger.RestartPolicy) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("setExecutionRestartPolicy: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("setExecutionRestartPolicy: contextID cannot be empty")
	}

	// Update the restart policy in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.PolicyManager.UpdateRestartPolicy(contextID, policy); err != nil {
		return fmt.Errorf("setExecutionRestartPolicy: failed to update restart policy for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Restart policy updated for context %s successfully.", contextID)
	return nil
}



// validateContextSwitch checks if a switch between two contexts is permissible.
func validateContextSwitch(ledgerInstance *ledger.Ledger, fromContextID, toContextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("validateContextSwitch: ledger instance cannot be nil")
	}
	if fromContextID == "" || toContextID == "" {
		return fmt.Errorf("validateContextSwitch: context IDs cannot be empty")
	}

	// Check eligibility for context switching
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.CheckContextSwitchEligibility(fromContextID, toContextID); err != nil {
		return fmt.Errorf("validateContextSwitch: context switch validation failed from %s to %s: %w", fromContextID, toContextID, err)
	}

	// Log success
	log.Printf("Context switch validated from %s to %s.", fromContextID, toContextID)
	return nil
}



// delegateExecutionTask assigns a task from one context to another.
func delegateExecutionTask(ledgerInstance *ledger.Ledger, fromContextID, toContextID, taskID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("delegateExecutionTask: ledger instance cannot be nil")
	}
	if fromContextID == "" || toContextID == "" || taskID == "" {
		return fmt.Errorf("delegateExecutionTask: context IDs and taskID cannot be empty")
	}

	// Assign the task to the new context
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.AssignTaskToContext(taskID, toContextID); err != nil {
		return fmt.Errorf("delegateExecutionTask: failed to delegate task %s to context %s: %w", taskID, toContextID, err)
	}

	// Log success
	log.Printf("Task %s delegated from context %s to %s successfully.", taskID, fromContextID, toContextID)
	return nil
}


// retractExecutionDelegation reassigns a task back to its original context.
func retractExecutionDelegation(ledgerInstance *ledger.Ledger, taskID, originalContextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("retractExecutionDelegation: ledger instance cannot be nil")
	}
	if taskID == "" || originalContextID == "" {
		return fmt.Errorf("retractExecutionDelegation: taskID and originalContextID cannot be empty")
	}

	// Reassign the task to the original context
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.ReassignTaskToContext(taskID, originalContextID); err != nil {
		return fmt.Errorf("retractExecutionDelegation: failed to reassign task %s to context %s: %w", taskID, originalContextID, err)
	}

	// Log success
	log.Printf("Task %s retracted to original context %s successfully.", taskID, originalContextID)
	return nil
}


// registerContextObserver registers an observer for a specific context.
func registerContextObserver(ledgerInstance *ledger.Ledger, contextID, observerID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("registerContextObserver: ledger instance cannot be nil")
	}
	if contextID == "" || observerID == "" {
		return fmt.Errorf("registerContextObserver: contextID and observerID cannot be empty")
	}

	// Add the observer to the context
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ObserverManager.AddContextObserver(contextID, observerID); err != nil {
		return fmt.Errorf("registerContextObserver: failed to register observer %s for context %s: %w", observerID, contextID, err)
	}

	// Log success
	log.Printf("Observer %s registered for context %s successfully.", observerID, contextID)
	return nil
}


// unregisterContextObserver removes an observer from a specific context.
func unregisterContextObserver(ledgerInstance *ledger.Ledger, contextID, observerID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("unregisterContextObserver: ledger instance cannot be nil")
	}
	if contextID == "" || observerID == "" {
		return fmt.Errorf("unregisterContextObserver: contextID and observerID cannot be empty")
	}

	// Remove the observer from the context
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ObserverManager.RemoveContextObserver(contextID, observerID); err != nil {
		return fmt.Errorf("unregisterContextObserver: failed to unregister observer %s from context %s: %w", observerID, contextID, err)
	}

	// Log success
	log.Printf("Observer %s unregistered from context %s successfully.", observerID, contextID)
	return nil
}


// trackContextDependency adds a dependency between two contexts.
func trackContextDependency(ledgerInstance *ledger.Ledger, contextID, dependentContextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("trackContextDependency: ledger instance cannot be nil")
	}
	if contextID == "" || dependentContextID == "" {
		return fmt.Errorf("trackContextDependency: contextID and dependentContextID cannot be empty")
	}

	// Add the dependency in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.DependencyManager.AddContextDependency(contextID, dependentContextID); err != nil {
		return fmt.Errorf("trackContextDependency: failed to track dependency from context %s to %s: %w", contextID, dependentContextID, err)
	}

	// Log success
	log.Printf("Dependency tracked from context %s to %s successfully.", contextID, dependentContextID)
	return nil
}


// cleanupExecutionResiduals clears residuals left by a context's execution.
func cleanupExecutionResiduals(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("cleanupExecutionResiduals: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("cleanupExecutionResiduals: contextID cannot be empty")
	}

	// Clear residuals in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.CleanupManager.ClearResiduals(contextID); err != nil {
		return fmt.Errorf("cleanupExecutionResiduals: failed to clean up residuals for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution residuals cleaned up for context %s successfully.", contextID)
	return nil
}

// extendContextLifespan extends the lifespan of a context.
func extendContextLifespan(ledgerInstance *ledger.Ledger, contextID string, additionalTime time.Duration) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("extendContextLifespan: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("extendContextLifespan: contextID cannot be empty")
	}

	// Extend the lifespan in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ContextManager.ExtendContextDuration(contextID, additionalTime); err != nil {
		return fmt.Errorf("extendContextLifespan: failed to extend lifespan for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Context %s lifespan extended by %v successfully.", contextID, additionalTime)
	return nil
}

// reduceContextLifespan decreases the lifespan of a specific context by a given duration.
func reduceContextLifespan(ledgerInstance *ledger.Ledger, contextID string, reductionTime time.Duration) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("reduceContextLifespan: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("reduceContextLifespan: contextID cannot be empty")
	}
	if reductionTime <= 0 {
		return fmt.Errorf("reduceContextLifespan: reductionTime must be greater than zero")
	}

	// Reduce the context duration in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ContextManager.ReduceContextDuration(contextID, reductionTime); err != nil {
		return fmt.Errorf("reduceContextLifespan: failed to reduce lifespan for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Lifespan of context %s reduced by %v successfully.", contextID, reductionTime)
	return nil
}


// restrictContextAccess applies restrictions to a specific context's access.
func restrictContextAccess(ledgerInstance *ledger.Ledger, contextID string, restrictions ledger.AccessRestrictions) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("restrictContextAccess: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("restrictContextAccess: contextID cannot be empty")
	}

	// Apply the access restrictions in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.AccessManager.SetContextAccessRestrictions(contextID, restrictions); err != nil {
		return fmt.Errorf("restrictContextAccess: failed to set access restrictions for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Access restrictions applied to context %s successfully.", contextID)
	return nil
}


// grantContextAccess grants access rights to a specific context.
func grantContextAccess(ledgerInstance *ledger.Ledger, contextID string, accessRights ledger.AccessRights) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("grantContextAccess: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("grantContextAccess: contextID cannot be empty")
	}

	// Update the access rights in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.AccessManager.UpdateContextAccessRights(contextID, accessRights); err != nil {
		return fmt.Errorf("grantContextAccess: failed to update access rights for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Access rights granted to context %s successfully.", contextID)
	return nil
}

