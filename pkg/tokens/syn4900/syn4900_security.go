package syn4900

import (
	"errors"
	"sync"
	"time"
)

// Syn4900Security defines the security aspects of the Syn4900 token, including ownership, validation, and transfer.
type Syn4900Security struct {
	TokenID            string    `json:"token_id"`
	Owner              string    `json:"owner"`
	OwnershipVerified  bool      `json:"ownership_verified"`
	LastVerified       time.Time `json:"last_verified"`
	EncryptedMetadata  string    `json:"encrypted_metadata"`
}

// SecurityManager handles the security, verification, and ownership validation of Syn4900 tokens.
type SecurityManager struct {
	mutex             sync.Mutex
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	tokenSecurity     map[string]*Syn4900Security
}

// NewSecurityManager creates a new instance of the SecurityManager for managing Syn4900 token security.
func NewSecurityManager(ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *SecurityManager {
	return &SecurityManager{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
		tokenSecurity:     make(map[string]*Syn4900Security),
	}
}

// VerifyOwnership verifies the ownership of a given Syn4900 token and records it in the ledger.
func (sm *SecurityManager) VerifyOwnership(tokenID, ownerID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the token's security details.
	security, exists := sm.tokenSecurity[tokenID]
	if !exists {
		return errors.New("token not found for ownership verification")
	}

	// Verify the owner matches the provided ownerID.
	if security.Owner != ownerID {
		return errors.New("ownership verification failed: owner mismatch")
	}

	// Mark the ownership as verified.
	security.OwnershipVerified = true
	security.LastVerified = time.Now()

	// Log the ownership verification in the ledger.
	if err := sm.ledgerService.LogEvent("OwnershipVerified", time.Now(), tokenID); err != nil {
		return err
	}

	// Validate the ownership verification using Synnergy Consensus.
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	// Update the token's security in the internal map.
	sm.tokenSecurity[tokenID] = security

	return nil
}

// TransferOwnership securely transfers ownership of the Syn4900 token to a new owner.
func (sm *SecurityManager) TransferOwnership(tokenID, newOwnerID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the token's security details.
	security, exists := sm.tokenSecurity[tokenID]
	if !exists {
		return errors.New("token not found for ownership transfer")
	}

	// Ensure ownership is verified before transfer.
	if !security.OwnershipVerified {
		return errors.New("cannot transfer ownership: token ownership not verified")
	}

	// Update the owner of the token.
	security.Owner = newOwnerID
	security.OwnershipVerified = false // Reset ownership verification for new owner.

	// Encrypt the updated metadata.
	encryptedMetadata, err := sm.encryptionService.EncryptData(security)
	if err != nil {
		return err
	}
	security.EncryptedMetadata = encryptedMetadata.(string)

	// Log the ownership transfer in the ledger.
	if err := sm.ledgerService.LogEvent("OwnershipTransferred", time.Now(), tokenID); err != nil {
		return err
	}

	// Validate the ownership transfer using Synnergy Consensus.
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	// Update the token's security in the internal map.
	sm.tokenSecurity[tokenID] = security

	return nil
}

// EncryptTokenMetadata encrypts the token's metadata to ensure data security.
func (sm *SecurityManager) EncryptTokenMetadata(tokenID string, metadata interface{}) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the token's security details.
	security, exists := sm.tokenSecurity[tokenID]
	if !exists {
		return errors.New("token not found for metadata encryption")
	}

	// Encrypt the metadata.
	encryptedMetadata, err := sm.encryptionService.EncryptData(metadata)
	if err != nil {
		return err
	}

	// Store the encrypted metadata in the security structure.
	security.EncryptedMetadata = encryptedMetadata.(string)

	// Log the encryption event in the ledger.
	if err := sm.ledgerService.LogEvent("MetadataEncrypted", time.Now(), tokenID); err != nil {
		return err
	}

	// Update the token's security in the internal map.
	sm.tokenSecurity[tokenID] = security

	return nil
}

// DecryptTokenMetadata decrypts the token's metadata for secure access.
func (sm *SecurityManager) DecryptTokenMetadata(tokenID string) (interface{}, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the token's security details.
	security, exists := sm.tokenSecurity[tokenID]
	if !exists {
		return nil, errors.New("token not found for metadata decryption")
	}

	// Decrypt the metadata.
	decryptedMetadata, err := sm.encryptionService.DecryptData(security.EncryptedMetadata)
	if err != nil {
		return nil, err
	}

	return decryptedMetadata, nil
}

// ValidateTokenSecurity performs a full security check on the token's metadata and owner.
func (sm *SecurityManager) ValidateTokenSecurity(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the token's security details.
	security, exists := sm.tokenSecurity[tokenID]
	if !exists {
		return errors.New("token not found for security validation")
	}

	// Perform encryption and ownership verification validation.
	if security.EncryptedMetadata == "" || !security.OwnershipVerified {
		return errors.New("token security validation failed")
	}

	// Log the security validation in the ledger.
	if err := sm.ledgerService.LogEvent("TokenSecurityValidated", time.Now(), tokenID); err != nil {
		return err
	}

	// Validate the security validation using Synnergy Consensus.
	return sm.consensusService.ValidateSubBlock(tokenID)
}

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// generateUniqueTokenID generates a cryptographically secure, time-based unique identifier for a new token.
func generateUniqueTokenID() string {
	// Create a buffer for random bytes
	randomBytes := make([]byte, 16)

	// Read cryptographically secure random bytes
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Fallback to a timestamp-based ID if random generation fails
		return "token_" + time.Now().Format("20060102150405")
	}

	// Convert the random bytes to a hex string
	randomString := hex.EncodeToString(randomBytes)

	// Append current timestamp for extra uniqueness (and readability)
	timestamp := time.Now().UnixNano()

	// Generate final unique token ID
	uniqueTokenID := randomString + "_" + string(timestamp)

	return uniqueTokenID
}

