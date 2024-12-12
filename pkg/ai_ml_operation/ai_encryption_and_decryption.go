package ai_ml_operation

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

// AIModelDecrypt securely decrypts AI model data using AES-GCM.
func AIModelDecrypt(encryptedHex string, key string) ([]byte, error) {
	logAction("AIModelDecrypt - Start", fmt.Sprintf("EncryptedHex: %s", encryptedHex))

	// Validate inputs
	if encryptedHex == "" {
		log.Printf("AIModelDecrypt - Validation Failed: EncryptedHex is empty")
		return nil, errors.New("encrypted hex string cannot be empty")
	}
	if key == "" {
		log.Printf("AIModelDecrypt - Validation Failed: Key is empty")
		return nil, errors.New("key cannot be empty")
	}

	// Generate a 32-byte hash key using SHA-256
	hashKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hashKey[:])
	if err != nil {
		log.Printf("AIModelDecrypt - Cipher Block Creation Failed: %v", err)
		return nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	// Create GCM mode for decryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("AIModelDecrypt - GCM Block Creation Failed: %v", err)
		return nil, fmt.Errorf("failed to create GCM block: %w", err)
	}

	// Decode the hex-encoded encrypted string
	encryptedData, err := hex.DecodeString(encryptedHex)
	if err != nil {
		log.Printf("AIModelDecrypt - Hex Decoding Failed: %v", err)
		return nil, fmt.Errorf("invalid hex-encoded string: %w", err)
	}

	// Ensure the encrypted data contains a valid nonce
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		log.Printf("AIModelDecrypt - Validation Failed: Encrypted data length too short")
		return nil, errors.New("invalid encrypted data: insufficient length")
	}

	// Split the encrypted data into nonce and ciphertext
	nonce, cipherText := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Attempt decryption
	decryptedData, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Printf("AIModelDecrypt - Decryption Failed: %v", err)
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	log.Printf("AIModelDecrypt - Success: Decryption completed")
	return decryptedData, nil
}


