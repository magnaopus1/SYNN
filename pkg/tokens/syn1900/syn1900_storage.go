package syn1900

import (
	"errors"
	"time"
)


// TokenStorageService provides storage-related operations for SYN1900 tokens.
type TokenStorageService struct {
	tokenStore TokenStoreInterface // Interface for interacting with the storage backend
	mutex      sync.Mutex          // To handle concurrent access to the storage
}

// TokenStoreInterface defines methods for interacting with the storage layer.
type TokenStoreInterface interface {
	SaveToken(token common.SYN1900Token) error
	GetTokenByID(tokenID string) (common.SYN1900Token, error)
	UpdateToken(token common.SYN1900Token) error
	DeleteToken(tokenID string) error
	ListAllTokens() ([]common.SYN1900Token, error)
}

// SaveToken stores a new SYN1900 token into the storage backend.
func (s *TokenStorageService) SaveToken(token common.SYN1900Token) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Encrypt sensitive metadata before saving
	encryptedMetadata, err := encryption.Encrypt([]byte(token.Metadata), "your-encryption-key")
	if err != nil {
		return errors.New("failed to encrypt token metadata")
	}
	token.EncryptedMetadata = encryptedMetadata
	token.Metadata = "" // Clear the plain metadata for security

	// Save the token to storage
	err = s.tokenStore.SaveToken(token)
	if err != nil {
		return errors.New("failed to save token to storage")
	}

	return nil
}

// GetToken retrieves a SYN1900 token by its ID from the storage backend.
func (s *TokenStorageService) GetToken(tokenID string, decryptionKey string) (common.SYN1900Token, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Retrieve the token from the store
	token, err := s.tokenStore.GetTokenByID(tokenID)
	if err != nil {
		return common.SYN1900Token{}, errors.New("token not found in storage")
	}

	// Decrypt the metadata
	decryptedMetadata, err := encryption.Decrypt(token.EncryptedMetadata, decryptionKey)
	if err != nil {
		return common.SYN1900Token{}, errors.New("failed to decrypt token metadata")
	}
	token.Metadata = string(decryptedMetadata)

	return token, nil
}

// UpdateToken updates the details of an existing SYN1900 token in the storage backend.
func (s *TokenStorageService) UpdateToken(token common.SYN1900Token) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Re-encrypt metadata before updating
	encryptedMetadata, err := encryption.Encrypt([]byte(token.Metadata), "your-encryption-key")
	if err != nil {
		return errors.New("failed to encrypt token metadata")
	}
	token.EncryptedMetadata = encryptedMetadata
	token.Metadata = "" // Clear the plain metadata for security

	// Update the token in storage
	err = s.tokenStore.UpdateToken(token)
	if err != nil {
		return errors.New("failed to update token in storage")
	}

	return nil
}

// DeleteToken removes a SYN1900 token from the storage backend.
func (s *TokenStorageService) DeleteToken(tokenID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Delete the token from the store
	err := s.tokenStore.DeleteToken(tokenID)
	if err != nil {
		return errors.New("failed to delete token from storage")
	}

	return nil
}

// ListAllTokens retrieves a list of all SYN1900 tokens stored in the backend.
func (s *TokenStorageService) ListAllTokens() ([]common.SYN1900Token, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Retrieve all tokens from the storage
	tokens, err := s.tokenStore.ListAllTokens()
	if err != nil {
		return nil, errors.New("failed to list tokens from storage")
	}

	return tokens, nil
}

// RevokeToken updates the token's status as revoked and stores the reason and timestamp.
func (s *TokenStorageService) RevokeToken(tokenID, reason string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Retrieve the token from storage
	token, err := s.tokenStore.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("token not found in storage")
	}

	// Mark the token as revoked
	token.Revoked = true
	token.RevocationReason = reason
	token.RevocationDate = time.Now()

	// Update the token in storage
	err = s.tokenStore.UpdateToken(token)
	if err != nil {
		return errors.New("failed to revoke token in storage")
	}

	return nil
}
