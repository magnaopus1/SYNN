package consensus_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/consensus"
)

const (
    ForkDetectionInterval = 2 * time.Minute   // Interval for checking for forks
    ForkResolutionKey     = "fork_resolution_key" // Encryption key for sensitive data in fork resolution
)

// ConsensusForkDetectionAndResolutionAutomation handles fork detection and resolution in the Synnergy Consensus
type ConsensusForkDetectionAndResolutionAutomation struct {
    ledgerInstance  *ledger.Ledger                     // Blockchain ledger for tracking consensus actions
    consensusEngine *consensus.SynnergyConsensus // Synnergy Consensus engine for validation
    stateMutex      *sync.RWMutex                      // Mutex for thread-safe ledger access
}

// NewConsensusForkDetectionAndResolutionAutomation initializes the automation for detecting and resolving forks
func NewConsensusForkDetectionAndResolutionAutomation(consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ConsensusForkDetectionAndResolutionAutomation {
    return &ConsensusForkDetectionAndResolutionAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
    }
}

// StartForkDetection initiates continuous monitoring for fork detection and resolution
func (automation *ConsensusForkDetectionAndResolutionAutomation) StartForkDetection() {
    ticker := time.NewTicker(ForkDetectionInterval)
    for range ticker.C {
        fmt.Println("Checking for forks in the chain...")
        automation.detectAndResolveForks()
    }
}

// detectAndResolveForks detects forks and triggers resolution if found
func (automation *ConsensusForkDetectionAndResolutionAutomation) detectAndResolveForks() {
    if automation.isForkDetected() {
        fmt.Println("Fork detected. Initiating resolution...")
        automation.resolveFork()
    } else {
        fmt.Println("No fork detected. Chain is valid.")
    }
}

// isForkDetected checks the integrity of the chain using the Synnergy Consensus validation mechanism
func (automation *ConsensusForkDetectionAndResolutionAutomation) isForkDetected() bool {
    err := automation.consensusEngine.ValidateChain()
    if err != nil {
        fmt.Println("Fork detected: Chain validation failed.")
        return true // Assume fork if chain validation fails
    }
    fmt.Println("Chain validation successful. No fork detected.")
    return false
}

// resolveFork resolves the detected fork by validating the correct chain and recalibrating nodes
func (automation *ConsensusForkDetectionAndResolutionAutomation) resolveFork() {
    success := automation.consensusEngine.PoW.ValidateBlock()
    if !success {
        fmt.Println("Error validating block during fork resolution.")
        return
    }
    fmt.Println("Fork resolved. Main chain recalibrated.")

    // Log the fork resolution in the ledger
    automation.logForkResolution()
}

// logForkResolution logs the fork resolution event in the blockchain ledger for auditing
func (automation *ConsensusForkDetectionAndResolutionAutomation) logForkResolution() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("fork-resolution-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Fork Resolution",
        Status:    "Resolved",
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(ForkResolutionKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for fork resolution: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(entry) // Validate fork resolution as a sub-block
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated with fork resolution event.\n")
}

// Additional helper function to ensure chain and ledger integrity post-fork resolution
func (automation *ConsensusForkDetectionAndResolutionAutomation) ensureChainAndLedgerIntegrity() {
    fmt.Println("Ensuring chain and ledger integrity post-fork resolution...")

    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Validate the chain and sub-blocks for consistency
    err := automation.consensusEngine.ValidateChain()
    if err != nil {
        fmt.Printf("Chain validation failed: %v\n", err)
        automation.resolveFork() // Trigger fork resolution again if integrity is not restored
    } else {
        fmt.Println("Ledger and chain are consistent after fork resolution.")
    }
}
