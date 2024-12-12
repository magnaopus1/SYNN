package marketplace

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// VerifyNFTAuthenticity checks the authenticity of an NFT to confirm it is genuine.
func VerifyNFTAuthenticity(nftID string) (bool, error) {
    ledger := Ledger{}
    authenticity, err := ledger.GetNFTAuthenticity(nftID)
    if err != nil {
        return false, fmt.Errorf("failed to verify NFT authenticity: %v", err)
    }
    return authenticity.IsGenuine, nil
}

// TrackNFTOwnershipHistory logs and retrieves the history of ownership changes for an NFT.
func TrackNFTOwnershipHistory(nftID string) ([]NFTOwnershipHistory, error) {
    ledger := Ledger{}
    return ledger.GetNFTOwnershipHistory(nftID)
}

// RecordNFTTransferEvent logs the transfer of ownership of an NFT.
func RecordNFTTransferEvent(transfer NFTTransferEvent) error {
    ledger := Ledger{}
    transfer.Timestamp = time.Now()
    return ledger.RecordNFTTransferEvent(transfer)
}

// ApproveNFTListing approves a listing for an NFT in the marketplace.
func ApproveNFTListing(nftID string) error {
    ledger := Ledger{}
    return ledger.UpdateNFTListingStatus(nftID, "Approved")
}

// DenyNFTListing denies a listing request for an NFT in the marketplace.
func DenyNFTListing(nftID string) error {
    ledger := Ledger{}
    return ledger.UpdateNFTListingStatus(nftID, "Denied")
}

// ScheduleNFTListing schedules an NFT listing for a future time.
func ScheduleNFTListing(nftID string, listingTime time.Time) error {
    ledger := Ledger{}
    listing := NFTListing{
        NFTID:       nftID,
        ScheduledAt: listingTime,
        Status:      "Scheduled",
    }
    return ledger.RecordNFTListing(listing)
}

// EnableNFTStaking enables staking functionality for a specific NFT.
func EnableNFTStaking(nftID string) error {
    ledger := Ledger{}
    return ledger.SetNFTStakingStatus(nftID, true)
}

// StakeNFT stakes an NFT for a specified duration.
func StakeNFT(stake NFTStaking) error {
    ledger := Ledger{}
    stake.StakedAt = time.Now()
    return ledger.RecordNFTStake(stake)
}

// UnstakeNFT releases an NFT from staking.
func UnstakeNFT(nftID string, userID string) error {
    ledger := Ledger{}
    unstake := NFTUnstake{
        NFTID:      nftID,
        UserID:     userID,
        UnstakedAt: time.Now(),
    }
    return ledger.RecordNFTUnstake(unstake)
}

// LogNFTStakingEvent records events related to NFT staking.
func LogNFTStakingEvent(nftID string, eventDescription string) error {
    ledger := Ledger{}
    encryptedDescription, err := encryption.EncryptData(eventDescription)
    if err != nil {
        return fmt.Errorf("failed to encrypt staking event description: %v", err)
    }
    event := NFTStakingEvent{
        NFTID:       nftID,
        Description: encryptedDescription,
        Timestamp:   time.Now(),
    }
    return ledger.RecordNFTStakingEvent(event)
}

// SetNFTRoyaltyPercentage sets the royalty percentage for secondary sales of an NFT.
func SetNFTRoyaltyPercentage(nftID string, percentage float64) error {
    ledger := Ledger{}
    royalty := NFTRoyalty{
        NFTID:      nftID,
        Percentage: percentage,
        SetAt:      time.Now(),
    }
    return ledger.RecordNFTRoyalty(royalty)
}

// DistributeNFTRoyalties distributes royalties to the NFT creator on each secondary sale.
func DistributeNFTRoyalties(nftID string, saleAmount float64) error {
    ledger := Ledger{}
    royalty, err := ledger.GetNFTRoyalty(nftID)
    if err != nil {
        return fmt.Errorf("failed to retrieve NFT royalty information: %v", err)
    }
    royaltyAmount := saleAmount * royalty.Percentage / 100
    distribution := RoyaltyDistribution{
        NFTID:         nftID,
        Amount:        royaltyAmount,
        DistributedAt: time.Now(),
    }
    return ledger.RecordRoyaltyDistribution(distribution)
}

// GenerateRoyaltyReport generates a report on royalty distributions for a specified period.
func GenerateRoyaltyReport(nftID string, from, to time.Time) (RoyaltyReport, error) {
    ledger := Ledger{}
    distributions, err := ledger.GetRoyaltyDistributions(nftID, from, to)
    if err != nil {
        return RoyaltyReport{}, fmt.Errorf("failed to retrieve royalty distributions: %v", err)
    }
    return RoyaltyReport{
        NFTID:         nftID,
        Distributions: distributions,
        Period:        fmt.Sprintf("%s - %s", from.Format(time.RFC3339), to.Format(time.RFC3339)),
    }, nil
}


// VerifyNFTRoyaltyDistribution checks if the royalties have been properly distributed.
func VerifyNFTRoyaltyDistribution(nftID string, expectedAmount float64) (bool, error) {
    ledger := Ledger{}
    totalDistributed, err := ledger.GetTotalRoyaltyDistribution(nftID)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve total royalty distribution: %v", err)
    }
    return totalDistributed >= expectedAmount, nil
}

// EnableFractionalOwnershipForNFT enables fractional ownership for an NFT.
func EnableFractionalOwnershipForNFT(nftID string) error {
    ledger := Ledger{}
    return ledger.SetFractionalOwnershipStatus(nftID, true)
}

// TrackFractionalOwnership logs and retrieves details of fractional ownership for an NFT.
func TrackFractionalOwnership(nftID string) ([]FractionalOwnership, error) {
    ledger := Ledger{}
    return ledger.GetFractionalOwnershipDetails(nftID)
}

// LogFractionalOwnershipChange records changes to fractional ownership of an NFT.
func LogFractionalOwnershipChange(nftID string, ownershipChange FractionalOwnershipChange) error {
    ledger := Ledger{}
    ownershipChange.NFTID = nftID
    ownershipChange.Timestamp = time.Now()
    return ledger.RecordFractionalOwnershipChange(ownershipChange)
}

// EnableNFTEscrow enables escrow services for transactions involving the NFT.
func EnableNFTEscrow(nftID string) error {
    ledger := Ledger{}
    return ledger.SetNFTEscrowStatus(nftID, true)
}

// ReleaseNFTEscrow releases an NFT from escrow, completing the transaction.
func ReleaseNFTEscrow(escrowID string) error {
    ledger := Ledger{}
    release := NFTEscrowRelease{
        EscrowID:    escrowID,
        ReleasedAt:  time.Now(),
    }
    return ledger.RecordEscrowRelease(release)
}
