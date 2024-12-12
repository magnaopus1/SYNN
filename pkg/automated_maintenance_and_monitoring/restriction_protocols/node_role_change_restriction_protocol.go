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
	NodeRoleCheckInterval      = 10 * time.Second // Interval for checking node role changes
	MaxAllowedRoleChanges      = 3               // Maximum number of role changes allowed per node
	RoleChangeCoolDownPeriod   = 24 * time.Hour  // Cool-down period for role changes
)

// NodeRoleChangeRestrictionAutomation monitors and restricts node role changes across the network
type NodeRoleChangeRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	nodeRoleChangeCount    map[string]int       // Tracks the number of role changes per node
	nodeLastRoleChangeTime map[string]time.Time // Tracks the last role change timestamp per node
}

// NewNodeRoleChangeRestrictionAutomation initializes and returns an instance of NodeRoleChangeRestrictionAutomation
func NewNodeRoleChangeRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *NodeRoleChangeRestrictionAutomation {
	return &NodeRoleChangeRestrictionAutomation{
		consensusSystem:        consensusSystem,
		ledgerInstance:         ledgerInstance,
		stateMutex:             stateMutex,
		nodeRoleChangeCount:    make(map[string]int),
		nodeLastRoleChangeTime: make(map[string]time.Time),
	}
}

// StartRoleChangeMonitoring starts continuous monitoring of node role changes
func (automation *NodeRoleChangeRestrictionAutomation) StartRoleChangeMonitoring() {
	ticker := time.NewTicker(NodeRoleCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorNodeRoleChanges()
		}
	}()
}

// monitorNodeRoleChanges checks for node role changes and enforces restrictions if necessary
func (automation *NodeRoleChangeRestrictionAutomation) monitorNodeRoleChanges() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch node role change data from Synnergy Consensus
	nodeData := automation.consensusSystem.GetNodeRoleChangeData()

	for nodeID, roleChangeDetails := range nodeData {
		lastChange := automation.nodeLastRoleChangeTime[nodeID]
		if time.Since(lastChange) < RoleChangeCoolDownPeriod {
			automation.flagRoleChangeViolation(nodeID, "Role change attempted during cooldown period")
			continue
		}

		// Check if the node has exceeded the allowed number of role changes
		if automation.nodeRoleChangeCount[nodeID] > MaxAllowedRoleChanges {
			automation.flagRoleChangeViolation(nodeID, "Exceeded allowed number of role changes")
		}
	}
}

// flagRoleChangeViolation flags a node's role change violation and logs it in the ledger
func (automation *NodeRoleChangeRestrictionAutomation) flagRoleChangeViolation(nodeID string, reason string) {
	fmt.Printf("Node role change violation: Node ID %s, Reason: %s\n", nodeID, reason)

	// Log the violation in the ledger
	automation.logRoleChangeViolation(nodeID, reason)
}

// logRoleChangeViolation logs the flagged node role change violation into the ledger with full details
func (automation *NodeRoleChangeRestrictionAutomation) logRoleChangeViolation(nodeID string, violationReason string) {
	// Create a ledger entry for node role change violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("node-role-change-violation-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Node Role Change Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Node %s violated role change rules. Reason: %s", nodeID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptRoleChangeData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log node role change violation:", err)
	} else {
		fmt.Println("Node role change violation logged.")
	}
}

// encryptRoleChangeData encrypts the role change data before logging for security
func (automation *NodeRoleChangeRestrictionAutomation) encryptRoleChangeData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting role change data:", err)
		return data
	}
	return string(encryptedData)
}
