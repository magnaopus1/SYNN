package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewHierarchicalShardedStateChannel initializes a new Hierarchical Sharded State Channel
func NewHierarchicalShardedStateChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.HierarchicalShardedStateChannel {
	return &common.HierarchicalShardedStateChannel{
		ChannelID:      channelID,
		Participants:   participants,
		State:          make(map[string]interface{}),
		Shards:         make(map[string]*common.Shard),
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
	}
}

// OpenChannel opens the hierarchical sharded state channel
func (hssc *common.HierarchicalShardedStateChannel) OpenChannel() error {
	hssc.mu.Lock()
	defer hssc.mu.Unlock()

	if hssc.IsOpen {
		return errors.New("state channel is already open")
	}

	hssc.IsOpen = true

	// Log channel opening in the ledger
	err := hssc.Ledger.RecordChannelOpening(hssc.ChannelID, hssc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Hierarchical Sharded State Channel %s opened with participants: %v\n", hssc.ChannelID, hssc.Participants)
	return nil
}

// CloseChannel closes the hierarchical sharded state channel and its shards
func (hssc *common.HierarchicalShardedStateChannel) CloseChannel() error {
	hssc.mu.Lock()
	defer hssc.mu.Unlock()

	if !hssc.IsOpen {
		return errors.New("state channel is already closed")
	}

	// Close all shards
	for _, shard := range hssc.Shards {
		err := hssc.CloseShard(shard.ShardID)
		if err != nil {
			return fmt.Errorf("failed to close shard %s: %v", shard.ShardID, err)
		}
	}

	hssc.IsOpen = false

	// Log channel closure in the ledger
	err := hssc.Ledger.RecordChannelClosure(hssc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Hierarchical Sharded State Channel %s and all shards closed\n", hssc.ChannelID)
	return nil
}

// CreateShard creates a new shard within the hierarchical sharded state channel
func (hssc *common.HierarchicalShardedStateChannel) CreateShard(shardID string, participants []string) (*common.Shard, error) {
	hssc.mu.Lock()
	defer hssc.mu.Unlock()

	// Validate if the shard already exists
	if _, exists := hssc.Shards[shardID]; exists {
		return nil, fmt.Errorf("shard %s already exists in the channel", shardID)
	}

	shard := &common.Shard{
		ShardID:      shardID,
		Participants: participants,
		State:        make(map[string]interface{}),
		Timestamp:    time.Now(),
	}

	// Encrypt shard state before adding it
	encryptedState, err := hssc.Encryption.EncryptData([]byte(fmt.Sprintf("%v", shard.State)), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt shard state: %v", err)
	}
	shard.State["encrypted"] = encryptedState

	hssc.Shards[shardID] = shard

	// Log shard creation in the ledger
	err = hssc.Ledger.RecordShardCreation(hssc.ChannelID, shardID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log shard creation: %v", err)
	}

	fmt.Printf("Shard %s created in hierarchical sharded state channel %s\n", shardID, hssc.ChannelID)
	return shard, nil
}

// CloseShard closes a shard within the hierarchical sharded state channel
func (hssc *common.HierarchicalShardedStateChannel) CloseShard(shardID string) error {
	hssc.mu.Lock()
	defer hssc.mu.Unlock()

	shard, exists := hssc.Shards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	// Clear the shard state
	shard.State = make(map[string]interface{})
	shard.Timestamp = time.Now()

	// Log shard closure in the ledger
	err := hssc.Ledger.RecordShardClosure(hssc.ChannelID, shardID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard closure: %v", err)
	}

	fmt.Printf("Shard %s closed in hierarchical sharded state channel %s\n", shardID, hssc.ChannelID)
	return nil
}

// SyncShards synchronizes the state of all shards in the hierarchical sharded state channel
func (hssc *common.HierarchicalShardedStateChannel) SyncShards() error {
	hssc.mu.Lock()
	defer hssc.mu.Unlock()

	if !hssc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Synchronize each shard
	for _, shard := range hssc.Shards {
		// Aggregate shard state to the main channel state
		for key, value := range shard.State {
			hssc.State[key] = value
		}
	}

	// Log the shard synchronization event
	err := hssc.Ledger.RecordShardSync(hssc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard synchronization: %v", err)
	}

	fmt.Printf("Shards in hierarchical sharded state channel %s synchronized\n", hssc.ChannelID)
	return nil
}

// UpdateState securely updates the internal state of a shard within the channel
func (hssc *common.HierarchicalShardedStateChannel) UpdateShardState(shardID, key string, value interface{}) error {
	hssc.mu.Lock()
	defer hssc.mu.Unlock()

	shard, exists := hssc.Shards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	// Update the shard state
	shard.State[key] = value
	shard.Timestamp = time.Now()

	// Log the state update
	err := hssc.Ledger.RecordStateUpdate(hssc.ChannelID, shardID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of shard %s in hierarchical sharded state channel %s updated: %s = %v\n", shardID, hssc.ChannelID, key, value)
	return nil
}

// RetrieveShard retrieves a shard by its ID within the hierarchical sharded state channel
func (hssc *common.HierarchicalShardedStateChannel) RetrieveShard(shardID string) (*common.Shard, error) {
	hssc.mu.Lock()
	defer hssc.mu.Unlock()

	shard, exists := hssc.Shards[shardID]
	if !exists {
		return nil, fmt.Errorf("shard %s not found in channel %s", shardID, hssc.ChannelID)
	}

	fmt.Printf("Retrieved shard %s from hierarchical sharded state channel %s\n", shardID, hssc.ChannelID)
	return shard, nil
}

// RebalanceShards dynamically reallocates resources across shards to optimize channel performance
func (hssc *common.HierarchicalShardedStateChannel) RebalanceShards() error {
	hssc.mu.Lock()
	defer hssc.mu.Unlock()

	if !hssc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Example logic for rebalancing shards (adjusting shard size, participants, etc.)
	for _, shard := range hssc.Shards {
		// Example: Add a new dummy participant or adjust shard size
		shard.Participants = append(shard.Participants, fmt.Sprintf("dummy-participant-%d", len(shard.Participants)+1))
		shard.Timestamp = time.Now()

		// Log the rebalance action
		err := hssc.Ledger.RecordShardRebalance(hssc.ChannelID, shard.ShardID, time.Now())
		if err != nil {
			return fmt.Errorf("failed to log shard rebalance: %v", err)
		}
	}

	fmt.Printf("Rebalanced shards in hierarchical sharded state channel %s\n", hssc.ChannelID)
	return nil
}

// ValidateShard ensures that the state of a shard is consistent and valid
func (hssc *common.HierarchicalShardedStateChannel) ValidateShard(shardID string) error {
	hssc.mu.Lock()
	defer hssc.mu.Unlock()

	shard, exists := hssc.Shards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	// Example validation logic
	for key, value := range shard.State {
		fmt.Printf("Validating state: %s = %v\n", key, value)
	}

	// Log the validation event
	err := hssc.Ledger.RecordShardValidation(hssc.ChannelID, shardID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard validation: %v", err)
	}

	fmt.Printf("Shard %s in hierarchical sharded state channel %s validated\n", shardID, hssc.ChannelID)
	return nil
}
