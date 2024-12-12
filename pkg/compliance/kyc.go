package compliance

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "time"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/common"
)

// KYCStatus represents the status of KYC verification
type KYCStatus string

const (
    Pending  KYCStatus = "Pending"
    Verified KYCStatus = "Verified"
    Rejected KYCStatus = "Rejected"
)


// NewKYCManager initializes a new KYC Manager
func NewKYCManager(ledgerInstance *ledger.Ledger) *KYCManager {
    return &KYCManager{
        Records:        make(map[string]KYCRecord),
        LedgerInstance: ledgerInstance,
    }
}

// SubmitKYC allows a user to submit their KYC data for verification
func (km *KYCManager) SubmitKYC(userID, kycData string) error {
    km.mutex.Lock()
    defer km.mutex.Unlock()

    if _, exists := km.Records[userID]; exists {
        return errors.New("KYC already submitted for this user")
    }

    // Hash the KYC data for record-keeping
    dataHash := hashData(kycData)

    // Encrypt the KYC data using the encryption package
    encryptionInstance := &common.Encryption{}
    encryptedKYC, err := encryptionInstance.EncryptData("AES", []byte(kycData), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt KYC data: %v", err)
    }

    // Store the encrypted KYC data as []byte
    km.Records[userID] = KYCRecord{
        UserID:       userID,
        Status:       Pending,
        VerifiedAt:   time.Time{},
        DataHash:     dataHash,
        EncryptedKYC: encryptedKYC, // Storing encrypted data as []byte
    }

    fmt.Printf("KYC data submitted for user %s\n", userID)
    return nil
}




// VerifyKYC verifies the KYC data for a user
func (km *KYCManager) VerifyKYC(userID string) error {
    km.mutex.Lock()
    defer km.mutex.Unlock()

    record, exists := km.Records[userID]
    if !exists {
        return errors.New("no KYC data found for this user")
    }

    if record.Status != Pending {
        return errors.New("KYC verification not in pending state")
    }

    // Update KYC record status
    record.Status = Verified
    record.VerifiedAt = time.Now()

    km.Records[userID] = record

    // Encrypt the KYC record for storage in the ledger
    encryptionInstance := &common.Encryption{}
    encryptedRecord, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%+v", record)), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt verified KYC record: %v", err)
    }

    // Convert record.Status (KYCStatus) to string and pass it to RecordKYC
    recordResult, err := km.LedgerInstance.ComplianceLedger.RecordKYC(userID, string(encryptedRecord), string(record.Status))
    if err != nil {
        return fmt.Errorf("failed to record KYC verification in ledger: %v", err)
    }

    fmt.Printf("KYC verified for user %s. Ledger record: %s\n", userID, recordResult)
    return nil
}



// RejectKYC rejects the KYC data for a user
func (km *KYCManager) RejectKYC(userID string) error {
    km.mutex.Lock()
    defer km.mutex.Unlock()

    record, exists := km.Records[userID]
    if !exists {
        return errors.New("no KYC data found for this user")
    }

    if record.Status != Pending {
        return errors.New("KYC verification not in pending state")
    }

    record.Status = Rejected
    km.Records[userID] = record

    fmt.Printf("KYC rejected for user %s\n", userID)
    return nil
}

// RetrieveKYC retrieves and decrypts the KYC record of a user
func (km *KYCManager) RetrieveKYC(userID string) (*KYCRecord, error) {
    record, exists := km.Records[userID]
    if !exists {
        return nil, errors.New("no KYC data found for this user")
    }

    // Decrypt the KYC data using the encryption package
    encryptionInstance := &common.Encryption{}
    decryptedKYC, err := encryptionInstance.DecryptData(record.EncryptedKYC, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt KYC data: %v", err)
    }

    fmt.Printf("KYC data retrieved for user %s\n", userID)
    record.EncryptedKYC = decryptedKYC // Replace with decrypted data for viewing
    return &record, nil
}


// ValidateKYC ensures that the KYC data matches the submitted information
func (km *KYCManager) ValidateKYC(userID, kycData string) error {
    record, exists := km.Records[userID]
    if !exists {
        return errors.New("no KYC data found for this user")
    }

    dataHash := hashData(kycData)
    if record.DataHash != dataHash {
        return errors.New("KYC validation failed: hash mismatch")
    }

    fmt.Printf("KYC data validated for user %s\n", userID)
    return nil
}

// hashData generates a SHA-256 hash for the KYC data
func hashData(data string) string {
    hash := sha256.New()
    hash.Write([]byte(data))
    return hex.EncodeToString(hash.Sum(nil))
}
