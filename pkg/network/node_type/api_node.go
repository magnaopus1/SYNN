package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"net/http"
	"encoding/json"

	"synnergy_network/pkg/common"     // Shared components like encryption, consensus, sub-blocks, etc.
	"synnergy_network/pkg/ledger"     // Blockchain and ledger-related components
)

// APINode represents an API node in the Synnergy Network, facilitating high-throughput API requests and blockchain interaction.
type APINode struct {
	NodeID            string                        // Unique identifier for the node
	Blockchain        *ledger.Blockchain            // Local copy of the blockchain for data retrieval and transaction processing
	ConsensusEngine   *common.SynnergyConsensus     // Consensus engine for validating transactions and data queries
	EncryptionService *common.Encryption            // Encryption service for secure communication
	NetworkManager    *common.NetworkManager        // Network manager for communication with other nodes
	APIManager        *common.APIManager            // API manager for handling external API requests
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with the blockchain network
	SNVM              *common.SynnergyVirtualMachine // Virtual Machine instance for executing smart contracts
	ActiveRequests    map[string]*common.APIRequest  // Map of active API requests for tracking and management
	RequestTimeout    time.Duration                 // Timeout for handling API requests
	LoadBalancer      *common.LoadBalancer          // Load balancer to distribute API request loads
}

// NewAPINode initializes a new API node for the Synnergy Network.
func NewAPINode(nodeID string, blockchain *ledger.Blockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *common.NetworkManager, apiManager *common.APIManager, syncInterval time.Duration, requestTimeout time.Duration) *APINode {
	return &APINode{
		NodeID:            nodeID,
		Blockchain:        blockchain,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		APIManager:        apiManager,
		ActiveRequests:    make(map[string]*common.APIRequest),
		SyncInterval:      syncInterval,
		RequestTimeout:    requestTimeout,
		LoadBalancer:      common.NewLoadBalancer(),  // Initialize the load balancer
		SNVM:              common.NewSynnergyVirtualMachine(), // Initialize the virtual machine
	}
}

// StartNode starts the API node's operations, including syncing, handling API requests, and validating transactions.
func (an *APINode) StartNode() error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Start syncing with the blockchain network and listening for API requests.
	go an.syncWithOtherNodes()
	go an.listenForAPIRequests()

	fmt.Printf("API node %s started successfully.\n", an.NodeID)
	return nil
}

// syncWithOtherNodes handles syncing the blockchain with other nodes at regular intervals.
func (an *APINode) syncWithOtherNodes() {
	ticker := time.NewTicker(an.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		an.mutex.Lock()
		otherNodes := an.NetworkManager.DiscoverOtherNodes(an.NodeID)
		for _, node := range otherNodes {
			// Sync blockchain from the other node.
			an.syncBlockchainFromNode(node)
		}
		an.mutex.Unlock()
	}
}

// syncBlockchainFromNode syncs the blockchain from a peer node.
func (an *APINode) syncBlockchainFromNode(peerNode string) {
	peerBlockchain, err := an.NetworkManager.RequestBlockchain(peerNode)
	if err != nil {
		fmt.Printf("Failed to sync blockchain from node %s: %v\n", peerNode, err)
		return
	}

	// Validate the blockchain and update the local copy if necessary.
	if an.ConsensusEngine.ValidateBlockchain(peerBlockchain) {
		an.Blockchain = an.Blockchain.MergeWith(peerBlockchain)
		fmt.Printf("Blockchain synced successfully from node %s.\n", peerNode)
	} else {
		fmt.Printf("Blockchain sync from node %s failed validation.\n", peerNode)
	}
}

// listenForAPIRequests listens for incoming API requests from external clients.
func (an *APINode) listenForAPIRequests() {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// Parse the incoming API request.
		var apiRequest common.APIRequest
		err := json.NewDecoder(r.Body).Decode(&apiRequest)
		if err != nil {
			http.Error(w, "Invalid API request", http.StatusBadRequest)
			return
		}

		// Process the API request.
		response, err := an.processAPIRequest(&apiRequest)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to process request: %v", err), http.StatusInternalServerError)
			return
		}

		// Send the response back to the client.
		json.NewEncoder(w).Encode(response)
	})

	// Start the HTTP server.
	http.ListenAndServe(":8080", nil)
}

// processAPIRequest processes and validates an incoming API request.
func (an *APINode) processAPIRequest(apiRequest *common.APIRequest) (*common.APIResponse, error) {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Validate the API request using the consensus engine.
	if valid, err := an.ConsensusEngine.ValidateAPIRequest(apiRequest); err != nil || !valid {
		return nil, fmt.Errorf("invalid API request: %v", err)
	}

	// Process the request based on its type.
	switch apiRequest.RequestType {
	case "query":
		return an.handleQueryRequest(apiRequest)
	case "transaction":
		return an.handleTransactionRequest(apiRequest)
	default:
		return nil, errors.New("unsupported API request type")
	}
}

// handleQueryRequest handles blockchain data queries.
func (an *APINode) handleQueryRequest(apiRequest *common.APIRequest) (*common.APIResponse, error) {
	// Decrypt the query parameters.
	decryptedParams, err := an.EncryptionService.DecryptData(apiRequest.Params, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt query parameters: %v", err)
	}

	// Execute the query on the blockchain.
	data, err := an.Blockchain.QueryData(decryptedParams)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	// Encrypt the response data.
	encryptedData, err := an.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt response data: %v", err)
	}

	return &common.APIResponse{
		RequestID: apiRequest.RequestID,
		Status:    "success",
		Data:      encryptedData,
	}, nil
}

// handleTransactionRequest handles transaction submissions via the API.
func (an *APINode) handleTransactionRequest(apiRequest *common.APIRequest) (*common.APIResponse, error) {
	// Decrypt the transaction data.
	decryptedTx, err := an.EncryptionService.DecryptData(apiRequest.Params, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt transaction data: %v", err)
	}

	// Validate the transaction using the consensus engine.
	tx := common.ParseTransaction(decryptedTx)
	if valid, err := an.ConsensusEngine.ValidateTransaction(tx); err != nil || !valid {
		return nil, fmt.Errorf("invalid transaction: %v", err)
	}

	// Add the transaction to a sub-block.
	subBlockID := common.GenerateSubBlockID()
	subBlock := an.createSubBlock(tx, subBlockID)

	// Store the sub-block and attempt to validate it into a block.
	an.SubBlocks[subBlockID] = subBlock
	if err := an.tryValidateSubBlock(subBlock); err != nil {
		return nil, fmt.Errorf("failed to validate sub-block: %v", err)
	}

	return &common.APIResponse{
		RequestID: apiRequest.RequestID,
		Status:    "success",
		Data:      []byte(fmt.Sprintf("Transaction %s processed successfully", tx.TransactionID)),
	}, nil
}

// createSubBlock creates a sub-block from a validated transaction.
func (an *APINode) createSubBlock(tx *ledger.Transaction, subBlockID string) *common.SubBlock {
	return &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: []*ledger.Transaction{tx},
		Timestamp:    time.Now(),
		NodeID:       an.NodeID,
	}
}

// tryValidateSubBlock tries to validate a sub-block into a block.
func (an *APINode) tryValidateSubBlock(subBlock *common.SubBlock) error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Check if the sub-block is valid using the consensus mechanism.
	if an.ConsensusEngine.ValidateSubBlock(subBlock) {
		// Add the sub-block to the blockchain as part of a full block.
		block := an.Blockchain.AddSubBlock(subBlock)

		// Notify the network of the new block.
		err := an.NetworkManager.BroadcastNewBlock(block)
		if err != nil {
			return fmt.Errorf("failed to broadcast new block: %v", err)
		}

		fmt.Printf("Sub-block %s validated and added to blockchain.\n", subBlock.SubBlockID)
		return nil
	}

	return errors.New("failed to validate sub-block")
}

// Security and Load Balancing

// ApplySecurityProtocols applies all necessary security measures for API requests.
func (an *APINode) ApplySecurityProtocols() error {
	// Implement security protocols such as OAuth, JWT, and SSL/TLS encryption.
	err := an.APIManager.ApplySecurity(an.NodeID)
	if err != nil {
		return fmt.Errorf("failed to apply security protocols: %v", err)
	}

	fmt.Printf("Security protocols applied successfully for API node %s.\n", an.NodeID)
	return nil
}

// DistributeLoad manages load balancing across API nodes to handle high traffic.
func (an *APINode) DistributeLoad(requests []common.APIRequest) error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Distribute the requests across available API nodes using the load balancer.
	err := an.LoadBalancer.DistributeRequests(requests, an.NodeID)
	if err != nil {
		return fmt.Errorf("failed to distribute load: %v", err)
	}

	fmt.Printf("Load distributed successfully across nodes by API node %s.\n", an.NodeID)
	return nil
}
