package resource_management

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
)


// NewResourceManager creates and returns a new ResourceManager instance.
func NewResourceManager(ledgerInstance *ledger.Ledger) *common.ResourceManager {
    return &common.ResourceManager{
        Resources:       make(map[string]common.Resource),
        LedgerInstance:  ledgerInstance,
        AllocationQueue: []common.ResourceRequest{},
    }
}

// RegisterResource registers a new resource within the system.
func (rm *common.ResourceManager) RegisterResource(resourceID string, resourceType string, availableUnits int, limit float64) error {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    // Encrypt resource data before storing it
    resource := common.Resource{
        ID:            resourceID,
        Type:          resourceType,
        AvailableUnits: availableUnits,
        Limit:         limit,
        CreatedAt:     time.Now(),
    }

    encryptedResource, err := encryption.EncryptResource(resource, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt resource data: %v", err)
    }

    rm.Resources[resourceID] = encryptedResource
    fmt.Printf("Resource %s of type %s registered with %d available units.\n", resourceID, resourceType, availableUnits)

    // Record the resource in the ledger
    err = rm.LedgerInstance.RecordResourceRegistration(resourceID, encryptedResource)
    if err != nil {
        return fmt.Errorf("failed to record resource registration in the ledger: %v", err)
    }

    return nil
}

// AllocateResource allocates units from a specific resource to a requesting node.
func (rm *common.ResourceManager) AllocateResource(resourceID string, request common.ResourceRequest) error {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    resource, exists := rm.Resources[resourceID]
    if !exists {
        return fmt.Errorf("resource %s not found", resourceID)
    }

    if resource.AvailableUnits < request.RequiredUnits {
        // Add to the allocation queue if resource units are insufficient
        rm.AllocationQueue = append(rm.AllocationQueue, request)
        fmt.Printf("Resource %s allocation deferred for %d units. Added to the queue.\n", resourceID, request.RequiredUnits)
        return fmt.Errorf("insufficient units for resource allocation")
    }

    // Deduct the requested units
    resource.AvailableUnits -= request.RequiredUnits
    rm.Resources[resourceID] = resource

    // Log allocation to ledger
    encryptedRequest, err := encryption.EncryptResourceRequest(request, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt resource allocation request: %v", err)
    }

    err = rm.LedgerInstance.RecordResourceAllocation(request.NodeID, encryptedRequest)
    if err != nil {
        return fmt.Errorf("failed to log resource allocation in the ledger: %v", err)
    }

    fmt.Printf("Allocated %d units of resource %s to node %s.\n", request.RequiredUnits, resourceID, request.NodeID)
    return nil
}

// ReleaseResource releases previously allocated units back into the resource pool.
func (rm *common.ResourceManager) ReleaseResource(resourceID string, units int) error {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    resource, exists := rm.Resources[resourceID]
    if !exists {
        return fmt.Errorf("resource %s not found", resourceID)
    }

    // Add units back to the pool
    resource.AvailableUnits += units
    rm.Resources[resourceID] = resource

    // Record release in ledger
    err := rm.LedgerInstance.RecordResourceRelease(resourceID, units)
    if err != nil {
        return fmt.Errorf("failed to log resource release in the ledger: %v", err)
    }

    fmt.Printf("Released %d units of resource %s back into the pool.\n", units, resourceID)
    rm.processAllocationQueue() // Process any pending allocations
    return nil
}

// processAllocationQueue processes pending resource allocation requests.
func (rm *common.ResourceManager) processAllocationQueue() {
    remainingQueue := []common.ResourceRequest{}
    for _, request := range rm.AllocationQueue {
        err := rm.AllocateResource(request.ResourceID, request)
        if err != nil {
            remainingQueue = append(remainingQueue, request)
        }
    }
    rm.AllocationQueue = remainingQueue
}

// MonitorResource monitors usage and ensures it stays within set limits.
func (rm *common.ResourceManager) MonitorResource(resourceID string) {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    resource, exists := rm.Resources[resourceID]
    if !exists {
        fmt.Printf("Resource %s not found for monitoring.\n", resourceID)
        return
    }

    if resource.Usage > resource.Limit {
        fmt.Printf("Warning: Resource %s exceeds usage limit.\n", resourceID)
        err := rm.LedgerInstance.RecordResourceOveruse(resourceID, resource)
        if err != nil {
            fmt.Printf("Failed to record resource overuse for %s in the ledger: %v\n", resourceID, err)
        }
    }
}

// AdjustResourceLimit dynamically changes the limit for a resource.
func (rm *common.ResourceManager) AdjustResourceLimit(resourceID string, newLimit float64) error {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    resource, exists := rm.Resources[resourceID]
    if !exists {
        return fmt.Errorf("resource %s not found", resourceID)
    }

    resource.Limit = newLimit
    rm.Resources[resourceID] = resource

    // Encrypt and update in the ledger
    encryptedResource, err := encryption.EncryptResource(resource, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt updated resource: %v", err)
    }

    err = rm.LedgerInstance.RecordResourceAdjustment(resourceID, encryptedResource)
    if err != nil {
        return fmt.Errorf("failed to log resource adjustment for %s in the ledger: %v", resourceID, err)
    }

    fmt.Printf("Resource %s limit adjusted to %.2f.\n", resourceID, newLimit)
    return nil
}

// GetResourceStatus retrieves the current status of all registered resources.
func (rm *common.ResourceManager) GetResourceStatus() {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    fmt.Println("Resource Status:")
    for resourceID, resource := range rm.Resources {
        fmt.Printf("Resource %s: Available Units: %d, Usage: %.2f, Limit: %.2f\n", resourceID, resource.AvailableUnits, resource.Usage, resource.Limit)
    }
}
