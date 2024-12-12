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
	BridgeAccessCheckInterval      = 15 * time.Second  // Interval for checking cross-chain bridge access
	MaxBridgeAccessPerUser         = 100               // Maximum number of bridge accesses allowed per user in a set period
	BridgeAccessTimeWindow         = 24 * 7 * time.Hour // Time window for counting bridge accesses (1 week)
)

// CrossChainBridgeAccessRestrictionAutomation monitors and restricts access to cross-chain bridges across the network
type CrossChainBridgeAccessRestrictionAutomation struct {
	consensusSystem         *consensus.SynnergyConsensus
	ledgerInstance          *ledger.Ledger
	stateMutex              *sync.RWMutex
	userBridgeAccessCount   map[string]int // Tracks bridge access count per user
}

// NewCrossChainBridgeAccessRestrictionAutomation initializes and returns an instance of CrossChainBridgeAccessRestrictionAutomation
func NewCrossChainBridgeAccessRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CrossChainBridgeAccessRestrictionAutomation {
	return &CrossChainBridgeAccessRestrictionAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		userBridgeAccessCount: make(map[string]int),
	}
}

// StartBridgeAccessMonitoring begins continuous monitoring of cross-chain bridge access
func (automation *CrossChainBridgeAccessRestrictionAutomation) StartBridgeAccessMonitoring() {
	ticker := time.NewTicker(BridgeAccessCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorBridgeAccess()
		}
	}()
}

// monitorBridgeAccess checks recent bridge accesses and enforces access limits
func (automation *CrossChainBridgeAccessRestrictionAutomation) monitorBridgeAccess() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent bridge accesses from Synnergy Consensus
	recentAccesses := automation.consensusSystem.GetRecentBridgeAccesses()

	for _, access := range recentAccesses {
		// Validate bridge access limits
		if !automation.validateBridgeAccessLimit(access) {
			automation.flagBridgeAccessViolation(access, "Exceeded maximum number of bridge accesses for this user")
		}
	}
}

// validateBridgeAccessLimit checks if a user has exceeded the bridge access limit within the time window
func (automation *CrossChainBridgeAccessRestrictionAutomation) validateBridgeAccessLimit(access common.BridgeAccess) bool {
	currentAccessCount := automation.userBridgeAccessCount[access.UserID]
	if currentAccessCount+1 > MaxBridgeAccessPerUser {
		return false
	}

	// Update the access count for the user
	automation.userBridgeAccessCount[access.UserID]++
	return true
}

// flagBridgeAccessViolation flags a bridge access that violates the system's rules and logs it in the ledger
func (automation *CrossChainBridgeAccessRestrictionAutomation) flagBridgeAccessViolation(access common.BridgeAccess, reason string) {
	fmt.Printf("Bridge access violation: User %s, Reason: %s\n", access.UserID, reason)

	// Log the violation into the ledger
	automation.logBridgeAccessViolation(access, reason)
}

// logBridgeAccessViolation logs the flagged bridge access violation into the ledger with full details
func (automation *CrossChainBridgeAccessRestrictionAutomation) logBridgeAccessViolation(access common.BridgeAccess, violationReason string) {
	// Encrypt the bridge access data
	encryptedData := automation.encryptBridgeAccessData(access)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("bridge-access-violation-%s-%d", access.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Cross-Chain Bridge Access Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for bridge access violation. Reason: %s. Encrypted Data: %s", access.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log bridge access violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Bridge access violation logged for user: %s\n", access.UserID)
	}
}

// encryptBridgeAccessData encrypts bridge access data before logging for security
func (automation *CrossChainBridgeAccessRestrictionAutomation) encryptBridgeAccessData(access common.BridgeAccess) string {
	data := fmt.Sprintf("User ID: %s, Timestamp: %d", access.UserID, access.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting bridge access data:", err)
		return data
	}
	return string(encryptedData)
}
