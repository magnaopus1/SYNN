package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewDynamicStateFragmentationChannel initializes a new dynamic state fragmentation channel
func NewDynamicStateFragmentationChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.DynamicStateFragmentationChannel {
	return &common.DynamicStateFragmentationChannel{
		ChannelID:     channelID,
		Participants:  participants,
		Fragments:     make(map[int]*common.StateFragment),
		State:         make(map[string]interface{}),
		FragmentCount: 0,  // Initial fragment count
		Ledger:        ledgerInstance,
		Encryption:    encryptionService,
		NetworkManager: networkManager,
	}
}

// FragmentState breaks the state into smaller fragments and distributes them
func (dsfc *common.DynamicStateFragmentationChannel) FragmentState(stateData string, numFragments int) error {
	dsfc.mu.Lock()
	defer dsfc.mu.Unlock()

	// Divide the state data into smaller pieces (fragments)
	stateFragments, err := dsfc.createStateFragments(stateData, numFragments)
	if err != nil {
		return fmt.Errorf("failed to fragment state: %v", err)
	}

	// Encrypt and store each fragment
	for i, fragment := range stateFragments {
		encryptedFragment, err := dsfc.Encryption.EncryptData([]byte(fragment), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt fragment: %v", err)
		}

		// Create and store the state fragment
		dsfc.Fragments[i] = &common.StateFragment{
			FragmentID:    fmt.Sprintf("%s-%d", dsfc.ChannelID, i),
			ChannelID:     dsfc.ChannelID,
			FragmentIndex: i,
			Data:          string(encryptedFragment),
			Timestamp:     time.Now(),
		}
	}

	dsfc.FragmentCount = numFragments

	// Log the fragmentation event in the ledger
	err = dsfc.Ledger.RecordStateFragmentation(dsfc.ChannelID, numFragments, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state fragmentation: %v", err)
	}

	fmt.Printf("State fragmented into %d pieces for channel %s\n", numFragments, dsfc.ChannelID)
	return nil
}

// createStateFragments splits the state data into smaller pieces
func (dsfc *common.DynamicStateFragmentationChannel) createStateFragments(stateData string, numFragments int) ([]string, error) {
	stateLength := len(stateData)
	fragmentSize := stateLength / numFragments

	// Check for valid fragmentation
	if fragmentSize == 0 {
		return nil, errors.New("invalid fragment size")
	}

	var fragments []string
	for i := 0; i < numFragments; i++ {
		start := i * fragmentSize
		end := start + fragmentSize
		if i == numFragments-1 {
			end = stateLength
		}
		fragments = append(fragments, stateData[start:end])
	}

	return fragments, nil
}

// ReassembleState takes the fragments and reassembles the original state
func (dsfc *common.DynamicStateFragmentationChannel) ReassembleState() (string, error) {
	dsfc.mu.Lock()
	defer dsfc.mu.Unlock()

	if len(dsfc.Fragments) != dsfc.FragmentCount {
		return "", errors.New("missing state fragments")
	}

	var reassembledState string
	for i := 0; i < dsfc.FragmentCount; i++ {
		fragment, exists := dsfc.Fragments[i]
		if !exists {
			return "", fmt.Errorf("fragment %d missing", i)
		}

		// Decrypt the fragment
		decryptedData, err := dsfc.Encryption.DecryptData([]byte(fragment.Data), common.EncryptionKey)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt fragment %d: %v", i, err)
		}

		reassembledState += string(decryptedData)
	}

	// Log the reassembly event in the ledger
	err := dsfc.Ledger.RecordStateReassembly(dsfc.ChannelID, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to log state reassembly: %v", err)
	}

	fmt.Printf("State reassembled for channel %s\n", dsfc.ChannelID)
	return reassembledState, nil
}

// DistributeFragment distributes a specific state fragment to a participant
func (dsfc *common.DynamicStateFragmentationChannel) DistributeFragment(fragmentIndex int, participant string) error {
	dsfc.mu.Lock()
	defer dsfc.mu.Unlock()

	fragment, exists := dsfc.Fragments[fragmentIndex]
	if !exists {
		return fmt.Errorf("fragment %d not found", fragmentIndex)
	}

	// Send the fragment to the participant
	err := dsfc.NetworkManager.SendDataToParticipant(participant, fragment.Data)
	if err != nil {
		return fmt.Errorf("failed to send fragment %d to participant %s: %v", fragmentIndex, participant, err)
	}

	// Log the fragment distribution in the ledger
	err = dsfc.Ledger.RecordFragmentDistribution(dsfc.ChannelID, fragmentIndex, participant, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log fragment distribution: %v", err)
	}

	fmt.Printf("Fragment %d sent to participant %s\n", fragmentIndex, participant)
	return nil
}

// RetrieveFragment retrieves a specific fragment by index
func (dsfc *common.DynamicStateFragmentationChannel) RetrieveFragment(fragmentIndex int) (*common.StateFragment, error) {
	dsfc.mu.Lock()
	defer dsfc.mu.Unlock()

	fragment, exists := dsfc.Fragments[fragmentIndex]
	if !exists {
		return nil, fmt.Errorf("fragment %d not found", fragmentIndex)
	}

	fmt.Printf("Retrieved fragment %d for channel %s\n", fragmentIndex, dsfc.ChannelID)
	return fragment, nil
}

// UpdateState securely updates the internal state of the channel
func (dsfc *common.DynamicStateFragmentationChannel) UpdateState(key string, value interface{}) error {
	dsfc.mu.Lock()
	defer dsfc.mu.Unlock()

	// Encrypt the state data before updating
	encryptedValue, err := dsfc.Encryption.EncryptData([]byte(fmt.Sprintf("%v", value)), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt state value: %v", err)
	}

	// Update the channel state
	dsfc.State[key] = string(encryptedValue)

	// Log the state update in the ledger
	err = dsfc.Ledger.RecordStateUpdate(dsfc.ChannelID, key, encryptedValue, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of channel %s updated: %s = %v (encrypted)\n", dsfc.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the channel
func (dsfc *common.DynamicStateFragmentationChannel) RetrieveState(key string) (interface{}, error) {
	dsfc.mu.Lock()
	defer dsfc.mu.Unlock()

	value, exists := dsfc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	// Decrypt the state value
	decryptedValue, err := dsfc.Encryption.DecryptData([]byte(value.(string)), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt state value: %v", err)
	}

	fmt.Printf("Retrieved state from channel %s: %s = %v\n", dsfc.ChannelID, key, string(decryptedValue))
	return string(decryptedValue), nil
}
