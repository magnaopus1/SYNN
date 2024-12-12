package security_automations

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
    HashCheckInterval         = 10 * time.Second // Interval for checking the security of hash functions
    MaxHashFailures           = 5                // Maximum allowable hash failures before triggering protocol
)

// HashFunctionSecurityAutomation automates the process of ensuring hash function security
type HashFunctionSecurityAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging hash security events
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    hashFailureCount  map[string]int               // Counter for tracking hash failures
}

// NewHashFunctionSecurityAutomation initializes the automation for hash function security
func NewHashFunctionSecurityAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *HashFunctionSecurityAutomation {
    return &HashFunctionSecurityAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        hashFailureCount: make(map[string]int),
    }
}

// StartHashFunctionSecurityMonitoring starts continuous monitoring of hash function security
func (automation *HashFunctionSecurityAutomation) StartHashFunctionSecurityMonitoring() {
    ticker := time.NewTicker(HashCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkHashFunctionIntegrity()
        }
    }()
}

// checkHashFunctionIntegrity verifies that all hash functions are operating securely
func (automation *HashFunctionSecurityAutomation) checkHashFunctionIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    hashStatus := automation.consensusSystem.CheckHashFunctions()

    for id, status := range hashStatus {
        if !status {
            automation.hashFailureCount[id]++
            fmt.Printf("Hash function failure detected for ID %s. Failure count: %d\n", id, automation.hashFailureCount[id])
            automation.handleHashFailure(id)
        } else {
            fmt.Printf("Hash function for ID %s is secure.\n", id)
            automation.resetHashFailureCount(id)
        }
    }
}

// handleHashFailure processes hash function failures
func (automation *HashFunctionSecurityAutomation) handleHashFailure(hashID string) {
    if automation.hashFailureCount[hashID] >= MaxHashFailures {
        fmt.Printf("Max hash failure threshold reached for ID %s. Triggering protocol.\n", hashID)
        automation.triggerHashSecurityProtocol(hashID)
    }
}

// triggerHashSecurityProtocol triggers a response when hash failures exceed the allowable limit
func (automation *HashFunctionSecurityAutomation) triggerHashSecurityProtocol(hashID string) {
    success := automation.consensusSystem.ActivateHashFailureProtocol(hashID)

    if success {
        fmt.Printf("Hash function security protocol successfully activated for hash ID %s.\n", hashID)
        automation.logHashSecurityEvent(hashID, "Security Protocol Activated")
    } else {
        fmt.Printf("Error activating hash function security protocol for hash ID %s.\n", hashID)
        automation.logHashSecurityEvent(hashID, "Protocol Activation Failed")
    }
}

// resetHashFailureCount resets the hash failure count for a given hash function
func (automation *HashFunctionSecurityAutomation) resetHashFailureCount(hashID string) {
    automation.hashFailureCount[hashID] = 0
}

// logHashSecurityEvent logs hash function security events to the ledger
func (automation *HashFunctionSecurityAutomation) logHashSecurityEvent(hashID string, status string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("hash-security-event-%s", hashID),
        Timestamp: time.Now().Unix(),
        Type:      "Hash Function Security",
        Status:    status,
        Details:   fmt.Sprintf("Hash function security event for hash ID %s with status: %s", hashID, status),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with hash security event for hash ID %s.\n", hashID)
}

// ensureHashFunctionIntegrity checks the integrity of hash function protocols
func (automation *HashFunctionSecurityAutomation) ensureHashFunctionIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateHashFunctionIntegrity()
    if !integrityValid {
        fmt.Println("Hash function integrity breach detected. Re-checking hashes.")
        automation.checkHashFunctionIntegrity()
    } else {
        fmt.Println("Hash functions are secure.")
    }
}
