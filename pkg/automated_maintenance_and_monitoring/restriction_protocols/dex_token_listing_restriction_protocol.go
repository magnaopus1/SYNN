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
	TokenListingCheckInterval      = 10 * time.Second  // Interval for checking token listing activities
	MaxListingsPerUser             = 3                 // Maximum number of token listings allowed per user
	ListingAuditRequirementEnabled = true              // If true, all new token listings must pass audit
)

// DexTokenListingRestrictionAutomation monitors and restricts new token listings on decentralized exchanges (DEX)
type DexTokenListingRestrictionAutomation struct {
	consensusSystem           *consensus.SynnergyConsensus
	ledgerInstance            *ledger.Ledger
	stateMutex                *sync.RWMutex
	userTokenListingCount     map[string]int // Tracks token listings per user
}

// NewDexTokenListingRestrictionAutomation initializes and returns an instance of DexTokenListingRestrictionAutomation
func NewDexTokenListingRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DexTokenListingRestrictionAutomation {
	return &DexTokenListingRestrictionAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		userTokenListingCount: make(map[string]int),
	}
}

// StartTokenListingMonitoring starts continuous monitoring of token listing activities
func (automation *DexTokenListingRestrictionAutomation) StartTokenListingMonitoring() {
	ticker := time.NewTicker(TokenListingCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorTokenListingActivities()
		}
	}()
}

// monitorTokenListingActivities checks recent token listing activities and enforces restrictions on token listings
func (automation *DexTokenListingRestrictionAutomation) monitorTokenListingActivities() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent token listings from Synnergy Consensus
	recentListings := automation.consensusSystem.GetRecentTokenListings()

	for _, listing := range recentListings {
		// Validate token listing limits and audit requirements
		if !automation.validateListingLimit(listing) {
			automation.flagTokenListingViolation(listing, "Exceeded maximum number of token listings per user")
		} else if ListingAuditRequirementEnabled && !automation.validateTokenAudit(listing) {
			automation.flagTokenListingViolation(listing, "Token listing failed audit requirement")
		}
	}
}

// validateListingLimit checks if a user has exceeded the token listing limit
func (automation *DexTokenListingRestrictionAutomation) validateListingLimit(listing common.TokenListing) bool {
	currentListingCount := automation.userTokenListingCount[listing.UserID]
	if currentListingCount+1 > MaxListingsPerUser {
		return false
	}

	// Update the listing count for the user
	automation.userTokenListingCount[listing.UserID]++
	return true
}

// validateTokenAudit checks if the token passed the audit requirement for listing
func (automation *DexTokenListingRestrictionAutomation) validateTokenAudit(listing common.TokenListing) bool {
	return listing.AuditPassed
}

// flagTokenListingViolation flags a token listing activity that violates system rules and logs it in the ledger
func (automation *DexTokenListingRestrictionAutomation) flagTokenListingViolation(listing common.TokenListing, reason string) {
	fmt.Printf("DEX token listing violation: User %s, Reason: %s\n", listing.UserID, reason)

	// Log the violation into the ledger
	automation.logTokenListingViolation(listing, reason)
}

// logTokenListingViolation logs the flagged token listing violation into the ledger with full details
func (automation *DexTokenListingRestrictionAutomation) logTokenListingViolation(listing common.TokenListing, violationReason string) {
	// Encrypt the token listing data before logging
	encryptedData := automation.encryptTokenListingData(listing)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("token-listing-violation-%s-%d", listing.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "DEX Token Listing Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for token listing violation. Reason: %s. Encrypted Data: %s", listing.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log token listing violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Token listing violation logged for user: %s\n", listing.UserID)
	}
}

// encryptTokenListingData encrypts token listing data before logging for security
func (automation *DexTokenListingRestrictionAutomation) encryptTokenListingData(listing common.TokenListing) string {
	data := fmt.Sprintf("User ID: %s, Token ID: %s, Timestamp: %d, Audit Status: %t", listing.UserID, listing.TokenID, listing.Timestamp, listing.AuditPassed)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting token listing data:", err)
		return data
	}
	return string(encryptedData)
}
