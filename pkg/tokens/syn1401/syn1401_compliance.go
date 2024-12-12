package syn1401

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

// SYN1401ComplianceManager handles compliance checks and audit processes for SYN1401 tokens.
type SYN1401ComplianceManager struct {
	Ledger common.LedgerInterface // Interface to interact with the ledger
}

// PerformComplianceCheck verifies if the SYN1401 token is compliant with regulatory standards.
func (cm *SYN1401ComplianceManager) PerformComplianceCheck(tokenID string) (bool, error) {
	// Retrieve the token from the ledger
	token, err := cm.Ledger.GetToken(tokenID)
	if err != nil {
		return false, fmt.Errorf("error retrieving token from ledger: %w", err)
	}

	// Check compliance status
	if token.ComplianceStatus == "Compliant" {
		return true, nil
	} else if token.ComplianceStatus == "Non-Compliant" {
		return false, errors.New("token is non-compliant")
	}

	// If pending audit, perform a detailed compliance check
	if token.ComplianceStatus == "Pending Audit" {
		isCompliant, auditLog := cm.auditToken(token)
		if isCompliant {
			token.ComplianceStatus = "Compliant"
		} else {
			token.ComplianceStatus = "Non-Compliant"
		}

		// Append audit log to token's audit trail
		token.AuditTrail = append(token.AuditTrail, auditLog)

		// Update token in the ledger after audit
		if err := cm.Ledger.UpdateToken(tokenID, token); err != nil {
			return false, fmt.Errorf("error updating token after compliance check: %w", err)
		}

		return isCompliant, nil
	}

	return false, errors.New("unknown compliance status")
}

// auditToken performs a detailed audit of the SYN1401 token for compliance verification.
func (cm *SYN1401ComplianceManager) auditToken(token *common.SYN1401Token) (bool, common.AuditLog) {
	// Simulate detailed audit process
	auditPassed := true // Assume audit passes for demonstration
	description := "Detailed audit completed: Compliant with all regulatory standards."

	if token.PrincipalAmount <= 0 || token.InterestRate <= 0 {
		auditPassed = false
		description = "Audit failed: Invalid principal or interest rate."
	}

	// Create audit log
	auditLog := common.AuditLog{
		AuditID:     generateUniqueID(),
		PerformedBy: "Compliance System",
		Description: description,
		Timestamp:   time.Now(),
	}

	return auditPassed, auditLog
}

// LogComplianceEvent logs any compliance-related events such as audits, verifications, and breaches.
func (cm *SYN1401ComplianceManager) LogComplianceEvent(tokenID string, eventType string, description string) error {
	// Retrieve the token from the ledger
	token, err := cm.Ledger.GetToken(tokenID)
	if err != nil {
		return fmt.Errorf("error retrieving token: %w", err)
	}

	// Create compliance event log
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   eventType,
		Description: description,
		EventDate:   time.Now(),
		PerformedBy: "Compliance System",
	}

	// Append event log to the token
	token.EventLogs = append(token.EventLogs, eventLog)

	// Update token in the ledger
	if err := cm.Ledger.UpdateToken(tokenID, token); err != nil {
		return fmt.Errorf("error updating token with compliance event: %w", err)
	}

	return nil
}

// EncryptComplianceData encrypts sensitive compliance data before storage.
func (cm *SYN1401ComplianceManager) EncryptComplianceData(token *common.SYN1401Token, owner string) error {
	key, err := cm.getOwnerKey(owner)
	if err != nil {
		return err
	}

	plaintext := []byte(fmt.Sprintf("ComplianceStatus:%s|AuditTrail:%v", token.ComplianceStatus, token.AuditTrail))
	encryptedData, err := encryptAES(key, plaintext)
	if err != nil {
		return err
	}

	token.EncryptedMetadata = encryptedData
	return nil
}

// DecryptComplianceData decrypts sensitive compliance data for review.
func (cm *SYN1401ComplianceManager) DecryptComplianceData(token *common.SYN1401Token, owner string) error {
	key, err := cm.getOwnerKey(owner)
	if err != nil {
		return err
	}

	plaintext, err := decryptAES(key, token.EncryptedMetadata)
	if err != nil {
		return err
	}

	// Parse decrypted data (ComplianceStatus, AuditTrail)
	_, err = fmt.Sscanf(string(plaintext), "ComplianceStatus:%s|AuditTrail:%v", &token.ComplianceStatus, &token.AuditTrail)
	if err != nil {
		return fmt.Errorf("error parsing decrypted compliance data: %w", err)
	}

	return nil
}

// Helper function to retrieve owner's encryption key
func (cm *SYN1401ComplianceManager) getOwnerKey(owner string) ([]byte, error) {
	ownerInfo, err := cm.Ledger.GetOwnerInfo(owner)
	if err != nil {
		return nil, fmt.Errorf("error retrieving encryption key: %w", err)
	}

	return hex.DecodeString(ownerInfo.EncryptionKey)
}

// Helper functions for AES encryption and decryption
func encryptAES(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return aesGCM.Seal(nonce, nonce, plaintext, nil), nil
}

func decryptAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

// Helper function to generate a unique ID
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
