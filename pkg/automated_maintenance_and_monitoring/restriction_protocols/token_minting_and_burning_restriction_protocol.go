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
	MintBurnCheckInterval       = 5 * time.Second // Interval for checking minting and burning events
	MaxMintBurnTransactionLimit = 1000            // Example max mint/burn transaction per block
)

// TokenMintingAndBurningRestrictionAutomation handles minting and burning restrictions on the network
type TokenMintingAndBurningRestrictionAutomation struct {
	consensusSystem   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	transactionMutex  *sync.RWMutex
	mintBurnTracking  map[string]int  // Tracks token minting/burning per wallet
	restrictedWallets map[string]bool // Tracks wallets restricted from minting/burning
}

// NewTokenMintingAndBurningRestrictionAutomation initializes the automation for token minting and burning restriction
func NewTokenMintingAndBurningRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, transactionMutex *sync.RWMutex) *TokenMintingAndBurningRestrictionAutomation {
	return &TokenMintingAndBurningRestrictionAutomation{
		consensusSystem:   consensusSystem,
		ledgerInstance:    ledgerInstance,
		transactionMutex:  transactionMutex,
		mintBurnTracking:  make(map[string]int),
		restrictedWallets: make(map[string]bool),
	}
}

// StartMonitoring starts continuous monitoring of token minting and burning activities
func (automation *TokenMintingAndBurningRestrictionAutomation) StartMonitoring() {
	ticker := time.NewTicker(MintBurnCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateMintBurnTransactions()
		}
	}()
}

// evaluateMintBurnTransactions checks for excessive minting or burning transactions and restricts if necessary
func (automation *TokenMintingAndBurningRestrictionAutomation) evaluateMintBurnTransactions() {
	automation.transactionMutex.Lock()
	defer automation.transactionMutex.Unlock()

	// Get minting and burning transactions data from the consensus system
	mintBurnData := automation.consensusSystem.GetMintBurnTransactionData()

	for walletID, transactionCount := range mintBurnData {
		if transactionCount > MaxMintBurnTransactionLimit {
			automation.logExcessiveMintBurn(walletID, transactionCount)
			automation.restrictMintBurn(walletID)
		}
	}
}

// logExcessiveMintBurn logs excessive minting or burning attempts in the ledger
func (automation *TokenMintingAndBurningRestrictionAutomation) logExcessiveMintBurn(walletID string, transactionCount int) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("excessive-mint-burn-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Excessive Token Minting/Burning",
		Status:    "Critical",
		Details:   fmt.Sprintf("Wallet %s exceeded minting/burning limit with %d transactions.", walletID, transactionCount),
	}

	// Encrypt the details before storing in the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log excessive minting/burning:", err)
	} else {
		fmt.Println("Excessive minting/burning logged for wallet:", walletID)
	}
}

// restrictMintBurn restricts wallets from further minting or burning due to excessive activity
func (automation *TokenMintingAndBurningRestrictionAutomation) restrictMintBurn(walletID string) {
	fmt.Printf("Minting/Burning restricted for wallet %s due to exceeding limits.\n", walletID)

	// Log the restriction in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-restriction-mint-burn-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Token Minting/Burning Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Wallet %s restricted from minting/burning tokens due to policy violations.", walletID),
	}

	// Encrypt the restriction details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log mint/burn restriction:", err)
	} else {
		fmt.Println("Mint/Burn restriction logged for:", walletID)
	}

	// Update the consensus system to restrict minting/burning for the wallet
	automation.consensusSystem.RestrictMintBurnAccess(walletID)
}

// encryptData encrypts sensitive information before storing in the ledger
func (automation *TokenMintingAndBurningRestrictionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting mint/burn data:", err)
		return data
	}
	return string(encryptedData)
}

