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

// Configuration for operational continuity enforcement automation
const (
	ContinuityCheckInterval        = 20 * time.Second // Interval to check operational continuity
	MaxAllowedDowntimePerNode      = 2 * time.Minute  // Maximum allowed downtime per node before action
	MinOperationalNodeCount        = 5                // Minimum number of operational nodes for network stability
)

// OperationalContinuityEnforcementAutomation monitors and enforces operational continuity requirements
type OperationalContinuityEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	downtimeMap          map[string]time.Duration // Tracks downtime duration for each node
	operationalNodeCount int                      // Tracks number of currently operational nodes
}

// NewOperationalContinuityEnforcementAutomation initializes the operational continuity enforcement automation
func NewOperationalContinuityEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *OperationalContinuityEnforcementAutomation {
	return &OperationalContinuityEnforcementAutomation{
		networkManager:       networkManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		downtimeMap:          make(map[string]time.Duration),
		operationalNodeCount: 0,
	}
}

// StartOperationalContinuityEnforcement begins continuous monitoring and enforcement of network continuity
func (automation *OperationalContinuityEnforcementAutomation) StartOperationalContinuityEnforcement() {
	ticker := time.NewTicker(ContinuityCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkOperationalContinuity()
		}
	}()
}

// checkOperationalContinuity monitors the network’s operational status and enforces continuity requirements
func (automation *OperationalContinuityEnforcementAutomation) checkOperationalContinuity() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.evaluateNodeStatus()
	automation.ensureMinimumOperationalNodes()
}

// evaluateNodeStatus updates the downtime status of nodes and identifies those with excessive downtime
func (automation *OperationalContinuityEnforcementAutomation) evaluateNodeStatus() {
	operationalNodes := 0

	for _, nodeID := range automation.networkManager.GetAllNodes() {
		isOperational := automation.networkManager.IsNodeOperational(nodeID)

		if isOperational {
			automation.downtimeMap[nodeID] = 0
			operationalNodes++
		} else {
			automation.downtimeMap[nodeID] += ContinuityCheckInterval
			if automation.downtimeMap[nodeID] >= MaxAllowedDowntimePerNode {
				fmt.Printf("Operational continuity enforcement triggered for node %s due to excessive downtime.\n", nodeID)
				automation.applyRestriction(nodeID, "Excessive Downtime")
			}
		}
	}

	automation.operationalNodeCount = operationalNodes
}

// ensureMinimumOperationalNodes checks if the network meets minimum operational node requirements and takes action if necessary
func (automation *OperationalContinuityEnforcementAutomation) ensureMinimumOperationalNodes() {
	if automation.operationalNodeCount < MinOperationalNodeCount {
		fmt.Println("Operational continuity enforcement triggered due to insufficient operational nodes.")
		automation.activateStandbyNodes()
		automation.logContinuityAction("Activated Standby Nodes", "Network Stabilization")
	}
}

// applyRestriction restricts non-operational nodes that fail to meet continuity standards
func (automation *OperationalContinuityEnforcementAutomation) applyRestriction(nodeID, reason string) {
	err := automation.networkManager.RestrictNode(nodeID)
	if err != nil {
		fmt.Printf("Failed to restrict non-operational node %s: %v\n", nodeID, err)
		automation.logContinuityAction(fmt.Sprintf("Restriction Failed for Node %s", nodeID), reason)
	} else {
		fmt.Printf("Node %s restricted due to %s.\n", nodeID, reason)
		automation.logContinuityAction(fmt.Sprintf("Restricted Node %s", nodeID), reason)
	}
}

// activateStandbyNodes activates standby nodes to ensure network stability when operational nodes are insufficient
func (automation *OperationalContinuityEnforcementAutomation) activateStandbyNodes() {
	err := automation.networkManager.ActivateStandbyNodes(MinOperationalNodeCount - automation.operationalNodeCount)
	if err != nil {
		fmt.Println("Failed to activate standby nodes:", err)
		automation.logContinuityAction("Failed to Activate Standby Nodes", "Network Stabilization")
	} else {
		fmt.Println("Standby nodes activated to maintain network continuity.")
	}
}

// logContinuityAction securely logs actions related to operational continuity enforcement
func (automation *OperationalContinuityEnforcementAutomation) logContinuityAction(action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, Reason: %s", action, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("continuity-enforcement-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Operational Continuity Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log operational continuity enforcement action: %v\n", err)
	} else {
		fmt.Println("Operational continuity enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *OperationalContinuityEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualNodeActivation allows administrators to manually activate standby nodes for operational continuity
func (automation *OperationalContinuityEnforcementAutomation) TriggerManualNodeActivation() {
	fmt.Println("Manually triggering standby node activation for operational continuity.")

	automation.activateStandbyNodes()
	automation.logContinuityAction("Manual Standby Node Activation", "Administrator Action")
}
