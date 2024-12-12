package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"

// Event represents a logged event in the SYN20 token system.
type Event struct {
	EventID      string    // Unique event ID
	Timestamp    time.Time // Timestamp of the event
	EventType    string    // Type of event (e.g., "Mint", "Burn", "Transfer")
	ContractID   string    // SYN20 contract ID related to the event
	Details      string    // Event-specific details (e.g., recipient, amount)
	EncryptedLog string    // Encrypted version of the event log
}

// EventLogger is responsible for logging events related to SYN20 tokens.
type EventLogger struct {
	mutex       sync.Mutex      // Thread safety for concurrent logging
	Ledger      *ledger.Ledger  // Reference to the ledger for logging events
	Encryption  *encryption.Encryption // Encryption service for securing event data
	EventLog    map[string]*Event // In-memory map of logged events
}

// NewEventLogger initializes a new event logger for SYN20 token events.
func NewEventLogger(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *EventLogger {
	return &EventLogger{
		Ledger:     ledgerInstance,
		Encryption: encryptionService,
		EventLog:   make(map[string]*Event),
	}
}

// LogEvent creates a new event log entry and stores it in the ledger.
func (logger *EventLogger) LogEvent(eventType, contractID, details string) (string, error) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	// Create a unique event ID
	eventID := common.GenerateUniqueID()

	// Create the event structure
	event := &Event{
		EventID:    eventID,
		Timestamp:  time.Now(),
		EventType:  eventType,
		ContractID: contractID,
		Details:    details,
	}

	// Encrypt the event details before storing
	encryptedDetails, err := logger.Encryption.EncryptData(event.Details, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt event details: %v", err)
	}
	event.EncryptedLog = encryptedDetails

	// Store the event in the ledger
	if err := logger.Ledger.RecordEvent(eventID, event.EncryptedLog); err != nil {
		return "", fmt.Errorf("failed to store event in ledger: %v", err)
	}

	// Save the event in memory
	logger.EventLog[eventID] = event

	fmt.Printf("Event logged successfully with ID: %s, Type: %s, Contract: %s\n", eventID, eventType, contractID)
	return eventID, nil
}

// GetEvent retrieves an event by its ID.
func (logger *EventLogger) GetEvent(eventID string) (*Event, error) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	// Check if the event exists in the log
	event, exists := logger.EventLog[eventID]
	if !exists {
		return nil, fmt.Errorf("event with ID %s not found", eventID)
	}

	// Decrypt the event details
	decryptedDetails, err := logger.Encryption.DecryptData(event.EncryptedLog, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt event details: %v", err)
	}
	event.Details = decryptedDetails

	return event, nil
}

// ListEvents returns all logged events for a specific contract ID.
func (logger *EventLogger) ListEvents(contractID string) ([]*Event, error) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	var events []*Event
	for _, event := range logger.EventLog {
		if event.ContractID == contractID {
			// Decrypt the event details before returning
			decryptedDetails, err := logger.Encryption.DecryptData(event.EncryptedLog, common.EncryptionKey)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt event details for event %s: %v", event.EventID, err)
			}
			event.Details = decryptedDetails
			events = append(events, event)
		}
	}

	return events, nil
}

// DeleteEvent removes an event from the log.
func (logger *EventLogger) DeleteEvent(eventID string) error {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	// Check if the event exists in the log
	if _, exists := logger.EventLog[eventID]; !exists {
		return fmt.Errorf("event with ID %s not found", eventID)
	}

	// Remove the event from memory and ledger
	delete(logger.EventLog, eventID)
	if err := logger.Ledger.DeleteEvent(eventID); err != nil {
		return fmt.Errorf("failed to remove event from ledger: %v", err)
	}

	fmt.Printf("Event %s deleted successfully.\n", eventID)
	return nil
}

// StartEventLogger starts a background service that listens for specific blockchain events to log.
func (logger *EventLogger) StartEventLogger() {
	// This function could listen for blockchain events (e.g., transactions, contract interactions)
	// and automatically log them. It could integrate with other components of the system like the
	// transaction pool or the virtual machine to capture important events.
	fmt.Println("Event logger is now running and listening for blockchain events...")
}

