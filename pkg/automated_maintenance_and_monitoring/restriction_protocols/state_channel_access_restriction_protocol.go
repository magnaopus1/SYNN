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

// Constants for State Channel access checks
const (
	StateChannelCheckInterval       = 5 * time.Second // Check interval
	UnauthorizedAccessWarning       = "Unauthorized access detected for state channel"
	ExcessiveAccessAttemptWarning   = "Excessive access attempts detected"
	MaxAccessAttempts               = 5 // Maximum access attempts before restriction
)

// StateChannelAccessRestrictionAutomation manages access restrictions to state channels
type StateChannelAccessRestrictionAutomation struct {
	consensusSystem    *consensus.SynnergyConsensus
	ledgerInstance     *ledger.Ledger
	stateMutex         *sync.RWMutex
	accessAttempts     map[string]int    // Track access attempts per wallet
	restrictedWallets  map[string]bool   // Track restricted wallets
}

// NewStateChannelAccessRestrictionAutomation initializes the state channel restriction automation
func NewStateChannelAccessRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *StateChannelAccessRestrictionAutomation {
	return &StateChannelAccessRestrictionAutomation{
		consensusSystem:   consensusSystem,
		ledgerInstance:    ledgerInstance,
		stateMutex:        stateMutex,
		accessAttempts:    make(map[string]int),
		restrictedWallets: make(map[string]bool),
	}
}

// StartAccessMonitoring begins continuous monitoring of state channel access attempts
func (automation *StateChannelAccessRestrictionAutomation) StartAccessMonitoring() {
	ticker := time.NewTicker(StateChannelCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateStateChannelAccess()
		}
	}()
}

// evaluateStateChannelAccess monitors access to state channels and enforces restrictions
func (automation *StateChannelAccessRestrictionAutomation) evaluateStateChannelAccess() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	accessData := automation.consensusSystem.GetStateChannelAccessData()

	for walletID, accessCount := range accessData {
		if accessCount > MaxAccessAttempts {
			automation.logExcessiveAccessAttempts(walletID, accessCount)
			automation.restrictWallet(walletID)
		} else if automation.consensusSystem.IsUnauthorizedStateChannelAccess(walletID) {
			automation.logUnauthorizedAccess(walletID)
			automation.restrictWallet(walletID)
		}
	}
}

// logExcessiveAccessAttempts logs instances of excessive access attempts in the ledger
func (automation *StateChannelAccessRestrictionAutomation) logExcessiveAccessAttempts(walletID string, accessCount int) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("excessive-access-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Excessive State Channel Access",
		Status:    "Warning",
		Details:   fmt.Sprintf("Wallet %s attempted to access state channel %d times, exceeding the limit.", walletID, accessCount),
	}

	// Encrypt the details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log excessive access attempts:", err)
	} else {
		fmt.Println("Excessive state channel access logged for wallet:", walletID)
	}
}

// logUnauthorizedAccess logs unauthorized state channel access in the ledger
func (automation *StateChannelAccessRestrictionAutomation) logUnauthorizedAccess(walletID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("unauthorized-access-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Unauthorized State Channel Access",
		Status:    "Critical",
		Details:   fmt.Sprintf("Unauthorized access to a state channel detected for wallet %s.", walletID),
	}

	// Encrypt the unauthorized access details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log unauthorized access:", err)
	} else {
		fmt.Println("Unauthorized access logged for wallet:", walletID)
	}
}

// restrictWallet restricts access for wallets that breach state channel access policies
func (automation *StateChannelAccessRestrictionAutomation) restrictWallet(walletID string) {
	fmt.Printf("Wallet %s has been restricted due to state channel access violations.\n", walletID)

	// Log the restriction in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-restriction-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "State Channel Access Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Wallet %s restricted from accessing state channels due to excessive attempts or unauthorized activity.", walletID),
	}

	// Encrypt the restriction details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log wallet restriction:", err)
	} else {
		fmt.Println("Wallet restriction logged for:", walletID)
	}

	// Update the consensus system to restrict access for the wallet
	automation.consensusSystem.RestrictStateChannelAccess(walletID)
}

// encryptData encrypts sensitive information before storing in the ledger
func (automation *StateChannelAccessRestrictionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting access data:", err)
		return data
	}
	return string(encryptedData)
}

