package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"    // Shared components like encryption, consensus, and sub-blocks
	"synnergy_network/pkg/ledger"    // Blockchain and ledger-related components
	"synnergy_network/pkg/cross_chain" // Cross-chain bridge and interaction components
)

// CrossChainNode represents a node responsible for cross-chain interactions.
type CrossChainNode struct {
	NodeID            string                        // Unique identifier for the node
	Blockchain        *ledger.Blockchain            // Local copy of the blockchain ledger
	ConsensusEngine   *common.SynnergyConsensus     // Consensus engine for validating transactions
	EncryptionService *common.Encryption            // Encryption service for secure communication
	NetworkManager    *common.NetworkManager        // Network manager for communicating with other nodes
	SNVM              *common.VirtualMachine        // Virtual Machine for contract execution
	CrossChainManager *cross_chain.Manager          // Cross-chain manager for bridge and cross-chain operations
	Bridges           map[string]*cross_chain.Bridge // Active cross-chain bridges for other blockchains
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SyncInterval      time.Duration                 // Interval for syncing with other full nodes
}

// NewCrossChainNode initializes a new cross-chain node in the Synnergy Network.
func NewCrossChainNode(nodeID string, blockchain *ledger.Blockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *common.NetworkManager, crossChainManager *cross_chain.Manager, syncInterval time.Duration) *CrossChainNode {
	return &CrossChainNode{
		NodeID:            nodeID,
		Blockchain:        blockchain,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SNVM:              common.NewVirtualMachine(), // Initialize the virtual machine for contract execution
		CrossChainManager: crossChainManager,
		Bridges:           make(map[string]*cross_chain.Bridge),
		SyncInterval:      syncInterval,
	}
}

// StartNode begins the cross-chain node's operations, syncing with other chains, managing bridges, and processing cross-chain transactions.
func (ccn *CrossChainNode) StartNode() error {
	ccn.mutex.Lock()
	defer ccn.mutex.Unlock()

	// Begin syncing with other nodes.
	go ccn.syncWithOtherNodes()

	// Listen for incoming cross-chain transactions.
	go ccn.listenForTransactions()

	// Set up cross-chain bridges.
	go ccn.setupCrossChainBridges()

	fmt.Printf("Cross-chain node %s started successfully.\n", ccn.NodeID)
	return nil
}

// syncWithOtherNodes syncs the blockchain with other nodes at regular intervals.
func (ccn *CrossChainNode) syncWithOtherNodes() {
	ticker := time.NewTicker(ccn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		ccn.mutex.Lock()
		otherNodes := ccn.NetworkManager.DiscoverOtherNodes(ccn.NodeID)
		for _, node := range otherNodes {
			ccn.syncBlockchainFromNode(node)
		}
		ccn.mutex.Unlock()
	}
}

// syncBlockchainFromNode syncs the blockchain from a peer node to ensure it is up to date.
func (ccn *CrossChainNode) syncBlockchainFromNode(peerNode string) {
	peerBlockchain, err := ccn.NetworkManager.RequestBlockchain(peerNode)
	if err != nil {
		fmt.Printf("Failed to sync blockchain from node %s: %v\n", peerNode, err)
		return
	}

	// Validate the blockchain and update the local copy if necessary.
	if ccn.ConsensusEngine.ValidateBlockchain(peerBlockchain) {
		ccn.Blockchain = peerBlockchain
		fmt.Printf("Blockchain synced successfully from node %s.\n", peerNode)
	} else {
		fmt.Printf("Blockchain sync from node %s failed validation.\n", peerNode)
	}
}

// listenForTransactions listens for incoming cross-chain transactions and processes them.
func (ccn *CrossChainNode) listenForTransactions() {
	for {
		transaction, err := ccn.NetworkManager.ReceiveTransaction()
		if err != nil {
			fmt.Printf("Error receiving transaction: %v\n", err)
			continue
		}

		// Validate and process the cross-chain transaction.
		err = ccn.processTransaction(transaction)
		if err != nil {
			fmt.Printf("Transaction processing failed: %v\n", err)
		}
	}
}

// processTransaction processes and validates incoming cross-chain transactions.
func (ccn *CrossChainNode) processTransaction(tx *ledger.Transaction) error {
	ccn.mutex.Lock()
	defer ccn.mutex.Unlock()

	// Validate the transaction using the consensus engine.
	if valid, err := ccn.ConsensusEngine.ValidateTransaction(tx); err != nil || !valid {
		return fmt.Errorf("invalid transaction: %v", err)
	}

	// Handle cross-chain transactions.
	if tx.IsCrossChain {
		err := ccn.processCrossChainTransaction(tx)
		if err != nil {
			return fmt.Errorf("cross-chain transaction failed: %v", err)
		}
	}

	// Add the transaction to a sub-block and validate it.
	subBlockID := common.GenerateSubBlockID()
	subBlock := ccn.createSubBlock(tx, subBlockID)

	ccn.SubBlocks[subBlockID] = subBlock
	ccn.tryValidateSubBlock(subBlock)

	fmt.Printf("Transaction %s processed successfully.\n", tx.TransactionID)
	return nil
}

// processCrossChainTransaction handles transactions that involve cross-chain transfers.
func (ccn *CrossChainNode) processCrossChainTransaction(tx *ledger.Transaction) error {
	bridge, exists := ccn.Bridges[tx.CrossChainID]
	if !exists {
		return errors.New("cross-chain bridge not found")
	}

	// Execute cross-chain logic and transfer assets/data.
	err := bridge.Transfer(tx)
	if err != nil {
		return fmt.Errorf("cross-chain transfer failed: %v", err)
	}

	// Record the transaction in the ledger.
	err = ccn.Blockchain.RecordCrossChainTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to record cross-chain transaction: %v", err)
	}

	fmt.Printf("Cross-chain transaction %s processed successfully through bridge %s.\n", tx.TransactionID, bridge.BridgeID)
	return nil
}

// createSubBlock creates a sub-block from the validated transaction.
func (ccn *CrossChainNode) createSubBlock(tx *ledger.Transaction, subBlockID string) *common.SubBlock {
	return &common.SubBlock{
		SubBlockID:   subBlockID,
		Transactions: []*ledger.Transaction{tx},
		Timestamp:    time.Now(),
		NodeID:       ccn.NodeID,
	}
}

// tryValidateSubBlock attempts to validate a sub-block into a full block.
func (ccn *CrossChainNode) tryValidateSubBlock(subBlock *common.SubBlock) error {
	ccn.mutex.Lock()
	defer ccn.mutex.Unlock()

	// Validate the sub-block using the consensus mechanism.
	if ccn.ConsensusEngine.ValidateSubBlock(subBlock) {
		// Add the sub-block to the blockchain and broadcast the new block.
		block := ccn.Blockchain.AddSubBlock(subBlock)
		err := ccn.NetworkManager.BroadcastNewBlock(block)
		if err != nil {
			return fmt.Errorf("failed to broadcast new block: %v", err)
		}

		fmt.Printf("Sub-block %s validated and added to blockchain.\n", subBlock.SubBlockID)
		return nil
	}
	return errors.New("failed to validate sub-block")
}

// setupCrossChainBridges initializes cross-chain bridges to other blockchains.
func (ccn *CrossChainNode) setupCrossChainBridges() {
	ccn.mutex.Lock()
	defer ccn.mutex.Unlock()

	bridges := ccn.CrossChainManager.GetAvailableBridges()
	for _, bridge := range bridges {
		ccn.Bridges[bridge.BridgeID] = bridge
		fmt.Printf("Cross-chain bridge %s set up for interaction.\n", bridge.BridgeID)
	}
}

// Cross-Chain Contract Deployment and Execution

// DeployCrossChainContract deploys a smart contract across chains using the node’s virtual machine.
func (ccn *CrossChainNode) DeployCrossChainContract(contractCode []byte, contractOwner string) (string, error) {
	ccn.mutex.Lock()
	defer ccn.mutex.Unlock()

	// Encrypt the contract code before deployment.
	encryptedCode, err := ccn.EncryptionService.EncryptData(contractCode, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt contract code: %v", err)
	}

	// Generate a unique contract ID.
	contractID := common.GenerateUniqueID()

	// Deploy the contract using the cross-chain virtual machine.
	err = ccn.SNVM.DeployContract(contractID, encryptedCode, contractOwner)
	if err != nil {
		return "", fmt.Errorf("failed to deploy cross-chain contract: %v", err)
	}

	// Record the deployment in the ledger.
	err = ccn.Blockchain.RecordContractDeployment(contractID, contractOwner, encryptedCode)
	if err != nil {
		return "", fmt.Errorf("failed to record contract deployment in ledger: %v", err)
	}

	fmt.Printf("Cross-chain contract %s deployed successfully on node %s.\n", contractID, ccn.NodeID)
	return contractID, nil
}

// ExecuteCrossChainContract executes a smart contract that spans across different chains using the cross-chain node’s virtual machine.
func (ccn *CrossChainNode) ExecuteCrossChainContract(contractID string, args []byte) ([]byte, error) {
	ccn.mutex.Lock()
	defer ccn.mutex.Unlock()

	// Retrieve the contract code from the ledger.
	contractCode, err := ccn.Blockchain.GetContractCode(contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve contract code: %v", err)
	}

	// Decrypt the contract code.
	decryptedCode, err := ccn.EncryptionService.DecryptData(contractCode, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt contract code: %v", err)
	}

	// Execute the contract using the virtual machine.
	result, err := ccn.SNVM.ExecuteContract(decryptedCode, args)
	if err != nil {
		return nil, fmt.Errorf("cross-chain contract execution failed: %v", err)
	}

	// Record the execution result in the ledger.
	err = ccn.Blockchain.RecordContractExecution(contractID, result)
	if err != nil {
		return nil, fmt.Errorf("failed to record cross-chain contract execution in ledger: %v", err)
	}

	fmt.Printf("Cross-chain contract %s executed successfully on node %s.\n", contractID, ccn.NodeID)
	return result, nil
}

// Cross-Chain Asset and Data Transfer

// TransferAssetsAcrossChains transfers assets between two chains through cross-chain bridges.
func (ccn *CrossChainNode) TransferAssetsAcrossChains(assetID string, destinationChain string, amount uint64) error {
	ccn.mutex.Lock()
	defer ccn.mutex.Unlock()

	bridge, exists := ccn.Bridges[destinationChain]
	if !exists {
		return fmt.Errorf("no bridge available for chain %s", destinationChain)
	}

	// Perform the asset transfer through the cross-chain bridge.
	err := bridge.TransferAssets(assetID, amount)
	if err != nil {
		return fmt.Errorf("asset transfer across chains failed: %v", err)
	}

	fmt.Printf("Asset %s transferred successfully to chain %s via bridge.\n", assetID, destinationChain)
	return nil
}

// TransferDataAcrossChains transfers data securely between two blockchain networks.
func (ccn *CrossChainNode) TransferDataAcrossChains(data []byte, destinationChain string) error {
	ccn.mutex.Lock()
	defer ccn.mutex.Unlock()

	bridge, exists := ccn.Bridges[destinationChain]
	if !exists {
		return fmt.Errorf("no bridge available for chain %s", destinationChain)
	}

	// Encrypt the data before transfer.
	encryptedData, err := ccn.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %v", err)
	}

	// Transfer the encrypted data through the cross-chain bridge.
	err = bridge.TransferData(encryptedData)
	if err != nil {
		return fmt.Errorf("data transfer across chains failed: %v", err)
	}

	fmt.Printf("Data transferred successfully to chain %s via bridge.\n", destinationChain)
	return nil
}

// Security and Encryption

// ApplySecurityProtocols ensures all communications and data transfers are encrypted.
func (ccn *CrossChainNode) ApplySecurityProtocols() error {
	err := ccn.EncryptionService.ApplySecurity(ccn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to apply security protocols: %v", err)
	}

	fmt.Printf("Security protocols applied successfully for cross-chain node %s.\n", ccn.NodeID)
	return nil
}

// HandleCrossChainFailures

// HandleBridgeFailure handles failures in cross-chain bridges, ensuring continuity of operations.
func (ccn *CrossChainNode) HandleBridgeFailure(bridgeID string) error {
	ccn.mutex.Lock()
	defer ccn.mutex.Unlock()

	bridge, exists := ccn.Bridges[bridgeID]
	if !exists {
		return fmt.Errorf("bridge %s not found", bridgeID)
	}

	// Attempt to recover or reroute the bridge operations.
	err := bridge.RecoverFromFailure()
	if err != nil {
		return fmt.Errorf("failed to recover from bridge failure: %v", err)
	}

	fmt.Printf("Bridge %s recovered from failure successfully.\n", bridgeID)
	return nil
}

// Monitoring and Auditing

// MonitorCrossChainTransactions monitors all cross-chain transactions to ensure accuracy and security.
func (ccn *CrossChainNode) MonitorCrossChainTransactions() error {
	ccn.mutex.Lock()
	defer ccn.mutex.Unlock()

	transactions, err := ccn.CrossChainManager.GetPendingCrossChainTransactions()
	if err != nil {
		return fmt.Errorf("failed to retrieve cross-chain transactions: %v", err)
	}

	for _, tx := range transactions {
		// Validate each transaction using the consensus engine.
		if valid, err := ccn.ConsensusEngine.ValidateTransaction(tx); err != nil || !valid {
			fmt.Printf("Cross-chain transaction %s failed validation.\n", tx.TransactionID)
		} else {
			fmt.Printf("Cross-chain transaction %s validated successfully.\n", tx.TransactionID)
		}
	}

	return nil
}

// PerformCrossChainAudit performs a regular audit of all cross-chain operations to ensure transparency and integrity.
func (ccn *CrossChainNode) PerformCrossChainAudit() error {
	ccn.mutex.Lock()
	defer ccn.mutex.Unlock()

	auditResults, err := ccn.CrossChainManager.PerformAudit()
	if err != nil {
		return fmt.Errorf("cross-chain audit failed: %v", err)
	}

	// Log the audit results.
	for _, result := range auditResults {
		fmt.Printf("Cross-chain audit result: %s - Status: %s\n", result.TransactionID, result.Status)
	}

	return nil
}

