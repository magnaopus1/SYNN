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

const (
	EnergyCheckInterval         = 20 * time.Second // Interval for checking renewable energy usage
	MaxAllowedNonRenewableUsage = 15               // Maximum allowed percentage of non-renewable energy usage per node
)

// RenewableEnergyRestrictionAutomation monitors and enforces the use of renewable energy sources across the network
type RenewableEnergyRestrictionAutomation struct {
	consensusSystem         *consensus.SynnergyConsensus
	ledgerInstance          *ledger.Ledger
	stateMutex              *sync.RWMutex
	nonRenewableUsageCount  map[string]float64 // Tracks non-renewable energy usage percentage per node
}

// NewRenewableEnergyRestrictionAutomation initializes and returns an instance of RenewableEnergyRestrictionAutomation
func NewRenewableEnergyRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RenewableEnergyRestrictionAutomation {
	return &RenewableEnergyRestrictionAutomation{
		consensusSystem:        consensusSystem,
		ledgerInstance:         ledgerInstance,
		stateMutex:             stateMutex,
		nonRenewableUsageCount: make(map[string]float64),
	}
}

// StartEnergyMonitoring starts continuous monitoring of renewable energy usage compliance
func (automation *RenewableEnergyRestrictionAutomation) StartEnergyMonitoring() {
	ticker := time.NewTicker(EnergyCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorEnergyUsage()
		}
	}()
}

// monitorEnergyUsage checks the network's energy usage and enforces restrictions if necessary
func (automation *RenewableEnergyRestrictionAutomation) monitorEnergyUsage() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch energy usage data from Synnergy Consensus
	energyData := automation.consensusSystem.GetNodeEnergyUsage()

	for nodeID, nonRenewableUsage := range energyData {
		// Check if the node exceeds the allowed non-renewable energy usage percentage
		if nonRenewableUsage > MaxAllowedNonRenewableUsage {
			automation.flagEnergyViolation(nodeID, nonRenewableUsage, "Exceeded allowed non-renewable energy usage")
		}
	}
}

// flagEnergyViolation flags a node's non-compliant energy usage and logs it in the ledger
func (automation *RenewableEnergyRestrictionAutomation) flagEnergyViolation(nodeID string, nonRenewableUsage float64, reason string) {
	fmt.Printf("Energy usage violation: Node ID %s, Reason: %s\n", nodeID, reason)

	// Log the violation in the ledger
	automation.logEnergyViolation(nodeID, nonRenewableUsage, reason)
}

// logEnergyViolation logs the flagged energy usage violation into the ledger with full details
func (automation *RenewableEnergyRestrictionAutomation) logEnergyViolation(nodeID string, nonRenewableUsage float64, violationReason string) {
	// Create a ledger entry for energy usage violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("energy-violation-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Energy Usage Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Node %s violated renewable energy usage rules. Non-Renewable Usage: %.2f%%. Reason: %s", nodeID, nonRenewableUsage, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptEnergyData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log energy usage violation:", err)
	} else {
		fmt.Println("Energy usage violation logged.")
	}
}

// encryptEnergyData encrypts the energy usage data before logging for security
func (automation *RenewableEnergyRestrictionAutomation) encryptEnergyData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting energy data:", err)
		return data
	}
	return string(encryptedData)
}
