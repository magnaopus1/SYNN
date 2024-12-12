package utility

import (
	"crypto/rand"
	"errors"
	"math/big"
	"sync"
	"synnergy_network/pkg/common"
	"time"
)

var safeModeEnabled = false
var environmentVars sync.Map

// SET_SYSTEM_HEALTH_CHECK: Periodically checks and logs system health
func SET_SYSTEM_HEALTH_CHECK(interval time.Duration) {
    go func() {
        for {
            if !safeModeEnabled {
                healthStatus := common.GetSystemHealth() // Assume a health check method in common
                LOG_MESSAGE("System Health Check", healthStatus)
            }
            time.Sleep(interval)
        }
    }()
}

// EXECUTE_IN_SAFE_MODE: Enters safe mode, restricting operations
func EXECUTE_IN_SAFE_MODE() {
    safeModeEnabled = true
    LOG_MESSAGE("System", "Entered Safe Mode for secure operations")
}

// EXIT_SAFE_MODE: Exits safe mode, resuming full operations
func EXIT_SAFE_MODE() {
    safeModeEnabled = false
    LOG_MESSAGE("System", "Exited Safe Mode")
}

// UPDATE_SYSTEM_FLAG: Updates specified system flags within the ledger
func UPDATE_SYSTEM_FLAG(flagName string, value interface{}) error {
    return common.ledger.UpdateFlag(flagName, value)
}

// GET_ENVIRONMENT_INFO: Retrieves and logs environment-specific information
func GET_ENVIRONMENT_INFO() (string, error) {
    envInfo := common.GetEnvironmentDetails() // Assuming system call to gather environment details
    LOG_MESSAGE("Environment Info", envInfo)
    return envInfo, nil
}

// SET_ENV_VARIABLE: Securely sets an environment variable in encrypted form
func SET_ENV_VARIABLE(key, value string) error {
    encryptedValue, err := encryption.Encrypt([]byte(value))
    if err != nil {
        return err
    }
    environmentVars.Store(key, encryptedValue)
    return nil
}

// LOG_MESSAGE: Securely logs a message to the ledger
func LOG_MESSAGE(context, message string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent(context, encryptedMessage)
}

// DELAY_EXECUTION: Delays execution for a specified duration
func DELAY_EXECUTION(duration time.Duration) {
    time.Sleep(duration)
}

// GENERATE_RANDOM: Generates a secure random number up to a max value
func GENERATE_RANDOM(max int64) (int64, error) {
    n, err := rand.Int(rand.Reader, big.NewInt(max))
    if err != nil {
        return 0, err
    }
    return n.Int64(), nil
}

// TRIGGER_EVENT: Triggers an event within the blockchain environment
func TRIGGER_EVENT(eventName string, data interface{}) error {
    return common.ledger.TriggerEvent(eventName, data)
}

// VALIDATE_SUB_BLOCK: Validates a sub-block through the Synnergy Consensus
func VALIDATE_SUB_BLOCK(subBlock common.SubBlock) error {
    // Encrypt transaction data before validation
    encryptedSubBlock, err := encryption.EncryptSubBlock(subBlock)
    if err != nil {
        return err
    }
    isValid, err := common.synnergyConsensus.ValidateSubBlock(encryptedSubBlock)
    if err != nil || !isValid {
        return errors.New("sub-block validation failed")
    }
    LOG_MESSAGE("Validation", "Sub-block validated successfully")
    return nil
}

// REVERT_OPERATION: Reverts a specific operation within the blockchain ledger
func REVERT_OPERATION(operationID string) error {
    err := common.ledger.RevertOperation(operationID)
    if err != nil {
        return err
    }
    LOG_MESSAGE("Revert Operation", "Operation " + operationID + " reverted successfully")
    return nil
}

// GET_CONTRACT_ID: Retrieves contract ID from the ledger for tracking and auditing
func GET_CONTRACT_ID(contractName string) (string, error) {
    contractID, err := common.ledger.GetContractID(contractName)
    if err != nil {
        return "", err
    }
    return contractID, nil
}
