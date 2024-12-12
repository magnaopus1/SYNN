package cryptography

import (
	"errors"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

var keyCache = make(map[string][]byte) // Cache for storing cryptographic keys
var cacheLock sync.Mutex

// ClearKeyCache: Clears all keys from the key cache securely
func ClearKeyCache() error {
    cacheLock.Lock()
    defer cacheLock.Unlock()

    // Ensure there is data in the cache before attempting to clear
    if len(keyCache) == 0 {
        LogClearOperation("ClearKeyCache", "Key cache is already empty")
        return errors.New("key cache is already empty")
    }

    // Clear the key cache securely
    for keyID := range keyCache {
        delete(keyCache, keyID)
    }

    LogClearOperation("ClearKeyCache", "Key cache cleared successfully")
    return nil
}

// Helper Functions

// LogClearOperation: Logs cache clearing operations with encryption
func LogClearOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("ClearOperation", encryptedMessage)
}
