package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"    // Shared components for encryption, consensus, and more
	"synnergy_network/pkg/ledger"    // Blockchain and ledger-related components
	"synnergy_network/pkg/geospatial" // Geospatial processing and data management components
	"synnergy_network/pkg/network"   // Network and communication management
	"synnergy_network/pkg/synnergy_vm" // Synnergy Virtual Machine for smart contract execution
)

// GeospatialNode represents a node that processes and validates transactions involving geospatial data.
type GeospatialNode struct {
	NodeID            string                      // Unique identifier for the node
	Blockchain        *ledger.Blockchain          // Local copy of the blockchain ledger
	ConsensusEngine   *common.SynnergyConsensus   // Consensus engine for validating geospatial data transactions
	EncryptionService *common.Encryption          // Encryption service for securing sensitive data
	NetworkManager    *network.NetworkManager     // Handles communication with other nodes and data sources
	GeoProcessor      *geospatial.Processor       // Processes geospatial data and computations
	SNVM              *synnergy_vm.VirtualMachine // Virtual Machine for executing smart contracts
	mutex             sync.Mutex                  // Mutex for thread-safe operations
	SyncInterval      time.Duration               // Interval for syncing with geospatial data sources
}

// NewGeospatialNode initializes a new geospatial node for handling geographical data and smart contract execution.
func NewGeospatialNode(nodeID string, blockchain *ledger.Blockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *network.NetworkManager, geoProcessor *geospatial.Processor, snvm *synnergy_vm.VirtualMachine, syncInterval time.Duration) *GeospatialNode {
	return &GeospatialNode{
		NodeID:            nodeID,
		Blockchain:        blockchain,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		GeoProcessor:      geoProcessor,
		SNVM:              snvm,
		SyncInterval:      syncInterval,
	}
}

// StartNode begins the geospatial node's operations, syncing geospatial data, processing transactions, and validating geospatial triggers.
func (gn *GeospatialNode) StartNode() error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	// Begin syncing geospatial data and monitoring transactions.
	go gn.syncGeospatialData()
	go gn.monitorGeospatialTransactions()

	fmt.Printf("Geospatial Node %s started successfully.\n", gn.NodeID)
	return nil
}

// syncGeospatialData handles syncing geospatial data sources with the blockchain.
func (gn *GeospatialNode) syncGeospatialData() {
	ticker := time.NewTicker(gn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		gn.mutex.Lock()
		err := gn.pullGeospatialData()
		if err != nil {
			fmt.Printf("Error syncing geospatial data: %v\n", err)
		}
		gn.mutex.Unlock()
	}
}

// pullGeospatialData retrieves and processes data from geospatial sources such as satellites, sensors, and APIs.
func (gn *GeospatialNode) pullGeospatialData() error {
	// Retrieve geospatial data from external systems.
	geoData, err := gn.NetworkManager.FetchGeospatialData()
	if err != nil {
		return fmt.Errorf("failed to retrieve geospatial data: %v", err)
	}

	// Process geospatial data for validation.
	processedData, err := gn.GeoProcessor.ProcessData(geoData)
	if err != nil {
		return fmt.Errorf("failed to process geospatial data: %v", err)
	}

	// Encrypt the processed data before recording it.
	encryptedData, err := gn.EncryptionService.EncryptData(processedData, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt geospatial data: %v", err)
	}

	// Record encrypted data in the blockchain ledger.
	err = gn.Blockchain.RecordData(encryptedData)
	if err != nil {
		return fmt.Errorf("failed to record geospatial data to blockchain: %v", err)
	}

	fmt.Printf("Geospatial data successfully recorded by node %s.\n", gn.NodeID)
	return nil
}

// monitorGeospatialTransactions monitors blockchain transactions related to geospatial data, validating them based on geographical triggers.
func (gn *GeospatialNode) monitorGeospatialTransactions() {
	for {
		transaction, err := gn.NetworkManager.ReceiveTransaction()
		if err != nil {
			fmt.Printf("Error receiving transaction: %v\n", err)
			continue
		}

		// Validate and process the transaction.
		err = gn.processGeospatialTransaction(transaction)
		if err != nil {
			fmt.Printf("Geospatial transaction processing failed: %v\n", err)
		}
	}
}

// processGeospatialTransaction validates and processes transactions involving geospatial data, including executing smart contracts via the virtual machine.
func (gn *GeospatialNode) processGeospatialTransaction(tx *ledger.Transaction) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	// Validate the transaction using the consensus engine.
	if valid, err := gn.ConsensusEngine.ValidateTransaction(tx); err != nil || !valid {
		return fmt.Errorf("invalid transaction: %v", err)
	}

	// Check for geospatial triggers in the transaction.
	if tx.HasGeospatialTriggers() {
		err := gn.executeGeospatialSmartContract(tx)
		if err != nil {
			return fmt.Errorf("geospatial smart contract execution failed: %v", err)
		}
	}

	// Add the transaction to a sub-block for validation.
	subBlockID := common.GenerateSubBlockID()
	subBlock := gn.createSubBlock(tx, subBlockID)

	// Validate and store the sub-block.
	gn.Blockchain.AddSubBlock(subBlock)
	gn.validateSubBlock(subBlock)

	fmt.Printf("Geospatial transaction %s processed successfully by node %s.\n", tx.TransactionID, gn.NodeID)
	return nil
}

// executeGeospatialSmartContract executes a smart contract based on geospatial triggers using the virtual machine.
func (gn *GeospatialNode) executeGeospatialSmartContract(tx *ledger.Transaction) error {
	// Encrypt the contract code or data before execution.
	encryptedArgs, err := gn.EncryptionService.EncryptData(tx.Args, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt contract arguments: %v", err)
	}

	// Execute the contract on the virtual machine.
	result, err := gn.SNVM.ExecuteContract(tx.ContractID, encryptedArgs)
	if err != nil {
		return fmt.Errorf("smart contract execution failed: %v", err)
	}

	// Record the result of the contract execution.
	err = gn.Blockchain.RecordContractExecution(tx.ContractID, result)
	if err != nil {
		return fmt.Errorf("failed to record contract execution: %v", err)
	}

	fmt.Printf("Geospatial smart contract executed successfully for transaction %s.\n", tx.TransactionID)
	return nil
}

// createSubBlock creates a sub-block from a validated transaction.
func (gn *GeospatialNode) createSubBlock(tx *ledger.Transaction, subBlockID string) *common.SubBlock {
	return &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: []*ledger.Transaction{tx},
		Timestamp:    time.Now(),
		NodeID:       gn.NodeID,
	}
}

// validateSubBlock validates a sub-block into a full block in the blockchain.
func (gn *GeospatialNode) validateSubBlock(subBlock *common.SubBlock) error {
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

// Geospatial Data Security and Compliance

// applyGeospatialSecurity ensures encryption and security protocols are applied to all geospatial data and transactions.
func (gn *GeospatialNode) applyGeospatialSecurity() error {
	err := gn.EncryptionService.ApplySecurity(gn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to apply encryption security for geospatial data: %v", err)
	}
	fmt.Printf("Encryption security applied successfully for geospatial node %s.\n", gn.NodeID)
	return nil
}

// ensureCompliance verifies that geospatial transactions comply with regulatory and legal frameworks.
func (gn *GeospatialNode) ensureCompliance(tx *ledger.Transaction) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	// Check compliance for the geospatial transaction.
	err := gn.GeoProcessor.CheckCompliance(tx)
	if err != nil {
		return fmt.Errorf("compliance check failed for transaction %s: %v", tx.TransactionID, err)
	}

	fmt.Printf("Transaction %s passed compliance check on geospatial node %s.\n", tx.TransactionID, gn.NodeID)
	return nil
}
