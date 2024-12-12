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
	WalletRecoveryCheckInterval = 5 * time.Second  // Interval for checking wallet recovery attempts
	MaxWalletRecoveryAttempts   = 3                // Maximum allowed unauthorized wallet recovery attempts
)

// RestrictedWalletRecoveryAutomation enforces restrictions on unauthorized wallet recovery attempts
type RestrictedWalletRecoveryAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	walletRecoveryAttempts map[string]int // Tracks unauthorized wallet recovery attempts by user
}

// NewRestrictedWalletRecoveryAutomation initializes RestrictedWalletRecoveryAutomation
func NewRestrictedWalletRecoveryAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RestrictedWalletRecoveryAutomation {
	return &RestrictedWalletRecoveryAutomation{
		consensusSystem:        consensusSystem,
		ledgerInstance:         ledgerInstance,
		stateMutex:             stateMutex,
		walletRecoveryAttempts: make(map[string]int),
	}
}

// StartWalletRecoveryMonitoring starts continuous monitoring of wallet recovery attempts
func (automation *RestrictedWalletRecoveryAutomation) StartWalletRecoveryMonitoring() {
	ticker := time.NewTicker(WalletRecoveryCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorWalletRecoveryAttempts()
		}
	}()
}

// monitorWalletRecoveryAttempts checks for unauthorized wallet recovery attempts and enforces restrictions
func (automation *RestrictedWalletRecoveryAutomation) monitorWalletRecoveryAttempts() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch wallet recovery attempt data from Synnergy Consensus
	recoveryData := automation.consensusSystem.GetWalletRecoveryData()

	for userID, recoveryStatus := range recoveryData {
		// Check if the wallet recovery attempt is unauthorized
		if recoveryStatus == "unauthorized" {
			automation.flagWalletRecoveryViolation(userID, "Unauthorized wallet recovery attempt detected")
		}
	}
}

// flagWalletRecoveryViolation flags an unauthorized wallet recovery attempt and logs it in the ledger
func (automation *RestrictedWalletRecoveryAutomation) flagWalletRecoveryViolation(userID string, reason string) {
	fmt.Printf("Wallet recovery violation: User ID %s, Reason: %s\n", userID, reason)

	// Increment the violation count for the user
	automation.walletRecoveryAttempts[userID]++

	// Log the violation in the ledger
	automation.logWalletRecoveryViolation(userID, reason)

	// Check if the user has exceeded the allowed number of wallet recovery violations
	if automation.walletRecoveryAttempts[userID] >= MaxWalletRecoveryAttempts {
		automation.restrictWalletRecovery(userID)
	}
}

// logWalletRecoveryViolation logs the flagged wallet recovery violation into the ledger with details
func (automation *RestrictedWalletRecoveryAutomation) logWalletRecoveryViolation(userID string, violationReason string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-recovery-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Wallet Recovery Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated wallet recovery restrictions. Reason: %s", userID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptWalletRecoveryData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log wallet recovery violation:", err)
	} else {
		fmt.Println("Wallet recovery violation logged.")
	}
}

// restrictWalletRecovery restricts wallet recovery for a user after exceeding violations
func (automation *RestrictedWalletRecoveryAutomation) restrictWalletRecovery(userID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-recovery-restriction-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Wallet Recovery Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("User %s has been restricted from recovering wallets due to repeated violations.", userID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptWalletRecoveryData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log wallet recovery restriction:", err)
	} else {
		fmt.Println("Wallet recovery restriction applied.")
	}
}

// encryptWalletRecoveryData encrypts the wallet recovery data before logging for security
func (automation *RestrictedWalletRecoveryAutomation) encryptWalletRecoveryData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting wallet recovery data:", err)
		return data
	}
	return string(encryptedData)
}
