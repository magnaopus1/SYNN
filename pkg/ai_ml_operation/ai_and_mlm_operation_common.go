package ai_ml_operation

import (
	"synnergy_network/pkg/ledger"
	"time"
)

// InferenceRecord holds inference details.
type InferenceRecord struct {
	ModelID    string
	Timestamp  time.Time
	Result     string
	NodeID     string
	Processed  bool
}

// AnalysisRecord for AI analysis sessions.
type AnalysisRecord struct {
	AnalysisID   string
	ModelID      string
	NodeID       string
	StartTime    time.Time
	StopTime     time.Time
	Status       string
}

// Recommendation struct for AI recommendations.
type Recommendation struct {
	ModelID    string
	Timestamp  time.Time
	Content    string
	Updated    bool
}


// Container represents a deployed container.
type Container struct {
	ContainerID string
	ModelID     string
	Status      string
	NodeID      string
}

// AiService represents a deployed AI service.
type AiService struct {
	ServiceID     string
	Status        string
	Metrics       ServiceMetrics
	LastUpdated   time.Time
}

// DataProcessingLog logs data processing activities.
type DataProcessingLog struct {
	ProcessID   string
	ModelID     string
	NodeID      string
	StartTime   time.Time
	EndTime     time.Time
	Status      string
}

// TrafficRecord logs load balancing and traffic operations.
type TrafficRecord struct {
	ModelID     string
	NodeID      string
	Action      string
	Timestamp   time.Time
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
	ModelID     string
	Restricted  bool
	Reason      string
	Timestamp   time.Time
}

// ModelPermissions represents permissions for a model.
type ModelPermissions struct {
	ModelID       string
	AllowedUsers  []string
}

// AccessToken represents access token for model operations.
type AccessToken struct {
	TokenID      string
	ModelID      string
	GrantedTo    string
	Expiry       time.Time
	Permissions  string
}

// ModelCheckpoint represents a model checkpoint.
type ModelCheckpoint struct {
	ModelID      string
	Version      int
	CreatedAt    time.Time
	DataHash     string
}

// AccessLog represents a log entry for model access.
type AccessLog struct {
	ModelID     string
	UserID      string
	Action      string
	Timestamp   time.Time
}

// UsageStatistics represents model usage statistics.
type UsageStatistics struct {
	ModelID       string
	UsageCount    int
	LastUsedAt    time.Time
	UsageDuration time.Duration
}

// PerformanceMetrics represents model performance metrics.
type PerformanceMetrics struct {
	ModelID       string
	Accuracy      float64
	Loss          float64
	LastUpdated   time.Time
}

// ComplianceCheck records model compliance checks.
type ComplianceCheck struct {
	ModelID       string
	Status        string
	Timestamp     time.Time
	Details       string
}

// SecurityAudit represents a security audit entry.
type SecurityAudit struct {
	ModelID       string
	Passed        bool
	Timestamp     time.Time
	Findings      string
}

// ResourceAllocation represents allocated resources.
type ResourceAllocation struct {
	ResourceID    string
	ModelID       string
	Amount        float64
	AllocatedAt   time.Time
}

// StorageAllocation represents allocated storage.
type StorageAllocation struct {
	StorageID     string
	ModelID       string
	SizeMB        int
	AllocatedAt   time.Time
}

// CacheData represents cached data.
type CacheData struct {
	ModelID       string
	DataID        string
	Data          string
	CreatedAt     time.Time
}

// DeploymentCheck records deployment checks.
type DeploymentCheck struct {
	DeploymentID  string
	Status        string
	Timestamp     time.Time
	Details       string
}

// TrainingStatus represents the training status of a model.
type TrainingStatus struct {
	ModelID       string
	Status        string
	LastUpdated   time.Time
}

// RunStatus represents the running status of a model.
type RunStatus struct {
	ModelID       string
	Status        string
	LastChecked   time.Time
}

type Model struct {
	ModelID          string            // Unique identifier for the model
	ModelName        string            // Name of the model
	Version          string            // Version of the model
	LastUpdated      time.Time         // Timestamp of the last update
	IsDeployed       bool              // Deployment status
	Capabilities     []string          // Capabilities like ["ImageClassification", "ObjectDetection"]
	ProcessingModes  []string          // Modes such as ["HighAccuracy", "FastProcessing"]
	Status           string            // Current model status (e.g., "idle", "running", "error")
	ModelType        string            // Type of model (e.g., "Deep Learning", "Machine Learning")
	Classification   string            // Classification category of the model (e.g., "Convolutional Neural Network")
	Purpose          string            // Model's intended purpose (e.g., "Image Recognition")
	Description      string            // Brief description of the model's functionality
	DocumentationURL string            // URL to detailed model documentation
	Permissions      map[string]string // Permissions, e.g., {"read": "public", "write": "restricted"}
	TrainingLevel    string            // Level of training applied (e.g., "Basic", "Advanced", "Specialized")
	LedgerInstance   *ledger.Ledger    // Instance of the ledger to record transactions
	ResourceUsage      map[string]float64 // Resource usage stats, e.g., {"CPU": 75.0, "Memory": 60.5}
	PerformanceMetrics map[string]float64 // Performance metrics, e.g., {"Latency": 1.2, "Throughput": 200}
	CurrentScale int    // Current scale level
	MaxScale    int     // Maximum allowed scale level
	MinScale    int     // Minimum allowed scale level
	RecommendationCache map[string][]byte // Caches recommendations by transaction ID
	InferenceCount   int // Track the number of inferences
	AnalysisSessions int // Track the number of analysis sessions
	PredictionCount  int // Track the number of predictions
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

// ModelIndex struct for AI/ML models.
type ModelIndex struct {
	ModelID      string    // Unique identifier for the model
	ModelName    string    // Name of the model
	NodeLocation string    // Node or server location where the model resides
	Load         int       // Load level or usage of the model
	Status       string    // Current status of the model in the network
	IPFSLink     string    // IPFS/Swarm link to the model’s data
	Timestamp    time.Time // Time when the model was added to the index
	TrainingHistory []string         // History of training sessions
	Permissions     map[string]string // Access permissions, e.g., {"user1": "read"}
	OptimizationHistory []string          // Optimization logs
	VersionHash         string            // Current version hash
	DeployedAt          time.Time         // Deployment timestamp
	OffChainLink string    // Off-chain storage link for model data (new field added)
	TestHistory   []string          // History of test result IDs
	AccessLogs    []time.Time       // Log of access times for monitoring

}

// RecommendationCriteria represents the structure for criteria used to generate recommendations.
type RecommendationCriteria struct {
	UserID             string            `json:"user_id"`
	PreferenceScore    float64           `json:"preference_score"`
	MaxRecommendations int               `json:"max_recommendations"`
	Preference         string            `json:"preference"`
	Threshold          float64           `json:"threshold"`
	MaxSuggestions     int               `json:"max_suggestions"`
	Tags               []string          `json:"tags"`
	ExcludeList        []string          `json:"exclude_list"`
	UserHistoryWeight  float64           `json:"user_history_weight"`
	ContextualFactors  map[string]string `json:"contextual_factors"`
}

// RecommendationUpdateData represents the new data for updating an existing recommendation.
type RecommendationUpdateData struct {
	UserID        string  `json:"user_id"`
	FeedbackScore float64 `json:"feedback_score"`
	NewCriteria   []byte  `json:"new_criteria"`
	UpdateReason  string  `json:"update_reason"`
}

