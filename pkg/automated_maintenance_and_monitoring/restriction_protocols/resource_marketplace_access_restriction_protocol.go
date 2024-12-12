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
	AccessCheckInterval     = 30 * time.Second // Interval for checking marketplace access
	MaxAllowedAccessViolations = 5             // Maximum allowed access violations before restriction
)

// ResourceMarketplaceAccessRestrictionAutomation monitors and enforces access rules for the resource marketplace
type ResourceMarketplaceAccessRestrictionAutomation struct {
	consensusSystem          *consensus.SynnergyConsensus
	ledgerInstance           *ledger.Ledger
	stateMutex               *sync.RWMutex
	accessViolationCount     map[string]int // Tracks access violation counts per user or node
}

// NewResourceMarketplaceAccessRestrictionAutomation initializes and returns an instance of ResourceMarketplaceAccessRestrictionAutomation
func NewResourceMarketplaceAccessRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ResourceMarketplaceAccessRestrictionAutomation {
	return &ResourceMarketplaceAccessRestrictionAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		accessViolationCount:  make(map[string]int),
	}
}

// StartAccessMonitoring starts continuous monitoring of marketplace access compliance
func (automation *ResourceMarketplaceAccessRestrictionAutomation) StartAccessMonitoring() {
	ticker := time.NewTicker(AccessCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorMarketplaceAccess()
		}
	}()
}

// monitorMarketplaceAccess checks for access violations and enforces restrictions if necessary
func (automation *ResourceMarketplaceAccessRestrictionAutomation) monitorMarketplaceAccess() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch access data from Synnergy Consensus
	accessData := automation.consensusSystem.GetMarketplaceAccessData()

	for userID, accessStatus := range accessData {
		// Check if the user has violated access rules
		if accessStatus == "violation" {
			automation.flagAccessViolation(userID, "Unauthorized marketplace access")
		}
	}
}

// flagAccessViolation flags a user's access violation and logs it in the ledger
func (automation *ResourceMarketplaceAccessRestrictionAutomation) flagAccessViolation(userID string, reason string) {
	fmt.Printf("Marketplace access violation: User ID %s, Reason: %s\n", userID, reason)

	// Increment the violation count for the user
	automation.accessViolationCount[userID]++

	// Log the violation in the ledger
	automation.logAccessViolation(userID, reason)

	// Check if the user has exceeded the allowed number of access violations
	if automation.accessViolationCount[userID] >= MaxAllowedAccessViolations {
		automation.restrictMarketplaceAccess(userID)
	}
}

// logAccessViolation logs the flagged marketplace access violation into the ledger with full details
func (automation *ResourceMarketplaceAccessRestrictionAutomation) logAccessViolation(userID string, violationReason string) {
	// Create a ledger entry for access violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("marketplace-access-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Marketplace Access Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated marketplace access rules. Reason: %s", userID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptAccessData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log marketplace access violation:", err)
	} else {
		fmt.Println("Marketplace access violation logged.")
	}
}

// restrictMarketplaceAccess restricts access for a user after exceeding allowed violations
func (automation *ResourceMarketplaceAccessRestrictionAutomation) restrictMarketplaceAccess(userID string) {
	// Add restriction details to the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("marketplace-access-restriction-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Marketplace Access Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("User %s has been restricted from accessing the marketplace due to repeated violations.", userID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptAccessData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log marketplace access restriction:", err)
	} else {
		fmt.Println("Marketplace access restriction applied.")
	}
}

// encryptAccessData encrypts the marketplace access data before logging for security
func (automation *ResourceMarketplaceAccessRestrictionAutomation) encryptAccessData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting marketplace access data:", err)
		return data
	}
	return string(encryptedData)
}
