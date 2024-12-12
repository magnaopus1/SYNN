package utility

import (
	"log"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

var executionPaused bool
var systemUptime time.Time
var overflowFlag, signFlag, interruptFlag, parityFlag, auxiliaryFlag bool
var intervalExecutionMap sync.Map

// WAIT_UNTIL_BLOCK: Pauses execution until a specified block number is reached
func WAIT_UNTIL_BLOCK(blockNumber int) error {
    for {
        currentBlock, err := common.ledger.GetCurrentBlockNumber()
        if err != nil {
            return err
        }
        if currentBlock >= blockNumber {
            break
        }
        time.Sleep(time.Second)
    }
    return nil
}

// WAIT_UNTIL_TIME: Pauses execution until a specified time is reached
func WAIT_UNTIL_TIME(targetTime time.Time) {
    for time.Now().Before(targetTime) {
        time.Sleep(time.Second)
    }
}

// SCHEDULE_EVENT_AT_BLOCK: Schedules an event to execute at a specific block number
func SCHEDULE_EVENT_AT_BLOCK(blockNumber int, event func()) {
    go func() {
        WAIT_UNTIL_BLOCK(blockNumber)
        event()
    }()
}

// SCHEDULE_EVENT_AT_TIME: Schedules an event to execute at a specific time
func SCHEDULE_EVENT_AT_TIME(targetTime time.Time, event func()) {
    go func() {
        WAIT_UNTIL_TIME(targetTime)
        event()
    }()
}

// SET_INTERVAL_EXECUTION: Sets a function to execute at regular intervals
func SET_INTERVAL_EXECUTION(interval time.Duration, functionID string, fn func()) {
    intervalExecutionMap.Store(functionID, interval)
    go func() {
        for {
            if executionPaused {
                continue
            }
            fn()
            time.Sleep(interval)
        }
    }()
}

// SELF_DESTRUCT: Completely stops the program, removing all active state
func SELF_DESTRUCT() {
    common.ledger.ClearAllState()
    log.Fatal("System self-destructed")
}

// DESTROY_CONTRACT: Destroys a specified contract and logs the action
func DESTROY_CONTRACT(contractID string) error {
    err := common.ledger.DeleteContract(contractID)
    if err != nil {
        return err
    }
    LogDiagnostic("Contract Destruction", "Contract " + contractID + " destroyed")
    return nil
}

// STOP_EXECUTION: Stops all scheduled executions until resumed
func STOP_EXECUTION() {
    executionPaused = true
    LogDiagnostic("Execution Control", "Execution stopped")
}

// EXECUTE_AFTER_BLOCK: Executes a function after reaching a specific block
func EXECUTE_AFTER_BLOCK(blockNumber int, fn func()) {
    go func() {
        WAIT_UNTIL_BLOCK(blockNumber)
        if !executionPaused {
            fn()
        }
    }()
}

// EXECUTE_AFTER_TIME: Executes a function after a specified duration
func EXECUTE_AFTER_TIME(duration time.Duration, fn func()) {
    time.Sleep(duration)
    if !executionPaused {
        fn()
    }
}

// PAUSE_EXECUTION: Pauses execution of all scheduled events
func PAUSE_EXECUTION() {
    executionPaused = true
    LogDiagnostic("Execution Paused", "All executions are paused")
}

// RESUME_EXECUTION: Resumes execution of all scheduled events
func RESUME_EXECUTION() {
    executionPaused = false
    LogDiagnostic("Execution Resumed", "All executions are resumed")
}

// RESTART_SYSTEM: Restarts the blockchain system, resetting state
func RESTART_SYSTEM() error {
    err := common.ledger.RestartSystem()
    if err != nil {
        return err
    }
    systemUptime = time.Now()
    LogDiagnostic("System Restart", "System restarted successfully")
    return nil
}

// CHECK_SYSTEM_UPTIME: Returns the duration since the last system restart
func CHECK_SYSTEM_UPTIME() time.Duration {
    return time.Since(systemUptime)
}

// GET_EXECUTION_INTERVAL: Retrieves the set execution interval for a specific function
func GET_EXECUTION_INTERVAL(functionID string) (time.Duration, error) {
    interval, ok := intervalExecutionMap.Load(functionID)
    if !ok {
        return 0, errors.New("function ID not found")
    }
    return interval.(time.Duration), nil
}

// ResumeProgram: Resumes execution if previously paused
func ResumeProgram() {
    RESUME_EXECUTION()
}

// ProgramStatusCheck: Checks if the program is currently paused
func ProgramStatusCheck() bool {
    return executionPaused
}

// SetOverflowFlag: Sets the overflow flag
func SetOverflowFlag() {
    overflowFlag = true
    LogDiagnostic("Overflow Flag", "Overflow flag set")
}

// ClearOverflowFlag: Clears the overflow flag
func ClearOverflowFlag() {
    overflowFlag = false
    LogDiagnostic("Overflow Flag", "Overflow flag cleared")
}

// CheckOverflowFlag: Checks the status of the overflow flag
func CheckOverflowFlag() bool {
    return overflowFlag
}

// SetSignFlag: Sets the sign flag
func SetSignFlag() {
    signFlag = true
    LogDiagnostic("Sign Flag", "Sign flag set")
}

// ClearSignFlag: Clears the sign flag
func ClearSignFlag() {
    signFlag = false
    LogDiagnostic("Sign Flag", "Sign flag cleared")
}

// CheckSignFlag: Checks the status of the sign flag
func CheckSignFlag() bool {
    return signFlag
}

// SetInterruptFlag: Sets the interrupt flag
func SetInterruptFlag() {
    interruptFlag = true
    LogDiagnostic("Interrupt Flag", "Interrupt flag set")
}

// ClearInterruptFlag: Clears the interrupt flag
func ClearInterruptFlag() {
    interruptFlag = false
    LogDiagnostic("Interrupt Flag", "Interrupt flag cleared")
}

// CheckInterruptFlag: Checks the status of the interrupt flag
func CheckInterruptFlag() bool {
    return interruptFlag
}

// SetParityFlag: Sets the parity flag
func SetParityFlag() {
    parityFlag = true
    LogDiagnostic("Parity Flag", "Parity flag set")
}

// ClearParityFlag: Clears the parity flag
func ClearParityFlag() {
    parityFlag = false
    LogDiagnostic("Parity Flag", "Parity flag cleared")
}

// CheckParityFlag: Checks the status of the parity flag
func CheckParityFlag() bool {
    return parityFlag
}

// SetAuxiliaryFlag: Sets the auxiliary flag
func SetAuxiliaryFlag() {
    auxiliaryFlag = true
    LogDiagnostic("Auxiliary Flag", "Auxiliary flag set")
}
