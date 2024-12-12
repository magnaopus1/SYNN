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
    SecurityCheckInterval     = 10 * time.Second // Interval for checking security threats
    SecurityEnforcementKey    = "security_enforcement_key" // Encryption key for security logs
)

// ConsensusSecurityMonitoringAutomation automates security checks for Synnergy Consensus
type ConsensusSecurityMonitoringAutomation struct {
    consensusSystem *consensus.SynnergyConsensus // Reference to the Synnergy Consensus system
    ledgerInstance  *ledger.Ledger               // Ledger instance for storing security logs
    stateMutex      *sync.RWMutex                // Mutex for thread-safe access
}

// NewConsensusSecurityMonitoringAutomation initializes the security monitoring automation
func NewConsensusSecurityMonitoringAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ConsensusSecurityMonitoringAutomation {
    return &ConsensusSecurityMonitoringAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
    }
}

// StartSecurityMonitoring initiates the continuous monitoring of security threats
func (automation *ConsensusSecurityMonitoringAutomation) StartSecurityMonitoring() {
    ticker := time.NewTicker(SecurityCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorPoHActivity()
            automation.monitorPoSValidatorActivity()
            automation.monitorPoWBlockValidation()
            automation.logSecurityMetrics()
        }
    }()
}

// monitorPoHActivity checks for irregularities or malicious activity in PoH transactions
func (automation *ConsensusSecurityMonitoringAutomation) monitorPoHActivity() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    // Simulating the retrieval of PoH activity logs
    proofs := automation.consensusSystem.PoH.GenerateMultipleProofs(10)
    
    for _, proof := range proofs {
        if !automation.consensusSystem.PoH.ValidatePoHProof(proof, "monitor") {
            fmt.Printf("Security Alert: Invalid PoH proof detected: %s\n", proof.Hash)
            automation.flagMaliciousPoH(proof.Hash)
        }
    }
}

// monitorPoSValidatorActivity monitors PoS validator activity to detect potential Sybil attacks or malicious behavior
func (automation *ConsensusSecurityMonitoringAutomation) monitorPoSValidatorActivity() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    activeValidators := automation.consensusSystem.PoS.State.Validators

    // Check for unusual staking patterns or Sybil attacks
    for _, validator := range activeValidators {
        if validator.Stake < automation.consensusSystem.PoS.State.TotalStake*0.01 { // Example: suspiciously low stake
            fmt.Printf("Security Alert: Suspicious PoS validator detected: %s\n", validator.Address)
            automation.blockMaliciousValidator(validator.Address)
        }
    }
}

// monitorPoWBlockValidation checks for irregularities in PoW block mining, preventing potential 51% attacks
func (automation *ConsensusSecurityMonitoringAutomation) monitorPoWBlockValidation() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    // Simulating the retrieval of the last block mined
    lastBlock := automation.consensusSystem.LedgerInstance.GetLastBlock()

    if !automation.consensusSystem.PoW.ValidateBlock(&lastBlock) {
        fmt.Printf("Security Alert: Invalid PoW block detected: %s\n", lastBlock.Hash)
        automation.flagMaliciousPoWBlock(lastBlock.Hash)
    }
}

// flagMaliciousPoH logs and handles malicious PoH activity
func (automation *ConsensusSecurityMonitoringAutomation) flagMaliciousPoH(proofHash string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Log the malicious activity
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("security-log-poh-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "PoH Security Alert",
        Status:    fmt.Sprintf("Invalid PoH proof detected: %s", proofHash),
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(SecurityEnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting security log for PoH: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("PoH security log stored for proof: %s\n", proofHash)
}

// blockMaliciousValidator blocks a malicious PoS validator from further participation
func (automation *ConsensusSecurityMonitoringAutomation) blockMaliciousValidator(validatorAddress string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Disqualify and freeze the malicious validator's stake
    automation.consensusSystem.PoS.RemoveStake(validatorAddress, automation.consensusSystem.PoS.GetValidatorStake(validatorAddress))
    fmt.Printf("Validator %s has been blocked and removed from validation.\n", validatorAddress)

    // Log the action in the ledger
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("security-log-pos-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "PoS Security Alert",
        Status:    fmt.Sprintf("Validator %s blocked due to suspicious activity.", validatorAddress),
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(SecurityEnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting security log for PoS: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("PoS security log stored for validator: %s\n", validatorAddress)
}

// flagMaliciousPoWBlock logs and handles malicious PoW activity
func (automation *ConsensusSecurityMonitoringAutomation) flagMaliciousPoWBlock(blockHash string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Log the malicious block
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("security-log-pow-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "PoW Security Alert",
        Status:    fmt.Sprintf("Invalid PoW block detected: %s", blockHash),
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(SecurityEnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting security log for PoW: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("PoW security log stored for block: %s\n", blockHash)
}

// logSecurityMetrics logs security metrics and stores them in the ledger
func (automation *ConsensusSecurityMonitoringAutomation) logSecurityMetrics() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("security-metrics-log-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Security Metrics",
        Status:    "Security metrics recorded",
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(SecurityEnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting security metrics: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Security metrics logged and stored in the ledger.")
}
