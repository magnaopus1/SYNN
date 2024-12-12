package layer2_consensus

import (
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewDynamicConsensusManager initializes the dynamic consensus hopping manager
func NewDynamicConsensusManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *DynamicConsensusManager {
	return &DynamicConsensusManager{
		Strategies:        make(map[string]*ConsensusStrategy),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// AddConsensusStrategy adds a new consensus strategy to the manager
func (dcm *DynamicConsensusManager) AddConsensusStrategy(strategyID, strategyType string) (*ConsensusStrategy, error) {
	dcm.mu.Lock()
	defer dcm.mu.Unlock()

	// Encrypt strategy data
	strategyData := fmt.Sprintf("StrategyID: %s, Type: %s", strategyID, strategyType)

	// Encrypt the data with the required arguments (remove iv)
	encryptedData, err := dcm.EncryptionService.EncryptData(strategyID, []byte(strategyData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt strategy data: %v", err)
	}

	// Use or store the encrypted data (at least log it to avoid the unused variable error)
	fmt.Printf("Encrypted strategy data: %x\n", encryptedData)

	// Create the new consensus strategy
	strategy := &ConsensusStrategy{
		StrategyID:   strategyID,
		StrategyType: strategyType,
		CurrentUsage: 0,
		LastHopped:   time.Now(),
		Active:       false,
		HopCount:     0,
	}

	// Add strategy to the manager
	dcm.Strategies[strategyID] = strategy

	// Log the new strategy addition in the ledger
	dcm.Ledger.BlockchainConsensusCoinLedger.RecordStrategyAddition(strategyID, strategyType) // Removed the third argument

	fmt.Printf("Consensus strategy %s of type %s added\n", strategyID, strategyType)
	return strategy, nil
}



// ActivateConsensusStrategy switches the active consensus mechanism to the specified strategy
func (dcm *DynamicConsensusManager) ActivateConsensusStrategy(strategyID string) error {
	dcm.mu.Lock()
	defer dcm.mu.Unlock()

	// Retrieve the consensus strategy
	strategy, exists := dcm.Strategies[strategyID]
	if !exists {
		return fmt.Errorf("consensus strategy %s not found", strategyID)
	}

	// Deactivate the currently active strategy, if any
	if dcm.ActiveStrategy != nil {
		dcm.ActiveStrategy.Active = false
	}

	// Activate the new strategy
	strategy.Active = true
	strategy.HopCount++
	strategy.LastHopped = time.Now()
	dcm.ActiveStrategy = strategy

	// Log the strategy hop in the ledger (no need to assign to err)
	dcm.Ledger.BlockchainConsensusCoinLedger.RecordStrategyHop(strategyID, strategy.StrategyType)

	fmt.Printf("Consensus strategy %s is now active\n", strategyID)
	return nil
}



// MonitorStrategyUsage monitors the load on a specific consensus strategy and triggers a hop if needed
func (dcm *DynamicConsensusManager) MonitorStrategyUsage(strategyID string, currentUsage float64) error {
	dcm.mu.Lock()
	defer dcm.mu.Unlock()

	// Retrieve the consensus strategy
	strategy, exists := dcm.Strategies[strategyID]
	if !exists {
		return fmt.Errorf("consensus strategy %s not found", strategyID)
	}

	// Update the strategy's current load
	strategy.CurrentUsage = currentUsage

	// Trigger a hop if the load exceeds the threshold
	if currentUsage > 0.75 {
		fmt.Printf("Load on %s is high (%f), considering a hop...\n", strategy.StrategyType, currentUsage)
		for id, strat := range dcm.Strategies {
			if strat.CurrentUsage < 0.5 && strat.StrategyID != strategyID {
				return dcm.ActivateConsensusStrategy(id)
			}
		}
	}

	return nil
}

// GetActiveConsensusStrategy returns the currently active consensus strategy
func (dcm *DynamicConsensusManager) GetActiveConsensusStrategy() (*ConsensusStrategy, error) {
	dcm.mu.Lock()
	defer dcm.mu.Unlock()

	if dcm.ActiveStrategy == nil {
		return nil, errors.New("no active consensus strategy")
	}

	return dcm.ActiveStrategy, nil
}

