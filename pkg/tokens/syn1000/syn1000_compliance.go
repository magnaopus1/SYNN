package syn1000

import (
	"errors"
	"sync"
	"time"

)

// ComplianceStatus defines the compliance information for SYN1000 tokens
type ComplianceStatus struct {
	KYCVerified          bool      `json:"kyc_verified"`
	AMLVerified          bool      `json:"aml_verified"`
	ApprovedJurisdiction string    `json:"approved_jurisdiction"`
	ComplianceDate       time.Time `json:"compliance_date"`
}

// SYN1000ComplianceManager manages the compliance and regulatory functions for SYN1000 tokens
type SYN1000ComplianceManager struct {
	mutex              sync.Mutex
	Ledger             *ledger.Ledger                // Ledger for storing compliance data
	ConsensusEngine    *consensus.SynnergyConsensus  // Consensus engine for validating compliance actions
	EncryptionService  *encryption.EncryptionService // Encryption service for securing sensitive compliance data
	ComplianceRecords  map[string]ComplianceStatus   // Compliance records for each token
}

// NewSYN1000ComplianceManager creates a new SYN1000ComplianceManager
func NewSYN1000ComplianceManager(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *SYN1000ComplianceManager {
	return &SYN1000ComplianceManager{
		Ledger:            ledger,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		ComplianceRecords: make(map[string]ComplianceStatus),
	}
}

// AddComplianceRecord adds a compliance record for a SYN1000 token
func (cm *SYN1000ComplianceManager) AddComplianceRecord(tokenID, owner string, jurisdiction string, kycVerified, amlVerified bool) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Validate if the jurisdiction is approved
	if !cm.isJurisdictionApproved(jurisdiction) {
		return errors.New("unapproved jurisdiction for token operations")
	}

	// Create the compliance status
	compliance := ComplianceStatus{
		KYCVerified:          kycVerified,
		AMLVerified:          amlVerified,
		ApprovedJurisdiction: jurisdiction,
		ComplianceDate:       time.Now(),
	}

	// Encrypt compliance data
	complianceData, encryptionKey, err := cm.EncryptionService.EncryptData([]byte(common.StructToString(compliance)))
	if err != nil {
		return errors.New("failed to encrypt compliance data")
	}

	// Validate the compliance record using Synnergy Consensus
	if err := cm.ConsensusEngine.ValidateCompliance(tokenID, compliance); err != nil {
		return errors.New("compliance validation failed via Synnergy Consensus")
	}

	// Store the compliance record in the ledger
	if err := cm.Ledger.StoreComplianceData(tokenID, complianceData, encryptionKey); err != nil {
		return errors.New("failed to store compliance data in the ledger")
	}

	// Add compliance to the in-memory store
	cm.ComplianceRecords[tokenID] = compliance

	return nil
}

// RetrieveComplianceRecord retrieves the compliance record for a SYN1000 token
func (cm *SYN1000ComplianceManager) RetrieveComplianceRecord(tokenID string) (*ComplianceStatus, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Retrieve encrypted compliance data from the ledger
	encryptedData, encryptionKey, err := cm.Ledger.GetComplianceData(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve compliance data from ledger")
	}

	// Decrypt the compliance data
	decryptedData, err := cm.EncryptionService.DecryptData([]byte(encryptedData), encryptionKey)
	if err != nil {
		return nil, errors.New("failed to decrypt compliance data")
	}

	// Parse the decrypted data into a ComplianceStatus struct
	compliance := common.StringToStruct(string(decryptedData), ComplianceStatus{}).(ComplianceStatus)

	return &compliance, nil
}

// VerifyCompliance checks if a SYN1000 token complies with KYC/AML and jurisdiction rules
func (cm *SYN1000ComplianceManager) VerifyCompliance(tokenID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	compliance, exists := cm.ComplianceRecords[tokenID]
	if !exists {
		return errors.New("compliance record not found")
	}

	// Check if KYC and AML verification passed
	if !compliance.KYCVerified || !compliance.AMLVerified {
		return errors.New("token does not meet KYC/AML requirements")
	}

	return nil
}

// isJurisdictionApproved checks if the provided jurisdiction is approved for SYN1000 token operations
func (cm *SYN1000ComplianceManager) isJurisdictionApproved(jurisdiction string) bool {
	// In a real-world implementation, this could check against a dynamic list of approved jurisdictions.
	approvedJurisdictions := []string{"Global", "USA", "EU", "UK", "Singapore"}
	for _, approved := range approvedJurisdictions {
		if approved == jurisdiction {
			return true
		}
	}
	return false
}

// UpdateComplianceRecord updates an existing compliance record for a SYN1000 token
func (cm *SYN1000ComplianceManager) UpdateComplianceRecord(tokenID string, kycVerified, amlVerified bool, jurisdiction string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Retrieve current compliance record
	compliance, err := cm.RetrieveComplianceRecord(tokenID)
	if err != nil {
		return err
	}

	// Update compliance details
	compliance.KYCVerified = kycVerified
	compliance.AMLVerified = amlVerified
	compliance.ApprovedJurisdiction = jurisdiction
	compliance.ComplianceDate = time.Now()

	// Encrypt updated compliance data
	encryptedData, encryptionKey, err := cm.EncryptionService.EncryptData([]byte(common.StructToString(*compliance)))
	if err != nil {
		return errors.New("failed to encrypt updated compliance data")
	}

	// Validate compliance update via Synnergy Consensus
	if err := cm.ConsensusEngine.ValidateCompliance(tokenID, *compliance); err != nil {
		return errors.New("compliance validation failed via Synnergy Consensus")
	}

	// Update the ledger with the new compliance record
	if err := cm.Ledger.UpdateComplianceData(tokenID, encryptedData, encryptionKey); err != nil {
		return errors.New("failed to update compliance data in the ledger")
	}

	return nil
}

// RemoveComplianceRecord deletes a compliance record from both in-memory and ledger storage
func (cm *SYN1000ComplianceManager) RemoveComplianceRecord(tokenID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Remove the compliance record from the ledger
	if err := cm.Ledger.DeleteComplianceData(tokenID); err != nil {
		return errors.New("failed to delete compliance data from ledger")
	}

	// Remove the record from in-memory storage
	delete(cm.ComplianceRecords, tokenID)

	return nil
}
