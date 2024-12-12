package ledger

import (
	"fmt"
	"time"
)


// RecordResourceAllocation allocates a resource to an account
func (ledger *ResourceManagementLedger) RecordResourceAllocation(resourceID, accountID string, limit float64) {
	ledger.Lock()
	defer ledger.Unlock()

	newResource := Resource{
		ID:          resourceID,
		AllocatedTo: accountID,
		Limit:       limit,
		UsedAmount:  0,
		CreatedAt:   time.Now(),
	}
	ledger.ResourceManagementLedgerState.Resources[resourceID] = newResource
	fmt.Printf("Resource %s allocated to account %s with limit %.2f.\n", resourceID, accountID, limit)
}

// RemoveResourceAllocation removes a resource allocation
func (ledger *ResourceManagementLedger) RemoveResourceAllocation(resourceID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if _, exists := ledger.ResourceManagementLedgerState.Resources[resourceID]; exists {
		delete(ledger.ResourceManagementLedgerState.Resources, resourceID)
		fmt.Printf("Resource %s allocation removed.\n", resourceID)
	} else {
		fmt.Printf("Resource %s not found.\n", resourceID)
	}
}

// RecordResourceUsageIssue records any issue related to resource usage.
func (ledger *ResourceManagementLedger) RecordResourceUsageIssue(resourceID, issue string) {
    ledger.Lock()
    defer ledger.Unlock()

    // Initialize the map if it hasn't been created yet
    if ledger.ResourceIssues == nil {
        ledger.ResourceIssues = make(map[string][]ResourceIssue)
    }

    // Generate a unique ID for the issue
    issueID := fmt.Sprintf("%s_%d", resourceID, time.Now().UnixNano())

    // Create a new ResourceIssue instance
    newIssue := ResourceIssue{
        IssueID:     issueID,
        ResourceID:  resourceID,
        Description: issue,
        ReportedAt:  time.Now(), // Correct field name
        Resolved:    false,      // Defaults to unresolved
    }

    // Append the issue to the list for the given resource ID
    ledger.ResourceIssues[resourceID] = append(ledger.ResourceIssues[resourceID], newIssue)

    fmt.Printf("Resource %s usage issue recorded: %s\n", resourceID, issue)
}



// RecordResourceAddition adds a new resource to the ledger
func (ledger *ResourceManagementLedger) RecordResourceAddition(resourceID string, limit float64) {
	ledger.Lock()
	defer ledger.Unlock()

	newResource := Resource{
		ID:         resourceID,
		Limit:      limit,
		UsedAmount: 0,
		CreatedAt:  time.Now(),
	}
	ledger.ResourceManagementLedgerState.Resources[resourceID] = newResource
	fmt.Printf("New resource %s added with limit %.2f.\n", resourceID, limit)
}

// RecordResourceRelease releases a resource from allocation
func (ledger *ResourceManagementLedger) RecordResourceRelease(resourceID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if resource, exists := ledger.ResourceManagementLedgerState.Resources[resourceID]; exists {
		resource.AllocatedTo = ""
		ledger.ResourceManagementLedgerState.Resources[resourceID] = resource
		fmt.Printf("Resource %s released from its allocation.\n", resourceID)
	} else {
		fmt.Printf("Resource %s not found.\n", resourceID)
	}
}

// RecordResourceOveruse records an overuse of a resource
func (ledger *ResourceManagementLedger) RecordResourceOveruse(resourceID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if resource, exists := ledger.ResourceManagementLedgerState.Resources[resourceID]; exists {
		resource.OveruseReported = true
		ledger.ResourceManagementLedgerState.Resources[resourceID] = resource
		fmt.Printf("Overuse reported for resource %s.\n", resourceID)
	} else {
		fmt.Printf("Resource %s not found.\n", resourceID)
	}
}

// RecordResourceAdjustment adjusts the usage limit of a resource
func (ledger *ResourceManagementLedger) RecordResourceAdjustment(resourceID string, newLimit float64) {
	ledger.Lock()
	defer ledger.Unlock()

	if resource, exists := ledger.ResourceManagementLedgerState.Resources[resourceID]; exists {
		resource.Limit = newLimit
		ledger.ResourceManagementLedgerState.Resources[resourceID] = resource
		fmt.Printf("Resource %s limit adjusted to %.2f.\n", resourceID, newLimit)
	} else {
		fmt.Printf("Resource %s not found.\n", resourceID)
	}
}

// GetMarketplaceWallet retrieves the wallet address for the marketplace operations.
func (ledger *ResourceManagementLedger) GetMarketplaceWallet() string {
    ledger.Lock()
    defer ledger.Unlock()

    // Check if the MarketplaceWallet field is initialized
    if ledger.MarketplaceWallet == "" {
        // Default or error handling if no wallet address is set
        fmt.Println("Marketplace wallet address is not set.")
        return ""
    }

    fmt.Printf("Marketplace wallet retrieved: %s\n", ledger.MarketplaceWallet)
    return ledger.MarketplaceWallet
}


// RecordResourceLease logs and stores the lease of a resource.
func (ledger *ResourceManagementLedger) RecordResourceLease(resourceID, lesseeID string, leaseAmount float64) error {
    ledger.Lock()
    defer ledger.Unlock()

    // Validate inputs
    if resourceID == "" || lesseeID == "" {
        return fmt.Errorf("resourceID and lesseeID cannot be empty")
    }
    if leaseAmount <= 0 {
        return fmt.Errorf("leaseAmount must be greater than zero")
    }

    // Initialize the ResourceLeases map if not already initialized
    if ledger.ResourceLeases == nil {
        ledger.ResourceLeases = make(map[string]LeaseRecord)
    }

    // Generate a unique LeaseID
    leaseID := fmt.Sprintf("%s_%s_%d", resourceID, lesseeID, time.Now().UnixNano())

    // Create a new lease record
    leaseStartTime := time.Now()
    leaseDuration := time.Hour * 24 * 30 // Example: default duration is 30 days
    leaseRecord := LeaseRecord{
        LeaseID:        leaseID,
        ResourceID:     resourceID,
        LesseeID:       lesseeID,
        LeaseAmount:    leaseAmount,
        LeaseDuration:  leaseDuration,
        LeaseStartTime: leaseStartTime,
        LeaseEndTime:   leaseStartTime.Add(leaseDuration),
        Status:         "Active",
        AdditionalInfo: map[string]string{"CreatedBy": "System"}, // Example additional info
    }

    // Store the lease record in the ledger
    ledger.ResourceLeases[resourceID] = leaseRecord

    // Log the lease
    fmt.Printf("Resource %s leased to %s for %.2f, starting at %s and ending at %s.\n",
        resourceID, lesseeID, leaseAmount,
        leaseRecord.LeaseStartTime.Format(time.RFC3339),
        leaseRecord.LeaseEndTime.Format(time.RFC3339),
    )
    return nil
}



// RecordResourceRegistration registers a resource to the ledger
func (ledger *ResourceManagementLedger) RecordResourceRegistration(resourceID, ownerID string, usageLimit float64) {
	ledger.Lock()
	defer ledger.Unlock()

	newResource := Resource{
		ID:          resourceID,
		AllocatedTo: ownerID,
		Limit:       usageLimit,
		UsedAmount:  0,
		CreatedAt:   time.Now(),
	}
	ledger.ResourceManagementLedgerState.Resources[resourceID] = newResource
	fmt.Printf("Resource %s registered to account %s with limit %.2f.\n", resourceID, ownerID, usageLimit)
}

// EnableResourcePooling enables resource pooling.
func (l *ResourceManagementLedger) EnableResourcePooling() error {
    l.ResourcePool.Enabled = true
    l.ResourcePool.LastUpdated = time.Now()
    return nil
}

// DisableResourcePooling disables resource pooling.
func (l *ResourceManagementLedger) DisableResourcePooling() error {
    l.ResourcePool.Enabled = false
    l.ResourcePool.LastUpdated = time.Now()
    return nil
}

// SetResourcePoolingPolicy sets the resource pooling policy.
func (l *ResourceManagementLedger) SetResourcePoolingPolicy(policy string) error {
    l.ResourcePoolingPolicy = ResourcePoolingPolicy{
        Policy:      policy,
        LastUpdated: time.Now(),
    }
    return nil
}

// GetResourcePoolingPolicy retrieves the resource pooling policy.
func (l *ResourceManagementLedger) GetResourcePoolingPolicy() (string, error) {
    if l.ResourcePoolingPolicy.Policy == "" {
        return "", fmt.Errorf("resource pooling policy is not set")
    }
    return l.ResourcePoolingPolicy.Policy, nil
}
