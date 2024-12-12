package plasma

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewPlasmaCore initializes the PlasmaCore for the childchain
func NewPlasmaCore(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.PlasmaCore {
	return &common.PlasmaCore{
		Blocks:          make(map[string]*common.PlasmaBlock),
		SubBlocks:       make(map[string]*common.PlasmaSubBlock),
		Ledger:          ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// CreateSubBlock handles the creation of a new Plasma sub-block
func (pc *common.PlasmaCore) CreateSubBlock(subBlockID string, txs []*common.PlasmaUTXO, subBlockIndex int) (*common.PlasmaSubBlock, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// Validate if the sub-block already exists
	if _, exists := pc.SubBlocks[subBlockID]; exists {
		return nil, errors.New("sub-block already exists")
	}

	// Create the sub-block with transactions
	subBlock := &common.PlasmaSubBlock{
		SubBlockID:    subBlockID,
		Transactions:  txs,
		SubBlockIndex: subBlockIndex,
		Timestamp:     time.Now(),
		MerkleRoot:    common.GenerateMerkleRoot(txs), // Assuming a helper function that generates Merkle root
	}

	// Log the creation of the sub-block
	err := pc.Ledger.RecordSubBlockCreation(subBlockID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log sub-block creation: %v", err)
	}

	// Encrypt transactions
	for _, tx := range txs {
		encryptedTx, err := pc.EncryptionService.EncryptData([]byte(tx.UTXOID), common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt transaction: %v", err)
		}
		tx.UTXOID = string(encryptedTx)
	}

	pc.SubBlocks[subBlockID] = subBlock
	fmt.Printf("Sub-block %s created with %d transactions\n", subBlockID, len(txs))
	return subBlock, nil
}

// CreateBlock aggregates 1000 sub-blocks into a full block
func (pc *common.PlasmaCore) CreateBlock(blockID, previousBlock string, subBlocks []*common.PlasmaSubBlock) (*common.PlasmaBlock, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if len(subBlocks) != 1000 {
		return nil, errors.New("a Plasma block must contain exactly 1000 sub-blocks")
	}

	// Create the block and associate the sub-blocks
	block := &common.PlasmaBlock{
		BlockID:       blockID,
		SubBlocks:     subBlocks,
		PreviousBlock: previousBlock,
		Timestamp:     time.Now(),
		MerkleRoot:    common.GenerateMerkleRoot(subBlocks),
	}

	// Log the block creation in the ledger
	err := pc.Ledger.RecordBlockCreation(blockID, previousBlock, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log block creation: %v", err)
	}

	pc.Blocks[blockID] = block
	fmt.Printf("Block %s created with 1000 sub-blocks\n", blockID)
	return block, nil
}

// ValidateSubBlock checks the validity of a Plasma sub-block using Synnergy Consensus
func (pc *common.PlasmaCore) ValidateSubBlock(subBlockID string) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	subBlock, exists := pc.SubBlocks[subBlockID]
	if !exists {
		return fmt.Errorf("sub-block %s not found", subBlockID)
	}

	// Validate transactions in the sub-block
	for _, tx := range subBlock.Transactions {
		if tx.IsSpent {
			return fmt.Errorf("invalid transaction found: UTXO %s is already spent", tx.UTXOID)
		}
	}

	// Log the validation event
	err := pc.Ledger.RecordSubBlockValidation(subBlockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sub-block validation: %v", err)
	}

	fmt.Printf("Sub-block %s validated\n", subBlockID)
	return nil
}

// ValidateBlock validates a block containing 1000 sub-blocks
func (pc *common.PlasmaCore) ValidateBlock(blockID string) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	block, exists := pc.Blocks[blockID]
	if !exists {
		return fmt.Errorf("block %s not found", blockID)
	}

	// Validate each sub-block within the block
	for _, subBlock := range block.SubBlocks {
		err := pc.ValidateSubBlock(subBlock.SubBlockID)
		if err != nil {
			return fmt.Errorf("validation failed for sub-block %s: %v", subBlock.SubBlockID, err)
		}
	}

	// Log the block validation event
	err := pc.Ledger.RecordBlockValidation(blockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log block validation: %v", err)
	}

	fmt.Printf("Block %s validated\n", blockID)
	return nil
}

// RetrieveSubBlock retrieves a sub-block from the Plasma childchain
func (pc *common.PlasmaCore) RetrieveSubBlock(subBlockID string) (*common.PlasmaSubBlock, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	subBlock, exists := pc.SubBlocks[subBlockID]
	if !exists {
		return nil, fmt.Errorf("sub-block %s not found", subBlockID)
	}

	fmt.Printf("Retrieved sub-block %s\n", subBlockID)
	return subBlock, nil
}

// RetrieveBlock retrieves a block from the Plasma childchain
func (pc *common.PlasmaCore) RetrieveBlock(blockID string) (*common.PlasmaBlock, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	block, exists := pc.Blocks[blockID]
	if !exists {
		return nil, fmt.Errorf("block %s not found", blockID)
	}

	fmt.Printf("Retrieved block %s\n", blockID)
	return block, nil
}

// FinalizeTransaction marks a UTXO transaction as spent
func (pc *common.PlasmaCore) FinalizeTransaction(txID string) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// Search for the UTXO across sub-blocks
	for _, subBlock := range pc.SubBlocks {
		for _, tx := range subBlock.Transactions {
			if tx.UTXOID == txID && !tx.IsSpent {
				tx.IsSpent = true
				tx.SpentTime = time.Now()

				// Log the transaction spend event
				err := pc.Ledger.RecordTransactionSpent(txID, time.Now())
				if err != nil {
					return fmt.Errorf("failed to log transaction spend: %v", err)
				}

				fmt.Printf("Transaction %s finalized (spent)\n", txID)
				return nil
			}
		}
	}

	return fmt.Errorf("transaction %s not found or already spent", txID)
}
