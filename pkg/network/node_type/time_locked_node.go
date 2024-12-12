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

// TimeLockedNode represents a node responsible for managing time-locked transactions and contracts.
type TimeLockedNode struct {
	NodeID            string                        // Unique identifier for the Time-Locked Node
	PartialLedger     map[string]*ledger.Block      // Stores a subset of blocks (no full ledger stored)
	TimeLocks         map[string]*common.TimeLock   // Time-locked transactions and contracts
	ConsensusEngine   *synnergy_consensus.Engine    // Consensus engine for transaction validation
	EncryptionService *encryption.Encryption        // Encryption service for secure data management
	NetworkManager    *network.NetworkManager       // Manages communication with other nodes
	Ledger            *ledger.Ledger                // Ledger for recording time-locked transactions
	SNVM              *synnergy_vm.VirtualMachine   // Synnergy Virtual Machine for executing smart contracts
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with full nodes
	FullNodes         []string                      // List of full nodes to interact with
}

// NewTimeLockedNode initializes a new Time-Locked Node.
func NewTimeLockedNode(nodeID string, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, ledgerInstance *ledger.Ledger, snvm *synnergy_vm.VirtualMachine, syncInterval time.Duration, fullNodes []string) *TimeLockedNode {
	return &TimeLockedNode{
		NodeID:            nodeID,
		PartialLedger:     make(map[string]*ledger.Block),
		TimeLocks:         make(map[string]*common.TimeLock),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		Ledger:            ledgerInstance,
		SNVM:              snvm,
		SyncInterval:      syncInterval,
		FullNodes:         fullNodes,
	}
}

// StartNode starts the Time-Locked Node, initiating syncing with full nodes and listening for time-locked transaction requests.
func (tn *TimeLockedNode) StartNode() error {
	tn.mutex.Lock()
	defer tn.mutex.Unlock()

	// Start syncing with full nodes in the network.
	go tn.syncWithFullNodes()

	// Listen for incoming time-locked transaction requests.
	go tn.listenForTimeLockRequests()

	fmt.Printf("Time-Locked Node %s started successfully.\n", tn.NodeID)
	return nil
}

// syncWithFullNodes syncs the partial ledger with full nodes at regular intervals.
func (tn *TimeLockedNode) syncWithFullNodes() {
	ticker := time.NewTicker(tn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		tn.mutex.Lock()
		for _, fullNode := range tn.FullNodes {
			tn.syncPartialLedger(fullNode)
		}
		tn.mutex.Unlock()
	}
}

// syncPartialLedger syncs blocks from full nodes to maintain a partial ledger.
func (tn *TimeLockedNode) syncPartialLedger(fullNodeID string) {
	partialBlocks, err := tn.NetworkManager.RequestPartialBlockchain(fullNodeID)
	if err != nil {
		fmt.Printf("Failed to sync partial blockchain from node %s: %v\n", fullNodeID, err)
		return
	}

	// Validate and store the blocks in the partial ledger.
	for _, block := range partialBlocks {
		if tn.ConsensusEngine.ValidateBlock(block) {
			tn.PartialLedger[block.BlockID] = block
			fmt.Printf("Block %s synced from node %s.\n", block.BlockID, fullNodeID)
		} else {
			fmt.Printf("Block %s from node %s failed validation.\n", block.BlockID, fullNodeID)
		}
	}
}

// listenForTimeLockRequests listens for incoming time-locked transaction requests and processes them.
func (tn *TimeLockedNode) listenForTimeLockRequests() {
	for {
		request, err := tn.NetworkManager.ReceiveTransactionRequest()
		if err != nil {
			fmt.Printf("Error receiving transaction request: %v\n", err)
			continue
		}

		err = tn.processTimeLockTransactionRequest(request)
		if err != nil {
			fmt.Printf("Failed to process time-locked transaction: %v\n", err)
		}
	}
}

// processTimeLockTransactionRequest processes and stores time-locked transactions in the ledger.
func (tn *TimeLockedNode) processTimeLockTransactionRequest(request *network.TransactionRequest) error {
	tn.mutex.Lock()
	defer tn.mutex.Unlock()

	// Validate the transaction using the consensus engine.
	valid, err := tn.ConsensusEngine.ValidateTransaction(request.Transaction)
	if err != nil || !valid {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Store the transaction as time-locked.
	timeLock := &common.TimeLock{
		TransactionID: request.Transaction.TransactionID,
		LockTime:      request.TimeLock,
		Transaction:   request.Transaction,
	}

	tn.TimeLocks[timeLock.TransactionID] = timeLock

	// Log the time-locked transaction in the ledger.
	err = tn.Ledger.RecordTimeLockedTransaction(timeLock)
	if err != nil {
		return fmt.Errorf("failed to record time-locked transaction: %v", err)
	}

	fmt.Printf("Time-locked transaction %s stored successfully.\n", request.Transaction.TransactionID)
	return nil
}

// ExecuteTimeLockedTransactions checks if time-locked transactions are ready for execution and executes them.
func (tn *TimeLockedNode) ExecuteTimeLockedTransactions() {
	tn.mutex.Lock()
	defer tn.mutex.Unlock()

	currentTime := time.Now()

	for _, timeLock := range tn.TimeLocks {
		if currentTime.After(timeLock.LockTime) {
			err := tn.executeTransaction(timeLock.Transaction)
			if err != nil {
				fmt.Printf("Failed to execute time-locked transaction %s: %v\n", timeLock.TransactionID, err)
			} else {
				// Remove the transaction from the time-lock queue after execution.
				delete(tn.TimeLocks, timeLock.TransactionID)
				fmt.Printf("Time-locked transaction %s executed.\n", timeLock.TransactionID)
			}
		}
	}
}

// executeTransaction processes and finalizes the time-locked transaction.
func (tn *TimeLockedNode) executeTransaction(transaction *ledger.Transaction) error {
	// Execute the transaction using Synnergy Consensus and Virtual Machine.
	err := tn.ConsensusEngine.ExecuteTransaction(transaction, tn.SNVM)
	if err != nil {
		return fmt.Errorf("transaction execution failed: %v", err)
	}

	// Log the execution in the ledger.
	err = tn.Ledger.RecordTransactionExecution(transaction)
	if err != nil {
		return fmt.Errorf("failed to log transaction execution: %v", err)
	}

	return nil
}

// EncryptData encrypts sensitive data before transmission or storage.
func (tn *TimeLockedNode) EncryptData(data []byte) ([]byte, error) {
	encryptedData, err := tn.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return encryptedData, nil
}

// DecryptData decrypts incoming encrypted data from other nodes.
func (tn *TimeLockedNode) DecryptData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := tn.EncryptionService.DecryptData(encryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decryptedData, nil
}

// RequestBlock queries a full node for a specific block.
func (tn *TimeLockedNode) RequestBlock(blockID string) (*ledger.Block, error) {
	tn.mutex.Lock()
	defer tn.mutex.Unlock()

	// Check if the block is already in the partial ledger.
	block, exists := tn.PartialLedger[blockID]
	if exists {
		return block, nil
	}

	// Request the block from a full node.
	fullNodeID := tn.selectRandomFullNode()
	block, err := tn.NetworkManager.RequestBlockFromFullNode(fullNodeID, blockID)
	if err != nil {
		return nil, fmt.Errorf("failed to request block from full node: %v", err)
	}

	// Validate and store the block if valid.
	if tn.ConsensusEngine.ValidateBlock(block) {
		tn.PartialLedger[blockID] = block
		return block, nil
	}

	return nil, errors.New("block validation failed")
}

// selectRandomFullNode selects a random full node for forwarding requests.
func (tn *TimeLockedNode) selectRandomFullNode() string {
	tn.mutex.Lock()
	defer tn.mutex.Unlock()

	if len(tn.FullNodes) == 0 {
		fmt.Println("No full nodes available.")
		return ""
	}

	randomIndex := common.GenerateRandomInt(len(tn.FullNodes))
	return tn.FullNodes[randomIndex]
}
