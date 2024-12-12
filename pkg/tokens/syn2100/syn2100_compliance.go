package syn2100

import (
	"errors"
	"time"

)

// CheckComplianceStatus checks the current compliance status of a SYN2100 token.
func CheckComplianceStatus(token *common.SYN2100Token) (string, error) {
	// Retrieve compliance data from the ledger
	complianceStatus, err := ledger.GetTokenComplianceStatus(token.TokenID)
	if err != nil {
		return "", errors.New("failed to retrieve compliance status from ledger")
	}

	// Return the compliance status of the token
	return complianceStatus, nil
}

// UpdateComplianceStatus updates the compliance status of a SYN2100 token.
func UpdateComplianceStatus(token *common.SYN2100Token, status string, updatedBy string) error {
	// Validate the new status
	if status != "Compliant" && status != "Non-Compliant" && status != "Pending" {
		return errors.New("invalid compliance status")
	}

	// Update the token's compliance status
	token.ComplianceStatus = status
	token.AddAuditLog("Compliance Update", updatedBy, "Token compliance status updated to "+status)

	// Update the ledger with the new compliance status
	err := ledger.UpdateTokenComplianceStatus(token.TokenID, status)
	if err != nil {
		return errors.New("failed to update compliance status in ledger")
	}

	return nil
}

// VerifyComplianceDocument verifies and encrypts a compliance document for the token.
func VerifyComplianceDocument(token *common.SYN2100Token, documentID string, verifier string) error {
	// Check if the document ID exists in the token's compliance records
	found := false
	for _, doc := range token.ComplianceRecords {
		if doc.DocumentID == documentID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("compliance document not found")
	}

	// Encrypt the compliance document to secure sensitive information
	encryptedDocument, err := encryption.Encrypt([]byte(documentID))
	if err != nil {
		return errors.New("failed to encrypt compliance document")
	}

	// Store the encrypted document back into the compliance record
	token.AddComplianceDocumentRecord(documentID, encryptedDocument, verifier)

	// Log the verification event
	token.AddAuditLog("Compliance Verified", verifier, "Compliance document verified and encrypted")

	// Update the ledger with the verification status of the document
	err = ledger.VerifyComplianceDocument(token.TokenID, documentID)
	if err != nil {
		return errors.New("failed to verify compliance document in ledger")
	}

	return nil
}

// AddComplianceDocumentRecord adds a new compliance document record to the token.
func (token *common.SYN2100Token) AddComplianceDocumentRecord(documentID string, encryptedDocument []byte, verifiedBy string) {
	newRecord := common.ComplianceDocumentRecord{
		DocumentID:        documentID,
		EncryptedDocument: encryptedDocument,
		VerifiedBy:        verifiedBy,
		VerificationDate:  time.Now(),
	}

	token.ComplianceRecords = append(token.ComplianceRecords, newRecord)
}

// RemoveComplianceViolation clears the compliance violation record of a token.
func RemoveComplianceViolation(token *common.SYN2100Token, violationID string) error {
	// Find the violation in the compliance records
	var newRecords []common.ComplianceDocumentRecord
	for _, record := range token.ComplianceRecords {
		if record.DocumentID != violationID {
			newRecords = append(newRecords, record)
		}
	}

	if len(newRecords) == len(token.ComplianceRecords) {
		return errors.New("violation record not found")
	}

	// Update token compliance records
	token.ComplianceRecords = newRecords

	// Log the removal of the violation
	token.AddAuditLog("Compliance Violation Cleared", token.Owner, "Compliance violation cleared for "+violationID)

	// Update ledger to reflect the removal of the violation
	err := ledger.RemoveComplianceViolation(token.TokenID, violationID)
	if err != nil {
		return errors.New("failed to remove compliance violation from ledger")
	}

	return nil
}

// AddAuditLog adds a new audit log entry for the token's compliance actions.
func (token *common.SYN2100Token) AddAuditLog(actionType, performedBy, description string) {
	logEntry := common.AuditLog{
		ActionType:  actionType,
		PerformedBy: performedBy,
		Timestamp:   time.Now(),
		Description: description,
	}

	token.AuditTrail = append(token.AuditTrail, logEntry)
}
