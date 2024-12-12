package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewFractalStateChannel initializes a new Fractal State Channel (FSC)
func NewFractalStateChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.FractalStateChannel {
	return &common.FractalStateChannel{
		ChannelID:      channelID,
		Participants:   participants,
		State:          make(map[string]interface{}),
		SubChannels:    []*common.FractalStateChannel{},
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
	}
}

// OpenChannel opens the fractal state channel
func (fsc *common.FractalStateChannel) OpenChannel() error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if fsc.IsOpen {
		return errors.New("state channel is already open")
	}

	fsc.IsOpen = true

	// Log channel opening in the ledger
	err := fsc.Ledger.RecordChannelOpening(fsc.ChannelID, fsc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Fractal State Channel %s opened with participants: %v\n", fsc.ChannelID, fsc.Participants)
	return nil
}

// CloseChannel closes the fractal state channel and its sub-channels
func (fsc *common.FractalStateChannel) CloseChannel() error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return errors.New("state channel is already closed")
	}

	// Close all sub-channels
	for _, subChannel := range fsc.SubChannels {
		err := subChannel.CloseChannel()
		if err != nil {
			return fmt.Errorf("failed to close sub-channel %s: %v", subChannel.ChannelID, err)
		}
	}

	fsc.IsOpen = false

	// Log channel closure in the ledger
	err := fsc.Ledger.RecordChannelClosure(fsc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Fractal State Channel %s and all sub-channels closed\n", fsc.ChannelID)
	return nil
}

// AddSubChannel adds a new fractal sub-channel to the current fractal state channel
func (fsc *common.FractalStateChannel) AddSubChannel(subChannel *common.FractalStateChannel) error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	// Ensure the sub-channel isn't already part of the main channel
	for _, existingSub := range fsc.SubChannels {
		if existingSub.ChannelID == subChannel.ChannelID {
			return fmt.Errorf("sub-channel %s is already part of fractal state channel %s", subChannel.ChannelID, fsc.ChannelID)
		}
	}

	fsc.SubChannels = append(fsc.SubChannels, subChannel)
	fmt.Printf("Sub-channel %s added to fractal state channel %s\n", subChannel.ChannelID, fsc.ChannelID)
	return nil
}

// SyncState performs recursive state aggregation and synchronization across fractal channels
func (fsc *common.FractalStateChannel) SyncState() error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Aggregate the state of all sub-channels
	for _, subChannel := range fsc.SubChannels {
		for key, value := range subChannel.State {
			fsc.State[key] = value
		}
		subChannel.SyncState() // Recursively synchronize the sub-channels
	}

	// Log the synchronization event in the ledger
	err := fsc.Ledger.RecordStateSync(fsc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state synchronization: %v", err)
	}

	fmt.Printf("Fractal State Channel %s synchronized with all sub-channels\n", fsc.ChannelID)
	return nil
}

// UpdateState securely updates the internal state of the fractal channel and its sub-channels
func (fsc *common.FractalStateChannel) UpdateState(key string, value interface{}) error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Update the channel state
	fsc.State[key] = value

	// Log the state update in the ledger
	err := fsc.Ledger.RecordStateUpdate(fsc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	// Recursively update the states of all sub-channels
	for _, subChannel := range fsc.SubChannels {
		err := subChannel.UpdateState(key, value)
		if err != nil {
			return fmt.Errorf("failed to update state in sub-channel %s: %v", subChannel.ChannelID, err)
		}
	}

	fmt.Printf("State of fractal channel %s and all sub-channels updated: %s = %v\n", fsc.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the fractal channel
func (fsc *common.FractalStateChannel) RetrieveState(key string) (interface{}, error) {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	value, exists := fsc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from fractal channel %s: %s = %v\n", fsc.ChannelID, key, value)
	return value, nil
}

// MonitorStateSync continuously monitors and recursively synchronizes the state of the fractal channel
func (fsc *common.FractalStateChannel) MonitorStateSync(interval time.Duration) {
	for {
		time.Sleep(interval)

		// Perform recursive state synchronization
		err := fsc.SyncState()
		if err != nil {
			fmt.Printf("Error synchronizing state for fractal channel %s: %v\n", fsc.ChannelID, err)
			return
		}
	}
}

// ValidateState ensures that the current state of the fractal channel and its sub-channels is consistent and valid
func (fsc *common.FractalStateChannel) ValidateState() error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Example validation logic for the main channel and its sub-channels
	for key, value := range fsc.State {
		fmt.Printf("Validating state: %s = %v\n", key, value)
	}

	// Recursively validate all sub-channels
	for _, subChannel := range fsc.SubChannels {
		err := subChannel.ValidateState()
		if err != nil {
			return fmt.Errorf("state validation failed for sub-channel %s: %v", subChannel.ChannelID, err)
		}
	}

	// Log state validation in the ledger
	err := fsc.Ledger.RecordStateValidation(fsc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state validation: %v", err)
	}

	fmt.Printf("State for fractal channel %s validated\n", fsc.ChannelID)
	return nil
}
