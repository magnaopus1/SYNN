package storage

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
)

// StorageAllocateColdStorage allocates space in cold storage with encryption
func StorageAllocateColdStorage(data []byte, location string) error {
    encryptedData, err := encryption.EncryptData(data, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("encryption failed: %v", err)
    }
    if err := common.SaveToColdStorage(location, encryptedData); err != nil {
        return fmt.Errorf("cold storage allocation failed: %v", err)
    }
    ledger.RecordColdStorageAllocation(location, time.Now())
    fmt.Printf("Data allocated in cold storage at %s\n", location)
    return nil
}

// StorageDeallocateColdStorage removes data from cold storage
func StorageDeallocateColdStorage(location string) error {
    if err := common.DeleteFromColdStorage(location); err != nil {
        return fmt.Errorf("cold storage deallocation failed: %v", err)
    }
    ledger.RecordColdStorageDeallocation(location, time.Now())
    fmt.Printf("Cold storage deallocated at %s\n", location)
    return nil
}

// StorageMonitorColdStorageUsage monitors cold storage utilization
func StorageMonitorColdStorageUsage() (float64, error) {
    usage, err := common.MonitorColdStorageUsage()
    if err != nil {
        return 0, fmt.Errorf("cold storage usage monitoring failed: %v", err)
    }
    fmt.Printf("Current cold storage usage: %.2f%%\n", usage)
    return usage, nil
}

// StorageFetchColdStorageHistory retrieves the allocation history for cold storage
func StorageFetchColdStorageHistory() ([]string, error) {
    history, err := ledger.FetchColdStorageHistory()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch cold storage history: %v", err)
    }
    fmt.Println("Cold Storage History:\n", history)
    return history, nil
}

// StorageAuditColdStorage audits data integrity in cold storage
func StorageAuditColdStorage(location string) error {
    if !common.VerifyColdStorageIntegrity(location) {
        return fmt.Errorf("cold storage integrity check failed for %s", location)
    }
    fmt.Printf("Cold storage at %s passed integrity check\n", location)
    return nil
}

// StorageSetColdStoragePolicy defines policies for cold storage usage
func StorageSetColdStoragePolicy(policy string) {
    common.SetColdStoragePolicy(policy)
    fmt.Printf("Cold storage policy set: %s\n", policy)
}

// StorageGetColdStoragePolicy retrieves the current cold storage policy
func StorageGetColdStoragePolicy() string {
    policy := common.GetColdStoragePolicy()
    fmt.Printf("Current cold storage policy: %s\n", policy)
    return policy
}

// StorageOptimizeColdStorage performs optimization for cold storage usage
func StorageOptimizeColdStorage() error {
    if err := common.OptimizeColdStorage(); err != nil {
        return fmt.Errorf("cold storage optimization failed: %v", err)
    }
    fmt.Println("Cold storage optimized for better efficiency")
    return nil
}

// StorageAllocateHotStorage allocates space in hot storage with encryption
func StorageAllocateHotStorage(data []byte, location string) error {
    encryptedData, err := encryption.EncryptData(data, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("encryption failed: %v", err)
    }
    if err := common.SaveToHotStorage(location, encryptedData); err != nil {
        return fmt.Errorf("hot storage allocation failed: %v", err)
    }
    ledger.RecordHotStorageAllocation(location, time.Now())
    fmt.Printf("Data allocated in hot storage at %s\n", location)
    return nil
}

// StorageDeallocateHotStorage removes data from hot storage
func StorageDeallocateHotStorage(location string) error {
    if err := common.DeleteFromHotStorage(location); err != nil {
        return fmt.Errorf("hot storage deallocation failed: %v", err)
    }
    ledger.RecordHotStorageDeallocation(location, time.Now())
    fmt.Printf("Hot storage deallocated at %s\n", location)
    return nil
}

// StorageMonitorHotStorageUsage monitors hot storage utilization
func StorageMonitorHotStorageUsage() (float64, error) {
    usage, err := common.MonitorHotStorageUsage()
    if err != nil {
        return 0, fmt.Errorf("hot storage usage monitoring failed: %v", err)
    }
    fmt.Printf("Current hot storage usage: %.2f%%\n", usage)
    return usage, nil
}

// StorageFetchHotStorageHistory retrieves the allocation history for hot storage
func StorageFetchHotStorageHistory() ([]string, error) {
    history, err := ledger.FetchHotStorageHistory()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch hot storage history: %v", err)
    }
    fmt.Println("Hot Storage History:\n", history)
    return history, nil
}

// StorageAuditHotStorage audits data integrity in hot storage
func StorageAuditHotStorage(location string) error {
    if !common.VerifyHotStorageIntegrity(location) {
        return fmt.Errorf("hot storage integrity check failed for %s", location)
    }
    fmt.Printf("Hot storage at %s passed integrity check\n", location)
    return nil
}

// StorageSetHotStoragePolicy defines policies for hot storage usage
func StorageSetHotStoragePolicy(policy string) {
    common.SetHotStoragePolicy(policy)
    fmt.Printf("Hot storage policy set: %s\n", policy)
}

// StorageGetHotStoragePolicy retrieves the current hot storage policy
func StorageGetHotStoragePolicy() string {
    policy := common.GetHotStoragePolicy()
    fmt.Printf("Current hot storage policy: %s\n", policy)
    return policy
}

// StorageOptimizeHotStorage performs optimization for hot storage usage
func StorageOptimizeHotStorage() error {
    if err := common.OptimizeHotStorage(); err != nil {
        return fmt.Errorf("hot storage optimization failed: %v", err)
    }
    fmt.Println("Hot storage optimized for better efficiency")
    return nil
}

// StorageSetEnergyConsumptionLimit sets a limit for storage energy consumption
func StorageSetEnergyConsumptionLimit(limit float64) {
    common.SetEnergyConsumptionLimit(limit)
    fmt.Printf("Energy consumption limit set to %.2f kWh\n", limit)
}

// StorageGetEnergyConsumptionLimit retrieves the current energy consumption limit
func StorageGetEnergyConsumptionLimit() float64 {
    limit := common.GetEnergyConsumptionLimit()
    fmt.Printf("Current energy consumption limit: %.2f kWh\n", limit)
    return limit
}

// StorageMonitorEnergyUsage monitors energy usage for storage operations
func StorageMonitorEnergyUsage() (float64, error) {
    usage, err := common.MonitorEnergyUsage()
    if err != nil {
        return 0, fmt.Errorf("energy usage monitoring failed: %v", err)
    }
    fmt.Printf("Current energy usage: %.2f kWh\n", usage)
    return usage, nil
}

// StorageAuditEnergyUsage audits energy usage to ensure compliance with limits
func StorageAuditEnergyUsage() error {
    if !common.VerifyEnergyUsageCompliance() {
        return fmt.Errorf("energy usage audit failed: limit exceeded")
    }
    fmt.Println("Energy usage within set limits")
    return nil
}

// StorageOptimizeEnergyUsage optimizes energy consumption for storage
func StorageOptimizeEnergyUsage() error {
    if err := common.OptimizeEnergyConsumption(); err != nil {
        return fmt.Errorf("energy optimization failed: %v", err)
    }
    fmt.Println("Energy consumption optimized for storage operations")
    return nil
}

// StorageTrackCostEfficiency monitors cost efficiency for storage operations
func StorageTrackCostEfficiency() (float64, error) {
    efficiency, err := common.CalculateCostEfficiency()
    if err != nil {
        return 0, fmt.Errorf("cost efficiency tracking failed: %v", err)
    }
    fmt.Printf("Current cost efficiency: %.2f\n", efficiency)
    return efficiency, nil
}
