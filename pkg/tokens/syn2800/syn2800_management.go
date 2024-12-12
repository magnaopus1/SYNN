package syn2800

import (
    "errors"
    "time"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
    "sync"
)

// PolicyManager manages all life insurance policies and related operations.
type PolicyManager struct {
	mutex sync.Mutex
}

// NewPolicyManager creates a new instance of PolicyManager.
func NewPolicyManager() *PolicyManager {
	return &PolicyManager{}
}

// AddPolicy creates and adds a new life insurance policy to the system.
func (pm *PolicyManager) AddPolicy(policy common.SYN2800Token) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Validate policy data
	if policy.TokenID == "" || policy.PolicyHolder == "" || policy.Beneficiary == "" {
		return fmt.Errorf("invalid policy data: required fields are missing")
	}

	// Set initial policy status and other details
	policy.IssueDate = time.Now()
	policy.LastPaymentDate = time.Now()
	policy.ActiveStatus = true

	// Encrypt and store the policy in the ledger
	encryptedPolicyData, err := pm.encryptPolicyData(&policy)
	if err != nil {
		return fmt.Errorf("failed to encrypt policy data: %v", err)
	}
	if err := ledger.StoreToken(policy.TokenID, encryptedPolicyData); err != nil {
		return fmt.Errorf("failed to store policy in ledger: %v", err)
	}

	log.Printf("New life insurance policy added with Token ID: %s", policy.TokenID)
	return nil
}

// UpdatePolicy updates the details of an existing life insurance policy.
func (pm *PolicyManager) UpdatePolicy(tokenID string, updateData map[string]interface{}) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Retrieve and decrypt the policy from the ledger
	policy, err := pm.retrieveAndDecryptPolicy(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve policy for update: %v", err)
	}

	// Update policy details based on the provided data
	if beneficiary, ok := updateData["beneficiary"].(string); ok {
		policy.Beneficiary = beneficiary
	}
	if coverageAmount, ok := updateData["coverageAmount"].(float64); ok {
		policy.CoverageAmount = coverageAmount
	}

	// Encrypt and store the updated policy back into the ledger
	encryptedPolicyData, err := pm.encryptPolicyData(policy)
	if err != nil {
		return fmt.Errorf("failed to encrypt policy data after update: %v", err)
	}
	if err := ledger.StoreToken(policy.TokenID, encryptedPolicyData); err != nil {
		return fmt.Errorf("failed to store updated policy in ledger: %v", err)
	}

	log.Printf("Policy updated for Token ID: %s", tokenID)
	return nil
}

// RemovePolicy removes a life insurance policy from the system.
func (pm *PolicyManager) RemovePolicy(tokenID string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Check if the policy exists in the ledger
	_, err := ledger.RetrieveToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve policy for removal: %v", err)
	}

	// Remove the policy from the ledger
	if err := ledger.RemoveToken(tokenID); err != nil {
		return fmt.Errorf("failed to remove policy from ledger: %v", err)
	}

	log.Printf("Life insurance policy removed with Token ID: %s", tokenID)
	return nil
}

// RenewPolicy renews a life insurance policy, extending its coverage.
func (pm *PolicyManager) RenewPolicy(tokenID string, newEndDate time.Time) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Retrieve and decrypt the policy from the ledger
	policy, err := pm.retrieveAndDecryptPolicy(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve policy for renewal: %v", err)
	}

	// Update the policy's end date and status
	policy.EndDate = newEndDate
	policy.ActiveStatus = true

	// Encrypt and store the renewed policy back into the ledger
	encryptedPolicyData, err := pm.encryptPolicyData(policy)
	if err != nil {
		return fmt.Errorf("failed to encrypt policy data after renewal: %v", err)
	}
	if err := ledger.StoreToken(policy.TokenID, encryptedPolicyData); err != nil {
		return fmt.Errorf("failed to store renewed policy in ledger: %v", err)
	}

	log.Printf("Policy renewed for Token ID: %s", tokenID)
	return nil
}

// GetPolicy retrieves a life insurance policy and returns its decrypted data.
func (pm *PolicyManager) GetPolicy(tokenID string) (*common.SYN2800Token, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Retrieve and decrypt the policy from the ledger
	policy, err := pm.retrieveAndDecryptPolicy(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve policy: %v", err)
	}

	return policy, nil
}

// Helper to retrieve and decrypt the policy from the ledger.
func (pm *PolicyManager) retrieveAndDecryptPolicy(tokenID string) (*common.SYN2800Token, error) {
	encryptedData, err := ledger.RetrieveToken(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve policy from ledger: %v", err)
	}
	return pm.decryptPolicyData(encryptedData)
}

// Encrypt and store policy data.
func (pm *PolicyManager) encryptPolicyData(policy *common.SYN2800Token) ([]byte, error) {
	key := generateEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	policyData := serializePolicyData(policy)
	return gcm.Seal(nonce, nonce, policyData, nil), nil
}

// Decrypt the policy data.
func (pm *PolicyManager) decryptPolicyData(encryptedData []byte) (*common.SYN2800Token, error) {
	key := generateEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return deserializePolicyData(decryptedData), nil
}

// Helper function to generate an encryption key.
func generateEncryptionKey() []byte {
	return []byte("your-secure-256-bit-key")
}

// Helper function to serialize policy data.
func serializePolicyData(policy *common.SYN2800Token) []byte {
	data, err := json.Marshal(policy)
	if err != nil {
		log.Fatalf("failed to serialize policy data: %v", err)
	}
	return data
}

// Helper function to deserialize policy data after decryption.
func deserializePolicyData(data []byte) *common.SYN2800Token {
	var policy common.SYN2800Token
	if err := json.Unmarshal(data, &policy); err != nil {
		log.Fatalf("failed to deserialize policy data: %v", err)
	}
	return &policy
}
