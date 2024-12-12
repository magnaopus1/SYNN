package common

import (
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)

// Resource represents a system resource managed by the blockchain network.
type Resource struct {
	ID            string    // Unique identifier for the resource
	Type          string    // Type of resource (e.g., "CPU", "Memory", "Storage")
	AvailableUnits int       // Number of units available for allocation
	Usage         float64   // Current usage level
	Limit         float64   // Usage limit for the resource
	CreatedAt     time.Time // Time when the resource was registered
}

// Lease represents the details of a leased resource in the marketplace.
type Lease struct {
	LeaseID       string    // Unique identifier for the lease
	ResourceID    string    // ID of the resource being leased
	Leaser        string    // Wallet address of the leaser
	LeaseDuration time.Duration // Duration of the lease
	StartTime     time.Time // Time when the lease started
	EndTime       time.Time // Time when the lease ends
	Price         float64   // Price paid for the lease
}

// Purchase represents the details of a purchased resource in the marketplace.
type Purchase struct {
	PurchaseID    string    // Unique identifier for the purchase
	ResourceID    string    // ID of the resource being purchased
	Buyer         string    // Wallet address of the buyer
	PurchaseTime  time.Time // Time when the purchase was made
	Price         float64   // Price paid for the purchase
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
	LedgerInstance     *ledger.Ledger      // Ledger for resource tracking and auditing
}

// ResourceManager handles all resource-related operations across the Synnergy Network.
type ResourceManager struct {
	Resources       map[string]Resource         // The resources available for management.
	mutex           sync.Mutex                  // Mutex for thread-safe resource management.
	LedgerInstance  *ledger.Ledger              // The ledger instance to record resource operations.
	AllocationQueue []ResourceRequest           // Queue for pending resource allocation requests.
}

// ResourceMarketplace manages the leasing, purchasing, and selling of resources.
type ResourceMarketplace struct {
	AvailableResources map[string]Resource     // List of available resources for leasing or purchasing
	LeasedResources    map[string]Lease        // Leased resources with details
	Purchases          map[string]Purchase     // Completed resource purchases
	mutex              sync.Mutex              // Mutex for thread-safe operations
	LedgerInstance     *ledger.Ledger          // Ledger instance for tracking resource transactions
}
