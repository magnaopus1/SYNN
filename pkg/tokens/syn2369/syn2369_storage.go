package syn2369

import (
	"time"
	"errors"
)


// StoreToken stores a SYN2369Token in the blockchain ledger after validating and encrypting the data.
func StoreToken(token common.SYN2369Token) error {
	// Validate token through Synnergy Consensus
	err := synnergy.ValidateTokenStorage(token)
	if err != nil {
		return errors.New("token validation failed: " + err.Error())
	}

	// Encrypt sensitive token data before storage (e.g., metadata)
	encryptedMetadata, err := encryption.EncryptData(token.Metadata)
	if err != nil {
		return errors.New("failed to encrypt token metadata: " + err.Error())
	}
	token.Metadata = encryptedMetadata

	// Store the token in the ledger
	err = ledger.StoreToken(token)
	if err != nil {
		return errors.New("failed to store token in ledger: " + err.Error())
	}

	// Log the event of token storage
	err = LogStorageEvent(token, "TokenStored", "Token ID "+token.TokenID+" stored in ledger.")
	if err != nil {
		return err
	}

	return nil
}

// RetrieveToken retrieves a SYN2369Token from the ledger and decrypts the sensitive data.
func RetrieveToken(tokenID string) (common.SYN2369Token, error) {
	// Retrieve the token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2369Token{}, errors.New("token not found: " + err.Error())
	}

	// Decrypt the token's metadata
	decryptedMetadata, err := encryption.DecryptData(token.Metadata)
	if err != nil {
		return common.SYN2369Token{}, errors.New("failed to decrypt token metadata: " + err.Error())
	}
	token.Metadata = decryptedMetadata

	return token, nil
}

// UpdateToken updates an existing SYN2369Token in the blockchain ledger after re-encrypting any sensitive data.
func UpdateToken(token common.SYN2369Token) error {
	// Validate the update through Synnergy Consensus
	err := synnergy.ValidateTokenUpdate(token)
	if err != nil {
		return errors.New("token update validation failed: " + err.Error())
	}

	// Encrypt sensitive token data before updating in ledger
	encryptedMetadata, err := encryption.EncryptData(token.Metadata)
	if err != nil {
		return errors.New("failed to encrypt token metadata: " + err.Error())
	}
	token.Metadata = encryptedMetadata

	// Update the token in the ledger
	err = ledger.UpdateToken(token)
	if err != nil {
		return errors.New("failed to update token in ledger: " + err.Error())
	}

	// Log the event of token update
	err = LogStorageEvent(token, "TokenUpdated", "Token ID "+token.TokenID+" updated in ledger.")
	if err != nil {
		return err
	}

	return nil
}

// DeleteToken removes a SYN2369Token from the blockchain ledger.
func DeleteToken(tokenID string) error {
	// Retrieve the token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("token not found: " + err.Error())
	}

	// Perform any pre-deletion checks
	err = synnergy.ValidateTokenDeletion(token)
	if err != nil {
		return errors.New("token deletion validation failed: " + err.Error())
	}

	// Delete the token from the ledger
	err = ledger.DeleteToken(tokenID)
	if err != nil {
		return errors.New("failed to delete token from ledger: " + err.Error())
	}

	// Log the event of token deletion
	err = LogStorageEvent(token, "TokenDeleted", "Token ID "+tokenID+" deleted from ledger.")
	if err != nil {
		return err
	}

	return nil
}

// ListTokens retrieves a list of all SYN2369Tokens currently stored in the ledger.
func ListTokens() ([]common.SYN2369Token, error) {
	// Retrieve the list of tokens from the ledger
	tokens, err := ledger.ListTokensByType("SYN2369")
	if err != nil {
		return nil, errors.New("failed to retrieve tokens: " + err.Error())
	}

	// Decrypt metadata for each token
	for i := range tokens {
		decryptedMetadata, err := encryption.DecryptData(tokens[i].Metadata)
		if err != nil {
			return nil, errors.New("failed to decrypt token metadata: " + err.Error())
		}
		tokens[i].Metadata = decryptedMetadata
	}

	return tokens, nil
}

// LogStorageEvent logs the storage-related event for the SYN2369 token in the ledger.
func LogStorageEvent(token common.SYN2369Token, eventType, eventDescription string) error {
	eventLog := common.SYN2369Event{
		TokenID:          token.TokenID,
		EventType:        eventType,
		EventDescription: eventDescription,
		EventTime:        time.Now(),
	}

	// Store the event in the ledger
	err := ledger.StoreEvent(eventLog)
	if err != nil {
		return err
	}

	return nil
}

// GetTokenHistory retrieves the transaction history and event logs for a specific SYN2369 token.
func GetTokenHistory(tokenID string) ([]common.SYN2369Event, error) {
	// Retrieve the token's event logs from the ledger
	events, err := ledger.GetTokenEventLogs(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve event logs: " + err.Error())
	}

	return events, nil
}
