package sidechains

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)


// NewSidechainConsensus initializes the SidechainConsensus mechanism
func NewSidechainConsensus(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensusEngine *common.SynnergyConsensus) *common.SidechainConsensus {
	return &common.SidechainConsensus{
		Nodes:           make(map[string]*common.SidechainNode),
		PendingBlocks:   make(map[string]*common.Block),
		Ledger:          ledgerInstance,
		Encryption:      encryptionService,
		ConsensusEngine: consensusEngine,
	}
}

// AddNode adds a new SidechainNode to the sidechain consensus group
func (sc *common.SidechainConsensus) AddNode(node *common.SidechainNode) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if _, exists := sc.Nodes[node.NodeID]; exists {
		return errors.New("node already exists in the sidechain consensus")
	}

	sc.Nodes[node.NodeID] = node

	// Log the node addition in the ledger
	err := sc.Ledger.RecordNodeAddition(node.NodeID, node.IPAddress, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node addition: %v", err)
	}

	fmt.Printf("Node %s added to sidechain consensus\n", node.NodeID)
	return nil
}

// RemoveNode removes a SidechainNode from the sidechain consensus group
func (sc *common.SidechainConsensus) RemoveNode(nodeID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if _, exists := sc.Nodes[nodeID]; !exists {
		return errors.New("node not found in the sidechain consensus")
	}

	delete(sc.Nodes, nodeID)

	// Log the node removal in the ledger
	err := sc.Ledger.RecordNodeRemoval(nodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node removal: %v", err)
	}

	fmt.Printf("Node %s removed from sidechain consensus\n", nodeID)
	return nil
}

// ValidateBlock validates a block through Synnergy Consensus
func (sc *common.SidechainConsensus) ValidateBlock(block *common.Block) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Validate the block using Synnergy Consensus
	err := sc.ConsensusEngine.ValidateBlock(block.BlockID, block.SubBlocks)
	if err != nil {
		return fmt.Errorf("block %s validation failed: %v", block.BlockID, err)
	}

	// Log the block validation in the ledger
	err = sc.Ledger.RecordBlockValidation(block.BlockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log block validation: %v", err)
	}

	fmt.Printf("Block %s validated\n", block.BlockID)
	return nil
}

// ProposeBlock proposes a block for validation within the sidechain consensus
func (sc *common.SidechainConsensus) ProposeBlock(block *common.Block) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if _, exists := sc.PendingBlocks[block.BlockID]; exists {
		return errors.New("block is already proposed and awaiting validation")
	}

	sc.PendingBlocks[block.BlockID] = block

	// Encrypt the proposed block data
	encryptedBlockData, err := sc.Encryption.EncryptData([]byte(block.BlockID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt block data: %v", err)
	}

	// Broadcast the proposed block to all nodes for validation
	for _, node := range sc.Nodes {
		err := sc.broadcastToNode(node, encryptedBlockData)
		if err != nil {
			return fmt.Errorf("failed to broadcast proposed block to node %s: %v", node.NodeID, err)
		}
	}

	// Log the block proposal in the ledger
	err = sc.Ledger.RecordBlockProposal(block.BlockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log block proposal: %v", err)
	}

	fmt.Printf("Block %s proposed for consensus\n", block.BlockID)
	return nil
}

// broadcastToNode broadcasts the proposed block to a node for validation
func (sc *common.SidechainConsensus) broadcastToNode(node *common.SidechainNode, blockData []byte) error {
	// Use network manager to send the encrypted block data to the node
	err := sc.ConsensusEngine.NetworkManager.SendData(node, blockData)
	if err != nil {
		return fmt.Errorf("failed to send block data to node %s: %v", node.NodeID, err)
	}
	return nil
}

// ReconfigureNodes dynamically reconfigures the nodes in the sidechain consensus
func (sc *common.SidechainConsensus) ReconfigureNodes(newTopology map[string]*common.SidechainNode) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.Nodes = newTopology

	// Log the node reconfiguration in the ledger
	err := sc.Ledger.RecordNodeReconfiguration(time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node reconfiguration: %v", err)
	}

	fmt.Println("Node reconfiguration complete in sidechain consensus")
	return nil
}

// FinalizeBlock finalizes the block after consensus validation
func (sc *common.SidechainConsensus) FinalizeBlock(blockID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	block, exists := sc.PendingBlocks[blockID]
	if !exists {
		return fmt.Errorf("block %s not found in pending blocks", blockID)
	}

	// Finalize the block
	err := sc.ConsensusEngine.consensus.FinalizeBlock(block.BlockID)
	if err != nil {
		return fmt.Errorf("block %s finalization failed: %v", blockID, err)
	}

	// Remove the block from pending blocks
	delete(sc.PendingBlocks, blockID)

	// Log the block finalization in the ledger
	err = sc.Ledger.RecordBlockFinalization(blockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log block finalization: %v", err)
	}

	fmt.Printf("Block %s finalized in sidechain consensus\n", blockID)
	return nil
}

// ValidateSubBlock validates a sub-block in the sidechain consensus
func (sc *common.SidechainConsensus) ValidateSubBlock(subBlockID string, transactions []*common.Transaction) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Validate the sub-block using Synnergy Consensus
	err := sc.ConsensusEngine.consensus.ValidateSubBlock(subBlockID, transactions)
	if err != nil {
		return fmt.Errorf("sub-block %s validation failed: %v", subBlockID, err)
	}

	// Log the sub-block validation in the ledger
	err = sc.Ledger.RecordSubBlockValidation(subBlockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sub-block validation: %v", err)
	}

	fmt.Printf("Sub-block %s validated\n", subBlockID)
	return nil
}
