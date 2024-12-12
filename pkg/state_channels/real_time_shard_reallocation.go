package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewRealTimeShardReallocationChannel initializes a new Real-Time Shard Reallocation Channel (RTSR)
func NewRealTimeShardReallocationChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.RealTimeShardReallocationChannel {
	return &common.RealTimeShardReallocationChannel{
		ChannelID:      channelID,
		Shards:         make(map[string]*common.Shard),
		Participants:   participants,
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
	}
}

// OpenChannel opens the real-time shard reallocation state channel
func (rtsr *common.RealTimeShardReallocationChannel) OpenChannel() error {
	rtsr.mu.Lock()
	defer rtsr.mu.Unlock()

	if rtsr.IsOpen {
		return errors.New("real-time shard reallocation channel is already open")
	}

	rtsr.IsOpen = true

	// Log channel opening in the ledger
	err := rtsr.Ledger.RecordChannelOpening(rtsr.ChannelID, rtsr.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Real-Time Shard Reallocation Channel %s opened with participants: %v\n", rtsr.ChannelID, rtsr.Participants)
	return nil
}

// CloseChannel closes the real-time shard reallocation state channel
func (rtsr *common.RealTimeShardReallocationChannel) CloseChannel() error {
	rtsr.mu.Lock()
	defer rtsr.mu.Unlock()

	if !rtsr.IsOpen {
		return errors.New("real-time shard reallocation channel is already closed")
	}

	rtsr.IsOpen = false

	// Log channel closure in the ledger
	err := rtsr.Ledger.RecordChannelClosure(rtsr.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Real-Time Shard Reallocation Channel %s closed\n", rtsr.ChannelID)
	return nil
}

// AllocateShard allocates a new shard to the channel
func (rtsr *common.RealTimeShardReallocationChannel) AllocateShard(shardID string, data []byte) error {
	rtsr.mu.Lock()
	defer rtsr.mu.Unlock()

	if !rtsr.IsOpen {
		return errors.New("channel is closed")
	}

	if _, exists := rtsr.Shards[shardID]; exists {
		return errors.New("shard already allocated")
	}

	// Encrypt the shard data
	encryptedData, err := rtsr.Encryption.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt shard data: %v", err)
	}

	// Create the shard
	shard := &common.Shard{
		ShardID:   shardID,
		Data:      encryptedData,
		Timestamp: time.Now(),
	}

	// Allocate the shard to the channel
	rtsr.Shards[shardID] = shard

	// Log the shard allocation in the ledger
	err = rtsr.Ledger.RecordShardAllocation(rtsr.ChannelID, shardID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard allocation: %v", err)
	}

	fmt.Printf("Shard %s allocated to channel %s\n", shardID, rtsr.ChannelID)
	return nil
}

// ReallocateShard reallocates a shard to another channel or participant based on network load
func (rtsr *common.RealTimeShardReallocationChannel) ReallocateShard(shardID string, newChannelID string) error {
	rtsr.mu.Lock()
	defer rtsr.mu.Unlock()

	if !rtsr.IsOpen {
		return errors.New("channel is closed")
	}

	shard, exists := rtsr.Shards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found in channel %s", shardID, rtsr.ChannelID)
	}

	// Perform the reallocation logic
	err := rtsr.NetworkManager.ReallocateShard(shardID, rtsr.ChannelID, newChannelID)
	if err != nil {
		return fmt.Errorf("failed to reallocate shard %s: %v", shardID, err)
	}

	// Update the shard status
	shard.Reallocated = true
	shard.Timestamp = time.Now()

	// Log the shard reallocation in the ledger
	err = rtsr.Ledger.RecordShardReallocation(shardID, rtsr.ChannelID, newChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard reallocation: %v", err)
	}

	fmt.Printf("Shard %s reallocated from channel %s to channel %s\n", shardID, rtsr.ChannelID, newChannelID)
	return nil
}

// RetrieveShard retrieves a shard by its ID from the channel
func (rtsr *common.RealTimeShardReallocationChannel) RetrieveShard(shardID string) (*common.Shard, error) {
	rtsr.mu.Lock()
	defer rtsr.mu.Unlock()

	shard, exists := rtsr.Shards[shardID]
	if !exists {
		return nil, fmt.Errorf("shard %s not found in channel %s", shardID, rtsr.ChannelID)
	}

	fmt.Printf("Retrieved shard %s from channel %s\n", shardID, rtsr.ChannelID)
	return shard, nil
}

// MonitorReallocation monitors shards and reallocates them based on network load and congestion
func (rtsr *common.RealTimeShardReallocationChannel) MonitorReallocation(interval time.Duration) {
	for {
		time.Sleep(interval)
		for shardID := range rtsr.Shards {
			// Dummy check for reallocation, in a real system this would be based on metrics like load, performance, etc.
			if rtsr.Shards[shardID].Reallocated {
				continue
			}
			err := rtsr.ReallocateShard(shardID, "some-other-channel-id") // Placeholder for actual channel ID
			if err != nil {
				fmt.Printf("Error reallocating shard %s: %v\n", shardID, err)
				continue
			}
		}
	}
}

// UpdateState securely updates the internal state of the channel
func (rtsr *common.RealTimeShardReallocationChannel) UpdateState(key string, value interface{}) error {
	rtsr.mu.Lock()
	defer rtsr.mu.Unlock()

	if !rtsr.IsOpen {
		return errors.New("state channel is closed")
	}

	// Update the channel state
	rtsr.State[key] = value

	// Log the state update in the ledger
	err := rtsr.Ledger.RecordStateUpdate(rtsr.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of channel %s updated: %s = %v\n", rtsr.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the channel
func (rtsr *common.RealTimeShardReallocationChannel) RetrieveState(key string) (interface{}, error) {
	rtsr.mu.Lock()
	defer rtsr.mu.Unlock()

	if !rtsr.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	value, exists := rtsr.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from channel %s: %s = %v\n", rtsr.ChannelID, key, value)
	return value, nil
}
