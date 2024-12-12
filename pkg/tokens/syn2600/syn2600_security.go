package syn2600

import (

	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

)

// SYN2600TokenSecurity provides security functions for SYN2600 tokens.
type SYN2600TokenSecurity struct {
	TokenID        string
	AssetDetails   string
	Owner          string
	EncryptedToken string
	Signature      string // Encryption signature for verification
	Timestamp      time.Time
}

// ValidateTokenIntegrity validates the integrity of a token by comparing its hash signature.
func ValidateTokenIntegrity(tokenID string) (bool, error) {
	// Fetch token from the ledger
	encryptedToken, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return false, errors.New("failed to fetch token for integrity validation")
	}

	// Decrypt token data
	decryptedToken, err := encryption.DecryptTokenData(encryptedToken)
	if err != nil {
		return false, errors.New("failed to decrypt token data")
	}

	// Generate a new signature and compare with the stored one
	expectedSignature := generateTokenSignature(decryptedToken.TokenID, decryptedToken.AssetDetails, decryptedToken.Owner)
	if decryptedToken.Signature != expectedSignature {
		return false, errors.New("token integrity check failed: signatures do not match")
	}

	return true, nil
}

// EncryptTokenData encrypts token data using secure encryption techniques before storage.
func EncryptTokenData(token *SYN2600TokenSecurity) (string, error) {
	// Perform encryption using predefined encryption module
	encryptedToken, err := encryption.EncryptTokenData(token)
	if err != nil {
		return "", errors.New("encryption failed")
	}

	// Update encrypted token string
	token.EncryptedToken = encryptedToken
	return encryptedToken, nil
}

// DecryptTokenData decrypts token data for secure access.
func DecryptTokenData(encryptedToken string) (*SYN2600TokenSecurity, error) {
	// Decrypt using predefined encryption module
	decryptedToken, err := encryption.DecryptTokenData(encryptedToken)
	if err != nil {
		return nil, errors.New("decryption failed")
	}

	return decryptedToken, nil
}

// ValidateOwnership verifies that the specified user is the actual owner of the token.
func ValidateOwnership(tokenID string, owner string) (bool, error) {
	// Fetch encrypted token from the ledger
	encryptedToken, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return false, errors.New("failed to fetch token for ownership validation")
	}

	// Decrypt token data
	decryptedToken, err := encryption.DecryptTokenData(encryptedToken)
	if err != nil {
		return false, errors.New("failed to decrypt token data")
	}

	// Compare stored owner with the provided one
	if decryptedToken.Owner != owner {
		return false, errors.New("ownership validation failed: owner does not match")
	}

	return true, nil
}

// RecordSecurityEvent logs any security-related event in the ledger for audit purposes.
func RecordSecurityEvent(tokenID string, eventType string, details string, affectedParty string) (string, error) {
	// Create a timestamped security event record
	eventID := common.GenerateUniqueID()
	eventTimestamp := time.Now()

	// Build the event structure
	securityEvent := common.SecurityEvent{
		EventID:       eventID,
		TokenID:       tokenID,
		EventType:     eventType,
		Details:       details,
		AffectedParty: affectedParty,
		Timestamp:     eventTimestamp,
	}

	// Store the event in the ledger
	err := ledger.RecordEvent(securityEvent)
	if err != nil {
		return "", errors.New("failed to record security event in the ledger")
	}

	return eventID, nil
}

// RevokeToken revokes a SYN2600 token by deactivating it and recording the event.
func RevokeToken(tokenID string) (string, error) {
	// Fetch the token from the ledger
	encryptedToken, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return "", errors.New("failed to fetch token for revocation")
	}

	// Decrypt token data
	decryptedToken, err := encryption.DecryptTokenData(encryptedToken)
	if err != nil {
		return "", errors.New("failed to decrypt token data for revocation")
	}

	// Deactivate the token
	decryptedToken.Owner = ""
	decryptedToken.Signature = generateTokenSignature(decryptedToken.TokenID, decryptedToken.AssetDetails, "")
	decryptedToken.Timestamp = time.Now()

	// Re-encrypt the token
	reEncryptedToken, err := encryption.EncryptTokenData(decryptedToken)
	if err != nil {
		return "", errors.New("failed to re-encrypt token after revocation")
	}

	// Update the token in the ledger
	err = ledger.UpdateInvestorToken(reEncryptedToken)
	if err != nil {
		return "", errors.New("failed to update token after revocation in the ledger")
	}

	// Record the revocation event
	eventID, err := RecordSecurityEvent(tokenID, "REVOKE", "Token revoked", decryptedToken.Owner)
	if err != nil {
		return "", errors.New("failed to record revocation event")
	}

	return eventID, nil
}

// generateTokenSignature generates a signature for each token to ensure authenticity and track security.
func generateTokenSignature(tokenID string, assetDetails string, owner string) string {
	signatureInput := tokenID + assetDetails + owner
	hash := sha256.Sum256([]byte(signatureInput))
	return hex.EncodeToString(hash[:])
}
