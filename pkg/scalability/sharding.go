package scalability

import (
	"errors"
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)


// NewShardManager initializes the Shard Manager
func NewShardManager(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.ShardManager {
	return &common.ShardManager{
		Shards:           make(map[string]*common.Shard),
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// CreateShard creates a new shard in the network
func (sm *common.ShardManager) CreateShard(shardID, parentShardID string, data []byte) (*common.Shard, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.Shards[shardID]; exists {
		return nil, fmt.Errorf("shard %s already exists", shardID)
	}

	// Encrypt the shard data
	encryptedData, err := sm.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt shard data: %v", err)
	}

	// Create the new shard
	shard := &common.Shard{
		ShardID:         shardID,
		ParentShardID:   parentShardID,
		Data:            encryptedData,
		LastMergedSplit: time.Now(),
		IsAvailable:     true,
	}

	sm.Shards[shardID] = shard

	// Log shard creation
	err = sm.Ledger.RecordShardCreation(shardID, parentShardID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log shard creation: %v", err)
	}

	fmt.Printf("Shard %s created with parent shard %s\n", shardID, parentShardID)
	return shard, nil
}

// CrossShardCommunication enables communication between two shards
func (sm *common.ShardManager) CrossShardCommunication(sourceShardID, targetShardID string, data []byte) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sourceShard, exists := sm.Shards[sourceShardID]
	if !exists {
		return fmt.Errorf("source shard %s not found", sourceShardID)
	}

	targetShard, exists := sm.Shards[targetShardID]
	if !exists {
		return fmt.Errorf("target shard %s not found", targetShardID)
	}

	// Encrypt the communication data
	encryptedData, err := sm.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt cross-shard communication data: %v", err)
	}

	// Simulate communication by storing encrypted data in the target shard
	targetShard.Data = append(targetShard.Data, encryptedData...)

	// Log cross-shard communication
	err = sm.Ledger.RecordCrossShardCommunication(sourceShardID, targetShardID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log cross-shard communication: %v", err)
	}

	fmt.Printf("Communication from shard %s to shard %s completed\n", sourceShardID, targetShardID)
	return nil
}

// HorizontalSharding distributes data across multiple shards horizontally
func (sm *common.ShardManager) HorizontalSharding(parentShardID string, dataChunks [][]byte) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	parentShard, exists := sm.Shards[parentShardID]
	if !exists {
		return fmt.Errorf("parent shard %s not found", parentShardID)
	}

	// Encrypt and distribute the data across shards
	for _, chunk := range dataChunks {
		encryptedChunk, err := sm.EncryptionService.EncryptData(chunk, common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt data chunk: %v", err)
		}
		parentShard.Data = append(parentShard.Data, encryptedChunk...)
	}

	// Log horizontal sharding
	err := sm.Ledger.RecordHorizontalSharding(parentShardID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log horizontal sharding: %v", err)
	}

	fmt.Printf("Horizontal sharding completed for parent shard %s\n", parentShardID)
	return nil
}

// VerticalSharding distributes data across shards vertically (e.g., by columns)
func (sm *common.ShardManager) VerticalSharding(parentShardID string, columns [][]byte) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	parentShard, exists := sm.Shards[parentShardID]
	if !exists {
		return fmt.Errorf("parent shard %s not found", parentShardID)
	}

	// Encrypt and distribute the columns across shards
	for _, column := range columns {
		encryptedColumn, err := sm.EncryptionService.EncryptData(column, common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt column data: %v", err)
		}
		parentShard.Data = append(parentShard.Data, encryptedColumn...)
	}

	// Log vertical sharding
	err := sm.Ledger.RecordVerticalSharding(parentShardID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log vertical sharding: %v", err)
	}

	fmt.Printf("Vertical sharding completed for parent shard %s\n", parentShardID)
	return nil
}

// MergeShards merges two shards into one
func (sm *common.ShardManager) MergeShards(sourceShardID, targetShardID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sourceShard, exists := sm.Shards[sourceShardID]
	if !exists {
		return fmt.Errorf("source shard %s not found", sourceShardID)
	}

	targetShard, exists := sm.Shards[targetShardID]
	if !exists {
		return fmt.Errorf("target shard %s not found", targetShardID)
	}

	// Merge data from source shard to target shard
	targetShard.Data = append(targetShard.Data, sourceShard.Data...)

	// Mark the source shard as unavailable
	sourceShard.IsAvailable = false

	// Log the merge action
	err := sm.Ledger.RecordShardMerge(sourceShardID, targetShardID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard merge: %v", err)
	}

	fmt.Printf("Shard %s merged into shard %s\n", sourceShardID, targetShardID)
	return nil
}

// SplitShard splits a shard into two, creating a new shard with a subset of the data
func (sm *common.ShardManager) SplitShard(sourceShardID, newShardID string, splitRatio float64) (*common.Shard, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sourceShard, exists := sm.Shards[sourceShardID]
	if !exists {
		return nil, fmt.Errorf("source shard %s not found", sourceShardID)
	}

	// Determine the split point based on the split ratio
	splitPoint := int(float64(len(sourceShard.Data)) * splitRatio)

	// Create a new shard with a portion of the data
	newShard := &common.Shard{
		ShardID:         newShardID,
		ParentShardID:   sourceShardID,
		Data:            sourceShard.Data[:splitPoint],
		LastMergedSplit: time.Now(),
		IsAvailable:     true,
	}

	// Adjust data in the original shard
	sourceShard.Data = sourceShard.Data[splitPoint:]

	sm.Shards[newShardID] = newShard

	// Log the split action
	err := sm.Ledger.RecordShardSplit(sourceShardID, newShardID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log shard split: %v", err)
	}

	fmt.Printf("Shard %s split into shard %s\n", sourceShardID, newShardID)
	return newShard, nil
}

// ManageShardAvailability updates the availability status of a shard
func (sm *common.ShardManager) ManageShardAvailability(shardID string, isAvailable bool) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	shard, exists := sm.Shards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	shard.IsAvailable = isAvailable

	// Log the availability change
	err := sm.Ledger.RecordShardAvailabilityChange(shardID, isAvailable, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard availability change: %v", err)
	}

	fmt.Printf("Shard %s availability updated to %v\n", shardID, isAvailable)
	returnnil
}

// RetrieveShardDetails retrieves details of a shard, including its current data and availability status
func (sm *common.ShardManager) RetrieveShardDetails(shardID string) (*common.Shard, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	shard, exists := sm.Shards[shardID]
	if !exists {
		return nil, fmt.Errorf("shard %s not found", shardID)
	}

	fmt.Printf("Shard %s retrieved: Available: %v\n", shardID, shard.IsAvailable)
	return shard, nil
}

// ListAvailableShards returns a list of all shards currently available for operations
func (sm *common.ShardManager) ListAvailableShards() ([]*common.Shard, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var availableShards []*common.Shard
	for _, shard := range sm.Shards {
		if shard.IsAvailable {
			availableShards = append(availableShards, shard)
		}
	}

	if len(availableShards) == 0 {
		return nil, errors.New("no available shards found")
	}

	fmt.Printf("Retrieved %d available shards\n", len(availableShards))
	return availableShards, nil
}

