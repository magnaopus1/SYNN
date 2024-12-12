package syn200

import (
	"errors"
	"sync"
	"time"
)

// CreateBatchTokens creates multiple SYN200 tokens in a batch, logging each issuance in the ledger.
func CreateBatchTokens(batchMetadata []common.CarbonCreditMetadata, issuer common.IssuerRecord) ([]*common.SYN200Token, error) {
    tokens := []*common.SYN200Token{}
    for _, metadata := range batchMetadata {
        token, err := CreateSYN200Token(metadata, issuer)
        if err != nil {
            return nil, fmt.Errorf("failed to create token in batch: %v", err)
        }
        
        // Log the issuance for each token in the ledger
        if err := LogTokenIssuance(token); err != nil {
            return nil, fmt.Errorf("failed to log issuance in batch: %v", err)
        }
        tokens = append(tokens, token)
    }
    return tokens, nil
}

// UpdateTokenDetails updates specific fields of a SYN200 token, encrypts, and logs the update.
func UpdateTokenDetails(tokenID string, updatedFields map[string]string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    // Update token fields based on the provided map
    for field, newValue := range updatedFields {
        switch field {
        case "ValidityStatus":
            token.ValidityStatus = newValue
        case "Issuer":
            token.Issuer.Name = newValue
        case "CO2OffsetAmount":
            token.CreditMetadata.CO2OffsetAmount = parseCO2Offset(newValue)
        // Additional fields can be added as necessary
        default:
            return fmt.Errorf("unknown field: %s", field)
        }
    }

    // Encrypt and log the update event in the ledger
    return LogTokenUpdate(tokenID, updatedFields)
}

// RetireToken retires a SYN200 token by marking it as inactive and recording the retirement.
func RetireToken(tokenID string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    if token.ValidityStatus == "Retired" {
        return fmt.Errorf("token is already retired")
    }

    token.ValidityStatus = "Retired"

    retirementEvent := common.TokenEvent{
        EventID:        generateEventID(),
        TokenID:        tokenID,
        EventType:      "Retirement",
        EventTimestamp: time.Now(),
        Description:    fmt.Sprintf("Token %s has been retired", tokenID),
    }

    // Encrypt and record the retirement in the ledger
    encryptedEvent, err := encryption.EncryptMetadata(retirementEvent)
    if err != nil {
        return fmt.Errorf("failed to encrypt retirement event: %v", err)
    }

    if err := ledger.RecordEvent(token, encryptedEvent); err != nil {
        return fmt.Errorf("failed to record retirement in ledger: %v", err)
    }

    return nil
}

// MonitorTokenStatus continuously checks the status of a SYN200 token and triggers expiration if due.
func MonitorTokenStatus(tokenID string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    if token.ExpirationDate != nil && time.Now().After(*token.ExpirationDate) {
        return LogTokenExpiration(tokenID)
    }

    return nil
}

// TransferOwnership transfers a SYN200 token to a new owner and logs the transfer in the ledger.
func TransferOwnership(tokenID, newOwnerID, transferMethod string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    ownershipRecord := common.OwnershipRecord{
        OwnerID:        newOwnerID,
        OwnershipDate:  time.Now(),
        TransferMethod: transferMethod,
    }

    // Append the new ownership record
    token.OwnershipHistory = append(token.OwnershipHistory, ownershipRecord)

    // Log the transfer event
    return LogTokenTransfer(tokenID, token.Issuer.Name, newOwnerID)
}

// generateEventID creates a unique identifier for events.
func generateEventID() string {
    b := make([]byte, 12)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}

// parseCO2Offset converts a CO2 offset string into a float64.
func parseCO2Offset(offset string) float64 {
    value, err := strconv.ParseFloat(offset, 64)
    if err != nil {
        panic(fmt.Sprintf("failed to parse CO2 offset: %v", err))
    }
    return value
}
