package network

import (
	"container/list"
	"crypto/rsa"
	"encoding/json"
	"net"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// Client represents an RPC client for interacting with the Synnergy Network
type RPCClient struct {
	ServerAddress string
	PublicKey     string
	PrivateKey    string
}

// ConnectionEvent represents an event related to wallet connections.
type ConnectionEvent struct {
	EventID        string    // Unique identifier for the event
	ConnectionID   string
	WalletID       string    // Wallet associated with the event
	EventType      string    // Type of event (e.g., "connection", "disconnection")
	EventTime      time.Time // Timestamp of the event
	Details        string    // Additional details related to the event
}


// NetworkRequest represents an incoming network request
type NetworkRequest struct {
	Type string `json:"type"` // The type of the request (e.g., "TRANSACTION_SUBMISSION")
	Data []byte `json:"data"` // The raw data (e.g., serialized transaction)
	Payload json.RawMessage `json:"payload"` // The raw payload for further processing

}


// NetworkResponse represents the structure of an outgoing network response
type NetworkResponse struct {
    Type    string `json:"type"`    // The type of the response (e.g., "PONG", "ERROR")
    Message string `json:"message"` // The message or payload returned in the response
}




// FaultToleranceManager manages the fault tolerance mechanisms in the network
type FaultToleranceManager struct {
    Nodes             []string        // List of nodes in the network
    NodeStatus        map[string]bool // Status of each node (alive or down)
    QuorumThreshold   int             // Minimum number of nodes required for consensus
    ledger            *ledger.Ledger  // Ledger for logging fault tolerance events
    mutex             sync.Mutex      // Mutex for thread-safe operations
	NodeState     map[string]string        // Map of node to its latest synced block hash
}

// Firewall manages and filters network traffic, securing the network from malicious actors
type Firewall struct {
	allowedIPs      map[string]bool         // A whitelist of allowed IPs
	blockedIPs      map[string]time.Time    // A list of blocked IPs and the time they were blocked
	blockDuration   time.Duration           // Duration for which an IP will remain blocked
	mutex           sync.Mutex              // Mutex for thread-safe operations
	ledgerInstance  *ledger.Ledger          // Pointer to the ledger for logging firewall events
}

// FirewallManager manages multiple firewall types including normal, dynamic, stateless, and stateful firewalls
type FirewallManager struct {
	NormalFirewall    *Firewall        // Normal firewall (whitelisting/blocking)
	DynamicFirewall   *DynamicFirewall // Dynamic firewall (adjusts based on traffic patterns)
	StatelessFirewall *StatelessFirewall // Stateless firewall (simple packet-filtering firewall)
	StatefulFirewall  *StatefulFirewall  // Stateful firewall (tracks connection state)
	ledgerInstance    *ledger.Ledger
}

// DynamicFirewall adjusts its rules dynamically based on traffic patterns
type DynamicFirewall struct {
	allowedIPs     map[string]bool
	blockedIPs     map[string]time.Time
	blockDuration  time.Duration
	ledgerInstance *ledger.Ledger
	encryptionService *common.Encryption
	mutex          sync.Mutex
}

// StatelessFirewall represents a simple firewall that filters based on rules without tracking connection states
type StatelessFirewall struct {
	allowedPorts   []int
	blockedIPs     map[string]time.Time
	ledgerInstance *ledger.Ledger
	encryptionService *common.Encryption
	mutex          sync.Mutex
}

// StatefulFirewall keeps track of connections and the state of each connection (open, closed, etc.)
type StatefulFirewall struct {
	allowedConnections map[string]string // IP -> State (open/closed)
	blockedIPs         map[string]time.Time
	ledgerInstance     *ledger.Ledger
	encryptionService  *common.Encryption
	mutex              sync.Mutex
}

// FlowControlManager handles network traffic, ensuring smooth flow of transactions and blocks
type FlowControlManager struct {
	MaxPendingTransactions int                    // Max number of pending transactions allowed
	MaxSubBlockSize        int                    // Max number of transactions per sub-block
	MaxBlockSize           int                    // Max number of sub-blocks per block
	PendingTransactions    []common.Transaction   // Pool of pending transactions
	PendingSubBlocks       []common.SubBlock      // Pool of pending sub-blocks waiting to be added to a block
	mutex                  sync.Mutex             // Mutex for thread-safe operations
	ledgerInstance         *ledger.Ledger         // Pointer to the ledger for tracking flow control events
}

// GeoLocation represents the latitude and longitude of a node
type GeoLocation struct {
    Latitude  float64
    Longitude float64
}

// GeoLocationManager manages the geolocation data for nodes in the network
type GeoLocationManager struct {
    NodeLocations   map[string]GeoLocation // Maps node IDs to their geolocation
    ledger          *ledger.Ledger         // Reference to the ledger for logging
}

// Handshake represents the handshake mechanism for secure communication between nodes
type Handshake struct {
    PrivateKey *rsa.PrivateKey       // The private key of the node
    PublicKey  *rsa.PublicKey        // The public key of the node
    mutex      sync.Mutex            // Mutex for thread-safe operations
    ledger     *ledger.Ledger        // Pointer to the ledger for logging handshake events
}

// KademliaNode represents a node in the Kademlia DHT
type KademliaNode struct {
    ID         string          // Unique ID of the node
    IPAddress  string          // Node's IP address
    LastActive time.Time       // Last active timestamp
    Location   GeoLocation // Geolocation of the node
}

// KademliaDHT represents the Kademlia Distributed Hash Table
type KademliaDHT struct {
    NodeID         string                // ID of the local node
    KBucket        map[string]KademliaNode // Kademlia node bucket (ID to node mapping)
    ledger         *ledger.Ledger        // Reference to ledger for logging
    lock           sync.Mutex            // Mutex for thread-safe operations
}

// MessageQueue represents the queue of messages waiting to be processed
type MessageQueue struct {
	queue          *list.List        // Doubly-linked list for efficient message queuing
	maxQueueSize   int               // Maximum size of the message queue
	lock           sync.Mutex        // Mutex to ensure thread-safe operations
	ledgerInstance *ledger.Ledger    // Reference to the ledger for logging message events
}

// Message represents a message between nodes in the network
type Message struct {
	ID        string             // Unique message ID
	Timestamp time.Time          // Time when the message was created
	From      string             // Node ID of the sender
	To        string             // Node ID of the receiver
	Content   string             // Message content
	Hash      string             // SHA-256 hash of the message content
	Encrypted bool               // Whether the message content is encrypted
}

// NATTraversalManager handles NAT traversal tasks for nodes behind routers
type NATTraversalManager struct {
	publicIP       string             // Public IP address of the node
	privateIP      string             // Private IP address of the node
	peerMap        map[string]string  // Map of peers' public IPs to private IPs
	lock           sync.Mutex         // Mutex to ensure thread-safe operations
	ledgerInstance *ledger.Ledger     // Reference to the ledger for logging NAT traversal events
	connectionPool  *ConnectionPool            // Pool to manage active connections
	webrtcManager *WebRTCManager
}

// NetworkManager handles all networking activities, including peer discovery, message passing, and encryption
type NetworkManager struct {
	nodeAddress     string
	peers           map[string]*PeerConnection // Active peers connected to this node
	ledgerInstance  *ledger.Ledger             // Reference to the ledger for logging network events
	lock            sync.Mutex                 // Mutex for thread-safe operations
	connectionPool  *ConnectionPool            // Pool to manage active connections
	PeerDiscovery   *PeerDiscoveryManager 
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

// BlockRequest represents a request for retrieving a block by index
type BlockRequest struct {
	BlockIndex int `json:"block_index"`
}


// Packet represents a network packet that needs to be routed
type Packet struct {
    SourceID      string // ID of the source peer
    DestinationID string // ID of the destination peer
    Data          []byte // Data to be sent in the packet
}

// RPCRequest represents an RPC call request
type RPCRequest struct {
	Method         string `json:"method"`
	Payload        string `json:"payload"`
	SenderPublicKey string `json:"sender_public_key"`
}

// RPCResponse represents the response from an RPC call
type RPCResponse struct {
	Data string `json:"data"`
}

// NodeType defines the possible types of nodes in the network
type NodeType string

const (
    DefaultNodeType   NodeType = "default"   // Default node type (general-purpose node)
    ValidatorNodeType NodeType = "validator" // Node that participates in PoS validation
    MinerNodeType     NodeType = "miner"     // Node that performs PoW mining
)




// PeerKey represents the public key of a peer
type PeerKey []byte


// PeerDiscoveryManager manages the discovery of new peers
type PeerDiscoveryManager struct {
	Peers         map[string]*Peer       // List of discovered peers
	Ledger        *ledger.Ledger         // Instance of the ledger for validation
	LocalAddress  string                 // Local node's address
	LocalPublicKey string                // Local node's public key
	LocalPrivateKey  string  // Add this field for the private key
	KnownNetworks []string               // List of known networks for discovering peers
}

// QoSManager manages the quality of service across the P2P network
type QoSManager struct {
	BandwidthLimit  int                 // Maximum bandwidth allowed for communication in KBps
	PriorityQueues  map[int][]*QoSPacket // Queues for different packet priorities
	Mutex           sync.Mutex          // Mutex to handle concurrency
	LedgerInstance  *ledger.Ledger      // Instance of the ledger to record QoS-related events
	NodePublicKey   string              // Public key of the local node for encryption
	ConnectedPeers  map[string]*Peer    // List of connected peers
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
	Routes           map[string]string   // Map of node IDs to their IP addresses
	Mutex            sync.Mutex          // Mutex for thread-safe operations
	LedgerInstance   *ledger.Ledger      // Ledger instance to store routing events
	EncryptedRoutes  map[string][]byte   // Encrypted routes for added security
	NodePublicKey    string              // Public key of the local node for encryption
	Peers            map[string]*Peer    // List of active peers with their addresses and public keys
	ConnectionPool   *ConnectionPool     // Connection pool for managing network connections
}

// RPCServer represents the RPC server for Synnergy Network
type RPCServer struct {
	Address        string
	LedgerInstance *ledger.Ledger // Ledger integration for storing and retrieving data
	Router         *Router        // Router to manage routing and peer connections
}

// SDNController represents the Software-Defined Network controller, managing the network nodes
type SDNController struct {
	Nodes          map[string]*SDNNode  // Connected nodes in the network
	NodeLock       sync.Mutex           // Mutex for thread-safe node operations
	LedgerInstance *ledger.Ledger       // Integration with ledger for network state tracking
	EncryptionKey  string               // Key for encrypting SDN controller communications
}

// SDNNode represents a node in the SDN-controlled network
type SDNNode struct {
	NodeID    string         // Unique ID of the node
	Address   net.IP         // IP address of the node
	Status    string         // Status (active, inactive)
	LastCheck time.Time      // Last heartbeat or health check
}

// Server represents a blockchain server node responsible for handling incoming connections
type Server struct {
	Address          string                // Server address
	LedgerInstance   *ledger.Ledger        // Ledger for storing blocks and transactions
	ConsensusEngine  *common.SynnergyConsensus // Consensus mechanism
	EncryptionKey    string                // Encryption key for secure communication
	connections      map[string]net.Conn   // Map of active connections
	connectionLock   sync.Mutex            // Mutex for thread-safe operations
}

// SSLHandshakeManager handles SSL/TLS secure handshakes for encrypted communication
type SSLHandshakeManager struct {
	CertFile   string        // Path to the certificate file
	KeyFile    string        // Path to the private key file
	CAFile     string        // Path to the Certificate Authority (CA) file
	Connections map[string]net.Conn // Map of active connections
	lock       sync.Mutex    // Mutex to ensure thread-safe operations
}

// TLSHandshakeManager manages secure TLS handshakes for encrypted communication between nodes
type TLSHandshakeManager struct {
	CertFile     string               // Path to the TLS certificate file
	KeyFile      string               // Path to the private key file
	CAFile       string               // Path to the CA certificate
	Connections  map[string]net.Conn  // Active TLS connections
	lock         sync.Mutex           // Mutex to protect concurrent access
	Ledger       *ledger.Ledger       // Integration with the ledger for secure transaction handling
}



// TopologyManager manages the network topology and routing for communication between nodes
type TopologyManager struct {
	Nodes          map[string]*common.Node     // A map of all active nodes in the network
	NetworkGraph   *NetworkGraph // The network graph for representing connections between nodes
	lock           sync.Mutex           // Mutex for thread-safe operations
	LedgerInstance *ledger.Ledger       // Integration with the ledger for node data persistence
}

// WebRTCManager handles WebRTC connections for peer-to-peer communication across the Synnergy Network
type WebRTCManager struct {
	Peers         map[string]*PeerConnection // Map of active peer connections
	lock          sync.Mutex                 // Mutex for thread-safe operations
	LedgerInstance *ledger.Ledger            // Ledger for tracking peer connections
}

// EncryptedMessage represents an encrypted payload to be sent between nodes.
type EncryptedMessage struct {
    CipherText  []byte    // The encrypted message content
    Hash        string    // Hash of the original message for integrity check
    CreatedAt   time.Time // When the message was encrypted (Renamed from Timestamp to CreatedAt)
}


// GraphNode represents a node in the network graph.
type GraphNode struct {
	NodeID      string              // Unique identifier of the node
	NodeInfo    *NodeInfo           // Information about the node
	Connections []*GraphEdge        // Edges representing connections to other nodes
}

// GraphEdge represents a connection between two nodes in the network graph.
type GraphEdge struct {
	FromNodeID string    // ID of the source node
	ToNodeID   string    // ID of the destination node
	Weight     float64   // Weight or cost of the connection (can represent latency, bandwidth, etc.)
}

// NodeKey represents the public-private key pair for a node, used for encryption and identity verification.
type NodeKey struct {
	PublicKey  string    // The public key of the node
	PrivateKey string    // The private key of the node (kept secret)
}

// NodeInfo represents the information related to a network node.
type NodeInfo struct {
	NodeID         string       // Unique identifier of the node
	Address 	   string
	IPAddress      string       // IP address of the node
	Port           int          // Port number the node is listening on
	NodeType       NodeType     // Type of the node (default, validator, miner)
	GeoLocation    GeoLocation  // Geographical location of the node (latitude and longitude)
	LastActiveTime time.Time    // Last active timestamp of the node
	IsOnline       bool         // Status indicating if the node is currently online
}

// NodeHealthStatus represents the health status of a node in the network
type NodeHealthStatus string

const (
    Healthy    NodeHealthStatus = "healthy"    // Node is operating normally
    Unhealthy  NodeHealthStatus = "unhealthy"  // Node is experiencing issues or errors
    Degraded   NodeHealthStatus = "degraded"   // Node is operational but performing below expectations
    Offline    NodeHealthStatus = "offline"    // Node is not reachable or non-functional
)

// NodeState represents the operational state of a node in the network
type NodeState string

const (
    Active     NodeState = "active"      // Node is actively participating in the network
    Inactive   NodeState = "inactive"    // Node is currently inactive or idle
    Syncing    NodeState = "syncing"     // Node is synchronizing data with the network
    Validating NodeState = "validating"  // Node is validating blocks or transactions
    Mining     NodeState = "mining"      // Node is mining (for PoW nodes)
    Staking    NodeState = "staking"     // Node is staking (for PoS nodes)
    ShuttingDown NodeState = "shutting_down" // Node is in the process of shutting down
)

// NodeMetrics represents performance and resource usage metrics for a node in the network.
type NodeMetrics struct {
	NodeID             string    // Unique identifier for the node
	Uptime             time.Duration // Uptime of the node
	TotalTransactions  int       // Total number of transactions processed by the node
	BlockProcessingTime time.Duration // Average time taken to process a block
	SubBlockProcessingTime time.Duration // Average time taken to process a sub-block (if applicable)
	MemoryUsage        int64     // Current memory usage of the node (in bytes)
	CPUUsage           float64   // CPU usage percentage
	NetworkLatency     time.Duration // Average network latency for the node
	TotalBlocks        int       // Total number of blocks processed by the node
	ErrorsEncountered  int       // Total number of errors encountered by the node
	LastUpdated        time.Time // The last time the metrics were updated
	PeerCount          int       // Number of peers connected to the node
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


// Define the WebRTCConn struct with an EncryptionKey field
type WebRTCConn struct {
    ConnectionID   string
    EncryptionKey  []byte // Added EncryptionKey field
}

// Define WebRTCConnection struct in the ledger package (if needed)
type WebRTCConnection struct {
    ConnectionID string
    PeerID       string
    Timestamp    time.Time // Renamed from "ConnectedAt" to "Timestamp"
}

