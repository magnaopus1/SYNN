package sustainability

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewRenewableEnergyIntegrationSystem initializes the renewable energy integration system
func NewRenewableEnergyIntegrationSystem(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *RenewableEnergyIntegrationSystem {
	return &RenewableEnergyIntegrationSystem{
		EnergySources:     make(map[string]*RenewableEnergySource),
		TotalEnergy:       0,
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// RegisterEnergySource registers a new renewable energy source in the system
func (reis *RenewableEnergyIntegrationSystem) RegisterEnergySource(sourceID, sourceType string, energyProduced float64) (*RenewableEnergySource, error) {
	reis.mu.Lock()
	defer reis.mu.Unlock()

	// Encrypt energy source data
	sourceData := fmt.Sprintf("SourceID: %s, SourceType: %s, EnergyProduced: %f", sourceID, sourceType, energyProduced)
	encryptedData, err := reis.EncryptionService.EncryptData([]byte(sourceData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt energy source data: %v", err)
	}

	// Register the energy source
	source := &RenewableEnergySource{
		SourceID:        sourceID,
		SourceType:      sourceType,
		EnergyProduced:  energyProduced,
		IntegrationDate: time.Now(),
	}
	reis.EnergySources[sourceID] = source
	reis.TotalEnergy += energyProduced

	// Log the energy source registration in the ledger
	err = reis.Ledger.RecordRenewableEnergySourceRegistration(sourceID, sourceType, energyProduced, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log renewable energy source registration: %v", err)
	}

	fmt.Printf("Renewable energy source %s (%s) registered with energy production of %f kWh\n", sourceID, sourceType, energyProduced)
	return source, nil
}

// TrackEnergyContribution tracks the energy contributed by renewable sources over time
func (reis *RenewableEnergyIntegrationSystem) TrackEnergyContribution(sourceID string, additionalEnergy float64) error {
	reis.mu.Lock()
	defer reis.mu.Unlock()

	// Retrieve the energy source
	source, exists := reis.EnergySources[sourceID]
	if !exists {
		return fmt.Errorf("renewable energy source %s not found", sourceID)
	}

	// Update the energy contribution
	source.EnergyProduced += additionalEnergy
	reis.TotalEnergy += additionalEnergy

	// Log the energy contribution in the ledger
	err := reis.Ledger.RecordRenewableEnergyContribution(sourceID, additionalEnergy, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log renewable energy contribution: %v", err)
	}

	fmt.Printf("Additional energy of %f kWh added to source %s. Total energy produced: %f kWh\n", additionalEnergy, sourceID, source.EnergyProduced)
	return nil
}

// ViewTotalRenewableEnergy returns the total renewable energy contributed to the network
func (reis *RenewableEnergyIntegrationSystem) ViewTotalRenewableEnergy() float64 {
	reis.mu.Lock()
	defer reis.mu.Unlock()

	return reis.TotalEnergy
}

// GenerateEnergyImpactReport generates a report on the impact of renewable energy integration
func (reis *RenewableEnergyIntegrationSystem) GenerateEnergyImpactReport() {
	reis.mu.Lock()
	defer reis.mu.Unlock()

	fmt.Printf("Total renewable energy contributed to the network: %f kWh\n", reis.TotalEnergy)
	for _, source := range reis.EnergySources {
		fmt.Printf("SourceID: %s, SourceType: %s, EnergyProduced: %f kWh\n", source.SourceID, source.SourceType, source.EnergyProduced)
	}
}
