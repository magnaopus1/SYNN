package advanced_data_and_resource_management

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewAVANManager initializes the Autonomous Validator Allocation Network manager.
func NewAVANManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *AVANManager {
	return &AVANManager{
		Validators:        make(map[string]*Validator),
		PendingTasks:      []*ValidationTask{},
		CompletedTasks:    []*ValidationTask{},
		EncryptionService: encryptionService,
		Ledger:            ledgerInstance,
	}
}

// RegisterValidator adds a new validator node to the network.
func (avan *AVANManager) RegisterValidator(nodeID string) (*Validator, error) {
	avan.mu.Lock()
	defer avan.mu.Unlock()

	if len(nodeID) == 0 {
		return nil, errors.New("nodeID cannot be empty")
	}

	// Generate a unique validator ID
	validatorID := generateUniqueID()

	// Encrypt the node ID using AES and the EncryptionKey
	encryptedNodeID, err := avan.EncryptionService.EncryptData("AES", []byte(nodeID), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt node ID: %w", err)
	}

	// Encode encryptedNodeID as a base64 string
	encodedNodeID := base64.StdEncoding.EncodeToString(encryptedNodeID)

	// Adjust the Validator struct initialization
	validator := &Validator{
		ValidatorID:     validatorID,
		NodeID:          nodeID,
		Allocated:       false,
		EncryptedNodeID: []byte(encodedNodeID), // Ensure the field type matches []byte
	}

	// Add the validator to the network's list of validators
	avan.Validators[validatorID] = validator

	// Log the validator registration in the ledger
	ledgerMsg, err := avan.Ledger.BlockchainConsensusCoinLedger.RecordValidatorRegistration(validatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to log validator registration in ledger: %w", err)
	}
	log.Printf("Validator %s registered successfully. Ledger Message: %s", validatorID, ledgerMsg)

	return validator, nil
}


// AllocateValidator dynamically allocates a free validator to a pending validation task.
func (avan *AVANManager) AllocateValidator(taskID, subBlockID, blockID string) (*ValidationTask, error) {
	avan.mu.Lock()
	defer avan.mu.Unlock()

	if len(taskID) == 0 || len(subBlockID) == 0 || len(blockID) == 0 {
		return nil, errors.New("taskID, subBlockID, and blockID must not be empty")
	}

	// Find an available validator
	var availableValidator *Validator
	for _, validator := range avan.Validators {
		if !validator.Allocated {
			availableValidator = validator
			break
		}
	}

	if availableValidator == nil {
		return nil, errors.New("no available validators for allocation")
	}

	// Mark the validator as allocated
	availableValidator.Allocated = true
	availableValidator.AllocationTime = time.Now()

	// Create a new validation task
	validationTask := &ValidationTask{
		TaskID:            taskID,
		SubBlockID:        subBlockID,
		BlockID:           blockID,
		AssignedValidator: availableValidator.ValidatorID,
		Status:            "Pending",
	}

	// Add the task to the pending queue
	avan.PendingTasks = append(avan.PendingTasks, validationTask)

	// Log the task allocation in the ledger
	err := avan.Ledger.UtilityLedger.RecordTaskAllocation(taskID, availableValidator.ValidatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to log task allocation in ledger: %w", err)
	}
	log.Printf("Task %s allocated to validator %s (sub-block: %s, block: %s)", taskID, availableValidator.ValidatorID, subBlockID, blockID)

	return validationTask, nil
}

// CompleteTask marks a validation task as completed and logs it in the ledger.
func (avan *AVANManager) CompleteTask(taskID string, result string) error {
	avan.mu.Lock()
	defer avan.mu.Unlock()

	if len(taskID) == 0 || len(result) == 0 {
		return errors.New("taskID and result cannot be empty")
	}

	// Find the task in the pending list
	var task *ValidationTask
	taskIndex := -1
	for i, t := range avan.PendingTasks {
		if t.TaskID == taskID {
			task = t
			taskIndex = i
			break
		}
	}
	if task == nil {
		return fmt.Errorf("task %s not found in pending tasks", taskID)
	}

	// Update task fields (ensure the struct has these fields or adjust accordingly)
	task.Status = "Completed"
	task.Result = result       // Ensure `Result` is a valid field in ValidationTask
	task.ValidationTime = time.Now() // Ensure `CompletedAt` is a valid field in ValidationTask

	// Release the validator assigned to the task
	validator, exists := avan.Validators[task.AssignedValidator]
	if !exists {
		return fmt.Errorf("validator %s not found for task %s", task.AssignedValidator, taskID)
	}
	validator.Allocated = false
	validator.AllocationTime = time.Time{} // Reset allocation time

	// Move the task to the completed list and remove from pending
	avan.CompletedTasks = append(avan.CompletedTasks, task)
	avan.PendingTasks = append(avan.PendingTasks[:taskIndex], avan.PendingTasks[taskIndex+1:]...)

	// Log the task completion in the ledger
	err := avan.Ledger.UtilityLedger.RecordTaskCompletion(taskID, task.AssignedValidator, result) // Adjust to include all required arguments
	if err != nil {
		return fmt.Errorf("failed to log task completion in ledger: %w", err)
	}
	log.Printf("Task %s completed successfully with result: %s", taskID, result)

	return nil
}


// CompleteValidationTask marks a task as completed and releases the validator.
func (avan *AVANManager) CompleteValidationTask(taskID string) error {
	avan.mu.Lock()
	defer avan.mu.Unlock()

	// Find the validation task
	var task *ValidationTask
	var taskIndex int
	found := false
	for i, t := range avan.PendingTasks {
		if t.TaskID == taskID {
			task = t
			taskIndex = i
			found = true
			break
		}
	}

	if !found || task == nil {
		return fmt.Errorf("validation task %s not found", taskID)
	}

	// Mark task as completed
	task.Status = "Completed"
	task.ValidationTime = time.Now()

	// Encrypt the validation data for security
	plainText := []byte(fmt.Sprintf("%s:%s", task.SubBlockID, task.BlockID))
	encryptedData, err := avan.EncryptionService.EncryptData("AES", plainText, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt validation data for task %s: %w", taskID, err)
	}
	task.EncryptedData = base64.StdEncoding.EncodeToString(encryptedData)

	// Release the validator
	validator, ok := avan.Validators[task.AssignedValidator]
	if !ok {
		return fmt.Errorf("validator %s not found for task %s", task.AssignedValidator, taskID)
	}
	validator.Allocated = false
	validator.AllocationTime = time.Time{} // Reset allocation time

	// Remove task from pending queue and add to completed tasks
	avan.PendingTasks = append(avan.PendingTasks[:taskIndex], avan.PendingTasks[taskIndex+1:]...)
	avan.CompletedTasks = append(avan.CompletedTasks, task)

	// Log task completion in the ledger
	validationTimeStr := task.ValidationTime.Format(time.RFC3339)
	err = avan.Ledger.UtilityLedger.RecordTaskCompletion(taskID, validator.ValidatorID, validationTimeStr)
	if err != nil {
		return fmt.Errorf("failed to log task completion in ledger for task %s: %w", taskID, err)
	}

	log.Printf("Task %s completed by validator %s at %s", taskID, validator.ValidatorID, validationTimeStr)
	return nil
}

// MonitorTaskStatus checks the status of a specific validation task.
func (avan *AVANManager) MonitorTaskStatus(taskID string) (string, error) {
	avan.mu.Lock()
	defer avan.mu.Unlock()

	// Search for the task in pending tasks
	for _, task := range avan.PendingTasks {
		if task.TaskID == taskID {
			log.Printf("Task %s status: %s (Pending)", taskID, task.Status)
			return task.Status, nil
		}
	}

	// Search for the task in completed tasks
	for _, task := range avan.CompletedTasks {
		if task.TaskID == taskID {
			log.Printf("Task %s status: %s (Completed)", taskID, task.Status)
			return task.Status, nil
		}
	}

	return "", fmt.Errorf("task %s not found", taskID)
}

