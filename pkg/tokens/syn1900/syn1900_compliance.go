package syn1900

import (
	"errors"
	"time"
)

// ComplianceManager handles the compliance of SYN1900 education tokens.
type ComplianceManager struct {
	ledger LedgerInterface // Interface for interacting with the ledger
}

// LedgerInterface defines methods for interacting with the ledger.
type LedgerInterface interface {
	GetTokenByID(tokenID string) (common.SYN1900Token, error)
	UpdateToken(token common.SYN1900Token) error
	AddTransaction(transaction common.Transaction) error
	GetComplianceLogByToken(tokenID string) ([]common.ComplianceLog, error)
	AddComplianceLog(log common.ComplianceLog) error
}

// VerifyTokenCompliance ensures that a given SYN1900 token complies with regulatory or institutional standards.
func (compliance *ComplianceManager) VerifyTokenCompliance(tokenID string) error {
	// Fetch the token from the ledger
	token, err := compliance.ledger.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("token not found in the ledger")
	}

	// Check if the token has been revoked
	if token.RevocationStatus {
		return errors.New("token has been revoked and is non-compliant")
	}

	// Perform compliance checks (this can be expanded with institution-specific rules)
	if token.Issuer.IssuerType == "University" && len(token.CreditMetadata.DigitalSignature) == 0 {
		return errors.New("token is missing a valid digital signature and is non-compliant")
	}

	// Log compliance verification to the ledger
	complianceLog := common.ComplianceLog{
		TokenID:     token.TokenID,
		CheckedBy:   "Compliance Officer", // Placeholder for actual compliance officer ID or name
		CheckDate:   time.Now(),
		Status:      "Compliant",
		Description: "Token verified to comply with regulatory and institutional standards",
	}

	err = compliance.ledger.AddComplianceLog(complianceLog)
	if err != nil {
		return errors.New("failed to log compliance check to the ledger")
	}

	return nil
}

// RevokeNonCompliantToken revokes a token if it does not meet compliance standards.
func (compliance *ComplianceManager) RevokeNonCompliantToken(tokenID, reason string) error {
	// Fetch the token from the ledger
	token, err := compliance.ledger.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("token not found in the ledger")
	}

	// Check if the token is already revoked
	if token.RevocationStatus {
		return errors.New("token is already revoked")
	}

	// Mark the token as revoked and set the reason
	token.RevocationStatus = true
	token.RevokeCredit(reason)

	// Update the token in the ledger
	err = compliance.ledger.UpdateToken(token)
	if err != nil {
		return errors.New("failed to update token status in the ledger")
	}

	// Log the revocation in the ledger
	complianceLog := common.ComplianceLog{
		TokenID:     token.TokenID,
		CheckedBy:   "Compliance Officer", // Placeholder for actual compliance officer ID or name
		CheckDate:   time.Now(),
		Status:      "Revoked",
		Description: "Token revoked due to non-compliance: " + reason,
	}
	err = compliance.ledger.AddComplianceLog(complianceLog)
	if err != nil {
		return errors.New("failed to log revocation in the ledger")
	}

	return nil
}

// AuditTokenCompliance audits all tokens related to a specific institution or recipient to ensure full compliance.
func (compliance *ComplianceManager) AuditTokenCompliance(issuerID, recipientID string) ([]common.ComplianceLog, error) {
	var complianceLogs []common.ComplianceLog

	// Retrieve all tokens related to the issuer or recipient (simplified search logic)
	// This can be expanded with specific ledger search/filter logic.
	tokens, err := compliance.ledger.GetTokensByIssuerOrRecipient(issuerID, recipientID)
	if err != nil {
		return nil, errors.New("failed to retrieve tokens for audit")
	}

	// Iterate through tokens and check their compliance
	for _, token := range tokens {
		err := compliance.VerifyTokenCompliance(token.TokenID)
		if err == nil {
			// Fetch compliance log for this token
			logs, err := compliance.ledger.GetComplianceLogByToken(token.TokenID)
			if err != nil {
				return nil, errors.New("failed to fetch compliance logs")
			}
			complianceLogs = append(complianceLogs, logs...)
		}
	}

	return complianceLogs, nil
}

// GenerateComplianceReport generates a comprehensive compliance report for an institution or individual.
func (compliance *ComplianceManager) GenerateComplianceReport(issuerID, recipientID string) (string, error) {
	// Audit the compliance of all tokens for the issuer or recipient
	complianceLogs, err := compliance.AuditTokenCompliance(issuerID, recipientID)
	if err != nil {
		return "", err
	}

	// Generate the report
	report := "Compliance Report\n"
	report += "----------------------\n"
	report += "Issuer/Recipient ID: " + issuerID + "\n"
	report += "Audit Date: " + time.Now().Format(time.RFC3339) + "\n\n"
	report += "Tokens Audited: " + string(len(complianceLogs)) + "\n\n"
	for _, log := range complianceLogs {
		report += "Token ID: " + log.TokenID + "\n"
		report += "Status: " + log.Status + "\n"
		report += "Description: " + log.Description + "\n"
		report += "Checked By: " + log.CheckedBy + "\n"
		report += "Date: " + log.CheckDate.Format(time.RFC3339) + "\n\n"
	}

	return report, nil
}

// Helper function to encrypt compliance report data if needed.
func (compliance *ComplianceManager) EncryptComplianceReport(report string) ([]byte, error) {
	encryptedData, err := encryption.Encrypt([]byte(report))
	if err != nil {
		return nil, errors.New("failed to encrypt compliance report")
	}
	return encryptedData, nil
}
