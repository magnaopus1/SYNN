package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/zkproof"
)

// Constants for ZK Proof submission limits
const (
	ZKProofSubmissionInterval = 3 * time.Second // Interval to check ZK Proof submission requests
	MaxZKProofSubmissions     = 50              // Max allowed ZK Proof submissions per day per wallet
)

// ZKProofSubmissionAutomation monitors and restricts ZK Proof submissions
type ZKProofSubmissionAutomation struct {
	consensusSystem  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	transactionMutex *sync.RWMutex
}

// NewZKProofSubmissionAutomation initializes the automation for ZK Proof submission restriction
func NewZKProofSubmissionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, transactionMutex *sync.RWMutex) *ZKProofSubmissionAutomation {
	return &ZKProofSubmissionAutomation{
		consensusSystem:  consensusSystem,
		ledgerInstance:   ledgerInstance,
		transactionMutex: transactionMutex,
	}
}

// StartMonitoring continuously monitors ZK Proof submission limits for wallets
func (automation *ZKProofSubmissionAutomation) StartMonitoring() {
	ticker := time.NewTicker(ZKProofSubmissionInterval)

	go func() {
		for range ticker.C {
			automation.evaluateZKProofSubmissionLimits()
		}
	}()
}

// evaluateZKProofSubmissionLimits checks if wallets have exceeded the ZK Proof submission limits
func (automation *ZKProofSubmissionAutomation) evaluateZKProofSubmissionLimits() {
	automation.transactionMutex.Lock()
	defer automation.transactionMutex.Unlock()

	wallets := automation.consensusSystem.GetWallets()
	for _, wallet := range wallets {
		if !automation.isWithinZKProofLimit(wallet) {
			automation.enforceLimit(wallet)
		}
	}
}

// isWithinZKProofLimit checks if a wallet is within the allowed ZK Proof submission limits
func (automation *ZKProofSubmissionAutomation) isWithinZKProofLimit(wallet consensus.Wallet) bool {
	if wallet.DailyZKProofSubmissions > MaxZKProofSubmissions {
		automation.logLimitBreach(wallet.ID, "Exceeded ZK Proof submission limit", wallet.DailyZKProofSubmissions)
		return false
	}
	return true
}

// enforceLimit takes actions when a wallet exceeds the ZK Proof submission limit
func (automation *ZKProofSubmissionAutomation) enforceLimit(wallet consensus.Wallet) {
	// Implement logic to freeze or restrict further ZK Proof submissions
	err := automation.consensusSystem.RestrictZKProofSubmissions(wallet.ID)
	if err != nil {
		fmt.Println("Failed to enforce ZK Proof submission limit for wallet:", wallet.ID, "Error:", err)
		return
	}

	// Log the enforcement in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("zkproof-limit-enforcement-%s-%d", wallet.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "ZK Proof Submission Limit Enforcement",
		Status:    "Enforced",
		Details:   fmt.Sprintf("ZK Proof submission limits enforced for wallet %s due to exceeding limits.", wallet.ID),
	}

	// Encrypt the enforcement details
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log ZK Proof submission limit enforcement:", err)
	} else {
		fmt.Println("ZK Proof submission limit enforcement logged for wallet:", wallet.ID)
	}
}

// logLimitBreach logs a breach of the ZK Proof submission limit into the ledger
func (automation *ZKProofSubmissionAutomation) logLimitBreach(walletID string, reason string, value interface{}) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("zkproof-limit-breach-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "ZK Proof Limit Breach",
		Status:    "Breach",
		Details:   fmt.Sprintf("Wallet %s breached ZK Proof limits due to %s: %v.", walletID, reason, value),
	}

	// Encrypt details before storing in the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log ZK Proof submission limit breach:", err)
	} else {
		fmt.Println("ZK Proof submission limit breach logged for wallet:", walletID)
	}
}

// encryptData encrypts sensitive information before storing it in the ledger
func (automation *ZKProofSubmissionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
