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

// IndexingNode represents a node responsible for indexing blockchain data and handling queries efficiently.
type IndexingNode struct {
	NodeID            string                        // Unique identifier for the node
	Blockchain        *ledger.Blockchain            // Local copy of the blockchain ledger
	ConsensusEngine   *common.SynnergyConsensus     // Consensus engine for validating transactions and sub-blocks
	EncryptionService *common.Encryption            // Encryption service for secure data handling
	NetworkManager    *network.NetworkManager       // Network manager for communication with other nodes
	SNVM              *common.VMInterface   // The Synnergy Network Virtual Machine
	Index             map[string]interface{}        // Index map for efficient data retrieval
	SubBlocks         map[string]*common.SubBlock   // Sub-block storage before validation into blocks
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with other full nodes
	QueryCache        map[string]interface{}        // Cache for frequently queried data
}

// NewIndexingNode initializes a new indexing node in the blockchain network.
func NewIndexingNode(nodeID string, blockchain *ledger.Blockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *network.NetworkManager, syncInterval time.Duration) *IndexingNode {
	return &IndexingNode{
		NodeID:            nodeID,
		Blockchain:        blockchain,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SNVM:              synnergy_vm.NewVirtualMachine(), // Initialize the virtual machine for smart contract execution
		Index:             make(map[string]interface{}),    // Initialize the index for blockchain data
		SubBlocks:         make(map[string]*common.SubBlock),
		SyncInterval:      syncInterval,
		QueryCache:        make(map[string]interface{}),    // Initialize cache for frequently queried data
	}
}

// StartNode starts the indexing node's operations, including syncing, handling transactions, and processing queries.
func (in *IndexingNode) StartNode() error {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Begin syncing the blockchain with other nodes.
	go in.syncWithOtherNodes()

	// Listen for incoming transactions and process them.
	go in.listenForTransactions()

	// Start handling queries from the network.
	go in.handleQueries()

	fmt.Printf("Indexing Node %s started successfully.\n", in.NodeID)
	return nil
}

// syncWithOtherNodes handles syncing the blockchain with other full nodes at regular intervals.
func (in *IndexingNode) syncWithOtherNodes() {
	ticker := time.NewTicker(in.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		in.mutex.Lock()
		// Discover other nodes and sync blockchain data.
		otherNodes := in.NetworkManager.DiscoverOtherNodes(in.NodeID)
		for _, node := range otherNodes {
			err := in.syncBlockchainFromNode(node)
			if err != nil {
				fmt.Printf("Failed to sync blockchain from node %s: %v\n", node, err)
			}
		}
		in.mutex.Unlock()
	}
}

// syncBlockchainFromNode syncs the blockchain state from a peer node.
func (in *IndexingNode) syncBlockchainFromNode(peerNode string) error {
	peerBlockchain, err := in.NetworkManager.RequestBlockchain(peerNode)
	if err != nil {
		return fmt.Errorf("failed to request blockchain from node %s: %v", peerNode, err)
	}

	// Validate the blockchain and update the local copy if necessary.
	if in.ConsensusEngine.ValidateBlockchain(peerBlockchain) {
		in.Blockchain = peerBlockchain
		fmt.Printf("Blockchain synced successfully from node %s.\n", peerNode)
	} else {
		fmt.Printf("Blockchain sync from node %s failed validation.\n", peerNode)
	}
	return nil
}

// listenForTransactions listens for incoming transactions and processes them.
func (in *IndexingNode) listenForTransactions() {
	for {
		transaction, err := in.NetworkManager.ReceiveTransaction()
		if err != nil {
			fmt.Printf("Error receiving transaction: %v\n", err)
			continue
		}

		// Process and validate the transaction.
		err = in.processTransaction(transaction)
		if err != nil {
			fmt.Printf("Transaction processing failed: %v\n", err)
		}
	}
}

// processTransaction processes and validates an incoming transaction.
func (in *IndexingNode) processTransaction(tx *ledger.Transaction) error {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Validate the transaction using the consensus engine.
	if valid, err := in.ConsensusEngine.ValidateTransaction(tx); err != nil || !valid {
		return fmt.Errorf("invalid transaction: %v", err)
	}

	// Add the transaction to a sub-block.
	subBlockID := common.GenerateSubBlockID()
	subBlock := in.createSubBlock(tx, subBlockID)

	// Store the sub-block and attempt to validate it into a block.
	in.SubBlocks[subBlockID] = subBlock
	err := in.tryValidateSubBlock(subBlock)
	if err != nil {
		return fmt.Errorf("failed to validate sub-block: %v", err)
	}

	// Add the transaction to the index for efficient querying.
	in.Index[tx.TransactionID] = tx
	fmt.Printf("Transaction %s processed and indexed.\n", tx.TransactionID)
	return nil
}

// createSubBlock creates a sub-block from a validated transaction.
func (in *IndexingNode) createSubBlock(tx *ledger.Transaction, subBlockID string) *common.SubBlock {
	return &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: []*ledger.Transaction{tx},
		Timestamp:    time.Now(),
		NodeID:       in.NodeID,
	}
}

// tryValidateSubBlock attempts to validate a sub-block into a full block.
func (in *IndexingNode) tryValidateSubBlock(subBlock *common.SubBlock) error {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Check if the sub-block is valid using the consensus mechanism.
	if in.ConsensusEngine.ValidateSubBlock(subBlock) {
		// Add the sub-block to a full block and broadcast it.
		block := in.Blockchain.AddSubBlock(subBlock)
		err := in.NetworkManager.BroadcastNewBlock(block)
		if err != nil {
			return fmt.Errorf("failed to broadcast new block: %v", err)
		}

		fmt.Printf("Sub-block %s validated and added to blockchain.\n", subBlock.SubBlockID)
		return nil
	}

	return errors.New("failed to validate sub-block")
}

// handleQueries listens for and processes queries for indexed data.
func (in *IndexingNode) handleQueries() {
	for {
		query, err := in.NetworkManager.ReceiveQuery()
		if err != nil {
			fmt.Printf("Error receiving query: %v\n", err)
			continue
		}

		// Process the query and respond with indexed data.
		err = in.processQuery(query)
		if err != nil {
			fmt.Printf("Query processing failed: %v\n", err)
		}
	}
}

// processQuery handles a query by searching the index and returning the requested data.
func (in *IndexingNode) processQuery(query *common.Query) error {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Check the query cache first for faster results.
	if result, found := in.QueryCache[query.QueryID]; found {
		in.NetworkManager.SendQueryResult(query, result)
		return nil
	}

	// If not cached, search the index for the relevant data.
	result, exists := in.Index[query.QueryID]
	if !exists {
		return fmt.Errorf("query result not found")
	}

	// Cache the result for future queries.
	in.QueryCache[query.QueryID] = result

	// Return the result.
	in.NetworkManager.SendQueryResult(query, result)
	fmt.Printf("Query %s processed successfully.\n", query.QueryID)
	return nil
}

// DeployContract deploys a smart contract to the indexing node's virtual machine.
func (in *IndexingNode) DeployContract(contractCode []byte, contractOwner string) (string, error) {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Encrypt the contract code before deployment.
	encryptedCode, err := in.EncryptionService.EncryptData(contractCode, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt contract code: %v", err)
	}

	// Generate a unique contract ID.
	contractID := common.GenerateUniqueID()

	// Deploy the contract on the virtual machine.
	err = in.SNVM.DeployContract(contractID, encryptedCode, contractOwner)
	if err != nil {
		return "", fmt.Errorf("failed to deploy contract: %v", err)
	}

	// Record the contract deployment in the ledger.
	err = in.Blockchain.RecordContractDeployment(contractID, contractOwner, encryptedCode)
	if err != nil {
		return "", fmt.Errorf("failed to record contract deployment: %v", err)
	}

	fmt.Printf("Contract %s deployed successfully on indexing node %s.\n", contractID, in.NodeID)
	return contractID, nil
}

// ExecuteContract executes a smart contract on the indexing node’s virtual machine.
func (in *IndexingNode) ExecuteContract(contractID string, args []byte) ([]byte, error) {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Encrypt the contract arguments.
	encryptedArgs, err := in.EncryptionService.EncryptData(args, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt contract arguments: %v", err)
	}

	// Execute the contract on the virtual machine.
	result, err := in.SNVM.ExecuteContract(contractID, encryptedArgs)
	if err != nil {
		return nil, fmt.Errorf("contract execution failed: %v", err)
	}

	// Record the execution result in the ledger.
	err = in.Blockchain.RecordContractExecution(contractID, result)
	if err != nil {
		return nil, fmt.Errorf("failed to record contract execution: %v", err)
	}

	fmt.Printf("Contract %s executed successfully on indexing node %s.\n", contractID, in.NodeID)
	return result, nil
}
