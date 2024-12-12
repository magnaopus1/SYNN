package node_type

import (
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/ledger"           // Blockchain ledger components
	"synnergy_network_demo/encryption"       // Encryption service for securing data
	"synnergy_network_demo/synnergy_consensus" // Synnergy Consensus engine
	"synnergy_network_demo/network"          // Network management
	"synnergy_network_demo/common"           // Common utilities for encryption, random generation, etc.
)

// OrphanNode represents a node responsible for handling orphan blocks in the blockchain.
type OrphanNode struct {
	NodeID            string                        // Unique identifier for the node
	OrphanBlocks      map[string]*ledger.Block      // Stores detected orphan blocks
	ConsensusEngine   *synnergy_consensus.Engine    // Consensus engine for validating blocks and resolving conflicts
	EncryptionService *encryption.Encryption        // Encryption service for securing data transfers
	NetworkManager    *network.NetworkManager       // Network manager for communicating with other nodes
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with other nodes
	FullNodes         []string                      // List of full nodes for syncing and querying
	SNVM              *common.VMInterface   // The Synnergy Network Virtual Machine
	ArchivedBlocks    map[string]*ledger.Block      // Storage for archived orphan blocks for future reference
}

// NewOrphanNode initializes a new orphan node in the blockchain network.
func NewOrphanNode(nodeID string, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, syncInterval time.Duration, fullNodes []string) *OrphanNode {
	return &OrphanNode{
		NodeID:            nodeID,
		OrphanBlocks:      make(map[string]*ledger.Block),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SyncInterval:      syncInterval,
		FullNodes:         fullNodes,
		ArchivedBlocks:    make(map[string]*ledger.Block),
	}
}

// StartNode begins the operations of the Orphan Node, including detection, analysis, and resource recovery.
func (on *OrphanNode) StartNode() error {
	on.mutex.Lock()
	defer on.mutex.Unlock()

	// Start syncing with full nodes periodically based on the sync interval.
	go on.syncWithFullNodes()

	// Start listening for orphan block detection.
	go on.listenForOrphanBlocks()

	fmt.Printf("Orphan node %s started successfully.\n", on.NodeID)
	return nil
}

// syncWithFullNodes syncs the Orphan Node's block data with full nodes to detect orphan blocks.
func (on *OrphanNode) syncWithFullNodes() {
	ticker := time.NewTicker(on.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		on.mutex.Lock()
		for _, fullNode := range on.FullNodes {
			// Sync block data with full nodes to detect any orphan blocks.
			on.detectAndHandleOrphanBlocks(fullNode)
		}
		on.mutex.Unlock()
	}
}

// detectAndHandleOrphanBlocks detects orphan blocks by comparing blockchains from full nodes.
func (on *OrphanNode) detectAndHandleOrphanBlocks(fullNodeID string) {
	blockchainData, err := on.NetworkManager.RequestBlockchainData(fullNodeID)
	if err != nil {
		fmt.Printf("Failed to sync blockchain data from full node %s: %v\n", fullNodeID, err)
		return
	}

	for _, block := range blockchainData {
		// Check if block is orphaned (valid block but not part of the longest chain)
		if on.ConsensusEngine.IsOrphanBlock(block) {
			on.OrphanBlocks[block.BlockID] = block
			fmt.Printf("Detected orphan block %s from full node %s.\n", block.BlockID, fullNodeID)

			// Analyze and attempt to recover resources from the orphan block
			on.analyzeAndRecoverOrphanBlock(block)
		}
	}
}

// listenForOrphanBlocks listens for incoming orphan block alerts from other nodes.
func (on *OrphanNode) listenForOrphanBlocks() {
	for {
		alert, err := on.NetworkManager.ReceiveOrphanBlockAlert()
		if err != nil {
			fmt.Printf("Error receiving orphan block alert: %v\n", err)
			continue
		}

		// Process the orphan block alert.
		err = on.processOrphanBlockAlert(alert)
		if err != nil {
			fmt.Printf("Orphan block alert processing failed: %v\n", err)
		}
	}
}

// processOrphanBlockAlert processes incoming orphan block alerts and handles them.
func (on *OrphanNode) processOrphanBlockAlert(alert *network.OrphanBlockAlert) error {
	on.mutex.Lock()
	defer on.mutex.Unlock()

	// Fetch the orphan block and analyze it.
	block, err := on.NetworkManager.RequestBlockFromFullNode(alert.FullNodeID, alert.BlockID)
	if err != nil {
		return fmt.Errorf("failed to request orphan block from full node: %v", err)
	}

	if on.ConsensusEngine.IsOrphanBlock(block) {
		on.OrphanBlocks[block.BlockID] = block
		fmt.Printf("Orphan block %s processed and added to orphan block pool.\n", block.BlockID)

		// Analyze the orphan block and recover resources.
		on.analyzeAndRecoverOrphanBlock(block)
	}

	return nil
}

// analyzeAndRecoverOrphanBlock analyzes an orphan block and attempts to recover its resources.
func (on *OrphanNode) analyzeAndRecoverOrphanBlock(block *ledger.Block) {
	// Validate transactions and determine which can be reintegrated into the pool.
	for _, tx := range block.Transactions {
		if on.ConsensusEngine.ValidateTransaction(tx) {
			fmt.Printf("Reintegrating valid transaction %s from orphan block %s.\n", tx.TransactionID, block.BlockID)
			err := on.NetworkManager.ReintroduceTransactionToPool(tx)
			if err != nil {
				fmt.Printf("Failed to reintegrate transaction %s: %v\n", tx.TransactionID, err)
			}
		} else {
			fmt.Printf("Transaction %s from orphan block %s is invalid or conflicting.\n", tx.TransactionID, block.BlockID)
		}
	}

	// Archive the orphan block for future reference or analysis.
	on.ArchiveOrphanBlock(block)
}

// ArchiveOrphanBlock archives an orphan block for future reference or analysis.
func (on *OrphanNode) ArchiveOrphanBlock(block *ledger.Block) {
	on.mutex.Lock()
	defer on.mutex.Unlock()

	on.ArchivedBlocks[block.BlockID] = block
	fmt.Printf("Orphan block %s archived for future reference.\n", block.BlockID)
}

// EncryptData encrypts sensitive data related to orphan blocks before transmission.
func (on *OrphanNode) EncryptData(data []byte) ([]byte, error) {
	encryptedData, err := on.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return encryptedData, nil
}

// DecryptData decrypts incoming encrypted data related to orphan blocks.
func (on *OrphanNode) DecryptData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := on.EncryptionService.DecryptData(encryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decryptedData, nil
}

// ReclaimResources reclaims computational resources from orphan blocks.
func (on *OrphanNode) ReclaimResources() {
	on.mutex.Lock()
	defer on.mutex.Unlock()

	for blockID, block := range on.OrphanBlocks {
		fmt.Printf("Reclaiming resources from orphan block %s.\n", blockID)
		// Free up memory and storage resources
		delete(on.OrphanBlocks, blockID)
	}
}
