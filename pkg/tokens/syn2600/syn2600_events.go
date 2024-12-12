package syn2600

import (

	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

)

// InvestorTokenEvent represents an event related to SYN2600 tokens, such as transfers, dividends, and returns.
type InvestorTokenEvent struct {
	EventID      string
	TokenID      string
	EventType    string // Possible values: "TRANSFER", "DIVIDEND_PAYMENT", "RETURN_UPDATE"
	Timestamp    time.Time
	Details      string
	Owner        string
	Signature    string // Encrypted event signature to verify authenticity
	Validated    bool
}

// RecordEvent stores and processes any event related to SYN2600 tokens.
func RecordEvent(tokenID string, eventType string, details string, owner string) (string, error) {
	// Fetch token details from the ledger
	token, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return "", errors.New("failed to fetch the token from the ledger for event recording")
	}

	// Decrypt token data before event handling
	decryptedToken, err := encryption.DecryptTokenData(token)
	if err != nil {
		return "", errors.New("failed to decrypt token data for event recording")
	}

	// Create a unique EventID
	eventID := generateEventID(tokenID, eventType)

	// Create the event
	event := InvestorTokenEvent{
		EventID:   eventID,
		TokenID:   tokenID,
		EventType: eventType,
		Timestamp: time.Now(),
		Details:   details,
		Owner:     owner,
		Signature: generateEventSignature(eventID, details, owner),
		Validated: false,
	}

	// Encrypt event data before storing
	encryptedEvent, err := encryption.EncryptEventData(&event)
	if err != nil {
		return "", errors.New("failed to encrypt event data before storing")
	}

	// Store the encrypted event in the ledger
	err = ledger.StoreInvestorTokenEvent(encryptedEvent)
	if err != nil {
		return "", errors.New("failed to store the event in the ledger")
	}

	// Validate event through Synnergy Consensus and update validation status
	err = synconsensus.ValidateSubBlock(eventID)
	if err != nil {
		return "", errors.New("event validation failed in Synnergy Consensus")
	}
	event.Validated = true

	// Update the event's validated status in the ledger
	encryptedEvent, err = encryption.EncryptEventData(&event)
	if err != nil {
		return "", errors.New("failed to re-encrypt event after validation")
	}
	err = ledger.UpdateInvestorTokenEvent(encryptedEvent)
	if err != nil {
		return "", errors.New("failed to update event validation status in the ledger")
	}

	return eventID, nil
}

// FetchEvents retrieves all events related to a specific investor token.
func FetchEvents(tokenID string) ([]InvestorTokenEvent, error) {
	// Fetch all related events from the ledger
	encryptedEvents, err := ledger.FetchAllEventsByTokenID(tokenID)
	if err != nil {
		return nil, errors.New("failed to fetch events from the ledger")
	}

	// Decrypt each event
	var events []InvestorTokenEvent
	for _, encryptedEvent := range encryptedEvents {
		decryptedEvent, err := encryption.DecryptEventData(encryptedEvent)
		if err != nil {
			return nil, errors.New("failed to decrypt event data")
		}
		events = append(events, *decryptedEvent)
	}

	return events, nil
}

// generateEventID generates a unique ID for each event.
func generateEventID(tokenID string, eventType string) string {
	hashInput := tokenID + eventType + time.Now().String()
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

// generateEventSignature creates a signature to ensure event authenticity.
func generateEventSignature(eventID string, details string, owner string) string {
	signatureInput := eventID + details + owner
	hash := sha256.Sum256([]byte(signatureInput))
	return hex.EncodeToString(hash[:])
}

// HandleTransferEvent processes token transfers and records a transfer event.
func HandleTransferEvent(tokenID string, newOwner string) (string, error) {
	// Update token ownership
	token, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return "", errors.New("failed to fetch token for transfer event")
	}

	// Decrypt token before transferring ownership
	decryptedToken, err := encryption.DecryptTokenData(token)
	if err != nil {
		return "", errors.New("failed to decrypt token data for ownership transfer")
	}

	decryptedToken.Owner = newOwner

	// Re-encrypt and update the token in the ledger
	encryptedToken, err := encryption.EncryptTokenData(decryptedToken)
	if err != nil {
		return "", errors.New("failed to encrypt token data after ownership transfer")
	}
	err = ledger.UpdateInvestorToken(encryptedToken)
	if err != nil {
		return "", errors.New("failed to update the ledger with new ownership")
	}

	// Record the transfer event
	eventID, err := RecordEvent(tokenID, "TRANSFER", "Ownership transferred to "+newOwner, newOwner)
	if err != nil {
		return "", err
	}

	return eventID, nil
}

// HandleDividendPayment processes dividends and records a dividend payment event.
func HandleDividendPayment(tokenID string, dividendAmount float64) (string, error) {
	// Fetch token details
	token, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return "", errors.New("failed to fetch token for dividend event")
	}

	// Decrypt token data before processing dividends
	decryptedToken, err := encryption.DecryptTokenData(token)
	if err != nil {
		return "", errors.New("failed to decrypt token data for dividend event")
	}

	// Ensure that dividends are paid to the correct owner
	eventDetails := "Dividend of " + common.FloatToString(dividendAmount) + " paid to owner " + decryptedToken.Owner

	// Record the dividend payment event
	eventID, err := RecordEvent(tokenID, "DIVIDEND_PAYMENT", eventDetails, decryptedToken.Owner)
	if err != nil {
		return "", err
	}

	return eventID, nil
}
