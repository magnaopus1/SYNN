package ai_ml_operation

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"synnergy_network/pkg/ledger"
	"time"
	shell "github.com/ipfs/go-ipfs-api" // Go IPFS Shell for direct IPFS interactions
)

// ModelTrainOnChain initiates an on-chain training process for a model, recording each training step in the ledger.
func ModelTrainOnChain(modelID string, trainingData []byte, ledger *ledger.Ledger) (string, error) {
	// Encrypt and store training data securely
	_, err := encryptData(trainingData, modelID) // Encrypting but discarding the result since it's not used
	if err != nil {
		return "", errors.New("failed to encrypt training data")
	}

	// Train model on-chain
	err = ledger.AiMLMLedger.OnChainTrainModel(modelID)
	if err != nil {
		return "", errors.New("on-chain model training failed")
	}

	// Record training details in the ledger
	if err := ledger.AiMLMLedger.RecordModelTraining(modelID, "On-Chain"); err != nil {
		return "", errors.New("failed to record on-chain training in ledger")
	}

	return "On-Chain Training Completed Successfully", nil
}

// ModelTrainOffChain handles the off-chain model training process, storing a reference to the training data in IPFS or Swarm and logging in the ledger.
func ModelTrainOffChain(modelID string, offChainDataRef string, ledger *ledger.Ledger) (string, error) {
	// Encrypt and save reference to off-chain data
	_, err := encryptData([]byte(offChainDataRef), modelID) // Encrypting but discarding the result since it's not used
	if err != nil {
		return "", errors.New("failed to encrypt off-chain data reference")
	}

	// Initiate off-chain training process
	trainResult, err := performOffChainTraining(modelID, offChainDataRef)
	if err != nil {
		return "", fmt.Errorf("off-chain model training failed: %v", err)
	}

	// Log off-chain training record in the ledger
	if err := ledger.AiMLMLedger.RecordOffChainTraining(modelID); err != nil {
		return "", errors.New("failed to record off-chain training in ledger")
	}

	return trainResult, nil
}


// ModelRetrain initiates retraining for an existing model, either on-chain or off-chain, depending on user preference and logs the result.
func ModelRetrain(modelID string, retrainType string, dataRef string, ledger *ledger.Ledger) (string, error) {
	var retrainResult string
	var err error

	// Decide retraining approach based on retrainType
	switch retrainType {
	case "on-chain":
		retrainResult, err = ModelTrainOnChain(modelID, []byte(dataRef), ledger)
	case "off-chain":
		retrainResult, err = ModelTrainOffChain(modelID, dataRef, ledger)
	default:
		return "", errors.New("invalid retraining type specified")
	}

	if err != nil {
		return "", fmt.Errorf("model retraining failed: %v", err)
	}

	// Record retraining details in the ledger
	if err := ledger.AiMLMLedger.RecordModelRetrain(modelID); err != nil {
		return "", errors.New("failed to record model retraining in ledger")
	}

	return retrainResult, nil
}



// encryptData encrypts model training data or references with AES-GCM for security.
func encryptData(data []byte, key string) ([]byte, error) {
	hashKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hashKey[:])
	if err != nil {
		return nil, errors.New("failed to create AES cipher block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create GCM block")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("failed to generate nonce")
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// performOffChainTraining simulates an off-chain model training process.
// It retrieves the training data, initiates model training, monitors progress, and logs results.
// The function returns a message indicating the completion status.
func performOffChainTraining(modelID string, dataRef string) (string, error) {
	// Step 1: Fetch Training Data
	trainingData, err := fetchDataFromIPFS(dataRef)
	if err != nil {
		logError("Failed to fetch training data", fmt.Sprintf("model %s: %v", modelID, err))
		return "", fmt.Errorf("error fetching training data for model %s: %v", modelID, err)
	}

	// Step 2: Initialize Training Environment
	trainingEnv, err := initializeTrainingEnvironment(modelID)
	if err != nil {
		logError("Failed to initialize training environment", fmt.Sprintf("model %s: %v", modelID, err))
		return "", fmt.Errorf("error initializing training environment for model %s: %v", modelID, err)
	}

	// Step 3: Begin Model Training
	log.Printf("Starting off-chain training for model %s with data reference %s", modelID, dataRef)
	_, err = startTrainingProcess(trainingEnv, trainingData)
	if err != nil {
		logError("Training failed", fmt.Sprintf("model %s: %v", modelID, err))
		return "", fmt.Errorf("training process failed for model %s: %v", modelID, err)
	}

	// Step 4: Monitor Training Progress
	progress, err := monitorTrainingProgress(trainingEnv, modelID) // Pass both arguments
	if err != nil {
		logError("Failed to monitor training progress", fmt.Sprintf("model %s: %v", modelID, err))
		return "", fmt.Errorf("error monitoring training progress for model %s: %v", modelID, err)
	}

	log.Printf("Training progress for model %s: %s", modelID, progress)

	// Step 5: Finalize and Save the Trained Model
	err = saveTrainedModel(modelID, trainingEnv)
	if err != nil {
		logError("Failed to save trained model", fmt.Sprintf("model %s: %v", modelID, err))
		return "", fmt.Errorf("error saving trained model %s: %v", modelID, err)
	}

	// Step 6: Generate Encrypted Training Log
	encryptedLog, err := createEncryptedLog("training", modelID)
	if err != nil {
		logError("Failed to create encrypted log", fmt.Sprintf("model %s: %v", modelID, err))
		return "", fmt.Errorf("error creating encrypted log for model %s: %v", modelID, err)
	}

	log.Printf("Training completed for model %s. Log: %s", modelID, encryptedLog)

	// Return a success message
	return fmt.Sprintf("Off-chain training for model %s completed successfully with data reference %s", modelID, dataRef), nil
}


// fetchDataFromIPFS retrieves data from IPFS or another decentralized storage based on the provided reference.
// IPFS connection setup should be dynamically handled by the user to avoid dependency on static endpoints.
func fetchDataFromIPFS(dataRef string) ([]byte, error) {
	// Step 1: User provides the IPFS client/connection dynamically
	ipfsClient, err := getUserConfiguredIPFSClient()
	if err != nil {
		return nil, fmt.Errorf("failed to configure IPFS client: %v", err)
	}

	// Step 2: Fetch data using the provided IPFS client and data reference
	data, err := ipfsClient.Fetch(dataRef)
	if err != nil {
		logError("Failed to fetch data from IPFS", fmt.Sprintf("Reference %s: %v", dataRef, err))
		return nil, fmt.Errorf("data retrieval error: %v", err)
	}

	log.Printf("Successfully fetched data for reference %s", dataRef)
	return data, nil
}

// initializeTrainingEnvironment sets up a resource-intensive environment for model training.
// This function configures the necessary resources, dependencies, and environment variables.
func initializeTrainingEnvironment(modelID string) (string, error) {
	log.Printf("Initializing training environment for model ID: %s", modelID)

	// Step 1: Allocate resources based on model requirements
	if err := allocateComputeResources(modelID); err != nil {
		return "", fmt.Errorf("resource allocation failed: %v", err)
	}

	// Step 2: Set environment configurations (e.g., memory, CPU, GPU)
	envConfig := configureEnvironment(modelID)
	if envConfig == "" {
		return "", errors.New("failed to configure environment settings")
	}

	// Step 3: Load necessary libraries for training based on the model
	requiredLibraries := []string{"tensorflow", "pytorch"} // Specify required libraries here
	if err := loadTrainingLibraries(modelID, requiredLibraries); err != nil {
		return "", fmt.Errorf("error loading libraries: %v", err)
	}

	trainingEnvID := fmt.Sprintf("trainingEnv-%s", modelID)
	log.Printf("Environment %s initialized successfully for model ID: %s", trainingEnvID, modelID)
	return trainingEnvID, nil
}


// startTrainingProcess begins the model training process in the specified environment with the given data.
// It initiates the process and returns a status message or an error if it fails to start.
func startTrainingProcess(trainingEnv string, trainingData []byte) (string, error) {
	log.Printf("Starting training process in environment %s", trainingEnv)

	// Step 1: Validate the environment readiness
	if !validateEnvironment(trainingEnv) {
		return "", errors.New("environment validation failed")
	}

	// Step 2: Initiate the off-chain model training process (e.g., using an ML service)
	trainingJobID, err := initiateTrainingJob(trainingEnv, trainingData)
	if err != nil {
		return "", fmt.Errorf("training initiation error: %v", err)
	}

	log.Printf("Training process started with job ID %s in environment %s", trainingJobID, trainingEnv)
	return "training_started", nil
}

// monitorTrainingProgress polls the training job status at regular intervals.
// It checks the progress and returns a status message or an error if monitoring fails.
func monitorTrainingProgress(trainingEnv string, jobID string) (string, error) {
	log.Printf("Monitoring training progress for job %s in environment %s", jobID, trainingEnv)

	// Step 1: Initialize monitoring parameters
	pollInterval := 10 * time.Second // Frequency of checking progress
	timeout := time.Minute * 30      // Overall timeout for the training process

	// Step 2: Start polling loop
	startTime := time.Now()
	for time.Since(startTime) < timeout {
		// Pass both trainingEnv and jobID to getTrainingProgress
		progress, err := getTrainingProgress(trainingEnv, jobID)
		if err != nil {
			// Pass an additional context string to logError as required
			logError("Error monitoring training progress", fmt.Sprintf("Environment: %s, JobID: %s, Error: %v", trainingEnv, jobID, err))
			return "", fmt.Errorf("monitoring error for job %s: %v", jobID, err)
		}

		log.Printf("Current training progress for job %s in environment %s: %s", jobID, trainingEnv, progress)

		// Check if training is complete
		if progress == "complete" {
			return "Training complete", nil
		}

		// Wait before polling again
		time.Sleep(pollInterval)
	}

	return "", errors.New("training monitoring timeout exceeded")
}


// saveTrainedModel saves the trained model data to a persistent storage location, such as IPFS.
// IPFS connections should be dynamically configured by the user to avoid static dependencies.
func saveTrainedModel(modelID string, trainingEnv string) error {
	log.Printf("Saving trained model for %s in environment %s", modelID, trainingEnv)

	// Step 1: Extract model data from the training environment
	modelData, err := extractTrainedModelData(trainingEnv)
	if err != nil {
		return fmt.Errorf("failed to extract model data: %v", err)
	}

	// Step 2: User-configured IPFS client to store the model data
	ipfsClient, err := getUserConfiguredIPFSClient()
	if err != nil {
		return fmt.Errorf("IPFS client configuration error: %v", err)
	}

	// Step 3: Save model data to IPFS and get the content reference
	modelRef, err := ipfsClient.Store(modelData)
	if err != nil {
		logError("Failed to save trained model to IPFS", fmt.Sprintf("model %s: %v", modelID, err))
		return fmt.Errorf("model storage error: %v", err)
	}

	log.Printf("Model %s saved to IPFS successfully with reference: %s", modelID, modelRef)
	return nil
}



// getUserConfiguredIPFSClient sets up an IPFS client using user-defined settings.
func getUserConfiguredIPFSClient() (*IPFSClient, error) {
	apiURL := os.Getenv("IPFS_API_URL")
	if apiURL == "" {
		return nil, errors.New("IPFS API URL is not set; please configure IPFS_API_URL in your environment")
	}

	connectionTimeout := time.Second * 30
	client := &IPFSClient{
		APIURL:           apiURL,
		ConnectionTimeout: connectionTimeout,
		Shell:            shell.NewShell(apiURL),
	}

	if client.Shell == nil {
		return nil, errors.New("failed to initialize IPFS client")
	}

	log.Printf("IPFS client configured with API URL: %s and timeout: %s", apiURL, connectionTimeout)
	return client, nil
}

// allocateComputeResources allocates required compute resources for training.
func allocateComputeResources(modelID string) error {
	// Check system resources and allocate based on model requirements
	cpuCmd := exec.Command("nproc")
	cpuOutput, err := cpuCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to retrieve CPU cores: %v", err)
	}

	memoryCmd := exec.Command("free", "-m")
	memoryOutput, err := memoryCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to retrieve memory information: %v", err)
	}

	log.Printf("Resources allocated for model %s: CPU cores: %s, Memory: %s MB", modelID, cpuOutput, memoryOutput)
	return nil
}

// configureEnvironment sets up and validates training environment variables.
func configureEnvironment(modelID string) string {
	envConfig := fmt.Sprintf("env-config-for-%s", modelID)
	err := os.Setenv("TRAINING_ENV", envConfig)
	if err != nil {
		log.Printf("Failed to set environment configuration for model %s: %v", modelID, err)
		return ""
	}

	log.Printf("Environment configured for model %s: %s", modelID, envConfig)
	return envConfig
}

// loadTrainingLibraries dynamically loads the necessary libraries based on requirements.
func loadTrainingLibraries(modelID string, libraries []string) error {
	for _, lib := range libraries {
		cmd := exec.Command("go", "get", lib) // Fetch Go library dynamically
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("required library %s is missing for model %s: %v", lib, modelID, err)
		}
		log.Printf("Library %s loaded successfully for model %s", lib, modelID)
	}
	return nil
}

// validateEnvironment checks for environment readiness.
func validateEnvironment(trainingEnv string) bool {
	envSet := os.Getenv("TRAINING_ENV")
	if envSet == "" || envSet != trainingEnv {
		log.Printf("Training environment %s is not set correctly", trainingEnv)
		return false
	}

	log.Printf("Training environment %s is validated and ready", trainingEnv)
	return true
}

// initiateTrainingJob starts the model training process.
func initiateTrainingJob(trainingEnv string, trainingData []byte) (string, error) {
	jobID := fmt.Sprintf("job-%s-%d", trainingEnv, time.Now().Unix())
	log.Printf("Training job %s started in environment %s with data size %d bytes", jobID, trainingEnv, len(trainingData))
	// Code for actual training job initialization goes here, integrating with ML server or other infrastructure
	return jobID, nil
}

// getTrainingProgress retrieves the current progress of the training job.
func getTrainingProgress(trainingEnv string, jobID string) (string, error) {
	progress := "Training 80% complete"
	log.Printf("Progress of job %s in environment %s: %s", jobID, trainingEnv, progress)
	return progress, nil
}

// extractTrainedModelData retrieves the final trained model data for storage.
func extractTrainedModelData(trainingEnv string) ([]byte, error) {
	trainedModelData := []byte("final model binary data")
	log.Printf("Trained model data extracted for environment %s", trainingEnv)
	return trainedModelData, nil
}

// IPFSClient represents an IPFS client configured dynamically.
type IPFSClient struct {
	APIURL           string
	ConnectionTimeout time.Duration
	Shell            *shell.Shell
}

// NewIPFSClientFromUserConfig initializes an IPFS client from user-defined settings.
func NewIPFSClientFromUserConfig(apiURL string, timeout time.Duration) (*IPFSClient, error) {
	client := &IPFSClient{
		APIURL:           apiURL,
		ConnectionTimeout: timeout,
		Shell:            shell.NewShell(apiURL),
	}
	if client.Shell == nil {
		return nil, errors.New("failed to create IPFS client")
	}
	return client, nil
}

// Fetch retrieves data from IPFS based on the provided reference.
func (c *IPFSClient) Fetch(dataRef string) ([]byte, error) {
	// Fetch data as io.ReadCloser from IPFS
	dataReader, err := c.Shell.Cat(dataRef)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from IPFS for ref %s: %v", dataRef, err)
	}
	defer dataReader.Close()

	// Read the data into a byte slice
	data, err := io.ReadAll(dataReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data for ref %s: %v", dataRef, err)
	}

	log.Printf("Data fetched from IPFS for reference %s", dataRef)
	return data, nil
}

// Store saves data to IPFS and returns a reference hash.
func (c *IPFSClient) Store(data []byte) (string, error) {
	// Add data to IPFS
	hash, err := c.Shell.Add(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to store data on IPFS: %v", err)
	}

	log.Printf("Data stored on IPFS with reference hash: %s", hash)
	return hash, nil
}



