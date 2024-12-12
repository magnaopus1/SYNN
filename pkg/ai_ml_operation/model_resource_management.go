package ai_ml_operation

import (
	"encoding/json"
	"errors"
	"synnergy_network/pkg/ledger"
	"time"
)

// ModelResourceAllocation manages allocation of compute resources for AI/ML model execution, with secure logging in the ledger.
func ModelResourceAllocation(modelID string, resourceID string, amount float64, ledgerInstance *ledger.Ledger) error {
	allocationLog := generateResourceLogEntry(modelID, map[string]interface{}{
		"ResourceID": resourceID,
		"Amount":     amount,
	})

	// Convert log to JSON for encryption
	allocationLogData, err := json.Marshal(allocationLog)
	if err != nil {
		return errors.New("failed to encode resource allocation log to JSON")
	}

	// Encrypt the allocation log (for security) even if it’s not stored in this function
	_, err = encryptData(allocationLogData, modelID)
	if err != nil {
		return errors.New("failed to encrypt resource allocation log")
	}

	// Record in the ledger using the resourceID, modelID, and amount for logging purposes
	if err := ledgerInstance.AiMLMLedger.RecordResourceAllocation(resourceID, modelID, amount); err != nil {
		return errors.New("failed to record resource allocation in ledger")
	}

	return nil
}


// ModelStorageAllocate reserves storage for model data, ensuring secure, traceable access through ledger recording.
func ModelStorageAllocate(modelID string, storageSize int, ledgerInstance *ledger.Ledger) error {
	storageLog := map[string]interface{}{
		"ModelID":     modelID,
		"StorageSize": storageSize,
		"Timestamp":   time.Now(),
	}

	// Convert log to JSON for encryption
	storageLogData, err := json.Marshal(storageLog)
	if err != nil {
		return errors.New("failed to encode storage allocation log to JSON")
	}

	// Encrypt the storage log and store it in the ledger
	encryptedLog, err := encryptData(storageLogData, modelID)
	if err != nil {
		return errors.New("failed to encrypt storage allocation log")
	}

	if err := ledgerInstance.AiMLMLedger.RecordCacheData(modelID, "storageAllocation", string(encryptedLog)); err != nil {
		return errors.New("failed to record storage allocation in ledger")
	}

	return nil
}

// ModelResourceRelease releases allocated resources and updates the ledger.
func ModelResourceRelease(modelID string, resourceID string, ledgerInstance *ledger.Ledger) error {
	releaseLog := map[string]interface{}{
		"ModelID":    modelID,
		"ResourceID": resourceID,
		"Action":     "Release",
		"Timestamp":  time.Now(),
	}

	// Convert log to JSON for encryption
	releaseLogData, err := json.Marshal(releaseLog)
	if err != nil {
		return errors.New("failed to encode resource release log to JSON")
	}

	// Encrypt the release log (for security) even though it’s not used in the ledger function
	_, err = encryptData(releaseLogData, modelID)
	if err != nil {
		return errors.New("failed to encrypt resource release log")
	}

	// Record the resource release action in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordResourceRelease(resourceID); err != nil {
		return errors.New("failed to record resource release in ledger")
	}

	return nil
}


// ModelCache stores model data in a secure cache, enabling quick access for high-performance operations.
func ModelCache(modelID string, data interface{}, ledgerInstance *ledger.Ledger) error {
	// Convert data to JSON format before encryption
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return errors.New("failed to encode data to JSON for caching")
	}

	// Encrypt the data and store it in the ledger
	encryptedData, err := encryptData(dataJSON, modelID)
	if err != nil {
		return errors.New("failed to encrypt data for caching")
	}

	if err := ledgerInstance.AiMLMLedger.RecordCacheData(modelID, modelID, string(encryptedData)); err != nil {
		return errors.New("failed to store data in cache")
	}

	return nil
}

// ModelCleanCache removes model data from the cache, logging the action in the ledger.
func ModelCleanCache(modelID string, ledgerInstance *ledger.Ledger) error {
	cacheClearLog := map[string]interface{}{
		"ModelID":   modelID,
		"Action":    "CacheClear",
		"Timestamp": time.Now(),
	}

	// Convert log to JSON for encryption
	cacheClearLogData, err := json.Marshal(cacheClearLog)
	if err != nil {
		return errors.New("failed to encode cache clearing log to JSON")
	}

	// Encrypt the cache clear log (even if unused)
	_, err = encryptData(cacheClearLogData, modelID)
	if err != nil {
		return errors.New("failed to encrypt cache clearing log")
	}

	// Record cache clearing in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordCacheClear(modelID); err != nil {
		return errors.New("failed to clear cache data")
	}

	return nil
}

// ModelShutdown initiates a safe shutdown of a model instance, ensuring the state is stored and resources are freed.
func ModelShutdown(modelID string, ledgerInstance *ledger.Ledger) error {
	shutdownLog := map[string]interface{}{
		"ModelID":   modelID,
		"Action":    "Shutdown",
		"Timestamp": time.Now(),
	}

	// Convert log to JSON for encryption
	shutdownLogData, err := json.Marshal(shutdownLog)
	if err != nil {
		return errors.New("failed to encode shutdown log to JSON")
	}

	// Encrypt the shutdown log (even if unused)
	_, err = encryptData(shutdownLogData, modelID)
	if err != nil {
		return errors.New("failed to encrypt shutdown log")
	}

	// Record shutdown in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordModelShutdown(modelID); err != nil {
		return errors.New("failed to record model shutdown in ledger")
	}

	return nil
}

// ModelRestart restarts a model instance, securely logging the action and verifying it through Synnergy Consensus.
func ModelRestart(modelID string, ledgerInstance *ledger.Ledger) error {
	restartLog := map[string]interface{}{
		"ModelID":   modelID,
		"Action":    "Restart",
		"Timestamp": time.Now(),
	}

	// Convert log to JSON for encryption
	restartLogData, err := json.Marshal(restartLog)
	if err != nil {
		return errors.New("failed to encode restart log to JSON")
	}

	// Encrypt the restart log (even if unused)
	_, err = encryptData(restartLogData, modelID)
	if err != nil {
		return errors.New("failed to encrypt restart log")
	}

	// Record restart in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordModelRestart(modelID); err != nil {
		return errors.New("failed to record model restart in ledger")
	}

	return nil
}


// generateResourceLogEntry creates a standardized log entry for resource allocation and release operations.
func generateResourceLogEntry(modelID string, resources map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"ModelID":   modelID,
		"Resources": resources,
		"Timestamp": time.Now(),
	}
}
