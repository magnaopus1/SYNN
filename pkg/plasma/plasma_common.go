package common

import (
	"sync"
	"time"

	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/encryption"
)

// PlasmaClient represents a client in the Plasma childchain network
type PlasmaClient struct {
	ClientID          string                   // Unique identifier for the client
	WalletAddress     string                   // Wallet address associated with the client
	Ledger            *ledger.Ledger           // Reference to the ledger for transaction recording
	EncryptionService *encryption.Encryption   // Encryption service for securing client transactions
	PlasmaChain       *PlasmaChain             // Reference to the Plasma childchain
	NetworkManager    *NetworkManager  // Network manager for interacting with nodes and the childchain
}

// PlasmaCore represents the core logic of the Plasma childchain
type PlasmaCore struct {
	Blocks           map[string]*PlasmaBlock    // Collection of all blocks in the childchain
	SubBlocks        map[string]*PlasmaSubBlock // Collection of sub-blocks within the childchain
	Ledger           *ledger.Ledger             // Reference to the ledger for recording chain events
	EncryptionService *encryption.Encryption    // Encryption service to secure transactions
	mu               sync.Mutex                 // Mutex for handling concurrency
}

// PlasmaBlock represents a block in the Plasma childchain
type PlasmaBlock struct {
	BlockID        string                      // Unique identifier for the block
	PreviousBlock  string                      // ID of the previous block in the chain
	SubBlocks      map[string]*PlasmaSubBlock  // Collection of sub-blocks in the block
	Timestamp      time.Time                   // Timestamp of when the block was created
	ValidatorID    string                      // Validator responsible for the block
}

// PlasmaSubBlock represents a sub-block in the Plasma childchain
type PlasmaSubBlock struct {
	SubBlockID     string                      // Unique identifier for the sub-block
	ParentBlockID  string                      // ID of the parent block
	Transactions   []Transaction        // List of transactions included in the sub-block
	Timestamp      time.Time                   // Timestamp of when the sub-block was created
	ValidatorID    string                      // Validator responsible for the sub-block
}

// PlasmaChain represents a Plasma childchain, containing blocks and managing the network
type PlasmaChain struct {
	ChainID        string                      // Unique identifier for the Plasma childchain
	GenesisBlock   *PlasmaBlock                // The genesis block of the chain
	CurrentBlockID string                      // ID of the current block
	Blocks         map[string]*PlasmaBlock     // All blocks in the chain
	Ledger         *ledger.Ledger              // Ledger for recording events
	Core           *PlasmaCore                 // Core logic for sub-block and block validation
	Encryption     *encryption.Encryption      // Encryption service for securing the chain
	NetworkManager *NetworkManager     // Network manager for handling nodes and chain state
}

// PlasmaChainConfig represents the configuration for creating a Plasma childchain
type PlasmaChainConfig struct {
	ChainID           string                    // Unique ID for the Plasma childchain
	GenesisBlockID    string                    // ID of the genesis block
	GenesisTimestamp  time.Time                 // Timestamp for the genesis block
	PlasmaCore        *PlasmaCore               // Reference to the Plasma core logic
	Ledger            *ledger.Ledger            // Reference to the ledger for recording events
	NetworkManager    *NetworkManager   // Network manager for handling nodes
	EncryptionService *encryption.Encryption    // Encryption service for securing the chain
}

// PlasmaCrossChain represents the logic for handling cross-chain operations
type PlasmaCrossChain struct {
	Blocks          map[string]*PlasmaBlock     // Collection of all blocks in the childchain
	SubBlocks       map[string]*PlasmaSubBlock  // Collection of sub-blocks within the childchain
	Ledger          *ledger.Ledger              // Reference to the ledger for recording chain events
	Encryption      *encryption.Encryption      // Encryption service to secure transactions
	CrossChainComm  *CrossChainCommunication  // Cross-chain communication logic
	mu              sync.Mutex                  // Mutex for concurrency handling
}

// PlasmaNetwork represents the Plasma childchain network operations
type PlasmaNetwork struct {
	PlasmaNodes     map[string]*PlasmaNode // Collection of all participating Plasma nodes
	Blocks          map[string]*PlasmaBlock        // Collection of all blocks in the Plasma childchain
	SubBlocks       map[string]*PlasmaSubBlock     // Collection of all sub-blocks in the Plasma childchain
	Ledger          *ledger.Ledger                 // Reference to the ledger for recording chain events
	Encryption      *encryption.Encryption         // Encryption service to secure transactions
	NetworkManager  *NetworkManager               // Network manager to handle node communications
	mu              sync.Mutex                     // Mutex for concurrency handling
}

// PlasmaNode represents a node in the Plasma childchain network
type PlasmaNode struct {
	NodeID          string                      // Unique identifier for the node
	IPAddress       string                      // Node's IP address
	Encryption      *encryption.Encryption      // Encryption service used by the node
	LastActiveTime  time.Time                   // Timestamp of the node's last activity
	NodeHealth      NodeHealthStatus    // Current health status of the node
	NodeState       NodeState           // Node's current operational state
}
