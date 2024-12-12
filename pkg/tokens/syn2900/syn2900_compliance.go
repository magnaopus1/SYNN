package syn2900

import (

	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"time"
	"sync"
)

// ComplianceChecker ensures SYN2900 tokens comply with regulatory and internal rules.
type ComplianceChecker struct {
	mu sync.Mutex
}

// NewComplianceChecker creates a new instance of ComplianceChecker.
func NewComplianceChecker() *ComplianceChecker {
	return &ComplianceChecker{}
}

// ValidatePolicyCompliance checks whether an insurance policy token is compliant with specified regulations and standards.
func (c *ComplianceChecker) ValidatePolicyCompliance(tokenID string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Fetch the token from the ledger.
	encryptedToken, err := ledger.FetchToken(tokenID)
	if err != nil {
		return false, errors.New("failed to fetch token from ledger")
	}

	// Decrypt token data.
	token, err := decryptTokenData(encryptedToken)
	if err != nil {
		return false, errors.New("failed to decrypt token data")
	}

	// Perform various compliance checks.
	if err := c.checkActivePolicy(token); err != nil {
		return false, err
	}

	if err := c.checkCoverageLimits(token); err != nil {
		return false, err
	}

	if err := c.verifyRegulatoryRequirements(token); err != nil {
		return false, err
	}

	return true, nil
}

// checkActivePolicy verifies that the policy is active and not expired.
func (c *ComplianceChecker) checkActivePolicy(token common.SYN2900Token) error {
	if time.Now().After(token.EndDate) {
		return errors.New("policy has expired")
	}

	if !token.ActiveStatus {
		return errors.New("policy is no longer active")
	}

	return nil
}

// checkCoverageLimits ensures that the policy's coverage limits are valid and within allowed thresholds.
func (c *ComplianceChecker) checkCoverageLimits(token common.SYN2900Token) error {
	for _, coverage := range token.Coverages {
		if coverage.Limit < 0 {
			return errors.New("coverage limit cannot be negative")
		}
	}

	return nil
}

// verifyRegulatoryRequirements ensures that the token complies with jurisdiction-specific regulations.
func (c *ComplianceChecker) verifyRegulatoryRequirements(token common.SYN2900Token) error {
	// Check compliance with a series of jurisdiction-specific rules.
	// This can be expanded depending on specific insurance regulations in the real world.
	if token.Premium <= 0 {
		return errors.New("policy premium must be greater than zero")
	}

	if len(token.Owner) == 0 {
		return errors.New("policy must have a valid owner")
	}

	return nil
}

// AuditCompliance generates an audit trail for compliance checks and logs the results.
func (c *ComplianceChecker) AuditCompliance(tokenID string, compliant bool, details string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	log := common.ComplianceLog{
		TokenID:    tokenID,
		Timestamp:  time.Now(),
		Compliant:  compliant,
		Details:    details,
	}

	// Save the compliance log to the ledger.
	return ledger.SaveComplianceLog(log)
}

// Encrypt and decrypt functions using AES encryption for secure data handling.
var encryptionKey = []byte("your-encryption-key-32-bytes-long") // Replace with a secure key

// encryptTokenData encrypts the token data for secure storage.
func encryptTokenData(token common.SYN2900Token) (string, error) {
	plaintext, err := common.SerializeStruct(token)
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

// decryptTokenData decrypts the token data for use.
func decryptTokenData(encrypted string) (common.SYN2900Token, error) {
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

	var token common.SYN2900Token
	err = common.DeserializeStruct(plaintext, &token)
	if err != nil {
		return common.SYN2900Token{}, err
	}

	return token, nil
}
