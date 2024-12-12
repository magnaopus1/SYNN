package ai_ml_operation

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// AINetworkMaintenance performs routine maintenance tasks on all AI/ML models deployed on the network.
// It includes resource optimization, model update checks, performance diagnostics, and error handling.
func AINetworkMaintenance(l *ledger.Ledger, sc *common.SynnergyConsensus, transactionID string) error {
	log.Printf("AINetworkMaintenance started: TransactionID: %s, Timestamp: %s", transactionID, time.Now().Format(time.RFC3339))

	// Step 1: Fetch all registered models
	models, err := fetchAllModels()
	if err != nil {
		log.Printf("Error fetching models for maintenance: %v", err)
		return fmt.Errorf("failed to fetch models for network maintenance: %w", err)
	}

	// Step 2: Perform maintenance tasks on each model
	for _, model := range models {
		modelID := model.ModelID
		log.Printf("Starting maintenance for ModelID: %s", modelID)

		// Optimize resources
		if err := optimizeResources(model); err != nil {
			logError(fmt.Sprintf("Resource optimization failed for ModelID: %s", modelID), modelID)
			continue
		}

		// Check and apply updates
		if err := checkAndUpdateModel(model); err != nil {
			logError(fmt.Sprintf("Model update check failed for ModelID: %s", modelID), modelID)
			continue
		}

		// Run diagnostics
		if err := runDiagnostics(model); err != nil {
			logError(fmt.Sprintf("Diagnostics failed for ModelID: %s", modelID), modelID)
			continue
		}

		// Create and log maintenance action
		encryptedLog, err := createEncryptedLog("maintenance", modelID)
		if err != nil {
			log.Printf("Error encrypting maintenance log for ModelID: %s: %v", modelID, err)
			continue
		}
		log.Printf("Maintenance log for ModelID %s: %s", modelID, encryptedLog)

		// Record maintenance in the ledger
		if err := l.AiMLMLedger.RecordAiNetworkMaintenance(transactionID); err != nil {
			log.Printf("Error recording maintenance in ledger for ModelID: %s: %v", modelID, err)
			continue
		}

		log.Printf("Maintenance completed for ModelID: %s", modelID)
	}

	log.Printf("AINetworkMaintenance completed: TransactionID: %s", transactionID)
	return nil
}



// ModelRegistry holds all models in the network for easy access.
var (
	modelRegistry = make(map[string]Model) // maps modelID to Model
	mu            sync.RWMutex
)

// GetAllModels retrieves all models from the registry for maintenance or other operations.
func GetAllModels() ([]Model, error) {
	mu.RLock()
	defer mu.RUnlock()

	models := make([]Model, 0, len(modelRegistry))
	for _, model := range modelRegistry {
		models = append(models, model)
	}

	if len(models) == 0 {
		return nil, errors.New("no models available in the registry")
	}
	return models, nil
}

// fetchAllModels retrieves all AI/ML models from the common package for network maintenance.
func fetchAllModels() ([]Model, error) {
	mu.RLock()
	defer mu.RUnlock()

	if len(modelRegistry) == 0 {
		return nil, errors.New("no models available in the registry")
	}

	models := make([]Model, 0, len(modelRegistry))
	for _, model := range modelRegistry {
		models = append(models, model)
	}
	return models, nil
}

// Helper function to optimize resources for a given model
func optimizeResources(model Model) error {
	log.Printf("Optimizing resources for ModelID: %s", model.ModelID)

	if model.NeedsResourceAdjustment() {
		if err := model.AdjustResources(); err != nil {
			return fmt.Errorf("resource adjustment failed for ModelID: %s: %w", model.ModelID, err)
		}
	}
	return nil
}

// Helper function to check for available updates and apply them if present
func checkAndUpdateModel(model Model) error {
	log.Printf("Checking for updates for ModelID: %s", model.ModelID)

	updateAvailable, err := model.CheckForUpdates()
	if err != nil {
		return fmt.Errorf("error checking updates for ModelID: %s: %w", model.ModelID, err)
	}

	if updateAvailable {
		log.Printf("Applying updates for ModelID: %s", model.ModelID)
		if err := model.ApplyUpdate(); err != nil {
			return fmt.Errorf("failed to apply updates for ModelID: %s: %w", model.ModelID, err)
		}
	}
	return nil
}

// Helper function to run diagnostics on a model to ensure optimal performance
func runDiagnostics(model Model) error {
	log.Printf("Running diagnostics for ModelID: %s", model.ModelID)

	if err := model.RunPerformanceTests(); err != nil {
		return fmt.Errorf("performance testing failed for ModelID: %s: %w", model.ModelID, err)
	}

	if err := model.LogDiagnostics(); err != nil {
		return fmt.Errorf("logging diagnostics failed for ModelID: %s: %w", model.ModelID, err)
	}
	return nil
}


// Utility function to log errors during the maintenance process
func logError(message, modelID string) {
	log.Printf("Error: %s, ModelID: %s, Timestamp: %s", message, modelID, time.Now().Format(time.RFC3339))
}


// Utility function to create an encrypted log for each maintenance task
func createEncryptedLog(action, modelID string) (string, error) {
	logData := fmt.Sprintf("Action: %s, ModelID: %s, Timestamp: %s", action, modelID, time.Now().Format(time.RFC3339))
	return AIModelEncrypt([]byte(logData), modelID)
}


// Utility function for encrypting maintenance log data
func AIModelEncrypt(data []byte, key string) (string, error) {
	hashKey := sha256.Sum256([]byte(key)) // Generate a 32-byte AES key
	block, err := aes.NewCipher(hashKey[:])
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM block: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	encryptedData := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(encryptedData), nil
}
