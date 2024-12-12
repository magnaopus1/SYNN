package conditional_flags_and_programs_status

import (
	"fmt"
	"synnergy_network/pkg/ledger"
	"time"
)

// Constants for different flag types
const (
	LoopFlag       = "loop"
	EqualityFlag   = "equality"
	WaitFlag       = "wait"
	HaltFlag       = "halt"
	PauseFlag      = "pause"
	CompletionFlag = "completion"
	RecoveryFlag   = "recovery"
	RecoveryMode   = "recovery_mode"
	InterruptFlag  = "interrupt"
	ZeroFlag       = "zero"
	CarryFlag      = "carry"
	OverflowFlag   = "overflow"
	SignFlag       = "sign"
	ParityFlag     = "parity"
	AuxiliaryFlag  = "auxiliary"
	NegativeFlag   = "negative"
	PositiveFlag   = "positive"
)

type ExecutionPathEntry struct {
	Path      string    // The path in the program flow
	Timestamp time.Time // Time when this path was tracked
}

// enableConditionCheck activates a condition check within the program flow
func EnableConditionCheck(conditionID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Enable the specified condition
	if err := l.ConditionalFlagsLedger.EnableCondition(conditionID); err != nil {
		return fmt.Errorf("failed to enable condition check: %v", err)
	}
	return nil
}

// disableConditionCheck deactivates a specific condition check
func DisableConditionCheck(conditionID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Disable the specified condition
	if err := l.ConditionalFlagsLedger.DisableCondition(conditionID); err != nil {
		return fmt.Errorf("failed to disable condition check: %v", err)
	}
	return nil
}

// trackExecutionPath logs the execution path of the program for tracking and debugging
func TrackExecutionPath(path string) error {
	entry := ledger.ExecutionPathEntry{
		Path:      path,
		Timestamp: time.Now(),
	}

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Track the execution path
	if err := l.ConditionalFlagsLedger.TrackExecutionPath(entry); err != nil {
		return fmt.Errorf("failed to track execution path: %v", err)
	}
	return nil
}

// setLoopFlag sets a loop flag for a specific operation
func SetLoopFlag(loopID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Set the loop flag
	if err := l.ConditionalFlagsLedger.SetFlag(loopID, LoopFlag, true); err != nil {
		return fmt.Errorf("failed to set loop flag: %v", err)
	}
	return nil
}

// clearLoopFlag clears an existing loop flag
func ClearLoopFlag(loopID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Clear the loop flag
	if err := l.ConditionalFlagsLedger.SetFlag(loopID, LoopFlag, false); err != nil {
		return fmt.Errorf("failed to clear loop flag: %v", err)
	}
	return nil
}

// checkLoopFlag verifies the status of a loop flag
func CheckLoopFlag(loopID string) (bool, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Check the loop flag
	flag, err := l.ConditionalFlagsLedger.CheckFlag(loopID, LoopFlag)
	if err != nil {
		return false, fmt.Errorf("failed to check loop flag: %v", err)
	}
	return flag, nil
}

// setEqualityFlag sets an equality flag for a condition check
func SetEqualityFlag(flagID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Set the equality flag
	if err := l.ConditionalFlagsLedger.SetFlag(flagID, EqualityFlag, true); err != nil {
		return fmt.Errorf("failed to set equality flag: %v", err)
	}
	return nil
}

// clearEqualityFlag clears an equality flag for a condition check
func ClearEqualityFlag(flagID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Clear the equality flag
	if err := l.ConditionalFlagsLedger.SetFlag(flagID, EqualityFlag, false); err != nil {
		return fmt.Errorf("failed to clear equality flag: %v", err)
	}
	return nil
}

// checkEqualityFlag verifies the status of an equality flag
func CheckEqualityFlag(flagID string) (bool, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Check the equality flag
	flag, err := l.ConditionalFlagsLedger.CheckFlag(flagID, EqualityFlag)
	if err != nil {
		return false, fmt.Errorf("failed to check equality flag: %v", err)
	}
	return flag, nil
}

// setWaitFlag sets a wait flag, indicating the program should pause until cleared
func SetWaitFlag(waitID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Set the wait flag
	if err := l.ConditionalFlagsLedger.SetFlag(waitID, WaitFlag, true); err != nil {
		return fmt.Errorf("failed to set wait flag: %v", err)
	}
	return nil
}

// clearWaitFlag clears a wait flag, allowing the program to resume
func ClearWaitFlag(waitID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Clear the wait flag
	if err := l.ConditionalFlagsLedger.SetFlag(waitID, WaitFlag, false); err != nil {
		return fmt.Errorf("failed to clear wait flag: %v", err)
	}
	return nil
}

// checkWaitFlag verifies the status of a wait flag
func CheckWaitFlag(waitID string) (bool, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Check the wait flag
	flag, err := l.ConditionalFlagsLedger.CheckFlag(waitID, WaitFlag)
	if err != nil {
		return false, fmt.Errorf("failed to check wait flag: %v", err)
	}
	return flag, nil
}
