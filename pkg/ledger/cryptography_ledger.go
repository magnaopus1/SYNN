package ledger

import (
	"errors"
	"fmt"
	"time"
)

// StoreEncryptedKey stores the encrypted key in the ledger
func (l *CryptographyLedger) StoreEncryptedKey(keyID string, encryptedKey []byte) error {
	// Ensure "encryptedKeys" is initialized
	if l.EncryptedKeys == nil {
		l.EncryptedKeys = make(map[string][]byte)
	}

	// Check if the key already exists
	if _, exists := l.EncryptedKeys[keyID]; exists {
		return fmt.Errorf("key with ID %s already exists", keyID)
	}

	l.EncryptedKeys[keyID] = encryptedKey
	return nil
}


// RecordProofGeneration records the generation of a zk-proof for a specific transaction.
func (l *CryptographyLedger) RecordProofGeneration(proofID, transactionID, proofData string) error {
	l.Lock()
	defer l.Unlock()

	// Initialize zkProofRecords if it's nil
	if l.ZKProofRecords == nil {
		l.ZKProofRecords = make(map[string]*ZKProofRecord)
	}

	// Check if the proof already exists
	if _, exists := l.ZKProofRecords[proofID]; exists {
		return errors.New("zk-proof already exists")
	}

	// Record the zk-proof
	newProof := &ZKProofRecord{
		ProofID:       proofID,
		TransactionID: transactionID,
		ProofData:     proofData,
		GeneratedAt:   time.Now(),
		ProofStatus:   "generated",
	}

	// Store the proof
	l.ZKProofRecords[proofID] = newProof
	return nil
}

// RecordProofValidation logs the validation result of a zk-proof.
func (l *CryptographyLedger) RecordZKProofValidation(proofID, validatorID string, isValid bool) error {
	l.Lock()
	defer l.Unlock()

	// Initialize zkValidationRecords if it's nil
	if l.ZKValidationRecords == nil {
		l.ZKValidationRecords = make(map[string]*ZKValidationRecord)
	}

	// Check if the proof exists
	proof, exists := l.ZKProofRecords[proofID]
	if !exists {
		return errors.New("zk-proof not found")
	}

	// Update proof status
	if isValid {
		proof.ProofStatus = "validated"
	} else {
		proof.ProofStatus = "invalidated"
	}
	l.ZKProofRecords[proofID] = proof

	// Record the validation details
	validationRecord := &ZKValidationRecord{
		ProofID:     proofID,
		ValidatorID: validatorID,
		IsValid:     isValid,
		ValidatedAt: time.Now(),
	}

	// Store the validation record
	l.ZKValidationRecords[proofID] = validationRecord
	return nil
}

// GetProofByID retrieves a zk-proof by its proof ID.
func (l *CryptographyLedger) GetProofByID(proofID string) (*ZKProofRecord, error) {
	l.Lock()
	defer l.Unlock()

	// Check if the proof exists
	proof, exists := l.ZKProofRecords[proofID]
	if !exists {
		return nil, errors.New("zk-proof not found")
	}

	return proof, nil
}

// GetValidationRecordByProofID retrieves a zk-proof validation record by its proof ID.
func (l *CryptographyLedger) GetValidationRecordByProofID(proofID string) (*ZKValidationRecord, error) {
	l.Lock()
	defer l.Unlock()

	// Check if the validation record exists
	validationRecord, exists := l.ZKValidationRecords[proofID]
	if !exists {
		return nil, errors.New("validation record not found for this zk-proof")
	}

	return validationRecord, nil
}

// GetZkProofByID retrieves a ZKProofRecord by its ID from the ledger.
func (l *CryptographyLedger) GetZkProofByID(proofID string) (*ZKProofRecord, error) {
    // Assuming ZKProofRecords are stored in the ledger by their proofID
    proof, exists := l.ZKProofRecords[proofID]
    if !exists {
        return nil, fmt.Errorf("zk-proof with ID %s not found", proofID)
    }
    return proof, nil // Return the ZKProofRecord directly
}


// RecordKeySharing logs the sharing of encryption keys between entities.
func (l *CryptographyLedger) RecordKeySharing(entityID, sharedWith, key string) error {
	l.Lock()
	defer l.Unlock()

	// Initialize keySharings map if it's nil
	if l.KeySharings == nil {
		l.KeySharings = make(map[string][]KeySharing)
	}

	// Create a new KeySharing record
	keyShare := KeySharing{
		EntityID:   entityID,
		SharedWith: sharedWith,
		Key:        key,
		SharedAt:   time.Now(),
	}

	// Append the key sharing record to the ledger
	l.KeySharings[entityID] = append(l.KeySharings[entityID], keyShare)
	return nil
}

// LogHashEvent logs a hashing operation.
func (l *CryptographyLedger) LogHashEvent(eventType string, details []byte) error {
    logMessage := fmt.Sprintf("[%s] %s", eventType, string(details))
    l.HashLogs = append(l.HashLogs, logMessage)
    fmt.Println("Hash Event Logged:", logMessage)
    return nil
}

// LogEncryptionEvent logs an encryption operation.
func (l *CryptographyLedger) LogEncryptionEvent(eventType string, details []byte) error {
    logMessage := fmt.Sprintf("[%s] %s", eventType, string(details))
    l.EncryptionLogs = append(l.EncryptionLogs, logMessage)
    fmt.Println("Encryption Event Logged:", logMessage)
    return nil
}