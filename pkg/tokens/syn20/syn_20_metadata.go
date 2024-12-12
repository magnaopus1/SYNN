package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SYN20Metadata stores all the necessary metadata for an SYN20 token.
type SYN20Metadata struct {
	TokenName      string  // Name of the token
	TokenSymbol    string  // Symbol for the token (e.g., SYN)
	TotalSupply    float64 // Total supply of the token
	Decimals       uint8   // Number of decimal places
	TokenOwner     string  // The address of the token owner (creator)
	EncryptedData  string  // Encrypted metadata for security
	ValidatorID    string  // The ID of the validator that validated this metadata
	DecryptedData  string  // Decrypted data for local use
}

// SYN20Manager handles the management and creation of SYN20 token metadata.
type SYN20Manager struct {
	mutex        sync.Mutex                 // For thread-safe operations
	Metadata     *SYN20Metadata             // Metadata of the SYN20 token
	Ledger       *ledger.Ledger             // Reference to the ledger for storing token metadata
	Consensus    *synnergy_consensus.Engine // Synnergy Consensus engine
	Encryption   *encryption.Encryption     // Encryption service
}

// NewSYN20Manager initializes a new SYN20 token manager with the given parameters.
func NewSYN20Manager(ledgerInstance *ledger.Ledger, consensus *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *SYN20Manager {
	return &SYN20Manager{
		Ledger:     ledgerInstance,
		Consensus:  consensus,
		Encryption: encryptionService,
	}
}

// CreateToken initializes a new SYN20 token, encrypts its metadata, and stores it in the ledger.
func (sm *SYN20Manager) CreateToken(name, symbol string, totalSupply float64, decimals uint8, owner string) (*SYN20Metadata, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Create the metadata for the SYN20 token
	metadata := &SYN20Metadata{
		TokenName:   name,
		TokenSymbol: symbol,
		TotalSupply: totalSupply,
		Decimals:    decimals,
		TokenOwner:  owner,
	}

	// Encrypt the token metadata
	encryptedData, err := sm.Encryption.EncryptData(fmt.Sprintf("%v", metadata), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting token metadata: %v", err)
	}
	metadata.EncryptedData = encryptedData

	// Store the token metadata in the ledger
	if err := sm.Ledger.AddTokenMetadata(metadata); err != nil {
		return nil, fmt.Errorf("error storing token metadata in ledger: %v", err)
	}

	// Validate token creation using consensus
	if valid, err := sm.Consensus.ValidateTokenCreation(metadata); !valid || err != nil {
		return nil, fmt.Errorf("error validating token creation: %v", err)
	}

	fmt.Printf("Token %s (%s) successfully created with total supply: %f.\n", name, symbol, totalSupply)
	return metadata, nil
}

// UpdateTokenMetadata allows the token owner to update metadata (e.g., name or symbol).
func (sm *SYN20Manager) UpdateTokenMetadata(owner, name, symbol string) (*SYN20Metadata, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if the requester is the owner of the token
	if owner != sm.Metadata.TokenOwner {
		return nil, errors.New("only the token owner can update metadata")
	}

	// Update the token metadata
	sm.Metadata.TokenName = name
	sm.Metadata.TokenSymbol = symbol

	// Encrypt updated metadata
	encryptedData, err := sm.Encryption.EncryptData(fmt.Sprintf("%v", sm.Metadata), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting updated metadata: %v", err)
	}
	sm.Metadata.EncryptedData = encryptedData

	// Update metadata in the ledger
	if err := sm.Ledger.UpdateTokenMetadata(sm.Metadata); err != nil {
		return nil, fmt.Errorf("error updating token metadata in ledger: %v", err)
	}

	fmt.Printf("Token metadata updated to Name: %s, Symbol: %s.\n", name, symbol)
	return sm.Metadata, nil
}

// GetTokenMetadata retrieves and decrypts the token metadata.
func (sm *SYN20Manager) GetTokenMetadata() (*SYN20Metadata, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Decrypt the token metadata for display
	decryptedData, err := sm.Encryption.DecryptData(sm.Metadata.EncryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error decrypting token metadata: %v", err)
	}

	// Convert decrypted data back into SYN20Metadata structure (this would need proper unmarshalling in real-world systems)
	sm.Metadata.DecryptedData = decryptedData

	return sm.Metadata, nil
}

// IncreaseTotalSupply increases the total supply of the token by the specified amount.
func (sm *SYN20Manager) IncreaseTotalSupply(owner string, amount float64) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if the requester is the owner of the token
	if owner != sm.Metadata.TokenOwner {
		return errors.New("only the token owner can increase total supply")
	}

	// Increase the total supply
	sm.Metadata.TotalSupply += amount

	// Encrypt updated metadata
	encryptedData, err := sm.Encryption.EncryptData(fmt.Sprintf("%v", sm.Metadata), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting updated metadata: %v", err)
	}
	sm.Metadata.EncryptedData = encryptedData

	// Update metadata in the ledger
	if err := sm.Ledger.UpdateTokenMetadata(sm.Metadata); err != nil {
		return fmt.Errorf("error updating token metadata in ledger: %v", err)
	}

	fmt.Printf("Total supply increased by %f. New total supply: %f.\n", amount, sm.Metadata.TotalSupply)
	return nil
}

// BurnTokens reduces the total supply of the token by the specified amount.
func (sm *SYN20Manager) BurnTokens(owner string, amount float64) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if the requester is the owner of the token
	if owner != sm.Metadata.TokenOwner {
		return errors.New("only the token owner can burn tokens")
	}

	// Ensure that there are enough tokens to burn
	if sm.Metadata.TotalSupply < amount {
		return errors.New("insufficient token supply to burn")
	}

	// Burn tokens by reducing the total supply
	sm.Metadata.TotalSupply -= amount

	// Encrypt updated metadata
	encryptedData, err := sm.Encryption.EncryptData(fmt.Sprintf("%v", sm.Metadata), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting updated metadata: %v", err)
	}
	sm.Metadata.EncryptedData = encryptedData

	// Update metadata in the ledger
	if err := sm.Ledger.UpdateTokenMetadata(sm.Metadata); err != nil {
		return fmt.Errorf("error updating token metadata in ledger: %v", err)
	}

	fmt.Printf("Tokens burned: %f. New total supply: %f.\n", amount, sm.Metadata.TotalSupply)
	return nil
}

// TransferOwnership transfers the ownership of the token to a new address.
func (sm *SYN20Manager) TransferOwnership(currentOwner, newOwner string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if the requester is the current owner of the token
	if currentOwner != sm.Metadata.TokenOwner {
		return errors.New("only the current owner can transfer ownership")
	}

	// Transfer ownership to the new address
	sm.Metadata.TokenOwner = newOwner

	// Encrypt updated metadata
	encryptedData, err := sm.Encryption.EncryptData(fmt.Sprintf("%v", sm.Metadata), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting updated metadata: %v", err)
	}
	sm.Metadata.EncryptedData = encryptedData

	// Update metadata in the ledger
	if err := sm.Ledger.UpdateTokenMetadata(sm.Metadata); err != nil {
		return fmt.Errorf("error updating token metadata in ledger: %v", err)
	}

	fmt.Printf("Token ownership transferred from %s to %s.\n", currentOwner, newOwner)
	return nil
}
