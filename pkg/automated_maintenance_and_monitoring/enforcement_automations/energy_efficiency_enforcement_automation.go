package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/network"
)

// Configuration for energy efficiency enforcement automation
const (
	EnergyCheckInterval           = 20 * time.Second // Interval to check energy efficiency compliance
	MaxAllowableEnergyConsumption = 300              // Maximum allowable energy consumption (in kWh) per node per hour
	MinEfficiencyRating           = 80               // Minimum energy efficiency rating for node compliance
)

// EnergyEfficiencyEnforcementAutomation monitors and enforces energy efficiency standards across the network
type EnergyEfficiencyEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	nodeEnergyUsageMap   map[string]float64 // Tracks energy usage per node
	nodeEfficiencyRating map[string]int     // Tracks energy efficiency rating per node
}

// NewEnergyEfficiencyEnforcementAutomation initializes the energy efficiency enforcement automation
func NewEnergyEfficiencyEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *EnergyEfficiencyEnforcementAutomation {
	return &EnergyEfficiencyEnforcementAutomation{
		networkManager:      networkManager,
		consensusEngine:     consensusEngine,
		ledgerInstance:      ledgerInstance,
		enforcementMutex:    enforcementMutex,
		nodeEnergyUsageMap:  make(map[string]float64),
		nodeEfficiencyRating: make(map[string]int),
	}
}

// StartEnergyEfficiencyEnforcement begins continuous monitoring and enforcement of energy efficiency compliance
func (automation *EnergyEfficiencyEnforcementAutomation) StartEnergyEfficiencyEnforcement() {
	ticker := time.NewTicker(EnergyCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkEnergyEfficiencyCompliance()
		}
	}()
}

// checkEnergyEfficiencyCompliance monitors each node’s energy usage and efficiency rating for compliance
func (automation *EnergyEfficiencyEnforcementAutomation) checkEnergyEfficiencyCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.verifyEnergyConsumption()
	automation.verifyEfficiencyRating()
}

// verifyEnergyConsumption checks if each node’s energy consumption is within allowable limits
func (automation *EnergyEfficiencyEnforcementAutomation) verifyEnergyConsumption() {
	for _, nodeID := range automation.networkManager.GetAllNodes() {
		energyUsage := automation.networkManager.GetNodeEnergyConsumption(nodeID)

		if energyUsage > MaxAllowableEnergyConsumption {
			fmt.Printf("Energy consumption violation for node %s.\n", nodeID)
			automation.applyEfficiencyRestriction(nodeID, "Excessive Energy Consumption")
		}
	}
}

// verifyEfficiencyRating checks if each node meets the minimum energy efficiency rating requirement
func (automation *EnergyEfficiencyEnforcementAutomation) verifyEfficiencyRating() {
	for _, nodeID := range automation.networkManager.GetAllNodes() {
		efficiencyRating := automation.networkManager.GetNodeEfficiencyRating(nodeID)

		if efficiencyRating < MinEfficiencyRating {
			fmt.Printf("Energy efficiency rating violation for node %s.\n", nodeID)
			automation.applyEfficiencyRestriction(nodeID, "Insufficient Efficiency Rating")
		}
	}
}

// applyEfficiencyRestriction restricts nodes that fail to meet energy efficiency standards
func (automation *EnergyEfficiencyEnforcementAutomation) applyEfficiencyRestriction(nodeID, reason string) {
	err := automation.networkManager.RestrictNode(nodeID)
	if err != nil {
		fmt.Printf("Failed to restrict node %s due to energy efficiency violation: %v\n", nodeID, err)
		automation.logEnergyAction(nodeID, "Efficiency Restriction Failed", reason)
	} else {
		fmt.Printf("Node %s restricted due to %s.\n", nodeID, reason)
		automation.logEnergyAction(nodeID, "Efficiency Restricted", reason)
	}
}

// logEnergyAction securely logs actions related to energy efficiency enforcement
func (automation *EnergyEfficiencyEnforcementAutomation) logEnergyAction(nodeID, action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, Node ID: %s, Reason: %s", action, nodeID, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("energy-efficiency-enforcement-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Energy Efficiency Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log energy efficiency enforcement action for node %s: %v\n", nodeID, err)
	} else {
		fmt.Println("Energy efficiency enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *EnergyEfficiencyEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualEfficiencyCheck allows administrators to manually check energy efficiency compliance for a specific node
func (automation *EnergyEfficiencyEnforcementAutomation) TriggerManualEfficiencyCheck(nodeID string) {
	fmt.Printf("Manually triggering energy efficiency compliance check for node: %s\n", nodeID)

	energyUsage := automation.networkManager.GetNodeEnergyConsumption(nodeID)
	efficiencyRating := automation.networkManager.GetNodeEfficiencyRating(nodeID)

	if energyUsage > MaxAllowableEnergyConsumption {
		automation.applyEfficiencyRestriction(nodeID, "Manual Trigger: Excessive Energy Consumption")
	} else if efficiencyRating < MinEfficiencyRating {
		automation.applyEfficiencyRestriction(nodeID, "Manual Trigger: Insufficient Efficiency Rating")
	} else {
		fmt.Printf("Node %s is compliant with energy efficiency standards.\n", nodeID)
		automation.logEnergyAction(nodeID, "Manual Compliance Check Passed", "Energy Efficiency Compliance Verified")
	}
}
