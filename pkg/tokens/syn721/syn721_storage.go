package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// SYN721Storage defines the structure for storing ownership and metadata of SYN721 tokens (NFTs).
type SYN721Storage struct {
	mutex          sync.Mutex
	TokenOwners    map[string]string            // Maps tokenID -> owner address
	TokenMetadata  map[string]*SYN721Metadata   // Maps tokenID -> metadata
	Ledger         *ledger.Ledger               // Ledger for blockchain transactions
	Consensus      *synnergy_consensus.Engine   // Synnergy Consensus engine for validating changes
	Encryption     *encryption.Encryption       // Encryption service for securing token storage
}

// NewSYN721Storage initializes a new storage for SYN721 tokens.
func NewSYN721Storage(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *SYN721Storage {
	return &SYN721Storage{
		TokenOwners:   make(map[string]string),
		TokenMetadata: make(map[string]*SYN721Metadata),
		Ledger:        ledgerInstance,
		Consensus:     consensusEngine,
		Encryption:    encryptionService,
	}
}

// MintToken mints a new SYN721 token, assigns it to an owner, and stores its metadata.
func (storage *SYN721Storage) MintToken(tokenID, owner, tokenURI string) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	// Check if token already exists
	if _, exists := storage.TokenOwners[tokenID]; exists {
		return errors.New("token already exists")
	}

	// Create new token metadata
	metadata := &SYN721Metadata{
		TokenID:    tokenID,
		TokenURI:   tokenURI,
		TokenOwner: owner,
	}

	// Encrypt the token metadata
	encryptedData, err := storage.Encryption.EncryptData(fmt.Sprintf("%v", metadata), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting token metadata: %v", err)
	}
	metadata.EncryptedData = encryptedData

	// Validate token creation using Synnergy Consensus
	if valid, err := storage.Consensus.ValidateTokenCreation(metadata); !valid || err != nil {
		return fmt.Errorf("error validating token creation: %v", err)
	}

	// Store token owner and metadata
	storage.TokenOwners[tokenID] = owner
	storage.TokenMetadata[tokenID] = metadata

	// Log the token minting in the ledger
	if err := storage.Ledger.RecordTokenCreation(tokenID, owner, tokenURI); err != nil {
		return fmt.Errorf("error recording token minting in ledger: %v", err)
	}

	fmt.Printf("SYN721 token %s successfully minted for owner: %s\n", tokenID, owner)
	return nil
}

// TransferToken transfers ownership of a specific SYN721 token to a new owner.
func (storage *SYN721Storage) TransferToken(tokenID, currentOwner, newOwner string) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	// Check if token exists
	if owner, exists := storage.TokenOwners[tokenID]; !exists {
		return errors.New("token does not exist")
	} else if owner != currentOwner {
		return errors.New("only the current owner can transfer this token")
	}

	// Update token ownership
	storage.TokenOwners[tokenID] = newOwner

	// Update the metadata with the new owner
	metadata := storage.TokenMetadata[tokenID]
	metadata.TokenOwner = newOwner

	// Encrypt the updated metadata
	encryptedData, err := storage.Encryption.EncryptData(fmt.Sprintf("%v", metadata), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting updated metadata: %v", err)
	}
	metadata.EncryptedData = encryptedData

	// Log the transfer in the ledger
	if err := storage.Ledger.RecordOwnershipTransfer(tokenID, currentOwner, newOwner); err != nil {
		return fmt.Errorf("error recording ownership transfer in ledger: %v", err)
	}

	fmt.Printf("Ownership of SYN721 token %s transferred from %s to %s.\n", tokenID, currentOwner, newOwner)
	return nil
}

// BurnToken burns an SYN721 token, removing it from the ledger and storage.
func (storage *SYN721Storage) BurnToken(tokenID, owner string) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	// Check if token exists
	if tokenOwner, exists := storage.TokenOwners[tokenID]; !exists {
		return errors.New("token does not exist")
	} else if tokenOwner != owner {
		return errors.New("only the token owner can burn this token")
	}

	// Remove the token from storage
	delete(storage.TokenOwners, tokenID)
	delete(storage.TokenMetadata, tokenID)

	// Log the burn in the ledger
	if err := storage.Ledger.RecordTokenBurn(tokenID, owner); err != nil {
		return fmt.Errorf("error recording token burn in ledger: %v", err)
	}

	fmt.Printf("SYN721 token %s successfully burned by owner: %s.\n", tokenID, owner)
	return nil
}

// GetTokenOwner retrieves the current owner of an SYN721 token.
func (storage *SYN721Storage) GetTokenOwner(tokenID string) (string, error) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	if owner, exists := storage.TokenOwners[tokenID]; exists {
		return owner, nil
	}

	return "", errors.New("token does not exist")
}

// GetTokenMetadata retrieves and decrypts the metadata for an SYN721 token.
func (storage *SYN721Storage) GetTokenMetadata(tokenID string) (*SYN721Metadata, error) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	metadata, exists := storage.TokenMetadata[tokenID]
	if !exists {
		return nil, errors.New("token metadata does not exist")
	}

	// Decrypt the metadata
	decryptedData, err := storage.Encryption.DecryptData(metadata.EncryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error decrypting token metadata: %v", err)
	}

	fmt.Printf("Retrieved metadata for SYN721 token %s: %v\n", tokenID, decryptedData)
	return metadata, nil
}

// VerifyTokenExistence checks if a specific SYN721 token exists.
func (storage *SYN721Storage) VerifyTokenExistence(tokenID string) bool {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	_, exists := storage.TokenOwners[tokenID]
	return exists
}
