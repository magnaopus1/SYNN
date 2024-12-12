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

// StorageManager is responsible for securely storing and retrieving SYN2800 life insurance tokens.
type StorageManager struct {
	mutex sync.Mutex
}

// NewStorageManager creates a new instance of StorageManager.
func NewStorageManager() *StorageManager {
	return &StorageManager{}
}

// StoreToken securely stores the life insurance token in the ledger with encryption.
func (sm *StorageManager) StoreToken(token *common.SYN2800Token) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Encrypt token data before storing in the ledger
	encryptedTokenData, err := sm.encryptTokenData(token)
	if err != nil {
		return fmt.Errorf("failed to encrypt token data: %v", err)
	}

	// Store encrypted token data in the ledger
	if err := ledger.StoreToken(token.TokenID, encryptedTokenData); err != nil {
		return fmt.Errorf("failed to store token in ledger: %v", err)
	}

	log.Printf("Token %s securely stored in ledger", token.TokenID)
	return nil
}

// RetrieveToken securely retrieves and decrypts the life insurance token from the ledger.
func (sm *StorageManager) RetrieveToken(tokenID string) (*common.SYN2800Token, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve encrypted token data from the ledger
	encryptedTokenData, err := ledger.RetrieveToken(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token from ledger: %v", err)
	}

	// Decrypt the token data
	token, err := sm.decryptTokenData(encryptedTokenData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token data: %v", err)
	}

	log.Printf("Token %s successfully retrieved from ledger", tokenID)
	return token, nil
}

// RemoveToken securely removes a token from the ledger.
func (sm *StorageManager) RemoveToken(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Remove the token from the ledger
	if err := ledger.RemoveToken(tokenID); err != nil {
		return fmt.Errorf("failed to remove token from ledger: %v", err)
	}

	log.Printf("Token %s removed from ledger", tokenID)
	return nil
}

// ListAllTokens lists all SYN2800 tokens stored in the ledger.
func (sm *StorageManager) ListAllTokens() ([]*common.SYN2800Token, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve all token IDs stored in the ledger
	tokenIDs, err := ledger.ListAllTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token list: %v", err)
	}

	// Retrieve and decrypt each token
	var tokens []*common.SYN2800Token
	for _, tokenID := range tokenIDs {
		token, err := sm.RetrieveToken(tokenID)
		if err != nil {
			log.Printf("failed to retrieve token %s: %v", tokenID, err)
			continue
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// encryptTokenData encrypts the token data using AES encryption.
func (sm *StorageManager) encryptTokenData(token *common.SYN2800Token) ([]byte, error) {
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

// decryptTokenData decrypts the token data using AES encryption.
func (sm *StorageManager) decryptTokenData(encryptedData []byte) (*common.SYN2800Token, error) {
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

// Helper function to generate an encryption key.
func generateEncryptionKey() []byte {
	// This should be replaced with a secure key generation strategy in a production environment.
	return []byte("your-secure-256-bit-key")
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
