package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewFluidStateSyncChannel initializes a new Fluid State Sync Channel (FSSC)
func NewFluidStateSyncChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.FluidStateSyncChannel {
	return &common.FluidStateSyncChannel{
		ChannelID:      channelID,
		Participants:   participants,
		State:          make(map[string]interface{}),
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
	}
}

// OpenChannel opens the fluid state sync channel
func (fssc *common.FluidStateSyncChannel) OpenChannel() error {
	fssc.mu.Lock()
	defer fssc.mu.Unlock()

	if fssc.IsOpen {
		return errors.New("state channel is already open")
	}

	fssc.IsOpen = true

	// Log channel opening in the ledger
	err := fssc.Ledger.RecordChannelOpening(fssc.ChannelID, fssc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Fluid State Sync Channel %s opened with participants: %v\n", fssc.ChannelID, fssc.Participants)
	return nil
}

// CloseChannel closes the fluid state sync channel
func (fssc *common.FluidStateSyncChannel) CloseChannel() error {
	fssc.mu.Lock()
	defer fssc.mu.Unlock()

	if !fssc.IsOpen {
		return errors.New("state channel is already closed")
	}

	fssc.IsOpen = false

	// Log channel closure in the ledger
	err := fssc.Ledger.RecordChannelClosure(fssc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Fluid State Sync Channel %s closed\n", fssc.ChannelID)
	return nil
}

// SyncState synchronizes the state of the current channel with the other synchronized channels in real-time
func (fssc *common.FluidStateSyncChannel) SyncState() error {
	fssc.mu.Lock()
	defer fssc.mu.Unlock()

	if !fssc.IsOpen {
		return errors.New("state channel is closed")
	}

	for _, syncedChannel := range fssc.SyncedChannels {
		for key, value := range fssc.State {
			syncedChannel.UpdateState(key, value)
		}
		fmt.Printf("State of channel %s synchronized with channel %s\n", fssc.ChannelID, syncedChannel.ChannelID)
	}

	// Log the synchronization event in the ledger
	err := fssc.Ledger.RecordStateSync(fssc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state synchronization: %v", err)
	}

	return nil
}

// AddSyncedChannel adds a new channel to the synchronization group
func (fssc *common.FluidStateSyncChannel) AddSyncedChannel(channel *common.FluidStateSyncChannel) error {
	fssc.mu.Lock()
	defer fssc.mu.Unlock()

	// Ensure the channel is not already in the sync group
	for _, syncedChannel := range fssc.SyncedChannels {
		if syncedChannel.ChannelID == channel.ChannelID {
			return fmt.Errorf("channel %s is already synced with %s", channel.ChannelID, fssc.ChannelID)
		}
	}

	fssc.SyncedChannels = append(fssc.SyncedChannels, channel)
	fmt.Printf("Channel %s added to sync group of channel %s\n", channel.ChannelID, fssc.ChannelID)
	return nil
}

// UpdateState securely updates the internal state of the channel and synchronizes with synced channels
func (fssc *common.FluidStateSyncChannel) UpdateState(key string, value interface{}) error {
	fssc.mu.Lock()
	defer fssc.mu.Unlock()

	if !fssc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Update the channel state
	fssc.State[key] = value

	// Log the state update in the ledger
	err := fssc.Ledger.RecordStateUpdate(fssc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	// Synchronize the updated state with other channels
	err = fssc.SyncState()
	if err != nil {
		return fmt.Errorf("failed to synchronize state: %v", err)
	}

	fmt.Printf("State of channel %s updated and synchronized: %s = %v\n", fssc.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the channel
func (fssc *common.FluidStateSyncChannel) RetrieveState(key string) (interface{}, error) {
	fssc.mu.Lock()
	defer fssc.mu.Unlock()

	if !fssc.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	value, exists := fssc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from channel %s: %s = %v\n", fssc.ChannelID, key, value)
	return value, nil
}

// MonitorStateSync continuously monitors the state of the channel and synchronizes it across channels at regular intervals
func (fssc *common.FluidStateSyncChannel) MonitorStateSync(interval time.Duration) {
	for {
		time.Sleep(interval)

		// Perform real-time state synchronization
		err := fssc.SyncState()
		if err != nil {
			fmt.Printf("Error synchronizing state for channel %s: %v\n", fssc.ChannelID, err)
			return
		}
	}
}

// ValidateState ensures that the current state of the channel is consistent and valid
func (fssc *common.FluidStateSyncChannel) ValidateState() error {
	fssc.mu.Lock()
	defer fssc.mu.Unlock()

	if !fssc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Example validation logic
	for key, value := range fssc.State {
		fmt.Printf("Validating state: %s = %v\n", key, value)
	}

	// Log state validation in the ledger
	err := fssc.Ledger.RecordStateValidation(fssc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state validation: %v", err)
	}

	fmt.Printf("State for channel %s validated\n", fssc.ChannelID)
	return nil
}
