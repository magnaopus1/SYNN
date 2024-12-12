package ledger

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/json"
	"math/big"
	"net"
	"sync"
	"time"
)

// ************** Account And Balance Structs **************

// Supplementary Structures
type BalanceSnapshot struct {
	AccountID string
	Balance   float64
	Timestamp time.Time
}

// BalanceTransfer represents a single balance transfer between two accounts.
type BalanceTransfer struct {
	FromID string  // Source account ID
	ToID   string  // Destination account ID
	Amount float64 // Amount to be transferred
}

// Mnemonic represents a cryptographic mnemonic phrase used for generating keys or wallets.
type Mnemonic struct {
	Phrase         string            // The mnemonic phrase as a single string
	WordList       []string          // The list of words making up the mnemonic
	Language       string            // The language of the mnemonic (e.g., English, Japanese)
	CreatedAt      time.Time         // The timestamp when the mnemonic was generated
	Entropy        []byte            // The entropy used to generate the mnemonic
	Passphrase     string            // An optional passphrase for additional security
	Checksum       string            // The checksum for verifying the integrity of the mnemonic
	Purpose        string            // Purpose of the mnemonic (e.g., Wallet, Authentication)
	Version        string            // Version of the mnemonic standard (e.g., BIP-39)
	AssociatedData map[string]string // Additional metadata associated with the mnemonic
}

// TrustAccount represents a trust account within the ledger.
type TrustAccount struct {
	ID      string
	Balance *big.Float // Use *big.Float for high precision balance
}

// Account represents a user's balance and other account-related information.
type Account struct {
	Balance               float64
	HeldBalance           float64 // Add this field to hold a locked amount
	CreatedAt             time.Time
	Address               string
	Nonce                 uint64
	Verified              bool
	EncryptedKey          string
	ConnectionEvents      []ConnectionEvent
	ContractExecutionLogs []ContractExecutionLog
	MintRecords           []MintRecord
	BurnRecords           []BurnRecord
	CurrencyExchanges     []CurrencyExchange
	CustomName            string
	PrivateKey            string
	PublicKey             string
	Stake                 float64
	IsAdmin               bool
	Permissions           []string
	LockedBalances        []BalanceLock   // Holds locked balance entries
	ReservedBalance       float64         // Holds reserved balance
	FreezeUntil           time.Time       // Indicates until when the account balance is frozen
	IsSuspicious          bool            // Flag for suspicious accounts
	BalanceStatus         BalanceStatus   // Holds the balance status of the account
	IsFrozen              bool            // Indicates if the account is frozen
	EncryptedBalance      string          // Encrypted representation of the balance (if applicable)
	ExternalAccountID     string          // Link to an external account
	Authorizations        []Authorization // List of authorizations for this account
	VerificationHold      float64         // Holds funds for verification purposes
	LockedBalance         float64         // Balance that is locked
	Approvals             []Approval      // List of approved transactions
	LastTransactionID     string          // ID of the last transaction affecting this account
	TotalDeposited        float64         // Total deposited into the account
	TotalWithdrawn        float64         // Total withdrawn from the account
	LastUpdated           time.Time       // Last time the account was updated
	RequiresReview        bool            // Indicates if the account is flagged for review
	Allocations           []Allocation    // List of allocations for specific purposes
}

type Pool struct {
	ID      string
	Balance *big.Float // Updated to big.Float for precision
}

// BalanceLock represents a lock on a specified amount within an account until a specific time.
type BalanceLock struct {
	ID        string    // Unique identifier for the lock
	AccountID string    // The account to which this lock applies
	Amount    float64   // Amount locked
	UnlockAt  time.Time // Time when the lock can be released
}

// BalanceStatus represents the current status of an account's balance.
type BalanceStatus struct {
	AccountID    string    // Unique identifier for the account
	IsActive     bool      // Indicates if the balance is active
	IsFrozen     bool      // Indicates if the balance is frozen
	IsOnHold     bool      // Indicates if the balance is on hold
	FreezeReason string    // Reason for balance freeze, if applicable
	HoldReason   string    // Reason for balance hold, if applicable
	UpdatedAt    time.Time // Last time the status was updated
}

// Authorization represents permissions or access rights assigned to an account.
type Authorization struct {
	ID          string   // Unique ID for the authorization
	Description string   // Description of the authorization
	Permissions []string // List of permissions granted
}

// BalanceInfo represents detailed information about an account’s balance.
type BalanceInfo struct {
	AccountID         string        // Unique identifier for the account
	AvailableBalance  float64       // Available balance for transactions
	LockedBalance     float64       // Balance that is locked and cannot be used
	HeldBalance       float64       // Balance that is held (e.g., for pending transactions)
	LastTransactionID string        // ID of the last transaction affecting the balance
	TotalDeposited    float64       // Total deposited into the account
	TotalWithdrawn    float64       // Total withdrawn from the account
	Status            BalanceStatus // Current status of the balance
	LastUpdated       time.Time     // Last time the balance was updated
}

// Approval represents an approval granted for specific transactions or operations.
type Approval struct {
	ID            string    // Unique identifier for the approval
	AccountID     string    // ID of the account granting the approval
	ApproverID    string    // ID of the approving authority or user
	TransactionID string    // ID of the related transaction (if applicable)
	ApprovedAt    time.Time // Timestamp of approval
	ExpiresAt     time.Time // Expiration timestamp for time-limited approvals
	Status        string    // Status of the approval (e.g., "Pending", "Approved", "Rejected")
	Remarks       string    // Optional remarks or comments related to the approval
}

// Allocation represents an allocation of resources or funds.
type Allocation struct {
	ID           string    // Unique identifier for the allocation
	AccountID    string    // ID of the account to which the allocation is associated
	ResourceType string    // Type of resource being allocated (e.g., "funds", "CPU", "memory")
	AllocatedAt  time.Time // Timestamp when the allocation was made
	ExpiresAt    time.Time // Expiration time for the allocation (if applicable)
	Amount       float64   // Amount of resource allocated
	Status       string    // Current status of the allocation (e.g., "active", "expired", "pending")
	Remarks      string    // Optional remarks or comments related to the allocation
	ApprovalID   string    // ID of any approval associated with the allocation
	Allocated    bool      // Indicates if the allocation is active (not yet redeemed)

}

// ************** A.I & Machine Learning Structs **************

// InferenceRecord holds inference details.
type InferenceRecord struct {
	ModelID   string
	Timestamp time.Time
	Result    string
	NodeID    string
	Processed bool
}

// AccessList defines the list of entities or users with access to a model.
type AccessList struct {
	ModelID      string    // Unique identifier for the model
	AllowedUsers []string  // List of user IDs with access
	AllowedRoles []string  // List of roles with access
	Expiration   time.Time // Expiration of access, if any
	AccessType   string    // Type of access (e.g., Read, Write, Execute)
}

// PermissionRecord tracks permissions granted for a model.
type PermissionRecord struct {
	UserID      string    // Unique identifier of the user
	ModelID     string    // Unique identifier of the model
	Permissions []string  // List of permissions (e.g., "read", "write", "train")
	GrantedAt   time.Time // Timestamp when permissions were granted
	ExpiresAt   time.Time // Timestamp when permissions expire, if any
	GrantedBy   string    // Identifier of the entity that granted the permissions
}

// ContainerInfo represents information about a containerized AI/ML model or process.
type ContainerInfo struct {
	ContainerID   string    // Unique identifier for the container
	ModelID       string    // Unique identifier for the model
	Status        string    // Current status of the container (e.g., Running, Stopped)
	CreatedAt     time.Time // Timestamp when the container was created
	Resources     Resources // Resource usage by the container
	NodeID        string    // Node where the container is hosted
	ContainerLogs []string  // Logs generated by the container
}

// Resources represents the resources allocated to a container.
type Resources struct {
	CPU       float64 // CPU usage in cores
	Memory    float64 // Memory usage in GB
	DiskSpace float64 // Disk space usage in GB
	Network   float64 // Network usage in Mbps
}

// ServiceInfo represents a service associated with an AI/ML model.
type ServiceInfo struct {
	ServiceID        string            // Unique identifier for the service
	Name             string            // Name of the service
	Description      string            // Description of the service
	Status           string            // Current status of the service (e.g., Active, Deprecated)
	AssociatedModels []string          // List of model IDs associated with the service
	Metadata         map[string]string // Additional metadata for the service
	CreatedAt        time.Time         // Timestamp when the service was created
	UpdatedAt        time.Time         // Timestamp when the service was last updated
}

// RecommendationRecord stores recommendations generated by an AI/ML model.
type RecommendationRecord struct {
	RecommendationID string            // Unique identifier for the recommendation
	ModelID          string            // Model that generated the recommendation
	GeneratedAt      time.Time         // Timestamp when the recommendation was generated
	TargetEntity     string            // Entity for which the recommendation is intended
	Details          map[string]string // Details of the recommendation
	Status           string            // Current status of the recommendation (e.g., Pending, Accepted, Rejected)
}

// AnalysisRecord for AI analysis sessions.
type AnalysisRecord struct {
	AnalysisID string
	ModelID    string
	NodeID     string
	StartTime  time.Time
	StopTime   time.Time
	Status     string
}

// RecommendationCriteria represents the structure for criteria used to generate recommendations.
type RecommendationCriteria struct {
	UserID             string            `json:"user_id"`
	PreferenceScore    float64           `json:"preference_score"`
	MaxRecommendations int               `json:"max_recommendations"`
	Preference         string            `json:"preference"`          // User preference, e.g., "popular", "recent"
	Threshold          float64           `json:"threshold"`           // Minimum score threshold for recommendations
	MaxSuggestions     int               `json:"max_suggestions"`     // Maximum number of suggestions
	Tags               []string          `json:"tags"`                // Tags to filter recommendations, e.g., ["electronics", "books"]
	ExcludeList        []string          `json:"exclude_list"`        // Items to exclude from recommendations
	UserHistoryWeight  float64           `json:"user_history_weight"` // Weight given to user’s past interactions
	ContextualFactors  map[string]string `json:"contextual_factors"`  // Additional factors, e.g., {"time_of_day": "morning"}
}

// Recommendation struct for AI recommendations.
type Recommendation struct {
	ModelID   string
	Timestamp time.Time
	Content   string
	Updated   bool
}

// Container represents a deployed container.
type Container struct {
	ContainerID string
	ModelID     string
	Status      string
	NodeID      string
}

type AiService struct {
	ServiceID   string
	Status      string
	Metrics     ServiceMetrics
	LastUpdated time.Time
	Balance     float64 // New field to store the service balance
}

// DataProcessingLog logs data processing activities.
type DataProcessingLog struct {
	ProcessID string
	ModelID   string
	NodeID    string
	StartTime time.Time
	EndTime   time.Time
	Status    string
}

// TrafficRecord logs load balancing and traffic operations.
type TrafficRecord struct {
	ModelID   string
	NodeID    string
	Action    string
	Timestamp time.Time
}

// ServiceMetrics stores metadata and encrypted data for an AI service.
type ServiceMetrics struct {
	ServiceName   string    // Name of the service
	Owner         string    // Owner of the service
	EncryptedData []byte    // Encrypted service data
	CreatedAt     time.Time // Creation timestamp
	Status        string    // Current status ("Active", "Inactive", etc.)
}

// ModelRestriction represents model restrictions.
type ModelRestriction struct {
	ModelID    string
	Restricted bool
	Reason     string
	Timestamp  time.Time
}

// ModelPermissions represents permissions for a model.
type ModelPermissions struct {
	ModelID      string
	AllowedUsers []string
}

// AccessToken represents access token for model operations.
type AccessToken struct {
	TokenID    string
	ModelID    string
	GrantedTo  string
	Expiry     time.Time
	ermissions string
}

// EncryptionLog struct for storing encryption details
type EncryptionLog struct {
	TransactionID string
	EncryptedData string
	Timestamp     time.Time
}

// DecryptionLog struct for storing decryption details
type DecryptionLog struct {
	TransactionID string
	DecryptedData string
	Timestamp     time.Time
}

// ModelCheckpoint represents a model checkpoint.
type ModelCheckpoint struct {
	ModelID   string
	Version   int
	CreatedAt time.Time
	DataHash  string
}

// AccessLog represents a log entry for model access.
type ModelAccessLog struct {
	ModelID   string
	UserID    string
	Action    string
	Timestamp time.Time
}

// UsageStatistics represents model usage statistics.
type UsageStatistics struct {
	ModelID       string
	UsageCount    int
	LastUsedAt    time.Time
	UsageDuration time.Duration
}

// PerformanceMetrics represents model performance metrics.
type ModelPerformanceMetrics struct {
	ModelID     string
	Accuracy    float64
	Loss        float64
	LastUpdated time.Time
}

// ComplianceCheck records model compliance checks.
type ComplianceCheck struct {
	ModelID   string
	Status    string
	Timestamp time.Time
	Details   string
}

// SecurityAudit represents a security audit entry.
type SecurityAudit struct {
	ModelID   string
	Passed    bool
	Timestamp time.Time
	Findings  string
}

// ResourceAllocation represents allocated resources.
type AiModelResourceAllocation struct {
	ResourceID  string
	ModelID     string
	Amount      float64
	AllocatedAt time.Time
}

// StorageAllocation represents allocated storage.
type AiModelStorageAllocation struct {
	StorageID   string
	ModelID     string
	SizeMB      int
	AllocatedAt time.Time
}

// CacheData represents cached data.
type CacheData struct {
	ModelID   string
	DataID    string
	Data      string
	CreatedAt time.Time
}

// DeploymentCheck records deployment checks.
type DeploymentCheck struct {
	DeploymentID string
	Status       string
	Timestamp    time.Time
	Details      string
}

// TrainingStatus represents the training status of a model.
type TrainingStatus struct {
	ModelID     string
	Status      string
	LastUpdated time.Time
}

// RunStatus represents the running status of a model.
type RunStatus struct {
	ModelID     string
	Status      string
	LastChecked time.Time
}

type Model struct {
	ModelID             string             // Unique identifier for the model
	ModelName           string             // Name of the model
	Version             string             // Version of the model
	LastUpdated         time.Time          // Timestamp of the last update
	IsDeployed          bool               // Deployment status
	Capabilities        []string           // Capabilities like ["ImageClassification", "ObjectDetection"]
	ProcessingModes     []string           // Modes such as ["HighAccuracy", "FastProcessing"]
	Status              string             // Current model status (e.g., "idle", "running", "error")
	ModelType           string             // Type of model (e.g., "Deep Learning", "Machine Learning")
	Classification      string             // Classification category of the model (e.g., "Convolutional Neural Network")
	Purpose             string             // Model's intended purpose (e.g., "Image Recognition")
	Description         string             // Brief description of the model's functionality
	DocumentationURL    string             // URL to detailed model documentation
	Permissions         map[string]string  // Permissions, e.g., {"read": "public", "write": "restricted"}
	TrainingLevel       string             // Level of training applied (e.g., "Basic", "Advanced", "Specialized")
	LedgerInstance      *Ledger            // Instance of the ledger to record transactions
	ResourceUsage       map[string]float64 // Resource usage stats, e.g., {"CPU": 75.0, "Memory": 60.5}
	PerformanceMetrics  map[string]float64 // Performance metrics, e.g., {"Latency": 1.2, "Throughput": 200}
	CurrentScale        int                // Current scale level
	MaxScale            int                // Maximum allowed scale level
	MinScale            int                // Minimum allowed scale level
	RecommendationCache map[string][]byte  // Caches recommendations by transaction ID
	InferenceCount      int                // Track the number of inferences
	AnalysisSessions    int                // Track the number of analysis sessions
	PredictionCount     int                // Track the number of predictions
}

// ActionLog struct stores log details for actions performed by models.
type ActionLog struct {
	TransactionID string
	ModelID       string
	Action        string
	Timestamp     time.Time
}

// ModelActionRecord stores details of a model action (deployment, undeployment, update).
type ModelActionRecord struct {
	TransactionID string
	ModelID       string
	Action        string
	Timestamp     time.Time
	Description   string
}

// ModelIndex struct for AI/ML models.
type ModelIndex struct {
	ModelID             string            // Unique identifier for the model
	ModelName           string            // Name of the model
	NodeLocation        string            // Node or server location where the model resides
	Load                int               // Load level or usage of the model
	Status              string            // Current status of the model in the network
	IPFSLink            string            // IPFS/Swarm link to the model’s data
	Timestamp           time.Time         // Time when the model was added to the index
	TrainingHistory     []string          // History of training sessions
	Permissions         map[string]string // Access permissions, e.g., {"user1": "read"}
	OptimizationHistory []string          // Optimization logs
	VersionHash         string            // Current version hash
	DeployedAt          time.Time         // Deployment timestamp
	OffChainLink        string            // Off-chain storage link for model data (new field added)
	TestHistory         []string          // History of test result IDs
	AccessLogs          []time.Time       // Log of access times for monitoring

}

// ScalingLog struct for recording scaling actions
type ScalingLog struct {
	TransactionID string    // ID of the scaling transaction
	ModelID       string    // ID of the scaled model
	Direction     string    // "up" or "down" indicating the scaling direction
	EncryptedLog  string    // Encrypted log of the scaling action
	Timestamp     time.Time // Timestamp of the scaling action
}

// RecommendationUpdateData represents the new data for updating an existing recommendation.
type RecommendationUpdateData struct {
	UserID        string  `json:"user_id"`
	FeedbackScore float64 `json:"feedback_score"`
	NewCriteria   []byte  `json:"new_criteria"`
	UpdateReason  string  `json:"update_reason"`
}

// ModelExecutionResult represents the structure of the result after model execution.
type ModelExecutionResult struct {
	Status         string    `json:"status"`
	Output         string    `json:"output"`
	Executed       time.Time `json:"executed_at"`
	ProcessingTime string    `json:"processing_time"`
	Confidence     float64   `json:"confidence"` // Confidence level for the output (simulated)
}

// ModelTrainingResult represents the result of a model training session.
type ModelTrainingResult struct {
	Status          string    `json:"status"`
	TrainingID      string    `json:"training_id"`
	UpdatedAt       time.Time `json:"updated_at"`
	EpochsCompleted int       `json:"epochs_completed"`
	Loss            float64   `json:"loss"`
	Accuracy        float64   `json:"accuracy"`
	TrainingTime    string    `json:"training_time"` // Total time taken for training
}

// ************** Advanced Data And Resource Management Structs **************

// TaskRecord represents a task that is recorded in the ledger.
type TaskRecord struct {
	ID           string
	Creator      string
	Assignee     string
	Status       string
	CreatedAt    time.Time
	CompletedAt  *time.Time
	SubBlockID   string
	BlockID      string
	TaskDetails  map[string]interface{}
	ExecutorNode string
	ErrorMessage string
}

// OrchestrationRequest represents a resource orchestration request.
type OrchestrationRequest struct {
	ID        string
	Requester string
	Resources map[string]interface{}
	Status    string
	CreatedAt time.Time
}

// ************** Advanced Security Structs **************

type MitigationMetrics struct {
	IncidentReductionRate       float64 // Percentage of incident reduction
	PerformanceImprovementScore float64 // Score out of 100 for performance improvement
	LastEvaluation              time.Time
}

// ThreatDetectionManager handles threat detection systems.
type ThreatDetectionManager struct {
	Active bool
}

// Channel represents a channel created within the ledger.
type Channel struct {
	ID           string
	Creator      string
	Participants []string
	CreatedAt    time.Time
	ClosedAt     *time.Time
	Status       string
}

type AlertPolicy struct {
	Policy  string
	Version int
	SetAt   time.Time
}

// FirmwareMetadata holds expected firmware integrity data
type FirmwareMetadata struct {
	ExpectedHash string
	Version      string
}

type AnomalyEvent struct {
	Event     string
	Timestamp string // Use a string format for consistency
}

// AggregationValidation represents the validation of an aggregation in the ledger.
type AggregationValidation struct {
	ID         string
	Aggregator string
	Validated  bool
	Timestamp  time.Time
}

// NodeStatus represents the status of a node in the network.
type NodeStatus struct {
	NodeID   string
	IsActive bool
}

type SessionTimeoutLog struct {
	SessionID string
	UserID    string
	TimeoutAt time.Time
	LoggedAt  time.Time
}

// ThreatDetectionStatus represents the status of threat detection for a specific entity.
type ThreatDetectionStatus struct {
	ThreatID        string    // Unique identifier for the detected threat
	DetectedAt      time.Time // Timestamp when the threat was detected
	Severity        string    // Severity level (e.g., Low, Medium, High, Critical)
	Status          string    // Current status (e.g., Detected, Mitigated, Escalated)
	MitigationSteps []string  // List of mitigation actions taken
	DetectionSource string    // Source of detection (e.g., IDS, Antivirus, Manual)
	AssociatedNodes []string  // Nodes affected by the threat
}

// DetectionStatus represents the status of a detection mechanism.
type DetectionStatus struct {
	SystemID          string    // Unique identifier for the system
	LastChecked       time.Time // Last time the detection was checked
	Operational       bool      // Whether the detection system is operational
	DetectedAnomalies int       // Count of detected anomalies
}

// RateLimitingStatus represents the status of rate limiting for a node or API.
type RateLimitingStatus struct {
	NodeID       string    // Unique identifier for the node
	CurrentUsage int       // Current API calls or data usage
	Limit        int       // Maximum allowed usage
	ResetTime    time.Time // Time when the limit will reset
	Status       string    // Status of rate limiting (e.g., Active, Suspended)
}

// EventMonitoringStatus represents the status of monitoring events in the system.
type EventMonitoringStatus struct {
	MonitorID   string    // Unique identifier for the monitor
	Active      bool      // Whether monitoring is active
	LastEvent   string    // Description of the last monitored event
	LastUpdated time.Time // Last time the monitor was updated
	EventCount  int       // Total number of monitored events
}

// ActiveSession represents a currently active session in the system.
type ActiveSession struct {
	SessionID     string        // Unique identifier for the session
	UserID        string        // Identifier for the user associated with the session
	StartTime     time.Time     // Timestamp for when the session started
	LastActivity  time.Time     // Timestamp for the last recorded activity in the session
	IPAddress     string        // IP address from which the session was initiated
	DeviceDetails string        // Information about the device used (e.g., "iOS, Chrome Browser")
	SessionToken  string        // Token associated with the session for validation
	IsActive      bool          // Indicates if the session is currently active
	Permissions   []string      // List of permissions or roles granted during the session
	SessionType   string        // Type of session (e.g., "User", "Admin", "Service")
	Timeout       time.Duration // Session timeout duration
}

// NodeAccessLimitPolicy defines the access limits for a node.
type NodeAccessLimitPolicy struct {
	NodeID          string      // Unique identifier for the node
	MaxConnections  int         // Maximum allowed concurrent connections
	AccessTimes     []time.Time // Allowed access times
	RestrictedIPs   []string    // List of restricted IP addresses
	EnforcementType string      // Enforcement type (e.g., Hard, Soft)
}

// AccessFrequencyPolicy defines the policy for access frequency.
type AccessFrequencyPolicy struct {
	UserID          string        // Unique identifier for the user
	MaxRequests     int           // Maximum allowed requests per period
	Period          time.Duration // Time period for the requests
	CurrentRequests int           // Current count of requests
	Status          string        // Status of the policy (e.g., Active, Suspended)
}

// FirmwareCheckStatus represents the status of firmware checks on a device.
type FirmwareCheckStatus struct {
	DeviceID        string    // Unique identifier for the device
	LastChecked     time.Time // Last time the firmware was checked
	Version         string    // Current firmware version
	Outdated        bool      // Whether the firmware is outdated
	SecurityRisk    string    // Security risk level if outdated
	UpdateAvailable bool      // Whether an update is available
}

// EventLog represents a log of system events.
type EventLog struct {
	EventID       string    // Unique identifier for the event
	Timestamp     time.Time // Timestamp of the event
	EventType     string    // Type of event (e.g., Error, Warning, Info)
	Details       string    // Detailed description of the event
	Source        string    // Source of the event
	AffectedNodes []string  // Nodes affected by the event
}

// HardeningEvent represents an event related to application hardening.
type HardeningEvent struct {
	EventID         string    // Unique identifier for the event
	Timestamp       time.Time // Timestamp of the event
	HardeningType   string    // Type of hardening (e.g., Patch, Configuration)
	Impact          string    // Impact of the hardening
	AssociatedNodes []string  // Nodes affected by the event
}

// AlertResponse represents the response to a triggered alert.
type AlertResponse struct {
	AlertID         string    // Unique identifier for the alert
	RespondedAt     time.Time // Timestamp of the response
	ResponseActions []string  // List of actions taken in response
	Status          string    // Current status of the alert (e.g., Resolved, Escalated)
	Responder       string    // Identifier of the responder
}

// ThreatEvent represents an event related to a detected threat.
type ThreatEvent struct {
	EventID         string    // Unique identifier for the event
	Timestamp       time.Time // Timestamp of the event
	ThreatType      string    // Type of threat (e.g., DDoS, Malware)
	Description     string    // Detailed description of the event
	AffectedSystems []string  // Systems affected by the threat
	MitigationSteps []string  // Mitigation steps taken
}

// IncidentStatus represents the status of an incident.
type IncidentStatus struct {
	IncidentID      string    // Unique identifier for the incident
	ReportedAt      time.Time // Timestamp when the incident was reported
	Status          string    // Current status (e.g., Open, Resolved, Escalated)
	Severity        string    // Severity level of the incident
	ResolutionSteps []string  // Steps taken to resolve the incident
}

// HealthVerification represents the verification of system health.
type HealthVerification struct {
	SystemID          string    // Unique identifier for the system
	LastVerified      time.Time // Last time health was verified
	HealthStatus      string    // Current health status (e.g., Healthy, Degraded)
	VerificationLog   []string  // Log of verification steps
	VerificationScore int       // Score of the system's health verification
}

// SecurityPolicyRecord represents a record of a security level policy change.
type SecurityPolicyRecord struct {
	Policy    string
	Timestamp string
}

// ApplicationHardeningEvent records the status of application hardening.
type ApplicationHardeningEvent struct {
	Status    string
	Timestamp time.Time
}

type Alert struct {
	AlertID   string
	Status    string
	CreatedAt time.Time
}

type AlertStatusLog struct {
	AlertID  string
	Status   string
	LoggedAt time.Time
}

// Validator represents a validator registered in the ledger.
type Validator struct {
	ID           string
	RegisteredAt time.Time
}

type SuspiciousActivity struct {
	Details   string
	Timestamp string
}

type HealthLogEntry struct {
	Timestamp string
	Message   string
}

type DataTransferMetrics struct {
	RateMBps     int       `json:"rate_mbps"`
	PeakRateMBps int       `json:"peak_rate_mbps"`
	Timestamp    time.Time `json:"timestamp"`
}

type HealthStatusLog struct {
	Status    string
	Timestamp time.Time
}

// DataTransmission represents the transmission of data between entities.
type DataTransmission struct {
	ID            string
	Sender        string
	Receiver      string
	Data          map[string]interface{}
	TransmittedAt time.Time
}

type IncidentActivation struct {
	IncidentID string
	Timestamp  string
}

type IncidentReport struct {
	IncidentID string
	Details    string
	Timestamp  string
}

type IncidentResolution struct {
	IncidentID       string
	ResolutionStatus string
	Timestamp        string // Store timestamp as a string
}

type ThreatLevel struct {
	Level     int
	Timestamp string // Store timestamp as a string
}

// Supporting structs for mitigation plans, anomalies, and protocol deviations
type MitigationPlan struct {
	PlanID        string
	Description   string
	Activated     bool
	Effectiveness string
	SetAt         time.Time
}

// TrafficAnomaly represents details of a detected traffic anomaly
type TrafficAnomaly struct {
	Description string
	SourceIP    string
	DetectedAt  time.Time
	Severity    string
}

// TrafficData represents traffic metrics for analysis
type TrafficData struct {
	SourceIP        string
	RequestCount    int
	FailedLogins    int
	AvgRequestRate  float64 // Requests per second
	PeakRequestRate float64 // Peak request rate
	Timestamps      []time.Time
}

type TrafficPattern struct {
	PatternID  string
	Details    string
	RecordedAt time.Time
}

type ProtocolDeviation struct {
	DeviationID string
	Description string
	DetectedAt  time.Time
}

type APIUsageRecord struct {
	APIID      string
	UsageCount int
	LastUsedAt time.Time
}

type ConsensusAnomaly struct {
	AnomalyID   string
	Description string
	DetectedAt  time.Time
}

type SuspiciousActivityRecord struct {
	ActivityID  string
	Description string
	DetectedAt  time.Time
}

type IncidentProtocol struct {
	ProtocolID    string
	Description   string
	Activated     bool
	ActivatedAt   time.Time
	DeactivatedAt time.Time
}

type IncidentEvent struct {
	EventID     string
	Description string
	RecordedAt  time.Time
	Resolution  string
}

type RetentionPolicy struct {
	PolicyID       string
	Description    string
	SetAt          time.Time
	ComplianceDate time.Time
}

type ComplianceAudit struct {
	AuditID     string
	Details     string
	PerformedAt time.Time
	Status      string
}

type EscalationProtocol struct {
	ProtocolID  string
	Description string
	SetAt       time.Time
}

type IsolationIncident struct {
	IncidentID  string
	Description string
	IsolatedAt  time.Time
}

// Supporting structs for threats, activities, rate limits, health metrics, and maintenance events
type ThreatDetails struct {
	ThreatID    string
	Description string
	Severity    int
	DetectedAt  time.Time
}

type RateLimitPolicy struct {
	PolicyID    string
	MaxRequests int
	Period      time.Duration
	SetAt       time.Time
}

type TransferRateLimit struct {
	LimitID     string
	MaxTransfer int
	SetAt       time.Time
}

type DataTransferMonitor struct {
	MonitorID  string
	DataAmount int
	RecordedAt time.Time
}

type HealthMetric struct {
	MetricID   string
	Value      int
	RecordedAt time.Time
}

type MaintenanceEvent struct {
	EventID     string
	Description string
	PerformedAt time.Time
}

type HealthEvent struct {
	EventID     string
	Description string
	RecordedAt  time.Time
}

type BackupEvent struct {
	BackupID    string
	Description string
	PerformedAt time.Time
}

type SessionTimeoutRecord struct {
	SessionID string
	TimeoutAt time.Time
}

type MigrationCompliance struct {
	MigrationID      string
	Details          string
	ComplianceStatus string
	RecordedAt       time.Time
}

type AccessLimit struct {
	NodeID      string
	MaxAccesses int
	SetAt       time.Time
}

type AccessFrequency struct {
	AccessID   string
	Frequency  int
	LastAccess time.Time
}

type SecurityLevelPolicy struct {
	PolicyID string
	Level    string
	SetAt    time.Time
}

// TrapManager manages traps for detecting, logging, and responding to suspicious activities.
type TrapManager struct {
	TrapID          string                 // Unique identifier for the trap.
	TrapType        string                 // Type of the trap (e.g., anomaly detection, honeypot, etc.).
	ActivationTime  time.Time              // Timestamp when the trap was activated.
	IsActive        bool                   // Indicates if the trap is currently active.
	TriggerCount    int                    // Number of times the trap has been triggered.
	TriggerLogs     []TrapTriggerLog       // Logs of all triggers associated with the trap.
	ResponseActions []string               // Actions to be executed when the trap is triggered.
	CreatedBy       string                 // Identifier of the creator of the trap.
	Metadata        map[string]interface{} // Additional metadata related to the trap.
}

// TrapTriggerLog represents a log entry for a trap trigger.
type TrapTriggerLog struct {
	Timestamp      time.Time // Time when the trap was triggered.
	TriggeredBy    string    // Identifier of the entity or process that triggered the trap.
	TriggerReason  string    // Reason for the trigger (e.g., suspicious activity).
	ResponseStatus string    // Status of the response to the trigger (e.g., executed, pending).
	LogDetails     string    // Additional details about the trigger event.
}

// AlertManager oversees alert generation, escalation policies, and notification dispatch.
type AlertManager struct {
	AlertID                string           // Unique identifier for the alert.
	AlertType              string           // Type of the alert (e.g., critical, warning, informational).
	AlertDescription       string           // Detailed description of the alert.
	AlertPriority          int              // Priority level of the alert (e.g., 1 = High, 2 = Medium, 3 = Low).
	AffectedComponents     []string         // List of system components affected by the alert.
	NotificationRecipients []string         // List of recipients to notify when the alert is generated.
	IsAcknowledged         bool             // Indicates whether the alert has been acknowledged.
	AcknowledgedBy         string           // Identifier of the entity that acknowledged the alert.
	AcknowledgedAt         time.Time        // Timestamp of when the alert was acknowledged.
	EscalationPolicy       EscalationPolicy // Policy for escalating unacknowledged alerts.
	AlertLogs              []AlertLog       // Logs for tracking alert lifecycle and actions.
}

// EscalationPolicy defines the policy for escalating unacknowledged alerts.
type EscalationPolicy struct {
	EscalationLevel    int           // Current escalation level.
	MaxEscalationLevel int           // Maximum allowed escalation level.
	EscalationActions  []string      // Actions to perform during escalation.
	EscalationInterval time.Duration // Time interval between escalation levels.
	EscalationContacts []string      // List of contacts for each escalation level.
}

// AlertLog represents a log entry for an alert's lifecycle.
type AlertLog struct {
	Timestamp       time.Time // Time of the logged action.
	LogType         string    // Type of log (e.g., created, acknowledged, resolved).
	ActionPerformed string    // Action taken related to the alert.
	PerformedBy     string    // Identifier of the entity that performed the action.
	LogDetails      string    // Additional details about the action or event.
}

// ************** Authorization Structs **************

// AuthorizationLog represents a record of an authorization event within the ledger.
type AuthorizationLog struct {
	LogID       string    // Unique ID for the log entry
	OperationID string    // ID of the operation being approved or queried
	SignerID    string    // ID of the signer or user involved in the action
	Action      string    // Type of action, e.g., "Approved", "Rejected"
	Timestamp   time.Time // Time the action was recorded
	IPAddress   string    // IP address from where the action was initiated
}

// UnauthorizedAccessLog represents a log of unauthorized access attempts.
type UnauthorizedAccessLog struct {
	LogID           string    // Unique identifier for the log entry
	AttemptID       string    // ID of the unauthorized access attempt
	EntityID        string    // ID of the entity that attempted unauthorized access
	Timestamp       time.Time // Timestamp of the unauthorized access attempt
	Reason          string    // Reason for marking access as unauthorized
	Location        string    // Physical or virtual location of the attempt
	AffectedSystem  string    // System targeted by the unauthorized access
	ResponseActions []string  // Actions taken in response to the attempt
}

// PublicKeyRecord represents details about a public key used for authorization.
type PublicKeyRecord struct {
	KeyID           string    // Unique identifier for the public key
	OwnerID         string    // ID of the entity that owns the public key
	Key             string    // The public key itself
	CreatedAt       time.Time // Timestamp of when the key was created
	ExpiresAt       time.Time // Expiration time for the public key
	Usage           string    // Purpose of the key (e.g., Authentication, Encryption)
	AssociatedRoles []string  // Roles associated with this key
	Revoked         bool      // Whether the key has been revoked
	RevokedAt       time.Time // Timestamp when the key was revoked, if applicable
}

// SignerActivityRecord logs activities performed by authorized signers.
type SignerActivityRecord struct {
	ActivityID      string    // Unique identifier for the activity log
	SignerID        string    // ID of the authorized signer
	Timestamp       time.Time // Timestamp of the activity
	ActionPerformed string    // Action performed by the signer
	TargetEntity    string    // Entity or system targeted by the action
	Status          string    // Status of the activity (e.g., Successful, Failed)
	Details         string    // Additional details about the activity
	Logs            []string  // Associated logs or outputs from the activity
}

// RoleChangeRecord tracks changes to roles assigned to entities.
type RoleChangeRecord struct {
	ChangeID       string    // Unique identifier for the role change
	EntityID       string    // ID of the entity whose role was changed
	OldRole        string    // The previous role of the entity
	NewRole        string    // The new role assigned to the entity
	ChangedBy      string    // ID of the entity or user who made the change
	Timestamp      time.Time // Timestamp of the role change
	Reason         string    // Reason for the role change
	ApprovalStatus string    // Approval status of the change (e.g., Approved, Pending)
	EffectiveDate  time.Time // Date when the role change takes effect
}

// MicrochipAccessAttempt represents an attempt to access a system using a microchip.
type MicrochipAccessAttempt struct {
	AttemptID       string    // Unique identifier for the access attempt
	MicrochipID     string    // ID of the microchip used in the attempt
	EntityID        string    // ID of the entity associated with the microchip
	Timestamp       time.Time // Timestamp of the access attempt
	AccessType      string    // Type of access requested (e.g., Login, Read, Write)
	Location        string    // Physical or virtual location of the attempt
	Outcome         string    // Outcome of the attempt (e.g., Successful, Failed)
	FailureReason   string    // Reason for failure, if applicable
	ResponseActions []string  // Actions taken in response to the access attempt
}

// AuthorizedSigner represents an authorized signer entry in the ledger.
type AuthorizedSigner struct {
	SignerID    string        // Unique identifier of the signer
	Permissions PermissionSet // Permissions assigned to the signer
	AddedAt     time.Time     // Timestamp when the signer was added
}

// BiometricData represents the biometric information of a user.
type BiometricData struct {
	FingerprintHash []byte // Stores a hash of the user's fingerprint data
	FaceIDHash      []byte // Stores a hash of the user's facial recognition data
	IrisScanHash    []byte // Stores a hash of the user's iris scan data
	VoicePrintHash  []byte // Stores a hash of the user's voice print data
}

// BiometricUpdate represents an update to a user's biometric data.
type BiometricUpdate struct {
	UserID        string    // ID of the user
	EncryptedData []byte    // Encrypted biometric data
	UpdatedAt     time.Time // Timestamp of the update
}

// EntityType defines the type of entity in the system, such as User, Device, or Service.
type EntityType struct {
	TypeID   string // Unique identifier for the entity type
	Name     string // Name of the entity type, e.g., "User", "Device", "Service"
	Category string // Optional category for further classification
}

// IdentityData represents identity-related data for an entity, such as personal details or credentials.
type IdentityData struct {
	EntityID     string    // Unique identifier for the entity
	FullName     string    // Full name of the entity owner (if applicable)
	DateOfBirth  time.Time // Date of birth (optional, if applicable)
	NationalID   string    // National ID or similar identifier (optional)
	Address      string    // Address information (optional)
	ContactInfo  string    // Contact information, such as email or phone
	RegisteredAt time.Time // Date and time of registration
}

type MicrochipAuthorization struct {
	ChipID              string
	DeviceID            string // ID of the device associated with the microchip
	AuthorizedUser      string
	EncryptedKey        []byte
	IssuedAt            time.Time
	ExpiresAt           time.Time
	AuthorizationStatus string
	AuthorizationLevel  string // Added AuthorizationLevel if missing
}

// PermissionSet defines the permissions associated with an authorized signer in the system.
type PermissionSet struct {
	CanApproveTransactions    bool // Permission to approve transactions
	CanModifyLedgerEntries    bool // Permission to modify ledger entries
	CanAccessSensitiveData    bool // Permission to access sensitive data within the system
	CanAddOrRemoveSigners     bool // Permission to add or remove authorized signers
	CanSetAuthorizationLevel  bool // Permission to set authorization levels for other users
	CanFlagSuspiciousActivity bool // Permission to flag or review suspicious activities
}

// BiometricRegistration represents a biometric registration entry in the ledger.
type BiometricRegistration struct {
	UserID         string    // ID of the user
	EncryptedData  []byte    // Encrypted biometric data
	RegistrationAt time.Time // Timestamp of the registration
}

// DeviceInfo contains metadata and identifiers for an authorized device.
type DeviceInfo struct {
	DeviceName      string    // Friendly name for the device
	DeviceType      string    // Type of device, e.g., "Mobile", "Laptop", "Tablet"
	OperatingSystem string    // OS of the device, e.g., "iOS", "Android", "Windows"
	OSVersion       string    // Version of the operating system
	SerialNumber    string    // Device serial number for unique identification
	IPAddress       string    // IP address of the device
	MACAddress      string    // MAC address for network identification
	LastAccessedAt  time.Time // Timestamp of the last access by this device
}

type TemporaryAccess struct {
	DeviceID     string    // Unique identifier for the device granted temporary access
	AuthorizedAt time.Time // Timestamp when access was authorized
	ExpiresAt    time.Time // Expiration timestamp for temporary access
	AuthorizedBy string    // ID of the user or system that authorized the access
	AccessLevel  string    // Level of access granted (e.g., "read-only", "full-access")
	Reason       string    // Reason for granting temporary access
}

type AccessLog struct {
	AccessID      string    // Unique identifier for the access attempt log
	UserID        string    // ID of the user attempting access
	DeviceID      string    // ID of the device involved in the access attempt
	Timestamp     time.Time // Timestamp of the access attempt
	Success       bool      // Outcome of the attempt (true for success, false for failure)
	IPAddress     string    // IP address from where the attempt was made
	AccessType    string    // Type of access attempted (e.g., "login", "data retrieval")
	FailureReason string    // Reason for failure if the attempt was unsuccessful
}

type RoleChangeLog struct {
	LogID      string    // Unique identifier for the role change log entry
	UserID     string    // ID of the user whose role is being changed
	ChangedBy  string    // ID of the user or system making the role change
	OldRole    string    // Previous role of the user
	NewRole    string    // New role of the user
	ChangedAt  time.Time // Timestamp of when the role change was made
	Reason     string    // Reason for the role change
	ApprovedBy string    // ID of the approver if needed for compliance
}

type AuthorizationConstraints struct {
	ConstraintID       string    // Unique identifier for the constraint
	UserID             string    // ID of the user subject to these constraints
	MaxAccessLevel     string    // Maximum level of access the user can have
	AccessTimeLimits   []string  // Specific times or time ranges allowed for access (e.g., "09:00-17:00")
	AccessFrequency    int       // Maximum allowed access attempts per day or hour
	DeviceRestrictions []string  // List of allowed devices for the user
	ExpiryDate         time.Time // Expiration date of the constraints
	CreatedAt          time.Time // Timestamp when the constraints were created
	CreatedBy          string    // ID of the user or system that created the constraints
}

type AuthorizationKeyReset struct {
	ResetID     string    // Unique identifier for the key reset request
	UserID      string    // ID of the user whose key is being reset
	InitiatedBy string    // ID of the user or system initiating the reset
	ApprovedBy  string    // ID of the user or system approving the reset
	RequestedAt time.Time // Timestamp when the reset was requested
	ApprovedAt  time.Time // Timestamp when the reset was approved
	Reason      string    // Reason for the reset (e.g., "compromised key", "expiration")
	NewKeyHash  string    // Hash of the new authorization key (if applicable)
	Status      string    // Status of the reset request (e.g., "Pending", "Approved", "Completed")
}

// DelegatedAccess represents a record of delegated access to a device.
type DelegatedAccess struct {
	DeviceID    string    // ID of the device being accessed
	DelegatorID string    // ID of the user granting access
	DelegateID  string    // ID of the user receiving delegated access
	GrantedAt   time.Time // Timestamp when access was granted
	ExpiresAt   time.Time // Timestamp when access expires
}

// BiometricAccessLog records the outcome of a biometric authentication attempt.
type BiometricAccessLog struct {
	UserID    string
	Timestamp time.Time
	Success   bool
}

// AuthorizationEvent represents an event related to authorization actions.
type AuthorizationEvent struct {
	EventID   string    // Unique identifier for the authorization event
	Action    string    // Type of action, e.g., "Added", "Removed", "Modified"
	UserID    string    // User ID associated with the event
	Timestamp time.Time // Timestamp of the event
	Details   string    // Additional details about the event
}

// TrustedParty represents a trusted party in the system.
type TrustedParty struct {
	PartyID string    // Unique identifier for the trusted party
	AddedAt time.Time // Timestamp of when the party was added as trusted
	Flagged bool      // Indicates if the party is flagged for review
}

// AuthorizationData represents the structure for user authorization data.
type AuthorizationData struct {
	UserID             string    // Unique identifier for the user
	AuthorizationLevel int       // Authorization level for the user
	SetAt              time.Time // Timestamp of when the authorization level was set
}

type SignerActivity struct {
	SignerID  string
	Action    string
	Timestamp time.Time
}

type PermissionRequest struct {
	RequestID   string
	RequestedBy string
	Permission  string
	Status      string
	RequestedAt time.Time
}

// AuditTrail records details about an audit trail for traceability.
type AuditTrail struct {
	TrailID       string    // Unique identifier for the audit trail entry
	EventType     string    // Type of event (e.g., "transaction", "access", "authorization")
	UserID        string    // User ID associated with the action
	NodeID        string    // ID of the node where the event occurred
	Timestamp     time.Time // Timestamp of the event
	ActionDetails string    // Detailed description of the action
	Status        string    // Current status of the audit trail (e.g., "Enabled", "Disabled")
	OperationID   string    // Operation ID associated with the audit trail
}

// UnauthorizedAccess represents an unauthorized access attempt in the system.
type UnauthorizedAccess struct {
	OperationID string    // ID of the operation involved
	SignerID    string    // ID of the signer attempting access
	Details     string    // Details of the unauthorized access
	Timestamp   time.Time // Timestamp of the access attempt
}

// Supporting structs for permission revocations, sub-block validation authorization, biometrics, etc.
type PermissionRevocation struct {
	PermissionID string
	RevokedBy    string
	Reason       string
	RevokedAt    time.Time
}

type SubBlockValidationAuth struct {
	AuthID       string
	ValidatorID  string
	Status       string
	AuthorizedAt time.Time
}

type AccessAttempt struct {
	AttemptID   string
	EntityID    string
	AttemptedAt time.Time
	Result      string
}

type BiometricRecord struct {
	BiometricID  string
	DataHash     string
	RegisteredAt time.Time
	LastUpdated  time.Time
}

type DeviceAuthorization struct {
	DeviceID           string
	AuthorizedBy       string
	AuthorizationLevel string
	AuthorizedAt       time.Time
}

type DelegatedAccessRecord struct {
	DelegationID string
	GrantedBy    string
	GrantedTo    string
	AccessLevel  string
	GrantedAt    time.Time
}

type TemporaryAccessRecord struct {
	AccessID    string
	EntityID    string
	AccessLevel string
	ExpiresAt   time.Time
}

type RoleChange struct {
	RoleID    string
	ChangedBy string
	NewRole   string
	ChangedAt time.Time
}

type AuthorizationConstraint struct {
	ConstraintID string
	Description  string
	AppliedBy    string
	AppliedAt    time.Time
}

type KeyResetRecord struct {
	KeyID   string
	ResetBy string
	ResetAt time.Time
}

type MultiSigWallet struct {
	WalletID     string
	Owners       []string
	RequiredSigs int
	CreatedAt    time.Time
}

type RoleAssignment struct {
	AssignmentID string
	Role         string
	AssignedTo   string
	AssignedBy   string
	AssignedAt   time.Time
}

type Privilege struct {
	PrivilegeID string
	RoleID      string
	Permissions []string
	SetAt       time.Time
}

type TimeBasedAuthorization struct {
	AuthorizationID string
	EntityID        string
	AccessLevel     string
	ValidFrom       time.Time
	ExpiresAt       time.Time
}

type SignerPriority struct {
	SignerID      string
	PriorityLevel int
	SetAt         time.Time
}

type UserStatus struct {
	UserID    string
	Status    string
	UpdatedAt time.Time
}

type FeedbackComment struct {
	CommentID  string
	FeedbackID string
	UserID     string
	Content    string
	CreatedAt  time.Time
}

type Vote struct {
	VoteID  string
	UserID  string
	PollID  string
	Option  string
	VotedAt time.Time
}

type Favorite struct {
	UserID  string
	PostID  string
	AddedAt time.Time
}

type Follow struct {
	FollowerID  string
	FollowingID string
	FollowedAt  time.Time
}

type BlockedUser struct {
	BlockerID string
	BlockedID string
	BlockedAt time.Time
}

type Mute struct {
	MuterID string
	MutedID string
	MutedAt time.Time
}

// RoleManager manages user roles, their associated permissions, and hierarchical relationships.
type RoleManager struct {
	RoleID        string                 // Unique identifier for the role.
	RoleName      string                 // Human-readable name for the role.
	Permissions   []string               // List of permissions assigned to the role.
	AssignedUsers []string               // List of users assigned to this role.
	CreatedAt     time.Time              // Timestamp when the role was created.
	CreatedBy     string                 // Identifier of the entity that created the role.
	UpdatedAt     time.Time              // Timestamp of the last update to the role.
	ParentRoleID  string                 // Parent role for hierarchical permissions, if applicable.
	Metadata      map[string]interface{} // Additional metadata related to the role.
}

// RoleAssignmentLog keeps a history of role assignments and changes.
type RoleAssignmentLog struct {
	UserID      string    // Identifier of the user affected.
	RoleID      string    // Identifier of the role involved.
	Action      string    // Action performed (e.g., assigned, removed, updated).
	PerformedBy string    // Identifier of the entity performing the action.
	Timestamp   time.Time // Time the action was performed.
	LogDetails  string    // Additional details about the action.
}

// AccessManager handles access policies, rights, and restrictions across the system.
type AccessManager struct {
	PolicyID         string                 // Unique identifier for the access policy.
	PolicyName       string                 // Human-readable name for the access policy.
	AccessRules      []AccessRule           // List of rules defining access rights and restrictions.
	AffectedEntities []string               // List of entities (users, roles, or components) affected by the policy.
	CreatedAt        time.Time              // Timestamp when the policy was created.
	CreatedBy        string                 // Identifier of the entity that created the policy.
	UpdatedAt        time.Time              // Timestamp of the last update to the policy.
	IsActive         bool                   // Indicates if the access policy is currently active.
	Metadata         map[string]interface{} // Additional metadata related to the policy.
}

// AccessRule defines the rules for granting or restricting access.
type AccessRule struct {
	RuleID     string                 // Unique identifier for the access rule.
	Resource   string                 // Resource or component the rule applies to.
	Action     string                 // Action allowed or restricted (e.g., read, write, execute).
	Conditions map[string]interface{} // Conditions for the rule to apply.
	IsAllowed  bool                   // Indicates if the action is allowed (true) or restricted (false).
	CreatedAt  time.Time              // Timestamp when the rule was created.
	CreatedBy  string                 // Identifier of the entity that created the rule.
}

// ************** Blockchain Structs **************

// PrunedBlockchain represents a blockchain where older blocks and data have been removed to save storage space.
type PrunedBlockchain struct {
	LatestBlockHash   string            // Hash of the most recent block in the pruned blockchain
	BlockHeight       int               // Current block height of the pruned chain
	RetainedBlocks    map[string]*Block // Map of retained recent blocks, identified by block hash
	PrunedBlockHeight int               // Block height below which blocks have been pruned
	Ledger            *Ledger           // Reference to the ledger to ensure consistency across pruned and full chains
	SnapshotTimestamp time.Time         // Timestamp of the last snapshot taken before pruning
	ValidationHash    string            // Hash representing the current state of the pruned blockchain for validation
}

// BlockListener represents an entity listening to block-related events.
type BlockListener struct {
	ID           string                 // Unique identifier for the listener
	ListenerType string                 // Type of listener (e.g., Validator, Miner, Application)
	CallbackURL  string                 // URL or endpoint for event notifications
	Events       []string               // List of events the listener subscribes to (e.g., "BlockCreated", "BlockFinalized")
	Active       bool                   // Indicates if the listener is active
	LastUpdated  time.Time              // Timestamp of the last update to the listener
	Metadata     map[string]interface{} // Additional metadata about the listener
}

// Block represents a blockchain block.
type Block struct {
	BlockID     string     // Unique identifier for the block
	Index       int        // Block index
	Timestamp   time.Time  // Block creation time
	SubBlocks   []SubBlock // Sub-blocks inside the block
	PrevHash    string     // Previous block's hash
	Hash        string     // Current block's hash
	Nonce       int        // Nonce for PoW
	Difficulty  int        // Difficulty level for PoW
	MinerReward float64    // Reward given to the miner
	Validators  []string   // List of validators who contributed to the block validation
	Status      string     // Block status (new field)

}

// BlockSummary represents a simplified version of the block for light nodes.
type BlockSummary struct {
	BlockID   string    // Unique identifier for the block
	Index     int       // Block index
	Hash      string    // Current block's hash
	PrevHash  string    // Hash of the previous block
	Timestamp time.Time // Block creation time
	Status    string    // Block status (e.g., "Confirmed", "Pending")
}

// Blockchain represents the main blockchain with all blocks
type Blockchain struct {
	Chain               []Block       // The blockchain itself
	PendingTransactions []Transaction // Transactions waiting to be included in a block
	SubBlockChain       SubBlockChain // Sub-blockchain that handles sub-blocks
	Validators          []string      // List of validators
	OwnerWallet         string        // The owner's wallet address
	mutex               sync.Mutex    // Mutex for thread-safe operations
	Ledger              *Ledger       // Ledger to store blocks and transactions
}

// SubBlockChain holds a list of sub-blocks and manages them
type SubBlockChain struct {
	SubBlocks  []SubBlock // The actual chain of sub-blocks
	Validators []string   // List of validators for PoS
	Ledger     *Ledger    // Pointer to the ledger for storing validated sub-blocks
}

// SubBlockManager manages the creation, validation, and addition of sub-blocks to the ledger.
type SubBlockManager struct {
	LedgerInstance *Ledger     // Reference to the ledger to store sub-blocks
	Encryption     *Encryption // Encryption instance for data encryption
	mutex          sync.Mutex  // Mutex for thread-safe operations
}

type SubBlock struct {
	SubBlockID   string        // Sub-block ID (new field)
	Index        int           // Sub-block index
	Timestamp    time.Time     // Sub-block creation time (keeping it as time.Time here)
	Transactions []Transaction // List of transactions in the sub-block
	Validator    string        // Validator who validated the sub-block
	PrevHash     string        // Previous sub-block's hash
	Hash         string        // Current sub-block's hash
	PoHProof     PoHProof
	Status       string // Block status (new field)
	Signature    string // Signature from the validator (new field)

}

// ************** Consensus Structs **************
// SynnergyConsensus manages the combination of PoH, PoS, and PoW for sub-block validation
type SynnergyConsensus struct {
	PoH            *PoH           // Proof of History mechanism
	PoS            *PoS           // Proof of Stake mechanism
	PoW            *PoW           // Proof of Work mechanism
	RewardManager  *RewardManager // Reward Manager for PoH, PoS, and PoW rewards
	LedgerInstance *Ledger        // Ledger instance for tracking rewards, blocks, and transactions
	SubBlockCount  int            // Keeps track of sub-blocks processed for the current block
	Validators     []Validator    // List of validators in the system
	Encryption     *Encryption    // Encryption system for proof data

}

// PoH represents the Proof of History mechanism, now integrated with the Reward Manager and Ledger
type PoH struct {
	State          PoHState       // The current state of the PoH system
	LedgerInstance *Ledger        // Ledger to store PoH proofs
	RewardManager  *RewardManager // Reward manager for rewarding validators
}

// PoHState represents the state of the Proof of History (PoH) mechanism
type PoHState struct {
	Sequence      int       // The current sequence number in PoH
	LastTimestamp time.Time // Timestamp of the last PoH entry
	LastHash      string    // Hash of the previous PoH entry
}

// PoHProof represents the proof generated by PoH
type PoHProof struct {
	Sequence  int       // The sequence number for this proof
	Timestamp time.Time // The timestamp for this PoH proof
	Hash      string    // Hash generated by PoH for this entry
}

// PoS represents the Proof of Stake mechanism, integrated with the Ledger and Reward Manager
type PoS struct {
	State          PoSState       // Current state of the PoS system
	LedgerInstance *Ledger        // Ledger instance for tracking stake and rewards
	RewardManager  *RewardManager // Reward Manager to handle validator rewards
}

// PoSState represents the current state of the PoS system
type PoSState struct {
	Validators   []Validator // List of all validators participating in PoS
	TotalStake   float64     // Total amount of SYNN staked in the network
	LastSelected string      // Last validator selected to validate a sub-block
	Epoch        int         // Epoch number for selecting new validators
	LastUpdated  time.Time   // Timestamp of the last validator selection
}

// PoW represents the Proof of Work mechanism, integrated with the Ledger
type PoW struct {
	State          PoWState // The current state of the PoW system
	LedgerInstance *Ledger  // Ledger instance for recording blocks and rewards
}

// PoWState holds the current state of the Proof of Work system
type PoWState struct {
	Difficulty   int     // Current difficulty for PoW
	BlockReward  float64 // Block reward for the miner
	MinerAddress string  // Address of the miner
	Epoch        int     // Current epoch
	LastHash     string  // Hash of the last mined block
}

// RewardManager handles the distribution of rewards for PoH, PoS, and PoW
type RewardManager struct {
	PoSRewardRate      float64            // Percentage reward for staking in PoS
	PoHRewardRate      float64            // Percentage reward for participating in PoH
	PoWInitialReward   float64            // Initial block reward for PoW
	PoWHalvingInterval int                // Number of blocks before PoW reward halves
	CurrentBlockCount  int                // Current block count to track PoW halving
	RewardPool         float64            // Combined pool for validator and miner rewards
	LedgerInstance     *Ledger            // Instance of the ledger for reward tracking
	PunishmentManager  *PunishmentManager // Reference to the PunishmentManager for enforcing penalties
	mutex              sync.Mutex         // Mutex for thread-safe operations
}

type Punishment struct {
	Amount    float64
	Entity    string
	Timestamp time.Time
}

// PunishmentManager manages punishments for validators and miners
type PunishmentManager struct {
	PoSPunishmentThreshold  float64                 // Threshold for PoS violations (e.g., 24 hours of downtime)
	PoHPunishmentThreshold  float64                 // Threshold for PoH inactivity (e.g., 1 missed participation cycle)
	PoWPunishmentThreshold  float64                 // Threshold for PoW failure (e.g., 5 failed block attempts)
	PoSPunishmentRate       float64                 // Punishment rate for PoS violations (e.g., percentage reduction in stake)
	PoHPunishmentRate       float64                 // Punishment rate for PoH violations
	PoWPunishmentRate       float64                 // Punishment rate for PoW failures
	PunishmentHistory       map[string][]Punishment // Punishment history for each entity
	PunishmentResetInterval time.Duration           // Time interval for punishment reset (e.g., 90 days)
	LedgerInstance          *Ledger                 // Reference to the ledger for recording punishments
	mutex                   sync.Mutex              // Mutex for thread-safe operations
}

// ConsensusState represents the state of the consensus mechanism.
type ConsensusState struct {
	PoHProofs            []PoHProof            // List of PoH proofs generated in the network
	ValidatorPunishments map[string]Punishment // Punishments applied to validators, mapped by validator ID
	ValidatorRewards     map[string]float64    // Rewards assigned to validators, mapped by validator ID
	MinerRewards         map[string]float64    // Rewards assigned to miners, mapped by miner ID
	ParticipantRewards   map[string]float64    // Rewards assigned to non-validator participants, mapped by participant ID
	SubBlockCount        int                   // Number of sub-blocks processed
	FinalizedSubBlocks   []SubBlock            // List of finalized sub-blocks
	FinalizedBlocks      []Block               // List of finalized blocks
	PoWState             PoWState              // Current state of Proof of Work (PoW)
	PoSState             PoSState              // Current state of Proof of Stake (PoS)
	PoHState             PoHState              // Current state of Proof of History (PoH)
}

// ************** Charity Structs **************

// CharityPoolManagement handles the distribution of transaction fees into internal and external charity pools
type CharityPoolManagement struct {
	mutex               sync.Mutex
	InternalPoolBalance float64 // Balance for the internal charity pool
	ExternalPoolBalance float64 // Balance for the external charity pool
	LedgerInstance      *Ledger // Ledger instance for tracking pool activities
}

// CharityPool represents the external and internal charity pools and manages their balances
type CharityPool struct {
	externalPool   float64    // External charity pool balance
	internalPool   float64    // Internal charity pool balance
	totalBalance   float64    // Total balance to be distributed between both pools
	LedgerInstance *Ledger    // Ledger instance for tracking pool activities
	mutex          sync.Mutex // Mutex for thread-safe operations
}

// CharityProposal represents a charity that enters into the external charity pool
type CharityProposal struct {
	CharityID     string    // Unique ID for the charity
	Name          string    // Charity name
	CharityNumber string    // Registered charity number
	Description   string    // Charity description
	Website       string    // Charity website
	Addresses     []string  // Encrypted addresses to receive funds
	CreatedAt     time.Time // Timestamp of when the charity entered the pool
	VoteCount     int       // Number of votes received
	IsValid       bool      // Is the charity valid
}

// ExternalCharityPoolManager manages the external charity proposal process and fund distribution
type ExternalCharityPoolManager struct {
	mutex               sync.Mutex
	CurrentCycle        []*CharityProposal          // Charities in the current 90-day cycle
	CharityEntries      map[string]*CharityProposal // Entries for the current round
	ProposalStart       time.Time                   // Start of the proposal submission
	VotingEnd           time.Time                   // End of the voting period
	LedgerInstance      *Ledger                     // Ledger instance for tracking charity activity
	ExternalPoolBalance float64                     // Balance of the external charity pool
}

// InternalCharityPool manages the internal charity pool, distributing funds every 24 hours
type InternalCharityPool struct {
	mutex           sync.Mutex
	PoolBalance     float64            // Current balance of the internal charity pool
	WalletAddresses map[string]float64 // Map of charity wallet addresses to their respective balances
	OwnerAddress    string             // Blockchain owner's address (for access control)
	LedgerInstance  *Ledger            // Ledger instance for tracking all transactions and activities
	stopChan        chan bool          // Channel to stop the 24-hour distribution
}

// ************** Coin Structs **************

// ************** Community Engagement Structs **************

// ForumPost represents a post within the community forum.
type ForumPost struct {
	ID        string    // Unique ID of the post
	Author    string    // The author of the post
	Content   string    // The content of the post
	Timestamp time.Time // Time the post was created
	Replies   []Reply   // Replies to the post
	Hash      string    // Hash to ensure data integrity
}

// CommunityEvent represents a community event or activity.
type CommunityEvent struct {
	EventID      string            // Unique identifier for the event
	Title        string            // Title of the event
	Description  string            // Description of the event
	StartTime    time.Time         // Start time of the event
	EndTime      time.Time         // End time of the event
	Location     string            // Location of the event (optional for virtual events)
	Participants map[string]string // Map of participant IDs to their roles (e.g., attendee, organizer)
	IsVirtual    bool              // Indicates whether the event is virtual
	CreatedAt    time.Time         // Timestamp of event creation
	CreatedBy    string            // ID of the user who created the event
}

// Event struct to store information about an event in the blockchain network
type Event struct {
	ID              string    // Unique identifier for the event
	Name            string    // Name of the event
	Description     string    // Description of the event
	Location        string    // Location of the event
	Date            time.Time // Date and time of the event
	CreatorID       string    // User ID of the event creator
	MaxParticipants int       // Maximum number of participants allowed
	Participants    []string  // List of participant IDs
	CreatedAt       time.Time // Timestamp when the event was created
}

// Collection struct
type Collection struct {
	ID          string    // Unique identifier for the collection
	Name        string    // Name of the collection
	Description string    // Description of the collection
	OwnerID     string    // User ID of the collection's owner
	CreatedAt   time.Time // Timestamp for when the collection was created
	PostIDs     []string  // List of post IDs associated with the collection
}

// Reaction struct to store information about a reaction to a post or reply
type Reaction struct {
	ID        string    // Unique identifier for the reaction
	PostID    string    // ID of the post associated with the reaction (optional)
	ReplyID   string    // ID of the reply associated with the reaction (optional)
	UserID    string    // ID of the user who reacted
	Type      string    // Type of reaction (like, dislike, etc.)
	CreatedAt time.Time // Timestamp for when the reaction was created
}

// UserFavorites struct to store information about user's favorite posts
type UserFavorites struct {
	UserID    string   // ID of the user
	Favorites []string // List of post IDs marked as favorite by the user
}

// PostReport struct to store information about reported posts
type PostReport struct {
	PostID     string    // ID of the reported post
	ReporterID string    // ID of the user who reported the post
	Reason     string    // Reason for reporting
	Timestamp  time.Time // Time when the report was made
}

// Post struct to store information about posts
type Post struct {
	ID        string    // Unique identifier for the post
	AuthorID  string    // ID of the author who created the post
	Content   string    // Content of the post
	Tags      []string  // Tags associated with the post
	Timestamp time.Time // Timestamp when the post was created
	Upvotes   int       // Number of upvotes on the post
	Downvotes int       // Number of downvotes on the post
}

// UserProfile struct to store user profile information
type UserProfile struct {
	UserID   string   // Unique identifier for the user
	Username string   // Username of the user
	Bio      string   // Short bio or description
	Keywords []string // List of keywords for searchability
}

// UserBlock struct to store information about blocks between users
type UserBlock struct {
	RequesterID  string // ID of the user who set the block
	TargetUserID string // ID of the user who is blocked
}

// Reply struct to store information about replies to posts
type Reply struct {
	ID        string    // Unique identifier for the reply
	PostID    string    // ID of the post this reply is associated with
	AuthorID  string    // ID of the author who created the reply
	Content   string    // Content of the reply
	Timestamp time.Time // Timestamp when the reply was created
}

// ForumManager manages the community forum, posts, and replies.
type ForumManager struct {
	Posts          map[string]*ForumPost // All forum posts
	LedgerInstance *Ledger               // Reference to the ledger for audit and storage
	mutex          sync.Mutex            // Mutex for thread-safe operations
}

// Poll struct to store information about polls
type Poll struct {
	ID        string            // Unique identifier for the poll
	CreatorID string            // ID of the user who created the poll
	Question  string            // Poll question
	Options   []string          // Options for voting
	Expiry    time.Time         // Expiration time of the poll
	Open      bool              // Status indicating if the poll is open for voting
	Votes     map[string]int    // Vote counts for each option
	VoterList map[string]string // Records the option each user voted for to prevent double-voting
	CreatedAt time.Time         // Timestamp when the poll was created
}

// PrivateMessage struct to store information about private messages
type PrivateMessage struct {
	ID         string    // Unique identifier for the message
	SenderID   string    // ID of the user sending the message
	ReceiverID string    // ID of the user receiving the message
	Content    string    // Content of the message
	Timestamp  time.Time // Time when the message was sent
}

// FeedbackSystem manages the feedback from users.
type FeedbackSystem struct {
	Feedbacks      map[string]*Feedback // Store feedback data
	LedgerInstance *Ledger              // Ledger integration for secure feedback storage
	mutex          sync.Mutex           // Mutex for thread-safe operations
}

// Feedback struct to store information about user feedback
type Feedback struct {
	ID        string    // Unique identifier for the feedback
	UserID    string    // User ID of the feedback submitter
	Content   string    // Content of the feedback
	Submitted time.Time // Timestamp when the feedback was submitted
	Resolved  bool      // Flag indicating if the feedback is resolved
	Likes     int       // Number of likes for the feedback
	Dislikes  int       // Number of dislikes for the feedback
	Comments  []Comment // List of comments associated with the feedback
}

// Comment struct to store information about comments on feedback or posts
type Comment struct {
	ID        string    // Unique identifier for the comment
	UserID    string    // User ID of the commenter
	Content   string    // Content of the comment
	Submitted time.Time // Timestamp for when the comment was submitted
	Likes     int       // Number of likes the comment has received
	Dislikes  int       // Number of dislikes the comment has received
	ParentID  string    // Optional: ID of the parent comment for nested replies
}

// Report struct to store information about user reports
type Report struct {
	ID           string    // Unique identifier for the report
	ReporterID   string    // ID of the user who submitted the report
	ReportedUser string    // ID of the user being reported
	Reason       string    // Reason for the report
	DateReported time.Time // Date when the report was submitted
	Resolved     bool      // Flag indicating if the report is resolved
}

// User struct to store user information and status
type User struct {
	ID        string    // Unique identifier for the user
	Username  string    // Username of the user
	Email     string    // Email of the user
	Status    string    // Status of the user (e.g., "active", "banned")
	CreatedAt time.Time // Timestamp when the user was created
}

// ModerationLog struct to track moderation actions
type ModerationLog struct {
	ID         string    // Unique identifier for the log entry
	AdminID    string    // ID of the admin or moderator performing the action
	TargetID   string    // ID of the user or content being moderated
	Action     string    // Action taken (e.g., "ban", "unban", "moderate content")
	Reason     string    // Reason for the action
	DateLogged time.Time // Timestamp when the action was logged
}

// BugReport struct to store information about reported bugs
type BugReport struct {
	ID           string    // Unique identifier for the bug report
	UserID       string    // ID of the user reporting the bug
	Description  string    // Description of the bug
	DateReported time.Time // Timestamp when the bug was reported
	Resolved     bool      // Status indicating if the bug is resolved
}

// ************** Compliance Structs **************

// AMLSystem defines the Anti-Money Laundering (AML) system
type AMLSystem struct {
	SuspiciousActivityThreshold float64           // Threshold for suspicious activity
	BlockedWallets              map[string]bool   // List of blocked wallets
	ReportedTransactions        map[string]string // Map of reported transactions
	LedgerInstance              *Ledger           // Instance of the ledger for transaction logging
	mutex                       sync.Mutex        // Mutex for thread-safe operations
}

// AuditTrailEntry represents a single entry in the audit trail
type AuditTrailEntry struct {
	EventID   string    // Unique identifier for the event
	EventType string    // Type of event (transaction, system change, etc.)
	Timestamp time.Time // Time the event occurred
	UserID    string    // ID of the user who performed the action
	Details   string    // Description of the event
}

// PolicyManager oversees system-wide policies, including creation, enforcement, and lifecycle management.
type PolicyManager struct {
	PolicyID         string                 // Unique identifier for the policy.
	PolicyName       string                 // Human-readable name of the policy.
	PolicyType       string                 // Type of the policy (e.g., security, access, governance).
	Description      string                 // Description of the policy and its purpose.
	Rules            []PolicyRule           // List of rules associated with the policy.
	EnforcementLevel string                 // Level of enforcement (e.g., strict, moderate, lenient).
	CreatedAt        time.Time              // Timestamp when the policy was created.
	CreatedBy        string                 // Identifier of the entity that created the policy.
	UpdatedAt        time.Time              // Timestamp of the last update to the policy.
	IsActive         bool                   // Indicates whether the policy is currently active.
	Metadata         map[string]interface{} // Additional metadata about the policy.
}

// PolicyRule defines individual rules within a policy.
type PolicyRule struct {
	RuleID      string                 // Unique identifier for the rule.
	Resource    string                 // Resource or component the rule applies to.
	Action      string                 // Action allowed or restricted by the rule.
	Conditions  map[string]interface{} // Conditions that must be met for the rule to apply.
	IsMandatory bool                   // Indicates whether compliance with this rule is mandatory.
	CreatedAt   time.Time              // Timestamp when the rule was created.
	CreatedBy   string                 // Identifier of the entity that created the rule.
}

type DataRetentionPolicy struct {
	PolicyID        string        // Unique identifier for the policy
	Name            string        // Name of the retention policy
	Description     string        // Description of the policy
	RetentionPeriod time.Duration // Duration to retain data
	CreatedAt       time.Time     // Timestamp of policy creation
	UpdatedAt       time.Time     // Timestamp of last update
	CreatedBy       string        // User or system that created the policy
	IsActive        bool          // Whether the policy is currently active
}

type UserPrivacySetting struct {
	UserID                  string          // ID of the user associated with the setting
	AllowDataSharing        bool            // Indicates if the user allows data sharing
	DataVisibility          string          // Data visibility level (e.g., public, private, restricted)
	NotificationPreferences map[string]bool // Notification preferences (e.g., email, SMS)
	LastUpdated             time.Time       // Timestamp of the last update to settings
}

type RegulatoryReport struct {
	ReportID    string    // Unique identifier for the report
	Title       string    // Title of the regulatory report
	Description string    // Description or summary of the report
	CreatedAt   time.Time // Timestamp of report creation
	SubmittedBy string    // User or system that submitted the report
	Status      string    // Status of the report (e.g., pending, submitted, approved)
}

type SuspiciousTransaction struct {
	TransactionID string    // Unique identifier for the transaction
	DetectedBy    string    // System or user that flagged the transaction
	Reason        string    // Reason for marking the transaction as suspicious
	Amount        float64   // Amount involved in the transaction
	Timestamp     time.Time // Timestamp of the transaction
	IsResolved    bool      // Whether the suspicious activity has been resolved
}

type ExportLog struct {
	LogID      string    // Unique identifier for the export log
	ExportedBy string    // User or system that performed the export
	FileName   string    // Name of the exported file
	ExportType string    // Type of export (e.g., CSV, JSON)
	ExportedAt time.Time // Timestamp of export
	FileSize   int64     // Size of the exported file in bytes
	Status     string    // Status of the export (e.g., completed, failed)
}

type ImportedLog struct {
	LogID           string    // Unique identifier for the import log
	ImportedBy      string    // User or system that performed the import
	FileName        string    // Name of the imported file
	ImportType      string    // Type of import (e.g., CSV, JSON)
	ImportedAt      time.Time // Timestamp of import
	Status          string    // Status of the import (e.g., completed, failed)
	RecordsImported int       // Number of records successfully imported
}

type CacheMonitor struct {
	CacheID     string    // Unique identifier for the cache
	Name        string    // Name of the cache
	Usage       float64   // Current usage percentage
	Threshold   float64   // Usage threshold percentage
	LastChecked time.Time // Timestamp of the last check
	IsHealthy   bool      // Whether the cache is functioning within healthy parameters
}

type CacheUsage struct {
	CacheID      string             // Unique identifier for the cache
	UsageRecords []CacheUsageRecord // List of usage records
	LastUpdated  time.Time          // Timestamp of the last usage update
}

type CacheUsageRecord struct {
	Timestamp time.Time // Timestamp of the record
	Usage     float64   // Cache usage percentage at the time
}

type EncryptionStandard struct {
	StandardID  string    // Unique identifier for the encryption standard
	Name        string    // Name of the encryption standard
	Description string    // Description of the standard
	Version     string    // Version of the standard
	Strength    int       // Strength of the encryption (e.g., key size in bits)
	IsApproved  bool      // Whether the standard is approved for use
	LastUpdated time.Time // Timestamp of the last update
}

type Role struct {
	RoleID        string    // Unique identifier for the role
	Name          string    // Name of the role
	Permissions   []string  // List of permissions associated with the role
	AssignedUsers []string  // List of user IDs assigned to this role
	CreatedAt     time.Time // Timestamp of role creation
	UpdatedAt     time.Time // Timestamp of the last update
	IsActive      bool      // Whether the role is currently active
}

type AuditTask struct {
	TaskID   string
	EntityID string
	Interval time.Duration
	NextRun  time.Time
	Active   bool
}

type TransactionReversion struct {
	TransactionID string
	Reverted      bool
	Timestamp     time.Time
	Reason        string
}

type AdminNotification struct {
	NotificationID string
	Message        string
	Timestamp      time.Time
	Read           bool
}

// AuditAction represents a single action in an audit trail.
type AuditAction struct {
	ActionID  string    // Unique ID for the action
	Details   string    // Description of the action
	Timestamp time.Time // When the action was recorded
}

// IDDocument holds the ID verification document for a compliance check.
type IDDocument struct {
	DocumentID string    // Unique identifier for the ID document
	Issuer     string    // Issuer of the ID document
	ExpiryDate time.Time // Expiration date of the document
	Hash       string    // Hash of the document contents
}

// AccessControls defines the access rules for an entity.
type AccessControls struct {
	EntityID          string   // Entity the controls apply to
	AuthorizedRoles   []string // Roles that have access
	RestrictedActions []string // Actions restricted under the control rules
}

// ViolationRecord represents a record of compliance violations for an entity.
type ViolationRecord struct {
	EntityID         string    // Unique identifier of the entity
	ViolationDetails string    // Description of the compliance violation
	ReportedAt       time.Time // Timestamp when the violation was reported
	Resolved         bool      // Indicates if the violation has been resolved
}

// ComplianceSummary provides a summary of an entity's compliance history.
type ComplianceSummary struct {
	EntityID           string              // Unique identifier of the entity
	Violations         []ViolationRecord   // List of compliance violations
	Restrictions       []RestrictionRecord // List of restrictions applied
	LastAuditTimestamp time.Time           // Last audit time
	ComplianceScore    float64             // Overall compliance score
}

// ComplianceCertificate represents a certificate granted to an entity for compliance.
type ComplianceCertificate struct {
	CertID          string    // Unique identifier for the certificate
	EntityID        string    // Unique identifier of the entity
	IssuedAt        time.Time // Timestamp when the certificate was issued
	ExpiryDate      time.Time // Expiration date of the certificate
	ComplianceLevel string    // Level of compliance achieved
	IsValid         bool      // Indicates if the certificate is still valid
}

// RegulatoryRequest represents a request for compliance approval submitted to regulators.
type RegulatoryRequest struct {
	RequestID   string    // Unique identifier for the request
	EntityID    string    // Entity making the request
	Details     string    // Description of the request
	SubmittedAt time.Time // Timestamp of submission
	Status      string    // Current status of the request (e.g., pending, approved)
}

// RegulatoryFramework represents a set of regulatory standards or framework.
type RegulatoryFramework struct {
	Regulations   map[string]string // Key-value pair of regulation names and details
	EffectiveDate time.Time         // Date from which the new framework is effective
}

// ComplianceActionLog records actions taken for compliance-related purposes.
type ComplianceActionLog struct {
	EntityID      string    // Entity ID related to the compliance action
	ActionDetails string    // Details of the compliance action
	Timestamp     time.Time // Time of the action
}

// AccessRestriction represents access control details on compliance information.
type AccessRestriction struct {
	EntityID           string    // Entity the restriction applies to
	RestrictionDetails string    // Details of the restriction
	RestrictedAt       time.Time // Timestamp of restriction
}

// ComplianceRecord stores the compliance check data for a specific action or transaction
type ComplianceRecord struct {
	ActionID      string           // Unique identifier for the action or transaction
	Status        ComplianceStatus // Status of the compliance check
	CheckedBy     string           // Compliance officer or module responsible for the check
	EncryptedData string           // Field to hold encrypted data
}

// ComplianceAddition represents the compliance system managing the checks
type ComplianceAddition struct {
	ComplianceRules []string   // List of predefined rules to check
	LedgerInstance  *Ledger    // Reference to the ledger for storing compliance records
	mutex           sync.Mutex // Mutex for thread-safe operations
}

// ComplianceExecutionRecord holds information about compliance checks and their status
type ComplianceExecutionRecord struct {
	ID              string   // Unique identifier for the compliance record
	Status          string   // Current status of the compliance (e.g., "Pending", "Approved", "Rejected")
	TransactionID   string   // ID of the associated transaction
	ValidatorID     string   // ID of the validator or authority node that performed the compliance check
	Timestamp       int64    // Unix timestamp when the compliance check was performed
	ComplianceRules []string // List of compliance rules that were checked
	Comments        string   // Any additional comments or notes regarding the compliance check
	ResultHash      string   // Hash of the result for immutability verification
	ExpiryDate      int64    // Unix timestamp for expiry of compliance, if applicable
	IsFinalized     bool     // Whether the compliance check is finalized
	EncryptedData   string   // Field for storing encrypted compliance data as a string

}

// ComplianceHistory represents an entry in the compliance history of an entity.
type ComplianceHistory struct {
	EntityID    string    // ID of the entity being reviewed
	Status      string    // Compliance status at the time of review
	Description string    // Description of compliance event
	Timestamp   time.Time // Timestamp of the compliance event
}

// RegulatoryNotice represents a notice issued to an entity for regulatory issues.
type RegulatoryNotice struct {
	EntityID string    // ID of the entity receiving the notice
	Notice   string    // Content of the notice
	IssuedAt time.Time // Timestamp when the notice was issued
}

// LegalDocument represents a legal document required for compliance purposes.
type LegalDocument struct {
	DocumentID string    // Unique identifier for the document
	Content    string    // Content of the legal document
	CreatedAt  time.Time // Date when the document was created
	ExpiresAt  time.Time // Expiration date for the document
}

// RestrictionRule defines a rule for restricting certain actions based on compliance.
type RestrictionRule struct {
	RuleID      string // Unique identifier for the restriction rule
	Description string // Description of the restriction rule
	Conditions  string // Conditions under which the restriction applies
	AppliesTo   string // Entities affected by this rule
}

// RiskProfile represents the calculated risk profile of an entity.
type RiskProfile struct {
	EntityID     string    // Entity ID being assessed
	RiskLevel    string    // Risk level (e.g., Low, Medium, High)
	LastAssessed time.Time // Timestamp of the last assessment
	Comments     string    // Additional comments or insights
}

type AuditEntry struct {
	EntryID     string    // Unique identifier for the audit entry
	Signature   []byte    // Digital signature of the entry
	ContentHash string    // Hash of the content to validate integrity
	Validated   bool      // Whether the entry has been validated
	Timestamp   time.Time // Time of the audit entry creation
}

// ComplianceData contains specific compliance-related data for an entity.
type ComplianceData struct {
	EntityID        string    // ID of the entity
	ComplianceScore int       // Compliance score
	Violations      []string  // List of compliance violations
	LastReviewed    time.Time // Last review timestamp
}

// RegulatoryFeedback captures feedback from regulatory bodies.
type RegulatoryFeedback struct {
	FeedbackID string    // Unique identifier for the feedback
	EntityID   string    // Entity that the feedback applies to
	Comments   string    // Feedback comments
	ReviewedAt time.Time // Timestamp of when the feedback was reviewed
}

// EncryptionPolicy represents an encryption standard for compliance.
type EncryptionPolicy struct {
	EntityID   string    // ID of the entity associated with this policy
	Algorithm  string    // Algorithm used for encryption (e.g., AES)
	KeyLength  int       // Length of encryption key in bits
	ValidUntil time.Time // Expiration date of the encryption standard
}

// SecurityProfile represents a security profile with clearance levels.
type SecurityProfile struct {
	EntityID       string // ID of the entity
	ClearanceLevel int    // Security clearance level
	Roles          []Role // Roles assigned to the entity
	Active         bool   // Profile active status
}

// NodeActivityLog represents activity logs for a specific node.
type NodeActivityLog struct {
	NodeID    string    // Unique identifier of the node
	Timestamp time.Time // Time of the logged activity
	Action    string    // Action or event description
	Status    string    // Status of the activity (e.g., success, error)
	Details   string    // Additional details about the activity
}

// ComplianceIssue represents a compliance issue that may need escalation.
type ComplianceIssue struct {
	IssueID     string    // Unique identifier for the compliance issue
	EntityID    string    // ID of the entity associated with the issue
	Level       string    // Current escalation level
	Description string    // Description of the compliance issue
	Timestamp   time.Time // Time when the issue was reported
}

// AccessRequest represents a user access request for compliance review.
type AccessRequest struct {
	RequestID   string    // Unique ID of the access request
	UserID      string    // User ID associated with the request
	EntityID    string    // ID of the entity requesting access
	Status      string    // Approval status (e.g., pending, approved, denied)
	RequestedAt time.Time // Time the request was made
}

// UserPrivacySettings represents privacy settings for user compliance.
type UserPrivacySettings struct {
	UserID        string    // Unique ID of the user
	GDPRCompliant bool      // GDPR compliance flag
	CCPACompliant bool      // CCPA compliance flag
	LastReviewed  time.Time // Last time settings were reviewed
}

// RegulatoryResponse stores responses from regulatory bodies.
type RegulatoryResponse struct {
	ResponseID   string    // Unique ID for the regulatory response
	EntityID     string    // ID of the entity the response pertains to
	Feedback     string    // Feedback from the regulator
	DateReceived time.Time // Date when the response was received
}

// NodeComplianceMetric stores compliance metrics for a node.
type NodeComplianceMetric struct {
	NodeID          string    // Unique identifier for the node
	ComplianceScore int       // Score indicating compliance level
	LastChecked     time.Time // Last checked date for compliance
}

// RegulatoryAdjustments represents adjustments to compliance processes.
type RegulatoryAdjustments struct {
	AdjustmentID    string    // Unique identifier for adjustments
	Description     string    // Description of the adjustments made
	DateImplemented time.Time // When the adjustment was implemented
}

// License represents a license for an entity to operate.
type License struct {
	EntityID    string    // ID of the entity holding the license
	LicenseType string    // Type of license granted
	IssuedAt    time.Time // License issuance date
	ExpiryDate  time.Time // License expiration date
}

// ComplianceAlert represents an alert for potential compliance violations.
type ComplianceAlert struct {
	AlertID      string    // Unique identifier for the alert
	EntityID     string    // ID of the entity the alert pertains to
	AlertDetails string    // Details of the compliance issue
	DateCreated  time.Time // Alert creation timestamp
	Status       string    // Status of the alert (e.g., Open, Resolved)
}

// EncryptionStandards defines encryption parameters for compliance data.
type EncryptionStandards struct {
	StandardID string    // ID for the encryption standard
	Algorithm  string    // Encryption algorithm to be used
	KeyLength  int       // Key length in bits
	ValidUntil time.Time // Expiry of the encryption standard
}

type AuditIssue struct {
	IssueID     string
	Description string
	Resolved    bool
	ResolvedAt  time.Time
	Resolution  string
}

type AuditSummary struct {
	TotalIssues        int
	ResolvedIssues     int
	PendingIssues      int
	LastAuditTimestamp time.Time
}

type SuspiciousActivityReport struct {
	ReportID      string
	EntityID      string
	Description   string
	Timestamp     time.Time
	FlaggedIssues []string
}

type AuditRule struct {
	RuleID   string
	Criteria string
	Action   string
	Active   bool
}

type ContractDeploymentAudit struct {
	ContractID    string
	DeployedAt    time.Time
	Compliant     bool
	ComplianceLog string
}

type SystemAlert struct {
	AlertID     string
	Description string
	Timestamp   time.Time
	Resolved    bool
}

type AuditLog struct {
	EntryID   string
	Timestamp time.Time
	Action    string
	UserID    string
	Content   string
}

type AuditRecord struct {
	EntryID   string
	Timestamp time.Time
	Validated bool
	EntityID  string
	Action    string
	Details   map[string]string
}

type ExportOptions struct {
	Format        string
	IncludeAll    bool
	EncryptionKey string
}

type ImportOptions struct {
	Source        string
	Validation    bool
	EncryptionKey string
}

type ComplianceStatus struct {
	EntityID      string
	IsCompliant   bool
	LastCheckTime time.Time
	NextCheckTime time.Time
}

// RestrictionRecord holds information about restrictions within the ledger
type RestrictionRecord struct {
	ID               string   // Unique identifier for the restriction
	Type             string   // Type of restriction (e.g., "Transaction", "Account", "Address", "Node")
	Status           string   // Current status (e.g., "Active", "Inactive", "Pending")
	Metadata         string   // Additional metadata about the restriction
	Reason           string   // Reason for applying the restriction
	EffectiveDate    int64    // Unix timestamp when the restriction becomes effective
	ExpiryDate       int64    // Unix timestamp when the restriction expires, if applicable
	AppliedBy        string   // ID or name of the entity or node that applied the restriction
	AffectedEntity   string   // ID or address of the entity, transaction, or account affected by the restriction
	SeverityLevel    int      // Level of severity (e.g., 1 for low, 5 for high)
	RemediationSteps []string // Steps required to resolve or lift the restriction
	CreationTime     int64    // Unix timestamp when the restriction was created
	LastModified     int64    // Unix timestamp of the last modification
	EncryptedData    string   // Field to store the encrypted restriction result as a string

}

// ComplianceContract represents a smart contract for automated compliance enforcement
type ComplianceContract struct {
	ContractID      string     // Unique identifier for the contract
	Creator         string     // Creator of the contract (e.g., regulatory body or compliance authority)
	ComplianceRules []string   // List of compliance rules enforced by the contract
	LedgerInstance  *Ledger    // Reference to the ledger for recording actions and compliance results
	mutex           sync.Mutex // Mutex for thread-safe operations
}

// ComplianceResult stores the result of the compliance check executed by the contract
type ComplianceResult struct {
	ActionID  string    // Unique identifier for the action
	IsValid   bool      // Whether the action complies with rules
	Reason    string    // Reason for failure (if applicable)
	Timestamp time.Time // Timestamp of the compliance check
}

// ComplianceExecution represents a compliance execution process for an action
type ComplianceExecution struct {
	ExecutionID    string     // Unique identifier for the compliance execution
	ActionID       string     // ID of the action being validated for compliance
	Executor       string     // Address of the entity executing compliance (e.g., validator)
	RulesApplied   []string   // List of compliance rules applied
	Timestamp      time.Time  // Time when the compliance execution was initiated
	LedgerInstance *Ledger    // Reference to the ledger for recording results
	mutex          sync.Mutex // Mutex for thread-safe operations
}

// ComplianceExecutionResult holds the result of a compliance execution
type ComplianceExecutionResult struct {
	ExecutionID string    // ID of the execution process
	ActionID    string    // The action being validated
	IsValid     bool      // Whether the action complies with the rules
	Reason      string    // Reason for failure, if applicable
	Timestamp   time.Time // Timestamp of the result
}

// ComplianceRestrictions defines a set of rules and restrictions
type ComplianceRestrictions struct {
	RestrictionID    string     // Unique identifier for the restriction
	RestrictionRules []string   // Rules for the compliance restrictions
	CreatedAt        time.Time  // Timestamp of when the restriction was created
	EnforcedBy       string     // Address of the enforcer (e.g., validator)
	LedgerInstance   *Ledger    // Reference to the ledger for recording restrictions
	mutex            sync.Mutex // Mutex for thread-safe operations
}

// RestrictionResult defines the result of a restriction check
type RestrictionResult struct {
	RestrictionID string    // ID of the restriction applied
	ActionID      string    // ID of the action that was restricted
	IsRestricted  bool      // Whether the action is restricted
	Reason        string    // Reason for restriction
	Timestamp     time.Time // Timestamp of the restriction check
}

// DataProtectionPolicy defines policies to protect personal and sensitive data
type DataProtectionPolicy struct {
	PolicyID         string     // Unique ID for the data protection policy
	EncryptionMethod string     // Type of encryption method (e.g., AES, RSA)
	CreatedAt        time.Time  // Timestamp when the policy was created
	EnforcedBy       string     // Address of the enforcer (e.g., admin/validator)
	LedgerInstance   *Ledger    // Reference to the ledger for recording policies
	mutex            sync.Mutex // Mutex for thread-safe operations
}

// DataProtectionRecord logs information about data protection measures taken
type DataProtectionRecord struct {
	PolicyID      string    // ID of the data protection policy applied
	DataHash      string    // Hash of the protected data
	IsEncrypted   bool      // Whether the data is encrypted
	Timestamp     time.Time // Time when the protection was applied
	EncryptedData string    // Field to store the encrypted data as a string
}

// KYCRecord stores the details of a user's KYC verification
type KYCRecord struct {
	UserID       string    // Unique identifier of the user
	Status       KYCStatus // Status of the KYC verification
	VerifiedAt   time.Time // Timestamp of verification
	DataHash     string    // Hash of KYC data
	EncryptedKYC []byte
}

// KYCManager handles KYC verification and maintains records
type KYCManager struct {
	Records        map[string]KYCRecord // Stores KYC records by UserID
	LedgerInstance *Ledger              // Reference to the ledger for recording KYC actions
	mutex          sync.Mutex           // Mutex for thread-safe operations
}

// Restriction represents a restriction record imposed on accounts or transactions.
type Restriction struct {
	ID            string
	Target        string
	Reason        string
	EnforcedBy    string
	Timestamp     int64
	Expiration    int64 // Unix timestamp for expiration (if applicable)
	EncryptedData string
}

// KYCStatus defines the status of a KYC verification process.
type KYCStatus struct {
	IsVerified bool      // Indicates whether the KYC process was successfully verified
	Reason     string    // Reason for failure (if applicable), or description of the status
	VerifiedAt time.Time // Timestamp of when the KYC verification was completed
	VerifiedBy string    // The entity (e.g., compliance officer or system) that verified the KYC
}

// ComplianceEngine ensures that all blockchain activities adhere to regulatory standards.
type ComplianceEngine struct {
	ComplianceID       string                      // Unique identifier for the compliance engine instance
	Rules              []ComplianceRule            // List of compliance rules to be enforced
	ActiveMonitors     map[string]bool             // Monitors that are active on specific nodes, channels, or transactions
	ComplianceReports  []ComplianceReport          // Reports generated after compliance checks
	ViolationThreshold int                         // Threshold for violations before actions are taken
	ActionsTaken       map[string]ComplianceAction // Actions taken in response to non-compliance
	LoggingEnabled     bool                        // Flag indicating if logging of compliance checks is enabled
	mutex              sync.Mutex                  // Mutex to ensure thread-safe operations
}

// ComplianceRule represents a specific rule or standard that must be adhered to.
type ComplianceRule struct {
	RuleID      string    // Unique identifier for the compliance rule
	Description string    // Detailed description of the rule
	Severity    string    // Severity level if the rule is violated (Low, Medium, High)
	Enforcement string    // How the rule is enforced (e.g., automatic, manual review)
	CreatedAt   time.Time // Timestamp of when the rule was created
}

// ComplianceReport represents the outcome of a compliance check.
type ComplianceReport struct {
	ReportID      string    // Unique identifier for the compliance report
	EntityID      string    // Entity that the report applies to
	Content       string    // Detailed content of the compliance check
	CreatedAt     time.Time // Timestamp of report creation
	IntegrityHash string    // Hash for data integrity
}

// ComplianceAction represents an action taken due to non-compliance.
type ComplianceAction struct {
	ActionID    string    // Unique identifier for the action
	ActionType  string    // Type of action (e.g., "Alert", "Suspend", "Investigate")
	Description string    // Detailed description of the action
	Timestamp   time.Time // Timestamp of when the action was executed
}

// ************** Conditional Flag Structs **************

type ExecutionPathEntry struct {
	Path      string    // The path in the program flow
	Timestamp time.Time // Time when this path was tracked
}

// SystemErrorEntry represents a system error log entry
type SystemErrorEntry struct {
	ErrorID     string    // Unique identifier for the error
	Description string    // Description of the error
	Timestamp   time.Time // Time the error occurred
}

// ProgramStatusEntry represents a program status log entry
type ProgramStatusEntry struct {
	ProgramID string    // Unique identifier for the program
	Status    string    // Current status of the program
	ErrorCode int       // Error code if applicable
	Timestamp time.Time // Time the status was recorded
}

// ConditionLogEntry represents a condition log entry
type ConditionLogEntry struct {
	ConditionID string    // Unique identifier for the condition
	Status      string    // Status of the condition
	Timestamp   time.Time // Time the condition was logged
}

// ProgramLogEntry represents a log entry for program-specific operations.
type ProgramLogEntry struct {
	ProgramID      string                 // Unique identifier for the program
	Operation      string                 // Description of the operation
	Result         string                 // Result of the operation (e.g., "Success", "Failed")
	Timestamp      time.Time              // Time the operation occurred
	AdditionalData map[string]interface{} // Additional data related to the operation
}

// ************** Consensus Operations Structs **************

// PunitiveMeasureRecord represents a record of punitive actions taken against validators.
type PunitiveMeasureRecord struct {
	ActionID    string    // Unique identifier for the action
	ValidatorID string    // Validator subject to punitive action
	Reason      string    // Reason for the punitive action
	Timestamp   time.Time // Time the punitive action was taken
	Status      string    // Status of the action (e.g., active, reverted)
}

// PunishmentAdjustmentLog represents a log entry for adjustments made to punitive measures.
type PunishmentAdjustmentLog struct {
	AdjustmentID string    // Unique identifier for the adjustment
	ActionID     string    // Related punitive action ID
	AdjustedBy   string    // Who made the adjustment
	Timestamp    time.Time // Time the adjustment was made
	Details      string    // Details of the adjustment
}

// DifficultyAdjustmentLog records adjustments made to consensus difficulty.
type DifficultyAdjustmentLog struct {
	AdjustmentID       string
	Timestamp          time.Time
	NewDifficultyLevel int
	Reason             string
}

// BlockGenerationLog records block generation times.
type BlockGenerationLog struct {
	BlockID        string
	GenerationTime time.Duration
	Timestamp      time.Time
}

// ConsensusAuditLog tracks consensus audit details.
type ConsensusAuditLog struct {
	AuditID             string
	ValidatorID         string
	Timestamp           time.Time
	ParticipationStatus string
}

// FinalityCheckLog records finality checks for blocks.
type FinalityCheckLog struct {
	BlockID        string
	FinalityStatus bool
	Timestamp      time.Time
}

// ConditionManager manages system conditions, triggers, and dependencies for policies and events.
type ConditionManager struct {
	ConditionID     string                 // Unique identifier for the condition.
	ConditionName   string                 // Human-readable name of the condition.
	Description     string                 // Description of the condition and its purpose.
	Expression      string                 // Logical expression defining the condition (e.g., "ResourceUsage > 80").
	TargetResources []string               // List of resources or entities the condition applies to.
	TriggerActions  []string               // Actions triggered when the condition is met.
	CreatedAt       time.Time              // Timestamp when the condition was created.
	CreatedBy       string                 // Identifier of the entity that created the condition.
	IsActive        bool                   // Indicates whether the condition is currently active.
	Metadata        map[string]interface{} // Additional metadata about the condition.
}

// ConditionEvaluationLog keeps a history of condition evaluations and their outcomes.
type ConditionEvaluationLog struct {
	ConditionID      string    // Identifier of the evaluated condition.
	Timestamp        time.Time // Timestamp of the evaluation.
	Result           bool      // Result of the evaluation (true if the condition was met, false otherwise).
	TriggeredActions []string  // List of actions triggered by the condition.
	LogDetails       string    // Additional details about the evaluation process.
}

// PoHLog represents a log entry for Proof of History (PoH) validations.
type PoHLog struct {
	ValidatorID string
	Status      string
	Timestamp   time.Time
}

// StakeChangeRecord represents a record for encrypted stake change for a validator
type StakeChangeRecord struct {
	ValidatorID          string
	EncryptedStakeChange []byte
	Timestamp            time.Time
}

// PunishmentRecord represents a record of a punishment issued to a validator.
type PunishmentRecord struct {
	ValidatorID     string
	Reason          string
	Timestamp       time.Time
	PunishmentLevel int // Level of punishment severity
}

// RewardRecord represents a record of a reward issued to a validator.
type RewardRecord struct {
	ValidatorID     string
	EncryptedReward string // Encrypted reward amount
	Timestamp       time.Time
}

// ValidatorActivityLog tracks validator activity.
type ValidatorActivityLog struct {
	ValidatorID string
	Action      string
	Timestamp   time.Time
	Details     string
}

// ValidatorBanRecord represents a record of a banned validator, including the reason and timestamp.
type ValidatorBanRecord struct {
	ValidatorID string
	Reason      string
	Timestamp   time.Time
}

// RewardDistributionMode and ValidatorSelectionMode settings
type RewardDistributionMode struct {
	ModeID      string
	Description string
	Active      bool
}

type ValidatorSelectionMode struct {
	ModeID      string
	Description string
	Active      bool
}

// StakeChange represents a change in a validator's stake.
type StakeChange struct {
	ValidatorID  string
	ChangeAmount float64 // Encrypted stake change amount
	Timestamp    time.Time
}

// StakeLog represents a log entry for stake adjustments.
type StakeLog struct {
	ValidatorID string
	Adjustment  []byte // Store encrypted adjustment as []byte
	Timestamp   time.Time
}

// ValidatorPenalty represents a penalty record for a validator.
type ValidatorPenalty struct {
	ValidatorID   string
	PenaltyAmount []byte // Store encrypted penalty amount as []byte
	Timestamp     time.Time
}

// EpochLog represents a historical record of epoch changes for auditing.
type EpochLog struct {
	EpochID           string
	Timestamp         time.Time
	Duration          time.Duration // Add Duration for original unencrypted value
	EncryptedDuration []byte        // Store the encrypted duration here
}

// ReinforcementPolicy defines the rules for reinforcing consensus.
type ReinforcementPolicy struct {
	PolicyID    string
	Description string
	Details     map[string]interface{} // Specifics of the policy
	Timestamp   time.Time
}

// HealthLog represents health metrics of the consensus for analysis.
type HealthLog struct {
	HealthID       string
	Metric         string
	Value          float64
	EncryptedValue []byte // Stores encrypted health metric value
	Timestamp      time.Time
}

// ************** Cryptography Structs **************

// SignatureAggregation represents the aggregation of multiple cryptographic signatures.
type SignatureAggregation struct {
	AggregationID       string            // Unique identifier for the aggregation process
	Signatures          map[string][]byte // Map of participant IDs to their respective signatures
	AggregatedSignature []byte            // The final aggregated signature
	Participants        []string          // List of participant IDs involved in the aggregation
	Threshold           int               // Minimum number of signatures required for aggregation
	Algorithm           string            // Cryptographic algorithm used (e.g., BLS, ECDSA)
	Timestamp           time.Time         // Timestamp when the aggregation was created
	IsVerified          bool              // Whether the aggregated signature has been verified
	VerificationLogs    []VerificationLog // Logs of verification attempts
	Status              string            // Status of the aggregation (e.g., pending, completed, failed)
}

// VerificationLog represents a log entry for signature verification attempts.
type VerificationLog struct {
	LogID        string    // Unique identifier for the log entry
	VerifierID   string    // ID of the entity performing the verification
	VerifiedAt   time.Time // Timestamp of the verification attempt
	IsSuccessful bool      // Whether the verification was successful
	Error        string    // Error message if the verification failed
}

// ************** Dao Structs **************

// DAORecord holds all DAO-related information in the ledger.
type DAORecord struct {
	ID               string
	Members          map[string]DAOMember
	Proposals        map[string]DAOProposal
	Transactions     map[string]TransactionRecord
	GovernanceStakes map[string]float64
	RoleAssignments  map[string]string // Maps members to roles
}

// DAO represents a decentralized autonomous organization on the blockchain.
type DAO struct {
	DAOID           string                // Unique ID of the DAO
	Name            string                // Name of the DAO
	CreatorWallet   string                // Wallet of the DAO creator
	CreatedAt       time.Time             // Time of DAO creation
	Members         map[string]*DAOMember // Members of the DAO with roles and permissions
	FundsVault      *DAOFundVault         // DAO's fund vault
	VotingThreshold int                   // Minimum number of votes required for DAO decisions
	IsActive        bool                  // Is DAO active or deactivated
}

// DAOMember represents a member of a DAO with their role and permissions.
type DAOMember struct {
	WalletAddress string // Wallet address of the member
	Role          string // Role of the member (Admin, Member, Treasurer)
	VotingPower   int    // Voting power of the member
	IsAuthorized  bool   // Whether the member is authorized to perform DAO actions
}

// DAOManagement handles DAO creation, updates, and management.
type DAOManagement struct {
	mutex             sync.Mutex      // Mutex for thread-safe operations
	DAOs              map[string]*DAO // Map of DAO objects by DAO ID
	Ledger            *Ledger         // Ledger reference for recording DAO activities
	EncryptionService *Encryption     // Encryption service for securing DAO data
	Syn900Verifier    *Syn900Verifier // Verifier for DAO-related actions
}

// AccessControl is responsible for managing roles and permissions within the DAO.
type AccessControl struct {
	mutex             sync.Mutex        // For thread-safe operations
	DAOID             string            // ID of the DAO
	Members           map[string]string // Mapping of wallet addresses to roles
	Ledger            *Ledger           // Ledger instance for storing role assignments
	EncryptionService *Encryption       // Encryption service for securing role information
}

// DAOProposal defines the structure for a DAO proposal.
type DAOProposal struct {
	ProposalID   string    // Unique proposal identifier
	Title        string    // Proposal title
	Description  string    // Proposal description
	Author       string    // Author's wallet address
	CreationTime time.Time // Proposal creation time
	VoteCount    int       // Total number of votes
	ApproveCount int       // Number of approvals
	RejectCount  int       // Number of rejections
	Status       string    // "Pending", "Approved", "Rejected"
}

// DAOFundVault manages the funds for a DAO.
type DAOFundVault struct {
	mutex             sync.Mutex                   // For thread-safe operations
	DAOID             string                       // ID of the DAO
	Balance           float64                      // Current balance of the DAO vault
	Ledger            *Ledger                      // Ledger instance for recording transactions
	EncryptionService *Encryption                  // Encryption service for securing fund management
	Syn900Verifier    *tokenledgers.Syn900Verifier // Verifier for emergency access via Syn900
	TransactionLimit  float64                      // Daily transaction limit to ensure security
	LastTransactionAt time.Time                    // Timestamp of the last transaction
	TransactionQueue  []VaultTransaction           // Queue of pending transactions
	Admins            map[string]bool              // DAO admin addresses with access to funds
}

// VaultTransaction represents a transaction from the DAO vault.
type VaultTransaction struct {
	TransactionID string
	Amount        float64
	Recipient     string
	Timestamp     time.Time
	ApprovedBy    []string // List of admin approvals
	Status        string   // Pending, Approved, Rejected
}

// EmergencyAccessRequest represents an emergency procedure triggered by the Syn900 protocol.
type EmergencyAccessRequest struct {
	RequestID       string
	RequestedBy     string
	Reason          string
	Timestamp       time.Time
	Status          string // Pending, Approved, Rejected
	ApprovalConfirm []string
}

// GovernanceStake represents a user's governance staking record.
type GovernanceStake struct {
	StakerWallet   string    // Wallet address of the staker
	Amount         float64   // Amount of tokens staked for governance
	VotingPower    float64   // Derived voting power based on staked amount
	StakeTimestamp time.Time // Time when the stake was made
	IsActive       bool      // Whether the stake is currently active
}

// GovernanceStakingSystem manages the staking system for governance in a DAO.
type GovernanceStakingSystem struct {
	DAOID             string                      // DAO ID associated with this staking system
	TotalStakedTokens float64                     // Total amount of tokens staked in the DAO for governance
	StakingRecords    map[string]*GovernanceStake // User staking records
	MinStakeAmount    float64                     // Minimum amount required to stake for governance
	StakingDuration   time.Duration               // Lock-in period for governance staking
}

// StakingManager handles governance staking within the DAO.
type StakingManager struct {
	mutex             sync.Mutex                          // Mutex for thread-safe operations
	Ledger            *Ledger                             // Ledger reference for recording staking actions
	EncryptionService *Encryption                         // Encryption for secure staking transactions
	Syn900Verifier    *tokenledgers.Syn900Verifier        // Identity verification system using Syn900
	GovernanceStakes  map[string]*GovernanceStakingSystem // DAO governance staking systems
}

// GovernanceTokenVotingSystem manages the governance token-based voting system.
type GovernanceTokenVotingSystem struct {
	mutex             sync.Mutex                     // Mutex for thread-safe operations
	Proposals         map[string]*GovernanceProposal // Map of governance proposals by proposal ID
	Ledger            *Ledger                        // Ledger to store all voting records
	EncryptionService *Encryption                    // Encryption service for secure votes
	Syn800Token       *tokenledgers.SYN300Token      // Token contract for voting
	Syn900Verifier    *tokenledgers.Syn900Verifier   // Verifier for identity checks via Syn900
}

// QuadraticProposal represents a proposal for quadratic voting.
type QuadraticProposal struct {
	ProposalID   string             // Unique proposal identifier
	ProposalText string             // Description of the proposal
	CreationTime time.Time          // Time when the proposal was created
	Deadline     time.Time          // Voting deadline for the proposal
	SubmittedBy  string             // Wallet address of the proposer
	TotalVotes   float64            // Total tokens squared (expressed as votes)
	YesVotes     float64            // Total quadratic tokens voted "Yes"
	NoVotes      float64            // Total quadratic tokens voted "No"
	VoterRecords map[string]float64 // Tracks how many tokens each user has voted
	Status       string             // "Open", "Passed", "Rejected"
}

// QuadraticVotingSystem manages the quadratic voting system.
type QuadraticVotingSystem struct {
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	Proposals         map[string]*QuadraticProposal // Map of quadratic proposals by proposal ID
	Ledger            *Ledger                       // Ledger to store all voting records
	EncryptionService *Encryption                   // Encryption service for secure votes
	Syn800Token       *tokenledgers.SYN300Token     // Token contract for voting
	Syn900Verifier    *tokenledgers.Syn900Verifier  // Verifier for identity checks via Syn900
}

// ************** Data Management Structs **************

// MetadataManagement handles metadata and transaction summaries.
type MetadataManagement struct {
	MetadataRecords map[string]string // Map of metadata keys to their values
	LastUpdated     time.Time         // Last updated timestamp
}

// DataTransferRecord represents a record of data transfer.
type DataTransferRecord struct {
	RecordID   string      // Unique identifier for the transfer record
	ChannelID  string      // ID of the data transfer channel
	DataBlocks []DataBlock // List of data blocks in the transfer
	CreatedAt  time.Time   // Timestamp when the transfer was initiated
}

// ************** Defi Structs **************

// InsurancePolicy represents a decentralized insurance policy
type InsurancePolicy struct {
	PolicyID       string
	InsuredEntity  string
	Premium        float64
	Terms          string
	CoverageAmount float64
	Status         string
	Duration       time.Duration
	StartTime      time.Time
	EndTime        time.Time
	Frozen         bool
	Locked         bool
	AutoRenew      bool
	ClaimFee       float64
}

type YieldFarmingRecord struct {
	RecordID      string    // Unique identifier for the record
	ParticipantID string    // ID of the participant
	AmountStaked  float64   // Amount staked in farming
	RewardsEarned float64   // Rewards earned from farming
	RecordedAt    time.Time // Timestamp of the record
}

type ParticipantPrediction struct {
    UserID  string
    Amount  float64
    Odds    float64
    Payout  float64
    EventID string
    Status  string
}

type StakingSnapshot struct {
    ProgramID    string
    TotalStaked  float64
    ParticipantData map[string]float64 // userID -> staked amount
    Timestamp    time.Time
}

type LoanAuditRecord struct {
    LoanID       string
    AuditDetails string
    Timestamp    time.Time
}

type LatePaymentRecord struct {
    LoanID     string
    Amount     float64
    DueDate    time.Time
    PaidDate   time.Time
    PenaltyFee float64
}


type VolatilityRecord struct {
    AssetID        string
    VolatilityRate float64
    Timestamp      time.Time
}

type MarketCapRecord struct {
    AssetID      string
    MarketCap    float64
    Timestamp    time.Time
}

type YieldFarmPool struct {
    PoolID           string
    TotalLiquidity   float64
    StakedTokens     map[string]float64 // UserID -> Amount Staked
    RewardBalance    float64
    APY              float64
    IsLocked         bool
    LastDistributed  time.Time
}

type YieldFarmEarning struct {
    UserID        string
    PoolID        string
    EarnedRewards float64
    LastHarvest   time.Time
}

type PoolPerformanceMetrics struct {
    PoolID          string
    TotalLiquidity  float64
    TotalRewards    float64
    TotalParticipants int
    APY             float64
    LastUpdated     time.Time
}


type SyntheticAssetPriceChange struct {
    AssetID    string
    OldPrice   float64
    NewPrice   float64
    ChangeTime time.Time
}


type StakingProgram struct {
    ProgramID      string
    RewardRate     float64
    MinStake       float64
    TotalStaked    float64
    Status         string
    LockedTokens   map[string]float64 // userID -> amount
}

type StakingParticipant struct {
    UserID       string
    StakedAmount float64
    Rewards      float64
    ProgramID    string
    Locked       bool
}


type LiquidityPoolTransaction struct {
    PoolID        string
    TransactionID string
    Action        string
    Amount1       float64
    Amount2       float64
    Timestamp     time.Time
}

type LPStaking struct {
    PoolID   string
    UserID   string
    Amount   float64
    Rewards  float64
    StakedAt time.Time
}

type PredictionEvent struct {
    EventID      string
    EventDetails string
    Odds         float64
    Status       string
    Outcome      string
    TotalPool    float64
    EscrowFunds  float64
}

type Prediction struct {
    EventID string
    UserID  string
    Amount  float64
    Odds    float64
    Payout  float64
    Status  string
}


// TokenRecord represents the details and history of a specific token.
type TokenRecord struct {
	TokenID           string             // Unique identifier for the token
	TokenName         string             // Name of the token
	Symbol            string             // Token symbol (e.g., ETH, BTC)
	TotalSupply       float64            // Total supply of the token
	CirculatingSupply float64            // Current circulating supply
	Decimals          int                // Number of decimal places for the token
	Owners            map[string]float64 // Map of owner addresses to their token balances
	Transactions      []TokenTransaction // List of transactions involving this token
	IsMintable        bool               // Indicates whether the token can be minted
	IsBurnable        bool               // Indicates whether the token can be burned
	CreatedAt         time.Time          // Timestamp when the token was created
	LastUpdated       time.Time          // Timestamp of the last update to the token record
	Status            string             // Status of the token (e.g., active, suspended)
}

// TokenTransaction represents a transaction involving a specific token.
type TokenTransaction struct {
	TransactionID string    // Unique identifier for the transaction
	Sender        string    // Address of the sender
	Receiver      string    // Address of the receiver
	Amount        float64   // Amount of tokens transferred
	Timestamp     time.Time // Timestamp of the transaction
	Type          string    // Type of transaction (e.g., mint, burn, transfer)
	Remarks       string    // Optional remarks or metadata for the transaction
}

type Bet struct {
	BetID      string
	Event      string // Encrypted
	Expiration time.Time
	Status     string
	Odds       float64 // Encrypted
	Amount     float64 // Total bet amount
	Winner     string  // Encrypted winner (if applicable)
}

type BetParticipant struct {
	BetID  string
	User   string // Encrypted
	Amount float64
}

type BetHistoryRecord struct {
	BetID      string
	User       string // Encrypted user
	Amount     float64
	TimePlaced time.Time
}

type BetConfig struct {
	BettingPaused bool
}

type CrowdfundingCampaign struct {
	CampaignID     string
	Title          string // Encrypted
	Description    string // Encrypted
	GoalAmount     float64
	CollectedFunds float64
	EndTime        time.Time
	Status         string // "Active", "Closed", "Failed"
	CreatorID      string // Encrypted creator ID
}

type CrowdfundingContribution struct {
	CampaignID string
	UserID     string // Encrypted
	Amount     float64
	Time       time.Time
}

type CrowdfundingAuditRecord struct {
	CampaignID string
	Details    string // Information about the audit
}

type ContributionRecord struct {
	CampaignID string
	UserID     string // Encrypted user ID
	Amount     float64
	Time       time.Time
}

type ContributionLimits struct {
	Min float64
	Max float64
}

// InsuranceClaim represents a claim made on an insurance policy
type InsuranceClaim struct {
	ClaimID       string    // Unique ID for the claim
	PolicyID      string    // ID of the insurance policy
	ClaimAmount   float64   // The amount being claimed
	ClaimDate     time.Time // The date the claim was made
	ClaimStatus   string    // Claim status ("Pending", "Approved", "Rejected")
	EncryptedData string    // Encrypted claim data for security
}

// InsuranceManager manages DeFi insurance policies and claims
type InsuranceManager struct {
	Policies          map[string]*InsurancePolicy // Active insurance policies
	Claims            map[string]*InsuranceClaim  // Claims made by policyholders
	Ledger            *Ledger                     // Ledger instance for tracking policies and claims
	EncryptionService *Encryption                 // Encryption service for secure data handling
	mu                sync.Mutex                  // Mutex for concurrent access to policies and claims
}

// DeFiManagement represents the core management of decentralized finance operations
type DeFiManagement struct {
	LiquidityPools      map[string]*LiquidityPool  // Managed liquidity pools
	AssetPools          map[string]*AssetPool      // Managed asset pools for synthetic assets or other DeFi assets
	YieldFarmingRecords map[string]*FarmingRecord  // Yield farming records
	LoanManagement      map[string]*Loan           // DeFi loan management
	SyntheticAssets     map[string]*SyntheticAsset // Synthetic assets issued in the network
	Ledger              *Ledger                    // Ledger instance for tracking all DeFi activities
	EncryptionService   *Encryption                // Encryption service for securing all data
	mu                  sync.Mutex                 // Mutex for managing concurrent access
}

// LiquidityPool represents a liquidity pool in DeFi operations
type LiquidityPool struct {
	PoolID             string    // Unique ID for the liquidity pool
	TotalLiquidity     float64   // Total liquidity in the pool
	AvailableLiquidity float64   // Available liquidity for operations
	RewardRate         float64   // Reward rate for liquidity providers
	CreatedAt          time.Time // Creation timestamp
	Status             string    // Status of the pool ("Active", "Paused", etc.)
	TotalStaked        float64   // Total amount staked in the pool (added field)
	TokenRatio        float64
    FeeRate           float64
    WithdrawalFee     float64
    IsSwapsPaused     bool
    TotalBalance      float64
    LastCompoundTime  time.Time
    RebalancingActive bool
}

// AssetPool represents an asset pool for synthetic or DeFi assets
type AssetPool struct {
	PoolID      string    // Unique ID for the asset pool
	TotalAssets float64   // Total assets in the pool
	AssetType   string    // Type of asset (e.g., synthetic, native)
	RewardRate  float64   // Reward rate for asset providers
	CreatedAt   time.Time // Creation timestamp
	Status      string    // Status of the asset pool ("Active", "Paused", etc.)
}

// FarmingRecord represents user participation in yield farming
type FarmingRecord struct {
	FarmingID      string
	UserID         string
	AmountStaked   float64
	RewardsEarned  float64 // Field to track the rewards earned
	StakeTimestamp time.Time
	Status         string
}

// OracleData represents data provided by a DeFi oracle
type OracleData struct {
	OracleID         string    // Unique ID for the oracle data
	DataFeedID       string    // ID of the data feed being provided by the oracle
	DataPayload      string    // The actual data being provided
	Verified         bool      // Whether the data has been verified
	Timestamp        time.Time // Timestamp of the data submission
	HandlerNode      string    // Node handling the oracle data submission
	EncryptedPayload string    // Encrypted version of the data payload
}

// OracleManager manages the lifecycle of DeFi oracles
type OracleManager struct {
	OracleSubmissions   map[string]*OracleData // Active oracle submissions
	VerifiedSubmissions []*OracleData          // Log of verified submissions
	PendingSubmissions  []*OracleData          // Queue of pending oracle submissions
	Ledger              *Ledger                // Ledger instance for logging oracle activities
	EncryptionService   *Encryption            // Encryption service for secure data handling
	mu                  sync.Mutex             // Mutex for concurrent operations
}

// Loan represents a loan given by a lender to a borrower
type Loan struct {
	LoanID        string        // Unique loan identifier
	Lender        string        // Lender's wallet address
	Borrower      string        // Borrower's wallet address
	Amount        float64       // Loan amount
	Collateral    float64       // Collateral deposited by the borrower
	InterestRate  float64       // Interest rate applied to the loan
	Duration      time.Duration // Loan duration
	StartDate     time.Time     // When the loan started
	ExpiryDate    time.Time     // Loan expiry date
	Status        string        // Loan status ("Active", "Repaid", "Defaulted")
	EncryptedData string        // Encrypted loan data for security
}

// LendingPool represents a pool of assets available for lending
type LendingPool struct {
	PoolID         string  // Unique identifier for the lending pool
	TotalLiquidity float64 // Total liquidity in the pool
	InterestRate   float64 // Interest rate offered by the pool
	AvailableFunds float64 // Available funds for lending
	ActiveLoans    []*Loan // List of active loans
	EncryptedData  string  // Encrypted pool data for security
}

// LendingManager manages decentralized lending and borrowing
type LendingManager struct {
	LendingPools      map[string]*LendingPool // Available lending pools
	Loans             map[string]*Loan        // All active loans
	Ledger            *Ledger                 // Ledger instance for logging lending and borrowing activities
	EncryptionService *Encryption             // Encryption service for secure data handling
	mu                sync.Mutex              // Mutex for managing concurrent access
}

// SyntheticAsset represents a synthetic asset in the system
type SyntheticAsset struct {
	AssetID         string    // Unique identifier for the synthetic asset
	AssetName       string    // Name of the synthetic asset (e.g., sUSD, sBTC)
	UnderlyingAsset string    // Underlying asset that the synthetic asset represents (e.g., USD, BTC)
	Price           float64   // Current price of the synthetic asset
	CollateralRatio float64   // Collateral ratio required to mint the synthetic asset
	TotalSupply     float64   // Total supply of the synthetic asset
	CreatedAt       time.Time // Timestamp when the asset was created
	Status          string    // Status of the synthetic asset ("Active", "Paused", "Deprecated")
	EncryptedData   string    // Encrypted data for the synthetic asset
}

// SyntheticAssetManager manages the creation and trading of synthetic assets
type SyntheticAssetManager struct {
	Assets            map[string]*SyntheticAsset // Map of all synthetic assets
	Ledger            *Ledger                    // Ledger instance for logging synthetic asset actions
	EncryptionService *Encryption                // Encryption service for secure data handling
	mu                sync.Mutex                 // Mutex for managing concurrent access
}

type StakeRecord struct {
	PoolID       string
	UserID       string
	AmountStaked float64
	Timestamp    time.Time
}

// FarmingPool represents a liquidity pool used for yield farming
type FarmingPool struct {
	PoolID         string    // Unique identifier for the farming pool
	TokenPair      string    // The token pair used for liquidity (e.g., ETH/USDC)
	TotalLiquidity float64   // Total liquidity in the pool
	RewardRate     float64   // Reward rate for liquidity providers
	Rewards        float64   // Total rewards available for distribution
	CreatedAt      time.Time // Timestamp when the pool was created
	Status         string    // Pool status ("Active", "Inactive")
	EncryptedData  string    // Encrypted pool data for privacy and security
}

// StakingRecord represents the details of a user's liquidity stake in the farming pool
type StakingRecord struct {
	StakeID        string    // Unique identifier for the stake
	StakerAddress  string    // The address of the staker
	AmountStaked   float64   // The amount of liquidity provided
	StakeTimestamp time.Time // Timestamp when the stake was made
	RewardEarned   float64   // Rewards earned so far
	EncryptedData  string    // Encrypted data for the stake
}

// YieldFarmingManager manages the yield farming pools and staked liquidity
type YieldFarmingManager struct {
	FarmingPools      map[string]*FarmingPool   // Active farming pools
	StakingRecords    map[string]*StakingRecord // Active staking records
	Ledger            *Ledger                   // Ledger instance for tracking farming activities
	EncryptionService *Encryption               // Encryption service for securing data
	mu                sync.Mutex                // Mutex for managing concurrent access
}

// DeFiRecord keeps track of all DeFi-related activities in the ledger
type DeFiRecord struct {
	InsurancePolicies map[string]InsurancePolicy // Insurance policies issued in the network
	LiquidityPools    map[string]LiquidityPool   // Liquidity pools for DeFi operations
	AssetPools        map[string]AssetPool       // Asset pools (e.g., for synthetic assets)
	YieldFarmingPools map[string]FarmingPool     // Yield farming pools
	OracleSubmissions map[string]OracleData      // Oracle data submissions
	Loans             map[string]Loan            // Active loans in the DeFi system
	SyntheticAssets   map[string]SyntheticAsset  // Synthetic assets in the system
}

// OracleSubmission represents a submission of data from a decentralized oracle.
type OracleSubmission struct {
	ID        string                 // Unique ID for the submission
	OracleID  string                 // ID of the oracle submitting the data
	Data      map[string]interface{} // Data submitted by the oracle
	Submitted time.Time              // Timestamp of the submission
	Verified  bool                   // Whether the data has been verified
}

// ************** Environmnet And System Core Structs **************

type EventAction struct {
	EventID   string
	Action    string
	Timestamp time.Time
}

type EventCondition struct {
	EventID string
	IsValid func() bool            // A function pointer to validate the condition
	Details map[string]interface{} // Additional condition details
}

type EventProperties struct {
	EventID    string
	Properties map[string]interface{}
}

type RecurringEvent struct {
	EventID  string
	Interval time.Duration
	LastRun  time.Time
}

type EventDependency struct {
	EventID      string
	Dependencies []string
}

type AutomationSchedule struct {
	Events []RecurringEvent
}

type EventTriggerCondition struct {
	EventID       string
	Conditions    map[string]interface{}
	LastValidated time.Time
}

type ChainActivityLog struct {
	Timestamp    time.Time
	ActivityType string
	Description  string
}

type BlockFinality struct {
	BlockID   string
	IsFinal   bool
	CheckedAt time.Time
}

type NodeRole struct {
	NodeID     string
	Role       string
	AssignedAt time.Time
}

type BlockData struct {
	BlockID      string
	Transactions []string
	Miner        string
	Timestamp    time.Time
}

type NetworkMetrics struct {
	NodeCount       int
	TransactionRate float64
	BlockLatency    time.Duration
}

type NodeHealthLog struct {
	NodeID      string
	HealthScore int
	Timestamp   time.Time
}

type NodeConfig struct {
	NodeID    string
	Config    map[string]interface{}
	UpdatedAt time.Time
}

type EmergencyStatus struct {
	Active    bool
	UpdatedAt time.Time
}

type EnvironmentStatus struct {
	Healthy   bool
	CheckedAt time.Time
}

type NodeHealthScore struct {
	NodeID      string
	HealthScore int
	CheckedAt   time.Time
}

type MiningDifficulty struct {
	DifficultyLevel int
	UpdatedAt       time.Time
}

type ContextResources struct {
	ContextID    string
	ResourceType string
	Amount       int
	ReservedAt   time.Time
}

type ContextClock struct {
	ContextID string
	ClockTime time.Time
	UpdatedAt time.Time
}

type ContextMemory struct {
	ContextID   string
	MemoryLimit int
	UpdatedAt   time.Time
}

type ExecutionCapacity struct {
	ContextID string
	Capacity  int
	UpdatedAt time.Time
}

type ContextVariables struct {
	ContextID string
	Variables map[string]interface{}
	UpdatedAt time.Time
}

type ContextDiagnosticsLog struct {
	ContextID   string
	Timestamp   time.Time
	Diagnostics string
}
type ContextReport struct {
	ContextID      string
	ResourceUsage  map[string]int
	Concurrency    int
	CheckpointTime time.Time
	Variables      map[string]interface{}
	Diagnostics    []string
}

type ResourcePolicy struct {
	ContextID     string
	PolicyDetails map[string]interface{}
	LastChecked   time.Time
}

type ContextConcurrency struct {
	ContextID        string
	ConcurrencyLevel int
	UpdatedAt        time.Time
}

type ContextCheckpoint struct {
	ContextID string
	State     map[string]interface{}
	SavedAt   time.Time
}

type ContextCleanup struct {
	ContextID string
	MarkedAt  time.Time
}

type ContextTerminationLog struct {
	ContextID string
	Timestamp time.Time
}

type RestartPolicy struct {
	ContextID string
	Policy    map[string]interface{}
	UpdatedAt time.Time
}

type TaskDelegation struct {
	TaskID        string
	FromContextID string
	ToContextID   string
	DelegatedAt   time.Time
	RetractedAt   *time.Time
}

type ContextObserver struct {
	ContextID  string
	ObserverID string
	AddedAt    time.Time
}

type ContextDependency struct {
	ContextID          string
	DependentContextID string
	AddedAt            time.Time
}

type ContextLifespan struct {
	ContextID string
	Duration  time.Duration
	UpdatedAt time.Time
}

type AccessRestrictions struct {
	ContextID    string
	Restrictions map[string]interface{}
	UpdatedAt    time.Time
}

type AccessRights struct {
	ContextID string
	Rights    map[string]interface{}
	UpdatedAt time.Time
}

type EventTrigger struct {
	TriggerID  string
	Conditions string // Encrypted conditions
	Actions    []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type EventStatus struct {
	TriggerID   string
	Status      string
	LastUpdated time.Time
}

type AutomationTask struct {
	TaskID        string
	ScheduledTime time.Time
	TaskDetails   map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ConditionResult struct {
	ID          string
	Met         bool
	EvaluatedAt time.Time
}

type ConditionalEvent struct {
	ConditionID string
	TriggeredAt time.Time
	Status      string
	Details     map[string]interface{}
}

type AutomationEvent struct {
	EventID      string
	EventName    string
	EventDetails map[string]interface{}
	LoggedAt     time.Time
}

type EventPriority struct {
	EventID   string
	Priority  int
	UpdatedAt time.Time
}

type ScheduledEvent struct {
	EventID       string
	EventName     string
	ScheduledTime time.Time
	Details       map[string]interface{}
}

type EventListener struct {
	EventID      string
	ListenerID   string
	RegisteredAt time.Time
}

type ExecutionContext struct {
	ContextID string
	Resources string // Encrypted resource requirements
	Priority  int
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ResourceRequirements struct {
	CPU              int
	Memory           int
	Storage          int
	NetworkBandwidth int
}

type ExecutionStatus struct {
	ContextID   string
	Status      string
	LastUpdated time.Time
}

type ContextTimeout struct {
	ContextID string
	Timeout   time.Duration
	SetAt     time.Time
}

type ExecutionLogEntry struct {
	ContextID string
	Activity  string
	Timestamp time.Time
}

type ExecutionConstraints struct {
	Usage int
	Quota int
}

type EnvironmentConfig struct {
	Variables   map[string]string
	Constraints map[string]string
	Limits      map[string]int
}

type StateHash struct {
	SubBlockID string
	Hash       string
}

type Dispute struct {
	TransactionID string
	Details       string
	Status        string
}

type DisputeResolutionResult struct {
	Resolved  bool
	Outcome   string
	Timestamp time.Time
}

type FinalityLogEntry struct {
	EventID     string
	Description string
	Timestamp   time.Time
}

type ReconciliationStatus struct {
	ContextID string
	IsOngoing bool
	Progress  string
	Issues    []string
}

type ReconciliationLogEntry struct {
	ContextID string
	Result    string
	Timestamp time.Time
}

type FinalizationRecord struct {
	EntityID   string
	StartTime  time.Time
	IsPending  bool
	IsResolved bool
}

type TrapEvent struct {
	ErrorCode int
	Message   string
	Timestamp time.Time
}

type ExceptionLogEntry struct {
	ExceptionID string
	Description string
	Timestamp   time.Time
}

type InterruptHandler struct {
	ID      string
	Handler func() error
}

type SystemHaltLog struct {
	Reason    string
	Timestamp time.Time
}

// Interrupt represents an interrupt with a unique ID, priority, and status.
type Interrupt struct {
	InterruptID string
	Priority    int
	Status      string // e.g., "active", "inactive", "disabled"
}

// TrapCondition defines the conditions under which a trap is triggered.
type TrapCondition struct {
	ConditionID string
	Parameters  map[string]interface{}
}

// EmergencyAlert represents an active emergency alert.
type EmergencyAlert struct {
	Message   string
	IsActive  bool
	Timestamp time.Time
}

// SelfTestResult represents the results of a self-test operation.
type SelfTestResult struct {
	TestName  string
	Result    string // e.g., "success", "failure"
	Timestamp time.Time
}

// CriticalInterrupt represents a registered critical interrupt handler.
type CriticalInterrupt struct {
	InterruptID string
	Handler     func() error
}

// TrapTimeout represents the timeout settings for a specific trap.
type TrapTimeout struct {
	TrapID  string
	Timeout time.Duration
	SetTime time.Time
}

// SafeModeEntry logs the initiation of safe mode with a reason.
type SafeModeEntry struct {
	Reason    string
	Timestamp time.Time
}

// NodeHealth represents the health status of a node in the network.
type NodeHealth struct {
	NodeID        string
	SyncStatus    bool
	ResourceUsage string // e.g., "low", "medium", "high"
	LastChecked   time.Time
}

// NetworkStatus represents the overall health of the network.
type NetworkStatus struct {
	Latency      time.Duration
	NodeStatuses map[string]string // NodeID -> Status
	RecentEvents []string
}

// BlockEvent represents an event related to a specific block.
type BlockEvent struct {
	EventType string
	BlockID   string
	Message   string
	Timestamp time.Time
}

// BlockchainParameter represents a configurable parameter in the blockchain settings.
type BlockchainParameter struct {
	Name  string
	Value string
}

// NetworkTrafficReport represents a report on network traffic.
type NetworkTrafficReport struct {
	Report    string
	Timestamp time.Time
}

// BlockValidationResult represents the result of block validation.
type BlockValidationResult struct {
	BlockID   string
	IsValid   bool
	Timestamp time.Time
}

// ProcessState represents the state of a specific process in the system.
type ProcessState struct {
	ProcessID  string
	Status     string        // e.g., "paused", "running"
	Timeout    time.Duration // Timeout duration for the process.
	LastUpdate time.Time     // Last update timestamp.
}

// SystemLockState represents the lock status of the system.
type SystemLockState struct {
	Reason    string
	IsLocked  bool
	Timestamp time.Time
}

// SystemConstant represents a constant in the system with its value and metadata.
type SystemConstant struct {
	Name      string
	Value     string
	Timestamp time.Time
}

// MaintenanceLog represents a log entry for a system maintenance event.
type MaintenanceLog struct {
	Description string
	Timestamp   time.Time
}

// RecoveryEvent represents a recovery-related event for logging in the ledger.
type RecoveryEvent struct {
	EventID   string
	EventType string
	Timestamp time.Time
}

// DiagnosticEvent represents a diagnostic-related event for logging in the ledger.
type DiagnosticEvent struct {
	EventID   string
	EventType string
	Timestamp time.Time
}

// PanicHandler represents the configuration of a panic handler in the system.
type PanicHandler struct {
	HandlerID string
	Status    string
	Timestamp time.Time
}

// EmergencyOverride represents an emergency override event for logging in the ledger.
type EmergencyOverride struct {
	EventID   string
	Reason    string
	Status    string // e.g., "Activated", "Cancelled"
	Timestamp time.Time
}

// PanicStatus represents the panic status for system assessment.
type PanicStatus struct {
	Status    bool
	Timestamp time.Time
}

// FailedOperation represents a failed operation for cleanup tracking.
type FailedOperation struct {
	OperationID string
	Status      string // e.g., "Cleaned"
	Timestamp   time.Time
}

// ShutdownEvent represents a shutdown sequence event for logging in the ledger.
type ShutdownEvent struct {
	EventID   string
	EventType string
	Timestamp time.Time
}

// ExecutionPolicy represents a policy applied to manage resource allocations.
type ExecutionPolicy struct {
	PolicyID    string
	Description string
	AppliedAt   time.Time
}

// HandoverEvent logs a context handover between nodes.
type HandoverEvent struct {
	ContextID    string
	TargetNodeID string
	Status       string // e.g., "Initiated", "Completed"
	Timestamp    time.Time
}

type RecoveryProtocol struct {
	ProtocolID    string    `json:"protocol_id"`
	RecoverySteps []string  `json:"recovery_steps"`
	DefinedAt     time.Time `json:"defined_at"`
}

type SystemEvent struct {
	EventID     string    `json:"event_id"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

// MetricRecord represents a system metric stored in the ledger.
type MetricRecord struct {
	MetricName  string
	MetricValue float64
	Timestamp   time.Time
}

// SystemEvent represents a system event stored in the ledger.
type SystemEvent struct {
	EventName string
	EventData string
	Timestamp time.Time
}

// HookRecord represents a registered system hook in the ledger.
type HookRecord struct {
	EventID   string
	Priority  int
	Timestamp time.Time
}

// RoleRecord represents a system role and its permissions stored in the ledger.
type RoleRecord struct {
	RoleName    string
	Permissions []string
	Timestamp   time.Time
}

// ProfileRecord represents a system profile loaded into memory.
type ProfileRecord struct {
	ProfileID string
	Timestamp time.Time
}

// SystemProfileManager manages system profiles and configurations.
type SystemProfileManager struct {
	ProfileID          string            // Unique identifier for the profile.
	Configuration      map[string]string // Key-value pairs for profile configurations.
	CreatedAt          time.Time         // Timestamp when the profile was created.
	UpdatedAt          time.Time         // Timestamp for the last update to the profile.
	IsActive           bool              // Indicates if the profile is currently active.
	AssociatedNodes    []string          // List of nodes associated with this profile.
	DefaultPermissions []string          // Default permissions for the profile.
	AuditLogs          []string          // Log entries related to profile changes.
}

// SystemStateSynchronizer ensures synchronization of state across the system.
type SystemStateSynchronizer struct {
	SyncID               string         // Unique identifier for the synchronization session.
	LastSyncTimestamp    time.Time      // Timestamp of the last successful synchronization.
	NodesInSync          []string       // List of nodes currently in sync.
	PendingSyncNodes     []string       // List of nodes pending synchronization.
	SyncFailureLogs      []string       // Logs of synchronization failures.
	RetryPolicy          map[string]int // Retry policies for failed synchronizations (e.g., max retries).
	SyncIntervalDuration time.Duration  // Duration between synchronization attempts.
	IsEnabled            bool           // Indicates whether synchronization is active.
}

// SystemState tracks and stores the current system state.
type SystemState struct {
	StateID         string            // Unique identifier for the state.
	CurrentState    map[string]string // Current system state key-value pairs.
	LastUpdated     time.Time         // Timestamp of the last state update.
	AssociatedTasks []string          // Tasks associated with the current state.
	IsStable        bool              // Indicates if the state is stable or fluctuating.
	ErrorLogs       []string          // Logs of errors affecting the state.
	RecoveryPoint   map[string]string // Recovery point for reverting to a previous state.
}

// SystemManager oversees general system operations and management.
type SystemManager struct {
	ManagerID       string            // Unique identifier for the system manager instance.
	OperationalLogs []string          // Logs of system operations.
	ActiveProcesses []string          // List of currently active processes.
	DowntimeRecords []time.Time       // Records of system downtimes.
	Alerts          []string          // System alerts or warnings.
	MaintenanceMode bool              // Indicates if the system is in maintenance mode.
	Policies        map[string]string // Operational policies and configurations.
}

// OverrideManager handles override configurations and operations.
type OverrideManager struct {
	OverrideID       string            // Unique identifier for the override instance.
	OverrideRules    map[string]string // Rules for override configurations.
	AppliedOverrides []string          // List of overrides currently applied.
	LastOverrideTime time.Time         // Timestamp of the last override applied.
	IsOverrideActive bool              // Indicates if the override is active.
	AuditTrail       []string          // Logs of override actions.
}

// OperationManager manages system-level operations and processes.
type OperationManager struct {
	OperationID         string            // Unique identifier for the operation.
	ActiveOperations    []string          // List of currently active operations.
	CompletedOperations []string          // List of completed operations.
	ErrorLogs           []string          // Logs of errors during operations.
	QueueStatus         []string          // Status of operations in the queue.
	SchedulerPolicies   map[string]string // Scheduling policies for operations.
	LastOperationTime   time.Time         // Timestamp of the last operation execution.
}

// AutomationManager facilitates and monitors automation processes.
type AutomationManager struct {
	AutomationID       string            // Unique identifier for the automation process.
	ActiveAutomations  []string          // List of currently active automation tasks.
	FailedAutomations  []string          // List of failed automation attempts.
	AutomationPolicies map[string]string // Policies governing automation.
	Logs               []string          // Logs of automation events and actions.
	IsAutomationActive bool              // Indicates if automation is currently active.
	NextScheduledTask  time.Time         // Timestamp of the next scheduled automation task.
}

// ************** Governance Structs **************

// GovernanceProposal represents a governance proposal in the network
type GovernanceProposal struct {
	ProposalID       string    // Unique ID of the proposal
	Title            string    // Title of the proposal
	Description      string    // Description of the proposal
	Creator          string    // Address of the proposer
	CreatedAt        time.Time // Timestamp of creation
	Status           ProposalStatus
	VotesFor         int       // Number of votes in favor
	VotesAgainst     int       // Number of votes against
	ExpirationTime   time.Time // Proposal expiration timestamp
	EncryptedDetails string    // Encrypted proposal details
	CreationFee      float64   // Fee charged for proposal creation
}

// PolicyRecord represents a policy within the system, including metadata, rules, and status.
type PolicyRecord struct {
	PolicyID             string           // Unique identifier for the policy
	Name                 string           // Descriptive name of the policy
	Description          string           // Detailed description of the policy
	CreatedBy            string           // Creator of the policy (e.g., node, user, or entity)
	CreationTimestamp    time.Time        // Timestamp when the policy was created
	LastUpdatedBy        string           // Identifier of the last updater
	LastUpdatedTimestamp time.Time        // Timestamp of the last update
	Status               string           // Current status of the policy (e.g., active, suspended, revoked)
	Scope                string           // Scope of the policy (e.g., global, local, entity-specific)
	Rules                []PolicyRule     // List of rules or conditions within the policy
	EnforcementMechanism string           // Mechanism used for policy enforcement
	ViolationPenalties   []PenaltyRecord  // Penalties associated with policy violations
	AssociatedEntities   []string         // Entities (e.g., users, nodes) associated with the policy
	Dependencies         []string         // Dependencies or prerequisites for the policy
	RevisionHistory      []PolicyRevision // History of changes made to the policy
	AuditLogs            []PolicyAuditLog // Logs of policy-related activities or checks
	ExpirationDate       *time.Time       // Optional expiration date for the policy
	Tags                 []string         // Tags or categories for easy retrieval or classification
}

// PolicyRule represents a single rule or condition within a policy.
type PolicyRule struct {
	RuleID        string    // Unique identifier for the rule
	Condition     string    // Condition that must be met (e.g., "Balance > 1000")
	Action        string    // Action to take if the condition is met (e.g., "Allow", "Deny")
	Description   string    // Description of the rule
	Priority      int       // Priority of the rule (lower number = higher priority)
	Enabled       bool      // Whether the rule is currently enabled
	LastEvaluated time.Time // Timestamp of the last evaluation
}

// PenaltyRecord represents a penalty imposed for violating a policy.
type PenaltyRecord struct {
	PenaltyID      string    // Unique identifier for the penalty
	Type           string    // Type of penalty (e.g., "Fine", "Ban", "Warning")
	Amount         float64   // Amount or severity of the penalty
	Description    string    // Description of the penalty
	AssociatedRule string    // Rule that triggered the penalty
	Timestamp      time.Time // Timestamp of when the penalty was applied
}

// PolicyRevision represents a single revision in the policy's history.
type PolicyRevision struct {
	RevisionID    string    // Unique identifier for the revision
	UpdatedBy     string    // Entity that made the revision
	Timestamp     time.Time // Timestamp of the revision
	ChangeSummary string    // Summary of changes made in the revision
}

// PolicyAuditLog represents an audit log entry for policy-related activities.
type PolicyAuditLog struct {
	LogID       string    // Unique identifier for the log entry
	Action      string    // Action performed (e.g., "Created", "Updated", "Deleted")
	PerformedBy string    // Entity that performed the action
	Timestamp   time.Time // Timestamp of the action
	Details     string    // Additional details about the action
}

// ProposalManager manages the lifecycle of governance proposals
type ProposalManager struct {
	Proposals      map[string]*GovernanceProposal // Map of proposal ID to proposals
	mutex          sync.Mutex                     // Mutex for thread-safe operations
	LedgerInstance *Ledger                        // Ledger instance for tracking proposals
	FeePercentage  float64                        // Fee percentage based on transaction fees (0.25%)
}

// ProposalStatus represents the status of a governance proposal
type ProposalStatus string

const (
	Pending  ProposalStatus = "Pending"
	Approved ProposalStatus = "Approved"
	Rejected ProposalStatus = "Rejected"
)

// GovernanceVoting handles the voting mechanism for governance proposals
type GovernanceVoting struct {
	Votes          map[string]map[string]bool // Map[proposalID]map[voterID]bool
	LedgerInstance *Ledger                    // Ledger to store voting records
	mutex          sync.Mutex                 // Mutex for thread-safe operations
}

// ExecutionRecord keeps track of the governance actions to be executed after voting
type ExecutionRecord struct {
	ProposalID string    // ID of the proposal to be executed
	Executed   bool      // Whether the execution was successful
	Timestamp  time.Time // Time when the execution happened
}

// GovernanceExecution manages the execution of approved proposals
type GovernanceExecution struct {
	ExecutionQueue []ExecutionRecord // Queue of proposals to be executed
	mutex          sync.Mutex        // Mutex for thread-safe operations
	LedgerInstance *Ledger           // Ledger instance for tracking executed proposals
}

// GovernanceTimelock defines the timelock contract for governance
type GovernanceTimelock struct {
	PendingExecutions map[string]*TimelockExecution // Map of proposal ID to pending executions
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	LedgerInstance    *Ledger                       // Ledger instance for logging execution
}

// TimelockExecution represents a proposal pending execution after timelock
type TimelockExecution struct {
	ProposalID        string    // ID of the proposal
	ExecutionTime     time.Time // Time when the proposal can be executed
	EncryptedProposal string    // Encrypted details of the proposal
	Creator           string    // Address of the proposer
}

// GovernanceTracking manages the tracking of governance proposals and generates reports
type GovernanceTracking struct {
	ProposalHistory map[string]*GovernanceProposalStatus // Map of proposal ID to status
	LedgerInstance  *Ledger                              // Instance of the ledger
	mutex           sync.Mutex                           // Mutex for thread-safe operations
}

// GovernanceProposalStatus tracks the status of a single proposal
type GovernanceProposalStatus struct {
	ProposalID       string      // ID of the proposal
	Status           string      // Current status (e.g., "Pending", "Approved", "Executed", "Rejected")
	Timestamps       []time.Time // Important timestamps (e.g., creation, approval, execution)
	EncryptedDetails string      // Encrypted details of the proposal
}

// ReputationVoting represents the structure for reputation-based voting
type ReputationVoting struct {
	Votes          map[string]map[string]float64 // Map[proposalID]map[voterID]reputationScore
	LedgerInstance *Ledger                       // Ledger to store voting records
	mutex          sync.Mutex                    // Mutex for thread-safe operations
}

// GovernanceRecord holds governance-specific data for the ledger.
type GovernanceRecord struct {
	Proposals       map[string]GovernanceProposal // Governance proposals
	TransactionFees map[uint64]float64            // Tracks transaction fees per block
	Delegations     map[string]string             // Tracks delegations from one user to another
}

// ************** High Availability Structs **************

// DataBackupManager is responsible for backing up and restoring blockchain data
type DataBackupManager struct {
	LedgerInstance *Ledger                        // Ledger instance to back up
	Backups        map[string][]*BlockchainBackup // A map of node IDs to their list of backups
	BackupInterval time.Duration                  // Interval for automatic backups
	BackupLocation string                         // Directory where backups are stored
	mutex          sync.Mutex                     // Mutex for thread-safe backup operations
}

// FailoverThreshold represents thresholds for failover conditions.
type FailoverThreshold struct {
	MaxAllowedDowntime time.Duration // Maximum allowed downtime before failover
	MinHealthyNodes    int           // Minimum number of healthy nodes required
	FailureRate        float64       // Failure rate percentage threshold
}

// ChainForkManager handles the detection and resolution of blockchain forks
type ChainForkManager struct {
	LedgerInstance *Ledger    // The ledger instance to track the chain state
	ForkedChains   [][]Block  // List of forked chains detected
	mutex          sync.Mutex // Mutex for thread-safe operations
}

// DataCollectionManager is responsible for collecting and distributing blockchain data across nodes
type DataCollectionManager struct {
	Nodes                 []string      // List of nodes in the network
	CollectedTransactions []Transaction // Transactions collected from the network
	CollectedSubBlocks    []SubBlock    // Sub-blocks collected from the network
	CollectedBlocks       []Block       // Blocks collected from the network
	mutex                 sync.Mutex    // Mutex for thread-safe data collection
	DataRequestInterval   time.Duration // Interval for requesting data from nodes
}

// DataDistributionManager is responsible for distributing blockchain data to all nodes in the network
type DataDistributionManager struct {
	Nodes                []string   // List of nodes in the network
	DistributedSubBlocks []SubBlock // Sub-blocks that have been distributed
	DistributedBlocks    []Block    // Blocks that have been distributed
	mutex                sync.Mutex // Mutex for thread-safe data distribution
}

// DataReplicationManager is responsible for replicating blockchain data across nodes for redundancy and fault tolerance
type DataReplicationManager struct {
	Nodes               []string   // List of nodes in the network
	ReplicatedSubBlocks []SubBlock // Sub-blocks that have been replicated
	ReplicatedBlocks    []Block    // Blocks that have been replicated
	mutex               sync.Mutex // Mutex for thread-safe replication operations
}

// DataSynchronizationManager handles the synchronization of blockchain data between nodes
type DataSynchronizationManager struct {
	Nodes           []string   // List of nodes in the network
	LatestSubBlocks []SubBlock // Latest sub-blocks to synchronize
	LatestBlocks    []Block    // Latest blocks to synchronize
	mutex           sync.Mutex // Mutex for thread-safe operations
}

// DisasterRecoveryManager handles disaster recovery strategies for the blockchain
type DisasterRecoveryManager struct {
	BackupNodes       []string           // List of backup nodes for failover
	DataBackupManager *DataBackupManager // Manager responsible for handling backups
	RecoveryLog       []string           // Log of recovery operations
	mutex             sync.Mutex         // Mutex for thread-safe recovery operations
}

// HealthMonitoringManager monitors the health and performance of the blockchain network
type HealthMonitoringManager struct {
	Nodes            []string              // List of nodes to monitor
	NodeHealthStatus map[string]bool       // Health status of each node
	SubBlockLatency  map[int]time.Duration // Sub-block validation latency
	BlockLatency     map[int]time.Duration // Block validation latency
	LedgerInstance   *Ledger               // Instance of the ledger for validation tracking
	mutex            sync.Mutex            // Mutex for thread-safe operations
}

// HeartbeatService is responsible for sending and monitoring heartbeats between nodes
type HeartbeatService struct {
	Nodes          []string             // List of nodes to monitor
	HeartbeatLogs  map[string]time.Time // Records of last heartbeat received from each node
	Interval       time.Duration        // Interval between heartbeat checks
	mutex          sync.Mutex           // Mutex for thread-safe operations
	LedgerInstance *Ledger              // Ledger instance for storing heartbeat data
}

// RecoveryManager handles system recovery processes, including checkpoints and state restoration.
type RecoveryManager struct {
	RecoveryID       string                 // Unique identifier for the recovery instance.
	RecoveryType     string                 // Type of recovery (e.g., crash recovery, checkpoint recovery).
	Description      string                 // Description of the recovery process and its purpose.
	Checkpoints      []RecoveryCheckpoint   // List of recovery checkpoints.
	LastRecoveryTime time.Time              // Timestamp of the last successful recovery.
	CreatedAt        time.Time              // Timestamp when the recovery instance was created.
	Metadata         map[string]interface{} // Additional metadata related to the recovery process.
}

// RecoveryCheckpoint represents a checkpoint in the recovery process.
type RecoveryCheckpoint struct {
	CheckpointID string    // Unique identifier for the checkpoint.
	StateData    []byte    // Serialized state data captured at the checkpoint.
	CreatedAt    time.Time // Timestamp when the checkpoint was created.
	CreatedBy    string    // Identifier of the entity that created the checkpoint.
	IsValid      bool      // Indicates whether the checkpoint is valid for recovery.
}

// FallbackManager manages fallback strategies for systems during failures or disruptions.
type FallbackManager struct {
	FallbackID    string                 // Unique identifier for the fallback instance.
	Description   string                 // Description of the fallback mechanism.
	Strategies    []FallbackStrategy     // List of fallback strategies.
	LastTriggered time.Time              // Timestamp of the last fallback event.
	CreatedAt     time.Time              // Timestamp when the fallback instance was created.
	Metadata      map[string]interface{} // Additional metadata related to the fallback process.
}

// FallbackStrategy represents an individual fallback mechanism.
type FallbackStrategy struct {
	StrategyID       string    // Unique identifier for the fallback strategy.
	TriggerCondition string    // Condition that triggers the fallback strategy.
	Actions          []string  // Actions to execute as part of the fallback.
	Priority         int       // Priority level of the fallback strategy.
	CreatedAt        time.Time // Timestamp when the strategy was created.
}

// SystemBackupManager handles system backup processes for disaster recovery and data integrity.
type SystemBackupManager struct {
	BackupID       string                 // Unique identifier for the backup instance.
	BackupType     string                 // Type of backup (e.g., full, incremental).
	TargetLocation string                 // Destination for the backup data.
	Schedule       BackupSchedule         // Backup schedule and frequency.
	LastBackupTime time.Time              // Timestamp of the last completed backup.
	CreatedAt      time.Time              // Timestamp when the backup instance was created.
	Metadata       map[string]interface{} // Additional metadata related to the backup process.
}

// BackupSchedule represents the scheduling details for backups.
type BackupSchedule struct {
	ScheduleID     string    // Unique identifier for the schedule.
	Frequency      string    // Frequency of the backup (e.g., daily, weekly).
	NextBackupTime time.Time // Timestamp for the next scheduled backup.
	CreatedAt      time.Time // Timestamp when the schedule was created.
}

// MirroringManager manages data mirroring for redundancy and fault tolerance.
type MirroringManager struct {
	MirrorID       string                 // Unique identifier for the mirroring instance.
	SourceNode     string                 // Source node for the mirroring process.
	TargetNodes    []string               // Target nodes receiving mirrored data.
	Status         string                 // Current status of the mirroring process.
	LastSyncedTime time.Time              // Timestamp of the last successful synchronization.
	CreatedAt      time.Time              // Timestamp when the mirroring instance was created.
	Metadata       map[string]interface{} // Additional metadata related to the mirroring process.
}

// NodeFailoverManager handles the failover mechanism in case a node goes down
type NodeFailoverManager struct {
	PrimaryNodes     []string        // List of primary nodes
	BackupNodes      []string        // List of backup nodes to failover to
	NodeHealthStatus map[string]bool // Health status of each node
	CurrentPrimary   string          // The current active primary node
	LedgerInstance   *Ledger         // Ledger instance for managing state and transactions
	mutex            sync.Mutex      // Mutex for thread-safe operations
}

// NodeMonitoringService monitors node performance and health for high availability
type NodeMonitoringService struct {
	Nodes          map[string]*NodeMetrics // Map of node addresses to their metrics
	LedgerInstance *Ledger                 // Ledger for recording node health data
	mutex          sync.Mutex              // Mutex for thread-safe operations
	CheckInterval  time.Duration           // Interval between each health check
	FaultThreshold int                     // Threshold for marking a node as faulty
}

// BlockchainBackup represents a backup of the blockchain state at a given point in time
type BlockchainBackup struct {
	BackupID     string    // Unique identifier for the backup
	Timestamp    time.Time // The time the backup was created
	Blocks       []Block   // List of blocks included in the backup
	NodeID       string    // Node ID that created the backup
	BackupSize   int64     // Size of the backup in bytes
	BackupHash   string    // Hash to verify the integrity of the backup
	IsCompressed bool      // Whether the backup is compressed
}

// ReplicationLog holds data about log replication events.
type ReplicationLog struct {
	NodeID    string
	Timestamp time.Time
	Status    string
}

// NodeMetric holds metrics data for a node.
type NodeMetric struct {
	NodeID           string
	LastReportedTime time.Time
	CPUUsage         float64
	MemoryUsage      float64
	DiskUsage        float64
}

// Snapshot represents a data snapshot for high availability.
type Snapshot struct {
	SnapshotID string
	CreatedAt  time.Time
	Data       string
	Metadata   map[string]string
}

// SnapshotStatus represents the current status of ongoing snapshots.
type SnapshotStatus struct {
	ActiveSnapshots  int
	LastSnapshotID   string
	LastSnapshotTime time.Time
}

// MirroringStatus represents the status of data mirroring processes.
type MirroringStatus struct {
	IsActive          bool
	LastMirroringTime time.Time
	MirroringRate     int
}

// QuorumPolicy represents the policy configuration for quorum decisions.
type QuorumPolicy struct {
	PolicyID  string
	Policy    string
	Encrypted bool
}

// QuorumStatus represents the status of quorum activities.
type QuorumStatus struct {
	IsActive       bool
	LastActivity   time.Time
	DecisionCount  int
	CurrentMembers []string
}

// LoadBalancerStatus represents the status of the load balancer.
type LoadBalancerStatus struct {
	IsActive         bool
	ActiveNodesCount int
	CurrentPolicy    string
}

// RecoveryTimeoutConfig represents the recovery timeout configuration.
type RecoveryTimeoutConfig struct {
	TimeoutSeconds int
	ConfiguredAt   time.Time
}

// ArchiveRetentionPolicy represents the retention policy for archived data.
type ArchiveRetentionPolicy struct {
	PolicyName   string
	ConfiguredAt time.Time
}

// ConsistencyCheckResult represents the result of a consistency check.
type ConsistencyCheckResult struct {
	CheckID     string
	Timestamp   time.Time
	IssuesFound int
	Resolved    bool
}

// PredictiveScalingPolicy represents the policy configuration for predictive scaling.
type PredictiveScalingPolicy struct {
	PolicyName   string
	ConfiguredAt time.Time
}

// PredictiveFailoverConfig represents the configuration for predictive failover.
type PredictiveFailoverConfig struct {
	IsEnabled    bool
	ConfiguredAt time.Time
}

// PredictiveFailoverPolicy represents the failover policy configuration.
type PredictiveFailoverPolicy struct {
	PolicyName   string
	ConfiguredAt time.Time
}

// AdaptiveResourcePolicy represents the adaptive resource management policy.
type AdaptiveResourcePolicy struct {
	PolicyName   string
	ConfiguredAt time.Time
}

// DisasterRecoveryPlan holds the disaster recovery plan details.
type DisasterRecoveryPlan struct {
	PlanName     string
	ConfiguredAt time.Time
}

// DisasterRecoveryBackup represents a disaster recovery backup.
type DisasterRecoveryBackup struct {
	BackupName string
	CreatedAt  time.Time
	Data       string // Encrypted data for backup
}

// WriteAheadLogConfig holds configuration for write-ahead logging.
type WriteAheadLogConfig struct {
	Enabled       bool
	RetentionDays int
}

// LogRetentionConfig represents the configuration for log retention.
type LogRetentionConfig struct {
	RetentionPeriod int
	ConfiguredAt    time.Time
}

// HighAvailabilityMode represents the current high-availability mode and configuration.
type HighAvailabilityMode struct {
	Mode         string
	ConfiguredAt time.Time
}

// FailoverConfig holds the timeout and settings for failover.
type FailoverConfig struct {
	Timeout      int
	ConfiguredAt time.Time
}

// FailoverStatus holds the status of the failover process.
type FailoverStatus struct {
	CurrentStatus string
	LastUpdated   time.Time
}

// AutoScalingConfig represents the configuration for auto-scaling.
type AutoScalingConfig struct {
	Policy       string
	Enabled      bool
	ConfiguredAt time.Time
}

// RecoveryPoint represents a recovery point in the system.
type RecoveryPoint struct {
	PointID     string
	CreatedAt   time.Time
	Description string
}

// BackupStatus represents the current status of a backup.
type BackupStatus struct {
	BackupName  string
	Status      string // e.g., "InProgress", "Completed", "Failed"
	LastUpdated time.Time
}

// SnapshotDetails holds the details of a snapshot.
type SnapshotDetails struct {
	SnapshotName string
	CreatedAt    time.Time
	Size         int64 // Size in bytes
	Description  string
}

// FailoverGroup represents a failover group with members and policies.
type FailoverGroup struct {
	GroupID     string
	Members     []string
	Policy      string
	LastUpdated time.Time
}

// HAProxyConfig represents the configuration for HA Proxy services.
type HAProxyConfig struct {
	Policy      string
	Enabled     bool
	LastUpdated time.Time
}

// ResourcePool represents the resource pooling configuration.
type ResourcePool struct {
	Enabled     bool
	LastUpdated time.Time
}

// ResourcePoolingPolicy represents the policy for resource pooling.
type ResourcePoolingPolicy struct {
	Policy      string
	LastUpdated time.Time
}

// GeoRedundancyPolicy represents the policy for geographic redundancy.
type GeoRedundancyPolicy struct {
	Policy      string
	LastUpdated time.Time
}

// DisasterSimulationConfig represents the configuration for disaster simulations.
type DisasterSimulationConfig struct {
	Mode        string
	Parameters  string
	LastUpdated time.Time
}

// ReplicationConfig stores replication settings.
type ReplicationConfig struct {
	ReplicationFactor int
	LastUpdated       time.Time
}

// ClusterConfig represents the cluster configuration.
type ClusterConfig struct {
	Nodes       []string
	LastUpdated time.Time
}

// HighAvailabilityConfig represents the overall HA configuration.
type HighAvailabilityConfig struct {
	LoadBalancingEnabled bool
	ReplicationEnabled   bool
	ClusteringEnabled    bool
	LastUpdated          time.Time
}

// ClusterPolicy represents a clustering policy with configuration details.
type ClusterPolicy struct {
	PolicyName  string
	Rules       string
	LastUpdated time.Time
}

// HeartbeatConfig stores configuration for heartbeat monitoring.
type HeartbeatConfig struct {
	Interval    int
	Enabled     bool
	LastUpdated time.Time
}

// HealthCheckConfig stores configuration for health checks.
type HealthCheckConfig struct {
	Interval    int
	Enabled     bool
	LastUpdated time.Time
}

// ReplicaConfig stores configuration details for replication.
type ReplicaConfig struct {
	Count            int
	RedundancyLevel  int
	CompressionLevel int
	LastUpdated      time.Time
}

// SynchronizationConfig stores details for data synchronization.
type SynchronizationConfig struct {
	Interval    int
	IsEnabled   bool
	LastUpdated time.Time
}

// CompressionConfig stores details for data compression.
type CompressionConfig struct {
	CompressionLevel int
	IsEnabled        bool
	LastUpdated      time.Time
}

// RedundancyConfig stores redundancy configuration details.
type RedundancyConfig struct {
	Level       int
	LastUpdated time.Time
}

// DeduplicationConfig stores deduplication settings.
type DeduplicationConfig struct {
	Policy      string
	IsEnabled   bool
	LastUpdated time.Time
}

// StandbyConfig stores settings for standby modes.
type StandbyConfig struct {
	Mode        string // "hot", "cold", or "none"
	Policy      string
	IsEnabled   bool
	LastUpdated time.Time
}

// SelfHealingConfig represents self-healing settings and their intervals.
type SelfHealingConfig struct {
	IsEnabled        bool
	Interval         int // Interval in seconds
	FailbackPriority int
	LastUpdated      time.Time
}

// ArchivedData represents metadata for archived data in the ledger.
type ArchivedData struct {
	ID         string
	Data       string
	ArchivedAt time.Time
}

// SimulationResult represents simulation results stored in the ledger.
type SimulationResult struct {
	ID        string
	Results   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ResourceQuotaConfig represents resource quota settings.
type ResourceQuotaConfig struct {
	Limits      string
	IsEnabled   bool
	LastUpdated time.Time
}

// ResourceScalingConfig represents resource scaling settings.
type ResourceScalingConfig struct {
	ScalingLimits string
	Policy        string
	LastUpdated   time.Time
}

// ************** Identity Services Structs **************

// AccessControlManager manages access control for authority nodes, users, and stakeholders
type AccessControlManager struct {
	AuthorityNodes    map[string]*AuthorityNodeTypes // Map of NodeID to AuthorityNode
	UserAccess        map[string]string              // Map of UserID to their access level (encrypted)
	StakeholderAccess map[string]string              // Map of StakeholderID to their access level (encrypted)
	OwnerAccess       string                         // Owner's access secret
	LedgerInstance    *Ledger                        // Ledger for recording access changes
	mutex             sync.Mutex                     // Mutex for thread-safe access control
}

// Identity represents an individual's identity on the blockchain.
type Identity struct {
	IdentityID    string       // Unique ID for the identity
	IdentityType  IdentityType // The type of identity (Syn900 or DecentralizedID)
	Owner         string       // Owner of the identity (wallet address or public key)
	CreatedAt     time.Time    // Timestamp of identity creation
	EncryptedData string       // Encrypted identity details
	IsVerified    bool         // Whether the identity has been verified
	RecoverySetup string       // Recovery setup information for the wallet (new field)
	WalletID      string       // Wallet ID associated with the identity
}

// IdentityType defines the types of identities in the system
type IdentityType string

// IdentityVerificationManager handles creation and verification of identities
type IdentityVerificationManager struct {
	Identities     map[string]*Identity // Map of identityID to identity
	LedgerInstance *Ledger              // Ledger for recording identity-related actions
	mutex          sync.Mutex           // Mutex for thread-safe identity operations
}

// PrivacySettings defines the settings for user privacy
type PrivacySettings struct {
	UserID            string    // User's ID (wallet address, public key, or other unique identifier)
	DataEncryption    bool      // Whether the user's data is encrypted
	PermissionToShare bool      // Whether the user permits sharing their data
	LastUpdated       time.Time // The last time the privacy settings were updated
}

// PrivacyManager handles the management of user privacy settings
type PrivacyManager struct {
	PrivacyRecords map[string]*PrivacySettings // Map of user ID to privacy settings
	mutex          sync.Mutex                  // Mutex for thread-safe privacy operations
	LedgerInstance *Ledger                     // Ledger instance for logging privacy actions
}

// IdentityLog tracks identity-related actions in the ledger.
type IdentityLog struct {
	NodeID    string
	Action    string
	Timestamp time.Time
	Details   string
}

// ************** Integration Structs **************

// ServiceProvider represents a service provider entity
type ServiceProvider struct {
	ProviderID   string
	Name         string
	Details      string
	EncryptedKey string
}

// ApplicationUpdate represents an update to an application
type ApplicationUpdate struct {
	Version     string
	Description string
	UpdateData  string
}

// IntegrationState represents the state of an integration or proxy
type IntegrationState struct {
	ProxyID    string
	Status     string
	LastUpdate time.Time
}

// APIProxyConfig represents an API proxy configuration
type APIProxyConfig struct {
	ProxyID    string
	ConfigData string
	CreatedAt  time.Time
}

// DependencyManager manages dependencies between system components, processes, or modules.
type DependencyManager struct {
	DependencyID    string                 // Unique identifier for the dependency instance.
	SourceComponent string                 // Name or ID of the source component.
	TargetComponent string                 // Name or ID of the dependent component.
	DependencyType  string                 // Type of dependency (e.g., hard, soft, optional).
	Status          string                 // Current status of the dependency (e.g., active, inactive).
	LastUpdated     time.Time              // Timestamp of the last update to the dependency.
	CreatedAt       time.Time              // Timestamp when the dependency instance was created.
	Metadata        map[string]interface{} // Additional metadata related to the dependency.
}

// DependencyLink represents a specific link between two components.
type DependencyLink struct {
	LinkID     string    // Unique identifier for the dependency link.
	Source     string    // Name or ID of the source component.
	Target     string    // Name or ID of the target component.
	Dependency string    // Description of the dependency relationship.
	CreatedAt  time.Time // Timestamp when the dependency link was created.
	IsCritical bool      // Indicates if the dependency is critical for operation.
}

// HandlerManager manages event handlers, process handlers, and system interaction points.
type HandlerManager struct {
	HandlerID       string                 // Unique identifier for the handler instance.
	HandlerType     string                 // Type of handler (e.g., event, process, request).
	Description     string                 // Description of the handler's purpose.
	AssignedProcess string                 // ID of the process or component assigned to this handler.
	Status          string                 // Current status of the handler (e.g., active, inactive).
	LastExecuted    time.Time              // Timestamp of the last execution of the handler.
	CreatedAt       time.Time              // Timestamp when the handler instance was created.
	Metadata        map[string]interface{} // Additional metadata related to the handler.
}

// HandlerConfiguration represents configuration details for handlers.
type HandlerConfiguration struct {
	ConfigID   string                 // Unique identifier for the handler configuration.
	HandlerID  string                 // ID of the handler the configuration applies to.
	Parameters map[string]interface{} // Configuration parameters for the handler.
	CreatedAt  time.Time              // Timestamp when the configuration was created.
	IsDefault  bool                   // Indicates if this is the default configuration for the handler.
}

// ServiceConfig represents a service configuration.
type ServiceConfig struct {
	ConfigID       string
	ConfigData     string
	LastUpdated    time.Time
	EncryptionHash string
}

// WebhookConfig represents a webhook configuration for a service.
type WebhookConfig struct {
	WebhookID      string
	URL            string
	EventType      string
	Authentication string
	CreatedAt      time.Time
}

// CustomFunction represents a custom function for a service.
type CustomFunction struct {
	FunctionID     string
	FunctionCode   string
	Description    string
	CreatedAt      time.Time
	EncryptionHash string
}

// AnalyticsConfig represents an analytics tool configuration.
type AnalyticsConfig struct {
	AnalyticsID     string
	ToolName        string
	ToolVersion     string
	IntegrationDate time.Time
	EncryptionHash  string
}

// AppConfig represents the configuration for an application.
type AppConfig struct {
	ConfigID       string
	Parameters     map[string]string
	LastUpdated    time.Time
	EncryptionHash string
}

// IntegrationParams represent parameters for integrating an application.
type IntegrationParams struct {
	ParamID        string
	Values         map[string]string
	LastUpdated    time.Time
	EncryptionHash string
}

// APISchema represents the schema definition for API validation.
type APISchema struct {
	SchemaID    string
	Version     string
	Definition  string
	LastUpdated time.Time
}

// Extension represents an extension for a DApp.
type Extension struct {
	ExtensionID   string
	Name          string
	Version       string
	Description   string
	EncryptedData string
}

// EventHandler represents an event handler configuration for applications.
type EventHandler struct {
	HandlerID     string
	HandlerName   string
	TriggerEvent  string
	HandlerAction string
	RegisteredAt  time.Time
}

// Library represents a library integrated into an application.
type Library struct {
	LibraryID    string
	Name         string
	Version      string
	Description  string
	Encrypted    string
	IntegratedAt time.Time
}

// APIKeys represents a set of API keys for an application.
type APIKeys struct {
	KeyID       string
	Keys        map[string]string
	LastUpdated time.Time
	Encrypted   string
}

// Service represents a service integration for an application.
type Service struct {
	ServiceID   string
	Name        string
	Description string
	Status      string
	ActivatedAt time.Time
}

// IntegrationStatus represents the status of application integration.
type IntegrationStatus struct {
	AppID          string    // Unique identifier for the application or DApp
	IsIntegrated   bool      // Indicates whether the app is successfully integrated
	LastChecked    time.Time // Timestamp of the last integration status check
	IntegrationLog string    // Log message or details about the integration status
}

// DappMetadata holds metadata about a decentralized application.
type DappMetadata struct {
	DappID      string
	Name        string
	Description string
	Version     string
	Owner       string
	CreatedAt   time.Time
	Encrypted   string // Encrypted metadata
}

// APIEndpoint represents an API endpoint linked to a DApp.
type APIEndpoint struct {
	EndpointID  string
	Name        string
	URL         string
	Description string
	Encrypted   string // Encrypted endpoint data
}

// CLICommand represents a CLI command for a DApp.
type CLICommand struct {
	CommandID   string
	Name        string
	Description string
	Script      string
	Encrypted   string // Encrypted command data
}

// Functionality represents new functionality added to a DApp.
type Functionality struct {
	FunctionalityID string
	Name            string
	Description     string
	Code            string
	Encrypted       string // Encrypted functionality data
}

// FeatureToggle represents a feature toggle for a DApp.
type FeatureToggle struct {
	DappID      string
	FeatureName string
	Enabled     bool
	Timestamp   time.Time
}

// ExternalService represents an external service integrated with a DApp.
type ExternalService struct {
	ServiceID   string
	Name        string
	Description string
	APIEndpoint string
	Encrypted   string // Encrypted service details
}

// Opcode represents a bytecode-level opcode for a DApp.
type Opcode struct {
	OpcodeID    string
	Name        string
	Description string
	Execution   string // Bytecode or execution logic
}

// AppComponent represents an application component for a DApp.
type AppComponent struct {
	ComponentID string
	Name        string
	Version     string
	Description string
	Encrypted   string // Encrypted component data
}

// Feature represents a specific feature available for an application.
type Feature struct {
	FeatureID     string
	Name          string
	Description   string
	Dependencies  []string
	Enabled       bool
	LastValidated time.Time
}

// Dependency represents a dependency for a feature.
type Dependency struct {
	DependencyID string
	Name         string
	Type         string // Could be service, feature, etc.
	Version      string
	Encrypted    string
}

// IntegrationMapping represents cross-application data sharing mappings.
type IntegrationMapping struct {
	SourceAppID  string
	TargetAppID  string
	MappingRules string
	Encrypted    string
}

// WorkflowConfig represents the configuration for an integration workflow.
type WorkflowConfig struct {
	WorkflowID string
	Name       string
	Steps      []string // List of orchestrated steps
	Encrypted  string
}

// SecurityReview represents the security assessment of an application's integration.
type SecurityReview struct {
	AppID         string
	Findings      []string
	ReviewDate    time.Time
	Reviewer      string
	OverallStatus string // Example: "Secure", "Issues Found"
}

// ActivityLog represents the activity log for integration events.
type ActivityLog struct {
	LogID     string
	AppID     string
	Timestamp time.Time
	Event     string
	Details   string
	Encrypted string
}

// CrossAppFunction represents a function invoked between DApps.
type CrossAppFunction struct {
	FunctionID string
	Name       string
	Parameters map[string]interface{}
	Encrypted  string
}

// Module represents a module dependency for a feature.
type Module struct {
	ModuleID    string
	Name        string
	Version     string
	Description string
	Encrypted   string
}

// Policy represents a policy governing service integration.
type Policy struct {
	PolicyID     string
	Name         string
	Description  string
	CreatedAt    time.Time
	LastModified time.Time
	Encrypted    string
}

// IntegrationEvent represents an event related to service integration.
type IntegrationEvent struct {
	EventID      string
	ServiceID    string
	Timestamp    time.Time
	EventDetails string
	Encrypted    string
}

// AccessLevel represents the access permissions for a service integration.
type AccessLevel struct {
	LevelName string // Example: "Read-Only", "Full-Access"
	Details   string
}

// LogLevel represents the logging level for integration activities.
type LogLevel struct {
	LevelName string // Example: "Info", "Debug", "Error"
}

// IntegrationLog represents a log entry for a service integration.
type IntegrationLog struct {
	LogID      string
	ServiceID  string
	Timestamp  time.Time
	LogDetails string
	Severity   string // Example: "Info", "Warning", "Error"
}

// HealthStatus represents the health status of a service integration.
type HealthStatus struct {
	ServiceID string
	Status    string // Example: "Healthy", "Degraded", "Down"
	LastCheck time.Time
}

// TestConfig represents configuration for an integration test.
type TestConfig struct {
	TestID   string
	Params   map[string]string
	Expected string
}

// CLITool represents a CLI tool for integration purposes.
type CLITool struct {
	Name          string
	Version       string
	Compatibility string
	Description   string
}

// ************** Interoperability Structs **************

// CrossChainManager handles operations related to managing transactions and communication between multiple blockchain networks.
type CrossChainManager struct {
	Bridges          map[string]*Bridge             // Map of active cross-chain bridges by bridge ID
	PendingTransfers map[string]*CrossChainTransfer // Map of pending cross-chain transfers, identified by transfer ID
	ActiveNetworks   []string                       // List of active blockchain networks the manager interacts with
	SyncInterval     time.Duration                  // Frequency at which cross-chain syncing operations are performed
	TransferFee      float64                        // Fee applied to cross-chain transfers
	mutex            sync.Mutex                     // Mutex for ensuring thread-safe cross-chain operations
}

// ValidationLog represents a log entry for cross-chain validation activities.
type ValidationLog struct {
	LogID       string    // Unique identifier for the log
	ValidatorID string    // ID of the validator performing the validation
	Details     string    // Details of the validation activity
	Timestamp   time.Time // Timestamp of the validation
	Status      string    // Status of the validation (e.g., "success", "failure")
}

// InteroperabilityLog tracks cross-chain and atomic swap events in the ledger.
type InteroperabilityLog struct {
	EventType string
	Timestamp time.Time
	Details   string
	Status    string
}

type DataFeed struct {
    FeedID      string
    Data        string
    Timestamp   time.Time
    Validated   bool
    SourceChain string
}

type ExternalData struct {
    DataID   string
    Data     string
    Accuracy bool
}

type DisputeEvidence struct {
    EvidenceID   string
    DisputeID    string
    Content      string
    Validated    bool
    SubmittedAt  time.Time
    ValidatedAt  *time.Time
}

type ArbitrationSummary struct {
    SummaryID    string
    DisputeID    string
    Summary      string
    GeneratedAt  time.Time
}

type CrossChainAssetLog struct {
    AssetID         string
    ChainID         string
    TransactionType string
    Details         string
    Timestamp       time.Time
    Status          string
}

type AssetHistory struct {
    AssetID        string
    TransactionID  string
    TransactionDetails string
    Timestamp      time.Time
}

type AssetHistoryRecord struct {
    AssetID      string
    TransactionID string
    Details      string
    Timestamp    time.Time
    Verified     bool
}

type CrossChainEvent struct {
    EventID      string
    AssetID      string
    EventType    string
    Details      string
    Timestamp    time.Time
}

type CrossChainState struct {
    StateID       string
    TargetChainID string
    SyncStatus    string
    Timestamp     time.Time
}

type CrossChainSettlement struct {
    SettlementID      string
    SourceChainID     string
    DestinationChainID string
    Amount            float64
    Timestamp         time.Time
    Status            string
}

type CrossChainActivity struct {
    ActivityID string
    Status     string
    Reason     string
    Timestamp  time.Time
}

type NodeLatency struct {
    NodeID    string
    Latency   time.Duration
    Timestamp time.Time
}

type CrossChainAssetTransfer struct {
    TransferID    string
    AssetID       string
    SourceChainID string
    TargetChainID string
    Amount        float64
    Status        string
    Timestamp     time.Time
}

type CrossChainEscrow struct {
    EscrowID      string
    AssetID       string
    SourceChainID string
    TargetChainID string
    Amount        float64
    Status        string
    Timestamp     time.Time
}

type CrossChainAssetSwap struct {
    SwapID        string
    AssetID1      string
    ChainID1      string
    AssetID2      string
    ChainID2      string
    Amount1       float64
    Amount2       float64
    Status        string
    Timestamp     time.Time
}


type CrossChainVerification struct {
    RequestID      string
    ActivityID     string
    TargetChainID  string
    RequestDetails string
    ResponseDetails string
    Status         string
    Timestamp      time.Time
}


type DisputeEvent struct {
    DisputeID string
    EventType string
    Details   string
    Timestamp time.Time
}

type MediatorAssignment struct {
    DisputeID  string
    MediatorID string
    AssignedAt time.Time
}


type ChainStatus struct {
    ChainID    string
    Status     string
    LastUpdate time.Time
}

type EscrowEvent struct {
    TransactionID string
    EventType     string
    Details       string
    Timestamp     time.Time
}

type DataFeedEvent struct {
    FeedID    string
    EventType string
    Details   string
    Timestamp time.Time
}

type DataEvent struct {
    DataID    string
    EventType string
    Details   string
    Timestamp time.Time
}

type CrossChainActionRollback struct {
    ActionID   string
    Details    string
    Status     string
    Timestamp  time.Time
}

type CrossChainBalance struct {
    AssetID   string
    ChainID   string
    Balance   float64
    Timestamp time.Time
}

type CrossChainContract struct {
    ContractID   string
    IsValid      bool
    ValidationDetails string
    Timestamp    time.Time
}

type InterchainAgreement struct {
    AgreementID   string
    IsValid       bool
    ValidationDetails string
    Timestamp     time.Time
}


// AtomicSwap represents an atomic swap operation for cross-chain token exchange.
type AtomicSwap struct {
	SwapID         string       // Unique ID for the swap
	TokenA         SYN1200Token // Syn1200 token on Chain A
	TokenB         SYN1200Token // Syn1200 token on Chain B
	AmountA        float64      // Amount of Token A to swap
	AmountB        float64      // Amount of Token B to swap
	ChainAAddress  string       // Address on Chain A
	ChainBAddress  string       // Address on Chain B
	SecretHash     string       // Hash of the secret
	Secret         string       // The actual secret
	ExpirationTime time.Time    // Swap expiration time
	SwapInitiator  string       // Address initiating the swap
	SwapResponder  string       // Address responding to the swap
	Status         string       // Status of the swap (pending, completed, expired)
	mutex          sync.Mutex   // Mutex for thread-safe operations
	LedgerInstance *Ledger      // Ledger instance to track the swap operations
}

// AtomicSwapManager manages atomic swaps for cross-chain token exchanges
type AtomicSwapManager struct {
	ActiveSwaps    map[string]*AtomicSwap // Active swaps indexed by swap ID
	mutex          sync.Mutex             // Mutex for thread-safe operations
	LedgerInstance *Ledger                // Ledger instance for recording transactions
}

// BlockchainAgnosticProtocol represents the core protocol for cross-chain interaction
type BlockchainAgnosticProtocol struct {
	SupportedChains []string    // List of supported blockchain networks
	Validators      []Validator // Validators for cross-chain validation
	LedgerInstance  *Ledger     // Ledger instance for logging protocol activities
	mutex           sync.Mutex  // Mutex for thread-safe operations
}

// CrossChainTransaction represents a transaction processed across multiple blockchains
type CrossChainTransaction struct {
	TransactionID  string    // Unique transaction ID
	FromChain      string    // Originating chain
	ToChain        string    // Destination chain
	Amount         float64   // Amount being transferred
	TokenSymbol    string    // Token symbol being used
	FromAddress    string    // Sender's address
	ToAddress      string    // Recipient's address
	Timestamp      time.Time // Timestamp of the transaction
	ValidationHash string    // Validation hash for security
	Status         string    // Transaction status (pending, completed, failed)
}

// BlockchainAgnosticManager manages the cross-chain transactions using blockchain-agnostic protocols
type BlockchainAgnosticManager struct {
	ActiveTransactions map[string]*CrossChainTransaction // Active cross-chain transactions
	mutex              sync.Mutex                        // Mutex for thread-safe operations
	LedgerInstance     *Ledger                           // Ledger instance to track transactions
	Protocol           *BlockchainAgnosticProtocol       // Blockchain-agnostic protocol instance
}

// Bridge represents the structure for a cross-chain bridge
type Bridge struct {
	SupportedChains []string           // List of supported blockchain networks for the bridge
	Validators      []Validator        // Validators for cross-chain transactions
	LedgerInstance  *Ledger            // Ledger instance for logging bridge operations
	BridgeBalance   map[string]float64 // Bridge balance for each supported token
	mutex           sync.Mutex         // Mutex for thread-safe operations
}

// CrossChainTransfer represents a transfer processed by the bridge
type CrossChainTransfer struct {
	TransferID     string    // Unique transfer ID
	FromChain      string    // Originating blockchain network
	ToChain        string    // Destination blockchain network
	Amount         float64   // Amount being transferred
	TokenSymbol    string    // Token symbol being used
	FromAddress    string    // Sender's address
	ToAddress      string    // Recipient's address
	Timestamp      time.Time // Timestamp of the transfer
	Status         string    // Transfer status (pending, completed, failed)
	ValidationHash string    // Validation hash for security
}

// CrossChainMessage represents a message for communication between two blockchains
type CrossChainMessage struct {
	MessageID      string    // Unique message ID
	FromChain      string    // Originating blockchain network
	ToChain        string    // Destination blockchain network
	Payload        string    // The message payload (encrypted)
	Timestamp      time.Time // Timestamp of the message
	ValidationHash string    // Hash to validate the message's authenticity
	Status         string    // Message status (sent, received, confirmed)
}

// CrossChainCommunication represents the communication system between chains
type CrossChainCommunication struct {
	SupportedChains []string                     // List of supported blockchain networks
	Validators      []Validator                  // Validators for cross-chain message validation
	LedgerInstance  *Ledger                      // Ledger instance for logging cross-chain communications
	mutex           sync.Mutex                   // Mutex for thread-safe operations
	MessagePool     map[string]CrossChainMessage // Pool to store pending messages
}

// CrossChainSetup manages the configuration for cross-chain connections with other blockchains
type CrossChainSetup struct {
	Connections    map[string]string // Map of blockchain names to connection URLs
	LedgerInstance *Ledger           // Instance of the ledger for recording connections
	mutex          sync.Mutex        // Mutex for thread-safe operations
}

// CrossChainConnection handles establishing and managing connections between blockchains for cross-chain transactions
type CrossChainConnection struct {
	ConnectedChains []string       // List of connected blockchain networks
	LedgerInstance  *Ledger        // Instance of the ledger to track cross-chain transactions
	SubBlockPool    *SubBlockChain // Sub-block pool for transaction validation
	mutex           sync.Mutex     // Mutex for thread-safe operations
}

// OracleService represents an oracle that brings external data into the blockchain
type OracleService struct {
	DataSources    map[string]OracleDataSource // Data source name to URL
	LedgerInstance *Ledger                     // Ledger instance to record oracle data
	mutex          sync.Mutex                  // Mutex for thread-safe operations
}

// OracleDataSource represents a source from which the oracle fetches external data
type OracleDataSource struct {
	SourceID     string    // Unique ID for the data source
	Name         string    // Name of the data source
	URL          string    // The URL or API endpoint for fetching data
	Description  string    // A brief description of the data source
	IsActive     bool      // Whether the data source is currently active
	LastUpdated  time.Time // Timestamp of the last time the data source was updated
	DataFormat   string    // Format of the data (e.g., JSON, XML, CSV)
	AuthRequired bool      // Whether authentication is required to access the data
	ApiKey       string    // API key for accessing the data, if needed
	EncryptedKey string    // Encrypted version of the API key for security
}

// ************** Layer 2 Consensus Structs **************

// ConsensusMechanism represents a consensus mechanism with its configuration
type ConsensusMechanism struct {
	MechanismID     string    // Unique identifier for the consensus mechanism
	MechanismType   string    // Type of the consensus mechanism (e.g., "PoS", "PoH", "Synnergy")
	CurrentLoad     float64   // Current load on the consensus mechanism
	TransitionCount int       // Number of transitions to this mechanism
	Active          bool      // Whether the mechanism is currently active
	LastTransition  time.Time // Timestamp of the last transition to this mechanism
}

// CrossConsensusScalingManager manages transitions between different consensus mechanisms
type CrossConsensusScalingManager struct {
	ConsensusMechanisms map[string]*ConsensusMechanism // Available consensus mechanisms
	ActiveMechanism     *ConsensusMechanism            // Currently active consensus mechanism
	Ledger              *Ledger                        // Ledger instance for tracking consensus transitions
	EncryptionService   *Encryption                    // Encryption service for securing consensus-related data
	mu                  sync.Mutex                     // Mutex for concurrent management
}

// ConsensusStrategy defines the parameters of a consensus mechanism
type ConsensusStrategy struct {
	StrategyID   string    // Unique identifier for the strategy
	StrategyType string    // Type of consensus mechanism (e.g., "PoS", "PoH", "Synnergy")
	CurrentUsage float64   // Current resource usage (load) of the strategy
	LastHopped   time.Time // Last time this strategy was used
	Active       bool      // Whether this strategy is currently active
	HopCount     int       // Number of hops to this strategy
}

// DynamicConsensusManager manages dynamic consensus switching (hopping) between multiple strategies
type DynamicConsensusManager struct {
	Strategies        map[string]*ConsensusStrategy // All available consensus strategies
	ActiveStrategy    *ConsensusStrategy            // Currently active consensus strategy
	Ledger            *Ledger                       // Ledger instance for tracking consensus hops
	EncryptionService *Encryption                   // Encryption service for securing strategy-related data
	mu                sync.Mutex                    // Mutex for thread-safe management
}

// ConsensusLayer represents a consensus layer with adaptive properties
type ConsensusLayer struct {
	LayerID         string    // Unique identifier for the consensus layer
	LayerType       string    // Type of consensus layer (e.g., "PoS", "PoH", "Synnergy")
	CurrentLoad     float64   // Current load on the consensus layer
	MaxLoad         float64   // Maximum load before transitioning to a new layer
	TransitionTime  time.Time // Time of the last transition to this layer
	Active          bool      // Whether this layer is currently active
	TransitionCount int       // Number of transitions made to this layer
}

// ElasticConsensusManager manages the adaptive transitions between elastic consensus layers
type ElasticConsensusManager struct {
	ConsensusLayers   map[string]*ConsensusLayer // Available consensus layers
	ActiveLayer       *ConsensusLayer            // Currently active consensus layer
	Ledger            *Ledger                    // Ledger instance for tracking transitions and decisions
	EncryptionService *Encryption                // Encryption service for securing consensus data
	mu                sync.Mutex                 // Mutex for concurrent management
}

// CollaborationTask represents an off-chain computational task that requires collaboration
type CollaborationTask struct {
	TaskID            string    // Unique ID for the collaboration task
	AssignedNodes     []string  // List of nodes assigned to the task
	ComputationResult string    // Result of the off-chain computation
	CompletionStatus  string    // Status of the task ("Pending", "Completed")
	AssignedTime      time.Time // Time when the task was assigned
	CompletedTime     time.Time // Time when the task was completed
	EncryptedData     string    // Encrypted task details for security
}

// CollaborationNode represents a node that participates in Proof-of-Collaboration
type CollaborationNode struct {
	NodeID         string    // Unique identifier for the node
	Reputation     float64   // Reputation score based on successful collaboration tasks
	LastCollabTime time.Time // Time of the last collaboration
	Active         bool      // Whether the node is actively participating
}

// ProofOfCollaborationManager manages collaborative tasks and rewards in the PoCol network
type ProofOfCollaborationManager struct {
	Nodes             map[string]*CollaborationNode // Available nodes for collaboration
	ActiveTasks       map[string]*CollaborationTask // Active collaboration tasks
	Ledger            *Ledger                       // Ledger for recording PoCol actions
	EncryptionService *Encryption                   // Encryption service for securing collaboration data
	mu                sync.Mutex                    // Mutex for concurrent task management
}

// Layer2ConsensusLog tracks events related to Layer 2 consensus mechanisms, strategies, and layers
type Layer2ConsensusLog struct {
	EventType string    // Type of the consensus event (e.g., "StrategyAddition", "ConsensusTransition")
	Timestamp time.Time // Time the event occurred
	Details   string    // Details of the event
	Status    string    // Status of the event (e.g., "Added", "Transitioned", "Hopped")
}

// ************** Loanpool Structs **************

// LoanPool represents the structure of the main loan pool.
type LoanPool struct {
	mutex                 sync.Mutex         // For thread safety
	MainFund              *big.Int           // Main loanpool fund
	PersonalGrantFund     *big.Int           // 25%
	EcosystemGrantFund    *big.Int           // 25%
	EducationFund         *big.Int           // 5%
	HealthcareSupportFund *big.Int           // 5%
	PovertyFund           *big.Int           // 5%
	SecuredFund           *big.Int           // 15%
	BusinessGrantFund     *big.Int           // 25%
	UnsecuredLoanFund     *big.Int           // 15%
	EnvironmentalFund     *big.Int           // 5%
	Ledger                *Ledger            // Reference to the ledger for transaction logging
	Consensus             *SynnergyConsensus // Synnergy Consensus engine for validating fund transfers
	Encryption            *Encryption        // Encryption service for secure data handling
	EncryptedData         string             // Encrypted data for privacy and security

}

// LoanPoolManager provides functions to view the balances of the loan pool and its associated sub-funds.
type LoanPoolManager struct {
	mutex    sync.Mutex // For thread-safe operations
	LoanPool *LoanPool  // Reference to the LoanPool structure
	Ledger   *Ledger    // Reference to the ledger for transaction logging
}

// GrantApprovalProcess handles the two-stage approval process of a business personal grant proposal.
type BusinessPersonalGrantApprovalProcess struct {
	mutex             sync.Mutex                                        // Mutex for thread safety
	Ledger            *Ledger                                           // Reference to the ledger
	Consensus         *SynnergyConsensus                                // Synnergy Consensus engine
	Proposals         map[string]*BusinessPersonalGrantProposalApproval // Map to hold grant proposals by proposal ID
	AuthorityNodes    []AuthorityNodeTypes                              // List of valid authority node types (bank, government, central bank, etc.)
	PublicVotePeriod  time.Duration                                     // Time allowed for public voting
	AuthorityVoteTime time.Duration                                     // Time window for authority nodes to vote
}

// GrantProposalApproval represents a grant proposal along with its voting data
type BusinessPersonalGrantProposalApproval struct {
	Proposal          *BusinessPersonalGrantProposal // Reference to the grant proposal
	PublicVotes       map[string]bool                // Map of public votes (address -> vote)
	Stage             ApprovalStage                  // Current approval stage
	AuthorityVotes    map[string]bool                // Authority node votes
	VoteStartTime     time.Time                      // Time when voting starts
	ConfirmationCount int                            // Count of authority confirmations
	RejectionCount    int                            // Count of authority rejections
}

// BusinessPersonalGrantFund holds the details of the fund such as balance and distributed grants.
type BusinessPersonalGrantFund struct {
	mutex             sync.Mutex                                // Mutex for thread safety
	TotalBalance      *big.Int                                  // Total balance available in the fund
	GrantsDistributed *big.Int                                  // Total amount of grants distributed
	Ledger            *Ledger                                   // Reference to the ledger for storing transactions
	EncryptedData     string                                    // Encrypted data for privacy and security
	Proposals         map[string]*BusinessPersonalGrantProposal // Tracks loan proposals
}

// BusinessPersonalGrantDisbursementQueueEntry represents a proposal waiting for funds to be disbursed for the Business Personal Grant Fund.
type BusinessPersonalGrantDisbursementQueueEntry struct {
	ProposalID        string    // The proposal ID
	ProposerWallet    string    // Wallet address of the proposer
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
}

// BusinessPersonalGrantDisbursementManager manages the disbursement of funds for confirmed proposals in the Business Personal Grant Fund.
type BusinessPersonalGrantDisbursementManager struct {
	mutex             sync.Mutex                                     // Mutex for thread safety
	Ledger            *Ledger                                        // Reference to the ledger
	Consensus         *SynnergyConsensus                             // Synnergy Consensus engine
	FundBalance       float64                                        // Current balance of the Business Personal Grant Fund
	DisbursementQueue []*BusinessPersonalGrantDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime      time.Duration                                  // Maximum time a proposal can wait in the queue (30 days)
}

// GrantProposal represents the structure of the business grant application.
type BusinessPersonalGrantProposal struct {
	BusinessName        string            // Name of the business
	BusinessAddress     string            // Address of the business
	RegistrationNumber  string            // Business registration number
	Country             string            // Country of registration
	Website             string            // Business website (optional)
	BusinessActivities  string            // Description of business activities
	ApplicantName       string            // Name of the acting member applying for the funds
	WalletAddress       string            // The wallet address of the applicant
	AmountAppliedFor    float64           // Amount of grant funds being applied for
	UsageDescription    string            // Full description of how the funds will be used
	FinancialPosition   string            // Financial position or last submitted accounts (or state if it's a startup)
	SubmissionTimestamp time.Time         // Timestamp of proposal submission
	VerifiedBySyn900    bool              // Whether the proposal has been verified with syn900
	Status              string            // Proposal status (e.g., Pending, Approved, Rejected)
	Comments            []ProposalComment // Comments made on the proposal
	LastUpdated         time.Time         // Last update timestamp for the proposal
}

// ProposalManager manages the submission and verification of grant proposals.
type BusinessPersonalGrantProposalManager struct {
	mutex           sync.Mutex                                // Mutex for thread safety
	Ledger          *Ledger                                   // Reference to the ledger for storing proposals
	Proposals       map[string]*BusinessPersonalGrantProposal // Map of proposals by business name
	Syn900Validator *Syn900Validator                          // Reference to syn900 validator for wallet verification
	Encryption      *Encryption                               // Encryption service for secure proposal data
}

// EcosystemGrantApprovalProcess handles the two-stage approval process for an ecosystem grant proposal.
type EcosystemGrantApprovalProcess struct {
	mutex             sync.Mutex                                 // Mutex for thread safety
	Ledger            *Ledger                                    // Reference to the ledger
	Consensus         *SynnergyConsensus                         // Synnergy Consensus engine
	Proposals         map[string]*EcosystemGrantProposalApproval // Map to hold grant proposals by proposal ID
	AuthorityNodes    []AuthorityNodeTypes                       // List of valid authority node types (bank, government, central bank, etc.)
	PublicVotePeriod  time.Duration                              // Time allowed for public voting
	AuthorityVoteTime time.Duration                              // Time window for authority nodes to vote
}

// GrantProposalApproval represents a grant proposal along with its voting data.
type EcosystemGrantProposalApproval struct {
	Proposal          *EcosystemGrantProposal // Reference to the grant proposal
	PublicVotes       map[string]bool         // Map of public votes (address -> vote)
	Stage             ApprovalStage           // Current approval stage
	AuthorityVotes    map[string]bool         // Authority node votes
	VoteStartTime     time.Time               // Time when voting starts
	ConfirmationCount int                     // Count of authority confirmations
	RejectionCount    int                     // Count of authority rejections
}

// EcosystemGrantFund holds the details of the fund such as balance and distributed grants.
type EcosystemGrantFund struct {
	mutex             sync.Mutex                         // Mutex for thread safety
	TotalBalance      *big.Int                           // Total balance available in the fund
	GrantsDistributed *big.Int                           // Total amount of grants distributed
	Ledger            *Ledger                            // Reference to the ledger for storing transactions
	EncryptedData     string                             // Encrypted data for privacy and security
	Proposals         map[string]*EcosystemGrantProposal // Tracks loan proposals
}

// EcosystemGrantDisbursementQueueEntry represents a proposal waiting for funds to be disbursed for the Ecosystem Grant Fund.
type EcosystemGrantDisbursementQueueEntry struct {
	ProposalID        string    // The proposal ID
	ProposerWallet    string    // Wallet address of the proposer
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
}

// EcosystemGrantDisbursementManager manages the disbursement of funds for confirmed proposals in the Ecosystem Grant Fund.
type EcosystemGrantDisbursementManager struct {
	mutex             sync.Mutex                              // Mutex for thread safety
	Ledger            *Ledger                                 // Reference to the ledger
	Consensus         *SynnergyConsensus                      // Synnergy Consensus engine
	FundBalance       float64                                 // Current balance of the Ecosystem Grant Fund
	DisbursementQueue []*EcosystemGrantDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime      time.Duration                           // Maximum time a proposal can wait in the queue (30 days)
}

// EcosystemGrantProposal represents the structure of the ecosystem grant application.
type EcosystemGrantProposal struct {
	BusinessName         string            // Name of the business
	BusinessAddress      string            // Address of the business
	RegistrationNumber   string            // Business registration number
	Country              string            // Country of registration
	Website              string            // Business website (optional)
	BusinessActivities   string            // Description of business activities
	ApplicantName        string            // Name of the acting member applying for the funds
	WalletAddress        string            // The wallet address of the applicant
	AmountAppliedFor     float64           // Amount of grant funds being applied for
	UsageDescription     string            // Full description of how the funds will be used
	EcosystemApplication string            // Specific description of how the funds will be used within the ecosystem
	FinancialPosition    string            // Financial position or last submitted accounts (or state if it's a startup)
	SubmissionTimestamp  time.Time         // Timestamp of proposal submission
	VerifiedBySyn900     bool              // Whether the proposal has been verified with syn900
	Status               string            // Proposal status (e.g., Pending, Approved, Rejected)
	Comments             []ProposalComment // Comments made on the proposal
	LastUpdated          time.Time         // Last update timestamp for the proposal
}

// ProposalManager manages the submission and verification of ecosystem grant proposals.
type EcosystemGrantProposalManager struct {
	mutex           sync.Mutex                         // Mutex for thread safety
	Ledger          *Ledger                            // Reference to the ledger for storing proposals
	Proposals       map[string]*EcosystemGrantProposal // Map of proposals by business name
	Syn900Validator *tokenledgers.Syn900Validator      // Reference to syn900 validator for wallet verification
	Encryption      *Encryption                        // Encryption service for secure proposal data
}

// EducationFundApprovalProcess manages the approval process for education fund proposals.
type EducationFundApprovalProcess struct {
	mutex             sync.Mutex
	Ledger            *Ledger                                 // Ledger to store proposal and approval status
	Nodes             []*AuthorityNode                        // List of all authority nodes in the network
	ActiveProposals   map[string]*EducationFundActiveProposal // Map of active proposals being reviewed
	EncryptionService *Encryption                             // Encryption service for secure transmission
	NetworkManager    *NetworkManager                         // Network manager to handle proposal transmission
	RequeueDuration   time.Duration                           // Duration before a proposal is requeued
	MaxConfirmations  int                                     // Required confirmations for proposal approval
	MaxRejections     int                                     // Required rejections for proposal rejection
}

// ActiveProposal keeps track of confirmations, rejections, and assigned nodes for a proposal.
type EducationFundActiveProposal struct {
	ProposalID       string                    // Unique proposal ID
	ProposalData     *EducationFundProposal    // The education fund proposal details
	ConfirmedNodes   map[string]bool           // Nodes that confirmed the proposal
	RejectedNodes    map[string]bool           // Nodes that rejected the proposal
	AssignedNodes    map[string]*AuthorityNode // Nodes currently assigned for review
	Status           string                    // Status of the proposal (Pending, Approved, Rejected)
	LastDistribution time.Time                 // Timestamp of last node distribution
	ProposalDeadline time.Time                 // Deadline for the proposal before requeuing
}

// EducationFund holds the details of the fund, including balance and distributed grants.
type EducationFund struct {
	mutex             sync.Mutex                         // Mutex for thread safety
	TotalBalance      *big.Int                           // Total balance available in the fund
	GrantsDistributed *big.Int                           // Total amount of grants distributed
	Ledger            *Ledger                            // Reference to the ledger for storing transactions
	Consensus         *SynnergyConsensus                 // Synnergy Consensus engine for transaction validation
	EncryptionService *Encryption                        // Encryption service for securing sensitive data
	EncryptedData     string                             // Encrypted data for privacy and security
	Proposals         map[string]*EcosystemGrantProposal // Tracks loan proposals
}

// EducationFundDisbursementQueueEntry represents a proposal waiting for funds to be disbursed for the Education Fund.
type EducationFundDisbursementQueueEntry struct {
	ProposalID        string    // The proposal ID
	ProposerWallet    string    // Wallet address of the proposer
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
}

// EducationFundDisbursementManager manages the disbursement of funds for confirmed proposals in the Education Fund.
type EducationFundDisbursementManager struct {
	mutex             sync.Mutex                             // Mutex for thread safety
	Ledger            *Ledger                                // Reference to the ledger
	Consensus         *SynnergyConsensus                     // Synnergy Consensus engine
	FundBalance       float64                                // Current balance of the Education Fund
	DisbursementQueue []*EducationFundDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime      time.Duration                          // Maximum time a proposal can wait in the queue (30 days)
	EncryptionService *Encryption                            // Encryption service for secure data
}

// EducationFundProposal represents the structure of the education fund proposal.
type EducationFundProposal struct {
	ApplicantName       string            // Name of the applicant
	ApplicantContact    string            // Applicant's contact information
	WalletAddress       string            // The wallet address of the applicant
	InstitutionName     string            // Name of the educational institution
	CourseName          string            // Name of the course
	CourseLevel         string            // Level of the course (e.g., Bachelor's, Master's)
	ApplicationEvidence string            // Evidence of course application or acceptance
	PersonalStatement   string            // Personal statement of the applicant
	AmountAppliedFor    float64           // Amount of funds being applied for
	SubmissionTimestamp time.Time         // Timestamp of proposal submission
	SponsorName         string            // Name of the sponsor (if applicable)
	SponsorContactInfo  string            // Contact information of the sponsor (if applicable)
	VerifiedBySyn900    bool              // Whether the proposal has been verified with syn900
	Status              string            // Proposal status (e.g., Pending, Approved, Rejected)
	Comments            []ProposalComment // Comments made on the proposal
	LastUpdated         time.Time         // Last update timestamp for the proposal
}

// ProposalManager manages the submission and verification of education fund proposals.
type EducationFundProposalManager struct {
	mutex           sync.Mutex                        // Mutex for thread safety
	Ledger          *Ledger                           // Reference to the ledger for storing proposals
	Proposals       map[string]*EducationFundProposal // Map of proposals by applicant name
	Syn900Validator *tokenledgers.Syn900Validator     // Reference to syn900 validator for wallet verification
	Encryption      *Encryption                       // Encryption service for secure proposal data
}

// HealthcareSupportFundApprovalProcess manages the approval process for healthcare fund proposals.
type HealthcareSupportFundApprovalProcess struct {
	mutex             sync.Mutex
	Ledger            *Ledger                                         // Ledger to store proposal and approval status
	Nodes             []*AuthorityNode                                // List of all authority nodes in the network
	ActiveProposals   map[string]*HealthcareSupportFundActiveProposal // Map of active proposals being reviewed
	EncryptionService *Encryption                                     // Encryption service for secure transmission
	NetworkManager    *NetworkManager                                 // Network manager to handle proposal transmission
	RequeueDuration   time.Duration                                   // Duration before a proposal is requeued (7 days)
	MaxConfirmations  int                                             // Required confirmations for proposal approval
	MaxRejections     int                                             // Required rejections for proposal rejection
}

// ActiveProposal keeps track of confirmations, rejections, and assigned nodes for a proposal.
type HealthcareSupportFundActiveProposal struct {
	ProposalID       string                         // Unique proposal ID
	ProposalData     *HealthcareSupportFundProposal // The healthcare support fund proposal details
	ConfirmedNodes   map[string]bool                // Nodes that confirmed the proposal
	RejectedNodes    map[string]bool                // Nodes that rejected the proposal
	AssignedNodes    map[string]*AuthorityNode      // Nodes currently assigned for review
	Status           string                         // Status of the proposal (Pending, Approved, Rejected)
	LastDistribution time.Time                      // Timestamp of last node distribution
	ProposalDeadline time.Time                      // Deadline for the proposal before requeuing
}

// HealthcareSupportFundDisbursementQueueEntry represents a proposal waiting for funds to be disbursed for the Healthcare Support Fund.
type HealthcareSupportFundDisbursementQueueEntry struct {
	ProposalID        string    // The proposal ID
	ProposerWallet    string    // Wallet address of the proposer
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
}

// HealthcareSupportFundDisbursementManager manages the disbursement of funds for confirmed proposals in the Healthcare Support Fund.
type HealthcareSupportFundDisbursementManager struct {
	mutex             sync.Mutex                                     // Mutex for thread safety
	Ledger            *Ledger                                        // Reference to the ledger
	Consensus         *SynnergyConsensus                             // Synnergy Consensus engine
	FundBalance       float64                                        // Current balance of the Healthcare Support Fund
	DisbursementQueue []*HealthcareSupportFundDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime      time.Duration                                  // Maximum time a proposal can wait in the queue (7 days)
	EncryptionService *Encryption                                    // Encryption service for secure data
}

// HealthcareSupportFund manages the details of the fund, including balance and distributed healthcare grants.
type HealthcareSupportFund struct {
	mutex             sync.Mutex                                // Mutex for thread safety
	TotalBalance      *big.Int                                  // Total balance available in the fund
	GrantsDistributed *big.Int                                  // Total amount of healthcare grants distributed
	Ledger            *Ledger                                   // Reference to the ledger for storing transactions
	Consensus         *SynnergyConsensus                        // Synnergy Consensus engine for transaction validation
	EncryptionService *Encryption                               // Encryption service for securing sensitive data
	EncryptedData     string                                    // Encrypted data for privacy and security
	Proposals         map[string]*HealthcareSupportFundProposal // Tracks loan proposals
}

// HealthcareSupportFundProposal represents the structure of the healthcare support fund proposal.
type HealthcareSupportFundProposal struct {
	ApplicantName              string            // Name of the applicant (person requiring medical treatment)
	ApplicantContact           string            // Applicant's contact information
	MedicalProfessionalName    string            // Name of the supporting medical professional
	MedicalProfessionalContact string            // Contact information for the medical professional
	WalletAddress              string            // The wallet address of the applicant
	HospitalName               string            // Name of the hospital, medical practice, or provider
	MedicalProcedure           string            // Details of the medical procedure required
	CostBreakdownEvidence      string            // Evidence of the cost breakdown of the medical treatment
	HospitalAddress            string            // Full address of the hospital or medical provider
	HospitalContactInfo        string            // Contact information for the hospital or provider
	AmountAppliedFor           float64           // Amount of funds being applied for
	SubmissionTimestamp        time.Time         // Timestamp of proposal submission
	VerifiedBySyn900           bool              // Whether the proposal has been verified with syn900
	Status                     string            // Proposal status (e.g., Pending, Approved, Rejected)
	Comments                   []ProposalComment // Comments made on the proposal
	LastUpdated                time.Time         // Last update timestamp for the proposal
}

// ProposalManager manages the submission and verification of healthcare fund proposals.
type HealthcareSupportFundProposalManager struct {
	mutex           sync.Mutex                                // Mutex for thread safety
	Ledger          *Ledger                                   // Reference to the ledger for storing proposals
	Proposals       map[string]*HealthcareSupportFundProposal // Map of proposals by applicant name
	Syn900Validator *tokenledgers.Syn900Validator             // Reference to syn900 validator for wallet verification
	Encryption      *Encryption                               // Encryption service for secure proposal data
}

// PovertyFundApprovalProcess manages the approval process for poverty fund proposals.
type PovertyFundApprovalProcess struct {
	mutex             sync.Mutex
	Ledger            *Ledger                               // Ledger to store proposal and approval status
	Nodes             []*AuthorityNode                      // List of all authority nodes in the network
	ActiveProposals   map[string]*PovertyFundActiveProposal // Map of active proposals being reviewed
	EncryptionService *Encryption                           // Encryption service for secure transmission
	NetworkManager    *NetworkManager                       // Network manager to handle proposal transmission
	RequeueDuration   time.Duration                         // Duration before a proposal is requeued (7 days)
	MaxConfirmations  int                                   // Required confirmations for proposal approval
	MaxRejections     int                                   // Required rejections for proposal rejection
}

// EventRecord represents a log entry for an event
type EventRecord struct {
	EventType   string
	Description string // Add this line if Description is missing
	Timestamp   time.Time
}

// Define LedgerEvent struct in ledger package if it doesn’t exist
type LedgerEvent struct {
	Timestamp time.Time
	EventType string
	Details   string
}

// ActiveProposal keeps track of confirmations, rejections, and assigned nodes for a proposal.
type PovertyFundActiveProposal struct {
	ProposalID       string                         // Unique proposal ID
	ProposalData     *PovertyFundProposal           // The poverty fund proposal details
	ConfirmedNodes   map[string]bool                // Nodes that confirmed the proposal
	RejectedNodes    map[string]bool                // Nodes that rejected the proposal
	AssignedNodes    map[string]*AuthorityNodeTypes // Nodes currently assigned for review
	Status           string                         // Status of the proposal (Pending, Approved, Rejected)
	LastDistribution time.Time                      // Timestamp of last node distribution
	ProposalDeadline time.Time                      // Deadline for the proposal before requeuing
}

// PovertyFundDisbursementQueueEntry represents a proposal waiting for funds to be disbursed for the Poverty Fund.
type PovertyFundDisbursementQueueEntry struct {
	ProposalID        string    // The proposal ID
	ProposerWallet    string    // Wallet address of the proposer
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
}

// PovertyFundDisbursementManager manages the disbursement of funds for confirmed proposals in the Poverty Fund.
type PovertyFundDisbursementManager struct {
	mutex             sync.Mutex                           // Mutex for thread safety
	Ledger            *Ledger                              // Reference to the ledger
	Consensus         *SynnergyConsensus                   // Synnergy Consensus engine
	FundBalance       float64                              // Current balance of the Poverty Fund
	DisbursementQueue []*PovertyFundDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime      time.Duration                        // Maximum time a proposal can wait in the queue (48 hours)
	EncryptionService *Encryption                          // Encryption service for secure data
}

// PovertyFund manages the details of the fund, including balance and distributed grants.
type PovertyFund struct {
	mutex             sync.Mutex                      // Mutex for thread safety
	TotalBalance      *big.Int                        // Total balance available in the fund
	GrantsDistributed *big.Int                        // Total amount of grants distributed
	Ledger            *Ledger                         // Reference to the ledger for storing transactions
	Consensus         *SynnergyConsensus              // Synnergy Consensus engine for transaction validation
	EncryptionService *Encryption                     // Encryption service for securing sensitive data
	EncryptedData     string                          // Encrypted data for privacy and security
	Proposals         map[string]*PovertyFundProposal // Tracks loan proposals
}

// PovertyFundProposal represents the structure of the poverty fund proposal.
type PovertyFundProposal struct {
	ApplicantName       string            // Name of the applicant
	ApplicantContact    string            // Applicant's contact information
	IncomeDetails       string            // Income details of the applicant
	BankBalanceDetails  string            // Current bank balance details of the applicant
	IncomeEvidence      []byte            // Attachment: Evidence of income (encrypted)
	BankBalanceEvidence []byte            // Attachment: Evidence of bank balance (encrypted)
	StatementOfReason   string            // Statement of reason for the request
	BenefitStatus       string            // Current benefit status of the applicant
	WalletAddress       string            // The wallet address of the applicant
	AmountAppliedFor    float64           // Amount of funds being applied for
	SubmissionTimestamp time.Time         // Timestamp of proposal submission
	VerifiedBySyn900    bool              // Whether the proposal has been verified with syn900
	Status              string            // Proposal status (e.g., Pending, Approved, Rejected)
	Comments            []ProposalComment // Comments made on the proposal
	LastUpdated         time.Time         // Last update timestamp for the proposal
}

// ProposalComment represents a comment added to a proposal.
type ProposalComment struct {
	CommentID string    // Unique ID for the comment
	Commenter string    // Name or wallet address of the commenter
	Comment   string    // The content of the comment
	CreatedAt time.Time // Timestamp of when the comment was added
}

// ProposalManager manages the submission and verification of poverty fund proposals.
type PovertyFundProposalManager struct {
	mutex           sync.Mutex                      // Mutex for thread safety
	Ledger          *Ledger                         // Reference to the ledger for storing proposals
	Proposals       map[string]*PovertyFundProposal // Map of proposals by applicant name
	Syn900Validator *tokenledgers.Syn900Validator   // Reference to syn900 validator for wallet verification
	Encryption      *Encryption                     // Encryption service for secure proposal data
}

// SecuredLoanApprovalProcess manages the approval process for secured loan proposals.
type SecuredLoanApprovalProcess struct {
	mutex             sync.Mutex
	Ledger            *Ledger                    // Ledger to store proposal and approval status
	Nodes             []*AuthorityNode           // List of all authority nodes in the network
	ActiveProposals   map[string]*ActiveProposal // Map of active proposals being reviewed
	EncryptionService *Encryption                // Encryption service for secure transmission
	NetworkManager    *NetworkManager            // Network manager to handle proposal transmission
	RequeueDuration   time.Duration              // Duration before a proposal is requeued (48 hours)
	MaxConfirmations  int                        // Required confirmations for proposal approval
	MaxRejections     int                        // Required rejections for proposal rejection
}

// ActiveProposal keeps track of confirmations, rejections, interest rates, and assigned nodes for a proposal.
type SecuredLoanActiveProposal struct {
	ProposalID       string                         // Unique proposal ID
	ProposalData     *SecuredLoanProposal           // The loan proposal details
	ConfirmedNodes   map[string]bool                // Nodes that confirmed the proposal
	RejectedNodes    map[string]bool                // Nodes that rejected the proposal
	AssignedNodes    map[string]*AuthorityNodeTypes // Nodes currently assigned for review
	Status           string                         // Status of the proposal (Pending, Approved, Rejected)
	InterestRates    []float64                      // List of interest rates submitted by authority nodes
	LastDistribution time.Time                      // Timestamp of last node distribution
	ProposalDeadline time.Time                      // Deadline for the proposal before requeuing
	AverageInterest  float64                        // Running average of the interest rates
	AllDocsOpened    bool                           // Whether all nodes have opened the attached documents
}

// Collateral represents an asset that is used as security for a loan.
type Collateral struct {
	CollateralID     string    // Unique identifier for the collateral
	OwnerID          string    // Owner's identifier (e.g., wallet address)
	AssetType        string    // Type of asset (e.g., property, cryptocurrency, etc.)
	AssetValue       float64   // Value of the asset in monetary terms
	LockedValue      float64   // Amount locked for the loan
	LoanID           string    // Associated loan ID
	Status           string    // Status of the collateral (e.g., "Locked", "Released")
	SubmissionDate   time.Time // Date when the collateral was submitted
	VerificationDate time.Time // Date when the collateral was verified
	IsVerified       bool      // Whether the collateral has been verified
}

// CollateralSubmission represents the structure for collateral proof submission.
type CollateralSubmission struct {
	LoanID           string    // Unique loan ID linked to the collateral
	ProposerWallet   string    // Wallet address of the proposer (borrower)
	CollateralType   string    // Type of collateral (e.g., property, car, assets)
	CollateralValue  float64   // Value of the collateral being offered
	CollateralProof  []byte    // Digital document providing proof of collateral (e.g., title deed)
	IOULegalDocument []byte    // Digital IOU or legal agreement document
	SubmissionTime   time.Time // Time of collateral submission
	ApprovalStatus   string    // Status of the collateral (Pending, Approved, Rejected)
	ApprovedBy       string    // Authority node that approved/rejected the collateral
	ApprovalTime     time.Time // Time of approval or rejection
}

// CollateralManager manages the submission and approval process for collateral in a secured loan.
type CollateralManager struct {
	mutex             sync.Mutex                       // Mutex for thread safety
	Ledger            *Ledger                          // Ledger to record collateral submissions and approvals
	Consensus         *SynnergyConsensus               // Synnergy Consensus engine for validating collateral
	EncryptionService *Encryption                      // Encryption service for securing sensitive collateral data
	Submissions       map[string]*CollateralSubmission // Map of collateral submissions by loan ID
	ApprovalQueue     []*CollateralSubmission          // Queue of submissions pending approval
}

// LoanTerms represents the structure of customized loan terms.
type LoanTerms struct {
	RepaymentLength      int     // Number of months to repay the loan
	AmountBorrowed       float64 // Total loan amount
	InterestRate         float64 // Interest rate applied (unless Islamic terms)
	IslamicFinance       bool    // If true, switches to Islamic terms (no interest, fee applied)
	FeeOnTop             float64 // Flat fee applied if Islamic finance is selected
	TotalRepaymentAmount float64 // Total amount to be repaid (calculated)
}

// SecuredLoanTermManager manages customization of secured loan terms.
type SecuredLoanTermManager struct {
	mutex             sync.Mutex            // Mutex for thread safety
	Ledger            *Ledger               // Reference to the ledger for storing loan term records
	Consensus         *SynnergyConsensus    // Synnergy Consensus engine for validating terms
	LoanTermRecords   map[string]*LoanTerms // Records of loan terms by loan ID
	EncryptionService *Encryption           // Encryption service for securing sensitive data
}

// SecuredLoanDisbursementQueueEntry represents a loan waiting for disbursement in the secured loan pool.
type SecuredLoanDisbursementQueueEntry struct {
	ProposalID        string    // The loan proposal ID
	ProposerWallet    string    // Wallet address of the borrower
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
	AverageInterest   float64   // Average interest rate to be applied to the loan
}

// SecuredLoanDisbursementManager manages the disbursement of approved loans in the secured loan pool.
type SecuredLoanDisbursementManager struct {
	mutex               sync.Mutex                           // Mutex for thread safety
	Ledger              *Ledger                              // Reference to the ledger for recording disbursements
	Consensus           *SynnergyConsensus                   // Synnergy Consensus engine
	FundBalance         float64                              // Current balance of the loan pool
	DisbursementQueue   []*SecuredLoanDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime        time.Duration                        // Maximum time a proposal can wait in the queue
	EncryptionService   *Encryption                          // Encryption service for securing sensitive data
	IssuerFeePercentage float64                              // Issuer fee (0.5%)
}

// SecuredLoanProposal represents the full loan application process.
type SecuredLoanProposal struct {
	LoanID              string             // Unique loan ID
	ApplicantName       string             // Name of the applicant
	ApplicantID         string             // Unique applicant ID (validated by syn900)
	WalletAddress       string             // Wallet address of the applicant
	SubmissionTimestamp time.Time          // Time when the proposal was submitted
	ProposalStatus      string             // Status of the proposal (e.g., Pending, Approved, Rejected)
	LastUpdated         time.Time          // Timestamp of last update
	ApprovalStage       string             // Current stage of the approval process (e.g., Application, CreditCheck, Affordability, Collateral, Terms)
	CreditScore         float64            // Applicant's credit score (from decentralized credit check)
	AffordabilityStatus string             // Result of affordability check (e.g., Approved, Rejected)
	CollateralStatus    string             // Result of collateral submission
	TermsCustomization  bool               // Whether terms customization has been completed
	Status              string             // Status of the proposal ("Pending", "Approved", "Rejected")
	Repayments          []Repayment        // List of repayments made by the borrower
	LoanTerms           LoanTerms          // Terms of the loan (interest rate, duration, etc.)
	InterestPayments    []InterestPayment  // Interest payments made by the borrower
	AffordabilityCheck  AffordabilityCheck // Results of the affordability check for the borrower
	FundsDisbursed      float64            // Amount of funds disbursed for the loan
	Collateral          []Collateral       // Slice of Collateral structs to store multiple collateral records

}

// Repayment represents a repayment made on a loan
type Repayment struct {
	Amount    float64   // Amount repaid
	Date      time.Time // Date of the repayment
	Remaining float64   // Remaining balance after the repayment
}

// InterestPayment represents an interest payment made on a loan
type InterestPayment struct {
	Amount float64   // Amount of interest paid
	Date   time.Time // Date of the interest payment
}

// ProposalManager handles the overall proposal process for secured loans.
type SecuredLoanProposalManager struct {
	mutex             sync.Mutex                      // Mutex for thread safety
	Ledger            *Ledger                         // Ledger for storing proposal data
	Consensus         *SynnergyConsensus              // Consensus engine for validation
	Syn900Validator   *Syn900Validator                // Syn900 validator for ID validation
	EncryptionService *Encryption                     // Encryption service for secure proposal data
	CreditChecker     *CreditCheckManager             // Decentralized Credit Check manager
	AffordabilityMgr  *AffordabilityManager           // Affordability Check manager
	CollateralMgr     *CollateralManager              // Collateral submission manager
	TermsManager      *TermsCustomizationManager      // Customization of loan terms
	Proposals         map[string]*SecuredLoanProposal // Proposals mapped by LoanID
}

// LoanRepaymentDetails stores repayment information for a loan.
type SecuredLoanRepaymentDetails struct {
	LoanID            string      // Unique loan ID
	ProposerWallet    string      // Wallet address of the borrower
	TotalAmount       float64     // Total amount to be repaid
	RemainingAmount   float64     // Remaining amount to be repaid
	InterestRate      float64     // Interest rate applied to the loan
	RepaymentDates    []time.Time // Scheduled payment dates
	NextPaymentDue    time.Time   // Next payment due date
	Status            string      // Loan status (Active, Defaulted, Satisfied)
	DefaultedAt       *time.Time  // If loan defaulted, record the default date
	CollateralContact string      // Contact email for collateral request
	AuthorityWallets  []string    // Wallet addresses of authority nodes
}

// SecuredLoanRepaymentManager manages the repayment and settlement process for secured loans.
type SecuredLoanRepaymentManager struct {
	mutex             sync.Mutex                              // Mutex for thread safety
	Ledger            *Ledger                                 // Reference to the ledger
	EncryptionService *Encryption                             // Encryption service for secure data
	Syn900Registry    *SYN900Registry                         // Reference to Syn900 for record keeping
	LoanRepayments    map[string]*SecuredLoanRepaymentDetails // Map of loan repayments by loan ID
	DefaultThreshold  time.Duration                           // Time duration for default (e.g., 6 months)
}

// SecuredLoanManagement handles management tasks such as authority node updates, borrower detail changes, and term change requests.
type SecuredLoanManagement struct {
	mutex              sync.Mutex
	Ledger             *Ledger                               // Ledger reference for recording updates
	EncryptionService  *Encryption                           // Encryption service for data security
	ConsensusEngine    *SynnergyConsensus                    // Consensus engine for approval processes
	NetworkManager     *NetworkManager                       // Network manager for sending requests
	AuthorityNodes     map[string]*AuthorityNodeTypes        // Map of authority nodes that can manage loans
	LoanBorrowerInfo   map[string]*BorrowerDetails           // Stores borrower details by loan ID
	TermChangeRequests map[string]*BorrowerTermChangeRequest // Stores term change requests by loan ID
}

// BorrowerTermChangeRequest represents a request for changing loan terms.
type BorrowerTermChangeRequest struct {
	LoanID           string
	RequestedTerms   string                         // New terms requested by the borrower
	ApprovalStatus   string                         // Approval status (Pending, Accepted, Rejected)
	ConfirmedNodes   map[string]bool                // Nodes that confirmed the request
	RejectedNodes    map[string]bool                // Nodes that rejected the request
	AssignedNodes    map[string]*AuthorityNodeTypes // Nodes currently assigned for review
	LastDistribution time.Time                      // Last node distribution time
	RequeueDeadline  time.Time                      // Requeue if not processed within this time
}

// SecuredLoanPool manages the details of the fund, including balance, loan disbursements, repayments, and defaults.
type SecuredLoanPool struct {
	mutex             sync.Mutex                      // Mutex for thread safety
	TotalBalance      *big.Int                        // Total balance available in the loan pool
	LoansDistributed  *big.Int                        // Total amount of loans distributed
	LoansRepaid       *big.Int                        // Total amount of loans repaid
	LoansDefaulted    *big.Int                        // Total amount of loans that have defaulted
	Ledger            *Ledger                         // Reference to the ledger for storing loan transactions
	Consensus         *SynnergyConsensus              // Synnergy Consensus engine for transaction validation
	EncryptionService *Encryption                     // Encryption service for securing sensitive loan data
	LoanRecords       map[string]*LoanRecord          // Map of loan records by applicant wallet
	EncryptedData     string                          // Encrypted data for privacy and security
	Proposals         map[string]*SecuredLoanProposal // Tracks loan proposals
	FundsAvailable    float64                         // Funds currently available for new loans

}

// SmallBusinessGrantApprovalProcess handles the two-stage approval process of a small business grant proposal.
type SmallBusinessGrantApprovalProcess struct {
	mutex             sync.Mutex                                     // Mutex for thread safety
	Ledger            *Ledger                                        // Reference to the ledger
	Consensus         *SynnergyConsensus                             // Synnergy Consensus engine
	Proposals         map[string]*SmallBusinessGrantProposalApproval // Map to hold grant proposals by proposal ID
	AuthorityNodes    []AuthorityNodeTypes                           // List of valid authority node types (bank, government, central bank, etc.)
	PublicVotePeriod  time.Duration                                  // Time allowed for public voting
	AuthorityVoteTime time.Duration                                  // Time window for authority nodes to vote
}

// SmallBusinessGrantProposalApproval represents a grant proposal along with its voting data.
type SmallBusinessGrantProposalApproval struct {
	Proposal          *SmallBusinessGrantProposal // Reference to the grant proposal
	PublicVotes       map[string]bool             // Map of public votes (address -> vote)
	Stage             ApprovalStage               // Current approval stage
	AuthorityVotes    map[string]bool             // Authority node votes
	VoteStartTime     time.Time                   // Time when voting starts
	ConfirmationCount int                         // Count of authority confirmations
	RejectionCount    int                         // Count of authority rejections
}

// SmallBusinessGrantDisbursementQueueEntry represents a proposal waiting for funds to be disbursed for the Small Business Grant Fund.
type SmallBusinessGrantDisbursementQueueEntry struct {
	ProposalID        string    // The proposal ID
	ProposerWallet    string    // Wallet address of the proposer
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
}

// SmallBusinessGrantDisbursementManager manages the disbursement of funds for confirmed proposals in the Small Business Grant Fund.
type SmallBusinessGrantDisbursementManager struct {
	mutex             sync.Mutex                                  // Mutex for thread safety
	Ledger            *Ledger                                     // Reference to the ledger
	Consensus         *SynnergyConsensus                          // Synnergy Consensus engine
	FundBalance       float64                                     // Current balance of the Small Business Grant Fund
	DisbursementQueue []*SmallBusinessGrantDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime      time.Duration                               // Maximum time a proposal can wait in the queue (30 days)
}

// SmallBusinessGrantProposal represents the structure of the small business grant application.
type SmallBusinessGrantProposal struct {
	BusinessName        string            // Name of the business
	BusinessAddress     string            // Address of the business
	RegistrationNumber  string            // Business registration number
	Country             string            // Country of registration
	Website             string            // Business website (optional)
	BusinessActivities  string            // Description of business activities
	ApplicantName       string            // Name of the acting member applying for the funds
	WalletAddress       string            // The wallet address of the applicant
	AmountAppliedFor    float64           // Amount of grant funds being applied for
	UsageDescription    string            // Full description of how the funds will be used
	FinancialPosition   string            // Financial position or last submitted accounts (or state if it's a startup)
	SubmissionTimestamp time.Time         // Timestamp of proposal submission
	VerifiedBySyn900    bool              // Whether the proposal has been verified with syn900
	Status              string            // Proposal status (e.g., Pending, Approved, Rejected)
	Comments            []ProposalComment // Comments made on the proposal
	LastUpdated         time.Time         // Last update timestamp for the proposal
	Startup             bool              // Is the business a startup?
	EmployeeCount       int               // Number of employees (required if not a startup)
}

// SmallBusinessProposalManager manages the submission and verification of small business grant proposals.
type SmallBusinessGrantProposalManager struct {
	mutex           sync.Mutex                             // Mutex for thread safety
	Ledger          *Ledger                                // Reference to the ledger for storing proposals
	Proposals       map[string]*SmallBusinessGrantProposal // Map of proposals by business name
	Syn900Validator *tokenledgers.Syn900Validator          // Reference to syn900 validator for wallet verification
	Encryption      *Encryption                            // Encryption service for secure proposal data
}

// SmallBusinessGrantFund holds the details of the fund such as balance and distributed grants.
type SmallBusinessGrantFund struct {
	mutex             sync.Mutex                             // Mutex for thread safety
	TotalBalance      *big.Int                               // Total balance available in the fund
	GrantsDistributed *big.Int                               // Total amount of grants distributed
	Ledger            *Ledger                                // Reference to the ledger for storing transactions
	EncryptedData     string                                 // Encrypted data for privacy and security
	Proposals         map[string]*SmallBusinessGrantProposal // Tracks loan proposals
}

// AffordabilityCheck represents the structure for an affordability assessment.
type AffordabilityCheck struct {
	LoanID          string    // Unique loan ID
	ApplicantWallet string    // Wallet address of the applicant
	Income          float64   // Applicant's monthly income
	Expenses        float64   // Applicant's monthly expenses
	Dependents      int       // Number of dependents the applicant has
	DependentCosts  float64   // Monthly cost for dependents
	WorkingStatus   string    // Employment status (e.g., employed, self-employed, unemployed)
	OtherDebts      float64   // Total amount of other debts the applicant owes
	WorkProof       []byte    // Proof of employment (e.g., employment contract or income statement)
	SubmissionTime  time.Time // Timestamp of affordability check submission
	ApprovalStatus  string    // Status of the affordability check (Pending, Approved, Rejected)
	ApprovedBy      string    // Authority node that approved or rejected the check
	ApprovalTime    time.Time // Time of approval or rejection
	Status          string    // Status of the check ("Approved", "Rejected")

}

// AffordabilityManager handles the submission and approval process for affordability checks.
type AffordabilityManager struct {
	mutex             sync.Mutex                     // Mutex for thread safety
	Ledger            *Ledger                        // Ledger for recording affordability checks
	Consensus         *SynnergyConsensus             // Synnergy Consensus engine for validation
	EncryptionService *Encryption                    // Encryption service for securing sensitive financial data
	Submissions       map[string]*AffordabilityCheck // Map of affordability submissions by loan ID
	ApprovalQueue     []*AffordabilityCheck          // Queue of submissions pending approval
}

// UnsecuredLoanApprovalProcess manages the approval process for unsecured loan proposals.
type UnsecuredLoanApprovalProcess struct {
	mutex             sync.Mutex
	Ledger            *Ledger                    // Ledger to store proposal and approval status
	Nodes             []*AuthorityNodeTypes      // List of all authority nodes in the network
	ActiveProposals   map[string]*ActiveProposal // Map of active proposals being reviewed
	EncryptionService *Encryption                // Encryption service for secure transmission
	NetworkManager    *NetworkManager            // Network manager to handle proposal transmission
	RequeueDuration   time.Duration              // Duration before a proposal is requeued (48 hours)
	MaxConfirmations  int                        // Required confirmations for proposal approval
	MaxRejections     int                        // Required rejections for proposal rejection
}

// ActiveProposal tracks confirmations, rejections, interest rates, and assigned nodes for a proposal.
type ActiveProposal struct {
	ProposalID       string                         // Unique proposal ID
	ProposalData     *UnsecuredLoanProposal         // The loan proposal details
	ConfirmedNodes   map[string]bool                // Nodes that confirmed the proposal
	RejectedNodes    map[string]bool                // Nodes that rejected the proposal
	AssignedNodes    map[string]*AuthorityNodeTypes // Nodes currently assigned for review
	Status           string                         // Status of the proposal (Pending, Approved, Rejected)
	InterestRates    []float64                      // List of interest rates submitted by authority nodes
	LastDistribution time.Time                      // Timestamp of last node distribution
	ProposalDeadline time.Time                      // Deadline for the proposal before requeuing
	AverageInterest  float64                        // Running average of the interest rates
	AllDocsOpened    bool                           // Whether all nodes have opened the attached documents
}

// UnsecuredLoanTermManager manages customization of unsecured loan terms.
type UnsecuredLoanTermManager struct {
	mutex             sync.Mutex            // Mutex for thread safety
	Ledger            *Ledger               // Reference to the ledger for storing loan term records
	Consensus         *SynnergyConsensus    // Synnergy Consensus engine for validating terms
	LoanTermRecords   map[string]*LoanTerms // Records of loan terms by loan ID
	EncryptionService *Encryption           // Encryption service for securing sensitive data
}

// DecentralizedCreditCheck tracks spending for wallets and stores credit score documents.
type DecentralizedCreditCheck struct {
	mutex                 sync.Mutex                 // Mutex for thread safety
	Ledger                *Ledger                    // Ledger for storing credit check and transaction data
	Consensus             *SynnergyConsensus         // Consensus engine for validation
	WalletSpendingRecords map[string]*SpendingRecord // Stores spending records by wallet address
	CreditScoreDocuments  map[string][]byte          // Stores encrypted credit score documents by wallet address
	EncryptionService     *Encryption                // Encryption service for securing sensitive data
}

// SpendingRecord represents a record of wallet spending and associated transactions.
type SpendingRecord struct {
	WalletAddress string         // Address of the wallet being tracked
	TotalSpent    float64        // Total amount spent from the wallet
	Transactions  []*Transaction // List of transactions made from the wallet
	LastUpdated   time.Time      // Timestamp of the last update
}

// UnsecuredLoanDisbursementQueueEntry represents a loan waiting for disbursement in the unsecured loan pool.
type UnsecuredLoanDisbursementQueueEntry struct {
	ProposalID        string    // The loan proposal ID
	ProposerWallet    string    // Wallet address of the borrower
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
	AverageInterest   float64   // Average interest rate to be applied to the loan
}

// UnsecuredLoanDisbursementManager manages the disbursement of approved loans in the unsecured loan pool.
type UnsecuredLoanDisbursementManager struct {
	mutex               sync.Mutex                             // Mutex for thread safety
	Ledger              *Ledger                                // Reference to the ledger for recording disbursements
	Consensus           *SynnergyConsensus                     // Synnergy Consensus engine
	FundBalance         float64                                // Current balance of the loan pool
	DisbursementQueue   []*UnsecuredLoanDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime        time.Duration                          // Maximum time a proposal can wait in the queue
	EncryptionService   *Encryption                            // Encryption service for securing sensitive data
	IssuerFeePercentage float64                                // Issuer fee (0.5%)
}

// UnsecuredLoanProposal represents the full loan application process.
type UnsecuredLoanProposal struct {
	LoanID              string             // Unique loan ID
	ApplicantName       string             // Name of the applicant
	ApplicantID         string             // Unique applicant ID (validated by syn900)
	WalletAddress       string             // Wallet address of the applicant
	SubmissionTimestamp time.Time          // Time when the proposal was submitted
	ProposalStatus      string             // Status of the proposal (e.g., Pending, Approved, Rejected)
	LastUpdated         time.Time          // Timestamp of last update
	ApprovalStage       string             // Current stage of the approval process (e.g., Application, CreditCheck, Affordability, Terms)
	CreditScore         float64            // Applicant's credit score (from decentralized credit check)
	AffordabilityStatus string             // Result of affordability check (e.g., Approved, Rejected)
	TermsCustomization  bool               // Whether terms customization has been completed
	Status              string             // Status of the proposal ("Pending", "Approved", "Rejected")
	Repayments          []Repayment        // List of repayments made by the borrower
	LoanTerms           LoanTerms          // Terms of the loan (interest rate, duration, etc.)
	InterestPayments    []InterestPayment  // Interest payments made by the borrower
	AffordabilityCheck  AffordabilityCheck // Results of the affordability check for the borrower
	FundsDisbursed      float64            // Amount of funds disbursed for the loan
}

// ProposalManager handles the overall proposal process for unsecured loans.
type UnsecuredLoanProposalManager struct {
	mutex             sync.Mutex                        // Mutex for thread safety
	Ledger            *Ledger                           // Ledger for storing proposal data
	Consensus         *SynnergyConsensus                // Consensus engine for validation
	Syn900Validator   *tokenledgers.Syn900Validator     // Syn900 validator for ID validation
	EncryptionService *Encryption                       // Encryption service for secure proposal data
	CreditChecker     *CreditCheckManager               // Decentralized Credit Check manager
	AffordabilityMgr  *AffordabilityManager             // Affordability Check manager
	TermsManager      *TermsCustomizationManager        // Customization of loan terms
	Proposals         map[string]*UnsecuredLoanProposal // Proposals mapped by LoanID
}

// LoanRepaymentDetails stores repayment information for a loan.
type UnsecuredLoanRepaymentDetails struct {
	LoanID           string      // Unique loan ID
	ProposerWallet   string      // Wallet address of the borrower
	TotalAmount      float64     // Total amount to be repaid
	RemainingAmount  float64     // Remaining amount to be repaid
	InterestRate     float64     // Interest rate applied to the loan
	RepaymentDates   []time.Time // Scheduled payment dates
	NextPaymentDue   time.Time   // Next payment due date
	Status           string      // Loan status (Active, Defaulted, Satisfied)
	DefaultedAt      *time.Time  // If loan defaulted, record the default date
	DefaultContact   string      // Contact email for default notification
	AuthorityWallets []string    // Wallet addresses of authority nodes
}

// SecuredLoanRepaymentManager manages the repayment and settlement process for secured loans.
type UnsecuredLoanRepaymentManager struct {
	mutex             sync.Mutex                                // Mutex for thread safety
	Ledger            *Ledger                                   // Reference to the ledger
	EncryptionService *Encryption                               // Encryption service for secure data
	Syn900Registry    *SYN900Registry                           // Reference to Syn900 for record keeping
	LoanRepayments    map[string]*UnsecuredLoanRepaymentDetails // Map of loan repayments by loan ID
	DefaultThreshold  time.Duration                             // Time duration for default (e.g., 6 months)
}

// BorrowerDetails represents the borrower information in a loan.
type BorrowerDetails struct {
	LoanID          string
	BorrowerName    string
	BorrowerEmail   string
	BorrowerContact string
	WalletAddress   string // Borrower's wallet address
}

// UnsecuredLoanManagement handles management tasks such as authority node updates, borrower detail changes, and term change requests.
type UnsecuredLoanManagement struct {
	mutex              sync.Mutex
	Ledger             *Ledger                               // Ledger reference for recording updates
	EncryptionService  *Encryption                           // Encryption service for data security
	ConsensusEngine    *SynnergyConsensus                    // Consensus engine for approval processes
	NetworkManager     *NetworkManager                       // Network manager for sending requests
	AuthorityNodes     map[string]*AuthorityNodeTypes        // Map of authority nodes that can manage loans
	LoanBorrowerInfo   map[string]*BorrowerDetails           // Stores borrower details by loan ID
	TermChangeRequests map[string]*BorrowerTermChangeRequest // Stores term change requests by loan ID
}

// Syn900Registry represents the registry for associating a loan with an SYN900 token.
type Syn900Registry struct {
	mutex         sync.Mutex // For thread safety
	LoanID        string     // Unique identifier for the loan
	SYN900TokenID string     // The SYN900 token ID associated with the borrower
	Ledger        *Ledger    // Ledger reference to track loan transactions
}

// UnsecuredLoanPool manages the details of the unsecured loan pool, including balance, loan disbursements, repayments, and defaults.
type UnsecuredLoanPool struct {
	mutex             sync.Mutex                        // Mutex for thread safety
	TotalBalance      *big.Int                          // Total balance available in the loan pool
	LoansDistributed  *big.Int                          // Total amount of loans distributed
	LoansRepaid       *big.Int                          // Total amount of loans repaid
	LoansDefaulted    *big.Int                          // Total amount of loans that have defaulted
	Ledger            *Ledger                           // Reference to the ledger for storing loan transactions
	Consensus         *SynnergyConsensus                // Synnergy Consensus engine for transaction validation
	EncryptionService *Encryption                       // Encryption service for securing sensitive loan data
	LoanRecords       map[string]*LoanRecord            // Map of loan records by applicant wallet
	EncryptedData     string                            // Encrypted data for privacy and security
	Proposals         map[string]*UnsecuredLoanProposal // Tracks loan proposals
	FundsAvailable    float64                           // Funds currently available for new loans

}

// LoanRecord tracks details for each unsecured loan.
type LoanRecord struct {
	ProposalID       string              // Unique ID for the loan proposal
	ApplicantWallet  string              // Wallet address of the loan applicant
	LoanAmount       *big.Int            // Amount of the loan distributed
	LoanRepaid       *big.Int            // Amount repaid so far
	LoanStatus       string              // Status of the loan: Active, Repaid, Defaulted
	RepaymentDueDate time.Time           // Due date for full repayment
	Transactions     []TransactionRecord // Add this field to store transactions
}

// ApprovalStage represents a stage in a multi-step approval process.
type ApprovalStage struct {
	StageID         string    // Unique identifier for the approval stage
	Description     string    // Description of the approval stage (e.g., "KYC Verification", "AML Check")
	IsApproved      bool      // Whether this stage has been approved
	ApproverID      string    // The identifier of the entity or individual who approved this stage
	ApprovalTime    time.Time // The timestamp when the stage was approved
	RejectionReason string    // Reason for rejection, if applicable
	IsRejected      bool      // Whether this stage has been rejected
}

// ************** Marketplace Frameworks Structs **************
// OrderType defines the type of the order (buy or sell)
type OrderType string

const (
	BuyOrder  OrderType = "buy"
	SellOrder OrderType = "sell"
)

// TradeCompletion represents the completion of a trade in the marketplace.
type TradeCompletion struct {
	TradeID   string    // Unique identifier for the trade
	BuyerID   string    // ID of the buyer
	SellerID  string    // ID of the seller
	ItemID    string    // ID of the item being traded
	Timestamp time.Time // Time the trade was completed
}

// IllegalItemReport represents a report of an illegal item in the marketplace.
type IllegalItemReport struct {
	ReportID    string    // Unique identifier for the report
	ItemID      string    // ID of the illegal item
	ListingID   string    // ID of the reported listing (Add this field)
	ReporterID  string    // ID of the user who reported the item
	Description string    // Description of the illegal activity
	Timestamp   time.Time // Time the report was submitted
}

// ResourceListing represents an item listed in the marketplace.
type ResourceListing struct {
	ListingID  string    // Unique identifier for the listing
	OwnerID    string    // ID of the item owner
	ResourceID string    // ID of the resource being listed
	Price      float64   // Price of the resource
	Lease      bool      // Whether the resource is for lease
	Timestamp  time.Time // Time the listing was created

}

// AMMManager manages liquidity pools and trades within the AMM
type AMMManager struct {
	Pools             map[string]*LiquidityPool // Active liquidity pools
	Ledger            *Ledger                   // Ledger instance for tracking trades and liquidity actions
	EncryptionService *Encryption               // Encryption service for securing trade and liquidity data
	mu                sync.Mutex                // Mutex for concurrent management of pools and trades
}

// TokenOrder represents a buy or sell order on the exchange
type TokenOrder struct {
	OrderID    string    // Unique identifier for the order
	TokenID    string    // The token being traded
	OrderType  OrderType // Whether this is a buy or sell order
	Amount     float64   // Amount of tokens to buy or sell
	Price      float64   // Price per token
	Trader     string    // Wallet address of the trader
	Timestamp  time.Time // Time when the order was created
	IsExecuted bool      // Whether the order has been executed
}

// TokenPair represents the trading pair (token being traded for another)
type TokenPair struct {
	TokenID   string        // The token being traded
	BaseToken string        // The base token (e.g., ETH, USDC) against which the token is traded
	Orders    []*TokenOrder // The list of active orders for this token pair
}

// CentralizedTokenExchange manages buy and sell orders for tokens
type CentralizedTokenExchange struct {
	TokenPairs        map[string]*TokenPair // Active token pairs and their respective orders
	Ledger            *Ledger               // Ledger for logging all transactions
	EncryptionService *Encryption           // Encryption for securing sensitive data
	mu                sync.Mutex            // Mutex for concurrency control
}

// Escrow holds funds for transactions in the marketplace until conditions are met
type Escrow struct {
	EscrowID       string    // Unique identifier for the escrow
	Buyer          string    // Wallet address of the buyer
	Seller         string    // Wallet address of the seller
	Amount         float64   // Amount held in escrow
	ResourceID     string    // ID of the resource for this escrow
	CompletionTime time.Time // Timestamp when the escrow is completed
	IsReleased     bool      // Whether the funds have been released
	IsDisputed     bool      // Whether the transaction is in dispute
	Status         string    // Current status of the escrow (e.g., active, completed)
	Timestamp      time.Time // Timestamp when the escrow was created
}

// ComputerResourceMarketplace manages the listing, rental, and escrow of computing resources
type ComputerResourceMarketplace struct {
	Resources         map[string]*Resource // List of available resources
	Escrows           map[string]*Escrow   // Active escrows for resource transactions
	Ledger            *Ledger              // Ledger instance for recording transactions
	EncryptionService *Encryption          // Encryption for securing sensitive data
	mu                sync.Mutex           // Mutex for concurrent operations
}

// Order represents a trading order in the decentralized exchange
type Order struct {
	OrderID       string    // Unique identifier for the order
	Trader        string    // Wallet address of the trader
	AssetIn       string    // The asset being traded from
	AmountIn      float64   // Amount of asset being offered
	AssetOut      string    // The asset being traded for
	AmountOut     float64   // Amount of asset expected in return
	OrderType     string    // "Buy" or "Sell"
	OrderTime     time.Time // Time when the order was placed
	IsFilled      bool      // Whether the order has been filled
	TransactionID string    // Associated transaction ID when the order is filled
	BuyerID       string    // ID of the buyer placing the order
	SellerID      string    // ID of the seller
	ItemID        string    // ID of the item being ordered
	Quantity      int       // Quantity of the item
	Price         float64   // Price per item
	Status        string    // Order status (e.g., pending, completed)
	Timestamp     time.Time // Time the order was placed
}

// DEXManager manages decentralized trading orders and matches orders
type DEXManager struct {
	Orders            map[string]*Order // Active orders in the DEX
	CompletedOrders   map[string]*Order // Filled orders
	Ledger            *Ledger           // Ledger for logging trades and order completion
	EncryptionService *Encryption       // Encryption service for securing order data
	mu                sync.Mutex        // Mutex for concurrent order management
}

// Listing represents a product or service listed on the decentralized marketplace
type Listing struct {
	ListingID   string    // Unique identifier for the listing
	ItemName    string    // Name of the item or service
	Description string    // Description of the item or service
	Price       float64   // Price for buying the item
	LeasePrice  float64   // Price for leasing the item per day
	Owner       string    // Wallet address of the item owner
	Available   bool      // Whether the item is available for purchase or lease
	Category    string    // Category of the item
	ListedTime  time.Time // Timestamp when the item was listed
	Leasing     bool      // Whether leasing is allowed
	IsLegal     bool      // Flag to ensure the listing complies with legal standards
}

// GeneralMarketplace manages the decentralized listing, buying, leasing, and escrow
type GeneralMarketplace struct {
	Listings          map[string]*Listing // Available listings in the marketplace
	Escrows           map[string]*Escrow  // Active escrows for transactions
	Ledger            *Ledger             // Ledger for logging all transactions
	EncryptionService *Encryption         // Encryption service for securing sensitive data
	mu                sync.Mutex          // Mutex for concurrency control
}

// NFTListing represents an NFT listed for sale in the marketplace
type NFTListing struct {
	ListingID   string    // Unique identifier for the listing
	OwnerID     string    // ID of the owner listing the NFT (Add this field)
	TokenID     string    // The unique token ID of the NFT
	Standard    string    // The token standard, either "Syn721" or "Syn1155"
	MetadataURI string    // The URI that points to the NFT's metadata
	Price       float64   // Price of the NFT
	Owner       string    // Wallet address of the current owner
	Available   bool      // Whether the NFT is available for sale
	ListedTime  time.Time // Timestamp when the NFT was listed
}

// NFTMarketplace manages the listing, buying, and selling of NFTs
type NFTMarketplace struct {
	Listings          map[string]*NFTListing // Active NFT listings in the marketplace
	Ledger            *Ledger                // Ledger for logging all NFT transactions
	EncryptionService *Encryption            // Encryption for securing sensitive data
	mu                sync.Mutex             // Mutex for concurrency control
}

// StakingPool represents a pool where users can stake tokens for a project
type StakingPool struct {
	PoolID              string               // Unique identifier for the staking pool
	ProjectName         string               // Name of the project the pool supports
	TokenAddress        string               // Address of the token being staked
	Owner               string               // Owner or creator of the staking pool
	StakedAmount        float64              // Total amount of tokens staked in the pool
	RewardRate          float64              // Reward rate for stakers (e.g., percentage per day)
	StartTime           time.Time            // Start time of the staking period
	EndTime             time.Time            // End time of the staking period
	IsActive            bool                 // Whether the pool is currently active
	Participants        map[string]float64   // Tracks each participant's staked amount
	Unstakes            []UnstakeTransaction // List of unstaking transactions (add this field)
	RewardDistributions []RewardDistribution // List of reward distributions (add this field)

}

// StakingLaunchpad manages staking pools and allows users to stake tokens
type StakingLaunchpad struct {
	Pools             map[string]*StakingPool // Active staking pools in the launchpad
	Ledger            *Ledger                 // Ledger for logging all staking transactions
	EncryptionService *Encryption             // Encryption for securing sensitive data
	mu                sync.Mutex              // Mutex for concurrency control
}

// MarketplaceManager manages various marketplace activities.
type MarketplaceManager struct {
	Orders            map[string]*Order            // Tracks marketplace orders
	Trades            map[string]*TradeExecution   // Tracks executed trades
	ItemListings      map[string]*Listing          // Tracks listings of general items
	ResourceListings  map[string]*Resource         // Tracks resource listings
	ResourcePurchases map[string]*ResourcePurchase // Tracks resource purchases
	ResourceRentals   map[string]*ResourceRental   // Tracks rentals of resources
	EscrowDisputes    map[string]*EscrowDispute    // Tracks escrow-related disputes
	Escrows           map[string]*Escrow           // Tracks escrows
}

// MarketplaceState represents the state of the general marketplace.
type MarketplaceState struct {
	Orders             map[string]*Order            // Active orders in the marketplace
	Trades             map[string]*TradeExecution   // Executed trades
	ItemListings       map[string]*ItemListing      // Adjust type to store ItemListing pointers
	ResourceListings   map[string]*ResourceListing  // Store resource listings
	ResourcePurchases  map[string]*ResourcePurchase // Resource purchases
	ResourceRentals    map[string]*ResourceRental   // Resource rentals
	EscrowDisputes     map[string]*EscrowDispute    // Disputes related to escrow transactions
	Escrows            map[string]*Escrow           // Escrows
	TradeCompletions   map[string]*TradeCompletion  // Change this to store pointers
	ItemPurchases      map[string]*ItemPurchase     // Change this to store ItemPurchase pointers
	ItemLeases         map[string]*ItemLease        // Store ItemLease pointers
	IllegalItemReports map[string]IllegalItemReport // Store values, not pointers

}

// NFTMarketplaceState represents the state of the NFT marketplace.
type NFTMarketplaceState struct {
	Listings  map[string]*NFTListing  // Active NFT listings
	Purchases map[string]*NFTPurchase // NFT purchases
}

// DEXManagerState represents the state of the decentralized exchange.
type DEXManagerState struct {
	Orders          map[string]*Order // Orders placed in the DEX
	CompletedOrders map[string]*Order // Completed orders in the DEX
}

// AMMManagerState represents the state of the Automated Market Maker (AMM).
type AMMManagerState struct {
	Pools map[string]*LiquidityPool // Active liquidity pools for the AMM
}

// ComputerResourceMarketState represents the state of the computer resource marketplace.
type ComputerResourceMarketState struct {
	Resources map[string]*Resource // Available resources in the market
	Escrows   map[string]*Escrow   // Escrows for computer resource transactions
}

// CentralizedExchangeState represents the state of the centralized token exchange.
type CentralizedExchangeState struct {
	TokenPairs map[string]*TokenPair // Active token pairs for trading
}

// TradeExecution represents the execution of a trade between two parties.
type TradeExecution struct {
	TradeID   string    // Unique identifier for the trade
	BuyerID   string    // Wallet address of the buyer
	SellerID  string    // Wallet address of the seller
	AssetID   string    // ID of the asset being traded
	Amount    float64   // Amount of the asset traded
	Price     float64   // Price per unit of the asset
	Timestamp time.Time // Timestamp of when the trade was executed
}

// ResourcePurchase represents a purchase of a resource in the marketplace.
type ResourcePurchase struct {
	PurchaseID   string    // Unique identifier for the purchase
	BuyerID      string    // Wallet address of the buyer
	ResourceID   string    // ID of the resource being purchased
	Amount       float64   // Amount paid for the resource
	Timestamp    time.Time // Timestamp of when the purchase was made
	ListingID    string    // ID of the resource listing being purchased
	PurchaseTime time.Time // Time when the purchase was made
}

// ResourceRental represents the rental of a resource in the marketplace.
type ResourceRental struct {
	RentalID     string        // Unique identifier for the rental
	RenterID     string        // Wallet address of the renter
	ResourceID   string        // ID of the resource being rented
	RentalPeriod time.Duration // Duration of the rental
	Amount       float64       // Amount paid for the rental period
	Timestamp    time.Time     // Timestamp of when the rental transaction occurred
	ListingID    string        // ID of the resource listing being rented
	RentalTime   time.Time     // Time when the rental agreement was made
}

// EscrowDispute represents a dispute raised regarding an escrow transaction.
type EscrowDispute struct {
	DisputeID     string    // Unique identifier for the dispute
	TransactionID string    // ID of the escrow transaction being disputed
	PartyA        string    // Wallet address of the first party in the dispute
	PartyB        string    // Wallet address of the second party in the dispute
	DisputeReason string    // Reason for the dispute
	Timestamp     time.Time // Timestamp of when the dispute was raised
}

// NFTPurchase represents the purchase of an NFT in the marketplace.
type NFTPurchase struct {
	PurchaseID string    // Unique identifier for the NFT purchase
	BuyerID    string    // Wallet address of the buyer
	ListingID  string    // ID of the NFT listing
	TokenID    string    // ID of the NFT token being purchased
	Price      float64   // Price paid for the NFT
	Timestamp  time.Time // Timestamp of when the purchase was made
}

// ItemListing represents a general listing in the marketplace (for sale or rent).
type ItemListing struct {
	ListingID  string
	ItemName   string
	Price      float64
	LeasePrice float64
	OwnerID    string
	Available  bool
	ListedTime time.Time
}

// ItemPurchase represents a purchase of an item.
type ItemPurchase struct {
	PurchaseID string
	BuyerID    string
	ListingID  string
	Price      float64
	Timestamp  time.Time
}

// ItemLease represents a lease of an item.
type ItemLease struct {
	LeaseID     string
	LeaserID    string
	ListingID   string
	LeasePrice  float64
	LeasePeriod time.Duration
	Timestamp   time.Time
}

// StakeTransaction represents a staking transaction.
type StakeTransaction struct {
	TransactionID string
	ParticipantID string
	AmountStaked  float64
	Timestamp     time.Time
}

// UnstakeTransaction represents an unstaking transaction.
type UnstakeTransaction struct {
	TransactionID string
	ParticipantID string
	Timestamp     time.Time
	Amount        float64
}

// RewardDistribution represents the distribution of staking rewards.
type RewardDistribution struct {
	DistributionID string
	Amount         float64
	Timestamp      time.Time
}

// PacketEvent represents a packet-related event in the ledger
type PacketEvent struct {
	PacketID    string    // Unique identifier for the packet
	Timestamp   time.Time // When the packet event occurred
	Data        []byte    // Packet data (encrypted)
	Destination string    // Destination peer address
	Size        int       // Size of the payload
	Priority    int       // Priority level of the packet
}


type AIMarketplaceConfig struct {
    InitialModules []AIModule
    CreatedAt      time.Time
}

type AIModule struct {
    ID                  string
    Name                string
    Description         string
    OwnerID             string
    UsagePrice          float64
    ForSale             bool
    SalePrice           float64
    RegisteredAt        time.Time
    UpdatedAt           time.Time
    UsageTrackingEnabled bool
}

type AIRental struct {
    ModuleID   string
    RenterID   string
    StartTime  time.Time
    EndTime    time.Time
}

type UsageStats struct {
    UserID     string
    UsageCount int
    Duration   time.Duration
}

type AILog struct {
    ModuleID  string
    Usage     UsageStats
    Timestamp time.Time
}

type AIResourceRequest struct {
    ModuleID    string
    Resources   NetworkResources
    RequestTime time.Time
}

type NetworkResources struct {
    CPU    int
    Memory int
    Disk   int
}
type AIPermission struct {
    ReadAccess   bool
    WriteAccess  bool
    ExecuteAccess bool
    AdminAccess  bool
}

type AITransaction struct {
    ModuleID        string
    TransactionType string
    Details         TransactionDetails
    Timestamp       time.Time
}

type TransactionDetails struct {
    Amount   float64
    UserID   string
    Metadata string
}

type AIModuleReport struct {
    Module       AIModule
    Usage        []AILog
    Transactions []AITransaction
}

type AIInputData struct {
    Parameters map[string]interface{}
    Metadata   string
}

type AIOutputData struct {
    Result    string
    Timestamp time.Time
}

type AIUsageSchedule struct {
    ModuleID     string
    Frequency    time.Duration
    ScheduledAt  time.Time
    StartTime    time.Time
    EndTime      time.Time
}

type AIEventLog struct {
    ModuleID    string
    EventType   string
    Description string
    Timestamp   time.Time
}

type AIResourceAllocation struct {
    ModuleID  string
    CPU       int
    Memory    int
    Disk      int
    Timestamp time.Time
}

type AIMetrics struct {
    ModuleID    string
    Performance string
    Accuracy    float64
    Latency     time.Duration
    Timestamp   time.Time
}

type AITrainingData struct {
    ModuleID  string
    DataHash  string
    UpdatedAt time.Time
}

type AIModelVersion struct {
    ModuleID  string
    Version   string
    Changes   AIModelChanges
    Timestamp time.Time
}

type AIModelChanges struct {
    AddedFeatures   []string
    RemovedFeatures []string
    Optimization    string
}

type AITask struct {
    TaskID      string
    ModuleID    string
    Description string
    AssignedAt  time.Time
    Completed   bool
    CompletedAt time.Time
}

type AIReward struct {
    ModuleID string
    Amount   float64
    IssuedAt time.Time
    Reason   string
}

type AIPenalty struct {
    ModuleID string
    Amount   float64
    IssuedAt time.Time
    Reason   string
}

type AIDatasetLink struct {
    ModuleID  string
    DatasetID string
    LinkedAt  time.Time
}

type LiquidityFee struct {
    PairID    string
    FeeRate   float64
    Timestamp time.Time
}

type UserLiquidity struct {
    UserID    string
    PairID    string
    Amount    float64
    Timestamp time.Time
}

type TradeReport struct {
    PairID      string
    Trades      []TradeDetails
    Period      string
    GeneratedAt time.Time
}

type SlippageSettings struct {
    PairID     string
    Tolerance  float64
    Timestamp  time.Time
}

type TradeVolume struct {
    PairID    string
    Volume    float64
    Timestamp time.Time
}

type TradeExpiry struct {
    PairID         string
    ExpiryDuration time.Duration
    SetAt          time.Time
}

type PriceFluctuation struct {
    PairID    string
    Change    float64
    Direction string
    Timestamp time.Time
}

type OrderBookDepth struct {
    PairID      string
    BidDepth    float64
    AskDepth    float64
    Timestamp   time.Time
}

type FeeStructure struct {
    MakerFee float64
    TakerFee float64
    PairID   string
    Timestamp time.Time
}

type FeeHistory struct {
    PairID    string
    FeeRate   float64
    Timestamp time.Time
}

type PoolTokenRatio struct {
    PairID      string
    TokenARatio float64
    TokenBRatio float64
    Timestamp   time.Time
}

type LiquidityProvision struct {
    UserID    string
    PairID    string
    Amount    float64
    Timestamp time.Time
}

type LiquidityReport struct {
    PairID      string
    Provisions  []LiquidityProvision
    Withdrawals []LiquidityWithdrawal
    Period      string
}

type LiquidityWithdrawal struct {
    UserID    string
    PairID    string
    Amount    float64
    Timestamp time.Time
}

type LiquidityYield struct {
    PairID    string
    Yield     float64
    Period    time.Duration
    Timestamp time.Time
}

type DEXConfig struct {
    Name           string
    Owner          string
    InitializedAt  time.Time
    SupportedPairs []string
}

type TradingPair struct {
    PairID    string
    TokenA    string
    TokenB    string
    CreatedAt time.Time
}

type LiquidityPool struct {
    ProviderID string
    PairID     string
    TokenA     float64
    TokenB     float64
    AddedAt    time.Time
    RemovedAt  time.Time
}

type Swap struct {
    PairID     string
    AmountIn   float64
    TokenIn    string
    AmountOut  float64
    TokenOut   string
    ExecutedAt time.Time
}

type Order struct {
    OrderID   string
    PairID    string
    Type      string
    Amount    float64
    Price     float64
    Status    string
    PlacedAt  time.Time
}

type TradingFee struct {
    PairID    string
    FeeRate   float64
    Timestamp time.Time
}

type TradeExecution struct {
    TradeID    string
    PairID     string
    Amount     float64
    Price      float64
    ExecutedAt time.Time
}

type PriceImpact struct {
    PairID       string
    TradeAmount  float64
    ImpactPercent float64
}

type OrderCancellation struct {
    OrderID     string
    CancelledAt time.Time
}

type MinimumTradeAmount struct {
    PairID string
    Amount float64
    SetAt  time.Time
}

type PoolReward struct {
    PairID         string
    TotalReward    float64
    DistributedAt  time.Time
    ProviderRewards map[string]float64 // ProviderID to Reward
}

type NFTTradeDenial struct {
    TradeID          string
    Reason           string
    EncryptedReason  string
    DeniedAt         time.Time
}

type NFTTrade struct {
    TradeID     string
    Completed   bool
    CompletedAt time.Time
}

type NFTExchangeRate struct {
    NFTID   string
    Rate    float64
    SetAt   time.Time
}

type ExchangeTransaction struct {
    NFTID      string
    BuyerID    string
    SellerID   string
    Price      float64
    Timestamp  time.Time
}

type NFTExchangeReport struct {
    Transactions []ExchangeTransaction
    Period       string
    GeneratedAt  time.Time
}

type NFTMintingLimit struct {
    NFTID   string
    Limit   int
    SetAt   time.Time
}

type NFTMintingEvent struct {
    NFTID     string
    Amount    int
    MintedAt  time.Time
}

type NFTMintingReport struct {
    NFTID  string
    Events []NFTMintingEvent
    Period string
}

// MintingAuthorization struct for recording authorization details
type MintingAuthorization struct {
    RequestID       string
    Status          string
    Reason          string
    EncryptedReason string
    UpdatedAt       time.Time
}

// NFTCustomizationOptions struct for storing customization options
type NFTCustomizationOptions struct {
    NFTID  string
    Options map[string]interface{}
    SetAt  time.Time
}

// NFTCustomization struct for tracking NFT customizations
type NFTCustomization struct {
    NFTID      string
    Details    map[string]interface{}
    Timestamp  time.Time
}

// CustomizationEvent struct for logging NFT customization events
type CustomizationEvent struct {
    NFTID               string
    Description         string
    EncryptedDescription string
    LoggedAt            time.Time
}

// StakeReward struct for distributing stake rewards
type StakeReward struct {
    HolderID     string
    Amount       float64
    DistributedAt time.Time
    NFTID        string
}

// StakeReport struct for generating reports
type StakeReport struct {
    NFTID      string
    Rewards    []StakeReward
    Period     string
}

// StakeDistribution struct to store stake distribution details
type StakeDistribution struct {
    NFTID     string
    IsValid   bool
    Details   map[string]interface{}
    Timestamp time.Time
}

// CrossMarketplaceRate struct to store exchange rates for cross-marketplace trades
type CrossMarketplaceRate struct {
    NFTID     string
    Rate      float64
    Timestamp time.Time
}

// CrossMarketplaceTrade struct to log cross-marketplace trades
type CrossMarketplaceTrade struct {
    NFTID     string
    BuyerID   string
    SellerID  string
    Amount    float64
    Timestamp time.Time
}

// CrossMarketplaceMetrics struct to track metrics for cross-marketplace trades
type CrossMarketplaceMetrics struct {
    NFTID     string
    Metrics   map[string]interface{}
    Timestamp time.Time
}

// CrossMarketplaceStatus struct to verify if an NFT is listed across marketplaces
type CrossMarketplaceStatus struct {
    NFTID   string
    Listed  bool
    Details string
}

// UserRating struct for storing user ratings for NFTs
type UserRating struct {
    NFTID     string
    UserID    string
    Rating    int
    Timestamp time.Time
}

// UserFeedback struct to record encrypted user feedback for NFT transactions
type UserFeedback struct {
    NFTID     string
    UserID    string
    Feedback  string
    Timestamp time.Time
}

// RatingSummary struct to summarize user ratings for an NFT
type RatingSummary struct {
    NFTID         string
    AverageRating float64
    TotalRatings  int
}

// RatingActivity struct to log user activity related to NFT ratings
type RatingActivity struct {
    UserID    string
    NFTID     string
    Action    string
    Timestamp time.Time
}

// NFTInheritanceRights struct for recording inheritance rights of an NFT
type NFTInheritanceRights struct {
    NFTID         string
    BeneficiaryID string
    SetAt         time.Time
}

// InheritanceActivity struct for logging inheritance-related activities
type InheritanceActivity struct {
    NFTID       string
    Description string
    Timestamp   time.Time
}

// InheritanceReport struct for generating inheritance reports
type InheritanceReport struct {
    NFTID      string
    Activities []InheritanceActivity
    Period     string
}

// NFTBundle struct for creating bundles of NFTs
type NFTBundle struct {
    BundleID  string
    NFTs      []string
    CreatedAt time.Time
}

// NFTBundleEntry struct for adding NFTs to an existing bundle
type NFTBundleEntry struct {
    BundleID string
    NFTID    string
    AddedAt  time.Time
}

// EscrowStatus struct for tracking escrow status of an NFT
type EscrowStatus struct {
    NFTID      string
    Status     string
    Timestamp  time.Time
}

// EscrowReport struct for generating escrow activity reports
type EscrowReport struct {
    NFTID    string
    Statuses []EscrowStatus
    Period   string
}

// NFTRentalTerms struct for storing rental terms of an NFT
type NFTRentalTerms struct {
    NFTID     string
    Terms     map[string]interface{}
    SetAt     time.Time
}

// RentalPayment struct for tracking payments made towards NFT rentals
type RentalPayment struct {
    NFTID      string
    Amount     float64
    PaidBy     string
    Timestamp  time.Time
}

// RentalActivity struct for logging rental-related activities
type RentalActivity struct {
    NFTID       string
    Description string
    Timestamp   time.Time
}

// RentalContract struct to verify rental contracts for NFTs
type RentalContract struct {
    ContractID string
    IsValid    bool
}

// VerificationBadge struct for NFT verification badges
type VerificationBadge struct {
    NFTID      string
    GrantedAt  time.Time
    RevokedAt  time.Time
}

// NFTCollectionEntry struct for managing NFT collections
type NFTCollectionEntry struct {
    CollectionID string
    NFTID        string
    AddedAt      time.Time
}

// CollectionOwnershipChange struct for tracking ownership changes within a collection
type CollectionOwnershipChange struct {
    CollectionID string
    NFTID        string
    NewOwnerID   string
    Timestamp    time.Time
}

// CollectionActivity struct for recording activities related to NFT collections
type CollectionActivity struct {
    CollectionID string
    Description  string
    Timestamp    time.Time
}

// CollectionReport struct for generating reports on NFT collection activities
type CollectionReport struct {
    CollectionID string
    Activities   []CollectionActivity
    Period       string
}

// NFTTradeEvent struct for logging trade events of NFTs
type NFTTradeEvent struct {
    TradeID     string
    NFTID       string
    BuyerID     string
    SellerID    string
    Amount      float64
    Timestamp   time.Time
}

// NFTTradeStatus struct for managing the status of NFT trades
type NFTTradeStatus struct {
    TradeID string
    Status  string
    UpdatedAt time.Time
}

// NFTMarketplaceConfig struct for initializing the NFT marketplace
type NFTMarketplaceConfig struct {
    MarketplaceName string
    OwnerID         string
    InitializedAt   time.Time
}

// NFT struct for storing NFT metadata
type NFT struct {
    NFTID     string
    Metadata  map[string]interface{}
    MintedAt  time.Time
}

// NFTBurnRecord struct for recording burned NFTs
type NFTBurnRecord struct {
    NFTID    string
    BurnedAt time.Time
}

// NFTOwnershipTransfer struct for recording ownership transfer of NFTs
type NFTOwnershipTransfer struct {
    NFTID         string
    NewOwnerID    string
    TransferredAt time.Time
}

// NFTSale struct for listing NFTs for sale
type NFTSale struct {
    NFTID      string
    Price      float64
    ListedAt   time.Time
    UpdatedAt  time.Time
}

// NFTBid struct for recording bids on NFTs
type NFTBid struct {
    BidID     string
    NFTID     string
    BidderID  string
    Amount    float64
    Timestamp time.Time
}

// NFTAuction struct for starting NFT auctions
type NFTAuction struct {
    AuctionID string
    NFTID     string
    StartTime time.Time
    StartedAt time.Time
}

// NFTAuctionEnd struct for concluding NFT auctions
type NFTAuctionEnd struct {
    AuctionID string
    EndedAt   time.Time
}

// NFTAuctionStatus struct for tracking auction status
type NFTAuctionStatus struct {
    AuctionID  string
    Status     string
    HighestBid float64
    EndsAt     time.Time
}

// NFTAuctionEvent struct for logging events related to NFT auctions
type NFTAuctionEvent struct {
    AuctionID   string
    Description string
    Timestamp   time.Time
}

// NFTAuctionReport struct for generating auction reports
type NFTAuctionReport struct {
    AuctionID string
    Events    []NFTAuctionEvent
    Bids      []NFTBid
    Period    string
}

// NFTOwnership struct for tracking the ownership of NFTs
type NFTOwnership struct {
    NFTID   string
    OwnerID string
}

// NFTMetadata struct for storing metadata of NFTs
type NFTMetadata struct {
    NFTID      string
    Data       map[string]interface{}
    UpdatedAt  time.Time
}

// NFTAuthenticity struct for tracking the authenticity of an NFT
type NFTAuthenticity struct {
    NFTID     string
    IsGenuine bool
    VerifiedAt time.Time
}

// NFTOwnershipHistory struct for storing ownership history of NFTs
type NFTOwnershipHistory struct {
    NFTID     string
    OwnerID   string
    ChangedAt time.Time
}

// NFTTransferEvent struct for logging NFT transfer events
type NFTTransferEvent struct {
    NFTID      string
    FromUserID string
    ToUserID   string
    Timestamp  time.Time
}

// NFTListing struct for managing NFT listings in the marketplace
type NFTListing struct {
    NFTID       string
    ScheduledAt time.Time
    Status      string
}

// NFTStaking struct for recording NFT staking details
type NFTStaking struct {
    NFTID     string
    UserID    string
    Duration  int64
    StakedAt  time.Time
}

// NFTUnstake struct for tracking the unstaking of NFTs
type NFTUnstake struct {
    NFTID      string
    UserID     string
    UnstakedAt time.Time
}

// NFTStakingEvent struct for logging staking events
type NFTStakingEvent struct {
    NFTID      string
    Description string
    Timestamp  time.Time
}

// NFTRoyalty struct for managing royalties on NFT sales
type NFTRoyalty struct {
    NFTID      string
    Percentage float64
    SetAt      time.Time
}

// RoyaltyDistribution struct for recording royalty distributions
type RoyaltyDistribution struct {
    NFTID         string
    Amount        float64
    DistributedAt time.Time
}

// RoyaltyReport struct for generating reports on royalty distributions
type RoyaltyReport struct {
    NFTID         string
    Distributions []RoyaltyDistribution
    Period        string
}

// RoyaltyDistribution struct for tracking total royalty distributions
type RoyaltyDistribution struct {
    NFTID         string
    Amount        float64
    DistributedAt time.Time
}

// FractionalOwnership struct for tracking fractional ownership details
type FractionalOwnership struct {
    NFTID       string
    OwnerID     string
    Fraction    float64 // Represents percentage of ownership
    AssignedAt  time.Time
}

// FractionalOwnershipChange struct for recording changes to fractional ownership
type FractionalOwnershipChange struct {
    NFTID      string
    OwnerID    string
    Fraction   float64
    Timestamp  time.Time
}

// NFTEscrowRelease struct for recording the release of NFTs from escrow
type NFTEscrowRelease struct {
    EscrowID    string
    ReleasedAt  time.Time
}



// ************** Monitoring and Maintenance Structs **************

// CleanupManager manages cleanup tasks and operations to ensure system stability and resource optimization.
type CleanupManager struct {
	TaskID          string                 // Unique identifier for the cleanup task.
	TargetComponent string                 // Component or resource targeted for cleanup.
	Status          string                 // Status of the cleanup task (e.g., pending, in-progress, completed).
	ScheduledTime   time.Time              // Time when the cleanup task is scheduled.
	LastRunTime     time.Time              // Timestamp of the last execution of the cleanup task.
	CreatedAt       time.Time              // Timestamp when the cleanup task was created.
	Metadata        map[string]interface{} // Additional metadata related to the cleanup operation.
}

// MonitoringSystem manages the network monitoring process
type MonitoringSystem struct {
	mutex      sync.Mutex
	active     bool
	Nodes      []string
	trafficLog []TrafficData
}

// CleanupTask defines the details of an individual cleanup operation.
type CleanupTask struct {
	TaskID      string    // Unique identifier for the cleanup task.
	ResourceID  string    // ID of the resource being cleaned up.
	Description string    // Description of the cleanup task.
	CreatedAt   time.Time // Timestamp when the cleanup task was created.
	Status      string    // Current status of the cleanup task.
}

// LogManager handles the creation, storage, and retrieval of system logs.
type LogManager struct {
	LogID       string                 // Unique identifier for the log entry.
	LogType     string                 // Type of log (e.g., error, event, debug).
	Description string                 // Description of the log entry.
	Timestamp   time.Time              // Time when the log entry was created.
	Severity    string                 // Severity level of the log (e.g., info, warning, error).
	ComponentID string                 // Component or process associated with the log entry.
	Metadata    map[string]interface{} // Additional metadata related to the log entry.
}

// DiagnosticManager oversees diagnostic tasks, tests, and health checks.
type DiagnosticManager struct {
	DiagnosticID string                 // Unique identifier for the diagnostic task.
	TestType     string                 // Type of diagnostic test (e.g., health check, performance test).
	ComponentID  string                 // Component or system being tested.
	Status       string                 // Status of the diagnostic task (e.g., pending, completed).
	LastRunTime  time.Time              // Timestamp of the last diagnostic run.
	CreatedAt    time.Time              // Timestamp when the diagnostic task was created.
	Results      map[string]interface{} // Results or findings from the diagnostic tests.
}

// DiagnosticResult captures the result of a diagnostic operation.
type DiagnosticResult struct {
	ResultID  string    // Unique identifier for the diagnostic result.
	TestID    string    // ID of the diagnostic test.
	Outcome   string    // Outcome of the test (e.g., pass, fail).
	Details   string    // Detailed description of the diagnostic result.
	Timestamp time.Time // Timestamp of the diagnostic result.
}

// MaintenanceManager manages scheduled and ad-hoc maintenance operations.
type MaintenanceManager struct {
	MaintenanceID   string                 // Unique identifier for the maintenance operation.
	TargetComponent string                 // Component or system targeted for maintenance.
	MaintenanceType string                 // Type of maintenance (e.g., scheduled, emergency).
	Status          string                 // Status of the maintenance operation (e.g., pending, in-progress, completed).
	ScheduledTime   time.Time              // Time when the maintenance is scheduled.
	CreatedAt       time.Time              // Timestamp when the maintenance task was created.
	Metadata        map[string]interface{} // Additional metadata related to the maintenance operation.
}

// MaintenanceRecord captures the details of a maintenance activity.
type MaintenanceRecord struct {
	RecordID    string    // Unique identifier for the maintenance record.
	Description string    // Description of the maintenance activity.
	StartTime   time.Time // Start time of the maintenance activity.
	EndTime     time.Time // End time of the maintenance activity.
	Status      string    // Status of the maintenance (e.g., completed, failed).
}

// ObserverManager tracks and manages observers monitoring system components or processes.
type ObserverManager struct {
	ObserverID      string                 // Unique identifier for the observer.
	TargetComponent string                 // Component or process being monitored.
	ObserverType    string                 // Type of observer (e.g., real-time, scheduled).
	Status          string                 // Status of the observer (e.g., active, inactive).
	CreatedAt       time.Time              // Timestamp when the observer was registered.
	LastActivity    time.Time              // Timestamp of the last activity by the observer.
	Metadata        map[string]interface{} // Additional metadata related to the observer.
}

// ObserverLog represents an activity log of an observer.
type ObserverLog struct {
	LogID      string    // Unique identifier for the observer log entry.
	ObserverID string    // ID of the observer.
	Activity   string    // Description of the observer's activity.
	Timestamp  time.Time // Timestamp of the activity.
}

type SystemCheck struct {
    CheckID   string
    Status    string
    Timestamp time.Time
}

type DiagnosticTest struct {
    TestID    string
    Status    string
    Timestamp time.Time
}

type RebootSchedule struct {
    RebootID  string
    Scheduled time.Time
    Status    string
}

type StorageOptimization struct {
    OptimizationID string
    Timestamp      time.Time
    Status         string
}

type DiskHealth struct {
    DiskID     string
    Status     string
    Timestamp  time.Time
}

type Backup struct {
    BackupID       string
    EncryptedData  []byte
    Timestamp      time.Time
    IntegrityCheck bool
}

type HardwareStatus struct {
    ComponentID string
    Status      string
    Timestamp   time.Time
}

type TemporaryFileRecord struct {
    FileID     string
    Status     string
    Timestamp  time.Time
}

type DatabaseCleanupRecord struct {
    CleanupID  string
    Status     string
    Timestamp  time.Time
}

type DefragmentationRecord struct {
    DefragID   string
    Status     string
    Timestamp  time.Time
}

type MaintenanceSchedule struct {
    MaintenanceID string
    ScheduledTime time.Time
    Status        string
}

type SystemHealthStatus struct {
    HealthID  string
    Status    string
    Timestamp time.Time
}

type SystemUpdateRecord struct {
    UpdateID  string
    Timestamp time.Time
    Status    string
}

type ServiceStatus struct {
    StatusID   string
    Status     string
    Timestamp  time.Time
}

type ConfigurationValidation struct {
    ValidationID string
    Status       string
    Timestamp    time.Time
}

type FirmwareUpdate struct {
    UpdateID    string
    Scheduled   time.Time
    Status      string
}

type CPUHealth struct {
    HealthID   string
    Status     string
    Timestamp  time.Time
}

type EncryptionKeyUpdate struct {
    KeyID      string
    NewKey     string
    Timestamp  time.Time
}

type BackupFrequency struct {
    FrequencyID string
    Timestamp   time.Time
}

type ErrorLog struct {
    LogID      string
    Timestamp  time.Time
}

type SecurityCheck struct {
    CheckID    string
    Timestamp  time.Time
}

type NetworkRouteValidation struct {
    ValidationID string
    Timestamp    time.Time
}

type RedundantSystemTest struct {
    TestID    string
    Timestamp time.Time
    Status    string
}

type DataMigrationSchedule struct {
    MigrationID string
    Scheduled   time.Time
    Status      string
}

type ResourceConsumption struct {
    ConsumptionID string
    Timestamp     time.Time
}

type SystemUptime struct {
    UptimeID  string
    Timestamp time.Time
}

type SystemAlert struct {
    AlertID   string
    Timestamp time.Time
}

type SystemSnapshot struct {
    SnapshotID string
    Timestamp  time.Time
}

type ActivityLog struct {
    LogID      string
    Timestamp  time.Time
}

type StressTest struct {
    TestID    string
    Timestamp time.Time
    Result    string
}

type MaintenanceHistory struct {
    EventID   string
    Timestamp time.Time
}

type EnergyConsumption struct {
    ConsumptionID string
    Timestamp     time.Time
}

type NodeSyncStatus struct {
    Status      string
    Timestamp   time.Time
}

type FailoverEvent struct {
    EventID    string
    Timestamp  time.Time
}

type LatencyChange struct {
    Latency    float64
    Timestamp  time.Time
}

type BandwidthUsage struct {
    Bandwidth  float64
    Timestamp  time.Time
}

type ConfigUpdate struct {
    UpdateID   string
    Timestamp  time.Time
}

type DatabaseConnectionStatus struct {
    Status     string
    Timestamp  time.Time
}

type ConsistencyCheck struct {
    Consistency string
    Timestamp   time.Time
}

type ResourceLimit struct {
    Usage      float64
    Timestamp  time.Time
}

type UpdateSchedule struct {
    ScheduleID string
    Scheduled  time.Time
}

type MaintenanceWindow struct {
    WindowID   string
    StartTime  time.Time
    EndTime    time.Time
}

type EncryptionCompliance struct {
    Compliance string
    Timestamp  time.Time
}

type LicenseCompliance struct {
    Compliance string
    Timestamp  time.Time
}

type FirmwareCheck struct {
    Compliance string
    Timestamp  time.Time
}

type StorageUtilization struct {
    Usage      float64
    Timestamp  time.Time
}

type TransactionLoad struct {
    Load       int
    Timestamp  time.Time
}

type RetentionPolicyCompliance struct {
    Compliance string
    Timestamp  time.Time
}

type AntiVirusScanResult struct {
    Results    string
    Timestamp  time.Time
}

type NetworkConfigUpdate struct {
    UpdateID   string
    Timestamp  time.Time
}

type CompressionCompliance struct {
    Compliance string
    Timestamp  time.Time
}

type ProcessExecution struct {
    ExecutionID string
    Timestamp   time.Time
}

type MemoryUsage struct {
    Usage      float64
    Timestamp  time.Time
}

type MemoryCleanupSchedule struct {
    CleanupID  string
    Scheduled  time.Time
}

type EventQueueStatus struct {
    Status     string
    Timestamp  time.Time
}

type HighAvailabilityTest struct {
    Result    string
    Timestamp time.Time
}

type ReplicationCompliance struct {
    Compliance string
    Timestamp  time.Time
}

type SnapshotStatus struct {
    Status    string
    Timestamp time.Time
}

type ProcessLifecycle struct {
    EventID   string
    Timestamp time.Time
}

type NodeRedundancy struct {
    Compliance string
    Timestamp  time.Time
}

type NodeFailoverValidation struct {
    Result    string
    Timestamp time.Time
}

type HeartbeatCheck struct {
    Status    string
    Timestamp time.Time
}

type AlertQueueStatus struct {
    Status    string
    Timestamp time.Time
}

type FileIntegrity struct {
    Status    string
    Timestamp time.Time
}

type FirmwareCompliance struct {
    Compliance string
    Timestamp  time.Time
}

type DataValidation struct {
    ValidationResult string
    Timestamp        time.Time
}

type DisasterRecoverySetup struct {
    SetupID   string
    Timestamp time.Time
}

type SystemShutdown struct {
    Timestamp time.Time
}

type AuditTrailCompliance struct {
    Compliance string
    Timestamp  time.Time
}

type SmartContractStatus struct {
    Status    string
    Timestamp time.Time
}

type DatabaseBackupSchedule struct {
    BackupTime time.Time
}

type ServiceAvailability struct {
    Availability string
    Timestamp    time.Time
}

type ServiceFailures struct {
    Failures  []string
    Timestamp time.Time
}

type LicenseUsage struct {
    Usage     string
    Timestamp time.Time
}

type DataRetentionCompliance struct {
    Compliance string
    Timestamp  time.Time
}

type NetworkTopology struct {
    Topology  string
    Timestamp time.Time
}

type FileLockStatus struct {
    Status    string
    Timestamp time.Time
}

type SmartContractIntegrity struct {
    Integrity string
    Timestamp time.Time
}

type CompressionRatio struct {
    Ratio     float64
    Timestamp time.Time
}

type NetworkDiagnostics struct {
    Diagnostics string
    Timestamp   time.Time
}

type HealthCheckRoutine struct {
    Interval  time.Duration
    Timestamp time.Time
}

type BandwidthTestSchedule struct {
    TestTime time.Time
}

type DataRedundancyStatus struct {
    Status    string
    Timestamp time.Time
}

type ConfigurationChange struct {
    Changes   string
    Timestamp time.Time
}

type SystemUpdateReapplication struct {
    Timestamp time.Time
}

type LogRotationStatus struct {
    Status    string
    Timestamp time.Time
}

type BackupValidation struct {
    Status    string
    Timestamp time.Time
}

type EncryptionUpdateSchedule struct {
    UpdateTime time.Time
}

type ErrorCorrection struct {
    Corrections string
    Timestamp   time.Time
}

type NodeMemoryHealth struct {
    Status    string
    Timestamp time.Time
}

type BackupScheduleValidation struct {
    Status    string
    Timestamp time.Time
}

type SecurityPatchStatus struct {
    Status    string
    Timestamp time.Time
}

type NodeConnectivityStatus struct {
    Status    string
    Timestamp time.Time
}

type APIComplianceStatus struct {
    Compliance string
    Timestamp  time.Time
}

type APIEndpointHealth struct {
    EndpointStatus string
    Timestamp      time.Time
}

type SystemAuditReport struct {
    Report    string
    Timestamp time.Time
}

type DatabaseRebuildSchedule struct {
    RebuildTime time.Time
}

type DataIngestionValidation struct {
    Status    string
    Timestamp time.Time
}

type NodePerformanceMetrics struct {
    Metrics   string
    Timestamp time.Time
}

type SoftwareComplianceStatus struct {
    Compliance string
    Timestamp  time.Time
}

type NodeReinitialization struct {
    Timestamp time.Time
}

type SystemHardening struct {
    Timestamp time.Time
}

type FirewallValidation struct {
    RuleStatus string
    Timestamp  time.Time
}

type SystemWarmupStatus struct {
    Status    string
    Timestamp time.Time
}

type RedundancyValidation struct {
    Status    string
    Timestamp time.Time
}

type SecurityAuditFrequency struct {
    Frequency time.Duration
    Timestamp time.Time
}

type DatabaseTransactionValidation struct {
    Status    string
    Timestamp time.Time
}

type SecurityScanSchedule struct {
    ScanTime time.Time
}

type DataLossPreventionStatus struct {
    Status    string
    Timestamp time.Time
}

type SystemRollback struct {
    Details   string
    Timestamp time.Time
}

type CloudBackupStatus struct {
    Status    string
    Timestamp time.Time
}

type ConnectionPoolValidation struct {
    Status    string
    Timestamp time.Time
}

type SmartContractLoad struct {
    Load      string
    Timestamp time.Time
}

type ResourceDeallocation struct {
    Deallocations string
    Timestamp     time.Time
}

type PermissionIntegrity struct {
    IntegrityStatus string
    Timestamp       time.Time
}

type DataCleanupFrequency struct {
    Frequency string
    Timestamp time.Time
}

type ConfigurationDrift struct {
    Drifts    string
    Timestamp time.Time
}

type LogSizeLimits struct {
    Size       string
    Timestamp  time.Time
}

type SessionPersistenceStatus struct {
    Status     string
    Timestamp  time.Time
}

type ApplicationUpdate struct {
    Timestamp  time.Time
}

type RoleAssignmentChanges struct {
    Changes    string
    Timestamp  time.Time
}

type NodeUpdateStatus struct {
    Status     string
    Timestamp  time.Time
}

type APIComplianceStatus struct {
    Compliance string
    Timestamp  time.Time
}

type LogAccessAttempts struct {
    Attempts   string
    Timestamp  time.Time
}

type APIRateLimits struct {
    Status     string
    Timestamp  time.Time
}

type NodeRebootSchedule struct {
    RebootTime time.Time
}

type TokenDistribution struct {
    Distribution string
    Timestamp    time.Time
}

type SystemSelfRepair struct {
    Status     string
    Timestamp  time.Time
}

// PerformanceLog struct for logging system performance metrics
type PerformanceLog struct {
    Metric    string
    Value     string
    Timestamp time.Time
}

// OptimizationSetting struct for managing system optimization levels
type OptimizationSetting struct {
    Level     int
    Timestamp time.Time
}

// PerformanceLog struct for logging system performance metrics
type PerformanceLog struct {
    Metric    string
    Value     string
    Timestamp time.Time
}

// ResourceLimit struct for setting resource usage limits
type ResourceLimit struct {
    Resource  string
    Limit     float64
    Timestamp time.Time
}

// SystemLoad struct for monitoring overall system resource usage
type SystemLoad struct {
    CPUUsage    float64
    MemoryUsage float64
    DiskUsage   float64
}

// DiskCacheConfig struct for configuring disk cache size
type DiskCacheConfig struct {
    SizeMB    int
    Timestamp time.Time
}

// ResourceSharingConfig struct for managing dynamic resource scaling
type ResourceSharingConfig struct {
    DynamicScalingEnabled bool
    Timestamp             time.Time
}

// NetworkConfig struct for managing network bandwidth configuration
type NetworkConfig struct {
    BandwidthLimitMBps float64
    Timestamp          time.Time
}

// HealthMetrics struct for monitoring blockchain health
type HealthMetrics struct {
    NodeCount       int
    ActiveNodes     int
    TransactionRate float64
    AvgLatency      float64
    Timestamp       time.Time
}

// ErrorLog struct for storing error logs
type ErrorLog struct {
    ErrorMessage string
    Timestamp    time.Time
}

// Alert struct for managing performance alerts
type Alert struct {
    Metric    string
    Threshold float64
    Active    bool
    Timestamp time.Time
}

// CompressionRateLog struct for tracking compression rates
type CompressionRateLog struct {
    Rate      float64
    Timestamp time.Time
}

// FileTransferStatusLog struct for tracking file transfer statuses
type FileTransferStatusLog struct {
    Status    string
    Timestamp time.Time
}

// KeyRotationLog struct for tracking encryption key rotations
type KeyRotationLog struct {
    Status    string
    Timestamp time.Time
}

// HardwareStatusLog struct for tracking hardware health
type HardwareStatusLog struct {
    Status    string
    Timestamp time.Time
}

// SessionDurationLog struct for logging user session durations
type SessionDurationLog struct {
    SessionDurations map[string]time.Duration
    Timestamp        time.Time
}

// RBACLog struct for role-based access control monitoring
type RBACLog struct {
    Status    string
    Timestamp time.Time
}

// LogIntegrityLog struct for tracking log integrity
type LogIntegrityLog struct {
    Integrity bool
    Timestamp time.Time
}

// MultiFactorAuthStatusLog struct for tracking MFA status
type MultiFactorAuthStatusLog struct {
    Status    string
    Timestamp time.Time
}

// TokenUsageLog struct for tracking token usage
type TokenUsageLog struct {
    UsageMetrics map[string]int
    Timestamp    time.Time
}

// ConsensusEfficiencyLog struct for monitoring consensus efficiency
type ConsensusEfficiencyLog struct {
    Efficiency float64
    Timestamp  time.Time
}

// AlertResponseTimeLog struct for logging alert response times
type AlertResponseTimeLog struct {
    ResponseTimes map[string]time.Duration
    Timestamp     time.Time
}

// UserPermissionsStatusLog struct for tracking user permissions compliance
type UserPermissionsStatusLog struct {
    Status    string
    Timestamp time.Time
}

// NodeReconnectionLog struct for logging node reconnection events
type NodeReconnectionLog struct {
    ReconnectionCount int
    Timestamp         time.Time
}

// DataAccessPatternLog struct for monitoring data access patterns
type DataAccessPatternLog struct {
    Patterns  map[string]interface{}
    Timestamp time.Time
}

// TransactionVolumeLog struct for tracking transaction volumes
type TransactionVolumeLog struct {
    Volume    int
    Timestamp time.Time
}

// ContractExecutionLog struct for logging contract execution metrics
type ContractExecutionLog struct {
    Metrics   map[string]interface{}
    Timestamp time.Time
}

// FunctionExecutionTimeLog struct for tracking function execution times
type FunctionExecutionTimeLog struct {
    ExecutionTimes map[string]time.Duration
    Timestamp      time.Time
}

// APICallVolumeLog struct for monitoring API call volume
type APICallVolumeLog struct {
    Volume    int
    Timestamp time.Time
}

// ResourceUsageTrendLog struct for tracking resource usage trends
type ResourceUsageTrendLog struct {
    Trends    map[string]interface{}
    Timestamp time.Time
}

// SecurityPatchStatusLog struct for monitoring security patch statuses
type SecurityPatchStatusLog struct {
    Status    string
    Timestamp time.Time
}

// PerformanceSummary for blockchain performance metrics
type PerformanceSummary struct {
    AvgCPUUsage       float64
    AvgMemoryUsage    float64
    AvgDiskIO         float64
    NetworkThroughput float64
    TotalTransactions int
    AverageLatency    float64
    Timestamp         time.Time
}

// ResourceReport for resource utilization and efficiency
type ResourceReport struct {
    CPUUtilization       float64
    MemoryUtilization    float64
    DiskSpaceUtilization float64
    NetworkUsage         float64
    EnergyConsumption    float64
    EncryptedData        string
    Timestamp            time.Time
}

// ResourceReallocation for logging resource reallocation details
type ResourceReallocation struct {
    ResourceType       string
    AmountReallocated  float64
    Reason             string
    EncryptedData      string
    Timestamp          time.Time
}

// CostReport for logging resource allocation costs
type CostReport struct {
    ResourceType       string
    AllocationDuration time.Duration
    RatePerUnit        float64
    TotalCost          float64
    EncryptedData      string
    Timestamp          time.Time
}

type FirmwareStatus struct {
    ComplianceDetails string
    Timestamp         time.Time
}

type RoleChange struct {
    UserID    string
    OldRole   string
    NewRole   string
    Timestamp time.Time
}

type NodeReputation struct {
    NodeID        string
    ReputationScore float64
    Timestamp      time.Time
}

type AccessViolation struct {
    ViolationDetails string
    Timestamp        time.Time
}

type IntrusionAttempt struct {
    AttemptDetails string
    Timestamp      time.Time
}

type ProtocolCompliance struct {
    ComplianceDetails string
    Timestamp         time.Time
}

type ThreatLevel struct {
    Level     string
    Timestamp time.Time
}

type RetentionCompliance struct {
    ComplianceDetails string
    Timestamp         time.Time
}

type TrafficVolume struct {
    Volume    float64
    Timestamp time.Time
}

type BandwidthUsage struct {
    Usage     float64
    Timestamp time.Time
}

type NodeMigration struct {
    NodeID     string
    MigrationDetails string
    Timestamp   time.Time
}

type ServiceResponseTime struct {
    ServiceName   string
    ResponseTime  float64
    Timestamp     time.Time
}

type UserLoginAttempt struct {
    UserID       string
    IPAddress    string
    Timestamp    time.Time
    IsSuccessful bool
}

type ComplianceAuditResult struct {
    AuditDetails string
    Timestamp    time.Time
}

type BlockchainUpdate struct {
    Version       string
    UpdateDetails string
    Timestamp     time.Time
}

type EnergyConsumption struct {
    Amount    float64
    Timestamp time.Time
}

type NodeFailureRate struct {
    NodeID       string
    FailureCount int
    Timestamp    time.Time
}

type APIThrottleLimit struct {
    Endpoint    string
    Limit       float64
    CurrentUsage float64
    Timestamp   time.Time
}

type DatabaseHealth struct {
    Status     string
    Details    string
    Timestamp  time.Time
}

type SystemConfigurationChange struct {
    ConfigName   string
    OldValue     string
    NewValue     string
    ChangedBy    string
    Timestamp    time.Time
}

type CacheUsage struct {
    CacheType  string
    Usage      float64
    Timestamp  time.Time
}

type APIUsage struct {
    Endpoint   string
    Calls      int
    Timestamp  time.Time
}

type SessionTimeout struct {
    SessionID  string
    Duration   time.Duration
    Timestamp  time.Time
}

type AccessFrequency struct {
    UserID     string
    Frequency  int
    Timestamp  time.Time
}

type RateLimitCompliance struct {
    Endpoint      string
    Compliant     bool
    Timestamp     time.Time
}

type ThreatDetection struct {
    ThreatType    string
    DetectedAt    time.Time
    Severity      string
}

type AlertStatus struct {
    AlertType  string
    Active     bool
    Timestamp  time.Time
}

type AnomalyDetection struct {
    AnomalyType  string
    Details      string
    Timestamp    time.Time
}

type EventFrequency struct {
    EventName  string
    Frequency  int
    Timestamp  time.Time
}


type DataTransferRate struct {
    NodeID    string
    Rate      float64
    Timestamp time.Time
}

type DataRetrievalTime struct {
    RetrievalID string
    TimeTaken   float64
    Timestamp   time.Time
}

type TransactionLatency struct {
    TransactionID string
    Latency       float64
    Timestamp     time.Time
}

type StorageQuotaUsage struct {
    UserID   string
    QuotaUsed float64
    Timestamp time.Time
}

type DiskSpeed struct {
    ReadSpeed  float64
    WriteSpeed float64
    Timestamp  time.Time
}

type NetworkResilience struct {
    Metric      string
    Resilience  float64
    Timestamp   time.Time
}

type BlockchainIntegrity struct {
    Status     string
    Timestamp  time.Time
}

type EncryptionCompliance struct {
    Compliant   bool
    Timestamp   time.Time
}

type SessionActivity struct {
    SessionID string
    Details   string
    Timestamp time.Time
}

type AccessControlStatus struct {
    Status     string
    Timestamp  time.Time
}

type SystemHealth struct {
    Status    string
    Timestamp time.Time
}

type NodeStatus struct {
    NodeID    string
    Status    string
    Timestamp time.Time
}

type ResourceUsage struct {
    ResourceType string
    Usage        float64
    Timestamp    time.Time
}

type NetworkLatency struct {
    Latency   float64
    Timestamp time.Time
}

type DataThroughput struct {
    Throughput float64
    Timestamp  time.Time
}

type TransactionRate struct {
    Rate      float64
    Timestamp time.Time
}

type BlockPropagationTime struct {
    Time      float64
    Timestamp time.Time
}

type ConsensusStatus struct {
    Status    string
    Timestamp time.Time
}

type SubBlockValidation struct {
    ValidationID string
    Status       string
    Timestamp    time.Time
}

type SubBlockCompletion struct {
    BlockID    string
    Completion float64
    Timestamp  time.Time
}

type PeerConnectionStatus struct {
    PeerID    string
    Status    string
    Timestamp time.Time
}

type DataSyncStatus struct {
    Status    string
    Timestamp time.Time
}

type NodeAvailability struct {
    NodeID    string
    Available bool
    Timestamp time.Time
}

type ShardHealth struct {
    ShardID   string
    Health    string
    Timestamp time.Time
}

type DiskUsage struct {
    TotalSpace  float64
    UsedSpace   float64
    FreeSpace   float64
    Timestamp   time.Time
}

type MemoryUsage struct {
    TotalMemory float64
    UsedMemory  float64
    FreeMemory  float64
    Timestamp   time.Time
}

type CPUUtilization struct {
    Usage      float64
    Timestamp  time.Time
}

type NodeDowntime struct {
    NodeID    string
    Downtime  float64
    Timestamp time.Time
}

type NetworkBandwidth struct {
    Bandwidth float64
    Timestamp time.Time
}

type ErrorRate struct {
    Rate      float64
    Timestamp time.Time
}

type UserActivity struct {
    Activity  string
    Timestamp time.Time
}

type ComplianceStatus struct {
    Status    string
    Timestamp time.Time
}

type AuditLog struct {
    Logs      string
    Timestamp time.Time
}

type ThreatResponseTime struct {
    ResponseTime float64
    Timestamp    time.Time
}

type SystemUptime struct {
    Uptime    time.Duration
    Timestamp time.Time
}

type TrafficPattern struct {
    Pattern   string
    Timestamp time.Time
}

type SuspiciousActivity struct {
    Activity  string
    Timestamp time.Time
}

type LoadBalancingStatus struct {
    Status    string
    Timestamp time.Time
}

type HealthThreshold struct {
    Component  string
    Threshold  float64
    Timestamp  time.Time
}

type IncidentResponseTime struct {
    ResponseTime float64
    Timestamp    time.Time
}

type APIResponseTime struct {
    ResponseTime float64
    Timestamp    time.Time
}

type DataRequestVolume struct {
    RequestVolume int
    Timestamp     time.Time
}

type SessionDataUsage struct {
    DataUsage   float64
    Timestamp   time.Time
}

type RateLimitExceedance struct {
    Exceedances int
    Timestamp   time.Time
}

type EventLog struct {
    Logs      string
    Timestamp time.Time
}

type SystemAlert struct {
    Alert     string
    Timestamp time.Time
}

type ResourceAllocation struct {
    Resource  string
    Status    string
    Timestamp time.Time
}

type EncryptionStatus struct {
    Status    string
    Timestamp time.Time
}

type ConsensusAnomaly struct {
    Anomaly   string
    Timestamp time.Time
}

type SecurityPolicyCompliance struct {
    ComplianceStatus string
    Timestamp        time.Time
}

type OptimizationPolicy struct {
    CachingEnabled bool
    Timestamp      time.Time
}

type PerformanceLog struct {
    Metric    string
    Value     string // Encrypted value
    Timestamp time.Time
}

type SystemOverhead struct {
    CPUOverhead       float64
    MemoryOverhead    float64
    DiskOverhead      float64
    EventDescription  string // For logging specific events
    Timestamp         time.Time
}

type PriorityMode struct {
    Mode      string // High, Medium, or Low
    Timestamp time.Time
}

type OptimizationPolicy struct {
    CachingEnabled        bool
    AutoAdjustmentEnabled bool
    PriorityMode          string // High, Medium, or Low
    Timestamp             time.Time
}

type PerformanceLog struct {
    Metric    string
    Value     string // Encrypted value
    Timestamp time.Time
}

type ResourceConsumption struct {
    CPUUsage     float64
    MemoryUsage  float64
    DiskUsage    float64
    NetworkUsage float64
    Timestamp    time.Time
}

type UtilizationRates struct {
    CPUUtilization    float64
    MemoryUtilization float64
    DiskUtilization   float64
    Timestamp         time.Time
}

type ThreadPoolConfig struct {
    MaxSize   int
    Timestamp time.Time
}

type ResourceAlert struct {
    Metric    string
    Threshold float64
    Active    bool
    Timestamp time.Time
}

type PerformanceGoal struct {
    Metric    string
    Target    float64
    Timestamp time.Time
}

type UptimeLog struct {
    Uptime    time.Duration
    Timestamp time.Time
}

type PerformanceLog struct {
    Metric    string
    Value     string
    Timestamp time.Time
}

type UsageStats struct {
    CPUUsage    float64
    MemoryUsage float64
    DiskUsage   float64
    NetworkLoad float64
    Timestamp   time.Time
}

type ScalingEvent struct {
    Description   string
    ScalingFactor float64
    Timestamp     time.Time
}

type ResourceAlert struct {
    Metric    string
    Threshold float64
    Active    bool
    Timestamp time.Time
}



// ************** Network Structs **************

// PeerInfo represents the information about a peer in the network.
type PeerInfo struct {
	PeerID    string // Unique identifier for the peer
	Address   string // Peer address
	PublicKey string // Public key of the peer
}

// BlockRequest represents a request for retrieving a block by index
type BlockRequest struct {
	BlockIndex int `json:"block_index"`
}

// RoutingEvent represents a routing-related event in the network.
type RoutingEvent struct {
	RouteID    string    // Unique identifier for the route
	SourceNode string    // Source node of the route
	DestNode   string    // Destination node of the route
	Timestamp  time.Time // When the event occurred
}

// Define the WebRTCConn struct with an EncryptionKey field
type WebRTCConn struct {
	ConnectionID  string
	EncryptionKey []byte // Added EncryptionKey field
}

// Define WebRTCConnection struct in the ledger package (if needed)
type WebRTCConnection struct {
	ConnectionID string
	PeerID       string
	Timestamp    time.Time // Renamed from "ConnectedAt" to "Timestamp"
}

// Node represents a network node for monitoring and communication.
type Node struct {
	ID           string       // Unique identifier for the node
	Address      string       // URL or address of the node
	AuthToken    string       // Authorization token for secure communication
	Name         string       // Name of the node
	NodeCategory NodeCategory // Node's category
	NodeType     NodeType     // Node's type
	NodeKey      *NodeKey     // Node's public key
	IsActive     bool         // Node activity status
	Endpoint     string       // Full endpoint URL for communication
}

// NodeCategory represents the category of the node (e.g., authority, standard).
type NodeCategory string

const (
	AuthorityCategory NodeCategory = "Authority"
	StandardCategory  NodeCategory = "Standard"
)

// BootstrapNode represents the node responsible for bootstrapping the network
type BootstrapNode struct {
	ActiveNodes    map[string]NodeInfo // Map of active nodes in the network
	mutex          sync.Mutex          // Mutex for safe concurrent access
	ListenAddress  string              // The IP address and port to listen on
	LedgerInstance *Ledger             // Pointer to the ledger for recording network activities
	NodeType       string              // Type of node (e.g., Bootstrap, Validator)
}

// ConnectionPool manages reusable connections between nodes
type ConnectionPool struct {
	connections map[string]net.Conn // Active network connections keyed by node ID
	mutex       sync.Mutex          // Mutex for thread-safe access
	maxIdleTime time.Duration       // Maximum idle time before closing a connection
	ActiveConns []net.Conn          // Use a list of net.Conn to store active connections
}

// NodeManager manages node operations, roles, health, and configurations within the network.
type NodeManager struct {
	NodeID         string                 // Unique identifier for the node.
	NodeType       string                 // Type of node (e.g., full, light, archival, validator).
	Status         string                 // Current status of the node (e.g., active, suspended, decommissioned).
	Role           string                 // Role assigned to the node (e.g., validator, observer).
	HealthScore    int                    // Health score of the node (0-100).
	Configurations map[string]interface{} // Configuration settings specific to the node.
	LastActive     time.Time              // Timestamp of the last activity of the node.
	RegisteredAt   time.Time              // Timestamp when the node was registered in the network.
	Metrics        NodeMetrics            // Performance metrics of the node.
	Dependencies   []string               // Dependencies on other nodes or components.
}

// NodeConfiguration represents a detailed configuration structure for nodes.
type NodeConfiguration struct {
	MaxConnections int                    // Maximum allowed connections to the node.
	SecurityConfig map[string]interface{} // Security-related configurations (e.g., firewalls, TLS).
	ResourceLimits map[string]interface{} // Resource limits (e.g., CPU, memory, disk usage).
	CustomSettings map[string]interface{} // Custom settings specific to the node.
}

// NodeHealthReport provides a snapshot of the node's health and status.
type NodeHealthReport struct {
	NodeID          string    // Unique identifier for the node.
	HealthScore     int       // Current health score of the node.
	LastCheck       time.Time // Timestamp of the last health check.
	Status          string    // Status of the node (e.g., healthy, degraded).
	Recommendations []string  // Suggested actions to improve node health.
}

// DistributedNetworkCoordinator manages the coordination of nodes in a decentralized network
type DistributedNetworkCoordinator struct {
	Nodes          map[string]*NodeInfo       // All nodes participating in the network
	ActivePeers    map[string]*PeerConnection // Active peer connections
	mutex          sync.Mutex                 // For thread-safe access
	ledgerInstance *Ledger                    // Pointer to the ledger for logging and auditing
	Blockchain     *Blockchain                // Pointer to the blockchain (add this field)
	Consensus      *SynnergyConsensus
}

// Peer represents a connected peer in the network
type Peer struct {
	Address   string    // Peer address (IP/Port)
	Conn      net.Conn  // Peer connection
	PublicKey string    // Public key of the peer (as a string, representing the AES key)
	LastSeen  time.Time // Last time the peer was seen (used for timeouts)
}

// P2PNetwork manages the peer-to-peer network and communication
type P2PNetwork struct {
	Address          string             // The node's own address
	MessageQueue     chan P2PMessage    // Queue to store outgoing messages
	PendingBlocks    []Block            // Blocks pending validation
	Peers            map[string]*Peer   // Connected peers
	IncomingMessages chan P2PMessage    // Incoming message channel
	OutgoingMessages chan P2PMessage    // Outgoing message channel
	PeerLock         sync.Mutex         // Synchronizes access to peers
	NodeKey          *NodeKey           // Node's public-private keypair for encryption
	LedgerInstance   *Ledger            // Pointer to the ledger for transaction and block management
	ConsensusEngine  *SynnergyConsensus // Pointer to the Synnergy Consensus engine
}

// P2PMessage represents the structure of a peer-to-peer message
type P2PMessage struct {
	Sender       string    // Sender's public key (encoded as string)
	Recipient    string    // Recipient peer address
	EncryptedKey string    // Base64 encoded AES key
	Content      string    // Base64 encoded encrypted message content
	Timestamp    time.Time // Timestamp of the message
}

// PeerAdvertisement contains information about a peer that is being advertised in the network
type PeerAdvertisement struct {
	PeerAddress string    // IP address or domain of the peer
	PublicKey   string    // Public key of the peer for encryption
	Timestamp   time.Time // Timestamp when the advertisement was created
	Signature   string    // Encrypted signature to ensure authenticity
}

// PeerAdvertiser handles advertising and discovering peers in the network
type PeerAdvertiser struct {
	Peers         map[string]PeerAdvertisement // Known peers in the network
	Mutex         sync.Mutex                   // Mutex for thread-safe access
	Network       *P2PNetwork                  // Pointer to the P2P network
	AdvertiseFreq time.Duration                // Frequency to broadcast peer advertisements
}

// PeerConnectionManager manages all active peer connections in the network
type PeerConnectionManager struct {
	Connections map[string]*PeerConnection // Active peer connections
	Mutex       sync.Mutex                 // Thread-safe access to peer connections
	Ledger      *Ledger                    // Pointer to the ledger for recording transactions
}

// PeerConnection represents a connection to a peer node, using a ConnectionPool
type PeerConnection struct {
	NodeID        string
	Connection    *ConnectionPool // Custom connection pool instead of net.Conn
	LastPingTime  time.Time       // Last time the peer was pinged
	IsAlive       bool            // Status of the connection
	EncryptionKey []byte          // Encryption key used for secure communication
}

// Client represents an RPC client for interacting with the Synnergy Network
type RPCClient struct {
	ServerAddress string
	PublicKey     string
	PrivateKey    string
}

// NetworkRequest represents an incoming network request
type NetworkRequest struct {
	Type    string          `json:"type"`    // The type of the request (e.g., "TRANSACTION_SUBMISSION")
	Data    []byte          `json:"data"`    // The raw data (e.g., serialized transaction)
	Payload json.RawMessage `json:"payload"` // The raw payload for further processing

}

// FaultToleranceManager manages the fault tolerance mechanisms in the network
type FaultToleranceManager struct {
	Nodes           []string          // List of nodes in the network
	NodeStatus      map[string]bool   // Status of each node (alive or down)
	QuorumThreshold int               // Minimum number of nodes required for consensus
	ledger          *Ledger           // Ledger for logging fault tolerance events
	mutex           sync.Mutex        // Mutex for thread-safe operations
	NodeState       map[string]string // Map of node to its latest synced block hash
}

// Firewall manages and filters network traffic, securing the network from malicious actors
type Firewall struct {
	allowedIPs     map[string]bool      // A whitelist of allowed IPs
	blockedIPs     map[string]time.Time // A list of blocked IPs and the time they were blocked
	blockDuration  time.Duration        // Duration for which an IP will remain blocked
	mutex          sync.Mutex           // Mutex for thread-safe operations
	ledgerInstance *Ledger              // Pointer to the ledger for logging firewall events
}

// FirewallManager manages multiple firewall types including normal, dynamic, stateless, and stateful firewalls
type FirewallManager struct {
	NormalFirewall    *Firewall          // Normal firewall (whitelisting/blocking)
	DynamicFirewall   *DynamicFirewall   // Dynamic firewall (adjusts based on traffic patterns)
	StatelessFirewall *StatelessFirewall // Stateless firewall (simple packet-filtering firewall)
	StatefulFirewall  *StatefulFirewall  // Stateful firewall (tracks connection state)
	ledgerInstance    *Ledger
}

// DynamicFirewall adjusts its rules dynamically based on traffic patterns
type DynamicFirewall struct {
	allowedIPs        map[string]bool
	blockedIPs        map[string]time.Time
	blockDuration     time.Duration
	ledgerInstance    *Ledger
	encryptionService *Encryption
	mutex             sync.Mutex
}

// StatelessFirewall represents a simple firewall that filters based on rules without tracking connection states
type StatelessFirewall struct {
	allowedPorts      []int
	blockedIPs        map[string]time.Time
	ledgerInstance    *Ledger
	encryptionService *Encryption
	mutex             sync.Mutex
}

// StatefulFirewall keeps track of connections and the state of each connection (open, closed, etc.)
type StatefulFirewall struct {
	allowedConnections map[string]string // IP -> State (open/closed)
	blockedIPs         map[string]time.Time
	ledgerInstance     *Ledger
	encryptionService  *Encryption
	mutex              sync.Mutex
}

// FlowControlManager handles network traffic, ensuring smooth flow of transactions and blocks
type FlowControlManager struct {
	MaxPendingTransactions int           // Max number of pending transactions allowed
	MaxSubBlockSize        int           // Max number of transactions per sub-block
	MaxBlockSize           int           // Max number of sub-blocks per block
	PendingTransactions    []Transaction // Pool of pending transactions
	PendingSubBlocks       []SubBlock    // Pool of pending sub-blocks waiting to be added to a block
	mutex                  sync.Mutex    // Mutex for thread-safe operations
	ledgerInstance         *Ledger       // Pointer to the ledger for tracking flow control events
}

// GeoLocation represents the latitude and longitude of a node
type GeoLocation struct {
	Latitude  float64
	Longitude float64
}

// GeoLocationManager manages the geolocation data for nodes in the network
type GeoLocationManager struct {
	NodeLocations map[string]GeoLocation // Maps node IDs to their geolocation
	ledger        *Ledger                // Reference to the ledger for logging
}

// Handshake represents the handshake mechanism for secure communication between nodes
type Handshake struct {
	PrivateKey *rsa.PrivateKey // The private key of the node
	PublicKey  *rsa.PublicKey  // The public key of the node
	mutex      sync.Mutex      // Mutex for thread-safe operations
	ledger     *Ledger         // Pointer to the ledger for logging handshake events
}

// KademliaNode represents a node in the Kademlia DHT
type KademliaNode struct {
	ID         string      // Unique ID of the node
	IPAddress  string      // Node's IP address
	LastActive time.Time   // Last active timestamp
	Location   GeoLocation // Geolocation of the node
}

// KademliaDHT represents the Kademlia Distributed Hash Table
type KademliaDHT struct {
	NodeID  string                  // ID of the local node
	KBucket map[string]KademliaNode // Kademlia node bucket (ID to node mapping)
	ledger  *Ledger                 // Reference to ledger for logging
	lock    sync.Mutex              // Mutex for thread-safe operations
}

// MessageQueue represents the queue of messages waiting to be processed
type MessageQueue struct {
	queue          *MessageList // Doubly-linked list for efficient message queuing
	maxQueueSize   int          // Maximum size of the message queue
	lock           sync.Mutex   // Mutex to ensure thread-safe operations
	ledgerInstance *Ledger      // Reference to the ledger for logging message events
}

// Message represents a message between nodes in the network
type Message struct {
	ID        string    // Unique message ID
	Timestamp time.Time // Time when the message was created
	From      string    // Node ID of the sender
	To        string    // Node ID of the receiver
	Content   string    // Message content
	Hash      string    // SHA-256 hash of the message content
	Encrypted bool      // Whether the message content is encrypted
}

// NATTraversalManager handles NAT traversal tasks for nodes behind routers
type NATTraversalManager struct {
	publicIP       string            // Public IP address of the node
	privateIP      string            // Private IP address of the node
	peerMap        map[string]string // Map of peers' public IPs to private IPs
	lock           sync.Mutex        // Mutex to ensure thread-safe operations
	ledgerInstance *Ledger           // Reference to the ledger for logging NAT traversal events
	webrtcManager  *WebRTCManager
}

// NetworkManager handles all networking activities, including peer discovery, message passing, and encryption
type NetworkManager struct {
	nodeAddress    string
	peers          map[string]*PeerConnection // Active peers connected to this node
	ledgerInstance *Ledger                    // Reference to the ledger for logging network events
	lock           sync.Mutex                 // Mutex for thread-safe operations
	connectionPool *ConnectionPool            // Pool to manage active connections
	PeerDiscovery  *PeerDiscoveryManager
}

// NodeType defines the possible types of nodes in the network
type NodeType string

const (
	DefaultNodeType   NodeType = "default"   // Default node type (general-purpose node)
	ValidatorNodeType NodeType = "validator" // Node that participates in PoS validation
	MinerNodeType     NodeType = "miner"     // Node that performs PoW mining
)

// Packet represents a network packet that needs to be routed
type Packet struct {
	SourceID      string // ID of the source peer
	DestinationID string // ID of the destination peer
	Data          []byte // Data to be sent in the packet
}

// RPCRequest represents an RPC call request
type RPCRequest struct {
	Method          string `json:"method"`
	Payload         string `json:"payload"`
	SenderPublicKey string `json:"sender_public_key"`
}

// BalanceRequest represents a request to retrieve the balance of a wallet
type BalanceRequest struct {
	WalletAddress string `json:"wallet_address"`
}

// BalanceResponse represents the response containing the wallet balance
type BalanceResponse struct {
	WalletAddress string  `json:"wallet_address"`
	Balance       float64 `json:"balance"`
}

// RPCResponse represents the response from an RPC call
type RPCResponse struct {
	Data string `json:"data"`
}

// PeerKey represents the public key of a peer
type PeerKey []byte

// PeerDiscoveryManager manages the discovery of new peers
type PeerDiscoveryManager struct {
	Peers          map[string]*Peer // List of discovered peers
	Ledger         *Ledger          // Instance of the ledger for validation
	LocalAddress   string           // Local node's address
	LocalPublicKey string           // Local node's public key
	KnownNetworks  []string         // List of known networks for discovering peers
}

// QoSManager manages the quality of service across the P2P network
type QoSManager struct {
	BandwidthLimit int                  // Maximum bandwidth allowed for communication in KBps
	PriorityQueues map[int][]*QoSPacket // Queues for different packet priorities
	Mutex          sync.Mutex           // Mutex to handle concurrency
	LedgerInstance *Ledger              // Instance of the ledger to record QoS-related events
	NodePublicKey  string               // Public key of the local node for encryption
	ConnectedPeers map[string]*Peer     // List of connected peers
}

// QoSPacket represents a network packet with a priority level
type QoSPacket struct {
	Payload     []byte    // The actual data in the packet
	Priority    int       // Priority of the packet (0 = highest, 10 = lowest)
	Timestamp   time.Time // Timestamp when the packet was created
	Destination string    // Destination peer address
}

// Router represents the main routing manager for peer-to-peer communications
type Router struct {
	Routes          map[string]string // Map of node IDs to their IP addresses
	Mutex           sync.Mutex        // Mutex for thread-safe operations
	LedgerInstance  *Ledger           // Ledger instance to store routing events
	EncryptedRoutes map[string][]byte // Encrypted routes for added security
	NodePublicKey   string            // Public key of the local node for encryption
	Peers           map[string]*Peer  // List of active peers with their addresses and public keys
	ConnectionPool  *ConnectionPool   // Connection pool for managing network connections
}

// RPCServer represents the RPC server for Synnergy Network
type RPCServer struct {
	Address        string
	LedgerInstance *Ledger // Ledger integration for storing and retrieving data
	Router         *Router // Router to manage routing and peer connections
}

// SDNController represents the Software-Defined Network controller, managing the network nodes
type SDNController struct {
	Nodes          map[string]*SDNNode // Connected nodes in the network
	NodeLock       sync.Mutex          // Mutex for thread-safe node operations
	LedgerInstance *Ledger             // Integration with ledger for network state tracking
	EncryptionKey  string              // Key for encrypting SDN controller communications
}

// SDNNode represents a node in the SDN-controlled network
type SDNNode struct {
	NodeID    string    // Unique ID of the node
	Address   net.IP    // IP address of the node
	Status    string    // Status (active, inactive)
	LastCheck time.Time // Last heartbeat or health check
}

// Server represents a blockchain server node responsible for handling incoming connections
type Server struct {
	Address         string              // Server address
	LedgerInstance  *Ledger             // Ledger for storing blocks and transactions
	ConsensusEngine *SynnergyConsensus  // Consensus mechanism
	EncryptionKey   string              // Encryption key for secure communication
	connections     map[string]net.Conn // Map of active connections
	connectionLock  sync.Mutex          // Mutex for thread-safe operations
}

// SSLHandshakeManager handles SSL/TLS secure handshakes for encrypted communication
type SSLHandshakeManager struct {
	CertFile    string              // Path to the certificate file
	KeyFile     string              // Path to the private key file
	CAFile      string              // Path to the Certificate Authority (CA) file
	Connections map[string]net.Conn // Map of active connections
	lock        sync.Mutex          // Mutex to ensure thread-safe operations
}

// TLSHandshakeManager manages secure TLS handshakes for encrypted communication between nodes
type TLSHandshakeManager struct {
	CertFile    string              // Path to the TLS certificate file
	KeyFile     string              // Path to the private key file
	CAFile      string              // Path to the CA certificate
	Connections map[string]net.Conn // Active TLS connections
	lock        sync.Mutex          // Mutex to protect concurrent access
	Ledger      *Ledger             // Integration with the ledger for secure transaction handling
}

// TopologyManager manages the network topology and routing for communication between nodes
type TopologyManager struct {
	Nodes          map[string]*Node // A map of all active nodes in the network
	NetworkGraph   *NetworkGraph    // The network graph for representing connections between nodes
	lock           sync.Mutex       // Mutex for thread-safe operations
	LedgerInstance *Ledger          // Integration with the ledger for node data persistence
}

// WebRTCManager handles WebRTC connections for peer-to-peer communication across the Synnergy Network
type WebRTCManager struct {
	Peers          map[string]*PeerConnection // Map of active peer connections
	lock           sync.Mutex                 // Mutex for thread-safe operations
	LedgerInstance *Ledger                    // Ledger for tracking peer connections
}

// NetworkGraph represents the network topology, containing nodes and edges for connections between nodes.
type NetworkGraph struct {
	Nodes map[string]*GraphNode   // Nodes in the network, keyed by node ID
	Edges map[string][]*GraphEdge // Connections between nodes, keyed by node ID (adjacency list)
}

// GraphNode represents a node in the network graph.
type GraphNode struct {
	NodeID      string       // Unique identifier of the node
	NodeInfo    *NodeInfo    // Information about the node
	Connections []*GraphEdge // Edges representing connections to other nodes
}

// GraphEdge represents a connection between two nodes in the network graph.
type GraphEdge struct {
	FromNodeID string  // ID of the source node
	ToNodeID   string  // ID of the destination node
	Weight     float64 // Weight or cost of the connection (can represent latency, bandwidth, etc.)
}

// NodeInfo represents the information related to a network node.
type NodeInfo struct {
	NodeID         string // Unique identifier of the node
	Address        string
	IPAddress      string      // IP address of the node
	Port           int         // Port number the node is listening on
	NodeType       NodeType    // Type of the node (default, validator, miner)
	GeoLocation    GeoLocation // Geographical location of the node (latitude and longitude)
	LastActiveTime time.Time   // Last active timestamp of the node
	IsOnline       bool        // Status indicating if the node is currently online
}

// NodeHealthStatus represents the health status of a node in the network
type NodeHealthStatus string

const (
	Healthy   NodeHealthStatus = "healthy"   // Node is operating normally
	Unhealthy NodeHealthStatus = "unhealthy" // Node is experiencing issues or errors
	Degraded  NodeHealthStatus = "degraded"  // Node is operational but performing below expectations
	Offline   NodeHealthStatus = "offline"   // Node is not reachable or non-functional
)

// NodeState represents the operational state of a node in the network
type NodeState string

const (
	Active       NodeState = "active"        // Node is actively participating in the network
	Inactive     NodeState = "inactive"      // Node is currently inactive or idle
	Syncing      NodeState = "syncing"       // Node is synchronizing data with the network
	Validating   NodeState = "validating"    // Node is validating blocks or transactions
	Mining       NodeState = "mining"        // Node is mining (for PoW nodes)
	Staking      NodeState = "staking"       // Node is staking (for PoS nodes)
	ShuttingDown NodeState = "shutting_down" // Node is in the process of shutting down
)

// NodeMetrics represents performance and resource usage metrics for a node in the network.
type NodeMetrics struct {
	NodeID                 string        // Unique identifier for the node
	Uptime                 time.Duration // Uptime of the node
	TotalTransactions      int           // Total number of transactions processed by the node
	BlockProcessingTime    time.Duration // Average time taken to process a block
	SubBlockProcessingTime time.Duration // Average time taken to process a sub-block (if applicable)
	MemoryUsage            int64         // Current memory usage of the node (in bytes)
	CPUUsage               float64       // CPU usage percentage
	NetworkLatency         time.Duration // Average network latency for the node
	TotalBlocks            int           // Total number of blocks processed by the node
	ErrorsEncountered      int           // Total number of errors encountered by the node
	LastUpdated            time.Time     // The last time the metrics were updated
	PeerCount              int           // Number of peers connected to the node
	DiskUsage              float64       // Disk usage in MB
	NetworkTraffic         float64       // Network traffic in MB/s
	FaultCount             int           // Count of faults for the node
	LastChecked            time.Time     // Timestamp of the last health check
	IsHealthy              bool          // Indicates if the node is healthy
	Faulty                 bool          // Indicates if the node is faulty
	RecoveryState          bool          // Indicates if the node is in recovery state
	Latency                float64       // Network latency in milliseconds
}

// NodeData represents detailed information about a node in the network.
type NodeData struct {
	NodeID          string            // Unique identifier for the node
	NodeType        string            // Type of node (e.g., Validator, Full Node, Light Node)
	Owner           string            // Owner of the node (typically a wallet address or entity)
	Status          string            // Current status of the node (e.g., Active, Inactive, Pending)
	ConnectedPeers  []string          // List of node IDs for connected peers
	LastBlockHash   string            // Hash of the last block processed by this node
	LastBlockHeight int               // Height of the last block processed by the node
	DataStore       map[string]string // In-memory store of key-value data associated with the node
	SyncStatus      string            // Synchronization status of the node (e.g., Synced, Syncing, Not Synced)
	Version         string            // Software version the node is running
	StartupTime     time.Time         // Timestamp of when the node was started
}

// ************** Plasma Structs **************

// NetworkReconfig represents a network reconfiguration event.
type NetworkReconfig struct {
	Details       string    // Description of the reconfiguration event
	Timestamp     time.Time // Time when the reconfiguration occurred
	AffectedNodes []string  // List of nodes affected by the reconfiguration
}

// PlasmaClient represents a client in the Plasma childchain network
type PlasmaClient struct {
	ClientID          string          // Unique identifier for the client
	WalletAddress     string          // Wallet address associated with the client
	Ledger            *Ledger         // Reference to the ledger for transaction recording
	EncryptionService *Encryption     // Encryption service for securing client transactions
	PlasmaChain       *PlasmaChain    // Reference to the Plasma childchain
	NetworkManager    *NetworkManager // Network manager for interacting with nodes and the childchain
}

// PlasmaCore represents the core logic of the Plasma childchain
type PlasmaCore struct {
	Blocks            map[string]*PlasmaBlock    // Collection of all blocks in the childchain
	SubBlocks         map[string]*PlasmaSubBlock // Collection of sub-blocks within the childchain
	Ledger            *Ledger                    // Reference to the ledger for recording chain events
	EncryptionService *Encryption                // Encryption service to secure transactions
	mu                sync.Mutex                 // Mutex for handling concurrency
}

// PlasmaBlock represents a block in the Plasma childchain
type PlasmaBlock struct {
	BlockID       string                     // Unique identifier for the block
	PreviousBlock string                     // ID of the previous block in the chain
	SubBlocks     map[string]*PlasmaSubBlock // Collection of sub-blocks in the block
	Timestamp     time.Time                  // Timestamp of when the block was created
	ValidatorID   string                     // Validator responsible for the block
}

// PlasmaSubBlock represents a sub-block in the Plasma childchain
type PlasmaSubBlock struct {
	SubBlockID    string              // Unique identifier for the sub-block
	ParentBlockID string              // ID of the parent block
	Transactions  []TransactionRecord // List of transactions included in the sub-block
	Timestamp     time.Time           // Timestamp of when the sub-block was created
	ValidatorID   string              // Validator responsible for the sub-block
}

// PlasmaChain represents a Plasma childchain, containing blocks and managing the network
type PlasmaChain struct {
	ChainID        string                  // Unique identifier for the Plasma childchain
	GenesisBlock   *PlasmaBlock            // The genesis block of the chain
	CurrentBlockID string                  // ID of the current block
	Blocks         map[string]*PlasmaBlock // All blocks in the chain
	Ledger         *Ledger                 // Ledger for recording events
	Core           *PlasmaCore             // Core logic for sub-block and block validation
	Encryption     *Encryption             // Encryption service for securing the chain
	NetworkManager *NetworkManager         // Network manager for handling nodes and chain state
}

// PlasmaChainConfig represents the configuration for creating a Plasma childchain
type PlasmaChainConfig struct {
	ChainID           string          // Unique ID for the Plasma childchain
	GenesisBlockID    string          // ID of the genesis block
	GenesisTimestamp  time.Time       // Timestamp for the genesis block
	PlasmaCore        *PlasmaCore     // Reference to the Plasma core logic
	Ledger            *Ledger         // Reference to the ledger for recording events
	NetworkManager    *NetworkManager // Network manager for handling nodes
	EncryptionService *Encryption     // Encryption service for securing the chain
}

// PlasmaCrossChain represents the logic for handling cross-chain operations
type PlasmaCrossChain struct {
	Blocks         map[string]*PlasmaBlock    // Collection of all blocks in the childchain
	SubBlocks      map[string]*PlasmaSubBlock // Collection of sub-blocks within the childchain
	Ledger         *Ledger                    // Reference to the ledger for recording chain events
	Encryption     *Encryption                // Encryption service to secure transactions
	CrossChainComm *CrossChainCommunication   // Cross-chain communication logic
	mu             sync.Mutex                 // Mutex for concurrency handling
}

// PlasmaNetwork represents the Plasma childchain network operations
type PlasmaNetwork struct {
	PlasmaNodes    map[string]*PlasmaNode     // Collection of all participating Plasma nodes
	Blocks         map[string]*PlasmaBlock    // Collection of all blocks in the Plasma childchain
	SubBlocks      map[string]*PlasmaSubBlock // Collection of all sub-blocks in the Plasma childchain
	Ledger         *Ledger                    // Reference to the ledger for recording chain events
	Encryption     *Encryption                // Encryption service to secure transactions
	NetworkManager *NetworkManager            // Network manager to handle node communications
	mu             sync.Mutex                 // Mutex for concurrency handling
}

// PlasmaNode represents a node in the Plasma childchain network
type PlasmaNode struct {
	NodeID         string           // Unique identifier for the node
	IPAddress      string           // Node's IP address
	Encryption     *Encryption      // Encryption service used by the node
	LastActiveTime time.Time        // Timestamp of the node's last activity
	NodeHealth     NodeHealthStatus // Current health status of the node
	NodeState      NodeState        // Node's current operational state
}

// RootChainTxRecord represents the record of a transaction on the root chain.
type RootChainTxRecord struct {
	TxHash string
}

// ************** Resource Management Structs **************

// Resource represents a system resource managed by the blockchain network.
type Resource struct {
	ID              string    // Unique identifier for the resource
	Type            string    // Type of resource (e.g., "CPU", "Memory", "Storage")
	AvailableUnits  int       // Number of units available for allocation
	Usage           float64   // Current usage level
	Limit           float64   // Usage limit for the resource
	AllocatedTo     string    // ID of the account or node the resource is allocated to
	UsedAmount      float64   // The amount of resource that has been used
	OveruseReported bool      // Indicates if the resource has been overused
	CreatedAt       time.Time // Time when the resource was registered
}

// Lease represents the details of a leased resource in the marketplace.
type Lease struct {
	LeaseID       string        // Unique identifier for the lease
	ResourceID    string        // ID of the resource being leased
	Leaser        string        // Wallet address of the leaser
	LeaseDuration time.Duration // Duration of the lease
	StartTime     time.Time     // Time when the lease started
	EndTime       time.Time     // Time when the lease ends
	Price         float64       // Price paid for the lease
}

// QuotaManager manages resource quotas for nodes, users, and processes.
type QuotaManager struct {
	QuotaID             string              // Unique identifier for the quota.
	TargetEntity        string              // The entity the quota is applied to (e.g., user, node, process).
	ResourceType        string              // Type of resource (e.g., CPU, memory, storage).
	AllocatedQuota      float64             // The allocated quota for the resource.
	UsedQuota           float64             // The currently used amount of the resource.
	QuotaLimit          float64             // The maximum allowable quota.
	EnforcementPolicies []string            // Policies for enforcing quota limits.
	LastUpdated         time.Time           // Timestamp of the last quota update.
	Notifications       []QuotaNotification // Notifications triggered by quota events.
}

// QuotaNotification defines a notification related to quota usage.
type QuotaNotification struct {
	NotificationID string    // Unique identifier for the notification.
	QuotaID        string    // The associated quota ID.
	Message        string    // The notification message.
	TriggeredAt    time.Time // When the notification was triggered.
}

// ResourceManager oversees resource allocation, optimization, and scaling for the network.
type ResourceManager struct {
	ResourceID         string             // Unique identifier for the resource.
	ResourceType       string             // Type of resource (e.g., compute, memory, storage).
	TotalCapacity      float64            // Total capacity of the resource in the system.
	AllocatedCapacity  float64            // Total allocated capacity.
	AvailableCapacity  float64            // Remaining available capacity.
	AllocationPolicies []ResourcePolicy   // Policies governing resource allocation.
	OptimizationRules  []OptimizationRule // Rules for optimizing resource usage.
	LastOptimized      time.Time          // Timestamp of the last optimization operation.
}

// ResourcePolicy defines policies for allocating resources.
type ResourcePolicy struct {
	PolicyID        string                 // Unique identifier for the policy.
	Description     string                 // Description of the policy.
	AllocationRules map[string]interface{} // Rules for resource allocation (e.g., priority-based, weighted).
}

// OptimizationRule represents a rule for optimizing resource usage.
type OptimizationRule struct {
	RuleID      string    // Unique identifier for the rule.
	Description string    // Description of the optimization rule.
	LastApplied time.Time // Timestamp when the rule was last applied.
}

// Purchase represents the details of a purchased resource in the marketplace.
type Purchase struct {
	PurchaseID   string    // Unique identifier for the purchase
	ResourceID   string    // ID of the resource being purchased
	Buyer        string    // Wallet address of the buyer
	PurchaseTime time.Time // Time when the purchase was made
	Price        float64   // Price paid for the purchase
}

// ResourceRequest represents a request for resource allocation.
type ResourceRequest struct {
	RequestID    string    // Unique identifier for the request
	ResourceType string    // Type of resource being requested
	Units        int       // Number of units requested
	Requester    string    // Wallet address of the requester
	Timestamp    time.Time // Time when the request was made
}

// ResourceAllocationManager handles the allocation and management of system resources across the Synnergy Network.
type ResourceAllocationManager struct {
	AllocatedResources map[string]Resource // Resources allocated to nodes
	mutex              sync.Mutex          // Mutex for thread-safe operations
	LedgerInstance     *Ledger             // Ledger for resource tracking and auditing
}

// ResourceManager handles all resource-related operations across the Synnergy Network.
type ResourceManager struct {
	Resources       map[string]Resource // The resources available for management.
	mutex           sync.Mutex          // Mutex for thread-safe resource management.
	LedgerInstance  *Ledger             // The ledger instance to record resource operations.
	AllocationQueue []ResourceRequest   // Queue for pending resource allocation requests.
}

// ResourceMarketplace manages the leasing, purchasing, and selling of resources.
type ResourceMarketplace struct {
	AvailableResources map[string]Resource // List of available resources for leasing or purchasing
	LeasedResources    map[string]Lease    // Leased resources with details
	Purchases          map[string]Purchase // Completed resource purchases
	mutex              sync.Mutex          // Mutex for thread-safe operations
	LedgerInstance     *Ledger             // Ledger instance for tracking resource transactions
}

// ResourceAllocation represents resource allocation across the channels.
type ResourceAllocation struct {
	ChannelID     string    // Channel ID
	Resources     int       // Amount of resources allocated to the channel
	ReallocatedAt time.Time // Timestamp when resources were last reallocated
}

// ************** Rollups Structs **************

// MultiRollupOracle represents an oracle that can serve multiple rollups asynchronously
type MultiRollupOracle struct {
	OracleID       string                 // Unique identifier for the oracle
	Rollups        map[string]*Rollup     // Map of rollups the oracle is serving
	DataSources    map[string]interface{} // Data sources the oracle is fetching from
	Ledger         *Ledger                // Reference to the ledger for logging events
	Encryption     *Encryption            // Encryption for securing oracle data
	NetworkManager *NetworkManager        // Network manager for handling communications
	mu             sync.Mutex             // Mutex for concurrency handling
}

// DecentralizedGovernanceRollup represents a rollup with integrated decentralized governance
type DecentralizedGovernanceRollup struct {
	RollupID        string                         // Unique identifier for the rollup
	Transactions    []*Transaction                 // Transactions included in the rollup
	StateRoot       string                         // Root hash of the final state after the rollup
	IsFinalized     bool                           // Whether the rollup is finalized
	Ledger          *Ledger                        // Reference to the ledger for recording rollup events
	Encryption      *Encryption                    // Encryption service for securing governance data
	NetworkManager  *NetworkManager                // Network manager for communications
	VotingProposals map[string]*GovernanceProposal // Proposals in the governance system
	Participants    []string                       // List of participants in the governance system
	mu              sync.Mutex                     // Mutex for concurrency control
}

// HierarchicalProofRollup represents a rollup with hierarchical proof compression and verification
type HierarchicalProofRollup struct {
	RollupID       string                       // Unique identifier for the rollup
	ProofHierarchy map[string]*RollupProofLayer // Hierarchical layers of proofs
	Transactions   []*Transaction               // Transactions in the rollup
	StateRoot      string                       // Root hash of the final state
	IsFinalized    bool                         // Whether the rollup has been finalized
	Ledger         *Ledger                      // Reference to the ledger for recording events
	Encryption     *Encryption                  // Encryption service for securing data
	Consensus      *SynnergyConsensus           // Consensus system for verifying hierarchical proofs
	NetworkManager *NetworkManager              // Network manager for communications
	mu             sync.Mutex                   // Mutex for concurrency control
}

// RollupProofLayer represents a hierarchical layer in the proof structure
type RollupProofLayer struct {
	LayerID     string // Unique identifier for the proof layer
	Proof       string // Proof generated for this layer
	ParentLayer string // Parent layer in the hierarchy
	IsVerified  bool   // Whether the proof for this layer has been verified
}

// HyperLayeredRollupFramework represents a multi-layer rollup orchestration system
type HyperLayeredRollupFramework struct {
	FrameworkID    string                  // Unique identifier for the framework
	RollupLayers   map[string]*RollupLayer // Rollup layers in the framework
	Ledger         *Ledger                 // Reference to the ledger for recording events
	Encryption     *Encryption             // Encryption service for securing data
	Consensus      *SynnergyConsensus      // Consensus system for verifying proofs between layers
	NetworkManager *NetworkManager         // Network manager for handling communications between layers
	mu             sync.Mutex              // Mutex for concurrency control
}

// RollupLayer represents an individual rollup layer within the hyper-layered framework
type RollupLayer struct {
	LayerID     string             // Unique identifier for the rollup layer
	Rollups     map[string]*Rollup // Rollups within this layer
	StateRoot   string             // Root hash of the state of the layer
	IsFinalized bool               // Whether the layer has been finalized
}

// InteroperableRollupLayer represents a rollup layer that handles cross-rollup communication and shared state
type InteroperableRollupLayer struct {
	LayerID        string                 // Unique identifier for the interoperable rollup layer
	Rollups        map[string]*Rollup     // Rollups within the layer
	SharedState    map[string]interface{} // Shared state across rollups
	Ledger         *Ledger                // Reference to the ledger for recording rollup events
	Encryption     *Encryption            // Encryption service for securing cross-rollup data
	Consensus      *SynnergyConsensus     // Consensus system for validating state changes
	NetworkManager *NetworkManager        // Network manager for cross-rollup communication
	mu             sync.Mutex             // Mutex for concurrency control
}

// LiquidRollupPool manages dynamic liquidity pools with automated yield redistribution
type LiquidRollupPool struct {
	PoolID         string                        // Unique identifier for the liquidity pool
	Assets         map[string]map[string]float64 // Mapping of assets to participants' balances
	YieldRates     map[string]float64            // Yield rates for each asset
	Transactions   []*Transaction                // Transactions associated with the pool
	IsFinalized    bool                          // Whether the pool's current rollup cycle is finalized
	Ledger         *Ledger                       // Reference to the ledger for recording pool events
	Encryption     *Encryption                   // Encryption service for securing data
	Consensus      *SynnergyConsensus            // Consensus system for pool validation
	NetworkManager *NetworkManager               // Network manager for handling communication
	mu             sync.Mutex                    // Mutex for concurrency control
}

// MultiAssetLiquidityRollup (MALR) handles cross-asset liquidity bridging and reconciliation within a rollup.
type MultiAssetLiquidityRollup struct {
	RollupID       string                        // Unique identifier for the rollup
	LiquidityPools map[string]map[string]float64 // Mapping of assets to liquidity pools (e.g., asset -> [participant -> balance])
	Transactions   []*Transaction                // List of transactions in the rollup
	IsFinalized    bool                          // Whether the rollup is finalized
	Ledger         *Ledger                       // Ledger for recording events
	Encryption     *Encryption                   // Encryption service for securing data
	Consensus      *SynnergyConsensus            // Consensus system for validating rollups
	NetworkManager *NetworkManager               // Network manager for communications
	mu             sync.Mutex                    // Mutex for concurrency control
}

// SubmissionRecord tracks the submission of batches, rollups, or other data.
type SubmissionRecord struct {
	SubmissionID string    // Unique identifier for the submission
	RollupID     string    // Associated rollup ID
	SubmittedAt  time.Time // Timestamp of submission
	SubmittedBy  string    // ID of the submitting entity (e.g., node, user)
	Status       string    // Status of the submission (e.g., "Pending", "Finalized")
	Details      string    // Additional details or metadata
}

// VerificationRecord tracks the verification of transactions, rollups, or proofs.
type VerificationRecord struct {
	VerificationID string    // Unique identifier for the verification
	RollupID       string    // Associated rollup ID
	VerifiedAt     time.Time // Timestamp of verification
	VerifiedBy     string    // ID of the verifying entity (e.g., node, validator)
	Result         string    // Result of the verification (e.g., "Success", "Failure")
	Details        string    // Additional verification metadata or logs
}

// OracleValidationLog tracks validation results of Oracle data.
type OracleValidationLog struct {
	OracleID       string    // Oracle identifier
	RollupID       string    // Associated rollup ID
	ValidationTime time.Time // Timestamp of validation
	Result         string    // Validation result (e.g., "Valid", "Invalid")
	Details        string    // Additional information or metadata
}

// GovernanceUpdate represents updates to governance rules or policies.
type GovernanceUpdate struct {
	UpdateID  string    // Unique identifier for the governance update
	RollupID  string    // Associated rollup ID
	UpdatedAt time.Time // Timestamp of the update
	UpdatedBy string    // ID of the updating entity (e.g., DAO, admin)
	Details   string    // Description of the update
}

// GovernanceApplication represents the application of a governance decision.
type GovernanceApplication struct {
	ApplicationID string    // Unique identifier for the application
	RollupID      string    // Associated rollup ID
	AppliedAt     time.Time // Timestamp of the application
	AppliedBy     string    // ID of the applying entity (e.g., DAO, admin)
	Details       string    // Description of the application
}

// GovernanceMonitor tracks ongoing governance monitoring actions.
type GovernanceMonitor struct {
	MonitorID   string    // Unique identifier for the monitoring action
	RollupID    string    // Associated rollup ID
	MonitoredAt time.Time // Timestamp of the monitoring
	Status      string    // Current status (e.g., "Active", "Inactive")
	Details     string    // Additional information or metadata
}

// ZKProof represents a zero-knowledge proof.
type ZKProof struct {
	ProofID     string    // Unique identifier for the proof
	RollupID    string    // Associated rollup ID
	ProofData   string    // Encoded proof data
	GeneratedAt time.Time // Timestamp of generation
	GeneratedBy string    // ID of the generating entity
	Description string    // Description or purpose of the proof
}

// ZKProofVerification represents the verification of a zero-knowledge proof.
type ZKProofVerification struct {
	VerificationID string    // Unique identifier for the verification
	ProofID        string    // Associated proof ID
	VerifiedAt     time.Time // Timestamp of verification
	VerifiedBy     string    // ID of the verifying entity
	Result         string    // Verification result (e.g., "Valid", "Invalid")
	Details        string    // Additional verification metadata
}

// Proof represents a generic proof in the rollup ledger.
type Proof struct {
	ProofID     string    // Unique identifier for the proof
	RollupID    string    // Associated rollup ID
	ProofType   string    // Type of proof (e.g., "ZK", "Fraud")
	Data        string    // Encoded proof data
	GeneratedAt time.Time // Timestamp of proof generation
	GeneratedBy string    // ID of the generating entity
}

// ProofAggregation represents the aggregation of multiple proofs.
type ProofAggregation struct {
	AggregationID string    // Unique identifier for the aggregation
	RollupID      string    // Associated rollup ID
	AggregatedAt  time.Time // Timestamp of aggregation
	AggregatedBy  string    // ID of the aggregating entity
	Proofs        []string  // List of proof IDs included in the aggregation
	Result        string    // Aggregation result (e.g., "Successful", "Failed")
	Details       string    // Additional aggregation metadata
}

// CrossRollupTransaction tracks transactions between rollups.
type CrossRollupTransaction struct {
	TransactionID string    // Unique identifier for the transaction
	SourceRollup  string    // ID of the source rollup
	TargetRollup  string    // ID of the target rollup
	Amount        float64   // Amount transferred
	Fee           float64   // Transaction fee
	Timestamp     time.Time // Timestamp of the transaction
	Status        string    // Status of the transaction (e.g., "Pending", "Completed")
	Details       string    // Additional metadata
}

// LayerVerification represents verification details for a rollup layer.
type LayerVerification struct {
	LayerID    string    // Unique identifier for the layer
	RollupID   string    // Associated rollup ID
	VerifiedAt time.Time // Timestamp of verification
	VerifiedBy string    // ID of the verifying entity
	Result     string    // Verification result (e.g., "Valid", "Invalid")
	Details    string    // Additional verification information
}

// LayerFinalization tracks the finalization of a layer within a rollup.
type LayerFinalization struct {
	LayerID     string    // Unique identifier for the layer
	RollupID    string    // Associated rollup ID
	FinalizedAt time.Time // Timestamp of finalization
	FinalizedBy string    // ID of the finalizing entity
	Details     string    // Metadata about the finalization process
}

// ResultSync represents the synchronization of results across layers or rollups.
type ResultSync struct {
	SyncID   string    // Unique identifier for the synchronization
	RollupID string    // Associated rollup ID
	SyncedAt time.Time // Timestamp of synchronization
	SyncedBy string    // ID of the entity performing synchronization
	Status   string    // Sync status (e.g., "Success", "Failed")
	Details  string    // Additional information about the sync
}

// VerificationResult represents the result of a validation or proof verification.
type VerificationResult struct {
	VerificationID string    // Unique identifier for the verification
	RollupID       string    // Associated rollup ID
	Result         string    // Result (e.g., "Valid", "Invalid")
	VerifiedAt     time.Time // Timestamp of verification
	VerifiedBy     string    // ID of the verifying entity
	Details        string    // Additional metadata
}

// SharedStateValidation tracks validation of shared state between rollups.
type SharedStateValidation struct {
	ValidationID  string    // Unique identifier for the validation
	SharedStateID string    // ID of the shared state
	ValidatedAt   time.Time // Timestamp of validation
	ValidatorID   string    // ID of the validating entity
	Result        string    // Validation result
	Details       string    // Metadata or notes about the validation
}

// BridgeTransaction represents a transaction between two rollups via a bridge.
type BridgeTransaction struct {
	TransactionID string    // Unique identifier for the transaction
	SourceRollup  string    // ID of the source rollup
	TargetRollup  string    // ID of the target rollup
	Amount        float64   // Amount transferred
	Timestamp     time.Time // Timestamp of the transaction
	Status        string    // Transaction status
	Details       string    // Additional metadata
}

// BridgeFinalization logs the finalization of a bridge-related operation.
type BridgeFinalization struct {
	BridgeID    string    // Unique identifier for the bridge
	FinalizedAt time.Time // Timestamp of finalization
	FinalizedBy string    // ID of the finalizing entity
	Details     string    // Metadata about the finalization process
}

// ChallengeRecord tracks challenges and their resolutions.
type ChallengeRecord struct {
	ChallengeID string    // Unique identifier for the challenge
	RollupID    string    // Associated rollup ID
	SubmittedBy string    // ID of the entity submitting the challenge
	SubmittedAt time.Time // Timestamp of submission
	Status      string    // Status of the challenge (e.g., "Pending", "Resolved")
	Resolution  string    // Description of the resolution
	Details     string    // Metadata about the challenge
}

// NodeConnectionRecord tracks connections or disconnections of nodes.
type NodeConnectionRecord struct {
	NodeID    string    // Unique identifier for the node
	RollupID  string    // Associated rollup ID
	Action    string    // Action (e.g., "Connected", "Disconnected")
	Timestamp time.Time // Timestamp of the action
	Details   string    // Additional metadata
}

// SyncRecord tracks synchronization events between nodes or rollups.
type SyncRecord struct {
	SyncID   string    // Unique identifier for the sync
	Source   string    // Source node or rollup
	Target   string    // Target node or rollup
	SyncedAt time.Time // Timestamp of synchronization
	Status   string    // Sync status (e.g., "Completed", "Failed")
	Details  string    // Additional metadata
}

// FeeRecord tracks fee-related operations in the rollup ledger.
type FeeRecord struct {
	FeeID     string    // Unique identifier for the fee record
	RollupID  string    // Associated rollup ID
	Amount    float64   // Fee amount
	Timestamp time.Time // Timestamp of the fee event
	Action    string    // Action (e.g., "Applied", "Refunded")
	Details   string    // Additional metadata
}

// BaseFeeRecord tracks updates to base fees for rollups.
type BaseFeeRecord struct {
	BaseFeeID  string    // Unique identifier for the base fee record
	RollupID   string    // Associated rollup ID
	NewBaseFee float64   // Updated base fee value
	UpdatedAt  time.Time // Timestamp of the update
	UpdatedBy  string    // ID of the updating entity
	Details    string    // Additional metadata
}

// ScalingRecord logs scaling operations in the rollup ledger.
type ScalingRecord struct {
	ScalingID   string    // Unique identifier for the scaling record
	RollupID    string    // Associated rollup ID
	Action      string    // Scaling action (e.g., "Up", "Down")
	ScaleFactor float64   // Factor by which scaling occurred
	Timestamp   time.Time // Timestamp of the scaling event
	Details     string    // Additional metadata
}

// YieldDistributionRecord tracks yield distribution in rollup pools.
type YieldDistributionRecord struct {
	DistributionID string    // Unique identifier for the yield distribution
	RollupID       string    // Associated rollup ID
	Amount         float64   // Yield amount distributed
	DistributedAt  time.Time // Timestamp of distribution
	DistributedBy  string    // ID of the distributing entity
	Details        string    // Metadata about the distribution
}

// Batch represents a batch in a rollup.
type Batch struct {
	BatchID      string    // Unique identifier for the batch
	RollupID     string    // Associated rollup ID
	Transactions []string  // List of transaction IDs in the batch
	CreatedAt    time.Time // Timestamp of batch creation
	Status       string    // Batch status (e.g., "Pending", "Submitted")
	Details      string    // Metadata about the batch
}

// BatchValidation logs the validation of a batch.
type BatchValidation struct {
	ValidationID string    // Unique identifier for the validation
	BatchID      string    // Associated batch ID
	ValidatedAt  time.Time // Timestamp of validation
	ValidatorID  string    // ID of the validator
	Result       string    // Validation result (e.g., "Valid", "Invalid")
	Details      string    // Additional metadata
}

// BatchBroadcast logs the broadcasting of a batch.
type BatchBroadcast struct {
	BroadcastID string    // Unique identifier for the broadcast
	BatchID     string    // Associated batch ID
	RollupID    string    // Associated rollup ID
	BroadcastAt time.Time // Timestamp of broadcasting
	Status      string    // Status of the broadcast
	Details     string    // Metadata about the broadcast
}

// ResourceManagementLog represents a log entry for resource management operations.
type ResourceManagementLog struct {
	LogID      string            // Unique identifier for the log entry
	ResourceID string            // ID of the resource involved
	Action     string            // Description of the action performed (e.g., "Allocated", "Released")
	Timestamp  time.Time         // Time when the action was performed
	Details    map[string]string // Additional details about the resource operation
}

// ResourceIssue represents an issue related to a resource's usage or management.
type ResourceIssue struct {
	IssueID       string    // Unique identifier for the issue
	ResourceID    string    // ID of the resource associated with the issue
	Description   string    // Detailed description of the issue
	ReportedAt    time.Time // Timestamp when the issue was reported
	SeverityLevel string    // Severity of the issue (e.g., "Low", "Medium", "High", "Critical")
	Resolved      bool      // Indicates if the issue has been resolved
	Resolution    string    // Description of how the issue was resolved (if applicable)
	ResolvedAt    time.Time // Timestamp when the issue was resolved (if applicable)
}

// LeaseRecord represents the details of a leased resource.
type LeaseRecord struct {
	LeaseID        string            // Unique identifier for the lease
	ResourceID     string            // ID of the resource being leased
	LesseeID       string            // ID of the entity leasing the resource
	LeaseAmount    float64           // Amount charged for the lease
	LeaseDuration  time.Duration     // Duration of the lease
	LeaseStartTime time.Time         // Start time of the lease
	LeaseEndTime   time.Time         // End time of the lease
	Status         string            // Status of the lease (e.g., "Active", "Expired", "Cancelled")
	AdditionalInfo map[string]string // Additional metadata or details about the lease
}

// BatchSubmission logs the submission of a batch.
type BatchSubmission struct {
	SubmissionID string    // Unique identifier for the submission
	BatchID      string    // Associated batch ID
	RollupID     string    // Associated rollup ID
	SubmittedAt  time.Time // Timestamp of submission
	SubmittedBy  string    // ID of the submitting entity
	Status       string    // Submission status
	Details      string    // Metadata about the submission
}

// ContractRecord represents a smart contract associated with a rollup.
type ContractRecord struct {
	ContractID string    // Unique identifier for the contract
	RollupID   string    // Associated rollup ID
	DeployedAt time.Time // Timestamp of contract deployment
	DeployedBy string    // ID of the deploying entity
	Status     string    // Status of the contract (e.g., "Active", "Terminated")
	Details    string    // Metadata about the contract
}

// CancelledTransaction tracks a cancelled transaction in the rollup ledger.
type CancelledTransaction struct {
	TransactionID string    // Unique identifier of the transaction
	RollupID      string    // Associated rollup ID
	CancelledAt   time.Time // Timestamp of cancellation
	CancelledBy   string    // ID of the cancelling entity
	Reason        string    // Reason for cancellation
}

// PruningRecord tracks the pruning of data in the rollup ledger.
type PruningRecord struct {
	PruningID string    // Unique identifier for the pruning action
	RollupID  string    // Associated rollup ID
	PrunedAt  time.Time // Timestamp of pruning
	PrunedBy  string    // ID of the entity performing pruning
	DataSize  int64     // Size of the data pruned (e.g., in bytes)
	Details   string    // Additional pruning metadata
}

// MultiDimensionalCompressionRollup (MDCR) represents a rollup system that applies multi-dimensional compression to rollup data
type MultiDimensionalCompressionRollup struct {
	RollupID        string                // Unique identifier for the rollup
	Transactions    []*Transaction        // Transactions within the rollup
	CompressedData  []byte                // Compressed representation of the transactions
	CompressionAlgo *CompressionAlgorithm // Compression algorithm used for multi-dimensional compression
	IsCompressed    bool                  // Flag to indicate whether the rollup data is compressed
	Ledger          *Ledger               // Ledger for recording compression events
	Consensus       *SynnergyConsensus    // Consensus system for validating compression
	Encryption      *Encryption           // Encryption service for securing compressed data
	NetworkManager  *NetworkManager       // Network manager for communication
	mu              sync.Mutex            // Mutex for concurrency control
}

// MultiLayerTransactionPruning (MLTP) represents a multi-layer system for pruning old rollup data.
type MultiLayerTransactionPruning struct {
	RollupID        string                 // Unique identifier for the rollup
	Transactions    []*Transaction         // Transactions within the rollup
	PrunedLayers    map[int][]*Transaction // Map of pruned layers (by pruning level) and their transactions
	PruningInterval time.Duration          // Time interval to trigger pruning
	MaxRetention    time.Duration          // Maximum age of transactions before pruning
	Ledger          *Ledger                // Ledger for recording pruning events
	Consensus       *SynnergyConsensus     // Consensus system for validating pruning
	Encryption      *Encryption            // Encryption service for securing transaction data
	NetworkManager  *NetworkManager        // Network manager for broadcast
	mu              sync.Mutex             // Mutex for concurrency control
}

// ParallelExecutionRollupLayer (PERL) handles parallel processing and transaction bundling within rollups.
type ParallelExecutionRollupLayer struct {
	LayerID        string             // Unique identifier for the parallel execution layer
	RollupID       string             // Rollup ID this layer belongs to
	Transactions   []*Transaction     // Transactions to be processed in parallel
	IsFinalized    bool               // Whether the layer's execution is finalized
	Encryption     *Encryption        // Encryption service for securing transactions
	Ledger         *Ledger            // Ledger for recording all actions
	Consensus      *SynnergyConsensus // Consensus system for validating transactions
	NetworkManager *NetworkManager    // Network manager for broadcasting
	mu             sync.Mutex         // Mutex for handling concurrency
}

// PoSTRollup represents a Proof-of-Space-Time (PoSTR) enabled rollup in the network.
type PoSTRollup struct {
	RollupID          string             // Unique identifier for the rollup
	Transactions      []*Transaction     // Transactions included in the rollup
	StateRoot         string             // Root hash of the final state after the rollup
	SpaceTimeProof    *SpaceTimeProof    // Proof of space-time for data storage
	IsFinalized       bool               // Whether the rollup is finalized
	Ledger            *Ledger            // Reference to the ledger for recording rollup events
	Encryption        *Encryption        // Encryption service for securing data
	NetworkManager    *NetworkManager    // Network manager for communication
	SynnergyConsensus *SynnergyConsensus // Consensus for validation
	mu                sync.Mutex         // Mutex for concurrency control
}

// RecursiveProofAggregation represents the RPA system for zk-SNARK proof aggregation in rollups
type RecursiveProofAggregation struct {
	AggregationID   string             // Unique identifier for the recursive proof aggregation
	Rollups         map[string]*Rollup // Collection of rollups being aggregated
	AggregatedProof *ZkProof           // Aggregated zk-SNARK proof
	IsFinalized     bool               // Whether the aggregation process is finalized
	Ledger          *Ledger            // Reference to the ledger for recording proof events
	Encryption      *Encryption        // Encryption service for securing proof data
	Consensus       *SynnergyConsensus // Consensus for validating proofs
	NetworkManager  *NetworkManager    // Network manager for proof broadcasting
	mu              sync.Mutex         // Mutex for concurrency control
}

// Rollup represents a rollup in the network, which aggregates multiple transactions into a single batch.
type Rollup struct {
	RollupID         string                          // Unique identifier for the rollup
	Transactions     []*Transaction                  // Transactions included in the rollup
	StateRoot        string                          // Root hash of the final state after the rollup
	IsFinalized      bool                            // Whether the rollup is finalized
	Ledger           *Ledger                         // Reference to the ledger for recording rollup events
	Encryption       *Encryption                     // Encryption service for securing data
	NetworkManager   *NetworkManager                 // Network manager for communications
	CancellationMgr  *TransactionCancellationManager // Transaction cancellation manager
	mu               sync.Mutex                      // Mutex for concurrency control
	TotalFees        float64                         // Total fees aggregated within the rollup
	ValidatorAddress string                          // Validator responsible for processing the rollup
	CreationTime     time.Time                       // Timestamp for when the rollup was created
	OracleSources    []string                        // List of oracle sources associated with the rollup (new field)
	DataSources      []string                        // List of data sources associated with the rollup (new field)

}

// RollupNetwork represents a network handling the communication and coordination of rollups
type RollupNetwork struct {
	Nodes          map[string]*Node   // Collection of all participating nodes
	Rollups        map[string]*Rollup // Collection of rollups managed by the network
	Ledger         *Ledger            // Reference to the ledger for recording rollup events
	Encryption     *Encryption        // Encryption service for securing data
	NetworkManager *NetworkManager    // Network manager for node communication
	mu             sync.Mutex         // Mutex for handling concurrency
}

// AssetType defines the type of asset being bridged (either Token or Main SYNN Coin)
type AssetType string

const (
	TokenType AssetType = "TOKEN" // Bridging specific tokens
	SYNNType  AssetType = "SYNN"  // Bridging the main SYNN coin
)

// RollupBridge represents a bridge that facilitates asset transfer and communication between rollups and other chains
type RollupBridge struct {
	BridgeID       string               // Unique identifier for the bridge
	RollupID       string               // Rollup being connected to
	TargetChainID  string               // ID of the target blockchain/rollup to bridge with
	Transactions   map[string]*BridgeTx // Transactions flowing through the bridge
	Ledger         *Ledger              // Reference to the ledger for recording bridge events
	Encryption     *Encryption          // Encryption service for securing transactions
	NetworkManager *NetworkManager      // Network manager to handle cross-chain communications
	Consensus      *SynnergyConsensus   // Consensus for bridge validation
	mu             sync.Mutex           // Mutex for handling concurrency
}

// BridgeTx represents a transaction that flows through the bridge
type BridgeTx struct {
	TxID        string    // Unique transaction ID
	SourceChain string    // Source chain of the transaction
	Destination string    // Destination chain or rollup
	Amount      float64   // Amount being transferred
	AssetType   AssetType // Type of asset (TOKEN or SYNN coin)
	Timestamp   time.Time // Transaction creation time
	IsFinalized bool      // Whether the transaction is finalized
}

// RollupChallenge represents a mechanism to challenge invalid rollup states or transactions
type RollupChallenge struct {
	ChallengeID     string             // Unique identifier for the challenge
	RollupID        string             // Rollup being challenged
	Challenger      string             // Address of the challenger
	ChallengedBlock string             // Block ID or state that is being challenged
	Timestamp       time.Time          // Time when the challenge was created
	IsResolved      bool               // Whether the challenge has been resolved
	Ledger          *Ledger            // Reference to the ledger for recording challenge events
	Encryption      *Encryption        // Encryption service for securing challenge data
	Consensus       *SynnergyConsensus // Consensus mechanism to resolve the challenge
	mu              sync.Mutex         // Mutex for handling concurrency
}

// RollupContract represents a smart contract deployed in the rollup
type RollupContract struct {
	ContractID    string                 // Unique identifier for the contract
	ContractOwner string                 // Owner of the contract
	ContractState map[string]interface{} // Current state of the contract
	Transactions  []*Transaction         // Transactions interacting with the contract
	IsDeployed    bool                   // Whether the contract is deployed on the rollup
	Ledger        *Ledger                // Reference to the ledger for recording contract events
	Encryption    *Encryption            // Encryption service for securing contract data
	Consensus     *SynnergyConsensus     // Consensus mechanism for contract validation
	mu            sync.Mutex             // Mutex for handling concurrency
}

// RollupFeeManager manages fees for transactions within a rollup
type RollupFeeManager struct {
	FeeID           string             // Unique identifier for the fee instance
	TransactionFees map[string]float64 // Mapping of transaction IDs to their associated fees
	BaseFee         float64            // Base fee for transactions within the rollup
	TotalFees       float64            // Accumulated fees within the rollup
	Ledger          *Ledger            // Reference to the ledger for recording fee events
	Encryption      *Encryption        // Encryption service for securing fee data
	Consensus       *SynnergyConsensus // Consensus mechanism for fee validation
	mu              sync.Mutex         // Mutex for handling concurrency
}

// RollupNode represents a node within the rollup network
type RollupNode struct {
	NodeID         string                 // Unique identifier for the rollup node
	IPAddress      string                 // IP address of the node
	NodeType       NodeType               // Type of node (e.g., aggregator, validator)
	ConnectedNodes map[string]*RollupNode // Connected nodes in the rollup network
	Ledger         *Ledger                // Ledger for recording transactions
	Encryption     *Encryption            // Encryption service for secure communication
	Consensus      *SynnergyConsensus     // Consensus used for transaction validation
	mu             sync.Mutex             // Mutex for concurrency control
}

// RollupBatch represents a batch of aggregated transactions in the rollup network
type RollupBatch struct {
	BatchID      string         // Unique identifier for the transaction batch
	Transactions []*Transaction // Transactions aggregated in the rollup
	MerkleRoot   string         // Merkle root for the batch of transactions
	Timestamp    time.Time      // Timestamp when the batch was created
}

// RollupOperator represents an operator responsible for managing rollup processes
type RollupOperator struct {
	OperatorID     string                  // Unique identifier for the rollup operator
	NodeID         string                  // Associated rollup node ID
	IPAddress      string                  // IP address of the operator
	ManagedBatches map[string]*RollupBatch // Collection of batches managed by this operator
	Ledger         *Ledger                 // Ledger for recording actions
	Encryption     *Encryption             // Encryption service for securing data
	Consensus      *SynnergyConsensus      // Consensus engine for batch validation
	NetworkManager *NetworkManager         // Network manager for node communication
	mu             sync.Mutex              // Mutex for concurrency control
}

// OptimisticRollup represents an optimistic rollup process
type OptimisticRollup struct {
	RollupID        string                          // Unique identifier for the rollup
	NodeID          string                          // Rollup node handling this process
	Transactions    []*Transaction                  // Transactions being processed
	Ledger          *Ledger                         // Reference to the ledger for recording events
	Encryption      *Encryption                     // Encryption service to secure transaction data
	NetworkManager  *NetworkManager                 // Network manager for communicating with other nodes
	Consensus       *SynnergyConsensus              // Consensus for transaction validation (optimistic)
	SubmittedProofs map[string]*FraudProof          // Submitted fraud proofs for disputed transactions
	CancellationMgr *TransactionCancellationManager // Transaction cancellation manager for processing frauds
	mu              sync.Mutex                      // Mutex for concurrency control
}

// FraudProof represents a proof that challenges the validity of a transaction in an optimistic rollup
type FraudProof struct {
	ProofID    string    // Unique identifier for the fraud proof
	TxID       string    // The transaction being challenged
	Challenger string    // The node that submitted the fraud proof
	Evidence   string    // Evidence provided to dispute the transaction's validity
	Timestamp  time.Time // Timestamp of when the fraud proof was submitted
	IsResolved bool      // Indicates whether the fraud proof has been resolved
}

// RollupScalingManager handles dynamic scaling of rollups based on network load, transaction volume, and other factors
type RollupScalingManager struct {
	RollupID        string                          // Unique identifier for the rollup
	ScalingFactor   float64                         // Current scaling factor for the rollup
	MaxScalingLimit float64                         // Maximum allowable scaling factor
	MinScalingLimit float64                         // Minimum allowable scaling factor
	Transactions    []*Transaction                  // Transactions handled by the rollup
	Ledger          *Ledger                         // Reference to the ledger for recording events
	Encryption      *Encryption                     // Encryption service to secure transaction data
	NetworkManager  *NetworkManager                 // Network manager for communicating with nodes
	CancellationMgr *TransactionCancellationManager // Transaction cancellation manager
	mu              sync.Mutex                      // Mutex for concurrency control
}

// RollupVerifier handles the verification of rollups, ensuring validity of transactions and state roots.
type RollupVerifier struct {
	VerifierID        string             // Unique identifier for the verifier
	Ledger            *Ledger            // Reference to the ledger for recording verification events
	Encryption        *Encryption        // Encryption service for securing data
	NetworkManager    *NetworkManager    // Network manager for communication
	SynnergyConsensus *SynnergyConsensus // Synnergy Consensus for verification
	mu                sync.Mutex         // Mutex for concurrency control
}

// ZKRollup represents a zero-knowledge proof (ZKP) enabled rollup in the network.
type ZKRollup struct {
	RollupID          string             // Unique identifier for the rollup
	Transactions      []*Transaction     // Transactions included in the rollup
	StateRoot         string             // Root hash of the final state after the rollup
	ZKProof           *ZkProof           // Zero-knowledge proof for validating the rollup
	IsFinalized       bool               // Whether the rollup is finalized
	Ledger            *Ledger            // Reference to the ledger for recording rollup events
	Encryption        *Encryption        // Encryption service for securing data
	NetworkManager    *NetworkManager    // Network manager for communication
	SynnergyConsensus *SynnergyConsensus // Consensus for validation
	mu                sync.Mutex         // Mutex for concurrency control
}

// SelfGoverningRollupEcosystem represents a rollup system with self-governing capabilities.
type SelfGoverningRollupEcosystem struct {
	EcosystemID     string             // Unique identifier for the ecosystem
	Rollups         map[string]*Rollup // Collection of rollups in the ecosystem
	GovernanceRules *GovernanceRules   // Automated governance rules for the ecosystem
	Ledger          *Ledger            // Reference to the ledger for recording events
	Encryption      *Encryption        // Encryption service for securing governance data
	NetworkManager  *NetworkManager    // Network manager for off-chain communication
	Consensus       *SynnergyConsensus // Consensus mechanism for governance validation
	mu              sync.Mutex         // Mutex for concurrency control
}

// GovernanceRules represents the self-regulating parameters for the rollup ecosystem.
type GovernanceRules struct {
	MaxTransactionsPerRollup int       // Maximum number of transactions per rollup
	FeeStructure             float64   // Fee structure for processing transactions
	ScalingFactor            float64   // Scaling factor to adjust rollup capacity
	LastUpdated              time.Time // Timestamp when the governance rules were last updated
}

// SmartContractCoProcessingLayer represents the co-processing layer for smart contracts, enabling off-chain execution for speed.
type SmartContractCoProcessingLayer struct {
	LayerID        string                 // Unique identifier for the processing layer
	RollupID       string                 // Associated Rollup ID
	Contracts      []*SmartContract       // Smart contracts processed in the layer
	Results        map[string]interface{} // Results of the contract execution
	IsFinalized    bool                   // Whether the processing is finalized
	Ledger         *Ledger                // Reference to the ledger for recording events
	Encryption     *Encryption            // Encryption service for securing contract data
	NetworkManager *NetworkManager        // Network manager for off-chain communication
	Consensus      *SynnergyConsensus     // Consensus mechanism for verification
	mu             sync.Mutex             // Mutex for concurrency control
}

// TemporalRollup represents a time-based rollup in the network.
type TemporalRollup struct {
	RollupID          string             // Unique identifier for the rollup
	Transactions      []*Transaction     // Transactions in the rollup
	StateRoot         string             // Root hash of the final state after the rollup
	CreationTime      time.Time          // Timestamp when the rollup was created
	PruneThreshold    time.Duration      // Duration after which older data is pruned
	IsFinalized       bool               // Whether the rollup is finalized
	Ledger            *Ledger            // Reference to the ledger for recording rollup events
	Encryption        *Encryption        // Encryption service for securing data
	NetworkManager    *NetworkManager    // Network manager for communications
	SynnergyConsensus *SynnergyConsensus // Consensus for rollup validation
	mu                sync.Mutex         // Mutex for concurrency control
}

// ZeroLatencyRollupBridge represents a zero-latency bridge between rollups, allowing instant synchronization.
type ZeroLatencyRollupBridge struct {
	BridgeID            string             // Unique identifier for the bridge
	SourceRollupID      string             // ID of the source rollup
	DestinationRollupID string             // ID of the destination rollup
	Transactions        []*Transaction     // Transactions transferred across the bridge
	StateRoot           string             // Final state root after synchronization
	IsFinalized         bool               // Whether the bridge sync is finalized
	Ledger              *Ledger            // Reference to the ledger for recording bridge events
	Encryption          *Encryption        // Encryption service for securing bridge data
	NetworkManager      *NetworkManager    // Network manager for bridge communication
	SynnergyConsensus   *SynnergyConsensus // Consensus mechanism for bridge verification
	mu                  sync.Mutex         // Mutex for concurrency control
}

// ************** Scalability Structs **************

// CompressionSystem manages various compression methods for blocks, transactions, and files
type CompressionSystem struct {
	Ledger            *Ledger     // Ledger for logging compression actions
	EncryptionService *Encryption // Encryption for securing data before compression
	mu                sync.Mutex  // Mutex for concurrency control
}

type ShardRecord struct {
	ID           string              // Unique identifier for the shard
	StateChannel string              // Associated state channel, if applicable
	Status       string              // Current status of the shard (e.g., "active", "closed")
	Transactions []TransactionRecord // List of transactions processed within this shard
	Resources    int                 // Resource allocation for the shard
	IsAvailable  bool                // Availability status of the shard
	ValidatedBy  string              // ID of the validator that last validated this shard
	ValidatedAt  time.Time           // Timestamp of the last validation
	CreatedAt    time.Time           // Timestamp when the shard was created
	ClosedAt     time.Time           // Timestamp when the shard was closed
	SyncedAt     time.Time           // Timestamp of the last synchronization
	Notes        string              // Optional notes or metadata about the shard
}

// HandoverManager manages the transfer of tasks, responsibilities, or data between nodes or systems.
type HandoverManager struct {
	HandoverID          string           // Unique identifier for the handover operation.
	SourceNode          string           // The node initiating the handover.
	TargetNode          string           // The node receiving the handover.
	HandoverType        string           // Type of handover (e.g., task, data, responsibility).
	Status              string           // Current status of the handover (e.g., pending, in-progress, completed).
	AssociatedResources []string         // Resources involved in the handover process.
	Policies            []HandoverPolicy // Policies governing the handover process.
	LastUpdated         time.Time        // Timestamp of the last status update.
	Logs                []HandoverLog    // Logs related to the handover operation.
}

// HandoverPolicy defines policies governing handover operations.
type HandoverPolicy struct {
	PolicyID    string                 // Unique identifier for the policy.
	Description string                 // Description of the policy.
	Conditions  map[string]interface{} // Conditions for triggering or completing a handover.
}

// HandoverLog captures events or updates during the handover process.
type HandoverLog struct {
	LogID      string    // Unique identifier for the log entry.
	HandoverID string    // Associated handover operation ID.
	Message    string    // Log message.
	Timestamp  time.Time // Timestamp of the log entry.
}

// LoadBalancer manages the distribution of workloads across nodes or systems.
type LoadBalancer struct {
	BalancerID       string              // Unique identifier for the load balancer.
	DistributionType string              // Type of distribution algorithm (e.g., round-robin, weighted).
	Nodes            []LoadBalancerNode  // List of nodes under the balancer's control.
	Metrics          LoadBalancerMetrics // Performance metrics for the load balancer.
	Policies         []BalancingPolicy   // Policies governing load balancing.
	LastUpdated      time.Time           // Timestamp of the last balancing operation.
}

// LoadBalancerNode represents a node managed by the load balancer.
type LoadBalancerNode struct {
	NodeID      string  // Unique identifier for the node.
	CurrentLoad float64 // Current workload on the node.
	MaxCapacity float64 // Maximum capacity of the node.
	Status      string  // Status of the node (e.g., active, overloaded, offline).
}

// LoadBalancerMetrics contains performance data for the load balancer.
type LoadBalancerMetrics struct {
	TotalRequests    int       // Total number of requests handled.
	AverageLatency   float64   // Average response time in milliseconds.
	SuccessfulRoutes int       // Number of successful routing operations.
	FailedRoutes     int       // Number of failed routing operations.
	LastUpdated      time.Time // Timestamp of the last metrics update.
}

// BalancingPolicy defines policies for load distribution.
type BalancingPolicy struct {
	PolicyID    string                 // Unique identifier for the policy.
	Description string                 // Description of the policy.
	Rules       map[string]interface{} // Rules governing workload distribution.
}

type PartitionRecord struct {
	ID               string              // Unique identifier for the partition
	Status           string              // Current status of the partition (e.g., "active", "rebalance")
	Configuration    map[string]string   // Configuration details for the partition (e.g., size, resource allocation)
	CreatedAt        time.Time           // Timestamp when the partition was created
	LastUpdated      time.Time           // Timestamp of the last update to the partition
	RebalanceLog     []TransactionRecord // Log of rebalance actions taken on the partition
	AssociatedShards []string            // List of shard IDs associated with this partition
	Notes            string              // Optional notes or metadata about the partition
}

// CacheEntry represents a cache entry in the system
type CacheEntry struct {
	Key        string    // Key for the cache entry
	Data       []byte    // Cached data
	Timestamp  time.Time // Timestamp of when the data was cached
	Expiration time.Time // Expiration time for the cached entry
}

// DataRetrievalSystem manages cache retrieval and prefetching for efficient data access
type DataRetrievalSystem struct {
	Cache             map[string]*CacheEntry // In-memory cache
	CacheTTL          time.Duration          // Time-to-live for cached data
	PrefetchKeys      []string               // Keys to prefetch
	Ledger            *Ledger                // Ledger for logging cache activities
	EncryptionService *Encryption            // Encryption for securing data
	mu                sync.Mutex             // Mutex for concurrency control
}

// DelegatedAdaptiveLayerSharding manages dynamic shard delegation and adaptation for rollups and state channels
type DelegatedAdaptiveLayerSharding struct {
	Shards            map[string]*Shard   // Collection of shards
	NodeShards        map[string][]string // Map of nodes to the shards they manage
	Ledger            *Ledger             // Ledger for logging shard actions
	EncryptionService *Encryption         // Encryption for securing state data
	mu                sync.Mutex          // Mutex for concurrency control
}

// DistributionNode represents a node that can receive transactions or tasks for distribution purposes
type DistributionNode struct {
	NodeID   string    // Unique identifier for the node
	NodeType string    // Type of node (e.g., validator, execution, storage, etc.)
	Load     int       // Current load of the node
	Weight   int       // Weight assigned to the node (used in weighted distribution)
	LastTask time.Time // Timestamp when the node last received a task
}

// DistributionSystem manages the distribution of transactions or tasks using different strategies
type DistributionSystem struct {
	Nodes             []*DistributionNode // List of available distribution nodes
	Ledger            *Ledger             // Ledger for logging distribution activities
	EncryptionService *Encryption         // Encryption for securing data transfers
	mu                sync.Mutex          // Mutex for concurrency control
}

// GossipNode represents a node in the gossip network with various types (validator, execution, etc.)
type GossipNode struct {
	NodeID       string        // Unique identifier for the node
	NodeType     string        // Type of node (e.g., validator, execution, storage)
	LastSyncTime time.Time     // Last time the node was synchronized
	Neighbors    []*GossipNode // Neighboring nodes that receive gossip updates
}

// GossipMessage represents a message being gossiped in the network
type GossipMessage struct {
	MessageID  string    // Unique identifier for the message
	Data       []byte    // The message data (encrypted)
	Timestamp  time.Time // The time the message was created
	OriginNode string    // The node that originally created the message
}

// GossipSystem implements the gossip protocol, redundancy protocol, and sync protocol
type GossipSystem struct {
	Nodes             []*GossipNode // List of all nodes participating in gossip
	Ledger            *Ledger       // Ledger for logging gossip activities
	EncryptionService *Encryption   // Encryption for securing messages
	mu                sync.Mutex    // Mutex for concurrency control
}

// LiquidShard represents a dynamic shard in the Liquid-State Sharding system
type LiquidShard struct {
	ShardID          string    // Unique identifier for the shard
	ShardType        string    // Type of shard (e.g., cross-chain, rollup, etc.)
	AllocatedChains  []string  // List of chains this shard is allocated to
	ShardState       []byte    // The state data of the shard (encrypted)
	LastAdjustment   time.Time // Timestamp of the last shard adjustment
	ReallocationTime time.Time // Time when the shard was last reallocated
}

// LiquidStateSharding manages the dynamic allocation and adjustment of shards across chains
type LiquidStateSharding struct {
	LiquidShards      map[string]*LiquidShard // Collection of shards
	Ledger            *Ledger                 // Ledger for logging shard activities
	EncryptionService *Encryption             // Encryption for secure shard state handling
	mu                sync.Mutex              // Mutex for concurrency control
}

// MetaLayerShard represents a shard participating in the meta-layer orchestration
type MetaLayerShard struct {
	ShardID      string    // Unique identifier for the shard
	ShardType    string    // Type of shard (e.g., execution, storage, validation)
	LastActivity time.Time // Last time the shard participated in an orchestration
}

// MetaLayerOrchestrator manages cross-layer orchestration between shards
type MetaLayerOrchestrator struct {
	Shards            map[string]*MetaLayerShard // Collection of all active shards
	Ledger            *Ledger                    // Ledger to log orchestration activities
	EncryptionService *Encryption                // Encryption service for securing orchestrated data
	mu                sync.Mutex                 // Mutex for concurrency control
}

// Partition represents a partition for data and load management
type Partition struct {
	PartitionID    string    // Unique identifier for the partition
	PartitionType  string    // Type of partition (horizontal or vertical)
	Data           []byte    // Encrypted data in this partition
	LastRebalanced time.Time // Timestamp of the last rebalancing
	LastAdjusted   time.Time // Timestamp of the last dynamic adjustment
	ToleranceLimit int       // Tolerance limit for partition overload
}

// PartitionManager handles the management of partitions in the network
type PartitionManager struct {
	Partitions        map[string]*Partition // Collection of partitions
	Ledger            *Ledger               // Ledger for logging partition activities
	EncryptionService *Encryption           // Encryption service for partition data
	mu                sync.Mutex            // Mutex for concurrency control
}

// Shard represents a shard within the system, handling part of the blockchain state.
type Shard struct {
	ShardID         string              // Unique identifier for the shard
	LayerID         string              // Associated layer ID (rollup or state channel)
	StateChannel    string              // State channel this shard is currently serving
	RollupID        string              // Rollup ID this shard is associated with
	AssignedNodes   []string            // Nodes assigned to this shard
	StateData       []byte              // State data managed by the shard
	LastUpdate      time.Time           // Last time the shard was updated
	ShardSize       int                 // The size of the shard (measured in data blocks)
	LastReallocated time.Time           // Last time the shard was reallocated
	ParentShardID   string              // ID of the parent shard for hierarchical sharding (if applicable)
	Data            []byte              // Encrypted data within the shard
	LastMergedSplit time.Time           // Timestamp for last merge or split
	IsAvailable     bool                // Availability of shard for new transactions
	CreatedAt       time.Time           // Timestamp of shard creation
	ClosedAt        time.Time           // Timestamp of shard closure
	Status          string              // Shard status (active, closed, etc.)
	SyncedAt        time.Time           // Timestamp of the last sync operation
	Resources       int                 // Resource allocation for the shard
	ValidatedBy     string              // Validator that validated the shard
	ValidatedAt     time.Time           // Timestamp of the shard validation
	Transactions    []TransactionRecord // List of transactions within the shard
}

// ShardReallocationManager manages real-time shard reallocation
type ShardReallocationManager struct {
	Shards            map[string]*Shard // All shards available for reallocation
	Ledger            *Ledger           // Ledger for logging shard activities
	EncryptionService *Encryption       // Encryption service for securing shard data
	mu                sync.Mutex        // Mutex for concurrency control
}

// ShardManager manages all shard operations, including cross-shard communication, hierarchical, horizontal, and vertical sharding
type ShardManager struct {
	Shards            map[string]*Shard // All shards in the network
	Ledger            *Ledger           // Ledger for logging shard operations
	EncryptionService *Encryption       // Encryption service for securing shard data
	mu                sync.Mutex        // Mutex for concurrency control
}

// CompressionAlgorithm defines various compression algorithms available for the system
type CompressionAlgorithm struct {
	AlgorithmName string                       // Name of the compression algorithm
	Compress      func([]byte) ([]byte, error) // Compression function
	Decompress    func([]byte) ([]byte, error) // Decompression function
}

// ************** Security Structs **************

// ThreatDetection handles real-time threat detection and anomaly analysis in the blockchain.
type ThreatDetection struct {
	DetectionID       string            // Unique identifier for each detection instance
	ActiveMonitors    map[string]bool   // Monitors activated for specific nodes, channels, or transactions
	DetectedAnomalies []AnomalyReport   // A list of detected anomalies (transactions, behaviors, etc.)
	AlertThreshold    int               // Threshold for raising alerts (number of anomalies before action)
	ActionTaken       map[string]Action // Actions taken in response to detected threats
	LoggingEnabled    bool              // Flag indicating if logging of threats is enabled
	mutex             sync.Mutex        // Mutex to ensure thread-safe operations
}

// AnomalyReport represents a detected anomaly or suspicious activity.
type AnomalyReport struct {
	ReportID       string    // Unique identifier for the anomaly report
	NodeID         string    // Node that reported the anomaly
	TransactionID  string    // Related transaction ID if applicable
	Description    string    // Detailed description of the anomaly
	Severity       string    // Severity level (Low, Medium, High, Critical)
	Timestamp      time.Time // Timestamp of when the anomaly was detected
	ActionRequired bool      // Flag to indicate if further action is needed
}

// Action represents the action taken in response to a detected threat.
type Action struct {
	ActionID    string    // Unique identifier for the action
	ActionType  string    // Type of action (e.g., "Block", "Alert", "Investigate")
	Description string    // Detailed description of the action taken
	Timestamp   time.Time // Timestamp of when the action was executed
}

// SecurityManager handles all security-related operations for the blockchain network
type SecurityManager struct {
	LedgerInstance      *Ledger             // Ledger for recording security events
	ActiveSessions      map[string]*Session // Active sessions mapped by session ID
	mutex               sync.Mutex          // Mutex for thread-safe operations
	WhitelistedIPs      map[string]bool     // Whitelisted IP addresses
	BlacklistedIPs      map[string]bool     // Blacklisted IP addresses
	NodeKeys            map[string]string   // Holds public keys for authorized nodes
	MasterSecurityKey   string              // Master key for system-level security management
	EncryptionService   *Encryption         // Encryption service for securing sensitive data
	FailedLoginAttempts map[string]int      // Tracks failed login attempts by IP or UserID
	AlertThreshold      int                 // Threshold for raising alerts (e.g., failed login attempts)
}

// Session represents an active session within the blockchain network
type Session struct {
	SessionID     string    // Unique ID for the session
	UserID        string    // User ID associated with the session
	StartTime     time.Time // Session start time
	LastActivity  time.Time // Last activity timestamp
	IPAddress     string    // IP address associated with the session
	EncryptionKey string    // Encryption key for securing session data
}

// SecurityEvent represents a security-related event logged within the blockchain network
type SecurityEvent struct {
	EventID   string    // Unique ID for the event
	EventType string    // Type of event (e.g., "LOGIN", "ACCESS_DENIED", "FAILED_LOGIN")
	Timestamp time.Time // Time when the event occurred
	UserID    string    // User ID associated with the event
	IPAddress string    // IP address involved in the event
	Details   string    // Additional details regarding the event
}

// LoginAttempt represents an attempt to log into the network
type LoginAttempt struct {
	AttemptID  string    // Unique ID for the login attempt
	UserID     string    // User ID making the login attempt
	Timestamp  time.Time // Time of the login attempt
	IPAddress  string    // IP address from which the attempt was made
	Successful bool      // Whether the login attempt was successful
}

// IntrusionDetectionSystem handles the detection and prevention of suspicious activities in the network
type IntrusionDetectionSystem struct {
	SuspiciousIPs     map[string]int  // Tracks suspicious IP addresses and their activity count
	BlacklistedIPs    map[string]bool // Blacklisted IP addresses
	Threshold         int             // Threshold for flagging an IP as suspicious
	LedgerInstance    *Ledger         // Ledger for logging detected intrusions
	EncryptionService *Encryption     // Encryption service for securing logs
	mutex             sync.Mutex      // Mutex for thread-safe access
}

// SessionLog represents a log entry for a session.
type SessionLog struct {
	SessionID string
	Message   string
	Timestamp time.Time
}

// SecurityAlert represents a security alert raised within the network
type SecurityAlert struct {
	AlertID     string    // Unique ID for the alert
	AlertType   string    // Type of alert (e.g., "INTRUSION", "FAILED_LOGIN")
	Timestamp   time.Time // Time the alert was raised
	Severity    string    // Severity level of the alert (e.g., "LOW", "MEDIUM", "HIGH")
	Description string    // Description of the alert
	Resolved    bool      // Whether the alert has been resolved
}

// EncryptionManager handles all encryption and decryption operations in the system
type EncryptionManager struct {
	EncryptionService *Encryption // Encryption service for managing cryptographic operations
	LedgerInstance    *Ledger     // Ledger for logging encryption-related events
}

// NodeSecurityManager manages node-specific security settings and access control
type NodeSecurityManager struct {
	NodeID            string          // Unique identifier for the node
	NodePublicKey     string          // Public key of the node for secure communication
	Whitelist         map[string]bool // Whitelisted IP addresses for the node
	Blacklist         map[string]bool // Blacklisted IP addresses for the node
	EncryptionService *Encryption     // Encryption service for securing node data
	LedgerInstance    *Ledger         // Ledger for recording node security events
}

// ************** Sidechain Structs **************

// SidechainCoin represents a native coin for a specific sidechain
type SidechainCoin struct {
	CoinID      string  // Unique identifier for the coin
	Name        string  // Name of the sidechain coin (e.g., SidechainToken)
	Symbol      string  // Symbol for the coin (e.g., SCT)
	TotalSupply float64 // Total supply of the coin
	Decimals    int     // Number of decimal places for the coin
}

// Coin represents the details of a coin in the sidechain.
type Coin struct {
	CoinID            string            // Unique identifier for the coin
	Name              string            // Name of the coin
	Symbol            string            // Symbol for the coin
	TotalSupply       float64           // Total supply of the coin
	CirculatingSupply float64           // Circulating supply of the coin
	Decimals          int               // Number of decimal places
	Creator           string            // Address of the coin creator
	CreationDate      time.Time         // Timestamp of coin creation
	Status            string            // Status of the coin (e.g., active, retired)
	Metadata          map[string]string // Additional metadata for the coin
}

// SidechainLogEntry represents a log entry for a sidechain-related event.
type SidechainLogEntry struct {
	EventID     string            // Unique identifier for the event
	SidechainID string            // Associated sidechain ID
	EventType   string            // Type of event (e.g., creation, upgrade)
	Description string            // Description of the event
	Initiator   string            // Entity that initiated the event
	Timestamp   time.Time         // Timestamp of the event
	Metadata    map[string]string // Additional metadata for the event
}

// StateUpdateLog represents a log entry for state updates.
type StateUpdateLog struct {
	UpdateID   string            // Unique identifier for the update
	StateID    string            // State being updated
	UpdateType string            // Type of update (e.g., addition, modification)
	Details    string            // Detailed description of the update
	Initiator  string            // Entity that initiated the update
	Timestamp  time.Time         // Timestamp of the update
	Metadata   map[string]string // Additional metadata for the update
}

// BlockStateLog represents a log entry for block state changes.
type BlockStateLog struct {
	BlockID       string            // Unique identifier for the block
	PreviousState string            // Previous state of the block
	NewState      string            // New state of the block
	ChangeType    string            // Type of change (e.g., finalized, reverted)
	Initiator     string            // Entity that initiated the state change
	Timestamp     time.Time         // Timestamp of the state change
	Metadata      map[string]string // Additional metadata for the block state change
}

// StateSyncLog represents a log entry for state synchronization events.
type StateSyncLog struct {
	SyncID    string            // Unique identifier for the synchronization event
	SourceID  string            // Source state or entity
	TargetID  string            // Target state or entity
	Status    string            // Status of the synchronization (e.g., completed, failed)
	Timestamp time.Time         // Timestamp of the synchronization
	Metadata  map[string]string // Additional metadata for the synchronization event
}

// ProposedBlock represents a block proposed for inclusion in the ledger.
type ProposedBlock struct {
	BlockID      string            // Unique identifier for the block
	Proposer     string            // Entity proposing the block
	Transactions []Transaction     // Transactions included in the proposed block
	Timestamp    time.Time         // Timestamp of block proposal
	Status       string            // Status of the proposed block (e.g., pending, accepted)
	Metadata     map[string]string // Additional metadata for the block proposal
}

// SidechainRegistry represents a registry entry for a sidechain.
type SidechainRegistry struct {
	SidechainID  string            // Unique identifier for the sidechain
	Name         string            // Name of the sidechain
	Description  string            // Description of the sidechain
	Creator      string            // Entity that created the sidechain
	CreationDate time.Time         // Timestamp of sidechain creation
	Status       string            // Status of the sidechain (e.g., active, deprecated)
	Metadata     map[string]string // Additional metadata for the sidechain
}

// BlockLogEntry represents a log entry for block-related actions.
type BlockLogEntry struct {
	BlockID    string            // Unique identifier for the block
	ActionType string            // Type of action (e.g., creation, validation)
	Initiator  string            // Entity that performed the action
	Timestamp  time.Time         // Timestamp of the action
	Metadata   map[string]string // Additional metadata for the block action
}

// SubBlockLogEntry represents a log entry for sub-block-related actions.
type SubBlockLogEntry struct {
	SubBlockID string            // Unique identifier for the sub-block
	ActionType string            // Type of action (e.g., validation, broadcast)
	Initiator  string            // Entity that performed the action
	Timestamp  time.Time         // Timestamp of the action
	Metadata   map[string]string // Additional metadata for the sub-block action
}

// NodeLogEntry represents a log entry for node-related events.
type NodeLogEntry struct {
	NodeID     string            // Unique identifier for the node
	ActionType string            // Type of action (e.g., addition, removal)
	Initiator  string            // Entity that performed the action
	Timestamp  time.Time         // Timestamp of the action
	Metadata   map[string]string // Additional metadata for the node action
}

// BlockValidationLog represents a log entry for block validation events.
type BlockValidationLog struct {
	BlockID   string            // Unique identifier for the block
	Validator string            // Entity that validated the block
	Status    string            // Validation status (e.g., successful, failed)
	Timestamp time.Time         // Timestamp of the validation
	Metadata  map[string]string // Additional metadata for the validation
}

// SecurityLogEntry represents a log entry for security-related events.
type SecurityLogEntry struct {
	EventID     string            // Unique identifier for the security event
	Description string            // Description of the security event
	Severity    string            // Severity level of the event
	Initiator   string            // Entity that triggered the event
	Timestamp   time.Time         // Timestamp of the event
	Metadata    map[string]string // Additional metadata for the event
}

// AssetTransferLog represents a log entry for asset transfer events.
type AssetTransferLog struct {
	TransferID string            // Unique identifier for the asset transfer
	AssetID    string            // Identifier for the transferred asset
	Sender     string            // Sender of the asset
	Receiver   string            // Receiver of the asset
	Amount     float64           // Amount of asset transferred
	Timestamp  time.Time         // Timestamp of the transfer
	Metadata   map[string]string // Additional metadata for the transfer
}

// TransactionValidationLog represents a log entry for transaction validation.
type TransactionValidationLog struct {
	TransactionID string            // Unique identifier for the transaction
	Validator     string            // Entity that validated the transaction
	Status        string            // Validation status (e.g., successful, failed)
	Timestamp     time.Time         // Timestamp of the validation
	Metadata      map[string]string // Additional metadata for the validation
}

// BlockSyncLog represents a log entry for block synchronization events.
type BlockSyncLog struct {
	BlockID   string            // Unique identifier for the block
	Source    string            // Source of the block
	Target    string            // Target of the synchronization
	Status    string            // Synchronization status
	Timestamp time.Time         // Timestamp of the synchronization
	Metadata  map[string]string // Additional metadata for the event
}

// SubBlockStateLog represents a log entry for sub-block state updates.
type SubBlockStateLog struct {
	SubBlockID    string            // Unique identifier for the sub-block
	PreviousState string            // Previous state of the sub-block
	NewState      string            // New state of the sub-block
	Initiator     string            // Entity that updated the state
	Timestamp     time.Time         // Timestamp of the update
	Metadata      map[string]string // Additional metadata for the update
}

// StateValidationLog represents a log entry for state validation events.
type StateValidationLog struct {
	StateID   string            // Unique identifier for the state
	Validator string            // Entity that validated the state
	Status    string            // Validation status (e.g., valid, invalid)
	Timestamp time.Time         // Timestamp of the validation
	Metadata  map[string]string // Additional metadata for the validation
}

// TransactionLogEntry represents a log entry for transaction-related events.
type TransactionLogEntry struct {
	TransactionID string            // Unique identifier for the transaction
	ActionType    string            // Type of action (e.g., creation, validation)
	Initiator     string            // Entity that performed the action
	Timestamp     time.Time         // Timestamp of the action
	Metadata      map[string]string // Additional metadata for the transaction
}

// UpgradeLog represents a log entry for upgrade events.
type UpgradeLog struct {
	UpgradeID   string            // Unique identifier for the upgrade
	Description string            // Description of the upgrade
	Status      string            // Current status of the upgrade (e.g., applied, pending)
	Initiator   string            // Entity that initiated the upgrade
	Timestamp   time.Time         // Timestamp of the upgrade
	Metadata    map[string]string // Additional metadata for the upgrade
}

// SidechainCoinSetup manages the operations and state of a sidechain coin
type SidechainCoinSetup struct {
	Coins        map[string]*SidechainCoin     // All coins in the sidechain
	Balances     map[string]map[string]float64 // User balances (map of coinID -> userID -> balance)
	Transactions []*Transaction                // List of all transactions involving the sidechain coin
	Ledger       *SidechainLedger              // Reference to the ledger for recording coin events
	Encryption   *Encryption                   // Encryption service for secure transactions
	Consensus    *SidechainConsensus           // Consensus mechanism for validating transactions
	mu           sync.Mutex                    // Mutex for concurrency control
}

// SidechainConsensus represents the consensus mechanism for the sidechain
type SidechainConsensus struct {
	Nodes           map[string]*SidechainNode // Nodes participating in the sidechain consensus
	PendingBlocks   map[string]*Block         // Blocks pending validation by consensus
	Ledger          *SidechainLedger          // Reference to the ledger for consensus events
	Encryption      *Encryption               // Encryption service to secure consensus-related communications
	ConsensusEngine *SynnergyConsensus        // Synnergy Consensus engine for validation
	mu              sync.Mutex                // Mutex for concurrency handling
}

// SidechainNode represents a node participating in the sidechain consensus
type SidechainNode struct {
	NodeID    string   // Unique identifier for the node
	IPAddress string   // IP address of the node
	NodeType  NodeType // Type of node (validator, full node, etc.)
}

// Sidechain represents a sidechain in the network
type Sidechain struct {
	ChainID       string                // Unique identifier for the sidechain
	ParentChainID string                // Identifier of the parent chain (usually the main chain)
	Blocks        map[string]*SideBlock // All blocks within the sidechain
	SubBlocks     map[string]*SubBlock  // All sub-blocks within the sidechain
	CoinSetup     *SidechainCoinSetup   // Coin setup for the sidechain
	Consensus     *SidechainConsensus   // Reference to Synnergy Consensus for validation
	Ledger        *SidechainLedger      // Reference to the ledger for sidechain recording
	Encryption    *Encryption           // Encryption service for securing data
	mu            sync.Mutex            // Mutex for handling concurrent operations
	Status        string                // Status of the sidechain (e.g., active, terminated)
	CreatedAt     time.Time             // Time when the sidechain was created

}

// SideBlock represents a block in the sidechain
type SideBlock struct {
	BlockID     string      // Unique identifier for the block
	SubBlocks   []*SubBlock // Sub-blocks included in this block
	ParentBlock string      // Parent block ID for chain continuity
	Timestamp   time.Time   // Block creation timestamp
	MerkleRoot  string      // Merkle root of the block
}

// SidechainNetwork represents a separate network for a sidechain
type SidechainNetwork struct {
	Nodes      map[string]*SidechainNode // Sidechain-specific nodes
	Blocks     map[string]*SideBlock     // Blocks in the sidechain
	SubBlocks  map[string]*SubBlock      // Sub-blocks in the sidechain
	Ledger     *SidechainLedger          // Ledger reference for recording sidechain events
	Encryption *Encryption               // Encryption service for sidechain security
	Consensus  *SidechainConsensus       // Consensus mechanism for the sidechain
	mu         sync.Mutex                // Mutex for concurrency handling
}

// SidechainDeployment manages sidechain deployment and its network
type SidechainDeployment struct {
	DeployedNetworks map[string]*SidechainNetwork // Deployed sidechain networks
	Ledger           *Ledger                      // Ledger for sidechain recording
	Encryption       *Encryption                  // Encryption service for securing sidechain data
	Consensus        *SynnergyConsensus           // Synnergy Consensus for validation
	mu               sync.Mutex                   // Mutex for concurrency handling
}

// SidechainInteroperability handles cross-chain and sidechain-to-mainchain interactions
type SidechainInteroperability struct {
	MainChain         *Blockchain                  // Reference to the main chain
	SidechainNetworks map[string]*SidechainNetwork // Deployed sidechains
	Ledger            *SidechainLedger             // Reference to the ledger for recording cross-chain events
	Encryption        *Encryption                  // Encryption service for securing data transfer
	Consensus         *SidechainConsensus          // Synnergy Consensus for validation across chains
	mu                sync.Mutex                   // Mutex for concurrency handling
}

// SidechainManager handles the overall management of sidechains, including creation, deployment, and monitoring
type SidechainManager struct {
	Sidechains       map[string]*Sidechain // List of active sidechains
	Ledger           *SidechainLedger      // Reference to the ledger for logging events
	Encryption       *Encryption           // Encryption service for securing sidechain operations
	SidechainNetwork *SidechainNetwork     // Network management for sidechain node communication
	Consensus        *SidechainConsensus   // Consensus system for validating transactions and sub-blocks
	mu               sync.Mutex            // Mutex for concurrency management
}

// SidechainNodeManager manages all SidechainNodes within the sidechain network
type SidechainNodeManager struct {
	Nodes      map[string]*SidechainNode // Collection of nodes participating in the sidechain
	mu         sync.Mutex                // Mutex for handling concurrency
	Consensus  *SidechainConsensus       // Reference to the Synnergy Consensus mechanism
	Ledger     *SidechainLedger          // Reference to the ledger for recording node actions
	Encryption *Encryption               // Encryption service for securing node interactions
}

// SidechainSecurityManager manages security protocols for the sidechain
type SidechainSecurityManager struct {
	Consensus        *SidechainConsensus // Consensus mechanism for transaction validation
	Encryption       *Encryption         // Encryption service for securing data
	Ledger           *SidechainLedger    // Reference to the ledger for logging security events
	SidechainNetwork *SidechainNetwork   // Network to handle secure communication between nodes
	mu               sync.Mutex          // Mutex for concurrent security operations
}

// SidechainState represents the state of a sidechain in the system
type SidechainState struct {
	ChainID        string                    // Unique identifier for the sidechain
	StateData      map[string]*StateObject   // Current state data for the sidechain
	BlockStates    map[string]*BlockState    // States for blocks within the sidechain
	SubBlockStates map[string]*SubBlockState // States for sub-blocks within the sidechain
	Ledger         *SidechainLedger          // Ledger for recording state changes
	Encryption     *Encryption               // Encryption service to secure state data
	Consensus      *SidechainConsensus       // Reference to the consensus mechanism
	mu             sync.Mutex                // Mutex for concurrency handling
}

// BlockState represents the state data for a specific block in the sidechain
type BlockState struct {
	BlockID   string                  // Unique identifier for the block
	StateData map[string]*StateObject // State data for the block
	Timestamp time.Time               // Timestamp of the last state update
}

// SubBlockState represents the state data for a specific sub-block in the sidechain
type SubBlockState struct {
	SubBlockID string                  // Unique identifier for the sub-block
	StateData  map[string]*StateObject // State data for the sub-block
	Timestamp  time.Time               // Timestamp of the last state update
}

// SidechainTransaction represents a transaction within the sidechain
type SidechainTransaction struct {
	TxID        string    // Unique identifier for the transaction
	From        string    // Sender's wallet address
	To          string    // Recipient's wallet address
	Amount      float64   // Amount being transferred
	Fee         float64   // Transaction fee
	Timestamp   time.Time // Time the transaction was created
	IsValidated bool      // Whether the transaction has been validated
	IsFinalized bool      // Whether the transaction has been finalized (completed)
}

// SidechainTransactionPool manages pending transactions
type SidechainTransactionPool struct {
	PendingTransactions map[string]*SidechainTransaction // Pool of transactions awaiting validation
	Ledger              *SidechainLedger                 // Ledger reference for recording transactions
	Consensus           *SidechainConsensus              // Consensus mechanism for transaction validation
	Encryption          *Encryption                      // Encryption service for securing transactions
	mu                  sync.Mutex                       // Mutex for concurrency handling
}

// SidechainUpgrade represents an upgrade for a sidechain (e.g., protocol upgrade, consensus upgrade)
type SidechainUpgrade struct {
	UpgradeID      string              // Unique identifier for the upgrade
	Description    string              // Description of the upgrade
	UpgradeTime    time.Time           // Time the upgrade was initiated
	IsApplied      bool                // Whether the upgrade has been applied
	Consensus      *SidechainConsensus // Consensus method upgrade
	Encryption     *Encryption         // Encryption upgrade (if necessary)
	NetworkChanges bool                // Whether network topology is altered during the upgrade
}

// SidechainUpgradeManager manages sidechain upgrades
type SidechainUpgradeManager struct {
	PendingUpgrades   map[string]*SidechainUpgrade // Collection of pending upgrades
	CompletedUpgrades map[string]*SidechainUpgrade // Collection of applied upgrades
	Ledger            *SidechainLedger             // Ledger reference for recording upgrade events
	NetworkManager    *NetworkManager              // Network manager for node communication
	mu                sync.Mutex                   // Mutex for handling concurrency
}

// StateObject represents the state of an individual object in the blockchain network
type StateObject struct {
	ObjectID     string                 // Unique identifier for the state object
	OwnerID      string                 // ID of the owner of the object
	StateData    map[string]interface{} // Key-value pairs representing the state data
	LastModified time.Time              // Timestamp of the last modification to the state
	IsActive     bool                   // Indicates if the object is currently active
	Permissions  map[string]bool        // Permissions associated with the object (e.g., read, write access)
}

// ************** Smart Contract Structs **************

// SmartContract represents a smart contract in the Synnergy blockchain.
type SmartContract struct {
	ID             string                 // Unique ID of the contract
	Code           string                 // Smart contract code (e.g., bytecode or source)
	Parameters     map[string]interface{} // Default parameters (initial values or configuration)
	Bytecode       string                 // Compiled bytecode of the smart contrac
	State          map[string]interface{} // Current state (storage variables, balances, etc.)
	Owner          string                 // Owner of the smart contract (e.g., the creator’s wallet address)
	Executions     []ContractExecution    // History of contract executions
	CreationTime   time.Time              // Time the contract was created
	LastModified   time.Time              // Last time the contract was modified or executed
	IsActive       bool                   // Indicates if the contract is active or not
	mutex          sync.Mutex             // Mutex for safe concurrency
	LedgerInstance *Ledger                // Ledger instance for storing contract data and actions
}

// ContractInteraction represents an interaction between smart contracts on the blockchain.
type ContractInteraction struct {
	InteractionID            string                 // Unique identifier for the interaction
	InitiatingContractID     string                 // ID of the contract initiating the interaction
	ReceivingContractID      string                 // ID of the contract receiving the interaction
	Caller                   string                 // Address of the caller who initiated the interaction
	FunctionName             string                 // Name of the function called in the contract
	Parameters               map[string]interface{} // Parameters passed to the function
	Timestamp                time.Time              // Timestamp of when the interaction took place
	ExecutionResult          string                 // Result of the execution (e.g., success, failure)
	InteractionData          string                 // Raw data of the interaction
	EncryptedInteractionData string                 // Encrypted data of the interaction
}

// ContractInteractionManager manages the interactions between different smart contracts on the blockchain.
type ContractInteractionManager struct {
	Interactions   map[string]*ContractInteraction // A map of contract interactions by interaction ID
	LedgerInstance *Ledger                         // Ledger instance for storing interaction records
	mutex          sync.Mutex                      // Mutex for thread-safe operations
}

// ContractMigrationRecord tracks details of smart contract migrations.
type ContractMigrationRecord struct {
	MigrationID            string    // Unique identifier for the migration process
	ContractID             string    // Unique identifier of the smart contract being migrated
	SourceChain            string    // Source blockchain or network
	DestinationChain       string    // Destination blockchain or network
	MigrationInitiator     string    // Address or entity that initiated the migration
	Timestamp              time.Time // Time when the migration was initiated
	Status                 string    // Status of the migration (e.g., "Pending", "Completed", "Failed")
	ErrorDetails           string    // Details of any errors encountered during migration
	OriginalContractState  string    // Snapshot of the contract's state before migration
	NewContractState       string    // Snapshot of the contract's state after migration
	ExecutionLogs          []string  // Logs of actions taken during the migration
	AssociatedTransactions []string  // List of transaction IDs related to the migration
	VerificationProof      string    // Proof of successful migration (e.g., cryptographic proof)
	MigrationFee           float64   // Fee incurred during the migration process
}

// StoredContract represents the data for a stored contract in the blockchain.
type StoredContract struct {
	ContractID     string                 // Unique identifier for the contract
	Code           string                 // Smart contract code (e.g., bytecode or source)
	Owner          string                 // Address of the contract owner
	State          map[string]interface{} // Current state of the contract
	CreationTime   time.Time              // Time when the contract was deployed
	EncryptedCode  []byte                 // Encrypted smart contract code
	EncryptedState []byte                 // Encrypted state or parameters of the contract
}

// ContractStorageManager manages the storage of contract data on the blockchain.
type ContractStorageManager struct {
	Contracts      map[string]*StoredContract // A map to store contract data by contract ID
	LedgerInstance *Ledger                    // Ledger instance for storing contract data
	mutex          sync.Mutex                 // Mutex for thread-safe operations
}

// ContractDeployment represents the deployment record of a smart contract
type ContractDeployment struct {
	ContractID   string    // Unique ID of the deployed contract
	Deployer     string    // Address of the user deploying the contract
	ContractCode string    // Contract code
	DeployedAt   time.Time // Timestamp of deployment
	StoredHash   string    // Hash of the deployed contract code
	Status       string    // Current status of the contract (e.g., "open", "closed")
	ClosedBy     string    // Address of the user who closed the contract
	ClosedAt     time.Time // Timestamp when the contract was closed
}

// ContractMigration represents a migration of a smart contract to a new version
type ContractMigration struct {
	OldContractAddress string    // Address of the old contract
	NewContractAddress string    // Address of the new contract
	MigratorAddress    string    // Address of the user migrating the contract
	Timestamp          time.Time // Timestamp of migration
}

// ContractSignature represents a record of contract signatures
type ContractSignature struct {
	ContractAddress string    // Address of the contract
	SignerAddress   string    // Address of the signer
	Signature       string    // Signature provided by the signer
	Timestamp       time.Time // Timestamp of the signature
}

// CrossChainSmartContract represents a cross-chain smart contract that interacts with multiple blockchains.
type CrossChainSmartContract struct {
	ID              string                 // Unique ID of the contract
	Code            string                 // Smart contract code (bytecode or source)
	Parameters      map[string]interface{} // Contract parameters
	State           map[string]interface{} // Current state of the contract
	Owner           string                 // Owner of the smart contract
	ConnectedChains []string               // List of connected blockchains
	Executions      []ContractExecution    // History of contract executions
	mutex           sync.Mutex             // Mutex for concurrency safety
	LedgerInstance  *Ledger                // Ledger for storing contract data
}

// CrossChainContractManager manages cross-chain smart contracts.
type CrossChainContractManager struct {
	Contracts      map[string]*CrossChainSmartContract // All deployed cross-chain contracts
	LedgerInstance *Ledger                             // Ledger for contract deployments and executions
	mutex          sync.Mutex                          // Mutex for concurrency safety
}

// MigrationManager handles smart contract migrations, ensuring seamless transitions between contract versions.
type MigrationManager struct {
	Contracts      map[string]*MigratedContract // A map to track migrated contracts by ID
	LedgerInstance *Ledger                      // Ledger instance to store migration data
	mutex          sync.Mutex                   // Mutex for thread-safe operations
}

// RicardianContract represents a Ricardian contract in the Synnergy blockchain.
type RicardianContract struct {
	ID              string                 // Unique ID of the contract
	HumanReadable   string                 // Human-readable legal terms of the contract
	MachineReadable string                 // Machine-readable executable code
	PartiesInvolved []string               // Parties involved in the contract
	Signatures      map[string]string      // Digital signatures of the involved parties
	State           map[string]interface{} // Current state of the contract
	Owner           string                 // Owner or issuer of the contract
	Executions      []ContractExecution    // History of contract executions
	mutex           sync.Mutex             // Mutex for safe concurrency
	LedgerInstance  *Ledger                // Ledger instance for storing contract data
}

// RicardianContractManager manages multiple Ricardian contracts.
type RicardianContractManager struct {
	Contracts      map[string]*RicardianContract // All deployed Ricardian contracts
	LedgerInstance *Ledger                       // Ledger for recording contract deployments and executions
	mutex          sync.Mutex                    // Mutex for safe concurrency
}

// SmartContractTemplate represents a template for a smart contract in the marketplace.
type SmartContractTemplate struct {
	ID            string    // Unique ID of the template
	Name          string    // Template name
	Description   string    // Template description
	Creator       string    // Creator of the template
	Code          string    // Contract code or bytecode
	Price         float64   // Price to purchase the template
	Timestamp     time.Time // Time when the template was created
	EncryptedCode string    // Encrypted code for security
}

// SmartContractTemplateMarketplace manages the marketplace for smart contract templates.
type SmartContractTemplateMarketplace struct {
	Templates      map[string]*SmartContractTemplate // Map of templates by ID
	Escrows        map[string]*Escrow                // Escrows for purchases
	EscrowFee      float64                           // Marketplace fee (in percentage)
	mutex          sync.Mutex                        // Mutex for safe concurrency
	LedgerInstance *Ledger                           // Ledger for recording transactions
}

// SmartContractManager manages multiple smart contracts.
type SmartContractManager struct {
	Contracts      map[string]*SmartContract // All deployed smart contracts
	LedgerInstance *Ledger                   // Ledger for recording contract deployments and executions
	mutex          sync.Mutex                // Mutex for safe concurrency
}

// SmartLegalContract represents a legal contract in the Synnergy blockchain.
type SmartLegalContract struct {
	ID              string                 // Unique ID of the legal contract
	ContractTerms   string                 // The terms of the legal contract
	PartiesInvolved []string               // Parties involved in the contract (addresses)
	Signatures      map[string]string      // Digital signatures of the involved parties
	State           map[string]interface{} // Current state of the contract
	Owner           string                 // Owner or issuer of the legal contract
	Executions      []ContractExecution    // History of contract executions
	LegallyBinding  bool                   // Indicates if the contract is legally binding
	mutex           sync.Mutex             // Mutex for safe concurrency
	LedgerInstance  *Ledger                // Ledger instance for storing contract data
}

// ContractExecution represents an execution instance of a contract function.
type ContractExecution struct {
	ExecutionID   string                 // Unique identifier for the execution
	Executor      string                 // Address of the entity executing the contract
	ContractID    string                 // ID of the contract being executed
	FunctionName  string                 // Name of the function executed
	Parameters    map[string]interface{} // Parameters passed during the execution
	Timestamp     time.Time              // Time when the execution occurred
	Result        map[string]interface{} // The result of the execution (success, error, etc.)
	ExecutionTime time.Time              // Time of the execution
	GasUsed       float64                // Amount of gas used for the execution
}

// MigratedContract represents a migrated contract from one version to another.
type MigratedContract struct {
	OldContractID string                 // ID of the old contract
	NewContractID string                 // ID of the new contract after migration
	MigrationTime time.Time              // Time when the migration occurred
	Reason        string                 // Reason for migration (e.g., upgrade, bug fix)
	NewState      map[string]interface{} // State of the new contract post-migration
}

// ContractStorage represents the storage data for a smart contract.
type ContractStorage struct {
	ContractAddress string
	Data            map[string]string
	Timestamp       time.Time
}

// ************** Space-Time Proof Structs **************

// SpaceTimeProof represents a proof that links storage to the space-time consensus.
type SpaceTimeProof struct {
	ProofID        string    // Unique proof ID
	StorageID      string    // Associated storage ID
	Validator      string    // Address of the validator who validated the proof
	ValidationTime time.Time // When the proof was validated
	Status         string    // Status of the proof (valid, invalid, revalidated)
}

// SpaceTimeProofInvalidation represents the invalidation details of a space-time proof.
type SpaceTimeProofInvalidation struct {
	ProofID     string    // ID of the invalidated proof
	Invalidator string    // Entity or person who invalidated the proof
	Reason      string    // Reason for invalidation
	Timestamp   time.Time // Time of invalidation
}

// STValidationResult represents the result of a space-time proof validation.
type STValidationResult struct {
	ProofID          string    // Unique identifier for the space-time proof
	ValidatorID      string    // ID of the validator who performed the validation
	ValidationTime   time.Time // Timestamp of the validation process
	ValidationStatus string    // Status of the validation (e.g., "Valid", "Invalid")
	ValidationLogs   string    // Logs or details of the validation process
	ValidationScore  float64   // Score or metric for the validation, if applicable
}

// STInvalidationRecord tracks details of an invalidated space-time proof.
type STInvalidationRecord struct {
	ProofID            string    // Unique identifier for the invalidated space-time proof
	Reason             string    // Reason for the invalidation
	InvalidatedBy      string    // ID of the entity or validator responsible for the invalidation
	InvalidationTime   time.Time // Timestamp of when the proof was invalidated
	ResolutionRequired bool      // Indicates if resolution is required for this invalidation
	ResolutionDeadline time.Time // Deadline for resolving the invalidation, if applicable
}

// STRevalidationRecord tracks attempts to revalidate a previously invalidated space-time proof.
type STRevalidationRecord struct {
	ProofID            string    // Unique identifier for the space-time proof being revalidated
	RevalidationTime   time.Time // Timestamp of the revalidation process
	RevalidationStatus string    // Status of the revalidation (e.g., "Successful", "Failed")
	RevalidatedBy      string    // ID of the entity or validator who performed the revalidation
	RevalidationLogs   string    // Logs or details of the revalidation process
	RetryCount         int       // Number of revalidation attempts made for this proof
	NextRetryAllowed   time.Time // Timestamp indicating when the next retry is allowed, if applicable
}

// SpaceTimeProofRecord represents a record of proof validation, invalidation, or revalidation.
type SpaceTimeProofRecord struct {
	ProofID   string
	Validator string
	Timestamp time.Time
	Status    string // Can be "validated", "invalidated", "revalidated"
}

// ************** State-Channel Structs **************

// StateChannel represents a state channel in the ledger.
type StateChannel struct {
	ChannelID       string                       // Unique ID for the state channel
	Participants    []string                     // Addresses of participants in the state channel
	Status          string                       // Status of the channel (open, closed, etc.)
	CreatedAt       time.Time                    // Timestamp when the channel was opened
	ClosedAt        time.Time                    // Timestamp when the channel was closed (if applicable)
	State           string                       // Current state of the channel
	Transactions    []TransactionRecord          // List of transactions within the state channel
	DataTransfers   map[string]*DataBlock        // Data transfers in the channel
	Flexibility     float64                      // Flexibility setting for the state channel
	LastValidatedBy string                       // Last validator who validated the channel
	LastValidatedAt time.Time                    // Timestamp of the last validation
	LastSyncedAt    time.Time                    // Timestamp when the channel was last synced
	IsOpen          bool                         // Indicates if the channel is open or closed
	Collateral      map[string]*CollateralRecord // Map of participant addresses to their collateral records

}

// KeySharing represents a log of an encryption key being shared between entities.
type KeySharing struct {
	EntityID   string    // The entity sharing the key
	SharedWith string    // The entity with whom the key is shared
	Key        string    // The encryption key being shared
	SharedAt   time.Time // The timestamp when the key was shared
}

// PerformanceMetrics represents the performance metrics for a state channel or entity.
type PerformanceMetrics struct {
	EntityID  string    // The entity ID (channel or shard)
	Metrics   string    // Performance metrics in a string format
	Timestamp time.Time // Timestamp of when the metrics were recorded
}

// LoadMetric represents the load metrics associated with a state channel.
type LoadMetric struct {
	ChannelID string    // ID of the state channel
	Load      int       // Current load metric of the channel (e.g., number of transactions)
	Timestamp time.Time // Timestamp of the load measurement
}

// ScalingEvent represents an event related to scaling the state channels.
type ScalingEvent struct {
	EventID      string    // Unique ID for the scaling event
	ChannelID    string    // Channel ID being scaled
	OldResources int       // Previous resource allocation
	NewResources int       // New resource allocation
	Timestamp    time.Time // Timestamp of the scaling event
}

// FragmentedState represents a fragmented state for load balancing.
type FragmentedState struct {
	FragmentID   string    // Unique ID of the state fragment
	ChannelID    string    // Channel ID that the state belongs to
	FragmentData string    // Data contained in the fragment
	Timestamp    time.Time // Timestamp when the fragment was created
}

// AdaptiveLoadBalancingChannel represents a state channel with dynamic load-balancing capabilities
type AdaptiveLoadBalancingChannel struct {
	ChannelID      string                 // Unique identifier for the adaptive load-balancing channel
	Participants   []string               // Participants in the state channel
	State          map[string]interface{} // Current state of the channel
	LoadMetrics    map[string]float64     // Tracks load metrics for each participant or node
	IsOpen         bool                   // Indicates if the channel is open or closed
	Ledger         *Ledger                // Reference to the ledger for recording events
	Encryption     *Encryption            // Encryption service for securing state data
	NetworkManager *NetworkManager        // Network manager for distributing load
	mu             sync.Mutex             // Mutex for concurrency control
}

// ContinuousTimeScalingChannel (CTSC) represents a state channel with real-time, continuous optimization.
type ContinuousTimeScalingChannel struct {
	ChannelID      string                 // Unique identifier for the state channel
	Participants   []string               // Addresses of participants in the channel
	State          map[string]interface{} // Current state of the channel
	IsOpen         bool                   // Whether the channel is open or closed
	Ledger         *Ledger                // Reference to the ledger for recording events
	Encryption     *Encryption            // Encryption service for securing state data
	NetworkManager *NetworkManager        // Network manager to handle communications and scaling
	ScalingFactor  float64                // Current scaling factor for channel operations
	mu             sync.Mutex             // Mutex for concurrency control
}

// DynamicLoadBalancingChannel represents a state channel with dynamic load-balancing capabilities
type DynamicLoadBalancingChannel struct {
	ChannelID      string                 // Unique identifier for the state channel
	Participants   []string               // Participants in the state channel
	State          map[string]interface{} // Current state of the channel
	IsOpen         bool                   // Indicates if the channel is open or closed
	Ledger         *Ledger                // Reference to the ledger for logging events
	Encryption     *Encryption            // Encryption service for securing state data
	NetworkManager *NetworkManager        // Network manager to handle communication and load balancing
	LoadThreshold  float64                // Threshold beyond which load-balancing is triggered
	CurrentLoad    float64                // Current load of the channel
	mu             sync.Mutex             // Mutex for concurrency control
}

// DynamicResourceAllocationChannel (DRAC) represents a state channel with dynamic resource allocation
type DynamicResourceAllocationChannel struct {
	ChannelID      string                 // Unique identifier for the state channel
	Participants   []string               // Addresses of participants in the channel
	State          map[string]interface{} // Current state of the channel
	IsOpen         bool                   // Whether the channel is open or closed
	Ledger         *Ledger                // Reference to the ledger for recording events
	Encryption     *Encryption            // Encryption service for securing state data
	NetworkManager *NetworkManager        // Network manager to handle communication and resource allocation
	ResourceUsage  float64                // Current resource usage of the channel
	LoadThreshold  float64                // Threshold beyond which resource reallocation is triggered
	mu             sync.Mutex             // Mutex for concurrency control
}

// StateFragment represents a piece of fragmented state data
type StateFragment struct {
	FragmentID    string    // Unique identifier for the state fragment
	ChannelID     string    // State channel to which this fragment belongs
	FragmentIndex int       // Index of the fragment in the fragmented state data
	Data          string    // Encrypted fragment data
	Timestamp     time.Time // Time when the fragment was created
}

// DynamicStateFragmentationChannel represents a state channel with dynamic state fragmentation
type DynamicStateFragmentationChannel struct {
	ChannelID      string                 // Unique identifier for the state channel
	Participants   []string               // Participants in the state channel
	Fragments      map[int]*StateFragment // Map of state fragments, indexed by fragment index
	State          map[string]interface{} // General state of the channel
	FragmentCount  int                    // Total number of state fragments
	Ledger         *Ledger                // Reference to the ledger for recording transactions
	Encryption     *Encryption            // Encryption service for fragment encryption
	NetworkManager *NetworkManager        // Network manager to handle fragment distribution
	mu             sync.Mutex             // Mutex for concurrency control
}

// CollateralRecord represents a participant's collateral in the state channel
type CollateralRecord struct {
	Participant      string    // The participant's address
	Collateral       float64   // The amount of collateral deposited
	LastUpdated      time.Time // Timestamp of the last update to the collateral
	Amount           float64   // Amount of collateral deposited
	ValidationStatus string    // Status of collateral validation
	Validator        string    // Address of the validator
	ValidatedAt      time.Time // Time when the collateral was validated
}

// ElasticCollateralChannel represents a state channel with adaptive collateral management
type ElasticCollateralChannel struct {
	ChannelID      string                       // Unique identifier for the state channel
	Participants   map[string]*CollateralRecord // Map of participants and their collateral records
	State          map[string]interface{}       // General state of the channel
	IsOpen         bool                         // Indicates if the channel is open or closed
	Ledger         *Ledger                      // Reference to the ledger for recording transactions
	Encryption     *Encryption                  // Encryption service for securing collateral data
	NetworkManager *NetworkManager              // Network manager for managing participant communications
	CollateralCap  float64                      // Maximum allowed collateral for the channel
	mu             sync.Mutex                   // Mutex for concurrency control
}

// FlexibleStateChannel represents a state channel with enhanced flexibility for dynamic operations
type FlexibleStateChannel struct {
	ChannelID      string                  // Unique identifier for the state channel
	Participants   []string                // Addresses of participants in the channel
	State          map[string]interface{}  // Current state of the channel
	DataTransfers  map[string]*DataBlock   // Data transfers in the channel
	Liquidity      map[string]float64      // Liquidity holdings of participants
	Transactions   map[string]*Transaction // Transactions in the channel
	IsOpen         bool                    // Whether the channel is open or closed
	Flexibility    float64                 // Flexibility factor for dynamic channel adjustments
	Ledger         *Ledger                 // Reference to the ledger for recording events
	Encryption     *Encryption             // Encryption service for securing state data
	NetworkManager *NetworkManager         // Network manager for managing participant communications
	mu             sync.Mutex              // Mutex for concurrency control
}

// FluidStateSyncChannel represents a state channel with real-time state synchronization across channels
type FluidStateSyncChannel struct {
	ChannelID      string                   // Unique identifier for the fluid state channel
	Participants   []string                 // Addresses of participants in the channel
	State          map[string]interface{}   // Current state of the channel
	IsOpen         bool                     // Whether the channel is open or closed
	SyncedChannels []*FluidStateSyncChannel // Other channels to synchronize state with
	Ledger         *Ledger                  // Reference to the ledger for recording state events
	Encryption     *Encryption              // Encryption service for securing state data
	NetworkManager *NetworkManager          // Network manager for real-time state synchronization
	mu             sync.Mutex               // Mutex for concurrency control
}

// FractalStateChannel represents a recursive state channel with state aggregation and fractal syncing
type FractalStateChannel struct {
	ChannelID      string                 // Unique identifier for the fractal state channel
	Participants   []string               // Addresses of participants in the channel
	State          map[string]interface{} // Current state of the fractal channel
	SubChannels    []*FractalStateChannel // Nested fractal sub-channels
	IsOpen         bool                   // Whether the channel is open or closed
	Ledger         *Ledger                // Reference to the ledger for recording events
	Encryption     *Encryption            // Encryption service for securing state data
	NetworkManager *NetworkManager        // Network manager for recursive state synchronization
	mu             sync.Mutex             // Mutex for concurrency control
}

// HierarchicalShardedStateChannel represents a hierarchical sharded state channel with recursive sharding
type HierarchicalShardedStateChannel struct {
	ChannelID      string                 // Unique identifier for the hierarchical state channel
	Participants   []string               // Participants in the main state channel
	State          map[string]interface{} // State of the main channel
	Shards         map[string]*Shard      // Shards within the main channel
	IsOpen         bool                   // Whether the channel is open or closed
	Ledger         *Ledger                // Reference to the ledger for event recording
	Encryption     *Encryption            // Encryption service to secure state and shard data
	NetworkManager *NetworkManager        // Network manager for state syncing and shard communications
	mu             sync.Mutex             // Mutex for concurrency control
}

// InstantFinalityChannel represents a state channel with near-instant consensus and state finality
type InstantFinalityChannel struct {
	ChannelID      string                 // Unique identifier for the state channel
	Participants   []string               // Addresses of participants in the channel
	State          map[string]interface{} // Current state of the channel
	Transactions   []*Transaction         // Transactions within the channel
	IsOpen         bool                   // Whether the channel is open or closed
	Ledger         *Ledger                // Reference to the ledger for recording events
	Encryption     *Encryption            // Encryption service for securing state data
	NetworkManager *NetworkManager        // Network manager to handle communications
	Finalized      bool                   // Whether the channel has achieved finality
	mu             sync.Mutex             // Mutex for concurrency control
}

// IdentityVerificationChannel (IVC) represents a state channel for identity verification
type IdentityVerificationChannel struct {
	ChannelID      string                    // Unique identifier for the identity verification channel
	Participants   []string                  // Addresses of participants in the channel
	IdentityProofs map[string]*IdentityProof // Mapping of participant addresses to their identity proofs
	IsOpen         bool                      // Whether the channel is open or closed
	Ledger         *Ledger                   // Reference to the ledger for recording events
	Encryption     *Encryption               // Encryption service for securing identity data
	NetworkManager *NetworkManager           // Network manager for communications and verification
	mu             sync.Mutex                // Mutex for concurrency control
}

// IdentityProof represents a participant's proof of identity
type IdentityProof struct {
	Participant string    // Address of the participant
	ProofHash   string    // Hash of the identity proof document
	Timestamp   time.Time // Timestamp when the proof was submitted
	Verified    bool      // Whether the proof has been verified
}

// RealTimeShardReallocationChannel (RTSR) represents a state channel with dynamic shard reallocation
type RealTimeShardReallocationChannel struct {
	ChannelID      string            // Unique identifier for the channel
	Shards         map[string]*Shard // Shards allocated to this channel
	Participants   []string          // Addresses of participants in the channel
	IsOpen         bool              // Whether the channel is open or closed
	Ledger         *Ledger           // Reference to the ledger for recording events
	Encryption     *Encryption       // Encryption service for securing shard data
	NetworkManager *NetworkManager   // Network manager for communication and shard reallocation
	mu             sync.Mutex        // Mutex for concurrency control
}

// DataChannel represents a state channel for managing data transfers between participants
type DataChannel struct {
	ChannelID     string                 // Unique identifier for the data channel
	Participants  []string               // Addresses of participants in the channel
	DataState     map[string]interface{} // Current state of data in the channel
	DataTransfers map[string]*DataBlock  // Blocks of data transferred within the channel
	IsOpen        bool                   // Whether the channel is open or closed
	Ledger        *Ledger                // Reference to the ledger for recording data transfer events
	Encryption    *Encryption            // Encryption service for securing data
	mu            sync.Mutex             // Mutex for handling concurrency
}

// DataBlock represents a block of data transferred between participants
type DataBlock struct {
	BlockID    string    // Unique identifier for the data block
	Data       string    // Data content being transferred
	Timestamp  time.Time // Timestamp when the data block was created
	MerkleRoot string    // Merkle root of the data for validation
}

// StateChannelInteroperability handles interoperability between different state channels
type StateChannelInteroperability struct {
	ChannelID                 string                  // Unique identifier for the state channel
	LinkedChannels            map[string]*InteropLink // Linked state channels for interoperability
	IsInteroperabilityEnabled bool                    // Indicates whether cross-channel operations are allowed
	Participants              []string                // Addresses of participants in the channel
	State                     map[string]interface{}  // Current state of the channel
	Ledger                    *Ledger                 // Reference to the ledger for recording interoperability events
	Encryption                *Encryption             // Encryption service for securing data exchanges
	mu                        sync.Mutex              // Mutex for concurrency control
	NetworkManager            *NetworkManager         // Handles cross-network communication for state channels
}

// InteropLink represents a link between the current channel and another interoperable state channel
type InteropLink struct {
	LinkID          string                 // Unique identifier for the interoperability link
	TargetChannelID string                 // The ID of the linked state channel
	SharedData      map[string]interface{} // Data shared between channels
	Timestamp       time.Time              // Timestamp when the interoperability was established
}

// LiquidityStateChannel represents a state channel for managing liquidity between participants
type LiquidityStateChannel struct {
	ChannelID    string                 // Unique identifier for the liquidity channel
	Participants []string               // Participants in the liquidity pool
	Liquidity    map[string]float64     // Mapping of participants to their liquidity contributions
	State        map[string]interface{} // General state of the channel
	IsOpen       bool                   // Indicates if the channel is open or closed
	Ledger       *Ledger                // Reference to the ledger for recording transactions
	Encryption   *Encryption            // Encryption service for securing liquidity data
	mu           sync.Mutex             // Mutex for concurrency control
}

// MicroChannelMode represents the mode that enables a state channel to act as a microchannel
type MicroChannelMode struct {
	Enabled        bool          // Whether the state channel is in microchannel mode
	MaxTransaction int           // Maximum number of transactions before closing the channel
	Timeout        time.Duration // Time before automatically closing the microchannel
	StartTime      time.Time     // Time when the microchannel mode was enabled
}

// OffChainSettlementChannel represents a state channel for managing off-chain settlements between participants
type OffChainSettlementChannel struct {
	ChannelID    string                 // Unique identifier for the settlement channel
	Participants []string               // Addresses of participants in the channel
	State        map[string]interface{} // Current state of the channel
	Balances     map[string]float64     // Mapping of participant addresses to their balances
	IsOpen       bool                   // Whether the channel is open or closed
	Ledger       *Ledger                // Reference to the ledger for recording settlement events
	Encryption   *Encryption            // Encryption service for securing settlement data
	mu           sync.Mutex             // Mutex for handling concurrency
}

// PaymentStateChannel represents a state channel for managing payment transactions between participants
type PaymentStateChannel struct {
	ChannelID    string                 // Unique identifier for the payment channel
	Participants []string               // Addresses of participants in the payment channel
	Balances     map[string]float64     // Mapping of participant addresses to their current balances
	State        map[string]interface{} // General state of the channel
	IsOpen       bool                   // Whether the payment channel is open or closed
	Ledger       *Ledger                // Reference to the ledger for recording payment events
	Encryption   *Encryption            // Encryption service for securing payment data
	mu           sync.Mutex             // Mutex for handling concurrency
}

// StateChannelPerformance handles performance tracking and optimization for a state channel
type StateChannelPerformance struct {
	ChannelID        string                 // Unique identifier for the state channel
	Participants     []string               // Participants in the state channel
	TransactionTimes []time.Duration        // Array to store transaction times for performance metrics
	Throughput       float64                // Throughput of transactions (transactions per second)
	Latency          time.Duration          // Average latency in processing transactions
	State            map[string]interface{} // General state of the channel
	Ledger           *Ledger                // Reference to the ledger for logging performance data
	Encryption       *Encryption            // Encryption service for securing performance data
	mu               sync.Mutex             // Mutex for concurrency handling
}

// PrivacyStateChannel ensures the privacy and confidentiality of state channel data
type PrivacyStateChannel struct {
	ChannelID      string            // Unique identifier for the state channel
	Participants   []string          // Participants in the state channel
	EncryptedState map[string]string // Encrypted state data of the channel
	IsOpen         bool              // Whether the channel is open or closed
	Ledger         *Ledger           // Reference to the ledger for logging privacy-related events
	Encryption     *Encryption       // Encryption service for securing data
	DecryptionKeys map[string]string // Store participants' decryption keys
	mu             sync.Mutex        // Mutex for concurrency control
}

// SecurityStateChannel represents a secure state channel with enhanced security measures
type SecurityStateChannel struct {
	ChannelID       string                 // Unique identifier for the state channel
	Participants    []string               // Participants in the state channel
	State           map[string]interface{} // Current state of the channel
	IsOpen          bool                   // Whether the channel is open or closed
	Ledger          *Ledger                // Reference to the ledger for logging events
	Encryption      *Encryption            // Encryption service for securing state data
	SecurityModule  *StateChannelSecurity  // Security module for additional authentication and verification
	ParticipantKeys map[string]string      // Public keys of participants for signature verification
	mu              sync.Mutex             // Mutex for concurrency control
}

// SmartContractStateChannel represents a state channel with smart contract functionality
type SmartContractStateChannel struct {
	ChannelID      string                    // Unique identifier for the state channel
	Participants   []string                  // Participants in the channel
	State          map[string]interface{}    // Current state of the channel
	SmartContracts map[string]*SmartContract // Deployed smart contracts in the channel
	IsOpen         bool                      // Whether the channel is open or closed
	Ledger         *Ledger                   // Reference to the ledger for logging
	Encryption     *Encryption               // Encryption service for securing contract and state data
	mu             sync.Mutex                // Mutex for concurrency control
}

// StateChannelSecurity represents the security module for the state channel, handling authentication and verification.
type StateChannelSecurity struct {
	AuthProtocols    map[string]SecurityProtocol // Authentication protocols applied (e.g., multi-factor authentication)
	VerificationKeys map[string]string           // Public keys for verifying participant identities
	EncryptionKeys   map[string]string           // Encryption keys for securing communication
	ActiveThreats    []string                    // List of active security threats or attacks being mitigated
	SecurityLogs     []*SecurityLog              // Logs of security events and actions
	Ledger           *Ledger                     // Ledger reference for recording security-related events
	mu               sync.Mutex                  // Mutex for concurrent security operations
}

// SecurityProtocol defines different security protocols for authentication and verification.
type SecurityProtocol struct {
	ProtocolID   string    // Unique identifier for the security protocol
	ProtocolName string    // Name of the protocol (e.g., MFA, signature verification)
	IsActive     bool      // Indicates if the protocol is currently active
	Timestamp    time.Time // Timestamp of the last protocol update or usage
}

// SecurityLog represents a log entry for a security-related event.
type SecurityLog struct {
	LogID       string    // Unique identifier for the log entry
	EventType   string    // Type of event (e.g., authentication attempt, threat detected)
	Participant string    // Participant involved in the event
	Timestamp   time.Time // Time when the event occurred
	Description string    // Detailed description of the event
	IsResolved  bool      // Whether the event or issue has been resolved
}

// ************** Storage Structs **************

// CacheRecord represents a single cache entry in the ledger.
type CacheRecord struct {
	Key       string    // Cache key
	Value     string    // Cached value
	AddedAt   time.Time // Time the cache was added
	ExpiresAt time.Time // Expiration time for the cache entry
}

// FileOperation represents a file operation in the ledger.
type FileOperation struct {
	OperationID string    // Unique ID for the file operation
	FilePath    string    // Path to the file being operated on
	Action      string    // Action taken on the file (e.g., read, write, delete)
	Timestamp   time.Time // Timestamp of the file operation
}

// StorageEvent represents events where storage-related actions occur.
type StorageEvent struct {
	EventID   string    // Unique storage event ID
	StorageID string    // Associated storage ID
	Action    string    // Action performed (e.g., storage added, updated)
	Details   string    // Additional details about the event
	Timestamp time.Time // Timestamp of the storage event
}

// SystemCacheManager manages system-wide caching for improved performance and resource optimization.
type SystemCacheManager struct {
	CacheID          string                // Unique identifier for the cache instance.
	CacheSize        int64                 // Total size of the cache in bytes.
	CachePolicy      string                // Policy for cache management (e.g., LRU, LFU, FIFO).
	ActiveCaches     map[string]CacheEntry // Active cache entries indexed by key.
	CacheHitCount    int                   // Number of successful cache hits.
	CacheMissCount   int                   // Number of cache misses.
	LastEvictionTime time.Time             // Timestamp of the last cache eviction.
	Logs             []CacheLog            // Logs related to cache operations.
}

// CacheEntry represents a single cache item.
type CacheEntry struct {
	Key        string    // Unique key for the cache entry.
	Value      []byte    // Cached value.
	Size       int64     // Size of the cached value in bytes.
	CreatedAt  time.Time // Timestamp when the cache entry was created.
	LastAccess time.Time // Timestamp of the last access.
}

// CacheLog captures events or updates in the caching process.
type CacheLog struct {
	LogID     string    // Unique identifier for the log entry.
	CacheID   string    // Associated cache instance ID.
	Message   string    // Log message.
	Timestamp time.Time // Timestamp of the log entry.
}

// MemoryManager oversees system memory allocation, usage, and optimization.
type MemoryManager struct {
	TotalMemory        int64              // Total memory available in bytes.
	UsedMemory         int64              // Memory currently in use in bytes.
	FreeMemory         int64              // Memory available for allocation in bytes.
	MemoryPools        []MemoryPool       // List of memory pools managed by the system.
	AllocationPolicies []AllocationPolicy // Policies governing memory allocation and usage.
	LastGarbageCollect time.Time          // Timestamp of the last garbage collection.
	Logs               []MemoryLog        // Logs related to memory operations.
}

// MemoryPool represents a segment of memory allocated for specific purposes.
type MemoryPool struct {
	PoolID      string    // Unique identifier for the memory pool.
	PoolSize    int64     // Total size of the memory pool in bytes.
	UsedSize    int64     // Size of memory currently in use in the pool.
	Allocations int       // Number of active memory allocations in the pool.
	PoolType    string    // Type of memory pool (e.g., stack, heap, shared).
	LastAccess  time.Time // Timestamp of the last access to this pool.
}

// AllocationPolicy defines policies for memory allocation and usage.
type AllocationPolicy struct {
	PolicyID    string                 // Unique identifier for the policy.
	Description string                 // Description of the allocation policy.
	Rules       map[string]interface{} // Rules governing memory allocation and optimization.
}

// MemoryLog captures memory-related events or updates.
type MemoryLog struct {
	LogID     string    // Unique identifier for the log entry.
	Message   string    // Log message.
	Timestamp time.Time // Timestamp of the log entry.
}

// CacheManager handles the caching mechanism for blockchain operations
type CacheManager struct {
	CacheEntries   map[string]*CacheEntry // Cache entries stored by key
	mutex          sync.Mutex             // Mutex for thread-safe operations
	LedgerInstance *Ledger                // Ledger for recording cache activity
	MaxCacheSize   int                    // Maximum number of cache entries allowed
	DefaultTTL     time.Duration          // Time-to-live for cache entries
}

// FileEntry represents a file stored in the system
type FileEntry struct {
	FileName     string    // The original file name
	FilePath     string    // The full path of the file in the storage system
	Encrypted    bool      // Whether the file is encrypted
	UploadedAt   time.Time // The timestamp when the file was uploaded
	LastAccessed time.Time // The timestamp of the last access
}

// FileManager manages file storage, retrieval, and encryption
type FileManager struct {
	Files          map[string]*FileEntry // A map of file hashes to FileEntry
	storageDir     string                // Directory where files are stored
	mutex          sync.Mutex            // Mutex for thread-safe operations
	LedgerInstance *Ledger               // Ledger for tracking file operations
}

// IPFSManager handles file storage and retrieval from IPFS with encryption
type IPFSManager struct {
	mutex          sync.Mutex // Mutex for thread-safe operations
	LedgerInstance *Ledger    // Ledger for tracking IPFS file operations
}

// FileIndex represents metadata for files stored in the system
type FileIndex struct {
	FileID       string    // Unique file identifier (encrypted)
	FileName     string    // Original file name
	FileSize     int64     // Size of the file in bytes
	UploadedAt   time.Time // Timestamp of when the file was uploaded
	Owner        string    // Owner address or user who uploaded the file
	EncryptedCID string    // Encrypted IPFS CID for file retrieval
}

// FileIndexer manages file indexing and metadata
type FileIndexer struct {
	mutex          sync.Mutex            // Mutex for thread-safe operations
	Indexes        map[string]*FileIndex // Indexed file metadata by FileID
	LedgerInstance *Ledger               // Ledger instance for tracking file operations
}

// EscrowAccount represents an escrow account to hold funds temporarily during a transaction
type EscrowAccount struct {
	EscrowID   string    // Unique escrow account ID
	Buyer      string    // Buyer involved in the transaction
	Seller     string    // Seller involved in the transaction
	Amount     float64   // Amount of funds in escrow
	CreatedAt  time.Time // Timestamp of when the escrow account was created
	IsReleased bool      // Whether the funds have been released
}

// StorageListing represents a storage listing on the marketplace
type StorageListing struct {
	ListingID        string    // Unique ID of the listing
	Owner            string    // Owner or creator of the listing
	CapacityGB       int       // Storage capacity in GB
	PricePerGB       float64   // Price per GB in SYNN
	LeaseDuration    int       // Lease duration in days
	PostedAt         time.Time // Timestamp of when the listing was created
	EncryptedDetails string    // Encrypted storage details (description, terms, etc.)
	Active           bool      // Whether the listing is active
}

// StorageMarketplace manages the listing and leasing of storage on the blockchain
type StorageMarketplace struct {
	mutex          sync.Mutex                 // Mutex for thread-safe operations
	Listings       map[string]*StorageListing // Map of ListingID to StorageListing
	LedgerInstance *Ledger                    // Ledger instance for tracking storage transactions
	EscrowAccounts map[string]*EscrowAccount  // Map of EscrowID to EscrowAccount
}

// OffchainStorage represents an off-chain storage unit
type OffchainStorage struct {
	StorageID        string // Unique identifier for the off-chain storage unit
	Owner            string // Owner or manager of the off-chain storage
	Location         string // Physical or cloud location of the storage
	CapacityGB       int    // Total capacity of the storage in GB
	UsedCapacityGB   int    // Capacity used in GB
	EncryptedDetails string // Encrypted details of the storage
	Active           bool   // Whether the storage is active or not
}

// OffchainStorageManager manages off-chain storage listings
type OffchainStorageManager struct {
	mutex          sync.Mutex                  // Mutex for thread-safe operations
	StorageUnits   map[string]*OffchainStorage // Map of storage units
	LedgerInstance *Ledger                     // Ledger instance for tracking storage transactions
}

// StorageRetrievalManager handles the retrieval of on-chain and off-chain storage data
type StorageRetrievalManager struct {
	OnChainData    map[string][]byte           // On-chain data stored by transaction hash
	OffChainUnits  map[string]*OffchainStorage // Map of off-chain storage units
	LedgerInstance *Ledger                     // Ledger instance for verification
	mutex          sync.Mutex                  // Mutex for thread-safe operations
}

// StorageSanitizationManager handles sanitizing data before storage and removing sensitive or irrelevant data from the system
type StorageSanitizationManager struct {
	LedgerInstance *Ledger    // Ledger instance for recording sanitized data transactions
	mutex          sync.Mutex // Mutex for thread-safe operations
}

// StorageManager manages the creation, retrieval, and deletion of storage entries
type StorageManager struct {
	LedgerInstance *Ledger                 // Ledger instance for transaction recording
	mutex          sync.Mutex              // Mutex for thread-safe operations
	StorageMap     map[string]StorageEntry // Map of storage entries, identified by storage ID
}

// SwarmManager manages the interaction with the decentralized storage swarm
type SwarmManager struct {
	LedgerInstance *Ledger                 // Ledger instance for transaction logging
	mutex          sync.Mutex              // Mutex for thread-safe operations
	SwarmNodes     map[string]SwarmNode    // Map of swarm nodes connected to the system
	StorageMap     map[string]StorageEntry // Map of storage entries, identified by storage ID
}

// TimestampManager manages the process of timestamping data and logging it to the ledger
type TimestampManager struct {
	LedgerInstance *Ledger // Ledger for logging timestamped data
}

// StorageEntry represents a storage entry in the system
type StorageEntry struct {
	StorageID    string    // Unique identifier for the storage entry
	Data         []byte    // Encrypted data stored in the entry
	Owner        string    // Owner of the storage entry (address or ID)
	CreatedAt    time.Time // Timestamp of when the storage entry was created
	LastAccessed time.Time // Timestamp of when the storage entry was last accessed
	Expiration   time.Time // Expiration time for the storage entry (optional)
	IsActive     bool      // Indicates if the storage entry is active or expired
}

// SwarmNode represents a node in the decentralized storage swarm
type SwarmNode struct {
	NodeID         string            // Unique identifier for the swarm node
	IPAddress      string            // IP address of the swarm node
	Location       string            // Physical or cloud location of the swarm node
	CapacityGB     int               // Total storage capacity in GB
	UsedCapacityGB int               // Used storage capacity in GB
	Status         string            // Status of the node (active, inactive, etc.)
	LastActive     time.Time         // Timestamp of the last time the node was active
	EncryptedData  map[string][]byte // Encrypted data stored on the node
}

// ************** Sustainability Structs **************

// CarbonCredit represents a carbon credit in the ledger.
type CarbonCredit struct {
	CreditID  string    // Unique ID for the carbon credit
	Issuer    string    // The issuer of the carbon credit
	Owner     string    // The current owner of the credit
	Amount    float64   // Amount of carbon offset (in tons)
	Status    string    // Status of the credit (active, retired)
	IssuedAt  time.Time // Time when the credit was issued
	RetiredAt time.Time // Time when the credit was retired (if applicable)
}

// SustainabilityState represents the overall state of sustainability efforts.
type SustainabilityState struct {
	TotalCarbonCredits   float64           // Total carbon credits in the system
	AvailableCredits     float64           // Total available carbon credits for offset
	ActivePrograms       int               // Count of active sustainability programs
	CertifiedEntities    int               // Count of entities with eco-certifications
	OptimizationProgress map[string]string // Tracks optimization progress by entity
	LastUpdated          time.Time         // Timestamp of the last state update
}

// EcoFriendlyCertificate represents a certificate awarded for eco-friendly practices.
type EcoFriendlyCertificate struct {
	CertificateID      string     // Unique ID for the certificate
	Recipient          string     // Entity or individual receiving the certificate
	AwardedBy          string     // Issuing authority of the certificate
	AwardDate          time.Time  // Date of awarding the certificate
	Description        string     // Description of the eco-friendly practices recognized
	ExpirationDate     *time.Time // Optional expiration date of the certificate
	CertificationLevel string     // Level or grade of the certification
}

// CircularEconomyProgram represents a registered circular economy program.
type CircularEconomyProgram struct {
	ProgramID     string             // Unique identifier for the program
	Description   string             // Description of the program
	Initiator     string             // Entity or individual that started the program
	StartDate     time.Time          // Start date of the program
	Active        bool               // Whether the program is currently active
	ImpactMetrics map[string]float64 // Metrics for assessing the program's impact (e.g., waste reduced, resources reused)
	Participants  []string           // List of participating entities
}

// EcoFriendlySoftware represents eco-friendly software and its registration details.
type EcoFriendlySoftware struct {
	SoftwareID       string    // Unique ID for the software
	EntityID         string    // Entity responsible for the software
	RegistrationDate time.Time // Date of software registration
	Description      string    // Description of eco-friendly aspects (e.g., energy efficiency, resource optimization)
	Version          string    // Software version
	Verified         bool      // Whether the software has been verified for eco-friendliness
	Certifications   []string  // List of certifications associated with the software
}

// EnergyEfficiencyRating represents an energy efficiency rating for hardware or software.
type EnergyEfficiencyRating struct {
	EntityID   string    // Unique ID of the rated entity
	Rating     float64   // Energy efficiency rating (e.g., out of 5 or 10)
	RatedBy    string    // Entity or individual providing the rating
	RatingDate time.Time // Date of the rating
	Comments   string    // Optional comments about the rating
	Verified   bool      // Whether the rating has been independently verified
}

// EcoCertificate represents an eco-friendly certificate.
type EcoCertificate struct {
	CertificateID string    // Unique ID for the certificate
	Recipient     string    // Entity receiving the certificate
	Type          string    // Type of certificate (e.g., green hardware, energy efficient)
	IssuedAt      time.Time // Time when the certificate was issued
	RenewedAt     time.Time // Time when the certificate was last renewed
	RevokedAt     time.Time // Time when the certificate was revoked (if applicable)
}

// EnergyUsage represents an energy usage record.
type EnergyUsage struct {
	UsageID    string    // Unique ID for the energy usage record
	EntityID   string    // Entity being tracked for energy usage
	EnergyUsed float64   // Amount of energy used (in kWh)
	RecordedAt time.Time // Time when the energy usage was recorded
}

// ConservationInitiative represents a conservation program or initiative.
type ConservationInitiative struct {
	InitiativeID string    // Unique ID for the initiative
	Name         string    // Name of the conservation initiative
	Description  string    // Description of the initiative
	LaunchedAt   time.Time // Time when the initiative was launched
}

// OptimizationRecord represents an energy consumption optimization.
type OptimizationRecord struct {
	OptimizationID string    // Unique ID for the optimization event
	Details        string    // Details of the optimization
	OptimizedAt    time.Time // Time when the optimization occurred
}

// CarbonCreditSystem manages carbon credit issuance, trading, and retiring
type CarbonCreditSystem struct {
	Tokens            map[string]*tokenledgers.SYN1800Token // Carbon credits in the system
	Ledger            *Ledger                               // Ledger for logging all transactions
	EncryptionService *Encryption                           // Encryption for securing sensitive data
	mu                sync.Mutex                            // Mutex for concurrency control
}

// OffsetRequest represents an entity requesting to offset their carbon emissions
type OffsetRequest struct {
	RequestID     string    // Unique identifier for the offset request
	Requester     string    // Wallet address of the requester
	OffsetAmount  float64   // Amount of carbon emissions to offset (in tons of CO2)
	RequestedTime time.Time // Timestamp when the request was made
	IsFulfilled   bool      // Whether the offset request has been fulfilled
	FulfilledTime time.Time // Timestamp when the request was fulfilled
}

// CarbonOffsetMatch represents a match between an offset request and available carbon credits
type CarbonOffsetMatch struct {
	MatchID       string    // Unique identifier for the offset match
	RequestID     string    // Offset request being fulfilled
	CreditID      string    // Carbon credit used to fulfill the request
	MatchedAmount float64   // Amount of carbon credits matched
	MatchedTime   time.Time // Timestamp when the match was made
}

// CarbonOffsetMatchingSystem manages matching offset requests with available carbon credits
type CarbonOffsetMatchingSystem struct {
	Requests          map[string]*OffsetRequest             // Active offset requests
	Matches           map[string]*CarbonOffsetMatch         // Successful offset matches
	AvailableCredits  map[string]*tokenledgers.SYN1800Token // Available carbon credits in the system
	Ledger            *Ledger                               // Ledger for logging all matching transactions
	EncryptionService *Encryption                           // Encryption for securing sensitive data
	mu                sync.Mutex                            // Mutex for concurrency control
}

// EcoFriendlyNodeCertificate represents a certificate for eco-friendly nodes
type EcoFriendlyNodeCertificate struct {
	CertificateID string    // Unique identifier for the certificate
	NodeID        string    // Node ID being certified
	Owner         string    // Owner of the node (wallet address)
	Issuer        string    // Issuing authority or organization
	IssueDate     time.Time // Date the certificate was issued
	ExpiryDate    time.Time // Expiry date of the certificate
	IsRevoked     bool      // Whether the certificate has been revoked
	RevokedDate   time.Time // Date the certificate was revoked (if applicable)
}

// EcoFriendlyNodeCertificationSystem manages the certification process for eco-friendly nodes
type EcoFriendlyNodeCertificationSystem struct {
	Certificates      map[string]*EcoFriendlyNodeCertificate // Map of node certificates
	Ledger            *Ledger                                // Ledger for logging all transactions
	EncryptionService *Encryption                            // Encryption for securing sensitive data
	mu                sync.Mutex                             // Mutex for concurrency control
}

// EfficiencyRating represents an energy efficiency rating for a node
type EfficiencyRating struct {
	RatingID    string    // Unique identifier for the rating
	NodeID      string    // Node being rated
	Owner       string    // Owner of the node (wallet address)
	Issuer      string    // Issuing authority or organization
	Rating      float64   // Energy efficiency rating (on a scale, e.g., 1-10)
	IssueDate   time.Time // Date the rating was issued
	ExpiryDate  time.Time // Expiry date of the rating
	IsRevoked   bool      // Whether the rating has been revoked
	RevokedDate time.Time // Date the rating was revoked (if applicable)
	EnergyUsage float64   // Energy consumption (kWh)
}

// EnergyEfficiencyRatingSystem manages the issuance and tracking of energy efficiency ratings
type EnergyEfficiencyRatingSystem struct {
	Ratings           map[string]*EfficiencyRating // Map of node energy efficiency ratings
	Ledger            *Ledger                      // Ledger for logging all rating transactions
	EncryptionService *Encryption                  // Encryption for securing sensitive data
	mu                sync.Mutex                   // Mutex for concurrency control
}

// EnergyUsageRecord represents the energy consumption data for a node
type EnergyUsageRecord struct {
	RecordID    string    // Unique identifier for the energy usage record
	NodeID      string    // Node being monitored
	Owner       string    // Owner of the node (wallet address)
	EnergyUsage float64   // Energy consumption (kWh) during the monitoring period
	PeriodStart time.Time // Start of the monitoring period
	PeriodEnd   time.Time // End of the monitoring period
	LoggedTime  time.Time // Timestamp when the record was logged
}

// EnergyUsageMonitoringSystem manages the tracking and logging of energy consumption data
type EnergyUsageMonitoringSystem struct {
	UsageRecords      map[string]*EnergyUsageRecord // Map of energy usage records
	Ledger            *Ledger                       // Ledger for logging all usage records
	EncryptionService *Encryption                   // Encryption for securing sensitive data
	mu                sync.Mutex                    // Mutex for concurrency control
}

// GreenHardware represents eco-friendly hardware registered in the system
type GreenHardware struct {
	HardwareID     string    // Unique identifier for the hardware
	Manufacturer   string    // Manufacturer of the hardware
	Model          string    // Model of the hardware
	EnergyRating   float64   // Energy efficiency rating of the hardware
	RegisteredDate time.Time // Date when the hardware was registered
}

// GreenTechnologySystem manages green technology registration, efficiency calculations, and eco-friendly initiatives
type GreenTechnologySystem struct {
	HardwareInventory map[string]*GreenHardware // Inventory of registered green hardware
	SoftwareInventory map[string]string         // Software registered in the system (softwareID -> description)
	Programs          map[string]string         // Circular economy programs (programID -> description)
	NodeCertificates  map[string]string         // Eco-friendly certificates for nodes (nodeID -> certificate details)
	Conservation      []string                  // List of conservation initiatives
	Ledger            *Ledger                   // Ledger for logging transactions
	EncryptionService *Encryption               // Encryption for securing data
	mu                sync.Mutex                // Mutex for concurrency control
}

// RenewableEnergySource represents a renewable energy source integrated into the system
type RenewableEnergySource struct {
	SourceID        string    // Unique identifier for the energy source
	SourceType      string    // Type of renewable energy (e.g., solar, wind, hydro)
	EnergyProduced  float64   // Amount of energy produced (kWh)
	IntegrationDate time.Time // Date the source was integrated into the system
}

// RenewableEnergyIntegrationSystem manages the integration of renewable energy sources into the network
type RenewableEnergyIntegrationSystem struct {
	EnergySources     map[string]*RenewableEnergySource // Map of renewable energy sources
	TotalEnergy       float64                           // Total renewable energy contributed to the network
	Ledger            *Ledger                           // Ledger for logging energy integration
	EncryptionService *Encryption                       // Encryption for securing sensitive data
	mu                sync.Mutex                        // Mutex for concurrency control
}

// ************** Testnet Structs **************

// ContractDeployment represents the deployment details of a smart contract in the testnet.
type TestnetContractDeployment struct {
	ContractID   string    // Unique ID for the deployed contract
	Deployer     string    // Address of the account deploying the contract
	DeployedAt   time.Time // Timestamp of the contract deployment
	ContractCode string    // The contract code deployed
}

type TestnetContractExecution struct {
	ContractID string    // ID of the contract being executed
	Executor   string    // Address or ID of the executor
	InputData  string    // Input data provided to the contract
	ExecutedAt time.Time // Timestamp of execution
	Status     string    // Execution status (e.g., "success", "failed")
}

// TestnetFaucetClaim represents a faucet claim in the testnet.
type TestnetFaucetClaim struct {
	ClaimID   string    // Unique identifier for the claim
	Claimer   string    // Identifier of the claimer
	ClaimedAt time.Time // Timestamp when the claim was made
	Amount    uint64    // Amount claimed
}

// TestnetMetrics captures performance and usage data for the testnet.
type TestnetMetrics struct {
	TotalTransactions       int       // Total number of transactions processed
	ActiveParticipants      int       // Number of active participants in the testnet
	ContractDeploymentCount int       // Number of smart contracts deployed
	BlockHeight             int       // Current block height
	SubBlockHeight          int       // Current sub-block height
	AverageTransactionTime  float64   // Average time for transactions to be processed
	NetworkThroughput       float64   // Transactions per second (TPS) metric
	LastUpdated             time.Time // Timestamp of the last metrics update
}

// TestnetEvent represents an event occurring within the testnet.
type TestnetEvent struct {
	EventID       string    // Unique identifier for the event
	EventType     string    // Type of event (e.g., Transaction, Block, ContractDeployment)
	Description   string    // Description of the event
	Timestamp     time.Time // Time of the event
	RelatedEntity string    // Associated entity (e.g., contract ID, participant ID)
	Status        string    // Status of the event (e.g., Completed, Pending, Failed)
}

// TestnetConfiguration defines the settings and parameters for the testnet.
type TestnetConfiguration struct {
	NetworkName        string        // Name of the testnet
	GenesisBlockHash   string        // Hash of the genesis block
	ConsensusMechanism string        // Consensus mechanism used (e.g., PoW, PoS)
	TransactionLimit   int           // Maximum number of transactions per block
	BlockInterval      time.Duration // Time interval for block generation
	MaxParticipants    int           // Maximum allowed participants in the testnet
	AllowedTokenTypes  []string      // List of allowed token standards
	LastUpdated        time.Time     // Timestamp of the last configuration update
}

// TestnetLogEntry represents a log entry for activities in the testnet.
type TestnetLogEntry struct {
	LogID        string    // Unique identifier for the log entry
	Timestamp    time.Time // Time the log entry was created
	LogType      string    // Type of log (e.g., Transaction, Block, Error)
	Description  string    // Description of the log entry
	AssociatedID string    // ID of the associated entity (e.g., transaction ID)
	Severity     string    // Severity of the log (e.g., Info, Warning, Error)
}

// ParticipantRecord captures details of a participant in the testnet.
type ParticipantRecord struct {
	ParticipantID    string    // Unique identifier for the participant
	PublicKey        string    // Public key of the participant
	Role             string    // Role of the participant (e.g., Validator, Developer, User)
	JoinDate         time.Time // Date the participant joined the testnet
	TransactionCount int       // Number of transactions performed by the participant
	Status           string    // Current status (e.g., Active, Inactive, Banned)
	ReputationScore  float64   // Reputation score of the participant
}

// TokenDeployment represents the details of a token deployed in the testnet.
type TokenDeployment struct {
	TokenID       string    // Unique ID of the token
	Deployer      string    // Address of the account deploying the token
	TokenSymbol   string    // Symbol of the token (e.g., TST for Test Token)
	InitialSupply uint64    // The initial supply of the token
	DeployedAt    time.Time // Timestamp of the token deployment
}

// FaucetClaimRecord represents the details of a faucet claim in the testnet.
type FaucetClaimRecord struct {
	ClaimID   string    // Unique ID of the claim
	Claimer   string    // Address of the account that claimed the faucet
	ClaimedAt time.Time // Timestamp of the claim
	Amount    uint64    // Amount of tokens claimed from the faucet
}

// ************** Token Utility Structs **************

// LeaseTerms represents the terms for leasing an asset.
type LeaseTerms struct {
	LeaseID               string        // Unique identifier for the lease agreement
	AssetID               string        // ID of the leased asset
	Lessee                string        // Wallet address or ID of the lessee
	Lessor                string        // Wallet address or ID of the lessor
	LeaseDuration         time.Duration // Duration of the lease
	LeaseStartDate        time.Time     // Start date of the lease
	LeaseEndDate          time.Time     // End date of the lease
	PaymentTerms          string        // Details of the payment terms (e.g., monthly, lump sum)
	PaymentAmount         float64       // Total payment amount for the lease
	LateFee               float64       // Fee for late payments
	TerminationConditions string        // Conditions under which the lease may be terminated early
	AutoRenewal           bool          // Indicates whether the lease auto-renews
	Signatories           []string      // IDs or addresses of the parties who have signed the agreement
}

// LicenseTerms represents the terms for licensing an asset.
type LicenseTerms struct {
	LicenseID             string        // Unique identifier for the license agreement
	AssetID               string        // ID of the licensed asset
	Licensee              string        // Wallet address or ID of the licensee
	Licensor              string        // Wallet address or ID of the licensor
	LicenseDuration       time.Duration // Duration of the license agreement
	LicenseStartDate      time.Time     // Start date of the license
	LicenseEndDate        time.Time     // End date of the license
	UsageRights           string        // Description of the rights granted (e.g., reproduction, distribution)
	RoyaltyPercentage     float64       // Percentage of royalties, if applicable
	LicenseFee            float64       // Total fee for the license
	RenewalOptions        bool          // Whether the license is renewable
	TerminationConditions string        // Conditions under which the license may be terminated early
	Signatories           []string      // IDs or addresses of the parties who have signed the agreement
}

// RentalTerms represents the terms for renting an asset.
type RentalTerms struct {
	RentalID              string        // Unique identifier for the rental agreement
	AssetID               string        // ID of the rented asset
	Renter                string        // Wallet address or ID of the renter
	Owner                 string        // Wallet address or ID of the asset owner
	RentalDuration        time.Duration // Duration of the rental
	RentalStartDate       time.Time     // Start date of the rental
	RentalEndDate         time.Time     // End date of the rental
	RentalFee             float64       // Total rental fee
	DepositAmount         float64       // Deposit required for the rental
	MaintenanceTerms      string        // Terms regarding maintenance responsibility
	LateReturnFee         float64       // Fee for returning the asset late
	TerminationConditions string        // Conditions under which the rental may be terminated early
	Signatories           []string      // IDs or addresses of the parties who have signed the agreement
}

// ProvenanceRecord tracks the history of ownership and transactions for an asset.
type ProvenanceRecord struct {
	RecordID      string    // Unique identifier for the provenance record
	AssetID       string    // ID of the asset being tracked
	PreviousOwner string    // ID or wallet address of the previous owner
	NewOwner      string    // ID or wallet address of the new owner
	TransferDate  time.Time // Date of the ownership transfer
	TransactionID string    // ID of the transaction that facilitated the transfer
	Description   string    // Description of the transfer or event
}

// OwnershipLedger tracks the ownership of assets.
type OwnershipLedger struct {
	AssetID          string             // Unique identifier for the asset
	CurrentOwner     string             // Current owner of the asset
	OwnershipHistory []ProvenanceRecord // List of past ownership changes
}

// TransactionLedger records transactions related to assets.
type TransactionLedger struct {
	TransactionID   string    // Unique identifier for the transaction
	AssetID         string    // ID of the asset involved
	Sender          string    // Sender's wallet address
	Receiver        string    // Receiver's wallet address
	Amount          float64   // Amount transferred (if applicable)
	Timestamp       time.Time // Time when the transaction occurred
	TransactionType string    // Type of transaction (e.g., "Lease", "Sale", "Transfer")
}

// AssetValuationManager handles the valuation of assets.
type AssetValuationManager struct {
	AssetID          string            // Unique identifier for the asset
	CurrentValue     float64           // Current value of the asset
	ValuationHistory []ValuationRecord // Record of past valuations
}

// ValuationRecord represents a historical valuation of an asset.
type ValuationRecord struct {
	ValuationDate   time.Time // Date of the valuation
	ValuationAmount float64   // Value assigned to the asset
	Valuer          string    // ID or wallet address of the individual/entity who conducted the valuation
}

// Notifier handles sending notifications for leases, licenses, and other events.
type Notifier struct {
	NotificationID   string    // Unique identifier for the notification
	Recipient        string    // ID or wallet address of the notification recipient
	NotificationType string    // Type of notification (e.g., "Lease Expiry", "Payment Due")
	Message          string    // The message content of the notification
	SentDate         time.Time // Timestamp when the notification was sent
}

// LeaseManagement manages lease agreements.
type LeaseManagement struct {
	LeaseID      string                // Unique identifier for the lease agreement
	Terms        LeaseTerms            // Lease terms associated with the agreement
	ActiveLeases map[string]LeaseTerms // Map of active lease agreements by asset ID
}

// LicenseManagement manages license agreements.
type LicenseManagement struct {
	LicenseID      string                  // Unique identifier for the license agreement
	Terms          LicenseTerms            // License terms associated with the agreement
	ActiveLicenses map[string]LicenseTerms // Map of active license agreements by asset ID
}

// RentalManagement manages rental agreements.
type RentalManagement struct {
	RentalID      string                 // Unique identifier for the rental agreement
	Terms         RentalTerms            // Rental terms associated with the agreement
	ActiveRentals map[string]RentalTerms // Map of active rental agreements by asset ID
}

// CoOwnershipManagement handles agreements where multiple parties co-own an asset.
type CoOwnershipManagement struct {
	AssetID        string             // ID of the co-owned asset
	CoOwners       []string           // List of wallet addresses of the co-owners
	OwnershipShare map[string]float64 // Map of ownership shares for each co-owner
	DecisionRules  string             // Rules for decision-making among co-owners
}

// SYN10AuditLog represents an audit log entry within the system.
type SYN10AuditLog struct {
	LogID            string    // Unique identifier for the audit log entry
	EventType        string    // Type of event being logged (e.g., "Transaction", "System Change", "Access")
	Timestamp        time.Time // Timestamp of when the event occurred
	ActorID          string    // ID or wallet address of the entity responsible for the event
	ActionDetails    string    // Detailed description of the action or event
	AffectedEntities []string  // List of entities (e.g., accounts, contracts) affected by the event
	Status           string    // Status of the audit log (e.g., "Pending", "Completed", "Failed")
	Hash             string    // Cryptographic hash of the log entry for verification
}

// QuantitativeEasingMechanism represents a mechanism for injecting liquidity into the economy.
type QuantitativeEasingMechanism struct {
	EasingID          string         // Unique identifier for the easing event
	AssetType         string         // Type of assets being purchased (e.g., bonds, equities)
	PurchaseAmount    float64        // Total amount of assets being purchased
	PurchaseDate      time.Time      // Date when the quantitative easing began
	LiquidityInjected float64        // Total liquidity injected into the system
	Duration          time.Duration  // Duration over which easing will occur
	CentralBankNode   string         // Central bank node responsible for executing the easing
	ImpactAssessment  []ImpactRecord // Records tracking the economic impact of easing
}

// ImpactRecord represents the economic impact record of the quantitative easing event.
type ImpactRecord struct {
	AssessmentDate   time.Time // Date of the impact assessment
	InflationRate    float64   // Inflation rate at the time of assessment
	GDPGrowth        float64   // GDP growth at the time of assessment
	UnemploymentRate float64   // Unemployment rate at the time of assessment
}

// MonetaryTighteningMechanism represents a mechanism for reducing liquidity in the economy.
type MonetaryTighteningMechanism struct {
	TighteningID       string         // Unique identifier for the tightening event
	AssetType          string         // Type of assets being sold or withdrawn (e.g., bonds, equities)
	SaleAmount         float64        // Total amount of assets being sold or withdrawn
	SaleDate           time.Time      // Date when the monetary tightening began
	LiquidityWithdrawn float64        // Total liquidity withdrawn from the system
	Duration           time.Duration  // Duration over which tightening will occur
	CentralBankNode    string         // Central bank node responsible for executing the tightening
	ImpactAssessment   []ImpactRecord // Records tracking the economic impact of tightening
}

// AuditLogs is a map of audit log entries keyed by unique log ID.
type AuditLogs map[string]SYN10AuditLog

// EasingMechanism is a reference to the quantitative easing mechanism in the system.
type EasingMechanism struct {
	Mechanism *QuantitativeEasingMechanism
}

// TighteningMechanism is a reference to the monetary tightening mechanism in the system.
type TighteningMechanism struct {
	Mechanism *MonetaryTighteningMechanism
}

// SYN131EventListener represents an event listener for handling SYN131 events within the system.
type SYN131EventListener struct {
	ListenerID    string            // Unique identifier for the event listener
	EventName     string            // Name of the event being listened for (e.g., "TransactionComplete", "ContractExecuted")
	EventHandler  func(interface{}) // Function or handler to execute when the event occurs
	Active        bool              // Indicates whether the listener is currently active
	RegisteredAt  time.Time         // Timestamp when the listener was registered
	LastTriggered time.Time         // Timestamp when the listener was last triggered
	TriggerCount  int               // Number of times the listener has been triggered
	Mutex         sync.Mutex        // Mutex for ensuring thread-safe operations
}

// PeggedAsset represents an asset pegged to an external value or index.
type PeggedAsset struct {
	AssetID         string    // ID of the pegged asset
	ExternalIndex   string    // External index or value the asset is pegged to
	CurrentPegValue float64   // Current value pegged to the external index
	PeggedSince     time.Time // Timestamp when the asset was pegged
}

// TrackedAsset represents a tangible asset that is being tracked on the blockchain.
type TrackedAsset struct {
	AssetID           string           // ID of the tracked asset
	TrackingMethod    string           // Method of tracking (e.g., GPS, RFID, Barcode)
	LastKnownLocation string           // Last known location of the asset
	TrackingHistory   []TrackingRecord // Record of past tracking data
}

// AssetStatus represents the current status of a tangible or intangible asset.
type AssetStatus struct {
	AssetID       string               // ID of the asset
	CurrentStatus string               // Current status (e.g., "Active", "Inactive", "Under Maintenance")
	StatusHistory []StatusChangeRecord // Record of past status changes
}

// TrackingRecord represents a point in the tracking history of an asset.
type TrackingRecord struct {
	Timestamp time.Time // Time when the tracking data was recorded
	Location  string    // Location of the asset at the time of tracking
}

// StatusChangeRecord represents a point in the status history of an asset.
type StatusChangeRecord struct {
	Timestamp time.Time // Time when the status changed
	OldStatus string    // Previous status
	NewStatus string    // New status
}

// LicenseAgreement represents the legal agreement for licensing an asset.
type LicenseAgreement struct {
	AgreementID    string       // Unique identifier for the license agreement
	Licensee       string       // Wallet address or ID of the licensee
	Licensor       string       // Wallet address or ID of the licensor
	LicenseTerms   LicenseTerms // Terms of the license agreement
	EffectiveDate  time.Time    // Date when the license becomes effective
	ExpirationDate time.Time    // Date when the license expires
}

// RoyaltyInfo represents information about royalties for an asset.
type RoyaltyInfo struct {
	RoyaltyID      string          // Unique identifier for the royalty record
	AssetID        string          // ID of the asset generating royalties
	Percentage     float64         // Royalty percentage
	Payee          string          // Wallet address of the royalty recipient
	LastPayment    time.Time       // Timestamp of the last royalty payment
	PaymentHistory []PaymentRecord // Record of past royalty payments
}

// ModeChangeLog tracks changes in the operating mode of an asset or system.
type ModeChangeLog struct {
	ChangeID   string    // Unique identifier for the mode change
	AssetID    string    // ID of the asset whose mode changed
	OldMode    string    // Previous mode
	NewMode    string    // New mode
	ChangeDate time.Time // Date when the mode was changed
}

// ************** Transactions Structs **************

// TransactionRecord keeps track of all transactions in the ledger.
type TransactionRecord struct {
	From         string    // Sender's address (for financial transactions)
	To           string    // Recipient's address (for financial transactions)
	Amount       float64   // Amount transferred (if applicable)
	Fee          float64   // Transaction fee (if applicable)
	Hash         string    // Transaction hash (unique ID)
	Status       string    // Status of the transaction (e.g., "pending", "confirmed")
	BlockIndex   int       // Block in which the transaction was confirmed
	Timestamp    time.Time // Timestamp when the transaction was created or confirmed
	BlockHeight  int       // Height of the block containing the transaction
	ValidatorID  string    // ID of the validator who processed the transaction
	ID           string    // ID of the transaction or shard-related activity
	Action       string    // Action performed, e.g., "ShardCreated", "ShardUpdated", etc.
	Delegator    string    // Delegator involved in shard delegation (if applicable)
	NodeID       string    // Node ID for shard reallocation (if applicable)
	Orchestrator string    // Orchestrator for orchestrated transactions (if applicable)
	Details      string    // Additional details (if applicable)
}

// BlockMetric tracks the metrics for a block.
type BlockMetric struct {
	BlockID      string    // Unique ID of the block
	Timestamp    time.Time // Time the metrics were recorded
	Metrics      string    // Performance metrics or other relevant data
	BlockSize    int64     // Size of the block in bytes
	Transactions int       // Number of transactions in the block
	ValidatorID  string    // ID of the validator that validated the block
}

// SubBlockMetric tracks the metrics for a sub-block.
type SubBlockMetric struct {
	SubBlockID   string    // Unique ID of the sub-block
	Timestamp    time.Time // Time the metrics were recorded
	Metrics      string    // Performance metrics or other relevant data
	SubBlockSize int64     // Size of the sub-block in bytes
	Transactions int       // Number of transactions in the sub-block
	ParentBlock  string    // ID of the parent block
}

// EncryptedMessage represents an encrypted payload to be sent between nodes.
type EncryptedMessage struct {
	CipherText []byte    // The encrypted message content
	Hash       string    // Hash of the original message for integrity check
	CreatedAt  time.Time // When the message was encrypted (Renamed from Timestamp to CreatedAt)
}

// TransactionPool manages the unconfirmed transactions waiting to be validated.
type TransactionPool struct {
	transactions      map[string]*Transaction   // Map of transaction ID to Transaction
	pendingSubBlocks  map[string][]*Transaction // Map of sub-block IDs to transactions
	mu                sync.Mutex
	maxPoolSize       int
	ledger            *Ledger
	encryptionService *Encryption
}

// EscrowTransaction holds the details of the escrow agreement.
type EscrowTransaction struct {
	EscrowID     string
	SenderID     string
	ReceiverID   string
	Amount       float64
	Status       EscrowStatus
	CreationTime time.Time
	ReleaseTime  time.Time
	Condition    string // Optional: Condition to release funds (e.g., service completion)
}

// EscrowManager manages escrow transactions and ensures the proper release or cancellation of funds.
type EscrowManager struct {
	ledgerInstance *Ledger
	mutex          sync.Mutex
	escrows        map[string]*EscrowTransaction
}

// SmartLegalContractManager manages multiple smart legal contracts.
type SmartLegalContractManager struct {
	Contracts      map[string]*SmartLegalContract // All deployed legal contracts
	LedgerInstance *Ledger                        // Ledger for recording contract deployments and executions
	mutex          sync.Mutex                     // Mutex for safe concurrency
}

// TransactionCancellationManager handles transaction cancellation requests and processing.
type TransactionCancellationManager struct {
	mutex               sync.Mutex
	Consensus           *SynnergyConsensus
	TimeoutPeriod       time.Duration // Time allowed for requesting a cancellation
	ResponseTimeout     time.Duration // Time allowed for node response
	NotificationService Notification  // Notifies involved parties
	Ledger              *Ledger
	Encryption          *Encryption
}

// TransactionReversalManager manages transaction reversals.
type TransactionReversalManager struct {
	mutex               sync.Mutex
	Consensus           *SynnergyConsensus
	Ledger              *Ledger
	ReversalTimeLimit   time.Duration // Time allowed for requesting reversal (within 28 days)
	NotificationService Notification  // Notifies involved parties
	Encryption          *Encryption
}

// ChannelState represents the state of a payment or communication channel between two or more parties.
type ChannelState struct {
	ChannelID    string             // Unique identifier for the channel
	Participants []string           // List of participants in the channel (wallet addresses or node IDs)
	Balance      map[string]float64 // Balance of each participant in the channel
	ChannelType  string             // Type of channel (e.g., "Payment", "State", "DataStream")
	Expiry       time.Time          // Expiration time of the channel (if applicable)
	Status       string             // Current status of the channel (e.g., "Open", "Closed", "Suspended")
	Mutex        sync.Mutex         // Mutex for thread-safe channel operations
}

// Timelock represents a mechanism to lock transactions, contracts, or assets until a specific time is reached.
type Timelock struct {
	LockID       string    // Unique identifier for the timelock
	LockedUntil  time.Time // Timestamp when the lock expires and assets are released
	LockedValue  float64   // Amount of tokens or assets locked
	LockedEntity string    // The entity (contract, transaction) under the timelock
	Status       string    // Current status of the timelock (e.g., "Active", "Expired")
	IsRevocable  bool      // Whether the timelock can be revoked before expiry
}

// GeoProcessor handles the processing and validation of geospatial data in the blockchain.
type GeoProcessor struct {
	ProcessorID       string             // Unique identifier for the geospatial processor
	GeospatialData    map[string]float64 // Map of geospatial data (coordinates, regions, etc.)
	LastProcessed     time.Time          // Timestamp of the last geospatial data processed
	Accuracy          float64            // Accuracy of the geospatial data in meters
	EncryptionService *Encryption        // Service to secure the geospatial data during transmission and storage
}

// HolographicData represents data that has been holographically encoded for storage across a distributed system.
type HolographicData struct {
	DataID           string   // Unique identifier for the holographic data
	EncodedData      []byte   // The holographically encoded data
	DataFragments    [][]byte // Fragments of the encoded data stored across different nodes
	StorageNodes     []string // List of node IDs where fragments are stored
	RedundancyFactor int      // Redundancy factor for data replication across nodes
}

// ChainAdapter facilitates interaction between different blockchains by enabling cross-chain transactions and data exchanges.
type ChainAdapter struct {
	AdapterID         string      // Unique identifier for the chain adapter
	SourceChain       string      // Source blockchain for the interaction
	DestinationChain  string      // Destination blockchain for the interaction
	SupportedAssets   []string    // List of assets (tokens, NFTs) supported by the adapter
	TransactionLogs   []string    // Log of transactions processed by the adapter
	EncryptionService *Encryption // Encryption service for securing cross-chain interactions
}

// APIGateway manages API requests and serves as an interface between external systems and the blockchain.
type APIGateway struct {
	GatewayID      string                 // Unique identifier for the API gateway
	ActiveRequests map[string]*APIRequest // Map of active API requests being handled
	MaxConcurrent  int                    // Maximum number of concurrent requests allowed
	RateLimit      int                    // Rate limit for handling API requests (requests per second)
	Timeout        time.Duration          // Timeout for each API request
	LoadBalancer   *LoadBalancer          // Load balancer for distributing API requests across nodes
}

// Oracle represents an external data provider that feeds information into the blockchain for use in smart contracts.
type Oracle struct {
	OracleID          string             // Unique identifier for the oracle
	DataFeeds         map[string]float64 // Data feeds provided by the oracle (e.g., prices, weather data)
	LastUpdated       time.Time          // Timestamp of the last update to the data feeds
	ReliabilityScore  float64            // Score representing the reliability of the oracle (0.0 to 1.0)
	DataSource        string             // External source of the data (e.g., API, IoT devices)
	Frequency         time.Duration      // Frequency of data updates from the oracle
	EncryptionService *Encryption        // Encryption service for securing the data feed during transmission
}

// TransactionMetricsManager manages the collection of transaction metrics in the blockchain.
type TransactionMetricsManager struct {
	totalTransactions     int
	totalSubBlocks        int
	totalBlocks           int
	totalGasConsumed      int
	totalFeesCollected    float64
	transactionThroughput float64 // transactions per second (TPS)
	gasEfficiency         float64 // ratio of gas used vs gas limit

	metricsLock sync.Mutex
	ledger      *Ledger
}

// ArchiveManager handles archiving and managing old blockchain data.
type ArchiveManager struct {
	ledgerInstance *Ledger           // Reference to the ledger instance
	archivePath    string            // Path to store archived data
	mutex          sync.Mutex        // Mutex for thread-safe operations
	memoryStorage  map[string][]byte // Store encrypted archives in memory
}

// PrivateTransaction defines the structure of a private transaction.
type PrivateTransaction struct {
	TransactionID   string            // Unique identifier for the transaction
	Sender          string            // Sender of the transaction
	Receiver        string            // Receiver of the transaction
	Amount          float64           // Transaction amount
	TokenType       string            // Token type (optional, defaults to "SYNN")
	TokenID         string            // Token ID (optional)
	IsPrivate       bool              // Flag indicating if this transaction is private
	EncryptedData   string            // Encrypted transaction details
	AuthorizedNodes map[string]string // Nodes allowed to view private details
	Fee             float64           // Fee for converting to private transaction
}

// PrivateTransactionManager manages the creation and conversion of private transactions.
type PrivateTransactionManager struct {
	mutex              sync.Mutex                     // For thread-safe operations
	Transactions       map[string]*PrivateTransaction // List of all private transactions
	Ledger             *Ledger                        // Ledger reference for transaction logging
	Consensus          *SynnergyConsensus             // Consensus engine for validation
	Encryption         *Encryption                    // Encryption service for securing transaction details
	AuthorityNodeTypes []string                       // List of authority node types allowed to manage private transactions
	TransactionPool    *TransactionPool               // Pool for holding unconfirmed transactions
}

// CancellationRequest represents a request to cancel a transaction.
type CancellationRequest struct {
	ID                string          // Unique request ID
	TransactionID     string          // ID of the transaction to cancel
	UserID            string          // ID of the user requesting cancellation
	Timestamp         time.Time       // Timestamp of the cancellation request
	Status            string          // Current status (e.g., "Pending", "Approved", "Rejected")
	Reason            string          // Reason for the cancellation request
	DocumentEvidence  string          // Any evidence provided for the cancellation
	ContactDetails    string          // Contact details of the user requesting cancellation
	RequiredApprovals int             // Number of required approvals for cancellation
	ApprovalNodes     map[string]bool // Nodes that have approved
	RejectionNodes    map[string]bool // Nodes that have rejected
}

// FeeManager manages the calculation and enforcement of transaction fee ceilings and floors.
type FeeManager struct {
	ledgerInstance      *Ledger    // Reference to the ledger instance
	mutex               sync.Mutex // Mutex for thread-safe operations
	BaseFee             float64    // The base fee for transactions
	NetworkLoad         float64    // Current network load percentage (0 to 1)
	PendingTransactions int        // Current number of pending transactions
}

// TransactionReceiptManager manages the creation, storage, and validation of transaction receipts.
type TransactionReceiptManager struct {
	receipts          map[string]*TransactionReceipt // Map of transaction IDs to their receipts
	encryptionService *Encryption                    // Encryption service for receipt integrity
}

// EscrowStatus represents the status of an escrow transaction.
type EscrowStatus string

const (
	EscrowStatusPending   EscrowStatus = "Pending"
	EscrowStatusReleased  EscrowStatus = "Released"
	EscrowStatusCancelled EscrowStatus = "Cancelled"
)

// IsFinalized checks if the escrow status is finalized (either Released or Cancelled)
func (status EscrowStatus) IsFinalized() bool {
	return status == EscrowStatusReleased || status == EscrowStatusCancelled
}

// Transaction represents a blockchain transaction with all the required fields.
type Transaction struct {
	TransactionID     string    // Unique identifier for the transaction
	FromAddress       string    // Sender's address
	ToAddress         string    // Receiver's address
	Amount            float64   // Amount being transferred
	Fee               float64   // Transaction fee
	TokenStandard     string    // Token standard (e.g., ERC20, Syn700)
	TokenID           string    // Unique Token ID for token transactions
	Timestamp         time.Time // Timestamp when the transaction was created
	SubBlockID        string    // Associated sub-block ID
	BlockID           string    // Associated block ID
	ValidatorID       string    // Validator who validated this transaction
	Signature         string    // Transaction signature
	Status            string    // Transaction status (e.g., pending, confirmed, failed)
	EncryptedData     string    // Encrypted transaction data
	DecryptedData     string    // Decrypted transaction data (used internally)
	ExecutionResult   string    // Result after executing the transaction (e.g., success, failure reason)
	FrozenAmount      float64   // Amount that is frozen in the transaction (if applicable)
	RefundAmount      float64   // Amount refunded in case of a reversal or error
	ReversalRequested bool      // Whether a reversal has been requested (Add this field)

}

// TransactionMetric stores performance and analytics data for transactions.
type TransactionMetric struct {
	TransactionID   string    // Unique ID of the transaction
	ExecutionTime   float64   // Time taken to execute the transaction (in milliseconds)
	GasUsed         float64   // Amount of gas consumed by the transaction
	GasPrice        float64   // Gas price for the transaction
	Fee             float64   // Total fee paid for the transaction
	ValidationTime  float64   // Time taken to validate the transaction (in milliseconds)
	RejectionReason string    // Reason for rejection (if applicable)
	ConsensusTime   float64   // Time taken to achieve consensus for the transaction (in milliseconds)
	Timestamp       time.Time // Timestamp of the transaction
	Success         bool      // Indicates if the transaction was successful
	NodeID          string    // ID of the node that processed the transaction
	RetryCount      int       // Number of retries for the transaction
}

// EscrowLog represents logs for escrow transactions.
type EscrowLog struct {
	EscrowID           string                 // Unique identifier for the escrow
	TransactionID      string                 // ID of the associated transaction
	Initiator          string                 // Address or ID of the initiator
	Beneficiary        string                 // Address or ID of the beneficiary
	EscrowAmount       float64                // Amount locked in escrow
	Status             string                 // Status of the escrow (e.g., "Pending", "Released", "Canceled")
	Conditions         map[string]interface{} // Conditions required to release the escrow
	Timestamp          time.Time              // Timestamp of the escrow creation
	ReleaseTimestamp   *time.Time             // Timestamp when escrow was released (if applicable)
	CancellationReason string                 // Reason for cancellation (if applicable)
	Metadata           map[string]interface{} // Additional metadata about the escrow
}

// Distribution percentages
const (
	DevelopmentPoolPercentage              = 0.05
	CharityPoolPercentage                  = 0.10
	LoanPoolPercentage                     = 0.05
	PassiveIncomePoolPercentage            = 0.05
	ValidatorMinerRewardPoolPercentage     = 0.70
	AuthorityNodeHostsRewardPoolPercentage = 0.05
)

// TransactionDistributionManager manages the distribution of transaction fees and rewards.
type TransactionDistributionManager struct {
	ledgerInstance *Ledger
	mutex          sync.Mutex
}

// TransactionManager handles transaction creation, validation, encryption, and ledger integration.
type TransactionManager struct {
	Ledger     *Ledger            // Reference to the blockchain ledger
	Consensus  *SynnergyConsensus // Consensus engine for Synnergy Consensus
	Encryption *Encryption        // Encryption service
}

// TransactionReceipt represents a receipt for a processed transaction.
type TransactionReceipt struct {
	TransactionID      string    // Unique ID of the transaction
	BlockID            string    // ID of the block the transaction is part of
	SubBlockID         string    // ID of the sub-block the transaction is part of
	Status             string    // Transaction status: "SUCCESS", "FAILURE"
	Timestamp          time.Time // Timestamp when the transaction was processed
	GasUsed            uint64    // Amount of gas used by the transaction
	TransactionOutput  string    // Output data from the transaction execution
	ValidatorSignature string    // Validator signature confirming the transaction inclusion
	EncryptionHash     string    // Hash of the encrypted transaction data
}

// ReceiptManager manages the creation and handling of transaction receipts.
type ReceiptManager struct {
	receipts          map[string]*TransactionReceipt // Stores transaction receipts by transaction ID
	encryptionService Encryption                     // Encryption service for handling encryption
}

// TransactionScheduler manages the scheduling of blockchain transactions.
type TransactionScheduler struct {
	mutex           sync.Mutex                      // For thread safety
	ScheduleList    map[string]*TransactionSchedule // Scheduled transactions list
	Ledger          *Ledger                         // Reference to the ledger for transaction logging
	ContractManager *SmartContractManager           // Reference to the SmartContractManager
}

// TransactionSchedule defines a scheduled transaction with its conditions.
type TransactionSchedule struct {
	TransactionID string    // Unique ID of the transaction
	ScheduledTime time.Time // Time when the transaction is scheduled to be executed
	BlockHeight   uint64    // Block height at which the transaction will be executed (optional)
	Condition     string    // Condition to trigger the transaction (e.g., contract state change)
	Recurring     bool      // Whether this transaction is recurring
	NextExecution time.Time // Next execution time for recurring transactions
	EncryptedTx   string    // Encrypted transaction data
	ValidatorID   string    // Validator responsible for the transaction execution
	Executed      bool      // Whether the transaction has already been executed
}

// Condition defines the structure for the transaction condition
type Condition struct {
	ContractID    string      `json:"contractID"`
	ConditionKey  string      `json:"conditionKey"`
	ExpectedValue interface{} `json:"expectedValue"`
}

// TransactionSearchService handles searching transactions based on different criteria.
type TransactionSearchService struct {
	Ledger       *Ledger                   // Reference to the ledger to search transactions
	Cache        map[string]*Transaction   // Cache to store frequently accessed transactions
	Index        map[string][]*Transaction // Index to improve search performance
	cacheMutex   sync.RWMutex              // Mutex for thread-safe cache access
	mutex        sync.RWMutex              // Mutex for safe read/write operations
	cacheEnabled bool                      // Toggle cache usage
}

// TransactionSearchCriteria defines the criteria for searching transactions.
type TransactionSearchCriteria struct {
	TransactionID string    // Filter by Transaction ID
	SenderID      string    // Filter by Sender ID
	RecipientID   string    // Filter by Recipient ID
	DateFrom      time.Time // Filter by start date (inclusive)
	DateTo        time.Time // Filter by end date (inclusive)
	MinAmount     float64   // Filter by minimum transaction amount
	MaxAmount     float64   // Filter by maximum transaction amount
	Status        string    // Filter by transaction status (e.g., "Pending", "Confirmed")
}

// ************** Virtual Machine Structs **************

// VirtualMachine represents the core VM that handles contract execution, bytecode generation, and sandbox management.
type VirtualMachine struct {
	LedgerInstance       *Ledger // Ledger to log transactions and contract execution results
	Contracts            map[string]*SmartContract
	BytecodeInterpreter  *BytecodeInterpreter  // Interprets and executes bytecode
	BytecodeGenerator    *BytecodeGenerator    // Generates bytecode from smart contracts
	CodeQualityAssurance *CodeQualityAssurance // Handles code quality validation
	CompilationDebugger  *CompilationDebugger  // Compiles and debugs bytecode
	SandboxManager       *SandboxManager       // Handles isolated execution environments
	SyntaxChecker        *SyntaxChecker        // Checks contract syntax before execution
	GasManager           *GasManager           // Handles gas calculation for transactions
	SubBlockManager      *SubBlockManager      // Manages sub-block validation and aggregation
	SolidityCompiler     *SolidityCompiler     // Compiler for Solidity contracts
	RustCompiler         *RustCompiler         // Compiler for Rust contracts
	YulCompiler          *YulCompiler          // Compiler for Yul contracts
	JavascriptCompiler   *JavaScriptCompiler   // Compiler for JavaScript contracts
	GolangCompiler       *GoContractCompiler   // Compiler for Go contracts
	SoliditySupport      *SoliditySupport      // Solidity support instance for contract execution
	RustSupport          *RustSupport          // Rust support instance for contract execution
	YulSupport           *YulSupport           // Yul support instance for contract execution
	GolangSupport        *GoSupport            // Go support instance for contract execution
	JavascriptSupport    *JavaScriptSupport    // JavaScript support instance for contract execution
	mutex                sync.Mutex            // Mutex for ensuring thread-safety
}

// Bytecode represents the bytecode deployed for a contract.
type Bytecode struct {
	ContractID string
	Code       []byte
	DeployedAt time.Time
}

// ContractState represents the state of a smart contract.
type ContractState struct {
	ContractID string
	StateData  map[string]interface{}
	UpdatedAt  time.Time
}

// LogEntry represents a log of an operation or event.
type LogEntry struct {
	EntryID     string
	LogType     string
	Description string
	Timestamp   time.Time
}

// IsolationManager handles system isolation to ensure secure and segmented operations.
type IsolationManager struct {
	IsolatedProcesses  map[string]ProcessDetails // Processes currently running in isolation.
	ResourceAllocation map[string]int            // Resources allocated to isolated environments.
	Logs               []IsolationLog            // Logs of isolation events.
}

// ProcessDetails captures details about an isolated process.
type ProcessDetails struct {
	ProcessID      string    // Unique identifier for the process.
	IsolationLevel string    // Level of isolation (e.g., sandbox, VM, container).
	ResourceUsage  int       // Amount of resources allocated to the process.
	Status         string    // Current status of the process (e.g., running, paused, terminated).
	CreatedAt      time.Time // Timestamp when the process was isolated.
}

// IsolationLog captures events related to isolation management.
type IsolationLog struct {
	LogID     string    // Unique identifier for the log entry.
	ProcessID string    // Associated process ID.
	Action    string    // Action performed (e.g., isolated, resumed, terminated).
	Timestamp time.Time // Timestamp of the action.
	Message   string    // Additional information or message.
}

// ExecutionManager handles task execution and process lifecycle management.
type ExecutionManager struct {
	ActiveExecutions map[string]ExecutionDetails // Details of currently active executions.
	ExecutionQueue   []string                    // Queue of pending executions.
	Logs             []ExecutionLog              // Logs of execution events.
}

// ExecutionDetails captures details about an execution task or process.
type ExecutionDetails struct {
	ExecutionID  string    // Unique identifier for the execution.
	TaskID       string    // Associated task or workflow ID.
	Status       string    // Current status (e.g., in-progress, completed, failed).
	AssignedNode string    // Node or resource assigned to handle the execution.
	StartTime    time.Time // Timestamp when the execution started.
	EndTime      time.Time // Timestamp when the execution ended (if applicable).
}

// ExecutionLog captures events during execution management.
type ExecutionLog struct {
	LogID       string    // Unique identifier for the log entry.
	ExecutionID string    // Associated execution ID.
	Action      string    // Action performed (e.g., started, paused, resumed).
	Timestamp   time.Time // Timestamp of the action.
	Message     string    // Additional information or message.
}

// ProcessManager oversees the lifecycle and states of system processes.
type ProcessManager struct {
	ActiveProcesses map[string]ProcessState // Map of currently active processes and their states.
	ProcessQueue    []string                // Queue of processes waiting for execution.
	DependencyGraph map[string][]string     // Dependencies between processes.
	Logs            []ProcessLog            // Logs of process-related events.
}

// ProcessLog captures events during process management.
type ProcessLog struct {
	LogID     string    // Unique identifier for the log entry.
	ProcessID string    // Associated process ID.
	Action    string    // Action performed (e.g., started, terminated, updated).
	Timestamp time.Time // Timestamp of the action.
	Message   string    // Additional information or message.
}

// ContextManager manages execution contexts and their associated data.
type ContextManager struct {
	ActiveContexts map[string]ExecutionContext // Details of currently active contexts.
	ContextQueue   []string                    // Queue of pending context initializations.
	Logs           []ContextLog                // Logs of context-related events.
}

// ContextLog captures events related to execution context management.
type ContextLog struct {
	LogID     string    // Unique identifier for the log entry.
	ContextID string    // Associated context ID.
	Action    string    // Action performed (e.g., created, updated, terminated).
	Timestamp time.Time // Timestamp of the action.
	Message   string    // Additional information or message.
}

// Constants for connection events
const (
	EventTypeConnected    = "CONNECTED"
	EventTypeDisconnected = "DISCONNECTED"
)

// InteractiveCodeEditor provides an interactive environment for writing, debugging, and deploying smart contracts.
type InteractiveCodeEditor struct {
	CurrentCode     string         // The current code being written or edited
	CurrentContract *SmartContract // The contract being worked on
	LedgerInstance  *Ledger        // Ledger to store deployed contracts
	mutex           sync.Mutex     // Mutex for thread-safe operations
	Bytecode        string         // Compiled bytecode of the smart contract
}

// SnapshotManager handles the creation and restoration of blockchain snapshots.
type SnapshotManager struct {
	LedgerInstance  *Ledger           // Ledger instance to record all snapshots
	SnapshotStorage map[string][]byte // Stores encrypted snapshots
	CurrentState    *BlockchainState  // Current blockchain state
	mutex           sync.Mutex        // Ensures thread-safe snapshot creation and restoration
}

// BlockchainState represents the state of the blockchain at a given point in time.
type BlockchainState struct {
	BlockHeight int64     // The block height at the time of the snapshot
	Hash        string    // The hash of the blockchain state
	Timestamp   time.Time // The timestamp of the blockchain snapshot
}

// BytecodeGenerator generates bytecode for smart contracts and virtual machine execution.
type BytecodeGenerator struct {
	LedgerInstance *Ledger    // Ledger instance for logging bytecode deployments
	mutex          sync.Mutex // Mutex for thread-safe bytecode generation
}

// VirtualMachineConcurrency handles the concurrency mechanisms within the virtual machine.
type VirtualMachineConcurrency struct {
	LedgerInstance *Ledger    // Ledger to log transactions, sub-blocks, and blocks
	mutex          sync.Mutex // Ensures thread-safety across operations
}

type BytecodeInterpreter struct {
	LedgerInstance *Ledger                // Ledger instance to log contract executions
	mutex          sync.Mutex             // Mutex for thread-safe bytecode execution
	State          map[string]interface{} // State of the contract, storing key-value pairs
}

// CodeQualityAssurance handles the verification and validation of smart contract code quality.
type CodeQualityAssurance struct {
	LedgerInstance *Ledger    // Ledger instance for logging quality checks
	mutex          sync.Mutex // Mutex for thread-safe quality assurance
}

// CompilationDebugger handles the compilation and debugging of smart contract bytecode.
type CompilationDebugger struct {
	LedgerInstance *Ledger    // Ledger instance for logging compilation and debugging
	mutex          sync.Mutex // Mutex for thread-safe compilation and debugging
}

// GasManager manages gas fees for smart contract executions in the virtual machine.
type GasManager struct {
	LedgerInstance  *Ledger    // Ledger instance to track gas usage and fees
	mutex           sync.Mutex // Mutex for thread-safe operations
	GasPrice        float64    // Current gas price (in terms of native currency)
	ConsensusEngine *SynnergyConsensus
}

// GoContractCompiler manages the compilation and deployment of Go-based smart contracts.
type GoContractCompiler struct {
	LedgerInstance *Ledger           // Ledger instance for logging contract deployments
	CompiledCode   map[string]string // Stores compiled bytecode for each contract
	mutex          sync.Mutex        // Mutex for thread-safe operations
}

// GoSupport manages Go contract execution and validation within the Synnergy Network's virtual machine.
type GoSupport struct {
	LedgerInstance *Ledger     // Ledger for logging transactions and contract activities
	SubBlocks      []*SubBlock // List of sub-blocks pending full block aggregation
	mutex          sync.Mutex  // Mutex for thread-safe operations
}

// JavaScriptCompiler manages the compilation, execution, and deployment of JavaScript smart contracts.
type JavaScriptCompiler struct {
	LedgerInstance *Ledger           // Ledger instance for logging contract activities
	CompiledCode   map[string]string // Stores compiled bytecode for each contract
	mutex          sync.Mutex        // Mutex for thread-safe operations
}

// JavaScriptSupport manages JavaScript contract execution and validation within the Synnergy Network.
type JavaScriptSupport struct {
	LedgerInstance *Ledger     // Ledger for logging transactions and contract activities
	SubBlocks      []*SubBlock // List of sub-blocks pending full block aggregation
	mutex          sync.Mutex  // Mutex for thread-safe operations
}

// ReversalRequest represents a request to reverse a transaction.
type ReversalRequest struct {
	TransactionID string    // Unique transaction ID
	WalletID      string    // Wallet ID requesting the reversal
	Reason        string    // Reason for reversal
	Status        string    // Current status of the request (e.g., pending, approved, rejected)
	Timestamp     time.Time // Time when the request was created
}

// RustCompiler manages the compilation, execution, and deployment of Rust-based smart contracts.
type RustCompiler struct {
	LedgerInstance *Ledger           // Ledger instance for logging contract activities
	CompiledCode   map[string]string // Stores compiled bytecode for each contract
	mutex          sync.Mutex        // Mutex for thread-safe operations
}

// RustSupport manages Rust-based contract execution and validation within the Synnergy Network.
type RustSupport struct {
	LedgerInstance *Ledger     // Ledger for logging transactions and contract activities
	SubBlocks      []*SubBlock // List of sub-blocks pending full block aggregation
	mutex          sync.Mutex  // Mutex for thread-safe operations
}

// Sandbox is an isolated virtual environment for testing and executing smart contracts.
type Sandbox struct {
	ID               string                         // Unique identifier for the sandbox instance
	SmartContracts   map[string]*SmartContract      // Loaded smart contracts in the sandbox
	ExecutionHistory []SmartContractExecutionRecord // Records of contract executions within the sandbox
	LedgerInstance   *Ledger                        // Sandbox-associated ledger for tracking state
	VirtualMachine   *VirtualMachine                // Virtual machine for executing smart contracts
	mutex            sync.Mutex                     // Mutex for thread-safe operations
	IsActive         bool                           // Indicates if the sandbox is active or paused
	CreatedAt        time.Time                      // Time of sandbox creation
}

// VMResourceUsage tracks the resource usage for a virtual machine.
type VMResourceUsage struct {
	CPUUsage      float64            // CPU usage percentage
	MemoryUsage   float64            // Memory usage in MB
	DiskUsage     float64            // Disk usage in MB
	NetworkUsage  float64            // Network bandwidth usage in Mbps
	Uptime        time.Duration      // Total uptime of the VM
	PeakResources map[string]float64 // Peak resource usage metrics
}

// ExecutionState represents the execution state of a virtual machine.
type ExecutionState struct {
	State          string    // Current state (e.g., "idle", "running", "error")
	LastCheckpoint string    // ID of the last successful checkpoint
	ExecutionLogs  []string  // Logs of recent executions
	LastUpdated    time.Time // Timestamp of the last state update
	ErrorMessage   string    // Error message, if any
}

// FaultRecoveryRecord stores information about fault recovery mechanisms.
type FaultRecoveryRecord struct {
	RecoveryID      string    // Unique ID for the recovery instance
	FaultType       string    // Type of fault (e.g., "crash", "timeout", "memory overflow")
	RecoveryActions []string  // Actions taken to recover
	RecoveryStatus  string    // Status of the recovery (e.g., "successful", "failed", "in-progress")
	Timestamp       time.Time // Timestamp of the fault
	ResolvedBy      string    // Identifier of the responsible recovery process
}

// VMConfiguration contains the configuration details for a virtual machine.
type VMConfiguration struct {
	VMID            string                 // Unique identifier for the VM
	AllocatedCPU    float64                // Allocated CPU cores
	AllocatedMemory float64                // Allocated memory in MB
	AllocatedDisk   float64                // Allocated disk space in GB
	NetworkLimits   float64                // Network bandwidth limits in Mbps
	CustomSettings  map[string]interface{} // Additional custom settings
	CreatedAt       time.Time              // Timestamp when the configuration was created
}

// VMPerformanceMetrics represents performance metrics of a virtual machine.
type VMPerformanceMetrics struct {
	Throughput         float64   // Transactions or operations per second
	Latency            float64   // Average latency in milliseconds
	ErrorRate          float64   // Percentage of errors
	CPUUtilization     float64   // Average CPU utilization percentage
	MemoryUtilization  float64   // Average memory utilization percentage
	ResourceEfficiency float64   // Efficiency score for resource utilization
	LastUpdated        time.Time // Timestamp of the last metrics update
}

// VMLogEntry captures operational details or events for a virtual machine.
type VMLogEntry struct {
	LogID     string    // Unique ID for the log entry
	Timestamp time.Time // Timestamp of the log event
	VMID      string    // ID of the associated VM
	Event     string    // Description of the logged event
	Severity  string    // Severity level (e.g., "info", "warning", "error")
	Details   string    // Additional details about the event
}

// VMEventLog records significant events related to a virtual machine.
type VMEventLog struct {
	EventID     string    // Unique identifier for the event
	VMID        string    // ID of the VM related to the event
	EventType   string    // Type of event (e.g., "start", "stop", "error")
	Timestamp   time.Time // Time when the event occurred
	Initiator   string    // Entity or user who triggered the event
	Description string    // Additional description of the event
	Impact      string    // Impact of the event (e.g., "high", "low", "none")
}

// SandboxManager manages multiple sandboxes for isolated smart contract testing and execution.
type SandboxManager struct {
	LedgerInstance *Ledger             // Ledger for logging sandbox actions
	Sandboxes      map[string]*Sandbox // Map of sandbox instances
	mutex          sync.Mutex          // Mutex for thread-safety
}

// SolidityCompiler manages the compilation and deployment of Solidity smart contracts.
type SolidityCompiler struct {
	LedgerInstance *Ledger           // Ledger instance to log contract deployment
	CompiledCode   map[string]string // Stores the compiled bytecode for each contract
	mutex          sync.Mutex        // Mutex for thread-safety during compilation and deployment
}

// SoliditySupport manages Solidity contract execution and validation within the virtual machine.
type SoliditySupport struct {
	LedgerInstance *Ledger    // Ledger for logging transactions and contract activities
	mutex          sync.Mutex // Mutex for thread-safe operations
}

// SyntaxChecker is responsible for verifying the syntax of smart contracts in multiple languages
// (Solidity, JavaScript, Rust, Golang, Yul) and ensuring transaction compliance.
type SyntaxChecker struct {
	LedgerInstance *Ledger    // Ledger for logging results and transaction histories
	mutex          sync.Mutex // Mutex for thread-safety
}

// YulCompiler manages the compilation, execution, and deployment of Yul-based smart contracts.
type YulCompiler struct {
	LedgerInstance *Ledger           // Ledger instance for logging contract activities
	CompiledCode   map[string]string // Stores compiled bytecode for each contract
	mutex          sync.Mutex        // Mutex for thread-safe operations
}

// YulSupport manages Yul contract execution and validation within the Synnergy Network.
type YulSupport struct {
	LedgerInstance *Ledger     // Ledger for logging transactions and contract activities
	SubBlocks      []*SubBlock // List of sub-blocks pending full block aggregation
	Blocks         []Block     // List of aggregated blocks
	mutex          sync.Mutex  // Mutex for thread-safe operations
}

// ************** Wallet Structs **************

// Wallet represents a user's wallet containing private and public keys.
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey // The private key used to sign transactions.
	PublicKey  []byte            // The public key used to receive funds.
	Address    string            // Wallet address derived from the public key.
	Ledger     *Ledger           // Ledger instance for interaction.
}

// MultiSigWallets represents multi-signature wallet management in the ledger.
type MultiSigWallets struct {
	WalletID           string               // Unique identifier for the multi-signature wallet
	Owners             map[string]float64   // Wallet owners and their voting weights (address -> weight)
	RequiredSignatures int                  // Minimum number of signatures required to execute a transaction
	Transactions       []PendingTransaction // List of pending transactions requiring approval
	CreatedAt          time.Time            // Timestamp of wallet creation
	UpdatedAt          time.Time            // Timestamp of last modification
}

// PendingTransaction represents a transaction waiting for approval in a multi-signature wallet.
type PendingTransaction struct {
	TransactionID string          // Unique identifier for the transaction
	Initiator     string          // Address of the user who initiated the transaction
	Signatures    map[string]bool // Map of addresses to their approval status (true if signed)
	Payload       string          // Encoded payload of the transaction
	CreatedAt     time.Time       // Timestamp of transaction initiation
	ExpiresAt     time.Time       // Timestamp when the transaction expires
	Status        string          // Status of the transaction (e.g., "Pending", "Approved", "Rejected")
}

// WalletData represents data stored for a user's wallet.
type WalletData struct {
	WalletID           string              // Unique identifier for the wallet
	OwnerAddress       string              // Address of the wallet owner
	Balance            *big.Int            // Balance of the wallet in the main currency
	TokenBalances      map[string]*big.Int // Token balances
	TransactionHistory []TransactionRecord // Transaction history associated with the wallet
	IsVerified         bool                // Whether the wallet is verified
	CreatedAt          time.Time           // Time the wallet was created
	TokenID            string              // Token ID associated with the wallet (Add this field)
	Verified           bool                // Whether the wallet has been verified (Add this field)

}

// WalletConnection represents the connection details between a wallet and a system.
type WalletConnection struct {
	ConnectionID string    // Unique identifier for the connection
	WalletID     string    // Wallet associated with the connection
	ConnectedAt  time.Time // Timestamp of the connection
	Status       string    // Connection status (e.g., "active", "disconnected")
	IPAddress    string    // IP address from which the connection originated
}

// CurrencyExchange represents a currency exchange transaction.
type CurrencyExchange struct {
	ExchangeID      string    // Unique identifier for the exchange
	FromCurrency    string    // Currency being exchanged
	ToCurrency      string    // Currency being received
	Amount          *big.Int  // Amount being exchanged
	ExchangedAmount *big.Int  // Amount received after the exchange
	ExchangeRate    float64   // Exchange rate applied
	ExecutedAt      time.Time // Timestamp of the exchange
}

// ConnectionEvent represents an event related to wallet connections.
type ConnectionEvent struct {
	EventID      string // Unique identifier for the event
	ConnectionID string
	WalletID     string    // Wallet associated with the event
	EventType    string    // Type of event (e.g., "connection", "disconnection")
	EventTime    time.Time // Timestamp of the event
	Details      string    // Additional details related to the event
}

// MintRecord represents a record of token minting.
type MintRecord struct {
	RecordID  string    // Unique identifier for the mint record
	TokenID   string    // ID of the token being minted
	Amount    *big.Int  // Amount of tokens minted
	MintedBy  string    // Address of the entity that minted the tokens
	Timestamp time.Time // Timestamp of the minting
}

// BurnRecord represents a record of token burning.
type BurnRecord struct {
	RecordID  string    // Unique identifier for the burn record
	TokenID   string    // ID of the token being burned
	Amount    *big.Int  // Amount of tokens burned
	BurnedBy  string    // Address of the entity that burned the tokens
	Timestamp time.Time // Timestamp of the burning
}

// ContractExecutionLog represents a log of a smart contract execution.
type ContractExecutionLog struct {
	LogID         string    // Unique identifier for the log
	ContractID    string    // ID of the contract being executed
	ExecutedBy    string    // Address of the entity executing the contract
	ExecutionTime time.Time // Timestamp of the execution
	InputData     string    // Input data for the contract execution
	OutputData    string    // Output data from the contract execution
	Status        string    // Status of the execution (e.g., "success", "failure")
}

// ************** ZK-Proof Structs **************

// ZKProofRecord stores a record of the generated zero-knowledge proof.
type ZKProofRecord struct {
	ProofID       string    // Unique ID of the proof
	TransactionID string    // Transaction ID linked to the proof
	ProofData     string    // The actual zero-knowledge proof data (could be in a specific format)
	GeneratedAt   time.Time // Timestamp when the proof was generated
	ProofStatus   string    // Status of the proof (e.g., "generated", "validated", "invalidated")
}

// ZKValidationRecord stores a record of zk-proof validation.
type ZKValidationRecord struct {
	ProofID     string    // ID of the proof being validated
	ValidatorID string    // ID of the entity that validated the proof
	IsValid     bool      // Whether the proof passed validation or not
	ValidatedAt time.Time // Timestamp of the validation
}

// ZkProof represents a zero-knowledge proof used in the blockchain system
type ZkProof struct {
	ProofID    string             // Unique identifier for the zk-proof
	ProverID   string             // ID of the prover generating the proof
	VerifierID string             // ID of the verifier validating the proof
	ProofData  []byte             // The actual zk-proof data
	IsValid    bool               // Whether the proof has been validated
	VerifiedAt time.Time          // Time when the proof was verified
	Ledger     *Ledger            // Reference to the ledger for recording proof events
	Encryption *Encryption        // Encryption service for securing proof data
	Consensus  *SynnergyConsensus // Consensus system for validating zk-proofs
	mu         sync.Mutex         // Mutex for concurrency control
}

type ZKOracleData struct {
	OracleID    string    `json:"oracle_id"`
	DataFeedID  string    `json:"data_feed_id"`
	ZKProof     []byte    `json:"zk_proof"` // Change to []byte
	DataPayload string    `json:"data_payload"`
	Verified    bool      `json:"verified"`
	Timestamp   time.Time `json:"timestamp"`
	HandlerNode string    `json:"handler_node"`
}

// ************** Utility Structs **************

// SystemHookRegistry stores and manages hooks for system events and transitions.
type SystemHookRegistry struct {
	HookID         string          // Unique identifier for the hook.
	HookType       string          // Type of hook (e.g., pre-action, post-action, error-handler).
	TargetEvent    string          // Event or action the hook is associated with.
	ExecutionOrder int             // Order of execution for the hook.
	ActiveHooks    map[string]Hook // Active hooks by ID.
	Logs           []HookLog       // Logs for hook executions.
}

type MetadataSummary struct {
	ID          string            // Unique identifier for the metadata
	Name        string            // Name of the metadata
	Description string            // Description of the metadata
	Tags        []string          // Associated tags for categorization
	Attributes  map[string]string // Key-value pairs for metadata attributes
	CreatedAt   time.Time         // Timestamp when the metadata was created
	UpdatedAt   time.Time         // Timestamp when the metadata was last updated
	Version     int               // Version number for tracking changes
}

type BlockHeader struct {
	BlockID          string    // Unique identifier for the block
	ParentBlockID    string    // ID of the parent block
	Height           int       // Block height in the chain
	Timestamp        time.Time // Timestamp when the block was created
	MerkleRoot       string    // Merkle root of the block's transactions
	Validator        string    // Validator ID that created the block
	Signature        string    // Cryptographic signature of the block
	TransactionCount int       // Number of transactions included in the block
	Difficulty       int       // Difficulty level for mining or validation
	Nonce            string    // Nonce used in Proof-of-Work or other consensus
}

// Hook represents a single system hook.
type Hook struct {
	ID        string    // Unique identifier for the hook.
	Function  func()    // Function to execute when the hook is triggered.
	CreatedAt time.Time // Timestamp when the hook was created.
}

// HookLog captures the execution details of a system hook.
type HookLog struct {
	LogID     string    // Unique identifier for the log entry.
	HookID    string    // Associated hook ID.
	Message   string    // Log message.
	Timestamp time.Time // Timestamp of the log entry.
}

// InterruptManager manages system interrupts and their responses.
type InterruptManager struct {
	InterruptQueue []InterruptEvent  // Queue of active interrupts.
	Policies       map[string]string // Policies for handling interrupts.
	Logs           []InterruptLog    // Logs for interrupt handling.
}

// InterruptEvent represents a single interrupt event.
type InterruptEvent struct {
	ID          string    // Unique identifier for the interrupt.
	Description string    // Description of the interrupt.
	Severity    string    // Severity level of the interrupt.
	Timestamp   time.Time // Timestamp of the interrupt occurrence.
}

// InterruptLog captures the handling details of interrupts.
type InterruptLog struct {
	LogID       string    // Unique identifier for the log entry.
	InterruptID string    // Associated interrupt ID.
	Action      string    // Action taken in response to the interrupt.
	Timestamp   time.Time // Timestamp of the log entry.
}

// ShutdownManager oversees system shutdowns and restarts.
type ShutdownManager struct {
	ShutdownID    string        // Unique identifier for the shutdown process.
	ScheduledTime time.Time     // Scheduled time for the shutdown.
	Reason        string        // Reason for the shutdown.
	Status        string        // Status of the shutdown (e.g., pending, in-progress, completed).
	RecoverySteps []string      // List of recovery steps after the shutdown.
	Logs          []ShutdownLog // Logs for the shutdown process.
}

// ShutdownLog captures events during the shutdown process.
type ShutdownLog struct {
	LogID     string    // Unique identifier for the log entry.
	Message   string    // Log message.
	Timestamp time.Time // Timestamp of the log entry.
}

// ConstantsManager stores and manages system-wide constants.
type ConstantsManager struct {
	Constants     map[string]interface{} // Key-value pairs of constants.
	LastUpdated   time.Time              // Timestamp of the last update to constants.
	UpdateHistory []ConstantUpdateLog    // Logs of updates to constants.
}

// ConstantUpdateLog records updates to constants.
type ConstantUpdateLog struct {
	LogID     string      // Unique identifier for the log entry.
	Key       string      // Name of the constant updated.
	OldValue  interface{} // Previous value of the constant.
	NewValue  interface{} // New value of the constant.
	UpdatedAt time.Time   // Timestamp of the update.
}

// DebugManager oversees system debugging sessions and logs.
type DebugManager struct {
	DebugSessions []DebugSession // Active and completed debug sessions.
	Logs          []DebugLog     // Logs for debugging activities.
}

// DebugSession represents a single debugging session.
type DebugSession struct {
	SessionID   string    // Unique identifier for the debug session.
	Description string    // Description of the debug session.
	Status      string    // Status of the debug session (e.g., active, completed).
	StartTime   time.Time // Timestamp when the session started.
	EndTime     time.Time // Timestamp when the session ended (if completed).
}

// DebugLog captures events during debugging.
type DebugLog struct {
	LogID     string    // Unique identifier for the log entry.
	SessionID string    // Associated debug session ID.
	Message   string    // Log message.
	Timestamp time.Time // Timestamp of the log entry.
}

// SelfTestManager manages self-tests and health diagnostics.
type SelfTestManager struct {
	TestSchedules []TestSchedule   // List of scheduled tests.
	TestResults   []SelfTestResult // Results of completed tests.
	Logs          []SelfTestLog    // Logs for self-test operations.
}

// TestSchedule represents a scheduled self-test.
type TestSchedule struct {
	TestID      string    // Unique identifier for the test.
	TestName    string    // Name of the self-test.
	ScheduledAt time.Time // Timestamp of the scheduled test.
	Status      string    // Status of the test (e.g., pending, completed).
}

// SelfTestLog records events related to self-tests.
type SelfTestLog struct {
	LogID     string    // Unique identifier for the log entry.
	TestID    string    // Associated test ID.
	Message   string    // Log message.
	Timestamp time.Time // Timestamp of the log entry.
}

// EventManager oversees the registration and execution of system events.
type EventManager struct {
	RegisteredEvents map[string]EventDetails // Map of registered events by ID.
	EventTriggers    []EventTrigger          // List of active triggers for events.
	Logs             []EventLog              // Logs for event-related activities.
}

// EventDetails captures details about a registered event.
type EventDetails struct {
	EventID       string    // Unique identifier for the event.
	Description   string    // Description of the event.
	CreatedAt     time.Time // Timestamp when the event was created.
	LastTriggered time.Time // Timestamp of the last event trigger.
}

// TaskManager manages system tasks and their execution.
type TaskManager struct {
	Tasks     map[string]TaskDetails // Map of tasks by ID.
	TaskQueue []string               // Queue of pending tasks.
	Logs      []TaskLog              // Logs for task operations.
}

// TaskDetails captures details about a system task.
type TaskDetails struct {
	TaskID      string    // Unique identifier for the task.
	Description string    // Description of the task.
	AssignedTo  string    // Entity or module assigned to the task.
	Status      string    // Status of the task (e.g., pending, in-progress, completed).
	CreatedAt   time.Time // Timestamp when the task was created.
	CompletedAt time.Time // Timestamp when the task was completed (if applicable).
}

// TaskLog captures events during task execution.
type TaskLog struct {
	LogID     string    // Unique identifier for the log entry.
	TaskID    string    // Associated task ID.
	Message   string    // Log message.
	Timestamp time.Time // Timestamp of the log entry.
}

// EventScheduler manages the scheduling of system events.
type EventScheduler struct {
	ScheduledEvents map[string]ScheduledEvent // Map of scheduled events by ID.
	Logs            []ScheduleLog             // Logs for scheduling operations.
}

// ScheduleLog captures events during event scheduling.
type ScheduleLog struct {
	LogID     string    // Unique identifier for the log entry.
	EventID   string    // Associated event ID.
	Message   string    // Log message.
	Timestamp time.Time // Timestamp of the log entry.
}

// TimeManager manages system time and time-based operations.
type TimeManager struct {
	CurrentTime  time.Time         // Current synchronized system time.
	LastSyncTime time.Time         // Timestamp of the last time synchronization.
	TimeZones    map[string]string // Supported time zones and offsets.
	Logs         []TimeLog         // Logs for time-related operations.
}

// TimeLog captures events related to time synchronization and updates.
type TimeLog struct {
	LogID     string    // Unique identifier for the log entry.
	Message   string    // Log message.
	Timestamp time.Time // Timestamp of the log entry.
}

// BackupManager manages the creation, storage, and restoration of blockchain backups to ensure data integrity and system recovery.
type BackupManager struct {
	BackupID          string             // Unique identifier for the backup
	BackupLocation    string             // Location where the backup is stored (e.g., cloud storage, external server)
	BackupFrequency   time.Duration      // Frequency of automatic backups (e.g., daily, weekly)
	BackupRetention   int                // Number of backups to retain (older backups will be deleted)
	LastBackup        time.Time          // Timestamp of the last completed backup
	BackupSize        int64              // Size of the most recent backup in bytes
	RedundancyFactor  int                // Level of redundancy in the backup (e.g., how many copies are made)
	EncryptionService *Encryption        // Encryption service for securing backups
	BackupLogs        []BlockchainBackup // Log of all backup events, including successes and failures
	mutex             sync.Mutex         // Mutex for thread-safe operations
}

// RecoveryPlans define the strategy and protocols for recovering the blockchain in the event of a failure or disaster.
type RecoveryPlans struct {
	PlanID               string         // Unique identifier for the recovery plan
	RecoverySteps        []string       // List of steps to be followed during recovery
	BackupManager        *BackupManager // Reference to the BackupManager for restoration
	ContactList          []string       // List of contacts to notify during recovery
	IncidentResponseTeam []string       // Team responsible for managing recovery efforts
	RecoveryTests        []time.Time    // Dates when recovery tests were performed
	LastTestResult       string         // Result of the most recent recovery test (e.g., "Success", "Failure")
	EstimatedDowntime    time.Duration  // Estimated time for full recovery
	EncryptionService    *Encryption    // Encryption service for securing sensitive recovery data
	RecoveryLogs         []RecoveryLog  // Log of recovery events and actions
}

// RecoveryLog represents an event or action that occurred during the recovery process.
type RecoveryLog struct {
	EventID     string    // Unique identifier for the recovery event
	EventType   string    // Type of event (e.g., "Test", "Actual Recovery", "Update")
	Timestamp   time.Time // Time when the event occurred
	Description string    // Detailed description of the event or action
	Status      string    // Status of the event (e.g., "Success", "Failed", "In Progress")
}

// SensorManager handles the integration, monitoring, and management of environmental sensors in the network.
type SensorManager struct {
	SensorID          string             // Unique identifier for the sensor
	SensorType        string             // Type of sensor (e.g., "Temperature", "Humidity", "Air Quality")
	Location          string             // Physical or logical location of the sensor
	SensorReadings    map[string]float64 // Map of sensor readings (e.g., timestamped data points)
	LastReadingTime   time.Time          // Timestamp of the last reading taken from the sensor
	Thresholds        map[string]float64 // Warning or critical thresholds for sensor data
	AlertsEnabled     bool               // Whether alerts are enabled for the sensor
	EncryptionService *Encryption        // Encryption service for securing sensor data
	NetworkManager    *NetworkManager    // Manages communication with other nodes and systems
	SensorLogs        []SensorLog        // Log of sensor events (e.g., readings, alerts, maintenance)
	mutex             sync.Mutex         // Mutex for thread-safe sensor operations

}

// SensorLog represents a log entry for a sensor event, such as a reading or alert.
type SensorLog struct {
	EventID   string    // Unique identifier for the sensor event
	EventType string    // Type of event (e.g., "Reading", "Alert", "Maintenance")
	Timestamp time.Time // Time when the event occurred
	Value     float64   // Sensor reading or event value
	Status    string    // Status of the event (e.g., "Normal", "Alert", "Critical")
}

// Content represents large or complex data types that are stored, managed, and linked to blockchain transactions.
type Content struct {
	ContentID         string            // Unique identifier for the content
	ContentType       string            // Type of content (e.g., "Video", "Image", "Document")
	OwnerWallet       string            // Wallet address of the owner of the content
	EncryptedData     []byte            // Encrypted data representing the content
	StorageLocation   string            // Location where the content is stored (e.g., decentralized storage)
	AccessLogs        []AccessLog       // Log of access events for the content
	Metadata          map[string]string // Metadata about the content (e.g., title, description, file size)
	EncryptionService *Encryption       // Encryption service for securing content data
	AccessControl     AccessControl     // Access control rules for the content
}

// AlertSystem manages alerts for anomalies, security threats, and other important events in the blockchain network.
type AlertSystem struct {
	AlertID         string     // Unique identifier for the alert
	AlertType       string     // Type of alert (e.g., "Security", "Network", "Performance")
	Severity        string     // Severity of the alert (e.g., "Low", "Medium", "High", "Critical")
	TriggeredBy     string     // The entity (node, transaction, etc.) that triggered the alert
	Timestamp       time.Time  // Timestamp when the alert was triggered
	Description     string     // Detailed description of the alert
	Status          string     // Status of the alert (e.g., "Open", "Resolved", "In Progress")
	ResponseActions []string   // Actions taken in response to the alert
	EscalationLevel int        // Level of escalation (e.g., 1 = low, 5 = high)
	Logs            []AlertLog // Log of events related to the alert
}

// AlertLog represents an event or action related to an alert.
type AlertLog struct {
	LogID     string    // Unique identifier for the log entry
	EventType string    // Type of event (e.g., "Alert Raised", "Response Action")
	Timestamp time.Time // Timestamp when the event occurred
	Details   string    // Details of the event or action
}

// RedundantStorage manages the distribution of data across multiple storage locations to ensure redundancy and availability.
type RedundantStorage struct {
	StorageID         string           // Unique identifier for the redundant storage system
	PrimaryLocation   string           // Primary storage location
	BackupLocations   []string         // List of backup storage locations for redundancy
	ReplicationFactor int              // Number of redundant copies of data to be maintained
	EncryptionService *Encryption      // Encryption service for securing stored data
	StorageUsage      map[string]int64 // Map of storage locations and their current usage
	SyncInterval      time.Duration    // Interval for syncing data between locations
	StorageLogs       []StorageLog     // Log of storage events and activities
}

// StorageLog represents an event or action related to the storage system.
type StorageLog struct {
	LogID     string    // Unique identifier for the log entry
	EventType string    // Type of event (e.g., "Data Replicated", "Data Restored", "Error")
	Timestamp time.Time // Timestamp when the event occurred
	Details   string    // Details of the event or action
}

// NotaryService provides certified notarization services for blockchain transactions and records.
type NotaryService struct {
	NotaryID              string                        // Unique identifier for the notary service
	CertifiedBlocks       map[string]*Block             // Certified blocks that have been notarized
	CertifiedTransactions map[string]*TransactionRecord // Certified transactions that have been notarized
	Timestamp             time.Time                     // Timestamp of the last notarization
	NotaryLog             []NotaryLog                   // Log of notarized events and actions
	EncryptionService     *Encryption                   // Encryption service for securing notarized records
}

// NotaryLog represents an event where a transaction or block was notarized.
type NotaryLog struct {
	LogID     string    // Unique identifier for the log entry
	EventType string    // Type of event (e.g., "Block Notarized", "Transaction Notarized")
	Timestamp time.Time // Timestamp when the event occurred
	BlockID   string    // ID of the block that was notarized
	TxID      string    // ID of the transaction that was notarized
}

// LoadBalancer distributes network traffic and requests across nodes to ensure even distribution of load.
type LoadBalancer struct {
	BalancerID      string                 // Unique identifier for the load balancer
	Strategy        string                 // Load-balancing strategy (e.g., "Round-Robin", "Weighted", "Random")
	NodeMetrics     map[string]NodeMetrics // Performance metrics for each node used in decision-making
	ActiveNodes     []string               // List of currently active nodes in the network
	RequestQueue    []*APIRequest          // Queue of pending API requests to be distributed
	SyncInterval    time.Duration          // Interval for syncing load metrics with nodes
	LoadBalancerLog []LoadBalancerLog      // Log of load-balancing events
}

// LoadBalancerLog represents an event related to load-balancing operations.
type LoadBalancerLog struct {
	LogID     string    // Unique identifier for the log entry
	EventType string    // Type of event (e.g., "Request Routed", "Node Added", "Node Removed")
	Timestamp time.Time // Timestamp when the event occurred
	Details   string    // Details of the event or action
}

// APIRequest represents a request made to the blockchain network through the API gateway.
type APIRequest struct {
	RequestID    string                 // Unique identifier for the API request
	Endpoint     string                 // API endpoint that is being accessed (e.g., "/transaction", "/block")
	RequestData  map[string]interface{} // Data being sent with the API request
	Requestor    string                 // Wallet address or node ID making the request
	Timestamp    time.Time              // Timestamp when the request was made
	Status       string                 // Status of the request (e.g., "Pending", "Completed", "Failed")
	ResponseData map[string]interface{} // Data returned by the API request (if applicable)
	ErrorMessage string                 // Error message in case of failure
	ResponseTime time.Duration          // Time taken to process the request
}

// APIManager manages incoming API requests and ensures proper routing, validation, and processing within the blockchain network.
type APIManager struct {
	ManagerID      string                 // Unique identifier for the API manager
	ActiveRequests map[string]*APIRequest // Map of currently active API requests
	RateLimit      int                    // Maximum number of API requests allowed per second
	Timeout        time.Duration          // Timeout for API requests
	LoadBalancer   *LoadBalancer          // Load balancer for distributing API requests across nodes
	RequestLogs    []APIRequestLog        // Log of API requests and their outcomes
}

// APIRequestLog represents a log entry for an API request.
type APIRequestLog struct {
	LogID        string        // Unique identifier for the log entry
	RequestID    string        // ID of the API request
	Timestamp    time.Time     // Time when the request was logged
	Status       string        // Status of the request (e.g., "Success", "Failed")
	ResponseTime time.Duration // Time taken to process the request
	Details      string        // Additional details or error messages related to the request
}

// Notification represents a structure for sending notifications
type Notification struct {
	Recipient   string           // Email address or webhook URL
	Message     string           // Message to be sent
	Type        NotificationType // Type of notification (Email, Webhook, etc.)
	SentAt      time.Time        // Timestamp of when the notification was sent
	IsDelivered bool             // Status whether the notification was successfully delivered
}

// CreditCheckManager manages the decentralized credit checking process.
type CreditCheckManager struct {
	CreditScores     map[string]float64        // Mapping of wallet addresses to credit scores
	CreditHistory    map[string][]CreditRecord // Credit history for each participant
	CreditBureaus    []string                  // List of trusted decentralized credit bureaus
	VerificationLogs []CreditVerificationLog   // Logs of all credit checks performed
	mutex            sync.Mutex                // Mutex for thread-safe operations
}

// CreditRecord stores an individual's credit record.
type CreditRecord struct {
	DateChecked time.Time // Date when the credit check was performed
	Score       float64   // Credit score at the time of checking
	Report      string    // Detailed credit report
}

// CreditVerificationLog stores logs for credit checks.
type CreditVerificationLog struct {
	WalletAddress    string    // Wallet address being verified
	CreditScore      float64   // The score verified
	VerificationDate time.Time // Date of verification
}

// TermsCustomizationManager manages loan term customization and negotiations.
type TermsCustomizationManager struct {
	AvailableTerms  []LoanTerm                       // List of available loan terms
	Customizations  map[string]LoanTermCustomization // Customization for each loan proposal
	NegotiationLogs []TermNegotiationLog             // Records of all term customizations and negotiations
	mutex           sync.Mutex                       // Mutex for thread-safe operations
}

// LoanTerm represents basic loan terms.
type LoanTerm struct {
	InterestRate      float64 // Interest rate for the loan
	Duration          int     // Duration of the loan in months
	RepaymentSchedule string  // Repayment schedule (e.g., Monthly, Quarterly)
}

// LoanTermCustomization stores customized terms for a loan proposal.
type LoanTermCustomization struct {
	ProposalID        string    // ID of the loan proposal
	CustomTerms       LoanTerm  // Customized terms agreed between the lender and borrower
	CustomizationDate time.Time // Date when the customization was agreed upon
}

// TermNegotiationLog logs the customization negotiation process.
type TermNegotiationLog struct {
	ProposalID      string    // Loan proposal ID
	OriginalTerms   LoanTerm  // Original terms of the loan
	CustomTerms     LoanTerm  // Customized terms after negotiation
	NegotiationDate time.Time // Date of negotiation
}

// SYN900Registry stores and manages SYN900 record-keeping and verifications.
type SYN900Registry struct {
	RegisteredEntities map[string]bool         // Wallets or entities registered in SYN900
	VerificationLogs   []SYN900VerificationLog // Logs of all verifications
	mutex              sync.Mutex              // Mutex for thread-safe operations
}

// SYN900VerificationLog logs SYN900 verifications.
type SYN900VerificationLog struct {
	WalletAddress string    // Address being verified
	VerifiedAt    time.Time // Time when verification took place
	VerifiedBy    string    // Verifier's wallet or entity
}

// MessageList represents a doubly linked list for efficient message queuing.
type MessageList struct {
	Head  *MessageNode // Points to the first message in the list
	Tail  *MessageNode // Points to the last message in the list
	Count int          // Number of messages in the list
}

// MessageNode represents a node in the message queue.
type MessageNode struct {
	Data interface{}  // Message data
	Prev *MessageNode // Pointer to the previous node
	Next *MessageNode // Pointer to the next node
}

// NotificationType represents the type of notification (Email, Webhook, etc.).
type NotificationType string

const (
	EmailNotification   NotificationType = "Email"
	WebhookNotification NotificationType = "Webhook"
	SMSNotification     NotificationType = "SMS"
	PushNotification    NotificationType = "Push"
)

// SmartContractExecutionRecord stores details of smart contract executions.
type SmartContractExecutionRecord struct {
	ContractID      string    // Unique identifier of the smart contract
	ExecutedAt      time.Time // Timestamp when the contract was executed
	GasUsed         int64     // Amount of gas used during execution
	ExecutionResult string    // Result or output of the contract execution
	Executor        string    // Wallet address of the executor
	SandboxMode     bool      // Whether the execution occurred in sandbox mode
}

// TaskHistory represents the history of a task.
type TaskHistory struct {
	TaskID  string      // Unique identifier for the task
	Updates []TaskEvent // List of events related to the task
}

// TaskEvent represents an event or update in a task's lifecycle.
type TaskEvent struct {
	Timestamp time.Time // Timestamp of the event
	Status    string    // Status of the task at this event
	Comment   string    // Additional comments or details
}

// OrchestrationRecord represents a record of an orchestration action.
type OrchestrationRecord struct {
	RecordID  string    // Unique identifier for the record
	Action    string    // Action performed during orchestration
	Timestamp time.Time // Timestamp of the action
	Details   string    // Additional details about the orchestration action
}
