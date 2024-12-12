package syn2400

import (
	"errors"
	"time"

)

// SYN2400Storage provides secure storage-related functions for managing SYN2400 tokens
type SYN2400Storage struct {
	Ledger      ledger.LedgerInterface          // Interface for interacting with the blockchain ledger
	Encrypt     encryption.EncryptionInterface   // Interface for encryption
	Compress    compression.CompressionInterface // Interface for compression
	Audit       audit.AuditInterface             // Interface for auditing activities
}

// NewSYN2400Storage initializes a new instance of SYN2400Storage
func NewSYN2400Storage(ledger ledger.LedgerInterface, encrypt encryption.EncryptionInterface, compress compression.CompressionInterface, audit audit.AuditInterface) *SYN2400Storage {
	return &SYN2400Storage{
		Ledger:   ledger,
		Encrypt:  encrypt,
		Compress: compress,
		Audit:    audit,
	}
}

// StoreDataToken securely stores a SYN2400 token in the ledger after encryption and optional compression
func (storage *SYN2400Storage) StoreDataToken(token common.SYN2400Token, compress bool) error {
	// Encrypt the token data using the encryption interface
	encryptedToken, err := storage.Encrypt.EncryptTokenData(token)
	if err != nil {
		return errors.New("failed to encrypt data token: " + err.Error())
	}

	// Optionally compress the encrypted token data before storing
	if compress {
		encryptedToken, err = storage.Compress.CompressData(encryptedToken)
		if err != nil {
			return errors.New("failed to compress encrypted data token: " + err.Error())
		}
	}

	// Store the encrypted (and possibly compressed) token in the ledger
	err = storage.Ledger.StoreToken(token.TokenID, encryptedToken)
	if err != nil {
		return errors.New("failed to store token in ledger: " + err.Error())
	}

	// Audit the token storage event
	storage.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "StoreDataToken",
		PerformedBy: token.Owner,
		Timestamp:   time.Now(),
		Details:     "Stored encrypted token " + token.TokenID,
	})

	return nil
}

// RetrieveDataToken retrieves and decrypts a SYN2400 token from the ledger
func (storage *SYN2400Storage) RetrieveDataToken(tokenID string, owner string, decompress bool) (common.SYN2400Token, error) {
	// Retrieve the encrypted (and possibly compressed) token from the ledger
	encryptedToken, err := storage.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, errors.New("failed to retrieve token from ledger: " + err.Error())
	}

	// Optionally decompress the encrypted token data
	if decompress {
		encryptedToken, err = storage.Compress.DecompressData(encryptedToken)
		if err != nil {
			return common.SYN2400Token{}, errors.New("failed to decompress encrypted data token: " + err.Error())
		}
	}

	// Decrypt the token data using the encryption interface
	decryptedToken, err := storage.Encrypt.DecryptTokenData(encryptedToken)
	if err != nil {
		return common.SYN2400Token{}, errors.New("failed to decrypt data token: " + err.Error())
	}

	// Verify ownership before granting access to the decrypted token
	if decryptedToken.Owner != owner {
		return common.SYN2400Token{}, errors.New("ownership verification failed, unauthorized access attempt")
	}

	// Audit the token retrieval event
	storage.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "RetrieveDataToken",
		PerformedBy: owner,
		Timestamp:   time.Now(),
		Details:     "Retrieved token " + tokenID,
	})

	return decryptedToken, nil
}

// UpdateDataToken updates the SYN2400 token stored in the ledger after encryption and optional compression
func (storage *SYN2400Storage) UpdateDataToken(token common.SYN2400Token, compress bool) error {
	// Encrypt the token data using the encryption interface
	encryptedToken, err := storage.Encrypt.EncryptTokenData(token)
	if err != nil {
		return errors.New("failed to encrypt data token: " + err.Error())
	}

	// Optionally compress the encrypted token data before storing
	if compress {
		encryptedToken, err = storage.Compress.CompressData(encryptedToken)
		if err != nil {
			return errors.New("failed to compress encrypted data token: " + err.Error())
		}
	}

	// Update the encrypted (and possibly compressed) token in the ledger
	err = storage.Ledger.UpdateToken(token.TokenID, encryptedToken)
	if err != nil {
		return errors.New("failed to update token in ledger: " + err.Error())
	}

	// Audit the token update event
	storage.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "UpdateDataToken",
		PerformedBy: token.Owner,
		Timestamp:   time.Now(),
		Details:     "Updated encrypted token " + token.TokenID,
	})

	return nil
}

// VerifyDataIntegrity verifies the integrity of the data token stored in the ledger
func (storage *SYN2400Storage) VerifyDataIntegrity(tokenID string) (bool, error) {
	// Retrieve the encrypted token from the ledger
	encryptedToken, err := storage.Ledger.GetToken(tokenID)
	if err != nil {
		return false, errors.New("failed to retrieve token from ledger: " + err.Error())
	}

	// Verify the integrity of the encrypted token using the encryption interface
	isValid, err := storage.Encrypt.VerifyTokenIntegrity(encryptedToken)
	if err != nil {
		return false, errors.New("failed to verify data token integrity: " + err.Error())
	}

	// Audit the integrity verification event
	storage.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "VerifyDataIntegrity",
		PerformedBy: "System",
		Timestamp:   time.Now(),
		Details:     "Verified integrity for token " + tokenID + " - Valid: " + boolToString(isValid),
	})

	return isValid, nil
}

// DeleteDataToken securely deletes a SYN2400 token from the ledger
func (storage *SYN2400Storage) DeleteDataToken(tokenID string, owner string) error {
	// Verify that the owner is authorized to delete the token
	isOwner, err := storage.ValidateOwnership(tokenID, owner)
	if err != nil || !isOwner {
		return errors.New("unauthorized deletion attempt, ownership validation failed")
	}

	// Delete the token from the ledger
	err = storage.Ledger.DeleteToken(tokenID)
	if err != nil {
		return errors.New("failed to delete token from ledger: " + err.Error())
	}

	// Audit the token deletion event
	storage.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "DeleteDataToken",
		PerformedBy: owner,
		Timestamp:   time.Now(),
		Details:     "Deleted token " + tokenID,
	})

	return nil
}

// ValidateOwnership verifies that the provided owner is indeed the owner of the token before allowing deletion or updates
func (storage *SYN2400Storage) ValidateOwnership(tokenID string, owner string) (bool, error) {
	// Retrieve the encrypted token from the ledger
	encryptedToken, err := storage.Ledger.GetToken(tokenID)
	if err != nil {
		return false, err
	}

	// Decrypt the token to verify ownership
	decryptedToken, err := storage.Encrypt.DecryptTokenData(encryptedToken)
	if err != nil {
		return false, err
	}

	if decryptedToken.Owner == owner {
		return true, nil
	}

	return false, errors.New("ownership validation failed")
}

// boolToString converts a boolean to string for logging purposes
func boolToString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}
