package syn1900

import (
	"errors"
	"time"
)


// ManagementService handles management operations for SYN1900 tokens (education credits).
type ManagementService struct {
	ledger LedgerInterface // Interface for interacting with the ledger
}

// LedgerInterface defines methods for interacting with the ledger.
type LedgerInterface interface {
	GetTokenByID(tokenID string) (common.SYN1900Token, error)
	UpdateToken(token common.SYN1900Token) error
}

// IssueNewCredit issues a new educational credit token and adds it to the ledger.
func (ms *ManagementService) IssueNewCredit(
	recipientID string,
	courseID string,
	courseName string,
	issuer string,
	creditValue float64,
	expirationDate time.Time,
	metadata string,
) (string, error) {
	// Create new token
	newToken := common.SYN1900Token{
		TokenID:       generateUniqueID(),
		RecipientID:   recipientID,
		CourseID:      courseID,
		CourseName:    courseName,
		Issuer:        issuer,
		CreditValue:   creditValue,
		IssueDate:     time.Now(),
		ExpirationDate: expirationDate,
		Metadata:      metadata,
		Signature:     encryption.GenerateDigitalSignature(issuer),
	}

	// Add the new token to the ledger
	err := ms.ledger.UpdateToken(newToken)
	if err != nil {
		return "", errors.New("failed to issue and add new credit to the ledger")
	}

	return newToken.TokenID, nil
}

// TransferCredit transfers a credit token from one recipient to another.
func (ms *ManagementService) TransferCredit(tokenID, fromID, toID string) error {
	// Fetch the token from the ledger
	token, err := ms.ledger.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("token not found in the ledger")
	}

	// Ensure that the transfer is authorized
	if token.RecipientID != fromID {
		return errors.New("only the current recipient can transfer the token")
	}

	// Update the recipient
	token.RecipientID = toID

	// Update the ledger
	err = ms.ledger.UpdateToken(token)
	if err != nil {
		return errors.New("failed to transfer the token in the ledger")
	}

	return nil
}

// RevokeCredit revokes a credit token due to invalidation or academic dishonesty.
func (ms *ManagementService) RevokeCredit(tokenID, revocationReason string) error {
	// Fetch the token from the ledger
	token, err := ms.ledger.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("token not found in the ledger")
	}

	// Set the revocation status
	token.Revoked = true
	token.RevocationReason = revocationReason
	token.RevocationDate = time.Now()

	// Update the ledger
	err = ms.ledger.UpdateToken(token)
	if err != nil {
		return errors.New("failed to revoke the token in the ledger")
	}

	return nil
}

// ValidateCredit verifies the authenticity and status of a credit token.
func (ms *ManagementService) ValidateCredit(tokenID string) (bool, error) {
	// Fetch the token from the ledger
	token, err := ms.ledger.GetTokenByID(tokenID)
	if err != nil {
		return false, errors.New("token not found in the ledger")
	}

	// Check if the token is valid (not revoked and not expired)
	if token.Revoked || token.ExpirationDate.Before(time.Now()) {
		return false, nil
	}

	// Verify the digital signature of the token issuer
	isValidSignature := encryption.VerifyDigitalSignature(token.Issuer, token.Signature)
	if !isValidSignature {
		return false, errors.New("failed to verify the issuer's digital signature")
	}

	return true, nil
}

// UpdateMetadata updates the metadata of an existing token.
func (ms *ManagementService) UpdateMetadata(tokenID string, newMetadata string) error {
	// Fetch the token from the ledger
	token, err := ms.ledger.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("token not found in the ledger")
	}

	// Update the metadata
	token.Metadata = newMetadata

	// Update the ledger
	err = ms.ledger.UpdateToken(token)
	if err != nil {
		return errors.New("failed to update token metadata in the ledger")
	}

	return nil
}

// RenewCredit extends the expiration date of an educational credit.
func (ms *ManagementService) RenewCredit(tokenID string, newExpirationDate time.Time) error {
	// Fetch the token from the ledger
	token, err := ms.ledger.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("token not found in the ledger")
	}

	// Update the expiration date
	token.ExpirationDate = newExpirationDate

	// Update the ledger
	err = ms.ledger.UpdateToken(token)
	if err != nil {
		return errors.New("failed to renew the token in the ledger")
	}

	return nil
}

// GenerateSummary generates a summary of the token's details, including status and key events.
func (ms *ManagementService) GenerateSummary(tokenID string) (string, error) {
	// Fetch the token from the ledger
	token, err := ms.ledger.GetTokenByID(tokenID)
	if err != nil {
		return "", errors.New("token not found in the ledger")
	}

	// Generate the summary
	summary := "Education Credit Summary:\n"
	summary += "Token ID: " + token.TokenID + "\n"
	summary += "Recipient ID: " + token.RecipientID + "\n"
	summary += "Course ID: " + token.CourseID + "\n"
	summary += "Course Name: " + token.CourseName + "\n"
	summary += "Issuer: " + token.Issuer + "\n"
	summary += "Credit Value: " + string(token.CreditValue) + "\n"
	summary += "Issue Date: " + token.IssueDate.Format(time.RFC3339) + "\n"
	summary += "Expiration Date: " + token.ExpirationDate.Format(time.RFC3339) + "\n"
	summary += "Revoked: " + string(token.Revoked) + "\n"
	if token.Revoked {
		summary += "Revocation Reason: " + token.RevocationReason + "\n"
		summary += "Revocation Date: " + token.RevocationDate.Format(time.RFC3339) + "\n"
	}
	summary += "Metadata: " + token.Metadata + "\n"

	return summary, nil
}
