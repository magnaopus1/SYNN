package marketplace

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// GeneralMarketplace manages the decentralized listing, buying, leasing, and escrow
type GeneralMarketplace struct {
    Stores            map[string]*Store         // Decentralized store entries
    Listings          map[string]*Listing       // Marketplace listings
    Escrows           map[string]*Escrow        // Active escrows
    Ledger            *ledger.Ledger           // Blockchain ledger instance
    EncryptionService *encryption.Encryption   // Encryption service for data security
    mu                sync.Mutex               // Mutex for concurrency control
}

type Store struct {
    Owner            string                  // Wallet address of the store owner
    Name             string                  // Store name
    Description      string                  // Description of the store
    Category         string                  // Category of the store
    Listings         map[string]*Listing    // Items listed in the store
    Escrows          map[string]*Escrow     // Escrows associated with the store
    Ledger           *ledger.Ledger         // Blockchain ledger instance
    EncryptionService *encryption.Encryption // Encryption service for securing data
    mu               sync.Mutex             // Mutex for concurrency control
}

// Listing represents a product or service listed on the decentralized marketplace
type Listing struct {
    ListingID         string    // Unique identifier
    ItemName          string    // Item or service name
    Description       string    // Item description
    SalePrice         float64   // Sale price
    RentalPrices      map[string]float64 // Rental prices (per day, month, quarter, year)
    Owner             string    // Wallet address of the owner
    Available         bool      // Availability for sale or lease
    Category          string    // Item category
    ListedTime        time.Time // Listing timestamp
    Leasing           bool      // Flag for leasing
    IsLegal           bool      // Flag for compliance
    PhysicalOrDigital string    // Type of item (physical/digital)
    AvailableToRent   bool      // Availability for rent
    Tags              []string  // Tags for better categorization
}


// NewGeneralMarketplace initializes the general decentralized marketplace
func NewGeneralMarketplace(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) (*GeneralMarketplace, error) {
    // Validate dependencies
    if ledgerInstance == nil {
        return nil, fmt.Errorf("ledger instance cannot be nil")
    }
    if encryptionService == nil {
        return nil, fmt.Errorf("encryption service cannot be nil")
    }

    // Initialize the GeneralMarketplace instance
    marketplace := &GeneralMarketplace{
        Listings:          make(map[string]*Listing),
        Escrows:           make(map[string]*Escrow),
        Stores:            make(map[string]*Store),
        Ledger:            ledgerInstance,
        EncryptionService: encryptionService,
    }

    log.Printf("[INFO] General Marketplace initialized successfully")
    return marketplace, nil
}

//Create a General Marketplace
func CreateGeneralMarketplace(
    name string,
    description string,
    owner string,
    category string,
    ledgerInstance *ledger.Ledger,
    encryptionService *encryption.Encryption,
) (*Store, error) {
    if name == "" || owner == "" || ledgerInstance == nil || encryptionService == nil {
        return nil, fmt.Errorf("all parameters (name, owner, ledgerInstance, encryptionService) are required")
    }

    // Generate a unique store ID
    storeID := fmt.Sprintf("store-%s-%d", owner, time.Now().UnixNano())

    // Encrypt the store data
    storeData := fmt.Sprintf("StoreID: %s, Name: %s, Owner: %s", storeID, name, owner)
    encryptedData, err := encryptionService.EncryptData("AES", []byte(storeData), encryptionService.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt store data: %w", err)
    }

    // Create the Store instance
    store := &Store{
        Owner:             owner,
        Name:              name,
        Description:       description,
        Category:          category,
        Listings:          make(map[string]*Listing),
        Escrows:           make(map[string]*Escrow),
        Ledger:            ledgerInstance,
        EncryptionService: encryptionService,
    }

    // Log the store creation in the blockchain ledger
    err = ledgerInstance.RecordStoreCreation(storeID, owner, name, category)
    if err != nil {
        return nil, fmt.Errorf("failed to log store creation: %w", err)
    }

    log.Printf("[INFO] Store %s created successfully by %s", name, owner)
    return store, nil
}


func (l *ledger.Ledger) RecordStoreCreation(storeID, owner, name, category string) error {
    if storeID == "" || owner == "" || name == "" || category == "" {
        return fmt.Errorf("all parameters are required for store creation logging")
    }

    logEntry := fmt.Sprintf("Store Created: ID=%s, Owner=%s, Name=%s, Category=%s", storeID, owner, name, category)
    err := l.Log(logEntry)
    if err != nil {
        return fmt.Errorf("failed to log store creation: %w", err)
    }

    log.Printf("[LEDGER] Store creation logged successfully: %s", logEntry)
    return nil
}


// ListNewItem allows a user to list an item or service on the marketplace
func (gm *common.GeneralMarketplace) ListNewItem(listingID, itemName, description, owner, category string, price, leasePrice float64, leasing bool) (*common.Listing, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Check for legality (assumed implemented as part of legal compliance checks)
	if !isLegalListing(itemName, description) {
		return nil, errors.New("listing does not comply with legal standards")
	}

	// Encrypt listing data
	listingData := fmt.Sprintf("ListingID: %s, Owner: %s, ItemName: %s", listingID, owner, itemName)
	encryptedData, err := gm.EncryptionService.EncryptData([]byte(listingData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt listing data: %v", err)
	}

	// Create the listing
	listing := &common.Listing{
		ListingID:     listingID,
		ItemName:      itemName,
		Description:   description,
		Price:         price,
		LeasePrice:    leasePrice,
		Owner:         owner,
		Available:     true,
		Category:      category,
		ListedTime:    time.Now(),
		Leasing:       leasing,
		IsLegal:       true,
	}

	// Add the listing to the marketplace
	gm.Listings[listingID] = listing

	// Log the listing in the ledger
	err = gm.Ledger.RecordNewListing(listingID, itemName, owner, category, price, leasePrice, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log listing: %v", err)
	}

	fmt.Printf("Item %s listed by %s with price %f\n", itemName, owner, price)
	return listing, nil
}

// BuyItem allows a user to buy an item listed on the marketplace
func (gm *common.GeneralMarketplace) BuyItem(listingID, buyer string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Retrieve the listing
	listing, exists := gm.Listings[listingID]
	if !exists || !listing.Available {
		return fmt.Errorf("listing %s is not available", listingID)
	}

	// Create an escrow for the purchase
	escrowID := generateUniqueID()
	escrow := &common.Escrow{
		EscrowID:        escrowID,
		Buyer:           buyer,
		Seller:          listing.Owner,
		Amount:          listing.Price,
		ListingID:       listingID,
		TransactionType: "Buy",
		IsReleased:      false,
		IsDisputed:      false,
	}

	// Add the escrow to the system
	gm.Escrows[escrowID] = escrow

	// Mark the item as sold (not available)
	listing.Available = false

	// Log the purchase in the ledger
	err := gm.Ledger.RecordItemPurchase(listingID, buyer, listing.Price, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log item purchase: %v", err)
	}

	fmt.Printf("Item %s purchased by %s for %f\n", listing.ItemName, buyer, listing.Price)
	return nil
}

// LeaseItem allows a user to lease an item listed on the marketplace
func (gm *common.GeneralMarketplace) LeaseItem(listingID, lessee string, leaseDuration time.Duration) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Retrieve the listing
	listing, exists := gm.Listings[listingID]
	if !exists || !listing.Available || !listing.Leasing {
		return fmt.Errorf("listing %s is not available for lease", listingID)
	}

	// Calculate lease amount
	leaseAmount := listing.LeasePrice * float64(leaseDuration.Hours()) / 24.0

	// Create an escrow for the lease
	escrowID := generateUniqueID()
	escrow := &common.Escrow{
		EscrowID:        escrowID,
		Buyer:           lessee,
		Seller:          listing.Owner,
		Amount:          leaseAmount,
		ListingID:       listingID,
		TransactionType: "Lease",
		IsReleased:      false,
		IsDisputed:      false,
	}

	// Add the escrow to the system
	gm.Escrows[escrowID] = escrow

	// Mark the item as unavailable during the lease
	listing.Available = false

	// Log the lease in the ledger
	err := gm.Ledger.RecordItemLease(listingID, lessee, leaseAmount, leaseDuration, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log item lease: %v", err)
	}

	fmt.Printf("Item %s leased by %s for %f for %v duration\n", listing.ItemName, lessee, leaseAmount, leaseDuration)
	return nil
}

// ReleaseEscrow releases funds from escrow after a transaction is completed
func (gm *common.GeneralMarketplace) ReleaseEscrow(escrowID string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Retrieve the escrow
	escrow, exists := gm.Escrows[escrowID]
	if !exists || escrow.IsReleased {
		return fmt.Errorf("escrow %s not found or already released", escrowID)
	}

	// Mark the escrow as released
	escrow.IsReleased = true
	escrow.CompletionTime = time.Now()

	// Log the escrow release in the ledger
	err := gm.Ledger.RecordEscrowRelease(escrowID, escrow.Buyer, escrow.Seller, escrow.Amount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log escrow release: %v", err)
	}

	fmt.Printf("Escrow %s released, funds transferred to seller %s\n", escrowID, escrow.Seller)
	return nil
}

// SearchListings allows users to search and filter available listings
func (gm *common.GeneralMarketplace) SearchListings(query, category string) ([]*common.Listing, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	var results []*common.Listing
	for _, listing := range gm.Listings {
		if listing.Available && (strings.Contains(strings.ToLower(listing.ItemName), strings.ToLower(query)) || strings.Contains(strings.ToLower(listing.Description), strings.ToLower(query))) {
			if category == "" || strings.ToLower(listing.Category) == strings.ToLower(category) {
				results = append(results, listing)
			}
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no listings found matching the query")
	}

	return results, nil
}

// ReportIllegalItem allows users to report listings that violate laws or marketplace rules
func (gm *common.GeneralMarketplace) ReportIllegalItem(listingID, reporter string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Retrieve the listing
	listing, exists := gm.Listings[listingID]
	if !exists {
		return fmt.Errorf("listing %s not found", listingID)
	}

	// Mark the listing as illegal
	listing.IsLegal = false
	listing.Available = false

	// Log the report in the ledger
	err := gm.Ledger.RecordIllegalItemReport(listingID, reporter, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log illegal item report: %v", err)
	}

	fmt.Printf("Listing %s reported by %s and marked as illegal\n", listingID, reporter)
	return nil
}

// RemoveIllegalItem removes a reported illegal item from the marketplace
func (gm *common.GeneralMarketplace) RemoveIllegalItem(listingID string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Retrieve the listing
	listing, exists := gm.Listings[listingID]
	if !exists {
		return fmt.Errorf("listing %s not found", listingID)
	}

	// Check if the listing is marked as illegal
	if listing.IsLegal {
		return fmt.Errorf("listing %s is not marked as illegal", listingID)
	}

	// Remove the illegal item from the marketplace
	delete(gm.Listings, listingID)

	// Log the removal of the illegal item in the ledger
	err := gm.Ledger.RecordItemRemoval(listingID, "Illegal Item Removal", time.Now())
	if err != nil {
		return fmt.Errorf("failed to log illegal item removal: %v", err)
	}

	fmt.Printf("Illegal listing %s has been removed from the marketplace\n", listingID)
	return nil
}

// RemoveItem allows the owner to remove their item from the marketplace
func (gm *common.GeneralMarketplace) RemoveItem(listingID, owner string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Retrieve the listing
	listing, exists := gm.Listings[listingID]
	if !exists {
		return fmt.Errorf("listing %s not found", listingID)
	}

	// Ensure only the owner can remove the item
	if listing.Owner != owner {
		return fmt.Errorf("you are not the owner of this listing")
	}

	// Remove the item from the marketplace
	delete(gm.Listings, listingID)

	// Log the item removal in the ledger
	err := gm.Ledger.RecordItemRemoval(listingID, "Owner Removal", time.Now())
	if err != nil {
		return fmt.Errorf("failed to log item removal: %v", err)
	}

	fmt.Printf("Listing %s has been removed by the owner %s\n", listingID, owner)
	return nil
}


// isLegalListing simulates a legal compliance check for items or services listed
func isLegalListing(itemName, description string) bool {
	// For the sake of this demo, we'll just check if the item is flagged with illegal keywords
	illegalKeywords := []string{"illegal", "contraband", "banned"}
	for _, keyword := range illegalKeywords {
		if strings.Contains(strings.ToLower(itemName), keyword) || strings.Contains(strings.ToLower(description), keyword) {
			return false
		}
	}
	return true
}
