package syn1100

import (
	"errors"
	"sync"
	"time"

)

// SYN1100SecurityManager handles the security aspects of SYN1100 tokens, ensuring data protection and compliance.
type SYN1100SecurityManager struct {
	Ledger            *ledger.Ledger                // Ledger for storing token and security info
	ConsensusEngine   *consensus.SynnergyConsensus  // Synnergy Consensus for validation
	EncryptionService *encryption.EncryptionService // Encryption service for secure token operations
	mutex             sync.Mutex                    // Mutex for concurrency
}

// TokenSecurityLog represents security-related logs for healthcare tokens
type TokenSecurityLog struct {
	LogID        string    `json:"log_id"`
	TokenID      string    `json:"token_id"`
	EntryType    string    `json:"entry_type"` // e.g., "ACCESS", "MODIFICATION", "BREACH_ATTEMPT"
	Details      string    `json:"details"`    // Description of the security event
	Timestamp    time.Time `json:"timestamp"`
	RecordedBy   string    `json:"recorded_by"` // User or system who initiated the event
}

// EncryptTokenData encrypts healthcare data within the SYN1100 token
func (sm *SYN1100SecurityManager) EncryptTokenData(tokenID string, rawData string) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	encryptionKey := sm.EncryptionService.GenerateKey()
	encryptedData, err := sm.EncryptionService.EncryptData([]byte(rawData), encryptionKey)
	if err != nil {
		return "", errors.New("failed to encrypt healthcare data")
	}

	// Store the encryption key securely in the ledger
	if err := sm.Ledger.StoreEncryptionKey(tokenID, encryptionKey); err != nil {
		return "", errors.New("failed to store encryption key in the ledger")
	}

	return string(encryptedData), nil
}

// DecryptTokenData decrypts healthcare data within the SYN1100 token
func (sm *SYN1100SecurityManager) DecryptTokenData(tokenID string, encryptedData string) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve encryption key from the ledger
	encryptionKey, err := sm.Ledger.GetEncryptionKey(tokenID)
	if err != nil {
		return "", errors.New("failed to retrieve encryption key from ledger")
	}

	// Decrypt the data
	rawData, err := sm.EncryptionService.DecryptData([]byte(encryptedData), encryptionKey)
	if err != nil {
		return "", errors.New("failed to decrypt healthcare data")
	}

	return string(rawData), nil
}

// LogSecurityEvent logs a security event related to SYN1100 tokens
func (sm *SYN1100SecurityManager) LogSecurityEvent(tokenID, entryType, details, recordedBy string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Create a security log entry
	logID := common.GenerateUUID()
	logEntry := TokenSecurityLog{
		LogID:      logID,
		TokenID:    tokenID,
		EntryType:  entryType,
		Details:    details,
		Timestamp:  time.Now(),
		RecordedBy: recordedBy,
	}

	// Store the security log in the ledger
	logEntryData := common.StructToString(logEntry)
	if err := sm.Ledger.StoreSecurityLog(logID, logEntryData); err != nil {
		return errors.New("failed to store security log in ledger")
	}

	// Validate the event with Synnergy Consensus to ensure it's a valid log
	if err := sm.ConsensusEngine.ValidateLog(logID); err != nil {
		return errors.New("log validation failed via Synnergy Consensus")
	}

	return nil
}

// ValidateAccessRights ensures that the access rights for the token's healthcare data are properly enforced
func (sm *SYN1100SecurityManager) ValidateAccessRights(tokenID, userID string) (bool, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve access control from the ledger
	accessControl, err := sm.Ledger.GetAccessControl(tokenID)
	if err != nil {
		return false, errors.New("failed to retrieve access control data from ledger")
	}

	// Check if the user has access rights
	accessLevel, ok := accessControl[userID]
	if !ok || accessLevel == "" {
		return false, nil // No access for the user
	}

	// Log the access validation attempt
	if err := sm.LogSecurityEvent(tokenID, "ACCESS", "Access validation successful", userID); err != nil {
		return false, err
	}

	return true, nil
}

// RevokeAccess revokes access rights for a specific user on a token
func (sm *SYN1100SecurityManager) RevokeAccess(tokenID, userID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve access control from the ledger
	accessControl, err := sm.Ledger.GetAccessControl(tokenID)
	if err != nil {
		return errors.New("failed to retrieve access control data from ledger")
	}

	// Revoke the access
	delete(accessControl, userID)

	// Update the ledger with new access control settings
	if err := sm.Ledger.UpdateAccessControl(tokenID, accessControl); err != nil {
		return errors.New("failed to update access control in ledger")
	}

	// Log the access revocation
	if err := sm.LogSecurityEvent(tokenID, "MODIFICATION", "Access rights revoked", "system"); err != nil {
		return err
	}

	return nil
}

// MonitorBreachAttempts monitors and logs potential security breach attempts on tokens
func (sm *SYN1100SecurityManager) MonitorBreachAttempts(tokenID string, suspiciousActivity string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Log the potential breach attempt
	if err := sm.LogSecurityEvent(tokenID, "BREACH_ATTEMPT", suspiciousActivity, "system"); err != nil {
		return err
	}

	// Trigger validation through Synnergy Consensus for an audit
	if err := sm.ConsensusEngine.ValidateLog(tokenID); err != nil {
		return errors.New("failed to validate potential breach via Synnergy Consensus")
	}

	return nil
}

// AuditTokenSecurity performs an audit on the token's security logs to ensure compliance with standards
func (sm *SYN1100SecurityManager) AuditTokenSecurity(tokenID string) ([]TokenSecurityLog, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve all security logs for the token from the ledger
	securityLogs, err := sm.Ledger.GetAllSecurityLogsByTokenID(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve security logs from ledger")
	}

	var logs []TokenSecurityLog
	for _, logData := range securityLogs {
		var logEntry TokenSecurityLog
		if err := common.StringToStruct(logData, &logEntry); err != nil {
			return nil, errors.New("failed to unmarshal security log data")
		}
		logs = append(logs, logEntry)
	}

	// Validate the audit logs with Synnergy Consensus
	if err := sm.ConsensusEngine.ValidateLogs(tokenID); err != nil {
		return nil, errors.New("audit log validation failed via Synnergy Consensus")
	}

	return logs, nil
}
