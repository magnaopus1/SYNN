package authority_node

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
)

// AuthorityNodeKey represents a generated key for a specific node type.
type AuthorityNodeKey struct {
	KeyID          string    // Unique key identifier
	NodeType       string    // Type of node the key is associated with
	OwnerName      string    // Owner of the key
	OwnerWallet    string    // Wallet address associated with the key
	ExpirationDate time.Time // Key expiration date (30 months from creation)
	MaxNodes       int       // Maximum number of nodes this key can start
	UsedNodes      int       // Number of nodes started with this key
	IsExpired      bool      // Whether the key has expired
}

// AuthorityNodeAccountManager manages authority node key ownership, wallet registration, key cancellation, and refreshing.
type AuthorityNodeAccountManager struct {
	mutex             sync.Mutex                 // Mutex for thread-safe operations
	Ledger            *ledger.Ledger             // Reference to the ledger for storing key details
	EncryptionService *encryption.Encryption     // Encryption service for secure key management
	AuthorityKeys     map[string]*AuthorityNodeKey // Map of keys by key ID
}

// NewAuthorityNodeAccountManager initializes a new AuthorityNodeAccountManager.
func NewAuthorityNodeAccountManager(ledger *ledger.Ledger, encryptionService *encryption.Encryption) *AuthorityNodeAccountManager {
	return &AuthorityNodeAccountManager{
		Ledger:            ledger,
		EncryptionService: encryptionService,
		AuthorityKeys:     make(map[string]*AuthorityNodeKey),
	}
}

// RegisterWalletToKey associates a wallet with the key ownership.
func (am *AuthorityNodeAccountManager) RegisterWalletToKey(keyID, ownerName, walletAddress string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	key, exists := am.AuthorityKeys[keyID]
	if !exists {
		return errors.New("key not found")
	}

	// Check if the wallet address is valid and not already registered.
	if key.OwnerWallet != "" {
		return errors.New("wallet already registered to this key")
	}

	// Register the wallet address and store the update in the ledger.
	key.OwnerWallet = walletAddress
	key.OwnerName = ownerName

	// Update the ledger with the new ownership details.
	err := am.Ledger.RecordKeyOwnershipUpdate(keyID, walletAddress, ownerName)
	if err != nil {
		return fmt.Errorf("failed to register wallet: %v", err)
	}

	fmt.Printf("Wallet %s registered to key %s for owner %s.\n", walletAddress, keyID, ownerName)
	return nil
}

// ChangeOwnershipDetails allows changing the key ownership details, including wallet and owner name.
func (am *AuthorityNodeAccountManager) ChangeOwnershipDetails(keyID, newOwnerName, newWalletAddress string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	key, exists := am.AuthorityKeys[keyID]
	if !exists {
		return errors.New("key not found")
	}

	// Check if the new wallet address is different from the existing one.
	if newWalletAddress == key.OwnerWallet {
		return errors.New("new wallet address is the same as the existing one")
	}

	// Update the ownership details.
	key.OwnerWallet = newWalletAddress
	key.OwnerName = newOwnerName

	// Store the updated details in the ledger.
	err := am.Ledger.RecordKeyOwnershipUpdate(keyID, newWalletAddress, newOwnerName)
	if err != nil {
		return fmt.Errorf("failed to change ownership details: %v", err)
	}

	fmt.Printf("Ownership details changed for key %s. New owner: %s, New wallet: %s.\n", keyID, newOwnerName, newWalletAddress)
	return nil
}

// CancelKey cancels a key, marking it as expired and removing its validity in the ledger.
func (am *AuthorityNodeAccountManager) CancelKey(keyID string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	key, exists := am.AuthorityKeys[keyID]
	if !exists {
		return errors.New("key not found")
	}

	// Mark the key as expired and update the ledger.
	key.IsExpired = true
	err := am.Ledger.RecordKeyCancellation(keyID)
	if err != nil {
		return fmt.Errorf("failed to cancel key: %v", err)
	}

	fmt.Printf("Key %s has been successfully canceled.\n", keyID)
	return nil
}

// RefreshKey refreshes a key at the 29-month mark, extending its expiration by another 30 months.
func (am *AuthorityNodeAccountManager) RefreshKey(keyID string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	key, exists := am.AuthorityKeys[keyID]
	if !exists {
		return errors.New("key not found")
	}

	// Check if the key is close to expiration (at the 29-month mark).
	if time.Until(key.ExpirationDate) > 31*24*time.Hour {
		return errors.New("key is not yet eligible for refresh")
	}

	// Extend the key's expiration date by another 30 months.
	key.ExpirationDate = time.Now().Add(30 * 24 * time.Hour * 30)
	key.IsExpired = false

	// Update the ledger with the new expiration date.
	err := am.Ledger.UpdateKeyExpiration(keyID, key.ExpirationDate)
	if err != nil {
		return fmt.Errorf("failed to refresh key expiration in ledger: %v", err)
	}

	fmt.Printf("Key %s has been successfully refreshed for another 30 months.\n", keyID)
	return nil
}

// ViewKeyDetails retrieves the full details of a specific key.
func (am *AuthorityNodeAccountManager) ViewKeyDetails(keyID string) (*AuthorityNodeKey, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	key, exists := am.AuthorityKeys[keyID]
	if !exists {
		return nil, errors.New("key not found")
	}

	return key, nil
}

// FileNodeActivity logs activity or updates related to the authority node key.
func (am *AuthorityNodeAccountManager) FileNodeActivity(keyID, activityDescription string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Retrieve the key.
	key, exists := am.AuthorityKeys[keyID]
	if !exists {
		return errors.New("key not found")
	}

	// Encrypt the activity log.
	encryptedActivity, err := am.EncryptionService.EncryptData([]byte(activityDescription), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt activity log: %v", err)
	}

	// Record the activity in the ledger.
	err = am.Ledger.RecordNodeActivity(keyID, encryptedActivity)
	if err != nil {
		return fmt.Errorf("failed to log node activity: %v", err)
	}

	fmt.Printf("Activity logged for key %s: %s\n", keyID, activityDescription)
	return nil
}

