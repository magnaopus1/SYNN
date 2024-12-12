package authorization

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"synnergy_network/pkg/ledger"
	"time"
)

// EntityType defines the type of entity in the system, such as User, Device, or Service.
type EntityType struct {
    TypeID   string // Unique identifier for the entity type
    Name     string // Name of the entity type, e.g., "User", "Device", "Service"
    Category string // Optional category for further classification
}

// IdentityData represents identity-related data for an entity, such as personal details or credentials.
type IdentityData struct {
    EntityID     string    // Unique identifier for the entity
    FullName     string    // Full name of the entity owner (if applicable)
    DateOfBirth  time.Time // Date of birth (optional, if applicable)
    NationalID   string    // National ID or similar identifier (optional)
    Address      string    // Address information (optional)
    ContactInfo  string    // Contact information, such as email or phone
    RegisteredAt time.Time // Date and time of registration
	Metadata    string   
}

// MicrochipAuthorization holds authorization details for a device with a microchip.
type MicrochipAuthorization struct {
    ChipID          string    // Unique identifier for the microchip
    DeviceID        string    // ID of the device associated with the microchip
    AuthorizedUser  string    // ID of the authorized user for the device
    IssuedAt        time.Time // Date and time when the authorization was issued
    ExpiresAt       time.Time // Expiration date and time of the authorization
    EncryptedKey    []byte    // Encrypted authorization key
    AuthorizationStatus string // Status of the authorization, e.g., "Active", "Revoked"
	AuthorizationLevel  string // Added AuthorizationLevel if missing

}


// RegisterMicrochip stores information about a new microchip in the ledger with encryption.
func RegisterMicrochip(ledgerInstance *ledger.Ledger, chipID string, identityData IdentityData) error {
	// Encrypt the identity data for secure storage
	encryptedIdentity, err := encryptIdentityData(identityData)
	if err != nil {
		return fmt.Errorf("encryption of microchip identity data failed: %v", err)
	}

	// Create a record for the microchip authorization
	chipRecord := ledger.MicrochipAuthorization{
		ChipID:             chipID,
		AuthorizedUser:     identityData.EntityID,       // assign the owner of the microchip
		EncryptedKey:       []byte(encryptedIdentity),   // ensure correct data type for EncryptedKey
		IssuedAt:           time.Now(),
		ExpiresAt:          time.Now().AddDate(1, 0, 0), // default expiration of one year
		AuthorizationStatus: "Active",
	}

	// Record the authorization directly in the ledger
	ledgerInstance.AuthorizationLedger.RecordMicrochipAuthorization(chipRecord.ChipID, chipRecord.AuthorizedUser, chipRecord.AuthorizationStatus)

	return nil
}


// AuthorizeMicrochipAccess validates if a microchip has permission to access certain functionalities.
func AuthorizeMicrochipAccess(ledgerInstance *ledger.Ledger, chipID string, requiredLevel string) (bool, error) {
    // Fetch the stored microchip authorization record from the ledger
    chipRecord, err := ledgerInstance.AuthorizationLedger.FetchMicrochipAuthorization(chipID)
    if err != nil {
        return false, fmt.Errorf("error fetching microchip authorization from ledger: %v", err)
    }

    // Verify that AuthorizationStatus and AuthorizationLevel meet requirements
    if chipRecord.AuthorizationStatus != "Active" || chipRecord.AuthorizationLevel != requiredLevel {
        // Log failed access attempt without stopping the function
        _ = logAccessAttempt(ledgerInstance, chipID, false)
        return false, fmt.Errorf("authorization level mismatch or inactive status")
    }

    // Log a successful access attempt
    _ = logAccessAttempt(ledgerInstance, chipID, true)
    return true, nil
}




// UpdateMicrochipAuthorizationLevel changes the authorization level of a microchip.
func UpdateMicrochipAuthorizationLevel(ledgerInstance *ledger.Ledger, chipID string, newLevel string) error {
	// Retrieve the current microchip authorization record to check existence.
	_, err := ledgerInstance.AuthorizationLedger.FetchMicrochipAuthorization(chipID)
	if err != nil {
		return fmt.Errorf("retrieving microchip authorization from ledger failed: %v", err)
	}

	// Update the authorization level in the ledger directly
	ledgerInstance.AuthorizationLedger.UpdateMicrochipAuthorization(chipID, newLevel)

	// Log the role change event
	if err := logRoleChangeEvent(ledgerInstance, chipID, newLevel); err != nil {
		return fmt.Errorf("failed to log role change event: %v", err)
	}

	return nil
}



// RevokeMicrochipAuthorization revokes access for a specific microchip.
func RevokeMicrochipAuthorization(ledgerInstance *ledger.Ledger, chipID string) error {
	// Delete the microchip authorization record from the ledger
	ledgerInstance.AuthorizationLedger.DeleteMicrochipAuthorization(chipID)

	// Log the access revocation attempt as a failed access
	if err := logAccessAttempt(ledgerInstance, chipID, false); err != nil {
		return fmt.Errorf("failed to log revocation access attempt: %v", err)
	}

	return nil
}




// QueryMicrochipStatus retrieves the status of a specific microchip, including authorization level and entity type.
func QueryMicrochipStatus(ledgerInstance *ledger.Ledger, chipID string) (ledger.MicrochipAuthorization, error) {
	chipRecord, err := ledgerInstance.AuthorizationLedger.FetchMicrochipAuthorization(chipID)
	if err != nil {
		return ledger.MicrochipAuthorization{}, fmt.Errorf("failed to retrieve microchip status: %v", err)
	}

	return chipRecord, nil
}


// LogMicrochipAccessAttempt records access attempts made by microchips in the ledger.
func LogMicrochipAccessAttempt(ledgerInstance *ledger.Ledger, chipID string, success bool) error {
	accessLog := ledger.AccessLog{ // Use ledger.AccessLog explicitly
		DeviceID:   chipID,
		Timestamp:  time.Now(),
		Success:    success,
		AccessType: "microchip_access",
		FailureReason: func() string {
			if !success {
				return "authorization level mismatch or inactive status"
			}
			return ""
		}(),
	}

	if err := ledgerInstance.AuthorizationLedger.RecordAccessAttempt(accessLog); err != nil {
		return errors.New("failed to record access attempt in ledger")
	}

	return nil
}






// encryptIdentityData encrypts identity data using AES-GCM for secure storage.
func encryptIdentityData(data IdentityData) (string, error) {
	key := sha256.Sum256([]byte(data.EntityID))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", errors.New("failed to create AES cipher block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("failed to create GCM block")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.New("failed to generate nonce")
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(data.Metadata), nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// decryptIdentityData decrypts encrypted identity data for authorization checks.
func decryptIdentityData(encryptedData string) (IdentityData, error) {
	// Decode the base64-encoded encrypted data
	decoded, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return IdentityData{}, errors.New("failed to decode base64 encrypted data")
	}

	// Generate the decryption key based on the decoded data
	key := sha256.Sum256([]byte("encryption-key")) // Replace "encryption-key" with a secure key source
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return IdentityData{}, errors.New("failed to create AES cipher block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return IdentityData{}, errors.New("failed to create GCM block")
	}

	// Verify and separate nonce and ciphertext
	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return IdentityData{}, errors.New("ciphertext too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return IdentityData{}, errors.New("decryption failed")
	}

	// Convert decrypted data to string and return as part of IdentityData
	return IdentityData{Metadata: string(decrypted)}, nil
}


// logAccessAttempt logs an access attempt for a microchip authorization request.
func logAccessAttempt(ledgerInstance *ledger.Ledger, chipID string, success bool) error {
	// Generate a unique attempt ID
	attemptID := generateUniqueAttemptID()

	// Set the result as a string based on success
	result := "success"
	if !success {
		result = "failure"
	}

	// Record the access attempt in the ledger with compatible parameters
	ledgerInstance.AuthorizationLedger.RecordMicrochipAccessAttempt(attemptID, chipID, result)

	return nil // No error handling is needed as the method does not return a value
}

// logRoleChangeEvent records a change in authorization level or role.
func logRoleChangeEvent(ledgerInstance *ledger.Ledger, chipID, newLevel string) error {
	// Generate a unique role change ID
	roleChangeID := generateUniqueRoleChangeID()

	// Record the role change event in the ledger with compatible parameters
	ledgerInstance.AuthorizationLedger.RecordRoleChange(roleChangeID, chipID, newLevel)

	return nil // No error handling is needed as the method does not return a value
}


func generateUniqueAttemptID() string {
	return fmt.Sprintf("attempt-%d", time.Now().UnixNano())
}

func generateUniqueRoleChangeID() string {
	return fmt.Sprintf("rolechange-%d", time.Now().UnixNano())
}
