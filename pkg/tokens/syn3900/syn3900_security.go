package syn3900

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
	"sync"

)

// SecurityManager manages the security of SYN3900 benefit tokens, including breach monitoring and access revocation.
type SecurityManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// NewSecurityManager creates a new SecurityManager.
func NewSecurityManager(ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *SecurityManager {
	return &SecurityManager{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// MonitorSecurityBreach listens for potential security breach events, logs them, and takes necessary actions.
func (sm *SecurityManager) MonitorSecurityBreach(tokenID string, eventType string, details string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Log the security breach event in the ledger
	if err := sm.ledgerService.LogEvent(eventType, time.Now(), tokenID); err != nil {
		return err
	}

	// Trigger an alert or security action
	err := sm.triggerSecurityAlert(tokenID, eventType, details)
	if err != nil {
		return err
	}

	// Depending on severity, freeze or flag the token for manual review
	if severity := sm.assessBreachSeverity(eventType, details); severity == "high" {
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
	// This function alerts administrators or other stakeholders via email, SMS, or dashboard notification.
	alertMessage := fmt.Sprintf("Security breach detected: %s for Token ID: %s | Details: %s", eventType, tokenID, details)
	if err := sm.ledgerService.AlertNetwork(alertMessage); err != nil {
		return err
	}

	// Log the alert in the ledger
	if err := sm.ledgerService.LogEvent("SecurityAlertTriggered", time.Now(), tokenID); err != nil {
		return err
	}

	return nil
}

// freezeToken temporarily freezes the token if the breach is deemed severe enough to warrant suspension.
func (sm *SecurityManager) freezeToken(tokenID string, details string) error {
	// Update the token status to "Frozen" in the ledger to prevent further actions.
	freezeDetails := fmt.Sprintf("Token frozen due to security breach | Details: %s", details)
	if err := sm.ledgerService.UpdateTokenStatus(tokenID, "Frozen", freezeDetails); err != nil {
		return err
	}

	// Log the freeze action in the ledger
	if err := sm.ledgerService.LogEvent("TokenFrozen", time.Now(), tokenID); err != nil {
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

	// Additional step: take actual action to revoke access (e.g., update token status)
	if err := sm.revokeToken(tokenID); err != nil {
		return err
	}

	return nil
}

// revokeToken performs the actual revocation of token access, preventing further usage.
func (sm *SecurityManager) revokeToken(tokenID string) error {
	// Example action: set token status to "Revoked" and mark it as unusable.
	revocationDetails := fmt.Sprintf("Token revoked due to security breach or ownership dispute | Token ID: %s", tokenID)
	if err := sm.ledgerService.UpdateTokenStatus(tokenID, "Revoked", revocationDetails); err != nil {
		return err
	}

	// Log the revocation in the ledger
	if err := sm.ledgerService.LogEvent("TokenRevoked", time.Now(), tokenID); err != nil {
		return err
	}

	return nil
}

// assessBreachSeverity determines the severity of a security breach based on the event type and details.
func (sm *SecurityManager) assessBreachSeverity(eventType, details string) string {
	// Example logic to assess severity - could be more complex based on event patterns, prior incidents, etc.
	if eventType == "UnauthorizedAccess" || eventType == "TamperingDetected" {
		return "high"
	}
	return "low"
}

// signTransaction signs a transaction using the sender's private key, ensuring authenticity.
func (sm *SecurityManager) signTransaction(tx *Syn3900Transaction, privateKey *ecdsa.PrivateKey) (string, error) {
	// Hash the transaction data to create a digest
	txData := tx.TransactionID + tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Amount)
	hash := sha256.Sum256([]byte(txData))

	// Use the private key to sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}

	// Concatenate the r and s values as the signature
	signature := append(r.Bytes(), s.Bytes()...)
	return hex.EncodeToString(signature), nil
}

// validateTransactionSignature validates the digital signature of a transaction.
func (sm *SecurityManager) validateTransactionSignature(tx *Syn3900Transaction, publicKey *ecdsa.PublicKey) bool {
	// Recreate the transaction data to get the hash
	txData := tx.TransactionID + tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Amount)
	hash := sha256.Sum256([]byte(txData))

	// Decode the signature from hex
	signatureBytes, err := hex.DecodeString(tx.Signature)
	if err != nil || len(signatureBytes) < 64 {
		return false
	}

	// Split the signature into r and s values
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])

	// Verify the signature using the sender's public key
	return ecdsa.Verify(publicKey, hash[:], r, s)
}

// getPrivateKeyFromVault securely retrieves the private key for a given sender from a vault.
func (sm *SecurityManager) getPrivateKeyFromVault(sender string) (*ecdsa.PrivateKey, error) {
	// Securely retrieve the private key (from a key vault or HSM)
	privateKey, err := securekeyvault.GetPrivateKey("key-id-for-" + sender)
	if err != nil {
		return nil, fmt.Errorf("error retrieving private key for sender %s: %w", sender, err)
	}
	return privateKey, nil
}

// getPublicKeyFromVault securely retrieves the public key for a given sender from a vault.
func (sm *SecurityManager) getPublicKeyFromVault(sender string) (*ecdsa.PublicKey, error) {
	// Securely retrieve the public key (from a key vault or HSM)
	publicKey, err := securekeyvault.GetPublicKey("key-id-for-" + sender)
	if err != nil {
		return nil, fmt.Errorf("error retrieving public key for sender %s: %w", sender, err)
	}
	return publicKey, nil
}
