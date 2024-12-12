package syn10

import (
    "encoding/json"
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NewEventManager creates a new EventManager instance with references to the ledger and consensus engine.
func NewEventManager(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption) *SYN10EventManager {
    return &SYN10EventManager{
        Ledger:     ledgerInstance,
        Consensus:  consensusEngine,
        Encryption: encryptionService,
        Events:     make([]SYN10Event, 0),
    }
}

// LogEvent logs and validates a new event. This includes encrypting the event details and updating the ledger.
func (em *SYN10EventManager) LogEvent(eventType, tokenID, from, to string, amount uint64, exchangeRate float64, details string) error {
    em.mutex.Lock()
    defer em.mutex.Unlock()

    // Create the event struct
    event := SYN10Event{
        EventType:    eventType,
        Timestamp:    time.Now(),
        TokenID:      tokenID,
        FromAddress:  from,
        ToAddress:    to,
        Amount:       amount,
        ExchangeRate: exchangeRate,
        Details:      details,
        Encrypted:    false,
    }

    // Encrypt event details
    encryptedDetails, err := em.Encryption.EncryptData(details, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("error encrypting event details: %v", err)
    }
    event.Details = encryptedDetails
    event.Encrypted = true

    // Add event to the event log
    em.Events = append(em.Events, event)

    // Validate the event using Synnergy Consensus
    if valid, err := em.Consensus.ValidateEvent(event); !valid || err != nil {
        return fmt.Errorf("error validating event: %v", err)
    }

    // Update the ledger with the event
    eventData, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("error marshalling event data: %v", err)
    }
    if err := em.Ledger.AddEvent(event.TokenID, eventData); err != nil {
        return fmt.Errorf("error storing event in ledger: %v", err)
    }

    return nil
}

// GetEvents retrieves all logged events for the token.
func (em *SYN10EventManager) GetEvents(tokenID string) ([]SYN10Event, error) {
    var tokenEvents []SYN10Event
    for _, event := range em.Events {
        if event.TokenID == tokenID {
            tokenEvents = append(tokenEvents, event)
        }
    }

    if len(tokenEvents) == 0 {
        return nil, fmt.Errorf("no events found for token ID %s", tokenID)
    }

    return tokenEvents, nil
}
