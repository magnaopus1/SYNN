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

// Constants
const (
	TokenTransferCheckInterval       = 3 * time.Second  // Interval for checking token transfers
	MaxTokenTransferLimitPerBlock    = 10000            // Maximum tokens allowed to be transferred per block
	RestrictedTransferTransactionFee = 0.02             // Increased transaction fee for restricted wallets (2%)
)

// TokenTransferRestrictionAutomation handles restrictions and checks for token transfers in the network
type TokenTransferRestrictionAutomation struct {
	consensusSystem   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	transactionMutex  *sync.RWMutex
	transferTracking  map[string]int   // Tracks token transfer activity per wallet
	restrictedWallets map[string]bool  // Tracks restricted wallets for token transfer
}

// NewTokenTransferRestrictionAutomation initializes automation for token transfer restriction
func NewTokenTransferRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, transactionMutex *sync.RWMutex) *TokenTransferRestrictionAutomation {
	return &TokenTransferRestrictionAutomation{
		consensusSystem:   consensusSystem,
		ledgerInstance:    ledgerInstance,
		transactionMutex:  transactionMutex,
		transferTracking:  make(map[string]int),
		restrictedWallets: make(map[string]bool),
	}
}

// StartMonitoring begins continuous monitoring of token transfers
func (automation *TokenTransferRestrictionAutomation) StartMonitoring() {
	ticker := time.NewTicker(TokenTransferCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateTokenTransfers()
		}
	}()
}

// evaluateTokenTransfers checks for excessive or suspicious token transfer activity
func (automation *TokenTransferRestrictionAutomation) evaluateTokenTransfers() {
	automation.transactionMutex.Lock()
	defer automation.transactionMutex.Unlock()

	// Get token transfer data from the consensus system
	transferData := automation.consensusSystem.GetTokenTransferData()

	for walletID, transferAmount := range transferData {
		if transferAmount > MaxTokenTransferLimitPerBlock {
			automation.logExcessiveTransfer(walletID, transferAmount)
			automation.restrictTokenTransfers(walletID)
		} else if automation.isWalletRestricted(walletID) {
			automation.applyRestrictedTransferFee(walletID)
		}
	}
}

// logExcessiveTransfer logs excessive token transfers in the ledger
func (automation *TokenTransferRestrictionAutomation) logExcessiveTransfer(walletID string, transferAmount int) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("excessive-token-transfer-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Excessive Token Transfer",
		Status:    "Critical",
		Details:   fmt.Sprintf("Wallet %s exceeded token transfer limit with %d tokens transferred.", walletID, transferAmount),
	}

	// Encrypt the details before storing in the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log excessive token transfer:", err)
	} else {
		fmt.Println("Excessive token transfer logged for wallet:", walletID)
	}
}

// restrictTokenTransfers restricts wallets that have exceeded transfer limits
func (automation *TokenTransferRestrictionAutomation) restrictTokenTransfers(walletID string) {
	fmt.Printf("Token transfer restricted for wallet %s due to exceeding transfer limits.\n", walletID)

	// Log the restriction in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-restriction-transfer-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Token Transfer Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Wallet %s restricted from token transfers due to policy violations.", walletID),
	}

	// Encrypt the restriction details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log transfer restriction:", err)
	} else {
		fmt.Println("Transfer restriction logged for:", walletID)
	}

	// Update the consensus system to restrict transfers for the wallet
	automation.consensusSystem.RestrictTokenTransfers(walletID)
	automation.restrictedWallets[walletID] = true
}

// isWalletRestricted checks if a wallet is restricted from transferring tokens
func (automation *TokenTransferRestrictionAutomation) isWalletRestricted(walletID string) bool {
	_, restricted := automation.restrictedWallets[walletID]
	return restricted
}

// applyRestrictedTransferFee applies a higher transaction fee for restricted wallets
func (automation *TokenTransferRestrictionAutomation) applyRestrictedTransferFee(walletID string) {
	fmt.Printf("Applying increased transaction fee for restricted wallet: %s\n", walletID)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-restricted-transfer-fee-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Increased Transaction Fee",
		Status:    "Applied",
		Details:   fmt.Sprintf("Wallet %s charged with a restricted transaction fee due to violation.", walletID),
	}

	// Encrypt the fee application details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log restricted transfer fee:", err)
	} else {
		fmt.Println("Restricted transfer fee applied for wallet:", walletID)
	}

	// Implement fee change in consensus
	automation.consensusSystem.ApplyIncreasedTransactionFee(walletID, RestrictedTransferTransactionFee)
}

// encryptData encrypts sensitive information before storing in the ledger
func (automation *TokenTransferRestrictionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting transfer data:", err)
		return data
	}
	return string(encryptedData)
}
