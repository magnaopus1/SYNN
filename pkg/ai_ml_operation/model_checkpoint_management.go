package ai_ml_operation

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// ModelSaveCheckpoint saves a checkpoint of a model's state, securely encrypts it, and records it in the ledger for retrieval.
func ModelSaveCheckpoint(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, transactionID, nodeID string, modelState interface{}) (string, error) {

	// Convert modelState to []byte for encryption
	modelStateBytes, err := json.Marshal(modelState)
	if err != nil {
		return "", errors.New("failed to serialize model state")
	}

	// Encrypt model state for secure storage
	encryptedState, err := encryptData(modelStateBytes, modelID)
	if err != nil {
		return "", errors.New("failed to encrypt model checkpoint")
	}

	// Generate a version number for the checkpoint
	checkpointVersion := generateCheckpointVersion(l, modelID)

	// Record encrypted checkpoint in the ledger with metadata
	if err := l.AiMLMLedger.RecordModelCheckpoint(modelID, checkpointVersion, string(encryptedState)); err != nil {
		return "", errors.New("failed to record model checkpoint in ledger")
	}

	return fmt.Sprintf("%s-%d", modelID, checkpointVersion), nil
}

// generateCheckpointVersion creates a new version number for the checkpoint based on the ledger state.
func generateCheckpointVersion(l *ledger.Ledger, modelID string) int {
	// Retrieve the last checkpoint version for this model, if exists, and increment it
	if checkpoint, exists := l.AiMLMLedger.Checkpoints[modelID]; exists {
		return checkpoint.Version + 1
	}
	// Start with version 1 if no checkpoint exists
	return 1
}


// ModelLoadCheckpoint retrieves and decrypts a model checkpoint, allowing the model to revert to a saved state.
func ModelLoadCheckpoint(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, transactionID, nodeID string, checkpointID string) (interface{}, error) {

	// Retrieve the encrypted checkpoint data from the ledger
	checkpoint, err := l.AiMLMLedger.GetModelCheckpoint(modelID)
	if err != nil {
		return nil, errors.New("failed to retrieve model checkpoint from ledger")
	}

	// Decrypt the checkpoint data for model loading
	decryptedState, err := decryptData([]byte(checkpoint.DataHash), modelID)
	if err != nil {
		return nil, errors.New("failed to decrypt model checkpoint")
	}

	// Ensure decryptedState is explicitly a byte slice
	decryptedBytes, ok := decryptedState.([]byte)
	if !ok {
		return nil, errors.New("decrypted data is not in expected format")
	}

	// Deserialize the decrypted state back into the original format
	var modelState interface{}
	if err := json.Unmarshal(decryptedBytes, &modelState); err != nil {
		return nil, errors.New("failed to deserialize model checkpoint data")
	}

	return modelState, nil
}



// decryptData decrypts AES-GCM encrypted data for secure access
func decryptData(encryptedData []byte, key string) (interface{}, error) {
	hashKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hashKey[:])
	if err != nil {
		return nil, errors.New("failed to create AES cipher block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create GCM block")
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	data, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt data")
	}

	return data, nil
}

// generateCheckpointID generates a unique checkpoint identifier based on model ID and timestamp
func generateCheckpointID(modelID string) string {
	timestamp := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", modelID, timestamp)))
	return hex.EncodeToString(hash[:])
}
