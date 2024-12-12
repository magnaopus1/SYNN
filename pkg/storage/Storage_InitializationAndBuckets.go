package storage

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
)

// StorageInit initializes the storage system
func StorageInit() error {
    if err := common.InitializeStorage(); err != nil {
        return fmt.Errorf("storage initialization failed: %v", err)
    }
    ledger.RecordStorageInitialization(time.Now())
    fmt.Println("Storage system initialized successfully")
    return nil
}

// StorageCreateBucket creates a new storage bucket
func StorageCreateBucket(bucketName string) error {
    if err := common.CreateBucket(bucketName); err != nil {
        return fmt.Errorf("bucket creation failed: %v", err)
    }
    ledger.RecordBucketCreation(bucketName, time.Now())
    fmt.Printf("Bucket %s created successfully\n", bucketName)
    return nil
}

// StorageDeleteBucket deletes a storage bucket
func StorageDeleteBucket(bucketName string) error {
    if err := common.DeleteBucket(bucketName); err != nil {
        return fmt.Errorf("bucket deletion failed: %v", err)
    }
    ledger.RecordBucketDeletion(bucketName, time.Now())
    fmt.Printf("Bucket %s deleted successfully\n", bucketName)
    return nil
}

// StorageListBuckets lists all available storage buckets
func StorageListBuckets() ([]string, error) {
    buckets, err := common.ListBuckets()
    if err != nil {
        return nil, fmt.Errorf("failed to list buckets: %v", err)
    }
    fmt.Printf("Buckets available: %v\n", buckets)
    return buckets, nil
}

// StorageAllocateSpace allocates storage space in a specified bucket
func StorageAllocateSpace(bucketName string, space int64) error {
    if err := common.AllocateSpace(bucketName, space); err != nil {
        return fmt.Errorf("space allocation failed: %v", err)
    }
    ledger.RecordSpaceAllocation(bucketName, space, time.Now())
    fmt.Printf("Allocated %d bytes in bucket %s\n", space, bucketName)
    return nil
}

// StorageDeallocateSpace deallocates storage space from a specified bucket
func StorageDeallocateSpace(bucketName string, space int64) error {
    if err := common.DeallocateSpace(bucketName, space); err != nil {
        return fmt.Errorf("space deallocation failed: %v", err)
    }
    ledger.RecordSpaceDeallocation(bucketName, space, time.Now())
    fmt.Printf("Deallocated %d bytes in bucket %s\n", space, bucketName)
    return nil
}

// StorageIncreaseQuota increases the storage quota for a bucket
func StorageIncreaseQuota(bucketName string, increment int64) error {
    if err := common.IncreaseQuota(bucketName, increment); err != nil {
        return fmt.Errorf("quota increase failed: %v", err)
    }
    ledger.RecordQuotaIncrease(bucketName, increment, time.Now())
    fmt.Printf("Increased quota for bucket %s by %d bytes\n", bucketName, increment)
    return nil
}

// StorageDecreaseQuota decreases the storage quota for a bucket
func StorageDecreaseQuota(bucketName string, decrement int64) error {
    if err := common.DecreaseQuota(bucketName, decrement); err != nil {
        return fmt.Errorf("quota decrease failed: %v", err)
    }
    ledger.RecordQuotaDecrease(bucketName, decrement, time.Now())
    fmt.Printf("Decreased quota for bucket %s by %d bytes\n", bucketName, decrement)
    return nil
}

// StorageFetchQuota retrieves the current quota for a specified bucket
func StorageFetchQuota(bucketName string) (int64, error) {
    quota, err := common.GetQuota(bucketName)
    if err != nil {
        return 0, fmt.Errorf("failed to fetch quota: %v", err)
    }
    fmt.Printf("Current quota for bucket %s: %d bytes\n", bucketName, quota)
    return quota, nil
}

// StorageSetQuotaLimit sets a quota limit for a specified bucket
func StorageSetQuotaLimit(bucketName string, limit int64) error {
    if err := common.SetQuotaLimit(bucketName, limit); err != nil {
        return fmt.Errorf("setting quota limit failed: %v", err)
    }
    ledger.RecordQuotaLimitSet(bucketName, limit, time.Now())
    fmt.Printf("Quota limit set for bucket %s to %d bytes\n", bucketName, limit)
    return nil
}

// StorageMonitorQuotaUsage monitors quota usage for a bucket
func StorageMonitorQuotaUsage(bucketName string) (float64, error) {
    usage, err := common.MonitorQuotaUsage(bucketName)
    if err != nil {
        return 0, fmt.Errorf("quota usage monitoring failed: %v", err)
    }
    fmt.Printf("Quota usage for bucket %s: %.2f%%\n", bucketName, usage)
    return usage, nil
}

// StorageTrackSpaceUsage tracks space usage over time in a bucket
func StorageTrackSpaceUsage(bucketName string) error {
    if err := common.TrackSpaceUsage(bucketName); err != nil {
        return fmt.Errorf("tracking space usage failed: %v", err)
    }
    ledger.RecordSpaceUsageTracking(bucketName, time.Now())
    fmt.Printf("Tracking space usage in bucket %s\n", bucketName)
    return nil
}

// StorageFetchSpaceUsageHistory retrieves the space usage history for a bucket
func StorageFetchSpaceUsageHistory(bucketName string) ([]string, error) {
    history, err := ledger.FetchSpaceUsageHistory(bucketName)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch space usage history: %v", err)
    }
    fmt.Println("Space Usage History:\n", history)
    return history, nil
}

// StorageSetAutoExpansion enables auto-expansion for a bucket's storage
func StorageSetAutoExpansion(bucketName string) error {
    if err := common.EnableAutoExpansion(bucketName); err != nil {
        return fmt.Errorf("auto-expansion setting failed: %v", err)
    }
    ledger.RecordAutoExpansionEnabled(bucketName, time.Now())
    fmt.Printf("Auto-expansion enabled for bucket %s\n", bucketName)
    return nil
}

// StorageDisableAutoExpansion disables auto-expansion for a bucket's storage
func StorageDisableAutoExpansion(bucketName string) error {
    if err := common.DisableAutoExpansion(bucketName); err != nil {
        return fmt.Errorf("disabling auto-expansion failed: %v", err)
    }
    ledger.RecordAutoExpansionDisabled(bucketName, time.Now())
    fmt.Printf("Auto-expansion disabled for bucket %s\n", bucketName)
    return nil
}

// StorageAuditUsage audits storage usage for a bucket
func StorageAuditUsage(bucketName string) error {
    if !common.VerifyStorageUsageCompliance(bucketName) {
        return fmt.Errorf("storage usage audit failed for bucket %s", bucketName)
    }
    fmt.Printf("Storage usage for bucket %s passed audit\n", bucketName)
    return nil
}

// StorageSetRetentionPolicy sets a retention policy for data in a bucket
func StorageSetRetentionPolicy(bucketName, policy string) error {
    if err := common.SetRetentionPolicy(bucketName, policy); err != nil {
        return fmt.Errorf("setting retention policy failed: %v", err)
    }
    ledger.RecordRetentionPolicySet(bucketName, policy, time.Now())
    fmt.Printf("Retention policy set for bucket %s\n", bucketName)
    return nil
}

// StorageFetchRetentionPolicy retrieves the retention policy for a bucket
func StorageFetchRetentionPolicy(bucketName string) (string, error) {
    policy, err := common.GetRetentionPolicy(bucketName)
    if err != nil {
        return "", fmt.Errorf("failed to fetch retention policy: %v", err)
    }
    fmt.Printf("Retention policy for bucket %s: %s\n", bucketName, policy)
    return policy, nil
}

// StorageApplyRetentionPolicy applies the retention policy to data in a bucket
func StorageApplyRetentionPolicy(bucketName string) error {
    if err := common.ApplyRetentionPolicy(bucketName); err != nil {
        return fmt.Errorf("applying retention policy failed: %v", err)
    }
    ledger.RecordRetentionPolicyApplied(bucketName, time.Now())
    fmt.Printf("Retention policy applied for bucket %s\n", bucketName)
    return nil
}

// StorageRemoveRetentionPolicy removes the retention policy from a bucket
func StorageRemoveRetentionPolicy(bucketName string) error {
    if err := common.RemoveRetentionPolicy(bucketName); err != nil {
        return fmt.Errorf("removing retention policy failed: %v", err)
    }
    ledger.RecordRetentionPolicyRemoved(bucketName, time.Now())
    fmt.Printf("Retention policy removed from bucket %s\n", bucketName)
    return nil
}

// StorageCompressData compresses data within a bucket
func StorageCompressData(bucketName string, data []byte) ([]byte, error) {
    compressedData, err := common.CompressData(data)
    if err != nil {
        return nil, fmt.Errorf("data compression failed: %v", err)
    }
    ledger.RecordDataCompression(bucketName, time.Now())
    fmt.Printf("Data compressed for bucket %s\n", bucketName)
    return compressedData, nil
}

// StorageDecompressData decompresses data within a bucket
func StorageDecompressData(bucketName string, data []byte) ([]byte, error) {
    decompressedData, err := common.DecompressData(data)
    if err != nil {
        return nil, fmt.Errorf("data decompression failed: %v", err)
    }
    fmt.Printf("Data decompressed for bucket %s\n", bucketName)
    return decompressedData, nil
}

// StorageEnableCompression enables compression for a bucket
func StorageEnableCompression(bucketName string) error {
    if err := common.EnableBucketCompression(bucketName); err != nil {
        return fmt.Errorf("enabling compression failed: %v", err)
    }
    ledger.RecordCompressionEnabled(bucketName, time.Now())
    fmt.Printf("Compression enabled for bucket %s\n", bucketName)
    return nil
}

// StorageDisableCompression disables compression for a bucket
func StorageDisableCompression(bucketName string) error {
    if err := common.DisableBucketCompression(bucketName); err != nil {
        return fmt.Errorf("disabling compression failed: %v", err)
    }
    ledger.RecordCompressionDisabled(bucketName, time.Now())
    fmt.Printf("Compression disabled for bucket %s\n", bucketName)
    return nil
}

// StorageFetchCompressionStatus checks if compression is enabled for a bucket
func StorageFetchCompressionStatus(bucketName string) bool {
    status := common.CheckCompressionStatus(bucketName)
    fmt.Printf("Compression status for bucket %s: %v\n", bucketName, status)
    return status
}
