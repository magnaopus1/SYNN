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
    LatencyCheckInterval = 2 * time.Minute         // Interval for latency checks between PoH, PoS, and PoW
    LatencyThreshold     = 500 * time.Millisecond  // Threshold for acceptable latency between stages
    LatencyResolutionKey = "latency_resolution_key" // Encryption key for latency-related data
)

// ConsensusLatencyMonitoringAutomation monitors latency between consensus stages (PoH, PoS, PoW)
type ConsensusLatencyMonitoringAutomation struct {
    ledgerInstance  *ledger.Ledger                   // Blockchain ledger for tracking consensus actions
    consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validation
    stateMutex      *sync.RWMutex                    // Mutex for thread-safe ledger access
}

// NewConsensusLatencyMonitoringAutomation initializes the automation for tracking consensus latency
func NewConsensusLatencyMonitoringAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ConsensusLatencyMonitoringAutomation {
    return &ConsensusLatencyMonitoringAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
    }
}

// StartLatencyMonitoring initiates continuous monitoring of latency between PoH, PoS, and PoW
func (automation *ConsensusLatencyMonitoringAutomation) StartLatencyMonitoring() {
    ticker := time.NewTicker(LatencyCheckInterval)
    for range ticker.C {
        fmt.Println("Checking latency between consensus stages...")
        automation.monitorConsensusLatency()
    }
}

// monitorConsensusLatency tracks the time spent between PoH, PoS, and PoW transitions and ensures latency is within thresholds
func (automation *ConsensusLatencyMonitoringAutomation) monitorConsensusLatency() {
    pohLatency := automation.trackPoHLatency()
    posLatency := automation.trackPoSLatency()
    powLatency := automation.trackPoWLatency()

    if pohLatency > LatencyThreshold || posLatency > LatencyThreshold || powLatency > LatencyThreshold {
        fmt.Println("Consensus latency exceeded acceptable thresholds.")
        automation.triggerLatencyOptimization(pohLatency, posLatency, powLatency)
    } else {
        fmt.Println("Consensus latency is within acceptable thresholds.")
    }
}

// trackPoHLatency measures the latency of the PoH stage
func (automation *ConsensusLatencyMonitoringAutomation) trackPoHLatency() time.Duration {
    start := time.Now()
    success := automation.consensusEngine.PoH.ValidatePoHProof()
    if !success {
        fmt.Println("Error tracking PoH latency.")
        return time.Duration(0)
    }
    latency := time.Since(start)
    fmt.Printf("PoH latency: %v\n", latency)
    return latency
}

// trackPoSLatency measures the latency of the PoS stage
func (automation *ConsensusLatencyMonitoringAutomation) trackPoSLatency() time.Duration {
    start := time.Now()
    validator := automation.consensusEngine.PoS.SelectValidator()
    if validator == nil {
        fmt.Println("Error tracking PoS latency.")
        return time.Duration(0)
    }

    success := automation.consensusEngine.PoS.ValidateSubBlock(validator)
    if !success {
        fmt.Println("Error tracking PoS latency.")
        return time.Duration(0)
    }
    latency := time.Since(start)
    fmt.Printf("PoS latency: %v\n", latency)
    return latency
}

// trackPoWLatency measures the latency of the PoW stage
func (automation *ConsensusLatencyMonitoringAutomation) trackPoWLatency() time.Duration {
    start := time.Now()
    success := automation.consensusEngine.PoW.MineBlock()
    if !success {
        fmt.Println("Error tracking PoW latency.")
        return time.Duration(0)
    }
    latency := time.Since(start)
    fmt.Printf("PoW latency: %v\n", latency)
    return latency
}

// triggerLatencyOptimization optimizes consensus performance when latency exceeds thresholds
func (automation *ConsensusLatencyMonitoringAutomation) triggerLatencyOptimization(pohLatency, posLatency, powLatency time.Duration) {
    fmt.Printf("Optimizing consensus performance due to high latency: PoH: %v, PoS: %v, PoW: %v\n", pohLatency, posLatency, powLatency)

    // Re-process transactions to optimize latency
    automation.processTransactionsForLatency()

    // Finalize the block as part of performance optimization
    automation.finalizeBlockForLatency()

    // Log the latency optimization event in the ledger
    automation.logLatencyOptimization(pohLatency, posLatency, powLatency)
}

// processTransactionsForLatency handles processing transactions to reduce latency
func (automation *ConsensusLatencyMonitoringAutomation) processTransactionsForLatency() {
    success := automation.consensusEngine.ProcessTransactions()
    if !success {
        fmt.Println("Error processing transactions for latency optimization.")
        return
    }
    fmt.Println("Transactions processed successfully for latency optimization.")
}

// finalizeBlockForLatency finalizes the block to improve consensus performance and latency
func (automation *ConsensusLatencyMonitoringAutomation) finalizeBlockForLatency() {
    success := automation.consensusEngine.PoW.MineBlock()
    if !success {
        fmt.Println("Error finalizing block for latency optimization.")
        return
    }
    fmt.Println("Block finalized successfully for latency optimization.")
}

// logLatencyOptimization logs the latency optimization event in the blockchain ledger for auditing
func (automation *ConsensusLatencyMonitoringAutomation) logLatencyOptimization(pohLatency, posLatency, powLatency time.Duration) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("latency-optimization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Latency Optimization",
        Status:    "Optimized",
        Details:   fmt.Sprintf("PoH: %v, PoS: %v, PoW: %v", pohLatency, posLatency, powLatency),
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(LatencyResolutionKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for latency optimization: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated with latency optimization event.\n")
}

// ensureLatencyConsistency checks and ensures latency consistency across consensus stages
func (automation *ConsensusLatencyMonitoringAutomation) ensureLatencyConsistency() {
    fmt.Println("Ensuring latency consistency across consensus stages...")

    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Validate consensus consistency after latency optimizations
    success := automation.consensusEngine.ValidateConsensusChain()
    if !success {
        fmt.Printf("Consensus validation failed. Re-optimizing...\n")
        automation.triggerLatencyOptimization(time.Duration(0), time.Duration(0), time.Duration(0)) // Trigger optimization again if consistency fails
    } else {
        fmt.Println("Consensus performance and latency are consistent.")
    }
}
