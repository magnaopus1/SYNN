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
	"synnergy_network_demo/security"
)

// Configuration for slashing penalty enforcement automation
const (
	SlashingCheckInterval           = 15 * time.Second // Interval to check for slashing triggers
	MaxViolationThreshold           = 3                // Maximum allowed violations before slashing penalty
	SlashingPenaltyAmount           = 100              // Amount of tokens deducted as penalty
)

// SlashingPenaltyEnforcementAutomation monitors and enforces slashing penalties for nodes violating network rules
type SlashingPenaltyEnforcementAutomation struct {
	securityManager     *security.SecurityManager
	consensusEngine     *consensus.SynnergyConsensus
	ledgerInstance      *ledger.Ledger
	enforcementMutex    *sync.RWMutex
	nodeViolationMap    map[string]int // Tracks violation count for each node
}

// NewSlashingPenaltyEnforcementAutomation initializes the slashing penalty enforcement automation
func NewSlashingPenaltyEnforcementAutomation(securityManager *security.SecurityManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *SlashingPenaltyEnforcementAutomation {
	return &SlashingPenaltyEnforcementAutomation{
		securityManager:     securityManager,
		consensusEngine:     consensusEngine,
		ledgerInstance:      ledgerInstance,
		enforcementMutex:    enforcementMutex,
		nodeViolationMap:    make(map[string]int),
	}
}

// StartSlashingPenaltyEnforcement begins continuous monitoring and enforcement of slashing penalties
func (automation *SlashingPenaltyEnforcementAutomation) StartSlashingPenaltyEnforcement() {
	ticker := time.NewTicker(SlashingCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkForSlashingPenalties()
		}
	}()
}

// checkForSlashingPenalties monitors each node for rule violations and applies slashing penalties where necessary
func (automation *SlashingPenaltyEnforcementAutomation) checkForSlashingPenalties() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.evaluateNodeViolations()
	automation.applySlashingPenalties()
}

// evaluateNodeViolations checks each node for rule violations and flags nodes that reach the violation threshold
func (automation *SlashingPenaltyEnforcementAutomation) evaluateNodeViolations() {
	for _, nodeID := range automation.securityManager.GetMonitoredNodes() {
		violations := automation.securityManager.GetNodeViolations(nodeID)
		automation.nodeViolationMap[nodeID] = violations

		if violations >= MaxViolationThreshold {
			fmt.Printf("Node %s has reached maximum violation threshold.\n", nodeID)
		}
	}
}

// applySlashingPenalties enforces slashing penalties on nodes that exceed the violation threshold
func (automation *SlashingPenaltyEnforcementAutomation) applySlashingPenalties() {
	for nodeID, violations := range automation.nodeViolationMap {
		if violations >= MaxViolationThreshold {
			fmt.Printf("Applying slashing penalty to node %s for exceeding violation threshold.\n", nodeID)
			automation.slashNode(nodeID)
		}
	}
}

// slashNode deducts tokens as a penalty from nodes that breach the violation limit
func (automation *SlashingPenaltyEnforcementAutomation) slashNode(nodeID string) {
	err := automation.ledgerInstance.DeductTokens(nodeID, SlashingPenaltyAmount)
	if err != nil {
		fmt.Printf("Failed to apply slashing penalty to node %s: %v\n", nodeID, err)
		automation.logSlashingAction(nodeID, "Penalty Application Failed", fmt.Sprintf("Violation Count: %d", automation.nodeViolationMap[nodeID]))
	} else {
		fmt.Printf("Slashing penalty of %d tokens applied to node %s.\n", SlashingPenaltyAmount, nodeID)
		automation.logSlashingAction(nodeID, "Penalty Applied", fmt.Sprintf("Tokens Deducted: %d", SlashingPenaltyAmount))
	}
}

// logSlashingAction securely logs actions related to slashing penalty enforcement
func (automation *SlashingPenaltyEnforcementAutomation) logSlashingAction(nodeID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Node ID: %s, Details: %s", action, nodeID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("slashing-penalty-enforcement-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Slashing Penalty Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log slashing penalty enforcement action for node %s: %v\n", nodeID, err)
	} else {
		fmt.Println("Slashing penalty enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *SlashingPenaltyEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualSlashingPenalty allows administrators to manually apply a slashing penalty to a node
func (automation *SlashingPenaltyEnforcementAutomation) TriggerManualSlashingPenalty(nodeID string) {
	fmt.Printf("Manually applying slashing penalty to node: %s\n", nodeID)

	if automation.nodeViolationMap[nodeID] >= MaxViolationThreshold {
		automation.slashNode(nodeID)
	} else {
		fmt.Printf("Node %s does not meet the violation threshold for a slashing penalty.\n", nodeID)
		automation.logSlashingAction(nodeID, "Manual Penalty Skipped", "Insufficient Violations")
	}
}
