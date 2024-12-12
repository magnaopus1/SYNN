package marketplace

import (
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// VerifyStakeDistribution verifies that the stake distribution is correctly assigned.
func VerifyStakeDistribution(nftID string) (bool, error) {
    ledger := Ledger{}
    distribution, err := ledger.GetStakeDistribution(nftID)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve stake distribution: %v", err)
    }
    return distribution.IsValid, nil
}

// EnableCrossMarketplaceTrade enables trading of NFTs across multiple marketplaces.
func EnableCrossMarketplaceTrade() error {
    ledger := Ledger{}
    return ledger.SetCrossMarketplaceTradeStatus(true)
}

// SetCrossMarketplaceRate sets the exchange rate for cross-marketplace NFT trades.
func SetCrossMarketplaceRate(nftID string, rate float64) error {
    ledger := Ledger{}
    exchangeRate := CrossMarketplaceRate{
        NFTID:     nftID,
        Rate:      rate,
        Timestamp: time.Now(),
    }
    return ledger.RecordCrossMarketplaceRate(exchangeRate)
}

// ApproveCrossMarketplaceListing approves an NFT for listing across marketplaces.
func ApproveCrossMarketplaceListing(nftID string) error {
    ledger := Ledger{}
    return ledger.UpdateCrossMarketplaceListingStatus(nftID, "Approved")
}

// LogCrossMarketplaceTrade logs a trade that occurs across multiple marketplaces.
func LogCrossMarketplaceTrade(trade CrossMarketplaceTrade) error {
    ledger := Ledger{}
    trade.Timestamp = time.Now()
    return ledger.RecordCrossMarketplaceTrade(trade)
}

// TrackCrossMarketplaceMetrics tracks metrics associated with cross-marketplace trades.
func TrackCrossMarketplaceMetrics(nftID string, metrics CrossMarketplaceMetrics) error {
    ledger := Ledger{}
    metrics.NFTID = nftID
    metrics.Timestamp = time.Now()
    return ledger.RecordCrossMarketplaceMetrics(metrics)
}

// VerifyCrossMarketplaceStatus checks if an NFT is listed and tradable across marketplaces.
func VerifyCrossMarketplaceStatus(nftID string) (bool, error) {
    ledger := Ledger{}
    status, err := ledger.GetCrossMarketplaceStatus(nftID)
    if err != nil {
        return false, fmt.Errorf("failed to verify cross-marketplace status: %v", err)
    }
    return status.Listed, nil
}

// EnableUserRatingSystem enables the rating system for NFT transactions.
func EnableUserRatingSystem() error {
    ledger := Ledger{}
    return ledger.SetUserRatingSystemStatus(true)
}

// SetUserRatingForNFT assigns a rating from a user for a specific NFT.
func SetUserRatingForNFT(nftID string, userID string, rating int) error {
    ledger := Ledger{}
    userRating := UserRating{
        NFTID:     nftID,
        UserID:    userID,
        Rating:    rating,
        Timestamp: time.Now(),
    }
    return ledger.RecordUserRating(userRating)
}

// TrackUserFeedback records feedback from a user about an NFT transaction.
func TrackUserFeedback(nftID string, userID string, feedback string) error {
    ledger := Ledger{}
    encryptedFeedback, err := encryption.EncryptData(feedback)
    if err != nil {
        return fmt.Errorf("failed to encrypt feedback: %v", err)
    }
    userFeedback := UserFeedback{
        NFTID:     nftID,
        UserID:    userID,
        Feedback:  encryptedFeedback,
        Timestamp: time.Now(),
    }
    return ledger.RecordUserFeedback(userFeedback)
}

// GenerateRatingSummary generates a summary of user ratings for a specific NFT.
func GenerateRatingSummary(nftID string) (RatingSummary, error) {
    ledger := Ledger{}
    return ledger.GetRatingSummary(nftID)
}

// LogUserRatingActivity logs a user's activity related to rating an NFT.
func LogUserRatingActivity(userID string, activity RatingActivity) error {
    ledger := Ledger{}
    activity.UserID = userID
    activity.Timestamp = time.Now()
    return ledger.RecordUserRatingActivity(activity)
}


// EnableNFTInheritance enables inheritance rights for an NFT.
func EnableNFTInheritance(nftID string) error {
    ledger := Ledger{}
    return ledger.SetNFTInheritanceStatus(nftID, true)
}

// SetNFTInheritanceRights assigns inheritance rights for an NFT.
func SetNFTInheritanceRights(nftID string, beneficiaryID string) error {
    ledger := Ledger{}
    inheritance := NFTInheritanceRights{
        NFTID:         nftID,
        BeneficiaryID: beneficiaryID,
        SetAt:         time.Now(),
    }
    return ledger.RecordNFTInheritanceRights(inheritance)
}

// VerifyInheritanceClaim verifies a claim to an inherited NFT.
func VerifyInheritanceClaim(nftID string, claimantID string) (bool, error) {
    ledger := Ledger{}
    inheritance, err := ledger.GetNFTInheritanceRights(nftID)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve inheritance rights: %v", err)
    }
    return inheritance.BeneficiaryID == claimantID, nil
}

// LogInheritanceActivity logs activities related to inheritance claims for an NFT.
func LogInheritanceActivity(nftID string, activityDescription string) error {
    ledger := Ledger{}
    encryptedDescription, err := encryption.EncryptData(activityDescription)
    if err != nil {
        return fmt.Errorf("failed to encrypt inheritance activity description: %v", err)
    }
    activity := InheritanceActivity{
        NFTID:       nftID,
        Description: encryptedDescription,
        Timestamp:   time.Now(),
    }
    return ledger.RecordInheritanceActivity(activity)
}

// GenerateInheritanceReport generates a report of inheritance activities for a specified NFT.
func GenerateInheritanceReport(nftID string, from, to time.Time) (InheritanceReport, error) {
    ledger := Ledger{}
    activities, err := ledger.GetInheritanceActivities(nftID, from, to)
    if err != nil {
        return InheritanceReport{}, fmt.Errorf("failed to retrieve inheritance activities: %v", err)
    }
    return InheritanceReport{
        NFTID:      nftID,
        Activities: activities,
        Period:     fmt.Sprintf("%s - %s", from.Format(time.RFC3339), to.Format(time.RFC3339)),
    }, nil
}

// EnableNFTBundleListing enables the listing of NFT bundles in the marketplace.
func EnableNFTBundleListing() error {
    ledger := Ledger{}
    return ledger.SetNFTBundleListingStatus(true)
}

// CreateNFTBundle initializes a new NFT bundle for grouping multiple NFTs.
func CreateNFTBundle(bundle NFTBundle) error {
    ledger := Ledger{}
    bundle.CreatedAt = time.Now()
    return ledger.RecordNFTBundleCreation(bundle)
}

// AddNFTToBundle adds an NFT to an existing bundle.
func AddNFTToBundle(bundleID string, nftID string) error {
    ledger := Ledger{}
    nftBundle := NFTBundleEntry{
        BundleID: bundleID,
        NFTID:    nftID,
        AddedAt:  time.Now(),
    }
    return ledger.AddNFTToBundle(nftBundle)
}
