package authorization

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"synnergy_network/pkg/ledger"
	"time"
)

// DeviceInfo contains metadata and identifiers for an authorized device.
type DeviceInfo struct {
	DeviceName      string    // Friendly name for the device
	DeviceType      string    // Type of device, e.g., "Mobile", "Laptop", "Tablet"
	OperatingSystem string    // OS of the device, e.g., "iOS", "Android", "Windows"
	OSVersion       string    // Version of the operating system
	SerialNumber    string    // Device serial number for unique identification
	IPAddress       string    // IP address of the device
	MACAddress      string    // MAC address for network identification
	LastAccessedAt  time.Time // Timestamp of the last access by this device
}

// DelegatedAccess represents a record of delegated access to a device.
type DelegatedAccess struct {
	DeviceID    string    // ID of the device being accessed
	DelegatorID string    // ID of the user granting access
	DelegateID  string    // ID of the user receiving delegated access
	GrantedAt   time.Time // Timestamp when access was granted
	ExpiresAt   time.Time // Timestamp when access expires
}

type TemporaryAccess struct {
    DeviceID      string    // Unique identifier for the device granted temporary access
    AuthorizedAt  time.Time // Timestamp when access was authorized
    ExpiresAt     time.Time // Expiration timestamp for temporary access
    AuthorizedBy  string    // ID of the user or system that authorized the access
    AccessLevel   string    // Level of access granted (e.g., "read-only", "full-access")
    Reason        string    // Reason for granting temporary access
}

type AccessLog struct {
    AccessID      string    // Unique identifier for the access attempt log
    UserID        string    // ID of the user attempting access
    DeviceID      string    // ID of the device involved in the access attempt
    Timestamp     time.Time // Timestamp of the access attempt
    Success       bool      // Outcome of the attempt (true for success, false for failure)
    IPAddress     string    // IP address from where the attempt was made
    AccessType    string    // Type of access attempted (e.g., "login", "data retrieval")
    FailureReason string    // Reason for failure if the attempt was unsuccessful
}

type RoleChangeLog struct {
    LogID         string    // Unique identifier for the role change log entry
    UserID        string    // ID of the user whose role is being changed
    ChangedBy     string    // ID of the user or system making the role change
    OldRole       string    // Previous role of the user
    NewRole       string    // New role of the user
    ChangedAt     time.Time // Timestamp of when the role change was made
    Reason        string    // Reason for the role change
    ApprovedBy    string    // ID of the approver if needed for compliance
}


type AuthorizationConstraints struct {
    ConstraintID       string    // Unique identifier for the constraint
    UserID             string    // ID of the user subject to these constraints
    MaxAccessLevel     string    // Maximum level of access the user can have
    AccessTimeLimits   []string  // Specific times or time ranges allowed for access (e.g., "09:00-17:00")
    AccessFrequency    int       // Maximum allowed access attempts per day or hour
    DeviceRestrictions []string  // List of allowed devices for the user
    ExpiryDate         time.Time // Expiration date of the constraints
    CreatedAt          time.Time // Timestamp when the constraints were created
    CreatedBy          string    // ID of the user or system that created the constraints
}

type AuthorizationKeyReset struct {
    ResetID         string    // Unique identifier for the key reset request
    UserID          string    // ID of the user whose key is being reset
    InitiatedBy     string    // ID of the user or system initiating the reset
    ApprovedBy      string    // ID of the user or system approving the reset
    RequestedAt     time.Time // Timestamp when the reset was requested
    ApprovedAt      time.Time // Timestamp when the reset was approved
    Reason          string    // Reason for the reset (e.g., "compromised key", "expiration")
    NewKeyHash      string    // Hash of the new authorization key (if applicable)
    Status          string    // Status of the reset request (e.g., "Pending", "Approved", "Completed")
}



// AddAuthorizedDevice registers a device with authorized access in the ledger, encrypting the device information for secure storage.
func AddAuthorizedDevice(ledgerInstance *ledger.Ledger, deviceID, userID string, deviceInfo DeviceInfo) error {
	// Encrypt the device information before storage
	_, err := encryptDeviceData(deviceInfo)
	if err != nil {
		return fmt.Errorf("encryption of device data failed: %v", err)
	}

	// Record the device authorization in the ledger
	if err := ledgerInstance.AuthorizationLedger.RecordDeviceAuthorization(deviceID, userID, "default"); err != nil {
		return fmt.Errorf("failed to record device authorization in ledger: %v", err)
	}

	return nil
}



// RemoveAuthorizedDevice revokes authorization for a device and removes its record from the ledger.
func RemoveAuthorizedDevice(ledgerInstance *ledger.Ledger, deviceID string) error {
	// Delete the device authorization record from the ledger
	ledgerInstance.AuthorizationLedger.DeleteDeviceAuthorization(deviceID)
	return nil
}



// RemoveDelegatedAccess revokes delegated access for a user on a specific device.
func RemoveDelegatedAccess(ledgerInstance *ledger.Ledger, deviceID, delegateID string) error {
	// Delete the delegated access record from the ledger
	if err := ledgerInstance.AuthorizationLedger.DeleteDelegatedAccess(deviceID, delegateID); err != nil {
		return fmt.Errorf("failed to remove delegated access from ledger: %v", err)
	}

	return nil
}


// ResetAuthorizationKeys resets the authorization keys for a device.
func ResetAuthorizationKeys(ledgerInstance *ledger.Ledger, deviceID string, resetBy string) error {
	// Generate a new authorization key, but don't store or use it if not needed
	_, err := generateAuthorizationKey()
	if err != nil {
		return fmt.Errorf("key generation failed: %v", err)
	}

	// Record the key reset in the ledger with the correct parameters
	ledgerInstance.AuthorizationLedger.RecordKeyReset(deviceID, resetBy)

	return nil
}


// Assuming hashKey is implemented as follows:
func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}


// encryptDeviceData encrypts device data using AES-GCM.
func encryptDeviceData(data DeviceInfo) (string, error) {
	key := sha256.Sum256([]byte(data.SerialNumber)) // Use SerialNumber as the unique key for encryption
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

	// Assume that the data to be encrypted is a JSON encoding of the entire DeviceInfo struct
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", errors.New("failed to marshal device data to JSON")
	}

	encrypted := gcm.Seal(nonce, nonce, dataBytes, nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// decryptDeviceData decrypts encrypted device data.
func decryptDeviceData(encryptedData string) (DeviceInfo, error) {
	decoded, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return DeviceInfo{}, errors.New("failed to decode base64 encrypted data")
	}

	key := sha256.Sum256([]byte("encryption_key")) // Example key
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return DeviceInfo{}, errors.New("failed to create AES cipher block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return DeviceInfo{}, errors.New("failed to create GCM block")
	}

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return DeviceInfo{}, errors.New("ciphertext too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return DeviceInfo{}, errors.New("decryption failed")
	}

	// Assume the decrypted data is JSON encoded, and unmarshal it into DeviceInfo
	var deviceInfo DeviceInfo
	if err := json.Unmarshal(decrypted, &deviceInfo); err != nil {
		return DeviceInfo{}, errors.New("failed to unmarshal decrypted data")
	}

	return deviceInfo, nil
}

// matchDeviceInfo checks if provided device information matches stored data.
func matchDeviceInfo(provided, stored DeviceInfo) bool {
	return provided.DeviceName == stored.DeviceName &&
		provided.DeviceType == stored.DeviceType &&
		provided.OperatingSystem == stored.OperatingSystem &&
		provided.OSVersion == stored.OSVersion &&
		provided.SerialNumber == stored.SerialNumber &&
		provided.IPAddress == stored.IPAddress &&
		provided.MACAddress == stored.MACAddress
}
// generateAuthorizationKey generates a new authorization key for device access.
func generateAuthorizationKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", errors.New("failed to generate new authorization key")
	}

	return base64.StdEncoding.EncodeToString(key), nil
}
