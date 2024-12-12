package syn1301

import (
	"errors"
	"time"

)

// SYN1301EventManager manages events and logging for SYN1301 tokens in the supply chain.
type SYN1301EventManager struct {
	Ledger            *ledger.Ledger                // Ledger instance for managing event-related records
	EncryptionService *encryption.EncryptionService // Encryption service for securing event data
	Consensus         *synnergy_consensus.Consensus // Synnergy Consensus system for event validation
}

// LogEvent records an event for a SYN1301 token in the ledger.
func (em *SYN1301EventManager) LogEvent(tokenID string, eventType string, eventDetails map[string]string, userID string) error {
	// Step 1: Encrypt the event details
	encryptedDetails, err := em.EncryptionService.Encrypt(eventDetails)
	if err != nil {
		return errors.New("failed to encrypt event details: " + err.Error())
	}

	// Step 2: Create the event log structure
	eventLog := ledger.EventLog{
		EventType:   eventType,
		TokenID:     tokenID,
		UserID:      userID,
		Description: "Event: " + eventType,
		Metadata:    encryptedDetails,
		Timestamp:   time.Now(),
	}

	// Step 3: Record the event in the ledger
	err = em.Ledger.LogEvent(eventLog)
	if err != nil {
		return errors.New("failed to record event in ledger: " + err.Error())
	}

	// Step 4: Optionally validate the event in Synnergy Consensus (for critical events)
	if eventType == "CRITICAL_UPDATE" {
		err = em.Consensus.ValidateEvent(eventLog)
		if err != nil {
			return errors.New("failed to validate event in consensus: " + err.Error())
		}
	}

	return nil
}

// ValidateEventTransaction processes event-related transactions through Synnergy Consensus.
func (em *SYN1301EventManager) ValidateEventTransaction(tokenID string, eventDetails map[string]string) (bool, error) {
	// Step 1: Validate transaction into sub-block via Synnergy Consensus
	subBlock, err := em.Consensus.ValidateTransactionIntoSubBlock(tokenID, eventDetails)
	if err != nil {
		return false, errors.New("failed to validate transaction into sub-block: " + err.Error())
	}

	// Step 2: Validate sub-block into a full block
	_, err = em.Consensus.ValidateSubBlockIntoBlock(subBlock)
	if err != nil {
		return false, errors.New("failed to validate sub-block into block: " + err.Error())
	}

	// Transaction and event validated
	return true, nil
}

// RecordEvent logs important actions or transitions for supply chain tokens in the event log.
func (em *SYN1301EventManager) RecordEvent(tokenID string, eventType string, metadata map[string]string, userID string) error {
	// Step 1: Encrypt the event data
	encryptedMetadata, err := em.EncryptionService.Encrypt(metadata)
	if err != nil {
		return errors.New("failed to encrypt event metadata: " + err.Error())
	}

	// Step 2: Create event log entry
	event := ledger.EventLog{
		EventType:   eventType,
		TokenID:     tokenID,
		UserID:      userID,
		Description: "Event logged: " + eventType,
		Metadata:    encryptedMetadata,
		Timestamp:   time.Now(),
	}

	// Step 3: Log the event in the ledger
	err = em.Ledger.LogEvent(event)
	if err != nil {
		return errors.New("failed to record event in ledger: " + err.Error())
	}

	// Step 4: Trigger validation through consensus for critical events
	if eventType == "CRITICAL_EVENT" {
		err = em.Consensus.ValidateEvent(event)
		if err != nil {
			return errors.New("failed to validate event via consensus: " + err.Error())
		}
	}

	return nil
}

// QueryEventHistory retrieves the event history for a SYN1301 token.
func (em *SYN1301EventManager) QueryEventHistory(tokenID string) ([]ledger.EventLog, error) {
	// Step 1: Retrieve the event logs from the ledger for the given token ID
	eventLogs, err := em.Ledger.GetEventsByToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve event history from ledger: " + err.Error())
	}

	// Step 2: Return the event logs
	return eventLogs, nil
}

// UpdateEvent records a change or update to an event log entry.
func (em *SYN1301EventManager) UpdateEvent(eventID string, newDetails map[string]string, userID string) error {
	// Step 1: Encrypt the new event details
	encryptedDetails, err := em.EncryptionService.Encrypt(newDetails)
	if err != nil {
		return errors.New("failed to encrypt new event details: " + err.Error())
	}

	// Step 2: Retrieve the event from the ledger
	event, err := em.Ledger.GetEvent(eventID)
	if err != nil {
		return errors.New("failed to retrieve event for update: " + err.Error())
	}

	// Step 3: Update the event details and re-encrypt the metadata
	event.Metadata = encryptedDetails
	event.Timestamp = time.Now()

	// Step 4: Update the event in the ledger
	err = em.Ledger.UpdateEvent(eventID, event)
	if err != nil {
		return errors.New("failed to update event in ledger: " + err.Error())
	}

	// Step 5: Log the update as a new event in the ledger
	err = em.Ledger.LogEvent(ledger.EventLog{
		EventType:   "EVENT_UPDATED",
		TokenID:     event.TokenID,
		UserID:      userID,
		Description: "Event updated: " + eventID,
		Metadata:    encryptedDetails,
		Timestamp:   time.Now(),
	})
	if err != nil {
		return errors.New("failed to log event update: " + err.Error())
	}

	return nil
}

// MarkEventAsCritical flags an event as critical for validation through consensus.
func (em *SYN1301EventManager) MarkEventAsCritical(eventID string, userID string) error {
	// Step 1: Retrieve the event from the ledger
	event, err := em.Ledger.GetEvent(eventID)
	if err != nil {
		return errors.New("failed to retrieve event: " + err.Error())
	}

	// Step 2: Flag the event as critical and update the ledger
	event.EventType = "CRITICAL_EVENT"
	event.Timestamp = time.Now()

	err = em.Ledger.UpdateEvent(eventID, event)
	if err != nil {
		return errors.New("failed to flag event as critical: " + err.Error())
	}

	// Step 3: Trigger validation through Synnergy Consensus
	err = em.Consensus.ValidateEvent(event)
	if err != nil {
		return errors.New("failed to validate critical event: " + err.Error())
	}

	// Step 4: Log the critical event marking in the ledger
	err = em.Ledger.LogEvent(ledger.EventLog{
		EventType:   "EVENT_MARKED_CRITICAL",
		TokenID:     event.TokenID,
		UserID:      userID,
		Description: "Event marked as critical: " + eventID,
		Timestamp:   time.Now(),
	})
	if err != nil {
		return errors.New("failed to log event marking as critical: " + err.Error())
	}

	return nil
}
