package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// SYN721TransactionManager handles transactions involving SYN721 tokens.
type SYN721TransactionManager struct {
	mutex       sync.Mutex                 // For thread safety
	Storage     *SYN721Storage             // Reference to the token storage
	Ledger      *ledger.Ledger             // Ledger for recording transactions
	Consensus   *synnergy_consensus.Engine // Synnergy Consensus engine for validating transactions
	Encryption  *encryption.Encryption     // Encryption service for securing transactions
}

// NewSYN721TransactionManager initializes a new transaction manager for SYN721 tokens.
func NewSYN721TransactionManager(storage *SYN721Storage, ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *SYN721TransactionManager {
	return &SYN721TransactionManager{
		Storage:    storage,
		Ledger:     ledgerInstance,
		Consensus:  consensusEngine,
		Encryption: encryptionService,
	}
}

// MintToken mints a new SYN721 token (NFT) and assigns it to the specified owner.
func (tm *SYN721TransactionManager) MintToken(tokenID, owner, tokenURI string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Encrypt the token data for security
	encryptedURI, err := tm.Encryption.EncryptData(tokenURI, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting token URI: %v", err)
	}

	// Validate the minting operation using Synnergy Consensus
	if valid, err := tm.Consensus.ValidateTokenMint(tokenID, owner, encryptedURI); !valid || err != nil {
		return fmt.Errorf("error validating token mint: %v", err)
	}

	// Mint the token and store it in the storage
	err = tm.Storage.MintToken(tokenID, owner, encryptedURI)
	if err != nil {
		return fmt.Errorf("error minting token: %v", err)
	}

	// Log the minting transaction in the ledger
	err = tm.Ledger.RecordMintTransaction(tokenID, owner)
	if err != nil {
		return fmt.Errorf("error recording mint transaction in ledger: %v", err)
	}

	fmt.Printf("Token %s minted successfully for owner %s.\n", tokenID, owner)
	return nil
}

// TransferToken transfers an SYN721 token from one owner to another.
func (tm *SYN721TransactionManager) TransferToken(tokenID, currentOwner, newOwner string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate the transfer operation using Synnergy Consensus
	if valid, err := tm.Consensus.ValidateTokenTransfer(tokenID, currentOwner, newOwner); !valid || err != nil {
		return fmt.Errorf("error validating token transfer: %v", err)
	}

	// Transfer the token and update ownership in the storage
	err := tm.Storage.TransferToken(tokenID, currentOwner, newOwner)
	if err != nil {
		return fmt.Errorf("error transferring token: %v", err)
	}

	// Log the transfer transaction in the ledger
	err = tm.Ledger.RecordTransferTransaction(tokenID, currentOwner, newOwner)
	if err != nil {
		return fmt.Errorf("error recording transfer transaction in ledger: %v", err)
	}

	fmt.Printf("Token %s transferred from %s to %s.\n", tokenID, currentOwner, newOwner)
	return nil
}

// BurnToken burns (destroys) an SYN721 token, removing it from the system.
func (tm *SYN721TransactionManager) BurnToken(tokenID, owner string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate the burn operation using Synnergy Consensus
	if valid, err := tm.Consensus.ValidateTokenBurn(tokenID, owner); !valid || err != nil {
		return fmt.Errorf("error validating token burn: %v", err)
	}

	// Burn the token and remove it from the storage
	err := tm.Storage.BurnToken(tokenID, owner)
	if err != nil {
		return fmt.Errorf("error burning token: %v", err)
	}

	// Log the burn transaction in the ledger
	err = tm.Ledger.RecordBurnTransaction(tokenID, owner)
	if err != nil {
		return fmt.Errorf("error recording burn transaction in ledger: %v", err)
	}

	fmt.Printf("Token %s burned successfully by owner %s.\n", tokenID, owner)
	return nil
}

// GetTokenOwner retrieves the current owner of an SYN721 token.
func (tm *SYN721TransactionManager) GetTokenOwner(tokenID string) (string, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	owner, err := tm.Storage.GetTokenOwner(tokenID)
	if err != nil {
		return "", fmt.Errorf("error retrieving token owner: %v", err)
	}

	return owner, nil
}

// GetTokenMetadata retrieves the metadata (URI) of an SYN721 token.
func (tm *SYN721TransactionManager) GetTokenMetadata(tokenID string) (string, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	metadata, err := tm.Storage.GetTokenMetadata(tokenID)
	if err != nil {
		return "", fmt.Errorf("error retrieving token metadata: %v", err)
	}

	// Decrypt the token metadata
	decryptedURI, err := tm.Encryption.DecryptData(metadata.EncryptedData, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("error decrypting token metadata: %v", err)
	}

	return decryptedURI, nil
}

// VerifyTokenExistence checks if an SYN721 token exists in the storage.
func (tm *SYN721TransactionManager) VerifyTokenExistence(tokenID string) (bool, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	exists := tm.Storage.VerifyTokenExistence(tokenID)
	return exists, nil
}
