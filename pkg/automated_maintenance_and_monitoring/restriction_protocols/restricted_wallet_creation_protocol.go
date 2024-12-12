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
	WalletCreationCheckInterval = 10 * time.Second // Interval for checking restricted wallet creations
	MaxWalletCreationAttempts   = 5                // Maximum allowed unauthorized wallet creation attempts
)

// RestrictedWalletCreationAutomation enforces restrictions on unauthorized wallet creation attempts
type RestrictedWalletCreationAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	walletCreationAttempt map[string]int // Tracks unauthorized wallet creation attempts by user
}

// NewRestrictedWalletCreationAutomation initializes RestrictedWalletCreationAutomation
func NewRestrictedWalletCreationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RestrictedWalletCreationAutomation {
	return &RestrictedWalletCreationAutomation{
		consensusSystem:     consensusSystem,
		ledgerInstance:      ledgerInstance,
		stateMutex:          stateMutex,
		walletCreationAttempt: make(map[string]int),
	}
}

// StartWalletCreationMonitoring starts continuous monitoring of wallet creation attempts
func (automation *RestrictedWalletCreationAutomation) StartWalletCreationMonitoring() {
	ticker := time.NewTicker(WalletCreationCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorWalletCreations()
		}
	}()
}

// monitorWalletCreations checks for unauthorized wallet creation attempts and enforces restrictions
func (automation *RestrictedWalletCreationAutomation) monitorWalletCreations() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch wallet creation data from Synnergy Consensus
	walletCreationData := automation.consensusSystem.GetWalletCreationData()

	for userID, creationStatus := range walletCreationData {
		// Check if the wallet creation attempt is unauthorized
		if creationStatus == "unauthorized" {
			automation.flagWalletCreationViolation(userID, "Unauthorized wallet creation attempt detected")
		}
	}
}

// flagWalletCreationViolation flags an unauthorized wallet creation attempt and logs it in the ledger
func (automation *RestrictedWalletCreationAutomation) flagWalletCreationViolation(userID string, reason string) {
	fmt.Printf("Wallet creation violation: User ID %s, Reason: %s\n", userID, reason)

	// Increment the violation count for the user
	automation.walletCreationAttempt[userID]++

	// Log the violation in the ledger
	automation.logWalletCreationViolation(userID, reason)

	// Check if the user has exceeded the allowed number of wallet creation violations
	if automation.walletCreationAttempt[userID] >= MaxWalletCreationAttempts {
		automation.restrictWalletCreation(userID)
	}
}

// logWalletCreationViolation logs the flagged wallet creation violation into the ledger with details
func (automation *RestrictedWalletCreationAutomation) logWalletCreationViolation(userID string, violationReason string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-creation-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Wallet Creation Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated wallet creation restrictions. Reason: %s", userID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptWalletCreationData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log wallet creation violation:", err)
	} else {
		fmt.Println("Wallet creation violation logged.")
	}
}

// restrictWalletCreation restricts wallet creation for a user after exceeding violations
func (automation *RestrictedWalletCreationAutomation) restrictWalletCreation(userID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-creation-restriction-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Wallet Creation Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("User %s has been restricted from creating wallets due to repeated violations.", userID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptWalletCreationData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log wallet creation restriction:", err)
	} else {
		fmt.Println("Wallet creation restriction applied.")
	}
}

// encryptWalletCreationData encrypts the wallet creation data before logging for security
func (automation *RestrictedWalletCreationAutomation) encryptWalletCreationData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting wallet creation data:", err)
		return data
	}
	return string(encryptedData)
}
