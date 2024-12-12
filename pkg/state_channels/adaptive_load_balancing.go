package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewAdaptiveLoadBalancingChannel initializes a new state channel with adaptive load-balancing
func NewAdaptiveLoadBalancingChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.AdaptiveLoadBalancingChannel {
	return &common.AdaptiveLoadBalancingChannel{
		ChannelID:      channelID,
		Participants:   participants,
		State:          make(map[string]interface{}),
		LoadMetrics:    make(map[string]float64),
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
	}
}

// OpenChannel opens the adaptive load-balancing state channel
func (alc *common.AdaptiveLoadBalancingChannel) OpenChannel() error {
	alc.mu.Lock()
	defer alc.mu.Unlock()

	if alc.IsOpen {
		return errors.New("state channel is already open")
	}

	alc.IsOpen = true

	// Log channel opening in the ledger
	err := alc.Ledger.RecordChannelOpening(alc.ChannelID, alc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Adaptive Load-Balancing Channel %s opened with participants: %v\n", alc.ChannelID, alc.Participants)
	return nil
}

// CloseChannel closes the adaptive load-balancing state channel
func (alc *common.AdaptiveLoadBalancingChannel) CloseChannel() error {
	alc.mu.Lock()
	defer alc.mu.Unlock()

	if !alc.IsOpen {
		return errors.New("state channel is already closed")
	}

	alc.IsOpen = false

	// Log channel closure in the ledger
	err := alc.Ledger.RecordChannelClosure(alc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Adaptive Load-Balancing Channel %s closed\n", alc.ChannelID)
	return nil
}

// MonitorLoad updates the load metrics for each participant in the channel
func (alc *common.AdaptiveLoadBalancingChannel) MonitorLoad(participant string, loadValue float64) error {
	alc.mu.Lock()
	defer alc.mu.Unlock()

	if !alc.IsOpen {
		return errors.New("state channel is closed")
	}

	alc.LoadMetrics[participant] = loadValue

	// Log the load monitoring in the ledger
	err := alc.Ledger.RecordLoadMetrics(alc.ChannelID, participant, loadValue, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log load metrics: %v", err)
	}

	fmt.Printf("Load for participant %s updated to %f in channel %s\n", participant, loadValue, alc.ChannelID)
	return nil
}

// ReallocateLoad dynamically reallocates load to balance the channel's state
func (alc *common.AdaptiveLoadBalancingChannel) ReallocateLoad() error {
	alc.mu.Lock()
	defer alc.mu.Unlock()

	if !alc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Logic to calculate load balancing based on current metrics
	for participant, load := range alc.LoadMetrics {
		// Adjust load by distributing it across participants/nodes based on some adaptive algorithm
		adjustedLoad := alc.NetworkManager.AdjustLoad(participant, load)
		alc.LoadMetrics[participant] = adjustedLoad

		// Log the load reallocation
		err := alc.Ledger.RecordLoadReallocation(alc.ChannelID, participant, adjustedLoad, time.Now())
		if err != nil {
			return fmt.Errorf("failed to log load reallocation: %v", err)
		}
		fmt.Printf("Load for participant %s reallocated to %f in channel %s\n", participant, adjustedLoad, alc.ChannelID)
	}

	return nil
}

// RetrieveLoad retrieves the current load metrics for a participant
func (alc *common.AdaptiveLoadBalancingChannel) RetrieveLoad(participant string) (float64, error) {
	alc.mu.Lock()
	defer alc.mu.Unlock()

	load, exists := alc.LoadMetrics[participant]
	if !exists {
		return 0, fmt.Errorf("load for participant %s not found in channel %s", participant, alc.ChannelID)
	}

	fmt.Printf("Retrieved load for participant %s: %f in channel %s\n", participant, load, alc.ChannelID)
	return load, nil
}

// RetrieveState retrieves the current state of the channel
func (alc *common.AdaptiveLoadBalancingChannel) RetrieveState(key string) (interface{}, error) {
	alc.mu.Lock()
	defer alc.mu.Unlock()

	if !alc.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	value, exists := alc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from channel %s: %s = %v\n", alc.ChannelID, key, value)
	return value, nil
}

// UpdateState securely updates the internal state of the channel
func (alc *common.AdaptiveLoadBalancingChannel) UpdateState(key string, value interface{}) error {
	alc.mu.Lock()
	defer alc.mu.Unlock()

	if !alc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Update the channel state
	alc.State[key] = value

	// Log the state update in the ledger
	err := alc.Ledger.RecordStateUpdate(alc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of channel %s updated: %s = %v\n", alc.ChannelID, key, value)
	return nil
}
