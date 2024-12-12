package utility

import (
	"errors"
	"synnergy_network/pkg/common"
)

var criticalErrorFlag, zeroFlag, carryFlag bool
var alertTriggered bool
var diagnosticLog []string

// ClearCriticalErrorFlag: Clears the critical error flag and logs the action
func ClearCriticalErrorFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    criticalErrorFlag = false
    LogDiagnostic("Critical Error Flag", "Critical error flag cleared")
}

// TriggerAlertOnError: Triggers an alert if a critical error flag is set
func TriggerAlertOnError() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    if criticalErrorFlag && !alertTriggered {
        alertTriggered = true
        LogDiagnostic("Error Alert", "Critical error alert triggered")
        // Hypothetical alerting system integration
        common.alertSystem.SendAlert("Critical error detected in the system")
    }
}

// ValidateErrorHandler: Validates the presence of an error handler
func ValidateErrorHandler() bool {
    return common.ledger.CheckErrorHandler() // Assumes ledger has a method for checking error handler
}

// RollbackResourceAllocation: Rolls back resources allocated to a transaction on error
func RollbackResourceAllocation(txID string) error {
    err := common.ledger.RevertAllocation(txID)
    if err != nil {
        LogDiagnostic("Resource Rollback", "Failed to rollback resources for transaction: " + txID)
        return err
    }
    LogRollbackEvent(txID, "Resource allocation rolled back")
    return nil
}

// ResetDiagnosticLog: Clears the diagnostic log for fresh logging
func ResetDiagnosticLog() {
    diagnosticLog = []string{}
    LogDiagnostic("Diagnostic Log", "Diagnostic log reset")
}

// RetryOnError: Retries a failed transaction if an error occurs
func RetryOnError(txID string) error {
    err := common.ledger.RetryTransaction(txID)
    if err != nil {
        LogDiagnostic("Transaction Retry", "Failed to retry transaction: " + txID)
        return err
    }
    LogDiagnostic("Transaction Retry", "Transaction " + txID + " retried successfully")
    return nil
}

// LogRollbackEvent: Logs a rollback event to the ledger securely
func LogRollbackEvent(txID, message string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("RollbackEvent-"+txID, encryptedMessage)
}

// ValidateTransactionIntegrity: Checks if a transaction maintains integrity in Synnergy Consensus
func ValidateTransactionIntegrity(txID string) error {
    isValid, err := common.synnergyConsensus.ValidateTransactionIntegrity(txID)
    if err != nil {
        LogDiagnostic("Transaction Integrity", "Failed to validate transaction integrity: " + txID)
        return err
    }
    if !isValid {
        return errors.New("transaction integrity validation failed")
    }
    LogDiagnostic("Transaction Integrity", "Transaction " + txID + " integrity validated")
    return nil
}

// SetZeroFlag: Sets the zero flag to indicate a zero-state condition
func SetZeroFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    zeroFlag = true
    LogDiagnostic("Zero Flag", "Zero flag set")
}

// ClearZeroFlag: Clears the zero flag
func ClearZeroFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    zeroFlag = false
    LogDiagnostic("Zero Flag", "Zero flag cleared")
}

// SetCarryFlag: Sets the carry flag for carrying data over operations
func SetCarryFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    carryFlag = true
    LogDiagnostic("Carry Flag", "Carry flag set")
}

// ClearCarryFlag: Clears the carry flag
func ClearCarryFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    carryFlag = false
    LogDiagnostic("Carry Flag", "Carry flag cleared")
}

// CheckZeroFlag: Checks the status of the zero flag
func CheckZeroFlag() bool {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    return zeroFlag
}

// CheckCarryFlag: Checks the status of the carry flag
func CheckCarryFlag() bool {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    return carryFlag
}

// HaltProgram: Halts program execution completely on critical failure
func HaltProgram(message string) {
    CaptureError("Program Halt", message)
    log.Fatal("System halted due to critical error: " + message)
}

// LogDiagnostic: Helper function for encrypted diagnostic logging
func LogDiagnostic(context, message string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogDiagnostic(context, encryptedMessage)
}

// CaptureError: Records and logs an error securely
func CaptureError(context, message string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogError(context, encryptedMessage)
}
