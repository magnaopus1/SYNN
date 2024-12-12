package syn200

import (
	"errors"
	"sync"
	"time"
)

// IssueTokenTransaction initiates a transaction to issue a new SYN200 token.
func IssueTokenTransaction(metadata common.CarbonCreditMetadata, issuer common.IssuerRecord) (*common.SYN200Token, error) {
    token, err := CreateSYN200Token(metadata, issuer)
    if err != nil {
        return nil, fmt.Errorf("failed to create token: %v", err)
    }

    // Log issuance event and validate in Synnergy Consensus
    if err := LogTokenIssuance(token); err != nil {
        return nil, fmt.Errorf("failed to log issuance: %v", err)
    }
    if err := ValidateTransactionForConsensus(token.TokenID); err != nil {
        return nil, fmt.Errorf("sub-block validation failed for issuance: %v", err)
    }

    return token, nil
}

// TransferTokenTransaction transfers a SYN200 token to a new owner.
func TransferTokenTransaction(tokenID, newOwnerID, transferMethod string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    ownershipRecord := common.OwnershipRecord{
        OwnerID:        newOwnerID,
        OwnershipDate:  time.Now(),
        TransferMethod: transferMethod,
    }

    token.OwnershipHistory = append(token.OwnershipHistory, ownershipRecord)

    // Log transfer event and validate in Synnergy Consensus
    if err := LogTokenTransfer(tokenID, token.Issuer.Name, newOwnerID); err != nil {
        return fmt.Errorf("failed to log transfer: %v", err)
    }
    if err := ValidateTransactionForConsensus(tokenID); err != nil {
        return fmt.Errorf("sub-block validation failed for transfer: %v", err)
    }

    return nil
}

// RetireTokenTransaction retires a SYN200 token, marking it as inactive.
func RetireTokenTransaction(tokenID string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    if token.ValidityStatus == "Retired" {
        return fmt.Errorf("token is already retired")
    }

    token.ValidityStatus = "Retired"

    // Log retirement event and validate in Synnergy Consensus
    if err := RetireToken(tokenID); err != nil {
        return fmt.Errorf("failed to retire token: %v", err)
    }
    if err := ValidateTransactionForConsensus(tokenID); err != nil {
        return fmt.Errorf("sub-block validation failed for retirement: %v", err)
    }

    return nil
}

// ReverseTokenTransaction reverses a previously conducted transaction on a SYN200 token.
func ReverseTokenTransaction(tokenID, reason string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    reversalRequest := common.ReversalRequest{
        RequestID:    generateRequestID(),
        TokenID:      tokenID,
        Reason:       reason,
        RequestedAt:  time.Now(),
    }

    // Encrypt and log the reversal request in the ledger
    encryptedRequest, err := encryption.EncryptMetadata(reversalRequest)
    if err != nil {
        return fmt.Errorf("failed to encrypt reversal request: %v", err)
    }

    if err := ledger.RecordReversalRequest(encryptedRequest); err != nil {
        return fmt.Errorf("failed to log reversal request: %v", err)
    }

    // Perform the reversal and log the result
    reversalEvent := common.TokenEvent{
        EventID:        generateEventID(),
        TokenID:        tokenID,
        EventType:      "Reversal",
        EventTimestamp: time.Now(),
        Description:    fmt.Sprintf("Transaction on token %s reversed for reason: %s", tokenID, reason),
    }

    encryptedEvent, err := encryption.EncryptMetadata(reversalEvent)
    if err != nil {
        return fmt.Errorf("failed to encrypt reversal event: %v", err)
    }

    if err := ledger.RecordEvent(token, encryptedEvent); err != nil {
        return fmt.Errorf("failed to log reversal event: %v", err)
    }
    if err := ValidateTransactionForConsensus(tokenID); err != nil {
        return fmt.Errorf("sub-block validation failed for reversal: %v", err)
    }

    return nil
}

// ValidateTransactionForConsensus validates a transaction within a sub-block under Synnergy Consensus.
func ValidateTransactionForConsensus(tokenID string) error {
    validationStatus := ledger.ValidateInSubBlock(tokenID)
    if !validationStatus {
        return fmt.Errorf("validation in Synnergy Consensus failed for token ID %s", tokenID)
    }

    // Once validated in the sub-block, the transaction is ready for block validation
    return nil
}

// generateEventID creates a unique identifier for events.
func generateEventID() string {
    b := make([]byte, 12)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}

// generateRequestID creates a unique identifier for reversal requests.
func generateRequestID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
