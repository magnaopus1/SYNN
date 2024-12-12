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
	CrossChainTransactionCheckInterval = 10 * time.Second  // Interval for checking cross-chain transactions
	MaxTransactionsPerUser             = 100               // Maximum cross-chain transactions allowed per user in a set period
	CrossChainTransactionWindow        = 24 * 7 * time.Hour // Time window for counting cross-chain transactions (1 week)
	MaxTransactionAmount               = 100000.0           // Maximum transaction amount allowed per transaction
)

// CrossChainTransactionRestrictionAutomation monitors and restricts cross-chain transactions across the network
type CrossChainTransactionRestrictionAutomation struct {
	consensusSystem             *consensus.SynnergyConsensus
	ledgerInstance              *ledger.Ledger
	stateMutex                  *sync.RWMutex
	userCrossChainTransactionCount map[string]int // Tracks cross-chain transaction count per user
}

// NewCrossChainTransactionRestrictionAutomation initializes and returns an instance of CrossChainTransactionRestrictionAutomation
func NewCrossChainTransactionRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CrossChainTransactionRestrictionAutomation {
	return &CrossChainTransactionRestrictionAutomation{
		consensusSystem:             consensusSystem,
		ledgerInstance:              ledgerInstance,
		stateMutex:                  stateMutex,
		userCrossChainTransactionCount: make(map[string]int),
	}
}

// StartTransactionMonitoring begins continuous monitoring of cross-chain transactions
func (automation *CrossChainTransactionRestrictionAutomation) StartTransactionMonitoring() {
	ticker := time.NewTicker(CrossChainTransactionCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorCrossChainTransactions()
		}
	}()
}

// monitorCrossChainTransactions checks recent cross-chain transactions and enforces transaction limits
func (automation *CrossChainTransactionRestrictionAutomation) monitorCrossChainTransactions() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent cross-chain transactions from Synnergy Consensus
	recentTransactions := automation.consensusSystem.GetRecentCrossChainTransactions()

	for _, transaction := range recentTransactions {
		// Validate transaction limits
		if !automation.validateTransactionLimit(transaction) {
			automation.flagTransactionViolation(transaction, "Exceeded maximum number of cross-chain transactions for this user")
		} else if !automation.validateTransactionAmount(transaction) {
			automation.flagTransactionViolation(transaction, "Transaction amount exceeds the maximum allowed limit")
		}
	}
}

// validateTransactionLimit checks if a user has exceeded the cross-chain transaction limit within the time window
func (automation *CrossChainTransactionRestrictionAutomation) validateTransactionLimit(transaction common.CrossChainTransaction) bool {
	currentTransactionCount := automation.userCrossChainTransactionCount[transaction.UserID]
	if currentTransactionCount+1 > MaxTransactionsPerUser {
		return false
	}

	// Update the transaction count for the user
	automation.userCrossChainTransactionCount[transaction.UserID]++
	return true
}

// validateTransactionAmount checks if a transaction amount exceeds the maximum allowed amount
func (automation *CrossChainTransactionRestrictionAutomation) validateTransactionAmount(transaction common.CrossChainTransaction) bool {
	return transaction.Amount <= MaxTransactionAmount
}

// flagTransactionViolation flags a cross-chain transaction that violates system rules and logs it in the ledger
func (automation *CrossChainTransactionRestrictionAutomation) flagTransactionViolation(transaction common.CrossChainTransaction, reason string) {
	fmt.Printf("Cross-chain transaction violation: User %s, Reason: %s\n", transaction.UserID, reason)

	// Log the violation into the ledger
	automation.logTransactionViolation(transaction, reason)
}

// logTransactionViolation logs the flagged cross-chain transaction violation into the ledger with full details
func (automation *CrossChainTransactionRestrictionAutomation) logTransactionViolation(transaction common.CrossChainTransaction, violationReason string) {
	// Encrypt the cross-chain transaction data
	encryptedData := automation.encryptTransactionData(transaction)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("cross-chain-transaction-violation-%s-%d", transaction.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Cross-Chain Transaction Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for cross-chain transaction violation. Reason: %s. Encrypted Data: %s", transaction.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log cross-chain transaction violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Cross-chain transaction violation logged for user: %s\n", transaction.UserID)
	}
}

// encryptTransactionData encrypts cross-chain transaction data before logging for security
func (automation *CrossChainTransactionRestrictionAutomation) encryptTransactionData(transaction common.CrossChainTransaction) string {
	data := fmt.Sprintf("User ID: %s, Amount: %.2f, Timestamp: %d", transaction.UserID, transaction.Amount, transaction.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting transaction data:", err)
		return data
	}
	return string(encryptedData)
}
