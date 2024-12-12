package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/validators"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/common"
)

const (
	EnergyEfficiencyCheckInterval = 10 * time.Minute // Interval to check validator energy efficiency
	EnergyThreshold               = 0.75             // Threshold for triggering energy efficiency optimizations
	EnergyEfficiencyRatingLimit   = 0.85             // Rating below which optimizations are triggered
)

// ValidatorEnergyEfficiencyAutomation handles monitoring and optimizing validator energy efficiency
type ValidatorEnergyEfficiencyAutomation struct {
	consensusEngine  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	validatorManager *validators.ValidatorManager
	efficiencyMutex  *sync.RWMutex
}

// NewValidatorEnergyEfficiencyAutomation initializes the validator energy efficiency automation
func NewValidatorEnergyEfficiencyAutomation(consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, validatorManager *validators.ValidatorManager, efficiencyMutex *sync.RWMutex) *ValidatorEnergyEfficiencyAutomation {
	return &ValidatorEnergyEfficiencyAutomation{
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		validatorManager: validatorManager,
		efficiencyMutex:  efficiencyMutex,
	}
}

// StartEnergyEfficiencyMonitor starts monitoring validator energy efficiency in a continuous loop
func (automation *ValidatorEnergyEfficiencyAutomation) StartEnergyEfficiencyMonitor() {
	ticker := time.NewTicker(EnergyEfficiencyCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkEnergyEfficiency()
		}
	}()
}

// checkEnergyEfficiency checks energy efficiency ratings for validators and applies optimizations as needed
func (automation *ValidatorEnergyEfficiencyAutomation) checkEnergyEfficiency() {
	automation.efficiencyMutex.Lock()
	defer automation.efficiencyMutex.Unlock()

	// Get the current list of active validators and their energy efficiency ratings
	validatorList := automation.validatorManager.GetActiveValidators()

	for _, validator := range validatorList {
		energyRating := automation.validatorManager.GetEnergyEfficiencyRating(validator.ID)

		// If a validator's energy efficiency rating is below the threshold, trigger optimization
		if energyRating < EnergyEfficiencyRatingLimit {
			fmt.Printf("Validator %s has low energy efficiency: %.2f\n", validator.ID, energyRating)
			automation.optimizeValidatorEfficiency(validator.ID, energyRating)
			automation.logEnergyEvent(validator.ID, energyRating)
		}
	}
}

// optimizeValidatorEfficiency applies optimizations for a validator to improve energy efficiency
func (automation *ValidatorEnergyEfficiencyAutomation) optimizeValidatorEfficiency(validatorID string, currentRating float64) {
	fmt.Printf("Optimizing energy efficiency for Validator %s. Current rating: %.2f\n", validatorID, currentRating)

	// Call the ValidatorManager to apply energy-saving optimizations
	err := automation.validatorManager.OptimizeEnergyEfficiency(validatorID)
	if err != nil {
		fmt.Printf("Error optimizing energy efficiency for Validator %s: %v\n", validatorID, err)
	}
}

// logEnergyEvent logs the energy efficiency event in the ledger for auditing
func (automation *ValidatorEnergyEfficiencyAutomation) logEnergyEvent(validatorID string, energyRating float64) {
	entryDetails := fmt.Sprintf("Validator %s had an energy efficiency rating of %.2f, optimization triggered.", validatorID, energyRating)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := ledger.LedgerEntry{
		ID:        fmt.Sprintf("energy-efficiency-%s-%d", validatorID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Energy Efficiency Optimization",
		Status:    "Triggered",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log energy efficiency event in the ledger for Validator %s: %v\n", validatorID, err)
	} else {
		fmt.Println("Energy efficiency event successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ValidatorEnergyEfficiencyAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualEnergyEfficiencyOptimization allows administrators to manually trigger energy efficiency optimization
func (automation *ValidatorEnergyEfficiencyAutomation) TriggerManualEnergyEfficiencyOptimization(validatorID string) {
	fmt.Printf("Manually optimizing energy efficiency for Validator %s.\n", validatorID)
	automation.optimizeValidatorEfficiency(validatorID, automation.validatorManager.GetEnergyEfficiencyRating(validatorID))
}
