package syn2900

import (

	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"time"
	"sync"
)

// StorageManager handles storage and retrieval operations for SYN2900 tokens.
type StorageManager struct {
	mu sync.Mutex
}

// NewStorageManager creates a new instance of StorageManager.
func NewStorageManager() *StorageManager {
	return &StorageManager{}
}

// SavePolicy stores a new insurance policy in the blockchain ledger.
func (sm *StorageManager) SavePolicy(policy common.InsurancePolicy) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Convert the policy to JSON format for storage
	data, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	// Encrypt the policy data before storing
	encryptedData, err := encryptData(data)
	if err != nil {
		return err
	}

	// Save encrypted policy to the ledger
	err = ledger.Store("policy_"+policy.PolicyID, encryptedData)
	if err != nil {
		return err
	}

	return nil
}

// GetPolicy retrieves an insurance policy by its Policy ID.
func (sm *StorageManager) GetPolicy(policyID string) (*common.InsurancePolicy, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Retrieve encrypted data from the ledger
	encryptedData, err := ledger.Retrieve("policy_" + policyID)
	if err != nil {
		return nil, errors.New("policy not found")
	}

	// Decrypt the policy data
	decryptedData, err := decryptData(encryptedData)
	if err != nil {
		return nil, err
	}

	// Convert JSON data back into InsurancePolicy struct
	var policy common.InsurancePolicy
	err = json.Unmarshal(decryptedData, &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

// UpdatePolicy updates an existing insurance policy.
func (sm *StorageManager) UpdatePolicy(policy common.InsurancePolicy) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Convert the updated policy to JSON format
	data, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	// Encrypt the updated data
	encryptedData, err := encryptData(data)
	if err != nil {
		return err
	}

	// Update the policy in the ledger
	err = ledger.Store("policy_"+policy.PolicyID, encryptedData)
	if err != nil {
		return err
	}

	return nil
}

// DeletePolicy deletes an insurance policy from the ledger.
func (sm *StorageManager) DeletePolicy(policyID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Delete the policy from the ledger
	err := ledger.Delete("policy_" + policyID)
	if err != nil {
		return errors.New("failed to delete policy")
	}

	return nil
}

// SaveClaim stores a new insurance claim for a specific policy.
func (sm *StorageManager) SaveClaim(claim common.InsuranceClaim) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Convert the claim to JSON format for storage
	data, err := json.Marshal(claim)
	if err != nil {
		return err
	}

	// Encrypt the claim data before storing
	encryptedData, err := encryptData(data)
	if err != nil {
		return err
	}

	// Save encrypted claim to the ledger
	err = ledger.Store("claim_"+claim.ClaimID, encryptedData)
	if err != nil {
		return err
	}

	return nil
}

// GetClaim retrieves an insurance claim by its Claim ID.
func (sm *StorageManager) GetClaim(claimID string) (*common.InsuranceClaim, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Retrieve encrypted claim data from the ledger
	encryptedData, err := ledger.Retrieve("claim_" + claimID)
	if err != nil {
		return nil, errors.New("claim not found")
	}

	// Decrypt the claim data
	decryptedData, err := decryptData(encryptedData)
	if err != nil {
		return nil, err
	}

	// Convert JSON data back into InsuranceClaim struct
	var claim common.InsuranceClaim
	err = json.Unmarshal(decryptedData, &claim)
	if err != nil {
		return nil, err
	}

	return &claim, nil
}

// UpdateClaim updates an existing insurance claim.
func (sm *StorageManager) UpdateClaim(claim common.InsuranceClaim) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Convert the updated claim to JSON format
	data, err := json.Marshal(claim)
	if err != nil {
		return err
	}

	// Encrypt the updated data
	encryptedData, err := encryptData(data)
	if err != nil {
		return err
	}

	// Update the claim in the ledger
	err = ledger.Store("claim_"+claim.ClaimID, encryptedData)
	if err != nil {
		return err
	}

	return nil
}

// DeleteClaim deletes an insurance claim from the ledger.
func (sm *StorageManager) DeleteClaim(claimID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Delete the claim from the ledger
	err := ledger.Delete("claim_" + claimID)
	if err != nil {
		return errors.New("failed to delete claim")
	}

	return nil
}

// encryptData encrypts the data using AES encryption.
func encryptData(data []byte) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext), nil
}

// decryptData decrypts the encrypted data using AES decryption.
func decryptData(encrypted string) ([]byte, error) {
	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("invalid ciphertext")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// encryptionKey is the 32-byte key used for AES encryption (must be kept secure).
var encryptionKey = []byte("your-encryption-key-32-bytes-long") // Replace with a secure key
