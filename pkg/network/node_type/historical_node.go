package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"   // Shared components for encryption, consensus, and more
	"synnergy_network/pkg/ledger"   // Blockchain and ledger-related components
	"synnergy_network/pkg/network"  // Network and communication management
	"synnergy_network/pkg/synnergy_vm" // Synnergy Virtual Machine for contract execution
)

// HistoricalNode represents a node that stores and manages the entire blockchain history for audit and compliance purposes.
type HistoricalNode struct {
	NodeID            string                      // Unique identifier for the node
	Blockchain        *ledger.Blockchain          // Local copy of the blockchain ledger
	ConsensusEngine   *common.SynnergyConsensus   // Consensus engine for validating historical data
	EncryptionService *common.Encryption          // Encryption service for securing data
	NetworkManager    *network.NetworkManager     // Network manager for syncing with other nodes
	SNVM              *synnergy_vm.VirtualMachine // Virtual Machine for executing smart contracts
	BackupManager     *common.BackupManager       // Manages backups and data restoration
	mutex             sync.Mutex                  // Mutex for thread-safe operations
	SyncInterval      time.Duration               // Interval for syncing with other historical nodes
	BackupLocations   []string                    // Multiple backup locations (on-site, off-site, cloud)
}

// NewHistoricalNode initializes a new historical node responsible for archiving the blockchain’s entire history.
func NewHistoricalNode(nodeID string, blockchain *ledger.Blockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *network.NetworkManager, snvm *synnergy_vm.VirtualMachine, backupManager *common.BackupManager, syncInterval time.Duration, backupLocations []string) *HistoricalNode {
	return &HistoricalNode{
		NodeID:            nodeID,
		Blockchain:        blockchain,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SNVM:              snvm,
		BackupManager:     backupManager,
		SyncInterval:      syncInterval,
		BackupLocations:   backupLocations,
	}
}

// StartNode starts the historical node's operations, including syncing historical data, managing backups, and validating historical records.
func (hn *HistoricalNode) StartNode() error {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Begin syncing historical data and validating blocks.
	go hn.syncHistoricalData()

	// Backup management
	go hn.manageBackups()

	fmt.Printf("Historical Node %s started successfully.\n", hn.NodeID)
	return nil
}

// syncHistoricalData handles syncing the full blockchain history with other historical nodes and ensuring data integrity.
func (hn *HistoricalNode) syncHistoricalData() {
	ticker := time.NewTicker(hn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		hn.mutex.Lock()
		err := hn.retrieveAndValidateHistoricalData()
		if err != nil {
			fmt.Printf("Error syncing historical data: %v\n", err)
		}
		hn.mutex.Unlock()
	}
}

// retrieveAndValidateHistoricalData retrieves blockchain data from other historical nodes and validates it for integrity.
func (hn *HistoricalNode) retrieveAndValidateHistoricalData() error {
	// Retrieve blockchain history from other nodes.
	historicalData, err := hn.NetworkManager.FetchHistoricalData()
	if err != nil {
		return fmt.Errorf("failed to retrieve historical data: %v", err)
	}

	// Validate the historical data using the consensus engine.
	valid, err := hn.ConsensusEngine.ValidateHistoricalData(historicalData)
	if err != nil || !valid {
		return fmt.Errorf("invalid historical data: %v", err)
	}

	// Encrypt and store the validated historical data.
	encryptedData, err := hn.EncryptionService.EncryptData(historicalData, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt historical data: %v", err)
	}

	err = hn.Blockchain.RecordHistoricalData(encryptedData)
	if err != nil {
		return fmt.Errorf("failed to record historical data: %v", err)
	}

	fmt.Printf("Historical data successfully recorded by node %s.\n", hn.NodeID)
	return nil
}

// manageBackups ensures regular backups of historical data to multiple locations for redundancy and disaster recovery.
func (hn *HistoricalNode) manageBackups() {
	for {
		// Backup blockchain data at scheduled intervals.
		err := hn.BackupManager.ScheduleBackup(hn.BackupLocations)
		if err != nil {
			fmt.Printf("Backup failed for node %s: %v\n", hn.NodeID, err)
		}
		time.Sleep(hn.SyncInterval)
	}
}

// fetchAndRestoreBackup handles the restoration of historical data in case of node failure or data loss.
func (hn *HistoricalNode) fetchAndRestoreBackup(backupLocation string) error {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Retrieve backup data from the specified location.
	backupData, err := hn.BackupManager.FetchBackup(backupLocation)
	if err != nil {
		return fmt.Errorf("failed to fetch backup from %s: %v", backupLocation, err)
	}

	// Decrypt the backup data before restoring.
	decryptedData, err := hn.EncryptionService.DecryptData(backupData, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt backup data: %v", err)
	}

	// Restore the blockchain data from the decrypted backup.
	err = hn.Blockchain.RestoreData(decryptedData)
	if err != nil {
		return fmt.Errorf("failed to restore backup data: %v", err)
	}

	fmt.Printf("Backup from %s restored successfully for node %s.\n", backupLocation, hn.NodeID)
	return nil
}

// processHistoricalTransaction processes and validates a transaction, adding it to the historical ledger.
func (hn *HistoricalNode) processHistoricalTransaction(tx *ledger.Transaction) error {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Validate the transaction using the consensus engine.
	if valid, err := hn.ConsensusEngine.ValidateTransaction(tx); err != nil || !valid {
		return fmt.Errorf("invalid transaction: %v", err)
	}

	// Add the transaction to a sub-block for validation.
	subBlockID := common.GenerateSubBlockID()
	subBlock := hn.createSubBlock(tx, subBlockID)

	// Validate and store the sub-block.
	hn.Blockchain.AddSubBlock(subBlock)
	hn.validateSubBlock(subBlock)

	fmt.Printf("Transaction %s processed successfully by node %s.\n", tx.TransactionID, hn.NodeID)
	return nil
}

// createSubBlock creates a sub-block from a validated transaction.
func (hn *HistoricalNode) createSubBlock(tx *ledger.Transaction, subBlockID string) *common.SubBlock {
	return &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: []*ledger.Transaction{tx},
		Timestamp:    time.Now(),
		NodeID:       hn.NodeID,
	}
}

// validateSubBlock validates a sub-block into a full block in the blockchain.
func (hn *HistoricalNode) validateSubBlock(subBlock *common.SubBlock) error {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Validate the sub-block using the consensus mechanism.
	if hn.ConsensusEngine.ValidateSubBlock(subBlock) {
		block := hn.Blockchain.AddSubBlock(subBlock)
		hn.NetworkManager.BroadcastNewBlock(block)
		fmt.Printf("Sub-block %s validated and added to blockchain.\n", subBlock.SubBlockID)
		return nil
	}

	return errors.New("failed to validate sub-block")
}

// Virtual Machine (VM) Smart Contract Functions

// DeployContract deploys a smart contract to the historical node’s virtual machine.
func (hn *HistoricalNode) DeployContract(contractCode []byte, contractOwner string) (string, error) {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Encrypt the contract code before deployment.
	encryptedCode, err := hn.EncryptionService.EncryptData(contractCode, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt contract code: %v", err)
	}

	// Generate a unique contract ID.
	contractID := common.GenerateUniqueID()

	// Deploy the contract on the virtual machine.
	err = hn.SNVM.DeployContract(contractID, encryptedCode, contractOwner)
	if err != nil {
		return "", fmt.Errorf("failed to deploy contract: %v", err)
	}

	// Record the contract deployment in the ledger.
	err = hn.Blockchain.RecordContractDeployment(contractID, contractOwner, encryptedCode)
	if err != nil {
		return "", fmt.Errorf("failed to record contract deployment in ledger: %v", err)
	}

	fmt.Printf("Contract %s deployed successfully on historical node %s by owner %s.\n", contractID, hn.NodeID, contractOwner)
	return contractID, nil
}

// ExecuteContract executes a smart contract on the historical node’s virtual machine.
func (hn *HistoricalNode) ExecuteContract(contractID string, args []byte) ([]byte, error) {
	hn.mutex.Lock()
	defer hn.mutex.Unlock()

	// Encrypt the contract arguments before execution.
	encryptedArgs, err := hn.EncryptionService.EncryptData(args, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt contract arguments: %v", err)
	}

	// Execute the contract on the virtual machine.
	result, err := hn.SNVM.ExecuteContract(contractID, encryptedArgs)
	if err != nil {
		return nil, fmt.Errorf("contract execution failed: %v", err)
	}

	// Record the execution result in the ledger.
	err = hn.Blockchain.RecordContractExecution(contractID, result)
	if err != nil {
		return nil, fmt.Errorf("failed to record contract execution in ledger: %v", err)
	}

	fmt.Printf("Contract %s executed successfully on historical node %s.\n", contractID, hn.NodeID)
	return result, nil
}
