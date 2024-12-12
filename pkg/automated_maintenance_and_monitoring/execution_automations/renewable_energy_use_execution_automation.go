package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/energy"
)

const (
	EnergyCheckInterval           = 15 * time.Minute // Interval for monitoring energy usage
	RenewableEnergyThreshold      = 0.95             // 95% of the energy must come from renewable sources
	EnergyUsageLedgerEntryType    = "Energy Use Event"
	EnergyComplianceLedgerEntryType = "Energy Compliance Event"
)

// RenewableEnergyUseExecutionAutomation ensures that all nodes comply with renewable energy usage requirements
type RenewableEnergyUseExecutionAutomation struct {
	consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validation
	ledgerInstance    *ledger.Ledger                        // Ledger instance for logging events
	energyManager     *energy.Manager                       // Energy management system
	executionMutex    *sync.RWMutex                         // Mutex for thread-safe execution
}

// NewRenewableEnergyUseExecutionAutomation initializes the automation for renewable energy monitoring
func NewRenewableEnergyUseExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, energyManager *energy.Manager, executionMutex *sync.RWMutex) *RenewableEnergyUseExecutionAutomation {
	return &RenewableEnergyUseExecutionAutomation{
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		energyManager:     energyManager,
		executionMutex:    executionMutex,
	}
}

// StartEnergyUseMonitor starts the continuous monitoring of energy usage by nodes in the network
func (automation *RenewableEnergyUseExecutionAutomation) StartEnergyUseMonitor() {
	ticker := time.NewTicker(EnergyCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkAndEnforceEnergyCompliance()
		}
	}()
}

// checkAndEnforceEnergyCompliance monitors energy usage and ensures compliance with renewable energy requirements
func (automation *RenewableEnergyUseExecutionAutomation) checkAndEnforceEnergyCompliance() {
	automation.executionMutex.Lock()
	defer automation.executionMutex.Unlock()

	// Fetch all nodes' energy usage data from the energy manager
	nodeEnergyData := automation.energyManager.GetNodeEnergyUsage()

	for _, node := range nodeEnergyData {
		automation.evaluateEnergyCompliance(node)
	}
}

// evaluateEnergyCompliance checks if a node's energy usage complies with the renewable energy threshold
func (automation *RenewableEnergyUseExecutionAutomation) evaluateEnergyCompliance(node *energy.NodeEnergyUsage) {
	renewablePercentage := node.GetRenewableEnergyPercentage()

	if renewablePercentage < RenewableEnergyThreshold {
		automation.flagNonCompliance(node)
	} else {
		automation.logEnergyCompliance(node)
	}
}

// flagNonCompliance triggers actions for nodes that are not compliant with renewable energy use
func (automation *RenewableEnergyUseExecutionAutomation) flagNonCompliance(node *energy.NodeEnergyUsage) {
	fmt.Printf("Node %s is non-compliant with renewable energy usage (%.2f%% renewable energy)\n", node.ID, node.GetRenewableEnergyPercentage()*100)

	// Penalize non-compliant nodes through the consensus engine
	err := automation.consensusEngine.PenalizeNodeForEnergyNonCompliance(node.ID)
	if err != nil {
		fmt.Printf("Failed to penalize node %s for energy non-compliance: %v\n", node.ID, err)
		return
	}

	// Log the non-compliance event into the ledger
	automation.logNonComplianceInLedger(node)
}

// logNonComplianceInLedger securely logs a non-compliance event in the ledger
func (automation *RenewableEnergyUseExecutionAutomation) logNonComplianceInLedger(node *energy.NodeEnergyUsage) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("energy-non-compliance-%s-%d", node.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      EnergyComplianceLedgerEntryType,
		Status:    "Non-Compliance",
		Details:   fmt.Sprintf("Node %s was non-compliant with renewable energy use: %.2f%% renewable energy.", node.ID, node.GetRenewableEnergyPercentage()*100),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log energy non-compliance for node %s: %v\n", node.ID, err)
	} else {
		fmt.Printf("Energy non-compliance logged successfully for node %s.\n", node.ID)
	}
}

// logEnergyCompliance securely logs an energy compliance event in the ledger
func (automation *RenewableEnergyUseExecutionAutomation) logEnergyCompliance(node *energy.NodeEnergyUsage) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("energy-compliance-%s-%d", node.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      EnergyUsageLedgerEntryType,
		Status:    "Compliant",
		Details:   fmt.Sprintf("Node %s is compliant with renewable energy use: %.2f%% renewable energy.", node.ID, node.GetRenewableEnergyPercentage()*100),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log energy compliance for node %s: %v\n", node.ID, err)
	} else {
		fmt.Printf("Energy compliance logged successfully for node %s.\n", node.ID)
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *RenewableEnergyUseExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualEnergyCheck allows administrators to manually trigger an energy compliance check for a specific node
func (automation *RenewableEnergyUseExecutionAutomation) TriggerManualEnergyCheck(nodeID string) {
	fmt.Printf("Manually triggering energy compliance check for node %s...\n", nodeID)

	node := automation.energyManager.GetNodeByID(nodeID)
	if node != nil {
		automation.evaluateEnergyCompliance(node)
	} else {
		fmt.Printf("Node %s not found.\n", nodeID)
	}
}
