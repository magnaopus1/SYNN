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
	NodeAccessCheckInterval    = 10 * time.Second // Interval for checking restricted node access
	MaxUnauthorizedAttempts    = 3                // Maximum allowed unauthorized node access attempts before restriction
)

// RestrictedNodeAccessAutomation enforces restrictions on unauthorized node access
type RestrictedNodeAccessAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	accessViolationCount  map[string]int // Tracks unauthorized node access attempts per node or entity
}

// NewRestrictedNodeAccessAutomation initializes and returns an instance of RestrictedNodeAccessAutomation
func NewRestrictedNodeAccessAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RestrictedNodeAccessAutomation {
	return &RestrictedNodeAccessAutomation{
		consensusSystem:      consensusSystem,
		ledgerInstance:       ledgerInstance,
		stateMutex:           stateMutex,
		accessViolationCount: make(map[string]int),
	}
}

// StartNodeAccessMonitoring starts continuous monitoring of node access violations
func (automation *RestrictedNodeAccessAutomation) StartNodeAccessMonitoring() {
	ticker := time.NewTicker(NodeAccessCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorNodeAccess()
		}
	}()
}

// monitorNodeAccess checks for unauthorized node access attempts and enforces restrictions if necessary
func (automation *RestrictedNodeAccessAutomation) monitorNodeAccess() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch node access data from Synnergy Consensus
	nodeAccessData := automation.consensusSystem.GetNodeAccessData()

	for nodeID, accessStatus := range nodeAccessData {
		// Check if the access attempt is unauthorized
		if accessStatus == "unauthorized" {
			automation.flagAccessViolation(nodeID, "Unauthorized node access attempt detected")
		}
	}
}

// flagAccessViolation flags a node's unauthorized access attempt and logs it in the ledger
func (automation *RestrictedNodeAccessAutomation) flagAccessViolation(nodeID string, reason string) {
	fmt.Printf("Node access violation: Node ID %s, Reason: %s\n", nodeID, reason)

	// Increment the violation count for the node
	automation.accessViolationCount[nodeID]++

	// Log the violation in the ledger
	automation.logAccessViolation(nodeID, reason)

	// Check if the node has exceeded the allowed number of access violations
	if automation.accessViolationCount[nodeID] >= MaxUnauthorizedAttempts {
		automation.restrictNodeAccess(nodeID)
	}
}

// logAccessViolation logs the flagged node access violation into the ledger with full details
func (automation *RestrictedNodeAccessAutomation) logAccessViolation(nodeID string, violationReason string) {
	// Create a ledger entry for node access violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("node-access-violation-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Node Access Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Node %s violated node access restrictions. Reason: %s", nodeID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptAccessData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log node access violation:", err)
	} else {
		fmt.Println("Node access violation logged.")
	}
}

// restrictNodeAccess restricts node access after exceeding allowed violations
func (automation *RestrictedNodeAccessAutomation) restrictNodeAccess(nodeID string) {
	// Add restriction details to the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("node-access-restriction-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Node Access Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Node %s has been restricted from network access due to repeated violations.", nodeID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptAccessData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log node access restriction:", err)
	} else {
		fmt.Println("Node access restriction applied.")
	}
}

// encryptAccessData encrypts the node access data before logging for security
func (automation *RestrictedNodeAccessAutomation) encryptAccessData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting node access data:", err)
		return data
	}
	return string(encryptedData)
}
