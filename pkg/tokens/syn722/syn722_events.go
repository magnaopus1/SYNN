package syn722

import (
	"errors"
	"sync"
	"time"

)

// SYN722Event represents a blockchain event related to SYN722 tokens, such as transfers, creation, and compliance events
type SYN722Event struct {
	ID          string                 `json:"id"`
	EventType   string                 `json:"event_type"` // "transfer", "creation", "compliance", "mode_change", etc.
	TokenID     string                 `json:"token_id"`
	Initiator   string                 `json:"initiator"`  // The user who initiated the event
	Timestamp   time.Time              `json:"timestamp"`
	Details     map[string]interface{} `json:"details"`    // Event-specific details
	Encrypted   bool                   `json:"encrypted"`
	EncryptedData string               `json:"encrypted_data,omitempty"`
	EncryptionKey string               `json:"encryption_key,omitempty"`
}

// SYN722EventManager manages all SYN722-related events, ensuring events are securely logged, validated through Synnergy Consensus, and stored in the ledger
type SYN722EventManager struct {
	Ledger            *ledger.Ledger                // Ledger for recording events
	ConsensusEngine   *consensus.SynnergyConsensus  // Synnergy Consensus for event validation
	EncryptionService *encryption.EncryptionService // Encryption service for securing event data
	mutex             sync.Mutex                    // Mutex for safe concurrent access
}

// NewSYN722EventManager initializes a new instance of SYN722EventManager
func NewSYN722EventManager(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *SYN722EventManager {
	return &SYN722EventManager{
		Ledger:            ledger,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
	}
}

// LogEvent logs a new SYN722 event (transfer, creation, compliance, etc.) and stores it in the ledger
func (em *SYN722EventManager) LogEvent(eventType, tokenID, initiator string, details map[string]interface{}, encrypt bool) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Create the event structure
	eventID := generateEventID(eventType, tokenID)
	event := &SYN722Event{
		ID:        eventID,
		EventType: eventType,
		TokenID:   tokenID,
		Initiator: initiator,
		Timestamp: time.Now(),
		Details:   details,
		Encrypted: encrypt,
	}

	// Encrypt event details if requested
	if encrypt {
		encryptedData, encryptionKey, err := em.EncryptionService.EncryptData([]byte(common.MapToString(details)))
		if err != nil {
			return errors.New("failed to encrypt event data")
		}
		event.EncryptedData = encryptedData
		event.EncryptionKey = encryptionKey
	}

	// Validate the event through Synnergy Consensus
	if err := em.ConsensusEngine.ValidateEvent(event); err != nil {
		return errors.New("event validation failed via Synnergy Consensus")
	}

	// Record the event in the ledger
	if err := em.Ledger.RecordEvent(event); err != nil {
		return errors.New("failed to log event in the ledger")
	}

	return nil
}

// RetrieveEvent retrieves a specific event by its ID, decrypting it if necessary
func (em *SYN722EventManager) RetrieveEvent(eventID string) (*SYN722Event, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Retrieve the event from the ledger
	event, err := em.Ledger.GetEvent(eventID)
	if err != nil {
		return nil, errors.New("failed to retrieve event from ledger")
	}

	// If the event is encrypted, decrypt the data
	if event.Encrypted {
		decryptedData, err := em.EncryptionService.DecryptData([]byte(event.EncryptedData), event.EncryptionKey)
		if err != nil {
			return nil, errors.New("failed to decrypt event data")
		}

		// Convert decrypted data back to the original details map
		event.Details = common.StringToMap(string(decryptedData))
	}

	return event, nil
}

// ListEventsByToken retrieves all events related to a specific token
func (em *SYN722EventManager) ListEventsByToken(tokenID string) ([]*SYN722Event, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Fetch all events for the given tokenID from the ledger
	events, err := em.Ledger.GetEventsByToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve events for token from ledger")
	}

	// Decrypt any encrypted events
	for _, event := range events {
		if event.Encrypted {
			decryptedData, err := em.EncryptionService.DecryptData([]byte(event.EncryptedData), event.EncryptionKey)
			if err != nil {
				return nil, errors.New("failed to decrypt event data for event ID: " + event.ID)
			}
			event.Details = common.StringToMap(string(decryptedData))
		}
	}

	return events, nil
}

// generateEventID creates a unique event ID based on the event type and token ID
func generateEventID(eventType, tokenID string) string {
	return eventType + "_" + tokenID + "_" + time.Now().Format("20060102150405")
}

