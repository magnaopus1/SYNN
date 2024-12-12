package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewInstantFinalityChannel initializes a new Instant Finality Channel (IFC)
func NewInstantFinalityChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.InstantFinalityChannel {
	return &common.InstantFinalityChannel{
		ChannelID:      channelID,
		Participants:   participants,
		State:          make(map[string]interface{}),
		Transactions:   make([]*common.Transaction, 0),
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
		Finalized:      false,
	}
}

// OpenChannel opens the Instant Finality Channel
func (ifc *common.InstantFinalityChannel) OpenChannel() error {
	ifc.mu.Lock()
	defer ifc.mu.Unlock()

	if ifc.IsOpen {
		return errors.New("channel is already open")
	}

	ifc.IsOpen = true

	// Log channel opening in the ledger
	err := ifc.Ledger.RecordChannelOpening(ifc.ChannelID, ifc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Instant Finality Channel %s opened with participants: %v\n", ifc.ChannelID, ifc.Participants)
	return nil
}

// CloseChannel closes the Instant Finality Channel and settles final balances
func (ifc *common.InstantFinalityChannel) CloseChannel() error {
	ifc.mu.Lock()
	defer ifc.mu.Unlock()

	if !ifc.IsOpen {
		return errors.New("channel is already closed")
	}

	// Ensure finality is reached before closing
	if !ifc.Finalized {
		err := ifc.ReachFinality()
		if err != nil {
			return fmt.Errorf("failed to achieve finality before closing: %v", err)
		}
	}

	ifc.IsOpen = false

	// Log channel closure in the ledger
	err := ifc.Ledger.RecordChannelClosure(ifc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Instant Finality Channel %s closed\n", ifc.ChannelID)
	return nil
}

// AddTransaction adds a new transaction to the Instant Finality Channel
func (ifc *common.InstantFinalityChannel) AddTransaction(tx *common.Transaction) error {
	ifc.mu.Lock()
	defer ifc.mu.Unlock()

	if !ifc.IsOpen {
		return errors.New("channel is closed")
	}

	// Encrypt the transaction data
	encryptedTx, err := ifc.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	ifc.Transactions = append(ifc.Transactions, tx)

	// Log the transaction in the ledger
	err = ifc.Ledger.RecordTransaction(ifc.ChannelID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction: %v", err)
	}

	fmt.Printf("Transaction %s added to Instant Finality Channel %s\n", tx.TxID, ifc.ChannelID)
	return nil
}

// ReachFinality achieves consensus and marks the state as finalized
func (ifc *common.InstantFinalityChannel) ReachFinality() error {
	ifc.mu.Lock()
	defer ifc.mu.Unlock()

	if ifc.Finalized {
		return errors.New("channel already finalized")
	}

	// Simulate finality by instantly confirming all transactions
	for _, tx := range ifc.Transactions {
		fmt.Printf("Transaction %s in channel %s is confirmed.\n", tx.TxID, ifc.ChannelID)
	}

	// Set finality
	ifc.Finalized = true

	// Log finality achievement in the ledger
	err := ifc.Ledger.RecordFinality(ifc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log finality: %v", err)
	}

	fmt.Printf("Channel %s reached instant finality\n", ifc.ChannelID)
	return nil
}

// UpdateState updates the internal state of the Instant Finality Channel
func (ifc *common.InstantFinalityChannel) UpdateState(key string, value interface{}) error {
	ifc.mu.Lock()
	defer ifc.mu.Unlock()

	if !ifc.IsOpen {
		return errors.New("channel is closed")
	}

	if ifc.Finalized {
		return errors.New("cannot update state after finality is reached")
	}

	// Update the channel state
	ifc.State[key] = value

	// Log the state update in the ledger
	err := ifc.Ledger.RecordStateUpdate(ifc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of channel %s updated: %s = %v\n", ifc.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the Instant Finality Channel
func (ifc *common.InstantFinalityChannel) RetrieveState(key string) (interface{}, error) {
	ifc.mu.Lock()
	defer ifc.mu.Unlock()

	if !ifc.IsOpen {
		return nil, errors.New("channel is closed")
	}

	value, exists := ifc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from channel %s: %s = %v\n", ifc.ChannelID, key, value)
	return value, nil
}

// FinalizeTransaction marks a transaction as finalized
func (ifc *common.InstantFinalityChannel) FinalizeTransaction(txID string) error {
	ifc.mu.Lock()
	defer ifc.mu.Unlock()

	if !ifc.IsOpen {
		return errors.New("channel is closed")
	}

	// Search for the transaction and finalize it
	for _, tx := range ifc.Transactions {
		if tx.TxID == txID {
			fmt.Printf("Transaction %s in channel %s is finalized.\n", txID, ifc.ChannelID)
			// Log the transaction finalization
			err := ifc.Ledger.RecordTransactionFinalization(ifc.ChannelID, txID, time.Now())
			if err != nil {
				return fmt.Errorf("failed to log transaction finalization: %v", err)
			}
			return nil
		}
	}

	return fmt.Errorf("transaction %s not found", txID)
}

// RetrieveTransaction retrieves a transaction by its ID from the Instant Finality Channel
func (ifc *common.InstantFinalityChannel) RetrieveTransaction(txID string) (*common.Transaction, error) {
	ifc.mu.Lock()
	defer ifc.mu.Unlock()

	for _, tx := range ifc.Transactions {
		if tx.TxID == txID {
			fmt.Printf("Retrieved transaction %s from channel %s\n", txID, ifc.ChannelID)
			return tx, nil
		}
	}

	return nil, fmt.Errorf("transaction %s not found", txID)
}
