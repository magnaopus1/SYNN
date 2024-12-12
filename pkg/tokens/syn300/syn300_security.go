package syn300

import (
	"errors"
	"sync"
	"time"
)

// Syn300SecurityManager manages the security-related aspects of SYN300 tokens, including encryption, verification, and consensus validation.
type Syn300SecurityManager struct {
	mutex   sync.RWMutex
	Ledger  *ledger.Ledger
	Keys    map[string]string // Stores encrypted keys for each token holder
}

// NewSyn300SecurityManager creates a new security manager for SYN300 tokens.
func NewSyn300SecurityManager(ledger *ledger.Ledger) *Syn300SecurityManager {
	return &Syn300SecurityManager{
		Ledger: ledger,
		Keys:   make(map[string]string),
	}
}

// EncryptTokenData encrypts sensitive data of the SYN300 token for storage and transmission.
func (sm *Syn300SecurityManager) EncryptTokenData(data string, key string) (string, error) {
	encryptedData, err := encryption.Encrypt(data, key)
	if err != nil {
		return "", errors.New("failed to encrypt data")
	}
	return encryptedData, nil
}

// DecryptTokenData decrypts encrypted token data for secure retrieval.
func (sm *Syn300SecurityManager) DecryptTokenData(encryptedData string, key string) (string, error) {
	decryptedData, err := encryption.Decrypt(encryptedData, key)
	if err != nil {
		return "", errors.New("failed to decrypt data")
	}
	return decryptedData, nil
}

// StoreEncryptionKey securely stores an encryption key for a specific token holder in the ledger.
func (sm *Syn300SecurityManager) StoreEncryptionKey(holderID string, key string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	encryptedKey, err := encryption.Encrypt(key, holderID)
	if err != nil {
		return errors.New("failed to encrypt the key")
	}

	// Store the encrypted key in the manager's key map
	sm.Keys[holderID] = encryptedKey

	// Optionally, store in the ledger for full auditability
	if err := sm.Ledger.StoreKey(holderID, encryptedKey); err != nil {
		return errors.New("failed to store encryption key in the ledger")
	}

	return nil
}

// RetrieveEncryptionKey retrieves and decrypts the encryption key for a specific token holder.
func (sm *Syn300SecurityManager) RetrieveEncryptionKey(holderID string) (string, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	encryptedKey, exists := sm.Keys[holderID]
	if !exists {
		return "", errors.New("encryption key not found")
	}

	// Decrypt the key using the holder's ID
	decryptedKey, err := encryption.Decrypt(encryptedKey, holderID)
	if err != nil {
		return "", errors.New("failed to decrypt the encryption key")
	}

	return decryptedKey, nil
}

// ValidateTransactionSignature validates the cryptographic signature of a governance token transaction.
func (sm *Syn300SecurityManager) ValidateTransactionSignature(tx GovernanceTransaction) (bool, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Retrieve the token holder's key from the ledger or local store
	key, err := sm.RetrieveEncryptionKey(tx.From)
	if err != nil {
		return false, errors.New("failed to retrieve encryption key for signature validation")
	}

	// Validate the signature using the retrieved key
	isValid, err := encryption.VerifySignature(tx.Signature, key)
	if err != nil {
		return false, errors.New("failed to verify transaction signature")
	}

	return isValid, nil
}

// VerifyConsensus performs a validation using Synnergy Consensus to ensure the transaction is legitimate.
func (sm *Syn300SecurityManager) VerifyConsensus(txID string) (bool, error) {
	// Using Synnergy Consensus to validate the transaction
	isValid, err := consensus.ValidateTransaction(txID, "")
	if err != nil || !isValid {
		return false, errors.New("failed consensus validation for transaction")
	}
	return true, nil
}

// AuditSecurityLogs audits all the security logs to check for any breaches or anomalies.
func (sm *Syn300SecurityManager) AuditSecurityLogs() ([]string, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Fetch security logs from the ledger
	logs, err := sm.Ledger.GetSecurityLogs()
	if err != nil {
		return nil, errors.New("failed to retrieve security logs")
	}

	return logs, nil
}

// MonitorSuspiciousActivity monitors and reports suspicious activities on SYN300 token transactions.
func (sm *Syn300SecurityManager) MonitorSuspiciousActivity(tx GovernanceTransaction) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Implement rules to detect suspicious activity (e.g., large transactions, rapid transfers)
	if tx.Amount > 1000000 { // Example threshold
		// Log suspicious activity
		err := sm.Ledger.LogSuspiciousActivity(tx.ID, tx.From, tx.To, tx.Amount, time.Now())
		if err != nil {
			return errors.New("failed to log suspicious activity")
		}
	}
	return nil
}

// RevokeTokenAccess revokes access to a token in case of security breaches or non-compliance.
func (sm *Syn300SecurityManager) RevokeTokenAccess(holderID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Revoke access by deleting the holder's encryption key from both the local store and ledger
	delete(sm.Keys, holderID)

	if err := sm.Ledger.RemoveKey(holderID); err != nil {
		return errors.New("failed to revoke token access from ledger")
	}

	return nil
}
