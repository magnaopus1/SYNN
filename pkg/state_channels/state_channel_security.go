package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewSecurityStateChannel initializes a new secure state channel
func NewSecurityStateChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, securityModule *common.SecurityManager) *common.SecurityStateChannel {
	return &common.SecurityStateChannel{
		ChannelID:       channelID,
		Participants:    participants,
		State:           make(map[string]interface{}),
		IsOpen:          true,
		Ledger:          ledgerInstance,
		Encryption:      encryptionService,
		SecurityModule:  securityModule,
		ParticipantKeys: make(map[string]string),
	}
}

// OpenChannel opens the secure state channel
func (sc *common.SecurityStateChannel) OpenChannel() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("state channel is already open")
	}

	sc.IsOpen = true

	// Log the channel opening in the ledger
	err := sc.Ledger.RecordChannelOpening(sc.ChannelID, sc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Security state channel %s opened with participants: %v\n", sc.ChannelID, sc.Participants)
	return nil
}

// CloseChannel closes the secure state channel and performs security checks
func (sc *common.SecurityStateChannel) CloseChannel() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("state channel is already closed")
	}

	// Perform security checks before closing
	err := sc.SecurityModule.ValidateChannelSecurity(sc.ChannelID, sc.State)
	if err != nil {
		return fmt.Errorf("security validation failed for channel %s: %v", sc.ChannelID, err)
	}

	sc.IsOpen = false

	// Log the channel closure in the ledger
	err = sc.Ledger.RecordChannelClosure(sc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Security state channel %s closed\n", sc.ChannelID)
	return nil
}

// UpdateState securely updates the internal state of the channel
func (sc *common.SecurityStateChannel) UpdateState(key string, value interface{}, signature string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Verify the participant's signature
	valid, err := sc.SecurityModule.VerifySignature(key, signature, sc.ParticipantKeys)
	if err != nil || !valid {
		return errors.New("invalid signature, update rejected")
	}

	// Update the channel state
	sc.State[key] = value

	// Log the state update in the ledger
	err = sc.Ledger.RecordStateUpdate(sc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of channel %s updated securely: %s = %v\n", sc.ChannelID, key, value)
	return nil
}

// EncryptState encrypts and stores the sensitive state data securely
func (sc *common.SecurityStateChannel) EncryptState(key string, value interface{}) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Convert value to bytes and encrypt
	data := fmt.Sprintf("%v", value)
	encryptedData, err := sc.Encryption.EncryptData([]byte(data), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt state data: %v", err)
	}

	// Store the encrypted data in the state
	sc.State[key] = encryptedData

	// Log the encryption event in the ledger
	err = sc.Ledger.RecordStateEncryption(sc.ChannelID, key, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log encryption event: %v", err)
	}

	fmt.Printf("State data for key %s encrypted in security state channel %s\n", key, sc.ChannelID)
	return nil
}

// ValidateState ensures the current state data is valid and secure
func (sc *common.SecurityStateChannel) ValidateState() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Perform security validation of the state
	err := sc.SecurityModule.ValidateState(sc.ChannelID, sc.State)
	if err != nil {
		return fmt.Errorf("state validation failed for channel %s: %v", sc.ChannelID, err)
	}

	// Log the state validation in the ledger
	err = sc.Ledger.RecordStateValidation(sc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state validation: %v", err)
	}

	fmt.Printf("State validated successfully for channel %s\n", sc.ChannelID)
	return nil
}

// AddParticipantKey adds a participant's public key for signature verification
func (sc *common.SecurityStateChannel) AddParticipantKey(participant string, publicKey string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.ParticipantKeys[participant] = publicKey

	// Log the key addition in the ledger
	err := sc.Ledger.RecordKeyAddition(sc.ChannelID, participant, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log key addition: %v", err)
	}

	fmt.Printf("Public key for participant %s added to security state channel %s\n", participant, sc.ChannelID)
	return nil
}

// RetrieveEncryptedState retrieves encrypted state data for auditing
func (sc *common.SecurityStateChannel) RetrieveEncryptedState(key string) (string, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	encryptedData, exists := sc.State[key]
	if !exists {
		return "", fmt.Errorf("no encrypted data found for key %s", key)
	}

	fmt.Printf("Retrieved encrypted state data for key %s in security state channel %s\n", key, sc.ChannelID)
	return fmt.Sprintf("%v", encryptedData), nil
}

// RetrieveState retrieves the current state of the channel securely
func (sc *common.SecurityStateChannel) RetrieveState(key string) (interface{}, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	value, exists := sc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from security channel %s: %s = %v\n", sc.ChannelID, key, value)
	return value, nil
}
