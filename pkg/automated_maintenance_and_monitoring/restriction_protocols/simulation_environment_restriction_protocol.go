package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

// Constants for simulation monitoring
const (
	SimulationCheckInterval   = 15 * time.Second // Interval to monitor simulation environments
	SimulationErrorThreshold  = 3                // Max allowed errors in simulation environments
	UnauthorizedAccessMessage = "Unauthorized access detected in simulation environment"
)

// SimulationEnvironmentRestrictionAutomation regulates simulation environments and prevents unauthorized access or misconfigurations.
type SimulationEnvironmentRestrictionAutomation struct {
	consensusSystem         *consensus.SynnergyConsensus
	ledgerInstance          *ledger.Ledger
	stateMutex              *sync.RWMutex
	simulationErrorCount    map[string]int // Tracks errors in simulation environments
	simulationStatus        map[string]bool // Status of each simulation environment (Active/Restricted)
}

// NewSimulationEnvironmentRestrictionAutomation initializes the automation
func NewSimulationEnvironmentRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SimulationEnvironmentRestrictionAutomation {
	return &SimulationEnvironmentRestrictionAutomation{
		consensusSystem:      consensusSystem,
		ledgerInstance:       ledgerInstance,
		stateMutex:           stateMutex,
		simulationErrorCount: make(map[string]int),
		simulationStatus:     make(map[string]bool),
	}
}

// StartSimulationEnvironmentMonitoring continuously monitors the simulation environments
func (automation *SimulationEnvironmentRestrictionAutomation) StartSimulationEnvironmentMonitoring() {
	ticker := time.NewTicker(SimulationCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateSimulationEnvironments()
		}
	}()
}

// evaluateSimulationEnvironments checks the status of all simulation environments and triggers restrictions as needed
func (automation *SimulationEnvironmentRestrictionAutomation) evaluateSimulationEnvironments() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	simulationEnvironments := automation.consensusSystem.GetSimulationEnvironments()

	for simulationID, isActive := range simulationEnvironments {
		if !isActive {
			automation.simulationErrorCount[simulationID]++
			automation.logSimulationError(simulationID)

			// Trigger restriction if error count exceeds the threshold
			if automation.simulationErrorCount[simulationID] >= SimulationErrorThreshold {
				automation.restrictSimulationEnvironment(simulationID)
			}
		} else {
			// Reset error count if simulation environment is functioning normally
			automation.simulationErrorCount[simulationID] = 0
		}
	}
}

// logSimulationError logs unauthorized access or error events into the ledger
func (automation *SimulationEnvironmentRestrictionAutomation) logSimulationError(simulationID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("simulation-error-%s-%d", simulationID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Simulation Environment Error",
		Status:    "Error",
		Details:   fmt.Sprintf("Issue detected in simulation environment %s.", simulationID),
	}

	// Encrypt the error details before logging into the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log simulation error:", err)
	} else {
		fmt.Println("Simulation error logged for environment:", simulationID)
	}
}

// restrictSimulationEnvironment restricts a simulation environment due to repeated errors or unauthorized access
func (automation *SimulationEnvironmentRestrictionAutomation) restrictSimulationEnvironment(simulationID string) {
	fmt.Printf("Simulation environment %s has exceeded the error threshold. Access restricted.\n", simulationID)

	// Log restriction in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("simulation-restriction-%s-%d", simulationID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Simulation Environment Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Simulation environment %s restricted due to exceeding error threshold.", simulationID),
	}

	// Encrypt restriction details before adding to the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log simulation restriction:", err)
	} else {
		fmt.Println("Simulation environment restricted:", simulationID)
	}

	// Inform the consensus system to restrict further activity in the simulation environment
	automation.consensusSystem.RestrictSimulationEnvironment(simulationID)
}

// encryptData encrypts sensitive data before storing it in the ledger
func (automation *SimulationEnvironmentRestrictionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting simulation error details:", err)
		return data
	}
	return string(encryptedData)
}

// SimulateUnauthorizedAccess simulates unauthorized access for testing purposes
func (automation *SimulationEnvironmentRestrictionAutomation) SimulateUnauthorizedAccess(simulationID string) {
	fmt.Println(UnauthorizedAccessMessage)

	// Log unauthorized access in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("unauthorized-access-%s-%d", simulationID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Unauthorized Access",
		Status:    "Critical",
		Details:   fmt.Sprintf("Unauthorized access detected in simulation environment %s.", simulationID),
	}

	// Encrypt the access details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log unauthorized access:", err)
	} else {
		fmt.Println("Unauthorized access logged for environment:", simulationID)
	}

	// Immediately restrict simulation environment due to security breach
	automation.restrictSimulationEnvironment(simulationID)
}
