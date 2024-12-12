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

// Configuration for network protocol enforcement automation
const (
	ProtocolCheckInterval         = 30 * time.Second // Interval to check network protocol compliance
	MaxProtocolViolationThreshold = 5                // Max allowed protocol violations before action
)

// NetworkProtocolEnforcementAutomation monitors and enforces compliance with network protocol standards
type NetworkProtocolEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	protocolViolationMap map[string]int // Tracks protocol violations for each node
	compliantNodeMap     map[string]bool // Tracks compliance status of nodes
}

// NewNetworkProtocolEnforcementAutomation initializes the network protocol enforcement automation
func NewNetworkProtocolEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *NetworkProtocolEnforcementAutomation {
	return &NetworkProtocolEnforcementAutomation{
		networkManager:       networkManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		protocolViolationMap: make(map[string]int),
		compliantNodeMap:     make(map[string]bool),
	}
}

// StartProtocolEnforcement begins continuous monitoring and enforcement of protocol compliance
func (automation *NetworkProtocolEnforcementAutomation) StartProtocolEnforcement() {
	ticker := time.NewTicker(ProtocolCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkProtocolCompliance()
		}
	}()
}

// checkProtocolCompliance monitors each node’s compliance with protocol standards and restricts non-compliant nodes
func (automation *NetworkProtocolEnforcementAutomation) checkProtocolCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.verifyProtocolAdherence()
	automation.enforceProtocolRestrictions()
}

// verifyProtocolAdherence checks if each node complies with network protocol standards
func (automation *NetworkProtocolEnforcementAutomation) verifyProtocolAdherence() {
	for _, nodeID := range automation.networkManager.GetAllNodes() {
		isCompliant := automation.networkManager.CheckNodeProtocolCompliance(nodeID)
		automation.compliantNodeMap[nodeID] = isCompliant

		if !isCompliant {
			automation.protocolViolationMap[nodeID]++
		} else {
			automation.protocolViolationMap[nodeID] = 0 // Reset if compliant
		}
	}
}

// enforceProtocolRestrictions restricts nodes with excessive protocol violations
func (automation *NetworkProtocolEnforcementAutomation) enforceProtocolRestrictions() {
	for nodeID, violations := range automation.protocolViolationMap {
		if violations > MaxProtocolViolationThreshold && !automation.compliantNodeMap[nodeID] {
			fmt.Printf("Protocol enforcement triggered for node %s due to repeated violations.\n", nodeID)
			automation.applyRestriction(nodeID, "Exceeded Protocol Violation Threshold")
		}
	}
}

// applyRestriction restricts access for nodes that repeatedly violate network protocol standards
func (automation *NetworkProtocolEnforcementAutomation) applyRestriction(nodeID, reason string) {
	err := automation.networkManager.RestrictNode(nodeID)
	if err != nil {
		fmt.Printf("Failed to restrict non-compliant node %s: %v\n", nodeID, err)
		automation.logProtocolAction(nodeID, "Restriction Failed", reason)
	} else {
		fmt.Printf("Node %s restricted due to %s.\n", nodeID, reason)
		automation.logProtocolAction(nodeID, "Restricted", reason)
	}
}

// logProtocolAction securely logs actions related to network protocol enforcement
func (automation *NetworkProtocolEnforcementAutomation) logProtocolAction(nodeID, action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, Node ID: %s, Reason: %s", action, nodeID, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("protocol-enforcement-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Network Protocol Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log protocol enforcement action for node %s: %v\n", nodeID, err)
	} else {
		fmt.Println("Protocol enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *NetworkProtocolEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualProtocolEnforcement allows administrators to manually restrict a non-compliant node for protocol violations
func (automation *NetworkProtocolEnforcementAutomation) TriggerManualProtocolEnforcement(nodeID string) {
	fmt.Printf("Manually enforcing protocol compliance for node: %s\n", nodeID)

	if automation.compliantNodeMap[nodeID] {
		fmt.Printf("Node %s is already protocol compliant.\n", nodeID)
		automation.logProtocolAction(nodeID, "Manual Enforcement Skipped - Already Compliant", "Manual Check")
	} else {
		automation.applyRestriction(nodeID, "Manual Trigger: Protocol Non-Compliance Restriction")
	}
}
