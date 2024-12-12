package syn300

import (
	"errors"
	"sync"
	"time"
)

// ComplianceStatus defines the possible compliance statuses of an address or entity
type ComplianceStatus string

const (
	Compliant    ComplianceStatus = "Compliant"
	NonCompliant ComplianceStatus = "NonCompliant"
	Pending      ComplianceStatus = "Pending"
)

// ComplianceRecord stores compliance data for addresses or entities interacting with the token
type ComplianceRecord struct {
	Address          string
	Status           ComplianceStatus
	LastChecked      time.Time
	Reason           string
	EncryptedDetails string
}

// syn300Compliance manages the compliance-related functions for the SYN300 token standard
type syn300Compliance struct {
	Ledger           *ledger.Ledger
	ComplianceRecords map[string]ComplianceRecord
	mutex            sync.RWMutex
}

// NewSyn300Compliance creates a new compliance manager for SYN300 tokens
func NewSyn300Compliance(ledger *ledger.Ledger) *syn300Compliance {
	return &syn300Compliance{
		Ledger:            ledger,
		ComplianceRecords: make(map[string]ComplianceRecord),
	}
}

// CheckCompliance verifies the compliance status of an address
func (c *syn300Compliance) CheckCompliance(address string) (ComplianceStatus, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	record, exists := c.ComplianceRecords[address]
	if !exists {
		return Pending, errors.New("compliance record not found")
	}
	return record.Status, nil
}

// UpdateCompliance updates the compliance status of an address
func (c *syn300Compliance) UpdateCompliance(address string, status ComplianceStatus, reason string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	encryptedDetails, err := encryption.Encrypt(reason)
	if err != nil {
		return errors.New("failed to encrypt compliance details")
	}

	c.ComplianceRecords[address] = ComplianceRecord{
		Address:          address,
		Status:           status,
		LastChecked:      time.Now(),
		Reason:           reason,
		EncryptedDetails: encryptedDetails,
	}

	// Store compliance data in the ledger for transparency
	if err := c.Ledger.StoreComplianceRecord(address, status); err != nil {
		return errors.New("failed to store compliance record in the ledger")
	}

	return nil
}

// GetComplianceDetails retrieves compliance details for an address (with decryption)
func (c *syn300Compliance) GetComplianceDetails(address string) (ComplianceRecord, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	record, exists := c.ComplianceRecords[address]
	if !exists {
		return ComplianceRecord{}, errors.New("compliance record not found")
	}

	decryptedDetails, err := encryption.Decrypt(record.EncryptedDetails)
	if err != nil {
		return ComplianceRecord{}, errors.New("failed to decrypt compliance details")
	}

	record.Reason = decryptedDetails
	return record, nil
}

// IsCompliant checks if an address is compliant and returns true if it is
func (c *syn300Compliance) IsCompliant(address string) (bool, error) {
	status, err := c.CheckCompliance(address)
	if err != nil {
		return false, err
	}
	return status == Compliant, nil
}

// VerifyAndValidateCompliance runs through a compliance check during transaction validation
func (c *syn300Compliance) VerifyAndValidateCompliance(txID string, address string) error {
	// Ensure compliance check before transaction validation
	status, err := c.CheckCompliance(address)
	if err != nil {
		return err
	}

	if status != Compliant {
		return errors.New("transaction denied: address is non-compliant")
	}

	// Ensure the transaction gets validated through Synnergy Consensus
	if err := consensus.ValidateTransaction(txID); err != nil {
		return errors.New("transaction validation failed under consensus")
	}

	return nil
}

// StoreComplianceHistory stores the compliance status of addresses in the ledger
func (c *syn300Compliance) StoreComplianceHistory() error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for address, record := range c.ComplianceRecords {
		if err := c.Ledger.StoreComplianceRecord(address, record.Status); err != nil {
			return err
		}
	}
	return nil
}

// RemoveComplianceRecord removes a compliance record from the system
func (c *syn300Compliance) RemoveComplianceRecord(address string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.ComplianceRecords[address]; !exists {
		return errors.New("compliance record not found")
	}

	delete(c.ComplianceRecords, address)

	// Remove from ledger as well
	if err := c.Ledger.RemoveComplianceRecord(address); err != nil {
		return errors.New("failed to remove compliance record from ledger")
	}

	return nil
}

