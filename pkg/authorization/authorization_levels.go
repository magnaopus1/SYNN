package authorization

import (
	"fmt"
	"synnergy_network/pkg/ledger"
	"time"
)

// AuthorizationData represents the structure for user authorization data.
type AuthorizationData struct {
	UserID            string    // Unique identifier for the user
	AuthorizationLevel int       // Authorization level for the user
	SetAt             time.Time // Timestamp of when the authorization level was set
}

// TrustedParty represents a trusted party in the system.
type TrustedParty struct {
	PartyID string    // Unique identifier for the trusted party
	AddedAt time.Time // Timestamp of when the party was added as trusted
	Flagged bool      // Indicates if the party is flagged for review
}

// AuthorizationEvent represents an event related to authorization actions.
type AuthorizationEvent struct {
	EventID    string    // Unique identifier for the authorization event
	Action     string    // Type of action, e.g., "Added", "Removed", "Modified"
	UserID     string    // User ID associated with the event
	Timestamp  time.Time // Timestamp of the event
	Details    string    // Additional details about the event
}

// SetAuthorizationLevel sets an authorization level for a user, storing it in the ledger.
func SetAuthorizationLevel(ledgerInstance *ledger.Ledger, userID string, level int) error {
	// Create an instance of ledger.AuthorizationData to match the expected type
	authData := ledger.AuthorizationData{
		UserID:             userID,
		AuthorizationLevel: level,
		SetAt:              time.Now(),
	}

	// Record authorization level in the ledger
	if err := ledgerInstance.AuthorizationLedger.RecordAuthorizationLevel(authData); err != nil {
		return fmt.Errorf("failed to record authorization level in ledger: %v", err)
	}

	return nil
}

// GetAuthorizationLevel retrieves the authorization level for a given user from the ledger.
func GetAuthorizationLevel(ledgerInstance *ledger.Ledger, userID string) (int, error) {

	level, err := ledgerInstance.AuthorizationLedger.FetchAuthorizationLevel(userID)
	if err != nil {
		return 0, fmt.Errorf("error fetching authorization level from ledger: %v", err)
	}

	return level, nil
}

// VerifySignerIdentity confirms that the signer's identity matches the authorization level required.
func VerifySignerIdentity(ledgerInstance *ledger.Ledger, signerID string, requiredLevel int) (bool, error) {
	level, err := ledgerInstance.AuthorizationLedger.FetchAuthorizationLevel(signerID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve signer authorization level: %v", err)
	}
	return level >= requiredLevel, nil
}

// AddTrustedParty adds a trusted party to the system, logging the addition in the ledger.
func AddTrustedParty(ledgerInstance *ledger.Ledger, partyID string) error {
	// Create an instance of ledger.TrustedParty to match the expected type
	trustedParty := ledger.TrustedParty{
		PartyID: partyID,
		AddedAt: time.Now(),
		Flagged: false,
	}

	if err := ledgerInstance.AuthorizationLedger.RecordTrustedParty(trustedParty); err != nil {
		return fmt.Errorf("failed to add trusted party to ledger: %v", err)
	}

	return nil
}

// RemoveTrustedParty removes a trusted party from the ledger, revoking any authorization.
func RemoveTrustedParty(ledgerInstance *ledger.Ledger, partyID string) error {
	if err := ledgerInstance.AuthorizationLedger.DeleteTrustedParty(partyID); err != nil {
		return fmt.Errorf("failed to remove trusted party from ledger: %v", err)
	}
	return nil
}

// LogAuthorizationEvent records an event related to authorization actions in the ledger.
func LogAuthorizationEvent(ledgerInstance *ledger.Ledger, event AuthorizationEvent) error {
	// Ensure the event uses ledger.AuthorizationEvent type
	ledgerEvent := ledger.AuthorizationEvent{
		EventID:   event.EventID,
		UserID:    event.UserID,
		Timestamp: time.Now(),
		Details:   event.Details,
	}

	if err := ledgerInstance.AuthorizationLedger.RecordAuthorizationEvent(ledgerEvent); err != nil {
		return fmt.Errorf("failed to log authorization event in ledger: %v", err)
	}

	return nil
}

// SetTrustedPartyFlag marks a trusted party as flagged (e.g., due to suspicious activity), updating the ledger.
func SetTrustedPartyFlag(ledgerInstance *ledger.Ledger, partyID string, flag bool) error {

	err := ledgerInstance.AuthorizationLedger.UpdateTrustedPartyFlag(partyID, flag)
	if err != nil {
		return fmt.Errorf("failed to update trusted party flag in ledger: %v", err)
	}

	return nil
}
