package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// SYN721Burner handles the burning (destruction) of SYN721 tokens (NFTs).
type SYN721Burner struct {
	mutex      sync.Mutex                 // For thread-safe operations
	Ledger     *ledger.Ledger             // Reference to the ledger for recording token burns
	Consensus  *synnergy_consensus.Engine // Synnergy Consensus engine for validation
	Encryption *encryption.Encryption     // Encryption service for securing burn data
	Storage    *SYN721Storage             // Storage for SYN721 token data
}

// NewSYN721Burner initializes a new burner for SYN721 tokens.
func NewSYN721Burner(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, storage *SYN721Storage) *SYN721Burner {
	return &SYN721Burner{
		Ledger:     ledgerInstance,
		Consensus:  consensusEngine,
		Encryption: encryptionService,
		Storage:    storage,
	}
}

// BurnToken burns (destroys) an SYN721 token and removes it from the ledger.
func (sb *SYN721Burner) BurnToken(tokenID, owner string) error {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()

	// Retrieve the token from storage
	token, exists := sb.Storage.GetTokenData(tokenID)
	if !exists {
		return errors.New("token not found")
	}

	// Ensure the caller is the owner of the token
	if token.Owner != owner {
		return errors.New("only the owner can burn the token")
	}

	// Encrypt burn transaction data for security
	encryptedData, err := sb.Encryption.EncryptData(fmt.Sprintf("Burn token %s owned by %s", tokenID, owner), "")
	if err != nil {
		return fmt.Errorf("error encrypting burn transaction: %v", err)
	}

	// Validate the burn operation using Synnergy Consensus
	if valid, err := sb.Consensus.ValidateTokenBurn(tokenID, owner); !valid || err != nil {
		return fmt.Errorf("token burn failed consensus validation: %v", err)
	}

	// Remove the token from storage
	err = sb.Storage.RemoveTokenData(tokenID)
	if err != nil {
		return fmt.Errorf("failed to remove token from storage: %v", err)
	}

	// Record the burn operation in the ledger
	err = sb.Ledger.RecordTokenBurn(tokenID, owner)
	if err != nil {
		return fmt.Errorf("failed to record token burn in ledger: %v", err)
	}

	fmt.Printf("Token %s successfully burned by owner %s.\n", tokenID, owner)
	return nil
}

// BurnMultipleTokens burns multiple SYN721 tokens and removes them from the ledger.
func (sb *SYN721Burner) BurnMultipleTokens(tokenIDs []string, owner string) error {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()

	for _, tokenID := range tokenIDs {
		// Retrieve the token from storage
		token, exists := sb.Storage.GetTokenData(tokenID)
		if !exists {
			return fmt.Errorf("token %s not found", tokenID)
		}

		// Ensure the caller is the owner of the token
		if token.Owner != owner {
			return fmt.Errorf("only the owner can burn the token %s", tokenID)
		}

		// Encrypt burn transaction data for security
		encryptedData, err := sb.Encryption.EncryptData(fmt.Sprintf("Burn token %s owned by %s", tokenID, owner), "")
		if err != nil {
			return fmt.Errorf("error encrypting burn transaction for token %s: %v", tokenID, err)
		}

		// Validate the burn operation using Synnergy Consensus
		if valid, err := sb.Consensus.ValidateTokenBurn(tokenID, owner); !valid || err != nil {
			return fmt.Errorf("token %s burn failed consensus validation: %v", tokenID, err)
		}

		// Remove the token from storage
		err = sb.Storage.RemoveTokenData(tokenID)
		if err != nil {
			return fmt.Errorf("failed to remove token %s from storage: %v", tokenID, err)
		}

		// Record the burn operation in the ledger
		err = sb.Ledger.RecordTokenBurn(tokenID, owner)
		if err != nil {
			return fmt.Errorf("failed to record token %s burn in ledger: %v", tokenID, err)
		}

		fmt.Printf("Token %s successfully burned by owner %s.\n", tokenID, owner)
	}

	return nil
}

// ValidateBurn checks if a token is eligible to be burned based on consensus rules.
func (sb *SYN721Burner) ValidateBurn(tokenID, owner string) error {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()

	// Retrieve the token from storage
	token, exists := sb.Storage.GetTokenData(tokenID)
	if !exists {
		return errors.New("token not found")
	}

	// Ensure the caller is the owner of the token
	if token.Owner != owner {
		return errors.New("only the owner can burn the token")
	}

	// Validate the burn transaction using Synnergy Consensus
	valid, err := sb.Consensus.ValidateTokenBurn(tokenID, owner)
	if !valid || err != nil {
		return fmt.Errorf("token burn failed consensus validation: %v", err)
	}

	fmt.Printf("Token %s burn validated successfully.\n", tokenID)
	return nil
}
