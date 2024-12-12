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
	MaliciousActivityCheckInterval = 5 * time.Second // Interval for checking malicious activities
	MaxAllowedViolations           = 3               // Maximum allowed malicious activity violations before action is triggered
)

// MaliciousActivityRestrictionAutomation monitors and restricts malicious activities across the network
type MaliciousActivityRestrictionAutomation struct {
	consensusSystem          *consensus.SynnergyConsensus
	ledgerInstance           *ledger.Ledger
	stateMutex               *sync.RWMutex
	maliciousActivityCount   map[string]int // Tracks the number of malicious activities per user
}

// NewMaliciousActivityRestrictionAutomation initializes and returns an instance of MaliciousActivityRestrictionAutomation
func NewMaliciousActivityRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *MaliciousActivityRestrictionAutomation {
	return &MaliciousActivityRestrictionAutomation{
		consensusSystem:        consensusSystem,
		ledgerInstance:         ledgerInstance,
		stateMutex:             stateMutex,
		maliciousActivityCount: make(map[string]int),
	}
}

// StartMaliciousActivityMonitoring starts continuous monitoring of malicious activities across the network
func (automation *MaliciousActivityRestrictionAutomation) StartMaliciousActivityMonitoring() {
	ticker := time.NewTicker(MaliciousActivityCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorMaliciousActivities()
		}
	}()
}

// monitorMaliciousActivities checks for malicious activities by users and enforces restrictions if necessary
func (automation *MaliciousActivityRestrictionAutomation) monitorMaliciousActivities() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch malicious activity data from Synnergy Consensus
	maliciousData := automation.consensusSystem.GetMaliciousActivityData()

	for userID, activityDetails := range maliciousData {
		// Check if the user has exceeded the allowed number of malicious activities
		if automation.maliciousActivityCount[userID] > MaxAllowedViolations {
			automation.flagMaliciousActivityViolation(userID, activityDetails, "Exceeded allowed malicious activity violations")
		}
	}
}

// flagMaliciousActivityViolation flags a user's malicious activity violation and logs it in the ledger
func (automation *MaliciousActivityRestrictionAutomation) flagMaliciousActivityViolation(userID string, activityDetails string, reason string) {
	fmt.Printf("Malicious activity violation: User ID %s, Reason: %s\n", userID, reason)

	// Log the violation in the ledger
	automation.logMaliciousActivityViolation(userID, activityDetails, reason)
}

// logMaliciousActivityViolation logs the flagged malicious activity violation into the ledger with full details
func (automation *MaliciousActivityRestrictionAutomation) logMaliciousActivityViolation(userID string, activityDetails string, violationReason string) {
	// Create a ledger entry for the malicious activity violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("malicious-activity-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Malicious Activity Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s performed malicious activity. Details: %s. Reason: %s", userID, activityDetails, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptMaliciousData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log malicious activity violation:", err)
	} else {
		fmt.Println("Malicious activity violation logged.")
	}
}

// encryptMaliciousData encrypts the malicious activity data before logging for security
func (automation *MaliciousActivityRestrictionAutomation) encryptMaliciousData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting malicious activity data:", err)
		return data
	}
	return string(encryptedData)
}
