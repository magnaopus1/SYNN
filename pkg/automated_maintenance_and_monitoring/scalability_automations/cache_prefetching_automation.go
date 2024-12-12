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
    CachePrefetchingInterval  = 5 * time.Minute // Interval for cache prefetching
    CachePrefetchingThreshold = 75              // Cache usage threshold for initiating prefetching
    MaxPrefetchingRetries     = 3               // Maximum number of retry attempts for prefetching
)

// CachePrefetchingAutomation automates the process of prefetching cache data
type CachePrefetchingAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Synnergy Consensus integration
    ledgerInstance      *ledger.Ledger               // Ledger integration
    stateMutex          *sync.RWMutex                // Mutex for thread-safe operations
    prefetchingRetries  map[string]int               // Retry count per node for cache prefetching
    cacheCycleCount     int                          // Cache prefetching cycle counter
}

// NewCachePrefetchingAutomation initializes the cache prefetching automation
func NewCachePrefetchingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CachePrefetchingAutomation {
    return &CachePrefetchingAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        prefetchingRetries: make(map[string]int),
        cacheCycleCount:    0,
    }
}

// StartCachePrefetchingMonitoring starts the continuous loop for cache prefetching
func (automation *CachePrefetchingAutomation) StartCachePrefetchingMonitoring() {
    ticker := time.NewTicker(CachePrefetchingInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndPrefetchCache()
        }
    }()
}

// monitorAndPrefetchCache checks the cache usage status and triggers prefetching when necessary
func (automation *CachePrefetchingAutomation) monitorAndPrefetchCache() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch cache usage reports from Synnergy Consensus
    cacheReports := automation.consensusSystem.FetchCacheUsageReports()

    for _, report := range cacheReports {
        if report.Utilization <= CachePrefetchingThreshold {
            fmt.Printf("Cache prefetching needed for node %s (Usage: %d%%). Initiating prefetching.\n", report.NodeID, report.Utilization)
            automation.prefetchCache(report)
        } else {
            fmt.Printf("Cache usage on node %s is above the prefetching threshold: %d%%.\n", report.NodeID, report.Utilization)
        }
    }

    automation.cacheCycleCount++
    fmt.Printf("Cache prefetching cycle #%d completed.\n", automation.cacheCycleCount)

    if automation.cacheCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizePrefetchingCycle()
    }
}

// prefetchCache attempts to prefetch cache when usage is below the threshold
func (automation *CachePrefetchingAutomation) prefetchCache(report common.CacheUsageReport) {
    encryptedData := automation.encryptPrefetchData(report)

    // Try to prefetch cache through Synnergy Consensus
    success := automation.consensusSystem.InitiateCachePrefetching(report.NodeID, encryptedData)
    if success {
        fmt.Printf("Cache prefetching completed successfully for node %s.\n", report.NodeID)
        automation.logPrefetchingEvent(report, "Prefetched")
        automation.resetPrefetchingRetries(report.NodeID)
    } else {
        fmt.Printf("Cache prefetching failed for node %s. Retrying...\n", report.NodeID)
        automation.retryPrefetching(report)
    }
}

// retryPrefetching retries cache prefetching if the initial attempt fails
func (automation *CachePrefetchingAutomation) retryPrefetching(report common.CacheUsageReport) {
    automation.prefetchingRetries[report.NodeID]++
    if automation.prefetchingRetries[report.NodeID] < MaxPrefetchingRetries {
        automation.prefetchCache(report)
    } else {
        fmt.Printf("Max retry attempts reached for node %s. Cache prefetching failed.\n", report.NodeID)
        automation.logPrefetchingFailure(report)
    }
}

// resetPrefetchingRetries resets the retry counter for cache prefetching attempts
func (automation *CachePrefetchingAutomation) resetPrefetchingRetries(nodeID string) {
    automation.prefetchingRetries[nodeID] = 0
}

// finalizePrefetchingCycle finalizes the cache prefetching cycle and logs it in the ledger
func (automation *CachePrefetchingAutomation) finalizePrefetchingCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizePrefetchingCycle()
    if success {
        fmt.Println("Cache prefetching cycle finalized successfully.")
        automation.logCycleFinalization()
    } else {
        fmt.Println("Error finalizing cache prefetching cycle.")
    }
}

// logPrefetchingEvent logs successful cache prefetching in the ledger
func (automation *CachePrefetchingAutomation) logPrefetchingEvent(report common.CacheUsageReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-prefetching-%s-%s", report.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Prefetching Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s cache %s successfully.", report.NodeID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with cache prefetching event for node %s.\n", report.NodeID)
}

// logPrefetchingFailure logs failed cache prefetching attempts in the ledger
func (automation *CachePrefetchingAutomation) logPrefetchingFailure(report common.CacheUsageReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-prefetching-failure-%s", report.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Prefetching Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Cache prefetching failed for node %s after maximum retries.", report.NodeID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with cache prefetching failure for node %s.\n", report.NodeID)
}

// logCycleFinalization logs the finalization of a cache prefetching cycle in the ledger
func (automation *CachePrefetchingAutomation) logCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-prefetching-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Prefetching Cycle",
        Status:    "Finalized",
        Details:   "Cache prefetching cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with cache prefetching cycle finalization.")
}

// encryptPrefetchData encrypts cache data before prefetching action is performed
func (automation *CachePrefetchingAutomation) encryptPrefetchData(report common.CacheUsageReport) []byte {
    encryptedData, err := encryption.EncryptData(report.Data)
    if err != nil {
        fmt.Printf("Error encrypting cache data for node %s: %v\n", report.NodeID, err)
        return report.Data
    }

    fmt.Printf("Cache data successfully encrypted for node %s.\n", report.NodeID)
    return encryptedData
}
