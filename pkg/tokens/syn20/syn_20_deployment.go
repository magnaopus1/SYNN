package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SYN20DeploymentManager manages the deployment of SYN20 token contracts.
type SYN20DeploymentManager struct {
	mutex       sync.Mutex
	Ledger      *ledger.Ledger              // Reference to the ledger for recording deployed contracts
	Consensus   *synnergy_consensus.Engine  // Consensus engine for validation
	Encryption  *encryption.Encryption      // Encryption service
	Deployments map[string]*SYN20Contract   // Deployed SYN20 contracts
}

// NewSYN20DeploymentManager initializes a new deployment manager for SYN20 tokens.
func NewSYN20DeploymentManager(ledgerInstance *ledger.Ledger, consensus *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *SYN20DeploymentManager {
	return &SYN20DeploymentManager{
		Ledger:      ledgerInstance,
		Consensus:   consensus,
		Encryption:  encryptionService,
		Deployments: make(map[string]*SYN20Contract),
	}
}

// DeployTokenContract deploys a new SYN20 token contract on the blockchain.
func (manager *SYN20DeploymentManager) DeployTokenContract(tokenName string, symbol string, initialSupply uint64, owner string) (string, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	// Validate owner address with the Synnergy Consensus
	valid, err := manager.Consensus.ValidateAddress(owner)
	if !valid || err != nil {
		return "", fmt.Errorf("owner address validation failed: %v", err)
	}

	// Create a unique contract ID
	contractID := common.GenerateUniqueID()

	// Initialize the new SYN20 contract
	contract := NewSYN20Contract(tokenName, symbol, initialSupply, owner)

	// Encrypt the contract details
	contractData := fmt.Sprintf("SYN20 Token: %s (Symbol: %s), Initial Supply: %d, Owner: %s", tokenName, symbol, initialSupply, owner)
	encryptedData, err := manager.Encryption.EncryptData(contractData, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("error encrypting contract data: %v", err)
	}

	// Register the contract in the ledger
	err = manager.Ledger.RegisterNewContract(contractID, encryptedData)
	if err != nil {
		return "", fmt.Errorf("error registering new token contract in the ledger: %v", err)
	}

	// Store the contract
	manager.Deployments[contractID] = contract

	fmt.Printf("SYN20 token contract successfully deployed with ID: %s\n", contractID)
	return contractID, nil
}

// GetContract retrieves an existing SYN20 contract by its contract ID.
func (manager *SYN20DeploymentManager) GetContract(contractID string) (*SYN20Contract, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	contract, exists := manager.Deployments[contractID]
	if !exists {
		return nil, errors.New("contract not found")
	}

	return contract, nil
}

// ListDeployedContracts lists all deployed SYN20 contracts.
func (manager *SYN20DeploymentManager) ListDeployedContracts() map[string]string {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	contracts := make(map[string]string)
	for contractID, contract := range manager.Deployments {
		contracts[contractID] = fmt.Sprintf("Token: %s, Symbol: %s, Total Supply: %d, Owner: %s",
			contract.TokenName, contract.Symbol, contract.TotalSupply, contract.Owner)
	}
	return contracts
}

// SYN20Contract represents a SYN20 token contract on the blockchain.
type SYN20Contract struct {
	TokenName   string  // Name of the token
	Symbol      string  // Token symbol
	TotalSupply uint64  // Total supply of the token
	Owner       string  // Owner of the token contract
}

// NewSYN20Contract initializes a new SYN20 token contract.
func NewSYN20Contract(tokenName, symbol string, initialSupply uint64, owner string) *SYN20Contract {
	return &SYN20Contract{
		TokenName:   tokenName,
		Symbol:      symbol,
		TotalSupply: initialSupply,
		Owner:       owner,
	}
}

// TransferTokens transfers tokens from the contract's owner to a recipient.
func (contract *SYN20Contract) TransferTokens(recipient string, amount uint64) error {
	if amount > contract.TotalSupply {
		return errors.New("insufficient supply for transfer")
	}

	contract.TotalSupply -= amount
	// Simulate transfer logic (in a real system, it would integrate with a ledger or token system)
	fmt.Printf("Transferred %d tokens to %s from %s\n", amount, recipient, contract.Owner)
	return nil
}

// MintTokens allows the contract owner to mint new tokens.
func (contract *SYN20Contract) MintTokens(amount uint64) {
	contract.TotalSupply += amount
	fmt.Printf("Minted %d new tokens. Total supply is now %d.\n", amount, contract.TotalSupply)
}

// BurnTokens allows the contract owner to burn tokens, reducing the total supply.
func (contract *SYN20Contract) BurnTokens(amount uint64) error {
	if amount > contract.TotalSupply {
		return errors.New("insufficient tokens to burn")
	}

	contract.TotalSupply -= amount
	fmt.Printf("Burned %d tokens. Total supply is now %d.\n", amount, contract.TotalSupply)
	return nil
}
