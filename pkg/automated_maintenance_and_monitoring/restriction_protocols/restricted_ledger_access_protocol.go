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
	LedgerAccessCheckInterval    = 10 * time.Second // Interval for checking restricted ledger access
	MaxAllowedAccessAttempts     = 5                // Maximum allowed unauthorized access attempts before restriction
)

// RestrictedLedgerAccessAutomation enforces restrictions on unauthorized ledger access
type RestrictedLedgerAccessAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	accessViolationCount   map[string]int // Tracks unauthorized ledger access attempts per user or entity
}

// NewRestrictedLedgerAccessAutomation initializes and returns an instance of RestrictedLedgerAccessAutomation
func NewRestrictedLedgerAccessAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RestrictedLedgerAccessAutomation {
	return &RestrictedLedgerAccessAutomation{
		consensusSystem:      consensusSystem,
		ledgerInstance:       ledgerInstance,
		stateMutex:           stateMutex,
		accessViolationCount: make(map[string]int),
	}
}

// StartLedgerAccessMonitoring starts continuous monitoring of ledger access violations
func (automation *RestrictedLedgerAccessAutomation) StartLedgerAccessMonitoring() {
	ticker := time.NewTicker(LedgerAccessCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorLedgerAccess()
		}
	}()
}

// monitorLedgerAccess checks for unauthorized ledger access attempts and enforces restrictions if necessary
func (automation *RestrictedLedgerAccessAutomation) monitorLedgerAccess() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch ledger access data from Synnergy Consensus
	accessData := automation.consensusSystem.GetLedgerAccessData()

	for userID, accessStatus := range accessData {
		// Check if the access attempt is unauthorized
		if accessStatus == "unauthorized" {
			automation.flagAccessViolation(userID, "Unauthorized ledger access attempt detected")
		}
	}
}

// flagAccessViolation flags a user's unauthorized ledger access attempt and logs it in the ledger
func (automation *RestrictedLedgerAccessAutomation) flagAccessViolation(userID string, reason string) {
	fmt.Printf("Ledger access violation: User ID %s, Reason: %s\n", userID, reason)

	// Increment the violation count for the user
	automation.accessViolationCount[userID]++

	// Log the violation in the ledger
	automation.logAccessViolation(userID, reason)

	// Check if the user has exceeded the allowed number of access violations
	if automation.accessViolationCount[userID] >= MaxAllowedAccessAttempts {
		automation.restrictLedgerAccess(userID)
	}
}

// logAccessViolation logs the flagged ledger access violation into the ledger with full details
func (automation *RestrictedLedgerAccessAutomation) logAccessViolation(userID string, violationReason string) {
	// Create a ledger entry for ledger access violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("ledger-access-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Ledger Access Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated ledger access restrictions. Reason: %s", userID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptAccessData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log ledger access violation:", err)
	} else {
		fmt.Println("Ledger access violation logged.")
	}
}

// restrictLedgerAccess restricts ledger access for a user after exceeding allowed violations
func (automation *RestrictedLedgerAccessAutomation) restrictLedgerAccess(userID string) {
	// Add restriction details to the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("ledger-access-restriction-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Ledger Access Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("User %s has been restricted from ledger access due to repeated violations.", userID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptAccessData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log ledger access restriction:", err)
	} else {
		fmt.Println("Ledger access restriction applied.")
	}
}

// encryptAccessData encrypts the ledger access data before logging for security
func (automation *RestrictedLedgerAccessAutomation) encryptAccessData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting ledger access data:", err)
		return data
	}
	return string(encryptedData)
}
