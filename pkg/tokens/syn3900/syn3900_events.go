package syn3900

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
	"sync"

)

// BenefitEventManager handles all the events related to the SYN3900 tokens, including allocations, claims, and expirations.
type BenefitEventManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// BenefitEvent represents a general structure for events tied to SYN3900 benefit tokens.
type BenefitEvent struct {
	EventID       string    `json:"event_id"`
	TokenID       string    `json:"token_id"`
	EventType     string    `json:"event_type"`    // e.g., "Allocation", "Claim", "Expiration"
	EventDetails  string    `json:"event_details"` // Description of the event
	Timestamp     time.Time `json:"timestamp"`     // Time when the event was logged
}

// NewBenefitEventManager initializes a new BenefitEventManager.
func NewBenefitEventManager(ledgerService *ledger.LedgerService, encryptionService *encryption.Encryptor, consensusService *consensus.SynnergyConsensus) *BenefitEventManager {
	return &BenefitEventManager{
		ledgerService:     ledgerService,
		encryptionService: encryptionService,
		consensusService:  consensusService,
	}
}

// LogAllocationEvent logs an allocation event for a SYN3900 benefit token.
func (bem *BenefitEventManager) LogAllocationEvent(tokenID string, recipient string, amount float64) error {
	bem.mutex.Lock()
	defer bem.mutex.Unlock()

	// Create allocation event
	event := BenefitEvent{
		EventID:      generateUniqueEventID(),
		TokenID:      tokenID,
		EventType:    "Allocation",
		EventDetails: "Allocated " + recipient + " an amount of " + fmt.Sprintf("%f", amount),
		Timestamp:    time.Now(),
	}

	// Store the event in the ledger
	if err := bem.logEventToLedger(event); err != nil {
		return err
	}

	// Validate the event in Synnergy Consensus
	if err := bem.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// LogClaimEvent logs a claim event for a SYN3900 benefit token.
func (bem *BenefitEventManager) LogClaimEvent(tokenID string, claimant string, amount float64) error {
	bem.mutex.Lock()
	defer bem.mutex.Unlock()

	// Create claim event
	event := BenefitEvent{
		EventID:      generateUniqueEventID(),
		TokenID:      tokenID,
		EventType:    "Claim",
		EventDetails: claimant + " claimed an amount of " + fmt.Sprintf("%f", amount),
		Timestamp:    time.Now(),
	}

	// Store the event in the ledger
	if err := bem.logEventToLedger(event); err != nil {
		return err
	}

	// Validate the event in Synnergy Consensus
	if err := bem.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// LogExpirationEvent logs an expiration event for a SYN3900 benefit token.
func (bem *BenefitEventManager) LogExpirationEvent(tokenID string) error {
	bem.mutex.Lock()
	defer bem.mutex.Unlock()

	// Create expiration event
	event := BenefitEvent{
		EventID:      generateUniqueEventID(),
		TokenID:      tokenID,
		EventType:    "Expiration",
		EventDetails: "Benefit token expired",
		Timestamp:    time.Now(),
	}

	// Store the event in the ledger
	if err := bem.logEventToLedger(event); err != nil {
		return err
	}

	// Validate the event in Synnergy Consensus
	if err := bem.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// LogConditionalReleaseEvent logs an event when funds are conditionally released based on predefined criteria.
func (bem *BenefitEventManager) LogConditionalReleaseEvent(tokenID string, releaseDetails string) error {
	bem.mutex.Lock()
	defer bem.mutex.Unlock()

	// Create conditional release event
	event := BenefitEvent{
		EventID:      generateUniqueEventID(),
		TokenID:      tokenID,
		EventType:    "ConditionalRelease",
		EventDetails: releaseDetails,
		Timestamp:    time.Now(),
	}

	// Store the event in the ledger
	if err := bem.logEventToLedger(event); err != nil {
		return err
	}

	// Validate the event in Synnergy Consensus
	if err := bem.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// logEventToLedger securely logs the benefit event to the ledger.
func (bem *BenefitEventManager) logEventToLedger(event BenefitEvent) error {
	// Encrypt the event data before storing
	encryptedEvent, err := bem.encryptionService.EncryptData(event)
	if err != nil {
		return err
	}

	// Log the encrypted event in the ledger
	if err := bem.ledgerService.StoreData(event.EventID, encryptedEvent); err != nil {
		return err
	}

	return nil
}

// generateUniqueEventID generates a unique identifier for benefit events.
func generateUniqueEventID() string {
	return "event-" + time.Now().Format("20060102150405")
}

// RetrieveEventHistory retrieves the full event history for a specific token.
func (bem *BenefitEventManager) RetrieveEventHistory(tokenID string) ([]BenefitEvent, error) {
	// Query the ledger for event data
	eventIDs, err := bem.ledgerService.RetrieveEventIDs(tokenID)
	if err != nil {
		return nil, err
	}

	// Decrypt and load events
	var events []BenefitEvent
	for _, eventID := range eventIDs {
		encryptedData, err := bem.ledgerService.RetrieveData(eventID)
		if err != nil {
			return nil, err
		}

		decryptedEvent, err := bem.encryptionService.DecryptData(encryptedData)
		if err != nil {
			return nil, err
		}

		events = append(events, decryptedEvent.(BenefitEvent))
	}

	return events, nil
}
