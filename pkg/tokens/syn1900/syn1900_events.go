package syn1900

import (
	"errors"
	"time"
)


// EventManager handles event logging and event-related operations for SYN1900 tokens.
type EventManager struct {
	ledger LedgerInterface // Interface for interacting with the ledger
}

// LedgerInterface defines methods for interacting with the ledger.
type LedgerInterface interface {
	GetTokenByID(tokenID string) (common.SYN1900Token, error)
	UpdateToken(token common.SYN1900Token) error
	AddEventLog(eventLog common.EventLog) error
	GetEventLogsByToken(tokenID string) ([]common.EventLog, error)
}

// LogEvent creates and stores an event log for a specific token action.
func (em *EventManager) LogEvent(tokenID, eventType, description string) error {
	// Fetch the token from the ledger
	token, err := em.ledger.GetTokenByID(tokenID)
	if err != nil {
		return errors.New("token not found in the ledger")
	}

	// Create a new event log
	eventLog := common.EventLog{
		TokenID:     token.TokenID,
		EventType:   eventType,
		Description: description,
		Timestamp:   time.Now(),
	}

	// Log the event to the ledger
	err = em.ledger.AddEventLog(eventLog)
	if err != nil {
		return errors.New("failed to log event to the ledger")
	}

	return nil
}

// FetchTokenEvents retrieves all event logs associated with a specific token.
func (em *EventManager) FetchTokenEvents(tokenID string) ([]common.EventLog, error) {
	// Fetch event logs from the ledger
	eventLogs, err := em.ledger.GetEventLogsByToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve event logs from the ledger")
	}
	return eventLogs, nil
}

// GenerateEventReport generates a detailed report of all events associated with a specific token.
func (em *EventManager) GenerateEventReport(tokenID string) (string, error) {
	// Fetch all events related to the token
	eventLogs, err := em.FetchTokenEvents(tokenID)
	if err != nil {
		return "", err
	}

	// Generate the event report
	report := "Event Report for Token ID: " + tokenID + "\n"
	report += "----------------------------------\n"
	for _, log := range eventLogs {
		report += "Event Type: " + log.EventType + "\n"
		report += "Description: " + log.Description + "\n"
		report += "Timestamp: " + log.Timestamp.Format(time.RFC3339) + "\n\n"
	}

	return report, nil
}

// EncryptEventReport encrypts the event report before sending it to external systems.
func (em *EventManager) EncryptEventReport(report string) ([]byte, error) {
	encryptedReport, err := encryption.Encrypt([]byte(report))
	if err != nil {
		return nil, errors.New("failed to encrypt event report")
	}
	return encryptedReport, nil
}

// AddEventForCompletion logs a "Completion" event for a course that has been completed by a recipient.
func (em *EventManager) AddEventForCompletion(tokenID string, recipientID string) error {
	description := "Course completed by recipient: " + recipientID
	return em.LogEvent(tokenID, "Completion", description)
}

// AddEventForTransfer logs a "Transfer" event when educational credits are transferred between parties.
func (em *EventManager) AddEventForTransfer(tokenID, fromID, toID string) error {
	description := "Token transferred from " + fromID + " to " + toID
	return em.LogEvent(tokenID, "Transfer", description)
}

// AddEventForRevocation logs a "Revocation" event when an educational credit is revoked.
func (em *EventManager) AddEventForRevocation(tokenID, reason string) error {
	description := "Token revoked: " + reason
	return em.LogEvent(tokenID, "Revocation", description)
}
