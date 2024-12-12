package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// MaxTransactionsPerSubBlock defines the limit of transactions per sub-block
const MaxTransactionsPerSubBlock = 10000


// SubBlockChain holds a list of sub-blocks and manages them.
type SubBlockChain struct {
	SubBlocks  []SubBlock // The actual chain of sub-blocks
	Validators []string   // List of validators for PoS
	Ledger     *ledger.Ledger    // Pointer to the ledger for storing validated sub-blocks
	mutex      sync.Mutex // Mutex for thread-safe operations
}

// SubBlockManager manages the creation, validation, and addition of sub-blocks to the ledger.
type SubBlockManager struct {
	LedgerInstance *ledger.Ledger    // Reference to the ledger to store sub-blocks
	Encryption     *Encryption // Encryption instance for data encryption
	mutex          sync.Mutex             // Mutex for thread-safe operations
}

// NewSubBlockManager initializes a new SubBlockManager instance.
func NewSubBlockManager(ledgerInstance *ledger.Ledger, encryptionModule *Encryption) *SubBlockManager {
	return &SubBlockManager{
		LedgerInstance: ledgerInstance,
		Encryption:     encryptionModule,
		mutex:          sync.Mutex{},
	}
}


type SubBlock struct {
	SubBlockID   string        // Sub-block ID (new field)
    Index        int           // Sub-block index
    Timestamp    time.Time     // Sub-block creation time (keeping it as time.Time here)
    Transactions []Transaction // List of transactions in the sub-block
    Validator    string        // Validator who validated the sub-block
    PrevHash     string        // Previous sub-block's hash
    Hash         string        // Current sub-block's hash
	PoHProof     PoHProof   
	Status      string       // Block status (new field)
	Signature    string        // Signature from the validator (new field)

}

// ConvertToLedgerSubBlock converts a SubBlock to a ledger.SubBlock
func ConvertToLedgerSubBlock(subBlock SubBlock) ledger.SubBlock {
    return ledger.SubBlock{
        SubBlockID:   subBlock.SubBlockID,
        Index:        subBlock.Index,
        Timestamp:    subBlock.Timestamp,
        Transactions: ConvertTransactions(subBlock.Transactions), // Convert transactions
        Validator:    subBlock.Validator,
        PrevHash:     subBlock.PrevHash,
        Hash:         subBlock.Hash,
        PoHProof:     ConvertPoHProof(subBlock.PoHProof),         // Convert PoHProof
        Status:       subBlock.Status,
    }
}



// NewSubBlock creates a new sub-block with transactions and validator info
func (bc *SubBlockChain) NewSubBlock(transactions []Transaction, prevHash string) SubBlock {
    validator := bc.selectValidator() // Select PoS validator

    subBlock := SubBlock{
        Index:        len(bc.SubBlocks),
        Timestamp:    time.Now(),
        Transactions: transactions,
        Validator:    validator,
        PrevHash:     prevHash,
    }

    subBlock.Hash = calculateSubBlockHash(subBlock)

    // Add the sub-block to the blockchain
    bc.SubBlocks = append(bc.SubBlocks, subBlock)

    // Check if the Ledger is available and call the AddSubBlock method
    if bc.Ledger != nil {
        ledgerSubBlock := ConvertToLedgerSubBlock(subBlock) // Convert to ledger.SubBlock
        err := bc.Ledger.BlockchainConsensusCoinLedger.AddSubBlock(ledgerSubBlock)        // Use the converted value
        if err != nil {
            fmt.Printf("Error adding sub-block to ledger: %v\n", err)
        }
    }

    fmt.Printf("SubBlock %d created with validator: %s\n", subBlock.Index, validator)
    return subBlock
}



// selectValidator randomly selects a validator using a PoS-like mechanism
func (bc *SubBlockChain) selectValidator() string {
    rand.Seed(time.Now().UnixNano())
    return bc.Validators[rand.Intn(len(bc.Validators))]
}

// calculateSubBlockHash generates a SHA-256 hash for the sub-block
func calculateSubBlockHash(subBlock SubBlock) string {
    hashInput := fmt.Sprintf("%d%s%s%s", 
        subBlock.Index, 
        subBlock.Timestamp.String(),
        subBlock.PrevHash,
        subBlock.Validator)

    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// AddTransactions adds new transactions to a sub-block and finalizes it
func (bc *SubBlockChain) AddTransactions(transactions []Transaction) {
    prevHash := ""
    if len(bc.SubBlocks) > 0 {
        prevHash = bc.SubBlocks[len(bc.SubBlocks)-1].Hash
    }

    for len(transactions) > 0 {
        var txBatch []Transaction

        if len(transactions) > MaxTransactionsPerSubBlock {
            txBatch = transactions[:MaxTransactionsPerSubBlock]
            transactions = transactions[MaxTransactionsPerSubBlock:]
        } else {
            txBatch = transactions
            transactions = nil
        }

        // Create a new sub-block with a batch of transactions
        bc.NewSubBlock(txBatch, prevHash)
    }
}

// ValidateSubBlock validates the sub-block using PoH (timestamping) and PoS (validator)
func (bc *SubBlockChain) ValidateSubBlock(subBlock SubBlock) bool {
    // Simulating Proof of History (PoH): Check if the timestamp is valid
    if time.Since(subBlock.Timestamp) < 0 {
        fmt.Println("SubBlock rejected due to invalid timestamp.")
        return false
    }

    // Simulating Proof of Stake (PoS): Check if the validator is valid
    isValidValidator := false
    for _, validator := range bc.Validators {
        if subBlock.Validator == validator {
            isValidValidator = true
            break
        }
    }

    if !isValidValidator {
        fmt.Println("SubBlock rejected due to invalid validator.")
        return false
    }

    fmt.Printf("SubBlock %d successfully validated.\n", subBlock.Index)
    return true
}

// AddTransactionToSubBlock adds a transaction to the current sub-block.
func (sbm *SubBlockManager) AddTransactionToSubBlock(transaction *Transaction) error {
	sbm.mutex.Lock()
	defer sbm.mutex.Unlock()

	// Retrieve the current sub-block or create a new one
	subBlock, err := sbm.GetOrCreateCurrentSubBlock()
	if err != nil {
		return fmt.Errorf("failed to get or create current sub-block: %v", err)
	}

	// Add the transaction to the sub-block
	subBlock.Transactions = append(subBlock.Transactions, *transaction)

	// Check if the sub-block has reached the maximum number of transactions
	if len(subBlock.Transactions) >= MaxTransactionsPerSubBlock {
		// Finalize and log the sub-block
		err := sbm.FinalizeSubBlock(subBlock)
		if err != nil {
			return fmt.Errorf("failed to finalize sub-block: %v", err)
		}
		// Create a new sub-block for further transactions
		sbm.CreateNewSubBlock()
	}

	return nil
}


// GetOrCreateCurrentSubBlock retrieves the current sub-block or creates a new one if none exists.
func (sbm *SubBlockManager) GetOrCreateCurrentSubBlock() (*SubBlock, error) {
	// Use the ledger to get the last sub-block
	ledgerSubBlock, err := sbm.LedgerInstance.BlockchainConsensusCoinLedger.GetLastSubBlock()
	if err != nil || ledgerSubBlock == nil {
		// If no sub-block is available, create a new one
		return sbm.CreateNewSubBlock(), nil
	}

	// Convert ledger.SubBlock to common.SubBlock
	commonSubBlock := ConvertLedgerSubBlockToCommonSubBlock(ledgerSubBlock)

	// Return the current sub-block in the common format
	return commonSubBlock, nil
}


// CreateNewSubBlock creates and returns a new sub-block.
func (sbm *SubBlockManager) CreateNewSubBlock() *SubBlock {
    newSubBlock := &SubBlock{
        Index:        sbm.LedgerInstance.BlockchainConsensusCoinLedger.GetSubBlockCount(), // Sub-block index based on the ledger's sub-block count
        Timestamp:    time.Now(),
        Transactions: []Transaction{},
        Validator:    "", // Validator will be assigned later during PoS validation
        PrevHash:     sbm.getLastSubBlockHash(),
    }

    // Convert *SubBlock to ledger.SubBlock by dereferencing newSubBlock
    ledgerSubBlock := ConvertToLedgerSubBlock(*newSubBlock) // Dereference the pointer

    // Add the new sub-block to the ledger
    err := sbm.LedgerInstance.BlockchainConsensusCoinLedger.AddSubBlock(ledgerSubBlock)
    if err != nil {
        fmt.Printf("Error adding sub-block to ledger: %v\n", err)
    }

    return newSubBlock
}



// FinalizeSubBlock finalizes and logs a sub-block in the ledger.
func (sbm *SubBlockManager) FinalizeSubBlock(subBlock *SubBlock) error {
	// Finalize the sub-block: Assign a validator, calculate the hash, etc.
	subBlock.Hash = calculateSubBlockHash(*subBlock)

	// Convert common SubBlock to ledger SubBlock
	ledgerSubBlock := ConvertCommonSubBlockToLedgerSubBlock(subBlock)

	// Log the finalized sub-block in the ledger
	err := sbm.LedgerInstance.BlockchainConsensusCoinLedger.LogSubBlock(ledgerSubBlock)
	if err != nil {
		return fmt.Errorf("failed to log sub-block: %v", err)
	}

	return nil
}


// getLastSubBlockHash retrieves the hash of the last sub-block in the ledger.
func (sbm *SubBlockManager) getLastSubBlockHash() string {
	subBlock, err := sbm.LedgerInstance.BlockchainConsensusCoinLedger.GetLastSubBlock()
	if err == nil && subBlock != nil {
		return subBlock.Hash
	}
	return ""
}

// ConvertCommonSubBlockToLedgerSubBlock converts a common SubBlock to a ledger SubBlock.
func ConvertCommonSubBlockToLedgerSubBlock(subBlock *SubBlock) *ledger.SubBlock {
	return &ledger.SubBlock{
		Index:        subBlock.Index,
		Timestamp:    subBlock.Timestamp, // Directly assign the time.Time value
		Transactions: convertTransactionsToLedgerFormat(subBlock.Transactions), // Fix transaction conversion
		Validator:    subBlock.Validator,
		PrevHash:     subBlock.PrevHash,
		Hash:         subBlock.Hash,
	}
}

// ConvertLedgerSubBlockToCommonSubBlock converts a ledger SubBlock back to a common SubBlock.
func ConvertLedgerSubBlockToCommonSubBlock(ledgerSubBlock *ledger.SubBlock) *SubBlock {
	return &SubBlock{
		Index:        ledgerSubBlock.Index,
		Timestamp:    ledgerSubBlock.Timestamp, // Directly assign the time.Time value
		Transactions: convertLedgerTransactionsToCommonFormat(ledgerSubBlock.Transactions), // Fix transaction conversion
		Validator:    ledgerSubBlock.Validator,
		PrevHash:     ledgerSubBlock.PrevHash,
		Hash:         ledgerSubBlock.Hash,
	}
}

// GetCurrentSubBlock retrieves the current sub-block.
func (sbm *SubBlockManager) GetCurrentSubBlock() (*SubBlock, error) {
	// Use the ledger to get the last sub-block
	ledgerSubBlock, err := sbm.LedgerInstance.BlockchainConsensusCoinLedger.GetLastSubBlock()
	if err != nil || ledgerSubBlock == nil {
		return nil, fmt.Errorf("no current sub-block found")
	}

	// Convert ledger.SubBlock to common.SubBlock (assuming you have a conversion function)
	return ConvertLedgerSubBlockToCommonSubBlock(ledgerSubBlock), nil
}

// Helper function to convert []Transaction to []ledger.Transaction
func convertTransactionsToLedgerFormat(transactions []Transaction) []ledger.Transaction {
	ledgerTransactions := make([]ledger.Transaction, len(transactions))
	for i, tx := range transactions {
		ledgerTransactions[i] = ledger.Transaction{
			TransactionID:   tx.TransactionID,
			FromAddress:     tx.FromAddress,
			ToAddress:       tx.ToAddress,
			Amount:          tx.Amount,
			Fee:             tx.Fee,
			TokenStandard:   tx.TokenStandard,
			TokenID:         tx.TokenID,
			Timestamp:       tx.Timestamp, // Use time.Time directly
			SubBlockID:      tx.SubBlockID,
			BlockID:         tx.BlockID,
			ValidatorID:     tx.ValidatorID,
			Signature:       tx.Signature,
			Status:          tx.Status,
			EncryptedData:   tx.EncryptedData,
			DecryptedData:   tx.DecryptedData,
			ExecutionResult: tx.ExecutionResult,
			FrozenAmount:    tx.FrozenAmount,
			RefundAmount:    tx.RefundAmount,
			ReversalRequested: tx.ReversalRequested,
		}
	}
	return ledgerTransactions
}

// Helper function to convert []ledger.Transaction to []Transaction
func convertLedgerTransactionsToCommonFormat(ledgerTransactions []ledger.Transaction) []Transaction {
	transactions := make([]Transaction, len(ledgerTransactions))
	for i, ledgerTx := range ledgerTransactions {
		transactions[i] = Transaction{
			TransactionID:   ledgerTx.TransactionID,
			FromAddress:     ledgerTx.FromAddress,
			ToAddress:       ledgerTx.ToAddress,
			Amount:          ledgerTx.Amount,
			Fee:             ledgerTx.Fee,
			TokenStandard:   ledgerTx.TokenStandard,
			TokenID:         ledgerTx.TokenID,
			Timestamp:       ledgerTx.Timestamp, // Use time.Time directly
			SubBlockID:      ledgerTx.SubBlockID,
			BlockID:         ledgerTx.BlockID,
			ValidatorID:     ledgerTx.ValidatorID,
			Signature:       ledgerTx.Signature,
			Status:          ledgerTx.Status,
			EncryptedData:   ledgerTx.EncryptedData,
			DecryptedData:   ledgerTx.DecryptedData,
			ExecutionResult: ledgerTx.ExecutionResult,
			FrozenAmount:    ledgerTx.FrozenAmount,
			RefundAmount:    ledgerTx.RefundAmount,
			ReversalRequested: ledgerTx.ReversalRequested,
		}
	}
	return transactions
}
