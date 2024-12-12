package compliance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewDataProtectionPolicy initializes a new data protection policy
func NewDataProtectionPolicy(policyID, encryptionMethod, enforcer string, ledgerInstance *ledger.Ledger) *DataProtectionPolicy {
    return &DataProtectionPolicy{
        PolicyID:         generatePolicyID(policyID, enforcer),
        EncryptionMethod: encryptionMethod,
        CreatedAt:        time.Now(),
        EnforcedBy:       enforcer,
        LedgerInstance:   ledgerInstance,
    }
}

// ApplyDataProtection encrypts data according to the policy and records the action in the ledger
func (dpp *DataProtectionPolicy) ApplyDataProtection(data string) (*DataProtectionRecord, error) {
    dpp.mutex.Lock()
    defer dpp.mutex.Unlock()

    fmt.Printf("Applying data protection with policy ID %s by %s\n", dpp.PolicyID, dpp.EnforcedBy)

    // Hash the data for recording and encryption reference
    dataHash := hashData(data)

    // Log the data protection action
    protectionRecord := &DataProtectionRecord{
        PolicyID:    dpp.PolicyID,
        DataHash:    dataHash,
        IsEncrypted: true,
        Timestamp:   time.Now(),
    }

    // Encrypt and store the protection record in the ledger
    encryptionInstance := &common.Encryption{}
    encryptedRecord, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%+v", protectionRecord)), common.EncryptionKey)
    if err != nil {
        return protectionRecord, fmt.Errorf("failed to encrypt protection record: %v", err)
    }

    // RecordDataProtection now takes three arguments: policyID, encryptedRecord, and dataHash
    recordResult, err := dpp.LedgerInstance.ComplianceLedger.RecordDataProtection(dpp.PolicyID, string(encryptedRecord), dataHash)
    if err != nil {
        return protectionRecord, fmt.Errorf("failed to record data protection in ledger: %v", err)
    }

    fmt.Printf("Data protection applied successfully and logged for Policy ID %s. Result: %s\n", dpp.PolicyID, recordResult)
    return protectionRecord, nil
}



// RetrieveDataProtection retrieves and decrypts a data protection record from the ledger
func (dpp *DataProtectionPolicy) RetrieveDataProtection(policyID string) (*DataProtectionRecord, error) {
    // Get the encrypted data protection record from the ledger
    encryptedRecord, err := dpp.LedgerInstance.ComplianceLedger.GetDataProtectionRecord(policyID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve data protection record: %v", err)
    }

    // Assuming DataProtectionRecord has an EncryptedData field that stores the encrypted data
    encryptedData := []byte(encryptedRecord.EncryptedData) // Convert EncryptedData to []byte if it's a string

    // Decrypt the record using the encryption package
    encryptionInstance := &common.Encryption{}
    decryptedRecord, err := encryptionInstance.DecryptData(encryptedData, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt data protection record: %v", err)
    }

    // Parse the decrypted data into a DataProtectionRecord struct
    var record DataProtectionRecord
    if err := json.Unmarshal(decryptedRecord, &record); err != nil {
        return nil, fmt.Errorf("failed to parse data protection record: %v", err)
    }

    return &record, nil
}



// ValidateDataProtection ensures that the data protection policy is in place for specific data
func (dpp *DataProtectionPolicy) ValidateDataProtection(data string) error {
    dataHash := hashData(data)

    // Check the ledger for an existing protection record for this data
    record, err := dpp.RetrieveDataProtection(dpp.PolicyID)
    if err != nil {
        return err
    }

    if record.DataHash != dataHash {
        return errors.New("data protection validation failed: hash mismatch")
    }

    if !record.IsEncrypted {
        return errors.New("data is not encrypted")
    }

    fmt.Printf("Data protection validated for Policy ID %s\n", dpp.PolicyID)
    return nil
}


// generatePolicyID creates a unique identifier for a data protection policy
func generatePolicyID(policyID, enforcer string) string {
    input := fmt.Sprintf("%s-%s-%d", policyID, enforcer, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(input))
    return hex.EncodeToString(hash.Sum(nil))
}
