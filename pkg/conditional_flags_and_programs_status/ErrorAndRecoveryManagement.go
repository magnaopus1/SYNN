package conditional_flags_and_programs_status

import (
	"fmt"
	"synnergy_network/pkg/ledger"
	"time"
)

// SystemErrorEntry represents a system error log entry
type SystemErrorEntry struct {
	ErrorID     string    // Unique identifier for the error
	Description string    // Description of the error
	Timestamp   time.Time // Time the error occurred
}

// ProgramStatusEntry represents a program status log entry
type ProgramStatusEntry struct {
	ProgramID string    // Unique identifier for the program
	Status    string    // Current status of the program
	ErrorCode int       // Error code if applicable
	Timestamp time.Time // Time the status was recorded
}

// ConditionLogEntry represents a condition log entry
type ConditionLogEntry struct {
	ConditionID string    // Unique identifier for the condition
	Status      string    // Status of the condition
	Timestamp   time.Time // Time the condition was logged
}

// flagSystemError logs a system error with a unique error identifier and description
func FlagSystemError(errorID string, description string) error {
	entry := ledger.SystemErrorEntry{
		ErrorID:     errorID,
		Description: description,
		Timestamp:   time.Now(),
	}

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Record the system error
	if err := l.ConditionalFlagsLedger.RecordSystemError(entry); err != nil {
		return fmt.Errorf("failed to flag system error: %v", err)
	}
	return nil
}

// resetErrorFlag clears an error flag, indicating the system is back to a normal state
func ResetErrorFlag(errorID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Reset the error flag
	if err := l.ConditionalFlagsLedger.ResetErrorFlag(errorID); err != nil {
		return fmt.Errorf("failed to reset error flag: %v", err)
	}
	return nil
}

// activateRecoveryMode enables recovery mode to stabilize the system during errors
func ActivateRecoveryMode() error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Set the global recovery mode flag
	if err := l.ConditionalFlagsLedger.SetGlobalFlag(RecoveryMode, true); err != nil {
		return fmt.Errorf("failed to activate recovery mode: %v", err)
	}
	return nil
}

// deactivateRecoveryMode disables recovery mode once the system stabilizes
func DeactivateRecoveryMode() error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Clear the global recovery mode flag
	if err := l.ConditionalFlagsLedger.SetGlobalFlag(RecoveryMode, false); err != nil {
		return fmt.Errorf("failed to deactivate recovery mode: %v", err)
	}
	return nil
}

// lockProgramStatus locks the program status to prevent further changes
func LockProgramStatus(programID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Lock the program status
	if err := l.ConditionalFlagsLedger.LockStatus(programID); err != nil {
		return fmt.Errorf("failed to lock program status: %v", err)
	}
	return nil
}

// unlockProgramStatus unlocks the program status to allow modifications
func UnlockProgramStatus(programID string) error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.UnlockStatus(programID); err != nil {
		return fmt.Errorf("failed to unlock program status: %v", err)
	}
	return nil
}

// resetProgramFlags resets all flags related to a particular program
func ResetProgramFlags(programID string) error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.ResetFlags(programID); err != nil {
		return fmt.Errorf("failed to reset program flags: %v", err)
	}
	return nil
}

// checkStatusFlags checks all status flags related to a program
func CheckStatusFlags(programID string) (map[string]bool, error) {
	l := &ledger.Ledger{}
	flags, err := l.ConditionalFlagsLedger.CheckStatusFlags(programID)
	if err != nil {
		return nil, fmt.Errorf("failed to check status flags: %v", err)
	}
	return flags, nil
}

// logProgramStatus records the current status of a program, including any error codes
func LogProgramStatus(programID string, status string, errorCode int) error {
	entry := ledger.ProgramStatusEntry{
		ProgramID: programID,
		Status:    status,
		ErrorCode: errorCode,
		Timestamp: time.Now(),
	}
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.LogProgramStatus(entry); err != nil {
		return fmt.Errorf("failed to log program status: %v", err)
	}
	return nil
}

// logConditionCheck logs a specific condition check, useful for debugging
func LogConditionCheck(conditionID string, status string) error {
	entry := ledger.ConditionLogEntry{
		ConditionID: conditionID,
		Status:      status,
		Timestamp:   time.Now(),
	}
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.LogCondition(entry); err != nil {
		return fmt.Errorf("failed to log condition check: %v", err)
	}
	return nil
}

// programErrorStatus checks the current error status of a specific program
func ProgramErrorStatus(programID string) (bool, error) {
	l := &ledger.Ledger{}
	status, err := l.ConditionalFlagsLedger.CheckErrorStatus(programID)
	if err != nil {
		return false, fmt.Errorf("failed to check program error status: %v", err)
	}
	return status, nil
}

// toggleInterrupt toggles an interrupt flag for immediate system attention
func ToggleInterrupt(interruptID string, enabled bool) error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(interruptID, InterruptFlag, enabled); err != nil {
		return fmt.Errorf("failed to toggle interrupt flag: %v", err)
	}
	return nil
}

// setRecoveryFlag sets a recovery flag for a specific issue, marking it for recovery
func SetRecoveryFlag(flagID string) error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(flagID, RecoveryFlag, true); err != nil {
		return fmt.Errorf("failed to set recovery flag: %v", err)
	}
	return nil
}

// clearRecoveryFlag clears an existing recovery flag, marking it as resolved
func ClearRecoveryFlag(flagID string) error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(flagID, RecoveryFlag, false); err != nil {
		return fmt.Errorf("failed to clear recovery flag: %v", err)
	}
	return nil
}

// checkRecoveryFlag verifies if a specific recovery flag is set
func CheckRecoveryFlag(flagID string) (bool, error) {
	l := &ledger.Ledger{}
	flag, err := l.ConditionalFlagsLedger.CheckFlag(flagID, RecoveryFlag)
	if err != nil {
		return false, fmt.Errorf("failed to check recovery flag: %v", err)
	}
	return flag, nil
}
