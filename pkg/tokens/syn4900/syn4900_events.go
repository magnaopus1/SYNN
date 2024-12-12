package syn4900

import (
	"errors"
	"sync"
	"time"
)

// Syn4900Event represents a significant event in the lifecycle of a SYN4900 token.
type Syn4900Event struct {
	EventID      string    `json:"event_id"`
	TokenID      string    `json:"token_id"`
	EventType    string    `json:"event_type"`  // Types like "Creation", "Transfer", "OwnershipChange", "AssetLinkVerified"
	Details      string    `json:"details"`     // Additional event details or metadata
	Timestamp    time.Time `json:"timestamp"`
	Verified     bool      `json:"verified"`
	VerifierID   string    `json:"verifier_id"` // Validator or entity verifying the event
}

// EventManager handles the logging and retrieval of token-related events.
type EventManager struct {
	mutex             sync.Mutex
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	eventRecords      map[string]*Syn4900Event
}

// NewEventManager creates a new instance of EventManager.
func NewEventManager(ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *EventManager {
	return &EventManager{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
		eventRecords:      make(map[string]*Syn4900Event),
	}
}

// LogEvent logs a new SYN4900 token event and stores it securely.
func (em *EventManager) LogEvent(tokenID, eventType, details, verifierID string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Generate a unique ID for the event.
	eventID := generateUniqueEventID()

	// Create the SYN4900 event.
	event := &Syn4900Event{
		EventID:    eventID,
		TokenID:    tokenID,
		EventType:  eventType,
		Details:    details,
		Timestamp:  time.Now(),
		Verified:   true, // Assuming events are verified in this implementation.
		VerifierID: verifierID,
	}

	// Encrypt the event before storing.
	encryptedEvent, err := em.encryptionService.EncryptData(event)
	if err != nil {
		return err
	}

	// Store the encrypted event in memory.
	em.eventRecords[eventID] = encryptedEvent.(*Syn4900Event)

	// Log the event in the ledger.
	if err := em.ledgerService.LogEvent("Syn4900EventLogged", time.Now(), eventID); err != nil {
		return err
	}

	// Validate the event using Synnergy Consensus.
	return em.consensusService.ValidateSubBlock(eventID)
}

// RetrieveEvent retrieves a specific SYN4900 event by its event ID.
func (em *EventManager) RetrieveEvent(eventID string) (*Syn4900Event, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Check if the event exists.
	event, exists := em.eventRecords[eventID]
	if !exists {
		return nil, errors.New("event not found for event ID: " + eventID)
	}

	// Decrypt the event before returning it.
	decryptedEvent, err := em.encryptionService.DecryptData(event)
	if err != nil {
		return nil, err
	}

	return decryptedEvent.(*Syn4900Event), nil
}


