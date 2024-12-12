package scalability

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewShardReallocationManager initializes the Real-Time Shard Reallocation system
func NewShardReallocationManager(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.ShardReallocationManager {
	return &common.ShardReallocationManager{
		Shards:           make(map[string]*common.Shard),
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// RegisterShard registers a new shard for reallocation in the system
func (srm *common.ShardReallocationManager) RegisterShard(shardID, stateChannel, rollupID string) (*common.Shard, error) {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	if _, exists := srm.Shards[shardID]; exists {
		return nil, fmt.Errorf("shard %s is already registered", shardID)
	}

	// Create the shard and set its initial state
	shard := &common.Shard{
		ShardID:        shardID,
		StateChannel:   stateChannel,
		RollupID:       rollupID,
		LastReallocated: time.Now(),
		IsAvailable:    true,
	}

	srm.Shards[shardID] = shard

	// Log the shard registration in the ledger
	err := srm.Ledger.RecordShardRegistration(shardID, stateChannel, rollupID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log shard registration: %v", err)
	}

	fmt.Printf("Shard %s registered with state channel %s and rollup %s\n", shardID, stateChannel, rollupID)
	return shard, nil
}

// ReallocateShard reallocates a shard to a different state channel and rollup on demand
func (srm *common.ShardReallocationManager) ReallocateShard(shardID, newStateChannel, newRollupID string) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	shard, exists := srm.Shards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	if !shard.IsAvailable {
		return fmt.Errorf("shard %s is not available for reallocation", shardID)
	}

	// Update shard details for reallocation
	shard.StateChannel = newStateChannel
	shard.RollupID = newRollupID
	shard.LastReallocated = time.Now()

	// Log the reallocation in the ledger
	err := srm.Ledger.RecordShardReallocation(shardID, newStateChannel, newRollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard reallocation: %v", err)
	}

	fmt.Printf("Shard %s reallocated to state channel %s and rollup %s\n", shardID, newStateChannel, newRollupID)
	return nil
}

// MarkShardUnavailable marks a shard as unavailable for reallocation (e.g., during maintenance or load balancing)
func (srm *common.ShardReallocationManager) MarkShardUnavailable(shardID string) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	shard, exists := srm.Shards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	shard.IsAvailable = false

	// Log shard unavailability in the ledger
	err := srm.Ledger.RecordShardAvailability(shardID, false, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard unavailability: %v", err)
	}

	fmt.Printf("Shard %s marked as unavailable\n", shardID)
	return nil
}

// MarkShardAvailable marks a shard as available for reallocation
func (srm *common.ShardReallocationManager) MarkShardAvailable(shardID string) error {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	shard, exists := srm.Shards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	shard.IsAvailable = true

	// Log shard availability in the ledger
	err := srm.Ledger.RecordShardAvailability(shardID, true, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard availability: %v", err)
	}

	fmt.Printf("Shard %s marked as available\n", shardID)
	return nil
}

// RetrieveShardDetails retrieves details of a specific shard, including its current state channel and rollup
func (srm *common.ShardReallocationManager) RetrieveShardDetails(shardID string) (*common.Shard, error) {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	shard, exists := srm.Shards[shardID]
	if !exists {
		return nil, fmt.Errorf("shard %s not found", shardID)
	}

	fmt.Printf("Retrieved details for shard %s: StateChannel: %s, RollupID: %s\n", shardID, shard.StateChannel, shard.RollupID)
	return shard, nil
}

// RetrieveAvailableShards returns a list of all shards available for reallocation
func (srm *common.ShardReallocationManager) RetrieveAvailableShards() ([]*common.Shard, error) {
	srm.mu.Lock()
	defer srm.mu.Unlock()

	var availableShards []*common.Shard
	for _, shard := range srm.Shards {
		if shard.IsAvailable {
			availableShards = append(availableShards, shard)
		}
	}

	if len(availableShards) == 0 {
		return nil, errors.New("no shards available for reallocation")
	}

	fmt.Printf("Retrieved %d available shards for reallocation\n", len(availableShards))
	return availableShards, nil
}
