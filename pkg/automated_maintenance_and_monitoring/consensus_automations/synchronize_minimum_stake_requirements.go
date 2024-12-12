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
    StakeCheckInterval         = 10 * time.Minute // Interval for stake synchronization checks
    MinStakeEncryptionKey      = "min_stake_key"  // Encryption key for logging stake changes
    DefaultMinStake            = 1000             // Default minimum stake
    StakeAdjustmentThreshold   = 0.75             // Threshold for stake adjustments
    ValidatorStakeSyncRate     = 0.05             // Rate for stake adjustments in synchronization
    StakeOverridePriority      = 100              // Priority for stake override in the hierarchy
)

// GlobalStakeRegistry is a centralized mechanism that manages the global minimum stake value
type GlobalStakeRegistry struct {
    MinStake       float64          // Global minimum stake value
    stateMutex     sync.RWMutex     // Mutex for thread-safe access to the registry
    ledgerInstance *ledger.Ledger   // Ledger for logging stake changes
    changeRequests chan StakeChange // Channel for registering stake adjustment requests
}

// StakeChange holds data for requested changes to the minimum stake
type StakeChange struct {
    RequestedBy   string  // The process requesting the stake change
    NewMinStake   float64 // The new minimum stake being requested
    Override      bool    // Whether the change should override the current stake
    Priority      int     // The priority of the change request
}

// NewGlobalStakeRegistry initializes the global stake registry with default values
func NewGlobalStakeRegistry(ledgerInstance *ledger.Ledger) *GlobalStakeRegistry {
    return &GlobalStakeRegistry{
        MinStake:       DefaultMinStake,
        ledgerInstance: ledgerInstance,
        changeRequests: make(chan StakeChange, 10), // Allow up to 10 stake change requests concurrently
    }
}

// ApplyStakeChange applies a change to the global minimum stake based on the priority and conditions of the request
func (registry *GlobalStakeRegistry) ApplyStakeChange(change StakeChange) {
    registry.stateMutex.Lock()
    defer registry.stateMutex.Unlock()

    if change.Override || change.Priority >= StakeOverridePriority {
        fmt.Printf("Applying override stake change to %.2f requested by %s.\n", change.NewMinStake, change.RequestedBy)
        registry.MinStake = change.NewMinStake
    } else if change.NewMinStake > registry.MinStake {
        fmt.Printf("Increasing stake to %.2f as requested by %s.\n", change.NewMinStake, change.RequestedBy)
        registry.MinStake = change.NewMinStake
    } else {
        fmt.Printf("Stake change request by %s denied (current stake is higher).\n", change.RequestedBy)
    }

    registry.logStakeChange(change)
}

// logStakeChange logs the stake adjustment in the ledger
func (registry *GlobalStakeRegistry) logStakeChange(change StakeChange) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("stake-change-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Stake Change",
        Status:    "Applied",
        Details:   fmt.Sprintf("Stake changed to %.2f by %s (Priority: %d)", change.NewMinStake, change.RequestedBy, change.Priority),
    }

    // Encrypt and log the entry
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(MinStakeEncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting stake change log: %v\n", err)
        return
    }

    registry.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Stake change event logged.")
}

// SynchronizeMinimumStakeAutomation ensures that minimum stake is synchronized globally
type SynchronizeMinimumStakeAutomation struct {
    registry       *GlobalStakeRegistry               // Reference to the global stake registry
    consensusEngine *synnergy_consensus.SynnergyConsensus // Consensus engine for accessing and adjusting stake
    stateMutex     *sync.RWMutex                      // Mutex for thread-safe access
}

// NewSynchronizeMinimumStakeAutomation initializes the automation for synchronizing stake
func NewSynchronizeMinimumStakeAutomation(registry *GlobalStakeRegistry, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *SynchronizeMinimumStakeAutomation {
    return &SynchronizeMinimumStakeAutomation{
        registry:       registry,
        consensusEngine: consensusEngine,
        stateMutex:     stateMutex,
    }
}

// StartStakeSynchronization begins monitoring and synchronizing stake requirements across PoH and PoS
func (automation *SynchronizeMinimumStakeAutomation) StartStakeSynchronization() {
    ticker := time.NewTicker(StakeCheckInterval)
    go func() {
        for range ticker.C {
            fmt.Println("Synchronizing and adjusting minimum stake globally...")
            automation.synchronizeStakeRequirements()
        }
    }()
}

// synchronizeStakeRequirements ensures that the PoH and PoS minimum stakes are synchronized globally
func (automation *SynchronizeMinimumStakeAutomation) synchronizeStakeRequirements() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the current minimum stake from the global registry
    currentMinStake := automation.registry.MinStake

    // Fetch PoH and PoS stakes
    pohMinStake := automation.consensusEngine.GetPoHMinimumStake()
    posMinStake := automation.consensusEngine.GetPoSMinimumStake()

    // If either PoH or PoS stake differs from the global minimum, update it
    if pohMinStake != currentMinStake {
        automation.consensusEngine.SetPoHMinimumStake(currentMinStake)
        fmt.Printf("Updated PoH minimum stake to %.2f to match global setting.\n", currentMinStake)
    }

    if posMinStake != currentMinStake {
        automation.consensusEngine.SetPoSMinimumStake(currentMinStake)
        fmt.Printf("Updated PoS minimum stake to %.2f to match global setting.\n", currentMinStake)
    }

    // Reconcile any changes from other processes
    automation.reconcileStakeChanges(pohMinStake, posMinStake)
}

// reconcileStakeChanges ensures changes from other functions are synergized with the global registry
func (automation *SynchronizeMinimumStakeAutomation) reconcileStakeChanges(pohStake, posStake float64) {
    if pohStake > automation.registry.MinStake || posStake > automation.registry.MinStake {
        highestStake := max(pohStake, posStake)
        fmt.Printf("Reconciling global minimum stake to the highest detected stake: %.2f\n", highestStake)
        automation.registry.ApplyStakeChange(StakeChange{
            RequestedBy: "Reconciliation Process",
            NewMinStake: highestStake,
            Override:    false,
            Priority:    90, // Reconciliation priority
        })
    }
}

// max returns the maximum value between two floats
func max(a, b float64) float64 {
    if a > b {
        return a
    }
    return b
}

// RequestStakeChange allows other processes to request changes to the global minimum stake
func (automation *SynchronizeMinimumStakeAutomation) RequestStakeChange(requestedBy string, newMinStake float64, override bool, priority int) {
    automation.registry.ApplyStakeChange(StakeChange{
        RequestedBy:   requestedBy,
        NewMinStake:   newMinStake,
        Override:      override,
        Priority:      priority,
    })
}

// MonitorOtherStakeFunctions monitors other stake functions for conflicts and ensures synchronization overrides them
func (automation *SynchronizeMinimumStakeAutomation) MonitorOtherStakeFunctions() {
    // Continuously monitor for any stake changes or adjustments made by other components
    ticker := time.NewTicker(StakeCheckInterval)
    for range ticker.C {
        fmt.Println("Monitoring for conflicting stake adjustments...")
        automation.synchronizeStakeRequirements()
    }
}
