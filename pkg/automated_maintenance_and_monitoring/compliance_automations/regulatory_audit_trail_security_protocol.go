package compliance_automations

import (
    "bytes"
    "crypto/sha256"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"

    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    AuditCheckInterval   = 30 * time.Minute // Interval for audit trail verification
    MaxAuditEntryBacklog = 1000              // Max number of audit entries before triggering an integrity check
)

// RegulatoryAuditTrailAutomation manages the security and integrity of the audit trail
type RegulatoryAuditTrailAutomation struct {
    ledgerInstance  *ledger.Ledger               // Blockchain ledger instance
    consensusEngine *synnergy_consensus.Consensus // Synnergy Consensus engine for audit validation
    stateMutex      *sync.RWMutex                // Mutex for thread-safe access
    apiURL          string                       // API URL for audit trail endpoints
    auditBacklog    int                          // Counter to track the number of audit entries since the last integrity check
}

// NewRegulatoryAuditTrailAutomation initializes the automation for audit trail security
func NewRegulatoryAuditTrailAutomation(apiURL string, ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Consensus, stateMutex *sync.RWMutex) *RegulatoryAuditTrailAutomation {
    return &RegulatoryAuditTrailAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        apiURL:          apiURL,
        auditBacklog:    0,
    }
}

// StartAuditTrailLogging starts the continuous logging and verification of audit entries
func (automation *RegulatoryAuditTrailAutomation) StartAuditTrailLogging() {
    ticker := time.NewTicker(AuditCheckInterval)
    for range ticker.C {
        fmt.Println("Starting audit trail verification...")
        automation.verifyAuditTrailIntegrity()
    }
}

// AddAuditEntry securely adds an entry to the audit trail
func (automation *RegulatoryAuditTrailAutomation) AddAuditEntry(entry common.AuditEntry) error {
    // Encrypt the audit entry before sending it to the audit trail
    encryptedEntry, err := encryption.EncryptAuditEntry(entry, []byte(EncryptionKey))
    if err != nil {
        return fmt.Errorf("failed to encrypt audit entry: %v", err)
    }

    // Add the entry to the audit trail via the API
    url := fmt.Sprintf("%s/api/compliance/audit/add_entry", automation.apiURL)
    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedEntry))
    if err != nil || resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to add audit entry: %v", err)
    }

    fmt.Printf("Audit entry added for transaction ID: %s\n", entry.TransactionID)
    
    // Increment the audit backlog and check if it exceeds the limit
    automation.auditBacklog++
    if automation.auditBacklog >= MaxAuditEntryBacklog {
        automation.verifyAuditTrailIntegrity()
        automation.auditBacklog = 0
    }

    return nil
}

// verifyAuditTrailIntegrity retrieves and verifies the integrity of audit entries from the audit trail
func (automation *RegulatoryAuditTrailAutomation) verifyAuditTrailIntegrity() {
    url := fmt.Sprintf("%s/api/compliance/audit/list_entries", automation.apiURL)
    resp, err := http.Get(url)
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error retrieving audit entries: %v\n", err)
        return
    }
    defer resp.Body.Close()

    var auditEntries []common.AuditEntry
    json.NewDecoder(resp.Body).Decode(&auditEntries)

    for _, entry := range auditEntries {
        if !automation.isAuditEntryValid(entry) {
            fmt.Printf("Audit entry for transaction ID %s has been tampered with.\n", entry.TransactionID)
            automation.flagTamperedEntry(entry)
        }
    }

    fmt.Println("Audit trail integrity verification completed.")
}

// isAuditEntryValid checks if the retrieved audit entry matches its hash and is untampered
func (automation *RegulatoryAuditTrailAutomation) isAuditEntryValid(entry common.AuditEntry) bool {
    // Recalculate the hash of the audit entry data and compare it to the stored hash
    hash := sha256.New()
    hash.Write([]byte(entry.TransactionID + entry.Action + entry.Timestamp))

    calculatedHash := fmt.Sprintf("%x", hash.Sum(nil))
    if calculatedHash != entry.Hash {
        fmt.Printf("Audit entry hash mismatch for transaction ID: %s\n", entry.TransactionID)
        return false
    }

    return true
}

// flagTamperedEntry flags a tampered audit entry and logs it for regulatory purposes
func (automation *RegulatoryAuditTrailAutomation) flagTamperedEntry(entry common.AuditEntry) {
    // Mark the entry as tampered in the ledger
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    tamperedEntry := common.LedgerEntry{
        ID:        entry.TransactionID,
        Timestamp: time.Now().Unix(),
        Type:      "Audit Entry Tampered",
        Status:    "Tampered",
    }

    // Encrypt the tampered entry before adding it to the ledger
    encryptedEntry, err := encryption.EncryptLedgerEntry(tamperedEntry, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting tampered audit entry: %v\n", err)
        return
    }

    // Add the tampered entry to the ledger with consensus validation
    automation.consensusEngine.ValidateSubBlock(tamperedEntry)
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Tampered audit entry flagged for transaction ID: %s\n", entry.TransactionID)
}

// RetrieveAuditEntry retrieves a specific audit entry based on the transaction ID
func (automation *RegulatoryAuditTrailAutomation) RetrieveAuditEntry(transactionID string) (common.AuditEntry, error) {
    url := fmt.Sprintf("%s/api/compliance/audit/retrieve_entry?transactionID=%s", automation.apiURL, transactionID)
    resp, err := http.Get(url)
    if err != nil || resp.StatusCode != http.StatusOK {
        return common.AuditEntry{}, fmt.Errorf("failed to retrieve audit entry for transaction ID %s: %v", transactionID, err)
    }
    defer resp.Body.Close()

    var entry common.AuditEntry
    err = json.NewDecoder(resp.Body).Decode(&entry)
    if err != nil {
        return common.AuditEntry{}, fmt.Errorf("failed to decode audit entry: %v", err)
    }

    return entry, nil
}

// ContinuousLogging continuously logs significant regulatory events to the audit trail
func (automation *RegulatoryAuditTrailAutomation) ContinuousLogging() {
    ticker := time.NewTicker(5 * time.Minute) // Log regulatory events every 5 minutes
    for range ticker.C {
        fmt.Println("Logging regulatory events...")

        events := automation.ledgerInstance.GetRegulatoryEvents()
        for _, event := range events {
            entry := common.AuditEntry{
                TransactionID: event.TransactionID,
                Action:        event.Action,
                Timestamp:     event.Timestamp,
                Hash:          automation.calculateAuditHash(event),
            }

            if err := automation.AddAuditEntry(entry); err != nil {
                fmt.Printf("Error logging regulatory event for transaction ID %s: %v\n", entry.TransactionID, err)
            }
        }
    }
}

// calculateAuditHash calculates a unique hash for an audit entry
func (automation *RegulatoryAuditTrailAutomation) calculateAuditHash(event common.RegulatoryEvent) string {
    hash := sha256.New()
    hash.Write([]byte(event.TransactionID + event.Action + event.Timestamp))
    return fmt.Sprintf("%x", hash.Sum(nil))
}
