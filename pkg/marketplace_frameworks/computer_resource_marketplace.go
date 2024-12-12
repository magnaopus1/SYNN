package marketplace

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// NewComputerResourceMarketplace initializes the marketplace
func NewComputerResourceMarketplace(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.ComputerResourceMarketplace {
	return &common.ComputerResourceMarketplace{
		Resources:        make(map[string]*common.Resource),
		Escrows:          make(map[string]*common.Escrow),
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// ListResource allows a user to list a computing resource for rent or purchase
func (crm *common.ComputerResourceMarketplace) ListResource(resourceID, resourceType, owner string, price float64) (*common.Resource, error) {
	crm.mu.Lock()
	defer crm.mu.Unlock()

	// Encrypt resource data
	resourceData := fmt.Sprintf("ResourceID: %s, Owner: %s, Price: %f", resourceID, owner, price)
	encryptedData, err := crm.EncryptionService.EncryptData([]byte(resourceData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt resource data: %v", err)
	}

	// Create a new resource
	resource := &common.Resource{
		ResourceID:   resourceID,
		ResourceType: resourceType,
		Owner:        owner,
		Price:        price,
		Available:    true,
		ListedTime:   time.Now(),
	}

	// Add the resource to the marketplace
	crm.Resources[resourceID] = resource

	// Log the resource listing in the ledger
	err = crm.Ledger.RecordResourceListing(resourceID, resourceType, owner, price, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log resource listing: %v", err)
	}

	fmt.Printf("Resource %s of type %s listed by %s at price %f\n", resourceID, resourceType, owner, price)
	return resource, nil
}

// RentResource allows a user to rent a computing resource
func (crm *common.ComputerResourceMarketplace) RentResource(resourceID, renter string, duration time.Duration) error {
	crm.mu.Lock()
	defer crm.mu.Unlock()

	// Retrieve the resource
	resource, exists := crm.Resources[resourceID]
	if !exists || !resource.Available {
		return fmt.Errorf("resource %s is not available for rent", resourceID)
	}

	// Create an escrow for the transaction
	escrowID := generateUniqueID()
	escrowAmount := resource.Price * float64(duration.Hours())
	escrow := &common.Escrow{
		EscrowID:   escrowID,
		Buyer:      renter,
		Seller:     resource.Owner,
		Amount:     escrowAmount,
		ResourceID: resourceID,
		IsReleased: false,
		IsDisputed: false,
	}

	// Add escrow to the system
	crm.Escrows[escrowID] = escrow

	// Mark the resource as unavailable
	resource.Available = false

	// Log the resource rental in the ledger
	err := crm.Ledger.RecordResourceRental(resourceID, renter, escrowAmount, duration, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log resource rental: %v", err)
	}

	fmt.Printf("Resource %s rented by %s for duration %v, escrow amount: %f\n", resourceID, renter, duration, escrowAmount)
	return nil
}

// PurchaseResource allows a user to purchase a computing resource
func (crm *common.ComputerResourceMarketplace) PurchaseResource(resourceID, buyer string) error {
	crm.mu.Lock()
	defer crm.mu.Unlock()

	// Retrieve the resource
	resource, exists := crm.Resources[resourceID]
	if !exists || !resource.Available {
		return fmt.Errorf("resource %s is not available for purchase", resourceID)
	}

	// Create an escrow for the purchase
	escrowID := generateUniqueID()
	escrowAmount := resource.Price
	escrow := &common.Escrow{
		EscrowID:   escrowID,
		Buyer:      buyer,
		Seller:     resource.Owner,
		Amount:     escrowAmount,
		ResourceID: resourceID,
		IsReleased: false,
		IsDisputed: false,
	}

	// Add escrow to the system
	crm.Escrows[escrowID] = escrow

	// Mark the resource as sold (no longer available)
	resource.Available = false

	// Log the resource purchase in the ledger
	err := crm.Ledger.RecordResourcePurchase(resourceID, buyer, escrowAmount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log resource purchase: %v", err)
	}

	fmt.Printf("Resource %s purchased by %s for amount %f\n", resourceID, buyer, escrowAmount)
	return nil
}

// ReleaseEscrow releases the funds from escrow after transaction completion
func (crm *common.ComputerResourceMarketplace) ReleaseEscrow(escrowID string) error {
	crm.mu.Lock()
	defer crm.mu.Unlock()

	// Retrieve the escrow
	escrow, exists := crm.Escrows[escrowID]
	if !exists || escrow.IsReleased {
		return fmt.Errorf("escrow %s not found or already released", escrowID)
	}

	// Mark the escrow as released
	escrow.IsReleased = true
	escrow.CompletionTime = time.Now()

	// Log the escrow release in the ledger
	err := crm.Ledger.RecordEscrowRelease(escrowID, escrow.Buyer, escrow.Seller, escrow.Amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log escrow release: %v", err)
	}

	fmt.Printf("Escrow %s released, funds transferred to seller %s\n", escrowID, escrow.Seller)
	return nil
}

// DisputeEscrow marks an escrow as disputed
func (crm *common.ComputerResourceMarketplace) DisputeEscrow(escrowID string) error {
	crm.mu.Lock()
	defer crm.mu.Unlock()

	// Retrieve the escrow
	escrow, exists := crm.Escrows[escrowID]
	if !exists || escrow.IsReleased {
		return fmt.Errorf("escrow %s not found or already released", escrowID)
	}

	// Mark the escrow as disputed
	escrow.IsDisputed = true

	// Log the dispute in the ledger
	err := crm.Ledger.RecordEscrowDispute(escrowID, escrow.Buyer, escrow.Seller, escrow.Amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log escrow dispute: %v", err)
	}

	fmt.Printf("Escrow %s disputed by buyer %s\n", escrowID, escrow.Buyer)
	return nil
}

