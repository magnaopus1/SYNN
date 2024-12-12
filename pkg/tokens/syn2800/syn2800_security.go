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

// SecurityManager handles all security-related operations for the SYN2800 tokens.
type SecurityManager struct {
	mutex sync.Mutex
}

// NewSecurityManager creates a new instance of SecurityManager.
func NewSecurityManager() *SecurityManager {
	return &SecurityManager{}
}

// EncryptTokenData encrypts the life insurance token data using AES encryption.
func (sm *SecurityManager) EncryptTokenData(token *common.SYN2800Token) ([]byte, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

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

	tokenData := serializeTokenData(token)
	return gcm.Seal(nonce, nonce, tokenData, nil), nil
}

// DecryptTokenData decrypts the life insurance token data using AES encryption.
func (sm *SecurityManager) DecryptTokenData(encryptedData []byte) (*common.SYN2800Token, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

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

	return deserializeTokenData(decryptedData), nil
}

// ValidateTokenOwnership verifies that the provided user is the legitimate owner of the life insurance token.
func (sm *SecurityManager) ValidateTokenOwnership(tokenID, userID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the token from the ledger
	encryptedTokenData, err := ledger.RetrieveToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token: %v", err)
	}

	// Decrypt the token data
	token, err := sm.DecryptTokenData(encryptedTokenData)
	if err != nil {
		return fmt.Errorf("failed to decrypt token data: %v", err)
	}

	// Validate ownership
	if token.PolicyHolder != userID {
		return fmt.Errorf("ownership validation failed: user %s is not the owner of token %s", userID, tokenID)
	}

	return nil
}

// GenerateTokenSignature generates a digital signature for a life insurance token using HMAC with SHA-256.
func (sm *SecurityManager) GenerateTokenSignature(token *common.SYN2800Token) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	h := hmac.New(sha256.New, []byte("secure-shared-secret-key"))
	tokenData := serializeTokenData(token)
	_, err := h.Write(tokenData)
	if err != nil {
		return "", err
	}

	signature := hex.EncodeToString(h.Sum(nil))
	return signature, nil
}

// VerifyTokenSignature verifies the integrity of a token using its HMAC signature.
func (sm *SecurityManager) VerifyTokenSignature(token *common.SYN2800Token, signature string) (bool, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	expectedSignature, err := sm.GenerateTokenSignature(token)
	if err != nil {
		return false, err
	}

	return hmac.Equal([]byte(signature), []byte(expectedSignature)), nil
}

// LogSecurityEvent logs security-related events like unauthorized access attempts or token tampering.
func (sm *SecurityManager) LogSecurityEvent(tokenID string, eventType string, details string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	eventLog := common.SecurityEventLog{
		EventType: eventType,
		Details:   details,
		Timestamp: time.Now(),
		TokenID:   tokenID,
	}

	// Log the event in the token's security logs in the ledger
	token, err := sm.retrieveAndDecryptToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token for logging event: %v", err)
	}
	token.SecurityLogs = append(token.SecurityLogs, eventLog)

	// Encrypt and store the updated token back in the ledger
	encryptedTokenData, err := sm.EncryptTokenData(token)
	if err != nil {
		return fmt.Errorf("failed to encrypt token data after logging event: %v", err)
	}
	if err := ledger.StoreToken(tokenID, encryptedTokenData); err != nil {
		return fmt.Errorf("failed to store updated token in ledger after logging event: %v", err)
	}

	log.Printf("Security event logged for Token ID: %s, Event: %s", tokenID, eventType)
	return nil
}

// generateEncryptionKey generates a 256-bit encryption key.
func generateEncryptionKey() []byte {
	// You may want to generate this key using a secure key management system in production
	return []byte("your-secure-256-bit-key")
}

// Helper function to retrieve and decrypt the token from the ledger.
func (sm *SecurityManager) retrieveAndDecryptToken(tokenID string) (*common.SYN2800Token, error) {
	encryptedData, err := ledger.RetrieveToken(tokenID)
	if err != nil {
		return nil, err
	}
	return sm.DecryptTokenData(encryptedData)
}

// Helper function to serialize token data.
func serializeTokenData(token *common.SYN2800Token) []byte {
	data, err := json.Marshal(token)
	if err != nil {
		log.Fatalf("failed to serialize token data: %v", err)
	}
	return data
}

// Helper function to deserialize token data after decryption.
func deserializeTokenData(data []byte) *common.SYN2800Token {
	var token common.SYN2800Token
	if err := json.Unmarshal(data, &token); err != nil {
		log.Fatalf("failed to deserialize token data: %v", err)
	}
	return &token
}
