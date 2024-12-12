package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewFlexibleStateChannel initializes a new flexible state channel
func NewFlexibleStateChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.FlexibleStateChannel {
	return &common.FlexibleStateChannel{
		ChannelID:    channelID,
		Participants: participants,
		State:        make(map[string]interface{}),
		DataTransfers: make(map[string]*common.DataBlock),
		Liquidity:    make(map[string]float64),
		Transactions: make(map[string]*common.Transaction),
		IsOpen:       true,
		Ledger:       ledgerInstance,
		Encryption:   encryptionService,
	}
}

// OpenChannel opens the flexible state channel for various interactions
func (fsc *common.FlexibleStateChannel) OpenChannel() error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if fsc.IsOpen {
		return errors.New("state channel is already open")
	}

	fsc.IsOpen = true

	// Log the channel opening in the ledger
	err := fsc.Ledger.RecordChannelOpening(fsc.ChannelID, fsc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Flexible state channel %s opened with participants: %v\n", fsc.ChannelID, fsc.Participants)
	return nil
}

// CloseChannel closes the flexible state channel after settling all interactions
func (fsc *common.FlexibleStateChannel) CloseChannel() error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return errors.New("state channel is already closed")
	}

	fsc.IsOpen = false

	// Log the channel closure in the ledger
	err := fsc.Ledger.RecordChannelClosure(fsc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Flexible state channel %s closed\n", fsc.ChannelID)
	return nil
}

// AddLiquidity adds liquidity for a participant in the state channel
func (fsc *common.FlexibleStateChannel) AddLiquidity(participant string, amount float64) error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Update participant's liquidity
	fsc.Liquidity[participant] += amount

	// Log the liquidity addition in the ledger
	err := fsc.Ledger.RecordLiquidityAddition(fsc.ChannelID, participant, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log liquidity addition: %v", err)
	}

	fmt.Printf("Participant %s added %f liquidity to flexible channel %s\n", participant, amount, fsc.ChannelID)
	return nil
}

// TransferData creates and transfers a block of data in the channel
func (fsc *common.FlexibleStateChannel) TransferData(blockID string, data string) (*common.DataBlock, error) {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	// Encrypt the data before transfer
	encryptedData, err := fsc.Encryption.EncryptData([]byte(data), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}

	// Create the data block
	dataBlock := &common.DataBlock{
		BlockID:   blockID,
		Data:      string(encryptedData),
		Timestamp: time.Now(),
		MerkleRoot: common.GenerateMerkleRoot([]byte(data)),
	}

	// Store the data block
	fsc.DataTransfers[blockID] = dataBlock

	// Log the data transfer in the ledger
	err = fsc.Ledger.RecordDataTransfer(fsc.ChannelID, blockID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log data transfer: %v", err)
	}

	fmt.Printf("Data block %s transferred in flexible channel %s\n", blockID, fsc.ChannelID)
	return dataBlock, nil
}

// ProcessTransaction adds a new transaction to the state channel
func (fsc *common.FlexibleStateChannel) ProcessTransaction(transactionID string, tx *common.Transaction) error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Check if the transaction already exists
	if _, exists := fsc.Transactions[transactionID]; exists {
		return errors.New("transaction already exists")
	}

	// Encrypt the transaction data
	encryptedTx, err := fsc.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add transaction to the channel
	fsc.Transactions[transactionID] = tx

	// Log the transaction in the ledger
	err = fsc.Ledger.RecordTransaction(fsc.ChannelID, transactionID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction: %v", err)
	}

	fmt.Printf("Transaction %s processed in flexible channel %s\n", transactionID, fsc.ChannelID)
	return nil
}

// UpdateState updates the state of the flexible channel
func (fsc *common.FlexibleStateChannel) UpdateState(key string, value interface{}) error {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	if !fsc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Update the state
	fsc.State[key] = value

	// Log the state update in the ledger
	err := fsc.Ledger.RecordStateUpdate(fsc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of flexible channel %s updated: %s = %v\n", fsc.ChannelID, key, value)
	return nil
}

// RetrieveLiquidity retrieves the liquidity of a participant in the channel
func (fsc *common.FlexibleStateChannel) RetrieveLiquidity(participant string) (float64, error) {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	liquidity, exists := fsc.Liquidity[participant]
	if !exists {
		return 0, fmt.Errorf("participant %s not found in channel", participant)
	}

	fmt.Printf("Participant %s has %f liquidity in flexible channel %s\n", participant, liquidity, fsc.ChannelID)
	return liquidity, nil
}

// RetrieveDataBlock retrieves a data block by its ID from the flexible state channel
func (fsc *common.FlexibleStateChannel) RetrieveDataBlock(blockID string) (*common.DataBlock, error) {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	dataBlock, exists := fsc.DataTransfers[blockID]
	if !exists {
		return nil, fmt.Errorf("data block %s not found", blockID)
	}

	fmt.Printf("Retrieved data block %s from flexible channel %s\n", blockID, fsc.ChannelID)
	return dataBlock, nil
}

// RetrieveTransaction retrieves a transaction by its ID from the flexible state channel
func (fsc *common.FlexibleStateChannel) RetrieveTransaction(transactionID string) (*common.Transaction, error) {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	transaction, exists := fsc.Transactions[transactionID]
	if !exists {
		return nil, fmt.Errorf("transaction %s not found", transactionID)
	}

	fmt.Printf("Retrieved transaction %s from flexible channel %s\n", transactionID, fsc.ChannelID)
	return transaction, nil
}

// RetrieveState retrieves the current state of the channel by key
func (fsc *common.FlexibleStateChannel) RetrieveState(key string) (interface{}, error) {
	fsc.mu.Lock()
	defer fsc.mu.Unlock()

	value, exists := fsc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from flexible channel %s: %s = %v\n", fsc.ChannelID, key, value)
	return value, nil
}
