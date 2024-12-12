package syn12

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"

)


// SYN12EventType represents the different types of events for SYN12 tokens.
type SYN12EventType string

const (
	EventTokenIssued      SYN12EventType = "TokenIssued"
	EventTokenRedeemed    SYN12EventType = "TokenRedeemed"
	EventInterestAccrued  SYN12EventType = "InterestAccrued"
	EventTokenMatured     SYN12EventType = "TokenMatured"
	EventEarlyRedemption  SYN12EventType = "EarlyRedemption"
	EventRateUpdated      SYN12EventType = "RateUpdated"
)

// SYN12Event represents a structure for logging events related to SYN12 tokens.
type SYN12Event struct {
	EventID        string          // Unique event identifier
	TokenID        string          // Token associated with the event
	EventType      SYN12EventType  // Type of the event
	Timestamp      time.Time       // Time when the event occurred
	Details        string          // Additional details about the event
	Encrypted      bool            // Indicates if the event data is encrypted
}

// SYN12EventManager handles the lifecycle and management of token-related events.
type SYN12EventManager struct {
	ledgerManager     *ledger.LedgerManager         // Ledger for event logging
	encryptionService *encryption.EncryptionService // Encryption service for secure event data
	mutex             sync.Mutex                    // Mutex for concurrency control
}

// NewSYN12EventManager initializes a new SYN12EventManager.
func NewSYN12EventManager(ledgerManager *ledger.LedgerManager, encryptionService *encryption.EncryptionService) *SYN12EventManager {
	return &SYN12EventManager{
		ledgerManager:     ledgerManager,
		encryptionService: encryptionService,
	}
}

// LogEvent logs a new event to the ledger for SYN12 tokens.
func (em *SYN12EventManager) LogEvent(tokenID string, eventType SYN12EventType, details string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Create a new event
	eventID := common.GenerateUUID() // Generate a unique event ID
	event := SYN12Event{
		EventID:    eventID,
		TokenID:    tokenID,
		EventType:  eventType,
		Timestamp:  time.Now(),
		Details:    details,
		Encrypted:  false,
	}

	// Encrypt the event details if encryption is enabled
	encryptedDetails, err := em.encryptionService.Encrypt([]byte(details))
	if err != nil {
		return fmt.Errorf("failed to encrypt event details: %v", err)
	}
	event.Details = string(encryptedDetails)
	event.Encrypted = true

	// Log the event to the ledger
	if err := em.ledgerManager.RecordEvent(event); err != nil {
		return fmt.Errorf("failed to log event: %v", err)
	}

	return nil
}

// GetEvent retrieves an event by its ID.
func (em *SYN12EventManager) GetEvent(eventID string) (*SYN12Event, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Retrieve the event from the ledger
	event, err := em.ledgerManager.GetEvent(eventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %v", err)
	}

	// Decrypt the event details if necessary
	if event.Encrypted {
		decryptedDetails, err := em.encryptionService.Decrypt([]byte(event.Details))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt event details: %v", err)
		}
		event.Details = string(decryptedDetails)
		event.Encrypted = false
	}

	return event, nil
}

// ListEventsByToken lists all events associated with a specific token.
func (em *SYN12EventManager) ListEventsByToken(tokenID string) ([]SYN12Event, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Retrieve the list of events from the ledger
	events, err := em.ledgerManager.GetEventsByTokenID(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %v", err)
	}

	// Decrypt event details if necessary
	for i, event := range events {
		if event.Encrypted {
			decryptedDetails, err := em.encryptionService.Decrypt([]byte(event.Details))
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt event details for event ID %s: %v", event.EventID, err)
			}
			events[i].Details = string(decryptedDetails)
			events[i].Encrypted = false
		}
	}

	return events, nil
}

// ListEventsByType lists all events of a specific type.
func (em *SYN12EventManager) ListEventsByType(eventType SYN12EventType) ([]SYN12Event, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Retrieve the list of events by type from the ledger
	events, err := em.ledgerManager.GetEventsByType(string(eventType))
	if err != nil {
		return nil, fmt.Errorf("failed to list events by type: %v", err)
	}

	// Decrypt event details if necessary
	for i, event := range events {
		if event.Encrypted {
			decryptedDetails, err := em.encryptionService.Decrypt([]byte(event.Details))
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt event details for event ID %s: %v", event.EventID, err)
			}
			events[i].Details = string(decryptedDetails)
			events[i].Encrypted = false
		}
	}

	return events, nil
}
