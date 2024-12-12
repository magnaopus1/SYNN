package authorization

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"strings"
	"synnergy_network/pkg/ledger"
	"time"

	"github.com/google/uuid"
)

// PermissionSet defines the permissions associated with an authorized signer in the system.
type PermissionSet struct {
    CanApproveTransactions   bool // Permission to approve transactions
    CanModifyLedgerEntries   bool // Permission to modify ledger entries
    CanAccessSensitiveData   bool // Permission to access sensitive data within the system
    CanAddOrRemoveSigners    bool // Permission to add or remove authorized signers
    CanSetAuthorizationLevel bool // Permission to set authorization levels for other users
    CanFlagSuspiciousActivity bool // Permission to flag or review suspicious activities
}

// PermissionRequest represents a request for specific permissions by a user.
type PermissionRequest struct {
	RequestID            string        // Unique ID for the permission request
	UserID               string        // ID of the user requesting permissions
	RequestedPermissions PermissionSet // Permissions being requested
	RequestedAt          time.Time     // Time of the request
	Status               string        // Status of the request, e.g., "Pending", "Approved", "Denied"
}

// PermissionRevocation represents the revocation of specific permissions for a user.
type PermissionRevocation struct {
	UserID             string        // ID of the user whose permissions are revoked
	RevokedPermissions PermissionSet // Permissions being revoked
	RevokedAt          time.Time     // Time of the revocation
	Reason             string        // Reason for revocation, if applicable
}


// SubBlockValidationAuthorization records the authorization information for a user or node to validate sub-blocks in the network.
type SubBlockValidationAuthorization struct {
    AuthorizationID    string    // Unique identifier for the validation authorization record
    ValidatorID        string    // ID of the validator or node authorized for sub-block validation
    BlockID            string    // ID of the block or sub-block for which validation is authorized
    GrantedAt          time.Time // Timestamp of when the validation authorization was granted
    ValidUntil         time.Time // Expiration timestamp for the authorization
    AuthorizationLevel int       // Level of authorization required to perform validation
    GrantedBy          string    // ID of the user or entity granting the validation authorization
    Comments           string    // Additional notes or comments about the validation authorization
}

// AccessAttempt represents an attempt to access a system resource.
type AccessAttempt struct {
	AttemptID  string    // Unique identifier for the access attempt
	UserID     string    // ID of the user attempting access
	Action     string    // Action attempted by the user (e.g., "login", "dataAccess")
	Result     string    // Result of the attempt (e.g., "success", "failure")
	Timestamp  time.Time // Timestamp of the attempt
}


// AddAuthorizedSigner adds an authorized signer to the system, recording this in the ledger.
func AddAuthorizedSigner(ledgerInstance *ledger.Ledger, signerID string, permissions ledger.PermissionSet) error {
	// Define the authorized signer with the appropriate fields
	authSigner := ledger.AuthorizedSigner{
		SignerID:    signerID,
		Permissions: permissions,
		AddedAt:     time.Now(),
	}

	// Record the authorized signer in the ledger (no return value expected)
	ledgerInstance.AuthorizationLedger.RecordAuthorizedSigner(authSigner)

	return nil
}



// RemoveAuthorizedSigner removes an authorized signer from the system, updating the ledger.
func RemoveAuthorizedSigner(ledgerInstance *ledger.Ledger, signerID string) error {
	// Delete the authorized signer from the ledger (no value expected in return)
	ledgerInstance.AuthorizationLedger.DeleteAuthorizedSigner(signerID)

	return nil
}


// CheckAuthorization verifies that a signer has the necessary permissions to perform an action.
func CheckAuthorization(ledgerInstance *ledger.Ledger, signerID string, requiredPermissions ledger.PermissionSet) (bool, error) {
	// Fetch the signer's authorization details from the ledger
	signer, err := ledgerInstance.AuthorizationLedger.FetchAuthorizedSigner(signerID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve signer authorization from ledger: %v", err)
	}

	// Check if the signer has all required permissions
	if !signer.Permissions.HasAll(requiredPermissions) {
		return false, errors.New("signer does not have the required permissions")
	}

	return true, nil
}

// RequestUserPermissions allows users to request specific permissions, logging the request in the ledger.
func RequestUserPermissions(ledgerInstance *ledger.Ledger, userID string, requestedPermissions PermissionSet) error {
	// Create a new permission request record
	permissionRequest := PermissionRequest{
		RequestID:            generateRequestID(),
		UserID:               userID,
		RequestedPermissions: requestedPermissions,
		RequestedAt:          time.Now(),
		Status:               "Pending",
	}

	// Record the permission request in the ledger by passing individual fields
	if err := ledgerInstance.AuthorizationLedger.RecordPermissionRequest(
		permissionRequest.RequestID,
		permissionRequest.UserID,
		permissionRequest.RequestedPermissions.String(), // Now this should work
		permissionRequest.Status,
	); err != nil {
		return fmt.Errorf("failed to record permission request in ledger: %v", err)
	}

	return nil
}


// String converts the PermissionSet to a single string representation.
func (p PermissionSet) String() string {
    // Represent each permission in a "key=true/false" format
    var permissions []string

    permissions = append(permissions, fmt.Sprintf("CanApproveTransactions=%v", p.CanApproveTransactions))
    permissions = append(permissions, fmt.Sprintf("CanModifyLedgerEntries=%v", p.CanModifyLedgerEntries))
    permissions = append(permissions, fmt.Sprintf("CanAccessSensitiveData=%v", p.CanAccessSensitiveData))
    permissions = append(permissions, fmt.Sprintf("CanAddOrRemoveSigners=%v", p.CanAddOrRemoveSigners))
    permissions = append(permissions, fmt.Sprintf("CanSetAuthorizationLevel=%v", p.CanSetAuthorizationLevel))
    permissions = append(permissions, fmt.Sprintf("CanFlagSuspiciousActivity=%v", p.CanFlagSuspiciousActivity))

    // Join the permission strings with commas for a compact representation
    return strings.Join(permissions, ", ")
}


// RevokeUserPermissions revokes certain permissions for a user, recording the revocation in the ledger.
func RevokeUserPermissions(ledgerInstance *ledger.Ledger, userID string, permissions PermissionSet, reason string) error {
	// Create a PermissionRevocation record with the provided permissions
	revocation := PermissionRevocation{
		UserID:             userID,
		RevokedPermissions: permissions,
		RevokedAt:          time.Now(),
		Reason:             reason,
	}

	// Convert the PermissionSet to a string format for ledger storage
	permissionsString := permissions.String()

	// Record the revocation in the ledger using the required parameters
	ledgerInstance.AuthorizationLedger.RecordPermissionRevocation(revocation.UserID, permissionsString, revocation.Reason)

	return nil
}


// generateRequestID generates a unique ID for permission requests.
func generateRequestID() string {
	return fmt.Sprintf("REQ-%d", time.Now().UnixNano())
}


// AuthorizeSubBlockValidation grants permission to validate sub-blocks, recording the authorization in the ledger.
func AuthorizeSubBlockValidation(ledgerInstance *ledger.Ledger, signerID string) error {
	// Generate an authorization ID (assuming a function exists for unique ID generation)
	authID := generateAuthID()

	// Status for the authorization action, e.g., "Authorized"
	status := "Authorized"

	// Record the authorization in the ledger with specific parameters
	ledgerInstance.AuthorizationLedger.RecordSubBlockValidationAuthorization(authID, signerID, status)

	return nil
}

// generateAuthID generates a unique authorization ID by combining a UUID with a timestamp.
func generateAuthID() string {
	// Generate a UUID
	uuidPart := uuid.New()

	// Append a timestamp to enhance uniqueness and traceability
	timestampPart := time.Now().UnixNano()

	// Format the final ID
	authID := fmt.Sprintf("%s-%d", uuidPart.String(), timestampPart)

	return authID
}

// LogAccessAttempt records an access attempt in the ledger, noting success or failure.
func LogAccessAttempt(ledgerInstance *ledger.Ledger, userID string, action string, success bool) error {
	// Create an AccessAttempt record
	accessAttempt := ledger.AccessLog{ // Ensure AccessLog is from ledger package
		AccessID:   generateAuthID(),
		UserID:     userID,
		DeviceID:   "", // Set DeviceID if relevant or leave empty if it's not applicable
		Timestamp:  time.Now(),
		Success:    success,
		AccessType: action,
		FailureReason: func() string {
			if !success {
				return "authorization level mismatch or inactive status"
			}
			return ""
		}(),
	}

	// Log the access attempt in the ledger by passing the entire AccessLog object
	if err := ledgerInstance.AuthorizationLedger.RecordAccessAttempt(accessAttempt); err != nil {
		return fmt.Errorf("failed to record access attempt in ledger: %v", err)
	}

	return nil
}


// SetAccessControlFlag sets an access control flag for a user, storing it in the ledger.
func SetAccessControlFlag(ledgerInstance *ledger.Ledger, userID string, flag bool) error {
	// Update the access control flag for the user in the ledger
	ledgerInstance.AuthorizationLedger.UpdateAccessControlFlag(userID, flag)

	return nil
}



// CheckAccessControlFlag checks the access control flag for a user from the ledger.
func CheckAccessControlFlag(ledgerInstance *ledger.Ledger, userID string) (bool, error) {
	flag, err := ledgerInstance.AuthorizationLedger.FetchAccessControlFlag(userID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve access control flag for user %s from ledger: %v", userID, err)
	}

	return flag, nil
}



// encryptData encrypts sensitive data, such as authorization information, using AES-GCM for secure storage.
func encryptData(data []byte, key string) ([]byte, error) {
	hashKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hashKey[:])
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

	return gcm.Seal(nonce, nonce, data, nil), nil
}
