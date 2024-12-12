package utility

import (
	"log"
	"sync"
	"synnergy_network/pkg/common"
)

var errorFlag bool
var rollbackPoint common.BlockchainState // Hypothetical struct for rollback checkpoints
var lastError string
var errorHandlerEnabled bool
var handlerMutex sync.Mutex



// RevertTransaction: Reverts a specific transaction within the ledger
func RevertTransaction(txID string) error {
    err := common.ledger.RevertTransaction(txID)
    if err != nil {
        CaptureError("RevertTransaction", "Failed to revert transaction: " + err.Error())
        return err
    }
    LogDiagnostic("Transaction Reverted", "Transaction " + txID + " successfully reverted")
    return nil
}



// SetErrorFlag: Sets an error flag for error tracking
func SetErrorFlag() {
    errorFlag = true
    LogDiagnostic("ErrorFlag", "Error flag set")
}

// ClearErrorFlag: Clears the error flag
func ClearErrorFlag() {
    errorFlag = false
    LogDiagnostic("ErrorFlag", "Error flag cleared")
}

// TriggerPanic: Forces a panic for critical failures
func TriggerPanic(message string) {
    CaptureError("PanicTriggered", message)
    log.Panic(message)
}

// SetRollbackPoint: Establishes a rollback checkpoint
func SetRollbackPoint(state common.BlockchainState) {
    rollbackPoint = state
    LogDiagnostic("Rollback Point", "Rollback point set")
}

// RollbackToCheckpoint: Rolls back system to a checkpoint
func RollbackToCheckpoint() error {
    if rollbackPoint == nil {
        return CaptureError("RollbackToCheckpoint", "No rollback point set")
    }
    err := common.ledger.RestoreState(rollbackPoint)
    if err != nil {
        CaptureError("RollbackToCheckpoint", "Failed to rollback: " + err.Error())
        return err
    }
    LogDiagnostic("Rollback", "Rolled back to checkpoint")
    return nil
}

// RetrieveLastError: Retrieves the last recorded error message
func RetrieveLastError() string {
    return lastError
}

// ValidateConsensusStatus: Validates consensus for the current transaction state
func ValidateConsensusStatus(txID string) error {
    isValid, err := common.synnergyConsensus.CheckStatus(txID)
    if err != nil {
        CaptureError("ValidateConsensusStatus", "Failed consensus validation for transaction " + txID)
        return err
    }
    if !isValid {
        CaptureError("ValidateConsensusStatus", "Consensus invalid for transaction " + txID)
        return errors.New("consensus validation failed")
    }
    LogDiagnostic("Consensus Validation", "Transaction " + txID + " validated by consensus")
    return nil
}

// RetryTransaction: Attempts to retry a failed transaction
func RetryTransaction(txID string) error {
    err := common.ledger.RetryTransaction(txID)
    if err != nil {
        CaptureError("RetryTransaction", "Failed to retry transaction: " + err.Error())
        return err
    }
    LogDiagnostic("Transaction Retry", "Transaction " + txID + " retried successfully")
    return nil
}

// RecordErrorEvent: Logs an error event to the ledger with encryption
func RecordErrorEvent(context, errorMsg string) error {
    encryptedError, err := encryption.Encrypt([]byte(errorMsg))
    if err != nil {
        return err
    }
    return common.ledger.RecordEvent(context, encryptedError)
}

// SetErrorHandler: Enables a custom error handler function
func SetErrorHandler(handler func(error)) {
    handlerMutex.Lock()
    defer handlerMutex.Unlock()
    errorHandlerEnabled = true
    log.Println("Custom error handler enabled")
}

// DisableErrorHandler: Disables the custom error handler
func DisableErrorHandler() {
    handlerMutex.Lock()
    defer handlerMutex.Unlock()
    errorHandlerEnabled = false
    log.Println("Custom error handler disabled")
}

// LogErrorDetails: Logs detailed error information securely
func LogErrorDetails(context, details string) error {
    encryptedDetails, err := encryption.Encrypt([]byte(details))
    if err != nil {
        return err
    }
    return common.ledger.LogErrorDetails(context, encryptedDetails)
}
