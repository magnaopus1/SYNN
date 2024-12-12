package syn3200

import (
	"time"
	"errors"
	"sync"

)

// BillSecurity represents the security structure for a SYN3200 bill token.
type BillSecurity struct {
	TokenID      string    `json:"token_id"`
	SecurityHash string    `json:"security_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SecurityManager handles security checks and encryption for bill tokens.
type SecurityManager struct {
	securityRecords    map[string]*BillSecurity
	ledgerService      *ledger.Ledger
	encryptionService  *encryption.Encryptor
	consensusService   *consensus.SynnergyConsensus
	mutex              sync.Mutex
}

// NewSecurityManager creates a new instance of SecurityManager.
func NewSecurityManager(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *SecurityManager {
	return &SecurityManager{
		securityRecords:   make(map[string]*BillSecurity),
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// GenerateSecurityHash generates a security hash for a bill token.
func (sm *SecurityManager) GenerateSecurityHash(tokenID string) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Simulate the creation of a security hash for the token.
	securityHash, err := sm.encryptionService.GenerateHash(tokenID)
	if err != nil {
		return "", err
	}

	// Create a new security record for the bill token.
	securityRecord := &BillSecurity{
		TokenID:      tokenID,
		SecurityHash: securityHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Encrypt the security record.
	encryptedRecord, err := sm.encryptionService.EncryptData(securityRecord)
	if err != nil {
		return "", err
	}

	// Store the encrypted security record.
	sm.securityRecords[tokenID] = encryptedRecord.(*BillSecurity)

	// Log the creation of the security hash in the ledger.
	if err := sm.ledgerService.LogEvent("SecurityHashGenerated", time.Now(), tokenID); err != nil {
		return "", err
	}

	// Validate the security hash creation with consensus.
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return "", err
	}

	return securityHash, nil
}

// ValidateSecurityHash validates the security hash of a bill token.
func (sm *SecurityManager) ValidateSecurityHash(tokenID, providedHash string) (bool, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the security record.
	securityRecord, exists := sm.securityRecords[tokenID]
	if !exists {
		return false, errors.New("security record not found")
	}

	// Decrypt the security record.
	decryptedRecord, err := sm.encryptionService.DecryptData(securityRecord)
	if err != nil {
		return false, err
	}

	// Check if the provided hash matches the stored hash.
	if decryptedRecord.(*BillSecurity).SecurityHash == providedHash {
		// Log the validation in the ledger.
		if err := sm.ledgerService.LogEvent("SecurityHashValidated", time.Now(), tokenID); err != nil {
			return false, err
		}

		// Validate the security hash with consensus.
		if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

// UpdateSecurityHash updates the security hash for a bill token.
func (sm *SecurityManager) UpdateSecurityHash(tokenID string) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Generate a new security hash for the token.
	newHash, err := sm.encryptionService.GenerateHash(tokenID)
	if err != nil {
		return "", err
	}

	// Retrieve the current security record.
	securityRecord, exists := sm.securityRecords[tokenID]
	if !exists {
		return "", errors.New("security record not found")
	}

	// Update the security hash and timestamps.
	securityRecord.SecurityHash = newHash
	securityRecord.UpdatedAt = time.Now()

	// Encrypt the updated security record.
	encryptedRecord, err := sm.encryptionService.EncryptData(securityRecord)
	if err != nil {
		return "", err
	}

	// Store the updated security record.
	sm.securityRecords[tokenID] = encryptedRecord.(*BillSecurity)

	// Log the update in the ledger.
	if err := sm.ledgerService.LogEvent("SecurityHashUpdated", time.Now(), tokenID); err != nil {
		return "", err
	}

	// Validate the update with consensus.
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return "", err
	}

	return newHash, nil
}

// GetSecurityRecord retrieves the security record for a bill token.
func (sm *SecurityManager) GetSecurityRecord(tokenID string) (*BillSecurity, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the security record.
	securityRecord, exists := sm.securityRecords[tokenID]
	if !exists {
		return nil, errors.New("security record not found")
	}

	// Decrypt the security record before returning.
	decryptedRecord, err := sm.encryptionService.DecryptData(securityRecord)
	if err != nil {
		return nil, err
	}

	return decryptedRecord.(*BillSecurity), nil
}
