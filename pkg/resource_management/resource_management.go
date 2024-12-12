package resource_management

import (
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/encryption"
)

// NewResourceManager initializes a new resource manager with an empty resource pool and ledger integration.
func NewResourceManager() *common.ResourceManager {
	return &common.ResourceManager{
		Resources:       make(map[string]common.Resource),
		AllocationQueue: []common.ResourceRequest{},
	}
}

// AddResource adds a new resource to the resource pool.
func (rm *common.ResourceManager) AddResource(resourceID string, resource common.Resource) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Encrypt the resource details for security before storing
	encryptedResource, err := encryption.EncryptResource(resource, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("Failed to encrypt resource %s: %v", resourceID, err)
	}

	rm.Resources[resourceID] = encryptedResource
	fmt.Printf("Resource %s added to the pool.\n", resourceID)

	// Record the new resource in the ledger
	err = rm.LedgerInstance.RecordResourceAddition(resourceID, encryptedResource)
	if err != nil {
		return fmt.Errorf("Failed to record resource %s in the ledger: %v", resourceID, err)
	}

	return nil
}

// AllocateResource handles resource allocation for a specific request.
func (rm *common.ResourceManager) AllocateResource(request common.ResourceRequest) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Check if requested resource is available
	resource, exists := rm.Resources[request.ResourceID]
	if !exists || resource.AvailableUnits < request.RequiredUnits {
		// If resource is not available or insufficient, enqueue the request
		rm.AllocationQueue = append(rm.AllocationQueue, request)
		fmt.Printf("Resource %s allocation delayed. Added to the queue.\n", request.ResourceID)
		return fmt.Errorf("resource allocation delayed, insufficient units or not available")
	}

	// Deduct the required units from the resource pool
	resource.AvailableUnits -= request.RequiredUnits
	rm.Resources[request.ResourceID] = resource

	// Encrypt and log the resource allocation in the ledger
	encryptedRequest, err := encryption.EncryptResourceRequest(request, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("Failed to encrypt resource request: %v", err)
	}

	err = rm.LedgerInstance.RecordResourceAllocation(request.NodeID, encryptedRequest)
	if err != nil {
		return fmt.Errorf("Failed to record resource allocation in the ledger: %v", err)
	}

	fmt.Printf("Allocated %d units of resource %s to node %s.\n", request.RequiredUnits, request.ResourceID, request.NodeID)
	return nil
}

// ReleaseResource releases allocated resources back into the resource pool.
func (rm *common.ResourceManager) ReleaseResource(resourceID string, units int) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	resource, exists := rm.Resources[resourceID]
	if !exists {
		return fmt.Errorf("resource %s not found", resourceID)
	}

	// Add the units back to the resource pool
	resource.AvailableUnits += units
	rm.Resources[resourceID] = resource
	fmt.Printf("Released %d units of resource %s.\n", units, resourceID)

	// Update the ledger with the released resource units
	err := rm.LedgerInstance.RecordResourceRelease(resourceID, units)
	if err != nil {
		return fmt.Errorf("Failed to record resource release in the ledger: %v", err)
	}

	// Process the allocation queue to handle any pending requests
	rm.processAllocationQueue()

	return nil
}

// processAllocationQueue processes pending resource allocation requests in the queue.
func (rm *common.ResourceManager) processAllocationQueue() {
	var remainingQueue []common.ResourceRequest

	for _, request := range rm.AllocationQueue {
		// Attempt to allocate the requested resources
		err := rm.AllocateResource(request)
		if err != nil {
			// If allocation fails, keep the request in the queue
			remainingQueue = append(remainingQueue, request)
		}
	}

	rm.AllocationQueue = remainingQueue
}

// MonitorResourceUsage continuously monitors resource usage and flags any overuse.
func (rm *common.ResourceManager) MonitorResourceUsage() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	for resourceID, resource := range rm.Resources {
		if resource.Usage > resource.Limit {
			fmt.Printf("Resource %s exceeding usage limit. Current usage: %.2f, Limit: %.2f.\n", resourceID, resource.Usage, resource.Limit)

			// Log the overuse event to the ledger
			err := rm.LedgerInstance.RecordResourceOveruse(resourceID, resource)
			if err != nil {
				fmt.Printf("Failed to log resource overuse for %s: %v\n", resourceID, err)
			}
		}
	}
}

// AdjustResource dynamically adjusts the limit of a resource.
func (rm *common.ResourceManager) AdjustResource(resourceID string, newLimit float64) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	resource, exists := rm.Resources[resourceID]
	if !exists {
		return fmt.Errorf("resource %s not found", resourceID)
	}

	resource.Limit = newLimit
	rm.Resources[resourceID] = resource
	fmt.Printf("Resource %s limit adjusted to %.2f.\n", resourceID, newLimit)

	// Update the ledger with the new resource allocation
	encryptedResource, err := encryption.EncryptResource(resource, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("Failed to encrypt adjusted resource: %v", err)
	}

	err = rm.LedgerInstance.RecordResourceAdjustment(resourceID, encryptedResource)
	if err != nil {
		return fmt.Errorf("Failed to update resource allocation in ledger for resource %s: %v", resourceID, err)
	}

	return nil
}

// ReportResourceStatus provides a summary of all available and allocated resources.
func (rm *common.ResourceManager) ReportResourceStatus() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	fmt.Println("Resource Management Status Report:")
	for resourceID, resource := range rm.Resources {
		fmt.Printf("Resource %s: Available Units: %d, Usage: %.2f, Limit: %.2f\n", resourceID, resource.AvailableUnits, resource.Usage, resource.Limit)
	}
}
