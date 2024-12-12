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
    ValidatorIncentiveCheckInterval    = 1500 * time.Millisecond // Interval for checking validator eco-friendliness
    MaxAllowedEnergyUsage              = 500                     // Maximum allowed energy usage in watts for eco-friendly incentive
    SubBlocksPerBlock                  = 1000                    // Number of sub-blocks in a block
    IncentiveBonus                     = 10                      // Incentive bonus for eco-friendly validators
)

// EcoFriendlyValidatorIncentiveAutomation automates the incentivization of eco-friendly validators based on energy efficiency
type EcoFriendlyValidatorIncentiveAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store incentive-related logs
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    incentiveTriggerCount  int                          // Counter for eco-friendly validator incentives
}

// NewEcoFriendlyValidatorIncentiveAutomation initializes the automation for incentivizing eco-friendly validators
func NewEcoFriendlyValidatorIncentiveAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *EcoFriendlyValidatorIncentiveAutomation {
    return &EcoFriendlyValidatorIncentiveAutomation{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        incentiveTriggerCount: 0,
    }
}

// StartEcoFriendlyIncentiveAutomation starts the continuous loop for monitoring validator energy efficiency
func (automation *EcoFriendlyValidatorIncentiveAutomation) StartEcoFriendlyIncentiveAutomation() {
    ticker := time.NewTicker(ValidatorIncentiveCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndIncentivizeValidators()
        }
    }()
}

// monitorAndIncentivizeValidators checks validator energy usage and triggers incentives for eco-friendly behavior
func (automation *EcoFriendlyValidatorIncentiveAutomation) monitorAndIncentivizeValidators() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch current validator energy usage data from the consensus system
    validators := automation.consensusSystem.CheckValidatorEnergyUsage()

    for _, validator := range validators {
        if validator.EnergyUsage <= MaxAllowedEnergyUsage {
            fmt.Printf("Validator %s is eco-friendly with energy usage %d watts. Triggering incentive.\n", validator.Address, validator.EnergyUsage)
            automation.triggerIncentive(validator)
        } else {
            fmt.Printf("Validator %s exceeds eco-friendly energy usage (%d watts).\n", validator.Address, validator.EnergyUsage)
        }
    }

    automation.incentiveTriggerCount++
    fmt.Printf("Eco-friendly validator incentive cycle #%d executed.\n", automation.incentiveTriggerCount)

    if automation.incentiveTriggerCount%SubBlocksPerBlock == 0 {
        automation.finalizeIncentiveCycle()
    }
}

// triggerIncentive grants incentives to validators with eco-friendly energy usage
func (automation *EcoFriendlyValidatorIncentiveAutomation) triggerIncentive(validator common.Validator) {
    // Encrypt the validator's incentive data before applying the incentive
    encryptedValidator := automation.AddEncryptionToValidatorData(validator)

    // Provide the incentive through Synnergy Consensus
    incentiveSuccess := automation.consensusSystem.GrantIncentive(encryptedValidator, IncentiveBonus)
    if incentiveSuccess {
        fmt.Printf("Incentive of %d tokens granted to eco-friendly validator %s.\n", IncentiveBonus, validator.Address)
        automation.logIncentiveAction(validator)
    } else {
        fmt.Printf("Error granting incentive to validator %s.\n", validator.Address)
    }
}

// finalizeIncentiveCycle finalizes the incentive cycle and logs the result in the ledger
func (automation *EcoFriendlyValidatorIncentiveAutomation) finalizeIncentiveCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeIncentiveCycle()
    if success {
        fmt.Println("Eco-friendly validator incentive cycle finalized successfully.")
        automation.logIncentiveCycleFinalization()
    } else {
        fmt.Println("Error finalizing eco-friendly validator incentive cycle.")
    }
}

// logIncentiveAction logs each validator incentive action into the ledger for traceability
func (automation *EcoFriendlyValidatorIncentiveAutomation) logIncentiveAction(validator common.Validator) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("eco-friendly-incentive-%s", validator.Address),
        Timestamp: time.Now().Unix(),
        Type:      "Eco-Friendly Incentive",
        Status:    "Incentive Granted",
        Details:   fmt.Sprintf("Incentive of %d tokens granted.", IncentiveBonus),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with eco-friendly incentive action for ValidatorID %s.\n", validator.Address)
}

// logIncentiveCycleFinalization logs the finalization of an eco-friendly validator incentive cycle into the ledger
func (automation *EcoFriendlyValidatorIncentiveAutomation) logIncentiveCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("eco-friendly-incentive-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Incentive Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with eco-friendly validator incentive cycle finalization.")
}

// AddEncryptionToValidatorData encrypts validator data before applying the incentive
func (automation *EcoFriendlyValidatorIncentiveAutomation) AddEncryptionToValidatorData(validator common.Validator) common.Validator {
    encryptedData, err := encryption.EncryptData(validator)
    if err != nil {
        fmt.Println("Error encrypting validator data:", err)
        return validator
    }
    validator.EncryptedData = encryptedData
    fmt.Println("Validator data successfully encrypted.")
    return validator
}

// ensureValidatorEnergyIntegrity checks the integrity of validator energy usage data and triggers incentives if necessary
func (automation *EcoFriendlyValidatorIncentiveAutomation) ensureValidatorEnergyIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateEnergyUsageIntegrity()
    if !integrityValid {
        fmt.Println("Validator energy usage integrity breach detected. Triggering eco-friendly incentives.")
        automation.monitorAndIncentivizeValidators()
    } else {
        fmt.Println("Validator energy usage integrity is valid.")
    }
}
