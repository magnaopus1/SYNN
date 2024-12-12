package rollups

import (
	"errors"
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewSelfGoverningRollupEcosystem initializes a new self-governing rollup ecosystem.
func NewSelfGoverningRollupEcosystem(ecosystemID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, consensus *common.SynnergyConsensus) *common.SelfGoverningRollupEcosystem {
	// Initialize with default governance rules
	defaultRules := &common.GovernanceRules{
		MaxTransactionsPerRollup: 1000,  // Default max transactions per rollup
		FeeStructure:             0.01,  // Default fee structure (1%)
		ScalingFactor:            1.0,   // Default scaling factor
		LastUpdated:              time.Now(),
	}

	return &common.SelfGoverningRollupEcosystem{
		EcosystemID:    ecosystemID,
		Rollups:        make(map[string]*common.Rollup),
		GovernanceRules: defaultRules,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
		Consensus:      consensus,
	}
}

// AddRollup adds a new rollup to the self-governing ecosystem.
func (sgre *common.SelfGoverningRollupEcosystem) AddRollup(rollupID string) error {
	sgre.mu.Lock()
	defer sgre.mu.Unlock()

	if _, exists := sgre.Rollups[rollupID]; exists {
		return errors.New("rollup already exists in the ecosystem")
	}

	newRollup := NewRollup(rollupID, sgre.Ledger, sgre.Encryption, sgre.NetworkManager)
	sgre.Rollups[rollupID] = newRollup

	// Log the addition of a new rollup
	err := sgre.Ledger.RecordRollupAddition(sgre.EcosystemID, rollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup addition: %v", err)
	}

	fmt.Printf("Rollup %s added to ecosystem %s\n", rollupID, sgre.EcosystemID)
	return nil
}

// UpdateGovernanceRules allows the ecosystem to update the self-regulating governance parameters.
func (sgre *common.SelfGoverningRollupEcosystem) UpdateGovernanceRules(maxTxPerRollup int, fee float64, scalingFactor float64) error {
	sgre.mu.Lock()
	defer sgre.mu.Unlock()

	// Apply the new rules
	sgre.GovernanceRules.MaxTransactionsPerRollup = maxTxPerRollup
	sgre.GovernanceRules.FeeStructure = fee
	sgre.GovernanceRules.ScalingFactor = scalingFactor
	sgre.GovernanceRules.LastUpdated = time.Now()

	// Log the update in the ledger
	err := sgre.Ledger.RecordGovernanceUpdate(sgre.EcosystemID, maxTxPerRollup, fee, scalingFactor, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log governance rule update: %v", err)
	}

	fmt.Printf("Governance rules updated for ecosystem %s: MaxTx: %d, Fee: %.2f, ScalingFactor: %.2f\n", sgre.EcosystemID, maxTxPerRollup, fee, scalingFactor)
	return nil
}

// ApplyGovernanceRules applies the governance rules to a rollup and validates it using consensus.
func (sgre *common.SelfGoverningRollupEcosystem) ApplyGovernanceRules(rollupID string) error {
	sgre.mu.Lock()
	defer sgre.mu.Unlock()

	rollup, exists := sgre.Rollups[rollupID]
	if !exists {
		return fmt.Errorf("rollup %s not found in ecosystem %s", rollupID, sgre.EcosystemID)
	}

	// Enforce maximum transactions per rollup
	if len(rollup.Transactions) > sgre.GovernanceRules.MaxTransactionsPerRollup {
		return fmt.Errorf("rollup %s exceeds maximum allowed transactions", rollupID)
	}

	// Apply fee structure to each transaction
	for _, tx := range rollup.Transactions {
		tx.Fee = sgre.GovernanceRules.FeeStructure * tx.Amount
	}

	// Use Synnergy Consensus to validate the governance compliance
	valid, err := sgre.Consensus.ValidateGovernanceCompliance(sgre.EcosystemID, rollupID, sgre.GovernanceRules)
	if err != nil || !valid {
		return fmt.Errorf("failed to validate governance rules for rollup %s: %v", rollupID, err)
	}

	// Log the governance application in the ledger
	err = sgre.Ledger.RecordGovernanceApplication(sgre.EcosystemID, rollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log governance application: %v", err)
	}

	fmt.Printf("Governance rules applied and validated for rollup %s in ecosystem %s\n", rollupID, sgre.EcosystemID)
	return nil
}

// RemoveRollup removes a rollup from the self-governing ecosystem.
func (sgre *common.SelfGoverningRollupEcosystem) RemoveRollup(rollupID string) error {
	sgre.mu.Lock()
	defer sgre.mu.Unlock()

	if _, exists := sgre.Rollups[rollupID]; !exists {
		return fmt.Errorf("rollup %s not found in ecosystem %s", rollupID, sgre.EcosystemID)
	}

	delete(sgre.Rollups, rollupID)

	// Log the rollup removal in the ledger
	err := sgre.Ledger.RecordRollupRemoval(sgre.EcosystemID, rollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup removal: %v", err)
	}

	fmt.Printf("Rollup %s removed from ecosystem %s\n", rollupID, sgre.EcosystemID)
	return nil
}

// RetrieveRollup retrieves a rollup by its ID from the ecosystem.
func (sgre *common.SelfGoverningRollupEcosystem) RetrieveRollup(rollupID string) (*common.Rollup, error) {
	sgre.mu.Lock()
	defer sgre.mu.Unlock()

	rollup, exists := sgre.Rollups[rollupID]
	if !exists {
		return nil, fmt.Errorf("rollup %s not found in ecosystem %s", rollupID, sgre.EcosystemID)
	}

	fmt.Printf("Retrieved rollup %s from ecosystem %s\n", rollupID, sgre.EcosystemID)
	return rollup, nil
}

// MonitorGovernance continuously monitors the rollup ecosystem and automatically adjusts parameters based on network conditions.
func (sgre *common.SelfGoverningRollupEcosystem) MonitorGovernance(interval time.Duration) {
	for {
		time.Sleep(interval)

		// Example governance monitoring logic
		// Adjust scaling factor based on the number of active rollups
		sgre.mu.Lock()
		if len(sgre.Rollups) > 100 {
			sgre.GovernanceRules.ScalingFactor *= 1.1
		} else {
			sgre.GovernanceRules.ScalingFactor = 1.0
		}
		sgre.mu.Unlock()

		// Log the monitoring event in the ledger
		err := sgre.Ledger.RecordGovernanceMonitoring(sgre.EcosystemID, sgre.GovernanceRules.ScalingFactor, time.Now())
		if err != nil {
			fmt.Printf("Failed to log governance monitoring event: %v\n", err)
		}

		fmt.Printf("Governance monitoring adjusted scaling factor to %.2f for ecosystem %s\n", sgre.GovernanceRules.ScalingFactor, sgre.EcosystemID)
	}
}
