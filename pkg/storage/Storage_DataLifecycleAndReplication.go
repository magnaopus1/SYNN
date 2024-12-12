package storage

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
)

// StorageSetDataLifecyclePolicy sets a lifecycle policy for data storage
func StorageSetDataLifecyclePolicy(policyID string, policyDetails []byte) error {
    encryptedPolicy, err := encryption.EncryptData(policyDetails, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("encryption failed: %v", err)
    }
    if err := common.SaveLifecyclePolicy(policyID, encryptedPolicy); err != nil {
        return fmt.Errorf("failed to set lifecycle policy: %v", err)
    }
    ledger.RecordLifecyclePolicySet(policyID, time.Now())
    fmt.Printf("Lifecycle policy %s set\n", policyID)
    return nil
}

// StorageFetchLifecyclePolicy retrieves a data lifecycle policy by ID
func StorageFetchLifecyclePolicy(policyID string) ([]byte, error) {
    encryptedPolicy, err := common.GetLifecyclePolicy(policyID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch lifecycle policy: %v", err)
    }
    policy, err := encryption.DecryptData(encryptedPolicy, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %v", err)
    }
    fmt.Printf("Lifecycle policy %s retrieved\n", policyID)
    return policy, nil
}

// StorageApplyLifecyclePolicy applies a lifecycle policy to data
func StorageApplyLifecyclePolicy(dataID, policyID string) error {
    if err := common.ApplyPolicyToData(dataID, policyID); err != nil {
        return fmt.Errorf("applying lifecycle policy failed: %v", err)
    }
    ledger.RecordLifecyclePolicyApplication(dataID, policyID, time.Now())
    fmt.Printf("Lifecycle policy %s applied to data %s\n", policyID, dataID)
    return nil
}

// StorageRemoveLifecyclePolicy removes a lifecycle policy from data
func StorageRemoveLifecyclePolicy(dataID, policyID string) error {
    if err := common.RemovePolicyFromData(dataID, policyID); err != nil {
        return fmt.Errorf("removing lifecycle policy failed: %v", err)
    }
    ledger.RecordLifecyclePolicyRemoval(dataID, policyID, time.Now())
    fmt.Printf("Lifecycle policy %s removed from data %s\n", policyID, dataID)
    return nil
}

// StorageAuditLifecyclePolicy audits data lifecycle policy integrity
func StorageAuditLifecyclePolicy(policyID string) error {
    if !common.VerifyLifecyclePolicyIntegrity(policyID) {
        return fmt.Errorf("lifecycle policy audit failed for %s", policyID)
    }
    fmt.Printf("Lifecycle policy %s passed audit\n", policyID)
    return nil
}

// StorageMonitorLifecycleEvents monitors lifecycle events for all data
func StorageMonitorLifecycleEvents() error {
    events, err := common.MonitorLifecycleEvents()
    if err != nil {
        return fmt.Errorf("monitoring lifecycle events failed: %v", err)
    }
    fmt.Printf("Current lifecycle events: %v\n", events)
    return nil
}

// StorageFetchLifecycleEventLog fetches event log for data lifecycle events
func StorageFetchLifecycleEventLog() ([]string, error) {
    eventLog, err := ledger.FetchLifecycleEventLog()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch lifecycle event log: %v", err)
    }
    fmt.Println("Lifecycle Event Log:\n", eventLog)
    return eventLog, nil
}

// StorageOptimizeDataLifecycle optimizes data based on lifecycle policies
func StorageOptimizeDataLifecycle() error {
    if err := common.OptimizeLifecycleData(); err != nil {
        return fmt.Errorf("data lifecycle optimization failed: %v", err)
    }
    fmt.Println("Data lifecycle optimized based on policies")
    return nil
}

// StorageAllocateReplication allocates replication for data
func StorageAllocateReplication(dataID string, replicationFactor int) error {
    if err := common.AllocateDataReplication(dataID, replicationFactor); err != nil {
        return fmt.Errorf("replication allocation failed: %v", err)
    }
    ledger.RecordReplicationAllocation(dataID, replicationFactor, time.Now())
    fmt.Printf("Replication allocated for data %s with factor %d\n", dataID, replicationFactor)
    return nil
}

// StorageDeallocateReplication removes replication for data
func StorageDeallocateReplication(dataID string) error {
    if err := common.DeallocateDataReplication(dataID); err != nil {
        return fmt.Errorf("replication deallocation failed: %v", err)
    }
    ledger.RecordReplicationDeallocation(dataID, time.Now())
    fmt.Printf("Replication deallocated for data %s\n", dataID)
    return nil
}

// StorageMonitorReplicationUsage monitors usage of replicated data storage
func StorageMonitorReplicationUsage() error {
    usage, err := common.MonitorReplicationUsage()
    if err != nil {
        return fmt.Errorf("replication usage monitoring failed: %v", err)
    }
    fmt.Printf("Current replication storage usage: %.2f%%\n", usage)
    return nil
}

// StorageFetchReplicationHistory retrieves the history of replication events
func StorageFetchReplicationHistory() ([]string, error) {
    history, err := ledger.FetchReplicationHistory()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch replication history: %v", err)
    }
    fmt.Println("Replication History:\n", history)
    return history, nil
}

// StorageSetReplicationFactor sets replication factor for data
func StorageSetReplicationFactor(dataID string, factor int) error {
    if err := common.SetReplicationFactor(dataID, factor); err != nil {
        return fmt.Errorf("setting replication factor failed: %v", err)
    }
    ledger.RecordReplicationFactorSet(dataID, factor, time.Now())
    fmt.Printf("Replication factor set for data %s to %d\n", dataID, factor)
    return nil
}

// StorageGetReplicationFactor retrieves the replication factor for data
func StorageGetReplicationFactor(dataID string) (int, error) {
    factor, err := common.GetReplicationFactor(dataID)
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve replication factor: %v", err)
    }
    fmt.Printf("Current replication factor for data %s is %d\n", dataID, factor)
    return factor, nil
}

// StorageAuditReplication audits replication consistency across nodes
func StorageAuditReplication(dataID string) error {
    if !common.VerifyReplicationConsistency(dataID) {
        return fmt.Errorf("replication audit failed for data %s", dataID)
    }
    fmt.Printf("Replication for data %s passed consistency check\n", dataID)
    return nil
}

// StorageOptimizeReplication optimizes replication data for storage efficiency
func StorageOptimizeReplication() error {
    if err := common.OptimizeReplicationStorage(); err != nil {
        return fmt.Errorf("replication optimization failed: %v", err)
    }
    fmt.Println("Replication optimized for storage efficiency")
    return nil
}

// StorageAllocateDataShards allocates data shards for high availability
func StorageAllocateDataShards(dataID string, numShards int) error {
    if err := common.AllocateShards(dataID, numShards); err != nil {
        return fmt.Errorf("data shard allocation failed: %v", err)
    }
    ledger.RecordShardAllocation(dataID, numShards, time.Now())
    fmt.Printf("Allocated %d shards for data %s\n", numShards, dataID)
    return nil
}

// StorageDeallocateDataShards deallocates data shards
func StorageDeallocateDataShards(dataID string) error {
    if err := common.DeallocateShards(dataID); err != nil {
        return fmt.Errorf("shard deallocation failed: %v", err)
    }
    ledger.RecordShardDeallocation(dataID, time.Now())
    fmt.Printf("Deallocated shards for data %s\n", dataID)
    return nil
}

// StorageMonitorShardUsage monitors usage of data shards
func StorageMonitorShardUsage() (float64, error) {
    usage, err := common.MonitorShardUsage()
    if err != nil {
        return 0, fmt.Errorf("shard usage monitoring failed: %v", err)
    }
    fmt.Printf("Current shard usage: %.2f%%\n", usage)
    return usage, nil
}

// StorageFetchShardHistory retrieves the history of shard allocations
func StorageFetchShardHistory() ([]string, error) {
    history, err := ledger.FetchShardHistory()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch shard history: %v", err)
    }
    fmt.Println("Shard Allocation History:\n", history)
    return history, nil
}

// StorageAuditShardUsage audits shard usage for data integrity
func StorageAuditShardUsage(dataID string) error {
    if !common.VerifyShardIntegrity(dataID) {
        return fmt.Errorf("shard usage audit failed for data %s", dataID)
    }
    fmt.Printf("Shards for data %s passed integrity check\n", dataID)
    return nil
}

// StorageOptimizeShardUsage optimizes the distribution and usage of data shards
func StorageOptimizeShardUsage() error {
    if err := common.OptimizeShardDistribution(); err != nil {
        return fmt.Errorf("shard optimization failed: %v", err)
    }
    fmt.Println("Shard distribution optimized for efficiency")
    return nil
}

// StorageSetShardRedundancy sets the redundancy level for data shards
func StorageSetShardRedundancy(dataID string, redundancyLevel int) error {
    if err := common.SetShardRedundancy(dataID, redundancyLevel); err != nil {
        return fmt.Errorf("setting shard redundancy level failed: %v", err)
    }
    ledger.RecordShardRedundancySet(dataID, redundancyLevel, time.Now())
    fmt.Printf("Shard redundancy set for data %s to level %d\n", dataID, redundancyLevel)
    return nil
}

// StorageGetShardRedundancy retrieves the redundancy level of data shards
func StorageGetShardRedundancy(dataID string) (int, error) {
    level, err := common.GetShardRedundancy(dataID)
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve shard redundancy level: %v", err)
    }
    fmt.Printf("Current shard redundancy level for data %s is %d\n", dataID, level)
    return level, nil
}

// StorageTrackRedundantData monitors and tracks redundant data across shards
func StorageTrackRedundantData() (map[string]int, error) {
    redundancyMap, err := common.TrackRedundantData()
    if err != nil {
        return nil, fmt.Errorf("tracking redundant data failed: %v", err)
    }
    fmt.Println("Redundant Data Tracking:\n", redundancyMap)
    return redundancyMap, nil
}

// StorageFetchRedundantDataLog retrieves logs related to redundant data
func StorageFetchRedundantDataLog() ([]string, error) {
    log, err := ledger.FetchRedundantDataLog()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch redundant data log: %v", err)
    }
    fmt.Println("Redundant Data Log:\n", log)
    return log, nil
}
