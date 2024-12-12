package syn4300

import (
	"errors"
	"sync"
	"time"
)

// Syn4300Transaction represents a transaction for SYN4300 tokens (energy assets, RECs, or carbon credits).
type Syn4300Transaction struct {
	TransactionID     string    `json:"transaction_id"`   // Unique identifier for the transaction
	TokenID           string    `json:"token_id"`         // ID of the token being transacted
	Sender            string    `json:"sender"`           // Sender/Owner of the token
	Receiver          string    `json:"receiver"`         // Recipient of the token
	Quantity          float64   `json:"quantity"`         // Quantity of energy/asset being transacted
	Timestamp         time.Time `json:"timestamp"`        // Time when the transaction took place
	TransactionStatus string    `json:"transaction_status"`// Status of the transaction (e.g., pending, completed)
	Signature         string    `json:"signature"`        // Digital signature of the transaction for verification
}

// TransactionManager handles the creation, validation, and execution of SYN4300 transactions.
type TransactionManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// NewTransactionManager creates a new instance of TransactionManager.
func NewTransactionManager(ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *TransactionManager {
	return &TransactionManager{
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// CreateTransaction creates a new transaction for transferring SYN4300 tokens.
func (tm *TransactionManager) CreateTransaction(sender, receiver string, token *Syn4300Token, quantity float64) (*Syn4300Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Ensure sufficient quantity of the asset is available for the transaction
	if quantity > token.Metadata.Quantity {
		return nil, errors.New("insufficient token quantity for the transaction")
	}

	// Ensure the token is active
	if token.Metadata.Status != "active" {
		return nil, errors.New("token is not active, transaction cannot proceed")
	}

	// Create a new transaction
	tx := &Syn4300Transaction{
		TransactionID:     generateUniqueTransactionID(),
		TokenID:           token.TokenID,
		Sender:            sender,
		Receiver:          receiver,
		Quantity:          quantity,
		Timestamp:         time.Now(),
		TransactionStatus: "pending",
	}

	// Sign the transaction
	signature, err := tm.signTransaction(tx)
	if err != nil {
		return nil, err
	}
	tx.Signature = signature

	// Encrypt the transaction data
	encryptedTx, err := tm.encryptionService.EncryptData(tx)
	if err != nil {
		return nil, err
	}

	// Store the transaction in the ledger
	if err := tm.ledgerService.StoreTransaction(tx.TransactionID, encryptedTx); err != nil {
		return nil, err
	}

	// Log the transaction event in the ledger
	if err := tm.ledgerService.LogEvent("TransactionCreated", time.Now(), tx.TransactionID); err != nil {
		return nil, err
	}

	// Validate the transaction with Synnergy Consensus
	if err := tm.consensusService.ValidateSubBlock(tx.TransactionID); err != nil {
		return nil, err
	}

	return tx, nil
}

// ValidateTransaction validates the authenticity of a transaction.
func (tm *TransactionManager) ValidateTransaction(txID string) (*Syn4300Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the transaction from the ledger
	encryptedTx, err := tm.ledgerService.RetrieveTransaction(txID)
	if err != nil {
		return nil, err
	}

	// Decrypt the transaction data
	tx, err := tm.encryptionService.DecryptData(encryptedTx)
	if err != nil {
		return nil, err
	}

	// Validate the digital signature of the transaction
	if !tm.validateTransactionSignature(tx.(*Syn4300Transaction)) {
		return nil, errors.New("transaction signature validation failed")
	}

	// Validate with Synnergy Consensus
	if err := tm.consensusService.ValidateSubBlock(txID); err != nil {
		return nil, err
	}

	return tx.(*Syn4300Transaction), nil
}

// ExecuteTransaction finalizes and executes a transaction, updating token ownership and status.
func (tm *TransactionManager) ExecuteTransaction(txID string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve and validate the transaction
	tx, err := tm.ValidateTransaction(txID)
	if err != nil {
		return err
	}

	// Retrieve the token associated with the transaction
	token, err := tm.ledgerService.RetrieveToken(tx.TokenID)
	if err != nil {
		return err
	}

	// Update the token quantity and ownership
	token.Metadata.Quantity -= tx.Quantity
	token.Metadata.Owner = tx.Receiver

	// Re-encrypt and store the updated token
	encryptedToken, err := tm.encryptionService.EncryptData(token)
	if err != nil {
		return err
	}
	if err := tm.ledgerService.UpdateToken(token.(*Syn4300Token).TokenID, encryptedToken); err != nil {
		return err
	}

	// Update the transaction status to "completed"
	tx.TransactionStatus = "completed"
	encryptedTx, err := tm.encryptionService.EncryptData(tx)
	if err != nil {
		return err
	}
	if err := tm.ledgerService.UpdateTransaction(txID, encryptedTx); err != nil {
		return err
	}

	// Log the transaction execution in the ledger
	if err := tm.ledgerService.LogEvent("TransactionExecuted", time.Now(), txID); err != nil {
		return err
	}

	// Validate the transaction execution with Synnergy Consensus
	if err := tm.consensusService.ValidateSubBlock(txID); err != nil {
		return err
	}

	return nil
}

// signTransaction signs the transaction for verification purposes using a private key securely stored in a key vault.
func (tm *TransactionManager) signTransaction(tx *Syn4300Transaction) (string, error) {
	// Serialize the transaction into a hashable format (e.g., using the transaction ID, sender, receiver, and quantity).
	txData := tx.TransactionID + tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Quantity)
	hash := sha256.Sum256([]byte(txData))

	// Retrieve the sender's private key from a secure key vault or HSM
	privateKey, err := tm.getPrivateKeyFromVault(tx.Sender)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve private key: %w", err)
	}

	// Sign the transaction hash using ECDSA
	r, s, err := ecdsa.Sign(nil, privateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Concatenate the r and s values as the signature
	signature := append(r.Bytes(), s.Bytes()...)
	return hex.EncodeToString(signature), nil
}

// validateTransactionSignature validates the digital signature of the transaction using the public key retrieved securely.
func (tm *TransactionManager) validateTransactionSignature(tx *Syn4300Transaction) (bool, error) {
	// Recreate the transaction data to get the hash
	txData := tx.TransactionID + tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Quantity)
	hash := sha256.Sum256([]byte(txData))

	// Decode the signature from hex
	signatureBytes, err := hex.DecodeString(tx.Signature)
	if err != nil || len(signatureBytes) < 64 {
		return false, fmt.Errorf("invalid signature format")
	}

	// Split the signature into r and s values
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])

	// Retrieve the sender's public key from a secure source (e.g., key vault or HSM)
	publicKey, err := tm.getPublicKeyFromVault(tx.Sender)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve public key: %w", err)
	}

	// Verify the signature using ECDSA
	isValid := ecdsa.Verify(publicKey, hash[:], r, s)
	return isValid, nil
}

// getPrivateKeyFromVault retrieves the private key securely from a key vault or HSM (real-world implementation).
func (tm *TransactionManager) getPrivateKeyFromVault(sender string) (*ecdsa.PrivateKey, error) {
	// Fetch the private key from a secure key vault service (e.g., AWS KMS, Azure Key Vault, or HSM)
	keyID := "key-id-for-" + sender
	privateKeyPEM, err := securekeyvault.GetPrivateKey(keyID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving private key from vault for sender %s: %w", sender, err)
	}

	// Decode the PEM-encoded private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("invalid PEM format for private key")
	}

	// Parse the private key to ECDSA format
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key for sender %s: %w", sender, err)
	}

	return privateKey, nil
}

// getPublicKeyFromVault retrieves the public key securely from a key vault or HSM (real-world implementation).
func (tm *TransactionManager) getPublicKeyFromVault(sender string) (*ecdsa.PublicKey, error) {
	// Fetch the public key from a secure key vault service (e.g., AWS KMS, Azure Key Vault, or HSM)
	keyID := "key-id-for-" + sender

	// Example: Assume securekeyvault.GetPublicKey retrieves PEM-encoded public key from key vault
	publicKeyPEM, err := securekeyvault.GetPublicKey(keyID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving public key from vault for sender %s: %w", sender, err)
	}

	// Decode the PEM-encoded public key
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid PEM format for public key")
	}

	// Parse the public key to ECDSA format
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key for sender %s: %w", sender, err)
	}

	// Cast the public key to ECDSA format
	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("retrieved key is not an ECDSA public key")
	}

	return publicKey, nil
}

// generateUniqueTransactionID generates a unique identifier for a transaction.
func generateUniqueTransactionID() string {
	return "tx-id-" + time.Now().Format("20060102150405")
}
