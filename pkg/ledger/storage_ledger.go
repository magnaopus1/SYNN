package ledger

import (
	"errors"
	"fmt"
	"time"
)

// RecordCacheAddition logs the addition of cache data.
func (l *StorageLedger) RecordCacheAddition(key, value string, duration time.Duration) error {
	l.Lock()
	defer l.Unlock()

	// Check if the cache key already exists
	if _, exists := l.CacheRecords[key]; exists {
		return errors.New("cache key already exists")
	}

	// Add the cache record
	cacheRecord := CacheRecord{
		Key:       key,
		Value:     value,
		AddedAt:   time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	l.CacheRecords[key] = cacheRecord
	l.CacheExpirationTimes[key] = cacheRecord.ExpiresAt
	return nil
}

// SetCacheMonitoring enables or disables cache monitoring.
func (l *StorageLedger) SetCacheMonitoring(enabled bool) {
    l.Lock()
    defer l.Unlock()

    l.CacheMonitoring = enabled
}


// RecordCacheUsage records usage for a specific cache ID.
func (l *StorageLedger) RecordCacheUsage(cacheID string, usage int) error {
    l.Lock()
    defer l.Unlock()

    if usage < 0 {
        return fmt.Errorf("usage cannot be negative")
    }
    l.CacheUsageHistory[cacheID] = usage
    return nil
}

// GetCacheUsageHistory retrieves the historical cache usage data.
func (l *StorageLedger) GetCacheUsageHistory() (map[string]int, error) {
    l.Lock()
    defer l.Unlock()

    if len(l.CacheUsageHistory) == 0 {
        return nil, fmt.Errorf("no cache usage history available")
    }
    return l.CacheUsageHistory, nil
}

// RecordCacheRemoval logs the removal of cache data.
func (l *StorageLedger) RecordCacheRemoval(key string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the cache key exists
	if _, exists := l.CacheRecords[key]; !exists {
		return errors.New("cache key not found")
	}

	// Remove the cache record
	delete(l.CacheRecords, key)
	delete(l.CacheExpirationTimes, key)
	return nil
}

// RecordCacheExpiration checks for expired cache and removes it.
func (l *StorageLedger) RecordCacheExpiration() error {
	l.Lock()
	defer l.Unlock()

	// Iterate through the expiration times and remove expired caches
	for key, expiresAt := range l.CacheExpirationTimes {
		if time.Now().After(expiresAt) {
			delete(l.CacheRecords, key)
			delete(l.CacheExpirationTimes, key)
		}
	}
	return nil
}

// RecordCachePurge purges all cache records from the system.
func (l *StorageLedger) RecordCachePurge() error {
	l.Lock()
	defer l.Unlock()

	// Purge all cache records
	l.CacheRecords = make(map[string]CacheRecord)
	l.CacheExpirationTimes = make(map[string]time.Time)
	return nil
}

// RecordFileOperation logs a file operation (read, write, delete).
func (l *StorageLedger) RecordFileOperation(operationID, filePath, action string) error {
	l.Lock()
	defer l.Unlock()

	// Log the file operation
	fileOp := FileOperation{
		OperationID: operationID,
		FilePath:    filePath,
		Action:      action,
		Timestamp:   time.Now(),
	}

	l.FileOperations[operationID] = fileOp
	return nil
}


// RecordProofValidation logs the validation of a space-time proof in the ledger.
func (l *StorageLedger) RecordSTProofValidation(proofID, storageID, validator string) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.SpaceTimeProofs[proofID]; exists {
		return errors.New("proof already exists")
	}

	proof := SpaceTimeProof{
		ProofID:       proofID,
		StorageID:     storageID,
		Validator:     validator,
		ValidationTime: time.Now(),
		Status:        "valid",
	}

	// Store the proof in the ledger
	l.SpaceTimeProofs[proofID] = proof
	return nil
}

// RecordStorageEvent logs a storage-related event in the ledger.
func (l *StorageLedger) RecordStorageEvent(eventID, storageID, action, details string) error {
	l.Lock()
	defer l.Unlock()

	event := StorageEvent{
		EventID:   eventID,
		StorageID: storageID,
		Action:    action,
		Details:   details,
		Timestamp: time.Now(),
	}

	// Store the storage event in the ledger
	l.StorageEvents[eventID] = event
	return nil
}

// RecordProofRevalidation logs the revalidation of a space-time proof.
func (l *StorageLedger) RecordSTProofRevalidation(proofID, validator string) error {
	l.Lock()
	defer l.Unlock()

	proof, exists := l.SpaceTimeProofs[proofID]
	if !exists {
		return errors.New("proof not found")
	}

	// Update proof status and log revalidation
	proof.Status = "revalidated"
	proof.ValidationTime = time.Now()
	l.SpaceTimeProofs[proofID] = proof

	return nil
}

// RecordProofInvalidation logs the invalidation of a space-time proof.
func (l *StorageLedger) RecordSTProofInvalidation(proofID, invalidator, reason string) error {
	l.Lock()
	defer l.Unlock()

	proof, exists := l.SpaceTimeProofs[proofID]
	if !exists {
		return errors.New("proof not found")
	}

	// Invalidate the proof and record the invalidation
	proof.Status = "invalid"
	l.SpaceTimeProofs[proofID] = proof

	// Log the invalidation
	invalidRecord := SpaceTimeProofInvalidation{
		ProofID:    proofID,
		Invalidator: invalidator,
		Reason:     reason,
		Timestamp:  time.Now(),
	}
	l.InvalidProofs[proofID] = invalidRecord

	return nil
}

// RecordProofValidation records the validation of a space-time proof in the ledger.
func (l *StorageLedger) RecordProofValidation(proofID, validator string) error {
    l.Lock()
    defer l.Unlock()

    // Initialize the SpaceTimeProofValidations map if it's nil
    if l.SpaceTimeProofValidations == nil {
        l.SpaceTimeProofValidations = make(map[string][]SpaceTimeProofRecord)
    }

    // Create a new proof record
    proofRecord := SpaceTimeProofRecord{
        ProofID:   proofID,
        Validator: validator,
        Timestamp: time.Now(),
        Status:    "validated",
    }

    // Store the proof validation in the ledger
    l.SpaceTimeProofValidations[proofID] = append(l.SpaceTimeProofValidations[proofID], proofRecord)

    fmt.Printf("Space-time proof validation recorded: ProofID %s by Validator %s\n", proofID, validator)
    return nil
}


// RecordProofInvalidation records the invalidation of a space-time proof in the ledger.
func (l *StorageLedger) RecordProofInvalidation(proofID, validator string) error {
    l.Lock()
    defer l.Unlock()

    // Initialize the SpaceTimeProofInvalidations map if it's nil
    if l.SpaceTimeProofInvalidations == nil {
        l.SpaceTimeProofInvalidations = make(map[string][]SpaceTimeProofRecord)
    }

    // Create a new proof record
    proofRecord := SpaceTimeProofRecord{
        ProofID:   proofID,
        Validator: validator,
        Timestamp: time.Now(),
        Status:    "invalidated",
    }

    // Store the proof invalidation in the ledger
    l.SpaceTimeProofInvalidations[proofID] = append(l.SpaceTimeProofInvalidations[proofID], proofRecord)

    fmt.Printf("Space-time proof invalidation recorded: ProofID %s by Validator %s\n", proofID, validator)
    return nil
}

// RecordProofRevalidation records the revalidation of a space-time proof in the ledger.
func (l *StorageLedger) RecordProofRevalidation(proofID, validator string) error {
    l.Lock()
    defer l.Unlock()

    // Initialize the SpaceTimeProofRevalidations map if it's nil
    if l.SpaceTimeProofRevalidations == nil {
        l.SpaceTimeProofRevalidations = make(map[string][]SpaceTimeProofRecord)
    }

    // Create a new proof record
    proofRecord := SpaceTimeProofRecord{
        ProofID:   proofID,
        Validator: validator,
        Timestamp: time.Now(),
        Status:    "revalidated",
    }

    // Store the proof revalidation in the ledger
    l.SpaceTimeProofRevalidations[proofID] = append(l.SpaceTimeProofRevalidations[proofID], proofRecord)

    fmt.Printf("Space-time proof revalidation recorded: ProofID %s by Validator %s\n", proofID, validator)
    return nil
}

