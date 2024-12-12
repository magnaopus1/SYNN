package utility

import (
	"errors"
	"log"
	"sync"
	"synnergy_network/pkg/common"
)

var diagnosticMode bool
var errorFrequencyMap sync.Map
var recoveryPolicy string
var memoryCheckpoint SystemState // Hypothetical struct for memory checkpoint
var errorLog []string

// CheckErrorThreshold: Checks if the error count exceeds a threshold
func CheckErrorThreshold(processID string, threshold int) bool {
    count, _ := errorFrequencyMap.LoadOrStore(processID, 0)
    if count.(int) >= threshold {
        LogDiagnostic("Error Threshold Exceeded", processID)
        return true
    }
    return false
}

// MarkTransactionInvalid: Flags a transaction as invalid in the ledger
func MarkTransactionInvalid(txID string) error {
    err := common.ledger.FlagTransactionInvalid(txID)
    if err != nil {
        CaptureError("MarkTransactionInvalid", "Failed to mark transaction as invalid: " + err.Error())
        return err
    }
    LogDiagnostic("Transaction Invalid", "Transaction " + txID + " marked as invalid")
    return nil
}

// ClearAllErrors: Clears all tracked errors and resets logs
func ClearAllErrors() {
    errorFrequencyMap.Range(func(key, value interface{}) bool {
        errorFrequencyMap.Store(key, 0)
        return true
    })
    errorLog = []string{}
    LogDiagnostic("Clear Errors", "All errors have been cleared")
}

// TrackErrorFrequency: Increments error frequency for a given process
func TrackErrorFrequency(processID string) {
    count, _ := errorFrequencyMap.LoadOrStore(processID, 0)
    errorFrequencyMap.Store(processID, count.(int)+1)
}

// TriggerSoftRevert: Softly reverts specific non-critical actions without a full rollback
func TriggerSoftRevert(actionID string) error {
    err := common.ledger.SoftRevertAction(actionID)
    if err != nil {
        CaptureError("TriggerSoftRevert", "Failed to soft revert action: " + err.Error())
        return err
    }
    LogDiagnostic("Soft Revert", "Soft revert executed for action " + actionID)
    return nil
}

// SetRecoveryPolicy: Configures the system's recovery policy
func SetRecoveryPolicy(policy string) {
    recoveryPolicy = policy
    LogDiagnostic("Recovery Policy", "Recovery policy set to " + policy)
}

// ValidateTransactionStatus: Validates a transaction’s status within Synnergy Consensus
func ValidateTransactionStatus(txID string) error {
    status, err := common.synnergyConsensus.ValidateTransaction(txID)
    if err != nil {
        CaptureError("ValidateTransactionStatus", "Failed to validate transaction: " + err.Error())
        return err
    }
    if !status {
        return errors.New("transaction validation failed")
    }
    LogDiagnostic("Transaction Validation", "Transaction " + txID + " validated successfully")
    return nil
}

// EnableDiagnosticMode: Activates diagnostic mode for enhanced logging
func EnableDiagnosticMode() {
    diagnosticMode = true
    LogDiagnostic("Diagnostic Mode", "Diagnostic mode enabled")
}

// DisableDiagnosticMode: Deactivates diagnostic mode
func DisableDiagnosticMode() {
    diagnosticMode = false
    LogDiagnostic("Diagnostic Mode", "Diagnostic mode disabled")
}

// CheckpointMemoryState: Creates a memory checkpoint for potential recovery
func CheckpointMemoryState(state common.SystemState) {
    memoryCheckpoint = state
    LogDiagnostic("Memory Checkpoint", "Memory checkpoint established")
}

// RevertToMemoryCheckpoint: Restores system state to a prior memory checkpoint
func RevertToMemoryCheckpoint() error {
    if memoryCheckpoint == nil {
        return CaptureError("RevertToMemoryCheckpoint", "No memory checkpoint available")
    }
    err := common.ledger.RestoreState(memoryCheckpoint)
    if err != nil {
        CaptureError("RevertToMemoryCheckpoint", "Failed to revert to memory checkpoint: " + err.Error())
        return err
    }
    LogDiagnostic("Memory Revert", "Reverted to memory checkpoint")
    return nil
}

// InitiateErrorRecovery: Begins the error recovery process based on policy
func InitiateErrorRecovery() {
    if recoveryPolicy == "Auto" {
        LogDiagnostic("Error Recovery", "Automatic recovery initiated")
        // Hypothetical recovery logic based on ledger integration
        common.ledger.InitiateAutoRecovery()
    } else {
        LogDiagnostic("Error Recovery", "Manual recovery required as per policy")
    }
}

// PauseOnError: Pauses operations on error, logging the issue
func PauseOnError(message string) {
    CaptureError("PauseOnError", message)
    log.Println("Operations paused due to error:", message)
}

// ResumeAfterError: Resumes operations after error resolution
func ResumeAfterError() {
    log.Println("Operations resumed after error resolution")
    LogDiagnostic("Error Resolution", "Resumed operations")
}

// QueryErrorLog: Retrieves recent error log entries
func QueryErrorLog(limit int) []string {
    if limit > len(errorLog) {
        limit = len(errorLog)
    }
    return errorLog[len(errorLog)-limit:]
}
