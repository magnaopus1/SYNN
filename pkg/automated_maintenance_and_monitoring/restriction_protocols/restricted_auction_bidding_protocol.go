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
	AuctionCheckInterval      = 30 * time.Second // Interval for checking auction bidding restrictions
	MaxAllowedBidViolations   = 5                // Maximum allowed bidding violations before restriction
)

// RestrictedAuctionBiddingAutomation enforces bidding restrictions in the auction system
type RestrictedAuctionBiddingAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	bidViolationCount     map[string]int // Tracks auction bidding violations per user or entity
}

// NewRestrictedAuctionBiddingAutomation initializes and returns an instance of RestrictedAuctionBiddingAutomation
func NewRestrictedAuctionBiddingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RestrictedAuctionBiddingAutomation {
	return &RestrictedAuctionBiddingAutomation{
		consensusSystem:     consensusSystem,
		ledgerInstance:      ledgerInstance,
		stateMutex:          stateMutex,
		bidViolationCount:   make(map[string]int),
	}
}

// StartBidMonitoring starts continuous monitoring of restricted auction bidding activity
func (automation *RestrictedAuctionBiddingAutomation) StartBidMonitoring() {
	ticker := time.NewTicker(AuctionCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorAuctionBidding()
		}
	}()
}

// monitorAuctionBidding checks for auction bidding violations and enforces restrictions if necessary
func (automation *RestrictedAuctionBiddingAutomation) monitorAuctionBidding() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch bidding data from Synnergy Consensus
	biddingData := automation.consensusSystem.GetAuctionBiddingData()

	for userID, bidStatus := range biddingData {
		// Check if the user has violated auction bidding rules
		if bidStatus == "violation" {
			automation.flagBidViolation(userID, "Unauthorized auction bidding activity detected")
		}
	}
}

// flagBidViolation flags a user's auction bidding violation and logs it in the ledger
func (automation *RestrictedAuctionBiddingAutomation) flagBidViolation(userID string, reason string) {
	fmt.Printf("Auction bidding violation: User ID %s, Reason: %s\n", userID, reason)

	// Increment the violation count for the user
	automation.bidViolationCount[userID]++

	// Log the violation in the ledger
	automation.logBidViolation(userID, reason)

	// Check if the user has exceeded the allowed number of bidding violations
	if automation.bidViolationCount[userID] >= MaxAllowedBidViolations {
		automation.restrictAuctionBidding(userID)
	}
}

// logBidViolation logs the flagged auction bidding violation into the ledger with full details
func (automation *RestrictedAuctionBiddingAutomation) logBidViolation(userID string, violationReason string) {
	// Create a ledger entry for bidding violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("auction-bid-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Auction Bidding Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated auction bidding rules. Reason: %s", userID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptBidData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log auction bidding violation:", err)
	} else {
		fmt.Println("Auction bidding violation logged.")
	}
}

// restrictAuctionBidding restricts auction bidding access for a user after exceeding allowed violations
func (automation *RestrictedAuctionBiddingAutomation) restrictAuctionBidding(userID string) {
	// Add restriction details to the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("auction-bid-restriction-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Auction Bidding Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("User %s has been restricted from auction bidding due to repeated violations.", userID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptBidData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log auction bidding restriction:", err)
	} else {
		fmt.Println("Auction bidding restriction applied.")
	}
}

// encryptBidData encrypts the auction bidding data before logging for security
func (automation *RestrictedAuctionBiddingAutomation) encryptBidData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting bidding data:", err)
		return data
	}
	return string(encryptedData)
}
