package syn1900

import (
	"errors"
	"time"
)


// SYN1900Token represents an education credit token under the SYN1900 standard.
type SYN1900Token struct {
	TokenID              string               // Unique identifier for the token
	CreditMetadata       EducationCreditMetadata // Detailed metadata about the education credit
	Issuer               IssuerRecord         // Record of the issuer of the credit (e.g., institution)
	Recipient            RecipientRecord      // Record of the recipient of the credit
	VerificationLogs     []VerificationLog    // Logs of verification activities for this credit
	OwnershipHistory     []OwnershipRecord    // Historical records of ownership or credit transfers
	CertificationStatus  string               // Certification status of the credit (e.g., "Active", "Expired", "Revoked")
	ExpirationDate       *time.Time           // Optional expiration date for the credit
	RevocationStatus     bool                 // Indicates whether the credit has been revoked
	ImmutableRecords     []ImmutableRecord    // Immutable records for compliance and transparency
	EncryptedMetadata    []byte               // Encrypted sensitive metadata (e.g., certifications, agreements)
	TransferRestrictions bool                 // Whether there are transfer restrictions on this credit
	ApprovalRequired     bool                 // Whether certain actions require approval (e.g., high-value transfers)
}

// EducationCreditMetadata holds detailed information about an education credit.
type EducationCreditMetadata struct {
	CreditID       string    // Unique identifier for the credit
	CourseID       string    // Course ID associated with the credit
	CourseName     string    // Name of the course or program
	CreditValue    float64   // Value of the educational credit (e.g., credit hours, units)
	IssueDate      time.Time // Date when the credit was issued
	ExpirationDate *time.Time // Optional expiration date for the credit
	Metadata       string    // Additional metadata (e.g., certification level, grade)
	DigitalSignature []byte  // Digital signature to verify the authenticity of the credit
}

// IssuerRecord contains information about the issuer of an educational credit.
type IssuerRecord struct {
	IssuerID      string    // Unique identifier for the issuer
	IssuerName    string    // Name of the issuing organization or platform
	IssuerType    string    // Type of issuer (e.g., "University", "Online Platform")
	IssueDate     time.Time // Date the issuer created the credit
	ContactInfo   string    // Contact information for the issuer
	VerificationStatus string // Verification status of the issuer (e.g., "Verified", "Pending")
}

// RecipientRecord contains information about the recipient of an educational credit.
type RecipientRecord struct {
	RecipientID   string    // Unique identifier for the recipient
	RecipientName string    // Name of the recipient
	ReceivedDate  time.Time // Date when the recipient earned the credit
}

// OwnershipRecord captures the history of ownership or credit transfers.
type OwnershipRecord struct {
	PreviousOwner string    // ID of the previous owner of the credit
	NewOwner      string    // ID of the new owner
	TransferDate  time.Time // Date of the transfer
	TransferType  string    // Type of transfer (e.g., "Institutional Transfer", "Credit Transfer")
}

// VerificationLog captures verification activities related to an educational credit.
type VerificationLog struct {
	VerificationID string    // Unique identifier for the verification entry
	Verifier       string    // ID of the entity or person verifying the credit
	VerificationDate time.Time // Date of the verification activity
	Description    string    // Description of the verification process
	Status         string    // Status of the verification (e.g., "Verified", "Pending")
}

// ImmutableRecord stores immutable records for compliance and transparency.
type ImmutableRecord struct {
	RecordID    string    // Unique identifier for the record
	Description string    // Description of the record (e.g., "Credit Issued", "Transfer")
	Timestamp   time.Time // Time when the record was created
}

// Methods to interact with the SYN1900Token struct

// AddVerificationLog adds a new verification log to the credit's record.
func (token *SYN1900Token) AddVerificationLog(verifier string, description string, status string) {
	newLog := VerificationLog{
		VerificationID:  generateUniqueID(),
		Verifier:        verifier,
		VerificationDate: time.Now(),
		Description:     description,
		Status:          status,
	}
	token.VerificationLogs = append(token.VerificationLogs, newLog)
}

// TransferOwnership transfers the ownership of an educational credit to a new owner.
func (token *SYN1900Token) TransferOwnership(newOwner string, transferType string) {
	newOwnershipRecord := OwnershipRecord{
		PreviousOwner: token.Recipient.RecipientID,
		NewOwner:      newOwner,
		TransferDate:  time.Now(),
		TransferType:  transferType,
	}
	token.OwnershipHistory = append(token.OwnershipHistory, newOwnershipRecord)
	token.Recipient.RecipientID = newOwner
}

// RevokeCredit revokes the educational credit, preventing it from being used.
func (token *SYN1900Token) RevokeCredit(reason string) {
	token.RevocationStatus = true
	newRecord := ImmutableRecord{
		RecordID:    generateUniqueID(),
		Description: fmt.Sprintf("Credit revoked: %s", reason),
		Timestamp:   time.Now(),
	}
	token.ImmutableRecords = append(token.ImmutableRecords, newRecord)
}

// ExpireCredit sets the expiration date of an educational credit.
func (token *SYN1900Token) ExpireCredit(expirationDate time.Time) {
	token.ExpirationDate = &expirationDate
	newRecord := ImmutableRecord{
		RecordID:    generateUniqueID(),
		Description: "Credit expiration set",
		Timestamp:   time.Now(),
	}
	token.ImmutableRecords = append(token.ImmutableRecords, newRecord)
}

// Helper function to generate unique IDs.
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}


// TokenFactory is responsible for creating and managing SYN1900 tokens.
type TokenFactory struct {
	ledger ledger.Ledger // Interface for interacting with the ledger
}

// CreateSYN1900Token creates a new SYN1900 education credit token.
func (factory *TokenFactory) CreateSYN1900Token(issuerID, recipientID, courseID, courseName string, creditValue float64, metadata string, signature []byte, expirationDate *time.Time) (common.SYN1900Token, error) {
	// Validate input
	if issuerID == "" || recipientID == "" || courseID == "" || creditValue <= 0 {
		return common.SYN1900Token{}, errors.New("invalid input for creating SYN1900Token")
	}

	// Generate a unique token ID
	tokenID := generateUniqueID()

	// Create a new SYN1900Token
	token := common.SYN1900Token{
		TokenID: tokenID,
		CreditMetadata: common.EducationCreditMetadata{
			CreditID:        generateUniqueID(),
			CourseID:        courseID,
			CourseName:      courseName,
			CreditValue:     creditValue,
			IssueDate:       time.Now(),
			ExpirationDate:  expirationDate,
			Metadata:        metadata,
			DigitalSignature: signature,
		},
		Issuer: common.IssuerRecord{
			IssuerID:      issuerID,
			IssuerName:    "Example University", // You can get this from a DB or config
			IssuerType:    "University",
			IssueDate:     time.Now(),
			VerificationStatus: "Verified",
		},
		Recipient: common.RecipientRecord{
			RecipientID:   recipientID,
			RecipientName: "John Doe", // Replace with actual recipient name
			ReceivedDate:  time.Now(),
		},
		CertificationStatus: "Active",
		ExpirationDate:      expirationDate,
		RevocationStatus:    false,
	}

	// Encrypt sensitive metadata (e.g., certifications)
	encryptedMetadata, err := encryption.Encrypt([]byte(metadata))
	if err != nil {
		return common.SYN1900Token{}, errors.New("failed to encrypt metadata")
	}
	token.EncryptedMetadata = encryptedMetadata

	// Add the new token to the ledger
	err = factory.ledger.AddTransaction(common.Transaction{
		Type:        "CreateToken",
		TokenID:     tokenID,
		IssuerID:    issuerID,
		RecipientID: recipientID,
		Timestamp:   time.Now(),
	})
	if err != nil {
		return common.SYN1900Token{}, err
	}

	// Return the created token
	return token, nil
}

// RevokeSYN1900Token revokes an existing SYN1900 token.
func (factory *TokenFactory) RevokeSYN1900Token(tokenID, reason string) error {
	// Fetch the token from the ledger
	token, err := factory.ledger.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token")
	}

	// Check if the token is already revoked
	if token.RevocationStatus {
		return errors.New("token is already revoked")
	}

	// Revoke the token and update its status
	token.RevocationStatus = true
	token.RevokeCredit(reason)

	// Update the token in the ledger
	err = factory.ledger.UpdateToken(token)
	if err != nil {
		return errors.New("failed to update token in ledger")
	}

	// Log the revocation in the ledger
	err = factory.ledger.AddTransaction(common.Transaction{
		Type:      "RevokeToken",
		TokenID:   tokenID,
		IssuerID:  token.Issuer.IssuerID,
		Timestamp: time.Now(),
		Details:   reason,
	})
	if err != nil {
		return errors.New("failed to log revocation in ledger")
	}

	return nil
}

// TransferSYN1900Token transfers the ownership of an SYN1900 token to a new recipient.
func (factory *TokenFactory) TransferSYN1900Token(tokenID, newRecipientID, transferType string) error {
	// Fetch the token from the ledger
	token, err := factory.ledger.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token")
	}

	// Update the ownership
	token.TransferOwnership(newRecipientID, transferType)

	// Update the token in the ledger
	err = factory.ledger.UpdateToken(token)
	if err != nil {
		return errors.New("failed to update token in ledger")
	}

	// Log the transfer in the ledger
	err = factory.ledger.AddTransaction(common.Transaction{
		Type:        "TransferToken",
		TokenID:     tokenID,
		IssuerID:    token.Issuer.IssuerID,
		RecipientID: newRecipientID,
		Timestamp:   time.Now(),
		Details:     "Ownership transfer",
	})
	if err != nil {
		return errors.New("failed to log transfer in ledger")
	}

	return nil
}

// Helper function to generate unique IDs.
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
