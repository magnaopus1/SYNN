package utility

import (
    "synnergy_network/pkg/common"
)

var conditionCheckEnabled bool
var completionFlag bool
var loopFlag bool
var systemErrorFlag bool
var recoveryModeEnabled bool


// LogProgramStatus: Logs the current status of the program securely
func LogProgramStatus(statusMessage string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(statusMessage))
    if err != nil {
        return err
    }
    return common.ledger.LogProgramStatus("ProgramStatus", encryptedMessage)
}

// EnableConditionCheck: Enables condition-based checks within the program
func EnableConditionCheck() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    conditionCheckEnabled = true
    LogProgramStatus("Condition check enabled")
}

// DisableConditionCheck: Disables condition-based checks within the program
func DisableConditionCheck() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    conditionCheckEnabled = false
    LogProgramStatus("Condition check disabled")
}

// ProgramErrorStatus: Retrieves and logs the current error status
func ProgramErrorStatus() bool {
    LogProgramStatus("Error status: " + boolToString(systemErrorFlag))
    return systemErrorFlag
}

// CheckProgramCompletion: Checks if the completion flag is set
func CheckProgramCompletion() bool {
    LogProgramStatus("Program completion checked: " + boolToString(completionFlag))
    return completionFlag
}

// SetCompletionFlag: Sets the completion flag to indicate the program has completed
func SetCompletionFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    completionFlag = true
    LogProgramStatus("Completion flag set")
}

// ClearCompletionFlag: Clears the completion flag, resetting program state
func ClearCompletionFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    completionFlag = false
    LogProgramStatus("Completion flag cleared")
}

// SetLoopFlag: Sets the loop flag, enabling loop functionality in the program
func SetLoopFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    loopFlag = true
    LogProgramStatus("Loop flag set")
}

// ClearLoopFlag: Clears the loop flag, disabling loop functionality
func ClearLoopFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    loopFlag = false
    LogProgramStatus("Loop flag cleared")
}

// CheckLoopFlag: Checks the status of the loop flag
func CheckLoopFlag() bool {
    LogProgramStatus("Loop flag checked: " + boolToString(loopFlag))
    return loopFlag
}

// FlagSystemError: Flags a system error and logs it securely
func FlagSystemError(errorMessage string) {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    systemErrorFlag = true
    encryptedMessage, err := encryption.Encrypt([]byte(errorMessage))
    if err == nil {
        common.ledger.LogSystemError("SystemErrorFlag", encryptedMessage)
    }
}

// ResetErrorFlag: Clears the system error flag
func ResetErrorFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    systemErrorFlag = false
    LogProgramStatus("System error flag reset")
}

// CheckRecoveryFlag: Checks the current status of the recovery mode flag
func CheckRecoveryFlag() bool {
    LogProgramStatus("Recovery flag checked: " + boolToString(recoveryModeEnabled))
    return recoveryModeEnabled
}

// ActivateRecoveryMode: Activates recovery mode for the program
func ActivateRecoveryMode() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    recoveryModeEnabled = true
    LogProgramStatus("Recovery mode activated")
}

// DeactivateRecoveryMode: Deactivates recovery mode for the program
func DeactivateRecoveryMode() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    recoveryModeEnabled = false
    LogProgramStatus("Recovery mode deactivated")
}

