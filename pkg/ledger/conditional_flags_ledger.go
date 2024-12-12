package ledger

import "fmt"

func (l *ConditionalFlagsLedger) EnableCondition(conditionID string) error {
	l.Conditions[conditionID] = true
	return nil
}

func (l *ConditionalFlagsLedger) DisableCondition(conditionID string) error {
	l.Conditions[conditionID] = false
	return nil
}

func (l *ConditionalFlagsLedger) TrackExecutionPath(entry ExecutionPathEntry) error {
	l.ExecutionPaths = append(l.ExecutionPaths, entry)
	return nil
}

func (l *ConditionalFlagsLedger) SetFlag(id string, flagType string, state bool) error {
	if l.Flags[id] == nil {
		l.Flags[id] = make(map[string]bool)
	}
	l.Flags[id][flagType] = state
	return nil
}


func (l *ConditionalFlagsLedger) CheckFlag(id string, flagType string) (bool, error) {
	// Check if the flag exists for the given ID and type
	if flagMap, exists := l.Flags[id]; exists {
		if flag, found := flagMap[flagType]; found {
			return flag, nil
		}
	}
	return false, fmt.Errorf("flag of type %s for ID %s not found", flagType, id)
}

func (l *ConditionalFlagsLedger) RecordSystemError(entry SystemErrorEntry) error {
	l.Lock()
	defer l.Unlock()
	l.SystemErrors[entry.ErrorID] = entry
	return nil
}

func (l *ConditionalFlagsLedger) ResetErrorFlag(errorID string) error {
	l.Lock()
	defer l.Unlock()
	delete(l.SystemErrors, errorID)
	return nil
}


func (l *ConditionalFlagsLedger) SetGlobalFlag(flagType string, state bool) error {
	l.Lock()
	defer l.Unlock()
	l.GlobalFlags[flagType] = state
	return nil
}

func (l *ConditionalFlagsLedger) LockStatus(programID string) error {
	l.Lock()
	defer l.Unlock()
	l.StatusLocks[programID] = true
	return nil
}

func (l *ConditionalFlagsLedger) UnlockStatus(programID string) error {
	l.Lock()
	defer l.Unlock()
	l.StatusLocks[programID] = false
	return nil
}

func (l *ConditionalFlagsLedger) ResetFlags(programID string) error {
	l.Lock()
	defer l.Unlock()
	delete(l.ProgramFlags, programID)
	return nil
}

func (l *ConditionalFlagsLedger) CheckStatusFlags(programID string) (map[string]bool, error) {
	l.Lock()
	defer l.Unlock()
	flags, exists := l.ProgramFlags[programID]
	if !exists {
		return nil, fmt.Errorf("no flags found for program ID: %s", programID)
	}
	return flags, nil
}

func (l *ConditionalFlagsLedger) LogProgramStatus(entry ProgramStatusEntry) error {
	l.Lock()
	defer l.Unlock()
	l.ProgramLogs[entry.ProgramID] = entry
	return nil
}

func (l *ConditionalFlagsLedger) LogCondition(entry ConditionLogEntry) error {
	l.Lock()
	defer l.Unlock()
	l.ConditionLogs[entry.ConditionID] = entry
	return nil
}

func (l *ConditionalFlagsLedger) CheckErrorStatus(programID string) (bool, error) {
	l.Lock()
	defer l.Unlock()
	entry, exists := l.ProgramLogs[programID]
	if !exists {
		return false, fmt.Errorf("no error status found for program ID: %s", programID)
	}
	return entry.ErrorCode != 0, nil
}

