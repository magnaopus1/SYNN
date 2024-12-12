package storage

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
)

// StorageAllocateVirtualDisk allocates a new virtual disk with a specified size
func StorageAllocateVirtualDisk(diskID string, size int64) error {
    if err := common.AllocateVirtualDisk(diskID, size); err != nil {
        return fmt.Errorf("virtual disk allocation failed: %v", err)
    }
    ledger.RecordVirtualDiskAllocation(diskID, size, time.Now())
    fmt.Printf("Virtual disk %s allocated with size %d bytes\n", diskID, size)
    return nil
}

// StorageDeallocateVirtualDisk deallocates a virtual disk
func StorageDeallocateVirtualDisk(diskID string) error {
    if err := common.DeallocateVirtualDisk(diskID); err != nil {
        return fmt.Errorf("virtual disk deallocation failed: %v", err)
    }
    ledger.RecordVirtualDiskDeallocation(diskID, time.Now())
    fmt.Printf("Virtual disk %s deallocated\n", diskID)
    return nil
}

// StorageMonitorVirtualDiskUsage monitors usage of a specified virtual disk
func StorageMonitorVirtualDiskUsage(diskID string) (float64, error) {
    usage, err := common.MonitorVirtualDiskUsage(diskID)
    if err != nil {
        return 0, fmt.Errorf("virtual disk usage monitoring failed: %v", err)
    }
    fmt.Printf("Current usage for virtual disk %s: %.2f%%\n", diskID, usage)
    return usage, nil
}

// StorageFetchVirtualDiskHistory retrieves the usage history for a virtual disk
func StorageFetchVirtualDiskHistory(diskID string) ([]string, error) {
    history, err := ledger.FetchVirtualDiskHistory(diskID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch virtual disk history: %v", err)
    }
    fmt.Println("Virtual Disk Usage History:\n", history)
    return history, nil
}

// StorageOptimizeVirtualDisk optimizes storage allocation on a virtual disk
func StorageOptimizeVirtualDisk(diskID string) error {
    if err := common.OptimizeVirtualDisk(diskID); err != nil {
        return fmt.Errorf("virtual disk optimization failed: %v", err)
    }
    fmt.Printf("Virtual disk %s optimized for storage efficiency\n", diskID)
    return nil
}

// StorageSetVirtualDiskSize sets the size of a virtual disk
func StorageSetVirtualDiskSize(diskID string, size int64) error {
    if err := common.SetVirtualDiskSize(diskID, size); err != nil {
        return fmt.Errorf("setting virtual disk size failed: %v", err)
    }
    ledger.RecordVirtualDiskSizeChange(diskID, size, time.Now())
    fmt.Printf("Virtual disk %s size set to %d bytes\n", diskID, size)
    return nil
}

// StorageFetchVirtualDiskSize fetches the current size of a virtual disk
func StorageFetchVirtualDiskSize(diskID string) (int64, error) {
    size, err := common.GetVirtualDiskSize(diskID)
    if err != nil {
        return 0, fmt.Errorf("failed to fetch virtual disk size: %v", err)
    }
    fmt.Printf("Current size for virtual disk %s: %d bytes\n", diskID, size)
    return size, nil
}

// StorageEnableVirtualDiskCompression enables compression for a virtual disk
func StorageEnableVirtualDiskCompression(diskID string) error {
    if err := common.EnableDiskCompression(diskID); err != nil {
        return fmt.Errorf("enabling disk compression failed: %v", err)
    }
    ledger.RecordDiskCompressionEnabled(diskID, time.Now())
    fmt.Printf("Compression enabled for virtual disk %s\n", diskID)
    return nil
}

// StorageDisableVirtualDiskCompression disables compression for a virtual disk
func StorageDisableVirtualDiskCompression(diskID string) error {
    if err := common.DisableDiskCompression(diskID); err != nil {
        return fmt.Errorf("disabling disk compression failed: %v", err)
    }
    ledger.RecordDiskCompressionDisabled(diskID, time.Now())
    fmt.Printf("Compression disabled for virtual disk %s\n", diskID)
    return nil
}

// StorageFetchVirtualDiskCompressionStatus fetches the compression status of a virtual disk
func StorageFetchVirtualDiskCompressionStatus(diskID string) bool {
    status := common.CheckDiskCompressionStatus(diskID)
    fmt.Printf("Compression status for virtual disk %s: %v\n", diskID, status)
    return status
}

// StorageAllocatePartitions allocates partitions within a virtual disk
func StorageAllocatePartitions(diskID string, numPartitions int, partitionSize int64) error {
    if err := common.AllocatePartitions(diskID, numPartitions, partitionSize); err != nil {
        return fmt.Errorf("partition allocation failed: %v", err)
    }
    ledger.RecordPartitionAllocation(diskID, numPartitions, partitionSize, time.Now())
    fmt.Printf("Allocated %d partitions on disk %s, each with size %d bytes\n", numPartitions, diskID, partitionSize)
    return nil
}

// StorageDeallocatePartitions deallocates partitions within a virtual disk
func StorageDeallocatePartitions(diskID string) error {
    if err := common.DeallocatePartitions(diskID); err != nil {
        return fmt.Errorf("partition deallocation failed: %v", err)
    }
    ledger.RecordPartitionDeallocation(diskID, time.Now())
    fmt.Printf("Partitions deallocated on disk %s\n", diskID)
    return nil
}

// StorageMonitorPartitionUsage monitors usage of a specific partition
func StorageMonitorPartitionUsage(diskID, partitionID string) (float64, error) {
    usage, err := common.MonitorPartitionUsage(diskID, partitionID)
    if err != nil {
        return 0, fmt.Errorf("partition usage monitoring failed: %v", err)
    }
    fmt.Printf("Current usage for partition %s on disk %s: %.2f%%\n", partitionID, diskID, usage)
    return usage, nil
}

// StorageFetchPartitionHistory retrieves the usage history for partitions on a virtual disk
func StorageFetchPartitionHistory(diskID string) ([]string, error) {
    history, err := ledger.FetchPartitionHistory(diskID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch partition history: %v", err)
    }
    fmt.Println("Partition Usage History:\n", history)
    return history, nil
}

// StorageSetPartitionSize sets the size of a specified partition
func StorageSetPartitionSize(diskID, partitionID string, size int64) error {
    if err := common.SetPartitionSize(diskID, partitionID, size); err != nil {
        return fmt.Errorf("setting partition size failed: %v", err)
    }
    ledger.RecordPartitionSizeChange(diskID, partitionID, size, time.Now())
    fmt.Printf("Partition %s on disk %s size set to %d bytes\n", partitionID, diskID, size)
    return nil
}

// StorageGetPartitionSize fetches the current size of a specified partition
func StorageGetPartitionSize(diskID, partitionID string) (int64, error) {
    size, err := common.GetPartitionSize(diskID, partitionID)
    if err != nil {
        return 0, fmt.Errorf("failed to fetch partition size: %v", err)
    }
    fmt.Printf("Current size for partition %s on disk %s: %d bytes\n", partitionID, diskID, size)
    return size, nil
}

// StorageAuditPartitions audits partitions within a virtual disk for compliance
func StorageAuditPartitions(diskID string) error {
    if !common.VerifyPartitionCompliance(diskID) {
        return fmt.Errorf("partition audit failed for disk %s", diskID)
    }
    fmt.Printf("Partitions on disk %s passed audit\n", diskID)
    return nil
}

// StorageOptimizePartitionUsage optimizes usage and space within partitions
func StorageOptimizePartitionUsage(diskID string) error {
    if err := common.OptimizePartitionUsage(diskID); err != nil {
        return fmt.Errorf("partition usage optimization failed: %v", err)
    }
    fmt.Printf("Partition usage optimized on disk %s\n", diskID)
    return nil
}

// StorageSetPartitionEncryption enables encryption for a specified partition
func StorageSetPartitionEncryption(diskID, partitionID string) error {
    if err := common.EnablePartitionEncryption(diskID, partitionID); err != nil {
        return fmt.Errorf("partition encryption failed: %v", err)
    }
    ledger.RecordPartitionEncryptionEnabled(diskID, partitionID, time.Now())
    fmt.Printf("Encryption enabled for partition %s on disk %s\n", partitionID, diskID)
    return nil
}

// StorageRemovePartitionEncryption disables encryption for a specified partition
func StorageRemovePartitionEncryption(diskID, partitionID string) error {
    if err := common.DisablePartitionEncryption(diskID, partitionID); err != nil {
        return fmt.Errorf("disabling partition encryption failed: %v", err)
    }
    ledger.RecordPartitionEncryptionDisabled(diskID, partitionID, time.Now())
    fmt.Printf("Encryption disabled for partition %s on disk %s\n", partitionID, diskID)
    return nil
}

// StorageFetchPartitionEncryptionStatus retrieves the encryption status for a partition
func StorageFetchPartitionEncryptionStatus(diskID, partitionID string) bool {
    status := common.CheckPartitionEncryptionStatus(diskID, partitionID)
    fmt.Printf("Encryption status for partition %s on disk %s: %v\n", partitionID, diskID, status)
    return status
}

// StorageAllocateTape allocates a tape storage for long-term backup
func StorageAllocateTape(tapeID string, size int64) error {
    if err := common.AllocateTape(tapeID, size); err != nil {
        return fmt.Errorf("tape allocation failed: %v", err)
    }
    ledger.RecordTapeAllocation(tapeID, size, time.Now())
    fmt.Printf("Tape %s allocated with size %d bytes\n", tapeID, size)
    return nil
}

// StorageDeallocateTape deallocates tape storage
func StorageDeallocateTape(tapeID string) error {
    if err := common.DeallocateTape(tapeID); err != nil {
        return fmt.Errorf("tape deallocation failed: %v", err)
    }
    ledger.RecordTapeDeallocation(tapeID, time.Now())
    fmt.Printf("Tape %s deallocated\n", tapeID)
    return nil
}

// StorageMonitorTapeUsage monitors the usage of a specific tape storage
func StorageMonitorTapeUsage(tapeID string) (float64, error) {
    usage, err := common.MonitorTapeUsage(tapeID)
    if err != nil {
        return 0, fmt.Errorf("tape usage monitoring failed: %v", err)
    }
    fmt.Printf("Current usage for tape %s: %.2f%%\n", tapeID, usage)
    return usage, nil
}

// StorageFetchTapeUsageHistory retrieves the usage history of a specific tape
func StorageFetchTapeUsageHistory(tapeID string) ([]string, error) {
    history, err := ledger.FetchTapeUsageHistory(tapeID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch tape usage history: %v", err)
    }
    fmt.Println("Tape Usage History:\n", history)
    return history, nil
}

// StorageAuditTapeUsage audits the tape usage for compliance and integrity
func StorageAuditTapeUsage(tapeID string) error {
    if !common.VerifyTapeUsageCompliance(tapeID) {
        return fmt.Errorf("tape usage audit failed for tape %s", tapeID)
    }
    fmt.Printf("Tape %s passed usage audit\n", tapeID)
    return nil
}
