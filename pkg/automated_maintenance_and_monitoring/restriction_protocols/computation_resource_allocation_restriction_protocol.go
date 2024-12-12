package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

const (
	ComputationResourceCheckInterval = 10 * time.Second  // Interval for monitoring resource allocation
	MaxResourceAllocationPerUser     = 500.0             // Maximum computation resources allocated per user in the system
	MinResourceAllocationRequest     = 0.1               // Minimum request for resource allocation to avoid system abuse
)

// ComputationResourceAllocationRestrictionAutomation handles restrictions for resource allocation across the network
type ComputationResourceAllocationRestrictionAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	userResourceUsage     map[string]float64 // Tracks computation resource usage by user
}

// NewComputationResourceAllocationRestrictionAutomation initializes and returns an instance of ComputationResourceAllocationRestrictionAutomation
func NewComputationResourceAllocationRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ComputationResourceAllocationRestrictionAutomation {
	return &ComputationResourceAllocationRestrictionAutomation{
		consensusSystem:   consensusSystem,
		ledgerInstance:    ledgerInstance,
		stateMutex:        stateMutex,
		userResourceUsage: make(map[string]float64),
	}
}

// StartResourceAllocationMonitoring begins continuous monitoring of resource allocation across the network
func (automation *ComputationResourceAllocationRestrictionAutomation) StartResourceAllocationMonitoring() {
	ticker := time.NewTicker(ComputationResourceCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorResourceAllocations()
		}
	}()
}

// monitorResourceAllocations continuously checks for valid resource allocation requests and flags any violations
func (automation *ComputationResourceAllocationRestrictionAutomation) monitorResourceAllocations() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent resource allocation requests from Synnergy Consensus
	recentAllocations := automation.consensusSystem.GetRecentResourceAllocations()

	for _, allocation := range recentAllocations {
		// Validate each resource allocation request
		if !automation.validateResourceAllocation(allocation) {
			automation.flagResourceAllocationViolation(allocation, "Exceeded maximum resource allocation per user")
			continue
		}

		if !automation.validateMinResourceRequest(allocation) {
			automation.flagResourceAllocationViolation(allocation, "Requested resource allocation below the minimum allowed")
		}
	}
}

// validateResourceAllocation ensures that the resource allocation does not exceed the maximum per user
func (automation *ComputationResourceAllocationRestrictionAutomation) validateResourceAllocation(allocation common.ResourceAllocation) bool {
	currentUsage := automation.userResourceUsage[allocation.UserID]
	if currentUsage+allocation.ResourcesRequested > MaxResourceAllocationPerUser {
		return false
	}

	// Update the usage for the user
	automation.userResourceUsage[allocation.UserID] += allocation.ResourcesRequested
	return true
}

// validateMinResourceRequest ensures that resource requests meet the minimum required to prevent system abuse
func (automation *ComputationResourceAllocationRestrictionAutomation) validateMinResourceRequest(allocation common.ResourceAllocation) bool {
	return allocation.ResourcesRequested >= MinResourceAllocationRequest
}

// flagResourceAllocationViolation flags an allocation request that violates system rules and logs it in the ledger
func (automation *ComputationResourceAllocationRestrictionAutomation) flagResourceAllocationViolation(allocation common.ResourceAllocation, reason string) {
	fmt.Printf("Resource allocation violation: User %s, Reason: %s\n", allocation.UserID, reason)

	// Log the violation into the ledger
	automation.logResourceAllocationViolation(allocation, reason)
}

// logResourceAllocationViolation logs the flagged resource allocation violation into the ledger with full details
func (automation *ComputationResourceAllocationRestrictionAutomation) logResourceAllocationViolation(allocation common.ResourceAllocation, violationReason string) {
	// Encrypt the resource allocation data
	encryptedData := automation.encryptResourceAllocationData(allocation)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("resource-allocation-violation-%s-%d", allocation.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Resource Allocation Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for resource allocation violation. Reason: %s. Encrypted Data: %s", allocation.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log resource allocation violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Resource allocation violation logged for user: %s\n", allocation.UserID)
	}
}

// encryptResourceAllocationData encrypts resource allocation data before logging for security
func (automation *ComputationResourceAllocationRestrictionAutomation) encryptResourceAllocationData(allocation common.ResourceAllocation) string {
	data := fmt.Sprintf("User ID: %s, Resources Requested: %.2f", allocation.UserID, allocation.ResourcesRequested)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting resource allocation data:", err)
		return data
	}
	return string(encryptedData)
}
