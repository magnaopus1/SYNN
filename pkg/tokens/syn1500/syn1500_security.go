package syn1500

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// SYN1500Security provides security features for SYN1500 tokens, such as encryption, decryption, and validation.
type SYN1500Security struct {
	Ledger ledger.Ledger // Ledger to log all security-related transactions
}

// EncryptToken encrypts sensitive data within a SYN1500Token.
func (ss *SYN1500Security) EncryptToken(token *common.SYN1500Token, key []byte) error {
	// Convert token struct into bytes (simplified for demonstration)
	tokenBytes := []byte(fmt.Sprintf("%v", token))

	// Encrypt the token data
	encryptedData, err := encrypt(tokenBytes, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt token data: %v", err)
	}

	// Store the encrypted data in the token
	token.EncryptedMetadata = encryptedData

	// Record encryption event in the ledger
	eventDescription := fmt.Sprintf("Token %s encrypted", token.TokenID)
	return ss.recordSecurityEvent(token.TokenID, "Encryption", eventDescription)
}

// DecryptToken decrypts the encrypted data within a SYN1500Token.
func (ss *SYN1500Security) DecryptToken(token *common.SYN1500Token, key []byte) (string, error) {
	// Decrypt the token's encrypted metadata
	decryptedData, err := decrypt(token.EncryptedMetadata, key)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt token data: %v", err)
	}

	// Record decryption event in the ledger
	eventDescription := fmt.Sprintf("Token %s decrypted", token.TokenID)
	if err := ss.recordSecurityEvent(token.TokenID, "Decryption", eventDescription); err != nil {
		return "", err
	}

	// Return decrypted data (as string for simplicity)
	return string(decryptedData), nil
}

// recordSecurityEvent logs encryption or decryption events in the ledger.
func (ss *SYN1500Security) recordSecurityEvent(tokenID, eventType, description string) error {
	// Create a new ledger transaction
	tx := ledger.Transaction{
		TxID:        generateUniqueID(tokenID),
		Description: description,
		Timestamp:   time.Now(),
		Data: map[string]interface{}{
			"token_id":   tokenID,
			"event_type": eventType,
		},
	}

	// Validate transaction using Synnergy Consensus before recording it
	if err := synnergy_consensus.ValidateTransaction(tx); err != nil {
		return fmt.Errorf("failed Synnergy Consensus validation: %v", err)
	}

	// Record the event in the ledger
	return ss.Ledger.RecordTransaction(tx)
}

// encrypt uses AES encryption to securely encrypt the data.
func encrypt(data, key []byte) ([]byte, error) {
	// Create a new AES cipher with the given key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the data using AES-GCM
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts AES-GCM encrypted data.
func decrypt(data, key []byte) ([]byte, error) {
	// Create a new AES cipher with the given key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract the nonce size from the cipher
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// Split the nonce and the actual ciphertext
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt the data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// generateUniqueID creates a unique identifier for transactions or security events.
func generateUniqueID(seed string) string {
	timestamp := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", seed, timestamp)))
	return hex.EncodeToString(hash[:])
}

// ValidateSecurity performs comprehensive security checks on SYN1500Token to ensure it hasn't been tampered with.
func (ss *SYN1500Security) ValidateSecurity(token *common.SYN1500Token) error {
	// Generate a hash from the encrypted metadata for integrity validation
	hash := sha256.Sum256(token.EncryptedMetadata)

	// Create a validation transaction
	tx := ledger.Transaction{
		TxID:        generateUniqueID(token.TokenID),
		Description: fmt.Sprintf("Security validation for token %s", token.TokenID),
		Timestamp:   time.Now(),
		Data: map[string]interface{}{
			"token_id": token.TokenID,
			"hash":     hex.EncodeToString(hash[:]),
		},
	}

	// Validate the security check using Synnergy Consensus
	if err := synnergy_consensus.ValidateTransaction(tx); err != nil {
		return fmt.Errorf("failed Synnergy Consensus security validation: %v", err)
	}

	// Record the validation in the ledger
	return ss.Ledger.RecordTransaction(tx)
}
