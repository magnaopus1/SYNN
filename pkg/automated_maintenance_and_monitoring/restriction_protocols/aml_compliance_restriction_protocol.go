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
	AMLMonitoringInterval        = 1 * time.Hour  // Regular interval for AML compliance checks
	HighRiskTransactionThreshold = 10000.0        // High-risk transaction threshold in tokens
	SuspiciousActivityThreshold  = 3              // Number of suspicious activities before flagging a user
)

// AMLComplianceAutomation manages the AML compliance protocols by automating the monitoring and enforcement of transaction restrictions.
type AMLComplianceAutomation struct {
	consensusSystem    *consensus.SynnergyConsensus
	ledgerInstance     *ledger.Ledger
	stateMutex         *sync.RWMutex
	suspiciousActivity map[string]int // Tracking suspicious activities by user address
}

// NewAMLComplianceAutomation initializes and returns an instance of AMLComplianceAutomation with proper dependencies
func NewAMLComplianceAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AMLComplianceAutomation {
	return &AMLComplianceAutomation{
		consensusSystem:    consensusSystem,
		ledgerInstance:     ledgerInstance,
		stateMutex:         stateMutex,
		suspiciousActivity: make(map[string]int),
	}
}

// StartAMLMonitoring initiates continuous AML compliance checks that run at regular intervals
func (automation *AMLComplianceAutomation) StartAMLMonitoring() {
	ticker := time.NewTicker(AMLMonitoringInterval)

	go func() {
		for range ticker.C {
			automation.monitorTransactions()
		}
	}()
}

// monitorTransactions reviews all recent transactions to identify high-risk or suspicious activity
func (automation *AMLComplianceAutomation) monitorTransactions() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Retrieve the list of recent transactions
	recentTransactions := automation.consensusSystem.GetRecentTransactions()

	for _, tx := range recentTransactions {
		if tx.Amount > HighRiskTransactionThreshold {
			automation.handleHighRiskTransaction(tx)
		}
		if automation.isSuspicious(tx) {
			automation.trackSuspiciousActivity(tx.Sender)
		}
	}
}

// handleHighRiskTransaction processes a high-risk transaction, flags it, and records it in the ledger
func (automation *AMLComplianceAutomation) handleHighRiskTransaction(tx common.Transaction) {
	fmt.Printf("High-risk transaction flagged: Sender: %s, Amount: %.2f tokens.\n", tx.Sender, tx.Amount)

	// Encrypt transaction data for security before logging it
	encryptedTx := automation.encryptTransaction(tx)
	automation.logAMLViolation(encryptedTx, "High-Risk Transaction")
}

// isSuspicious determines whether a transaction exhibits suspicious activity based on predefined criteria
func (automation *AMLComplianceAutomation) isSuspicious(tx common.Transaction) bool {
	// Check if the user has exceeded the threshold for suspicious activity
	return tx.Amount > HighRiskTransactionThreshold && automation.suspiciousActivity[tx.Sender] >= SuspiciousActivityThreshold
}

// trackSuspiciousActivity increments the count of suspicious transactions by a user and logs if thresholds are exceeded
func (automation *AMLComplianceAutomation) trackSuspiciousActivity(sender string) {
	automation.suspiciousActivity[sender]++
	if automation.suspiciousActivity[sender] >= SuspiciousActivityThreshold {
		fmt.Printf("User %s flagged for repeated suspicious activities.\n", sender)
		automation.logAMLViolation(common.Transaction{Sender: sender}, "Suspicious Activity")
	}
}

// encryptTransaction securely encrypts transaction data before it is logged or flagged
func (automation *AMLComplianceAutomation) encryptTransaction(tx common.Transaction) common.Transaction {
	encryptedData, err := encryption.EncryptData([]byte(fmt.Sprintf("Sender: %s, Amount: %.2f tokens", tx.Sender, tx.Amount)))
	if err != nil {
		fmt.Println("Error encrypting transaction data:", err)
		return tx
	}
	tx.EncryptedData = encryptedData
	return tx
}

// logAMLViolation logs flagged transactions and suspicious activity in the ledger with a detailed entry
func (automation *AMLComplianceAutomation) logAMLViolation(tx common.Transaction, violationType string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("aml-violation-%s-%d", violationType, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "AML Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("AML violation (%s) for transaction: %s.", violationType, tx.EncryptedData),
	}

	automation.ledgerInstance.AddEntry(entry)
	fmt.Printf("AML violation logged: %s for transaction by %s.\n", violationType, tx.Sender)
}

