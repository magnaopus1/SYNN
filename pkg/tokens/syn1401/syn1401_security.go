package syn1401

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

// SYN1401Security handles the security, encryption, and audit functionalities for SYN1401 tokens.
type SYN1401Security struct {
	Ledger common.LedgerInterface // Interface to interact with the blockchain ledger.
	Keys   map[string][]byte      // A map holding encryption keys for token owners.
}

// EncryptSensitiveData encrypts sensitive token data (like ownership, principal) before it is stored in the ledger.
func (s *SYN1401Security) EncryptSensitiveData(owner string, token *common.SYN1401Token) error {
	// Retrieve encryption key for the token owner
	key, exists := s.Keys[owner]
	if !exists {
		return errors.New("encryption key for owner not found")
	}

	// Prepare the data to encrypt
	plaintext := []byte(fmt.Sprintf("%s|%f|%s", token.TokenID, token.PrincipalAmount, token.Owner))

	// Encrypt the plaintext using AES
	encryptedData, err := s.encryptAES(key, plaintext)
	if err != nil {
		return fmt.Errorf("error encrypting token data: %w", err)
	}

	// Assign the encrypted data to the token
	token.EncryptedMetadata = encryptedData
	return nil
}

// DecryptSensitiveData decrypts the encrypted metadata of the SYN1401Token.
func (s *SYN1401Security) DecryptSensitiveData(owner string, encryptedData []byte) (string, error) {
	// Retrieve decryption key for the token owner
	key, exists := s.Keys[owner]
	if !exists {
		return "", errors.New("decryption key for owner not found")
	}

	// Decrypt the data
	plaintext, err := s.decryptAES(key, encryptedData)
	if err != nil {
		return "", fmt.Errorf("error decrypting token data: %w", err)
	}

	return string(plaintext), nil
}

// ValidateToken performs a full security audit and validation of the SYN1401 token before it is processed.
func (s *SYN1401Security) ValidateToken(tokenID string) error {
	token, err := s.Ledger.GetToken(tokenID)
	if err != nil {
		return fmt.Errorf("error retrieving token from ledger: %w", err)
	}

	// Ensure the token is compliant and not tampered with
	if token.ComplianceStatus != "Compliant" {
		return errors.New("token is non-compliant, cannot proceed")
	}

	// Perform a basic security check on the token's metadata
	if len(token.EncryptedMetadata) == 0 {
		return errors.New("token metadata is not encrypted")
	}

	// Check if the redemption status is valid
	if token.RedemptionStatus == "Redeemed" {
		return errors.New("token already redeemed, invalid for further transactions")
	}

	// Create a validation audit log
	auditLog := common.AuditLog{
		AuditID:     generateUniqueID(),
		PerformedBy: "System",
		Description: fmt.Sprintf("Security audit performed for token: %s", tokenID),
		Timestamp:   time.Now(),
	}
	token.AuditTrail = append(token.AuditTrail, auditLog)

	// Update the ledger with the new audit log
	err = s.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return fmt.Errorf("error updating token in ledger after validation: %w", err)
	}

	return nil
}

// MonitorTransactionSecurity monitors transactions related to SYN1401 tokens and logs any suspicious activities.
func (s *SYN1401Security) MonitorTransactionSecurity(txID string) error {
	// Retrieve the transaction from the ledger
	transaction, err := s.Ledger.GetTransaction(txID)
	if err != nil {
		return fmt.Errorf("error retrieving transaction from ledger: %w", err)
	}

	// Check for suspicious activity (e.g., unusually large transfers, early redemption attempts)
	if transaction.Amount > 1000000 { // Example threshold for suspicious amount
		s.logSuspiciousActivity(transaction)
	}

	// Perform basic encryption validation
	if transaction.IsEncrypted && len(transaction.EncryptedData) == 0 {
		return errors.New("transaction data flagged: expected encrypted content is missing")
	}

	return nil
}

// logSuspiciousActivity logs any suspicious transaction activities related to SYN1401 tokens.
func (s *SYN1401Security) logSuspiciousActivity(transaction *common.Transaction) {
	auditLog := common.AuditLog{
		AuditID:     generateUniqueID(),
		PerformedBy: "System",
		Description: fmt.Sprintf("Suspicious activity detected in transaction: %s", transaction.TxID),
		Timestamp:   time.Now(),
	}

	// Add the audit log to the related token (if applicable)
	if tokenID := transaction.TokenID; tokenID != "" {
		token, err := s.Ledger.GetToken(tokenID)
		if err == nil {
			token.AuditTrail = append(token.AuditTrail, auditLog)
			s.Ledger.UpdateToken(tokenID, token) // Update the token in the ledger
		}
	}
}

// Generate a unique ID for transactions, events, or logs.
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// AES encryption
func (s *SYN1401Security) encryptAES(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// AES decryption
func (s *SYN1401Security) decryptAES(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// hex encoding/decoding for encrypted data (optional)
func encodeToHex(data []byte) string {
	return hex.EncodeToString(data)
}

func decodeFromHex(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}
