package utility

import (
	"synnergy_network/pkg/common"
)

var recoveryFlag bool
var statusLocked bool
var interruptEnabled bool
var executionPath []string

// SetRecoveryFlag: Activates the recovery flag and logs the event
func SetRecoveryFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    recoveryFlag = true
    LogRecoveryEvent("Recovery flag set")
}

// ClearRecoveryFlag: Deactivates the recovery flag and logs the event
func ClearRecoveryFlag() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    recoveryFlag = false
    LogRecoveryEvent("Recovery flag cleared")
}

// CheckStatusFlags: Logs and retrieves the status of recovery and interrupt flags
func CheckStatusFlags() map[string]bool {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    status := map[string]bool{
        "recoveryFlag":    recoveryFlag,
        "interruptEnabled": interruptEnabled,
    }
    LogRecoveryEvent("Status flags checked: " + formatFlagStatus(status))
    return status
}

// LockProgramStatus: Locks the program's current status, preventing changes
func LockProgramStatus() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    statusLocked = true
    LogRecoveryEvent("Program status locked")
}

// UnlockProgramStatus: Unlocks the program's status, allowing changes
func UnlockProgramStatus() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    statusLocked = false
    LogRecoveryEvent("Program status unlocked")
}

// TrackExecutionPath: Records the opcode or function executed in the program path
func TrackExecutionPath(step string) {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    executionPath = append(executionPath, step)
    LogRecoveryEvent("Execution path tracked: " + step)
}

// LogConditionCheck: Logs the result of a conditional check securely
func LogConditionCheck(condition string, result bool) error {
    message := "Condition checked: " + condition + " - Result: " + boolToString(result)
    return LogRecoveryEvent(message)
}

// ToggleInterrupt: Toggles the interrupt flag and logs the current state
func ToggleInterrupt() {
    flagMutex.Lock()
    defer flagMutex.Unlock()
    interruptEnabled = !interruptEnabled
    LogRecoveryEvent("Interrupt flag toggled: " + boolToString(interruptEnabled))
}

// Helper Functions

// LogRecoveryEvent: Helper function to log recovery events securely
func LogRecoveryEvent(message string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("RecoveryEvent", encryptedMessage)
}

// formatFlagStatus: Formats flag statuses for logging
func formatFlagStatus(status map[string]bool) string {
    formatted := ""
    for key, value := range status {
        formatted += key + ": " + boolToString(value) + "; "
    }
    return formatted
}


