package syn12

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// SYN12ComplianceManager handles the compliance checks for SYN12 tokens, including KYC/AML.
type SYN12ComplianceManager struct {
	ledgerManager     *ledger.LedgerManager         // Ledger manager for compliance recording
	encryptionService *encryption.EncryptionService // Encryption service for secure data handling
	consensus         *consensus.SynnergyConsensus  // Consensus engine for validation
	approvedEntities  map[string]bool               // Approved entities for compliance
	mutex             sync.Mutex                    // Mutex for concurrency
}

// NewSYN12ComplianceManager initializes the compliance manager with ledger, consensus, and encryption services.
func NewSYN12ComplianceManager(ledgerManager *ledger.LedgerManager, encryptionService *encryption.EncryptionService, consensus *consensus.SynnergyConsensus) *SYN12ComplianceManager {
	return &SYN12ComplianceManager{
		ledgerManager:     ledgerManager,
		encryptionService: encryptionService,
		consensus:         consensus,
		approvedEntities:  make(map[string]bool),
	}
}

// VerifyIssuer checks if the issuer of a SYN12 token is compliant with KYC/AML requirements.
func (cm *SYN12ComplianceManager) VerifyIssuer(issuerID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Ensure the issuer is in the approved list
	if !cm.approvedEntities[issuerID] {
		return fmt.Errorf("issuer %s is not KYC/AML compliant", issuerID)
	}

	// Validate using consensus
	if err := cm.consensus.ValidateIssuer(issuerID); err != nil {
		return fmt.Errorf("issuer validation failed: %v", err)
	}

	return nil
}

// ApproveEntity adds an entity to the approved list after passing compliance.
func (cm *SYN12ComplianceManager) ApproveEntity(entityID string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.approvedEntities[entityID] = true
}

// IsEntityApproved checks if an entity is approved.
func (cm *SYN12ComplianceManager) IsEntityApproved(entityID string) bool {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	return cm.approvedEntities[entityID]
}

// ValidateTransaction ensures that a SYN12 token transaction complies with regulatory standards.
func (cm *SYN12ComplianceManager) ValidateTransaction(tokenID, senderID, receiverID string, amount uint64) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Verify sender and receiver compliance
	if !cm.IsEntityApproved(senderID) {
		return errors.New("sender is not KYC/AML compliant")
	}

	if !cm.IsEntityApproved(receiverID) {
		return errors.New("receiver is not KYC/AML compliant")
	}

	// Validate transaction through consensus
	if err := cm.consensus.ValidateTransaction(tokenID, senderID, receiverID, amount); err != nil {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	return nil
}

// LogCompliance stores a compliance check event in the ledger.
func (cm *SYN12ComplianceManager) LogCompliance(tokenID, entityID, complianceStatus string) error {
	eventID := common.GenerateUUID()

	complianceEvent := fmt.Sprintf("Compliance Check: Token %s, Entity %s, Status: %s", tokenID, entityID, complianceStatus)

	// Encrypt the event for secure storage
	encryptedEvent, err := cm.encryptionService.Encrypt([]byte(complianceEvent))
	if err != nil {
		return fmt.Errorf("failed to encrypt compliance event: %v", err)
	}

	// Record the event in the ledger
	if err := cm.ledgerManager.RecordEvent(eventID, tokenID, string(encryptedEvent)); err != nil {
		return fmt.Errorf("failed to log compliance event: %v", err)
	}

	return nil
}

// VerifyRedeemer ensures that the redeemer of a SYN12 token is compliant before processing redemption.
func (cm *SYN12ComplianceManager) VerifyRedeemer(tokenID, redeemerID string) error {
	// Check if the redeemer is approved for KYC/AML compliance
	if !cm.IsEntityApproved(redeemerID) {
		return fmt.Errorf("redeemer %s is not KYC/AML compliant", redeemerID)
	}

	// Use consensus to verify the redeemer's legitimacy
	if err := cm.consensus.ValidateRedeemer(tokenID, redeemerID); err != nil {
		return fmt.Errorf("redeemer validation failed: %v", err)
	}

	return nil
}

// ValidateMetadataUpdate checks compliance for updates to SYN12 token metadata.
func (cm *SYN12ComplianceManager) ValidateMetadataUpdate(oldMetadata, newMetadata common.TokenMetadata, updaterID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Ensure updater is compliant
	if !cm.IsEntityApproved(updaterID) {
		return errors.New("updater is not KYC/AML compliant")
	}

	// Log the compliance check
	if err := cm.LogCompliance(newMetadata.TokenID, updaterID, "Metadata Updated"); err != nil {
		return fmt.Errorf("failed to log metadata update compliance: %v", err)
	}

	// Consensus validation for metadata update
	if err := cm.consensus.ValidateMetadataUpdate(oldMetadata, newMetadata); err != nil {
		return fmt.Errorf("consensus validation failed for metadata update: %v", err)
	}

	return nil
}
