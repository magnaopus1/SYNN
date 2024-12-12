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
	DataAccessPrivilegeCheckInterval = 10 * time.Second  // Interval for checking data access privileges
	MaxDataAccessPerUser             = 100               // Maximum data access requests per user in a set period
	DataAccessPrivilegeWindow        = 24 * 7 * time.Hour // Time window for counting data access privileges (1 week)
)

// DataAccessPrivilegeRestrictionAutomation monitors and restricts data access privileges across the network
type DataAccessPrivilegeRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	userDataAccessCount    map[string]int // Tracks data access requests per user
}

// NewDataAccessPrivilegeRestrictionAutomation initializes and returns an instance of DataAccessPrivilegeRestrictionAutomation
func NewDataAccessPrivilegeRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataAccessPrivilegeRestrictionAutomation {
	return &DataAccessPrivilegeRestrictionAutomation{
		consensusSystem:     consensusSystem,
		ledgerInstance:      ledgerInstance,
		stateMutex:          stateMutex,
		userDataAccessCount: make(map[string]int),
	}
}

// StartDataAccessMonitoring starts continuous monitoring of data access privileges
func (automation *DataAccessPrivilegeRestrictionAutomation) StartDataAccessMonitoring() {
	ticker := time.NewTicker(DataAccessPrivilegeCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorDataAccessPrivileges()
		}
	}()
}

// monitorDataAccessPrivileges checks recent data access requests and enforces privilege limits
func (automation *DataAccessPrivilegeRestrictionAutomation) monitorDataAccessPrivileges() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent data access requests from Synnergy Consensus
	recentAccesses := automation.consensusSystem.GetRecentDataAccessRequests()

	for _, access := range recentAccesses {
		// Validate data access limits
		if !automation.validateDataAccessLimit(access) {
			automation.flagDataAccessViolation(access, "Exceeded maximum data access requests for this user")
		}
	}
}

// validateDataAccessLimit checks if a user has exceeded the data access limit within the time window
func (automation *DataAccessPrivilegeRestrictionAutomation) validateDataAccessLimit(access common.DataAccessRequest) bool {
	currentAccessCount := automation.userDataAccessCount[access.UserID]
	if currentAccessCount+1 > MaxDataAccessPerUser {
		return false
	}

	// Update the access count for the user
	automation.userDataAccessCount[access.UserID]++
	return true
}

// flagDataAccessViolation flags a data access request that violates system rules and logs it in the ledger
func (automation *DataAccessPrivilegeRestrictionAutomation) flagDataAccessViolation(access common.DataAccessRequest, reason string) {
	fmt.Printf("Data access violation: User %s, Reason: %s\n", access.UserID, reason)

	// Log the violation into the ledger
	automation.logDataAccessViolation(access, reason)
}

// logDataAccessViolation logs the flagged data access violation into the ledger with full details
func (automation *DataAccessPrivilegeRestrictionAutomation) logDataAccessViolation(access common.DataAccessRequest, violationReason string) {
	// Encrypt the data access request details
	encryptedData := automation.encryptDataAccessRequest(access)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("data-access-violation-%s-%d", access.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Data Access Privilege Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for data access violation. Reason: %s. Encrypted Data: %s", access.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log data access violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Data access violation logged for user: %s\n", access.UserID)
	}
}

// encryptDataAccessRequest encrypts data access request details before logging for security
func (automation *DataAccessPrivilegeRestrictionAutomation) encryptDataAccessRequest(access common.DataAccessRequest) string {
	data := fmt.Sprintf("User ID: %s, Data Access Timestamp: %d", access.UserID, access.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data access request:", err)
		return data
	}
	return string(encryptedData)
}
