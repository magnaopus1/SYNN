package compliance

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)


// NewAuditTrail initializes a new audit trail system
func NewAuditTrail(ledgerInstance *ledger.Ledger) *AuditTrail {
    return &AuditTrail{
        Entries:       []AuditTrailEntry{},
        LedgerInstance: ledgerInstance,
    }
}

// AddEntry creates a new audit trail entry and records it in the system
func (at *AuditTrail) AddEntry(eventType, userID, details string) error {
    at.mutex.Lock()
    defer at.mutex.Unlock()

    eventID := generateEventID(eventType, userID)
    entry := AuditTrailEntry{
        EventID:   eventID,
        EventType: eventType,
        Timestamp: time.Now(),
        UserID:    userID,
        Details:   details,
    }

    at.Entries = append(at.Entries, entry)
    
    // Encrypt and store the entry in the ledger
    encryptionInstance := &common.Encryption{} // Assuming you have an instance of Encryption
    encryptedEntry, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%+v", entry)), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt audit entry: %v", err)
    }

    // RecordAuditEntry now has three arguments: eventID, encryptedEntry, and eventType
    result, err := at.LedgerInstance.ComplianceLedger.RecordAuditEntry(eventID, string(encryptedEntry), eventType)
    if err != nil {
        return fmt.Errorf("failed to record audit entry in ledger: %v", err)
    }

    fmt.Printf("Audit trail entry added: %s, Result: %v\n", eventID, result)
    return nil
}



// ListEntries lists all audit trail entries for review
func (at *AuditTrail) ListEntries() {
    fmt.Println("Audit Trail Entries:")
    for _, entry := range at.Entries {
        fmt.Printf("Event ID: %s | Type: %s | Timestamp: %s | User: %s | Details: %s\n", 
            entry.EventID, entry.EventType, entry.Timestamp.String(), entry.UserID, entry.Details)
    }
}

// RetrieveEntry retrieves and decrypts a specific audit entry from the ledger
func (at *AuditTrail) RetrieveEntry(eventID string) (AuditTrailEntry, error) {
    for _, entry := range at.Entries {
        if entry.EventID == eventID {
            encryptionInstance := &common.Encryption{}
            
            // Decrypt the details field
            decryptedDetails, err := encryptionInstance.DecryptData([]byte(entry.Details), common.EncryptionKey)
            if err != nil {
                return AuditTrailEntry{}, fmt.Errorf("failed to decrypt audit entry details: %v", err)
            }

            // Set the decrypted details back to the entry before returning it
            entry.Details = string(decryptedDetails)

            return entry, nil
        }
    }
    return AuditTrailEntry{}, fmt.Errorf("entry with EventID %s not found", eventID)
}


// generateEventID generates a unique ID for each audit event based on its type and user
func generateEventID(eventType, userID string) string {
    hashInput := fmt.Sprintf("%s%s%d", eventType, userID, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}
