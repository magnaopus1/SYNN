package syn4700

import (
	"errors"
	"sync"
	"time"

)

// Syn4700Event represents an event tied to legal contract management, such as agreement signing, contract termination, etc.
type Syn4700Event struct {
	EventID     string    `json:"event_id"`
	TokenID     string    `json:"token_id"`
	EventType   string    `json:"event_type"`   // Type of event (e.g., AgreementSigned, ContractTerminated)
	Details     string    `json:"details"`      // Additional details about the event
	Timestamp   time.Time `json:"timestamp"`    // Time the event occurred
	Signatures  map[string]string `json:"signatures"` // Party signatures involved in the event
}

// EventManager handles the lifecycle of Syn4700 events.
type EventManager struct {
	events           map[string][]Syn4700Event // Store events linked by token ID
	ledgerService    *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService *consensus.SynnergyConsensus
	mutex            sync.Mutex
}

// NewEventManager creates a new instance of the EventManager.
func NewEventManager(ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *EventManager {
	return &EventManager{
		events:           make(map[string][]Syn4700Event),
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService: consensus,
	}
}

// AddEvent creates a new event tied to a Syn4700 token and logs it in the system.
func (em *EventManager) AddEvent(tokenID, eventType, details string, signatures map[string]string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Create a new event
	event := Syn4700Event{
		EventID:    generateUniqueEventID(),
		TokenID:    tokenID,
		EventType:  eventType,
		Details:    details,
		Timestamp:  time.Now(),
		Signatures: signatures,
	}

	// Encrypt the event before storing it
	encryptedEvent, err := em.encryptionService.EncryptData(event)
	if err != nil {
		return err
	}

	// Store the event in the event history
	em.events[tokenID] = append(em.events[tokenID], event)

	// Log the event in the ledger
	if err := em.ledgerService.LogEvent(eventType, time.Now(), tokenID); err != nil {
		return err
	}

	// Validate the event using Synnergy Consensus
	if err := em.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// RetrieveEvents retrieves all events associated with a specific token ID.
func (em *EventManager) RetrieveEvents(tokenID string) ([]Syn4700Event, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	events, exists := em.events[tokenID]
	if !exists {
		return nil, ErrNoEventsFound
	}

	return events, nil
}

// generateUniqueEventID generates a unique identifier for an event.
func generateUniqueEventID() string {
	// Use a UUID or timestamp-based unique ID generation
	return "event-id-" + time.Now().Format("20060102150405")
}

// ErrNoEventsFound is returned when no events are found for a specific token ID.
var ErrNoEventsFound = errors.New("no events found for the given token ID")

