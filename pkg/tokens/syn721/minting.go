package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// SYN721MintingManager handles the minting operations for SYN721 tokens.
type SYN721MintingManager struct {
	mutex      sync.Mutex                 // For thread safety
	Ledger     *ledger.Ledger             // Reference to the ledger for recording transactions
	Consensus  *synnergy_consensus.Engine // Synnergy Consensus engine for validation
	Encryption *encryption.Encryption     // Encryption service for securing transaction data
	Tokens     map[string]*SYN721Token    // Map of all SYN721 tokens by their ID
	TotalMinted uint64                    // Total number of minted tokens (NFTs)
	MaxSupply  uint64                     // Maximum supply of tokens that can be minted
}

// NewSYN721MintingManager initializes a new minting manager for SYN721 tokens.
func NewSYN721MintingManager(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, maxSupply uint64) *SYN721MintingManager {
	return &SYN721MintingManager{
		Ledger:      ledgerInstance,
		Consensus:   consensusEngine,
		Encryption:  encryptionService,
		Tokens:      make(map[string]*SYN721Token),
		MaxSupply:   maxSupply,
		TotalMinted: 0,
	}
}

// MintToken mints a new SYN721 token (NFT) and records it in the ledger.
func (mm *SYN721MintingManager) MintToken(tokenID, owner, tokenURI string) (*SYN721Token, error) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	// Check if the maximum supply is reached
	if mm.TotalMinted >= mm.MaxSupply {
		return nil, errors.New("maximum supply reached")
	}

	// Check if the token ID already exists
	if _, exists := mm.Tokens[tokenID]; exists {
		return nil, errors.New("token ID already exists")
	}

	// Create the new token
	token := &SYN721Token{
		TokenID:    tokenID,
		TokenURI:   tokenURI,
		Owner:      owner,
		Ledger:     mm.Ledger,
		Consensus:  mm.Consensus,
		Encryption: mm.Encryption,
	}

	// Encrypt token data
	encryptedData, err := mm.Encryption.EncryptData(fmt.Sprintf("%v", token), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting token data: %v", err)
	}

	// Validate the minting transaction using Synnergy Consensus
	if valid, err := mm.Consensus.ValidateTokenMint(tokenID, owner, encryptedData); !valid || err != nil {
		return nil, fmt.Errorf("token minting failed consensus validation: %v", err)
	}

	// Add the token to the token map
	mm.Tokens[tokenID] = token
	mm.TotalMinted++

	// Record the minting transaction in the ledger
	err = mm.Ledger.RecordTokenMint(tokenID, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to record mint transaction in ledger: %v", err)
	}

	fmt.Printf("Token %s successfully minted for owner %s.\n", tokenID, owner)
	return token, nil
}

// BatchMint allows minting multiple SYN721 tokens (NFTs) at once.
func (mm *SYN721MintingManager) BatchMint(tokens []TokenMintingRequest) ([]*SYN721Token, error) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	if len(tokens) == 0 {
		return nil, errors.New("no tokens to mint")
	}

	mintedTokens := make([]*SYN721Token, 0)

	for _, request := range tokens {
		// Check if the maximum supply is reached
		if mm.TotalMinted >= mm.MaxSupply {
			return mintedTokens, errors.New("maximum supply reached during batch minting")
		}

		// Check if the token ID already exists
		if _, exists := mm.Tokens[request.TokenID]; exists {
			return mintedTokens, fmt.Errorf("token ID %s already exists", request.TokenID)
		}

		// Create the new token
		token := &SYN721Token{
			TokenID:    request.TokenID,
			TokenURI:   request.TokenURI,
			Owner:      request.Owner,
			Ledger:     mm.Ledger,
			Consensus:  mm.Consensus,
			Encryption: mm.Encryption,
		}

		// Encrypt token data
		encryptedData, err := mm.Encryption.EncryptData(fmt.Sprintf("%v", token), common.EncryptionKey)
		if err != nil {
			return mintedTokens, fmt.Errorf("error encrypting token data for token %s: %v", request.TokenID, err)
		}

		// Validate the minting transaction using Synnergy Consensus
		if valid, err := mm.Consensus.ValidateTokenMint(request.TokenID, request.Owner, encryptedData); !valid || err != nil {
			return mintedTokens, fmt.Errorf("token minting failed consensus validation for token %s: %v", request.TokenID, err)
		}

		// Add the token to the token map
		mm.Tokens[request.TokenID] = token
		mm.TotalMinted++

		// Record the minting transaction in the ledger
		err = mm.Ledger.RecordTokenMint(request.TokenID, request.Owner)
		if err != nil {
			return mintedTokens, fmt.Errorf("failed to record mint transaction for token %s: %v", request.TokenID, err)
		}

		mintedTokens = append(mintedTokens, token)
	}

	fmt.Printf("Batch minting of %d tokens completed successfully.\n", len(mintedTokens))
	return mintedTokens, nil
}

// TokenMintingRequest represents the details required for minting a new SYN721 token.
type TokenMintingRequest struct {
	TokenID   string // Unique identifier for the token (NFT)
	Owner     string // Address of the token owner
	TokenURI  string // URI linking to the token's metadata (image, etc.)
}

// GetTotalMinted returns the total number of minted SYN721 tokens.
func (mm *SYN721MintingManager) GetTotalMinted() uint64 {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	return mm.TotalMinted
}

// GetMaxSupply returns the maximum supply of SYN721 tokens that can be minted.
func (mm *SYN721MintingManager) GetMaxSupply() uint64 {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	return mm.MaxSupply
}
