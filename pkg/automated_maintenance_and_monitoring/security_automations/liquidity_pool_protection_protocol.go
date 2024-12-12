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
    LiquidityPoolMonitoringInterval = 20 * time.Second // Interval for monitoring liquidity pools
    MaxPoolProtectionRetries        = 5                // Maximum retries for liquidity pool protection
    SubBlocksPerBlock               = 1000             // Number of sub-blocks in a block
    PoolProtectionThreshold         = 0.01         // Minimum liquidity value for activating protection
)

// LiquidityPoolProtectionProtocol handles the protection and security enforcement of liquidity pools
type LiquidityPoolProtectionProtocol struct {
    consensusSystem          *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance           *ledger.Ledger               // Ledger for logging liquidity pool-related events
    stateMutex               *sync.RWMutex                // Mutex for thread-safe access
    protectionRetryCount     map[string]int               // Counter for retrying protection application
    protectionCycleCount     int                          // Counter for protection monitoring cycles
    protectedLiquidityPools  map[string]bool              // Tracks liquidity pools currently under protection
}

// NewLiquidityPoolProtectionProtocol initializes the automation for liquidity pool protection
func NewLiquidityPoolProtectionProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *LiquidityPoolProtectionProtocol {
    return &LiquidityPoolProtectionProtocol{
        consensusSystem:         consensusSystem,
        ledgerInstance:          ledgerInstance,
        stateMutex:              stateMutex,
        protectionRetryCount:    make(map[string]int),
        protectedLiquidityPools: make(map[string]bool),
        protectionCycleCount:    0,
    }
}

// StartLiquidityPoolProtectionMonitoring starts the continuous loop for monitoring and enforcing liquidity pool protection
func (protocol *LiquidityPoolProtectionProtocol) StartLiquidityPoolProtectionMonitoring() {
    ticker := time.NewTicker(LiquidityPoolMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorAndEnforceLiquidityProtection()
        }
    }()
}

// monitorAndEnforceLiquidityProtection monitors liquidity pools and applies protection when necessary
func (protocol *LiquidityPoolProtectionProtocol) monitorAndEnforceLiquidityProtection() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch active liquidity pools and check their status
    liquidityPools := protocol.consensusSystem.FetchActiveLiquidityPools()

    for _, pool := range liquidityPools {
        if pool.LiquidityValue >= PoolProtectionThreshold && !protocol.protectedLiquidityPools[pool.ID] {
            fmt.Printf("Activating protection for liquidity pool ID: %s\n", pool.ID)
            protocol.activateProtection(pool)
        } else if pool.LiquidityValue < PoolProtectionThreshold && protocol.protectedLiquidityPools[pool.ID] {
            fmt.Printf("Deactivating protection for liquidity pool ID: %s due to low liquidity.\n", pool.ID)
            protocol.deactivateProtection(pool)
        } else {
            fmt.Printf("Liquidity pool ID: %s is stable. No protection changes required.\n", pool.ID)
        }
    }

    protocol.protectionCycleCount++
    fmt.Printf("Liquidity pool protection cycle #%d completed.\n", protocol.protectionCycleCount)

    if protocol.protectionCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeProtectionCycle()
    }
}

// activateProtection applies protection to a liquidity pool based on its value
func (protocol *LiquidityPoolProtectionProtocol) activateProtection(pool common.LiquidityPool) {
    encryptedPoolData := protocol.encryptPoolData(pool)

    // Attempt to apply protection through the Synnergy Consensus system
    protectionSuccess := protocol.consensusSystem.ProtectLiquidityPool(pool, encryptedPoolData)

    if protectionSuccess {
        protocol.protectedLiquidityPools[pool.ID] = true
        fmt.Printf("Protection activated for liquidity pool ID: %s.\n", pool.ID)
        protocol.logProtectionEvent(pool, "Activated")
        protocol.resetProtectionRetry(pool.ID)
    } else {
        fmt.Printf("Error activating protection for liquidity pool ID: %s. Retrying...\n", pool.ID)
        protocol.retryPoolProtection(pool)
    }
}

// deactivateProtection removes protection from a liquidity pool due to decreased liquidity
func (protocol *LiquidityPoolProtectionProtocol) deactivateProtection(pool common.LiquidityPool) {
    encryptedPoolData := protocol.encryptPoolData(pool)

    // Attempt to remove protection through the Synnergy Consensus system
    removalSuccess := protocol.consensusSystem.RemoveLiquidityPoolProtection(pool, encryptedPoolData)

    if removalSuccess {
        delete(protocol.protectedLiquidityPools, pool.ID)
        fmt.Printf("Protection deactivated for liquidity pool ID: %s.\n", pool.ID)
        protocol.logProtectionEvent(pool, "Deactivated")
    } else {
        fmt.Printf("Error deactivating protection for liquidity pool ID: %s.\n", pool.ID)
    }
}

// retryPoolProtection retries protection activation in case of a failure
func (protocol *LiquidityPoolProtectionProtocol) retryPoolProtection(pool common.LiquidityPool) {
    protocol.protectionRetryCount[pool.ID]++
    if protocol.protectionRetryCount[pool.ID] < MaxPoolProtectionRetries {
        protocol.activateProtection(pool)
    } else {
        fmt.Printf("Max retries reached for protecting liquidity pool ID: %s. Protection failed.\n", pool.ID)
        protocol.logProtectionFailure(pool)
    }
}

// resetProtectionRetry resets the retry count for a liquidity pool protection attempt
func (protocol *LiquidityPoolProtectionProtocol) resetProtectionRetry(poolID string) {
    protocol.protectionRetryCount[poolID] = 0
}

// finalizeProtectionCycle finalizes the liquidity protection cycle and logs the result in the ledger
func (protocol *LiquidityPoolProtectionProtocol) finalizeProtectionCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeProtectionCycle()
    if success {
        fmt.Println("Liquidity protection cycle finalized successfully.")
        protocol.logProtectionCycleFinalization()
    } else {
        fmt.Println("Error finalizing liquidity protection cycle.")
    }
}

// logProtectionEvent logs a liquidity pool protection event into the ledger
func (protocol *LiquidityPoolProtectionProtocol) logProtectionEvent(pool common.LiquidityPool, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("liquidity-pool-protection-%s-%s", pool.ID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Liquidity Pool Protection Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Liquidity pool ID: %s had protection %s.", pool.ID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with liquidity pool protection event for pool ID: %s.\n", pool.ID)
}

// logProtectionFailure logs the failure of liquidity pool protection into the ledger
func (protocol *LiquidityPoolProtectionProtocol) logProtectionFailure(pool common.LiquidityPool) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("liquidity-pool-protection-failure-%s", pool.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Liquidity Pool Protection Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to protect liquidity pool ID: %s after maximum retries.", pool.ID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with protection failure for pool ID: %s.\n", pool.ID)
}

// logProtectionCycleFinalization logs the finalization of a liquidity pool protection cycle into the ledger
func (protocol *LiquidityPoolProtectionProtocol) logProtectionCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("liquidity-protection-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Protection Cycle Finalization",
        Status:    "Finalized",
        Details:   "Liquidity protection cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with liquidity protection cycle finalization.")
}

// encryptPoolData encrypts liquidity pool data before applying protection or removal
func (protocol *LiquidityPoolProtectionProtocol) encryptPoolData(pool common.LiquidityPool) common.LiquidityPool {
    encryptedData, err := encryption.EncryptData(pool.Data)
    if err != nil {
        fmt.Println("Error encrypting liquidity pool data:", err)
        return pool
    }

    pool.EncryptedData = encryptedData
    fmt.Println("Liquidity pool data successfully encrypted for pool ID:", pool.ID)
    return pool
}

// triggerEmergencyProtectionLock triggers emergency protection of a liquidity pool in case of suspicious activity
func (protocol *LiquidityPoolProtectionProtocol) triggerEmergencyProtectionLock(poolID string) {
    fmt.Printf("Emergency protection triggered for liquidity pool ID: %s.\n", poolID)
    pool := protocol.consensusSystem.GetLiquidityPoolByID(poolID)
    encryptedData := protocol.encryptPoolData(pool)

    success := protocol.consensusSystem.TriggerEmergencyLiquidityPoolLock(poolID, encryptedData)

    if success {
        protocol.logProtectionEvent(pool, "Emergency Locked")
        fmt.Println("Emergency liquidity pool protection executed successfully.")
    } else {
        fmt.Println("Emergency liquidity pool protection failed.")
    }
}
