package syn2400

import (
	"errors"
	"time"

)

// SYN2400Security provides security-related functions for managing SYN2400 tokens
type SYN2400Security struct {
	Ledger   ledger.LedgerInterface        // Interface for interacting with the blockchain ledger
	Encrypt  encryption.EncryptionInterface // Interface for encryption
	Audit    audit.AuditInterface           // Interface for auditing activities
}

// NewSYN2400Security initializes a new instance of SYN2400Security
func NewSYN2400Security(ledger ledger.LedgerInterface, encrypt encryption.EncryptionInterface, audit audit.AuditInterface) *SYN2400Security {
	return &SYN2400Security{
		Ledger:  ledger,
		Encrypt: encrypt,
		Audit:   audit,
	}
}

// EncryptDataToken ensures the data associated with a SYN2400 token is encrypted before storage or transfer
func (security *SYN2400Security) EncryptDataToken(token common.SYN2400Token) ([]byte, error) {
	// Encrypt the token data using the encryption interface
	encryptedData, err := security.Encrypt.EncryptTokenData(token)
	if err != nil {
		return nil, errors.New("failed to encrypt data token: " + err.Error())
	}
	
	// Audit the encryption event
	security.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "EncryptDataToken",
		PerformedBy: token.Owner,
		Timestamp:   time.Now(),
		Details:     "Encrypted token " + token.TokenID,
	})

	return encryptedData, nil
}

// DecryptDataToken decrypts the encrypted data of a SYN2400 token for viewing or transfer
func (security *SYN2400Security) DecryptDataToken(encryptedData []byte, owner string) (common.SYN2400Token, error) {
	// Decrypt the token data using the encryption interface
	decryptedToken, err := security.Encrypt.DecryptTokenData(encryptedData)
	if err != nil {
		return common.SYN2400Token{}, errors.New("failed to decrypt data token: " + err.Error())
	}

	// Verify ownership before granting access to the decrypted data
	if decryptedToken.Owner != owner {
		return common.SYN2400Token{}, errors.New("ownership verification failed, unauthorized access attempt")
	}

	// Audit the decryption event
	security.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "DecryptDataToken",
		PerformedBy: owner,
		Timestamp:   time.Now(),
		Details:     "Decrypted token " + decryptedToken.TokenID,
	})

	return decryptedToken, nil
}

// ValidateOwnership verifies that the provided owner is indeed the owner of the token before transactions or updates
func (security *SYN2400Security) ValidateOwnership(tokenID string, owner string) (bool, error) {
	// Retrieve the encrypted token from the ledger
	encryptedToken, err := security.Ledger.GetToken(tokenID)
	if err != nil {
		return false, err
	}

	// Decrypt the token to verify ownership
	decryptedToken, err := security.DecryptDataToken(encryptedToken, owner)
	if err != nil {
		return false, err
	}

	if decryptedToken.Owner == owner {
		return true, nil
	}

	// Audit the ownership validation attempt
	security.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "ValidateOwnership",
		PerformedBy: owner,
		Timestamp:   time.Now(),
		Details:     "Ownership validation failed for token " + tokenID,
	})

	return false, errors.New("ownership validation failed")
}

// DetectFraudulentActivity checks for any suspicious activity in SYN2400 token transactions and flags them
func (security *SYN2400Security) DetectFraudulentActivity(tokenID string, activityType string) error {
	// Retrieve the encrypted token from the ledger
	encryptedToken, err := security.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	// Decrypt the token to access its metadata and transaction logs
	decryptedToken, err := security.DecryptDataToken(encryptedToken, decryptedToken.Owner)
	if err != nil {
		return err
	}

	// Apply fraud detection algorithms (for illustration purposes)
	if isSuspiciousActivity(decryptedToken, activityType) {
		// Flag the token as suspicious
		decryptedToken.Status = "Suspicious"

		// Encrypt the updated token
		encryptedToken, err := security.Encrypt.EncryptTokenData(decryptedToken)
		if err != nil {
			return err
		}

		// Store the flagged token back in the ledger
		if err := security.Ledger.UpdateToken(tokenID, encryptedToken); err != nil {
			return err
		}

		// Audit the fraud detection event
		security.Audit.LogAuditEvent(audit.AuditRecord{
			Action:      "DetectFraudulentActivity",
			PerformedBy: decryptedToken.Owner,
			Timestamp:   time.Now(),
			Details:     "Fraudulent activity detected for token " + tokenID,
		})

		return errors.New("fraudulent activity detected for token " + tokenID)
	}

	return nil
}

// SecureTransferToken ensures secure ownership transfer by validating both sender and receiver and enforcing encryption
func (security *SYN2400Security) SecureTransferToken(
	tokenID string,
	sender string,
	receiver string) (common.SYN2400Token, error) {

	// Validate sender's ownership of the token
	isOwner, err := security.ValidateOwnership(tokenID, sender)
	if err != nil || !isOwner {
		return common.SYN2400Token{}, errors.New("sender does not own the token")
	}

	// Retrieve the encrypted token from the ledger
	encryptedToken, err := security.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Decrypt the token for transfer
	decryptedToken, err := security.DecryptDataToken(encryptedToken, sender)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Transfer ownership
	decryptedToken.Owner = receiver
	decryptedToken.UpdateDate = time.Now()

	// Encrypt the updated token
	encryptedToken, err = security.Encrypt.EncryptTokenData(decryptedToken)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Store the updated token in the ledger
	if err := security.Ledger.UpdateToken(tokenID, encryptedToken); err != nil {
		return common.SYN2400Token{}, err
	}

	// Audit the transfer event
	security.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "SecureTransferToken",
		PerformedBy: sender,
		Timestamp:   time.Now(),
		Details:     "Transferred token " + tokenID + " to " + receiver,
	})

	return decryptedToken, nil
}

// isSuspiciousActivity is a dummy function to simulate fraud detection algorithms
func isSuspiciousActivity(token common.SYN2400Token, activityType string) bool {
	// In real-world implementations, this function would run sophisticated algorithms
	// to detect fraudulent activities based on transaction patterns, metadata anomalies, etc.
	return false
}
