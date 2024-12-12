package sustainability

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewEnergyEfficiencyRatingSystem initializes a new energy efficiency rating system
func NewEnergyEfficiencyRatingSystem(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *EnergyEfficiencyRatingSystem {
	return &EnergyEfficiencyRatingSystem{
		Ratings:          make(map[string]*EfficiencyRating),
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// IssueRating issues a new energy efficiency rating to a node
func (eers *EnergyEfficiencyRatingSystem) IssueRating(ratingID, nodeID, owner, issuer string, rating, energyUsage float64, validityPeriod time.Duration) (*EfficiencyRating, error) {
	eers.mu.Lock()
	defer eers.mu.Unlock()

	// Encrypt rating data
	ratingData := fmt.Sprintf("RatingID: %s, NodeID: %s, Owner: %s, Issuer: %s, Rating: %f, EnergyUsage: %f", ratingID, nodeID, owner, issuer, rating, energyUsage)
	encryptedData, err := eers.EncryptionService.EncryptData([]byte(ratingData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt rating data: %v", err)
	}

	// Create the energy efficiency rating
	ratingRecord := &EfficiencyRating{
		RatingID:    ratingID,
		NodeID:      nodeID,
		Owner:       owner,
		Issuer:      issuer,
		Rating:      rating,
		IssueDate:   time.Now(),
		ExpiryDate:  time.Now().Add(validityPeriod),
		IsRevoked:   false,
		EnergyUsage: energyUsage,
	}

	// Add the rating to the system
	eers.Ratings[ratingID] = ratingRecord

	// Log the rating issuance in the ledger
	err = eers.Ledger.RecordEnergyEfficiencyRating(ratingID, nodeID, owner, issuer, rating, energyUsage, time.Now(), ratingRecord.ExpiryDate)
	if err != nil {
		return nil, fmt.Errorf("failed to log rating issuance: %v", err)
	}

	fmt.Printf("Energy efficiency rating %s issued to node %s (owner: %s) with a rating of %f\n", ratingID, nodeID, owner, rating)
	return ratingRecord, nil
}

// RevokeRating revokes an energy efficiency rating before its expiry
func (eers *EnergyEfficiencyRatingSystem) RevokeRating(ratingID string, revocationReason string) error {
	eers.mu.Lock()
	defer eers.mu.Unlock()

	// Retrieve the rating
	rating, exists := eers.Ratings[ratingID]
	if !exists {
		return fmt.Errorf("rating %s not found", ratingID)
	}

	// Ensure the rating has not already been revoked
	if rating.IsRevoked {
		return fmt.Errorf("rating %s has already been revoked", ratingID)
	}

	// Revoke the rating
	rating.IsRevoked = true
	rating.RevokedDate = time.Now()

	// Log the revocation in the ledger
	err := eers.Ledger.RecordRatingRevocation(ratingID, rating.NodeID, revocationReason, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rating revocation: %v", err)
	}

	fmt.Printf("Energy efficiency rating %s revoked for node %s\n", ratingID, rating.NodeID)
	return nil
}

// RenewRating renews an energy efficiency rating by extending its validity period
func (eers *EnergyEfficiencyRatingSystem) RenewRating(ratingID string, additionalPeriod time.Duration) (*EfficiencyRating, error) {
	eers.mu.Lock()
	defer eers.mu.Unlock()

	// Retrieve the rating
	rating, exists := eers.Ratings[ratingID]
	if !exists {
		return nil, fmt.Errorf("rating %s not found", ratingID)
	}

	// Ensure the rating has not been revoked
	if rating.IsRevoked {
		return nil, fmt.Errorf("rating %s has been revoked and cannot be renewed", ratingID)
	}

	// Extend the expiry date
	oldExpiry := rating.ExpiryDate
	rating.ExpiryDate = rating.ExpiryDate.Add(additionalPeriod)

	// Log the renewal in the ledger
	err := eers.Ledger.RecordRatingRenewal(ratingID, rating.NodeID, oldExpiry, rating.ExpiryDate, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log rating renewal: %v", err)
	}

	fmt.Printf("Energy efficiency rating %s renewed for node %s (new expiry: %s)\n", ratingID, rating.NodeID, rating.ExpiryDate)
	return rating, nil
}

// ViewRating allows viewing of the details of a specific energy efficiency rating
func (eers *EnergyEfficiencyRatingSystem) ViewRating(ratingID string) (*EfficiencyRating, error) {
	eers.mu.Lock()
	defer eers.mu.Unlock()

	// Retrieve the rating
	rating, exists := eers.Ratings[ratingID]
	if !exists {
		return nil, fmt.Errorf("rating %s not found", ratingID)
	}

	return rating, nil
}
