package syn4700

import (
	"errors"
	"sync"
	"time"

)

// SecurityManager handles all security-related operations for Syn4700 tokens.
type SecurityManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// NewSecurityManager creates a new instance of SecurityManager.
func NewSecurityManager(ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *SecurityManager {
	return &SecurityManager{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// EncryptTokenData encrypts the token data for secure storage and transmission.
func (sm *SecurityManager) EncryptTokenData(token *Syn4700Token) ([]byte, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Encrypt the token data
	encryptedData, err := sm.encryptionService.EncryptData(token)
	if err != nil {
		return nil, err
	}

	// Log encryption action in the ledger
	if err := sm.ledgerService.LogEvent("TokenEncrypted", time.Now(), token.TokenID); err != nil {
		return nil, err
	}

	// Validate using consensus
	if err := sm.consensusService.ValidateSubBlock(token.TokenID); err != nil {
		return nil, err
	}

	return encryptedData, nil
}

// DecryptTokenData decrypts the token data for authorized access.
func (sm *SecurityManager) DecryptTokenData(encryptedData []byte) (*Syn4700Token, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Decrypt the token data
	decryptedData, err := sm.encryptionService.DecryptData(encryptedData)
	if err != nil {
		return nil, err
	}

	token := decryptedData.(*Syn4700Token)

	// Log decryption action in the ledger
	if err := sm.ledgerService.LogEvent("TokenDecrypted", time.Now(), token.TokenID); err != nil {
		return nil, err
	}

	// Validate using consensus
	if err := sm.consensusService.ValidateSubBlock(token.TokenID); err != nil {
		return nil, err
	}

	return token, nil
}

// VerifyOwnership ensures the rightful ownership of a Syn4700 token using signature validation.
func (sm *SecurityManager) VerifyOwnership(token *Syn4700Token, signature string, partyID string) (bool, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Verify the signature of the party
	expectedSignature, exists := token.Metadata.Signatures[partyID]
	if !exists || expectedSignature != signature {
		return false, errors.New("invalid ownership: signature does not match")
	}

	// Log the ownership verification event
	if err := sm.ledgerService.LogEvent("OwnershipVerified", time.Now(), token.TokenID); err != nil {
		return false, err
	}

	// Validate the verification with Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(token.TokenID); err != nil {
		return false, err
	}

	return true, nil
}
// MonitorSecurityBreach listens for potential security breach events, logs them, and takes necessary actions.
func (sm *SecurityManager) MonitorSecurityBreach(tokenID string, eventType string, details string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Log the security breach event in the ledger
	if err := sm.ledgerService.LogEvent(eventType, time.Now(), tokenID); err != nil {
		return err
	}

	// Trigger an alert or take appropriate security action
	if err := sm.triggerSecurityAlert(tokenID, eventType, details); err != nil {
		return err
	}

	// Assess the severity of the breach and take further action
	severity := sm.assessBreachSeverity(eventType, details)
	if severity == "high" {
		if err := sm.freezeToken(tokenID, details); err != nil {
			return err
		}
	}

	// Validate the breach event using Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// triggerSecurityAlert triggers appropriate security actions based on the breach event.
func (sm *SecurityManager) triggerSecurityAlert(tokenID, eventType, details string) error {
	alertMessage := "Security breach detected: " + eventType + " for Token ID: " + tokenID + " | Details: " + details
	
	// Send network-wide alert to notify stakeholders
	if err := sm.ledgerService.AlertNetwork(alertMessage); err != nil {
		return err
	}

	// Log the alert action in the ledger
	if err := sm.ledgerService.LogEvent("SecurityAlertTriggered", time.Now(), tokenID); err != nil {
		return err
	}

	// Notify system administrators via internal alerts
	if err := sm.notifyAdministrators(tokenID, eventType, details); err != nil {
		return err
	}

	return nil
}

// assessBreachSeverity determines the severity of a security breach based on the event type and details.
// The severity helps determine the necessary actions (e.g., freeze, revoke).
func (sm *SecurityManager) assessBreachSeverity(eventType, details string) string {
	
	case "UnauthorizedAccess":
		return "high"
	case "TamperingDetected":
		return "high"
	case "SuspiciousActivity":
		return "medium"
	case "MultipleFailedAttempts":
		return "medium"
	default:
		return "low"
	}
}

func (sm *SecurityManager) notifyAdministrators(tokenID, eventType, details string) error {
	// Retrieve the list of administrators from the ledger or system configuration
	admins, err := sm.ledgerService.GetAdministrators()
	if err != nil {
		return err
	}

	// Construct the alert message
	notification := "Admin alert: Security breach detected for Token ID: " + tokenID + " | Event: " + eventType + " | Details: " + details

	// Notify each administrator via their preferred communication method (email, SMS, etc.)
	for _, admin := range admins {
		if err := sm.sendNotification(admin, notification); err != nil {
			return err
		}
	}

	// Log the notification to the administrators in the ledger
	if err := sm.ledgerService.LogEvent("AdminAlertSent", time.Now(), tokenID); err != nil {
		return err
	}

	return nil
}

// freezeToken freezes the token temporarily if the breach is deemed severe enough to warrant suspension.
func (sm *SecurityManager) freezeToken(tokenID string, details string) error {
	// Freeze the token, preventing any transactions or access
	freezeDetails := "Token frozen due to security breach | Details: " + details
	if err := sm.ledgerService.UpdateTokenStatus(tokenID, "Frozen", freezeDetails); err != nil {
		return err
	}

	// Log the freeze action in the ledger
	if err := sm.ledgerService.LogEvent("TokenFrozen", time.Now(), tokenID); err != nil {
		return err
	}

	// Notify stakeholders about the freeze action
	if err := sm.notifyStakeholders(tokenID, "TokenFrozen", freezeDetails); err != nil {
		return err
	}

	return nil
}

func (sm *SecurityManager) notifyStakeholders(tokenID, status, details string) error {
	// Gather a list of stakeholders for the token from the ledger
	stakeholders, err := sm.ledgerService.GetTokenStakeholders(tokenID)
	if err != nil {
		return err
	}

	// Construct the notification message
	notification := "Stakeholder update: Token ID " + tokenID + " has been " + status + ". Details: " + details

	// Notify each stakeholder via their registered communication methods (email, SMS, app notifications, etc.)
	for _, stakeholder := range stakeholders {
		if err := sm.sendNotification(stakeholder, notification); err != nil {
			return err
		}
	}

	// Log the notification event in the ledger
	if err := sm.ledgerService.LogEvent("StakeholderNotificationSent", time.Now(), tokenID); err != nil {
		return err
	}

	return nil
}

// RevokeTokenAccess revokes access to a token if a security breach is detected or ownership is disputed.
func (sm *SecurityManager) RevokeTokenAccess(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Log the revocation event in the ledger
	if err := sm.ledgerService.LogEvent("TokenAccessRevoked", time.Now(), tokenID); err != nil {
		return err
	}

	// Invalidate the token within the consensus system
	if err := sm.consensusService.InvalidateSubBlock(tokenID); err != nil {
		return err
	}

	// Perform actual revocation of the token and notify stakeholders
	if err := sm.revokeToken(tokenID); err != nil {
		return err
	}

	return nil
}

// revokeToken performs the actual revocation of token access, preventing further usage.
func (sm *SecurityManager) revokeToken(tokenID string) error {
	// Mark token as revoked and log the event
	revocationDetails := "Token revoked due to security breach or ownership dispute"
	if err := sm.ledgerService.UpdateTokenStatus(tokenID, "Revoked", revocationDetails); err != nil {
		return err
	}

	// Log the revocation in the ledger
	if err := sm.ledgerService.LogEvent("TokenRevoked", time.Now(), tokenID); err != nil {
		return err
	}

	// Invalidate token within the consensus system, ensuring it cannot be used or recovered
	if err := sm.consensusService.InvalidateSubBlock(tokenID); err != nil {
		return err
	}

	// Notify stakeholders about the revocation
	if err := sm.notifyStakeholders(tokenID, "Revoked", revocationDetails); err != nil {
		return err
	}

	return nil
}
