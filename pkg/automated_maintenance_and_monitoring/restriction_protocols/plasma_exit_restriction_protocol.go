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
	PlasmaExitCheckInterval   = 20 * time.Second // Interval for checking Plasma exit requests
	MaxAllowedPendingExits    = 10               // Maximum number of pending Plasma exits allowed per user
)

// PlasmaExitRestrictionAutomation monitors and restricts Plasma exit requests across the network
type PlasmaExitRestrictionAutomation struct {
	consensusSystem      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	stateMutex           *sync.RWMutex
	pendingExitCount     map[string]int // Tracks pending Plasma exit requests per user
}

// NewPlasmaExitRestrictionAutomation initializes and returns an instance of PlasmaExitRestrictionAutomation
func NewPlasmaExitRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *PlasmaExitRestrictionAutomation {
	return &PlasmaExitRestrictionAutomation{
		consensusSystem:  consensusSystem,
		ledgerInstance:   ledgerInstance,
		stateMutex:       stateMutex,
		pendingExitCount: make(map[string]int),
	}
}

// StartExitMonitoring starts continuous monitoring of Plasma exit requests
func (automation *PlasmaExitRestrictionAutomation) StartExitMonitoring() {
	ticker := time.NewTicker(PlasmaExitCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorPlasmaExits()
		}
	}()
}

// monitorPlasmaExits checks for pending Plasma exit requests and enforces restrictions if necessary
func (automation *PlasmaExitRestrictionAutomation) monitorPlasmaExits() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch pending Plasma exit data from Synnergy Consensus
	exitData := automation.consensusSystem.GetPendingPlasmaExitData()

	for userID, exitCount := range exitData {
		// Check if the user has exceeded the allowed number of pending Plasma exits
		if automation.pendingExitCount[userID] > MaxAllowedPendingExits {
			automation.flagExitViolation(userID, exitCount, "Exceeded allowed number of pending Plasma exits")
		}
	}
}

// flagExitViolation flags a user's Plasma exit violation and logs it in the ledger
func (automation *PlasmaExitRestrictionAutomation) flagExitViolation(userID string, exitCount int, reason string) {
	fmt.Printf("Plasma exit violation: User ID %s, Reason: %s\n", userID, reason)

	// Log the violation in the ledger
	automation.logExitViolation(userID, exitCount, reason)
}

// logExitViolation logs the flagged Plasma exit violation into the ledger with full details
func (automation *PlasmaExitRestrictionAutomation) logExitViolation(userID string, exitCount int, violationReason string) {
	// Create a ledger entry for Plasma exit violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("plasma-exit-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Plasma Exit Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated Plasma exit rules. Exit Requests: %d. Reason: %s", userID, exitCount, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptExitData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log Plasma exit violation:", err)
	} else {
		fmt.Println("Plasma exit violation logged.")
	}
}

// encryptExitData encrypts the Plasma exit data before logging for security
func (automation *PlasmaExitRestrictionAutomation) encryptExitData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting Plasma exit data:", err)
		return data
	}
	return string(encryptedData)
}
