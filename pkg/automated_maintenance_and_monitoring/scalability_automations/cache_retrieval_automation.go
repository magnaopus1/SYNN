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
    CacheRetrievalInterval   = 3 * time.Minute // Interval for cache retrieval monitoring
    CacheRetrievalThreshold  = 80              // Cache usage threshold to trigger retrieval
    MaxCacheRetrievalRetries = 3               // Maximum number of retry attempts
)

// CacheRetrievalAutomation handles the automated cache retrieval process
type CacheRetrievalAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Integration with Synnergy Consensus
    ledgerInstance      *ledger.Ledger               // Ledger integration
    stateMutex          *sync.RWMutex                // Mutex for thread-safe operations
    retrievalRetries    map[string]int               // Retry count per node
    retrievalCycleCount int                          // Counter for cache retrieval cycles
}

// NewCacheRetrievalAutomation initializes the cache retrieval automation
func NewCacheRetrievalAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CacheRetrievalAutomation {
    return &CacheRetrievalAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        retrievalRetries:   make(map[string]int),
        retrievalCycleCount: 0,
    }
}

// StartCacheRetrievalMonitoring starts the continuous loop for cache retrieval monitoring
func (automation *CacheRetrievalAutomation) StartCacheRetrievalMonitoring() {
    ticker := time.NewTicker(CacheRetrievalInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndRetrieveCache()
        }
    }()
}

// monitorAndRetrieveCache monitors the cache usage status and triggers cache retrieval if necessary
func (automation *CacheRetrievalAutomation) monitorAndRetrieveCache() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch cache usage reports from Synnergy Consensus
    cacheReports := automation.consensusSystem.FetchCacheUsageReports()

    for _, report := range cacheReports {
        if report.Utilization >= CacheRetrievalThreshold {
            fmt.Printf("Cache retrieval needed for node %s (Usage: %d%%). Initiating retrieval.\n", report.NodeID, report.Utilization)
            automation.retrieveCache(report)
        } else {
            fmt.Printf("Cache usage on node %s is below retrieval threshold: %d%%.\n", report.NodeID, report.Utilization)
        }
    }

    automation.retrievalCycleCount++
    fmt.Printf("Cache retrieval cycle #%d completed.\n", automation.retrievalCycleCount)

    if automation.retrievalCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeRetrievalCycle()
    }
}

// retrieveCache handles the cache retrieval process when the usage exceeds the threshold
func (automation *CacheRetrievalAutomation) retrieveCache(report common.CacheUsageReport) {
    encryptedData := automation.encryptCacheData(report)

    // Try to retrieve cache through Synnergy Consensus
    success := automation.consensusSystem.InitiateCacheRetrieval(report.NodeID, encryptedData)
    if success {
        fmt.Printf("Cache retrieval completed successfully for node %s.\n", report.NodeID)
        automation.logRetrievalEvent(report, "Retrieved")
        automation.resetRetrievalRetries(report.NodeID)
    } else {
        fmt.Printf("Cache retrieval failed for node %s. Retrying...\n", report.NodeID)
        automation.retryCacheRetrieval(report)
    }
}

// retryCacheRetrieval retries cache retrieval if the initial attempt fails
func (automation *CacheRetrievalAutomation) retryCacheRetrieval(report common.CacheUsageReport) {
    automation.retrievalRetries[report.NodeID]++
    if automation.retrievalRetries[report.NodeID] < MaxCacheRetrievalRetries {
        automation.retrieveCache(report)
    } else {
        fmt.Printf("Max retry attempts reached for node %s. Cache retrieval failed.\n", report.NodeID)
        automation.logRetrievalFailure(report)
    }
}

// resetRetrievalRetries resets the retry count for a specific node
func (automation *CacheRetrievalAutomation) resetRetrievalRetries(nodeID string) {
    automation.retrievalRetries[nodeID] = 0
}

// finalizeRetrievalCycle finalizes the cache retrieval cycle and logs it in the ledger
func (automation *CacheRetrievalAutomation) finalizeRetrievalCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeCacheRetrievalCycle()
    if success {
        fmt.Println("Cache retrieval cycle finalized successfully.")
        automation.logCycleFinalization()
    } else {
        fmt.Println("Error finalizing cache retrieval cycle.")
    }
}

// logRetrievalEvent logs successful cache retrieval in the ledger
func (automation *CacheRetrievalAutomation) logRetrievalEvent(report common.CacheUsageReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-retrieval-%s-%s", report.NodeID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Retrieval Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s cache %s successfully.", report.NodeID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with cache retrieval event for node %s.\n", report.NodeID)
}

// logRetrievalFailure logs cache retrieval failure in the ledger
func (automation *CacheRetrievalAutomation) logRetrievalFailure(report common.CacheUsageReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-retrieval-failure-%s", report.NodeID),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Retrieval Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Cache retrieval failed for node %s after maximum retries.", report.NodeID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with cache retrieval failure for node %s.\n", report.NodeID)
}

// logCycleFinalization logs the finalization of the cache retrieval cycle in the ledger
func (automation *CacheRetrievalAutomation) logCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cache-retrieval-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Cache Retrieval Cycle",
        Status:    "Finalized",
        Details:   "Cache retrieval cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with cache retrieval cycle finalization.")
}

// encryptCacheData encrypts cache data before retrieval operations
func (automation *CacheRetrievalAutomation) encryptCacheData(report common.CacheUsageReport) []byte {
    encryptedData, err := encryption.EncryptData(report.Data)
    if err != nil {
        fmt.Printf("Error encrypting cache data for node %s: %v\n", report.NodeID, err)
        return report.Data
    }

    fmt.Printf("Cache data successfully encrypted for node %s.\n", report.NodeID)
    return encryptedData
}
