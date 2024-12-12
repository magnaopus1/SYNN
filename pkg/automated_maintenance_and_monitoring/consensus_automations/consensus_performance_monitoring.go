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
    PerformanceCheckInterval = 5 * time.Second  // Interval for checking performance metrics
    PerformanceKey           = "performance_monitor_key" // Encryption key for performance logs
)

// ConsensusPerformanceMonitoringAutomation monitors the performance of Synnergy Consensus (PoH, PoS, and PoW)
type ConsensusPerformanceMonitoringAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger to store performance-related data
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    startTime         time.Time                    // Start time for measuring performance
    transactionCount  int                          // Counter for the number of transactions processed
}

// NewConsensusPerformanceMonitoringAutomation initializes the performance monitoring automation for Synnergy Consensus
func NewConsensusPerformanceMonitoringAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ConsensusPerformanceMonitoringAutomation {
    return &ConsensusPerformanceMonitoringAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        startTime:        time.Now(),
        transactionCount: 0,
    }
}

// StartPerformanceMonitoring begins monitoring the performance of PoH, PoS, and PoW
func (automation *ConsensusPerformanceMonitoringAutomation) StartPerformanceMonitoring() {
    ticker := time.NewTicker(PerformanceCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorPoHPerformance()
            automation.monitorPoSPerformance()
            automation.monitorPoWPerformance()
            automation.trackTransactionThroughput()
            automation.trackBlockFinalizationTime()
            automation.trackValidatorPerformance()
            automation.logPerformanceMetrics()
        }
    }()
}

// monitorPoHPerformance measures the performance of PoH proof generation under load
func (automation *ConsensusPerformanceMonitoringAutomation) monitorPoHPerformance() {
    // Measure performance of PoH proof generation from SynnergyConsensus
    performance := automation.consensusSystem.PoH.MeasurePerformance()
    fmt.Printf("PoH performance metric: %.2f.\n", performance)
}

// monitorPoSPerformance checks validator efficiency and total stake in PoS
func (automation *ConsensusPerformanceMonitoringAutomation) monitorPoSPerformance() {
    // Measure performance and total stake from PoS validators
    totalStake := automation.consensusSystem.PoS.GetTotalStake()
    validatorCount := automation.consensusSystem.PoS.GetActiveValidatorCount()
    fmt.Printf("PoS Performance: Total Stake: %.2f, Active Validators: %d.\n", totalStake, validatorCount)
}

// monitorPoWPerformance measures block mining time and difficulty adjustment for PoW
func (automation *ConsensusPerformanceMonitoringAutomation) monitorPoWPerformance() {
    // Measure performance of PoW block mining
    miningTime := automation.consensusSystem.PoW.MeasureBlockMiningTime()
    difficulty := automation.consensusSystem.PoW.GetDifficulty()
    fmt.Printf("PoW Performance: Block Mining Time: %.2f seconds, Difficulty: %d.\n", miningTime, difficulty)
}

// trackTransactionThroughput measures the rate at which transactions are processed in the system
func (automation *ConsensusPerformanceMonitoringAutomation) trackTransactionThroughput() {
    // Increment transaction count from the consensus system
    processedTransactions := automation.consensusSystem.LedgerInstance.GetTransactionCount()
    automation.transactionCount += processedTransactions

    // Calculate throughput
    elapsed := time.Since(automation.startTime).Seconds()
    throughput := float64(automation.transactionCount) / elapsed
    fmt.Printf("Transaction Throughput: %.2f transactions per second.\n", throughput)
}

// trackBlockFinalizationTime measures the average time to finalize a block
func (automation *ConsensusPerformanceMonitoringAutomation) trackBlockFinalizationTime() {
    // Measure block finalization time from the consensus system
    blockTime := automation.consensusSystem.MeasureBlockFinalizationTime()
    fmt.Printf("Block finalization time: %.2f seconds.\n", blockTime)
}

// trackValidatorPerformance logs the number of sub-blocks each validator has validated
func (automation *ConsensusPerformanceMonitoringAutomation) trackValidatorPerformance() {
    // Get validator performance from the consensus system
    validatorPerformance := automation.consensusSystem.PoS.GetValidatorPerformance()

    for validator, subBlockCount := range validatorPerformance {
        fmt.Printf("Validator %s validated %d sub-blocks.\n", validator, subBlockCount)
    }
}

// logPerformanceMetrics logs the performance data in the ledger and encrypts the logs
func (automation *ConsensusPerformanceMonitoringAutomation) logPerformanceMetrics() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("performance-log-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Performance Metrics",
        Status:    "Logged",
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(PerformanceKey))
    if err != nil {
        fmt.Printf("Error encrypting performance log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Performance metrics logged and stored in the ledger.")
}

// ensureConsensusConsistency checks the overall consensus chain for consistency
func (automation *ConsensusPerformanceMonitoringAutomation) ensureConsensusConsistency() {
    valid := automation.consensusSystem.ValidateConsensusChain()
    if !valid {
        fmt.Println("Consensus validation failed during performance monitoring.")
        return
    }

    fmt.Println("Consensus chain is valid and consistent.")
}
