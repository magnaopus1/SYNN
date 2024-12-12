package marketplace

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// TrackEscrowStatus logs the status of an escrow account for an NFT.
func TrackEscrowStatus(nftID string, status EscrowStatus) error {
    ledger := Ledger{}
    status.NFTID = nftID
    status.Timestamp = time.Now()
    return ledger.RecordEscrowStatus(status)
}

// GenerateEscrowReport generates a report detailing escrow activity for a specific NFT.
func GenerateEscrowReport(nftID string, from, to time.Time) (EscrowReport, error) {
    ledger := Ledger{}
    statuses, err := ledger.GetEscrowHistory(nftID, from, to)
    if err != nil {
        return EscrowReport{}, fmt.Errorf("failed to retrieve escrow history: %v", err)
    }
    return EscrowReport{
        NFTID:    nftID,
        Statuses: statuses,
        Period:   fmt.Sprintf("%s - %s", from.Format(time.RFC3339), to.Format(time.RFC3339)),
    }, nil
}

// EnableNFTRental enables the rental option for an NFT.
func EnableNFTRental(nftID string) error {
    ledger := Ledger{}
    return ledger.SetNFTRentalStatus(nftID, true)
}

// SetNFTRentalTerms sets the terms for renting an NFT.
func SetNFTRentalTerms(nftID string, terms NFTRentalTerms) error {
    ledger := Ledger{}
    terms.NFTID = nftID
    terms.SetAt = time.Now()
    return ledger.RecordNFTRentalTerms(terms)
}

// TrackRentalPayments logs each payment made towards the rental of an NFT.
func TrackRentalPayments(nftID string, payment RentalPayment) error {
    ledger := Ledger{}
    payment.NFTID = nftID
    payment.Timestamp = time.Now()
    return ledger.RecordRentalPayment(payment)
}

// LogRentalActivity logs rental-related activity for an NFT.
func LogRentalActivity(nftID string, activityDescription string) error {
    ledger := Ledger{}
    encryptedDescription, err := encryption.EncryptData(activityDescription)
    if err != nil {
        return fmt.Errorf("failed to encrypt rental activity description: %v", err)
    }
    activity := RentalActivity{
        NFTID:       nftID,
        Description: encryptedDescription,
        Timestamp:   time.Now(),
    }
    return ledger.RecordRentalActivity(activity)
}

// VerifyRentalContract verifies the rental contract for an NFT.
func VerifyRentalContract(contractID string) (bool, error) {
    ledger := Ledger{}
    contract, err := ledger.GetRentalContract(contractID)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve rental contract: %v", err)
    }
    return contract.IsValid, nil
}

// EnableNFTVerificationBadge enables verification badges for NFTs.
func EnableNFTVerificationBadge() error {
    ledger := Ledger{}
    return ledger.SetNFTVerificationBadgeStatus(true)
}

// GrantNFTVerificationBadge assigns a verification badge to an NFT.
func GrantNFTVerificationBadge(nftID string) error {
    ledger := Ledger{}
    badge := VerificationBadge{
        NFTID:     nftID,
        GrantedAt: time.Now(),
    }
    return ledger.RecordVerificationBadgeGrant(badge)
}

// RevokeNFTVerificationBadge removes a verification badge from an NFT.
func RevokeNFTVerificationBadge(nftID string) error {
    ledger := Ledger{}
    badge := VerificationBadge{
        NFTID:     nftID,
        RevokedAt: time.Now(),
    }
    return ledger.RecordVerificationBadgeRevoke(badge)
}

// TrackBadgeHistory records the history of verification badges for an NFT.
func TrackBadgeHistory(nftID string) ([]VerificationBadge, error) {
    ledger := Ledger{}
    return ledger.GetBadgeHistory(nftID)
}

// EnableNFTCollection enables the creation and management of NFT collections.
func EnableNFTCollection() error {
    ledger := Ledger{}
    return ledger.SetNFTCollectionStatus(true)
}

// AddToNFTCollection adds an NFT to a specific collection.
func AddToNFTCollection(collectionID string, nftID string) error {
    ledger := Ledger{}
    entry := NFTCollectionEntry{
        CollectionID: collectionID,
        NFTID:        nftID,
        AddedAt:      time.Now(),
    }
    return ledger.AddToNFTCollection(entry)
}

// RemoveFromNFTCollection removes an NFT from a specific collection.
func RemoveFromNFTCollection(collectionID string, nftID string) error {
    ledger := Ledger{}
    return ledger.RemoveFromNFTCollection(collectionID, nftID)
}

// TrackCollectionOwnership logs ownership changes for an NFT within a collection.
func TrackCollectionOwnership(collectionID string, nftID string, newOwnerID string) error {
    ledger := Ledger{}
    ownershipChange := CollectionOwnershipChange{
        CollectionID: collectionID,
        NFTID:        nftID,
        NewOwnerID:   newOwnerID,
        Timestamp:    time.Now(),
    }
    return ledger.RecordCollectionOwnershipChange(ownershipChange)
}

// LogCollectionActivity records an activity related to a specific NFT collection.
func LogCollectionActivity(collectionID string, activityDescription string) error {
    ledger := Ledger{}
    encryptedDescription, err := encryption.EncryptData(activityDescription)
    if err != nil {
        return fmt.Errorf("failed to encrypt collection activity description: %v", err)
    }
    activity := CollectionActivity{
        CollectionID: collectionID,
        Description:  encryptedDescription,
        Timestamp:    time.Now(),
    }
    return ledger.RecordCollectionActivity(activity)
}

// GenerateNFTCollectionReport generates a report of all activities for an NFT collection.
func GenerateNFTCollectionReport(collectionID string, from, to time.Time) (CollectionReport, error) {
    ledger := Ledger{}
    activities, err := ledger.GetCollectionActivities(collectionID, from, to)
    if err != nil {
        return CollectionReport{}, fmt.Errorf("failed to retrieve collection activities: %v", err)
    }
    return CollectionReport{
        CollectionID: collectionID,
        Activities:   activities,
        Period:       fmt.Sprintf("%s - %s", from.Format(time.RFC3339), to.Format(time.RFC3339)),
    }, nil
}

// EnableNFTTrading enables trading of NFTs within the marketplace.
func EnableNFTTrading() error {
    ledger := Ledger{}
    return ledger.SetNFTTradingStatus(true)
}

// LogNFTTradeEvent logs a specific trade event for an NFT.
func LogNFTTradeEvent(trade NFTTradeEvent) error {
    ledger := Ledger{}
    trade.Timestamp = time.Now()
    return ledger.RecordNFTTradeEvent(trade)
}

// ApproveNFTTrade approves a pending NFT trade.
func ApproveNFTTrade(tradeID string) error {
    ledger := Ledger{}
    return ledger.UpdateNFTTradeStatus(tradeID, "Approved")
}
