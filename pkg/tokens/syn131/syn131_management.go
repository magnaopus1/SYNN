package syn131

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)



// NewAgreementManager creates a new AgreementManager.
func NewAgreementManager() *AgreementManager {
	return &AgreementManager{
		RentalAgreements:      make(map[string]*RentalAgreement),
		LeaseAgreements:       make(map[string]*LeaseAgreement),
		CoOwnershipAgreements: make(map[string]*CoOwnershipAgreement),
	}
}

// CreateRentalAgreement creates a new rental agreement.
func (am *AgreementManager) CreateRentalAgreement(assetID, lessor, lessee, terms string, startDate, endDate time.Time) (*RentalAgreement, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Encrypt the rental agreement terms
	encryptedTerms, encryptionKey, err := encryption.EncryptData([]byte(terms))
	if err != nil {
		return nil, errors.New("failed to encrypt rental agreement terms")
	}

	agreement := &RentalAgreement{
		ID:             generateAgreementID(assetID, lessee),
		AssetID:        assetID,
		Lessor:         lessor,
		Lessee:         lessee,
		StartDate:      startDate,
		EndDate:        endDate,
		Terms:          terms,
		EncryptedTerms: encryptedTerms,
		EncryptionKey:  encryptionKey,
		Status:         "active",
		PaymentSchedule: "monthly",
	}

	// Store the agreement
	am.RentalAgreements[agreement.ID] = agreement
	return agreement, nil
}

// CreateLeaseAgreement creates a new lease agreement.
func (am *AgreementManager) CreateLeaseAgreement(assetID, lessor, lessee, terms string, startDate, endDate time.Time, paymentSchedule map[time.Time]float64) (*LeaseAgreement, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Encrypt the lease agreement terms
	encryptedTerms, encryptionKey, err := encryption.EncryptData([]byte(terms))
	if err != nil {
		return nil, errors.New("failed to encrypt lease agreement terms")
	}

	agreement := &LeaseAgreement{
		ID:              generateAgreementID(assetID, lessee),
		AssetID:         assetID,
		Lessor:          lessor,
		Lessee:          lessee,
		StartDate:       startDate,
		EndDate:         endDate,
		Terms:           terms,
		EncryptedTerms:  encryptedTerms,
		EncryptionKey:   encryptionKey,
		Status:          "active",
		PaymentSchedule: paymentSchedule,
	}

	// Store the agreement
	am.LeaseAgreements[agreement.ID] = agreement
	return agreement, nil
}

// CreateCoOwnershipAgreement creates a new co-ownership agreement.
func (am *AgreementManager) CreateCoOwnershipAgreement(assetID string, owners map[string]float64, terms string) (*CoOwnershipAgreement, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	agreement := &CoOwnershipAgreement{
		AgreementID:     generateAgreementID(assetID, ""),
		AssetID:         assetID,
		Owners:          owners,
		CreationDate:    time.Now(),
		ModificationDate: time.Now(),
		Terms:           terms,
		Status:          "active",
	}

	// Store the agreement
	am.CoOwnershipAgreements[agreement.AgreementID] = agreement
	return agreement, nil
}

// ValuationManager handles asset valuation updates.
func (vm *ValuationManager) UpdateValuation(assetID string, newValue float64) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	valuation, exists := vm.Valuations[assetID]
	if !exists {
		return errors.New("asset valuation not found")
	}

	// Record historical valuation
	valuation.HistoricalRecords = append(valuation.HistoricalRecords, HistoricalValuation{
		Value:     newValue,
		Timestamp: time.Now(),
	})

	// Update current value and last updated time
	valuation.CurrentValue = newValue
	valuation.LastUpdated = time.Now()

	return nil
}

// GenerateAgreementID generates a unique agreement ID based on asset and party.
func generateAgreementID(assetID, party string) string {
	return assetID + "_" + party + "_" + time.Now().Format("20060102150405")
}

// TransferAssetOwnership handles the transfer of ownership of an intangible asset.
func (am *AssetManager) TransferAssetOwnership(assetID, newOwner string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	asset, exists := am.Assets[assetID]
	if !exists {
		return errors.New("asset not found")
	}

	// Add transfer record to history
	asset.TransferHistory = append(asset.TransferHistory, TransferRecord{
		From:      asset.Owner,
		To:        newOwner,
		Timestamp: time.Now(),
	})

	// Update asset owner
	asset.Owner = newOwner
	return nil
}

// UpdateAssetMetadata allows updating the metadata of an intangible asset.
func (am *AssetManager) UpdateAssetMetadata(assetID string, newMetadata assets.AssetMetadata) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	asset, exists := am.Assets[assetID]
	if !exists {
		return errors.New("asset not found")
	}

	// Update the asset's metadata
	asset.Metadata = newMetadata
	return nil
}
