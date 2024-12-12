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

// BiometricRegistration represents a biometric registration entry in the ledger.
type BiometricRegistration struct {
	UserID         string    // ID of the user
	EncryptedData  []byte    // Encrypted biometric data
	RegistrationAt time.Time // Timestamp of the registration
}


// BiometricAccessLog records the outcome of a biometric authentication attempt.
type BiometricAccessLog struct {
	UserID    string
	Timestamp time.Time
	Success   bool
}

// BiometricData represents the biometric information of a user.
type BiometricData struct {
    FingerprintHash []byte // Stores a hash of the user's fingerprint data
    FaceIDHash      []byte // Stores a hash of the user's facial recognition data
    IrisScanHash    []byte // Stores a hash of the user's iris scan data
    VoicePrintHash  []byte // Stores a hash of the user's voice print data
	UserID        string // ID of the user associated with the biometric data
    BiometricInfo string // The actual biometric information to be encrypted
}

// BiometricUpdate represents an update to a user's biometric data.
type BiometricUpdate struct {
    UserID         string    // ID of the user
    EncryptedData  []byte    // Encrypted biometric data
    UpdatedAt      time.Time // Timestamp of the update
}


// RegisterBiometricData registers a user's biometric data in the ledger after encrypting it.
func RegisterBiometricData(ledgerInstance *ledger.Ledger, userID string, biometricData BiometricData) error {
	// Encrypt the biometric data
	encryptedData, err := encryptBiometricData(biometricData)
	if err != nil {
		return fmt.Errorf("encryption of biometric data failed: %v", err)
	}

	// Convert encryptedData to []byte if necessary (for compatibility)
	encryptedDataBytes := []byte(encryptedData) // Ensure []byte type for EncryptedData field

	// Create a registration record
	registration := ledger.BiometricRegistration{
		UserID:         userID,
		EncryptedData:  encryptedDataBytes, // Use the []byte variable here
		RegistrationAt: time.Now(),
	}

	// Record the biometric registration in the ledger
	if err := ledgerInstance.AuthorizationLedger.RecordBiometricRegistration(registration); err != nil {
		return fmt.Errorf("recording biometric registration in ledger failed: %v", err)
	}

	return nil
}



// AuthenticateBiometricData authenticates a user's biometric data against stored data in the ledger.
func AuthenticateBiometricData(ledgerInstance *ledger.Ledger, userID string, providedData BiometricData) (bool, error) {
	// Retrieve the stored encrypted biometric data from the ledger
	storedData, err := ledgerInstance.AuthorizationLedger.FetchBiometricData(userID)
	if err != nil {
		return false, fmt.Errorf("retrieving biometric data from ledger failed: %v", err)
	}

	// Decrypt the stored biometric data
	decryptedData, err := decryptBiometricData(string(storedData.EncryptedData)) // Convert []byte to string
	if err != nil {
		return false, fmt.Errorf("decryption of biometric data failed: %v", err)
	}

	// Compare the provided biometric data with the stored data
	if matchBiometricData(providedData, decryptedData) {
		if err := logBiometricAccess(ledgerInstance, userID, true); err != nil {
			return false, fmt.Errorf("failed to log biometric access: %v", err)
		}
		return true, nil
	}

	// Log a failed access attempt
	if err := logBiometricAccess(ledgerInstance, userID, false); err != nil {
		return false, fmt.Errorf("failed to log biometric access: %v", err)
	}
	return false, nil
}





// RevokeBiometricAccess revokes biometric access for a user by removing their biometric data from the ledger.
func RevokeBiometricAccess(ledgerInstance *ledger.Ledger, userID string) error {

    if err := ledgerInstance.AuthorizationLedger.DeleteBiometricData(userID); err != nil {
        return fmt.Errorf("failed to remove biometric data from ledger: %v", err)
    }

    return nil
}

// UpdateBiometricData updates a user's biometric data in the ledger, re-encrypting it for security.
func UpdateBiometricData(ledgerInstance *ledger.Ledger, userID string, newBiometricData BiometricData) error {
    encryptedData, err := encryptBiometricData(newBiometricData)
    if err != nil {
        return fmt.Errorf("encryption of updated biometric data failed: %v", err)
    }

    update := ledger.BiometricUpdate{
        UserID:        userID,
        EncryptedData: encryptedData,
        UpdatedAt:     time.Now(),
    }

    if err := ledgerInstance.AuthorizationLedger.RecordBiometricUpdate(update); err != nil {
        return fmt.Errorf("failed to record updated biometric data in ledger: %v", err)
    }

    return nil
}

// encryptBiometricData encrypts the user's biometric data using AES-GCM.
func encryptBiometricData(data BiometricData) ([]byte, error) {
	key := sha256.Sum256([]byte(data.UserID)) // Generate key based on user ID for consistency
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, errors.New("failed to create AES cipher block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create GCM block")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("failed to generate nonce")
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(data.BiometricInfo), nil)
	return encrypted, nil
}

// decryptBiometricData decrypts the user's biometric data using AES-GCM.
func decryptBiometricData(encryptedData string) (BiometricData, error) {
	decoded, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return BiometricData{}, errors.New("failed to decode base64 encrypted data")
	}

	key := sha256.Sum256([]byte(decoded)) // Using a predefined key based on userID or other unique identifier
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return BiometricData{}, errors.New("failed to create AES cipher block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return BiometricData{}, errors.New("failed to create GCM block")
	}

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return BiometricData{}, errors.New("ciphertext too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return BiometricData{}, errors.New("decryption failed")
	}

	return BiometricData{BiometricInfo: string(decrypted)}, nil
}

// matchBiometricData checks if the provided and stored biometric data match.
func matchBiometricData(provided, stored BiometricData) bool {
	return provided.BiometricInfo == stored.BiometricInfo
}

// logBiometricAccess logs biometric access attempts and results in the ledger.
func logBiometricAccess(ledgerInstance *ledger.Ledger, userID string, success bool) error {
	// Define the access result as a string based on the boolean success
	action := "failure"
	if success {
		action = "success"
	}

	// Create the access log in the ledger using individual fields
	if err := ledgerInstance.AuthorizationLedger.RecordBiometricAccess(userID, action); err != nil {
		return fmt.Errorf("failed to record biometric access log in ledger: %v", err)
	}

	return nil
}

