package storage

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "time"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/common"
)

// NewCacheManager initializes the cache manager with a maximum size and TTL
func NewCacheManager(ledgerInstance *ledger.Ledger, maxSize int, ttl time.Duration) *CacheManager {
    return &CacheManager{
        CacheEntries:   make(map[string]*CacheEntry),
        LedgerInstance: ledgerInstance,
        MaxCacheSize:   maxSize,
        DefaultTTL:     ttl,
    }
}

// AddToCache adds a new entry to the cache, encrypts the data, and records the activity in the ledger
func (cm *CacheManager) AddToCache(key string, data string) error {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()

    // Check if the cache size limit has been reached
    if len(cm.CacheEntries) >= cm.MaxCacheSize {
        return errors.New("cache limit reached, cannot add new entry")
    }

    // Encrypt the data using the Encryption instance from the common package
    encryptedData, err := common.EncryptData("AES", []byte(data), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt cache data: %v", err)
    }

    // Create the cache entry
    cacheEntry := &CacheEntry{
        Key:          key,
        Data:         encryptedData,
        CachedAt:     time.Now(),
        ExpiresAt:    time.Now().Add(cm.DefaultTTL),
        LastAccessed: time.Now(),
    }

    // Add to cache
    cm.CacheEntries[key] = cacheEntry
    fmt.Printf("Added data to cache with key %s.\n", key)

    // Log the cache addition to the ledger
    return cm.LedgerInstance.RecordCacheAddition(key, encryptedData)
}

// GetFromCache retrieves data from the cache if it exists and is not expired, then decrypts it
func (cm *CacheManager) GetFromCache(key string) (string, error) {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()

    entry, exists := cm.CacheEntries[key]
    if !exists {
        return "", errors.New("cache entry not found")
    }

    // Check if the cache entry is expired
    if time.Now().After(entry.Expiration) {
        delete(cm.CacheEntries, key)
        fmt.Printf("Cache entry %s expired and removed.\n", key)
        return "", errors.New("cache entry expired")
    }

    entry.LastAccessed = time.Now()

    decryptedData, err := common.DecryptData(entry.Data, common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt cache data: %v", err)
    }

    fmt.Printf("Retrieved data from cache with key %s.\n", key)
    return decryptedData, nil
}

// RemoveFromCache removes an entry from the cache and records the removal in the ledger
func (cm *CacheManager) RemoveFromCache(key string) error {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()

    if _, exists := cm.CacheEntries[key]; !exists {
        return errors.New("cache entry not found")
    }

    delete(cm.CacheEntries, key)

    fmt.Printf("Removed data from cache with key %s.\n", key)
    return cm.LedgerInstance.RecordCacheRemoval(key)
}

// CleanupExpiredEntries removes expired entries from the cache and logs them in the ledger
func (cm *CacheManager) CleanupExpiredEntries() {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()

    for key, entry := range cm.CacheEntries {
        if time.Now().After(entry.Expiration) {
            delete(cm.CacheEntries, key)
            fmt.Printf("Cache entry %s expired and removed during cleanup.\n", key)
            cm.LedgerInstance.RecordCacheExpiration(key)
        }
    }
}

// PurgeCache clears all entries from the cache and logs the action
func (cm *CacheManager) PurgeCache() {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()

    cm.CacheEntries = make(map[string]*CacheEntry)
    fmt.Println("Cache purged completely.")
    cm.LedgerInstance.RecordCachePurge()
}

// HashCacheKey generates a hash of a cache key for secure storage
func (cm *CacheManager) HashCacheKey(key string) string {
    hash := sha256.New()
    hash.Write([]byte(key))
    return hex.EncodeToString(hash.Sum(nil))
}

// CacheSize returns the current size of the cache
func (cm *CacheManager) CacheSize() int {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()

    return len(cm.CacheEntries)
}
