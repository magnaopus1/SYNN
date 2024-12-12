package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewContinuousTimeScalingChannel initializes a new Continuous Time Scaling Channel (CTSC)
func NewContinuousTimeScalingChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.ContinuousTimeScalingChannel {
	return &common.ContinuousTimeScalingChannel{
		ChannelID:      channelID,
		Participants:   participants,
		State:          make(map[string]interface{}),
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
		ScalingFactor:  1.0, // Default scaling factor
	}
}

// OpenChannel opens the continuous-time scaling state channel
func (ctsc *common.ContinuousTimeScalingChannel) OpenChannel() error {
	ctsc.mu.Lock()
	defer ctsc.mu.Unlock()

	if ctsc.IsOpen {
		return errors.New("state channel is already open")
	}

	ctsc.IsOpen = true

	// Log channel opening in the ledger
	err := ctsc.Ledger.RecordChannelOpening(ctsc.ChannelID, ctsc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Continuous-Time Scaling Channel %s opened with participants: %v\n", ctsc.ChannelID, ctsc.Participants)
	return nil
}

// CloseChannel closes the continuous-time scaling state channel
func (ctsc *common.ContinuousTimeScalingChannel) CloseChannel() error {
	ctsc.mu.Lock()
	defer ctsc.mu.Unlock()

	if !ctsc.IsOpen {
		return errors.New("state channel is already closed")
	}

	// Final optimization before closure
	err := ctsc.OptimizeScaling()
	if err != nil {
		return fmt.Errorf("failed to optimize scaling before closing: %v", err)
	}

	ctsc.IsOpen = false

	// Log channel closure in the ledger
	err = ctsc.Ledger.RecordChannelClosure(ctsc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Continuous-Time Scaling Channel %s closed\n", ctsc.ChannelID)
	return nil
}

// OptimizeScaling adjusts the scaling factor of the channel in real-time.
func (ctsc *common.ContinuousTimeScalingChannel) OptimizeScaling() error {
	ctsc.mu.Lock()
	defer ctsc.mu.Unlock()

	if !ctsc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Here you can use real-time metrics, load balancing, or other network performance indicators
	// to adjust the scaling factor. This example increases it by 10%:
	ctsc.ScalingFactor *= 1.1

	// Log the optimization event in the ledger
	err := ctsc.Ledger.RecordScalingEvent(ctsc.ChannelID, ctsc.ScalingFactor, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log scaling optimization: %v", err)
	}

	fmt.Printf("Scaling factor for channel %s optimized to %f\n", ctsc.ChannelID, ctsc.ScalingFactor)
	return nil
}

// MonitorScaling continuously monitors and adjusts the scaling in real-time at specified intervals.
func (ctsc *common.ContinuousTimeScalingChannel) MonitorScaling(interval time.Duration) {
	for {
		time.Sleep(interval)
		err := ctsc.OptimizeScaling()
		if err != nil {
			fmt.Printf("Error optimizing scaling for channel %s: %v\n", ctsc.ChannelID, err)
			return
		}
	}
}

// UpdateState securely updates the internal state of the channel.
func (ctsc *common.ContinuousTimeScalingChannel) UpdateState(key string, value interface{}) error {
	ctsc.mu.Lock()
	defer ctsc.mu.Unlock()

	if !ctsc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Encrypt state value before updating
	encryptedValue, err := ctsc.Encryption.EncryptData([]byte(fmt.Sprintf("%v", value)), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt state value: %v", err)
	}

	ctsc.State[key] = string(encryptedValue)

	// Log the state update in the ledger
	err = ctsc.Ledger.RecordStateUpdate(ctsc.ChannelID, key, encryptedValue, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of channel %s updated: %s = %v (encrypted)\n", ctsc.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the channel.
func (ctsc *common.ContinuousTimeScalingChannel) RetrieveState(key string) (interface{}, error) {
	ctsc.mu.Lock()
	defer ctsc.mu.Unlock()

	if !ctsc.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	value, exists := ctsc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	// Decrypt state value before returning
	decryptedValue, err := ctsc.Encryption.DecryptData([]byte(value.(string)), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt state value: %v", err)
	}

	fmt.Printf("Retrieved state from channel %s: %s = %v\n", ctsc.ChannelID, key, string(decryptedValue))
	return string(decryptedValue), nil
}

// ReallocateResources dynamically reallocates resources to optimize performance in real-time.
func (ctsc *common.ContinuousTimeScalingChannel) ReallocateResources() error {
	ctsc.mu.Lock()
	defer ctsc.mu.Unlock()

	if !ctsc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Dynamically adjust scaling factor (for example, reduce by 5% under heavy load)
	adjustedScaling := ctsc.ScalingFactor * 0.95

	ctsc.ScalingFactor = adjustedScaling

	// Log the reallocation event in the ledger
	err := ctsc.Ledger.RecordResourceReallocation(ctsc.ChannelID, ctsc.ScalingFactor, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log resource reallocation: %v", err)
	}

	fmt.Printf("Resources for channel %s reallocated, new scaling factor: %f\n", ctsc.ChannelID, ctsc.ScalingFactor)
	return nil
}
