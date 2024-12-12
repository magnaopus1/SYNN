package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"      // Shared components for encryption, consensus, etc.
	"synnergy_network/pkg/ledger"      // Blockchain and ledger-related components
	"synnergy_network/pkg/network"     // Network and communication management
	"synnergy_network/pkg/synnergy_vm" // Synnergy Virtual Machine for smart contract execution
	"synnergy_network/pkg/api"         // API integration framework
	"synnergy_network/pkg/oracle"      // Smart contract oracles for external data
	"synnergy_network/pkg/interoperability" // Cross-chain communication & chain adapters
)

// IntegrationNode represents a node responsible for integrating external systems, APIs, and other blockchains.
type IntegrationNode struct {
	NodeID            string                         // Unique identifier for the node
	Blockchain        *ledger.Blockchain             // Local copy of the blockchain ledger
	ConsensusEngine   *common.SynnergyConsensus      // Consensus engine for validating transactions and sub-blocks
	EncryptionService *common.Encryption             // Encryption service for secure communication and data handling
	NetworkManager    *network.NetworkManager        // Network manager for communication with other nodes and external systems
	SNVM              *common.VMInterface   // The Synnergy Network Virtual Machine
	ChainAdapters     map[string]*interoperability.ChainAdapter // Chain adapters for cross-blockchain operations
	APIGateways       map[string]*api.APIGateway     // API gateways for external system integrations
	Oracles           map[string]*oracle.Oracle      // Oracles for smart contract external data handling
	mutex             sync.Mutex                     // Mutex for thread-safe operations
	SyncInterval      time.Duration                  // Interval for syncing with other nodes
	NewFeatures       map[string]interface{}         // Safe integration of new features
}

// NewIntegrationNode initializes a new integration node in the blockchain network.
func NewIntegrationNode(nodeID string, blockchain *ledger.Blockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *network.NetworkManager, syncInterval time.Duration) *IntegrationNode {
	return &IntegrationNode{
		NodeID:            nodeID,
		Blockchain:        blockchain,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SNVM:              synnergy_vm.NewVirtualMachine(), // Initialize the virtual machine for smart contract execution
		ChainAdapters:     make(map[string]*interoperability.ChainAdapter), // Cross-chain adapters for other blockchains
		APIGateways:       make(map[string]*api.APIGateway), // API integration gateways
		Oracles:           make(map[string]*oracle.Oracle),  // Oracle for smart contract external data
		SyncInterval:      syncInterval,
		NewFeatures:       make(map[string]interface{}),     // Safe handling of new features and functions
	}
}

// StartNode starts the integration node's operations, including syncing, cross-chain handling, API management, and contract execution.
func (in *IntegrationNode) StartNode() error {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Begin syncing the blockchain with other nodes.
	go in.syncWithOtherNodes()

	// Set up API gateways and external connections.
	go in.setupAPIGateways()

	// Set up chain adapters for cross-blockchain operations.
	go in.setupChainAdapters()

	// Set up oracles for smart contract external data.
	go in.setupOracles()

	fmt.Printf("Integration Node %s started successfully.\n", in.NodeID)
	return nil
}

// syncWithOtherNodes handles syncing the blockchain with other full nodes at regular intervals.
func (in *IntegrationNode) syncWithOtherNodes() {
	ticker := time.NewTicker(in.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		in.mutex.Lock()
		// Discover other nodes and sync blockchain data.
		otherNodes := in.NetworkManager.DiscoverOtherNodes(in.NodeID)
		for _, node := range otherNodes {
			err := in.syncBlockchainFromNode(node)
			if err != nil {
				fmt.Printf("Failed to sync blockchain from node %s: %v\n", node, err)
			}
		}
		in.mutex.Unlock()
	}
}

// syncBlockchainFromNode syncs the blockchain state from a peer node.
func (in *IntegrationNode) syncBlockchainFromNode(peerNode string) error {
	peerBlockchain, err := in.NetworkManager.RequestBlockchain(peerNode)
	if err != nil {
		return fmt.Errorf("failed to request blockchain from node %s: %v", peerNode, err)
	}

	// Validate the blockchain and update the local copy if necessary.
	if in.ConsensusEngine.ValidateBlockchain(peerBlockchain) {
		in.Blockchain = peerBlockchain
		fmt.Printf("Blockchain synced successfully from node %s.\n", peerNode)
	} else {
		fmt.Printf("Blockchain sync from node %s failed validation.\n", peerNode)
	}
	return nil
}

// setupAPIGateways sets up API gateways for external system connections.
func (in *IntegrationNode) setupAPIGateways() {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Example API setup: Connect to financial services, IoT, or external data providers
	externalAPIs := []string{"FinanceAPI", "IoTAPI", "WeatherAPI"}
	for _, apiName := range externalAPIs {
		apiGateway := api.NewAPIGateway(apiName)
		in.APIGateways[apiName] = apiGateway
		fmt.Printf("API Gateway %s set up.\n", apiName)
	}
}

// setupChainAdapters sets up chain adapters for cross-blockchain integration.
func (in *IntegrationNode) setupChainAdapters() {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Example chain adapters: Connect to Ethereum, Bitcoin, and other networks
	externalChains := []string{"Ethereum", "Bitcoin", "Polkadot"}
	for _, chainName := range externalChains {
		chainAdapter := interoperability.NewChainAdapter(chainName)
		in.ChainAdapters[chainName] = chainAdapter
		fmt.Printf("Chain Adapter for %s set up.\n", chainName)
	}
}

// setupOracles sets up oracles to integrate external data for smart contracts.
func (in *IntegrationNode) setupOracles() {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Example oracles: Financial data, weather data, etc.
	externalDataSources := []string{"StockPriceOracle", "WeatherOracle"}
	for _, oracleName := range externalDataSources {
		oracle := oracle.NewOracle(oracleName)
		in.Oracles[oracleName] = oracle
		fmt.Printf("Oracle %s set up.\n", oracleName)
	}
}

// processExternalAPICall handles external API data, validates, and integrates into blockchain.
func (in *IntegrationNode) processExternalAPICall(apiName string, data []byte) error {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Validate and process the data from the external API
	apiGateway, exists := in.APIGateways[apiName]
	if !exists {
		return errors.New("API Gateway not found")
	}

	validatedData, err := apiGateway.ValidateData(data)
	if err != nil {
		return fmt.Errorf("validation failed for API %s: %v", apiName, err)
	}

	// Encrypt the data before it enters the blockchain.
	encryptedData, err := in.EncryptionService.EncryptData(validatedData, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}

	// Add to blockchain transaction.
	tx := ledger.NewTransaction(encryptedData)
	err = in.Blockchain.AddTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to add transaction: %v", err)
	}

	fmt.Printf("External data from API %s processed and added to blockchain.\n", apiName)
	return nil
}

// processCrossChainTransaction handles cross-blockchain transactions through adapters.
func (in *IntegrationNode) processCrossChainTransaction(chainName string, tx *ledger.Transaction) error {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	chainAdapter, exists := in.ChainAdapters[chainName]
	if !exists {
		return fmt.Errorf("Chain Adapter not found for chain %s", chainName)
	}

	// Use the chain adapter to handle cross-chain transaction
	err := chainAdapter.ProcessTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to process cross-chain transaction on %s: %v", chainName, err)
	}

	fmt.Printf("Cross-chain transaction processed on chain %s.\n", chainName)
	return nil
}

// DeployContract deploys a smart contract to the integration node's virtual machine.
func (in *IntegrationNode) DeployContract(contractCode []byte, contractOwner string) (string, error) {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Encrypt the contract code before deployment.
	encryptedCode, err := in.EncryptionService.EncryptData(contractCode, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt contract code: %v", err)
	}

	// Generate a unique contract ID.
	contractID := common.GenerateUniqueID()

	// Deploy the contract on the virtual machine.
	err = in.SNVM.DeployContract(contractID, encryptedCode, contractOwner)
	if err != nil {
		return "", fmt.Errorf("failed to deploy contract: %v", err)
	}

	// Record the contract deployment in the blockchain.
	err = in.Blockchain.RecordContractDeployment(contractID, contractOwner, encryptedCode)
	if err != nil {
		return "", fmt.Errorf("failed to record contract deployment: %v", err)
	}

	fmt.Printf("Contract %s deployed successfully on integration node %s.\n", contractID, in.NodeID)
	return contractID, nil
}

// ExecuteContract executes a smart contract on the integration node’s virtual machine.
func (in *IntegrationNode) ExecuteContract(contractID string, args []byte) ([]byte, error) {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Encrypt the contract arguments.
	encryptedArgs, err := in.EncryptionService.EncryptData(args, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt contract arguments: %v", err)
	}

	// Execute the contract on the virtual machine.
	result, err := in.SNVM.ExecuteContract(contractID, encryptedArgs)
	if err != nil {
		return nil, fmt.Errorf("contract execution failed: %v", err)
	}

	// Record the execution result in the blockchain.
	err = in.Blockchain.RecordContractExecution(contractID, result)
	if err != nil {
		return nil, fmt.Errorf("failed to record contract execution: %v", err)
	}

	fmt.Printf("Contract %s executed successfully on integration node %s.\n", contractID, in.NodeID)
	return result, nil
}

// AddNewFeature safely integrates a new feature into the node.
func (in *IntegrationNode) AddNewFeature(featureName string, feature interface{}) error {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	// Add and activate the new feature.
	in.NewFeatures[featureName] = feature
	fmt.Printf("New feature %s integrated successfully on integration node %s.\n", featureName, in.NodeID)
	return nil
}
