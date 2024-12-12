package identity_services

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "errors"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger")


const (
    Syn900Identity IdentityType = "Syn900"
    DecentralizedID IdentityType = "DecentralizedID"
)

// NewIdentityVerificationManager initializes a new IdentityVerificationManager
func NewIdentityVerificationManager(ledgerInstance *ledger.Ledger) *IdentityVerificationManager {
    return &IdentityVerificationManager{
        Identities:     make(map[string]*Identity),
        LedgerInstance: ledgerInstance,
    }
}

// CreateIdentity creates a new identity (Syn900 or Decentralized) for a given owner
// CreateIdentity creates a new identity (Syn900 or Decentralized) for a given owner
func (ivm *IdentityVerificationManager) CreateIdentity(owner, identityData string, identityType IdentityType) (string, error) {
    ivm.mutex.Lock()
    defer ivm.mutex.Unlock()

    identityID := ivm.generateIdentityID(owner, identityData, identityType)

    if _, exists := ivm.Identities[identityID]; exists {
        return "", fmt.Errorf("identity with ID %s already exists", identityID)
    }

    // Encrypt the identity data using the ivm.Encryption instance
    encryptedData, err := ivm.Encryption.EncryptData("AES", []byte(identityData), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt identity data: %v", err)
    }

    // Convert encryptedData from []byte to string
    encryptedDataStr := string(encryptedData)

    // Create and store the new identity
    newIdentity := &Identity{
        IdentityID:    identityID,
        IdentityType:  identityType,
        Owner:         owner,
        CreatedAt:     time.Now(),
        EncryptedData: encryptedDataStr,  // Assign the converted string
        IsVerified:    false,
    }

    ivm.Identities[identityID] = newIdentity

    // Log the creation of the new identity in the ledger
    err = ivm.logIdentityToLedger(newIdentity, "Identity Created")
    if err != nil {
        return "", fmt.Errorf("failed to log identity creation to ledger: %v", err)
    }

    fmt.Printf("New %s identity created for owner: %s with ID: %s\n", identityType, owner, identityID)
    return identityID, nil
}



// VerifyIdentity verifies an identity based on identityID
func (ivm *IdentityVerificationManager) VerifyIdentity(identityID, verifier string) error {
    ivm.mutex.Lock()
    defer ivm.mutex.Unlock()

    identity, exists := ivm.Identities[identityID]
    if !exists {
        return errors.New("identity not found")
    }

    if identity.IsVerified {
        return errors.New("identity is already verified")
    }

    // Mark the identity as verified
    identity.IsVerified = true

    // Log the verification of the identity in the ledger
    err := ivm.logIdentityToLedger(identity, "Identity Verified by "+verifier)
    if err != nil {
        return fmt.Errorf("failed to log identity verification to ledger: %v", err)
    }

    fmt.Printf("Identity %s verified by %s.\n", identityID, verifier)
    return nil
}

// GetIdentity retrieves the details of an identity by its ID
func (ivm *IdentityVerificationManager) GetIdentity(identityID string) (*Identity, error) {
    ivm.mutex.Lock()
    defer ivm.mutex.Unlock()

    identity, exists := ivm.Identities[identityID]
    if !exists {
        return nil, errors.New("identity not found")
    }

    return identity, nil
}

// logIdentityToLedger logs identity-related actions in the ledger
func (ivm *IdentityVerificationManager) logIdentityToLedger(identity *Identity, action string) error {
    identityDetails := fmt.Sprintf("%+v", identity)

    // Encrypt the identity details using the ivm.Encryption instance
    encryptedDetails, err := ivm.Encryption.EncryptData("AES", []byte(identityDetails), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt identity details: %v", err)
    }

    // Record the identity action in the ledger
    ivm.LedgerInstance.RecordIdentityAction(action, string(encryptedDetails))

    fmt.Printf("Identity action logged: %s for IdentityID: %s\n", action, identity.IdentityID)
    return nil
}



// generateIdentityID generates a unique ID for an identity based on the owner, data, and identity type
func (ivm *IdentityVerificationManager) generateIdentityID(owner, identityData string, identityType IdentityType) string {
    hashInput := fmt.Sprintf("%s%s%s%d", owner, identityData, identityType, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// RemoveUnverifiedIdentities removes identities that are not verified and older than a specific duration
func (ivm *IdentityVerificationManager) RemoveUnverifiedIdentities(duration time.Duration) {
    ivm.mutex.Lock()
    defer ivm.mutex.Unlock()

    currentTime := time.Now()
    for id, identity := range ivm.Identities {
        if !identity.IsVerified && currentTime.Sub(identity.CreatedAt) > duration {
            delete(ivm.Identities, id)
            fmt.Printf("Removed unverified identity with ID: %s\n", id)
        }
    }
}
