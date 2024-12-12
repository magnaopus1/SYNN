package syn1100

import (
	"errors"
	"sync"
	"time"

)

// SYN1100TransactionManager manages all transaction operations for SYN1100 tokens
type SYN1100TransactionManager struct {
	Ledger            *ledger.Ledger                // Ledger for storing transaction records
	ConsensusEngine   *consensus.SynnergyConsensus  // Synnergy Consensus engine for validation
	EncryptionService *encryption.EncryptionService // Encryption service for secure data handling
	mutex             sync.Mutex                    // Mutex for concurrency control
}

// SYN1100TokenTransaction represents a transaction related to a SYN1100 healthcare token
type SYN1100ransaction struct {
	TransactionID      string               `json:"transaction_id"`
	TokenID            string               `json:"token_id"`
	PatientID          string               `json:"patient_id"`
	DoctorID           string               `json:"doctor_id,omitempty"`
	HealthcareData     string               `json:"encrypted_healthcare_data"`
	AccessPermissions  map[string]string    `json:"access_permissions"`
	Timestamp          time.Time            `json:"timestamp"`
	Status             string               `json:"status"`
}

// NewTransaction creates a new healthcare token transaction, encrypts the data, and stores it in the ledger
func (tm *SYN1100TransactionManager) NewTransaction(tokenID, patientID, doctorID, healthcareData string, accessPermissions map[string]string) (string, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Generate a transaction ID
	transactionID := common.GenerateID()

	// Encrypt healthcare data
	encryptedData, err := tm.EncryptHealthcareData(tokenID, healthcareData)
	if err != nil {
		return "", err
	}

	// Create a new transaction
	transaction := SYN1100TokenTransaction{
		TransactionID:     transactionID,
		TokenID:           tokenID,
		PatientID:         patientID,
		DoctorID:          doctorID,
		HealthcareData:    encryptedData,
		AccessPermissions: accessPermissions,
		Timestamp:         time.Now(),
		Status:            "pending",
	}

	// Serialize the transaction
	serializedTransaction := common.StructToString(transaction)

	// Store transaction in the ledger
	if err := tm.Ledger.StoreTransaction(transactionID, serializedTransaction); err != nil {
		return "", errors.New("failed to store transaction in the ledger")
	}

	// Validate transaction using Synnergy Consensus
	if err := tm.ConsensusEngine.ValidateTransaction(transactionID); err != nil {
		tm.UpdateTransactionStatus(transactionID, "failed")
		return "", errors.New("transaction validation failed via Synnergy Consensus")
	}

	// Update transaction status to "validated"
	if err := tm.UpdateTransactionStatus(transactionID, "validated"); err != nil {
		return "", err
	}

	return transactionID, nil
}

// EncryptHealthcareData encrypts healthcare data before storage in the transaction
func (tm *SYN1100TransactionManager) EncryptHealthcareData(tokenID, healthcareData string) (string, error) {
	encryptionKey := tm.EncryptionService.GenerateKey()
	encryptedData, err := tm.EncryptionService.EncryptData([]byte(healthcareData), encryptionKey)
	if err != nil {
		return "", errors.New("failed to encrypt healthcare data")
	}

	// Store the encryption key securely in the ledger
	if err := tm.Ledger.StoreEncryptionKey(tokenID, encryptionKey); err != nil {
		return "", errors.New("failed to store encryption key in ledger")
	}

	return string(encryptedData), nil
}

// DecryptHealthcareData decrypts healthcare data from a transaction
func (tm *SYN1100TransactionManager) DecryptHealthcareData(tokenID, encryptedData string) (string, error) {
	// Retrieve encryption key from the ledger
	encryptionKey, err := tm.Ledger.GetEncryptionKey(tokenID)
	if err != nil {
		return "", errors.New("failed to retrieve encryption key from ledger")
	}

	// Decrypt the data
	decryptedData, err := tm.EncryptionService.DecryptData([]byte(encryptedData), encryptionKey)
	if err != nil {
		return "", errors.New("failed to decrypt healthcare data")
	}

	return string(decryptedData), nil
}

// GetTransaction retrieves and decrypts a healthcare token transaction
func (tm *SYN1100TransactionManager) GetTransaction(transactionID string, requesterID string) (*SYN1100TokenTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the transaction from the ledger
	transactionData, err := tm.Ledger.GetTransaction(transactionID)
	if err != nil {
		return nil, errors.New("failed to retrieve transaction from ledger")
	}

	var transaction SYN1100TokenTransaction
	if err := common.StringToStruct(transactionData, &transaction); err != nil {
		return nil, errors.New("failed to deserialize transaction data")
	}

	// Verify access permissions
	if accessLevel, ok := transaction.AccessPermissions[requesterID]; !ok || accessLevel == "" {
		return nil, errors.New("access denied for the requested transaction")
	}

	// Decrypt healthcare data
	decryptedData, err := tm.DecryptHealthcareData(transaction.TokenID, transaction.HealthcareData)
	if err != nil {
		return nil, errors.New("failed to decrypt healthcare data")
	}

	// Replace encrypted data with decrypted data
	transaction.HealthcareData = decryptedData

	return &transaction, nil
}

// UpdateTransactionStatus updates the status of a SYN1100 transaction
func (tm *SYN1100TransactionManager) UpdateTransactionStatus(transactionID, newStatus string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the transaction from the ledger
	transactionData, err := tm.Ledger.GetTransaction(transactionID)
	if err != nil {
		return errors.New("failed to retrieve transaction from ledger")
	}

	var transaction SYN1100TokenTransaction
	if err := common.StringToStruct(transactionData, &transaction); err != nil {
		return errors.New("failed to deserialize transaction data")
	}

	// Update the transaction status
	transaction.Status = newStatus
	transaction.Timestamp = time.Now()

	// Serialize and store the updated transaction
	updatedTransactionData := common.StructToString(transaction)
	if err := tm.Ledger.UpdateTransaction(transactionID, updatedTransactionData); err != nil {
		return errors.New("failed to update transaction in the ledger")
	}

	return nil
}

// ValidateAccess checks if a user has the appropriate access level for a transaction
func (tm *SYN1100TransactionManager) ValidateAccess(transactionID, userID string) (bool, error) {
	// Retrieve the transaction from the ledger
	transactionData, err := tm.Ledger.GetTransaction(transactionID)
	if err != nil {
		return false, errors.New("failed to retrieve transaction from ledger")
	}

	var transaction SYN1100TokenTransaction
	if err := common.StringToStruct(transactionData, &transaction); err != nil {
		return false, errors.New("failed to deserialize transaction data")
	}

	// Check access level
	if accessLevel, ok := transaction.AccessPermissions[userID]; ok && accessLevel != "" {
		return true, nil
	}

	return false, errors.New("user does not have the required access level for the transaction")
}

// AddAccessPermission adds or updates a user's access permissions for a healthcare token transaction
func (tm *SYN1100TransactionManager) AddAccessPermission(transactionID, userID, accessLevel string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the transaction from the ledger
	transactionData, err := tm.Ledger.GetTransaction(transactionID)
	if err != nil {
		return errors.New("failed to retrieve transaction from ledger")
	}

	var transaction SYN1100TokenTransaction
	if err := common.StringToStruct(transactionData, &transaction); err != nil {
		return errors.New("failed to deserialize transaction data")
	}

	// Add or update access permissions
	transaction.AccessPermissions[userID] = accessLevel
	transaction.Timestamp = time.Now()

	// Serialize and store the updated transaction
	updatedTransactionData := common.StructToString(transaction)
	if err := tm.Ledger.UpdateTransaction(transactionID, updatedTransactionData); err != nil {
		return errors.New("failed to update transaction in the ledger")
	}

	return nil
}
