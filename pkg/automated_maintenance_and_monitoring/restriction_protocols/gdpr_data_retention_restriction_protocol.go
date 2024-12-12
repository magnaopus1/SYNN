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
	DataRetentionCheckInterval  = 24 * time.Hour  // Interval for checking data retention policies
	MaxDataRetentionPeriod      = 365 * 24 * time.Hour // Maximum data retention period (1 year)
	DataPurgeNotificationWindow = 30 * 24 * time.Hour  // Notify users 30 days before their data is purged
)

// GdprDataRetentionRestrictionAutomation monitors and enforces GDPR data retention compliance across the network
type GdprDataRetentionRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	dataRetentionRecords   map[string]time.Time // Tracks data storage times per user
}

// NewGdprDataRetentionRestrictionAutomation initializes and returns an instance of GdprDataRetentionRestrictionAutomation
func NewGdprDataRetentionRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *GdprDataRetentionRestrictionAutomation {
	return &GdprDataRetentionRestrictionAutomation{
		consensusSystem:      consensusSystem,
		ledgerInstance:       ledgerInstance,
		stateMutex:           stateMutex,
		dataRetentionRecords: make(map[string]time.Time),
	}
}

// StartDataRetentionMonitoring starts continuous monitoring of GDPR data retention compliance
func (automation *GdprDataRetentionRestrictionAutomation) StartDataRetentionMonitoring() {
	ticker := time.NewTicker(DataRetentionCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorDataRetentionPolicies()
		}
	}()
}

// monitorDataRetentionPolicies checks data retention records and enforces GDPR compliance
func (automation *GdprDataRetentionRestrictionAutomation) monitorDataRetentionPolicies() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch current time
	currentTime := time.Now()

	for userID, storedTime := range automation.dataRetentionRecords {
		dataAge := currentTime.Sub(storedTime)
		
		if dataAge > MaxDataRetentionPeriod {
			automation.purgeDataForUser(userID)
		} else if dataAge > (MaxDataRetentionPeriod - DataPurgeNotificationWindow) {
			automation.notifyUserAboutPurge(userID)
		}
	}
}

// purgeDataForUser handles the purging of data for a user who exceeds the maximum data retention period
func (automation *GdprDataRetentionRestrictionAutomation) purgeDataForUser(userID string) {
	fmt.Printf("Purging data for user: %s due to GDPR data retention policy.\n", userID)

	// Log the data purge event into the ledger
	automation.logDataPurge(userID)

	// Purge data in the consensus system
	automation.consensusSystem.PurgeUserData(userID)

	// Remove the user's retention record
	delete(automation.dataRetentionRecords, userID)
}

// notifyUserAboutPurge notifies a user that their data is due for purging in 30 days
func (automation *GdprDataRetentionRestrictionAutomation) notifyUserAboutPurge(userID string) {
	fmt.Printf("Notifying user: %s about upcoming data purge due to GDPR compliance.\n", userID)

	// Log the notification event in the ledger
	automation.logDataPurgeNotification(userID)

	// Notify the user through the system (add email or notification system integration as needed)
	automation.consensusSystem.NotifyUserOfUpcomingDataPurge(userID)
}

// logDataPurge logs the data purge event into the ledger with full details
func (automation *GdprDataRetentionRestrictionAutomation) logDataPurge(userID string) {
	// Create a ledger entry for the data purge
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("data-purge-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Data Purge",
		Status:    "Purged",
		Details:   fmt.Sprintf("Data for user %s was purged due to exceeding GDPR data retention limits.", userID),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log data purge event into the ledger: %v\n", err)
	} else {
		fmt.Printf("Data purge event logged for user: %s\n", userID)
	}
}

// logDataPurgeNotification logs the notification event into the ledger
func (automation *GdprDataRetentionRestrictionAutomation) logDataPurgeNotification(userID string) {
	// Create a ledger entry for the notification
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("data-purge-notification-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Data Purge Notification",
		Status:    "Notified",
		Details:   fmt.Sprintf("User %s was notified of upcoming data purge due to GDPR retention policy.", userID),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log data purge notification into the ledger: %v\n", err)
	} else {
		fmt.Printf("Data purge notification logged for user: %s\n", userID)
	}
}

