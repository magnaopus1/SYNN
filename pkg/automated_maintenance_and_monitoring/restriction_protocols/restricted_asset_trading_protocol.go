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
	TradingCheckInterval       = 30 * time.Second // Interval for checking asset trading restrictions
	MaxAllowedTradeViolations  = 5                // Maximum allowed trading violations before restriction
)

// RestrictedAssetTradingAutomation enforces asset trading restrictions based on predefined criteria
type RestrictedAssetTradingAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	tradeViolationCount   map[string]int // Tracks asset trading violations per user or entity
}

// NewRestrictedAssetTradingAutomation initializes and returns an instance of RestrictedAssetTradingAutomation
func NewRestrictedAssetTradingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RestrictedAssetTradingAutomation {
	return &RestrictedAssetTradingAutomation{
		consensusSystem:     consensusSystem,
		ledgerInstance:      ledgerInstance,
		stateMutex:          stateMutex,
		tradeViolationCount: make(map[string]int),
	}
}

// StartTradeMonitoring starts continuous monitoring of restricted asset trading activity
func (automation *RestrictedAssetTradingAutomation) StartTradeMonitoring() {
	ticker := time.NewTicker(TradingCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorAssetTrading()
		}
	}()
}

// monitorAssetTrading checks for asset trading violations and enforces restrictions if necessary
func (automation *RestrictedAssetTradingAutomation) monitorAssetTrading() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch trading data from Synnergy Consensus
	tradingData := automation.consensusSystem.GetTradingActivity()

	for userID, tradeStatus := range tradingData {
		// Check if the user has violated trading restrictions
		if tradeStatus == "violation" {
			automation.flagTradeViolation(userID, "Unapproved asset trading activity detected")
		}
	}
}

// flagTradeViolation flags a user's asset trading violation and logs it in the ledger
func (automation *RestrictedAssetTradingAutomation) flagTradeViolation(userID string, reason string) {
	fmt.Printf("Asset trading violation: User ID %s, Reason: %s\n", userID, reason)

	// Increment the violation count for the user
	automation.tradeViolationCount[userID]++

	// Log the violation in the ledger
	automation.logTradeViolation(userID, reason)

	// Check if the user has exceeded the allowed number of trading violations
	if automation.tradeViolationCount[userID] >= MaxAllowedTradeViolations {
		automation.restrictAssetTrading(userID)
	}
}

// logTradeViolation logs the flagged trading violation into the ledger with full details
func (automation *RestrictedAssetTradingAutomation) logTradeViolation(userID string, violationReason string) {
	// Create a ledger entry for trading violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("asset-trade-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Asset Trading Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated asset trading restrictions. Reason: %s", userID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptTradeData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log asset trading violation:", err)
	} else {
		fmt.Println("Asset trading violation logged.")
	}
}

// restrictAssetTrading restricts asset trading access for a user after exceeding allowed violations
func (automation *RestrictedAssetTradingAutomation) restrictAssetTrading(userID string) {
	// Add restriction details to the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("asset-trade-restriction-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Asset Trading Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("User %s has been restricted from asset trading due to repeated violations.", userID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptTradeData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log asset trading restriction:", err)
	} else {
		fmt.Println("Asset trading restriction applied.")
	}
}

// encryptTradeData encrypts the trading data before logging for security
func (automation *RestrictedAssetTradingAutomation) encryptTradeData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting trading data:", err)
		return data
	}
	return string(encryptedData)
}
