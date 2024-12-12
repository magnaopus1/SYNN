package compliance_automations

import (
    "bytes"
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
    DataRetentionCheckInterval = 30 * time.Minute // Interval for checking data retention and deletion
    DeletionKey                = "data_deletion_key" // Encryption key for deletion operations
)

// DataRetentionAndDeletionComplianceAutomation automates data retention and deletion compliance checks
type DataRetentionAndDeletionComplianceAutomation struct {
    ledgerInstance  *ledger.Ledger               // Blockchain ledger for data retention management
    consensusEngine *synnergy_consensus.Consensus // Synnergy Consensus Engine for validating data deletion operations
    stateMutex      *sync.RWMutex                // Mutex for thread-safe ledger access
    apiURL          string                       // API URL for data retention and deletion endpoints
}

// NewDataRetentionAndDeletionComplianceAutomation initializes the data retention and deletion compliance handler
func NewDataRetentionAndDeletionComplianceAutomation(apiURL string, ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Consensus, stateMutex *sync.RWMutex) *DataRetentionAndDeletionComplianceAutomation {
    return &DataRetentionAndDeletionComplianceAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        apiURL:          apiURL,
    }
}

// StartRetentionMonitoring initiates continuous monitoring for data retention and deletion compliance
func (automation *DataRetentionAndDeletionComplianceAutomation) StartRetentionMonitoring() {
    ticker := time.NewTicker(DataRetentionCheckInterval)
    for range ticker.C {
        fmt.Println("Checking data retention and deletion compliance...")
        automation.enforceDataRetentionPolicies()
        automation.validateDataRetention()
    }
}

// enforceDataRetentionPolicies applies retention policies and deletes expired data
func (automation *DataRetentionAndDeletionComplianceAutomation) enforceDataRetentionPolicies() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    records := automation.ledgerInstance.GetAllDataRecords() // Retrieve data records from ledger
    for _, record := range records {
        if automation.isDataExpired(record) {
            fmt.Printf("Data record ID %s is expired, applying deletion...\n", record.ID)
            automation.applyDeletion(record)
        }
    }
}

// isDataExpired checks if a data record is past its retention period
func (automation *DataRetentionAndDeletionComplianceAutomation) isDataExpired(record common.DataRecord) bool {
    return time.Now().After(record.RetentionExpiry) // Check if the retention period has expired
}

// applyDeletion securely deletes the expired data from the ledger and system
func (automation *DataRetentionAndDeletionComplianceAutomation) applyDeletion(record common.DataRecord) {
    url := fmt.Sprintf("%s/api/compliance/data_protection/apply", automation.apiURL)
    body, _ := json.Marshal(record)

    // Encrypt the data record before deletion
    encryptedBody, err := encryption.Encrypt(body, []byte(DeletionKey))
    if err != nil {
        fmt.Printf("Error encrypting data record ID %s for deletion: %v\n", record.ID, err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error applying deletion for data record ID %s: %v\n", record.ID, err)
        return
    }

    fmt.Printf("Data record ID %s successfully deleted.\n", record.ID)
    automation.updateLedgerForDeletion(record)
}

// updateLedgerForDeletion updates the ledger to track the deletion of the data record
func (automation *DataRetentionAndDeletionComplianceAutomation) updateLedgerForDeletion(record common.DataRecord) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        record.ID,
        Timestamp: time.Now().Unix(),
        Type:      "Data Deletion",
        Status:    "Deleted",
    }

    // Encrypt the ledger entry for additional security
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(DeletionKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for data deletion: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(record) // Synnergy Consensus validation
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated for data deletion on record ID: %s\n", record.ID)
}

// validateDataRetention ensures that data retention policies are applied consistently and that expired data is deleted
func (automation *DataRetentionAndDeletionComplianceAutomation) validateDataRetention() {
    url := fmt.Sprintf("%s/api/compliance/data_protection/validate", automation.apiURL)
    resp, err := http.Get(url)
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error validating data retention policies: %v\n", err)
        return
    }

    fmt.Println("Data retention validation successful.")
}

// retrieveRetentionRecords retrieves the details of retention records for auditing or compliance purposes
func (automation *DataRetentionAndDeletionComplianceAutomation) retrieveRetentionRecords(recordID string) {
    url := fmt.Sprintf("%s/api/compliance/data_protection/retrieve", automation.apiURL)
    body, _ := json.Marshal(map[string]string{"record_id": recordID})

    encryptedBody, err := encryption.Encrypt(body, []byte(DeletionKey))
    if err != nil {
        fmt.Printf("Error encrypting request for retention record details: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error retrieving retention records for record ID %s: %v\n", recordID, err)
        return
    }

    var retentionDetails common.DataRetentionRecord
    json.NewDecoder(resp.Body).Decode(&retentionDetails)
    fmt.Printf("Data retention details for record ID %s: %v\n", recordID, retentionDetails)
}

// addRetentionException adds an exception for retention policies and logs it in the ledger
func (automation *DataRetentionAndDeletionComplianceAutomation) addRetentionException(record common.DataRecord, reason string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        record.ID,
        Timestamp: time.Now().Unix(),
        Type:      "Retention Exception",
        Status:    "Exception Applied",
        Details:   reason,
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(DeletionKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for retention exception: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(record) // Synnergy Consensus validation
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated for data retention exception on record ID: %s, Reason: %s\n", record.ID, reason)
}
