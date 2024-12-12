package syn1500

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// SYN1500Events manages the lifecycle events of the SYN1500 reputation tokens.
type SYN1500Events struct {
	Ledger ledger.Ledger // Reference to the ledger for logging token events
}

// ReputationEvent represents an event related to SYN1500Token (e.g., endorsement, penalty, or transfer).
type ReputationEvent struct {
	EventID       string              // Unique identifier for the event
	TokenID       string              // The token involved in the event
	EventType     string              // Type of event (e.g., "Endorsement", "Penalty", "Transfer")
	Description   string              // Detailed description of the event
	PerformedBy   string              // ID of the user or system that triggered the event
	Timestamp     time.Time           // Timestamp of when the event occurred
	ImpactOnScore float64             // Impact (positive or negative) on the reputation score
	EncryptedData []byte              // Encrypted metadata related to the event
}

// RecordEvent logs an event into the ledger and applies the Synnergy Consensus mechanism.
func (se *SYN1500Events) RecordEvent(token *common.SYN1500Token, eventType, description, performedBy string, impactOnScore float64) error {
	// Create a new event and generate its unique ID
	eventID := generateUniqueID(token.TokenID)
	event := ReputationEvent{
		EventID:       eventID,
		TokenID:       token.TokenID,
		EventType:     eventType,
		Description:   description,
		PerformedBy:   performedBy,
		Timestamp:     time.Now(),
		ImpactOnScore: impactOnScore,
		EncryptedData: encryptEventDetails(description, performedBy, impactOnScore),
	}

	// Update the token's reputation score based on the event
	token.ReputationScore += impactOnScore

	// Append the event to the token's event logs
	token.ReputationEvents = append(token.ReputationEvents, event)

	// Record the event in the ledger with Synnergy Consensus validation
	tx := ledger.Transaction{
		TxID:        event.EventID,
		Description: fmt.Sprintf("Reputation event '%s' for token %s by %s", eventType, token.TokenID, performedBy),
		Timestamp:   event.Timestamp,
		Data:        event,
	}

	// Validate transaction using Synnergy Consensus before recording it
	if err := synnergy_consensus.ValidateTransaction(tx); err != nil {
		return fmt.Errorf("failed Synnergy Consensus validation: %v", err)
	}

	// Record the event into the ledger
	if err := se.Ledger.RecordTransaction(tx); err != nil {
		return errors.New("failed to record event in the ledger")
	}

	// Encrypt the metadata for the token after the event
	token.EncryptedMetadata = encryptMetadata(token)

	return nil
}

// ValidateSubBlocks performs sub-block validation for reputation-related transactions, following Synnergy Consensus.
func (se *SYN1500Events) ValidateSubBlocks(subBlocks []ledger.SubBlock) error {
	// Iterate through sub-blocks for validation
	for _, subBlock := range subBlocks {
		// Validate the sub-block under Synnergy Consensus
		if err := synnergy_consensus.ValidateSubBlock(subBlock); err != nil {
			return fmt.Errorf("sub-block validation failed for sub-block %s: %v", subBlock.SubBlockID, err)
		}
	}

	// Sub-blocks successfully validated
	return nil
}

// encryptEventDetails encrypts sensitive information related to a reputation event.
func encryptEventDetails(description, performedBy string, impactOnScore float64) []byte {
	data := fmt.Sprintf("Description: %s, PerformedBy: %s, ImpactOnScore: %.2f", description, performedBy, impactOnScore)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// encryptMetadata encrypts sensitive metadata for the SYN1500Token.
func encryptMetadata(token *common.SYN1500Token) []byte {
	data, _ := json.Marshal(token)
	hash := sha256.Sum256(data)
	return hash[:]
}

// generateUniqueID creates a unique identifier for events and logs.
func generateUniqueID(seed string) string {
	timestamp := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", seed, timestamp)))
	return hex.EncodeToString(hash[:])
}

// HandleEndorsement processes an endorsement event for the SYN1500Token.
func (se *SYN1500Events) HandleEndorsement(token *common.SYN1500Token, performedBy string, rating float64, review string) error {
	eventType := "Endorsement"
	description := fmt.Sprintf("Endorsement from %s: %s", performedBy, review)
	impactOnScore := rating * 0.1 // Impact on reputation score based on rating (example logic)

	return se.RecordEvent(token, eventType, description, performedBy, impactOnScore)
}

// HandlePenalty processes a penalty event for the SYN1500Token.
func (se *SYN1500Events) HandlePenalty(token *common.SYN1500Token, performedBy string, reason string, penaltyScore float64) error {
	eventType := "Penalty"
	description := fmt.Sprintf("Penalty from %s: %s", performedBy, reason)
	impactOnScore := -penaltyScore // Negative impact on reputation score

	return se.RecordEvent(token, eventType, description, performedBy, impactOnScore)
}

// HandleTransfer logs a transfer event for the SYN1500Token.
func (se *SYN1500Events) HandleTransfer(token *common.SYN1500Token, newOwner, performedBy string) error {
	eventType := "Transfer"
	description := fmt.Sprintf("Token transferred to %s by %s", newOwner, performedBy)
	impactOnScore := 0.0 // No impact on reputation score for a transfer

	// Update the token owner
	token.Owner = newOwner

	return se.RecordEvent(token, eventType, description, performedBy, impactOnScore)
}
