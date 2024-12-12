package ledger

import (
	"sync"
	"time"
)




/// ****** Standard Nodes ****** ///

// Standard node types
const (
    LightningNodeType NodeType = "LightningNode"
    LightNodeType     NodeType = "LightNode"
    FullNodeType      NodeType = "FullNode"
)


// APINode represents an API node in the Synnergy Network, facilitating high-throughput API requests and blockchain interaction.
type APINode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	APIManager        *APIManager
	mutex             sync.Mutex
	SyncInterval      time.Duration
	SNVM              *VirtualMachine
	ActiveRequests    map[string]*APIRequest
	RequestTimeout    time.Duration
	LoadBalancer      *LoadBalancer
}

// ArchivalWitnessNode represents an archival witness node in the Synnergy Network, providing certified archival services and historical accuracy.
type ArchivalWitnessNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	SubBlocks         map[string]*SubBlock
	mutex             sync.Mutex
	SyncInterval      time.Duration
	NotaryService     *NotaryService
	RedundantStorage  *RedundantStorage
	RequestTimeout    time.Duration
	SNVM              *VirtualMachine
}

// AuditNode represents an audit node in the Synnergy Network that continuously monitors and verifies blockchain activities.
type AuditNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	AuditTrail        map[string]*AuditEntry
	mutex             sync.Mutex
	SyncInterval      time.Duration
	AlertSystem       *AlertSystem
	SNVM              *VirtualMachine
}



// ConsensusSpecificNode represents a node that participates in the Synnergy Consensus using one specific mechanism.
type ConsensusSpecificNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	ConsensusMechanism ConsensusMechanism
	SubBlocks         map[string]*SubBlock
	mutex             sync.Mutex
	SyncInterval      time.Duration
	SNVM              *VirtualMachine
}

// ContentNode represents a specialized node designed to handle large data types linked to blockchain transactions.
type ContentNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	StorageManager    *StorageManager
	ContentCache      map[string]*Content
	mutex             sync.Mutex
	SyncInterval      time.Duration
	SNVM              *VirtualMachine
}

// CrossChainNode represents a node responsible for cross-chain interactions.
type CrossChainNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	SNVM              *VirtualMachine
	CrossChainManager *CrossChainManager
	Bridges           map[string]*Bridge
	mutex             sync.Mutex
	SyncInterval      time.Duration
}

// CustodialNode represents a node responsible for managing and safeguarding digital assets for users.
type CustodialNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	AssetManager      *AssetManager
	mutex             sync.Mutex
	SyncInterval      time.Duration
	Storage           *StorageManager
	ComplianceManager *ComplianceManager
	SNVM              *VirtualMachine
}

// DisasterRecoveryNode represents a node responsible for maintaining backups and ensuring recovery of the blockchain in case of failures.
type DisasterRecoveryNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	BackupManager     *BackupManager
	NetworkManager    *NetworkManager
	StorageManager    *StorageManager
	mutex             sync.Mutex
	BackupInterval    time.Duration
	RecoveryPlans     *RecoveryPlans
	SNVM              *VirtualMachine
}

// EnvironmentalMonitoringNode represents a node that integrates real-world environmental data with blockchain operations.
type EnvironmentalMonitoringNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	SensorManager     *SensorManager
	mutex             sync.Mutex
	SyncInterval      time.Duration
	SNVM              *VirtualMachine
}

// ForensicNode represents a node responsible for conducting forensic analysis on blockchain data.
type ForensicNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	ThreatDetection   *ThreatDetection
	ComplianceEngine  *ComplianceEngine
	mutex             sync.Mutex
	SyncInterval      time.Duration
	SNVM              *VirtualMachine
}

// ArchivalFullNode represents an archival full node in the blockchain network.
type ArchivalFullNode struct {
	NodeID             string
	CompleteBlockchain *Blockchain
	ConsensusEngine    *SynnergyConsensus
	EncryptionService  *Encryption
	NetworkManager     *NetworkManager
	SubBlocks          map[string]*SubBlock
	mutex              sync.Mutex
	SyncInterval       time.Duration
	SNVM               *VirtualMachine
}

// FullPrunedNode represents a full pruned node in the blockchain network.
type FullPrunedNode struct {
	NodeID            string
	PrunedBlockchain  *PrunedBlockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	SubBlocks         map[string]*SubBlock
	mutex             sync.Mutex
	SyncInterval      time.Duration
	SNVM              *VirtualMachine
}

// GatewayNode represents a node responsible for facilitating cross-chain interactions and integrating external data.
type GatewayNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	CrossChainManager *CrossChainManager
	Bridges           map[string]*Bridge
	mutex             sync.Mutex
	SyncInterval      time.Duration
	SNVM              *VirtualMachine
}

// GeospatialNode represents a node that processes and validates transactions involving geospatial data.
type GeospatialNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	GeoProcessor      *GeoProcessor
	SNVM              *VirtualMachine
	mutex             sync.Mutex
	SyncInterval      time.Duration
}

// HistoricalNode represents a node that stores and manages the entire blockchain history for audit and compliance purposes.
type HistoricalNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	SNVM              *VirtualMachine
	BackupManager     *DataBackupManager
	mutex             sync.Mutex
	SyncInterval      time.Duration
	BackupLocations   []string
}

// HolographicNode represents a node that distributes and stores blockchain data using holographic encoding.
type HolographicNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	SNVM              *VirtualMachine
	DataStorage       map[string]*HolographicData
	mutex             sync.Mutex
	SyncInterval      time.Duration
}

// HybridNode represents a versatile node that performs multiple functions within the blockchain, including validation, indexing, and transaction handling.
type HybridNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	SNVM              *VirtualMachine
	DataIndex         map[string]interface{}
	SubBlocks         map[string]*SubBlock
	mutex             sync.Mutex
	SyncInterval      time.Duration
}

// IndexingNode represents a node responsible for indexing blockchain data and handling queries efficiently.
type IndexingNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	SNVM              *VirtualMachine
	Index             map[string]interface{}
	SubBlocks         map[string]*SubBlock
	mutex             sync.Mutex
	SyncInterval      time.Duration
	QueryCache        map[string]interface{}
}

// IntegrationNode represents a node responsible for integrating external systems, APIs, and other blockchains.
type IntegrationNode struct {
	NodeID            string
	Blockchain        *Blockchain
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	SNVM              *VirtualMachine
	ChainAdapters     map[string]*ChainAdapter
	APIGateways       map[string]*APIGateway
	Oracles           map[string]*Oracle
	mutex             sync.Mutex
	SyncInterval      time.Duration
	NewFeatures       map[string]interface{}
}

// LightNode represents a lightweight node in the blockchain network.
type LightNode struct {
    NodeID            string                      // Unique identifier for the node
    BlockSummaries    map[string]*BlockSummary    // Stores simplified block summaries instead of full block headers
    ConsensusEngine   *SynnergyConsensus          // Consensus engine for validation and querying
    EncryptionService *Encryption                 // Encryption service for data encryption and decryption
    NetworkManager    *NetworkManager             // Network manager to communicate with full nodes
    FullNodes         []string                    // List of available full nodes for interaction
    SyncInterval      time.Duration               // Interval for syncing block summaries with full nodes
    mutex             sync.Mutex                  // Mutex to handle thread-safe operations
}


// LightningNode represents a node responsible for handling off-chain transactions using payment channels in the Lightning Network.
type LightningNode struct {
	NodeID            string
	PaymentChannels   map[string]*PaymentChannel
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	LiquidityPool     map[string]float64
	mutex             sync.Mutex
	SyncInterval      time.Duration
	FullNodes         []string
}


// PaymentChannel represents an active payment channel managed by the Lightning Node.
type PaymentChannel struct {
	ChannelID      string                      // Unique identifier for the payment channel
	ParticipantA   string                      // First participant in the channel
	ParticipantB   string                      // Second participant in the channel
	BalanceA       float64                     // Balance of Participant A in the channel
	BalanceB       float64                     // Balance of Participant B in the channel
	ChannelState   *ChannelState        // Current state of the payment channel
	ChannelTimeout time.Time                   // Expiry time for the channel
	IsActive       bool                        // Indicates whether the channel is active
}

// MobileNode represents a node designed to run on mobile devices (iOS, Android, etc.), providing a lightweight interface to the blockchain.
type MobileNode struct {
	NodeID            string
	PartialLedger     map[string]*Block
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	mutex             sync.Mutex
	SyncInterval      time.Duration
	FullNodes         []string
	BatteryOptimized  bool
}

// OptimizationNode represents a node responsible for optimizing transaction processing and improving network performance.
type OptimizationNode struct {
	NodeID            string
	PartialLedger     map[string]*Block
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	mutex             sync.Mutex
	SyncInterval      time.Duration
	FullNodes         []string
	OptimizationAlgo  OptimizationAlgorithm
}

// OptimizationAlgorithm defines the structure of the algorithm used for dynamic transaction optimization.
type OptimizationAlgorithm struct {
	PriorityModel   string // Defines the model used for prioritizing transactions (e.g., FIFO, fee-based priority)
	LoadBalancer    string // Defines the load-balancing strategy (e.g., round-robin, weighted distribution)
}

// OrphanNode represents a node responsible for handling orphan blocks in the blockchain.
type OrphanNode struct {
	NodeID            string
	OrphanBlocks      map[string]*Block
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	mutex             sync.Mutex
	SyncInterval      time.Duration
	FullNodes         []string
	ArchivedBlocks    map[string]*Block
}

// StoredFile represents the structure of a file stored on the storage node.
type StoredFile struct {
	FileID          string
	FileName        string
	OwnerWallet     string
	EncryptedData   []byte
	StoredAt        time.Time
	FileSize        int64
	IPFSHash        string
	SwarmHash       string
}

// StorageNode represents a storage node responsible for securely storing files and data.
type StorageNode struct {
	NodeID            string
	StorageCapacity   int64
	UsedStorage       int64
	StoredFiles       map[string]*StoredFile
	mutex             sync.Mutex
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	Ledger            *Ledger
	SNVM              *VirtualMachine
	IPFSService       *IPFSManager
	SwarmService      *SwarmManager
	CacheService      *CacheManager
}

// SuperNode represents a high-capacity node responsible for handling transaction routing, data storage, executing smart contracts, and privacy features.
type SuperNode struct {
	NodeID            string
	StorageCapacity   int64
	UsedStorage       int64
	StoredBlocks      map[string]*Block
	mutex             sync.Mutex
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	Ledger            *Ledger
	SNVM              *VirtualMachine
}

// TestnetNode represents a node designed specifically to run in the testnet environment. It validates test transactions and sub-blocks exclusively in the testnet.
type TestnetNode struct {
	NodeID            string
	PartialLedger     map[string]*Block
	TestBlocks        map[string]*SubBlock
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	Ledger            *Ledger
	SyncInterval      time.Duration
	FullNodes         []string
	mutex             sync.Mutex
}

// TimeLockedNode represents a node responsible for managing time-locked transactions and contracts.
type TimeLockedNode struct {
	NodeID            string
	PartialLedger     map[string]*Block
	TimeLocks         map[string]*Timelock
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	Ledger            *Ledger
	SNVM              *VirtualMachine
	mutex             sync.Mutex
	SyncInterval      time.Duration
	FullNodes         []string
}

// ValidatorNode represents a node responsible for validating transactions into sub-blocks but does not mine full blocks.
type ValidatorNode struct {
	NodeID            string
	SubBlocks         map[string]*SubBlock
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	mutex             sync.Mutex
	SyncInterval      time.Duration
	SNVM              *VirtualMachine
	ConsensusType     string
}

// WatchtowerNode represents a node responsible for overseeing transactions and ensuring compliance with smart contracts, particularly for off-chain transactions.
type WatchtowerNode struct {
	NodeID            string
	MonitoredChannels map[string]*Channel
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	mutex             sync.Mutex
	SNVM              *VirtualMachine
	SyncInterval      time.Duration
	Logs              map[string]string
}

// ZKPNode represents a Zero-Knowledge Proof node that handles privacy-preserving transactions and verifies proofs without revealing any sensitive data.
type ZKPNode struct {
	NodeID            string
	ZKProofs          map[string]*ZkProof
	ConsensusEngine   *SynnergyConsensus
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	mutex             sync.Mutex
	SyncInterval      time.Duration
	SNVM              *VirtualMachine
}


// AuthorityNode represents an authority node with unique permissions
type AuthorityNodeVersion struct {
	NodeID            string    // Unique ID for the authority node
	SecretKey         string    // Secret key for node access control
	CreatedAt         time.Time // Timestamp of node creation
	EncryptedKey      string    // Encrypted form of the secret key
	AuthorityNodeType AuthorityNodeTypes  // Assuming 'nodeType' is a valid type or enum defined elsewhere
}


// Permission defines the structure for specific permissions granted to an authority node
type Permission struct {
	PermissionID string    // Unique ID for the permission
	Description  string    // Description of the permission granted
	GrantedAt    time.Time // Time when the permission was granted
	GrantedBy    string    // ID of the entity granting the permission
}

// AuthorityNodeType defines different types of authority nodes in the network
type AuthorityNodeTypes string

// Enumeration of Authority Node Types
const (
	ElectedAuthorityNodeType AuthorityNodeTypes = "Elected Authority Node"
	MilitaryNodeType         AuthorityNodeTypes = "Military Node"
	BankingNodeType          AuthorityNodeTypes = "Banking Node"
	CentralBankNodeType      AuthorityNodeTypes = "Central Bank Node"
	ExchangeNodeType         AuthorityNodeTypes = "Exchange Node"
	GovernmentNodeType       AuthorityNodeTypes = "Government Node"
	RegulatorNodeType        AuthorityNodeTypes = "Regulator Node"
)


// AuthorityNodePermissions define what each authority node type is allowed to do.
type AuthorityNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
	VerifySyn900Deployment    bool
}

// AuthorityNode represents a full node with extended permissions and functionality for authority nodes.
type AuthorityNode struct {
	NodeID            string
	KeyManager        *KeyManager
	Ledger            *Ledger
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	RequestList       map[string]*Request
	Permissions       AuthorityNodePermissions
	mutex             sync.Mutex
}

// BankNodePermissions define what each bank node is allowed to do.
type BankNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
	DisburseLoans             bool
	FreezeWallets             bool
	VerifySyn900Deployment    bool
}

// BankNode represents a bank node with extended permissions and functionality specific to banking operations.
type BankNode struct {
	NodeID            string
	KeyManager        *KeyManager
	Ledger            *Ledger
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	LoanPool          *LoanPool
	RequestList       map[string]*Request
	Permissions       BankNodePermissions
	WalletFreezeList  map[string]bool
	mutex             sync.Mutex
}

// CentralBankNodePermissions defines what each central bank node is allowed to do.
type CentralBankNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
	DisburseLoans             bool
	FreezeWallets             bool
	VerifySyn900Deployment    bool
}

// CentralBankNode represents a central bank node with extended permissions and functionalities specific to central bank operations.
type CentralBankNode struct {
	NodeID            string
	KeyManager        *KeyManager
	Ledger            *Ledger
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	LoanPool          *LoanPool
	RequestList       map[string]*Request
	Permissions       CentralBankNodePermissions
	WalletFreezeList  map[string]bool
	mutex             sync.Mutex
}

// CreditProviderNodePermissions defines what each credit provider node is allowed to do.
type CreditProviderNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
	DisburseLoans             bool
	VerifySyn900Deployment    bool
}

// CreditProviderNode represents a credit provider node with extended permissions and functionalities.
type CreditProviderNode struct {
	NodeID            string
	KeyManager        *KeyManager
	Ledger            *Ledger
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	LoanPool          *LoanPool
	RequestList       map[string]*Request
	Permissions       CreditProviderNodePermissions
	mutex             sync.Mutex
}

// ElectedAuthorityNodePermissions defines what each elected authority node can do.
type ElectedAuthorityNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
}

// ElectedAuthorityNode represents an elected authority node with specific permissions and functionalities.
type ElectedAuthorityNode struct {
	NodeID            string
	KeyManager        *KeyManager
	Ledger            *Ledger
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	RequestList       map[string]*Request
	Permissions       ElectedAuthorityNodePermissions
	mutex             sync.Mutex
}

// ExchangeNodePermissions defines what each exchange node can do.
type ExchangeNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
	VerifySyn900ID            bool
}

// ExchangeNode represents an exchange node with specific permissions and functionalities.
type ExchangeNode struct {
	NodeID            string
	KeyManager        *KeyManager
	Ledger            *Ledger
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	RequestList       map[string]*Request
	Permissions       ExchangeNodePermissions
	Syn900Verifier    *tokenledgers.Syn900Verifier
	mutex             sync.Mutex
}

// GovernmentNodePermissions defines the immutable permissions for a government node.
type GovernmentNodePermissions struct {
	ViewRequestList            bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal    bool
	ReportAuthorityNode        bool
	ViewPrivateTransactions    bool
	DisburseLoans              bool
	FreezeWallets              bool
	ExecuteComplianceContracts bool
	AddCompliance              bool
	RemoveCompliance           bool
	VerifySyn900ID             bool
}

// GovernmentNode represents a government node with specific permissions and functionalities.
type GovernmentNode struct {
	NodeID            string
	KeyManager        *KeyManager
	Ledger            *Ledger
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	RequestList       map[string]*Request
	Permissions       GovernmentNodePermissions
	Syn900Verifier    *tokenledgers.Syn900Verifier
	mutex             sync.Mutex
}

// MilitaryNodePermissions defines the immutable permissions for a military node.
type MilitaryNodePermissions struct {
	ViewRequestList            bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal    bool
	ReportAuthorityNode        bool
	ViewPrivateTransactions    bool
	VerifySyn900ID             bool
}

// MilitaryNode represents a military node with specific permissions and functionalities.
type MilitaryNode struct {
	NodeID            string
	KeyManager        *KeyManager
	Ledger            *Ledger
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	RequestList       map[string]*Request
	Permissions       MilitaryNodePermissions
	Syn900Verifier    *tokenledgers.Syn900Verifier
	mutex             sync.Mutex
}

// RegulatorNodePermissions defines the immutable permissions for a regulator node.
type RegulatorNodePermissions struct {
	ViewRequestList            bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal    bool
	ReportAuthorityNode        bool
	ViewPrivateTransactions    bool
	ViewAndDisburseLoans       bool
	FreezeWallets              bool
	ExecuteComplianceContracts bool
	AddCompliance              bool
	RemoveCompliance           bool
	VerifySyn900ID             bool
}

// RegulatorNode represents a regulator node with specific permissions and functionalities.
type RegulatorNode struct {
	NodeID            string
	KeyManager        *KeyManager
	Ledger            *Ledger
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	RequestList       map[string]*Request
	Permissions       RegulatorNodePermissions
	Syn900Verifier    *tokenledgers.Syn900Verifier
	mutex             sync.Mutex
}

// VotingResult represents the result of the voting process.
type VotingResult struct {
	ConfirmedVotes int
	RejectedVotes  int
	Status         string // "Accepted" or "Rejected"
}

// AuthorityNodeAcceptanceManager manages the process for voting on authority node proposals.
type AuthorityNodeAcceptanceManager struct {
	mutex             sync.Mutex
	Ledger            *Ledger
	EncryptionService *Encryption
	NetworkManager    *NetworkManager
	IdentityVerifier  *tokenledgers.Syn900Verifier
	KeyDisburser      *KeyDisburser
}

// AuthorityNodeKey represents a generated key for a specific node type.
type AuthorityNodeKey struct {
	KeyID          string
	NodeType       string
	OwnerName      string
	OwnerWallet    string
	ExpirationDate time.Time
	MaxNodes       int
	UsedNodes      int
	IsExpired      bool
}

// AuthorityNodeAccountManager manages authority node key ownership, wallet registration, key cancellation, and refreshing.
type AuthorityNodeAccountManager struct {
	mutex             sync.Mutex
	Ledger            *Ledger
	EncryptionService *Encryption
	AuthorityKeys     map[string]*AuthorityNodeKey
}

// AuthorityNodeReport represents a report filed against an authority node.
type AuthorityNodeReport struct {
	NodeID           string
	ReportType       string
	ReporterID       string
	ReportTimestamp  time.Time
}

// AuthorityNodeReportManager manages the reporting system for authority nodes.
type AuthorityNodeReportManager struct {
	mutex             sync.Mutex
	Ledger            *Ledger
	EncryptionService *Encryption
	Reports           map[string][]*AuthorityNodeReport
}

// AuthorityNodeRewardPool manages and distributes the reward pool for authority nodes.
type AuthorityNodeRewardPool struct {
	TotalPoolAmount    float64
	AuthorityNodes     map[string]*AuthorityNode
	Ledger             *Ledger
	NetworkManager     *NetworkManager
	mutex              sync.Mutex
	DistributionPeriod time.Duration
}

// KeyManager manages authority node key creation, usage, and expiration.
type KeyManager struct {
	mutex             sync.Mutex
	Ledger            *Ledger
	EncryptionService *Encryption
	GeneratedKeys     map[string]*AuthorityNodeKey
	GenesisCreated    bool
	KeyDisburser      *KeyDisburser
}


// KeyDisburser manages the distribution, tracking, and revocation of keys used by nodes in the network.
type KeyDisburser struct {
    DisbursedKeys   map[string]*AuthorityNodeKey // Map of keys that have been disbursed, by their key ID
    TotalKeysIssued int                          // Total number of keys that have been issued by the system
    MaxKeysAllowed  int                          // Maximum number of keys that can be disbursed in total
    RevokedKeys     map[string]*AuthorityNodeKey // Map of keys that have been revoked, by their key ID
    Ledger          *Ledger                      // Reference to the ledger for recording disbursement and revocation
    mutex           sync.Mutex                   // Mutex for thread-safe operations
}

// Request represents a generic request made to a node or authority for processing transactions, data queries, or system actions.
type Request struct {
    RequestID     string    // Unique identifier for the request
    RequestType   string    // Type of request (e.g., "Transaction", "DataQuery", "Cancellation", "OracleUpdate")
    RequestorID   string    // ID of the entity making the request (could be a wallet address, node ID, or authority ID)
    RequestedData string    // Data or resource being requested
    Timestamp     time.Time // Timestamp when the request was made
    Status        string    // Status of the request (e.g., "Pending", "Approved", "Rejected")
    Priority      int       // Priority level of the request (e.g., 1 = high, 10 = low)
    AdditionalData string   // Optional field for any additional metadata or details about the request
}


