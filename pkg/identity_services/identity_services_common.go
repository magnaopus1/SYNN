package identity_services

import (
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// AccessControlManager manages access control for authority nodes, users, and stakeholders
type AccessControlManager struct {
    AuthorityNodes map[string]*AuthorityNodeStruct  // Change the value type to *AuthorityNodeStruct
	UserAccess       map[string]string         // Map of UserID to their access level (encrypted)
	StakeholderAccess map[string]string        // Map of StakeholderID to their access level (encrypted)
	OwnerAccess      string                    // Owner's access secret
	LedgerInstance   *ledger.Ledger            // Ledger for recording access changes
	mutex            sync.Mutex                // Mutex for thread-safe access control
	Encryption      *common.Encryption    // Reference to the encryption service

}

// Identity represents an individual's identity on the blockchain
type Identity struct {
	IdentityID    string       // Unique ID for the identity
	IdentityType  IdentityType // The type of identity (Syn900 or DecentralizedID)
	Owner         string       // Owner of the identity (wallet address or public key)
	CreatedAt     time.Time    // Timestamp of identity creation
	EncryptedData string       // Encrypted identity details
	IsVerified    bool         // Whether the identity has been verified
}

// IdentityType defines the types of identities in the system
type IdentityType string

// IdentityVerificationManager handles creation and verification of identities
type IdentityVerificationManager struct {
	Identities     map[string]*Identity  // Map of identityID to identity
	LedgerInstance *ledger.Ledger        // Ledger for recording identity-related actions
	mutex          sync.Mutex            // Mutex for thread-safe identity operations
	Encryption      *common.Encryption    // Reference to the encryption service

}

// PrivacySettings defines the settings for user privacy
type PrivacySettings struct {
	UserID            string    // User's ID (wallet address, public key, or other unique identifier)
	DataEncryption    bool      // Whether the user's data is encrypted
	PermissionToShare bool      // Whether the user permits sharing their data
	LastUpdated       time.Time // The last time the privacy settings were updated
}

// PrivacyManager handles the management of user privacy settings
type PrivacyManager struct {
	PrivacyRecords  map[string]*PrivacySettings // Map of user ID to privacy settings
	mutex           sync.Mutex                  // Mutex for thread-safe privacy operations
	LedgerInstance  *ledger.Ledger              // Ledger instance for logging privacy actions
	Encryption      *common.Encryption    // Reference to the encryption service

}
