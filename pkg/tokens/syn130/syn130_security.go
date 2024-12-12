package syn130

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SecurityManager is responsible for enforcing all security measures for SYN130 tokens.
type SecurityManager struct {
	SecurityLedger    *ledger.SecurityLedger         // Ledger to store security-related events
	Consensus         *consensus.SynnergyConsensus   // Consensus engine to ensure security validation
	EncryptionService *encryption.EncryptionService  // Encryption service to secure sensitive information
	mutex             sync.Mutex                     // Mutex for thread-safe security operations
}

// NewSecurityManager initializes a new SecurityManager.
func NewSecurityManager(ledger *ledger.SecurityLedger, encryptionService *encryption.EncryptionService, consensusEngine *consensus.SynnergyConsensus) *SecurityManager {
	return &SecurityManager{
		SecurityLedger:    ledger,
		EncryptionService: encryptionService,
		Consensus:         consensusEngine,
	}
}

// ValidateTransactionSecurity ensures the integrity and security of a transaction.
func (sm *SecurityManager) ValidateTransactionSecurity(transactionID, from, to string, amount float64) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Validate the transaction through Synnergy Consensus
	if err := sm.Consensus.ValidateTransaction(transactionID, from, to, amount); err != nil {
		return fmt.Errorf("transaction security validation failed: %v", err)
	}

	// Encrypt and log the security event in the ledger
	event := ledger.SecurityEvent{
		EventID:    transactionID,
		From:       from,
		To:         to,
		Amount:     amount,
		Timestamp:  time.Now(),
		EventType:  "Transaction Security Validation",
	}

	encryptedDetails, err := sm.EncryptionService.EncryptData(event.String())
	if err != nil {
		return fmt.Errorf("failed to encrypt security event: %v", err)
	}
	event.Details = encryptedDetails

	if err := sm.SecurityLedger.RecordSecurityEvent(&event); err != nil {
		return fmt.Errorf("failed to record security event in ledger: %v", err)
	}

	return nil
}

// ValidateOwnershipTransfer ensures the secure transfer of ownership between parties.
func (sm *SecurityManager) ValidateOwnershipTransfer(assetID, currentOwner, newOwner string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Validate ownership transfer using Synnergy Consensus
	if err := sm.Consensus.ValidateOwnershipTransfer(assetID, currentOwner, newOwner); err != nil {
		return fmt.Errorf("ownership transfer security validation failed: %v", err)
	}

	// Encrypt and log the security event in the ledger
	event := ledger.SecurityEvent{
		EventID:    assetID,
		From:       currentOwner,
		To:         newOwner,
		Timestamp:  time.Now(),
		EventType:  "Ownership Transfer Security Validation",
	}

	encryptedDetails, err := sm.EncryptionService.EncryptData(event.String())
	if err != nil {
		return fmt.Errorf("failed to encrypt security event: %v", err)
	}
	event.Details = encryptedDetails

	if err := sm.SecurityLedger.RecordSecurityEvent(&event); err != nil {
		return fmt.Errorf("failed to record security event in ledger: %v", err)
	}

	return nil
}

// ValidateLeaseAgreementSecurity ensures that lease agreements are secure and valid.
func (sm *SecurityManager) ValidateLeaseAgreementSecurity(leaseID, assetID, lessor, lessee string, startDate, endDate time.Time) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Validate the lease agreement using Synnergy Consensus
	if err := sm.Consensus.ValidateLeaseAgreement(assetID, lessor, lessee); err != nil {
		return fmt.Errorf("lease agreement security validation failed: %v", err)
	}

	// Encrypt and log the security event in the ledger
	event := ledger.SecurityEvent{
		EventID:    leaseID,
		From:       lessor,
		To:         lessee,
		Timestamp:  time.Now(),
		EventType:  "Lease Agreement Security Validation",
	}

	encryptedDetails, err := sm.EncryptionService.EncryptData(event.String())
	if err != nil {
		return fmt.Errorf("failed to encrypt security event: %v", err)
	}
	event.Details = encryptedDetails

	if err := sm.SecurityLedger.RecordSecurityEvent(&event); err != nil {
		return fmt.Errorf("failed to record security event in ledger: %v", err)
	}

	return nil
}

// ValidateLicenseAgreementSecurity ensures that license agreements are secure and valid.
func (sm *SecurityManager) ValidateLicenseAgreementSecurity(licenseID, assetID, licensor, licensee string, startDate, endDate time.Time, licenseFee float64) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Validate the license agreement using Synnergy Consensus
	if err := sm.Consensus.ValidateLicenseAgreement(assetID, licensor, licensee); err != nil {
		return fmt.Errorf("license agreement security validation failed: %v", err)
	}

	// Encrypt and log the security event in the ledger
	event := ledger.SecurityEvent{
		EventID:    licenseID,
		From:       licensor,
		To:         licensee,
		Timestamp:  time.Now(),
		EventType:  "License Agreement Security Validation",
	}

	encryptedDetails, err := sm.EncryptionService.EncryptData(event.String())
	if err != nil {
		return fmt.Errorf("failed to encrypt security event: %v", err)
	}
	event.Details = encryptedDetails

	if err := sm.SecurityLedger.RecordSecurityEvent(&event); err != nil {
		return fmt.Errorf("failed to record security event in ledger: %v", err)
	}

	return nil
}

// MonitorAssetSecurity continuously monitors the security status of a given asset.
func (sm *SecurityManager) MonitorAssetSecurity(assetID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Monitor the asset through the consensus engine and trigger necessary alerts
	if err := sm.Consensus.MonitorAsset(assetID); err != nil {
		return fmt.Errorf("asset security monitoring failed: %v", err)
	}

	// Log the monitoring event
	event := ledger.SecurityEvent{
		EventID:   assetID,
		EventType: "Asset Security Monitoring",
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Security monitoring initiated for asset %s", assetID),
	}

	encryptedDetails, err := sm.EncryptionService.EncryptData(event.String())
	if err != nil {
		return fmt.Errorf("failed to encrypt security event: %v", err)
	}
	event.Details = encryptedDetails

	if err := sm.SecurityLedger.RecordSecurityEvent(&event); err != nil {
		return fmt.Errorf("failed to record security event in ledger: %v", err)
	}

	return nil
}

// ValidateAssetDestruction ensures the secure destruction of an asset.
func (sm *SecurityManager) ValidateAssetDestruction(assetID, destroyerID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Validate the destruction process through Synnergy Consensus
	if err := sm.Consensus.ValidateAssetDestruction(assetID, destroyerID); err != nil {
		return fmt.Errorf("asset destruction security validation failed: %v", err)
	}

	// Encrypt and log the security event in the ledger
	event := ledger.SecurityEvent{
		EventID:   assetID,
		From:      destroyerID,
		EventType: "Asset Destruction Security Validation",
		Timestamp: time.Now(),
	}

	encryptedDetails, err := sm.EncryptionService.EncryptData(event.String())
	if err != nil {
		return fmt.Errorf("failed to encrypt security event: %v", err)
	}
	event.Details = encryptedDetails

	if err := sm.SecurityLedger.RecordSecurityEvent(&event); err != nil {
		return fmt.Errorf("failed to record security event in ledger: %v", err)
	}

	return nil
}

// AuditSecurityEvent allows for the review of security events by decrypting the event details.
func (sm *SecurityManager) AuditSecurityEvent(event *ledger.SecurityEvent) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Decrypt the security event details
	decryptedDetails, err := sm.EncryptionService.DecryptData(event.Details)
	if err != nil {
		return fmt.Errorf("failed to decrypt security event details: %v", err)
	}

	event.Details = decryptedDetails
	return nil
}
