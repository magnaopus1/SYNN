package syn4300

import (
	"errors"
	"sync"
	"time"
)

// EventManager manages the events associated with SYN4300 tokens.
type EventManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// NewEventManager initializes a new event manager.
func NewEventManager(ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *EventManager {
	return &EventManager{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// LogTokenCreationEvent logs an event when a SYN4300 token is created.
func (em *EventManager) LogTokenCreationEvent(tokenID string, owner string, metadata Syn4300Metadata) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Create event details
	eventDetails := "Token created. Owner: " + owner + ", Token ID: " + tokenID

	// Log event in the ledger
	if err := em.ledgerService.LogEvent("TokenCreation", time.Now(), tokenID, eventDetails); err != nil {
		return err
	}

	// Validate the event using Synnergy Consensus
	if err := em.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// LogTokenTransferEvent logs an event when a SYN4300 token is transferred.
func (em *EventManager) LogTokenTransferEvent(tokenID string, fromOwner string, toOwner string, quantity float64) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Create event details
	eventDetails := "Token transferred. From: " + fromOwner + " to " + toOwner + ". Quantity: " + fmt.Sprintf("%f", quantity)

	// Log event in the ledger
	if err := em.ledgerService.LogEvent("TokenTransfer", time.Now(), tokenID, eventDetails); err != nil {
		return err
	}

	// Validate the event using Synnergy Consensus
	if err := em.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// LogTradeEvent logs an event when a trade involving a SYN4300 token is executed.
func (em *EventManager) LogTradeEvent(tokenID string, tradeDetails string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Log event in the ledger
	if err := em.ledgerService.LogEvent("TradeExecuted", time.Now(), tokenID, tradeDetails); err != nil {
		return err
	}

	// Validate the trade event using Synnergy Consensus
	if err := em.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// LogTokenStatusChangeEvent logs an event when the status of a SYN4300 token is updated.
func (em *EventManager) LogTokenStatusChangeEvent(tokenID string, newStatus string, reason string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Create event details
	eventDetails := "Token status changed. Token ID: " + tokenID + ", New Status: " + newStatus + ". Reason: " + reason

	// Log event in the ledger
	if err := em.ledgerService.LogEvent("TokenStatusChange", time.Now(), tokenID, eventDetails); err != nil {
		return err
	}

	// Validate the status change event using Synnergy Consensus
	if err := em.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// LogEnergyProductionEvent logs an event when energy production is recorded for a SYN4300 token.
func (em *EventManager) LogEnergyProductionEvent(tokenID string, energyProduced float64, unit string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Create event details
	eventDetails := fmt.Sprintf("Energy produced. Token ID: %s, Amount: %.2f %s", tokenID, energyProduced, unit)

	// Log event in the ledger
	if err := em.ledgerService.LogEvent("EnergyProduction", time.Now(), tokenID, eventDetails); err != nil {
		return err
	}

	// Validate the energy production event using Synnergy Consensus
	if err := em.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// LogEnergyCertificationUpdateEvent logs an event when the certification of a SYN4300 token is updated.
func (em *EventManager) LogEnergyCertificationUpdateEvent(tokenID string, certificationDetails string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Log event in the ledger
	if err := em.ledgerService.LogEvent("CertificationUpdate", time.Now(), tokenID, certificationDetails); err != nil {
		return err
	}

	// Validate the certification update using Synnergy Consensus
	if err := em.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}
