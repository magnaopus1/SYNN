package sustainability

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// AllocateResourcesBasedOnEnergyEfficiency allocates resources to nodes based on their energy efficiency.
func AllocateResourcesBasedOnEnergyEfficiency(nodeID string, efficiencyLevel float64) error {
	allocated, err := ledger.AllocateNodeResources(nodeID, efficiencyLevel)
	if err != nil {
		log.Printf("Error allocating resources based on energy efficiency for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Resources allocated to node %s with efficiency level %.2f.\n", nodeID, efficiencyLevel)
	return allocated
}

// LogEnergyConservationStatus records the energy conservation status for reporting.
func LogEnergyConservationStatus(nodeID string, status string) error {
	encryptedStatus, err := encryption.EncryptString(status)
	if err != nil {
		log.Printf("Error encrypting energy conservation status: %v\n", err)
		return err
	}
	err = ledger.RecordEnergyConservation(nodeID, encryptedStatus)
	if err != nil {
		log.Printf("Error logging energy conservation status for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Energy conservation status logged for node %s.\n", nodeID)
	return nil
}

// MeasureLongTermEnergyImpact calculates the long-term impact of energy usage for analysis.
func MeasureLongTermEnergyImpact(nodeID string) (float64, error) {
	impact, err := ledger.CalculateEnergyImpact(nodeID)
	if err != nil {
		log.Printf("Error measuring long-term energy impact for node %s: %v\n", nodeID, err)
		return 0, err
	}
	log.Printf("Long-term energy impact for node %s: %.2f.\n", nodeID, impact)
	return impact, nil
}

// AuditNodePowerUsage performs a detailed audit on the power usage of a node.
func AuditNodePowerUsage(nodeID string) error {
	report, err := ledger.GeneratePowerAudit(nodeID)
	if err != nil {
		log.Printf("Error auditing power usage for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Power usage audit completed for node %s: %s\n", nodeID, report)
	return nil
}

// ValidateResourceEfficiencyLevels checks if a node meets minimum resource efficiency standards.
func ValidateResourceEfficiencyLevels(nodeID string) (bool, error) {
	isEfficient, err := ledger.CheckEfficiencyCompliance(nodeID)
	if err != nil {
		log.Printf("Error validating resource efficiency for node %s: %v\n", nodeID, err)
		return false, err
	}
	log.Printf("Resource efficiency validation for node %s: %t.\n", nodeID, isEfficient)
	return isEfficient, nil
}

// ScheduleEnergyOptimizationTasks plans tasks to optimize energy usage for a node.
func ScheduleEnergyOptimizationTasks(nodeID string, taskDetails string) error {
	encryptedTask, err := encryption.EncryptString(taskDetails)
	if err != nil {
		log.Printf("Error encrypting task details: %v\n", err)
		return err
	}
	err = ledger.ScheduleOptimization(nodeID, encryptedTask)
	if err != nil {
		log.Printf("Error scheduling optimization tasks for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Energy optimization tasks scheduled for node %s.\n", nodeID)
	return nil
}

// ReduceDataIntensity adjusts data processing to reduce energy demand.
func ReduceDataIntensity(nodeID string) error {
	err := ledger.ReduceNodeDataIntensity(nodeID)
	if err != nil {
		log.Printf("Error reducing data intensity for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Data intensity reduced for node %s.\n", nodeID)
	return nil
}

// EnforceMinRenewableUsage enforces a minimum renewable energy usage percentage.
func EnforceMinRenewableUsage(nodeID string, minUsage float64) error {
	err := ledger.SetMinRenewableUsage(nodeID, minUsage)
	if err != nil {
		log.Printf("Error enforcing min renewable usage for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Minimum renewable usage of %.2f%% enforced for node %s.\n", minUsage, nodeID)
	return nil
}

// EnableGreenPowerPriority gives priority to green power resources for a node.
func EnableGreenPowerPriority(nodeID string) error {
	err := ledger.SetGreenPowerPriority(nodeID, true)
	if err != nil {
		log.Printf("Error enabling green power priority for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Green power priority enabled for node %s.\n", nodeID)
	return nil
}

// DisableGreenPowerPriority removes green power priority from a node.
func DisableGreenPowerPriority(nodeID string) error {
	err := ledger.SetGreenPowerPriority(nodeID, false)
	if err != nil {
		log.Printf("Error disabling green power priority for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Green power priority disabled for node %s.\n", nodeID)
	return nil
}

// GeneratePowerUsageForecast provides a forecast of power usage for planning purposes.
func GeneratePowerUsageForecast(nodeID string) (string, error) {
	forecast, err := ledger.CreatePowerForecast(nodeID)
	if err != nil {
		log.Printf("Error generating power usage forecast for node %s: %v\n", nodeID, err)
		return "", err
	}
	log.Printf("Power usage forecast generated for node %s.\n", nodeID)
	return forecast, nil
}

// AnalyzeEnergyOverusePatterns detects and analyzes patterns of excessive energy use.
func AnalyzeEnergyOverusePatterns(nodeID string) (string, error) {
	patternReport, err := ledger.DetectOverusePatterns(nodeID)
	if err != nil {
		log.Printf("Error analyzing energy overuse patterns for node %s: %v\n", nodeID, err)
		return "", err
	}
	log.Printf("Energy overuse patterns analyzed for node %s: %s\n", nodeID, patternReport)
	return patternReport, nil
}

// IssueEnergyEfficiencyCertificate grants an efficiency certificate to a compliant node.
func IssueEnergyEfficiencyCertificate(nodeID string) error {
	err := ledger.GrantEfficiencyCertificate(nodeID)
	if err != nil {
		log.Printf("Error issuing energy efficiency certificate for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Energy efficiency certificate issued for node %s.\n", nodeID)
	return nil
}

// EstablishEnergyConservationMode activates energy conservation mode on a node.
func EstablishEnergyConservationMode(nodeID string) error {
	err := ledger.EnableConservationMode(nodeID)
	if err != nil {
		log.Printf("Error establishing energy conservation mode for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Energy conservation mode established for node %s.\n", nodeID)
	return nil
}

// TrackAnnualEnergyUsage logs the yearly energy usage of a node.
func TrackAnnualEnergyUsage(nodeID string) error {
	yearlyUsage, err := ledger.GetYearlyEnergyUsage(nodeID)
	if err != nil {
		log.Printf("Error tracking annual energy usage for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Annual energy usage for node %s: %.2f units.\n", nodeID, yearlyUsage)
	return nil
}
