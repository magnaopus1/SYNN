package syn11

import (
	"errors"
	"sync"
	"time"
  	"synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// Event Types for Syn11 tokens
const (
	EventIssuance    = "Issuance"     // Token issuance event
	EventTransfer    = "Transfer"     // Token transfer event
	EventBurn        = "Burn"         // Token burning event
	EventRedemption  = "Redemption"   // Token redemption event
	EventValidation  = "Validation"   // Transaction validation event
	EventCompliance  = "Compliance"   // Compliance (KYC/AML) event
	EventAudit       = "Audit"        // Audit-related event
	EventOwnership   = "Ownership"    // Token ownership transfer event
)

// Syn11Event represents an event associated with the Syn11 token.
type Syn11Event struct {
	EventID        string    // Unique identifier for the event
	TokenID        string    // Associated token ID
	EventType      string    // Type of event (e.g., Issuance, Transfer, Redemption)
	Initiator      string    // Address of the event initiator
	Recipient      string    // Address of the event recipient (if applicable)
	Amount         uint64    // Amount of tokens involved
	Timestamp      time.Time // Time the event occurred
	EncryptedData  string    // Encrypted event data
	Signature      string    // Digital signature for the event
}

// EventManager handles the logging, validation, and encryption of Syn11 token events.
type EventManager struct {
	mutex      sync.Mutex                  // Ensures thread-safe operations
	Ledger     *ledger.Ledger              // Reference to the ledger for storing events
	Consensus  *consensus.SynnergyConsensus // Reference to Synnergy Consensus for event validation
	Encryption *encryption.EncryptionService // Encryption service for secure event logging
	Events     []Syn11Event                // List of all events for the token
}

// NewEventManager creates a new EventManager instance.
func NewEventManager(ledgerInstance *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *EventManager {
	return &EventManager{
		Ledger:     ledgerInstance,
		Consensus:  consensusEngine,
		Encryption: encryptionService,
		Events:     []Syn11Event{},
	}
}

// LogEvent logs a new event for the Syn11 token.
func (em *EventManager) LogEvent(tokenID, eventType, initiator, recipient string, amount uint64) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Generate a unique event ID
	eventID := common.GenerateUniqueID()

	// Create the event
	event := Syn11Event{
		EventID:    eventID,
		TokenID:    tokenID,
		EventType:  eventType,
		Initiator:  initiator,
		Recipient:  recipient,
		Amount:     amount,
		Timestamp:  time.Now(),
	}

	// Encrypt the event data
	encryptedData, err := em.Encryption.Encrypt([]byte(event.EventType + event.Initiator + event.Recipient))
	if err != nil {
		return err
	}
	event.EncryptedData = string(encryptedData)

	// Generate a digital signature for the event
	event.Signature = common.GenerateSignature(eventID)

	// Validate the event through Synnergy Consensus
	if err := em.Consensus.ValidateEvent(event.EventID, event.TokenID, event.Initiator, event.Recipient, event.Amount); err != nil {
		return err
	}

	// Store the event in the ledger
	if err := em.Ledger.StoreEvent(event); err != nil {
		return err
	}

	// Append the event to the list
	em.Events = append(em.Events, event)

	return nil
}

// GetEvents returns the list of events for a specific token ID.
func (em *EventManager) GetEvents(tokenID string) ([]Syn11Event, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	var tokenEvents []Syn11Event
	for _, event := range em.Events {
		if event.TokenID == tokenID {
			tokenEvents = append(tokenEvents, event)
		}
	}

	if len(tokenEvents) == 0 {
		return nil, errors.New("no events found for the specified token ID")
	}

	return tokenEvents, nil
}

// ExportEvents exports all events for a given token, decrypting the details for external use.
func (em *EventManager) ExportEvents(tokenID string) ([]Syn11Event, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	events, err := em.GetEvents(tokenID)
	if err != nil {
		return nil, err
	}

	// Decrypt the event details before exporting
	for i, event := range events {
		decryptedData, err := em.Encryption.Decrypt([]byte(event.EncryptedData))
		if err != nil {
			return nil, err
		}
		events[i].EncryptedData = string(decryptedData)
	}

	return events, nil
}

// ValidateEvent verifies the authenticity and integrity of an event.
func (em *EventManager) ValidateEvent(eventID string) (bool, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	for _, event := range em.Events {
		if event.EventID == eventID {
			// Verify the event's digital signature
			if !common.VerifySignature(event.Signature, eventID) {
				return false, errors.New("invalid event signature")
			}
			return true, nil
		}
	}
	return false, errors.New("event not found")
}
