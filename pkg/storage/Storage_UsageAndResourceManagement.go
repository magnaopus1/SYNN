package storage

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
)

// StorageTrackDataUsage tracks data usage within the storage system
func StorageTrackDataUsage(dataID string, usage int64) error {
    if err := common.TrackDataUsage(dataID, usage); err != nil {
        return fmt.Errorf("data usage tracking failed: %v", err)
    }
    ledger.RecordDataUsage(dataID, usage, time.Now())
    fmt.Printf("Data usage for %s tracked: %d bytes\n", dataID, usage)
    return nil
}

// StorageGenerateUsageReport generates a usage report for all storage data
func StorageGenerateUsageReport() (string, error) {
    report, err := ledger.GenerateDataUsageReport()
    if err != nil {
        return "", fmt.Errorf("usage report generation failed: %v", err)
    }
    fmt.Println("Data Usage Report:\n", report)
    return report, nil
}

// StorageAuditDataUsage audits data usage for compliance and accuracy
func StorageAuditDataUsage(dataID string) error {
    if !common.VerifyDataUsageCompliance(dataID) {
        return fmt.Errorf("data usage audit failed for %s", dataID)
    }
    fmt.Printf("Data usage for %s passed audit\n", dataID)
    return nil
}

// StorageFetchDataUsageHistory retrieves the usage history for a specific data item
func StorageFetchDataUsageHistory(dataID string) ([]string, error) {
    history, err := ledger.FetchDataUsageHistory(dataID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch data usage history: %v", err)
    }
    fmt.Println("Data Usage History:\n", history)
    return history, nil
}

// StorageAllocateCache allocates cache memory for storage
func StorageAllocateCache(cacheID string, size int64) error {
    if err := common.AllocateCache(cacheID, size); err != nil {
        return fmt.Errorf("cache allocation failed: %v", err)
    }
    ledger.RecordCacheAllocation(cacheID, size, time.Now())
    fmt.Printf("Cache %s allocated with size %d bytes\n", cacheID, size)
    return nil
}

// StorageDeallocateCache deallocates cache memory
func StorageDeallocateCache(cacheID string) error {
    if err := common.DeallocateCache(cacheID); err != nil {
        return fmt.Errorf("cache deallocation failed: %v", err)
    }
    ledger.RecordCacheDeallocation(cacheID, time.Now())
    fmt.Printf("Cache %s deallocated\n", cacheID)
    return nil
}

// StorageMonitorCacheUsage monitors the usage of a specified cache
func StorageMonitorCacheUsage(cacheID string) (float64, error) {
    usage, err := common.MonitorCacheUsage(cacheID)
    if err != nil {
        return 0, fmt.Errorf("cache usage monitoring failed: %v", err)
    }
    fmt.Printf("Cache usage for %s: %.2f%%\n", cacheID, usage)
    return usage, nil
}

// StorageFetchCacheUsageHistory retrieves the usage history of a specific cache
func StorageFetchCacheUsageHistory(cacheID string) ([]string, error) {
    history, err := ledger.FetchCacheUsageHistory(cacheID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch cache usage history: %v", err)
    }
    fmt.Println("Cache Usage History:\n", history)
    return history, nil
}

// StorageSetCacheSize sets the size of a specified cache
func StorageSetCacheSize(cacheID string, size int64) error {
    if err := common.SetCacheSize(cacheID, size); err != nil {
        return fmt.Errorf("setting cache size failed: %v", err)
    }
    ledger.RecordCacheSizeChange(cacheID, size, time.Now())
    fmt.Printf("Cache %s size set to %d bytes\n", cacheID, size)
    return nil
}

// StorageGetCacheSize retrieves the current size of a specified cache
func StorageGetCacheSize(cacheID string) (int64, error) {
    size, err := common.GetCacheSize(cacheID)
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve cache size: %v", err)
    }
    fmt.Printf("Current size of cache %s: %d bytes\n", cacheID, size)
    return size, nil
}

// StorageAuditCache audits a cache for compliance and integrity
func StorageAuditCache(cacheID string) error {
    if !common.VerifyCacheCompliance(cacheID) {
        return fmt.Errorf("cache audit failed for %s", cacheID)
    }
    fmt.Printf("Cache %s passed audit\n", cacheID)
    return nil
}

// StorageAllocateResourcePool allocates resources to a resource pool
func StorageAllocateResourcePool(poolID string, resources int64) error {
    if err := common.AllocateResourcePool(poolID, resources); err != nil {
        return fmt.Errorf("resource pool allocation failed: %v", err)
    }
    ledger.RecordResourcePoolAllocation(poolID, resources, time.Now())
    fmt.Printf("Resource pool %s allocated with resources %d\n", poolID, resources)
    return nil
}

// StorageDeallocateResourcePool deallocates resources from a resource pool
func StorageDeallocateResourcePool(poolID string) error {
    if err := common.DeallocateResourcePool(poolID); err != nil {
        return fmt.Errorf("resource pool deallocation failed: %v", err)
    }
    ledger.RecordResourcePoolDeallocation(poolID, time.Now())
    fmt.Printf("Resource pool %s deallocated\n", poolID)
    return nil
}

// StorageMonitorResourcePoolUsage monitors the usage of a specified resource pool
func StorageMonitorResourcePoolUsage(poolID string) (float64, error) {
    usage, err := common.MonitorResourcePoolUsage(poolID)
    if err != nil {
        return 0, fmt.Errorf("resource pool usage monitoring failed: %v", err)
    }
    fmt.Printf("Resource pool usage for %s: %.2f%%\n", poolID, usage)
    return usage, nil
}

// StorageFetchResourcePoolHistory retrieves the usage history of a resource pool
func StorageFetchResourcePoolHistory(poolID string) ([]string, error) {
    history, err := ledger.FetchResourcePoolHistory(poolID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch resource pool history: %v", err)
    }
    fmt.Println("Resource Pool Usage History:\n", history)
    return history, nil
}

// StorageSetResourcePoolLimit sets a resource limit for a resource pool
func StorageSetResourcePoolLimit(poolID string, limit int64) error {
    if err := common.SetResourcePoolLimit(poolID, limit); err != nil {
        return fmt.Errorf("setting resource pool limit failed: %v", err)
    }
    ledger.RecordResourcePoolLimitSet(poolID, limit, time.Now())
    fmt.Printf("Resource limit set for pool %s to %d\n", poolID, limit)
    return nil
}

// StorageGetResourcePoolLimit retrieves the current resource limit for a resource pool
func StorageGetResourcePoolLimit(poolID string) (int64, error) {
    limit, err := common.GetResourcePoolLimit(poolID)
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve resource pool limit: %v", err)
    }
    fmt.Printf("Current resource pool limit for %s: %d\n", poolID, limit)
    return limit, nil
}

// StorageAuditResourcePoolUsage audits a resource pool for compliance and integrity
func StorageAuditResourcePoolUsage(poolID string) error {
    if !common.VerifyResourcePoolCompliance(poolID) {
        return fmt.Errorf("resource pool audit failed for %s", poolID)
    }
    fmt.Printf("Resource pool %s passed audit\n", poolID)
    return nil
}

// StorageAllocateTemporaryResources allocates temporary resources for short-term use
func StorageAllocateTemporaryResources(resourceID string, amount int64) error {
    if err := common.AllocateTemporaryResources(resourceID, amount); err != nil {
        return fmt.Errorf("temporary resource allocation failed: %v", err)
    }
    ledger.RecordTemporaryResourceAllocation(resourceID, amount, time.Now())
    fmt.Printf("Temporary resources %s allocated with amount %d\n", resourceID, amount)
    return nil
}

// StorageDeallocateTemporaryResources deallocates temporary resources
func StorageDeallocateTemporaryResources(resourceID string) error {
    if err := common.DeallocateTemporaryResources(resourceID); err != nil {
        return fmt.Errorf("temporary resource deallocation failed: %v", err)
    }
    ledger.RecordTemporaryResourceDeallocation(resourceID, time.Now())
    fmt.Printf("Temporary resources %s deallocated\n", resourceID)
    return nil
}

// StorageMonitorTemporaryUsage monitors the usage of temporary resources
func StorageMonitorTemporaryUsage(resourceID string) (float64, error) {
    usage, err := common.MonitorTemporaryUsage(resourceID)
    if err != nil {
        return 0, fmt.Errorf("temporary usage monitoring failed: %v", err)
    }
    fmt.Printf("Temporary resource usage for %s: %.2f%%\n", resourceID, usage)
    return usage, nil
}

// StorageFetchTemporaryUsageHistory retrieves the usage history of temporary resources
func StorageFetchTemporaryUsageHistory(resourceID string) ([]string, error) {
    history, err := ledger.FetchTemporaryUsageHistory(resourceID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch temporary usage history: %v", err)
    }
    fmt.Println("Temporary Resource Usage History:\n", history)
    return history, nil
}

// StorageAuditTemporaryUsage audits temporary resource usage for compliance
func StorageAuditTemporaryUsage(resourceID string) error {
    if !common.VerifyTemporaryUsageCompliance(resourceID) {
        return fmt.Errorf("temporary resource usage audit failed for %s", resourceID)
    }
    fmt.Printf("Temporary resource usage for %s passed audit\n", resourceID)
    return nil
}

// StorageSetResourcePriority sets priority for a specific resource
func StorageSetResourcePriority(resourceID string, priorityLevel int) error {
    if err := common.SetResourcePriority(resourceID, priorityLevel); err != nil {
        return fmt.Errorf("setting resource priority failed: %v", err)
    }
    ledger.RecordResourcePrioritySet(resourceID, priorityLevel, time.Now())
    fmt.Printf("Priority level for resource %s set to %d\n", resourceID, priorityLevel)
    return nil
}

// StorageGetResourcePriority retrieves the current priority level of a resource
func StorageGetResourcePriority(resourceID string) (int, error) {
    priority, err := common.GetResourcePriority(resourceID)
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve resource priority: %v", err)
    }
    fmt.Printf("Current priority level for resource %s: %d\n", resourceID, priority)
    return priority, nil
}
