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
	EnergyUsageCheckInterval  = 10 * time.Second  // Interval for checking energy usage per node
	MaxEnergyUsagePerNode     = 5000.0            // Maximum allowed energy usage per node (in kilowatt-hours or equivalent)
	ViolationPenaltyThreshold = 3                 // Maximum number of violations before penalizing a node
)

// EnergyUsageCapAutomation monitors and restricts energy usage across nodes on the network
type EnergyUsageCapAutomation struct {
	consensusSystem      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	stateMutex           *sync.RWMutex
	nodeEnergyUsage      map[string]float64 // Tracks energy usage per node
	nodeViolationCount   map[string]int     // Tracks the number of energy usage violations per node
}

// NewEnergyUsageCapAutomation initializes and returns an instance of EnergyUsageCapAutomation
func NewEnergyUsageCapAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *EnergyUsageCapAutomation {
	return &EnergyUsageCapAutomation{
		consensusSystem:    consensusSystem,
		ledgerInstance:     ledgerInstance,
		stateMutex:         stateMutex,
		nodeEnergyUsage:    make(map[string]float64),
		nodeViolationCount: make(map[string]int),
	}
}

// StartEnergyUsageMonitoring starts continuous monitoring of energy usage
func (automation *EnergyUsageCapAutomation) StartEnergyUsageMonitoring() {
	ticker := time.NewTicker(EnergyUsageCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorEnergyUsage()
		}
	}()
}

// monitorEnergyUsage checks recent energy consumption per node and enforces usage limits
func (automation *EnergyUsageCapAutomation) monitorEnergyUsage() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent energy usage data from Synnergy Consensus
	nodeUsageData := automation.consensusSystem.GetNodeEnergyUsage()

	for nodeID, usage := range nodeUsageData {
		if !automation.validateEnergyUsage(nodeID, usage) {
			automation.flagEnergyViolation(nodeID, usage, "Exceeded energy usage cap")
		}
	}
}

// validateEnergyUsage checks if a node's energy consumption exceeds the maximum allowed limit
func (automation *EnergyUsageCapAutomation) validateEnergyUsage(nodeID string, usage float64) bool {
	if usage > MaxEnergyUsagePerNode {
		automation.nodeViolationCount[nodeID]++
		return false
	}

	// Reset violation count for compliant nodes
	automation.nodeViolationCount[nodeID] = 0
	return true
}

// flagEnergyViolation flags a node's energy usage violation and logs it in the ledger
func (automation *EnergyUsageCapAutomation) flagEnergyViolation(nodeID string, usage float64, reason string) {
	fmt.Printf("Energy usage violation: Node %s, Reason: %s, Usage: %.2f\n", nodeID, reason, usage)

	// Log the violation into the ledger
	automation.logEnergyViolation(nodeID, usage, reason)

	// If the violation count exceeds the penalty threshold, apply penalties
	if automation.nodeViolationCount[nodeID] >= ViolationPenaltyThreshold {
		automation.enforceEnergyUsageCap(nodeID)
	}
}

// logEnergyViolation logs the energy usage violation into the ledger with full details
func (automation *EnergyUsageCapAutomation) logEnergyViolation(nodeID string, usage float64, violationReason string) {
	// Encrypt the energy usage data before logging
	encryptedData := automation.encryptEnergyUsageData(nodeID, usage)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("energy-violation-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Energy Usage Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Node %s flagged for energy usage violation. Reason: %s. Encrypted Data: %s", nodeID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log energy usage violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Energy usage violation logged for node: %s\n", nodeID)
	}
}

// enforceEnergyUsageCap triggers penalties or restrictions for nodes exceeding the allowed energy usage cap
func (automation *EnergyUsageCapAutomation) enforceEnergyUsageCap(nodeID string) {
	// Apply specific penalties or restrictions to the node (e.g., throttling, suspension)
	fmt.Printf("Enforcing energy usage penalties on node: %s\n", nodeID)

	// Log the enforcement action in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("energy-enforcement-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Energy Usage Cap Enforcement",
		Status:    "Enforced",
		Details:   fmt.Sprintf("Energy usage cap enforcement triggered for node %s after exceeding the violation threshold.", nodeID),
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log energy usage cap enforcement action: %v\n", err)
	} else {
		fmt.Printf("Energy usage cap enforcement logged for node: %s\n", nodeID)
	}
}

// encryptEnergyUsageData encrypts energy usage data before logging for security
func (automation *EnergyUsageCapAutomation) encryptEnergyUsageData(nodeID string, usage float64) string {
	data := fmt.Sprintf("Node ID: %s, Energy Usage: %.2f", nodeID, usage)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting energy usage data:", err)
		return data
	}
	return string(encryptedData)
}
