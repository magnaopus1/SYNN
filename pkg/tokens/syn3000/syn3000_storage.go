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

// SYN3000Storage handles storage operations for rental property tokens.
type SYN3000Storage struct {
    security *SYN3000Security // Reference to security layer for encryption/decryption
}

// NewSYN3000Storage initializes storage management for SYN3000 tokens.
func NewSYN3000Storage(security *SYN3000Security) *SYN3000Storage {
    return &SYN3000Storage{
        security: security,
    }
}

// RentalProperty represents a property being rented or leased.
type RentalProperty struct {
    PropertyID       string
    OwnerID          string
    TenantID         string
    LeaseStartDate   time.Time
    LeaseEndDate     time.Time
    MonthlyRent      float64
    Deposit          float64
    LastUpdated      time.Time
    Status           string
}

// StoreRentalProperty securely stores rental property data in the ledger.
func (s *SYN3000Storage) StoreRentalProperty(property *RentalProperty) error {
    if property == nil {
        return errors.New("rental property data is empty")
    }

    // Generate a unique token ID
    tokenID := common.GenerateTokenID(property.PropertyID)

    // Encrypt the rental property data
    err := s.security.EncryptAndStore(tokenID, property)
    if err != nil {
        return fmt.Errorf("failed to store rental property: %v", err)
    }

    return nil
}

// GetRentalProperty retrieves and decrypts rental property data from the ledger.
func (s *SYN3000Storage) GetRentalProperty(propertyID string) (*RentalProperty, error) {
    // Generate token ID for fetching the property
    tokenID := common.GenerateTokenID(propertyID)

    // Retrieve encrypted data from the ledger
    encryptedData, err := ledger.GetData(tokenID)
    if err != nil {
        return nil, fmt.Errorf("rental property not found: %v", err)
    }

    // Decrypt the data
    var property RentalProperty
    err = s.security.RetrieveSecureData(encryptedData, &property)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt rental property data: %v", err)
    }

    return &property, nil
}

// UpdateRentalProperty updates rental property details in the ledger.
func (s *SYN3000Storage) UpdateRentalProperty(propertyID string, updatedData *RentalProperty) error {
    if updatedData == nil {
        return errors.New("updated rental property data is empty")
    }

    // Retrieve the existing property
    property, err := s.GetRentalProperty(propertyID)
    if err != nil {
        return fmt.Errorf("failed to retrieve property for update: %v", err)
    }

    // Update fields
    property.LeaseStartDate = updatedData.LeaseStartDate
    property.LeaseEndDate = updatedData.LeaseEndDate
    property.MonthlyRent = updatedData.MonthlyRent
    property.Deposit = updatedData.Deposit
    property.LastUpdated = time.Now()
    property.Status = updatedData.Status

    // Store updated property details securely
    err = s.StoreRentalProperty(property)
    if err != nil {
        return fmt.Errorf("failed to update rental property: %v", err)
    }

    return nil
}

// RemoveRentalProperty deletes a rental property from the ledger.
func (s *SYN3000Storage) RemoveRentalProperty(propertyID string) error {
    // Generate token ID for deletion
    tokenID := common.GenerateTokenID(propertyID)

    // Remove data from ledger
    err := ledger.RemoveData(tokenID)
    if err != nil {
        return fmt.Errorf("failed to remove rental property: %v", err)
    }

    return nil
}

// ListPropertiesByOwner fetches all rental properties owned by a given owner.
func (s *SYN3000Storage) ListPropertiesByOwner(ownerID string) ([]RentalProperty, error) {
    // Fetch all property IDs associated with the owner from the ledger
    propertyIDs, err := ledger.GetDataByOwner(ownerID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve properties for owner: %v", err)
    }

    properties := []RentalProperty{}
    for _, propertyID := range propertyIDs {
        // Retrieve and decrypt each property
        property, err := s.GetRentalProperty(propertyID)
        if err != nil {
            return nil, fmt.Errorf("failed to retrieve property: %v", err)
        }
        properties = append(properties, *property)
    }

    return properties, nil
}

// TrackPaymentHistory stores payment history for a rental property.
func (s *SYN3000Storage) TrackPaymentHistory(propertyID string, payment *Payment) error {
    if payment == nil {
        return errors.New("payment data is empty")
    }

    // Generate a payment token ID
    paymentID := common.GeneratePaymentID(propertyID, payment.TransactionDate)

    // Encrypt and store payment data in the ledger
    err := s.security.EncryptAndStore(paymentID, payment)
    if err != nil {
        return fmt.Errorf("failed to store payment history: %v", err)
    }

    return nil
}

// GetPaymentHistory retrieves payment history for a rental property.
func (s *SYN3000Storage) GetPaymentHistory(propertyID string) ([]Payment, error) {
    // Fetch all payment IDs associated with the property from the ledger
    paymentIDs, err := ledger.GetPaymentIDs(propertyID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve payment history: %v", err)
    }

    payments := []Payment{}
    for _, paymentID := range paymentIDs {
        // Retrieve and decrypt each payment
        var payment Payment
        encryptedData, err := ledger.GetData(paymentID)
        if err != nil {
            return nil, fmt.Errorf("failed to retrieve payment data: %v", err)
        }

        err = s.security.RetrieveSecureData(encryptedData, &payment)
        if err != nil {
            return nil, fmt.Errorf("failed to decrypt payment data: %v", err)
        }
        payments = append(payments, payment)
    }

    return payments, nil
}

// Payment represents a rent payment made by a tenant.
type Payment struct {
    PropertyID      string
    TenantID        string
    Amount          float64
    TransactionDate time.Time
    Status          string
}
