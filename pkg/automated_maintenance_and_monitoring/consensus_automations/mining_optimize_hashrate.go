package consensus_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    HashrateCheckInterval    = 10 * time.Second // Interval for checking validator hashrate
    MinimumHashrateThreshold = 0.75             // Threshold for flagging slow hashrate
    HashrateOptimizationKey  = "hashrate_optimization_key" // Encryption key for hashrate logs
)

// MiningOptimizeHashrateAutomation manages dynamic optimization of validator hashrate and mining difficulty
type MiningOptimizeHashrateAutomation struct {
    ledgerInstance  *ledger.Ledger                    // Blockchain ledger for logging optimization changes
    consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for monitoring and optimization
    stateMutex      *sync.RWMutex                     // Mutex for thread-safe ledger access
    currentDifficulty float64                         // Current PoW mining difficulty
}

// NewMiningOptimizeHashrateAutomation initializes the automation for hashrate optimization
func NewMiningOptimizeHashrateAutomation(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *MiningOptimizeHashrateAutomation {
    return &MiningOptimizeHashrateAutomation{
        ledgerInstance:   ledgerInstance,
        consensusEngine:  consensusEngine,
        stateMutex:       stateMutex,
        currentDifficulty: DefaultMiningDifficulty,
    }
}

// StartHashrateMonitoring starts the continuous monitoring of validator hashrate and optimizes difficulty as needed
func (automation *MiningOptimizeHashrateAutomation) StartHashrateMonitoring() {
    ticker := time.NewTicker(HashrateCheckInterval)
    go func() {
        for range ticker.C {
            fmt.Println("Monitoring validator hashrate and adjusting mining difficulty...")
            automation.monitorAndOptimizeHashrate()
        }
    }()
}

// monitorAndOptimizeHashrate checks the hashrate of validators and dynamically adjusts difficulty or optimizes their performance
func (automation *MiningOptimizeHashrateAutomation) monitorAndOptimizeHashrate() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch all validator hashrates from the consensus engine
    hashrates := automation.consensusEngine.GetValidatorHashrates()

    for validator, hashrate := range hashrates {
        if hashrate < MinimumHashrateThreshold {
            fmt.Printf("Validator %s is struggling with a hashrate of %.2f. Optimizing...\n", validator, hashrate)
            automation.optimizeHashrate(validator)
        }
    }

    // Adjust overall mining difficulty based on network conditions and validator performance
    automation.adjustMiningDifficulty(hashrates)
}

// optimizeHashrate redistributes computational load or adjusts hashrate optimization for a specific validator
func (automation *MiningOptimizeHashrateAutomation) optimizeHashrate(validator string) {
    // Dynamically distribute mining tasks to assist the struggling validator
    success := automation.consensusEngine.DistributeMiningWorkload(validator)

    if success {
        fmt.Printf("Successfully optimized hashrate for validator %s.\n", validator)
        automation.logHashrateOptimization(validator)
    } else {
        fmt.Printf("Failed to optimize hashrate for validator %s.\n", validator)
    }
}

// adjustMiningDifficulty dynamically adjusts PoW mining difficulty based on average network hashrate
func (automation *MiningOptimizeHashrateAutomation) adjustMiningDifficulty(hashrates map[string]float64) {
    totalHashrate := float64(0)
    numValidators := len(hashrates)

    for _, hashrate := range hashrates {
        totalHashrate += hashrate
    }

    averageHashrate := totalHashrate / float64(numValidators)

    // Dynamically adjust difficulty if the average hashrate is below the threshold
    if averageHashrate < MinimumHashrateThreshold {
        automation.currentDifficulty *= (1 - DifficultyAdjustmentRate)
        fmt.Printf("Decreasing mining difficulty to %.2f due to low average hashrate.\n", automation.currentDifficulty)
    } else {
        automation.currentDifficulty *= (1 + DifficultyAdjustmentRate)
        fmt.Printf("Increasing mining difficulty to %.2f due to sufficient hashrate.\n", automation.currentDifficulty)
    }

    // Apply the new difficulty setting to the consensus engine
    automation.consensusEngine.SetMiningDifficulty(automation.currentDifficulty)

    // Log the difficulty adjustment
    automation.logDifficultyAdjustment()
}

// logHashrateOptimization logs the successful optimization of a validator's hashrate in the ledger
func (automation *MiningOptimizeHashrateAutomation) logHashrateOptimization(validator string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("hashrate-optimization-%s-%d", validator, time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Hashrate Optimization",
        Status:    "Optimized",
        Details:   fmt.Sprintf("Validator %s hashrate optimized.", validator),
    }

    // Encrypt the log entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(HashrateOptimizationKey))
    if err != nil {
        fmt.Printf("Error encrypting hashrate optimization log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Hashrate optimization log stored in the ledger.")
}

// logDifficultyAdjustment logs the dynamic adjustment of mining difficulty in the ledger
func (automation *MiningOptimizeHashrateAutomation) logDifficultyAdjustment() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("difficulty-adjustment-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Mining Difficulty Adjustment",
        Status:    "Adjusted",
        Details:   fmt.Sprintf("Mining difficulty adjusted to %.2f.", automation.currentDifficulty),
    }

    // Encrypt the log entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(HashrateOptimizationKey))
    if err != nil {
        fmt.Printf("Error encrypting mining difficulty log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Mining difficulty adjustment log stored in the ledger.")
}

// TriggerManualOptimization allows for manually triggering hashrate optimization for a specific validator
func (automation *MiningOptimizeHashrateAutomation) TriggerManualOptimization(validator string) {
    fmt.Printf("Manually triggering optimization for validator %s.\n", validator)
    automation.optimizeHashrate(validator)
}
