package syn2400

import (
	"errors"
	"time"

)

// SYN2400Token represents a data token in a decentralized data marketplace.
type SYN2400Token struct {
	TokenID             string              // Unique identifier for the data token
	Owner               string              // The owner of the data token
	DataHash            string              // Hash of the data (ensures integrity)
	Description         string              // Description of the data set
	AccessRights        AccessRights        // Access permissions associated with the data
	CreationDate        time.Time           // Timestamp for when the data token was created
	UpdateDate          time.Time           // Timestamp for the last update
	Price               float64             // Price of the data token in marketplace
	Status              string              // Status of the token (e.g., "Active", "Pending", "Sold")
	TransactionHistory  []TransactionRecord // List of previous transactions
	AuditTrail          []AuditRecord       // List of audit records for compliance
	EncryptedMetadata   []byte              // Encrypted metadata (optional details like contracts or additional descriptions)
	ImmutableRecords    []ImmutableRecord   // Immutable records for transparency and traceability
	ComplianceRecords   []ComplianceRecord  // Records of compliance with data regulations and standards
	DynamicPricingModel bool                // Whether the token follows a dynamic pricing model
}

// AccessRights defines the permissions for the data associated with the token.
type AccessRights struct {
	CanView       bool      // If the data can be viewed
	CanEdit       bool      // If the data can be edited
	CanShare      bool      // If the data can be shared with others
	ValidUntil    time.Time // Time limit for accessing the data
}

// TransactionRecord logs the details of a transaction involving the data token.
type TransactionRecord struct {
	TransactionID   string    // Unique identifier for the transaction
	Buyer           string    // Buyer of the data token
	Seller          string    // Seller of the data token
	TransactionDate time.Time // Date of transaction
	Price           float64   // Price at which the data was sold
	Status          string    // Status of the transaction (e.g., "Completed", "Pending")
}

// AuditRecord tracks actions taken on the data token for regulatory and compliance purposes.
type AuditRecord struct {
	Action     string    // Description of the action taken (e.g., "Updated Metadata", "Transferred Ownership")
	PerformedBy string   // The entity or user that performed the action
	Timestamp  time.Time // When the action occurred
	Details    string    // Additional details or notes on the action
}

// ImmutableRecord holds immutable records for compliance, transparency, and traceability.
type ImmutableRecord struct {
	RecordID    string    // Unique identifier for the immutable record
	Description string    // Description of the immutable record (e.g., "Data Created", "Ownership Transferred")
	Timestamp   time.Time // Time of the event being logged
}

// ComplianceRecord keeps track of the compliance status for the data token.
type ComplianceRecord struct {
	ComplianceID   string    // Unique identifier for the compliance record
	ComplianceType string    // Type of compliance (e.g., "GDPR", "CCPA")
	VerifiedBy     string    // The organization or authority verifying compliance
	VerificationDate time.Time // Date the compliance was verified
	Status         string    // Status of compliance (e.g., "Compliant", "Pending", "Non-Compliant")
}

// GetPrice returns the current price of the data token, factoring in dynamic pricing if enabled.
func (token *SYN2400Token) GetPrice() float64 {
	if token.DynamicPricingModel {
		// Implement dynamic pricing logic here (e.g., based on demand, time, etc.)
		// Example: increase price by 1% every day
		daysSinceCreation := time.Since(token.CreationDate).Hours() / 24
		return token.Price * (1 + 0.01*daysSinceCreation)
	}
	return token.Price
}

// UpdateMetadata updates the encrypted metadata for the data token and logs the action in the audit trail.
func (token *SYN2400Token) UpdateMetadata(newMetadata []byte, performedBy string) {
	token.EncryptedMetadata = newMetadata
	token.UpdateDate = time.Now()

	// Log the update in the audit trail
	token.AuditTrail = append(token.AuditTrail, AuditRecord{
		Action:     "Updated Metadata",
		PerformedBy: performedBy,
		Timestamp:  time.Now(),
		Details:    "Metadata updated for token ID " + token.TokenID,
	})
}

// TransferOwnership transfers ownership of the data token to a new owner and logs the transaction.
func (token *SYN2400Token) TransferOwnership(newOwner string, price float64, performedBy string) error {
	// Create a transaction record for the transfer
	transaction := TransactionRecord{
		TransactionID:   generateUniqueID(),
		Buyer:           newOwner,
		Seller:          token.Owner,
		TransactionDate: time.Now(),
		Price:           price,
		Status:          "Completed",
	}

	// Add the transaction to the transaction history
	token.TransactionHistory = append(token.TransactionHistory, transaction)

	// Update the owner of the token
	token.Owner = newOwner
	token.UpdateDate = time.Now()

	// Log the ownership transfer in the audit trail
	token.AuditTrail = append(token.AuditTrail, AuditRecord{
		Action:     "Transferred Ownership",
		PerformedBy: performedBy,
		Timestamp:  time.Now(),
		Details:    "Ownership transferred to " + newOwner + " for token ID " + token.TokenID,
	})

	return nil
}

// VerifyCompliance checks if the data token complies with the specified regulation and logs the result.
func (token *SYN2400Token) VerifyCompliance(complianceType, verifier string) error {
	// Create a compliance record
	complianceRecord := ComplianceRecord{
		ComplianceID:    generateUniqueID(),
		ComplianceType:  complianceType,
		VerifiedBy:      verifier,
		VerificationDate: time.Now(),
		Status:          "Compliant",
	}

	// Add the compliance record to the list
	token.ComplianceRecords = append(token.ComplianceRecords, complianceRecord)

	// Log the compliance verification in the audit trail
	token.AuditTrail = append(token.AuditTrail, AuditRecord{
		Action:     "Compliance Verified",
		PerformedBy: verifier,
		Timestamp:  time.Now(),
		Details:    "Compliance verified for " + complianceType + " by " + verifier,
	})

	return nil
}

// GetTransactionHistory returns the transaction history for the data token.
func (token *SYN2400Token) GetTransactionHistory() []TransactionRecord {
	return token.TransactionHistory
}

// GetAuditTrail returns the audit trail for the data token.
func (token *SYN2400Token) GetAuditTrail() []AuditRecord {
	return token.AuditTrail
}

// SYN2400Factory handles the creation, management, and encryption of SYN2400 tokens
type SYN2400Factory struct {
	Ledger   ledger.LedgerInterface        // Interface for interacting with the blockchain ledger
	Encrypt  encryption.EncryptionInterface // Interface for encryption
}

// NewSYN2400Factory initializes and returns a new instance of SYN2400Factory
func NewSYN2400Factory(ledger ledger.LedgerInterface, encrypt encryption.EncryptionInterface) *SYN2400Factory {
	return &SYN2400Factory{
		Ledger:  ledger,
		Encrypt: encrypt,
	}
}

// CreateToken creates a new SYN2400 token and stores it in the ledger
func (factory *SYN2400Factory) CreateToken(
	owner string,
	dataHash string,
	description string,
	accessRights common.AccessRights,
	price float64,
	dynamicPricingModel bool) (common.SYN2400Token, error) {

	// Generate unique token ID
	tokenID := generateUniqueID()

	// Create new token struct
	newToken := common.SYN2400Token{
		TokenID:             tokenID,
		Owner:               owner,
		DataHash:            dataHash,
		Description:         description,
		AccessRights:        accessRights,
		CreationDate:        time.Now(),
		UpdateDate:          time.Now(),
		Price:               price,
		Status:              "Active",
		DynamicPricingModel: dynamicPricingModel,
		TransactionHistory:  []common.TransactionRecord{},
		AuditTrail:          []common.AuditRecord{},
		ImmutableRecords:    []common.ImmutableRecord{},
		ComplianceRecords:   []common.ComplianceRecord{},
	}

	// Encrypt the token's metadata
	encryptedMetadata, err := factory.Encrypt.EncryptData([]byte(description))
	if err != nil {
		return common.SYN2400Token{}, err
	}
	newToken.EncryptedMetadata = encryptedMetadata

	// Store the token in the ledger
	if err := factory.Ledger.StoreToken(newToken.TokenID, newToken); err != nil {
		return common.SYN2400Token{}, err
	}

	// Log the creation event in the audit trail
	newToken.AuditTrail = append(newToken.AuditTrail, common.AuditRecord{
		Action:     "Token Created",
		PerformedBy: owner,
		Timestamp:  time.Now(),
		Details:    "Created SYN2400 token with ID: " + tokenID,
	})

	return newToken, nil
}

// UpdateToken allows the owner to update an existing SYN2400 token's metadata
func (factory *SYN2400Factory) UpdateToken(
	tokenID string,
	updatedDescription string,
	owner string) (common.SYN2400Token, error) {

	// Retrieve the token from the ledger
	token, err := factory.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Check ownership before updating
	if token.Owner != owner {
		return common.SYN2400Token{}, errors.New("unauthorized: only the token owner can update")
	}

	// Encrypt the new metadata
	encryptedMetadata, err := factory.Encrypt.EncryptData([]byte(updatedDescription))
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Update token metadata
	token.Description = updatedDescription
	token.EncryptedMetadata = encryptedMetadata
	token.UpdateDate = time.Now()

	// Log the update in the audit trail
	token.AuditTrail = append(token.AuditTrail, common.AuditRecord{
		Action:     "Token Updated",
		PerformedBy: owner,
		Timestamp:  time.Now(),
		Details:    "Updated description of token with ID: " + tokenID,
	})

	// Store updated token in ledger
	if err := factory.Ledger.UpdateToken(tokenID, token); err != nil {
		return common.SYN2400Token{}, err
	}

	return token, nil
}

// TransferToken handles ownership transfer of SYN2400 tokens between users
func (factory *SYN2400Factory) TransferToken(
	tokenID string,
	newOwner string,
	price float64,
	seller string) (common.SYN2400Token, error) {

	// Retrieve the token from the ledger
	token, err := factory.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Check ownership before transfer
	if token.Owner != seller {
		return common.SYN2400Token{}, errors.New("unauthorized: only the token owner can transfer ownership")
	}

	// Transfer ownership
	token.Owner = newOwner
	token.UpdateDate = time.Now()

	// Record the transfer in the transaction history
	transactionRecord := common.TransactionRecord{
		TransactionID:   generateUniqueID(),
		Buyer:           newOwner,
		Seller:          seller,
		TransactionDate: time.Now(),
		Price:           price,
		Status:          "Completed",
	}
	token.TransactionHistory = append(token.TransactionHistory, transactionRecord)

	// Log the transfer event in the audit trail
	token.AuditTrail = append(token.AuditTrail, common.AuditRecord{
		Action:     "Token Transferred",
		PerformedBy: seller,
		Timestamp:  time.Now(),
		Details:    "Transferred ownership of token to: " + newOwner,
	})

	// Store updated token in ledger
	if err := factory.Ledger.UpdateToken(tokenID, token); err != nil {
		return common.SYN2400Token{}, err
	}

	return token, nil
}

// VerifyCompliance verifies that a SYN2400 token meets a specific compliance standard
func (factory *SYN2400Factory) VerifyCompliance(
	tokenID string,
	complianceType string,
	verifier string) (common.SYN2400Token, error) {

	// Retrieve the token from the ledger
	token, err := factory.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Create compliance record
	complianceRecord := common.ComplianceRecord{
		ComplianceID:    generateUniqueID(),
		ComplianceType:  complianceType,
		VerifiedBy:      verifier,
		VerificationDate: time.Now(),
		Status:          "Compliant",
	}

	// Append to compliance records and log the action
	token.ComplianceRecords = append(token.ComplianceRecords, complianceRecord)
	token.AuditTrail = append(token.AuditTrail, common.AuditRecord{
		Action:     "Compliance Verified",
		PerformedBy: verifier,
		Timestamp:  time.Now(),
		Details:    "Verified compliance with " + complianceType,
	})

	// Store updated token in ledger
	if err := factory.Ledger.UpdateToken(tokenID, token); err != nil {
		return common.SYN2400Token{}, err
	}

	return token, nil
}

// Helper function to generate a unique ID for token creation and transactions
func generateUniqueID() string {
	return "SYN2400-" + time.Now().Format("20060102150405") + "-" + randomString(8)
}

// Helper function to generate a random string for unique identifiers
func randomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
	}
	return string(b)
}
