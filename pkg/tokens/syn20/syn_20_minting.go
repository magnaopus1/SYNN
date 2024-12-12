package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// MintingManager handles minting of new SYN20 tokens.
type MintingManager struct {
	mutex      sync.Mutex                 // Thread safety
	TotalSupply uint64                    // Total supply of SYN20 tokens
	Ledger     *ledger.Ledger             // Reference to the blockchain ledger
	Consensus  *synnergy_consensus.Engine // Consensus engine for validation
	Encryption *encryption.Encryption     // Encryption service
	Owner      string                     // Owner of the contract who is allowed to mint tokens
}

// NewMintingManager initializes a new MintingManager for SYN20 tokens.
func NewMintingManager(initialSupply uint64, owner string, ledgerInstance *ledger.Ledger, consensus *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *MintingManager {
	return &MintingManager{
		TotalSupply: initialSupply,
		Ledger:      ledgerInstance,
		Consensus:   consensus,
		Encryption:  encryptionService,
		Owner:       owner,
	}
}

// MintTokens mints new SYN20 tokens and records the transaction in the ledger.
func (mm *MintingManager) MintTokens(caller string, amount uint64, recipient string) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	// Validate the caller is the owner of the contract
	if caller != mm.Owner {
		return errors.New("only the contract owner can mint tokens")
	}

	// Validate the recipient address via Synnergy Consensus
	valid, err := mm.Consensus.ValidateAddress(recipient)
	if !valid || err != nil {
		return fmt.Errorf("recipient address validation failed: %v", err)
	}

	// Update total supply
	newTotalSupply := mm.TotalSupply + amount
	if newTotalSupply < mm.TotalSupply { // Check for overflow
		return errors.New("overflow detected, invalid minting request")
	}
	mm.TotalSupply = newTotalSupply

	// Record the minting transaction in the ledger
	mintingID := common.GenerateTransactionID()
	mintingRecord := fmt.Sprintf("Minted %d SYN20 tokens for %s", amount, recipient)
	encryptedRecord, err := mm.Encryption.EncryptData(mintingRecord, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting minting record: %v", err)
	}

	if err := mm.Ledger.RecordTransaction(mintingID, encryptedRecord); err != nil {
		return fmt.Errorf("error recording minting transaction in the ledger: %v", err)
	}

	fmt.Printf("Successfully minted %d SYN20 tokens for %s.\n", amount, recipient)
	return nil
}

// GetTotalSupply retrieves the total supply of SYN20 tokens.
func (mm *MintingManager) GetTotalSupply() uint64 {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	return mm.TotalSupply
}

// BurnTokens burns SYN20 tokens and updates the total supply.
func (mm *MintingManager) BurnTokens(caller string, amount uint64) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	// Validate the caller is the owner of the contract
	if caller != mm.Owner {
		return errors.New("only the contract owner can burn tokens")
	}

	// Check if enough supply exists to burn the requested amount
	if amount > mm.TotalSupply {
		return errors.New("not enough tokens to burn")
	}

	// Update the total supply
	mm.TotalSupply -= amount

	// Record the burn transaction in the ledger
	burnID := common.GenerateTransactionID()
	burnRecord := fmt.Sprintf("Burned %d SYN20 tokens", amount)
	encryptedRecord, err := mm.Encryption.EncryptData(burnRecord, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting burn record: %v", err)
	}

	if err := mm.Ledger.RecordTransaction(burnID, encryptedRecord); err != nil {
		return fmt.Errorf("error recording burn transaction in the ledger: %v", err)
	}

	fmt.Printf("Successfully burned %d SYN20 tokens.\n", amount)
	return nil
}

// ValidateMintingRequest validates if a minting request can proceed based on network conditions and consensus.
func (mm *MintingManager) ValidateMintingRequest(amount uint64, recipient string) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	// Validate recipient address via Synnergy Consensus
	valid, err := mm.Consensus.ValidateAddress(recipient)
	if !valid || err != nil {
		return fmt.Errorf("recipient validation failed: %v", err)
	}

	// Check if minting would cause overflow
	if mm.TotalSupply+amount < mm.TotalSupply {
		return errors.New("minting request would cause supply overflow")
	}

	return nil
}

// ScheduleMinting schedules the minting of new tokens for a future time or block height.
func (mm *MintingManager) ScheduleMinting(caller string, amount uint64, recipient string, scheduleTime time.Time) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	// Validate the caller is the owner
	if caller != mm.Owner {
		return errors.New("only the contract owner can schedule minting")
	}

	// Validate the recipient's address
	valid, err := mm.Consensus.ValidateAddress(recipient)
	if !valid || err != nil {
		return fmt.Errorf("recipient address validation failed: %v", err)
	}

	// Create the schedule record
	scheduledMintingID := common.GenerateTransactionID()
	mintingSchedule := fmt.Sprintf("Scheduled minting of %d SYN20 tokens for %s at %v", amount, recipient, scheduleTime)
	encryptedRecord, err := mm.Encryption.EncryptData(mintingSchedule, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting scheduled minting record: %v", err)
	}

	// Record the scheduled minting in the ledger
	if err := mm.Ledger.RecordScheduledTransaction(scheduledMintingID, encryptedRecord, scheduleTime); err != nil {
		return fmt.Errorf("error recording scheduled minting in the ledger: %v", err)
	}

	fmt.Printf("Scheduled minting of %d SYN20 tokens for %s at %v.\n", amount, recipient, scheduleTime)
	return nil
}
