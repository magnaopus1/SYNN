package advanced_data_and_resource_management

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)



func logAction(action, detail string) {
	log.Printf("[%s] Action: %s, Detail: %s", time.Now().Format(time.RFC3339), action, detail)
}

// NewAENManager initializes a new AENManager.
func NewAENManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *AENManager {
	logAction("InitializeAENManager", "New instance created")
	return &AENManager{
		Tasks:             make(map[string]*AENTask),
		TaskQueue:         []string{},
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
		mu:                sync.Mutex{},
	}
}

// CreateTask generates a new AEN task and adds it to the queue.
func (aen *AENManager) CreateTask(taskType string, payload string, executorNode string) (*AENTask, error) {
	aen.mu.Lock()
	defer aen.mu.Unlock()

	logAction("CreateTask", fmt.Sprintf("TaskType: %s, ExecutorNode: %s", taskType, executorNode))

	// Validate input
	if taskType == "" || payload == "" || executorNode == "" {
		return nil, errors.New("taskType, payload, and executorNode cannot be empty")
	}

	// Create a unique task ID
	taskID := generateUniqueID()

	// Encrypt the payload
	encryptedPayload, err := aen.EncryptionService.EncryptData("AES", []byte(payload), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt task payload: %w", err)
	}

	encryptedPayloadStr := base64.StdEncoding.EncodeToString(encryptedPayload)

	// Initialize the task
	newTask := &AENTask{
		TaskID:           taskID,
		TaskType:         taskType,
		Payload:          payload,
		TaskStatus:       "Pending",
		ExecutorNode:     executorNode,
		CreationTime:     time.Now(),
		EncryptedPayload: encryptedPayloadStr,
	}

	// Add task to queue and map
	aen.TaskQueue = append(aen.TaskQueue, taskID)
	aen.Tasks[taskID] = newTask

	// Log task creation in the ledger
	_, err = aen.Ledger.UtilityLedger.RecordTaskCreation(taskID, map[string]interface{}{
		"TaskType":     taskType,
		"ExecutorNode": executorNode,
		"CreationTime": newTask.CreationTime.Format(time.RFC3339),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to log task creation in the ledger: %w", err)
	}

	log.Printf("Task %s created successfully by node %s", taskID, executorNode)
	return newTask, nil
}

// ExecuteTask processes a task from the task queue with real-world logic.
func (aen *AENManager) ExecuteTask() error {
	aen.mu.Lock()
	defer aen.mu.Unlock()

	logAction("ExecuteTask", "Starting task execution")

	if len(aen.TaskQueue) == 0 {
		return errors.New("no tasks in the queue")
	}

	// Get the next task from the queue
	taskID := aen.TaskQueue[0]
	aen.TaskQueue = aen.TaskQueue[1:]

	// Retrieve the task
	task, exists := aen.Tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	task.TaskStatus = "Processing"

	// Step 1: Fetch external data
	externalData, err := fetchExternalData(taskID)
	if err != nil {
		task.TaskStatus = "Failed"
		_ = aen.Ledger.UtilityLedger.RecordTaskFailure(taskID, task.ExecutorNode, err.Error())
		return fmt.Errorf("failed to fetch external data for task %s: %w", taskID, err)
	}

	// Step 2: Process task payload with external data
	processedResult, err := processTaskData(task.Payload, externalData)
	if err != nil {
		task.TaskStatus = "Failed"
		_ = aen.Ledger.UtilityLedger.RecordTaskFailure(taskID, task.ExecutorNode, err.Error())
		return fmt.Errorf("failed to process task %s: %w", taskID, err)
	}

	// Mark task as completed
	task.TaskStatus = "Completed"
	task.CompletionTime = time.Now()

	// Log task completion in the ledger
	err = aen.Ledger.UtilityLedger.RecordTaskCompletion(taskID, task.ExecutorNode, task.CompletionTime.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to log task completion in the ledger: %w", err)
	}

	log.Printf("Task %s executed successfully by node %s. Result: %s", taskID, task.ExecutorNode, processedResult)
	return nil
}

// MonitorTaskStatus checks the status of a specific task.
func (aen *AENManager) MonitorTaskStatus(taskID string) (string, error) {
	aen.mu.Lock()
	defer aen.mu.Unlock()

	logAction("MonitorTaskStatus", fmt.Sprintf("TaskID: %s", taskID))

	task, exists := aen.Tasks[taskID]
	if !exists {
		return "", fmt.Errorf("task %s not found", taskID)
	}

	log.Printf("Task %s status: %s", taskID, task.TaskStatus)
	return task.TaskStatus, nil
}

// Ledger Integration
func (aen *AENManager) integrateWithLedger(task *AENTask) error {
    // Create a TaskRecord instance
    taskRecord := ledger.TaskRecord{
        ID:          task.TaskID,
        Creator:     task.CreatorNode,      // Assuming the task struct has a CreatorNode field
        Assignee:    task.ExecutorNode,     // The node executing the task
        Status:      task.TaskStatus,
        CreatedAt:   task.CreationTime,
        CompletedAt: &task.CompletionTime,  // Assuming task.CompletionTime exists and is a time.Time value
        TaskDetails: map[string]interface{}{
            "Payload": task.Payload, // Add additional details you want to track
        },
    }

    // Call the Ledger's RecordTask method with the TaskRecord
    return aen.Ledger.UtilityLedger.RecordTask(taskRecord)
}

// fetchExternalData retrieves external data from a trusted source, such as an API or a database.
// The function uses secure protocols and verifies data integrity before returning.
func fetchExternalData(taskID string) (string, error) {
	// Define the external API endpoint or data source
	apiEndpoint := fmt.Sprintf("https://trusted-data-provider.com/api/tasks/%s", taskID)

	// Make an HTTP GET request to fetch the external data
	client := &http.Client{
		Timeout: 10 * time.Second, // Set timeout for the request
	}

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add necessary headers for authentication or API access
	req.Header.Set("Authorization", "Bearer YOUR_API_KEY")
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch external data: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Validate data integrity using a checksum
	expectedChecksum := resp.Header.Get("X-Data-Checksum")
	calculatedChecksum := GenerateChecksum(string(body))
	if expectedChecksum != "" && calculatedChecksum != expectedChecksum {
		return "", errors.New("data integrity check failed: checksum mismatch")
	}

	log.Printf("Successfully fetched external data for TaskID: %s", taskID)
	return string(body), nil
}

// processTaskData processes the payload using the fetched external data.
// It applies business logic to transform and validate the results.
func processTaskData(payload string, externalData string) (string, error) {
	// Validate input data
	if payload == "" {
		return "", errors.New("task payload cannot be empty")
	}
	if externalData == "" {
		return "", errors.New("external data cannot be empty")
	}

	// Parse the payload and external data into structured formats (e.g., JSON)
	var payloadData map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &payloadData); err != nil {
		return "", fmt.Errorf("failed to parse task payload: %w", err)
	}

	var externalDataMap map[string]interface{}
	if err := json.Unmarshal([]byte(externalData), &externalDataMap); err != nil {
		return "", fmt.Errorf("failed to parse external data: %w", err)
	}

	// Merge payload and external data based on business rules
	processedResult := make(map[string]interface{})
	for key, value := range payloadData {
		processedResult[key] = value
	}
	for key, value := range externalDataMap {
		if _, exists := processedResult[key]; !exists {
			processedResult[key] = value
		}
	}

	// Validate the merged result
	if len(processedResult) == 0 {
		return "", errors.New("processed result is empty after merging data")
	}

	// Convert the processed result to JSON for storage or further processing
	resultJSON, err := json.Marshal(processedResult)
	if err != nil {
		return "", fmt.Errorf("failed to serialize processed result: %w", err)
	}

	log.Printf("Task data processed successfully. Result: %s", string(resultJSON))
	return string(resultJSON), nil
}

// integrateWithEncryption encrypts the task payload using AES encryption and Base64 encodes the result.
func (aen *AENManager) integrateWithEncryption(task *AENTask) (string, error) {
	// Ensure the encryption service and key are properly initialized
	if aen.EncryptionService == nil {
		return "", errors.New("encryption service is not initialized")
	}

	// Define the AES encryption key (must be 32 bytes for AES-256)
	const aesKey = "your-32-byte-key-for-aes-encryption"
	if len(aesKey) != 32 {
		return "", errors.New("invalid encryption key length: must be 32 bytes for AES-256")
	}

	// Encrypt the task payload
	encryptedPayload, err := aen.EncryptionService.EncryptData("AES", []byte(task.Payload), []byte(aesKey))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt payload: %w", err)
	}

	// Encode the encrypted payload to Base64 for storage or transfer
	encryptedPayloadStr := base64.StdEncoding.EncodeToString(encryptedPayload)

	log.Printf("Payload successfully encrypted for task %s", task.TaskID)
	return encryptedPayloadStr, nil
}



// validateTaskPayload ensures the task payload adheres to expected format and rules.
func validateTaskPayload(taskPayload string) error {
	if len(taskPayload) == 0 {
		return errors.New("task payload cannot be empty")
	}

	// Example validation: Ensure the payload contains only alphanumeric characters
	if match, _ := regexp.MatchString(`^[a-zA-Z0-9\-_.]+$`, taskPayload); !match {
		return errors.New("task payload contains invalid characters")
	}

	log.Printf("Task payload validation successful: %s", taskPayload)
	return nil
}

// validateExternalData ensures external data is valid and consistent with expectations.
func validateExternalData(externalData string) error {
	if len(externalData) == 0 {
		return errors.New("external data cannot be empty")
	}

	// Example validation: Ensure external data contains only alphanumeric characters
	if match, _ := regexp.MatchString(`^[a-zA-Z0-9\-_.]+$`, externalData); !match {
		return errors.New("external data contains invalid characters")
	}

	log.Printf("External data validation successful: %s", externalData)
	return nil
}

// applyBusinessRules applies logic to combine and transform task payload and external data.
func applyBusinessRules(taskPayload string, externalData string) (string, error) {
	// Business Rule: Concatenate payload and external data with a separator
	if len(taskPayload) < 5 {
		return "", errors.New("task payload is too short to process")
	}

	processedData := fmt.Sprintf("%s|%s", taskPayload, externalData)

	// Additional transformations or computations (as needed)
	// Example: Hash the combined data for integrity
	hashedData := GenerateChecksum(processedData)
	log.Printf("Business rules applied. Result: %s (Hash: %s)", processedData, hashedData)

	return processedData, nil
}

// GenerateChecksum computes a SHA-256 checksum of the input string.
func GenerateChecksum(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}
