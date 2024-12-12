package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SYN20TokenFactory is responsible for creating and managing new SYN20 tokens.
type SYN20TokenFactory struct {
	mutex        sync.Mutex                 // Thread safety
	Ledger       *ledger.Ledger             // Reference to the ledger
	Consensus    *synnergy_consensus.Engine // Consensus engine for validation
	Encryption   *encryption.Encryption     // Encryption service
	DeployedTokens map[string]*MintingManager // Deployed SYN20 token contracts
}

// NewSYN20TokenFactory initializes a new SYN20 token factory.
func NewSYN20TokenFactory(ledgerInstance *ledger.Ledger, consensus *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *SYN20TokenFactory {
	return &SYN20TokenFactory{
		Ledger:        ledgerInstance,
		Consensus:     consensus,
		Encryption:    encryptionService,
		DeployedTokens: make(map[string]*MintingManager),
	}
}

// CreateNewSYN20Token creates and registers a new SYN20 token contract.
func (factory *SYN20TokenFactory) CreateNewSYN20Token(tokenName string, symbol string, initialSupply uint64, owner string) (string, error) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Validate owner address using Synnergy Consensus
	valid, err := factory.Consensus.ValidateAddress(owner)
	if !valid || err != nil {
		return "", fmt.Errorf("token owner validation failed: %v", err)
	}

	// Generate a unique token contract ID
	contractID := common.GenerateUniqueID()

	// Initialize new SYN20 token with minting functionality
	mintingManager := NewMintingManager(initialSupply, owner, factory.Ledger, factory.Consensus, factory.Encryption)

	// Register the token in the ledger
	tokenData := fmt.Sprintf("SYN20 Token: %s (Symbol: %s), Initial Supply: %d, Owner: %s", tokenName, symbol, initialSupply, owner)
	encryptedTokenData, err := factory.Encryption.EncryptData(tokenData, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("error encrypting token data: %v", err)
	}

	if err := factory.Ledger.RegisterNewContract(contractID, encryptedTokenData); err != nil {
		return "", fmt.Errorf("error registering new SYN20 token contract: %v", err)
	}

	// Store the minting manager for further operations
	factory.DeployedTokens[contractID] = mintingManager

	fmt.Printf("Successfully deployed SYN20 token contract with ID %s.\n", contractID)
	return contractID, nil
}

// MintTokens mints new tokens for a deployed SYN20 contract.
func (factory *SYN20TokenFactory) MintTokens(contractID, caller string, amount uint64, recipient string) error {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Retrieve the minting manager for the token
	mintingManager, exists := factory.DeployedTokens[contractID]
	if !exists {
		return errors.New("token contract not found")
	}

	// Proceed with the minting operation
	if err := mintingManager.MintTokens(caller, amount, recipient); err != nil {
		return fmt.Errorf("error minting tokens: %v", err)
	}

	fmt.Printf("Minted %d tokens for recipient %s on contract %s.\n", amount, recipient, contractID)
	return nil
}

// BurnTokens burns tokens for a deployed SYN20 contract.
func (factory *SYN20TokenFactory) BurnTokens(contractID, caller string, amount uint64) error {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Retrieve the minting manager for the token
	mintingManager, exists := factory.DeployedTokens[contractID]
	if !exists {
		return errors.New("token contract not found")
	}

	// Proceed with the burning operation
	if err := mintingManager.BurnTokens(caller, amount); err != nil {
		return fmt.Errorf("error burning tokens: %v", err)
	}

	fmt.Printf("Burned %d tokens on contract %s.\n", amount, contractID)
	return nil
}

// GetTotalSupply retrieves the total supply of a deployed SYN20 token contract.
func (factory *SYN20TokenFactory) GetTotalSupply(contractID string) (uint64, error) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Retrieve the minting manager for the token
	mintingManager, exists := factory.DeployedTokens[contractID]
	if !exists {
		return 0, errors.New("token contract not found")
	}

	// Get the total supply of tokens
	totalSupply := mintingManager.GetTotalSupply()
	return totalSupply, nil
}

// ValidateTokenContract ensures that the contract exists and is valid.
func (factory *SYN20TokenFactory) ValidateTokenContract(contractID string) error {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Check if the contract ID exists in the deployed tokens
	if _, exists := factory.DeployedTokens[contractID]; !exists {
		return errors.New("token contract not found")
	}

	// Additional logic to validate the integrity of the contract could be added here
	return nil
}

// ListDeployedTokens returns a list of all deployed SYN20 tokens.
func (factory *SYN20TokenFactory) ListDeployedTokens() map[string]string {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	deployedTokens := make(map[string]string)
	for contractID, mintingManager := range factory.DeployedTokens {
		owner := mintingManager.Owner
		totalSupply := mintingManager.GetTotalSupply()
		deployedTokens[contractID] = fmt.Sprintf("Owner: %s, Total Supply: %d", owner, totalSupply)
	}

	return deployedTokens
}
