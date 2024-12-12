package sustainability

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewEnergyUsageMonitoringSystem initializes a new energy usage monitoring system
func NewEnergyUsageMonitoringSystem(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *EnergyUsageMonitoringSystem {
	return &EnergyUsageMonitoringSystem{
		UsageRecords:      make(map[string]*EnergyUsageRecord),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// LogEnergyUsage logs the energy consumption of a node for a specific time period
func (eums *EnergyUsageMonitoringSystem) LogEnergyUsage(recordID, nodeID, owner string, energyUsage float64, periodStart, periodEnd time.Time) (*EnergyUsageRecord, error) {
	eums.mu.Lock()
	defer eums.mu.Unlock()

	// Encrypt energy usage data
	usageData := fmt.Sprintf("RecordID: %s, NodeID: %s, Owner: %s, EnergyUsage: %f, PeriodStart: %s, PeriodEnd: %s", recordID, nodeID, owner, energyUsage, periodStart, periodEnd)
	encryptedData, err := eums.EncryptionService.EncryptData([]byte(usageData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt energy usage data: %v", err)
	}

	// Create the energy usage record
	record := &EnergyUsageRecord{
		RecordID:    recordID,
		NodeID:      nodeID,
		Owner:       owner,
		EnergyUsage: energyUsage,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		LoggedTime:  time.Now(),
	}

	// Add the record to the system
	eums.UsageRecords[recordID] = record

	// Log the energy usage in the ledger
	err = eums.Ledger.RecordEnergyUsage(recordID, nodeID, owner, energyUsage, periodStart, periodEnd, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log energy usage in the ledger: %v", err)
	}

	fmt.Printf("Energy usage of %f kWh logged for node %s (owner: %s) for the period from %s to %s\n", energyUsage, nodeID, owner, periodStart, periodEnd)
	return record, nil
}

// ViewEnergyUsage allows viewing of the energy usage details for a specific record
func (eums *EnergyUsageMonitoringSystem) ViewEnergyUsage(recordID string) (*EnergyUsageRecord, error) {
	eums.mu.Lock()
	defer eums.mu.Unlock()

	// Retrieve the energy usage record
	record, exists := eums.UsageRecords[recordID]
	if !exists {
		return nil, fmt.Errorf("energy usage record %s not found", recordID)
	}

	return record, nil
}

// UpdateEnergyUsage allows updating of an existing energy usage record if corrections are needed
func (eums *EnergyUsageMonitoringSystem) UpdateEnergyUsage(recordID string, newEnergyUsage float64, newPeriodEnd time.Time) (*EnergyUsageRecord, error) {
	eums.mu.Lock()
	defer eums.mu.Unlock()

	// Retrieve the energy usage record
	record, exists := eums.UsageRecords[recordID]
	if !exists {
		return nil, fmt.Errorf("energy usage record %s not found", recordID)
	}

	// Update the record
	oldUsage := record.EnergyUsage
	record.EnergyUsage = newEnergyUsage
	record.PeriodEnd = newPeriodEnd

	// Log the update in the ledger
	err := eums.Ledger.RecordEnergyUsageUpdate(recordID, record.NodeID, record.Owner, oldUsage, newEnergyUsage, record.PeriodStart, newPeriodEnd, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log energy usage update in the ledger: %v", err)
	}

	fmt.Printf("Energy usage record %s updated for node %s: new usage %f kWh (old: %f kWh)\n", recordID, record.NodeID, newEnergyUsage, oldUsage)
	return record, nil
}
