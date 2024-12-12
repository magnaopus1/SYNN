package syn1967

import (
	"fmt"
	"time"
	"sync"
)

// TokenStorageManager handles all storage, retrieval, and encryption-related tasks for SYN1967 tokens.
type TokenStorageManager struct {
	mu sync.Mutex
}

// StoreToken securely stores a SYN1967 token in the distributed storage after encryption.
func (tsm *TokenStorageManager) StoreToken(token *common.SYN1967Token) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Encrypt the token data before storing it
	encryptedData, err := encryption.Encrypt(token)
	if err != nil {
		return fmt.Errorf("failed to encrypt token data: %v", err)
	}

	// Store the encrypted token in the distributed storage system
	storageKey := generateStorageKey(token.TokenID)
	err = storage.Store(storageKey, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to store token in storage: %v", err)
	}

	// Log the storage in the ledger for traceability
	err = ledger.RecordTransaction(token.TokenID, encryptedData, "Token Storage")
	if err != nil {
		return fmt.Errorf("failed to record token storage in the ledger: %v", err)
	}

	return nil
}

// RetrieveToken retrieves and decrypts a SYN1967 token from the distributed storage.
func (tsm *TokenStorageManager) RetrieveToken(tokenID string) (*common.SYN1967Token, error) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Retrieve the encrypted token data from the storage
	storageKey := generateStorageKey(tokenID)
	encryptedData, err := storage.Retrieve(storageKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token from storage: %v", err)
	}

	// Decrypt the token data
	token, err := encryption.Decrypt(encryptedData, &common.SYN1967Token{})
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token data: %v", err)
	}

	return token, nil
}

// DeleteToken removes a SYN1967 token from the distributed storage, and records it in the ledger.
func (tsm *TokenStorageManager) DeleteToken(tokenID string) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Remove the token from the storage system
	storageKey := generateStorageKey(tokenID)
	err := storage.Delete(storageKey)
	if err != nil {
		return fmt.Errorf("failed to delete token from storage: %v", err)
	}

	// Log the deletion in the ledger
	err = ledger.RecordTransaction(tokenID, nil, "Token Deletion")
	if err != nil {
		return fmt.Errorf("failed to record token deletion in the ledger: %v", err)
	}

	return nil
}

// ValidateStorage ensures that a SYN1967 token is stored securely and is part of the validated sub-block.
func (tsm *TokenStorageManager) ValidateStorage(tokenID string, subBlockID string) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Validate that the token transaction is part of the correct sub-block
	valid, err := subblock.ValidateSubBlockTransaction(subBlockID, tokenID)
	if err != nil {
		return fmt.Errorf("failed to validate sub-block transaction: %v", err)
	}
	if !valid {
		return fmt.Errorf("token %s is not part of the validated sub-block", tokenID)
	}

	return nil
}

// UpdateToken allows for updating token information (like price or certification) in storage, with full encryption and ledger recording.
func (tsm *TokenStorageManager) UpdateToken(token *common.SYN1967Token) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Encrypt the updated token data
	encryptedData, err := encryption.Encrypt(token)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated token data: %v", err)
	}

	// Update the token in the distributed storage system
	storageKey := generateStorageKey(token.TokenID)
	err = storage.Update(storageKey, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to update token in storage: %v", err)
	}

	// Log the update in the ledger
	err = ledger.RecordTransaction(token.TokenID, encryptedData, "Token Update")
	if err != nil {
		return fmt.Errorf("failed to record token update in the ledger: %v", err)
	}

	return nil
}

// GenerateAuditReport generates a report of all actions related to the token storage (creation, retrieval, updates, deletion).
func (tsm *TokenStorageManager) GenerateAuditReport(tokenID string) (*common.AuditReport, error) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Retrieve the audit trail from the ledger
	auditRecords, err := ledger.RetrieveAuditRecords(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve audit records for token: %v", err)
	}

	// Compile the audit report
	auditReport := &common.AuditReport{
		TokenID:      tokenID,
		Timestamp:    time.Now(),
		AuditRecords: auditRecords,
	}

	return auditReport, nil
}

// generateStorageKey generates a unique key for storing tokens in the distributed storage.
func generateStorageKey(tokenID string) string {
	return fmt.Sprintf("SYN1967-%s", tokenID)
}
