package syn4900

import (
	"errors"
	"sync"
	"time"
)

// Event represents a specific action or occurrence related to SYN5000 tokens.
type Event struct {
	EventID      string    // Unique identifier for the event.
	TokenID      string    // Token involved in the event.
	EventType    string    // Type of event (e.g., "TokenMinted", "TokenBurned", "BetPlaced", "GameStarted").
	Timestamp    time.Time // Time when the event occurred.
	EventData    string    // Additional metadata or information about the event.
	EventHash    string    // Secure hash of the event for verification.
	Validated    bool      // Whether the event has been validated by Synnergy Consensus.
}

// EventTypes
const (
	EventMinted      = "TokenMinted"
	EventBurned      = "TokenBurned"
	EventBetPlaced   = "BetPlaced"
	EventGameStarted = "GameStarted"
	EventGameEnded   = "GameEnded"
	EventWinnings    = "WinningsDistributed"
	EventLoss        = "LossRecorded"
)

// EventManager handles the tracking, storing, and validation of events.
type EventManager struct {
	mu        sync.RWMutex
	events    map[string]*Event            // In-memory storage of events.
	ledger    *ledger.EventLedger          // Ledger for storing event data permanently.
	security  *encryption.Security         // Encryption for securing event data.
	consensus *consensus.SynnergyConsensus // Synnergy Consensus for validating events.
}

// NewEventManager creates a new instance of EventManager.
func NewEventManager(ledger *ledger.EventLedger, security *encryption.Security, consensus *consensus.SynnergyConsensus) *EventManager {
	return &EventManager{
		events:    make(map[string]*Event),
		ledger:    ledger,
		security:  security,
		consensus: consensus,
	}
}

// CreateEvent creates and logs a new event for a SYN5000 token.
func (em *EventManager) CreateEvent(tokenID, eventType, eventData string) (*Event, error) {
	em.mu.Lock()
	defer em.mu.Unlock()

	// Generate unique EventID and secure hash.
	eventID := generateUniqueID()
	timestamp := time.Now()
	eventHash := em.security.GenerateHash(tokenID + eventType + eventData + timestamp.String())

	// Create the event instance.
	event := &Event{
		EventID:   eventID,
		TokenID:   tokenID,
		EventType: eventType,
		Timestamp: timestamp,
		EventData: eventData,
		EventHash: eventHash,
		Validated: false, // Validation by Synnergy Consensus will happen later.
	}

	// Store the event in memory and ledger.
	em.events[eventID] = event
	em.ledger.StoreEvent(event)

	return event, nil
}

// ValidateEvent validates an event using Synnergy Consensus.
func (em *EventManager) ValidateEvent(eventID string) (*Event, error) {
	em.mu.Lock()
	defer em.mu.Unlock()

	// Retrieve the event.
	event, exists := em.events[eventID]
	if !exists {
		return nil, errors.New("event not found")
	}

	// Validate the event using Synnergy Consensus.
	subBlockHash, err := em.consensus.ValidateSubBlock(event.EventHash)
	if err != nil {
		return nil, err
	}

	// Mark the event as validated once 1000 sub-blocks form a full block.
	if err := em.validateBlock(subBlockHash); err != nil {
		return nil, err
	}

	// Update the event status.
	event.Validated = true
	em.events[eventID] = event

	// Store the updated event in the ledger.
	em.ledger.StoreEvent(event)

	return event, nil
}

// validateBlock validates a block of events after 1000 sub-blocks are validated.
func (em *EventManager) validateBlock(subBlockHash string) error {
	blockHash, err := em.consensus.ValidateBlock(subBlockHash)
	if err != nil {
		return err
	}

	// Store the validated block in the event ledger.
	return em.ledger.StoreValidatedBlock(blockHash)
}

// GetEvent retrieves an event by its ID.
func (em *EventManager) GetEvent(eventID string) (*Event, error) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	event, exists := em.events[eventID]
	if !exists {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// GetEventsByTokenID retrieves all events related to a specific token.
func (em *EventManager) GetEventsByTokenID(tokenID string) ([]*Event, error) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	var tokenEvents []*Event
	for _, event := range em.events {
		if event.TokenID == tokenID {
			tokenEvents = append(tokenEvents, event)
		}
	}

	if len(tokenEvents) == 0 {
		return nil, errors.New("no events found for the specified token")
	}

	return tokenEvents, nil
}

// GenerateUniqueID generates a unique ID for an event using Argon2.
func generateUniqueID() string {
	// Generate a unique identifier based on the current timestamp.
	return encryption.GenerateArgon2Hash(time.Now().String())
}
