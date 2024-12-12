package cryptography

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"sync"
	"synnergy_network/pkg/common"
)

var macOpsLock sync.Mutex

// CreateMAC: Creates a MAC using HMAC with SHA-256
func CreateMAC(key, data []byte) ([]byte, error) {
    mac := hmac.New(sha256.New, key)
    mac.Write(data)
    result := mac.Sum(nil)
    LogMACOperation("CreateMAC", "MAC created with HMAC-SHA256")
    return result, nil
}

// VerifyMAC: Verifies a MAC using HMAC with SHA-256
func VerifyMAC(key, data, macToVerify []byte) (bool, error) {
    mac := hmac.New(sha256.New, key)
    mac.Write(data)
    computedMAC := mac.Sum(nil)
    isValid := hmac.Equal(computedMAC, macToVerify)
    if !isValid {
        LogMACOperation("VerifyMAC", "MAC verification failed")
        return false, errors.New("MAC verification failed")
    }
    LogMACOperation("VerifyMAC", "MAC verified successfully")
    return true, nil
}

// PSSPaddingEncode: Encodes a message with PSS padding for RSA signatures
func PSSPaddingEncode(message []byte, priv *rsa.PrivateKey) ([]byte, error) {
    hashed := sha256.Sum256(message)
    signature, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA256, hashed[:], &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto})
    if err != nil {
        LogMACOperation("PSSPaddingEncode", "PSS padding encoding failed")
        return nil, err
    }
    LogMACOperation("PSSPaddingEncode", "Message encoded with PSS padding")
    return signature, nil
}

// PSSPaddingDecode: Verifies an RSA PSS-padded signature
func PSSPaddingDecode(message, signature []byte, pub *rsa.PublicKey) error {
    hashed := sha256.Sum256(message)
    err := rsa.VerifyPSS(pub, crypto.SHA256, hashed[:], signature, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto})
    if err != nil {
        LogMACOperation("PSSPaddingDecode", "PSS padding decoding failed")
        return errors.New("PSS verification failed")
    }
    LogMACOperation("PSSPaddingDecode", "PSS-padded signature verified")
    return nil
}

// OFBModeEncrypt: Encrypts data using AES in OFB mode
func OFBModeEncrypt(key, iv, plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        LogMACOperation("OFBModeEncrypt", "AES block cipher initialization failed")
        return nil, err
    }
    ofb := cipher.NewOFB(block, iv)
    ciphertext := make([]byte, len(plaintext))
    ofb.XORKeyStream(ciphertext, plaintext)
    LogMACOperation("OFBModeEncrypt", "Data encrypted using AES in OFB mode")
    return ciphertext, nil
}

// OFBModeDecrypt: Decrypts AES data encrypted in OFB mode
func OFBModeDecrypt(key, iv, ciphertext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        LogMACOperation("OFBModeDecrypt", "AES block cipher initialization failed")
        return nil, err
    }
    ofb := cipher.NewOFB(block, iv)
    plaintext := make([]byte, len(ciphertext))
    ofb.XORKeyStream(plaintext, ciphertext)
    LogMACOperation("OFBModeDecrypt", "Data decrypted using AES in OFB mode")
    return plaintext, nil
}

// LogMACOperation: Logs MAC, PSS, and OFB operations securely
func LogMACOperation(operation, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("MACOperation", encryptedMessage)
}
