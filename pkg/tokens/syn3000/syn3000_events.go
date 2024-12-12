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

// SYN3000Events handles events related to SYN3000 tokens.
type SYN3000Events struct{}

// NewSYN3000Events initializes the event system for SYN3000 tokens.
func NewSYN3000Events() *SYN3000Events {
    return &SYN3000Events{}
}

// Event types for SYN3000 tokens.
const (
    EventLeaseCreated      = "LeaseCreated"
    EventRentPayment       = "RentPayment"
    EventLeaseTransferred  = "LeaseTransferred"
    EventLeaseTerminated   = "LeaseTerminated"
    EventOwnershipTransferred = "OwnershipTransferred"
)

// LeaseCreatedEvent triggers when a new lease agreement is created.
func (e *SYN3000Events) LeaseCreated(tokenID string, propertyID string, tenantID string, leaseStart time.Time, leaseEnd time.Time, rent float64) error {
    event := common.Event{
        EventType:   EventLeaseCreated,
        TokenID:     tokenID,
        PropertyID:  propertyID,
        TenantID:    tenantID,
        LeaseStart:  leaseStart,
        LeaseEnd:    leaseEnd,
        Rent:        rent,
        Timestamp:   time.Now(),
    }
    
    // Store the event in the ledger
    encryptedEvent, err := e.encryptEventData(event)
    if err != nil {
        return err
    }
    return ledger.StoreEvent(encryptedEvent)
}

// RentPaymentEvent triggers when rent is paid.
func (e *SYN3000Events) RentPayment(tokenID string, propertyID string, tenantID string, amount float64) error {
    event := common.Event{
        EventType:   EventRentPayment,
        TokenID:     tokenID,
        PropertyID:  propertyID,
        TenantID:    tenantID,
        AmountPaid:  amount,
        Timestamp:   time.Now(),
    }
    
    // Store the event in the ledger
    encryptedEvent, err := e.encryptEventData(event)
    if err != nil {
        return err
    }
    return ledger.StoreEvent(encryptedEvent)
}

// LeaseTransferredEvent triggers when a lease is transferred to a new tenant.
func (e *SYN3000Events) LeaseTransferred(tokenID string, oldTenantID string, newTenantID string) error {
    event := common.Event{
        EventType:   EventLeaseTransferred,
        TokenID:     tokenID,
        OldTenantID: oldTenantID,
        NewTenantID: newTenantID,
        Timestamp:   time.Now(),
    }
    
    // Store the event in the ledger
    encryptedEvent, err := e.encryptEventData(event)
    if err != nil {
        return err
    }
    return ledger.StoreEvent(encryptedEvent)
}

// LeaseTerminatedEvent triggers when a lease is terminated.
func (e *SYN3000Events) LeaseTerminated(tokenID string, tenantID string, reason string) error {
    event := common.Event{
        EventType:   EventLeaseTerminated,
        TokenID:     tokenID,
        TenantID:    tenantID,
        Reason:      reason,
        Timestamp:   time.Now(),
    }
    
    // Store the event in the ledger
    encryptedEvent, err := e.encryptEventData(event)
    if err != nil {
        return err
    }
    return ledger.StoreEvent(encryptedEvent)
}

// OwnershipTransferredEvent triggers when property ownership is transferred.
func (e *SYN3000Events) OwnershipTransferred(tokenID string, oldOwnerID string, newOwnerID string) error {
    event := common.Event{
        EventType:   EventOwnershipTransferred,
        TokenID:     tokenID,
        OldOwnerID:  oldOwnerID,
        NewOwnerID:  newOwnerID,
        Timestamp:   time.Now(),
    }
    
    // Store the event in the ledger
    encryptedEvent, err := e.encryptEventData(event)
    if err != nil {
        return err
    }
    return ledger.StoreEvent(encryptedEvent)
}

// RetrieveEvents retrieves all events related to the SYN3000 tokens from the ledger.
func (e *SYN3000Events) RetrieveEvents(tokenID string) ([]common.Event, error) {
    // Fetch encrypted events from the ledger
    encryptedEvents, err := ledger.GetEvents(tokenID)
    if err != nil {
        return nil, errors.New("unable to retrieve events from ledger")
    }

    // Decrypt the events
    var events []common.Event
    for _, encryptedEvent := range encryptedEvents {
        decryptedEvent, err := e.decryptEventData(encryptedEvent)
        if err != nil {
            return nil, err
        }
        events = append(events, *decryptedEvent)
    }

    return events, nil
}

// Encryption and Decryption of Event Data

// Encrypts event data before storing it in the ledger.
func (e *SYN3000Events) encryptEventData(event common.Event) (string, error) {
    key := []byte(common.GetEncryptionKey()) // Fetch encryption key from common package
    plaintext, err := common.Serialize(event)
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

// Decrypts event data retrieved from the ledger.
func (e *SYN3000Events) decryptEventData(encryptedData string) (*common.Event, error) {
    key := []byte(common.GetEncryptionKey())
    ciphertext, err := hex.DecodeString(encryptedData)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    if len(ciphertext) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    var event common.Event
    if err := common.Deserialize(ciphertext, &event); err != nil {
        return nil, err
    }

    return &event, nil
}
