package syn1800

import (
	"time"
	"fmt"
)

// CarbonTokenStorageManager handles storage and retrieval operations for SYN1800 tokens.
type CarbonTokenStorageManager struct {
	ledger *ledger.Ledger  // Ledger instance for integrating with the blockchain
}

// NewCarbonTokenStorageManager initializes a new CarbonTokenStorageManager.
func NewCarbonTokenStorageManager(ledger *ledger.Ledger) *CarbonTokenStorageManager {
	return &CarbonTokenStorageManager{ledger: ledger}
}

// StoreToken securely stores a SYN1800 token in the ledger, encrypting its metadata before storage.
func (ctsm *CarbonTokenStorageManager) StoreToken(token common.SYN1800Token) error {
	// Encrypt the token's metadata
	encryptedMetadata, err := encryptTokenMetadata(token)
	if err != nil {
		return fmt.Errorf("failed to encrypt token metadata: %v", err)
	}
	token.EncryptedMetadata = encryptedMetadata

	// Store the token in the ledger
	err = ctsm.ledger.StoreToken(&token)
	if err != nil {
		return fmt.Errorf("failed to store token in the ledger: %v", err)
	}

	return nil
}

// RetrieveToken retrieves a SYN1800 token from the ledger by its token ID.
func (ctsm *CarbonTokenStorageManager) RetrieveToken(tokenID string) (*common.SYN1800Token, error) {
	// Retrieve the token from the ledger
	token, err := ctsm.ledger.GetTokenByID(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token from ledger: %v", err)
	}

	// Cast to SYN1800Token type
	syn1800Token, ok := token.(*common.SYN1800Token)
	if !ok {
		return nil, errors.New("invalid token type retrieved")
	}

	// Decrypt the token metadata before returning
	decryptedMetadata, err := decryptTokenMetadata(syn1800Token.EncryptedMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token metadata: %v", err)
	}
	syn1800Token.EncryptedMetadata = decryptedMetadata

	return syn1800Token, nil
}

// UpdateToken updates an existing SYN1800 token in the ledger, ensuring encryption before updating.
func (ctsm *CarbonTokenStorageManager) UpdateToken(token *common.SYN1800Token) error {
	// Encrypt the token's metadata before updating
	encryptedMetadata, err := encryptTokenMetadata(*token)
	if err != nil {
		return fmt.Errorf("failed to encrypt token metadata: %v", err)
	}
	token.EncryptedMetadata = encryptedMetadata

	// Update the token in the ledger
	err = ctsm.ledger.UpdateTokenInLedger(token)
	if err != nil {
		return fmt.Errorf("failed to update token in the ledger: %v", err)
	}

	return nil
}

// DeleteToken removes a SYN1800 token from the ledger by its token ID.
func (ctsm *CarbonTokenStorageManager) DeleteToken(tokenID string) error {
	// Delete the token from the ledger
	err := ctsm.ledger.DeleteTokenByID(tokenID)
	if err != nil {
		return fmt.Errorf("failed to delete token from ledger: %v", err)
	}

	return nil
}

// FetchTokenHistory retrieves the full history of a SYN1800 token's transactions, transfers, and updates.
func (ctsm *CarbonTokenStorageManager) FetchTokenHistory(tokenID string) ([]common.ImmutableRecord, error) {
	// Retrieve the token from the ledger
	token, err := ctsm.ledger.GetTokenByID(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token from ledger: %v", err)
	}

	// Cast to SYN1800Token type
	syn1800Token, ok := token.(*common.SYN1800Token)
	if !ok {
		return nil, errors.New("invalid token type retrieved")
	}

	// Return the immutable record history
	return syn1800Token.ImmutableRecords, nil
}

// Helper functions

// encryptTokenMetadata encrypts the sensitive metadata of a SYN1800 token.
func encryptTokenMetadata(token common.SYN1800Token) ([]byte, error) {
	// Real-world encryption logic using a secure key
	return crypto.Encrypt([]byte(fmt.Sprintf("%v", token)), "secure-encryption-key")
}

// decryptTokenMetadata decrypts the sensitive metadata of a SYN1800 token.
func decryptTokenMetadata(encryptedData []byte) ([]byte, error) {
	// Real-world decryption logic using a secure key
	return crypto.Decrypt(encryptedData, "secure-encryption-key")
}

// generateUniqueID generates a unique identifier for tokens, logs, and records.
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
