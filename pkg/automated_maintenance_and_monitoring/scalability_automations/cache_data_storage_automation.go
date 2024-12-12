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
    CacheMonitoringInterval   = 1 * time.Minute // Interval for cache storage monitoring
    CacheStorageLimit         = 85              // Cache storage limit threshold as a percentage
    MaxCacheRetryAttempts     = 3               // Maximum retry attempts for cache storage actions
)

// CacheDataStorageAutomation automates the process of monitoring and managing cache storage
type CacheDataStorageAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to the SynnergyConsensus system
    ledgerInstance   *ledger.Ledger               // Ledger for logging cache storage-related events
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    retryAttempt     map[string]int               // Retry attempt tracker for failed cache storage operations
    cacheCycleCount  int                          // Cycle count for cache storage management
}

// NewCacheDataStorageAutomation initializes the cache data storage automation
func NewCacheDataStorageAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CacheDataStorageAutomation {
    return &CacheDataStorageAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        retryAttempt:     make(map[string]int),
        cacheCycleCount:  0,
    }
}

// StartCacheStorageMonitoring begins the continuous loop for cache data storage monitoring
func (automation *CacheDataStorageAutomation) StartCacheStorageMonitoring() {
    ticker := time.NewTicker(CacheMonitoringInterval)

    go func() {
        for range ticker.C {
            automation.monitorCacheStorage()
        }
    }()
}

// monitorCacheStorage checks cache usage and triggers actions to manage storage if needed
func (automation *CacheDataStorageAutomation) monitorCacheStorage() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Retrieve the cache storage report from the consensus system
    cacheReports := automation.consensusSystem.FetchCacheUsageReports()

    for _, report := range cacheReports {
        if report.Utilization >= CacheStorageLimit {
            fmt.Printf("Cache storage limit reached on node %s with %d%% usage. Taking action.\n", report.NodeID, report.Utilization)
            automation.manageCacheStorage(report)
        } else {
            fmt.Printf("Cache usage on node %s is %d%%, within limits.\n", report.NodeID, report.Utilization)
        }
    }

    automation.cacheCycleCount++
    fmt.Printf("Cache storage management cycle #%d completed.\n", automation.cacheCycleCount)

    if automation.cacheCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeCacheCycle()
    }
}

// manageCacheStorage attempts to take action if cache storage is exceeding the limit
func (automation *CacheDataStorageAutomation) manageCacheStorage(report common.CacheUsageReport) {
    encryptedData := automation.encryptCacheData(report)

    // Try to optimize or clean cache storage through Synnergy Consensus
    success := automation.consensusSystem.InitiateCacheStorageOptimization(report.NodeID, encryptedData)
    if success {
        fmt.Printf("Cache storage optimized successfully for node %s.\n", report.NodeID)
        automation.logCacheStorageEvent(report, "Optimized")
        automation.resetCacheRetry(report.NodeID)
    } else {
        fmt.Printf("Cache storage optimization failed for node %s. Retrying...\n", report.NodeID)
        automation.retryCacheStorage(report)
    }
}

// retryCacheStorage retries cache storage management if initial attempts fail
func (automation *CacheDataStorageAutomation) retryCacheStorage(report common.CacheUsageReport) {
    automation.retryAttempt[report.NodeID]++
    if automation.retryAttempt[report.NodeID] < MaxCacheRetryAttempts {
        automation.manageCacheStorage(report)
    } else {
        fmt.Printf("Max retry attempts reached for node %s. Cache management failed.\n", report.NodeID)
        automation.logCacheFailure(report)
    }
}

// resetCacheRetry resets the retry counter for cache storage actions
func (automation *CacheDataStorageAutomation) resetCacheRetry(nodeID string) {
    automation.retryAttempt[nodeID] = 0
}

// finalizeCacheCycle finalizes the cache storage cycle and logs the event in the ledger
func (automation *CacheDataStorageAutomation) finalizeCacheCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeCacheCycle()
    if success {
        fmt.Println("Cache storage management cycle finalized successfully.")
        automation.logCacheCycleFinalization()
    } else {
        fmt.Println("Error finalizing cache storage management cycle.")
    }
}

// logCacheStorageEvent logs successful cache storage actions in the ledger
func (automation *CacheDataStorageAutomation) logCacheStorageEvent(report common.CacheUsageReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-storage-%s-%s", report.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Storage Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s cache storage %s successfully.", report.NodeID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with cache storage event for node %s.\n", report.NodeID)
}

// logCacheFailure logs the failure of cache storage actions into the ledger
func (automation *CacheDataStorageAutomation) logCacheFailure(report common.CacheUsageReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-failure-%s", report.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Storage Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Cache storage action failed for node %s after maximum retries.", report.NodeID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with cache storage failure for node %s.\n", report.NodeID)
}

// logCacheCycleFinalization logs the finalization of a cache management cycle into the ledger
func (automation *CacheDataStorageAutomation) logCacheCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Management Cycle",
        Status:    "Finalized",
        Details:   "Cache management cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with cache management cycle finalization.")
}

// encryptCacheData encrypts cache data before performing actions
func (automation *CacheDataStorageAutomation) encryptCacheData(report common.CacheUsageReport) []byte {
    encryptedData, err := encryption.EncryptData(report.Data)
    if err != nil {
        fmt.Printf("Error encrypting cache data for node %s: %v\n", report.NodeID, err)
        return report.Data
    }

    fmt.Printf("Cache data successfully encrypted for node %s.\n", report.NodeID)
    return encryptedData
}
