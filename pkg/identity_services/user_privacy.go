package identity_services

import (
    "fmt"
    "time"
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/common"
)

// NewPrivacyManager initializes a new PrivacyManager
func NewPrivacyManager(ledgerInstance *ledger.Ledger) *PrivacyManager {
    return &PrivacyManager{
        PrivacyRecords:  make(map[string]*PrivacySettings),
        LedgerInstance:  ledgerInstance,
    }
}

// UpdatePrivacySettings allows a user to update their privacy settings
func (pm *PrivacyManager) UpdatePrivacySettings(userID string, dataEncryption bool, permissionToShare bool) error {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    privacySettings := &PrivacySettings{
        UserID:            userID,
        DataEncryption:    dataEncryption,
        PermissionToShare: permissionToShare,
        LastUpdated:       time.Now(),
    }

    pm.PrivacyRecords[userID] = privacySettings

    // Log the privacy settings update in the ledger
    err := pm.logPrivacySettingsToLedger(privacySettings)
    if err != nil {
        return fmt.Errorf("failed to log privacy settings to ledger: %v", err)
    }

    fmt.Printf("Privacy settings for user %s updated.\n", userID)
    return nil
}

// GetPrivacySettings retrieves the privacy settings for a given user
func (pm *PrivacyManager) GetPrivacySettings(userID string) (*PrivacySettings, error) {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    privacySettings, exists := pm.PrivacyRecords[userID]
    if !exists {
        return nil, errors.New("user privacy settings not found")
    }

    return privacySettings, nil
}

// logPrivacySettingsToLedger logs the privacy settings changes to the ledger for accountability
func (pm *PrivacyManager) logPrivacySettingsToLedger(privacySettings *PrivacySettings) error {
    privacyDetails := fmt.Sprintf("%+v", privacySettings)

    // Encrypt the privacy settings details using the pm.Encryption instance
    encryptedDetails, err := pm.Encryption.EncryptData("AES", []byte(privacyDetails), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt privacy settings details: %v", err)
    }

    // Record the privacy action in the ledger
    pm.LedgerInstance.RecordPrivacyAction(privacySettings.UserID, string(encryptedDetails)) // Convert encrypted details to string

    fmt.Printf("Privacy settings for user %s logged to the ledger.\n", privacySettings.UserID)
    return nil
}



// generatePrivacyHash generates a hash for the privacy settings update
func (pm *PrivacyManager) generatePrivacyHash(userID string, lastUpdated time.Time) string {
    hashInput := fmt.Sprintf("%s%d", userID, lastUpdated.UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// RequestDataDeletion allows a user to request deletion of their data, which is logged in the ledger
func (pm *PrivacyManager) RequestDataDeletion(userID string) error {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    // Check if the user exists
    _, exists := pm.PrivacyRecords[userID]
    if !exists {
        return errors.New("user not found")
    }

    // Log the data deletion request in the ledger
    err := pm.logPrivacyActionToLedger(userID, "Data Deletion Requested")
    if err != nil {
        return fmt.Errorf("failed to log data deletion request to ledger: %v", err)
    }

    fmt.Printf("Data deletion request logged for user %s.\n", userID)
    return nil
}

// logPrivacyActionToLedger logs general privacy actions (like data deletion requests) in the ledger
func (pm *PrivacyManager) logPrivacyActionToLedger(userID string, action string) error {
    // Encrypt the privacy action using the pm.Encryption instance
    encryptedAction, err := pm.Encryption.EncryptData("AES", []byte(action), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt privacy action: %v", err)
    }

    // Record the privacy action in the ledger
    pm.LedgerInstance.RecordPrivacyAction(userID, string(encryptedAction))

    fmt.Printf("Privacy action '%s' for user %s logged to the ledger.\n", action, userID)
    return nil
}



// RevokeDataSharing revokes a user's consent to share data, updating their privacy settings and logging it
func (pm *PrivacyManager) RevokeDataSharing(userID string) error {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    privacySettings, exists := pm.PrivacyRecords[userID]
    if !exists {
        return errors.New("user not found")
    }

    // Update the privacy settings to revoke permission to share
    privacySettings.PermissionToShare = false
    privacySettings.LastUpdated = time.Now()

    // Log the update to the ledger
    err := pm.logPrivacySettingsToLedger(privacySettings)
    if err != nil {
        return fmt.Errorf("failed to log privacy revocation to ledger: %v", err)
    }

    fmt.Printf("Data sharing revoked for user %s.\n", userID)
    return nil
}
