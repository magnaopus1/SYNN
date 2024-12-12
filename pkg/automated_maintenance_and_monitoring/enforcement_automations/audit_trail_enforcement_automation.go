package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
)

const (
	AuditCheckInterval   = 60 * time.Second // Interval for generating audit logs
	MaxAuditLogRetention = 30 * 24 * time.Hour // Retain logs for 30 days
)

// AuditTrailEnforcementAutomation enforces audit trail generation and verification
type AuditTrailEnforcementAutomation struct {
	ledgerInstance   *ledger.Ledger
	consensusEngine  *consensus.SynnergyConsensus
	enforcementMutex *sync.RWMutex
}

// NewAuditTrailEnforcementAutomation initializes the audit trail enforcement automation
func NewAuditTrailEnforcementAutomation(ledgerInstance *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, enforcementMutex *sync.RWMutex) *AuditTrailEnforcementAutomation {
	return &AuditTrailEnforcementAutomation{
		ledgerInstance:   ledgerInstance,
		consensusEngine:  consensusEngine,
		enforcementMutex: enforcementMutex,
	}
}

// StartAuditTrailEnforcement starts the audit trail automation, generating and validating logs
func (automation *AuditTrailEnforcementAutomation) StartAuditTrailEnforcement() {
	ticker := time.NewTicker(AuditCheckInterval)

	go func() {
		for range ticker.C {
			automation.generateAuditLogs()
			automation.cleanUpOldAuditLogs()
		}
	}()
}

// generateAuditLogs generates audit logs for transactions and network events
func (automation *AuditTrailEnforcementAutomation) generateAuditLogs() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	// Get recent transactions from the consensus engine
	transactions := automation.consensusEngine.GetRecentTransactions()

	for _, tx := range transactions {
		auditEntry := common.LedgerEntry{
			ID:        fmt.Sprintf("audit-trail-%s-%d", tx.ID, time.Now().Unix()),
			Timestamp: time.Now().Unix(),
			Type:      "Audit Trail",
			Status:    "Generated",
			Details:   fmt.Sprintf("Transaction ID: %s, Amount: %f, Sender: %s, Receiver: %s", tx.ID, tx.Amount, tx.Sender, tx.Receiver),
		}

		// Encrypt audit details
		encryptedDetails := automation.encryptData(auditEntry.Details)
		auditEntry.Details = encryptedDetails

		// Log the audit trail into the ledger
		err := automation.ledgerInstance.AddEntry(auditEntry)
		if err != nil {
			fmt.Printf("Failed to add audit log for transaction %s: %v\n", tx.ID, err)
		} else {
			fmt.Printf("Audit log generated for transaction %s.\n", tx.ID)
		}
	}
}

// cleanUpOldAuditLogs removes audit logs that exceed the maximum retention period
func (automation *AuditTrailEnforcementAutomation) cleanUpOldAuditLogs() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	// Retrieve all audit trail entries
	allEntries, err := automation.ledgerInstance.GetEntriesByType("Audit Trail")
	if err != nil {
		fmt.Printf("Failed to retrieve audit logs: %v\n", err)
		return
	}

	// Iterate through entries and remove old logs
	for _, entry := range allEntries {
		if time.Now().Unix()-entry.Timestamp > int64(MaxAuditLogRetention.Seconds()) {
			err := automation.ledgerInstance.RemoveEntry(entry.ID)
			if err != nil {
				fmt.Printf("Failed to remove old audit log %s: %v\n", entry.ID, err)
			} else {
				fmt.Printf("Old audit log %s removed successfully.\n", entry.ID)
			}
		}
	}
}

// logAuditVerification logs the verification of audit trails into the ledger
func (automation *AuditTrailEnforcementAutomation) logAuditVerification(auditID, status string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("audit-verification-%s-%d", auditID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Audit Verification",
		Status:    status,
		Details:   fmt.Sprintf("Audit Trail ID: %s verification status: %s", auditID, status),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log audit verification for audit ID %s: %v\n", auditID, err)
	} else {
		fmt.Println("Audit verification log successfully added to ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *AuditTrailEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualAudit allows administrators to manually trigger an audit trail for a specific transaction
func (automation *AuditTrailEnforcementAutomation) TriggerManualAudit(transactionID string) {
	fmt.Printf("Manually generating audit trail for transaction: %s\n", transactionID)

	// Retrieve transaction details
	tx, err := automation.consensusEngine.GetTransactionByID(transactionID)
	if err != nil {
		fmt.Printf("Failed to retrieve transaction %s: %v\n", transactionID, err)
		return
	}

	auditEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("manual-audit-trail-%s-%d", tx.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Manual Audit Trail",
		Status:    "Generated",
		Details:   fmt.Sprintf("Transaction ID: %s, Amount: %f, Sender: %s, Receiver: %s", tx.ID, tx.Amount, tx.Sender, tx.Receiver),
	}

	// Encrypt audit details
	encryptedDetails := automation.encryptData(auditEntry.Details)
	auditEntry.Details = encryptedDetails

	// Log the audit trail into the ledger
	err = automation.ledgerInstance.AddEntry(auditEntry)
	if err != nil {
		fmt.Printf("Failed to manually add audit log for transaction %s: %v\n", tx.ID, err)
	} else {
		fmt.Printf("Manual audit log generated for transaction %s.\n", tx.ID)
	}
}
