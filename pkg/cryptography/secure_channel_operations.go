package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"synnergy_network/pkg/common"

	"golang.org/x/crypto/hkdf"
)

var secureChannelLock sync.Mutex

// SecureChannel represents a secure communication channel with symmetric encryption
type SecureChannel struct {
    sessionKey []byte
    aesGCM     cipher.AEAD
    isActive   bool
}

// EstablishSecureChannel: Establishes a secure channel by deriving a session key using HKDF
func EstablishSecureChannel(sharedSecret, salt []byte) (*SecureChannel, error) {
    if len(sharedSecret) == 0 {
        LogSecureChannelOperation("EstablishSecureChannel", "Invalid shared secret provided")
        return nil, errors.New("shared secret must not be empty")
    }

    // Derive session key using HKDF with SHA-256
    hkdf := hkdf.New(sha256.New, sharedSecret, salt, nil)
    sessionKey := make([]byte, 32) // 256-bit key for AES-GCM
    if _, err := hkdf.Read(sessionKey); err != nil {
        LogSecureChannelOperation("EstablishSecureChannel", "Failed to derive session key")
        return nil, err
    }

    // Initialize AES-GCM cipher with the derived session key
    block, err := aes.NewCipher(sessionKey)
    if err != nil {
        LogSecureChannelOperation("EstablishSecureChannel", "AES cipher initialization failed")
        return nil, err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        LogSecureChannelOperation("EstablishSecureChannel", "AES-GCM initialization failed")
        return nil, err
    }

    secureChannel := &SecureChannel{
        sessionKey: sessionKey,
        aesGCM:     aesGCM,
        isActive:   true,
    }

    LogSecureChannelOperation("EstablishSecureChannel", "Secure channel established successfully")
    return secureChannel, nil
}

// CloseSecureChannel: Closes a secure channel and securely erases the session key
func (sc *SecureChannel) CloseSecureChannel() error {
    if !sc.isActive {
        LogSecureChannelOperation("CloseSecureChannel", "Secure channel is already closed")
        return errors.New("secure channel is already closed")
    }

    // Securely erase session key
    for i := range sc.sessionKey {
        sc.sessionKey[i] = 0
    }
    sc.isActive = false

    LogSecureChannelOperation("CloseSecureChannel", "Secure channel closed and session key erased")
    return nil
}

// EncryptMessage: Encrypts a message using AES-GCM in the secure channel
func (sc *SecureChannel) EncryptMessage(plaintext []byte) (string, error) {
    if !sc.isActive {
        LogSecureChannelOperation("EncryptMessage", "Attempted to encrypt with closed secure channel")
        return "", errors.New("secure channel is closed")
    }

    nonce := make([]byte, sc.aesGCM.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        LogSecureChannelOperation("EncryptMessage", "Nonce generation failed")
        return "", err
    }

    ciphertext := sc.aesGCM.Seal(nonce, nonce, plaintext, nil)
    LogSecureChannelOperation("EncryptMessage", "Message encrypted successfully")
    return hex.EncodeToString(ciphertext), nil
}

// DecryptMessage: Decrypts a message using AES-GCM in the secure channel
func (sc *SecureChannel) DecryptMessage(ciphertextHex string) ([]byte, error) {
    if !sc.isActive {
        LogSecureChannelOperation("DecryptMessage", "Attempted to decrypt with closed secure channel")
        return nil, errors.New("secure channel is closed")
    }

    ciphertext, err := hex.DecodeString(ciphertextHex)
    if err != nil {
        LogSecureChannelOperation("DecryptMessage", "Failed to decode ciphertext")
        return nil, err
    }

    nonceSize := sc.aesGCM.NonceSize()
    if len(ciphertext) < nonceSize {
        LogSecureChannelOperation("DecryptMessage", "Ciphertext too short for nonce")
        return nil, errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := sc.aesGCM.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        LogSecureChannelOperation("DecryptMessage", "Decryption failed")
        return nil, err
    }

    LogSecureChannelOperation("DecryptMessage", "Message decrypted successfully")
    return plaintext, nil
}

// LogSecureChannelOperation: Logs secure channel operations with encryption
func LogSecureChannelOperation(operation, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("SecureChannelOperation", encryptedMessage)
}
