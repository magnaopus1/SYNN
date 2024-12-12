package sidechains

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/consensus"
)


// NewSidechain creates a new sidechain with coin setup
func NewSidechain(chainID, parentChainID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus) *common.Sidechain {
	coinSetup := NewSidechainCoinSetup(ledgerInstance, encryptionService, consensus)
	return &common.Sidechain{
		ChainID:       chainID,
		ParentChainID: parentChainID,
		Blocks:        make(map[string]*common.SideBlock),
		SubBlocks:     make(map[string]*common.SubBlock),
		CoinSetup:     coinSetup,
		Consensus:     consensus,
		Ledger:        ledgerInstance,
		Encryption:    encryptionService,
	}
}

// CreateSubBlock creates a new sub-block in the sidechain
func (sc *common.Sidechain) CreateSubBlock(subBlockID string, transactions []*common.Transaction) (*common.SubBlock, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Validate if the sub-block already exists
	if _, exists := sc.SubBlocks[subBlockID]; exists {
		return nil, errors.New("sub-block already exists")
	}

	// Create the sub-block with transactions
	subBlock := &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: transactions,
		Timestamp:    time.Now(),
		MerkleRoot:   common.GenerateMerkleRoot(transactions), // Assuming a helper function for Merkle root generation
	}

	// Encrypt the sub-block transactions
	for _, tx := range transactions {
		encryptedTx, err := sc.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt transaction: %v", err)
		}
		tx.TxID = string(encryptedTx)
	}

	sc.SubBlocks[subBlockID] = subBlock

	fmt.Printf("Sub-block %s created with %d transactions\n", subBlockID, len(transactions))
	return subBlock, nil
}

// CreateBlock aggregates sub-blocks into a block on the sidechain
func (sc *common.Sidechain) CreateBlock(blockID, parentBlock string, subBlocks []*common.SubBlock) (*common.SideBlock, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Create the block and assign the sub-blocks
	block := &common.SideBlock{
		BlockID:     blockID,
		SubBlocks:   subBlocks,
		ParentBlock: parentBlock,
		Timestamp:   time.Now(),
		MerkleRoot:  common.GenerateMerkleRoot(subBlocks),
	}

	sc.Blocks[blockID] = block

	// Log the block creation in the ledger
	err := sc.Ledger.RecordBlockCreation(blockID, parentBlock, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log block creation: %v", err)
	}

	fmt.Printf("Block %s created with %d sub-blocks\n", blockID, len(subBlocks))
	return block, nil
}

// MintSidechainCoins allows the sidechain to mint coins
func (sc *common.Sidechain) MintSidechainCoins(coinID, issuerID string, amount float64) error {
	return sc.CoinSetup.MintCoins(coinID, issuerID, amount)
}

// TransferSidechainCoins transfers coins within the sidechain
func (sc *common.Sidechain) TransferSidechainCoins(coinID, senderID, receiverID string, amount float64) error {
	return sc.CoinSetup.TransferCoin(coinID, senderID, receiverID, amount)
}

// ValidateSubBlock validates a sub-block on the sidechain
func (sc *common.Sidechain) ValidateSubBlock(subBlockID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	subBlock, exists := sc.SubBlocks[subBlockID]
	if !exists {
		return fmt.Errorf("sub-block %s not found", subBlockID)
	}

	// Use Synnergy Consensus to validate the sub-block
	err := sc.Consensus.consensus.ValidateSubBlock(subBlock.SubBlockID, subBlock.Transactions)
	if err != nil {
		return fmt.Errorf("sub-block validation failed: %v", err)
	}

	fmt.Printf("Sub-block %s validated\n", subBlockID)
	return nil
}

// ValidateBlock validates a block on the sidechain
func (sc *common.Sidechain) ValidateBlock(blockID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	block, exists := sc.Blocks[blockID]
	if !exists {
		return fmt.Errorf("block %s not found", blockID)
	}

	// Validate each sub-block
	for _, subBlock := range block.SubBlocks {
		err := sc.ValidateSubBlock(subBlock.SubBlockID)
		if err != nil {
			return fmt.Errorf("block validation failed at sub-block %s: %v", subBlock.SubBlockID, err)
		}
	}

	fmt.Printf("Block %s validated successfully\n", blockID)
	return nil
}

// RetrieveSubBlock retrieves a sub-block by ID
func (sc *common.Sidechain) RetrieveSubBlock(subBlockID string) (*common.SubBlock, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	subBlock, exists := sc.SubBlocks[subBlockID]
	if !exists {
		return nil, fmt.Errorf("sub-block %s not found", subBlockID)
	}

	fmt.Printf("Retrieved sub-block %s\n", subBlockID)
	return subBlock, nil
}

// RetrieveBlock retrieves a block by ID
func (sc *common.Sidechain) RetrieveBlock(blockID string) (*common.SideBlock, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	block, exists := sc.Blocks[blockID]
	if !exists {
		return nil, fmt.Errorf("block %s not found", blockID)
	}

	fmt.Printf("Retrieved block %s\n", blockID)
	return block, nil
}

// FinalizeTransaction finalizes a transaction on the sidechain
func (sc *common.Sidechain) FinalizeTransaction(txID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Search for the transaction across sub-blocks
	for _, subBlock := range sc.SubBlocks {
		for _, tx := range subBlock.Transactions {
			if tx.TxID == txID && !tx.IsSpent {
				tx.IsSpent = true
				tx.SpentTime = time.Now()

				// Log the transaction as spent
				err := sc.Ledger.RecordTransactionSpent(txID, time.Now())
				if err != nil {
					return fmt.Errorf("failed to log transaction spent: %v", err)
				}

				fmt.Printf("Transaction %s finalized (spent)\n", txID)
				return nil
			}
		}
	}

	return fmt.Errorf("transaction %s not found or already spent", txID)
}

// DeleteSubBlock deletes a sub-block from the sidechain
func (sc *common.Sidechain) DeleteSubBlock(subBlockID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	_, exists := sc.SubBlocks[subBlockID]
	if !exists {
		return fmt.Errorf("sub-block %s not found", subBlockID)
	}

	delete(sc.SubBlocks, subBlockID)
	fmt.Printf("Sub-block %s deleted\n", subBlockID)
	return nil
}

// DeleteBlock deletes a block from the sidechain
func (sc *common.Sidechain) DeleteBlock(blockID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	_, exists := sc.Blocks[blockID]
	if !exists {
		return fmt.Errorf("block %s not found", blockID)
	}

	delete(sc.Blocks, blockID)
	fmt.Printf("Block %s deleted\n", blockID)
	return nil
}
