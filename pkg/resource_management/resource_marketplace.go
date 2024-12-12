package resource_management

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/coin"
    "synnergy_network/pkg/network"
)


// NewResourceMarketplace creates a new resource marketplace
func NewResourceMarketplace(ledgerInstance *ledger.Ledger) *common.ResourceMarketplace {
    return &common.ResourceMarketplace{
        AvailableResources: make(map[string]common.Resource),
        LeasedResources:    make(map[string]common.Lease),
        Purchases:          make(map[string]common.Purchase),
        LedgerInstance:     ledgerInstance,
    }
}

// ListResource adds a resource to the marketplace for leasing or purchasing with a price
func (rm *common.ResourceMarketplace) ListResource(resource common.Resource, pricePerUnit float64) error {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    resource.PricePerUnit = pricePerUnit
    encryptedResource, err := encryption.EncryptResource(resource, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt resource: %v", err)
    }

    rm.AvailableResources[resource.ID] = encryptedResource
    fmt.Printf("Resource %s listed for %s with %d units available at %.2f SYNN per unit.\n", resource.ID, resource.Type, resource.AvailableUnits, pricePerUnit)

    // Record the listing in the ledger
    err = rm.LedgerInstance.RecordResourceListing(resource.ID, encryptedResource)
    if err != nil {
        return fmt.Errorf("failed to record resource listing in the ledger: %v", err)
    }

    return nil
}

// LeaseResource leases a specific resource for a certain period using Synnergy Coin (SYNN) for payment
func (rm *common.ResourceMarketplace) LeaseResource(leaseRequest common.LeaseRequest, paymentAddress string, paymentAmount float64) error {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    resource, exists := rm.AvailableResources[leaseRequest.ResourceID]
    if !exists {
        return fmt.Errorf("resource %s not found", leaseRequest.ResourceID)
    }

    requiredAmount := resource.PricePerUnit * float64(leaseRequest.RequiredUnits)
    if paymentAmount < requiredAmount {
        return fmt.Errorf("insufficient SYNN sent for lease. Required: %.2f, Sent: %.2f", requiredAmount, paymentAmount)
    }

    if resource.AvailableUnits < leaseRequest.RequiredUnits {
        return fmt.Errorf("insufficient units available for leasing")
    }

    // Perform the payment using SYNN coin
    err := coin.TransferCoins(paymentAddress, rm.LedgerInstance.GetMarketplaceWallet(), paymentAmount)
    if err != nil {
        return fmt.Errorf("failed to complete payment: %v", err)
    }

    // Deduct the leased units from the resource
    resource.AvailableUnits -= leaseRequest.RequiredUnits
    rm.AvailableResources[resource.ID] = resource

    lease := common.Lease{
        LeaseID:       fmt.Sprintf("%s-lease-%d", leaseRequest.ResourceID, time.Now().UnixNano()),
        ResourceID:    leaseRequest.ResourceID,
        NodeID:        leaseRequest.NodeID,
        RequiredUnits: leaseRequest.RequiredUnits,
        LeaseStart:    time.Now(),
        LeaseEnd:      time.Now().Add(leaseRequest.Duration),
    }

    encryptedLease, err := encryption.EncryptLease(lease, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt lease details: %v", err)
    }

    rm.LeasedResources[lease.LeaseID] = encryptedLease
    fmt.Printf("Resource %s leased to node %s for %d units until %s.\n", leaseRequest.ResourceID, leaseRequest.NodeID, leaseRequest.RequiredUnits, lease.LeaseEnd)

    // Record the lease in the ledger
    err = rm.LedgerInstance.RecordResourceLease(lease.LeaseID, encryptedLease)
    if err != nil {
        return fmt.Errorf("failed to record resource lease in the ledger: %v", err)
    }

    return nil
}

// PurchaseResource handles the full purchase of a resource by a node using SYNN
func (rm *common.ResourceMarketplace) PurchaseResource(purchaseRequest common.PurchaseRequest, paymentAddress string, paymentAmount float64) error {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    resource, exists := rm.AvailableResources[purchaseRequest.ResourceID]
    if !exists {
        return fmt.Errorf("resource %s not found", purchaseRequest.ResourceID)
    }

    requiredAmount := resource.PricePerUnit * float64(purchaseRequest.RequiredUnits)
    if paymentAmount < requiredAmount {
        return fmt.Errorf("insufficient SYNN sent for purchase. Required: %.2f, Sent: %.2f", requiredAmount, paymentAmount)
    }

    // Perform the payment using SYNN coin
    err := coin.TransferCoins(paymentAddress, rm.LedgerInstance.GetMarketplaceWallet(), paymentAmount)
    if err != nil {
        return fmt.Errorf("failed to complete payment: %v", err)
    }

    if resource.AvailableUnits < purchaseRequest.RequiredUnits {
        return fmt.Errorf("insufficient units available for purchase")
    }

    // Deduct the purchased units from the resource
    resource.AvailableUnits -= purchaseRequest.RequiredUnits
    rm.AvailableResources[resource.ID] = resource

    purchase := common.Purchase{
        PurchaseID:    fmt.Sprintf("%s-purchase-%d", purchaseRequest.ResourceID, time.Now().UnixNano()),
        ResourceID:    purchaseRequest.ResourceID,
        NodeID:        purchaseRequest.NodeID,
        RequiredUnits: purchaseRequest.RequiredUnits,
        PurchaseTime:  time.Now(),
    }

    encryptedPurchase, err := encryption.EncryptPurchase(purchase, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt purchase details: %v", err)
    }

    rm.Purchases[purchase.PurchaseID] = encryptedPurchase
    fmt.Printf("Resource %s purchased by node %s for %d units.\n", purchaseRequest.ResourceID, purchaseRequest.NodeID, purchaseRequest.RequiredUnits)

    // Record the purchase in the ledger
    err = rm.LedgerInstance.RecordResourcePurchase(purchase.PurchaseID, encryptedPurchase)
    if err != nil {
        return fmt.Errorf("failed to record resource purchase in the ledger: %v", err)
    }

    return nil
}

// QueryResource allows querying a resource by sending a message over the P2P network
func (rm *common.ResourceMarketplace) QueryResource(resourceID, requestingNodeID string) (*common.Resource, error) {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    resource, exists := rm.AvailableResources[resourceID]
    if !exists {
        return nil, fmt.Errorf("resource %s not found", resourceID)
    }

    // Encrypt and send a query message using P2P communication
    encryptedQuery, err := encryption.EncryptMessage(fmt.Sprintf("Querying resource %s", resourceID), common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt query message: %v", err)
    }

    err = p2p.SendMessage(requestingNodeID, encryptedQuery)
    if err != nil {
        return nil, fmt.Errorf("failed to send P2P message: %v", err)
    }

    return &resource, nil
}
