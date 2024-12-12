package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SYN20BatchTransferManager handles the batch transfer of SYN20 tokens.
type SYN20BatchTransferManager struct {
	mutex      sync.Mutex
	Ledger     *ledger.Ledger              // Reference to the ledger for transaction logging
	Consensus  *synnergy_consensus.Engine  // Consensus engine for validating transfers
	Encryption *encryption.Encryption      // Encryption service for secure data handling
	Contracts  map[string]*SYN20Contract   // Token contracts managed for transfers
}

// NewSYN20BatchTransferManager initializes a new manager for batch transfers of SYN20 tokens.
func NewSYN20BatchTransferManager(ledgerInstance *ledger.Ledger, consensus *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *SYN20BatchTransferManager {
	return &SYN20BatchTransferManager{
		Ledger:     ledgerInstance,
		Consensus:  consensus,
		Encryption: encryptionService,
		Contracts:  make(map[string]*SYN20Contract),
	}
}

// BatchTransfer performs a batch transfer of SYN20 tokens from a single sender to multiple recipients.
func (btm *SYN20BatchTransferManager) BatchTransfer(contractID string, sender string, transfers map[string]uint64) error {
	btm.mutex.Lock()
	defer btm.mutex.Unlock()

	// Retrieve the contract
	contract, exists := btm.Contracts[contractID]
	if !exists {
		return errors.New("token contract not found")
	}

	// Validate the sender's ownership and balance using Synnergy Consensus
	valid, err := btm.Consensus.ValidateOwnership(contract.Owner, sender)
	if !valid || err != nil {
		return fmt.Errorf("sender validation failed: %v", err)
	}

	// Calculate the total amount to transfer
	totalTransferAmount := uint64(0)
	for _, amount := range transfers {
		totalTransferAmount += amount
	}

	// Check if the sender has enough balance to cover the total amount
	if contract.Balances[sender] < totalTransferAmount {
		return errors.New("insufficient balance for batch transfer")
	}

	// Perform the transfers
	for recipient, amount := range transfers {
		if amount == 0 {
			continue // Skip zero-amount transfers
		}

		// Update sender's and recipient's balances
		contract.Balances[sender] -= amount
		contract.Balances[recipient] += amount

		// Log the transfer in the ledger
		transferEvent := fmt.Sprintf("Transferred %d tokens from %s to %s in contract %s", amount, sender, recipient, contractID)
		encryptedTransferEvent, err := btm.Encryption.EncryptData(transferEvent, common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("error encrypting transfer event for recipient %s: %v", recipient, err)
		}

		err = btm.Ledger.RecordTransaction(contractID, encryptedTransferEvent)
		if err != nil {
			return fmt.Errorf("error recording transfer event for recipient %s: %v", recipient, err)
		}

		fmt.Printf("Transferred %d tokens from %s to %s in contract %s\n", amount, sender, recipient, contractID)
	}

	return nil
}

// BatchTransferFromMultipleSenders performs a batch transfer of tokens from multiple senders to a single recipient.
func (btm *SYN20BatchTransferManager) BatchTransferFromMultipleSenders(contractID string, transfers map[string]uint64, recipient string) error {
	btm.mutex.Lock()
	defer btm.mutex.Unlock()

	// Retrieve the contract
	contract, exists := btm.Contracts[contractID]
	if !exists {
		return errors.New("token contract not found")
	}

	totalTransferAmount := uint64(0)

	// Process each sender's transfer
	for sender, amount := range transfers {
		if amount == 0 {
			continue // Skip zero-amount transfers
		}

		// Validate each sender's ownership and balance
		valid, err := btm.Consensus.ValidateOwnership(contract.Owner, sender)
		if !valid || err != nil {
			return fmt.Errorf("sender validation failed for %s: %v", sender, err)
		}

		if contract.Balances[sender] < amount {
			return fmt.Errorf("insufficient balance for sender %s", sender)
		}

		// Update balances
		contract.Balances[sender] -= amount
		contract.Balances[recipient] += amount
		totalTransferAmount += amount

		// Log the transfer in the ledger
		transferEvent := fmt.Sprintf("Transferred %d tokens from %s to %s in contract %s", amount, sender, recipient, contractID)
		encryptedTransferEvent, err := btm.Encryption.EncryptData(transferEvent, common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("error encrypting transfer event for sender %s: %v", sender, err)
		}

		err = btm.Ledger.RecordTransaction(contractID, encryptedTransferEvent)
		if err != nil {
			return fmt.Errorf("error recording transfer event for sender %s: %v", sender, err)
		}
	}

	fmt.Printf("Successfully transferred a total of %d tokens to %s from multiple senders in contract %s.\n", totalTransferAmount, recipient, contractID)
	return nil
}

// RegisterContract registers a new SYN20 contract for batch transfer operations.
func (btm *SYN20BatchTransferManager) RegisterContract(contract *SYN20Contract) error {
	btm.mutex.Lock()
	defer btm.mutex.Unlock()

	if _, exists := btm.Contracts[contract.TokenName]; exists {
		return errors.New("contract already registered")
	}

	btm.Contracts[contract.TokenName] = contract
	fmt.Printf("Contract %s registered for batch transfer operations.\n", contract.TokenName)
	return nil
}
