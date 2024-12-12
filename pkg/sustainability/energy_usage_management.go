package sustainability

import (
	"common"
	"ledger"
	"encryption"
	"log"

)

// MeasureCurrentEnergyUsage calculates the current energy usage of a specific node.
func MeasureCurrentEnergyUsage(nodeID string) (float64, error) {
	usage, err := ledger.GetNodeEnergyUsage(nodeID)
	if err != nil {
		log.Printf("Error measuring energy usage for node %s: %v\n", nodeID, err)
		return 0, err
	}
	log.Printf("Current energy usage for node %s: %.2f units.\n", nodeID, usage)
	return usage, nil
}

// SetMaxEnergyUsage sets a maximum limit on energy usage for a node.
func SetMaxEnergyUsage(nodeID string, maxUsage float64) error {
	err := ledger.SetEnergyUsageLimit(nodeID, maxUsage)
	if err != nil {
		log.Printf("Error setting max energy usage for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Max energy usage set to %.2f units for node %s.\n", maxUsage, nodeID)
	return nil
}

// EnableLowPowerMode activates low-power mode on a node to reduce energy usage.
func EnableLowPowerMode(nodeID string) error {
	err := ledger.ActivateLowPowerMode(nodeID)
	if err != nil {
		log.Printf("Error enabling low-power mode for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Low-power mode enabled for node %s.\n", nodeID)
	return nil
}

// DisableLowPowerMode deactivates low-power mode on a node.
func DisableLowPowerMode(nodeID string) error {
	err := ledger.DeactivateLowPowerMode(nodeID)
	if err != nil {
		log.Printf("Error disabling low-power mode for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Low-power mode disabled for node %s.\n", nodeID)
	return nil
}

// ReportEmissionLevels logs the emission levels for a node.
func ReportEmissionLevels(nodeID string, emissions float64) error {
	err := ledger.LogEmissionLevels(nodeID, emissions)
	if err != nil {
		log.Printf("Error reporting emission levels for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Emission levels of %.2f reported for node %s.\n", emissions, nodeID)
	return nil
}

// OffsetEnergyConsumption applies an offset to energy consumption.
func OffsetEnergyConsumption(nodeID string, offsetAmount float64) error {
	err := ledger.ApplyEnergyOffset(nodeID, offsetAmount)
	if err != nil {
		log.Printf("Error offsetting energy consumption for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Energy consumption offset by %.2f units for node %s.\n", offsetAmount, nodeID)
	return nil
}

// AdjustDynamicResourceAllocation modifies resource allocation based on energy use.
func AdjustDynamicResourceAllocation(nodeID string) error {
	resourceMetrics, err := ledger.GetResourceMetrics(nodeID)
	if err != nil {
		log.Printf("Error retrieving resource metrics for node %s: %v\n", nodeID, err)
		return err
	}
	err = ledger.DynamicResourceAllocation(nodeID, resourceMetrics)
	if err != nil {
		log.Printf("Error adjusting resource allocation for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Dynamic resource allocation adjusted for node %s.\n", nodeID)
	return nil
}

// GenerateEnergyUsageReport creates a report of energy usage for a node.
func GenerateEnergyUsageReport(nodeID string) (string, error) {
	reportData, err := ledger.GenerateNodeEnergyReport(nodeID)
	if err != nil {
		log.Printf("Error generating energy usage report for node %s: %v\n", nodeID, err)
		return "", err
	}
	report := formatEnergyReport(reportData)
	log.Printf("Energy usage report generated for node %s.\n", nodeID)
	return report, nil
}

// TrackResourceEfficiency logs the resource efficiency for a node.
func TrackResourceEfficiency(nodeID string, efficiency float64) error {
	err := ledger.RecordResourceEfficiency(nodeID, efficiency)
	if err != nil {
		log.Printf("Error tracking resource efficiency for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Resource efficiency of %.2f logged for node %s.\n", efficiency, nodeID)
	return nil
}

// SetEmissionThreshold sets a threshold for allowable emissions.
func SetEmissionThreshold(nodeID string, threshold float64) error {
	err := ledger.SetEmissionLimit(nodeID, threshold)
	if err != nil {
		log.Printf("Error setting emission threshold for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Emission threshold set to %.2f for node %s.\n", threshold, nodeID)
	return nil
}

// ManualLogEnergyEvent allows a manual entry of an energy event for logging purposes.
func ManualLogEnergyEvent(nodeID, eventDesc string) error {
	encryptedEvent, err := encryption.EncryptString(eventDesc)
	if err != nil {
		log.Printf("Error encrypting energy event: %v\n", err)
		return err
	}
	err = ledger.LogEnergyEvent(nodeID, encryptedEvent)
	if err != nil {
		log.Printf("Error logging energy event for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Energy event manually logged for node %s.\n", nodeID)
	return nil
}

// MeasureRenewableEnergyUsage tracks the amount of renewable energy used.
func MeasureRenewableEnergyUsage(nodeID string) (float64, error) {
	renewableUsage, err := ledger.GetRenewableEnergyUsage(nodeID)
	if err != nil {
		log.Printf("Error measuring renewable energy usage for node %s: %v\n", nodeID, err)
		return 0, err
	}
	log.Printf("Renewable energy usage for node %s: %.2f units.\n", nodeID, renewableUsage)
	return renewableUsage, nil
}

// IncentivizeRenewableUsage provides rewards for nodes using renewable energy.
func IncentivizeRenewableUsage(nodeID string) error {
	err := ledger.GrantRenewableIncentive(nodeID)
	if err != nil {
		log.Printf("Error incentivizing renewable usage for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Renewable usage incentive granted for node %s.\n", nodeID)
	return nil
}

// ReducePeakEnergyConsumption implements measures to reduce energy during peak times.
func ReducePeakEnergyConsumption(nodeID string) error {
	err := ledger.ReducePeakConsumption(nodeID)
	if err != nil {
		log.Printf("Error reducing peak energy consumption for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Peak energy consumption reduced for node %s.\n", nodeID)
	return nil
}

// TrackCarbonOffsets logs carbon offset data for a node.
func TrackCarbonOffsets(nodeID string, offsetAmount float64) error {
	err := ledger.LogCarbonOffset(nodeID, offsetAmount)
	if err != nil {
		log.Printf("Error tracking carbon offset for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Carbon offset of %.2f logged for node %s.\n", offsetAmount, nodeID)
	return nil
}

// Helper function to format the energy usage report
func formatEnergyReport(data common.EnergyReportData) string {
	return "Energy Usage Report:\n" + data.String()
}
