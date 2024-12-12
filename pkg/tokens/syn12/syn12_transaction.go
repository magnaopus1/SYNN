package syn12

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// SYN12Transaction represents the structure for a SYN12 token transaction.
type SYN12Transaction struct {
	TokenID         string  // Unique ID of the SYN12 token
	FromAddress     string  // Sender's address
	ToAddress       string  // Receiver's address
	Amount          uint64  // Amount of tokens being transferred
	Timestamp       int64   // Timestamp of the transaction
	TransactionHash string  // Hash of the transaction
	IsEncrypted     bool    // Indicates if the transaction is encrypted
}

// SYN12TransactionManager handles SYN12 token transactions.
type SYN12TransactionManager struct {
	ledgerManager     *ledger.LedgerManager         // Ledger for recording transactions
	encryptionService *encryption.EncryptionService // Encryption service for secure transactions
	consensus         *consensus.SynnergyConsensus  // Consensus engine for transaction validation
	mutex             sync.Mutex                    // Mutex for transaction handling
}

// NewSYN12TransactionManager initializes the transaction manager for SYN12 tokens.
func NewSYN12TransactionManager(ledgerManager *ledger.LedgerManager, encryptionService *encryption.EncryptionService, consensus *consensus.SynnergyConsensus) *SYN12TransactionManager {
	return &SYN12TransactionManager{
		ledgerManager:     ledgerManager,
		encryptionService: encryptionService,
		consensus:         consensus,
	}
}

// CreateTransaction creates a new transaction for SYN12 tokens.
func (tm *SYN12TransactionManager) CreateTransaction(tokenID, fromAddress, toAddress string, amount uint64) (*SYN12Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate the transaction inputs
	if tokenID == "" || fromAddress == "" || toAddress == "" || amount <= 0 {
		return nil, errors.New("invalid transaction parameters")
	}

	// Validate the transaction through the consensus engine
	if err := tm.consensus.ValidateTransfer(fromAddress, toAddress, amount); err != nil {
		return nil, fmt.Errorf("transaction validation failed: %v", err)
	}

	// Create the transaction structure
	transaction := &SYN12Transaction{
		TokenID:     tokenID,
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
		Timestamp:   common.GetCurrentTimestamp(),
	}

	// Generate a transaction hash for uniqueness
	transaction.TransactionHash = tm.generateTransactionHash(transaction)

	// Encrypt the transaction for security
	encryptedTransaction, err := tm.encryptTransaction(transaction)
	if err != nil {
		return nil, fmt.Errorf("transaction encryption failed: %v", err)
	}
	transaction.IsEncrypted = true

	// Store the encrypted transaction in the ledger
	if err := tm.ledgerManager.RecordTransaction(transaction.TokenID, encryptedTransaction); err != nil {
		return nil, fmt.Errorf("ledger recording failed: %v", err)
	}

	// Log the transaction event
	if err := tm.ledgerManager.LogEvent(transaction.TokenID, common.EventTokenTransferred); err != nil {
		return nil, fmt.Errorf("failed to log transaction event: %v", err)
	}

	fmt.Printf("Transaction created successfully. Hash: %s\n", transaction.TransactionHash)
	return transaction, nil
}

// ValidateTransaction verifies the validity of a transaction.
func (tm *SYN12TransactionManager) ValidateTransaction(transaction *SYN12Transaction) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Decrypt the transaction if encrypted
	if transaction.IsEncrypted {
		decryptedTransaction, err := tm.decryptTransaction(transaction)
		if err != nil {
			return fmt.Errorf("transaction decryption failed: %v", err)
		}
		transaction = decryptedTransaction
	}

	// Validate the transaction through consensus
	if err := tm.consensus.ValidateTransfer(transaction.FromAddress, transaction.ToAddress, transaction.Amount); err != nil {
		return fmt.Errorf("consensus validation failed: %v", err)
	}

	// Ensure the transaction hash matches the computed hash
	computedHash := tm.generateTransactionHash(transaction)
	if transaction.TransactionHash != computedHash {
		return errors.New("transaction hash mismatch")
	}

	fmt.Printf("Transaction with Hash: %s validated successfully.\n", transaction.TransactionHash)
	return nil
}

// RetrieveTransaction fetches and decrypts a transaction by its hash.
func (tm *SYN12TransactionManager) RetrieveTransaction(transactionHash string) (*SYN12Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the encrypted transaction from the ledger
	encryptedTransaction, err := tm.ledgerManager.GetTransactionByHash(transactionHash)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %v", err)
	}

	// Decrypt the transaction
	transaction, err := tm.decryptTransaction(encryptedTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt transaction: %v", err)
	}

	fmt.Printf("Transaction with Hash: %s retrieved successfully.\n", transactionHash)
	return transaction, nil
}

// GenerateSubBlock creates a sub-block for multiple transactions and validates it through the consensus.
func (tm *SYN12TransactionManager) GenerateSubBlock(transactions []*SYN12Transaction) (*common.SubBlock, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate all transactions through the consensus
	for _, transaction := range transactions {
		if err := tm.ValidateTransaction(transaction); err != nil {
			return nil, fmt.Errorf("transaction validation failed: %v", err)
		}
	}

	// Create a sub-block from the validated transactions
	subBlock := common.SubBlock{
		Transactions: transactions,
		BlockID:      common.GenerateUniqueID(),
		Timestamp:    common.GetCurrentTimestamp(),
	}

	// Store the sub-block in the ledger
	if err := tm.ledgerManager.StoreSubBlock(subBlock); err != nil {
		return nil, fmt.Errorf("failed to store sub-block in the ledger: %v", err)
	}

	// Log the sub-block creation event
	if err := tm.ledgerManager.LogEvent(subBlock.BlockID, common.EventSubBlockCreated); err != nil {
		return nil, fmt.Errorf("failed to log sub-block event: %v", err)
	}

	fmt.Printf("Sub-block with ID: %s created successfully.\n", subBlock.BlockID)
	return &subBlock, nil
}

// generateTransactionHash generates a unique hash for the transaction.
func (tm *SYN12TransactionManager) generateTransactionHash(transaction *SYN12Transaction) string {
	transactionData := fmt.Sprintf("%s:%s:%s:%d:%d", transaction.TokenID, transaction.FromAddress, transaction.ToAddress, transaction.Amount, transaction.Timestamp)
	return common.HashString(transactionData)
}

// encryptTransaction encrypts a SYN12 transaction using the encryption service.
func (tm *SYN12TransactionManager) encryptTransaction(transaction *SYN12Transaction) ([]byte, error) {
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %v", err)
	}

	encryptedData, err := tm.encryptionService.Encrypt(transactionJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt transaction: %v", err)
	}

	return encryptedData, nil
}

// decryptTransaction decrypts an encrypted SYN12 transaction.
func (tm *SYN12TransactionManager) decryptTransaction(encryptedData []byte) (*SYN12Transaction, error) {
	decryptedData, err := tm.encryptionService.Decrypt(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt transaction: %v", err)
	}

	var transaction SYN12Transaction
	if err := json.Unmarshal(decryptedData, &transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	return &transaction, nil
}
