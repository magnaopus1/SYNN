package syn1401

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

// SYN1401EventManager handles event logging for SYN1401 tokens.
type SYN1401EventManager struct {
	Ledger common.LedgerInterface // Interface to interact with the ledger
}

// LogTokenCreation logs the event of SYN1401 token creation.
func (em *SYN1401EventManager) LogTokenCreation(tokenID string, creator string) error {
	event := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Creation",
		Description: fmt.Sprintf("SYN1401 Token %s was created by %s", tokenID, creator),
		EventDate:   time.Now(),
		PerformedBy: creator,
	}

	return em.logEvent(tokenID, event)
}

// LogTokenTransfer logs the event of a token being transferred from one owner to another.
func (em *SYN1401EventManager) LogTokenTransfer(tokenID, from, to string) error {
	event := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Transfer",
		Description: fmt.Sprintf("SYN1401 Token %s was transferred from %s to %s", tokenID, from, to),
		EventDate:   time.Now(),
		PerformedBy: from,
	}

	return em.logEvent(tokenID, event)
}

// LogInterestPayment logs the event of an interest payment on the SYN1401 token.
func (em *SYN1401EventManager) LogInterestPayment(tokenID string, paymentAmount float64, payer string) error {
	event := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Interest Payment",
		Description: fmt.Sprintf("Interest payment of %.2f was made for SYN1401 Token %s by %s", paymentAmount, tokenID, payer),
		EventDate:   time.Now(),
		PerformedBy: payer,
	}

	return em.logEvent(tokenID, event)
}

// LogTokenRedemption logs the redemption event of a SYN1401 token.
func (em *SYN1401EventManager) LogTokenRedemption(tokenID string, redemptionType string, redeemer string, principalPaid float64, interestPaid float64) error {
	event := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Redemption",
		Description: fmt.Sprintf("SYN1401 Token %s was redeemed (%s) by %s with principal: %.2f and interest: %.2f", tokenID, redemptionType, redeemer, principalPaid, interestPaid),
		EventDate:   time.Now(),
		PerformedBy: redeemer,
	}

	return em.logEvent(tokenID, event)
}

// LogCustomEvent allows logging custom events associated with a SYN1401 token.
func (em *SYN1401EventManager) LogCustomEvent(tokenID string, eventType string, description string, performedBy string) error {
	event := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   eventType,
		Description: description,
		EventDate:   time.Now(),
		PerformedBy: performedBy,
	}

	return em.logEvent(tokenID, event)
}

// logEvent is a helper function to log events into the ledger and token's event log.
func (em *SYN1401EventManager) logEvent(tokenID string, event common.EventLog) error {
	// Retrieve the token from the ledger
	token, err := em.Ledger.GetToken(tokenID)
	if err != nil {
		return fmt.Errorf("error retrieving token from ledger: %w", err)
	}

	// Append the event log to the token
	token.EventLogs = append(token.EventLogs, event)

	// Update the token in the ledger
	if err := em.Ledger.UpdateToken(tokenID, token); err != nil {
		return fmt.Errorf("error updating token with event log: %w", err)
	}

	return nil
}

// generateUniqueID generates a unique ID for events.
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
