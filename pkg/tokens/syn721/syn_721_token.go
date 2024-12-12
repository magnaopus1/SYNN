package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// SYN721Token represents the structure of an NFT token under the SYN721 standard.
type SYN721Token struct {
	mutex         sync.Mutex                 // For thread safety
	TokenID       string                     // Unique identifier for the token (NFT)
	TokenURI      string                     // URI linking to the token's metadata (image, etc.)
	Owner         string                     // Owner of the NFT
	Approved      string                     // Approved address for transfer rights
	Ledger        *ledger.Ledger             // Reference to the ledger for recording transactions
	Consensus     *synnergy_consensus.Engine // Synnergy Consensus engine for validation
	Encryption    *encryption.Encryption     // Encryption service for secure transaction data
	Metadata      *SYN721Metadata            // Metadata for the NFT token
	RoyaltyInfo          common.RoyaltyInfo       `json:"royalty_info"`

}

// SYN721TokenManager manages the operations of SYN721 tokens.
type SYN721TokenManager struct {
	mutex       sync.Mutex                    // For thread safety
	Ledger      *ledger.Ledger                // Reference to the ledger for recording transactions
	Tokens      map[string]*SYN721Token       // Map of all tokens by TokenID
	Consensus   *synnergy_consensus.Engine    // Consensus engine for validation
	Encryption  *encryption.Encryption        // Encryption service for securing transactions
	Storage     *SYN721Storage                // Storage for SYN721 token data
	TotalSupply int                           // Total supply of NFTs minted
	MaxSupply   int                           // Maximum supply of NFTs allowed to be minted
}

// NewSYN721TokenManager initializes a new manager for SYN721 tokens.
func NewSYN721TokenManager(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, storage *SYN721Storage, maxSupply int) *SYN721TokenManager {
	return &SYN721TokenManager{
		Ledger:      ledgerInstance,
		Tokens:      make(map[string]*SYN721Token),
		Consensus:   consensusEngine,
		Encryption:  encryptionService,
		Storage:     storage,
		TotalSupply: 0,
		MaxSupply:   maxSupply,
	}
}

// Mint creates a new SYN721 token and records it in the ledger.
func (tm *SYN721TokenManager) Mint(tokenID, owner, tokenURI string) (*SYN721Token, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Check if the token ID already exists
	if _, exists := tm.Tokens[tokenID]; exists {
		return nil, errors.New("token ID already exists")
	}

	// Check if max supply is reached
	if tm.TotalSupply >= tm.MaxSupply {
		return nil, errors.New("maximum supply of tokens reached")
	}

	// Create the token
	token := &SYN721Token{
		TokenID:    tokenID,
		TokenURI:   tokenURI,
		Owner:      owner,
		Ledger:     tm.Ledger,
		Consensus:  tm.Consensus,
		Encryption: tm.Encryption,
	}

	// Encrypt token data
	encryptedData, err := tm.Encryption.EncryptData(fmt.Sprintf("%v", token), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting token data: %v", err)
	}

	// Validate the minting transaction through consensus
	if valid, err := tm.Consensus.ValidateTokenMint(tokenID, owner, encryptedData); !valid || err != nil {
		return nil, fmt.Errorf("token minting failed consensus validation: %v", err)
	}

	// Store the token in storage
	tm.Tokens[tokenID] = token
	tm.Storage.StoreTokenData(tokenID, token)

	// Record the minting transaction in the ledger
	err = tm.Ledger.RecordTokenMint(tokenID, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to log mint transaction in the ledger: %v", err)
	}

	// Increase total supply
	tm.TotalSupply++

	fmt.Printf("Token %s successfully minted for owner %s.\n", tokenID, owner)
	return token, nil
}

// Transfer transfers ownership of an SYN721 token from one address to another.
func (tm *SYN721TokenManager) Transfer(tokenID, from, to string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	token, exists := tm.Tokens[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	// Ensure the sender owns the token or is approved to transfer it
	if token.Owner != from && token.Approved != from {
		return errors.New("transfer not authorized")
	}

	// Perform the transfer
	token.Owner = to
	token.Approved = "" // Clear any approvals

	// Encrypt transfer transaction data
	encryptedData, err := tm.Encryption.EncryptData(fmt.Sprintf("Transfer %s from %s to %s", tokenID, from, to), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting transfer transaction: %v", err)
	}

	// Validate the transfer transaction using Synnergy Consensus
	if valid, err := tm.Consensus.ValidateTokenTransfer(tokenID, from, to); !valid || err != nil {
		return fmt.Errorf("token transfer failed consensus validation: %v", err)
	}

	// Update the token data in storage
	tm.Storage.UpdateTokenData(tokenID, token)

	// Record the transfer in the ledger
	err = tm.Ledger.RecordTokenTransfer(tokenID, from, to)
	if err != nil {
		return fmt.Errorf("failed to log transfer in the ledger: %v", err)
	}

	fmt.Printf("Token %s transferred from %s to %s.\n", tokenID, from, to)
	return nil
}

// Approve allows another address to transfer the token on the owner's behalf.
func (tm *SYN721TokenManager) Approve(tokenID, owner, approved string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	token, exists := tm.Tokens[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	// Ensure only the owner can approve transfers
	if token.Owner != owner {
		return errors.New("only the owner can approve transfers")
	}

	// Set the approved address
	token.Approved = approved

	// Encrypt approval transaction data
	encryptedData, err := tm.Encryption.EncryptData(fmt.Sprintf("Approved %s to transfer token %s", approved, tokenID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting approval transaction: %v", err)
	}

	// Validate the approval using Synnergy Consensus
	if valid, err := tm.Consensus.ValidateApproval(tokenID, owner, approved); !valid || err != nil {
		return fmt.Errorf("approval failed consensus validation: %v", err)
	}

	// Update token approval in storage
	tm.Storage.UpdateTokenData(tokenID, token)

	// Record the approval in the ledger
	err = tm.Ledger.RecordTokenApproval(tokenID, owner, approved)
	if err != nil {
		return fmt.Errorf("failed to log approval in the ledger: %v", err)
	}

	fmt.Printf("Token %s approved for transfer by %s to %s.\n", tokenID, owner, approved)
	return nil
}

// Burn destroys an SYN721 token and removes it from the ledger.
func (tm *SYN721TokenManager) Burn(tokenID, owner string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	token, exists := tm.Tokens[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	// Ensure only the owner can burn the token
	if token.Owner != owner {
		return errors.New("only the owner can burn the token")
	}

	// Encrypt burn transaction data
	encryptedData, err := tm.Encryption.EncryptData(fmt.Sprintf("Burn token %s owned by %s", tokenID, owner), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting burn transaction: %v", err)
	}

	// Validate the burn transaction using Synnergy Consensus
	if valid, err := tm.Consensus.ValidateTokenBurn(tokenID, owner); !valid || err != nil {
		return fmt.Errorf("token burn failed consensus validation: %v", err)
	}

	// Remove the token from storage and the token list
	tm.Storage.RemoveTokenData(tokenID)
	delete(tm.Tokens, tokenID)

	// Record the burn transaction in the ledger
	err = tm.Ledger.RecordTokenBurn(tokenID, owner)
	if err != nil {
		return fmt.Errorf("failed to log burn transaction in the ledger: %v", err)
	}

	// Decrease total supply
	tm.TotalSupply--

	fmt.Printf("Token %s burned by owner %s.\n", tokenID, owner)
	return nil
}

// GetToken retrieves the details of an SYN721 token.
func (tm *SYN721TokenManager) GetToken(tokenID string) (*SYN721Token, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	token, exists := tm.Tokens[tokenID]
	if !exists {
		return nil, errors.New("token not found")
	}

	return token, nil
}

// ListAllTokens returns a list of all SYN721 tokens managed by the token manager.
func (tm *SYN721TokenManager) ListAllTokens() map[string]*SYN721Token {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	return tm.Tokens
}
