package syn2900

import (

	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"time"
	"sync"
)

// TokenManager manages the SYN2900 insurance tokens including policy issuance, claims, and premium payments.
type TokenManager struct {
	mu sync.Mutex
}

// NewTokenManager creates a new TokenManager instance.
func NewTokenManager() *TokenManager {
	return &TokenManager{}
}

// IssuePolicy issues a new insurance policy and records it in the ledger.
func (tm *TokenManager) IssuePolicy(policy common.SYN2900Token) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Validate the policy data
	if policy.TokenID == "" || policy.Owner == "" {
		return errors.New("invalid policy data")
	}

	// Encrypt sensitive policy data
	encryptedPolicy, err := encryptPolicyData(policy)
	if err != nil {
		return err
	}

	// Save the policy to the ledger
	err = ledger.SavePolicy(encryptedPolicy)
	if err != nil {
		return err
	}

	// Trigger PolicyIssued event
	eventManager := NewEventManager()
	err = eventManager.PolicyIssuedEvent(policy)
	if err != nil {
		return err
	}

	return nil
}

// ProcessClaim processes an insurance claim and updates the ledger accordingly.
func (tm *TokenManager) ProcessClaim(claim common.Claim) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Verify if the claim is valid
	valid, err := verifyClaim(claim)
	if err != nil || !valid {
		return errors.New("invalid claim")
	}

	// Encrypt claim data
	encryptedClaim, err := encryptClaimData(claim)
	if err != nil {
		return err
	}

	// Save claim to the ledger
	err = ledger.SaveClaim(encryptedClaim)
	if err != nil {
		return err
	}

	// Trigger ClaimCreated event
	eventManager := NewEventManager()
	err = eventManager.ClaimCreatedEvent(claim)
	if err != nil {
		return err
	}

	return nil
}

// PayPremium processes a premium payment and updates the policy balance.
func (tm *TokenManager) PayPremium(policy common.SYN2900Token, amount float64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Update the premium balance
	policy.Balance += amount

	// Encrypt updated policy data
	encryptedPolicy, err := encryptPolicyData(policy)
	if err != nil {
		return err
	}

	// Update the policy in the ledger
	err = ledger.UpdatePolicy(encryptedPolicy)
	if err != nil {
		return err
	}

	// Trigger PremiumPaid event
	eventManager := NewEventManager()
	err = eventManager.PremiumPaidEvent(policy, amount)
	if err != nil {
		return err
	}

	return nil
}

// ExpirePolicy handles the expiration of a policy by updating its status and notifying the ledger.
func (tm *TokenManager) ExpirePolicy(policy common.SYN2900Token) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Update the policy status to expired
	policy.ActiveStatus = false

	// Encrypt the updated policy data
	encryptedPolicy, err := encryptPolicyData(policy)
	if err != nil {
		return err
	}

	// Save the expired policy to the ledger
	err = ledger.UpdatePolicy(encryptedPolicy)
	if err != nil {
		return err
	}

	// Trigger PolicyExpired event
	eventManager := NewEventManager()
	err = eventManager.PolicyExpiredEvent(policy)
	if err != nil {
		return err
	}

	return nil
}

// VerifyClaim verifies whether a claim is valid based on policy coverage and terms.
func verifyClaim(claim common.Claim) (bool, error) {
	// Get policy details from the ledger
	policy, err := ledger.GetPolicy(claim.TokenID)
	if err != nil {
		return false, err
	}

	// Check if the policy covers the claim
	if policy.CoverageAmount < claim.ClaimAmount || !policy.ActiveStatus {
		return false, nil
	}

	// Further validation could involve checking dates, limits, etc.
	return true, nil
}

// Encrypt and decrypt functions using AES encryption for secure policy and claim data handling.
var encryptionKey = []byte("your-encryption-key-32-bytes-long") // Replace with a secure key

// encryptPolicyData encrypts policy data for secure storage.
func encryptPolicyData(policy common.SYN2900Token) (string, error) {
	plaintext, err := common.SerializeStruct(policy)
	if err != nil {
		return "", err
	}

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

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

// decryptPolicyData decrypts policy data for usage.
func decryptPolicyData(encrypted string) (common.SYN2900Token, error) {
	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return common.SYN2900Token{}, err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return common.SYN2900Token{}, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return common.SYN2900Token{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return common.SYN2900Token{}, errors.New("invalid ciphertext")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return common.SYN2900Token{}, err
	}

	var policy common.SYN2900Token
	err = common.DeserializeStruct(plaintext, &policy)
	if err != nil {
		return common.SYN2900Token{}, err
	}

	return policy, nil
}

// encryptClaimData encrypts claim data for secure storage.
func encryptClaimData(claim common.Claim) (string, error) {
	plaintext, err := common.SerializeStruct(claim)
	if err != nil {
		return "", err
	}

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

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

// decryptClaimData decrypts claim data for usage.
func decryptClaimData(encrypted string) (common.Claim, error) {
	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return common.Claim{}, err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return common.Claim{}, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return common.Claim{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return common.Claim{}, errors.New("invalid ciphertext")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return common.Claim{}, err
	}

	var claim common.Claim
	err = common.DeserializeStruct(plaintext, &claim)
	if err != nil {
		return common.Claim{}, err
	}

	return claim, nil
}
