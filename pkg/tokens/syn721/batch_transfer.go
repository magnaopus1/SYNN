package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// SYN721BatchTransfer manages the batch transfer of SYN721 tokens (NFTs).
type SYN721BatchTransfer struct {
	mutex      sync.Mutex                 // For thread-safe operations
	Ledger     *ledger.Ledger             // Reference to the ledger for recording transfers
	Consensus  *synnergy_consensus.Engine // Synnergy Consensus engine for validation
	Encryption *encryption.Encryption     // Encryption service for securing transfer data
	Storage    *SYN721Storage             // Storage for SYN721 token data
}

// NewSYN721BatchTransfer initializes a new batch transfer manager for SYN721 tokens.
func NewSYN721BatchTransfer(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, storage *SYN721Storage) *SYN721BatchTransfer {
	return &SYN721BatchTransfer{
		Ledger:     ledgerInstance,
		Consensus:  consensusEngine,
		Encryption: encryptionService,
		Storage:    storage,
	}
}

// BatchTransfer performs the batch transfer of multiple SYN721 tokens from one address to multiple recipients.
func (bt *SYN721BatchTransfer) BatchTransfer(owner string, transfers map[string]string) error {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	// Iterate over each token ID and its respective recipient
	for tokenID, recipient := range transfers {
		// Retrieve the token from storage
		token, exists := bt.Storage.GetTokenData(tokenID)
		if !exists {
			return fmt.Errorf("token %s not found", tokenID)
		}

		// Ensure the caller is the owner of the token or has the appropriate approval
		if token.Owner != owner && token.Approved != owner {
			return fmt.Errorf("transfer not authorized for token %s", tokenID)
		}

		// Perform the transfer
		token.Owner = recipient
		token.Approved = "" // Clear any approvals

		// Encrypt transfer transaction data
		encryptedData, err := bt.Encryption.EncryptData(fmt.Sprintf("Transfer token %s from %s to %s", tokenID, owner, recipient), "")
		if err != nil {
			return fmt.Errorf("error encrypting transfer transaction for token %s: %v", tokenID, err)
		}

		// Validate the transfer using Synnergy Consensus
		if valid, err := bt.Consensus.ValidateTokenTransfer(tokenID, owner, recipient); !valid || err != nil {
			return fmt.Errorf("token %s transfer failed consensus validation: %v", tokenID, err)
		}

		// Update the token data in storage
		err = bt.Storage.UpdateTokenData(tokenID, token)
		if err != nil {
			return fmt.Errorf("failed to update token %s data: %v", tokenID, err)
		}

		// Record the transfer in the ledger
		err = bt.Ledger.RecordTokenTransfer(tokenID, owner, recipient)
		if err != nil {
			return fmt.Errorf("failed to record transfer for token %s in the ledger: %v", tokenID, err)
		}

		fmt.Printf("Token %s transferred from %s to %s.\n", tokenID, owner, recipient)
	}

	return nil
}

// BatchApprove allows a batch of tokens to be approved for transfer by another address.
func (bt *SYN721BatchTransfer) BatchApprove(owner string, approvals map[string]string) error {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	// Iterate over each token ID and the approved address
	for tokenID, approvedAddress := range approvals {
		// Retrieve the token from storage
		token, exists := bt.Storage.GetTokenData(tokenID)
		if !exists {
			return fmt.Errorf("token %s not found", tokenID)
		}

		// Ensure only the owner can approve the token
		if token.Owner != owner {
			return fmt.Errorf("only the owner can approve token %s", tokenID)
		}

		// Set the approved address for transfer rights
		token.Approved = approvedAddress

		// Encrypt approval transaction data
		encryptedData, err := bt.Encryption.EncryptData(fmt.Sprintf("Approve %s to transfer token %s", approvedAddress, tokenID), "")
		if err != nil {
			return fmt.Errorf("error encrypting approval transaction for token %s: %v", tokenID, err)
		}

		// Validate the approval using Synnergy Consensus
		if valid, err := bt.Consensus.ValidateApproval(tokenID, owner, approvedAddress); !valid || err != nil {
			return fmt.Errorf("approval for token %s failed consensus validation: %v", tokenID, err)
		}

		// Update token approval in storage
		err = bt.Storage.UpdateTokenData(tokenID, token)
		if err != nil {
			return fmt.Errorf("failed to update approval for token %s: %v", tokenID, err)
		}

		// Record the approval in the ledger
		err = bt.Ledger.RecordTokenApproval(tokenID, owner, approvedAddress)
		if err != nil {
			return fmt.Errorf("failed to record approval for token %s in the ledger: %v", tokenID, err)
		}

		fmt.Printf("Token %s approved for transfer by %s to %s.\n", tokenID, owner, approvedAddress)
	}

	return nil
}

// BatchBurn burns a batch of SYN721 tokens.
func (bt *SYN721BatchTransfer) BatchBurn(owner string, tokenIDs []string) error {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	// Iterate over each token to be burned
	for _, tokenID := range tokenIDs {
		// Retrieve the token from storage
		token, exists := bt.Storage.GetTokenData(tokenID)
		if !exists {
			return fmt.Errorf("token %s not found", tokenID)
		}

		// Ensure only the owner can burn the token
		if token.Owner != owner {
			return fmt.Errorf("only the owner can burn token %s", tokenID)
		}

		// Encrypt burn transaction data
		encryptedData, err := bt.Encryption.EncryptData(fmt.Sprintf("Burn token %s owned by %s", tokenID, owner), "")
		if err != nil {
			return fmt.Errorf("error encrypting burn transaction for token %s: %v", tokenID, err)
		}

		// Validate the burn using Synnergy Consensus
		if valid, err := bt.Consensus.ValidateTokenBurn(tokenID, owner); !valid || err != nil {
			return fmt.Errorf("token %s burn failed consensus validation: %v", tokenID, err)
		}

		// Remove the token from storage
		err = bt.Storage.RemoveTokenData(tokenID)
		if err != nil {
			return fmt.Errorf("failed to remove token %s from storage: %v", tokenID, err)
		}

		// Record the burn in the ledger
		err = bt.Ledger.RecordTokenBurn(tokenID, owner)
		if err != nil {
			return fmt.Errorf("failed to record token %s burn in the ledger: %v", tokenID, err)
		}

		fmt.Printf("Token %s burned by owner %s.\n", tokenID, owner)
	}

	return nil
}
