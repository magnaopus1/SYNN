package syn2100

import (
	"errors"
	"time"

)

// SYN2100Token represents a supply chain financing token under the SYN2100 standard
type SYN2100Token struct {
    TokenID               string              // Unique identifier for the token
    DocumentMetadata      FinancialDocumentMetadata // Metadata related to the financial document (invoice, purchase order, etc.)
    Owner                 string              // The current owner of the token (e.g., supplier, financier)
    Status                string              // Status of the token (e.g., "Issued", "Settled", "In Trade")
    DynamicDiscountRate   float64             // Dynamic discount rate for early payment
    IssueDate             time.Time           // The date when the financial document was tokenized
    DueDate               time.Time           // The due date for the financial document payment
    Amount                float64             // Amount associated with the financial document
    LiquidityPoolID       string              // The ID of the liquidity pool where this token is traded (if applicable)
    CollateralStatus      string              // Whether the token is backed by collateral (e.g., "Secured", "Unsecured")
    DocumentType          string              // The type of financial document (e.g., "Invoice", "Purchase Order")
    Tradeable             bool                // Whether the tokenized document can be traded in secondary markets
    AuditTrail            []AuditLog          // Detailed audit trail for all actions related to this token
    VerificationStatus    string              // Status of document verification (e.g., "Verified", "Pending", "Rejected")
    RestrictedTransfers   bool                // Indicates whether transfers of this token are restricted
    ApprovalRequired      bool                // Indicates whether transfers or high-value transactions require third-party approval
    RevenueShareModel     RevenueShare        // Information regarding the revenue sharing model for financiers
    RiskRating            string              // Risk rating based on issuer’s creditworthiness
    EncryptedMetadata     []byte              // Encrypted metadata containing sensitive information (e.g., terms of finance, contracts)
}

// FinancialDocumentMetadata represents metadata of the financial document that the SYN2100 token is based on.
type FinancialDocumentMetadata struct {
    DocumentID     string    // Unique identifier for the financial document
    Issuer         string    // Entity issuing the financial document (e.g., supplier, company)
    Recipient      string    // Entity receiving the financial document (e.g., customer, buyer)
    Amount         float64   // Total value of the financial document
    Description    string    // Description of the financial document or transaction
    IssueDate      time.Time // Date the document was issued
    DueDate        time.Time // Due date for payment
    DocumentType   string    // Type of document (e.g., invoice, purchase order)
    SourcePlatform string    // Platform or system from which the document originates
}

// AuditLog represents a record of actions taken on the tokenized financial document
type AuditLog struct {
    LogID          string    // Unique identifier for the audit log
    Action         string    // Action performed (e.g., "Issued", "Transferred", "Verified")
    PerformedBy    string    // Entity that performed the action
    Timestamp      time.Time // Time when the action was performed
    Description    string    // Additional details or notes on the action
}

// RevenueShare captures details on how revenue from financed invoices is shared with financiers
type RevenueShare struct {
    FinancierID    string    // ID of the financier participating in revenue sharing
    SharePercentage float64  // Percentage of revenue that the financier will earn from the tokenized invoice
    RevenueEarned  float64   // Total revenue earned by the financier from this token
    PayoutSchedule string    // Schedule of payouts (e.g., "Monthly", "Quarterly")
    PayoutDate     time.Time // Date of the last revenue payout
}

// AddAuditLog adds a new audit log entry to the token
func (token *SYN2100Token) AddAuditLog(action, performedBy, description string) {
    auditLog := AuditLog{
        LogID:       generateUniqueID(),
        Action:      action,
        PerformedBy: performedBy,
        Timestamp:   time.Now(),
        Description: description,
    }
    token.AuditTrail = append(token.AuditTrail, auditLog)
}

// AddRevenueShare adds revenue sharing details to the token for a financier
func (token *SYN2100Token) AddRevenueShare(financierID string, percentage float64, schedule string) {
    revenueShare := RevenueShare{
        FinancierID:    financierID,
        SharePercentage: percentage,
        PayoutSchedule: schedule,
        PayoutDate:     time.Now(),
    }
    token.RevenueShareModel = revenueShare
}

// UpdateVerificationStatus updates the verification status of the financial document
func (token *SYN2100Token) UpdateVerificationStatus(status string) {
    token.VerificationStatus = status
}

// GenerateUniqueID generates a unique identifier for tokens or logs
func generateUniqueID() string {
    return fmt.Sprintf("ID-%d", time.Now().UnixNano())
}


// CreateSYN2100Token creates and initializes a new SYN2100 token for a financial document
func CreateSYN2100Token(document common.FinancialDocumentMetadata, owner string, discountRate float64) (*common.SYN2100Token, error) {
	// Validate input data
	if document.DocumentID == "" || document.Issuer == "" || document.Recipient == "" || document.Amount <= 0 {
		return nil, errors.New("invalid document metadata")
	}

	// Create a new SYN2100 token
	token := &common.SYN2100Token{
		TokenID:             generateUniqueID(),
		DocumentMetadata:    document,
		Owner:               owner,
		Status:              "Issued",
		DynamicDiscountRate: discountRate,
		IssueDate:           time.Now(),
		DueDate:             document.DueDate,
		Amount:              document.Amount,
		Tradeable:           true,
		RestrictedTransfers: false,
		VerificationStatus:  "Pending",
	}

	// Encrypt sensitive data
	token.EncryptedMetadata = encryption.Encrypt([]byte("Sensitive contract data"))

	// Log the creation of the token
	token.AddAuditLog("Created", owner, "SYN2100 token created for financial document")

	// Integrate the token with the ledger for validation and tracking
	err := ledger.RecordNewToken(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// TransferSYN2100Token transfers ownership of the token to a new owner
func TransferSYN2100Token(token *common.SYN2100Token, newOwner string) error {
	if !token.Tradeable {
		return errors.New("token is not tradeable")
	}

	// Check if transfers are restricted and require approval
	if token.RestrictedTransfers && !token.ApprovalRequired {
		return errors.New("transfers require approval")
	}

	// Update ownership
	oldOwner := token.Owner
	token.Owner = newOwner

	// Log the transfer
	token.AddAuditLog("Transferred", oldOwner, "Token transferred to new owner")

	// Update ledger with new ownership details
	err := ledger.UpdateTokenOwnership(token.TokenID, newOwner)
	if err != nil {
		return err
	}

	return nil
}

// AddDynamicDiscount adds or updates the dynamic discount rate for early payment
func AddDynamicDiscount(token *common.SYN2100Token, newRate float64) error {
	if newRate < 0 {
		return errors.New("invalid discount rate")
	}

	token.DynamicDiscountRate = newRate
	token.AddAuditLog("Updated Discount", token.Owner, "Dynamic discount rate updated")

	// Record the update in the ledger
	err := ledger.UpdateTokenDiscountRate(token.TokenID, newRate)
	if err != nil {
		return err
	}

	return nil
}

// SetRestrictedTransfers updates whether the token allows restricted transfers
func SetRestrictedTransfers(token *common.SYN2100Token, restricted bool) {
	token.RestrictedTransfers = restricted
	token.AddAuditLog("Updated Restrictions", token.Owner, "Transfer restrictions updated")
}

// VerifyToken sets the verification status of the token
func VerifyToken(token *common.SYN2100Token, verifiedBy string) error {
	token.VerificationStatus = "Verified"
	token.AddAuditLog("Verified", verifiedBy, "Token verification completed")

	// Update ledger with verification status
	err := ledger.VerifyToken(token.TokenID)
	if err != nil {
		return err
	}

	return nil
}

// generateUniqueID generates a unique identifier for tokens or logs
func generateUniqueID() string {
	return fmt.Sprintf("ID-%d", time.Now().UnixNano())
}
