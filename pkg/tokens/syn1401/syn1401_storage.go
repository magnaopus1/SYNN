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

// SYN1401Storage manages token storage, ledger integration, encryption, and data retrieval for SYN1401 tokens.
type SYN1401Storage struct {
	Ledger common.LedgerInterface // Interface to interact with the ledger
}

// StoreToken securely stores the SYN1401 token in the blockchain ledger after encrypting sensitive data.
func (s *SYN1401Storage) StoreToken(owner string, token *common.SYN1401Token) error {
	// Encrypt sensitive token data before storage
	if err := s.encryptTokenData(owner, token); err != nil {
		return fmt.Errorf("error encrypting token data: %w", err)
	}

	// Store the encrypted token in the ledger
	if err := s.Ledger.StoreToken(token); err != nil {
		return fmt.Errorf("error storing token in ledger: %w", err)
	}

	// Log the storage event
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Token Storage",
		Description: fmt.Sprintf("Token %s securely stored for owner %s", token.TokenID, owner),
		EventDate:   time.Now(),
		PerformedBy: "System",
	}
	token.EventLogs = append(token.EventLogs, eventLog)

	return nil
}

// RetrieveToken retrieves and decrypts a SYN1401 token from the ledger.
func (s *SYN1401Storage) RetrieveToken(owner string, tokenID string) (*common.SYN1401Token, error) {
	// Fetch token from the ledger
	token, err := s.Ledger.GetToken(tokenID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving token from ledger: %w", err)
	}

	// Decrypt sensitive token data
	if err := s.decryptTokenData(owner, token); err != nil {
		return nil, fmt.Errorf("error decrypting token data: %w", err)
	}

	// Log the retrieval event
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Token Retrieval",
		Description: fmt.Sprintf("Token %s retrieved by owner %s", token.TokenID, owner),
		EventDate:   time.Now(),
		PerformedBy: "System",
	}
	token.EventLogs = append(token.EventLogs, eventLog)

	return token, nil
}

// UpdateToken updates the details of a SYN1401 token in the ledger after validating and encrypting sensitive fields.
func (s *SYN1401Storage) UpdateToken(owner string, token *common.SYN1401Token) error {
	// Encrypt sensitive token data before updating
	if err := s.encryptTokenData(owner, token); err != nil {
		return fmt.Errorf("error encrypting token data for update: %w", err)
	}

	// Update the token in the ledger
	if err := s.Ledger.UpdateToken(token.TokenID, token); err != nil {
		return fmt.Errorf("error updating token in ledger: %w", err)
	}

	// Log the update event
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Token Update",
		Description: fmt.Sprintf("Token %s updated for owner %s", token.TokenID, owner),
		EventDate:   time.Now(),
		PerformedBy: "System",
	}
	token.EventLogs = append(token.EventLogs, eventLog)

	return nil
}

// DeleteToken removes a token from the ledger and logs the deletion.
func (s *SYN1401Storage) DeleteToken(owner string, tokenID string) error {
	// Delete the token from the ledger
	if err := s.Ledger.DeleteToken(tokenID); err != nil {
		return fmt.Errorf("error deleting token from ledger: %w", err)
	}

	// Log the deletion event
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Token Deletion",
		Description: fmt.Sprintf("Token %s deleted for owner %s", tokenID, owner),
		EventDate:   time.Now(),
		PerformedBy: "System",
	}
	// Optionally, store the event log in the audit trail (if deletion events are to be recorded)

	return nil
}

// encryptTokenData encrypts sensitive fields of the SYN1401 token.
func (s *SYN1401Storage) encryptTokenData(owner string, token *common.SYN1401Token) error {
	// Get encryption key for the owner
	key, err := s.getOwnerKey(owner)
	if err != nil {
		return err
	}

	// Prepare sensitive data for encryption (TokenID, Owner, PrincipalAmount)
	plaintext := []byte(fmt.Sprintf("%s|%f|%s", token.TokenID, token.PrincipalAmount, token.Owner))

	// Encrypt the data using AES
	encryptedData, err := s.encryptAES(key, plaintext)
	if err != nil {
		return err
	}

	token.EncryptedMetadata = encryptedData
	return nil
}

// decryptTokenData decrypts sensitive fields of the SYN1401 token.
func (s *SYN1401Storage) decryptTokenData(owner string, token *common.SYN1401Token) error {
	// Get decryption key for the owner
	key, err := s.getOwnerKey(owner)
	if err != nil {
		return err
	}

	// Decrypt the token's encrypted metadata
	plaintext, err := s.decryptAES(key, token.EncryptedMetadata)
	if err != nil {
		return err
	}

	// Parse decrypted data (TokenID, PrincipalAmount, Owner)
	_, err = fmt.Sscanf(string(plaintext), "%s|%f|%s", &token.TokenID, &token.PrincipalAmount, &token.Owner)
	if err != nil {
		return fmt.Errorf("error parsing decrypted token data: %w", err)
	}

	return nil
}

// encryptAES performs AES encryption on plaintext using the provided key.
func (s *SYN1401Storage) encryptAES(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return aesGCM.Seal(nonce, nonce, plaintext, nil), nil
}

// decryptAES performs AES decryption on ciphertext using the provided key.
func (s *SYN1401Storage) decryptAES(key, ciphertext []byte) ([]byte, error) {
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

// getOwnerKey retrieves the encryption/decryption key for a given owner.
func (s *SYN1401Storage) getOwnerKey(owner string) ([]byte, error) {
	// Fetch the owner's key (for example, from a secure storage)
	// In this example, we assume the key is stored in the ledger for simplicity.
	ownerInfo, err := s.Ledger.GetOwnerInfo(owner)
	if err != nil {
		return nil, fmt.Errorf("error fetching encryption key for owner %s: %w", owner, err)
	}

	return hex.DecodeString(ownerInfo.EncryptionKey)
}

// generateUniqueID generates a unique ID for events or transactions.
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
