package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// SYN721Metadata defines the structure of the metadata for an SYN721 token (NFT).
type SYN721Metadata struct {
	TokenID       string  // Unique ID for the NFT token
	TokenURI      string  // URI pointing to the metadata (e.g., IPFS link)
	TokenOwner    string  // Address of the token owner
	TokenAttributes    string  // Address of the token owner
	EncryptedData string  // Encrypted metadata for security
	ValidatorID   string  // ID of the validator that validated this metadata
	MetadataHash  string  // Hash of the metadata for integrity verification
    Name            string                 `json:"name"`
    Description     string                 `json:"description"`
    DocumentHash    string                 `json:"document_hash,omitempty"` // Used for non-fungible mode
    Properties      map[string]interface{} `json:"properties,omitempty"` // Additional metadata for non-fungible tokens
}

// SYN721Manager handles the creation, encryption, and management of SYN721 metadata.
type SYN721Manager struct {
	mutex        sync.Mutex                 // For thread-safe operations
	Ledger       *ledger.Ledger             // Reference to the ledger for storing token metadata
	Consensus    *synnergy_consensus.Engine // Synnergy Consensus engine for validating metadata
	Encryption   *encryption.Encryption     // Encryption service for securing metadata
	Metadata     map[string]*SYN721Metadata // Stores metadata for each token by TokenID
}

// NewSYN721Manager initializes a new SYN721 metadata manager.
func NewSYN721Manager(ledgerInstance *ledger.Ledger, consensus *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *SYN721Manager {
	return &SYN721Manager{
		Ledger:     ledgerInstance,
		Consensus:  consensus,
		Encryption: encryptionService,
		Metadata:   make(map[string]*SYN721Metadata),
	}
}

// CreateMetadata generates and stores the metadata for an SYN721 token (NFT).
func (sm *SYN721Manager) CreateMetadata(tokenID, tokenURI, owner string) (*SYN721Metadata, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if the token metadata already exists
	if _, exists := sm.Metadata[tokenID]; exists {
		return nil, errors.New("token metadata already exists")
	}

	// Create the metadata structure
	metadata := &SYN721Metadata{
		TokenID:    tokenID,
		TokenURI:   tokenURI,
		TokenOwner: owner,
	}

	// Encrypt the metadata for security
	encryptedData, err := sm.Encryption.EncryptData(fmt.Sprintf("%v", metadata), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting token metadata: %v", err)
	}
	metadata.EncryptedData = encryptedData

	// Validate token creation using Synnergy Consensus
	if valid, err := sm.Consensus.ValidateTokenCreation(metadata); !valid || err != nil {
		return nil, fmt.Errorf("error validating token creation: %v", err)
	}

	// Store metadata hash for verification
	metadata.MetadataHash = common.GenerateHash(fmt.Sprintf("%v", metadata))

	// Log the metadata in the ledger
	if err := sm.Ledger.AddTokenMetadata(tokenID, metadata); err != nil {
		return nil, fmt.Errorf("error storing token metadata in ledger: %v", err)
	}

	// Store metadata in memory
	sm.Metadata[tokenID] = metadata

	fmt.Printf("SYN721 token %s created with metadata at URI: %s\n", tokenID, tokenURI)
	return metadata, nil
}

// UpdateMetadata allows the owner to update the metadata of the SYN721 token.
func (sm *SYN721Manager) UpdateMetadata(owner, tokenID, newURI string) (*SYN721Metadata, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the metadata for the token
	metadata, exists := sm.Metadata[tokenID]
	if !exists {
		return nil, errors.New("metadata for the token does not exist")
	}

	// Verify ownership
	if metadata.TokenOwner != owner {
		return nil, errors.New("only the token owner can update metadata")
	}

	// Update the metadata URI
	metadata.TokenURI = newURI

	// Encrypt updated metadata
	encryptedData, err := sm.Encryption.EncryptData(fmt.Sprintf("%v", metadata), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting updated metadata: %v", err)
	}
	metadata.EncryptedData = encryptedData

	// Update metadata hash
	metadata.MetadataHash = common.GenerateHash(fmt.Sprintf("%v", metadata))

	// Log the updated metadata in the ledger
	if err := sm.Ledger.UpdateTokenMetadata(tokenID, metadata); err != nil {
		return nil, fmt.Errorf("error updating token metadata in ledger: %v", err)
	}

	fmt.Printf("Metadata for SYN721 token %s updated. New URI: %s\n", tokenID, newURI)
	return metadata, nil
}

// TransferOwnership transfers the ownership of an SYN721 token to a new address.
func (sm *SYN721Manager) TransferOwnership(currentOwner, newOwner, tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the metadata for the token
	metadata, exists := sm.Metadata[tokenID]
	if !exists {
		return errors.New("metadata for the token does not exist")
	}

	// Verify ownership
	if metadata.TokenOwner != currentOwner {
		return errors.New("only the token owner can transfer ownership")
	}

	// Transfer ownership
	metadata.TokenOwner = newOwner

	// Encrypt updated metadata
	encryptedData, err := sm.Encryption.EncryptData(fmt.Sprintf("%v", metadata), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting updated metadata: %v", err)
	}
	metadata.EncryptedData = encryptedData

	// Log the ownership transfer in the ledger
	if err := sm.Ledger.UpdateTokenMetadata(tokenID, metadata); err != nil {
		return fmt.Errorf("error updating token ownership in ledger: %v", err)
	}

	fmt.Printf("Ownership of SYN721 token %s transferred from %s to %s.\n", tokenID, currentOwner, newOwner)
	return nil
}

// GetMetadata retrieves and decrypts the metadata for an SYN721 token.
func (sm *SYN721Manager) GetMetadata(tokenID string) (*SYN721Metadata, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the metadata for the token
	metadata, exists := sm.Metadata[tokenID]
	if !exists {
		return nil, errors.New("metadata for the token does not exist")
	}

	// Decrypt the metadata
	decryptedData, err := sm.Encryption.DecryptData(metadata.EncryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error decrypting token metadata: %v", err)
	}

	fmt.Printf("Retrieved metadata for SYN721 token %s: %v\n", tokenID, decryptedData)
	return metadata, nil
}

// VerifyMetadataHash verifies the integrity of the metadata by checking the stored hash.
func (sm *SYN721Manager) VerifyMetadataHash(tokenID string) (bool, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the metadata for the token
	metadata, exists := sm.Metadata[tokenID]
	if !exists {
		return false, errors.New("metadata for the token does not exist")
	}

	// Generate a new hash and compare with the stored hash
	calculatedHash := common.GenerateHash(fmt.Sprintf("%v", metadata))
	if calculatedHash != metadata.MetadataHash {
		return false, errors.New("metadata integrity check failed")
	}

	fmt.Printf("Metadata hash for SYN721 token %s verified successfully.\n", tokenID)
	return true, nil
}
