package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// SYN721OwnershipManager handles the ownership management for SYN721 tokens.
type SYN721OwnershipManager struct {
	mutex      sync.Mutex                 // For thread safety
	Ledger     *ledger.Ledger             // Reference to the ledger for transaction recording
	Consensus  *synnergy_consensus.Engine // Synnergy Consensus engine for validation
	Encryption *encryption.Encryption     // Encryption service for securing transaction data
	Tokens     map[string]*SYN721Token    // Map of all SYN721 tokens by their ID
}

// NewSYN721OwnershipManager initializes a new SYN721OwnershipManager.
func NewSYN721OwnershipManager(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *SYN721OwnershipManager {
	return &SYN721OwnershipManager{
		Ledger:     ledgerInstance,
		Consensus:  consensusEngine,
		Encryption: encryptionService,
		Tokens:     make(map[string]*SYN721Token),
	}
}

// TransferOwnership transfers ownership of an SYN721 token from one owner to another.
func (om *SYN721OwnershipManager) TransferOwnership(tokenID, currentOwner, newOwner string) error {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	// Retrieve the token
	token, exists := om.Tokens[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	// Verify that the current owner is the actual owner of the token
	if token.Owner != currentOwner {
		return errors.New("transfer failed: current owner does not match")
	}

	// Encrypt the ownership transfer data
	encryptedData, err := om.Encryption.EncryptData(fmt.Sprintf("Transfer ownership of token %s from %s to %s", tokenID, currentOwner, newOwner), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting ownership transfer data: %v", err)
	}

	// Validate the ownership transfer using Synnergy Consensus
	if valid, err := om.Consensus.ValidateOwnershipTransfer(tokenID, currentOwner, newOwner, encryptedData); !valid || err != nil {
		return fmt.Errorf("ownership transfer failed consensus validation: %v", err)
	}

	// Update the owner of the token
	token.Owner = newOwner

	// Record the ownership transfer in the ledger
	err = om.Ledger.RecordOwnershipTransfer(tokenID, currentOwner, newOwner)
	if err != nil {
		return fmt.Errorf("failed to log ownership transfer in the ledger: %v", err)
	}

	fmt.Printf("Ownership of token %s transferred from %s to %s.\n", tokenID, currentOwner, newOwner)
	return nil
}

// ApproveOwnershipTransfer allows the current owner to approve another address to transfer ownership on their behalf.
func (om *SYN721OwnershipManager) ApproveOwnershipTransfer(tokenID, owner, approvedAddress string) error {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	// Retrieve the token
	token, exists := om.Tokens[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	// Verify that the owner is the actual owner of the token
	if token.Owner != owner {
		return errors.New("only the owner can approve ownership transfer")
	}

	// Set the approved address for ownership transfer
	token.Approved = approvedAddress

	// Encrypt the approval data
	encryptedData, err := om.Encryption.EncryptData(fmt.Sprintf("Approve ownership transfer for token %s by owner %s to %s", tokenID, owner, approvedAddress), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting approval data: %v", err)
	}

	// Validate the approval using Synnergy Consensus
	if valid, err := om.Consensus.ValidateOwnershipApproval(tokenID, owner, approvedAddress, encryptedData); !valid || err != nil {
		return fmt.Errorf("approval failed consensus validation: %v", err)
	}

	// Record the approval in the ledger
	err = om.Ledger.RecordOwnershipApproval(tokenID, owner, approvedAddress)
	if err != nil {
		return fmt.Errorf("failed to log approval in the ledger: %v", err)
	}

	fmt.Printf("Ownership transfer for token %s approved by %s for %s.\n", tokenID, owner, approvedAddress)
	return nil
}

// RevokeOwnershipApproval allows the current owner to revoke a previously approved transfer address.
func (om *SYN721OwnershipManager) RevokeOwnershipApproval(tokenID, owner string) error {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	// Retrieve the token
	token, exists := om.Tokens[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	// Verify that the owner is the actual owner of the token
	if token.Owner != owner {
		return errors.New("only the owner can revoke transfer approval")
	}

	// Clear the approved address for ownership transfer
	token.Approved = ""

	// Encrypt the revocation data
	encryptedData, err := om.Encryption.EncryptData(fmt.Sprintf("Revoke ownership transfer approval for token %s by owner %s", tokenID, owner), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting revocation data: %v", err)
	}

	// Validate the revocation using Synnergy Consensus
	if valid, err := om.Consensus.ValidateOwnershipRevocation(tokenID, owner, encryptedData); !valid || err != nil {
		return fmt.Errorf("revocation failed consensus validation: %v", err)
	}

	// Record the revocation in the ledger
	err = om.Ledger.RecordOwnershipRevocation(tokenID, owner)
	if err != nil {
		return fmt.Errorf("failed to log revocation in the ledger: %v", err)
	}

	fmt.Printf("Ownership transfer approval for token %s revoked by owner %s.\n", tokenID, owner)
	return nil
}

// GetOwner retrieves the current owner of a specific SYN721 token.
func (om *SYN721OwnershipManager) GetOwner(tokenID string) (string, error) {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	// Retrieve the token
	token, exists := om.Tokens[tokenID]
	if !exists {
		return "", errors.New("token not found")
	}

	return token.Owner, nil
}

// IsApproved checks if an address is approved to transfer ownership of a specific SYN721 token.
func (om *SYN721OwnershipManager) IsApproved(tokenID, address string) (bool, error) {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	// Retrieve the token
	token, exists := om.Tokens[tokenID]
	if !exists {
		return false, errors.New("token not found")
	}

	return token.Approved == address, nil
}

// ListOwnedTokens lists all SYN721 tokens owned by a specific address.
func (om *SYN721OwnershipManager) ListOwnedTokens(owner string) ([]string, error) {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	var ownedTokens []string

	// Iterate through all tokens and check ownership
	for tokenID, token := range om.Tokens {
		if token.Owner == owner {
			ownedTokens = append(ownedTokens, tokenID)
		}
	}

	return ownedTokens, nil
}
