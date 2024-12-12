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

// Constants for transaction fee ceiling
const (
	TransactionFeeCheckInterval = 3 * time.Second  // Interval for checking transaction fees
	TransactionFeeCeiling       = 0.0025           // The fee ceiling, e.g., 5%
)

// TransactionFeeCeilingAutomation ensures no transaction exceeds the fee ceiling
type TransactionFeeCeilingAutomation struct {
	consensusSystem  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	transactionMutex *sync.RWMutex
}

// NewTransactionFeeCeilingAutomation initializes the automation for monitoring transaction fees
func NewTransactionFeeCeilingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, transactionMutex *sync.RWMutex) *TransactionFeeCeilingAutomation {
	return &TransactionFeeCeilingAutomation{
		consensusSystem:  consensusSystem,
		ledgerInstance:   ledgerInstance,
		transactionMutex: transactionMutex,
	}
}

// StartMonitoring begins continuous monitoring of transaction fees
func (automation *TransactionFeeCeilingAutomation) StartMonitoring() {
	ticker := time.NewTicker(TransactionFeeCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateTransactionFees()
		}
	}()
}

// evaluateTransactionFees checks if any transactions exceed the fee ceiling
func (automation *TransactionFeeCeilingAutomation) evaluateTransactionFees() {
	automation.transactionMutex.Lock()
	defer automation.transactionMutex.Unlock()

	// Fetch transaction fee data from the consensus system
	transactions := automation.consensusSystem.GetPendingTransactions()

	for _, tx := range transactions {
		if tx.Fee > TransactionFeeCeiling {
			automation.logTransactionFeeViolation(tx.ID, tx.Fee)
			automation.applyFeeReduction(tx.ID)
		}
	}
}

// logTransactionFeeViolation logs transactions that exceed the fee ceiling in the ledger
func (automation *TransactionFeeCeilingAutomation) logTransactionFeeViolation(txID string, fee float64) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("fee-violation-%s-%d", txID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Transaction Fee Violation",
		Status:    "Violation",
		Details:   fmt.Sprintf("Transaction %s exceeded the fee ceiling with a fee of %.2f%%.", txID, fee),
	}

	// Encrypt the violation details before storing in the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log transaction fee violation:", err)
	} else {
		fmt.Println("Transaction fee violation logged for transaction:", txID)
	}
}

// applyFeeReduction reduces the fee for transactions exceeding the ceiling and logs the adjustment
func (automation *TransactionFeeCeilingAutomation) applyFeeReduction(txID string) {
	fmt.Printf("Applying fee reduction for transaction %s to adhere to the fee ceiling.\n", txID)

	// Adjust the fee to meet the ceiling in the consensus system
	err := automation.consensusSystem.AdjustTransactionFee(txID, TransactionFeeCeiling)
	if err != nil {
		fmt.Println("Failed to adjust transaction fee for:", txID, "Error:", err)
		return
	}

	// Log the fee adjustment in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("fee-adjustment-%s-%d", txID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Transaction Fee Adjustment",
		Status:    "Adjusted",
		Details:   fmt.Sprintf("Transaction %s fee reduced to %.2f%% due to exceeding the fee ceiling.", txID, TransactionFeeCeiling),
	}

	// Encrypt the adjustment details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log fee adjustment:", err)
	} else {
		fmt.Println("Fee adjustment logged for transaction:", txID)
	}
}

// encryptData encrypts sensitive information before storing in the ledger
func (automation *TransactionFeeCeilingAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
