package environment_and_system_core

import (
	"fmt"
	"log"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// Mutex for concurrent management
var execMgmtLock sync.Mutex

// restartExecutionContext restarts a suspended execution context.
func restartExecutionContext(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("restartExecutionContext: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("restartExecutionContext: contextID cannot be empty")
	}

	// Lock to ensure concurrency safety
	execMgmtLock.Lock()
	defer execMgmtLock.Unlock()

	// Check the status of the execution context
	status, err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.GetContextStatus(contextID)
	if err != nil {
		return fmt.Errorf("restartExecutionContext: failed to retrieve context status: %w", err)
	}
	if status != "Suspended" {
		return fmt.Errorf("restartExecutionContext: context %s must be suspended to restart", contextID)
	}

	// Restart the execution context
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.UpdateContextStatus(contextID, "Active"); err != nil {
		return fmt.Errorf("restartExecutionContext: failed to restart execution context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution context %s successfully restarted.", contextID)
	return nil
}


// forceTerminateContext forcibly terminates an execution context.
func forceTerminateContext(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("forceTerminateContext: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("forceTerminateContext: contextID cannot be empty")
	}

	// Lock to ensure concurrency safety
	execMgmtLock.Lock()
	defer execMgmtLock.Unlock()

	// Delete the execution context
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.DeleteExecutionContext(contextID); err != nil {
		return fmt.Errorf("forceTerminateContext: failed to terminate execution context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution context %s forcibly terminated.", contextID)
	return nil
}


// setExecutionQuota sets an execution quota for a context.
func setExecutionQuota(ledgerInstance *ledger.Ledger, contextID string, quota ledger.ExecutionConstraints) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("setExecutionQuota: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("setExecutionQuota: contextID cannot be empty")
	}

	// Encrypt the quota for security
	encryptedQuota, err := common.EncryptData(quota)
	if err != nil {
		return fmt.Errorf("setExecutionQuota: failed to encrypt quota for context %s: %w", contextID, err)
	}

	// Update the execution quota in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.QuotaManager.UpdateExecutionQuota(contextID, encryptedQuota); err != nil {
		return fmt.Errorf("setExecutionQuota: failed to set execution quota for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution quota successfully set for context %s.", contextID)
	return nil
}


// checkExecutionConstraints checks if execution constraints are violated for a context.
func checkExecutionConstraints(ledgerInstance *ledger.Ledger, contextID string) (bool, error) {
	// Validate inputs
	if ledgerInstance == nil {
		return false, fmt.Errorf("checkExecutionConstraints: ledger instance cannot be nil")
	}
	if contextID == "" {
		return false, fmt.Errorf("checkExecutionConstraints: contextID cannot be empty")
	}

	// Retrieve execution constraints
	constraints, err := ledgerInstance.EnvironmentSystemCoreLedger.QuotaManager.GetContextConstraints(contextID)
	if err != nil {
		return false, fmt.Errorf("checkExecutionConstraints: failed to retrieve execution constraints for context %s: %w", contextID, err)
	}

	// Check if constraints are violated
	if constraints.Usage > constraints.Quota {
		return false, fmt.Errorf("checkExecutionConstraints: execution quota exceeded for context %s", contextID)
	}

	return true, nil
}


// allocateDynamicMemory allocates dynamic memory for an execution context.
func allocateDynamicMemory(ledgerInstance *ledger.Ledger, contextID string, amount int) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("allocateDynamicMemory: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("allocateDynamicMemory: contextID cannot be empty")
	}
	if amount <= 0 {
		return fmt.Errorf("allocateDynamicMemory: amount must be greater than 0")
	}

	// Adjust memory allocation
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ResourceManager.AdjustMemoryAllocation(contextID, amount); err != nil {
		return fmt.Errorf("allocateDynamicMemory: failed to allocate memory for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Allocated %d units of memory for context %s.", amount, contextID)
	return nil
}


// freeDynamicMemory releases dynamic memory from an execution context.
func freeDynamicMemory(ledgerInstance *ledger.Ledger, contextID string, amount int) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("freeDynamicMemory: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("freeDynamicMemory: contextID cannot be empty")
	}
	if amount <= 0 {
		return fmt.Errorf("freeDynamicMemory: amount must be greater than 0")
	}

	// Adjust memory allocation to release memory
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ResourceManager.AdjustMemoryAllocation(contextID, -amount); err != nil {
		return fmt.Errorf("freeDynamicMemory: failed to free memory for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Freed %d units of memory from context %s.", amount, contextID)
	return nil
}


// reclaimMemoryResources reclaims all memory resources for a context.
func reclaimMemoryResources(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("reclaimMemoryResources: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("reclaimMemoryResources: contextID cannot be empty")
	}

	// Reset memory allocation
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ResourceManager.ResetMemoryAllocation(contextID); err != nil {
		return fmt.Errorf("reclaimMemoryResources: failed to reclaim memory resources for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Memory resources reclaimed for context %s.", contextID)
	return nil
}


// logExecutionActivity logs an execution activity in the ledger.
func logExecutionActivity(ledgerInstance *ledger.Ledger, contextID, activity string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("logExecutionActivity: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("logExecutionActivity: contextID cannot be empty")
	}
	if activity == "" {
		return fmt.Errorf("logExecutionActivity: activity cannot be empty")
	}

	// Create the execution log entry
	logEntry := ledger.ExecutionLogEntry{
		ContextID: contextID,
		Activity:  activity,
		Timestamp: time.Now().UTC(),
	}

	// Record the activity in the ledger
	if err := ledgerInstance.EnvironmentSystemCoreLedger.LogManager.RecordExecutionActivity(logEntry); err != nil {
		return fmt.Errorf("logExecutionActivity: failed to log execution activity for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution activity logged for context %s: %s.", contextID, activity)
	return nil
}


// monitorExecutionPerformance retrieves performance metrics for an execution context.
func monitorExecutionPerformance(ledgerInstance *ledger.Ledger, contextID string) (ledger.PerformanceMetrics, error) {
	// Validate inputs
	if ledgerInstance == nil {
		return ledger.PerformanceMetrics{}, fmt.Errorf("monitorExecutionPerformance: ledger instance cannot be nil")
	}
	if contextID == "" {
		return ledger.PerformanceMetrics{}, fmt.Errorf("monitorExecutionPerformance: contextID cannot be empty")
	}

	// Retrieve performance metrics
	performanceData, err := ledgerInstance.EnvironmentSystemCoreLedger.DiagnosticManager.GetPerformanceMetrics(contextID)
	if err != nil {
		return ledger.PerformanceMetrics{}, fmt.Errorf("monitorExecutionPerformance: failed to monitor performance for context %s: %w", contextID, err)
	}

	return performanceData, nil
}


// throttleExecutionRate sets a rate limit for execution within a specific context.
func throttleExecutionRate(ledgerInstance *ledger.Ledger, contextID string, rateLimit int) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("throttleExecutionRate: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("throttleExecutionRate: contextID cannot be empty")
	}
	if rateLimit <= 0 {
		return fmt.Errorf("throttleExecutionRate: rateLimit must be greater than 0")
	}

	// Set the rate limit
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.SetRateLimit(contextID, rateLimit); err != nil {
		return fmt.Errorf("throttleExecutionRate: failed to throttle execution rate for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution rate for context %s throttled to %d requests per second.", contextID, rateLimit)
	return nil
}


// setContextRecoveryPoint sets a recovery point for the given context with encrypted state data.
func setContextRecoveryPoint(ledgerInstance *ledger.Ledger, contextID string, recoveryData ledger.ExecutionState) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("setContextRecoveryPoint: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("setContextRecoveryPoint: contextID cannot be empty")
	}

	// Encrypt recovery data
	encryptedData, err := common.EncryptData(recoveryData)
	if err != nil {
		return fmt.Errorf("setContextRecoveryPoint: failed to encrypt recovery data for context %s: %w", contextID, err)
	}

	// Save recovery point
	if err := ledgerInstance.EnvironmentSystemCoreLedger.RecoveryManager.SaveRecoveryPoint(contextID, encryptedData); err != nil {
		return fmt.Errorf("setContextRecoveryPoint: failed to set recovery point for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Recovery point set successfully for context %s.", contextID)
	return nil
}



// rollbackToRecoveryPoint rolls back a context to its most recent recovery point.
func rollbackToRecoveryPoint(ledgerInstance *ledger.Ledger, contextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("rollbackToRecoveryPoint: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("rollbackToRecoveryPoint: contextID cannot be empty")
	}

	// Retrieve encrypted recovery state
	encryptedState, err := ledgerInstance.EnvironmentSystemCoreLedger.RecoveryManager.GetRecoveryPoint(contextID)
	if err != nil {
		return fmt.Errorf("rollbackToRecoveryPoint: failed to retrieve recovery point for context %s: %w", contextID, err)
	}

	// Decrypt the recovery state
	recoveryState, err := common.DecryptData(encryptedState)
	if err != nil {
		return fmt.Errorf("rollbackToRecoveryPoint: failed to decrypt recovery data for context %s: %w", contextID, err)
	}

	// Restore the execution state
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.RestoreExecutionState(contextID, recoveryState); err != nil {
		return fmt.Errorf("rollbackToRecoveryPoint: failed to rollback context %s to recovery point: %w", contextID, err)
	}

	// Log success
	log.Printf("Context %s rolled back to recovery point successfully.", contextID)
	return nil
}


// validateExecutionEnvironment validates the configuration and security of an execution environment.
func validateExecutionEnvironment(ledgerInstance *ledger.Ledger, contextID string) (bool, error) {
	// Validate inputs
	if ledgerInstance == nil {
		return false, fmt.Errorf("validateExecutionEnvironment: ledger instance cannot be nil")
	}
	if contextID == "" {
		return false, fmt.Errorf("validateExecutionEnvironment: contextID cannot be empty")
	}

	// Check environment status
	envStatus, err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.CheckEnvironmentStatus(contextID)
	if err != nil {
		return false, fmt.Errorf("validateExecutionEnvironment: failed to validate environment for context %s: %w", contextID, err)
	}

	// Return status
	return envStatus.IsConfigured && envStatus.IsSecure, nil
}


// updateExecutionEnvironment updates the configuration of an execution environment.
func updateExecutionEnvironment(ledgerInstance *ledger.Ledger, contextID string, newConfig ledger.EnvironmentConfig) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("updateExecutionEnvironment: ledger instance cannot be nil")
	}
	if contextID == "" {
		return fmt.Errorf("updateExecutionEnvironment: contextID cannot be empty")
	}

	// Encrypt the new environment configuration
	encryptedConfig, err := common.EncryptData(newConfig)
	if err != nil {
		return fmt.Errorf("updateExecutionEnvironment: failed to encrypt environment configuration for context %s: %w", contextID, err)
	}

	// Update environment configuration
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.UpdateEnvironmentConfig(contextID, encryptedConfig); err != nil {
		return fmt.Errorf("updateExecutionEnvironment: failed to update environment for context %s: %w", contextID, err)
	}

	// Log success
	log.Printf("Execution environment updated successfully for context %s.", contextID)
	return nil
}

// createSubExecutionContext creates a sub-execution context under a parent context with specified resources.
func createSubExecutionContext(ledgerInstance *ledger.Ledger, parentContextID, subContextID string, resources ledger.ResourceRequirements) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("createSubExecutionContext: ledger instance cannot be nil")
	}
	if parentContextID == "" || subContextID == "" {
		return fmt.Errorf("createSubExecutionContext: parentContextID and subContextID cannot be empty")
	}

	// Lock to ensure concurrency safety
	execMgmtLock.Lock()
	defer execMgmtLock.Unlock()

	// Encrypt resource requirements
	encryptedResources, err := common.EncryptData(resources)
	if err != nil {
		return fmt.Errorf("createSubExecutionContext: failed to encrypt resources for sub-context %s: %w", subContextID, err)
	}

	// Create sub-context
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.CreateSubContext(parentContextID, subContextID, encryptedResources); err != nil {
		return fmt.Errorf("createSubExecutionContext: failed to create sub-context %s: %w", subContextID, err)
	}

	// Log success
	log.Printf("Sub-context %s created successfully under parent context %s.", subContextID, parentContextID)
	return nil
}


// mergeSubExecutionContext merges two sub-execution contexts into a new merged context.
func mergeSubExecutionContext(ledgerInstance *ledger.Ledger, mainContextID, subContextID1, subContextID2, mergedContextID string) error {
	// Validate inputs
	if ledgerInstance == nil {
		return fmt.Errorf("mergeSubExecutionContext: ledger instance cannot be nil")
	}
	if mainContextID == "" || subContextID1 == "" || subContextID2 == "" || mergedContextID == "" {
		return fmt.Errorf("mergeSubExecutionContext: all context IDs must be provided")
	}

	// Retrieve states of the sub-contexts
	subState1, err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.RetrieveSubContextState(subContextID1)
	if err != nil {
		return fmt.Errorf("mergeSubExecutionContext: failed to retrieve state for sub-context %s: %w", subContextID1, err)
	}

	subState2, err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.RetrieveSubContextState(subContextID2)
	if err != nil {
		return fmt.Errorf("mergeSubExecutionContext: failed to retrieve state for sub-context %s: %w", subContextID2, err)
	}

	// Merge the states
	mergedState := MergeExecutionStates(subState1, subState2)

	// Save the merged state
	if err := ledgerInstance.EnvironmentSystemCoreLedger.ExecutionManager.SaveMergedContextState(mainContextID, mergedContextID, mergedState); err != nil {
		return fmt.Errorf("mergeSubExecutionContext: failed to save merged state for context %s: %w", mergedContextID, err)
	}

	// Log success
	log.Printf("Sub-contexts %s and %s merged successfully into %s under main context %s.", subContextID1, subContextID2, mergedContextID, mainContextID)
	return nil
}

