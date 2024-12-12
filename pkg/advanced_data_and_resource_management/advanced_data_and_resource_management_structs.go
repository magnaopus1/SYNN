package advanced_data_and_resource_management

import (
	"sync"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// AENTask represents an autonomous off-chain task
type AENTask struct {
	TaskID           string    // Unique identifier for the task
	TaskType         string    // Type of the off-chain task (e.g., data processing, AI computation)
	Payload          string    // Data to be processed by the task
	TaskStatus       string    // Status of the task (e.g., "Pending", "Completed", "Failed")
	ExecutorNode     string    // Node executing the task
	CreationTime     time.Time // Task creation timestamp
	CompletionTime   time.Time // Task completion timestamp
	EncryptedPayload string    // Encrypted payload for secure transmission
	CreatorNode    	 string
}

// AENManager handles the execution and management of tasks across nodes in the network
type AENManager struct {
	Tasks             map[string]*AENTask   // Active tasks in the system
	TaskQueue         []string              // Queue of task IDs for processing
	Ledger            *ledger.Ledger        // Ledger instance for task logging
	EncryptionService *common.Encryption // Encryption service for secure task data
	mu                sync.Mutex            // Mutex for concurrent access to task data
}

// AMLARollup represents a rollup of multi-layer data for aggregation
type AMLARollup struct {
	RollupID         string    // Unique ID for the rollup
	Layer            int       // The specific layer being aggregated
	DataPayload      string    // Data being rolled up from the layer
	AggregatorNode   string    // Node responsible for aggregating the rollup
	AggregationTime  time.Time // Time of aggregation
	EncryptedPayload string    // Encrypted rollup payload for security
}

// AMLAData represents a complete multi-layer aggregation from different layers
type AMLAData struct {
	DataID          string         // Unique ID for the aggregation data
	Rollups         []*AMLARollup  // Collection of rollup data from various layers
	ValidatorNode   string         // Node validating the aggregation
	ValidationTime  time.Time      // Time of validation
	EncryptedDataID string         // Encrypted ID for validation
	Status          string         // Aggregation status: "Pending", "Validated", "Failed"
}

// AMLAManager manages rollup aggregation and validation for multi-layer data
type AMLAManager struct {
	Aggregations     map[string]*AMLAData    // Stores aggregations for validation
	PendingRollups   []*AMLARollup           // Stores incoming rollups waiting for aggregation
	CompletedData    []*AMLAData             // Stores validated aggregations
	EncryptionService *common.Encryption  // Encryption service
	Ledger           *ledger.Ledger          // Ledger instance for logging aggregation actions
	mu               sync.Mutex              // Mutex for safe concurrent access
}

// ValidationTask represents a validation task for a block or sub-block
type ValidationTask struct {
	TaskID           string    // Unique ID of the validation task
	SubBlockID       string    // Sub-block ID being validated
	BlockID          string    // Block ID if task involves a full block
	AssignedValidator string   // Validator assigned to this task
	Status           string    // Status of the task ("Pending", "In Progress", "Completed", "Failed")
	ValidationTime   time.Time // Time of task validation
	EncryptedData    string    // Encrypted validation data for security
	Result            string    // Add this field
}

// AVANManager manages the dynamic allocation of validators for validation tasks
type AVANManager struct {
	Validators       map[string]*Validator    // Map of all validators in the network
	PendingTasks     []*ValidationTask        // Queue of pending validation tasks
	CompletedTasks   []*ValidationTask        // Log of completed validation tasks
	EncryptionService *common.Encryption   // Encryption service for security
	Ledger           *ledger.Ledger           // Ledger instance for logging validations
	mu               sync.Mutex               // Mutex for concurrent task management
}


// OrchestratedData represents the data for a given orchestration task.
type OrchestratedData struct {
    DataID            string
    AppID             string
    DataPayload       string
    OrchestrationTime time.Time
    Status            string
    HandlerNode       string
    EncryptedPayload  string
    Result            string // Add this field to store the result of orchestration
}


// DDOLManager manages decentralized data orchestration for off-chain dApp data
type DDOLManager struct {
	Orchestrations     map[string]*OrchestratedData  // Active orchestrations in progress
	CompletedData      []*OrchestratedData           // Log of completed orchestrations
	PendingQueue       []string                     // Queue of pending orchestration requests
	Ledger             *ledger.Ledger               // Ledger instance for recording orchestration logs
	EncryptionService  *common.Encryption        // Encryption service for secure data handling
	mu                 sync.Mutex                   // Mutex for concurrent orchestration operations
}

// Validator represents a validator node in the network
type Validator struct {
	ValidatorID     string // Unique ID of the validator
	NodeID          string // Node ID associated with the validator
	Allocated       bool   // Whether the validator is currently allocated
	EncryptedNodeID []byte // Encrypted version of the NodeID for security
	AllocationTime time.Time
}
