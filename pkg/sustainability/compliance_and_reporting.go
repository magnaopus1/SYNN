package sustainability

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// ReviewEnergyConsumptionDeviations examines energy usage against expected benchmarks and flags deviations.
func ReviewEnergyConsumptionDeviations(nodeID string) error {
	energyUsage, err := ledger.GetNodeEnergyConsumption(nodeID)
	if err != nil {
		return err
	}
	if energyUsage.DeviatesFromTarget() {
		log.Printf("Energy consumption deviation found for node %s.\n", nodeID)
		return ledger.FlagEnergyDeviation(nodeID)
	}
	log.Printf("No significant deviations in energy consumption for node %s.\n", nodeID)
	return nil
}

// AllocateCarbonCredits allocates carbon credits based on eco-friendly contributions.
func AllocateCarbonCredits(nodeID string, credits int) error {
	encryptedCredits, err := encryption.EncryptInt(credits)
	if err != nil {
		return err
	}
	if err := ledger.AddCarbonCredits(nodeID, encryptedCredits); err != nil {
		return err
	}
	log.Printf("Allocated %d carbon credits to node %s.\n", credits, nodeID)
	return nil
}

// MonitorRenewableCompliance checks compliance with renewable energy targets.
func MonitorRenewableCompliance(nodeID string) (bool, error) {
	complianceStatus, err := ledger.CheckRenewableCompliance(nodeID)
	if err != nil {
		return false, err
	}
	log.Printf("Renewable compliance for node %s: %v\n", nodeID, complianceStatus)
	return complianceStatus, nil
}

// EnableAutoPowerAdjustment activates automatic adjustments for optimized power usage.
func EnableAutoPowerAdjustment(nodeID string) error {
	if err := ledger.SetAutoPowerAdjustment(nodeID, true); err != nil {
		return err
	}
	log.Printf("Automatic power adjustment enabled for node %s.\n", nodeID)
	return nil
}

// DisableAutoPowerAdjustment deactivates automatic power adjustments.
func DisableAutoPowerAdjustment(nodeID string) error {
	if err := ledger.SetAutoPowerAdjustment(nodeID, false); err != nil {
		return err
	}
	log.Printf("Automatic power adjustment disabled for node %s.\n", nodeID)
	return nil
}

// PromoteEnergyOptimization logs eco-friendly initiatives to optimize energy usage.
func PromoteEnergyOptimization(nodeID string, action common.EcoAction) error {
	encryptedAction, err := encryption.EncryptEcoAction(action)
	if err != nil {
		return err
	}
	if err := ledger.LogEcoFriendlyAction(nodeID, encryptedAction); err != nil {
		return err
	}
	log.Printf("Energy optimization action logged for node %s.\n", nodeID)
	return nil
}

// LogEcoFriendlyActions records initiatives aimed at reducing environmental impact.
func LogEcoFriendlyActions(nodeID string, actions []common.EcoAction) error {
	for _, action := range actions {
		if err := PromoteEnergyOptimization(nodeID, action); err != nil {
			return err
		}
	}
	log.Printf("Eco-friendly actions logged for node %s.\n", nodeID)
	return nil
}

// EstablishEcoPerformanceGoals sets eco-performance targets for nodes to achieve.
func EstablishEcoPerformanceGoals(nodeID string, goals common.EcoGoals) error {
	encryptedGoals, err := encryption.EncryptEcoGoals(goals)
	if err != nil {
		return err
	}
	if err := ledger.SetEcoPerformanceGoals(nodeID, encryptedGoals); err != nil {
		return err
	}
	log.Printf("Eco-performance goals established for node %s.\n", nodeID)
	return nil
}

// EvaluateEnergyConservationSuccess measures success against energy-saving targets.
func EvaluateEnergyConservationSuccess(nodeID string) (bool, error) {
	success, err := ledger.CheckEnergyConservation(nodeID)
	if err != nil {
		return false, err
	}
	log.Printf("Energy conservation success for node %s: %v\n", nodeID, success)
	return success, nil
}

// TrackGreenPowerConsumption monitors green energy usage and records it in the ledger.
func TrackGreenPowerConsumption(nodeID string, usage int) error {
	encryptedUsage, err := encryption.EncryptInt(usage)
	if err != nil {
		return err
	}
	if err := ledger.RecordGreenPowerUsage(nodeID, encryptedUsage); err != nil {
		return err
	}
	log.Printf("Green power consumption recorded for node %s: %d units.\n", nodeID, usage)
	return nil
}

// SetEnergyEfficiencyTargets assigns energy efficiency targets for nodes to aim for.
func SetEnergyEfficiencyTargets(nodeID string, targets common.EfficiencyTargets) error {
	encryptedTargets, err := encryption.EncryptEfficiencyTargets(targets)
	if err != nil {
		return err
	}
	if err := ledger.SetEfficiencyTargets(nodeID, encryptedTargets); err != nil {
		return err
	}
	log.Printf("Energy efficiency targets set for node %s.\n", nodeID)
	return nil
}

// ValidateGreenComplianceStatus confirms compliance with green standards.
func ValidateGreenComplianceStatus(nodeID string) (bool, error) {
	isCompliant, err := ledger.VerifyGreenCompliance(nodeID)
	if err != nil {
		return false, err
	}
	log.Printf("Green compliance status for node %s: %v\n", nodeID, isCompliant)
	return isCompliant, nil
}

// RecordRenewableEnergyAllocation documents renewable energy allocation in the ledger.
func RecordRenewableEnergyAllocation(nodeID string, allocation int) error {
	encryptedAllocation, err := encryption.EncryptInt(allocation)
	if err != nil {
		return err
	}
	if err := ledger.RecordRenewableAllocation(nodeID, encryptedAllocation); err != nil {
		return err
	}
	log.Printf("Renewable energy allocation recorded for node %s: %d units.\n", nodeID, allocation)
	return nil
}

// ReportPowerUsageReduction logs reductions in power usage for audit and compliance purposes.
func ReportPowerUsageReduction(nodeID string, reduction int) error {
	encryptedReduction, err := encryption.EncryptInt(reduction)
	if err != nil {
		return err
	}
	if err := ledger.RecordPowerReduction(nodeID, encryptedReduction); err != nil {
		return err
	}
	log.Printf("Power usage reduction reported for node %s: %d units.\n", nodeID, reduction)
	return nil
}

// LogRenewableIncentiveDistribution records incentives for using renewable energy sources.
func LogRenewableIncentiveDistribution(nodeID string, incentiveAmount int) error {
	encryptedIncentive, err := encryption.EncryptInt(incentiveAmount)
	if err != nil {
		return err
	}
	if err := ledger.RecordRenewableIncentive(nodeID, encryptedIncentive); err != nil {
		return err
	}
	log.Printf("Renewable incentive distributed to node %s: %d credits.\n", nodeID, incentiveAmount)
	return nil
}

// ReviewNodeEcoFriendliness evaluates and logs the eco-friendliness of nodes based on metrics.
func ReviewNodeEcoFriendliness(nodeID string) (common.EcoScore, error) {
	ecoScore, err := ledger.GetEcoScore(nodeID)
	if err != nil {
		return common.EcoScore{}, err
	}
	log.Printf("Eco-friendliness review for node %s: %v\n", nodeID, ecoScore)
	return ecoScore, nil
}
