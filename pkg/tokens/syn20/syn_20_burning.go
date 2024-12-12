package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SYN20BurningManager manages the burning of SYN20 tokens on the blockchain.
type SYN20BurningManager struct {
	mutex       sync.Mutex
	Ledger      *ledger.Ledger              // Reference to the ledger for transaction recording
	Consensus   *synnergy_consensus.Engine  // Consensus engine for transaction validation
	Encryption  *encryption.Encryption      // Encryption service for secure data handling
	Contracts   map[string]*SYN20Contract   // Token contracts managed for burning operations
}

// NewSYN20BurningManager initializes a new manager for burning SYN20 tokens.
func NewSYN20BurningManager(ledgerInstance *ledger.Ledger, consensus *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *SYN20BurningManager {
	return &SYN20BurningManager{
		Ledger:      ledgerInstance,
		Consensus:   consensus,
		Encryption:  encryptionService,
		Contracts:   make(map[string]*SYN20Contract),
	}
}

// BurnTokens burns the specified amount of tokens from the owner's supply and updates the total supply.
func (bm *SYN20BurningManager) BurnTokens(contractID, owner string, amount uint64) error {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	// Retrieve the contract
	contract, exists := bm.Contracts[contractID]
	if !exists {
		return errors.New("token contract not found")
	}

	// Validate the owner's authority to burn tokens using the Synnergy Consensus
	valid, err := bm.Consensus.ValidateOwnership(contract.Owner, owner)
	if !valid || err != nil {
		return fmt.Errorf("owner validation failed: %v", err)
	}

	// Perform the burning of tokens
	if amount > contract.TotalSupply {
		return errors.New("insufficient tokens to burn")
	}
	contract.TotalSupply -= amount

	// Record the burning event in the ledger
	burnEvent := fmt.Sprintf("Burned %d tokens from contract %s by owner %s", amount, contractID, owner)
	encryptedBurnEvent, err := bm.Encryption.EncryptData(burnEvent, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting burn event: %v", err)
	}

	err = bm.Ledger.RecordBurnEvent(contractID, encryptedBurnEvent)
	if err != nil {
		return fmt.Errorf("error recording burn event in the ledger: %v", err)
	}

	fmt.Printf("Successfully burned %d tokens from contract %s. New total supply: %d\n", amount, contractID, contract.TotalSupply)
	return nil
}

// GetBurningHistory retrieves the burn event history for a specific SYN20 contract.
func (bm *SYN20BurningManager) GetBurningHistory(contractID string) ([]string, error) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	// Retrieve the burn history from the ledger
	history, err := bm.Ledger.GetBurnHistory(contractID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving burn history: %v", err)
	}

	// Decrypt the burn history
	var decryptedHistory []string
	for _, event := range history {
		decryptedEvent, err := bm.Encryption.DecryptData(event, common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("error decrypting burn event: %v", err)
		}
		decryptedHistory = append(decryptedHistory, decryptedEvent)
	}

	return decryptedHistory, nil
}

// RegisterContract registers a new SYN20 contract for burning operations.
func (bm *SYN20BurningManager) RegisterContract(contract *SYN20Contract) error {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if _, exists := bm.Contracts[contract.TokenName]; exists {
		return errors.New("contract already registered")
	}

	bm.Contracts[contract.TokenName] = contract
	fmt.Printf("Contract %s registered for burning operations.\n", contract.TokenName)
	return nil
}

// BurnTokensFromMultipleOwners allows burning tokens from multiple owners at once.
func (bm *SYN20BurningManager) BurnTokensFromMultipleOwners(contractID string, owners map[string]uint64) error {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	// Retrieve the contract
	contract, exists := bm.Contracts[contractID]
	if !exists {
		return errors.New("token contract not found")
	}

	totalBurnAmount := uint64(0)

	// Loop through all owners and burn their respective amounts
	for owner, amount := range owners {
		// Validate ownership
		valid, err := bm.Consensus.ValidateOwnership(contract.Owner, owner)
		if !valid || err != nil {
			return fmt.Errorf("owner validation failed for %s: %v", owner, err)
		}

		// Perform the burn
		if amount > contract.TotalSupply {
			return fmt.Errorf("insufficient tokens to burn for owner %s", owner)
		}
		contract.TotalSupply -= amount
		totalBurnAmount += amount

		// Log burn event for each owner
		burnEvent := fmt.Sprintf("Burned %d tokens from contract %s by owner %s", amount, contractID, owner)
		encryptedBurnEvent, err := bm.Encryption.EncryptData(burnEvent, common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("error encrypting burn event for owner %s: %v", owner, err)
		}

		err = bm.Ledger.RecordBurnEvent(contractID, encryptedBurnEvent)
		if err != nil {
			return fmt.Errorf("error recording burn event for owner %s in the ledger: %v", owner, err)
		}
	}

	fmt.Printf("Successfully burned a total of %d tokens from contract %s by multiple owners. New total supply: %d\n", totalBurnAmount, contractID, contract.TotalSupply)
	return nil
}
