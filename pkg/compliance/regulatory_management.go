package compliance

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SetComplianceThreshold adjusts the threshold for regulatory compliance checks.
func SetComplianceThreshold(ledger *ledger.Ledger, threshold int) error {
    if err := ledger.ComplianceLedger.SetThreshold(threshold); err != nil {
        return fmt.Errorf("failed to set compliance threshold: %v", err)
    }
    return nil
}

// GenerateAuditTrail creates an audit trail for all compliance-related actions.
func GenerateAuditTrail(ledger *ledger.Ledger, entityID string) (*ledger.AuditTrail, error) {
    trail, err := ledger.ComplianceLedger.CreateAuditTrail(entityID)
    if err != nil {
        return nil, fmt.Errorf("failed to generate audit trail: %v", err)
    }
    return trail, nil
}



// EncryptComplianceData securely encrypts sensitive compliance information.
func EncryptComplianceData(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create encryption cipher: %v", err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return ciphertext, nil
}

// DecryptComplianceData decrypts encrypted compliance information.
func DecryptComplianceData(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create decryption cipher: %v", err)
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	data := ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return data, nil
}

// InitiateDueDiligence begins the due diligence process for an entity.
func InitiateDueDiligence(ledger *ledger.Ledger, entityID string) error {
    if err := ledger.ComplianceLedger.BeginDueDiligence(entityID); err != nil {
        return fmt.Errorf("failed to initiate due diligence: %v", err)
    }
    return nil
}

// FinalizeDueDiligence completes the due diligence review.
func FinalizeDueDiligence(ledger *ledger.Ledger, entityID string) error {
    if err := ledger.ComplianceLedger.CompleteDueDiligence(entityID); err != nil {
        return fmt.Errorf("failed to finalize due diligence: %v", err)
    }
    return nil
}

// ReviewSanctionList checks if an entity appears on the sanction list.
func ReviewSanctionList(ledger *ledger.Ledger, entityID string) (bool, error) {
    listed, err := ledger.ComplianceLedger.CheckSanctionList(entityID)
    if err != nil || !listed {
        return false, fmt.Errorf("entity not on sanction list: %v", err)
    }
    return true, nil
}

// SetAccessControl configures access controls based on regulatory needs.
func SetAccessControl(ledger *ledger.Ledger, entityID string, controls ledger.AccessControls) error {
    if err := ledger.ComplianceLedger.ApplyAccessControls(entityID, controls); err != nil {
        return fmt.Errorf("failed to set access controls: %v", err)
    }
    return nil
}


// RestrictTransaction blocks a specific transaction based on compliance checks.
func RestrictTransaction(ledger *ledger.Ledger, transactionID string) error {
    if err := ledger.ComplianceLedger.BlockTransaction(transactionID); err != nil {
        return fmt.Errorf("failed to restrict transaction: %v", err)
    }
    return nil
}


// ReleaseTransaction lifts a restriction on a previously blocked transaction.
func ReleaseTransaction(ledger *ledger.Ledger, transactionID string) error {
    if err := ledger.ComplianceLedger.UnblockTransaction(transactionID); err != nil {
        return fmt.Errorf("failed to release transaction: %v", err)
    }
    return nil
}


// BlockUserAccess restricts a user’s access for regulatory violations.
func BlockUserAccess(ledger *ledger.Ledger, userID string) error {
    if err := ledger.ComplianceLedger.BlockUser(userID); err != nil {
        return fmt.Errorf("failed to block user access: %v", err)
    }
    return nil
}


// UnblockUserAccess restores a restricted user's access.
func UnblockUserAccess(ledger *ledger.Ledger, userID string) error {
    if err := ledger.ComplianceLedger.UnblockUser(userID); err != nil {
        return fmt.Errorf("failed to unblock user access: %v", err)
    }
    return nil
}


// ReviewComplianceHistory retrieves an entity's historical compliance data.
func ReviewComplianceHistory(ledger *ledger.Ledger, entityID string) ([]ledger.ComplianceHistory, error) {
    history, err := ledger.ComplianceLedger.FetchComplianceHistory(entityID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve compliance history: %v", err)
    }
    return history, nil
}


// IssueRegulatoryNotice sends a notice to an entity for regulatory issues.
func IssueRegulatoryNotice(ledger *ledger.Ledger, entityID string, notice string) error {
    if err := ledger.ComplianceLedger.SendRegulatoryNotice(entityID, notice); err != nil {
        return fmt.Errorf("failed to issue regulatory notice: %v", err)
    }
    return nil
}


// SubmitComplianceUpdate submits an updated compliance report.
func SubmitComplianceUpdate(ledger *ledger.Ledger, entityID string, report ledger.ComplianceReport) error {
    if err := ledger.ComplianceLedger.UpdateComplianceReport(entityID, report); err != nil {
        return fmt.Errorf("failed to submit compliance update: %v", err)
    }
    return nil
}


// RetrieveLegalDocument fetches a legal document required for compliance.
func RetrieveLegalDocument(ledger *ledger.Ledger, docID string) (*ledger.LegalDocument, error) {
    doc, err := ledger.ComplianceLedger.FetchLegalDocument(docID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve legal document: %v", err)
    }
    return &doc, nil
}


// EnforcePrivacyPolicy enforces privacy policies on an entity’s data.
func EnforcePrivacyPolicy(ledger *ledger.Ledger, entityID string) error {
    if err := ledger.ComplianceLedger.ApplyPrivacyPolicy(entityID); err != nil {
        return fmt.Errorf("failed to enforce privacy policy: %v", err)
    }
    return nil
}


// CreateRestrictionRule adds a new restriction rule to the compliance system.
func CreateRestrictionRule(ledger *ledger.Ledger, rule ledger.RestrictionRule) error {
    if err := ledger.ComplianceLedger.AddRestrictionRule(rule); err != nil {
        return fmt.Errorf("failed to create restriction rule: %v", err)
    }
    return nil
}


// RemoveRestrictionRule deletes an existing restriction rule.
func RemoveRestrictionRule(ledger *ledger.Ledger, ruleID string) error {
    if err := ledger.ComplianceLedger.DeleteRestrictionRule(ruleID); err != nil {
        return fmt.Errorf("failed to remove restriction rule: %v", err)
    }
    return nil
}


// SetDataRetentionPolicy configures the data retention policy for compliance data.
func SetDataRetentionPolicy(ledger *ledger.Ledger, policy ledger.RetentionPolicy) error {
    if err := ledger.ComplianceLedger.ApplyDataRetentionPolicy(policy); err != nil {
        return fmt.Errorf("failed to set data retention policy: %v", err)
    }
    return nil
}

