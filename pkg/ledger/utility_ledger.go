package ledger

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "time"
)



// RecordTaskCreation records the creation of a new task in the ledger.
func (l *UtilityLedger) RecordTaskCreation(creator string, details map[string]interface{}) (string, error) {
    l.Lock()
    defer l.Unlock()

    // Initialize Tasks map if nil
    if l.Tasks == nil {
        l.Tasks = make(map[string]TaskRecord)
    }

    // Initialize TaskHistory map if nil
    if l.UtilityLedgerState.TaskHistory == nil {
        l.UtilityLedgerState.TaskHistory = make(map[string]TaskHistory)
    }

    // Generate a unique task ID
    id := l.generateTaskID(creator, details)
    task := TaskRecord{
        ID:          id,
        Creator:     creator,
        Status:      "created",
        CreatedAt:   time.Now(),
        TaskDetails: details,
    }

    // Store the task in the Tasks map
    l.Tasks[id] = task

    // Update the task history for the newly created task
    l.UtilityLedgerState.TaskHistory[id] = TaskHistory{
        TaskID: id,
        Updates: []TaskEvent{
            {
                Timestamp: time.Now(),
                Status:    "created",
                Comment:   "Task created successfully.",
            },
        },
    }

    return id, nil
}



// RecordTaskAllocation assigns a task to a specific entity in the ledger.
func (l *UtilityLedger) RecordTaskAllocation(taskID, assignee string) error {
    l.Lock()
    defer l.Unlock()

    task, exists := l.Tasks[taskID]
    if !exists {
        return errors.New("task not found")
    }

    task.Assignee = assignee
    task.Status = "allocated"
    l.Tasks[taskID] = task

    return nil
}


// RecordTaskCompletion records the completion of a task in the ledger.
func (l *UtilityLedger) RecordTaskCompletion(taskID string, subBlockID string, blockID string) error {
    l.Lock()
    defer l.Unlock()

    task, exists := l.Tasks[taskID]
    if !exists {
        return errors.New("task not found")
    }

    now := time.Now()
    task.CompletedAt = &now
    task.Status = "completed"
    task.SubBlockID = subBlockID
    task.BlockID = blockID
    l.Tasks[taskID] = task

    return nil
}


// RecordOrchestrationRequest logs a new orchestration request into the ledger.
func (l *UtilityLedger) RecordOrchestrationRequest(requester string, resources map[string]interface{}) (string, error) {
	l.Lock()
	defer l.Unlock()

	// Initialize OrchestrationRequests if nil
	if l.OrchestrationRequests == nil {
		l.OrchestrationRequests = make(map[string]OrchestrationRequest)
	}

	// Initialize OrchestrationHistory if nil
	if l.UtilityLedgerState.OrchestrationHistory == nil {
		l.UtilityLedgerState.OrchestrationHistory = make(map[string]OrchestrationRecord)
	}

	// Generate a unique orchestration request ID
	id := l.generateOrchestrationRequestID(requester, resources)
	request := OrchestrationRequest{
		ID:        id,
		Requester: requester,
		Resources: resources,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	// Store the orchestration request in OrchestrationRequests
	l.OrchestrationRequests[id] = request

	// Log the request in OrchestrationHistory
	l.UtilityLedgerState.OrchestrationHistory[id] = OrchestrationRecord{
		RecordID:  id,
		Action:    "create",
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Orchestration request created by %s", requester),
	}

	return id, nil
}


// RecordTask handles the general recording of a task within the ledger, ensuring task integrity.
func (l *UtilityLedger) RecordTask(task TaskRecord) error {
    l.Lock()
    defer l.Unlock()

    // Initialize Tasks map if nil
    if l.Tasks == nil {
        l.Tasks = make(map[string]TaskRecord)
    }

    // Initialize TaskHistory map if nil
    if l.UtilityLedgerState.TaskHistory == nil {
        l.UtilityLedgerState.TaskHistory = make(map[string]TaskHistory)
    }

    // Check if a task with the same ID already exists
    if _, exists := l.Tasks[task.ID]; exists {
        return errors.New("task with this ID already exists")
    }

    // Store the task in the Tasks map
    l.Tasks[task.ID] = task

    // Add a new entry to TaskHistory
    l.UtilityLedgerState.TaskHistory[task.ID] = TaskHistory{
        TaskID: task.ID,
        Updates: []TaskEvent{
            {
                Timestamp: time.Now(),
                Status:    "created",
                Comment:   "Task recorded successfully.",
            },
        },
    }

    return nil
}



// generateTaskID generates a unique ID for a task based on the creator and task details.
func (l *UtilityLedger) generateTaskID(creator string, details map[string]interface{}) string {
    input := fmt.Sprintf("%s-%v-%d", creator, details, time.Now().UnixNano())
    hash := sha256.Sum256([]byte(input))
    return hex.EncodeToString(hash[:])
}

// generateOrchestrationRequestID generates a unique ID for an orchestration request.
func (l *UtilityLedger) generateOrchestrationRequestID(requester string, resources map[string]interface{}) string {
    input := fmt.Sprintf("%s-%v-%d", requester, resources, time.Now().UnixNano())
    hash := sha256.Sum256([]byte(input))
    return hex.EncodeToString(hash[:])
}

// RecordOrchestrationCompletion records the completion of an orchestration request in the ledger.
func (l *UtilityLedger) RecordOrchestrationCompletion(dataID string, details map[string]interface{}) error {
    l.Lock()
    defer l.Unlock()

    if l.OrchestrationRecords == nil {
        l.OrchestrationRecords = make(map[string]map[string]interface{})
    }

    // Ensure the record exists before marking it as completed
    record, exists := l.OrchestrationRecords[dataID]
    if !exists {
        return fmt.Errorf("orchestration request with dataID %s not found", dataID)
    }

    // Mark the orchestration as completed by adding the completion details
    for key, value := range details {
        record[key] = value
    }
    record["Status"] = "Completed"

    // Update the ledger record
    l.OrchestrationRecords[dataID] = record
    return nil
}




// RecordTaskFailure records the failure of a task in the ledger.
func (l *UtilityLedger) RecordTaskFailure(taskID string, executorNode string, errorMessage string) error {
    l.Lock()
    defer l.Unlock()

    if l.Tasks == nil {
        l.Tasks = make(map[string]TaskRecord)
    }

    // Check if task already exists in the ledger, if not create a new one
    taskRecord, exists := l.Tasks[taskID]
    if !exists {
        taskRecord = TaskRecord{
            ID:           taskID,
            ExecutorNode: executorNode,
            Status:       "Failed",
        }
    }

    // Update task failure details
    taskRecord.ErrorMessage = errorMessage
    taskRecord.Status = "Failed"

    // Update ledger
    l.Tasks[taskID] = taskRecord

    // Simulate logging to persistent storage (e.g., database, file, etc.)
    fmt.Printf("Task %s failed by %s with error: %s\n", taskID, executorNode, errorMessage)

    return nil
}

