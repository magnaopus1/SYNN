package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"     // Shared components like encryption, consensus, sub-blocks, etc.
	"synnergy_network/pkg/ledger"     // Blockchain and ledger-related components
)

// ConsensusMechanism represents the available consensus mechanisms.
type ConsensusMechanism string

const (
	ProofOfWork   ConsensusMechanism = "PoW"
	ProofOfStake  ConsensusMechanism = "PoS"
	ProofOfHistory ConsensusMechanism = "PoH"
)

// ConsensusSpecificNode represents a node that participates in the Synnergy Consensus using one specific mechanism.
type ConsensusSpecificNode struct {
	NodeID            string                        // Unique identifier for the node
	Blockchain        *ledger.Blockchain            // Local copy of the blockchain ledger
	ConsensusEngine   *common.SynnergyConsensus     // Consensus engine for validating transactions and blocks
	EncryptionService *common.Encryption            // Encryption service for secure communication
	NetworkManager    *common.NetworkManager        // Network manager for communication with other nodes
	ConsensusMechanism ConsensusMechanism           // The consensus mechanism used by this node (PoW, PoS, or PoH)
	SubBlocks         map[string]*common.SubBlock   // Sub-blocks that are part of blocks in the blockchain
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with the blockchain network
	SNVM              *synnergy_vm.VirtualMachine // Virtual Machine for executing smart contracts

}

// NewConsensusSpecificNode initializes a new consensus-specific node in the Synnergy Network.
func NewConsensusSpecificNode(nodeID string, blockchain *ledger.Blockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *common.NetworkManager, consensusMechanism ConsensusMechanism, syncInterval time.Duration) *ConsensusSpecificNode {
	return &ConsensusSpecificNode{
		NodeID:             nodeID,
		Blockchain:         blockchain,
		ConsensusEngine:    consensusEngine,
		EncryptionService:  encryptionService,
		NetworkManager:     networkManager,
		ConsensusMechanism: consensusMechanism,
		SubBlocks:          make(map[string]*common.SubBlock),
		SyncInterval:       syncInterval,
	}
}

// StartNode starts the consensus-specific node’s operations based on the selected consensus mechanism.
func (csn *ConsensusSpecificNode) StartNode() error {
	csn.mutex.Lock()
	defer csn.mutex.Unlock()

	// Start syncing with the blockchain network and listening for transactions.
	go csn.syncWithOtherNodes()
	go csn.listenForTransactions()

	fmt.Printf("Consensus-specific node %s started using %s mechanism.\n", csn.NodeID, csn.ConsensusMechanism)
	return nil
}

// syncWithOtherNodes handles syncing the blockchain with other nodes at regular intervals.
func (csn *ConsensusSpecificNode) syncWithOtherNodes() {
	ticker := time.NewTicker(csn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		csn.mutex.Lock()
		otherNodes := csn.NetworkManager.DiscoverOtherNodes(csn.NodeID)
		for _, node := range otherNodes {
			// Sync blockchain from each node.
			csn.syncBlockchainFromNode(node)
		}
		csn.mutex.Unlock()
	}
}

// syncBlockchainFromNode syncs the blockchain from a peer node.
func (csn *ConsensusSpecificNode) syncBlockchainFromNode(peerNode string) {
	peerBlockchain, err := csn.NetworkManager.RequestBlockchain(peerNode)
	if err != nil {
		fmt.Printf("Failed to sync blockchain from node %s: %v\n", peerNode, err)
		return
	}

	// Validate the blockchain and update the local copy if necessary.
	if csn.ConsensusEngine.ValidateBlockchain(peerBlockchain, csn.ConsensusMechanism) {
		csn.Blockchain = csn.Blockchain.MergeWith(peerBlockchain)
		fmt.Printf("Blockchain synced successfully from node %s using %s mechanism.\n", peerNode, csn.ConsensusMechanism)
	} else {
		fmt.Printf("Blockchain sync from node %s failed validation using %s mechanism.\n", peerNode, csn.ConsensusMechanism)
	}
}

// listenForTransactions listens for incoming transactions and processes them.
func (csn *ConsensusSpecificNode) listenForTransactions() {
	for {
		transaction, err := csn.NetworkManager.ReceiveTransaction()
		if err != nil {
			fmt.Printf("Error receiving transaction: %v\n", err)
			continue
		}

		// Validate and process the transaction based on the consensus mechanism.
		err = csn.processTransaction(transaction)
		if err != nil {
			fmt.Printf("Transaction processing failed: %v\n", err)
		}
	}
}

// processTransaction processes and validates an incoming transaction based on the selected consensus mechanism.
func (csn *ConsensusSpecificNode) processTransaction(tx *ledger.Transaction) error {
	csn.mutex.Lock()
	defer csn.mutex.Unlock()

	// Validate the transaction using the selected consensus mechanism.
	if valid, err := csn.ConsensusEngine.ValidateTransaction(tx, csn.ConsensusMechanism); err != nil || !valid {
		return fmt.Errorf("invalid transaction: %v", err)
	}

	// Add the transaction to a sub-block.
	subBlockID := common.GenerateSubBlockID()
	subBlock := csn.createSubBlock(tx, subBlockID)

	// Store the sub-block and attempt to validate it into a block.
	csn.SubBlocks[subBlockID] = subBlock
	csn.tryValidateSubBlock(subBlock)

	fmt.Printf("Transaction %s processed successfully using %s mechanism.\n", tx.TransactionID, csn.ConsensusMechanism)
	return nil
}

// createSubBlock creates a sub-block from a validated transaction.
func (csn *ConsensusSpecificNode) createSubBlock(tx *ledger.Transaction, subBlockID string) *common.SubBlock {
	return &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: []*ledger.Transaction{tx},
		Timestamp:    time.Now(),
		NodeID:       csn.NodeID,
	}
}

// tryValidateSubBlock tries to validate a sub-block into a full block using the selected consensus mechanism.
func (csn *ConsensusSpecificNode) tryValidateSubBlock(subBlock *common.SubBlock) error {
	csn.mutex.Lock()
	defer csn.mutex.Unlock()

	// Validate the sub-block using the selected consensus mechanism.
	if csn.ConsensusEngine.ValidateSubBlock(subBlock, csn.ConsensusMechanism) {
		// Add the sub-block to the blockchain as part of a full block.
		block := csn.Blockchain.AddSubBlock(subBlock)

		// Notify the network of the new block.
		err := csn.NetworkManager.BroadcastNewBlock(block)
		if err != nil {
			return fmt.Errorf("failed to broadcast new block: %v", err)
		}

		fmt.Printf("Sub-block %s validated and added to blockchain using %s mechanism.\n", subBlock.SubBlockID, csn.ConsensusMechanism)
		return nil
	}

	return errors.New("failed to validate sub-block")
}

// Security and Encryption

// ApplySecurityProtocols applies the necessary encryption and security measures for node communications and data.
func (csn *ConsensusSpecificNode) ApplySecurityProtocols() error {
	// Apply end-to-end encryption for all data communications.
	err := csn.EncryptionService.ApplySecurity(csn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to apply security protocols: %v", err)
	}

	fmt.Printf("Security protocols applied successfully for consensus-specific node %s.\n", csn.NodeID)
	return nil
}

// SwitchConsensusMechanism allows the node to switch between consensus mechanisms if needed.
func (csn *ConsensusSpecificNode) SwitchConsensusMechanism(newMechanism ConsensusMechanism) error {
	csn.mutex.Lock()
	defer csn.mutex.Unlock()

	// Update the consensus mechanism.
	csn.ConsensusMechanism = newMechanism
	fmt.Printf("Consensus mechanism switched to %s for node %s.\n", newMechanism, csn.NodeID)
	return nil
}
