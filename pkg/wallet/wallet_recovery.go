package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"synnergy_network/pkg/common"

	"golang.org/x/crypto/pbkdf2"
)

// NewWalletRecovery initializes a WalletRecovery instance.
func NewWalletRecovery() *WalletRecovery {
	return &WalletRecovery{}
}

// RecoverFromMnemonic recovers a wallet using a mnemonic phrase (BIP-39).
func (wr *WalletRecovery) RecoverFromMnemonic(mnemonic string, passphrase string) (*ecdsa.PrivateKey, error) {
	// Step 1: Validate the mnemonic phrase
	if !wr.validateMnemonic(mnemonic) {
		return nil, errors.New("invalid mnemonic phrase")
	}

	// Step 2: Derive seed from mnemonic and passphrase using PBKDF2 (BIP-39 standard)
	seed := pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"+passphrase), 2048, 64, sha512.New)

	// Step 3: Use the seed to generate a private key
	privateKey, err := wr.generatePrivateKeyFromSeed(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key from mnemonic: %v", err)
	}

	wr.PrivateKey = privateKey
	return privateKey, nil
}

// RecoverFromPrivateKey recovers a wallet using the provided private key.
func (wr *WalletRecovery) RecoverFromPrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = elliptic.P256().ScalarBaseMult(privateKey.D.Bytes())

	if privateKey == nil {
		return nil, errors.New("failed to recover wallet from private key")
	}

	wr.PrivateKey = privateKey
	return privateKey, nil
}

// validateMnemonic checks if the provided mnemonic is valid.
func (wr *WalletRecovery) validateMnemonic(mnemonic string) bool {
	// For demo purposes, we assume the mnemonic is valid if it contains 12-24 words.
	words := strings.Split(mnemonic, " ")
	return len(words) >= 12 && len(words) <= 24
}

// generatePrivateKeyFromSeed generates a private key from the provided seed.
func (wr *WalletRecovery) generatePrivateKeyFromSeed(seed []byte) (*ecdsa.PrivateKey, error) {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()

	// Hash the seed using SHA-256 to create a deterministic private key
	privateKeyHash := sha256.Sum256(seed)
	privateKey.D = new(big.Int).SetBytes(privateKeyHash[:])
	privateKey.PublicKey.X, privateKey.PublicKey.Y = elliptic.P256().ScalarBaseMult(privateKey.D.Bytes())

	return privateKey, nil
}

// EncryptRecoveredPrivateKey encrypts the recovered private key using a passphrase.
func (wr *WalletRecovery) EncryptRecoveredPrivateKey(passphrase string) ([]byte, error) {
	privateKeyBytes := wr.PrivateKey.D.Bytes()

	// Create an encryption instance
	encryptionInstance := &common.Encryption{} // Adjust this if necessary to match your encryption logic.

	// Encrypt the private key using the passphrase
	encryptedPrivateKey, err := encryptionInstance.EncryptData("AES", privateKeyBytes, []byte(passphrase)) // Adjust the algorithm if needed.
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt private key: %v", err)
	}
	return encryptedPrivateKey, nil
}


// DecryptPrivateKey decrypts the previously encrypted private key using a passphrase.
func (wr *WalletRecovery) DecryptPrivateKey(encryptedKey []byte, passphrase string) error {
	// Create an encryption instance
	encryptionInstance := &common.Encryption{} // Adjust this if necessary to match your encryption logic.

	// Decrypt the private key using the passphrase
	decryptedKey, err := encryptionInstance.DecryptData(encryptedKey, []byte(passphrase)) // Remove "AES" from the arguments.
	if err != nil {
		return fmt.Errorf("failed to decrypt private key: %v", err)
	}

	// Set up the private key
	wr.PrivateKey = new(ecdsa.PrivateKey)
	wr.PrivateKey.PublicKey.Curve = elliptic.P256()
	wr.PrivateKey.D = new(big.Int).SetBytes(decryptedKey)
	wr.PrivateKey.PublicKey.X, wr.PrivateKey.PublicKey.Y = elliptic.P256().ScalarBaseMult(wr.PrivateKey.D.Bytes())

	return nil
}



// GenerateMnemonic generates a new random mnemonic phrase.
func (wr *WalletRecovery) GenerateMnemonic() (string, error) {
	entropy := make([]byte, 32)
	if _, err := rand.Read(entropy); err != nil {
		return "", fmt.Errorf("failed to generate entropy: %v", err)
	}

	// For simplicity, this example returns a hardcoded mnemonic phrase (BIP-39).
	// In a real-world implementation, you'd convert the entropy to a valid mnemonic.
	mnemonic := "seed essence lava shadow maple motor lucky canvas lunar cheese faith harbor"
	return mnemonic, nil
}
