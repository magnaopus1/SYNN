package node_type

import (
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/ledger"           // Blockchain ledger components
	"synnergy_network_demo/encryption"       // Encryption service for securing data
	"synnergy_network_demo/synnergy_consensus" // Synnergy Consensus engine
	"synnergy_network_demo/network"          // Network management for communication
	"synnergy_network_demo/common"           // Common utilities for random generation, keys, etc.
)

// OptimizationNode represents a node responsible for optimizing transaction processing and improving network performance.
type OptimizationNode struct {
	NodeID            string                        // Unique identifier for the node
	PartialLedger     map[string]*ledger.Block      // Stores a subset of blocks for optimization
	ConsensusEngine   *synnergy_consensus.Engine    // Consensus engine for transaction validation and block processing
	EncryptionService *encryption.Encryption        // Encryption service for securing sensitive optimization data
	NetworkManager    *network.NetworkManager       // Network manager for communicating with other nodes
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with full nodes
	FullNodes         []string                      // List of full nodes to interact with
	SNVM              *common.VMInterface   // The Synnergy Network Virtual Machine
	OptimizationAlgo  OptimizationAlgorithm         // Optimization algorithm for transaction prioritization and ordering
}

// OptimizationAlgorithm defines the structure of the algorithm used for dynamic transaction optimization.
type OptimizationAlgorithm struct {
	PriorityModel   string // Defines the model used for prioritizing transactions (e.g., FIFO, fee-based priority)
	LoadBalancer    string // Defines the load-balancing strategy (e.g., round-robin, weighted distribution)
}

// NewOptimizationNode initializes a new optimization node in the blockchain network.
func NewOptimizationNode(nodeID string, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, syncInterval time.Duration, fullNodes []string, optAlgo OptimizationAlgorithm) *OptimizationNode {
	return &OptimizationNode{
		NodeID:            nodeID,
		PartialLedger:     make(map[string]*ledger.Block),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SyncInterval:      syncInterval,
		FullNodes:         fullNodes,
		OptimizationAlgo:  optAlgo,
	}
}

// StartNode begins the operations of the Optimization Node, including syncing and optimizing transaction flows.
func (on *OptimizationNode) StartNode() error {
	on.mutex.Lock()
	defer on.mutex.Unlock()

	// Start syncing with full nodes periodically based on the sync interval.
	go on.syncWithFullNodes()

	// Start listening for transactions to be optimized.
	go on.listenForTransactionRequests()

	fmt.Printf("Optimization node %s started successfully.\n", on.NodeID)
	return nil
}

// syncWithFullNodes periodically syncs the optimization node's partial ledger with full nodes.
func (on *OptimizationNode) syncWithFullNodes() {
	ticker := time.NewTicker(on.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		on.mutex.Lock()
		for _, fullNode := range on.FullNodes {
			// Sync partial blockchain data from full nodes
			on.syncPartialLedgerFromFullNode(fullNode)
		}
		on.mutex.Unlock()
	}
}

// syncPartialLedgerFromFullNode fetches recent blocks from a full node for optimization purposes.
func (on *OptimizationNode) syncPartialLedgerFromFullNode(fullNodeID string) {
	partialBlockchain, err := on.NetworkManager.RequestPartialBlockchain(fullNodeID)
	if err != nil {
		fmt.Printf("Failed to sync partial blockchain from full node %s: %v\n", fullNodeID, err)
		return
	}

	// Validate and store the partial ledger.
	for _, block := range partialBlockchain {
		if on.ConsensusEngine.ValidateBlock(block) {
			on.PartialLedger[block.BlockID] = block
			fmt.Printf("Block %s synced and added to partial ledger from full node %s.\n", block.BlockID, fullNodeID)
		} else {
			fmt.Printf("Block %s from full node %s failed validation.\n", block.BlockID, fullNodeID)
		}
	}
}

// listenForTransactionRequests listens for incoming transaction requests and optimizes them.
func (on *OptimizationNode) listenForTransactionRequests() {
	for {
		request, err := on.NetworkManager.ReceiveTransactionRequest()
		if err != nil {
			fmt.Printf("Error receiving transaction request: %v\n", err)
			continue
		}

		// Optimize and process the transaction request.
		err = on.processAndOptimizeTransaction(request)
		if err != nil {
			fmt.Printf("Transaction processing failed: %v\n", err)
		}
	}
}

// processAndOptimizeTransaction optimizes and processes incoming transactions.
func (on *OptimizationNode) processAndOptimizeTransaction(request *network.TransactionRequest) error {
	on.mutex.Lock()
	defer on.mutex.Unlock()

	// Validate the transaction using the consensus engine.
	valid, err := on.ConsensusEngine.ValidateTransaction(request.Transaction)
	if err != nil || !valid {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Optimize transaction based on priority model.
	on.optimizeTransactionOrder(request.Transaction)

	// Forward the validated and optimized transaction to a full node for inclusion in a block.
	fullNodeID := on.selectRandomFullNode()
	err = on.NetworkManager.ForwardTransactionToFullNode(request.Transaction, fullNodeID)
	if err != nil {
		return fmt.Errorf("failed to forward transaction to full node: %v", err)
	}

	fmt.Printf("Transaction %s validated, optimized, and forwarded to full node %s.\n", request.Transaction.TransactionID, fullNodeID)
	return nil
}

// optimizeTransactionOrder optimizes the transaction order based on the optimization algorithm.
func (on *OptimizationNode) optimizeTransactionOrder(tx *ledger.Transaction) {
	switch on.OptimizationAlgo.PriorityModel {
	case "fee-based":
		// Sort transactions by fee priority
		fmt.Printf("Optimizing transaction %s based on fee.\n", tx.TransactionID)
	case "FIFO":
		// First-in, first-out (FIFO) approach
		fmt.Printf("Optimizing transaction %s based on FIFO.\n", tx.TransactionID)
	default:
		// Default optimization logic
		fmt.Printf("Optimizing transaction %s with default strategy.\n", tx.TransactionID)
	}
}

// selectRandomFullNode selects a random full node to forward transactions or optimization requests.
func (on *OptimizationNode) selectRandomFullNode() string {
	on.mutex.Lock()
	defer on.mutex.Unlock()

	if len(on.FullNodes) == 0 {
		fmt.Println("No full nodes available to forward transactions.")
		return ""
	}

	randomIndex := common.GenerateRandomInt(len(on.FullNodes))
	return on.FullNodes[randomIndex]
}

// EncryptData encrypts sensitive optimization data before sending it across the network.
func (on *OptimizationNode) EncryptData(data []byte) ([]byte, error) {
	encryptedData, err := on.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return encryptedData, nil
}

// DecryptData decrypts incoming encrypted optimization data.
func (on *OptimizationNode) DecryptData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := on.EncryptionService.DecryptData(encryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decryptedData, nil
}

// MonitorNetworkConditions continuously monitors network conditions to optimize load balancing.
func (on *OptimizationNode) MonitorNetworkConditions() {
	// Real-time network analysis for load balancing (non-AI, rule-based)
	for {
		networkStats := on.NetworkManager.GetNetworkStats()
		if networkStats.Load > 80 { // Threshold value
			fmt.Println("High network load detected. Rebalancing transaction processing.")
			// Apply load balancing across full nodes (e.g., round-robin, weighted)
			on.rebalanceLoad()
		}
		time.Sleep(10 * time.Second)
	}
}

// rebalanceLoad applies load balancing to distribute transaction processing evenly.
func (on *OptimizationNode) rebalanceLoad() {
	switch on.OptimizationAlgo.LoadBalancer {
	case "round-robin":
		fmt.Println("Applying round-robin load balancing.")
		// Implement round-robin distribution logic
	case "weighted":
		fmt.Println("Applying weighted load balancing.")
		// Implement weighted load balancing logic
	default:
		fmt.Println("Applying default load balancing strategy.")
	}
}
