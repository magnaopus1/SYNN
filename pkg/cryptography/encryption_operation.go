package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"sync"

	"golang.org/x/crypto/blake2b"
)

var encryptionLock sync.Mutex

// AESEncrypt: Encrypts data using AES-GCM
func AESEncrypt(key, plaintext []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        LogEncryptionOperation("AESEncrypt", "AES encryption failed")
        return "", errors.New("AES encryption failed")
    }
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonce := make([]byte, aesgcm.NonceSize())
    _, err = rand.Read(nonce)
    if err != nil {
        return "", err
    }
    ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
    encryptedData := base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))
    LogEncryptionOperation("AESEncrypt", "Data encrypted with AES-GCM")
    return encryptedData, nil
}

// AESDecrypt: Decrypts AES-GCM encrypted data
func AESDecrypt(key []byte, ciphertext string) ([]byte, error) {
    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return nil, err
    }
    block, err := aes.NewCipher(key)
    if err != nil {
        LogEncryptionOperation("AESDecrypt", "AES decryption failed")
        return nil, errors.New("AES decryption failed")
    }
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    nonce, ciphertext := data[:aesgcm.NonceSize()], data[aesgcm.NonceSize():]
    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }
    LogEncryptionOperation("AESDecrypt", "Data decrypted with AES-GCM")
    return plaintext, nil
}

// RSAEncrypt: Encrypts data using RSA
func RSAEncrypt(publicKey *rsa.PublicKey, plaintext []byte) ([]byte, error) {
    ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, nil)
    if err != nil {
        LogEncryptionOperation("RSAEncrypt", "RSA encryption failed")
        return nil, errors.New("RSA encryption failed")
    }
    LogEncryptionOperation("RSAEncrypt", "Data encrypted with RSA")
    return ciphertext, nil
}

// RSADecrypt: Decrypts RSA encrypted data
func RSADecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
    plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
    if err != nil {
        LogEncryptionOperation("RSADecrypt", "RSA decryption failed")
        return nil, errors.New("RSA decryption failed")
    }
    LogEncryptionOperation("RSADecrypt", "Data decrypted with RSA")
    return plaintext, nil
}

// Blake2bHash: Hashes data using Blake2b
func Blake2bHash(data []byte) ([]byte, error) {
    hash := blake2b.Sum256(data)
    LogEncryptionOperation("Blake2bHash", "Data hashed with Blake2b")
    return hash[:], nil
}

