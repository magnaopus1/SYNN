package sustainability

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// ExecuteEnergyConservationCheck performs a thorough analysis of energy usage across nodes, 
// identifies deviations from conservation goals, and logs results to the ledger.
func ExecuteEnergyConservationCheck(nodeID string) error {
	// Step 1: Retrieve energy usage data for the node
	energyUsage, err := ledger.GetNodeEnergyConsumption(nodeID)
	if err != nil {
		log.Printf("Error retrieving energy consumption data for node %s: %v\n", nodeID, err)
		return err
	}
	
	// Step 2: Retrieve predefined energy conservation targets from the ledger
	conservationTargets, err := ledger.GetEnergyConservationTargets(nodeID)
	if err != nil {
		log.Printf("Error retrieving conservation targets for node %s: %v\n", nodeID, err)
		return err
	}
	
	// Step 3: Compare actual usage with targets
	isCompliant, deviation := compareEnergyUsage(energyUsage, conservationTargets)
	if isCompliant {
		log.Printf("Node %s meets energy conservation targets.\n", nodeID)
	} else {
		log.Printf("Node %s deviates from energy targets by %d units.\n", nodeID, deviation)
		
		// Step 4: Encrypt deviation data and record in the ledger
		encryptedDeviation, encErr := encryption.EncryptInt(deviation)
		if encErr != nil {
			log.Printf("Error encrypting deviation data for node %s: %v\n", nodeID, encErr)
			return encErr
		}
		if err := ledger.LogEnergyDeviation(nodeID, encryptedDeviation); err != nil {
			log.Printf("Error logging energy deviation for node %s: %v\n", nodeID, err)
			return err
		}
	}
	
	// Step 5: Log conservation check completion in ledger
	if err := ledger.RecordEnergyConservationCheck(nodeID, isCompliant); err != nil {
		log.Printf("Error logging conservation check completion for node %s: %v\n", nodeID, err)
		return err
	}
	log.Printf("Energy conservation check executed successfully for node %s.\n", nodeID)
	return nil
}

// compareEnergyUsage is a helper function that checks if the energy usage meets conservation targets.
func compareEnergyUsage(actual common.EnergyUsage, targets common.EnergyTargets) (bool, int) {
	// Calculate deviation
	deviation := actual.Consumption - targets.MaxConsumption
	isCompliant := deviation <= 0
	return isCompliant, deviation
}
