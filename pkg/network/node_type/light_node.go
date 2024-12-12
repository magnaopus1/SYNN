package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/ledger"           // Blockchain ledger-related components
	"synnergy_network_demo/encryption"       // Encryption service for secure data handling
	"synnergy_network_demo/synnergy_consensus" // Synnergy Consensus engine
	"synnergy_network_demo/network"          // Network and communication management
	"synnergy_network_demo/common"           // Common utilities for random generation, keys, etc.
)

// LightNode represents a lightweight node in the blockchain network. 
// It validates transactions, syncs only headers, and interacts with full nodes for detailed queries.
type LightNode struct {
	NodeID            string                        // Unique identifier for the node
	BlockHeaders      map[string]*ledger.BlockHeader // Stores block headers, not full blocks
	ConsensusEngine   *synnergy_consensus.Engine    // Consensus engine for validation and querying
	EncryptionService *encryption.Encryption        // Encryption service for data encryption and decryption
	NetworkManager    *network.NetworkManager       // Network manager to communicate with full nodes
	FullNodes         []string                      // List of available full nodes for interaction
	SyncInterval      time.Duration                 // Interval for syncing block headers with full nodes
	mutex             sync.Mutex                    // Mutex to handle thread-safe operations
	SNVM              *common.VMInterface   // The Synnergy Network Virtual Machine
}

// NewLightNode initializes a new light node in the blockchain network.
func NewLightNode(nodeID string, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, syncInterval time.Duration, fullNodes []string) *LightNode {
	return &LightNode{
		NodeID:            nodeID,
		BlockHeaders:      make(map[string]*ledger.BlockHeader),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		FullNodes:         fullNodes,
		SyncInterval:      syncInterval,
	}
}

// StartNode begins the light node operations such as syncing block headers and listening for transactions.
func (ln *LightNode) StartNode() error {
	ln.mutex.Lock()
	defer ln.mutex.Unlock()

	// Start syncing block headers from full nodes.
	go ln.syncWithFullNodes()

	// Listen for incoming transaction requests and validation.
	go ln.listenForTransactionRequests()

	fmt.Printf("Light node %s started successfully.\n", ln.NodeID)
	return nil
}

// syncWithFullNodes syncs the block headers with full nodes at regular intervals.
func (ln *LightNode) syncWithFullNodes() {
	ticker := time.NewTicker(ln.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		ln.mutex.Lock()
		for _, fullNode := range ln.FullNodes {
			// Request latest block headers from full nodes.
			ln.syncBlockHeadersFromFullNode(fullNode)
		}
		ln.mutex.Unlock()
	}
}

// syncBlockHeadersFromFullNode fetches block headers from a full node.
func (ln *LightNode) syncBlockHeadersFromFullNode(fullNodeID string) {
	headers, err := ln.NetworkManager.RequestBlockHeaders(fullNodeID)
	if err != nil {
		fmt.Printf("Failed to sync block headers from full node %s: %v\n", fullNodeID, err)
		return
	}

	// Validate headers and update local store.
	for _, header := range headers {
		if ln.ConsensusEngine.ValidateBlockHeader(header) {
			ln.BlockHeaders[header.BlockID] = header
			fmt.Printf("Block header %s synced and validated from full node %s.\n", header.BlockID, fullNodeID)
		} else {
			fmt.Printf("Block header %s from full node %s failed validation.\n", header.BlockID, fullNodeID)
		}
	}
}

// listenForTransactionRequests listens for incoming transaction validation requests from the network.
func (ln *LightNode) listenForTransactionRequests() {
	for {
		transactionRequest, err := ln.NetworkManager.ReceiveTransactionRequest()
		if err != nil {
			fmt.Printf("Error receiving transaction request: %v\n", err)
			continue
		}

		// Process and validate the transaction request.
		err = ln.processTransactionRequest(transactionRequest)
		if err != nil {
			fmt.Printf("Transaction request processing failed: %v\n", err)
		}
	}
}

// processTransactionRequest handles transaction validation and forwards it to full nodes for block inclusion.
func (ln *LightNode) processTransactionRequest(request *network.TransactionRequest) error {
	ln.mutex.Lock()
	defer ln.mutex.Unlock()

	// Validate transaction using the consensus engine.
	valid, err := ln.ConsensusEngine.ValidateTransaction(request.Transaction)
	if err != nil || !valid {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Encrypt the transaction before sending it to a full node.
	encryptedTransaction, err := ln.EncryptTransaction(request.Transaction)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}

	// Forward the validated transaction to a full node for inclusion in the blockchain.
	fullNodeID := ln.selectRandomFullNode()
	err = ln.NetworkManager.ForwardTransactionToFullNode(encryptedTransaction, fullNodeID)
	if err != nil {
		return fmt.Errorf("failed to forward transaction to full node: %v", err)
	}

	fmt.Printf("Transaction %s validated and forwarded to full node %s.\n", request.Transaction.TransactionID, fullNodeID)
	return nil
}

// selectRandomFullNode selects a random full node to forward transactions or requests to.
func (ln *LightNode) selectRandomFullNode() string {
	ln.mutex.Lock()
	defer ln.mutex.Unlock()

	if len(ln.FullNodes) == 0 {
		fmt.Println("No full nodes available for transaction forwarding.")
		return ""
	}

	randomIndex := common.GenerateRandomInt(len(ln.FullNodes))
	return ln.FullNodes[randomIndex]
}

// RequestBlock queries a specific block from a full node using the block header information.
func (ln *LightNode) RequestBlock(blockID string) (*ledger.Block, error) {
	ln.mutex.Lock()
	defer ln.mutex.Unlock()

	// Check if the block header is available locally.
	header, exists := ln.BlockHeaders[blockID]
	if !exists {
		return nil, errors.New("block header not found in light node")
	}

	// Request the full block from a full node.
	fullNodeID := ln.selectRandomFullNode()
	block, err := ln.NetworkManager.RequestBlockFromFullNode(fullNodeID, blockID)
	if err != nil {
		return nil, fmt.Errorf("failed to request block from full node %s: %v", fullNodeID, err)
	}

	// Validate the block with the header and return it.
	if ln.ConsensusEngine.ValidateBlockWithHeader(block, header) {
		return block, nil
	}

	return nil, errors.New("block validation failed with the header")
}

// EncryptTransaction encrypts sensitive transaction data before sending it to full nodes.
func (ln *LightNode) EncryptTransaction(transaction *ledger.Transaction) (*ledger.Transaction, error) {
	encryptedData, err := ln.EncryptionService.EncryptData(transaction.Data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt transaction data: %v", err)
	}

	transaction.Data = encryptedData
	return transaction, nil
}

// DecryptTransaction decrypts the incoming encrypted transaction data from full nodes.
func (ln *LightNode) DecryptTransaction(encryptedTransaction *ledger.Transaction) (*ledger.Transaction, error) {
	decryptedData, err := ln.EncryptionService.DecryptData(encryptedTransaction.Data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt transaction data: %v", err)
	}

	encryptedTransaction.Data = decryptedData
	return encryptedTransaction, nil
}

