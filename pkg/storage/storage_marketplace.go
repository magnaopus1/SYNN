package storage

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)


// NewStorageMarketplace initializes a new StorageMarketplace
func NewStorageMarketplace(ledgerInstance *ledger.Ledger) *StorageMarketplace {
    return &StorageMarketplace{
        Listings:       make(map[string]*StorageListing),
        LedgerInstance: ledgerInstance,
        EscrowAccounts: make(map[string]*EscrowAccount),
    }
}

// PostStorageListing allows a user to post a new storage listing to the marketplace
func (sm *StorageMarketplace) PostStorageListing(owner string, capacityGB int, pricePerGB float64, leaseDuration int, details string) (string, error) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    // Generate a unique ListingID for the storage
    listingID := sm.generateListingID(owner, capacityGB)

    // Encrypt the details of the listing
    encryptedDetails, err := encryption.EncryptData(details, common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt listing details: %v", err)
    }

    // Create the new storage listing
    newListing := &StorageListing{
        ListingID:       listingID,
        Owner:           owner,
        CapacityGB:      capacityGB,
        PricePerGB:      pricePerGB,
        LeaseDuration:   leaseDuration,
        PostedAt:        time.Now(),
        EncryptedDetails: encryptedDetails,
        Active:          true,
    }

    // Add the listing to the marketplace
    sm.Listings[listingID] = newListing

    // Log the listing to the ledger as a transaction
    err = sm.logListingToLedger(newListing, "post")
    if err != nil {
        return "", fmt.Errorf("failed to log storage listing to ledger: %v", err)
    }

    fmt.Printf("Storage listing %s posted by %s.\n", listingID, owner)
    return listingID, nil
}

// PurchaseStorage allows a user to purchase or lease storage from the marketplace using escrow
func (sm *StorageMarketplace) PurchaseStorage(listingID, buyer string, capacityGB int) (string, error) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    listing, exists := sm.Listings[listingID]
    if !exists || !listing.Active {
        return "", fmt.Errorf("listing not found or inactive")
    }

    // Check if the requested capacity exceeds the available capacity
    if capacityGB > listing.CapacityGB {
        return "", fmt.Errorf("requested capacity exceeds available capacity")
    }

    totalCost := float64(capacityGB) * listing.PricePerGB

    // Create an escrow account for this transaction
    escrowID, err := sm.createEscrowAccount(buyer, listing.Owner, totalCost)
    if err != nil {
        return "", fmt.Errorf("failed to create escrow account: %v", err)
    }

    // Deduct the requested capacity from the listing
    listing.CapacityGB -= capacityGB
    if listing.CapacityGB == 0 {
        listing.Active = false // Deactivate the listing if capacity is exhausted
    }

    fmt.Printf("Storage of %dGB purchased by %s from listing %s via escrow %s.\n", capacityGB, buyer, listingID, escrowID)
    return escrowID, nil
}

// ReleaseEscrow releases the funds in escrow to the seller
func (sm *StorageMarketplace) ReleaseEscrow(escrowID string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    escrow, exists := sm.EscrowAccounts[escrowID]
    if !exists || escrow.IsReleased {
        return fmt.Errorf("escrow account not found or already released")
    }

    // Log the release of escrow funds to the ledger
    err := sm.logEscrowReleaseToLedger(escrow)
    if err != nil {
        return fmt.Errorf("failed to log escrow release to ledger: %v", err)
    }

    // Mark the escrow funds as released
    escrow.IsReleased = true

    fmt.Printf("Escrow %s released to seller %s.\n", escrowID, escrow.Seller)
    return nil
}

// ViewListing allows users to view the details of a storage listing
func (sm *StorageMarketplace) ViewListing(listingID string) (*StorageListing, error) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    listing, exists := sm.Listings[listingID]
    if !exists {
        return nil, fmt.Errorf("listing not found")
    }

    fmt.Printf("Listing %s viewed.\n", listingID)
    return listing, nil
}

// SearchListings allows users to search for active storage listings that meet a specific capacity and price range
func (sm *StorageMarketplace) SearchListings(minCapacity int, maxPricePerGB float64) ([]*StorageListing, error) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    var matchingListings []*StorageListing

    for _, listing := range sm.Listings {
        if listing.Active && listing.CapacityGB >= minCapacity && listing.PricePerGB <= maxPricePerGB {
            matchingListings = append(matchingListings, listing)
        }
    }

    fmt.Printf("%d matching listings found.\n", len(matchingListings))
    return matchingListings, nil
}

// createEscrowAccount creates a new escrow account for a transaction
func (sm *StorageMarketplace) createEscrowAccount(buyer, seller string, amount float64) (string, error) {
    escrowID := sm.generateEscrowID(buyer, seller, amount)

    // Create and store the new escrow account
    escrowAccount := &EscrowAccount{
        EscrowID:   escrowID,
        Buyer:      buyer,
        Seller:     seller,
        Amount:     amount,
        CreatedAt:  time.Now(),
        IsReleased: false,
    }

    sm.EscrowAccounts[escrowID] = escrowAccount

    // Log the creation of the escrow account in the ledger
    err := sm.logEscrowCreationToLedger(escrowAccount)
    if err != nil {
        return "", fmt.Errorf("failed to log escrow account to ledger: %v", err)
    }

    fmt.Printf("Escrow account %s created for buyer %s and seller %s.\n", escrowID, buyer, seller)
    return escrowID, nil
}

// generateListingID creates a unique ListingID for each storage listing
func (sm *StorageMarketplace) generateListingID(owner string, capacityGB int) string {
    hashInput := fmt.Sprintf("%s%d%d", owner, capacityGB, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// generateEscrowID creates a unique EscrowID based on buyer, seller, and amount
func (sm *StorageMarketplace) generateEscrowID(buyer, seller string, amount float64) string {
    hashInput := fmt.Sprintf("%s%s%.2f%d", buyer, seller, amount, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// logListingToLedger logs the storage listing to the ledger as a transaction
func (sm *StorageMarketplace) logListingToLedger(listing *StorageListing, action string) error {
    listingDetails := fmt.Sprintf("ListingID: %s, Owner: %s, Capacity: %dGB, Price: %.2f SYNN/GB, LeaseDuration: %d days",
        listing.ListingID, listing.Owner, listing.CapacityGB, listing.PricePerGB, listing.LeaseDuration)

    encryptedDetails, err := encryption.EncryptData(listingDetails, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt listing details: %v", err)
    }

    err = sm.LedgerInstance.RecordTransaction(listing.ListingID, action, listing.Owner, encryptedDetails)
    if err != nil {
        return fmt.Errorf("failed to log listing to ledger: %v", err)
    }

    fmt.Printf("Listing %s logged to the ledger.\n", listing.ListingID)
    return nil
}

// logEscrowCreationToLedger logs the creation of an escrow account in the ledger
func (sm *StorageMarketplace) logEscrowCreationToLedger(escrow *EscrowAccount) error {
    escrowDetails := fmt.Sprintf("EscrowID: %s, Buyer: %s, Seller: %s, Amount: %.2f, CreatedAt: %s",
        escrow.EscrowID, escrow.Buyer, escrow.Seller, escrow.Amount, escrow.CreatedAt)

    encryptedDetails, err := encryption.EncryptData(escrowDetails, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt escrow details: %v", err)
    }

    err = sm.LedgerInstance.RecordTransaction(escrow.EscrowID, "escrow_create", escrow.Buyer, encryptedDetails)
    if err != nil {
        return fmt.Errorf("failed to log escrow creation to ledger: %v", err)
    }

    fmt.Printf("Escrow %s logged to the ledger.\n", escrow.EscrowID)
    return nil
}

// logEscrowReleaseToLedger logs the release of escrow funds to the ledger
func (sm *StorageMarketplace) logEscrowReleaseToLedger(escrow *EscrowAccount) error {
    escrowReleaseDetails := fmt.Sprintf("EscrowID: %s, Buyer: %s, Seller: %s, Amount: %.2f, ReleasedAt: %s",
        escrow.EscrowID, escrow.Buyer, escrow.Seller, escrow.Amount, time.Now())

    encryptedDetails, err := encryption.EncryptData(escrowReleaseDetails, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt escrow release details: %v", err)
    }

    err = sm.LedgerInstance.RecordTransaction(escrow.EscrowID, "escrow_release", escrow.Seller, encryptedDetails)
    if err != nil {
        return fmt.Errorf("failed to log escrow release to ledger: %v", err)
    }

    fmt.Printf("Escrow release for %s logged to the ledger.\n", escrow.EscrowID)
    return nil
}
