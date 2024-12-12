package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"     // Shared components like encryption, consensus, sub-blocks, etc.
	"synnergy_network/pkg/ledger"     // Blockchain and ledger-related components
)

// ArchivalFullNode represents an archival full node in the blockchain network.
type ArchivalFullNode struct {
	NodeID            string                        // Unique identifier for the node
	CompleteBlockchain *ledger.Blockchain           // Local complete copy of the blockchain ledger from genesis
	ConsensusEngine   *common.SynnergyConsensus     // Consensus engine for validating transactions and smart contracts
	EncryptionService *common.Encryption            // Encryption service for secure communication
	NetworkManager    *common.NetworkManager        // Network manager for communicating with other nodes
	SubBlocks         map[string]*common.SubBlock   // Sub-blocks that are part of blocks in the blockchain
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with other full nodes
	SNVM              *common.SynnergyVirtualMachine // Virtual Machine for smart contract execution
}

// NewArchivalFullNode initializes a new archival full node in the blockchain network.
func NewArchivalFullNode(nodeID string, completeBlockchain *ledger.Blockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *common.NetworkManager, syncInterval time.Duration) *ArchivalFullNode {
	return &ArchivalFullNode{
		NodeID:             nodeID,
		CompleteBlockchain: completeBlockchain,      // Store the full blockchain, no pruning
		ConsensusEngine:    consensusEngine,
		EncryptionService:  encryptionService,
		NetworkManager:     networkManager,
		SubBlocks:          make(map[string]*common.SubBlock),
		SyncInterval:       syncInterval,
		SNVM:               common.NewSynnergyVirtualMachine(), // Initialize the Virtual Machine for contract execution
	}
}

// StartNode starts the archival full node's operations, including syncing, validating transactions, and running smart contracts.
func (afn *ArchivalFullNode) StartNode() error {
	afn.mutex.Lock()
	defer afn.mutex.Unlock()

	// Start syncing with other nodes and listening for transactions.
	go afn.syncWithOtherNodes()
	go afn.listenForTransactions()

	fmt.Printf("Archival full node %s started successfully.\n", afn.NodeID)
	return nil
}

// syncWithOtherNodes handles syncing the complete blockchain with other archival or full nodes at regular intervals.
func (afn *ArchivalFullNode) syncWithOtherNodes() {
	ticker := time.NewTicker(afn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		afn.mutex.Lock()
		otherNodes := afn.NetworkManager.DiscoverOtherFullNodes(afn.NodeID)
		for _, node := range otherNodes {
			// Sync blockchain from each node.
			afn.syncBlockchainFromNode(node)
		}
		afn.mutex.Unlock()
	}
}

// syncBlockchainFromNode syncs the full blockchain from a peer node.
func (afn *ArchivalFullNode) syncBlockchainFromNode(peerNode string) {
	peerBlockchain, err := afn.NetworkManager.RequestBlockchain(peerNode)
	if err != nil {
		fmt.Printf("Failed to sync blockchain from node %s: %v\n", peerNode, err)
		return
	}

	if afn.ConsensusEngine.ValidateBlockchain(peerBlockchain) {
		afn.CompleteBlockchain = afn.CompleteBlockchain.MergeWith(peerBlockchain)
		fmt.Printf("Complete blockchain synced successfully from node %s.\n", peerNode)
	} else {
		fmt.Printf("Blockchain sync from node %s failed validation.\n", peerNode)
	}
}

// listenForTransactions listens for incoming transactions and processes them.
func (afn *ArchivalFullNode) listenForTransactions() {
	for {
		transaction, err := afn.NetworkManager.ReceiveTransaction()
		if err != nil {
			fmt.Printf("Error receiving transaction: %v\n", err)
			continue
		}

		// Validate and process the transaction.
		err = afn.processTransaction(transaction)
		if err != nil {
			fmt.Printf("Transaction processing failed: %v\n", err)
		}
	}
}

// processTransaction processes and validates an incoming transaction, adding it to the blockchain if valid.
func (afn *ArchivalFullNode) processTransaction(tx *ledger.Transaction) error {
	afn.mutex.Lock()
	defer afn.mutex.Unlock()

	// Validate the transaction using the consensus engine.
	if valid, err := afn.ConsensusEngine.ValidateTransaction(tx); err != nil || !valid {
		return fmt.Errorf("invalid transaction: %v", err)
	}

	// Add the transaction to a sub-block.
	subBlockID := common.GenerateSubBlockID()
	subBlock := afn.createSubBlock(tx, subBlockID)

	// Store the sub-block and attempt to validate it into a block.
	afn.SubBlocks[subBlockID] = subBlock
	return afn.tryValidateSubBlock(subBlock)
}

// createSubBlock creates a sub-block from a validated transaction.
func (afn *ArchivalFullNode) createSubBlock(tx *ledger.Transaction, subBlockID string) *common.SubBlock {
	return &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: []*ledger.Transaction{tx},
		Timestamp:    time.Now(),
		NodeID:       afn.NodeID,
	}
}

// tryValidateSubBlock attempts to validate a sub-block and add it to the complete blockchain.
func (afn *ArchivalFullNode) tryValidateSubBlock(subBlock *common.SubBlock) error {
	afn.mutex.Lock()
	defer afn.mutex.Unlock()

	// Check if the sub-block is valid using the consensus mechanism.
	if afn.ConsensusEngine.ValidateSubBlock(subBlock) {
		// Add sub-block to the blockchain as part of a full block.
		block := afn.CompleteBlockchain.AddSubBlock(subBlock)

		// Notify the network of the new block.
		err := afn.NetworkManager.BroadcastNewBlock(block)
		if err != nil {
			return fmt.Errorf("failed to broadcast new block: %v", err)
		}

		fmt.Printf("Sub-block %s validated and added to blockchain.\n", subBlock.SubBlockID)
		return nil
	}

	return errors.New("failed to validate sub-block")
}

// Virtual Machine (VM) Smart Contract Functions

// DeployContract deploys a smart contract to the archival full node's virtual machine.
func (afn *ArchivalFullNode) DeployContract(contractCode []byte, contractOwner string) (string, error) {
	afn.mutex.Lock()
	defer afn.mutex.Unlock()

	// Encrypt the contract code before deployment.
	encryptedCode, err := afn.EncryptionService.EncryptData(contractCode, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt contract code: %v", err)
	}

	// Generate a unique contract ID.
	contractID := common.GenerateUniqueID()

	// Deploy the contract on the virtual machine.
	err = afn.SNVM.DeployContract(contractID, encryptedCode, contractOwner)
	if err != nil {
		return "", fmt.Errorf("failed to deploy contract: %v", err)
	}

	// Record the contract deployment in the ledger.
	err = afn.CompleteBlockchain.RecordContractDeployment(contractID, contractOwner, encryptedCode)
	if err != nil {
		return "", fmt.Errorf("failed to record contract deployment in ledger: %v", err)
	}

	fmt.Printf("Contract %s deployed successfully on archival node %s by owner %s.\n", contractID, afn.NodeID, contractOwner)
	return contractID, nil
}

// ExecuteContract executes a function of a deployed contract on the archival full node's virtual machine.
func (afn *ArchivalFullNode) ExecuteContract(contractID string, functionName string, args []byte) ([]byte, error) {
	afn.mutex.Lock()
	defer afn.mutex.Unlock()

	// Encrypt the function arguments.
	encryptedArgs, err := afn.EncryptionService.EncryptData(args, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt function arguments: %v", err)
	}

	// Execute the contract function using the virtual machine.
	result, err := afn.SNVM.ExecuteContractFunction(contractID, functionName, encryptedArgs)
	if err != nil {
		return nil, fmt.Errorf("contract execution failed: %v", err)
	}

	// Decrypt the execution result.
	decryptedResult, err := afn.EncryptionService.DecryptData(result, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt contract result: %v", err)
	}

	// Record the contract execution in the ledger.
	err = afn.CompleteBlockchain.RecordContractExecution(contractID, functionName, args, decryptedResult)
	if err != nil {
		return nil, fmt.Errorf("failed to record contract execution in ledger: %v", err)
	}

	fmt.Printf("Contract %s executed function %s on archival node %s.\n", contractID, functionName, afn.NodeID)
	return decryptedResult, nil
}

// HandleIncomingBlock handles a new block received from the network.
func (afn *ArchivalFullNode) HandleIncomingBlock(block *ledger.Block) error {
	afn.mutex.Lock()
	defer afn.mutex.Unlock()

	// Validate the incoming block using the consensus engine.
	if afn.ConsensusEngine.ValidateBlock(block) {
		// Add the validated block to the blockchain.
		err := afn.CompleteBlockchain.AddBlock(block)
		if err != nil {
			return fmt.Errorf("failed to add block to blockchain: %v", err)
		}

		fmt.Printf("Block %s added to blockchain.\n", block.BlockID)
		return nil
	}

	return errors.New("failed to validate incoming block")
}

// Historical Data Retrieval

// ProvideHistoricalData allows querying the blockchain for historical data.
func (afn *ArchivalFullNode) ProvideHistoricalData(query common.HistoricalQuery) ([]common.BlockData, error) {
	afn.mutex.Lock()
	defer afn.mutex.Unlock()

	// Fetch historical data based on the query parameters.
	history, err := afn.CompleteBlockchain.FetchHistoricalData(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve historical data: %v", err)
	}

	return history, nil
}
