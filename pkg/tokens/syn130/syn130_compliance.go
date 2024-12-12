package syn130

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"

)

// ComplianceManager handles all compliance checks and processes related to SYN130 tokens.
type ComplianceManager struct {
	ComplianceLedger  *ledger.ComplianceLedger        // Ledger to track compliance records
	Consensus         *consensus.SynnergyConsensus    // Synnergy consensus engine for compliance validation
	EncryptionService *encryption.EncryptionService   // Encryption service for secure compliance handling
	mutex             sync.Mutex                      // Mutex for concurrent compliance handling
}

// NewComplianceManager initializes a new ComplianceManager.
func NewComplianceManager(ledger *ledger.ComplianceLedger, encryptionService *encryption.EncryptionService, consensusEngine *consensus.SynnergyConsensus) *ComplianceManager {
	return &ComplianceManager{
		ComplianceLedger:  ledger,
		EncryptionService: encryptionService,
		Consensus:         consensusEngine,
	}
}

// ValidateOwnership ensures the ownership of the asset complies with the SYN130 standard.
func (cm *ComplianceManager) ValidateOwnership(assetID, ownerID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Use Synnergy Consensus to validate ownership
	if err := cm.Consensus.ValidateOwnership(assetID, ownerID); err != nil {
		return fmt.Errorf("ownership validation failed: %v", err)
	}

	// Encrypt and record the compliance check in the ledger
	record := ledger.ComplianceRecord{
		EntityID:  assetID,
		OwnerID:   ownerID,
		Timestamp: time.Now(),
		Details:   "Ownership validated",
	}

	encryptedDetails, err := cm.EncryptionService.EncryptData(record.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt compliance details: %v", err)
	}
	record.Details = encryptedDetails

	if err := cm.ComplianceLedger.RecordCompliance(&record); err != nil {
		return fmt.Errorf("failed to record compliance in ledger: %v", err)
	}

	return nil
}

// ValidateAssetTransfer ensures the transfer of the asset complies with SYN130 token standards.
func (cm *ComplianceManager) ValidateAssetTransfer(assetID, from, to string, amount float64) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Use Synnergy Consensus to validate asset transfer
	if err := cm.Consensus.ValidateTransfer(assetID, from, to, amount); err != nil {
		return fmt.Errorf("asset transfer validation failed: %v", err)
	}

	// Encrypt and record the compliance check in the ledger
	record := ledger.ComplianceRecord{
		EntityID:  assetID,
		OwnerID:   from,
		NewOwnerID: to,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Asset transfer of %s from %s to %s for amount %.2f validated", assetID, from, to, amount),
	}

	encryptedDetails, err := cm.EncryptionService.EncryptData(record.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt compliance details: %v", err)
	}
	record.Details = encryptedDetails

	if err := cm.ComplianceLedger.RecordCompliance(&record); err != nil {
		return fmt.Errorf("failed to record compliance in ledger: %v", err)
	}

	return nil
}

// ValidateLeaseAgreement ensures lease agreements comply with SYN130 token standards.
func (cm *ComplianceManager) ValidateLeaseAgreement(assetID, lessor, lessee string, startDate, endDate time.Time) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Use Synnergy Consensus to validate the lease agreement
	if err := cm.Consensus.ValidateLeaseAgreement(assetID, lessor, lessee); err != nil {
		return fmt.Errorf("lease agreement validation failed: %v", err)
	}

	// Encrypt and record the compliance check in the ledger
	record := ledger.ComplianceRecord{
		EntityID:  assetID,
		OwnerID:   lessor,
		NewOwnerID: lessee,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Lease agreement for asset %s validated between %s (lessor) and %s (lessee) from %s to %s", assetID, lessor, lessee, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339)),
	}

	encryptedDetails, err := cm.EncryptionService.EncryptData(record.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt compliance details: %v", err)
	}
	record.Details = encryptedDetails

	if err := cm.ComplianceLedger.RecordCompliance(&record); err != nil {
		return fmt.Errorf("failed to record compliance in ledger: %v", err)
	}

	return nil
}

// ValidateLicenseAgreement ensures license agreements comply with SYN130 token standards.
func (cm *ComplianceManager) ValidateLicenseAgreement(assetID, licensor, licensee string, startDate, endDate time.Time, licenseFee float64) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Use Synnergy Consensus to validate the license agreement
	if err := cm.Consensus.ValidateLicenseAgreement(assetID, licensor, licensee); err != nil {
		return fmt.Errorf("license agreement validation failed: %v", err)
	}

	// Encrypt and record the compliance check in the ledger
	record := ledger.ComplianceRecord{
		EntityID:  assetID,
		OwnerID:   licensor,
		NewOwnerID: licensee,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("License agreement for asset %s validated between %s (licensor) and %s (licensee) from %s to %s, license fee: %.2f", assetID, licensor, licensee, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), licenseFee),
	}

	encryptedDetails, err := cm.EncryptionService.EncryptData(record.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt compliance details: %v", err)
	}
	record.Details = encryptedDetails

	if err := cm.ComplianceLedger.RecordCompliance(&record); err != nil {
		return fmt.Errorf("failed to record compliance in ledger: %v", err)
	}

	return nil
}

// ValidateRentalAgreement ensures rental agreements comply with SYN130 token standards.
func (cm *ComplianceManager) ValidateRentalAgreement(assetID, lessor, lessee string, startDate, endDate time.Time, rentalRate float64) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Use Synnergy Consensus to validate the rental agreement
	if err := cm.Consensus.ValidateRentalAgreement(assetID, lessor, lessee); err != nil {
		return fmt.Errorf("rental agreement validation failed: %v", err)
	}

	// Encrypt and record the compliance check in the ledger
	record := ledger.ComplianceRecord{
		EntityID:  assetID,
		OwnerID:   lessor,
		NewOwnerID: lessee,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Rental agreement for asset %s validated between %s (lessor) and %s (lessee) from %s to %s, rental rate: %.2f", assetID, lessor, lessee, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), rentalRate),
	}

	encryptedDetails, err := cm.EncryptionService.EncryptData(record.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt compliance details: %v", err)
	}
	record.Details = encryptedDetails

	if err := cm.ComplianceLedger.RecordCompliance(&record); err != nil {
		return fmt.Errorf("failed to record compliance in ledger: %v", err)
	}

	return nil
}

// ValidateAssetDestruction ensures the destruction of an asset complies with SYN130 token standards.
func (cm *ComplianceManager) ValidateAssetDestruction(assetID, destroyer string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Use Synnergy Consensus to validate the destruction
	if err := cm.Consensus.ValidateAssetDestruction(assetID, destroyer); err != nil {
		return fmt.Errorf("asset destruction validation failed: %v", err)
	}

	// Encrypt and record the compliance check in the ledger
	record := ledger.ComplianceRecord{
		EntityID:  assetID,
		OwnerID:   destroyer,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Asset %s destroyed by %s", assetID, destroyer),
	}

	encryptedDetails, err := cm.EncryptionService.EncryptData(record.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt compliance details: %v", err)
	}
	record.Details = encryptedDetails

	if err := cm.ComplianceLedger.RecordCompliance(&record); err != nil {
		return fmt.Errorf("failed to record compliance in ledger: %v", err)
	}

	return nil
}

// AuditCompliance allows the review of compliance checks by decrypting the event details.
func (cm *ComplianceManager) AuditCompliance(record *ledger.ComplianceRecord) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Decrypt the compliance details
	decryptedDetails, err := cm.EncryptionService.DecryptData(record.Details)
	if err != nil {
		return fmt.Errorf("failed to decrypt compliance details: %v", err)
	}

	record.Details = decryptedDetails
	return nil
}
