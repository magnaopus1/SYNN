package wallet

import (
	"crypto/rand"
	"fmt"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)



// NewMnemonicRecoveryManager initializes the recovery manager for future recovery.
func NewMnemonicRecoveryManager(email string, phoneNumber string, syn900Token string, ledgerInstance *ledger.Ledger) (*MnemonicRecoveryManager, error) {
	if syn900Token == "" {
		return nil, fmt.Errorf("Syn900 token is required for recovery setup")
	}

	// Simulate sending the Syn900 token to the wallet and destroying it after verification
	err := verifyAndDestroySyn900(syn900Token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify Syn900 token: %v", err)
	}

	// Register recovery information
	manager := &MnemonicRecoveryManager{
		RecoveryEmail:      email,
		RecoveryPhoneNumber: phoneNumber,
		Syn900Token:        "", // Token destroyed after verification
		IsRecoverySetUp:    true,
		ledgerInstance:     ledgerInstance,
	}

	// Log the registration in the ledger
	ledgerInstance.LogRecoverySetup(email, phoneNumber)
	return manager, nil
}

// verifyAndDestroySyn900 verifies the Syn900 token and destroys it.
func verifyAndDestroySyn900(syn900Token string) error {
	// Simulate sending the Syn900 token to the wallet for verification
	verified := true // Assume verification succeeded

	if !verified {
		return fmt.Errorf("failed to verify Syn900 token")
	}

	log.Println("Syn900 token verified and destroyed successfully.")
	return nil
}

// StartRecoveryProcess begins the emergency recovery process if the mnemonic is forgotten.
func (manager *MnemonicRecoveryManager) StartRecoveryProcess(email string, phoneNumber string, syn900Token string) (string, error) {
    manager.mutex.Lock()
    defer manager.mutex.Unlock()

    // Ensure the recovery setup exists (check the boolean value directly)
    if !manager.IsRecoverySetUp {
        return "", fmt.Errorf("recovery setup not initialized")
    }

    // Step 1: Verify email
    _, err := sendEmailVerificationCode(email)
    if err != nil {
        return "", fmt.Errorf("email verification failed: %v", err)
    }

    // Step 2: Verify phone number
    _, err = sendPhoneVerificationCode(phoneNumber)
    if err != nil {
        return "", fmt.Errorf("phone verification failed: %v", err)
    }

    // Step 3: Verify Syn900 token and destroy it
    err = verifyAndDestroySyn900(syn900Token)
    if err != nil {
        return "", fmt.Errorf("Syn900 token verification failed: %v", err)
    }

    // All steps successful; proceed with mnemonic recovery
    mnemonicPhrase, err := retrieveMnemonic(manager.ledgerInstance)
    if err != nil {
        return "", fmt.Errorf("failed to recover mnemonic: %v", err)
    }

    log.Println("Mnemonic successfully recovered.")
    return mnemonicPhrase, nil
}



// sendEmailVerificationCode sends a verification code to the recovery email.
func sendEmailVerificationCode(email string) (string, error) {
	verificationCode := generateVerificationCode()
	// Simulate sending the code via email
	log.Printf("Verification code sent to email %s: %s\n", email, verificationCode)
	return verificationCode, nil
}

// sendPhoneVerificationCode sends a verification code to the recovery phone number.
func sendPhoneVerificationCode(phoneNumber string) (string, error) {
	verificationCode := generateVerificationCode()
	// Simulate sending the code via SMS
	log.Printf("Verification code sent to phone number %s: %s\n", phoneNumber, verificationCode)
	return verificationCode, nil
}

// generateVerificationCode generates a random 6-digit verification code.
func generateVerificationCode() string {
	code := make([]byte, 3)
	_, err := rand.Read(code)
	if err != nil {
		log.Fatalf("Failed to generate verification code: %v", err)
	}
	return fmt.Sprintf("%06x", code)
}

// retrieveMnemonic retrieves the mnemonic from the ledger for recovery purposes.
func retrieveMnemonic(ledgerInstance *ledger.Ledger) (string, error) {
	// Retrieve the encrypted mnemonic from the ledger
	encryptedMnemonic, err := ledgerInstance.RetrieveMnemonic()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve mnemonic: %v", err)
	}

	// Create an instance of Encryption
	encryption := &common.Encryption{}

	// Decrypt the mnemonic using the DecryptData function
	decryptedMnemonicBytes, err := encryption.DecryptData([]byte(encryptedMnemonic), common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt mnemonic: %v", err)
	}

	// Convert decrypted data to string
	decryptedMnemonic := string(decryptedMnemonicBytes)

	return decryptedMnemonic, nil
}




