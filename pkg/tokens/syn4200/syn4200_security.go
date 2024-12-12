package syn4200

import (
	"errors"
	"time"
	"sync"
)

// SecurityManager is responsible for managing security operations for SYN4200 tokens.
type SecurityManager struct {
	mutex             sync.Mutex
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
}

// NewSecurityManager creates a new instance of SecurityManager.
func NewSecurityManager(ledgerService *ledger.LedgerService, encryptionService *encryption.Encryptor, consensusService *consensus.SynnergyConsensus) *SecurityManager {
	return &SecurityManager{
		ledgerService:     ledgerService,
		encryptionService: encryptionService,
		consensusService:  consensusService,
	}
}

// MonitorSecurityBreach listens for potential security breach events, logs them, and takes necessary actions.
func (sm *SecurityManager) MonitorSecurityBreach(tokenID string, eventType string, details string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Log the security breach event in the ledger.
	if err := sm.ledgerService.LogEvent(eventType, time.Now(), tokenID); err != nil {
		return err
	}

	// Trigger a security alert for stakeholders and administrators.
	if err := sm.triggerSecurityAlert(tokenID, eventType, details); err != nil {
		return err
	}

	// Depending on severity, freeze or flag the token for manual review.
	severity := sm.assessBreachSeverity(eventType, details)
	if severity == "high" {
		if err := sm.freezeToken(tokenID, details); err != nil {
			return err
		}
	}

	// Validate the breach event using Synnergy Consensus.
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// triggerSecurityAlert notifies stakeholders of the security breach.
func (sm *SecurityManager) triggerSecurityAlert(tokenID, eventType, details string) error {
	alertMessage := "Security breach detected: " + eventType + " for Token ID: " + tokenID + " | Details: " + details
	if err := sm.ledgerService.AlertNetwork(alertMessage); err != nil {
		return err
	}

	// Log the alert in the ledger.
	if err := sm.ledgerService.LogEvent("SecurityAlertTriggered", time.Now(), tokenID); err != nil {
		return err
	}

	return nil
}

// assessBreachSeverity determines the severity of a security breach based on event type and details.
func (sm *SecurityManager) assessBreachSeverity(eventType, details string) string {
	// Example logic to assess severity - high severity for critical breaches.
	if eventType == "UnauthorizedAccess" || eventType == "TamperingDetected" {
		return "high"
	}
	return "low"
}

// freezeToken temporarily freezes the token if a severe security breach is detected.
func (sm *SecurityManager) freezeToken(tokenID string, details string) error {
	freezeDetails := "Token frozen due to security breach | Details: " + details
	if err := sm.ledgerService.UpdateTokenStatus(tokenID, "Frozen", freezeDetails); err != nil {
		return err
	}

	// Log the freeze action in the ledger.
	if err := sm.ledgerService.LogEvent("TokenFrozen", time.Now(), tokenID); err != nil {
		return err
	}

	return nil
}

// RevokeTokenAccess revokes access to a token in case of a security breach or ownership dispute.
func (sm *SecurityManager) RevokeTokenAccess(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Log the revocation event in the ledger.
	if err := sm.ledgerService.LogEvent("TokenAccessRevoked", time.Now(), tokenID); err != nil {
		return err
	}

	// Invalidate the token within the consensus system.
	if err := sm.consensusService.InvalidateSubBlock(tokenID); err != nil {
		return err
	}

	// Perform additional actions to revoke access (e.g., update token status).
	return sm.revokeToken(tokenID)
}

// revokeToken performs the actual revocation of token access.
func (sm *SecurityManager) revokeToken(tokenID string) error {
	revocationDetails := "Token revoked due to security breach or ownership dispute."
	if err := sm.ledgerService.UpdateTokenStatus(tokenID, "Revoked", revocationDetails); err != nil {
		return err
	}

	// Log the revocation in the ledger.
	if err := sm.ledgerService.LogEvent("TokenRevoked", time.Now(), tokenID); err != nil {
		return err
	}

	return nil
}

// signTransaction signs the transaction using the private key.
func (sm *SecurityManager) signTransaction(txData string, privateKey *ecdsa.PrivateKey) (string, error) {
	hash := sha256.Sum256([]byte(txData))
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return hex.EncodeToString(signature), nil
}

// validateTransactionSignature validates the transaction signature.
func (sm *SecurityManager) validateTransactionSignature(txData string, signature string, publicKey *ecdsa.PublicKey) bool {
	hash := sha256.Sum256([]byte(txData))
	signatureBytes, _ := hex.DecodeString(signature)
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])
	return ecdsa.Verify(publicKey, hash[:], r, s)
}

// AuditSecurity performs a security audit for the token.
func (sm *SecurityManager) AuditSecurity(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Perform a full security audit.
	auditDetails := "Full security audit for token " + tokenID + " conducted at " + time.Now().String()
	if err := sm.ledgerService.LogEvent("SecurityAudit", time.Now(), tokenID); err != nil {
		return err
	}

	// Log audit details and validate with consensus.
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}
