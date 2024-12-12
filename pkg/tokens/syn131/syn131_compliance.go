package syn131

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)



// NewComplianceManager creates a new instance of ComplianceManager.
func NewComplianceManager(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus) *ComplianceManager {
	return &ComplianceManager{
		complianceRecords: make(map[string]*ComplianceRecord),
		ledger:            ledger,
		consensusEngine:   consensusEngine,
	}
}

// AddComplianceRecord adds a new compliance record for a given token.
func (cm *ComplianceManager) AddComplianceRecord(tokenID string, status string, complianceDetails map[string]string) (*ComplianceRecord, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Encrypt compliance details for secure storage
	encryptedData, encryptionKey, err := encryption.EncryptData(complianceDetails)
	if err != nil {
		return nil, err
	}

	record := &ComplianceRecord{
		TokenID:         tokenID,
		Status:          status,
		ComplianceDate:  time.Now(),
		ComplianceDetails: complianceDetails,
		EncryptedData:   encryptedData,
		EncryptionKey:   encryptionKey,
	}

	// Store the compliance record
	cm.complianceRecords[tokenID] = record

	// Validate the compliance record via Synnergy Consensus
	err = cm.consensusEngine.ValidateComplianceRecord(record)
	if err != nil {
		return nil, err
	}

	// Store the record in the ledger for persistence
	err = cm.ledger.RecordCompliance(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// GetComplianceRecord retrieves the compliance record for a given token.
func (cm *ComplianceManager) GetComplianceRecord(tokenID string) (*ComplianceRecord, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	record, exists := cm.complianceRecords[tokenID]
	if !exists {
		return nil, errors.New("compliance record not found")
	}

	// Decrypt compliance details before returning
	decryptedDetails, err := encryption.DecryptData(record.EncryptedData, record.EncryptionKey)
	if err != nil {
		return nil, err
	}

	record.ComplianceDetails = decryptedDetails

	return record, nil
}

// UpdateComplianceStatus updates the compliance status of a token.
func (cm *ComplianceManager) UpdateComplianceStatus(tokenID string, newStatus string, updatedDetails map[string]string) (*ComplianceRecord, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	record, exists := cm.complianceRecords[tokenID]
	if !exists {
		return nil, errors.New("compliance record not found")
	}

	// Encrypt new compliance details
	encryptedData, encryptionKey, err := encryption.EncryptData(updatedDetails)
	if err != nil {
		return nil, err
	}

	// Update record details
	record.Status = newStatus
	record.ComplianceDate = time.Now()
	record.ComplianceDetails = updatedDetails
	record.EncryptedData = encryptedData
	record.EncryptionKey = encryptionKey

	// Validate updated compliance record via Synnergy Consensus
	err = cm.consensusEngine.ValidateComplianceRecord(record)
	if err != nil {
		return nil, err
	}

	// Update the record in the ledger
	err = cm.ledger.UpdateComplianceRecord(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// RemoveComplianceRecord removes a compliance record from the system.
func (cm *ComplianceManager) RemoveComplianceRecord(tokenID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	_, exists := cm.complianceRecords[tokenID]
	if !exists {
		return errors.New("compliance record not found")
	}

	// Remove from local records
	delete(cm.complianceRecords, tokenID)

	// Remove from ledger
	err := cm.ledger.RemoveComplianceRecord(tokenID)
	if err != nil {
		return err
	}

	return nil
}

// ValidateCompliance performs validation of compliance rules on a token.
func (cm *ComplianceManager) ValidateCompliance(tokenID string, complianceRules map[string]string) (bool, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	record, exists := cm.complianceRecords[tokenID]
	if !exists {
		return false, errors.New("compliance record not found")
	}

	// Decrypt the compliance details
	decryptedDetails, err := encryption.DecryptData(record.EncryptedData, record.EncryptionKey)
	if err != nil {
		return false, err
	}

	// Compare compliance rules with the details
	for key, expectedValue := range complianceRules {
		if decryptedDetails[key] != expectedValue {
			return false, errors.New("compliance rule violation: " + key)
		}
	}

	return true, nil
}

// TriggerComplianceCheck creates a compliance event and validates it.
func (cm *ComplianceManager) TriggerComplianceCheck(tokenID string) (*ComplianceRecord, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	record, exists := cm.complianceRecords[tokenID]
	if !exists {
		return nil, errors.New("compliance record not found")
	}

	// Validate compliance through the Synnergy Consensus
	err := cm.consensusEngine.ValidateComplianceRecord(record)
	if err != nil {
		return nil, err
	}

	// Store validation result in ledger
	err = cm.ledger.RecordCompliance(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// GenerateComplianceReport generates a compliance report for a token.
func (cm *ComplianceManager) GenerateComplianceReport(tokenID string) (map[string]interface{}, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	record, exists := cm.complianceRecords[tokenID]
	if !exists {
		return nil, errors.New("compliance record not found")
	}

	// Decrypt compliance details
	decryptedDetails, err := encryption.DecryptData(record.EncryptedData, record.EncryptionKey)
	if err != nil {
		return nil, err
	}

	// Generate report based on the compliance record
	report := map[string]interface{}{
		"TokenID":          record.TokenID,
		"Status":           record.Status,
		"ComplianceDate":   record.ComplianceDate,
		"ComplianceDetails": decryptedDetails,
	}

	return report, nil
}

// generateUniqueComplianceID generates a unique ID for a compliance record.
func generateUniqueComplianceID() string {
	return "comp_" + time.Now().Format("20060102150405") + "_" + generateRandomString(8)
}

// generateRandomString generates a random string of n length.
func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(bytes)
}
