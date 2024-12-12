package consensus_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
)

const (
    RecalibrationCheckInterval = 10 * time.Second // Interval for recalibration checks
    RecalibrationKey           = "recalibration_key" // Encryption key for recalibration logs
)

// ConsensusRecalibrationMaintenanceAutomation automates recalibration of Synnergy Consensus parameters
type ConsensusRecalibrationMaintenanceAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger to store recalibration-related data
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    recalibrationState recalibrationMetrics        // Current recalibration metrics
}

// recalibrationMetrics stores metrics for PoH, PoS, and PoW recalibration
type recalibrationMetrics struct {
    PoHTimestampRate   float64 `json:"poh_timestamp_rate"`   // Current PoH timestamp generation rate
    PoSMinStake        float64 `json:"pos_min_stake"`        // Current PoS minimum stake requirement
    PoWDifficultyLevel int     `json:"pow_difficulty_level"` // Current PoW difficulty level
}

// NewConsensusRecalibrationMaintenanceAutomation initializes recalibration automation
func NewConsensusRecalibrationMaintenanceAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ConsensusRecalibrationMaintenanceAutomation {
    return &ConsensusRecalibrationMaintenanceAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        recalibrationState: recalibrationMetrics{
            PoHTimestampRate:   1.0,   // Default PoH timestamp generation rate
            PoSMinStake:        100.0, // Default PoS minimum stake requirement
            PoWDifficultyLevel: 5,     // Default PoW difficulty level
        },
    }
}

// StartRecalibrationMonitoring starts the recalibration process based on real-time performance data
func (automation *ConsensusRecalibrationMaintenanceAutomation) StartRecalibrationMonitoring() {
    ticker := time.NewTicker(RecalibrationCheckInterval)

    go func() {
        for range ticker.C {
            automation.adjustPoHPerformance()
            automation.adjustPoSPerformance()
            automation.adjustPoWPerformance()
            automation.logRecalibrationMetrics()
        }
    }()
}

// adjustPoHPerformance dynamically adjusts the PoH timestamp generation rate based on system performance
func (automation *ConsensusRecalibrationMaintenanceAutomation) adjustPoHPerformance() {
    performanceData := automation.getPoHPerformanceMetrics()

    // Adjust PoH timestamp generation rate based on load
    if performanceData > 1.5 {
        automation.recalibrationState.PoHTimestampRate *= 0.9 // Reduce timestamp rate if load is high
        fmt.Println("Reduced PoH timestamp generation rate due to high load.")
    } else {
        automation.recalibrationState.PoHTimestampRate *= 1.1 // Increase rate if load is low
        fmt.Println("Increased PoH timestamp generation rate due to low load.")
    }

    // Apply recalibration by adjusting PoH internal parameters
    automation.consensusSystem.PoH.SetTimestampRate(automation.recalibrationState.PoHTimestampRate)
}

// adjustPoSPerformance dynamically adjusts PoS validator selection and stake requirements
func (automation *ConsensusRecalibrationMaintenanceAutomation) adjustPoSPerformance() {
    activeValidators := automation.getPoSValidatorMetrics()

    // Adjust PoS stake requirement based on number of validators
    if activeValidators > 100 {
        automation.recalibrationState.PoSMinStake += 10.0
        fmt.Printf("Increased PoS minimum stake to %.2f due to high validator count.\n", automation.recalibrationState.PoSMinStake)
    } else {
        automation.recalibrationState.PoSMinStake -= 5.0
        fmt.Printf("Decreased PoS minimum stake to %.2f due to low validator count.\n", automation.recalibrationState.PoSMinStake)
    }

    // Apply recalibration by updating stake requirements within the PoS system
    automation.consensusSystem.PoS.SetMinimumStake(automation.recalibrationState.PoSMinStake)
}

// adjustPoWPerformance dynamically adjusts PoW difficulty level based on block finalization time
func (automation *ConsensusRecalibrationMaintenanceAutomation) adjustPoWPerformance() {
    blockTime := automation.getPoWBlockMetrics()

    // Adjust difficulty based on block finalization time
    if blockTime < 5.0 {
        automation.recalibrationState.PoWDifficultyLevel += 1 // Increase difficulty if block mining is too fast
        fmt.Printf("Increased PoW difficulty to %d.\n", automation.recalibrationState.PoWDifficultyLevel)
    } else if blockTime > 15.0 {
        automation.recalibrationState.PoWDifficultyLevel -= 1 // Decrease difficulty if block mining is too slow
        fmt.Printf("Decreased PoW difficulty to %d.\n", automation.recalibrationState.PoWDifficultyLevel)
    }

    // Apply recalibration by updating PoW difficulty level
    automation.consensusSystem.PoW.SetDifficulty(automation.recalibrationState.PoWDifficultyLevel)
}

// logRecalibrationMetrics logs recalibration data and encrypts it for security
func (automation *ConsensusRecalibrationMaintenanceAutomation) logRecalibrationMetrics() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("recalibration-log-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Recalibration Metrics",
        Status:    "Logged",
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(RecalibrationKey))
    if err != nil {
        fmt.Printf("Error encrypting recalibration log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Recalibration metrics logged and stored in the ledger.")
}

// Method to fetch PoH performance metrics from the internal system
func (automation *ConsensusRecalibrationMaintenanceAutomation) getPoHPerformanceMetrics() float64 {
    return automation.consensusSystem.PoH.MeasurePerformance() // Directly measure PoH performance
}

// Method to fetch PoS validator count metrics from the internal system
func (automation *ConsensusRecalibrationMaintenanceAutomation) getPoSValidatorMetrics() int {
    return automation.consensusSystem.PoS.GetActiveValidatorCount() // Directly retrieve PoS active validator count
}

// Method to fetch PoW block finalization time metrics from the internal system
func (automation *ConsensusRecalibrationMaintenanceAutomation) getPoWBlockMetrics() float64 {
    return automation.consensusSystem.PoW.MeasureBlockFinalizationTime() // Directly measure PoW block mining time
}
