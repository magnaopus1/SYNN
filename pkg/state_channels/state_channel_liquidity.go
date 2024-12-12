package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewLiquidityStateChannel initializes a new liquidity state channel
func NewLiquidityStateChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.LiquidityStateChannel {
	return &common.LiquidityStateChannel{
		ChannelID:    channelID,
		Participants: participants,
		Liquidity:    make(map[string]float64),
		State:        make(map[string]interface{}),
		IsOpen:       true,
		Ledger:       ledgerInstance,
		Encryption:   encryptionService,
	}
}

// AddLiquidity allows a participant to add liquidity to the state channel
func (lc *common.LiquidityStateChannel) AddLiquidity(participant string, amount float64) error {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if !lc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Update participant's liquidity
	lc.Liquidity[participant] += amount

	// Log the liquidity addition in the ledger
	err := lc.Ledger.RecordLiquidityAddition(lc.ChannelID, participant, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log liquidity addition: %v", err)
	}

	fmt.Printf("Participant %s added %f liquidity to channel %s\n", participant, amount, lc.ChannelID)
	return nil
}

// RemoveLiquidity allows a participant to remove liquidity from the state channel
func (lc *common.LiquidityStateChannel) RemoveLiquidity(participant string, amount float64) error {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if !lc.IsOpen {
		return errors.New("state channel is closed")
	}

	if lc.Liquidity[participant] < amount {
		return errors.New("insufficient liquidity for removal")
	}

	// Update participant's liquidity
	lc.Liquidity[participant] -= amount

	// Log the liquidity removal in the ledger
	err := lc.Ledger.RecordLiquidityRemoval(lc.ChannelID, participant, amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log liquidity removal: %v", err)
	}

	fmt.Printf("Participant %s removed %f liquidity from channel %s\n", participant, amount, lc.ChannelID)
	return nil
}

// CloseChannel closes the liquidity state channel and settles balances
func (lc *common.LiquidityStateChannel) CloseChannel() error {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if !lc.IsOpen {
		return errors.New("state channel is already closed")
	}

	lc.IsOpen = false

	// Log channel closure in the ledger
	err := lc.Ledger.RecordChannelClosure(lc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Liquidity state channel %s closed\n", lc.ChannelID)
	return nil
}

// RetrieveLiquidity retrieves the liquidity of a participant in the channel
func (lc *common.LiquidityStateChannel) RetrieveLiquidity(participant string) (float64, error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	liquidity, exists := lc.Liquidity[participant]
	if !exists {
		return 0, fmt.Errorf("participant %s not found in channel", participant)
	}

	fmt.Printf("Participant %s has %f liquidity in channel %s\n", participant, liquidity, lc.ChannelID)
	return liquidity, nil
}

// UpdateState updates the internal state of the liquidity channel
func (lc *common.LiquidityStateChannel) UpdateState(key string, value interface{}) error {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if !lc.IsOpen {
		return errors.New("state channel is closed")
	}

	// Update the channel state
	lc.State[key] = value

	// Log the state update in the ledger
	err := lc.Ledger.RecordStateUpdate(lc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of liquidity channel %s updated: %s = %v\n", lc.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the liquidity channel
func (lc *common.LiquidityStateChannel) RetrieveState(key string) (interface{}, error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	value, exists := lc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from liquidity channel %s: %s = %v\n", lc.ChannelID, key, value)
	return value, nil
}

