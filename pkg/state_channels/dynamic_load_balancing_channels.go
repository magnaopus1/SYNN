package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewDynamicLoadBalancingChannel initializes a new dynamic load-balancing state channel
func NewDynamicLoadBalancingChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.DynamicLoadBalancingChannel {
	return &common.DynamicLoadBalancingChannel{
		ChannelID:      channelID,
		Participants:   participants,
		State:          make(map[string]interface{}),
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
		LoadThreshold:  0.75, // Default load threshold
		CurrentLoad:    0.0,  // Initially no load
	}
}

// OpenChannel opens the dynamic load-balancing state channel
func (dlbc *common.DynamicLoadBalancingChannel) OpenChannel() error {
	dlbc.mu.Lock()
	defer dlbc.mu.Unlock()

	if dlbc.IsOpen {
		return errors.New("state channel is already open")
	}

	dlbc.IsOpen = true

	// Log channel opening in the ledger
	err := dlbc.Ledger.RecordChannelOpening(dlbc.ChannelID, dlbc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Dynamic Load-Balancing Channel %s opened with participants: %v\n", dlbc.ChannelID, dlbc.Participants)
	return nil
}

// CloseChannel closes the dynamic load-balancing state channel
func (dlbc *common.DynamicLoadBalancingChannel) CloseChannel() error {
	dlbc.mu.Lock()
	defer dlbc.mu.Unlock()

	if !dlbc.IsOpen {
		return errors.New("state channel is already closed")
	}

	dlbc.IsOpen = false

	// Log channel closure in the ledger
	err := dlbc.Ledger.RecordChannelClosure(dlbc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Dynamic Load-Balancing Channel %s closed\n", dlbc.ChannelID)
	return nil
}

// UpdateState updates the state of the channel and performs load-balancing if necessary
func (dlbc *common.DynamicLoadBalancingChannel) UpdateState(key string, value interface{}) error {
	dlbc.mu.Lock()
	defer dlbc.mu.Unlock()

	if !dlbc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Encrypt the state data before updating
	encryptedValue, err := dlbc.Encryption.EncryptData([]byte(fmt.Sprintf("%v", value)), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt state value: %v", err)
	}

	// Update the state
	dlbc.State[key] = string(encryptedValue)

	// Log the state update in the ledger
	err = dlbc.Ledger.RecordStateUpdate(dlbc.ChannelID, key, encryptedValue, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of channel %s updated: %s = %v (encrypted)\n", dlbc.ChannelID, key, value)

	// Trigger load balancing if the current load exceeds the threshold
	dlbc.CurrentLoad = dlbc.CalculateCurrentLoad()
	if dlbc.CurrentLoad > dlbc.LoadThreshold {
		err = dlbc.PerformLoadBalancing()
		if err != nil {
			return fmt.Errorf("failed to perform load balancing: %v", err)
		}
	}

	return nil
}

// RetrieveState retrieves the current state of the channel
func (dlbc *common.DynamicLoadBalancingChannel) RetrieveState(key string) (interface{}, error) {
	dlbc.mu.Lock()
	defer dlbc.mu.Unlock()

	if !dlbc.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	value, exists := dlbc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	// Decrypt the state value
	decryptedValue, err := dlbc.Encryption.DecryptData([]byte(value.(string)), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt state value: %v", err)
	}

	fmt.Printf("Retrieved state from channel %s: %s = %v\n", dlbc.ChannelID, key, string(decryptedValue))
	return string(decryptedValue), nil
}

// CalculateCurrentLoad calculates the current load on the channel
func (dlbc *common.DynamicLoadBalancingChannel) CalculateCurrentLoad() float64 {
	// In a real-world scenario, the load would be calculated based on various metrics
	// For simplicity, we'll assume it increases by a small fraction with every state update
	// In practice, you'd use actual metrics like transaction volume, network latency, etc.
	loadIncrease := 0.1 // Simulated load increase
	newLoad := dlbc.CurrentLoad + loadIncrease
	if newLoad > 1.0 {
		newLoad = 1.0 // Max load is 100%
	}
	fmt.Printf("Current load for channel %s calculated as %f\n", dlbc.ChannelID, newLoad)
	return newLoad
}

// PerformLoadBalancing dynamically reallocates resources to balance the load across participants
func (dlbc *common.DynamicLoadBalancingChannel) PerformLoadBalancing() error {
	// Example logic: Redistribute channel resources based on the calculated load
	fmt.Printf("Performing load balancing for channel %s. Current load: %f\n", dlbc.ChannelID, dlbc.CurrentLoad)

	// Simulated reallocation of resources:
	newLoad := dlbc.CurrentLoad * 0.8 // Example: reduce load by 20%

	// Update current load after balancing
	dlbc.CurrentLoad = newLoad

	// Log the load balancing event in the ledger
	err := dlbc.Ledger.RecordLoadBalancingEvent(dlbc.ChannelID, dlbc.CurrentLoad, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log load balancing event: %v", err)
	}

	fmt.Printf("Load balancing complete for channel %s. New load: %f\n", dlbc.ChannelID, dlbc.CurrentLoad)
	return nil
}

// MonitorLoad continuously monitors the channel load and triggers balancing as needed
func (dlbc *common.DynamicLoadBalancingChannel) MonitorLoad(interval time.Duration) {
	for {
		time.Sleep(interval)
		dlbc.mu.Lock()
		currentLoad := dlbc.CalculateCurrentLoad()
		if currentLoad > dlbc.LoadThreshold {
			err := dlbc.PerformLoadBalancing()
			if err != nil {
				fmt.Printf("Error performing load balancing for channel %s: %v\n", dlbc.ChannelID, err)
			}
		}
		dlbc.mu.Unlock()
	}
}
