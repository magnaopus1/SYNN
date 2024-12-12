package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/network"
	"synnergy_network_demo/common"
	"synnergy_network_demo/synnergy_vm"
)

// ValidatorNode represents a node responsible for validating transactions into sub-blocks but does not mine full blocks.
type ValidatorNode struct {
	NodeID            string                        // Unique identifier for the validator node
	SubBlocks         map[string]*common.SubBlock   // Sub-blocks validated by this node
	ConsensusEngine   *synnergy_consensus.Engine    // Consensus engine for validating transactions
	EncryptionService *encryption.Encryption        // Encryption service for securing communications
	NetworkManager    *network.NetworkManager       // Manages communication with other nodes
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with other validator nodes
	SNVM              *synnergy_vm.VirtualMachine   // Synnergy Network Virtual Machine instance for executing smart contracts
	ConsensusType     string                        // Specifies the consensus mechanism used (PoS or PoH)
}

// NewValidatorNode initializes a new ValidatorNode for the network.
func NewValidatorNode(nodeID string, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, syncInterval time.Duration, snvm *synnergy_vm.VirtualMachine, consensusType string) *ValidatorNode {
	if consensusType != "PoS" && consensusType != "PoH" {
		panic("Invalid consensus type. Must be 'PoS' or 'PoH'.")
	}
	return &ValidatorNode{
		NodeID:            nodeID,
		SubBlocks:         make(map[string]*common.SubBlock),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SyncInterval:      syncInterval,
		SNVM:              snvm,
		ConsensusType:     consensusType,
	}
}

// StartNode starts the validator node and initiates syncing with other nodes.
func (vn *ValidatorNode) StartNode() error {
	vn.mutex.Lock()
	defer vn.mutex.Unlock()

	// Start syncing with other validator nodes and listening for transactions.
	go vn.syncWithOtherNodes()
	go vn.listenForTransactions()

	fmt.Printf("Validator node %s started successfully using %s consensus.\n", vn.NodeID, vn.ConsensusType)
	return nil
}

// syncWithOtherNodes syncs sub-blocks with other validator nodes at regular intervals.
func (vn *ValidatorNode) syncWithOtherNodes() {
	ticker := time.NewTicker(vn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		vn.mutex.Lock()
		otherNodes := vn.NetworkManager.DiscoverOtherValidatorNodes(vn.NodeID)
		for _, node := range otherNodes {
			vn.syncSubBlocksFromNode(node)
		}
		vn.mutex.Unlock()
	}
}

// syncSubBlocksFromNode fetches sub-blocks from a peer validator node for synchronization.
func (vn *ValidatorNode) syncSubBlocksFromNode(peerNode string) {
	peerSubBlocks, err := vn.NetworkManager.RequestSubBlocks(peerNode)
	if err != nil {
		fmt.Printf("Failed to sync sub-blocks from node %s: %v\n", peerNode, err)
		return
	}

	// Validate and store the fetched sub-blocks.
	for _, subBlock := range peerSubBlocks {
		if vn.ConsensusEngine.ValidateSubBlock(subBlock) {
			vn.SubBlocks[subBlock.SubBlockID] = subBlock
			fmt.Printf("Sub-block %s synced successfully from node %s.\n", subBlock.SubBlockID, peerNode)
		} else {
			fmt.Printf("Sub-block %s from node %s failed validation.\n", subBlock.SubBlockID, peerNode)
		}
	}
}

// listenForTransactions listens for incoming transactions and processes them into sub-blocks.
func (vn *ValidatorNode) listenForTransactions() {
	for {
		transaction, err := vn.NetworkManager.ReceiveTransaction()
		if err != nil {
			fmt.Printf("Error receiving transaction: %v\n", err)
			continue
		}

		// Process the received transaction into a sub-block.
		err = vn.processTransaction(transaction)
		if err != nil {
			fmt.Printf("Transaction processing failed: %v\n", err)
		}
	}
}

// processTransaction validates the transaction and processes it into a sub-block.
func (vn *ValidatorNode) processTransaction(tx *ledger.Transaction) error {
	vn.mutex.Lock()
	defer vn.mutex.Unlock()

	// Validate the transaction using the consensus engine (PoS or PoH).
	if valid, err := vn.ConsensusEngine.ValidateTransaction(tx, vn.ConsensusType); err != nil || !valid {
		return fmt.Errorf("invalid transaction: %v", err)
	}

	// Create and store a new sub-block from the validated transaction.
	subBlockID := common.GenerateSubBlockID()
	subBlock := vn.createSubBlock(tx, subBlockID)

	// Store the sub-block for further network-wide validation.
	vn.SubBlocks[subBlockID] = subBlock
	fmt.Printf("Transaction %s processed into sub-block %s.\n", tx.TransactionID, subBlockID)
	return nil
}

// createSubBlock generates a new sub-block for a validated transaction.
func (vn *ValidatorNode) createSubBlock(tx *ledger.Transaction, subBlockID string) *common.SubBlock {
	return &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: []*ledger.Transaction{tx},
		Timestamp:    time.Now(),
		NodeID:       vn.NodeID,
	}
}

// tryBroadcastSubBlock broadcasts a validated sub-block to other nodes for further validation.
func (vn *ValidatorNode) tryBroadcastSubBlock(subBlock *common.SubBlock) error {
	vn.mutex.Lock()
	defer vn.mutex.Unlock()

	// Broadcast the sub-block across the network.
	err := vn.NetworkManager.BroadcastSubBlock(subBlock)
	if err != nil {
		return fmt.Errorf("failed to broadcast sub-block: %v", err)
	}

	fmt.Printf("Sub-block %s broadcasted to the network.\n", subBlock.SubBlockID)
	return nil
}

// EncryptData encrypts sensitive data before it is transmitted to other nodes.
func (vn *ValidatorNode) EncryptData(data []byte) ([]byte, error) {
	encryptedData, err := vn.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return encryptedData, nil
}

// DecryptData decrypts received encrypted data from other nodes.
func (vn *ValidatorNode) DecryptData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := vn.EncryptionService.DecryptData(encryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decryptedData, nil
}
