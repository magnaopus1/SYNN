package syn2369

import (
	"time"
	"errors"
)


// SYN2369Event defines an event related to SYN2369Token.
type SYN2369Event struct {
	EventID     string    // Unique identifier for the event
	TokenID     string    // Associated token ID
	EventType   string    // Type of event (e.g., "Creation", "Transfer", "Update", "Compliance Check")
	Description string    // Detailed description of the event
	Timestamp   time.Time // Timestamp when the event occurred
	Owner       string    // Current owner of the token at the time of the event
	Encrypted   bool      // Indicates if the event data is encrypted
}

// LogEvent logs a new event for a SYN2369Token and stores it in the ledger.
func LogEvent(token common.SYN2369Token, eventType, description string) error {
	event := SYN2369Event{
		EventID:     generateEventID(),
		TokenID:     token.TokenID,
		EventType:   eventType,
		Description: description,
		Timestamp:   time.Now(),
		Owner:       token.Owner,
	}

	// Encrypt the event details if sensitive
	if shouldEncrypt(eventType) {
		encryptedData, err := encryptEventDetails(event)
		if err != nil {
			return err
		}
		event.Description = encryptedData
		event.Encrypted = true
	} else {
		event.Encrypted = false
	}

	// Add the event to the token's history
	token.EventHistory = append(token.EventHistory, event)

	// Store event in the ledger for immutability and traceability
	err := ledger.StoreEvent(event)
	if err != nil {
		return err
	}

	return nil
}

// GenerateTokenCreationEvent generates an event for token creation.
func GenerateTokenCreationEvent(token common.SYN2369Token) error {
	description := "Token created with ID: " + token.TokenID + " by owner: " + token.Owner
	return LogEvent(token, "Creation", description)
}

// GenerateTokenTransferEvent generates an event for token transfer.
func GenerateTokenTransferEvent(token common.SYN2369Token, oldOwner, newOwner string) error {
	description := "Token transferred from " + oldOwner + " to " + newOwner
	return LogEvent(token, "Transfer", description)
}

// GenerateMetadataUpdateEvent generates an event for a metadata update.
func GenerateMetadataUpdateEvent(token common.SYN2369Token, updateDescription string) error {
	description := "Metadata updated: " + updateDescription
	return LogEvent(token, "Metadata Update", description)
}

// GenerateComplianceCheckEvent logs an event for a compliance check on the token.
func GenerateComplianceCheckEvent(token common.SYN2369Token, complianceStatus string) error {
	description := "Compliance check performed. Status: " + complianceStatus
	return LogEvent(token, "Compliance Check", description)
}

// GenerateCustomEvent generates a custom event for any interaction with the SYN2369Token.
func GenerateCustomEvent(token common.SYN2369Token, eventType, customDescription string) error {
	return LogEvent(token, eventType, customDescription)
}

// shouldEncrypt determines if event data needs encryption based on event type.
func shouldEncrypt(eventType string) bool {
	// Encrypt sensitive events like Transfer or Metadata Update
	sensitiveEventTypes := []string{"Transfer", "Metadata Update", "Compliance Check"}
	for _, e := range sensitiveEventTypes {
		if eventType == e {
			return true
		}
	}
	return false
}

// encryptEventDetails encrypts event details to ensure privacy and security.
func encryptEventDetails(event SYN2369Event) (string, error) {
	// Assuming encryption.EncryptData is a method to encrypt event details
	encryptedData, err := encryption.EncryptData(event.Description)
	if err != nil {
		return "", errors.New("failed to encrypt event details")
	}
	return encryptedData, nil
}

// generateEventID generates a unique event ID.
func generateEventID() string {
	return encryption.GenerateRandomID()
}

// FetchTokenEventHistory fetches the event history of a given token from the ledger.
func FetchTokenEventHistory(tokenID string) ([]SYN2369Event, error) {
	// Fetch the event history for the token from the ledger
	eventHistory, err := ledger.GetEventHistory(tokenID)
	if err != nil {
		return nil, err
	}
	return eventHistory, nil
}

// RevokeEvent revokes or marks an event as invalid due to a compliance issue or error.
func RevokeEvent(eventID string) error {
	// Mark event as revoked in the ledger and in token history
	err := ledger.MarkEventRevoked(eventID)
	if err != nil {
		return err
	}
	return nil
}
