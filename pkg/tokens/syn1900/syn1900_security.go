package syn1900

import (
	"errors"
	"time"
)


// SecurityService provides security-related operations for the SYN1900 tokens.
type SecurityService struct {
	ledger LedgerInterface // Interface for interacting with the ledger
}

// LedgerInterface defines methods for interacting with the ledger.
type LedgerInterface interface {
	GetTokenByID(tokenID string) (common.SYN1900Token, error)
	UpdateToken(token common.SYN1900Token) error
}

// VerifyDigitalSignature verifies the digital signature on an education credit token.
func (ss *SecurityService) VerifyDigitalSignature(tokenID string) (bool, error) {
	// Fetch the token from the ledger
	token, err := ss.ledger.GetTokenByID(tokenID)
	if err != nil {
		return false, errors.New("token not found in the ledger")
	}

	// Verify the issuer's digital signature
	isValid := encryption.VerifyDigitalSignature(token.Issuer, token.Signature)
	if !isValid {
		return false, errors.New("invalid digital signature for token")
	}

	return true, nil
}

// EncryptTokenMetadata encrypts the metadata of the educational credit token.
func (ss *SecurityService) EncryptTokenMetadata(tokenID string, encryptionKey string) error {
	// Fetch the token from the ledger
	token, err := ss.ledger.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("token not found in the ledger")
	}

	// Encrypt the metadata
	encryptedMetadata, err := encryption.Encrypt([]byte(token.Metadata), encryptionKey)
	if err != nil {
		return errors.New("failed to encrypt token metadata")
	}

	// Update the token with encrypted metadata
	token.EncryptedMetadata = encryptedMetadata
	token.Metadata = "" // Clear plain text metadata

	// Update the ledger with the encrypted metadata
	err = ss.ledger.UpdateToken(token)
	if err != nil {
		return errors.New("failed to update token in the ledger")
	}

	return nil
}

// DecryptTokenMetadata decrypts the metadata of the educational credit token.
func (ss *SecurityService) DecryptTokenMetadata(tokenID string, decryptionKey string) (string, error) {
	// Fetch the token from the ledger
	token, err := ss.ledger.GetTokenByID(tokenID)
	if err != nil {
		return "", errors.New("token not found in the ledger")
	}

	// Decrypt the metadata
	decryptedMetadata, err := encryption.Decrypt(token.EncryptedMetadata, decryptionKey)
	if err != nil {
		return "", errors.New("failed to decrypt token metadata")
	}

	return string(decryptedMetadata), nil
}

// HashToken generates a cryptographic hash of the token's key details for integrity checks.
func (ss *SecurityService) HashToken(tokenID string) (string, error) {
	// Fetch the token from the ledger
	token, err := ss.ledger.GetTokenByID(tokenID)
	if err != nil {
		return "", errors.New("token not found in the ledger")
	}

	// Create a hash of the key token details
	hashInput := token.TokenID + token.CourseID + token.CourseName + token.RecipientID + token.Issuer + token.IssueDate.Format(time.RFC3339)
	hash := sha256.Sum256([]byte(hashInput))

	// Return the hash as a hexadecimal string
	return hex.EncodeToString(hash[:]), nil
}

// ValidateTokenIntegrity compares the token’s current hash with a previously stored hash to ensure integrity.
func (ss *SecurityService) ValidateTokenIntegrity(tokenID string, storedHash string) (bool, error) {
	// Generate the current hash for the token
	currentHash, err := ss.HashToken(tokenID)
	if err != nil {
		return false, err
	}

	// Compare the current hash with the stored hash
	if currentHash != storedHash {
		return false, errors.New("token integrity check failed: hashes do not match")
	}

	return true, nil
}

// RevokeToken revokes the token if it is determined that its security has been compromised.
func (ss *SecurityService) RevokeToken(tokenID, reason string) error {
	// Fetch the token from the ledger
	token, err := ss.ledger.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("token not found in the ledger")
	}

	// Set the token's revocation status
	token.Revoked = true
	token.RevocationReason = reason
	token.RevocationDate = time.Now()

	// Update the ledger with the revoked status
	err = ss.ledger.UpdateToken(token)
	if err != nil {
		return errors.New("failed to revoke token in the ledger")
	}

	return nil
}
