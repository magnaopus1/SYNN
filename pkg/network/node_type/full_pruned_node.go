package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"      // All shared components go here
	"synnergy_network/pkg/ledger"      // Blockchain and ledger-related components
	"synnergy_network/pkg/network" 
)

// FullPrunedNode represents a full pruned node in the blockchain network.
type FullPrunedNode struct {
	NodeID            string                         // Unique identifier for the node
	PrunedBlockchain  *ledger.PrunedBlockchain       // Local pruned version of the blockchain ledger
	ConsensusEngine   *common.SynnergyConsensus      // Consensus engine for validating transactions and sub-blocks
	EncryptionService *common.Encryption             // Encryption service for secure communication
	NetworkManager    *network.NetworkManager         // Network manager for communication with other nodes
	SubBlocks         map[string]*common.SubBlock    // Sub-blocks waiting to be validated
	mutex             sync.Mutex                     // Mutex for thread-safe operations
	SyncInterval      time.Duration                  // Interval for syncing with other full nodes
	SNVM              *common.VirtualMachine // Virtual Machine for executing smart contracts
}

// NewFullPrunedNode initializes a new full pruned node in the blockchain network.
func NewFullPrunedNode(nodeID string, prunedBlockchain *ledger.PrunedBlockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *network.NetworkManager, syncInterval time.Duration) *FullPrunedNode {
	return &FullPrunedNode{
		NodeID:            nodeID,
		PrunedBlockchain:  prunedBlockchain,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SubBlocks:         make(map[string]*common.SubBlock),
		SyncInterval:      syncInterval,
		SNVM:              common.NewVirtualMachine(), // Initialize the VM for smart contract execution
	}
}

// StartNode starts the full pruned node's operations, including syncing, validating, and running contracts.
func (fn *FullPrunedNode) StartNode() error {
	fn.mutex.Lock()
	defer fn.mutex.Unlock()

	// Start syncing the blockchain and listening for transactions.
	go fn.syncWithOtherNodes()
	go fn.listenForTransactions()

	fmt.Printf("Full pruned node %s started successfully.\n", fn.NodeID)
	return nil
}

// syncWithOtherNodes handles syncing the pruned blockchain with other full nodes at regular intervals.
func (fn *FullPrunedNode) syncWithOtherNodes() {
	ticker := time.NewTicker(fn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		fn.mutex.Lock()
		otherNodes := fn.NetworkManager.DiscoverOtherFullNodes(fn.NodeID)
		for _, node := range otherNodes {
			fn.syncBlockchainFromNode(node)
		}
		fn.mutex.Unlock()
	}
}

// syncBlockchainFromNode syncs the blockchain from a peer node, ensuring that only the latest and necessary data is stored.
func (fn *FullPrunedNode) syncBlockchainFromNode(peerNode string) {
	peerPrunedBlockchain, err := fn.NetworkManager.RequestPrunedBlockchain(peerNode)
	if err != nil {
		fmt.Printf("Failed to sync pruned blockchain from node %s: %v\n", peerNode, err)
		return
	}

	if fn.ConsensusEngine.ValidatePrunedBlockchain(peerPrunedBlockchain) {
		fn.PrunedBlockchain = fn.PrunedBlockchain.MergeWith(peerPrunedBlockchain)
		fmt.Printf("Pruned blockchain synced successfully from node %s.\n", peerNode)
	} else {
		fmt.Printf("Pruned blockchain sync from node %s failed validation.\n", peerNode)
	}
}

// listenForTransactions listens for incoming transactions and processes them.
func (fn *FullPrunedNode) listenForTransactions() {
	for {
		transaction, err := fn.NetworkManager.ReceiveTransaction()
		if err != nil {
			fmt.Printf("Error receiving transaction: %v\n", err)
			continue
		}

		err = fn.processTransaction(transaction)
		if err != nil {
			fmt.Printf("Transaction processing failed: %v\n", err)
		}
	}
}

// processTransaction processes and validates an incoming transaction, adding it to the pruned blockchain.
func (fn *FullPrunedNode) processTransaction(tx *ledger.Transaction) error {
	fn.mutex.Lock()
	defer fn.mutex.Unlock()

	if valid, err := fn.ConsensusEngine.ValidateTransaction(tx); err != nil || !valid {
		return fmt.Errorf("invalid transaction: %v", err)
	}

	subBlockID := common.GenerateSubBlockID()
	subBlock := fn.createSubBlock(tx, subBlockID)

	fn.SubBlocks[subBlockID] = subBlock
	return fn.tryValidateSubBlock(subBlock)
}

// createSubBlock creates a sub-block from a validated transaction.
func (fn *FullPrunedNode) createSubBlock(tx *ledger.Transaction, subBlockID string) *common.SubBlock {
	return &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: []*ledger.Transaction{tx},
		Timestamp:    time.Now(),
		NodeID:       fn.NodeID,
	}
}

// tryValidateSubBlock tries to validate a sub-block into a pruned block.
func (fn *FullPrunedNode) tryValidateSubBlock(subBlock *common.SubBlock) error {
	fn.mutex.Lock()
	defer fn.mutex.Unlock()

	if fn.ConsensusEngine.ValidateSubBlock(subBlock) {
		prunedBlock := fn.PrunedBlockchain.AddSubBlock(subBlock)
		err := fn.NetworkManager.BroadcastNewPrunedBlock(prunedBlock)
		if err != nil {
			return fmt.Errorf("failed to broadcast new pruned block: %v", err)
		}

		fmt.Printf("Sub-block %s validated and added to pruned blockchain.\n", subBlock.SubBlockID)
		return nil
	}

	return errors.New("failed to validate sub-block")
}

// Pruning Mechanism - Automatically prunes old transactions and blocks.
func (fn *FullPrunedNode) PruneOldBlocks() error {
	fn.mutex.Lock()
	defer fn.mutex.Unlock()

	// Custom logic to prune old blocks based on certain conditions (e.g., block age, disk space usage).
	prunedCount, err := fn.PrunedBlockchain.Prune()
	if err != nil {
		return fmt.Errorf("error during pruning: %v", err)
	}

	fmt.Printf("%d blocks pruned from the blockchain.\n", prunedCount)
	return nil
}
