package high_availability

import (
    "fmt"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    
)


// NewChainForkManager initializes a new fork manager with a reference to the ledger
func NewChainForkManager(ledgerInstance *ledger.Ledger) *ChainForkManager {
    return &ChainForkManager{
        LedgerInstance: ledgerInstance,
        ForkedChains:   [][]common.Block{},
    }
}

// DetectFork checks if there is a fork in the blockchain based on differing sub-block hashes
func (fm *ChainForkManager) DetectFork(newBlock common.Block) bool {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    // Get the last block and check for errors
    lastBlock, err := fm.LedgerInstance.GetLastBlock() // Assuming this returns *Block
    if err != nil {
        fmt.Printf("Error retrieving last block: %v\n", err)
        return false
    }

    // Ensure lastBlock is not nil before accessing its fields
    if lastBlock == nil {
        fmt.Println("No last block found, cannot detect fork.")
        return false
    }

    // Dereference lastBlock if necessary
    lastBlockValue := *lastBlock // Dereference the pointer

    // If the previous block hashes are different, it indicates a fork
    if lastBlockValue.Hash != newBlock.PrevHash {
        fmt.Printf("Fork detected: Block %d has a different previous hash than expected.\n", newBlock.Index)
        fm.handleFork(newBlock)
        return true
    }

    return false
}

// ConvertTransactions converts []Blockchain.Transaction to []Ledger.Transaction.
func ConvertTransactions(bcTransactions []common.Transaction) []ledger.Transaction {
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

// ConvertLedgerTransactionsToCommon converts []ledger.Transaction to []common.Transaction.
func ConvertLedgerTransactionsToCommon(ledgerTransactions []ledger.Transaction) []common.Transaction {
	var commonTransactions []common.Transaction
	for _, tx := range ledgerTransactions {
		commonTransactions = append(commonTransactions, common.Transaction{
			TransactionID:  tx.TransactionID,
			FromAddress:    tx.FromAddress,
			ToAddress:      tx.ToAddress,
			Amount:         tx.Amount,
			Timestamp:      tx.Timestamp,
			Signature:      tx.Signature,
			Status:         tx.Status,
			TokenID:        tx.TokenID,
			ExecutionResult: tx.ExecutionResult,
			Fee:            tx.Fee,
			FrozenAmount:   tx.FrozenAmount,
			RefundAmount:   tx.RefundAmount,
			ReversalRequested: tx.ReversalRequested,
		})
	}
	return commonTransactions
}

// ConvertLedgerPoHProofToCommon converts ledger.PoHProof to common.PoHProof.
func ConvertLedgerPoHProofToCommon(ledgerPoHProof ledger.PoHProof) common.PoHProof {
	return common.PoHProof{
		Sequence:  ledgerPoHProof.Sequence,
		Timestamp: ledgerPoHProof.Timestamp,
		Hash:      ledgerPoHProof.Hash,
	}
}




// ConvertSubBlocks converts []Blockchain.SubBlock to []Ledger.SubBlock.
func ConvertSubBlocks(bcSubBlocks []common.SubBlock) []ledger.SubBlock {
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
func ConvertToLedgerBlock(bcBlock common.Block) ledger.Block {
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
func ConvertPoHProof(bcPoHProof common.PoHProof) ledger.PoHProof {
	return ledger.PoHProof{
		Sequence:  bcPoHProof.Sequence,   // Map the sequence number
		Timestamp: bcPoHProof.Timestamp,  // Map the timestamp
		Hash:      bcPoHProof.Hash,       // Map the hash generated by PoH
	}
}




/// handleFork manages the newly detected fork by storing the alternative chain
func (fm *ChainForkManager) handleFork(newBlock common.Block) {
    fmt.Printf("Handling fork for Block %d...\n", newBlock.Index)

    // Get the blocks from the ledger (assuming GetBlocks returns []ledger.Block)
    ledgerBlocks := fm.LedgerInstance.GetBlocks() // No second return value

    // Convert ledgerBlocks to common.Block if necessary
    var forkChain []common.Block
    for _, blk := range ledgerBlocks {
        // Convert each ledger.Block to common.Block
        forkChain = append(forkChain, ConvertToCommonBlock(blk)) // Use the new conversion function
    }

    // Add the new block to the forked chain
    forkChain = append(forkChain, newBlock)

    // Store the forked chain in ForkedChains
    fm.ForkedChains = append(fm.ForkedChains, forkChain)
    fmt.Printf("Stored forked chain with %d blocks.\n", len(forkChain))
}


// ConvertToCommonBlock converts a Ledger.Block to a Blockchain.Block (common.Block).
func ConvertToCommonBlock(ledgerBlock ledger.Block) common.Block {
	return common.Block{
		BlockID:     ledgerBlock.BlockID,
		Index:       ledgerBlock.Index,
		Timestamp:   ledgerBlock.Timestamp,
		SubBlocks:   ConvertLedgerSubBlocksToCommon(ledgerBlock.SubBlocks), // Assuming a sub-block conversion function
		PrevHash:    ledgerBlock.PrevHash,
		Hash:        ledgerBlock.Hash,
		Nonce:       ledgerBlock.Nonce,
		Difficulty:  ledgerBlock.Difficulty,
		MinerReward: ledgerBlock.MinerReward,
		Validators:  ledgerBlock.Validators,
		Status:      ledgerBlock.Status, // Ensure that both types have this field or remove if unnecessary
	}
}

// ConvertLedgerSubBlocksToCommon converts []ledger.SubBlock to []common.SubBlock.
func ConvertLedgerSubBlocksToCommon(ledgerSubBlocks []ledger.SubBlock) []common.SubBlock {
	var commonSubBlocks []common.SubBlock
	for _, sb := range ledgerSubBlocks {
		commonSubBlocks = append(commonSubBlocks, common.SubBlock{
			SubBlockID:   sb.SubBlockID,
			Index:        sb.Index,
			Timestamp:    sb.Timestamp,
			Transactions: ConvertLedgerTransactionsToCommon(sb.Transactions), // Convert transactions
			Validator:    sb.Validator,
			PrevHash:     sb.PrevHash,
			Hash:         sb.Hash,
			PoHProof:     ConvertLedgerPoHProofToCommon(sb.PoHProof),         // Convert PoHProof
			Status:       sb.Status,
			Signature:    sb.Signature,
		})
	}
	return commonSubBlocks
}




// ResolveFork resolves a detected fork by selecting the longest valid chain
func (fm *ChainForkManager) ResolveFork() {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    fmt.Println("Resolving fork...")

    // Get the current chain (assuming it returns []ledger.Block)
    currentChain := fm.LedgerInstance.GetBlocks()

    // Convert []ledger.Block to []*common.Block if necessary
    longestChain := make([]*common.Block, len(currentChain))
    for i := range currentChain {
        longestChain[i] = ConvertLedgerBlockToCommon(&currentChain[i]) // Convert to *common.Block
    }

    // Iterate through the forked chains to find the longest valid chain
    for _, forkedChain := range fm.ForkedChains {
        // Convert []common.Block (forkedChain) to []*common.Block
        forkedChainPtr := make([]*common.Block, len(forkedChain))
        for i := range forkedChain {
            forkedChainPtr[i] = &forkedChain[i] // Take the address of each block
        }

        if len(forkedChainPtr) > len(longestChain) {
            fmt.Printf("Forked chain with %d blocks is longer than the current chain with %d blocks.\n", len(forkedChainPtr), len(longestChain))
            longestChain = forkedChainPtr // Update longestChain to the forked chain
        }
    }

    // If the longest chain is not the current ledger's chain, replace it
    if len(longestChain) > len(currentChain) {
        fmt.Println("Switching to the longest valid chain.")
        
        // Convert longestChain []*common.Block to []*ledger.Block
        ledgerChain := make([]*ledger.Block, len(longestChain))
        for i := range longestChain {
            ledgerChain[i] = ConvertCommonBlockToLedger(longestChain[i]) // Convert to *ledger.Block
        }

        fm.LedgerInstance.SetChain(ledgerChain) // Assuming SetChain accepts []*ledger.Block
        fmt.Printf("Ledger updated to the longest chain with %d blocks.\n", len(longestChain))
    } else {
        fmt.Println("Current chain is already the longest. No fork resolution needed.")
    }

    // Clear the forked chains after resolving
    fm.ForkedChains = nil
}


// ConvertLedgerBlockToCommon converts a ledger.Block to a common.Block
func ConvertLedgerBlockToCommon(ledgerBlock *ledger.Block) *common.Block {
    return &common.Block{
        BlockID:     ledgerBlock.BlockID,
        Index:       ledgerBlock.Index,
        Timestamp:   ledgerBlock.Timestamp,
        SubBlocks:   ConvertLedgerSubBlocksToCommon(ledgerBlock.SubBlocks),
        PrevHash:    ledgerBlock.PrevHash,
        Hash:        ledgerBlock.Hash,
        Nonce:       ledgerBlock.Nonce,
        Difficulty:  ledgerBlock.Difficulty,
        MinerReward: ledgerBlock.MinerReward,
        Validators:  ledgerBlock.Validators,
        Status:      ledgerBlock.Status,
    }
}

// ConvertCommonBlockToLedger converts a common.Block to a ledger.Block
func ConvertCommonBlockToLedger(commonBlock *common.Block) *ledger.Block {
    return &ledger.Block{
        BlockID:     commonBlock.BlockID,
        Index:       commonBlock.Index,
        Timestamp:   commonBlock.Timestamp,
        SubBlocks:   ConvertCommonSubBlocksToLedger(commonBlock.SubBlocks),
        PrevHash:    commonBlock.PrevHash,
        Hash:        commonBlock.Hash,
        Nonce:       commonBlock.Nonce,
        Difficulty:  commonBlock.Difficulty,
        MinerReward: commonBlock.MinerReward,
        Validators:  commonBlock.Validators,
        Status:      commonBlock.Status,
    }
}

// ConvertCommonSubBlocksToLedger converts []common.SubBlock to []ledger.SubBlock.
func ConvertCommonSubBlocksToLedger(commonSubBlocks []common.SubBlock) []ledger.SubBlock {
	var ledgerSubBlocks []ledger.SubBlock
	for _, sb := range commonSubBlocks {
		ledgerSubBlocks = append(ledgerSubBlocks, ledger.SubBlock{
			SubBlockID:   sb.SubBlockID,
			Index:        sb.Index,
			Timestamp:    sb.Timestamp,
			Transactions: ConvertCommonTransactionsToLedger(sb.Transactions), // Convert transactions
			Validator:    sb.Validator,
			PrevHash:     sb.PrevHash,
			Hash:         sb.Hash,
			PoHProof:     ConvertCommonPoHProofToLedger(sb.PoHProof),         // Convert PoHProof
			Status:       sb.Status,
			Signature:    sb.Signature,
		})
	}
	return ledgerSubBlocks
}


// ConvertCommonTransactionsToLedger converts []common.Transaction to []ledger.Transaction.
func ConvertCommonTransactionsToLedger(commonTransactions []common.Transaction) []ledger.Transaction {
	var ledgerTransactions []ledger.Transaction
	for _, tx := range commonTransactions {
		ledgerTransactions = append(ledgerTransactions, ledger.Transaction{
			TransactionID:   tx.TransactionID,
			FromAddress:     tx.FromAddress,
			ToAddress:       tx.ToAddress,
			Amount:          tx.Amount,
			Timestamp:       tx.Timestamp,
			Signature:       tx.Signature,
			Status:          tx.Status,
			TokenID:         tx.TokenID,
			ExecutionResult: tx.ExecutionResult,
			Fee:             tx.Fee,
			FrozenAmount:    tx.FrozenAmount,
			RefundAmount:    tx.RefundAmount,
			ReversalRequested: tx.ReversalRequested,
		})
	}
	return ledgerTransactions
}

// ConvertCommonPoHProofToLedger converts common.PoHProof to ledger.PoHProof.
func ConvertCommonPoHProofToLedger(commonPoHProof common.PoHProof) ledger.PoHProof {
	return ledger.PoHProof{
		Sequence:  commonPoHProof.Sequence,
		Timestamp: commonPoHProof.Timestamp,
		Hash:      commonPoHProof.Hash,
	}
}




// ValidateChain verifies the validity of a forked chain before resolving
func (fm *ChainForkManager) ValidateChain(forkedChain []common.Block) bool {
    fmt.Printf("Validating forked chain with %d blocks...\n", len(forkedChain))

    // Get the last block from the ledger and handle any potential errors
    lastBlock, err := fm.LedgerInstance.GetLastBlock()
    if err != nil {
        fmt.Printf("Error retrieving last block: %v\n", err)
        return false
    }

    // Ensure the first block in the fork matches the current chain's last block
    if forkedChain[0].PrevHash != lastBlock.Hash {
        fmt.Println("Forked chain validation failed: Mismatch in previous hash.")
        return false
    }

    // Validate each block in the fork
    for i := 1; i < len(forkedChain); i++ {
        currentBlock := forkedChain[i]
        prevBlock := forkedChain[i-1]

        // Check if the current block's previous hash matches the previous block's hash
        if currentBlock.PrevHash != prevBlock.Hash {
            fmt.Printf("Forked chain validation failed at block %d: Hash mismatch.\n", currentBlock.Index)
            return false
        }
    }

    fmt.Println("Forked chain is valid.")
    return true
}


// ForkRecovery handles recovery from forks, switching to the longest valid chain
func (fm *ChainForkManager) ForkRecovery() {
    fmt.Println("Attempting to recover from chain fork...")
    
    if len(fm.ForkedChains) == 0 {
        fmt.Println("No forks detected. No recovery needed.")
        return
    }

    // Validate all forked chains and resolve
    for _, forkedChain := range fm.ForkedChains {
        if fm.ValidateChain(forkedChain) {
            fm.ResolveFork()
            return
        }
    }

    fmt.Println("No valid forks found. Retaining the current chain.")
}
