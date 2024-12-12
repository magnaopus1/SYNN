package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/ledger"           // Blockchain ledger components
	"synnergy_network_demo/encryption"       // Encryption service for secure data handling
	"synnergy_network_demo/synnergy_consensus" // Synnergy Consensus engine
	"synnergy_network_demo/network"          // Network management
	"synnergy_network_demo/common"           // Common utilities for random generation, keys, etc.
)

// MobileNode represents a node designed to run on mobile devices (iOS, Android, etc.), providing a lightweight interface to the blockchain.
type MobileNode struct {
	NodeID            string                        // Unique identifier for the node
	PartialLedger     map[string]*ledger.Block      // Stores a subset of blocks to reduce storage usage
	ConsensusEngine   *synnergy_consensus.Engine    // Consensus engine for validating transactions and syncing with full nodes
	EncryptionService *encryption.Encryption        // Encryption service for securing transactions and data
	NetworkManager    *network.NetworkManager       // Network manager for communicating with other nodes
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with full nodes
	FullNodes         []string                      // List of full nodes for syncing and querying
	SNVM              *common.VMInterface   // The Synnergy Network Virtual Machine
	BatteryOptimized  bool                          // Flag to determine whether battery-saving measures are active
}

// NewMobileNode initializes a new mobile node for a mobile device in the blockchain network.
func NewMobileNode(nodeID string, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, syncInterval time.Duration, fullNodes []string, batteryOptimized bool) *MobileNode {
	return &MobileNode{
		NodeID:            nodeID,
		PartialLedger:     make(map[string]*ledger.Block),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SyncInterval:      syncInterval,
		FullNodes:         fullNodes,
		BatteryOptimized:  batteryOptimized,
	}
}

// StartNode begins the operations of the Mobile Node, including syncing with full nodes and processing transactions.
func (mn *MobileNode) StartNode() error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	// Sync with the full nodes periodically based on mobile capabilities.
	go mn.syncWithFullNodes()

	// Listen for incoming transaction requests or data syncing needs.
	go mn.listenForTransactionRequests()

	fmt.Printf("Mobile node %s started successfully.\n", mn.NodeID)
	return nil
}

// syncWithFullNodes handles syncing the partial blockchain with full nodes at regular intervals, optimized for mobile devices.
func (mn *MobileNode) syncWithFullNodes() {
	ticker := time.NewTicker(mn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		mn.mutex.Lock()
		for _, fullNode := range mn.FullNodes {
			// Sync recent blocks from the full node to the partial ledger.
			mn.syncPartialLedgerFromFullNode(fullNode)
		}
		mn.mutex.Unlock()
	}
}

// syncPartialLedgerFromFullNode requests recent blockchain data (block headers) from a full node.
func (mn *MobileNode) syncPartialLedgerFromFullNode(fullNodeID string) {
	partialBlockchain, err := mn.NetworkManager.RequestPartialBlockchain(fullNodeID)
	if err != nil {
		fmt.Printf("Failed to sync partial blockchain from full node %s: %v\n", fullNodeID, err)
		return
	}

	// Validate and store the partial ledger.
	for _, block := range partialBlockchain {
		if mn.ConsensusEngine.ValidateBlock(block) {
			mn.PartialLedger[block.BlockID] = block
			fmt.Printf("Block %s synced and added to partial ledger from full node %s.\n", block.BlockID, fullNodeID)
		} else {
			fmt.Printf("Block %s from full node %s failed validation.\n", block.BlockID, fullNodeID)
		}
	}
}

// listenForTransactionRequests listens for incoming transaction validation requests from other nodes or mobile apps.
func (mn *MobileNode) listenForTransactionRequests() {
	for {
		request, err := mn.NetworkManager.ReceiveTransactionRequest()
		if err != nil {
			fmt.Printf("Error receiving transaction request: %v\n", err)
			continue
		}

		// Process the transaction request.
		err = mn.processTransactionRequest(request)
		if err != nil {
			fmt.Printf("Transaction request processing failed: %v\n", err)
		}
	}
}

// processTransactionRequest processes a transaction validation request received from the network or mobile apps.
func (mn *MobileNode) processTransactionRequest(request *network.TransactionRequest) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	// Validate the transaction using the consensus engine.
	valid, err := mn.ConsensusEngine.ValidateTransaction(request.Transaction)
	if err != nil || !valid {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Forward the validated transaction to a full node for inclusion in a block.
	fullNodeID := mn.selectRandomFullNode()
	err = mn.NetworkManager.ForwardTransactionToFullNode(request.Transaction, fullNodeID)
	if err != nil {
		return fmt.Errorf("failed to forward transaction to full node: %v", err)
	}

	fmt.Printf("Transaction %s validated and forwarded to full node %s.\n", request.Transaction.TransactionID, fullNodeID)
	return nil
}

// selectRandomFullNode selects a random full node to forward transactions or requests to, optimizing for mobile constraints.
func (mn *MobileNode) selectRandomFullNode() string {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if len(mn.FullNodes) == 0 {
		fmt.Println("No full nodes available to forward transactions.")
		return ""
	}

	randomIndex := common.GenerateRandomInt(len(mn.FullNodes))
	return mn.FullNodes[randomIndex]
}

// RequestBlock queries a full node for a specific block, optimizing bandwidth usage and performance for mobile devices.
func (mn *MobileNode) RequestBlock(blockID string) (*ledger.Block, error) {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	// Check if the block is already in the partial ledger.
	block, exists := mn.PartialLedger[blockID]
	if exists {
		return block, nil
	}

	// Request the block from a full node.
	fullNodeID := mn.selectRandomFullNode()
	block, err := mn.NetworkManager.RequestBlockFromFullNode(fullNodeID, blockID)
	if err != nil {
		return nil, fmt.Errorf("failed to request block from full node: %v", err)
	}

	// Validate and store the block if valid.
	if mn.ConsensusEngine.ValidateBlock(block) {
		mn.PartialLedger[blockID] = block
		return block, nil
	}

	return nil, errors.New("failed to validate block from full node")
}

// EncryptData encrypts sensitive data before sending it across the network from a mobile device.
func (mn *MobileNode) EncryptData(data []byte) ([]byte, error) {
	encryptedData, err := mn.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return encryptedData, nil
}

// DecryptData decrypts incoming encrypted data from the network, optimizing for mobile device capabilities.
func (mn *MobileNode) DecryptData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := mn.EncryptionService.DecryptData(encryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decryptedData, nil
}

// BatteryOptimization enables or disables battery-saving measures for mobile operations.
func (mn *MobileNode) BatteryOptimization(enable bool) {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()
	mn.BatteryOptimized = enable
	if enable {
		fmt.Println("Battery optimization enabled. Adjusting sync intervals and resource usage.")
		mn.SyncInterval = time.Minute * 10 // Longer sync intervals to save battery
	} else {
		fmt.Println("Battery optimization disabled. Running full operations.")
		mn.SyncInterval = time.Minute * 1 // Regular sync intervals for optimal operations
	}
}
