package compliance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewComplianceAddition initializes a new compliance addition system
func NewComplianceAddition(rules []string, ledgerInstance *ledger.Ledger) *ComplianceAddition {
    return &ComplianceAddition{
        ComplianceRules: rules,
        LedgerInstance:  ledgerInstance,
    }
}

// CheckCompliance performs compliance checks for a specific action or transaction
func (ca *ComplianceAddition) CheckCompliance(actionID, checkedBy string, actionData string) (*ComplianceRecord, error) {
    ca.mutex.Lock()
    defer ca.mutex.Unlock()

    // Run compliance checks
    status := ca.runComplianceChecks(actionData)
    complianceRecord := &ComplianceRecord{
        ActionID:  actionID,
        Status:    status,
        CheckedBy: checkedBy,
    }

    // If the action is not compliant, return an error
    if !status.IsCompliant {
        return complianceRecord, fmt.Errorf("compliance check failed for Action ID %s", actionID)
    }

    // Encrypt and store the compliance record in the ledger
    encryptionInstance := &common.Encryption{} // Create an instance of Encryption
    encryptedRecord, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%+v", complianceRecord)), common.EncryptionKey)
    if err != nil {
        return complianceRecord, fmt.Errorf("failed to encrypt compliance record: %v", err)
    }

    // Record the compliance result in the ledger
    result, err := ca.LedgerInstance.ComplianceLedger.RecordCompliance(actionID, string(encryptedRecord))
    if err != nil {
        return complianceRecord, fmt.Errorf("failed to store compliance record in ledger: %v", err)
    }

    fmt.Printf("Compliance check for Action ID %s passed and recorded. Result: %v\n", actionID, result)
    return complianceRecord, nil
}



// runComplianceChecks runs a series of predefined checks on the action data.
func (ca *ComplianceAddition) runComplianceChecks(actionData string) ComplianceStatus {
    isCompliant := true
    lastCheckTime := time.Now()

    for _, rule := range ca.ComplianceRules {
        if !checkRule(rule, actionData) {
            isCompliant = false
            break
        }
    }

    return ComplianceStatus{
        IsCompliant:   isCompliant,
        LastCheckTime: lastCheckTime,
        NextCheckTime: lastCheckTime.Add(24 * time.Hour), // example: schedule next check in 24 hours
    }
}




// RetrieveComplianceRecord retrieves a compliance record from the ledger and decrypts it
func (ca *ComplianceAddition) RetrieveComplianceRecord(actionID string) (*ComplianceRecord, error) {
    // Retrieve the encrypted record from the ledger
    complianceRecord, err := ca.LedgerInstance.ComplianceLedger.GetComplianceRecord(actionID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve compliance record: %v", err)
    }

    // Ensure you have an instance of the Encryption struct
    encryptionInstance := &common.Encryption{}

    // Decrypt the EncryptedData field using AES decryption
    decryptedRecordBytes, err := encryptionInstance.DecryptData([]byte(complianceRecord.EncryptedData), common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt compliance record: %v", err)
    }

    // Deserialize the decrypted data back into the original struct
    var decryptedComplianceRecord ComplianceRecord
    err = json.Unmarshal(decryptedRecordBytes, &decryptedComplianceRecord)
    if err != nil {
        return nil, fmt.Errorf("failed to deserialize compliance record: %v", err)
    }

    return &decryptedComplianceRecord, nil
}

// generateActionID generates a unique ID for the compliance action based on its data
func generateActionID(actionData string) string {
    hashInput := fmt.Sprintf("%s%d", actionData, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}
