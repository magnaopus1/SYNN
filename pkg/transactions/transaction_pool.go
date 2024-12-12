package transactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewTransactionPool initializes a new transaction pool.
func NewTransactionPool(maxPoolSize int, ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *TransactionPool {
	return &TransactionPool{
		transactions:      make(map[string]*common.Transaction),
		pendingSubBlocks:  make(map[string][]*common.Transaction),
		maxPoolSize:       maxPoolSize,
		ledger:            ledgerInstance,
		encryptionService: encryptionService,  // Now passing a pointer
	}
}

// AddTransaction adds a new transaction to the pool.
func (tp *TransactionPool) AddTransaction(tx *common.Transaction) error {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	// Check if the transaction pool is full
	if len(tp.transactions) >= tp.maxPoolSize {
		return errors.New("transaction pool is full")
	}

	// Serialize the transaction before encryption (assuming JSON encoding)
	serializedTx, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to serialize transaction: %v", err)
	}

	// Encrypt the serialized transaction for security
	encryptedData, err := tp.encryptionService.EncryptData("AES", serializedTx, common.EncryptionKey) // Pass "AES" as the encryption algorithm
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}

	// Ensure transaction isn't already in the pool
	if _, exists := tp.transactions[tx.TransactionID]; exists {
		return errors.New("transaction already exists in pool")
	}

	// Store the encrypted data in the transaction's EncryptedData field
	tx.EncryptedData = string(encryptedData)

	// Add the transaction to the pool
	tp.transactions[tx.TransactionID] = tx

	return nil
}




// GetTransaction retrieves a transaction from the pool by ID.
func (tp *TransactionPool) GetTransaction(txID string) (*common.Transaction, error) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	tx, exists := tp.transactions[txID]
	if !exists {
		return nil, errors.New("transaction not found in pool")
	}

	// Decrypt the EncryptedData field of the transaction
	decryptedData, err := tp.encryptionService.DecryptData([]byte(tx.EncryptedData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt transaction: %v", err)
	}

	// Unmarshal the decrypted data back into a common.Transaction struct
	var decryptedTx common.Transaction
	err = json.Unmarshal(decryptedData, &decryptedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal decrypted transaction: %v", err)
	}

	return &decryptedTx, nil
}

// RemoveTransaction removes a transaction from the pool.
func (tp *TransactionPool) RemoveTransaction(txID string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	delete(tp.transactions, txID)
}

// CreateSubBlock creates a new sub-block from a set of pending transactions.
func (tp *TransactionPool) CreateSubBlock(subBlockID string, txCount int) (*common.SubBlock, error) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	if len(tp.transactions) == 0 {
		return nil, errors.New("no transactions in pool")
	}

	subBlock := &common.SubBlock{
		SubBlockID:   subBlockID,  // Fixed field name
		Timestamp:    time.Now(),
		Transactions: []common.Transaction{},  // Fixed type for transactions
	}

	// Select transactions for the sub-block
	count := 0
	for txID, tx := range tp.transactions {
		if count >= txCount {
			break
		}

		// Append the transaction by dereferencing the pointer
		subBlock.Transactions = append(subBlock.Transactions, *tx)
		tp.pendingSubBlocks[subBlockID] = append(tp.pendingSubBlocks[subBlockID], tx)
		delete(tp.transactions, txID) // Remove from pool
		count++
	}

	return subBlock, nil
}


// AddSubBlockToLedger validates the sub-block and adds it to the ledger.
func (tp *TransactionPool) AddSubBlockToLedger(subBlock *common.SubBlock) error {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	// Validate transactions in the sub-block before adding to the ledger
	for _, tx := range subBlock.Transactions {
		err := tp.ledger.ValidateTransaction(tx.TransactionID) // Pass TransactionID for validation
		if err != nil {
			return fmt.Errorf("transaction %s validation failed: %v", tx.TransactionID, err)
		}
	}

	// Convert common.Transaction to ledger.Transaction
	ledgerTransactions := make([]ledger.Transaction, len(subBlock.Transactions))
	for i, tx := range subBlock.Transactions {
		ledgerTransactions[i] = ledger.Transaction{
			TransactionID:   tx.TransactionID,
			FromAddress:     tx.FromAddress,   // Correct field name for sender
			ToAddress:       tx.ToAddress,     // Correct field name for receiver
			Amount:          tx.Amount,
			Fee:             tx.Fee,
			SubBlockID:      tx.SubBlockID,    // Correct field for sub-block association
			Status:          tx.Status,
			Timestamp:       tx.Timestamp,
			ValidatorID:     tx.ValidatorID,
			Signature:       tx.Signature,
			ExecutionResult: tx.ExecutionResult,
		}
	}

	// Convert *common.SubBlock to ledger.SubBlock, excluding PoHProof for now if the fields are mismatched
	ledgerSubBlock := ledger.SubBlock{
		SubBlockID:    subBlock.SubBlockID,
		Index:         subBlock.Index,
		Timestamp:     subBlock.Timestamp,
		Transactions:  ledgerTransactions, // Use the converted transactions
		Validator:     subBlock.Validator,
		PrevHash:      subBlock.PrevHash,
		Hash:          subBlock.Hash,
		Status:        subBlock.Status,
		Signature:     subBlock.Signature,
	}

	// Add sub-block to the ledger
	err := tp.ledger.AddSubBlock(ledgerSubBlock) // Pass the converted ledgerSubBlock
	if err != nil {
		return fmt.Errorf("failed to add sub-block to ledger: %v", err)
	}

	// Remove sub-block from the pool
	delete(tp.pendingSubBlocks, subBlock.SubBlockID)

	return nil
}



// PoolSize returns the current number of transactions in the pool.
func (tp *TransactionPool) PoolSize() int {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	return len(tp.transactions)
}

// ClearPool clears all transactions in the pool (used for system resets or testing).
func (tp *TransactionPool) ClearPool() {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	tp.transactions = make(map[string]*common.Transaction)
}

// ListTransactions returns a list of all transactions in the pool.
func (tp *TransactionPool) ListTransactions() []*common.Transaction {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	txs := []*common.Transaction{}
	for _, tx := range tp.transactions {
		// Decrypt the EncryptedData field (which should be a string of the encrypted data)
		decryptedData, err := tp.encryptionService.DecryptData([]byte(tx.EncryptedData), common.EncryptionKey)
		if err == nil {
			// Unmarshal the decrypted data back into a Transaction struct
			var decryptedTx common.Transaction
			err := json.Unmarshal(decryptedData, &decryptedTx)
			if err == nil {
				txs = append(txs, &decryptedTx)
			}
		}
	}
	return txs
}


// ListPendingSubBlocks returns a list of all pending sub-blocks in the pool.
func (tp *TransactionPool) ListPendingSubBlocks() []*common.SubBlock {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	subBlocks := []*common.SubBlock{}
	for subBlockID, txs := range tp.pendingSubBlocks {
		// Convert []*common.Transaction to []common.Transaction
		transactions := make([]common.Transaction, len(txs))
		for i, tx := range txs {
			transactions[i] = *tx // Dereference the pointer to get the value
		}

		// Create the sub-block with the correct SubBlockID and Transactions
		subBlock := &common.SubBlock{
			SubBlockID:   subBlockID,      // Use SubBlockID instead of ID
			Timestamp:    time.Now(),
			Transactions: transactions,    // Use the converted transactions
		}
		subBlocks = append(subBlocks, subBlock)
	}
	return subBlocks
}

