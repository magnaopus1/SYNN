package syn300

import (
	"errors"
	"sync"
	"time"
)

// Syn300EventType defines the type of events in the SYN300 token ecosystem
type Syn300EventType string

const (
	Minted    Syn300EventType = "Minted"
	Burned    Syn300EventType = "Burned"
	Transferred Syn300EventType = "Transferred"
	Delegated  Syn300EventType = "Delegated"
	Voted      Syn300EventType = "Voted"
)

// Syn300Event represents an event in the SYN300 token ecosystem
type Syn300Event struct {
	ID         string           `json:"id"`
	Type       Syn300EventType  `json:"type"`
	Timestamp  time.Time        `json:"timestamp"`
	Details    string           `json:"details"`
	Signature  string           `json:"signature"`
	Validated  bool             `json:"validated"`
	Encrypted  bool             `json:"encrypted"`
}

// Syn300EventManager manages events for the SYN300 token standard
type Syn300EventManager struct {
	Ledger *ledger.Ledger
	Events map[string]Syn300Event
	mutex  sync.RWMutex
}

// NewSyn300EventManager creates a new event manager for SYN300 tokens
func NewSyn300EventManager(ledger *ledger.Ledger) *Syn300EventManager {
	return &Syn300EventManager{
		Ledger: ledger,
		Events: make(map[string]Syn300Event),
	}
}

// CreateEvent generates a new event for the SYN300 token
func (em *Syn300EventManager) CreateEvent(eventType Syn300EventType, details string, signature string) (string, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	eventID := generateEventID()
	encryptedDetails, err := encryption.Encrypt(details)
	if err != nil {
		return "", errors.New("failed to encrypt event details")
	}

	event := Syn300Event{
		ID:         eventID,
		Type:       eventType,
		Timestamp:  time.Now(),
		Details:    encryptedDetails,
		Signature:  signature,
		Validated:  false,
		Encrypted:  true,
	}

	// Store the event
	em.Events[eventID] = event

	// Log the event to the ledger
	if err := em.Ledger.StoreEvent(event); err != nil {
		return "", errors.New("failed to store event in the ledger")
	}

	return eventID, nil
}

// ValidateEvent validates an event using Synnergy Consensus and marks it as validated
func (em *Syn300EventManager) ValidateEvent(eventID string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	event, exists := em.Events[eventID]
	if !exists {
		return errors.New("event not found")
	}

	if event.Validated {
		return errors.New("event already validated")
	}

	// Validate the event using Synnergy Consensus
	if err := consensus.ValidateEvent(eventID, event.Signature); err != nil {
		return errors.New("failed to validate event under Synnergy Consensus")
	}

	// Mark the event as validated
	event.Validated = true
	em.Events[eventID] = event

	// Update the ledger with the validated event
	if err := em.Ledger.UpdateEvent(event); err != nil {
		return errors.New("failed to update event in the ledger")
	}

	return nil
}

// GetEvent retrieves an event by its ID
func (em *Syn300EventManager) GetEvent(eventID string) (Syn300Event, error) {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	event, exists := em.Events[eventID]
	if !exists {
		return Syn300Event{}, errors.New("event not found")
	}

	// Decrypt event details before returning
	decryptedDetails, err := encryption.Decrypt(event.Details)
	if err != nil {
		return Syn300Event{}, errors.New("failed to decrypt event details")
	}

	event.Details = decryptedDetails
	return event, nil
}

// GetValidatedEvents retrieves all validated events
func (em *Syn300EventManager) GetValidatedEvents() ([]Syn300Event, error) {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	var validatedEvents []Syn300Event
	for _, event := range em.Events {
		if event.Validated {
			decryptedDetails, err := encryption.Decrypt(event.Details)
			if err != nil {
				return nil, errors.New("failed to decrypt event details")
			}
			event.Details = decryptedDetails
			validatedEvents = append(validatedEvents, event)
		}
	}

	return validatedEvents, nil
}

// DeleteEvent removes an event from the system and the ledger
func (em *Syn300EventManager) DeleteEvent(eventID string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	_, exists := em.Events[eventID]
	if !exists {
		return errors.New("event not found")
	}

	// Remove from the events list
	delete(em.Events, eventID)

	// Remove from ledger as well
	if err := em.Ledger.RemoveEvent(eventID); err != nil {
		return errors.New("failed to remove event from ledger")
	}

	return nil
}

// StoreAllEvents logs all events to the ledger
func (em *Syn300EventManager) StoreAllEvents() error {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	for _, event := range em.Events {
		if err := em.Ledger.StoreEvent(event); err != nil {
			return err
		}
	}
	return nil
}

// Helper function to generate unique event IDs
func generateEventID() string {
	// This is a placeholder. Replace it with a proper unique ID generator.
	return "event_" + time.Now().Format("20060102150405")
}

