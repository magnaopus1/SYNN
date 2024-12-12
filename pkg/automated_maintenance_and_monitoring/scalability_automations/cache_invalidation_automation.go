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
    CacheInvalidationInterval  = 5 * time.Minute // Interval for monitoring and invalidating cache
    CacheInvalidationThreshold = 90              // Cache usage threshold as percentage for invalidation
    MaxInvalidationRetries     = 3               // Maximum number of retry attempts for invalidation
)

// CacheInvalidationAutomation automates the process of invalidating cache data
type CacheInvalidationAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Synnergy Consensus integration
    ledgerInstance      *ledger.Ledger               // Ledger integration
    stateMutex          *sync.RWMutex                // Mutex for thread-safe operations
    invalidationRetries map[string]int               // Retry count per node for cache invalidation
    cacheCycleCount     int                          // Cache invalidation cycle counter
}

// NewCacheInvalidationAutomation initializes the cache invalidation automation
func NewCacheInvalidationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CacheInvalidationAutomation {
    return &CacheInvalidationAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        invalidationRetries: make(map[string]int),
        cacheCycleCount:     0,
    }
}

// StartCacheInvalidationMonitoring starts the continuous loop for cache invalidation
func (automation *CacheInvalidationAutomation) StartCacheInvalidationMonitoring() {
    ticker := time.NewTicker(CacheInvalidationInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndInvalidateCache()
        }
    }()
}

// monitorAndInvalidateCache checks the cache status and triggers invalidation when necessary
func (automation *CacheInvalidationAutomation) monitorAndInvalidateCache() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch cache usage report from Synnergy Consensus
    cacheReports := automation.consensusSystem.FetchCacheUsageReports()

    for _, report := range cacheReports {
        if report.Utilization >= CacheInvalidationThreshold {
            fmt.Printf("Cache invalidation needed for node %s (Usage: %d%%). Initiating invalidation.\n", report.NodeID, report.Utilization)
            automation.invalidateCache(report)
        } else {
            fmt.Printf("Cache usage on node %s within acceptable limits: %d%%.\n", report.NodeID, report.Utilization)
        }
    }

    automation.cacheCycleCount++
    fmt.Printf("Cache invalidation cycle #%d completed.\n", automation.cacheCycleCount)

    if automation.cacheCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeInvalidationCycle()
    }
}

// invalidateCache attempts to invalidate cache if the usage exceeds the threshold
func (automation *CacheInvalidationAutomation) invalidateCache(report common.CacheUsageReport) {
    encryptedData := automation.encryptCacheData(report)

    // Try to invalidate cache through Synnergy Consensus
    success := automation.consensusSystem.InitiateCacheInvalidation(report.NodeID, encryptedData)
    if success {
        fmt.Printf("Cache invalidated successfully for node %s.\n", report.NodeID)
        automation.logInvalidationEvent(report, "Invalidated")
        automation.resetInvalidationRetries(report.NodeID)
    } else {
        fmt.Printf("Cache invalidation failed for node %s. Retrying...\n", report.NodeID)
        automation.retryInvalidation(report)
    }
}

// retryInvalidation retries cache invalidation if initial attempts fail
func (automation *CacheInvalidationAutomation) retryInvalidation(report common.CacheUsageReport) {
    automation.invalidationRetries[report.NodeID]++
    if automation.invalidationRetries[report.NodeID] < MaxInvalidationRetries {
        automation.invalidateCache(report)
    } else {
        fmt.Printf("Max retry attempts reached for node %s. Cache invalidation failed.\n", report.NodeID)
        automation.logInvalidationFailure(report)
    }
}

// resetInvalidationRetries resets the retry counter for cache invalidation attempts
func (automation *CacheInvalidationAutomation) resetInvalidationRetries(nodeID string) {
    automation.invalidationRetries[nodeID] = 0
}

// finalizeInvalidationCycle finalizes the cache invalidation cycle and logs it in the ledger
func (automation *CacheInvalidationAutomation) finalizeInvalidationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeInvalidationCycle()
    if success {
        fmt.Println("Cache invalidation cycle finalized successfully.")
        automation.logCycleFinalization()
    } else {
        fmt.Println("Error finalizing cache invalidation cycle.")
    }
}

// logInvalidationEvent logs successful cache invalidation in the ledger
func (automation *CacheInvalidationAutomation) logInvalidationEvent(report common.CacheUsageReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-invalidation-%s-%s", report.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Invalidation Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s cache %s successfully.", report.NodeID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with cache invalidation event for node %s.\n", report.NodeID)
}

// logInvalidationFailure logs failed cache invalidation attempts in the ledger
func (automation *CacheInvalidationAutomation) logInvalidationFailure(report common.CacheUsageReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-invalidation-failure-%s", report.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Invalidation Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Cache invalidation failed for node %s after maximum retries.", report.NodeID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with cache invalidation failure for node %s.\n", report.NodeID)
}

// logCycleFinalization logs the finalization of a cache invalidation cycle in the ledger
func (automation *CacheInvalidationAutomation) logCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-invalidation-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Invalidation Cycle",
        Status:    "Finalized",
        Details:   "Cache invalidation cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with cache invalidation cycle finalization.")
}

// encryptCacheData encrypts cache data before invalidation action is performed
func (automation *CacheInvalidationAutomation) encryptCacheData(report common.CacheUsageReport) []byte {
    encryptedData, err := encryption.EncryptData(report.Data)
    if err != nil {
        fmt.Printf("Error encrypting cache data for node %s: %v\n", report.NodeID, err)
        return report.Data
    }

    fmt.Printf("Cache data successfully encrypted for node %s.\n", report.NodeID)
    return encryptedData
}
