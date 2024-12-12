package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"    // Shared components like encryption, consensus, and storage
	"synnergy_network/pkg/ledger"    // Blockchain and ledger-related components
	"synnergy_network/pkg/cross_chain" // Cross-chain and external system-related components
	"synnergy_network/pkg/network"   // Network management and communication
)

// GatewayNode represents a node responsible for facilitating cross-chain interactions and integrating external data.
type GatewayNode struct {
	NodeID            string                        // Unique identifier for the node
	Blockchain        *ledger.Blockchain            // Local copy of the blockchain ledger
	ConsensusEngine   *common.SynnergyConsensus     // Consensus engine for validating transactions and data
	EncryptionService *common.Encryption            // Encryption service for secure data handling
	NetworkManager    *network.NetworkManager       // Manages communication with other nodes and external systems
	CrossChainManager *cross_chain.Manager          // Manages cross-chain bridges and transactions
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with other systems
	Bridges           map[string]*cross_chain.Bridge // Active cross-chain bridges for interoperability
	SNVM              *synnergy_vm.VirtualMachine // Virtual Machine for executing smart contracts

}

// NewGatewayNode initializes a new gateway node in the Synnergy Network.
func NewGatewayNode(nodeID string, blockchain *ledger.Blockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *network.NetworkManager, crossChainManager *cross_chain.Manager, syncInterval time.Duration) *GatewayNode {
	return &GatewayNode{
		NodeID:            nodeID,
		Blockchain:        blockchain,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		CrossChainManager: crossChainManager,
		SyncInterval:      syncInterval,
		Bridges:           make(map[string]*cross_chain.Bridge),
	}
}

// StartNode begins the gateway node's operations, syncing with other systems, managing cross-chain transactions, and integrating external data.
func (gn *GatewayNode) StartNode() error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	// Begin syncing external data and cross-chain activities.
	go gn.syncExternalData()
	go gn.manageCrossChainInteractions()

	fmt.Printf("Gateway Node %s started successfully.\n", gn.NodeID)
	return nil
}

// syncExternalData handles syncing external data sources with the blockchain.
func (gn *GatewayNode) syncExternalData() {
	ticker := time.NewTicker(gn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		gn.mutex.Lock()
		err := gn.pullExternalData()
		if err != nil {
			fmt.Printf("Error syncing external data: %v\n", err)
		}
		gn.mutex.Unlock()
	}
}

// pullExternalData retrieves data from external sources such as IoT devices or traditional databases, encrypts it, and records it on the blockchain.
func (gn *GatewayNode) pullExternalData() error {
	// Retrieve data from external systems (IoT devices, APIs, etc.).
	externalData, err := gn.NetworkManager.FetchExternalData()
	if err != nil {
		return fmt.Errorf("failed to retrieve external data: %v", err)
	}

	// Encrypt the external data before recording.
	encryptedData, err := gn.EncryptionService.EncryptData(externalData, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt external data: %v", err)
	}

	// Record encrypted data in the blockchain ledger.
	err = gn.Blockchain.RecordData(encryptedData)
	if err != nil {
		return fmt.Errorf("failed to record external data to blockchain: %v", err)
	}

	fmt.Printf("External data successfully recorded to the blockchain by node %s.\n", gn.NodeID)
	return nil
}

// manageCrossChainInteractions manages cross-chain transactions and communication between different blockchain protocols.
func (gn *GatewayNode) manageCrossChainInteractions() {
	for {
		transaction, err := gn.NetworkManager.ReceiveTransaction()
		if err != nil {
			fmt.Printf("Error receiving transaction: %v\n", err)
			continue
		}

		// Validate and process the cross-chain transaction.
		err = gn.processCrossChainTransaction(transaction)
		if err != nil {
			fmt.Printf("Cross-chain transaction processing failed: %v\n", err)
		}
	}
}

// processCrossChainTransaction processes and validates a cross-chain transaction.
func (gn *GatewayNode) processCrossChainTransaction(tx *ledger.Transaction) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	// Validate the transaction using the consensus engine.
	if valid, err := gn.ConsensusEngine.ValidateTransaction(tx); err != nil || !valid {
		return fmt.Errorf("invalid transaction: %v", err)
	}

	// Handle cross-chain transaction if applicable.
	if tx.IsCrossChain {
		err := gn.executeCrossChainTransfer(tx)
		if err != nil {
			return fmt.Errorf("cross-chain transaction failed: %v", err)
		}
	}

	// Add the transaction to a sub-block for validation.
	subBlockID := common.GenerateSubBlockID()
	subBlock := gn.createSubBlock(tx, subBlockID)

	// Store the sub-block and validate it into a full block.
	gn.Blockchain.AddSubBlock(subBlock)
	gn.validateSubBlock(subBlock)

	fmt.Printf("Transaction %s processed successfully by gateway node %s.\n", tx.TransactionID, gn.NodeID)
	return nil
}

// executeCrossChainTransfer handles the transfer of assets between different blockchain networks using a cross-chain bridge.
func (gn *GatewayNode) executeCrossChainTransfer(tx *ledger.Transaction) error {
	// Identify the destination blockchain and initiate the cross-chain bridge.
	bridge, exists := gn.Bridges[tx.CrossChainID]
	if !exists {
		return errors.New("cross-chain bridge not found")
	}

	// Execute the asset transfer or data transfer across the chains.
	err := bridge.Transfer(tx)
	if err != nil {
		return fmt.Errorf("cross-chain transfer failed: %v", err)
	}

	// Log the cross-chain transaction in the ledger.
	err = gn.Blockchain.RecordCrossChainTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to record cross-chain transaction: %v", err)
	}

	fmt.Printf("Cross-chain transaction %s successfully processed via bridge %s.\n", tx.TransactionID, bridge.BridgeID)
	return nil
}

// createSubBlock creates a sub-block from a validated transaction.
func (gn *GatewayNode) createSubBlock(tx *ledger.Transaction, subBlockID string) *common.SubBlock {
	return &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: []*ledger.Transaction{tx},
		Timestamp:    time.Now(),
		NodeID:       gn.NodeID,
	}
}

// validateSubBlock validates a sub-block into a full block within the blockchain.
func (gn *GatewayNode) validateSubBlock(subBlock *common.SubBlock) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	// Validate the sub-block using the consensus mechanism.
	if gn.ConsensusEngine.ValidateSubBlock(subBlock) {
		block := gn.Blockchain.AddSubBlock(subBlock)
		gn.NetworkManager.BroadcastNewBlock(block)
		fmt.Printf("Sub-block %s validated and added to blockchain.\n", subBlock.SubBlockID)
		return nil
	}

	return errors.New("failed to validate sub-block")
}

// Secure Data Handling and Cross-Chain Management

// setupCrossChainBridges sets up cross-chain bridges for managing transactions between different blockchain networks.
func (gn *GatewayNode) setupCrossChainBridges() {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	// Establish bridges to other blockchain networks.
	bridges := gn.CrossChainManager.GetAvailableBridges()
	for _, bridge := range bridges {
		gn.Bridges[bridge.BridgeID] = bridge
		fmt.Printf("Cross-chain bridge %s established.\n", bridge.BridgeID)
	}
}

// Data Encryption and Security Protocols

// applyEncryptionSecurity ensures encryption protocols are applied to all external and cross-chain data before recording or transfer.
func (gn *GatewayNode) applyEncryptionSecurity() error {
	err := gn.EncryptionService.ApplySecurity(gn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to apply encryption security for gateway node: %v", err)
	}
	fmt.Printf("Encryption security applied successfully for gateway node %s.\n", gn.NodeID)
	return nil
}

// Compliance and Regulatory Monitoring

// ensureCompliance verifies that cross-chain and external transactions comply with regulatory frameworks.
func (gn *GatewayNode) ensureCompliance(tx *ledger.Transaction) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	// Check compliance for the cross-chain transaction.
	err := gn.CrossChainManager.CheckCompliance(tx)
	if err != nil {
		return fmt.Errorf("compliance check failed for transaction %s: %v", tx.TransactionID, err)
	}

	fmt.Printf("Transaction %s passed compliance check on gateway node %s.\n", tx.TransactionID, gn.NodeID)
	return nil
}
