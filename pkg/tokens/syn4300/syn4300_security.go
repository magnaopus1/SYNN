package syn4300

import (
	"errors"
	"sync"
	"time"
)

// SecurityManager manages the security operations for SYN4300 tokens.
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

// FreezeToken temporarily freezes a SYN4300 token, suspending all transactions.
func (sm *SecurityManager) FreezeToken(tokenID string, reason string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Log the freeze event in the ledger
	if err := sm.ledgerService.LogEvent("TokenFrozen", time.Now(), tokenID, "Token frozen due to: "+reason); err != nil {
		return err
	}

	// Update the token status to "frozen" in the ledger
	if err := sm.ledgerService.UpdateTokenStatus(tokenID, "frozen", reason); err != nil {
		return err
	}

	// Validate the freeze event with Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// UnfreezeToken unfreezes a SYN4300 token, allowing transactions to resume.
func (sm *SecurityManager) UnfreezeToken(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Log the unfreeze event in the ledger
	if err := sm.ledgerService.LogEvent("TokenUnfrozen", time.Now(), tokenID, "Token unfrozen"); err != nil {
		return err
	}

	// Update the token status to "active" in the ledger
	if err := sm.ledgerService.UpdateTokenStatus(tokenID, "active", "Token unfrozen and active"); err != nil {
		return err
	}

	// Validate the unfreeze event with Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// RevokeToken revokes a SYN4300 token permanently, preventing any future transactions.
func (sm *SecurityManager) RevokeToken(tokenID string, reason string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Log the revocation event in the ledger
	if err := sm.ledgerService.LogEvent("TokenRevoked", time.Now(), tokenID, "Token revoked due to: "+reason); err != nil {
		return err
	}

	// Update the token status to "revoked" in the ledger
	if err := sm.ledgerService.UpdateTokenStatus(tokenID, "revoked", reason); err != nil {
		return err
	}

	// Invalidate the token in the consensus
	if err := sm.consensusService.InvalidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// VerifyOwnership verifies the ownership of a SYN4300 token using the ledger.
func (sm *SecurityManager) VerifyOwnership(tokenID, owner string) (bool, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the token from the ledger
	tokenData, err := sm.ledgerService.RetrieveToken(tokenID)
	if err != nil {
		return false, err
	}

	// Decrypt the token data
	decryptedToken, err := sm.encryptionService.DecryptData(tokenData)
	if err != nil {
		return false, err
	}

	token := decryptedToken.(*Syn4300Token)

	// Check if the owner matches
	if token.Metadata.Owner != owner {
		return false, errors.New("ownership verification failed")
	}

	return true, nil
}

// MonitorSecurityBreach listens for potential security breaches and takes appropriate action.
func (sm *SecurityManager) MonitorSecurityBreach(tokenID, eventType, details string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Log the security breach event
	if err := sm.ledgerService.LogEvent(eventType, time.Now(), tokenID, details); err != nil {
		return err
	}

	// Depending on severity, freeze the token
	if err := sm.FreezeToken(tokenID, "Security breach detected: "+details); err != nil {
		return err
	}

	// Validate the breach event with Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// AuditTokenSecurity performs a security audit on a SYN4300 token.
func (sm *SecurityManager) AuditTokenSecurity(tokenID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Log the audit event in the ledger
	if err := sm.ledgerService.LogEvent("SecurityAudit", time.Now(), tokenID, "Token security audit performed"); err != nil {
		return err
	}

	// Validate the audit with Synnergy Consensus
	if err := sm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	return nil
}

// EncryptTokenData encrypts token-related data before storing it in the ledger.
func (sm *SecurityManager) EncryptTokenData(tokenData *Syn4300Token) ([]byte, error) {
	// Use the encryption service to encrypt token data
	return sm.encryptionService.EncryptData(tokenData)
}

// DecryptTokenData decrypts token-related data retrieved from the ledger.
func (sm *SecurityManager) DecryptTokenData(encryptedData []byte) (*Syn4300Token, error) {
	// Use the encryption service to decrypt token data
	decryptedData, err := sm.encryptionService.DecryptData(encryptedData)
	if err != nil {
		return nil, err
	}

	return decryptedData.(*Syn4300Token), nil
}
