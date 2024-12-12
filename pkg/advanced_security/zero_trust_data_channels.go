package advanced_security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// ZTDCChannel represents a secure zero-trust data channel
type ZTDCChannel struct {
	ChannelID      string    // Unique ID of the channel
	ParticipantA   string    // Address of the first participant
	ParticipantB   string    // Address of the second participant
	EncryptionKey  string    // Encryption key used for secure communication
	CreatedAt      time.Time // Timestamp of channel creation
	LastActivity   time.Time // Timestamp of the last activity in the channel
	Status         string    // Status of the channel ("Active", "Closed")
	EncryptedData  string    // Most recent encrypted data sent over the channel
}

// ZTDCManager handles the creation and management of zero-trust data channels
type ZTDCManager struct {
	Channels         map[string]*ZTDCChannel   // Active data channels
	Ledger           *ledger.Ledger            // Ledger instance for logging channel events
	EncryptionService *common.Encryption   // Encryption service for secure data handling
	mu               sync.Mutex                // Mutex for concurrent access to channels
}

// NewZTDCManager initializes a new Zero-Trust Data Channels manager
func NewZTDCManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *ZTDCManager {
	return &ZTDCManager{
		Channels:          make(map[string]*ZTDCChannel),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// CreateChannel initializes a new zero-trust data channel between two participants
func (ztdc *ZTDCManager) CreateChannel(participantA, participantB string) (*ZTDCChannel, error) {
	ztdc.mu.Lock()
	defer ztdc.mu.Unlock()

	// Generate a unique channel ID
	channelID := generateUniqueID()

	// Generate an encryption key for the channel
	encryptionKey, err := ztdc.generateEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %v", err)
	}

	// Create the channel
	channel := &ZTDCChannel{
		ChannelID:     channelID,
		ParticipantA:  participantA,
		ParticipantB:  participantB,
		EncryptionKey: encryptionKey,
		CreatedAt:     time.Now(),
		Status:        "Active",
	}

	// Add the channel to the active channel map
	ztdc.Channels[channelID] = channel

	// Log the channel creation in the ledger (combine participants into a []string)
	participants := []string{participantA, participantB}
	successMsg, err := ztdc.Ledger.StateChannelLedger.RecordChannelCreation(channelID, participants)
	if err != nil {
		return nil, fmt.Errorf("failed to log channel creation in the ledger: %v", err)
	}

	// Optionally log the success message
	fmt.Printf("Ledger Record Success: %s\n", successMsg)
	fmt.Printf("Zero-trust data channel created between %s and %s\n", participantA, participantB)

	return channel, nil
}



// SendEncryptedData allows participants to send encrypted data over the zero-trust channel
func (ztdc *ZTDCManager) SendEncryptedData(channelID, sender, data string) error {
	ztdc.mu.Lock()
	defer ztdc.mu.Unlock()

	// Retrieve the channel
	channel, exists := ztdc.Channels[channelID]
	if !exists {
		return fmt.Errorf("channel %s not found", channelID)
	}

	// Ensure the sender is one of the channel participants
	if sender != channel.ParticipantA && sender != channel.ParticipantB {
		return fmt.Errorf("sender %s is not authorized to send data over channel %s", sender, channelID)
	}

	// Encrypt the data with the channel's encryption key (pass the encryption algorithm)
	encryptedData, err := ztdc.EncryptionService.EncryptData("AES", []byte(data), []byte(channel.EncryptionKey))
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %v", err)
	}

	// Convert encryptedData to string if necessary
	encryptedDataStr := string(encryptedData)

	// Update the channel with the new encrypted data
	channel.EncryptedData = encryptedDataStr
	channel.LastActivity = time.Now()

	// Prepare a map for logging the transmission data in the ledger
	transmissionDetails := map[string]interface{}{
		"timestamp": time.Now(),
	}

	// Log the data transmission in the ledger (pass channelID, sender, and transmissionDetails)
	successMsg, err := ztdc.Ledger.DataManagementLedger.RecordDataTransmission(channelID, sender, transmissionDetails)
	if err != nil {
		return fmt.Errorf("failed to log data transmission in the ledger: %v", err)
	}

	// Optionally log the success message
	fmt.Printf("Ledger Record Success: %s\n", successMsg)
	fmt.Printf("Data sent over channel %s by %s\n", channelID, sender)

	return nil
}




// CloseChannel closes the zero-trust data channel
func (ztdc *ZTDCManager) CloseChannel(channelID string) error {
	ztdc.mu.Lock()
	defer ztdc.mu.Unlock()

	// Retrieve the channel
	channel, exists := ztdc.Channels[channelID]
	if !exists {
		return fmt.Errorf("channel %s not found", channelID)
	}

	// Mark the channel as closed
	channel.Status = "Closed"
	channel.LastActivity = time.Now()

	// Log the channel closure in the ledger (only pass channelID)
	err := ztdc.Ledger.StateChannelLedger.RecordChannelClosure(channelID)
	if err != nil {
		return fmt.Errorf("failed to log channel closure in the ledger: %v", err)
	}

	fmt.Printf("Channel %s has been closed\n", channelID)
	return nil
}


// MonitorChannelStatus checks the current status of a specific channel
func (ztdc *ZTDCManager) MonitorChannelStatus(channelID string) (string, error) {
	ztdc.mu.Lock()
	defer ztdc.mu.Unlock()

	// Retrieve the channel
	channel, exists := ztdc.Channels[channelID]
	if !exists {
		return "", fmt.Errorf("channel %s not found", channelID)
	}

	return channel.Status, nil
}

// generateEncryptionKey generates a cryptographically secure encryption key
func (ztdc *ZTDCManager) generateEncryptionKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate encryption key: %v", err)
	}
	return hex.EncodeToString(key), nil
}

// generateUniqueID creates a cryptographically secure unique ID
func generateUniqueID() string {
	id := make([]byte, 16)
	rand.Read(id)
	return hex.EncodeToString(id)
}

