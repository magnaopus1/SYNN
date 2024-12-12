package syn2369

import (
	"time"
	"errors"
)


// CreateSYN2369Token creates a new SYN2369 token and stores it in the ledger.
func CreateSYN2369Token(token common.SYN2369Token) error {
	// Validate the token through Synnergy Consensus
	err := synnergy.ValidateNewToken(token)
	if err != nil {
		return errors.New("token validation failed: " + err.Error())
	}

	// Encrypt sensitive metadata before storing
	if err := encryptTokenMetadata(&token); err != nil {
		return errors.New("metadata encryption failed: " + err.Error())
	}

	// Store the new token in the ledger
	if err := ledger.StoreToken(token); err != nil {
		return errors.New("storing token in ledger failed: " + err.Error())
	}

	// Log token creation event
	if err := LogEvent(token, "Creation", "New SYN2369 token created with ID: "+token.TokenID); err != nil {
		return err
	}

	return nil
}

// UpdateSYN2369Token allows updates to existing token metadata or attributes.
func UpdateSYN2369Token(tokenID string, updates map[string]interface{}) error {
	// Retrieve the token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("token not found: " + err.Error())
	}

	// Apply the updates to the token's attributes or metadata
	for key, value := range updates {
		if err := updateTokenAttribute(&token, key, value); err != nil {
			return err
		}
	}

	// Encrypt updated metadata
	if err := encryptTokenMetadata(&token); err != nil {
		return errors.New("metadata encryption failed: " + err.Error())
	}

	// Update the token in the ledger
	if err := ledger.UpdateToken(token); err != nil {
		return errors.New("updating token in ledger failed: " + err.Error())
	}

	// Log the update event
	if err := LogEvent(token, "Update", "Token with ID "+tokenID+" updated"); err != nil {
		return err
	}

	return nil
}

// TransferSYN2369Token transfers ownership of the token from one user to another.
func TransferSYN2369Token(tokenID, newOwner string) error {
	// Retrieve the token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("token not found: " + err.Error())
	}

	// Check if the token allows transfers
	if token.RestrictedTransfers {
		return errors.New("transfers restricted for this token")
	}

	// Validate the transfer through Synnergy Consensus
	err = synnergy.ValidateTransfer(token, token.Owner, newOwner)
	if err != nil {
		return errors.New("transfer validation failed: " + err.Error())
	}

	// Update the owner of the token
	token.Owner = newOwner

	// Update the token in the ledger
	if err := ledger.UpdateToken(token); err != nil {
		return errors.New("updating token owner in ledger failed: " + err.Error())
	}

	// Log the transfer event
	if err := GenerateTokenTransferEvent(token, token.Owner, newOwner); err != nil {
		return err
	}

	return nil
}

// BurnSYN2369Token removes a token permanently from circulation.
func BurnSYN2369Token(tokenID string) error {
	// Retrieve the token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("token not found: " + err.Error())
	}

	// Perform any necessary compliance checks before burning
	err = synnergy.ValidateTokenBurn(token)
	if err != nil {
		return errors.New("burn validation failed: " + err.Error())
	}

	// Remove the token from the ledger
	err = ledger.RemoveToken(tokenID)
	if err != nil {
		return errors.New("failed to remove token from ledger: " + err.Error())
	}

	// Log the burn event
	if err := LogEvent(token, "Burn", "Token with ID "+tokenID+" burned"); err != nil {
		return err
	}

	return nil
}

// encryptTokenMetadata encrypts sensitive metadata for the SYN2369 token.
func encryptTokenMetadata(token *common.SYN2369Token) error {
	encryptedData, err := encryption.EncryptData(token.Metadata)
	if err != nil {
		return err
	}
	token.Metadata = encryptedData
	return nil
}

// updateTokenAttribute updates a specific attribute of the token.
func updateTokenAttribute(token *common.SYN2369Token, key string, value interface{}) error {
	switch key {
	case "Name":
		token.Name, _ = value.(string)
	case "Attributes":
		if attributes, ok := value.(map[string]interface{}); ok {
			token.Attributes = attributes
		} else {
			return errors.New("invalid attributes value")
		}
	case "Metadata":
		if metadata, ok := value.(map[string]string); ok {
			token.Metadata = metadata
		} else {
			return errors.New("invalid metadata value")
		}
	default:
		return errors.New("unknown token attribute: " + key)
	}
	return nil
}

// ViewToken retrieves the details of a token by its ID.
func ViewToken(tokenID string) (common.SYN2369Token, error) {
	// Fetch the token from the ledger
	token, err := ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2369Token{}, err
	}

	// Decrypt the metadata for viewing
	decryptedMetadata, err := encryption.DecryptData(token.Metadata)
	if err != nil {
		return common.SYN2369Token{}, errors.New("failed to decrypt token metadata")
	}
	token.Metadata = decryptedMetadata

	return token, nil
}

// GetTokenHistory returns the event history of the SYN2369 token.
func GetTokenHistory(tokenID string) ([]SYN2369Event, error) {
	// Retrieve the event history from the ledger
	events, err := FetchTokenEventHistory(tokenID)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// GenerateAuditLog generates an audit log for SYN2369Token compliance and interactions.
func GenerateAuditLog(token common.SYN2369Token, auditDetails string) error {
	auditLog := common.AuditLog{
		AuditID:      generateAuditID(),
		TokenID:      token.TokenID,
		PerformedBy:  "Auditor",
		AuditDate:    time.Now(),
		AuditDetails: auditDetails,
	}

	// Store the audit log in the ledger
	err := ledger.AddAuditLog(token.TokenID, auditLog)
	if err != nil {
		return err
	}

	// Log audit event in the token's history
	return LogEvent(token, "Audit", "Audit performed: "+auditDetails)
}

// generateAuditID generates a unique audit ID.
func generateAuditID() string {
	return encryption.GenerateRandomID()
}
