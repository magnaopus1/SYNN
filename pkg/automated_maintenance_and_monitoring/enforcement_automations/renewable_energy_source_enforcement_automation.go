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

// Configuration for renewable energy source enforcement automation
const (
	EnergyCheckInterval               = 30 * time.Second // Interval to check energy compliance
	NonRenewableUsageLimit            = 10               // Max allowed percentage of non-renewable energy usage
	EnergyComplianceViolationThreshold = 3               // Allowed violations before restricting a node
)

// RenewableEnergySourceEnforcementAutomation monitors and enforces renewable energy compliance for nodes
type RenewableEnergySourceEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	energyComplianceMap  map[string]int // Tracks renewable compliance violations for each node
	nodeEnergySourceMap  map[string]int // Tracks percentage of renewable energy usage for each node
}

// NewRenewableEnergySourceEnforcementAutomation initializes the renewable energy source enforcement automation
func NewRenewableEnergySourceEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *RenewableEnergySourceEnforcementAutomation {
	return &RenewableEnergySourceEnforcementAutomation{
		networkManager:       networkManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		energyComplianceMap:  make(map[string]int),
		nodeEnergySourceMap:  make(map[string]int),
	}
}

// StartRenewableEnergyEnforcement begins continuous monitoring and enforcement of renewable energy compliance
func (automation *RenewableEnergySourceEnforcementAutomation) StartRenewableEnergyEnforcement() {
	ticker := time.NewTicker(EnergyCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkEnergySourceCompliance()
		}
	}()
}

// checkEnergySourceCompliance monitors each node's energy source usage and restricts non-compliant nodes
func (automation *RenewableEnergySourceEnforcementAutomation) checkEnergySourceCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.evaluateNodeEnergySources()
	automation.enforceEnergySourceRestrictions()
}

// evaluateNodeEnergySources checks each node’s renewable energy usage percentage and flags non-compliance
func (automation *RenewableEnergySourceEnforcementAutomation) evaluateNodeEnergySources() {
	for _, nodeID := range automation.networkManager.GetAllNodes() {
		renewableUsage := automation.networkManager.GetNodeRenewableEnergyPercentage(nodeID)
		automation.nodeEnergySourceMap[nodeID] = renewableUsage

		if renewableUsage < 100-NonRenewableUsageLimit {
			automation.energyComplianceMap[nodeID]++
		} else {
			automation.energyComplianceMap[nodeID] = 0 // Reset if compliant
		}
	}
}

// enforceEnergySourceRestrictions restricts nodes that repeatedly violate renewable energy compliance
func (automation *RenewableEnergySourceEnforcementAutomation) enforceEnergySourceRestrictions() {
	for nodeID, violations := range automation.energyComplianceMap {
		if violations >= EnergyComplianceViolationThreshold && automation.nodeEnergySourceMap[nodeID] < 100-NonRenewableUsageLimit {
			fmt.Printf("Energy source enforcement triggered for node %s due to renewable compliance violation.\n", nodeID)
			automation.applyRestriction(nodeID, "Insufficient Renewable Energy Usage")
		}
	}
}

// applyRestriction restricts non-compliant nodes that do not meet renewable energy requirements
func (automation *RenewableEnergySourceEnforcementAutomation) applyRestriction(nodeID, reason string) {
	err := automation.networkManager.RestrictNode(nodeID)
	if err != nil {
		fmt.Printf("Failed to restrict non-compliant node %s: %v\n", nodeID, err)
		automation.logEnergyComplianceAction(nodeID, "Restriction Failed", reason)
	} else {
		fmt.Printf("Node %s restricted due to %s.\n", nodeID, reason)
		automation.logEnergyComplianceAction(nodeID, "Restricted", reason)
	}
}

// logEnergyComplianceAction securely logs actions related to renewable energy enforcement
func (automation *RenewableEnergySourceEnforcementAutomation) logEnergyComplianceAction(nodeID, action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, Node ID: %s, Reason: %s", action, nodeID, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("renewable-energy-enforcement-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Renewable Energy Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log renewable energy enforcement action for node %s: %v\n", nodeID, err)
	} else {
		fmt.Println("Renewable energy enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *RenewableEnergySourceEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualEnergyComplianceCheck allows administrators to manually enforce energy compliance for a specific node
func (automation *RenewableEnergySourceEnforcementAutomation) TriggerManualEnergyComplianceCheck(nodeID string) {
	fmt.Printf("Manually enforcing renewable energy compliance for node: %s\n", nodeID)

	renewableUsage := automation.networkManager.GetNodeRenewableEnergyPercentage(nodeID)
	if renewableUsage < 100-NonRenewableUsageLimit {
		automation.applyRestriction(nodeID, "Manual Trigger: Insufficient Renewable Energy Usage")
	} else {
		fmt.Printf("Node %s meets renewable energy compliance.\n", nodeID)
		automation.logEnergyComplianceAction(nodeID, "Manual Compliance Check Passed", "Renewable Energy Verified")
	}
}
