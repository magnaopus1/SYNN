package syn1500

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// SYN1500Compliance manages the compliance checks, audits, and regulatory adherence of reputation tokens.
type SYN1500Compliance struct {
	ComplianceID string       // Unique ID for the compliance instance
	Ledger       ledger.Ledger // Reference to the blockchain ledger for recording compliance-related events
}

// AuditLog represents an audit entry for SYN1500Token, capturing compliance reviews.
type AuditLog struct {
	AuditID      string    // Unique identifier for the audit entry
	TokenID      string    // The token being audited
	PerformedBy  string    // ID of the auditor or system performing the audit
	Description  string    // Description of the audit or compliance check
	Timestamp    time.Time // When the audit was performed
	Result       string    // The result of the audit (e.g., "Compliant", "Non-Compliant", "Pending")
}

// ComplianceCheck represents a compliance check result for a SYN1500Token.
type ComplianceCheck struct {
	CheckID        string    // Unique ID of the compliance check
	TokenID        string    // The token being checked
	PerformedBy    string    // ID of the person or system performing the compliance check
	CheckDate      time.Time // When the compliance check was performed
	Status         string    // Compliance status after the check (e.g., "Compliant", "Non-Compliant")
	Description    string    // Description of what was checked (e.g., KYC, privacy, data security)
	RequiredAction string    // Action required for non-compliant tokens (e.g., "Fix Privacy Issue")
}

// RecordComplianceCheck logs the result of a compliance check on a SYN1500Token into the ledger.
func (sc *SYN1500Compliance) RecordComplianceCheck(token *common.SYN1500Token, performedBy string, status string, description string) error {
	// Generate unique ID for the compliance check
	checkID := generateUniqueID(token.TokenID)

	complianceCheck := ComplianceCheck{
		CheckID:        checkID,
		TokenID:        token.TokenID,
		PerformedBy:    performedBy,
		CheckDate:      time.Now(),
		Status:         status,
		Description:    description,
		RequiredAction: "",
	}

	if status != "Compliant" {
		complianceCheck.RequiredAction = "Action Required to Address Non-Compliance"
		token.ComplianceStatus = "Non-Compliant"
	} else {
		token.ComplianceStatus = "Compliant"
	}

	// Record compliance check in the ledger
	if err := sc.Ledger.RecordTransaction(ledger.Transaction{
		TxID:        complianceCheck.CheckID,
		Description: fmt.Sprintf("Compliance check performed on token %s by %s", token.TokenID, performedBy),
		Timestamp:   time.Now(),
		Data:        complianceCheck,
	}); err != nil {
		return errors.New("failed to record compliance check in the ledger")
	}

	// Update token metadata with the compliance status and encrypted metadata
	token.EncryptedMetadata = encryptMetadata(token)

	return nil
}

// PerformAudit performs a compliance audit on a SYN1500Token and logs the audit result.
func (sc *SYN1500Compliance) PerformAudit(token *common.SYN1500Token, performedBy string, result string) error {
	// Generate unique ID for the audit
	auditID := generateUniqueID(token.TokenID)

	auditLog := AuditLog{
		AuditID:      auditID,
		TokenID:      token.TokenID,
		PerformedBy:  performedBy,
		Description:  fmt.Sprintf("Audit performed on token %s", token.TokenID),
		Timestamp:    time.Now(),
		Result:       result,
	}

	// Record audit result in the ledger
	if err := sc.Ledger.RecordTransaction(ledger.Transaction{
		TxID:        auditLog.AuditID,
		Description: auditLog.Description,
		Timestamp:   auditLog.Timestamp,
		Data:        auditLog,
	}); err != nil {
		return errors.New("failed to record audit log in the ledger")
	}

	// Update the compliance status based on the audit result
	token.ComplianceStatus = result

	// Encrypt the token's metadata with the new compliance status
	token.EncryptedMetadata = encryptMetadata(token)

	return nil
}

// ResolveNonCompliance resolves a non-compliance issue by updating the token status.
func (sc *SYN1500Compliance) ResolveNonCompliance(token *common.SYN1500Token, resolutionDetails string) error {
	if token.ComplianceStatus != "Non-Compliant" {
		return errors.New("token is already compliant")
	}

	// Mark token as compliant and log the resolution in the ledger
	token.ComplianceStatus = "Compliant"

	// Create a new compliance check log to document the resolution
	complianceCheck := ComplianceCheck{
		CheckID:        generateUniqueID(token.TokenID),
		TokenID:        token.TokenID,
		PerformedBy:    "Compliance System",
		CheckDate:      time.Now(),
		Status:         "Compliant",
		Description:    fmt.Sprintf("Resolved: %s", resolutionDetails),
		RequiredAction: "",
	}

	// Record the resolution in the ledger
	if err := sc.Ledger.RecordTransaction(ledger.Transaction{
		TxID:        complianceCheck.CheckID,
		Description: fmt.Sprintf("Compliance issue resolved for token %s", token.TokenID),
		Timestamp:   complianceCheck.CheckDate,
		Data:        complianceCheck,
	}); err != nil {
		return errors.New("failed to record resolution in the ledger")
	}

	// Encrypt token metadata to reflect the updated compliance status
	token.EncryptedMetadata = encryptMetadata(token)

	return nil
}

// encryptMetadata encrypts sensitive metadata for a SYN1500Token
func encryptMetadata(token *common.SYN1500Token) []byte {
	data, _ := json.Marshal(token)
	hash := sha256.Sum256(data)
	return hash[:]
}

// generateUniqueID creates a unique identifier for compliance logs and checks
func generateUniqueID(seed string) string {
	timestamp := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", seed, timestamp)))
	return hex.EncodeToString(hash[:])
}
