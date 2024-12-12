package syn2100

import (
	"errors"
	"time"

)

// CreateSYN2100Token creates a new SYN2100 token representing a financial document in the supply chain financing system.
func CreateSYN2100Token(document common.FinancialDocumentMetadata, owner string) (*common.SYN2100Token, error) {
	// Encrypt sensitive document information
	encryptedDocID, err := encryption.Encrypt([]byte(document.DocumentID))
	if err != nil {
		return nil, errors.New("failed to encrypt document ID")
	}

	// Create a new token
	token := &common.SYN2100Token{
		TokenID:           generateUniqueID(),
		DocumentMetadata:  document,
		Owner:             owner,
		IssuedDate:        time.Now(),
		Status:            "Active",
		AuditTrail:        []common.AuditLog{},
		EncryptedDocID:    encryptedDocID,
	}

	// Record the creation event in the token's audit trail
	creationEvent := common.AuditLog{
		EventID:     generateUniqueID(),
		EventType:   "Token Creation",
		Description: "SYN2100 Token created for financial document: " + document.DocumentID,
		PerformedBy: owner,
		EventDate:   time.Now(),
	}
	token.AuditTrail = append(token.AuditTrail, creationEvent)

	// Store the token creation in the ledger
	err = ledger.RecordTokenCreation(token.TokenID, token)
	if err != nil {
		return nil, errors.New("failed to record token creation in ledger")
	}

	return token, nil
}

// TransferSYN2100Token transfers the ownership of a SYN2100 token from one party to another.
func TransferSYN2100Token(token *common.SYN2100Token, newOwner string) error {
	// Validate token status before transfer
	if token.Status != "Active" {
		return errors.New("token is not active and cannot be transferred")
	}

	// Record the transfer event
	transferEvent := common.AuditLog{
		EventID:     generateUniqueID(),
		EventType:   "Token Transfer",
		Description: "Transferred token ownership to: " + newOwner,
		PerformedBy: token.Owner,
		EventDate:   time.Now(),
	}

	// Update the token's audit trail
	token.AuditTrail = append(token.AuditTrail, transferEvent)

	// Update the ledger to record the transfer
	err := ledger.RecordTransfer(token.TokenID, token.Owner, newOwner)
	if err != nil {
		return errors.New("failed to record transfer in ledger")
	}

	// Change the owner in the token
	token.Owner = newOwner

	return nil
}

// ApplyDynamicDiscounting applies a dynamic discount to the tokenized invoice based on early payment.
func ApplyDynamicDiscounting(token *common.SYN2100Token, originalAmount float64, discountRate float64, issuer string) error {
	// Calculate the discounted amount
	discountedAmount := originalAmount - (originalAmount * discountRate / 100)

	// Record the dynamic discounting event in the audit trail
	discountingEvent := common.AuditLog{
		EventID:     generateUniqueID(),
		EventType:   "Dynamic Discounting Applied",
		Description: "Applied dynamic discounting. Original amount: " + formatAmount(originalAmount) + ", discounted to: " + formatAmount(discountedAmount),
		PerformedBy: issuer,
		EventDate:   time.Now(),
	}

	// Add the event to the token's audit trail
	token.AuditTrail = append(token.AuditTrail, discountingEvent)

	// Update the ledger with the dynamic discounting event
	err := ledger.RecordEvent(token.TokenID, "Dynamic Discounting", discountingEvent)
	if err != nil {
		return errors.New("failed to record discounting event in ledger")
	}

	return nil
}

// SettleSYN2100Token finalizes the settlement of a tokenized invoice.
func SettleSYN2100Token(token *common.SYN2100Token, settlementAmount float64, settledBy string) error {
	// Ensure the token is active and eligible for settlement
	if token.Status != "Active" {
		return errors.New("token is not active and cannot be settled")
	}

	// Record the settlement event
	settlementEvent := common.AuditLog{
		EventID:     generateUniqueID(),
		EventType:   "Token Settlement",
		Description: "Settled token for amount: " + formatAmount(settlementAmount),
		PerformedBy: settledBy,
		EventDate:   time.Now(),
	}

	// Update the token's audit trail
	token.AuditTrail = append(token.AuditTrail, settlementEvent)

	// Mark the token as settled
	token.Status = "Settled"

	// Update the ledger with the settlement event
	err := ledger.RecordEvent(token.TokenID, "Token Settlement", settlementEvent)
	if err != nil {
		return errors.New("failed to record settlement event in ledger")
	}

	return nil
}

// RevokeSYN2100Token revokes a token, preventing further actions from being taken on it.
func RevokeSYN2100Token(token *common.SYN2100Token, revoker string, reason string) error {
	// Ensure the token is active before revocation
	if token.Status != "Active" {
		return errors.New("only active tokens can be revoked")
	}

	// Record the revocation event
	revocationEvent := common.AuditLog{
		EventID:     generateUniqueID(),
		EventType:   "Token Revocation",
		Description: "Token revoked by " + revoker + ". Reason: " + reason,
		PerformedBy: revoker,
		EventDate:   time.Now(),
	}

	// Update the token's audit trail
	token.AuditTrail = append(token.AuditTrail, revocationEvent)

	// Mark the token as revoked
	token.Status = "Revoked"

	// Update the ledger to reflect the revocation
	err := ledger.RecordEvent(token.TokenID, "Token Revocation", revocationEvent)
	if err != nil {
		return errors.New("failed to record revocation in ledger")
	}

	return nil
}

// AuditSYN2100Token retrieves the full audit trail for a SYN2100 token.
func AuditSYN2100Token(token *common.SYN2100Token) ([]common.AuditLog, error) {
	// Return the complete audit trail
	return token.AuditTrail, nil
}

// Utility function to generate a unique ID (to be replaced with a real implementation)
func generateUniqueID() string {
	// This is a placeholder for unique ID generation
	return "unique-id"
}

// Utility function to format amounts (for better readability)
func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
