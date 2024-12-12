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
    ProtocolCheckInterval = 10 * time.Second // Interval for checking protocol compliance
    ProtocolKey           = "protocol_monitor_key" // Encryption key for protocol logs
)

// ConsensusProtocolMonitoringAutomation ensures Synnergy Consensus adheres to protocol rules
type ConsensusProtocolMonitoringAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // SynnergyConsensus struct
    ledgerInstance   *ledger.Ledger               // Ledger to store protocol-related data
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
}

// NewConsensusProtocolMonitoringAutomation initializes the protocol monitoring automation
func NewConsensusProtocolMonitoringAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ConsensusProtocolMonitoringAutomation {
    return &ConsensusProtocolMonitoringAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
    }
}

// StartProtocolMonitoring begins the continuous monitoring of protocol adherence
func (automation *ConsensusProtocolMonitoringAutomation) StartProtocolMonitoring() {
    ticker := time.NewTicker(ProtocolCheckInterval)

    go func() {
        for range ticker.C {
            automation.validatePoH()
            automation.validatePoS()
            automation.validatePoW()
            automation.ensureProtocolConsistency()
            automation.logProtocolMetrics()
        }
    }()
}

// validatePoH validates that PoH proofs follow correct timestamping and sequencing
func (automation *ConsensusProtocolMonitoringAutomation) validatePoH() {
    // Validate PoH proofs within the consensus system directly
    valid := automation.consensusSystem.PoH.ValidatePoHProofs()
    if !valid {
        fmt.Println("PoH protocol validation failed.")
        return
    }
    fmt.Println("PoH protocol validated successfully.")
}

// validatePoS ensures that PoS validators meet the stake requirements and integrity
func (automation *ConsensusProtocolMonitoringAutomation) validatePoS() {
    // Validate PoS validators directly within the consensus system
    valid := automation.consensusSystem.PoS.ValidateStakeRequirements()
    if !valid {
        fmt.Println("PoS protocol validation failed.")
        return
    }
    fmt.Println("PoS protocol validated successfully.")
}

// validatePoW ensures that PoW block difficulty and validation meet protocol rules
func (automation *ConsensusProtocolMonitoringAutomation) validatePoW() {
    // Validate PoW blocks directly within the consensus system
    valid := automation.consensusSystem.PoW.ValidateBlockDifficulty()
    if !valid {
        fmt.Println("PoW protocol validation failed.")
        return
    }
    fmt.Println("PoW protocol validated successfully.")
}

// ensureProtocolConsistency checks that the overall consensus process follows protocol sequencing
func (automation *ConsensusProtocolMonitoringAutomation) ensureProtocolConsistency() {
    valid := automation.consensusSystem.ValidateConsensusChain()
    if !valid {
        fmt.Println("Protocol sequencing validation failed.")
        return
    }
    fmt.Println("Protocol sequencing is consistent.")
}

// logProtocolMetrics logs the protocol validation data and encrypts the logs
func (automation *ConsensusProtocolMonitoringAutomation) logProtocolMetrics() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("protocol-log-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Protocol Validation",
        Status:    "Logged",
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(ProtocolKey))
    if err != nil {
        fmt.Printf("Error encrypting protocol log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Protocol metrics logged and stored in the ledger.")
}
