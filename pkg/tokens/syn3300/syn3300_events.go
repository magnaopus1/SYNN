package syn3300

import (
	"sync"
	"time"

)

// ETFEvent represents a significant event in the lifecycle of an ETF within the SYN3300 standard.
type ETFEvent struct {
	EventID        string    `json:"event_id"`        // Unique identifier for the event.
	ETFID          string    `json:"etf_id"`          // Associated ETF identifier.
	EventType      string    `json:"event_type"`      // Type of the event (e.g., "Creation", "Redemption", "Transfer").
	Timestamp      time.Time `json:"timestamp"`       // Timestamp when the event occurred.
	Details        string    `json:"details"`         // Additional details about the event.
	TransactionIDs []string  `json:"transaction_ids"` // Associated transactions for this event.
}

// ETFEventManager manages all ETF-related events and their logging.
type ETFEventManager struct {
	events           map[string]*ETFEvent
	ledgerService    *ledger.Ledger
	encryptionService *encryption.Encryptor
	consensusService *consensus.SynnergyConsensus
	mutex            sync.Mutex
}

// NewETFEventManager creates a new ETFEventManager with the necessary services.
func NewETFEventManager(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *ETFEventManager {
	return &ETFEventManager{
		events:           make(map[string]*ETFEvent),
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService: consensus,
	}
}

// LogEvent logs an ETF event, encrypts the data, and stores it in the ledger.
func (eem *ETFEventManager) LogEvent(etfID, eventType, details string, transactionIDs []string) (*ETFEvent, error) {
	eem.mutex.Lock()
	defer eem.mutex.Unlock()

	// Generate a unique EventID.
	eventID := eem.generateUniqueEventID()

	// Create the new event.
	event := &ETFEvent{
		EventID:        eventID,
		ETFID:          etfID,
		EventType:      eventType,
		Timestamp:      time.Now(),
		Details:        details,
		TransactionIDs: transactionIDs,
	}

	// Encrypt the event for security.
	encryptedEvent, err := eem.encryptionService.EncryptData(event)
	if err != nil {
		return nil, err
	}

	// Log the event in the ledger.
	if err := eem.ledgerService.LogEvent("ETFEventLogged", time.Now(), eventID); err != nil {
		return nil, err
	}

	// Validate the event with Synnergy Consensus.
	if err := eem.consensusService.ValidateSubBlock(eventID); err != nil {
		return nil, err
	}

	// Store the encrypted event in the manager.
	eem.events[eventID] = encryptedEvent.(*ETFEvent)

	return event, nil
}

// RetrieveEvent retrieves an ETF event by its EventID.
func (eem *ETFEventManager) RetrieveEvent(eventID string) (*ETFEvent, error) {
	eem.mutex.Lock()
	defer eem.mutex.Unlock()

	// Retrieve the event.
	event, exists := eem.events[eventID]
	if !exists {
		return nil, errors.New("event not found")
	}

	// Decrypt the event before returning it.
	decryptedEvent, err := eem.encryptionService.DecryptData(event)
	if err != nil {
		return nil, err
	}

	return decryptedEvent.(*ETFEvent), nil
}

// generateUniqueEventID generates a unique ID for an ETF event.
func (eem *ETFEventManager) generateUniqueEventID() string {
	// In a real-world scenario, this would be a more complex unique identifier generation logic.
	return "EVENT-" + time.Now().Format("20060102150405")
}

// ListEvents lists all events for a specific ETF.
func (eem *ETFEventManager) ListEvents(etfID string) ([]*ETFEvent, error) {
	eem.mutex.Lock()
	defer eem.mutex.Unlock()

	// Retrieve all events related to the specified ETF.
	var eventList []*ETFEvent
	for _, event := range eem.events {
		if event.ETFID == etfID {
			// Decrypt each event before adding to the list.
			decryptedEvent, err := eem.encryptionService.DecryptData(event)
			if err != nil {
				return nil, err
			}
			eventList = append(eventList, decryptedEvent.(*ETFEvent))
		}
	}

	return eventList, nil
}

// ListAllEvents returns all events logged in the system.
func (eem *ETFEventManager) ListAllEvents() ([]*ETFEvent, error) {
	eem.mutex.Lock()
	defer eem.mutex.Unlock()

	// Retrieve all events stored in the manager.
	var allEvents []*ETFEvent
	for _, event := range eem.events {
		// Decrypt each event before adding to the list.
		decryptedEvent, err := eem.encryptionService.DecryptData(event)
		if err != nil {
			return nil, err
		}
		allEvents = append(allEvents, decryptedEvent.(*ETFEvent))
	}

	return allEvents, nil
}

// DeleteEvent deletes an event from the system.
func (eem *ETFEventManager) DeleteEvent(eventID string) error {
	eem.mutex.Lock()
	defer eem.mutex.Unlock()

	// Check if the event exists.
	if _, exists := eem.events[eventID]; !exists {
		return errors.New("event not found")
	}

	// Remove the event from storage.
	delete(eem.events, eventID)

	// Log the deletion in the ledger.
	if err := eem.ledgerService.LogEvent("ETFEventDeleted", time.Now(), eventID); err != nil {
		return err
	}

	return nil
}
