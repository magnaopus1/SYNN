package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewPaymentStateChannel initializes a new payment state channel
func NewPaymentStateChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.PaymentStateChannel {
	return &common.PaymentStateChannel{
		ChannelID:    channelID,
		Participants: participants,
		Balances:     make(map[string]float64),
		State:        make(map[string]interface{}),
		IsOpen:       true,
		Ledger:       ledgerInstance,
		Encryption:   encryptionService,
	}
}

// OpenChannel opens the payment state channel
func (pc *common.PaymentStateChannel) OpenChannel() error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.IsOpen {
		return errors.New("payment channel is already open")
	}

	pc.IsOpen = true

	// Log channel opening in the ledger
	err := pc.Ledger.RecordChannelOpening(pc.ChannelID, pc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Payment channel %s opened with participants: %v\n", pc.ChannelID, pc.Participants)
	return nil
}

// CloseChannel closes the payment state channel and settles final balances
func (pc *common.PaymentStateChannel) CloseChannel() error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if !pc.IsOpen {
		return errors.New("payment channel is already closed")
	}

	pc.IsOpen = false

	// Log channel closure in the ledger
	err := pc.Ledger.RecordChannelClosure(pc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Payment channel %s closed\n", pc.ChannelID)
	return nil
}

// MakePayment processes a payment between participants within the state channel
func (pc *common.PaymentStateChannel) MakePayment(sender string, receiver string, amount float64) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if !pc.IsOpen {
		return errors.New("payment channel is closed")
	}

	// Check if the sender has enough balance
	if pc.Balances[sender] < amount {
		return fmt.Errorf("insufficient funds: sender %s has %f but tried to send %f", sender, pc.Balances[sender], amount)
	}

	// Process payment
	pc.Balances[sender] -= amount
	pc.Balances[receiver] += amount

	// Log the payment in the ledger
	err := pc.Ledger.RecordPayment(pc.ChannelID, sender, receiver, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log payment: %v", err)
	}

	fmt.Printf("Payment of %f from %s to %s recorded in channel %s\n", amount, sender, receiver, pc.ChannelID)
	return nil
}

// AddFunds adds funds to a participant's balance in the payment state channel
func (pc *common.PaymentStateChannel) AddFunds(participant string, amount float64) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if !pc.IsOpen {
		return errors.New("payment channel is closed")
	}

	// Update participant's balance
	pc.Balances[participant] += amount

	// Log the funds addition in the ledger
	err := pc.Ledger.RecordFundsAddition(pc.ChannelID, participant, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log funds addition: %v", err)
	}

	fmt.Printf("Participant %s added %f funds to channel %s\n", participant, amount, pc.ChannelID)
	return nil
}

// RetrieveBalance retrieves the current balance of a participant in the payment channel
func (pc *common.PaymentStateChannel) RetrieveBalance(participant string) (float64, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	balance, exists := pc.Balances[participant]
	if !exists {
		return 0, fmt.Errorf("participant %s not found in channel", participant)
	}

	fmt.Printf("Participant %s has a balance of %f in channel %s\n", participant, balance, pc.ChannelID)
	return balance, nil
}

// UpdateState updates the internal state of the payment channel
func (pc *common.PaymentStateChannel) UpdateState(key string, value interface{}) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if !pc.IsOpen {
		return errors.New("payment channel is closed")
	}

	// Update the channel state
	pc.State[key] = value

	// Log the state update in the ledger
	err := pc.Ledger.RecordStateUpdate(pc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of payment channel %s updated: %s = %v\n", pc.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the payment channel
func (pc *common.PaymentStateChannel) RetrieveState(key string) (interface{}, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	value, exists := pc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from payment channel %s: %s = %v\n", pc.ChannelID, key, value)
	return value, nil
}
