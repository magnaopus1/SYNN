package syn3900

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
	"sync"

)

// Syn3900Transaction represents a transaction related to SYN3900 benefit tokens.
type Syn3900Transaction struct {
	TransactionID string    `json:"transaction_id"`
	Sender        string    `json:"sender"`           // Sender of the transaction
	Receiver      string    `json:"receiver"`         // Receiver of the transaction
	Amount        float64   `json:"amount"`           // Amount of benefit transferred
	Timestamp     time.Time `json:"timestamp"`        // Time of the transaction
	Signature     string    `json:"signature"`        // Digital signature for transaction verification
}

// TransactionManager handles SYN3900 token transactions.
type TransactionManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// NewTransactionManager creates a new TransactionManager instance.
func NewTransactionManager(ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *TransactionManager {
	return &TransactionManager{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// InitiateTransaction initiates a new SYN3900 transaction between sender and receiver.
func (tm *TransactionManager) InitiateTransaction(sender, receiver string, amount float64) (*Syn3900Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Create a new transaction
	tx := &Syn3900Transaction{
		TransactionID: generateUniqueTransactionID(),
		Sender:        sender,
		Receiver:      receiver,
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// Sign the transaction
	signature, err := tm.signTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}
	tx.Signature = signature

	// Store transaction in the ledger
	if err := tm.ledgerService.StoreTransaction(tx.TransactionID, tx); err != nil {
		return nil, fmt.Errorf("failed to store transaction in ledger: %w", err)
	}

	// Log the transaction initiation
	if err := tm.ledgerService.LogEvent("TransactionInitiated", time.Now(), tx.TransactionID); err != nil {
		return nil, fmt.Errorf("failed to log transaction: %w", err)
	}

	// Validate transaction using Synnergy Consensus
	if err := tm.consensusService.ValidateSubBlock(tx.TransactionID); err != nil {
		return nil, fmt.Errorf("failed to validate transaction: %w", err)
	}

	return tx, nil
}

// ValidateTransactionSignature validates the digital signature of a transaction.
func (tm *TransactionManager) ValidateTransactionSignature(tx *Syn3900Transaction) (bool, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Recreate the transaction data to get the hash
	txData := tx.TransactionID + tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Amount)
	hash := sha256.Sum256([]byte(txData))

	// Decode the signature from hex
	signatureBytes, err := hex.DecodeString(tx.Signature)
	if err != nil || len(signatureBytes) < 64 {
		return false, fmt.Errorf("invalid signature format: %w", err)
	}

	// Split the signature into r and s values
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])

	// Retrieve the public key for the sender
	publicKey, err := tm.getPublicKeyFromVault(tx.Sender)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve public key: %w", err)
	}

	// Verify the signature
	isValid := ecdsa.Verify(publicKey, hash[:], r, s)
	return isValid, nil
}

// signTransaction signs the transaction using the sender's private key.
func (tm *TransactionManager) signTransaction(tx *Syn3900Transaction) (string, error) {
	// Serialize the transaction into a hashable format (e.g., using the transaction ID, sender, receiver, and amount).
	txData := tx.TransactionID + tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Amount)
	hash := sha256.Sum256([]byte(txData))

	// Use the private key to sign the transaction hash
	privateKey, err := tm.getPrivateKeyFromVault(tx.Sender)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve private key: %w", err)
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Concatenate the r and s values as the signature
	signature := append(r.Bytes(), s.Bytes()...)
	return hex.EncodeToString(signature), nil
}

// getPrivateKeyFromVault retrieves the sender's private key securely from a key vault.
func (tm *TransactionManager) getPrivateKeyFromVault(sender string) (*ecdsa.PrivateKey, error) {
	// Retrieve the private key from a secure key vault (e.g., AWS KMS, Azure Key Vault, HSM, etc.)
	privateKey, err := securekeyvault.GetPrivateKey("key-id-for-" + sender)
	if err != nil {
		return nil, fmt.Errorf("error retrieving private key: %w", err)
	}
	return privateKey, nil
}

// getPublicKeyFromVault retrieves the sender's public key securely from a key vault.
func (tm *TransactionManager) getPublicKeyFromVault(sender string) (*ecdsa.PublicKey, error) {
	// Retrieve the public key from a secure key vault (e.g., AWS KMS, Azure Key Vault, HSM, etc.)
	publicKey, err := securekeyvault.GetPublicKey("key-id-for-" + sender)
	if err != nil {
		return nil, fmt.Errorf("error retrieving public key: %w", err)
	}
	return publicKey, nil
}

// RevokeTransaction revokes a SYN3900 transaction if it is found to be invalid or fraudulent.
func (tm *TransactionManager) RevokeTransaction(transactionID string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Log the transaction revocation in the ledger
	if err := tm.ledgerService.LogEvent("TransactionRevoked", time.Now(), transactionID); err != nil {
		return fmt.Errorf("failed to log transaction revocation: %w", err)
	}

	// Remove the transaction from the ledger
	if err := tm.ledgerService.DeleteTransaction(transactionID); err != nil {
		return fmt.Errorf("failed to revoke transaction: %w", err)
	}

	// Invalidate the transaction in the Synnergy Consensus
	if err := tm.consensusService.InvalidateSubBlock(transactionID); err != nil {
		return fmt.Errorf("failed to invalidate transaction with consensus: %w", err)
	}

	return nil
}

// generateUniqueTransactionID generates a unique identifier for each transaction.
func generateUniqueTransactionID() string {
	// Generate a unique transaction ID (could use UUID, timestamp, etc.)
	return "txn-" + time.Now().Format("20060102150405")
}
