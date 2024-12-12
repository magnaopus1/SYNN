package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)


// NewDynamicResourceAllocationChannel initializes a new dynamic resource allocation state channel
func NewDynamicResourceAllocationChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.DynamicResourceAllocationChannel {
	return &common.DynamicResourceAllocationChannel{
		ChannelID:      channelID,
		Participants:   participants,
		State:          make(map[string]interface{}),
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
		LoadThreshold:  0.75, // Default load threshold
		ResourceUsage:  0.0,  // Initial resource usage
	}
}

// OpenChannel opens the dynamic resource allocation state channel
func (drac *common.DynamicResourceAllocationChannel) OpenChannel() error {
	drac.mu.Lock()
	defer drac.mu.Unlock()

	if drac.IsOpen {
		return errors.New("state channel is already open")
	}

	drac.IsOpen = true

	// Log channel opening in the ledger
	err := drac.Ledger.RecordChannelOpening(drac.ChannelID, drac.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Dynamic Resource Allocation Channel %s opened with participants: %v\n", drac.ChannelID, drac.Participants)
	return nil
}

// CloseChannel closes the dynamic resource allocation state channel
func (drac *common.DynamicResourceAllocationChannel) CloseChannel() error {
	drac.mu.Lock()
	defer drac.mu.Unlock()

	if !drac.IsOpen {
		return errors.New("state channel is already closed")
	}

	// Final resource reallocation before closure
	err := drac.ReallocateResources()
	if err != nil {
		return fmt.Errorf("failed to reallocate resources before closing: %v", err)
	}

	drac.IsOpen = false

	// Log channel closure in the ledger
	err = drac.Ledger.RecordChannelClosure(drac.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Dynamic Resource Allocation Channel %s closed\n", drac.ChannelID)
	return nil
}

// UpdateState securely updates the internal state of the channel
func (drac *common.DynamicResourceAllocationChannel) UpdateState(key string, value interface{}) error {
	drac.mu.Lock()
	defer drac.mu.Unlock()

	if !drac.IsOpen {
		return errors.New("state channel is closed")
	}

	// Encrypt the state data before updating
	encryptedValue, err := drac.Encryption.EncryptData([]byte(fmt.Sprintf("%v", value)), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt state value: %v", err)
	}

	// Update the state
	drac.State[key] = string(encryptedValue)

	// Log the state update in the ledger
	err = drac.Ledger.RecordStateUpdate(drac.ChannelID, key, encryptedValue, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of channel %s updated: %s = %v (encrypted)\n", drac.ChannelID, key, value)

	// Check for resource reallocation based on the updated load
	drac.ResourceUsage = drac.CalculateResourceUsage()
	if drac.ResourceUsage > drac.LoadThreshold {
		err = drac.ReallocateResources()
		if err != nil {
			return fmt.Errorf("failed to perform resource reallocation: %v", err)
		}
	}

	return nil
}

// RetrieveState retrieves the current state of the channel
func (drac *common.DynamicResourceAllocationChannel) RetrieveState(key string) (interface{}, error) {
	drac.mu.Lock()
	defer drac.mu.Unlock()

	if !drac.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	value, exists := drac.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	// Decrypt the state value
	decryptedValue, err := drac.Encryption.DecryptData([]byte(value.(string)), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt state value: %v", err)
	}

	fmt.Printf("Retrieved state from channel %s: %s = %v\n", drac.ChannelID, key, string(decryptedValue))
	return string(decryptedValue), nil
}

// CalculateResourceUsage calculates the current resource usage of the channel
func (drac *common.DynamicResourceAllocationChannel) CalculateResourceUsage() float64 {
	// Simulated calculation based on channel state
	// In a real-world scenario, this would be based on metrics like network traffic, transaction volume, etc.
	usageIncrease := 0.1 // Simulated increase
	newUsage := drac.ResourceUsage + usageIncrease
	if newUsage > 1.0 {
		newUsage = 1.0 // Max resource usage is 100%
	}
	fmt.Printf("Resource usage for channel %s calculated as %f\n", drac.ChannelID, newUsage)
	return newUsage
}

// ReallocateResources dynamically reallocates resources to optimize channel performance
func (drac *common.DynamicResourceAllocationChannel) ReallocateResources() error {
	drac.mu.Lock()
	defer drac.mu.Unlock()

	if !drac.IsOpen {
		return errors.New("state channel is closed")
	}

	// Example of resource reallocation logic
	fmt.Printf("Performing resource reallocation for channel %s. Current resource usage: %f\n", drac.ChannelID, drac.ResourceUsage)

	// Simulated reallocation of resources:
	newUsage := drac.ResourceUsage * 0.85 // Example: reduce usage by 15%

	// Update resource usage after reallocation
	drac.ResourceUsage = newUsage

	// Log the resource reallocation event in the ledger
	err := drac.Ledger.RecordResourceReallocation(drac.ChannelID, drac.ResourceUsage, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log resource reallocation: %v", err)
	}

	fmt.Printf("Resource reallocation complete for channel %s. New resource usage: %f\n", drac.ChannelID, drac.ResourceUsage)
	return nil
}

// MonitorResourceUsage continuously monitors resource usage and triggers reallocation as needed
func (drac *common.DynamicResourceAllocationChannel) MonitorResourceUsage(interval time.Duration) {
	for {
		time.Sleep(interval)
		drac.mu.Lock()
		resourceUsage := drac.CalculateResourceUsage()
		if resourceUsage > drac.LoadThreshold {
			err := drac.ReallocateResources()
			if err != nil {
				fmt.Printf("Error reallocating resources for channel %s: %v\n", drac.ChannelID, err)
			}
		}
		drac.mu.Unlock()
	}
}
