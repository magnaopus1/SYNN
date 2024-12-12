package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewStateChannelInteroperability initializes a new interoperability handler for state channels
func NewStateChannelInteroperability(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.Manager) *common.StateChannelInteroperability {
	return &common.StateChannelInteroperability{
		ChannelID:       channelID,
		LinkedChannels:  make(map[string]*common.InteropLink),
		IsInteroperabilityEnabled: true,
		Participants:    participants,
		State:           make(map[string]interface{}),
		Ledger:          ledgerInstance,
		Encryption:      encryptionService,
		NetworkManager:  networkManager,
	}
}

// EnableInteroperability enables cross-channel operations for the state channel
func (sci *common.StateChannelInteroperability) EnableInteroperability() error {
	sci.mu.Lock()
	defer sci.mu.Unlock()

	if sci.IsInteroperabilityEnabled {
		return errors.New("interoperability is already enabled")
	}

	sci.IsInteroperabilityEnabled = true
	err := sci.Ledger.RecordInteroperabilityStatus(sci.ChannelID, true, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log interoperability status: %v", err)
	}

	fmt.Printf("Interoperability enabled for state channel %s\n", sci.ChannelID)
	return nil
}

// DisableInteroperability disables cross-channel operations for the state channel
func (sci *common.StateChannelInteroperability) DisableInteroperability() error {
	sci.mu.Lock()
	defer sci.mu.Unlock()

	if !sci.IsInteroperabilityEnabled {
		return errors.New("interoperability is already disabled")
	}

	sci.IsInteroperabilityEnabled = false
	err := sci.Ledger.RecordInteroperabilityStatus(sci.ChannelID, false, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log interoperability status: %v", err)
	}

	fmt.Printf("Interoperability disabled for state channel %s\n", sci.ChannelID)
	return nil
}

// LinkStateChannel creates a link between two state channels for interoperability
func (sci *common.StateChannelInteroperability) LinkStateChannel(targetChannelID string) (*common.InteropLink, error) {
	sci.mu.Lock()
	defer sci.mu.Unlock()

	if !sci.IsInteroperabilityEnabled {
		return nil, errors.New("interoperability is disabled")
	}

	// Check if the link already exists
	for _, link := range sci.LinkedChannels {
		if link.TargetChannelID == targetChannelID {
			return nil, errors.New("a link to the target channel already exists")
		}
	}

	linkID := common.GenerateUniqueID()
	interopLink := &common.InteropLink{
		LinkID:          linkID,
		TargetChannelID: targetChannelID,
		SharedData:      make(map[string]interface{}),
		Timestamp:       time.Now(),
	}

	// Add the new link to the map
	sci.LinkedChannels[linkID] = interopLink

	// Log the channel linking event in the ledger
	err := sci.Ledger.RecordStateChannelLink(sci.ChannelID, targetChannelID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log state channel linking: %v", err)
	}

	fmt.Printf("State channel %s linked with %s\n", sci.ChannelID, targetChannelID)
	return interopLink, nil
}

// ShareDataAcrossChannels shares data between linked state channels
func (sci *common.StateChannelInteroperability) ShareDataAcrossChannels(linkID string, key string, value interface{}) error {
	sci.mu.Lock()
	defer sci.mu.Unlock()

	link, exists := sci.LinkedChannels[linkID]
	if !exists {
		return errors.New("interoperability link not found")
	}

	// Encrypt the data before sharing
	encryptedValue, err := sci.Encryption.EncryptData([]byte(fmt.Sprintf("%v", value)), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %v", err)
	}

	// Share the data across the linked channel
	link.SharedData[key] = string(encryptedValue)

	// Log the data sharing event
	err = sci.Ledger.RecordDataSharing(sci.ChannelID, link.TargetChannelID, key, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log data sharing: %v", err)
	}

	fmt.Printf("Data shared between channels %s and %s: %s = %v\n", sci.ChannelID, link.TargetChannelID, key, value)
	return nil
}

// RetrieveSharedData retrieves shared data from the linked channels
func (sci *common.StateChannelInteroperability) RetrieveSharedData(linkID string, key string) (interface{}, error) {
	sci.mu.Lock()
	defer sci.mu.Unlock()

	link, exists := sci.LinkedChannels[linkID]
	if !exists {
		return nil, errors.New("interoperability link not found")
	}

	value, exists := link.SharedData[key]
	if !exists {
		return nil, fmt.Errorf("shared data for key %s not found", key)
	}

	// Decrypt the data
	decryptedValue, err := sci.Encryption.DecryptData([]byte(value), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}

	fmt.Printf("Retrieved shared data for key %s: %s\n", key, decryptedValue)
	return string(decryptedValue), nil
}

// RemoveLink removes a link between two state channels
func (sci *common.StateChannelInteroperability) RemoveLink(linkID string) error {
	sci.mu.Lock()
	defer sci.mu.Unlock()

	_, exists := sci.LinkedChannels[linkID]
	if !exists {
		return fmt.Errorf("link %s not found", linkID)
	}

	delete(sci.LinkedChannels, linkID)

	// Log the link removal in the ledger
	err := sci.Ledger.RecordLinkRemoval(sci.ChannelID, linkID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log link removal: %v", err)
	}

	fmt.Printf("Link %s removed from channel %s\n", linkID, sci.ChannelID)
	return nil
}

// RetrieveLinkedChannels retrieves all the linked channels
func (sci *common.StateChannelInteroperability) RetrieveLinkedChannels() ([]*common.InteropLink, error) {
	sci.mu.Lock()
	defer sci.mu.Unlock()

	var links []*common.InteropLink
	for _, link := range sci.LinkedChannels {
		links = append(links, link)
	}

	fmt.Printf("Retrieved %d linked channels for channel %s\n", len(links), sci.ChannelID)
	return links, nil
}
