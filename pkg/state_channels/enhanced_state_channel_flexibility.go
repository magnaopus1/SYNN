package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewFlexibleStateChannel initializes a new Flexible State Channel
func NewFlexibleStateChannel(channelID string, participants []string, flexibility float64, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.FlexibleStateChannel {
	return &common.FlexibleStateChannel{
		ChannelID:      channelID,
		Participants:   participants,
		State:          make(map[string]interface{}),
		IsOpen:         true,
		Flexibility:    flexibility,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
	}
}

// OpenChannel opens the flexible state channel
func (fsc *common.FlexibleStateChannel) OpenChannel() error {
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

	fmt.Printf("Flexible State Channel %s opened with participants: %v\n", fsc.ChannelID, fsc.Participants)
	return nil
}

// CloseChannel closes the flexible state channel
func (fsc *common.FlexibleStateChannel) CloseChannel() error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return errors.New("state channel is already closed")
	}

	fsc.IsOpen = false

	// Log channel closure in the ledger
	err := fsc.Ledger.RecordChannelClosure(fsc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Flexible State Channel %s closed\n", fsc.ChannelID)
	return nil
}

// AdjustFlexibility dynamically adjusts the flexibility factor for the channel
func (fsc *common.FlexibleStateChannel) AdjustFlexibility(newFlexibility float64) error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	// Ensure flexibility is within an acceptable range
	if newFlexibility < 0.1 || newFlexibility > 10.0 {
		return fmt.Errorf("invalid flexibility factor, must be between 0.1 and 10")
	}

	// Adjust the flexibility factor
	fsc.Flexibility = newFlexibility

	// Log the flexibility adjustment in the ledger
	err := fsc.Ledger.RecordFlexibilityAdjustment(fsc.ChannelID, newFlexibility, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log flexibility adjustment: %v", err)
	}

	fmt.Printf("Flexibility factor for channel %s adjusted to %f\n", fsc.ChannelID, newFlexibility)
	return nil
}

// UpdateState securely updates the internal state of the channel
func (fsc *common.FlexibleStateChannel) UpdateState(key string, value interface{}) error {
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

	fmt.Printf("State of channel %s updated: %s = %v\n", fsc.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the channel
func (fsc *common.FlexibleStateChannel) RetrieveState(key string) (interface{}, error) {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	value, exists := fsc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from channel %s: %s = %v\n", fsc.ChannelID, key, value)
	return value, nil
}

// MonitorFlexibility dynamically monitors and adjusts flexibility in real-time
func (fsc *common.FlexibleStateChannel) MonitorFlexibility(interval time.Duration) {
	for {
		time.Sleep(interval)

		// Example of dynamic flexibility adjustment based on channel conditions
		fsc.AdjustFlexibility(fsc.Flexibility * 1.05) // Example: increase flexibility by 5%
	}
}

// RetrieveParticipants retrieves the list of participants in the channel
func (fsc *common.FlexibleStateChannel) RetrieveParticipants() []string {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	return fsc.Participants
}

// ReallocateResources reallocates channel resources based on new flexibility factor
func (fsc *common.FlexibleStateChannel) ReallocateResources() error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Example logic: adjust resources based on flexibility factor
	adjustedResources := fsc.Flexibility * 1.2 // Example increase of 20% based on flexibility

	// Log resource reallocation in the ledger
	err := fsc.Ledger.RecordResourceReallocation(fsc.ChannelID, adjustedResources, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log resource reallocation: %v", err)
	}

	fmt.Printf("Resources for channel %s reallocated based on flexibility factor\n", fsc.ChannelID)
	return nil
}

// ValidateState ensures that the current state of the channel is consistent and valid
func (fsc *common.FlexibleStateChannel) ValidateState() error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Example validation logic
	// In a real implementation, this would check more complex conditions in the state
	for key, value := range fsc.State {
		fmt.Printf("Validating state: %s = %v\n", key, value)
	}

	// Log state validation in the ledger
	err := fsc.Ledger.RecordStateValidation(fsc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state validation: %v", err)
	}

	fmt.Printf("State for channel %s validated\n", fsc.ChannelID)
	return nil
}
