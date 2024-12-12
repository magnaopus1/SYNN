package automations

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
    LiquidityLockMonitoringInterval = 30 * time.Second // Interval for monitoring liquidity locks
    MaxLiquidityUnlockRetries       = 3                // Maximum retries for liquidity unlocking
    SubBlocksPerBlock               = 1000             // Number of sub-blocks in a block
    LiquidityLockPeriod             = 7 * 24 * time.Hour // Liquidity locking period
)

// LiquidityLockingSecurityProtocol handles the security and enforcement of liquidity locking and unlocking mechanisms
type LiquidityLockingSecurityProtocol struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger for logging liquidity locking-related events
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    lockedLiquidityPools   map[string]time.Time         // Tracks locked liquidity pools and their unlock times
    unlockRetryCount       map[string]int               // Counter for retrying liquidity unlocks
    lockCycleCount         int                          // Counter for liquidity locking cycles
}

// NewLiquidityLockingSecurityProtocol initializes the automation for liquidity locking and unlocking
func NewLiquidityLockingSecurityProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *LiquidityLockingSecurityProtocol {
    return &LiquidityLockingSecurityProtocol{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        lockedLiquidityPools: make(map[string]time.Time),
        unlockRetryCount:     make(map[string]int),
        lockCycleCount:       0,
    }
}

// StartLiquidityLockingMonitoring starts the continuous loop for monitoring and enforcing liquidity locking policies
func (protocol *LiquidityLockingSecurityProtocol) StartLiquidityLockingMonitoring() {
    ticker := time.NewTicker(LiquidityLockMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorAndEnforceLiquidityLocking()
        }
    }()
}

// monitorAndEnforceLiquidityLocking monitors locked liquidity pools and triggers unlocks or lock renewals based on conditions
func (protocol *LiquidityLockingSecurityProtocol) monitorAndEnforceLiquidityLocking() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    currentTime := time.Now()

    // Iterate over locked liquidity pools and check if they need to be unlocked
    for poolID, unlockTime := range protocol.lockedLiquidityPools {
        if currentTime.After(unlockTime) {
            fmt.Printf("Unlocking liquidity for pool ID: %s\n", poolID)
            protocol.unlockLiquidity(poolID)
        } else {
            fmt.Printf("Liquidity for pool ID: %s remains locked. Unlock time: %s\n", poolID, unlockTime)
        }
    }

    protocol.lockCycleCount++
    fmt.Printf("Liquidity locking cycle #%d completed.\n", protocol.lockCycleCount)

    if protocol.lockCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeLiquidityLockCycle()
    }
}

// lockLiquidity locks a specific liquidity pool for a set duration
func (protocol *LiquidityLockingSecurityProtocol) lockLiquidity(poolID string) {
    lockEndTime := time.Now().Add(LiquidityLockPeriod)

    // Encrypt liquidity data before locking
    encryptedLiquidityData := protocol.encryptLiquidityData(poolID)

    // Lock liquidity in the pool via Synnergy Consensus
    lockSuccess := protocol.consensusSystem.LockLiquidityPool(poolID, encryptedLiquidityData)

    if lockSuccess {
        protocol.lockedLiquidityPools[poolID] = lockEndTime
        fmt.Printf("Liquidity locked for pool ID: %s until %s.\n", poolID, lockEndTime)
        protocol.logLiquidityLockEvent(poolID, "Locked")
    } else {
        fmt.Printf("Failed to lock liquidity for pool ID: %s.\n", poolID)
    }
}

// unlockLiquidity attempts to unlock liquidity from a pool after the lock period expires
func (protocol *LiquidityLockingSecurityProtocol) unlockLiquidity(poolID string) {
    encryptedUnlockData := protocol.encryptLiquidityData(poolID)

    // Attempt to unlock the liquidity pool via Synnergy Consensus
    unlockSuccess := protocol.consensusSystem.UnlockLiquidityPool(poolID, encryptedUnlockData)

    if unlockSuccess {
        fmt.Printf("Liquidity unlocked for pool ID: %s.\n", poolID)
        protocol.logLiquidityUnlockEvent(poolID, "Unlocked")
        delete(protocol.lockedLiquidityPools, poolID)
        protocol.resetUnlockRetry(poolID)
    } else {
        fmt.Printf("Failed to unlock liquidity for pool ID: %s. Retrying...\n", poolID)
        protocol.retryLiquidityUnlock(poolID)
    }
}

// retryLiquidityUnlock retries liquidity unlocking if the initial unlock attempt fails
func (protocol *LiquidityLockingSecurityProtocol) retryLiquidityUnlock(poolID string) {
    protocol.unlockRetryCount[poolID]++
    if protocol.unlockRetryCount[poolID] < MaxLiquidityUnlockRetries {
        protocol.unlockLiquidity(poolID)
    } else {
        fmt.Printf("Max retries reached for unlocking liquidity in pool ID: %s. Unlock failed.\n", poolID)
        protocol.logLiquidityUnlockFailure(poolID)
    }
}

// resetUnlockRetry resets the retry counter for unlocking liquidity in a specific pool
func (protocol *LiquidityLockingSecurityProtocol) resetUnlockRetry(poolID string) {
    protocol.unlockRetryCount[poolID] = 0
}

// finalizeLiquidityLockCycle finalizes the liquidity locking cycle and logs the result in the ledger
func (protocol *LiquidityLockingSecurityProtocol) finalizeLiquidityLockCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeLiquidityLockCycle()
    if success {
        fmt.Println("Liquidity locking cycle finalized successfully.")
        protocol.logLiquidityLockCycleFinalization()
    } else {
        fmt.Println("Error finalizing liquidity locking cycle.")
    }
}

// logLiquidityLockEvent logs a liquidity lock event into the ledger
func (protocol *LiquidityLockingSecurityProtocol) logLiquidityLockEvent(poolID, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("liquidity-lock-%s-%s", poolID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Liquidity Lock Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Liquidity for pool ID: %s was %s.", poolID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with liquidity lock event for pool ID: %s.\n", poolID)
}

// logLiquidityUnlockEvent logs a liquidity unlock event into the ledger
func (protocol *LiquidityLockingSecurityProtocol) logLiquidityUnlockEvent(poolID, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("liquidity-unlock-%s-%s", poolID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Liquidity Unlock Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Liquidity for pool ID: %s was %s.", poolID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with liquidity unlock event for pool ID: %s.\n", poolID)
}

// logLiquidityUnlockFailure logs the failure of a liquidity unlock into the ledger
func (protocol *LiquidityLockingSecurityProtocol) logLiquidityUnlockFailure(poolID string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("liquidity-unlock-failure-%s", poolID),
        Timestamp: time.Now().Unix(),
        Type:      "Liquidity Unlock Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to unlock liquidity for pool ID: %s after maximum retries.", poolID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with liquidity unlock failure for pool ID: %s.\n", poolID)
}

// logLiquidityLockCycleFinalization logs the finalization of a liquidity lock cycle into the ledger
func (protocol *LiquidityLockingSecurityProtocol) logLiquidityLockCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("liquidity-lock-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Liquidity Lock Cycle Finalization",
        Status:    "Finalized",
        Details:   "Liquidity lock cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with liquidity lock cycle finalization.")
}

// encryptLiquidityData encrypts liquidity data before locking or unlocking
func (protocol *LiquidityLockingSecurityProtocol) encryptLiquidityData(poolID string) string {
    encryptedData, err := encryption.EncryptData([]byte(poolID))
    if err != nil {
        fmt.Println("Error encrypting liquidity data:", err)
        return poolID
    }

    fmt.Println("Liquidity data successfully encrypted for pool ID:", poolID)
    return string(encryptedData)
}

// triggerEmergencyLiquidityLocking triggers the emergency locking of liquidity in case of security breaches
func (protocol *LiquidityLockingSecurityProtocol) triggerEmergencyLiquidityLocking(poolID string) {
    fmt.Printf("Emergency liquidity locking triggered for pool ID: %s.\n", poolID)
    encryptedData := protocol.encryptLiquidityData(poolID)
    success := protocol.consensusSystem.TriggerEmergencyLiquidityLock(poolID, encryptedData)

    if success {
        protocol.logLiquidityLockEvent(poolID, "Emergency Locked")
        fmt.Println("Emergency liquidity locking executed successfully.")
    } else {
        fmt.Println("Emergency liquidity locking failed.")
    }
}
