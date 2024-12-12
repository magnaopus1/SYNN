package state_channels

import (
	"errors"
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewDataChannel initializes a new data transfer channel
func NewDataChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.DataChannel {
	return &common.DataChannel{
		ChannelID:    channelID,
		Participants: participants,
		DataState:    make(map[string]interface{}),
		DataTransfers: make(map[string]*common.DataBlock),
		IsOpen:       true,
		Ledger:       ledgerInstance,
		Encryption:   encryptionService,
	}
}

// OpenChannel opens the data channel for data transfers
func (dc *common.DataChannel) OpenChannel() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if dc.IsOpen {
		return errors.New("data channel is already open")
	}

	dc.IsOpen = true

	// Log the channel opening in the ledger
	err := dc.Ledger.RecordChannelOpening(dc.ChannelID, dc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Data channel %s opened with participants: %v\n", dc.ChannelID, dc.Participants)
	return nil
}

// CloseChannel closes the data channel after completing data transfers
func (dc *common.DataChannel) CloseChannel() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if !dc.IsOpen {
		return errors.New("data channel is already closed")
	}

	dc.IsOpen = false

	// Log the channel closure in the ledger
	err := dc.Ledger.RecordChannelClosure(dc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Data channel %s closed\n", dc.ChannelID)
	return nil
}

// TransferData creates and transfers a block of data between participants
func (dc *common.DataChannel) TransferData(blockID string, data string) (*common.DataBlock, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if !dc.IsOpen {
		return nil, errors.New("data channel is closed")
	}

	// Check if the block already exists
	if _, exists := dc.DataTransfers[blockID]; exists {
		return nil, errors.New("data block already exists")
	}

	// Encrypt the data before transfer
	encryptedData, err := dc.Encryption.EncryptData([]byte(data), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}

	// Create the data block
	dataBlock := &common.DataBlock{
		BlockID:   blockID,
		Data:      string(encryptedData),
		Timestamp: time.Now(),
		MerkleRoot: common.GenerateMerkleRoot([]byte(data)), // Assuming a helper function for Merkle root generation
	}

	// Store the data block
	dc.DataTransfers[blockID] = dataBlock

	// Log the data transfer in the ledger
	err = dc.Ledger.RecordDataTransfer(dc.ChannelID, blockID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log data transfer: %v", err)
	}

	fmt.Printf("Data block %s transferred in data channel %s\n", blockID, dc.ChannelID)
	return dataBlock, nil
}

// ValidateDataBlock validates a data block transferred in the channel
func (dc *common.DataChannel) ValidateDataBlock(blockID string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	dataBlock, exists := dc.DataTransfers[blockID]
	if !exists {
		return fmt.Errorf("data block %s not found", blockID)
	}

	// Validate the Merkle root of the data block
	expectedMerkleRoot := common.GenerateMerkleRoot([]byte(dataBlock.Data))
	if dataBlock.MerkleRoot != expectedMerkleRoot {
		return errors.New("data block validation failed due to mismatched Merkle root")
	}

	// Log the data block validation in the ledger
	err := dc.Ledger.RecordDataBlockValidation(dc.ChannelID, blockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log data block validation: %v", err)
	}

	fmt.Printf("Data block %s validated in data channel %s\n", blockID, dc.ChannelID)
	return nil
}

// UpdateDataState updates the state of the data in the channel
func (dc *common.DataChannel) UpdateDataState(key string, value interface{}) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if !dc.IsOpen {
		return errors.New("data channel is closed")
	}

	// Update the data state
	dc.DataState[key] = value

	// Log the state update in the ledger
	err := dc.Ledger.RecordStateUpdate(dc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log data state update: %v", err)
	}

	fmt.Printf("Data state in channel %s updated: %s = %v\n", dc.ChannelID, key, value)
	return nil
}

// RetrieveDataState retrieves the current state of the data in the channel
func (dc *common.DataChannel) RetrieveDataState(key string) (interface{}, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if !dc.IsOpen {
		return nil, errors.New("data channel is closed")
	}

	value, exists := dc.DataState[key]
	if !exists {
		return nil, fmt.Errorf("data state key %s not found", key)
	}

	fmt.Printf("Retrieved data state from channel %s: %s = %v\n", dc.ChannelID, key, value)
	return value, nil
}

// RetrieveDataBlock retrieves a data block by its ID from the channel
func (dc *common.DataChannel) RetrieveDataBlock(blockID string) (*common.DataBlock, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	dataBlock, exists := dc.DataTransfers[blockID]
	if !exists {
		return nil, fmt.Errorf("data block %s not found", blockID)
	}

	fmt.Printf("Retrieved data block %s from data channel %s\n", blockID, dc.ChannelID)
	return dataBlock, nil
}
