package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Difficulty defines the number of leading zeroes required in a valid block hash
const Difficulty = 1

// Block represents a blockchain block.
type Block struct {
    BlockID     string       // Unique identifier for the block
    Index       int          // Block index
    Timestamp   time.Time    // Block creation time
    SubBlocks   []SubBlock   // Sub-blocks inside the block
    PrevHash    string       // Previous block's hash
    Hash        string       // Current block's hash
    Nonce       int          // Nonce for PoW
    Difficulty  int          // Difficulty level for PoW
    MinerReward float64      // Reward given to the miner
    Validators  []string     // List of validators who contributed to the block validation
	Status      string       // Block status (new field)
}



// NewBlock creates a new block with 1000 sub-blocks and mines it using PoW
func (bc *Blockchain) NewBlock(subBlocks []SubBlock, prevHash string, pow *PoW) (Block, error) {
    // Create the new block
    block := Block{
        Index:      len(bc.Chain),
        Timestamp:  time.Now(),
        SubBlocks:  subBlocks,
        PrevHash:   prevHash,
        Nonce:      0,         // Start with 0, will be updated during PoW mining
        Difficulty: pow.State.Difficulty,
    }

    // Call MineBlock to mine the new block
    err := pow.MineBlock(&block)
    if err != nil {
        return block, fmt.Errorf("failed to mine block: %v", err)
    }

    // Add the mined block to the blockchain
    bc.Chain = append(bc.Chain, block)

    // Convert Blockchain.Block to ledger.Block before adding to the ledger
    ledgerBlock := ConvertToLedgerBlock(block)

    // Add the block to the ledger for persistence
    if err := bc.Ledger.BlockchainConsensusCoinLedger.AddBlock(ledgerBlock); err != nil {
        return block, fmt.Errorf("failed to add block to ledger: %v", err)
    }

    fmt.Printf("Block %d created with hash: %s\n", block.Index, block.Hash)
    return block, nil
}


// mineBlock performs Proof of Work (PoW) on the block to find a valid hash
func (bc *Blockchain) mineBlock(block Block) string {
    for {
        // Calculate the block's hash based on its current state
        hash := calculateBlockHash(block)

        // Check if the hash meets the required difficulty (leading zeroes)
        if strings.HasPrefix(hash, strings.Repeat("0", block.Difficulty)) {
            fmt.Printf("Block mined successfully with nonce %d: %s\n", block.Nonce, hash)
            return hash
        }

        // Increment the nonce and try again
        block.Nonce++
    }
}



// ValidateBlock verifies the block's hash and ensures it matches the difficulty
func (bc *Blockchain) ValidateBlock(block Block) bool {
    expectedHash := calculateBlockHash(block)
    if block.Hash != expectedHash {
        fmt.Printf("Block %d validation failed: Invalid hash.\n", block.Index)
        return false
    }

    if !strings.HasPrefix(block.Hash, strings.Repeat("0", block.Difficulty)) {
        fmt.Printf("Block %d validation failed: Hash does not meet difficulty.\n", block.Index)
        return false
    }

    // Convert Blockchain.Block to ledger.Block before validation
    ledgerBlock := ConvertToLedgerBlock(block)

    // Ensure the block exists in the ledger
    if err := bc.Ledger.BlockchainConsensusCoinLedger.ValidateBlock(ledgerBlock); err != nil { 
        fmt.Printf("Block %d validation failed: Not present in the ledger (%v).\n", block.Index, err)
        return false
    }

    fmt.Printf("Block %d successfully validated.\n", block.Index)
    return true
}



// AddSubBlocksToBlock adds sub-blocks to a new block and mines it using PoW
func (bc *Blockchain) AddSubBlocksToBlock(subBlocks []SubBlock, pow *PoW) {
    prevHash := ""
    if len(bc.Chain) > 0 { // Use bc.Chain instead of bc.Blocks
        prevHash = bc.Chain[len(bc.Chain)-1].Hash
    }

    if len(subBlocks) > 1000 {
        subBlocks = subBlocks[:1000] // Limit to 1000 sub-blocks per block
    }

    // Create a new block with sub-blocks and mine it
    bc.NewBlock(subBlocks, prevHash, pow) // Pass the PoW instance as the third argument
}


// calculateBlockHash calculates the SHA-256 hash of the block's content.
func calculateBlockHash(block Block) string {
	// Convert block data into a string to be hashed
	blockData := fmt.Sprintf(
		"%s%d%s%s%d%d%f%s", // Format: BlockID, Index, PrevHash, Timestamp, Nonce, Difficulty, MinerReward, Validators (joined)
		block.BlockID,
		block.Index,
		block.PrevHash,
		block.Timestamp.String(),   // Convert timestamp to string
		block.Nonce,
		block.Difficulty,
		block.MinerReward,
		joinValidators(block.Validators), // Convert validators list to a single string
	)

	// Create a SHA-256 hash of the block data
	hash := sha256.New()
	hash.Write([]byte(blockData))
	hashed := hash.Sum(nil)

	// Convert the hashed bytes to a hexadecimal string and return
	return hex.EncodeToString(hashed)
}

// Helper function to join validators into a single string
func joinValidators(validators []string) string {
	return fmt.Sprintf("%v", validators) // Join validators into a single string
}