package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"synnergy_network/pkg/ledger"
	"time"
)

// Transaction represents a blockchain transaction with all the required fields.
type Transaction struct {
    TransactionID   string    // Unique identifier for the transaction
    FromAddress     string    // Sender's address
    ToAddress       string    // Receiver's address
    Amount          float64   // Amount being transferred
    Fee             float64   // Transaction fee
    TokenStandard   string    // Token standard (e.g., ERC20, Syn700)
    TokenID         string    // Unique Token ID for token transactions
    Timestamp       time.Time // Timestamp when the transaction was created
    SubBlockID      string    // Associated sub-block ID
    BlockID         string    // Associated block ID
    ValidatorID     string    // Validator who validated this transaction
    Signature       string    // Transaction signature
    Status          string    // Transaction status (e.g., pending, confirmed, failed)
    EncryptedData   string    // Encrypted transaction data
    DecryptedData   string    // Decrypted transaction data (used internally)
    ExecutionResult string    // Result after executing the transaction (e.g., success, failure reason)
	FrozenAmount float64 // Amount that is frozen in the transaction (if applicable)
    RefundAmount float64 // Amount refunded in case of a reversal or error
	ReversalRequested bool    // Whether a reversal has been requested (Add this field)
}

type CrossChainTransaction struct {
    TransactionID  string    // Unique transaction ID
    FromChain      string    // Originating chain
    ToChain        string    // Destination chain
    Amount         float64   // Amount being transferred
    TokenSymbol    string    // Token symbol being used
    FromAddress    string    // Sender's address
    ToAddress      string    // Recipient's address
    Timestamp      time.Time // Timestamp of the transaction
    ValidationHash string    // Validation hash for security
    Status         string    // Transaction status (pending, completed, failed)
    Data           string    // Additional data for the transaction (e.g., payload or metadata)
}

// TransactionRecord keeps track of all transactions in the ledger.
type TransactionRecord struct {
    From         string    // Sender's address (for financial transactions)
    To           string    // Recipient's address (for financial transactions)
    Amount       float64   // Amount transferred (if applicable)
    Fee          float64   // Transaction fee (if applicable)
    Hash         string    // Transaction hash (unique ID)
    Status       string    // Status of the transaction (e.g., "pending", "confirmed")
    BlockIndex   int       // Block in which the transaction was confirmed
    Timestamp    time.Time // Timestamp when the transaction was created or confirmed
    BlockHeight  int       // Height of the block containing the transaction
    ValidatorID  string    // ID of the validator who processed the transaction
    ID           string    // ID of the transaction or shard-related activity
    Action       string    // Action performed, e.g., "ShardCreated", "ShardUpdated", etc.
    Delegator    string    // Delegator involved in shard delegation (if applicable)
    NodeID       string    // Node ID for shard reallocation (if applicable)
    Orchestrator string    // Orchestrator for orchestrated transactions (if applicable)
    Details      string    // Additional details (if applicable)
}


// TransactionManager handles transaction creation, validation, encryption, and ledger integration.
type TransactionManager struct {
    Ledger      *ledger.Ledger              // Reference to the blockchain ledger
    Consensus   *SynnergyConsensus  // Consensus engine for Synnergy Consensus
    Encryption  *Encryption      // Encryption service
}

// NewTransactionManager initializes a new transaction manager.
func NewTransactionManager(ledgerInstance *ledger.Ledger, consensus *SynnergyConsensus, encryptionService *Encryption) *TransactionManager {
	return &TransactionManager{
		Ledger:     ledgerInstance,
		Consensus:  consensus,
		Encryption: encryptionService,
	}
}

// CreateTransaction initializes a new transaction and encrypts the data using the encryption instance.
func (tm *TransactionManager) CreateTransaction(fromAddress, toAddress string, amount, fee float64) (*Transaction, error) {
	// Generate a transaction ID
	txID := GenerateTransactionID()
	timestamp := time.Now()

	// Create the transaction struct
	transaction := &Transaction{
		TransactionID: txID,
		FromAddress:   fromAddress,
		ToAddress:     toAddress,
		Amount:        amount,
		Fee:           fee,
		Timestamp:     timestamp,
		Status:        "pending",
	}

	// Convert transaction details to bytes for encryption
	transactionBytes := []byte(fmt.Sprintf("%v", transaction))

	// Create a new encryption instance
	encryptionInstance := &Encryption{}

	// Define the encryption key for AES encryption (as in your encryption package)
	encryptionKey := []byte("your-32-byte-key-for-aes-encryption") // Example 32-byte key

	// Use the encryption instance to encrypt the transaction data (algorithm, data, and key)
	encryptedData, err := encryptionInstance.EncryptData("AES", transactionBytes, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting transaction: %v", err)
	}

	// Store the encrypted data as a string
	transaction.EncryptedData = string(encryptedData)

	// Add the transaction to the ledger: string (fromAddress), string (toAddress), float64 (amount)
	if err := tm.Ledger.BlockchainConsensusCoinLedger.AddTransaction(transaction.FromAddress, transaction.ToAddress, transaction.Amount); err != nil {
		return nil, fmt.Errorf("error adding transaction to ledger: %v", err)
	}

	return transaction, nil
}

// GenerateTransactionID generates a unique transaction ID using a hash of the current time and random data.
func GenerateTransactionID() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	data := append([]byte(fmt.Sprintf("%d", timestamp)), randomBytes...)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}


// ValidateTransaction performs validation on a transaction.
func (tm *TransactionManager) ValidateTransaction(transaction *Transaction) error {
    // Ensure that the sender has enough balance
    senderBalance, err := tm.Ledger.GetBalance(transaction.FromAddress)
    if err != nil {
        return err
    }
    if senderBalance < transaction.Amount+transaction.Fee {
        return errors.New("insufficient funds for transaction")
    }

    // Process the transaction with consensus (no return value)
    tm.Consensus.ProcessTransactions([]Transaction{*transaction}) // Pass the transaction in a slice

    // Mark transaction as validated (Assume validation logic is internal to ProcessTransactions)
    transaction.Status = "validated"

    return nil
}


// ExecuteTransaction executes the transaction and records the result in the ledger.
func (tm *TransactionManager) ExecuteTransaction(transaction *Transaction) error {
    // Define the encryption key
    encryptionKey := []byte("your-32-byte-key-for-aes-encryption") // Replace with your actual encryption key

    // Convert the EncryptedData from string to []byte for decryption
    encryptedDataBytes := []byte(transaction.EncryptedData)

    // Create the encryption instance
    encryptionInstance := &Encryption{}

    // Decrypt transaction data for execution
    decryptedData, err := encryptionInstance.DecryptData(encryptedDataBytes, encryptionKey)
    if err != nil {
        return fmt.Errorf("error decrypting transaction: %v", err)
    }

    // Convert the decrypted data back to a string for storage
    transaction.DecryptedData = string(decryptedData)

    // Ensure transaction is valid before execution
    if transaction.Status != "validated" {
        return errors.New("transaction not validated for execution")
    }

    // Simulate execution (this would be more complex in real-world systems, handling smart contracts, etc.)
    transaction.ExecutionResult = "success"
    transaction.Status = "executed"

    // Deduct from sender's balance and add to receiver's balance
    if err := tm.Ledger.BlockchainConsensusCoinLedger.UpdateBalances(transaction.TransactionID); err != nil {
        transaction.ExecutionResult = fmt.Sprintf("execution failed: %v", err)
        transaction.Status = "failed"
        return err
    }

    // Record the transaction execution in the ledger with only the transaction ID
    if err := tm.Ledger.BlockchainConsensusCoinLedger.RecordTransactionExecution(transaction.TransactionID); err != nil {
        return fmt.Errorf("error logging transaction execution: %v", err)
    }

    fmt.Printf("Transaction %s executed successfully.\n", transaction.TransactionID)
    return nil
}




// AssignTransactionToSubBlock assigns the transaction to a sub-block within the Synnergy Consensus system.
func (tm *TransactionManager) AssignTransactionToSubBlock(transaction *Transaction) error {
    // Extract transaction details and process it into a sub-block.
    transactions := []Transaction{*transaction}
    pohProof := tm.Consensus.PoH.GeneratePoHProof() // Generate the PoH proof as needed.
    
    // Create the sub-block with transactions and the PoH hash.
    subBlock := tm.Consensus.createSubBlock(transactions, pohProof.Hash)
    
    // Convert subBlock to ledger.SubBlock and add it to the ledger.
    ledgerSubBlock := ConvertSubBlockToLedgerSubBlock(subBlock)
    err := tm.Consensus.LedgerInstance.BlockchainConsensusCoinLedger.AddSubBlock(ledgerSubBlock)
    if err != nil {
        return fmt.Errorf("failed to add sub-block to ledger: %v", err)
    }
    
    // Assign sub-block ID to the transaction.
    transaction.SubBlockID = fmt.Sprintf("subblock-%d", subBlock.Index)
    
    return nil
}


// ConvertSubBlockToLedgerSubBlock converts a blockchain.SubBlock to a ledger.SubBlock.
func ConvertSubBlockToLedgerSubBlock(subBlock SubBlock) ledger.SubBlock {
    return ledger.SubBlock{
        Index:        subBlock.Index,
        Timestamp:    subBlock.Timestamp, // Keep as time.Time
        Transactions: convertTransactionsToLedger(subBlock.Transactions), // Convert transactions properly
        Validator:    subBlock.Validator,
        PrevHash:     subBlock.PrevHash,
        Hash:         subBlock.Hash,
    }
}

// ConvertTransactionsToLedger converts []Transaction to []ledger.Transaction.
func convertTransactionsToLedger(transactions []Transaction) []ledger.Transaction {
    ledgerTransactions := make([]ledger.Transaction, len(transactions))
    for i, tx := range transactions {
        ledgerTransactions[i] = ledger.Transaction{
            TransactionID: tx.TransactionID,
            FromAddress:   tx.FromAddress,
            ToAddress:     tx.ToAddress,
            Amount:        tx.Amount,
            Fee:           tx.Fee,
            Status:        tx.Status,
        }
    }
    return ledgerTransactions
}


// ConfirmTransactionInBlock confirms the transaction after sub-block validation.
func (tm *TransactionManager) ConfirmTransactionInBlock(transaction *Transaction) error {
    // Retrieve the sub-block from the ledger using its ID.
    ledgerSubBlock, err := tm.Consensus.LedgerInstance.BlockchainConsensusCoinLedger.GetSubBlockByID(transaction.SubBlockID)
    if err != nil {
        return fmt.Errorf("error retrieving sub-block: %v", err)
    }

    // Convert the ledger sub-block to the blockchain sub-block.
    subBlock := ConvertLedgerToBlockchainSubBlock(ledgerSubBlock)

    // If sub-block is validated, finalize the block.
    if tm.Consensus.PoS.ValidateSubBlock(subBlock) {
        tm.Consensus.FinalizeBlock() // Finalize the block when the sub-block is valid.
        transaction.BlockID = fmt.Sprintf("block-%d", tm.Consensus.LedgerInstance.BlockchainConsensusCoinLedger.GetBlockCount())
        transaction.Status = "confirmed"
    } else {
        return fmt.Errorf("sub-block %s validation failed", transaction.SubBlockID)
    }
    
    return nil
}

// ConvertLedgerToBlockchainSubBlock converts a ledger.SubBlock to a blockchain.SubBlock.
func ConvertLedgerToBlockchainSubBlock(ledgerSubBlock ledger.SubBlock) SubBlock {
    return SubBlock{
        Index:        ledgerSubBlock.Index,
        Timestamp:    ledgerSubBlock.Timestamp, // Keep it as time.Time
        Transactions: convertLedgerTransactionsToBlockchain(ledgerSubBlock.Transactions), // Convert transactions properly
        Validator:    ledgerSubBlock.Validator,
        PrevHash:     ledgerSubBlock.PrevHash,
        Hash:         ledgerSubBlock.Hash,
    }
}

// ConvertLedgerTransactionsToBlockchain converts []ledger.Transaction to []Transaction.
func convertLedgerTransactionsToBlockchain(ledgerTransactions []ledger.Transaction) []Transaction {
    transactions := make([]Transaction, len(ledgerTransactions))
    for i, ledgerTx := range ledgerTransactions {
        transactions[i] = Transaction{
            TransactionID: ledgerTx.TransactionID,
            FromAddress:   ledgerTx.FromAddress,
            ToAddress:     ledgerTx.ToAddress,
            Amount:        ledgerTx.Amount,
            Fee:           ledgerTx.Fee,
            Status:        ledgerTx.Status,
            // Add any other fields if needed
        }
    }
    return transactions
}



// GetTransactionByID retrieves a transaction by its ID (without encryption fields).
func (tm *TransactionManager) GetTransactionByID(txID string, userKey string) (*Transaction, error) {
    // Call to Ledger with the transaction ID
    transactionRecord, err := tm.Ledger.BlockchainConsensusCoinLedger.GetTransactionByID(txID)
    if err != nil {
        return nil, fmt.Errorf("transaction %s not found: %v", txID, err)
    }

    // Create a new Transaction instance and populate it with data from the transactionRecord
    transaction := &Transaction{
        TransactionID:   transactionRecord.ID,
        FromAddress:     transactionRecord.From,
        ToAddress:       transactionRecord.To,
        Amount:          transactionRecord.Amount,
        Fee:             transactionRecord.Fee,
        TokenStandard:   "", // Add this if applicable
        TokenID:         "", // Add this if applicable
        Timestamp:       transactionRecord.Timestamp,      // Ensure Timestamp is properly handled
        SubBlockID:      "", // Add this if applicable
        BlockID:         fmt.Sprintf("%d", transactionRecord.BlockIndex),
        ValidatorID:     transactionRecord.ValidatorID,
        Signature:       "", // Add this if applicable
        Status:          transactionRecord.Status,
        EncryptedData:   "", // Not applicable here
        DecryptedData:   "", // Not applicable here
        ExecutionResult: "", // Add this if applicable
        FrozenAmount:    0.0, // Add this if applicable
        RefundAmount:    0.0, // Add this if applicable
        ReversalRequested: false, // Add this if applicable
    }

    return transaction, nil
}





// CancelTransaction cancels a pending or unvalidated transaction.
func (tm *TransactionManager) CancelTransaction(txID string, userKey string) error {
    // Call GetTransactionByID with only txID
    transaction, err := tm.Ledger.BlockchainConsensusCoinLedger.GetTransactionByID(txID)
    if err != nil {
        return fmt.Errorf("transaction %s not found: %v", txID, err)
    }

    // Check if the transaction can be cancelled
    if transaction.Status == "executed" || transaction.Status == "confirmed" {
        return errors.New("cannot cancel an already executed or confirmed transaction")
    }

    // Remove the transaction from the ledger
    if err := tm.Ledger.BlockchainConsensusCoinLedger.RemoveTransaction(txID); err != nil {
        return fmt.Errorf("error cancelling transaction: %v", err)
    }

    // Update the transaction status
    transaction.Status = "cancelled"
    fmt.Printf("Transaction %s has been successfully cancelled.\n", txID)
    
    return nil
}



// GetTransactionHistory retrieves the transaction history for a specific address.
func (tm *TransactionManager) GetTransactionHistory(address string) ([]*Transaction, error) {
	// Get the encryption instance and the key (ensure the key is properly defined)
	encryptionInstance := &Encryption{}
	encryptionKey := []byte("your-32-byte-key-for-aes-encryption") // Define a 32-byte key for AES encryption

	// Get the transaction history from the ledger (returns []*ledger.Transaction)
	ledgerTransactions, err := tm.Ledger.BlockchainConsensusCoinLedger.GetTransactionHistoryByAddress(address)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transaction history: %v", err)
	}

	// Create a slice to hold the converted transactions
	var transactions []*Transaction

	// Loop over ledger transactions, decrypt and convert them
	for _, ledgerTx := range ledgerTransactions {
		decryptedData, err := encryptionInstance.DecryptData([]byte(ledgerTx.EncryptedData), encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("error decrypting transaction %s: %v", ledgerTx.TransactionID, err)
		}

		// Convert ledger.Transaction to Transaction and assign decrypted data
		tx := &Transaction{
			TransactionID:   ledgerTx.TransactionID,
			FromAddress:     ledgerTx.FromAddress,
			ToAddress:       ledgerTx.ToAddress,
			Amount:          ledgerTx.Amount,
			Fee:             ledgerTx.Fee,
			TokenStandard:   ledgerTx.TokenStandard,
			TokenID:         ledgerTx.TokenID,
			Timestamp:       ledgerTx.Timestamp,
			SubBlockID:      ledgerTx.SubBlockID,
			BlockID:         ledgerTx.BlockID,
			ValidatorID:     ledgerTx.ValidatorID,
			Signature:       ledgerTx.Signature,
			Status:          ledgerTx.Status,
			EncryptedData:   ledgerTx.EncryptedData,
			DecryptedData:   string(decryptedData), // Convert decrypted bytes to a string
			ExecutionResult: ledgerTx.ExecutionResult,
			FrozenAmount:    ledgerTx.FrozenAmount,
			RefundAmount:    ledgerTx.RefundAmount,
			ReversalRequested: ledgerTx.ReversalRequested,
		}

		// Add the converted transaction to the slice
		transactions = append(transactions, tx)
	}

	return transactions, nil
}


// StringForSigning creates a string representation of the transaction for signing purposes.
func (tx *Transaction) StringForSigning() string {
	return fmt.Sprintf("%s:%s:%.2f:%d", 
		tx.FromAddress, 
		tx.ToAddress, 
		tx.Amount, 
		tx.Timestamp.Unix()) // Customize the fields as per your signing requirements
}

// String provides a string representation of the transaction.
func (tx *Transaction) String() string {
	return fmt.Sprintf("Transaction ID: %s, From: %s, To: %s, Amount: %.2f, Status: %s", 
		tx.TransactionID, tx.FromAddress, tx.ToAddress, tx.Amount, tx.Status)
}


// SignTransaction signs a transaction by appending r and s as bytes.
func (tx *Transaction) SignTransaction(r, s *big.Int) []byte {
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

