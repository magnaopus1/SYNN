package scalability

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewPartitionManager initializes the partitioning manager
func NewPartitionManager(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.PartitionManager {
	return &common.PartitionManager{
		Partitions:       make(map[string]*common.Partition),
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// CreatePartition initializes a new partition with the specified type and tolerance limits
func (pm *common.PartitionManager) CreatePartition(partitionID, partitionType string, toleranceLimit int, partitionData []byte) (*common.Partition, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.Partitions[partitionID]; exists {
		return nil, fmt.Errorf("partition %s already exists", partitionID)
	}

	// Encrypt the data
	encryptedData, err := pm.EncryptionService.EncryptData(partitionData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt partition data: %v", err)
	}

	// Create a new partition
	partition := &common.Partition{
		PartitionID:      partitionID,
		PartitionType:    partitionType,
		Data:             encryptedData,
		LastRebalanced:   time.Now(),
		LastAdjusted:     time.Now(),
		ToleranceLimit:   toleranceLimit,
	}

	pm.Partitions[partitionID] = partition

	// Log the creation in the ledger
	err = pm.Ledger.RecordPartitionCreation(partitionID, partitionType, toleranceLimit, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log partition creation: %v", err)
	}

	fmt.Printf("Partition %s of type %s created with tolerance limit %d\n", partitionID, partitionType, toleranceLimit)
	return partition, nil
}

// RebalancePartition rebalances the load across a partition to ensure it operates within optimal limits
func (pm *common.PartitionManager) RebalancePartition(partitionID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	partition, exists := pm.Partitions[partitionID]
	if !exists {
		return fmt.Errorf("partition %s not found", partitionID)
	}

	// Simulate rebalancing logic
	partition.LastRebalanced = time.Now()

	// Log the rebalancing action
	err := pm.Ledger.RecordPartitionRebalance(partitionID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log partition rebalancing: %v", err)
	}

	fmt.Printf("Partition %s has been rebalanced\n", partitionID)
	return nil
}

// DynamicPartition adjusts a partition's tolerance and parameters based on network performance
func (pm *common.PartitionManager) DynamicPartition(partitionID string, newToleranceLimit int) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	partition, exists := pm.Partitions[partitionID]
	if !exists {
		return fmt.Errorf("partition %s not found", partitionID)
	}

	// Update tolerance limits
	partition.ToleranceLimit = newToleranceLimit
	partition.LastAdjusted = time.Now()

	// Log dynamic partitioning
	err := pm.Ledger.RecordPartitionAdjustment(partitionID, newToleranceLimit, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log partition adjustment: %v", err)
	}

	fmt.Printf("Partition %s dynamically adjusted with new tolerance limit %d\n", partitionID, newToleranceLimit)
	return nil
}

// HorizontalPartition splits data horizontally across multiple partitions
func (pm *common.PartitionManager) HorizontalPartition(partitionID string, dataChunks [][]byte) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	partition, exists := pm.Partitions[partitionID]
	if !exists {
		return fmt.Errorf("partition %s not found", partitionID)
	}

	// Encrypt each data chunk and append it to the partition
	for _, chunk := range dataChunks {
		encryptedChunk, err := pm.EncryptionService.EncryptData(chunk, common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt data chunk: %v", err)
		}
		partition.Data = append(partition.Data, encryptedChunk...)
	}

	// Log horizontal partitioning
	err := pm.Ledger.RecordHorizontalPartitioning(partitionID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log horizontal partitioning: %v", err)
	}

	fmt.Printf("Partition %s horizontally partitioned\n", partitionID)
	return nil
}

// VerticalPartition splits data vertically across multiple partitions
func (pm *common.PartitionManager) VerticalPartition(partitionID string, columns [][]byte) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	partition, exists := pm.Partitions[partitionID]
	if !exists {
		return fmt.Errorf("partition %s not found", partitionID)
	}

	// Encrypt each column of data
	for _, column := range columns {
		encryptedColumn, err := pm.EncryptionService.EncryptData(column, common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt column data: %v", err)
		}
		partition.Data = append(partition.Data, encryptedColumn...)
	}

	// Log vertical partitioning
	err := pm.Ledger.RecordVerticalPartitioning(partitionID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log vertical partitioning: %v", err)
	}

	fmt.Printf("Partition %s vertically partitioned\n", partitionID)
	return nil
}

// EnsurePartitionTolerance checks whether a partition is within its tolerance limits
func (pm *common.PartitionManager) EnsurePartitionTolerance(partitionID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	partition, exists := pm.Partitions[partitionID]
	if !exists {
		return fmt.Errorf("partition %s not found", partitionID)
	}

	// Check if partition exceeds tolerance limit
	if partition.ToleranceLimit > 100 {
		return fmt.Errorf("partition %s exceeds tolerance limit", partitionID)
	}

	fmt.Printf("Partition %s is operating within tolerance limits\n", partitionID)
	return nil
}

// RetrievePartitionData retrieves decrypted data from a partition
func (pm *common.PartitionManager) RetrievePartitionData(partitionID string) ([]byte, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	partition, exists := pm.Partitions[partitionID]
	if !exists {
		return nil, fmt.Errorf("partition %s not found", partitionID)
	}

	// Decrypt the data before returning
	decryptedData, err := pm.EncryptionService.DecryptData(partition.Data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt partition data: %v", err)
	}

	fmt.Printf("Data retrieved from partition %s\n", partitionID)
	return decryptedData, nil
}
