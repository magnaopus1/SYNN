package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
)

// Constants for wallet transaction limits
const (
	TransactionCheckInterval = 5 * time.Second // Interval to check wallet transaction limits
	MaxTransactionsPerDay    = 100000             // Maximum allowed transactions per wallet per day
	MaxTransactionVolume     = 500000000000000            // Maximum allowed transaction volume per day (in tokens)
)

// WalletTransactionLimitAutomation monitors wallet transaction limits
type WalletTransactionLimitAutomation struct {
	consensusSystem  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	transactionMutex *sync.RWMutex
}

// NewWalletTransactionLimitAutomation initializes the automation for wallet transaction limits
func NewWalletTransactionLimitAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, transactionMutex *sync.RWMutex) *WalletTransactionLimitAutomation {
	return &WalletTransactionLimitAutomation{
		consensusSystem:  consensusSystem,
		ledgerInstance:   ledgerInstance,
		transactionMutex: transactionMutex,
	}
}

// StartMonitoring continuously monitors transaction limits for all wallets
func (automation *WalletTransactionLimitAutomation) StartMonitoring() {
	ticker := time.NewTicker(TransactionCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateTransactionLimits()
		}
	}()
}

// evaluateTransactionLimits checks the transaction limits for all wallets in the network
func (automation *WalletTransactionLimitAutomation) evaluateTransactionLimits() {
	automation.transactionMutex.Lock()
	defer automation.transactionMutex.Unlock()

	wallets := automation.consensusSystem.GetWallets()
	for _, wallet := range wallets {
		if !automation.isWithinLimit(wallet) {
			automation.enforceLimit(wallet)
		}
	}
}

// isWithinLimit checks if a wallet is within the defined transaction limits
func (automation *WalletTransactionLimitAutomation) isWithinLimit(wallet consensus.Wallet) bool {
	// Check if the wallet's total transactions exceed the daily limit
	if wallet.DailyTransactions > MaxTransactionsPerDay {
		automation.logLimitBreach(wallet.ID, "Exceeded daily transaction count", wallet.DailyTransactions)
		return false
	}

	// Check if the wallet's total transaction volume exceeds the daily limit
	if wallet.DailyTransactionVolume > MaxTransactionVolume {
		automation.logLimitBreach(wallet.ID, "Exceeded daily transaction volume", wallet.DailyTransactionVolume)
		return false
	}

	return true
}

// enforceLimit restricts a wallet that exceeds the transaction limit
func (automation *WalletTransactionLimitAutomation) enforceLimit(wallet consensus.Wallet) {
	// Implement logic to freeze wallet or reject transactions exceeding limits
	err := automation.consensusSystem.FreezeWallet(wallet.ID)
	if err != nil {
		fmt.Println("Failed to enforce transaction limit for wallet:", wallet.ID, "Error:", err)
		return
	}

	// Log the enforcement in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-limit-enforcement-%s-%d", wallet.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Wallet Transaction Limit Enforcement",
		Status:    "Enforced",
		Details:   fmt.Sprintf("Transaction limits enforced for wallet %s due to exceeding limits.", wallet.ID),
	}

	// Encrypt the enforcement details
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log transaction limit enforcement:", err)
	} else {
		fmt.Println("Transaction limit enforcement logged for wallet:", wallet.ID)
	}
}

// logLimitBreach logs a breach of transaction limits into the ledger
func (automation *WalletTransactionLimitAutomation) logLimitBreach(walletID string, reason string, value interface{}) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("limit-breach-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Transaction Limit Breach",
		Status:    "Breach",
		Details:   fmt.Sprintf("Wallet %s breached limits due to %s: %v.", walletID, reason, value),
	}

	// Encrypt details before storing in the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log limit breach:", err)
	} else {
		fmt.Println("Transaction limit breach logged for wallet:", walletID)
	}
}

// encryptData encrypts sensitive information before storing it in the ledger
func (automation *WalletTransactionLimitAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
