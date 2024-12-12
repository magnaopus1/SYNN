package plasma

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)


// CreatePlasmaChain initializes the Plasma childchain and creates the genesis block
func CreatePlasmaChain(config *common.PlasmaChainConfig) (*common.PlasmaChain, error) {
	// Validate that all necessary configurations are provided
	if config.ChainID == "" || config.GenesisBlockID == "" || config.PlasmaCore == nil || config.Ledger == nil {
		return nil, errors.New("missing configuration details for Plasma chain creation")
	}

	// Create the genesis block with an empty list of sub-blocks
	genesisBlock := &common.PlasmaBlock{
		BlockID:       config.GenesisBlockID,
		SubBlocks:     []*common.PlasmaSubBlock{}, // No sub-blocks in genesis block
		PreviousBlock: "",
		Timestamp:     config.GenesisTimestamp,
		MerkleRoot:    "", // No transactions, so no Merkle root needed
	}

	// Record the genesis block in the ledger
	err := config.Ledger.RecordBlockCreation(config.GenesisBlockID, "", config.GenesisTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to record genesis block: %v", err)
	}

	// Create the Plasma chain with the genesis block
	plasmaChain := &common.PlasmaChain{
		ChainID:        config.ChainID,
		GenesisBlock:   genesisBlock,
		CurrentBlockID: config.GenesisBlockID,
		Blocks:         map[string]*common.PlasmaBlock{config.GenesisBlockID: genesisBlock},
		Ledger:         config.Ledger,
		Core:           config.PlasmaCore,
		Encryption:     config.EncryptionService,
		NetworkManager: config.NetworkManager,
	}

	fmt.Printf("Plasma chain %s created with genesis block %s\n", config.ChainID, config.GenesisBlockID)
	return plasmaChain, nil
}

// AddNodeToPlasmaChain adds a new node to the Plasma childchain network
func (pc *common.PlasmaChain) AddNodeToPlasmaChain(nodeID string) error {
	err := pc.NetworkManager.AddNode(nodeID)
	if err != nil {
		return fmt.Errorf("failed to add node to Plasma chain: %v", err)
	}

	fmt.Printf("Node %s added to Plasma chain %s\n", nodeID, pc.ChainID)
	return nil
}

// SyncChainStateWithNetwork ensures the Plasma chain is in sync with the network state
func (pc *common.PlasmaChain) SyncChainStateWithNetwork() error {
	err := pc.NetworkManager.SyncChildChainState()
	if err != nil {
		return fmt.Errorf("failed to sync Plasma chain with network: %v", err)
	}

	fmt.Printf("Plasma chain %s synced with network\n", pc.ChainID)
	return nil
}


// FinalizePlasmaBlock finalizes a block once it reaches 1000 sub-blocks
func (pc *common.PlasmaChain) FinalizePlasmaBlock(blockID string) error {
	block, exists := pc.Blocks[blockID]
	if !exists {
		return fmt.Errorf("block %s not found", blockID)
	}

	if len(block.SubBlocks) != 1000 {
		return errors.New("block cannot be finalized until it contains 1000 sub-blocks")
	}

	// Mark the block as finalized in the ledger
	err := pc.Ledger.RecordBlockFinalization(blockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to finalize block: %v", err)
	}

	fmt.Printf("Block %s finalized with 1000 sub-blocks\n", blockID)
	return nil
}

// AddSubBlockToPlasmaChain adds a sub-block to a block in the Plasma chain
func (pc *common.PlasmaChain) AddSubBlockToPlasmaChain(blockID string, subBlock *common.PlasmaSubBlock) error {
	block, exists := pc.Blocks[blockID]
	if !exists {
		return fmt.Errorf("block %s not found", blockID)
	}

	if len(block.SubBlocks) >= 1000 {
		return errors.New("cannot add more sub-blocks, block already contains 1000 sub-blocks")
	}

	block.SubBlocks = append(block.SubBlocks, subBlock)

	// Log the addition of the sub-block in the ledger
	err := pc.Ledger.RecordSubBlockAddition(blockID, subBlock.SubBlockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sub-block addition: %v", err)
	}

	fmt.Printf("Sub-block %s added to block %s\n", subBlock.SubBlockID, blockID)
	return nil
}

// GetPlasmaBlock retrieves the current block based on block ID
func (pc *common.PlasmaChain) GetPlasmaBlock(blockID string) (*common.PlasmaBlock, error) {
	block, exists := pc.Blocks[blockID]
	if !exists {
		return nil, fmt.Errorf("block %s not found", blockID)
	}

	fmt.Printf("Block %s retrieved from Plasma chain %s\n", blockID, pc.ChainID)
	return block, nil
}

// EncryptSubBlockData encrypts sub-block data before submission
func (pc *common.PlasmaChain) EncryptSubBlockData(subBlock *common.PlasmaSubBlock) error {
	for _, tx := range subBlock.Transactions {
		encryptedTx, err := pc.Encryption.EncryptData([]byte(tx.UTXOID), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt transaction data: %v", err)
		}
		tx.UTXOID = string(encryptedTx)
	}

	fmt.Printf("Sub-block %s data encrypted\n", subBlock.SubBlockID)
	return nil
}
