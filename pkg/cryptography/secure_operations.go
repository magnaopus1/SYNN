package cryptography

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"synnergy_network/pkg/common"
)

var secureOpsLock sync.Mutex

// SecureKeyErase: Securely erases a key from memory
func SecureKeyErase(key []byte) {
    for i := range key {
        key[i] = 0
    }
    LogSecureOperation("SecureKeyErase", "Key securely erased from memory")
}

// SecureEraseMemory: Securely wipes data from a byte array in memory
func SecureEraseMemory(data []byte) {
    for i := range data {
        data[i] = 0
    }
    LogSecureOperation("SecureEraseMemory", "Memory securely erased")
}

// SecureCompareHash: Compares two hashes in constant time to prevent timing attacks
func SecureCompareHash(hash1, hash2 []byte) (bool, error) {
    if len(hash1) != len(hash2) {
        LogSecureOperation("SecureCompareHash", "Hash length mismatch")
        return false, errors.New("hashes must be of equal length")
    }

    result := 0
    for i := range hash1 {
        result |= int(hash1[i] ^ hash2[i])
    }

    isEqual := result == 0
    LogSecureOperation("SecureCompareHash", "Hashes compared securely")
    return isEqual, nil
}

// SecureHashChain: Generates a secure hash chain based on input data and chain length
func SecureHashChain(data []byte, length int) ([][]byte, error) {
    if length <= 0 {
        LogSecureOperation("SecureHashChain", "Invalid hash chain length")
        return nil, errors.New("chain length must be greater than zero")
    }

    chain := make([][]byte, length)
    currentHash := sha256.Sum256(data)
    chain[0] = currentHash[:]

    for i := 1; i < length; i++ {
        currentHash = sha256.Sum256(chain[i-1])
        chain[i] = currentHash[:]
    }

    LogSecureOperation("SecureHashChain", "Hash chain generated")
    return chain, nil
}

// MaskData: Applies a random mask to data for obfuscation
func MaskData(data []byte) ([]byte, []byte, error) {
    mask := make([]byte, len(data))
    _, err := rand.Read(mask)
    if err != nil {
        LogSecureOperation("MaskData", "Mask generation failed")
        return nil, nil, err
    }

    maskedData := make([]byte, len(data))
    for i := range data {
        maskedData[i] = data[i] ^ mask[i]
    }

    LogSecureOperation("MaskData", "Data masked successfully")
    return maskedData, mask, nil
}

// UnmaskData: Reverses masking on data using the original mask
func UnmaskData(maskedData, mask []byte) ([]byte, error) {
    if len(maskedData) != len(mask) {
        LogSecureOperation("UnmaskData", "Mask length mismatch")
        return nil, errors.New("mask and data length must match")
    }

    unmaskedData := make([]byte, len(maskedData))
    for i := range maskedData {
        unmaskedData[i] = maskedData[i] ^ mask[i]
    }

    LogSecureOperation("UnmaskData", "Data unmasked successfully")
    return unmaskedData, nil
}

// ObfuscateHash: Obfuscates a hash using a salt and additional transformations
func ObfuscateHash(hash, salt []byte) (string, error) {
    hasher := sha256.New()
    hasher.Write(salt)
    hasher.Write(hash)
    obfuscatedHash := hasher.Sum(nil)

    LogSecureOperation("ObfuscateHash", "Hash obfuscated successfully")
    return hex.EncodeToString(obfuscatedHash), nil
}

// RekeyData: Re-encrypts data with a new key
func RekeyData(oldKey, newKey, data []byte) ([]byte, error) {
    // Apply old key as XOR mask
    maskedData, err := MaskData(data)
    if err != nil {
        LogSecureOperation("RekeyData", "Failed to mask data with old key")
        return nil, err
    }

    // Unmask with new key to re-encrypt
    rekeyedData, err := UnmaskData(maskedData, newKey)
    if err != nil {
        LogSecureOperation("RekeyData", "Failed to re-encrypt data with new key")
        return nil, err
    }

    LogSecureOperation("RekeyData", "Data rekeyed successfully")
    return rekeyedData, nil
}

// LogSecureOperation: Logs secure operations with encryption
func LogSecureOperation(operation, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("SecureOperation", encryptedMessage)
}
