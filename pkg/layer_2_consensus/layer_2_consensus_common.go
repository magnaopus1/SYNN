package layer2_consensus

import (
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/common"

)

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
	Ledger              *ledger.Ledger                 // Ledger instance for tracking consensus transitions
	EncryptionService   *common.Encryption         // Encryption service for securing consensus-related data
	mu                  sync.Mutex                     // Mutex for concurrent management
}

// ConsensusStrategy defines the parameters of a consensus mechanism
type ConsensusStrategy struct {
	StrategyID     string    // Unique identifier for the strategy
	StrategyType   string    // Type of consensus mechanism (e.g., "PoS", "PoH", "Synnergy")
	CurrentUsage   float64   // Current resource usage (load) of the strategy
	LastHopped     time.Time // Last time this strategy was used
	Active         bool      // Whether this strategy is currently active
	HopCount       int       // Number of hops to this strategy
}

// DynamicConsensusManager manages dynamic consensus switching (hopping) between multiple strategies
type DynamicConsensusManager struct {
	Strategies        map[string]*ConsensusStrategy // All available consensus strategies
	ActiveStrategy    *ConsensusStrategy            // Currently active consensus strategy
	Ledger            *ledger.Ledger                // Ledger instance for tracking consensus hops
	EncryptionService *common.Encryption        // Encryption service for securing strategy-related data
	mu                sync.Mutex                    // Mutex for thread-safe management
}

// ConsensusLayer represents a consensus layer with adaptive properties
type ConsensusLayer struct {
	LayerID        string    // Unique identifier for the consensus layer
	LayerType      string    // Type of consensus layer (e.g., "PoS", "PoH", "Synnergy")
	CurrentLoad    float64   // Current load on the consensus layer
	MaxLoad        float64   // Maximum load before transitioning to a new layer
	TransitionTime time.Time // Time of the last transition to this layer
	Active         bool      // Whether this layer is currently active
	TransitionCount int      // Number of transitions made to this layer
}

// ElasticConsensusManager manages the adaptive transitions between elastic consensus layers
type ElasticConsensusManager struct {
	ConsensusLayers   map[string]*ConsensusLayer // Available consensus layers
	ActiveLayer       *ConsensusLayer            // Currently active consensus layer
	Ledger            *ledger.Ledger             // Ledger instance for tracking transitions and decisions
	EncryptionService *common.Encryption     // Encryption service for securing consensus data
	mu                sync.Mutex                 // Mutex for concurrent management
}

// CollaborationTask represents an off-chain computational task that requires collaboration
type CollaborationTask struct {
	TaskID          string    // Unique ID for the collaboration task
	AssignedNodes   []string  // List of nodes assigned to the task
	ComputationResult string  // Result of the off-chain computation
	CompletionStatus string   // Status of the task ("Pending", "Completed")
	AssignedTime    time.Time // Time when the task was assigned
	CompletedTime   time.Time // Time when the task was completed
	EncryptedData   string    // Encrypted task details for security
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
	Ledger            *ledger.Ledger                // Ledger for recording PoCol actions
	EncryptionService *common.Encryption        // Encryption service for securing collaboration data
	mu                sync.Mutex                    // Mutex for concurrent task management
}
