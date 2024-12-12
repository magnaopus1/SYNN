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

// SYN3000Management handles management of rental properties and leases.
type SYN3000Management struct{}

// NewSYN3000Management initializes the management system for SYN3000 tokens.
func NewSYN3000Management() *SYN3000Management {
    return &SYN3000Management{}
}

// RentalProperty represents a property being rented.
type RentalProperty struct {
    PropertyID      string
    OwnerID         string
    Address         string
    Description     string
    Bedrooms        int
    Bathrooms       int
    SquareFootage   int
    Availability    bool
}

// LeaseAgreement represents a rental lease.
type LeaseAgreement struct {
    TokenID        string
    PropertyID     string
    TenantID       string
    LeaseStart     time.Time
    LeaseEnd       time.Time
    RentAmount     float64
    DepositAmount  float64
    Status         string // Active, Terminated, Expired
}

// CreateRentalProperty creates a new rental property.
func (m *SYN3000Management) CreateRentalProperty(property RentalProperty) error {
    if property.PropertyID == "" || property.OwnerID == "" || property.Address == "" {
        return errors.New("property details are incomplete")
    }
    
    // Encrypt property data before storing
    encryptedData, err := m.encryptData(property)
    if err != nil {
        return err
    }

    // Store property in ledger
    return ledger.StoreData(property.PropertyID, encryptedData)
}

// CreateLeaseAgreement creates a new lease agreement for a property.
func (m *SYN3000Management) CreateLeaseAgreement(lease LeaseAgreement) error {
    if lease.TokenID == "" || lease.PropertyID == "" || lease.TenantID == "" || lease.LeaseStart.IsZero() || lease.LeaseEnd.IsZero() {
        return errors.New("lease details are incomplete")
    }
    
    // Encrypt lease data before storing
    encryptedData, err := m.encryptData(lease)
    if err != nil {
        return err
    }

    // Store lease in ledger
    return ledger.StoreData(lease.TokenID, encryptedData)
}

// TerminateLease terminates an active lease agreement.
func (m *SYN3000Management) TerminateLease(tokenID string, reason string) error {
    lease, err := m.GetLease(tokenID)
    if err != nil {
        return err
    }

    if lease.Status != "Active" {
        return errors.New("lease is not active")
    }

    // Update lease status
    lease.Status = "Terminated"

    // Encrypt updated lease data
    encryptedData, err := m.encryptData(lease)
    if err != nil {
        return err
    }

    // Store updated lease in ledger
    return ledger.StoreData(tokenID, encryptedData)
}

// TransferPropertyOwnership transfers ownership of a property.
func (m *SYN3000Management) TransferPropertyOwnership(propertyID string, newOwnerID string) error {
    property, err := m.GetProperty(propertyID)
    if err != nil {
        return err
    }

    property.OwnerID = newOwnerID

    // Encrypt updated property data
    encryptedData, err := m.encryptData(property)
    if err != nil {
        return err
    }

    // Store updated property in ledger
    return ledger.StoreData(propertyID, encryptedData)
}

// GetProperty retrieves property details.
func (m *SYN3000Management) GetProperty(propertyID string) (*RentalProperty, error) {
    encryptedData, err := ledger.GetData(propertyID)
    if err != nil {
        return nil, errors.New("property not found in ledger")
    }

    // Decrypt property data
    property := &RentalProperty{}
    err = m.decryptData(encryptedData, property)
    if err != nil {
        return nil, err
    }

    return property, nil
}

// GetLease retrieves lease agreement details.
func (m *SYN3000Management) GetLease(tokenID string) (*LeaseAgreement, error) {
    encryptedData, err := ledger.GetData(tokenID)
    if err != nil {
        return nil, errors.New("lease not found in ledger")
    }

    // Decrypt lease data
    lease := &LeaseAgreement{}
    err = m.decryptData(encryptedData, lease)
    if err != nil {
        return nil, err
    }

    return lease, nil
}

// Encrypt data before storing in the ledger.
func (m *SYN3000Management) encryptData(data interface{}) (string, error) {
    key := []byte(common.GetEncryptionKey()) // Fetch encryption key from common package
    plaintext, err := common.Serialize(data)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return hex.EncodeToString(ciphertext), nil
}

// Decrypt data retrieved from the ledger.
func (m *SYN3000Management) decryptData(encryptedData string, out interface{}) error {
    key := []byte(common.GetEncryptionKey())
    ciphertext, err := hex.DecodeString(encryptedData)
    if err != nil {
        return err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    if len(ciphertext) < aes.BlockSize {
        return errors.New("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return common.Deserialize(ciphertext, out)
}
