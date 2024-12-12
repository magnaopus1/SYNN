package scalability

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewDataRetrievalSystem initializes the data retrieval system
func NewDataRetrievalSystem(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, cacheTTL time.Duration) *common.DataRetrievalSystem {
	return &common.DataRetrievalSystem{
		Cache:             make(map[string]*common.CacheEntry),
		CacheTTL:          cacheTTL,
		PrefetchKeys:      []string{},
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// CacheRetrieve retrieves data from the cache if available and valid
func (drs *common.DataRetrievalSystem) CacheRetrieve(key string) ([]byte, error) {
	drs.mu.Lock()
	defer drs.mu.Unlock()

	entry, exists := drs.Cache[key]
	if !exists {
		return nil, errors.New("cache miss: key not found")
	}

	if time.Now().After(entry.Expiration) {
		delete(drs.Cache, key)
		return nil, errors.New("cache miss: entry expired")
	}

	fmt.Printf("Cache hit for key: %s\n", key)
	return entry.Data, nil
}

// CacheData caches data with a specified key
func (drs *common.DataRetrievalSystem) CacheData(key string, data []byte) error {
	drs.mu.Lock()
	defer drs.mu.Unlock()

	// Encrypt the data before caching
	encryptedData, err := drs.EncryptionService.EncryptData(data, "encryption-key")
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %v", err)
	}

	expiration := time.Now().Add(drs.CacheTTL)
	drs.Cache[key] = &common.CacheEntry{
		Key:        key,
		Data:       encryptedData,
		Timestamp:  time.Now(),
		Expiration: expiration,
	}

	// Log the caching action in the ledger
	err = drs.Ledger.RecordCacheAction("cache", key, time.Now(), expiration)
	if err != nil {
		return fmt.Errorf("failed to log cache action: %v", err)
	}

	fmt.Printf("Data cached with key: %s (expires at %s)\n", key, expiration)
	return nil
}

// PrefetchCache preloads certain cache entries before they are requested
func (drs *common.DataRetrievalSystem) PrefetchCache(keys []string) error {
	drs.mu.Lock()
	defer drs.mu.Unlock()

	for _, key := range keys {
		// Assume the data source is external, fetch and cache the data
		data, err := drs.fetchFromDataSource(key)
		if err != nil {
			return fmt.Errorf("failed to prefetch key %s: %v", key, err)
		}

		// Cache the prefetched data
		err = drs.CacheData(key, data)
		if err != nil {
			return fmt.Errorf("failed to cache prefetched key %s: %v", key, err)
		}
	}

	// Log the prefetching in the ledger
	err := drs.Ledger.RecordPrefetch(keys, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log prefetch action: %v", err)
	}

	fmt.Printf("Prefetching complete for keys: %v\n", keys)
	return nil
}

// InvalidateCache invalidates a cache entry by key
func (drs *common.DataRetrievalSystem) InvalidateCache(key string) error {
	drs.mu.Lock()
	defer drs.mu.Unlock()

	_, exists := drs.Cache[key]
	if !exists {
		return fmt.Errorf("cache entry not found for key %s", key)
	}

	delete(drs.Cache, key)
	fmt.Printf("Cache entry invalidated for key: %s\n", key)

	// Log the cache invalidation in the ledger
	err := drs.Ledger.RecordCacheInvalidation(key, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log cache invalidation: %v", err)
	}

	return nil
}

// fetchFromDataSource simulates retrieving data from a source (e.g., a database or external API)
func (drs *common.DataRetrievalSystem) fetchFromDataSource(key string) ([]byte, error) {
	// Simulate data retrieval based on key
	data := []byte(fmt.Sprintf("Data for key: %s", key))
	fmt.Printf("Data fetched from source for key: %s\n", key)
	return data, nil
}
