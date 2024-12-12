package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/cache"
)

// Constants for Adaptive Caching
const (
	CacheCheckInterval         = 10 * time.Second // Interval for cache monitoring and updates
	CacheMaxSize               = 1000000          // Max cache size (example in bytes)
	CacheThresholdPercentage   = 0.8              // When cache reaches 80%, start cleanup
	CacheHighDemandThreshold   = 5000             // High demand access count threshold
	CacheCleanupPercentage     = 0.2              // Cleanup 20% of cache during cleanup cycles
)

// AdaptiveCachingExecutionAutomation manages adaptive cache operations
type AdaptiveCachingExecutionAutomation struct {
	cacheInstance     *cache.Cache
	consensusSystem   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	cacheAccessMutex  *sync.RWMutex
}

// NewAdaptiveCachingExecutionAutomation initializes adaptive caching automation
func NewAdaptiveCachingExecutionAutomation(cacheInstance *cache.Cache, consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, cacheAccessMutex *sync.RWMutex) *AdaptiveCachingExecutionAutomation {
	return &AdaptiveCachingExecutionAutomation{
		cacheInstance:    cacheInstance,
		consensusSystem:  consensusSystem,
		ledgerInstance:   ledgerInstance,
		cacheAccessMutex: cacheAccessMutex,
	}
}

// StartAdaptiveCachingMonitor starts the adaptive caching automation in a continuous loop
func (automation *AdaptiveCachingExecutionAutomation) StartAdaptiveCachingMonitor() {
	ticker := time.NewTicker(CacheCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateCacheUsage()
		}
	}()
}

// evaluateCacheUsage checks the cache size and triggers cleanup or optimization as necessary
func (automation *AdaptiveCachingExecutionAutomation) evaluateCacheUsage() {
	automation.cacheAccessMutex.Lock()
	defer automation.cacheAccessMutex.Unlock()

	cacheSize := automation.cacheInstance.GetCurrentSize()
	cacheCapacity := CacheMaxSize

	if float64(cacheSize)/float64(cacheCapacity) >= CacheThresholdPercentage {
		automation.initiateCacheCleanup()
	} else {
		automation.optimizeCacheBasedOnDemand()
	}
}

// initiateCacheCleanup clears low-demand data when cache exceeds threshold
func (automation *AdaptiveCachingExecutionAutomation) initiateCacheCleanup() {
	err := automation.cacheInstance.Cleanup(CacheCleanupPercentage)
	if err != nil {
		fmt.Println("Failed to clean up cache:", err)
		return
	}

	// Log cache cleanup into ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("adaptive-cache-cleanup-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Cache Cleanup",
		Status:    "Success",
		Details:   "Cache cleaned up due to high usage.",
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log cache cleanup:", err)
	} else {
		fmt.Println("Cache cleanup successfully logged in the ledger.")
	}
}

// optimizeCacheBasedOnDemand reorganizes cache based on data access frequency
func (automation *AdaptiveCachingExecutionAutomation) optimizeCacheBasedOnDemand() {
	highDemandItems := automation.cacheInstance.GetHighDemandItems(CacheHighDemandThreshold)
	lowDemandItems := automation.cacheInstance.GetLowDemandItems()

	// Move high-demand items to faster storage or prioritize access
	for _, item := range highDemandItems {
		err := automation.cacheInstance.Prioritize(item)
		if err != nil {
			fmt.Println("Failed to prioritize cache item:", err)
		}
	}

	// De-prioritize low-demand items
	for _, item := range lowDemandItems {
		err := automation.cacheInstance.Deprioritize(item)
		if err != nil {
			fmt.Println("Failed to de-prioritize cache item:", err)
		}
	}

	// Log optimization process into ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("adaptive-cache-optimization-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Cache Optimization",
		Status:    "Success",
		Details:   "Cache optimized based on data access patterns.",
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log cache optimization:", err)
	} else {
		fmt.Println("Cache optimization successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *AdaptiveCachingExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
