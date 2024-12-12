package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewOffChainSettlementChannel initializes a new off-chain settlement channel
func NewOffChainSettlementChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.OffChainSettlementChannel {
	return &common.OffChainSettlementChannel{
		ChannelID:    channelID,
		Participants: participants,
		State:        make(map[string]interface{}),
		Balances:     make(map[string]float64),
		IsOpen:       true,
		Ledger:       ledgerInstance,
		Encryption:   encryptionService,
	}
}

// OpenChannel opens the off-chain settlement state channel
func (sc *common.OffChainSettlementChannel) OpenChannel() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.IsOpen {
		return errors.New("settlement channel is already open")
	}

	sc.IsOpen = true

	// Log channel opening in the ledger
	err := sc.Ledger.RecordChannelOpening(sc.ChannelID, sc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Off-chain settlement channel %s opened with participants: %v\n", sc.ChannelID, sc.Participants)
	return nil
}

// CloseChannel closes the off-chain settlement state channel and settles final balances
func (sc *common.OffChainSettlementChannel) CloseChannel() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("settlement channel is already closed")
	}

	sc.IsOpen = false

	// Log channel closure in the ledger
	err := sc.Ledger.RecordChannelClosure(sc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Off-chain settlement channel %s closed\n", sc.ChannelID)
	return nil
}

// AddFunds adds funds to a participant's balance in the settlement channel
func (sc *common.OffChainSettlementChannel) AddFunds(participant string, amount float64) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("settlement channel is closed")
	}

	// Update participant's balance
	sc.Balances[participant] += amount

	// Log the funds addition in the ledger
	err := sc.Ledger.RecordFundsAddition(sc.ChannelID, participant, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log funds addition: %v", err)
	}

	fmt.Printf("Participant %s added %f funds to settlement channel %s\n", participant, amount, sc.ChannelID)
	return nil
}

// MakeSettlement processes an off-chain settlement between participants
func (sc *common.OffChainSettlementChannel) MakeSettlement(sender string, receiver string, amount float64) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("settlement channel is closed")
	}

	// Check if the sender has enough balance
	if sc.Balances[sender] < amount {
		return fmt.Errorf("insufficient funds: sender %s has %f but tried to send %f", sender, sc.Balances[sender], amount)
	}

	// Process settlement
	sc.Balances[sender] -= amount
	sc.Balances[receiver] += amount

	// Log the settlement in the ledger
	err := sc.Ledger.RecordPayment(sc.ChannelID, sender, receiver, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log settlement: %v", err)
	}

	fmt.Printf("Settlement of %f from %s to %s recorded in channel %s\n", amount, sender, receiver, sc.ChannelID)
	return nil
}

// RetrieveBalance retrieves the balance of a participant in the settlement channel
func (sc *common.OffChainSettlementChannel) RetrieveBalance(participant string) (float64, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	balance, exists := sc.Balances[participant]
	if !exists {
		return 0, fmt.Errorf("participant %s not found in channel", participant)
	}

	fmt.Printf("Participant %s has a balance of %f in settlement channel %s\n", participant, balance, sc.ChannelID)
	return balance, nil
}

// UpdateState updates the internal state of the settlement channel
func (sc *common.OffChainSettlementChannel) UpdateState(key string, value interface{}) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("settlement channel is closed")
	}

	// Update the channel state
	sc.State[key] = value

	// Log the state update in the ledger
	err := sc.Ledger.RecordStateUpdate(sc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of settlement channel %s updated: %s = %v\n", sc.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the settlement channel
func (sc *common.OffChainSettlementChannel) RetrieveState(key string) (interface{}, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	value, exists := sc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from settlement channel %s: %s = %v\n", sc.ChannelID, key, value)
	return value, nil
}
