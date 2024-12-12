package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewPrivacyStateChannel initializes a new privacy-focused state channel
func NewPrivacyStateChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.PrivacyStateChannel {
	return &common.PrivacyStateChannel{
		ChannelID:      channelID,
		Participants:   participants,
		EncryptedState: make(map[string]string),
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		DecryptionKeys: make(map[string]string),
	}
}

// OpenChannel opens the privacy-enabled state channel
func (pc *common.PrivacyStateChannel) OpenChannel() error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if !pc.IsOpen {
		return errors.New("state channel is already open")
	}

	pc.IsOpen = true

	// Log channel opening in the ledger
	err := pc.Ledger.RecordChannelOpening(pc.ChannelID, pc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Privacy state channel %s opened with participants: %v\n", pc.ChannelID, pc.Participants)
	return nil
}

// CloseChannel closes the privacy-enabled state channel
func (pc *common.PrivacyStateChannel) CloseChannel() error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if !pc.IsOpen {
		return errors.New("state channel is already closed")
	}

	pc.IsOpen = false

	// Log channel closure in the ledger
	err := pc.Ledger.RecordChannelClosure(pc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Privacy state channel %s closed\n", pc.ChannelID)
	return nil
}

// EncryptState encrypts and stores sensitive state data
func (pc *common.PrivacyStateChannel) EncryptState(key string, value interface{}) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if !pc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Convert value to bytes and encrypt
	data := fmt.Sprintf("%v", value)
	encryptedData, err := pc.Encryption.EncryptData([]byte(data), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt state data: %v", err)
	}

	// Store the encrypted data
	pc.EncryptedState[key] = string(encryptedData)

	// Log the encryption event in the ledger
	err = pc.Ledger.RecordStateEncryption(pc.ChannelID, key, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log encryption event: %v", err)
	}

	fmt.Printf("State data for key %s encrypted in privacy channel %s\n", key, pc.ChannelID)
	return nil
}

// DecryptState decrypts and retrieves sensitive state data
func (pc *common.PrivacyStateChannel) DecryptState(key string) (interface{}, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	encryptedData, exists := pc.EncryptedState[key]
	if !exists {
		return nil, fmt.Errorf("no encrypted data found for key %s", key)
	}

	// Decrypt the data
	decryptedData, err := pc.Encryption.DecryptData([]byte(encryptedData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt state data: %v", err)
	}

	fmt.Printf("State data for key %s decrypted in privacy channel %s\n", key, pc.ChannelID)
	return string(decryptedData), nil
}

// ShareDecryptionKey allows a participant to share a decryption key with another participant
func (pc *common.PrivacyStateChannel) ShareDecryptionKey(participant string, decryptionKey string) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// Store the decryption key for the participant
	pc.DecryptionKeys[participant] = decryptionKey

	// Log the key-sharing event in the ledger
	err := pc.Ledger.RecordKeySharing(pc.ChannelID, participant, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log key sharing: %v", err)
	}

	fmt.Printf("Decryption key shared with participant %s in privacy channel %s\n", participant, pc.ChannelID)
	return nil
}

// ValidateStateEncryption ensures the state data is properly encrypted
func (pc *common.PrivacyStateChannel) ValidateStateEncryption() error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// Simulate encryption validation logic
	for key, encryptedData := range pc.EncryptedState {
		_, err := pc.Encryption.DecryptData([]byte(encryptedData), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("validation failed for key %s: %v", key, err)
		}
	}

	// Log the encryption validation in the ledger
	err := pc.Ledger.RecordStateEncryptionValidation(pc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log encryption validation: %v", err)
	}

	fmt.Printf("State encryption validated for channel %s\n", pc.ChannelID)
	return nil
}

// RetrieveEncryptedState retrieves the encrypted data for auditing purposes
func (pc *common.PrivacyStateChannel) RetrieveEncryptedState(key string) (string, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	encryptedData, exists := pc.EncryptedState[key]
	if !exists {
		return "", fmt.Errorf("no encrypted data found for key %s", key)
	}

	fmt.Printf("Retrieved encrypted state data for key %s in privacy channel %s\n", key, pc.ChannelID)
	return encryptedData, nil
}
