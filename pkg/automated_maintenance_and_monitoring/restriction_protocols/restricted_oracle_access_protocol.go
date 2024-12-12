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
	OracleAccessCheckInterval = 10 * time.Second // Interval for checking restricted oracle access
	MaxOracleUnauthorizedAttempts = 3            // Maximum allowed unauthorized oracle access attempts
)

// RestrictedOracleAccessAutomation enforces restrictions on unauthorized oracle access
type RestrictedOracleAccessAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	oracleAccessViolation map[string]int // Tracks unauthorized oracle access attempts by user
}

// NewRestrictedOracleAccessAutomation initializes an instance of RestrictedOracleAccessAutomation
func NewRestrictedOracleAccessAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RestrictedOracleAccessAutomation {
	return &RestrictedOracleAccessAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		oracleAccessViolation: make(map[string]int),
	}
}

// StartOracleAccessMonitoring starts continuous monitoring of oracle access
func (automation *RestrictedOracleAccessAutomation) StartOracleAccessMonitoring() {
	ticker := time.NewTicker(OracleAccessCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorOracleAccess()
		}
	}()
}

// monitorOracleAccess checks for unauthorized oracle access attempts and enforces restrictions
func (automation *RestrictedOracleAccessAutomation) monitorOracleAccess() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch oracle access data from Synnergy Consensus
	oracleAccessData := automation.consensusSystem.GetOracleAccessData()

	for userID, accessStatus := range oracleAccessData {
		// Check if the access attempt is unauthorized
		if accessStatus == "unauthorized" {
			automation.flagOracleAccessViolation(userID, "Unauthorized oracle access attempt detected")
		}
	}
}

// flagOracleAccessViolation flags an unauthorized oracle access attempt and logs it in the ledger
func (automation *RestrictedOracleAccessAutomation) flagOracleAccessViolation(userID string, reason string) {
	fmt.Printf("Oracle access violation: User ID %s, Reason: %s\n", userID, reason)

	// Increment the violation count for the user
	automation.oracleAccessViolation[userID]++

	// Log the violation in the ledger
	automation.logOracleAccessViolation(userID, reason)

	// Check if the user has exceeded the allowed number of access violations
	if automation.oracleAccessViolation[userID] >= MaxOracleUnauthorizedAttempts {
		automation.restrictOracleAccess(userID)
	}
}

// logOracleAccessViolation logs the flagged oracle access violation into the ledger with details
func (automation *RestrictedOracleAccessAutomation) logOracleAccessViolation(userID string, violationReason string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("oracle-access-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Oracle Access Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated oracle access restrictions. Reason: %s", userID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptAccessData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log oracle access violation:", err)
	} else {
		fmt.Println("Oracle access violation logged.")
	}
}

// restrictOracleAccess restricts oracle access for a user after exceeding violations
func (automation *RestrictedOracleAccessAutomation) restrictOracleAccess(userID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("oracle-access-restriction-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Oracle Access Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("User %s has been restricted from oracle access due to repeated violations.", userID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptAccessData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log oracle access restriction:", err)
	} else {
		fmt.Println("Oracle access restriction applied.")
	}
}

// encryptAccessData encrypts the oracle access data before logging for security
func (automation *RestrictedOracleAccessAutomation) encryptAccessData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting oracle access data:", err)
		return data
	}
	return string(encryptedData)
}
