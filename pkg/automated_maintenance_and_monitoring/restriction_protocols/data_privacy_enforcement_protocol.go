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
	DataPrivacyCheckInterval       = 10 * time.Second  // Interval for checking data privacy compliance
	MaxPrivacyViolationsAllowed    = 3                 // Maximum number of privacy violations allowed per user
	DataPrivacyEnforcementWindow   = 24 * 7 * time.Hour // Time window for tracking privacy violations (1 week)
)

// DataPrivacyEnforcementAutomation monitors and enforces data privacy rules across the network
type DataPrivacyEnforcementAutomation struct {
	consensusSystem         *consensus.SynnergyConsensus
	ledgerInstance          *ledger.Ledger
	stateMutex              *sync.RWMutex
	userPrivacyViolationCount map[string]int // Tracks privacy violations per user
}

// NewDataPrivacyEnforcementAutomation initializes and returns an instance of DataPrivacyEnforcementAutomation
func NewDataPrivacyEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataPrivacyEnforcementAutomation {
	return &DataPrivacyEnforcementAutomation{
		consensusSystem:         consensusSystem,
		ledgerInstance:          ledgerInstance,
		stateMutex:              stateMutex,
		userPrivacyViolationCount: make(map[string]int),
	}
}

// StartPrivacyEnforcement starts continuous monitoring and enforcement of data privacy rules
func (automation *DataPrivacyEnforcementAutomation) StartPrivacyEnforcement() {
	ticker := time.NewTicker(DataPrivacyCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorPrivacyCompliance()
		}
	}()
}

// monitorPrivacyCompliance checks recent privacy activities and enforces privacy rules
func (automation *DataPrivacyEnforcementAutomation) monitorPrivacyCompliance() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent data access activities from Synnergy Consensus
	recentActivities := automation.consensusSystem.GetRecentDataActivities()

	for _, activity := range recentActivities {
		// Validate privacy compliance
		if !automation.validatePrivacyCompliance(activity) {
			automation.flagPrivacyViolation(activity, "Data privacy violation detected")
		}
	}
}

// validatePrivacyCompliance checks if the data access activity violates privacy rules
func (automation *DataPrivacyEnforcementAutomation) validatePrivacyCompliance(activity common.DataActivity) bool {
	currentViolationCount := automation.userPrivacyViolationCount[activity.UserID]
	if currentViolationCount+1 > MaxPrivacyViolationsAllowed {
		return false
	}

	// If no violation, reset violation count
	automation.userPrivacyViolationCount[activity.UserID] = 0
	return true
}

// flagPrivacyViolation flags a data access activity that violates privacy rules and logs it in the ledger
func (automation *DataPrivacyEnforcementAutomation) flagPrivacyViolation(activity common.DataActivity, reason string) {
	fmt.Printf("Data privacy violation: User %s, Reason: %s\n", activity.UserID, reason)

	// Log the violation into the ledger
	automation.logPrivacyViolation(activity, reason)
}

// logPrivacyViolation logs the flagged data privacy violation into the ledger with full details
func (automation *DataPrivacyEnforcementAutomation) logPrivacyViolation(activity common.DataActivity, violationReason string) {
	// Encrypt the privacy violation data
	encryptedData := automation.encryptPrivacyData(activity)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("privacy-violation-%s-%d", activity.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Data Privacy Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for privacy violation. Reason: %s. Encrypted Data: %s", activity.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log privacy violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Privacy violation logged for user: %s\n", activity.UserID)
	}
}

// encryptPrivacyData encrypts the privacy violation data before logging for security
func (automation *DataPrivacyEnforcementAutomation) encryptPrivacyData(activity common.DataActivity) string {
	data := fmt.Sprintf("User ID: %s, Activity Timestamp: %d", activity.UserID, activity.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting privacy violation data:", err)
		return data
	}
	return string(encryptedData)
}
