package utility

import (
	"log"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

var errorCounter sync.Map
var failsafeMode bool
var errorThreshold int
var lastKnownGoodState common.BlockchainState // Hypothetical struct representing last stable state

// IncrementErrorCounter: Increments the error counter for a given process
func IncrementErrorCounter(processID string) {
    count, _ := errorCounter.LoadOrStore(processID, 0)
    errorCounter.Store(processID, count.(int)+1)
}

// ResetErrorCounter: Resets the error counter for a specified process
func ResetErrorCounter(processID string) {
    errorCounter.Store(processID, 0)
}

// RollbackOnFailure: Initiates a rollback in case of failure
func RollbackOnFailure(operationID string) error {
    err := ledger.RevertOperation(operationID)
    if err != nil {
        return err
    }
    LogErrorStackTrace(operationID, "Operation reverted due to failure")
    return nil
}

// ValidateRecoveryCompletion: Checks if recovery from an error was successful
func ValidateRecoveryCompletion(operationID string) bool {
    status, err := ledger.GetOperationStatus(operationID)
    if err != nil || status != "Recovered" {
        return false
    }
    return true
}

// SetFailsafeMode: Enables failsafe mode for secure error handling
func SetFailsafeMode() {
    failsafeMode = true
    log.Println("Failsafe mode enabled")
}

// DisableFailsafeMode: Disables failsafe mode
func DisableFailsafeMode() {
    failsafeMode = false
    log.Println("Failsafe mode disabled")
}

// MarkForRetry: Flags an operation to retry on failure
func MarkForRetry(operationID string) {
    common.ledger.MarkForRetry(operationID)
    LogErrorStackTrace(operationID, "Marked for retry")
}

// RevertToLastKnownGoodState: Reverts the system to the last stable state
func RevertToLastKnownGoodState() error {
    err := common.ledger.RestoreState(lastKnownGoodState)
    if err != nil {
        return err
    }
    LogErrorStackTrace("System", "Reverted to last known good state")
    return nil
}

// LogErrorStackTrace: Logs an encrypted stack trace of the error
func LogErrorStackTrace(context, message string) {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        log.Println("Error encrypting stack trace:", err)
        return
    }
    common.ledger.LogEvent(context, encryptedMessage)
}

// CheckRecoveryStatus: Checks the recovery status of an operation
func CheckRecoveryStatus(operationID string) string {
    status, err := common.ledger.GetRecoveryStatus(operationID)
    if err != nil {
        LogErrorStackTrace(operationID, "Failed to retrieve recovery status")
        return "Unknown"
    }
    return status
}

// DefineErrorThreshold: Sets a threshold for error count before triggers
func DefineErrorThreshold(threshold int) {
    errorThreshold = threshold
    log.Println("Error threshold set to", threshold)
}

// ValidateRollbackSuccess: Verifies if rollback was successful for an operation
func ValidateRollbackSuccess(operationID string) bool {
    success, err := common.ledger.CheckRollbackSuccess(operationID)
    if err != nil || !success {
        LogErrorStackTrace(operationID, "Rollback validation failed")
        return false
    }
    return true
}

// InitializeErrorMonitor: Starts error monitoring for specified processes
func InitializeErrorMonitor(processes []string) {
    for _, processID := range processes {
        go func(id string) {
            for {
                count, _ := errorCounter.LoadOrStore(id, 0)
                if count.(int) >= errorThreshold {
                    LogErrorStackTrace(id, "Error threshold reached, initiating rollback")
                    RollbackOnFailure(id)
                }
                time.Sleep(1 * time.Minute)
            }
        }(processID)
    }
}

// DeactivateErrorMonitor: Disables error monitoring and resets counters
func DeactivateErrorMonitor(processes []string) {
    for _, processID := range processes {
        ResetErrorCounter(processID)
    }
    log.Println("Error monitoring deactivated")
}

// SetCriticalErrorFlag: Sets a flag for critical errors that require immediate attention
func SetCriticalErrorFlag() {
    criticalErrorFlag = true
    log.Println("Critical error flag set")
}
