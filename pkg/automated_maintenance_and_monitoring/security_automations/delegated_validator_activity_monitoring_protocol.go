package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    ValidatorCheckInterval     = 10 * time.Second // Interval for monitoring validator activity
    MaxValidatorRetries        = 3               // Maximum retries for validator performance check
    MaxInactivityThreshold     = 30 * time.Second // Threshold for inactivity before alerting
)

// DelegatedValidatorActivityMonitoringAutomation automates the process of monitoring delegated validator activity
type DelegatedValidatorActivityMonitoringAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger for logging validator activity
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    validatorRetryMap  map[string]int               // Counter for retrying failed validator checks
    validatorInactivity map[string]time.Time        // Timestamp map for validator inactivity
}

// NewDelegatedValidatorActivityMonitoringAutomation initializes the automation for delegated validator activity monitoring
func NewDelegatedValidatorActivityMonitoringAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DelegatedValidatorActivityMonitoringAutomation {
    return &DelegatedValidatorActivityMonitoringAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        validatorRetryMap:  make(map[string]int),
        validatorInactivity: make(map[string]time.Time),
    }
}

// StartValidatorActivityMonitoring starts the continuous loop for monitoring validator activity at intervals
func (automation *DelegatedValidatorActivityMonitoringAutomation) StartValidatorActivityMonitoring() {
    ticker := time.NewTicker(ValidatorCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorValidators()
        }
    }()
}

// monitorValidators checks the activity of delegated validators
func (automation *DelegatedValidatorActivityMonitoringAutomation) monitorValidators() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch list of active delegated validators
    validatorList := automation.consensusSystem.GetActiveDelegatedValidators()

    if len(validatorList) > 0 {
        for _, validator := range validatorList {
            fmt.Printf("Checking activity for validator %s.\n", validator.ID)
            automation.checkValidatorActivity(validator)
        }
    } else {
        fmt.Println("No active delegated validators found at this time.")
    }
}

// checkValidatorActivity checks if the validator is active and performing as expected
func (automation *DelegatedValidatorActivityMonitoringAutomation) checkValidatorActivity(validator common.Validator) {
    // Encrypt validator data before activity check
    encryptedValidatorData := automation.encryptValidatorData(validator)

    // Trigger activity check through the Synnergy Consensus system
    activitySuccess := automation.consensusSystem.CheckValidatorActivity(validator, encryptedValidatorData)

    if activitySuccess {
        fmt.Printf("Validator %s is active and performing as expected.\n", validator.ID)
        automation.resetValidatorInactivity(validator.ID)
        automation.logValidatorActivity(validator)
        automation.resetValidatorRetry(validator.ID)
    } else {
        fmt.Printf("Validator %s did not respond. Retrying...\n", validator.ID)
        automation.retryValidatorCheck(validator)
    }
}

// retryValidatorCheck retries the validator activity check if failed
func (automation *DelegatedValidatorActivityMonitoringAutomation) retryValidatorCheck(validator common.Validator) {
    automation.validatorRetryMap[validator.ID]++
    if automation.validatorRetryMap[validator.ID] < MaxValidatorRetries {
        automation.checkValidatorActivity(validator)
    } else {
        fmt.Printf("Max retries reached for validator %s. Logging inactivity and alerting.\n", validator.ID)
        automation.logValidatorInactivity(validator)
    }
}

// resetValidatorRetry resets the retry count for validator activity check
func (automation *DelegatedValidatorActivityMonitoringAutomation) resetValidatorRetry(validatorID string) {
    automation.validatorRetryMap[validatorID] = 0
}

// logValidatorActivity logs the successful activity of a validator into the ledger
func (automation *DelegatedValidatorActivityMonitoringAutomation) logValidatorActivity(validator common.Validator) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("validator-activity-%s", validator.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Validator Activity",
        Status:    "Active",
        Details:   fmt.Sprintf("Validator %s is actively performing duties.", validator.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with validator activity for validator %s.\n", validator.ID)
}

// logValidatorInactivity logs the inactivity of a validator into the ledger and takes action
func (automation *DelegatedValidatorActivityMonitoringAutomation) logValidatorInactivity(validator common.Validator) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Log the inactivity event
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("validator-inactivity-%s", validator.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Validator Inactivity",
        Status:    "Inactive",
        Details:   fmt.Sprintf("Validator %s has been inactive for longer than the threshold.", validator.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with validator inactivity for validator %s.\n", validator.ID)

    // Trigger inactivity alerts or actions
    automation.alertInactivityAction(validator)
}

// alertInactivityAction triggers alerts or performs actions due to validator inactivity
func (automation *DelegatedValidatorActivityMonitoringAutomation) alertInactivityAction(validator common.Validator) {
    // Logic for alerting or taking actions on inactive validators
    fmt.Printf("Alerting inactivity for validator %s. Further action may be required.\n", validator.ID)
    // This could involve sending alerts, reassigning validator duties, or flagging the validator for further review.
}

// resetValidatorInactivity resets the inactivity tracking for a validator
func (automation *DelegatedValidatorActivityMonitoringAutomation) resetValidatorInactivity(validatorID string) {
    automation.validatorInactivity[validatorID] = time.Now()
}

// encryptValidatorData encrypts the validator data before performing activity checks
func (automation *DelegatedValidatorActivityMonitoringAutomation) encryptValidatorData(validator common.Validator) common.Validator {
    encryptedData, err := encryption.EncryptData(validator.Data)
    if err != nil {
        fmt.Println("Error encrypting validator data:", err)
        return validator
    }

    validator.EncryptedData = encryptedData
    fmt.Println("Validator data successfully encrypted.")
    return validator
}
