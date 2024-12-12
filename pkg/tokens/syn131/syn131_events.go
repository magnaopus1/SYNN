package syn131

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)



// EventListener defines an interface for listening to events.
type EventListener interface {
	OnEvent(event *Event)
}

// NewEventManager creates a new instance of EventManager.
func NewEventManager(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus) *EventManager {
	return &EventManager{
		events:          make(map[string]*Event),
		eventListeners:  make(map[string][]EventListener),
		ledger:          ledger,
		consensusEngine: consensusEngine,
	}
}

// CreateEvent creates and stores a new event in the system.
func (em *EventManager) CreateEvent(eventType, description, source string, payload map[string]interface{}) (*Event, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	eventID := generateUniqueID()
	event := &Event{
		ID:          eventID,
		Type:        eventType,
		Description: description,
		Source:      source,
		Timestamp:   time.Now(),
		Payload:     payload,
	}

	em.events[eventID] = event

	// Trigger consensus validation if necessary
	err := em.consensusEngine.ValidateEvent(event)
	if err != nil {
		return nil, err
	}

	// Store the event in the ledger
	err = em.ledger.RecordEvent(event)
	if err != nil {
		return nil, err
	}

	// Notify all registered listeners
	em.notifyListeners(event)

	return event, nil
}

// notifyListeners sends an event to all registered listeners.
func (em *EventManager) notifyListeners(event *Event) {
	listeners := em.eventListeners[event.Type]
	for _, listener := range listeners {
		go listener.OnEvent(event)
	}
}

// AddEventListener registers a new event listener for a specific event type.
func (em *EventManager) AddEventListener(eventType string, listener EventListener) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	em.eventListeners[eventType] = append(em.eventListeners[eventType], listener)
}

// RemoveEventListener removes an event listener from the system.
func (em *EventManager) RemoveEventListener(eventType string, listener EventListener) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	listeners, exists := em.eventListeners[eventType]
	if !exists {
		return errors.New("event type not found")
	}

	for i, l := range listeners {
		if l == listener {
			em.eventListeners[eventType] = append(listeners[:i], listeners[i+1:]...)
			return nil
		}
	}
	return errors.New("listener not found")
}

// TriggerOwnershipChangeEvent triggers an event for ownership changes.
func (em *EventManager) TriggerOwnershipChangeEvent(tokenID, previousOwner, newOwner string) (*Event, error) {
	payload := map[string]interface{}{
		"token_id":       tokenID,
		"previous_owner": previousOwner,
		"new_owner":      newOwner,
	}

	return em.CreateEvent("OwnershipChange", "Ownership of token changed", previousOwner, payload)
}

// TriggerPaymentEvent triggers an event for payments (rental, lease, purchase).
func (em *EventManager) TriggerPaymentEvent(tokenID, payer, payee string, amount float64, paymentType string) (*Event, error) {
	payload := map[string]interface{}{
		"token_id":     tokenID,
		"payer":        payer,
		"payee":        payee,
		"amount":       amount,
		"payment_type": paymentType,
	}

	return em.CreateEvent("Payment", paymentType+" payment processed", payer, payload)
}

// TriggerComplianceEvent triggers an event for compliance-related activities.
func (em *EventManager) TriggerComplianceEvent(tokenID, complianceStatus string) (*Event, error) {
	payload := map[string]interface{}{
		"token_id":         tokenID,
		"compliance_status": complianceStatus,
	}

	return em.CreateEvent("Compliance", "Compliance status updated", tokenID, payload)
}

// TriggerShardedOwnershipEvent triggers an event for changes in sharded ownership.
func (em *EventManager) TriggerShardedOwnershipEvent(tokenID string, shardDetails map[string]float64) (*Event, error) {
	payload := map[string]interface{}{
		"token_id":      tokenID,
		"shard_details": shardDetails,
	}

	return em.CreateEvent("ShardedOwnership", "Sharded ownership changed", tokenID, payload)
}

// StoreEvent stores an event securely with encryption and consensus validation.
func (em *EventManager) StoreEvent(event *Event) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Encrypt event data
	encryptedPayload, encryptionKey, err := encryption.EncryptData(event.Payload)
	if err != nil {
		return err
	}

	event.Payload["encrypted_data"] = encryptedPayload
	event.Payload["encryption_key"] = encryptionKey

	// Validate event using consensus
	err = em.consensusEngine.ValidateEvent(event)
	if err != nil {
		return err
	}

	// Store event in ledger
	err = em.ledger.RecordEvent(event)
	if err != nil {
		return err
	}

	return nil
}

// LoadEvent retrieves and decrypts an event from the ledger.
func (em *EventManager) LoadEvent(eventID string) (*Event, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Retrieve the event from the ledger
	event, err := em.ledger.GetEvent(eventID)
	if err != nil {
		return nil, err
	}

	// Decrypt the event data
	encryptionKey := event.Payload["encryption_key"].(string)
	encryptedData := event.Payload["encrypted_data"].([]byte)

	decryptedPayload, err := encryption.DecryptData(encryptedData, encryptionKey)
	if err != nil {
		return nil, err
	}

	// Replace encrypted payload with decrypted data
	event.Payload = decryptedPayload

	return event, nil
}

// generateUniqueID generates a unique event ID.
func generateUniqueID() string {
	return "evt_" + time.Now().Format("20060102150405") + "_" + generateRandomString(8)
}

// generateRandomString generates a random string of n length.
func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(bytes)
}
