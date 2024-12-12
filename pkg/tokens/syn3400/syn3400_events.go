package syn3400

import (
	"errors"
	"sync"
	"time"

)

// ForexEvent represents an event that occurs in the Forex trading ecosystem.
type ForexEvent struct {
	EventID       string    `json:"event_id"`
	PairID        string    `json:"pair_id"`
	EventType     string    `json:"event_type"`
	EventDetails  string    `json:"event_details"`
	Timestamp     time.Time `json:"timestamp"`
	Encrypted     bool      `json:"encrypted"`
	SecureHash    string    `json:"secure_hash"`
	ValidationStatus string `json:"validation_status"`
}

// ForexEventManager manages the creation, validation, and encryption of Forex events.
type ForexEventManager struct {
	events     map[string]*ForexEvent
	ledger     *ledger.Ledger
	encryptor  *encryption.Encryptor
	consensus  *consensus.SynnergyConsensus
	mutex      sync.Mutex
}

// NewForexEventManager creates a new instance of ForexEventManager.
func NewForexEventManager(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *ForexEventManager {
	return &ForexEventManager{
		events:    make(map[string]*ForexEvent),
		ledger:    ledger,
		encryptor: encryptor,
		consensus: consensus,
	}
}

// CreateEvent creates a new Forex event, encrypts it, logs it in the ledger, and validates it through consensus.
func (fem *ForexEventManager) CreateEvent(pairID, eventType, eventDetails string) (*ForexEvent, error) {
	fem.mutex.Lock()
	defer fem.mutex.Unlock()

	// Validate input data
	if pairID == "" || eventType == "" {
		return nil, errors.New("invalid event: pairID and eventType must not be empty")
	}

	// Create new Forex event
	eventID := generateUniqueEventID()
	timestamp := time.Now()
	event := &ForexEvent{
		EventID:      eventID,
		PairID:       pairID,
		EventType:    eventType,
		EventDetails: eventDetails,
		Timestamp:    timestamp,
		Encrypted:    false,
	}

	// Encrypt the event data
	encryptedEvent, err := fem.encryptor.EncryptData(event)
	if err != nil {
		return nil, err
	}
	event = encryptedEvent.(*ForexEvent)
	event.Encrypted = true

	// Generate secure hash for event integrity
	event.SecureHash = generateSecureHash(eventID, pairID, eventType, eventDetails, timestamp)

	// Store the event
	fem.events[event.EventID] = event

	// Log the event in the ledger
	fem.ledger.LogEvent("ForexEventCreated", timestamp, event.EventID)

	// Validate the event using Synnergy Consensus
	err = fem.consensus.ValidateSubBlock(event.EventID)
	if err != nil {
		event.ValidationStatus = "Failed"
	} else {
		event.ValidationStatus = "Success"
	}

	// Return the newly created event
	return event, nil
}

// GetEvent retrieves an event by its ID, ensuring it is decrypted.
func (fem *ForexEventManager) GetEvent(eventID string) (*ForexEvent, error) {
	fem.mutex.Lock()
	defer fem.mutex.Unlock()

	// Retrieve the event
	event, exists := fem.events[eventID]
	if !exists {
		return nil, errors.New("event not found")
	}

	// Decrypt the event data if necessary
	if event.Encrypted {
		decryptedEvent, err := fem.encryptor.DecryptData(event)
		if err != nil {
			return nil, err
		}
		event = decryptedEvent.(*ForexEvent)
		event.Encrypted = false
	}

	return event, nil
}

// ListEvents lists all stored Forex events.
func (fem *ForexEventManager) ListEvents() ([]*ForexEvent, error) {
	fem.mutex.Lock()
	defer fem.mutex.Unlock()

	eventsList := []*ForexEvent{}
	for _, event := range fem.events {
		eventsList = append(eventsList, event)
	}

	return eventsList, nil
}

// DeleteEvent removes an event from the system and logs the deletion in the ledger.
func (fem *ForexEventManager) DeleteEvent(eventID string) error {
	fem.mutex.Lock()
	defer fem.mutex.Unlock()

	// Check if the event exists
	if _, exists := fem.events[eventID]; !exists {
		return errors.New("event does not exist")
	}

	// Delete the event
	delete(fem.events, eventID)

	// Log the event deletion in the ledger
	fem.ledger.LogEvent("ForexEventDeleted", time.Now(), eventID)

	return nil
}

// generateUniqueEventID generates a unique identifier for Forex events.
func generateUniqueEventID() string {
	return generateUniqueID() // Assumes generateUniqueID exists elsewhere in your project
}

// generateSecureHash generates a secure hash to ensure event integrity.
func generateSecureHash(eventID, pairID, eventType, eventDetails string, timestamp time.Time) string {
	hashData := eventID + pairID + eventType + eventDetails + timestamp.String()
	return encryption.GenerateHash(hashData) // Assumes GenerateHash exists in the encryption package
}
