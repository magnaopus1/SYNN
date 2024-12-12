package syn2800

import (
    "errors"
    "time"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
    "sync"
)

// TransactionManager handles the lifecycle of SYN2800 token transactions.
type TransactionManager struct {
	mutex sync.Mutex
}

// NewTransactionManager creates a new instance of TransactionManager.
func NewTransactionManager() *TransactionManager {
	return &TransactionManager{}
}

// CreateTransaction initiates a new transaction for a life insurance policy token.
func (tm *TransactionManager) CreateTransaction(tokenID string, transactionType string, amount float64, initiator string) (*common.Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate the token and retrieve it
	token, err := tm.retrieveAndDecryptToken(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token: %v", err)
	}

	// Create the transaction object
	transaction := &common.Transaction{
		TransactionID:     generateTransactionID(),
		TokenID:           tokenID,
		TransactionType:   transactionType,
		Amount:            amount,
		Initiator:         initiator,
		Timestamp:         time.Now(),
		SubBlockValidated: false,
		BlockValidated:    false,
	}

	// Encrypt transaction data and store in the ledger
	encryptedTxData, err := tm.encryptTransactionData(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt transaction data: %v", err)
	}
	if err := ledger.StoreTransaction(transaction.TransactionID, encryptedTxData); err != nil {
		return nil, fmt.Errorf("failed to store transaction in ledger: %v", err)
	}

	log.Printf("Transaction %s created for token %s", transaction.TransactionID, tokenID)
	return transaction, nil
}

// ValidateSubBlock validates the transaction into a sub-block, using Synnergy Consensus.
func (tm *TransactionManager) ValidateSubBlock(transactionID string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the transaction from the ledger
	encryptedTxData, err := ledger.RetrieveTransaction(transactionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction: %v", err)
	}

	// Decrypt transaction data
	transaction, err := tm.decryptTransactionData(encryptedTxData)
	if err != nil {
		return fmt.Errorf("failed to decrypt transaction data: %v", err)
	}

	// Validate the transaction as part of a sub-block
	subBlockID := synnergy.ValidateIntoSubBlock(transaction)
	transaction.SubBlockValidated = true
	transaction.SubBlockID = subBlockID

	// Encrypt and store the updated transaction data
	encryptedTxData, err = tm.encryptTransactionData(transaction)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction data after sub-block validation: %v", err)
	}
	if err := ledger.StoreTransaction(transaction.TransactionID, encryptedTxData); err != nil {
		return fmt.Errorf("failed to store updated transaction in ledger: %v", err)
	}

	log.Printf("Transaction %s validated into sub-block %s", transaction.TransactionID, subBlockID)
	return nil
}

// ValidateBlock finalizes the validation of transactions into a block after sub-block validation.
func (tm *TransactionManager) ValidateBlock(subBlockID string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve all transactions in the sub-block
	transactionIDs, err := synnergy.GetTransactionsInSubBlock(subBlockID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transactions in sub-block: %v", err)
	}

	for _, txID := range transactionIDs {
		// Retrieve and decrypt the transaction
		encryptedTxData, err := ledger.RetrieveTransaction(txID)
		if err != nil {
			return fmt.Errorf("failed to retrieve transaction %s: %v", txID, err)
		}

		transaction, err := tm.decryptTransactionData(encryptedTxData)
		if err != nil {
			return fmt.Errorf("failed to decrypt transaction %s: %v", txID, err)
		}

		// Mark the transaction as validated into the final block
		transaction.BlockValidated = true
		transaction.BlockID = subBlockID

		// Encrypt and store the updated transaction
		encryptedTxData, err = tm.encryptTransactionData(transaction)
		if err != nil {
			return fmt.Errorf("failed to encrypt transaction data after block validation: %v", err)
		}
		if err := ledger.StoreTransaction(txID, encryptedTxData); err != nil {
			return fmt.Errorf("failed to store updated transaction in ledger: %v", err)
		}

		log.Printf("Transaction %s validated into block %s", txID, subBlockID)
	}

	return nil
}

// CancelTransaction allows the cancellation of a pending transaction before sub-block validation.
func (tm *TransactionManager) CancelTransaction(transactionID string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the transaction
	encryptedTxData, err := ledger.RetrieveTransaction(transactionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction: %v", err)
	}

	// Decrypt transaction data
	transaction, err := tm.decryptTransactionData(encryptedTxData)
	if err != nil {
		return fmt.Errorf("failed to decrypt transaction data: %v", err)
	}

	// Ensure the transaction has not yet been validated
	if transaction.SubBlockValidated {
		return fmt.Errorf("transaction %s cannot be canceled, it is already validated into a sub-block", transactionID)
	}

	// Remove the transaction from the ledger
	if err := ledger.RemoveTransaction(transactionID); err != nil {
		return fmt.Errorf("failed to remove transaction: %v", err)
	}

	log.Printf("Transaction %s canceled successfully", transactionID)
	return nil
}

// Helper function to generate a unique transaction ID.
func generateTransactionID() string {
	return fmt.Sprintf("TX-%d", time.Now().UnixNano())
}

// Helper function to retrieve and decrypt the token from the ledger.
func (tm *TransactionManager) retrieveAndDecryptToken(tokenID string) (*common.SYN2800Token, error) {
	encryptedTokenData, err := ledger.RetrieveToken(tokenID)
	if err != nil {
		return nil, err
	}
	return tm.decryptTokenData(encryptedTokenData)
}

// Helper function to encrypt transaction data.
func (tm *TransactionManager) encryptTransactionData(transaction *common.Transaction) ([]byte, error) {
	key := generateEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	transactionData := serializeTransactionData(transaction)
	return gcm.Seal(nonce, nonce, transactionData, nil), nil
}

// Helper function to decrypt transaction data.
func (tm *TransactionManager) decryptTransactionData(encryptedData []byte) (*common.Transaction, error) {
	key := generateEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return deserializeTransactionData(decryptedData), nil
}

// Helper function to serialize transaction data.
func serializeTransactionData(transaction *common.Transaction) []byte {
	data, err := json.Marshal(transaction)
	if err != nil {
		log.Fatalf("failed to serialize transaction data: %v", err)
	}
	return data
}

// Helper function to deserialize transaction data after decryption.
func deserializeTransactionData(data []byte) *common.Transaction {
	var transaction common.Transaction
	if err := json.Unmarshal(data, &transaction); err != nil {
		log.Fatalf("failed to deserialize transaction data: %v", err)
	}
	return &transaction
}

// Helper function to generate an encryption key.
func generateEncryptionKey() []byte {
	return []byte("your-secure-256-bit-key")
}
