package syn845

import (
	"errors"
	"sync"
	"time"

)

// EventType defines the types of events for the SYN845 token standard
type EventType string

const (
	IssuanceEvent          EventType = "Issuance"
	RepaymentEvent         EventType = "Repayment"
	RefinancingEvent       EventType = "Refinancing"
	OwnershipTransferEvent EventType = "OwnershipTransfer"
	PenaltyEvent           EventType = "Penalty"
	DefaultEvent           EventType = "Default"
	InterestAdjustmentEvent EventType = "InterestAdjustment"
)

// Event represents a blockchain event for SYN845 tokens
type Event struct {
	ID            string                 `json:"id"`
	Type          EventType              `json:"type"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          map[string]interface{} `json:"data"`
	EncryptedData string                 `json:"encrypted_data,omitempty"`
	EncryptionKey string                 `json:"encryption_key,omitempty"`
}

// EventLog handles logging of SYN845 token events, including encryption and validation
type EventLog struct {
	events            []Event
	mutex             sync.RWMutex
	Ledger            *ledger.Ledger                // Ledger for recording events
	ConsensusEngine   *consensus.SynnergyConsensus  // Synnergy Consensus for validating events
	EncryptionService *encryption.EncryptionService // Encryption service for securing event data
}

// NewEventLog creates a new EventLog instance
func NewEventLog(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *EventLog {
	return &EventLog{
		events:            make([]Event, 0),
		Ledger:            ledger,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
	}
}

// LogEvent logs an event to the event log and securely stores it in the ledger
func (el *EventLog) LogEvent(eventType EventType, data map[string]interface{}, encrypt bool) error {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	eventID := generateUniqueID()
	timestamp := time.Now()

	// Encrypt event data if encryption is enabled
	var encryptedData, encryptionKey string
	var err error
	if encrypt {
		encryptedData, encryptionKey, err = el.EncryptionService.EncryptData([]byte(common.MapToString(data)))
		if err != nil {
			return errors.New("failed to encrypt event data")
		}
	}

	// Create the event object
	event := Event{
		ID:            eventID,
		Type:          eventType,
		Timestamp:     timestamp,
		Data:          data,
		EncryptedData: encryptedData,
		EncryptionKey: encryptionKey,
	}

	// Validate the event using Synnergy Consensus
	if err := el.ConsensusEngine.ValidateEvent(event); err != nil {
		return errors.New("event validation failed via Synnergy Consensus")
	}

	// Record the event in the ledger
	if err := el.Ledger.RecordEvent(event.ID, event); err != nil {
		return errors.New("failed to record event in the ledger")
	}

	el.events = append(el.events, event)
	return nil
}

// GetEvents returns all logged events
func (el *EventLog) GetEvents() ([]Event, error) {
	el.mutex.RLock()
	defer el.mutex.RUnlock()

	return el.events, nil
}

// GetEventsByType returns all events of a specific type
func (el *EventLog) GetEventsByType(eventType EventType) ([]Event, error) {
	el.mutex.RLock()
	defer el.mutex.RUnlock()

	var filteredEvents []Event
	for _, event := range el.events {
		if event.Type == eventType {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents, nil
}

// GetEventsByDateRange returns all events within a specific date range
func (el *EventLog) GetEventsByDateRange(startDate, endDate time.Time) ([]Event, error) {
	el.mutex.RLock()
	defer el.mutex.RUnlock()

	var filteredEvents []Event
	for _, event := range el.events {
		if event.Timestamp.After(startDate) && event.Timestamp.Before(endDate) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents, nil
}

// RetrieveEvent retrieves a specific event by its ID and decrypts it if necessary
func (el *EventLog) RetrieveEvent(eventID string) (*Event, error) {
	el.mutex.RLock()
	defer el.mutex.RUnlock()

	// Retrieve event from ledger
	event, err := el.Ledger.GetEvent(eventID)
	if err != nil {
		return nil, errors.New("failed to retrieve event from ledger")
	}

	// Decrypt event data if it was encrypted
	if event.EncryptedData != "" {
		decryptedData, err := el.EncryptionService.DecryptData([]byte(event.EncryptedData), event.EncryptionKey)
		if err != nil {
			return nil, errors.New("failed to decrypt event data")
		}
		event.Data = common.StringToMap(string(decryptedData))
	}

	return &event, nil
}

// generateUniqueID generates a unique ID for an event
func generateUniqueID() string {
	return time.Now().Format("20060102150405") + "_" + common.GenerateRandomString(8)
}
