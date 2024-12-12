package syn1100

import (
	"errors"
	"sync"
	"time"

)

// SYN1100StorageManager manages storage operations for SYN1100 healthcare data tokens
type SYN1100StorageManager struct {
	Ledger            *ledger.Ledger                // Ledger to store token information
	ConsensusEngine   *consensus.SynnergyConsensus  // Synnergy Consensus engine for validation
	EncryptionService *encryption.EncryptionService // Encryption service for secure data storage
	mutex             sync.Mutex                    // Mutex for concurrency control
}

// SYN1100TokenStorage represents the stored data for SYN1100 tokens
type SYN1100TokenStorage struct {
	TokenID       string          `json:"token_id"`
	Owner         string          `json:"owner"`
	EncryptedData string          `json:"encrypted_data"` // Encrypted healthcare data
	AccessControl map[string]string `json:"access_control"` // Access control: userID -> access level
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// StoreHealthcareData securely stores the encrypted healthcare data for a SYN1100 token
func (sm *SYN1100StorageManager) StoreHealthcareData(tokenID, owner, rawData string) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Encrypt the healthcare data
	encryptedData, err := sm.EncryptHealthcareData(tokenID, rawData)
	if err != nil {
		return "", err
	}

	// Create storage entry for the token
	tokenStorage := SYN1100TokenStorage{
		TokenID:       tokenID,
		Owner:         owner,
		EncryptedData: encryptedData,
		AccessControl: map[string]string{
			owner: "full",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store the encrypted data in the ledger
	storageData := common.StructToString(tokenStorage)
	if err := sm.Ledger.StoreTokenData(tokenID, storageData); err != nil {
		return "", errors.New("failed to store token data in the ledger")
	}

	// Validate the stored data with Synnergy Consensus
	if err := sm.ConsensusEngine.ValidateData(tokenID); err != nil {
		return "", errors.New("data validation failed via Synnergy Consensus")
	}

	return encryptedData, nil
}

// RetrieveHealthcareData securely retrieves the decrypted healthcare data for a SYN1100 token
func (sm *SYN1100StorageManager) RetrieveHealthcareData(tokenID, userID string) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve token data from the ledger
	tokenData, err := sm.Ledger.GetTokenData(tokenID)
	if err != nil {
		return "", errors.New("failed to retrieve token data from the ledger")
	}

	var tokenStorage SYN1100TokenStorage
	if err := common.StringToStruct(tokenData, &tokenStorage); err != nil {
		return "", errors.New("failed to unmarshal token data")
	}

	// Check if the user has access to the healthcare data
	accessLevel, hasAccess := tokenStorage.AccessControl[userID]
	if !hasAccess || accessLevel == "" {
		return "", errors.New("user does not have access to the data")
	}

	// Decrypt the healthcare data
	decryptedData, err := sm.DecryptHealthcareData(tokenID, tokenStorage.EncryptedData)
	if err != nil {
		return "", errors.New("failed to decrypt healthcare data")
	}

	return decryptedData, nil
}

// EncryptHealthcareData encrypts healthcare data using the encryption service
func (sm *SYN1100StorageManager) EncryptHealthcareData(tokenID, rawData string) (string, error) {
	encryptionKey := sm.EncryptionService.GenerateKey()
	encryptedData, err := sm.EncryptionService.EncryptData([]byte(rawData), encryptionKey)
	if err != nil {
		return "", errors.New("failed to encrypt healthcare data")
	}

	// Store the encryption key securely in the ledger
	if err := sm.Ledger.StoreEncryptionKey(tokenID, encryptionKey); err != nil {
		return "", errors.New("failed to store encryption key in ledger")
	}

	return string(encryptedData), nil
}

// DecryptHealthcareData decrypts healthcare data using the encryption service
func (sm *SYN1100StorageManager) DecryptHealthcareData(tokenID, encryptedData string) (string, error) {
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

// UpdateHealthcareData securely updates the healthcare data for a SYN1100 token
func (sm *SYN1100StorageManager) UpdateHealthcareData(tokenID, newData, userID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve token data from the ledger
	tokenData, err := sm.Ledger.GetTokenData(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token data from the ledger")
	}

	var tokenStorage SYN1100TokenStorage
	if err := common.StringToStruct(tokenData, &tokenStorage); err != nil {
		return errors.New("failed to unmarshal token data")
	}

	// Check if the user has permission to update the data
	if accessLevel, hasAccess := tokenStorage.AccessControl[userID]; !hasAccess || accessLevel != "full" {
		return errors.New("user does not have permission to update the data")
	}

	// Encrypt the new healthcare data
	encryptedData, err := sm.EncryptHealthcareData(tokenID, newData)
	if err != nil {
		return err
	}

	// Update token storage
	tokenStorage.EncryptedData = encryptedData
	tokenStorage.UpdatedAt = time.Now()

	// Store the updated token data in the ledger
	updatedData := common.StructToString(tokenStorage)
	if err := sm.Ledger.UpdateTokenData(tokenID, updatedData); err != nil {
		return errors.New("failed to update token data in the ledger")
	}

	// Validate the updated data with Synnergy Consensus
	if err := sm.ConsensusEngine.ValidateData(tokenID); err != nil {
		return errors.New("data validation failed via Synnergy Consensus")
	}

	return nil
}

// AddAccessControl adds or updates access control for a healthcare token
func (sm *SYN1100StorageManager) AddAccessControl(tokenID, userID, accessLevel string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve token data from the ledger
	tokenData, err := sm.Ledger.GetTokenData(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token data from the ledger")
	}

	var tokenStorage SYN1100TokenStorage
	if err := common.StringToStruct(tokenData, &tokenStorage); err != nil {
		return errors.New("failed to unmarshal token data")
	}

	// Add or update access control
	tokenStorage.AccessControl[userID] = accessLevel
	tokenStorage.UpdatedAt = time.Now()

	// Store the updated access control data in the ledger
	updatedData := common.StructToString(tokenStorage)
	if err := sm.Ledger.UpdateTokenData(tokenID, updatedData); err != nil {
		return errors.New("failed to update access control data in the ledger")
	}

	// Validate the updated data with Synnergy Consensus
	if err := sm.ConsensusEngine.ValidateData(tokenID); err != nil {
		return errors.New("data validation failed via Synnergy Consensus")
	}

	return nil
}

// RemoveAccessControl removes a user's access to a healthcare token
func (sm *SYN1100StorageManager) RemoveAccessControl(tokenID, userID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve token data from the ledger
	tokenData, err := sm.Ledger.GetTokenData(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token data from the ledger")
	}

	var tokenStorage SYN1100TokenStorage
	if err := common.StringToStruct(tokenData, &tokenStorage); err != nil {
		return errors.New("failed to unmarshal token data")
	}

	// Remove access control for the user
	delete(tokenStorage.AccessControl, userID)
	tokenStorage.UpdatedAt = time.Now()

	// Store the updated token data in the ledger
	updatedData := common.StructToString(tokenStorage)
	if err := sm.Ledger.UpdateTokenData(tokenID, updatedData); err != nil {
		return errors.New("failed to update access control data in the ledger")
	}

	// Validate the updated data with Synnergy Consensus
	if err := sm.ConsensusEngine.ValidateData(tokenID); err != nil {
		return errors.New("data validation failed via Synnergy Consensus")
	}

	return nil
}
