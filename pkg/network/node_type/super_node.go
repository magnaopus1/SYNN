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

// SuperNode represents a high-capacity node responsible for handling transaction routing, data storage, executing smart contracts, and privacy features.
type SuperNode struct {
	NodeID            string                        // Unique identifier for the super node
	StorageCapacity   int64                         // Maximum storage capacity of the node in bytes
	UsedStorage       int64                         // Amount of used storage in bytes
	StoredBlocks      map[string]*ledger.Block      // Map of stored blocks by block ID
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	ConsensusEngine   *synnergy_consensus.Engine    // Consensus engine for validating transactions and blocks
	EncryptionService *encryption.Encryption        // Encryption service for securing transaction data
	NetworkManager    *network.NetworkManager       // Network manager for communicating with other nodes
	Ledger            *ledger.Ledger                // Reference to the ledger for logging transactions
	SNVM              *common.VMInterface   // The Synnergy Network Virtual Machine
}

// NewSuperNode initializes a new SuperNode.
func NewSuperNode(nodeID string, storageCapacity int64, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, ledgerInstance *ledger.Ledger, snvm *synnergy_vm.VirtualMachine) *SuperNode {
	return &SuperNode{
		NodeID:            nodeID,
		StorageCapacity:   storageCapacity,
		UsedStorage:       0,
		StoredBlocks:      make(map[string]*ledger.Block),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		Ledger:            ledgerInstance,
		SNVM:              snvm,
	}
}

// StartNode starts the SuperNode’s operations, routing transactions, syncing data, and listening for smart contract requests.
func (sn *SuperNode) StartNode() error {
	sn.mutex.Lock()
	defer sn.mutex.Unlock()

	// Start syncing with the network.
	go sn.syncWithNetwork()

	// Listen for incoming transactions and routing requests.
	go sn.listenForTransactionRequests()

	// Listen for smart contract execution requests.
	go sn.listenForSmartContractRequests()

	fmt.Printf("SuperNode %s started successfully.\n", sn.NodeID)
	return nil
}

// syncWithNetwork syncs blocks and transactions with other full nodes to stay updated.
func (sn *SuperNode) syncWithNetwork() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sn.mutex.Lock()
		// Request updates from other full nodes in the network.
		sn.syncWithFullNodes()
		sn.mutex.Unlock()
	}
}

// syncWithFullNodes syncs data with other full nodes in the network.
func (sn *SuperNode) syncWithFullNodes() {
	for _, fullNodeID := range sn.NetworkManager.GetFullNodes() {
		// Request recent blocks and transactions.
		blocks, err := sn.NetworkManager.RequestBlocks(fullNodeID)
		if err != nil {
			fmt.Printf("Failed to sync with full node %s: %v\n", fullNodeID, err)
			continue
		}

		for _, block := range blocks {
			if sn.ConsensusEngine.ValidateBlock(block) {
				sn.StoredBlocks[block.BlockID] = block
				fmt.Printf("Block %s synced from full node %s.\n", block.BlockID, fullNodeID)
			}
		}
	}
}

// listenForTransactionRequests listens for transaction routing requests from other nodes.
func (sn *SuperNode) listenForTransactionRequests() {
	for {
		request, err := sn.NetworkManager.ReceiveTransactionRequest()
		if err != nil {
			fmt.Printf("Error receiving transaction request: %v\n", err)
			continue
		}

		err = sn.processTransactionRequest(request)
		if err != nil {
			fmt.Printf("Transaction request processing failed: %v\n", err)
		}
	}
}

// processTransactionRequest processes an incoming transaction request, validates it, and forwards it.
func (sn *SuperNode) processTransactionRequest(request *network.TransactionRequest) error {
	sn.mutex.Lock()
	defer sn.mutex.Unlock()

	// Validate the transaction using the consensus engine.
	valid, err := sn.ConsensusEngine.ValidateTransaction(request.Transaction)
	if err != nil || !valid {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Forward the validated transaction to other nodes or add it to the block.
	err = sn.forwardTransaction(request.Transaction)
	if err != nil {
		return fmt.Errorf("failed to forward transaction: %v", err)
	}

	fmt.Printf("Transaction %s validated and processed by SuperNode %s.\n", request.Transaction.TransactionID, sn.NodeID)
	return nil
}

// forwardTransaction forwards a validated transaction to other full nodes for inclusion in the next block.
func (sn *SuperNode) forwardTransaction(transaction *ledger.Transaction) error {
	fullNodeID := sn.NetworkManager.SelectFullNode()
	err := sn.NetworkManager.ForwardTransaction(transaction, fullNodeID)
	if err != nil {
		return fmt.Errorf("failed to forward transaction to full node %s: %v", fullNodeID, err)
	}

	fmt.Printf("Transaction %s forwarded to full node %s.\n", transaction.TransactionID, fullNodeID)
	return nil
}

// listenForSmartContractRequests listens for incoming smart contract execution requests.
func (sn *SuperNode) listenForSmartContractRequests() {
	for {
		request, err := sn.NetworkManager.ReceiveSmartContractRequest()
		if err != nil {
			fmt.Printf("Error receiving smart contract request: %v\n", err)
			continue
		}

		err = sn.executeSmartContract(request)
		if err != nil {
			fmt.Printf("Smart contract execution failed: %v\n", err)
		}
	}
}

// executeSmartContract executes a smart contract using the Synnergy Virtual Machine.
func (sn *SuperNode) executeSmartContract(request *network.SmartContractRequest) error {
	sn.mutex.Lock()
	defer sn.mutex.Unlock()

	// Execute the smart contract using the Synnergy Virtual Machine.
	err := sn.SNVM.ExecuteSmartContract(request.Contract)
	if err != nil {
		return fmt.Errorf("smart contract execution failed: %v", err)
	}

	// Log the execution result in the ledger.
	err = sn.Ledger.RecordSmartContractExecution(request.Contract.ContractID)
	if err != nil {
		return fmt.Errorf("failed to log smart contract execution: %v", err)
	}

	fmt.Printf("Smart contract %s executed successfully on SuperNode %s.\n", request.Contract.ContractID, sn.NodeID)
	return nil
}

// StoreBlock stores a block on the SuperNode after validation.
func (sn *SuperNode) StoreBlock(block *ledger.Block) error {
	sn.mutex.Lock()
	defer sn.mutex.Unlock()

	if sn.UsedStorage+block.Size > sn.StorageCapacity {
		return errors.New("insufficient storage capacity")
	}

	// Store the block if it is valid.
	if sn.ConsensusEngine.ValidateBlock(block) {
		sn.StoredBlocks[block.BlockID] = block
		sn.UsedStorage += block.Size
		fmt.Printf("Block %s stored successfully on SuperNode %s.\n", block.BlockID, sn.NodeID)
	} else {
		return fmt.Errorf("block validation failed for block %s", block.BlockID)
	}

	return nil
}

// RetrieveBlock retrieves a block by its ID.
func (sn *SuperNode) RetrieveBlock(blockID string) (*ledger.Block, error) {
	sn.mutex.Lock()
	defer sn.mutex.Unlock()

	block, exists := sn.StoredBlocks[blockID]
	if !exists {
		return nil, errors.New("block not found")
	}

	return block, nil
}

// EncryptData encrypts sensitive data before transmission or storage.
func (sn *SuperNode) EncryptData(data []byte) ([]byte, error) {
	encryptedData, err := sn.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return encryptedData, nil
}

// DecryptData decrypts incoming encrypted data from other nodes.
func (sn *SuperNode) DecryptData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := sn.EncryptionService.DecryptData(encryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decryptedData, nil
}
