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
	NodeDataAccessCheckInterval = 10 * time.Second  // Interval for checking cross-node data access
	MaxDataRequestsPerNode      = 500               // Maximum data requests allowed per node in a set period
	NodeAccessTimeWindow        = 24 * 7 * time.Hour // Time window for tracking cross-node data requests (1 week)
	MaxDataAccessSize           = 5000.0            // Maximum data access size in bytes allowed per request
)

// CrossNodeDataAccessRestrictionAutomation monitors and restricts cross-node data access across the network
type CrossNodeDataAccessRestrictionAutomation struct {
	consensusSystem          *consensus.SynnergyConsensus
	ledgerInstance           *ledger.Ledger
	stateMutex               *sync.RWMutex
	nodeDataAccessCount      map[string]int // Tracks data access requests per node
}

// NewCrossNodeDataAccessRestrictionAutomation initializes and returns an instance of CrossNodeDataAccessRestrictionAutomation
func NewCrossNodeDataAccessRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CrossNodeDataAccessRestrictionAutomation {
	return &CrossNodeDataAccessRestrictionAutomation{
		consensusSystem:         consensusSystem,
		ledgerInstance:          ledgerInstance,
		stateMutex:              stateMutex,
		nodeDataAccessCount:     make(map[string]int),
	}
}

// StartDataAccessMonitoring begins continuous monitoring of cross-node data access
func (automation *CrossNodeDataAccessRestrictionAutomation) StartDataAccessMonitoring() {
	ticker := time.NewTicker(NodeDataAccessCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorNodeDataAccess()
		}
	}()
}

// monitorNodeDataAccess checks recent data access requests and enforces access limits
func (automation *CrossNodeDataAccessRestrictionAutomation) monitorNodeDataAccess() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent data access requests from Synnergy Consensus
	recentDataRequests := automation.consensusSystem.GetRecentDataRequests()

	for _, request := range recentDataRequests {
		// Validate data access limits
		if !automation.validateDataAccessLimit(request) {
			automation.flagDataAccessViolation(request, "Exceeded maximum data requests for this node")
		} else if !automation.validateDataAccessSize(request) {
			automation.flagDataAccessViolation(request, "Data access size exceeds the maximum allowed limit")
		}
	}
}

// validateDataAccessLimit checks if a node has exceeded the data access request limit within the time window
func (automation *CrossNodeDataAccessRestrictionAutomation) validateDataAccessLimit(request common.DataAccessRequest) bool {
	currentAccessCount := automation.nodeDataAccessCount[request.NodeID]
	if currentAccessCount+1 > MaxDataRequestsPerNode {
		return false
	}

	// Update the access count for the node
	automation.nodeDataAccessCount[request.NodeID]++
	return true
}

// validateDataAccessSize checks if the data access size exceeds the maximum allowed size
func (automation *CrossNodeDataAccessRestrictionAutomation) validateDataAccessSize(request common.DataAccessRequest) bool {
	return request.DataSize <= MaxDataAccessSize
}

// flagDataAccessViolation flags a data access request that violates system rules and logs it in the ledger
func (automation *CrossNodeDataAccessRestrictionAutomation) flagDataAccessViolation(request common.DataAccessRequest, reason string) {
	fmt.Printf("Data access violation: Node %s, Reason: %s\n", request.NodeID, reason)

	// Log the violation into the ledger
	automation.logDataAccessViolation(request, reason)
}

// logDataAccessViolation logs the flagged data access violation into the ledger with full details
func (automation *CrossNodeDataAccessRestrictionAutomation) logDataAccessViolation(request common.DataAccessRequest, violationReason string) {
	// Encrypt the data access request details
	encryptedData := automation.encryptDataAccessRequest(request)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("data-access-violation-%s-%d", request.NodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Cross-Node Data Access Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Node %s flagged for data access violation. Reason: %s. Encrypted Data: %s", request.NodeID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log data access violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Data access violation logged for node: %s\n", request.NodeID)
	}
}

// encryptDataAccessRequest encrypts data access request details before logging for security
func (automation *CrossNodeDataAccessRestrictionAutomation) encryptDataAccessRequest(request common.DataAccessRequest) string {
	data := fmt.Sprintf("Node ID: %s, Data Size: %.2f, Timestamp: %d", request.NodeID, request.DataSize, request.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data access request:", err)
		return data
	}
	return string(encryptedData)
}
