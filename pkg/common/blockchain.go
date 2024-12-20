package common

import (
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// Blockchain represents the main blockchain with all blocks
type Blockchain struct {
	Chain               []Block        // The blockchain itself
	PendingTransactions []Transaction  // Transactions waiting to be included in a block
	SubBlockChain       SubBlockChain  // Sub-blockchain that handles sub-blocks
	Validators          []string       // List of validators
	OwnerWallet         string         // The owner's wallet address
	mutex               sync.Mutex     // Mutex for thread-safe operations
	Ledger              *ledger.Ledger // Ledger to store blocks and transactions
}

// ConvertTransactions converts []Blockchain.Transaction to []Ledger.Transaction.
func ConvertTransactions(bcTransactions []Transaction) []ledger.Transaction {
	var ledgerTransactions []ledger.Transaction
	for _, tx := range bcTransactions {
		ledgerTransactions = append(ledgerTransactions, ledger.Transaction{
			TransactionID:  tx.TransactionID,    // Map TransactionID correctly
			FromAddress:    tx.FromAddress,      // Ensure names match
			ToAddress:      tx.ToAddress,        // Same for the receiver
			Amount:         tx.Amount,           // Amount field should match
			Timestamp:      tx.Timestamp,        // Timestamp matches
			Signature:      tx.Signature,        // Signature mapping
			Status:         tx.Status,           // Status matches
			TokenID:        tx.TokenID,          // Mapping Token ID correctly
			ExecutionResult: tx.ExecutionResult, // Map execution result
			Fee:            tx.Fee,              // Transaction Fee
			FrozenAmount:   tx.FrozenAmount,     // Frozen amount
			RefundAmount:   tx.RefundAmount,     // Refund amount
			ReversalRequested: tx.ReversalRequested, // Reversal requested status
		})
	}
	return ledgerTransactions
}


// ConvertSubBlocks converts []Blockchain.SubBlock to []Ledger.SubBlock.
func ConvertSubBlocks(bcSubBlocks []SubBlock) []ledger.SubBlock {
	var ledgerSubBlocks []ledger.SubBlock
	for _, sb := range bcSubBlocks {
		ledgerSubBlocks = append(ledgerSubBlocks, ledger.SubBlock{
			SubBlockID:   sb.SubBlockID,
			Transactions: ConvertTransactions(sb.Transactions), // Convert transactions
			Timestamp:    sb.Timestamp,         // Timestamps match
			Validator:    sb.Validator,         // Validator field matches
			PrevHash:     sb.PrevHash,          // Previous hash mapping
			Hash:         sb.Hash,              // Current hash mapping
			PoHProof:     ConvertPoHProof(sb.PoHProof),  // Convert PoHProof
			Status:       sb.Status,            // Status field matches
		})
	}
	return ledgerSubBlocks
}




// ConvertToLedgerBlock converts a Blockchain.Block to a Ledger.Block.
func ConvertToLedgerBlock(bcBlock Block) ledger.Block {
	return ledger.Block{
		BlockID:     bcBlock.BlockID,
		Index:       bcBlock.Index,
		Timestamp:   bcBlock.Timestamp,
		SubBlocks:   ConvertSubBlocks(bcBlock.SubBlocks), // Convert SubBlocks from Blockchain to Ledger format
		PrevHash:    bcBlock.PrevHash,
		Hash:        bcBlock.Hash,
		Nonce:       bcBlock.Nonce,
		Difficulty:  bcBlock.Difficulty,
		MinerReward: bcBlock.MinerReward,
		Validators:  bcBlock.Validators,
		Status:      bcBlock.Status, // Ensure that both types have this field or remove it if unnecessary
	}
}

// ConvertPoHProof converts Blockchain.PoHProof to Ledger.PoHProof.
func ConvertPoHProof(bcPoHProof PoHProof) ledger.PoHProof {
	return ledger.PoHProof{
		Sequence:  bcPoHProof.Sequence,   // Map the sequence number
		Timestamp: bcPoHProof.Timestamp,  // Map the timestamp
		Hash:      bcPoHProof.Hash,       // Map the hash generated by PoH
	}
}



// AddBlock adds a new block to the main blockchain.
func (bc *Blockchain) AddBlock(newBlock Block) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Compute the hash of the previous block
	if len(bc.Chain) > 0 {
		newBlock.PrevHash = bc.Chain[len(bc.Chain)-1].Hash
	}

	// Add the new block to the chain
	bc.Chain = append(bc.Chain, newBlock)

	// Convert the Blockchain.Block to a Ledger.Block before adding it to the ledger
	if bc.Ledger != nil {
		ledgerBlock := ConvertToLedgerBlock(newBlock) // Convert the block
		err := bc.Ledger.BlockchainConsensusCoinLedger.AddBlock(ledgerBlock)        // Pass the converted block to the ledger
		if err != nil {
			fmt.Printf("Error adding block to ledger: %v\n", err)
		}
	}
}

// AddTransaction adds a new transaction to the list of pending transactions.
func (bc *Blockchain) AddTransaction(tx Transaction) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Add the transaction to the list of pending transactions
	bc.PendingTransactions = append(bc.PendingTransactions, tx)
}

// MineBlock mines a new block and includes the pending transactions.
func (bc *Blockchain) MineBlock(minerAddress string) Block {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Create a new block with the pending transactions
	newBlock := Block{
		BlockID:    generateBlockID(), // Custom function to generate a unique block ID
		Timestamp:  time.Now(),
		SubBlocks:  []SubBlock{},
		Nonce:      0, // Start with a nonce of 0
		Difficulty: 1, // Example difficulty level
		MinerReward: 50, // Example reward for the miner
	}

	// Add transactions from the pending transactions list to the new block
	newBlock.SubBlocks = bc.SubBlockChain.MineSubBlocks(bc.PendingTransactions)

	// Clear pending transactions after they've been included in the block
	bc.PendingTransactions = []Transaction{}

	// Add the new block to the blockchain
	bc.AddBlock(newBlock)

	return newBlock
}



// MineSubBlocks mines sub-blocks for a block using PoS or PoW.
func (sbc *SubBlockChain) MineSubBlocks(transactions []Transaction) []SubBlock {
	sbc.mutex.Lock()
	defer sbc.mutex.Unlock()

	var minedSubBlocks []SubBlock

	// Mining logic for sub-blocks (can be PoS or PoW)
	// Example: Group transactions into sub-blocks
	for len(transactions) > 0 {
		subBlockTxs := transactions[:10] // Example: each sub-block has 10 transactions
		transactions = transactions[10:]

		subBlock := SubBlock{
			SubBlockID:  generateSubBlockID(),
			Transactions: subBlockTxs,
			Timestamp:   time.Now(),
			Validator:   sbc.Validators[0], // Example: first validator handles this sub-block
		}

		// Add sub-block to list of mined sub-blocks
		minedSubBlocks = append(minedSubBlocks, subBlock)
	}

	return minedSubBlocks
}


// Helper function to generate unique block IDs.
func generateBlockID() string {
	return "block_" + time.Now().Format("20060102150405")
}

// Helper function to generate unique sub-block IDs.
func generateSubBlockID() string {
	return "subblock_" + time.Now().Format("20060102150405")
}

// Ledger represents a general ledger to store blockchain and sub-blockchain data.
type Ledger struct {
	mutex       sync.Mutex
	Blocks      []Block
	SubBlocks   []SubBlock
	Transactions []Transaction
}

// StoreBlock stores a block in the ledger.
func (l *Ledger) StoreBlock(block Block) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.Blocks = append(l.Blocks, block)
}

// StoreSubBlock stores a sub-block in the ledger.
func (l *Ledger) StoreSubBlock(subBlock SubBlock) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.SubBlocks = append(l.SubBlocks, subBlock)
}

// StoreTransaction stores a transaction in the ledger.
func (l *Ledger) StoreTransaction(tx Transaction) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.Transactions = append(l.Transactions, tx)
}

