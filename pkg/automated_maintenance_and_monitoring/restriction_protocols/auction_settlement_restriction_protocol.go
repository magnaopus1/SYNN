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
	AuctionSettlementCheckInterval = 5 * time.Second // Interval for checking auction settlements
	MaxAllowedSettlementTime       = 24 * time.Hour  // Maximum time allowed for auction settlements
	MinSettlementAmount            = 10.0            // Minimum allowable settlement amount
)

// AuctionSettlementRestrictionAutomation manages and enforces auction settlement restrictions across the network
type AuctionSettlementRestrictionAutomation struct {
	consensusSystem      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	stateMutex           *sync.RWMutex
	flaggedSettlements   map[string]int // Tracks flagged auctions by ID
}

// NewAuctionSettlementRestrictionAutomation initializes and returns an instance of AuctionSettlementRestrictionAutomation
func NewAuctionSettlementRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AuctionSettlementRestrictionAutomation {
	return &AuctionSettlementRestrictionAutomation{
		consensusSystem:    consensusSystem,
		ledgerInstance:     ledgerInstance,
		stateMutex:         stateMutex,
		flaggedSettlements: make(map[string]int),
	}
}

// StartAuctionSettlementMonitoring starts the continuous loop for monitoring auction settlements
func (automation *AuctionSettlementRestrictionAutomation) StartAuctionSettlementMonitoring() {
	ticker := time.NewTicker(AuctionSettlementCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorAuctionSettlements()
		}
	}()
}

// monitorAuctionSettlements checks the recent auction settlements and enforces restrictions
func (automation *AuctionSettlementRestrictionAutomation) monitorAuctionSettlements() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent auction settlements from Synnergy Consensus
	recentSettlements := automation.consensusSystem.GetRecentAuctionSettlements()

	for _, settlement := range recentSettlements {
		// Check for rule violations such as minimum settlement amount and settlement time
		if !automation.validateSettlementAmount(settlement) {
			automation.flagSettlement(settlement, "Settlement amount below the minimum threshold")
			continue
		}

		if !automation.validateSettlementTime(settlement) {
			automation.flagSettlement(settlement, "Settlement time exceeded the maximum allowable time")
		}
	}
}

// validateSettlementAmount checks if the auction settlement meets the minimum amount
func (automation *AuctionSettlementRestrictionAutomation) validateSettlementAmount(settlement common.AuctionSettlement) bool {
	return settlement.Amount >= MinSettlementAmount
}

// validateSettlementTime checks if the auction settlement occurred within the maximum allowable time
func (automation *AuctionSettlementRestrictionAutomation) validateSettlementTime(settlement common.AuctionSettlement) bool {
	settlementDuration := time.Since(settlement.AuctionEndTime)
	return settlementDuration <= MaxAllowedSettlementTime
}

// flagSettlement flags a settlement that violates auction settlement rules and logs it in the ledger
func (automation *AuctionSettlementRestrictionAutomation) flagSettlement(settlement common.AuctionSettlement, reason string) {
	fmt.Printf("Auction settlement flagged: %s, Reason: %s\n", settlement.AuctionID, reason)

	// Track the number of flags for this auction settlement
	automation.flaggedSettlements[settlement.AuctionID]++

	// Log the violation into the ledger
	automation.logSettlementViolation(settlement, reason)
}

// logSettlementViolation logs the flagged settlement violation into the ledger with full details
func (automation *AuctionSettlementRestrictionAutomation) logSettlementViolation(settlement common.AuctionSettlement, violationReason string) {
	// Encrypt settlement details before logging
	encryptedData := automation.encryptSettlementData(settlement)

	// Create ledger entry with the settlement violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("auction-settlement-violation-%s-%d", settlement.AuctionID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Auction Settlement Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Auction (%s) flagged for settlement violation. Reason: %s. Encrypted Data: %s", settlement.AuctionID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log settlement violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Settlement violation logged for auction: %s\n", settlement.AuctionID)
	}
}

// encryptSettlementData encrypts settlement data for security and compliance before logging
func (automation *AuctionSettlementRestrictionAutomation) encryptSettlementData(settlement common.AuctionSettlement) string {
	data := fmt.Sprintf("Auction ID: %s, Buyer: %s, Amount: %.2f, End Time: %s", settlement.AuctionID, settlement.Buyer, settlement.Amount, settlement.AuctionEndTime)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting settlement data:", err)
		return data
	}
	return string(encryptedData)
}
