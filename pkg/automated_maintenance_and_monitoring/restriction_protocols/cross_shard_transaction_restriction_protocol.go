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
	CrossShardTransactionCheckInterval = 10 * time.Second  // Interval for checking cross-shard transactions
	MaxTransactionsPerShard            = 500               // Maximum cross-shard transactions allowed per shard in a set period
	CrossShardTransactionWindow        = 24 * 7 * time.Hour // Time window for counting cross-shard transactions (1 week)
	MaxTransactionSize                 = 10000.0           // Maximum transaction size in bytes allowed per cross-shard transaction
)

// CrossShardTransactionRestrictionAutomation monitors and restricts cross-shard transactions across the network
type CrossShardTransactionRestrictionAutomation struct {
	consensusSystem           *consensus.SynnergyConsensus
	ledgerInstance            *ledger.Ledger
	stateMutex                *sync.RWMutex
	shardTransactionCount     map[string]int // Tracks cross-shard transaction count per shard
}

// NewCrossShardTransactionRestrictionAutomation initializes and returns an instance of CrossShardTransactionRestrictionAutomation
func NewCrossShardTransactionRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CrossShardTransactionRestrictionAutomation {
	return &CrossShardTransactionRestrictionAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		shardTransactionCount: make(map[string]int),
	}
}

// StartTransactionMonitoring starts continuous monitoring of cross-shard transactions
func (automation *CrossShardTransactionRestrictionAutomation) StartTransactionMonitoring() {
	ticker := time.NewTicker(CrossShardTransactionCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorCrossShardTransactions()
		}
	}()
}

// monitorCrossShardTransactions checks recent cross-shard transactions and enforces transaction limits
func (automation *CrossShardTransactionRestrictionAutomation) monitorCrossShardTransactions() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent cross-shard transactions from Synnergy Consensus
	recentTransactions := automation.consensusSystem.GetRecentCrossShardTransactions()

	for _, transaction := range recentTransactions {
		// Validate transaction limits
		if !automation.validateTransactionLimit(transaction) {
			automation.flagTransactionViolation(transaction, "Exceeded maximum number of transactions for this shard")
		} else if !automation.validateTransactionSize(transaction) {
			automation.flagTransactionViolation(transaction, "Transaction size exceeds the maximum allowed limit")
		}
	}
}

// validateTransactionLimit checks if a shard has exceeded the transaction limit within the time window
func (automation *CrossShardTransactionRestrictionAutomation) validateTransactionLimit(transaction common.CrossShardTransaction) bool {
	currentTransactionCount := automation.shardTransactionCount[transaction.ShardID]
	if currentTransactionCount+1 > MaxTransactionsPerShard {
		return false
	}

	// Update the transaction count for the shard
	automation.shardTransactionCount[transaction.ShardID]++
	return true
}

// validateTransactionSize checks if the transaction size exceeds the maximum allowed size
func (automation *CrossShardTransactionRestrictionAutomation) validateTransactionSize(transaction common.CrossShardTransaction) bool {
	return transaction.TransactionSize <= MaxTransactionSize
}

// flagTransactionViolation flags a cross-shard transaction that violates system rules and logs it in the ledger
func (automation *CrossShardTransactionRestrictionAutomation) flagTransactionViolation(transaction common.CrossShardTransaction, reason string) {
	fmt.Printf("Cross-shard transaction violation: Shard %s, Reason: %s\n", transaction.ShardID, reason)

	// Log the violation into the ledger
	automation.logTransactionViolation(transaction, reason)
}

// logTransactionViolation logs the flagged cross-shard transaction violation into the ledger with full details
func (automation *CrossShardTransactionRestrictionAutomation) logTransactionViolation(transaction common.CrossShardTransaction, violationReason string) {
	// Encrypt the cross-shard transaction data
	encryptedData := automation.encryptTransactionData(transaction)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("cross-shard-transaction-violation-%s-%d", transaction.ShardID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Cross-Shard Transaction Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Shard %s flagged for transaction violation. Reason: %s. Encrypted Data: %s", transaction.ShardID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log cross-shard transaction violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Cross-shard transaction violation logged for shard: %s\n", transaction.ShardID)
	}
}

// encryptTransactionData encrypts cross-shard transaction data before logging for security
func (automation *CrossShardTransactionRestrictionAutomation) encryptTransactionData(transaction common.CrossShardTransaction) string {
	data := fmt.Sprintf("Shard ID: %s, Transaction Size: %.2f, Timestamp: %d", transaction.ShardID, transaction.TransactionSize, transaction.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting transaction data:", err)
		return data
	}
	return string(encryptedData)
}
