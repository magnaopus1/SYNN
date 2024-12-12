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
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// AIImageProcessing processes image data using off-chain AI models hosted on IPFS or Swarm.
// The function loads the model, validates the transaction, processes the image, and stores the result in the ledger.
func AIImageProcessing(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, nodeID, mode string, imageData []byte, transactionID string) (string, error) {
	startTime := time.Now()
	log.Printf("AIImageProcessing started | ModelID: %s, NodeID: %s, TransactionID: %s, Timestamp: %s", modelID, nodeID, transactionID, startTime)

	// Step 1: Load the model from IPFS
	model, err := loadModelFromIPFS(modelID)
	if err != nil {
		log.Printf("Error: Failed to load model | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("failed to load image processing model: %w", err)
	}

	// Step 2: Validate the mode
	if !model.SupportsMode(mode) {
		log.Printf("Error: Unsupported mode | ModelID: %s, Mode: %s", modelID, mode)
		return "", fmt.Errorf("mode %s is not supported by the model", mode)
	}

	// Step 3: Process the image data
	result, err := model.ProcessImage(imageData, mode)
	if err != nil {
		log.Printf("Error: Image processing failed | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("image processing failed: %w", err)
	}

	// Step 4: Encrypt the result
	encryptedResult, err := encryptResult(result, modelID)
	if err != nil {
		log.Printf("Error: Result encryption failed | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("failed to encrypt image processing result: %w", err)
	}

	// Step 5: Record the processing in the ledger
	err = l.AiMLMLedger.RecordProcessing(transactionID, modelID, nodeID, "image_processing", encryptedResult)
	if err != nil {
		log.Printf("Error: Ledger recording failed | ModelID: %s, TransactionID: %s, Error: %s", modelID, transactionID, err.Error())
		return "", fmt.Errorf("error recording image processing in ledger: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("AIImageProcessing completed | ModelID: %s, TransactionID: %s, Duration: %s", modelID, transactionID, duration)
	return encryptedResult, nil
}


// AIAudioProcessing processes audio data using off-chain AI models hosted on IPFS or Swarm.
// This function loads the model, validates the transaction, processes the audio, and stores the result in the ledger.
func AIAudioProcessing(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, nodeID, mode string, audioData []byte, transactionID string) (string, error) {
	startTime := time.Now()
	log.Printf("AIAudioProcessing started | ModelID: %s, NodeID: %s, TransactionID: %s, Timestamp: %s", modelID, nodeID, transactionID, startTime)

	// Step 1: Load the model from IPFS
	model, err := loadModelFromIPFS(modelID)
	if err != nil {
		log.Printf("Error: Failed to load model | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("failed to load audio processing model: %w", err)
	}

	// Step 2: Validate the mode
	if !model.SupportsMode(mode) {
		log.Printf("Error: Unsupported mode | ModelID: %s, Mode: %s", modelID, mode)
		return "", fmt.Errorf("mode %s is not supported by the model", mode)
	}

	// Step 3: Process the audio data
	result, err := model.ProcessAudio(audioData, mode)
	if err != nil {
		log.Printf("Error: Audio processing failed | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("audio processing failed: %w", err)
	}

	// Step 4: Encrypt the result
	encryptedResult, err := encryptResult(result, modelID)
	if err != nil {
		log.Printf("Error: Result encryption failed | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("failed to encrypt audio processing result: %w", err)
	}

	// Step 5: Record the processing in the ledger
	err = l.AiMLMLedger.RecordProcessing(transactionID, modelID, nodeID, "audio_processing", encryptedResult)
	if err != nil {
		log.Printf("Error: Ledger recording failed | ModelID: %s, TransactionID: %s, Error: %s", modelID, transactionID, err.Error())
		return "", fmt.Errorf("error recording audio processing in ledger: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("AIAudioProcessing completed | ModelID: %s, TransactionID: %s, Duration: %s", modelID, transactionID, duration)
	return encryptedResult, nil
}


// AITextAnalysis performs text analysis on the input text using off-chain AI models hosted on IPFS or Swarm.
// This function validates the transaction, analyzes the text, and logs the analysis result in the ledger.
func AITextAnalysis(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, nodeID, textData, transactionID string) (string, error) {
	startTime := time.Now()
	log.Printf("AITextAnalysis started | ModelID: %s, NodeID: %s, TransactionID: %s, Timestamp: %s", modelID, nodeID, transactionID, startTime)

	// Step 1: Load the model from IPFS
	model, err := loadModelFromIPFS(modelID)
	if err != nil {
		log.Printf("Error: Failed to load model | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("failed to load text analysis model: %w", err)
	}

	// Step 2: Analyze the text data
	result, err := model.AnalyzeText(textData, nodeID)
	if err != nil {
		log.Printf("Error: Text analysis failed | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("text analysis failed: %w", err)
	}

	// Step 3: Encrypt the analysis result
	encryptedResult, err := encryptResult(result, modelID)
	if err != nil {
		log.Printf("Error: Result encryption failed | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("failed to encrypt text analysis result: %w", err)
	}

	// Step 4: Record the processing in the ledger
	err = l.AiMLMLedger.RecordProcessing(transactionID, modelID, nodeID, "text_analysis", encryptedResult)
	if err != nil {
		log.Printf("Error: Ledger recording failed | ModelID: %s, TransactionID: %s, Error: %s", modelID, transactionID, err.Error())
		return "", fmt.Errorf("error recording text analysis in ledger: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("AITextAnalysis completed | ModelID: %s, TransactionID: %s, Duration: %s", modelID, transactionID, duration)
	return encryptedResult, nil
}

// AIVideoAnalysis processes video data using off-chain AI models hosted on IPFS or Swarm.
// This function validates the transaction, analyzes the video, and logs the results in the ledger.
func AIVideoAnalysis(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, nodeID string, videoData []byte, transactionID string) (string, error) {
	startTime := time.Now()
	log.Printf("AIVideoAnalysis started | ModelID: %s, NodeID: %s, TransactionID: %s, Timestamp: %s", modelID, nodeID, transactionID, startTime)

	// Step 1: Load the model from IPFS
	model, err := loadModelFromIPFS(modelID)
	if err != nil {
		log.Printf("Error: Failed to load model | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("failed to load video analysis model: %w", err)
	}

	// Step 2: Analyze the video data
	result, err := model.AnalyzeVideo(videoData, nodeID)
	if err != nil {
		log.Printf("Error: Video analysis failed | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("video analysis failed: %w", err)
	}

	// Step 3: Encrypt the analysis result
	encryptedResult, err := encryptResult(result, modelID)
	if err != nil {
		log.Printf("Error: Result encryption failed | ModelID: %s, Error: %s", modelID, err.Error())
		return "", fmt.Errorf("failed to encrypt video analysis result: %w", err)
	}

	// Step 4: Record the processing in the ledger
	err = l.AiMLMLedger.RecordProcessing(transactionID, modelID, nodeID, "video_analysis", encryptedResult)
	if err != nil {
		log.Printf("Error: Ledger recording failed | ModelID: %s, TransactionID: %s, Error: %s", modelID, transactionID, err.Error())
		return "", fmt.Errorf("error recording video analysis in ledger: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("AIVideoAnalysis completed | ModelID: %s, TransactionID: %s, Duration: %s", modelID, transactionID, duration)
	return encryptedResult, nil
}





// Utility function to load AI models from IPFS/Swarm based on modelID
func loadModelFromIPFS(modelID string) (Model, error) {
	log.Printf("Loading model from IPFS | ModelID: %s", modelID)

	// Retrieve the model index to get the IPFS/Swarm link
	modelIndex, err := GetModelIndex(modelID)
	if err != nil {
		log.Printf("Error: Model index retrieval failed | ModelID: %s, Error: %s", modelID, err.Error())
		return Model{}, fmt.Errorf("failed to retrieve model index: %w", err)
	}

	// Fetch the model's raw data
	modelData, err := fetchModelFromStorage(modelIndex.IPFSLink)
	if err != nil {
		log.Printf("Error: Failed to fetch model data | ModelID: %s, Link: %s, Error: %s", modelID, modelIndex.IPFSLink, err.Error())
		return Model{}, fmt.Errorf("failed to fetch model data: %w", err)
	}

	// Deserialize the data into a structured Model object
	model, err := deserializeModelData(modelData)
	if err != nil {
		log.Printf("Error: Failed to deserialize model data | ModelID: %s, Error: %s", modelID, err.Error())
		return Model{}, fmt.Errorf("failed to deserialize model data: %w", err)
	}

	log.Printf("Model loaded successfully | ModelID: %s", modelID)
	return model, nil
}




// Utility function to encrypt the processing results before storing in the ledger
func encryptResult(data string, key string) (string, error) {
	log.Printf("Encrypting result | Key: %s", key)

	// Generate a 256-bit hash-based encryption key
	hashKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hashKey[:])
	if err != nil {
		log.Printf("Error: Failed to create AES cipher block | Error: %s", err.Error())
		return "", fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	// Create a GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("Error: Failed to create GCM block | Error: %s", err.Error())
		return "", fmt.Errorf("failed to create GCM block: %w", err)
	}

	// Generate a unique nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("Error: Failed to generate nonce | Error: %s", err.Error())
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	encryptedData := gcm.Seal(nonce, nonce, []byte(data), nil)
	log.Printf("Encryption successful")
	return hex.EncodeToString(encryptedData), nil
}
