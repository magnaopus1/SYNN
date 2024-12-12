package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"   // Shared components for encryption, consensus, and more
	"synnergy_network/pkg/ledger"   // Blockchain and ledger-related components
	"synnergy_network/pkg/network"  // Network and communication management
	"synnergy_network/pkg/synnergy_vm" // Synnergy Virtual Machine for smart contract execution
)

// HybridNode represents a versatile node that performs multiple functions within the blockchain, including validation, indexing, and transaction handling.
type HybridNode struct {
	NodeID            string                        // Unique identifier for the node
	Blockchain        *ledger.Blockchain            // Local copy of the blockchain ledger
	ConsensusEngine   *common.SynnergyConsensus     // Consensus engine for validating transactions and blocks
	EncryptionService *common.Encryption            // Encryption service for secure communication
	NetworkManager    *network.NetworkManager       // Network manager for communication with other nodes
	SNVM              *common.VMInterface   // The Synnergy Network Virtual Machine
	DataIndex         map[string]interface{}        // Indexing structure for quick data retrieval
	SubBlocks         map[string]*common.SubBlock   // Sub-block storage before they are validated into full blocks
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with other nodes
}

// NewHybridNode initializes a new hybrid node in the blockchain network.
func NewHybridNode(nodeID string, blockchain *ledger.Blockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *network.NetworkManager, syncInterval time.Duration) *HybridNode {
	return &HybridNode{
		NodeID:            nodeID,
		Blockchain:        blockchain,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SNVM:              synnergy_vm.NewVirtualMachine(), // Initialize the virtual machine for smart contract execution
		DataIndex:         make(map[string]interface{}),    // Initialize the data index map
		SubBlocks:         make(map[string]*common.SubBlock),
		SyncInterval:      syncInterval,
	}
}

// StartNode starts the hybrid node’s operations, including syncing, validating transactions, handling data queries, and participating in the consensus.
func (hn *HybridNode) StartNode() error {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Begin syncing with other nodes.
	go hn.syncWithOtherNodes()

	// Listen for incoming transactions and process them.
	go hn.listenForTransactions()

	// Set up block proposal and validation participation.
	go hn.setupConsensusParticipation()

	fmt.Printf("Hybrid Node %s started successfully.\n", hn.NodeID)
	return nil
}

// syncWithOtherNodes handles syncing the blockchain with other nodes at regular intervals.
func (hn *HybridNode) syncWithOtherNodes() {
	ticker := time.NewTicker(hn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		hn.mutex.Lock()
		// Discover other nodes and request blockchain data.
		otherNodes := hn.NetworkManager.DiscoverOtherNodes(hn.NodeID)
		for _, node := range otherNodes {
			err := hn.syncBlockchainFromNode(node)
			if err != nil {
				fmt.Printf("Failed to sync blockchain from node %s: %v\n", node, err)
			}
		}
		hn.mutex.Unlock()
	}
}

// syncBlockchainFromNode syncs the blockchain state from a peer node.
func (hn *HybridNode) syncBlockchainFromNode(peerNode string) error {
	peerBlockchain, err := hn.NetworkManager.RequestBlockchain(peerNode)
	if err != nil {
		return fmt.Errorf("failed to request blockchain from node %s: %v", peerNode, err)
	}

	// Validate the blockchain and update the local copy if necessary.
	if hn.ConsensusEngine.ValidateBlockchain(peerBlockchain) {
		hn.Blockchain = peerBlockchain
		fmt.Printf("Blockchain synced successfully from node %s.\n", peerNode)
	} else {
		fmt.Printf("Blockchain sync from node %s failed validation.\n", peerNode)
	}
	return nil
}

// listenForTransactions listens for incoming transactions and processes them.
func (hn *HybridNode) listenForTransactions() {
	for {
		transaction, err := hn.NetworkManager.ReceiveTransaction()
		if err != nil {
			fmt.Printf("Error receiving transaction: %v\n", err)
			continue
		}

		// Process and validate the transaction.
		err = hn.processTransaction(transaction)
		if err != nil {
			fmt.Printf("Transaction processing failed: %v\n", err)
		}
	}
}

// processTransaction processes and validates an incoming transaction.
func (hn *HybridNode) processTransaction(tx *ledger.Transaction) error {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Validate the transaction using the consensus engine.
	if valid, err := hn.ConsensusEngine.ValidateTransaction(tx); err != nil || !valid {
		return fmt.Errorf("invalid transaction: %v", err)
	}

	// Add the transaction to a sub-block.
	subBlockID := common.GenerateSubBlockID()
	subBlock := hn.createSubBlock(tx, subBlockID)

	// Store the sub-block and attempt to validate it into a block.
	hn.SubBlocks[subBlockID] = subBlock
	err := hn.tryValidateSubBlock(subBlock)
	if err != nil {
		return fmt.Errorf("failed to validate sub-block: %v", err)
	}

	// Add transaction to data index for query handling.
	hn.DataIndex[tx.TransactionID] = tx

	fmt.Printf("Transaction %s processed and indexed.\n", tx.TransactionID)
	return nil
}

// createSubBlock creates a sub-block from a validated transaction.
func (hn *HybridNode) createSubBlock(tx *ledger.Transaction, subBlockID string) *common.SubBlock {
	return &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: []*ledger.Transaction{tx},
		Timestamp:    time.Now(),
		NodeID:       hn.NodeID,
	}
}

// tryValidateSubBlock attempts to validate a sub-block into a full block.
func (hn *HybridNode) tryValidateSubBlock(subBlock *common.SubBlock) error {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Check if the sub-block is valid using the consensus mechanism.
	if hn.ConsensusEngine.ValidateSubBlock(subBlock) {
		// Add the sub-block to a full block and broadcast it.
		block := hn.Blockchain.AddSubBlock(subBlock)
		err := hn.NetworkManager.BroadcastNewBlock(block)
		if err != nil {
			return fmt.Errorf("failed to broadcast new block: %v", err)
		}

		fmt.Printf("Sub-block %s validated and added to blockchain.\n", subBlock.SubBlockID)
		return nil
	}

	return errors.New("failed to validate sub-block")
}

// setupConsensusParticipation sets up the hybrid node’s participation in the Synnergy Consensus mechanism for block proposal and validation.
func (hn *HybridNode) setupConsensusParticipation() {
	for {
		select {
		case proposal := <-hn.ConsensusEngine.ProposalChannel:
			err := hn.validateAndProposeBlock(proposal)
			if err != nil {
				fmt.Printf("Block proposal validation failed: %v\n", err)
			}
		}
	}
}

// validateAndProposeBlock validates a block proposal and proposes a new block if necessary.
func (hn *HybridNode) validateAndProposeBlock(proposal *common.BlockProposal) error {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Validate the block proposal using the consensus engine.
	if valid, err := hn.ConsensusEngine.ValidateBlockProposal(proposal); err != nil || !valid {
		return fmt.Errorf("invalid block proposal: %v", err)
	}

	// Propose a new block based on the validated proposal.
	newBlock, err := hn.Blockchain.CreateBlockFromProposal(proposal)
	if err != nil {
		return fmt.Errorf("failed to create block from proposal: %v", err)
	}

	// Broadcast the new block to the network.
	err = hn.NetworkManager.BroadcastNewBlock(newBlock)
	if err != nil {
		return fmt.Errorf("failed to broadcast new block: %v", err)
	}

	fmt.Printf("Block proposed and validated by hybrid node %s.\n", hn.NodeID)
	return nil
}

// DeployContract deploys a smart contract to the hybrid node’s virtual machine.
func (hn *HybridNode) DeployContract(contractCode []byte, contractOwner string) (string, error) {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Encrypt the contract code before deployment.
	encryptedCode, err := hn.EncryptionService.EncryptData(contractCode, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt contract code: %v", err)
	}

	// Generate a unique contract ID.
	contractID := common.GenerateUniqueID()

	// Deploy the contract on the virtual machine.
	err = hn.SNVM.DeployContract(contractID, encryptedCode, contractOwner)
	if err != nil {
		return "", fmt.Errorf("failed to deploy contract: %v", err)
	}

	// Record the contract deployment in the ledger.
	err = hn.Blockchain.RecordContractDeployment(contractID, encryptedCode)
	if err != nil {
		return "", fmt.Errorf("failed to record contract deployment in the ledger: %v", err)
	}

	fmt.Printf("Contract %s deployed successfully on hybrid node %s.\n", contractID, hn.NodeID)
	return contractID, nil
}

// ExecuteContract executes a smart contract on the hybrid node’s virtual machine.
func (hn *HybridNode) ExecuteContract(contractID string, args []byte) ([]byte, error) {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Encrypt the function arguments before execution.
	encryptedArgs, err := hn.EncryptionService.EncryptData(args, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt contract arguments: %v", err)
	}

	// Execute the contract on the virtual machine.
	result, err := hn.SNVM.ExecuteContract(contractID, encryptedArgs)
	if err != nil {
		return nil, fmt.Errorf("contract execution failed: %v", err)
	}

	// Record the execution result in the ledger.
	err = hn.Blockchain.RecordContractExecution(contractID, result)
	if err != nil {
		return nil, fmt.Errorf("failed to record contract execution in ledger: %v", err)
	}

	fmt.Printf("Contract %s executed successfully on hybrid node %s.\n", contractID, hn.NodeID)
	return result, nil
}
