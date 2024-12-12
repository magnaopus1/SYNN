package syn4200

import (
	"errors"
	"time"
	"sync"
)

// Syn4200Transaction represents a transaction involving SYN4200 tokens.
type Syn4200Transaction struct {
	TransactionID string    `json:"transaction_id"`
	Sender        string    `json:"sender"`       // Address of the sender (donor)
	Receiver      string    `json:"receiver"`     // Address of the receiver (charity)
	Amount        float64   `json:"amount"`       // Amount of tokens transferred
	Timestamp     time.Time `json:"timestamp"`    // Time when the transaction occurred
	Signature     string    `json:"signature"`    // Digital signature for transaction integrity
	Status        string    `json:"status"`       // Transaction status (pending, completed, failed)
}

// TransactionManager handles SYN4200 token transactions, ensuring secure processing, signing, and validation.
type TransactionManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex // To ensure thread safety
}

// NewTransactionManager initializes a new TransactionManager.
func NewTransactionManager(ledgerService *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *TransactionManager {
	return &TransactionManager{
		ledgerService:     ledgerService,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// CreateTransaction creates a new Syn4200 transaction, signs it, and stores it securely.
func (tm *TransactionManager) CreateTransaction(sender, receiver string, amount float64) (*Syn4200Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Generate a unique transaction ID
	txID := generateUniqueTransactionID()

	// Create the transaction object
	transaction := &Syn4200Transaction{
		TransactionID: txID,
		Sender:        sender,
		Receiver:      receiver,
		Amount:        amount,
		Timestamp:     time.Now(),
		Status:        "pending",
	}

	// Sign the transaction using the sender's private key
	signature, err := tm.signTransaction(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}
	transaction.Signature = signature

	// Encrypt the transaction before storing it
	encryptedTx, err := tm.encryptionService.EncryptData(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt transaction: %w", err)
	}

	// Store the transaction in the ledger
	if err := tm.ledgerService.StoreData(txID, encryptedTx); err != nil {
		return nil, fmt.Errorf("failed to store transaction in ledger: %w", err)
	}

	// Validate the transaction with Synnergy Consensus
	if err := tm.consensusService.ValidateSubBlock(txID); err != nil {
		return nil, fmt.Errorf("transaction validation failed: %w", err)
	}

	// Update the status to completed
	transaction.Status = "completed"

	// Log the transaction event
	tm.logTransactionEvent(txID, "TransactionCreated")

	return transaction, nil
}

// RetrieveTransaction retrieves a Syn4200 transaction from the ledger and decrypts it.
func (tm *TransactionManager) RetrieveTransaction(transactionID string) (*Syn4200Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the encrypted transaction data from the ledger
	encryptedTx, err := tm.ledgerService.RetrieveData(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction: %w", err)
	}

	// Decrypt the transaction data
	decryptedTx, err := tm.encryptionService.DecryptData(encryptedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt transaction: %w", err)
	}

	// Cast the decrypted data back to Syn4200Transaction
	transaction := decryptedTx.(*Syn4200Transaction)
	return transaction, nil
}

// signTransaction signs the transaction for verification purposes using a private key.
func (tm *TransactionManager) signTransaction(tx *Syn4200Transaction) (string, error) {
	// Serialize the transaction into a hashable format (e.g., using the transaction ID, sender, receiver, and amount).
	txData := tx.TransactionID + tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Amount)
	hash := sha256.Sum256([]byte(txData))

	// Use the private key to sign the transaction hash
	privateKey, err := tm.getPrivateKeyFromVault(tx.Sender)
	if err != nil {
		return "", fmt.Errorf("error retrieving private key: %w", err)
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("error signing transaction: %w", err)
	}

	// Concatenate the r and s values as the signature
	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	return hex.EncodeToString(signature), nil
}

// validateTransactionSignature validates the digital signature of the transaction using the public key.
func (tm *TransactionManager) validateTransactionSignature(tx *Syn4200Transaction) bool {
	// Recreate the transaction data to get the hash
	txData := tx.TransactionID + tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Amount)
	hash := sha256.Sum256([]byte(txData))

	// Decode the signature from hex
	signatureBytes, err := hex.DecodeString(tx.Signature)
	if err != nil || len(signatureBytes) < 64 {
		return false
	}

	// Split the signature into r and s values
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])

	// Retrieve the public key for the sender
	publicKey, err := tm.getPublicKeyFromVault(tx.Sender)
	if err != nil {
		return false
	}

	// Verify the signature
	return ecdsa.Verify(publicKey, hash[:], r, s)
}

// getPrivateKeyFromVault retrieves the private key securely from a key vault or HSM (real-world implementation).
func (tm *TransactionManager) getPrivateKeyFromVault(sender string) (*ecdsa.PrivateKey, error) {
	// Replace with actual key vault or HSM retrieval logic
	privateKey, err := securekeyvault.GetPrivateKey("key-id-for-" + sender)
	if err != nil {
		return nil, fmt.Errorf("error retrieving private key: %w", err)
	}
	return privateKey, nil
}

// getPublicKeyFromVault retrieves the public key securely from a key vault or HSM (real-world implementation).
func (tm *TransactionManager) getPublicKeyFromVault(sender string) (*ecdsa.PublicKey, error) {
	// Replace with actual key vault or HSM retrieval logic
	publicKey, err := securekeyvault.GetPublicKey("key-id-for-" + sender)
	if err != nil {
		return nil, fmt.Errorf("error retrieving public key: %w", err)
	}
	return publicKey, nil
}

// generateUniqueTransactionID generates a unique identifier for the transaction.
func generateUniqueTransactionID() string {
	// Implement a robust ID generation logic (e.g., UUID or timestamp-based ID)
	return fmt.Sprintf("tx-%d", time.Now().UnixNano())
}

// logTransactionEvent logs transaction-related events in the ledger for auditing purposes.
func (tm *TransactionManager) logTransactionEvent(transactionID, eventType string) {
	_ = tm.ledgerService.LogEvent(eventType, time.Now(), transactionID)
}
