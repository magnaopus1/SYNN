package sidechains

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/consensus"
)

// NewSidechainDeployment initializes a new sidechain deployment manager
func NewSidechainDeployment(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus) *common.SidechainDeployment {
	return &common.SidechainDeployment{
		DeployedNetworks: make(map[string]*common.SidechainNetwork),
		Ledger:           ledgerInstance,
		Encryption:       encryptionService,
		Consensus:        consensus,
	}
}

// DeploySidechainNetwork deploys a new sidechain network
func (sd *common.SidechainDeployment) DeploySidechainNetwork(chainID string) (*common.SidechainNetwork, error) {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	// Check if the network already exists
	if _, exists := sd.DeployedNetworks[chainID]; exists {
		return nil, errors.New("sidechain network already deployed")
	}

	// Initialize the new sidechain network
	network := &common.SidechainNetwork{
		Nodes:      make(map[string]*common.SidechainNode),
		Blocks:     make(map[string]*common.Block),
		SubBlocks:  make(map[string]*common.SubBlock),
		Ledger:     sd.Ledger,
		Encryption: sd.Encryption,
		Consensus:  sd.Consensus,
	}

	// Record the deployment in the ledger
	err := sd.Ledger.RecordSidechainDeployment(chainID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log sidechain network deployment: %v", err)
	}

	sd.DeployedNetworks[chainID] = network
	fmt.Printf("Sidechain network %s deployed successfully\n", chainID)
	return network, nil
}

// AddNodeToSidechain adds a node to the sidechain's network
func (sd *common.SidechainDeployment) AddNodeToSidechain(chainID, nodeID, ipAddress string, nodeType common.NodeType) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	network, exists := sd.DeployedNetworks[chainID]
	if !exists {
		return fmt.Errorf("sidechain network %s not found", chainID)
	}

	// Add node to the sidechain network
	newNode := &common.SidechainNode{
		NodeID:    nodeID,
		IPAddress: ipAddress,
		NodeType:  nodeType,
	}

	network.Nodes[nodeID] = newNode

	// Log the node addition in the ledger
	err := network.Ledger.RecordNodeAddition(nodeID, ipAddress, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node addition: %v", err)
	}

	fmt.Printf("Node %s added to sidechain network %s\n", nodeID, chainID)
	return nil
}

// ValidateSidechainNetwork validates the entire sidechain network
func (sd *common.SidechainDeployment) ValidateSidechainNetwork(chainID string) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	network, exists := sd.DeployedNetworks[chainID]
	if !exists {
		return fmt.Errorf("sidechain network %s not found", chainID)
	}

	// Validate all blocks using Synnergy Consensus
	for _, block := range network.Blocks {
		err := sd.Consensus.ValidateBlock(block.BlockID)
		if err != nil {
			return fmt.Errorf("validation failed for block %s in sidechain network %s: %v", block.BlockID, chainID, err)
		}
	}

	// Log validation in the ledger
	err := sd.Ledger.RecordSidechainValidation(chainID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sidechain validation: %v", err)
	}

	fmt.Printf("Sidechain network %s validated successfully\n", chainID)
	return nil
}

// SecureSidechainNetwork encrypts sidechain data for security
func (sd *common.SidechainDeployment) SecureSidechainNetwork(chainID string) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	network, exists := sd.DeployedNetworks[chainID]
	if !exists {
		return fmt.Errorf("sidechain network %s not found", chainID)
	}

	// Encrypt all transaction data in the sidechain
	for _, block := range network.Blocks {
		for _, subBlock := range block.SubBlocks {
			for _, tx := range subBlock.Transactions {
				encryptedTxID, err := sd.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
				if err != nil {
					return fmt.Errorf("failed to encrypt transaction in block %s: %v", block.BlockID, err)
				}
				tx.TxID = string(encryptedTxID)
			}
		}
	}

	// Log the encryption event in the ledger
	err := sd.Ledger.RecordSidechainSecurity(chainID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sidechain encryption: %v", err)
	}

	fmt.Printf("Sidechain network %s secured successfully\n", chainID)
	return nil
}

// TerminateSidechainNetwork terminates a sidechain network and removes it
func (sd *common.SidechainDeployment) TerminateSidechainNetwork(chainID string) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	_, exists := sd.DeployedNetworks[chainID]
	if !exists {
		return fmt.Errorf("sidechain network %s not found", chainID)
	}

	// Remove the sidechain network
	delete(sd.DeployedNetworks, chainID)

	// Log the termination in the ledger
	err := sd.Ledger.RecordSidechainTermination(chainID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sidechain network termination: %v", err)
	}

	fmt.Printf("Sidechain network %s terminated successfully\n", chainID)
	return nil
}

// RetrieveSidechainNetwork retrieves a deployed sidechain network by ID
func (sd *common.SidechainDeployment) RetrieveSidechainNetwork(chainID string) (*common.SidechainNetwork, error) {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	network, exists := sd.DeployedNetworks[chainID]
	if !exists {
		return nil, fmt.Errorf("sidechain network %s not found", chainID)
	}

	fmt.Printf("Retrieved sidechain network %s\n", chainID)
	return network, nil
}
