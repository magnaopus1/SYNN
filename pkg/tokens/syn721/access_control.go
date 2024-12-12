package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// AccessControl manages access rights for SYN721 tokens.
type AccessControl struct {
	mutex      sync.Mutex                 // For thread safety
	Ledger     *ledger.Ledger             // Ledger reference for logging
	Consensus  *synnergy_consensus.Engine // Synnergy Consensus engine for validation
	Encryption *encryption.Encryption     // Encryption service for securing access control data
	Storage    *SYN721Storage             // Token storage reference
}

// NewAccessControl initializes a new access control manager for SYN721 tokens.
func NewAccessControl(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, storage *SYN721Storage) *AccessControl {
	return &AccessControl{
		Ledger:     ledgerInstance,
		Consensus:  consensusEngine,
		Encryption: encryptionService,
		Storage:    storage,
	}
}

// GrantApproval grants permission to another address to transfer a specific SYN721 token.
func (ac *AccessControl) GrantApproval(tokenID, owner, approved string) error {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	// Retrieve the token from storage
	token, exists := ac.Storage.GetTokenData(tokenID)
	if !exists {
		return errors.New("token not found")
	}

	// Ensure only the owner can grant approval
	if token.Owner != owner {
		return errors.New("only the owner can grant approval")
	}

	// Grant the approval
	token.Approved = approved

	// Encrypt approval transaction data
	encryptedData, err := ac.Encryption.EncryptData(fmt.Sprintf("Approved %s to transfer token %s", approved, tokenID), "")
	if err != nil {
		return fmt.Errorf("error encrypting approval transaction: %v", err)
	}

	// Validate the approval using Synnergy Consensus
	if valid, err := ac.Consensus.ValidateApproval(tokenID, owner, approved); !valid || err != nil {
		return fmt.Errorf("approval failed consensus validation: %v", err)
	}

	// Update approval in storage
	err = ac.Storage.UpdateTokenData(tokenID, token)
	if err != nil {
		return fmt.Errorf("failed to update token %s data: %v", tokenID, err)
	}

	// Log the approval in the ledger
	err = ac.Ledger.RecordTokenApproval(tokenID, owner, approved)
	if err != nil {
		return fmt.Errorf("failed to log approval in the ledger: %v", err)
	}

	fmt.Printf("Token %s approved for transfer by %s to %s.\n", tokenID, owner, approved)
	return nil
}

// RevokeApproval revokes the approval previously granted to an address for a specific SYN721 token.
func (ac *AccessControl) RevokeApproval(tokenID, owner string) error {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	// Retrieve the token from storage
	token, exists := ac.Storage.GetTokenData(tokenID)
	if !exists {
		return errors.New("token not found")
	}

	// Ensure only the owner can revoke approval
	if token.Owner != owner {
		return errors.New("only the owner can revoke approval")
	}

	// Revoke the approval
	token.Approved = ""

	// Encrypt revocation transaction data
	encryptedData, err := ac.Encryption.EncryptData(fmt.Sprintf("Revoked approval for token %s by owner %s", tokenID, owner), "")
	if err != nil {
		return fmt.Errorf("error encrypting revocation transaction: %v", err)
	}

	// Validate the revocation using Synnergy Consensus
	if valid, err := ac.Consensus.ValidateRevocation(tokenID, owner); !valid || err != nil {
		return fmt.Errorf("revocation failed consensus validation: %v", err)
	}

	// Update token data in storage
	err = ac.Storage.UpdateTokenData(tokenID, token)
	if err != nil {
		return fmt.Errorf("failed to update token %s data: %v", tokenID, err)
	}

	// Log the revocation in the ledger
	err = ac.Ledger.RecordTokenRevocation(tokenID, owner)
	if err != nil {
		return fmt.Errorf("failed to log revocation in the ledger: %v", err)
	}

	fmt.Printf("Approval for token %s revoked by owner %s.\n", tokenID, owner)
	return nil
}

// CheckApproval checks whether a specific address has approval to transfer a given token.
func (ac *AccessControl) CheckApproval(tokenID, address string) (bool, error) {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	// Retrieve the token from storage
	token, exists := ac.Storage.GetTokenData(tokenID)
	if !exists {
		return false, errors.New("token not found")
	}

	// Check if the address is approved to transfer the token
	if token.Approved == address {
		return true, nil
	}

	return false, nil
}

// TransferOwnership transfers ownership of an SYN721 token from one address to another.
func (ac *AccessControl) TransferOwnership(tokenID, currentOwner, newOwner string) error {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	// Retrieve the token from storage
	token, exists := ac.Storage.GetTokenData(tokenID)
	if !exists {
		return errors.New("token not found")
	}

	// Ensure the current owner is initiating the transfer
	if token.Owner != currentOwner {
		return errors.New("only the current owner can transfer ownership")
	}

	// Transfer ownership
	token.Owner = newOwner
	token.Approved = "" // Clear any approvals

	// Encrypt transfer transaction data
	encryptedData, err := ac.Encryption.EncryptData(fmt.Sprintf("Transfer token %s from %s to %s", tokenID, currentOwner, newOwner), "")
	if err != nil {
		return fmt.Errorf("error encrypting transfer transaction: %v", err)
	}

	// Validate the transfer using Synnergy Consensus
	if valid, err := ac.Consensus.ValidateTokenTransfer(tokenID, currentOwner, newOwner); !valid || err != nil {
		return fmt.Errorf("token transfer failed consensus validation: %v", err)
	}

	// Update the token data in storage
	err = ac.Storage.UpdateTokenData(tokenID, token)
	if err != nil {
		return fmt.Errorf("failed to update token %s data: %v", tokenID, err)
	}

	// Log the transfer in the ledger
	err = ac.Ledger.RecordTokenTransfer(tokenID, currentOwner, newOwner)
	if err != nil {
		return fmt.Errorf("failed to log transfer in the ledger: %v", err)
	}

	fmt.Printf("Token %s ownership transferred from %s to %s.\n", tokenID, currentOwner, newOwner)
	return nil
}
