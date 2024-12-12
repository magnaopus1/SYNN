package syn845

import (
	"errors"
	"sync"
	"time"

)

// SYN845ComplianceManager handles compliance for SYN845 tokens, including KYC, AML, and jurisdiction checks
type SYN845ComplianceManager struct {
	Ledger              *ledger.Ledger               // Ledger for recording compliance actions
	ConsensusEngine     *consensus.SynnergyConsensus  // Synnergy Consensus for validating compliance
	EncryptionService   *encryption.EncryptionService // Encryption service for securing compliance data
	ApprovedJurisdictions []string                    // List of approved jurisdictions for transactions
	mutex               sync.Mutex                    // Mutex for safe concurrent access
}

// KYCInfo holds the Know Your Customer information for a user
type KYCInfo struct {
	Verified         bool      `json:"verified"`
	VerificationDate time.Time `json:"verification_date"`
}

// AMLInfo holds the Anti-Money Laundering information for a user
type AMLInfo struct {
	Verified         bool      `json:"verified"`
	VerificationDate time.Time `json:"verification_date"`
	RiskLevel        string    `json:"risk_level"`
}

// TransactionLimit defines limits on transaction amounts for compliance purposes
type TransactionLimit struct {
	MaxAmount float64       `json:"max_amount"`
	TimeFrame time.Duration `json:"time_frame"`
}

// NewSYN845ComplianceManager initializes a new instance of SYN845ComplianceManager
func NewSYN845ComplianceManager(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService, approvedJurisdictions []string) *SYN845ComplianceManager {
	return &SYN845ComplianceManager{
		Ledger:               ledger,
		ConsensusEngine:      consensusEngine,
		EncryptionService:    encryptionService,
		ApprovedJurisdictions: approvedJurisdictions,
	}
}

// AddKYCInfo adds or updates KYC information for a user, securely storing the data in the ledger
func (cm *SYN845ComplianceManager) AddKYCInfo(userID string, verified bool, verificationDate time.Time) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	kyc := KYCInfo{Verified: verified, VerificationDate: verificationDate}

	// Encrypt KYC information
	encryptedData, encryptionKey, err := cm.EncryptionService.EncryptData([]byte(common.StructToString(kyc)))
	if err != nil {
		return errors.New("failed to encrypt KYC information")
	}

	// Validate and store the KYC information in the ledger
	if err := cm.ConsensusEngine.ValidateKYC(userID, kyc); err != nil {
		return errors.New("KYC validation failed via Synnergy Consensus")
	}

	if err := cm.Ledger.StoreKYCData(userID, encryptedData, encryptionKey); err != nil {
		return errors.New("failed to store KYC information in the ledger")
	}

	return nil
}

// AddAMLInfo adds or updates AML information for a user, securely storing the data in the ledger
func (cm *SYN845ComplianceManager) AddAMLInfo(userID string, verified bool, verificationDate time.Time, riskLevel string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	aml := AMLInfo{Verified: verified, VerificationDate: verificationDate, RiskLevel: riskLevel}

	// Encrypt AML information
	encryptedData, encryptionKey, err := cm.EncryptionService.EncryptData([]byte(common.StructToString(aml)))
	if err != nil {
		return errors.New("failed to encrypt AML information")
	}

	// Validate and store the AML information in the ledger
	if err := cm.ConsensusEngine.ValidateAML(userID, aml); err != nil {
		return errors.New("AML validation failed via Synnergy Consensus")
	}

	if err := cm.Ledger.StoreAMLData(userID, encryptedData, encryptionKey); err != nil {
		return errors.New("failed to store AML information in the ledger")
	}

	return nil
}

// SetTransactionLimit sets a transaction limit for a user, storing the limit in the ledger
func (cm *SYN845ComplianceManager) SetTransactionLimit(userID string, maxAmount float64, timeFrame time.Duration) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	limit := TransactionLimit{MaxAmount: maxAmount, TimeFrame: timeFrame}

	// Encrypt transaction limit data
	encryptedData, encryptionKey, err := cm.EncryptionService.EncryptData([]byte(common.StructToString(limit)))
	if err != nil {
		return errors.New("failed to encrypt transaction limit information")
	}

	// Validate and store the transaction limit in the ledger
	if err := cm.ConsensusEngine.ValidateTransactionLimit(userID, limit); err != nil {
		return errors.New("transaction limit validation failed via Synnergy Consensus")
	}

	if err := cm.Ledger.StoreTransactionLimit(userID, encryptedData, encryptionKey); err != nil {
		return errors.New("failed to store transaction limit in the ledger")
	}

	return nil
}

// VerifyKYC checks if a user meets KYC requirements by retrieving the encrypted KYC data and verifying it
func (cm *SYN845ComplianceManager) VerifyKYC(userID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Retrieve KYC information from the ledger
	encryptedData, encryptionKey, err := cm.Ledger.GetKYCData(userID)
	if err != nil {
		return errors.New("failed to retrieve KYC information from ledger")
	}

	// Decrypt the KYC information
	decryptedData, err := cm.EncryptionService.DecryptData([]byte(encryptedData), encryptionKey)
	if err != nil {
		return errors.New("failed to decrypt KYC information")
	}

	// Check if KYC is verified
	kyc := common.StringToStruct(decryptedData, KYCInfo{}).(KYCInfo)
	if !kyc.Verified {
		return errors.New("KYC verification required")
	}

	return nil
}

// VerifyAML checks if a user meets AML requirements by retrieving the encrypted AML data and verifying it
func (cm *SYN845ComplianceManager) VerifyAML(userID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Retrieve AML information from the ledger
	encryptedData, encryptionKey, err := cm.Ledger.GetAMLData(userID)
	if err != nil {
		return errors.New("failed to retrieve AML information from ledger")
	}

	// Decrypt the AML information
	decryptedData, err := cm.EncryptionService.DecryptData([]byte(encryptedData), encryptionKey)
	if err != nil {
		return errors.New("failed to decrypt AML information")
	}

	// Check if AML is verified
	aml := common.StringToStruct(decryptedData, AMLInfo{}).(AMLInfo)
	if !aml.Verified {
		return errors.New("AML verification required")
	}

	return nil
}

// CheckTransactionLimit verifies if a transaction is within the set limit by retrieving the encrypted limit data
func (cm *SYN845ComplianceManager) CheckTransactionLimit(userID string, amount float64) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Retrieve transaction limit information from the ledger
	encryptedData, encryptionKey, err := cm.Ledger.GetTransactionLimit(userID)
	if err != nil {
		return errors.New("failed to retrieve transaction limit information from ledger")
	}

	// Decrypt the transaction limit
	decryptedData, err := cm.EncryptionService.DecryptData([]byte(encryptedData), encryptionKey)
	if err != nil {
		return errors.New("failed to decrypt transaction limit information")
	}

	limit := common.StringToStruct(decryptedData, TransactionLimit{}).(TransactionLimit)
	if amount > limit.MaxAmount {
		return errors.New("transaction amount exceeds limit")
	}

	return nil
}

// IsJurisdictionApproved checks if a jurisdiction is approved for transactions
func (cm *SYN845ComplianceManager) IsJurisdictionApproved(jurisdiction string) bool {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for _, approved := range cm.ApprovedJurisdictions {
		if approved == jurisdiction {
			return true
		}
	}
	return false
}

// EnsureCompliance checks all compliance requirements before allowing a transaction to proceed
func (cm *SYN845ComplianceManager) EnsureCompliance(userID string, amount float64, jurisdiction string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if err := cm.VerifyKYC(userID); err != nil {
		return err
	}

	if err := cm.VerifyAML(userID); err != nil {
		return err
	}

	if !cm.IsJurisdictionApproved(jurisdiction) {
		return errors.New("transaction jurisdiction not approved")
	}

	if err := cm.CheckTransactionLimit(userID, amount); err != nil {
		return err
	}

	// Ensure compliance is validated through Synnergy Consensus
	if err := cm.ConsensusEngine.ValidateCompliance(userID, amount, jurisdiction); err != nil {
		return errors.New("compliance validation failed via Synnergy Consensus")
	}

	return nil
}
