package storage

import (
    "sync"
    "time"
    "synnergy_network/pkg/ledger"
)

// CacheEntry represents an entry in the cache
type CacheEntry struct {
    Key          string    // Cache key
    Data         []byte    // Encrypted cached data
    CachedAt     time.Time // The time when the entry was cached
    ExpiresAt    time.Time // Expiration time for the cache entry
    LastAccessed time.Time // Last accessed time
}



// CacheManager handles the caching mechanism for blockchain operations
type CacheManager struct {
    CacheEntries   map[string]*CacheEntry // Cache entries stored by key
    mutex          sync.Mutex             // Mutex for thread-safe operations
    LedgerInstance *ledger.Ledger         // Ledger for recording cache activity
    MaxCacheSize   int                    // Maximum number of cache entries allowed
    DefaultTTL     time.Duration          // Time-to-live for cache entries
}

// TimestampEntry represents a timestamped data entry
type TimestampEntry struct {
    DataHash  string
    Timestamp time.Time
    OwnerID   string
}


// FileEntry represents a file stored in the system
type FileEntry struct {
    FileName      string    // The original file name
    FilePath      string    // The full path of the file in the storage system
    Encrypted     bool      // Whether the file is encrypted
    UploadedAt    time.Time // The timestamp when the file was uploaded
    LastAccessed  time.Time // The timestamp of the last access
}

// FileManager manages file storage, retrieval, and encryption
type FileManager struct {
    Files          map[string]*FileEntry // A map of file hashes to FileEntry
    storageDir     string                // Directory where files are stored
    mutex          sync.Mutex            // Mutex for thread-safe operations
    LedgerInstance *ledger.Ledger        // Ledger for tracking file operations
}

// IPFSManager handles file storage and retrieval from IPFS with encryption
type IPFSManager struct {
    mutex          sync.Mutex         // Mutex for thread-safe operations
    LedgerInstance *ledger.Ledger     // Ledger for tracking IPFS file operations
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
    mutex          sync.Mutex                 // Mutex for thread-safe operations
    Indexes        map[string]*FileIndex      // Indexed file metadata by FileID
    LedgerInstance *ledger.Ledger             // Ledger instance for tracking file operations
}

// EscrowAccount represents an escrow account to hold funds temporarily during a transaction
type EscrowAccount struct {
    EscrowID      string    // Unique escrow account ID
    Buyer         string    // Buyer involved in the transaction
    Seller        string    // Seller involved in the transaction
    Amount        float64   // Amount of funds in escrow
    CreatedAt     time.Time // Timestamp of when the escrow account was created
    IsReleased    bool      // Whether the funds have been released
}

// StorageListing represents a storage listing on the marketplace
type StorageListing struct {
    ListingID      string    // Unique ID of the listing
    Owner          string    // Owner or creator of the listing
    CapacityGB     int       // Storage capacity in GB
    PricePerGB     float64   // Price per GB in SYNN
    LeaseDuration  int       // Lease duration in days
    PostedAt       time.Time // Timestamp of when the listing was created
    EncryptedDetails string  // Encrypted storage details (description, terms, etc.)
    Active         bool      // Whether the listing is active
}

// StorageMarketplace manages the listing and leasing of storage on the blockchain
type StorageMarketplace struct {
    mutex          sync.Mutex                    // Mutex for thread-safe operations
    Listings       map[string]*StorageListing    // Map of ListingID to StorageListing
    LedgerInstance *ledger.Ledger                // Ledger instance for tracking storage transactions
    EscrowAccounts map[string]*EscrowAccount     // Map of EscrowID to EscrowAccount
}

// OffchainStorage represents an off-chain storage unit
type OffchainStorage struct {
    StorageID      string    // Unique identifier for the off-chain storage unit
    Owner          string    // Owner or manager of the off-chain storage
    Location       string    // Physical or cloud location of the storage
    CapacityGB     int       // Total capacity of the storage in GB
    UsedCapacityGB int       // Capacity used in GB
    EncryptedDetails string  // Encrypted details of the storage
    Active         bool      // Whether the storage is active or not
}

// OffchainStorageManager manages off-chain storage listings
type OffchainStorageManager struct {
    mutex          sync.Mutex                   // Mutex for thread-safe operations
    StorageUnits   map[string]*OffchainStorage  // Map of storage units
    LedgerInstance *ledger.Ledger               // Ledger instance for tracking storage transactions
}

// StorageRetrievalManager handles the retrieval of on-chain and off-chain storage data
type StorageRetrievalManager struct {
    OnChainData    map[string][]byte               // On-chain data stored by transaction hash
    OffChainUnits  map[string]*OffchainStorage     // Map of off-chain storage units
    LedgerInstance *ledger.Ledger                  // Ledger instance for verification
    mutex          sync.Mutex                      // Mutex for thread-safe operations
}

// StorageSanitizationManager handles sanitizing data before storage and removing sensitive or irrelevant data from the system
type StorageSanitizationManager struct {
    LedgerInstance *ledger.Ledger // Ledger instance for recording sanitized data transactions
    mutex          sync.Mutex     // Mutex for thread-safe operations
}

// StorageManager manages the creation, retrieval, and deletion of storage entries
type StorageManager struct {
    LedgerInstance *ledger.Ledger // Ledger instance for transaction recording
    mutex          sync.Mutex     // Mutex for thread-safe operations
    StorageMap     map[string]StorageEntry // Map of storage entries, identified by storage ID
}

// SwarmManager manages the interaction with the decentralized storage swarm
type SwarmManager struct {
    LedgerInstance *ledger.Ledger           // Ledger instance for transaction logging
    mutex          sync.Mutex               // Mutex for thread-safe operations
    SwarmNodes     map[string]SwarmNode     // Map of swarm nodes connected to the system
    StorageMap     map[string]StorageEntry  // Map of storage entries, identified by storage ID
}

// TimestampManager manages the process of timestamping data and logging it to the ledger
type TimestampManager struct {
    LedgerInstance *ledger.Ledger   // Ledger for logging timestamped data
}

// StorageEntry represents a storage entry in the system
type StorageEntry struct {
    StorageID      string    // Unique identifier for the storage entry
    Data           []byte    // Encrypted data stored in the entry
    Owner          string    // Owner of the storage entry (address or ID)
    CreatedAt      time.Time // Timestamp of when the storage entry was created
    LastAccessed   time.Time // Timestamp of when the storage entry was last accessed
    Expiration     time.Time // Expiration time for the storage entry (optional)
    IsActive       bool      // Indicates if the storage entry is active or expired
}

// SwarmNode represents a node in the decentralized storage swarm
type SwarmNode struct {
    NodeID         string    // Unique identifier for the swarm node
    IPAddress      string    // IP address of the swarm node
    Location       string    // Physical or cloud location of the swarm node
    CapacityGB     int       // Total storage capacity in GB
    UsedCapacityGB int       // Used storage capacity in GB
    Status         string    // Status of the node (active, inactive, etc.)
    LastActive     time.Time // Timestamp of the last time the node was active
    EncryptedData  map[string][]byte // Encrypted data stored on the node
}
