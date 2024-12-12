package conditional_flags_and_programs_status

import (
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
)

// Global mutex to control program state changes
var programLock sync.Mutex

// haltProgram stops the program execution completely
func HaltProgram() error {
	programLock.Lock()
	defer programLock.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Set the halt flag to true
	if err := l.ConditionalFlagsLedger.SetFlag(HaltFlag, HaltFlag, true); err != nil {
		return fmt.Errorf("failed to halt program: %v", err)
	}
	return nil
}

// resumeProgram resumes a halted program
func ResumeProgram() error {
	programLock.Lock()
	defer programLock.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Set the halt flag to false
	if err := l.ConditionalFlagsLedger.SetFlag(HaltFlag, HaltFlag, false); err != nil {
		return fmt.Errorf("failed to resume program: %v", err)
	}
	return nil
}

// pauseExecution pauses the program temporarily
func PauseExecution() error {
	programLock.Lock()
	defer programLock.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Set the pause flag to true
	if err := l.ConditionalFlagsLedger.SetFlag(PauseFlag, PauseFlag, true); err != nil {
		return fmt.Errorf("failed to pause execution: %v", err)
	}
	return nil
}

// resumeExecution resumes execution after a pause
func ResumeExecution() error {
	programLock.Lock()
	defer programLock.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Set the pause flag to false
	if err := l.ConditionalFlagsLedger.SetFlag(PauseFlag, PauseFlag, false); err != nil {
		return fmt.Errorf("failed to resume execution: %v", err)
	}
	return nil
}

// programStatusCheck checks the current status of the program
func ProgramStatusCheck() (string, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Check the halt flag
	isHalted, err := l.ConditionalFlagsLedger.CheckFlag(HaltFlag, HaltFlag)
	if err != nil && err.Error() != fmt.Sprintf("flag of type %s for ID %s not found", HaltFlag, HaltFlag) {
		return "", fmt.Errorf("failed to check halt status: %v", err)
	}

	// Check the pause flag
	isPaused, err := l.ConditionalFlagsLedger.CheckFlag(PauseFlag, PauseFlag)
	if err != nil && err.Error() != fmt.Sprintf("flag of type %s for ID %s not found", PauseFlag, PauseFlag) {
		return "", fmt.Errorf("failed to check pause status: %v", err)
	}

	switch {
	case isHalted:
		return "Halted", nil
	case isPaused:
		return "Paused", nil
	default:
		return "Running", nil
	}
}

// checkProgramCompletion checks if the program has completed execution
func CheckProgramCompletion() (bool, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Check the completion flag
	isComplete, err := l.ConditionalFlagsLedger.CheckFlag(CompletionFlag, CompletionFlag)
	if err != nil {
		return false, fmt.Errorf("failed to check completion status: %v", err)
	}
	return isComplete, nil
}

// setCompletionFlag sets the completion flag to indicate the program has finished
func SetCompletionFlag() error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Set the completion flag to true
	if err := l.ConditionalFlagsLedger.SetFlag(CompletionFlag, CompletionFlag, true); err != nil {
		return fmt.Errorf("failed to set completion flag: %v", err)
	}
	return nil
}

// clearCompletionFlag clears the completion flag
func ClearCompletionFlag() error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Set the completion flag to false
	if err := l.ConditionalFlagsLedger.SetFlag(CompletionFlag, CompletionFlag, false); err != nil {
		return fmt.Errorf("failed to clear completion flag: %v", err)
	}
	return nil
}
