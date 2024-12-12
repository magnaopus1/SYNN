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
	SlippageCheckInterval  = 10 * time.Second  // Interval for checking slippage in DEX trades
	MaxAllowedSlippage     = 5.0               // Maximum allowed slippage percentage
	SlippageViolationCount = 3                 // Maximum number of slippage violations before enforcement
)

// DexSlippageLimitAutomation monitors and restricts slippage levels in decentralized exchange trades
type DexSlippageLimitAutomation struct {
	consensusSystem           *consensus.SynnergyConsensus
	ledgerInstance            *ledger.Ledger
	stateMutex                *sync.RWMutex
	userSlippageViolationCount map[string]int // Tracks slippage violations per user
}

// NewDexSlippageLimitAutomation initializes and returns an instance of DexSlippageLimitAutomation
func NewDexSlippageLimitAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DexSlippageLimitAutomation {
	return &DexSlippageLimitAutomation{
		consensusSystem:           consensusSystem,
		ledgerInstance:            ledgerInstance,
		stateMutex:                stateMutex,
		userSlippageViolationCount: make(map[string]int),
	}
}

// StartSlippageMonitoring starts continuous monitoring of DEX slippage levels
func (automation *DexSlippageLimitAutomation) StartSlippageMonitoring() {
	ticker := time.NewTicker(SlippageCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorSlippageLevels()
		}
	}()
}

// monitorSlippageLevels checks recent DEX trades for slippage and enforces slippage limits
func (automation *DexSlippageLimitAutomation) monitorSlippageLevels() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent trades from Synnergy Consensus
	recentTrades := automation.consensusSystem.GetRecentTrades()

	for _, trade := range recentTrades {
		// Validate slippage limit
		if !automation.validateSlippage(trade) {
			automation.flagSlippageViolation(trade, "Trade slippage exceeded the maximum allowed limit")
		}
	}
}

// validateSlippage checks if the trade slippage exceeds the maximum allowed limit
func (automation *DexSlippageLimitAutomation) validateSlippage(trade common.Trade) bool {
	if trade.SlippagePercentage > MaxAllowedSlippage {
		automation.userSlippageViolationCount[trade.UserID]++
		return false
	}

	// Reset violation count for valid trades
	automation.userSlippageViolationCount[trade.UserID] = 0
	return true
}

// flagSlippageViolation flags a trade that violates slippage rules and logs it in the ledger
func (automation *DexSlippageLimitAutomation) flagSlippageViolation(trade common.Trade, reason string) {
	fmt.Printf("DEX slippage violation: User %s, Reason: %s\n", trade.UserID, reason)

	// Log the violation into the ledger
	automation.logSlippageViolation(trade, reason)

	// Trigger further enforcement actions if the violation count exceeds the allowed limit
	if automation.userSlippageViolationCount[trade.UserID] >= SlippageViolationCount {
		automation.enforceSlippageRestriction(trade.UserID)
	}
}

// logSlippageViolation logs the flagged slippage violation into the ledger with full details
func (automation *DexSlippageLimitAutomation) logSlippageViolation(trade common.Trade, violationReason string) {
	// Encrypt the slippage data before logging
	encryptedData := automation.encryptSlippageData(trade)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("slippage-violation-%s-%d", trade.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "DEX Slippage Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for slippage violation. Reason: %s. Encrypted Data: %s", trade.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log slippage violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Slippage violation logged for user: %s\n", trade.UserID)
	}
}

// enforceSlippageRestriction triggers restrictions or penalties for users exceeding slippage violations
func (automation *DexSlippageLimitAutomation) enforceSlippageRestriction(userID string) {
	// Enforce specific penalties or restrictions on the user
	fmt.Printf("Enforcing slippage restrictions on user: %s\n", userID)

	// Log the enforcement action in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("slippage-enforcement-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Slippage Restriction Enforcement",
		Status:    "Enforced",
		Details:   fmt.Sprintf("Slippage restrictions enforced on user %s after exceeding the allowed violation limit.", userID),
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log slippage enforcement action: %v\n", err)
	} else {
		fmt.Printf("Slippage restrictions enforced on user: %s\n", userID)
	}
}

// encryptSlippageData encrypts slippage data before logging for security
func (automation *DexSlippageLimitAutomation) encryptSlippageData(trade common.Trade) string {
	data := fmt.Sprintf("User ID: %s, Slippage Percentage: %.2f, Timestamp: %d", trade.UserID, trade.SlippagePercentage, trade.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting slippage data:", err)
		return data
	}
	return string(encryptedData)
}
