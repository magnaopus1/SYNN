package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// OwnershipManager handles ownership transfer and management of SYN20 tokens.
type OwnershipManager struct {
	mutex       sync.Mutex                 // Thread-safe operations
	Owner       string                     // Current owner of the token contract
	Ledger      *ledger.Ledger             // Reference to the blockchain ledger
	Consensus   *synnergy_consensus.Engine // Synnergy Consensus engine for ownership validation
	Encryption  *encryption.Encryption     // Encryption service for securing data
}

// NewOwnershipManager initializes a new ownership manager for SYN20 tokens.
func NewOwnershipManager(initialOwner string, ledgerInstance *ledger.Ledger, consensus *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *OwnershipManager {
	return &OwnershipManager{
		Owner:      initialOwner,
		Ledger:     ledgerInstance,
		Consensus:  consensus,
		Encryption: encryptionService,
	}
}

// TransferOwnership transfers ownership of the SYN20 token contract to a new address.
func (om *OwnershipManager) TransferOwnership(currentOwner, newOwner string) error {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	// Validate that the current owner is allowed to transfer ownership
	if currentOwner != om.Owner {
		return errors.New("only the current owner can transfer ownership")
	}

	// Validate the new owner's address via Synnergy Consensus
	valid, err := om.Consensus.ValidateAddress(newOwner)
	if !valid || err != nil {
		return fmt.Errorf("ownership transfer validation failed: %v", err)
	}

	// Transfer ownership
	om.Owner = newOwner

	// Log ownership transfer in the ledger
	ownershipID := common.GenerateTransactionID()
	ownershipRecord := fmt.Sprintf("Ownership transferred from %s to %s", currentOwner, newOwner)
	encryptedRecord, err := om.Encryption.EncryptData(ownershipRecord, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting ownership transfer record: %v", err)
	}

	if err := om.Ledger.RecordOwnershipTransfer(ownershipID, encryptedRecord); err != nil {
		return fmt.Errorf("error storing ownership transfer in ledger: %v", err)
	}

	fmt.Printf("Ownership successfully transferred from %s to %s.\n", currentOwner, newOwner)
	return nil
}

// GetOwner retrieves the current owner of the SYN20 token contract.
func (om *OwnershipManager) GetOwner() string {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	return om.Owner
}

// ValidateOwnership ensures that the caller is the current owner of the token contract.
func (om *OwnershipManager) ValidateOwnership(caller string) error {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	if caller != om.Owner {
		return errors.New("caller is not the owner of the contract")
	}

	return nil
}

// RecordOwnershipChange stores an ownership change in the ledger securely.
func (om *OwnershipManager) RecordOwnershipChange(previousOwner, newOwner string) error {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	// Record the ownership change details
	record := fmt.Sprintf("Ownership transferred from %s to %s", previousOwner, newOwner)

	// Encrypt the ownership change record
	encryptedRecord, err := om.Encryption.EncryptData(record, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting ownership change record: %v", err)
	}

	// Log the ownership change in the ledger
	ownershipID := common.GenerateTransactionID()
	if err := om.Ledger.RecordOwnershipTransfer(ownershipID, encryptedRecord); err != nil {
		return fmt.Errorf("error storing ownership change in ledger: %v", err)
	}

	fmt.Printf("Ownership change from %s to %s recorded in ledger.\n", previousOwner, newOwner)
	return nil
}

// RevokeOwnership revokes the current ownership (typically in the case of a lost private key).
func (om *OwnershipManager) RevokeOwnership(currentOwner string) error {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	// Validate that the caller is the current owner
	if currentOwner != om.Owner {
		return errors.New("only the current owner can revoke ownership")
	}

	// Revoke the ownership by setting it to a zero address (or some default "revoked" state)
	om.Owner = "0x0" // or a special revoked address

	// Log the ownership revocation in the ledger
	ownershipID := common.GenerateTransactionID()
	revocationRecord := fmt.Sprintf("Ownership revoked by %s", currentOwner)
	encryptedRecord, err := om.Encryption.EncryptData(revocationRecord, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting ownership revocation record: %v", err)
	}

	if err := om.Ledger.RecordOwnershipTransfer(ownershipID, encryptedRecord); err != nil {
		return fmt.Errorf("error storing ownership revocation in ledger: %v", err)
	}

	fmt.Printf("Ownership revoked by %s.\n", currentOwner)
	return nil
}

// ReinstateOwnership reinstates ownership to a previously validated owner.
func (om *OwnershipManager) ReinstateOwnership(previousOwner, newOwner string) error {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	// Validate that the previous owner is eligible for ownership reinstatement
	valid, err := om.Consensus.ValidateAddress(previousOwner)
	if !valid || err != nil {
		return fmt.Errorf("reinstatement validation failed for previous owner %s: %v", previousOwner, err)
	}

	// Reassign ownership to the new owner
	om.Owner = newOwner

	// Log the reinstatement in the ledger
	ownershipID := common.GenerateTransactionID()
	reinstatementRecord := fmt.Sprintf("Ownership reinstated from %s to %s", previousOwner, newOwner)
	encryptedRecord, err := om.Encryption.EncryptData(reinstatementRecord, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting ownership reinstatement record: %v", err)
	}

	if err := om.Ledger.RecordOwnershipTransfer(ownershipID, encryptedRecord); err != nil {
		return fmt.Errorf("error storing ownership reinstatement in ledger: %v", err)
	}

	fmt.Printf("Ownership reinstated from %s to %s.\n", previousOwner, newOwner)
	return nil
}
