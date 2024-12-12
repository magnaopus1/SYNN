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
)

// TestnetNode represents a node designed specifically to run in the testnet environment. It validates test transactions and sub-blocks exclusively in the testnet.
type TestnetNode struct {
	NodeID            string                        // Unique identifier for the testnet node
	PartialLedger     map[string]*ledger.Block      // Stores a subset of blocks for testing purposes (no mainnet data)
	TestBlocks        map[string]*ledger.SubBlock   // Stores validated sub-blocks in the testnet
	ConsensusEngine   *synnergy_consensus.Engine    // Consensus engine for validating test transactions
	EncryptionService *encryption.Encryption        // Encryption service for securing testnet transactions
	NetworkManager    *network.NetworkManager       // Network manager for communicating with other testnet nodes
	Ledger            *ledger.Ledger                // Reference to the testnet ledger for logging test transactions
	SyncInterval      time.Duration                 // Interval for syncing with other testnet nodes
	FullNodes         []string                      // List of full testnet nodes to interact with
	mutex             sync.Mutex                    // Mutex for thread-safe operations
}

// NewTestnetNode initializes a new testnet node in the testnet environment.
func NewTestnetNode(nodeID string, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, ledgerInstance *ledger.Ledger, syncInterval time.Duration, fullNodes []string) *TestnetNode {
	return &TestnetNode{
		NodeID:            nodeID,
		PartialLedger:     make(map[string]*ledger.Block),
		TestBlocks:        make(map[string]*ledger.SubBlock),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		Ledger:            ledgerInstance,
		SyncInterval:      syncInterval,
		FullNodes:         fullNodes,
	}
}

// StartNode starts the testnet node's operations, including syncing with other testnet nodes and listening for test transaction requests.
func (tn *TestnetNode) StartNode() error {
	tn.mutex.Lock()
	defer tn.mutex.Unlock()

	// Start syncing with full testnet nodes.
	go tn.syncWithTestnetNodes()

	// Listen for incoming transactions and sub-block validation requests.
	go tn.listenForTestTransactions()

	fmt.Printf("Testnet node %s started successfully.\n", tn.NodeID)
	return nil
}

// syncWithTestnetNodes syncs the partial ledger with full testnet nodes at regular intervals.
func (tn *TestnetNode) syncWithTestnetNodes() {
	ticker := time.NewTicker(tn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		tn.mutex.Lock()
		for _, fullNode := range tn.FullNodes {
			// Request recent blocks from the testnet full node.
			tn.syncPartialLedgerFromTestnetFullNode(fullNode)
		}
		tn.mutex.Unlock()
	}
}

// syncPartialLedgerFromTestnetFullNode requests a subset of the testnet blockchain from a full testnet node.
func (tn *TestnetNode) syncPartialLedgerFromTestnetFullNode(fullNodeID string) {
	partialBlocks, err := tn.NetworkManager.RequestPartialBlockchain(fullNodeID)
	if err != nil {
		fmt.Printf("Failed to sync partial testnet blockchain from node %s: %v\n", fullNodeID, err)
		return
	}

	// Validate and store the blocks in the partial ledger.
	for _, block := range partialBlocks {
		if tn.ConsensusEngine.ValidateBlock(block) {
			tn.PartialLedger[block.BlockID] = block
			fmt.Printf("Block %s synced and added to testnet partial ledger from node %s.\n", block.BlockID, fullNodeID)
		} else {
			fmt.Printf("Block %s from testnet node %s failed validation.\n", block.BlockID, fullNodeID)
		}
	}
}

// listenForTestTransactions listens for incoming test transaction requests and sub-block validation requests.
func (tn *TestnetNode) listenForTestTransactions() {
	for {
		request, err := tn.NetworkManager.ReceiveTransactionRequest()
		if err != nil {
			fmt.Printf("Error receiving test transaction request: %v\n", err)
			continue
		}

		// Process the transaction request.
		err = tn.processTestTransactionRequest(request)
		if err != nil {
			fmt.Printf("Test transaction request processing failed: %v\n", err)
		}
	}
}

// processTestTransactionRequest processes a test transaction validation request received from the network.
func (tn *TestnetNode) processTestTransactionRequest(request *network.TransactionRequest) error {
	tn.mutex.Lock()
	defer tn.mutex.Unlock()

	// Validate the transaction using the testnet consensus engine.
	valid, err := tn.ConsensusEngine.ValidateTransaction(request.Transaction)
	if err != nil || !valid {
		return fmt.Errorf("test transaction validation failed: %v", err)
	}

	// Store the transaction in a new test sub-block.
	subBlock := &ledger.SubBlock{
		SubBlockID:    common.GenerateUniqueID(),
		Transactions:  []*ledger.Transaction{request.Transaction},
		CreationTime:  time.Now(),
	}

	// Validate and store the sub-block in the testnet.
	err = tn.storeTestSubBlock(subBlock)
	if err != nil {
		return fmt.Errorf("failed to store test sub-block: %v", err)
	}

	fmt.Printf("Test transaction %s validated and stored in test sub-block %s.\n", request.Transaction.TransactionID, subBlock.SubBlockID)
	return nil
}

// storeTestSubBlock stores a validated test sub-block in the testnet ledger.
func (tn *TestnetNode) storeTestSubBlock(subBlock *ledger.SubBlock) error {
	// Store the test sub-block in memory.
	tn.TestBlocks[subBlock.SubBlockID] = subBlock

	// Log the sub-block in the testnet ledger.
	err := tn.Ledger.RecordSubBlock(subBlock)
	if err != nil {
		return fmt.Errorf("failed to log test sub-block in ledger: %v", err)
	}

	fmt.Printf("Test sub-block %s stored successfully by testnet node %s.\n", subBlock.SubBlockID, tn.NodeID)
	return nil
}

// EncryptTestData encrypts sensitive data before sending it to other testnet nodes.
func (tn *TestnetNode) EncryptTestData(data []byte) ([]byte, error) {
	encryptedData, err := tn.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt test data: %v", err)
	}
	return encryptedData, nil
}

// DecryptTestData decrypts incoming encrypted data from other testnet nodes.
func (tn *TestnetNode) DecryptTestData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := tn.EncryptionService.DecryptData(encryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt test data: %v", err)
	}
	return decryptedData, nil
}

// RequestTestBlock queries a testnet full node for a specific test block.
func (tn *TestnetNode) RequestTestBlock(blockID string) (*ledger.Block, error) {
	tn.mutex.Lock()
	defer tn.mutex.Unlock()

	// Check if the block is already in the partial ledger.
	block, exists := tn.PartialLedger[blockID]
	if exists {
		return block, nil
	}

	// Request the block from a full testnet node.
	fullNodeID := tn.selectRandomFullNode()
	block, err := tn.NetworkManager.RequestBlockFromFullNode(fullNodeID, blockID)
	if err != nil {
		return nil, fmt.Errorf("failed to request test block from full node: %v", err)
	}

	// Validate and store the block if valid.
	if tn.ConsensusEngine.ValidateBlock(block) {
		tn.PartialLedger[blockID] = block
		return block, nil
	}

	return nil, errors.New("failed to validate test block from full node")
}

// selectRandomFullNode selects a random full testnet node to forward transactions or requests to.
func (tn *TestnetNode) selectRandomFullNode() string {
	tn.mutex.Lock()
	defer tn.mutex.Unlock()

	if len(tn.FullNodes) == 0 {
		fmt.Println("No full testnet nodes available to forward transactions.")
		return ""
	}

	randomIndex := common.GenerateRandomInt(len(tn.FullNodes))
	return tn.FullNodes[randomIndex]
}
