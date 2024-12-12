package syn1100

import (
	"errors"
	"sync"
	"time"

)

// SYN1100EventManager manages events related to healthcare data tokens (SYN1100)
type SYN1100EventManager struct {
	Ledger            *ledger.Ledger                // Ledger integration for event logging
	ConsensusEngine   *consensus.SynnergyConsensus  // Synnergy Consensus for event validation
	EncryptionService *encryption.EncryptionService // Encryption service for secure event handling
	mutex             sync.Mutex                    // Mutex for concurrency control
}

// HealthcareEvent represents an event related to healthcare data, such as patient interactions
type HealthcareEvent struct {
	EventID     string    `json:"event_id"`
	TokenID     string    `json:"token_id"`
	EventType   string    `json:"event_type"`      // Type of healthcare event
	Description string    `json:"description"`     // Detailed description of the event
	Timestamp   time.Time `json:"timestamp"`       // Time the event occurred
	EncryptedData string  `json:"encrypted_data"`  // Encrypted event data, if sensitive
}

// EventLogEntry represents a log entry for an event
type EventLogEntry struct {
	EventID      string    `json:"event_id"`
	TokenID      string    `json:"token_id"`
	EventType    string    `json:"event_type"`
	LoggedBy     string    `json:"logged_by"`      // ID of the individual or system logging the event
	Timestamp    time.Time `json:"timestamp"`
}

// RegisterEvent logs a new healthcare event for the SYN1100 token
func (em *SYN1100EventManager) RegisterEvent(tokenID, eventType, description string, sensitiveData string) (*HealthcareEvent, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Validate if token exists in the ledger
	encryptedToken, encryptionKey, err := em.Ledger.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve token from ledger")
	}

	// Decrypt token to ensure it is valid
	tokenData, err := em.EncryptionService.DecryptData([]byte(encryptedToken), encryptionKey)
	if err != nil {
		return nil, errors.New("failed to decrypt token data")
	}

	// Create a new event with a unique ID and encrypt the event data if sensitive
	eventID := common.GenerateUUID()
	encryptedData, err := em.EncryptionService.EncryptData([]byte(sensitiveData), encryptionKey)
	if err != nil {
		return nil, errors.New("failed to encrypt sensitive data for event")
	}

	event := &HealthcareEvent{
		EventID:      eventID,
		TokenID:      tokenID,
		EventType:    eventType,
		Description:  description,
		Timestamp:    time.Now(),
		EncryptedData: string(encryptedData),
	}

	// Log the event in the ledger for future audits
	eventLog := EventLogEntry{
		EventID:      eventID,
		TokenID:      tokenID,
		EventType:    eventType,
		LoggedBy:     "system", // System or individual logging this event
		Timestamp:    time.Now(),
	}

	// Store event in ledger
	if err := em.Ledger.StoreEventLog(eventLog.EventID, common.StructToString(eventLog), encryptionKey); err != nil {
		return nil, errors.New("failed to store event log in ledger")
	}

	// Validate the event registration using Synnergy Consensus
	if err := em.ConsensusEngine.ValidateEvent(event.EventID); err != nil {
		return nil, errors.New("event validation failed via Synnergy Consensus")
	}

	return event, nil
}

// RetrieveEvent retrieves a healthcare event by its ID
func (em *SYN1100EventManager) RetrieveEvent(eventID string) (*HealthcareEvent, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Retrieve event from the ledger
	encryptedEvent, encryptionKey, err := em.Ledger.GetEventLog(eventID)
	if err != nil {
		return nil, errors.New("failed to retrieve event from ledger")
	}

	// Decrypt event data
	eventData, err := em.EncryptionService.DecryptData([]byte(encryptedEvent), encryptionKey)
	if err != nil {
		return nil, errors.New("failed to decrypt event data")
	}

	// Unmarshal the event data
	var event HealthcareEvent
	if err := common.StringToStruct(string(eventData), &event); err != nil {
		return nil, errors.New("failed to unmarshal event data")
	}

	return &event, nil
}

// ListEventsByToken retrieves all healthcare events associated with a specific token
func (em *SYN1100EventManager) ListEventsByToken(tokenID string) ([]HealthcareEvent, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Fetch all event logs related to the token from the ledger
	eventLogs, err := em.Ledger.GetAllEventLogsByTokenID(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve event logs from ledger")
	}

	var events []HealthcareEvent
	for _, encryptedEvent := range eventLogs {
		// Decrypt each event log
		eventData, encryptionKey, err := em.Ledger.GetEventLog(encryptedEvent.EventID)
		if err != nil {
			return nil, errors.New("failed to retrieve and decrypt event data from ledger")
		}

		decryptedData, err := em.EncryptionService.DecryptData([]byte(eventData), encryptionKey)
		if err != nil {
			return nil, errors.New("failed to decrypt event data")
		}

		// Unmarshal the event
		var event HealthcareEvent
		if err := common.StringToStruct(string(decryptedData), &event); err != nil {
			return nil, errors.New("failed to unmarshal event data")
		}
		events = append(events, event)
	}

	return events, nil
}

// ValidateEvent ensures an event follows the required regulations and integrity before being recorded
func (em *SYN1100EventManager) ValidateEvent(eventID string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Validate through Synnergy Consensus
	if err := em.ConsensusEngine.ValidateEvent(eventID); err != nil {
		return errors.New("event validation failed via Synnergy Consensus")
	}

	return nil
}

