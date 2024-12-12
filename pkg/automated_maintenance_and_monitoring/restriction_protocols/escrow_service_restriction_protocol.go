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
	EscrowCheckInterval         = 10 * time.Second // Interval for checking escrow transactions
	MaxEscrowTransactionsPerUser = 5                // Maximum number of active escrow transactions allowed per user
	MaxEscrowAmount              = 100000.0         // Maximum amount allowed for an escrow transaction
)

// EscrowServiceRestrictionAutomation monitors and restricts escrow service transactions across the network
type EscrowServiceRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	userEscrowTransactionCount map[string]int // Tracks active escrow transactions per user
}

// NewEscrowServiceRestrictionAutomation initializes and returns an instance of EscrowServiceRestrictionAutomation
func NewEscrowServiceRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *EscrowServiceRestrictionAutomation {
	return &EscrowServiceRestrictionAutomation{
		consensusSystem:        consensusSystem,
		ledgerInstance:         ledgerInstance,
		stateMutex:             stateMutex,
		userEscrowTransactionCount: make(map[string]int),
	}
}

// StartEscrowMonitoring starts continuous monitoring of escrow transactions
func (automation *EscrowServiceRestrictionAutomation) StartEscrowMonitoring() {
	ticker := time.NewTicker(EscrowCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorEscrowTransactions()
		}
	}()
}

// monitorEscrowTransactions checks recent escrow transactions and enforces restrictions
func (automation *EscrowServiceRestrictionAutomation) monitorEscrowTransactions() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent escrow transactions from Synnergy Consensus
	recentEscrows := automation.consensusSystem.GetRecentEscrowTransactions()

	for _, escrow := range recentEscrows {
		// Validate escrow transaction limits and amounts
		if !automation.validateEscrowTransactionLimit(escrow) {
			automation.flagEscrowViolation(escrow, "Exceeded maximum number of active escrow transactions per user")
		} else if !automation.validateEscrowAmount(escrow) {
			automation.flagEscrowViolation(escrow, "Escrow amount exceeds the maximum allowed limit")
		}
	}
}

// validateEscrowTransactionLimit checks if a user has exceeded the active escrow transaction limit
func (automation *EscrowServiceRestrictionAutomation) validateEscrowTransactionLimit(escrow common.EscrowTransaction) bool {
	currentEscrowCount := automation.userEscrowTransactionCount[escrow.UserID]
	if currentEscrowCount+1 > MaxEscrowTransactionsPerUser {
		return false
	}

	// Update the escrow transaction count for the user
	automation.userEscrowTransactionCount[escrow.UserID]++
	return true
}

// validateEscrowAmount checks if the escrow amount exceeds the maximum allowed limit
func (automation *EscrowServiceRestrictionAutomation) validateEscrowAmount(escrow common.EscrowTransaction) bool {
	return escrow.Amount <= MaxEscrowAmount
}

// flagEscrowViolation flags an escrow transaction that violates system rules and logs it in the ledger
func (automation *EscrowServiceRestrictionAutomation) flagEscrowViolation(escrow common.EscrowTransaction, reason string) {
	fmt.Printf("Escrow transaction violation: User %s, Reason: %s\n", escrow.UserID, reason)

	// Log the violation into the ledger
	automation.logEscrowViolation(escrow, reason)
}

// logEscrowViolation logs the flagged escrow transaction violation into the ledger with full details
func (automation *EscrowServiceRestrictionAutomation) logEscrowViolation(escrow common.EscrowTransaction, violationReason string) {
	// Encrypt the escrow transaction data before logging
	encryptedData := automation.encryptEscrowData(escrow)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("escrow-violation-%s-%d", escrow.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Escrow Transaction Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for escrow transaction violation. Reason: %s. Encrypted Data: %s", escrow.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log escrow transaction violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Escrow transaction violation logged for user: %s\n", escrow.UserID)
	}
}

// encryptEscrowData encrypts escrow transaction data before logging for security
func (automation *EscrowServiceRestrictionAutomation) encryptEscrowData(escrow common.EscrowTransaction) string {
	data := fmt.Sprintf("User ID: %s, Escrow Amount: %.2f, Timestamp: %d", escrow.UserID, escrow.Amount, escrow.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting escrow data:", err)
		return data
	}
	return string(encryptedData)
}
