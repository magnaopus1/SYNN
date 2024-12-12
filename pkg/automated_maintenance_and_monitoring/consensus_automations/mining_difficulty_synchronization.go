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
    DifficultyCheckInterval    = 10 * time.Minute  // Interval for difficulty synchronization checks
    MinDifficultyEncryptionKey = "min_difficulty_key"  // Encryption key for logging difficulty changes
    DefaultMiningDifficulty    = 1                 // Default mining difficulty level (adjusted)
    DifficultyAdjustmentRate   = 0.05              // Rate for dynamic difficulty adjustments
    OverridePriority           = 100               // Priority for difficulty override in hierarchy
)

// GlobalDifficultyRegistry is a centralized mechanism that manages the global mining difficulty
type GlobalDifficultyRegistry struct {
    CurrentDifficulty float64          // Global mining difficulty value
    stateMutex        sync.RWMutex     // Mutex for thread-safe access to the registry
    ledgerInstance    *ledger.Ledger   // Ledger for logging difficulty changes
    changeRequests    chan DifficultyChange // Channel for registering difficulty adjustment requests
}

// DifficultyChange holds data for requested changes to the mining difficulty
type DifficultyChange struct {
    RequestedBy   string  // The process requesting the difficulty change
    NewDifficulty float64 // The new difficulty being requested
    Override      bool    // Whether the change should override the current difficulty
    Priority      int     // The priority of the change request
}

// NewGlobalDifficultyRegistry initializes the global difficulty registry with default values
func NewGlobalDifficultyRegistry(ledgerInstance *ledger.Ledger) *GlobalDifficultyRegistry {
    return &GlobalDifficultyRegistry{
        CurrentDifficulty: DefaultMiningDifficulty,
        ledgerInstance:    ledgerInstance,
        changeRequests:    make(chan DifficultyChange, 10), // Allow up to 10 difficulty change requests concurrently
    }
}

// ApplyDifficultyChange applies a change to the global mining difficulty based on the priority and conditions of the request
func (registry *GlobalDifficultyRegistry) ApplyDifficultyChange(change DifficultyChange) {
    registry.stateMutex.Lock()
    defer registry.stateMutex.Unlock()

    if change.Override || change.Priority >= OverridePriority {
        fmt.Printf("Applying override difficulty change to %.2f requested by %s.\n", change.NewDifficulty, change.RequestedBy)
        registry.CurrentDifficulty = change.NewDifficulty
    } else if change.NewDifficulty > registry.CurrentDifficulty {
        fmt.Printf("Increasing difficulty to %.2f as requested by %s.\n", change.NewDifficulty, change.RequestedBy)
        registry.CurrentDifficulty = change.NewDifficulty
    } else {
        fmt.Printf("Difficulty change request by %s denied (current difficulty is higher).\n", change.RequestedBy)
    }

    registry.logDifficultyChange(change)
}

// logDifficultyChange logs the difficulty adjustment in the ledger
func (registry *GlobalDifficultyRegistry) logDifficultyChange(change DifficultyChange) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("difficulty-change-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Difficulty Change",
        Status:    "Applied",
        Details:   fmt.Sprintf("Difficulty changed to %.2f by %s (Priority: %d)", change.NewDifficulty, change.RequestedBy, change.Priority),
    }

    // Encrypt and log the entry
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(MinDifficultyEncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting difficulty change log: %v\n", err)
        return
    }

    registry.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Difficulty change event logged.")
}

// SynchronizeMiningDifficultyAutomation ensures that mining difficulty is synchronized globally
type SynchronizeMiningDifficultyAutomation struct {
    registry       *GlobalDifficultyRegistry            // Reference to the global difficulty registry
    consensusEngine *synnergy_consensus.SynnergyConsensus // Consensus engine for accessing and adjusting difficulty
    stateMutex     *sync.RWMutex                        // Mutex for thread-safe access
}

// NewSynchronizeMiningDifficultyAutomation initializes the automation for synchronizing difficulty
func NewSynchronizeMiningDifficultyAutomation(registry *GlobalDifficultyRegistry, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *SynchronizeMiningDifficultyAutomation {
    return &SynchronizeMiningDifficultyAutomation{
        registry:       registry,
        consensusEngine: consensusEngine,
        stateMutex:     stateMutex,
    }
}

// StartDifficultySynchronization begins monitoring and synchronizing difficulty requirements across the network
func (automation *SynchronizeMiningDifficultyAutomation) StartDifficultySynchronization() {
    ticker := time.NewTicker(DifficultyCheckInterval)
    go func() {
        for range ticker.C {
            fmt.Println("Synchronizing and adjusting mining difficulty globally...")
            automation.synchronizeDifficultyRequirements()
        }
    }()
}

// synchronizeDifficultyRequirements ensures that the PoW mining difficulty is synchronized globally
func (automation *SynchronizeMiningDifficultyAutomation) synchronizeDifficultyRequirements() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the current mining difficulty from the global registry
    currentDifficulty := automation.registry.CurrentDifficulty

    // Fetch the current PoW difficulty
    powCurrentDifficulty := automation.consensusEngine.GetPoWMiningDifficulty()

    // If PoW difficulty differs from the global mining difficulty, update it
    if powCurrentDifficulty != currentDifficulty {
        automation.consensusEngine.SetPoWMiningDifficulty(currentDifficulty)
        fmt.Printf("Updated PoW mining difficulty to %.2f to match global setting.\n", currentDifficulty)
    }

    // Reconcile any changes from other processes
    automation.reconcileDifficultyChanges(powCurrentDifficulty)
}

// reconcileDifficultyChanges ensures changes from other functions are synergized with the global registry
func (automation *SynchronizeMiningDifficultyAutomation) reconcileDifficultyChanges(powDifficulty float64) {
    if powDifficulty > automation.registry.CurrentDifficulty {
        fmt.Printf("Reconciling global mining difficulty to the highest detected difficulty: %.2f\n", powDifficulty)
        automation.registry.ApplyDifficultyChange(DifficultyChange{
            RequestedBy:   "Reconciliation Process",
            NewDifficulty: powDifficulty,
            Override:      false,
            Priority:      90, // Reconciliation priority
        })
    }
}

// RequestDifficultyChange allows other processes to request changes to the global mining difficulty
func (automation *SynchronizeMiningDifficultyAutomation) RequestDifficultyChange(requestedBy string, newDifficulty float64, override bool, priority int) {
    automation.registry.ApplyDifficultyChange(DifficultyChange{
        RequestedBy:   requestedBy,
        NewDifficulty: newDifficulty,
        Override:      override,
        Priority:      priority,
    })
}

// MonitorOtherDifficultyFunctions monitors other difficulty-related functions for conflicts and ensures synchronization overrides them
func (automation *SynchronizeMiningDifficultyAutomation) MonitorOtherDifficultyFunctions() {
    // Continuously monitor for any difficulty changes or adjustments made by other components
    ticker := time.NewTicker(DifficultyCheckInterval)
    for range ticker.C {
        fmt.Println("Monitoring for conflicting difficulty adjustments...")
        automation.synchronizeDifficultyRequirements()
    }
}
