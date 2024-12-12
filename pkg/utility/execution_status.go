package utility

import (
	"sync"
	"synnergy_network/pkg/common"
)

var flagMutex sync.Mutex // Ensures thread-safe operations on flags

// ClearAuxiliaryFlag: Clears the auxiliary flag, logging the action securely
func ClearAuxiliaryFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    auxiliaryFlag = false
    LogDiagnostic("Auxiliary Flag", "Auxiliary flag cleared")
}

// SetNegativeFlag: Sets the negative flag to indicate a negative status
func SetNegativeFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    negativeFlag = true
    LogDiagnostic("Negative Flag", "Negative flag set")
}

// ClearNegativeFlag: Clears the negative flag
func ClearNegativeFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    negativeFlag = false
    LogDiagnostic("Negative Flag", "Negative flag cleared")
}

// CheckNegativeFlag: Returns the status of the negative flag
func CheckNegativeFlag() bool {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    return negativeFlag
}

// SetPositiveFlag: Sets the positive flag to indicate a positive status
func SetPositiveFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    positiveFlag = true
    LogDiagnostic("Positive Flag", "Positive flag set")
}

// ClearPositiveFlag: Clears the positive flag
func ClearPositiveFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    positiveFlag = false
    LogDiagnostic("Positive Flag", "Positive flag cleared")
}

// CheckPositiveFlag: Returns the status of the positive flag
func CheckPositiveFlag() bool {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    return positiveFlag
}

// SetEqualityFlag: Sets the equality flag when conditions are equal
func SetEqualityFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    equalityFlag = true
    LogDiagnostic("Equality Flag", "Equality flag set")
}

// ClearEqualityFlag: Clears the equality flag
func ClearEqualityFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    equalityFlag = false
    LogDiagnostic("Equality Flag", "Equality flag cleared")
}

// CheckEqualityFlag: Returns the status of the equality flag
func CheckEqualityFlag() bool {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    return equalityFlag
}

// PauseExecution: Pauses all executions, preventing further operations until resumed
func PauseExecution() {
    executionPaused = true
    LogDiagnostic("Execution Status", "Execution paused")
}

// ResumeExecution: Resumes all paused executions
func ResumeExecution() {
    executionPaused = false
    LogDiagnostic("Execution Status", "Execution resumed")
}

// SetWaitFlag: Sets the wait flag, indicating the system is in a wait state
func SetWaitFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    waitFlag = true
    LogDiagnostic("Wait Flag", "Wait flag set")
}

// ClearWaitFlag: Clears the wait flag, exiting the wait state
func ClearWaitFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    waitFlag = false
    LogDiagnostic("Wait Flag", "Wait flag cleared")
}

// CheckWaitFlag: Checks if the wait flag is currently set
func CheckWaitFlag() bool {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    return waitFlag
}

// ResetProgramFlags: Resets all program flags to their default state
func ResetProgramFlags() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    auxiliaryFlag, negativeFlag, positiveFlag, equalityFlag, waitFlag = false, false, false, false, false
    LogDiagnostic("Program Flags", "All program flags reset")
}

// LogDiagnostic: Helper function to log encrypted diagnostic messages
func LogDiagnostic(context, message string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogDiagnostic(context, encryptedMessage)
}
