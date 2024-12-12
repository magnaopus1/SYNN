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
	TestnetAccessCheckInterval  = 5 * time.Second // Interval for checking access attempts
	MaxUnauthorizedAttempts     = 2
	UnauthorizedTestnetAccess   = "Unauthorized testnet access detected"
	RestrictedAccessLogMessage  = "Testnet access has been restricted for wallet"
)

// TestnetAccessRestrictionAutomation handles the restriction and monitoring of testnet access.
type TestnetAccessRestrictionAutomation struct {
	consensusSystem   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	accessMutex       *sync.RWMutex
	accessAttempts    map[string]int  // Tracks the access attempts for each wallet
	restrictedWallets map[string]bool // Tracks restricted wallets
}

// NewTestnetAccessRestrictionAutomation initializes the automation for testnet access restriction
func NewTestnetAccessRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, accessMutex *sync.RWMutex) *TestnetAccessRestrictionAutomation {
	return &TestnetAccessRestrictionAutomation{
		consensusSystem:   consensusSystem,
		ledgerInstance:    ledgerInstance,
		accessMutex:       accessMutex,
		accessAttempts:    make(map[string]int),
		restrictedWallets: make(map[string]bool),
	}
}

// StartAccessMonitoring starts continuous monitoring of testnet access privileges
func (automation *TestnetAccessRestrictionAutomation) StartAccessMonitoring() {
	ticker := time.NewTicker(TestnetAccessCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateTestnetAccess()
		}
	}()
}

// evaluateTestnetAccess checks for unauthorized testnet access and enforces restrictions if necessary
func (automation *TestnetAccessRestrictionAutomation) evaluateTestnetAccess() {
	automation.accessMutex.Lock()
	defer automation.accessMutex.Unlock()

	accessData := automation.consensusSystem.GetTestnetAccessData()

	for walletID, accessCount := range accessData {
		if accessCount > MaxUnauthorizedAttempts {
			automation.logUnauthorizedAccess(walletID)
			automation.restrictTestnetAccess(walletID)
		}
	}
}

// logUnauthorizedAccess logs unauthorized testnet access attempts in the ledger
func (automation *TestnetAccessRestrictionAutomation) logUnauthorizedAccess(walletID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("unauthorized-testnet-access-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Unauthorized Testnet Access",
		Status:    "Critical",
		Details:   fmt.Sprintf("Wallet %s made unauthorized testnet access attempts.", walletID),
	}

	// Encrypt the details before storing in the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log unauthorized testnet access:", err)
	} else {
		fmt.Println("Unauthorized testnet access logged for wallet:", walletID)
	}
}

// restrictTestnetAccess restricts testnet access for wallets that breach access policies
func (automation *TestnetAccessRestrictionAutomation) restrictTestnetAccess(walletID string) {
	fmt.Printf("Testnet access restricted for wallet %s due to unauthorized attempts.\n", walletID)

	// Log the restriction in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-restriction-testnet-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Testnet Access Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Wallet %s restricted from testnet due to policy violations.", walletID),
	}

	// Encrypt the restriction details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log testnet access restriction:", err)
	} else {
		fmt.Println("Testnet access restriction logged for:", walletID)
	}

	// Update the consensus system to restrict testnet access for the wallet
	automation.consensusSystem.RestrictTestnetAccess(walletID)
}

// encryptData encrypts sensitive information before storing in the ledger
func (automation *TestnetAccessRestrictionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting testnet access data:", err)
		return data
	}
	return string(encryptedData)
}

