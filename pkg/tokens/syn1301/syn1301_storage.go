package syn1301

import (
	"errors"
	"time"

)

// SYN1301StorageManager handles the secure storage and retrieval of SYN1301 tokens.
type SYN1301StorageManager struct {
	Ledger            *ledger.Ledger                // Ledger system for storing tokens and assets
	EncryptionService *encryption.EncryptionService // Encryption service for data security
	SecurityService   *security.SecurityService     // Security service to manage access and verification
}

// StoreSYN1301Token securely stores a SYN1301 token with encryption.
func (sm *SYN1301StorageManager) StoreSYN1301Token(token SYN1301Token) error {
	// Encrypt token metadata before storage
	metadata := map[string]string{
		"asset_id":    token.AssetID,
		"description": token.Description,
		"location":    token.Location,
		"status":      token.Status,
		"owner":       token.Owner,
	}
	encryptedMetadata, err := sm.EncryptionService.Encrypt(metadata)
	if err != nil {
		return errors.New("encryption failed: " + err.Error())
	}
	token.EncryptedMetadata = encryptedMetadata

	// Perform security checks before storing the token
	err = sm.SecurityService.PerformSecurityChecks(token)
	if err != nil {
		return errors.New("security check failed: " + err.Error())
	}

	// Store token in the ledger
	err = sm.Ledger.StoreToken(token.TokenID, token)
	if err != nil {
		return errors.New("ledger storage failed: " + err.Error())
	}

	return nil
}

// RetrieveSYN1301Token retrieves a SYN1301 token from the ledger and decrypts it.
func (sm *SYN1301StorageManager) RetrieveSYN1301Token(tokenID string) (SYN1301Token, error) {
	// Retrieve token from the ledger
	token, err := sm.Ledger.GetToken(tokenID)
	if err != nil {
		return SYN1301Token{}, errors.New("token retrieval failed: " + err.Error())
	}

	// Decrypt metadata
	decryptedMetadata, err := sm.EncryptionService.Decrypt(token.EncryptedMetadata)
	if err != nil {
		return SYN1301Token{}, errors.New("decryption failed: " + err.Error())
	}

	// Populate the token with decrypted data
	token.Description = decryptedMetadata["description"]
	token.Location = decryptedMetadata["location"]
	token.Status = decryptedMetadata["status"]
	token.Owner = decryptedMetadata["owner"]

	return token, nil
}

// UpdateSYN1301TokenMetadata securely updates a token's metadata and stores the changes.
func (sm *SYN1301StorageManager) UpdateSYN1301TokenMetadata(tokenID string, updatedMetadata map[string]string) (SYN1301Token, error) {
	// Retrieve the current token
	token, err := sm.Ledger.GetToken(tokenID)
	if err != nil {
		return SYN1301Token{}, errors.New("failed to retrieve token for update: " + err.Error())
	}

	// Update token metadata
	if description, ok := updatedMetadata["description"]; ok {
		token.Description = description
	}
	if location, ok := updatedMetadata["location"]; ok {
		token.Location = location
	}
	if status, ok := updatedMetadata["status"]; ok {
		token.Status = status
	}
	if owner, ok := updatedMetadata["owner"]; ok {
		token.Owner = owner
	}

	// Encrypt the updated metadata
	encryptedMetadata, err := sm.EncryptionService.Encrypt(updatedMetadata)
	if err != nil {
		return SYN1301Token{}, errors.New("encryption failed: " + err.Error())
	}
	token.EncryptedMetadata = encryptedMetadata

	// Perform security checks after updating the token
	err = sm.SecurityService.PerformSecurityChecks(token)
	if err != nil {
		return SYN1301Token{}, errors.New("security check failed: " + err.Error())
	}

	// Update the token in the ledger
	err = sm.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return SYN1301Token{}, errors.New("failed to update token in ledger: " + err.Error())
	}

	return token, nil
}

// DeleteSYN1301Token securely removes a SYN1301 token from the ledger.
func (sm *SYN1301StorageManager) DeleteSYN1301Token(tokenID string) error {
	// Retrieve token from the ledger to ensure existence
	_, err := sm.Ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("token not found: " + err.Error())
	}

	// Perform security checks (e.g., authorization) before deletion
	err = sm.SecurityService.ValidateDeletionRights(tokenID)
	if err != nil {
		return errors.New("deletion rights validation failed: " + err.Error())
	}

	// Remove the token from the ledger
	err = sm.Ledger.DeleteToken(tokenID)
	if err != nil {
		return errors.New("failed to delete token from ledger: " + err.Error())
	}

	return nil
}

// ListSYN1301Tokens lists all tokens available in the storage.
func (sm *SYN1301StorageManager) ListSYN1301Tokens() ([]SYN1301Token, error) {
	// Retrieve all tokens from the ledger
	tokens, err := sm.Ledger.ListTokens()
	if err != nil {
		return nil, errors.New("failed to list tokens: " + err.Error())
	}

	// Decrypt metadata for each token
	for i := range tokens {
		decryptedMetadata, err := sm.EncryptionService.Decrypt(tokens[i].EncryptedMetadata)
		if err != nil {
			return nil, errors.New("failed to decrypt token metadata: " + err.Error())
		}

		// Populate the tokens with decrypted data
		tokens[i].Description = decryptedMetadata["description"]
		tokens[i].Location = decryptedMetadata["location"]
		tokens[i].Status = decryptedMetadata["status"]
		tokens[i].Owner = decryptedMetadata["owner"]
	}

	return tokens, nil
}
