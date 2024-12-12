package resource_management

import (
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/encryption"
)

// NewResourceAllocationManager initializes a new resource allocation manager
func NewResourceAllocationManager() *common.ResourceAllocationManager {
	return &common.ResourceAllocationManager{
		AllocatedResources: make(map[string]common.Resource),
	}
}

// AllocateResource allocates system resources to a node and records it in the ledger
func (ram *common.ResourceAllocationManager) AllocateResource(nodeID string, resource common.Resource) error {
	ram.mutex.Lock()
	defer ram.mutex.Unlock()

	// Encrypt the resource allocation information for security
	encryptedResource, err := encryption.EncryptResource(resource, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("Failed to encrypt resource: %v", err)
	}

	// Allocate the resource to the node
	ram.AllocatedResources[nodeID] = resource
	fmt.Printf("Allocated resource %s to node %s\n", resource.ResourceID, nodeID)

	// Record the allocation in the ledger
	err = ram.LedgerInstance.RecordResourceAllocation(nodeID, encryptedResource)
	if err != nil {
		return fmt.Errorf("Failed to record resource allocation for node %s in ledger: %v", nodeID, err)
	}

	return nil
}

// ReleaseResource releases a previously allocated resource and updates the ledger
func (ram *common.ResourceAllocationManager) ReleaseResource(nodeID string) error {
	ram.mutex.Lock()
	defer ram.mutex.Unlock()

	resource, exists := ram.AllocatedResources[nodeID]
	if !exists {
		return fmt.Errorf("Resource for node %s not found", nodeID)
	}

	// Remove the resource allocation
	delete(ram.AllocatedResources, nodeID)
	fmt.Printf("Released resource %s from node %s\n", resource.ResourceID, nodeID)

	// Record the resource release in the ledger
	err := ram.LedgerInstance.RemoveResourceAllocation(nodeID)
	if err != nil {
		return fmt.Errorf("Failed to record resource release for node %s in ledger: %v", nodeID, err)
	}

	return nil
}

// MonitorResourceUsage monitors the usage of allocated resources and flags any issues
func (ram *common.ResourceAllocationManager) MonitorResourceUsage() {
	ram.mutex.Lock()
	defer ram.mutex.Unlock()

	// Iterate over allocated resources and check their usage
	for nodeID, resource := range ram.AllocatedResources {
		// Simulate monitoring resource usage
		usage := resource.Usage
		if usage > resource.Limit {
			fmt.Printf("Node %s is exceeding resource limits (Usage: %.2f, Limit: %.2f)\n", nodeID, usage, resource.Limit)

			// Record the issue in the ledger
			err := ram.LedgerInstance.RecordResourceUsageIssue(nodeID, resource)
			if err != nil {
				fmt.Printf("Failed to record resource usage issue for node %s: %v\n", nodeID, err)
			}
		}
	}
}

// AdjustResource dynamically adjusts the allocated resources for a node
func (ram *common.ResourceAllocationManager) AdjustResource(nodeID string, newLimit float64) error {
	ram.mutex.Lock()
	defer ram.mutex.Unlock()

	resource, exists := ram.AllocatedResources[nodeID]
	if !exists {
		return fmt.Errorf("Resource for node %s not found", nodeID)
	}

	// Adjust the resource limit
	resource.Limit = newLimit
	ram.AllocatedResources[nodeID] = resource
	fmt.Printf("Adjusted resource limit for node %s to %.2f\n", nodeID, newLimit)

	// Update the ledger with the new resource allocation
	encryptedResource, err := encryption.EncryptResource(resource, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("Failed to encrypt adjusted resource: %v", err)
	}

	err = ram.LedgerInstance.RecordResourceAllocation(nodeID, encryptedResource)
	if err != nil {
		return fmt.Errorf("Failed to update resource allocation for node %s in ledger: %v", nodeID, err)
	}

	return nil
}

// EncryptAndSyncResources ensures all resource allocations are encrypted and synced with the ledger
func (ram *common.ResourceAllocationManager) EncryptAndSyncResources() {
	for nodeID, resource := range ram.AllocatedResources {
		// Encrypt the resource data
		encryptedResource, err := encryption.EncryptResource(resource, common.EncryptionKey)
		if err != nil {
			fmt.Printf("Failed to encrypt resource for node %s: %v\n", nodeID, err)
			continue
		}

		// Sync the encrypted resource with the ledger
		err = ram.LedgerInstance.RecordResourceAllocation(nodeID, encryptedResource)
		if err != nil {
			fmt.Printf("Failed to sync resource allocation for node %s with ledger: %v\n", nodeID, err)
		}
	}
}

// generateResourceHash creates a hash to validate resource integrity
func generateResourceHash(resource common.Resource) string {
	input := fmt.Sprintf("%s%f%f", resource.ResourceID, resource.Usage, resource.Limit)
	hash := sha256.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

// ReportResourceStatus logs the current status of all allocated resources
func (ram *common.ResourceAllocationManager) ReportResourceStatus() {
	ram.mutex.Lock()
	defer ram.mutex.Unlock()

	fmt.Println("Reporting resource allocation status...")
	for nodeID, resource := range ram.AllocatedResources {
		fmt.Printf("Node %s: Resource %s, Usage: %.2f, Limit: %.2f\n", nodeID, resource.ResourceID, resource.Usage, resource.Limit)
	}
}
