package syn3000

import (
    "errors"
    "time"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/hex"
    "io"
)

// SYN3000Transaction handles transaction operations for SYN3000 tokens.
type SYN3000Transaction struct {
    security *SYN3000Security // Security for encryption and decryption
    storage  *SYN3000Storage  // Storage layer for handling rental properties and payments
}

// NewSYN3000Transaction initializes the transaction management for SYN3000 tokens.
func NewSYN3000Transaction(security *SYN3000Security, storage *SYN3000Storage) *SYN3000Transaction {
    return &SYN3000Transaction{
        security: security,
        storage:  storage,
    }
}

// ProcessRentPayment processes a rent payment for a given rental property.
func (t *SYN3000Transaction) ProcessRentPayment(propertyID, tenantID string, amount float64) (*Payment, error) {
    // Fetch the rental property
    property, err := t.storage.GetRentalProperty(propertyID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve rental property: %v", err)
    }

    // Validate that the tenant is authorized to make the payment
    if property.TenantID != tenantID {
        return nil, errors.New("unauthorized: tenant ID does not match the rental property")
    }

    // Ensure the amount is correct (matches the monthly rent)
    if amount != property.MonthlyRent {
        return nil, errors.New("invalid payment amount: does not match the monthly rent")
    }

    // Record payment details
    payment := &Payment{
        PropertyID:      propertyID,
        TenantID:        tenantID,
        Amount:          amount,
        TransactionDate: time.Now(),
        Status:          "Completed",
    }

    // Encrypt and store the payment history
    err = t.storage.TrackPaymentHistory(propertyID, payment)
    if err != nil {
        return nil, fmt.Errorf("failed to store payment history: %v", err)
    }

    // Return the processed payment details
    return payment, nil
}

// ProcessSecurityDeposit processes a security deposit transaction for a rental property.
func (t *SYN3000Transaction) ProcessSecurityDeposit(propertyID, tenantID string, amount float64) (*Payment, error) {
    // Fetch the rental property
    property, err := t.storage.GetRentalProperty(propertyID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve rental property: %v", err)
    }

    // Validate that the tenant is authorized to make the deposit
    if property.TenantID != tenantID {
        return nil, errors.New("unauthorized: tenant ID does not match the rental property")
    }

    // Ensure the amount matches the required deposit
    if amount != property.Deposit {
        return nil, errors.New("invalid deposit amount: does not match the required deposit")
    }

    // Record deposit details
    deposit := &Payment{
        PropertyID:      propertyID,
        TenantID:        tenantID,
        Amount:          amount,
        TransactionDate: time.Now(),
        Status:          "Completed",
    }

    // Encrypt and store the deposit in the ledger
    err = t.storage.TrackPaymentHistory(propertyID, deposit)
    if err != nil {
        return nil, fmt.Errorf("failed to store security deposit: %v", err)
    }

    return deposit, nil
}

// TransferOwnership transfers ownership of a rental property to a new owner.
func (t *SYN3000Transaction) TransferOwnership(propertyID, newOwnerID string) error {
    // Fetch the rental property
    property, err := t.storage.GetRentalProperty(propertyID)
    if err != nil {
        return fmt.Errorf("failed to retrieve rental property: %v", err)
    }

    // Verify that the new owner is valid (this can be extended with further validation rules)
    if newOwnerID == "" {
        return errors.New("invalid new owner ID")
    }

    // Update property ownership
    property.OwnerID = newOwnerID
    property.LastUpdated = time.Now()

    // Encrypt and store the updated property data
    err = t.storage.StoreRentalProperty(property)
    if err != nil {
        return fmt.Errorf("failed to update rental property ownership: %v", err)
    }

    return nil
}

// ValidateTransaction validates a transaction for integrity, ownership, and authorization.
func (t *SYN3000Transaction) ValidateTransaction(transactionID, userID string) error {
    // Use the security layer to ensure the transaction is secure and valid
    err := t.security.ValidateTransactionSecurity(transactionID, userID)
    if err != nil {
        return fmt.Errorf("transaction validation failed: %v", err)
    }
    return nil
}

// RefundDeposit processes a security deposit refund for a rental property.
func (t *SYN3000Transaction) RefundDeposit(propertyID, tenantID string) (*Payment, error) {
    // Fetch the rental property
    property, err := t.storage.GetRentalProperty(propertyID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve rental property: %v", err)
    }

    // Validate that the tenant is authorized to receive the deposit
    if property.TenantID != tenantID {
        return nil, errors.New("unauthorized: tenant ID does not match the rental property")
    }

    // Ensure the property lease has ended before refunding the deposit
    if time.Now().Before(property.LeaseEndDate) {
        return nil, errors.New("lease has not ended: security deposit cannot be refunded")
    }

    // Record the refund details
    refund := &Payment{
        PropertyID:      propertyID,
        TenantID:        tenantID,
        Amount:          property.Deposit,
        TransactionDate: time.Now(),
        Status:          "Refunded",
    }

    // Encrypt and store the refund in the ledger
    err = t.storage.TrackPaymentHistory(propertyID, refund)
    if err != nil {
        return nil, fmt.Errorf("failed to store security deposit refund: %v", err)
    }

    return refund, nil
}

// RecordLeaseAgreement records the details of a rental lease agreement securely in the ledger.
func (t *SYN3000Transaction) RecordLeaseAgreement(propertyID string, lease *RentalProperty) error {
    if lease == nil {
        return errors.New("lease data is empty")
    }

    // Encrypt and store the lease agreement in the ledger
    err := t.storage.StoreRentalProperty(lease)
    if err != nil {
        return fmt.Errorf("failed to record lease agreement: %v", err)
    }

    return nil
}

// CancelLeaseAgreement cancels a rental lease agreement and removes it from the ledger.
func (t *SYN3000Transaction) CancelLeaseAgreement(propertyID, tenantID string) error {
    // Fetch the rental property
    property, err := t.storage.GetRentalProperty(propertyID)
    if err != nil {
        return fmt.Errorf("failed to retrieve rental property: %v", err)
    }

    // Validate that the tenant is authorized to cancel the lease
    if property.TenantID != tenantID {
        return errors.New("unauthorized: tenant ID does not match the rental property")
    }

    // Remove the lease agreement from the ledger
    err = t.storage.RemoveRentalProperty(propertyID)
    if err != nil {
        return fmt.Errorf("failed to cancel lease agreement: %v", err)
    }

    return nil
}
