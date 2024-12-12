package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// SYN721Factory is responsible for creating and managing multiple SYN721 token contracts.
type SYN721Factory struct {
	mutex       sync.Mutex
	Ledger      *ledger.Ledger             // Ledger reference for recording transactions
	Consensus   *synnergy_consensus.Engine // Consensus engine for validation
	Encryption  *encryption.Encryption     // Encryption service for securing token data
	TokenManagers map[string]*SYN721TokenManager // Map of token managers created by the factory
}

// NewSYN721Factory initializes a new factory for SYN721 tokens.
func NewSYN721Factory(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *SYN721Factory {
	return &SYN721Factory{
		Ledger:       ledgerInstance,
		Consensus:    consensusEngine,
		Encryption:   encryptionService,
		TokenManagers: make(map[string]*SYN721TokenManager),
	}
}

// CreateTokenManager creates a new token manager for managing SYN721 tokens.
func (factory *SYN721Factory) CreateTokenManager(maxSupply uint64, managerID string) (*SYN721TokenManager, error) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Check if the manager already exists
	if _, exists := factory.TokenManagers[managerID]; exists {
		return nil, errors.New("token manager with this ID already exists")
	}

	// Initialize storage for the token manager
	storage := NewSYN721Storage()

	// Create a new token manager for SYN721
	tokenManager := NewSYN721TokenManager(factory.Ledger, factory.Consensus, factory.Encryption, storage)
	tokenManager.MaxSupply = maxSupply

	// Store the token manager
	factory.TokenManagers[managerID] = tokenManager

	// Log the creation of the token manager in the ledger
	err := factory.Ledger.RecordTokenManagerCreation(managerID, maxSupply)
	if err != nil {
		return nil, fmt.Errorf("failed to record token manager creation: %v", err)
	}

	fmt.Printf("Token Manager %s successfully created with max supply: %d.\n", managerID, maxSupply)
	return tokenManager, nil
}

// GetTokenManager retrieves an existing token manager by its ID.
func (factory *SYN721Factory) GetTokenManager(managerID string) (*SYN721TokenManager, error) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	manager, exists := factory.TokenManagers[managerID]
	if !exists {
		return nil, errors.New("token manager not found")
	}

	return manager, nil
}

// MintToken mints a new SYN721 token using a specified token manager.
func (factory *SYN721Factory) MintToken(managerID, tokenID, owner, tokenURI string) (*SYN721Token, error) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	manager, err := factory.GetTokenManager(managerID)
	if err != nil {
		return nil, err
	}

	// Use the token manager to mint the token
	token, err := manager.Mint(tokenID, owner, tokenURI)
	if err != nil {
		return nil, fmt.Errorf("failed to mint token: %v", err)
	}

	fmt.Printf("Token %s successfully minted under Manager %s for owner %s.\n", tokenID, managerID, owner)
	return token, nil
}

// BatchMintTokens allows batch minting of multiple SYN721 tokens using a specified token manager.
func (factory *SYN721Factory) BatchMintTokens(managerID string, mintRequests []TokenMintingRequest) ([]*SYN721Token, error) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	manager, err := factory.GetTokenManager(managerID)
	if err != nil {
		return nil, err
	}

	// Use the token manager to batch mint the tokens
	tokens, err := manager.BatchMint(mintRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to batch mint tokens: %v", err)
	}

	fmt.Printf("Batch minting of %d tokens successfully completed under Manager %s.\n", len(tokens), managerID)
	return tokens, nil
}

// ListTokenManagers returns a list of all token managers created by the factory.
func (factory *SYN721Factory) ListTokenManagers() []string {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	managerIDs := make([]string, 0, len(factory.TokenManagers))
	for managerID := range factory.TokenManagers {
		managerIDs = append(managerIDs, managerID)
	}

	return managerIDs
}

// DestroyTokenManager deletes an existing token manager and its associated data.
func (factory *SYN721Factory) DestroyTokenManager(managerID string) error {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	manager, exists := factory.TokenManagers[managerID]
	if !exists {
		return errors.New("token manager not found")
	}

	// Remove the manager from the factory
	delete(factory.TokenManagers, managerID)

	// Clear the storage associated with this manager
	manager.Storage.Clear()

	// Log the deletion in the ledger
	err := factory.Ledger.RecordTokenManagerDeletion(managerID)
	if err != nil {
		return fmt.Errorf("failed to record token manager deletion: %v", err)
	}

	fmt.Printf("Token Manager %s successfully destroyed.\n", managerID)
	return nil
}
