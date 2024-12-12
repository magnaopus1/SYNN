package syn4900

import (
	"errors"
	"sync"
	"time"
)

// TokenSecurityManager handles all security-related operations for SYN5000 tokens.
type TokenSecurityManager struct {
	mu               sync.RWMutex
	keys             map[string]*rsa.PrivateKey // In-memory storage of private keys for tokens.
	tokenSecureHashes map[string]string         // Storage of secure hashes for each token.
	ledger           *ledger.SecurityLedger     // Security ledger to track security-related events.
	encryptor        *encryption.Encryptor      // Encryption module for securing token data.
}

// NewTokenSecurityManager creates a new instance of TokenSecurityManager.
func NewTokenSecurityManager(ledger *ledger.SecurityLedger, encryptor *encryption.Encryptor) *TokenSecurityManager {
	return &TokenSecurityManager{
		keys:              make(map[string]*rsa.PrivateKey),
		tokenSecureHashes: make(map[string]string),
		ledger:            ledger,
		encryptor:         encryptor,
	}
}

// GenerateKeys generates an RSA key pair for the specified token and stores it securely.
func (tsm *TokenSecurityManager) GenerateKeys(tokenID string) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Check if keys already exist for this token.
	if _, exists := tsm.keys[tokenID]; exists {
		return errors.New("keys already exist for this token")
	}

	// Generate a new RSA key pair.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Store the private key securely.
	tsm.keys[tokenID] = privateKey

	// Generate a secure hash for the token and store it.
	tokenHash := tsm.generateSecureHash(tokenID)
	tsm.tokenSecureHashes[tokenID] = tokenHash

	// Log the key generation event in the ledger.
	tsm.ledger.LogSecurityEvent("KeyGenerated", tokenID, time.Now(), tokenHash)

	return nil
}

// EncryptData encrypts the data associated with a token using the stored private key.
func (tsm *TokenSecurityManager) EncryptData(tokenID string, data []byte) (string, error) {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()

	// Get the private key for the token.
	privateKey, exists := tsm.keys[tokenID]
	if !exists {
		return "", errors.New("no keys found for the specified token")
	}

	// Encrypt the data using the public key.
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &privateKey.PublicKey, data, nil)
	if err != nil {
		return "", err
	}

	// Return the encrypted data as a hex string.
	return hex.EncodeToString(encryptedData), nil
}

// DecryptData decrypts encrypted data for the specified token using the private key.
func (tsm *TokenSecurityManager) DecryptData(tokenID, encryptedData string) ([]byte, error) {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()

	// Get the private key for the token.
	privateKey, exists := tsm.keys[tokenID]
	if !exists {
		return nil, errors.New("no keys found for the specified token")
	}

	// Decode the encrypted data from hex string.
	encryptedBytes, err := hex.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	// Decrypt the data using the private key.
	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedBytes, nil)
	if err != nil {
		return nil, err
	}

	// Return the decrypted data.
	return decryptedData, nil
}

// ValidateToken validates the secure hash of a token to ensure its integrity.
func (tsm *TokenSecurityManager) ValidateToken(tokenID string) (bool, error) {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()

	// Retrieve the stored secure hash for the token.
	storedHash, exists := tsm.tokenSecureHashes[tokenID]
	if !exists {
		return false, errors.New("no secure hash found for the specified token")
	}

	// Generate a new secure hash based on the current state of the token.
	currentHash := tsm.generateSecureHash(tokenID)

	// Compare the stored hash with the newly generated hash.
	if storedHash != currentHash {
		return false, nil
	}

	// If the hashes match, the token is valid.
	return true, nil
}

// generateSecureHash generates a secure hash for a token based on its ID and associated data.
func (tsm *TokenSecurityManager) generateSecureHash(tokenID string) string {
	hash := sha256.New()
	hash.Write([]byte(tokenID))
	hash.Write([]byte(time.Now().String())) // Using current time as part of the hash.
	return hex.EncodeToString(hash.Sum(nil))
}

// RevokeKeys revokes the keys for a specific token, making further transactions invalid.
func (tsm *TokenSecurityManager) RevokeKeys(tokenID string) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Check if the keys exist for this token.
	_, exists := tsm.keys[tokenID]
	if !exists {
		return errors.New("no keys found for the specified token")
	}

	// Remove the keys and hash from memory.
	delete(tsm.keys, tokenID)
	delete(tsm.tokenSecureHashes, tokenID)

	// Log the key revocation event in the ledger.
	tsm.ledger.LogSecurityEvent("KeyRevoked", tokenID, time.Now(), "")

	return nil
}

// LogTokenEvent logs a security event related to the token.
func (tsm *TokenSecurityManager) LogTokenEvent(eventType, tokenID, eventData string) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Generate a secure hash for the event data.
	eventHash := tsm.generateSecureHash(tokenID + eventData)

	// Log the event in the security ledger.
	return tsm.ledger.LogSecurityEvent(eventType, tokenID, time.Now(), eventHash)
}
