package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
)

// Constants for validator node eligibility checks
const (
	ValidatorEligibilityCheckInterval = 3 * time.Second // Interval for eligibility checks
	MinimumStakeRequired              = 1000            // Minimum tokens required to be eligible as a validator
	MinimumUptimePercentage           = 95              // Minimum uptime required for validators
)

// ValidatorNodeEligibilityAutomation ensures that nodes meet eligibility criteria to become validators
type ValidatorNodeEligibilityAutomation struct {
	consensusSystem  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	nodeMutex        *sync.RWMutex
}

// NewValidatorNodeEligibilityAutomation initializes the automation for validator eligibility
func NewValidatorNodeEligibilityAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, nodeMutex *sync.RWMutex) *ValidatorNodeEligibilityAutomation {
	return &ValidatorNodeEligibilityAutomation{
		consensusSystem:  consensusSystem,
		ledgerInstance:   ledgerInstance,
		nodeMutex:        nodeMutex,
	}
}

// StartMonitoring continuously monitors node eligibility for validation
func (automation *ValidatorNodeEligibilityAutomation) StartMonitoring() {
	ticker := time.NewTicker(ValidatorEligibilityCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateValidatorEligibility()
		}
	}()
}

// evaluateValidatorEligibility checks all nodes in the system to determine if they meet validator eligibility criteria
func (automation *ValidatorNodeEligibilityAutomation) evaluateValidatorEligibility() {
	automation.nodeMutex.Lock()
	defer automation.nodeMutex.Unlock()

	nodes := automation.consensusSystem.GetNodes()
	for _, node := range nodes {
		if !automation.isEligible(node) {
			automation.revokeValidatorStatus(node.ID)
		} else {
			automation.grantValidatorStatus(node.ID)
		}
	}
}

// isEligible checks if a node meets the criteria to become or remain a validator
func (automation *ValidatorNodeEligibilityAutomation) isEligible(node consensus.Node) bool {
	// Check if the node has the minimum stake required
	if node.Stake < MinimumStakeRequired {
		automation.logEligibilityFailure(node.ID, "Insufficient stake", node.Stake)
		return false
	}

	// Check if the node has the required uptime percentage
	if node.Uptime < MinimumUptimePercentage {
		automation.logEligibilityFailure(node.ID, "Low uptime percentage", node.Uptime)
		return false
	}

	return true
}

// logEligibilityFailure logs nodes that fail eligibility criteria into the ledger
func (automation *ValidatorNodeEligibilityAutomation) logEligibilityFailure(nodeID string, reason string, value interface{}) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("eligibility-failure-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Validator Eligibility Failure",
		Status:    "Failed",
		Details:   fmt.Sprintf("Node %s failed eligibility due to %s: %v.", nodeID, reason, value),
	}

	// Encrypt details before storing in the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log eligibility failure:", err)
	} else {
		fmt.Println("Validator eligibility failure logged for node:", nodeID)
	}
}

// revokeValidatorStatus revokes validator privileges from a node that no longer meets eligibility criteria
func (automation *ValidatorNodeEligibilityAutomation) revokeValidatorStatus(nodeID string) {
	err := automation.consensusSystem.RevokeValidatorStatus(nodeID)
	if err != nil {
		fmt.Println("Failed to revoke validator status for node:", nodeID, "Error:", err)
		return
	}

	// Log the validator status revocation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("validator-revocation-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Validator Status Revocation",
		Status:    "Revoked",
		Details:   fmt.Sprintf("Validator status revoked for node %s due to eligibility failure.", nodeID),
	}

	// Encrypt the revocation details
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log validator status revocation:", err)
	} else {
		fmt.Println("Validator status revocation logged for node:", nodeID)
	}
}

// grantValidatorStatus grants validator privileges to eligible nodes
func (automation *ValidatorNodeEligibilityAutomation) grantValidatorStatus(nodeID string) {
	err := automation.consensusSystem.GrantValidatorStatus(nodeID)
	if err != nil {
		fmt.Println("Failed to grant validator status for node:", nodeID, "Error:", err)
		return
	}

	// Log the validator status grant in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("validator-grant-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Validator Status Grant",
		Status:    "Granted",
		Details:   fmt.Sprintf("Validator status granted to node %s after eligibility confirmation.", nodeID),
	}

	// Encrypt the grant details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log validator status grant:", err)
	} else {
		fmt.Println("Validator status grant logged for node:", nodeID)
	}
}

// encryptData encrypts sensitive information before storing in the ledger
func (automation *ValidatorNodeEligibilityAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
