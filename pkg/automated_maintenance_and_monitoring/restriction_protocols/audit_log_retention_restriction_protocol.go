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
	AuditLogRetentionCheckInterval = 24 * time.Hour // Interval for checking audit log retention
	MaxRetentionPeriod             = 365 * 24 * time.Hour // Maximum retention period of 1 year
)

// AuditLogRetentionRestrictionAutomation manages audit log retention, ensuring logs are stored for the correct duration and removed after expiry
type AuditLogRetentionRestrictionAutomation struct {
	consensusSystem      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	stateMutex           *sync.RWMutex
}

// NewAuditLogRetentionRestrictionAutomation initializes and returns an instance of AuditLogRetentionRestrictionAutomation
func NewAuditLogRetentionRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AuditLogRetentionRestrictionAutomation {
	return &AuditLogRetentionRestrictionAutomation{
		consensusSystem:    consensusSystem,
		ledgerInstance:     ledgerInstance,
		stateMutex:         stateMutex,
	}
}

// StartAuditLogRetentionMonitoring starts the loop for monitoring and enforcing audit log retention
func (automation *AuditLogRetentionRestrictionAutomation) StartAuditLogRetentionMonitoring() {
	ticker := time.NewTicker(AuditLogRetentionCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorAuditLogRetention()
		}
	}()
}

// monitorAuditLogRetention checks the ledger for logs that exceed the retention period and triggers removal if necessary
func (automation *AuditLogRetentionRestrictionAutomation) monitorAuditLogRetention() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Retrieve all audit logs from the ledger
	auditLogs := automation.ledgerInstance.GetAuditLogs()

	for _, logEntry := range auditLogs {
		// Check if the log exceeds the retention period
		if time.Since(time.Unix(logEntry.Timestamp, 0)) > MaxRetentionPeriod {
			automation.triggerLogRemoval(logEntry)
		}
	}
}

// triggerLogRemoval removes the audit log that has exceeded the retention period and logs the action into the ledger
func (automation *AuditLogRetentionRestrictionAutomation) triggerLogRemoval(logEntry common.LedgerEntry) {
	// Encrypt the log removal action
	encryptedLogID := automation.encryptLogID(logEntry.ID)

	// Remove the log from the ledger
	err := automation.ledgerInstance.RemoveLog(logEntry.ID)
	if err != nil {
		fmt.Printf("Failed to remove expired audit log: %v\n", err)
		return
	}

	// Log the removal of the expired log
	automation.logLogRemovalToLedger(logEntry.ID, encryptedLogID)
}

// logLogRemovalToLedger logs the action of removing an expired audit log into the ledger
func (automation *AuditLogRetentionRestrictionAutomation) logLogRemovalToLedger(logID, encryptedLogID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("audit-log-removal-%s-%d", logID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Audit Log Removal",
		Status:    "Removed",
		Details:   fmt.Sprintf("Audit log %s has been removed after exceeding the retention period. Encrypted Log ID: %s", logID, encryptedLogID),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log audit log removal into ledger: %v\n", err)
	} else {
		fmt.Printf("Audit log removal logged: %s\n", logID)
	}
}

// encryptLogID encrypts the log ID before logging its removal for security purposes
func (automation *AuditLogRetentionRestrictionAutomation) encryptLogID(logID string) string {
	encryptedData, err := encryption.EncryptData([]byte(logID))
	if err != nil {
		fmt.Println("Error encrypting log ID:", err)
		return logID
	}
	return string(encryptedData)
}
