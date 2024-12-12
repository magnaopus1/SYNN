package syn1800

import (
	"time"
	"fmt"
)

// EventManager handles the events related to SYN1800 tokens, including emissions, offsets, rewards, and verifications.
type EventManager struct {
	ledger *ledger.Ledger // Ledger integration for event tracking and updates
}

// NewEventManager initializes a new EventManager.
func NewEventManager(ledger *ledger.Ledger) *EventManager {
	return &EventManager{ledger: ledger}
}

// LogEmissionEvent logs an emission event to the ledger, adding it to the SYN1800 token's event history.
func (em *EventManager) LogEmissionEvent(tokenID string, amount float64, description string, verifiedBy string) error {
	// Retrieve the token from the ledger
	token, err := em.ledger.GetTokenByID(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token: %v", err)
	}

	// Cast to SYN1800Token type
	syn1800Token, ok := token.(*common.SYN1800Token)
	if !ok {
		return fmt.Errorf("invalid token type")
	}

	// Create a new emission event
	emissionEvent := common.EmissionOffsetLog{
		LogID:        generateUniqueID(),
		ActivityType: "Emission",
		Amount:       amount,
		ActivityDate: time.Now(),
		Description:  description,
		VerifiedBy:   verifiedBy,
	}

	// Add the event to the token's log
	syn1800Token.CarbonFootprintLogs = append(syn1800Token.CarbonFootprintLogs, emissionEvent)

	// Update the ledger with the new emission event
	err = em.ledger.UpdateTokenInLedger(syn1800Token)
	if err != nil {
		return fmt.Errorf("failed to update ledger with emission event: %v", err)
	}

	return nil
}

// LogOffsetEvent logs an offset event to the ledger, adding it to the SYN1800 token's event history.
func (em *EventManager) LogOffsetEvent(tokenID string, amount float64, description string, verifiedBy string) error {
	// Retrieve the token from the ledger
	token, err := em.ledger.GetTokenByID(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token: %v", err)
	}

	// Cast to SYN1800Token type
	syn1800Token, ok := token.(*common.SYN1800Token)
	if !ok {
		return fmt.Errorf("invalid token type")
	}

	// Create a new offset event
	offsetEvent := common.EmissionOffsetLog{
		LogID:        generateUniqueID(),
		ActivityType: "Offset",
		Amount:       amount,
		ActivityDate: time.Now(),
		Description:  description,
		VerifiedBy:   verifiedBy,
	}

	// Add the event to the token's log
	syn1800Token.CarbonFootprintLogs = append(syn1800Token.CarbonFootprintLogs, offsetEvent)

	// Update the ledger with the new offset event
	err = em.ledger.UpdateTokenInLedger(syn1800Token)
	if err != nil {
		return fmt.Errorf("failed to update ledger with offset event: %v", err)
	}

	return nil
}

// LogRewardEvent logs a reward event to the ledger, incentivizing a carbon reduction activity.
func (em *EventManager) LogRewardEvent(tokenID string, recipientID string, amount float64, activity string, verifiedBy string) error {
	// Retrieve the token from the ledger
	token, err := em.ledger.GetTokenByID(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token: %v", err)
	}

	// Cast to SYN1800Token type
	syn1800Token, ok := token.(*common.SYN1800Token)
	if !ok {
		return fmt.Errorf("invalid token type")
	}

	// Create a new reward event
	rewardEvent := common.RewardLog{
		RewardID:     generateUniqueID(),
		RecipientID:  recipientID,
		RewardAmount: amount,
		Activity:     activity,
		RewardDate:   time.Now(),
		Verification: verifiedBy,
	}

	// Add the reward to the token's event log
	syn1800Token.RewardRecords = append(syn1800Token.RewardRecords, rewardEvent)

	// Update the ledger with the new reward event
	err = em.ledger.UpdateTokenInLedger(syn1800Token)
	if err != nil {
		return fmt.Errorf("failed to update ledger with reward event: %v", err)
	}

	return nil
}

// LogVerificationEvent logs a verification event for emissions or offsets, ensuring regulatory compliance.
func (em *EventManager) LogVerificationEvent(tokenID string, source string, verifiedAmount float64, verificationStatus string) error {
	// Retrieve the token from the ledger
	token, err := em.ledger.GetTokenByID(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token: %v", err)
	}

	// Cast to SYN1800Token type
	syn1800Token, ok := token.(*common.SYN1800Token)
	if !ok {
		return fmt.Errorf("invalid token type")
	}

	// Create a new verification event
	verificationEvent := common.VerificationLog{
		VerificationID: generateUniqueID(),
		Source:         source,
		VerificationDate: time.Now(),
		Description:    fmt.Sprintf("Verification for source %s with amount %f", source, verifiedAmount),
		VerifiedAmount: verifiedAmount,
		Status:         verificationStatus,
	}

	// Add the verification event to the token's verification log
	syn1800Token.SourceVerificationLog = append(syn1800Token.SourceVerificationLog, verificationEvent)

	// Update the ledger with the new verification event
	err = em.ledger.UpdateTokenInLedger(syn1800Token)
	if err != nil {
		return fmt.Errorf("failed to update ledger with verification event: %v", err)
	}

	return nil
}

// encryptEventMetadata encrypts event metadata for sensitive activities.
func encryptEventMetadata(token *common.SYN1800Token) ([]byte, error) {
	// Placeholder encryption logic. Replace with your real encryption implementation.
	return crypto.Encrypt([]byte(fmt.Sprintf("%v", token)), "encryption-key")
}

// generateUniqueID generates a unique ID for events, verifications, and rewards.
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
