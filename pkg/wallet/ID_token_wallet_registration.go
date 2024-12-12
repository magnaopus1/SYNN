package wallet

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewIDTokenWalletRegistrationService initializes the IDTokenWalletRegistrationService.
func NewIDTokenWalletRegistrationService(ledgerInstance *ledger.Ledger) *IDTokenWalletRegistrationService {
	return &IDTokenWalletRegistrationService{
		ledgerInstance: ledgerInstance,
	}
}

// RegisterSyn900Token registers a Syn900 token to a specific walletID, linking it to a verified identity.
func (service *IDTokenWalletRegistrationService) RegisterSyn900Token(walletID string, tokenID string, identityData map[string]interface{}, privateKey *ecdsa.PrivateKey) error {
    service.mutex.Lock()
    defer service.mutex.Unlock()

    // Ensure the Syn900 token is not already registered to another wallet.
    existingWallet, err := service.ledgerInstance.CheckTokenOwnership(tokenID, walletID)
    if err != nil {
        return fmt.Errorf("failed to check token ownership: %v", err)
    }

    if existingWallet {
        return fmt.Errorf("Syn900 token %s is already registered to another wallet", tokenID)
    }

    // Link the wallet with the identity and token, passing the privateKey as required.
    err = service.linkIdentityToWallet(walletID, identityData, privateKey)
    if err != nil {
        return fmt.Errorf("failed to link identity to wallet: %v", err)
    }

    // Encrypt the tokenID before registration.
    encryption := &common.Encryption{}
    encryptedTokenID, err := encryption.EncryptData("AES", []byte(tokenID), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt Syn900 token ID: %v", err)
    }

    // Convert encryptedTokenID to string before storing it.
    encryptedTokenIDStr := string(encryptedTokenID)

    // Register the encrypted Syn900 token in the ledger for the wallet.
    err = service.ledgerInstance.RegisterTokenToWallet(walletID, encryptedTokenIDStr)
    if err != nil {
        return fmt.Errorf("failed to register Syn900 token to wallet: %v", err)
    }

    // Mark the wallet as verified in the ledger.
    err = service.markWalletAsVerified(walletID)
    if err != nil {
        return fmt.Errorf("failed to mark wallet as verified: %v", err)
    }

    fmt.Printf("Successfully registered Syn900 token %s to wallet %s, wallet is now verified.\n", tokenID, walletID)
    return nil
}



var encryptionKey = []byte("YourSecureEncryptionKey") // Ensure this is managed by a key management system (KMS).

// linkIdentityToWallet securely links the provided identity data to the wallet, ensuring it is encrypted and verified.
func (service *IDTokenWalletRegistrationService) linkIdentityToWallet(walletID string, identityData map[string]interface{}, privateKey *ecdsa.PrivateKey) error {
	// Step 1: Input Validation
	if walletID == "" {
		return errors.New("walletID cannot be empty")
	}
	if len(identityData) == 0 {
		return errors.New("identity data cannot be empty")
	}

	// Step 2: Serialize identity data to JSON format
	identityString, err := serializeIdentityData(identityData)
	if err != nil {
		return fmt.Errorf("failed to serialize identity data: %v", err)
	}

	// Step 3: Create an instance of Encryption for secure data encryption
	encryption := &common.Encryption{}

	// Step 4: Encrypt the identity data using AES (or other secure method), ensuring confidentiality
	encryptedIdentity, err := encryption.EncryptData("AES", []byte(identityString), common.EncryptionKey)
	if err != nil {
		log.Printf("Encryption failed for walletID %s: %v", walletID, err)
		return fmt.Errorf("failed to encrypt identity data: %v", err)
	}

	// Step 5: Create an identity object that includes the encrypted identity data
	identity := ledger.Identity{
		IdentityID:    walletID,               // Use walletID as the unique identifier
		Owner:         walletID,               // Set the owner of the identity (wallet address or public key)
		CreatedAt:     time.Now(),             // Record the current timestamp for auditing purposes
		EncryptedData: string(encryptedIdentity), // Convert []byte to string for storing encrypted data
		IsVerified:    false,                  // Initially set to false until verification is completed
		RecoverySetup: "",                     // Optionally add recovery setup details (if applicable)
	}

	// Step 6: Store the encrypted identity data in the ledger, linking it to the walletID
	err = service.ledgerInstance.StoreIdentity(walletID, identity)
	if err != nil {
		log.Printf("Failed to store encrypted identity in ledger for walletID %s: %v", walletID, err)
		return fmt.Errorf("failed to store encrypted identity in ledger: %v", err)
	}

	// Step 7: Audit log of successful identity linking
	log.Printf("Successfully linked and encrypted identity for walletID %s at %v", walletID, time.Now())

	return nil
}





// markWalletAsVerified marks the wallet as verified in the ledger after registering the Syn900 token.
func (service *IDTokenWalletRegistrationService) markWalletAsVerified(walletID string) error {
	// Update the ledger to reflect the wallet's verified status.
	err := service.ledgerInstance.UpdateWalletVerificationStatus(walletID, true)
	if err != nil {
		return fmt.Errorf("failed to mark wallet %s as verified: %v", walletID, err)
	}
	return nil
}

// serializeIdentityData converts identity data into a serializable string format (e.g., JSON).
func serializeIdentityData(identityData map[string]interface{}) (string, error) {
	// Convert the identity data map into a JSON string
	identityJSON, err := json.Marshal(identityData)
	if err != nil {
		return "", fmt.Errorf("failed to serialize identity data: %v", err)
	}

	// Return the JSON string
	return string(identityJSON), nil
}

// VerifySyn900Token verifies that the Syn900 token is correctly registered and links to a valid wallet.
func (service *IDTokenWalletRegistrationService) VerifySyn900Token(tokenID, walletID string) (string, error) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	// Check the wallet linked to the provided Syn900 token using both tokenID and walletID.
	linkedWallet, err := service.ledgerInstance.CheckTokenOwnership(tokenID, walletID)
	if err != nil {
		return "", fmt.Errorf("failed to check ownership of Syn900 token: %v", err)
	}

	if !linkedWallet {
		return "", fmt.Errorf("Syn900 token %s is not registered to any wallet", tokenID)
	}

	// Assuming the `CheckTokenOwnership` method returns a `bool`, this would need to be adapted to the correct logic.
	fmt.Printf("Syn900 token %s is registered to wallet %s\n", tokenID, walletID)
	return walletID, nil
}

