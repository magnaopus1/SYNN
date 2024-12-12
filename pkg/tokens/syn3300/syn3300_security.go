package syn3300

import (
	"sync"
	"time"

)

// Syn3300SecurityManager handles security-related operations for SYN3300 tokens.
type Syn3300SecurityManager struct {
	ledgerService    *ledger.Ledger
	encryptionService *encryption.Encryptor
	consensusService *consensus.SynnergyConsensus
	mutex            sync.Mutex
	securityLogs     map[string]*SecurityLog
}

// SecurityLog represents a security log for key actions on the SYN3300 tokens.
type SecurityLog struct {
	LogID        string    `json:"log_id"`
	Action       string    `json:"action"`
	TokenID      string    `json:"token_id"`
	Timestamp    time.Time `json:"timestamp"`
	UserID       string    `json:"user_id"`
	SecureHash   string    `json:"secure_hash"`
}

// NewSyn3300SecurityManager creates a new instance of Syn3300SecurityManager.
func NewSyn3300SecurityManager(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *Syn3300SecurityManager {
	return &Syn3300SecurityManager{
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService: consensus,
		securityLogs:     make(map[string]*SecurityLog),
	}
}

// LogSecurityEvent logs a security-related event, encrypts the log, and stores it securely.
func (sm *Syn3300SecurityManager) LogSecurityEvent(action, tokenID, userID string) (*SecurityLog, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Generate a unique LogID.
	logID := sm.generateUniqueLogID()

	// Create the security log.
	log := &SecurityLog{
		LogID:     logID,
		Action:    action,
		TokenID:   tokenID,
		Timestamp: time.Now(),
		UserID:    userID,
	}

	// Generate a secure hash for integrity verification.
	log.SecureHash = sm.generateSecureHash(log)

	// Encrypt the log before storing it.
	encryptedLog, err := sm.encryptionService.EncryptData(log)
	if err != nil {
		return nil, err
	}

	// Log the event in the ledger.
	if err := sm.ledgerService.LogEvent("SecurityLogCreated", time.Now(), logID); err != nil {
		return nil, err
	}

	// Validate the log entry with Synnergy Consensus.
	if err := sm.consensusService.ValidateSubBlock(logID); err != nil {
		return nil, err
	}

	// Store the encrypted log.
	sm.securityLogs[logID] = encryptedLog.(*SecurityLog)

	return log, nil
}

// RetrieveSecurityLog retrieves a security log by its LogID.
func (sm *Syn3300SecurityManager) RetrieveSecurityLog(logID string) (*SecurityLog, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the log.
	log, exists := sm.securityLogs[logID]
	if !exists {
		return nil, errors.New("security log not found")
	}

	// Decrypt the log before returning it.
	decryptedLog, err := sm.encryptionService.DecryptData(log)
	if err != nil {
		return nil, err
	}

	return decryptedLog.(*SecurityLog), nil
}

// ListAllSecurityLogs lists all security-related logs in the system.
func (sm *Syn3300SecurityManager) ListAllSecurityLogs() ([]*SecurityLog, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve all security logs.
	var allLogs []*SecurityLog
	for _, log := range sm.securityLogs {
		// Decrypt each log before adding it to the list.
		decryptedLog, err := sm.encryptionService.DecryptData(log)
		if err != nil {
			return nil, err
		}
		allLogs = append(allLogs, decryptedLog.(*SecurityLog))
	}

	return allLogs, nil
}

// DeleteSecurityLog deletes a security log by its LogID.
func (sm *Syn3300SecurityManager) DeleteSecurityLog(logID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Check if the log exists.
	if _, exists := sm.securityLogs[logID]; !exists {
		return errors.New("security log not found")
	}

	// Remove the log from the securityLogs map.
	delete(sm.securityLogs, logID)

	// Log the deletion event in the ledger.
	if err := sm.ledgerService.LogEvent("SecurityLogDeleted", time.Now(), logID); err != nil {
		return err
	}

	return nil
}

// generateUniqueLogID generates a unique ID for a security log.
func (sm *Syn3300SecurityManager) generateUniqueLogID() string {
	// Generate a unique ID using the current timestamp.
	return "LOG-" + time.Now().Format("20060102150405")
}

// generateSecureHash generates a secure hash for a given security log to ensure data integrity.
func (sm *Syn3300SecurityManager) generateSecureHash(log *SecurityLog) string {
	// Generate a SHA-256 hash for the log.
	hash := sha256.New()
	hash.Write([]byte(log.LogID))
	hash.Write([]byte(log.Action))
	hash.Write([]byte(log.TokenID))
	hash.Write([]byte(log.Timestamp.String()))
	hash.Write([]byte(log.UserID))

	return hex.EncodeToString(hash.Sum(nil))
}

// VerifyLogIntegrity verifies the integrity of a security log by comparing its secure hash.
func (sm *Syn3300SecurityManager) VerifyLogIntegrity(logID string) (bool, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the log.
	log, exists := sm.securityLogs[logID]
	if !exists {
		return false, errors.New("security log not found")
	}

	// Generate the current hash.
	currentHash := sm.generateSecureHash(log)

	// Compare the current hash with the stored secure hash.
	return currentHash == log.SecureHash, nil
}
