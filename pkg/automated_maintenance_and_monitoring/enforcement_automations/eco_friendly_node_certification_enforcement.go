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

// Configuration for eco-friendly node certification enforcement
const (
	EcoCheckInterval               = 30 * time.Second // Interval to check eco-friendly compliance
	MinRenewableEnergyPercentage   = 70               // Minimum renewable energy usage percentage for certification
	MaxEnergyUsageThreshold        = 500              // Maximum allowable energy usage (in kWh) per hour per node
)

// EcoFriendlyNodeCertificationEnforcementAutomation monitors and enforces eco-friendly node compliance
type EcoFriendlyNodeCertificationEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	nodeEnergyUsageMap   map[string]float64 // Tracks energy usage per node
	nodeRenewableEnergy  map[string]int     // Tracks renewable energy usage percentage per node
}

// NewEcoFriendlyNodeCertificationEnforcementAutomation initializes eco-friendly node certification enforcement
func NewEcoFriendlyNodeCertificationEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *EcoFriendlyNodeCertificationEnforcementAutomation {
	return &EcoFriendlyNodeCertificationEnforcementAutomation{
		networkManager:      networkManager,
		consensusEngine:     consensusEngine,
		ledgerInstance:      ledgerInstance,
		enforcementMutex:    enforcementMutex,
		nodeEnergyUsageMap:  make(map[string]float64),
		nodeRenewableEnergy: make(map[string]int),
	}
}

// StartEcoCertificationEnforcement begins continuous monitoring and enforcement of eco-friendly node practices
func (automation *EcoFriendlyNodeCertificationEnforcementAutomation) StartEcoCertificationEnforcement() {
	ticker := time.NewTicker(EcoCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkEcoFriendlyCompliance()
		}
	}()
}

// checkEcoFriendlyCompliance monitors each node's energy usage and renewable energy compliance
func (automation *EcoFriendlyNodeCertificationEnforcementAutomation) checkEcoFriendlyCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.verifyEnergyUsage()
	automation.verifyRenewableEnergyUsage()
}

// verifyEnergyUsage checks if each node’s energy usage is within allowable thresholds
func (automation *EcoFriendlyNodeCertificationEnforcementAutomation) verifyEnergyUsage() {
	for _, nodeID := range automation.networkManager.GetAllNodes() {
		energyUsage := automation.networkManager.GetNodeEnergyUsage(nodeID)

		if energyUsage > MaxEnergyUsageThreshold {
			fmt.Printf("Energy usage violation for node %s.\n", nodeID)
			automation.applyCertificationRestriction(nodeID, "Excessive Energy Usage")
		}
	}
}

// verifyRenewableEnergyUsage checks if each node meets the renewable energy usage requirement
func (automation *EcoFriendlyNodeCertificationEnforcementAutomation) verifyRenewableEnergyUsage() {
	for _, nodeID := range automation.networkManager.GetAllNodes() {
		renewablePercentage := automation.networkManager.GetNodeRenewableEnergyPercentage(nodeID)

		if renewablePercentage < MinRenewableEnergyPercentage {
			fmt.Printf("Renewable energy usage violation for node %s.\n", nodeID)
			automation.applyCertificationRestriction(nodeID, "Insufficient Renewable Energy Usage")
		}
	}
}

// applyCertificationRestriction restricts nodes that fail to meet eco-friendly compliance requirements
func (automation *EcoFriendlyNodeCertificationEnforcementAutomation) applyCertificationRestriction(nodeID, reason string) {
	err := automation.networkManager.RestrictNode(nodeID)
	if err != nil {
		fmt.Printf("Failed to restrict node %s due to eco-friendly certification violation: %v\n", nodeID, err)
		automation.logEcoAction(nodeID, "Certification Restriction Failed", reason)
	} else {
		fmt.Printf("Node %s restricted due to %s.\n", nodeID, reason)
		automation.logEcoAction(nodeID, "Certification Restricted", reason)
	}
}

// logEcoAction securely logs actions related to eco-friendly node certification enforcement
func (automation *EcoFriendlyNodeCertificationEnforcementAutomation) logEcoAction(nodeID, action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, Node ID: %s, Reason: %s", action, nodeID, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("eco-certification-enforcement-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Eco-Friendly Certification Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log eco-friendly certification enforcement action for node %s: %v\n", nodeID, err)
	} else {
		fmt.Println("Eco-friendly certification enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *EcoFriendlyNodeCertificationEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualEcoCheck allows administrators to manually check eco-friendly compliance for a specific node
func (automation *EcoFriendlyNodeCertificationEnforcementAutomation) TriggerManualEcoCheck(nodeID string) {
	fmt.Printf("Manually triggering eco-friendly compliance check for node: %s\n", nodeID)

	energyUsage := automation.networkManager.GetNodeEnergyUsage(nodeID)
	renewablePercentage := automation.networkManager.GetNodeRenewableEnergyPercentage(nodeID)

	if energyUsage > MaxEnergyUsageThreshold {
		automation.applyCertificationRestriction(nodeID, "Manual Trigger: Excessive Energy Usage")
	} else if renewablePercentage < MinRenewableEnergyPercentage {
		automation.applyCertificationRestriction(nodeID, "Manual Trigger: Insufficient Renewable Energy Usage")
	} else {
		fmt.Printf("Node %s is compliant with eco-friendly certification requirements.\n", nodeID)
		automation.logEcoAction(nodeID, "Manual Compliance Check Passed", "Eco-Friendly Compliance Verified")
	}
}
