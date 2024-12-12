package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewStateChannel initializes a new state channel, with optional microchannel capabilities
func NewStateChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.StateChannel {
	return &common.StateChannel{
		ChannelID:    channelID,
		Participants: participants,
		State:        make(map[string]interface{}),
		Transactions: []*common.Transaction{},
		IsOpen:       true,
		Ledger:       ledgerInstance,
		Encryption:   encryptionService,
	}
}

// EnableMicroChannelMode enables microchannel mode on an existing state channel
func (sc *common.StateChannel) EnableMicroChannelMode(maxTransactions int, timeout time.Duration) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("cannot enable microchannel mode on a closed state channel")
	}

	sc.MicroChannel = &common.MicroChannelMode{
		Enabled:        true,
		MaxTransaction: maxTransactions,
		Timeout:        timeout,
		StartTime:      time.Now(),
	}

	fmt.Printf("Microchannel mode enabled on state channel %s with a max of %d transactions and timeout of %v\n", sc.ChannelID, maxTransactions, timeout)
	return nil
}

// AddTransaction adds a transaction to the state channel, handling both normal and microchannel modes
func (sc *common.StateChannel) AddTransaction(tx *common.Transaction) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Encrypt transaction data
	encryptedTx, err := sc.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add the transaction
	sc.Transactions = append(sc.Transactions, tx)

	// Log the transaction in the ledger
	err = sc.Ledger.RecordTransactionAddition(sc.ChannelID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to state channel %s\n", tx.TxID, sc.ChannelID)

	// Check if the channel is in microchannel mode
	if sc.MicroChannel != nil && sc.MicroChannel.Enabled {
		// Close the channel if the maximum transaction count or timeout is reached
		if len(sc.Transactions) >= sc.MicroChannel.MaxTransaction || time.Since(sc.MicroChannel.StartTime) >= sc.MicroChannel.Timeout {
			err := sc.CloseChannel()
			if err != nil {
				return fmt.Errorf("failed to close channel: %v", err)
			}
		}
	}

	return nil
}

// CloseChannel closes the state channel, finalizing transactions or state if in microchannel mode
func (sc *common.StateChannel) CloseChannel() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("state channel is already closed")
	}

	// Finalize state or transactions for microchannel mode
	if sc.MicroChannel != nil && sc.MicroChannel.Enabled {
		err := sc.finalizeMicroChannelState()
		if err != nil {
			return fmt.Errorf("failed to finalize microchannel state: %v", err)
		}
	}

	sc.IsOpen = false

	// Log the channel closure in the ledger
	err := sc.Ledger.RecordChannelClosure(sc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("State channel %s closed\n", sc.ChannelID)
	return nil
}

// finalizeMicroChannelState finalizes the state or transactions of a microchannel before closing
func (sc *common.StateChannel) finalizeMicroChannelState() error {
	// For simplicity, this could involve settling balances, processing transactions, etc.
	// The logic here would depend on your system's requirements.

	// Example: Sum up the transaction values and update state
	total := 0.0
	for _, tx := range sc.Transactions {
		total += tx.Value
	}

	sc.State["MicroChannelTotal"] = total

	// Log the final state in the ledger
	err := sc.Ledger.RecordStateUpdate(sc.ChannelID, "MicroChannelTotal", total, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log final microchannel state: %v", err)
	}

	fmt.Printf("Finalized microchannel state with total: %f\n", total)
	return nil
}

