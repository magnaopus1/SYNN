package syn130

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"

)

// Syn130EventManager manages the lifecycle and event handling for SYN130 tokens.
type Syn130EventManager struct {
	EventLedger      *ledger.EventLedger          // Ledger to track events
	Consensus        *consensus.SynnergyConsensus // Synnergy consensus engine for event validation
	EncryptionService *encryption.EncryptionService // Encryption service for secure event handling
	mutex            sync.Mutex                    // Mutex to handle concurrent event processing
}

// NewSyn130EventManager initializes a new event manager.
func NewSyn130EventManager(eventLedger *ledger.EventLedger, encryptionService *encryption.EncryptionService, consensusEngine *consensus.SynnergyConsensus) *Syn130EventManager {
	return &Syn130EventManager{
		EventLedger:      eventLedger,
		EncryptionService: encryptionService,
		Consensus:        consensusEngine,
	}
}

// AssetCreatedEvent logs the creation of a new asset.
func (em *Syn130EventManager) AssetCreatedEvent(assetID, creator string, metadata map[string]string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Validate the event using Synnergy Consensus
	if err := em.Consensus.ValidateAssetCreation(assetID, creator); err != nil {
		return fmt.Errorf("asset creation validation failed: %v", err)
	}

	// Create event data
	event := ledger.EventRecord{
		EventType: "AssetCreated",
		EntityID:  assetID,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Asset %s created by %s", assetID, creator),
	}

	// Encrypt event details before recording
	encryptedDetails, err := em.EncryptionService.EncryptData(event.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt event details: %v", err)
	}
	event.Details = encryptedDetails

	// Record event in the event ledger
	if err := em.EventLedger.RecordEvent(&event); err != nil {
		return fmt.Errorf("failed to record asset creation event: %v", err)
	}

	fmt.Printf("Asset creation event for asset %s created successfully\n", assetID)
	return nil
}

// AssetTransferredEvent logs the transfer of an asset.
func (em *Syn130EventManager) AssetTransferredEvent(assetID, from, to string, amount float64) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Validate the transfer using Synnergy Consensus
	if err := em.Consensus.ValidateAssetTransfer(assetID, from, to, amount); err != nil {
		return fmt.Errorf("asset transfer validation failed: %v", err)
	}

	// Create event data
	event := ledger.EventRecord{
		EventType: "AssetTransferred",
		EntityID:  assetID,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Asset %s transferred from %s to %s for amount %.2f", assetID, from, to, amount),
	}

	// Encrypt event details before recording
	encryptedDetails, err := em.EncryptionService.EncryptData(event.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt event details: %v", err)
	}
	event.Details = encryptedDetails

	// Record event in the event ledger
	if err := em.EventLedger.RecordEvent(&event); err != nil {
		return fmt.Errorf("failed to record asset transfer event: %v", err)
	}

	fmt.Printf("Asset transfer event for asset %s from %s to %s created successfully\n", assetID, from, to)
	return nil
}

// LeaseAgreementEvent logs a new lease agreement event for an asset.
func (em *Syn130EventManager) LeaseAgreementEvent(assetID, lessor, lessee string, startDate, endDate time.Time, payment float64) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Validate the lease agreement using Synnergy Consensus
	if err := em.Consensus.ValidateLeaseAgreement(assetID, lessor, lessee); err != nil {
		return fmt.Errorf("lease agreement validation failed: %v", err)
	}

	// Create event data
	event := ledger.EventRecord{
		EventType: "LeaseAgreement",
		EntityID:  assetID,
		Timestamp: time.Now(),
		Details: fmt.Sprintf("Lease agreement for asset %s between %s (lessor) and %s (lessee) from %s to %s, payment: %.2f",
			assetID, lessor, lessee, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), payment),
	}

	// Encrypt event details before recording
	encryptedDetails, err := em.EncryptionService.EncryptData(event.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt event details: %v", err)
	}
	event.Details = encryptedDetails

	// Record event in the event ledger
	if err := em.EventLedger.RecordEvent(&event); err != nil {
		return fmt.Errorf("failed to record lease agreement event: %v", err)
	}

	fmt.Printf("Lease agreement event for asset %s created successfully\n", assetID)
	return nil
}

// LicenseAgreementEvent logs a new license agreement event for an asset.
func (em *Syn130EventManager) LicenseAgreementEvent(assetID, licensor, licensee string, startDate, endDate time.Time, licenseFee float64) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Validate the license agreement using Synnergy Consensus
	if err := em.Consensus.ValidateLicenseAgreement(assetID, licensor, licensee); err != nil {
		return fmt.Errorf("license agreement validation failed: %v", err)
	}

	// Create event data
	event := ledger.EventRecord{
		EventType: "LicenseAgreement",
		EntityID:  assetID,
		Timestamp: time.Now(),
		Details: fmt.Sprintf("License agreement for asset %s between %s (licensor) and %s (licensee) from %s to %s, license fee: %.2f",
			assetID, licensor, licensee, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), licenseFee),
	}

	// Encrypt event details before recording
	encryptedDetails, err := em.EncryptionService.EncryptData(event.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt event details: %v", err)
	}
	event.Details = encryptedDetails

	// Record event in the event ledger
	if err := em.EventLedger.RecordEvent(&event); err != nil {
		return fmt.Errorf("failed to record license agreement event: %v", err)
	}

	fmt.Printf("License agreement event for asset %s created successfully\n", assetID)
	return nil
}

// RentalAgreementEvent logs a new rental agreement event for an asset.
func (em *Syn130EventManager) RentalAgreementEvent(assetID, lessor, lessee string, startDate, endDate time.Time, rentalRate float64) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Validate the rental agreement using Synnergy Consensus
	if err := em.Consensus.ValidateRentalAgreement(assetID, lessor, lessee); err != nil {
		return fmt.Errorf("rental agreement validation failed: %v", err)
	}

	// Create event data
	event := ledger.EventRecord{
		EventType: "RentalAgreement",
		EntityID:  assetID,
		Timestamp: time.Now(),
		Details: fmt.Sprintf("Rental agreement for asset %s between %s (lessor) and %s (lessee) from %s to %s, rental rate: %.2f",
			assetID, lessor, lessee, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), rentalRate),
	}

	// Encrypt event details before recording
	encryptedDetails, err := em.EncryptionService.EncryptData(event.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt event details: %v", err)
	}
	event.Details = encryptedDetails

	// Record event in the event ledger
	if err := em.EventLedger.RecordEvent(&event); err != nil {
		return fmt.Errorf("failed to record rental agreement event: %v", err)
	}

	fmt.Printf("Rental agreement event for asset %s created successfully\n", assetID)
	return nil
}

// AssetDestroyedEvent logs the destruction or retirement of an asset.
func (em *Syn130EventManager) AssetDestroyedEvent(assetID, destroyer string, reason string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Validate the destruction event using Synnergy Consensus
	if err := em.Consensus.ValidateAssetDestruction(assetID, destroyer); err != nil {
		return fmt.Errorf("asset destruction validation failed: %v", err)
	}

	// Create event data
	event := ledger.EventRecord{
		EventType: "AssetDestroyed",
		EntityID:  assetID,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Asset %s destroyed by %s. Reason: %s", assetID, destroyer, reason),
	}

	// Encrypt event details before recording
	encryptedDetails, err := em.EncryptionService.EncryptData(event.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt event details: %v", err)
	}
	event.Details = encryptedDetails

	// Record event in the event ledger
	if err := em.EventLedger.RecordEvent(&event); err != nil {
		return fmt.Errorf("failed to record asset destruction event: %v", err)
	}

	fmt.Printf("Asset destruction event for asset %s created successfully\n", assetID)
	return nil
}

// DecryptEvent decrypts an event's details for review or audit purposes.
func (em *Syn130EventManager) DecryptEvent(event *ledger.EventRecord) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Decrypt the event details
	decryptedDetails, err := em.EncryptionService.DecryptData(event.Details)
	if err != nil {
		return fmt.Errorf("failed to decrypt event details: %v", err)
	}

	// Update the event with decrypted details
	event.Details = decryptedDetails
	return nil
}
