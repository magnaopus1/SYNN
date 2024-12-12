package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewStateChannel initializes a new state channel
func NewStateChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.StateChannel {
	return &common.StateChannel{
		ChannelID:    channelID,
		Participants: participants,
		State:        make(map[string]interface{}),
		IsOpen:       true,
		Ledger:       ledgerInstance,
		Encryption:   encryptionService,
	}
}

// OpenChannel creates and opens a new state channel
func (sc *common.StateChannel) OpenChannel() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.IsOpen {
		return errors.New("state channel is already open")
	}

	sc.IsOpen = true

	// Log channel opening in the ledger
	err := sc.Ledger.RecordChannelOpening(sc.ChannelID, sc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("State channel %s opened with participants: %v\n", sc.ChannelID, sc.Participants)
	return nil
}

// CloseChannel closes the state channel and settles final balances
func (sc *common.StateChannel) CloseChannel() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("state channel is already closed")
	}

	sc.IsOpen = false

	// Log channel closure in the ledger
	err := sc.Ledger.RecordChannelClosure(sc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("State channel %s closed\n", sc.ChannelID)
	return nil
}

// UpdateState updates the internal state of the channel
func (sc *common.StateChannel) UpdateState(key string, value interface{}) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Update the channel state
	sc.State[key] = value

	// Log the state update in the ledger
	err := sc.Ledger.RecordStateUpdate(sc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of channel %s updated: %s = %v\n", sc.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the channel
func (sc *common.StateChannel) RetrieveState(key string) (interface{}, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	value, exists := sc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from channel %s: %s = %v\n", sc.ChannelID, key, value)
	return value, nil
}

// FinalizeTransaction marks a transaction as finalized within the state channel
func (sc *common.StateChannel) FinalizeTransaction(txID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Mark the transaction as finalized (no blocks are involved)
	// Encrypt the transaction data (assuming each tx is in state map)
	txData, exists := sc.State[txID]
	if !exists {
		return fmt.Errorf("transaction %s not found", txID)
	}

	encryptedTx, err := sc.Encryption.EncryptData([]byte(txID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction data: %v", err)
	}

	sc.State[txID] = string(encryptedTx)

	// Log the transaction finalization
	err = sc.Ledger.RecordTransactionFinalized(txID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction finalization: %v", err)
	}

	fmt.Printf("Transaction %s finalized in state channel %s\n", txID, sc.ChannelID)
	return nil
}

