package syn721

import (
	"errors"
	"fmt"
	"sync"
)

// EventLogger manages and logs the events related to SYN721 tokens.
type EventLogger struct {
	mutex      sync.Mutex                 // For thread-safe logging
	Ledger     *ledger.Ledger             // Reference to the ledger for storing event logs
	Consensus  *synnergy_consensus.Engine // Synnergy Consensus engine for validation
	Encryption *encryption.Encryption     // Encryption service for securing event data
}

// NewEventLogger initializes a new event logger for SYN721 tokens.
func NewEventLogger(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *EventLogger {
	return &EventLogger{
		Ledger:     ledgerInstance,
		Consensus:  consensusEngine,
		Encryption: encryptionService,
	}
}

// LogMintEvent logs an event when an SYN721 token is minted.
func (el *EventLogger) LogMintEvent(tokenID, owner, tokenURI string) error {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	// Create event message
	eventMessage := fmt.Sprintf("Minted token: %s, Owner: %s, TokenURI: %s", tokenID, owner, tokenURI)

	// Encrypt the event message
	encryptedEvent, err := el.Encryption.EncryptData(eventMessage, "")
	if err != nil {
		return fmt.Errorf("error encrypting mint event: %v", err)
	}

	// Record the event in the ledger
	err = el.Ledger.RecordEvent(tokenID, owner, "Mint", encryptedEvent, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log mint event in ledger: %v", err)
	}

	fmt.Printf("Mint event logged for token %s (owner: %s).\n", tokenID, owner)
	return nil
}

// LogTransferEvent logs an event when an SYN721 token is transferred.
func (el *EventLogger) LogTransferEvent(tokenID, from, to string) error {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	// Create event message
	eventMessage := fmt.Sprintf("Transferred token: %s, From: %s, To: %s", tokenID, from, to)

	// Encrypt the event message
	encryptedEvent, err := el.Encryption.EncryptData(eventMessage, "")
	if err != nil {
		return fmt.Errorf("error encrypting transfer event: %v", err)
	}

	// Record the event in the ledger
	err = el.Ledger.RecordEvent(tokenID, from, "Transfer", encryptedEvent, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transfer event in ledger: %v", err)
	}

	fmt.Printf("Transfer event logged for token %s (from: %s to: %s).\n", tokenID, from, to)
	return nil
}

// LogApprovalEvent logs an event when a token is approved for transfer by another address.
func (el *EventLogger) LogApprovalEvent(tokenID, owner, approved string) error {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	// Create event message
	eventMessage := fmt.Sprintf("Approved transfer of token: %s, Owner: %s, Approved: %s", tokenID, owner, approved)

	// Encrypt the event message
	encryptedEvent, err := el.Encryption.EncryptData(eventMessage, "")
	if err != nil {
		return fmt.Errorf("error encrypting approval event: %v", err)
	}

	// Record the event in the ledger
	err = el.Ledger.RecordEvent(tokenID, owner, "Approval", encryptedEvent, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log approval event in ledger: %v", err)
	}

	fmt.Printf("Approval event logged for token %s (owner: %s, approved: %s).\n", tokenID, owner, approved)
	return nil
}

// LogBurnEvent logs an event when an SYN721 token is burned.
func (el *EventLogger) LogBurnEvent(tokenID, owner string) error {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	// Create event message
	eventMessage := fmt.Sprintf("Burned token: %s, Owner: %s", tokenID, owner)

	// Encrypt the event message
	encryptedEvent, err := el.Encryption.EncryptData(eventMessage, "")
	if err != nil {
		return fmt.Errorf("error encrypting burn event: %v", err)
	}

	// Record the event in the ledger
	err = el.Ledger.RecordEvent(tokenID, owner, "Burn", encryptedEvent, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log burn event in ledger: %v", err)
	}

	fmt.Printf("Burn event logged for token %s (owner: %s).\n", tokenID, owner)
	return nil
}

// LogCustomEvent allows logging of a custom event related to an SYN721 token.
func (el *EventLogger) LogCustomEvent(tokenID, owner, eventType, eventDescription string) error {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	// Create event message
	eventMessage := fmt.Sprintf("Event: %s, TokenID: %s, Owner: %s, Description: %s", eventType, tokenID, owner, eventDescription)

	// Encrypt the event message
	encryptedEvent, err := el.Encryption.EncryptData(eventMessage, "")
	if err != nil {
		return fmt.Errorf("error encrypting custom event: %v", err)
	}

	// Record the event in the ledger
	err = el.Ledger.RecordEvent(tokenID, owner, eventType, encryptedEvent, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log custom event in ledger: %v", err)
	}

	fmt.Printf("Custom event (%s) logged for token %s (owner: %s).\n", eventType, tokenID, owner)
	return nil
}

// GetEventLogs retrieves the logs related to a specific SYN721 token.
func (el *EventLogger) GetEventLogs(tokenID string) ([]string, error) {
	el.mutex.Lock()
	defer el.mutex.Unlock()

	// Retrieve events from the ledger
	eventLogs, err := el.Ledger.GetEventsByTokenID(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve event logs for token %s: %v", tokenID, err)
	}

	// Decrypt the event logs
	decryptedLogs := make([]string, 0, len(eventLogs))
	for _, encryptedLog := range eventLogs {
		decryptedLog, err := el.Encryption.DecryptData(encryptedLog, "")
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt event log for token %s: %v", tokenID, err)
		}
		decryptedLogs = append(decryptedLogs, decryptedLog)
	}

	return decryptedLogs, nil
}
