package automations

import (
    "fmt"
    "sort"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    EnergyConsumptionCheckInterval = 2000 * time.Millisecond // Interval for checking energy consumption ranking
    SubBlocksPerBlock              = 1000                    // Number of sub-blocks in a block
    MaxAllowedEnergyUsage          = 1000                    // Maximum allowed energy usage for ranking in watts
)

// ValidatorEnergyRank represents a structure for ranking validators based on energy consumption
type ValidatorEnergyRank struct {
    Address     string // Validator address
    EnergyUsage int    // Energy usage in watts
}

// EnergyConsumptionRankingAutomation automates the ranking of validators based on energy consumption
type EnergyConsumptionRankingAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance      *ledger.Ledger               // Ledger to store energy consumption ranking-related logs
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    rankingTriggerCount int                          // Counter for how many times rankings have been calculated
}

// NewEnergyConsumptionRankingAutomation initializes the automation for ranking validators based on energy consumption
func NewEnergyConsumptionRankingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *EnergyConsumptionRankingAutomation {
    return &EnergyConsumptionRankingAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        rankingTriggerCount: 0,
    }
}

// StartEnergyConsumptionRankingAutomation starts the continuous loop for monitoring and ranking validators based on energy consumption
func (automation *EnergyConsumptionRankingAutomation) StartEnergyConsumptionRankingAutomation() {
    ticker := time.NewTicker(EnergyConsumptionCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndRankEnergyConsumption()
        }
    }()
}

// monitorAndRankEnergyConsumption checks energy usage across validators and ranks them
func (automation *EnergyConsumptionRankingAutomation) monitorAndRankEnergyConsumption() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of validators and their energy consumption data from Synnergy Consensus
    validators := automation.consensusSystem.GetValidatorsEnergyData()

    // Rank validators by their energy consumption
    rankedValidators := automation.rankValidatorsByEnergyUsage(validators)

    fmt.Println("Validators ranked by energy consumption:")
    for i, rank := range rankedValidators {
        fmt.Printf("Rank %d: Validator %s with energy usage %d watts.\n", i+1, rank.Address, rank.EnergyUsage)
    }

    // Log the ranked validators into the ledger
    automation.logEnergyRanking(rankedValidators)

    automation.rankingTriggerCount++
    fmt.Printf("Energy consumption ranking cycle #%d executed.\n", automation.rankingTriggerCount)

    if automation.rankingTriggerCount%SubBlocksPerBlock == 0 {
        automation.finalizeRankingCycle()
    }
}

// rankValidatorsByEnergyUsage ranks the validators by their energy consumption
func (automation *EnergyConsumptionRankingAutomation) rankValidatorsByEnergyUsage(validators []common.ValidatorEnergyData) []ValidatorEnergyRank {
    var rankedValidators []ValidatorEnergyRank

    for _, validator := range validators {
        if validator.EnergyUsage <= MaxAllowedEnergyUsage {
            rankedValidators = append(rankedValidators, ValidatorEnergyRank{
                Address:     validator.Address,
                EnergyUsage: validator.EnergyUsage,
            })
        }
    }

    // Sort the validators by their energy usage in ascending order (lowest energy consumption ranks highest)
    sort.Slice(rankedValidators, func(i, j int) bool {
        return rankedValidators[i].EnergyUsage < rankedValidators[j].EnergyUsage
    })

    return rankedValidators
}

// logEnergyRanking logs the validator energy consumption rankings into the ledger
func (automation *EnergyConsumptionRankingAutomation) logEnergyRanking(rankedValidators []ValidatorEnergyRank) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("energy-consumption-ranking-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Energy Consumption Ranking",
        Status:    "Ranked",
        Details:   fmt.Sprintf("Top-ranked validator: %s with energy usage %d watts", rankedValidators[0].Address, rankedValidators[0].EnergyUsage),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with energy consumption ranking.")
}

// finalizeRankingCycle finalizes the energy consumption ranking cycle and logs it into the ledger
func (automation *EnergyConsumptionRankingAutomation) finalizeRankingCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeEnergyRankingCycle()
    if success {
        fmt.Println("Energy consumption ranking cycle finalized successfully.")
        automation.logRankingCycleFinalization()
    } else {
        fmt.Println("Error finalizing energy consumption ranking cycle.")
    }
}

// logRankingCycleFinalization logs the finalization of an energy consumption ranking cycle into the ledger
func (automation *EnergyConsumptionRankingAutomation) logRankingCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("energy-consumption-ranking-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Energy Consumption Ranking Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with energy consumption ranking cycle finalization.")
}

// AddEncryptionToValidatorData encrypts validator energy consumption data before logging
func (automation *EnergyConsumptionRankingAutomation) AddEncryptionToValidatorData(validator common.ValidatorEnergyData) common.ValidatorEnergyData {
    encryptedData, err := encryption.EncryptData(validator)
    if err != nil {
        fmt.Println("Error encrypting validator data:", err)
        return validator
    }
    validator.EncryptedData = encryptedData
    fmt.Println("Validator energy consumption data successfully encrypted.")
    return validator
}

// ensureValidatorEnergyIntegrity checks the integrity of validator energy data and ensures it is ranked properly
func (automation *EnergyConsumptionRankingAutomation) ensureValidatorEnergyIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateEnergyUsageIntegrity()
    if !integrityValid {
        fmt.Println("Validator energy usage integrity breach detected. Triggering energy consumption ranking.")
        automation.monitorAndRankEnergyConsumption()
    } else {
        fmt.Println("Validator energy usage integrity is valid.")
    }
}
