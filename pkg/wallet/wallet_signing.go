package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewWalletSigner initializes a new WalletSigner with the provided private key.
func NewWalletSigner(privateKey *ecdsa.PrivateKey, ledgerInstance *ledger.Ledger) *WalletSigner {
	return &WalletSigner{
		PrivateKey: privateKey,
		Ledger:     ledgerInstance,
	}
}

// SignTransaction signs the given transaction and attaches the signature.
func (ws *WalletSigner) SignTransaction(tx *common.Transaction) error {
	txHash := sha256.Sum256([]byte(tx.String()))
	r, s, err := ecdsa.Sign(rand.Reader, ws.PrivateKey, txHash[:])
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Attach signature (r, s) as a string to the transaction
	signatureBytes := append(r.Bytes(), s.Bytes()...)
	tx.Signature = string(signatureBytes)
	return nil
}


// VerifyTransactionSignature verifies the signature of a given transaction.
func (ws *WalletSigner) VerifyTransactionSignature(tx *common.Transaction) (bool, error) {
	// Check if the transaction signature exists
	if tx.Signature == "" {
		return false, errors.New("transaction has no signature")
	}

	pubKey := &ws.PrivateKey.PublicKey
	txHash := sha256.Sum256([]byte(tx.String()))

	// Split the signature into r and s from the string (even-length assumption)
	signatureLen := len(tx.Signature)
	if signatureLen%2 != 0 {
		return false, errors.New("invalid signature length")
	}

	// Convert the signature to bytes for cryptographic verification
	r := new(big.Int).SetBytes([]byte(tx.Signature[:signatureLen/2]))
	s := new(big.Int).SetBytes([]byte(tx.Signature[signatureLen/2:]))

	// Verify the signature
	valid := ecdsa.Verify(pubKey, txHash[:], r, s)
	if !valid {
		return false, errors.New("invalid transaction signature")
	}

	return true, nil
}



// SignContractExecution signs the execution of a smart contract by hashing the contract data.
func (ws *WalletSigner) SignContractExecution(contractExecutionData []byte) ([]byte, error) {
	dataHash := sha256.Sum256(contractExecutionData)
	r, s, err := ecdsa.Sign(rand.Reader, ws.PrivateKey, dataHash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign contract execution: %v", err)
	}

	// Combine r and s into a single byte array
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

// VerifyContractSignature verifies the signature of a smart contract execution.
func (ws *WalletSigner) VerifyContractSignature(contractData []byte, signature []byte) (bool, error) {
	if signature == nil {
		return false, errors.New("contract execution has no signature")
	}

	pubKey := &ws.PrivateKey.PublicKey
	dataHash := sha256.Sum256(contractData)

	// Split the signature into r and s
	r := new(big.Int).SetBytes(signature[:len(signature)/2])
	s := new(big.Int).SetBytes(signature[len(signature)/2:])

	// Verify the signature
	valid := ecdsa.Verify(pubKey, dataHash[:], r, s)
	if !valid {
		return false, errors.New("invalid contract signature")
	}

	return true, nil
}

// EncryptPrivateKey encrypts the private key for secure storage.
func (ws *WalletSigner) EncryptPrivateKey(passphrase string) ([]byte, error) {
	// Create an encryption instance
	encryptionInstance := &common.Encryption{} // Adjust this based on your actual encryption implementation.

	// Convert the private key to bytes
	privateKeyBytes := ws.PrivateKey.D.Bytes()

	// Encrypt the private key using the passphrase and a chosen algorithm, e.g., AES
	encryptedPrivateKey, err := encryptionInstance.EncryptData("AES", []byte(passphrase), privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt private key: %v", err)
	}

	return encryptedPrivateKey, nil
}


// DecryptPrivateKey decrypts the private key using the given passphrase.
func (ws *WalletSigner) DecryptPrivateKey(encryptedData []byte, passphrase string) error {
	// Create an encryption instance
	encryptionInstance := &common.Encryption{} // Adjust this based on your actual encryption implementation.

	// Decrypt the private key using the passphrase
	decryptedData, err := encryptionInstance.DecryptData([]byte(passphrase), encryptedData)
	if err != nil {
		return fmt.Errorf("failed to decrypt private key: %v", err)
	}

	// Set up the private key
	ws.PrivateKey = new(ecdsa.PrivateKey)
	ws.PrivateKey.PublicKey.Curve = elliptic.P256()
	ws.PrivateKey.D = new(big.Int).SetBytes(decryptedData)

	return nil
}


// SignMessage signs an arbitrary message using the private key.
func (ws *WalletSigner) SignMessage(message []byte) ([]byte, error) {
	messageHash := sha256.Sum256(message)
	r, s, err := ecdsa.Sign(rand.Reader, ws.PrivateKey, messageHash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %v", err)
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

// VerifyMessageSignature verifies a message signature using the public key.
func (ws *WalletSigner) VerifyMessageSignature(message []byte, signature []byte) (bool, error) {
	if signature == nil {
		return false, errors.New("message has no signature")
	}

	pubKey := &ws.PrivateKey.PublicKey
	messageHash := sha256.Sum256(message)

	// Split the signature into r and s
	r := new(big.Int).SetBytes(signature[:len(signature)/2])
	s := new(big.Int).SetBytes(signature[len(signature)/2:])

	// Verify the signature
	valid := ecdsa.Verify(pubKey, messageHash[:], r, s)
	if !valid {
		return false, errors.New("invalid message signature")
	}

	return true, nil
}

// SendSignedTransaction sends a signed transaction to the ledger and validates it via Synnergy Consensus.
func (ws *WalletSigner) SendSignedTransaction(tx *common.Transaction) error {
	// Verify the transaction before sending it
	isValid, err := ws.VerifyTransactionSignature(tx)
	if err != nil {
		return fmt.Errorf("failed to verify transaction signature: %v", err)
	}
	if !isValid {
		return errors.New("transaction signature is invalid")
	}

	// Validate transaction through the Synnergy Consensus mechanism
	err = ws.Ledger.ValidateTransaction(tx.TransactionID)
	if err != nil {
		return fmt.Errorf("failed to validate transaction: %v", err)
	}

	// Log the transaction to the ledger (adjusted to pass string arguments)
	err = ws.Ledger.LogTransaction(tx.FromAddress, tx.TransactionID)
	if err != nil {
		return fmt.Errorf("failed to log transaction in the ledger: %v", err)
	}

	fmt.Printf("Transaction %s has been successfully sent to the ledger.\n", tx.TransactionID)
	return nil
}

