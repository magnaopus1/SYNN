package cryptography

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"sync"
	"synnergy_network/pkg/common"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/scrypt"
)

var keyMgmtLock sync.Mutex

// RSAKeyGen: Generates an RSA key pair
func RSAKeyGen(bits int) (*rsa.PrivateKey, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        LogKeyMgmtOperation("RSAKeyGen", "RSA key generation failed")
        return nil, err
    }
    LogKeyMgmtOperation("RSAKeyGen", "RSA key pair generated")
    return privateKey, nil
}

// RSAKeyExport: Exports an RSA private key as PEM
func RSAKeyExport(privateKey *rsa.PrivateKey) (string, error) {
    pemData := pem.EncodeToMemory(
        &pem.Block{
            Type:  "RSA PRIVATE KEY",
            Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
        },
    )
    LogKeyMgmtOperation("RSAKeyExport", "RSA private key exported as PEM")
    return string(pemData), nil
}

// RSAKeyImport: Imports an RSA private key from PEM format
func RSAKeyImport(pemData string) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode([]byte(pemData))
    if block == nil || block.Type != "RSA PRIVATE KEY" {
        LogKeyMgmtOperation("RSAKeyImport", "RSA PEM import failed")
        return nil, errors.New("failed to decode PEM block containing RSA private key")
    }
    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    LogKeyMgmtOperation("RSAKeyImport", "RSA private key imported from PEM")
    return privateKey, nil
}

// ED25519Sign: Signs a message using ED25519
func ED25519Sign(privateKey ed25519.PrivateKey, message []byte) ([]byte, error) {
    signature := ed25519.Sign(privateKey, message)
    LogKeyMgmtOperation("ED25519Sign", "Message signed using ED25519")
    return signature, nil
}

// ED25519Verify: Verifies an ED25519 signature
func ED25519Verify(publicKey ed25519.PublicKey, message, signature []byte) bool {
    isValid := ed25519.Verify(publicKey, message, signature)
    LogKeyMgmtOperation("ED25519Verify", "ED25519 signature verified")
    return isValid
}

// X25519KeyExchange: Performs X25519 key exchange
func X25519KeyExchange(privateKey, peerPublicKey []byte) ([]byte, error) {
    sharedKey, err := curve25519.X25519(privateKey, peerPublicKey)
    if err != nil {
        LogKeyMgmtOperation("X25519KeyExchange", "X25519 key exchange failed")
        return nil, err
    }
    LogKeyMgmtOperation("X25519KeyExchange", "X25519 key exchange completed")
    return sharedKey, nil
}

// ECDHKeyExchange: Performs ECDH key exchange
func ECDHKeyExchange(privateKey *ecdsa.PrivateKey, peerPublicKey *ecdsa.PublicKey) ([]byte, error) {
    sharedX, _ := privateKey.ScalarMult(peerPublicKey.X, peerPublicKey.Y, privateKey.D.Bytes())
    sharedKey := sharedX.Bytes()
    LogKeyMgmtOperation("ECDHKeyExchange", "ECDH key exchange completed")
    return sharedKey, nil
}

// DeriveKeyFromPassword: Derives a key from a password using scrypt
func DeriveKeyFromPassword(password, salt []byte, keyLen int) ([]byte, error) {
    key, err := scrypt.Key(password, salt, 32768, 8, 1, keyLen)
    if err != nil {
        LogKeyMgmtOperation("DeriveKeyFromPassword", "Password-based key derivation failed")
        return nil, err
    }
    LogKeyMgmtOperation("DeriveKeyFromPassword", "Key derived from password using scrypt")
    return key, nil
}

// GenerateRandomBytes: Generates random bytes
func GenerateRandomBytes(length int) ([]byte, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        LogKeyMgmtOperation("GenerateRandomBytes", "Random byte generation failed")
        return nil, err
    }
    LogKeyMgmtOperation("GenerateRandomBytes", "Random bytes generated")
    return bytes, nil
}

// NonceGenerate: Generates a random nonce
func NonceGenerate(length int) ([]byte, error) {
    nonce, err := GenerateRandomBytes(length)
    if err != nil {
        LogKeyMgmtOperation("NonceGenerate", "Nonce generation failed")
        return nil, err
    }
    LogKeyMgmtOperation("NonceGenerate", "Nonce generated")
    return nonce, nil
}

// KeyExpansion: Expands a key using HMAC
func KeyExpansion(key, info []byte, length int) ([]byte, error) {
    mac := hmac.New(sha256.New, key)
    mac.Write(info)
    expandedKey := mac.Sum(nil)[:length]
    LogKeyMgmtOperation("KeyExpansion", "Key expanded")
    return expandedKey, nil
}

// MaskKey: Masks a key with a mask
func MaskKey(key, mask []byte) ([]byte, error) {
    if len(key) != len(mask) {
        LogKeyMgmtOperation("MaskKey", "Key masking failed due to length mismatch")
        return nil, errors.New("key and mask must have the same length")
    }
    maskedKey := make([]byte, len(key))
    for i := range key {
        maskedKey[i] = key[i] ^ mask[i]
    }
    LogKeyMgmtOperation("MaskKey", "Key masked")
    return maskedKey, nil
}



// LogKeyMgmtOperation: Logs key management operations with encryption
func LogKeyMgmtOperation(operation, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("KeyMgmtOperation", encryptedMessage)
}
