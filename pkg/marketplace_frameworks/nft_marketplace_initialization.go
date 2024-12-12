package marketplace

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// InitializeNFTMarketplace sets up the NFT marketplace with initial configuration.
func InitializeNFTMarketplace(config NFTMarketplaceConfig) error {
    ledger := Ledger{}
    config.InitializedAt = time.Now()
    return ledger.RecordNFTMarketplaceInitialization(config)
}

// MintNFT mints a new NFT with specified metadata.
func MintNFT(nft NFT) error {
    ledger := Ledger{}
    nft.MintedAt = time.Now()
    return ledger.RecordNFTMinting(nft)
}

// BurnNFT burns an existing NFT, removing it from the ledger.
func BurnNFT(nftID string) error {
    ledger := Ledger{}
    burnRecord := NFTBurnRecord{
        NFTID:    nftID,
        BurnedAt: time.Now(),
    }
    return ledger.RecordNFTBurn(burnRecord)
}

// TransferNFTOwnership transfers ownership of an NFT to a new owner.
func TransferNFTOwnership(nftID string, newOwnerID string) error {
    ledger := Ledger{}
    transfer := NFTOwnershipTransfer{
        NFTID:         nftID,
        NewOwnerID:    newOwnerID,
        TransferredAt: time.Now(),
    }
    return ledger.RecordNFTOwnershipTransfer(transfer)
}

// ApproveNFTTransfer approves the transfer of an NFT to a new owner.
func ApproveNFTTransfer(nftID string, approvedOwnerID string) error {
    ledger := Ledger{}
    return ledger.SetNFTTransferApproval(nftID, approvedOwnerID, true)
}

// DenyNFTTransfer denies a pending transfer of an NFT.
func DenyNFTTransfer(nftID string) error {
    ledger := Ledger{}
    return ledger.SetNFTTransferApproval(nftID, "", false)
}

// ListNFTForSale lists an NFT for sale in the marketplace.
func ListNFTForSale(nftID string, saleDetails NFTSale) error {
    ledger := Ledger{}
    saleDetails.NFTID = nftID
    saleDetails.ListedAt = time.Now()
    return ledger.RecordNFTListingForSale(saleDetails)
}

// RemoveNFTFromSale removes an NFT from the marketplace sale listing.
func RemoveNFTFromSale(nftID string) error {
    ledger := Ledger{}
    return ledger.RemoveNFTFromSale(nftID)
}

// SetNFTSalePrice sets the sale price for an NFT.
func SetNFTSalePrice(nftID string, price float64) error {
    ledger := Ledger{}
    saleDetails := NFTSale{
        NFTID:    nftID,
        Price:    price,
        UpdatedAt: time.Now(),
    }
    return ledger.UpdateNFTSalePrice(saleDetails)
}

// BidOnNFT places a bid on an NFT that is listed for sale or auction.
func BidOnNFT(nftID string, bid NFTBid) error {
    ledger := Ledger{}
    bid.NFTID = nftID
    bid.Timestamp = time.Now()
    return ledger.RecordNFTBid(bid)
}

// AcceptNFTBid accepts a bid for an NFT.
func AcceptNFTBid(bidID string) error {
    ledger := Ledger{}
    return ledger.UpdateNFTBidStatus(bidID, "Accepted")
}

// RejectNFTBid rejects a bid for an NFT.
func RejectNFTBid(bidID string) error {
    ledger := Ledger{}
    return ledger.UpdateNFTBidStatus(bidID, "Rejected")
}

// EnableAuctionForNFT enables an auction for a specific NFT.
func EnableAuctionForNFT(nftID string, auctionDetails NFTAuction) error {
    ledger := Ledger{}
    auctionDetails.NFTID = nftID
    auctionDetails.StartedAt = time.Now()
    return ledger.RecordNFTAuctionStart(auctionDetails)
}

// EndAuctionForNFT concludes an active auction for an NFT.
func EndAuctionForNFT(auctionID string) error {
    ledger := Ledger{}
    auctionEnd := NFTAuctionEnd{
        AuctionID: auctionID,
        EndedAt:   time.Now(),
    }
    return ledger.RecordNFTAuctionEnd(auctionEnd)
}

// TrackNFTAuctionStatus retrieves the current status of an NFT auction.
func TrackNFTAuctionStatus(auctionID string) (NFTAuctionStatus, error) {
    ledger := Ledger{}
    return ledger.GetNFTAuctionStatus(auctionID)
}


// LogNFTAuctionEvent logs an event related to an NFT auction.
func LogNFTAuctionEvent(auctionID string, eventDescription string) error {
    ledger := Ledger{}
    encryptedDescription, err := encryption.EncryptData(eventDescription)
    if err != nil {
        return fmt.Errorf("failed to encrypt auction event description: %v", err)
    }
    event := NFTAuctionEvent{
        AuctionID:   auctionID,
        Description: encryptedDescription,
        Timestamp:   time.Now(),
    }
    return ledger.RecordNFTAuctionEvent(event)
}

// GenerateAuctionReport generates a report of all events and bids for a specified auction period.
func GenerateAuctionReport(auctionID string, from, to time.Time) (NFTAuctionReport, error) {
    ledger := Ledger{}
    events, err := ledger.GetNFTAuctionEvents(auctionID, from, to)
    if err != nil {
        return NFTAuctionReport{}, fmt.Errorf("failed to retrieve auction events: %v", err)
    }
    bids, err := ledger.GetNFTAuctionBids(auctionID, from, to)
    if err != nil {
        return NFTAuctionReport{}, fmt.Errorf("failed to retrieve auction bids: %v", err)
    }
    return NFTAuctionReport{
        AuctionID: auctionID,
        Events:    events,
        Bids:      bids,
        Period:    fmt.Sprintf("%s - %s", from.Format(time.RFC3339), to.Format(time.RFC3339)),
    }, nil
}

// CheckNFTOwnership checks if a specific user owns a particular NFT.
func CheckNFTOwnership(nftID string, userID string) (bool, error) {
    ledger := Ledger{}
    ownership, err := ledger.GetNFTOwner(nftID)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve NFT ownership: %v", err)
    }
    return ownership.OwnerID == userID, nil
}

// GetNFTMetadata retrieves the metadata of a specified NFT.
func GetNFTMetadata(nftID string) (NFTMetadata, error) {
    ledger := Ledger{}
    return ledger.GetNFTMetadata(nftID)
}

// UpdateNFTMetadata updates the metadata of a specified NFT.
func UpdateNFTMetadata(nftID string, newMetadata NFTMetadata) error {
    ledger := Ledger{}
    newMetadata.NFTID = nftID
    newMetadata.UpdatedAt = time.Now()
    return ledger.UpdateNFTMetadata(newMetadata)
}
