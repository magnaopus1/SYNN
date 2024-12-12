package marketplace

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewNFTMarketplace initializes the NFT marketplace
func NewNFTMarketplace(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.NFTMarketplace {
	return &common.NFTMarketplace{
		Listings:         make(map[string]*common.NFTListing),
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// ListNFT allows a user to list an NFT for sale in the marketplace
func (nm *common.NFTMarketplace) ListNFT(listingID, tokenID, standard, owner, metadataURI string, price float64) (*common.NFTListing, error) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	// Ensure that only Syn721 or Syn1155 standards are used
	if standard != "Syn721" && standard != "Syn1155" {
		return nil, fmt.Errorf("invalid token standard: %s", standard)
	}

	// Encrypt listing data
	listingData := fmt.Sprintf("ListingID: %s, TokenID: %s, Owner: %s, Price: %f", listingID, tokenID, owner, price)
	encryptedData, err := nm.EncryptionService.EncryptData([]byte(listingData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt listing data: %v", err)
	}

	// Create a new NFT listing
	listing := &common.NFTListing{
		ListingID:    listingID,
		TokenID:      tokenID,
		Standard:     standard,
		MetadataURI:  metadataURI,
		Price:        price,
		Owner:        owner,
		Available:    true,
		ListedTime:   time.Now(),
	}

	// Add the listing to the marketplace
	nm.Listings[listingID] = listing

	// Log the NFT listing in the ledger
	err = nm.Ledger.RecordNFTListing(listingID, tokenID, standard, owner, price, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log NFT listing: %v", err)
	}

	fmt.Printf("NFT %s (Token ID: %s) listed by %s for %f\n", standard, tokenID, owner, price)
	return listing, nil
}

// BuyNFT allows a user to purchase an NFT from the marketplace
func (nm *common.NFTMarketplace) BuyNFT(listingID, buyer string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	// Retrieve the listing
	listing, exists := nm.Listings[listingID]
	if !exists || !listing.Available {
		return fmt.Errorf("listing %s is not available", listingID)
	}

	// Log the purchase in the ledger
	err := nm.Ledger.RecordNFTPurchase(listingID, listing.TokenID, listing.Standard, buyer, listing.Owner, listing.Price, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log NFT purchase: %v", err)
	}

	// Transfer ownership to the buyer
	listing.Owner = buyer
	listing.Available = false

	fmt.Printf("NFT %s (Token ID: %s) purchased by %s for %f\n", listing.Standard, listing.TokenID, buyer, listing.Price)
	return nil
}

// CancelListing allows the owner to cancel their NFT listing
func (nm *common.NFTMarketplace) CancelListing(listingID, owner string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	// Retrieve the listing
	listing, exists := nm.Listings[listingID]
	if !exists {
		return fmt.Errorf("listing %s not found", listingID)
	}

	// Ensure only the owner can cancel the listing
	if listing.Owner != owner {
		return fmt.Errorf("you are not the owner of this listing")
	}

	// Mark the listing as unavailable
	listing.Available = false

	// Log the cancellation in the ledger
	err := nm.Ledger.RecordNFTListingCancellation(listingID, listing.TokenID, listing.Standard, owner, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log NFT listing cancellation: %v", err)
	}

	fmt.Printf("NFT %s (Token ID: %s) listing canceled by owner %s\n", listing.Standard, listing.TokenID, owner)
	return nil
}

// SearchNFTs allows users to search for NFTs based on keywords and filter by token standard
func (nm *common.NFTMarketplace) SearchNFTs(query, standard string) ([]*common.NFTListing, error) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	var results []*common.NFTListing
	for _, listing := range nm.Listings {
		if listing.Available && strings.Contains(strings.ToLower(listing.MetadataURI), strings.ToLower(query)) {
			if standard == "" || strings.ToLower(listing.Standard) == strings.ToLower(standard) {
				results = append(results, listing)
			}
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no NFTs found matching the query")
	}

	return results, nil
}


