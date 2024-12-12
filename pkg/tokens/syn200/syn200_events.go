package syn200

import (
	"errors"
	"sync"
	"time"
)

// LogTokenIssuance logs the issuance event of a new SYN200 token to the ledger.
func LogTokenIssuance(token *common.SYN200Token) error {
    issuanceEvent := common.TokenEvent{
        EventID:        generateEventID(),
        TokenID:        token.TokenID,
        EventType:      "Issuance",
        EventTimestamp: time.Now(),
        Description:    fmt.Sprintf("New SYN200 token issued with ID %s", token.TokenID),
    }

    // Encrypt the issuance event for secure logging
    encryptedEvent, err := encryption.EncryptMetadata(issuanceEvent)
    if err != nil {
        return fmt.Errorf("failed to encrypt issuance event: %v", err)
    }

    // Record event in the ledger
    if err := ledger.RecordEvent(token, encryptedEvent); err != nil {
        return fmt.Errorf("failed to log issuance event in ledger: %v", err)
    }

    return nil
}

// LogTokenTransfer logs the transfer event of a SYN200 token between parties in the ledger.
func LogTokenTransfer(tokenID string, fromOwner string, toOwner string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    transferEvent := common.TokenEvent{
        EventID:        generateEventID(),
        TokenID:        tokenID,
        EventType:      "Transfer",
        EventTimestamp: time.Now(),
        Description:    fmt.Sprintf("Token %s transferred from %s to %s", tokenID, fromOwner, toOwner),
    }

    // Encrypt the transfer event for secure logging
    encryptedEvent, err := encryption.EncryptMetadata(transferEvent)
    if err != nil {
        return fmt.Errorf("failed to encrypt transfer event: %v", err)
    }

    // Record event in the ledger and validate in Synnergy Consensus
    if err := ledger.RecordEvent(token, encryptedEvent); err != nil {
        return fmt.Errorf("failed to log transfer event in ledger: %v", err)
    }
    if err := ValidateTokenForSubBlock(tokenID); err != nil {
        return fmt.Errorf("sub-block validation for transfer failed: %v", err)
    }

    return nil
}

// LogTokenExpiration logs the expiration event for a SYN200 token, marking it as expired.
func LogTokenExpiration(tokenID string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    if token.ExpirationDate != nil && time.Now().After(*token.ExpirationDate) {
        token.ValidityStatus = "Expired"
        
        expirationEvent := common.TokenEvent{
            EventID:        generateEventID(),
            TokenID:        tokenID,
            EventType:      "Expiration",
            EventTimestamp: time.Now(),
            Description:    fmt.Sprintf("Token %s has expired", tokenID),
        }

        // Encrypt and record the expiration event
        encryptedEvent, err := encryption.EncryptMetadata(expirationEvent)
        if err != nil {
            return fmt.Errorf("failed to encrypt expiration event: %v", err)
        }

        if err := ledger.RecordEvent(token, encryptedEvent); err != nil {
            return fmt.Errorf("failed to log expiration event in ledger: %v", err)
        }

        return nil
    }

    return fmt.Errorf("token %s has not reached its expiration date", tokenID)
}

// LogTokenUpdate logs any updates made to a SYN200 token's data, ensuring secure event recording.
func LogTokenUpdate(tokenID string, updatedFields map[string]string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    updateDetails := ""
    for field, value := range updatedFields {
        updateDetails += fmt.Sprintf("%s: %s; ", field, value)
    }

    updateEvent := common.TokenEvent{
        EventID:        generateEventID(),
        TokenID:        tokenID,
        EventType:      "Update",
        EventTimestamp: time.Now(),
        Description:    fmt.Sprintf("Token %s updated with fields: %s", tokenID, updateDetails),
    }

    // Encrypt and record the update event
    encryptedEvent, err := encryption.EncryptMetadata(updateEvent)
    if err != nil {
        return fmt.Errorf("failed to encrypt update event: %v", err)
    }

    if err := ledger.RecordEvent(token, encryptedEvent); err != nil {
        return fmt.Errorf("failed to log update event in ledger: %v", err)
    }

    return nil
}

// generateEventID creates a unique identifier for token events.
func generateEventID() string {
    b := make([]byte, 12)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
